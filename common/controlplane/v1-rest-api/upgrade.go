package v1_rest_api

import "fmt"

func (cp CPv1RestApi) Upgrade(isUpgradingToUnstableBranch, isPartialRebuildDisableNeeded bool) (string, error) {
	// #TODO implement REST api for upgrade
	panic(fmt.Errorf("not implemented REST api for upgrade"))
}

func (cp CPv1RestApi) UpgradeWithSkipDataPlaneRestart(isUpgradingToUnstableBranch, isPartialRebuildDisableNeeded bool) error {
	// #TODO implement REST api for upgrade
	panic(fmt.Errorf("not implemented REST api for upgrade with skip dataplane restart flag"))
}

func (cp CPv1RestApi) UpgradeWithSkipReplicaRebuild(isUpgradingToUnstableBranch, isPartialRebuildDisableNeeded bool) error {
	//#TODO implement REST api for getting upgrade status
	panic(fmt.Errorf("not implemented REST api for upgrade with --skip-replica-rebuild flag"))
}

func (cp CPv1RestApi) UpgradeWithSkipSingleReplicaValidation(isUpgradingToUnstableBranch, isPartialRebuildDisableNeeded bool) error {
	// #TODO implement REST api for upgrade
	panic(fmt.Errorf("not implemented REST api for upgrade"))
}

func (cp CPv1RestApi) UpgradeWithSkipCordonNodeValidation(isUpgradingToUnstableBranch, isPartialRebuildDisableNeeded bool) error {
	// #TODO implement REST api for upgrade
	panic(fmt.Errorf("not implemented REST api for upgrade with skip cordon node validation flag"))
}

func (cp CPv1RestApi) GetUpgradeStatus() (string, error) {
	//#TODO implement REST api for getting upgrade status
	panic(fmt.Errorf("not implemented REST api for getting upgrade status"))
}

func (cp CPv1RestApi) GetToUpgradeVersion() (string, error) {
	//#TODO implement REST api for getting upgrade status
	panic(fmt.Errorf("not implemented REST api for getting upgrade status"))
}

func (cp CPv1RestApi) DeleteUpgrade() error {
	//#TODO implement REST api for getting upgrade status
	panic(fmt.Errorf("not implemented REST api for deleting upgrade resources"))
}
