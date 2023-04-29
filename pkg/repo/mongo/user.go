package mongo

import (
	"authentication-ms/pkg/model"
	"authentication-ms/pkg/repo/mongo/document"
	"authentication-ms/pkg/svc"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func (d *dal) UpdateWatchLater(ctx context.Context, userId string, movies []string) error {
	if userId == "" {
		log.Println("important field missing : userID ")
		return svc.ErrMissingImportantField
	}
	uid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Println("error in converting userId into objectId")
		return err
	}
	var mids []primitive.ObjectID
	for _, movieID := range movies {
		temp, err := primitive.ObjectIDFromHex(movieID)
		if err != nil {
			continue
		}
		mids = append(mids, temp)
	}
	filter := bson.M{"_id": uid}
	update := bson.M{
		"$set": bson.M{
			"watch_later": mids,
		},
	}
	res, err := d.collLogRec.UpdateOne(ctx, filter, update)

	if err != nil {
		log.Println("error in saving user watchLater details")
		return svc.ErrUnexpected
	}
	if res.ModifiedCount > 0 || res.UpsertedCount > 0 {
		log.Println("userID: ", userId, " := updated user's watchLater details successfully")
		return nil
	}
	log.Println("error in saving user watchLater details")
	return svc.ErrUnexpected

}

func (d *dal) GetWatchLater(ctx context.Context, userID string) ([]model.Movie, error) {
	if userID == "" {
		log.Println("userID missing .....")
		return nil, svc.ErrMissingImportantField
	}
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var docUser document.User
	filter := bson.M{
		"_id": uid,
	}
	projection := bson.D{
		{Key: "_id", Value: 1},
		{Key: "watch_later", Value: 1},
	}
	opts := options.FindOne().SetProjection(projection)
	res := d.collLogRec.FindOne(ctx, filter, opts)
	err = res.Err()
	if err == mongo.ErrNoDocuments {
		err = svc.ErrNoData
		return nil, err
	} else if err != nil {
		log.Println("error in getting user of this userID")
		return nil, err
	}
	err = res.Decode(&docUser)
	if err != nil {
		log.Println("error in decoding result")
		return nil, err
	}
	var movies []model.Movie
	for _, docMovie := range docUser.WatchLater {
		temp := model.Movie{
			ID:        docMovie.ID.Hex(),
			Name:      docMovie.Name,
			OverView:  docMovie.OverView,
			Url:       docMovie.Url,
			ImageUrl:  docMovie.ImageUrl,
			LeadActor: docMovie.LeadActor,
			Tags:      docMovie.Tags,
		}
		movies = append(movies, temp)
	}

	log.Println("got successfully user of userID")
	return movies, nil

}

func (d *dal) UpdateUserMood(ctx context.Context, userId string, mood []string) error {
	if userId == "" || len(mood) == 0 {
		log.Println("important field missing")
		return svc.ErrMissingImportantField
	}
	uid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Println(err)
		return err
	}
	filter := bson.M{"_id": uid}
	update := bson.M{
		"$set": bson.M{
			"previous_mood": mood,
		},
	}
	res, err := d.collLogRec.UpdateOne(ctx, filter, update)

	if err != nil {
		log.Println("error in saving user mood details")
		return svc.ErrUnexpected
	}
	if res.ModifiedCount > 0 || res.UpsertedCount > 0 {
		log.Println("userID: ", userId, " := updated user's mood details successfully")
		return nil
	}
	log.Println("error in saving user mood details")
	return svc.ErrUnexpected
}

func (d *dal) UpdateUserWatchedMovies(ctx context.Context, userID string, movieID string) error {
	if userID == "" || movieID == "" {
		log.Println("userID : ", userID, " or movieId : ", movieID, " is empty")
		return svc.ErrMissingImportantField
	}

	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println("error in converting userID into objectID")
		return err
	}
	mid, err := primitive.ObjectIDFromHex(movieID)
	if err != nil {
		log.Println("error in converting movieID into objectID")
		return err
	}
	filter := bson.M{
		"_id": uid,
	}
	update := bson.D{
		{"$push", bson.D{
			{"movies_watched", mid},
		}},
	}
	res, err := d.collLogRec.UpdateOne(ctx, filter, update)

	if err != nil {
		log.Println("error in saving user watched movie details")
		return svc.ErrUnexpected
	}
	if res.ModifiedCount > 0 || res.UpsertedCount > 0 {
		log.Println("userID: ", userID, " := updated user's watched movies details successfully", " movieID : ", mid)
		return nil
	}
	log.Println("error in saving user watched movie details")
	return svc.ErrUnexpected
}

func (d *dal) GetUserProfile(ctx context.Context, userID string) (user model.User, err error) {

	if userID == "" {
		log.Println("userID is missing")
		err = svc.ErrMissingImportantField
		return
	}
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println("error in converting userId into objectID")
		return
	}
	var docUser document.User
	filter := bson.M{
		"_id": uid,
	}

	res := d.collLogRec.FindOne(ctx, filter)
	err = res.Err()
	if err == mongo.ErrNoDocuments {
		err = svc.ErrNoData
		return
	} else if err != nil {
		log.Println("error in getting user of this userID")
		return
	}
	err = res.Decode(&docUser)
	if err != nil {
		log.Println("error in decoding result")
		return
	}
	var movies []string
	for _, docMovie := range docUser.MoviesWatched {
		temp := docMovie.Hex()
		movies = append(movies, temp)
	}
	user = model.User{
		UserID:         docUser.ID.Hex(),
		Username:       docUser.Username,
		Email:          docUser.Email,
		FullName:       docUser.FullName,
		MoviesWatched:  movies,
		MoodPreviously: docUser.MoodPreviously,
	}
	log.Println("got successfully user of userID")
	return
}

func (d *dal) AddMovieToWatchLater(ctx context.Context, userID string, movie model.Movie) error {
	if userID == "" {
		log.Println("userId empty in dao..")
		return svc.ErrMissingImportantField
	}
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println("error in converting userID to objectID")
		return err
	}
	if movie.ID == "" {
		movie.ID = primitive.NewObjectID().Hex()
	}
	mid, err := primitive.ObjectIDFromHex(movie.ID)
	if err != nil {
		log.Println("error in converting movieID to objectID")
		return err
	}

	// create a filter to find the user by ID
	filter := bson.M{"_id": userObjectID}

	// create an update to add the movie to the watch_later array
	docMovie := document.Movie{
		Name:      movie.Name,
		ID:        mid,
		ImageUrl:  movie.ImageUrl,
		Url:       movie.Url,
		LeadActor: movie.LeadActor,
		MovieID:   movie.MovieId,
		Tags:      movie.Tags,
	}
	update := bson.M{"$addToSet": bson.M{"watch_later": docMovie}}

	res, err := d.collLogRec.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println("error in updating watchLater")
		return err
	}
	if res.ModifiedCount == 0 {
		// no document was modified, user not found
		return svc.ErrNoData
	}
	return nil
}
