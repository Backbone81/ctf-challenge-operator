package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	"github.com/Backbone81/ctf-challenge-operator/internal/utils"
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
	Use:          "ctf-operator",
	Short:        "This operator helps in running CTFs.",
	Long:         `This operator helps in running CTFs.`,
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

		// +kubebuilder:scaffold:builder

		if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
			return fmt.Errorf("setting up health check: %w", err)
		}
		if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
			return fmt.Errorf("setting up ready check: %w", err)
		}
		return mgr.Start(ctrl.SetupSignalHandler())
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
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
		"config file (default is $HOME/.ctf-operator.yaml)",
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
		"The address the metrics endpoint binds to. Use :8443 for HTTPS or :8080 for HTTP, or leave as 0 to disable the metrics service.",
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
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.",
	)
	rootCmd.PersistentFlags().StringVar(
		&leaderElectionNamespace,
		"leader-election-namespace",
		"ctf-operator",
		"The namespace in which leader election should happen.",
	)
	rootCmd.PersistentFlags().StringVar(
		&leaderElectionId,
		"leader-election-id",
		"ctf-operator",
		"The ID to use for leader election.",
	)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".ctf-operator" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".ctf-operator")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
