// Package textodb implements functionality to read and store data into
// a very simple database stored in disk as text files.
package textodb

import (
	"log"
	"os"
	"path/filepath"
)

// TextoDb is the main object to access the database functionality.
type TextoDb struct {
	RootDir string // Directory where data will be stored
}

// InitTextoDb initializes a new TextoDb object.
func InitTextoDb(rootDir string) TextoDb {
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
// Uses today's date for the basis of the Id, the Id will be
// in the form yyyy-mm-dd-00000, e.g. "2022-03-25-00001"
func (db *TextoDb) NewEntry() (TextoEntry, error) {
	id := db.getNextId()
	entry := newTextoEntry(db, id)
	entry.Title = "new " + id
	entry.contentChanged = true
	return db.saveEntry(entry, true)
}

// NewEntryFor creates a new record for a specific date and time.
// This is useful when importing existing data as it uses the given date for the basis
// of the Id and sets the CreatedOn appropriately.
//
// Date is expected to be in the form yyyy-mm-dd
// Time is expected to be in the form HH:mm:ss.xxx
//
// As with NewEntry() the Id will be in the form yyyy-mm-dd-00000
// but using the date provided.
func (db *TextoDb) NewEntryFor(date string, time string) (TextoEntry, error) {
	id := db.getNextIdFor(date)
	entry := newTextoEntry(db, id)
	entry.Title = "new " + id
	entry.CreatedOn = date + " " + time
	entry.contentChanged = true
	return db.saveEntry(entry, true)
}

// Saves an existing entry
func (db *TextoDb) UpdateEntry(entry TextoEntry) (TextoEntry, error) {
	return db.saveEntry(entry, true)
}

// Saves an existing entry but honors the createdOn and updatedOn
// values already on the entry rather than re-calculating them.
// This is useful when importing existing data.
func (db *TextoDb) UpdateEntryHonorDates(entry TextoEntry) (TextoEntry, error) {
	return db.saveEntry(entry, false)
}

// Returns all the entries in the database.
func (db *TextoDb) All() []TextoEntry {
	entries := []TextoEntry{}
	err := filepath.Walk(db.RootDir, func(path string, info os.FileInfo, err error) error {
		if path == db.RootDir {
			return nil
		}
		if info.IsDir() {
			id := idFromPath(path)
			entry, err := loadTextoEntry(db, id)
			if err == nil {
				entries = append(entries, entry)
			}
		}
		return nil
	})

	if err != nil {
		logError("Error walking file system", "", err)
	}

	return entries
}

// Finds an entry by Id
func (db *TextoDb) FindById(id string) (TextoEntry, error) {
	return loadTextoEntry(db, id)
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

func (db *TextoDb) saveEntry(entry TextoEntry, setDates bool) (TextoEntry, error) {
	// Make sure the entry is linked to this database.
	entry.db = db

	err := entry.save(setDates)
	return entry, err
}
