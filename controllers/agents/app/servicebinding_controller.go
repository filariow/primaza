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
	"path"
	"strings"
	"time"

	"fmt"

	primazaiov1alpha1 "github.com/primaza/primaza/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// ServiceBindingReconciler reconciles a ServiceBinding object
type ServiceBindingReconciler struct {
	client.Client
	Scheme             *runtime.Scheme
	mountPathDir       string
	volumeNamePrefix   string
	volumeName         string
	unstructuredVolume map[string]interface{}
}

// ServiceBindingRoot points to the environment variable in the container
// which is used as the volume mount path.  In the absence of this
// environment variable, `/bindings` is used as the volume mount path.
// Refer: https://github.com/k8s-service-bindings/spec#reconciler-implementation
const ServiceBindingRoot = "SERVICE_BINDING_ROOT"

type errorList []error

//+kubebuilder:rbac:groups=primaza.io,resources=servicebindings,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=primaza.io,resources=servicebindings/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ServiceBinding object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *ServiceBindingReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	l.Info("starting reconciliation")
	fmt.Println("starting reconciliation----------------------")

	var serviceBinding primazaiov1alpha1.ServiceBinding

	l.Info("retrieving ServiceBinding object", "ServiceBinding", serviceBinding)
	if err := r.Get(ctx, req.NamespacedName, &serviceBinding); err != nil {
		l.Error(err, "unable to retrieve ServiceBinding")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	l.Info("ServiceBinding object retrieved", "ServiceBinding", serviceBinding)
	fmt.Println("ServiceBinding object retrieved", "ServiceBinding----:s[;[", serviceBinding)

	// examine DeletionTimestamp to determine if object is under deletion
	var secretName string
	if serviceBinding.Spec.ServiceEndpointDefinitionSecret != "" {
		secretName = serviceBinding.Spec.ServiceEndpointDefinitionSecret
	}
	// if !serviceBinding.HasDeletionTimestamp() {
	// 	// The object is not being deleted, so if it does not have the finalizer,
	// 	// then lets add the finalizer and update the object. This is equivalent
	// 	// registering the finalizer.
	// 	if err := r.Update(ctx, &serviceBinding); err != nil {
	// 		return ctrl.Result{}, err
	// 	}
	// } else {
	// 	l.Info("Deleted, unbind the application")
	// 	applications, result, err := r.getApplication(ctx, req, serviceBinding, secretName)
	// 	if err != nil {
	// 		return result, err
	// 	}
	// 	result, err = r.unbindApplications(ctx, req, serviceBinding, applications...)
	// 	if err != nil {
	// 		return result, err
	// 	}
	// 	if err := r.Update(ctx, &serviceBinding); err != nil {
	// 		return ctrl.Result{}, err
	// 	}

	// }
	volumeNamePrefix := serviceBinding.Name
	if len(volumeNamePrefix) > 56 {
		volumeNamePrefix = volumeNamePrefix[:56]
	}
	r.volumeName = volumeNamePrefix + "-" + secretName
	r.mountPathDir = serviceBinding.Name
	sp := &v1.SecretProjection{
		LocalObjectReference: v1.LocalObjectReference{
			Name: secretName,
		}}

	volumeProjection := &v1.Volume{
		Name: r.volumeName,
		VolumeSource: v1.VolumeSource{
			Projected: &v1.ProjectedVolumeSource{
				Sources: []v1.VolumeProjection{{Secret: sp}},
			},
		},
	}
	l.Info("converting the volumeProjection to an unstructured object", "Volume", volumeProjection)
	var err error
	r.unstructuredVolume, err = runtime.DefaultUnstructuredConverter.ToUnstructured(volumeProjection)
	if err != nil {
		l.Error(err, "unable to convert volumeProjection to an unstructured object")
		return ctrl.Result{}, err
	}
	applications, result, err := r.getApplication(ctx, req, serviceBinding, secretName)
	if err != nil {
		return result, err
	}
	secretLookupKey := client.ObjectKey{Name: serviceBinding.Spec.ServiceEndpointDefinitionSecret, Namespace: req.NamespacedName.Namespace}
	psSecret := &v1.Secret{}
	l.Info(fmt.Sprintf("Length of applications: %v", len(applications)))
	if err := r.Get(ctx, secretLookupKey, psSecret); err != nil {
		return ctrl.Result{}, err
	}
	return r.bindApplications(ctx, req, serviceBinding, psSecret, applications...)
}

func (r *ServiceBindingReconciler) bindApplications(ctx context.Context, req ctrl.Request,
	sb primazaiov1alpha1.ServiceBinding, psSecret *v1.Secret, applications ...unstructured.Unstructured) (ctrl.Result, error) {

	l := log.FromContext(ctx)

	var el errorList
	for index, application := range applications {
		containersPaths := [][]string{}
		envsPaths := [][]string{}
		volumeMountsPaths := [][]string{}
		volumesPath := []string{"spec", "template", "spec", "volumes"}
		envsOrVolumeMountsFound := false

		containersPaths = append(containersPaths,
			[]string{"spec", "template", "spec", "containers"},
			[]string{"spec", "template", "spec", "initContainers"},
		)
		l.Info("referencing the volume in an unstructured object")
		volumes, found, err := unstructured.NestedSlice(application.Object, volumesPath...)
		if !found {
			l.Info("volumes not found in the application object")
		}
		if err != nil {
			l.Error(err, "unable to reference the volumes in the application object")
			return ctrl.Result{}, err
		}
		l.Info("Volumes values", "volumes", volumes)

		volumeFound := false

		for i, volume := range volumes {
			l.Info("Volume", "volume", volume)
			if strings.HasPrefix(volume.(map[string]interface{})["name"].(string), r.volumeNamePrefix) {
				volumes[i] = r.unstructuredVolume
				volumeFound = true
			}
		}

		if !volumeFound {
			volumes = append(volumes, r.unstructuredVolume)
		}
		l.Info("setting the updated volumes into the application using the unstructured object")
		if err := unstructured.SetNestedSlice(application.Object, volumes, volumesPath...); err != nil {
			return ctrl.Result{}, err
		}
		l.Info("application object after setting the update volume", "Application", application)

		if !envsOrVolumeMountsFound {
			for _, containersPath := range containersPaths {
				l.Info("referencing containers in an unstructured object")
				containers, found, err := unstructured.NestedSlice(application.Object, containersPath...)
				if !found {
					e := &field.Error{Type: field.ErrorTypeRequired, Field: strings.Join(containersPath, "."), Detail: "no containers"}
					l.Info("containers not found in the application object", "error", e)
				}
				if err != nil {
					l.Error(err, "unable to reference containers in the application object")
					return ctrl.Result{}, err
				}

			CONTAINERS_OUTER:
				for i := range containers {
					container := &containers[i]
					l.Info("updating container", "container", container)
					c := &v1.Container{}
					u := (*container).(map[string]interface{})
					if err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, c); err != nil {
						return ctrl.Result{}, err
					}

					if len(sb.Spec.Application.Containers) > 0 {
						found := false
						count := 0
						for _, v := range sb.Spec.Application.Containers {
							l.Info("container", "container value", v, "name", c.Name)
							if v.StrVal == c.Name {
								break
							}
							found = true
							count++
						}
						if found && len(sb.Spec.Application.Containers) == count {
							continue CONTAINERS_OUTER
						}

					}

					for _, e := range sb.Spec.Env {
						c.Env = append(c.Env, v1.EnvVar{
							Name:  e.Name,
							Value: string(psSecret.Data[e.Key]),
						})

					}
					mountPath := ""
					for _, e := range c.Env {
						if e.Name == ServiceBindingRoot {
							mountPath = path.Join(e.Value, r.mountPathDir)
							break
						}
					}

					if mountPath == "" {
						mountPath = path.Join("/bindings", r.mountPathDir)
						c.Env = append(c.Env, v1.EnvVar{
							Name:  ServiceBindingRoot,
							Value: "/bindings",
						})
					}

					volumeMount := v1.VolumeMount{
						Name:      r.volumeName,
						MountPath: mountPath,
						ReadOnly:  true,
					}

					volumeMountFound := false
					for j, vm := range c.VolumeMounts {
						if strings.HasPrefix(vm.Name, r.volumeNamePrefix) {
							c.VolumeMounts[j] = volumeMount
							volumeMountFound = true
							break
						}
					}

					if !volumeMountFound {
						c.VolumeMounts = append(c.VolumeMounts, volumeMount)
					}

					nu, err := runtime.DefaultUnstructuredConverter.ToUnstructured(c)
					if err != nil {
						return ctrl.Result{}, err
					}

					containers[i] = nu
				}

				l.Info("updated cntainer with volume and volume mounts", "containers", containers)

				l.Info("setting the updated containers into the application using the unstructured object")
				if err := unstructured.SetNestedSlice(application.Object, containers, containersPath...); err != nil {
					return ctrl.Result{}, err
				}
				l.Info("application object after setting the updated containers", "Application", application)
			}
		} else {

			mountPath := ""

			for _, envsPath := range envsPaths {
				l.Info("referencing env in an unstructured object")
				env, found, err := unstructured.NestedMap(application.Object, envsPath...)
				if !found {
					e := &field.Error{Type: field.ErrorTypeRequired, Field: strings.Join(envsPath, "."), Detail: "empty env"}
					l.Info("env not found in the application object", "error", e)
				}
				if err != nil {
					l.Error(err, "unable to reference env in the application object")
					return ctrl.Result{}, err
				}

				ev := []v1.EnvVar{}

				if err := runtime.DefaultUnstructuredConverter.FromUnstructured(env, ev); err != nil {
					return ctrl.Result{}, err
				}

				for _, e := range ev {
					if e.Name == ServiceBindingRoot {
						mountPath = path.Join(e.Value, r.mountPathDir)
						break
					}
				}

				if mountPath == "" {
					mountPath = path.Join("/bindings", r.mountPathDir)
					ev = append(ev, v1.EnvVar{
						Name:  ServiceBindingRoot,
						Value: "/bindings",
					})
				}

				for _, e := range sb.Spec.Env {
					ev = append(ev, v1.EnvVar{
						Name:  e.Name,
						Value: string(psSecret.Data[e.Key]),
					})

				}

				evUnstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(ev)
				if err != nil {
					return ctrl.Result{}, err
				}

				l.Info("setting the updated envs into the application using the unstructured object")
				if err := unstructured.SetNestedMap(application.Object, evUnstructured, envsPath...); err != nil {
					return ctrl.Result{}, err
				}
				l.Info("application object after setting the updated envs", "Application", application)

			}

			for _, volumeMountsPath := range volumeMountsPaths {
				l.Info("referencing volumeMount in an unstructured object")
				volumeMount, found, err := unstructured.NestedMap(application.Object, volumeMountsPath...)
				if !found {
					e := &field.Error{Type: field.ErrorTypeRequired, Field: strings.Join(volumeMountsPath, "."), Detail: "empty volumeMount"}
					l.Info("volumeMount not found in the application object", "error", e)
				}
				if err != nil {
					l.Error(err, "unable to reference volumeMount in the application object")
					return ctrl.Result{}, err
				}

				vm := []v1.VolumeMount{}
				if err := runtime.DefaultUnstructuredConverter.FromUnstructured(volumeMount, vm); err != nil {
					return ctrl.Result{}, err
				}

				volumeMountSB := v1.VolumeMount{
					Name:      r.volumeName,
					MountPath: mountPath,
					ReadOnly:  true,
				}

				volumeMountFound := false
				for j, v := range vm {
					if strings.HasPrefix(v.Name, r.volumeNamePrefix) {
						vm[j] = volumeMountSB
						volumeMountFound = true
						break
					}
				}

				if !volumeMountFound {
					vm = append(vm, volumeMountSB)
				}

				vmUnstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(vm)
				if err != nil {
					return ctrl.Result{}, err
				}

				l.Info("setting the updated volumeMounts into the application using the unstructured object")
				if err := unstructured.SetNestedMap(application.Object, vmUnstructured, volumeMountsPath...); err != nil {
					return ctrl.Result{}, err
				}
				l.Info("application object after setting the updated volumeMounts", "Application", application)

			}
		}

		var conditionStatus metav1.ConditionStatus
		var reason, state string

		conditionStatus = "True"
		reason = "Success"
		state = "Ready"

		l.Info("updating the application with updated volumes and volumeMounts")
		if err := r.Update(ctx, &applications[index]); err != nil {
			l.Error(err, "unable to update the application", "application", application)
			conditionStatus = "False"
			state = "Malformed"
			reason = "failure"
		}
		l.Info("set the status of the service binding")
		_, err = r.setStatus(ctx, psSecret.Name, sb, conditionStatus, reason, state)
		if err != nil {
			el = append(el, err)
		}
	}
	if len(el) > 0 {
		return ctrl.Result{}, el
	}
	return ctrl.Result{}, nil
}

func (el errorList) Error() string {
	msg := ""
	for _, e := range el {
		msg += e.Error() + " "
	}
	return msg
}

func (r *ServiceBindingReconciler) setStatus(ctx context.Context, secretName string,
	sb primazaiov1alpha1.ServiceBinding, conditionStatus metav1.ConditionStatus, reason, state string) (ctrl.Result, error) {
	l := log.FromContext(ctx)
	conditionFound := false
	for k, cond := range sb.Status.Conditions {
		if cond.Type == primazaiov1alpha1.ServiceBindingConditionReady {
			cond.Status = conditionStatus
			sb.Status.Conditions[k].Status = cond.Status
			conditionFound = true
		}
	}

	if !conditionFound {
		c := metav1.Condition{
			LastTransitionTime: metav1.NewTime(time.Now()),
			Type:               primazaiov1alpha1.ConditionReady,
			Status:             conditionStatus,
			Reason:             reason,
		}
		sb.Status.Conditions = append(sb.Status.Conditions, c)
		sb.Status.State = state
	}

	l.Info("updating the service binding status")
	if err := r.Status().Update(ctx, &sb); err != nil {
		l.Error(err, "unable to update the service binding", "ServiceBinding", sb)
		return ctrl.Result{}, err
	}
	l.Info("service binding status updated", "ServiceBinding", sb)

	return ctrl.Result{}, nil
}

func (r *ServiceBindingReconciler) getApplication(ctx context.Context, req ctrl.Request,
	sb primazaiov1alpha1.ServiceBinding, secretName string) ([]unstructured.Unstructured, ctrl.Result, error) {
	var applications []unstructured.Unstructured
	var conditionStatus metav1.ConditionStatus
	var reason, state string
	l := log.FromContext(ctx)
	if sb.Spec.Application.Name != "" {
		applicationLookupKey := client.ObjectKey{Name: sb.Spec.Application.Name, Namespace: req.NamespacedName.Namespace}

		application := &unstructured.Unstructured{
			Object: map[string]interface{}{
				"kind":       sb.Spec.Application.Kind,
				"apiVersion": sb.Spec.Application.APIVersion,
				"metadata": map[string]interface{}{
					"name":      sb.Spec.Application.Name,
					"namespace": req.NamespacedName.Namespace,
				},
			},
		}

		l.Info("retrieving the application object", "Application", application)
		if err := r.Get(ctx, applicationLookupKey, application); err != nil {
			reason = "unableToRetrieveApplication"
			l.Error(err, reason)
			conditionStatus = "False"
			state = "Malformed"
			result, err := r.setStatus(ctx, secretName, sb, conditionStatus, reason, state)
			return []unstructured.Unstructured{}, result, err
		}
		l.Info("application object retrieved", "Application", application)
		applications = append(applications, *application)
	}

	if sb.Spec.Application.Selector != nil {
		applicationList := &unstructured.UnstructuredList{
			Object: map[string]interface{}{
				"kind":       sb.Spec.Application.Kind,
				"apiVersion": sb.Spec.Application.APIVersion,
			},
		}

		l.Info("retrieving the application objects", "Application", applicationList)
		opts := &client.ListOptions{
			LabelSelector: labels.Set(sb.Spec.Application.Selector.MatchLabels).AsSelector(),
			Namespace:     req.NamespacedName.Namespace,
		}

		if err := r.List(ctx, applicationList, opts); err != nil {
			reason = "unableToRetrieveApplication"
			l.Error(err, reason)
			conditionStatus = "False"
			state = "Malformed"
			result, err := r.setStatus(ctx, secretName, sb, conditionStatus, reason, state)
			return []unstructured.Unstructured{}, result, err
		}
		l.Info("application objects retrieved", "Application", applicationList)
		applications = append(applications, applicationList.Items...)
	}
	if len(applications) == 0 {
		// Requeue with a time interval is required as the applications is not available to reconcile
		// In future, probably watching for applications os specific types (Deployment, CronJob etc.) based
		// on label can be introduced or a webhook can detect application change and trigger reconciliation
		return applications, ctrl.Result{RequeueAfter: time.Minute * 1}, nil
	}
	return applications, ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ServiceBindingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&primazaiov1alpha1.ServiceBinding{}).
		Complete(r)
}
