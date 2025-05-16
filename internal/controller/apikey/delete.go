package apikey

import (
	"context"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/backbone81/ctf-challenge-operator/api/v1alpha1"
	"github.com/backbone81/ctf-challenge-operator/internal/utils"
)

// DeleteReconciler is responsible for deleting the APIKey when it is expired.
type DeleteReconciler struct {
	utils.DefaultSubReconciler
}

// NewDeleteReconciler creates a new sub-reconciler instance. The reconciler is initialized with the given client.
func NewDeleteReconciler(client client.Client) *DeleteReconciler {
	return &DeleteReconciler{
		DefaultSubReconciler: utils.NewDefaultSubReconciler(client),
	}
}

// Reconcile is the main reconciler function.
func (r *DeleteReconciler) Reconcile(ctx context.Context, apiKey *v1alpha1.APIKey) (ctrl.Result, error) {
	if apiKey.Status.ExpirationTimestamp.Time.Before(time.Now()) {
		if err := r.GetClient().Delete(ctx, apiKey); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	return ctrl.Result{RequeueAfter: time.Until(apiKey.Status.ExpirationTimestamp.Time)}, nil
}
