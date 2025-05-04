package utils

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func RequestFromObject(obj client.Object) reconcile.Request {
	return reconcile.Request{
		NamespacedName: client.ObjectKeyFromObject(obj),
	}
}
