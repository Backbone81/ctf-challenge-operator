package challengeinstance_test

import (
	"github.com/backbone81/ctf-challenge-operator/api/v1alpha1"
	"github.com/backbone81/ctf-challenge-operator/internal/controller/challengeinstance"
	"github.com/backbone81/ctf-challenge-operator/internal/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("NamespaceReconciler", func() {
	var reconciler *challengeinstance.Reconciler

	BeforeEach(func() {
		reconciler = challengeinstance.NewReconciler(k8sClient, challengeinstance.WithNamespaceReconciler())
	})

	AfterEach(func() {
		DeleteAllInstances()
	})

	It("should successfully create the namespace", func() {
		By("prepare test with all preconditions")
		instance := v1alpha1.ChallengeInstance{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-",
				Namespace:    corev1.NamespaceDefault,
			},
		}
		Expect(k8sClient.Create(ctx, &instance)).To(Succeed())
		var namespace corev1.Namespace
		Expect(k8sClient.Get(ctx, client.ObjectKey{
			Name: instance.Name,
		}, &namespace)).ToNot(Succeed())

		By("run the reconciler")
		result, err := reconciler.Reconcile(ctx, utils.RequestFromObject(&instance))
		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(BeZero())

		By("verify all postconditions")
		Expect(k8sClient.Get(ctx, client.ObjectKey{
			Name: instance.Name,
		}, &namespace)).To(Succeed())
	})

	It("should succeed if the namespace already exists", func() {
		By("prepare test with all preconditions")
		instance := v1alpha1.ChallengeInstance{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-",
				Namespace:    corev1.NamespaceDefault,
			},
		}
		Expect(k8sClient.Create(ctx, &instance)).To(Succeed())
		namespace := corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: instance.Name,
			},
		}
		Expect(k8sClient.Create(ctx, &namespace)).To(Succeed())

		By("run the reconciler")
		result, err := reconciler.Reconcile(ctx, utils.RequestFromObject(&instance))
		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(BeZero())

		By("verify all postconditions")
		Expect(k8sClient.Get(ctx, client.ObjectKey{
			Name: instance.Name,
		}, &namespace)).To(Succeed())
	})

	It("should delete the namespace on deletion", func() {
		By("prepare test with all preconditions")
		instance := v1alpha1.ChallengeInstance{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-",
				Namespace:    corev1.NamespaceDefault,
				Finalizers: []string{
					DoNotDeleteFinalizerName,
				},
			},
		}
		Expect(k8sClient.Create(ctx, &instance)).To(Succeed())
		Expect(k8sClient.Delete(ctx, &instance)).To(Succeed())
		Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(&instance), &instance)).To(Succeed())
		namespace := corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: instance.Name,
			},
		}
		Expect(k8sClient.Create(ctx, &namespace)).To(Succeed())

		By("run the reconciler")
		result, err := reconciler.Reconcile(ctx, utils.RequestFromObject(&instance))
		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(BeZero())

		By("verify all postconditions")
		Expect(k8sClient.Get(ctx, client.ObjectKey{
			Name: instance.Name,
		}, &namespace)).To(Succeed())
		Expect(namespace.DeletionTimestamp.IsZero()).To(BeFalse())
	})
})
