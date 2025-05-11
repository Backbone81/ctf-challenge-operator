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

func NewDeleteReconciler(client client.Client) *DeleteReconciler {
	return &DeleteReconciler{
		client: client,
	}
}

func (r *DeleteReconciler) SetupWithManager(ctrlBuilder *builder.Builder) *builder.Builder {
	return ctrlBuilder
}

func (r *DeleteReconciler) Reconcile(ctx context.Context, challengeInstance *v1alpha1.ChallengeInstance) (ctrl.Result, error) {
	if !challengeInstance.DeletionTimestamp.IsZero() {
		// We do not delete the resource when the resource is already being deleted.
		return ctrl.Result{}, nil
	}

	if challengeInstance.Status.ExpirationTimestamp.Time.Before(time.Now()) {
		if err := r.client.Delete(ctx, challengeInstance); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	return ctrl.Result{RequeueAfter: time.Until(challengeInstance.Status.ExpirationTimestamp.Time)}, nil
}
