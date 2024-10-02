package db

import (
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/vadim-ivlev/url-shortener/internal/config"
)

func TestMain(m *testing.M) {
	config.ParseCommandLine()
	os.Exit(m.Run())
}

func TestCreateDBPool(t *testing.T) {
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
			if err := Connect(tt.args.numAttempts); (err != nil) != tt.wantErr {
				t.Errorf("ConnectToDatabase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
