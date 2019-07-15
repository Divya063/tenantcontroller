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

package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	tenantv1alpha1 "github.com/runyontr/tenantcontroller/api/v1alpha1"
)

// TenantNamespaceReconciler reconciles a TenantNamespace object
type TenantNamespaceReconciler struct {
	client.Client
	Log    logr.Logger
	scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=tenant.multitenant.k8s.io,resources=tenantnamespaces,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.k8s.io,resources=networkpolicies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=tenant.multitenant.k8s.io,resources=tenantnamespaces/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=*,resources=namespaces,verbs=get;list;watch;create;update;patch;delete

// Reconcile here
func (r *TenantNamespaceReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("tenantnamespace", req.NamespacedName)

	// your logic here
	log.Info(fmt.Sprintf("Recieved reconcile request for tenantnamespace %v", req.NamespacedName))

	tns := tenantv1alpha1.TenantNamespace{}
	err := r.Get(ctx, req.NamespacedName, &tns)
	if err != nil {

		if errors.IsNotFound(err) {
			log.Info(fmt.Sprintf("TenantNamespace %v deleted", req.NamespacedName))
			return ctrl.Result{}, nil
		}
		log.Error(err, fmt.Sprintf("Error getting tenantnamespace %v: %v", req.NamespacedName, err))
		return ctrl.Result{}, err
	}

	tenant := tenantv1alpha1.Tenant{}
	err = r.Get(ctx, types.NamespacedName{Name: tns.Spec.Tenant}, &tenant)
	if err != nil {
		log.Error(err, fmt.Sprintf("Error getting tenant %v referenced in tenantnamespace %v", tns.Spec.Tenant, req.NamespacedName))
		return ctrl.Result{}, err
	}

	//Now ensure the namespace is created:
	ns := corev1.Namespace{}
	//The name of the namespace matches the name of the tns
	err = r.Get(ctx, types.NamespacedName{Name: tns.Name}, &ns)
	if err == nil {
		//already present,
		//TODO (@runyontr)
		// ensure labels to ns
		// ensure finalizer logic
		log.Info(fmt.Sprintf("Namespace %v already present.  Good work team", tns.Name), "namespace", tns.Name)
		return ctrl.Result{}, nil
	}
	labels := map[string]string{"tenant": tns.Spec.Tenant}
	//create a namespace
	ns.ObjectMeta = metav1.ObjectMeta{
		Name:   tns.Name,
		Labels: labels,
		// Finalizers: ,
		// Annotations: map[string]string{"multitenant.k8s.io/created-at",
	}
	// ns.Spec.Finalizers = []corev1.FinalizerName{corev1.FinalizerName(fmt.Sprintf("tenant.multitenant.k8s.io/%v", tenant.Name))}
	err = ctrl.SetControllerReference(&tns, &ns, r.scheme)
	if err != nil {
		log.Error(err, "Error setting owner ref of namespace", "namespace", tns.Name)
		return ctrl.Result{}, err
	}
	err = r.Create(ctx, &ns)

	if err != nil {
		return ctrl.Result{}, err
	}

	/*
		Network Policies  See docs/NetworkPolicy.md for more information
	*/
	// Make Namespace specific items:
	namespaceSelector := &metav1.LabelSelector{}
	for k, v := range labels {
		metav1.AddLabelToSelector(namespaceSelector, k, v)
	}

	netpol := networkingv1.NetworkPolicy{}

	netpol = networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "tenant-network-policy",
			Namespace: tns.Name,
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{},
			Ingress: []networkingv1.NetworkPolicyIngressRule{
				{
					From: []networkingv1.NetworkPolicyPeer{
						{
							PodSelector:       &metav1.LabelSelector{},
							NamespaceSelector: namespaceSelector,
						},
					},
				},
			},
			Egress: []networkingv1.NetworkPolicyEgressRule{
				{
					To: []networkingv1.NetworkPolicyPeer{
						{
							PodSelector:       &metav1.LabelSelector{},
							NamespaceSelector: namespaceSelector,
						},
					},
				},
			},
			PolicyTypes: []networkingv1.PolicyType{
				networkingv1.PolicyTypeEgress,
				networkingv1.PolicyTypeIngress,
			},
		},
	}
	err = r.Client.Create(ctx, &netpol)
	if err != nil {
		log.Error(err, "Could not create network policy in namespace"+tns.Name)
	}

	// Talk to namespaces with label multitenant.k8s.io=public
	namespaceSelector = &metav1.LabelSelector{}
	metav1.AddLabelToSelector(namespaceSelector, "multitenant.k8s.io", "public")
	netpol = networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "multitenant-public",
			Namespace: tns.Name,
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{},
			Ingress: []networkingv1.NetworkPolicyIngressRule{
				{
					From: []networkingv1.NetworkPolicyPeer{
						{
							PodSelector:       &metav1.LabelSelector{},
							NamespaceSelector: namespaceSelector,
						},
					},
				},
			},
			Egress: []networkingv1.NetworkPolicyEgressRule{
				{
					To: []networkingv1.NetworkPolicyPeer{
						{
							PodSelector:       &metav1.LabelSelector{},
							NamespaceSelector: namespaceSelector,
						},
					},
				},
			},
			PolicyTypes: []networkingv1.PolicyType{
				networkingv1.PolicyTypeEgress,
				networkingv1.PolicyTypeIngress,
			},
		},
	}
	err = r.Client.Create(ctx, &netpol)
	if err != nil {
		log.Error(err, "Could not create network policy in namespace"+tns.Name)
	}

	return ctrl.Result{}, nil
}

func (r *TenantNamespaceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.scheme = mgr.GetScheme()
	corev1.AddToScheme(r.scheme)
	networkingv1.AddToScheme(r.scheme)
	return ctrl.NewControllerManagedBy(mgr).
		For(&tenantv1alpha1.TenantNamespace{}).
		Complete(r)
}
