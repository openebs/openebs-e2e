package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"k8s.io/klog/v2"
)

// CreateHostPathDisk creates hostpath disk
func CreateHostPathDisk(w http.ResponseWriter, r *http.Request) {
	var device LoopDevice
	var msg string

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&device); err != nil {
		msg = fmt.Sprintf("failed to read JSON encoded data, Error: %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrJsonDecode, w)
		return
	}
	klog.Info("creates a hostpath device, data: %v", device)
	if device.DiskPath == "" {
		klog.Error("no device path passed")
		msg = "no device path passed"
		klog.Error(msg)
		WrapResult(msg, UnprocessableEntityErrorCode, w)
		return
	} else if device.MountPoint == "" {
		klog.Error("no mount point passed")
		msg = "no mount point passed"
		klog.Error(msg)
		WrapResult(msg, UnprocessableEntityErrorCode, w)
		return
	}

	// erase the disk
	err := wipeDisk(device.DiskPath)
	if err != nil {
		msg = fmt.Sprintf("failed to erase disk %s, Error: %s", device.DiskPath, err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrGeneral, w)
		return
	}

	// format a disk partition with the ext4 file system
	err = formatDisk(device.DiskPath)
	if err != nil {
		msg = fmt.Sprintf("failed to format disk partition %s with ext4, Error: %s", device.DiskPath, err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrGeneral, w)
		return
	}

	// make directory
	err = makeDir(device.MountPoint)
	if err != nil {
		msg = fmt.Sprintf("failed to create directory %s, Error: %s", device.MountPoint, err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrGeneral, w)
		return
	}

	// mount the partition
	err = mountPartition(device.DiskPath, device.MountPoint)
	if err != nil {
		msg = fmt.Sprintf("failed to mount partition %s to %s, Error: %s", device.MountPoint, device.DiskPath, err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrGeneral, w)
		return
	}
	WrapResult(msg, ErrNone, w)
}

// RemoveHostPathDisk removes hostpath disk
func RemoveHostPathDisk(w http.ResponseWriter, r *http.Request) {
	var device LoopDevice
	var msg string

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&device); err != nil {
		msg = fmt.Sprintf("failed to read JSON encoded data, Error: %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrJsonDecode, w)
		return
	}
	klog.Info("remove a hostpath device, data: %v", device)
	if device.DiskPath == "" {
		klog.Error("no device path passed")
		msg = "no device path passed"
		klog.Error(msg)
		WrapResult(msg, UnprocessableEntityErrorCode, w)
		return
	} else if device.MountPoint == "" {
		klog.Error("no mount point passed")
		msg = "no mount point passed"
		klog.Error(msg)
		WrapResult(msg, UnprocessableEntityErrorCode, w)
		return
	}

	// unmount the partition
	err := unMountPartition(device.MountPoint)
	if err != nil {
		msg = fmt.Sprintf("failed to unmount partition %s, Error: %s", device.MountPoint, err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrGeneral, w)
		return
	}

	// erase the disk
	err = wipeDisk(device.DiskPath)
	if err != nil {
		msg = fmt.Sprintf("failed to erase disk %s, Error: %s", device.DiskPath, err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrGeneral, w)
		return
	}

	// remove directory
	err = removeDir(device.MountPoint)
	if err != nil {
		msg = fmt.Sprintf("failed to remove directory %s, Error: %s", device.MountPoint, err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrGeneral, w)
		return
	}
	WrapResult(msg, ErrNone, w)
}
