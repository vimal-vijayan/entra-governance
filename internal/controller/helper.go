package controller

import (
	"context"

	entragov "github.com/vimal-vijayan/entra-governance/api/v1alpha1"
)

func (r *EntraAppRegistrationReconciler) updateReconilerStatus(ctx context.Context, entraApp entragov.EntraAppRegistration) error {

	if err := r.Status().Update(ctx, &entraApp); err != nil {
		return err
	}

	return nil
}
