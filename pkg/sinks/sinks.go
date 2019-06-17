package sinks

import (
	apiv1 "k8s.io/api/core/v1"
)

// Writer wraps a sink's Write method.
type Writer interface {
	Write(*apiv1.Event) error
}
