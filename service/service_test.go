package service_test

import (
	"errors"
	"testing"

	"github.com/jeffersongirao/source-ranges-controller/service"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	kubetesting "k8s.io/client-go/testing"
	"k8s.io/client-go/tools/record"
)

func TestEnforceSourceRangesToServiceWithNoRanges(t *testing.T) {
	k8sCli := fake.NewSimpleClientset()

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-config",
		},
		Data: map[string]string{
			"test": "123.123.123.123/32",
		},
	}
	k8sCli.CoreV1().ConfigMaps(metav1.NamespaceDefault).Create(cm)

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: metav1.NamespaceDefault,
			Annotations: map[string]string{
				"source-ranges.alpha.girao.net/config-map": "test-config",
			},
		},
	}
	k8sCli.CoreV1().Services(svc.ObjectMeta.Namespace).Create(svc)

	recorder := record.NewFakeRecorder(1)
	e := service.NewConfigMapSourceRangeEnforcer(k8sCli, recorder)

	e.EnforceSourceRangesToService(svc)

	new, _ := k8sCli.CoreV1().Services(svc.ObjectMeta.Namespace).Get(svc.ObjectMeta.Name, metav1.GetOptions{})
	assert.Equal(t, []string{"123.123.123.123/32"}, new.Spec.LoadBalancerSourceRanges)

	events := collectEvents(recorder.Events)
	if eventCount := len(events); eventCount != 1 {
		t.Errorf("Expected 1 event when service load balancer source ranges are updated but got %d", eventCount)
	}
}

func TestEnforceSourceRangesToServiceWithLessRanges(t *testing.T) {
	k8sCli := fake.NewSimpleClientset()

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-config",
		},
		Data: map[string]string{
			"test":  "123.123.123.123/32",
			"test2": "123.123.123.124/32",
		},
	}
	k8sCli.CoreV1().ConfigMaps(metav1.NamespaceDefault).Create(cm)

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: metav1.NamespaceDefault,
			Annotations: map[string]string{
				"source-ranges.alpha.girao.net/config-map": "test-config",
			},
		},
		Spec: corev1.ServiceSpec{
			LoadBalancerSourceRanges: []string{"123.123.123.123/32"},
		},
	}
	k8sCli.CoreV1().Services(svc.ObjectMeta.Namespace).Create(svc)

	recorder := record.NewFakeRecorder(1)
	e := service.NewConfigMapSourceRangeEnforcer(k8sCli, recorder)

	e.EnforceSourceRangesToService(svc)

	new, _ := k8sCli.CoreV1().Services(svc.ObjectMeta.Namespace).Get(svc.ObjectMeta.Name, metav1.GetOptions{})
	assert.ElementsMatch(t, []string{"123.123.123.123/32", "123.123.123.124/32"}, new.Spec.LoadBalancerSourceRanges)

	events := collectEvents(recorder.Events)
	if eventCount := len(events); eventCount != 1 {
		t.Errorf("Expected 1 event when service load balancer source ranges are updated but got %d", eventCount)
	}
}

func TestEnforceSourceRangesToServiceWithMoreRanges(t *testing.T) {
	k8sCli := fake.NewSimpleClientset()

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-config",
		},
		Data: map[string]string{
			"test": "123.123.123.123/32",
		},
	}
	k8sCli.CoreV1().ConfigMaps(metav1.NamespaceDefault).Create(cm)

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: metav1.NamespaceDefault,
			Annotations: map[string]string{
				"source-ranges.alpha.girao.net/config-map": "test-config",
			},
		},
		Spec: corev1.ServiceSpec{
			LoadBalancerSourceRanges: []string{"123.123.123.123/32", "123.123.123.124/32"},
		},
	}
	k8sCli.CoreV1().Services(svc.ObjectMeta.Namespace).Create(svc)

	recorder := record.NewFakeRecorder(1)
	e := service.NewConfigMapSourceRangeEnforcer(k8sCli, recorder)

	e.EnforceSourceRangesToService(svc)

	new, _ := k8sCli.CoreV1().Services(svc.ObjectMeta.Namespace).Get(svc.ObjectMeta.Name, metav1.GetOptions{})
	assert.ElementsMatch(t, []string{"123.123.123.123/32"}, new.Spec.LoadBalancerSourceRanges)

	events := collectEvents(recorder.Events)
	if eventCount := len(events); eventCount != 1 {
		t.Errorf("Expected 1 event when service load balancer source ranges are updated but got %d", eventCount)
	}
}

func TestEnforceSourceRangesToServiceWithNoChanges(t *testing.T) {
	k8sCli := fake.NewSimpleClientset()

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-config",
		},
		Data: map[string]string{
			"test": "123.123.123.123/32",
		},
	}
	k8sCli.CoreV1().ConfigMaps(metav1.NamespaceDefault).Create(cm)

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: metav1.NamespaceDefault,
			Annotations: map[string]string{
				"source-ranges.alpha.girao.net/config-map": "test-config",
			},
		},
		Spec: corev1.ServiceSpec{
			LoadBalancerSourceRanges: []string{"123.123.123.123/32"},
		},
	}

	k8sCli.CoreV1().Services(svc.ObjectMeta.Namespace).Create(svc)

	recorder := record.NewFakeRecorder(1)
	e := service.NewConfigMapSourceRangeEnforcer(k8sCli, recorder)

	e.EnforceSourceRangesToService(svc)

	new, _ := k8sCli.CoreV1().Services(svc.ObjectMeta.Namespace).Get(svc.ObjectMeta.Name, metav1.GetOptions{})
	assert.Equal(t, []string{"123.123.123.123/32"}, new.Spec.LoadBalancerSourceRanges)

	events := collectEvents(recorder.Events)
	if eventCount := len(events); eventCount != 0 {
		t.Errorf("Expected 0 event when service load balancer source ranges are up to date but got %d", eventCount)
	}
}

func TestEnforceSourceRangesToServiceWithNonExistingConfigMap(t *testing.T) {
	k8sCli := fake.NewSimpleClientset()

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: metav1.NamespaceDefault,
			Annotations: map[string]string{
				"source-ranges.alpha.girao.net/config-map": "test-config",
			},
		},
	}
	k8sCli.CoreV1().Services(svc.ObjectMeta.Namespace).Create(svc)

	recorder := record.NewFakeRecorder(1)
	e := service.NewConfigMapSourceRangeEnforcer(k8sCli, recorder)

	err := e.EnforceSourceRangesToService(svc)
	assert.NotNil(t, err)

	new, _ := k8sCli.CoreV1().Services(svc.ObjectMeta.Namespace).Get(svc.ObjectMeta.Name, metav1.GetOptions{})
	assert.Nil(t, new.Spec.LoadBalancerSourceRanges)

	events := collectEvents(recorder.Events)
	if eventCount := len(events); eventCount != 1 {
		t.Errorf("Expected 1 event when can't find ConfigMap specified by annotation but got %d", eventCount)
	}

	assert.Equal(t, "Warning SourceRangesEnforcementFailed could not read ConfigMap test-config: configmaps \"test-config\" not found", events[0])
}

func TestEnforceSourceRangesToServiceWhenError(t *testing.T) {
	k8sCli := &fake.Clientset{}

	k8sCli.AddReactor("update", "services", func(action kubetesting.Action) (bool, runtime.Object, error) {
		return true, nil, apierrors.NewInternalError(errors.New("API server down"))
	})

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-config",
		},
		Data: map[string]string{
			"test": "123.123.123.123/32",
		},
	}
	k8sCli.AddReactor("get", "configmaps", func(action kubetesting.Action) (bool, runtime.Object, error) {
		return true, cm, nil
	})

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: metav1.NamespaceDefault,
			Name:      "test-service",
			Annotations: map[string]string{
				"source-ranges.alpha.girao.net/config-map": "test-config",
			},
		},
	}
	k8sCli.CoreV1().Services(svc.ObjectMeta.Namespace).Create(svc)

	recorder := record.NewFakeRecorder(1)
	e := service.NewConfigMapSourceRangeEnforcer(k8sCli, recorder)

	err := e.EnforceSourceRangesToService(svc)
	assert.NotNil(t, err)

	new, _ := k8sCli.CoreV1().Services(svc.ObjectMeta.Namespace).Get(svc.ObjectMeta.Name, metav1.GetOptions{})
	assert.Nil(t, new.Spec.LoadBalancerSourceRanges)

	events := collectEvents(recorder.Events)
	if eventCount := len(events); eventCount != 1 {
		t.Errorf("Expected 1 event when unable to update Service but got %d", eventCount)
		return
	}

	assert.Equal(t, "Warning SourceRangesEnforcementFailed could not update Service test-service: Internal error occurred: API server down", events[0])
}

func TestEnforceSourceRangesToServiceWithoutAnnotation(t *testing.T) {
	k8sCli := fake.NewSimpleClientset()

	svc := corev1.Service{
		Spec: corev1.ServiceSpec{
			LoadBalancerSourceRanges: []string{"123.123.123.122/32"},
		},
	}
	k8sCli.CoreV1().Services(svc.ObjectMeta.Namespace).Create(&svc)

	recorder := record.NewFakeRecorder(1)
	e := service.NewConfigMapSourceRangeEnforcer(k8sCli, recorder)

	e.EnforceSourceRangesToService(&svc)

	new, _ := k8sCli.CoreV1().Services(svc.ObjectMeta.Namespace).Get(svc.ObjectMeta.Name, metav1.GetOptions{})
	assert.Equal(t, []string{"123.123.123.122/32"}, new.Spec.LoadBalancerSourceRanges)
}

func collectEvents(source <-chan string) []string {
	done := false
	events := make([]string, 0)
	for !done {
		select {
		case event := <-source:
			events = append(events, event)
		default:
			done = true
		}
	}
	return events
}
