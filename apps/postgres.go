package apps

import (
	"fmt"
	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/e2e_config"
	"github.com/openebs/openebs-e2e/common/k8stest"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

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
		replicaCount:            0,
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

// WithReleaseName Could specify helm release name - default = "maya-mongo"
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

func (pb *postgresBuilder) Build() (PostgresApp, error) {
	var latest Chart
	latest, err := GetLatestHelmChartVersion(e2e_config.GetConfig().Product.PostgresHelmRepo)
	if err != nil && pb.helmVersion == "" {
		logf.Log.Error(err, "switching to default bitnami/postgresql chart version", "defaultChart", e2e_config.GetConfig().Product.PostgresDefaultChartVersion)
		pb.helmVersion = e2e_config.GetConfig().Product.PostgresDefaultChartVersion
	}
	if pb.helmVersion == "" {
		pb.helmVersion = latest.Version
	}
	if pb.scName == "" {
		scName, err := CreatePostgresStorageClass(pb)
		if err != nil {
			return PostgresApp{}, err
		}
		pb.values["global.storageClass"] = scName
		pb.scName = scName
		logf.Log.Info("StorageClass has been created", "storageClassName", scName)
	}
	err = AddHelmRepository(e2e_config.GetConfig().Product.PostgresHelmRepoName, e2e_config.GetConfig().Product.PostgresHelmRepoUrl)
	if err != nil {
		return PostgresApp{}, err
	}
	err = InstallHelmChart(e2e_config.GetConfig().Product.PostgresHelmRepo, pb.helmVersion, pb.namespace, pb.releaseName, pb.values)
	if err != nil {
		return PostgresApp{}, err
	}

	pa := k8stest.PostgresApp{
		Namespace:    pb.namespace,
		ReleaseName:  pb.releaseName,
		ReplicaCount: pb.values["replicaCount"].(int),
		ScName:       pb.values["global.storageClass"].(string),
		Standalone:   pb.values["architecture"].(string) == Standalone.String(),
		PvcName:      pb.pvcName,
	}
	err = pa.PostgresInstallReady()
	if err != nil {
		return PostgresApp{}, err
	}

	pgBench := k8stest.NewPgBench()
	if pb.pgBench {
		pgBench.Name = "pgbench"
		pgBench.Namespace = pb.namespace
		pgBench.NodeSelector = pb.nodeSelector
	}
	postgresApp := PostgresApp{
		Postgres: pa,
		PgBench:  *pgBench,
	}
	return postgresApp, nil
}
