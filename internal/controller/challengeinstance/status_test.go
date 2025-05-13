package challengeinstance_test

import (
	"time"

	"github.com/backbone81/ctf-challenge-operator/internal/controller/challengeinstance"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/backbone81/ctf-challenge-operator/api/v1alpha1"
	"github.com/backbone81/ctf-challenge-operator/internal/utils"
)

var _ = Describe("StatusReconciler", func() {
	var reconciler *challengeinstance.Reconciler

	BeforeEach(func() {
		reconciler = challengeinstance.NewReconciler(k8sClient, challengeinstance.WithStatusReconciler())
	})

	AfterEach(func(ctx SpecContext) {
		DeleteAllInstances(ctx)
	})

	It("should set the default expiration time when not set", func(ctx SpecContext) {
		By("prepare test with all preconditions")
		instance := v1alpha1.ChallengeInstance{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-",
				Namespace:    corev1.NamespaceDefault,
			},
		}
		Expect(k8sClient.Create(ctx, &instance)).To(Succeed())
		Expect(instance.Status.ExpirationTimestamp).To(BeZero())

		By("run the reconciler")
		result, err := reconciler.Reconcile(ctx, utils.RequestFromObject(&instance))
		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(BeZero())

		By("verify all postconditions")
		Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(&instance), &instance)).To(Succeed())
		Expect(instance.Status.ExpirationTimestamp.Time).To(BeTemporally(
			"~",
			time.Now().Add(time.Duration(challengeinstance.DefaultExpirationSeconds)*time.Second),
			time.Second,
		))
	})

	It("should set the custom expiration time when set", func(ctx SpecContext) {
		By("prepare test with all preconditions")
		customExpirationSeconds := int64(120)
		instance := v1alpha1.ChallengeInstance{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-",
				Namespace:    corev1.NamespaceDefault,
			},
			Spec: v1alpha1.ChallengeInstanceSpec{
				ExpirationSeconds: ptr.To(customExpirationSeconds),
			},
		}
		Expect(k8sClient.Create(ctx, &instance)).To(Succeed())
		Expect(instance.Status.ExpirationTimestamp).To(BeZero())

		By("run the reconciler")
		result, err := reconciler.Reconcile(ctx, utils.RequestFromObject(&instance))
		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(BeZero())

		By("verify all postconditions")
		Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(&instance), &instance)).To(Succeed())
		Expect(instance.Status.ExpirationTimestamp.Time).To(BeTemporally(
			"~",
			time.Now().Add(time.Duration(customExpirationSeconds)*time.Second),
			time.Second,
		))
	})

	It("should not overwrite the expiration time when already set", func(ctx SpecContext) {
		By("prepare test with all preconditions")
		instance := v1alpha1.ChallengeInstance{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-",
				Namespace:    corev1.NamespaceDefault,
			},
		}
		Expect(k8sClient.Create(ctx, &instance)).To(Succeed())

		customExpirationTimestamp := metav1.NewTime(time.Now().Add(3 * time.Minute))
		instance.Status.ExpirationTimestamp = customExpirationTimestamp
		Expect(k8sClient.Status().Update(ctx, &instance)).To(Succeed())
		Expect(instance.Status.ExpirationTimestamp).ToNot(BeZero())

		By("run the reconciler")
		result, err := reconciler.Reconcile(ctx, utils.RequestFromObject(&instance))
		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(BeZero())

		By("verify all postconditions")
		Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(&instance), &instance)).To(Succeed())
		Expect(instance.Status.ExpirationTimestamp.Time).To(BeTemporally(
			"~",
			customExpirationTimestamp.Time,
			time.Second,
		))
	})

	It("should not set the expiration time when already deleted", func(ctx SpecContext) {
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
		Expect(instance.Status.ExpirationTimestamp).To(BeZero())

		By("run the reconciler")
		result, err := reconciler.Reconcile(ctx, utils.RequestFromObject(&instance))
		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(BeZero())

		By("verify all postconditions")
		Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(&instance), &instance)).To(Succeed())
		Expect(instance.Status.ExpirationTimestamp).To(BeZero())
	})
})
