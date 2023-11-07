/*


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

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/DoodleScheduling/gc-controller/internal/controllers"
	"github.com/fluxcd/pkg/runtime/client"
	helper "github.com/fluxcd/pkg/runtime/controller"
	"github.com/fluxcd/pkg/runtime/leaderelection"
	"github.com/fluxcd/pkg/runtime/logger"
	flag "github.com/spf13/pflag"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlcache "sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	// +kubebuilder:scaffold:imports
)

const controllerName = "gc-controller"

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)

	_ = corev1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

var (
	keep                    int
	maxAge                  time.Duration
	metricsAddr             string
	healthAddr              string
	concurrent              int
	gracefulShutdownTimeout time.Duration
	clientOptions           client.Options
	kubeConfigOpts          client.KubeConfigOptions
	logOptions              logger.Options
	leaderElectionOptions   leaderelection.Options
	rateLimiterOptions      helper.RateLimiterOptions
	watchOptions            helper.WatchOptions
)

func main() {
	flag.IntVar(&keep, "keep", 2,
		"The number of pods to keep for each workload.")
	flag.DurationVar(&maxAge, "max-age", time.Hour*24*7,
		"The number of pods to keep for each workload.")
	flag.StringVar(&metricsAddr, "metrics-addr", ":9556",
		"The address the metric endpoint binds to.")
	flag.StringVar(&healthAddr, "health-addr", ":9557",
		"The address the health endpoint binds to.")
	flag.IntVar(&concurrent, "concurrent", 4,
		"The number of concurrent Pod reconciles.")
	flag.DurationVar(&gracefulShutdownTimeout, "graceful-shutdown-timeout", 600*time.Second,
		"The duration given to the reconciler to finish before forcibly stopping.")

	clientOptions.BindFlags(flag.CommandLine)
	logOptions.BindFlags(flag.CommandLine)
	leaderElectionOptions.BindFlags(flag.CommandLine)
	rateLimiterOptions.BindFlags(flag.CommandLine)
	kubeConfigOpts.BindFlags(flag.CommandLine)
	watchOptions.BindFlags(flag.CommandLine)

	flag.Parse()
	logger.SetLogger(logger.NewLogger(logOptions))

	leaderElectionId := fmt.Sprintf("%s-%s", controllerName, "leader-election")
	if watchOptions.LabelSelector != "" {
		leaderElectionId = leaderelection.GenerateID(leaderElectionId, watchOptions.LabelSelector)
	}

	watchNamespace := ""
	if !watchOptions.AllNamespaces {
		watchNamespace = os.Getenv("RUNTIME_NAMESPACE")
	}

	opts := ctrl.Options{
		Scheme:                        scheme,
		MetricsBindAddress:            metricsAddr,
		HealthProbeBindAddress:        healthAddr,
		LeaderElection:                leaderElectionOptions.Enable,
		LeaderElectionReleaseOnCancel: leaderElectionOptions.ReleaseOnCancel,
		LeaseDuration:                 &leaderElectionOptions.LeaseDuration,
		RenewDeadline:                 &leaderElectionOptions.RenewDeadline,
		RetryPeriod:                   &leaderElectionOptions.RetryPeriod,
		GracefulShutdownTimeout:       &gracefulShutdownTimeout,
		Port:                          9443,
		LeaderElectionID:              leaderElectionId,
		Cache: ctrlcache.Options{
			Namespaces: []string{watchNamespace},
		},
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), opts)
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// Add liveness probe
	err = mgr.AddHealthzCheck("healthz", healthz.Ping)
	if err != nil {
		setupLog.Error(err, "Could not add liveness probe")
		os.Exit(1)
	}

	// Add readiness probe
	err = mgr.AddReadyzCheck("readyz", healthz.Ping)
	if err != nil {
		setupLog.Error(err, "Could not add readiness probe")
		os.Exit(1)
	}

	setReconciler := &controllers.PodReconciler{
		Keep:   keep,
		MaxAge: maxAge,
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Pod"),
		Scheme: mgr.GetScheme(),
	}

	if err = setReconciler.SetupWithManager(mgr, controllers.PodReconcilerOptions{
		MaxConcurrentReconciles: concurrent,
	}); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Pod")
		os.Exit(1)
	}

	// +kubebuilder:scaffold:builder
	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
