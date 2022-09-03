package k8s

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// createBaseService creates a K8S service without setting its service type
func createBaseService(nameSpace, name string, port int32) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: nameSpace,
			Name:      name,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{{
				Port: port,
			},
			},
			Selector: map[string]string{
				"app": name,
			},
		},
	}
	return service
}

// CreateCLUService creates a K8S Cluster Service
func CreateCLUService(nameSpace, name string, port int32) *corev1.Service {
	service := createBaseService(nameSpace, name, port)
	service.Spec.Type = corev1.ServiceTypeClusterIP
	return service
}

// CreateLBService creates a K8S LoadBalance Service
func CreateLBService(nameSpace, name string, port int32) *corev1.Service {
	service := createBaseService(nameSpace, name, port)
	service.Spec.Type = corev1.ServiceTypeLoadBalancer
	return service
}
