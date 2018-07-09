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

	Describe("The information given by 'kubectl get'", func() {
		It("Contains a URL for the pipeline", func() {
			kubectlCmd := exec.Command("kubectl", "describe", "pipeline", "pipeline.example.com")
			session, err := gexec.Start(kubectlCmd, GinkgoWriter, GinkgoWriter)
			Expect(err).ShouldNot(HaveOccurred())

			Eventually(session.Out).Should(gbytes.Say(`Spec:
  Jobs:
    Name:  successful-job
    Plan:
      Get:   git-repo
      File:  git-repo/tasks/successful/task.yml
      Task:  job-task
    Name:    failed-job
    Plan:
      Get:   git-repo
      File:  git-repo/tasks/failed/task.yml
      Task:  job-task
    Name:    flaky-job
    Plan:
      Get:   git-repo
      File:  git-repo/tasks/flaky/task.yml
      Task:  job-task
  Resources:
    Name:  git-repo
    Source:
      Uri:  https://github.com/concourse-courses/002-basics-of-reading-a-concourse-pipeline.git
    Type:   git
Status:
  Concourse API URL:         http://concourse-web.concourse.svc.cluster.local:8080
  Concourse Version:         3.14.1
  Concourse Worker Version:  2.1
  Paused:                    true
  Pipeline Set:              true
  Pipeline URL:              http://concourse-web.concourse.svc.cluster.local:8080/teams/main/pipelines/pipeline.example.com
  Public:                    false`))
		})
	})
})
