# textodb

`textodb` is a Go module to store data in text format. I use it to store the data in my small, personal, and simplistic websites (e.g. like https://hectorcorrea.com/blog/).

Currently `textodb` supports a very specific data structure for each record with a few predefined fields: `title`, `summary`, `content`, and `postedOn`. These fields are used to support a blog web site.

Data for each record is stored in two files: `metadata.xml` and `content.md`. All the fields (except `content`) are stored in the `metadata.xml` file. The `content` field is stored in the `content.md` file.

The `metadata.xml` file also stores a few pre-calculated fields for each record, including `slug`, `createdOn`, and `updatedOn`.

There is support to add custom single-value string fields to each record.


## Source code

* `textoDb.go` The main code for the module, functions to create, fetch, and update records is here.
* `textoEntry.go` Functions specific for individual entries. The code to perform most of the IO is here.
* `nextSequence.go` Calculates new sequence numbers for IDs.
* `util.go` Miscellanous utilities.

* Folder `demoServer/` is a tiny website showing how to use it from a Go client. File `server.go` shows the functionality to create, display, and update records from inside  a Go website.


## Examples
Below are a few code snippets showing the basic functionality create, fetch, and update records.

```
import "github.com/hectorcorrea/textodb"

// Initialize the database
// (creates data folder if it does not exist)
dataFolder := './data'
db = textodb.InitTextoDb(dataFolder)


// Create a new record
// (Id will be today's date plus a sequence number)
entry, err := db.NewEntry()
if err != nil {
    panic(err)
}
log.Printf("Created %s %s %s\r\n", entry.Id, entry.Slug, entry.CreatedOn)


// Find a record by Id
entry, err = db.FindById(entry.Id)


// Update it
entry.Title = "the updated title"
entry.SetContent("blah blah blah")
entry, err = db.UpdateEntry(entry)
if err != nil {
    panic(err)
}
log.Printf("Updated %s %s %s\r\n", entry.Id, entry.Slug, entry.UpdatedOn)


// Update a custom field (author)
entry.SetField("author", "hector")
entry, err = db.UpdateEntry(entry)
entry, err = db.FindById(entry.Id)
log.Printf("Id=%s, Title=%s, Author=%s\r\n", entry.Id, entry.Title, entry.GetField("author))
```

## Notes
Will this scale? Probably no.

Shouldn't I be using a well-known database instead? Probably yes.
