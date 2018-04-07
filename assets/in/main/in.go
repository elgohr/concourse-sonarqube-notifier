package main

import (
	"encoding/json"
	"errors"
	. "github.com/concourse-sonarqube-notifier/assets/shared"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type InRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

type InResponse struct {
	Version Version `json:"version"`
}

func main() {
	downloadDir := os.Args[1]
	err := run(os.Stdin, os.Stdout, downloadDir, new(Sonarqube))
	if HasError(err) {
		log.Fatalln(err)
		os.Exit(1)
	}
}

func run(stdIn io.Reader, stdOut io.Writer, downloadDir string, resultSource ResultSource) error {
	var input InRequest
	err := json.NewDecoder(stdIn).Decode(&input)
	if HasError(err) {
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
	if HasError(err) {
		return err
	}

	destinationPath := filepath.Join(downloadDir, "result.json")
	ioutil.WriteFile(destinationPath, result, os.ModePerm)

	response := InResponse{
		input.Version,
	}
	err = json.NewEncoder(stdOut).Encode(response)
	if HasError(err) {
		return err
	}

	return nil
}
