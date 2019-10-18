package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/rasha108/apiCargoRest.git/internal/app/store/sqlstore"

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

	config := api.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sqlx.Connect("postgres", config.DatabaseURL)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer db.Close()

	store := sqlstore.New(db)
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	srv := api.NewServer(store, sessionStore)

	err = http.ListenAndServe(config.BindAddr, srv)
	if err != nil {
		log.Fatal(err)
	}
}
