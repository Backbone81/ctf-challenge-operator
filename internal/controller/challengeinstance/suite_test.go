package challengeinstance_test

import (
	"context"
	"fmt"
	"math/rand"
	"path/filepath"
	"testing"

	"github.com/backbone81/ctf-challenge-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"

	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/backbone81/ctf-challenge-operator/internal/utils"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
)

var (
	ctx       context.Context
	cancel    context.CancelFunc
	testEnv   *envtest.Environment
	k8sClient client.Client
)

func TestReconciler(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "ChallengeInstance Suite")
}

var _ = BeforeSuite(func() {
	Expect(utils.MoveToProjectRoot()).To(Succeed())
	ctx, cancel = context.WithCancel(context.TODO()) //nolint:fatcontext // This does not lead to fat contexts.

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
	cancel()
	Expect(testEnv.Stop()).To(Succeed())
})

var (
	// DoNotDeleteFinalizerName provides a name for a finalizer which is not used by the reconcilers themselves. This finalizer
	// is used to prevent a resource from being deleted immediately when you want to test situations where you need to
	// inspect the behavior for deletion.
	DoNotDeleteFinalizerName = "ctf.backbone81/do-not-delete"
)

func DeleteAllInstances() {
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
		suffix[i] = charset[rand.Intn(len(charset))]
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
