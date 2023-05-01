package mongo

import (
	"authentication-ms/pkg/model"
	"authentication-ms/pkg/repo/mongo/document"
	"authentication-ms/pkg/svc"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

const (
	Limit = 50
)

func (d *dal) AddMovie(ctx context.Context, movie model.Movie) (string, error) {
	if movie.Url == "" || movie.Name == "" {
		log.Println("important field missing")
		return "", svc.ErrBadRequest
	}

	exist, err := d.checkMovieID(ctx, movie.MovieId)
	if err != nil {
		log.Println("error in checking movie existence ")
		return "", svc.ErrUnexpected
	}
	if exist {
		log.Println("movie already exist")
		return "", svc.ErrUnexpected
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
		return "", svc.ErrUnexpected
	}
	log.Println("_id", res.InsertedID, "movie inserted successfully")
	objId, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Println("error in converting insertedId interface to ObjectId   ")
		return "", svc.ErrUnexpected
	}
	return objId.Hex(), nil
}

func (d *dal) checkMovieID(ctx context.Context, movieID int) (bool, error) {
	if movieID == 0 {
		log.Println("important field missing ")
	}

	filter := bson.M{"movie_id": movieID}
	countMovie, err := d.collMovieRec.CountDocuments(ctx, filter)
	if err != nil {
		log.Fatal("error in checking document email :  ", err)
		return false, err
	}
	return countMovie > 0, nil
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
			ImageUrl:  movie.ImageUrl,
			MovieId:   movie.MovieID,
			LeadActor: movie.LeadActor,
			Tags:      movie.Tags,
		})
	}
	log.Println(movies, "  in repo")
	return movies, nil
}

func (d *dal) GetMovieByMovieID(ctx context.Context, movieID string) (model.Movie, error) {
	if movieID == "" {
		return model.Movie{}, svc.ErrBadRequest
	}

	mid, err := primitive.ObjectIDFromHex(movieID)
	if err != nil {
		log.Println("invalid hex : ", movieID)
		return model.Movie{}, err
	}
	filter := bson.M{
		"_id": mid,
	}

	//filter := bson.D{}
	//opts := options.Find().SetLimit(int64(Limit))
	var docMovie document.Movie
	err = d.collMovieRec.FindOne(ctx, filter).Decode(&docMovie)
	if err != nil {
		log.Println("error in Find query")
		return model.Movie{}, err
	}

	var movie model.Movie

	movie = model.Movie{
		ID:        docMovie.ID.Hex(),
		Name:      docMovie.Name,
		Url:       docMovie.Url,
		ImageUrl:  docMovie.ImageUrl,
		MovieId:   docMovie.MovieID,
		LeadActor: docMovie.LeadActor,
		Tags:      docMovie.Tags,
	}
	log.Println(movie, "  in repo")
	return movie, nil
}
