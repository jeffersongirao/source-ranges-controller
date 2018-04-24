package controller

import (
	"fmt"

	"github.com/jeffersongirao/source-ranges-controller/eventer"
	"github.com/jeffersongirao/source-ranges-controller/log"
	"github.com/jeffersongirao/source-ranges-controller/service"
	"github.com/spotahome/kooper/operator/controller"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
)

type Controller struct {
	controller.Controller
	config Config
}

func New(config Config, k8sCli kubernetes.Interface, logger log.Logger) (*Controller, error) {
	recorder := eventer.NewEventRecorder(k8sCli, logger, "source-ranges-controller")
	sourceRangeEnforcer := service.NewConfigMapSourceRangeEnforcer(k8sCli, recorder)
	handler := &handler{sourceRangeEnforcerSrv: sourceRangeEnforcer}
	retriever := NewServiceRetriever(k8sCli, config.Namespace)
	ctrl := controller.NewSequential(config.ResyncPeriod, handler, retriever, nil, logger)

	return &Controller{
		Controller: ctrl,
		config:     config,
	}, nil
}

type handler struct {
	sourceRangeEnforcerSrv service.SourceRangeEnforcer
}

func (h *handler) Add(obj runtime.Object) error {
	svc, ok := obj.(*corev1.Service)
	if !ok {
		return fmt.Errorf("%v is not a service object", obj.GetObjectKind())
	}

	h.sourceRangeEnforcerSrv.EnforceSourceRangesToService(svc)
	return nil
}

func (h *handler) Delete(string) error {
	return nil
}
