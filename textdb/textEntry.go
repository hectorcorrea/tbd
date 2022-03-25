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

func (entry *TextEntry) SetField(name string, value string) {
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

func (entry *TextEntry) GetField(name string) string {
	for _, field := range entry.Fields {
		if field.Name == name {
			return field.Value
		}
	}
	return ""
}

func (entry *TextEntry) setCalculatedValues(calculateDates bool) {
	entry.Slug = slug(entry.Title)
	if calculateDates {
		if entry.CreatedOn == "" {
			entry.CreatedOn = now()
		} else {
			entry.UpdatedOn = now()
		}
	}
}

func (entry *TextEntry) MarkAsPosted() {
	entry.PostedOn = now()
}

func (entry *TextEntry) MarkAsDraft() {
	entry.PostedOn = ""
}

func (entry TextEntry) IsDraft() bool {
	return entry.PostedOn == ""
}
