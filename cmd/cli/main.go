package main

import (
	"archive/zip"
	"crypto/sha1"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/scirelli/rom-files-cleanup/internal/pkg/log"
	"github.com/scirelli/rom-files-cleanup/internal/pkg/model"
)

const defaultDatabaseFiles = "./assets/databases"
const defaultROMsFolder = "test/data/ROMs"

// list of known compressed file MIME types
var compressedMIMETypes = [...]string{
	"application/zip",
	"application/x-tar",
	"application/x-gzip",
	"application/x-7z-compressed",
	"application/x-bzip2",
}

// list of known file extensions that are compressed
var compressedFileExtensions = [...]string{
	".zip",
	".tar",
	".gz",
	".7z",
	".bz2",
}

var skipExt = []string{
	".gbp",
	".pal",
	".txt",
	".exe",
	".dll",
	".rom",
}

var logger = log.New("Main", log.DEFAULT_LOG_LEVEL)

/*
	 TODO:
		* Need to handle .zip files
		  - https://stackoverflow.com/questions/18194688/how-can-i-determine-if-a-file-is-a-zip-file
		  - https://pkg.go.dev/archive/zip
*/

type RomDict map[string]model.DatRom

func main() {
	var romDict RomDict
	var notFound, totalChecked, dupesCnt int
	dupes := make(map[string]struct{})
	//config := NewConfig()

	var romPath *string = flag.String("rom-path", defaultROMsFolder, fmt.Sprintf("path to the ROMs directory. Default '%s'", defaultROMsFolder))
	var databasePath *string = flag.String("database-path", defaultDatabaseFiles, fmt.Sprintf("path to the ROM databases directory. Default '%s'", defaultDatabaseFiles))
	flag.Parse()

	romDict = loadAllDatabases(*databasePath)
	logger.Infof("%d ROMs loaded from databases\n", len(romDict))
	if err := filepath.WalkDir(*romPath, func(path string, d fs.DirEntry, err error) error {
		var hash []map[string]string

		if d.IsDir() || excludeFile(path) {
			return nil
		}

		if isArchive(path) {
			if r, err := zip.OpenReader(path); err == nil {
				defer r.Close()
				for _, f := range r.File {
					if excludeFile(f.FileHeader.Name) {
						continue
					}
					//logger.Infof("Decompressed from zip file: %s %s\n", f.Name, path)
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
			if _, ok := dupes[h["h"]]; ok {
				dupesCnt++
			} else {
				dupes[h["h"]] = struct{}{}
				if _, ok := romDict[h["h"]]; !ok {
					logger.Warnf("Not found: %s\t%s\n", h["f"], h["h"])
					notFound++
				}
			}
		}

		return nil
	}); err != nil {
		logger.Error(err)
	}
	totalChecked = len(dupes)
	logger.Warnf("\nCould not find %d ROMs out of %d in any loaded database. \n%.2f%% not found\nDuplicates: %d", notFound, totalChecked, (float64(notFound)/float64(totalChecked))*100, dupesCnt)
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

func loadDATFile(romDict RomDict, dataFile string) RomDict {
	var datomaticData model.DatoMaticDataFile
	if fileData, err := os.ReadFile(dataFile); err == nil {
		if err := xml.Unmarshal(fileData, &datomaticData); err == nil {
			for _, game := range datomaticData.Games {
				for _, rom := range game.Rom {
					romDict[strings.ToLower(rom.Sha1)] = rom
				}
			}
		} else {
			logger.Warnf("%v", err)
		}
	} else {
		logger.Warnf("%v", err)
	}

	return romDict
}

func loadHSIFile(romDict RomDict, dataFile string) RomDict {
	return romDict
}

func loadAllDatabases(databaseDir string) RomDict {
	var romDict = make(RomDict)

	if err := filepath.WalkDir(databaseDir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		switch strings.ToLower(filepath.Ext(path)) {
		case ".dat":
			logger.Infof("Loading database '%s'", path)
			romDict = loadDATFile(romDict, path)
		case ".hsi":
			logger.Infof("Loading database '%s'", path)
			romDict = loadHSIFile(romDict, path)
		}
		return nil
	}); err != nil {
		logger.Error(err)
	}

	return romDict
}

func fileNameWithoutExtension(fileName string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}

func excludeFile(path string) bool {
	return slices.Contains(skipExt, strings.ToLower(filepath.Ext(path)))
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
	if strings.EqualFold(mime.TypeByExtension(filepath.Ext(fileName)), filetype) {
		switch filetype {
		case "application/x-gzip", "application/zip":
			return true
		}
	}
	return false
}
