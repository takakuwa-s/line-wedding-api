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

type FileUploadGateway struct {
	lb *dto.LineBot
}

// Newコンストラクタ
func NewFileUploadGateway(lb *dto.LineBot) *FileUploadGateway {
	return &FileUploadGateway{lb: lb}
}

func (fug *FileUploadGateway) StartUploadingFiles(ids []string) error {
	baseUrl := os.Getenv("FILE_UPLOAD_API_BASE_URL")
	u, _ := url.Parse(baseUrl + "/api/file/list")
	q := u.Query()
	for _, id := range ids {
		q.Add("id", id)
	}
	u.RawQuery = q.Encode()
	req, err := http.NewRequest(
		"POST",
		u.String(),
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create the http request for the file upload api; %w", err)
	}
	token, err := fug.lb.GetToken()
	if err != nil {
		return err
	}
	req.Header["Authorization"] = []string{"Bearer " + token}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call the file upload api; %w", err)
	}
	if resp.StatusCode != http.StatusAccepted {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read the file upload api response ; %w", err)
		}
		var obj map[string]interface{}
		if err := json.Unmarshal(body, &obj); err != nil {
			return fmt.Errorf("failed to convert the file upload api response to object ; %w", err)
		}
		return fmt.Errorf("failed to call file upload api ; %s", obj["error_description"])
	}
	conf.Log.Info("Successfully calling uploading file binary api", zap.Strings("ids", ids))
	return nil
}
