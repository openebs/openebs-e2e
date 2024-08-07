package apps

import (
	"fmt"
	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/e2e_config"
	"github.com/openebs/openebs-e2e/common/k8stest"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

/*
Before you start, ensure that you have at least one node without huge pages and Mayastor installed. You can use the
SetupPostgresEnvironment function to prepare the environment, which will set up the necessary nodes for deploying the
PostgreSQL application. Use a node selector to ensure PostgreSQL is installed on the prepared node.

The steps involve creating a Mayastor storage class and volume, then using the Postgres builder to configure PostgreSQL
with this storage class and volume. Additionally, you can create a PgBench job to run load tests against PostgreSQL,
similar to how the fio pod is commonly used. If you want to run the PgBench job in the background or in parallel,
you can use a goroutine as shown in the example below.

// Step 1: Prepare the environment by setting up nodes - probably in BeforeSuite
unlabeledNodes, err := SetupPostgresEnvironment()
if err != nil {
    log.Fatalf("Failed to set up environment: %v", err)
}

// Step 2: Define storage class name and volume name

// Step 3: Build the PostgreSQL application configuration with PgBench
postgresApp, err := apps.NewPostgresBuilder().
    WithPgBench().
    WithOwnStorageClass(scName).
    WithPvc(volName).
    WithPrimaryNodeSelector(unlabeledNodes[0].Name).
    Build()
if err != nil {
    log.Fatalf("Failed to build PostgreSQL app: %v", err)
}

// Step 4: Initialize PgBench
err = postgresApp.PgBench.InitializePgBench(fmt.Sprintf("%s-postgresql", postgresApp.Postgres.ReleaseName))
if err != nil {
    log.Fatalf("Failed to initialize PgBench: %v", err)
}

// Step 5: Run PgBench in a goroutine to perform load testing
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    err := postgresApp.PgBench.RunPgBench(fmt.Sprintf("%s-postgresql", postgresApp.Postgres.ReleaseName))
    if err != nil {
        log.Fatalf("Failed to run PgBench: %v", err)
    }
}()

// Step 6: Wait for PgBench job to complete
wg.Wait()
logf.Log.Info("pgBench job completed")
*/

type postgresBuilder struct {
	architecture            Architecture
	CloneFsIdAsVolumeIdType common.CloneFsIdAsVolumeIdType
	filesystemType          common.FileSystemType
	helmVersion             string
	namespace               string
	nodeSelector            string
	pgBench                 bool
	provisioningType        common.ProvisioningType
	pvcName                 string
	pvcSize                 int
	releaseName             string
	replicaCount            int
	scName                  string
	values                  map[string]interface{}
}

type PostgresApp struct {
	Postgres k8stest.PostgresApp
	PgBench  k8stest.PgBenchApp
	builder  postgresBuilder
}

type Limits struct {
	Memory int `json:"memory"`
	Cpu    int `json:"cpu"`
}

// NewPostgresBuilder creates new postgres builder with default standalone settings
func NewPostgresBuilder() *postgresBuilder {
	def := make(map[string]interface{})
	def["global.postgresql.auth.postgresPassword"] = e2e_config.GetConfig().Product.PostgresAuthRootPassword
	def["global.postgresql.auth.username"] = e2e_config.GetConfig().Product.PostgresAuthUsername
	def["global.postgresql.auth.password"] = e2e_config.GetConfig().Product.PostgresAuthPassword
	def["global.postgresql.auth.database"] = e2e_config.GetConfig().Product.PostgresDatabaseName
	def["architecture"] = Standalone.String()
	def["replicaCount"] = 1

	return &postgresBuilder{
		architecture:            Standalone,
		CloneFsIdAsVolumeIdType: common.CloneFsIdAsVolumeIdNone,
		filesystemType:          common.Ext4FsType,
		namespace:               common.NSDefault,
		provisioningType:        common.ThinProvisioning,
		releaseName:             e2e_config.GetConfig().Product.PostgresReleaseName,
		values:                  def,
	}
}

// WithHaMode Enables HA mode
func (pb *postgresBuilder) WithHaMode() *postgresBuilder {
	pb.architecture = Replication
	pb.values["architecture"] = Replication.String()
	return pb
}

// WithReadReplicaCount available only for HA mode - default = 1
func (pb *postgresBuilder) WithReadReplicaCount(replicaCount int) *postgresBuilder {
	if arch, ok := pb.values["architecture"].(string); ok {
		if arch == Replication.String() {
			pb.values["readReplicas.replicaCount"] = replicaCount
			return pb
		}
		logf.Log.Info("cannot set readReplicas.replicaCount with standalone architecture")
	}
	return pb
}

// WithPostgresVersion Could specify postgres image version - default = "postgres image version"
func (pb *postgresBuilder) WithPostgresVersion(version string) *postgresBuilder {
	pb.values["image.tag"] = version
	return pb
}

// WithNamespace Could specify namespace default = "default"
func (pb *postgresBuilder) WithNamespace(namespace string) *postgresBuilder {
	pb.namespace = namespace
	return pb
}

// WithPrimaryNodeSelector able set a node selector for deploying an application to a specific node
func (pb *postgresBuilder) WithPrimaryNodeSelector(nodeName string) *postgresBuilder {
	pb.values["primary.nodeSelector.kubernetes\\.io/hostname"] = nodeName
	pb.nodeSelector = nodeName
	return pb
}

// WithNodeSelector able set a node selector for deploying an application to a specific node
func (pb *postgresBuilder) WithReplicaNodeSelector(nodeName string) *postgresBuilder {
	if arch, ok := pb.values["architecture"].(string); ok {
		if arch == Replication.String() {
			pb.values["readReplicas.nodeSelector"] = nodeName
			return pb
		}
		logf.Log.Info("cannot set statefulSet with replicaset architecture")
	}
	return pb
}

// WithPvcSize Could specify PVC size - default = 9Gi
func (pb *postgresBuilder) WithPvcSize(sizeGi int) *postgresBuilder {
	pb.pvcSize = sizeGi
	pb.values["persistence.size"] = fmt.Sprintf("%dGi", sizeGi)
	return pb
}

// WithReleaseName Could specify helm release name - default = "ms-postgres"
func (pb *postgresBuilder) WithReleaseName(releaseName string) *postgresBuilder {
	pb.releaseName = releaseName
	return pb
}

// WithHelmChartVersion Could specify helm chart version - default = "latest"
func (pb *postgresBuilder) WithHelmChartVersion(chartVersion string) *postgresBuilder {
	pb.helmVersion = chartVersion
	return pb
}

// WithFileSystemType Could specify volume provisioning type - default = "thin"
func (pb *postgresBuilder) WithFileSystemType(systemType common.FileSystemType) *postgresBuilder {
	pb.filesystemType = systemType
	return pb
}

// WithOwnStorageClass Could specify already created storage class name.
// If this option is selected, you must ensure that the storage class parameters are set correctly.
func (pb *postgresBuilder) WithOwnStorageClass(scName string) *postgresBuilder {
	pb.scName = scName
	pb.values["global.storageClass"] = scName
	return pb
}

// WithPvc Could specify pvc
// If this option is selected, you must ensure that the pvc is created correctly.
func (pb *postgresBuilder) WithPvc(pvcName string) *postgresBuilder {
	if arch, ok := pb.values["architecture"].(string); ok {
		if arch == Standalone.String() {
			pb.pvcName = pvcName
			pb.values["primary.persistence.existingClaim"] = pvcName
			return pb
		}
	}
	return pb
}

// WithAnotherValuesParameters Could specify another parameters from https://artifacthub.io/packages/helm/bitnami/postgresql
func (pb *postgresBuilder) WithAnotherValuesParameters(values map[string]interface{}) *postgresBuilder {
	for k, v := range values {
		pb.values[k] = v
	}
	return pb
}

// WithPgBench Also install pg_bench benchmark app
func (pb *postgresBuilder) WithPgBench() *postgresBuilder {
	pb.pgBench = true
	return pb
}

func (pb *postgresBuilder) Create() (PostgresApp, error) {
	var latest Chart
	var err error

	if pb.helmVersion == "" {
		latest, err = GetLatestHelmChartVersion(e2e_config.GetConfig().Product.PostgresHelmRepo)
		if err != nil {
			logf.Log.Error(err, "switching to default bitnami/postgresql chart version", "defaultChart", e2e_config.GetConfig().Product.PostgresDefaultChartVersion)
			pb.helmVersion = e2e_config.GetConfig().Product.PostgresDefaultChartVersion
		} else {
			pb.helmVersion = latest.Version
		}
	}

	if pb.scName == "" {
		scName, err := CreatePostgresStorageClass(pb)
		if err != nil {
			pbEmpty := postgresBuilder{}
			return &pbEmpty, err
		}
		pb.values["global.storageClass"] = scName
		pb.scName = scName
		logf.Log.Info("StorageClass has been created", "storageClassName", scName)
	}

	pgBench := &k8stest.PgBenchApp{}
	if pb.pgBench {
		pgBench = k8stest.NewPgBenchAppBuilder().
			WithName("pgbench").
			WithNamespace(pb.namespace).
			WithNodeSelector(pb.nodeSelector).
			Build()
	}

	

	postgresApp := PostgresApp{
		Postgres: k8stest.PostgresApp{
			Namespace:    pb.namespace,
			ReleaseName:  pb.releaseName,
			ReplicaCount: pb.values["replicaCount"].(int),
			ScName:       pb.values["global.storageClass"].(string),
			Standalone:   pb.values["architecture"].(string) == Standalone.String(),
			PvcName:      pb.pvcName,
		},
		PgBench: *pgBench,
		Builder: pb,
	}

	return postgresApp, nil
}

func (pa *PostgresApp) Install() error {
	err := AddHelmRepository(e2e_config.GetConfig().Product.PostgresHelmRepoName, e2e_config.GetConfig().Product.PostgresHelmRepoUrl)
	if err != nil {
		return err
	}

	err = InstallHelmChart(e2e_config.GetConfig().Product.PostgresHelmRepo, pa.builder.helmVersion, pa.builder.namespace, pa.builder.releaseName, pa.builder.values)
	if err != nil {
		return err
	}

	pa := k8stest.PostgresApp{
		Namespace:    pa.builder.namespace,
		ReleaseName:  pa.builder.releaseName,
		ReplicaCount: pa.builder.values["replicaCount"].(int),
		ScName:       pa.builder.values["global.storageClass"].(string),
		Standalone:   pa.builder.values["architecture"].(string) == Standalone.String(),
		PvcName:      pa.builder.pvcName,
	}

	err = pa.PostgresInstallReady()
	if err != nil {
		return err
	}

	return nil
}
