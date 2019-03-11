package main

import (
	"encoding/json"
	"errors"
	"github.com/elgohr/concourse-sonarqube-notifier/assets/shared"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

type CheckRequest struct {
	Source shared.Source `json:"source"`
}

type CheckResponse []shared.Version

type SonarResponse struct {
	Analyses []Analyses `json:"analyses"`
}

type Analyses struct {
	Key  string `json:"key"`
	Date string `json:"date"`
}

func main() {
	if err := run(os.Stdin, os.Stdout); err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
}

func run(stdIn io.Reader, stdOut io.Writer) error {
	var (
		input    CheckRequest
		response SonarResponse
	)
	if err := json.NewDecoder(stdIn).Decode(&input); err != nil {
		return err
	}

	if !input.Source.Valid() {
		return errors.New("mandatory field is missing")
	}

	result, err := getVersions(
		input.Source.Target,
		input.Source.SonarToken,
		input.Source.Component,
	)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(result, &response); err != nil {
		return err
	}

	var remoteVersions CheckResponse
	for _, a := range response.Analyses {
		remoteVersions = append([]shared.Version{{"timestamp": a.Date}}, remoteVersions...)
	}

	return json.NewEncoder(stdOut).Encode(remoteVersions)
}

func getVersions(baseUrl string, authToken string, component string) ([]byte, error) {
	fullUrl, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	fullUrl.Path += "/api/project_analyses/search"
	parameters := url.Values{}
	parameters.Add("project", component)
	fullUrl.RawQuery = parameters.Encode()

	req, err := http.NewRequest("GET", fullUrl.String(), nil)
	req.SetBasicAuth(authToken, "")

	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("Status " + strconv.Itoa(resp.StatusCode) + " : " + string(body))
	}
	return body, nil
}
