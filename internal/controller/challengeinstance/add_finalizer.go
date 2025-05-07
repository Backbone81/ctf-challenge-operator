package challengeinstance

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/backbone81/ctf-challenge-operator/api/v1alpha1"
)

// AddFinalizerReconciler is responsible for adding the finalizer to the challenge instance.
type AddFinalizerReconciler struct {
	client client.Client
}

// NewAddFinalizerReconciler creates a new sub-reconciler instance. The reconciler is initialized with the given client.
func NewAddFinalizerReconciler(client client.Client) *AddFinalizerReconciler {
	return &AddFinalizerReconciler{
		client: client,
	}
}

// SetupWithManager registers the sub-reconciler with the manager.
func (r *AddFinalizerReconciler) SetupWithManager(ctrlBuilder *builder.Builder) *builder.Builder {
	return ctrlBuilder
}

// Reconcile is the main reconciler function.
func (r *AddFinalizerReconciler) Reconcile(ctx context.Context, challengeInstance *v1alpha1.ChallengeInstance) (ctrl.Result, error) {
	if !challengeInstance.DeletionTimestamp.IsZero() {
		// We do not add the finalizer when the resource is already being deleted.
		return ctrl.Result{}, nil
	}

	if !controllerutil.AddFinalizer(challengeInstance, ChallengeInstanceFinalizerName) {
		// The finalizer is already on the resource, nothing to do.
		return ctrl.Result{}, nil
	}

	if err := r.client.Update(ctx, challengeInstance); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}
