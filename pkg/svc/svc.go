package svc

import (
	"authentication-ms/pkg/model"
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type svc struct {
	dao   Dao
	sdk   Sdk
	cache Cache
	mail  Mail
}

var (
	//these are public errors, please do not include any technical wordings here
	ErrNoData                = errors.New("no data found")
	ErrBadRequest            = errors.New("mandatory input data missing")
	ErrDeleteFailed          = errors.New("domain delete failed")
	ErrUnexpected            = errors.New("unexpected error")
	ErrEmailAlreadyInUse     = errors.New("email is already in use")
	ErrUserNameAlreadyInUse  = errors.New("username is already in use")
	ErrUserNotAuthorized     = errors.New("user not authenticated")
	ErrMissingImportantField = errors.New("important field missing")
)

func New(dao Dao, cache Cache, sdk Sdk, mail Mail) SVC {
	s := &svc{dao, sdk, cache, mail}
	return s
}

func (s *svc) Signup(ctx context.Context, user model.User) error {
	log.Println(user)
	if user.Email == "" || user.Username == "" || user.PasswordHash == "" || user.FullName == "" {
		log.Println("missing necessary field...")
		return ErrBadRequest
	}
	emailExist, userNameExist, err := s.dao.CheckEmailAndUserName(ctx, user)
	if err != nil {
		log.Println("error in checking email and username existence..")
		return err
	}
	if emailExist {
		log.Println("email already in use...")
		return ErrEmailAlreadyInUse
	}
	if userNameExist {
		log.Println("username already in use...")
		return ErrUserNameAlreadyInUse
	}
	user.PasswordHash, err = s.hashPassword(user.PasswordHash)
	if err != nil {
		log.Println("error in creating password hash: ", err)
		return ErrUnexpected
	}
	err = s.dao.CreateUser(ctx, user)
	if err != nil {
		log.Println("error in creating user...")
		return ErrUnexpected
	}
	log.Println("user created successfully")
	return nil
}

func (s *svc) hashPassword(password string) (string, error) {
	// Convert password string to byte slice
	var passwordBytes = []byte(password)

	// Hash password with Bcrypt's min cost
	hashedPasswordBytes, err := bcrypt.
		GenerateFromPassword(passwordBytes, bcrypt.MinCost)

	return string(hashedPasswordBytes), err
}

func (s *svc) passwordsMatch(hashedPassword, currPassword string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword), []byte(currPassword))
	return err == nil
}

func (s *svc) createJWt(userID string) (string, error) {
	secretKey := []byte("ussr")

	// Create a new token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set the claims for the token
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = userID                   // Subject (typically the user ID)
	claims["iat"] = time.Now().Unix()        // Issued At (current time)
	claims["exp"] = time.Now().Unix() + 3600 // Expiration Time (1 hour from now)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		fmt.Println("Error creating JWT:", err)
		return "", err
	}
	return tokenString, nil
}

func (s *svc) SignIn(ctx context.Context, email string, password string) (string, error) {
	if email == "" || password == "" {
		return "", ErrBadRequest
	}
	hashedPassword, userID, err := s.dao.GetUser(ctx, email)

	if err != nil {
		log.Println("error in getting user from email")
		return "", ErrUnexpected
	}

	matched := s.passwordsMatch(hashedPassword, password)
	if matched {
		log.Println("password matched...")
		log.Println("userID we got at time of login : ", userID)
		newJwt, err := s.createJWt(userID)
		if err != nil {
			log.Println("error in generating jwt")
			return "", err
		}
		err = s.cache.SetJwtInCache(newJwt, userID)
		if err != nil {
			log.Println("error in setting jwt and userID")
			return "", err
		}
		log.Println("new jwt : ", newJwt)
		return newJwt, nil
	}
	log.Println("password mismatched...")
	return "", ErrUserNotAuthorized
}

func (s *svc) ChangePassword(ctx context.Context, user model.User, newPassword string) error {
	if user.Email == "" || user.PasswordHash == "" || newPassword == "" {
		log.Println("necessary field in missing...")
		return ErrBadRequest
	}
	log.Println("start.....")
	newHashed, err := s.hashPassword(newPassword)
	if err != nil {
		log.Println("error in hashing new password..", err)
		return ErrUnexpected
	}
	log.Println("new password hashed")
	err = s.dao.UpdatePassword(ctx, user, newHashed)
	if err != nil {
		log.Println("error in updating password ..", err)
		return ErrUnexpected
	}
	log.Println("password successfully changed...")
	return nil
}

var sameMoodGenres = func() map[string][]string {
	return map[string][]string{
		"sad":   []string{"depressed", "thriller", "sad", "separation"},
		"happy": []string{"family", "action", "happy", "drama", "motivation"},
	}
}

var oppositeMoodGenres = func() map[string][]string {
	return map[string][]string{
		"sad":   []string{"family", "love", "mind-refreshing", "motivational", "relaxing", "animation", "comedy", "fun"},
		"happy": []string{"alone", "aggression", "cringe", "depressing", "horror"},
	}
}

func (s *svc) getMoodMatcherGenre(moods []string) []string {
	var resultingTags []string
	for _, mood := range moods {
		resultingTags = append(resultingTags, sameMoodGenres()[mood]...)
	}
	return resultingTags
}

func (s *svc) getMoodChangerGenre(moods []string) []string {
	var resultingTags []string
	for _, mood := range moods {
		resultingTags = append(resultingTags, oppositeMoodGenres()[mood]...)
	}
	return resultingTags
}

func (s *svc) GetMoviesAccordingToUserMood(ctx context.Context, jwt string, mood []string) ([]model.Movie, error) {
	if jwt == "" {
		log.Println("jwt is empty")
		return nil, ErrMissingImportantField
	}
	userID, err := s.cache.GetUserIDFromJwt(jwt)
	if err != nil {
		log.Println("error in getting userId from jwt", err)
		return nil, ErrBadRequest
	}
	err = s.dao.UpdateUserMood(ctx, userID, mood)
	if err != nil {
		log.Println("error in updating user mood")
	}

	tags := s.getMoodMatcherGenre(mood)
	log.Println("tags matching mood ", mood, " are : ", tags)
	if len(tags) == 0 {
		var movies []model.Movie
		for _, m := range mood {
			movie, err := s.sdk.GetMovieByKeyword(ctx, m)
			if err != nil {
				log.Println("error in getting movie of mood : ", m)
				continue
			}
			movies = append(movies, movie...)
		}
		return movies, nil

	}
	movies, err := s.dao.GetMoviesByTags(ctx, tags)
	if err != nil {
		log.Println("error in getting movies by tags")
		return nil, err
	}
	log.Println("successfully got movies")
	return movies, nil
}

func (s *svc) GetMoviesOppositeToUserMood(ctx context.Context, jwt string, mood []string) ([]model.Movie, error) {
	if jwt == "" {
		log.Println("jwt is empty")
		return nil, ErrMissingImportantField
	}
	userID, err := s.cache.GetUserIDFromJwt(jwt)
	if err != nil {
		log.Println("error in getting userId from jwt", err)
		return nil, ErrBadRequest
	}
	err = s.dao.UpdateUserMood(ctx, userID, mood)
	if err != nil {
		log.Println("error in updating user mood")
	}

	tags := s.getMoodChangerGenre(mood)
	log.Println("tags matching mood ", mood, " are : ", tags)
	if len(tags) == 0 {
		var movies []model.Movie
		for _, m := range mood {
			movie, err := s.sdk.GetMovieByKeyword(ctx, m)
			if err != nil {
				log.Println("error in getting movie of mood : ", m)
				continue
			}
			movies = append(movies, movie...)
		}
		return movies, nil

	}
	movies, err := s.dao.GetMoviesByTags(ctx, tags)
	if err != nil {
		log.Println("error in getting movies by tags")
		return nil, err
	}
	log.Println("successfully got movies")
	return movies, nil
}

func (s *svc) updateMoodOfUser(ctx context.Context, userID string, mood []string) error {
	if userID == "" || len(mood) == 0 {
		log.Println("empty fields : ", " userID : ", userID, " mood : ", mood)
		return ErrMissingImportantField
	}
	err := s.dao.UpdateUserMood(ctx, userID, mood)
	if err != nil {
		log.Println("error in updating mood of user : ", userID, " in svc : ", err)
		return err
	}
	log.Println("mood updated successfully")
	return nil
}

func (s *svc) UpdateUserWatchedMovies(ctx context.Context, jwt string, movieID string) error {
	if jwt == "" || movieID == "" {
		log.Println("empty fields :  jwt : ", jwt, " movieID : ", movieID)
		return ErrMissingImportantField
	}
	userID, err := s.cache.GetUserIDFromJwt(jwt)
	if err != nil {
		log.Println("error in getting userID from jwt in svc : ", err)
		return err
	}
	err = s.dao.UpdateUserWatchedMovies(ctx, userID, movieID)
	if err != nil {
		log.Println("error in updating watched movies of user : ", userID, " : ", userID)
		return err
	}
	return nil
}

func (s *svc) UpdateWatchedMovieByMovieID(ctx context.Context, jwt string, movieID string) error {
	if jwt == "" || movieID == "" {
		log.Println("empty fields :  jwt : ", jwt, " movieID : ", movieID)
		return ErrMissingImportantField
	}
	userID, err := s.cache.GetUserIDFromJwt(jwt)
	if err != nil {
		log.Println("error in getting userID from jwt in svc : ", err)
		return err
	}

	movie, err := s.sdk.GetMovieByID(movieID)
	if err != nil {
		log.Println("error in getting movie from movieID from sdk")
		return err
	}

	mid, err := s.dao.AddMovie(ctx, movie)
	if err != nil {
		log.Println("error in adding movie in db")
		return err
	}

	err = s.dao.UpdateUserWatchedMovies(ctx, userID, mid)
	if err != nil {
		log.Println("error in updating watch movies")
		return err
	}
	return nil
}

func (s *svc) ForgotPassword(ctx context.Context, user model.User) error {
	if user.Email == "" {
		log.Println("email is missing : ", ErrBadRequest)
		return ErrBadRequest
	}
	exist, err := s.dao.CheckEmailExist(ctx, user)
	if err != nil {
		log.Println("email exist check internal error : ", err)
		return ErrUnexpected
	}
	if !exist {
		log.Println("email doesn't exist: ")
		return ErrNoData
	}
	otp, err := s.GenerateOtp()
	if err != nil {
		log.Println("error in generating otp : ", err)
		return err
	}
	err = s.mail.SendMail(user, otp)
	if err != nil {
		log.Println("error in sending mail at svc : ", err)
		return err
	}
	log.Println("mail sent successfully")
	err = s.cache.SetInCache(user.Email, otp)
	if err != nil {
		log.Println("error in setting otp cache: ", err)
		return err
	}
	return nil
}

func (s *svc) ProcessOtp(user model.User, otp string) (bool, error) {
	if user.Email == "" || otp == "" {
		log.Println("missing necessary field: ", ErrMissingImportantField)
		return false, ErrMissingImportantField
	}
	cachedOtp, err := s.cache.GetFromCache(user.Email)
	if err != nil {
		log.Println("error in finding email: ", err)
		return false, ErrUnexpected
	}
	log.Println("cachedOtp : ", cachedOtp, " otp : ", otp)
	if otp == cachedOtp {
		log.Println("success in verifying otp")
		return true, nil
	}
	return false, nil
}

func (s *svc) GetUserProfile(ctx context.Context, jwt string) (user model.User, err error) {
	if jwt == "" {
		log.Println("jwt is empty")
		err = ErrMissingImportantField
		return
	}
	userID, err := s.cache.GetUserIDFromJwt(jwt)
	if err != nil {
		log.Println("error in getting userId from jwt", err)
		return
	}
	user, err = s.dao.GetUserProfile(ctx, userID)

	if err != nil {
		log.Println("error in getting user from userID")
		return
	}
	var movies []model.Movie
	for _, mid := range user.MoviesWatched {
		movie, err := s.dao.GetMovieByMovieID(ctx, mid)
		if err != nil {
			log.Println("error in getting movie of mid : ", mid)
			continue
		}
		movies = append(movies, movie)
	}
	user.MoviesWatchedInformation = movies
	return
}

func (s *svc) GetWatchLater(ctx context.Context, jwt string) ([]model.Movie, error) {
	if jwt == "" {
		log.Println("userId is empty")
		return nil, ErrMissingImportantField
	}
	userId, err := s.cache.GetUserIDFromJwt(jwt)
	if err != nil {
		log.Println("error in getting user from jwt")
		return nil, ErrBadRequest
	}
	movies, err := s.dao.GetWatchLater(ctx, userId)
	if err != nil {
		log.Println("error in getting watchLater")
		return nil, err
	}
	log.Println("successfully got watchLater")
	return movies, nil
}

func (s *svc) DeleteWatchLater(ctx context.Context, jwt string, movieID string) error {
	if jwt == "" {
		log.Println("userId is empty")
		return ErrMissingImportantField
	}
	userId, err := s.cache.GetUserIDFromJwt(jwt)
	if err != nil {
		log.Println("error in getting user from jwt")
		return ErrBadRequest
	}
	log.Println("userID  : ", userId)
	err = s.dao.DeleteWatchLater(ctx, userId, movieID)
	if err != nil {
		log.Println("error in deleting movie watchLater")
		return err
	}
	return nil
}

func (s *svc) UpdateWatchLater(ctx context.Context, jwt string, movieID string, isMovieDB bool, showType string) error {
	if jwt == "" {
		log.Println("userId is empty")
		return ErrMissingImportantField
	}
	userId, err := s.cache.GetUserIDFromJwt(jwt)
	if err != nil {
		log.Println("error in getting user from jwt")
		return ErrBadRequest
	}
	log.Println("userID  : ", userId)
	var movie model.Movie
	if isMovieDB {
		movie, err = s.sdk.GetMovieByID(movieID)
		if err != nil {
			log.Println("error in getting movie from movieID", err)
			return err
		}
	} else {
		movie, err = s.dao.GetMovieByMovieID(ctx, movieID)
		if err != nil {
			log.Println("error in getting movie from db by movieID", err)
			return err
		}
	}
	err = s.dao.AddMovieToWatchLater(ctx, userId, movie)
	if err != nil {
		log.Println("error in repo layer : ", err)
		return err
	}
	log.Println("successfully updated")
	return nil
}
