package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestUUPS(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UUPS Suite")
}
