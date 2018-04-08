package shared_test

import (
	"errors"
	"github.com/concourse-sonarqube-notifier/assets/shared"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Types", func() {

	Describe("Source", func() {
		It("valid - returns true when all mandatory fields are filled", func() {
			src := shared.Source{
				Target: "Target",
				SonarToken: "Token",
				Metrics: "Metrics",
				Component: "Component",
			}
			Expect(src.Valid()).To(Equal(true))
		})

		It("valid - requires a target", func() {
			src := shared.Source{
				Target: "",
				SonarToken: "Token",
				Metrics: "Metrics",
				Component: "Component",
			}
			Expect(src.Valid()).To(Equal(false))
		})

		It("valid - requires a sonartoken", func() {
			src := shared.Source{
				Target: "target",
				SonarToken: "",
				Metrics: "Metrics",
				Component: "Component",
			}
			Expect(src.Valid()).To(Equal(false))
		})

		It("valid - requires metrics", func() {
			src := shared.Source{
				Target: "target",
				SonarToken: "Token",
				Metrics: "",
				Component: "Component",
			}
			Expect(src.Valid()).To(Equal(false))
		})

		It("valid - requires a component", func() {
			src := shared.Source{
				Target: "target",
				SonarToken: "Token",
				Metrics: "Metrics",
				Component: "",
			}
			Expect(src.Valid()).To(Equal(false))
		})

	})

	Describe("hasError", func() {
		It("returns true if an error is present", func() {
			result := shared.HasError(errors.New("BAD"))
			Expect(result).To(Equal(true))
		})

		It("returns false if no error is present", func() {
			result := shared.HasError(nil)
			Expect(result).To(Equal(false))
		})
	})

})
