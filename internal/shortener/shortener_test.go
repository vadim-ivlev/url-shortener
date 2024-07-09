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
			wantKey: "short1",
		},
		{
			name:    "non-empty value",
			args:    args{value: "https://www.google.com"},
			wantKey: "short2",
		},
		{
			name:    "another non-empty value",
			args:    args{value: "https://www.youtube.com"},
			wantKey: "short3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey := Shorten(tt.args.value)
			assert.Equal(t, tt.wantKey, gotKey)
		})
	}
}
