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

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	kubepurgexyzv1 "github.com/muhyousri/kubepurge/api/v1"
	"github.com/robfig/cron/v3"
)

// PurgePolicyReconciler reconciles a PurgePolicy object
type PurgePolicyReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	cronJobs map[string]*cron.Cron
}

// +kubebuilder:rbac:groups=kubepurge.xyz,resources=purgepolicies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kubepurge.xyz,resources=purgepolicies/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=kubepurge.xyz,resources=purgepolicies/finalizers,verbs=update

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile

// TODO 2,3 fetch & Delete all resources except resources with a specific label
func Purge_resources(input string) (output string, err error) {
	output = input
	return output, nil
}

func (r *PurgePolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	
	if r.cronJobs == nil {
		r.cronJobs = make(map[string]*cron.Cron)
	}

	var purgepolicy kubepurgexyzv1.PurgePolicy
	err := r.Get(ctx, req.NamespacedName, &purgepolicy)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("PurgePolicy resource not found. Stopping cron job if exists.")
			if cronJob, exists := r.cronJobs[req.NamespacedName.String()]; exists {
				cronJob.Stop()
				delete(r.cronJobs, req.NamespacedName.String())
			}
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Unable to fetch PurgePolicy")
		return ctrl.Result{}, err
	}

	schedule := purgepolicy.Spec.Schedule
	resources := purgepolicy.Spec.Resources
	targetNamespace := purgepolicy.Spec.TargetNamespace

	logger.Info("Processing PurgePolicy", "schedule", schedule, "resources", resources, "targetNamespace", targetNamespace)

	cronJobKey := req.NamespacedName.String()
	
	if existingCron, exists := r.cronJobs[cronJobKey]; exists {
		existingCron.Stop()
	}

	c := cron.New()
	_, err = c.AddFunc(schedule, func() {
		logger.Info("Executing purge for policy", "policy", purgepolicy.Name)
		// TODO: Implement actual resource purging logic
	})
	
	if err != nil {
		logger.Error(err, "Failed to add cron job", "schedule", schedule)
		return ctrl.Result{}, err
	}
	
	c.Start()
	r.cronJobs[cronJobKey] = c

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PurgePolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kubepurgexyzv1.PurgePolicy{}).
		Named("purgepolicy").
		Complete(r)
}
