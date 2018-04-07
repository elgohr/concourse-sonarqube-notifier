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

	BeforeEach(func() {
		stdin = new(bytes.Buffer)
		stdout = new(bytes.Buffer)
		stderr = new(bytes.Buffer)
		fakeResultSource = new(sharedfakes.FakeResultSource)
	})

	It("returns the hash of the response content as the latest version", func() {
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
		fakeResponse := make([]byte, 36, 36)
		fakeResultSource.GetResultReturns(fakeResponse, nil)

		err := run(stdin, stdout, fakeResultSource)

		Expect(err).NotTo(HaveOccurred())

		Expect(fakeResultSource.GetResultCallCount()).To(Equal(1))
		url, token, component, metrics := fakeResultSource.GetResultArgsForCall(0)
		Expect(url).To(Equal("https://my.sonar.server"))
		Expect(token).To(Equal("token"))
		Expect(component).To(Equal("my:component"))
		Expect(metrics).To(Equal("ncloc,complexity,violations,coverage"))


		expectedResponse := `[{"ref":"81684c2e68ade2cd4bf9f2e8a67dd4fe"}]`
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
