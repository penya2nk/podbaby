package commands

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/danjac/podbaby/api"
	"github.com/danjac/podbaby/database"
	"github.com/danjac/podbaby/feedparser"
	"github.com/jmoiron/sqlx"
)

// should be settings
const (
	defaultStaticURL = "/static/"
	defaultStaticDir = "./static/"
	devStaticURL     = "http://localhost:8080/static/"
)

// Serve runs the webserver
func Serve(url string, port int, secretKey, env string) {

	log := logrus.New()

	log.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	}

	log.Info("Starting web service...")

	db := database.New(sqlx.MustConnect("postgres", url))
	defer db.Close()

	var dynamicContentURL string
	if env == "dev" {
		dynamicContentURL = devStaticURL
	} else {
		dynamicContentURL = defaultStaticURL
	}

	api := api.New(db, log, &api.Config{
		StaticURL:         defaultStaticURL,
		DynamicContentURL: dynamicContentURL,
		StaticDir:         defaultStaticDir,
		SecretKey:         secretKey,
	})

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), api.Handler()); err != nil {
		panic(err)
	}

}

// Fetch retrieves latest podcasts
func Fetch(url string) {

	db := database.New(sqlx.MustConnect("postgres", url))
	defer db.Close()

	log := logrus.New()

	log.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	}

	f := feedparser.New(db, log)
	if err := f.FetchAll(); err != nil {
		panic(err)
	}

}
