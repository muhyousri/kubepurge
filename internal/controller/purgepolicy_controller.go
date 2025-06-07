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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
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

// TODO 2,3 fetch & Delete all resources except resources with a specific label
func Purge_resources(input string) (output string, err error) {
	output = input
	return output, nil
}

func (r *PurgePolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	var purgepolicy kubepurgexyzv1.PurgePolicy

	err := r.Get(ctx, req.NamespacedName, &purgepolicy)
	if err != nil {

	}
	schedule := purgepolicy.Spec.Schedule
	resources := purgepolicy.Spec.Resources
	targetNamespace := purgepolicy.Spec.TargetNamespace

	fmt.Printf("schedule is %s, resource are %s, targetNamespace is %s", schedule, resources, targetNamespace)

	//
	// TODO 1- process cron format and compare with current date [Done]
	c := cron.New()
	c.AddFunc(schedule, func() {
		result, err := Purge_resources(purgepolicy.Name)
		if err != nil {
			fmt.Println("error")
		} else {
			fmt.Printf("%v", result)
		}
	})
	c.Start()

	return ctrl.Result{}, nil
}

func (r *PurgePolicyReconciler) createOrUpdatePurgeStatus(ctx context.Context, policy *kubepurgexyzv1.PurgePolicy, purgedResources map[string]string) error {
	statusName := fmt.Sprintf("%s-status", policy.Name)
	
	var purgeStatus kubepurgexyzv1.PurgeStatus
	err := r.Get(ctx, types.NamespacedName{Name: statusName, Namespace: policy.Namespace}, &purgeStatus)
	
	if err != nil && errors.IsNotFound(err) {
		purgeStatus = kubepurgexyzv1.PurgeStatus{
			ObjectMeta: metav1.ObjectMeta{
				Name:      statusName,
				Namespace: policy.Namespace,
				Labels: map[string]string{
					"kubepurge.xyz/policy": policy.Name,
				},
			},
			Spec: kubepurgexyzv1.PurgeStatusSpec{
				CleanedNamespace: policy.Spec.TargetNamespace,
				LastPurgeTime:    metav1.Now(),
				PurgedResources:  purgedResources,
			},
		}
		
		return r.Create(ctx, &purgeStatus)
	} else if err != nil {
		return err
	}
	
	purgeStatus.Spec.LastPurgeTime = metav1.Now()
	purgeStatus.Spec.PurgedResources = purgedResources
	purgeStatus.Spec.CleanedNamespace = policy.Spec.TargetNamespace
	
	return r.Update(ctx, &purgeStatus)
}

// SetupWithManager sets up the controller with the Manager.
func (r *PurgePolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kubepurgexyzv1.PurgePolicy{}).
		Named("purgepolicy").
		Complete(r)
}
