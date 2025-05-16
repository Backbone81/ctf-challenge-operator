package apikey

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/backbone81/ctf-challenge-operator/api/v1alpha1"
	"github.com/backbone81/ctf-challenge-operator/internal/utils"
)

const (
	DefaultExpirationSeconds = int64(12 * 60 * 60) // default is 12 hours
)

// StatusReconciler is responsible for reconciling the status of the APIKey resource.
type StatusReconciler struct {
	utils.DefaultSubReconciler
}

// NewStatusReconciler creates a new sub-reconciler instance. The reconciler is initialized with the given client.
func NewStatusReconciler(client client.Client) *StatusReconciler {
	return &StatusReconciler{
		DefaultSubReconciler: utils.NewDefaultSubReconciler(client),
	}
}

// Reconcile is the main reconciler function.
func (r *StatusReconciler) Reconcile(ctx context.Context, apiKey *v1alpha1.APIKey) (ctrl.Result, error) {
	if !apiKey.DeletionTimestamp.IsZero() {
		// We do not update the status when the resource is already being deleted.
		return ctrl.Result{}, nil
	}

	updateStatus := false

	// generate a key if needed
	if len(apiKey.Status.Key) == 0 {
		key, err := GenerateAPIKey()
		if err != nil {
			return ctrl.Result{}, err
		}
		apiKey.Status.Key = key
		updateStatus = true
	}

	// calculate expiration timestamp
	if apiKey.Status.ExpirationTimestamp.IsZero() {
		expirationSeconds := DefaultExpirationSeconds
		if apiKey.Spec.ExpirationSeconds != nil {
			expirationSeconds = *apiKey.Spec.ExpirationSeconds
		}
		apiKey.Status.ExpirationTimestamp = metav1.NewTime(time.Now().Add(time.Duration(expirationSeconds) * time.Second))
		updateStatus = true
	}

	if updateStatus {
		if err := r.GetClient().Status().Update(ctx, apiKey); err != nil {
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

const APIKeyLength = 32 // 256-bit key

// GenerateAPIKey generates a cryptographically secure random API key.
func GenerateAPIKey() (string, error) {
	apiKey := make([]byte, APIKeyLength)
	if _, err := rand.Read(apiKey); err != nil {
		return "", fmt.Errorf("reading crypto rand: %w", err)
	}
	return hex.EncodeToString(apiKey), nil
}
