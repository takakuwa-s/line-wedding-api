package gateway

import (
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"go.uber.org/zap"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/dto"
	"github.com/takakuwa-s/line-wedding-api/entity"
)

type FileRepository struct {
	cr *CommonRepository
	f  *dto.Firestore
}

// Newコンストラクタ
func NewFileRepository(cr *CommonRepository, f *dto.Firestore) *FileRepository {
	return &FileRepository{cr: cr, f: f}
}

func (fr *FileRepository) SaveFile(file *entity.File) error {
	return fr.cr.Save("files", file.Id, file)
}

func (fr *FileRepository) DeleteFileById(id string) error {
	return fr.cr.DeleteById("files", id)
}

func (fr *FileRepository) DeleteFileByIds(ids []string) error {
	return fr.cr.DeleteByIds("files", ids)
}

func (fr *FileRepository) UpdateForBrideAndGroomById(forBrideAndGroom bool, id string) error {
	if _, err := fr.f.Firestore.Collection("files").Doc(id).Update(conf.Ctx, []firestore.Update{
		{
			Path:  "ForBrideAndGroom",
			Value: forBrideAndGroom,
		},
		{
			Path:  "UpdatedAt",
			Value: time.Now(),
		},
	}); err != nil {
		return fmt.Errorf("failed to update the file; id =  %s, forBrideAndGroom = %t, err = %w", id, forBrideAndGroom, err)
	}
	conf.Log.Info("Successfully update the file", zap.String("id", id), zap.Bool("forBrideAndGroom", forBrideAndGroom))
	return nil
}

func (fr *FileRepository) UpdateFileStatusByIdIn(fileStatus entity.FileStatus, ids []string) error {
	batch := fr.f.Firestore.Batch()
	for _, list := range fr.cr.SplitSlice(ids) {
		query := fr.f.Firestore.Collection("files").Where("Id", "in", list)
		iter := query.Documents(conf.Ctx)
		for {
			dsnap, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return fmt.Errorf("failed to get a file metadata ; err = %w", err)
			}
			batch.Update(dsnap.Ref, []firestore.Update{
				{
					Path:  "FileStatus",
					Value: fileStatus,
				},
				{
					Path:  "UpdatedAt",
					Value: time.Now(),
				},
			})
		}
	}
	_, err := batch.Commit(conf.Ctx)
	if err != nil {
		return fmt.Errorf("failed to update the file; id =  %s, fileStatus = %s, err = %w", ids, fileStatus, err)
	}
	conf.Log.Info("Successfully update the file", zap.Any("ids", ids), zap.Any("fileStatus", fileStatus))
	return nil
}

func (fr *FileRepository) FindById(id string) (*entity.File, error) {
	dsnap, err := fr.cr.FindById("files", id)
	if err != nil {
		return nil, err
	}
	if dsnap == nil {
		return nil, nil
	}
	var file entity.File
	dsnap.DataTo(&file)
	return &file, nil
}

func (fr *FileRepository) executeQuery(query *firestore.Query) ([]entity.File, error) {
	var files []entity.File
	iter := query.Documents(conf.Ctx)
	for {
		dsnap, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get a file metadata ; err = %w", err)
		}
		var f entity.File
		dsnap.DataTo(&f)
		files = append(files, f)
	}
	if files == nil {
		return []entity.File{}, nil
	}
	return files, nil
}

func (fr *FileRepository) FindByIds(ids []string) ([]entity.File, error) {
	var files []entity.File
	for _, list := range fr.cr.SplitSlice(ids) {
		query := fr.f.Firestore.Collection("files").Where("Id", "in", list)
		f, err := fr.executeQuery(&query)
		if err != nil {
			return nil, err
		}
		files = append(files, f...)
	}
	conf.Log.Info("Successfully find the file metadata with", zap.Int("file count", len(files)), zap.Strings("ids", ids))
	return files, nil
}

func (fr *FileRepository) FindByIdsAndFileStatus(ids []string, fileStatus entity.FileStatus) ([]entity.File, error) {
	var files []entity.File
	for _, list := range fr.cr.SplitSlice(ids) {
		query := fr.f.Firestore.Collection("files").Where("Id", "in", list).Where("FileStatus", "==", fileStatus)
		f, err := fr.executeQuery(&query)
		if err != nil {
			return nil, err
		}
		files = append(files, f...)
	}
	conf.Log.Info("Successfully find the file metadata with", zap.Int("file count", len(files)), zap.Strings("ids", ids), zap.String("FileStatus", string(fileStatus)))
	return files, nil
}

func (fr *FileRepository) FindByLimitAndStartIdAndUserIdAndFileTypeAndForBrideAndGroomAndFileStatusIn(limit int, startId, userId, orderBy, fileType string, forBrideAndGroom *bool, statuses []string) ([]entity.File, error) {
	query := fr.f.Firestore.Collection("files").OrderBy(orderBy, firestore.Desc)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if startId != "" {
		dsnap, err := fr.f.Firestore.Collection("files").Doc(startId).Get(conf.Ctx)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				return []entity.File{}, nil
			} else {
				return nil, fmt.Errorf("failed to get the file metadata by startId; id =  %s err = %w", startId, err)
			}
		}
		query = query.StartAfter(dsnap)
	}
	if userId != "" {
		query = query.Where("Creater", "==", userId)
	}
	if fileType != "" {
		query = query.Where("FileType", "==", fileType)
	}
	if len(statuses) > 0 {
		query = query.Where("FileStatus", "in", statuses)
	}
	if forBrideAndGroom != nil {
		query = query.Where("ForBrideAndGroom", "==", &forBrideAndGroom)
	}
	files, err := fr.executeQuery(&query)
	if err != nil {
		return nil, err
	}
	conf.Log.Info("Successfully find the file metadata with", zap.Int("file count", len(files)), zap.Int("limit", limit), zap.String("startId", startId), zap.String("userId", userId))
	return files, nil
}

func (fr *FileRepository) FindByFileStatusIn(statuses []entity.FileStatus) ([]entity.File, error) {
	query := fr.f.Firestore.Collection("files").Where("FileStatus", "in", statuses).OrderBy("UpdatedAt", firestore.Asc)
	f, err := fr.executeQuery(&query)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fr *FileRepository) FindByFaceCountAndFileStatusAndFileTypeAndForBrideAndGroom(limit, faceCount int, fileStatus entity.FileStatus, forBrideAndGroom bool, fileType entity.FileType) ([]entity.File, error) {
	query := fr.f.Firestore.Collection("files").Where("FaceCount", ">=", faceCount).Where("FileStatus", "==", fileStatus).Where("ForBrideAndGroom", "==", forBrideAndGroom).Where("FileType", "==", string(fileType)).Limit(limit)
	f, err := fr.executeQuery(&query)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fr *FileRepository) FindByFileStatusAndFileTypeAndForBrideAndGroomAndDuration(limit int, fileStatus entity.FileStatus, forBrideAndGroom bool, fileType entity.FileType, duration int) ([]entity.File, error) {
	query := fr.f.Firestore.Collection("files").Where("FileStatus", "==", fileStatus).Where("ForBrideAndGroom", "==", forBrideAndGroom).Where("FileType", "==", string(fileType)).Where("Duration", ">=", duration).Limit(limit)
	f, err := fr.executeQuery(&query)
	if err != nil {
		return nil, err
	}
	return f, nil
}
