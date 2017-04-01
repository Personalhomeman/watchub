package main

import (
	"net/http"
	"time"

	"github.com/apex/httplog"
	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	"github.com/caarlos0/watchub/config"
	"github.com/caarlos0/watchub/datastore/database"
	"github.com/caarlos0/watchub/oauth"
	"github.com/caarlos0/watchub/scheduler"
	"github.com/caarlos0/watchub/shared/pages"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	log.SetHandler(text.Default)
	log.SetLevel(log.InfoLevel)
	log.Info("starting up...")

	var config = config.Get()
	var db = database.Connect(config.DatabaseURL)
	defer func() { _ = db.Close() }()
	var store = database.NewDatastore(db)

	// oauth
	var oauth = oauth.New(store, config)

	// schedulers
	var scheduler = scheduler.New(config, store, oauth)
	scheduler.Start()
	defer scheduler.Stop()

	// routes
	var mux = mux.NewRouter()
	mux.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.Methods("GET").Path("/").
		HandlerFunc(pages.New(config, "index").Handler)
	mux.Methods("GET").Path("/donate").
		HandlerFunc(pages.New(config, "donate").Handler)
	mux.Methods("GET").Path("/support").
		HandlerFunc(pages.New(config, "support").Handler)

	var loginMux = mux.Methods("GET").PathPrefix("/login").Subrouter()
	loginMux.Path("").HandlerFunc(oauth.LoginHandler())
	loginMux.Path("/callback").HandlerFunc(oauth.LoginCallbackHandler())

	// RUN!
	var server = &http.Server{
		Handler:      httplog.New(handlers.CompressHandler(mux)),
		Addr:         ":" + config.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.WithField("port", config.Port).Info("started")
	log.WithError(server.ListenAndServe()).Error("failed to start up server")
}
