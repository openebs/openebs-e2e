package k8stest

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/e2e_config"

	v1 "k8s.io/api/core/v1"

	logf "sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/pkg/errors"
	storagev1 "k8s.io/api/storage/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type StorageClass struct {
	object        *storagev1.StorageClass
	reclaimPolicy v1.PersistentVolumeReclaimPolicy
}

// ScBuilder enables building an instance of StorageClass
type ScBuilder struct {
	sc   *StorageClass
	errs []error
}

// NewScBuilder returns new instance of ScBuilder
func NewScBuilder() *ScBuilder {
	obj := ScBuilder{sc: &StorageClass{object: &storagev1.StorageClass{}}}
	// set default mayastor csi provisioner
	scObject := obj.WithProvisioner(e2e_config.GetConfig().Product.CsiProvisioner)

	// set default replicas value i.e 1
	scObject = scObject.WithReplicas(common.DefaultReplicaCount())

	if e2e_config.GetConfig().DefaultThinProvisioning {
		scObject = scObject.WithProvisioningType(common.ThinProvisioning)
	}

	return scObject
}

// WithName sets the Name field of storageclass with provided argument.
func (b *ScBuilder) WithName(name string) *ScBuilder {
	if len(name) == 0 {
		b.errs = append(b.errs, errors.New("failed to build storageclass: missing storageclass name"))
		return b
	}
	b.sc.object.Name = name
	return b
}

// WithNamespace sets the namespace field of storageclass with provided argument.
func (b *ScBuilder) WithNamespace(ns string) *ScBuilder {
	if len(ns) == 0 {
		b.errs = append(b.errs, errors.New("failed to build storageclass: missing storageclass namespace"))
		return b
	}
	b.sc.object.Namespace = ns
	return b
}

// WithGenerateName appends a random string after the name
func (b *ScBuilder) WithGenerateName(name string) *ScBuilder {
	b.sc.object.GenerateName = name + "-"
	return b
}

// WithAnnotations sets the Annotations field of storageclass with provided value.
func (b *ScBuilder) WithAnnotations(annotations map[string]string) *ScBuilder {
	if len(annotations) == 0 {
		b.errs = append(b.errs, errors.New("failed to build storageclass: missing annotations"))
	}
	b.sc.object.Annotations = annotations
	return b
}

// WithReplicas sets the replica parameter of storageclass with provided argument.
func (b *ScBuilder) WithReplicas(value int) *ScBuilder {
	if value == 0 {
		b.errs = append(b.errs, errors.New("failed to build storageclass: missing storageclass replicas"))
		return b
	}
	if b.sc.object.Parameters == nil {
		b.sc.object.Parameters = map[string]string{}
	}
	b.sc.object.Parameters[string(common.ScReplicas)] = strconv.Itoa(value)
	return b
}

func (b *ScBuilder) WithNodeHasTopologyKey(topologyKey string) *ScBuilder {
	if topologyKey == "" {
		b.errs = append(b.errs, errors.New("failed to build storageclass: missing topology key"))
		return b
	}
	if b.sc.object.Parameters == nil {
		b.sc.object.Parameters = map[string]string{}
	}
	// Assuming poolHasTopologyKey is a multi-line string in YAML format
	b.sc.object.Parameters[string(common.ScNodeHasTopologyKey)] = topologyKey
	return b
}

func (b *ScBuilder) WithNodeSpreadTopologyKey(topologyKey string) *ScBuilder {
	if topologyKey == "" {
		b.errs = append(b.errs, errors.New("failed to build storageclass: missing topology key"))
		return b
	}
	if b.sc.object.Parameters == nil {
		b.sc.object.Parameters = map[string]string{}
	}
	// Add the topology key to parameters
	b.sc.object.Parameters[string(common.ScNodeSpreadTopologyKey)] = topologyKey
	return b
}

func (b *ScBuilder) WithNodeAffinityTopologyLabel(topology map[string]string) *ScBuilder {
	if len(topology) == 0 {
		b.errs = append(b.errs, errors.New("failed to build storageclass: missing topology"))
		return b
	}
	if b.sc.object.Parameters == nil {
		b.sc.object.Parameters = map[string]string{}
	}
	// Convert the map to a single-line string in YAML format
	var yamlTopology strings.Builder
	for key, value := range topology {
		yamlTopology.WriteString(fmt.Sprintf("%s: %s\n", key, value))
	}
	b.sc.object.Parameters[string(common.ScNodeAffinityTopologyLabel)] = yamlTopology.String()
	return b
}

func (b *ScBuilder) WithPoolHasTopologyKey(topologyKey string) *ScBuilder {
	if topologyKey == "" {
		b.errs = append(b.errs, errors.New("failed to build storageclass: missing topology key"))
		return b
	}
	if b.sc.object.Parameters == nil {
		b.sc.object.Parameters = map[string]string{}
	}
	// Assuming poolHasTopologyKey is a multi-line string in YAML format
	b.sc.object.Parameters[string(common.ScPoolHasTopologyKey)] = topologyKey
	return b
}

func (b *ScBuilder) WithPoolAffinityTopologyLabel(topology map[string]string) *ScBuilder {
	if len(topology) == 0 {
		b.errs = append(b.errs, errors.New("failed to build storageclass: missing topology"))
		return b
	}
	if b.sc.object.Parameters == nil {
		b.sc.object.Parameters = map[string]string{}
	}
	// Convert the map to a single-line string in YAML format
	var yamlTopology strings.Builder
	for key, value := range topology {
		yamlTopology.WriteString(fmt.Sprintf("%s: %s\n", key, value))
	}
	b.sc.object.Parameters[string(common.ScPoolAffinityTopologyLabel)] = yamlTopology.String()
	return b
}

// WithFileType sets the fsType parameter of storageclass with provided argument.
func (b *ScBuilder) WithFileSystemType(value common.FileSystemType) *ScBuilder {
	if value != common.NoneFsType {
		if b.sc.object.Parameters == nil {
			b.sc.object.Parameters = map[string]string{}
		}
		b.sc.object.Parameters[string(common.ScFsType)] = string(value)
	}
	return b
}

// WithProtocol sets the protocol parameter of storageclass with provided argument.
func (b *ScBuilder) WithProtocol(value common.ShareProto) *ScBuilder {
	if b.sc.object.Parameters == nil {
		b.sc.object.Parameters = map[string]string{}
	}
	b.sc.object.Parameters[string(common.ScProtocol)] = string(value)
	return b
}

// WithCloneFsIdAsVolumeId sets the cloneFsIdAsVolumeId parameter of storageclass with provided argument.
func (b *ScBuilder) WithCloneFsIdAsVolumeId(value common.CloneFsIdAsVolumeIdType) *ScBuilder {
	if b.sc.object.Parameters == nil {
		b.sc.object.Parameters = map[string]string{}
	}
	if value == common.CloneFsIdAsVolumeIdEnable {
		b.sc.object.Parameters[common.ScCloneFsIdAsVolumeId] = "true"
	} else if value == common.CloneFsIdAsVolumeIdDisable {
		b.sc.object.Parameters[common.ScCloneFsIdAsVolumeId] = "false"
	} else if value == common.CloneFsIdAsVolumeIdNone {
		return b
	}
	return b
}

// WithIOTimeout sets the ioTimeout parameter of storageclass with provided argument.
func (b *ScBuilder) WithIOTimeout(value int) *ScBuilder {
	if b.sc.object.Parameters == nil {
		b.sc.object.Parameters = map[string]string{}
	}
	b.sc.object.Parameters[string(common.ScIOTimeout)] = strconv.Itoa(value)
	return b
}

// WithNvmeCtrlLossTmo sets the nvme controller loss timeout parameter of storageclass with provided argument.
func (b *ScBuilder) WithNvmeCtrlLossTmo(value int) *ScBuilder {
	if b.sc.object.Parameters == nil {
		b.sc.object.Parameters = map[string]string{}
	}
	b.sc.object.Parameters[string(common.ScNvmeCtrlLossTmo)] = strconv.Itoa(value)
	return b
}

// WithProvisioner sets the Provisioner field of storageclass with provided argument.
func (b *ScBuilder) WithProvisioner(provisioner string) *ScBuilder {
	if len(provisioner) == 0 {
		b.errs = append(b.errs, errors.New("failed to build storageclass: missing provisioner name"))
		return b
	}
	b.sc.object.Provisioner = provisioner
	return b
}

// WithReclaimPolicy sets the reclaim policy field of storageclass with provided argument.
func (b *ScBuilder) WithReclaimPolicy(reclaimPolicy v1.PersistentVolumeReclaimPolicy) *ScBuilder {
	b.sc.reclaimPolicy = reclaimPolicy
	b.sc.object.ReclaimPolicy = &b.sc.reclaimPolicy
	return b
}

// WithProvisioningType sets the thin provisioning field of storageclass with provided argument.
func (b *ScBuilder) WithProvisioningType(provisioningType common.ProvisioningType) *ScBuilder {
	if b.sc.object.Parameters == nil {
		b.sc.object.Parameters = map[string]string{}
	}
	if provisioningType == common.ThinProvisioning {
		b.sc.object.Parameters[common.ScThinProvisioning] = "true"
	}
	if provisioningType == common.ThickProvisioning {
		b.sc.object.Parameters[common.ScThinProvisioning] = "false"
	}
	return b
}

// WithStsAffinityGroup sets the stsAffinityGroup field of storageclass with provided argument.
func (b *ScBuilder) WithStsAffinityGroup(stsAffinity common.StsAffinityGroup) *ScBuilder {
	if b.sc.object.Parameters == nil {
		b.sc.object.Parameters = map[string]string{}
	}
	if stsAffinity == common.StsAffinityGroupEnable {
		b.sc.object.Parameters[common.ScStsAffinityGroup] = "true"
	}
	if stsAffinity == common.StsAffinityGroupDisable {
		b.sc.object.Parameters[common.ScStsAffinityGroup] = "false"
	}
	return b
}

// WithMaxSnapshots sets the maxSnapshots parameter of storageclass with provided argument.
func (b *ScBuilder) WithMaxSnapshots(value int) *ScBuilder {
	if b.sc.object.Parameters == nil {
		b.sc.object.Parameters = map[string]string{}
	}
	b.sc.object.Parameters[string(common.ScMaxSnapshots)] = strconv.Itoa(value)
	return b
}

// Build returns the StorageClass API instance
func (b *ScBuilder) Build() (*storagev1.StorageClass, error) {
	if len(b.errs) > 0 {
		return nil, errors.Errorf("%+v", b.errs)
	}
	logf.Log.Info("NewScBuilder:", "name", b.sc.object.Name,
		"parameters", b.sc.object.Parameters)
	return b.sc.object, nil
}

// WithVolumeBindingMode sets the volume binding mode of storageclass with
// provided argument.
// VolumeBindingMode indicates how PersistentVolumeClaims should be bound.
// VolumeBindingImmediate indicates that PersistentVolumeClaims should be
// immediately provisioned and bound. This is the default mode.
// VolumeBindingWaitForFirstConsumer indicates that PersistentVolumeClaims
// should not be provisioned and bound until the first Pod is created that
// references the PeristentVolumeClaim.  The volume provisioning and
// binding will occur during Pod scheduing.
func (b *ScBuilder) WithVolumeBindingMode(bindingMode storagev1.VolumeBindingMode) *ScBuilder {
	if bindingMode != "" {
		b.sc.object.VolumeBindingMode = &bindingMode
	}
	return b
}

func (b *ScBuilder) WithMountOptions(options []string) *ScBuilder {
	if len(options) == 0 {
		return b
	}
	if b.sc.object.MountOptions == nil {
		b.sc.object.MountOptions = []string{}
	}
	b.sc.object.MountOptions = append(b.sc.object.MountOptions, options...)
	return b
}

func (b *ScBuilder) WithMountOption(option string) *ScBuilder {
	return b.WithMountOptions([]string{option})
}

// WithVolumeExpansion sets the AllowVolumeExpansion field of storageclass.
func (b *ScBuilder) WithVolumeExpansion(value common.AllowVolumeExpansion) *ScBuilder {
	var volExpansion bool
	if value == common.AllowVolumeExpansionEnable {
		volExpansion = true
	} else if value == common.AllowVolumeExpansionDisable {
		volExpansion = false
	} else if value == common.AllowVolumeExpansionNone {
		return b
	}
	b.sc.object.AllowVolumeExpansion = &volExpansion
	return b
}

// Build and create the StorageClass
func (b *ScBuilder) BuildAndCreate() error {
	scObj, err := b.Build()
	if err == nil {
		err = CreateSc(scObj)
	}
	return err
}

// CreateSc creates storageclass with provided storageclass object
func CreateSc(obj *storagev1.StorageClass) error {
	thin, present := obj.Parameters[common.ScThinProvisioning]
	if !present {
		thin = "not set"
	}
	logf.Log.Info("Creating", "StorageClass", obj,
		"thin provisioning", thin)
	ScApi := gTestEnv.KubeInt.StorageV1().StorageClasses
	_, createErr := ScApi().Create(context.TODO(), obj, metaV1.CreateOptions{})
	return createErr
}
