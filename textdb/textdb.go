package textdb

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Metadata struct {
	Slug   string `xml:"slug"`
	Title  string `xml:"title"`
	Author string `xml:"author"`
}

type TextEntry struct {
	Metadata Metadata
	Content  string
	Path     string
}

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
func (db *TextDb) CreateNewEntry() error {
	metadata := Metadata{Title: "new", Author: "", Slug: "new-entry"}
	content := "(to be defined)"
	path := db.getNextPath()
	entry := TextEntry{Metadata: metadata, Content: content, Path: path}
	return db.SaveEntry(entry)
}

func (db *TextDb) SaveEntry(entry TextEntry) error {
	if !dirExist(entry.Path) {
		logInfo("Creating path", entry.Path)
		if err := os.MkdirAll(entry.Path, os.ModePerm); err != nil {
			logError("Error creating path", entry.Path, err)
			return err
		}
	}
	err := saveMetadata(entry)
	if err == nil {
		err = saveContent(entry)
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
			metadata := readMetadata(filepath.Join(path, "metadata.xml"))
			entry := TextEntry{
				Path:     path,
				Metadata: metadata,
				Content:  readContent(filepath.Join(path, "content.md")),
			}
			entries = append(entries, entry)
		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	return entries
}

func logInfo(message string, parameter string) {
	log.Printf("textdb: %s %s", message, parameter)
}

func logError(message string, parameter string, err error) {
	log.Printf("textdb: %s %s. ERROR: %s", message, parameter, err)
}
