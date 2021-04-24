package endpoints

import (
	"context"
	"uploadservice/service"

	"github.com/go-kit/kit/endpoint"
)

func makeUploadOneEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UploadOneRequest)
		result, err := s.UploadOne(ctx, req.File)
		response := UploadResponse{Result: result, Success: 1, Message: "Uploaded Successfully"}
		if err != nil {
			response.Success = 0
			response.Message = err.Error()
		}
		return response, nil
	}
}

func makeUploadManyEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UploadManyRequest)
		result, err := s.UploadMany(ctx, req.Files)
		response := UploadResponse{Result: result, Success: 1, Message: "Uploaded Successfully"}
		if err != nil {
			response.Success = 0
			response.Message = err.Error()
		}
		return response, nil
	}
}
