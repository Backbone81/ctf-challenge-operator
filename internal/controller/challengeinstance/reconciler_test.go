package challengeinstance_test

import (
	"github.com/backbone81/ctf-challenge-operator/internal/controller/challengeinstance"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/backbone81/ctf-challenge-operator/api/v1alpha1"
	"github.com/backbone81/ctf-challenge-operator/internal/utils"
)

var _ = Describe("ChallengeInstance Reconciler", func() {
	var reconciler *challengeinstance.Reconciler

	BeforeEach(func() {
		reconciler = challengeinstance.NewReconciler(k8sClient, challengeinstance.WithDefaultReconcilers())
	})

	AfterEach(func() {
		var challengeInstanceList v1alpha1.ChallengeInstanceList
		Expect(k8sClient.List(ctx, &challengeInstanceList)).To(Succeed())

		for _, challengeInstance := range challengeInstanceList.Items {
			Expect(k8sClient.Delete(ctx, &challengeInstance)).To(Succeed())
		}
	})

	It("should successfully reconcile the resource", func() {
		challengeInstance := v1alpha1.ChallengeInstance{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test",
				Namespace:    "default",
			},
			Spec: v1alpha1.ChallengeInstanceSpec{},
		}
		Expect(k8sClient.Create(ctx, &challengeInstance)).To(Succeed())

		result, err := reconciler.Reconcile(ctx, utils.RequestFromObject(&challengeInstance))
		Expect(err).ToNot(HaveOccurred())
		Expect(result).ToNot(BeZero())
	})

	It("should set an expiration", func() {
		challengeInstance := v1alpha1.ChallengeInstance{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test",
				Namespace:    "default",
			},
			Spec: v1alpha1.ChallengeInstanceSpec{},
		}
		Expect(k8sClient.Create(ctx, &challengeInstance)).To(Succeed())

		Expect(challengeInstance.Status.ExpirationTimestamp).To(BeZero())

		result, err := reconciler.Reconcile(ctx, utils.RequestFromObject(&challengeInstance))
		Expect(err).ToNot(HaveOccurred())
		Expect(result).ToNot(BeZero())

		Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(&challengeInstance), &challengeInstance)).To(Succeed())

		Expect(challengeInstance.Status.ExpirationTimestamp).ToNot(BeZero())
	})

	It("should delete an expired challenge instance", func() {
		challengeInstance := v1alpha1.ChallengeInstance{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test",
				Namespace:    "default",
			},
			Spec: v1alpha1.ChallengeInstanceSpec{},
		}
		Expect(k8sClient.Create(ctx, &challengeInstance)).To(Succeed())

		challengeInstance.Status.ExpirationTimestamp = metav1.NewTime(time.Now().Add(-1 * time.Second))
		Expect(k8sClient.Status().Update(ctx, &challengeInstance)).To(Succeed())

		result, err := reconciler.Reconcile(ctx, utils.RequestFromObject(&challengeInstance))
		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(BeZero())

		Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(&challengeInstance), &challengeInstance)).To(MatchError(ContainSubstring("not found")))
	})
})
