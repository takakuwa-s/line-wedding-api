package usecase

import (
	"fmt"
	"math"
	"time"

	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/dto"
	"github.com/takakuwa-s/line-wedding-api/entity"
	"github.com/takakuwa-s/line-wedding-api/usecase/igateway"
	"go.uber.org/zap"
)

type BackgroundProcessUsecase struct {
	lg igateway.ILineGateway
	fg igateway.IFaceGateway
	fr igateway.IFileRepository
	br igateway.IBinaryRepository
}

func NewBackgroundProcessUsecase(lg igateway.ILineGateway, fg igateway.IFaceGateway, fr igateway.IFileRepository, br igateway.IBinaryRepository) *BackgroundProcessUsecase {
	return &BackgroundProcessUsecase{lg: lg, fg: fg, fr: fr, br: br}
}

func (bpu *BackgroundProcessUsecase) RecoverFileProcessing() {
	conf.Log.Info("[BATCH] Start the recovery process")
	statuses := []entity.FileStatus{entity.New, entity.Uploaded}
	files, err := bpu.fr.FindByFileStatusIn(statuses)
	if err != nil {
		conf.Log.Error("[BATCH] failed to get the file metadata for uploading", zap.String("error", err.Error()))
		return
	}
	conf.Log.Info("[BATCH] Start uploading file binary", zap.Int("file count", len(files)))
	bpu.uploadFiles(files, "[BATCH]")

	statuses = []entity.FileStatus{entity.Deleted}
	files, err = bpu.fr.FindByFileStatusIn(statuses)
	if err != nil {
		conf.Log.Error("[BATCH] failed to get the file metadata for deleting", zap.String("error", err.Error()))
		return
	}
	conf.Log.Info("[BATCH] Start deleting file binary", zap.Int("file count", len(files)))
	bpu.deleteFiles(files)
	conf.Log.Info("[BATCH] Complete the recovery process")
}

func (bpu *BackgroundProcessUsecase) UploadFilesByIds(ids []string) {
	conf.Log.Info("[API] Start uploading file binary", zap.Int("len", len(ids)), zap.Strings("ids", ids))
	files, err := bpu.fr.FindByIdsAndFileStatus(ids, entity.New)
	if err != nil {
		conf.Log.Error("[API] failed to get the file metadata", zap.String("error", err.Error()))
		return
	}
	bpu.uploadFiles(files, "[API]")
	conf.Log.Info("[API] Complete uploading file binary", zap.Int("len", len(ids)), zap.Strings("ids", ids))
}

func (bpu *BackgroundProcessUsecase) uploadFiles(files []entity.File, triggerName string) {
	for idx, f := range files {
		switch f.FileType {
		case entity.Image:
			if err := bpu.uploadImage(f, idx); err != nil {
				conf.Log.Error(triggerName+" Failed to upload image", zap.Int("idx", idx), zap.String("id", f.Id), zap.Error(err))
			}
		case entity.Video:
			if err := bpu.uploadVideo(f, idx); err != nil {
				conf.Log.Error(triggerName+" Failed to upload video", zap.Int("idx", idx), zap.String("id", f.Id), zap.Error(err))
			}
		default:
			conf.Log.Error(triggerName+" Unknown fileType", zap.Int("idx", idx), zap.String("id", f.Id), zap.String("fileType", string(f.FileType)))
		}
	}
}

func (bpu *BackgroundProcessUsecase) uploadImage(f entity.File, i int) error {
	conf.Log.Info("Start uploading image", zap.Int("idx", i), zap.String("id", f.Id))

	// upload the file binary
	f1, err := bpu.uploadBinary(f, entity.Image)
	if err != nil {
		return err
	}
	f2, faceErr := bpu.updateFaceScore(*f1)

	// Save file data
	f2.UpdatedAt = time.Now()
	err = bpu.fr.SaveFile(f2)
	if err != nil {
		if faceErr != nil {
			return fmt.Errorf("failed to call face api and save image metadata; face err = %s ,file err = %w", faceErr, err)
		} else {
			return err
		}
	}
	if faceErr != nil {
		return faceErr
	}
	conf.Log.Info("Complete uploading image", zap.Int("idx", i), zap.String("id", f2.Id))
	return nil
}

func (bpu *BackgroundProcessUsecase) uploadVideo(f entity.File, i int) error {
	conf.Log.Info("Start uploading video", zap.Int("idx", i), zap.String("id", f.Id))

	// upload the file binary
	f1, err := bpu.uploadBinary(f, entity.Video)
	if err != nil {
		return err
	}

	// Save file data
	f1.UpdatedAt = time.Now()
	if err := bpu.fr.SaveFile(f1); err != nil {
		return err
	}
	conf.Log.Info("Complete uploading video", zap.Int("idx", i), zap.String("id", f1.Id))
	return nil
}

func (bpu *BackgroundProcessUsecase) uploadBinary(f entity.File, fileType entity.FileType) (*entity.File, error) {
	if f.FileStatus != entity.New {
		return &f, nil
	}
	// Get the file binary
	content, err := bpu.lg.GetFileContent(f.Id)
	if err != nil {
		return nil, err
	}
	defer content.Close()
	var file *entity.File
	switch f.FileType {
	case entity.Image:
		file, err = bpu.br.SaveImageBinary(f, content)
	case entity.Video:
		file, err = bpu.br.SaveVideoBinary(f, content)
	}
	if err != nil {
		return nil, err
	} else {
		return file, nil
	}
}

func (bpu *BackgroundProcessUsecase) updateFaceScore(f entity.File) (*entity.File, error) {
	if f.FileStatus != entity.Uploaded {
		return &f, nil
	}
	res, err := bpu.fg.GetFaceAnalysis(f.ContentUrl)
	if err != nil {
		return &f, err
	}
	return bpu.calculateFaceScore(res, f), nil
}

func (bpu *BackgroundProcessUsecase) calculateFaceScore(r []*dto.FaceResponse, f entity.File) *entity.File {
	faceCount := len(r)
	if faceCount <= 0 || faceCount > 10 {
		f.FaceCount = faceCount
		f.FaceHappinessLevel = 0
		f.FacePhotoBeauty = 0
		f.FaceScore = 0
		f.FileStatus = entity.Open
		return &f
	}
	var faceHappinessLevelSum float32
	var facePhotoBeautySum float32
	for _, f := range r {

		//---- calculate the face happiness level (max: 40) ---

		//---- Joy : faceHappinessLevelSum = 0 ~ 40
		if f.JoyLikelihood == 0 {
			facePhotoBeautySum -= 5
		} else {
			faceHappinessLevelSum += 10 * f.DetectionConfidence * float32(f.JoyLikelihood-1)
		}

		//---- Sorrow : faceHappinessLevelSum = -12 ~ 0
		if f.SorrowLikelihood == 0 {
			facePhotoBeautySum -= 5
		} else {
			faceHappinessLevelSum -= 3 * f.DetectionConfidence * float32((f.SorrowLikelihood - 1))
		}

		//---- Anger : faceHappinessLevelSum = -16 ~ 0
		if f.AngerLikelihood == 0 {
			facePhotoBeautySum -= 5
		} else {
			faceHappinessLevelSum -= 4 * f.DetectionConfidence * float32((f.AngerLikelihood - 1))
		}

		//---- Surprise : faceHappinessLevelSum = -8 ~ 0
		if f.SurpriseLikelihood == 0 {
			facePhotoBeautySum -= 5
		} else {
			faceHappinessLevelSum -= 2 * f.DetectionConfidence * float32((f.SurpriseLikelihood - 1))
		}

		//---- calculate the face photo beauty  (max: 40) ---
		facePhotoBeautySum += 40

		//---- UnderExposedLikelihood : facePhotoBeautySum = -20 ~ 0
		if f.UnderExposedLikelihood == 0 {
			facePhotoBeautySum -= 5
		} else {
			facePhotoBeautySum -= 5 * f.DetectionConfidence * float32((f.UnderExposedLikelihood - 1))
		}

		//---- BlurredLikelihood : facePhotoBeautySum = -20 ~ 0
		if f.BlurredLikelihood == 0 {
			facePhotoBeautySum -= 5
		} else {
			facePhotoBeautySum -= 5 * f.DetectionConfidence * float32((f.BlurredLikelihood - 1))
		}

		//---- RollAngle : facePhotoBeautySum = -18 ~ 0
		facePhotoBeautySum -= 0.1 * f.DetectionConfidence * float32(math.Abs(float64(f.RollAngle)))

		//---- PanAngle : facePhotoBeautySum = -18 ~ 0
		facePhotoBeautySum -= 0.1 * f.DetectionConfidence * float32(math.Abs(float64(f.PanAngle)))

		//---- TiltAngle : facePhotoBeautySum = -18 ~ 0
		facePhotoBeautySum -= 0.1 * f.DetectionConfidence * float32(math.Abs(float64(f.TiltAngle)))
	}

	if faceHappinessLevelSum < 0 {
		faceHappinessLevelSum = 0
	}
	if facePhotoBeautySum < 0 {
		facePhotoBeautySum = 0
	}

	// calculate the face count bonus point (max: 20)
	bonusPoint := 2 * float32(faceCount)

	f.FaceCount = faceCount
	f.FaceHappinessLevel = faceHappinessLevelSum / float32(faceCount)
	f.FacePhotoBeauty = facePhotoBeautySum / float32(faceCount)
	f.FaceScore = f.FaceHappinessLevel + f.FacePhotoBeauty + bonusPoint
	f.FileStatus = entity.Open
	return &f
}

func (bpu *BackgroundProcessUsecase) DeleteFilesByIds(ids []string) {
	conf.Log.Info("[API] Start deleting file binary", zap.Int("len", len(ids)), zap.Strings("ids", ids))
	files, err := bpu.fr.FindByIds(ids)
	if err != nil {
		conf.Log.Error("[API] failed to get the file metadata", zap.String("error", err.Error()))
		return
	}
	bpu.deleteFiles(files)
	conf.Log.Info("[API] Complete deleting file binary", zap.Int("len", len(ids)), zap.Strings("ids", ids))
}

func (bpu *BackgroundProcessUsecase) deleteFiles(files []entity.File) {
	for _, f := range files {
		if err := bpu.br.DeleteBinary(f.Name, string(f.FileType)); err != nil {
			conf.Log.Error("failed to delete the file binary", zap.Any("file", f))
		}
		if err := bpu.fr.DeleteFileById(f.Id); err != nil {
			conf.Log.Error("failed to delete the file metadata", zap.Any("file", f))
		}
	}
}
