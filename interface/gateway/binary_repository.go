package gateway

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/entity"
	"go.uber.org/zap"
	"google.golang.org/api/drive/v2"
)

type BinaryRepository struct {
	srv *drive.Service
}

// Newコンストラクタ
func NewBinaryRepository() *BinaryRepository {
	srv, err := drive.NewService(conf.Ctx)
	if err != nil {
		panic(fmt.Sprintf("Unable to retrieve Drive client; err = %v", err))
	}
	return &BinaryRepository{srv: srv}
}

func (br *BinaryRepository) SaveBinary(file *entity.File, content io.ReadCloser) (*entity.File, error) {
	parentReference := os.Getenv("GOOGLE_PARENT_REFERENCE")
	parent := &drive.ParentReference{
		Id: parentReference,
	}
	name := file.FileType.ToString() + "-" + time.Now().Format("2006-01-02-15:04:05.000000")
	f := &drive.File{
		Title:     name,
		Shareable: true,
		Parents:   []*drive.ParentReference{parent},
		HasThumbnail: true,
	}
	res, err := br.srv.Files.Insert(f).Media(content).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to insert the file; err = %w", err)
	}
	conf.Log.Info("Successfully insert the file", zap.Any("res", res))
	file.FileId = res.Id
	file.ContentUrl = "https://drive.google.com/uc?export=view&id=" + res.Id
	if file.FileType == entity.ImageType {
		file.ThumbnailUrl = res.ThumbnailLink
		file.Width = res.ImageMediaMetadata.Width
		file.Height = res.ImageMediaMetadata.Height
	}
	file.MimeType = res.MimeType
	file.IsUploaded = true
	file.Name = name
	return file, nil
}

func (br *BinaryRepository) DeleteBinary(id string) error {
	if err :=	br.srv.Files.Delete(id).Do(); err != nil {
		return fmt.Errorf("failed to delete the file binary; id = %s, err = %w", id, err)
	}
	conf.Log.Info("Successfully delete the file", zap.String("id", id))
	return nil
}
