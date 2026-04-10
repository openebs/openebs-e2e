package offline_pool_delete

import (
	"bytes"
	"fmt"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/controlplane"
	cpv1 "github.com/openebs/openebs-e2e/common/controlplane/v1"
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
	logf.Log.Info("Deleting offline node via plugin", "node", nodeName, "flags", flags)
	err := controlplane.DeleteOfflineNodeViaPlugin(nodeName, flags...)
	if err != nil {
		logf.Log.Error(err, "Failed to delete offline node via plugin", "node", nodeName, "flags", flags)
		return err
	}
	return nil
}

// DeleteOfflineNodeWithOutput runs the same kubectl mayastor delete node as DeleteOfflineNode and returns
// combined stdout/stderr so tests can assert on plugin output (for example data loss / snapshot loss details).
func DeleteOfflineNodeWithOutput(nodeName string, flags ...common.OfflinePoolDelete) (string, error) {
	args := []string{"-n", common.NSMayastor(), "delete", "node", nodeName}
	for _, flag := range flags {
		// common.OfflinePoolDelete.String() already returns the CLI flag name
		// (e.g. "purge", "yes", "accept-volume-loss", "accept-snapshot-loss").
		if s := flag.String(); s != "" {
			args = append(args, "--"+s)
		}
	}
	logf.Log.Info("Deleting offline node via plugin", "node", nodeName, "args", args)
	cmd := cpv1.GetMayastorPluginCmd(args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return out.String(), fmt.Errorf("plugin failed to delete node %s, error %v, output: %s", nodeName, err, out.String())
	}
	return out.String(), nil
}
