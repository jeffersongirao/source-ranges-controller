package eventer

import (
	"github.com/jeffersongirao/source-ranges-controller/log"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
)

func NewEventRecorder(client kubernetes.Interface, logger log.Logger, component string) record.EventRecorder {
	broadcaster := record.NewBroadcaster()
	broadcaster.StartEventWatcher(
		func(event *corev1.Event) {
			if _, err := client.CoreV1().Events(event.Namespace).Create(event); err != nil {
				logger.Errorf("%v\n", err)
			}
		},
	)
	return broadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: component})
}
