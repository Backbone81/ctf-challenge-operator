package challengeinstance

import (
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/backbone81/ctf-challenge-operator/api/v1alpha1"
	"github.com/backbone81/ctf-challenge-operator/internal/utils"
)

var FinalizerName = "ctf.backbone81/challenge-instance"

// +kubebuilder:rbac:groups=core.ctf.backbone81,resources=challengeinstances,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core.ctf.backbone81,resources=challengeinstances/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core.ctf.backbone81,resources=challengeinstances/finalizers,verbs=update

// +kubebuilder:rbac:groups=core.ctf.backbone81,resources=challengedescriptions,verbs=get;list;watch

// +kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=events,verbs=create;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

func NewReconciler(client client.Client, options ...utils.ReconcilerOption[*v1alpha1.ChallengeInstance]) *utils.Reconciler[*v1alpha1.ChallengeInstance] {
	return utils.NewReconciler[*v1alpha1.ChallengeInstance](
		client,
		func() *v1alpha1.ChallengeInstance {
			return &v1alpha1.ChallengeInstance{}
		},
		options...,
	)
}

// WithDefaultReconcilers returns a reconciler option which enables the default sub-reconcilers.
func WithDefaultReconcilers(recorder record.EventRecorder) utils.ReconcilerOption[*v1alpha1.ChallengeInstance] {
	return func(reconciler *utils.Reconciler[*v1alpha1.ChallengeInstance]) {
		WithAddFinalizerReconciler()(reconciler)
		WithStatusReconciler()(reconciler)
		WithNamespaceReconciler()(reconciler)
		WithManifestsReconciler(recorder)(reconciler)
		WithRemoveFinalizerReconciler()(reconciler)

		// The delete reconciler must be last, because the other reconcilers behave differently when the resource is
		// deleted.
		WithDeleteReconciler()(reconciler)
	}
}

func WithStatusReconciler() utils.ReconcilerOption[*v1alpha1.ChallengeInstance] {
	return func(reconciler *utils.Reconciler[*v1alpha1.ChallengeInstance]) {
		reconciler.AppendSubReconciler(NewStatusReconciler(reconciler.GetClient()))
	}
}

func WithDeleteReconciler() utils.ReconcilerOption[*v1alpha1.ChallengeInstance] {
	return func(reconciler *utils.Reconciler[*v1alpha1.ChallengeInstance]) {
		reconciler.AppendSubReconciler(NewDeleteReconciler(reconciler.GetClient()))
	}
}

func WithNamespaceReconciler() utils.ReconcilerOption[*v1alpha1.ChallengeInstance] {
	return func(reconciler *utils.Reconciler[*v1alpha1.ChallengeInstance]) {
		reconciler.AppendSubReconciler(NewNamespaceReconciler(reconciler.GetClient()))
	}
}

func WithAddFinalizerReconciler() utils.ReconcilerOption[*v1alpha1.ChallengeInstance] {
	return func(reconciler *utils.Reconciler[*v1alpha1.ChallengeInstance]) {
		reconciler.AppendSubReconciler(NewAddFinalizerReconciler(reconciler.GetClient()))
	}
}

func WithRemoveFinalizerReconciler() utils.ReconcilerOption[*v1alpha1.ChallengeInstance] {
	return func(reconciler *utils.Reconciler[*v1alpha1.ChallengeInstance]) {
		reconciler.AppendSubReconciler(NewRemoveFinalizerReconciler(reconciler.GetClient()))
	}
}

func WithManifestsReconciler(recorder record.EventRecorder) utils.ReconcilerOption[*v1alpha1.ChallengeInstance] {
	return func(reconciler *utils.Reconciler[*v1alpha1.ChallengeInstance]) {
		reconciler.AppendSubReconciler(NewManifestsReconciler(reconciler.GetClient(), recorder))
	}
}
