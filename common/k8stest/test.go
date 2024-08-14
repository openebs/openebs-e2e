package k8stest

import (
	"fmt"
	"os"
	"time"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/custom_resources"
	"github.com/openebs/openebs-e2e/common/e2e_config"
	"github.com/openebs/openebs-e2e/common/mayastorclient"

	clientset "github.com/kubernetes-csi/external-snapshotter/client/v6/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

type TestEnvironment struct {
	KubeInt kubernetes.Interface
	CsiInt  clientset.Interface
}

var gTestEnv TestEnvironment

func SetupK8sEnvBasic() error {

	restConfig := config.GetConfigOrDie()
	if restConfig == nil {
		return fmt.Errorf("failed to create *rest.Config for talking to a Kubernetes apiserver : GetConfigOrDie")
	}
	kubeInt := kubernetes.NewForConfigOrDie(restConfig)
	if kubeInt == nil {
		return fmt.Errorf("failed to create new Clientset for the given config : NewForConfigOrDie")
	}
	csiInt := clientset.NewForConfigOrDie(restConfig)
	if csiInt == nil {
		return fmt.Errorf("failed to create new csi Clientset for the given config : NewForConfigOrDie")
	}
	gTestEnv = TestEnvironment{
		KubeInt: kubeInt,
		CsiInt:  csiInt,
	}

	if _, ok := os.LookupEnv("e2e_product"); !ok {
		product, err := DiscoverProduct()
		if err != nil {
			return fmt.Errorf("e2e_product not set, failed to discover product on cluster")
		}
		if err = os.Setenv("e2e_product", product); err != nil {
			return fmt.Errorf("failed to set envvar e2e_product: %v", err)
		}
	}
	// Check if gRPC calls are possible and store the result
	// subsequent calls to mayastorClient.CanConnect retrieves
	// the result.
	mayastorclient.Initialise(GetMayastorNodeIPAddresses())
	return nil
}

func SetupK8sEnv() error {
	if e2e_config.GetConfig().ReplicatedEngine {
		err := CheckAndSetControlPlane()
		if err != nil {
			return fmt.Errorf("failed to setup control plane version : CheckAndSetControlPlane: %v", err)
		}

		// Fail the test setup if gRPC calls are mandated and
		// gRPC calls are not supported.
		if e2e_config.GetConfig().GrpcMandated {
			grpcCalls := mayastorclient.CanConnect()
			if !grpcCalls {
				return fmt.Errorf("gRPC calls to mayastor are disabled, but mandated by configuration : CanConnect: %v", grpcCalls)
			}
		}
	}
	return nil
}

// SetupK8sLib convenience function for applications using common code
// go code using ginkgo should not invoke this function.
func SetupK8sLib() error {
	e2e_config.SetContext(e2e_config.Library)

	err := SetupK8sEnvBasic()
	if err != nil {
		return err
	}
	return SetupK8sEnv()
}

func TeardownTestEnvNoCleanup() error {
	// TODO - remove this function when we are sure we won't need it
	return nil
}

func TeardownTestEnv() error {
	AfterSuiteCleanup()
	err := TeardownTestEnvNoCleanup()
	if err != nil {
		return err
	}
	return nil
}

// AfterSuiteCleanup  placeholder function for now
// To aid postmortem analysis for the most common CI use case
// namely cluster is retained aon failure, we do nothing
// For other situations behaviour should be configurable
func AfterSuiteCleanup() {
	logf.Log.Info("AfterSuiteCleanup")
}

func getMspUsage() (uint64, error) {
	var mspUsage uint64
	msPools, err := ListMsPools()
	if err != nil {
		logf.Log.Info("unable to list mayastor pools", "error", err)
	} else {
		mspUsage = 0
		for _, pool := range msPools {
			mspUsage += pool.Status.Used
		}
	}
	return mspUsage, err
}

// ResourceCheck  Fit for purpose checks
// - No pods
// - No PVCs
// - No PVs
// - No MSVs
// - Mayastor pods are all healthy
// - All mayastor pools are online
// and if e2e-agent is available
// - mayastor pools usage is 0
// - No nexuses
// - No replicas
func ResourceCheck(waitForPools bool) error {
	var errs = common.ErrorAccumulator{}

	// Check that Mayastor pods are healthy no restarts or fails.
	err := CheckTestPodsHealth(common.NSMayastor())
	if err != nil {
		if e2e_config.GetConfig().SelfTest {
			logf.Log.Info("SelfTesting, ignoring:", "", err)
		} else {
			errs.Accumulate(err)
		}
	}

	pods, err := CheckForTestPods()
	errs.Accumulate(err)
	if pods {
		errs.Accumulate(fmt.Errorf("found Pods"))
	}

	pvcs, err := CheckForPVCs()
	errs.Accumulate(err)
	if pvcs {
		errs.Accumulate(fmt.Errorf("found PersistentVolumeClaims"))
	}

	pvs, err := CheckForPVs()
	errs.Accumulate(err)

	if pvs {
		errs.Accumulate(fmt.Errorf("found PersistentVolumes"))
	}

	msvs, err := ListMsvs()
	if err != nil {
		errs.Accumulate(err)
	} else {
		if msvs != nil {
			if len(msvs) != 0 {
				errs.Accumulate(fmt.Errorf("found MayastorVolumes"))
			}
		} else {
			logf.Log.Info("Listing MSVs returned nil array")
		}
	}

	scs, err := CheckForStorageClasses()
	errs.Accumulate(err)
	if scs {
		errs.Accumulate(fmt.Errorf("found storage classes using mayastor"))
	}

	err = custom_resources.CheckAllMsPoolsAreOnline()
	if err != nil {
		errs.Accumulate(err)
		logf.Log.Info("ResourceCheck: not all pools are online")
	}

	{
		mspUsage, err := getMspUsage()
		// skip waiting if fail quick and errors already exist
		skip := e2e_config.GetConfig().FailQuick && errs.GetError() != nil
		if (err != nil || mspUsage != 0) && waitForPools && !skip {
			logf.Log.Info("Waiting for pool usage to be 0")
			const sleepTime = 10
			t0 := time.Now()
			// Wait for pool usage reported by CRS to drop to 0
			for ix := 0; ix < (60*sleepTime) && mspUsage != 0; ix += sleepTime {
				time.Sleep(sleepTime * time.Second)
				mspUsage, err = getMspUsage()
				if err != nil {
					logf.Log.Info("ResourceCheck: unable to list msps")
				}
			}
			logf.Log.Info("ResourceCheck:", "mspool Usage", mspUsage, "waiting time", time.Since(t0))
			errs.Accumulate(err)
		}
		if mspUsage != 0 {
			errs.Accumulate(fmt.Errorf("pool usage reported via custom resources %d", mspUsage))
		}
		logf.Log.Info("ResourceCheck:", "mspool Usage", mspUsage)
	}

	// gRPC calls can only be executed successfully is the e2e-agent daemonSet has been deployed successfully.
	if mayastorclient.CanConnect() {
		// check pools
		{
			poolUsage, err := GetPoolUsageInCluster()
			// skip waiting if fail quick and errors already exist
			skip := e2e_config.GetConfig().FailQuick && errs.GetError() != nil
			if (err != nil || poolUsage != 0) && waitForPools && !skip {
				logf.Log.Info("Waiting for pool usage to be 0 (gRPC)")
				const sleepTime = 2
				t0 := time.Now()
				// Wait for pool usage to drop to 0
				for ix := 0; ix < 120 && poolUsage != 0; ix += sleepTime {
					time.Sleep(sleepTime * time.Second)
					poolUsage, err = GetPoolUsageInCluster()
					if err != nil {
						logf.Log.Info("ResourceEachCheck: failed to retrieve pools usage")
					}
				}
				logf.Log.Info("ResourceCheck:", "poolUsage", poolUsage, "waiting time", time.Since(t0))
			}
			errs.Accumulate(err)
			if poolUsage != 0 {
				errs.Accumulate(fmt.Errorf("gRPC: pool usage reported via custom resources %d", poolUsage))
			}
			logf.Log.Info("ResourceCheck:", "poolUsage", poolUsage)
		}
		// check nexuses
		{
			nexuses, err := ListNexusesInCluster()
			if err != nil {
				errs.Accumulate(err)
				logf.Log.Info("ResourceEachCheck: failed to retrieve list of nexuses")
			}
			logf.Log.Info("ResourceCheck:", "num nexuses", len(nexuses))
			if len(nexuses) != 0 {
				errs.Accumulate(fmt.Errorf("gRPC: count of nexuses reported via mayastor client is %d", len(nexuses)))
			}
		}
		// check replicas
		{
			replicas, err := ListReplicasInCluster()
			if err != nil {
				errs.Accumulate(err)
				logf.Log.Info("ResourceEachCheck: failed to retrieve list of replicas")
			}
			logf.Log.Info("ResourceCheck:", "num replicas", len(replicas))
			if len(replicas) != 0 {
				errs.Accumulate(fmt.Errorf("gRPC: count of replicas reported via mayastor client is %d", len(replicas)))
			}
		}
		// check nvmeControllers
		{
			nvmeControllers, err := ListNvmeControllersInCluster()
			if err != nil {
				errs.Accumulate(err)
				logf.Log.Info("ResourceEachCheck: failed to retrieve list of nvme controllers")
			}
			logf.Log.Info("ResourceCheck:", "num nvme controllers", len(nvmeControllers))
			if len(nvmeControllers) != 0 {
				errs.Accumulate(fmt.Errorf("gRPC: count of replicas reported via mayastor client is %d", len(nvmeControllers)))
			}
		}
	} else {
		logf.Log.Info("WARNING: gRPC calls to mayastor are not enabled, all checks cannot be run")
	}
	return errs.GetError()
}

func ResourceK8sCheck() error {
	var errs = common.ErrorAccumulator{}

	pods, err := CheckForTestPods()
	errs.Accumulate(err)
	if pods {
		errs.Accumulate(fmt.Errorf("found Pods"))
	}

	pvcs, err := CheckForPVCs()
	errs.Accumulate(err)
	if pvcs {
		errs.Accumulate(fmt.Errorf("found PersistentVolumeClaims"))
	}

	pvs, err := CheckForPVs()
	errs.Accumulate(err)

	if pvs {
		errs.Accumulate(fmt.Errorf("found PersistentVolumes"))
	}

	return errs.GetError()
}
