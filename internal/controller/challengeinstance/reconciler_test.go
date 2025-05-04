package challengeinstance_test

import (
	"github.com/backbone81/ctf-challenge-operator/internal/controller/challengeinstance"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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
		Expect(result).To(BeZero())
	})
})
