package mongo

import (
	"authentication-ms/pkg/model"
	"authentication-ms/pkg/repo/mongo/document"
	"authentication-ms/pkg/svc"
	"context"
	"log"
	"time"
)

func (d *dal) AddMovie(ctx context.Context, movie model.Movie) error {
	if movie.Url == "" || movie.Name == "" {
		log.Println("important field missing")
		return svc.ErrBadRequest
	}

	movieDoc := document.Movie{
		Name:      movie.Name,
		MovieID:   movie.MovieId,
		OverView:  movie.OverView,
		Url:       movie.Url,
		ImageUrl:  movie.ImageUrl,
		LeadActor: movie.LeadActor,
		Tags:      movie.Tags,
		CreateTs:  time.Now(),
		UpdateTs:  time.Now(),
	}

	res, err := d.collMovieRec.InsertOne(ctx, movieDoc)
	if err != nil {
		log.Println("error in adding new movie in collection")
		return svc.ErrUnexpected
	}
	log.Println("_id", res.InsertedID, "movie inserted successfully")
	return nil
}
