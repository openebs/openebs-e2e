package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"k8s.io/klog/v2"
)

// FcHost names one SCSI host to rescan, e.g. "host3" (the basename under
// /sys/class/scsi_host). Channel, Target and Lun make it a targeted scan
// ("<channel> <target> <lun>" written to the host's scan attribute,
// resolved from a remote port -- see ListFcRemotePorts); left empty (or
// "-") each defaults to "-", the SCSI wildcard, matching the driver's own
// fallback when no remote port is visible yet.
type FcHost struct {
	ScsiHost string `json:"scsiHost"`
	Channel  string `json:"channel"`
	Target   string `json:"target"`
	Lun      string `json:"lun"`
}

// GetFcWwpns lists this node's Fibre Channel WWPNs, one per line, from
// online FC hosts only (/sys/class/fc_host/host*/port_name, filtered by
// port_state == "Online") -- matching csi.sansymphony.datacore.com's own
// src/dev/fc.rs local_wwpns, which skips a Linkdown host rather than
// reporting a WWPN nothing can register a port for. sysfs is kernel-global,
// so this needs no chroot, the same reasoning server.go's getDeviceState
// already relies on for /sys/block reads.
func GetFcWwpns(w http.ResponseWriter, r *http.Request) {
	script := `for h in /sys/class/fc_host/host*; do
  [ -r "$h/port_state" ] || continue
  state=$(cat "$h/port_state" 2>/dev/null)
  [ "$state" = "Online" ] || continue
  cat "$h/port_name" 2>/dev/null
done`
	cmd := exec.Command("bash", "-c", script)
	output, err := cmd.CombinedOutput()
	trimmed := strings.TrimSpace(string(output))
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, trimmed)
		klog.Error("failed to read Fibre Channel WWPNs, Error: ", err)
		return
	}
	_, _ = fmt.Fprint(w, trimmed)
	klog.Info("Fibre Channel WWPNs (online hosts): ", trimmed)
}

// FcRemotePort is one entry under /sys/class/fc_remote_ports, with the
// scsi_host/channel it resolves to -- what a caller needs to build the
// targeted "<channel> <target> <lun>" scan FcRescanHost takes, the same
// resolution csi.sansymphony.datacore.com's own src/dev/fc.rs rescan
// performs before falling back to a wildcard scan.
type FcRemotePort struct {
	Rport        string `json:"rport"`        // e.g. "rport-7:0-3"
	PortName     string `json:"portName"`     // raw WWPN as sysfs reports it
	ScsiTargetID string `json:"scsiTargetId"` // "-1" when not yet a SCSI target
	ScsiHost     string `json:"scsiHost"`     // "host7", from the rport name
	ScsiChannel  string `json:"scsiChannel"`  // from the rport name
}

// ListFcRemotePorts lists every /sys/class/fc_remote_ports/rport-* entry.
func ListFcRemotePorts(w http.ResponseWriter, r *http.Request) {
	script := `for rp in /sys/class/fc_remote_ports/rport-*; do
  [ -d "$rp" ] || continue
  name=$(basename "$rp")
  port_name=$(cat "$rp/port_name" 2>/dev/null)
  target_id=$(cat "$rp/scsi_target_id" 2>/dev/null)
  # rport-H:C-I -> host H, channel C.
  hc=${name#rport-}
  host=${hc%%:*}
  rest=${hc#*:}
  channel=${rest%%-*}
  printf '%s\t%s\t%s\thost%s\t%s\n' "$name" "$port_name" "$target_id" "$host" "$channel"
done`
	cmd := exec.Command("bash", "-c", script)
	output, err := cmd.CombinedOutput()
	trimmed := strings.TrimSpace(string(output))
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, trimmed)
		klog.Error("failed to list FC remote ports, Error: ", err)
		return
	}
	_, _ = fmt.Fprint(w, trimmed)
	klog.Info("FC remote ports: ", trimmed)
}

// FcRescanHost writes "<channel> <target> <lun>" to one SCSI host's scan
// sysfs attribute -- a targeted scan when Channel/Target/Lun are set
// (resolve them first via ListFcRemotePorts), or the wildcard "- - -" when
// they are left empty, matching the driver's own fallback for a remote
// port that is not visible yet.
func FcRescanHost(w http.ResponseWriter, r *http.Request) {
	var fc FcHost
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&fc); err != nil {
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
		return
	}
	if fc.ScsiHost == "" {
		w.WriteHeader(UnprocessableEntityErrorCode)
		_, _ = fmt.Fprint(w, "no scsiHost passed")
		klog.Error("no scsiHost passed")
		return
	}
	channel, target, lun := fc.Channel, fc.Target, fc.Lun
	if channel == "" {
		channel = "-"
	}
	if target == "" {
		target = "-"
	}
	if lun == "" {
		lun = "-"
	}
	path := fmt.Sprintf("/sys/class/scsi_host/%s/scan", fc.ScsiHost)
	cmdStr := fmt.Sprintf(`echo "%s %s %s" > %s`, channel, target, lun, path)
	klog.Info("running command ", cmdStr)
	cmd := exec.Command("bash", "-c", cmdStr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, string(output))
		klog.Error("failed to rescan SCSI host ", fc.ScsiHost, "Error: ", err)
		return
	}
	_, _ = fmt.Fprint(w, "rescanned "+fc.ScsiHost)
	klog.Info("rescanned SCSI host ", fc.ScsiHost, " with ", channel, target, lun)
}
