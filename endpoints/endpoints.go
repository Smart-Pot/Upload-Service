package endpoints

import (
	"mime/multipart"
	"uploadservice/service"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	UploadOne  endpoint.Endpoint
	UploadMany endpoint.Endpoint
}

type UploadResponse struct {
	Result  interface{}
	Success int32
	Message string
}

type UploadOneRequest struct {
	File multipart.File
}

type UploadManyRequest struct {
	Files []multipart.File
}

func MakeEndpoints(s service.Service) Endpoints {
	return Endpoints{
		UploadOne:  makeUploadOneEndpoint(s),
		UploadMany: makeUploadManyEndpoint(s),
	}
}
