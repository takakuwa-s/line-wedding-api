package gateway

import (
	"bytes"
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

type FaceGateway struct {
}

// Newコンストラクタ
func NewFaceGateway() *FaceGateway {
	return &FaceGateway{}
}

func (fg *FaceGateway) GetFaceAnalysis(imageUrl string) ([]*dto.FaceResponse, error) {
	jsonStr := `{"url":"` + imageUrl + `"}`
	baseUrl := os.Getenv("FACE_API_BASE_URL")
	u, _ := url.Parse(baseUrl + "/face/v1.0/detect")
	q := u.Query()
	q.Add("detectionModel", "detection_01")
	q.Add("returnFaceId", "true")
	q.Add("returnFaceLandmarks", "false")
	q.Add("returnFaceAttributes", "age,gender,smile,emotion,blur,exposure,noise,occlusion,headPose")
	q.Add("recognitionModel", "recognition_04")
	u.RawQuery = q.Encode()
	req, err := http.NewRequest(
		"POST",
		u.String(),
		bytes.NewBuffer([]byte(jsonStr)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create the http request for the face api; %w", err)
	}
	subscriptionKey := os.Getenv("OCP_APIM_SUBSCRIPTION_KEY")
	req.Header.Set("Ocp-Apim-Subscription-Key", subscriptionKey)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send the face api request; %w", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the face api response body; %w", err)
	}
	if resp.StatusCode == http.StatusOK {
		var obj []*dto.FaceResponse
		if err := json.Unmarshal(body, &obj); err != nil {
			return nil, fmt.Errorf("failed to json Unmarshal for the face api success response; %w", err)
		}
		conf.Log.Info("Successfully get the face api response", zap.Any("response", obj))
		return obj, nil
	} else {
		return nil, fmt.Errorf("face api returns error; %s", string(body))
	}
}
