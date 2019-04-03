/*
Copyright 2019 The Knative Authors

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

package v1beta1

import (
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/knative/pkg/apis"
	"github.com/knative/pkg/apis/duck"
)

// Conditions is a simple wrapper around apis.Conditions to implement duck.Implementable.
type Conditions apis.Conditions

// Conditions is an Implementable "duck type".
var _ duck.Implementable = (*Conditions)(nil)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KResource is a skeleton type wrapping Conditions in the manner we expect
// resource writers defining compatible resources to embed it.  We will
// typically use this type to deserialize Conditions ObjectReferences and
// access the Conditions data.  This is not a real resource.
type KResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Status Status `json:"status"`
}

// Status shows how we expect folks to embed Conditions in
// their Status field.
// WARNING: Adding fields to this struct will add them to all Knative resources.
type Status struct {
	// ObservedGeneration is the 'Generation' of the Service that
	// was last processed by the controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Conditions the latest available observations of a resource's current state.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions Conditions `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

var _ apis.ConditionsAccessor = (*Status)(nil)

// GetConditions implements apis.ConditionsAccessor
func (s *Status) GetConditions() apis.Conditions {
	return apis.Conditions(s.Conditions)
}

// SetConditions implements apis.ConditionsAccessor
func (s *Status) SetConditions(c apis.Conditions) {
	s.Conditions = Conditions(c)
}

// In order for Conditions to be Implementable, KResource must be Populatable.
var _ duck.Populatable = (*KResource)(nil)

// Ensure KResource satisfies apis.Listable
var _ apis.Listable = (*KResource)(nil)

// GetFullType implements duck.Implementable
func (*Conditions) GetFullType() duck.Populatable {
	return &KResource{}
}

// GetCondition fetches the condition of the specified type.
func (s *Status) GetCondition(t apis.ConditionType) *apis.Condition {
	for _, cond := range s.Conditions {
		if cond.Type == t {
			return &cond
		}
	}
	return nil
}

// Populate implements duck.Populatable
func (t *KResource) Populate() {
	t.Status.ObservedGeneration = 42
	t.Status.Conditions = Conditions{{
		// Populate ALL fields
		Type:               "Birthday",
		Status:             corev1.ConditionTrue,
		LastTransitionTime: apis.VolatileTime{Inner: metav1.NewTime(time.Date(1984, 02, 28, 18, 52, 00, 00, time.UTC))},
		Reason:             "Celebrate",
		Message:            "n3wScott, find your party hat :tada:",
	}}
}

// GetListType implements apis.Listable
func (r *KResource) GetListType() runtime.Object {
	return &KResourceList{}
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KResourceList is a list of KResource resources
type KResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []KResource `json:"items"`
}