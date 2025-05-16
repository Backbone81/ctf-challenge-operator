package challengeinstance

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/backbone81/ctf-challenge-operator/api/v1alpha1"
	"github.com/backbone81/ctf-challenge-operator/internal/utils"
)

// RemoveFinalizerReconciler is responsible for removing the finalizer from the challenge instance.
type RemoveFinalizerReconciler struct {
	utils.DefaultSubReconciler
}

func NewRemoveFinalizerReconciler(client client.Client) *RemoveFinalizerReconciler {
	return &RemoveFinalizerReconciler{
		DefaultSubReconciler: utils.NewDefaultSubReconciler(client),
	}
}

func (r *RemoveFinalizerReconciler) Reconcile(ctx context.Context, challengeInstance *v1alpha1.ChallengeInstance) (ctrl.Result, error) {
	if challengeInstance.DeletionTimestamp.IsZero() {
		// We do not remove the finalizer when the resource is not being deleted.
		return ctrl.Result{}, nil
	}

	if !controllerutil.RemoveFinalizer(challengeInstance, FinalizerName) {
		// The finalizer is already gone from the resource, nothing to do.
		return ctrl.Result{}, nil
	}

	if err := r.GetClient().Update(ctx, challengeInstance); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}
