package handler

import (
	"authentication-ms/pkg/model"
	"authentication-ms/pkg/svc"
	"authentication-ms/pkg/transport/middleware"
	"encoding/json"
	"log"
	"net/http"
)

type forgotRequest struct {
	Email string `json:"email"`
}

func ForgotPassword(s svc.SVC) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request forgotRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			log.Println("error in decoding request body...", err)
			middleware.WriteJsonHttpErrorResponse(w, http.StatusBadRequest, errBadRequest)
			return
		}
		log.Println(request)
		if request.Email == "" {
			log.Println("necessary field missing, ", err)
			middleware.WriteJsonHttpErrorResponse(w, http.StatusBadRequest, errBadRequest)
			return
		}
		user := model.User{
			Email: request.Email,
		}
		err = s.ForgotPassword(r.Context(), user)
		if err != nil {
			log.Println("error in forgot password : ", err)
			if err == svc.ErrNoData {
				middleware.WriteJsonHttpErrorResponse(w, http.StatusNotFound, svc.ErrNoData)
			} else {
				middleware.WriteJsonHttpErrorResponse(w, http.StatusInternalServerError, svc.ErrUnexpected)
			}
			return
		}
		middleware.WriteJsonHttpResponse(w, http.StatusOK, "successfully changed")
	}
}
