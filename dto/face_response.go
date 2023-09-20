package dto

type FaceResponse struct {
	// Detection confidence. Range [0, 1].
	DetectionConfidence float32

	// Roll angle, which indicates the amount of clockwise/anti-clockwise rotation
	// of the face relative to the image vertical about the axis perpendicular to
	// the face. Range [-180,180].
	RollAngle float32
	// Yaw angle, which indicates the leftward/rightward angle that the face is
	// pointing relative to the vertical plane perpendicular to the image. Range
	// [-180,180].
	PanAngle float32
	// Pitch angle, which indicates the upwards/downwards angle that the face is
	// pointing relative to the image's horizontal plane. Range [-180,180].
	TiltAngle float32

	// unknown = 0
	// very unlikely = 1
	// unlikely = 2
	// possible = 3
	// likely = 4
	// very likely = 5
	JoyLikelihood          int32
	SorrowLikelihood       int32
	AngerLikelihood        int32
	SurpriseLikelihood     int32
	UnderExposedLikelihood int32
	BlurredLikelihood      int32
	HeadwearLikelihood     int32
}
