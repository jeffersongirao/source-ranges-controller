package service

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"
)

// SourceRangeEnforcer enforces loadBalancerSourceRanges
type SourceRangeEnforcer interface {
	EnforceSourceRangesToService(svc *corev1.Service) error
}

// ConfigMapSourceRangeEnforcer enforces that loadBalancerSourceRanges to a Service
// from a ConfigMap specified by annotation
type ConfigMapSourceRangeEnforcer struct {
	client   kubernetes.Interface
	recorder record.EventRecorder
}

// EnforceSourceRangesToService enforces loadBalancerSourceRanges to a Service based on ConfigMap from annotation
func (c *ConfigMapSourceRangeEnforcer) EnforceSourceRangesToService(svc *corev1.Service) error {
	options := metav1.GetOptions{}
	cmName := svc.ObjectMeta.Annotations["net.girao.source-ranges-controller/source-ranges-config-map"]

	if cmName != "" {
		cm, err := c.client.CoreV1().ConfigMaps(svc.ObjectMeta.Namespace).Get(cmName, options)
		if err != nil {
			reason := "SourceRangesEnforcementFailed"
			message := fmt.Sprintf("could not read ConfigMap %s: %v", cmName, err)
			c.recorder.Eventf(svc, corev1.EventTypeWarning, reason, message)
			return nil
		}

		if len(difference(configMapValues(cm.Data), svc.Spec.LoadBalancerSourceRanges)) != 0 {
			svc.Spec.LoadBalancerSourceRanges = configMapValues(cm.Data)
			_, err = c.client.CoreV1().Services(svc.ObjectMeta.Namespace).Update(svc)
			if err != nil {
				reason := "SourceRangesEnforcementFailed"
				message := fmt.Sprintf("could not update Service %s: %v", svc.ObjectMeta.Name, err)
				c.recorder.Eventf(svc, corev1.EventTypeWarning, reason, message)
			} else {
				reason := "SourceRangesEnforcementSuccessful"
				message := fmt.Sprintf("Updated Service %s with LB source ranges: %v", svc.ObjectMeta.Name, configMapValues(cm.Data))
				c.recorder.Eventf(svc, corev1.EventTypeNormal, reason, message)
			}
		}
	}
	return nil
}

// NewConfigMapSourceRangeEnforcer returns a new ConfigMapSourceRangeEnforcer
func NewConfigMapSourceRangeEnforcer(k8sCli kubernetes.Interface, recorder record.EventRecorder) SourceRangeEnforcer {
	return &ConfigMapSourceRangeEnforcer{
		client:   k8sCli,
		recorder: recorder,
	}
}

func configMapValues(data map[string]string) []string {
	values := make([]string, 0, len(data))
	for value := range data {
		values = append(values, data[value])
	}
	return values
}

func difference(slice1 []string, slice2 []string) []string {
	var diff []string

	// Loop two times, first to find slice1 strings not in slice2,
	// second loop to find slice2 strings not in slice1
	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			// String not found. We add it to return slice
			if !found {
				diff = append(diff, s1)
			}
		}
		// Swap the slices, only if it was the first loop
		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}

	return diff
}
