package k8stest

import (
	"context"
	"fmt"
	"github.com/openebs/openebs-e2e/common/custom_resources"
	"github.com/openebs/openebs-e2e/common/e2e_config"
	"k8s.io/client-go/kubernetes"
	"strings"
	"time"

	batchV1 "k8s.io/api/batch/v1"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// PostgresApp represents a PostgreSQL application deployment configuration.
type PostgresApp struct {
	Namespace    string     // Kubernetes namespace where the PostgreSQL application is deployed.
	Pod          coreV1.Pod // Kubernetes Pod running the PostgreSQL application.
	VolUuid      string     // UUID of the provisioned volume for the PostgreSQL application.
	ReleaseName  string     // Helm release name for the PostgreSQL application.
	ReplicaCount int        // Number of replicas for the PostgreSQL StatefulSet.
	ScName       string     // StorageClass name used for provisioning volumes.
	Standalone   bool       // Indicates if the PostgreSQL deployment is standalone.
	PvcName      string     // Name of the PersistentVolumeClaim.
	PgBench      PgBenchApp // PgBench application configuration for benchmarking.
}

// PostgresInstallReady checks if the PostgreSQL application is installed and ready.
func (psql *PostgresApp) PostgresInstallReady() error {
	logf.Log.Info("checking postgres application to be installed")
	ready := false
	counter := 12

	if psql.Standalone {
		// List all StatefulSets with label "app.kubernetes.io/name=postgresql" in the namespace.
		statefulSets, err := gTestEnv.KubeInt.AppsV1().StatefulSets(psql.Namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: fmt.Sprintf("%s=%s", e2e_config.GetConfig().Product.PostgresK8sLabelName, e2e_config.GetConfig().Product.PostgresK8sLabelValue)})
		if err != nil {
			return err
		}
		if len(statefulSets.Items) != 1 {
			return fmt.Errorf("there should be 1 stateful set for postgres deployment")
		}

		for _, ss := range statefulSets.Items {
			ready = ss.Status.ReadyReplicas == ss.Status.Replicas
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

			ready = true // Assume all pods are ready until proven otherwise
			for _, pod := range pods.Items {
				if strings.Contains(pod.Name, psql.ReleaseName) {
					if pod.Status.Phase != coreV1.PodRunning {
						ready = false
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
				logf.Log.Info("all pods are running")
				uuid, err := VerifyVolumeProvision(psql.PvcName, psql.Namespace)
				if err != nil {
					return fmt.Errorf("failed to verify volume provisioning: %v", err)
				}
				psql.VolUuid = uuid
				logf.Log.Info("postgres standalone installation is ready")
				return nil
			} else {
				logf.Log.Info("not all apps are ready yet")
				time.Sleep(10 * time.Second)
				counter--
			}
		}
	} else {
		return fmt.Errorf("replicaset architecture install check is not implemented")
	}
	return nil
}

// Default benchmark parameters
const (
	defaultConcurrentDbClients = 4
	defaultDurationSeconds     = 60
)

// PgBenchApp represents a PgBench application configuration for benchmarking PostgreSQL.
type PgBenchApp struct {
	BenchmarkParams PgBenchmarkParams // Parameters for benchmarking.
	Name            string            // Name of the PgBench application.
	Namespace       string            // Kubernetes namespace for the PgBench application.
	NodeSelector    string            // Node selector for scheduling the PgBench Pod.
	Pod             *coreV1.Pod       // Kubernetes Pod running the PgBench application.
	clientset       *kubernetes.Clientset
}

type PgBenchmarkParams struct {
	ConcurrentClients int // Number of concurrent clients for PgBench.
	ThreadCount       int // Number of threads for PgBench.
	DurationSeconds   int // Duration of the benchmark in seconds.
}

type PgBenchAppBuilder struct {
	app *PgBenchApp
}

// NewPgBenchAppBuilder creates a new PgBenchApp instance with default benchmark parameters.
func NewPgBenchAppBuilder() *PgBenchAppBuilder {
	return &PgBenchAppBuilder{
		app: &PgBenchApp{
			BenchmarkParams: PgBenchmarkParams{
				ConcurrentClients: defaultConcurrentDbClients,
				ThreadCount:       defaultThreadCount,
				DurationSeconds:   defaultDurationSeconds,
			},
		},
	}
}

func (b *PgBenchAppBuilder) WithName(name string) *PgBenchAppBuilder {
	b.app.Name = name
	return b
}

func (b *PgBenchAppBuilder) WithNamespace(namespace string) *PgBenchAppBuilder {
	b.app.Namespace = namespace
	return b
}

func (b *PgBenchAppBuilder) WithNodeSelector(nodeSelector string) *PgBenchAppBuilder {
	b.app.NodeSelector = nodeSelector
	return b
}

func (b *PgBenchAppBuilder) WithBenchmarkParams(params PgBenchmarkParams) *PgBenchAppBuilder {
	b.app.BenchmarkParams = params
	return b
}

func (b *PgBenchAppBuilder) Build() *PgBenchApp {
	return b.app
}

// InitializePgBench initializes the PgBench database. Must be called before RunPgBench
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
							Image: e2e_config.GetConfig().Product.PgBenchImage, // Ensure this version includes pgbench
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

// RunPgBench runs the PgBench benchmark on the PostgreSQL database.
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
							Image: e2e_config.GetConfig().Product.PgBenchImage, // Ensure this version includes pgbench
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

// waitForJobCompletion waits for the Kubernetes Job to complete.
func (pgBench *PgBenchApp) waitForJobCompletion(jobName string) error {
	for {
		job, err := pgBench.GetPgBenchJob(jobName)
		if err != nil {
			return err
		}
		if job.Status.Succeeded > 0 {
			fmt.Printf("Job %s completed successfully!\n", jobName)
			break
		}
		fmt.Printf("Job %s is still running...\n", jobName)
		time.Sleep(10 * time.Second)
	}

	// Clean up the Job after completion
	err := pgBench.DeletePgBenchJob(jobName)
	if err != nil {
		return err
	}
	return nil
}

func (pgBench *PgBenchApp) GetPgBenchJob(jobName string) (*batchV1.Job, error) {
	jobsClient := gTestEnv.KubeInt.BatchV1().Jobs(pgBench.Namespace)
	job, err := jobsClient.Get(context.TODO(), jobName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("error getting job %s status: %v", jobName, err)
	}
	return job, nil
}

func (pgBench *PgBenchApp) DeletePgBenchJob(jobName string) error {
	jobsClient := gTestEnv.KubeInt.BatchV1().Jobs(pgBench.Namespace)
	deletePropagation := metav1.DeletePropagationBackground
	err := jobsClient.Delete(context.TODO(), pgBench.Name, metav1.DeleteOptions{
		PropagationPolicy: &deletePropagation,
	})
	if err != nil {
		return fmt.Errorf("error deleting job %s: %v", jobName, err)
	}
	fmt.Printf("Deleted job %s.\n", jobName)
	return nil
}

// SetupPostgresEnvironment sets up the environment for PostgreSQL by managing node labels and taints. Normally called in Before suit action. It will return a slice of nodes which are ready for postgres installation.
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
