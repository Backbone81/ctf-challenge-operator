package apikey_test

import (
	"context"
	"testing"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/backbone81/ctf-challenge-operator/api/v1alpha1"
	"github.com/backbone81/ctf-challenge-operator/internal/testutils"
)

var (
	testEnv   *envtest.Environment
	k8sClient client.Client
)

func TestReconciler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "APIKey Suite")
}

var _ = BeforeSuite(func() {
	testEnv, k8sClient = testutils.SetupTestEnv()
})

var _ = AfterSuite(func() {
	Expect(testEnv.Stop()).To(Succeed())
})

func DeleteAllInstances(ctx context.Context) {
	var apiKeyList v1alpha1.APIKeyList
	Expect(k8sClient.List(ctx, &apiKeyList)).To(Succeed())

	for _, apiKey := range apiKeyList.Items {
		Expect(k8sClient.Delete(ctx, &apiKey)).To(Succeed())
	}
}
