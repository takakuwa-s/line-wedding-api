package gateway

import (
	"fmt"

	vision "cloud.google.com/go/vision/apiv1"

	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/dto"
)

type FaceGateway struct {
}

// Newコンストラクタ
func NewFaceGateway() *FaceGateway {
	return &FaceGateway{}
}

func (fg *FaceGateway) GetFaceAnalysis(imageUrl string) ([]*dto.FaceResponse, error) {
	// https://cloud.google.com/vision/docs/detecting-faces?hl=ja
	client, err := vision.NewImageAnnotatorClient(conf.Ctx)
	if err != nil {
		return nil, fmt.Errorf("error cannot create the image annotator client for Google Cloud Vision API; %w", err)
	}
	image := vision.NewImageFromURI(imageUrl)

	// https://cloud.google.com/vision/docs/reference/rest/v1/AnnotateImageResponse#FaceAnnotation
	annotations, err := client.DetectFaces(conf.Ctx, image, nil, 10)
	if err != nil {
		return nil, fmt.Errorf("error occurs when excuting Google Cloud Vision API; %w", err)
	}
	var res []*dto.FaceResponse
	if len(annotations) != 0 {
		for _, annotation := range annotations {
			f := &dto.FaceResponse{}
			f.DetectionConfidence = annotation.DetectionConfidence
			f.RollAngle = annotation.RollAngle
			f.PanAngle = annotation.PanAngle
			f.TiltAngle = annotation.TiltAngle
			f.JoyLikelihood = int32(annotation.JoyLikelihood)
			f.SorrowLikelihood = int32(annotation.SorrowLikelihood)
			f.AngerLikelihood = int32(annotation.AngerLikelihood)
			f.SurpriseLikelihood = int32(annotation.SurpriseLikelihood)
			f.UnderExposedLikelihood = int32(annotation.UnderExposedLikelihood)
			f.BlurredLikelihood = int32(annotation.BlurredLikelihood)
			f.HeadwearLikelihood = int32(annotation.HeadwearLikelihood)
			res = append(res, f)
		}
	}
	return res, nil
}
