package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ChallengeDescriptionSpec defines the desired state of ChallengeDescription.
type ChallengeDescriptionSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of ChallengeDescription. Edit challengedescription_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// ChallengeDescriptionStatus defines the observed state of ChallengeDescription.
type ChallengeDescriptionStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

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
