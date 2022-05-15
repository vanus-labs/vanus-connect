package k8s

import (
	keda "github.com/kedacore/keda/v2/apis/keda/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateTriggerAuthentication(nameSpace, name string) *keda.TriggerAuthentication {

	triggerAuth := &keda.TriggerAuthentication{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: nameSpace,
		},
		Spec: keda.TriggerAuthenticationSpec{
			SecretTargetRef: []keda.AuthSecretTargetRef{},
		},
	}
	return triggerAuth
}
