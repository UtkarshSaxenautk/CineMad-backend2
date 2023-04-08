package handler

import (
	"authentication-ms/pkg/svc"
	"authentication-ms/pkg/transport/middleware"
	"errors"
	"net/http"
)

var (
	errBadRequest   = errors.New("bad request")
	errMissingField = errors.New("one of the required field is missing")
)

func GetHealth(s svc.SVC) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		middleware.WriteJsonHttpResponse(w, http.StatusOK, "hello world")

	}
}
