package textdb

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

func logInfo(message string, parameter string) {
	log.Printf("textdb: %s %s", message, parameter)
}

func logError(message string, parameter string, err error) {
	log.Printf("textdb: %s %s. ERROR: %s", message, parameter, err)
}

func now() string {
	const time_format_now string = "2006-01-02 15:04:05.000" // yyyy-mm-dd hh:mm:ss.xxx
	return time.Now().Format(time_format_now)
}

func today() string {
	const time_format_today string = "2006-01-02" // yyyy-mm-dd
	return time.Now().Format(time_format_today)
}

func validId(id string) error {
	for _, c := range id {
		isDigit := c >= '0' && c <= '9'
		valid := isDigit || c == '-'
		if !valid {
			return errors.New(fmt.Sprintf("Invalid Id: %s", id))
		}
	}
	return nil
}

func slug(title string) string {
	slug := strings.Trim(title, " ")
	slug = strings.ToLower(slug)
	slug = strings.Replace(slug, "c#", "c-sharp", -1)
	var chars []rune
	for _, c := range slug {
		isAlpha := c >= 'a' && c <= 'z'
		isDigit := c >= '0' && c <= '9'
		if isAlpha || isDigit {
			chars = append(chars, c)
		} else {
			chars = append(chars, '-')
		}
	}
	slug = string(chars)

	// remove double dashes
	for strings.Index(slug, "--") > -1 {
		slug = strings.Replace(slug, "--", "-", -1)
	}

	if len(slug) == 0 || slug == "-" {
		return ""
	}

	// make sure we don't end with a dash
	if slug[len(slug)-1] == '-' {
		return slug[0 : len(slug)-1]
	}

	return slug
}
