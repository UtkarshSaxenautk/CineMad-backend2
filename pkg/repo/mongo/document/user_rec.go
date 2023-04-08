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
	ID                primitive.ObjectID
	Email             string      `bson:"email,omitempty"`
	Username          string      `bson:"username,omitempty"`
	PasswordHash      string      `bson:"password_hash,omitempty"`
	FullName          string      `bson:"full_name,omitempty"`
	Role              string      `bson:"role"`
	Dob               DateOfBirth `bson:"date_of_birth,omitempty"`
	CreateTs          time.Time   `bson:"create_ts"`
	UpdateTs          time.Time   `bson:"update_ts"`
	LoginTs           time.Time   `bson:"login_ts"`
	PreviousPasswords []string    `bson:"previous_passwords"`
}
