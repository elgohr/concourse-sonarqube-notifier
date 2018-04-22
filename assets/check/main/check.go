package main

import (
	"encoding/json"
	"errors"
	. "github.com/concourse-sonarqube-notifier/assets/shared"
	"io"
	"log"
	"os"
)

type CheckRequest struct {
	Source Source `json:"source"`
}

type CheckResponse []Version

type SonarResponse struct {
	Analyses []Analyses `json:"analyses"`
}

type Analyses struct {
	Key  string `json:"key"`
	Date string `json:"date"`
}

func main() {
	err := run(os.Stdin, os.Stdout, new(Sonarqube))
	if HasError(err) {
		log.Fatalln(err)
		os.Exit(1)
	}
}

func run(stdIn io.Reader, stdOut io.Writer, resultSource ResultSource) error {
	var (
		input    CheckRequest
		response SonarResponse
	)
	err := json.NewDecoder(stdIn).Decode(&input)
	if HasError(err) {
		return err
	}

	if !input.Source.Valid() {
		return errors.New("mandatory field is missing")
	}

	result, err := resultSource.GetVersions(
		input.Source.Target,
		input.Source.SonarToken,
		input.Source.Component,
	)
	if HasError(err) {
		return err
	}

	err = json.Unmarshal(result, &response)
	if HasError(err) {
		return err
	}

	var remoteVersions CheckResponse

	for _, a := range response.Analyses {
		remoteVersions = append([]Version{{"timestamp":a.Date}}, remoteVersions...)
	}

	err = json.NewEncoder(stdOut).Encode(remoteVersions)
	if HasError(err) {
		return err
	}
	return nil
}
