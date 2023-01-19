/*
Copyright 2023.

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

package aws

import (
	"context"
    "reflect"
    "strconv"
    awsv1 "github.com/nilesh-hadalgi/vm-scheduler-operator-NH/apis/aws/v1"
	"github.com/go-logr/logr"
    batchv1 "k8s.io/api/batch/v1"
    "k8s.io/apimachinery/pkg/api/errors"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/apimachinery/pkg/types"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/client"
    corev1 "k8s.io/api/core/v1"
    v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    ctrllog "sigs.k8s.io/controller-runtime/pkg/log"

	)

// AWSVMSchedulerStartNHReconciler reconciles a AWSVMSchedulerStartNH object
type AWSVMSchedulerStartNHReconciler struct {
	client.Client
	Scheme *runtime.Scheme
    Log logr.Logger
}

//+kubebuilder:rbac:groups=aws.xyzcompany.com,resources=awsvmschedulerstartnhs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=aws.xyzcompany.com,resources=awsvmschedulerstartnhs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=aws.xyzcompany.com,resources=awsvmschedulerstartnhs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the AWSVMSchedulerStartNH object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *AWSVMSchedulerStartNHReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	//_ = log.FromContext(ctx)
    log := ctrllog.FromContext(ctx)
    log.Info("Reconciling AWSVMSchedulerStartNH")

	// TODO(user): your logic here
    awsVMScheduler := &awsv1.AWSVMSchedulerStartNH{}
    log.Info(req.NamespacedName.Name)
    err := r.Client.Get(ctx, req.NamespacedName, awsVMScheduler)
    if err != nil {
        if errors.IsNotFound(err) {
            // Request object not found, could have been deleted after reconcile request.
            // Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
            // Return and don't requeue
            log.Info("awsVMScheduler resource not found. Ignoring since object must be deleted.")
            return ctrl.Result{}, nil
        }
        // Error reading the object - requeue the request.
        log.Error(err, "Failed to get awsVMScheduler.")
        return ctrl.Result{}, err
    }

    log.Info(awsVMScheduler.Name)
    found := &batchv1.CronJob{}
    err = r.Client.Get(ctx, types.NamespacedName{Name: awsVMScheduler.Name, Namespace: awsVMScheduler.Namespace}, found)
    if err != nil && errors.IsNotFound(err) {
        // Define a new CronJob
        cron := r.cronJobForAWSVMSchedulerStartNH(awsVMScheduler)
        log.Info("Creating a new CronJob", "CronJob.Namespace", cron.Namespace, "CronJob.Name", cron.Name)
        err = r.Client.Create(ctx, cron)
        if err != nil {
            log.Error(err, "Failed to create new CronJob", "CronJob.Namespace", cron.Namespace, "CronJob.Name", cron.Name)
            return ctrl.Result{}, err
        }
        // Cronjob created successfully - return and requeue
        return ctrl.Result{Requeue: true}, nil
    } else if err != nil {
        log.Error(err, "Failed to get Cronjob")
        return ctrl.Result{}, err
    }
    applyChange := false
    // Ensure image name is correct, update image if required
    newInstanceIds := awsVMScheduler.Spec.InstanceIds
    log.Info(newInstanceIds)
    newStartSchedule := awsVMScheduler.Spec.StartSchedule
    log.Info(newStartSchedule)
    newImage := awsVMScheduler.Spec.Image
    log.Info(newImage)
    var currentImage string = ""
    var currentStartSchedule string = ""
    var currentInstanceIds string = ""

    if found.Spec.Schedule != "" {
        currentStartSchedule = found.Spec.Schedule
    }
    if newStartSchedule != currentStartSchedule {
        found.Spec.Schedule = newStartSchedule
        applyChange = true
    }

    if found.Spec.JobTemplate.Spec.Template.Spec.Containers != nil {
        currentImage = found.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Image
    }



}

// SetupWithManager sets up the controller with the Manager.
func (r *AWSVMSchedulerStartNHReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&awsv1.AWSVMSchedulerStartNH{}).
		Complete(r)
}
