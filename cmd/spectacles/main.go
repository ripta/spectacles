package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/ripta/spectacles/pkg/app"
	"k8s.io/klog"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "last resort logger: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cmd := app.NewCommand()
	klog.InitFlags(nil)
	defer klog.Flush()
	return cmd.Execute()
}
