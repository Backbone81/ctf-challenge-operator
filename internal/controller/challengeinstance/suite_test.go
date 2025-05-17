package challengeinstance_test

import (
	"context"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
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
	RunSpecs(t, "ChallengeInstance Suite")
}

var _ = BeforeSuite(func() {
	testEnv, k8sClient = testutils.SetupTestEnv()
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
