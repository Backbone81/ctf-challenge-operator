package apikey_test

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/backbone81/ctf-challenge-operator/api/v1alpha1"
	"github.com/backbone81/ctf-challenge-operator/internal/controller/apikey"
	"github.com/backbone81/ctf-challenge-operator/internal/utils"
)

var _ = Describe("APIKey Reconciler", func() {
	var reconciler *apikey.Reconciler

	BeforeEach(func() {
		reconciler = apikey.NewReconciler(k8sClient, apikey.WithDefaultReconcilers())
	})

	AfterEach(func() {
		var apiKeyList v1alpha1.APIKeyList
		Expect(k8sClient.List(ctx, &apiKeyList)).To(Succeed())

		for _, apiKey := range apiKeyList.Items {
			Expect(k8sClient.Delete(ctx, &apiKey)).To(Succeed())
		}
	})

	It("should successfully reconcile the resource", func() {
		apiKey := v1alpha1.APIKey{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test",
				Namespace:    "default",
			},
			Spec: v1alpha1.APIKeySpec{},
		}
		Expect(k8sClient.Create(ctx, &apiKey)).To(Succeed())

		result, err := reconciler.Reconcile(ctx, utils.RequestFromObject(&apiKey))
		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(BeZero())
	})
})
