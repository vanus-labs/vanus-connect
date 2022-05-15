package k8s

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func CreateLBService(nameSpace, name string, port int32) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: nameSpace,
			Name:      name,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{{
				//Port:       connector.Spec.Scaler.Metadata["svcPort"].IntVal,
				Port:       port,
				TargetPort: intstr.IntOrString{IntVal: 8080},
			},
			},
			Selector: map[string]string{
				"app": name,
			},
			Type: corev1.ServiceTypeLoadBalancer,
		},
	}
	return service
}
