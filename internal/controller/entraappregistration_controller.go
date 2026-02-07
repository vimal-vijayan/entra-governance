/*
Copyright 2026.

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

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	entragov "github.com/vimal-vijayan/entra-governance/api/v1alpha1"
	appregistration "github.com/vimal-vijayan/entra-governance/internal/services/applications"
)

// EntraAppRegistrationReconciler reconciles a EntraAppRegistration object
type EntraAppRegistrationReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	AppService *appregistration.Service
}

// +kubebuilder:rbac:groups=iam.entra.governance.com,resources=entraappregistrations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=iam.entra.governance.com,resources=entraappregistrations/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=iam.entra.governance.com,resources=entraappregistrations/finalizers,verbs=update

func (r *EntraAppRegistrationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("------------------ Reconciling EntraAppRegistration --------------------", "name", req.Name, "namespace", req.Namespace)

	entraAppReg := &entragov.EntraAppRegistration{}
	if err := r.Get(ctx, req.NamespacedName, entraAppReg); err != nil {
		logger.Error(err, "Failed to get EntraAppRegistration")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Ensure finalizer is present
	if err := EnsureFinalizer(ctx, r.Client, entraAppReg, entraAppRegistrationFinalizer); err != nil {
		return ctrl.Result{RequeueAfter: defaultRequeueDuration}, err
	}

	if !entraAppReg.DeletionTimestamp.IsZero() {
		logger.Info("EntraAppRegistration resource is being deleted. skipping reconciliation.")
		return r.deleteAppRegistration(ctx, entraAppReg)
	}

	if entraAppReg.Status.AppRegistrationID == "" && entraAppReg.Status.AppRegistrationObjID == "" {
		return r.createAppRegistration(ctx, entraAppReg)
	} else {
		logger.Info("EntraAppRegistration already exists in status. skipping creation.", "appName", entraAppReg.Name, "clientId", entraAppReg.Status.AppRegistrationID)
		logger.Info("Reconciling entra app registration attributes.")
		return reconcileAppRegistrationAttributes(ctx, entraAppReg)
	}

	// appregistration upsert behavior
	
	return ctrl.Result{RequeueAfter: defaultRequeueDuration}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *EntraAppRegistrationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&entragov.EntraAppRegistration{}).
		Complete(r)
}

func (r *EntraAppRegistrationReconciler) createAppRegistration(ctx context.Context, entraAppReg *entragov.EntraAppRegistration) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	clientId, principalId, err := r.AppService.Create(ctx, *entraAppReg)
	if err != nil {
		logger.Error(err, "Failed to create app registration in Entra", "appName", entraAppReg.Name)
		return ctrl.Result{RequeueAfter: defaultRequeueDuration}, err
	}

	entraAppReg.Status.AppRegistrationID = clientId
	entraAppReg.Status.AppRegistrationObjID = principalId
	entraAppReg.Status.Phase = "Available"
	entraAppReg.Status.AppRegistrationName = entraAppReg.Name
	entraAppReg.Status.ObservedGeneration = entraAppReg.Generation
	//TODO: Add conditions and messages in status
	// entraAppReg.Status.Message = "App registration created successfully in Entra"

	if err := r.Status().Update(ctx, entraAppReg); err != nil {
		logger.Error(err, "Failed to update EntraAppRegistration status after creation", "appName", entraAppReg.Name)
		return ctrl.Result{RequeueAfter: defaultRequeueDuration}, err
	}

	if err := r.Status().Update(ctx, entraAppReg); err != nil {
		logger.Error(err, "Failed to update EntraAppRegistration status after creation", "appName", entraAppReg.Name)
		return ctrl.Result{RequeueAfter: defaultRequeueDuration}, err
	}

	logger.Info("EntraAppRegistration created successfully in Entra", "appName", entraAppReg.Name, "clientId", clientId, "principalId", principalId)
	return ctrl.Result{Requeue: true}, nil
}

func (r *EntraAppRegistrationReconciler) deleteAppRegistration(ctx context.Context, entraAppReg *entragov.EntraAppRegistration) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	if entraAppReg.Status.AppRegistrationID == "" && entraAppReg.Status.AppRegistrationObjID == "" {
		logger.Info("App Registration ID is empty in status, assuming app has already been deleted", "appName", entraAppReg.Name)
		logger.Info("Removing finalizer for EntraAppRegistration", "appName", entraAppReg.Name)
		if err := RemoveFinalizer(ctx, r.Client, entraAppReg, entraAppRegistrationFinalizer); err != nil {
			logger.Error(err, "Failed to remove finalizer from EntraAppRegistration", "appName", entraAppReg.Name)
			return ctrl.Result{RequeueAfter: defaultRequeueDuration}, err
		}
		logger.Info("Finalizer removed successfully from EntraAppRegistration", "appName", entraAppReg.Name)
		return ctrl.Result{RequeueAfter: defaultRequeueDuration}, nil
	}

	err := r.AppService.Delete(ctx, entraAppReg.Status.AppRegistrationID, *entraAppReg)
	if err != nil {
		logger.Error(err, "Failed to delete app registration in Entra", "appName", entraAppReg.Name, "clientId", entraAppReg.Status.AppRegistrationID)
		return ctrl.Result{RequeueAfter: defaultRequeueDuration}, err
	}

	logger.Info("Entra App Registration deleted successfully in Entra", "appName", entraAppReg.Name, "clientId", entraAppReg.Status.AppRegistrationID)
	return ctrl.Result{Requeue: true}, nil
}

func reconcileAppRegistrationAttributes(ctx context.Context, entraAppReg *entragov.EntraAppRegistration) (ctrl.Result, error) {
	// Placeholder for future attribute reconciliation logic
	logger := log.FromContext(ctx)
	logger.Info("Reconciliation of EntraAppRegistration attributes is not yet implemented", "appName", entraAppReg.Name)

	// TODO: Implement attribute reconciliation logic here. This may include:
	// upsert : Create a new application if it doesn't exist
	// upsert : Update an existing application
	// DOC: https://learn.microsoft.com/en-us/graph/api/application-upsert?view=graph-rest-1.0&tabs=http#examples
	return ctrl.Result{RequeueAfter: defaultRequeueDuration}, nil
}