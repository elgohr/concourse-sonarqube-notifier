package main

import (
	"encoding/json"
	"errors"
	"github.com/elgohr/concourse-sonarqube-notifier/assets/shared"
	"io"
	"log"
	"os"
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
	if err := run(os.Stdin, os.Stdout, &shared.Sonarqube{}); err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
}

func run(stdIn io.Reader, stdOut io.Writer, resultSource shared.ResultSource) error {
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

	result, err := resultSource.GetVersions(
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
