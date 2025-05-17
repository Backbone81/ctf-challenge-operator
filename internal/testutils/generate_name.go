package testutils

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/onsi/gomega"
)

// GenerateName simulates the behavior of Kubernetes GenerateName for test situations where we need to know the name
// beforehand.
func GenerateName(prefix string) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	suffix := make([]byte, 5)
	for i := range suffix {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

		suffix[i] = charset[n.Int64()]
	}
	return fmt.Sprintf("%s%s", prefix, string(suffix))
}
