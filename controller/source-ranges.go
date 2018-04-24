package controller

import (
	"fmt"
	"net/http"

	"github.com/jeffersongirao/source-ranges-controller/eventer"
	"github.com/jeffersongirao/source-ranges-controller/log"
	"github.com/jeffersongirao/source-ranges-controller/service"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spotahome/kooper/monitoring/metrics"
	"github.com/spotahome/kooper/operator/controller"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
)

type Controller struct {
	controller.Controller
	config Config
}

const (
	metricsPrefix = "source_ranges"
	eventsPrefix  = "source-ranges-controller"
	metricsAddr   = ":7777"
)

func createPrometheusRecorder(logger log.Logger) metrics.Recorder {
	reg := prometheus.NewRegistry()
	m := metrics.NewPrometheus(metricsPrefix, reg)

	h := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	go func() {
		logger.Infof("serving metrics at %s", metricsAddr)
		http.ListenAndServe(metricsAddr, h)
	}()

	return m
}

func New(config Config, k8sCli kubernetes.Interface, logger log.Logger) (*Controller, error) {
	recorder := eventer.NewEventRecorder(k8sCli, logger, eventsPrefix)
	sourceRangeEnforcer := service.NewConfigMapSourceRangeEnforcer(k8sCli, recorder)
	handler := &handler{sourceRangeEnforcerSrv: sourceRangeEnforcer}
	retriever := NewServiceRetriever(k8sCli, config.Namespace)
	m := createPrometheusRecorder(logger)
	ctrl := controller.NewSequential(config.ResyncPeriod, handler, retriever, m, logger)

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
