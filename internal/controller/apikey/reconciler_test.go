package apikey_test

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

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
		Expect(result).ToNot(BeZero())
	})

	It("should create an API key with expiration", func() {
		apiKey := v1alpha1.APIKey{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test",
				Namespace:    "default",
			},
		}
		Expect(k8sClient.Create(ctx, &apiKey)).To(Succeed())

		Expect(apiKey.Status.Key).To(BeZero())
		Expect(apiKey.Status.ExpirationTimestamp).To(BeZero())

		result, err := reconciler.Reconcile(ctx, utils.RequestFromObject(&apiKey))
		Expect(err).ToNot(HaveOccurred())
		Expect(result).ToNot(BeZero())

		Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(&apiKey), &apiKey)).To(Succeed())

		Expect(apiKey.Status.Key).ToNot(BeZero())
		Expect(apiKey.Status.Key).To(HaveLen(64))
		Expect(apiKey.Status.ExpirationTimestamp).ToNot(BeZero())
	})

	It("should create different API keys", func() {
		apiKey1 := v1alpha1.APIKey{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test",
				Namespace:    "default",
			},
		}
		Expect(k8sClient.Create(ctx, &apiKey1)).To(Succeed())

		result, err := reconciler.Reconcile(ctx, utils.RequestFromObject(&apiKey1))
		Expect(err).ToNot(HaveOccurred())
		Expect(result).ToNot(BeZero())

		Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(&apiKey1), &apiKey1)).To(Succeed())

		apiKey2 := v1alpha1.APIKey{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test",
				Namespace:    "default",
			},
		}
		Expect(k8sClient.Create(ctx, &apiKey2)).To(Succeed())

		result, err = reconciler.Reconcile(ctx, utils.RequestFromObject(&apiKey2))
		Expect(err).ToNot(HaveOccurred())
		Expect(result).ToNot(BeZero())

		Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(&apiKey2), &apiKey2)).To(Succeed())

		Expect(apiKey1.Status.Key).ToNot(Equal(apiKey2.Status.Key))
	})

	It("should delete an expired API key", func() {
		apiKey := v1alpha1.APIKey{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test",
				Namespace:    "default",
			},
		}
		Expect(k8sClient.Create(ctx, &apiKey)).To(Succeed())

		apiKey.Status.Key = "foo"
		apiKey.Status.ExpirationTimestamp = metav1.NewTime(time.Now().Add(-1 * time.Second))
		Expect(k8sClient.Status().Update(ctx, &apiKey)).To(Succeed())

		result, err := reconciler.Reconcile(ctx, utils.RequestFromObject(&apiKey))
		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(BeZero())

		Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(&apiKey), &apiKey)).To(MatchError(ContainSubstring("not found")))
	})

	It("should not overwrite existing API key", func() {
		apiKey := v1alpha1.APIKey{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test",
				Namespace:    "default",
			},
		}
		Expect(k8sClient.Create(ctx, &apiKey)).To(Succeed())

		expectedKey := "foo"
		expectedTimestamp := metav1.NewTime(time.Date(2035, 5, 4, 16, 38, 22, 0, time.Local))

		apiKey.Status.Key = expectedKey
		apiKey.Status.ExpirationTimestamp = expectedTimestamp
		Expect(k8sClient.Status().Update(ctx, &apiKey)).To(Succeed())

		result, err := reconciler.Reconcile(ctx, utils.RequestFromObject(&apiKey))
		Expect(err).ToNot(HaveOccurred())
		Expect(result).ToNot(BeZero())

		Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(&apiKey), &apiKey)).To(Succeed())

		Expect(apiKey.Status.Key).To(Equal(expectedKey))
		Expect(apiKey.Status.ExpirationTimestamp).To(Equal(expectedTimestamp))
	})
})
