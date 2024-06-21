package v1

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/openebs/openebs-e2e/common/e2e_config"

	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

// GetClusterRestAPINodeIPs return a set of IP addresses which can be used to communicate
// with bolt/mayastor using the REST API
// returns test control nodes IP addresses
func GetClusterRestAPINodeIPs() ([]string, error) {
	// If environment variable is set that wins, this allows the invoker
	// to avoid invocation of functions which may not work in the deployed
	// environment - for example when run within a cluster context.
	if nodeIPs, ok := os.LookupEnv(`e2e_rest_api_nodes`); ok {
		return strings.Split(nodeIPs, ","), nil
	}

	restConfig := config.GetConfigOrDie()
	if restConfig == nil {
		return nil, fmt.Errorf("failed to create *rest.Config for talking to a Kubernetes apiserver : GetConfigOrDie")
	}
	kubeInt := kubernetes.NewForConfigOrDie(restConfig)
	if kubeInt == nil {
		return nil, fmt.Errorf("failed to create new Clientset for the given config : NewForConfigOrDie")
	}

	nodeList, err := kubeInt.CoreV1().Nodes().List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes %v", err)
	}

	var addrs []string
	// Prefer using control nodes for REST API, they are exempt from disruptive tests
	for _, k8sNode := range nodeList.Items {
		_, ok1 := k8sNode.Labels[e2e_config.GetConfig().TestControlNodeLabel]
		if ok1 {
			for _, addr := range k8sNode.Status.Addresses {
				if addr.Type == coreV1.NodeInternalIP {
					addrs = append(addrs, addr.Address)
				}
			}
		}
	}

	// FIXME: if there are no test control nodes then return all nodes in the cluster, not ideal
	if 0 == len(addrs) {
		for _, k8sNode := range nodeList.Items {
			for _, addr := range k8sNode.Status.Addresses {
				if addr.Type == coreV1.NodeInternalIP {
					addrs = append(addrs, addr.Address)
				}
			}
		}
	}
	return addrs, nil
}
