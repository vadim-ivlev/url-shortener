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

func TestCreateDBPool(t *testing.T) {
	skipCI(t)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "TestCreateDBPool",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreatePool(); (err != nil) != tt.wantErr {
				t.Errorf("CreateDBPool() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConnectToDatabase(t *testing.T) {
	skipCI(t)

	type args struct {
		numAttempts int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "TestConnectToDatabase 3 attempts",
			args:    args{numAttempts: 3},
			wantErr: false,
		},
		{
			name:    "TestConnectToDatabase 0 attempts",
			args:    args{numAttempts: 0},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := TryToConnect(tt.args.numAttempts); (err != nil) != tt.wantErr {
				t.Errorf("ConnectToDatabase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteKeys(t *testing.T) {
	skipCI(t)
	err := TryToConnect(1)
	assert.NoError(t, err)

	err = Clear()
	assert.NoError(t, err)

	err = AddRecord(apptypes.URLShortener{Idx: 10, ShortID: "short_id10", OriginalURL: "original_url10", UserID: "user_t", Deleted: 0})
	assert.NoError(t, err)
	err = AddRecord(apptypes.URLShortener{Idx: 11, ShortID: "short_id11", OriginalURL: "original_url11", UserID: "user_t", Deleted: 0})
	assert.NoError(t, err)
	err = AddRecord(apptypes.URLShortener{Idx: 12, ShortID: "short_id12", OriginalURL: "original_url12", UserID: "user_t", Deleted: 0})
	assert.NoError(t, err)

	err = DeleteShortIDs(context.Background(), "user_t", []any{"short_id10", "short_id11"})
	assert.NoError(t, err)

	record, _ := GetByShortID(context.Background(), "short_id10")
	assert.EqualValues(t, record.Deleted, 1)

	record, _ = GetByShortID(context.Background(), "short_id11")
	assert.EqualValues(t, record.Deleted, 1)

	record, _ = GetByShortID(context.Background(), "short_id12")
	assert.EqualValues(t, record.Deleted, 0)

}
