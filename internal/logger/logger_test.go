package logger

import (
	"net/http"
	"reflect"
	"testing"
)

func TestRequestLogger(t *testing.T) {
	type args struct {
		h http.Handler
	}
	tests := []struct {
		name string
		args args
		want http.Handler
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RequestLogger(tt.args.h); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RequestLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}
