package metrics

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"reflect"
	"strconv"
	"time"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/e2e_config"
	"github.com/openebs/openebs-e2e/common/k8s_portforward"
	"github.com/openebs/openebs-e2e/common/k8stest"

	"github.com/openebs/openebs-e2e/common/locations"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// Though value for polling cycle interval is 30 seconds from prometheus
// client side. But Verifying counter values in test just after
// 30 seconds can be flaky. So for tests we will use 40 seconds as this interval
// so that we can avoid hitting timing related issues.
const PollingCycleIntervalSec = 40 // in seconds

const (
	PoolStatusQuery              = "diskpool_status"
	PoolTotalSizeQuery           = "diskpool_total_size_bytes"
	PoolUsedSizeQuery            = "diskpool_used_size_bytes"
	PoolCommittedSizeQuery       = "diskpool_committed_size_bytes"
	PoolNumberOfReadOpsQuery     = "diskpool_num_read_ops"
	PoolNumberOfWriteOpsQuery    = "diskpool_num_write_ops"
	PoolTotalBytesReadQuery      = "diskpool_bytes_read"
	PoolTotalBytesWriteQuery     = "diskpool_bytes_written"
	PoolReadLatencyQuery         = "diskpool_read_latency_us"
	PoolWriteLatencyQuery        = "diskpool_write_latency_us"
	VolumeNumberOfReadOpsQuery   = "volume_num_read_ops"
	VolumeNumberOfWriteOpsQuery  = "volume_num_write_ops"
	VolumeTotalBytesReadQuery    = "volume_bytes_read"
	VolumeTotalBytesWriteQuery   = "volume_bytes_written"
	VolumeReadLatencyQuery       = "volume_read_latency_us"
	VolumeWriteLatencyQuery      = "volume_write_latency_us"
	ReplicaNumberOfReadOpsQuery  = "replica_num_read_ops"
	ReplicaNumberOfWriteOpsQuery = "replica_num_write_ops"
	ReplicaTotalBytesReadQuery   = "replica_bytes_read"
	ReplicaTotalBytesWriteQuery  = "replica_bytes_written"
	ReplicaReadLatencyQuery      = "replica_read_latency_us"
	ReplicaWriteLatencyQuery     = "replica_write_latency_us"
)

var PoolQueryList = []string{
	PoolStatusQuery,
	PoolTotalSizeQuery,
	PoolUsedSizeQuery,
	PoolCommittedSizeQuery,
}

var PoolIOStatsQueryList = []string{
	PoolNumberOfReadOpsQuery,
	PoolNumberOfWriteOpsQuery,
	PoolTotalBytesReadQuery,
	PoolTotalBytesWriteQuery,
	PoolReadLatencyQuery,
	PoolWriteLatencyQuery,
}

var VolumeIOStatsQueryList = []string{
	VolumeNumberOfReadOpsQuery,
	VolumeNumberOfWriteOpsQuery,
	VolumeTotalBytesReadQuery,
	VolumeTotalBytesWriteQuery,
	VolumeReadLatencyQuery,
	VolumeWriteLatencyQuery,
}

var ReplicaIOStatsQueryList = []string{
	ReplicaNumberOfReadOpsQuery,
	ReplicaNumberOfWriteOpsQuery,
	ReplicaTotalBytesReadQuery,
	ReplicaTotalBytesWriteQuery,
	ReplicaReadLatencyQuery,
	ReplicaWriteLatencyQuery,
}

type Metric struct {
	Metric_Name string `json:"__name__"`
	Container   string `json:"container"`
	Endpoint    string `json:"endpoint"`
	Instance    string `json:"instance"`
	Job         string `json:"job"`
	Name        string `json:"name"`
	PVName      string `json:"pv_name"`
	Namespace   string `json:"namespace"`
	Node        string `json:"node"`
	Pod         string `json:"pod"`
	Service     string `json:"service"`
}

// Throughout this file Resource is a generic name
// which can be used for Diskpool, volume and replica
type ResourceStatus struct {
	Status string `json:"status"`
	Data   Data   `json:"data"`
}

type Resource struct {
	Metric Metric        `json:"metric"`
	Value  []interface{} `json:"value"`
}

type Data struct {
	ResultType string     `json:"resultType"`
	Result     []Resource `json:"result"`
}

func installPrometheus(namespace string) error {
	// add helm prometheus repo
	err := addPrometheusHelmRepo()
	if err != nil {
		return fmt.Errorf("failed to add prometheus helm repo: : Error: %v", err)
	}

	// update helm prometheus repo
	err = updatePrometheusHelmRepo()
	if err != nil {
		return fmt.Errorf("failed to update helm repo: Error: %v", err)
	}

	cmdString := "prometheus.prometheusSpec.serviceMonitorSelector.matchLabels=null,nodeExporter.enabled=false,prometheus.service.type=NodePort,alertmanager.enabled=false,grafana.enabled=false"
	// install prometheus using helm
	cmd := exec.Command(
		"helm", "install", "e2e-prometheus", "prometheus-community/kube-prometheus-stack",
		"--set", cmdString, "-n", namespace)
	logf.Log.Info("About to execute: helm install e2e-prometheus prometheus-community/kube-prometheus-stack", " --set ",
		cmdString, "namespace", namespace)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install prometheus using helm chart: namespace: %s Output: %s : Error: %v", namespace, out, err)
	}
	return nil
}

func postInstallationSteps() error {
	err := CheckBoltInstallReady()
	if err != nil {
		return fmt.Errorf("prometheus/product deployment not in ready state, error: %v", err)
	}

	configDir := locations.GetE2EServiceMonitorPath()
	err = ApplyServiceMonitorYaml(
		"service_monitor.yaml",
		configDir,
		common.NSMayastor(),
	)

	if err != nil {
		return fmt.Errorf("failed to apply service monitor, directory: %s, namespace: %s , error: %v", configDir, common.NSMayastor(), err)
	}

	nodeIP := k8stest.GetMayastorNodeIPAddresses()

	poolsInCluster, err := k8stest.ListMsPools()
	if err != nil {
		return fmt.Errorf("failed to list pools in cluster, error: %v", err)
	}

	if len(poolsInCluster) == 0 {
		return fmt.Errorf("no disk pools found in cluster")
	}

	// Wait for disk pool status metrics
	for _, query := range PoolQueryList {
		if ready, err := VerifyResourceMetricGeneration(len(poolsInCluster), nodeIP, query, 5, 240); err != nil || !ready {
			return fmt.Errorf("failed to generate %s metrics for all pools, error: %v", query, err)
		}
	}

	// Wait for disk pool io stats metrics
	for _, query := range PoolIOStatsQueryList {
		if ready, err := VerifyResourceMetricGeneration(len(poolsInCluster), nodeIP, query, 5, 240); err != nil || !ready {
			return fmt.Errorf("failed to generate %s metrics for all pools, error: %v", query, err)
		}
	}

	return nil
}

// prometheus helm repo add
//
//	helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
func addPrometheusHelmRepo() error {
	cmd := exec.Command(
		"helm", "repo", "add", "prometheus-community",
		"https://prometheus-community.github.io/helm-charts",
	)
	logf.Log.Info("About to execute: helm repo add prometheus-community https://prometheus-community.github.io/helm-charts")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to add prometheus helm repo: Output: %s : Error: %v", out, err)
	}
	return nil
}

// prometheus helm repo update
func updatePrometheusHelmRepo() error {

	cmd := exec.Command(
		"helm", "repo", "update",
	)
	logf.Log.Info("About to execute: helm repo update")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to update helm repo: Output: %s : Error: %v", out, err)
	}
	return nil
}

func uninstallPrometheus(namespace string) error {
	cmd := exec.Command(
		"helm", "uninstall", "e2e-prometheus",
		"-n", namespace,
	)
	logf.Log.Info("About to execute: helm uninstall e2e-prometheus",
		"namespace", namespace)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to uninstall prometheus using helm chart: namespace: %s Output: %s : Error: %v", namespace, out, err)
	}
	return nil
}

func postUnInstallationSteps() error {
	configDir := locations.GetE2EServiceMonitorPath()

	err := DeleteServiceMonitorYaml(
		"service_monitor.yaml",
		configDir,
		common.NSMayastor(),
	)
	if err != nil {
		return fmt.Errorf("failed to delete service monitor, error: %v", err)
	}

	// err = k8sinstall.RollbackBoltHelmRelease()
	// if err != nil {
	// 	return fmt.Errorf("failed to rollback helm release, error: %v", err)
	// }

	return nil
}

func ApplyServiceMonitorYaml(filename string, dir string, namespace string) error {
	cmd := exec.Command("kubectl", "apply", "-f", filename, "-n", namespace)
	cmd.Dir = dir
	logf.Log.Info("kubectl apply ", "yaml file", filename, "path", cmd.Dir, "namespace", namespace)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to apply yaml file %s : Output: %s : Error: %v", filename, out, err)
	}
	return nil
}

func DeleteServiceMonitorYaml(filename string, dir string, namespace string) error {
	cmd := exec.Command("kubectl", "delete", "-f", filename, "-n", namespace)
	cmd.Dir = dir
	logf.Log.Info("kubectl delete ", "yaml file", filename, "path", cmd.Dir, "namespace", namespace)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to delete yaml file %s : Output: %s : Error: %v", filename, out, err)
	}
	return nil
}

func CheckBoltInstallReady() error {
	ready, err := k8stest.MayastorReady(2, 540)
	if err != nil {
		return err
	}
	if !ready {
		return fmt.Errorf("prometheus installation is not ready")
	}

	ready = k8stest.ControlPlaneReady(10, 180)
	if !ready {
		return fmt.Errorf("control plane installation is not ready")
	}

	return err
}

func VerifyResourceMetricGeneration(resourceCount int, node []string, metricsQuery string, sleepTime int, duration int) (bool, error) {
	count := (duration + sleepTime - 1) / sleepTime
	ready := false
	var err error
	var resourceStatusMetrics ResourceStatus
	for ix := 0; ix < count && !ready; ix++ {
		time.Sleep(time.Duration(sleepTime) * time.Second)
		resourceStatusMetrics, err = GetResourceMetrics(metricsQuery, node)
		ready = resourceCount == len(resourceStatusMetrics.Data.Result)
	}
	if err != nil {
		return ready, fmt.Errorf("failed to get resource metrics, error: %v", err)
	}
	return ready, err
}

func GetResourceMetrics(query string, address []string) (ResourceStatus, error) {
	if len(address) == 0 {
		return ResourceStatus{}, fmt.Errorf("nodes not found")
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
		}
		req.Header.Add("Accept", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			logf.Log.Info("Error while making GET request", "url", url, "error", err)
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
		return ResourceStatus{}, err
	}
	var response ResourceStatus
	err = json.Unmarshal(jsonResponse, &response)
	if err != nil {
		logf.Log.Info("Failed to unmarshal", "string", string(jsonResponse))
		return ResourceStatus{}, err
	}
	return response, nil
}

func GetDiskPoolMetricValue(poolName string, query string, address []string) (string, error) {
	var metricValue interface{}
	poolMetrics, err := GetResourceMetrics(query, address)
	if err != nil {
		return "", err
	}
	for _, pool := range poolMetrics.Data.Result {
		if poolName == pool.Metric.Name {
			metricValue = pool.Value[1]
			break
		}

	}
	return fmt.Sprint(metricValue), err
}

func GetVolumeMetricValue(volumeName string, query string, address []string) (string, error) {
	var metricValue interface{}
	volumeMetrics, err := GetResourceMetrics(query, address)
	if err != nil {
		return "", err
	}
	for _, vol := range volumeMetrics.Data.Result {
		if volumeName == vol.Metric.PVName {
			metricValue = vol.Value[1]
			break
		}

	}
	return fmt.Sprint(metricValue), err
}

func GetReplicaMetricValue(replicaName string, query string, address []string) (string, error) {
	var metricValue interface{}
	replicaMetrics, err := GetResourceMetrics(query, address)
	if err != nil {
		return "", err
	}
	for _, vol := range replicaMetrics.Data.Result {
		if replicaName == vol.Metric.Name {
			metricValue = vol.Value[1]
			break
		}

	}
	return fmt.Sprint(metricValue), err
}

func verifyMetrics(name string, node []string, queryList []string, verifyMetricFunc func(string, string, []string) (bool, error), logFields map[string]interface{}) (bool, error) {
	var results []bool
	var errs common.ErrorAccumulator

	for _, query := range queryList {
		result, err := verifyMetricFunc(name, query, node)
		if err != nil {
			errs.Accumulate(fmt.Errorf("failed to verify %s metrics for %s : %v", query, name, err))
		}
		results = append(results, result)
	}

	ready := true
	for _, result := range results {
		ready = ready && result
	}

	if errs.GetError() != nil {
		return ready, fmt.Errorf("failed to get %s metrics, error: %v", name, errs.GetError())
	}

	logFields["name"] = name
	for i, result := range results {
		logFields[queryList[i]] = result
	}

	logf.Log.Info("io stats", "metrics", logFields)

	return ready, nil
}

func VerifyPoolMetrics(poolName string, node []string) (bool, error) {
	logFields := make(map[string]interface{})
	return verifyMetrics(poolName, node, PoolQueryList, VerifyDiskPoolMetric, logFields)
}

func VerifyPoolIOStatsMetrics(poolName string, node []string) (bool, error) {
	logFields := make(map[string]interface{})
	return verifyMetrics(poolName, node, PoolIOStatsQueryList, VerifyDiskPoolMetric, logFields)
}

func VerifyVolumeIOStatsMetrics(volumeName string, node []string) (bool, error) {
	logFields := make(map[string]interface{})
	return verifyMetrics(volumeName, node, VolumeIOStatsQueryList, VerifyVolumeMetric, logFields)
}

func VerifyReplicaIOStatsMetrics(replicaName string, node []string) (bool, error) {
	logFields := make(map[string]interface{})
	return verifyMetrics(replicaName, node, ReplicaIOStatsQueryList, VerifyReplicaMetric, logFields)
}

func VerifyDiskPoolMetric(poolName string, queryParameter string, address []string) (bool, error) {
	poolStatusMetrics, err := GetResourceMetrics(queryParameter, address)
	if err != nil {
		return false, err
	}
	for _, pool := range poolStatusMetrics.Data.Result {
		if poolName == pool.Metric.Name {
			return true, nil
		}
	}
	return false, nil
}

func VerifyVolumeMetric(volumeName string, queryParameter string, address []string) (bool, error) {
	volumeStatusMetrics, err := GetResourceMetrics(queryParameter, address)
	if err != nil {
		return false, err
	}
	for _, vol := range volumeStatusMetrics.Data.Result {
		if volumeName == vol.Metric.PVName {
			return true, nil
		}
	}
	return false, nil
}

func VerifyReplicaMetric(replicaName string, queryParameter string, address []string) (bool, error) {
	replicaStatusMetrics, err := GetResourceMetrics(queryParameter, address)
	if err != nil {
		return false, err
	}
	for _, rep := range replicaStatusMetrics.Data.Result {
		if replicaName == rep.Metric.Name {
			return true, nil
		}
	}
	return false, nil
}

func VerifyDiskPoolStatusMetricsOnline(actualPoolOnline int, address []string) (bool, error) {
	poolStatusMetrics, err := GetResourceMetrics(PoolStatusQuery, address)
	if err != nil {
		return false, err
	}
	// count all pools which are online
	// metrics value = "1" means pool is online else it's offline
	// Value[0] will hold time series value for metrics
	// Value[1] will hold actual pool status metrics value
	// example: "value":[1653301161.138,"1"]
	onlineCount := 0
	for _, pool := range poolStatusMetrics.Data.Result {
		if pool.Value[1] == "1" {
			onlineCount++
		}
	}
	return actualPoolOnline == onlineCount, nil
}

func VerifyDiskPoolCapacityMetrics(capacity map[string]uint64, address []string) (bool, error) {
	var metricsTotalPoolSize = make(map[string]uint64)
	diskPoolTotalSize, err := GetResourceMetrics(PoolTotalSizeQuery, address)
	if err != nil {
		return false, err
	}
	// Value[0] will hold time series value for metrics
	// Value[1] will hold actual pool capacity metrics value
	for _, pool := range diskPoolTotalSize.Data.Result {
		capacity, err := strconv.ParseUint(pool.Value[1].(string), 10, 64)
		if err != nil {
			return false, fmt.Errorf("failed to parse pool capacity string %s to unsigned numbers: error: %v", pool.Value[1].(string), err)
		}
		metricsTotalPoolSize[pool.Metric.Name] = capacity
	}
	//compare total capacity of pools across crd and metrics
	return reflect.DeepEqual(capacity, metricsTotalPoolSize), err
}

func VerifyDiskPoolUsedSizeMetrics(crdPoolUsedSize map[string]uint64, address []string) (bool, error) {
	var metricsUsedPoolSize = make(map[string]uint64)
	diskPoolTotalSize, err := GetResourceMetrics(PoolUsedSizeQuery, address)
	if err != nil {
		return false, err
	}

	// Value[0] will hold time series value for metrics
	// Value[1] will hold actual pool used size metrics value
	for _, pool := range diskPoolTotalSize.Data.Result {
		used, err := strconv.ParseUint(pool.Value[1].(string), 10, 64)
		if err != nil {
			return false, fmt.Errorf("failed to parse pool used size string %s to unsigned numbers: error: %v", pool.Value[1].(string), err)
		}
		metricsUsedPoolSize[pool.Metric.Name] = used
	}
	//compare total used size of pools across crd and metrics
	return reflect.DeepEqual(crdPoolUsedSize, metricsUsedPoolSize), err
}

func VerifyDiskPoolCommittedSizeMetrics(crdPoolCommittedSize map[string]uint64, address []string) (bool, error) {
	var metricsCommittedPoolSize = make(map[string]uint64)
	diskPoolTotalSize, err := GetResourceMetrics(PoolCommittedSizeQuery, address)
	if err != nil {
		return false, err
	}

	// Value[0] will hold time series value for metrics
	// Value[1] will hold actual pool committed size metrics value
	for _, pool := range diskPoolTotalSize.Data.Result {
		committed, err := strconv.ParseUint(pool.Value[1].(string), 10, 64)
		if err != nil {
			return false, fmt.Errorf("failed to parse pool committed size string %s to unsigned numbers: error: %v", pool.Value[1].(string), err)
		}
		metricsCommittedPoolSize[pool.Metric.Name] = committed
	}
	//compare total committed size of pools across crd and metrics
	return reflect.DeepEqual(crdPoolCommittedSize, metricsCommittedPoolSize), err
}

func MetricsTestTeardown() error {
	//  Uninstall prometheus stack
	err := uninstallPrometheus(common.NSMayastor())
	if err != nil {
		return fmt.Errorf("failed to uninstall prometheus, error %v", err)
	}

	// Post uninstall cleanup
	err = postUnInstallationSteps()
	if err != nil {
		return fmt.Errorf("failed in post prometheus uninstall steps, error %v", err)
	}

	err = CheckBoltInstallReady()
	if err != nil {
		return fmt.Errorf("product deployment not in ready state, error %v", err)
	}

	return err
}

func MetricsTestSetup() error {
	// chartDir, err := locations.GetHelmChartsDir()
	// if err != nil {
	// 	return fmt.Errorf("failed to get helm chart directory, Error %v", err)
	// }

	// // Update helm chart with metrics poll time
	// err = k8sinstall.UpdateMetricsPollTimeInHelmChart(chartDir)
	// if err != nil {
	// 	return fmt.Errorf("failed to update helm chart with metrics polling interval, Error %v", err)
	// }

	// // Upgrade helm release
	// err = k8sinstall.UpgradeHelmRelease(chartDir)
	// if err != nil {
	// 	return fmt.Errorf("failed to upgrade helm chart with metrics polling interval, Error %v", err)
	// }

	// Install prometheus stack
	err := installPrometheus(common.NSMayastor())
	if err != nil {
		return fmt.Errorf("failed to install prometheus, error %v", err)
	}

	// Post install check and configuration
	err = postInstallationSteps()
	if err != nil {
		return fmt.Errorf("failed to verify/config after prometheus installation, error %f", err)
	}

	return err
}

// return timeout which is three times of MetricsPollingInterval
func GetMetricsTimeoutSec() time.Duration {
	timeout, err := time.ParseDuration(e2e_config.GetConfig().Product.MetricsPollingInterval)
	if err != nil {
		panic(fmt.Errorf("failed to parse timeout %s string , error: %v",
			e2e_config.GetConfig().Product.MetricsPollingInterval,
			err))
	}
	return timeout * 3
}
