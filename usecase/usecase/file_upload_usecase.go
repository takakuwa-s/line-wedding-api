package usecase

import (
	"fmt"
	"math"

	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/dto"
	"github.com/takakuwa-s/line-wedding-api/entity"
	"github.com/takakuwa-s/line-wedding-api/usecase/igateway"
	"go.uber.org/zap"
)

type FileUploadUsecase struct {
	lg igateway.ILineGateway
	fg igateway.IFaceGateway
	fr igateway.IFileRepository
	br igateway.IBinaryRepository
}

func NewFileUploadUsecase(lg igateway.ILineGateway, fg igateway.IFaceGateway, fr igateway.IFileRepository, br igateway.IBinaryRepository) *FileUploadUsecase {
	return &FileUploadUsecase{lg: lg, fg: fg, fr: fr, br: br}
}

func (fuu *FileUploadUsecase) UploadFiles(ids []string) error {
	conf.Log.Info("Start uploading file binary", zap.Strings("ids", ids))
	files, err := fuu.fr.FindByIdsAndUploaded(ids, false)
	if err != nil {
		return err
	}
	for idx, f := range files {
		if err := fuu.uploadFile(f, idx); err != nil {
			conf.Log.Error("Failed to upload file", zap.Int("idx", idx), zap.String("id", f.Id), zap.Error(err))
		}
	}
	conf.Log.Info("Complete uploading file binary", zap.Strings("ids", ids))
	return nil
}

func (fuu *FileUploadUsecase) uploadFile(f entity.File, i int) error {
	conf.Log.Info("Start uploading file", zap.Int("idx", i), zap.String("id", f.Id))

	// Get the file binary
	content, err := fuu.lg.GetFileContent(f.Id)
	if err != nil {
		return err
	}
	// upload the file binary
	file, err := fuu.br.SaveImageBinary(&f, content)
	if err != nil {
		return err
	}

	faceRes, faceErr := fuu.fg.GetFaceAnalysis(file.ContentUrl)
	if faceErr == nil {
		fuu.calcurateFaceScore(faceRes, file)
	}

	// Save file data
	err = fuu.fr.SaveFile(file)
	if err != nil {
		if faceErr != nil {
			return fmt.Errorf("failed to call face api and save file metadata; face err = %s ,file err = %w", faceErr, err)
		} else {
			return err
		}
	}
	if faceErr != nil {
		return faceErr
	}
	conf.Log.Info("Complete uploading file", zap.Int("idx", i), zap.String("id", f.Id))
	return nil
}

func (fuu *FileUploadUsecase) calcurateFaceScore(r []*dto.FaceResponse, f *entity.File) {
	if len(r) <= 0 || len(r) > 10 {
		f.FaceCount = 0
		f.FaceHappinessLevel = 0
		f.FacePhotoBeauty = 0
		f.FaceScore = 0
		return
	}
	faceCount := len(r)
	faceIds := make([]string, faceCount)
	var faceHappinessLevelSum float32
	var facePhotoBeautySum float32
	var hasMale bool
	var hasFemale bool
	var hasYoung bool
	var hasElderly bool
	for i, f := range r {
		faceIds[i] = f.FaceId

		// calculate the face happiness level (max: 40)
		faceHappinessLevelSum += 20 * f.FaceAttributes.Smile
		faceHappinessLevelSum -= 20 * f.FaceAttributes.Emotion.Anger
		faceHappinessLevelSum -= 10 * f.FaceAttributes.Emotion.Contempt
		faceHappinessLevelSum -= 15 * f.FaceAttributes.Emotion.Disgust
		faceHappinessLevelSum -= 5 * f.FaceAttributes.Emotion.Fear
		faceHappinessLevelSum += 20 * f.FaceAttributes.Emotion.Happiness
		faceHappinessLevelSum += 1 * f.FaceAttributes.Emotion.Neutral
		faceHappinessLevelSum += 5 * f.FaceAttributes.Emotion.Surprise

		// calculate the face photo beauty (max: 30)
		facePhotoBeautySum += 10 * (1 - f.FaceAttributes.Blur.Value)
		facePhotoBeautySum += 10 * (1 - f.FaceAttributes.Noise.Value)
		facePhotoBeautySum += 10 * (1 - 2*float32(math.Abs(0.5-float64(f.FaceAttributes.Exposure.Value))))

		if f.FaceAttributes.Occlusion.ForeheadOccluded {
			facePhotoBeautySum -= 2
		}
		if f.FaceAttributes.Occlusion.EyeOccluded {
			facePhotoBeautySum -= 4
		}
		if f.FaceAttributes.Occlusion.MouthOccluded {
			facePhotoBeautySum -= 2
		}

		// For bonus
		if f.FaceAttributes.Gender == "male" {
			hasMale = true
		}
		if f.FaceAttributes.Gender == "female" {
			hasFemale = true
		}
		if f.FaceAttributes.Age < 10 {
			hasYoung = true
		}
		if f.FaceAttributes.Age > 50 {
			hasElderly = true
		}
	}
	// calculate the face count bonus point (max: 20)
	bonusPoint := 2 * float32(faceCount)
	if hasMale && hasFemale {
		bonusPoint += 4
	}
	if hasYoung {
		bonusPoint += 3
	}
	if hasElderly {
		bonusPoint += 3
	}

	f.FaceCount = faceCount
	f.FaceHappinessLevel = faceHappinessLevelSum / float32(faceCount)
	f.FacePhotoBeauty = facePhotoBeautySum / float32(faceCount)
	f.FaceScore = f.FaceHappinessLevel + f.FacePhotoBeauty + bonusPoint
	f.Calculated = true
}
