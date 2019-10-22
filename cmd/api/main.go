package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/rasha108/apiCargoRest.git/internal/app/db"

	"github.com/sirupsen/logrus"

	"github.com/rasha108/apiCargoRest.git/internal/app/rabbitmq"

	"github.com/gorilla/sessions"
	"github.com/rasha108/apiCargoRest.git/internal/app/db/sqlstore"

	"github.com/rasha108/apiCargoRest.git/internal/app/api"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/config.yaml", "path to config file")
}

func main() {
	flag.Parse()

	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)

	conf, err := api.GetConfig(configPath)
	if err != nil {
		logger.Fatalf("get")
	}

	logger.Info("Start server...")

	dbConfig := conf.DbConfig

	connParams := db.NewConnectParams(
		dbConfig.Host,
		dbConfig.User,
		dbConfig.Pass,
		dbConfig.DbName,
		dbConfig.Port,
		dbConfig.MaxConnections,
	)

	connHandler, err := db.Connect(connParams, logger)
	if err != nil {
		logger.WithError(err).Fatal("db connect failed")
		return
	}
	defer func() {
		dbErr := db.Disconnect(connHandler, logger)
		if dbErr != nil {
			logger.WithError(dbErr).Fatal("db disconnect failed")
		}
	}()

	store := sqlstore.New(connHandler)
	sessionStore := sessions.NewCookieStore([]byte(conf.SessionKey))
	rabbitServer, err := rabbitmq.NewConnection(conf.MailConfig)
	if err != nil {
		logger.WithError(err).Error("create mail client failed")
		return
	}

	defer func() {
		err := rabbitServer.Close()
		if err != nil {
			logger.WithError(err).Error("rabbitmq disconnect failed")
		}
	}()

	srv := api.NewServer(store, sessionStore, conf, rabbitServer)

	bind := conf.APIConfig.Bind
	err = http.ListenAndServe(bind, srv)
	if err != nil {
		logger.WithError(err).Error("application aborted")
		return
	}
}
