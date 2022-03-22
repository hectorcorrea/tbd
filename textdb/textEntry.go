package textdb

// id is not public because we don't want the clients to
// mess with it.
//
// content is not public because we want to serialize it
// to a different document.
type TextEntry struct {
	id        string
	Title     string `xml:"title"`
	Slug      string `xml:"slug"`
	Summary   string `xml:"summary"`
	content   string
	CreatedOn string `xml:"createdOn"`
	UpdatedOn string `xml:"updatedOn"`
	PostedOn  string `xml:"postedOn"`
}

func NewTextEntry(id string) TextEntry {
	return TextEntry{id: id}
}

func (doc TextEntry) Id() string {
	return doc.id
}

func (doc *TextEntry) SetId(id string) {
	doc.id = id
}
func (doc TextEntry) Content() string {
	return doc.content
}

func (doc *TextEntry) SetContent(content string) {
	doc.content = content
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
