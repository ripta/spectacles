package app

import (
	"time"

	"github.com/ripta/spectacles/pkg/sinks"
	"github.com/spf13/pflag"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Options encapsulates all application options.
type Options struct {
	Kubeconfig string
	Master     string

	HealthzPort int32
	MetricsPort int32

	ResyncPeriod *metav1.Duration

	Sink  sinks.Writer
	Sinks map[string]sinks.Writer
}

const (
	defaultHealthzPort = 8080
	defaultMetricsPort = 8081

	defaultResyncPeriod = 15 * time.Second
)

// NewOptions initializes a structure to encapsulate application options.
func NewOptions() *Options {
	return &Options{
		ResyncPeriod: &metav1.Duration{},
		Sinks:        make(map[string]sinks.Writer),
	}
}

// AddFlags injects application options into the flagset.
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Kubeconfig, "kubeconfig", "", "Path to kubeconfig file for out-of-cluster operations.")
	fs.StringVar(&o.Master, "master", "", "The address of the Kubernetes API server for out-of-cluster operations.")

	fs.Int32Var(&o.HealthzPort, "healthz-port", defaultHealthzPort, "The port to bind health check server. Use 0 to disable.")
	fs.Int32Var(&o.MetricsPort, "metrics-port", defaultMetricsPort, "The port to bind metrics server. Use 0 to disable.")

	fs.DurationVar(&o.ResyncPeriod.Duration, "resync-period", defaultResyncPeriod, "The interval of how often event watches are resynced, e.g., '30s', '5m'.")
}

// Complete checks command line argument combinations.
func (o *Options) Complete() error {
	return nil
}

// Validate checks that the correct options are specified.
func (o *Options) Validate(args []string) error {
	return nil
}
