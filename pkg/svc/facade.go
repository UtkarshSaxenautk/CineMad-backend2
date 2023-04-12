package svc

import (
	"authentication-ms/pkg/model"
	"context"
)

//go:generate mockgen -destination=mock_svc.go -package=svc . Cache
type Cache interface {
	SetInCache(email string, otp string) error
	GetFromCache(email string) (string, error)
}

//go:generate mockgen -destination=mock_svc.go -package=svc . SVC
type SVC interface {
	Signup(ctx context.Context, user model.User) error
	SignIn(ctx context.Context, email string, password string) (bool, error)
	ChangePassword(ctx context.Context, user model.User, newPassword string) error
	GetMoviesByTag(ctx context.Context, tag string) ([]model.Movie, error)
	GetMoviesByTags(ctx context.Context, tags []string) ([]model.Movie, error)
	AddMovieInDB(ctx context.Context, movie model.Movie) error
	GetMoviesByTagsFromDB(ctx context.Context, tags []string) ([]model.Movie, error)
}

//go:generate mockgen -destination=mock_dao.go -package=svc . Dao
type Dao interface {
	CreateUser(ctx context.Context, user model.User) error
	CheckEmailAndUserName(ctx context.Context, user model.User) (emailExist bool, usernameExist bool, err error)
	GetUser(ctx context.Context, email string) (string, error)
	UpdatePassword(ctx context.Context, user model.User, nPass string) error
	AddMovie(ctx context.Context, movie model.Movie) error
	GetMoviesByTags(ctx context.Context, tags []string) ([]model.Movie, error)
}

type Sdk interface {
	GetMovie(ctx context.Context, tag string) ([]model.Movie, error)
	GetMovieByKeyword(ctx context.Context, tag string) ([]model.Movie, error)
}
