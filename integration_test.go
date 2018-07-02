package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os/exec"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Integration", func() {
	BeforeEach(func() {
		kubectlCmd := exec.Command("kubectl", "delete", "pipeline", "pipeline.example.com")
		session, err := gexec.Start(kubectlCmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ShouldNot(HaveOccurred())
		Eventually(session).Should(gexec.Exit())
	})

	It("Can set a pipeline", func() {
		kubectlCmd := exec.Command("kubectl", "apply", "-f", "config/example-pipeline.yaml")
		session, err := gexec.Start(kubectlCmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ShouldNot(HaveOccurred())

		Eventually(session.Out).Should(gbytes.Say("pipeline.concourse-ci.org/pipeline.example.com created"))
	})
})
