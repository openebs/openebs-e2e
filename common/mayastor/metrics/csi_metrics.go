package metrics

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/e2e_config"
	"github.com/openebs/openebs-e2e/common/k8s_portforward"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var DefWaitTimeout = 180 //in seconds

const (
	CsiVolumeAvailableBytes = "kubelet_volume_stats_available_bytes"
	CsiVolumeCapacityBytes  = "kubelet_volume_stats_capacity_bytes"
	CsiVolumeUsedBytes      = "kubelet_volume_stats_used_bytes"
	InodeCapacity           = "kubelet_volume_stats_inodes"
	InodeMetricFree         = "kubelet_volume_stats_inodes_free"
	InodeMetricUsed         = "kubelet_volume_stats_inodes_used"
)

type CsiVolumeMetrics struct {
	Status string  `json:"status"`
	Data   CsiData `json:"data"`
}
type CsiMetric struct {
	Name                  string `json:"__name__"`
	Endpoint              string `json:"endpoint"`
	Instance              string `json:"instance"`
	Job                   string `json:"job"`
	MetricsPath           string `json:"metrics_path"`
	Namespace             string `json:"namespace"`
	Node                  string `json:"node"`
	Persistentvolumeclaim string `json:"persistentvolumeclaim"`
	Service               string `json:"service"`
}

type CsiVolume struct {
	Metric CsiMetric     `json:"metric"`
	Value  []interface{} `json:"value"`
}

type CsiData struct {
	ResultType string      `json:"resultType"`
	Result     []CsiVolume `json:"result"`
}

// GetCsiVolumeMetrics query Prometheus rest endpoint to get CSI volume metrics
// query will be
// 1. kubelet_volume_stats_available_bytes
// 2. kubelet_volume_stats_capacity_bytes
// 3. kubelet_volume_stats_used_bytes
// address will be array of product(mayastor/bolt) node ip address
func GetCsiVolumeMetrics(query string, address []string) (CsiVolumeMetrics, error) {
	if len(address) == 0 {
		return CsiVolumeMetrics{}, fmt.Errorf("product nodes not found")
	}
	var jsonResponse []byte
	var err error
	for _, addr := range address {
		url := fmt.Sprintf("http://%s/api/v1/query?query=%s",
			k8s_portforward.TryPortForwardNode(addr, e2e_config.GetConfig().Product.PrometheusPort),
			query)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			logf.Log.Info("Error in GET request", "node IP", addr, "url", url, "error", err)
			return CsiVolumeMetrics{}, err
		}
		req.Header.Add("Accept", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			logf.Log.Info("Error while making GET request", "url", url, "error", err)
			return CsiVolumeMetrics{}, err
		} else if resp.StatusCode != 200 {
			return CsiVolumeMetrics{}, fmt.Errorf("prometheus api reponse code is not 200, response code reveived is %d", resp.StatusCode)
		} else {
			defer resp.Body.Close()
			jsonResponse, err = io.ReadAll(resp.Body)
			if err != nil {
				logf.Log.Info("Error while reading data", "error", err)
			} else {
				break
			}
		}

	}
	if err != nil {
		return CsiVolumeMetrics{}, err
	}
	var response CsiVolumeMetrics
	err = json.Unmarshal(jsonResponse, &response)
	if err != nil {
		logf.Log.Info("Failed to unmarshal", "string", string(jsonResponse))
		return CsiVolumeMetrics{}, err
	}
	return response, nil
}

// GetCsiVolumeMetricValue return used or available or capacity metric value depending on query for a volume
func GetCsiVolumeMetricValue(pvcName string, query string, address []string) (string, error) {
	var metricValue interface{}
	var volumeExist bool
	volumeMetrics, err := GetCsiVolumeMetrics(query, address)
	if err != nil {
		return "", err
	}
	for _, volume := range volumeMetrics.Data.Result {
		if pvcName == volume.Metric.Persistentvolumeclaim {
			volumeExist = true
			metricValue = volume.Value[1]
			break
		}
	}
	if metricValue == nil || !volumeExist {
		return "", fmt.Errorf("failed to find volume %s or %s metrics value for volume %s", pvcName, query, pvcName)
	}
	return fmt.Sprint(metricValue), err
}

// VerifyCsiVolumeMetricsGenerated verify available, used and capacity metrics generated for volume
func VerifyCsiVolumeMetricsGenerated(pvcName string, node []string) (bool, error) {
	ready := false
	var err error
	var errs common.ErrorAccumulator
	var availableMetrics, capacityMetrics, usedMetrics, capacityInode, freeInode, usedInode bool
	availableMetrics, err = VerifyCsiVolumeMetricGenerated(pvcName, CsiVolumeAvailableBytes, node)
	if err != nil {
		errs.Accumulate(fmt.Errorf("failed to verify csi volume available bytes metrics for pvc %s, error %v", pvcName, err))
	}
	capacityMetrics, err = VerifyCsiVolumeMetricGenerated(pvcName, CsiVolumeCapacityBytes, node)
	if err != nil {
		errs.Accumulate(fmt.Errorf("failed to verify csi volume capacity bytes metrics for pvc %s, error %v", pvcName, err))
	}
	usedMetrics, err = VerifyCsiVolumeMetricGenerated(pvcName, CsiVolumeUsedBytes, node)
	if err != nil {
		errs.Accumulate(fmt.Errorf("failed to verify csi volume used bytes metrics for pvc %s, error %v", pvcName, err))
	}
	capacityInode, err = VerifyCsiVolumeMetricGenerated(pvcName, InodeCapacity, node)
	if err != nil {
		errs.Accumulate(fmt.Errorf("failed to verify inode metric for pvc %s, error %v", pvcName, err))
	}
	freeInode, err = VerifyCsiVolumeMetricGenerated(pvcName, InodeMetricFree, node)
	if err != nil {
		errs.Accumulate(fmt.Errorf("failed to verify inode free metric for pvc %s, error %v", pvcName, err))
	}
	usedInode, err = VerifyCsiVolumeMetricGenerated(pvcName, InodeMetricUsed, node)
	if err != nil {
		errs.Accumulate(fmt.Errorf("failed to verify inode used metric for pvc %s, error %v", pvcName, err))
	}
	if errs.GetError() != nil {
		return ready, fmt.Errorf("failed to get csi volume metrics, error: %v", errs.GetError())
	}
	logf.Log.Info("metrics", "available", availableMetrics, "capacity", capacityMetrics, "used", usedMetrics, "capacityInode", capacityInode, "freeInode", freeInode, "usedInode", usedInode)
	return availableMetrics && capacityMetrics && usedMetrics && capacityInode && freeInode && usedInode, errs.GetError()
}

// VerifyCsiVolumeMetricGenerated verify respective metrics generated for volume depending on query parameter
func VerifyCsiVolumeMetricGenerated(pvcName string, queryParameter string, address []string) (bool, error) {
	csiVolumeMetric, err := GetCsiVolumeMetrics(queryParameter, address)
	if err != nil {
		return false, err
	}
	for _, volume := range csiVolumeMetric.Data.Result {
		if pvcName == volume.Metric.Persistentvolumeclaim {
			return true, nil
		}
	}
	return false, nil

}

func VerifyWriteDataSize(usedSizeBefore string, usedSizeAfter string, fioSizeInMb int) (bool, error) {
	usedSizeBeforeWrite, err := strconv.ParseUint(usedSizeBefore, 10, 64)
	if err != nil {
		return false, fmt.Errorf("failed to convert size string %s to uint, error %v", usedSizeBefore, err)
	}
	usedSizeAfterWrite, err := strconv.ParseUint(usedSizeAfter, 10, 64)
	if err != nil {
		return false, fmt.Errorf("failed to convert size string %s to uint, error %v", usedSizeAfter, err)
	}
	overhead := (usedSizeAfterWrite - usedSizeBeforeWrite) - uint64((fioSizeInMb * 1024 * 1024))
	if overhead == 4096 {
		return true, nil
	}
	return false, nil
}

// GetSizeInUnit64 return size in unit64
func GetSizeInUnit64(size string) (uint64, error) {
	sizeInUnit64, err := strconv.ParseUint(size, 10, 64)
	if err != nil {
		return sizeInUnit64, fmt.Errorf("failed to convert size string %s to uint64, error %v", size, err)
	}
	return sizeInUnit64, nil
}
