package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ChallengeInstanceSpec defines the desired state of ChallengeInstance.
type ChallengeInstanceSpec struct {
	// ExpirationSeconds is the requested duration of validity of the Challenge instance.
	// +optional
	ExpirationSeconds *int64 `json:"expirationSeconds"`
}

// ChallengeInstanceStatus defines the observed state of ChallengeInstance.
type ChallengeInstanceStatus struct {
	// ExpirationTimestamp is the time of expiration of the challenge instance.
	ExpirationTimestamp metav1.Time `json:"expirationTimestamp"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Expiration",type="string",format="date-time",JSONPath=".status.expirationTimestamp"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

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
