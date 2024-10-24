// Мигграции.
// Файлы миграций содержат SQL команды для создания объектов базы данных.
// Файлы миграций находятся в отдельной директории (./migrations), имеют расширение *.up.sql,
// сортируются по имени и выполняются в порядке возрастания.

package db

import (
	"errors"
	"io/fs"
	"path"

	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

// MigrateUp порождает объекты базы данных
// dirname - директория с файлами миграций
func MigrateUp(dirname string) error {
	files, err := os.ReadDir(dirname)
	if err != nil {
		log.Error().Err(err).Msg("MigrateUp error in directory")
		return err
	}
	return executeSQLFiles(files, dirname, ".up.sql")
}

// executeSQLFiles выполняем SQL команды в файлах с расширением filenameSuffix в директории dirname
func executeSQLFiles(files []fs.DirEntry, dirname, filenameSuffix string) error {

	for _, file := range files {
		fileName := file.Name()
		if strings.HasSuffix(fileName, filenameSuffix) {
			log.Info().Str("filename", fileName).Msg("Executing")
			if !IsConnected() {
				log.Error().Msg("executeSQLFiles(). No connection to DB")
				return errors.New("executeSQLFiles(). No connection to DB")
			}
			_, err := DB.Exec(ReadTextFile(path.Join(dirname, fileName)))
			if err != nil {
				log.Error().Err(err).Msg("Query Execution error")
				return err
			}
		}
	}

	return nil
}

// ReadTextFile  - читает текстовый файл и возвращает его содержимое в виде строки.
// fileName - имя файла
func ReadTextFile(fileName string) string {
	b, err := os.ReadFile(fileName)
	if err != nil {
		log.Error().Err(err).Msg("ReadTextFile error")
	}
	return string(b)
}
