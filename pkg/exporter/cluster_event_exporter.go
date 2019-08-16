package exporter

import (
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	informerv1 "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"

	"github.com/pkg/errors"
	"github.com/ripta/spectacles/pkg/sinks"
)

type clusterEventExporter struct {
	eventsHaveSynced cache.InformerSynced
	sinkses          map[string]sinks.Writer
}

// NewClusterEventExporter creates a new event exporter with a sink.
func NewClusterEventExporter(eventInformer informerv1.EventInformer, w sinks.Writer) *clusterEventExporter {
	c := NewUnsunkClusterEventExporter(eventInformer)
	c.AddSink("default", w)
	return c
}

// NewUnsunkClusterEventExporter creates a new event exporter. Sinks must be
// added manually by calling AddSink.
func NewUnsunkClusterEventExporter(eventInformer informerv1.EventInformer) *clusterEventExporter {
	c := &clusterEventExporter{
		eventsHaveSynced: eventInformer.Informer().HasSynced,
		sinkses:          make(map[string]sinks.Writer),
	}

	klog.V(4).Info("installing event handlers")
	eventInformer.Informer().AddEventHandler(c.resourceEventHandlers())
	return c
}

// AddSink adds an event sink to the exporter.
func (c *clusterEventExporter) AddSink(name string, w sinks.Writer) {
	klog.V(4).Infof("adding sink %s", name)
	c.sinkses[name] = w
}

// DeleteSink removes an event sink by name from the exporter.
func (c *clusterEventExporter) DeleteSink(name string) {
	klog.V(4).Infof("deleting sink %s", name)
	delete(c.sinkses, name)
}

// Run waits for the initial cache sync to complete and blocks until the main
// exporter loop returns, i.e., when the stopCh closes.
func (c *clusterEventExporter) Run(stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()

	klog.Info("waiting for caches to populate")
	if ok := cache.WaitForCacheSync(stopCh, c.eventsHaveSynced); !ok {
		return fmt.Errorf("failed to sync caches")
	}

	klog.Info("initial cache sync complete")
	<-stopCh

	klog.Info("shutting down mainloop")
	return nil
}

func (c *clusterEventExporter) resourceEventHandlers() cache.ResourceEventHandlerFuncs {
	return cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			evt, ok := obj.(*apiv1.Event)
			if !ok {
				utilruntime.HandleError(fmt.Errorf("expecting an event, but received %#v", obj))
				return
			}

			for n, s := range c.sinkses {
				if err := s.Write(evt); err != nil {
					err = errors.Wrap(err, fmt.Sprintf("could not write event %s to sink %s", evt.GetName(), n))
					utilruntime.HandleError(err)
				}
			}
		},
	}
}
