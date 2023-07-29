package db

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/logrusadapter"
	log "github.com/sirupsen/logrus"
)

type Store struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string

	db *sql.DB
}

func New(host string, port int, user string, password string, dbname string) *Store {
	s := Store{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Dbname:   dbname,
	}
	return &s
}

func (s *Store) Open() error {
	info := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		s.Host, s.Port, s.User, s.Password, s.Dbname)

	loggerdb := sqldblogger.OpenDriver(
		info,
		&pq.Driver{},
		logrusadapter.New(log.StandardLogger()),
	)

	s.db = loggerdb

	err := s.db.Ping()
	if err != nil {
		return err
	}
	return nil
}
