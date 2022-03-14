package textdb

type Metadata struct {
	Slug      string `xml:"slug"`
	Title     string `xml:"title"`
	Author    string `xml:"author"`
	CreatedOn string `xml:"createdOn"`
	UpdatedOn string `xml:"updatedOn"`
}

type TextEntry struct {
	Metadata Metadata
	Content  string
	Id       string
}
