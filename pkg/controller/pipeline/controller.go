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

package pipeline

import (
	"log"

	"github.com/kubernetes-sigs/kubebuilder/pkg/controller"
	"github.com/kubernetes-sigs/kubebuilder/pkg/controller/types"
	"k8s.io/client-go/tools/record"

	concoursev5alpha1 "github.com/topherbullock/knative-build-pipeline-poc/pkg/apis/concourse/v5alpha1"
	concoursev5alpha1client "github.com/topherbullock/knative-build-pipeline-poc/pkg/client/clientset/versioned/typed/concourse/v5alpha1"
	concoursev5alpha1informer "github.com/topherbullock/knative-build-pipeline-poc/pkg/client/informers/externalversions/concourse/v5alpha1"
	concoursev5alpha1lister "github.com/topherbullock/knative-build-pipeline-poc/pkg/client/listers/concourse/v5alpha1"

	"github.com/topherbullock/knative-build-pipeline-poc/pkg/inject/args"

	"github.com/concourse/go-concourse/concourse"

	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

var concourseClient concourse.Client

// EDIT THIS FILE
// This files was created by "kubebuilder create resource" for you to edit.
// Controller implementation logic for Pipeline resources goes here.

func (bc *PipelineController) Reconcile(k types.ReconcileKey) error {
	pipelineInK8s, err := bc.pipelineclient.Pipelines(k.Namespace).Get(k.Name, v1.GetOptions{IncludeUninitialized: false})
	if err != nil {
		log.Printf("Failed to get pipeline resource for key '%s': %s", k, err.Error())
		return err
	}

	concourseClient, err := ConcourseClient("http://concourse-web.concourse.svc.cluster.local:8080")
	if err != nil {
		log.Printf("Failed to set up a Concourse client for key '%s': %s", k, err.Error())
		return err
	}

	team := concourseClient.Team("main")

	_, _, currentVersionNumber, _, err := team.PipelineConfig(k.Name)
	if err != nil {
		log.Printf("There was an error retrieving the existing config for '%s': %+v", k.Name, err)
		return err
	}

	pipelineSpec, err := json.Marshal(pipelineInK8s.Spec)
	pipelineSpecYaml, err := yaml.JSONToYAML(pipelineSpec)
	_, _, warnings, err := team.CreateOrUpdatePipelineConfig(
		k.Name,
		currentVersionNumber,
		pipelineSpecYaml,
	)

	if len(warnings) > 0 {
		log.Printf("After setting the %s pipeline in Concourse, there were warnings: %+v", k.Name, warnings)
	}

	info, err := concourseClient.GetInfo()
	if err != nil {
		log.Printf("Failed to get Concourse server information for key '%s': %s", k, err.Error())
		return err
	}

	pipelineInConcourse, foundPipeline, err := team.Pipeline(k.Name)
	if err != nil {
		log.Printf("Failed to get Concourse pipeline information for key '%s': %s", k, err.Error())
		return err
	}

	pipelineUrl := fmt.Sprintf("%s/teams/%s/pipelines/%s", concourseClient.URL(), team.Name(), k.Name)

	pipelineInK8s.Status = concoursev5alpha1.PipelineStatus{
		PipelineUrl: pipelineUrl,
		Paused:      pipelineInConcourse.Paused,
		Public:      pipelineInConcourse.Public,

		ConcourseAPIUrl:        concourseClient.URL(),
		ConcourseVersion:       info.Version,
		ConcourseWorkerVersion: info.WorkerVersion,
	}

	if foundPipeline {
		pipelineInK8s.Status.PipelineSet = true
	} else {
		pipelineInK8s.Status.PipelineSet = false
	}

	_, err = bc.pipelineclient.Pipelines(k.Namespace).Update(pipelineInK8s)
	if err != nil {
		log.Printf("Failed to update pipeline status for key '%s': %s", k, err.Error())
		return err
	}

	return nil
}

// +kubebuilder:controller:group=concourse,version=v5alpha1,kind=Pipeline,resource=pipelines
// +kubebuilder:rbac:groups=concourse,resources=pipelines,verbs=get;watch;list;create;update
type PipelineController struct {
	// INSERT ADDITIONAL FIELDS HERE
	pipelineLister concoursev5alpha1lister.PipelineLister
	pipelineclient concoursev5alpha1client.ConcourseV5alpha1Interface
	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	pipelinerecorder record.EventRecorder
}

// ProvideController provides a controller that will be run at startup.  Kubebuilder will use codegeneration
// to automatically register this controller in the inject package
func ProvideController(arguments args.InjectArgs) (*controller.GenericController, error) {
	// INSERT INITIALIZATIONS FOR ADDITIONAL FIELDS HERE
	bc := &PipelineController{
		pipelineLister: arguments.ControllerManager.GetInformerProvider(&concoursev5alpha1.Pipeline{}).(concoursev5alpha1informer.PipelineInformer).Lister(),

		pipelineclient:   arguments.Clientset.ConcourseV5alpha1(),
		pipelinerecorder: arguments.CreateRecorder("PipelineController"),
	}

	// Create a new controller that will call PipelineController.Reconcile on changes to Pipelines
	gc := &controller.GenericController{
		Name:             "PipelineController",
		Reconcile:        bc.Reconcile,
		InformerRegistry: arguments.ControllerManager,
	}
	if err := gc.Watch(&concoursev5alpha1.Pipeline{}); err != nil {
		return gc, err
	}

	// IMPORTANT:
	// To watch additional resource types - such as those created by your controller - add gc.Watch* function calls here
	// Watch function calls will transform each object event into a Pipeline Key to be reconciled by the controller.
	//
	// **********
	// For any new Watched types, you MUST add the appropriate // +kubebuilder:informer and // +kubebuilder:rbac
	// annotations to the PipelineController and run "kubebuilder generate.
	// This will generate the code to start the informers and create the RBAC rules needed for running in a cluster.
	// See:
	// https://godoc.org/github.com/kubernetes-sigs/kubebuilder/pkg/gen/controller#example-package
	// **********

	return gc, nil
}
