package challengeinstance

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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
	if !challengeInstance.DeletionTimestamp.IsZero() {
		// We do not update the status when the resource is already being deleted.
		return ctrl.Result{}, nil
	}

	updateStatus := false

	// calculate expiration timestamp
	if challengeInstance.Status.ExpirationTimestamp.IsZero() {
		expirationSeconds := int64(15 * 60) // default is 15 minutes
		if challengeInstance.Spec.ExpirationSeconds != nil {
			expirationSeconds = *challengeInstance.Spec.ExpirationSeconds
		}
		challengeInstance.Status.ExpirationTimestamp = metav1.NewTime(time.Now().Add(time.Duration(expirationSeconds) * time.Second))
		updateStatus = true
	}

	if updateStatus {
		if err := r.client.Status().Update(ctx, challengeInstance); err != nil {
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}
