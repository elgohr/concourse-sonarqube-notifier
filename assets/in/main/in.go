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
	"path/filepath"
	"strconv"
)

type InRequest struct {
	Source  shared.Source  `json:"source"`
	Version shared.Version `json:"version"`
}

type InResponse struct {
	Version shared.Version `json:"version"`
}

func main() {
	downloadDir := os.Args[1]
	if err := run(os.Stdin, os.Stdout, downloadDir);
		err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
}

func run(stdIn io.Reader, stdOut io.Writer, downloadDir string) error {
	var input InRequest
	if err := json.NewDecoder(stdIn).Decode(&input); err != nil {
		return err
	}

	if !input.Source.Valid() {
		return errors.New("mandatory field is missing")
	}

	result, err := getResult(
		input.Source.Target,
		input.Source.SonarToken,
		input.Source.Component,
		input.Source.Metrics,
	)
	if err != nil {
		return err
	}

	destinationPath := filepath.Join(downloadDir, "result.json")
	ioutil.WriteFile(destinationPath, result, os.ModePerm)

	return json.
		NewEncoder(stdOut).
		Encode(InResponse{
			Version: input.Version,
		})
}

func getResult(baseUrl string, authToken string, component string, metrics string) ([]byte, error) {
	fullUrl, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	fullUrl.Path += "/api/measures/component"
	parameters := url.Values{}
	parameters.Add("component", component)
	parameters.Add("metricKeys", metrics)
	fullUrl.RawQuery = parameters.Encode()

	req, err := http.NewRequest("GET", fullUrl.String(), nil)
	req.SetBasicAuth(authToken, "")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("Status " + strconv.Itoa(resp.StatusCode) + " : " + string(body))
	}
	return body, nil
}
