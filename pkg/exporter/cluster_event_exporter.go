package exporter

import (
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	informerv1 "k8s.io/client-go/informers/core/v1"
	listerv1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"

	"github.com/ripta/spectacles/pkg/sinks"
)

type clusterEventExporter struct {
	eventsLister     listerv1.EventLister
	eventsHaveSynced cache.InformerSynced
	sink             sinks.Writer
}

func NewClusterEventExporter(eventInformer informerv1.EventInformer, w sinks.Writer) *clusterEventExporter {
	c := &clusterEventExporter{
		eventsLister:     eventInformer.Lister(),
		eventsHaveSynced: eventInformer.Informer().HasSynced,
		sink:             w,
	}

	c.patchInformer(eventInformer.Informer())
	return c
}

func (c *clusterEventExporter) EmitToSink(evt *apiv1.Event) error {
	return c.sink.Write(evt)
}

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

func (c *clusterEventExporter) patchInformer(inf cache.SharedIndexInformer) {
	handlers := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			evt, ok := obj.(*apiv1.Event)
			if !ok {
				utilruntime.HandleError(fmt.Errorf("expecting an event, but received %#v", obj))
				return
			}

			if err := c.EmitToSink(evt); err != nil {
				utilruntime.HandleError(fmt.Errorf("could not write %q to sink: %v", evt.GetName(), err))
			}
		},
	}

	klog.Info("installing event handlers")
	inf.AddEventHandler(handlers)
}
