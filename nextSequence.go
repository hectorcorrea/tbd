package textodb

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

// Gets the ID for a new record created today.
// For now all IDs are generated as today's date + sequence.
// In the future we might break that down by year or year + month.
func (db *TextoDb) getNextId() string {
	today := today()
	sequence := db.getNextSequence(today)
	id := fmt.Sprintf("%s-%05d", today, sequence) // yyyy-mm-dd-00000
	return id
}

// Gets the ID for a new record created on a given date.
// Date is expected to be in the form yyyy-mm-dd
func (db *TextoDb) getNextIdFor(date string) string {
	sequence := db.getNextSequence(date)
	id := fmt.Sprintf("%s-%05d", date, sequence) // yyyy-mm-dd-00000
	return id
}

// Gets the next sequence number for a given date.
func (db *TextoDb) getNextSequence(date string) int {
	// Get all the directories for the given date...
	mask := filepath.Join(db.RootDir, date) + "-*"
	directories, err := filepath.Glob(mask)
	if err != nil {
		// This is bad, stop the presses
		panic(err)
	}

	// ...and find the max sequence number from them
	maxSequence := 0
	prefix := filepath.Join(db.RootDir, date) + "-"
	for _, directory := range directories {
		sequenceStr := strings.TrimPrefix(directory, prefix)
		sequence, err := strconv.Atoi(sequenceStr)
		if err != nil {
			// Unexpected but not fatal
			logError("Unexpected directory found", directory, err)
		} else if sequence > maxSequence {
			maxSequence = sequence
		}
	}

	// ...increase the sequence number by one
	return maxSequence + 1
}
