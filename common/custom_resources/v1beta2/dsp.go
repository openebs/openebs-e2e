package v1beta2

import (
	"context"
	"fmt"
	"reflect"

	"github.com/openebs/openebs-e2e/common"

	logf "sigs.k8s.io/controller-runtime/pkg/log"

	v1beta2 "github.com/openebs/openebs-e2e/common/custom_resources/api/types/v1beta2"
	v1beta2Client "github.com/openebs/openebs-e2e/common/custom_resources/clientset/v1beta2"
	crtypes "github.com/openebs/openebs-e2e/common/custom_resources/types"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

var poolClientSet *v1beta2Client.DiskPoolV1Beta2Client

type v1beta2Ifc struct {
}

type v1beta2DSP struct {
	v1beta2 *v1beta2.DiskPool
}

func Initialise(config *rest.Config, extended bool) (crtypes.DiskPoolFunctions, error) {
	var err error
	_ = v1beta2.PoolAddToScheme(scheme.Scheme)
	poolClientSet, err = v1beta2Client.DspNewForConfig(config)
	return v1beta2Ifc{}, err
}

// DiskPool implementation

func (p v1beta2DSP) String() string {
	return fmt.Sprintf("%v", p.v1beta2)
}

func (p v1beta2DSP) GetType() reflect.Type {
	return reflect.TypeOf(p)
}

func (p v1beta2DSP) GetName() string {
	if p.v1beta2 != nil {
		return p.v1beta2.GetName()
	}
	// panic?
	return ""
}

// For v1beta2 map Status.State to pool status
func (p v1beta2DSP) GetPoolStatus() string {
	if p.v1beta2 != nil {
		return p.v1beta2.Status.PoolStatus
	}
	return ""
}

// For v1beta2 map Status.State to CR status
func (p v1beta2DSP) GetCRStatus() string {
	if p.v1beta2 != nil {
		return p.v1beta2.Status.CRStatus
	}
	return ""
}

func (p v1beta2DSP) GetStatusCapacity() uint64 {
	if p.v1beta2 != nil {
		return p.v1beta2.Status.Capacity
	}
	return 0
}

func (p v1beta2DSP) GetStatusUsed() uint64 {
	if p.v1beta2 != nil {
		return p.v1beta2.Status.Used
	}
	return 0
}

func (p v1beta2DSP) CompareStatus(otherP *crtypes.DiskPool) bool {
	other := *otherP
	if p.GetType() != other.GetType() {
		logf.Log.Info("comparison between different diskpool types is unsupported",
			"this", p.GetType(), "other", other.GetType())
		return false
	}
	dsp := other.(v1beta2DSP)
	if p.v1beta2 != nil && dsp.v1beta2 != nil {
		return reflect.DeepEqual(p.v1beta2.Status, dsp.v1beta2.Status)
	}
	return false
}

func (p v1beta2DSP) GetFinalizers() []string {
	if p.v1beta2 != nil {
		return p.v1beta2.Finalizers
	}
	// panic?
	return []string{}
}

func (p v1beta2DSP) SetFinalizers(finalizers []string) (crtypes.DiskPool, error) {
	var err error
	dsp := p
	if p.v1beta2 != nil {
		mspIn := *p.v1beta2
		mspIn.SetFinalizers(finalizers)
		dsp.v1beta2, err = poolClientSet.DiskPools().Update(context.TODO(), &mspIn, metaV1.UpdateOptions{})
		return dsp, err
	}
	return dsp, fmt.Errorf("uninitialised DiskPool")
}

func (p v1beta2DSP) GetSpecDisks() []string {
	if p.v1beta2 != nil {
		return p.v1beta2.Spec.Disks
	}
	return []string{}
}

func (p v1beta2DSP) SetSpecDisks(disks []string) (crtypes.DiskPool, error) {
	var err error
	dsp := p
	if p.v1beta2 != nil {
		mspIn := *p.v1beta2
		mspIn.Spec.Disks = disks
		dsp.v1beta2, err = poolClientSet.DiskPools().Update(context.TODO(), &mspIn, metaV1.UpdateOptions{})
		return dsp, err
	}
	return dsp, fmt.Errorf("uninitialised DiskPool")
}

func (p v1beta2DSP) SetSpecNode(node string) (crtypes.DiskPool, error) {
	var err error
	dsp := p
	if p.v1beta2 != nil {
		mspIn := *p.v1beta2
		mspIn.Spec.Node = node
		dsp.v1beta2, err = poolClientSet.DiskPools().Update(context.TODO(), &mspIn, metaV1.UpdateOptions{})
		return dsp, err
	}
	return dsp, fmt.Errorf("uninitialised DiskPool")
}

func (p v1beta2DSP) GetSpecNode() string {
	if p.v1beta2 != nil {
		return p.v1beta2.Spec.Node
	}
	return ""
}

//  DiskPoolFunctions implementation

func (ifc v1beta2Ifc) CreateMsPool(poolName string, node string, disks []string) (crtypes.DiskPool, error) {
	msp := v1beta2.DiskPool{
		TypeMeta: metaV1.TypeMeta{Kind: "DiskPool"},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      poolName,
			Namespace: common.NSMayastor(),
		},
		Spec: v1beta2.DiskPoolSpec{
			Node:  node,
			Disks: disks,
		},
	}
	mspOut, err := poolClientSet.DiskPools().Create(context.TODO(), &msp, metaV1.CreateOptions{})
	dsp := v1beta2DSP{mspOut}
	return dsp, err
}

func (ifc v1beta2Ifc) CreateMsPoolWithTopologySpec(poolName string, node string, disks []string, labels map[string]string) (crtypes.DiskPool, error) {

	topology := &v1beta2.Topology{
		Labelled: labels,
	}

	msp := v1beta2.DiskPool{
		TypeMeta: metaV1.TypeMeta{Kind: "DiskPool"},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      poolName,
			Namespace: common.NSMayastor(),
		},
		Spec: v1beta2.DiskPoolSpec{
			Node:     node,
			Disks:    disks,
			Topology: topology,
		},
	}

	mspOut, err := poolClientSet.DiskPools().Create(context.TODO(), &msp, metaV1.CreateOptions{})
	dsp := v1beta2DSP{mspOut}
	return dsp, err
}

func (ifc v1beta2Ifc) GetMsPool(poolName string) (crtypes.DiskPool, error) {
	msp := v1beta2.DiskPool{}
	res, err := poolClientSet.DiskPools().Get(context.TODO(), poolName, metaV1.GetOptions{})
	if res != nil && err == nil {
		msp = *res
	}
	return v1beta2DSP{v1beta2: &msp}, err
}

func (ifc v1beta2Ifc) DeleteMsPool(poolName string) error {
	err := poolClientSet.DiskPools().Delete(context.TODO(), poolName, metaV1.DeleteOptions{})
	return err
}

func (ifc v1beta2Ifc) ListMsPoolCrs() ([]crtypes.DiskPool, error) {
	var poolCrs []crtypes.DiskPool = make([]crtypes.DiskPool, 0)
	poolList, err := poolClientSet.DiskPools().List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return poolCrs, err
	}
	for _, poolCR := range poolList.Items {
		cr := poolCR
		poolCrs = append(poolCrs, v1beta2DSP{v1beta2: &cr})
	}
	return poolCrs, nil
}
