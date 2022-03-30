package textodb

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"os"
	"path/filepath"
)

func dirExist(name string) bool {
	file, err := os.Open(name)
	if os.IsNotExist(err) {
		return false
	}
	defer file.Close()
	return true
}

func readContent(filename string) string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		logError("Error reading content file", filename, err)
	}
	return string(content)
}

func readMetadata(filename string) TextoEntry {
	reader, err := os.Open(filename)
	if err != nil {
		logError("Error reading metadata file", filename, err)
	}
	defer reader.Close()

	// Read the bytes and unmarshall into our TextoEntry struct
	byteValue, _ := ioutil.ReadAll(reader)
	var entry TextoEntry
	xml.Unmarshal(byteValue, &entry)
	return entry
}

func saveContent(path string, entry TextoEntry) error {
	filename := filepath.Join(path, "content.md")
	return ioutil.WriteFile(filename, []byte(entry.Content), 0644)
}

func saveMetadata(path string, entry TextoEntry) error {
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
	filename := filepath.Join(path, "metadata.xml")
	return ioutil.WriteFile(filename, buffer.Bytes(), 0644)
}
