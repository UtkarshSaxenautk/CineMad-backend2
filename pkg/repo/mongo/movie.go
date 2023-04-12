package mongo

import (
	"authentication-ms/pkg/model"
	"authentication-ms/pkg/repo/mongo/document"
	"authentication-ms/pkg/svc"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

const (
	Limit = 50
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

func (d *dal) GetMoviesByTags(ctx context.Context, tags []string) ([]model.Movie, error) {
	if len(tags) == 0 {
		log.Println("tags are empty")
		return nil, svc.ErrMissingImportantField
	}
	filter := bson.M{
		"tags": bson.M{
			"$elemMatch": bson.M{
				"$in": tags,
			},
		},
	}

	//filter := bson.D{}
	opts := options.Find().SetLimit(int64(Limit))
	cursor, err := d.collMovieRec.Find(ctx, filter, opts)
	if err != nil {
		log.Println("error in Find query")
		return nil, err
	}
	var docMovies []document.Movie
	err = cursor.All(ctx, &docMovies)
	if err != nil {
		log.Println("error in cursor decode")
		return nil, err
	}
	var movies []model.Movie
	for _, movie := range docMovies {

		movies = append(movies, model.Movie{
			ID:        movie.ID.Hex(),
			Name:      movie.Name,
			Url:       movie.Url,
			ImageUrl:  movie.Url,
			MovieId:   movie.MovieID,
			LeadActor: movie.LeadActor,
			Tags:      movie.Tags,
		})
	}
	log.Println(movies, "  in repo")
	return movies, nil
}
