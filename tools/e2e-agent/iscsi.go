package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"k8s.io/klog/v2"
)

// iscsiadmSessionAlreadyExists is iscsiadm's exit code for "session already
// exists" on login -- a benign login race, not a failure. Matches
// csi.sansymphony.datacore.com's own src/dev/iscsiadm.rs ERR_SESS_EXISTS.
const iscsiadmSessionAlreadyExists = 15

// iscsiadmNoObjectsFound is iscsiadm's exit code for "no active sessions" /
// "no records found" -- an empty result, not a failure. Matches
// src/dev/iscsiadm.rs ERR_NO_OBJS_FOUND.
const iscsiadmNoObjectsFound = 21

// Iscsi carries the parameters for an iscsiadm discovery/login/logout call.
type Iscsi struct {
	Portal    string `json:"portal"`    // "address:port", e.g. "10.0.0.5:3260"
	TargetIqn string `json:"targetIqn"` // required for login/logout
}

// runOnHost runs a shell command inside the host's root filesystem via
// chroot, so it uses the host's own copy of a tool (iscsiadm, multipath)
// rather than requiring the agent image to bundle it -- the same approach
// server.go already uses for dmsetup and kernel modules. Unlike those,
// this always execs through "bash -c" so callers can use shell features
// (pipes, redirection) in shellCmd.
func runOnHost(shellCmd string) (string, error) {
	cmd := exec.Command("chroot", "/host", "bash", "-c", shellCmd)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

// runOnHostExitCode is runOnHost plus the process's exit code, for callers
// that need to distinguish specific iscsiadm exit codes (15 "session
// exists", 21 "no active sessions") from a real failure rather than
// pattern-matching the message text. -1 when the process could not be
// started at all.
func runOnHostExitCode(shellCmd string) (output string, exitCode int, err error) {
	cmd := exec.Command("chroot", "/host", "bash", "-c", shellCmd)
	out, runErr := cmd.CombinedOutput()
	output = strings.TrimSpace(string(out))
	if runErr == nil {
		return output, 0, nil
	}
	var exitErr *exec.ExitError
	if errors.As(runErr, &exitErr) {
		return output, exitErr.ExitCode(), runErr
	}
	return output, -1, runErr
}

// GetIscsiInitiatorName returns the node's iSCSI initiator IQN from
// /etc/iscsi/initiatorname.iscsi.
func GetIscsiInitiatorName(w http.ResponseWriter, r *http.Request) {
	output, err := runOnHost(`grep '^InitiatorName=' /etc/iscsi/initiatorname.iscsi 2>/dev/null | cut -d= -f2-`)
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, output)
		klog.Error("failed to read iSCSI initiator name, Error: ", err)
		return
	}
	_, _ = fmt.Fprint(w, output)
	klog.Info("iSCSI initiator name: ", output)
}

// GetNvmeHostNqn returns the node's NVMe host NQN from /etc/nvme/hostnqn.
func GetNvmeHostNqn(w http.ResponseWriter, r *http.Request) {
	output, err := runOnHost(`cat /etc/nvme/hostnqn 2>/dev/null`)
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, output)
		klog.Error("failed to read NVMe host NQN, Error: ", err)
		return
	}
	_, _ = fmt.Fprint(w, output)
	klog.Info("NVMe host NQN: ", output)
}

// IscsiadmDiscovery runs sendtargets discovery against a portal
// ("address:port"), listing the targets it offers.
func IscsiadmDiscovery(w http.ResponseWriter, r *http.Request) {
	var iscsi Iscsi
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&iscsi); err != nil {
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
		return
	}
	if iscsi.Portal == "" {
		w.WriteHeader(UnprocessableEntityErrorCode)
		_, _ = fmt.Fprint(w, "no portal passed")
		klog.Error("no portal passed")
		return
	}
	cmdStr := fmt.Sprintf("iscsiadm -m discovery -t sendtargets -p %s", iscsi.Portal)
	klog.Info("running command ", cmdStr)
	output, err := runOnHost(cmdStr)
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, output)
		klog.Error("failed to run iscsiadm discovery against ", iscsi.Portal, "Error: ", err)
		return
	}
	_, _ = fmt.Fprint(w, output)
	klog.Info(output)
}

// IscsiadmLogin logs in to a discovered target.
func IscsiadmLogin(w http.ResponseWriter, r *http.Request) {
	var iscsi Iscsi
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&iscsi); err != nil {
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
		return
	}
	if iscsi.Portal == "" || iscsi.TargetIqn == "" {
		w.WriteHeader(UnprocessableEntityErrorCode)
		_, _ = fmt.Fprint(w, "no portal or targetIqn passed")
		klog.Error("no portal or targetIqn passed")
		return
	}
	cmdStr := fmt.Sprintf("iscsiadm -m node -T %s -p %s --login", iscsi.TargetIqn, iscsi.Portal)
	klog.Info("running command ", cmdStr)
	output, exitCode, err := runOnHostExitCode(cmdStr)
	if err != nil && exitCode != iscsiadmSessionAlreadyExists {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, output)
		klog.Error("failed to log in to ", iscsi.TargetIqn, " via ", iscsi.Portal, "Error: ", err)
		return
	}
	if exitCode == iscsiadmSessionAlreadyExists {
		klog.Info("session to ", iscsi.TargetIqn, " via ", iscsi.Portal, " already exists")
	}
	_, _ = fmt.Fprint(w, output)
	klog.Info(output)
}

// IscsiadmRescan rescans one target's session for new or resized LUNs
// (`iscsiadm -m node -T <iqn> -p <portal> -R`).
func IscsiadmRescan(w http.ResponseWriter, r *http.Request) {
	var iscsi Iscsi
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&iscsi); err != nil {
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
		return
	}
	if iscsi.Portal == "" || iscsi.TargetIqn == "" {
		w.WriteHeader(UnprocessableEntityErrorCode)
		_, _ = fmt.Fprint(w, "no portal or targetIqn passed")
		klog.Error("no portal or targetIqn passed")
		return
	}
	cmdStr := fmt.Sprintf("iscsiadm -m node -T %s -p %s -R", iscsi.TargetIqn, iscsi.Portal)
	klog.Info("running command ", cmdStr)
	output, err := runOnHost(cmdStr)
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, output)
		klog.Error("failed to rescan ", iscsi.TargetIqn, " via ", iscsi.Portal, "Error: ", err)
		return
	}
	_, _ = fmt.Fprint(w, output)
	klog.Info(output)
}

// IscsiadmLogout logs out of a target.
func IscsiadmLogout(w http.ResponseWriter, r *http.Request) {
	var iscsi Iscsi
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&iscsi); err != nil {
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
		return
	}
	if iscsi.Portal == "" || iscsi.TargetIqn == "" {
		w.WriteHeader(UnprocessableEntityErrorCode)
		_, _ = fmt.Fprint(w, "no portal or targetIqn passed")
		klog.Error("no portal or targetIqn passed")
		return
	}
	cmdStr := fmt.Sprintf("iscsiadm -m node -T %s -p %s --logout", iscsi.TargetIqn, iscsi.Portal)
	klog.Info("running command ", cmdStr)
	output, err := runOnHost(cmdStr)
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, output)
		klog.Error("failed to log out of ", iscsi.TargetIqn, " via ", iscsi.Portal, "Error: ", err)
		return
	}
	_, _ = fmt.Fprint(w, output)
	klog.Info(output)
}

// IscsiSessions lists active iSCSI sessions (`iscsiadm -m session`).
// iscsiadm exits with code 21 ("no active sessions" / "no records found")
// when there are none -- reported here as an empty success rather than an
// error, matching csi.sansymphony.datacore.com's own
// src/dev/iscsiadm.rs sessions(), which treats the same code the same way.
func IscsiSessions(w http.ResponseWriter, r *http.Request) {
	output, exitCode, err := runOnHostExitCode("iscsiadm -m session")
	if err != nil && exitCode == iscsiadmNoObjectsFound {
		_, _ = fmt.Fprint(w, "")
		klog.Info("no active iSCSI sessions")
		return
	}
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, output)
		klog.Error("failed to list iSCSI sessions, Error: ", err)
		return
	}
	_, _ = fmt.Fprint(w, output)
	klog.Info(output)
}
