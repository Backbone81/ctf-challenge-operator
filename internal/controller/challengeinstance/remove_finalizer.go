package challengeinstance

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/backbone81/ctf-challenge-operator/api/v1alpha1"
)

// RemoveFinalizerReconciler is responsible for adding the finalizer to the challenge instance.
type RemoveFinalizerReconciler struct {
	client client.Client
}

// NewRemoveFinalizerReconciler creates a new sub-reconciler instance. The reconciler is initialized with the given client.
func NewRemoveFinalizerReconciler(client client.Client) *RemoveFinalizerReconciler {
	return &RemoveFinalizerReconciler{
		client: client,
	}
}

// SetupWithManager registers the sub-reconciler with the manager.
func (r *RemoveFinalizerReconciler) SetupWithManager(ctrlBuilder *builder.Builder) *builder.Builder {
	return ctrlBuilder
}

// Reconcile is the main reconciler function.
func (r *RemoveFinalizerReconciler) Reconcile(ctx context.Context, challengeInstance *v1alpha1.ChallengeInstance) (ctrl.Result, error) {
	if challengeInstance.DeletionTimestamp.IsZero() {
		// We do not remove the finalizer when the resource is not being deleted.
		return ctrl.Result{}, nil
	}

	if !controllerutil.RemoveFinalizer(challengeInstance, ChallengeInstanceFinalizerName) {
		// The finalizer is already gone from the resource, nothing to do.
		return ctrl.Result{}, nil
	}

	if err := r.client.Update(ctx, challengeInstance); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}
