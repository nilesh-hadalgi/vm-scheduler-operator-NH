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

package azure

import (
	"context"
	//"reflect"
	"strconv"

	"github.com/go-logr/logr"
	azurev1 "github.com/nilesh-hadalgi/vm-scheduler-operator-NH/apis/azure/v1"
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

// AzureVMSchedulerStopNHReconciler reconciles a AzureVMSchedulerStopNH object
type AzureVMSchedulerStopNHReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

//+kubebuilder:rbac:groups=azure.xyzcompany.com,resources=azurevmschedulerstopnhs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=azure.xyzcompany.com,resources=azurevmschedulerstopnhs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=azure.xyzcompany.com,resources=azurevmschedulerstopnhs/finalizers,verbs=update
// +kubebuilder:rbac:groups=batch,resources=cronjobs;jobs,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the AzureVMSchedulerStopNH object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *AzureVMSchedulerStopNHReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	//_ = log.FromContext(ctx)
	//_ = log.FromContext(ctx)
	log := ctrllog.FromContext(ctx)
	log.Info("Reconciling AzureVMSchedulerStopNH")

	// TODO(user): your logic here
	// Fetch the azureVMScheduler instance
	azureVMScheduler := &azurev1.AzureVMSchedulerStopNH{}
	log.Info(req.NamespacedName.Name)

	err := r.Client.Get(ctx, req.NamespacedName, azureVMScheduler)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("azureVMScheduler resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get azureVMScheduler.")
		return ctrl.Result{}, err
	}

	log.Info(azureVMScheduler.Name)

	// Add const values for mandatory specs ( if left blank)
	// log.Info("Adding azureVMScheduler mandatory specs")
	// utils.AddBackupMandatorySpecs(azureVMScheduler)
	// Check if the CronJob already exists, if not create a new one

	found := &batchv1.CronJob{}
	err = r.Client.Get(ctx, types.NamespacedName{Name: azureVMScheduler.Name, Namespace: azureVMScheduler.Namespace}, found)
	// log.Info(*found.)
	if err != nil && errors.IsNotFound(err) {
		// Define a new CronJob
		cron := r.cronJobForAzureVMSchedulerStopNH(azureVMScheduler)
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

	// Check for any updates for redeployment
	applyChange := false

	// Ensure image name is correct, update image if required

	newStopSchedule := azureVMScheduler.Spec.StopSchedule
	log.Info(newStopSchedule)

	newImage := azureVMScheduler.Spec.Image
	log.Info(newImage)

	newAzureVMName := azureVMScheduler.Spec.AzureVMName
	log.Info(newAzureVMName)

	newAzureRGName := azureVMScheduler.Spec.AzureRGName
	log.Info(newAzureRGName)

	var currentImage string = ""
	var currentStopSchedule string = ""
	// var currentAzureVMName string = ""
	// var currentAzureRGName string = ""

	// Check existing schedule
	if found.Spec.Schedule != "" {
		currentStopSchedule = found.Spec.Schedule
	}

	if newStopSchedule != currentStopSchedule {
		found.Spec.Schedule = newStopSchedule
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

	// // Check AzureVMName
	// if found.Spec.JobTemplate.Spec.Template.Spec.Containers != nil {
	// 	currentAzureVMName = found.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Env[0].Value
	// 	log.Info(currentAzureVMName)
	// }

	// if newAzureVMName != currentAzureVMName {
	// 	found.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Env[0].Value = newAzureVMName
	// 	applyChange = true
	// }

	// // Check AzureRGName
	// if found.Spec.JobTemplate.Spec.Template.Spec.Containers != nil {
	// 	currentAzureRGName = found.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Env[0].Value
	// 	log.Info(currentAzureRGName)
	// }

	// if newAzureRGName != currentAzureRGName {
	// 	found.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Env[0].Value = newAzureRGName
	// 	applyChange = true
	// }

	// log.Info(currentAzureVMName)
	// log.Info(currentAzureRGName)
	log.Info(currentImage)
	log.Info(currentStopSchedule)

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
	//	if !reflect.DeepEqual(currentAzureVMName, azureVMScheduler.Status.VMStopStatus) {
	//		azureVMScheduler.Status.VMStopStatus = currentAzureVMName
	//		// azureVMScheduler.Status.VMStopStatus = currentAzureVMName
	//		err := r.Client.Status().Update(ctx, azureVMScheduler)
	//		if err != nil {
	//			log.Error(err, "Failed to update azureVMScheduler status")
	//			return ctrl.Result{}, err
	//		}
	//	}
	return ctrl.Result{}, nil

}

// SetupWithManager sets up the controller with the Manager.
func (r *AzureVMSchedulerStopNHReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&azurev1.AzureVMSchedulerStopNH{}).
		Owns(&batchv1.CronJob{}).
		Complete(r)
}

func (r *AzureVMSchedulerStopNHReconciler) cronJobForAzureVMSchedulerStopNH(azureVMScheduler *azurev1.AzureVMSchedulerStopNH) *batchv1.CronJob {

	cron := &batchv1.CronJob{
		ObjectMeta: v1.ObjectMeta{
			Name:      azureVMScheduler.Name,
			Namespace: azureVMScheduler.Namespace,
			Labels:    AzureVMSchedulerStopNHLabels(azureVMScheduler, "azureVMScheduler"),
		},
		Spec: batchv1.CronJobSpec{
			Schedule: azureVMScheduler.Spec.StopSchedule,
			// TODO: Add Stop schedule
			//Schedule:  azureVMScheduler.Spec.StopSchedule,
			JobTemplate: batchv1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{{
								Name:  azureVMScheduler.Name,
								Image: azureVMScheduler.Spec.Image,
								Env: []corev1.EnvVar{
									{
										Name:  "command",
										Value: azureVMScheduler.Spec.Command,
									},
									{
										Name:  "AZURE_VM_NAME",
										Value: azureVMScheduler.Spec.AzureVMName,
									},
									{
										Name:  "AZURE_RG_NAME",
										Value: azureVMScheduler.Spec.AzureRGName,
									},

									{
										Name: "AZURE_TENANT_ID",
										ValueFrom: &corev1.EnvVarSource{
											SecretKeyRef: &corev1.SecretKeySelector{
												LocalObjectReference: corev1.LocalObjectReference{
													Name: "azure-secret",
												},
												Key: "azure-tenant-id",
											},
										},
									},
									{
										Name: "AZURE_CLIENT_ID",
										ValueFrom: &corev1.EnvVarSource{
											SecretKeyRef: &corev1.SecretKeySelector{
												LocalObjectReference: corev1.LocalObjectReference{
													Name: "azure-secret",
												},
												Key: "azure-client-id",
											},
										},
									},
									{
										Name: "AZURE_CLIENT_SECRET",
										ValueFrom: &corev1.EnvVarSource{
											SecretKeyRef: &corev1.SecretKeySelector{
												LocalObjectReference: corev1.LocalObjectReference{
													Name: "azure-secret",
												},
												Key: "azure-client-secret",
											},
										},
									},
									{
										Name: "AZURE_SUBSCRIPTION_ID",
										ValueFrom: &corev1.EnvVarSource{
											SecretKeyRef: &corev1.SecretKeySelector{
												LocalObjectReference: corev1.LocalObjectReference{
													Name: "azure-secret",
												},
												Key: "azure-subscription-id",
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
	ctrl.SetControllerReference(azureVMScheduler, cron, r.Scheme)
	return cron
}

func AzureVMSchedulerStopNHLabels(v *azurev1.AzureVMSchedulerStopNH, tier string) map[string]string {
	return map[string]string{
		"app":                       "AzureVMSchedulerStopNH",
		"AzureVMSchedulerStopNH_cr": v.Name,
		"tier":                      tier,
	}
}
