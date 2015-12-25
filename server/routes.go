package server

import (
	"database/sql"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/justinas/nosurf"
	"net/http"
)

func (s *Server) Router() *mux.Router {

	router := mux.NewRouter()
	// static routes

	router.PathPrefix(s.Config.StaticURL).Handler(
		http.StripPrefix(s.Config.StaticURL,
			http.FileServer(http.Dir(s.Config.StaticDir))))

	// front page
	router.HandleFunc("/", s.indexPage)

	// API

	api := router.PathPrefix("/api/").Subrouter()

	// authentication

	auth := api.PathPrefix("/auth/").Subrouter()

	auth.HandleFunc("/currentuser/", s.getCurrentUser).Methods("GET")
	auth.HandleFunc("/login/", s.login).Methods("POST")
	auth.HandleFunc("/signup/", s.signup).Methods("POST")
	auth.HandleFunc("/logout/", s.logout).Methods("DELETE")

	// podcasts

	podcasts := api.PathPrefix("/podcasts/").Subrouter()

	podcasts.HandleFunc("/latest/", s.getLatestPodcasts).Methods("GET")

	return router
}

func (s *Server) Abort(w http.ResponseWriter, r *http.Request, err error) {
	logger := s.Log.WithFields(logrus.Fields{
		"URL":    r.URL,
		"Method": r.Method,
		"Error":  err,
	})
	if err == sql.ErrNoRows {
		http.Error(w, err.Error(), http.StatusNotFound)
		logger.Debug("Not found:" + err.Error())
		return
	}

	var msg string

	switch e := err.(error).(type) {
	case Error:
		msg = "HTTP Error"
		http.Error(w, e.Error(), e.Status())
	default:
		msg = "Internal Server Error"
		http.Error(w, "Sorry, an error occurred", http.StatusInternalServerError)
	}
	logger.Error(msg)
}

func (s *Server) indexPage(w http.ResponseWriter, r *http.Request) {
	csrfToken := nosurf.Token(r)
	ctx := map[string]string{
		"staticURL": s.Config.StaticURL,
		"csrfToken": csrfToken,
	}
	s.Render.HTML(w, http.StatusOK, "index", ctx)
}