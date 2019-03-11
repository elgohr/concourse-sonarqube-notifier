package shared_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	component             = "my:component"
	authToken             = "SONAR_TOKEN"
	metrics               = "ncloc,complexity,violations,coverage"
	getResultSuffix       = "/api/measures/component?component=my%3Acomponent&metricKeys=ncloc%2Ccomplexity%2Cviolations%2Ccoverage"
	getResultMockResponse = `{
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
)

func TestGetResultReturnsBodyOfRemoteCall(t *testing.T) {
	var called bool
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		if _, err := w.Write([]byte(getResultMockResponse)); err != nil {
			t.Error(err)
		}
	}))
	defer s.Close()

	res, err := sonarqube.GetResult(s.URL, authToken, component, metrics)
	if err != nil {
		t.Error(err)
	}
	if !called {
		t.Error("Didn't call the mock")
	}
	if string(res) != getResultMockResponse {
		t.Errorf("Expected %v, but got %v", getResultMockResponse, string(res))
	}
}

func TestGetResultRequestCorrectUrl(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() != getResultSuffix {
			t.Errorf("Expected %v, but got %v", getResultSuffix, r.URL.String())
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer s.Close()

	_, _ = sonarqube.GetResult(s.URL, authToken, component, metrics)
}

func TestAddsAuthenticationToTheRequest(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth := r.Header.Get("Authorization")
		expAuth := getBasicHeader(authToken)
		if gotAuth != expAuth {
			t.Errorf("Expected %v, but got %v", expAuth, gotAuth)
		}
	}))
	defer s.Close()

	_, _ = sonarqube.GetResult(s.URL, authToken, component, metrics)
}

func TestReturnsErrorIfAuthenticationFails(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer s.Close()

	_, err := sonarqube.GetResult(s.URL, "INVALID", component, metrics)
	if err == nil {
		t.Error("Expected to get an error, but didnt")
	}
}

func TestReturnsErrorIfContentCouldNotBeFound(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		errorMessage := `{"errors":[{"msg":"Component key 'NOT_PRESENT' not found"}]}`
		if _, err := w.Write([]byte(errorMessage)); err != nil {
			t.Error(err)
		}
	}))
	defer s.Close()

	_, err := sonarqube.GetResult(s.URL, authToken, "NOT_PRESENT", metrics)
	if err == nil {
		t.Error("Expected to get an error, but didnt")
	}
}
