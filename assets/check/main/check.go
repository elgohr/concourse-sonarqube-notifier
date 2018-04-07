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
	Source  Source  `json:"source"`
}

type CheckResponse []Version

func main() {
	err := run(os.Stdin, os.Stdout, new(Sonarqube))
	if HasError(err) {
		log.Fatalln(err)
		os.Exit(1)
	}
}

func run(stdIn io.Reader, stdOut io.Writer, resultSource ResultSource) error {
	var input CheckRequest
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

	remoteVersion := make(map[string]string)
	remoteVersion["ref"] = Md5Hash(string(result))
	response := CheckResponse{
		remoteVersion,
	}
	err = json.NewEncoder(stdOut).Encode(response)
	if HasError(err) {
		return err
	}
	return nil
}
