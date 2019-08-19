package uflag

import (
	"github.com/spf13/pflag"
	"k8s.io/klog"
)

func PrintFlags(l klog.Verbose, fs *pflag.FlagSet) {
	fs.VisitAll(func(f *pflag.Flag) {
		l.Infof("FLAG: --%s=%q", f.Name, f.Value)
	})
}
