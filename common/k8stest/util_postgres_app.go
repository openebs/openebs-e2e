package k8stest

import (
	"context"
	"fmt"
	"github.com/openebs/openebs-e2e/common/custom_resources"
	"github.com/openebs/openebs-e2e/common/e2e_config"
	"strings"
	"time"

	batchV1 "k8s.io/api/batch/v1"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

type PostgresApp struct {
	Namespace    string
	Pod          coreV1.Pod
	VolUuid      string
	ReleaseName  string
	ReplicaCount int
	ScName       string
	Standalone   bool
	PvcName      string
	PgBench      PgBenchApp
}

func (psql *PostgresApp) PostgresInstallReady() error {
	logf.Log.Info("checking postgres application to be installed")
	ready := false
	counter := 12
	if psql.Standalone {
		stateful, err := gTestEnv.KubeInt.AppsV1().StatefulSets(psql.Namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: "app.kubernetes.io/name=postgresql"})
		if err != nil {
			return err
		}
		for _, ss := range stateful.Items {
			ready = ss.Status.ReadyReplicas == ss.Status.Replicas &&
				ss.Status.AvailableReplicas == ss.Status.Replicas
			logf.Log.Info("StatefulSet",
				"app", "Postgres",
				"ready", ready,
				"name", ss.Name,
				"availableReplicas", ss.Status.AvailableReplicas,
				"readyReplicas", ss.Status.ReadyReplicas,
				"currentReplicas", ss.Status.CurrentReplicas,
			)

		}
		for !ready && counter > 0 {
			pods, err := ListPod(psql.Namespace)
			if err != nil {
				return err
			}
			for _, pod := range pods.Items {
				if strings.Contains(pod.Name, psql.ReleaseName) {
					if pod.Status.Phase == coreV1.PodRunning {
						ready = true
						psql.Pod = pod
						logf.Log.Info("Pod",
							"app", "Postgres",
							"ready", ready,
							"name", pod.Name,
							"status", pod.Status.Phase,
						)
						break
					}
					logf.Log.Info("Pod",
						"app", "Postgres",
						"ready", ready,
						"name", pod.Name,
						"status", pod.Status.Phase,
					)
				}
			}
			if ready {
				var pvcName = psql.PvcName
				logf.Log.Info("Verify volume provision", "pvc name", pvcName, "namespace", psql.Namespace)
				uuid, err := VerifyVolumeProvision(pvcName, psql.Namespace)
				if err != nil {
					return fmt.Errorf("failed to verify volume provisioning")
				}
				psql.VolUuid = uuid
				logf.Log.Info("postgres standalone installation is ready")
				return nil
			} else {
				logf.Log.Info("not all apps are ready yet")
				time.Sleep(10 * time.Second)
				counter--
				continue
			}
		}
	}
	return nil
}

const (
	pgBenchImage               = "postgres:16"
	defaultConcurrentDbClients = 4
	defaultDurationSeconds     = 60
)

type PgBenchApp struct {
	BenchmarkParams PgBenchmarkParams
	Name            string
	Namespace       string
	NodeSelector    string
	Pod             *coreV1.Pod
}

type PgBenchmarkParams struct {
	ConcurrentClients int
	ThreadCount       int
	DurationSeconds   int
}

func NewPgBench() *PgBenchApp {
	return &PgBenchApp{
		BenchmarkParams: PgBenchmarkParams{
			ConcurrentClients: defaultConcurrentDbClients,
			ThreadCount:       defaultThreadCount,
			DurationSeconds:   defaultDurationSeconds,
		},
	}
}

func (pgBench *PgBenchApp) InitializePgBench(host string) error {
	jobName := fmt.Sprintf("%s-init", pgBench.Name)
	job := &batchV1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: pgBench.Namespace,
		},
		Spec: batchV1.JobSpec{
			Template: coreV1.PodTemplateSpec{
				Spec: coreV1.PodSpec{
					Containers: []coreV1.Container{
						{
							Name:  pgBench.Name,
							Image: pgBenchImage, // Ensure this version includes pgbench
							Command: []string{
								"pgbench",
								"-i", // Initialize the database
								"-h", host,
								"-p", fmt.Sprintf("%d", e2e_config.GetConfig().Product.PostgresDatabasePort),
								"-U", e2e_config.GetConfig().Product.PostgresAuthUsername,
								"-d", e2e_config.GetConfig().Product.PostgresDatabaseName,
							},
							Env: []coreV1.EnvVar{
								{
									Name:  "PGPASSWORD",
									Value: e2e_config.GetConfig().Product.PostgresAuthPassword,
								},
							},
						},
					},
					RestartPolicy: coreV1.RestartPolicyNever,
				},
			},
		},
	}

	jobsClient := gTestEnv.KubeInt.BatchV1().Jobs(pgBench.Namespace)
	fmt.Printf("Creating job %s...\n", jobName)
	_, err := jobsClient.Create(context.TODO(), job, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("error creating job %s: %v", jobName, err)
	}

	fmt.Printf("Job %s created successfully\n", jobName)
	return pgBench.waitForJobCompletion(jobName)
}

func (pgBench *PgBenchApp) RunPgBench(host string) error {
	jobName := fmt.Sprintf("%s-benchmark", pgBench.Name)
	job := &batchV1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: pgBench.Namespace,
		},
		Spec: batchV1.JobSpec{
			Template: coreV1.PodTemplateSpec{
				Spec: coreV1.PodSpec{
					Containers: []coreV1.Container{
						{
							Name:  pgBench.Name,
							Image: pgBenchImage, // Ensure this version includes pgbench
							Command: []string{
								"pgbench",
								"-h", host,
								"-p", fmt.Sprintf("%d", e2e_config.GetConfig().Product.PostgresDatabasePort),
								"-U", e2e_config.GetConfig().Product.PostgresAuthUsername,
								"-d", e2e_config.GetConfig().Product.PostgresDatabaseName,
								"-T", fmt.Sprintf("%d", pgBench.BenchmarkParams.DurationSeconds), // Run for specified time
								"-j", fmt.Sprintf("%d", pgBench.BenchmarkParams.ThreadCount), // Number of threads
								"-c", fmt.Sprintf("%d", pgBench.BenchmarkParams.ConcurrentClients),
							},
							Env: []coreV1.EnvVar{
								{
									Name:  "PGPASSWORD",
									Value: e2e_config.GetConfig().Product.PostgresAuthPassword,
								},
							},
						},
					},
					RestartPolicy: coreV1.RestartPolicyNever,
				},
			},
		},
	}

	jobsClient := gTestEnv.KubeInt.BatchV1().Jobs(pgBench.Namespace)
	fmt.Printf("Creating job %s...\n", jobName)
	_, err := jobsClient.Create(context.TODO(), job, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("error creating job %s: %v", jobName, err)
	}

	fmt.Printf("Job %s created successfully\n", jobName)
	return pgBench.waitForJobCompletion(jobName)
}

func (pgBench *PgBenchApp) waitForJobCompletion(jobName string) error {
	jobsClient := gTestEnv.KubeInt.BatchV1().Jobs(pgBench.Namespace)
	for {
		job, err := jobsClient.Get(context.TODO(), jobName, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("error getting job %s status: %v", jobName, err)
		}
		if job.Status.Succeeded > 0 {
			fmt.Printf("Job %s completed successfully!\n", jobName)
			break
		}
		fmt.Printf("Job %s is still running...\n", jobName)
		time.Sleep(10 * time.Second)
	}

	// Clean up the Job after completion
	deletePropagation := metav1.DeletePropagationBackground
	err := jobsClient.Delete(context.TODO(), jobName, metav1.DeleteOptions{
		PropagationPolicy: &deletePropagation,
	})
	if err != nil {
		return fmt.Errorf("error deleting job %s: %v", jobName, err)
	}
	fmt.Printf("Deleted job %s.\n", jobName)
	return nil
}

func SetupPostgresEnvironment() ([]coreV1.Node, error) {
	var unlabeledNodes []coreV1.Node
	nonMsNs, err := ListAllNonMsnNodes()
	if err != nil {
		return nil, err
	}
	if nonMsNs != nil && len(nonMsNs.Items) > 0 {
		for _, n := range nonMsNs.Items {
			err = UpdateNodeTaints(n.Name, "openebs-test-control:NoSchedule-")
			if err != nil {
				return unlabeledNodes, err
			}
			unlabeledNodes = append(unlabeledNodes, n)
		}
		return unlabeledNodes, nil
	}
	ns, err := ListIOEngineNodes()
	if err != nil {
		return nil, err
	}
	if ns != nil && len(ns.Items) < 4 {
		return nil, fmt.Errorf("not enough nodes to complete test needs at least %d found %d", 4, len(ns.Items))
	}
	count := len(ns.Items) - 3
	pools, err := ListMsPools()
	if err != nil {
		return nil, err
	}

	for i := 1; i <= count; i++ {
		node := ns.Items[len(ns.Items)-i]
		for _, pool := range pools {
			if strings.Contains(pool.Name, node.Name) {
				logf.Log.Info("deleting pool on unlabeled node", "pool", pool.Name, "node", node.Name)
				err = custom_resources.DeleteMsPool(pool.Name)
				if err != nil {
					return nil, err
				}
			}
		}
		logf.Log.Info("unlabeled node", "node", node.Name)
		err = UnlabelNode(node.Name, e2e_config.GetConfig().Product.EngineLabel)
		if err != nil {
			return nil, err
		}
		msg, err := ZeroNodeHugePages(node.Name)
		if err != nil {
			return nil, err
		}
		logf.Log.Info("message from e2e-agent", "msg", msg)
		unlabeledNodes = append(unlabeledNodes, node)
	}
	return unlabeledNodes, nil
}
