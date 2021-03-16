package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	appsv1 "k8s.io/api/apps/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

func init() {
	log.SetLogger(zap.New())
}

func main() {
	entryLog := log.Log.WithName("entrypoint")
	loadConfiguration()
	backupRegistry := viper.GetString("backupRegistry")
	dockerUser := viper.GetString("dockerUser")
	dockerToken := viper.GetString("dockerToken")

	entryLog.Info("Backup registry configured. Docker user: " + dockerUser + " backup registry " + backupRegistry)

	// Setup a Manager
	entryLog.Info("setting up manager")
	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{})
	if err != nil {
		entryLog.Error(err, "unable to set up overall controller manager")
		os.Exit(1)
	}

	// Setup a new controller to reconcile Deployments
	entryLog.Info("Setting up controller")
	r := NewRegistry(backupRegistry, dockerUser, dockerToken)
	cDeployment, err := controller.New("backup-controller-deployment", mgr, controller.Options{
		Reconciler: &reconcileBackup{client: mgr.GetClient(), registry: r, k8sType: deploymentType},
	})
	if err != nil {
		entryLog.Error(err, "unable to set up individual controller")
		os.Exit(1)
	}

	cDaemonset, err := controller.New("backup-controller-daemonset", mgr, controller.Options{
		Reconciler: &reconcileBackup{client: mgr.GetClient(), registry: r, k8sType: daemonsetType},
	})
	if err != nil {
		entryLog.Error(err, "unable to set up individual controller")
		os.Exit(1)
	}
	// Watch Deployment and enqueue Deployment object key
	if err := cDeployment.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForObject{}); err != nil {
		entryLog.Error(err, "unable to watch Deployment")
		os.Exit(1)
	}

	// Watch DaemonSet and enqueue DaemonSet object key
	if err := cDaemonset.Watch(&source.Kind{Type: &appsv1.DaemonSet{}}, &handler.EnqueueRequestForObject{}); err != nil {
		entryLog.Error(err, "unable to watch DaemonSet")
		os.Exit(1)
	}

	entryLog.Info("starting manager")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		entryLog.Error(err, "unable to run manager")
		os.Exit(1)
	}
}

func loadConfiguration() {
	viper.SetConfigFile("./config/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	viper.SetConfigFile("./secure-config/config.yaml")
	viper.MergeInConfig()
}
