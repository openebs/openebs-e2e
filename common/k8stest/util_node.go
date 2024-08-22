package k8stest

// Utility functions for manipulation of nodes.
import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/openebs/openebs-e2e/common/e2e_agent"

	"github.com/openebs/openebs-e2e/common/e2e_config"

	coordinationV1 "k8s.io/api/coordination/v1"
	coreV1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

type NodeLocation struct {
	NodeName        string
	IPAddress       string
	MayastorNode    bool
	K8sControlPlane bool
	ControlNode     bool
	ExtIPAddress    string
}

type IOEngineNodeLocation struct {
	NodeName  string
	IPAddress string
}

func GetIOEngineNodeLocsMap() (map[string]IOEngineNodeLocation, error) {
	NodeLocsMap := make(map[string]IOEngineNodeLocation)
	nodeLocs, err := GetIOEngineNodes()
	if err != nil {
		return nil, err
	}
	for _, node := range nodeLocs {
		NodeLocsMap[node.NodeName] = node
	}
	return NodeLocsMap, nil
}

// getNodeLocs returns vector of populated NodeLocation structs
func getNodeLocs() ([]NodeLocation, error) {
	nodeList, err := gTestEnv.KubeInt.CoreV1().Nodes().List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return nil, errors.New("failed to list nodes")
	}
	NodeLocs := make([]NodeLocation, 0, len(nodeList.Items))
	for _, k8snode := range nodeList.Items {
		addrstr := ""
		namestr := ""
		extaddrstr := ""
		mayastorNode := false
		isK8sControlPlane := false
		controlNode := false
		if value, ok := k8snode.Labels[e2e_config.GetConfig().Product.EngineLabel]; ok && value == e2e_config.GetConfig().Product.EngineLabelValue {
			mayastorNode = true
		}
		isK8sControlPlane = ContainsK8sControlPlaneLabels(k8snode.Labels)
		_, controlNode = k8snode.Labels[e2e_config.GetConfig().TestControlNodeLabel]

		// master nodes which are not control nodes are of no interest
		if isK8sControlPlane && !controlNode {
			continue
		}

		for _, addr := range k8snode.Status.Addresses {
			if addr.Type == coreV1.NodeInternalIP {
				addrstr = addr.Address
			}
			if addr.Type == coreV1.NodeExternalIP {
				extaddrstr = addr.Address
			}
			if addr.Type == coreV1.NodeHostName {
				namestr = addr.Address
			}
		}
		if namestr != "" && addrstr != "" {
			NodeLocs = append(NodeLocs, NodeLocation{
				NodeName:        namestr,
				IPAddress:       addrstr,
				MayastorNode:    mayastorNode,
				K8sControlPlane: isK8sControlPlane,
				ControlNode:     controlNode,
				ExtIPAddress:    extaddrstr,
			})
		} else {
			return nil, errors.New("node lacks expected fields")
		}
	}
	return NodeLocs, nil
}

func GetIOEngineNodes() ([]IOEngineNodeLocation, error) {
	IOEngineNodeLocs := make([]IOEngineNodeLocation, 0)
	allNodeLocs, err := getNodeLocs()
	if err == nil {
		for _, nodeLoc := range allNodeLocs {
			if nodeLoc.MayastorNode {
				IOEngineNodeLocs = append(IOEngineNodeLocs, IOEngineNodeLocation{
					NodeName:  nodeLoc.NodeName,
					IPAddress: nodeLoc.IPAddress,
				})
			}
		}
	}
	return IOEngineNodeLocs, err
}

func GetTestControlNodes() ([]NodeLocation, error) {
	NodeLocs := make([]NodeLocation, 0)
	allNodeLocs, err := getNodeLocs()
	if err == nil {
		for _, nodeLoc := range allNodeLocs {
			if nodeLoc.ControlNode {
				NodeLocs = append(NodeLocs, nodeLoc)
			}
		}
	}
	if len(NodeLocs) == 0 {
		err = fmt.Errorf("no test control nodes found")
	}
	return NodeLocs, err
}

// GetNodeIPAddress returns IP address of a node
func GetNodeIPAddress(nodeName string) (*string, error) {
	nodeLocs, err := getNodeLocs()
	if err != nil {
		return nil, err
	}
	for _, nl := range nodeLocs {
		if nodeName == nl.NodeName {
			return &nl.IPAddress, nil
		}
	}
	return nil, fmt.Errorf("node %s not found", nodeName)
}

// GetNodeIPAddress returns IP address of a node
func GetNodeExtIPAddress(nodeName string) (*string, error) {
	nodeLocs, err := getNodeLocs()
	if err != nil {
		return nil, err
	}
	for _, nl := range nodeLocs {
		if nodeName == nl.NodeName {
			return &nl.ExtIPAddress, nil
		}
	}
	return nil, fmt.Errorf("node %s not found", nodeName)
}

// GetMayastorNodeIPAddresses return an array of IP addresses for nodes
// running mayastor. On error an empty array is returned.
func GetMayastorNodeIPAddresses() []string {
	var addrs []string
	nodes, err := getNodeLocs()
	if err != nil {
		return addrs
	}

	for _, node := range nodes {
		if node.MayastorNode {
			addrs = append(addrs, node.IPAddress)
		}
	}
	return addrs
}

func GetMayastorNodeNames() ([]string, error) {
	var nodeNames []string
	nodes, err := getNodeLocs()
	if err != nil {
		return nodeNames, err
	}

	for _, node := range nodes {
		if node.MayastorNode {
			nodeNames = append(nodeNames, node.NodeName)
		}
	}
	return nodeNames, err
}

// LabelNode add a label to a node
// label is a string in the form "key=value"
// function still succeeds if label already present
func LabelNode(nodename string, label string, value string) error {
	// TODO remove dependency on kubectl
	labelAssign := fmt.Sprintf("%s=%s", label, value)
	cmd := exec.Command("kubectl", "label", "node", nodename, labelAssign, "--overwrite=true")
	cmd.Dir = ""
	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to label node %s, label: %s, error: %v", nodename, labelAssign, err)
	}
	return nil
}

// UnlabelNode remove a label from a node
// function still succeeds if label not present
func UnlabelNode(nodename string, label string) error {
	// TODO remove dependency on kubectl
	cmd := exec.Command("kubectl", "label", "node", nodename, label+"-")
	cmd.Dir = ""
	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to remove label from node %s, label: %s, error: %v", nodename, label, err)
	}
	return nil
}

func AddNoScheduleTaintOnNode(nodeName string) error {
	cmd := exec.Command("kubectl", "taint", "node", nodeName, "node-role.kubernetes.io/nodeName:NoSchedule")
	cmd.Dir = ""
	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to add no schedule taint to node %s, error: %v", nodeName, err)
	}
	return nil
}

func RemoveNoScheduleTaintFromNode(nodeName string) error {
	cmd := exec.Command("kubectl", "taint", "node", nodeName, "node-role.kubernetes.io/nodeName:NoSchedule"+"-")
	cmd.Dir = ""
	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to remove no schedule taint from node %s, error: %v", nodeName, err)
	}
	return nil
}

func AreNodesReady() (bool, error) {
	nodes, err := gTestEnv.KubeInt.CoreV1().Nodes().List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return false, err
	}
	for _, node := range nodes.Items {
		readyStatus, err := IsNodeReady(node.Name, &node)
		if err != nil {
			return false, err
		}
		if !readyStatus {
			return false, nil
		}
	}
	return true, nil
}

func IsNodeReady(nodeName string, node *v1.Node) (bool, error) {
	var err error
	if node == nil {
		node, err = gTestEnv.KubeInt.CoreV1().Nodes().Get(context.TODO(), nodeName, metaV1.GetOptions{})
		if err != nil {
			return false, err
		}
	}
	master := false
	taints := node.Spec.Taints
	for _, taint := range taints {
		if taint.Key == "node-role.kubernetes.io/master" {
			master = true
		}
		if taint.Key == "node-role.kubernetes.io/control-plane" {
			master = true
		}
	}
	for _, nodeCond := range node.Status.Conditions {
		if nodeCond.Reason == "KubeletReady" && nodeCond.Type == v1.NodeReady {
			return true, nil
		} else if master && nodeCond.Type == v1.NodeReady {
			return true, nil
		}
	}
	addrs := node.Status.Addresses
	nodeAddr := ""
	for _, addr := range addrs {
		if addr.Type == v1.NodeInternalIP {
			nodeAddr = addr.Address
		}
	}
	logf.Log.Info("Node not ready", nodeName, nodeAddr)
	return false, nil
}

// GetNodeName returns node name corresponding to IP address of a node
func GetNodeName(nodeIp string) (*string, error) {
	nodeLocs, err := getNodeLocs()
	if err != nil {
		return nil, err
	}
	for _, nl := range nodeLocs {
		if nodeIp == nl.IPAddress {
			return &nl.NodeName, nil
		}
	}
	return nil, fmt.Errorf("node with ip %s not found", nodeIp)
}

// GetNodeNqn return mayastor node nqn
func GetNodeNqn(nodeIp string) (string, error) {
	nodeName, err := GetNodeName(nodeIp)
	if err != nil || *nodeName == "" {
		return "", fmt.Errorf("failed to get node name corresponding to node IP %s, error: %v", nodeIp, err)
	}

	node, err := GetMSN(*nodeName)
	if err != nil {
		return "", err
	}
	logf.Log.Info("Node nqn", "GetNodeNqn", node.Spec.Node_nqn)
	return node.Spec.Node_nqn, err
}

func ListIOEngineNodes() (*v1.NodeList, error) {
	nodeList, err := gTestEnv.KubeInt.CoreV1().Nodes().List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return nil, errors.New("failed to list nodes")
	}

	nodeListWithIOEngineLabelPresent := &v1.NodeList{}
	for _, node := range nodeList.Items {
		if value, ok := node.Labels[e2e_config.GetConfig().Product.EngineLabel]; ok && value == e2e_config.GetConfig().Product.EngineLabelValue {
			nodeListWithIOEngineLabelPresent.Items = append(nodeListWithIOEngineLabelPresent.Items, node)
		}
	}
	return nodeListWithIOEngineLabelPresent, nil
}

func GetIOEngineHostNameLabel() (string, error) {
	nodes, err := ListIOEngineNodes()
	if err != nil {
		return "", err
	}
	if len(nodes.Items) == 0 {
		return "", errors.New("io-engine node not found")
	}
	for _, n := range nodes.Items {
		labels := n.Labels
		if _, ok := labels[e2e_config.GetConfig().Product.EngineLabel]; ok {
			return labels[K8sNodeLabelKeyHostname], nil
		}
	}
	return "", nil
}

// ZeroNodeHugePages sets huge pages to 0 on current node
func ZeroNodeHugePages(nodeName string) (string, error) {
	ip, err := GetNodeIPAddress(nodeName)
	if err != nil {
		return "", err
	}
	output, err := e2e_agent.ZeroHugePages(*ip)
	if err != nil {
		return "", err
	}
	msg, errCode, err := e2e_agent.UnwrapResult(output)
	if err != nil || errCode != 0 {
		return "", fmt.Errorf("errCode=%v ; err=%v", errCode, err)
	}
	return msg, nil
}

// ListAllNonMsnNodes list all nodes without io-engine label and without master node
func ListAllNonMsnNodes() (*v1.NodeList, error) {
	nodeList, err := gTestEnv.KubeInt.CoreV1().Nodes().List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return nil, errors.New("failed to list nodes")
	}

	nodeListWithoutIOEngineLabelPresent := &v1.NodeList{}
	for _, node := range nodeList.Items {
		if value, ok := node.Labels[e2e_config.GetConfig().Product.EngineLabel]; ok && value == e2e_config.GetConfig().Product.EngineLabelValue {
			continue
		} else if !strings.Contains(node.Name, "master") {
			nodeListWithoutIOEngineLabelPresent.Items = append(nodeListWithoutIOEngineLabelPresent.Items, node)
		}
	}
	return nodeListWithoutIOEngineLabelPresent, nil
}

func UpdateNodeTaints(nodeName string, taintKey string) error {
	// Get the node object
	node, err := gTestEnv.KubeInt.CoreV1().Nodes().Get(context.TODO(), nodeName, metaV1.GetOptions{})
	if err != nil {
		return fmt.Errorf("error getting node: %v", err)
	}

	// Iterate through the taints and remove the specified one
	var newTaints []coreV1.Taint
	for _, taint := range node.Spec.Taints {
		if taint.Key != taintKey {
			newTaints = append(newTaints, taint)
		}
	}
	node.Spec.Taints = newTaints

	// Update the node object with the new taints
	_, err = gTestEnv.KubeInt.CoreV1().Nodes().Update(context.TODO(), node, metaV1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("error updating node taints: %v", err)
	}

	fmt.Println("Node taints updated successfully")
	return nil
}

func ListWorkerNode() ([]NodeLocation, error) {
	return getNodeLocs()
}

// ListNodeWithoutNoScheduleTaint returns list of nodes which does not have NoSchedule taint
func ListNodeWithoutNoScheduleTaint() ([]string, error) {
	nodeList, err := gTestEnv.KubeInt.CoreV1().Nodes().List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return nil, errors.New("failed to list nodes")
	}
	nodes := make([]string, 0, len(nodeList.Items))
	for _, k8snode := range nodeList.Items {
		var isTaintNoScheduleExist bool
		for _, taint := range k8snode.Spec.Taints {
			if taint.Value == "NoSchedule" {
				isTaintNoScheduleExist = true
				break
			}
		}
		if !isTaintNoScheduleExist {
			nodes = append(nodes, k8snode.Name)
		}
	}
	return nodes, nil
}

// GetLease return requested lease
func GetLease(name, ns string) (*coordinationV1.Lease, error) {
	lease, err := gTestEnv.KubeInt.CoordinationV1().Leases(ns).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, errors.New("failed to get lease")
	}
	return lease, nil
}
