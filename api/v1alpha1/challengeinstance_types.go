package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ChallengeInstanceSpec defines the desired state of ChallengeInstance.
type ChallengeInstanceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of ChallengeInstance. Edit challengeinstance_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// ChallengeInstanceStatus defines the observed state of ChallengeInstance.
type ChallengeInstanceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ChallengeInstance is the Schema for the challengeinstances API.
type ChallengeInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ChallengeInstanceSpec   `json:"spec,omitempty"`
	Status ChallengeInstanceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ChallengeInstanceList contains a list of ChallengeInstance.
type ChallengeInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ChallengeInstance `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ChallengeInstance{}, &ChallengeInstanceList{})
}
