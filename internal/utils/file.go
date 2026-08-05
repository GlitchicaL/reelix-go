package utils

import (
	"encoding/xml"
	"io"
	"os"
	"reelix-go/internal/db"
)

func IsDirEmpty(path string) (bool, error) {
	f, err := os.Open(path)

	if err != nil {
		return true, err
	}

	defer f.Close()

	_, err = f.ReadDir(1)

	if err == io.EOF {
		return true, err
	}

	return false, nil
}

type VideoMetadata struct {
	Title  string     `xml:"title"`
	Studio string     `xml:"studio"`
	Tags   []string   `xml:"tag"`
	Actors []db.Actor `xml:"actor"`
}

func ParseNfoFile(nfoPath string) (VideoMetadata, error) {
	data, err := os.ReadFile(nfoPath)

	if err != nil {
		return VideoMetadata{}, err
	}

	var metadata VideoMetadata

	err = xml.Unmarshal(data, &metadata)

	if err != nil {
		return VideoMetadata{}, err
	}

	return metadata, nil
}
