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
	ClusterSize      string            `json:"cluster_size,omitempty"`
}

type DiskPoolStatus struct {
	Available   uint64 `json:"available"`
	Capacity    uint64 `json:"capacity"`
	Used        uint64 `json:"used"`
	CRStatus    string `json:"cr_state"`
	PoolStatus  string `json:"pool_status"`
	Encrypted   bool   `json:"encrypted"`
	CapacityQ   string `json:"capacity_q"`
	AvailableQ  string `json:"available_q"`
	UsedQ       string `json:"used_q"`
	ClusterSize string `json:"cluster_size"`
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
