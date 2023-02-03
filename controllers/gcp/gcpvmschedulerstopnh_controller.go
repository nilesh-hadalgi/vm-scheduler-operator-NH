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

package gcp

import (
	"context"
	"strconv"

	"github.com/go-logr/logr"
	gcpv1 "github.com/nilesh-hadalgi/vm-scheduler-operator-NH/apis/gcp/v1"
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

// GCPVMSchedulerStopNHReconciler reconciles a GCPVMSchedulerStopNH object
type GCPVMSchedulerStopNHReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

//+kubebuilder:rbac:groups=gcp.xyzcompany.com,resources=gcpvmschedulerstopnhs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=gcp.xyzcompany.com,resources=gcpvmschedulerstopnhs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=gcp.xyzcompany.com,resources=gcpvmschedulerstopnhs/finalizers,verbs=update
// +kubebuilder:rbac:groups=batch,resources=cronjobs;jobs,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the GCPVMSchedulerStopNH object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *GCPVMSchedulerStopNHReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrllog.FromContext(ctx)

	log.Info("Reconciling GCPVMSchedulerStopNH")

	// Fetch the GCPVMScheduler CR
	//gcpVMScheduler, err := services.FetchGCPVMSchedulerCR(req.Name, req.Namespace)

	// Fetch the GCPVMScheduler instance
	gcpVMScheduler := &gcpv1.GCPVMSchedulerStopNH{}
	log.Info(req.NamespacedName.Name)

	err := r.Client.Get(ctx, req.NamespacedName, gcpVMScheduler)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("gcpVMScheduler resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get gcpVMScheduler.")
		return ctrl.Result{}, err
	}

	log.Info(gcpVMScheduler.Name)
	// Add const values for mandatory specs ( if left blank)
	// log.Info("Adding gcpVMScheduler mandatory specs")
	// utils.AddBackupMandatorySpecs(gcpVMScheduler)
	// Check if the CronJob already exists, if not create a new one

	found := &batchv1.CronJob{}
	err = r.Client.Get(ctx, types.NamespacedName{Name: gcpVMScheduler.Name, Namespace: gcpVMScheduler.Namespace}, found)
	//log.Info(*found.)
	if err != nil && errors.IsNotFound(err) {
		// Define a new CronJob
		cron := r.cronJobForGCPVMSchedulerStopNH(gcpVMScheduler)
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

	//if err := r.createResources(gcpVMScheduler, req); err != nil {
	//	log.Error(err, "Failed to create and update the secondary resource required for the GCPVMScheduler CR")
	//	return ctrl.Result{}, err
	//}

	// Ensure the cron inputs are the same as the spec.
	// TODO : Add support for instanceIds and stopSchedule

	// Check for any updates for redeployment
	applyChange := false

	// Ensure image name is correct, update image if required
	newInstance := gcpVMScheduler.Spec.Instance
	log.Info(newInstance)

	newProjectId := gcpVMScheduler.Spec.ProjectId
	log.Info(newProjectId)

	newStopSchedule := gcpVMScheduler.Spec.StopSchedule
	log.Info(newStopSchedule)

	newImage := gcpVMScheduler.Spec.Image
	log.Info(newImage)

	var currentImage string = ""
	// var currentProjectId string = ""
	var currentStopSchedule string = ""
	// var currentInstance string = ""

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

	/* // Check instance
	if found.Spec.JobTemplate.Spec.Template.Spec.Containers != nil {
		currentInstance = found.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Env[0].Value
		log.Info(currentInstance)
	}

	if newInstance != currentInstance {
		found.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Env[0].Value = newInstance
		applyChange = true
	}

	// Check ProjectIds
	if found.Spec.JobTemplate.Spec.Template.Spec.Containers != nil {
		currentProjectId = found.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Env[0].Value
		log.Info(currentProjectId)
	}

	if newProjectId != currentProjectId {
		found.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Env[0].Value = newProjectId
		applyChange = true
	}
	*/

	// log.Info(currentInstance)
	// log.Info(currentProjectId)
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

	// Update the GCPVMScheduler status
	// TODO: Define what needs to be added in status. Currently adding just instanceIds
	// if !reflect.DeepEqual(currentInstance, gcpVMScheduler.Status.VMStopStatus) {
	//	gcpVMScheduler.Status.VMStopStatus = currentInstance
	//	// gcpVMScheduler.Status.VMStopStatus = currentInstance
	//	err := r.Client.Status().Update(ctx, gcpVMScheduler)
	//	if err != nil {
	//		log.Error(err, "Failed to update gcpVMScheduler status")
	//		return ctrl.Result{}, err
	//	}
	//}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GCPVMSchedulerStopNHReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gcpv1.GCPVMSchedulerStopNH{}).
		Owns(&batchv1.CronJob{}).
		Complete(r)
}

func (r *GCPVMSchedulerStopNHReconciler) cronJobForGCPVMSchedulerStopNH(gcpVMScheduler *gcpv1.GCPVMSchedulerStopNH) *batchv1.CronJob {

	cron := &batchv1.CronJob{
		ObjectMeta: v1.ObjectMeta{
			Name:      gcpVMScheduler.Name,
			Namespace: gcpVMScheduler.Namespace,
			Labels:    GCPVMSchedulerStopNHLabels(gcpVMScheduler, "gcpVMScheduler"),
		},
		Spec: batchv1.CronJobSpec{
			Schedule: gcpVMScheduler.Spec.StopSchedule,
			// TODO: Add Stop schedule
			//Schedule:  gcpVMScheduler.Spec.StartSchedule,
			JobTemplate: batchv1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{{
								Name:  gcpVMScheduler.Name,
								Image: gcpVMScheduler.Spec.Image,
								VolumeMounts: []corev1.VolumeMount{
									{
										Name:      "foo",
										MountPath: "/etc/foo",
									},
								},
								Env: []corev1.EnvVar{
									{
										Name:  "gcp_instance",
										Value: gcpVMScheduler.Spec.Instance,
									},
									{
										Name:  "gcp_command",
										Value: gcpVMScheduler.Spec.Command,
									},
									{
										Name:  "gcp_projectID",
										Value: gcpVMScheduler.Spec.ProjectId,
									},
									{
										Name:  "gcp_zone",
										Value: gcpVMScheduler.Spec.Zone,
									},
									{
										Name:  "GOOGLE_APPLICATION_CREDENTIALS",
										Value: "/etc/foo/confi.json",
									}},
							}},
							Volumes: []corev1.Volume{
								{
									Name: "foo",
									VolumeSource: corev1.VolumeSource{
										Secret: &corev1.SecretVolumeSource{
											SecretName: "gcp-secret",
										},
									},
								},
							},
							RestartPolicy: "OnFailure",
						},
					},
				},
			},
		},
	}
	// Set gcpVMScheduler instance as the owner and controller
	ctrl.SetControllerReference(gcpVMScheduler, cron, r.Scheme)
	return cron
}

// // createResources will create and update the  resource which are required
// func (r *GCPVMSchedulerReconciler) createResources(gcpVMScheduler *gcpv1.GCPVMScheduler, request ctrl.Request) error {

// 	log := r.Log.WithValues("GCPVMScheduler", request.NamespacedName)
// 	log.Info("Creating   resources ...")

// 	// Check if the cronJob is created, if not create one
// 	if err := r.createCronJob(gcpVMScheduler); err != nil {
// 		log.Error(err, "Failed to create the CronJob")
// 		return err
// 	}

// 	return nil
// }

// // Check if the cronJob is created, if not create one
// func (r *GCPVMSchedulerReconciler) createCronJob(gcpVMScheduler *gcpv1.GCPVMScheduler) error {
// 	if _, err := services.FetchCronJob(gcpVMScheduler.Name, gcpVMScheduler.Namespace); err != nil {
// 		if err := r.client.Create(context.TODO(), resources.NewGCPVMSchedulerCronJob(gcpVMScheduler)); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

func GCPVMSchedulerStopNHLabels(v *gcpv1.GCPVMSchedulerStopNH, tier string) map[string]string {
	return map[string]string{
		"app":                     "GCPVMSchedulerStopNH",
		"GCPVMSchedulerStopNH_cr": v.Name,
		"tier":                    tier,
	}
}
