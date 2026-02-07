package controller

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// EnsureFinalizer ensures that the finalizer is set on the resource.
func EnsureFinalizer(ctx context.Context, c client.Client, obj client.Object, finalizer string) error {
	if !controllerutil.ContainsFinalizer(obj, finalizer) {
		controllerutil.AddFinalizer(obj, finalizer)
		if err := c.Update(ctx, obj); err != nil {
			return err
		}
	}
	return nil
}

// RemoveFinalizer removes the finalizer from the resource.
func RemoveFinalizer(ctx context.Context, c client.Client, obj client.Object, finalizer string) error {
	if controllerutil.ContainsFinalizer(obj, finalizer) {
		controllerutil.RemoveFinalizer(obj, finalizer)
		if err := c.Update(ctx, obj); err != nil {
			return err
		}
	}
	return nil
}
