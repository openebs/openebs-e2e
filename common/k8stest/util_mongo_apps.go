package k8stest

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/openebs/openebs-e2e/common"

	"io"
	"os"
	"os/exec"

	"github.com/openebs/openebs-e2e/common/e2e_config"
	appsv1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

type MongoApp struct {
	Namespace    string
	Pod          coreV1.Pod
	VolUuid      string
	ReleaseName  string
	ReplicaCount int
	ScName       string
	Standalone   bool
	PvcName      string
	Ycsb         bool
	StsName      string
}

const (
	defaultPodTimeoutSecs       = 120
	defaultMongodbStsimeoutSecs = 180
)

func (mongo *MongoApp) MongoDump() (string, error) {
	log := logf.Log.WithName("mongo-dump")

	dumpFile := "/tmp/mongo.dump.archive"
	cmd := fmt.Sprintf(
		"rm -f %s && mongodump --host localhost --db %s --archive=%s",
		dumpFile,
		e2e_config.GetConfig().Product.MongoAuthDatabase,
		dumpFile,
	)

	_, stderr, err := ExecuteCommandInPod(
		mongo.Namespace,
		mongo.Pod.Name,
		cmd,
	)
	if err != nil {
		log.Error(err, "mongodump failed", "stderr", stderr)
		return "", err
	}

	p, err := common.GetTestCaseLogsPath()
	if err != nil {
		return "", err
	}

	localDumpPath := fmt.Sprintf("%s/tmp/%s-mongo.dump.archive", p, mongo.Pod.Name)

	copyCmd := exec.Command(
		"kubectl",
		"cp",
		"-n", mongo.Namespace,
		fmt.Sprintf("%s:%s", mongo.Pod.Name, dumpFile),
		localDumpPath,
	)

	output, err := copyCmd.CombinedOutput()
	if err != nil {
		log.Error(err, "kubectl cp failed", "output", string(output))
		return "", err
	}

	log.Info("Mongo dump copied", "path", localDumpPath)
	return localDumpPath, nil
}

func calculateChecksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	//nolint:errcheck
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func (mongo *MongoApp) CompareBSONChecksums(file1, file2 string) bool {
	checksum1, err := calculateChecksum(file1)
	if err != nil {
		logf.Log.Error(err, "error calculating checksum", "file", file1)
		return false
	}

	checksum2, err := calculateChecksum(file2)
	if err != nil {
		logf.Log.Error(err, "error calculating checksum", "file", file2)
		return false
	}

	return checksum1 == checksum2
}

func (mongo *MongoApp) MongoInstallReady() error {
	logf.Log.Info("Checking MongoDB application installation")

	err := WaitForStsReady(mongo.StsName, mongo.Namespace, time.Duration(defaultMongodbStsimeoutSecs)*time.Second)
	if err != nil {
		return fmt.Errorf("mongoDB StatefulSet %s not ready: %v", mongo.StsName, err)
	}

	// Get all pods in the StatefulSet
	pods, err := GetStsPodNames(mongo.StsName, mongo.Namespace)
	if err != nil {
		return fmt.Errorf("failed to list MongoDB StatefulSet pods: %v", err)
	}

	// Wait for volume provisioning and confirm MongoDB pods are ready
	for _, pod := range pods {
		var pvcName, uuid string

		pvcName, err = GetPvcNameFromPod(pod, mongo.Namespace)
		if err != nil {
			return fmt.Errorf("failed to get PVC name from pod %s: %v", pod, err)
		}
		if pvcName == "" {
			return fmt.Errorf("pvc name not found for pod %s", pod)
		}

		logf.Log.Info("Verifying volume provisioning", "pvcName", pvcName, "namespace", mongo.Namespace)
		uuid, err = VerifyVolumeProvision(pvcName, mongo.Namespace)
		if err != nil {
			return fmt.Errorf("failed to verify volume provisioning for %s: %v", pvcName, err)
		}

		logf.Log.Info("MongoDB pod ready",
			"pod", pod,
			"pvcName", pvcName,
			"volumeUUID", uuid,
		)

		mongo.Pod.Name = pod
		mongo.VolUuid = uuid
	}

	logf.Log.Info("MongoDB installation is ready", "statefulset", mongo.StsName)
	return nil
}

const (
	ycsbImage               = "openebs/e2e-ycsb:v1.0.0"
	appWorkdir              = "/app/ycsb"
	defaultThreadCount      = 8
	defaultRecordCount      = 5_000
	defaultWorkloadFileName = "workloada"
)

type YcsbApp struct {
	BenchmarkParams BenchmarkParams
	MongoConnUrl    string
	Name            string
	Namespace       string
	NodeSelector    string
	PodName         string
}

type BenchmarkParams struct {
	InsertCount      int
	InsertStart      int
	OperationCount   int
	RecordCount      int
	ThreadCount      int
	WorkloadFileName string
}

func NewYCSB() *YcsbApp {
	return &YcsbApp{
		BenchmarkParams: BenchmarkParams{
			RecordCount:      defaultRecordCount,
			ThreadCount:      defaultThreadCount,
			WorkloadFileName: defaultWorkloadFileName,
		},
	}
}

func (ycsb *YcsbApp) DeployYcsbApp(namePrefix string) (string, string, error) {
	deploymentsClient := gTestEnv.KubeInt.AppsV1().Deployments(ycsb.Namespace)
	name := fmt.Sprintf("%s-ycsb", namePrefix)
	helper := int32(1)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"app":  "ycsb",
				"role": "client",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &helper,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":  "ycsb",
					"role": "client",
				},
			},
			Template: coreV1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":  "ycsb",
						"role": "client",
					},
				},
				Spec: coreV1.PodSpec{
					Containers: []coreV1.Container{
						{
							Name:            "benchmark",
							Image:           ycsbImage,
							ImagePullPolicy: coreV1.PullIfNotPresent,
							Command:         []string{"sleep", "3650d"},
						},
					},
				},
			},
		},
	}

	if ycsb.NodeSelector != "" {
		deployment.Spec.Template.Spec.NodeSelector = map[string]string{
			"kubernetes.io/hostname": ycsb.NodeSelector,
		}
	}

	// Create Deployment
	logf.Log.Info("creating YCSB deployment", "namespace", ycsb.Namespace)
	ycsbDep, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return "", "", fmt.Errorf("failed to create ycsb deployment , error: %v", err)
	}

	// verify ycsb deployment and pod ready
	isReady := WaitForDeploymentReady(ycsbDep.Name, ycsb.Namespace, 5, 120)
	if !isReady {
		return "", "", fmt.Errorf("ycsb deployment %s not ready, ready status: %v", ycsb.Name, isReady)
	}
	pods, err := ListPod(ycsb.Namespace)
	if err != nil {
		return "", "", err
	}
	for _, pod := range pods.Items {
		if strings.Contains(pod.Name, ycsbDep.Name) {
			return name, pod.Name, nil
		}
	}
	return "", "", nil
}

func (ycsb *YcsbApp) UndeployYcsbApp() error {
	deploymentsClient := gTestEnv.KubeInt.AppsV1().Deployments(ycsb.Namespace)
	err := deploymentsClient.Delete(context.TODO(), ycsb.Name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	to := 120
	for i := 0; i < to; i++ {
		p, err := GetPod(ycsb.PodName, ycsb.Namespace)
		if p != nil {
			logf.Log.Info("ycsb pod still visible", "name", p.Name, "status", p.Status.Phase)
			time.Sleep(10 * time.Second)
			to -= 10
			continue
		}
		if err == nil || to < 0 {
			return errors.New("failed to undeploy ycsb deployment")
		}
	}
	logf.Log.Info("ycsb deployment removed", "pod name", ycsb.PodName, "namespace", ycsb.Namespace)
	return nil
}

func (ycsb *YcsbApp) LoadYcsbApp() error {
	var wg sync.WaitGroup
	outputChan := make(chan string, 1)
	errChan := make(chan error, 1)
	var builder strings.Builder
	//nolint:staticcheck // QF1012 - string concatenation is more readable for command construction
	builder.WriteString(fmt.Sprintf(
		"%s/bin/ycsb.sh load mongodb -P %s/workloads/%s -p %s -p recordcount=%d -p threadcount=%d",
		appWorkdir,
		appWorkdir,
		ycsb.BenchmarkParams.WorkloadFileName,
		ycsb.MongoConnUrl,
		ycsb.BenchmarkParams.RecordCount,
		ycsb.BenchmarkParams.ThreadCount,
	))
	if ycsb.BenchmarkParams.InsertStart > 0 {
		//nolint:staticcheck // QF1012 - string concatenation is more readable for command construction
		builder.WriteString(fmt.Sprintf(" -p insertstart=%d", ycsb.BenchmarkParams.InsertStart))
	}
	if ycsb.BenchmarkParams.InsertCount > 0 {
		//nolint:staticcheck // QF1012 - string concatenation is more readable for command construction
		builder.WriteString(fmt.Sprintf(" -p insertcount=%d", ycsb.BenchmarkParams.InsertCount))
	}
	if ycsb.BenchmarkParams.OperationCount > 0 {
		//nolint:staticcheck // QF1012 - string concatenation is more readable for command construction
		builder.WriteString(fmt.Sprintf(" -p operationcount=%d", ycsb.BenchmarkParams.OperationCount))
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		output, _, err := ExecuteCommandInPod(ycsb.Namespace, ycsb.PodName, builder.String())
		if err != nil {
			errChan <- err
			return
		}
		outputChan <- output
	}()
	wg.Wait()
	close(outputChan)
	close(errChan)

	err, hasError := <-errChan
	if hasError {
		return err
	}

	_, hasResult := <-outputChan
	if hasResult {
		logf.Log.Info("ycsb load completed")
	}
	return nil
}

func (ycsb *YcsbApp) RunYcsbApp(result *string) error {
	var wg sync.WaitGroup
	outputChan := make(chan string, 1)
	errChan := make(chan error, 1)
	var builder strings.Builder
	//nolint:staticcheck // QF1012 - string concatenation is more readable for command construction
	builder.WriteString(fmt.Sprintf(
		"%s/bin/ycsb.sh run mongodb -P %s/workloads/%s -p %s -p recordcount=%d -p threadcount=%d",
		appWorkdir,
		appWorkdir,
		ycsb.BenchmarkParams.WorkloadFileName,
		ycsb.MongoConnUrl,
		ycsb.BenchmarkParams.RecordCount,
		ycsb.BenchmarkParams.ThreadCount,
	))
	if ycsb.BenchmarkParams.InsertStart > 0 {
		//nolint:staticcheck // QF1012 - string concatenation is more readable for command construction
		builder.WriteString(fmt.Sprintf(" -p insertstart=%d", ycsb.BenchmarkParams.InsertStart))
	}
	if ycsb.BenchmarkParams.InsertCount > 0 {
		//nolint:staticcheck // QF1012 - string concatenation is more readable for command construction
		builder.WriteString(fmt.Sprintf(" -p insertcount=%d", ycsb.BenchmarkParams.InsertCount))
	}
	if ycsb.BenchmarkParams.OperationCount > 0 {
		//nolint:staticcheck // QF1012 - string concatenation is more readable for command construction
		builder.WriteString(fmt.Sprintf(" -p operationcount=%d", ycsb.BenchmarkParams.OperationCount))
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		output, _, err := ExecuteCommandInPod(ycsb.Namespace, ycsb.PodName, builder.String())
		if err != nil {
			errChan <- err
			return
		}
		outputChan <- output
	}()
	wg.Wait()
	close(outputChan)
	close(errChan)

	err, hasError := <-errChan
	if hasError {
		return err
	}

	out, hasResult := <-outputChan
	if hasResult {
		*result = out
	}
	return nil
}
