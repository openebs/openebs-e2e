package disk_failures

import (
	"errors"
	"time"
)

var (
	ErrPoolNil            = errors.New("pool CR is nil")
	ErrBadCRStatus        = errors.New("cr status is not creating")
	ErrBadReadyStatus     = errors.New("pool ready status is not false")
	ErrBadReason          = errors.New("pool ready reason is not diskioerror")
	ErrUnknownStatus      = errors.New("pool status is unknown")
	ErrPoolIOErrorCount   = errors.New("pool error count is not greater than 0")
	ErrPoolAlertAttention = errors.New("pool alert status is not attention")
	ErrPoolImport         = errors.New("pool error code is not pool import error")
	ErrPoolOffline        = errors.New("pool status is not offline")
	ErrPoolOnline         = errors.New("pool status is not online")
	ErrPoolSuspected      = errors.New("pool status is not suspected")
	ErrPoolAlertCritical  = errors.New("pool alert status is not critical")
	ErrPoolAlertHealthy   = errors.New("pool alert status is not healthy")
)

const (
	DiskPoolCRCheckTimeout  = 120 * time.Second
	DiskPoolCRCheckInterval = 5 * time.Second
	NodeReadyTimeout        = 5 * time.Minute
	NodeReadyInterval       = 10 * time.Second
	NodeRebootWait          = 30 * time.Second
	IoEngineTimeout         = 120
	ThresholdPath           = "io_engine.pool.alerts.errorThreshold"
	VolumeStateTimeout      = 120
	VolumeStateInterval     = "5s"
	RebuildTimeout          = 240
	StallDeadlineHelmValue  = "io_engine.pool.alerts.stallDeadline"
	StallDeadline           = "10s"
	PoolStateOnline         = "Online"
	PoolStateOffline        = "Offline"
	PoolStateSuspected      = "Suspected"
	PoolAlertStatusCritical = "Critical"
	PoolAlertStatusHealthy  = "Healthy"
)
