package transport

import (
	"authentication-ms/pkg/svc"
	"authentication-ms/pkg/transport/handler"
	"authentication-ms/pkg/transport/middleware"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
)

type Router struct {
	*mux.Router
	routePrefix string
	svc         svc.SVC
}

func NewRouter(routePrefix string, svc svc.SVC) *Router {
	return &Router{mux.NewRouter(), routePrefix, svc}
}

func (r *Router) Initialize() *Router {
	(*r).Use(middleware.TraceMiddleware)
	r.HandleFunc("/healthcheck", healthCheck()).Methods(http.MethodGet)
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Allow requests from any origin
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},
		AllowedHeaders: []string{"*"}, // Allow all headers
	})
	r.Use(c.Handler)
	//r.Use(func(next http.Handler) http.Handler {
	//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//		// Allow all origins
	//		w.Header().Set("Access-Control-Allow-Origin", "*")
	//		// Allow all headers
	//		w.Header().Set("Access-Control-Allow-Headers", "*")
	//		// Allow all methods
	//		w.Header().Set("Access-Control-Allow-Methods", "*")
	//		next.ServeHTTP(w, r)
	//	})
	//})
	cf := (*r).PathPrefix("/authenticate").Subrouter()
	(*cf).Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Allow all origins
			w.Header().Set("Access-Control-Allow-Origin", "*")
			// Allow all headers
			w.Header().Set("Access-Control-Allow-Headers", "*")
			// Allow all methods
			w.Header().Set("Access-Control-Allow-Methods", "*")
			next.ServeHTTP(w, r)
		})
	})
	(*cf).HandleFunc("/health", handler.GetHealth(r.svc)).Methods(http.MethodGet)
	(*cf).HandleFunc("/signup", handler.Signup(r.svc)).Methods(http.MethodPost, http.MethodOptions)
	(*cf).HandleFunc("/signIn", handler.SignIn(r.svc)).Methods(http.MethodPost, http.MethodOptions)
	(*cf).HandleFunc("/changepassword", handler.ChangePassword(r.svc)).Methods(http.MethodPost)
	(*cf).HandleFunc("/forgot", handler.ForgotPassword(r.svc)).Methods(http.MethodPost)
	(*cf).HandleFunc("/otp/{email}", handler.ProcessOtp(r.svc)).Methods(http.MethodPost)
	uf := (*r).PathPrefix("/user").Subrouter()
	(*uf).HandleFunc("/getMovieAccMood", handler.GetMovieAccordingToMood(r.svc)).Methods(http.MethodPost, http.MethodOptions)
	(*uf).HandleFunc("/getMovieOppMood", handler.GetMovieOppositeToMood(r.svc)).Methods(http.MethodPost, http.MethodOptions)
	(*uf).HandleFunc("/getWatchLater", handler.GetWatchLater(r.svc)).Methods(http.MethodPost, http.MethodOptions)
	(*uf).HandleFunc("/deleteWatchLater", handler.DeleteWatchLaterMovie(r.svc)).Methods(http.MethodPost, http.MethodOptions)
	(*uf).HandleFunc("/updateWatchedMovie", handler.UpdateUserWatchedMovie(r.svc)).Methods(http.MethodPost, http.MethodOptions)
	(*uf).HandleFunc("/updateWatchLater", handler.UpdateUserWatchLaterMovie(r.svc)).Methods(http.MethodPost, http.MethodOptions)
	(*uf).HandleFunc("/getprofile", handler.GetUserProfile(r.svc)).Methods(http.MethodPost, http.MethodOptions)
	mf := (*r).PathPrefix("/movie").Subrouter()
	(*mf).Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Allow all origins
			w.Header().Set("Access-Control-Allow-Origin", "*")
			// Allow all headers
			w.Header().Set("Access-Control-Allow-Headers", "*")
			// Allow all methods
			w.Header().Set("Access-Control-Allow-Methods", "*")
			next.ServeHTTP(w, r)
		})
	})
	(*mf).HandleFunc("/getMovieByTag/{tag}", handler.GetMovies(r.svc)).Methods(http.MethodGet)
	(*mf).HandleFunc("/getMovieByTags", handler.GetMoviesByTags(r.svc)).Methods(http.MethodGet)
	(*mf).HandleFunc("/addMovie", handler.AddMovie(r.svc)).Methods(http.MethodPost)
	//corsOpts := cors.New(cors.Options{
	//	AllowedOrigins: []string{"http://localhost:1234"}, //you service is available and allowed for this base url
	//	AllowedMethods: []string{
	//		http.MethodGet, //http methods for your app
	//		http.MethodPost,
	//		http.MethodPut,
	//		http.MethodPatch,
	//		http.MethodDelete,
	//		http.MethodOptions,
	//		http.MethodHead,
	//	},
	//
	//	AllowedHeaders: []string{
	//		"*", //or you can your header key values which you are using in your application
	//
	//	},
	//})
	//

	return r
}

func healthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello"))
	}
}
