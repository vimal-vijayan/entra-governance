package client

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

func getSecretData(secret *corev1.Secret, key string) (string, error) {

	if secret == nil {
		return "", fmt.Errorf("secret is nil")
	}

	value, exists := secret.Data[key]
	if !exists {
		return "", fmt.Errorf("key %s not found in secret %s/%s", key, secret.Namespace, secret.Name)
	}

	return string(value), nil
}
