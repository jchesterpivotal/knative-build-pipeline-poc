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

package args

import (
	"time"

	"github.com/kubernetes-sigs/kubebuilder/pkg/inject/args"
	"k8s.io/client-go/rest"

	clientset "github.com/jchesterpivotal/knative-build-pipeline-poc/pkg/client/clientset/versioned"
	informer "github.com/jchesterpivotal/knative-build-pipeline-poc/pkg/client/informers/externalversions"
)

// InjectArgs are the arguments need to initialize controllers
type InjectArgs struct {
	args.InjectArgs

	Clientset *clientset.Clientset
	Informers informer.SharedInformerFactory
}

// CreateInjectArgs returns new controller args
func CreateInjectArgs(config *rest.Config) InjectArgs {
	cs := clientset.NewForConfigOrDie(config)
	return InjectArgs{
		InjectArgs: args.CreateInjectArgs(config),

		Clientset: cs,
		Informers: informer.NewSharedInformerFactory(cs, 2*time.Minute),
	}
}
