/*
Copyright 2024.

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

package controller

import (
	"context"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	kubepurgexyzv1 "github.com/muhyousri/kubepurge/api/v1"
	"github.com/robfig/cron/v3"
)

// PurgePolicyReconciler reconciles a PurgePolicy object
type PurgePolicyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=kubepurge.xyz,resources=purgepolicies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kubepurge.xyz,resources=purgepolicies/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=kubepurge.xyz,resources=purgepolicies/finalizers,verbs=update

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile

func (r *PurgePolicyReconciler) purgeResources(ctx context.Context, policy *kubepurgexyzv1.PurgePolicy) map[string]string {
	logger := log.FromContext(ctx)
	purgedResources := make(map[string]string)

	for _, resourceType := range policy.Spec.Resources {
		switch strings.ToLower(resourceType) {
		case "pods":
			count, err := r.deleteResourcesByType(ctx, policy.Spec.TargetNamespace, "v1", "Pod")
			if err != nil {
				logger.Error(err, "Failed to delete pods")
				continue
			}
			purgedResources["pods"] = fmt.Sprintf("%d", count)

		case "deployments":
			count, err := r.deleteResourcesByType(ctx, policy.Spec.TargetNamespace, "apps/v1", "Deployment")
			if err != nil {
				logger.Error(err, "Failed to delete deployments")
				continue
			}
			purgedResources["deployments"] = fmt.Sprintf("%d", count)

		case "services":
			count, err := r.deleteResourcesByType(ctx, policy.Spec.TargetNamespace, "v1", "Service")
			if err != nil {
				logger.Error(err, "Failed to delete services")
				continue
			}
			purgedResources["services"] = fmt.Sprintf("%d", count)

		case "configmaps":
			count, err := r.deleteResourcesByType(ctx, policy.Spec.TargetNamespace, "v1", "ConfigMap")
			if err != nil {
				logger.Error(err, "Failed to delete configmaps")
				continue
			}
			purgedResources["configmaps"] = fmt.Sprintf("%d", count)

		case "secrets":
			count, err := r.deleteResourcesByType(ctx, policy.Spec.TargetNamespace, "v1", "Secret")
			if err != nil {
				logger.Error(err, "Failed to delete secrets")
				continue
			}
			purgedResources["secrets"] = fmt.Sprintf("%d", count)
		}
	}

	return purgedResources
}

func (r *PurgePolicyReconciler) deleteResourcesByType(ctx context.Context, namespace, apiVersion, kind string) (int, error) {
	gv, err := schema.ParseGroupVersion(apiVersion)
	if err != nil {
		return 0, err
	}

	gvk := schema.GroupVersionKind{
		Group:   gv.Group,
		Version: gv.Version,
		Kind:    kind,
	}

	list := &metav1.PartialObjectMetadataList{}
	list.SetGroupVersionKind(gvk)

	err = r.List(ctx, list, client.InNamespace(namespace))
	if err != nil {
		return 0, err
	}

	count := 0
	for _, item := range list.Items {
		if item.Labels["kubepurge.xyz/exclude"] == "true" {
			continue
		}

		obj := &metav1.PartialObjectMetadata{}
		obj.SetGroupVersionKind(gvk)
		obj.SetNamespace(item.Namespace)
		obj.SetName(item.Name)

		err := r.Delete(ctx, obj)
		if err != nil && !errors.IsNotFound(err) {
			return count, err
		}
		count++
	}

	return count, nil
}

func (r *PurgePolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	var purgepolicy kubepurgexyzv1.PurgePolicy

	err := r.Get(ctx, req.NamespacedName, &purgepolicy)
	if err != nil {
		return ctrl.Result{}, err
	}
	schedule := purgepolicy.Spec.Schedule
	resources := purgepolicy.Spec.Resources
	targetNamespace := purgepolicy.Spec.TargetNamespace

	fmt.Printf("schedule is %s, resource are %s, targetNamespace is %s", schedule, resources, targetNamespace)

	//
	// TODO 1- process cron format and compare with current date [Done]
	c := cron.New()
	_, err = c.AddFunc(schedule, func() {
		result := r.purgeResources(ctx, &purgepolicy)
		fmt.Printf("%v", result)
	})
	if err != nil {
		return ctrl.Result{}, err
	}
	c.Start()

	// TODO 4- create or patch a purge status

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PurgePolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kubepurgexyzv1.PurgePolicy{}).
		Named("purgepolicy").
		Complete(r)
}
