/*


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

package v1alpha4

import (
	etcdbp "github.com/mrajashree/etcdadm-bootstrap-provider/api/v1alpha4"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	EtcdCertsGeneratedCondition string = "EtcdCertsGenerated"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// EtcdadmClusterSpec defines the desired state of EtcdadmCluster
type EtcdadmClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Replicas *int32 `json:"replicas,omitempty"`

	// +optional
	Version string `json:"version,omitempty"`

	// InfrastructureTemplate is a required reference to a custom resource
	// offered by an infrastructure provider.
	InfrastructureTemplate corev1.ObjectReference `json:"infrastructureTemplate"`

	// +optional
	EtcdadmConfigSpec etcdbp.EtcdadmConfigSpec `json:"etcdadmConfigSpec"`
}

// EtcdadmClusterStatus defines the observed state of EtcdadmCluster
type EtcdadmClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Total number of non-terminated machines targeted by this etcd cluster
	// (their labels match the selector).
	// +optional
	ReadyReplicas int32 `json:"replicas,omitempty"`

	// +optional
	InitMachineAddress string `json:"initMachineAddress"`

	// +optional
	Initialized bool `json:"initialized"`

	// +optional
	CreationComplete bool `json:"creationComplete"`

	// +optional
	Endpoint string `json:"endpoint"`

	// Selector is the label selector in string format to avoid introspection
	// by clients, and is used to provide the CRD-based integration for the
	// scale subresource and additional integrations for things like kubectl
	// describe.. The string will be in the same format as the query-param syntax.
	// More info about label selectors: http://kubernetes.io/docs/user-guide/labels#label-selectors
	// +optional
	Selector string `json:"selector,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// EtcdadmCluster is the Schema for the etcdclusters API
type EtcdadmCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EtcdadmClusterSpec   `json:"spec,omitempty"`
	Status EtcdadmClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// EtcdadmClusterList contains a list of EtcdadmCluster
type EtcdadmClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EtcdadmCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EtcdadmCluster{}, &EtcdadmClusterList{})
}
