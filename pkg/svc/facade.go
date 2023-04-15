package svc

import (
	"authentication-ms/pkg/model"
	"context"
)

//go:generate mockgen -destination=mock_svc.go -package=svc . Cache
type Cache interface {
	SetInCache(email string, otp string) error
	GetFromCache(email string) (string, error)
	SetJwtInCache(jwt string, userID string) error
	GetUserIDFromJwt(jwt string) (string, error)
}

//go:generate mockgen -destination=mock_svc.go -package=svc . SVC
type SVC interface {
	Signup(ctx context.Context, user model.User) error
	SignIn(ctx context.Context, email string, password string) (string, error)
	ChangePassword(ctx context.Context, user model.User, newPassword string) error
	ForgotPassword(ctx context.Context, user model.User) error
	ProcessOtp(user model.User, otp string) (bool, error)
	GetMoviesByTag(ctx context.Context, tag string) ([]model.Movie, error)
	GetMoviesByTags(ctx context.Context, tags []string) ([]model.Movie, error)
	AddMovieInDB(ctx context.Context, movie model.Movie) error
	GetMoviesByTagsFromDB(ctx context.Context, tags []string) ([]model.Movie, error)
	UpdateUserWatchedMovies(ctx context.Context, jwt string, movieID string) error
}

//go:generate mockgen -destination=mock_dao.go -package=svc . Dao
type Dao interface {
	CreateUser(ctx context.Context, user model.User) error
	CheckEmailAndUserName(ctx context.Context, user model.User) (emailExist bool, usernameExist bool, err error)
	GetUser(ctx context.Context, email string) (string, string, error)
	UpdatePassword(ctx context.Context, user model.User, nPass string) error
	AddMovie(ctx context.Context, movie model.Movie) error
	GetMoviesByTags(ctx context.Context, tags []string) ([]model.Movie, error)
	UpdateUserMood(ctx context.Context, userId string, mood string) error
	UpdateUserWatchedMovies(ctx context.Context, userID string, movieID string) error
	CheckEmailExist(ctx context.Context, user model.User) (bool, error)
}

//go:generate mockgen -destination=mock_svc.go -package=svc . SVC
type Mail interface {
	SendMail(user model.User, otp string) error
}

type Sdk interface {
	GetMovie(ctx context.Context, tag string) ([]model.Movie, error)
	GetMovieByKeyword(ctx context.Context, tag string) ([]model.Movie, error)
}
