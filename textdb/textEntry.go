package textdb

// Metadata represents the user entered fields for a record.
// Keep these fields separate so that we can serialize
// them to XML files easily.
type Metadata struct {
	Slug      string `xml:"slug"`
	Title     string `xml:"title"`
	Summary   string `xml:"summary"`
	CreatedOn string `xml:"createdOn"`
	UpdatedOn string `xml:"updatedOn"`
	PostedOn  string `xml:"postedOn"`
}

// TextEntry represents a complete entry in the database,
// It includes the user provided Metadata plus the text
// content of the blog post.
type TextEntry struct {
	Metadata Metadata
	Content  string
	Id       string
}

func (entry *TextEntry) setSlug() {
	entry.Metadata.Slug = slug(entry.Metadata.Title)
}

func (entry *TextEntry) setUpdated() {
	entry.Metadata.UpdatedOn = now()
}

func (entry *TextEntry) MarkAsPosted() {
	entry.Metadata.PostedOn = now()
}

func (entry *TextEntry) MarkAsDraft() {
	entry.Metadata.PostedOn = ""
}

func (entry TextEntry) IsDraft() bool {
	return entry.Metadata.PostedOn == ""
}
