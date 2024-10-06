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
	name string
	args args
	want string
}{
	{
		name: "empty key & value",
		args: args{key: "", value: ""},
		want: "",
	},
	{
		name: "google",
		args: args{key: "F870F1E9", value: "https://www.google.com"},
		want: "F870F1E9",
	},
	{
		name: "youtube",
		args: args{key: "4AED1C05", value: "https://www.youtube.com"},
		want: "4AED1C05",
	},
	{
		name: "youtube2",
		args: args{key: "short4", value: "https://www.youtube.com"},
		want: "4AED1C05",
	},
	{
		name: "empty key youtube",
		args: args{key: "", value: "https://www.youtube.com"},
		want: "4AED1C05",
	},
	{
		name: "empty value",
		args: args{key: "empty0", value: ""},
		want: "",
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
	// Add test data
	urls := AddTestData()
	// Get data
	data := GetData()
	// Check if all URLs are in the data
	for shortID, url := range urls {
		assert.Equal(t, url, data[shortID])
	}
}

func AddTestData() map[string]string {
	// URLs to add
	urls := map[string]string{
		"google":  "https://www.google.com",
		"youtube": "https://www.youtube.com",
		"yandex":  "https://www.yandex.ru",
	}
	// Add URLs
	for shortID, url := range urls {
		Set(shortID, url)
	}
	return urls
}
