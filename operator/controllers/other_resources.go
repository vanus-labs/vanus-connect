package controllers

import (
	"context"
	"github.com/go-logr/logr"
	kedahttp "github.com/kedacore/http-add-on/operator/api/v1alpha1"
	keda "github.com/kedacore/keda/v2/apis/keda/v1alpha1"
	vance "github.com/linkall-labs/vance/operator/api/v1alpha1"
	k8s2 "github.com/linkall-labs/vance/operator/pkg/k8s"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrs "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// this method is used to create or update an API resource
func createOrUpdateAPIResources(
	ctx context.Context,
	r *ConnectorReconciler,
	connector *vance.Connector,
	resourceType vance.ResourceType,
) error {
	logger := r.Log.WithValues("createOrUpdateAPIResources type", resourceType)
	logger.Info("Start creating or updating an API resource ")
	switch resourceType {
	case vance.DeployResType:
		{
			return createOrUpdateDeployment(ctx, r, connector)
		}
	case vance.SoResType:
		{
			return createOrUpdateScaledObject(ctx, r, connector)
		}
	case vance.CLSvcResType:
		{
			return createOrUpdateService(ctx, r, connector, vance.CLSvcResType)
		}
	case vance.LBSvcResType:
		{
			return createOrUpdateService(ctx, r, connector, vance.LBSvcResType)
		}
	case vance.HttpSOResType:
		{
			return createOrUpdateHttpScaledObject(ctx, r, connector)
		}
	default:
		return nil
	}
}

func createOrUpdateHttpScaledObject(
	ctx context.Context,
	r *ConnectorReconciler,
	connector *vance.Connector) error {
	logger := r.Log.WithValues(
		"httpScaledObject name", connector.Name,
		"httpScaledObject namespace", connector.Namespace)
	logger.Info(" Start creating or updating a httpScaledObject ")
	var svcPort int32 = *connector.Spec.ExposePort

	httpso := k8s2.CreateHttpScaledObject(
		connector.Namespace, connector.Name,
		connector.Spec.ScalingRule.HTTPScaling.Host,
		svcPort)
	if connector.Spec.ScalingRule.MinReplicaCount != nil {
		httpso.Spec.Replicas.Min = *connector.Spec.ScalingRule.MinReplicaCount
	} else {
		httpso.Spec.Replicas.Min = 0
	}
	if connector.Spec.ScalingRule.MaxReplicaCount != nil {
		httpso.Spec.Replicas.Max = *connector.Spec.ScalingRule.MaxReplicaCount
	} else {
		httpso.Spec.Replicas.Max = 10
	}

	if connector.Spec.ScalingRule.HTTPScaling.PendingRequests != 0 {
		httpso.Spec.TargetPendingRequests = connector.Spec.ScalingRule.HTTPScaling.PendingRequests
	}
	if err := controllerutil.SetControllerReference(connector, httpso, r.Scheme); err != nil {
		logger.Error(err, "Set service ControllerReference error")
		return err
	}
	if err := createOrPatchObj(ctx, r, httpso, connector.Name,
		connector.Namespace, logger, vance.HttpSOResType); err != nil {
		return err
	}
	return nil
}

// createOrUpdateService is used to generate SVCs for HTTP Scaling rules
func createOrUpdateService(
	ctx context.Context,
	r *ConnectorReconciler,
	connector *vance.Connector,
	resourceType vance.ResourceType,
) error {
	logger := r.Log.WithValues(
		"service name", connector.Name,
		"service namespace", connector.Namespace)
	logger.Info(" Start creating or updating a Service ")

	var svcPort int32 = *connector.Spec.ExposePort
	var svc *corev1.Service
	switch resourceType {
	case vance.CLSvcResType:
		{
			svc = k8s2.CreateCLUService(connector.Namespace, connector.Name, svcPort)
		}
	case vance.LBSvcResType:
		{
			svc = k8s2.CreateLBService(connector.Namespace, connector.Name, svcPort)
		}
	}

	if err := controllerutil.SetControllerReference(connector, svc, r.Scheme); err != nil {
		logger.Error(err, "Set service ControllerReference error")
		return err
	}
	if err := createOrPatchObj(ctx, r, svc, connector.Name,
		connector.Namespace, logger, vance.SvcResType); err != nil {
		return err
	}
	return nil
}
func createOrUpdateDeployment(
	ctx context.Context,
	r *ConnectorReconciler,
	connector *vance.Connector) error {
	logger := r.Log.WithValues(
		"deployment name", connector.Name,
		"deployment namespace", connector.Namespace)
	logger.Info(" Start creating or updating a Deployment ")
	deployment := k8s2.CreateBaseDeployment(connector.Namespace, connector.Name)
	if err := controllerutil.SetControllerReference(connector, deployment, r.Scheme); err != nil {
		logger.Error(err, "Set deployment ControllerReference error")
		return err
	}
	userConfig := &corev1.ConfigMap{}
	var cmVolume corev1.Volume
	if connector.Spec.ConfigRef != "" {
		existedKey := client.ObjectKey{
			Namespace: connector.Namespace,
			Name:      connector.Spec.ConfigRef,
		}
		if err := r.Get(ctx, existedKey, userConfig); err != nil {
			if k8serrs.IsNotFound(err) {
				logger.Error(err, "no such ConfigMap ", "configRef", connector.Spec.ConfigRef,
					"namespace", connector.Namespace)
			} else {
				logger.Error(err, "fetch ConfigMap err")
			}
			return err
		}
		cmVolume = corev1.Volume{Name: connector.Name + "-cmv"}
		cmVolume.ConfigMap = &corev1.ConfigMapVolumeSource{
			LocalObjectReference: corev1.LocalObjectReference{
				Name: connector.Spec.ConfigRef,
			},
		}
	}
	if connector.Spec.Containers != nil {
		logger.Info("custom pod containers")
		deployment.Spec.Template.Spec.Containers = connector.Spec.Containers
	} else {
		logger.Info("simply provide image url")
		deployment.Spec.Template.Spec.Containers = []corev1.Container{
			{
				Name: connector.Name,
				// 用指定的镜像
				Image:           connector.Spec.Image,
				ImagePullPolicy: "IfNotPresent",
			},
		}
	}
	// Add configMap volume to Spec.Volumes if the cmVolume is ready
	// Also add a VolumeMount to Containers[0].VolumeMounts and set the MountPath as "/vance/config"
	if cmVolume.Name != "" {
		deployment.Spec.Template.Spec.Volumes = []corev1.Volume{
			cmVolume,
		}
		deployment.Spec.Template.Spec.Containers[0].VolumeMounts = []corev1.VolumeMount{
			{
				Name:      cmVolume.Name,
				MountPath: "/vance/config",
			},
		}
	}

	logger.Info("Create an in-memory Deployment ",
		"deployment", *deployment)
	if err := createOrPatchObj(ctx, r, deployment, connector.Name,
		connector.Namespace, logger, vance.DeployResType); err != nil {
		return err
	}

	return nil
}

func createOrPatchObj(
	ctx context.Context,
	r *ConnectorReconciler,
	obj client.Object,
	name, nameSpace string,
	logger logr.Logger,
	resourceType vance.ResourceType) error {
	logger.Info(" Create or patch obj ", "resourceType", resourceType)
	if err := r.Create(ctx, obj); err != nil {
		if k8serrs.IsAlreadyExists(err) {
			existedKey := client.ObjectKey{
				Namespace: nameSpace,
				Name:      name,
			}
			var fetchedObj client.Object
			switch resourceType {
			case vance.DeployResType:
				{
					fetchedObj = &appsv1.Deployment{}
				}
			case vance.SoResType:
				{
					fetchedObj = &keda.ScaledObject{}
				}
			case vance.SvcResType:
				{
					fetchedObj = &corev1.Service{}
				}
			case vance.HttpSOResType:
				{
					fetchedObj = &kedahttp.HTTPScaledObject{}
				}
			case vance.TAResType:
				{
					fetchedObj = &keda.TriggerAuthentication{}
				}
			}
			if err := r.Get(ctx, existedKey, fetchedObj); err != nil {
				logger.Error(
					err,
					"[ERROR MSG, failed to fetch the existing resource]",
					"object", resourceType,
				)
				return err
			}
			if err := r.Patch(ctx, obj, client.Merge); err != nil {
				logger.Error(
					err,
					"failed to patch existing resource",
					"object", resourceType,
				)
				return err
			} else {
				logger.Info(" update existing resource success ", "resourceType", resourceType)
			}
		} else {
			logger.Error(err, "create or patch occurs other errors")
		}
	} else {
		logger.Info(" Create new resource success ", "resourceType", resourceType)
	}
	return nil
}
func resourceName(name string, resourceType vance.ResourceType) string {
	switch resourceType {
	case vance.DeployResType:
		{
			return name
		}
	case vance.SoResType:
		{
			return name + "-vso"
		}
	case vance.SvcResType:
		{
			return name
		}
	case vance.HttpSOResType:
		{
			return name + "-vhso"
		}
	case vance.TAResType:
		{
			return name + "-vta"
		}
	}
	return ""
}
