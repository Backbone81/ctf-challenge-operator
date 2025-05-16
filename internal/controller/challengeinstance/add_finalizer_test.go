package challengeinstance_test

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/backbone81/ctf-challenge-operator/api/v1alpha1"
	"github.com/backbone81/ctf-challenge-operator/internal/controller/challengeinstance"
	"github.com/backbone81/ctf-challenge-operator/internal/utils"
)

var _ = Describe("AddFinalizerReconciler", func() {
	var reconciler *utils.Reconciler[*v1alpha1.ChallengeInstance]

	BeforeEach(func() {
		reconciler = challengeinstance.NewReconciler(k8sClient, challengeinstance.WithAddFinalizerReconciler())
	})

	AfterEach(func(ctx SpecContext) {
		DeleteAllInstances(ctx)
	})

	It("should successfully add the finalizer", func(ctx SpecContext) {
		By("prepare test with all preconditions")
		instance := v1alpha1.ChallengeInstance{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-",
				Namespace:    corev1.NamespaceDefault,
			},
		}
		Expect(k8sClient.Create(ctx, &instance)).To(Succeed())
		Expect(instance.DeletionTimestamp.IsZero()).To(BeTrue())
		Expect(controllerutil.ContainsFinalizer(&instance, challengeinstance.FinalizerName)).To(BeFalse())

		By("run the reconciler")
		result, err := reconciler.Reconcile(ctx, utils.RequestFromObject(&instance))
		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(BeZero())

		By("verify all postconditions")
		Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(&instance), &instance)).To(Succeed())
		Expect(controllerutil.ContainsFinalizer(&instance, challengeinstance.FinalizerName)).To(BeTrue())
	})

	It("should succeed if the finalizer already exists", func(ctx SpecContext) {
		By("prepare test with all preconditions")
		instance := v1alpha1.ChallengeInstance{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-",
				Namespace:    corev1.NamespaceDefault,
				Finalizers: []string{
					challengeinstance.FinalizerName,
				},
			},
		}
		Expect(k8sClient.Create(ctx, &instance)).To(Succeed())
		Expect(instance.DeletionTimestamp.IsZero()).To(BeTrue())
		Expect(controllerutil.ContainsFinalizer(&instance, challengeinstance.FinalizerName)).To(BeTrue())

		By("run the reconciler")
		result, err := reconciler.Reconcile(ctx, utils.RequestFromObject(&instance))
		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(BeZero())

		By("verify all postconditions")
		Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(&instance), &instance)).To(Succeed())
		Expect(controllerutil.ContainsFinalizer(&instance, challengeinstance.FinalizerName)).To(BeTrue())
	})

	It("should not add the finalizer when being deleted", func(ctx SpecContext) {
		By("prepare test with all preconditions")
		instance := v1alpha1.ChallengeInstance{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-",
				Namespace:    corev1.NamespaceDefault,
				Finalizers: []string{
					utils.DoNotDeleteFinalizerName,
				},
			},
		}
		Expect(k8sClient.Create(ctx, &instance)).To(Succeed())
		Expect(k8sClient.Delete(ctx, &instance)).To(Succeed())
		Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(&instance), &instance)).To(Succeed())
		Expect(instance.DeletionTimestamp.IsZero()).To(BeFalse())
		Expect(controllerutil.ContainsFinalizer(&instance, challengeinstance.FinalizerName)).To(BeFalse())

		By("run the reconciler")
		result, err := reconciler.Reconcile(ctx, utils.RequestFromObject(&instance))
		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(BeZero())

		By("verify all postconditions")
		Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(&instance), &instance)).To(Succeed())
		Expect(controllerutil.ContainsFinalizer(&instance, challengeinstance.FinalizerName)).To(BeFalse())
	})
})
