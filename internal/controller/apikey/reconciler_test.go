package apikey_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/backbone81/ctf-challenge-operator/api/v1alpha1"
	"github.com/backbone81/ctf-challenge-operator/internal/controller/apikey"
	"github.com/backbone81/ctf-challenge-operator/internal/utils"
)

var _ = Describe("APIKey Reconciler", func() {
	var reconciler *apikey.Reconciler

	BeforeEach(func() {
		reconciler = apikey.NewReconciler(k8sClient, apikey.WithDefaultReconcilers())
	})

	AfterEach(func(ctx SpecContext) {
		DeleteAllInstances(ctx)
	})

	It("should successfully reconcile the resource", func(ctx SpecContext) {
		apiKey := v1alpha1.APIKey{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test",
				Namespace:    "default",
			},
		}
		Expect(k8sClient.Create(ctx, &apiKey)).To(Succeed())

		result, err := reconciler.Reconcile(ctx, utils.RequestFromObject(&apiKey))
		Expect(err).ToNot(HaveOccurred())
		Expect(result).ToNot(BeZero())
	})
})
