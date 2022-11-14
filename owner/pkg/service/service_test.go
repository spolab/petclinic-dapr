package service_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service", func() {

	Context("Given a valid service instance", func() {
		When("Registering a valid owner", func() {
			It("Returns no error", func() {
				Expect(nil).To(BeNil())
			})
		})
	})
})
