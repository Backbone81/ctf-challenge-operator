package challengeinstance_test

import (
	"time"

	"github.com/backbone81/ctf-challenge-operator/internal/controller/challengeinstance"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/backbone81/ctf-challenge-operator/api/v1alpha1"
	"github.com/backbone81/ctf-challenge-operator/internal/utils"
)

var _ = Describe("DeleteReconciler", func() {
	var reconciler *challengeinstance.Reconciler

	BeforeEach(func() {
		reconciler = challengeinstance.NewReconciler(k8sClient, challengeinstance.WithDeleteReconciler())
	})

	AfterEach(func() {
		DeleteAllInstances()
	})

	It("should delete the instance when expiration is reached", func() {
		By("prepare test with all preconditions")
		instance := v1alpha1.ChallengeInstance{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-",
				Namespace:    corev1.NamespaceDefault,
			},
		}
		Expect(k8sClient.Create(ctx, &instance)).To(Succeed())

		instance.Status.ExpirationTimestamp = metav1.NewTime(time.Now().Add(-time.Minute))
		Expect(k8sClient.Status().Update(ctx, &instance)).To(Succeed())

		Expect(instance.Status.ExpirationTimestamp.Time.Before(time.Now())).To(BeTrue())

		By("run the reconciler")
		result, err := reconciler.Reconcile(ctx, utils.RequestFromObject(&instance))
		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(BeZero())

		By("verify all postconditions")
		Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(&instance), &instance)).To(MatchError(ContainSubstring("not found")))
	})

	It("should not delete the instance when expiration is not reached", func() {
		By("prepare test with all preconditions")
		instance := v1alpha1.ChallengeInstance{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-",
				Namespace:    corev1.NamespaceDefault,
			},
		}
		Expect(k8sClient.Create(ctx, &instance)).To(Succeed())

		instance.Status.ExpirationTimestamp = metav1.NewTime(time.Now().Add(time.Minute))
		Expect(k8sClient.Status().Update(ctx, &instance)).To(Succeed())

		Expect(instance.Status.ExpirationTimestamp.Time.Before(time.Now())).To(BeFalse())

		By("run the reconciler")
		result, err := reconciler.Reconcile(ctx, utils.RequestFromObject(&instance))
		Expect(err).ToNot(HaveOccurred())
		Expect(result).ToNot(BeZero())
		Expect(result.RequeueAfter).To(BeNumerically("~", time.Minute, time.Second))

		By("verify all postconditions")
		Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(&instance), &instance)).To(Succeed())
	})

	It("should not delete the instance when expiration is reached and instance is already deleted", func() {
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

		instance.Status.ExpirationTimestamp = metav1.NewTime(time.Now().Add(-time.Minute))
		Expect(k8sClient.Status().Update(ctx, &instance)).To(Succeed())

		Expect(instance.Status.ExpirationTimestamp.Time.Before(time.Now())).To(BeTrue())

		By("run the reconciler")
		result, err := reconciler.Reconcile(ctx, utils.RequestFromObject(&instance))
		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(BeZero())

		By("verify all postconditions")
		Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(&instance), &instance)).To(Succeed())
	})
})
