package main

import (
	"bytes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Out", func() {
	It("outputs an empty JSON array so that it satisfies the resource interface", func() {
		stdout := new(bytes.Buffer)

		status := run(stdout)

		Expect(status).To(Equal(0))
		Expect(stdout.String()).To(Equal(`[]`))
	})
})
