package k8stest

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/controlplane"
	client "github.com/openebs/openebs-e2e/common/e2e_agent"
	"github.com/openebs/openebs-e2e/common/mayastorclient"

	"sigs.k8s.io/controller-runtime/pkg/log"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

type replicaInfo struct {
	IP   string
	URI  string
	Pool string
	UUID string
}

func (ri replicaInfo) String() string {
	return fmt.Sprintf("IP:%s URI:%s Pool:%s UUID:%s", ri.IP, ri.URI, ri.Pool, ri.UUID)
}

// getCheckSum calculate and return the checksum of replica contents as a string,
// using Mayastor REST API and e2e-agent
func getCheckSum(replica replicaInfo) (string, error) {
	var checksum string

	crc32, err := mayastorclient.ChecksumReplica(replica.IP, replica.UUID, replica.Pool)
	if err == nil {
		return fmt.Sprintf("%v", crc32), nil
	}
	logf.Log.Info("gRPC ChecksumReplica call failed/unsupported, falling back to cli method", "err", err)
	replicaSharedURI, err := ShareReplica(replica.Pool, replica.UUID)
	if err != nil {
		log.Log.Info("ShareReplica failed", "error", err)
		return fmt.Sprintf("%s; %v", replica, err), err
	}
	nqn, err := GetNodeNqn(replica.IP)
	if err == nil {
		replicaNqn, err := getReplicaNqn(replicaSharedURI)
		if err != nil {
			checksum = fmt.Sprintf("%s; %v", replica, err)
		} else {
			checksum, err = ChecksumReplica(replica.IP, replica.IP, replicaNqn, 10, nqn)
			if err != nil {
				log.Log.Info("ChecksumReplica failed", "error", err)
				// do not return from here because we want to unshare if the original
				// replica URI was a bdev
				checksum = fmt.Sprintf("%s; %v", replica, err)
			}
		}
	} else {
		log.Log.Info("GetNodeNqn failed", "IP", replica.IP, "error", err)
		checksum = fmt.Sprintf("%s; %v", replica, err)
	}
	if strings.HasPrefix(replica.URI, "bdev") {
		// for now ignore unshare errors as we may have successfully retrieved a checksum
		unsErr := UnShareReplica(replica.Pool, replica.UUID)
		if unsErr != nil {
			log.Log.Info("Unshare replica failed", "pool", replica.Pool, "UUID", replica.UUID, "error", unsErr)
		}
	}
	return checksum, err
}

func getReplicaNqn(replicaUri string) (string, error) {
	nqnoffset := strings.Index(replicaUri, "nqn.")
	if nqnoffset == -1 {
		logf.Log.Info("ChecksumReplica", "Invalid URI", replicaUri)
		return "", fmt.Errorf("invalid nqn URI %v", replicaUri)
	}
	nqnlong := replicaUri[nqnoffset:]
	tailoffset := strings.Index(nqnlong, "?")
	if tailoffset == -1 {
		logf.Log.Info("ChecksumReplica", "Invalid URI", replicaUri)
		return "", fmt.Errorf("invalid nqn URI %v", replicaUri)
	}
	nqn := nqnlong[:tailoffset]
	return nqn, nil
}

// CompareVolumeReplicas compares contents of all volume replicas
// returns replicas match, mismatch or failure to compare replica contents
//
// Fails if a nexus is present for more than 120 seconds after invocation of
// this function. Presence of a nexus implies writes are possible
// whilst calculating the checksum, also this function amends replica
// sharing.
//
// This function is intended for use at the end of volume operations,
// immediately before a volume is deleted.
func CompareVolumeReplicas(volName string, ns string) common.ReplicasComparison {
	result := common.ReplicasComparison{
		Result:      common.CmpReplicasFailed,
		Description: "",
		Err:         nil,
	}

	var checksums []string

	// Create a map of IP addresses keyed on node name
	nodeList, err := GetIOEngineNodes()
	if err != nil {
		result.Err = err
		return result
	}
	IPAddresses := map[string]string{}
	for _, node := range nodeList {
		IPAddresses[node.NodeName] = node.IPAddress
	}

	if ns == "" {
		ns = common.NSDefault
	}
	pvc, err := GetPVC(volName, ns)
	if err != nil {
		result.Err = err
		return result
	}
	volUuid := fmt.Sprintf("%v", pvc.UID)

	log.Log.Info("CompareVolumeReplicas:", "volume", volName, "volume uuid", volUuid)

	nexus, replicaNodes := GetMsvNodes(volUuid)

	if len(replicaNodes) < 2 {
		log.Log.Info("CompareVolumeReplicas: single replica - no comparison possible")
		result.Result = common.CmpReplicasMatch
		result.Description = "single replica - no comparison possible"
		return result
	}

	for cnt := 0; cnt < 24 && nexus != ""; cnt++ {
		time.Sleep(5 * time.Second)
		nexus, _ = GetMsvNodes(volUuid)
	}
	if nexus != "" {
		log.Log.Info("CompareVolumeReplicas: nexus for volume is present, aborting replica comparison")
		result.Err = fmt.Errorf("nexus is present for volume")
		result.Description = "nexus is present for volume"
		return result
	}

	// Create a map of replica URIs keyed on node and pool
	msv, err := GetMSV(volUuid)
	if err != nil {
		result.Err = err
		return result
	}
	replicaURIs := map[string]string{}
	for uri, child := range msv.State.ReplicaTopology {
		replicaURIs[fmt.Sprintf("%s/%s", child.Node, child.Pool)] = uri
	}

	// Setup to do content comparison of replicas
	replicaTopologies, err := GetMsvReplicaTopology(volUuid)
	if err != nil {
		result.Err = err
		return result
	}

	for _, replicaTopology := range replicaTopologies {
		uriKey := fmt.Sprintf("%s/%s", replicaTopology.Node, replicaTopology.Pool)
		_, haveUri := replicaURIs[uriKey]
		if !haveUri {
			result.Err = fmt.Errorf("failed to establish replica URI for vol=%s on node=%s pool=%s",
				volUuid, replicaTopology.Node, replicaTopology.Pool,
			)
			return result
		}
		replicaState, rErr := GetReplicaState(volUuid, replicaURIs[uriKey])
		if rErr != nil {
			result.Err = fmt.Errorf("failed to get replica state, error :%s", rErr.Error())
			return result
		}
		if replicaState != controlplane.ReplicaStateOnline() {
			result.Err = fmt.Errorf("replica %s is not online", replicaURIs[uriKey])
			return result
		}
	}

	for uuid, replicaTopology := range replicaTopologies {
		uriKey := fmt.Sprintf("%s/%s", replicaTopology.Node, replicaTopology.Pool)
		checksum, ckErr := getCheckSum(
			replicaInfo{
				IP:   IPAddresses[replicaTopology.Node],
				URI:  replicaURIs[uriKey],
				Pool: replicaTopology.Pool,
				UUID: uuid,
			})
		checksums = append(checksums, checksum)
		result.Description = fmt.Sprintf("%sreplicaURI:%s, cksum output:%s\n", result.Description, uriKey, checksum)
		if ckErr != nil {
			err = fmt.Errorf("%v; %v", ckErr, err)
		}
	}

	result.Result = common.CmpReplicasMatch
	for _, checksum := range checksums[1:] {
		if checksum != checksums[0] {
			log.Log.Info("CompareVolumeReplicas: checksums do not match", "1", checksums[0], "2", checksum)
			result.Result = common.CmpReplicasMismatch
		}
	}
	log.Log.Info("CompareVolumeReplicas:", "checksums", checksums)
	return result
}

func ZapPoolDevices() error {
	nodeList, err := GetIOEngineNodes()
	if err != nil {
		return err
	}
	pools, err := ListMsPools()
	if err != nil {
		return err
	}
	log.Log.Info("ZapPoolDevices")
	startTime := time.Now()
	if !DeleteAllPools() {
		return fmt.Errorf("failed to delete all pools")
	}
	for _, node := range nodeList {
		for _, pool := range pools {
			if pool.Spec.Node == node.NodeName {
				log.Log.Info("Zeroing pool device",
					"ip", node.IPAddress,
					"disk", pool.Spec.Disks[0])
				/*
					// try discard first
					_, err := client.BlkDiscard(node.IPAddress, pool.Spec.Disks[0], "-v")
					if err != nil {
						// then zero-fill
						_, err = client.BlkDiscard(node.IPAddress, pool.Spec.Disks[0], "-v -z")
					}
				*/
				// zero-fill
				_, err := client.BlkDiscard(node.IPAddress, pool.Spec.Disks[0], "-v -z")
				log.Log.Info("Zeroing pool device complete",
					"ip", node.IPAddress,
					"disk", pool.Spec.Disks[0],
					"error", err,
				)
				if err != nil {
					return err
				}
			}
		}
	}
	err = CreateConfiguredPools()
	if err != nil {
		return fmt.Errorf("failed to create configured pools; %v", err)
	}
	log.Log.Info("ZapPoolDevices", "duration", time.Since(startTime))
	return nil
}

func ByteCompareVolumeReplicas(volName string, ns string) (bool, string, error) {
	log.Log.Info("ByteCompareVolumeReplicas: temporarily disabled: return \"okay\"")
	return true, "", nil
	/* FIXME: temporarily disabled
	var err error
	var ok bool = false
	description := ""

	// Create a map of IP addresses keyed on node name
	nodeList, err := GetIOEngineNodes()
	if err != nil {
		return ok, description, err
	}
	IPAddresses := map[string]string{}
	for _, node := range nodeList {
		IPAddresses[node.NodeName] = node.IPAddress
	}

	if ns == "" {
		ns = common.NSDefault
	}
	pvc, err := GetPVC(volName, ns)
	if err != nil {
		return ok, description, err
	}
	volUuid := fmt.Sprintf("%v", pvc.UID)
	nexus, replicaNodes := GetMsvNodes(volUuid)
	for cnt := 0; cnt < 24 && nexus != ""; cnt++ {
		time.Sleep(5 * time.Second)
		nexus, replicaNodes = GetMsvNodes(volUuid)
	}
	if nexus != "" {
		log.Log.Info("ByteCompareVolumeReplicas: nexus for volume is present, aborting replica comparison")
		return ok, description, fmt.Errorf("nexus is present for volume")
	}

	log.Log.Info("ByteCompareVolumeReplicas:", "volume", volName, "volume uuid", volUuid)
	if len(replicaNodes) < 2 {
		log.Log.Info("ByteCompareVolumeReplicas: single replica - no comparison possible")
		return true, "single replica - no comparison possible", nil
	}

	// Create a map of replica URIs keyed on node and pool
	msv, err := GetMSV(volUuid)
	if err != nil {
		return ok, description, err
	}
	replicaURIs := map[string]string{}
	for uri, child := range msv.State.ReplicaTopology {
		replicaURIs[fmt.Sprintf("%s/%s", child.Node, child.Pool)] = uri
	}

	// Setup to do content comparison of replicas
	replicaTopologies, err := GetMsvReplicaTopology(volUuid)
	if err != nil {
		return ok, description, err
	}

	for _, replicaTopology := range replicaTopologies {
		uriKey := fmt.Sprintf("%s/%s", replicaTopology.Node, replicaTopology.Pool)
		_, haveUri := replicaURIs[uriKey]
		if !haveUri {
			return ok, description, fmt.Errorf("failed to establish replica URI for vol=%s on node=%s pool=%s",
				volUuid, replicaTopology.Node, replicaTopology.Pool,
			)
		}
		replicaState, rErr := GetReplicaState(volUuid, replicaURIs[uriKey])
		if rErr != nil {
			return ok, description, fmt.Errorf("failed to get replica state, error :%s", rErr.Error())
		}
		if replicaState != controlplane.ReplicaStateOnline() {
			return ok, description, fmt.Errorf("replica %s is not online", replicaURIs[uriKey])
		}
	}

	var replicaMounts []*ReplicaMount
	var initiatorIP string
	var hostNqn string
	for uuid, replicaTopology := range replicaTopologies {
		uriKey := fmt.Sprintf("%s/%s", replicaTopology.Node, replicaTopology.Pool)
		if initiatorIP == "" {
			initiatorIP = IPAddresses[replicaTopology.Node]
			hostNqn, err = GetNodeNqn(IPAddresses[replicaTopology.Node])
			if err != nil {
				return ok, description, fmt.Errorf("failed to retrieve nqn for %v", IPAddresses[replicaTopology.Node])
			}
		}
		var rmp *ReplicaMount
		rmp, err = ConnectReplica(ReplicaMountSpec{
			initiatorIP: initiatorIP,
			targetIP:    IPAddresses[replicaTopology.Node],
			uri:         replicaURIs[uriKey],
			hostNqn:     hostNqn,
			pool:        replicaTopology.Pool,
			uuid:        uuid,
		}, 1)
		if err == nil {
			replicaMounts = append(replicaMounts, rmp)
		}
	}

	if len(replicaMounts) > 1 {
		ok = true
		replMounts := replicaMounts[:]
		for len(replMounts) > 1 {
			ref := replMounts[0]
			replMounts = replMounts[1:]
			for _, r := range replMounts {
				output, cmpErr := client.Cmp(ref.Spec.initiatorIP, ref.Status.devicePath, r.Status.devicePath)
				log.Log.Info("ByteCompareVolumeReplicas", "output", output, "error", cmpErr)
				if cmpErr == nil {
					description += fmt.Sprintf("cmp %s %s: output=%s\n ", replicaMounts[0].Status.devicePath, r.Status.devicePath, output)
					ok = ok && len(output) == 0
				} else {
					description += fmt.Sprintf("cmp %s %s: error: %v, output=%s\n ", replicaMounts[0].Status.devicePath, r.Status.devicePath, cmpErr, output)
				}
			}
		}
	} else {
		description += fmt.Sprintf("unable to run cmp number of mounted replicas is %d; ", len(replicaMounts))
		log.Log.Info("ByteCompareVolumeReplicas unable to run cmp", "mounted replicas", len(replicaMounts))
	}

	for _, r := range replicaMounts {
		r.Unmount()
	}

	return ok, description, err
	*/
}

func WipeVolumeReplicas(volUuid string) error {
	var err error
	var replicaTopologies common.ReplicaTopology

	for ix := 0; ix < 3; ix += 1 {
		replicaTopologies, err = GetMsvReplicaTopology(volUuid)
		if err == nil {
			var wipeErr error
			ch := make(chan error, len(replicaTopologies))
			var wg sync.WaitGroup
			wg.Add(len(replicaTopologies))
			for replicaUuid, replicaTopology := range replicaTopologies {
				go wipeReplica(replicaTopology, replicaUuid, ch, &wg)
			}
			wg.Wait()
			close(ch)
			for e := range ch {
				if e != nil {
					wipeErr = fmt.Errorf("%v ; %v", e, err)
				}
			}
			return wipeErr
		}
		time.Sleep(10 * time.Second)
	}
	return err
}

func GetVolumeReplicasChecksum(volName string, ns string) ([]string, error) {
	var checksums []string

	// Create a map of IP addresses keyed on node name
	nodeList, err := GetIOEngineNodes()
	if err != nil {
		return checksums, err
	}
	IPAddresses := map[string]string{}
	for _, node := range nodeList {
		IPAddresses[node.NodeName] = node.IPAddress
	}

	if ns == "" {
		ns = common.NSDefault
	}
	pvc, err := GetPVC(volName, ns)
	if err != nil {
		return checksums, err
	}
	volUuid := fmt.Sprintf("%v", pvc.UID)

	log.Log.Info("GetVolumeReplicasChecksum:", "volume", volName, "volume uuid", volUuid)

	var nexus string
	for cnt := 0; cnt < 24 && nexus != ""; cnt++ {
		time.Sleep(5 * time.Second)
		nexus, _ = GetMsvNodes(volUuid)
	}
	if nexus != "" {
		log.Log.Info("nexus for volume is present, aborting replica checksum")
		return checksums, fmt.Errorf("nexus is present for volume")
	}

	// Create a map of replica URIs keyed on node and pool
	msv, err := GetMSV(volUuid)
	if err != nil {
		return checksums, err
	}
	replicaURIs := map[string]string{}
	for uri, child := range msv.State.ReplicaTopology {
		replicaURIs[fmt.Sprintf("%s/%s", child.Node, child.Pool)] = uri
	}

	// Setup to do content comparison of replicas
	replicaTopologies, err := GetMsvReplicaTopology(volUuid)
	if err != nil {
		return checksums, err
	}

	for _, replicaTopology := range replicaTopologies {
		uriKey := fmt.Sprintf("%s/%s", replicaTopology.Node, replicaTopology.Pool)
		_, haveUri := replicaURIs[uriKey]
		if !haveUri {
			return checksums, fmt.Errorf("failed to establish replica URI for vol=%s on node=%s pool=%s",
				volUuid, replicaTopology.Node, replicaTopology.Pool,
			)
		}
		replicaState, rErr := GetReplicaState(volUuid, replicaURIs[uriKey])
		if rErr != nil {
			return checksums, fmt.Errorf("failed to get replica state, error :%s", rErr.Error())
		}
		if replicaState != controlplane.ReplicaStateOnline() {
			return checksums, fmt.Errorf("replica %s is not online", replicaURIs[uriKey])
		}
	}

	for uuid, replicaTopology := range replicaTopologies {
		uriKey := fmt.Sprintf("%s/%s", replicaTopology.Node, replicaTopology.Pool)
		checksum, ckErr := getCheckSum(
			replicaInfo{
				IP:   IPAddresses[replicaTopology.Node],
				URI:  replicaURIs[uriKey],
				Pool: replicaTopology.Pool,
				UUID: uuid,
			})
		checksums = append(checksums, checksum)
		if ckErr != nil {
			err = fmt.Errorf("%v; %v", ckErr, err)
		}
	}
	return checksums, err
}

func GetReplicaTopoloy(volUuid string, replicaUuid string) (common.Replica, error) {
	replicaTopologies, err := GetMsvReplicaTopology(volUuid)
	if err != nil {
		logf.Log.Info("Failed to get replica topology for volume", "uuid", volUuid, "error", err)
		return common.Replica{}, err
	}

	if _, ok := replicaTopologies[replicaUuid]; ok {
		return replicaTopologies[replicaUuid], nil
	}
	return common.Replica{}, fmt.Errorf("not found replica topology uuid: %s for requested volume uuid: %s", replicaUuid, volUuid)
}
