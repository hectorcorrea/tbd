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

// NewEntry creates a new record and initializes it.
// Uses today's date for the basis of the Id.
func (db *TextDb) NewEntry() (TextEntry, error) {
	id := db.getNextId()
	entry := NewTextEntry(id)
	entry.Title = "new " + id
	entry.Content = "(content goes here)"
	return db.saveEntry(entry)
}

// NewEntryFor creates a new record for a specific date and time.
// This is useful when importing existing data as it uses the given date for the basis
// of the Id and sets the CreatedOn appropriately.
// Date is expected to be in the form yyyy-mm-dd
// Time is expected to be in the form HH:mm:ss.xxx
func (db *TextDb) NewEntryFor(date string, time string) (TextEntry, error) {
	id := db.getNextIdFor(date)
	entry := NewTextEntry(id)
	entry.Title = "new " + id
	entry.Content = "(content goes here)"
	entry.CreatedOn = date + " " + time
	return db.saveEntry(entry)
}

// Saves an existing entry
func (db *TextDb) UpdateEntry(entry TextEntry) (TextEntry, error) {
	return db.saveEntry(entry)
}

func (db *TextDb) saveEntry(entry TextEntry) (TextEntry, error) {
	// Always set the slug before saving and make sure the Id
	// still is valid.
	err := validId(entry.Id)
	if err != nil {
		return entry, err
	}

	entry.Save()

	// Create the directory for it if it does not exist
	path := db.entryPath(entry)
	if !dirExist(path) {
		logInfo("Creating path", path)
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			logError("Error creating path", path, err)
			return entry, err
		}
	}
	// Save metadata + content
	err = saveMetadata(path, entry)
	if err == nil {
		err = saveContent(path, entry)
	}
	return entry, err
}

func (db *TextDb) All() []TextEntry {
	entries := []TextEntry{}
	err := filepath.Walk(db.RootDir, func(path string, info os.FileInfo, err error) error {
		if path == db.RootDir {
			return nil
		}
		if info.IsDir() {
			id := idFromPath(path)
			entry, err := db.readEntry(id)
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

func (db *TextDb) FindById(id string) (TextEntry, error) {
	err := validId(id)
	if err != nil {
		return TextEntry{}, err
	}
	return db.readEntry(id)
}

func (db *TextDb) FindBySlug(slug string) (TextEntry, bool) {
	for _, entry := range db.All() {
		if entry.Slug == slug {
			return entry, true
		}
	}
	return TextEntry{}, false
}

func (db *TextDb) readEntry(id string) (TextEntry, error) {
	path := filepath.Join(db.RootDir, id)
	if !dirExist(path) {
		logError("ReadEntry did not find path", path, nil)
		return TextEntry{}, errors.New("Path not found")
	}

	entry := readMetadata(filepath.Join(path, "metadata.xml"))
	entry.Id = idFromPath(path)
	entry.Content = readContent(filepath.Join(path, "content.md"))
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
