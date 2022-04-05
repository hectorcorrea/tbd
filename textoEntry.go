package textodb

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

// TextoEntry represents the information from database entry.
// It includes the information stored in the metadata file
// plus the informtion stored in the content file.
type TextoEntry struct {
	Id        string  `xml:"-"` // not stored
	Title     string  `xml:"title"`
	Slug      string  `xml:"slug"`
	Summary   string  `xml:"summary"`
	content   string  `xml:"-"` // stored in separate doc
	CreatedOn string  `xml:"createdOn"`
	UpdatedOn string  `xml:"updatedOn"`
	PostedOn  string  `xml:"postedOn"`
	Fields    []Field `xml:"fields"` // client managed fields
	db        *TextoDb
}

// Field represents a custom field in a a TextoEntry.
type Field struct {
	Name  string `xml:"name"`
	Value string `xml:"value"`
}

// newTextoEntry initializes an entry object.
func newTextoEntry(db *TextoDb, id string) TextoEntry {
	return TextoEntry{db: db, Id: id}
}

// loadTextoEntry loads a TextoEntry from disk
// Notice that we need a db object here to access the path
// since there is no entry yet.
func loadTextoEntry(db *TextoDb, id string) (TextoEntry, error) {
	// Make sure the Id is valid
	err := validId(id)
	if err != nil {
		logError("loadTextoEntry - invalid Id received", id, nil)
		return TextoEntry{}, errors.New("Invalid Id received")
	}

	// Make sure the path exists
	path := filepath.Join(db.RootDir, id)
	if !dirExist(path) {
		logError("loadTextoEntry - path not found", path, nil)
		return TextoEntry{}, errors.New("Path not found")
	}

	// Open the metadata file...
	filename := filepath.Join(path, "metadata.xml")
	reader, err := os.Open(filename)
	if err != nil {
		logError("loadTextoEntry - error reading metadata file", filename, err)
		return TextoEntry{}, err
	}
	defer reader.Close()

	// ...and read it into a TextoEntry struct
	byteValue, err := ioutil.ReadAll(reader)
	var entry TextoEntry
	xml.Unmarshal(byteValue, &entry)
	entry.Id = id
	entry.db = db
	return entry, err
}

func (entry *TextoEntry) Path() string {
	return filepath.Join(entry.db.RootDir, entry.Id)
}

func (entry *TextoEntry) metadataFile() string {
	return filepath.Join(entry.Path(), "metadata.xml")
}

func (entry *TextoEntry) contentFile() string {
	return filepath.Join(entry.Path(), "content.md")
}

func (entry TextoEntry) Content() string {
	// TODO: should we cache this value so that multiple calls to Content()
	// don't re-read the file on disk?
	content, err := ioutil.ReadFile(entry.contentFile())
	if err != nil {
		logError("Error reading content file", entry.contentFile(), err)
	}
	return string(content)
}

func (entry *TextoEntry) SetContent(content string) {
	entry.content = content
}

func (entry *TextoEntry) Save(setDates bool) error {
	// Make sure the Id is valid
	err := validId(entry.Id)
	if err != nil {
		return err
	}

	// Make sure entry is linked to a database.
	if entry.db == nil {
		err := errors.New("No db set on entry")
		logError("Cannot save entry without a db value", entry.Id, err)
		return err
	}

	entry.setCalculatedValues(setDates)

	// Create the directory for it if does not exist
	path := entry.Path()
	if !dirExist(path) {
		logInfo("Creating path", path)
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			logError("Error creating path", path, err)
			return err
		}
	}

	// Save metadata + content
	err = entry.saveMetadata()
	if err != nil {
		return err
	}
	return entry.saveContent()
}

func (entry *TextoEntry) saveMetadata() error {
	// Convert our TextoEntry struct to an XML string...
	xmlDeclaration := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\r\n"
	buffer := bytes.NewBufferString(xmlDeclaration)
	encoder := xml.NewEncoder(buffer)
	encoder.Indent("  ", "    ")

	err := encoder.Encode(entry)
	if err != nil {
		return err
	}
	// ... and save it.
	filename := entry.metadataFile()
	return ioutil.WriteFile(filename, buffer.Bytes(), 0644)
}

func (entry *TextoEntry) saveContent() error {
	filename := entry.contentFile()
	return ioutil.WriteFile(filename, []byte(entry.content), 0644)
}

// Sets the value for a custom field in an entry.
func (entry *TextoEntry) SetField(name string, value string) {
	for i, field := range entry.Fields {
		if field.Name == name {
			// replace the existing values
			entry.Fields[i].Value = value
			return
		}
	}
	// create the field
	field := Field{Name: name, Value: value}
	entry.Fields = append(entry.Fields, field)
}

// Gets the value for a custom field in an entry.
func (entry *TextoEntry) GetField(name string) string {
	for _, field := range entry.Fields {
		if field.Name == name {
			return field.Value
		}
	}
	return ""
}

func (entry *TextoEntry) setCalculatedValues(setDates bool) {
	entry.Slug = slug(entry.Title)
	if setDates {
		if entry.CreatedOn == "" {
			entry.CreatedOn = now()
		} else {
			entry.UpdatedOn = now()
		}
	}
}

// MarkAsPosted sets the posted on value of the entry.
func (entry *TextoEntry) MarkAsPosted() {
	entry.PostedOn = now()
}

// MarkAsDraft clears the posted on value of the entry.
func (entry *TextoEntry) MarkAsDraft() {
	entry.PostedOn = ""
}

// IsDraft returns true if the posted on value is empty.
func (entry TextoEntry) IsDraft() bool {
	return entry.PostedOn == ""
}
