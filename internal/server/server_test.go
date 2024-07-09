package server

import (
	"os"
	"testing"

	"github.com/vadim-ivlev/url-shortener/internal/config"
)

func TestMain(m *testing.M) {
	config.ParseCommandLine()
	os.Exit(m.Run())
}

func TestServe(t *testing.T) {
	tests := []struct {
		name string
	}{
		// No test cases for this function
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ServeChi(config.Address)
		})
	}
}
