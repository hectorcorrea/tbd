package textdb

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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

func dirExist(name string) bool {
	file, err := os.Open(name)
	if os.IsNotExist(err) {
		return false
	}
	defer file.Close()
	return true
}

// Creates a new record for today and initalizes it
func (db *TextDb) CreateNewEntry() error {
	metadata := Metadata{Title: "new", Author: "", Slug: "new-entry"}
	content := "(to be defined)"
	path := db.getNextPath()
	entry := TextEntry{Metadata: metadata, Content: content, Path: path}
	return db.SaveEntry(entry)
}

// Gets the path for a new record created today.
// For now all records are at the db.RootDir + date + sequence.
// In the future we might break that down by year or year + month.
func (db *TextDb) getNextPath() string {
	today := time.Now().Format("2006-01-02")
	sequence := db.getNextSequence(today)
	basePath := filepath.Join(db.RootDir, today)
	path := fmt.Sprintf("%s-%05d", basePath, sequence)
	return path
}

func (db *TextDb) getNextSequence(date string) int {
	// Get all the directories for the given date...
	mask := filepath.Join(db.RootDir, date) + "-*"
	directories, err := filepath.Glob(mask)
	if err != nil {
		// This is bad, stop the presses
		panic(err)
	}

	// ...and find the max sequence number from them
	maxSequence := 0
	prefix := filepath.Join(db.RootDir, date) + "-"
	for _, directory := range directories {
		sequenceStr := strings.TrimPrefix(directory, prefix)
		sequence, err := strconv.Atoi(sequenceStr)
		if err != nil {
			// Unexpected but not fatal
			logError("Unexpected directory found", directory, err)
		} else if sequence > maxSequence {
			maxSequence = sequence
		}
	}

	// ...increase the sequence number by one
	return maxSequence + 1
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

func saveMetadata(entry TextEntry) error {
	// Convert our Metadata struct to an XML string...
	xmlDeclaration := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\r\n"
	buffer := bytes.NewBufferString(xmlDeclaration)
	encoder := xml.NewEncoder(buffer)
	encoder.Indent("  ", "    ")
	err := encoder.Encode(entry.Metadata)
	if err != nil {
		return err
	}
	// ... and save it.
	filename := filepath.Join(entry.Path, "metadata.xml")
	return ioutil.WriteFile(filename, buffer.Bytes(), 0644)
}

func saveContent(entry TextEntry) error {
	filename := filepath.Join(entry.Path, "content.md")
	return ioutil.WriteFile(filename, []byte(entry.Content), 0644)
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

func readMetadata(filename string) Metadata {
	reader, err := os.Open(filename)
	if err != nil {
		logError("Error reading metadata file", filename, err)
	}
	defer reader.Close()

	// Read the bytes and unmarshall into our metadata struct
	byteValue, _ := ioutil.ReadAll(reader)
	var metadata Metadata
	xml.Unmarshal(byteValue, &metadata)
	return metadata
}

func readContent(filename string) string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		logError("Error reading content file", filename, err)
	}
	return string(content)
}

func logInfo(message string, parameter string) {
	log.Printf("textdb: %s %s", message, parameter)
}

func logError(message string, parameter string, err error) {
	log.Printf("textdb: %s %s. ERROR: %s", message, parameter, err)
}
