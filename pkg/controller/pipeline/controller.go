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

	concoursev5alpha1 "github.com/jchesterpivotal/knative-build-pipeline-poc/pkg/apis/concourse/v5alpha1"
	concoursev5alpha1client "github.com/jchesterpivotal/knative-build-pipeline-poc/pkg/client/clientset/versioned/typed/concourse/v5alpha1"
	concoursev5alpha1informer "github.com/jchesterpivotal/knative-build-pipeline-poc/pkg/client/informers/externalversions/concourse/v5alpha1"
	concoursev5alpha1lister "github.com/jchesterpivotal/knative-build-pipeline-poc/pkg/client/listers/concourse/v5alpha1"

	"github.com/jchesterpivotal/knative-build-pipeline-poc/pkg/inject/args"

	"github.com/concourse/go-concourse/concourse"

	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
	"os"
	"errors"
)

// EDIT THIS FILE
// This files was created by "kubebuilder create resource" for you to edit.
// Controller implementation logic for Pipeline resources goes here.

func (bc *PipelineController) Reconcile(k types.ReconcileKey) error {
	pipelineInK8s, err := bc.pipelineclient.Pipelines(k.Namespace).Get(k.Name, metav1.GetOptions{IncludeUninitialized: false})
	if err != nil {
		log.Printf("Failed to get pipeline resource for key '%s': %s", k, err.Error())
		return err
	}
	pipelineInK8s = pipelineInK8s.DeepCopy()

	_, _, _, pipelineAlreadyExists, err := bc.concourseTeam.PipelineConfig(k.Name)
	if err != nil {
		log.Printf("There was an error checking for an existing config for '%s': %+v", k.Name, err.Error())
		return err
	}

	if pipelineAlreadyExists {
		return bc.updateJobsInStatus(pipelineInK8s, k)
	}

	return bc.createNewPipelineInConcourse(pipelineInK8s, k)
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

	concourseClient concourse.Client
	concourseTeam   concourse.Team
}

// ProvideController provides a controller that will be run at startup.  Kubebuilder will use codegeneration
// to automatically register this controller in the inject package
func ProvideController(arguments args.InjectArgs) (*controller.GenericController, error) {
	bc := &PipelineController{
		pipelineLister: arguments.ControllerManager.GetInformerProvider(&concoursev5alpha1.Pipeline{}).(concoursev5alpha1informer.PipelineInformer).Lister(),

		pipelineclient:   arguments.Clientset.ConcourseV5alpha1(),
		pipelinerecorder: arguments.CreateRecorder("PipelineController"),
	}

	concourseUrl, found := os.LookupEnv("CONCOURSE_URL")
	if !found {
		errMsg := fmt.Sprintf("CONCOURSE_URL was not set")
		log.Printf(errMsg)
		return nil, errors.New(errMsg)
	}

	username, found := os.LookupEnv("CONCOURSE_USERNAME")
	if !found {
		errMsg := fmt.Sprintf("CONCOURSE_USERNAME was not set")
		log.Printf(errMsg)
		return nil, errors.New(errMsg)
	}

	password, found := os.LookupEnv("CONCOURSE_PASSWORD")
	if !found {
		errMsg := fmt.Sprintf("CONCOURSE_PASSWORD was not set")
		log.Printf(errMsg)
		return nil, errors.New(errMsg)
	}

	client, err := ConcourseClient(concourseUrl, username, password)
	if err != nil {
		log.Printf("Failed to set up a Concourse client for server '%s': %s", concourseUrl, err.Error())
		return nil, err
	}

	bc.concourseClient = client
	bc.concourseTeam = bc.concourseClient.Team("main")

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

func (bc *PipelineController) updateJobsInStatus(pipelineInK8s *concoursev5alpha1.Pipeline, k types.ReconcileKey) error {
	jobs, err := bc.concourseTeam.ListJobs(k.Name)
	if err != nil {
		log.Printf("failed to get jobs for '%s': %s", k.Name, err)
		return err
	}

	builds := make([]concoursev5alpha1.BuildStatus, 0)
	for _, j := range jobs {
		// TODO: do something smarter with pagination
		jobsInBuild, _, found, err := bc.concourseTeam.JobBuilds(k.Name, j.Name, concourse.Page{})
		if err != nil {
			log.Printf("failed to get builds for '%s/%s': %s", k.Name, j.Name, err.Error())
			return err
		}

		if !found {
			break
		}

		for _, b := range jobsInBuild {
			buildUrl := fmt.Sprintf("%s/teams/%s/pipelines/%s/jobs/%s/builds/%s", bc.concourseClient.URL(), bc.concourseTeam.Name(), k.Name, j.Name, b.Name)

			var start, end metav1.Time
			var duration metav1.Duration
			start = metav1.Unix(b.StartTime, 0)
			if b.EndTime == 0 {
				duration = metav1.Duration{Duration: time.Since(start.Time)}
			} else {
				end = metav1.Unix(b.EndTime, 0)
				duration = metav1.Duration{Duration: end.Time.Sub(start.Time)}
			}

			statBuild := concoursev5alpha1.BuildStatus{
				Url:       buildUrl,
				JobName:   b.JobName,
				Status:    b.Status,
				StartTime: start,
				EndTime:   end,
				Duration:  duration,
			}
			builds = append(builds, statBuild)
		}
	}

	pipelineInK8s.Status.Builds = builds

	_, err = bc.pipelineclient.Pipelines(k.Namespace).Update(pipelineInK8s)
	if err != nil {
		log.Printf("Failed to update pipeline status for key '%s': %s", k, err.Error())
		return err
	}

	return nil
}

func (bc *PipelineController) createNewPipelineInConcourse(pipelineInK8s *concoursev5alpha1.Pipeline, k types.ReconcileKey) error {
	pipelineSpec, err := json.Marshal(pipelineInK8s.Spec)
	pipelineSpecYaml, err := yaml.JSONToYAML(pipelineSpec)
	_, _, warnings, err := bc.concourseTeam.CreateOrUpdatePipelineConfig(
		k.Name,
		"",
		pipelineSpecYaml,
	)

	if len(warnings) > 0 {
		log.Printf("After setting the %s pipeline in Concourse, there were warnings: %+v", k.Name, warnings)
	}

	info, err := bc.concourseClient.GetInfo()
	if err != nil {
		log.Printf("Failed to get Concourse server information for key '%s': %s", k, err.Error())
		return err
	}

	pipelineInConcourse, foundPipeline, err := bc.concourseTeam.Pipeline(k.Name)
	if err != nil {
		log.Printf("Failed to get Concourse pipeline information for key '%s': %s", k, err.Error())
		return err
	}

	pipelineUrl := fmt.Sprintf("%s/teams/%s/pipelines/%s", bc.concourseClient.URL(), bc.concourseTeam.Name(), k.Name)

	pipelineInK8s.Status = concoursev5alpha1.PipelineStatus{
		PipelineUrl: pipelineUrl,
		Paused:      pipelineInConcourse.Paused,
		Public:      pipelineInConcourse.Public,

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
