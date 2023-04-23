package sdk

import (
	"authentication-ms/pkg/svc"
	"fmt"
	"github.com/ryanbradynd05/go-tmdb"
)

type sdk struct {
	movieDbUrl string
	tmDb       *tmdb.TMDb
}

//type movie struct {
//	Title       string `json:"title"`
//	ReleaseDate string `json:"release_date"`
//	Overview    string `json:"overview"`
//}
//
//// Response represents the structure of the API response
//type Response struct {
//	Results []movie `json:"results"`
//}
//
//type response struct {
//	Result model.Movie `json:"result"`
//}

func New() svc.Sdk {
	movieDBurl := GenerateMovieDbUrl()
	tmDb := InitTMDB()
	return &sdk{
		movieDbUrl: movieDBurl,
		tmDb:       tmDb,
	}
}

func GenerateMovieDbUrl() string {
	apiKey := "84ba864072e47e9b88790e9a06781353"
	baseURL := "https://api.themoviedb.org/3"
	url := fmt.Sprintf("%s/search/movie?api_key=%s", baseURL, apiKey)
	return url
}

//
//func (s *sdk) generateFindMovieUrl(id string) string {
//	apiKey := "84ba864072e47e9b88790e9a06781353"
//	baseURL := "https://api.themoviedb.org/3"
//	url := fmt.Sprintf("%s/movie/%s?api_key=%s", baseURL, id, apiKey)
//	return url
//}
//
//func (s *sdk) generateFindTvUrl(id string) string {
//	apiKey := "84ba864072e47e9b88790e9a06781353"
//	baseURL := "https://api.themoviedb.org/3"
//	url := fmt.Sprintf("%s/tv/%s?api_key=%s", baseURL, id, apiKey)
//	return url
//}
//
//func (s *sdk) generateFindPersonUrl(id string) string {
//	apiKey := "84ba864072e47e9b88790e9a06781353"
//	baseURL := "https://api.themoviedb.org/3"
//	url := fmt.Sprintf("%s/person/%s?api_key=%s", baseURL, id, apiKey)
//	return url
//}

//func (s *sdk) GetMovie(ctx context.Context, tag string) ([]model.Movie, error) {
//	url := s.movieDbUrl + fmt.Sprintf("&query=%s", tag)
//
//	// Send a GET request to the API
//	resp, err := http.Get(url)
//	if err != nil {
//		log.Println("error in getting movies from tag in movieDb....", err)
//		return nil, err
//	}
//	defer func(Body io.ReadCloser) {
//		err := Body.Close()
//		if err != nil {
//			log.Println("error in closing body in movieDb", err)
//		}
//	}(resp.Body)
//
//	// Parse the response body into a Response struct
//	var response Response
//	err = json.NewDecoder(resp.Body).Decode(&response)
//	if err != nil {
//		return nil, err
//	}
//
//	// Print the movie titles from the API response
//	var res []model.Movie
//	fmt.Println("Movies with mood tag", tag+":")
//	for _, movie := range response.Results {
//		tempMovie := model.Movie{
//			Name:     movie.Title,
//			OverView: movie.Overview,
//		}
//		res = append(res, tempMovie)
//		fmt.Println(movie.Title)
//	}
//	return res, nil
//}
