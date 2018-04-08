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
	const testUrl = "http://localhost/my.sonar.server"
	const authToken = "SONAR_TOKEN"

	Describe("GetResult", func() {
		const metrics = "ncloc,complexity,violations,coverage"
		const suffix = "/api/measures/component?component=my%3Acomponent&metricKeys=ncloc%2Ccomplexity%2Cviolations%2Ccoverage"
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
			httpmock.RegisterResponder("GET", testUrl+suffix,
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

	Describe("GetVersions", func() {
		const suffix = "/api/project_analyses/search?project=my%3Acomponent"
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

		BeforeEach(func() {
			sonarqube = new(shared.Sonarqube)
			httpmock.Activate()
			httpmock.RegisterResponder("GET", testUrl+suffix,
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
			response, err := sonarqube.GetVersions(testUrl, authToken, component)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(response)).To(Equal(mockResponse))
		})

		It("returns an error if the authentication is invalid", func() {
			_, err := sonarqube.GetVersions(testUrl, "INVALID", component)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("StatusCode:401"))
		})

		It("returns an error if the content could not be found", func() {
			_, err := sonarqube.GetVersions("http://localhost/INVALID", authToken, component)
			Expect(err).To(HaveOccurred())
		})

	})

})

func isNotAuthenticated(req *http.Request, authToken string) bool {
	return req.Header.Get("Authorization") != getBasicHeader(authToken)
}

func getBasicHeader(authToken string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(authToken+":"))
}
