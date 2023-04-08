package handler

import (
	"authentication-ms/pkg/svc"
	"authentication-ms/pkg/transport/middleware"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type MovieGetRequest struct {
	Tags []string `json:"tags"`
}

func GetMovies(s svc.SVC) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		vars := mux.Vars(r)
		tag := vars["tag"]
		log.Println(tag)
		if tag == "" {
			middleware.WriteJsonHttpErrorResponse(w, http.StatusBadRequest, errBadRequest)
			return
		}
		res, err := s.GetMoviesByTag(ctx, tag)
		if err != nil {
			log.Println("error in getting movies...", err)
			middleware.WriteJsonHttpErrorResponse(w, http.StatusInternalServerError, err)
			return
		}
		middleware.WriteJsonHttpResponse(w, http.StatusOK, res)

	}
}

func GetMoviesByTags(s svc.SVC) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var request MovieGetRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			log.Println("error in decoding request body...", err)
			middleware.WriteJsonHttpErrorResponse(w, http.StatusBadRequest, errBadRequest)
			return
		}

		log.Println(request)

		if request.Tags == nil {
			middleware.WriteJsonHttpErrorResponse(w, http.StatusBadRequest, errBadRequest)
			return
		}
		res, err := s.GetMoviesByTags(ctx, request.Tags)
		if err != nil {
			log.Println("error in getting movies...", err)
			middleware.WriteJsonHttpErrorResponse(w, http.StatusInternalServerError, err)
			return
		}
		middleware.WriteJsonHttpResponse(w, http.StatusOK, res)

	}
}
