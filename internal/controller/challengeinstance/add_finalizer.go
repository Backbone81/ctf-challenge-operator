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

func NewAddFinalizerReconciler(client client.Client) *AddFinalizerReconciler {
	return &AddFinalizerReconciler{
		client: client,
	}
}

func (r *AddFinalizerReconciler) SetupWithManager(ctrlBuilder *builder.Builder) *builder.Builder {
	return ctrlBuilder
}

func (r *AddFinalizerReconciler) Reconcile(ctx context.Context, challengeInstance *v1alpha1.ChallengeInstance) (ctrl.Result, error) {
	if !challengeInstance.DeletionTimestamp.IsZero() {
		// We do not add the finalizer when the resource is already being deleted.
		return ctrl.Result{}, nil
	}

	if !controllerutil.AddFinalizer(challengeInstance, FinalizerName) {
		// The finalizer is already on the resource, nothing to do.
		return ctrl.Result{}, nil
	}

	if err := r.client.Update(ctx, challengeInstance); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}
