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
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
// +genclient:nonNamespaced
// +kubebuilder:resource:path=tenantrolebindings,scope=Cluster
// TenantRoleBinding is the Schema for the tenantrolebindings API
type TenantRoleBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Subjects holds references to the objects the role applies to.
	// +optional
	Subjects []rbacv1.Subject `json:"subjects,omitempty" protobuf:"bytes,2,rep,name=subjects"`

	// Role subject has in Tenant:  admin, edit, view
	Role TenantRole `json:"roleRef" protobuf:"bytes,3,opt,name=roleRef"`

	// Tenant
	Tenant corev1.ObjectReference `json:"tenant"`
}

type TenantRole string

const (
	TenantRoleAdmin TenantRole = "admin"
	TenantRoleEdit  TenantRole = "edit"
	TenantRoleView  TenantRole = "view"
)

// +kubebuilder:object:root=true

// TenantRoleBindingList contains a list of TenantRoleBinding
type TenantRoleBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TenantRoleBinding `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TenantRoleBinding{}, &TenantRoleBindingList{})
}
