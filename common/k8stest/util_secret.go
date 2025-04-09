package k8stest

import (
	"context"
	"fmt"
	"sort"
	"strings"

	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

type EncryptionKeyValue struct {
	Value  string
	Length int
}

func CreateAesXtsSecret(secretName, namespace string, keyMap map[string]EncryptionKeyValue) (*corev1.Secret, error) {
	// get Secret parameter as a string
	encryptionParameters := GetEncryptionParameterJsonString("AesXts", secretName, keyMap)
	logf.Log.Info("Creating AesXts secret", "name", secretName, "namespace", namespace)
	// Create the secret
	return CreateSecret(secretName, namespace, encryptionParameters)
}

func GetEncryptionParameterJsonString(cipherName, keyName string, keyMap map[string]EncryptionKeyValue) string {

	// Start building the JSON string
	var b strings.Builder
	b.WriteString("{\n")
	b.WriteString(fmt.Sprintf("  \"cipher\": \"%s\",\n", cipherName))
	b.WriteString(fmt.Sprintf("  \"key_name\": \"%s\",\n", keyName))

	// Get sorted keys for consistent output
	var sortedKeys []string
	for k := range keyMap {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	// Add key/value pairs
	for i, k := range sortedKeys {
		v := keyMap[k]
		b.WriteString(fmt.Sprintf("  \"%s\": \"%s\",\n", k, v.Value))
		b.WriteString(fmt.Sprintf("  \"%s_len\": %d", k, v.Length))
		if i < len(sortedKeys)-1 {
			b.WriteString(",\n")
		} else {
			b.WriteString("\n")
		}
	}

	b.WriteString("}")

	// Final encryption_parameters JSON string
	encryptionParameters := b.String()
	return encryptionParameters
}

func boolPtr(b bool) *bool {
	return &b
}

// CreateSecret creates a Kubernetes Secret with the given name, namespace, and encryption parameters.
// It returns the created Secret object or an error if the creation fails.
// The Secret is marked as immutable and of type Opaque.
// The encryption parameters are stored in the StringData field of the Secret.
func CreateSecret(secretName, namespace, encryptionParameters string) (*corev1.Secret, error) {
	// Create the Secret object
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
		},
		Immutable: boolPtr(true),
		Type:      corev1.SecretTypeOpaque,
		StringData: map[string]string{
			"encryption_parameters": encryptionParameters,
		},
	}
	secretApi := gTestEnv.KubeInt.CoreV1().Secrets

	// Create the Secret in the cluster
	secret, err := secretApi(namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		return secret, fmt.Errorf("failed to create secret: %v", err)
	}

	logf.Log.Info("Secret created successfully", "name", secretName, "namespace", namespace)
	return secret, nil
}

// GetSecret retrieves a Kubernetes Secret with the given name and namespace.
// It returns the Secret object or an error if the retrieval fails.
func GetSecret(secretName, namespace string) (*corev1.Secret, error) {
	secretApi := gTestEnv.KubeInt.CoreV1().Secrets

	// Get the Secret in the cluster
	secret, err := secretApi(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get secret %s in namespace %s, error: %v", secretName, namespace, err)
	}
	logf.Log.Info("Secret retrieved successfully", "name", secretName, "namespace", namespace)
	// Check if the secret is nil
	if secret == nil {
		return secret, fmt.Errorf("secret %s in namespace %s is nil", secretName, namespace)
	}
	// Check if the secret is empty
	if len(secret.Data) == 0 {
		return secret, fmt.Errorf("secret %s in namespace %s is empty", secretName, namespace)
	}

	return secret, nil
}

func IsSecretPresent(secretName, namespace string) (bool, error) {
	secretApi := gTestEnv.KubeInt.CoreV1().Secrets

	// Get the Secret in the cluster
	_, err := secretApi(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if k8serrors.IsNotFound(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to get secret %s in namespace %s, error: %v", secretName, namespace, err)
	}
	logf.Log.Info("Secret retrieved successfully", "name", secretName, "namespace", namespace)
	return true, nil
}

// DeleteSecret deletes a Kubernetes Secret with the given name and namespace.
// It returns the deleted Secret object or an error if the deletion fails.
func DeleteSecret(secretName, namespace string) error {
	secretApi := gTestEnv.KubeInt.CoreV1().Secrets

	// Delete the Secret in the cluster
	err := secretApi(namespace).Delete(context.TODO(), secretName, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete secret %s in namespace %s, error: %v", secretName, namespace, err)
	}
	logf.Log.Info("Secret deleted successfully", "name", secretName, "namespace", namespace)
	return nil
}

// ListSecrets lists all Kubernetes Secrets in the given namespace.
// It returns a slice of Secret objects or an error if the listing fails.
func ListSecrets(namespace string) ([]corev1.Secret, error) {
	secretApi := gTestEnv.KubeInt.CoreV1().Secrets

	// List all Secrets in the namespace
	secrets, err := secretApi(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list secrets in namespace %s, error: %v", namespace, err)
	}

	logf.Log.Info("Secrets listed successfully", "namespace", namespace)
	return secrets.Items, nil
}

// GetSecretData retrieves the data from a Kubernetes Secret with the given name and namespace.
// It returns the data as a map of string keys to byte slices or an error if the retrieval fails.
func GetSecretData(secretName, namespace string) (map[string][]byte, error) {
	secretApi := gTestEnv.KubeInt.CoreV1().Secrets

	// Get the Secret in the cluster
	secret, err := secretApi(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get secret %s in namespace %s, error: %v", secretName, namespace, err)
	}

	logf.Log.Info("Secret data retrieved successfully", "name", secretName, "namespace", namespace)
	return secret.Data, nil
}

// GetSecretStringData retrieves the StringData from a Kubernetes Secret with the given name and namespace.
// It returns the StringData as a map of string keys to string values or an error if the retrieval fails.
func GetSecretStringData(secretName, namespace string) (map[string]string, error) {
	secretApi := gTestEnv.KubeInt.CoreV1().Secrets

	// Get the Secret in the cluster
	secret, err := secretApi(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get secret %s in namespace %s, error: %v", secretName, namespace, err)
	}

	logf.Log.Info("Secret StringData retrieved successfully", "name", secretName, "namespace", namespace)
	return secret.StringData, nil
}
