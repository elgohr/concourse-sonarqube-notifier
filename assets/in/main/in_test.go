package main

import (
	"bytes"
	"errors"
	"github.com/concourse-sonarqube-notifier/assets/shared/sharedfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"path/filepath"
)

var _ = Describe("In", func() {
	var (
		stdin            *bytes.Buffer
		stdout           *bytes.Buffer
		stderr           *bytes.Buffer
		fakeResultSource *sharedfakes.FakeResultSource
		tmpDownloadDir   string
	)

	BeforeEach(func() {
		stdin = new(bytes.Buffer)
		stdout = new(bytes.Buffer)
		stderr = new(bytes.Buffer)
		fakeResultSource = new(sharedfakes.FakeResultSource)
		var err error
		tmpDownloadDir, err = ioutil.TempDir("", "concourse-sonarqube")
		if err != nil {
			panic(1)
		}
	})

	It("writes the content into the download directory", func() {
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

		err := run(stdin, stdout, tmpDownloadDir, fakeResultSource)

		Expect(err).NotTo(HaveOccurred())

		Expect(fakeResultSource.GetResultCallCount()).To(Equal(1))
		url, token, component, metrics := fakeResultSource.GetResultArgsForCall(0)
		Expect(url).To(Equal("https://my.sonar.server"))
		Expect(token).To(Equal("token"))
		Expect(component).To(Equal("my:component"))
		Expect(metrics).To(Equal("ncloc,complexity,violations,coverage"))

		expectedDestination, err := ioutil.ReadDir(tmpDownloadDir)
		if err != nil {
			panic(1)
		}
		Expect(len(expectedDestination)).To(Equal(1))
		Expect(expectedDestination[0].Name()).To(Equal("result.json"))

		fullPath := filepath.Join(tmpDownloadDir, "result.json")
		content, err := ioutil.ReadFile(fullPath)
		if err != nil {
			panic(1)
		}
		Expect(content).To(Equal(fakeResponse))
	})

	It("writes the version to stdout", func() {
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

		err := run(stdin, stdout, tmpDownloadDir, fakeResultSource)

		Expect(err).NotTo(HaveOccurred())

		Expect(fakeResultSource.GetResultCallCount()).To(Equal(1))
		url, token, component, metrics := fakeResultSource.GetResultArgsForCall(0)
		Expect(url).To(Equal("https://my.sonar.server"))
		Expect(token).To(Equal("token"))
		Expect(component).To(Equal("my:component"))
		Expect(metrics).To(Equal("ncloc,complexity,violations,coverage"))


		expectedResponse := `{"version":{"ref":"61cebf"}}`
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

		err := run(stdin, stdout, tmpDownloadDir, fakeResultSource)

		Expect(fakeResultSource.GetResultCallCount()).To(Equal(1))
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

			err := run(stdin, stdout, tmpDownloadDir, fakeResultSource)

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

			err := run(stdin, stdout, tmpDownloadDir, fakeResultSource)

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

			err := run(stdin, stdout, tmpDownloadDir, fakeResultSource)

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

			err := run(stdin, stdout, tmpDownloadDir, fakeResultSource)

			Expect(err).To(HaveOccurred())
		})
	})
})
