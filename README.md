# texto

`texto` is a Go module to store data in text format. I use it to store the data in my small, personal, and simplistic websites (e.g. like https://hectorcorrea.com/blog/).

Currently `texto` supports a very specific data structure for each record with a few predefined fields: `title`, `summary`, `content`, and `postedOn`. These fields are used to support a blog web site.

Each record is stored in two files: `metadata.xml` and `content.md`. All the fields (except `content`) are stored in the `metadata.xml` file. The `content` field is stored in the `content.md` file.

The `metadata.xml` file also stores a few pre-calculated fields for each record, including `slug`, `createdOn`, and `updatedOn`.

There is support to add custom single value string fields to each record.


## Source code

* Package `textdb` is the core functionality.

* Package`demoServer` is a tiny website showing how to use it from a Go client. File `server.go` shows the functionality to create, display, and update records from inside  a Go website.


## Examples
Below are a few code snippets showing the basic functionality create, fetch, and update records.

```
import "github.com/hectorcorrea/texto/textdb"

// Initialize the database 
// (creates data folder if it does not exist)
dataFolder := './data'
db = textdb.InitTextDb(dataFolder)


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
entry.Content = "blah blah blah"
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

## textdb
(more stuff goes here)

## demoServer
(more stuff goes here)

## Notes
Will this scale? Probably no.
