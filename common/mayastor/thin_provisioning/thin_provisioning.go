package thin_provisioning

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	_ "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/k8stest"
	coreV1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	DefTimeoutSecs     = 600
	RebuildTimeoutSecs = 120
	DefPoolCommitment  = 2.2
)

func GetPools() []common.MayastorPool {
	pools, err := k8stest.ListMsPools()
	Expect(err).ToNot(HaveOccurred(), "failed to get mayastor pools")
	Expect(len(pools)).ToNot(BeZero())
	for _, pool := range pools {
		Expect(pool.Status.Used).To(BeZero())
	}
	// sorts the pools by capacity to make sure we can't get less than the smallest capacity.
	sort.Slice(pools, func(i, j int) bool {
		return pools[i].Status.Capacity < pools[j].Status.Capacity
	})
	return pools
}

func CreateSc(provisioningType common.ProvisioningType, replicaCount int) (string, error) {
	scName := strings.ToLower(fmt.Sprintf("%s-sc-%d-%s-%s-%s", provisioningType, replicaCount, common.ShareProtoNvmf, common.VolRawBlock, randomString(5)))
	err := k8stest.NewScBuilder().
		WithName(scName).
		WithNamespace(common.NSDefault).
		WithProtocol(common.ShareProtoNvmf).
		WithReplicas(replicaCount).
		WithVolumeBindingMode(storagev1.VolumeBindingImmediate).
		WithProvisioningType(provisioningType).
		BuildAndCreate()
	if err != nil {
		return scName, err
	}
	return scName, nil
}

func randomString(length int) string {
	alphabet := "abcdefghijklmnopqrstuvwxyz"
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < length; i++ {
		r := randomInt(1, 10)
		if r%2 == 0 {
			c := alphabet[randomInt(0, k-1)]
			sb.WriteByte(c)
		} else {
			sb.Write([]byte(strconv.Itoa(randomInt(0, 9))))
		}

	}
	return sb.String()
}

func randomInt(min, max int) int {
	rand.NewSource(time.Now().UnixNano())
	return min + rand.Intn(max-min+1)
}

func CreateVolume(scName, volBaseName string, volSizeMb int) (string, string, error) {
	// Create the volume
	volName := strings.ToLower(fmt.Sprintf("%s-vol-%d-%s-%s-%s", volBaseName, common.DefaultReplicaCount(), common.ShareProtoNvmf, common.VolRawBlock, randomString(5)))
	uid, err := k8stest.MkPVC(volSizeMb, volName, scName, common.VolRawBlock, common.NSDefault)
	if err != nil {
		return uid, volName, err
	}
	logf.Log.Info("Volume", "uid", uid, "name", volName)
	return uid, volName, nil
}

func CreateAndRunSizedFio(uuid string, sizeMiB int, volName string, expectError bool) {
	fioPodName := fmt.Sprintf("fio-%s", volName)
	pod := k8stest.CreateFioPodDef(fioPodName, volName, common.VolRawBlock, common.NSDefault)
	Expect(pod).ToNot(BeNil())

	var args = []string{
		"--",
	}
	args = append(args, fmt.Sprintf("--filename=%s", common.FioBlockFilename))
	args = append(args, fmt.Sprintf("--size=%dM", sizeMiB))
	args = append(args, common.GetFioArgs()...)
	logf.Log.Info("fio", "arguments", args)

	pod.Spec.Containers[0].Args = args
	pod, err := k8stest.CreatePod(pod, common.NSDefault)
	Expect(err).ToNot(HaveOccurred())
	Expect(pod).ToNot(BeNil())

	// Wait for the fio Pod to transition to running
	Eventually(func() bool {
		return k8stest.IsPodRunning(fioPodName, common.NSDefault)
	},
		DefTimeoutSecs,
		"1s",
	).Should(BeTrue())
	logf.Log.Info("fio test pod is running.")

	if !expectError {
		msvc_err := k8stest.MsvConsistencyCheck(uuid)
		Expect(msvc_err).ToNot(HaveOccurred(), "%v", msvc_err)
	}

	logf.Log.Info("Waiting for run to complete", "timeout", DefTimeoutSecs)
	tSecs := 0
	var phase coreV1.PodPhase
	var podLogSynopsis *common.E2eFioPodLogSynopsis
	for {
		if tSecs > DefTimeoutSecs {
			break
		}
		time.Sleep(1 * time.Second)
		tSecs += 1
		phase, podLogSynopsis, err = k8stest.CheckFioPodCompleted(fioPodName, common.NSDefault)
		Expect(err).To(BeNil(), "CheckPodComplete got error %s", err)
		if phase != coreV1.PodRunning {
			break
		}
	}
	if expectError {
		Expect(podLogSynopsis.CriticalFailure).To(BeFalse(), "%s", podLogSynopsis)
		Expect(phase == coreV1.PodFailed).To(BeTrue(), "fio pod phase is %s", phase)
		logf.Log.Info("fio failed", "duration", tSecs)

	} else {
		Expect(phase == coreV1.PodSucceeded).To(BeTrue(), "fio pod phase is %s, %s", phase, podLogSynopsis)
		logf.Log.Info("fio completed", "duration", tSecs)
	}
	err = k8stest.DeletePod(fioPodName, common.NSDefault)
	Expect(err).ToNot(HaveOccurred())
}

func CreateAndRunRunningFio(sizeMiB int, volName string) *coreV1.Pod {
	fioPodName := fmt.Sprintf("fio-%s-%s", volName, randomString(5))
	pod := k8stest.CreateFioPodDef(fioPodName, volName, common.VolRawBlock, common.NSDefault)
	Expect(pod).ToNot(BeNil())

	var args = []string{
		"--",
	}
	args = append(args, fmt.Sprintf("--filename=%s", common.FioBlockFilename))
	args = append(args, fmt.Sprintf("--size=%dM", sizeMiB))
	args = append(args, fmt.Sprintf("--runtime=%ds", DefTimeoutSecs))
	//args = append(args, "-d")
	args = append(args, common.GetFioArgs()...)
	logf.Log.Info("fio", "arguments", args)

	pod.Spec.Containers[0].Args = args
	pod, err := k8stest.CreatePod(pod, common.NSDefault)
	Expect(err).ToNot(HaveOccurred())
	Expect(pod).ToNot(BeNil())

	// Wait for the fio Pod to transition to running
	Eventually(func() bool {
		return k8stest.IsPodRunning(fioPodName, common.NSDefault)
	},
		DefTimeoutSecs,
		"1s",
	).Should(Equal(true))
	logf.Log.Info("fio test pod is running.")
	return pod
}

func CleanUp(scName string, volNames []string) error {
	for _, name := range volNames {
		err := k8stest.RmPVC(name, scName, common.NSDefault)
		if err != nil {
			return err
		}
	}
	err := k8stest.RmStorageClass(scName)
	if err != nil {
		return err
	}
	_, err = k8stest.CheckForPVs()
	if err != nil {
		return err
	}
	return nil
}

// returns the percentage of the disk pool in the selected units (1.00 = 100%)
func GetPoolSizeFraction(pool common.MayastorPool, percentCapacity float64, unit string) int {
	capacityInUnits := GetSizePerUnits(pool.Status.Capacity, unit)
	return int(capacityInUnits * percentCapacity)
}

func GetSizePerUnits(b uint64, unit string) float64 {
	m := map[string]float64{"": 0, "KiB": 1, "MiB": 2, "GiB": 3, "TiB": 4, "PiB": 5}
	bf := float64(b)
	bf /= math.Pow(1024, m[unit])
	return math.Round(bf)
}

func GetReplicasUuid(replicaTopology common.ReplicaTopology) []string {
	var uuids = make([]string, 0, len(replicaTopology))
	for uuid := range replicaTopology {
		uuids = append(uuids, uuid)
	}
	return uuids
}
