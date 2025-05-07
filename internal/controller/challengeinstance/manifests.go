package challengeinstance

import (
	"context"

	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/backbone81/ctf-challenge-operator/api/v1alpha1"
)

// ManifestsReconciler is responsible for creating the namespace for the challenge instance.
type ManifestsReconciler struct {
	client client.Client
}

// NewManifestsReconciler creates a new sub-reconciler instance. The reconciler is initialized with the given client.
func NewManifestsReconciler(client client.Client) *ManifestsReconciler {
	return &ManifestsReconciler{
		client: client,
	}
}

// SetupWithManager registers the sub-reconciler with the manager.
func (r *ManifestsReconciler) SetupWithManager(ctrlBuilder *builder.Builder) *builder.Builder {
	return ctrlBuilder
}

// Reconcile is the main reconciler function.
func (r *ManifestsReconciler) Reconcile(ctx context.Context, challengeInstance *v1alpha1.ChallengeInstance) (ctrl.Result, error) {
	var challengeDescription v1alpha1.ChallengeDescription
	if err := r.client.Get(ctx, client.ObjectKey{
		Namespace: challengeInstance.Namespace,
		Name:      challengeInstance.Spec.ChallengeDescription.Name,
	}, &challengeDescription); err != nil {
		return ctrl.Result{}, err
	}

	codecFactory := serializer.NewCodecFactory(clientgoscheme.Scheme)
	decoder := codecFactory.UniversalDeserializer()

	for _, raw := range challengeDescription.Spec.Manifests {
		var desiredSpec unstructured.Unstructured
		if _, _, err := decoder.Decode(raw.Raw, nil, &desiredSpec); err != nil {
			return ctrl.Result{}, err
		}

		// We need to make sure that we overwrite the target namespace to prevent challenge instances from placing
		// workload into unrelated namespaces.
		desiredSpec.SetNamespace(challengeInstance.Name)
		if result, err := r.reconcileManifest(ctx, &desiredSpec); err != nil || !result.IsZero() {
			return result, err
		}
	}
	return ctrl.Result{}, nil
}

func (r *ManifestsReconciler) reconcileManifest(ctx context.Context, desiredSpec *unstructured.Unstructured) (ctrl.Result, error) {
	currentSpec, err := r.getCurrentSpec(ctx, desiredSpec)
	if err != nil {
		return ctrl.Result{}, err
	}

	if currentSpec == nil {
		return r.reconcileManifestOnCreate(ctx, desiredSpec)
	}
	return r.reconcileManifestOnUpdate(ctx, desiredSpec, currentSpec)
}

func (r *ManifestsReconciler) reconcileManifestOnCreate(ctx context.Context, desiredSpec *unstructured.Unstructured) (ctrl.Result, error) {
	if err := r.client.Create(ctx, desiredSpec); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *ManifestsReconciler) reconcileManifestOnUpdate(ctx context.Context, desiredSpec *unstructured.Unstructured, currentSpec *unstructured.Unstructured) (ctrl.Result, error) {
	if equality.Semantic.DeepDerivative(desiredSpec.Object["spec"], currentSpec.Object["spec"]) {
		// The resources are identical. Nothing to do.
		return ctrl.Result{}, nil
	}

	currentSpec.Object["spec"] = desiredSpec.Object["spec"]
	if err := r.client.Update(ctx, currentSpec); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *ManifestsReconciler) getCurrentSpec(ctx context.Context, desiredSpec client.Object) (*unstructured.Unstructured, error) {
	var currentSpec unstructured.Unstructured
	currentSpec.SetGroupVersionKind(desiredSpec.GetObjectKind().GroupVersionKind())
	if err := r.client.Get(ctx, client.ObjectKey{
		Namespace: desiredSpec.GetNamespace(),
		Name:      desiredSpec.GetName(),
	}, &currentSpec); err != nil {
		return nil, client.IgnoreNotFound(err)
	}
	return &currentSpec, nil
}
