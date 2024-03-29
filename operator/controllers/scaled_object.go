package controllers

import (
	"context"
	"errors"
	keda "github.com/kedacore/keda/v2/apis/keda/v1alpha1"
	vance "github.com/linkall-labs/vance/operator/api/v1alpha1"
	k8s2 "github.com/linkall-labs/vance/operator/pkg/k8s"
	"github.com/linkall-labs/vance/operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func createOrUpdateScaledObject(
	ctx context.Context,
	r *ConnectorReconciler,
	connector *vance.Connector) error {
	soName := resourceName(connector.Name, vance.SoResType)
	logger := r.Log.WithValues(
		"ScaledObject name", soName,
		"ScaledObject namespace", connector.Namespace)
	logger.Info(" Start creating or updating a ScaledObject ")
	so := k8s2.CreateScaledObject(
		connector.Namespace, soName, connector.Name,
		connector.Spec.ScalingRule.CustomScaling.Triggers,
	)

	if connector.Spec.ScalingRule.CustomScaling.CheckInterval != nil {
		so.Spec.PollingInterval = connector.Spec.ScalingRule.CustomScaling.CheckInterval
	}
	if connector.Spec.ScalingRule.CustomScaling.CooldownPeriod != nil {
		so.Spec.CooldownPeriod = connector.Spec.ScalingRule.CustomScaling.CooldownPeriod
	}
	if connector.Spec.ScalingRule.MaxReplicaCount != nil {
		so.Spec.MaxReplicaCount = connector.Spec.ScalingRule.MaxReplicaCount
	}
	if connector.Spec.ScalingRule.MinReplicaCount != nil {
		so.Spec.MinReplicaCount = connector.Spec.ScalingRule.MinReplicaCount
	}
	userSecret := &corev1.Secret{}
	for _, trigger := range connector.Spec.ScalingRule.CustomScaling.Triggers {
		if _, existedKey := ScalerConf[trigger.Type+"-auth"]; !existedKey {
			err := errors.New("TriggerType not found error")
			logger.Error(err, "trigger type "+trigger.Type+" not supported")
			return err
		}
		var soTrigger = keda.ScaleTriggers{}
		soTrigger.Type = trigger.Type
		soTrigger.Metadata = trigger.Metadata
		if trigger.SecretRef != "" {
			existedKey := client.ObjectKey{
				Namespace: connector.Namespace,
				Name:      trigger.SecretRef,
			}
			var desiredAuth []string
			if err := r.Get(ctx, existedKey, userSecret); err != nil {
				if k8serrors.IsNotFound(err) {
					logger.Error(err, "no such secret ", "trigger.SecretRef", trigger.SecretRef,
						"namespace", connector.Namespace)
				} else {
					logger.Error(err, "fetch secret err")
				}
				return err
			} else {
				if ok, v := util.IsValidSecret(util.WrapSBM(userSecret.Data), ScalerConf[string(trigger.Type)+"-auth"]); !ok {
					err := errors.New("secret " + trigger.SecretRef + " misses required field")
					logger.Error(err, "secret misses required field")
					return err
				} else {
					desiredAuth = v
				}
			}
			// create a TriggerAuthentication if provided secret is valid
			taName := resourceName(connector.Name, vance.TAResType)
			ta := k8s2.CreateTriggerAuthentication(connector.Namespace, taName)
			logger.Info("scalerConfData", "map len", len(desiredAuth))
			for i, v := range desiredAuth {
				logger.Info("scalerConfMap", "value", v)
				ta.Spec.SecretTargetRef = append(ta.Spec.SecretTargetRef, keda.AuthSecretTargetRef{
					Parameter: v,
					Name:      trigger.SecretRef,
					Key:       v,
				})
				logger.Info("ta", "slice", ta.Spec.SecretTargetRef[i])
			}
			logger.Info("TA SecretTargetRef", "len of SecretTargetRef", len(ta.Spec.SecretTargetRef))
			logger = r.Log.WithValues("TriggerAuthentication name", taName,
				"TriggerAuthentication namespace", connector.Namespace)
			logger.Info("create a TriggerAuthentication", "ta", ta)
			if err := controllerutil.SetControllerReference(connector, ta, r.Scheme); err != nil {
				logger.Error(err, "Set TriggerAuthentication ControllerReference error")
				return err
			}
			if err := createOrPatchObj(ctx, r, ta, taName,
				connector.Namespace, logger, vance.TAResType); err != nil {
				return err
			}
			so.Spec.Triggers = append(so.Spec.Triggers, keda.ScaleTriggers{
				Type:     trigger.Type,
				Metadata: trigger.Metadata,
				AuthenticationRef: &keda.ScaledObjectAuthRef{
					Name: taName,
				},
			})

			if err := controllerutil.SetControllerReference(connector, so, r.Scheme); err != nil {
				logger.Error(err, "Set ScaledObject ControllerReference error")
				return err
			}
			if err := createOrPatchObj(ctx, r, so, soName,
				connector.Namespace, logger, vance.SoResType); err != nil {
				return err
			}
		}
	}

	return nil
}
