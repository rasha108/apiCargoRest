package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/rasha108/apiCargoRest.git/internal/app/rabbitclient"

	"github.com/gorilla/sessions"
	"github.com/rasha108/apiCargoRest.git/internal/app/db/sqlstore"

	"github.com/jmoiron/sqlx"

	"github.com/BurntSushi/toml"
	"github.com/rasha108/apiCargoRest.git/internal/app/api"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")
}

func main() {
	flag.Parse()
	logger := logrus.Logger{}

	config := api.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sqlx.Connect("postgres", config.DatabaseURL)
	if err != nil {
		logger.WithError(err).Error("db connect failed")
		return
	}

	defer db.Close()

	store := sqlstore.New(db)
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))

	mailServer, err := rabbitclient.NewConnection(config.MailConfig)
	if err != nil {
		logger.WithError(err).Error("create mail client failed")
		return
	}

	srv := api.NewServer(store, sessionStore, mailServer)

	err = http.ListenAndServe(config.BindAddr, srv)
	if err != nil {
		logger.WithError(err).Error("application aborted")
		return
	}
}
