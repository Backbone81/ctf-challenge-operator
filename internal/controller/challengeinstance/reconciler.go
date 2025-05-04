package challengeinstance

import (
	"context"

	"github.com/backbone81/ctf-challenge-operator/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
)

// +kubebuilder:rbac:groups=core.ctf.backbone81,resources=challengeinstances,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core.ctf.backbone81,resources=challengeinstances/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core.ctf.backbone81,resources=challengeinstances/finalizers,verbs=update

// Reconciler provides functionality for provisioning challenge instances.
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
		For(&v1alpha1.ChallengeInstance{})
	for _, subReconciler := range r.subReconcilers {
		ctrlBuilder = subReconciler.SetupWithManager(ctrlBuilder)
	}
	return ctrlBuilder.Complete(r)
}

// Reconcile is the main reconciler function.
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = ctrllog.FromContext(ctx)

	challengeInstance, err := r.getChallengeInstance(ctx, req)
	if err != nil {
		return ctrl.Result{}, err
	}
	if challengeInstance == nil {
		// The resource was deleted.
		return ctrl.Result{}, nil
	}

	for _, subReconciler := range r.subReconcilers {
		result, err := subReconciler.Reconcile(ctx, challengeInstance)
		if err != nil || !result.IsZero() {
			return result, err
		}
	}
	return ctrl.Result{}, nil
}

func (r *Reconciler) getChallengeInstance(ctx context.Context, req ctrl.Request) (*v1alpha1.ChallengeInstance, error) {
	var result v1alpha1.ChallengeInstance
	if err := r.client.Get(ctx, req.NamespacedName, &result); err != nil {
		return nil, client.IgnoreNotFound(err)
	}
	return &result, nil
}

// SubReconciler is the interface all sub-reconcilers need to implement.
type SubReconciler interface {
	Reconcile(ctx context.Context, challengeInstance *v1alpha1.ChallengeInstance) (ctrl.Result, error)
	SetupWithManager(builder *builder.Builder) *builder.Builder
}

// ReconcilerOption is an option which can be applied to the reconciler.
type ReconcilerOption func(reconciler *Reconciler)

// WithDefaultReconcilers returns a reconciler option which enables the default sub-reconcilers.
func WithDefaultReconcilers() ReconcilerOption {
	return func(reconciler *Reconciler) {
		WithStatusReconciler()(reconciler)
	}
}

func WithStatusReconciler() ReconcilerOption {
	return func(reconciler *Reconciler) {
		reconciler.subReconcilers = append(reconciler.subReconcilers, NewStatusReconciler(reconciler.client))
	}
}
