package stats

import (
	"encoding/json"
	"fmt"

	"github.com/openebs/openebs-e2e/common/k8stest"
)

//{"stats":"{\"pool\":{\"pool_created\":0,\"pool_deleted\":0},\"volume\":{\"volume_created\":0,\"volume_deleted\":0}}"}

type ConfigMapData struct {
	Stats ConfigMapStats `json:"stats"`
}

type ConfigMapStats struct {
	Pool   PoolStats   `json:"pool"`
	Volume VolumeStats `json:"volume"`
}
type PoolStats struct {
	PoolCreated int `json:"pool_created"`
	PoolDeleted int `json:"pool_deleted"`
}

type VolumeStats struct {
	VolumeCreated int `json:"volume_created"`
	VolumeDeleted int `json:"volume_deleted"`
}

func GetStatsConfigMapValue(name string, namespace string, statsType StatsType, statsAction StatsAction) (int, error) {
	configmap, err := k8stest.GetConfigMap(name, namespace)
	if err != nil {
		return 0, fmt.Errorf("failed to get configmap %s, error: %s", name, err.Error())
	}

	var stats ConfigMapStats
	data, exists := configmap.Data["stats"]
	if !exists {
		return 0, fmt.Errorf("failed to find stats data in %v", configmap.Data)
	}

	// The data is json-encoded
	err = json.Unmarshal([]byte(data), &stats)
	if err != nil {
		return 0, fmt.Errorf("Failed to unmarshall, data %s, error %s", configmap.Data["stats"], err.Error())
	}
	var val int
	switch statsType {
	case POOL:
		switch statsAction {
		case CREATED:
			val = stats.Pool.PoolCreated
		case DELETED:
			val = stats.Pool.PoolDeleted
		}
	case VOLUME:
		switch statsAction {
		case CREATED:
			val = stats.Volume.VolumeCreated
		case DELETED:
			val = stats.Volume.VolumeDeleted
		}
	}
	return val, err
}
