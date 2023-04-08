package handler

import (
	"authentication-ms/pkg/svc"
	"authentication-ms/pkg/transport/middleware"
	"encoding/json"
	"log"
	"net/http"
)

type SignInRequest struct {
	Email        string `json:"email"`
	PasswordHash string `json:"passwordHash"`
}

func SignIn(s svc.SVC) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request SignInRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			log.Println("error in decoding request body...", err)
			middleware.WriteJsonHttpErrorResponse(w, http.StatusBadRequest, errBadRequest)
			return
		}

		log.Println(request)
		if request.Email == "" || request.PasswordHash == "" {
			log.Println("necessary field missing, ", err)
			middleware.WriteJsonHttpErrorResponse(w, http.StatusBadRequest, errBadRequest)
			return
		}

		res, err := s.SignIn(r.Context(), request.Email, request.PasswordHash)
		if err != nil || res == false {
			log.Println("error in signIn...", err)
			middleware.WriteJsonHttpErrorResponse(w, http.StatusBadRequest, err)
			return
		}
		middleware.WriteJsonHttpResponse(w, http.StatusOK, "successfully signIn")

	}
}
