/*
Copyright 2023 The Primaza Authors.

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

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/google/uuid"
	primazaiov1alpha1 "github.com/primaza/primaza/api/v1alpha1"
)

// ServiceClaimReconciler reconciles a ServiceClaim object
type ServiceClaimReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=primaza.io,resources=serviceclaims,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=primaza.io,resources=serviceclaims/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=primaza.io,resources=serviceclaims/finalizers,verbs=update

// TODO: Remove this later once Primaza ServiceBinding is implemented.
//+kubebuilder:rbac:groups=servicebinding.io,resources=servicebindings,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=,resources=secrets,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ServiceClaim object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *ServiceClaimReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	l.V(0).Info("starting reconciliation")

	var sclaim primazaiov1alpha1.ServiceClaim

	if err := r.Get(ctx, req.NamespacedName, &sclaim); err != nil {
		l.V(0).Info("unable to retrieve ServiceClaim", "error", err)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var rsl primazaiov1alpha1.RegisteredServiceList
	if err := r.List(ctx, &rsl); err != nil {
		l.V(0).Info("unable to retrieve RegisteredServiceList", "error", err)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.NamespacedName.Name,
			Namespace: req.NamespacedName.Namespace,
		},
		StringData: map[string]string{},
	}

	var rsName string
	for _, rs := range rsl.Items {
		for _, sci := range rs.Spec.ServiceClassIdentity {
			for _, sclaimSCI := range sclaim.Spec.ServiceClassIdentity {
				if sclaimSCI == sci {
					for _, e := range rs.Spec.Constraints.Environments {
						if sclaim.Spec.EnvironmentTag == e {
							for _, sed := range rs.Spec.ServiceEndpointDefinition {
								if sed.Value != "" && itemContains(sclaim.Spec.ServiceEndpointDefinitionKeys, sed.Name) {
									secret.StringData[sed.Name] = sed.Value
								} else if sed.ValueFromSecret.Key != "" && itemContains(sclaim.Spec.ServiceEndpointDefinitionKeys, sed.Name) {
									sec := &corev1.Secret{
										ObjectMeta: metav1.ObjectMeta{
											Name:      sed.ValueFromSecret.Name,
											Namespace: req.NamespacedName.Namespace,
										},
									}
									if err := r.Get(ctx, req.NamespacedName, sec); err != nil {
										l.V(0).Info("unable to retrieve Secret", "error", err)
										return ctrl.Result{}, client.IgnoreNotFound(err)
									}
									secret.StringData[sed.Name] = sec.StringData[sed.ValueFromSecret.Key]
								}
							}
							rsName = rs.Name
						}
					}
				}
			}
		}
	}

	var cel primazaiov1alpha1.ClusterEnvironmentList
	if err := r.List(ctx, &cel); err != nil {
		l.V(0).Info("error fetching ClusterEnvironmentList", "error", err)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	for _, ce := range cel.Items {
		if sclaim.Spec.EnvironmentTag == ce.Spec.EnvironmentName {
			sn := ce.Spec.ClusterContextSecret
			k := client.ObjectKey{Namespace: ce.Namespace, Name: sn}
			var s corev1.Secret
			if err := r.Get(ctx, k, &s); err != nil {
				return ctrl.Result{}, err
			}

			kubeconfig := s.Data["kubeconfig"]
			cg, err := clientcmd.RESTConfigFromKubeConfig(kubeconfig)
			if err != nil {
				return ctrl.Result{}, err
			}

			cs, err := kubernetes.NewForConfig(cg)
			if err != nil {
				return ctrl.Result{}, err
			}

			_, err = cs.CoreV1().Secrets(ce.Namespace).Create(ctx, secret, metav1.CreateOptions{})
			if err != nil {
				return ctrl.Result{}, err
			}

			dynamicClient, err := dynamic.NewForConfig(cg)

			sb := &unstructured.Unstructured{
				Object: map[string]interface{}{
					"kind":       "ServiceBinding",
					"apiVersion": "servicebinding.io/v1beta1",
					"metadata": map[string]interface{}{
						"name":      req.NamespacedName.Name,
						"namespace": req.NamespacedName.Namespace,
					},
					"spec": map[string]interface{}{
						"service": map[string]interface{}{
							"apiVersion": "v1",
							"kind":       "Secret",
							"name":       req.NamespacedName.Name,
						},
						"workload": map[string]interface{}{
							"apiVersion": sclaim.Spec.Application.APIVersion,
							"kind":       sclaim.Spec.Application.Kind,
							"selector": map[string]interface{}{
								"matchLabels": sclaim.Spec.Application.Selector.MatchLabels,
							},
						},
					},
				}}
			gvr := schema.GroupVersionResource{
				Group:    "servicebinding.io",
				Version:  "v1beta1",
				Resource: "servicebindings",
			}
			_, err = dynamicClient.Resource(gvr).Namespace(req.NamespacedName.Namespace).Create(ctx, sb, metav1.CreateOptions{})
			if err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	sclaim.Status.State = "Resolved"
	sclaim.Status.ClaimID = uuid.New().String()
	sclaim.Status.RegisteredService = rsName
	sclaim.Status.ServiceBinding = req.NamespacedName.Name
	if err := r.Status().Update(ctx, &sclaim); err != nil {
		l.Error(err, "unable to update the ServiceClaim", "ServiceClaim", sclaim)
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func itemContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// SetupWithManager sets up the controller with the Manager.
func (r *ServiceClaimReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&primazaiov1alpha1.ServiceClaim{}).
		Complete(r)
}
