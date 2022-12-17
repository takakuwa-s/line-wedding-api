package dto

type ConfigResponce struct {
	FileFeatureAvailable       bool `json:"fileFeatureAvailable"`
	AttendanceFeatureAvailable bool `json:"attendanceFeatureAvailable"`
}

func NewConfigResponce(fileFeatureAvailable, attendanceFeatureAvailable bool) ConfigResponce {
	return ConfigResponce{
		FileFeatureAvailable:       fileFeatureAvailable,
		AttendanceFeatureAvailable: attendanceFeatureAvailable,
	}
}
