package k8s

import (
	kedahttp "github.com/kedacore/http-add-on/operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateHttpScaledObject(nameSpace, name, host string, port int32) *kedahttp.HTTPScaledObject {
	httpScaledObject := &kedahttp.HTTPScaledObject{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: nameSpace,
		},
		Spec: kedahttp.HTTPScaledObjectSpec{
			Host: host,
			ScaleTargetRef: &kedahttp.ScaleTargetRef{
				Deployment: name,
				Service:    name,
				Port:       port,
			},
		},
	}
	return httpScaledObject
}
