package gateway

import (
	"fmt"
	"io"
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
	parent := &drive.ParentReference{
		Id: "1yziSpir9Io-LjaGNz0_uOpfHT3YecFYN",
	}
	name := "image-" + time.Now().Format("2006-01-02-15:04:05.000000")
	f := &drive.File{
		Title:     name,
		Shareable: true,
		Parents:   []*drive.ParentReference{parent},
	}
	res, err := br.srv.Files.Insert(f).Media(content).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to insert the file; err = %w", err)
	}
	conf.Log.Info("2", zap.Any("OriginalFilename", res.OriginalFilename), zap.Any("Id", res.Id), zap.Any("DefaultOpenWithLink", res.DefaultOpenWithLink), zap.Any("WebContentLink", res.WebContentLink), zap.Any("WebViewLink", res.WebViewLink))
	file.FileId = res.Id
	file.Uri = res.DefaultOpenWithLink
	file.IsUploaded = true
	file.Name = name
	return file, nil
}

func (br *BinaryRepository) DeleteBinary(id string) error {
	if err :=	br.srv.Files.Delete(id).Do(); err != nil {
		return fmt.Errorf("failed to delete the file; err = %w", err)
	}
	conf.Log.Info("Successfully delete the file", zap.String("id", id))
	return nil
}
