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

package v5alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// EDIT THIS FILE!
// Created by "kubebuilder create resource" for you to implement the Pipeline resource schema definition
// as a go struct.
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PipelineSpec defines the desired state of Pipeline
type PipelineSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "kubebuilder generate" to regenerate code after modifying this file

	Resources []ConcourseResource `json:"resources"`
	Jobs      []ConcourseJob      `json:"jobs"`
}

type ConcourseJob struct {
	Name string                 `json:"name"`
	Plan runtime.RawExtension `json:"plan"`
}

type ConcourseResource struct {
	Name   string                 `json:"name"`
	Type   string                 `json:"type"`
	Source runtime.RawExtension `json:"source"`
}

// PipelineStatus defines the observed state of Pipeline
type PipelineStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "kubebuilder generate" to regenerate code after modifying this file

	PipelineSet bool   `json:"pipelineSet"`
	PipelineUrl string `json:"pipelineUrl"`
	Paused      bool   `json:"paused"`
	Public      bool   `json:"public"`

	ConcourseVersion       string `json:"concourseVersion"`
	ConcourseWorkerVersion string `json:"concourseWorkerVersion"`

	Builds []BuildStatus `json:"builds"`
}

// BuildStatus defines the observed Builds reflecting the history of the Pipeline
type BuildStatus struct {
	Url       string          `json:"url"`
	JobName   string          `json:"jobName"`
	Status    string          `json:"status"`
	StartTime metav1.Time     `json:"startTime"`
	EndTime   metav1.Time     `json:"endTime"`
	Duration  metav1.Duration `json:"duration"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Pipeline
// +k8s:openapi-gen=true
// +kubebuilder:resource:path=pipelines
// +kubebuilder:subresource:status
type Pipeline struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PipelineSpec   `json:"spec"`
	Status PipelineStatus `json:"status"`
}
