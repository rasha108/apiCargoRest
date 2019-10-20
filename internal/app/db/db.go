package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

//
//   Пока нигде не использовал, хочу заюзать
//
//

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
	DSNtpl := "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s"
	DSN := fmt.Sprintf(DSNtpl,
		params.host,
		params.port,
		params.user,
		params.pass,
		params.dbName,
	)

	logger.Tracef("connecting: host=%s pord=%d user=%s dbname=%s",
		params.host,
		params.port,
		params.user,
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
