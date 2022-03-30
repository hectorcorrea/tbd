package textodb

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type TextoDb struct {
	RootDir string
}

func InitTextDb(rootDir string) TextoDb {
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
	return TextoDb{RootDir: rootDir}
}

// NewEntry creates a new record and initializes it.
// Uses today's date for the basis of the Id.
func (db *TextoDb) NewEntry() (TextoEntry, error) {
	id := db.getNextId()
	entry := NewTextoEntry(id)
	entry.Title = "new " + id
	return db.saveEntry(entry, true)
}

// NewEntryFor creates a new record for a specific date and time.
// This is useful when importing existing data as it uses the given date for the basis
// of the Id and sets the CreatedOn appropriately.
// Date is expected to be in the form yyyy-mm-dd
// Time is expected to be in the form HH:mm:ss.xxx
func (db *TextoDb) NewEntryFor(date string, time string) (TextoEntry, error) {
	id := db.getNextIdFor(date)
	entry := NewTextoEntry(id)
	entry.Title = "new " + id
	entry.CreatedOn = date + " " + time
	return db.saveEntry(entry, true)
}

// Saves an existing entry
func (db *TextoDb) UpdateEntry(entry TextoEntry) (TextoEntry, error) {
	return db.saveEntry(entry, true)
}

// Saves an existing entry but honors the createdOn and updatedOn
// values already on the entry rather than re-calculating them.
func (db *TextoDb) UpdateEntryHonorDates(entry TextoEntry) (TextoEntry, error) {
	return db.saveEntry(entry, false)
}

func (db *TextoDb) saveEntry(entry TextoEntry, calculateDates bool) (TextoEntry, error) {
	err := validId(entry.Id)
	if err != nil {
		return entry, err
	}

	entry.setCalculatedValues(calculateDates)

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

func (db *TextoDb) All() []TextoEntry {
	entries := []TextoEntry{}
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

// Finds an entry by Id
func (db *TextoDb) FindById(id string) (TextoEntry, error) {
	err := validId(id)
	if err != nil {
		return TextoEntry{}, err
	}
	return db.readEntry(id)
}

// Finds an entry by Slug
func (db *TextoDb) FindBySlug(slug string) (TextoEntry, bool) {
	for _, entry := range db.All() {
		if entry.Slug == slug {
			return entry, true
		}
	}
	return TextoEntry{}, false
}

// Finds an entry by a user defined field/value
func (db *TextoDb) FindBy(field string, value string) (TextoEntry, bool) {
	for _, entry := range db.All() {
		if entry.GetField(field) == value {
			return entry, true
		}
	}
	return TextoEntry{}, false
}

func (db *TextoDb) readEntry(id string) (TextoEntry, error) {
	path := filepath.Join(db.RootDir, id)
	if !dirExist(path) {
		logError("ReadEntry did not find path", path, nil)
		return TextoEntry{}, errors.New("Path not found")
	}

	entry := readMetadata(filepath.Join(path, "metadata.xml"))
	entry.Id = idFromPath(path)
	entry.Content = readContent(filepath.Join(path, "content.md"))
	return entry, nil
}

// Returns the full path to an entry
func (db *TextoDb) entryPath(entry TextoEntry) string {
	return filepath.Join(db.RootDir, entry.Id)
}

// Returns the Id from a path (i.e. the last segment of the path)
func idFromPath(path string) string {
	pathTokens := strings.Split(path, string(os.PathSeparator))
	return pathTokens[len(pathTokens)-1]
}
