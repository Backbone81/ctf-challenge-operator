package challengeinstance_test

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/backbone81/ctf-challenge-operator/api/v1alpha1"
	"github.com/backbone81/ctf-challenge-operator/internal/controller/challengeinstance"
	"github.com/backbone81/ctf-challenge-operator/internal/testutils"
	"github.com/backbone81/ctf-challenge-operator/internal/utils"
)

var _ = Describe("ManifestsReconciler", func() {
	var reconciler *utils.Reconciler[*v1alpha1.ChallengeInstance]

	BeforeEach(func() {
		reconciler = challengeinstance.NewReconciler(k8sClient, challengeinstance.WithManifestsReconciler(record.NewFakeRecorder(5)))
	})

	AfterEach(func(ctx SpecContext) {
		DeleteAllInstances(ctx)
	})

	It("should successfully create the manifests", func(ctx SpecContext) {
		By("prepare test with all preconditions")
		configMapName := testutils.GenerateName("test-")
		configMap := corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name: configMapName,
			},
		}
		configMapRaw, err := ToRaw(&configMap)
		Expect(err).ToNot(HaveOccurred())

		description := v1alpha1.ChallengeDescription{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-",
				Namespace:    corev1.NamespaceDefault,
			},
			Spec: v1alpha1.ChallengeDescriptionSpec{
				Title:       "test",
				Description: "test",
				Flag:        "test",
				Manifests: []runtime.RawExtension{
					{
						Raw: configMapRaw,
					},
				},
			},
		}
		Expect(k8sClient.Create(ctx, &description)).To(Succeed())

		instance := v1alpha1.ChallengeInstance{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-",
				Namespace:    corev1.NamespaceDefault,
			},
			Spec: v1alpha1.ChallengeInstanceSpec{
				ChallengeDescriptionName: description.Name,
			},
		}
		Expect(k8sClient.Create(ctx, &instance)).To(Succeed())
		Expect(instance.Spec.ChallengeDescriptionName).ToNot(BeZero())

		namespace := corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: instance.Name,
			},
		}
		Expect(k8sClient.Create(ctx, &namespace)).To(Succeed())

		By("run the reconciler")
		result, err := reconciler.Reconcile(ctx, testutils.RequestFromObject(&instance))
		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(BeZero())

		By("verify all postconditions")
		Expect(k8sClient.Get(ctx, client.ObjectKey{
			Name:      configMap.Name,
			Namespace: instance.Name,
		}, &configMap)).To(Succeed())
	})

	It("should succeed if the manifests are already there", func(ctx SpecContext) {
		By("prepare test with all preconditions")
		configMapName := testutils.GenerateName("test-")
		configMap := corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name: configMapName,
			},
		}
		configMapRaw, err := ToRaw(&configMap)
		Expect(err).ToNot(HaveOccurred())

		description := v1alpha1.ChallengeDescription{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-",
				Namespace:    corev1.NamespaceDefault,
			},
			Spec: v1alpha1.ChallengeDescriptionSpec{
				Title:       "test",
				Description: "test",
				Flag:        "test",
				Manifests: []runtime.RawExtension{
					{
						Raw: configMapRaw,
					},
				},
			},
		}
		Expect(k8sClient.Create(ctx, &description)).To(Succeed())

		instance := v1alpha1.ChallengeInstance{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-",
				Namespace:    corev1.NamespaceDefault,
			},
			Spec: v1alpha1.ChallengeInstanceSpec{
				ChallengeDescriptionName: description.Name,
			},
		}
		Expect(k8sClient.Create(ctx, &instance)).To(Succeed())
		Expect(instance.Spec.ChallengeDescriptionName).ToNot(BeZero())

		namespace := corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: instance.Name,
			},
		}
		Expect(k8sClient.Create(ctx, &namespace)).To(Succeed())
		configMap.Namespace = instance.Name
		Expect(k8sClient.Create(ctx, &configMap)).To(Succeed())

		By("run the reconciler")
		result, err := reconciler.Reconcile(ctx, testutils.RequestFromObject(&instance))
		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(BeZero())

		By("verify all postconditions")
		Expect(k8sClient.Get(ctx, client.ObjectKey{
			Name:      configMap.Name,
			Namespace: instance.Name,
		}, &configMap)).To(Succeed())
	})

	It("should fail if the referenced challenge description is missing", func(ctx SpecContext) {
		By("prepare test with all preconditions")
		instance := v1alpha1.ChallengeInstance{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-",
				Namespace:    corev1.NamespaceDefault,
			},
			Spec: v1alpha1.ChallengeInstanceSpec{
				ChallengeDescriptionName: testutils.GenerateName("not-existing-"),
			},
		}
		Expect(k8sClient.Create(ctx, &instance)).To(Succeed())
		Expect(instance.Spec.ChallengeDescriptionName).ToNot(BeZero())

		namespace := corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: instance.Name,
			},
		}
		Expect(k8sClient.Create(ctx, &namespace)).To(Succeed())

		By("run the reconciler")
		result, err := reconciler.Reconcile(ctx, testutils.RequestFromObject(&instance))
		Expect(err).To(HaveOccurred())
		Expect(result).To(BeZero())
	})

	It("should fail with a malformed manifests", func(ctx SpecContext) {
		By("prepare test with all preconditions")
		description := v1alpha1.ChallengeDescription{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-",
				Namespace:    corev1.NamespaceDefault,
			},
			Spec: v1alpha1.ChallengeDescriptionSpec{
				Title:       "test",
				Description: "test",
				Flag:        "test",
				Manifests: []runtime.RawExtension{
					{
						Raw: []byte(`{"kind":"NotExisting"}`),
					},
				},
			},
		}
		Expect(k8sClient.Create(ctx, &description)).To(Succeed())

		instance := v1alpha1.ChallengeInstance{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-",
				Namespace:    corev1.NamespaceDefault,
			},
			Spec: v1alpha1.ChallengeInstanceSpec{
				ChallengeDescriptionName: description.Name,
			},
		}
		Expect(k8sClient.Create(ctx, &instance)).To(Succeed())
		Expect(instance.Spec.ChallengeDescriptionName).ToNot(BeZero())

		namespace := corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: instance.Name,
			},
		}
		Expect(k8sClient.Create(ctx, &namespace)).To(Succeed())

		By("run the reconciler")
		result, err := reconciler.Reconcile(ctx, testutils.RequestFromObject(&instance))
		Expect(err).To(HaveOccurred())
		Expect(result).To(BeZero())
	})

	It("should not create the manifests when the instance is deleted", func(ctx SpecContext) {
		By("prepare test with all preconditions")
		configMapName := testutils.GenerateName("test-")
		configMap := corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name: configMapName,
			},
		}
		configMapRaw, err := ToRaw(&configMap)
		Expect(err).ToNot(HaveOccurred())

		description := v1alpha1.ChallengeDescription{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-",
				Namespace:    corev1.NamespaceDefault,
			},
			Spec: v1alpha1.ChallengeDescriptionSpec{
				Title:       "test",
				Description: "test",
				Flag:        "test",
				Manifests: []runtime.RawExtension{
					{
						Raw: configMapRaw,
					},
				},
			},
		}
		Expect(k8sClient.Create(ctx, &description)).To(Succeed())

		instance := v1alpha1.ChallengeInstance{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-",
				Namespace:    corev1.NamespaceDefault,
				Finalizers: []string{
					challengeinstance.FinalizerName,
				},
			},
			Spec: v1alpha1.ChallengeInstanceSpec{
				ChallengeDescriptionName: description.Name,
			},
		}
		Expect(k8sClient.Create(ctx, &instance)).To(Succeed())
		Expect(k8sClient.Delete(ctx, &instance)).To(Succeed())
		Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(&instance), &instance)).To(Succeed())
		Expect(instance.DeletionTimestamp.IsZero()).To(BeFalse())
		Expect(instance.Spec.ChallengeDescriptionName).ToNot(BeZero())

		namespace := corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: instance.Name,
			},
		}
		Expect(k8sClient.Create(ctx, &namespace)).To(Succeed())

		By("run the reconciler")
		result, err := reconciler.Reconcile(ctx, testutils.RequestFromObject(&instance))
		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(BeZero())

		By("verify all postconditions")
		Expect(k8sClient.Get(ctx, client.ObjectKey{
			Name:      configMap.Name,
			Namespace: instance.Name,
		}, &configMap)).To(MatchError(ContainSubstring("not found")))
	})
})
