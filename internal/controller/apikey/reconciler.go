package apikey

import (
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/backbone81/ctf-challenge-operator/api/v1alpha1"
	"github.com/backbone81/ctf-challenge-operator/internal/utils"
)

// +kubebuilder:rbac:groups=core.ctf.backbone81,resources=apikeys,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core.ctf.backbone81,resources=apikeys/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core.ctf.backbone81,resources=apikeys/finalizers,verbs=update

func NewReconciler(client client.Client, options ...utils.ReconcilerOption[*v1alpha1.APIKey]) *utils.Reconciler[*v1alpha1.APIKey] {
	return utils.NewReconciler[*v1alpha1.APIKey](
		client,
		func() *v1alpha1.APIKey {
			return &v1alpha1.APIKey{}
		},
		options...,
	)
}

// WithDefaultReconcilers returns a reconciler option which enables the default sub-reconcilers.
func WithDefaultReconcilers() utils.ReconcilerOption[*v1alpha1.APIKey] {
	return func(reconciler *utils.Reconciler[*v1alpha1.APIKey]) {
		WithStatusReconciler()(reconciler)
		WithDeleteReconciler()(reconciler)
	}
}

func WithStatusReconciler() utils.ReconcilerOption[*v1alpha1.APIKey] {
	return func(reconciler *utils.Reconciler[*v1alpha1.APIKey]) {
		reconciler.AppendSubReconciler(NewStatusReconciler(reconciler.GetClient()))
	}
}

func WithDeleteReconciler() utils.ReconcilerOption[*v1alpha1.APIKey] {
	return func(reconciler *utils.Reconciler[*v1alpha1.APIKey]) {
		reconciler.AppendSubReconciler(NewDeleteReconciler(reconciler.GetClient()))
	}
}
