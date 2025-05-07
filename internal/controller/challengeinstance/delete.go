package challengeinstance

import (
	"context"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/backbone81/ctf-challenge-operator/api/v1alpha1"
)

// DeleteReconciler is responsible for deleting the challenge instance when it is expired.
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
func (r *DeleteReconciler) Reconcile(ctx context.Context, challengeInstance *v1alpha1.ChallengeInstance) (ctrl.Result, error) {
	if challengeInstance.Status.ExpirationTimestamp.Time.Before(time.Now()) {
		if err := r.client.Delete(ctx, challengeInstance); err != nil {
			// It can happen that the resource is already deleted. When the expiration is hit, all sub-reconcilers are
			// run and at last the delete reconciler will mark the resource for deletion. Deletion will not happen,
			// as we have a finalizer set. Another reconcile run is done and before the delete reconciler is run, the
			// finalizer is removed. Kubernetes then deletes the resource before the delete reconciler runs again.
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
		return ctrl.Result{}, nil
	}

	return ctrl.Result{RequeueAfter: time.Until(challengeInstance.Status.ExpirationTimestamp.Time)}, nil
}
