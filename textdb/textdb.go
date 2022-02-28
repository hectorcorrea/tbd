package textdb

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Metadata struct {
	Title  string `xml:"title"`
	Author string `xml:"author"`
}

type TextEntry struct {
	Title   string
	Content string
}

type TextDb struct {
	RootDir string
}

func (db *TextDb) ListAll() []TextEntry {
	rootDir := "./data/"
	entries := []TextEntry{}
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			metadata := readMetadata(path + "/metadata.xml")
			entry := TextEntry{
				Title:   metadata.Title,
				Content: readContent(path + "/content.md"),
			}
			entries = append(entries, entry)
		}
		// if !info.IsDir() {
		// 	entry := TextEntry{
		// 		Title:   path,
		// 		Content: readOne(path),
		// 	}
		// 	entries = append(entries, entry)
		// }
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	return entries
}

func readMetadata(filename string) Metadata {
	xmlFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened users.xml")
	// defer the closing of our xmlFile so that we can parse it later on
	defer xmlFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(xmlFile)

	// we initialize our Users array
	var metadata Metadata
	xml.Unmarshal(byteValue, &metadata)
	return metadata
}

func readContent(filename string) string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Err")
	}
	return string(content)
}
