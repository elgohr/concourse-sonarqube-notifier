package main

import (
	"bytes"
	"errors"
	"github.com/elgohr/concourse-sonarqube-notifier/assets/shared/sharedfakes"
	"io/ioutil"
	"path/filepath"
	"testing"
)

var (
	stdin            *bytes.Buffer
	stdout           *bytes.Buffer
	fakeResultSource *sharedfakes.FakeResultSource
	tmpDownloadDir   string
)

func setup(t *testing.T) {
	stdin = &bytes.Buffer{}
	stdout = &bytes.Buffer{}
	fakeResultSource = &sharedfakes.FakeResultSource{}
	var err error
	tmpDownloadDir, err = ioutil.TempDir("", "concourse-sonarqube")
	if err != nil {
		t.Error(err)
	}
}

func TestWritesContentIntoDownloadDirectory(t *testing.T) {
	setup(t)

	stdin.WriteString(`{
			"source": {
    			"target": "https://my.sonar.server",
				"sonartoken": "token",
    			"component": "my:component",
    			"metrics": "ncloc,complexity,violations,coverage"
  			},
  			"version": {
				"ref": "61cebf"
			}
		}`)

	fakeResponse := make([]byte, 36, 36)
	fakeResultSource.GetResultReturns(fakeResponse, nil)

	if err := run(stdin, stdout, tmpDownloadDir, fakeResultSource); err != nil {
		t.Error(err)
	}

	count := fakeResultSource.GetResultCallCount()
	if count != 1 {
		t.Errorf("Expected GetResult to be called 1 time, but was %v times", count)
	}

	url, token, component, metrics := fakeResultSource.GetResultArgsForCall(0)
	expUrl := "https://my.sonar.server"
	if url != expUrl {
		t.Errorf("Expected url to be %v, but was %v", expUrl, url)
	}
	expToken := "token"
	if token != expToken {
		t.Errorf("Expected token to be %v, but was %v", expToken, url)
	}
	expComponent := "my:component"
	if component != expComponent {
		t.Errorf("Expected component to be %v, but was %v", expComponent, url)
	}
	expMetrics := "ncloc,complexity,violations,coverage"
	if metrics != expMetrics {
		t.Errorf("Expected metrics to be %v, but was %v", expMetrics, url)

	}

	expectedDestination, err := ioutil.ReadDir(tmpDownloadDir)
	if err != nil {
		t.Error(err)
	}
	if len(expectedDestination) < 1 {
		t.Errorf("Expected a file to be created, but was not")
	}
	name := expectedDestination[0].Name()
	if name != "result.json" {
		t.Errorf("Expected the file to be named correctly, but was %v", name)
	}

	fullPath := filepath.Join(tmpDownloadDir, "result.json")
	content, err := ioutil.ReadFile(fullPath)
	if err != nil {
		t.Error(err)
	}
	if string(content) != string(fakeResponse) {
		t.Errorf("Expected content to be %v, but was %v", fakeResponse, content)
	}
}

func TestWritesVersionToStdout(t *testing.T) {
	setup(t)

	stdin.WriteString(`{
			"source": {
    			"target": "https://my.sonar.server",
				"sonartoken": "token",
    			"component": "my:component",
    			"metrics": "ncloc,complexity,violations,coverage"
  			},
  			"version": {
				"ref": "61cebf"
			}
		}`)
	fakeResponse := make([]byte, 36, 36)
	fakeResultSource.GetResultReturns(fakeResponse, nil)

	err := run(stdin, stdout, tmpDownloadDir, fakeResultSource)
	if err != nil {
		t.Error(err)
	}

	expectedResponse := `{"version":{"ref":"61cebf"}}`
	response := make([]byte, len(expectedResponse), len(expectedResponse))
	if _, err := stdout.Read(response); err != nil {
		t.Error(err)
	}
	if string(response) != string(expectedResponse) {
		t.Errorf("Expected content to be %v, but was %v", expectedResponse, response)
	}
}

func TestReturnsErrorIfContentCouldNotBeFetched(t *testing.T) {
	setup(t)

	stdin.WriteString(`{
			"source": {
    			"target": "https://my.sonar.server",
				"sonartoken": "token",
    			"component": "my-component",
    			"metrics": "ncloc,complexity,violations,coverage"
  			},
  			"version": {
				"ref": "61cebf"
			}
		}`)
	expErr := errors.New("something serious")
	fakeResultSource.GetResultReturns(nil, expErr)

	if err := run(stdin, stdout, tmpDownloadDir, fakeResultSource); err != expErr {
		t.Errorf("Expected error to be %v, but was %v", expErr, err)
	}

	if fakeResultSource.GetResultCallCount() < 1 {
		t.Error("Expected GetResult to be called, but wasn't")
	}
}

func TestErrorsWhenTargetIsMissing(t *testing.T) {
	setup(t)

	stdin.WriteString(`{
				"source": {
    				"missing_target": "https://my.sonar.server",
    				"sonartoken": "token",
    				"component": "my-component",
    				"metrics": "ncloc,complexity,violations,coverage"
  				},
  				"version": { 
					"ref": "61cebf" 
				}
			}`)
	fakeResponse := make([]byte, 36, 36)
	fakeResultSource.GetResultReturns(fakeResponse, nil)

	err := run(stdin, stdout, tmpDownloadDir, fakeResultSource)
	if err.Error() != "mandatory field is missing" {
		t.Errorf("Expected error to occure, but was %v", err)
	}
}

func TestErrorsWhenComponentIsMissing(t *testing.T) {
	setup(t)

	stdin.WriteString(`{
				"source": {
    				"target": "https://my.sonar.server",
    				"sonartoken": "token",
    				"missing_component": "my-component",
    				"metrics": "ncloc,complexity,violations,coverage"
  				},
  				"version": { 
					"ref": "61cebf" 
				}
			}`)
	fakeResponse := make([]byte, 36, 36)
	fakeResultSource.GetResultReturns(fakeResponse, nil)

	err := run(stdin, stdout, tmpDownloadDir, fakeResultSource)
	if err.Error() != "mandatory field is missing" {
		t.Errorf("Expected error to occure, but was %v", err)
	}
}

func TestErrorsWhenMetricsAreMissing(t *testing.T) {
	setup(t)

	stdin.WriteString(`{
				"source": {
    				"target": "https://my.sonar.server",
    				"sonartoken": "token",
    				"component": "my-component",
    				"missing_metrics": "ncloc,complexity,violations,coverage"
  				},
  				"version": { 
					"ref": "61cebf" 
				}
			}`)
	fakeResponse := make([]byte, 36, 36)
	fakeResultSource.GetResultReturns(fakeResponse, nil)

	err := run(stdin, stdout, tmpDownloadDir, fakeResultSource)
	if err.Error() != "mandatory field is missing" {
		t.Errorf("Expected error to occure, but was %v", err)
	}
}

func TestErrorsWhenSonartokenIsMissing(t *testing.T) {
	setup(t)

	stdin.WriteString(`{
				"source": {
    				"target": "https://my.sonar.server",
    				"missing_sonartoken": "token",
    				"component": "my-component",
    				"metrics": "ncloc,complexity,violations,coverage"
  				},
  				"version": { 
					"ref": "61cebf" 
				}
			}`)
	fakeResponse := make([]byte, 36, 36)
	fakeResultSource.GetResultReturns(fakeResponse, nil)

	err := run(stdin, stdout, tmpDownloadDir, fakeResultSource)
	if err.Error() != "mandatory field is missing" {
		t.Errorf("Expected error to occure, but was %v", err)
	}
}
