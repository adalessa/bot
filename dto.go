package main

import (
	"encoding/json"
	"strings"
	"time"
)

type ChapterResponse struct {
	Chapter Chapter `json:"data"`
}

type Chapter struct {
	Number       uint            `json:"number"`
	ReleaseDate  JsonReleaseDate `json:"release_date"`
	Title        string          `json:"title"`
	Cover        Cover           `json:"cover"`
	ShortSummary string          `json:"short_summary"`
	Summary      string          `json:"summary"`
	Links        []Link          `json:"links"`
}

type Cover struct {
	Image string `json:"image"`
	Text  string `json:"text"`
}

type Link struct {
	Site string `json:"site"`
	Url  string `json:"url"`
}

type SearchResponse struct {
	Encounter Encounter `json:"data"`
}
type Encounter struct {
	Times    uint      `json:"times"`
	Entities []Entity  `json:"entities"`
	Chapters []Chapter `json:"chapters"`
}

type Entity struct {
	Name    string  `json:"name"`
	Wiki    string  `json:"wiki"`
	Aliases []Alias `json:"aliases"`
}

type Alias struct {
	Name string `json:"name"`
}

type JsonReleaseDate time.Time

// Implement Marshaler and Unmarshaler interface
func (j *JsonReleaseDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*j = JsonReleaseDate(t)
	return nil
}

func (j JsonReleaseDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(j)
}

// Maybe a Format function for printing your date
func (j JsonReleaseDate) Format(s string) string {
	t := time.Time(j)
	return t.Format(s)
}
