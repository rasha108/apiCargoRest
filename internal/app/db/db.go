package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	// load PostgreSQL driver
	_ "gopkg.in/jackc/pgx.v3/stdlib"
)

const (
	// ConnectionType a driver name in connection args
	ConnectionType = "pgx"
)

// ConnectParams contains database connection parametrs
type ConnectParams struct {
	host           string
	user           string
	pass           string
	dbName         string
	port           int
	maxConnections int
}

// NewConnectParamas is a constructor for connectParams
func NewConnectParams(host, user, pass, dbName string, port, maxConnections int) *ConnectParams {
	return &ConnectParams{
		host:           host,
		user:           user,
		pass:           pass,
		dbName:         dbName,
		port:           port,
		maxConnections: maxConnections,
	}
}

// Connect creates database connection
// returns wrapper connection type
func Connect(params *ConnectParams, logger logrus.Logger) (*sqlx.DB, error) {
	DSNtpl := "host=%s dbname=%s sslmode=disable"
	DSN := fmt.Sprintf(DSNtpl,
		params.host,
		params.dbName,
	)
	logger.Tracef("connecting: host=%s dbname=%s",
		params.host,
		params.dbName,
	)
	connectionType := ConnectionType
	db, err := sqlx.Open(connectionType, DSN)
	if err != nil {
		logger.WithError(err).Warning("open connect failed")
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		logger.WithError(err).Warning("ping connections failed")
		return nil, err
	}

	db.SetMaxOpenConns(params.maxConnections)

	return db, nil
}

func Disconnect(rawConn *sqlx.DB, logger logrus.Logger) error {
	logger.Tracef("disconnecting db")
	err := rawConn.Close()
	if err != nil {
		logger.WithError(err).Warning("disconnect failed")
	}

	return err
}
