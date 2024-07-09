package storage

import (
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
		args: args{key: "short1", value: "https://www.google.com"},
		want: "short1",
	},
	{
		name: "youtube",
		args: args{key: "short2", value: "https://www.youtube.com"},
		want: "short2",
	},
	{
		name: "youtube2",
		args: args{key: "short4", value: "https://www.youtube.com"},
		want: "short2",
	},
	{
		name: "empty key youtube",
		args: args{key: "", value: "https://www.youtube.com"},
		want: "short2",
	},
	{
		name: "empty value",
		args: args{key: "short3", value: ""},
		want: "",
	},
}

func TestSet(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Set(tt.args.key, tt.args.value)
			assert.Equal(t, tt.want, got)
			// fmt.Printf("Name = %v, Key = %v, Value = %v, Want = %v, Got = %v\n", tt.name, tt.args.key, tt.args.value, tt.want, got)
		})
	}
}

func TestGet(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := Set(tt.args.key, tt.args.value)
			assert.Equal(t, tt.want, got)

			gotValue := Get(tt.want)
			assert.Equal(t, tt.args.value, gotValue)
			// fmt.Printf("Name = %v, PassedKey = %v, Wanted Value = %v, GotValue = %v\n", tt.name, tt.want, tt.args.value, gotValue)
		})
	}
}
