package apikey

import (
	"context"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/backbone81/ctf-challenge-operator/api/v1alpha1"
)

// DeleteReconciler is responsible for deleting the APIKey when it is expired.
type DeleteReconciler struct {
	client client.Client
}

// NewDeleteReconciler creates a new sub-reconciler instance. The reconciler is initialized with the given client.
func NewDeleteReconciler(client client.Client) *DeleteReconciler {
	return &DeleteReconciler{
		client: client,
	}
}

// SetupWithManager registers the sub-reconciler with the manager.
func (r *DeleteReconciler) SetupWithManager(ctrlBuilder *builder.Builder) *builder.Builder {
	return ctrlBuilder
}

// Reconcile is the main reconciler function.
func (r *DeleteReconciler) Reconcile(ctx context.Context, apiKey *v1alpha1.APIKey) (ctrl.Result, error) {
	if apiKey.Status.ExpirationTimestamp.Time.Before(time.Now()) {
		if err := r.client.Delete(ctx, apiKey); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	return ctrl.Result{RequeueAfter: time.Until(apiKey.Status.ExpirationTimestamp.Time)}, nil
}
