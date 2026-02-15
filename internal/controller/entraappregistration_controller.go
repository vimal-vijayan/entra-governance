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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	entragov "github.com/vimal-vijayan/entra-governance/api/v1alpha1"
	appregistration "github.com/vimal-vijayan/entra-governance/internal/services/applications"
	serviceprincipal "github.com/vimal-vijayan/entra-governance/internal/services/serviceprincipals"
)

// EntraAppRegistrationReconciler reconciles a EntraAppRegistration object
type EntraAppRegistrationReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	AppService *appregistration.Service
	SPService  *serviceprincipal.Service
}

// +kubebuilder:rbac:groups=iam.entra.governance.com,resources=entraappregistrations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=iam.entra.governance.com,resources=entraappregistrations/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=iam.entra.governance.com,resources=entraappregistrations/finalizers,verbs=update

func (r *EntraAppRegistrationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("------------------ Reconciling EntraAppRegistration --------------------", "name", req.Name, "namespace", req.Namespace)

	entraAppReg := &entragov.EntraAppRegistration{}
	if err := r.Get(ctx, req.NamespacedName, entraAppReg); err != nil {
		if client.IgnoreNotFound(err) != nil {
			logger.Error(err, "EntraAppRegistration resource not found. skipping reconciliation.")
		}
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

	if entraAppReg.Status.AppRegistrationID != "" && entraAppReg.Status.AppRegistrationObjID != "" {
		if err := r.updateAppRegistration(ctx, entraAppReg); err != nil {
			return ctrl.Result{RequeueAfter: defaultRequeueDuration}, err
		}
	}

	if entraAppReg.Status.AppRegistrationID == "" && entraAppReg.Status.AppRegistrationObjID == "" {
		return r.createAppRegistration(ctx, entraAppReg)
	}

	if err := r.ensureServicePrincipal(ctx, entraAppReg); err != nil {
		logger.Error(err, "Failed to ensure service principal for app registration", "appName", entraAppReg.Name)
		return ctrl.Result{RequeueAfter: defaultRequeueDuration}, err
	}

	return ctrl.Result{RequeueAfter: defaultRequeueDuration}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *EntraAppRegistrationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&entragov.EntraAppRegistration{}).
		Complete(r)
}

// ensure service principal
func (r *EntraAppRegistrationReconciler) ensureServicePrincipal(ctx context.Context, entraAppReg *entragov.EntraAppRegistration) error {
	logger := log.FromContext(ctx)

	if !entraAppReg.Spec.ServicePrincipal.Enabled {
		if entraAppReg.Status.ServicePrincipal == "Enabled" {
			logger.Info("Service principal exists but is now disabled in spec. Deleting service principal.", "appName", entraAppReg.Name)
			if err := r.SPService.Delete(ctx, entraAppReg.Status.ServicePrincipalID, *entraAppReg); err != nil {
				logger.Error(err, "Failed to delete service principal for app registration", "appName", entraAppReg.Name)
				return err
			}
			entraAppReg.Status.ServicePrincipalID = ""
			entraAppReg.Status.ServicePrincipal = "Disabled"
			if err := r.Status().Update(ctx, entraAppReg); err != nil {
				logger.Error(err, "Failed to update EntraAppRegistration status after deleting service principal", "appName", entraAppReg.Name)
				return err
			}
			logger.Info("Service principal deleted successfully for app registration", "appName", entraAppReg.Name)
		}
		logger.Info("Service principal creation is disabled for this app registration. Skipping service principal creation.", "appName", entraAppReg.Name)
		return nil
	}

	if entraAppReg.Status.ServicePrincipalID != "" {
		logger.Info("Service principal already exists for this app registration. Skipping service principal creation.", "appName", entraAppReg.Name)
		return nil
	}

	spId, err := r.SPService.Create(ctx, *entraAppReg)
	if err != nil {
		logger.Error(err, "Failed to create service principal for app registration", "appName", entraAppReg.Name)
		return err
	}

	entraAppReg.Status.ServicePrincipalID = spId
	entraAppReg.Status.ServicePrincipal = "Enabled"

	if err := r.updateReconilerStatus(ctx, *entraAppReg); err != nil {
		logger.Error(err, "Failed to update EntraAppRegistration status with service principal ID", "appName", entraAppReg.Name)
		return err
	}

	logger.Info("Service principal created successfully for app registration", "appName", entraAppReg.Name, "servicePrincipalID", spId)
	return nil
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
	entraAppReg.Status.Phase = "Pending"
	entraAppReg.Status.AppRegistrationName = entraAppReg.Name
	entraAppReg.Status.ObservedGeneration = entraAppReg.Generation
	//TODO: Add conditions and messages in status
	// entraAppReg.Status.Message = "App registration created successfully in Entra"

	if err := r.updateReconilerStatus(ctx, *entraAppReg); err != nil {
		logger.Error(err, "Failed to update EntraAppRegistration status after creation", "appName", entraAppReg.Name)
		return ctrl.Result{RequeueAfter: defaultRequeueDuration}, err
	}

	logger.Info("EntraAppRegistration created successfully in Entra", "appName", entraAppReg.Name, "clientId", clientId, "principalId", principalId)
	return ctrl.Result{RequeueAfter: defaultRequeueDuration}, nil
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

		return ctrl.Result{}, nil
	}

	err := r.AppService.Delete(ctx, entraAppReg.Status.AppRegistrationObjID, *entraAppReg)
	if err != nil {
		logger.Error(err, "Failed to delete app registration in Entra", "appName", entraAppReg.Name, "clientId", entraAppReg.Status.AppRegistrationID)
		return ctrl.Result{RequeueAfter: defaultRequeueDuration}, err
	}

	if err := RemoveFinalizer(ctx, r.Client, entraAppReg, entraAppRegistrationFinalizer); err != nil {
		logger.Error(err, "Failed to remove finalizer from EntraAppRegistration after deletion", "appName", entraAppReg.Name)
		return ctrl.Result{RequeueAfter: defaultRequeueDuration}, err
	}

	logger.Info("Entra App Registration deleted successfully in Entra", "appName", entraAppReg.Name, "clientId", entraAppReg.Status.AppRegistrationID)
	return ctrl.Result{}, nil
}

func (r *EntraAppRegistrationReconciler) updateAppRegistration(ctx context.Context, entraAppReg *entragov.EntraAppRegistration) error {
	logger := log.FromContext(ctx)

	// if entraAppReg.Status.ObservedGeneration >= entraAppReg.Generation {
	// 	return nil
	// }

	logger.Info("reconciling app registration drift", "appName", entraAppReg.Name, "observedGeneration", entraAppReg.Status.ObservedGeneration, "generation", entraAppReg.Generation)

	updated, err := r.AppService.GetAndPatch(ctx, *entraAppReg)
	if err != nil {
		logger.Error(err, "Failed to update app registration in Entra", "appName", entraAppReg.Name)
		return err
	}
	if updated {
		logger.Info("app registration patched in Entra", "appName", entraAppReg.Name)
	} else {
		logger.Info("app registration already in sync", "appName", entraAppReg.Name)
	}

	if err := r.AppService.UpdateOwners(ctx, entraAppReg.Status.AppRegistrationObjID, *entraAppReg); err != nil {
		logger.Error(err, "Failed to update app owners in Entra", "appName", entraAppReg.Name)
		entraAppReg.Status.Phase = "Warning"
		entraAppReg.Status.Message = "App registration is in sync with Entra but failed to update owners: " + err.Error()
		entraAppReg.Status.LastRun = metav1.Now()
		return r.updateReconilerStatus(ctx, *entraAppReg)
	}

	entraAppReg.Status.Phase = "Available"
	entraAppReg.Status.ObservedGeneration = entraAppReg.Generation
	entraAppReg.Status.Message = "App registration is in sync with Entra"
	entraAppReg.Status.LastRun = metav1.Now()
	if entraAppReg.Spec.Owners != nil {
		entraAppReg.Status.Owners = append([]string(nil), (*entraAppReg.Spec.Owners)...)
	} else {
		entraAppReg.Status.Owners = nil
	}

	return r.updateReconilerStatus(ctx, *entraAppReg)
}
