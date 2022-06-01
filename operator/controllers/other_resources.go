package controllers

import (
	"context"
	"errors"
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
	case vance.SvcResType:
		{
			return createOrUpdateService(ctx, r, connector)
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
	var svcPort int32
	if svcPort = connector.Spec.Scaler.Metadata["svcPort"].IntVal; svcPort == 0 {
		err := errors.New("Error: missing required field of metadata <svcPort>. ")
		logger.Error(err, err.Error())
		return err
	}
	var host string
	if host = connector.Spec.Scaler.Metadata["host"].StrVal; host == "" {
		err := errors.New("Error: missing required field of metadata <host>. ")
		logger.Error(err, err.Error())
		return err
	}

	httpso := k8s2.CreateHttpScaledObject(
		connector.Namespace, connector.Name,
		host,
		svcPort)
	if minReplica := connector.Spec.Scaler.Metadata["minReplica"].IntVal; minReplica != 0 {
		httpso.Spec.Replicas.Min = minReplica
	}
	if maxReplica := connector.Spec.Scaler.Metadata["maxReplica"].IntVal; maxReplica != 0 {
		httpso.Spec.Replicas.Max = maxReplica
	}
	if pendingRequests := connector.Spec.Scaler.Metadata["pendingRequests"].IntVal; pendingRequests != 0 {
		httpso.Spec.TargetPendingRequests = pendingRequests
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
func createOrUpdateService(
	ctx context.Context,
	r *ConnectorReconciler,
	connector *vance.Connector) error {
	logger := r.Log.WithValues(
		"service name", connector.Name,
		"service namespace", connector.Namespace)
	logger.Info(" Start creating or updating a Service ")

	var port int32
	if port = connector.Spec.Scaler.Metadata["svcPort"].IntVal; port == 0 {
		err := errors.New("Error: missing required field of metadata <svcPort>. ")
		logger.Error(err, err.Error())
		return err
	}
	service := k8s2.CreateLBService(connector.Namespace, connector.Name,
		port)
	if err := controllerutil.SetControllerReference(connector, service, r.Scheme); err != nil {
		logger.Error(err, "Set service ControllerReference error")
		return err
	}
	if err := createOrPatchObj(ctx, r, service, connector.Name,
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
	logger.Info(" Create an in-memory Deployment ",
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
