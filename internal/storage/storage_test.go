package storage

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	Create()
	os.Exit(m.Run())
}

type args struct {
	key   string
	value string
}

var tests = []struct {
	name   string
	args   args
	want   string
	exists bool
}{
	{
		name:   "empty key & value",
		args:   args{key: "", value: ""},
		want:   "",
		exists: true,
	},
	{
		name:   "google",
		args:   args{key: "F870F1E9", value: "https://www.google.com"},
		want:   "F870F1E9",
		exists: true,
	},
	{
		name:   "youtube",
		args:   args{key: "4AED1C05", value: "https://www.youtube.com"},
		want:   "4AED1C05",
		exists: true,
	},
	{
		name:   "youtube2",
		args:   args{key: "short4", value: "https://www.youtube.com"},
		want:   "4AED1C05",
		exists: false,
	},
	{
		name:   "empty key youtube",
		args:   args{key: "", value: "https://www.youtube.com"},
		want:   "4AED1C05",
		exists: false,
	},
	{
		name:   "empty value",
		args:   args{key: "empty0", value: ""},
		want:   "",
		exists: false,
	},
}

func TestSet(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := Set(tt.args.key, tt.args.value)
			assert.Equal(t, tt.want, got)
			fmt.Printf("Name = %v, Key = %v, Value = %v, Want = %v, Got = %v\n", tt.name, tt.args.key, tt.args.value, tt.want, got)
		})
	}
}

func TestGet(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, _ := Set(tt.args.key, tt.args.value)
			assert.Equal(t, tt.want, got)

			gotValue := Get(tt.want)
			assert.Equal(t, tt.args.value, gotValue)
		})
	}
}

func TestGetData(t *testing.T) {
	Clear()

	// Add test data
	urls := map[string]string{
		"google":  "https://www.google.com",
		"youtube": "https://www.youtube.com",
		"yandex":  "https://www.yandex.ru",
	}
	// Add URLs
	for shortID, url := range urls {
		Set(shortID, url)
	}

	// Get data
	data := GetData()

	// Check if all URLs are in the data
	for shortID, url := range urls {
		dataURL := data[shortID]
		assert.Equal(t, url, dataURL)
	}
}

func TestIsDeletedKey(t *testing.T) {
	type args struct {
		key        string
		keyToCheck string
		value      string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "empty key",
			args: args{key: "", keyToCheck: "", value: "@"},
			want: false,
		},
		{
			name: "deleted key",
			args: args{key: "-deletedKey", keyToCheck: "deletedKey", value: "@deletedValue"},
			want: true,
		},
		{
			name: "normal key",
			args: args{key: "key", keyToCheck: "key", value: "@value"},
			want: false,
		},
	}

	Clear()
	// Add test data
	for _, tt := range tests {
		Set(tt.args.key, tt.args.value)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDeletedKey(tt.args.keyToCheck); got != tt.want {
				t.Errorf("IsDeletedKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	type args struct {
		key   string
		value string
	}

	var tests = []struct {
		name   string
		args   args
		want   string
		exists bool
	}{
		{
			name:   "empty key & value",
			args:   args{key: "", value: "@"},
			want:   "",
			exists: true,
		},
		{
			name:   "google",
			args:   args{key: "F870F1E9", value: "@https://www.google.com"},
			want:   "F870F1E9",
			exists: true,
		},
		{
			name:   "youtube",
			args:   args{key: "4AED1C05", value: "@https://www.youtube.com"},
			want:   "4AED1C05",
			exists: true,
		},
		{
			name:   "youtube2",
			args:   args{key: "short4", value: "@https://www.youtube.com"},
			want:   "4AED1C05",
			exists: false,
		},
		{
			name:   "empty key youtube",
			args:   args{key: "", value: "@https://www.youtube.com"},
			want:   "4AED1C05",
			exists: false,
		},
		{
			name:   "empty value",
			args:   args{key: "empty0", value: "@"},
			want:   "",
			exists: false,
		},
	}

	Clear()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Set(tt.args.key, tt.args.value)
			err := Delete("", tt.args.key)
			if err != nil && tt.exists {
				t.Errorf("DeleteKey() error = %v, want nil", err)
			}
		})
	}
	PrintContent(10)
}

func TestDeleteKeys(t *testing.T) {
	type args struct {
		key   string
		value string
	}

	var tests = []struct {
		name   string
		args   args
		want   string
		exists bool
	}{
		{
			name:   "empty key & value",
			args:   args{key: "", value: "@"},
			want:   "",
			exists: true,
		},
		{
			name:   "google",
			args:   args{key: "F870F1E9", value: "@https://www.google.com"},
			want:   "F870F1E9",
			exists: true,
		},
		{
			name:   "youtube",
			args:   args{key: "4AED1C05", value: "@https://www.youtube.com"},
			want:   "4AED1C05",
			exists: true,
		},
		{
			name:   "youtube2",
			args:   args{key: "short4", value: "@https://www.youtube.com"},
			want:   "4AED1C05",
			exists: false,
		},
		{
			name:   "empty key youtube",
			args:   args{key: "", value: "@https://www.youtube.com"},
			want:   "4AED1C05",
			exists: false,
		},
		{
			name:   "empty value",
			args:   args{key: "empty0", value: "@"},
			want:   "",
			exists: false,
		},
	}

	Clear()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Set(tt.args.key, tt.args.value)
		})
	}

	err := DeleteKeys("", nil)
	if err != nil {
		t.Errorf("DeleteKeys() error = %v, want nil", err)
	}

	err = DeleteKeys("", []any{})
	if err != nil {
		t.Errorf("DeleteKeys() error = %v, want nil", err)
	}

	err = DeleteKeys("", []any{"F870F1E9", "4AED1C05"})
	if err != nil {
		t.Errorf("DeleteKeys() error = %v, want nil", err)
	}

	PrintContent(10)

}
