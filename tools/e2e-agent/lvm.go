package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"k8s.io/klog/v2"
)

type Lvm struct {
	Pv                          string `json:"pv"`                          // Physical volume
	Vg                          string `json:"vg"`                          // Volume group
	ThinPoolAutoExtendThreshold int    `json:"thinPoolAutoExtendThreshold"` // thin pool auto extend threshold
	ThinPoolAutoExtendPercent   int    `json:"thinPoolAutoExtendPercent"`   // thin pool auto extend percent
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

	lvmListCommand := "lvm vgs --reportformat json"
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

	lvmListCommand := "lvm pvs --reportformat json"
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

	lvmPvCreateCommand := fmt.Sprintf("lvm pvcreate %s", lvm.Pv)
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

	lvmVgCreateCommand := fmt.Sprintf("lvm vgcreate %s %s", lvm.Vg, lvm.Pv)
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

	lvmPvRemoveCommand := fmt.Sprintf("lvm pvremove %s", lvm.Pv)
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

	lvmVgRemoveCommand := fmt.Sprintf("lvm vgremove %s", lvm.Vg)
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
	lvmConfFile, err := getLvmConfFile()
	if err != nil {
		msg = fmt.Sprintf("lvm conf file verification error, error: %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrFileNotExist, w)
		return
	}
	if lvmConfFile == "" {
		msg := "lvm conf file not found"
		klog.Error(msg)
		WrapResult(msg, ErrFileNotExist, w)
		return
	}
	klog.Info("update %s  thin pool auto extend threshold value, data: %v", lvmConfFile, lvm)

	lvmThinPoolAutoExtentThresholdCommand := fmt.Sprintf("sed -i '/^[^#]*thin_pool_autoextend_threshold/ s/= .*/= %d/' %s",
		lvm.ThinPoolAutoExtendThreshold,
		lvmConfFile)
	output, err := bashLocal(lvmThinPoolAutoExtentThresholdCommand)
	if err != nil {
		msg = fmt.Sprintf("update %s thin pool auto extend threshold value, Error %s", lvmConfFile, err.Error())
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
	if lvm.ThinPoolAutoExtendPercent <= 0 {
		msg = "no thin pool auto extend percent passed"
		klog.Error(msg)
		WrapResult(msg, UnprocessableEntityErrorCode, w)
		return
	}
	lvmConfFile, err := getLvmConfFile()
	if err != nil {
		msg = fmt.Sprintf("lvm conf file verification error, error: %v", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrFileNotExist, w)
		return
	}
	if lvmConfFile == "" {
		msg := "lvm conf file not found"
		klog.Error(msg)
		WrapResult(msg, ErrFileNotExist, w)
		return
	}
	klog.Info("update %s  thin pool auto extend threshold value, data: %v", lvmConfFile, lvm)

	lvmThinPoolAutoExtendPercentCommand := fmt.Sprintf("sed -i '/^[^#]*thin_pool_autoextend_percent/ s/= .*/= %d/' %s",
		lvm.ThinPoolAutoExtendPercent,
		lvmConfFile)
	output, err := bashLocal(lvmThinPoolAutoExtendPercentCommand)
	if err != nil {
		msg = fmt.Sprintf("update %s thin pool auto extend percent value, Error %s", lvmConfFile, err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrExecFailed, w)
		return
	}
	WrapResult(output, ErrNone, w)
}

// If it's Github action based kind cluster, lvm conf file will be
// expected at /host/host/etc/lvm/lvm.conf otherwise it will be at /host/etc/lvm/lvm.conf
func getLvmConfFile() (string, error) {
	ghLvmConf := "/host/host/etc/lvm/lvm.conf"
	isPresent, err := isFilePresent(ghLvmConf)
	if err != nil {
		return "", err
	}
	if isPresent {
		return ghLvmConf, nil
	} else {
		lvmConf := "/host/etc/lvm/lvm.conf"
		isPresent, err := isFilePresent(lvmConf)
		if err != nil {
			return "", err
		}
		if isPresent {
			return lvmConf, nil
		}
	}
	return "", nil
}

func isFilePresent(file string) (bool, error) {
	_, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			klog.Info("file not present", "file", file)
			return false, nil
		}
		return false, fmt.Errorf("error while checking %s file existence, error: %v", file, err)
	}
	klog.Info("file present,", "file:", file)
	return true, nil
}

// LvmLvChangeMonitor monitor lvm lv
func LvmLvChangeMonitor(w http.ResponseWriter, r *http.Request) {
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
	klog.Info("monitor lvm lv, data: %v", lvm)

	lvmLvMonitorCommand := fmt.Sprintf("lvm lvchange --monitor y %s/%s_thinpool", lvm.Vg, lvm.Vg)
	output, err := bashLocal(lvmLvMonitorCommand)
	if err != nil {
		msg = fmt.Sprintf("cannot monitor lvm lv, Error %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrExecFailed, w)
		return
	}
	WrapResult(output, ErrNone, w)
}

// LvmLvRemoveThinPool lv thin pool
func LvmLvRemoveThinPool(w http.ResponseWriter, r *http.Request) {
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
	klog.Info("remove lvm thin pool lv, data: %v", lvm)

	lvmLvMonitorCommand := fmt.Sprintf("lvm lvremove -f --noudevsync %s/%s_thinpool", lvm.Vg, lvm.Vg)
	output, err := bashLocal(lvmLvMonitorCommand)
	if err != nil {
		msg = fmt.Sprintf("cannot remove lvm thin pool lv, Error %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrExecFailed, w)
		return
	}
	WrapResult(output, ErrNone, w)
}
