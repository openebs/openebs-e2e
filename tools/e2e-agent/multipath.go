package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"k8s.io/klog/v2"
)

// BlockDeviceQuery names a block device by its kernel name, e.g. "sdb" or
// "dm-3" -- not a full /dev path.
type BlockDeviceQuery struct {
	Device string `json:"device"`
}

// GetDeviceWwid returns the wwid backing a block device, read from sysfs
// directly (kernel-global, no chroot needed) -- the identity a CSI node
// plugin verifies a resolved device against, per bolt-specifications PR
// #133's "Devices are located by identity, never by name" scenario. Tries
// both locations the kernel is known to expose it under:
// /sys/block/<dev>/wwid (dm devices) and /sys/block/<dev>/device/wwid
// (scsi disks).
func GetDeviceWwid(w http.ResponseWriter, r *http.Request) {
	var q BlockDeviceQuery
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&q); err != nil {
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
		return
	}
	if q.Device == "" {
		w.WriteHeader(UnprocessableEntityErrorCode)
		_, _ = fmt.Fprint(w, "no device passed")
		klog.Error("no device passed")
		return
	}
	cmdStr := fmt.Sprintf(
		"cat /sys/block/%s/wwid 2>/dev/null || cat /sys/block/%s/device/wwid 2>/dev/null",
		q.Device, q.Device,
	)
	klog.Info("running command ", cmdStr)
	cmd := exec.Command("bash", "-c", cmdStr)
	output, err := cmd.CombinedOutput()
	trimmed := strings.TrimSpace(string(output))
	if err != nil || trimmed == "" {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprintf(w, "no wwid found for device %s", q.Device)
		klog.Error("failed to read wwid for ", q.Device, "Error: ", err)
		return
	}
	_, _ = fmt.Fprint(w, trimmed)
	klog.Info("wwid for ", q.Device, ": ", trimmed)
}

// MultipathList runs `multipath -ll`, listing every assembled dm-multipath
// map and its member paths.
func MultipathList(w http.ResponseWriter, r *http.Request) {
	output, err := runOnHost("multipath -ll")
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, output)
		klog.Error("failed to list multipath maps, Error: ", err)
		return
	}
	_, _ = fmt.Fprint(w, output)
	klog.Info(output)
}

// ListDMDevices lists every /sys/block/dm-* device with its device-mapper
// uuid and name, one "dm-N\tuuid\tname" line each -- the same sysfs fields
// csi.sansymphony.datacore.com's own src/dev/multipath.rs find_map/
// map_name read to locate a wwid's map (by matching dm/uuid against
// "mpath-3<hex>") and to address it via multipathd (by dm/name, which is
// the wwid-derived name only when user_friendly_names is off). This is
// ground truth the driver itself relies on, independent of `multipath -ll`
// output formatting.
func ListDMDevices(w http.ResponseWriter, r *http.Request) {
	script := `for d in /sys/block/dm-*; do
  [ -d "$d" ] || continue
  name=$(basename "$d")
  uuid=$(cat "$d/dm/uuid" 2>/dev/null)
  dmname=$(cat "$d/dm/name" 2>/dev/null)
  printf '%s\t%s\t%s\n' "$name" "$uuid" "$dmname"
done`
	cmd := exec.Command("bash", "-c", script)
	output, err := cmd.CombinedOutput()
	trimmed := strings.TrimSpace(string(output))
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, trimmed)
		klog.Error("failed to list dm devices, Error: ", err)
		return
	}
	_, _ = fmt.Fprint(w, trimmed)
	klog.Info("dm devices: ", trimmed)
}

// MultipathStatus runs `multipath -l <wwidOrMapName>`, reporting one map's
// status and paths -- empty output means no map exists for it, which is
// the condition "The multipath map is waited for and its own wwid
// verified" polls on.
func MultipathStatus(w http.ResponseWriter, r *http.Request) {
	var q BlockDeviceQuery
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&q); err != nil {
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
		return
	}
	if q.Device == "" {
		w.WriteHeader(UnprocessableEntityErrorCode)
		_, _ = fmt.Fprint(w, "no device passed")
		klog.Error("no device passed")
		return
	}
	cmdStr := fmt.Sprintf("multipath -l %s", q.Device)
	klog.Info("running command ", cmdStr)
	output, err := runOnHost(cmdStr)
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, output)
		klog.Error("failed to get multipath status for ", q.Device, "Error: ", err)
		return
	}
	_, _ = fmt.Fprint(w, output)
	klog.Info(output)
}
