package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os/exec"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/gbytes"
	"time"
	"fmt"
)

var _ = Describe("Integration", func() {
	BeforeSuite(func() {
		kubectlCmd := exec.Command("kubectl", "delete", "pipeline", "pipeline.example.com")
		session, err := gexec.Start(kubectlCmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ShouldNot(HaveOccurred())

		// This is gross, but we need to give time for kubernetes to update its view of the world
		// before proceeding to the rest of the test suite.
		time.Sleep(time.Millisecond*500)

		// It is useful to know whether we deleted something or not.
		fmt.Printf("BeforeSuite 'kubectl delete' response: %s %s\n\n", session.Out.Contents(), session.Err.Contents())
		Eventually(session).Should(gexec.Exit())
	})

	It("Can set a pipeline with kubectl", func() {
		kubectlCmd := exec.Command("kubectl", "apply", "-f", "config/example-pipeline.yaml")
		session, err := gexec.Start(kubectlCmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ShouldNot(HaveOccurred())

		Eventually(session.Out).Should(gbytes.Say("pipeline.concourse-ci.org/pipeline.example.com created"))
	})

	It("Can retrieve a pipeline with kubectl", func() {
		kubectlCmd := exec.Command("kubectl", "get", "pipelines")
		session, err := gexec.Start(kubectlCmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ShouldNot(HaveOccurred())

		Eventually(session.Out).Should(gbytes.Say(`NAME                   CREATED AT
pipeline.example.com   [\d+]s`))
	})
})
