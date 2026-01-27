package offline_pool_delete

import (
	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/controlplane"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// FixME : error messages related to offline Node delete

var (
	NodeOnlineState                = "node is not offline"
	DeleteNodeWithoutCordon        = "cannot delete node that is not cordoned"
	DeleteNodeWithResources        = "node has resources, cannot delete without confirm flag"
	DeleteNodeWithoutPurgeFlag     = "purge flag is required to delete Node, Node is not in online state"
	DeleteOfflineNodeWithSnapshots = "node has snapshots, cannot delete without confirm snapshot loss flag"
	DeleteOfflineNodeWithData      = "node has data, cannot delete without confirm data loss flag" // last healthy replica is scheduled on this node
)

// DeleteOfflineNode deletes the offline Node via control plane plugin
func DeleteOfflineNode(nodeName string, flags ...common.OfflinePoolDelete) error {
	logf.Log.Info("Deleting offline pool via plugin", "pool", nodeName, "flags", flags)
	err := controlplane.DeleteOfflineNodeViaPlugin(nodeName, flags...)
	if err != nil {
		logf.Log.Error(err, "Failed to delete offline pool via plugin", "pool", nodeName, "flags", flags)
		return err
	}
	return nil
}
