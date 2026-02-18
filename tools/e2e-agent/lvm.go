package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"k8s.io/klog/v2"
)

type Lvm struct {
	Pv                          string `json:"pv"`                          // Physical volume
	Vg                          string `json:"vg"`                          // Volume group
	ThinPoolAutoExtendThreshold int    `json:"thinPoolAutoExtendThreshold"` // thin pool auto extend threshold
	ThinPoolAutoExtendPercent   int    `json:"thinPoolAutoExtendPercent"`   // thin pool auto extend percent
}

const (
	thinPoolAutoextendThresholdKey = "thin_pool_autoextend_threshold"
	thinPoolAutoextendPercentKey   = "thin_pool_autoextend_percent"
)

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
	err = updateLVMConfig(lvmConfFile, thinPoolAutoextendThresholdKey, fmt.Sprintf("%d", lvm.ThinPoolAutoExtendThreshold))
	if err != nil {
		msg = fmt.Sprintf("update %s thin pool auto extend threshold value, Error %s", lvmConfFile, err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrExecFailed, w)
		return
	}
	output := fmt.Sprintf("Updated %s thin pool auto extend threshold value to %d", lvmConfFile, lvm.ThinPoolAutoExtendThreshold)
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
	klog.Info("update %s  thin pool auto extend percent value, data: %v", lvmConfFile, lvm)
	err = updateLVMConfig(lvmConfFile, thinPoolAutoextendPercentKey, fmt.Sprintf("%d", lvm.ThinPoolAutoExtendPercent))
	if err != nil {
		msg = fmt.Sprintf("update %s thin pool auto extend percent value, Error %s", lvmConfFile, err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrExecFailed, w)
		return
	}
	output := fmt.Sprintf("Updated %s thin pool auto extend percent value to %d", lvmConfFile, lvm.ThinPoolAutoExtendPercent)
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

// updateLVMConfig updates a specific configuration line in a file.
// It first checks for an uncommented line in the entire file.
// If an uncommented line exists, it updates only that line.
// If no uncommented line is found, it then checks for a commented line,
// updates its value, and uncomments it.
// If neither is found, it appends the new line to the file.
func updateLVMConfig(path, key, value string) error {
	klog.Infof("Updating %s in %s to %s", key, path, value)
	// Read the file content.
	input, err := os.ReadFile(path)
	if err != nil {
		// If the file does not exist, create it and add the new line.
		if os.IsNotExist(err) {
			klog.Info("Configuration file %s does not exist", path)
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}
		// Return any other file reading errors.
		return fmt.Errorf("failed to read file %s: %w", path, err)
	}

	// Split the file content into individual lines.
	lines := strings.Split(string(input), "\n")
	var outputLines []string // Slice to store the modified lines.

	// Regex to match an uncommented line for the key.
	uncommentedRegex := regexp.MustCompile(fmt.Sprintf(`^(\s*)(%s)(\s*=\s*).*$`, regexp.QuoteMeta(key)))
	// Regex to match a commented line for the key.
	commentedRegex := regexp.MustCompile(fmt.Sprintf(`^(\s*)(#)(\s*)(%s)(\s*=\s*).*$`, regexp.QuoteMeta(key)))

	// --- First Pass: Determine if an uncommented key exists anywhere in the file ---
	uncommentedKeyFoundInFile := false
	for _, line := range lines {
		if uncommentedRegex.MatchString(line) {
			uncommentedKeyFoundInFile = true
			break // Found an uncommented line, no need to check further.
		}
	}

	// --- Second Pass: Modify the lines based on the first pass result ---
	lineModified := false // Flag to track if the relevant line has been found and modified in this pass.

	for _, line := range lines {
		// If the target line has already been modified in this pass,
		// append the rest of the lines as they are without further checks for this key.
		if lineModified {
			outputLines = append(outputLines, line)
			continue
		}

		if uncommentedKeyFoundInFile {
			// If an uncommented key exists in the file, only try to update uncommented lines.
			if matches := uncommentedRegex.FindStringSubmatch(line); len(matches) > 0 {
				newLine := matches[1] + matches[2] + matches[3] + value
				outputLines = append(outputLines, newLine)
				lineModified = true // Mark that we've modified the line.
				klog.Infof("Updated uncommented line: %s", newLine)
				continue // Move to the next line.
			}
		} else {
			// If no uncommented key was found in the entire file,
			// then try to update and uncomment a commented line.
			if matches := commentedRegex.FindStringSubmatch(line); len(matches) > 0 {
				newLine := matches[1] + matches[4] + matches[5] + value // Remove '#' and leading whitespace
				outputLines = append(outputLines, newLine)
				lineModified = true // Mark that we've modified the line.
				klog.Infof("Updated and uncommented line: %s", newLine)
				continue // Move to the next line.
			}
		}

		// If no relevant match for the key (based on the logic above), keep the original line.
		outputLines = append(outputLines, line)
	}

	// If after iterating through the entire file, no line for the key was found (neither uncommented nor commented
	// according to the prioritization logic), append the new line to the end of the file.
	if !lineModified {
		newLine := fmt.Sprintf("%s = %s", key, value)
		outputLines = append(outputLines, newLine)
		klog.Infof("Added new line: %s", newLine)
	}

	// Join the modified lines back into a single string, separated by newlines.
	output := strings.Join(outputLines, "\n")

	// Write the modified content back to the file using a temporary file for atomic update.
	dir := filepath.Dir(path)
	tmpFile, err := os.CreateTemp(dir, filepath.Base(path)+".tmp")
	if err != nil {
		return fmt.Errorf("failed to create temporary file in %s: %w", dir, err)
	}
	//nolint:errcheck
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(output); err != nil {
		//nolint:errcheck
		tmpFile.Close()
		return fmt.Errorf("failed to write to temporary file %s: %w", tmpFile.Name(), err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temporary file %s: %w", tmpFile.Name(), err)
	}

	if err := os.Rename(tmpFile.Name(), path); err != nil {
		return fmt.Errorf("failed to rename temporary file %s to %s: %w", tmpFile.Name(), path, err)
	}
	return nil
}
