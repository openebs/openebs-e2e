package v1beta1

import (
	"context"
	"fmt"
	"reflect"

	"github.com/openebs/openebs-e2e/common"

	logf "sigs.k8s.io/controller-runtime/pkg/log"

	v1beta12 "github.com/openebs/openebs-e2e/common/custom_resources/api/types/v1beta1"
	v1beta1Client "github.com/openebs/openebs-e2e/common/custom_resources/clientset/v1beta1"
	crtypes "github.com/openebs/openebs-e2e/common/custom_resources/types"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

var poolClientSet *v1beta1Client.DiskPoolV1Beta1Client

type v1beta1Ifc struct {
}

type v1beta1DSP struct {
	v1beta1 *v1beta12.DiskPool
}

func Initialise(config *rest.Config, extended bool) (crtypes.DiskPoolFunctions, error) {
	var err error
	_ = v1beta12.PoolAddToScheme(scheme.Scheme)
	poolClientSet, err = v1beta1Client.DspNewForConfig(config)
	return v1beta1Ifc{}, err
}

// DiskPool implementation

func (p v1beta1DSP) String() string {
	return fmt.Sprintf("%v", p.v1beta1)
}

func (p v1beta1DSP) GetType() reflect.Type {
	return reflect.TypeOf(p)
}

func (p v1beta1DSP) GetName() string {
	if p.v1beta1 != nil {
		return p.v1beta1.GetName()
	}
	// panic?
	return ""
}

// For v1beta1 map Status.State to pool status
func (p v1beta1DSP) GetPoolStatus() string {
	if p.v1beta1 != nil {
		return p.v1beta1.Status.PoolStatus
	}
	return ""
}

// For v1beta1 map Status.State to CR status
func (p v1beta1DSP) GetCRStatus() string {
	if p.v1beta1 != nil {
		return p.v1beta1.Status.CRStatus
	}
	return ""
}

func (p v1beta1DSP) GetStatusCapacity() uint64 {
	if p.v1beta1 != nil {
		return p.v1beta1.Status.Capacity
	}
	return 0
}

func (p v1beta1DSP) GetStatusUsed() uint64 {
	if p.v1beta1 != nil {
		return p.v1beta1.Status.Used
	}
	return 0
}

func (p v1beta1DSP) CompareStatus(otherP *crtypes.DiskPool) bool {
	other := *otherP
	if p.GetType() != other.GetType() {
		logf.Log.Info("comparison between different diskpool types is unsupported",
			"this", p.GetType(), "other", other.GetType())
		return false
	}
	dsp := other.(v1beta1DSP)
	if p.v1beta1 != nil && dsp.v1beta1 != nil {
		return reflect.DeepEqual(p.v1beta1.Status, dsp.v1beta1.Status)
	}
	return false
}

func (p v1beta1DSP) GetFinalizers() []string {
	if p.v1beta1 != nil {
		return p.v1beta1.Finalizers
	}
	// panic?
	return []string{}
}

func (p v1beta1DSP) SetFinalizers(finalizers []string) (crtypes.DiskPool, error) {
	var err error
	dsp := p
	if p.v1beta1 != nil {
		mspIn := *p.v1beta1
		mspIn.SetFinalizers(finalizers)
		dsp.v1beta1, err = poolClientSet.DiskPools().Update(context.TODO(), &mspIn, metaV1.UpdateOptions{})
		return dsp, err
	}
	return dsp, fmt.Errorf("uninitialised DiskPool")
}

func (p v1beta1DSP) GetSpecDisks() []string {
	if p.v1beta1 != nil {
		return p.v1beta1.Spec.Disks
	}
	return []string{}
}

func (p v1beta1DSP) SetSpecDisks(disks []string) (crtypes.DiskPool, error) {
	var err error
	dsp := p
	if p.v1beta1 != nil {
		mspIn := *p.v1beta1
		mspIn.Spec.Disks = disks
		dsp.v1beta1, err = poolClientSet.DiskPools().Update(context.TODO(), &mspIn, metaV1.UpdateOptions{})
		return dsp, err
	}
	return dsp, fmt.Errorf("uninitialised DiskPool")
}

func (p v1beta1DSP) SetSpecNode(node string) (crtypes.DiskPool, error) {
	var err error
	dsp := p
	if p.v1beta1 != nil {
		mspIn := *p.v1beta1
		mspIn.Spec.Node = node
		dsp.v1beta1, err = poolClientSet.DiskPools().Update(context.TODO(), &mspIn, metaV1.UpdateOptions{})
		return dsp, err
	}
	return dsp, fmt.Errorf("uninitialised DiskPool")
}

func (p v1beta1DSP) GetSpecNode() string {
	if p.v1beta1 != nil {
		return p.v1beta1.Spec.Node
	}
	return ""
}

//  DiskPoolFunctions implementation

func (ifc v1beta1Ifc) CreateMsPool(poolName string, node string, disks []string) (crtypes.DiskPool, error) {
	msp := v1beta12.DiskPool{
		TypeMeta: metaV1.TypeMeta{Kind: "DiskPool"},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      poolName,
			Namespace: common.NSMayastor(),
		},
		Spec: v1beta12.DiskPoolSpec{
			Node:  node,
			Disks: disks,
		},
	}
	mspOut, err := poolClientSet.DiskPools().Create(context.TODO(), &msp, metaV1.CreateOptions{})
	dsp := v1beta1DSP{mspOut}
	return dsp, err
}

func (ifc v1beta1Ifc) CreateMsPoolWithTopologySpec(poolName string, node string, disks []string, labels map[string]string) (crtypes.DiskPool, error) {
	panic(fmt.Errorf("not implemented labelled topolgy pool in v1beta1"))
}

func (ifc v1beta1Ifc) GetMsPool(poolName string) (crtypes.DiskPool, error) {
	msp := v1beta12.DiskPool{}
	res, err := poolClientSet.DiskPools().Get(context.TODO(), poolName, metaV1.GetOptions{})
	if res != nil && err == nil {
		msp = *res
	}
	return v1beta1DSP{v1beta1: &msp}, err
}

func (ifc v1beta1Ifc) DeleteMsPool(poolName string) error {
	err := poolClientSet.DiskPools().Delete(context.TODO(), poolName, metaV1.DeleteOptions{})
	return err
}

func (ifc v1beta1Ifc) ListMsPoolCrs() ([]crtypes.DiskPool, error) {
	var poolCrs []crtypes.DiskPool = make([]crtypes.DiskPool, 0)
	poolList, err := poolClientSet.DiskPools().List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return poolCrs, err
	}
	for _, poolCR := range poolList.Items {
		cr := poolCR
		poolCrs = append(poolCrs, v1beta1DSP{v1beta1: &cr})
	}
	return poolCrs, nil
}
