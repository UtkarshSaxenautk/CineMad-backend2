package handler

import (
	"authentication-ms/pkg/svc"
	"authentication-ms/pkg/transport/middleware"
	"encoding/json"
	"log"
	"net/http"
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

		err = s.UpdateUserWatchedMovies(r.Context(), request.Jwt, request.Movie)
		if err != nil {
			log.Println("error in updating watched movie...", err)
			middleware.WriteJsonHttpErrorResponse(w, http.StatusBadRequest, err)
			return
		}
		middleware.WriteJsonHttpResponse(w, http.StatusOK, "success")

	}
}
