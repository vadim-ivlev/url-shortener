// Description: функции для работы с JWT-токенами.
package auth

import (
	"strings"
	"testing"
)

func TestBuildToken(t *testing.T) {
	type args struct {
		userID string
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "test1",
			args:    args{userID: "id1"},
			want:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jwtString, err := BuildToken(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.HasPrefix(jwtString, tt.want) {
				t.Errorf("BuildToken() \n JWT string \n %#v \n does not have the prefix \n %#v", jwtString, tt.want)

			}

		})
	}
}

func TestGetUserIDFromToken(t *testing.T) {
	type args struct {
		userID string
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "test1",
			args:    args{userID: "id1"},
			want:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jwtString, err := BuildToken(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			gotUserID, err := GetUserIDFromToken(jwtString)
			// assert.NoError(t, err, "GetUserID() error = %v", err)
			if err != nil {
				t.Errorf("GetUserID() error = %v", err)
			}
			if gotUserID != tt.args.userID {
				t.Errorf("GetUserID() got %v, want %v", gotUserID, tt.args.userID)
			}

		})
	}
}
