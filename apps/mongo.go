package apps

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/openebs/openebs-e2e/common/e2e_config"
	corev1 "k8s.io/api/core/v1"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/k8stest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	MongoPod            = "mongo"
	MongoPort           = 27017
	MongoReplicaSetName = "rs0"
	MongoMountPath      = "/data/db"
	MongoPVCName        = "mongo-data"
	MongoStsName        = "mongo"
	MongoSvcName        = "mongo"
	MongoHeadlessSvc    = "mongo-headless"
	MongoReadyTimeout   = 10 * time.Minute
)

type mongoBuilder struct {
	namespace      string
	replicaCount   int
	pvcSizeGi      int
	scName         string
	ycsb           bool
	releaseName    string
	useExistingPVC bool
	existingPVC    string
	nodeSelector   string
}

type MongoApp struct {
	Mongo k8stest.MongoApp
	Ycsb  k8stest.YcsbApp
}

func NewMongoBuilder() *mongoBuilder {
	return &mongoBuilder{
		namespace:    common.NSDefault,
		replicaCount: 1,
		pvcSizeGi:    5,
	}
}

func (mb *mongoBuilder) WithNamespace(ns string) *mongoBuilder {
	mb.namespace = ns
	return mb
}

func (mb *mongoBuilder) WithReplicaCount(n int) *mongoBuilder {
	mb.replicaCount = n
	return mb
}

func (mb *mongoBuilder) WithHaMode() *mongoBuilder {
	return mb
}

func (mb *mongoBuilder) WithPvcSize(sizeGi int) *mongoBuilder {
	mb.pvcSizeGi = sizeGi
	return mb
}

func (mb *mongoBuilder) WithOwnStorageClass(sc string) *mongoBuilder {
	mb.scName = sc
	return mb
}

func (mb *mongoBuilder) WithYcsb() *mongoBuilder {
	mb.ycsb = true
	return mb
}

func (mb *mongoBuilder) WithReleaseName(name string) *mongoBuilder {
	mb.releaseName = name
	return mb
}

func (mb *mongoBuilder) WithPvc(pvcName string) *mongoBuilder {
	if mb.replicaCount > 1 {
		logf.Log.Info("Ignoring WithPvc for HA Mongo", "pvc", pvcName)
		return mb
	}
	mb.useExistingPVC = true
	mb.existingPVC = pvcName
	return mb
}

func (mb *mongoBuilder) WithNodeSelector(nodeName string) *mongoBuilder {
	if mb.replicaCount > 1 {
		logf.Log.Info(
			"Ignoring WithNodeSelector: not supported for HA Mongo",
			"node", nodeName,
		)
		return mb
	}

	mb.nodeSelector = nodeName
	return mb
}

func (mb *mongoBuilder) Build() (MongoApp, error) {
	log := logf.Log.WithName("mongo-util")
	ns := mb.namespace

	release := MongoStsName
	if mb.releaseName != "" {
		release = mb.releaseName
	}

	if mb.scName == "" && !mb.useExistingPVC {
		sc, err := CreateStorageClass(mb)
		if err != nil {
			return MongoApp{}, err
		}
		mb.scName = sc
	}

	if err := applyMongoServices(ns); err != nil {
		return MongoApp{}, err
	}

	if err := createMongoStatefulSet(mb); err != nil {
		return MongoApp{}, err
	}

	if err := waitForMongoPods(ns, mb.replicaCount); err != nil {
		return MongoApp{}, err
	}

	mongoPodName := fmt.Sprintf("%s-0", MongoStsName)

	if mb.replicaCount > 1 {
		log.Info("Initializing MongoDB ReplicaSet")
		if err := initReplicaSet(ns, release, mb.replicaCount); err != nil {
			return MongoApp{}, err
		}
	}

	// Resolve Mayastor Volume UUID
	volUUID, err := getMongoVolumeUUID(mb)
	if err != nil {
		return MongoApp{}, err
	}

	var pvcName string

	if mb.useExistingPVC {
		pvcName = mb.existingPVC
	} else {
		pvcName = fmt.Sprintf("%s-%s-0", MongoPVCName, release)
	}

	mongo := k8stest.MongoApp{
		Namespace:    ns,
		ReleaseName:  release,
		ReplicaCount: mb.replicaCount,
		ScName:       mb.scName,
		Standalone:   mb.replicaCount == 1,
		StsName:      release,
		VolUuid:      volUUID,
		PvcName:      pvcName,
	}

	mongo.Pod.Name = mongoPodName
	mongo.Pod.Namespace = ns

	var ycsb k8stest.YcsbApp
	if mb.ycsb {
		y := k8stest.NewYCSB()
		y.Namespace = ns

		name, pod, err := y.DeployYcsbApp(mb.scName)
		if err != nil {
			return MongoApp{}, err
		}

		y.Name = name
		y.PodName = pod

		if mb.replicaCount > 1 {
			y.MongoConnUrl = fmt.Sprintf(
				"mongodb.url=mongodb://%s-0.%s.%s.svc.cluster.local:%d/?replicaSet=%s",
				release,
				MongoHeadlessSvc,
				ns,
				MongoPort,
				MongoReplicaSetName,
			)
		} else {
			y.MongoConnUrl = fmt.Sprintf(
				"mongodb.url=mongodb://%s.%s.svc.cluster.local:%d",
				MongoSvcName,
				ns,
				MongoPort,
			)
		}
		ycsb = *y
	}

	return MongoApp{Mongo: mongo, Ycsb: ycsb}, nil
}

func (mb *mongoBuilder) Upgrade(old *MongoApp) (MongoApp, error) {
	log := logf.Log.WithName("mongo-upgrade")
	ns := old.Mongo.Namespace

	if !mb.useExistingPVC {
		return MongoApp{}, fmt.Errorf("Upgrade requires WithPvc() to be set")
	}

	log.Info("Upgrading Mongo using existing PVC",
		"pvc", mb.existingPVC,
		"release", old.Mongo.ReleaseName,
	)

	newMongo := old.Mongo
	newMongo.PvcName = mb.existingPVC

	_ = k8stest.DeleteStatefulset(old.Mongo.StsName, ns)

	if err := createMongoStatefulSet(mb); err != nil {
		return MongoApp{}, err
	}

	if err := waitForMongoPods(ns, newMongo.ReplicaCount); err != nil {
		return MongoApp{}, err
	}

	volUUID, err := getMongoVolumeUUID(mb)
	if err != nil {
		return MongoApp{}, err
	}

	newMongo.VolUuid = volUUID
	newMongo.Pod.Name = fmt.Sprintf("%s-0", newMongo.StsName)
	newMongo.Pod.Namespace = ns

	return MongoApp{
		Mongo: newMongo,
		Ycsb:  old.Ycsb,
	}, nil
}

func applyMongoServices(namespace string) error {
	rootDir := e2e_config.GetConfig().OpenEbsE2eRootDir
	filePath := fmt.Sprintf("%s/configurations/service_mongo.yaml", rootDir)
	return k8stest.KubeCtlApplyYaml(filePath, rootDir)
}

func createMongoStatefulSet(mb *mongoBuilder) error {
	ns := mb.namespace
	replicas := int32(mb.replicaCount)

	command := []string{"mongod", "--bind_ip_all"}
	if mb.replicaCount > 1 {
		command = append(command, "--replSet", MongoReplicaSetName)
	}

	containerBuilder := k8stest.NewContainerBuilder().
		WithName(MongoPod).
		WithImage(e2e_config.GetConfig().Product.MongoImage).
		WithCommandNew(command).
		WithPortsNew([]corev1.ContainerPort{
			{ContainerPort: MongoPort},
		}).
		WithVolumeMountsNew([]corev1.VolumeMount{
			{
				Name:      MongoPVCName,
				MountPath: MongoMountPath,
			},
		})

	podTmplBuilder := k8stest.NewPodtemplatespecBuilder().
		WithLabelsNew(map[string]string{"app": MongoPod}).
		WithContainerBuildersNew(containerBuilder)

	if mb.useExistingPVC {
		volumeBuilder := k8stest.NewVolumeBuilder().
			WithName(MongoPVCName).
			WithPVCSource(mb.existingPVC)

		podTmplBuilder = podTmplBuilder.WithVolumeBuildersNew(volumeBuilder)
	}

	stsBuilder := k8stest.NewStatefulsetBuilder().
		WithName(MongoStsName).
		WithNamespace(ns).
		WithLabels(map[string]string{"app": MongoPod}).
		WithSelectorMatchLabels(map[string]string{"app": MongoPod}).
		WithReplicas(&replicas).
		WithPodTemplateSpecBuilder(podTmplBuilder)

	if !mb.useExistingPVC {
		stsBuilder = stsBuilder.WithVolumeClaimTemplate(
			MongoPVCName,
			fmt.Sprintf("%dGi", mb.pvcSizeGi),
			mb.scName,
			common.VolFileSystem,
		)
	}

	sts, err := stsBuilder.Build()
	if err != nil {
		return err
	}

	sts.Spec.ServiceName = MongoHeadlessSvc

	return k8stest.CreateStatefulset(sts)
}

func initReplicaSet(ns, name string, replicas int) error {
	members := []string{}
	for i := 0; i < replicas; i++ {
		members = append(members, fmt.Sprintf(
			`{ _id: %d, host: "%s-%d.%s.%s.svc.cluster.local:%d" }`,
			i, name, i, MongoHeadlessSvc, ns, MongoPort,
		))
	}

	cmd := fmt.Sprintf(
		`mongo --quiet --eval 'rs.initiate({_id: "%s", members: [%s]})'`,
		MongoReplicaSetName,
		strings.Join(members, ","),
	)

	_, _, err := k8stest.ExecuteCommandInPod(ns, fmt.Sprintf("%s-0", name), cmd)
	return err
}

func waitForMongoPods(ns string, replicas int) error {
	timeout := time.After(MongoReadyTimeout)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return fmt.Errorf("mongo pods not ready")
		case <-ticker.C:
			pods, err := k8stest.ListPodsWithLabel(ns, map[string]string{"app": "mongo"})
			if err != nil || len(pods.Items) != replicas {
				continue
			}
			sort.Slice(pods.Items, func(i, j int) bool {
				return pods.Items[i].Name < pods.Items[j].Name
			})
			allReady := true
			for _, p := range pods.Items {
				if p.Status.Phase != corev1.PodRunning {
					allReady = false
					break
				}
				for _, cs := range p.Status.ContainerStatuses {
					if !cs.Ready {
						allReady = false
						break
					}
				}
			}
			if allReady {
				return nil
			}
		}
	}
}

func getMongoVolumeUUID(mb *mongoBuilder) (string, error) {
	var pvcName string

	if mb.useExistingPVC {
		pvcName = mb.existingPVC
	} else {
		pvcName = fmt.Sprintf("%s-%s-0", MongoPVCName, MongoStsName)
	}

	pvc, err := k8stest.GetPVC(pvcName, mb.namespace)
	if err != nil {
		return "", fmt.Errorf("failed to get pvc %s: %w", pvcName, err)
	}

	if pvc.Spec.VolumeName == "" {
		return "", fmt.Errorf("pvc %s is not bound yet", pvcName)
	}

	pv, err := k8stest.GetPV(pvc.Spec.VolumeName)
	if err != nil {
		return "", fmt.Errorf("failed to get pv %s: %w", pvc.Spec.VolumeName, err)
	}

	if pv.Spec.CSI == nil || pv.Spec.CSI.VolumeHandle == "" {
		return "", fmt.Errorf("pv %s has no CSI volume handle", pv.Name)
	}

	return pv.Spec.CSI.VolumeHandle, nil
}
