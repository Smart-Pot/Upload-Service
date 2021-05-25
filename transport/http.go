package transport

import (
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"strconv"
	"uploadservice/endpoints"

	pkghttp "github.com/Smart-Pot/pkg/common/http"
	"github.com/Smart-Pot/pkg/common/perrors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

const (
	fileLimit = 10 << 20
)

func MakeHTTPHandlers(e endpoints.Endpoints, logger log.Logger) http.Handler {
	r := mux.NewRouter().PathPrefix("/upload").Subrouter()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(perrors.EncodeHTTPError),
	}

	r.Methods("POST").Path("/").Handler(httptransport.NewServer(
		e.UploadOne,
		decodeUploadOneRequest,
		encodeHTTPResponse,
		options...,
	))

	r.Methods("POST").Path("/many").Handler(httptransport.NewServer(
		e.UploadMany,
		decodeUploadManyRequest,
		encodeHTTPResponse,
		options...,
	))

	r.Methods("GET").PathPrefix("/").Handler(
		http.StripPrefix("/upload/", http.FileServer(
			http.Dir("./upload"),
		)))

	return pkghttp.EnableCORS(r)
}

func encodeHTTPResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func decodeUploadOneRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if err := r.ParseMultipartForm(fileLimit); err != nil {
		return nil, err
	}
	file, _, err := r.FormFile("image")
	if err != nil {
		return nil, err
	}
	return endpoints.UploadOneRequest{
		File: file,
	}, nil
}

func decodeUploadManyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if err := r.ParseMultipartForm(fileLimit); err != nil {
		return nil, err
	}
	size, err := strconv.Atoi(r.FormValue("size"))
	files := make([]multipart.File, size)
	if err != nil {
		return nil, err
	}
	for x := 0; x < size; x++ {
		file, _, err := r.FormFile("image" + strconv.Itoa(x+1))
		if err != nil {
			return nil, err
		}
		files[x] = file
	}
	return endpoints.UploadManyRequest{
		Files: files,
	}, nil
}
