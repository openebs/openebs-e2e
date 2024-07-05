package topology

import (
	"fmt"
	"time"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/e2e_config"
	"github.com/openebs/openebs-e2e/common/k8stest"
)

const (
	g_pollIntervalSeconds      = 5
	g_podAbsenceTimeoutSeconds = 60
)

var CPDeployment = []string{
	e2e_config.GetConfig().Product.ControlPlanePoolOperator,
	e2e_config.GetConfig().Product.ControlPlaneRestServer,
	e2e_config.GetConfig().Product.ControlPlaneCoreAgent,
	e2e_config.GetConfig().Product.ControlPlaneCsiController,
	e2e_config.GetConfig().Product.ControlPlaneLocalpvProvisioner,
	e2e_config.GetConfig().Product.ControlPlaneObsCallhome,
}

func MigrateControlPlane(fromNode string, toNode string) error {
	// reschedule control plane components to different node
	var err error

	for _, deploy := range CPDeployment {
		isPresent, err := k8stest.PodPresentOnNode(deploy, common.NSMayastor(), fromNode)
		if err != nil {
			return fmt.Errorf("failed to verify pod presence %v", err)
		}
		if isPresent {
			err = k8stest.ApplyNodeSelectorToDeployment(deploy, common.NSMayastor(), "kubernetes.io/hostname", toNode)
			if err != nil {
				return fmt.Errorf("failed to apply node selector %v", err)
			}
			endTime := time.Now().Add(g_podAbsenceTimeoutSeconds * time.Second)
			for endTime.After(time.Now()) {
				isPresent, err = k8stest.PodPresentOnNode(deploy, common.NSMayastor(), fromNode)
				if err != nil {
					return fmt.Errorf("failed to verify pod presence %v", err)
				}
				if !isPresent {
					break
				}
				time.Sleep(g_pollIntervalSeconds * time.Second)
			}
			if isPresent {
				return fmt.Errorf("failed to reschedule deployment %s %v", deploy, err)
			}
		}
	}

	// verify mayastor ready check
	ready, err := k8stest.MayastorReady(2, 360)
	if err != nil {
		return err
	}
	if !ready {
		return fmt.Errorf("mayastor installation not ready")
	}
	return err
}

func RemoveControlPlaneSelector() error {
	var err error

	for _, deploy := range CPDeployment {
		err = k8stest.RemoveAllNodeSelectorsFromDeployment(deploy, common.NSMayastor())
		if err != nil {
			return fmt.Errorf("failed to remove node selector %v", err)
		}
	}

	// verify mayastor ready check, in case anything starts moving back
	ready, err := k8stest.MayastorReady(2, 360)
	if err != nil {
		return err
	}
	if !ready {
		return fmt.Errorf("mayastor installation not ready")
	}
	return err
}
