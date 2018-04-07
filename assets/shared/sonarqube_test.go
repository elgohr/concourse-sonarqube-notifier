package shared_test

import (
	"encoding/base64"
	"github.com/concourse-sonarqube-notifier/assets/shared"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/jarcoal/httpmock.v1"
	"log"
	"net/http"
)

var _ = Describe("Sonarqube", func() {

	var sonarqube *shared.Sonarqube
	const component = "my:component"
	const metrics = "ncloc,complexity,violations,coverage"
	const testUrl = "http://localhost/my.sonar.server"
	const suffix = "/api/measures/component?component=my%3Acomponent&metricKeys=ncloc%2Ccomplexity%2Cviolations%2Ccoverage"
	const authToken = "SONAR_TOKEN"
	const mockResponse = `{
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

	BeforeEach(func() {
		sonarqube = new(shared.Sonarqube)
		httpmock.Activate()
		httpmock.RegisterResponder("GET", testUrl + suffix,
			func(req *http.Request) (*http.Response, error) {
				if isNotAuthenticated(req, authToken) {
					errMsg := "Expected:" + getBasicHeader(authToken) + " but was " + req.Header.Get("Authorization")
					log.Print(errMsg)
					return httpmock.NewStringResponse(401, errMsg), nil
				}
				resp := httpmock.NewStringResponse(200, mockResponse)
				return resp, nil
			})
	})

	AfterEach(func() {
		httpmock.DeactivateAndReset()
	})

	It("returns the body of the remote call with correct authentication", func() {
		response, err := sonarqube.GetResult(testUrl, authToken, component, metrics)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(response)).To(Equal(mockResponse))
	})

	It("returns an error if the authentication is invalid", func() {
		_, err := sonarqube.GetResult(testUrl, "INVALID", component, metrics)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("StatusCode:401"))
	})

	It("returns an error if the content could not be found", func() {
		_, err := sonarqube.GetResult("http://localhost/INVALID", authToken, component, metrics)
		Expect(err).To(HaveOccurred())
	})

})

func isNotAuthenticated(req *http.Request, authToken string) bool {
	return req.Header.Get("Authorization") != getBasicHeader(authToken)
}

func getBasicHeader(authToken string) string{
	return "Basic "+base64.StdEncoding.EncodeToString([]byte(authToken+":"))
}
