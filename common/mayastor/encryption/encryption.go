package encryption

import (
	"fmt"
	"os"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/custom_resources"
	"github.com/openebs/openebs-e2e/common/custom_resources/types"
	"github.com/openebs/openebs-e2e/common/k8stest"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var poolOnlineTimeoutSec = 120 // seconds

// CreateEncryptedPools creates encrypted pools
// It takes a map of node, disk and aesXtsSecret as input and returns a map of pool name and pool object
func CreateEncryptedPools(nodeDiskMap map[string]string, aesXtsSecret string) (map[string]types.DiskPool, error) {
	pools := make(map[string]types.DiskPool)
	for node, disk := range nodeDiskMap {
		poolName := "pool-" + node
		pool, err := custom_resources.CreateMsPoolWithEncryption(poolName, node, []string{disk}, aesXtsSecret)
		if err != nil {
			return pools, err
		}
		pools[poolName] = pool
	}
	return pools, nil
}

// CreateAndWaitForEncryptedPools creates encrypted pools and waits for them to be online
// It takes a map of node, disk and aesXtsSecret as input and returns a map of pool name and pool object
// It returns an error if the pools are not online within the given time
// It also refreshes the pool list after creating the pools
func CreateAndWaitForEncryptedPools(nodeDiskMap map[string]string, aesXtsSecret string) (map[string]types.DiskPool, error) {
	pools, err := CreateEncryptedPools(nodeDiskMap, aesXtsSecret)
	if err != nil {
		return pools, err
	}

	// wait for all pools to be online
	err = k8stest.WaitForPoolsToBeOnline(poolOnlineTimeoutSec)
	if err != nil {
		return pools, fmt.Errorf("all pools are not online, error: %v", err)
	}

	// refresh the pool list
	poolList, err := custom_resources.ListMsPools()
	if err != nil {
		return pools, fmt.Errorf("failed to list pools, error: %v", err)
	}

	for _, pool := range poolList {
		if pools[pool.GetName()] != nil {
			pools[pool.GetName()] = pool

		}
	}

	// check if all pools are encrypted
	for _, pool := range pools {
		if !pool.IsPoolEncrypted() {
			return pools, fmt.Errorf("pool %s is not encrypted", pool.GetName())
		}
	}
	return pools, nil
}

// CreateAesXtsEncryptionSecret creates a k8s secret with the given name and namespace
// The secret is created with the encryption keys from the environment variables
func CreateAesXtsEncryptionSecret(secretName string, namespace string) error {
	log.Log.Info("Creating k8s secret", "secret name", secretName)
	keys := GetEncryptionKeys()
	if len(keys) == 0 {
		return fmt.Errorf("no encryption key found")
	} else if len(keys) == 1 {
		return fmt.Errorf("not enough encryption key found, required two keys but got single key")
	}
	keyMap := map[string]k8stest.EncryptionKeyValue{
		"key": {
			Value:  keys[0],
			Length: 128,
		},
		"key2": {
			Value:  keys[1],
			Length: 128,
		},
	}

	_, err := k8stest.CreateAesXtsSecret(secretName, namespace, keyMap)
	if err != nil {
		return fmt.Errorf("failed to create aes xts secret, error: %v", err)
	}

	// Get secret
	secret, err := k8stest.GetSecret(secretName, common.NSMayastor())
	if err != nil {
		return fmt.Errorf("failed to get secret %s in %s namespace, error: %v", secretName, namespace, err)
	}
	if secret == nil {
		return fmt.Errorf("secret %s in %s namespace not found", secretName, namespace)
	}
	log.Log.Info("Secret Created", "secret name", secretName)
	return nil
}

// GetEncryptionKeys returns the encryption keys from the environment variables
// The keys are expected to be set in the environment variables encryption_key1 and encryption_key2
// The keys are expected to be 128 bits long
func GetEncryptionKeys() []string {
	keys := make([]string, 0)
	// Get first key
	key1, err := GetKeyFromEnv("encryption_key1")
	if err != nil {
		log.Log.Error(err, "Failed to get first key from env", "key name", key1)
		return keys
	}
	keys = append(keys, key1)
	// Get second key
	key2, err := GetKeyFromEnv("encryption_key2")
	if err != nil {
		log.Log.Error(err, "Failed to get second key from env", "key name", key2)
		return keys
	}
	keys = append(keys, key2)
	return keys
}

func GetKeyFromEnv(keyName string) (string, error) {
	key := os.Getenv(keyName)
	if key == "" {
		return "", fmt.Errorf("key %s not found in environment", keyName)
	}
	return key, nil
}

// GetNodeAndDiskOfCreatedPool returns the node and disk of the created pool
func GetNodeAndDiskOfCreatedPool(poolName string) (string, []string, error) {
	pool, err := custom_resources.GetMsPool(poolName)
	if err != nil {
		return "", []string{}, fmt.Errorf("failed to get pool %s, error: %v", poolName, err)
	}
	if pool == nil {
		return "", []string{}, fmt.Errorf("pool %s not found", poolName)
	}
	node := pool.GetSpecNode()
	disk := pool.GetSpecDisks()
	return node, disk, nil
}

// GetNodeDiskMapOfPools returns a map of node and disk for all pools
func GetNodeDiskMapOfPools() (map[string]string, error) {
	nodeDiskMap := make(map[string]string)
	poolList, err := custom_resources.ListMsPools()
	if err != nil {
		return nodeDiskMap, fmt.Errorf("failed to list pools, error: %v", err)
	}
	for _, pool := range poolList {
		node, disk, err := GetNodeAndDiskOfCreatedPool(pool.GetName())
		if err != nil {
			return nodeDiskMap, fmt.Errorf("failed to get node and disk of pool %s, error: %v", pool.GetName(), err)
		}
		if len(disk) == 0 {
			return nodeDiskMap, fmt.Errorf("no disk found for pool %s", pool.GetName())
		}
		nodeDiskMap[node] = disk[0]
	}
	return nodeDiskMap, nil
}
