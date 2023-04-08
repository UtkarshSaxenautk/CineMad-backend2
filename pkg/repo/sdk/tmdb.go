package sdk

import (
	"authentication-ms/pkg/model"
	"authentication-ms/pkg/svc"
	"context"
	"github.com/ryanbradynd05/go-tmdb"
	"log"
	"time"
)

func InitTMDB() *tmdb.TMDb {
	config := tmdb.Config{
		APIKey:   "84ba864072e47e9b88790e9a06781353",
		Proxies:  nil,
		UseProxy: false,
	}

	return tmdb.Init(config)
}

func (s *sdk) GetMovieByKeyword(ctx context.Context, tag string) ([]model.Movie, error) {
	if tag == "" {
		log.Println("missing tag")
		return nil, svc.ErrBadRequest
	}
	var movies []model.Movie
	releaseDateGTE := time.Date(2001, time.January, 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")

	options := map[string]string{
		"sort_by":          "release_date.asce", // Sort movies by release date in descending order
		"include_adult":    "false",             // Exclude adult movies
		"release_date.gte": releaseDateGTE,      // Set the minimum release date
	}
	res, err := s.tmDb.SearchCollection(tag, options)
	if err != nil {
		log.Println("error in getting movies for keyWord")
		return nil, err
	}
	log.Println(res.Results)
	for _, x := range res.Results {

		url, err := s.tmDb.GetMovieImages(x.ID, nil)
		if err != nil {
			log.Println("error in getting url of movie poster", x.ID, " : ", err)

		}
		var movie model.Movie
		if len(url.Posters) <= 0 {
			movie = model.Movie{
				MovieId: x.ID,
				Name:    x.Name,
			}
		} else {
			movie = model.Movie{
				MovieId: x.ID,
				Name:    x.Name,
				Url:     "https://image.tmdb.org/t/p/w500/" + url.Posters[0].FilePath,
			}
		}

		movies = append(movies, movie)
	}

	return movies, nil
}
