package dto

type HeadPose struct {
	Pitch float32 `json:"pitch"`
	Roll float32 `json:"roll"`
	Yaw float32 `json:"yaw"`
}

type Emotion struct {
	Anger float32 `json:"anger"`
	Contempt float32 `json:"contempt"`
	Disgust float32 `json:"disgust"`
	Fear float32 `json:"fear"`
	Happiness float32 `json:"happiness"`
	Neutral float32 `json:"neutral"`
	Surprise float32 `json:"surprise"`
}

/**
 * blur: face is blurry or not.
 * Level returns 'Low', 'Medium' or 'High'.
 * Value returns a number between [0,1], the larger the blurrier.
 */
type Blur struct {
	BlurLevel string `json:"blurLevel"`
	Value float32 `json:"value"`
}

/**
 * exposure: face exposure level.
 * Level returns 'GoodExposure', 'OverExposure' or 'UnderExposure'.
 */
type Exposure struct {
	ExposureLevel string `json:"exposureLevel"`
	Value float32 `json:"value"`
}

/**
 * noise: noise level of face pixels.
 * Level returns 'Low', 'Medium' and 'High'.
 * Value returns a number between [0,1], the larger the noisier
 */
type Noise struct {
	NoiseLevel string `json:"noiseLevel"`
	Value float32 `json:"value"`
}

/**
 * occlusion: whether each facial area is occluded.
 * including forehead, eyes and mouth.
 */
type Occlusion struct {
	ForeheadOccluded bool `json:"foreheadOccluded"`
	EyeOccluded bool `json:"eyeOccluded"`
	MouthOccluded bool `json:"mouthOccluded"`
}

type FaceAttributes struct {
	Smile float32 `json:"smile"`
	HeadPose HeadPose `json:"headPose"`
	Age	float32 `json:"age"`
	Gender string `json:"gender"`
	Emotion Emotion `json:"emotion"`
	Blur Blur `json:"blur"`
	Exposure Exposure `json:"exposure"`
	Noise Noise `json:"noise"`
	Occlusion Occlusion `json:"occlusion"`
}

type FaceResponse struct {
	FaceId string `json:"faceId"`
	FaceAttributes FaceAttributes `json:"faceAttributes"`
}