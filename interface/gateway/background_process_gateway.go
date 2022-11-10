package gateway

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/dto"
	"go.uber.org/zap"
)

type BackgroundProcessGateway struct {
	lb *dto.LineBot
}

// Newコンストラクタ
func NewBackgroundProcessGateway(lb *dto.LineBot) *BackgroundProcessGateway {
	return &BackgroundProcessGateway{lb: lb}
}

func (bpg *BackgroundProcessGateway) StartUploadingFiles(ids []string) error {
	baseUrl := os.Getenv("BACKGROUND_PROCESS_API_BASE_URL")
	u, _ := url.Parse(baseUrl + "/api/file/list")
	q := u.Query()
	for _, id := range ids {
		q.Add("id", id)
	}
	u.RawQuery = q.Encode()
	if err := bpg.executeRequest("POST", u.String()); err != nil {
		return err
	}
	conf.Log.Info("Successfully calling uploading file binary api", zap.Strings("ids", ids))
	return nil
}

func (bpg *BackgroundProcessGateway) StartDeletingFiles(ids []string) error {
	baseUrl := os.Getenv("BACKGROUND_PROCESS_API_BASE_URL")
	u, _ := url.Parse(baseUrl + "/api/file/list")
	q := u.Query()
	for _, id := range ids {
		q.Add("id", id)
	}
	u.RawQuery = q.Encode()
	if err := bpg.executeRequest("DELETE", u.String()); err != nil {
		return err
	}
	conf.Log.Info("Successfully calling deleting file binary api", zap.Strings("ids", ids))
	return nil
}

func (bpg *BackgroundProcessGateway) executeRequest(method, url string) error {
	req, err := http.NewRequest(
		method,
		url,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create the http request for the background process api; %w", err)
	}
	token, err := bpg.lb.GetToken()
	if err != nil {
		return err
	}
	req.Header["Authorization"] = []string{"Bearer " + token}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call the background process api; %w", err)
	}
	if resp.StatusCode != http.StatusAccepted {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read the background process api response ; %w", err)
		}
		var obj map[string]interface{}
		if err := json.Unmarshal(body, &obj); err != nil {
			return fmt.Errorf("failed to convert the background process api response to object ; %w", err)
		}
		return fmt.Errorf("failed to call background process api ; %s", obj["error_description"])
	}
	return nil
}
