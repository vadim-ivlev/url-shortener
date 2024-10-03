// Мигграции.
// Файлы миграций содержат SQL команды для создания объектов базы данных.
// Файлы миграций находятся в отдельной директории (./migrations), имеют расширение *.up.sql,
// сортируются по имени и выполняются в порядке возрастания.

package db

import (
	"testing"
)

func TestMigrateUp(t *testing.T) {
	// skipCI skips tests in CI environment
	skipCI(t)

	Connect(1)

	type args struct {
		dirname string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "TestMigrateUp",
			args:    args{dirname: "./migrations"},
			wantErr: false,
		},
		{
			name:    "TestMigrateUp nonexistent directory",
			args:    args{dirname: "./nonexistent_directory"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MigrateUp(tt.args.dirname); (err != nil) != tt.wantErr {
				t.Errorf("MigrateUp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
