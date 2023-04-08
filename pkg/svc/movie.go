package svc

import (
	"authentication-ms/pkg/model"
	"context"
	"log"
)

func (s *svc) GetMoviesByTag(ctx context.Context, tag string) ([]model.Movie, error) {
	if tag == "" {
		log.Println("tag is empty")
		return nil, ErrBadRequest
	}
	movies, err := s.sdk.GetMovie(ctx, tag)
	if err != nil {
		log.Println("error in getting movie by tag ", err)
		return nil, err
	}
	return movies, nil
}

func (s *svc) GetMoviesByTags(ctx context.Context, tags []string) ([]model.Movie, error) {
	if tags == nil {
		log.Println("tags are empty..")
		return nil, ErrBadRequest
	}
	var movies []model.Movie
	for _, tag := range tags {
		res, err := s.sdk.GetMovieByKeyword(ctx, tag)
		if err != nil {
			log.Println("error in getting movies for tag : ", tag)
			continue
		}
		movies = append(movies, res...)
	}
	return movies, nil
}

func (s *svc) AddMovieInDB(ctx context.Context, movie model.Movie) error {
	if movie.Name == "" || len(movie.Tags) <= 0 || movie.Url == "" || movie.ImageUrl == "" {
		log.Println("important field missing ")
		return ErrBadRequest
	}
	err := s.dao.AddMovie(ctx, movie)
	if err != nil {
		log.Println("error in adding movie to db ", err)
		return ErrUnexpected
	}
	log.Println("movie: ", movie, "inserted successfully")
	return nil
}
