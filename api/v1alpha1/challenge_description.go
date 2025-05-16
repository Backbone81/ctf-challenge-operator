package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// ChallengeDescriptionSpec defines the desired state of ChallengeDescription.
type ChallengeDescriptionSpec struct {
	// Title is the name of the challenge
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Title string `json:"title"`

	// Description is the content of the challenge
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Description string `json:"description"`

	// Category is the category this challenge belongs to.
	// +kubebuilder:validation:Optional
	Category string `json:"category"`

	// Value is the number of points which are added upon solving the challenge.
	// +kubebuilder:default=0
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=0
	Value int `json:"value"`

	// Hints provides a list of hints to help solve the challenge.
	// +kubebuilder:validation:Optional
	Hints []ChallengeHint `json:"hints"`

	// Manifests provide the Kubernetes manifests which should be created when a new instance of the challenge is
	// requested. The manifests are placed in a dedicated namespace. The namespace provided in those manifests is
	// overwritten.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	Manifests []runtime.RawExtension `json:"manifests"`
}

type ChallengeHint struct {
	// Description is the content of the hint.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Description string `json:"description"`

	// Cost is the number of points which are to be deducted from the overall score if this hint is being used.
	// +kubebuilder:default=0
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=0
	Cost int `json:"cost"`
}

// ChallengeDescriptionStatus defines the observed state of ChallengeDescription.
type ChallengeDescriptionStatus struct{}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Title",type="string",JSONPath=".spec.title"
// +kubebuilder:printcolumn:name="Category",type="string",JSONPath=".spec.category"
// +kubebuilder:printcolumn:name="Value",type="integer",JSONPath=".spec.value"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// ChallengeDescription is the Schema for the challengedescriptions API.
type ChallengeDescription struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ChallengeDescriptionSpec   `json:"spec,omitempty"`
	Status ChallengeDescriptionStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ChallengeDescriptionList contains a list of ChallengeDescription.
type ChallengeDescriptionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ChallengeDescription `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ChallengeDescription{}, &ChallengeDescriptionList{})
}
