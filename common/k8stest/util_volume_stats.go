package k8stest

import (
	"context"
	"encoding/json"
	"fmt"
)

// VolumeStats is one pod volume's usage as the kubelet's own /stats/summary
// API reports it -- the same numbers a CSI driver's NodeGetVolumeStats
// surfaces to the scheduler/metrics pipeline, read directly from the
// kubelet rather than inferred from inside the pod.
type VolumeStats struct {
	AvailableBytes uint64
	CapacityBytes  uint64
	UsedBytes      uint64
	InodesFree     uint64
	Inodes         uint64
	InodesUsed     uint64
}

// summaryResponse is the kubelet /stats/summary response, trimmed to the
// fields VolumeStats needs.
type summaryResponse struct {
	Pods []struct {
		PodRef struct {
			Name      string `json:"name"`
			Namespace string `json:"namespace"`
		} `json:"podRef"`
		VolumeStats []struct {
			Name           string  `json:"name"`
			AvailableBytes *uint64 `json:"availableBytes"`
			CapacityBytes  *uint64 `json:"capacityBytes"`
			UsedBytes      *uint64 `json:"usedBytes"`
			InodesFree     *uint64 `json:"inodesFree"`
			Inodes         *uint64 `json:"inodes"`
			InodesUsed     *uint64 `json:"inodesUsed"`
		} `json:"volumeStats"`
	} `json:"pods"`
}

// PodVolumeStats returns the kubelet's reported usage for one volume of one
// pod on nodeName, read via the apiserver's node proxy subresource
// (GET /api/v1/nodes/{node}/proxy/stats/summary) -- the same source a CSI
// driver's NodeGetVolumeStats reports from, independent of the driver under
// test. volumeName is the volume's name in the pod spec (pod.Spec.Volumes),
// not the PVC's name -- this repo's own pod builders all name it
// "ms-volume".
func PodVolumeStats(nodeName, podNamespace, podName, volumeName string) (VolumeStats, error) {
	raw, err := gTestEnv.KubeInt.CoreV1().RESTClient().Get().
		Resource("nodes").
		Name(nodeName).
		SubResource("proxy").
		Suffix("stats/summary").
		DoRaw(context.TODO())
	if err != nil {
		return VolumeStats{}, fmt.Errorf("fetching /stats/summary from node %s: %w", nodeName, err)
	}

	var summary summaryResponse
	if err := json.Unmarshal(raw, &summary); err != nil {
		return VolumeStats{}, fmt.Errorf("decoding /stats/summary from node %s: %w", nodeName, err)
	}

	for _, pod := range summary.Pods {
		if pod.PodRef.Namespace != podNamespace || pod.PodRef.Name != podName {
			continue
		}
		for _, v := range pod.VolumeStats {
			if v.Name != volumeName {
				continue
			}
			var stats VolumeStats
			if v.AvailableBytes != nil {
				stats.AvailableBytes = *v.AvailableBytes
			}
			if v.CapacityBytes != nil {
				stats.CapacityBytes = *v.CapacityBytes
			}
			if v.UsedBytes != nil {
				stats.UsedBytes = *v.UsedBytes
			}
			if v.InodesFree != nil {
				stats.InodesFree = *v.InodesFree
			}
			if v.Inodes != nil {
				stats.Inodes = *v.Inodes
			}
			if v.InodesUsed != nil {
				stats.InodesUsed = *v.InodesUsed
			}
			return stats, nil
		}
		return VolumeStats{}, fmt.Errorf("node %s: pod %s/%s has no volume stats for volume %q", nodeName, podNamespace, podName, volumeName)
	}
	return VolumeStats{}, fmt.Errorf("node %s reported no stats for pod %s/%s", nodeName, podNamespace, podName)
}
