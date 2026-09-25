package v1_rest_api

import (
	"fmt"
	"time"

	"github.com/openebs/openebs-e2e/common"
	cpV1 "github.com/openebs/openebs-e2e/common/controlplane/v1"
)

// pool drain - not implemented via direct REST calls, same as pool cordon.
// TODO: implement once the drain endpoints and generated client types exist.

func (cp CPv1RestApi) DrainPool(poolID string, unsafeRebuildOtherwiseEvict *time.Duration, opts ...common.PoolDrainOption) error {
	// #TODO implement REST api for pool drain
	panic(fmt.Errorf("not implemented REST api for pool drain"))
}

func (cp CPv1RestApi) AbortPoolDrain(poolID string) error {
	// #TODO implement REST api for aborting a pool drain
	panic(fmt.Errorf("not implemented REST api for aborting a pool drain"))
}

func (cp CPv1RestApi) GetPoolDrainProgress(poolID string) (*cpV1.PoolDrainRecord, error) {
	// #TODO implement REST api for getting pool drain progress
	panic(fmt.Errorf("not implemented REST api for getting pool drain progress"))
}

func (cp CPv1RestApi) GetPoolLiveUsage(poolID string) (*cpV1.PoolDrainUsage, error) {
	// #TODO implement REST api for getting a pool's live usage
	panic(fmt.Errorf("not implemented REST api for getting a pool's live usage"))
}
