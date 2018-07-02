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

package inject

import (
	"github.com/kubernetes-sigs/kubebuilder/pkg/inject/run"
	concoursev5alpha1 "github.com/topherbullock/knative-build-pipeline-poc/pkg/apis/concourse/v5alpha1"
	rscheme "github.com/topherbullock/knative-build-pipeline-poc/pkg/client/clientset/versioned/scheme"
	"github.com/topherbullock/knative-build-pipeline-poc/pkg/controller/pipeline"
	"github.com/topherbullock/knative-build-pipeline-poc/pkg/inject/args"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
)

func init() {
	rscheme.AddToScheme(scheme.Scheme)

	// Inject Informers
	Inject = append(Inject, func(arguments args.InjectArgs) error {
		Injector.ControllerManager = arguments.ControllerManager

		if err := arguments.ControllerManager.AddInformerProvider(&concoursev5alpha1.Pipeline{}, arguments.Informers.Concourse().V5alpha1().Pipelines()); err != nil {
			return err
		}

		// Add Kubernetes informers

		if c, err := pipeline.ProvideController(arguments); err != nil {
			return err
		} else {
			arguments.ControllerManager.AddController(c)
		}
		return nil
	})

	// Inject CRDs
	Injector.CRDs = append(Injector.CRDs, &concoursev5alpha1.PipelineCRD)
	// Inject PolicyRules
	Injector.PolicyRules = append(Injector.PolicyRules, rbacv1.PolicyRule{
		APIGroups: []string{"concourse.concourse-ci.org"},
		Resources: []string{"*"},
		Verbs:     []string{"*"},
	})
	Injector.PolicyRules = append(Injector.PolicyRules, rbacv1.PolicyRule{
		APIGroups: []string{
			"concourse",
		},
		Resources: []string{
			"pipelines",
		},
		Verbs: []string{
			"create", "get", "list", "update", "watch",
		},
	})
	// Inject GroupVersions
	Injector.GroupVersions = append(Injector.GroupVersions, schema.GroupVersion{
		Group:   "concourse.concourse-ci.org",
		Version: "v5alpha1",
	})
	Injector.RunFns = append(Injector.RunFns, func(arguments run.RunArguments) error {
		Injector.ControllerManager.RunInformersAndControllers(arguments)
		return nil
	})
}
