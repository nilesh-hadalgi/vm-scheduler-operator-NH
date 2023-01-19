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

	"github.com/go-logr/logr"
	awsv1 "github.com/nilesh-hadalgi/vm-scheduler-operator-NH/apis/aws/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
)

// AWSVMSchedulerStartNHReconciler reconciles a AWSVMSchedulerStartNH object
type AWSVMSchedulerStartNHReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

//+kubebuilder:rbac:groups=aws.xyzcompany.com,resources=awsvmschedulerstartnhs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=aws.xyzcompany.com,resources=awsvmschedulerstartnhs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=aws.xyzcompany.com,resources=awsvmschedulerstartnhs/finalizers,verbs=update
// +kubebuilder:rbac:groups=batch,resources=cronjobs;jobs,verbs=get;list;watch;create;update;patch;delete

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

	// Check existing schedule
	if found.Spec.Schedule != "" {
		currentStartSchedule = found.Spec.Schedule
	}

	if newStartSchedule != currentStartSchedule {
		found.Spec.Schedule = newStartSchedule
		applyChange = true
	}

	// Check existing image
	if found.Spec.JobTemplate.Spec.Template.Spec.Containers != nil {
		currentImage = found.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Image
	}

	if newImage != currentImage {
		found.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Image = newImage
		applyChange = true
	}

	// Check instanceIds
	if found.Spec.JobTemplate.Spec.Template.Spec.Containers != nil {
		currentInstanceIds = found.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Env[0].Value
		log.Info(currentInstanceIds)
	}

	if newInstanceIds != currentInstanceIds {
		found.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Env[0].Value = newInstanceIds
		applyChange = true
	}

	log.Info(currentInstanceIds)
	log.Info(currentImage)
	log.Info(currentStartSchedule)

	log.Info(strconv.FormatBool(applyChange))

	if applyChange {
		log.Info(strconv.FormatBool(applyChange))
		err = r.Client.Update(ctx, found)
		if err != nil {
			log.Error(err, "Failed to update CronJob", "CronJob.Namespace", found.Namespace, "CronJob.Name", found.Name)
			return ctrl.Result{}, err
		}
		// Spec updated - return and requeue
		return ctrl.Result{Requeue: true}, nil
	}

	// Update the AWSVMScheduler status
	// TODO: Define what needs to be added in status. Currently adding just instanceIds
	if !reflect.DeepEqual(currentInstanceIds, awsVMScheduler.Status.VMStartStatus) {
		awsVMScheduler.Status.VMStartStatus = currentInstanceIds
		// awsVMScheduler.Status.VMStopStatus = currentInstanceIds
		err := r.Client.Status().Update(ctx, awsVMScheduler)
		if err != nil {
			log.Error(err, "Failed to update awsVMScheduler status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil

}

// SetupWithManager sets up the controller with the Manager.
func (r *AWSVMSchedulerStartNHReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&awsv1.AWSVMSchedulerStartNH{}).
		Complete(r)
}

func (r *AWSVMSchedulerStartNHReconciler) cronJobForAWSVMSchedulerStartNH(awsVMScheduler *awsv1.AWSVMSchedulerStartNH) *batchv1.CronJob {

	cron := &batchv1.CronJob{
		ObjectMeta: v1.ObjectMeta{
			Name:      awsVMScheduler.Name,
			Namespace: awsVMScheduler.Namespace,
			Labels:    AWSVMSchedulerStartNHLabels(awsVMScheduler, "awsVMScheduler"),
		},
		Spec: batchv1.CronJobSpec{
			Schedule: awsVMScheduler.Spec.StartSchedule,
			// TODO: Add Stop schedule
			//Schedule:  awsVMScheduler.Spec.StopSchedule,
			JobTemplate: batchv1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{{
								Name:  awsVMScheduler.Name,
								Image: awsVMScheduler.Spec.Image,
								Env: []corev1.EnvVar{
									{
										Name:  "ec2_instanceIds",
										Value: awsVMScheduler.Spec.InstanceIds,
									},
									{
										Name: "AWS_ACCESS_KEY_ID",
										ValueFrom: &corev1.EnvVarSource{
											SecretKeyRef: &corev1.SecretKeySelector{
												LocalObjectReference: corev1.LocalObjectReference{
													Name: "aws-secret",
												},
												Key: "aws-access-key-id",
											},
										},
									},
									{
										Name: "AWS_SECRET_ACCESS_KEY",
										ValueFrom: &corev1.EnvVarSource{
											SecretKeyRef: &corev1.SecretKeySelector{
												LocalObjectReference: corev1.LocalObjectReference{
													Name: "aws-secret",
												},
												Key: "aws-secret-access-key",
											},
										},
									},
									{
										Name: "AWS_DEFAULT_REGION",
										ValueFrom: &corev1.EnvVarSource{
											SecretKeyRef: &corev1.SecretKeySelector{
												LocalObjectReference: corev1.LocalObjectReference{
													Name: "aws-secret",
												},
												Key: "region",
											},
										},
									}},
							}},
							RestartPolicy: "OnFailure",
						},
					},
				},
			},
		},
	}
	// Set awsVMScheduler instance as the owner and controller
	ctrl.SetControllerReference(awsVMScheduler, cron, r.Scheme)
	return cron
}

// // createResources will create and update the  resource which are required
// func (r *AWSVMSchedulerReconciler) createResources(awsVMScheduler *awsv1.AWSVMScheduler, request ctrl.Request) error {

// 	log := r.Log.WithValues("AWSVMScheduler", request.NamespacedName)
// 	log.Info("Creating   resources ...")

// 	// Check if the cronJob is created, if not create one
// 	if err := r.createCronJob(awsVMScheduler); err != nil {
// 		log.Error(err, "Failed to create the CronJob")
// 		return err
// 	}

// 	return nil
// }

// // Check if the cronJob is created, if not create one
// func (r *AWSVMSchedulerReconciler) createCronJob(awsVMScheduler *awsv1.AWSVMScheduler) error {
// 	if _, err := services.FetchCronJob(awsVMScheduler.Name, awsVMScheduler.Namespace); err != nil {
// 		if err := r.client.Create(context.TODO(), resources.NewAWSVMSchedulerCronJob(awsVMScheduler)); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

func AWSVMSchedulerStartNHLabels(v *awsv1.AWSVMSchedulerStartNH, tier string) map[string]string {
	return map[string]string{
		"app":                      "AWSVMSchedulerStartNH",
		"AWSVMSchedulerStartNH_cr": v.Name,
		"tier":                     tier,
	}
}
