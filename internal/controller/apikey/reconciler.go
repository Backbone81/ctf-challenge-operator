package apikey

import (
	"context"

	"github.com/backbone81/ctf-challenge-operator/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// +kubebuilder:rbac:groups=core.ctf.backbone81,resources=apikeys,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core.ctf.backbone81,resources=apikeys/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core.ctf.backbone81,resources=apikeys/finalizers,verbs=update

// Reconciler provides functionality for provisioning API keys.
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
	ctrlBuilder := ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.APIKey{})
	for _, subReconciler := range r.subReconcilers {
		ctrlBuilder = subReconciler.SetupWithManager(ctrlBuilder)
	}
	return ctrlBuilder.Complete(r)
}

// Reconcile is the main reconciler function.
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	apiKey, err := r.getAPIKey(ctx, req)
	if err != nil {
		return ctrl.Result{}, err
	}
	if apiKey == nil {
		// The resource was deleted.
		return ctrl.Result{}, nil
	}

	for _, subReconciler := range r.subReconcilers {
		result, err := subReconciler.Reconcile(ctx, apiKey)
		if err != nil || !result.IsZero() {
			return result, err
		}
	}
	return ctrl.Result{}, nil
}

func (r *Reconciler) getAPIKey(ctx context.Context, req ctrl.Request) (*v1alpha1.APIKey, error) {
	var result v1alpha1.APIKey
	if err := r.client.Get(ctx, req.NamespacedName, &result); err != nil {
		return nil, client.IgnoreNotFound(err)
	}
	return &result, nil
}

// SubReconciler is the interface all sub-reconcilers need to implement.
type SubReconciler interface {
	Reconcile(ctx context.Context, apiKey *v1alpha1.APIKey) (ctrl.Result, error)
	SetupWithManager(builder *builder.Builder) *builder.Builder
}

// ReconcilerOption is an option which can be applied to the reconciler.
type ReconcilerOption func(reconciler *Reconciler)

// WithDefaultReconcilers returns a reconciler option which enables the default sub-reconcilers.
func WithDefaultReconcilers() ReconcilerOption {
	return func(reconciler *Reconciler) {
		WithStatusReconciler()(reconciler)
		WithDeleteReconciler()(reconciler)
	}
}

func WithStatusReconciler() ReconcilerOption {
	return func(reconciler *Reconciler) {
		reconciler.subReconcilers = append(reconciler.subReconcilers, NewStatusReconciler(reconciler.client))
	}
}

func WithDeleteReconciler() ReconcilerOption {
	return func(reconciler *Reconciler) {
		reconciler.subReconcilers = append(reconciler.subReconcilers, NewDeleteReconciler(reconciler.client))
	}
}
