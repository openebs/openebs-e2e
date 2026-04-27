package v1beta3

import metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type Topology struct {
	Labelled map[string]string `json:"labelled,omitempty"`
}

type Secret struct {
	Name string `json:"name,omitempty"`
}

type Source struct {
	Secret Secret `json:"secret,omitempty"`
}

type EncryptionConfig struct {
	Source Source `json:"source,omitempty"`
}

type DiskPoolSpec struct {
	Disks            []string          `json:"disks"`
	Node             string            `json:"node"`
	Topology         *Topology         `json:"topology,omitempty"`
	EncryptionConfig *EncryptionConfig `json:"encryptionConfig,omitempty"`
	ClusterSize      string            `json:"clusterSize,omitempty"`
	// MaxExpansion is an optional limit for pool expansion.
	// Supported formats include absolute sizes (e.g. "20GiB") or
	// multiplicative factors (e.g. "2x", "1.2x").
	// This field is optional and omitting it preserves existing behavior.
	MaxExpansion string `json:"maxExpansion,omitempty"`
}

type Condition struct {
	Status string `json:"status"`
	Reason string `json:"reason,omitempty"`
}

type DiskPoolDiagErrors struct {
	IOErrors uint64 `json:"io_errors,omitempty"`
	Status   string `json:"status,omitempty"`
}

type DiskPoolDiag struct {
	Errors DiskPoolDiagErrors `json:"errors,omitempty"`
	Error  DiskPoolError      `json:"error,omitempty"`
}

type DiskPoolStatus struct {
	Available uint64 `json:"available"`
	Capacity  uint64 `json:"capacity"`
	Used      uint64 `json:"used"`
	CRStatus  string `json:"cr_state"`
	// Status is a generic pool state field used by some control-plane versions.
	// When present, it may carry values like "Offline" even if PoolStatus is empty.
	Status            string            `json:"status,omitempty"`
	PoolStatus        string            `json:"pool_status"`
	Encrypted         bool              `json:"encrypted"`
	CapacityQ         string            `json:"capacity_q"`
	AvailableQ        string            `json:"available_q"`
	UsedQ             string            `json:"used_q"`
	ClusterSize       string            `json:"clusterSize"`
	MaxExpandableSize string            `json:"maxExpandableSize,omitempty"`
	Conditions        []Condition       `json:"conditions,omitempty"`
	Diag              DiskPoolDiag      `json:"diag,omitempty"`
	Error             DiskPoolError     `json:"error,omitempty"`
	ErrorInfo         DiskPoolErrorInfo `json:"errorInfo,omitempty"`
}

type DiskPoolError struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type DiskPoolErrorInfo struct {
	IoStallTransitionCount uint64         `json:"ioStallTransitionCount,omitempty"`
	IoStalled              bool           `json:"ioStalled,omitempty"`
	Alerts                 DiskPoolAlerts `json:"alerts,omitempty"`
	IoErrorCount           uint64         `json:"ioErrorCount,omitempty"`
}

type DiskPool struct {
	metaV1.TypeMeta   `json:",inline"`
	metaV1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DiskPoolSpec   `json:"spec"`
	Status DiskPoolStatus `json:"status"`
}

type DiskPoolList struct {
	metaV1.TypeMeta `json:",inline"`
	metaV1.ListMeta `json:"metadata,omitempty"`

	Items []DiskPool `json:"items"`
}

type DiskPoolAlerts struct {
	Status    string   `json:"status,omitempty"`
	Notice    []string `json:"notice,omitempty"`
	Attention []string `json:"attention,omitempty"`
	Warning   []string `json:"warning,omitempty"`
	Critical  []string `json:"critical,omitempty"`
}
