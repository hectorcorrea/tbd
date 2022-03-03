package textdb

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Gets the path for a new record created today.
// For now all records are at the db.RootDir + date + sequence.
// In the future we might break that down by year or year + month.
func (db *TextDb) getNextPath() string {
	today := time.Now().Format("2006-01-02")
	sequence := db.getNextSequence(today)
	basePath := filepath.Join(db.RootDir, today)
	path := fmt.Sprintf("%s-%05d", basePath, sequence)
	return path
}

func (db *TextDb) getNextSequence(date string) int {
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
