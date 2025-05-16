package challengeinstance_test

import (
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/backbone81/ctf-challenge-operator/api/v1alpha1"
	"github.com/backbone81/ctf-challenge-operator/internal/controller/challengeinstance"
	"github.com/backbone81/ctf-challenge-operator/internal/utils"
)

var _ = Describe("Reconciler", func() {
	var reconciler *utils.Reconciler[*v1alpha1.ChallengeInstance]

	BeforeEach(func() {
		reconciler = challengeinstance.NewReconciler(k8sClient, challengeinstance.WithDefaultReconcilers(record.NewFakeRecorder(5)))
	})

	AfterEach(func(ctx SpecContext) {
		DeleteAllInstances(ctx)
	})

	It("should successfully reconcile the resource", func(ctx SpecContext) {
		By("prepare test with all preconditions")
		configMapName := GenerateName("test-")
		configMap := corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name: configMapName,
			},
		}
		configMapRaw, err := ToRaw(&configMap)
		Expect(err).ToNot(HaveOccurred())

		description := v1alpha1.ChallengeDescription{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-",
				Namespace:    corev1.NamespaceDefault,
			},
			Spec: v1alpha1.ChallengeDescriptionSpec{
				Title:       "test",
				Description: "test",
				Manifests: []runtime.RawExtension{
					{
						Raw: configMapRaw,
					},
				},
			},
		}
		Expect(k8sClient.Create(ctx, &description)).To(Succeed())

		instance := v1alpha1.ChallengeInstance{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-",
				Namespace:    corev1.NamespaceDefault,
			},
			Spec: v1alpha1.ChallengeInstanceSpec{
				ChallengeDescriptionName: description.Name,
			},
		}
		Expect(k8sClient.Create(ctx, &instance)).To(Succeed())
		Expect(instance.Spec.ChallengeDescriptionName).ToNot(BeZero())

		By("run the reconciler")
		result, err := reconciler.Reconcile(ctx, utils.RequestFromObject(&instance))
		Expect(err).ToNot(HaveOccurred())
		Expect(result.RequeueAfter).To(BeNumerically(
			"~",
			time.Duration(challengeinstance.DefaultExpirationSeconds)*time.Second,
			utils.DurationEpsilon,
		))

		By("verify all postconditions")
		var namespace corev1.Namespace
		Expect(k8sClient.Get(ctx, client.ObjectKey{
			Name: instance.Name,
		}, &namespace)).To(Succeed())
		Expect(k8sClient.Get(ctx, client.ObjectKey{
			Name:      configMap.Name,
			Namespace: instance.Name,
		}, &configMap)).To(Succeed())
	})
})
