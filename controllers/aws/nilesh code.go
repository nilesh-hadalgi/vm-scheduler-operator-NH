package aws
import (
    "context"
    "reflect"
    "strconv"
    awsv1 "github.com/ManojDhanorkar/vm-scheduler-operator/apis/aws/v1"
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

type AWSVMSchedulerReconciler struct {
    Client client.Client
    Scheme *runtime.Scheme
    Log    logr.Logger
}

func (r *AWSVMSchedulerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    log := ctrllog.FromContext(ctx)
	log.Info("Reconciling AWSVMScheduler")
	awsVMScheduler := &awsv1.AWSVMScheduler{}
	log.Info(req.NamespacedName.Name)
	err := r.Client.Get(ctx,req.NamespacedName.Name, awsVMScheduler)
	if err != nil{
		if errors.IsNotFound(err){
			log.Info("awsVMScheduler resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{},nil
        
		}
		log.Error(err, "Failed to get awsVMScheduler.")
        return ctrl.Result{},err
	}

	log.Info(awsVMScheduler.Name)

	found := &batchv1.CronJob{}
	err = r.Client.Get(ctx,types.NamespacedName{Name:awsVMScheduler.name,Namespace:awsVMScheduler.Namespace},found)
	if err !=nil && errors.IsNotFound(err){
		corn := r.cronJobForAWSVMScheduler(awsVMScheduler)
		log.Info("Creating a new CronJob", "CronJob.Namespace", cron.Namespace, "CronJob.Name", cron.Name)
        err = r.Client.create(ctx,corn)
		if err !=nil {
			log.Error(err,"Failed to create new CronJob", "CronJob.Namespace", cron.Namespace, "CronJob.Name", cron.Name)
			return ctrl.Result{},err
		}
		return ctrl.Result{Requeue: true}, nil
	} else if err!= nil {
		log.Error(err, "Failed to get Cronjob")
        return ctrl.Result{}, err
	}

	applyChange:= false
}