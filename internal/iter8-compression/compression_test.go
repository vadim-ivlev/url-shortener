package compression

import (
	"net/http"
	"reflect"
	"testing"
)

func Test_newCompressWriter(t *testing.T) {
	type args struct {
		w http.ResponseWriter
	}
	tests := []struct {
		name string
		args args
		want *compressWriter
	}{
		{
			name: "Test 1",
			args: args{
				w: nil,
			},
			want: &compressWriter{
				w:  nil,
				zw: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newCompressWriter(tt.args.w)
			if got.zw == nil {
				t.Errorf("newCompressWriter() zw = nil")
			}
		})
	}
}

func TestGzipMiddleware(t *testing.T) {
	type args struct {
		next http.Handler
	}
	tests := []struct {
		name string
		args args
		want http.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GzipMiddleware(tt.args.next); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GzipMiddleware() = %v, want %v", got, tt.want)
			}
		})
	}
}
