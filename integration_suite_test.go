package integration_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestKnativeBuildPipelinePoc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "KnativeBuildPipelinePoc Suite")
}
