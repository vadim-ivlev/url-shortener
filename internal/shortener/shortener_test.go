package shortener

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestShorten(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantKey string
	}{
		{
			name:    "empty value",
			args:    args{value: ""},
			wantKey: "811C9DC5",
		},
		{
			name:    "empty value 2",
			args:    args{value: ""},
			wantKey: "811C9DC5",
		},
		{
			name:    "google",
			args:    args{value: "https://www.google.com"},
			wantKey: "F870F1E9",
		},
		{
			name:    "youtube",
			args:    args{value: "https://www.youtube.com"},
			wantKey: "4AED1C05",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey := Shorten(tt.args.value)
			assert.Equal(t, tt.wantKey, gotKey)
		})
	}
}
