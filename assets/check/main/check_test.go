package main

import (
	"bytes"
	"errors"
	"github.com/concourse-sonarqube-notifier/assets/shared/sharedfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Check", func() {
	var (
		stdin            *bytes.Buffer
		stdout           *bytes.Buffer
		stderr           *bytes.Buffer
		fakeResultSource *sharedfakes.FakeResultSource
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

	BeforeEach(func() {
		stdin = new(bytes.Buffer)
		stdout = new(bytes.Buffer)
		stderr = new(bytes.Buffer)
		fakeResultSource = new(sharedfakes.FakeResultSource)
	})

	It("returns the versions of the last results", func() {
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

		err := run(stdin, stdout, fakeResultSource)

		Expect(err).NotTo(HaveOccurred())

		Expect(fakeResultSource.GetVersionsCallCount()).To(Equal(1))
		url, token, component := fakeResultSource.GetVersionsArgsForCall(0)
		Expect(url).To(Equal("https://my.sonar.server"))
		Expect(token).To(Equal("token"))
		Expect(component).To(Equal("my:component"))


		expectedResponse := `[{"AWIFz6Qd0iGqzMJL9y73":"2018-03-08T14:31:37+0100"},{"AWJOESP5NZwlownmr1uo":"2018-03-22T15:15:48+0100"},{"AWJhuKRVdrIzrRaH-JD8":"2018-03-26T11:51:30+0200"},{"AWKQ3B6rdrIzrRaH-Rt3":"2018-04-04T15:32:28+0200"},{"AWKa7VV9drIzrRaH-p_z":"2018-04-06T14:27:06+0200"}]`
		response := make([]byte, len(expectedResponse), len(expectedResponse))
		stdout.Read(response)
		Expect(string(response)).To(Equal(expectedResponse))
	})

	It("returns an error if the content could not be fetched", func() {
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

		err := run(stdin, stdout, fakeResultSource)

		Expect(err).To(HaveOccurred())
	})

	Describe("mandatory check", func() {
		It("errors when target is missing", func() {
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

			Expect(err).To(HaveOccurred())
		})

		It("errors when component is missing", func() {
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

			Expect(err).To(HaveOccurred())
		})

		It("errors when metrics are missing", func() {
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

			Expect(err).To(HaveOccurred())
		})

		It("errors when sonartoken is missing", func() {
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

			Expect(err).To(HaveOccurred())
		})
	})
})
