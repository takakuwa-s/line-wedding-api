package usecase

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/dto"
	"github.com/takakuwa-s/line-wedding-api/entity"
	"github.com/takakuwa-s/line-wedding-api/usecase/igateway"
	"github.com/takakuwa-s/line-wedding-api/utils"
	"go.uber.org/zap"
)

type SlideShowUsecase struct {
	mr  igateway.IMessageRepository
	fr  igateway.IFileRepository
	sg  igateway.ISlideShowGateway
	sr  igateway.ISlideShowRepository
	br  igateway.IBinaryRepository
	lpu *LinePushUsecase
}

// Newコンストラクタ
func NewSlideShowUsecase(
	mr igateway.IMessageRepository,
	fr igateway.IFileRepository,
	sg igateway.ISlideShowGateway,
	sr igateway.ISlideShowRepository,
	br igateway.IBinaryRepository,
	lpu *LinePushUsecase) *SlideShowUsecase {
	return &SlideShowUsecase{mr: mr, fr: fr, sg: sg, sr: sr, br: br, lpu: lpu}
}

func (su *SlideShowUsecase) CreateSlideShow() (*dto.SlideShowCreateResponce, error) {
	imageUrls, err := su.getImagesForSlideshow()
	if err != nil {
		return nil, err
	}
	videoUrls, err := su.getVideoesForSlideshow()
	if err != nil {
		return nil, err
	}
	res, err := su.sg.Render(imageUrls, videoUrls)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (su *SlideShowUsecase) getImagesForSlideshow() ([]string, error) {
	limit := 30
	images, err := su.fr.FindByFaceCountAndFileStatusAndFileTypeAndForBrideAndGroom(500, 1, entity.Open, false, entity.Image)
	if err != nil {
		return nil, err
	}
	if len(images) < limit {
		return nil, fmt.Errorf("the number of images is %d and less than %d", len(images), limit)
	}
	conf.Log.Info("images for slideshow are found", zap.Int("images size", len(images)))
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(images), func(i, j int) { images[i], images[j] = images[j], images[i] })

	urls := []string{}
	idxSet := []int{}
	keySet := []string{}
	prevIdx := 0
	for lap := 1; lap <= 3; lap++ {
		for i := 0; i < len(images) && len(urls) < limit; i++ {
			if utils.IntContains(idxSet, i) {
				continue
			}
			cur := images[i]
			switch lap {
			case 1:
				key := cur.Creater + "-" + strconv.Itoa(cur.FaceCount)
				if cur.Width < cur.Height || utils.StringContains(keySet, key) {
					continue
				}
				keySet = append(keySet, key)
				conf.Log.Info("Image is selected at lap 1", zap.Int("curIdx", i), zap.String("key", key))
			case 2:
				prev := images[prevIdx]
				if cur.Width < cur.Height || (cur.Creater == prev.Creater && cur.FaceCount == prev.FaceCount) {
					continue
				}
				conf.Log.Info("Image is selected at lap 2", zap.Int("curIdx", i), zap.Int("prevIdx", prevIdx))
			case 3:
				prev := images[prevIdx]
				if cur.Creater == prev.Creater && cur.FaceCount == prev.FaceCount {
					continue
				}
				conf.Log.Info("Image is selected at lap 3", zap.Int("curIdx", i), zap.Int("prevIdx", prevIdx))
			}
			urls = append(urls, cur.ContentUrl)
			idxSet = append(idxSet, i)
			prevIdx = i
		}
	}
	if len(urls) < limit {
		return nil, fmt.Errorf("too many same images, hence cannot select %d images; urls = %d; images = %d", limit, len(urls), len(images))
	}
	conf.Log.Info("All images for slideshow are selected")
	return urls[:limit], nil
}

func (su *SlideShowUsecase) getVideoesForSlideshow() ([]string, error) {
	return []string{}, nil
	// limit := 1
	// videoes, err := su.fr.FindByFileStatusAndFileTypeAndForBrideAndGroomAndDuration(300, entity.Open, false, entity.Video, 6000)
	// if err != nil {
	// 	return nil, err
	// }
	// if len(videoes) < limit {
	// 	return nil, fmt.Errorf("the number of videoes is %d and less than %d", len(videoes), limit)
	// }
	// var urls []string
	// for _, f := range videoes {
	// 	urls = append(urls, f.ContentUrl)
	// }
	// rand.Seed(time.Now().UnixNano())
	// rand.Shuffle(len(urls), func(i, j int) { urls[i], urls[j] = urls[j], urls[i] })
	// return urls[:limit], nil
}

func (su *SlideShowUsecase) UploadSlideshow(r dto.SlideshowWebhook) error {
	switch r.Status {
	case "done":
		err := su.uploadSlideshow(r)
		if err != nil {
			note := fmt.Sprintf("uploading failed; error = %s", err.Error())
			return su.lpu.SendSlideshowErrorNotification(note)
		}
	case "failed":
		note := fmt.Sprintf("id = %s, error = %s, type = %s, action = %s", r.Id, r.Error, r.Type, r.Action)
		return su.lpu.SendSlideshowErrorNotification(note)
	default:
		return nil
	}
	return nil
}

func (su *SlideShowUsecase) uploadSlideshow(r dto.SlideshowWebhook) error {
	if r.Type != "edit" || r.Action != "render" {
		return nil
	}
	s, err := su.sr.FindById(r.Id)
	if err != nil {
		return err
	}
	if s != nil && len(s.ContentUrl) > 0 {
		return nil
	}
	content, err := su.sg.DownloadContent(r.Url)
	if err != nil {
		return err
	}
	s = entity.NewSlideShow(r.Id)
	s1, err := su.br.SaveSlideShowBinary(*s, content)
	if err != nil {
		return err
	}
	if err := su.sr.SaveSlideShow(s1); err != nil {
		return err
	}
	return su.lpu.SendSlideshowSuccessNotification(s1.ContentUrl, s1.ThumbnailUrl)
}

func (su *SlideShowUsecase) ListSlideshow() ([]entity.SlideShow, error) {
	list, err := su.sr.FindAllOrderByUpdatedAt()
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (su *SlideShowUsecase) PatchSlideshow(id string, req dto.PatchSlideShowRequest) error {
	s, err := su.sr.FindById(id)
	if err != nil {
		return err
	}
	if s == nil {
		return fmt.Errorf("not found the slideshow with id = %s", id)
	}
	if err := su.sr.UpdateSelectedById(req.Selected, id); err != nil {
		return err
	}
	return nil
}

func (su *SlideShowUsecase) DeleteSlideshow(id string) error {
	s, err := su.sr.FindById(id)
	if err != nil {
		return err
	}
	if s == nil {
		return fmt.Errorf("not found the slideshow with id = %s", id)
	}
	if err := su.br.DeleteBinary(s.Name, "slideshow"); err != nil {
		return fmt.Errorf("failed to delete the slideshow binary; id = %s, err = %w", id, err)
	}
	if err := su.sr.DeleteById(s.Id); err != nil {
		return err
	}
	return nil
}

func (su *SlideShowUsecase) CreateSlideshowMessage() ([]map[string]interface{}, error) {
	list, err := su.sr.FindBySelectedOrderByUpdatedAt(true)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, fmt.Errorf("not found created slide show")
	}
	messages := su.mr.FindMessageByKey("slideshow")
	messages[2]["originalContentUrl"] = fmt.Sprintf(messages[2]["originalContentUrl"].(string), list[0].ContentUrl)
	messages[2]["previewImageUrl"] = fmt.Sprintf(messages[2]["previewImageUrl"].(string), list[0].ThumbnailUrl)
	return messages, nil
}
