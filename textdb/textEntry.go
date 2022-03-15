package textdb

// Metadata represents the user entered fields for a record.
// Keep these fields separate so that we can serialize
// them to XML files easily.
type Metadata struct {
	Slug      string `xml:"slug"`
	Title     string `xml:"title"`
	Author    string `xml:"author"`
	CreatedOn string `xml:"createdOn"`
	UpdatedOn string `xml:"updatedOn"`
}

// TextEntry represents a complete entry in the database,
// It includes the user provided Metadata plus the text
// content of the blog post.
type TextEntry struct {
	Metadata Metadata
	Content  string
	Id       string
}
