package main

import (
	"archive/zip"
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/scirelli/rom-files-cleanup/internal/pkg/log"
	"github.com/scirelli/rom-files-cleanup/internal/pkg/model"
)

const defaultDatabaseFiles = "./assets/databases"
const defaultROMsFolder = "test/data/ROMs"

var logger = log.New("Main", log.DEFAULT_LOG_LEVEL)

/*
	 TODO:
		* Need to handle .zip files
		  - https://stackoverflow.com/questions/18194688/how-can-i-determine-if-a-file-is-a-zip-file
		  - https://pkg.go.dev/archive/zip
*/
func main() {
	var romDict map[string]model.Game
	var notFound, totalChecked int
	//config := NewConfig()

	romDict = loadDatabases(defaultDatabaseFiles)
	fmt.Printf("%d ROMs loaded from databases\n", len(romDict))
	if err := filepath.WalkDir(defaultROMsFolder, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		// Open a zip archive for reading.
		r, err := zip.OpenReader(path)
		if err != nil {
			logger.Error(err)
			return nil
		}
		defer r.Close()

		//fmt.Printf("%s | %s\n", path, romDict[hashFile(path)].Name)
		if romDict[hashFile(path)].Name == "" {
			notFound++
		}
		totalChecked++
		return nil
	}); err != nil {
		logger.Error(err)
	}
	fmt.Printf("Could not find %d ROMs out of %d in any loaded database. %.2f%% not found\n", notFound, totalChecked, (float64(notFound)/float64(totalChecked))*100)
}

func hashFile(file string) string {
	f, err := os.Open(file)
	if err != nil {
		logger.Errorf("%s %s", file, err)
		return ""
	}
	defer f.Close()

	return hashBytes(f)
}

func hashBytes(data io.Reader) string {
	h := sha1.New()
	if _, err := io.Copy(h, data); err != nil {
		logger.Errorf("%s", err)
		return ""
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}

func loadDatabases(dataFiles string) map[string]model.Game {
	var romDict = make(map[string]model.Game)
	databaseFiles, err := os.ReadDir(dataFiles)
	if err != nil {
		logger.Error(err)
		return romDict
	}

	for _, file := range databaseFiles {
		var datomaticData model.DatoMaticDataFile
		if file.IsDir() {
			continue
		}
		fileData, err := os.ReadFile(path.Join(defaultDatabaseFiles, file.Name()))
		if err != nil {
			logger.Infof("%v", err)
			continue
		}
		if err := xml.Unmarshal(fileData, &datomaticData); err != nil {
			logger.Infof("%v", err)
			continue
		}
		for _, game := range datomaticData.Games {
			romDict[strings.ToLower(game.Rom.Sha1)] = game
		}
	}
	return romDict
}

func fileNameWithoutExtension(fileName string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}
