/*
Copyright 2026.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KubeHealthSpec defines the desired state of KubeHealth.
type KubeHealthSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of KubeHealth. Edit kubehealth_types.go to remove/update
	

	// Namespace to monitor
	Namespace string `json:"namespace"`

	// Restart threshold after which a pod is considered unhealthy
	RestartThreshold int32 `json:"restartThreshold"`
}

// KubeHealthStatus defines the observed state of KubeHealth.
type KubeHealthStatus struct {
	HealthyPods    int32       `json:"healthyPods,omitempty"`
	FailedPods     int32       `json:"failedPods,omitempty"`
	RestartingPods int32       `json:"restartingPods,omitempty"`
	ClusterStatus  string      `json:"clusterStatus,omitempty"`
	LastChecked    metav1.Time `json:"lastChecked,omitempty"`
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// KubeHealth is the Schema for the kubehealths API.
type KubeHealth struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KubeHealthSpec   `json:"spec,omitempty"`
	Status KubeHealthStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KubeHealthList contains a list of KubeHealth.
type KubeHealthList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KubeHealth `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KubeHealth{}, &KubeHealthList{})
}
