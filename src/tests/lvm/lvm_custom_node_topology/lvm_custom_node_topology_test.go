package lvm_custom_node_topology

import (
	"fmt"
	"testing"
	"time"

	"github.com/openebs/openebs-e2e/common/e2e_config"
	"github.com/openebs/openebs-e2e/common/e2e_ginkgo"
	"github.com/openebs/openebs-e2e/common/lvm"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/k8stest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	coreV1 "k8s.io/api/core/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var nodeConfig lvm.LvmNodesDevicePvVgConfig
var appInstances []*k8stest.FioApplication
var targetNode, key string
var csiNodeUpdateTime = 60 // in seconds

var defDaemonsetReadyTime = 120 // in seconds

/*
Background:
	Given that the product is installed on a kubernetes cluster
Scenario: Node custom node topology for immediate volume binding
	Given One of the worker node labeled with custom label
	And Minimum two worker nodes should exist in cluster
	When Lvm immediate binding volumes and applications (number of worker nodes + 1) are deployed using custom topology
	Then All volumes should get provisioned on only those node which was labeled prior to the provisioning
*/

func customTopologyImmediateTest(decor string, engine common.OpenEbsEngine, volType common.VolumeType, fstype common.FileSystemType, volBindModeWait bool) {
	workerNodes, err := lvm.ListLvmNode(common.NSOpenEBS())
	Expect(err).ToNot(HaveOccurred(), "failed to list worker node")

	// minimum worker nodes in cluster should be two
	Expect(len(workerNodes)).Should(BeNumerically(">=", 1),
		"test case requires are least 2 worker nodes, %d nodes found", len(workerNodes))

	key = "lvme2e/nodename"
	targetNode = workerNodes[0]

	// label worker node
	err = k8stest.LabelNode(targetNode, key, targetNode)
	Expect(err).ToNot(HaveOccurred(), "failed to label node %s", targetNode)

	lvmScTopology := k8stest.LvmOptions{
		VolGroup:      nodeConfig.VgName,
		Storage:       "lvm",
		ThinProvision: common.No,
		AllowedTopologies: []coreV1.TopologySelectorTerm{
			{
				MatchLabelExpressions: []coreV1.TopologySelectorLabelRequirement{
					{
						Key:    key,
						Values: []string{targetNode},
					},
				},
			},
		},
	}
	appInstances = []*k8stest.FioApplication{}

	for i := 0; i <= len(workerNodes); i++ {
		app := k8stest.FioApplication{
			Decor:                   fmt.Sprintf("%s-%d", decor, i),
			VolSizeMb:               1024,
			OpenEbsEngine:           engine,
			VolType:                 volType,
			FsType:                  fstype,
			Loops:                   1,
			VolWaitForFirstConsumer: volBindModeWait,
			Lvm:                     lvmScTopology,
		}
		appInstances = append(appInstances, &app)
	}

	// deploy all fio application
	for ix, app := range appInstances {
		logf.Log.Info(fmt.Sprintf("%d)", ix), "fio-pod", app.Decor, "err", err)
		err = app.DeployApplication()
		Expect(err).ToNot(HaveOccurred(), "failed to deploy %s, %v", app.Decor, err)
	}

	// verify all fio application are deployed on same labeled node
	for _, app := range appInstances {
		appPodName := app.GetPodName()
		logf.Log.Info("app pod", "fio-pod", appPodName)
		//get node name where where app is deployed
		node, err := k8stest.GetNodeForPodByPrefix(appPodName, common.NSDefault)
		Expect(err).ToNot(HaveOccurred(), "failed to node for app pod %s, %v", appPodName, err)
		Expect(node).Should(Equal(targetNode), "app pod %s does not scheduled on node %s", targetNode)
	}

	// remove all fio application
	for ix, app := range appInstances {
		logf.Log.Info(fmt.Sprintf("%d)", ix), "fio-pod", app.Decor, "err", err)
		// remove app pod, pvc,sc
		err = app.Cleanup()
		Expect(err).To(BeNil(), "failed to clean resources")
	}

	// Remove the labels from nodes after the end of test
	err = k8stest.UnlabelNode(targetNode, key)
	Expect(err).ToNot(HaveOccurred(), "failed to remove label from node %s", targetNode)
	targetNode = ""

	productConfig := e2e_config.GetConfig().Product
	label := fmt.Sprintf("%s=%s", productConfig.LocalEngineComponentPodLabelKey,
		productConfig.LvmEngineComponentDsPodLabelValue)

	// restart lvm daemonset pods so that topology key present in csinode kubernetes object
	// for local.csi.openebs.io plugin driver should be removed before starting new topology test
	// and to do so , daemonset pods need to restarted after removing node label with the key
	err = k8stest.DeletePodsByLabel(label, common.NSOpenEBS())
	Expect(err).To(BeNil(), "failed to restart lvm daemonset pods with label %s", label)

	// verify lvm daemonset to be ready
	Eventually(func() bool {
		return k8stest.DaemonSetReady(productConfig.LvmEngineDaemonSetName, common.NSOpenEBS())
	},
		defDaemonsetReadyTime,
		"5s",
	).Should(BeTrue())

	ready, err := k8stest.OpenEBSReady(10, 540)
	Expect(err).To(BeNil(), "failed to verify openebs pods running state")
	Expect(ready).To(BeTrue(), "some of the openebs pods are not running")

	// verify topology key in csi node
	var csiNodeErr error
	Eventually(func() bool {
		var isKeyFound bool
		isKeyFound, csiNodeErr = k8stest.CheckCsiNodeTopologyKeysPresent(workerNodes[0],
			productConfig.LvmEnginePluginDriverName,
			[]string{
				key,
			})
		if csiNodeErr != nil {
			logf.Log.Info("Failed to check csinode topology key",
				"driver", productConfig.LvmEnginePluginDriverName,
				"key", key,
				"node", workerNodes[0],
				"error", err)
		}
		return isKeyFound
	},
		csiNodeUpdateTime,
		"5s",
	).Should(BeFalse())

	Expect(csiNodeErr).ToNot(HaveOccurred(), "failed to get csi node %s, %v", workerNodes[0], csiNodeErr)
}

/*
Background:
	Given that the product is installed on a kubernetes cluster
Scenario: Node custom node topology for immediate volume binding
	Given One of the worker node labeled with custom label
	And Minimum two worker nodes should exist in cluster
	When Lvm WaitForFirstConsumer binding volumes and applications (number of worker nodes + 1) are deployed using custom topology
	Then All volumes should be in pending state
	When lvm node-daemonset pods are restarted
	Then Verify topology key is now available in csi node for local.csi.openebs.io plugin driver
	And All volumes should be transition from pending to bound state
	And All volumes should get provisioned on only those node which was labeled prior to the provisioning
*/

func customTopologyWfcTest(decor string, engine common.OpenEbsEngine, volType common.VolumeType, fstype common.FileSystemType, volBindModeWait bool) {

	workerNodes, err := lvm.ListLvmNode(common.NSOpenEBS())
	Expect(err).ToNot(HaveOccurred(), "failed to list worker node")

	// minimum worker nodes in cluster should be two
	Expect(len(workerNodes)).Should(BeNumerically(">=", 1),
		"test case requires are least 2 worker nodes, %d nodes found", len(workerNodes))

	key = "lvme2e/nodename"
	targetNode = workerNodes[0]

	// label worker node
	logf.Log.Info("Label node", "node", targetNode, "key", key)
	err = k8stest.LabelNode(targetNode, key, targetNode)
	Expect(err).ToNot(HaveOccurred(), "failed to label node %s", targetNode)

	lvmScOption := k8stest.LvmOptions{
		VolGroup:      nodeConfig.VgName,
		Storage:       "lvm",
		ThinProvision: common.No,
		AllowedTopologies: []coreV1.TopologySelectorTerm{
			{
				MatchLabelExpressions: []coreV1.TopologySelectorLabelRequirement{
					{
						Key:    key,
						Values: []string{targetNode},
					},
				},
			},
		},
	}
	appInstances = []*k8stest.FioApplication{}

	for i := 0; i <= len(workerNodes); i++ {
		app := k8stest.FioApplication{
			Decor:                          fmt.Sprintf("%s-%d", decor, i),
			VolSizeMb:                      1024,
			OpenEbsEngine:                  engine,
			VolType:                        volType,
			FsType:                         fstype,
			Loops:                          1,
			VolWaitForFirstConsumer:        volBindModeWait,
			Lvm:                            lvmScOption,
			SkipPvcVerificationAfterCreate: true,
		}
		appInstances = append(appInstances, &app)
	}

	// deploy all lvm volume
	for ix, app := range appInstances {
		logf.Log.Info(fmt.Sprintf("%d)", ix), "lvm-volume", app.Decor)
		err = app.CreateVolume()
		Expect(err).ToNot(HaveOccurred(), "failed to create volume %s, %v", app.Decor, err)
	}

	logf.Log.Info("Sleep for 30 seconds before verifying pvc's pending state")
	time.Sleep(30 * time.Second)

	for ix, app := range appInstances {
		logf.Log.Info(fmt.Sprintf("%d)", ix), "lvm-volume", app.Decor)
		// verify pvc to be in pending state
		pvcPhase, err := k8stest.GetPvcStatusPhase(app.GetPvcName(), common.NSDefault)
		Expect(err).ToNot(HaveOccurred(), "failed to get pvc phase")
		Expect(pvcPhase).Should(Equal(coreV1.ClaimPending), "pvc phase is not pending")
	}

	productConfig := e2e_config.GetConfig().Product
	label := fmt.Sprintf("%s=%s", productConfig.LocalEngineComponentPodLabelKey,
		productConfig.LvmEngineComponentDsPodLabelValue)

	// Restart lvm daemonset pods after applying node label with key
	// so that csinode kubernetes object for local.csi.openebs.io plugin driver picks
	// that particular topology key for scheduling volume
	err = k8stest.DeletePodsByLabel(label, common.NSOpenEBS())
	Expect(err).To(BeNil(), "failed to restart lvm daemonset pods with label %s", label)

	// verify lvm daemonset to be ready
	Eventually(func() bool {
		return k8stest.DaemonSetReady(productConfig.LvmEngineDaemonSetName, common.NSOpenEBS())
	},
		defDaemonsetReadyTime,
		"5s",
	).Should(BeTrue())

	ready, err := k8stest.OpenEBSReady(10, 540)
	Expect(err).To(BeNil(), "failed to verify openebs pods running state")
	Expect(ready).To(BeTrue(), "some of the openebs pods are not running")

	// verify topology key in csi node
	logf.Log.Info("verify topology key in csi node", "key", key, "node", targetNode)
	var csiNodeErr error
	Eventually(func() bool {
		var isKeyFound bool
		isKeyFound, csiNodeErr = k8stest.CheckCsiNodeTopologyKeysPresent(targetNode,
			productConfig.LvmEnginePluginDriverName,
			[]string{
				key,
			})
		if csiNodeErr != nil {
			logf.Log.Info("Failed to check csinode topology key",
				"driver", productConfig.LvmEnginePluginDriverName,
				"key", key,
				"node", targetNode,
				"error", err)
		}
		return isKeyFound
	},
		csiNodeUpdateTime,
		"5s",
	).Should(BeTrue())

	Expect(csiNodeErr).ToNot(HaveOccurred(), "failed to get csi node %s, %v", targetNode, err)

	// deploy fio pods for created lvm volumes
	for ix, app := range appInstances {
		logf.Log.Info(fmt.Sprintf("%d)", ix), "fio-pod", app.Decor)

		// deploy fio pod with created volume
		logf.Log.Info("deploy fio pod with created volume")
		err = app.DeployApplication()
		Expect(err).To(BeNil(), "failed to deploy app")

		//verify pvc and pv to be bound
		volUuid, err := k8stest.VerifyVolumeProvision(app.GetPvcName(), common.NSDefault)
		Expect(err).ToNot(HaveOccurred())
		Expect(volUuid).ToNot(BeEmpty())

		// use created PVC which is deployed as part of restore app
		err = app.RefreshVolumeState()
		Expect(err).ToNot(HaveOccurred())
	}

	// verify all fio application are deployed on same labeled node
	logf.Log.Info("verify all fio application are deployed on same labeled node", "labeled node", targetNode)
	for ix, app := range appInstances {
		logf.Log.Info(fmt.Sprintf("%d)", ix), "fio-pod", app.Decor)
		appPodName := app.GetPodName()
		logf.Log.Info("app pod", "fio-pod", appPodName)
		//get node name where where app is deployed
		node, err := k8stest.GetNodeForPodByPrefix(appPodName, common.NSDefault)
		Expect(err).ToNot(HaveOccurred(), "failed to get node for app pod %s, %v", appPodName, err)
		Expect(node).Should(Equal(targetNode), "app pod %s is not scheduled on node %s", targetNode)
	}

	// remove all fio application
	logf.Log.Info("remove all fio application")
	for ix, app := range appInstances {
		logf.Log.Info(fmt.Sprintf("%d)", ix), "fio-pod", app.Decor, "err", err)
		// remove app pod, pvc,sc
		err = app.Cleanup()
		Expect(err).To(BeNil(), "failed to clean resources")
	}

	// Remove the labels from nodes after the end of test
	err = k8stest.UnlabelNode(targetNode, key)
	Expect(err).ToNot(HaveOccurred(), "failed to remove label from node %s", targetNode)

	// restart lvm daemonset pods so that topology key present in csinode kubernetes object
	// for local.csi.openebs.io plugin driver should be removed before starting new topology test
	// and to do so , daemonset pods needs to restated after removing node label with the key
	err = k8stest.DeletePodsByLabel(label, common.NSOpenEBS())
	Expect(err).To(BeNil(), "failed to restart lvm daemonset pods with label %s", label)

	// verify lvm daemonset to be ready
	Eventually(func() bool {
		return k8stest.DaemonSetReady(productConfig.LvmEngineDaemonSetName, common.NSOpenEBS())
	},
		defDaemonsetReadyTime,
		"5s",
	).Should(BeTrue())

	ready, err = k8stest.OpenEBSReady(10, 540)
	Expect(err).To(BeNil(), "failed to verify openebs pods running state")
	Expect(ready).To(BeTrue(), "some of the openebs pods are not running")
}

func TestLvmCustomTopologyTest(t *testing.T) {
	// Initialise test and set class and file names for reports
	e2e_ginkgo.InitTesting(t, "lvm_custom_node_topology", "lvm_custom_node_topology")
}

var _ = Describe("lvm_custom_node_topology", func() {

	BeforeEach(func() {
		// Check ready to run
		err := e2e_ginkgo.BeforeEachK8sCheck()
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		// Check resource leakage.
		afterErr := e2e_ginkgo.AfterEachK8sCheck()
		// remove all fio application
		for ix, app := range appInstances {
			logf.Log.Info(fmt.Sprintf("%d)", ix), "fio-pod", app.Decor)
			// remove app pod, pvc,sc
			err := app.Cleanup()
			Expect(err).To(BeNil(), "failed to clean resources")
		}
		if targetNode != "" {
			// Remove the labels from nodes after the end of test
			err := k8stest.UnlabelNode(targetNode, key)
			Expect(err).ToNot(HaveOccurred(), "failed to remove label from node %s", targetNode)
			targetNode = ""
		}
		Expect(afterErr).ToNot(HaveOccurred())
	})

	// immediate binding
	It("lvm: should verify custom node topology for immediate binding mode", func() {
		customTopologyImmediateTest("lvm", common.Lvm, common.VolFileSystem, common.Ext4FsType, false)
	})

	// wait for first consumer
	It("lvm: should verify custom node topology for wait for first consumer binding mode", func() {
		customTopologyWfcTest("lvm", common.Lvm, common.VolFileSystem, common.Ext4FsType, true)
	})

})

var _ = BeforeSuite(func() {
	err := e2e_ginkgo.SetupTestEnv()
	Expect(err).ToNot(HaveOccurred(), "failed to setup test environment in BeforeSuite : SetupTestEnv %v", err)

	//setup nodes with lvm pv and vg
	nodeConfig, err = lvm.SetupLvmNodes("lvmvg", 10737418240)
	Expect(err).ToNot(HaveOccurred(), "failed to setup lvm pv and vg")
})

var _ = AfterSuite(func() {
	// NB This only tears down the local structures for talking to the cluster,
	// not the kubernetes cluster itself.	By("tearing down the test environment")
	logf.Log.Info("remove node with device and vg", "node config", nodeConfig)
	err := nodeConfig.RemoveConfiguredLvmNodesWithDeviceAndVg()
	Expect(err).ToNot(HaveOccurred(), "failed to cleanup node with device")

	err = k8stest.TeardownTestEnv()
	Expect(err).ToNot(HaveOccurred(), "failed to tear down test environment in AfterSuite : TeardownTestEnv %v", err)
})
