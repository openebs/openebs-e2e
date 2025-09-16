package etcd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/e2e_config"
	"github.com/openebs/openebs-e2e/common/k8stest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

type EtcdCtl struct {
	Endpoint string `json:"Endpoint"`
	Status   struct {
		Header struct {
			Revision int64 `json:"revision"`
		} `json:"header"`
		StorageVersion string `json:"storageVersion"`
	} `json:"Status"`
}

var etcdStsName = e2e_config.GetConfig().Product.ControlPlaneEtcd

// Get Etcd StatefulSet Pod Names
func GetEtcdPodNames() ([]string, error) {
	pods, err := k8stest.GetStsPodNames(etcdStsName, common.NSMayastor())
	if err != nil {
		return nil, fmt.Errorf("failed to list etcd pods by sts name: %v", err)
	}
	if len(pods) == 0 {
		return nil, fmt.Errorf("no etcd pods found for sts name: %s", etcdStsName)
	}
	return pods, nil
}

// Execute a command in all Etcd pods
func ExecuteInEtcdPods(command string) (map[string]string, error) {
	pods, err := GetEtcdPodNames()
	if err != nil {
		return nil, fmt.Errorf("failed to get etcd pods: %v", err)
	}
	results := make(map[string]string)
	for _, pod := range pods {
		output, _, err := k8stest.ExecuteCommandInPod(common.NSMayastor(), pod, command)
		if err != nil {
			return nil, fmt.Errorf("failed to execute command in pod %s: %v", pod, err)
		}
		results[pod] = output
	}
	return results, nil
}

// Get health of Etcd endpoints using etcdctl
func GetHealthEtcd() (map[string]string, error) {
	command := "etcdctl endpoint health --cluster"
	return ExecuteInEtcdPods(command)
}

// Verify if all Etcd pods are healthy
// Return true if all are healthy else false
func AreAllEtcdPodsHealthy() (bool, error) {
	healths, err := GetHealthEtcd()
	if err != nil {
		return false, fmt.Errorf("failed to get health of etcd pods: %v", err)
	}
	// Display raw output for debugging
	for podName, health := range healths {
		logf.Log.Info("Raw output from etcd pod", "pod", podName, "health", health)
	}
	// Check if all etcd pods are healthy
	for podName, health := range healths {
		if strings.Contains(health, "is unhealthy") {
			return false, fmt.Errorf("etcd pod %s is not healthy: %s", podName, health)
		}
	}
	return true, nil
}

// Get Status of all Etcd pods
func GetStatusInEtcdPods() (map[string]string, error) {
	command := "etcdctl endpoint status --cluster --write-out=json"
	results, err := ExecuteInEtcdPods(command)
	if err != nil {
		return nil, fmt.Errorf("failed to execute command in etcd pods: %v", err)
	}
	// Display raw output for debugging
	for podName, output := range results {
		logf.Log.Info("Raw output from etcd pod", "pod", podName, "output", output)
	}
	return results, nil
}

// Get Revision on all Etcd pods
func GetRevisionInEtcdPods() (map[string]int64, error) {
	podOutputs, err := GetStatusInEtcdPods()
	if err != nil {
		return nil, fmt.Errorf("failed to get status of etcd pods: %v", err)
	}
	revisions := make(map[string]int64)
	for pod, output := range podOutputs {
		var eps []EtcdCtl
		if err := json.Unmarshal([]byte(output), &eps); err != nil {
			return nil, fmt.Errorf("failed to unmarshal status output from pod %s: %v", pod, err)
		}
		for _, ep := range eps { // Assuming each pod's output contains status for its own endpoint
			revisions[pod] = ep.Status.Header.Revision
		}
		logf.Log.Info("Etcd pod revision", "pod", pod, "revision", revisions[pod])
	}
	return revisions, nil
}

// Verify if all Etcd pods have the same revision
// Return true if all have the same revision else false
func AreAllEtcdPodsSameRevision() (bool, error) {
	revisions, err := GetRevisionInEtcdPods()
	if err != nil {
		return false, fmt.Errorf("failed to get revision of etcd pods: %v", err)
	}
	var lastRevision int64
	for podName, revision := range revisions {
		if lastRevision == 0 {
			lastRevision = revision
		} else if lastRevision != revision && revision != 0 {
			return false, fmt.Errorf("etcd pod %s has different revision: %d", podName, revision)
		}
	}
	return true, nil
}

// Get Storage Version of all Etcd pods from all Etcd pods using etcdctl
func GetStorageVersionInEtcdPods() (map[string]string, error) {
	results, err := GetStatusInEtcdPods()
	if err != nil {
		return nil, fmt.Errorf("failed to execute command in etcd pods: %v", err)
	}
	for pod, output := range results {
		var eps []EtcdCtl
		if err := json.Unmarshal([]byte(output), &eps); err != nil {
			return nil, fmt.Errorf("failed to unmarshal status output from pod %s: %v", pod, err)
		}
		for _, ep := range eps {
			results[pod] = ep.Status.StorageVersion
		}
		logf.Log.Info("Etcd pod storage version", "pod", pod, "version", results[pod])
	}
	return results, nil
}

// Verify if all Etcd pods have the same storage version
// Return true if all have the same storage version else false
func AreAllEtcdPodsSameStorageVersion() (bool, error) {
	versions, err := GetStorageVersionInEtcdPods()
	if err != nil {
		return false, fmt.Errorf("failed to get storage version of etcd pods: %v", err)
	}
	var lastVersion string
	for podName, version := range versions {
		if lastVersion == "" {
			lastVersion = version
		} else if lastVersion != version && version != "" {
			return false, fmt.Errorf("etcd pod %s has different storage version: %s", podName, version)
		}
	}
	return true, nil
}
