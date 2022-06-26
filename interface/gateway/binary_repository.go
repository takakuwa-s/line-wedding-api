package gateway

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
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

func (br *BinaryRepository) getStoragePath(name string, fileType entity.FileType, isContent bool) string {
	if isContent {
		return string(fileType) + "/content/" + name
	} else {
		return string(fileType) + "/thumbnail/" + name
	}
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
	writer := br.getUploadWriter(br.getStoragePath(name, entity.Image, false))
	thumb := imaging.Resize(img, 200, 0, imaging.Lanczos)
	imaging.Encode(writer, thumb, imaging.JPEG)
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close the writer for image thumbnail; err = %w", err)
	}
	conf.Log.Info("Successfully upload the image thumbnail binary", zap.String("name", name), zap.Any("attrs", writer.Attrs()))
	return writer, nil
}

func (br *BinaryRepository) SaveImageBinary(file entity.File, content io.ReadCloser) (*entity.File, error) {
	defer content.Close()
	name := file.Id + "-" + time.Now().Format("2023-03-19-12:30:00")
	contentWriter := br.getUploadWriter(br.getStoragePath(name, entity.Image, true))
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
		br.deleteBinary(br.getStoragePath(name, entity.Image, false))
		return nil, fmt.Errorf("failed to close the writer for the image content binary; err = %w", err)
	}
	conf.Log.Info("Successfully upload the image binary content", zap.String("name", name), zap.Any("attrs", contentWriter.Attrs()))
	file.Width = img.Bounds().Dx()
	file.Height = img.Bounds().Dy()
	file.ThumbnailUrl = thumbWriter.Attrs().MediaLink
	file.ContentUrl = contentWriter.Attrs().MediaLink
	file.MimeType = contentWriter.Attrs().ContentType
	file.Uploaded = true
	file.Name = name
	return &file, nil
}

func (br *BinaryRepository) uploadVideoThumb(name, url string) (string, error) {
	tempDir, err := ioutil.TempDir("", "thumbnail*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory; err = %w", err)
	}
	outputFile := tempDir + "/thumbnail.jpg"
	cmd := fmt.Sprintf(`ffmpeg -y -i '%s' -ss 1 -vframes 1 '%s'`, url, outputFile)
	shellName := os.Getenv("FFMPEG_ENV")
	ffCmd := exec.Command(shellName, "-c", cmd)
	output, err := ffCmd.CombinedOutput()
	if err != nil {
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("failed to execute ffmpef cmd; cmd = %s, err = %w, output = %s", ffCmd.String(), err, string(output))
	}
	conf.Log.Info("ffmpeg is successfully completed", zap.Any("ffCmd", ffCmd))
	b, err := ioutil.ReadFile(outputFile)
	os.RemoveAll(tempDir)
	if err != nil {
		return "", fmt.Errorf("failed to read the created thumbnail image for video; err = %w", err)
	}
	content := bytes.NewReader(b)
	writer := br.getUploadWriter(br.getStoragePath(name, entity.Video, false))
	if _, err := io.Copy(writer, content); err != nil {
		return "", fmt.Errorf("failed to copy the video thumbnail binary; err = %w", err)
	}
	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close the writer for the content binary; err = %w", err)
	}
	return writer.Attrs().MediaLink, nil
}

func (br *BinaryRepository) SaveVideoBinary(file entity.File, content io.ReadCloser) (*entity.File, error) {
	defer content.Close()
	name := file.Id + "-" + time.Now().Format("2023-03-19-12:30:00")
	contentWriter := br.getUploadWriter(br.getStoragePath(name, entity.Video, true))
	if _, err := io.Copy(contentWriter, content); err != nil {
		return nil, fmt.Errorf("failed to copy the video content binary; err = %w", err)
	}
	if err := contentWriter.Close(); err != nil {
		return nil, fmt.Errorf("failed to close the writer for the content binary; err = %w", err)
	}
	contentUrl := contentWriter.Attrs().MediaLink
	thumbUrl, err := br.uploadVideoThumb(name, contentUrl)
	if err != nil {
		br.deleteBinary(br.getStoragePath(name, entity.Video, true))
		return nil, err
	}
	conf.Log.Info("Successfully upload the video binary content", zap.String("name", name), zap.Any("attrs", contentWriter.Attrs()))
	file.ThumbnailUrl = thumbUrl
	file.ContentUrl = contentUrl
	file.MimeType = contentWriter.Attrs().ContentType
	file.Uploaded = true
	file.Name = name
	return &file, nil
}

func (br *BinaryRepository) deleteBinary(name string) error {
	if err := br.f.Bucket.Object(name).Delete(conf.Ctx); err != nil {
		return fmt.Errorf("failed to delete the file binary; name = %s, err = %w", name, err)
	}
	conf.Log.Info("Successfully delete the file", zap.String("name", name))
	return nil
}

func (br *BinaryRepository) DeleteBinary(name string, fileType entity.FileType) error {
	if err := br.deleteBinary(br.getStoragePath(name, fileType, true)); err != nil {
		return err
	}
	if err := br.deleteBinary(br.getStoragePath(name, fileType, false)); err != nil {
		return fmt.Errorf("successfully deleted the binary content, but failed to delete the thumbnail binary, err = %w", err)
	}
	return nil
}
