package challengeinstance

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/backbone81/ctf-challenge-operator/api/v1alpha1"
	"github.com/backbone81/ctf-challenge-operator/internal/utils"
)

// NamespaceReconciler is responsible for creating the namespace for the challenge instance.
type NamespaceReconciler struct {
	utils.DefaultSubReconciler
}

func NewNamespaceReconciler(client client.Client) *NamespaceReconciler {
	return &NamespaceReconciler{
		DefaultSubReconciler: utils.NewDefaultSubReconciler(client),
	}
}

func (r *NamespaceReconciler) Reconcile(ctx context.Context, challengeInstance *v1alpha1.ChallengeInstance) (ctrl.Result, error) {
	namespace, err := r.getNamespace(ctx, challengeInstance)
	if err != nil {
		return ctrl.Result{}, err
	}

	if !challengeInstance.DeletionTimestamp.IsZero() {
		return r.reconcileOnDelete(ctx, namespace)
	}

	if namespace == nil {
		desiredSpec := r.getDesiredNamespaceSpec(challengeInstance)
		return r.reconcileOnCreate(ctx, desiredSpec)
	}
	return ctrl.Result{}, nil
}

func (r *NamespaceReconciler) reconcileOnCreate(ctx context.Context, desiredSpec *corev1.Namespace) (ctrl.Result, error) {
	if err := r.GetClient().Create(ctx, desiredSpec); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *NamespaceReconciler) reconcileOnDelete(ctx context.Context, currentSpec *corev1.Namespace) (ctrl.Result, error) {
	if currentSpec == nil {
		return ctrl.Result{}, nil
	}

	if err := r.GetClient().Delete(ctx, currentSpec); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *NamespaceReconciler) getNamespace(ctx context.Context, challengeInstance *v1alpha1.ChallengeInstance) (*corev1.Namespace, error) {
	var namespace corev1.Namespace
	if err := r.GetClient().Get(ctx, client.ObjectKey{
		Name: challengeInstance.Name,
	}, &namespace); err != nil {
		return nil, client.IgnoreNotFound(err)
	}
	return &namespace, nil
}

func (r *NamespaceReconciler) getDesiredNamespaceSpec(challengeInstance *v1alpha1.ChallengeInstance) *corev1.Namespace {
	return &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: challengeInstance.Name,
		},
	}
}
