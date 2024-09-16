package main

import (
	"archive/zip"
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
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

	romDict = loadAllDatabases(defaultDatabaseFiles)
	logger.Infof("%d ROMs loaded from databases\n", len(romDict))
	if err := filepath.WalkDir(defaultROMsFolder, func(path string, d fs.DirEntry, err error) error {
		var hash []map[string]string

		if d.IsDir() {
			return nil
		}

		if isArchive(path) {
			if r, err := zip.OpenReader(path); err == nil {
				defer r.Close()
				for _, f := range r.File {
					logger.Infof("Decompressed from zip file: %s\n", f.Name)
					if r, err := f.Open(); err == nil {
						if h := hashBytes(r); h != "" {
							hash = append(hash, map[string]string{"h": h, "f": f.Name})
						}
						r.Close()
					}
				}
			}
		} else {
			if h := hashFile(path); h != "" {
				hash = append(hash, map[string]string{"h": h, "f": path})
			}
		}

		for _, h := range hash {
			if romDict[h["h"]].Name == "" {
				logger.Warnf("Not found: %s\n", h["f"])
				notFound++
			} else {
				//logger.Info(romDict[h["h"]].Name)
			}
		}

		totalChecked++

		return nil
	}); err != nil {
		logger.Error(err)
	}
	logger.Warnf("Could not find %d ROMs out of %d in any loaded database. %.2f%% not found\n", notFound, totalChecked, (float64(notFound)/float64(totalChecked))*100)
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

func loadDatabase(romDict map[string]model.Game, dataFile string) map[string]model.Game {
	var datomaticData model.DatoMaticDataFile
	if fileData, err := os.ReadFile(dataFile); err == nil {
		if err := xml.Unmarshal(fileData, &datomaticData); err == nil {
			for _, game := range datomaticData.Games {
				romDict[strings.ToLower(game.Rom.Sha1)] = game
			}
		} else {
			logger.Warnf("%v", err)
		}
	} else {
		logger.Warnf("%v", err)
	}

	return romDict
}

func loadAllDatabases(databaseDir string) map[string]model.Game {
	var romDict = make(map[string]model.Game)

	if err := filepath.WalkDir(defaultDatabaseFiles, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		logger.Infof("Loading database '%s'", path)
		romDict = loadDatabase(romDict, path)
		return nil
	}); err != nil {
		logger.Error(err)
	}

	return romDict
}

func fileNameWithoutExtension(fileName string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}

func isArchive(fileName string) bool {
	file, err := os.Open(fileName)
	if err != nil {
		logger.Warn(err)
		return false
	}
	defer file.Close()

	buff := make([]byte, 512) // why 512 bytes ? see http://golang.org/pkg/net/http/#DetectContentType
	_, err = file.Read(buff)
	if err != nil {
		logger.Warn(err)
		return false
	}

	filetype := http.DetectContentType(buff)

	switch filetype {
	case "application/x-gzip", "application/zip":
		return true
	}
	return false
}
