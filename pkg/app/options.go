package app

import (
	"time"

	"github.com/spf13/pflag"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Options struct {
	Kubeconfig string
	Master     string

	HealthzPort int32
	MetricsPort int32

	ResyncPeriod *metav1.Duration
}

const (
	defaultHealthzPort = 8080
	defaultMetricsPort = 8081

	defaultResyncPeriod = 15 * time.Second
)

func NewOptions() *Options {
	return &Options{
		ResyncPeriod: &metav1.Duration{},
	}
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Kubeconfig, "kubeconfig", "", "Path to kubeconfig file for out-of-cluster operations.")
	fs.StringVar(&o.Master, "master", "", "The address of the Kubernetes API server for out-of-cluster operations.")

	fs.Int32Var(&o.HealthzPort, "healthz-port", defaultHealthzPort, "The port to bind health check server. Use 0 to disable.")
	fs.Int32Var(&o.MetricsPort, "metrics-port", defaultMetricsPort, "The port to bind metrics server. Use 0 to disable.")

	fs.DurationVar(&o.ResyncPeriod.Duration, "resync-period", defaultResyncPeriod, "The interval of how often event watches are resynced, e.g., '30s', '5m'.")
}

func (o *Options) Complete() error {
	return nil
}

func (o *Options) Validate(args []string) error {
	return nil
}
