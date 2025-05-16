package challengeinstance_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/backbone81/ctf-challenge-operator/api/v1alpha1"
	"github.com/backbone81/ctf-challenge-operator/internal/utils"
)

var (
	testEnv   *envtest.Environment
	k8sClient client.Client
)

func TestReconciler(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "ChallengeInstance Suite")
}

var _ = BeforeSuite(func() {
	Expect(utils.MoveToProjectRoot()).To(Succeed())

	logger := zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true))
	ctrllog.SetLogger(logger)

	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{"manifests/ctf-challenge-operator-crd.yaml"},
		ErrorIfCRDPathMissing: true,
		BinaryAssetsDirectory: "bin",
	}
	cfg, err := testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	k8sClient, err = client.New(cfg, client.Options{Scheme: clientgoscheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())
	k8sClient = utils.NewLoggingClient(k8sClient, logger)
})

var _ = AfterSuite(func() {
	Expect(testEnv.Stop()).To(Succeed())
})

func DeleteAllInstances(ctx context.Context) {
	var challengeInstanceList v1alpha1.ChallengeInstanceList
	Expect(k8sClient.List(ctx, &challengeInstanceList)).To(Succeed())

	for _, challengeInstance := range challengeInstanceList.Items {
		Expect(k8sClient.Delete(ctx, &challengeInstance)).To(Succeed())
	}
}

// GenerateName simulates the behavior of Kubernetes GenerateName for test situations where we need to know the name
// beforehand.
func GenerateName(prefix string) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	suffix := make([]byte, 5)
	for i := range suffix {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		Expect(err).ToNot(HaveOccurred())

		suffix[i] = charset[n.Int64()]
	}
	return fmt.Sprintf("%s%s", prefix, string(suffix))
}

// ToRaw converts a Kubernetes object into its JSON representation.
func ToRaw(obj client.Object) ([]byte, error) {
	codecFactory := serializer.NewCodecFactory(clientgoscheme.Scheme)
	encoder := codecFactory.LegacyCodec(getGroupVersionKind(obj).GroupVersion())
	return runtime.Encode(encoder, obj)
}

func getGroupVersionKind(obj client.Object) schema.GroupVersionKind {
	if !obj.GetObjectKind().GroupVersionKind().Empty() {
		return obj.GetObjectKind().GroupVersionKind()
	}
	gvks, _, err := clientgoscheme.Scheme.ObjectKinds(obj)
	Expect(err).ToNot(HaveOccurred())
	Expect(gvks).To(HaveLen(1))
	return gvks[0]
}
