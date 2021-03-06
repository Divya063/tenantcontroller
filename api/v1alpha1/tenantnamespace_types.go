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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TenantNamespaceSpec defines the desired state of TenantNamespace
type TenantNamespaceSpec struct {
	Tenant string `json:"tenant,omitempty"`
}

// TenantNamespaceStatus defines the observed state of TenantNamespace
type TenantNamespaceStatus struct {
	corev1.NamespaceStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +genclient:nonNamespaced
// +kubebuilder:resource:path=tenantnamespaces,scope=Cluster
// TenantNamespace is the Schema for the tenantnamespaces API
type TenantNamespace struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TenantNamespaceSpec   `json:"spec,omitempty"`
	Status TenantNamespaceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TenantNamespaceList contains a list of TenantNamespace
type TenantNamespaceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TenantNamespace `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TenantNamespace{}, &TenantNamespaceList{})
}
