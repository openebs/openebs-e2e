package v1beta3

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/openebs/openebs-e2e/common"

	logf "sigs.k8s.io/controller-runtime/pkg/log"

	v1beta3 "github.com/openebs/openebs-e2e/common/custom_resources/api/types/v1beta3"
	v1beta3Client "github.com/openebs/openebs-e2e/common/custom_resources/clientset/v1beta3"
	crtypes "github.com/openebs/openebs-e2e/common/custom_resources/types"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

var poolClientSet *v1beta3Client.DiskPoolV1Beta3Client

type v1beta3Ifc struct {
}

type v1beta3DSP struct {
	v1beta3 *v1beta3.DiskPool
}

func Initialise(config *rest.Config, extended bool) (crtypes.DiskPoolFunctions, error) {
	var err error
	_ = v1beta3.PoolAddToScheme(scheme.Scheme)
	poolClientSet, err = v1beta3Client.DspNewForConfig(config)
	return v1beta3Ifc{}, err
}

// DiskPool implementation

func (p v1beta3DSP) String() string {
	return fmt.Sprintf("%v", p.v1beta3)
}

func (p v1beta3DSP) GetType() reflect.Type {
	return reflect.TypeOf(p)
}

func (p v1beta3DSP) GetName() string {
	if p.v1beta3 != nil {
		return p.v1beta3.GetName()
	}
	// panic?
	return ""
}

// For v1beta3 map Status.State to pool status
func (p v1beta3DSP) GetPoolStatus() string {
	if p.v1beta3 != nil {
		return p.v1beta3.Status.PoolStatus
	}
	return ""
}

// For v1beta3 map Status.State to CR status
func (p v1beta3DSP) GetCRStatus() string {
	if p.v1beta3 != nil {
		return p.v1beta3.Status.CRStatus
	}
	return ""
}

func (p v1beta3DSP) GetPoolReadyStatus() string {
	if p.v1beta3 != nil {
		return p.v1beta3.Status.Conditions[0].Status
	}
	return ""
}

func (p v1beta3DSP) GetPoolReadyReason() string {
	if p.v1beta3 != nil {
		return p.v1beta3.Status.Conditions[0].Reason
	}
	return ""
}

func (p v1beta3DSP) GetPoolErrorCount() uint64 {
	if p.v1beta3 != nil {
		return p.v1beta3.Status.Diag.Errors.IOErrors
	}
	return 0
}

func (p v1beta3DSP) GetPoolAlertStatus() string {
	if p.v1beta3 != nil {
		return p.v1beta3.Status.Diag.Errors.Status
	}
	return ""
}

func (p v1beta3DSP) GetStatusCapacity() uint64 {
	if p.v1beta3 != nil {
		return p.v1beta3.Status.Capacity
	}
	return 0
}

func (p v1beta3DSP) GetStatusUsed() uint64 {
	if p.v1beta3 != nil {
		return p.v1beta3.Status.Used
	}
	return 0
}

func (p v1beta3DSP) CompareStatus(otherP *crtypes.DiskPool) bool {
	other := *otherP
	if p.GetType() != other.GetType() {
		logf.Log.Info("comparison between different diskpool types is unsupported",
			"this", p.GetType(), "other", other.GetType())
		return false
	}
	dsp := other.(v1beta3DSP)
	if p.v1beta3 != nil && dsp.v1beta3 != nil {
		return reflect.DeepEqual(p.v1beta3.Status, dsp.v1beta3.Status)
	}
	return false
}

func (p v1beta3DSP) GetFinalizers() []string {
	if p.v1beta3 != nil {
		return p.v1beta3.Finalizers
	}
	// panic?
	return []string{}
}

func (p v1beta3DSP) SetFinalizers(finalizers []string) (crtypes.DiskPool, error) {
	var err error
	dsp := p
	if p.v1beta3 != nil {
		mspIn := *p.v1beta3
		mspIn.SetFinalizers(finalizers)
		dsp.v1beta3, err = poolClientSet.DiskPools().Update(context.TODO(), &mspIn, metaV1.UpdateOptions{})
		return dsp, err
	}
	return dsp, fmt.Errorf("uninitialised DiskPool")
}

func (p v1beta3DSP) GetSpecDisks() []string {
	if p.v1beta3 != nil {
		return p.v1beta3.Spec.Disks
	}
	return []string{}
}

func (p v1beta3DSP) SetSpecDisks(disks []string) (crtypes.DiskPool, error) {
	var err error
	dsp := p
	if p.v1beta3 != nil {
		mspIn := *p.v1beta3
		mspIn.Spec.Disks = disks
		dsp.v1beta3, err = poolClientSet.DiskPools().Update(context.TODO(), &mspIn, metaV1.UpdateOptions{})
		return dsp, err
	}
	return dsp, fmt.Errorf("uninitialised DiskPool")
}

func (p v1beta3DSP) SetSpecNode(node string) (crtypes.DiskPool, error) {
	var err error
	dsp := p
	if p.v1beta3 != nil {
		mspIn := *p.v1beta3
		mspIn.Spec.Node = node
		dsp.v1beta3, err = poolClientSet.DiskPools().Update(context.TODO(), &mspIn, metaV1.UpdateOptions{})
		return dsp, err
	}
	return dsp, fmt.Errorf("uninitialised DiskPool")
}

func (p v1beta3DSP) GetSpecNode() string {
	if p.v1beta3 != nil {
		return p.v1beta3.Spec.Node
	}
	return ""
}

func (p v1beta3DSP) SetSpecEncryptionSecret(secretName string) (crtypes.DiskPool, error) {
	var err error
	dsp := p
	if p.v1beta3 != nil {
		mspIn := *p.v1beta3
		mspIn.Spec.EncryptionConfig.Source.Secret.Name = secretName
		dsp.v1beta3, err = poolClientSet.DiskPools().Update(context.TODO(), &mspIn, metaV1.UpdateOptions{})
		return dsp, err
	}
	return dsp, fmt.Errorf("uninitialised DiskPool")
}

func (p v1beta3DSP) GetSpecEncryptionSecret() string {
	if p.v1beta3 != nil {
		return p.v1beta3.Spec.EncryptionConfig.Source.Secret.Name
	}
	return ""
}

func (p v1beta3DSP) IsPoolEncrypted() bool {
	if p.v1beta3 != nil {
		return p.v1beta3.Status.Encrypted
	}
	return false
}

func (p v1beta3DSP) GetAnnotations() map[string]string {
	if p.v1beta3 != nil {
		return p.v1beta3.Annotations
	}
	return map[string]string{}
}

func (p v1beta3DSP) SetAnnotations(annotations map[string]string) (crtypes.DiskPool, error) {
	var err error
	dsp := p
	if p.v1beta3 != nil {
		mspIn := *p.v1beta3
		mspIn.Annotations = annotations
		dsp.v1beta3, err = poolClientSet.DiskPools().Update(context.TODO(), &mspIn, metaV1.UpdateOptions{})
		return dsp, err
	}
	return dsp, fmt.Errorf("uninitialised DiskPool")
}

func (p v1beta3DSP) SetClusterSize(clusterSize string) (crtypes.DiskPool, error) {
	var err error
	dsp := p
	if p.v1beta3 != nil {
		mspIn := *p.v1beta3
		mspIn.Spec.ClusterSize = clusterSize
		dsp.v1beta3, err = poolClientSet.DiskPools().Update(context.TODO(), &mspIn, metaV1.UpdateOptions{})
		return dsp, err
	}
	return dsp, fmt.Errorf("uninitialised DiskPool")
}

func (p v1beta3DSP) GetClusterSize() string {
	if p.v1beta3 != nil {
		return p.v1beta3.Status.ClusterSize
	}
	return ""
}

// Expose status.maxExpandableSize if available
func (p v1beta3DSP) GetStatusMaxExpandableSize() string {
	if p.v1beta3 != nil {
		return p.v1beta3.Status.MaxExpandableSize
	}
	return ""
}

func (p v1beta3DSP) GetMaxExpansion() string {
	if p.v1beta3 != nil {
		return p.v1beta3.Spec.MaxExpansion
	}
	return ""
}

func (p v1beta3DSP) SetMaxExpansion(maxExpansion string) (crtypes.DiskPool, error) {
	var err error
	dsp := p
	if p.v1beta3 != nil {
		mspIn := *p.v1beta3
		mspIn.Spec.MaxExpansion = maxExpansion
		dsp.v1beta3, err = poolClientSet.DiskPools().Update(context.TODO(), &mspIn, metaV1.UpdateOptions{})
		return dsp, err
	}
	return dsp, fmt.Errorf("uninitialised DiskPool")
}

//  DiskPoolFunctions implementation

func (ifc v1beta3Ifc) CreateMsPool(poolName string, node string, disks []string) (crtypes.DiskPool, error) {
	msp := v1beta3.DiskPool{
		TypeMeta: metaV1.TypeMeta{Kind: "DiskPool"},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      poolName,
			Namespace: common.NSMayastor(),
		},
		Spec: v1beta3.DiskPoolSpec{
			Node:  node,
			Disks: disks,
		},
	}
	mspOut, err := poolClientSet.DiskPools().Create(context.TODO(), &msp, metaV1.CreateOptions{})
	dsp := v1beta3DSP{mspOut}
	return dsp, err
}

func (ifc v1beta3Ifc) CreateMsPoolWithTopologySpec(poolName string, node string, disks []string, labels map[string]string) (crtypes.DiskPool, error) {

	topology := &v1beta3.Topology{
		Labelled: labels,
	}

	msp := v1beta3.DiskPool{
		TypeMeta: metaV1.TypeMeta{Kind: "DiskPool"},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      poolName,
			Namespace: common.NSMayastor(),
		},
		Spec: v1beta3.DiskPoolSpec{
			Node:     node,
			Disks:    disks,
			Topology: topology,
		},
	}

	mspOut, err := poolClientSet.DiskPools().Create(context.TODO(), &msp, metaV1.CreateOptions{})
	dsp := v1beta3DSP{mspOut}
	return dsp, err
}

func (ifc v1beta3Ifc) CreateMsPoolWithClusterSize(poolName string, node string, disks []string, clusterSize string) (crtypes.DiskPool, error) {
	msp := v1beta3.DiskPool{
		TypeMeta: metaV1.TypeMeta{Kind: "DiskPool"},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      poolName,
			Namespace: common.NSMayastor(),
		},
		Spec: v1beta3.DiskPoolSpec{
			Node:        node,
			Disks:       disks,
			ClusterSize: clusterSize,
		},
	}

	mspOut, err := poolClientSet.DiskPools().Create(context.TODO(), &msp, metaV1.CreateOptions{})
	dsp := v1beta3DSP{mspOut}
	return dsp, err
}

// CreateMsPoolWithMaxSizeAndClusterSize creates a DiskPool CR with both MaxExpansion and ClusterSize set.
func (ifc v1beta3Ifc) CreateMsPoolWithMaxSizeAndClusterSize(poolName string, node string, disks []string, maxSize string, clusterSize string) (crtypes.DiskPool, error) {
	// Default cluster size to 4MiB if not provided
	if clusterSize == "" {
		clusterSize = "4 MiB"
	}

	msp := v1beta3.DiskPool{
		TypeMeta: metaV1.TypeMeta{Kind: "DiskPool"},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      poolName,
			Namespace: common.NSMayastor(),
		},
		Spec: v1beta3.DiskPoolSpec{
			Node:         node,
			Disks:        disks,
			MaxExpansion: maxSize,
			ClusterSize:  clusterSize,
		},
	}
	mspOut, err := poolClientSet.DiskPools().Create(context.TODO(), &msp, metaV1.CreateOptions{})
	dsp := v1beta3DSP{mspOut}
	return dsp, err
}

// CreateMsPoolWithEncryptionAndMaxSize creates a DiskPool CR with encryption, MaxExpansion and ClusterSize set.
func (ifc v1beta3Ifc) CreateMsPoolWithEncryptionAndMaxSize(poolName string, node string, disks []string, encryptionSecretName string, maxSize string, clusterSize string) (crtypes.DiskPool, error) {
	// Default cluster size to 4MiB if not provided
	if clusterSize == "" {
		clusterSize = "4 MiB"
	}

	logf.Log.Info("Creating DiskPool with encryption, max size and cluster size", "poolName", poolName, "node", node, "disks", disks, "encryptionSecretName", encryptionSecretName, "maxSize", maxSize, "clusterSize", clusterSize)
	msp := v1beta3.DiskPool{
		TypeMeta: metaV1.TypeMeta{Kind: "DiskPool"},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      poolName,
			Namespace: common.NSMayastor(),
		},
		Spec: v1beta3.DiskPoolSpec{
			Node:         node,
			Disks:        disks,
			MaxExpansion: maxSize,
			ClusterSize:  clusterSize,
			EncryptionConfig: &v1beta3.EncryptionConfig{
				Source: v1beta3.Source{
					Secret: v1beta3.Secret{
						Name: encryptionSecretName,
					},
				},
			},
		},
	}
	mspOut, err := poolClientSet.DiskPools().Create(context.TODO(), &msp, metaV1.CreateOptions{})
	dsp := v1beta3DSP{mspOut}
	return dsp, err
}

// CreateMsPoolWithMaxSize creates a DiskPool CR with MaxExpansion set.
func (ifc v1beta3Ifc) CreateMsPoolWithMaxSize(poolName string, node string, disks []string, maxSize string) (crtypes.DiskPool, error) {
	msp := v1beta3.DiskPool{
		TypeMeta: metaV1.TypeMeta{Kind: "DiskPool"},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      poolName,
			Namespace: common.NSMayastor(),
		},
		Spec: v1beta3.DiskPoolSpec{
			Node:         node,
			Disks:        disks,
			MaxExpansion: maxSize,
		},
	}
	mspOut, err := poolClientSet.DiskPools().Create(context.TODO(), &msp, metaV1.CreateOptions{})
	dsp := v1beta3DSP{mspOut}
	return dsp, err
}

func (ifc v1beta3Ifc) GetMsPool(poolName string) (crtypes.DiskPool, error) {
	msp := v1beta3.DiskPool{}
	res, err := poolClientSet.DiskPools().Get(context.TODO(), poolName, metaV1.GetOptions{})
	if res != nil && err == nil {
		msp = *res
	}
	return v1beta3DSP{v1beta3: &msp}, err
}

func (ifc v1beta3Ifc) DeleteMsPool(poolName string) error {
	err := poolClientSet.DiskPools().Delete(context.TODO(), poolName, metaV1.DeleteOptions{})
	return err
}

func (ifc v1beta3Ifc) ListMsPoolCrs() ([]crtypes.DiskPool, error) {
	var poolCrs = make([]crtypes.DiskPool, 0)
	poolList, err := poolClientSet.DiskPools().List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return poolCrs, err
	}
	for _, poolCR := range poolList.Items {
		cr := poolCR
		poolCrs = append(poolCrs, v1beta3DSP{v1beta3: &cr})
	}
	return poolCrs, nil
}

func (ifc v1beta3Ifc) CreateMsPoolWithEncryption(poolName string, node string, disks []string, encryptionSecretName string) (crtypes.DiskPool, error) {
	logf.Log.Info("Creating DiskPool with encryption", "poolName", poolName, "node", node, "disks", disks, "encryptionSecretName", encryptionSecretName)
	msp := v1beta3.DiskPool{
		TypeMeta: metaV1.TypeMeta{Kind: "DiskPool"},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      poolName,
			Namespace: common.NSMayastor(),
		},
		Spec: v1beta3.DiskPoolSpec{
			Node:  node,
			Disks: disks,
			EncryptionConfig: &v1beta3.EncryptionConfig{
				Source: v1beta3.Source{
					Secret: v1beta3.Secret{
						Name: encryptionSecretName,
					},
				},
			},
		},
	}
	mspOut, err := poolClientSet.DiskPools().Create(context.TODO(), &msp, metaV1.CreateOptions{})
	dsp := v1beta3DSP{mspOut}
	return dsp, err
}

// AnnotatePoolForExpansion adds the openebs.io/expand: true annotation to enable pool expansion
func (ifc v1beta3Ifc) AnnotatePoolForExpansion(poolName string) error {
	res, err := poolClientSet.DiskPools().Get(context.TODO(), poolName, metaV1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get pool %s: %v", poolName, err)
	}
	if res == nil {
		return fmt.Errorf("pool %s not found", poolName)
	}

	// Add or update the expand annotation
	if res.Annotations == nil {
		res.Annotations = make(map[string]string)
	}
	res.Annotations["openebs.io/expand"] = "true"

	_, err = poolClientSet.DiskPools().Update(context.TODO(), res, metaV1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to annotate pool %s for expansion: %v", poolName, err)
	}

	logf.Log.Info("Successfully annotated pool for expansion", "pool", poolName)
	return nil
}

// VerifyPoolCapacityAndMaxExpansion verifies that the pool has the expected capacity and max expansion settings
func (ifc v1beta3Ifc) VerifyPoolCapacityAndMaxExpansion(poolName string, expectedCapacity uint64, expectedMaxExpansion string) error {
	pool, err := ifc.GetMsPool(poolName)
	if err != nil {
		return fmt.Errorf("failed to get pool %s: %v", poolName, err)
	}

	actualCapacity := pool.GetStatusCapacity()
	if actualCapacity != expectedCapacity {
		return fmt.Errorf("pool %s capacity mismatch: expected %d, got %d", poolName, expectedCapacity, actualCapacity)
	}

	actualMaxExpansion := pool.GetMaxExpansion()
	if actualMaxExpansion != expectedMaxExpansion {
		return fmt.Errorf("pool %s max expansion mismatch: expected %s, got %s", poolName, expectedMaxExpansion, actualMaxExpansion)
	}

	logf.Log.Info("Pool capacity and max expansion verified successfully",
		"pool", poolName,
		"capacity", actualCapacity,
		"maxExpansion", actualMaxExpansion)
	return nil
}

// AnnotateOfflinePoolForDelete adds the openebs.io/delete-opts annotation for offline pool deletion
// with specified options passed as strings
func (ifc v1beta3Ifc) AnnotateOfflinePoolForDelete(poolName string, opts ...string) error {
	res, err := poolClientSet.DiskPools().Get(context.TODO(), poolName, metaV1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get pool %s: %v", poolName, err)
	}
	if res == nil {
		return fmt.Errorf("pool %s not found", poolName)
	}

	// Build the annotation value based on provided options
	annotationValue := ifc.buildDeleteOptsAnnotation(opts...)

	if annotationValue == "" {
		return fmt.Errorf("at least one delete option must be specified")
	}

	// Add or update the delete-opts annotation for offline pool
	if res.Annotations == nil {
		res.Annotations = make(map[string]string)
	}
	res.Annotations[common.DeleteOptsAnnotation] = annotationValue

	_, err = poolClientSet.DiskPools().Update(context.TODO(), res, metaV1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to annotate offline pool %s for deletion: %v", poolName, err)
	}

	logf.Log.Info("Successfully annotated offline pool for deletion", "pool", poolName, "options", annotationValue)
	return nil
}

// buildDeleteOptsAnnotation constructs the YAML-formatted annotation value from option strings
func (ifc v1beta3Ifc) buildDeleteOptsAnnotation(opts ...string) string {
	var parts []string
	validOpts := map[string]string{
		"purge":                 "purge: true",
		"confirm":               "confirm: true",
		"confirm_data_loss":     "confirm_data_loss: true",
		"confirm_snapshot_loss": "confirm_snapshot_loss: true",
	}

	for _, opt := range opts {
		if value, exists := validOpts[opt]; exists {
			parts = append(parts, value)
		}
	}

	if len(parts) == 0 {
		return ""
	}

	return strings.Join(parts, "\n")
}
