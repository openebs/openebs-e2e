package k8stest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/controlplane"
	agent "github.com/openebs/openebs-e2e/common/e2e_agent"
	"github.com/openebs/openebs-e2e/common/e2e_config"
	"github.com/openebs/openebs-e2e/common/k8s_portforward"
	"github.com/openebs/openebs-e2e/common/mayastorclient"

	"k8s.io/apimachinery/pkg/util/uuid"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// var productName = e2e_config.GetConfig().Product.ProductName
var nvmeControllerModel = e2e_config.GetConfig().Product.NvmeControllerModel

// use the e2e-agent running on each non-nexus node:
//
//	for each non-nexus replica node
//	    nvme connect to its own target
//	    cksum /dev/nvme0n1p2
//	    nvme disconnect
//	compare the checksum results, they should match
func CompareReplicas(nexusIP string, replica_1_IP string, replica_2_IP string, volUUID string) error {

	// get nexus node nqn
	nexusNodeNqn, err := GetNodeNqn(nexusIP)
	logf.Log.Info("node", "ip", nexusIP, "nqn", nexusNodeNqn)
	if err != nil {
		return fmt.Errorf("failed to get valid nexus node nqn, error :%s", err.Error())
	}
	if nexusNodeNqn == "" {
		return fmt.Errorf("failed to get valid nexus node nqn")
	}

	// the first replica
	replicas, err := mayastorclient.ListReplicas([]string{replica_1_IP})
	if err != nil {
		return fmt.Errorf("failed to list replicas, error :%s", err.Error())
	}
	if len(replicas) != 1 {
		return fmt.Errorf("expected to find 1 replica, found %d", len(replicas))
	}
	uuid := replicas[0].GetUuid()
	uri := replicas[0].GetUri()
	logf.Log.Info("First Replica", "IP", replica_1_IP, "uuid", uuid, "Uri", uri)

	// get replica state
	replicaState, err := GetReplicaState(volUUID, uuid)
	if err != nil {
		return fmt.Errorf("failed to get replica state, error :%s", err.Error())
	}
	if replicaState != controlplane.ReplicaStateOnline() {
		return fmt.Errorf("replica %s is not online", uuid)
	}

	firstchecksum, err := ChecksumReplica(nexusIP, replica_1_IP, uri, 60, nexusNodeNqn)
	if err != nil {
		return fmt.Errorf("failed to checksum replica, error :%s", err.Error())
	}

	// the second replica
	replicas, err = mayastorclient.ListReplicas([]string{replica_2_IP})
	if err != nil {
		return fmt.Errorf("failed to list replicas, error :%s", err.Error())
	}
	if len(replicas) != 1 {
		return fmt.Errorf("expected to find 1 replica, found %d", len(replicas))
	}
	uuid = replicas[0].GetUuid()
	uri = replicas[0].GetUri()
	logf.Log.Info("Second Replica", "IP", replica_2_IP, "uuid", uuid, "Uri", uri)

	// get replica state
	replicaState, err = GetReplicaState(volUUID, uuid)
	if err != nil {
		return fmt.Errorf("failed to get replica state, error :%s", err.Error())
	}
	if replicaState != controlplane.ReplicaStateOnline() {
		return fmt.Errorf("replica %s is not online", uuid)
	}

	secondchecksum, err := ChecksumReplica(nexusIP, replica_2_IP, uri, 60, nexusNodeNqn)
	if err != nil {
		return fmt.Errorf("failed to checksum replica, error :%s", err.Error())
	}

	// verify that they match
	logf.Log.Info("match", "first", firstchecksum, "this", secondchecksum)
	if firstchecksum != secondchecksum {
		return fmt.Errorf("checksums do not match %s %s", firstchecksum, secondchecksum)
	}
	return nil
}

// for json deserialization of output from nvme list-subsys
type nvmeListSubsystemEntry struct {
	Name  string
	NQN   string
	Paths []map[string]string
}

func (n *nvmeListSubsystemEntry) String() string {
	return fmt.Sprintf("Name:%s, NQN:%s, Paths:%v", n.Name, n.NQN, n.Paths)
}

// for json deserialization of output from nvme list
type nvmeListDeviceEntry struct {
	DevicePath   string
	SerialNumber string
	ModelNumber  string
}

func (n *nvmeListDeviceEntry) String() string {
	return fmt.Sprintf("DevicePath:%s, ModelNumber:%s SerialNumber:%s", n.DevicePath, n.ModelNumber, n.SerialNumber)
}

// for json deserialization of output from nvme resv report
type regctlextEntry struct {
	Hostid string
}

func (r *regctlextEntry) String() string {
	return fmt.Sprintf("hostid: %s", r.Hostid)
}

// for json deserialization of output from nvme resv report
type reservationEntry struct {
	Regctlext []regctlextEntry
}

func (r *reservationEntry) String() string {
	s := "regctlext: ["
	for _, re := range r.Regctlext {
		s += re.String() + ", "
	}
	return s + "]"
}

// trimForJson simplistic removal of extraneous characters surrounding JSON formatted data,
// good enough for now.
func trimForJson(jstr string) string {
	var start, end int
	var endChar uint8
	for start = 0; start < len(jstr); start++ {
		if jstr[start] == '{' {
			endChar = '}'
			break
		}
		if jstr[start] == '[' {
			endChar = ']'
			break
		}
	}
	for end = len(jstr) - 1; end > 0; end-- {
		if jstr[end] == endChar {
			break
		}
	}
	return jstr[start : end+1]
}

func getNvmeDevice(initiatorIP string, maxRetries int, targetNqn string) (bool, string, string) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			logf.Log.Info("PANIC!", "error", panicErr)
		}
	}()
	logf.Log.Info("Check device path after running nvme connect with host id to replica")
	// fetch device path again because it might change between above multiple nvme disconnect and connect
	retryCount := 0
	found := false
	var deviceOnly, devicePath, hostId string

	//	for retryCount < maxRetries && !found {
	for retryCount < 1 && !found {
		retryCount++

		output, err := agent.NvmeList(initiatorIP)
		if output == "" || err != nil {
			logf.Log.Info("nvme list failed", "output", output, "err", err)
			continue
		}
		output = trimForJson(output)
		list := make(map[string][]nvmeListDeviceEntry)
		if err = json.Unmarshal([]byte(output), &list); err != nil {
			logf.Log.Info("Failed to unmarshal target", "output", output)
			continue
		}
		logf.Log.Info("nvme", "list", list)

		nqnPathNames := make(map[string][]string)
		output, err = agent.NvmeListSubSys(initiatorIP)
		if output == "" || err != nil {
			logf.Log.Info("nvme list-subsys failed", "output", output, "err", err)
			continue
		} else {
			subSysList := make(map[string][]nvmeListSubsystemEntry)
			output = trimForJson(output)
			if err = json.Unmarshal([]byte(output), &subSysList); err != nil {
				logf.Log.Info("Failed to unmarshal target", "error", err)
			} else {
				logf.Log.Info("nvme", "list-subsys", subSysList)
				for _, subsys := range subSysList[`Subsystems`] {
					var pathNames []string
					for _, path := range subsys.Paths {
						pathNames = append(pathNames, path[`Name`])
					}
					nqnPathNames[subsys.NQN] = pathNames
				}
			}
		}

		for nqnKey, pathNames := range nqnPathNames {
			if nqnKey == targetNqn {
				for _, pathName := range pathNames {
					for _, deviceData := range list["Devices"] {
						devicePath = deviceData.DevicePath
						// For mayastor the ModelNumber would be Mayastor NVMe controller
						// For openebspro the ModelNumber would be Mayastor-Pro NVMe controller
						// check if the device listed is mayastor/openebspro nvme controlled device
						if deviceData.ModelNumber != nvmeControllerModel {
							continue
						}
						if len(devicePath) > 5 {
							deviceOnly = devicePath[5:]
							if strings.HasPrefix(deviceOnly, pathName) {
								logf.Log.Info("Found device", "devicePath", devicePath, "pathName", pathName, "nqn", nqnKey)
								// get reservation report
								reservations, err := agent.ListReservation(initiatorIP, devicePath)
								if reservations == "" || err != nil {
									logf.Log.Info("nvme list reservation failed", "reservation", reservations, "error", err)
									continue
								}
								// sample nvme reservation output in json format
								/*
									{
										"gen" : 1,
										"rtype" : 2,
										"regctl" : 1,
										"ptpls" : 1,
										"regctlext" : [
											{
											"cntlid" : 61135,
											"rcsts" : 1,
											"rkey" : 13212321943525734355,
											"hostid" : "a12aae512c0e4b918c362d9fe529a6f0"
											}
										]
									}
								*/
								var reservationList reservationEntry
								if err = json.Unmarshal([]byte(reservations), &reservationList); err != nil {
									logf.Log.Info("Failed to unmarshal target", "reservations", reservations)
									continue
								}
								logf.Log.Info("reservation", "regctlext", reservationList)
								reservation := reservationList.Regctlext
								if len(reservation) > 0 {
									hostId = reservation[0].Hostid
								}
								return true, devicePath, hostId
							}
						}
					}
				}
			}
		}
		time.Sleep(5 * time.Second)
	}
	return false, "", ""
}

// ChecksumReplica checksums the nvme target defined by the given uri
// It uses the e2e agent and the nvme client to connect to the target.
// the returned format is derived from cksum: <checksum> <size>
// e.g. "924018992 61849088
func ChecksumReplica(initiatorIP, targetIP, nqn string, maxRetries int, nexusNodeNqn string) (string, error) {
	var err error
	logf.Log.Info("ChecksumReplica", "nexusIP", initiatorIP, "nodeIP", targetIP, "nqn", nqn)
	resp, err := agent.NvmeConnect(initiatorIP, targetIP, nqn, nexusNodeNqn)
	resp = strings.TrimSpace(resp)
	if err != nil {
		logf.Log.Info("Running agent failed", "error", err)
		return "", err
	}
	if resp != "" { // connect should be silent
		return "", fmt.Errorf("nvme connect returned with %s", resp)
	}

	// NOTE:
	// Beyond this point we must ALWAYS disconnect the previously established connection
	// prior to returning

	// nvme list returns all nvme devices in form n * [ ... <Node> ... <Model> ... \n]
	// we want the device (=Node) associated with Model = "Mayastor NVMe controller"
	devicePath := ""
	hostId := ""
	found := false
	found, devicePath, hostId = getNvmeDevice(initiatorIP, 60, nqn)

	if !found {
		_, _ = agent.NvmeDisconnect(initiatorIP, nqn)
		return "", fmt.Errorf("failed to get device path")
	}

	if hostId != "" {
		//format the hostid into something nvme-cli can accept
		// example: 28940560b05b47b299cc4098943eb1e5 -> 28940560-b05b-47b2-99cc-4098943eb1e5
		formattedHostId := fmt.Sprintf("%s-%s-%s-%s-%s", hostId[0:8], hostId[8:12], hostId[12:16], hostId[16:20], hostId[20:])
		// disconnect nvme
		resp, err = agent.NvmeDisconnect(initiatorIP, nqn)
		if err != nil {
			logf.Log.Info("Running agent failed", "error", err)
			return "", err
		}
		logf.Log.Info("Executed NvmeDisconnect", "nqn", nqn, "got", resp)

		// nvme connect with Host uuid
		resp, err = agent.NvmeConnectWithHostId(initiatorIP, targetIP, nqn, nexusNodeNqn, formattedHostId)
		resp = strings.TrimSpace(resp)
		if err != nil {
			logf.Log.Info("Running agent failed", "error", err)
			return "", err
		}
		if resp != "" { // connect should be silent
			return "", fmt.Errorf("nvme connect with host id returned with %s", resp)
		}
		logf.Log.Info("Executed NvmeConnectWithHostId", "nqn", nqn, "nexus node nqn", nexusNodeNqn, "hostid", formattedHostId, "got", resp)

		logf.Log.Info("Check device path after running nvme connect with host id to replica")
		// fetch device path again because it might change between above multiple nvme disconnect and connect
		found, devicePath, _ = getNvmeDevice(initiatorIP, 60, nqn)
		if !found {
			_, _ = agent.NvmeDisconnect(initiatorIP, nqn)
			return "", fmt.Errorf("failed to get device path 2nd time")
		}
	}

	deviceOnly := devicePath[5:] // remove the /dev/ prefix
	logf.Log.Info("only device", "Device", deviceOnly)

	logf.Log.Info("To be Executed ChecksumDevice", "devicePath", devicePath)
	// checksum the device
	// the returned format is <checksum> <size> <device>
	// e.g. "924018992 61849088 /dev/nvme0n1p2"
	cksumText, err := agent.ChecksumDevice(initiatorIP, devicePath)
	if err != nil {
		logf.Log.Info("Running agent failed", "error", err)
		_, _ = agent.NvmeDisconnect(initiatorIP, nqn)
		return "", err
	}
	cksumText = strings.TrimSpace(cksumText)
	// double check the response contains the device name
	if !strings.Contains(cksumText, deviceOnly) {
		_, _ = agent.NvmeDisconnect(initiatorIP, nqn)
		return "", fmt.Errorf("Unexpected result from cksum %v", cksumText)
	}
	cksumDevice := ``
	// discard the device name from the checksum
	fields := strings.Fields(cksumText)
	if len(fields) == 3 {
		cksumText = strings.Join(fields[:2], " ")
		cksumDevice = fields[2]
	}
	logf.Log.Info("Executed ChecksumDevice", "devicePath", devicePath, "got", cksumText, "device", cksumDevice)
	resp, err = agent.NvmeDisconnect(initiatorIP, nqn)
	if err != nil {
		logf.Log.Info("Running agent failed", "error", err)
		return "", err
	}
	logf.Log.Info("Executed NvmeDisconnect", "nqn", nqn, "got", resp)

	// check that the device no longer exists
	resp, err = agent.ListDevice(initiatorIP)
	if err != nil {
		logf.Log.Info("Running agent failed", "error", err)
		return "", err
	}
	//	logf.Log.Info("agent.ListDevice", "device only", deviceOnly, "resp", resp)
	if strings.Contains(resp, deviceOnly) {
		return "", fmt.Errorf("Device %s still exists", deviceOnly)
	}
	return cksumText, nil
}

// FsConsistentReplica verifies the filesystem consistency of the nvme target
// defined by the given uri.
// It uses the e2e agent and the nvme client to connect to the target.
// Returns the fsck or equivalent output and exit status code
func FsConsistentReplica(initiatorIP, targetIP, nqn string, maxRetries int, nexusNodeNqn string, fsType common.FileSystemType) (string, error) {
	var err error
	logf.Log.Info("FsConsistentReplica", "nexusIP", initiatorIP, "nodeIP", targetIP, "nqn", nqn)
	resp, err := agent.NvmeConnect(initiatorIP, targetIP, nqn, nexusNodeNqn)
	resp = strings.TrimSpace(resp)
	if err != nil {
		logf.Log.Info("Running agent failed", "error", err)
		return "", err
	}
	if resp != "" { // connect should be silent
		return "", fmt.Errorf("nvme connect returned with %s", resp)
	}

	// NOTE:
	// Beyond this point we must ALWAYS disconnect the previously established connection
	// prior to returning

	// nvme list returns all nvme devices in form n * [ ... <Node> ... <Model> ... \n]
	// we want the device (=Node) associated with Model = "Mayastor NVMe controller"
	devicePath := ""
	hostId := ""
	found := false
	found, devicePath, hostId = getNvmeDevice(initiatorIP, 60, nqn)

	if !found {
		_, _ = agent.NvmeDisconnect(initiatorIP, nqn)
		return "", fmt.Errorf("failed to get device path")
	}

	if hostId != "" {
		//format the hostid into something nvme-cli can accept
		// example: 28940560b05b47b299cc4098943eb1e5 -> 28940560-b05b-47b2-99cc-4098943eb1e5
		formattedHostId := fmt.Sprintf("%s-%s-%s-%s-%s", hostId[0:8], hostId[8:12], hostId[12:16], hostId[16:20], hostId[20:])
		// disconnect nvme
		resp, err = agent.NvmeDisconnect(initiatorIP, nqn)
		if err != nil {
			logf.Log.Info("Running agent failed", "error", err)
			return "", err
		}
		logf.Log.Info("Executed NvmeDisconnect", "nqn", nqn, "got", resp)

		// nvme connect with Host uuid
		resp, err = agent.NvmeConnectWithHostId(initiatorIP, targetIP, nqn, nexusNodeNqn, formattedHostId)
		resp = strings.TrimSpace(resp)
		if err != nil {
			logf.Log.Info("Running agent failed", "error", err)
			return "", err
		}
		if resp != "" { // connect should be silent
			return "", fmt.Errorf("nvme connect with host id returned with %s", resp)
		}
		logf.Log.Info("Executed NvmeConnectWithHostId successfully", "nqn", nqn, "nexus node nqn", nexusNodeNqn, "hostid", formattedHostId)

		logf.Log.Info("Check device path after running nvme connect with host id to replica")
		// fetch device path again because it might change between above multiple nvme disconnect and connect
		found, devicePath, _ = getNvmeDevice(initiatorIP, 60, nqn)
		if !found {
			_, _ = agent.NvmeDisconnect(initiatorIP, nqn)
			return "", fmt.Errorf("failed to get device path 2nd time")
		}
	}

	deviceOnly := devicePath[5:] // remove the /dev/ prefix
	logf.Log.Info("only device", "Device", deviceOnly)

	logf.Log.Info("To be Executed FsConsistentReplica", "devicePath", devicePath)

	fsckText, err := agent.FsCheckDevice(initiatorIP, devicePath, fsType)
	if err != nil {
		logf.Log.Info("FsCheckDevice failed", "error", err)
		_, _ = agent.NvmeDisconnect(initiatorIP, nqn)
		return "", err
	}
	logf.Log.Info("fsck Output for Snapshot", "fsckText", fsckText)
	fsckText = strings.TrimSpace(fsckText)
	logf.Log.Info("Executed fsckDevice", "devicePath", devicePath, "got", fsckText)

	resp, err = agent.NvmeDisconnect(initiatorIP, nqn)
	if err != nil {
		logf.Log.Info("Running agent failed", "error", err)
		return "", err
	}
	logf.Log.Info("Executed NvmeDisconnect", "nqn", nqn, "got", resp)

	// check that the device no longer exists
	resp, err = agent.ListDevice(initiatorIP)
	if err != nil {
		logf.Log.Info("Running agent failed", "error", err)
		return "", err
	}
	//	logf.Log.Info("agent.ListDevice", "device only", deviceOnly, "resp", resp)
	if strings.Contains(resp, deviceOnly) {
		return "", fmt.Errorf("Device %s still exists", deviceOnly)
	}
	return fsckText, nil
}

// FsFreezeReplica freezes  the filesystem of the nvme target
// defined by the given uri.
// It uses the e2e agent and the nvme client to connect to the target.
// Returns the exit status code
func FsFreezeReplica(nodeName string) (string, error) {
	nodeAddress, err := GetNodeIPAddress(nodeName)
	//Expect(err).ToNot(HaveOccurred())
	if err != nil {
		logf.Log.Info("unable to get Node IP Address", "error", err)
		return "", err
	}
	logf.Log.Info("Current", "nodeAddress", *nodeAddress)
	output, err := agent.NvmeList(*nodeAddress)
	if output == "" || err != nil {
		logf.Log.Info("nvme list failed", "output", output, "err", err)
		return "", err
	}
	output = trimForJson(output)
	list := make(map[string][]nvmeListDeviceEntry)
	if err = json.Unmarshal([]byte(output), &list); err != nil {
		logf.Log.Info("Failed to unmarshal target", "output", output)
		return "", err
	}
	logf.Log.Info("nvme", "list", list)
	var nvmeDevice string
	for _, deviceData := range list["Devices"] {
		// For mayastor the ModelNumber would be Mayastor NVMe controller
		// For openebspro the ModelNumber would be Mayastor-Pro NVMe controller
		// check if the device listed is mayastor/openebspro nvme controlled device
		if deviceData.ModelNumber != nvmeControllerModel {
			continue
		}
		nvmeDevice = deviceData.DevicePath
	}
	logf.Log.Info("nvmeDevice is ", "nvmeDevice", nvmeDevice)
	if nvmeDevice == "" {
		logf.Log.Info("nvmeDevice is empty. Not valid")
		return "", fmt.Errorf("failed to get nvme device")
	}
	fsFreezeText, err := agent.FsFreezeDevice(*nodeAddress, nvmeDevice, common.Ext4FsType)
	if err != nil {
		logf.Log.Info("FsFreezeDevice failed", "error", err)
		fsFreezeText = fmt.Sprintf("%s; %s", *nodeAddress, nvmeDevice)
	}
	return fsFreezeText, err
}

// FsUnfreezeReplica unfreezes the filesystem of the nvme target
// defined by the given uri.
// It uses the e2e agent and the nvme client to connect to the target.
// Returns the exit status code
func FsUnfreezeReplica(nodeName string) (string, error) {
	nodeAddress, err := GetNodeIPAddress(nodeName)
	if err != nil {
		logf.Log.Info("unable to get Node IP Address", "error", err)
		return "", err
	}
	logf.Log.Info("Current", "nodeAddress", *nodeAddress)
	output, err := agent.NvmeList(*nodeAddress)
	if output == "" || err != nil {
		logf.Log.Info("nvme list failed", "output", output, "err", err)
		return "", err
	}
	output = trimForJson(output)
	list := make(map[string][]nvmeListDeviceEntry)
	if err = json.Unmarshal([]byte(output), &list); err != nil {
		logf.Log.Info("Failed to unmarshal target", "output", output)
		return "", err
	}
	logf.Log.Info("nvme", "list", list)
	var nvmeDevice string
	for _, deviceData := range list["Devices"] {
		// For mayastor the ModelNumber would be Mayastor NVMe controller
		// For openebspro the ModelNumber would be Mayastor-Pro NVMe controller
		// check if the device listed is mayastor/openebspro nvme controlled device
		if deviceData.ModelNumber != nvmeControllerModel {
			continue
		}
		nvmeDevice = deviceData.DevicePath
	}
	logf.Log.Info("nvmeDevice is ", "nvmeDevice", nvmeDevice)
	if nvmeDevice == "" {
		logf.Log.Info("nvmeDevice is empty. Not valid")
		return "", fmt.Errorf("failed to get nvme device")
	}
	fsFreezeText, err := agent.FsUnfreezeDevice(*nodeAddress, nvmeDevice)
	if err != nil {
		logf.Log.Info("FsUnfreezeReplica failed", "error", err)
		fsFreezeText = fmt.Sprintf("%s; %s", *nodeAddress, nvmeDevice)
	}
	return fsFreezeText, err
}

// ExcludeNexusReplica - ensure the volume has no nexus-local replica
// This depends on there being an unused mayastor instance available so
// e.g. a 2-replica volume needs at least a 3-node cluster
func ExcludeNexusReplica(nexusIP string, nexusUuid string, volUuid string) (bool, error) {
	// get the nexus local device
	var nxlist []string
	nxlist = append(nxlist, nexusIP)
	if !mayastorclient.CanConnect() {
		return false, fmt.Errorf("gRPC calls not enabled")
	}
	nexusList, err := mayastorclient.ListNexuses(nxlist)
	if err != nil {
		return false, fmt.Errorf("Failed to list nexuses, err=%v", err)
	}
	if len(nexusList) == 0 {
		return false, fmt.Errorf("Expected to find at least 1 nexus")
	}

	nxChild := ""
	for _, nx := range nexusList {
		if nx.GetUuid() == nexusUuid {
			for _, ch := range nx.GetChildren() {
				if strings.HasPrefix(ch.Uri(), "bdev:///") {
					if nxChild != "" {
						return false, fmt.Errorf("More than 1 nexus local replica found")
					}
					nxChild = ch.Uri()
				}
			}
			if nxChild == "" { // there is no local replica so we are done
				return false, nil
			}
			break
		}
	}
	if nxChild == "" {
		return false, fmt.Errorf("failed to find the nexus")
	}

	// fault the replica
	logf.Log.Info("Faulting local replica", "replica", nxChild)
	err = mayastorclient.FaultNexusChild(nexusIP, nexusUuid, nxChild)
	if err != nil {
		return false, fmt.Errorf("Failed to fault child, err=%v", err)
	}

	// wait for the replica to disappear from the msv
	const sleepTime = 10
	const timeOut = 240
	var found bool
	for ix := 0; ix < (timeOut-1)/sleepTime; ix++ {
		found = false
		replicas, err := GetMsvReplicas(volUuid)
		if err != nil {
			return false, fmt.Errorf("Failed to get replicas, err=%v", err)
		}
		for _, replica := range replicas {
			if strings.HasPrefix(nxChild, replica.Uuid) {
				found = true
				break
			}
		}
		if !found {
			break
		}
		time.Sleep(sleepTime * time.Second)
	}
	if found {
		return true, fmt.Errorf("timed out waiting for faulted replica to be removed")
	}
	// wait for the msv to become healthy - now rebuilt with a non-nexus replica
	state := ""
	for ix := 0; ix < (timeOut-1)/sleepTime; ix++ {
		state, err = GetMsvState(volUuid)
		if err != nil {
			return false, fmt.Errorf("Failed to get state, err=%v", err)
		}
		if state == controlplane.VolStateHealthy() {
			break
		}
		time.Sleep(sleepTime * time.Second)
	}
	if state != controlplane.VolStateHealthy() {
		return true, fmt.Errorf("timed out waiting for volume to become healthy")
	}
	return true, nil
}

// Identify the nexus IP address,
// the uri of the replica , fault replica
func FaultReplica(volumeUuid string, replicaUuid string) error {
	nodeList, err := GetIOEngineNodes()
	if err != nil {
		return fmt.Errorf("Failed to get mayastor nodes, error %v", err)
	}

	nexus, _ := GetMsvNodes(volumeUuid)
	if err != nil {
		return fmt.Errorf("Failed to find nexus, error %v", err)
	}

	// identify the nexus IP address
	nexusIP := ""
	for _, node := range nodeList {
		if node.NodeName == nexus {
			nexusIP = node.IPAddress
		}

	}

	if nexusIP == "" {
		return fmt.Errorf("Nexus IP not found")
	}
	nexusNodeIP := nexusIP

	var nxList []string
	nxList = append(nxList, nexusIP)

	nexusList, _ := mayastorclient.ListNexuses(nxList)

	if len(nexusList) == 0 {
		return fmt.Errorf("failed to find at least 1 nexus")
	}
	nx := nexusList[0]

	// identify the local replica to be faulted
	nxChildUri := ""
	for _, ch := range nx.GetChildren() {
		if strings.Contains(ch.Uri(), replicaUuid) {
			nxChildUri = ch.Uri()
			break
		}
	}

	if nxChildUri == "" {
		return fmt.Errorf("Could not find nexus replica")
	}
	msv, err := GetMSV(volumeUuid)
	if err != nil {
		return fmt.Errorf("failed to retrieve MSV for volume %s, error %v", volumeUuid, err)
	}

	nexusUuid := msv.State.Target.Uuid
	if nexusUuid == "" {
		return fmt.Errorf("Could not find nexus uuid")
	}

	logf.Log.Info("faulting the replica")
	err = mayastorclient.FaultNexusChild(nexusNodeIP, nexusUuid, nxChildUri)

	if err != nil {
		return fmt.Errorf("failed to fault local replica, uri %s, error %v", nxChildUri, err)
	}
	return nil
}

// ScaleCsiDeployment scale up or scale down CSI deployment
func ScaleCsiDeployment(replication int32) error {
	logf.Log.Info("scaling CSI deployment", "replicas", replication)
	CsiDeployment := e2e_config.GetConfig().Product.ControlPlaneCsiController
	CsiNamespace := e2e_config.GetConfig().Product.ProductNamespace

	err := SetDeploymentReplication(CsiDeployment, CsiNamespace, &replication)
	if err != nil {
		return fmt.Errorf("failed to set csi-controller deployment replication to 1, error: %s", err.Error())
	}
	var count int
	for iter := 0; iter < 100; iter++ {
		count, err = DeploymentReadyCount(CsiDeployment, CsiNamespace)
		if err != nil {
			return fmt.Errorf("failed to get replica count, error: %s", err.Error())
		}
		if int32(count) == replication {
			logf.Log.Info("CSI controller scaling is done", "instances", replication)
			return nil
		}
		time.Sleep(time.Second)
	}
	return nil
}

func GetReplicaState(volumeUuid string, replicaUuid string) (string, error) {
	var replicaState string
	replicas, err := GetMsvReplicas(volumeUuid)
	if err != nil {
		return replicaState, fmt.Errorf("failed to list replica for volume %s, error: %v", volumeUuid, err)
	}
	var isFound bool
	for _, replica := range replicas {
		if replica.Uuid == replicaUuid {
			isFound = true
			replicaState = replica.Replica.State
			break
		}
	}
	if !isFound {
		return replicaState, fmt.Errorf("failed to find replica with uuid %s for volume %s", replicaUuid, volumeUuid)
	}
	return replicaState, err
}

func ShareReplica(poolId string, replicaUuid string) (string, error) {
	var addr string
	nodes, err := GetTestControlNodes()
	if err != nil {
		return "", err
	}
	addr = nodes[0].IPAddress

	url := fmt.Sprintf("http://%s/v0/pools/%s/replicas/%s/share/nvmf",
		k8s_portforward.TryPortForwardNode(addr, e2e_config.GetConfig().Product.KubectlPluginPort),
		poolId,
		replicaUuid)
	logf.Log.Info("", "URL", url)
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		logf.Log.Info("Error in PUT request", "node IP", addr, "url", url, "error", err)
		return "", err
	}
	req.Header.Add("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logf.Log.Info("Error while making PUT request", "url", url, "error", err)
		return "", err
	}
	jsonResponse, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		logf.Log.Info("Error while reading data", "error", err)
		return "", err
	}
	jsonStr := string(jsonResponse)
	logf.Log.Info("REPLICA SHARE ", "json response", jsonStr[1:])
	return jsonStr[1:], nil
}

func UnShareReplica(poolId string, replicaUuid string) error {
	var addr string
	nodes, err := GetTestControlNodes()
	if err != nil {
		return err
	}
	addr = nodes[0].IPAddress
	url := fmt.Sprintf("http://%s:%v/v0/pools/%s/replicas/%s/share",
		addr,
		e2e_config.GetConfig().Product.KubectlPluginPort,
		poolId,
		replicaUuid)
	logf.Log.Info("", "URL", url)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		logf.Log.Info("Error in DELETE request", "node IP", addr, "url", url, "error", err)
		return err
	}
	req.Header.Add("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logf.Log.Info("Error while making DELETE request", "url", url, "error", err)
		return err
	}
	jsonResponse, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		logf.Log.Info("Error while reading data", "error", err)
		return err
	}
	jsonStr := string(jsonResponse)
	logf.Log.Info("REPLICA UNSHARE", "json response", jsonStr)
	return nil
}

type ReplicaMountSpec struct {
	initiatorIP string
	targetIP    string
	//	uri         string
	hostNqn string
	pool    string
	uuid    string
}

type ReplicaMountStatus struct {
	devicePath string
	nqn        string
	mounted    bool
	sharedURI  string
}

type ReplicaMount struct {
	Spec   ReplicaMountSpec
	Status ReplicaMountStatus
}

func ConnectReplica(spec ReplicaMountSpec, maxRetries int) (*ReplicaMount, error) {
	var err error
	var status ReplicaMountStatus
	status.sharedURI, err = ShareReplica(spec.pool, spec.uuid)
	if err != nil {
		return nil, fmt.Errorf("failed to share replica %s %s", spec.pool, spec.uuid)
	}

	logf.Log.Info("ConnectReplica", "initiatorIP", spec.initiatorIP, "nodeIP", spec.targetIP, "shared uri", status.sharedURI)
	nqnoffset := strings.Index(status.sharedURI, "nqn.")
	if nqnoffset == -1 {
		logf.Log.Info("ConnectReplica", "Invalid URI", status.sharedURI)
		return nil, fmt.Errorf("invalid nqn URI %v", status.sharedURI)
	}
	nqnlong := status.sharedURI[nqnoffset:]
	tailoffset := strings.Index(nqnlong, "?")
	if tailoffset == -1 {
		logf.Log.Info("ConnectReplica", "Invalid URI", status.sharedURI)
		return nil, fmt.Errorf("invalid nqn URI %v", status.sharedURI)
	}
	status.nqn = nqnlong[:tailoffset]
	resp, err := agent.NvmeConnect(spec.initiatorIP, spec.targetIP, status.nqn, spec.hostNqn)
	resp = strings.TrimSpace(resp)
	if err != nil {
		logf.Log.Info("Running agent failed", "error", err)
		return nil, err
	}
	if resp != "" { // connect should be silent
		return nil, fmt.Errorf("nvme connect returned with %s", resp)
	}

	// NOTE:
	// Beyond this point we must ALWAYS disconnect the previously established connection
	// prior to returning

	// nvme list returns all nvme devices in form n * [ ... <Node> ... <Model> ... \n]
	// we want the device (=Node) associated with Model = "Mayastor NVMe controller"
	hostId := ""
	found := false
	found, status.devicePath, hostId = getNvmeDevice(spec.initiatorIP, 60, status.nqn)

	if !found {
		_, _ = agent.NvmeDisconnect(spec.initiatorIP, status.nqn)
		return nil, fmt.Errorf("failed to get device path")
	}

	if hostId != "" {
		//format the hostid into something nvme-cli can accept
		// example: 28940560b05b47b299cc4098943eb1e5 -> 28940560-b05b-47b2-99cc-4098943eb1e5
		formattedHostId := fmt.Sprintf("%s-%s-%s-%s-%s", hostId[0:8], hostId[8:12], hostId[12:16], hostId[16:20], hostId[20:])
		// disconnect nvme
		resp, err = agent.NvmeDisconnect(spec.initiatorIP, status.nqn)
		if err != nil {
			logf.Log.Info("Running agent failed", "error", err)
			return nil, err
		}
		logf.Log.Info("Executed NvmeDisconnect", "nqn", status.nqn, "got", resp)

		// nvme connect with Host uuid
		resp, err = agent.NvmeConnectWithHostId(spec.initiatorIP, spec.targetIP, status.nqn, spec.hostNqn, formattedHostId)
		resp = strings.TrimSpace(resp)
		if err != nil {
			logf.Log.Info("Running agent failed", "error", err)
			return nil, err
		}
		if resp != "" { // connect should be silent
			return nil, fmt.Errorf("nvme connect with host id returned with %s", resp)
		}
		logf.Log.Info("Executed NvmeConnectWithHostId", "nqn", status.nqn, "host nqn", spec.hostNqn, "hostid", formattedHostId, "got", resp)

		logf.Log.Info("Check device path after running nvme connect with host id to replica")
		// fetch device path again because it might change between above multiple nvme disconnect and connect
		found, status.devicePath, _ = getNvmeDevice(spec.initiatorIP, 60, status.nqn)
		if !found {
			_, _ = agent.NvmeDisconnect(spec.initiatorIP, status.nqn)
			return nil, fmt.Errorf("failed to get device path 2nd time")
		}
	}

	status.mounted = true
	logf.Log.Info("replica", "mount status", status)
	return &ReplicaMount{Spec: spec, Status: status}, nil
}

func (rm *ReplicaMount) Unmount() bool {
	if rm.Status.mounted {
		res, err := agent.NvmeDisconnect(rm.Spec.initiatorIP, rm.Status.nqn)
		logf.Log.Info("Unmount Replica", "res", res, "error", err)
		_ = UnShareReplica(rm.Spec.pool, rm.Spec.uuid)
		rm.Status.mounted = err == nil
	}
	return !rm.Status.mounted
}

// CheckReplicaComparisonPossible given a volume uuid, name, and type attempt to ensure
// that conditions replica comparison are met.
// Primarily ensure the the volume is not degraded on return from this function.
// Volumes should not be published/mounted when this function is called
func CheckReplicaComparisonPossible(volUuid string, volName string, volType common.VolumeType, timeout int) error {
	var err error
	var msv *common.MayastorVolume

	nexus, _ := GetMsvNodes(volUuid)
	for endTime := time.Now().Add(60 * time.Second); nexus != "" && time.Now().Before(endTime); time.Sleep(10 * time.Second) {
		logf.Log.Info("CheckReplicaComparisonPossible: volume has nexus: ", "volUuid", volUuid, "vol", volName, "nexus", nexus)
		nexus, _ = GetMsvNodes(volUuid)
	}
	if nexus != "" {
		return fmt.Errorf("attempting to verify replicas when nexus %v is present", nexus)
	}

	msv, err = GetMSV(volUuid)
	if err != nil {
		return fmt.Errorf("failed to retrieve mayastor volume for k8s volume %s uuid=%s, %v ", volName, volUuid, err)
	}
	if msv.State.Status != controlplane.VolStateHealthy() {
		return fmt.Errorf("mayastor volume for k8s volume %s uuid=%s, Status.State=%s, is not Online", volName, volUuid, msv.State.Status)
	}

	// double check that volume is healthy
	sleeperPodName := string(uuid.NewUUID())
	err = CreateSleepingFioPod(sleeperPodName, volName, volType)
	if err != nil {
		return fmt.Errorf("failed to create sleeping fio pod")
	}
	if !WaitPodRunning(sleeperPodName, common.NSDefault, 60) {
		_ = DeletePod(sleeperPodName, common.NSDefault)
		return fmt.Errorf("sleeping fio pod failed to start")
	}

	waitSecs := uint(timeout)

	for start := time.Now(); time.Since(start) < time.Duration(waitSecs)*time.Second; time.Sleep(10 * time.Second) {
		msv, err = GetMSV(volUuid)
		if err == nil {
			if msv.State.Status == controlplane.VolStateHealthy() {
				break
			} else {
				logf.Log.Info("CheckReplicaComparisonPossiblefioApp:", "volUuid", volUuid, "vol", volName, "msv.State.Status", msv.State.Status)
			}
		} else {
			logf.Log.Info("CheckReplicaComparisonPossible: failed to retrieve msv", "volUuid", volUuid, "vol", volName, "error", err)

			return err
		}
	}
	logf.Log.Info("CheckReplicaComparisonPossible:", "volUuid", volUuid, "vol", volName, "msv.State.Status", msv.State.Status)
	if msv.State.Status != controlplane.VolStateHealthy() {
		return fmt.Errorf("volume is not healthy aborting replica comparison msv.State.Status is %s", msv.State.Status)
	}
	err = DeletePod(sleeperPodName, common.NSDefault)
	if err != nil {
		return fmt.Errorf("failed to delete sleeping fio pod")
	}
	nexus, _ = GetMsvNodes(volUuid)
	for start := time.Now(); nexus != "" && time.Since(start) < time.Second*120; time.Sleep(20 * time.Second) {
		logf.Log.Info("CheckReplicaComparisonPossible: volume has: ", "volUuid", volUuid, "vol", volName, "sleeper pod nexus", nexus)
		nexus, _ = GetMsvNodes(volUuid)
	}
	if nexus != "" {
		return fmt.Errorf("attempting to verify replicas when nexus %v is present", nexus)
	}

	return err
}
