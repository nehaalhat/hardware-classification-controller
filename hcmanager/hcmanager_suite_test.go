package hcmanager_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHcmanager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Hcmanager Suite")
}
