package transport

import (
	"authentication-ms/pkg/svc"
	"authentication-ms/pkg/transport/handler"
	"authentication-ms/pkg/transport/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

type Router struct {
	*mux.Router
	routePrefix string

	svc svc.SVC
}

func NewRouter(routePrefix string, svc svc.SVC) *Router {
	return &Router{mux.NewRouter(), routePrefix, svc}
}

func (r *Router) Initialize() *Router {
	(*r).Use(middleware.TraceMiddleware)

	r.HandleFunc("/healthcheck", healthCheck()).Methods(http.MethodGet)
	cf := (*r).PathPrefix("/authenticate").Subrouter()
	(*cf).HandleFunc("/health", handler.GetHealth(r.svc)).Methods(http.MethodGet)
	(*cf).HandleFunc("/signup", handler.Signup(r.svc)).Methods(http.MethodPost)
	(*cf).HandleFunc("/signIn", handler.SignIn(r.svc)).Methods(http.MethodGet)
	(*cf).HandleFunc("/changepassword", handler.ChangePassword(r.svc)).Methods(http.MethodPost)
	mf := (*r).PathPrefix("/movie").Subrouter()
	(*mf).HandleFunc("/getMovieByTag/{tag}", handler.GetMovies(r.svc)).Methods(http.MethodGet)
	(*mf).HandleFunc("/getMovieByTags", handler.GetMoviesByTags(r.svc)).Methods(http.MethodGet)
	(*mf).HandleFunc("/addMovie", handler.AddMovie(r.svc)).Methods(http.MethodPost)
	return r
}

func healthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello"))
	}
}
