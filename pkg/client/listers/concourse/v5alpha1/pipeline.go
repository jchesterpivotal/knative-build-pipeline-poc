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

// Code generated by lister-gen. DO NOT EDIT.

package v5alpha1

import (
	v5alpha1 "github.com/topherbullock/knative-build-pipeline-poc/pkg/apis/concourse/v5alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// PipelineLister helps list Pipelines.
type PipelineLister interface {
	// List lists all Pipelines in the indexer.
	List(selector labels.Selector) (ret []*v5alpha1.Pipeline, err error)
	// Pipelines returns an object that can list and get Pipelines.
	Pipelines(namespace string) PipelineNamespaceLister
	PipelineListerExpansion
}

// pipelineLister implements the PipelineLister interface.
type pipelineLister struct {
	indexer cache.Indexer
}

// NewPipelineLister returns a new PipelineLister.
func NewPipelineLister(indexer cache.Indexer) PipelineLister {
	return &pipelineLister{indexer: indexer}
}

// List lists all Pipelines in the indexer.
func (s *pipelineLister) List(selector labels.Selector) (ret []*v5alpha1.Pipeline, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v5alpha1.Pipeline))
	})
	return ret, err
}

// Pipelines returns an object that can list and get Pipelines.
func (s *pipelineLister) Pipelines(namespace string) PipelineNamespaceLister {
	return pipelineNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// PipelineNamespaceLister helps list and get Pipelines.
type PipelineNamespaceLister interface {
	// List lists all Pipelines in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v5alpha1.Pipeline, err error)
	// Get retrieves the Pipeline from the indexer for a given namespace and name.
	Get(name string) (*v5alpha1.Pipeline, error)
	PipelineNamespaceListerExpansion
}

// pipelineNamespaceLister implements the PipelineNamespaceLister
// interface.
type pipelineNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Pipelines in the indexer for a given namespace.
func (s pipelineNamespaceLister) List(selector labels.Selector) (ret []*v5alpha1.Pipeline, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v5alpha1.Pipeline))
	})
	return ret, err
}

// Get retrieves the Pipeline from the indexer for a given namespace and name.
func (s pipelineNamespaceLister) Get(name string) (*v5alpha1.Pipeline, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v5alpha1.Resource("pipeline"), name)
	}
	return obj.(*v5alpha1.Pipeline), nil
}
