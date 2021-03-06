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
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: "concourse.concourse-ci.org", Version: "v5alpha1"}

// Kind takes an unqualified kind and returns back a Group qualified GroupKind
func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

// Adds the list of known types to Scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Pipeline{},
		&PipelineList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type PipelineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Pipeline `json:"items"`
}

// CRD Generation
func getFloat(f float64) *float64 {
	return &f
}

func getInt(i int64) *int64 {
	return &i
}

var (
	// Define CRDs for resources
	PipelineCRD = v1beta1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: "pipelines.concourse.concourse-ci.org",
		},
		Spec: v1beta1.CustomResourceDefinitionSpec{
			Group:   "concourse.concourse-ci.org",
			Version: "v5alpha1",
			Names: v1beta1.CustomResourceDefinitionNames{
				Kind:   "Pipeline",
				Plural: "pipelines",
			},
			Scope: "Namespaced",
			Validation: &v1beta1.CustomResourceValidation{
				OpenAPIV3Schema: &v1beta1.JSONSchemaProps{
					Type: "object",
					Properties: map[string]v1beta1.JSONSchemaProps{
						"apiVersion": v1beta1.JSONSchemaProps{
							Type: "string",
						},
						"kind": v1beta1.JSONSchemaProps{
							Type: "string",
						},
						"metadata": v1beta1.JSONSchemaProps{
							Type: "object",
						},
						"spec": v1beta1.JSONSchemaProps{
							Type: "object",
							Properties: map[string]v1beta1.JSONSchemaProps{
								"jobs": v1beta1.JSONSchemaProps{
									Type: "array",
									Items: &v1beta1.JSONSchemaPropsOrArray{
										Schema: &v1beta1.JSONSchemaProps{
											Type: "object",
											Properties: map[string]v1beta1.JSONSchemaProps{
												"name": v1beta1.JSONSchemaProps{
													Type: "string",
												},
												"plan": v1beta1.JSONSchemaProps{
													Type:       "object",
													Properties: map[string]v1beta1.JSONSchemaProps{},
												},
												"public": v1beta1.JSONSchemaProps{
													Type: "boolean",
												},
											},
											Required: []string{
												"name",
												"public",
												"plan",
											}},
									},
								},
								"resources": v1beta1.JSONSchemaProps{
									Type: "array",
									Items: &v1beta1.JSONSchemaPropsOrArray{
										Schema: &v1beta1.JSONSchemaProps{
											Type: "object",
											Properties: map[string]v1beta1.JSONSchemaProps{
												"name": v1beta1.JSONSchemaProps{
													Type: "string",
												},
												"source": v1beta1.JSONSchemaProps{
													Type:       "object",
													Properties: map[string]v1beta1.JSONSchemaProps{},
												},
												"type": v1beta1.JSONSchemaProps{
													Type: "string",
												},
											},
											Required: []string{
												"name",
												"type",
												"source",
											}},
									},
								},
							},
							Required: []string{
								"resources",
								"jobs",
							}},
						"status": v1beta1.JSONSchemaProps{
							Type: "object",
							Properties: map[string]v1beta1.JSONSchemaProps{
								"builds": v1beta1.JSONSchemaProps{
									Type: "array",
									Items: &v1beta1.JSONSchemaPropsOrArray{
										Schema: &v1beta1.JSONSchemaProps{
											Type: "object",
											Properties: map[string]v1beta1.JSONSchemaProps{
												"duration": v1beta1.JSONSchemaProps{
													Type:       "object",
													Properties: map[string]v1beta1.JSONSchemaProps{},
												},
												"endTime": v1beta1.JSONSchemaProps{
													Type:   "string",
													Format: "date-time",
												},
												"jobName": v1beta1.JSONSchemaProps{
													Type: "string",
												},
												"startTime": v1beta1.JSONSchemaProps{
													Type:   "string",
													Format: "date-time",
												},
												"status": v1beta1.JSONSchemaProps{
													Type: "string",
												},
												"url": v1beta1.JSONSchemaProps{
													Type: "string",
												},
											},
											Required: []string{
												"url",
												"jobName",
												"status",
												"startTime",
												"endTime",
												"duration",
											}},
									},
								},
								"concourseVersion": v1beta1.JSONSchemaProps{
									Type: "string",
								},
								"concourseWorkerVersion": v1beta1.JSONSchemaProps{
									Type: "string",
								},
								"paused": v1beta1.JSONSchemaProps{
									Type: "boolean",
								},
								"pipelineSet": v1beta1.JSONSchemaProps{
									Type: "boolean",
								},
								"pipelineUrl": v1beta1.JSONSchemaProps{
									Type: "string",
								},
								"public": v1beta1.JSONSchemaProps{
									Type: "boolean",
								},
							},
							Required: []string{
								"pipelineSet",
								"pipelineUrl",
								"paused",
								"public",
								"concourseVersion",
								"concourseWorkerVersion",
								"builds",
							}},
					},
					Required: []string{
						"spec",
						"status",
					}},
			},
			Subresources: &v1beta1.CustomResourceSubresources{
				Status: &v1beta1.CustomResourceSubresourceStatus{},
			},
		},
	}
)
