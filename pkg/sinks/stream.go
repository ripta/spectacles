package sinks

import (
	"io"

	"github.com/pkg/errors"

	apiv1 "k8s.io/api/core/v1"
)

type StreamSink struct {
	Stream  io.Writer
	Encoder Encoder
}

func (s *StreamSink) Write(evt *apiv1.Event) error {
	p, err := s.Encoder(evt)
	if err != nil {
		return errors.Wrap(err, "encoding an event")
	}

	p = append(p, '\n')
	_, err = s.Stream.Write(p)
	return errors.Wrap(err, "writing to stream")
}
