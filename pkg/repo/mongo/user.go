package mongo

import (
	"authentication-ms/pkg/svc"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

func (d *dal) UpdateUserMood(ctx context.Context, userId string, mood string) error {
	if userId == "" || mood == "" {
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
