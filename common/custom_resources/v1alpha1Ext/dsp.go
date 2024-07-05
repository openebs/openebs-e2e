package v1alpha1Ext

import (
	"context"
	"fmt"
	"reflect"

	"github.com/openebs/openebs-e2e/common"
	v1alpha12 "github.com/openebs/openebs-e2e/common/custom_resources/api/types/v1alpha1Ext"
	v1alpha1Client "github.com/openebs/openebs-e2e/common/custom_resources/clientset/v1alpha1Ext"
	crtypes "github.com/openebs/openebs-e2e/common/custom_resources/types"

	logf "sigs.k8s.io/controller-runtime/pkg/log"

	apiExtV1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

var poolClientSet *v1alpha1Client.DiskPoolV1Alpha1Client

type v1alpha1ExtIfc struct {
	extended bool
}

type v1alpha1ExtDSP struct {
	v1alpha1 *v1alpha12.DiskPool
	extended bool
}

func Initialise(config *rest.Config, extended bool) (crtypes.DiskPoolFunctions, error) {
	var err error
	_ = v1alpha12.PoolAddToScheme(scheme.Scheme)
	poolClientSet, err = v1alpha1Client.DspNewForConfig(config)
	return v1alpha1ExtIfc{extended: extended}, err
}

// DiskPool implementation

func (p v1alpha1ExtDSP) String() string {
	return fmt.Sprintf("%v", p.v1alpha1)
}

func (p v1alpha1ExtDSP) GetType() reflect.Type {
	return reflect.TypeOf(p)
}

func (p v1alpha1ExtDSP) GetName() string {
	if p.v1alpha1 != nil {
		return p.v1alpha1.GetName()
	}
	// panic?
	return ""
}

func (p v1alpha1ExtDSP) GetPoolStatus() string {
	if p.v1alpha1 != nil {
		if p.extended {
			return p.v1alpha1.Status.PoolStatus
		} else {
			// For v1alpha1 map Status.State to pool status
			return p.v1alpha1.Status.State
		}
	}
	return ""
}

func (p v1alpha1ExtDSP) GetCRStatus() string {
	if p.v1alpha1 != nil {
		if p.extended {
			return p.v1alpha1.Status.CRStatus
		} else {
			// For v1alpha1 map Status.State to CR status
			return p.v1alpha1.Status.State
		}
	}
	return ""
}

func (p v1alpha1ExtDSP) GetStatusCapacity() uint64 {
	if p.v1alpha1 != nil {
		return p.v1alpha1.Status.Capacity
	}
	return 0
}

func (p v1alpha1ExtDSP) GetStatusUsed() uint64 {
	if p.v1alpha1 != nil {
		return p.v1alpha1.Status.Used
	}
	return 0
}

func (p v1alpha1ExtDSP) CompareStatus(otherP *crtypes.DiskPool) bool {
	other := *otherP
	if p.GetType() != other.GetType() {
		logf.Log.Info("comparison between different diskpool types is unsupported",
			"this", p.GetType(), "other", other.GetType())
		return false
	}
	dsp := other.(v1alpha1ExtDSP)
	if p.v1alpha1 != nil && dsp.v1alpha1 != nil {
		return reflect.DeepEqual(p.v1alpha1.Status, dsp.v1alpha1.Status)
	}
	return false
}

func (p v1alpha1ExtDSP) GetFinalizers() []string {
	if p.v1alpha1 != nil {
		return p.v1alpha1.Finalizers
	}
	// panic?
	return []string{}
}

func (p v1alpha1ExtDSP) SetFinalizers(finalizers []string) (crtypes.DiskPool, error) {
	var err error
	dsp := p
	if p.v1alpha1 != nil {
		mspIn := *p.v1alpha1
		mspIn.SetFinalizers(finalizers)
		dsp.v1alpha1, err = poolClientSet.DiskPools().Update(context.TODO(), &mspIn, metaV1.UpdateOptions{})
		return dsp, err
	}
	return dsp, fmt.Errorf("uninitialised DiskPool")
}

func (p v1alpha1ExtDSP) GetSpecDisks() []string {
	if p.v1alpha1 != nil {
		return p.v1alpha1.Spec.Disks
	}
	return []string{}
}

func (p v1alpha1ExtDSP) SetSpecDisks(disks []string) (crtypes.DiskPool, error) {
	var err error
	dsp := p
	if p.v1alpha1 != nil {
		mspIn := *p.v1alpha1
		mspIn.Spec.Disks = disks
		dsp.v1alpha1, err = poolClientSet.DiskPools().Update(context.TODO(), &mspIn, metaV1.UpdateOptions{})
		return dsp, err
	}
	return dsp, fmt.Errorf("uninitialised DiskPool")
}

func (p v1alpha1ExtDSP) SetSpecNode(node string) (crtypes.DiskPool, error) {
	var err error
	dsp := p
	if p.v1alpha1 != nil {
		mspIn := *p.v1alpha1
		mspIn.Spec.Node = node
		dsp.v1alpha1, err = poolClientSet.DiskPools().Update(context.TODO(), &mspIn, metaV1.UpdateOptions{})
		return dsp, err
	}
	return dsp, fmt.Errorf("uninitialised DiskPool")
}

func (p v1alpha1ExtDSP) GetSpecNode() string {
	if p.v1alpha1 != nil {
		return p.v1alpha1.Spec.Node
	}
	return ""
}

//  DiskPoolFunctions implementation

func (ifc v1alpha1ExtIfc) CreateMsPool(poolName string, node string, disks []string) (crtypes.DiskPool, error) {
	msp := v1alpha12.DiskPool{
		TypeMeta: metaV1.TypeMeta{Kind: "DiskPool"},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      poolName,
			Namespace: common.NSMayastor(),
		},
		Spec: v1alpha12.DiskPoolSpec{
			Node:  node,
			Disks: disks,
		},
	}
	mspOut, err := poolClientSet.DiskPools().Create(context.TODO(), &msp, metaV1.CreateOptions{})
	dsp := v1alpha1ExtDSP{mspOut, ifc.extended}
	return dsp, err
}

func (ifc v1alpha1ExtIfc) CreateMsPoolWithTopologySpec(poolName string, node string, disks []string, labels map[string]string) (crtypes.DiskPool, error) {
	panic(fmt.Errorf("not implemented labelled topolgy pool in v1alpha1Ext"))
}

func (ifc v1alpha1ExtIfc) GetMsPool(poolName string) (crtypes.DiskPool, error) {
	msp := v1alpha12.DiskPool{}
	res, err := poolClientSet.DiskPools().Get(context.TODO(), poolName, metaV1.GetOptions{})
	if res != nil && err == nil {
		msp = *res
	}
	return v1alpha1ExtDSP{v1alpha1: &msp, extended: ifc.extended}, err
}

func (ifc v1alpha1ExtIfc) DeleteMsPool(poolName string) error {
	err := poolClientSet.DiskPools().Delete(context.TODO(), poolName, metaV1.DeleteOptions{})
	return err
}

func (ifc v1alpha1ExtIfc) ListMsPoolCrs() ([]crtypes.DiskPool, error) {
	var poolCrs []crtypes.DiskPool = make([]crtypes.DiskPool, 0)
	poolList, err := poolClientSet.DiskPools().List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return poolCrs, err
	}
	for _, poolCR := range poolList.Items {
		cr := poolCR
		poolCrs = append(poolCrs, v1alpha1ExtDSP{v1alpha1: &cr})
	}
	return poolCrs, nil
}

func IsExtended(jsProps map[string]apiExtV1.JSONSchemaProps) bool {
	// Create a dummy DiskPoolStatus initialising the fields we want to
	// access, so that we generate a compile-time failure if the names
	// change
	dsp := v1alpha12.DiskPoolStatus{State: "", CRStatus: "", PoolStatus: ""}
	rt := reflect.TypeOf(dsp)
	// the fields DO exist (see comment at the start),
	// retrieve ignoring presence checks
	fieldJsonTag := func(rt reflect.Type, fieldName string) string {
		fld, _ := rt.FieldByName(fieldName)
		return fld.Tag.Get("json")
	}
	extended := true
	for _, fieldName := range []string{"State", "CRStatus", "PoolStatus"} {
		jsonTag := fieldJsonTag(rt, fieldName)
		logf.Log.Info("v1alpha1Ext.IsExtended", "fieldName", fieldName, "jsonTag", jsonTag)
		_, present := jsProps[jsonTag]
		if !present {
			extended = false
		}
	}
	return extended
}
