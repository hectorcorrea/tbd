package textdb

type Metadata struct {
	Slug   string `xml:"slug"`
	Title  string `xml:"title"`
	Author string `xml:"author"`
}

type TextEntry struct {
	Metadata Metadata
	Content  string
	Id       string
}
