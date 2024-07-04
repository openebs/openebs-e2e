package v1beta1

import metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type DiskPoolSpec struct {
	Disks []string `json:"disks"`
	Node  string   `json:"node"`
}

type DiskPoolStatus struct {
	Available  uint64 `json:"available"`
	Capacity   uint64 `json:"capacity"`
	Used       uint64 `json:"used"`
	CRStatus   string `json:"cr_state"`
	PoolStatus string `json:"pool_status"`
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
