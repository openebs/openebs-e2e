package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"k8s.io/klog/v2"
)

type Lvm struct {
	Pv                          string `json:"pv"`                          // Physical volume
	Vg                          string `json:"vg"`                          // Volume group
	ThinPoolAutoExtendThreshold int    `json:"thinPoolAutoExtendThreshold"` // thin pool auto extend threshold
	ThinPoolAutoExtendPercent   int    `json:"ThinPoolAutoExtendPercent"`   // thin pool auto extend percent
}

// LvmVersion check lvm version installed on node
func LvmVersion(w http.ResponseWriter, r *http.Request) {
	var msg string
	klog.Info("Get lvm version installed")

	lvmVersionCommand := "lvm version"
	output, err := bashLocal(lvmVersionCommand)
	if err != nil {
		msg = fmt.Sprintf("cannot get lvm version. Error %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrExecFailed, w)
		return
	}
	WrapResult(output, ErrNone, w)
}

// LvmListVg list lvm vg
func LvmListVg(w http.ResponseWriter, r *http.Request) {
	var msg string
	klog.Info("List lvm vgs")

	lvmListCommand := "vgs --reportformat json"
	output, err := bashLocal(lvmListCommand)
	if err != nil {
		msg = fmt.Sprintf("cannot list lvm vgs Error %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrExecFailed, w)
		return
	}
	WrapResult(output, ErrNone, w)
}

// LvmListPv list lvm pv
func LvmListPv(w http.ResponseWriter, r *http.Request) {
	var msg string
	klog.Info("List lvm pvs")

	lvmListCommand := "pvs --reportformat json"
	output, err := bashLocal(lvmListCommand)
	if err != nil {
		msg = fmt.Sprintf("cannot list lvm pvs Error %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrExecFailed, w)
		return
	}
	WrapResult(output, ErrNone, w)
}

// LvmCreatePv create lvm pv
func LvmCreatePv(w http.ResponseWriter, r *http.Request) {
	var msg string
	var lvm Lvm
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&lvm); err != nil {
		msg = fmt.Sprintf("failed to read JSON encoded data, Error: %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrJsonDecode, w)
		return
	}
	if lvm.Pv == "" {
		msg = "no device path passed"
		klog.Error(msg)
		WrapResult(msg, UnprocessableEntityErrorCode, w)
		return
	}
	klog.Info("creates lvm pv, data: %v", lvm)

	lvmPvCreateCommand := fmt.Sprintf("pvcreate %s", lvm.Pv)
	output, err := bashLocal(lvmPvCreateCommand)
	if err != nil {
		msg = fmt.Sprintf("cannot create lvm pv Error %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrExecFailed, w)
		return
	}
	WrapResult(output, ErrNone, w)
}

// LvmCreateVg create lvm vg
func LvmCreateVg(w http.ResponseWriter, r *http.Request) {
	var msg string
	var lvm Lvm
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&lvm); err != nil {
		msg = fmt.Sprintf("failed to read JSON encoded data, Error: %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrJsonDecode, w)
		return
	}
	if lvm.Pv == "" {
		msg = "no device path passed"
		klog.Error(msg)
		WrapResult(msg, UnprocessableEntityErrorCode, w)
		return
	}
	if lvm.Vg == "" {
		msg = "no vg name passed"
		klog.Error(msg)
		WrapResult(msg, UnprocessableEntityErrorCode, w)
		return
	}
	klog.Info("creates lvm vg, data: %v", lvm)

	lvmVgCreateCommand := fmt.Sprintf("vgcreate %s %s", lvm.Vg, lvm.Pv)
	output, err := bashLocal(lvmVgCreateCommand)
	if err != nil {
		msg = fmt.Sprintf("cannot create lvm vg Error %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrExecFailed, w)
		return
	}
	WrapResult(output, ErrNone, w)
}

// LvmRemovePv remove lvm pv
func LvmRemovePv(w http.ResponseWriter, r *http.Request) {
	var msg string
	var lvm Lvm
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&lvm); err != nil {
		msg = fmt.Sprintf("failed to read JSON encoded data, Error: %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrJsonDecode, w)
		return
	}
	if lvm.Pv == "" {
		msg = "no device path passed"
		klog.Error(msg)
		WrapResult(msg, UnprocessableEntityErrorCode, w)
		return
	}
	klog.Info("remove lvm pv, data: %v", lvm)

	lvmPvRemoveCommand := fmt.Sprintf("pvremove %s", lvm.Pv)
	output, err := bashLocal(lvmPvRemoveCommand)
	if err != nil {
		msg = fmt.Sprintf("cannot remove lvm pv Error %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrExecFailed, w)
		return
	}
	WrapResult(output, ErrNone, w)
}

// LvmRemoveVg remove lvm vg
func LvmRemoveVg(w http.ResponseWriter, r *http.Request) {
	var msg string
	var lvm Lvm
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&lvm); err != nil {
		msg = fmt.Sprintf("failed to read JSON encoded data, Error: %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrJsonDecode, w)
		return
	}
	if lvm.Vg == "" {
		msg = "no vg name passed"
		klog.Error(msg)
		WrapResult(msg, UnprocessableEntityErrorCode, w)
		return
	}
	klog.Info("remove lvm vg, data: %v", lvm)

	lvmVgRemoveCommand := fmt.Sprintf("vgremove %s", lvm.Vg)
	output, err := bashLocal(lvmVgRemoveCommand)
	if err != nil {
		msg = fmt.Sprintf("cannot create lvm vg Error %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrExecFailed, w)
		return
	}
	WrapResult(output, ErrNone, w)
}

// LvmThinPoolAutoExtendThreshold update lvm.conf thin pool auto extend threshold value
func LvmThinPoolAutoExtendThreshold(w http.ResponseWriter, r *http.Request) {
	var msg string
	var lvm Lvm
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&lvm); err != nil {
		msg = fmt.Sprintf("failed to read JSON encoded data, Error: %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrJsonDecode, w)
		return
	}
	if lvm.ThinPoolAutoExtendThreshold <= 0 {
		msg = "no thin pool auto extent threshold passed"
		klog.Error(msg)
		WrapResult(msg, UnprocessableEntityErrorCode, w)
		return
	}
	klog.Info("update lvm.conf thin pool auto extend threshold value, data: %v", lvm)

	lvmThinPoolAutoExtentThresholdCommand := fmt.Sprintf("sed -i '/^[^#]*thin_pool_autoextend_threshold/ s/= .*/= %d/' /etc/lvm/lvm.conf",
		lvm.ThinPoolAutoExtendThreshold)
	output, err := bashLocal(lvmThinPoolAutoExtentThresholdCommand)
	if err != nil {
		msg = fmt.Sprintf("update lvm.conf thin pool auto extend threshold value, Error %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrExecFailed, w)
		return
	}
	WrapResult(output, ErrNone, w)
}

// LvmThinPoolAutoExtendPercent update lvm.conf thin pool auto extend percent value
func LvmThinPoolAutoExtendPercent(w http.ResponseWriter, r *http.Request) {
	var msg string
	var lvm Lvm
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&lvm); err != nil {
		msg = fmt.Sprintf("failed to read JSON encoded data, Error: %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrJsonDecode, w)
		return
	}
	if lvm.ThinPoolAutoExtendThreshold <= 0 {
		msg = "no thin pool auto extent percent passed"
		klog.Error(msg)
		WrapResult(msg, UnprocessableEntityErrorCode, w)
		return
	}
	klog.Info("update lvm.conf thin pool auto extend threshold value, data: %v", lvm)

	lvmThinPoolAutoExtentPercentCommand := fmt.Sprintf("sed -i '/^[^#]*thin_pool_autoextend_percent/ s/= .*/= %d/' /etc/lvm/lvm.conf",
		lvm.ThinPoolAutoExtendPercent)
	output, err := bashLocal(lvmThinPoolAutoExtentPercentCommand)
	if err != nil {
		msg = fmt.Sprintf("update lvm.conf thin pool auto extend percent value, Error %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrExecFailed, w)
		return
	}
	WrapResult(output, ErrNone, w)
}
