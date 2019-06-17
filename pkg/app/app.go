package app

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ripta/spectacles/pkg/sinks"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"

	"github.com/ripta/spectacles/pkg/exporter"

	"k8s.io/client-go/informers"
	kc "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
	utilflag "k8s.io/kubernetes/pkg/util/flag"
	"k8s.io/kubernetes/pkg/version/verflag"
)

func NewCommand() *cobra.Command {
	o := NewOptions()
	cmd := &cobra.Command{
		Use:  "spectacles",
		Long: `Export Kubernetes Events`,
		RunE: generateRunnerE(o),
	}

	o.AddFlags(cmd.Flags())
	cmd.MarkFlagFilename("config", "yaml", "json")
	return cmd
}

func generateRunnerE(o *Options) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		verflag.PrintAndExitIfRequested()
		utilflag.PrintFlags(cmd.Flags())

		if err := o.Complete(); err != nil {
			return errors.Wrap(err, "options were incomplete")
		}
		if err := o.Validate(args); err != nil {
			return errors.Wrap(err, "options did not validate")
		}

		klog.Infof("installing signal handlers")
		stopCh := setupSignalHandler()

		klog.Infof("setting up kubernetes client")
		cs, err := newKubeClientset(o.Master, o.Kubeconfig)
		if err != nil {
			return errors.Wrap(err, "creating new clientset to connect to Kubernetes")
		}

		go startMetricsServer(o.MetricsPort)

		fs := &sinks.StreamSink{
			Stream:  os.Stdout,
			Encoder: sinks.JSONEncoder,
		}
		inf := informers.NewSharedInformerFactory(cs, o.ResyncPeriod.Duration)

		klog.Info("booting up exporter")
		ex := exporter.NewClusterEventExporter(inf.Core().V1().Events(), fs)

		klog.Info("booting up informers")
		inf.Start(stopCh)

		klog.Info("starting main runloop")
		errCh := make(chan error)
		go func() {
			errCh <- ex.Run(stopCh)
		}()

		if err := <-errCh; err != nil {
			klog.Fatalf("exit: %v", err)
			return nil
		}

		return nil
	}
}

func newKubeClientset(master, kubeconfig string) (kc.Interface, error) {
	cfg, err := clientcmd.BuildConfigFromFlags(master, kubeconfig)
	if err != nil {
		return nil, errors.Wrap(err, "building configuration from options")
	}

	cs, err := kc.NewForConfig(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "initializing clientset")
	}

	return cs, nil
}

func setupSignalHandler() <-chan struct{} {
	stopCh := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-c
		klog.Infof("received signal %s; starting termination procedure", sig)
		close(stopCh)
		<-c
		os.Exit(0)
	}()
	return stopCh
}

func startMetricsServer(port int32) {
	http.Handle("/metrics", promhttp.Handler())

	portString := fmt.Sprintf(":%d", port)
	if err := http.ListenAndServe(portString, nil); err != nil {
		klog.Fatalf("metrics endpoint initialization failed: %v", err)
	}
}
