package handler

import (
	"authentication-ms/pkg/svc"
	"authentication-ms/pkg/transport/middleware"
	"encoding/json"
	"log"
	"net/http"
)

type GetMovieInRequest struct {
	Jwt  string   `json:"jwt"`
	Mood []string `json:"mood"`
}

func GetMovieAccordingToMood(s svc.SVC) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Body)
		//(w).Header().Set("Access-Control-Allow-Origin", "http://localhost:1234")
		var request GetMovieInRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			log.Println("error in decoding request body...", err)
			middleware.WriteJsonHttpErrorResponse(w, http.StatusBadRequest, errBadRequest)
			return
		}

		log.Println(request)
		if request.Jwt == "" {
			log.Println("necessary field missing, ", err)
			middleware.WriteJsonHttpErrorResponse(w, http.StatusBadRequest, errBadRequest)
			return
		}
		ctx := r.Context()
		log.Println("mood we got : ", request.Mood)
		res, err := s.GetMoviesAccordingToUserMood(ctx, request.Jwt, request.Mood)
		if err != nil {
			log.Println("error in getting userProfile...", err)
			middleware.WriteJsonHttpErrorResponse(w, http.StatusInternalServerError, err)
			return
		}
		middleware.WriteJsonHttpResponse(w, http.StatusOK, res)

	}
}
