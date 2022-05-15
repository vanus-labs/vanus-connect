package controllers

import (
	"context"
	"errors"
	keda "github.com/kedacore/keda/v2/apis/keda/v1alpha1"
	vance "github.com/linkall-labs/vance/api/v1alpha1"
	"github.com/linkall-labs/vance/pkg/k8s"
	"github.com/linkall-labs/vance/pkg/util"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"strconv"
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
	so := k8s.CreateScaledObject(
		connector.Namespace, soName, connector.Name,
		string(connector.Spec.ConnectorType),
	)

	if connector.Spec.Scaler.CheckInterval != nil {
		so.Spec.PollingInterval = connector.Spec.Scaler.CheckInterval
	}
	if connector.Spec.Scaler.CooldownPeriod != nil {
		so.Spec.CooldownPeriod = connector.Spec.Scaler.CooldownPeriod
	}
	if connector.Spec.Scaler.MaxReplicaCount != nil {
		so.Spec.MaxReplicaCount = connector.Spec.Scaler.MaxReplicaCount
	}
	if connector.Spec.Scaler.MinReplicaCount != nil {
		so.Spec.MinReplicaCount = connector.Spec.Scaler.MinReplicaCount
	}
	userSecret := &corev1.Secret{}
	secretName := connector.Spec.Scaler.Metadata["secret"].StrVal
	if secretName != "" {
		existedKey := client.ObjectKey{
			Namespace: connector.Namespace,
			Name:      secretName,
		}
		var desiredAuth []string
		if err := r.Get(ctx, existedKey, userSecret); err != nil {
			if k8serrors.IsNotFound(err) {
				logger.Error(err, "no such secret ", secretName,
					"namespace", connector.Namespace)
			} else {
				logger.Error(err, "fetch secret err")
			}
			return err
		} else {
			if ok, v := util.IsValidSecret(util.WrapSBM(userSecret.Data), ScalerConf[string(connector.Spec.ConnectorType)+"-auth"]); !ok {
				err := errors.New("secret " + secretName + " misses required field")
				logger.Error(err, "secret misses required field", "missing field", v)
				return err
			} else {
				desiredAuth = v
			}
		}
		// create a TriggerAuthentication if provided secret is valid
		taName := resourceName(connector.Name, vance.TAResType)
		ta := k8s.CreateTriggerAuthentication(connector.Namespace, taName)
		logger.Info("scalerConfData", "map len", len(desiredAuth))
		for i, v := range desiredAuth {
			logger.Info("scalerConfMap", "value", v)
			ta.Spec.SecretTargetRef = append(ta.Spec.SecretTargetRef, keda.AuthSecretTargetRef{
				Parameter: v,
				Name:      secretName,
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
		so.Spec.Triggers[0].Metadata = make(map[string]string)
		for k, v := range connector.Spec.Scaler.Metadata {
			if k == "secret" {
				continue
			}
			if v.StrVal != "" {
				so.Spec.Triggers[0].Metadata[k] = v.StrVal
			} else {
				so.Spec.Triggers[0].Metadata[k] = strconv.FormatInt(int64(v.IntVal), 10)
			}
		}
		logger = r.Log.WithValues(
			"ScaledObject name", soName,
			"ScaledObject namespace", connector.Namespace)
		logger.Info("ScaledObject triggers", "so metadata", so.Spec.Triggers[0].Metadata)
		so.Spec.Triggers[0].AuthenticationRef = &keda.ScaledObjectAuthRef{
			Name: taName,
		}
		if err := controllerutil.SetControllerReference(connector, so, r.Scheme); err != nil {
			logger.Error(err, "Set ScaledObject ControllerReference error")
			return err
		}
		if err := createOrPatchObj(ctx, r, so, soName,
			connector.Namespace, logger, vance.SoResType); err != nil {
			return err
		}
	}

	return nil
}
