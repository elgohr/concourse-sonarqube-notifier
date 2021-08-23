package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	mockResponse = `{
		  "paging": {
		    "pageIndex": 1,
		    "pageSize": 100,
		    "total": 12
		  },
		  "analyses": [
		    {
		      "key": "AWKa7VV9drIzrRaH-p_z",
		      "date": "2018-04-06T14:27:06+0200",
		      "events": [
		        {
		          "key": "AWKa7VsYdrIzrRaH-p_0",
		          "category": "VERSION",
		          "name": "0.0.1-SNAPSHOT"
		        }
		      ]
		    },
		    {
		      "key": "AWKQ3B6rdrIzrRaH-Rt3",
		      "date": "2018-04-04T15:32:28+0200",
		      "events": []
		    },
		    {
		      "key": "AWJhuKRVdrIzrRaH-JD8",
		      "date": "2018-03-26T11:51:30+0200",
		      "events": [
		        {
		          "key": "AWJhuKoddrIzrRaH-JD-",
		          "category": "QUALITY_GATE",
		          "name": "Green (was Red)",
		          "description": ""
		        }
		      ]
		    },
		    {
		      "key": "AWJOESP5NZwlownmr1uo",
		      "date": "2018-03-22T15:15:48+0100",
		      "events": []
		    },
		    {
		      "key": "AWIFz6Qd0iGqzMJL9y73",
		      "date": "2018-03-08T14:31:37+0100",
		      "events": []
		    }
		  ]
		}`
	suffix = "/api/project_analyses/search?project=my%3Acomponent"
)

func TestReturnsVersionsOfLastResult(t *testing.T) {
	stdin := &bytes.Buffer{}
	stdout := &bytes.Buffer{}

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

	stdin.WriteString(fmt.Sprintf(`{
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

	if err := run(stdin, stdout); err != nil {
		t.Error(err)
	}

	if !called {
		t.Error("Didn't call the remote service")
	}

	expectedResponse := `[{"timestamp":"2018-03-08T14:31:37+0100"},{"timestamp":"2018-03-22T15:15:48+0100"},{"timestamp":"2018-03-26T11:51:30+0200"},{"timestamp":"2018-04-04T15:32:28+0200"},{"timestamp":"2018-04-06T14:27:06+0200"}]`
	response := make([]byte, len(expectedResponse), len(expectedResponse))
	if _, err := stdout.Read(response); err != nil {
		t.Error(err)
	}
	if string(response) != string(expectedResponse) {
		t.Errorf("Expected content to be %v, but was %v", string(expectedResponse), string(response))
	}
}

func TestRequestsTheCorrectUrl(t *testing.T) {
	stdin := &bytes.Buffer{}
	stdout := &bytes.Buffer{}

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

	stdin.WriteString(fmt.Sprintf(`{
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

	if err := run(stdin, stdout); err != nil {
		t.Error(err)
	}

	if !called {
		t.Error("Didn't call the remote service")
	}
}

func TestAddsAuthenticationToTheRequest(t *testing.T) {
	stdin := &bytes.Buffer{}
	stdout := &bytes.Buffer{}

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

	stdin.WriteString(fmt.Sprintf(`{
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

	if err := run(stdin, stdout); err != nil {
		t.Error(err)
	}

	if !called {
		t.Error("Didn't call the remote service")
	}
}

func TestReturnsErrorIfContentCouldNotBeFetched(t *testing.T) {
	stdin := &bytes.Buffer{}
	stdout := &bytes.Buffer{}

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer s.Close()

	stdin.WriteString(fmt.Sprintf(`{
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

	if err := run(stdin, stdout); err == nil {
		t.Error("Expected error to occure, but didn't")
	}
}

func TestReturnsErrorIfUnauthorized(t *testing.T) {
	stdin := &bytes.Buffer{}
	stdout := &bytes.Buffer{}

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer s.Close()

	stdin.WriteString(fmt.Sprintf(`{
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

	if err := run(stdin, stdout); err == nil {
		t.Error("Expected error to occure, but didn't")
	}
}

func TestReturnsWhenErrorDuringCall(t *testing.T) {
	stdin := &bytes.Buffer{}
	stdout := &bytes.Buffer{}

	stdin.WriteString(`{
			"source": {
    			"target": "http://localhost",
				"sonartoken": "token",
    			"component": "my:component",
    			"metrics": "ncloc,complexity,violations,coverage"
  			},
  			"version": {
				"ref": "61cebf"
			}
		}`)

	if err := run(stdin, stdout); err == nil {
		t.Error("Expected error to occure, but didn't")
	}
}

func TestErrorsWhenTargetIsMissing(t *testing.T) {
	stdin := &bytes.Buffer{}
	stdout := &bytes.Buffer{}

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

	err := run(stdin, stdout)
	if err == nil || err.Error() != "mandatory field is missing" {
		t.Errorf("Expected error to occure, but was %v", err)
	}
}

func TestErrorsWhenComponentIsMissing(t *testing.T) {
	stdin := &bytes.Buffer{}
	stdout := &bytes.Buffer{}

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

	err := run(stdin, stdout)
	if err == nil || err.Error() != "mandatory field is missing" {
		t.Errorf("Expected error to occure, but was %v", err)
	}
}

func TestErrorsWhenMetricsAreMissing(t *testing.T) {
	stdin := &bytes.Buffer{}
	stdout := &bytes.Buffer{}

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

	err := run(stdin, stdout)
	if err == nil || err.Error() != "mandatory field is missing" {
		t.Errorf("Expected error to occure, but was %v", err)
	}
}

func TestErrorsWhenSonartokenIsMissing(t *testing.T) {
	stdin := &bytes.Buffer{}
	stdout := &bytes.Buffer{}

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

	err := run(stdin, stdout)
	if err == nil || err.Error() != "mandatory field is missing" {
		t.Errorf("Expected error to occure, but was %v", err)
	}
}

func getBasicHeader(authToken string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(authToken+":"))
}
