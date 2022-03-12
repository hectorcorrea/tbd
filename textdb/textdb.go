package textdb

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type TextDb struct {
	RootDir string
}

func InitTextDb(rootDir string) TextDb {
	rootDir, err := filepath.Abs(rootDir)
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.Open(rootDir)

	if os.IsNotExist(err) {
		logInfo("Creating data folder", rootDir)
		err = os.Mkdir(rootDir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	} else if err != nil {
		log.Fatal(err)
	} else {
		logInfo("Using data folder", rootDir)
	}

	defer file.Close()
	return TextDb{RootDir: rootDir}
}

// Creates a new record for today and initalizes it
func (db *TextDb) NewEntry() (TextEntry, error) {
	metadata := Metadata{Title: "new", Author: "", Slug: "new-entry"}
	content := "(to be defined)"
	id := db.getNextId()
	entry := TextEntry{Metadata: metadata, Content: content, Id: id}
	return entry, db.SaveEntry(entry)
}

func (db *TextDb) SaveEntry(entry TextEntry) error {
	path := db.entryPath(entry)
	if !dirExist(path) {
		logInfo("Creating path", path)
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			logError("Error creating path", path, err)
			return err
		}
	}
	err := saveMetadata(path, entry)
	if err == nil {
		err = saveContent(path, entry)
	}
	return err
}

func (db *TextDb) ListAll() []TextEntry {
	entries := []TextEntry{}
	err := filepath.Walk(db.RootDir, func(path string, info os.FileInfo, err error) error {
		if path == db.RootDir {
			return nil
		}
		if info.IsDir() {
			id := idFromPath(path)
			entry, err := db.ReadEntry(id)
			if err == nil {
				entries = append(entries, entry)
			}
		}
		return nil
	})

	if err != nil {
		logError("ListAll error walking file system", "", err)
	}

	return entries
}

func (db *TextDb) ReadEntry(id string) (TextEntry, error) {
	// TODO: validate the ID cannot walk paths
	path := filepath.Join(db.RootDir, id)
	if !dirExist(path) {
		logError("ReadEntry did not find path", path, nil)
		return TextEntry{}, errors.New("Path not found")
	}

	entry := TextEntry{
		Id:       idFromPath(path),
		Metadata: readMetadata(filepath.Join(path, "metadata.xml")),
		Content:  readContent(filepath.Join(path, "content.md")),
	}
	return entry, nil
}

// Returns the full path to an entry
func (db *TextDb) entryPath(entry TextEntry) string {
	return filepath.Join(db.RootDir, entry.Id)
}

// Returns the Id from a path (i.e. the last segment of the path)
func idFromPath(path string) string {
	pathTokens := strings.Split(path, string(os.PathSeparator))
	return pathTokens[len(pathTokens)-1]
}

func logInfo(message string, parameter string) {
	log.Printf("textdb: %s %s", message, parameter)
}

func logError(message string, parameter string, err error) {
	log.Printf("textdb: %s %s. ERROR: %s", message, parameter, err)
}
