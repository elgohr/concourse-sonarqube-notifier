package main

import (
	"bytes"
	"errors"
	"github.com/elgohr/concourse-sonarqube-notifier/assets/shared/sharedfakes"
	"testing"
)

const mockResponse = `{
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

var (
	stdin            *bytes.Buffer
	stdout           *bytes.Buffer
	fakeResultSource *sharedfakes.FakeResultSource
)

func setup() {
	stdin = &bytes.Buffer{}
	stdout = &bytes.Buffer{}
	fakeResultSource = &sharedfakes.FakeResultSource{}
}

func TestReturnsVersionsOfLastResult(t *testing.T) {
	setup()

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
	fakeResultSource.GetVersionsReturns([]byte(mockResponse), nil)

	if err := run(stdin, stdout, fakeResultSource); err != nil {
		t.Error(err)
	}

	count := fakeResultSource.GetVersionsCallCount()
	if count != 1 {
		t.Errorf("Expected GetVersions to be called 1 time, but was %v times", count)
	}

	url, token, component := fakeResultSource.GetVersionsArgsForCall(0)
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

	expectedResponse := `[{"timestamp":"2018-03-08T14:31:37+0100"},{"timestamp":"2018-03-22T15:15:48+0100"},{"timestamp":"2018-03-26T11:51:30+0200"},{"timestamp":"2018-04-04T15:32:28+0200"},{"timestamp":"2018-04-06T14:27:06+0200"}]`
	response := make([]byte, len(expectedResponse), len(expectedResponse))
	if _, err := stdout.Read(response); err != nil {
		t.Error(err)
	}
	if string(response) != string(expectedResponse) {
		t.Errorf("Expected content to be %v, but was %v", string(expectedResponse), string(response))
	}
}

func TestReturnsErrorIfContentCouldNotBeFetched(t *testing.T) {
	setup()

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
	fakeResultSource.GetResultReturns(nil, errors.New("something serious"))

	if err := run(stdin, stdout, fakeResultSource); err == nil {
		t.Error("Expected error to occure, but didn't")
	}
}

func TestErrorsWhenTargetIsMissing(t *testing.T) {
	setup()

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

	err := run(stdin, stdout, fakeResultSource)
	if err.Error() != "mandatory field is missing" {
		t.Errorf("Expected error to occure, but was %v", err)
	}
}

func TestErrorsWhenComponentIsMissing(t *testing.T) {
	setup()

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

	err := run(stdin, stdout, fakeResultSource)
	if err.Error() != "mandatory field is missing" {
		t.Errorf("Expected error to occure, but was %v", err)
	}
}

func TestErrorsWhenMetricsAreMissing(t *testing.T) {
	setup()

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

	err := run(stdin, stdout, fakeResultSource)
	if err.Error() != "mandatory field is missing" {
		t.Errorf("Expected error to occure, but was %v", err)
	}
}

func TestErrorsWhenSonartokenIsMissing(t *testing.T) {
	setup()

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

	err := run(stdin, stdout, fakeResultSource)
	if err.Error() != "mandatory field is missing" {
		t.Errorf("Expected error to occure, but was %v", err)
	}
}
