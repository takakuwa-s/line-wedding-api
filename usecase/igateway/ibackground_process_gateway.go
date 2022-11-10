package igateway

type IBackgroundProcessGateway interface {
	StartUploadingFiles(ids []string) error
	StartDeletingFiles(ids []string) error
}
