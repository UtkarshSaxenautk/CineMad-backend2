package document

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type DateOfBirth struct {
	Year      int    `bson:"year"`
	Month     int    `bson:"month"`
	MonthName string `bson:"month_name"`
	Date      int    `bson:"date"`
}

type User struct {
	ID                primitive.ObjectID   `bson:"_id,omitempty"`
	Email             string               `bson:"email,omitempty"`
	Username          string               `bson:"username,omitempty"`
	PasswordHash      string               `bson:"password_hash,omitempty"`
	FullName          string               `bson:"full_name,omitempty"`
	MoviesWatched     []primitive.ObjectID `bson:"movies_watched"`
	WatchLater        []Movie              `bson:"watch_later"`
	MoodPreviously    string               `bson:"previous_mood"`
	Dob               DateOfBirth          `bson:"date_of_birth,omitempty"`
	CreateTs          time.Time            `bson:"create_ts"`
	UpdateTs          time.Time            `bson:"update_ts"`
	LoginTs           time.Time            `bson:"login_ts"`
	PreviousPasswords []string             `bson:"previous_passwords"`
}
