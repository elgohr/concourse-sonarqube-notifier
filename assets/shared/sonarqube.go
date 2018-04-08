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
}

type Sonarqube struct {
}

func (s *Sonarqube) GetResult(baseUrl string, authToken string, component string, metrics string) ([]byte, error) {
	fullUrl, err := url.Parse(baseUrl)
	if HasError(err) {
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
	if resp.StatusCode == 401 {
		return nil, errors.New("StatusCode:" + strconv.Itoa(resp.StatusCode))
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
