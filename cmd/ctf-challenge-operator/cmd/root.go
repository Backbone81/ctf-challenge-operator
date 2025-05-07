package cmd

import (
	"fmt"
	"os"

	"github.com/backbone81/ctf-challenge-operator/internal/controller"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	"github.com/backbone81/ctf-challenge-operator/internal/utils"
)

var (
	cfgFile             string
	enableDeveloperMode bool
	logLevel            int

	metricsBindAddress      string
	healthProbeBindAddress  string
	leaderElectionEnabled   bool
	leaderElectionNamespace string
	leaderElectionId        string
)

var rootCmd = &cobra.Command{
	Use:          "ctf-challenge-operator",
	Short:        "This operator manages CTF challenge instances.",
	Long:         `This operator manages CTF challenge instances.`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger, zapLogger, err := utils.CreateLogger(logLevel, enableDeveloperMode)
		if err != nil {
			return fmt.Errorf("setting up logger: %w", err)
		}
		defer zapLogger.Sync() //nolint:errcheck // This is the logger we are flushing, no way to log the error here.

		if enableDeveloperMode {
			logger.Info("WARNING: Developer mode is enabled. This must not be used in production!")
		}

		ctrl.SetLogger(logger)
		restConfig, err := ctrl.GetConfig()
		if err != nil {
			return fmt.Errorf("setting up kubernetes config: %w", err)
		}
		mgr, err := ctrl.NewManager(
			restConfig,
			ctrl.Options{
				LeaderElection:          leaderElectionEnabled,
				LeaderElectionNamespace: leaderElectionNamespace,
				LeaderElectionID:        leaderElectionId,
				Metrics: metricsserver.Options{
					BindAddress: metricsBindAddress,
				},
				HealthProbeBindAddress: healthProbeBindAddress,
			},
		)
		if err != nil {
			return fmt.Errorf("setting up manager: %w", err)
		}

		reconciler := controller.NewReconciler(
			utils.NewLoggingClient(mgr.GetClient(), logger),
			controller.WithDefaultReconcilers(mgr.GetEventRecorderFor("challenge-instance")),
		)
		if err := reconciler.SetupWithManager(mgr); err != nil {
			return fmt.Errorf("setting up reconciler with manager: %w", err)
		}

		if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
			return fmt.Errorf("setting up health check: %w", err)
		}
		if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
			return fmt.Errorf("setting up ready check: %w", err)
		}
		return mgr.Start(ctrl.SetupSignalHandler())
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(
		&cfgFile,
		"config",
		"",
		"config file (default is $HOME/.ctf-challenge-operator.yaml)",
	)
	rootCmd.PersistentFlags().BoolVar(
		&enableDeveloperMode,
		"enable-developer-mode",
		false,
		"This option makes the log output friendlier to humans.",
	)
	rootCmd.PersistentFlags().IntVar(
		&logLevel,
		"log-level",
		0,
		"How verbose the logs are. Level 0 will show info, warning and error. Level 1 and up will show increasing details.",
	)

	rootCmd.PersistentFlags().StringVar(
		&metricsBindAddress,
		"metrics-bind-address",
		"0",
		"The address the metrics endpoint binds to. Use :8443 for HTTPS or :8080 for HTTP, or leave as 0 "+
			"to disable the metrics service.",
	)
	rootCmd.PersistentFlags().StringVar(
		&healthProbeBindAddress,
		"health-probe-bind-address",
		":8081",
		"The address the probe endpoint binds to.",
	)
	rootCmd.PersistentFlags().BoolVar(
		&leaderElectionEnabled,
		"leader-election-enabled",
		false,
		"Enable leader election for controller manager. Enabling this will ensure there is only one "+
			"active controller manager.",
	)
	rootCmd.PersistentFlags().StringVar(
		&leaderElectionNamespace,
		"leader-election-namespace",
		"ctf-challenge-operator",
		"The namespace in which leader election should happen.",
	)
	rootCmd.PersistentFlags().StringVar(
		&leaderElectionId,
		"leader-election-id",
		"ctf-challenge-operator",
		"The ID to use for leader election.",
	)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".ctf-challenge-operator")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
