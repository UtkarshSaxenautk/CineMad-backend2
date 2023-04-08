package mongo

import (
	"authentication-ms/pkg/model"
	"authentication-ms/pkg/repo/mongo/document"
	"authentication-ms/pkg/svc"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

const (
	CollLogRec   string = "authentication_log_rec"
	CollMovieRec string = "movie_rec"
)

type password struct {
	hash string
}

type prevPasswords struct {
	hashPasswords []string
}

type dal struct {
	db           *mongo.Database
	collLogRec   *mongo.Collection
	collMovieRec *mongo.Collection
}

func NewDal(db *mongo.Database) svc.Dao {
	CollLogRec := db.Collection(CollLogRec)
	CollMovieRec := db.Collection(CollMovieRec)
	return &dal{db, CollLogRec, CollMovieRec}
}

func (d *dal) CreateUser(ctx context.Context, user model.User) error {
	// Insert new User in db
	userDoc := bson.D{{"email", user.Email},
		{"username", user.Username},
		{"password_hash", user.PasswordHash},
		{"full_name", user.FullName},
		{"role", user.Role},
		{"date_of_birth", user.Dob},
		{"create_ts", time.Now()},
		{"update_ts", time.Now()}}

	_, err := d.collLogRec.InsertOne(ctx, userDoc)
	if err != nil {
		log.Println("error in saving userdata at signup ", err)
		return err
	}
	return nil
}

func (d *dal) CheckEmailAndUserName(ctx context.Context, user model.User) (emailExist bool, usernameExist bool, err error) {
	filter1 := bson.M{
		"email": user.Email,
	}
	filter2 := bson.M{
		"username": user.Username,
	}
	countEmail, err := d.collLogRec.CountDocuments(ctx, filter1)
	if err != nil {
		log.Fatal("error in checking document email :  ", err)
		return emailExist, usernameExist, err
	}
	countUserName, err := d.collLogRec.CountDocuments(ctx, filter2)
	if err != nil {
		log.Fatal("error in checking document username :  ", err)
		return emailExist, usernameExist, err
	}
	if countEmail > 0 {
		emailExist = true
	}
	if countUserName > 0 {
		usernameExist = true
	}
	return emailExist, usernameExist, nil
}

func (d *dal) checkInPreviousPasswords(ctx context.Context, user model.User, newPassword string) (bool, error) {
	var prev prevPasswords
	filter := bson.M{
		"email": user.Email,
	}
	err := d.collLogRec.FindOne(ctx, filter).Decode(&prev)
	if err != nil {
		log.Println("error in getting previous passwords : ", err)
		return false, err
	}
	for _, curr := range prev.hashPasswords {
		if curr == newPassword {
			return true, nil
		}
	}
	return false, nil
}

func (d *dal) UpdatePassword(ctx context.Context, user model.User, nPass string) error {
	check, err := d.checkInPreviousPasswords(ctx, user, nPass)
	if err != nil {
		log.Println("error in checking new password with last used :", err)
		return svc.ErrUnexpected
	}
	if check == true {
		log.Println("password is previously used! ")
		return svc.ErrBadRequest
	}
	filter := bson.M{
		"email": user.Email,
	}
	update := bson.M{
		"$set": bson.M{
			"password_hash": nPass,
		},
	}
	opts := options.Update().SetUpsert(true)
	res, err := d.collLogRec.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Println("error in updating login_ts ", err)
		return svc.ErrUnexpected
	}

	if res.ModifiedCount > 0 || res.UpsertedCount > 0 {
		return nil
	}
	return svc.ErrUnexpected
}

func (d *dal) GetUser(ctx context.Context, email string) (string, error) {
	var user document.User
	if email == "" {
		log.Println("empty email field in repo layer at login time ")
		return "", svc.ErrNoData
	}
	filter := bson.M{
		"email": email,
	}
	update := bson.M{
		"$set": bson.M{
			"login_ts": time.Now(),
		},
	}
	// First update login timestamp
	res, err := d.collLogRec.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println("error in updating login_ts ", err)
		return "", svc.ErrUnexpected
	}

	if res.ModifiedCount > 0 || res.UpsertedCount > 0 {
		// If no error get hashedPassword
		err = d.collLogRec.FindOne(ctx, filter).Decode(&user)
		if err != nil {
			log.Println("error in getting hash by email : ", err)
			return "", err
		}
		log.Println("pass : ", user.PasswordHash)
		return user.PasswordHash, nil
	}

	return "", svc.ErrUnexpected
}
