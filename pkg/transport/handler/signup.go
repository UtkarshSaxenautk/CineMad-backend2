package handler

import (
	"authentication-ms/pkg/model"
	"authentication-ms/pkg/svc"
	"authentication-ms/pkg/transport/middleware"
	"encoding/json"
	"log"
	"net/http"
)

type DobRequest struct {
	Year      int    `json:"year"`
	Month     int    `json:"month"`
	MonthName string `json:"monthName"`
	Date      int    `json:"date"`
}

type SignupRequest struct {
	Email        string     `json:"email"`
	Username     string     `json:"username"`
	PasswordHash string     `json:"passwordHash"`
	FullName     string     `json:"fullName"`
	Role         string     `json:"role"`
	Dob          DobRequest `json:"dob"`
}

func Signup(s svc.SVC) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request SignupRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			log.Println("error in decoding request body...", err)
			middleware.WriteJsonHttpErrorResponse(w, http.StatusBadRequest, errBadRequest)
			return
		}
		emptyDob := DobRequest{
			Date:      -1,
			Year:      -1,
			Month:     -1,
			MonthName: "",
		}
		log.Println(request)
		if request.Email == "" || request.PasswordHash == "" || request.FullName == "" || request.Dob == emptyDob {
			log.Println("necessary field missing, ", err)
			middleware.WriteJsonHttpErrorResponse(w, http.StatusBadRequest, errBadRequest)
			return
		}
		newUser := model.User{
			Email:        request.Email,
			PasswordHash: request.PasswordHash,
			Username:     request.Username,
			FullName:     request.FullName,
			Dob: model.DateOfBirth{
				Date:      request.Dob.Date,
				Month:     request.Dob.Month,
				Year:      request.Dob.Year,
				MonthName: request.Dob.MonthName,
			},
		}
		err = s.Signup(r.Context(), newUser)
		if err != nil {
			log.Println("error in signup...", err)
			middleware.WriteJsonHttpErrorResponse(w, http.StatusInternalServerError, err)
			return
		}
		middleware.WriteJsonHttpResponse(w, http.StatusOK, "successfully signup")

	}
}
