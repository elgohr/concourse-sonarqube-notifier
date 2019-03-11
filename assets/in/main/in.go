package main

import (
	"encoding/json"
	"errors"
	"github.com/elgohr/concourse-sonarqube-notifier/assets/shared"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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
	if err := run(os.Stdin, os.Stdout, downloadDir, new(shared.Sonarqube));
		err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
}

func run(stdIn io.Reader, stdOut io.Writer, downloadDir string, resultSource shared.ResultSource) error {
	var input InRequest
	if err := json.NewDecoder(stdIn).Decode(&input); err != nil {
		return err
	}

	if !input.Source.Valid() {
		return errors.New("mandatory field is missing")
	}

	result, err := resultSource.GetResult(
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
