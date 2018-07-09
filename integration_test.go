/*
Copyright 2018 The Knative Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package integration_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"os/exec"
	"time"
)

var _ = Describe("Integration", func() {
	BeforeSuite(func() {
		kubectlCmd := exec.Command("kubectl", "delete", "pipeline", "pipeline.example.com")
		session, err := gexec.Start(kubectlCmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ShouldNot(HaveOccurred())

		// This is gross, but we need to give time for kubernetes to update its view of the world
		// before proceeding to the rest of the test suite.
		time.Sleep(time.Millisecond * 500)

		// It is useful to know whether we deleted something or not.
		fmt.Printf("BeforeSuite 'kubectl delete' response: %s %s\n\n", session.Out.Contents(), session.Err.Contents())
		Eventually(session).Should(gexec.Exit())
	})

	It("Can set a pipeline with kubectl", func() {
		kubectlCmd := exec.Command("kubectl", "apply", "-f", "config/example-pipeline.yaml")
		session, err := gexec.Start(kubectlCmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ShouldNot(HaveOccurred())

		Eventually(session.Out).Should(gbytes.Say("pipeline.concourse.concourse-ci.org/pipeline.example.com created"))
	})

	Describe("The information given by 'kubectl describe'", func() {
		var session *gexec.Session
		var err error

		BeforeEach(func() {
			kubectlCmd := exec.Command("kubectl", "describe", "pipeline", "pipeline.example.com")
			session, err = gexec.Start(kubectlCmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

		})

		Describe("In the Spec", func() {
			It("Contains the Pipeline's Jobs", func() {
				Eventually(session.Out).Should(gbytes.Say("Jobs:\n    Name:  successful-job"))
			})
			It("Contains the Pipeline's Resources", func() {
				Eventually(session.Out).Should(gbytes.Say("Resources:\n    Name:  git-repo"))
			})
		})
		Describe("In the Status", func() {
			It("Contains an Pipeline URL", func() {
				Eventually(session.Out).Should(gbytes.Say("Pipeline URL:              http://concourse-web.concourse.svc.cluster.local:8080/teams/main/pipelines/pipeline.example.com"))
			})
			It("Contains a 'Pipeline Set' value", func() {
				Eventually(session.Out).Should(gbytes.Say("Pipeline Set:              true"))
			})
			It("Shows if the Pipeline is paused", func() {
				Eventually(session.Out).Should(gbytes.Say("Paused:                    true"))
			})
			It("Shows if the Pipeline is public", func() {
				Eventually(session.Out).Should(gbytes.Say("Public:                    false"))
			})
			It("Contains an API URL", func() {
				Eventually(session.Out).Should(gbytes.Say("Concourse API URL:         http://concourse-web.concourse.svc.cluster.local:8080"))
			})
			It("Contains a Concourse server/API version", func() {
				Eventually(session.Out).Should(gbytes.Say("Concourse Version:         3.14.1"))
			})
			It("Contains a Concourse worker version", func() {
				Eventually(session.Out).Should(gbytes.Say("Concourse Worker Version:  2.1"))
			})
		})
	})
})
