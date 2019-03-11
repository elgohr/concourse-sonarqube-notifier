package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

const (
	mockResponse = `{
		  "component": {
		    "id": "AWH_6osdce3G0HojaCW1",
		    "key": "my:component",
		    "name": "component-name",
		    "qualifier": "TRK",
		    "measures": [
		      {
		        "metric": "violations",
		        "value": "5",
		        "periods": [
		          {
		            "index": 1,
		            "value": "-6"
		          }
		        ]
		      },
		      {
		        "metric": "coverage",
		        "value": "91.2",
		        "periods": [
		          {
		            "index": 1,
		            "value": "40.5"
		          }
		        ]
		      },
		      {
		        "metric": "complexity",
		        "value": "84",
		        "periods": [
		          {
		            "index": 1,
		            "value": "21"
		          }
		        ]
		      },
		      {
		        "metric": "ncloc",
		        "value": "795",
		        "periods": [
		          {
		            "index": 1,
		            "value": "270"
		          }
		        ]
		      }
		    ]
		  }
		}`
	suffix = "/api/measures/component?component=my%3Acomponent&metricKeys=ncloc%2Ccomplexity%2Cviolations%2Ccoverage"
)

func setup(t *testing.T) (stdIn *bytes.Buffer, stdOut *bytes.Buffer, tmpDir string) {
	var err error
	tmpDir, err = ioutil.TempDir("", "concourse-sonarqube")
	if err != nil {
		t.Error(err)
	}
	return &bytes.Buffer{}, &bytes.Buffer{}, tmpDir
	
}

func TestWritesContentIntoDownloadDirectory(t *testing.T) {
	stdIn, stdOut, tmpDir := setup(t)

	var called bool
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		if _, err := w.Write([]byte(mockResponse)); err != nil {
			t.Error(err)
		}
	}))
	defer s.Close()

	stdIn.WriteString(fmt.Sprintf(`{
			"source": {
    			"target": "%v",
				"sonartoken": "token",
    			"component": "my:component",
    			"metrics": "ncloc,complexity,violations,coverage"
  			},
  			"version": {
				"ref": "61cebf"
			}
		}`, s.URL))

	if err := run(stdIn, stdOut, tmpDir); err != nil {
		t.Error(err)
	}

	if !called {
		t.Error("Didn't call the remote service")
	}

	expectedDestination, err := ioutil.ReadDir(tmpDir)
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

	fullPath := filepath.Join(tmpDir, "result.json")
	content, err := ioutil.ReadFile(fullPath)
	if err != nil {
		t.Error(err)
	}
	if string(content) != mockResponse {
		t.Errorf("Expected content to be %v, but was %v", mockResponse, string(content))
	}
}

func TestWritesVersionTostdOut(t *testing.T) {
	stdIn, stdOut, tmpDir := setup(t)

	var called bool
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		if _, err := w.Write([]byte(mockResponse)); err != nil {
			t.Error(err)
		}
	}))
	defer s.Close()

	stdIn.WriteString(fmt.Sprintf(`{
			"source": {
    			"target": "%v",
				"sonartoken": "token",
    			"component": "my:component",
    			"metrics": "ncloc,complexity,violations,coverage"
  			},
  			"version": {
				"ref": "61cebf"
			}
		}`, s.URL))

	err := run(stdIn, stdOut, tmpDir)
	if err != nil {
		t.Error(err)
	}

	if !called {
		t.Error("Didn't call the remote service")
	}

	expectedResponse := `{"version":{"ref":"61cebf"}}`
	response := make([]byte, len(expectedResponse), len(expectedResponse))
	if _, err := stdOut.Read(response); err != nil {
		t.Error(err)
	}
	if string(response) != string(expectedResponse) {
		t.Errorf("Expected content to be %v, but was %v", expectedResponse, response)
	}
}

func TestRequestsTheCorrectUrl(t *testing.T) {
	stdIn, stdOut, tmpDir := setup(t)

	var called bool
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		if r.URL.String() != suffix {
			t.Errorf("Expected %v, but got %v", suffix, r.URL.String())
		}
		if _, err := w.Write([]byte(mockResponse)); err != nil {
			t.Error(err)
		}
	}))
	defer s.Close()

	stdIn.WriteString(fmt.Sprintf(`{
			"source": {
    			"target": "%v",
				"sonartoken": "token",
    			"component": "my:component",
    			"metrics": "ncloc,complexity,violations,coverage"
  			},
  			"version": {
				"ref": "61cebf"
			}
		}`, s.URL))

	err := run(stdIn, stdOut, tmpDir)
	if err != nil {
		t.Error(err)
	}

	if !called {
		t.Error("Didn't call the remote service")
	}
}

func TestAddsAuthenticationToTheRequest(t *testing.T) {
	stdIn, stdOut, tmpDir := setup(t)

	authToken := "token"

	var called bool
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		gotAuth := r.Header.Get("Authorization")
		expAuth := getBasicHeader(authToken)
		if gotAuth != expAuth {
			t.Errorf("Expected %v, but got %v", expAuth, gotAuth)
		}
		if _, err := w.Write([]byte(mockResponse)); err != nil {
			t.Error(err)
		}
	}))
	defer s.Close()

	stdIn.WriteString(fmt.Sprintf(`{
			"source": {
    			"target": "%v",
				"sonartoken": "%v",
    			"component": "my:component",
    			"metrics": "ncloc,complexity,violations,coverage"
  			},
  			"version": {
				"ref": "61cebf"
			}
		}`, s.URL, authToken))

	err := run(stdIn, stdOut, tmpDir)
	if err != nil {
		t.Error(err)
	}

	if !called {
		t.Error("Didn't call the remote service")
	}
}

func TestReturnsErrorIfContentCouldNotBeFetched(t *testing.T) {
	stdIn, stdOut, tmpDir := setup(t)

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer s.Close()

	stdIn.WriteString(fmt.Sprintf(`{
			"source": {
    			"target": "%v",
				"sonartoken": "token",
    			"component": "my:component",
    			"metrics": "ncloc,complexity,violations,coverage"
  			},
  			"version": {
				"ref": "61cebf"
			}
		}`, s.URL))

	if err := run(stdIn, stdOut, tmpDir); err == nil {
		t.Error("Expected error, but didn't error")
	}
}

func TestReturnsErrorIfUnauthorized(t *testing.T) {
	stdIn, stdOut, tmpDir := setup(t)

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer s.Close()

	stdIn.WriteString(fmt.Sprintf(`{
			"source": {
    			"target": "%v",
				"sonartoken": "token",
    			"component": "my:component",
    			"metrics": "ncloc,complexity,violations,coverage"
  			},
  			"version": {
				"ref": "61cebf"
			}
		}`, s.URL))

	if err := run(stdIn, stdOut, tmpDir); err == nil {
		t.Error("Expected error, but didn't error")
	}
}

func TestErrorsWhenTargetIsMissing(t *testing.T) {
	stdIn, stdOut, tmpDir := setup(t)

	stdIn.WriteString(`{
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

	err := run(stdIn, stdOut, tmpDir)
	if err.Error() != "mandatory field is missing" {
		t.Errorf("Expected error to occure, but was %v", err)
	}
}

func TestErrorsWhenComponentIsMissing(t *testing.T) {
	stdIn, stdOut, tmpDir := setup(t)

	stdIn.WriteString(`{
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

	err := run(stdIn, stdOut, tmpDir)
	if err.Error() != "mandatory field is missing" {
		t.Errorf("Expected error to occure, but was %v", err)
	}
}

func TestErrorsWhenMetricsAreMissing(t *testing.T) {
	stdIn, stdOut, tmpDir := setup(t)

	stdIn.WriteString(`{
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

	err := run(stdIn, stdOut, tmpDir)
	if err.Error() != "mandatory field is missing" {
		t.Errorf("Expected error to occure, but was %v", err)
	}
}

func TestErrorsWhenSonartokenIsMissing(t *testing.T) {
	stdIn, stdOut, tmpDir := setup(t)

	stdIn.WriteString(`{
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

	err := run(stdIn, stdOut, tmpDir)
	if err.Error() != "mandatory field is missing" {
		t.Errorf("Expected error to occure, but was %v", err)
	}
}

func getBasicHeader(authToken string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(authToken+":"))
}
