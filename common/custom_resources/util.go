package custom_resources

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/e2e_config"

	apiExtV1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	crtypes "github.com/openebs/openebs-e2e/common/custom_resources/types"
	v1alpha1ExtIfc "github.com/openebs/openebs-e2e/common/custom_resources/v1alpha1Ext"
	v1beta1Ifc "github.com/openebs/openebs-e2e/common/custom_resources/v1beta1"
	v1beta2Ifc "github.com/openebs/openebs-e2e/common/custom_resources/v1beta2"

	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var csOnce sync.Once
var restConfig *rest.Config
var dspFuncsMap = make(map[string]*crtypes.DiskPoolFunctions, 0)
var selectedDspFuncs *crtypes.DiskPoolFunctions
var autoSelectCrdVersion = true

type extCRDVersion struct {
	version  string
	extended bool
	name     string
}

func (e extCRDVersion) String() string {
	return fmt.Sprintf("%s, extended:%v", e.version, e.extended)
}

/*
func contains(crdVers []extCRDVersion, ver extCRDVersion) bool {
	for _, elem := range crdVers {
		if elem == ver {
			return true
		}
	}
	return false
}
*/

func getDspCRDVersions(config *rest.Config, verbose bool) []extCRDVersion {
	clientSet, err := apiextensionsclientset.NewForConfig(config)
	crdVersions := make([]extCRDVersion, 0)
	if err == nil {
		var crDef *apiExtV1.CustomResourceDefinition
		crDef, err = clientSet.ApiextensionsV1().CustomResourceDefinitions().Get(context.TODO(), e2e_config.GetConfig().Product.PoolCrdName, metaV1.GetOptions{})
		if err == nil {
			serBytes, err := json.Marshal(crDef.Spec.Versions)
			if err == nil {
				fmt.Printf("DiskPool CRD.Spec.Versions:\n%s\n", string(serBytes))
			} else {
				if verbose {
					log.Log.Info("DiskPool", "CRD.Spec.Versions", crDef.Spec.Versions)
				}
			}
			for _, crdVer := range crDef.Spec.Versions {
				extended := false
				fullName := crdVer.Name
				if crdVer.Name == "v1alpha1" {
					extended = v1alpha1ExtIfc.IsExtended(crdVer.Schema.OpenAPIV3Schema.Properties["status"].Properties)
					if extended {
						fullName = fullName + "Ext"
					}
				}
				ecv := extCRDVersion{
					version:  crdVer.Name,
					extended: extended,
					name:     fullName,
				}
				crdVersions = append(crdVersions, ecv)
			}
		} else {
			log.Log.Info("Failed to retrieve CRD", "error", err)
		}
	}
	return crdVersions
}

func GetDiskpoolCrdVersions() ([]string, error) {
	var versions = make([]string, 0)
	useCluster := true
	testEnv := &envtest.Environment{
		UseExistingCluster: &useCluster,
	}
	config, err := testEnv.Start()
	if err == nil {
		crdVersions := getDspCRDVersions(config, true)
		for _, cv := range crdVersions {
			versions = append(versions, cv.version)
		}
	}
	return versions, err
}

// initialise the custom resources package
// Note: if the CRD versions defined in the k8s cluster
// are unsupported - we will panic on nil pointer de-reference
// For now that is deemed acceptable
func initialise() {
	csOnce.Do(func() {
		useCluster := true
		testEnv := &envtest.Environment{
			UseExistingCluster: &useCluster,
		}
		var err error
		restConfig, err = testEnv.Start()
		if err != nil {
			log.Log.Info("custom_resources: Initialise", "error", err)
		}
		// list of versions sorted in order of preference
		dspCrdVers := []struct {
			description  string
			version      extCRDVersion
			fnInitialise func(*rest.Config, bool) (crtypes.DiskPoolFunctions, error)
		}{
			{
				description: "v1beta2",
				version: extCRDVersion{
					"v1beta2", false, "v1beta2",
				},
				fnInitialise: v1beta2Ifc.Initialise,
			},
			{
				description: "v1beta1",
				version: extCRDVersion{
					"v1beta1", false, "v1beta1",
				},
				fnInitialise: v1beta1Ifc.Initialise,
			},
			{
				description: "v1alpha1 extended",
				version: extCRDVersion{
					"v1alpha1", true, "v1alpha1Ext",
				},
				fnInitialise: v1alpha1ExtIfc.Initialise,
			},
			{
				description: "v1alpha1",
				version: extCRDVersion{
					"v1alpha1", false, "v1alpha1",
				},
				fnInitialise: v1alpha1ExtIfc.Initialise,
			},
		}
		for _, dcv := range dspCrdVers {
			var dspFns crtypes.DiskPoolFunctions
			dspFns, err = dcv.fnInitialise(restConfig, dcv.version.extended)
			dspFuncsMap[dcv.version.name] = &dspFns
			if err != nil {
				log.Log.Info("Failed to setup for DiskPoolCrd", "version", dcv.version, "Error", err)
			} else {
				log.Log.Info("*** Initialised DiskPool CRD", "version", dcv.version, "description", dcv.description)
			}
		}
	},
	)
}

var crdVersionOrder = []string{"v1beta2", "v1beta1", "v1alpha1Ext", "v1alpha1"}

func getDspFuncs() crtypes.DiskPoolFunctions {
	if selectedDspFuncs != nil {
		return *selectedDspFuncs
	}
	initialise()
	crdVersions := getDspCRDVersions(restConfig, false)
	// select the latest version
	for _, crdVer := range crdVersions {
		for _, crdVerName := range crdVersionOrder {
			if crdVerName == crdVer.name {
				dspFns := dspFuncsMap[crdVerName]
				log.Log.Info(fmt.Sprintf(" Using DiskPool CRD version %s", crdVerName))
				if autoSelectCrdVersion {
					selectedDspFuncs = dspFns
				}
				return *dspFns
			}
		}
	}
	log.Log.Info(fmt.Sprintf("unknown DiskPool CRD %s", crdVersions[0].name))
	// unknown version return nil this will cause a panic - which is deemed acceptable
	return nil
}

// - AutoSelectDiskPoolCRD - select the pool CRD version and API to use
// based on the CRD(s) available on the target system
func AutoSelectDiskPoolCRD() {
	autoSelectCrdVersion = true
}

// ClearDiskPoolCRDSelection - "unselect" the pool CRD version and API to use,
// the process of selecting the Disk Pool CRD version to use is executed on every call,
// this may impact stress tests or tests with concurrency negatively.
func ClearDiskPoolCRDSelection() {
	log.Log.Info("Fixed DiskPool CRD is cleared")
	selectedDspFuncs = nil
	autoSelectCrdVersion = false
}

// == Mayastor Pool  ======================

func CreateMsPool(poolName string, node string, disks []string) (crtypes.DiskPool, error) {
	return getDspFuncs().CreateMsPool(poolName, node, disks)
}

func CreateMsPoolWithTopologySpec(poolName string, node string, disks []string, labels map[string]string) (crtypes.DiskPool, error) {
	return getDspFuncs().CreateMsPoolWithTopologySpec(poolName, node, disks, labels)
}

func GetMsPool(poolName string) (crtypes.DiskPool, error) {
	return getDspFuncs().GetMsPool(poolName)
}

func DeleteMsPool(poolName string) error {
	return getDspFuncs().DeleteMsPool(poolName)
}

func ListMsPools() ([]crtypes.DiskPool, error) {
	return getDspFuncs().ListMsPoolCrs()
}

// CheckAllMsPoolsAreOnline checks if all mayastor pools are online
func CheckAllMsPoolsAreOnline() error {
	initialise()
	allHealthy := true
	pools, err := ListMsPools()
	var synopsis string
	if err == nil && pools != nil && len(pools) != 0 {
		for _, pool := range pools {
			poolName := pool.GetName()
			poolStatus := pool.GetPoolStatus()
			synopsis = fmt.Sprintf("%spool=%s poolStatus=%s\n", synopsis, poolName, poolStatus)
			if strings.ToLower(poolStatus) != "online" {
				log.Log.Info("CheckAllMsPoolsAreOnline", "pool", poolName, "poolStatus", poolStatus)
				allHealthy = false
			}
		}
	}

	if !allHealthy {
		return fmt.Errorf("all pools were not online\n%s", synopsis)
	}
	return err
}

// CheckDiskPoolsIsOnline checks if mayastor pool is online
func CheckDiskPoolsIsOnline(poolName string) (bool, error) {
	initialise()
	online := false

	pool, err := GetMsPool(poolName)
	if err == nil {
		state := pool.GetPoolStatus()
		if strings.ToLower(state) == "online" {
			log.Log.Info("CheckMsPoolsAreOnline", "pool", poolName, "state", state)
			online = true
		}
	}

	return online, err
}

// CheckAllMsPoolFinalizers check
//  1. that finalizers exist for pools with replicas (used size != 0)
//  2. that finalizers DO NOT EXIST for pools with no replicas (used size == 0)
//
// Note this function should not be call if mayastor is deployed with control plane
// versions > 0
func CheckAllMsPoolFinalizers() error {
	initialise()
	var errs common.ErrorAccumulator
	pools, err := ListMsPools()
	if err != nil {
		return err
	}

	for _, pool := range pools {
		finalizer := pool.GetFinalizers()
		if finalizer != nil {
			errs.Accumulate(fmt.Errorf("finalizer set on pool: %s", pool.GetName()))
		}
	}

	return errs.GetError()
}

func GetDiskPoolCrStatus(poolName string) (string, error) {
	initialise()
	dsp, err := GetMsPool(poolName)
	return dsp.GetCRStatus(), err
}

func GetDiskPoolStatus(poolName string) (string, error) {
	initialise()
	dsp, err := GetMsPool(poolName)
	return dsp.GetPoolStatus(), err
}
