package v1_rest_api

import (
	"context"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/e2e_config"
	"github.com/openebs/openebs-e2e/common/k8s_portforward"

	//	"encoding/json"
	"fmt"
	"net/http"

	openapiClient "github.com/openebs/openebs-e2e/common/generated/openapi"
)

type OAClientWrapper struct {
	nodes   []string
	clients []*openapiClient.APIClient
	clindex uint
}

func makeWrapper(nodeIPs []string) OAClientWrapper {
	if len(nodeIPs) == 0 {
		panic(fmt.Errorf("unable to create any openapi clients nodes=%v", nodeIPs))
	}
	oacw := OAClientWrapper{
		nodes: nodeIPs,
	}
	cfg := openapiClient.NewConfiguration()

	var err error
	cfg.Host, err = k8s_portforward.PortForwardService(e2e_config.GetConfig().Product.RestApiService,
		e2e_config.GetConfig().Product.ProductNamespace,
		e2e_config.GetConfig().Product.RestApiPort)
	if err != nil {
		panic(err)
	}
	cfg.Scheme = `http`
	client := openapiClient.NewAPIClient(cfg)
	oacw.clients = append(oacw.clients, client)
	return oacw
}

func (oacw *OAClientWrapper) client() *openapiClient.APIClient {
	oacw.clindex = (1 + oacw.clindex) % uint(len(oacw.clients))
	return oacw.clients[oacw.clindex]
}

func (oacw OAClientWrapper) getVolume(uuid string) (openapiClient.Volume, error, int) {
	var statusCode int

	req := oacw.client().VolumesAPI.GetVolume(context.TODO(), uuid)
	volume, resp, err := req.Execute()
	if resp != nil {
		statusCode = resp.StatusCode
	}
	return *volume, err, statusCode
}

func (oacw OAClientWrapper) getVolumes() (openapiClient.Volumes, error, int) {
	var volumes *openapiClient.Volumes
	var statusCode int
	var resp *http.Response
	var err error

	// FIXME: Handle pagination
	req := oacw.client().VolumesAPI.GetVolumes(context.TODO()).StartingToken(0).MaxEntries(1000)
	volumes, resp, err = req.Execute()
	if resp != nil {
		statusCode = resp.StatusCode
	}
	return *volumes, err, statusCode
}

func (oacw OAClientWrapper) deleteVolume(uuid string) (error, int) {
	var statusCode int

	req := oacw.client().VolumesAPI.DelVolume(context.TODO(), uuid)
	resp, err := req.Execute()
	if resp != nil {
		statusCode = resp.StatusCode
	}
	return err, statusCode
}

func (oacw OAClientWrapper) putReplicaCount(uuid string, replicaCount int) (error, int) {
	var statusCode int

	req := oacw.client().VolumesAPI.PutVolumeReplicaCount(context.TODO(), uuid, int32(replicaCount))
	_, resp, err := req.Execute()
	if resp != nil {
		statusCode = resp.StatusCode
	}
	return err, statusCode
}

func (oacw OAClientWrapper) volToMsv(vol openapiClient.Volume) common.MayastorVolume {
	var nexusChildren []common.TargetChild
	var replicaTopology = make(common.ReplicaTopology)

	volSpec := vol.GetSpec()
	volState := vol.GetState()
	nexus, ok := volState.GetTargetOk()
	if !ok {
		nexus = &openapiClient.Nexus{}
	}

	for _, inChild := range nexus.Children {
		nexusChildren = append(nexusChildren, common.TargetChild{
			Uri:             inChild.Uri,
			State:           string(inChild.State),
			RebuildProgress: inChild.RebuildProgress,
		})
	}

	for k, topology := range volState.ReplicaTopology {
		replicaTopology[k] = common.Replica{
			Node:        topology.GetNode(),
			State:       string(topology.GetState()),
			Pool:        topology.GetPool(),
			ChildStatus: string(topology.GetChildStatus()),
			Usage: common.ReplicaUsage{
				Capacity:              topology.Usage.Capacity,
				Allocated:             topology.Usage.Allocated,
				AllocatedSnapshots:    topology.Usage.AllocatedSnapshots,
				AllocatedAllSnapshots: topology.Usage.AllocatedAllSnapshots,
			},
		}
	}

	return common.MayastorVolume{
		Spec: common.MsvSpec{
			Num_replicas: int(volSpec.NumReplicas),
			Size:         volSpec.Size,
			Status:       string(volSpec.Status),
			Target: common.SpecTarget{
				Protocol: string(vol.GetSpec().Target.GetProtocol()),
				Node:     string(vol.GetSpec().Target.GetNode()),
			},
			Uuid: volSpec.Uuid,
			Thin: volSpec.GetThin(),
			Policy: common.Policy{
				Self_heal: volSpec.GetPolicy().SelfHeal,
			},
			NumSnapshots: volSpec.NumSnapshots,
			AsThin:       *volSpec.AsThin,
			ContentSource: common.ContentSource{
				Snapshot: common.Snapshot(*volSpec.GetContentSource().Snapshot),
			},
			MaxSnapshots: *volSpec.MaxSnapshots,
		},
		State: common.MsvState{
			Target: common.StateTarget{
				Children:  nexusChildren,
				DeviceUri: nexus.DeviceUri,
				Node:      nexus.Node,
				Rebuilds:  nexus.Rebuilds,
				Protocol:  string(volSpec.Target.GetProtocol()),
				Size:      nexus.GetSize(),
				State:     string(nexus.GetState()),
				Uuid:      nexus.Uuid,
			},
			Size:            volState.GetSize(),
			Status:          string(volState.GetStatus()),
			Uuid:            volState.GetUuid(),
			ReplicaTopology: replicaTopology,
			Usage: common.Usage{
				Capacity:                volState.Usage.Capacity,
				Allocated:               volState.Usage.Allocated,
				AllocatedReplica:        volState.Usage.AllocatedReplica,
				AllocatedSnapshots:      volState.Usage.AllocatedSnapshots,
				AllocatedAllSnapshots:   volState.Usage.AllocatedAllSnapshots,
				TotalAllocated:          volState.Usage.TotalAllocated,
				TotalAllocatedSnapshots: volState.Usage.TotalAllocatedSnapshots,
				TotalAllocatedReplicas:  volState.Usage.TotalAllocatedReplicas.(int64),
			},
		},
	}
}

func (oacw OAClientWrapper) getNode(nodeName string) (openapiClient.Node, error) {
	req := oacw.client().NodesAPI.GetNode(context.TODO(), nodeName)
	node, _, err := req.Execute()
	return *node, err
}

func (oacw OAClientWrapper) getNodes() ([]openapiClient.Node, error) {
	req := oacw.client().NodesAPI.GetNodes(context.TODO())
	nodes, _, err := req.Execute()
	return nodes, err
}

func (oacw OAClientWrapper) getPool(poolName string) (openapiClient.Pool, error, int) {
	var statusCode int

	req := oacw.client().PoolsAPI.GetPool(context.TODO(), poolName)
	pool, resp, err := req.Execute()
	if resp != nil {
		statusCode = resp.StatusCode
	}
	return *pool, err, statusCode
}

func (oacw OAClientWrapper) getPools() ([]openapiClient.Pool, error, int) {
	var statusCode int

	req := oacw.client().PoolsAPI.GetPools(context.TODO())
	pools, resp, err := req.Execute()
	if resp != nil {
		statusCode = resp.StatusCode
	}
	return pools, err, statusCode
}

func (oacw OAClientWrapper) setVolumeMaxSnapshotCount(uuid string, maxSnapshotCount int32) (error, int) {
	var statusCode int

	req := oacw.client().VolumesAPI.PutVolumeProperty(context.TODO(), uuid)
	maxCount := maxSnapshotCount
	req.SetVolumePropertyBody(openapiClient.SetVolumePropertyBody{
		MaxSnapshots: &maxCount,
	})
	_, resp, err := req.Execute()
	if resp != nil {
		statusCode = resp.StatusCode
	}
	return err, statusCode
}

// { Commented out implementation of license support
//    so that we can use the latest Mayastor openapi spec - which does not include the License API
//
//func (oacw OAClientWrapper) getLicense() (openapiClient.License, error) {
//	req := oacw.client().LicenseApi.GetLicense(context.TODO())
//	license, _, err := req.Execute()
//	return license, err
//}
//
//func (oacw OAClientWrapper) installLicense(rawLicenseData []byte) (string, error, int) {
//	var result string
//	var err error
//	var licenseData map[string]interface{}
//	var statusCode int
//	var resp *http.Response
//
//	if err = json.Unmarshal(rawLicenseData, &licenseData); err != nil {
//		return "", err, 0
//	}
//
//	req := oacw.client().LicenseApi.PutLicense(context.TODO())
//	req = req.RequestBody(licenseData)
//	result, resp, err = req.Execute()
//	if resp != nil {
//		statusCode = resp.StatusCode
//	}
//	return result, err, statusCode
//}
//
//func (oacw OAClientWrapper) uninstallLicense(uuid string) error {
//	req := oacw.client().LicenseApi.DelLicense(context.TODO(), uuid)
//	_, err := req.Execute()
//	return err
//}
//} Commented out implementation of license support */
