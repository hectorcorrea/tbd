package textdb

// TextEntry represents the information from database entry.
// It includes the information stored in the metadata file
// plus the informtion stored in the content file.
type TextEntry struct {
	Id        string  `xml:"-"` // not stored
	Title     string  `xml:"title"`
	Slug      string  `xml:"slug"`
	Summary   string  `xml:"summary"`
	Content   string  `xml:"-"` // stored in separate doc
	CreatedOn string  `xml:"createdOn"`
	UpdatedOn string  `xml:"updatedOn"`
	PostedOn  string  `xml:"postedOn"`
	Fields    []Field `xml:"fields"` // client managed fields
}

type Field struct {
	Name  string `xml:"name"`
	Value string `xml:"value"`
}

func NewTextEntry(id string) TextEntry {
	return TextEntry{Id: id}
}

func (doc *TextEntry) SetField(name string, value string) {
	for _, field := range doc.Fields {
		if field.Name == name {
			// replace the existing values
			field.Value = value
			return
		}
	}
	// create the field
	field := Field{Name: name, Value: value}
	doc.Fields = append(doc.Fields, field)
}

func (doc *TextEntry) GetField(name string) string {
	for _, field := range doc.Fields {
		if field.Name == name {
			return field.Value
		}
	}
	return ""
}

func (doc *TextEntry) setCalculatedValues() {
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
