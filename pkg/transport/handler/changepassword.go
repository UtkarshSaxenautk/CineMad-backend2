package handler

import (
	"authentication-ms/pkg/model"
	"authentication-ms/pkg/svc"
	"authentication-ms/pkg/transport/middleware"
	"encoding/json"
	"log"
	"net/http"
)

type ChangePasswordRequest struct {
	Email        string `json:"email"`
	Username     string `json:"username"`
	PasswordHash string `json:"passwordHash"`
	NewPassword  string `json:"newPassword"`
}

func ChangePassword(s svc.SVC) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request ChangePasswordRequest
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
		user := model.User{
			Email:        request.Email,
			PasswordHash: request.PasswordHash,
		}
		authorized, err := s.SignIn(r.Context(), request.Email, request.PasswordHash)
		if err != nil {
			log.Println("error in checking signin : ", err)
			middleware.WriteJsonHttpErrorResponse(w, http.StatusInternalServerError, svc.ErrUnexpected)
			return
		}
		if authorized == false {
			log.Println("user not authorized to changed password")
			middleware.WriteJsonHttpErrorResponse(w, http.StatusForbidden, svc.ErrUserNotAuthorized)
			return
		}
		log.Println("calling svc layer to update")
		err = s.ChangePassword(r.Context(), user, request.NewPassword)

		if err != nil {
			log.Println("internal server error in changing password : ", err)
			middleware.WriteJsonHttpErrorResponse(w, http.StatusInternalServerError, svc.ErrUnexpected)
			return
		}
		middleware.WriteJsonHttpResponse(w, http.StatusOK, "successfully changed")
	}
}
