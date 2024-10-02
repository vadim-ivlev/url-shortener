package db

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/vadim-ivlev/url-shortener/internal/config"
)

// DB - пул соединений с базой данных
var DB *sqlx.DB = nil

// CreateDBPool - создает пул соединений с базой данных
func CreateDBPool() (err error) {
	DB, err = sqlx.Connect("postgres", config.Params.DatabaseDSN)
	return err
}

// ConnectToDatabase - Ожидает соединения с базой данных повторяя попытки в случае неудачи.
// numAttempts - количество попыток
func ConnectToDatabase(numAttempts int) {
	for i := 1; i <= numAttempts; i++ {
		err := CreateDBPool()
		if err == nil {
			log.Info().Msg("Connected to DB")
			return
		}
		log.Warn().Err(err).Msg(fmt.Sprintf("Waiting for db connection. Attempt # %d", i))
		time.Sleep(time.Second)
	}
	log.Error().Msg("Failed to connect to DB")
}
