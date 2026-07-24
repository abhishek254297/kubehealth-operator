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

package controller

import (
	"context"
	"log/slog"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1 "github.com/abhishek254297/kubehealth-operator/api/v1"
)

// KubeHealthReconciler reconciles a KubeHealth object
type KubeHealthReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apps.abhishek.dev,resources=kubehealths,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.abhishek.dev,resources=kubehealths/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps.abhishek.dev,resources=kubehealths/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch

// Reconcile is part of the main Kubernetes reconciliation loop.
func (r *KubeHealthReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	slog.Info("Reconciling KubeHealth", "resource", req.Name)

	// Fetch KubeHealth Custom Resource
	var kubeHealth appsv1.KubeHealth

	if err := r.Get(ctx, req.NamespacedName, &kubeHealth); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Namespace to monitor (default: same namespace as CR)
	namespace := kubeHealth.Spec.Namespace
	if namespace == "" {
		namespace = req.Namespace
	}

	// Restart threshold (default: 3)
	threshold := kubeHealth.Spec.RestartThreshold
	if threshold == 0 {
		threshold = 3
	}

	slog.Info(
		"Checking pod health",
		"namespace", namespace,
		"restartThreshold", threshold,
	)

	// Get all Pods
	var podList corev1.PodList
	if err := r.List(ctx, &podList, client.InNamespace(namespace)); err != nil {
		return ctrl.Result{}, err
	}

	var healthyPods int32
	var failedPods int32
	var restartingPods int32

	for _, pod := range podList.Items {

		switch pod.Status.Phase {

		case corev1.PodRunning:
			healthyPods++

		// Pending and Unknown pods are treated as unhealthy.
		case corev1.PodFailed, corev1.PodPending, corev1.PodUnknown:
			failedPods++
		}

		for _, container := range pod.Status.ContainerStatuses {
			if container.RestartCount >= threshold {
				restartingPods++
				break
			}
		}
	}

	// Decide cluster health
	if failedPods == 0 && restartingPods == 0 {
		kubeHealth.Status.ClusterStatus = "Healthy"
	} else {
		kubeHealth.Status.ClusterStatus = "Degraded"
	}

	// Update status
	kubeHealth.Status.HealthyPods = healthyPods
	kubeHealth.Status.FailedPods = failedPods
	kubeHealth.Status.RestartingPods = restartingPods
	kubeHealth.Status.LastChecked = metav1.Now()

	if err := r.Status().Update(ctx, &kubeHealth); err != nil {
		return ctrl.Result{}, err
	}

	slog.Info(
		"KubeHealth status updated",
		"resource", req.Name,
		"namespace", namespace,
		"healthyPods", healthyPods,
		"failedPods", failedPods,
		"restartingPods", restartingPods,
		"clusterStatus", kubeHealth.Status.ClusterStatus,
	)

	slog.Info("Reconcile completed", "resource", req.Name)

	// Reconcile again after 30 seconds
	return ctrl.Result{
		RequeueAfter: 30 * time.Second,
	}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KubeHealthReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.KubeHealth{}).
		Named("kubehealth").
		Complete(r)
}
