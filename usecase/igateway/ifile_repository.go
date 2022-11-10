package igateway

import (
	"github.com/takakuwa-s/line-wedding-api/entity"
)

type IFileRepository interface {
	SaveFile(file *entity.File) error
	UpdateForBrideAndGroomById(forBrideAndGroom bool, id string) error
	UpdateFileStatusByIdIn(fileStatus entity.FileStatus, ids []string) error
	DeleteFileById(id string) error
	DeleteFileByIds(ids []string) error
	FindById(id string) (*entity.File, error)
	FindByIds(ids []string) ([]entity.File, error)
	FindByIdsAndFileStatus(ids []string, fileStatus entity.FileStatus) ([]entity.File, error)
	FindByLimitAndStartIdAndUserIdAndFileTypeAndForBrideAndGroomAndFileStatusIn(limit int, startId, userId, orderBy, fileType string, forBrideAndGroom *bool, statuses []string) ([]entity.File, error)
	FindByFileStatusIn(statuses []entity.FileStatus) ([]entity.File, error)
	FindByFileStatusAndFileTypeAndForBrideAndGroom(limit int, fileStatus entity.FileStatus, forBrideAndGroom bool, fileType entity.FileType) ([]entity.File, error)
	FindByFileStatusAndFileTypeAndForBrideAndGroomAndDuration(limit int, fileStatus entity.FileStatus, forBrideAndGroom bool, fileType entity.FileType, duration int) ([]entity.File, error)
}
