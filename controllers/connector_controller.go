/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"github.com/go-logr/logr"
	vance "github.com/linkall-labs/vance/api/v1alpha1"
	"github.com/linkall-labs/vance/pkg/config"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func init() {
	ScalerConf = config.LoadScalerConfig()
}

var ScalerConf map[string][][]string

// ConnectorReconciler reconciles a Connector object
type ConnectorReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=vance.io,resources=connectors,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=vance.io,resources=connectors/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=vance.io,resources=connectors/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=keda.sh,resources=scaledobjects,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=keda.sh,resources=triggerauthentications,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=http.keda.sh,resources=httpscaledobjects,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Connector object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *ConnectorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := r.Log.WithValues("Connector.Namespace", req.Namespace, "Connector.Name", req.Name)
	logger.Info("Start reconciling")
	for k, v := range ScalerConf {
		logger.Info("scalerConfig", "key", k, "value", v)
	}
	connector := &vance.Connector{}
	err := r.Get(ctx, req.NamespacedName, connector)

	if err != nil {

		// If the connector doesn't exist, we just return to avoid reconciliation.
		if errors.IsNotFound(err) {
			logger.Info("Connector is not found. It might be deleted already.")
			return reconcile.Result{}, nil
		}

		logger.Error(err, "Getting connector failed")
		// 返回错误信息给外部
		return ctrl.Result{}, err
	}
	if err = createOrUpdateAPIResources(ctx, r, connector, vance.DeployResType); err != nil {
		return ctrl.Result{}, err
	}
	if connector.Spec.Scaler == nil {
		logger.Info("scaler is nil")
	} else {
		logger.Info("scaler is not nil")
		// build a http scaler if the type is http, otherwise build other scalers
		if connector.Spec.ConnectorType == vance.Http {
			if err = createOrUpdateAPIResources(ctx, r, connector, vance.SvcResType); err != nil {
				return ctrl.Result{}, err
			}
			if err = createOrUpdateAPIResources(ctx, r, connector, vance.HttpSOResType); err != nil {
				return ctrl.Result{}, err
			}
		} else {
			if err = createOrUpdateAPIResources(ctx, r, connector, vance.SoResType); err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ConnectorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		// Uncomment the following line adding a pointer to an instance of the controlled resource as an argument
		For(&vance.Connector{}).
		Complete(r)
}
