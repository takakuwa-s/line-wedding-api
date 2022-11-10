package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/dto"
	"go.uber.org/zap"
)

type SlideShowGateway struct {
}

// Newコンストラクタ
func NewSlideShowGateway() *SlideShowGateway {
	return &SlideShowGateway{}
}

func (sg *SlideShowGateway) Render(imageUrls, videoUrls []string) (*dto.SlideShowCreateResponce, error) {
	reqBody, err := sg.createBody(imageUrls, videoUrls)
	if err != nil {
		return nil, err
	}
	baseUrl := os.Getenv("SHOTSTACK_API_BASE_URL")
	u, _ := url.Parse(baseUrl + "templates/render")
	req, err := http.NewRequest(
		"POST",
		u.String(),
		reqBody,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create the http request for the shotstack api; %w", err)
	}
	apiKey := os.Getenv("SHOTSTACK_API_KEY")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send the shotstack api request; %w", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the shotstack api response body; %w", err)
	}
	if resp.StatusCode == http.StatusCreated {
		var obj *dto.SlideShowCreateResponce
		if err := json.Unmarshal(body, &obj); err != nil {
			return nil, fmt.Errorf("failed to json Unmarshal for the shotstack api success response; %w", err)
		}
		conf.Log.Info("Successfully get the shotstack api response", zap.Any("response", obj))
		return obj, nil
	} else {
		return nil, fmt.Errorf("shotstack api returns error; %s", string(body))
	}
}

func (sg *SlideShowGateway) createBody(imageUrls, videoUrls []string) (io.Reader, error) {
	id := os.Getenv("SHOTSTACK_TEMPLATE_ID")
	req := dto.NewTemplateRender(id)
	sg.appendUrl(req, imageUrls, "IMAGE_URL_")
	sg.appendUrl(req, videoUrls, "VIDEO_URL_")
	json, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to json Marshal for the shotstack api request; %w", err)
	}
	conf.Log.Info("Successfully create request body for shotstack api", zap.String("json", string(json)))
	return bytes.NewBuffer(json), nil
}

func (sg *SlideShowGateway) appendUrl(r *dto.TemplateRender, urls []string, findPrefix string) {
	for i, u := range urls {
		find := fmt.Sprintf("%s%d", findPrefix, i)
		m := dto.NewMergeField(find, u)
		r.ApendMerge(m)
	}
}

func (sg *SlideShowGateway) DownloadContent(url string) (io.Reader, error) {
	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get the content, url = %s, err = %w", url, err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the response body, url = %s, err = %w", url, err)
	}
	reader := bytes.NewReader(b)
	return reader, nil
}
