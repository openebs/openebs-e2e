package apps

import (
	"fmt"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/e2e_config"
	"github.com/openebs/openebs-e2e/common/k8stest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

type Architecture string

const (
	Standalone Architecture = "standalone"
	Replicaset Architecture = "replicaset"
)

func (a Architecture) String() string {
	return string(a)
}

type mongoBuilder struct {
	architecture            Architecture
	CloneFsIdAsVolumeIdType common.CloneFsIdAsVolumeIdType
	filesystemType          common.FileSystemType
	helmVersion             string
	namespace               string
	nodeSelector            string
	provisioningType        common.ProvisioningType
	pvcName                 string
	pvcSize                 int
	releaseName             string
	replicaCount            int
	scName                  string
	values                  map[string]interface{}
	ycsb                    bool
}

type MongoApp struct {
	Ycsb  k8stest.YcsbApp
	Mongo k8stest.MongoApp
}

// NewMongoBuilder creates new mongo builder with default standalone settings
func NewMongoBuilder() *mongoBuilder {
	def := make(map[string]interface{})
	def["auth.rootPassword"] = e2e_config.GetConfig().Product.MongoAuthRootPassword
	def["auth.username"] = e2e_config.GetConfig().Product.MongoAuthUsername
	def["auth.password"] = e2e_config.GetConfig().Product.MongoAuthPassword
	def["auth.database"] = e2e_config.GetConfig().Product.MongoAuthDatabase
	def["architecture"] = Standalone.String()
	def["replicaCount"] = 1

	return &mongoBuilder{
		architecture:            Standalone,
		CloneFsIdAsVolumeIdType: common.CloneFsIdAsVolumeIdNone,
		filesystemType:          common.Ext4FsType,
		namespace:               common.NSDefault,
		provisioningType:        common.ThinProvisioning,
		releaseName:             e2e_config.GetConfig().Product.MongoReleaseName,
		replicaCount:            0,
		values:                  def,
	}
}

// WithHaMode Enables HA mode
func (mb *mongoBuilder) WithHaMode() *mongoBuilder {
	mb.architecture = Replicaset
	mb.values["architecture"] = Replicaset.String()
	return mb
}

// WithReplicaCount available only for HA mode - default = 2
func (mb *mongoBuilder) WithReplicaCount(replicaCount int) *mongoBuilder {
	if arch, ok := mb.values["architecture"].(string); ok {
		if arch == Replicaset.String() {
			mb.values["replicaCount"] = replicaCount
			return mb
		}
		logf.Log.Info("cannot set mongo.replicaCount with standalone architecture")
	}
	return mb
}

// WithMongoVersion Could specify mongo image version - default = "mongo image version"
func (mb *mongoBuilder) WithMongoVersion(version string) *mongoBuilder {
	mb.values["image.tag"] = version
	return mb
}

// WithStatefulSet able to set up statefulSet mode for standalone architecture - default = false
func (mb *mongoBuilder) WithStatefulSet() *mongoBuilder {
	if arch, ok := mb.values["architecture"].(string); ok {
		if arch == Standalone.String() {
			mb.values["useStatefulSet"] = true
			n, err := k8stest.GetIOEngineHostNameLabel()
			if err != nil {
				panic(err)
			}
			mb.values["nodeSelector.kubernetes\\.io/hostname"] = n
			mb.nodeSelector = n
			return mb
		}
		logf.Log.Info("cannot set statefulSet with replicaset architecture")
	}
	return mb
}

// WithNamespace Could specify namespace default = "default"
func (mb *mongoBuilder) WithNamespace(namespace string) *mongoBuilder {
	mb.namespace = namespace
	return mb
}

// WithNodeSelector able set a node selector for deploying an application to a specific node
func (mb *mongoBuilder) WithNodeSelector(nodeName string) *mongoBuilder {
	if arch, ok := mb.values["architecture"].(string); ok {
		if arch == Standalone.String() {
			mb.values["nodeSelector.kubernetes\\.io/hostname"] = nodeName
			mb.nodeSelector = nodeName
			return mb
		}
		logf.Log.Info("cannot set statefulSet with replicaset architecture")
	}
	return mb
}

// WithPvcSize Could specify PVC size - default = 9Gi
func (mb *mongoBuilder) WithPvcSize(sizeGi int) *mongoBuilder {
	mb.pvcSize = sizeGi
	mb.values["persistence.size"] = fmt.Sprintf("%dGi", sizeGi)
	return mb
}

// WithReleaseName Could specify helm release name - default = "maya-mongo"
func (mb *mongoBuilder) WithReleaseName(releaseName string) *mongoBuilder {
	mb.releaseName = releaseName
	return mb
}

// WithHelmChartVersion Could specify helm chart version - default = "latest"
func (mb *mongoBuilder) WithHelmChartVersion(chartVersion string) *mongoBuilder {
	mb.helmVersion = chartVersion
	return mb
}

// WithFileSystemType Could specify volume provisioning type - default = "thin"
func (mb *mongoBuilder) WithFileSystemType(systemType common.FileSystemType) *mongoBuilder {
	mb.filesystemType = systemType
	return mb
}

// WithOwnStorageClass Could specify already created storage class name.
// If this option is selected, you must ensure that the storage class parameters are set correctly.
func (mb *mongoBuilder) WithOwnStorageClass(scName string) *mongoBuilder {
	mb.scName = scName
	mb.values["global.storageClass"] = scName
	return mb
}

// WithPvc Could specify pvc
// If this option is selected, you must ensure that the pvc is created correctly.
func (mb *mongoBuilder) WithPvc(pvcName string) *mongoBuilder {
	if arch, ok := mb.values["architecture"].(string); ok {
		if arch == Standalone.String() {
			mb.pvcName = pvcName
			mb.values["persistence.existingClaim"] = pvcName
			return mb
		}
	}
	return mb
}

// WithAnotherValuesParameters Could specify another parameters from https://artifacthub.io/packages/helm/bitnami/mongodb
func (mb *mongoBuilder) WithAnotherValuesParameters(values map[string]interface{}) *mongoBuilder {
	for k, v := range values {
		mb.values[k] = v
	}
	return mb
}

// WithYcsb Also install YCSB benchmark app
func (mb *mongoBuilder) WithYcsb() *mongoBuilder {
	mb.ycsb = true
	return mb
}

func (mb *mongoBuilder) Build() (MongoApp, error) {
	var latest Chart
	latest, err := GetLatestHelmChartVersion(e2e_config.GetConfig().Product.MongoHelmRepo)
	if err != nil && mb.helmVersion == "" {
		logf.Log.Error(err, "switching to default bitnami/mongo chart version", "defaultChart", e2e_config.GetConfig().Product.MongoDefaultChartVersion)
		mb.helmVersion = e2e_config.GetConfig().Product.MongoDefaultChartVersion
	}
	if mb.helmVersion == "" {
		mb.helmVersion = latest.Version
	}
	if mb.scName == "" {
		scName, err := CreateStorageClass(mb)
		if err != nil {
			return MongoApp{}, err
		}
		mb.values["global.storageClass"] = scName
		mb.scName = scName
		logf.Log.Info("StorageClass has been created", "storageClassName", scName)
	}
	err = AddHelmRepository(e2e_config.GetConfig().Product.MongoHelmRepoName, e2e_config.GetConfig().Product.MongoHelmRepoUrl)
	if err != nil {
		return MongoApp{}, err
	}
	err = InstallHelmChart(e2e_config.GetConfig().Product.MongoHelmRepo, mb.helmVersion, mb.namespace, mb.releaseName, mb.values)
	if err != nil {
		return MongoApp{}, err
	}

	ma := k8stest.MongoApp{
		Namespace:    mb.namespace,
		ReleaseName:  mb.releaseName,
		ReplicaCount: mb.values["replicaCount"].(int),
		ScName:       mb.values["global.storageClass"].(string),
		Standalone:   mb.values["architecture"].(string) == Standalone.String(),
		PvcName:      mb.pvcName,
	}
	err = ma.MongoInstallReady()
	if err != nil {
		return MongoApp{}, err
	}

	// deploy ycsb
	ycsb := k8stest.NewYCSB()
	if mb.ycsb {
		ycsb.Namespace = mb.namespace
		ycsb.NodeSelector = mb.nodeSelector
		var name, podName string
		if mb.pvcName != "" {
			name, podName, err = ycsb.DeployYcsbApp(mb.pvcName)
			if err != nil {
				return MongoApp{}, err
			}
		} else {
			name, podName, err = ycsb.DeployYcsbApp(mb.scName)
			if err != nil {
				return MongoApp{}, err
			}
		}
		ycsb.Name = name
		ycsb.PodName = podName
		ycsb.MongoConnUrl = fmt.Sprintf("mongodb.url=mongodb://%s:%s@%s-mongodb.%s.svc.cluster.local:%d/%s", e2e_config.GetConfig().Product.MongoAuthUsername, e2e_config.GetConfig().Product.MongoAuthPassword, mb.releaseName, mb.namespace, e2e_config.GetConfig().Product.MongoDatabasePort, e2e_config.GetConfig().Product.MongoAuthDatabase)
	}

	ma.Ycsb = mb.ycsb

	mongoApp := MongoApp{
		Ycsb:  *ycsb,
		Mongo: ma,
	}
	return mongoApp, nil
}

func (mb *mongoBuilder) Upgrade(app *MongoApp) (MongoApp, error) {
	err := UpgradeHelmChart(e2e_config.GetConfig().Product.MongoHelmRepo, mb.namespace, app.Mongo.ReleaseName, mb.values)
	if err != nil {
		return MongoApp{}, err
	}
	upgradedMongoApp := *app
	err = upgradedMongoApp.Mongo.MongoInstallReady()
	if err != nil {
		return MongoApp{}, err
	}
	ycsb := k8stest.NewYCSB()
	if mb.ycsb {
		ycsb.Namespace = mb.namespace
		ycsb.NodeSelector = mb.nodeSelector
		name, podName, err := ycsb.DeployYcsbApp(mb.pvcName)
		if err != nil {
			return MongoApp{}, err
		}
		ycsb.Name = name
		ycsb.PodName = podName
		ycsb.MongoConnUrl = fmt.Sprintf("mongodb.url=mongodb://%s:%s@%s-mongodb.%s.svc.cluster.local:%d/%s", e2e_config.GetConfig().Product.MongoAuthUsername, e2e_config.GetConfig().Product.MongoAuthPassword, mb.releaseName, mb.namespace, e2e_config.GetConfig().Product.MongoDatabasePort, e2e_config.GetConfig().Product.MongoAuthDatabase)
	}
	upgradedMongoApp.Ycsb = *ycsb
	upgradedMongoApp.Mongo.Namespace = mb.namespace
	upgradedMongoApp.Mongo.ReplicaCount = mb.values["replicaCount"].(int)
	upgradedMongoApp.Mongo.ScName = mb.values["global.storageClass"].(string)
	upgradedMongoApp.Mongo.PvcName = mb.pvcName
	return upgradedMongoApp, err
}
