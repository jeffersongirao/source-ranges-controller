package controller

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

type ServiceRetriever struct {
	client kubernetes.Interface
}

func NewServiceRetriever(client kubernetes.Interface) *ServiceRetriever {
	return &ServiceRetriever{
		client: client,
	}
}

func (s *ServiceRetriever) GetListerWatcher() cache.ListerWatcher {
	return &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			return s.client.CoreV1().Services(metav1.NamespaceAll).List(options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return s.client.CoreV1().Services(metav1.NamespaceAll).Watch(options)
		},
	}
}

func (s *ServiceRetriever) GetObject() runtime.Object {
	return &corev1.Service{}
}
