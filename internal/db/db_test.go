package db

import (
	"context"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/vadim-ivlev/url-shortener/internal/apptypes"
	"github.com/vadim-ivlev/url-shortener/internal/config"
)

// skipCI skips tests in CI environment
func skipCI(t *testing.T) {
	if os.Getenv("CI") != "" {
		log.Info().Msg("Skipping testing in CI environment")
		t.Skip("Skipping testing in CI environment")
	}
}

func TestMain(m *testing.M) {
	os.Chdir("../../")
	os.Setenv("DATABASE_DSN", "postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable")
	config.ParseCommandLine()
	config.ParseEnv()
	os.Exit(m.Run())
}

func TestDeleteKeys(t *testing.T) {
	skipCI(t)
	err := Connect()
	assert.NoError(t, err)

	err = Clear()
	assert.NoError(t, err)

	err = AddRecord(apptypes.URLShortener{Idx: 10, ShortID: "short_id10", OriginalURL: "original_url10", UserID: "user_t", Deleted: 0})
	assert.NoError(t, err)
	err = AddRecord(apptypes.URLShortener{Idx: 11, ShortID: "short_id11", OriginalURL: "original_url11", UserID: "user_t", Deleted: 0})
	assert.NoError(t, err)
	err = AddRecord(apptypes.URLShortener{Idx: 12, ShortID: "short_id12", OriginalURL: "original_url12", UserID: "user_t", Deleted: 0})
	assert.NoError(t, err)

	err = DeleteRecords(context.Background(), "user_t", []any{"short_id10", "short_id11"})
	assert.NoError(t, err)

	record, _ := GetRecordByShortID(context.Background(), "short_id10")
	assert.EqualValues(t, record.Deleted, 1)

	record, _ = GetRecordByShortID(context.Background(), "short_id11")
	assert.EqualValues(t, record.Deleted, 1)

	record, _ = GetRecordByShortID(context.Background(), "short_id12")
	assert.EqualValues(t, record.Deleted, 0)

}
