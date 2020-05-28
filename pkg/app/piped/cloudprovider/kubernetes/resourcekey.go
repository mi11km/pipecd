// Copyright 2020 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kubernetes

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var builtInApiVersions = map[string]struct{}{
	"admissionregistration.k8s.io/v1":      {},
	"admissionregistration.k8s.io/v1beta1": {},
	"apiextensions.k8s.io/v1":              {},
	"apiextensions.k8s.io/v1beta1":         {},
	"apiregistration.k8s.io/v1":            {},
	"apiregistration.k8s.io/v1beta1":       {},
	"apps/v1":                              {},
	"authentication.k8s.io/v1":             {},
	"authentication.k8s.io/v1beta1":        {},
	"authorization.k8s.io/v1":              {},
	"authorization.k8s.io/v1beta1":         {},
	"autoscaling/v1":                       {},
	"autoscaling/v2beta1":                  {},
	"autoscaling/v2beta2":                  {},
	"batch/v1":                             {},
	"batch/v1beta1":                        {},
	"certificates.k8s.io/v1beta1":          {},
	"coordination.k8s.io/v1":               {},
	"coordination.k8s.io/v1beta1":          {},
	"extensions/v1beta1":                   {},
	"internal.autoscaling.k8s.io/v1alpha1": {},
	"metrics.k8s.io/v1beta1":               {},
	"networking.k8s.io/v1":                 {},
	"networking.k8s.io/v1beta1":            {},
	"node.k8s.io/v1beta1":                  {},
	"policy/v1beta1":                       {},
	"rbac.authorization.k8s.io/v1":         {},
	"rbac.authorization.k8s.io/v1beta1":    {},
	"scheduling.k8s.io/v1":                 {},
	"scheduling.k8s.io/v1beta1":            {},
	"storage.k8s.io/v1":                    {},
	"storage.k8s.io/v1beta1":               {},
	"v1":                                   {},
}

type ResourceKey struct {
	APIVersion string
	Kind       string
	Namespace  string
	Name       string
}

func (k ResourceKey) String() string {
	return fmt.Sprintf("%s:%s:%s:%s", k.APIVersion, k.Kind, k.Namespace, k.Name)
}

func (k ResourceKey) IsZero() bool {
	return k.APIVersion == "" &&
		k.Kind == "" &&
		k.Namespace == "" &&
		k.Name == ""
}

func (k ResourceKey) IsKubernetesBuiltInResource() bool {
	_, ok := builtInApiVersions[k.APIVersion]
	return ok
}

func (k ResourceKey) IsDeployment() bool {
	if k.Kind != "Deployment" {
		return false
	}
	if !k.IsKubernetesBuiltInResource() {
		return false
	}
	return true
}

func (k ResourceKey) IsConfigMap() bool {
	if k.Kind != "ConfigMap" {
		return false
	}
	if !k.IsKubernetesBuiltInResource() {
		return false
	}
	return true
}

func (k ResourceKey) IsSecret() bool {
	if k.Kind != "Secret" {
		return false
	}
	if !k.IsKubernetesBuiltInResource() {
		return false
	}
	return true
}

func MakeResourceKey(obj *unstructured.Unstructured) ResourceKey {
	return ResourceKey{
		APIVersion: obj.GetAPIVersion(),
		Kind:       obj.GetKind(),
		Namespace:  obj.GetNamespace(),
		Name:       obj.GetName(),
	}
}

func DecodeResourceKey(key string) (ResourceKey, error) {
	parts := strings.Split(key, ":")
	if len(parts) != 4 {
		return ResourceKey{}, fmt.Errorf("malformed key")
	}
	return ResourceKey{
		APIVersion: parts[0],
		Kind:       parts[1],
		Namespace:  parts[2],
		Name:       parts[3],
	}, nil
}
