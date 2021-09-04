package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	GetChapterUri      string = "/api/chapters/%d"
	SearchEncounterUri        = "/api/entities/encounters/search"
)

type OpApi struct {
	host string
}

func NewOpApi(host string) *OpApi {
	return &OpApi{host: host}
}

func (e *OpApi) GetChapter(number uint) (*Chapter, error) {
	path := fmt.Sprintf(GetChapterUri, number)
	resp, err := e.request(http.MethodGet, path, nil)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Invalid response")
	}

	chapterResponse := new(ChapterResponse)

	return &chapterResponse.Chapter, json.Unmarshal(body, &chapterResponse)
}

type SearchRequest struct {
	Entities []string `json:"entities"`
	Type     string   `json:"type"`
}

func (e *OpApi) SearchEncounter(entities []string) (*Encounter, error) {
	request := new(SearchRequest)
	request.Entities = entities
	request.Type = "characters"

	resp, err := e.request(http.MethodPost, SearchEncounterUri, request)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Invalid response")
	}

	searchResponse := new(SearchResponse)

	return &searchResponse.Encounter, json.Unmarshal(body, &searchResponse)
}

func (e *OpApi) request(method string, path string, body interface{}) (*http.Response, error) {

	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, e.host+path, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	return client.Do(req)
}
