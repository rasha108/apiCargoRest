package api

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/rasha108/apiCargoRest.git/internal/app/store/sqlstore"

	_ "github.com/lib/pq"
)

func Start(config *Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}

	defer db.Close()

	store := sqlstore.New(db)
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	srv := NewServer(store, sessionStore)

	return http.ListenAndServe(config.BindAddr, srv)
}

func newDB(databaseURL string) (*sqlx.DB, error) {
	//Connect внутри себя вызывает ping
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	return db, nil
}
