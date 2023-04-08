package handler

import (
	"authentication-ms/pkg/model"
	"authentication-ms/pkg/svc"
	"authentication-ms/pkg/transport/middleware"
	"encoding/json"
	"log"
	"net/http"
)

type MovieRequest struct {
	Name      string   `json:"name"`
	MovieID   int      `json:"movie_id"`
	Tags      []string `json:"tags"`
	OverView  string   `json:"over_view"`
	LeadActor string   `json:"lead_actor"`
	Url       string   `json:"url"`
	ImageUrl  string   `json:"image_url"`
}

func AddMovie(s svc.SVC) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request MovieRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			log.Println("error in decoding request body...", err)
			middleware.WriteJsonHttpErrorResponse(w, http.StatusBadRequest, errBadRequest)
			return
		}
		log.Println(request)
		if request.Name == "" || request.Url == "" || request.ImageUrl == "" || len(request.Tags) <= 0 {
			log.Println("necessary field missing, ", err)
			middleware.WriteJsonHttpErrorResponse(w, http.StatusBadRequest, errBadRequest)
			return
		}
		newMovie := model.Movie{
			Name:      request.Name,
			MovieId:   request.MovieID,
			OverView:  request.OverView,
			LeadActor: request.LeadActor,
			Url:       request.Url,
			ImageUrl:  request.ImageUrl,
			Tags:      request.Tags,
		}
		err = s.AddMovieInDB(r.Context(), newMovie)
		if err != nil {
			log.Println("error in signup...", err)
			middleware.WriteJsonHttpErrorResponse(w, http.StatusInternalServerError, err)
			return
		}
		middleware.WriteJsonHttpResponse(w, http.StatusOK, "successfully added")

	}
}
