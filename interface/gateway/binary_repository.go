package gateway

import (
	"fmt"
	"image"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"github.com/disintegration/imaging"
	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/dto"
	"github.com/takakuwa-s/line-wedding-api/entity"
	"go.uber.org/zap"
)

type BinaryRepository struct {
	f *dto.Firestore
}

// Newコンストラクタ
func NewBinaryRepository(f *dto.Firestore) *BinaryRepository {
	return &BinaryRepository{f: f}
}

func (br *BinaryRepository) getUploadWriter(name string) *storage.Writer {
	writer := br.f.Bucket.Object(name).NewWriter(conf.Ctx)
	writer.ObjectAttrs.ACL = []storage.ACLRule{
		{
			Entity: storage.AllUsers,
			Role:   storage.RoleReader,
		},
	}
	return writer
}

func (br *BinaryRepository) uploadImageThumb(name string, img image.Image) (*storage.Writer, error) {
	writer := br.getUploadWriter("thumbnail/" + name)
	thumb := imaging.Resize(img, 200, 0, imaging.Lanczos)
	imaging.Encode(writer, thumb, imaging.JPEG)
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close the writer for thumbnail; err = %w", err)
	}
	conf.Log.Info("Successfully upload the thumbnail binary", zap.String("name", name), zap.Any("attrs", writer.Attrs()))
	return writer, nil
}

func (br *BinaryRepository) deleteBinary(name string) error {
	if err := br.f.Bucket.Object(name).Delete(conf.Ctx); err != nil {
		return fmt.Errorf("failed to delete the file binary; name = %s, err = %w", name, err)
	}
	conf.Log.Info("Successfully delete the file", zap.String("name", name))
	return nil
}

func (br *BinaryRepository) SaveImageBinary(file *entity.File, content io.ReadCloser) (*entity.File, error) {
	name := file.Id + "-" + time.Now().Format("2006-01-02-15:04:05")
	contentWriter := br.getUploadWriter("content/" + name)
	reader := io.TeeReader(content, contentWriter)
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image; err = %w", err)
	}
	thumbWriter, err := br.uploadImageThumb(name, img)
	if err != nil {
		return nil, err
	}
	if err := contentWriter.Close(); err != nil {
		br.deleteBinary("thumbnail/" + name)
		return nil, fmt.Errorf("failed to close the writer for the content binary; err = %w", err)
	}
	conf.Log.Info("Successfully upload the binary content", zap.String("name", name), zap.Any("attrs", contentWriter.Attrs()))
	file.Width = img.Bounds().Dx()
	file.Height = img.Bounds().Dy()
	file.ThumbnailUrl = thumbWriter.Attrs().MediaLink
	file.ContentUrl = contentWriter.Attrs().MediaLink
	file.MimeType = contentWriter.Attrs().ContentType
	file.Uploaded = true
	file.Name = name
	return file, nil
}

func (br *BinaryRepository) DeleteBinary(name string) error {
	if err := br.deleteBinary("content/" + name); err != nil {
		return err
	}
	if err := br.deleteBinary("thumbnail/" + name); err != nil {
		return fmt.Errorf("successfully deleted the binary content, but failed to delete the thumbnail binary, err = %w", err)
	}
	return nil
}
