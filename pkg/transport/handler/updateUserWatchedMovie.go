package handler

import (
	"authentication-ms/pkg/svc"
	"authentication-ms/pkg/transport/middleware"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type UpdateMoviesRequest struct {
	Jwt   string `json:"jwt"`
	Movie string `json:"movie"`
}

func UpdateUserWatchedMovie(s svc.SVC) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request UpdateMoviesRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			log.Println("error in decoding request body...", err)
			middleware.WriteJsonHttpErrorResponse(w, http.StatusBadRequest, errBadRequest)
			return
		}

		log.Println(request)
		if request.Jwt == "" || request.Movie == "" {
			log.Println("necessary field missing, ", err)
			middleware.WriteJsonHttpErrorResponse(w, http.StatusBadRequest, errBadRequest)
			return
		}
		ctx := r.Context()
		if _, err := strconv.ParseInt(request.Movie, 10, 64); err == nil {
			log.Println("looks like a number")
			err = s.UpdateWatchedMovieByMovieID(ctx, request.Jwt, request.Movie)
			if err != nil {
				log.Println("error in updating watched movie...", err)
				middleware.WriteJsonHttpErrorResponse(w, http.StatusBadRequest, err)
				return
			}
		} else {
			err = s.UpdateUserWatchedMovies(ctx, request.Jwt, request.Movie)
			if err != nil {
				log.Println("error in updating watched movie...", err)
				middleware.WriteJsonHttpErrorResponse(w, http.StatusBadRequest, err)
				return
			}
		}
		middleware.WriteJsonHttpResponse(w, http.StatusOK, "success")

	}
}
