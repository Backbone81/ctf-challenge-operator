package utils

import (
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type DefaultSubReconciler struct {
	client client.Client
}

func NewDefaultSubReconciler(client client.Client) DefaultSubReconciler {
	return DefaultSubReconciler{
		client: client,
	}
}

func (r *DefaultSubReconciler) GetClient() client.Client {
	return r.client
}

func (r *DefaultSubReconciler) SetupWithManager(ctrlBuilder *builder.Builder) *builder.Builder {
	return ctrlBuilder
}
