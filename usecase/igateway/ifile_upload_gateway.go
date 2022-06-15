package igateway

type IFileUploadGateway interface {
	StartUploadingFiles(ids []string) error
}
