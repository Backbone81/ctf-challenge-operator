package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// APIKeySpec defines the desired state of APIKey.
type APIKeySpec struct {
	// ExpirationSeconds is the requested duration of validity of the API key.
	// +optional
	ExpirationSeconds *int64 `json:"expirationSeconds"`
}

// APIKeyStatus defines the observed state of APIKey.
type APIKeyStatus struct {
	// Key is the opaque API key.
	// +optional
	Key string `json:"token"`

	// ExpirationTimestamp is the time of expiration of the returned API key.
	// +optional
	ExpirationTimestamp metav1.Time `json:"expirationTimestamp"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Expiration",type="string",format="date-time",JSONPath=".status.expirationTimestamp"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// APIKey is the Schema for the apikeys API.
type APIKey struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   APIKeySpec   `json:"spec,omitempty"`
	Status APIKeyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// APIKeyList contains a list of APIKey.
type APIKeyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []APIKey `json:"items"`
}

func init() {
	SchemeBuilder.Register(&APIKey{}, &APIKeyList{})
}
