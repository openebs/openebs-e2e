package pool_expansion

import (
	"fmt"
	"time"

	v1cp "github.com/openebs/openebs-e2e/common/controlplane/v1"
	"github.com/openebs/openebs-e2e/common/custom_resources"
	"github.com/openebs/openebs-e2e/common/k8stest"
	"github.com/openebs/openebs-e2e/common/mayastor/encryption"
	"github.com/openebs/openebs-e2e/common/platform"
	plattypes "github.com/openebs/openebs-e2e/common/platform/types"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// resizeVolumeByDevicePath extracts the Hetzner volume id from a device path and resizes it.
func resizeVolumeByDevicePath(plat plattypes.Platform, devicePath string, targetGiB int) error {
	volId, err := plat.ExtractVolumeIdFromDevicePath(devicePath)
	if err != nil {
		return err
	}
	return plat.ResizeVolume(volId, targetGiB)
}

// CreatePoolWithMaxAndCluster creates a DiskPool with both MaxExpansion and ClusterSize in one operation.
func CreatePoolWithMaxAndCluster(poolName, node, disk, maxExpansion, clusterSize string) error {
	_, err := custom_resources.CreateMsPoolWithMaxSizeAndClusterSize(poolName, node, []string{disk}, maxExpansion, clusterSize)
	return err
}

// CreateEncryptedPoolWithMaxAndCluster creates an encrypted DiskPool with both MaxExpansion and ClusterSize in one operation.
func CreateEncryptedPoolWithMaxAndCluster(poolName, node, disk, encryptionSecretName, maxExpansion, clusterSize string) error {
	_, err := custom_resources.CreateMsPoolWithEncryptionAndMaxSize(poolName, node, []string{disk}, encryptionSecretName, maxExpansion, clusterSize)
	return err
}

func VerifyMaxExpandableInOriginalUnit(poolName string, maxExpansionStr string, expectedCapacity uint64) error {
	// Read DiskPool CR (version-agnostic via custom_resources)
	poolCR, err := custom_resources.GetMsPool(poolName)
	if err != nil {
		return err
	}
	// Status value is a human-readable string (e.g., "255.8 GiB"). May be empty on older CRDs.
	statusMaxExpandable := poolCR.GetStatusMaxExpandableSize()
	if statusMaxExpandable == "" {
		return fmt.Errorf("status.maxExpandableSize not available on this DiskPool version")
	}

	expectedCapacityGiB, err := k8stest.ParseGiBOrBytesToGiB(maxExpansionStr)
	if err != nil {
		return err
	}
	actualCapacityGiB, err := k8stest.ParseGiBOrBytesToGiB(statusMaxExpandable)
	if err != nil {
		return err
	}

	// Check that actual maxExpandable is at least the parsed MaxExpansion value
	// (it may be larger due to rounding or minimum allocation units)
	if actualCapacityGiB < expectedCapacityGiB {
		return fmt.Errorf("max expandable too small: expected at least %.1fGiB got %.1fGiB", expectedCapacityGiB, actualCapacityGiB)
	}

	// Also verify actual maxExpandable in bytes is >= provided expectedCapacity (bytes)
	actualCapacityBytes := uint64(actualCapacityGiB * 1024 * 1024 * 1024)
	if actualCapacityBytes < expectedCapacity {
		return fmt.Errorf("max expandable too small: expected >= %d bytes (current capacity) got %d bytes", expectedCapacity, actualCapacityBytes)
	}
	return nil
}

// ResizeHetznerVolumes resizes the primary volume and optionally other node volumes to the target GiB.
func ResizeHetznerVolumes(nodes []string, targetGiB int) error {
	var plat plattypes.Platform = platform.Create()
	for _, node := range nodes {
		nodeDisks, err := k8stest.GetConfiguredNodePoolDevices(node)
		if err != nil {
			return fmt.Errorf("failed to get pool devices for node %s: %v", node, err)
		}
		if len(nodeDisks) == 0 {
			return fmt.Errorf("no configured devices for node %s", node)
		}
		if err := resizeVolumeByDevicePath(plat, nodeDisks[0], targetGiB); err != nil {
			return fmt.Errorf("failed to resize volume on node %s to %dGiB: %v", node, targetGiB, err)
		}
	}
	return nil
}

// ResizeAllVolumesBeyondPoolMax resizes each pool's backing Hetzner volume to MaxExpandableGiB + plusGiB.
// This uses control-plane pool state and encryption helpers to map pool -> node/disk -> volume id.
func ResizeAllVolumesBeyondPoolMax(pools []string, plusGiB int) error {
	var plat plattypes.Platform = platform.Create()
	for _, poolName := range pools {
		// Read control-plane pool to compute target size
		cpPool, err := v1cp.GetMayastorCpPool(poolName)
		if err != nil {
			return err
		}
		maxExpGiB := int(cpPool.State.MaxExpandableBytes / (1024 * 1024 * 1024))
		targetGiB := maxExpGiB + plusGiB

		// Map pool -> node -> disk
		node, disks, err := encryption.GetNodeAndDiskOfCreatedPool(poolName)
		if err != nil {
			return fmt.Errorf("failed to get node/disk for pool %s: %v", poolName, err)
		}
		if len(disks) == 0 {
			return fmt.Errorf("no disk found for pool %s", poolName)
		}
		if err := resizeVolumeByDevicePath(plat, disks[0], targetGiB); err != nil {
			return fmt.Errorf("failed to resize volume on node %s to %dGiB: %v", node, targetGiB, err)
		}
	}
	return nil
}

// VerifyExpandAnnotationCleared ensures the expand annotation is removed from the DiskPool CR.
func VerifyExpandAnnotationCleared(poolName string) error {
	rawDsp, err := custom_resources.GetMsPool(poolName)
	if err != nil {
		return err
	}
	annMap := rawDsp.GetAnnotations()
	if annMap["openebs.io/expand"] != "" {
		return fmt.Errorf("expand annotation not cleared")
	}
	return nil
}

// VerifyExpandAnnotattion is cleared for all provided poolionClearedForPools checks that the expand annotas.
func VerifyExpandAnnotationClearedForPools(pools []string) error {
	for _, poolName := range pools {
		if err := VerifyExpandAnnotationCleared(poolName); err != nil {
			return fmt.Errorf("pool %s: %w", poolName, err)
		}
	}
	return nil
}

// VerifyExpandAnnotationClearedWithRetry polls until the expand annotation is cleared or timeout elapses.
// It returns nil as soon as the annotation is cleared; otherwise an error after the timeout.
func VerifyExpandAnnotationClearedWithRetry(poolName string, timeout time.Duration, interval time.Duration) error {
	deadline := time.Now().Add(timeout)
	var lastErr error
	for {
		err := VerifyExpandAnnotationCleared(poolName)
		if err == nil {
			return nil
		}
		lastErr = err
		if time.Now().After(deadline) {
			return lastErr
		}
		time.Sleep(interval)
	}
}

// VerifyExpandAnnotationClearedForPoolsWithRetry polls all pools until their expand annotations are cleared or timeout elapses.
func VerifyExpandAnnotationClearedForPoolsWithRetry(pools []string, timeout time.Duration, interval time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		var lastErr error
		allCleared := true
		for _, poolName := range pools {
			if err := VerifyExpandAnnotationCleared(poolName); err != nil {
				allCleared = false
				lastErr = err
				break
			}
		}
		if allCleared {
			return nil
		}
		if time.Now().After(deadline) {
			if lastErr != nil {
				return lastErr
			}
			return fmt.Errorf("expand annotation not cleared on all pools before timeout")
		}
		time.Sleep(interval)
	}
}

// AnnotatePoolsAndWaitForExpansion annotates each pool for expansion and waits until capacity increases.
// Returns the expanded capacity of the primary pool once its expansion is observed.
func AnnotatePoolsAndWaitForExpansion(pools []string, primaryPool string, waitTimeout time.Duration, pollInterval time.Duration) (uint64, error) {
	var expandedPrimaryCap uint64
	for _, poolName := range pools {
		// Read initial capacity
		msp, err := k8stest.GetMsPool(poolName)
		if err != nil {
			return 0, fmt.Errorf("failed to get pool %s: %w", poolName, err)
		}
		initial := msp.Status.Capacity

		// Annotate for expansion
		if err := custom_resources.AnnotatePoolForExpansion(poolName); err != nil {
			return 0, fmt.Errorf("failed to annotate pool %s for expansion: %w", poolName, err)
		}
		logf.Log.Info("Annotated pool for expansion", "poolName", poolName, "initialCapacity", initial)

		// Wait for capacity to increase
		var capAfter uint64
		deadline := time.Now().Add(waitTimeout)
		for {
			latest, err := k8stest.GetMsPool(poolName)
			if err == nil && latest != nil {
				capAfter = latest.Status.Capacity
				if capAfter > initial {
					break
				}
			}
			if time.Now().After(deadline) {
				return 0, fmt.Errorf("pool %s capacity did not increase after expansion", poolName)
			}
			time.Sleep(pollInterval)
		}
		logf.Log.Info("Pool capacity increased after expansion", "poolName", poolName, "from", initial, "to", capAfter)
		if poolName == primaryPool {
			expandedPrimaryCap = capAfter
		}
	}
	if expandedPrimaryCap == 0 {
		return 0, fmt.Errorf("primary pool %s did not report expanded capacity", primaryPool)
	}
	return expandedPrimaryCap, nil
}

// CaptureDiskCapacities returns a map of poolName -> DiskCapacityBytes from control-plane.
func CaptureDiskCapacities(pools []string) (map[string]uint64, error) {
	res := make(map[string]uint64, len(pools))
	for _, poolName := range pools {
		cpPool, err := v1cp.GetMayastorCpPool(poolName)
		if err != nil {
			return nil, err
		}
		res[poolName] = cpPool.State.DiskCapacityBytes
	}
	return res, nil
}

// VerifyNoFurtherExpansionAndDiskBehavior verifies that pool capacity does not change after a second expansion
// annotation and that disk capacity behavior is correct: it never decreases, and increases only if previous disk
// capacity was below MaxExpandable.
func VerifyNoFurtherExpansionAndDiskBehavior(pools []string, waitTimeout time.Duration, pollInterval time.Duration, prevDiskCaps map[string]uint64) error {
	for _, poolName := range pools {
		// Baseline capacity before second annotation
		ms, err := k8stest.GetMsPool(poolName)
		if err != nil {
			return err
		}
		baseline := ms.Status.Capacity

		// Wait to ensure capacity remains unchanged
		deadline := time.Now().Add(waitTimeout)
		for {
			latest, err := k8stest.GetMsPool(poolName)
			if err == nil && latest != nil {
				// Fail fast if capacity changes during the observation window
				if latest.Status.Capacity != baseline {
					return fmt.Errorf("pool %s capacity changed during second expansion: baseline=%d now=%d", poolName, baseline, latest.Status.Capacity)
				}
			}
			if time.Now().After(deadline) {
				break
			}
			time.Sleep(pollInterval)
		}

		// Control-plane verification: capacity unchanged and disk behavior
		cpNow, err := v1cp.GetMayastorCpPool(poolName)
		if err != nil {
			return err
		}
		if cpNow.State.Capacity != baseline {
			return fmt.Errorf("pool %s capacity changed during second expansion: baseline=%d now=%d", poolName, baseline, cpNow.State.Capacity)
		}
		prevDisk := prevDiskCaps[poolName]
		if cpNow.State.DiskCapacityBytes < prevDisk {
			return fmt.Errorf("pool %s disk capacity decreased during second expansion: prev=%d now=%d", poolName, prevDisk, cpNow.State.DiskCapacityBytes)
		}
		if prevDisk < cpNow.State.MaxExpandableBytes && cpNow.State.DiskCapacityBytes == prevDisk {
			return fmt.Errorf("pool %s disk capacity did not increase though below maxExpandable: prev=%d max=%d now=%d", poolName, prevDisk, cpNow.State.MaxExpandableBytes, cpNow.State.DiskCapacityBytes)
		}
		logf.Log.Info("Verified capacity unchanged and disk behavior after second expansion", "poolName", poolName, "capacity", baseline, "diskCapacity", cpNow.State.DiskCapacityBytes, "maxExpandable", cpNow.State.MaxExpandableBytes)
	}
	return nil
}

// CreatePoolsWithMaxExpansionAndClusterSizeOnAllNodes creates a pool on each provided node using
// the given MaxExpansion and ClusterSize values. It returns the created pool names and the primary pool info.
func CreatePoolsWithMaxExpansionOnAllNodes(allNodes []string, maxExpansionStr string, clusterSizeArg string) (createdPools []string, primaryPoolName, primaryDisk string, err error) {
	createdPools = make([]string, 0)

	for _, n := range allNodes {
		// Get the first available disk on this node
		devs, derr := k8stest.GetConfiguredNodePoolDevices(n)
		if derr != nil || len(devs) == 0 {
			logf.Log.Info("Skipping node without configured pool device", "node", n, "err", derr)
			continue
		}
		disk := devs[0]
		pName := fmt.Sprintf("pool-expansion-test-%s", n)

		logf.Log.Info("Creating pool with MaxExpansion and ClusterSize", "poolName", pName, "clusterSize", clusterSizeArg)
		if err = CreatePoolWithMaxAndCluster(pName, n, disk, maxExpansionStr, clusterSizeArg); err != nil {
			return
		}
		createdPools = append(createdPools, pName)
		if primaryPoolName == "" {
			primaryPoolName = pName
			primaryDisk = disk
		}
	}
	if len(createdPools) == 0 {
		err = fmt.Errorf("no pools could be created on any node")
		return
	}
	return createdPools, primaryPoolName, primaryDisk, nil
}

// CreateEncryptedPoolsWithMaxExpansionOnAllNodes creates an encrypted pool on each provided node using
// the given encryption secret, MaxExpansion and ClusterSize values. It returns the created pool names and the primary pool info.
func CreateEncryptedPoolsWithMaxExpansionOnAllNodes(allNodes []string, encryptionSecretName string, maxExpansionStr string, clusterSizeArg string) (createdPools []string, primaryPoolName, primaryDisk string, err error) {
	createdPools = make([]string, 0)

	for _, n := range allNodes {
		// Get the first available disk on this node
		devs, derr := k8stest.GetConfiguredNodePoolDevices(n)
		if derr != nil || len(devs) == 0 {
			logf.Log.Info("Skipping node without configured pool device", "node", n, "err", derr)
			continue
		}
		disk := devs[0]
		pName := fmt.Sprintf("enc-pool-expansion-test-%s", n)

		logf.Log.Info("Creating encrypted pool with MaxExpansion and ClusterSize", "poolName", pName, "clusterSize", clusterSizeArg, "encryptionSecret", encryptionSecretName)
		if err = CreateEncryptedPoolWithMaxAndCluster(pName, n, disk, encryptionSecretName, maxExpansionStr, clusterSizeArg); err != nil {
			return
		}
		createdPools = append(createdPools, pName)
		if primaryPoolName == "" {
			primaryPoolName = pName
			primaryDisk = disk
		}
	}
	if len(createdPools) == 0 {
		err = fmt.Errorf("no pools could be created on any node")
		return
	}
	return createdPools, primaryPoolName, primaryDisk, nil
}
