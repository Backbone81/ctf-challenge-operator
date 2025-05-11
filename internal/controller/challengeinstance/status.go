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

const (
	DefaultExpirationSeconds = int64(15 * 60) // 15 minutes
)

// StatusReconciler is responsible for reconciling the status of the challenge instance.
type StatusReconciler struct {
	client client.Client
}

func NewStatusReconciler(client client.Client) *StatusReconciler {
	return &StatusReconciler{
		client: client,
	}
}

func (r *StatusReconciler) SetupWithManager(ctrlBuilder *builder.Builder) *builder.Builder {
	return ctrlBuilder
}

func (r *StatusReconciler) Reconcile(ctx context.Context, challengeInstance *v1alpha1.ChallengeInstance) (ctrl.Result, error) {
	if !challengeInstance.DeletionTimestamp.IsZero() {
		// We do not update the status when the resource is already being deleted.
		return ctrl.Result{}, nil
	}

	updateStatus := false

	// calculate expiration timestamp
	if challengeInstance.Status.ExpirationTimestamp.IsZero() {
		expirationSeconds := DefaultExpirationSeconds
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
