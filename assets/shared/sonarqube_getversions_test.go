package shared_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	getVersionsSuffix       = "/api/project_analyses/search?project=my%3Acomponent"
	getVersionsMockResponse = `{
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
)

func TestGetVersionsReturnsBodyOfRemoteCall(t *testing.T) {
	var called bool
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		if _, err := w.Write([]byte(getVersionsMockResponse)); err != nil {
			t.Error(err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer s.Close()

	res, err := sonarqube.GetVersions(s.URL, authToken, component)
	if err != nil {
		t.Error(err)
	}
	if !called {
		t.Error("Didn't call the mock")
	}
	if string(res) != getVersionsMockResponse {
		t.Errorf("Expected %v, but got %v", getVersionsMockResponse, string(res))
	}
}

func TestGetVersionsRequestCorrectUrl(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() != getVersionsSuffix {
			t.Errorf("Expected %v, but got %v", getVersionsSuffix, r.URL.String())
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer s.Close()

	_, _ = sonarqube.GetVersions(s.URL, authToken, component)
}

func TestGetVersionsAddsAuthenticationToTheRequest(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth := r.Header.Get("Authorization")
		expAuth := getBasicHeader(authToken)
		if gotAuth != expAuth {
			t.Errorf("Expected %v, but got %v", expAuth, gotAuth)
		}
	}))
	defer s.Close()

	_, _ = sonarqube.GetVersions(s.URL, authToken, component)
}

func TestGetVersionsReturnsErrorIfAuthenticationFails(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer s.Close()

	_, err := sonarqube.GetVersions(s.URL, "INVALID", component)
	if err == nil {
		t.Error("Expected to get an error, but didnt")
	}
}

func TestGetVersionsReturnsErrorIfContentCouldNotBeFound(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer s.Close()

	_, err := sonarqube.GetVersions(s.URL, authToken, "NOT_PRESENT")
	if err == nil {
		t.Error("Expected to get an error, but didnt")
	}
}
