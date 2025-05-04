package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/backbone81/ctf-challenge-operator/internal/controller/apikey"
	"github.com/backbone81/ctf-challenge-operator/internal/controller/challengeinstance"
)

// Reconciler is the main reconciler of this operator. It is responsible for registering and running all
// sub-reconcilers.
type Reconciler struct {
	client         client.Client
	subReconcilers []SubReconciler
}

// NewReconciler creates a new reconciler instance. The reconciler is initialized with the given client and applies
// the provided options to the reconciler.
func NewReconciler(client client.Client, options ...ReconcilerOption) *Reconciler {
	result := &Reconciler{
		client: client,
	}
	for _, option := range options {
		option(result)
	}
	return result
}

// SetupWithManager registers all enabled sub-reconcilers with the given manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	for _, subReconciler := range r.subReconcilers {
		if err := subReconciler.SetupWithManager(mgr); err != nil {
			return err
		}
	}
	return nil
}

// SubReconciler is the interface all sub-reconcilers need to implement.
type SubReconciler interface {
	reconcile.Reconciler
	SetupWithManager(mgr ctrl.Manager) error
}

// ReconcilerOption is an option which can be applied to the reconciler.
type ReconcilerOption func(reconciler *Reconciler)

// WithDefaultReconcilers returns a reconciler option which enables the default sub-reconcilers.
func WithDefaultReconcilers() ReconcilerOption {
	return func(reconciler *Reconciler) {
		WithAPIKeyReconciler()(reconciler)
		WithChallengeInstanceReconciler()(reconciler)
	}
}

// WithAPIKeyReconciler returns a reconciler option which enables the APIKey sub-reconciler.
func WithAPIKeyReconciler() ReconcilerOption {
	return func(reconciler *Reconciler) {
		reconciler.subReconcilers = append(
			reconciler.subReconcilers,
			apikey.NewReconciler(reconciler.client),
		)
	}
}

// WithChallengeInstanceReconciler returns a reconciler option which enables the ChallengeInstance sub-reconciler.
func WithChallengeInstanceReconciler() ReconcilerOption {
	return func(reconciler *Reconciler) {
		reconciler.subReconcilers = append(
			reconciler.subReconcilers,
			challengeinstance.NewReconciler(reconciler.client),
		)
	}
}
