package challengeinstance

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/backbone81/ctf-challenge-operator/api/v1alpha1"
)

// StatusReconciler is responsible for reconciling the status of the challenge instance.
type StatusReconciler struct {
	client client.Client
}

// NewStatusReconciler creates a new sub-reconciler instance. The reconciler is initialized with the given client.
func NewStatusReconciler(client client.Client) *StatusReconciler {
	return &StatusReconciler{
		client: client,
	}
}

// SetupWithManager registers the sub-reconciler with the manager.
func (r *StatusReconciler) SetupWithManager(ctrlBuilder *builder.Builder) *builder.Builder {
	return ctrlBuilder
}

// Reconcile is the main reconciler function.
func (r *StatusReconciler) Reconcile(ctx context.Context, challengeInstance *v1alpha1.ChallengeInstance) (ctrl.Result, error) {
	return ctrl.Result{}, nil
}
