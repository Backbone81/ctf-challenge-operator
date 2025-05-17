package apikey_test

import (
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/backbone81/ctf-challenge-operator/api/v1alpha1"
	"github.com/backbone81/ctf-challenge-operator/internal/controller/apikey"
	"github.com/backbone81/ctf-challenge-operator/internal/testutils"
	"github.com/backbone81/ctf-challenge-operator/internal/utils"
)

var _ = Describe("StatusReconciler", func() {
	var reconciler *utils.Reconciler[*v1alpha1.APIKey]

	BeforeEach(func() {
		reconciler = apikey.NewReconciler(k8sClient, apikey.WithStatusReconciler())
	})

	AfterEach(func(ctx SpecContext) {
		DeleteAllInstances(ctx)
	})

	Context("expiration time", func() {
		It("should set the default expiration time when not set", func(ctx SpecContext) {
			By("prepare test with all preconditions")
			instance := v1alpha1.APIKey{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "test-",
					Namespace:    corev1.NamespaceDefault,
				},
			}
			Expect(k8sClient.Create(ctx, &instance)).To(Succeed())
			Expect(instance.Status.ExpirationTimestamp).To(BeZero())

			By("run the reconciler")
			result, err := reconciler.Reconcile(ctx, testutils.RequestFromObject(&instance))
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(BeZero())

			By("verify all postconditions")
			Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(&instance), &instance)).To(Succeed())
			Expect(instance.Status.ExpirationTimestamp.Time).To(BeTemporally(
				"~",
				time.Now().Add(time.Duration(apikey.DefaultExpirationSeconds)*time.Second),
				testutils.DurationEpsilon,
			))
		})

		It("should set the custom expiration time when set", func(ctx SpecContext) {
			By("prepare test with all preconditions")
			customExpirationSeconds := int64(120)
			instance := v1alpha1.APIKey{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "test-",
					Namespace:    corev1.NamespaceDefault,
				},
				Spec: v1alpha1.APIKeySpec{
					ExpirationSeconds: ptr.To(customExpirationSeconds),
				},
			}
			Expect(k8sClient.Create(ctx, &instance)).To(Succeed())
			Expect(instance.Status.ExpirationTimestamp).To(BeZero())

			By("run the reconciler")
			result, err := reconciler.Reconcile(ctx, testutils.RequestFromObject(&instance))
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(BeZero())

			By("verify all postconditions")
			Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(&instance), &instance)).To(Succeed())
			Expect(instance.Status.ExpirationTimestamp.Time).To(BeTemporally(
				"~",
				time.Now().Add(time.Duration(customExpirationSeconds)*time.Second),
				testutils.DurationEpsilon,
			))
		})

		It("should not overwrite the expiration time when already set", func(ctx SpecContext) {
			By("prepare test with all preconditions")
			instance := v1alpha1.APIKey{
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
			result, err := reconciler.Reconcile(ctx, testutils.RequestFromObject(&instance))
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(BeZero())

			By("verify all postconditions")
			Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(&instance), &instance)).To(Succeed())
			Expect(instance.Status.ExpirationTimestamp.Time).To(BeTemporally(
				"~",
				customExpirationTimestamp.Time,
				testutils.DurationEpsilon,
			))
		})

		It("should not set the expiration time when already deleted", func(ctx SpecContext) {
			By("prepare test with all preconditions")
			instance := v1alpha1.APIKey{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "test-",
					Namespace:    corev1.NamespaceDefault,
					Finalizers: []string{
						testutils.DoNotDeleteFinalizerName,
					},
				},
			}
			Expect(k8sClient.Create(ctx, &instance)).To(Succeed())
			Expect(k8sClient.Delete(ctx, &instance)).To(Succeed())
			Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(&instance), &instance)).To(Succeed())
			Expect(instance.Status.ExpirationTimestamp).To(BeZero())

			By("run the reconciler")
			result, err := reconciler.Reconcile(ctx, testutils.RequestFromObject(&instance))
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(BeZero())

			By("verify all postconditions")
			Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(&instance), &instance)).To(Succeed())
			Expect(instance.Status.ExpirationTimestamp).To(BeZero())
		})
	})

	Context("API key", func() {
		It("should generate an API key when not set", func(ctx SpecContext) {
			By("prepare test with all preconditions")
			instance := v1alpha1.APIKey{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "test-",
					Namespace:    corev1.NamespaceDefault,
				},
			}
			Expect(k8sClient.Create(ctx, &instance)).To(Succeed())
			Expect(instance.Status.Key).To(BeZero())

			By("run the reconciler")
			result, err := reconciler.Reconcile(ctx, testutils.RequestFromObject(&instance))
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(BeZero())

			By("verify all postconditions")
			Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(&instance), &instance)).To(Succeed())
			Expect(instance.Status.Key).ToNot(BeZero())
			Expect(instance.Status.Key).To(HaveLen(64))
		})

		It("should not generate an API key when already set", func(ctx SpecContext) {
			By("prepare test with all preconditions")
			instance := v1alpha1.APIKey{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "test-",
					Namespace:    corev1.NamespaceDefault,
				},
			}
			Expect(k8sClient.Create(ctx, &instance)).To(Succeed())
			instance.Status.Key = "abc"
			Expect(k8sClient.Status().Update(ctx, &instance)).To(Succeed())
			Expect(instance.Status.Key).ToNot(BeZero())

			By("run the reconciler")
			result, err := reconciler.Reconcile(ctx, testutils.RequestFromObject(&instance))
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(BeZero())

			By("verify all postconditions")
			Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(&instance), &instance)).To(Succeed())
			Expect(instance.Status.Key).To(Equal("abc"))
		})

		It("should not generate two identical API keys", func(ctx SpecContext) {
			By("prepare test with all preconditions")
			instance1 := v1alpha1.APIKey{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "test-",
					Namespace:    corev1.NamespaceDefault,
				},
			}
			Expect(k8sClient.Create(ctx, &instance1)).To(Succeed())
			Expect(instance1.Status.Key).To(BeZero())

			instance2 := v1alpha1.APIKey{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "test-",
					Namespace:    corev1.NamespaceDefault,
				},
			}
			Expect(k8sClient.Create(ctx, &instance2)).To(Succeed())
			Expect(instance2.Status.Key).To(BeZero())

			By("run the reconciler")
			result, err := reconciler.Reconcile(ctx, testutils.RequestFromObject(&instance1))
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(BeZero())

			result, err = reconciler.Reconcile(ctx, testutils.RequestFromObject(&instance2))
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(BeZero())

			By("verify all postconditions")
			Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(&instance1), &instance1)).To(Succeed())
			Expect(instance1.Status.Key).ToNot(BeZero())
			Expect(instance1.Status.Key).To(HaveLen(64))

			Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(&instance2), &instance2)).To(Succeed())
			Expect(instance2.Status.Key).ToNot(BeZero())
			Expect(instance2.Status.Key).To(HaveLen(64))

			Expect(instance1.Status.Key).ToNot(Equal(instance2.Status.Key))
		})
	})
})
