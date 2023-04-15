package handler

import (
	"authentication-ms/pkg/model"
	"authentication-ms/pkg/svc"
	"authentication-ms/pkg/transport/middleware"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type otpRequest struct {
	Otp string
}

func ProcessOtp(s svc.SVC) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		email := vars["email"]
		var request otpRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			log.Println("error in decoding request body...", err)
			middleware.WriteJsonHttpErrorResponse(w, http.StatusBadRequest, errBadRequest)
			return
		}

		log.Println(request)
		if request.Otp == "" || email == "" {
			log.Println("necessary field missing, ", err)
			middleware.WriteJsonHttpErrorResponse(w, http.StatusBadRequest, errBadRequest)
			return
		}
		user := model.User{
			Email: email,
		}
		pass, err := s.ProcessOtp(user, request.Otp)
		if err != nil {
			log.Println("error in processing password : ", err)
			middleware.WriteJsonHttpErrorResponse(w, http.StatusNotFound, svc.ErrUnexpected)
			return
		}
		if !pass {
			log.Println("otp invalid ")
			middleware.WriteJsonHttpErrorResponse(w, http.StatusForbidden, svc.ErrUserNotAuthorized)
			return
		}

		middleware.WriteJsonHttpResponse(w, http.StatusOK, "success")
	}
}
