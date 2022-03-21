package textdb

// // Metadata represents the user entered fields for a record.
// // Keep these fields separate so that we can serialize
// // them to XML files easily.
// type Metadata struct {
// 	Slug      string `xml:"slug"`      // calc (always)
// 	CreatedOn string `xml:"createdOn"` // calc (except on import)
// 	UpdatedOn string `xml:"updatedOn"` // calc (except on import)
// 	Title     string `xml:"title"`     // user
// 	Summary   string `xml:"summary"`   // user
// 	PostedOn  string `xml:"postedOn"`  // user (via methods)
// }

// // TextEntry represents a complete entry in the database,
// // It includes the user provided Metadata plus the text
// // content of the blog post.
// type TextEntry struct {
// 	Metadata Metadata
// 	Content  string
// 	Id       string
// }

// ==============================

type TextEntryDisk struct {
	Title     string `xml:"title"`
	Summary   string `xml:"summary"`
	Slug      string `xml:"slug"`
	CreatedOn string `xml:"createdOn"`
	UpdatedOn string `xml:"updatedOn"`
	PostedOn  string `xml:"postedOn"`
}

type TextEntry struct {
	Id        string
	Title     string
	Slug      string
	Summary   string
	Content   string
	CreatedOn string
	UpdatedOn string
	PostedOn  string
}

func NewTextEntry(id string) TextEntry {
	return TextEntry{Id: id}
}

func (doc *TextEntry) Save() {
	doc.Slug = slug(doc.Title)
	if doc.CreatedOn == "" {
		doc.CreatedOn = now()
	} else {
		doc.UpdatedOn = now()
	}
}

func (doc *TextEntry) MarkAsPosted() {
	doc.PostedOn = now()
}

func (doc *TextEntry) MarkAsDraft() {
	doc.PostedOn = ""
}

func (doc TextEntry) IsDraft() bool {
	return doc.PostedOn == ""
}

func (doc TextEntry) ToTextEntryDisk() TextEntryDisk {
	ted := TextEntryDisk{
		Title:     doc.Title,
		Summary:   doc.Summary,
		Slug:      doc.Slug,
		CreatedOn: doc.CreatedOn,
		UpdatedOn: doc.UpdatedOn,
		PostedOn:  doc.PostedOn,
	}
	return ted
}

func NewTextEntryFromDisk(id string, ted TextEntryDisk) TextEntry {
	entry := TextEntry{
		Id:        id,
		Title:     ted.Title,
		Summary:   ted.Summary,
		Slug:      ted.Slug,
		CreatedOn: ted.CreatedOn,
		UpdatedOn: ted.UpdatedOn,
		PostedOn:  ted.PostedOn,
		Content:   "set outside here",
	}
	return entry
}
