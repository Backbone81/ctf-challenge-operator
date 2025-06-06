package challengeinstance

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/backbone81/ctf-challenge-operator/api/v1alpha1"
	"github.com/backbone81/ctf-challenge-operator/internal/utils"
)

// ManifestsReconciler is responsible for creating the manifests for the challenge instance.
type ManifestsReconciler struct {
	utils.DefaultSubReconciler
	recorder record.EventRecorder
}

func NewManifestsReconciler(client client.Client, recorder record.EventRecorder) *ManifestsReconciler {
	return &ManifestsReconciler{
		DefaultSubReconciler: utils.NewDefaultSubReconciler(client),
		recorder:             recorder,
	}
}

func (r *ManifestsReconciler) Reconcile(ctx context.Context, challengeInstance *v1alpha1.ChallengeInstance) (ctrl.Result, error) {
	if !challengeInstance.DeletionTimestamp.IsZero() {
		// We do not create manifests when the resource is already being deleted.
		return ctrl.Result{}, nil
	}

	var challengeDescription v1alpha1.ChallengeDescription
	if err := r.GetClient().Get(ctx, client.ObjectKey{
		Namespace: challengeInstance.Namespace,
		Name:      challengeInstance.Spec.ChallengeDescriptionName,
	}, &challengeDescription); err != nil {
		r.recorder.Eventf(
			challengeInstance,
			corev1.EventTypeWarning,
			"Creating",
			"ChallengeDescription could not be found at %s/%s: %s",
			challengeInstance.Namespace,
			challengeInstance.Spec.ChallengeDescriptionName,
			err,
		)
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

		if result, err := r.reconcileManifest(ctx, challengeInstance, &desiredSpec); err != nil || !result.IsZero() {
			return result, err
		}
	}
	return ctrl.Result{}, nil
}

func (r *ManifestsReconciler) reconcileManifest(ctx context.Context, challengeInstance *v1alpha1.ChallengeInstance, desiredSpec *unstructured.Unstructured) (ctrl.Result, error) {
	currentSpec, err := r.getCurrentSpec(ctx, desiredSpec)
	if err != nil {
		return ctrl.Result{}, err
	}

	if currentSpec == nil {
		return r.reconcileManifestOnCreate(ctx, challengeInstance, desiredSpec)
	}
	return r.reconcileManifestOnUpdate(ctx, challengeInstance, desiredSpec, currentSpec)
}

func (r *ManifestsReconciler) reconcileManifestOnCreate(ctx context.Context, challengeInstance *v1alpha1.ChallengeInstance, desiredSpec *unstructured.Unstructured) (ctrl.Result, error) {
	if err := r.GetClient().Create(ctx, desiredSpec); err != nil {
		r.recorder.Eventf(
			challengeInstance,
			corev1.EventTypeWarning,
			"Creating",
			"Failed to create %s at %s/%s: %s",
			desiredSpec.GroupVersionKind(),
			desiredSpec.GetNamespace(),
			desiredSpec.GetName(),
			err,
		)
		return ctrl.Result{}, err
	}
	r.recorder.Eventf(
		challengeInstance,
		corev1.EventTypeNormal,
		"Creating",
		"Created %s at %s/%s",
		desiredSpec.GroupVersionKind(),
		desiredSpec.GetNamespace(),
		desiredSpec.GetName(),
	)
	return ctrl.Result{}, nil
}

func (r *ManifestsReconciler) reconcileManifestOnUpdate(ctx context.Context, challengeInstance *v1alpha1.ChallengeInstance, desiredSpec *unstructured.Unstructured, currentSpec *unstructured.Unstructured) (ctrl.Result, error) {
	if equality.Semantic.DeepDerivative(desiredSpec.Object["spec"], currentSpec.Object["spec"]) {
		// The resources are identical. Nothing to do.
		return ctrl.Result{}, nil
	}

	currentSpec.Object["spec"] = desiredSpec.Object["spec"]
	if err := r.GetClient().Update(ctx, currentSpec); err != nil {
		r.recorder.Eventf(
			challengeInstance,
			corev1.EventTypeWarning,
			"Updating",
			"Failed to update %s at %s/%s: %s",
			desiredSpec.GroupVersionKind(),
			desiredSpec.GetNamespace(),
			desiredSpec.GetName(),
			err,
		)
		return ctrl.Result{}, err
	}
	r.recorder.Eventf(
		challengeInstance,
		corev1.EventTypeNormal,
		"Updating",
		"Updated %s at %s/%s",
		desiredSpec.GroupVersionKind(),
		desiredSpec.GetNamespace(),
		desiredSpec.GetName(),
	)
	return ctrl.Result{}, nil
}

func (r *ManifestsReconciler) getCurrentSpec(ctx context.Context, desiredSpec client.Object) (*unstructured.Unstructured, error) {
	var currentSpec unstructured.Unstructured
	currentSpec.SetGroupVersionKind(desiredSpec.GetObjectKind().GroupVersionKind())
	if err := r.GetClient().Get(ctx, client.ObjectKey{
		Namespace: desiredSpec.GetNamespace(),
		Name:      desiredSpec.GetName(),
	}, &currentSpec); err != nil {
		return nil, client.IgnoreNotFound(err)
	}
	return &currentSpec, nil
}
