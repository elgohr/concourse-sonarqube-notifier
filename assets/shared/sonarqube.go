package shared

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

//go:generate counterfeiter . ResultSource
type ResultSource interface {
	GetResult(url string, authToken string, component string, metrics string) ([]byte, error)
	GetVersions(baseUrl string, authToken string, component string) ([]byte, error)
}

type Sonarqube struct {
}

func (s *Sonarqube) GetResult(baseUrl string, authToken string, component string, metrics string) ([]byte, error) {
	fullUrl, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	fullUrl.Path += "/api/measures/component"
	parameters := url.Values{}
	parameters.Add("component", component)
	parameters.Add("metricKeys", metrics)
	fullUrl.RawQuery = parameters.Encode()

	req, err := http.NewRequest("GET", fullUrl.String(), nil)
	req.SetBasicAuth(authToken, "")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("Status " + strconv.Itoa(resp.StatusCode) + " : " + string(body))
	}
	return body, nil
}

func (s *Sonarqube) GetVersions(baseUrl string, authToken string, component string) ([]byte, error) {
	fullUrl, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	fullUrl.Path += "/api/project_analyses/search"
	parameters := url.Values{}
	parameters.Add("project", component)
	fullUrl.RawQuery = parameters.Encode()

	req, err := http.NewRequest("GET", fullUrl.String(), nil)
	req.SetBasicAuth(authToken, "")

	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("Status " + strconv.Itoa(resp.StatusCode) + " : " + string(body))
	}
	return body, nil
}
