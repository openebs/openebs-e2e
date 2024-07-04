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
	"regexp"

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
}

func (mongo *MongoApp) MongoDump() (string, error) {
	// Convert the output to a string
	_, stderr, err := ExecuteCommandInPod(mongo.Namespace, mongo.Pod.Name, fmt.Sprintf("mongodump --host localhost --username %s --password %s --db %s --out %s", e2e_config.GetConfig().Product.MongoAuthUsername, e2e_config.GetConfig().Product.MongoAuthPassword, e2e_config.GetConfig().Product.MongoAuthDatabase, "/tmp/dump/"))
	if err != nil {
		return "", err
	}
	// regexp looks for the name of the bson dump file that we get from the dump output
	re := regexp.MustCompile(`to (/[^ ]+\.bson)`)
	matches := re.FindStringSubmatch(stderr)
	if len(matches) < 2 {
		return "", errors.New("cannot find dump file path")
	}
	path := matches[1]
	p, err := common.GetTestCaseLogsPath()
	if err != nil {
		return "", err
	}
	dumpPath := fmt.Sprintf("%s/tmp/%s.bson", p, mongo.Pod.Name)
	// Convert the output to a string
	cmd := exec.Command("kubectl", "cp", "-n", mongo.Namespace, fmt.Sprintf("%s:%s", mongo.Pod.Name, path), dumpPath)
	// Convert the output to a string
	// Get the output of the command
	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		logf.Log.Error(err, string(outputBytes))
		return "", err
	}
	// Convert the output to a string and return
	logf.Log.Info(string(outputBytes))
	return dumpPath, nil
}

func calculateChecksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
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
	logf.Log.Info("checking mongoDB application to be installed")
	ready := false

	if !mongo.Standalone {
		arbiterReady := false
		stateful, err := gTestEnv.KubeInt.AppsV1().StatefulSets(mongo.Namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: "app.kubernetes.io/name=mongodb"})
		if err != nil {
			return err
		}
		for _, ss := range stateful.Items {
			if strings.Contains(ss.Name, "arbiter") {
				arbiterReady = ss.Status.ReadyReplicas == ss.Status.Replicas &&
					ss.Status.AvailableReplicas == ss.Status.Replicas
				logf.Log.Info("StatefulSet",
					"app", "MongoDB",
					"ready", arbiterReady,
					"name", ss.Name,
					"availableReplicas", ss.Status.AvailableReplicas,
					"readyReplicas", ss.Status.ReadyReplicas,
					"currentReplicas", ss.Status.CurrentReplicas,
				)
			} else {
				ready = ss.Status.ReadyReplicas == ss.Status.Replicas &&
					ss.Status.AvailableReplicas == ss.Status.Replicas
				logf.Log.Info("StatefulSet",
					"app", "MongoDB",
					"ready", ready,
					"name", ss.Name,
					"availableReplicas", ss.Status.AvailableReplicas,
					"readyReplicas", ss.Status.ReadyReplicas,
					"currentReplicas", ss.Status.CurrentReplicas,
				)
			}
		}
		if ready && arbiterReady {
			logf.Log.Info("mongoDB HA installation and all replicas are ready")
			return nil
		}
		logf.Log.Info("not all apps are ready yet")
		time.Sleep(10 * time.Second)
	} else {

		// verify mongo deployment and pod ready
		mongoDeployName := fmt.Sprintf("%s-mongodb", mongo.ReleaseName)
		ready = WaitForDeploymentReady(mongoDeployName, mongo.Namespace, 5, defaultTimeoutSecs)
		if !ready {
			return fmt.Errorf("mongo deployment %s not ready, ready status: %v", mongoDeployName, ready)
		}
		if ready {
			pods, err := ListPod(mongo.Namespace)
			if err != nil {
				return err
			}
			for _, pod := range pods.Items {
				if strings.Contains(pod.Name, mongoDeployName) && pod.Name != mongo.Pod.Name {
					logf.Log.Info("Pod",
						"app", "MongoDB",
						"ready", ready,
						"name", pod.Name,
						"status", pod.Status.Phase,
					)
					mongo.Pod = pod
					break
				}
			}

			// wait for volume to provision
			var pvcName = mongoDeployName
			if mongo.PvcName != "" {
				pvcName = mongo.PvcName
			}
			logf.Log.Info("Verify volume provision", "pvc name", pvcName, "namespace", mongo.Namespace)
			uuid, err := VerifyVolumeProvision(pvcName, mongo.Namespace)
			if err != nil {
				return fmt.Errorf("failed to verify volume provisioning")
			}

			mongo.VolUuid = uuid
			logf.Log.Info("mongoDB standalone installation is ready")
			return nil
		}

	}
	//}
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
		builder.WriteString(fmt.Sprintf(" -p insertstart=%d", ycsb.BenchmarkParams.InsertStart))
	}
	if ycsb.BenchmarkParams.InsertCount > 0 {
		builder.WriteString(fmt.Sprintf(" -p insertcount=%d", ycsb.BenchmarkParams.InsertCount))
	}
	if ycsb.BenchmarkParams.OperationCount > 0 {
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
		builder.WriteString(fmt.Sprintf(" -p insertstart=%d", ycsb.BenchmarkParams.InsertStart))
	}
	if ycsb.BenchmarkParams.InsertCount > 0 {
		builder.WriteString(fmt.Sprintf(" -p insertcount=%d", ycsb.BenchmarkParams.InsertCount))
	}
	if ycsb.BenchmarkParams.OperationCount > 0 {
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
