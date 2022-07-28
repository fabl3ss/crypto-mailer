package utils

import (
	"encoding/csv"
	"io"
	"os"
	"sort"
)

// ExistsInCsv returns sorted slice with
// inserted toFind if record not exists
// or nil if record already exists
func ExistsInCsv(path string, toFind string) ([]string, error) {
	f, err := os.OpenFile(
		path,
		os.O_RDONLY|os.O_CREATE,
		0644,
	)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var emails []string
	csvReader := csv.NewReader(f)

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		emails = append(emails, record[0])
	}

	index := sort.SearchStrings(emails, toFind)
	if index != len(emails) {
		if emails[index] == toFind {
			return nil, nil
		}
	}

	// insert into correct position
	emails = append(emails, "")
	copy(emails[index+1:], emails[index:])
	emails[index] = toFind

	return emails, nil
}
