package apikey_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/backbone81/ctf-challenge-operator/api/v1alpha1"

	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/backbone81/ctf-challenge-operator/internal/utils"
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
	Expect(utils.MoveToProjectRoot()).To(Succeed())

	logger := zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true))
	ctrllog.SetLogger(logger)

	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,
	}
	cfg, err := testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())
	k8sClient = utils.NewLoggingClient(k8sClient, logger)
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
