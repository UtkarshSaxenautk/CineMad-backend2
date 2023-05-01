package handler

import (
	"authentication-ms/pkg/svc"
	"authentication-ms/pkg/transport/middleware"
	"encoding/json"
	"log"
	"net/http"
)

type UpdateWatchLaterRequest struct {
	Jwt       string `json:"jwt"`
	ID        string `json:"id"`
	IsMovieDB bool   `json:"isMovieDB"`
	Type      string `json:"type"`
}

type deleteWatchLater struct {
	Jwt     string `json:"jwt"`
	MovieID string `json:"movie_id"`
}

func UpdateUserWatchLaterMovie(s svc.SVC) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request UpdateWatchLaterRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			log.Println("error in decoding request body...", err)
			middleware.WriteJsonHttpErrorResponse(w, http.StatusBadRequest, errBadRequest)
			return
		}

		log.Println(request)
		if request.Jwt == "" || request.ID == "" {
			log.Println("necessary field missing, ", err)
			middleware.WriteJsonHttpErrorResponse(w, http.StatusBadRequest, errBadRequest)
			return
		}

		err = s.UpdateWatchLater(r.Context(), request.Jwt, request.ID, request.IsMovieDB, request.Type)
		if err != nil {
			log.Println("error in updating watchLater movie...", err)
			middleware.WriteJsonHttpErrorResponse(w, http.StatusBadRequest, err)
			return
		}
		middleware.WriteJsonHttpResponse(w, http.StatusOK, "success")

	}
}

func DeleteWatchLaterMovie(s svc.SVC) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request deleteWatchLater
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			log.Println("error in decoding request body...", err)
			middleware.WriteJsonHttpErrorResponse(w, http.StatusBadRequest, errBadRequest)
			return
		}

		log.Println(request)
		if request.Jwt == "" || request.MovieID == "" {
			log.Println("necessary field missing, ", err)
			middleware.WriteJsonHttpErrorResponse(w, http.StatusBadRequest, errBadRequest)
			return
		}
		ctx := r.Context()
		err = s.DeleteWatchLater(ctx, request.Jwt, request.MovieID)
		if err != nil {
			log.Println("error in deleting watchLater movie...", err)
			middleware.WriteJsonHttpErrorResponse(w, http.StatusBadRequest, err)
			return
		}
		middleware.WriteJsonHttpResponse(w, http.StatusOK, "success")

	}
}
