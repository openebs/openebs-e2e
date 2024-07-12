package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"k8s.io/klog/v2"
)

const (
	// DiskImageSize is the default file size(1GB) used while creating backing image
	DiskImageNamePrefix = "openebs-disk"
)

// LoopDevice has the attributes of a virtual disk which is emulated for testing
type LoopDevice struct {
	// Size in bytes
	Size int64 `json:"size"`
	// The backing image name
	// eg: fake123
	ImageName string `json:"imageName"`
	// Image directory
	// eg: /tmp
	ImgDir string `json:"imgDir"`
	// the disk name
	// eg: /tmp/loop9002
	DiskPath string `json:"diskPath"`
}

func (disk *LoopDevice) createDiskImage() error {
	f, err := os.CreateTemp(disk.ImgDir, DiskImageNamePrefix+"-*.img")
	if err != nil {
		return fmt.Errorf("error creating disk image. Error : %v", err)
	}
	disk.ImageName = f.Name()
	err = f.Truncate(disk.Size)
	if err != nil {
		return fmt.Errorf("error truncating disk image. Error : %v", err)
	}

	return nil
}

// CreateLoopDevice creates a loop device if the disk is not present
func CreateLoopDevice(w http.ResponseWriter, r *http.Request) {
	var device LoopDevice
	var msg string

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&device); err != nil {
		msg = fmt.Sprintf("failed to read JSON encoded data, Error: %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrJsonDecode, w)
		return
	}
	klog.Info("creates a loop device, data: %v", device)
	if device.ImgDir == "" {
		klog.Error("no device directory passed")
		msg = "no device directory passed"
		klog.Error(msg)
		WrapResult(msg, UnprocessableEntityErrorCode, w)
		return
	} else if device.Size <= 0 {
		klog.Error("no device size passed")
		msg = "no device size passed"
		klog.Error(msg)
		WrapResult(msg, UnprocessableEntityErrorCode, w)
		return
	}

	err := device.createDiskImage()
	if err != nil {
		msg = fmt.Sprintf("failed to create disk image, Error: %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrGeneral, w)
		return
	}

	// create the loop device using losetup
	createLoopDeviceCommand := fmt.Sprintf("losetup -f %s --show", device.ImageName)

	device.DiskPath, err = bashLocal(createLoopDeviceCommand)
	if err != nil {
		msg = fmt.Sprintf("error creating loop device. Error : %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrExecFailed, w)
		return
	}

	// Marshal the struct to JSON
	jsonData, err := json.Marshal(device)
	if err != nil {
		msg = fmt.Sprintf("failed to marshal the device struct to JSON Error : %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrGeneral, w)
		return
	}

	WrapResult(string(jsonData), ErrNone, w)
}

// DeleteLoopDevice detaches the loop device from the backing
// image. Also deletes the backing image and block device file in /dev
func DeleteLoopDevice(w http.ResponseWriter, r *http.Request) {
	var device LoopDevice
	var msg string

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&device); err != nil {
		msg = fmt.Sprintf("failed to read JSON encoded data, Error: %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrJsonDecode, w)
		return
	}

	klog.Info("creates a loop device, data: %v", device)
	if device.DiskPath == "" {
		klog.Error("no device path passed")
		msg = "no device path passed"
		klog.Error(msg)
		WrapResult(msg, UnprocessableEntityErrorCode, w)
		return
	} else if device.ImageName == "" {
		klog.Error("no device image passed")
		msg = "no device image passed"
		klog.Error(msg)
		WrapResult(msg, UnprocessableEntityErrorCode, w)
		return
	}

	detachLoopCommand := "losetup -d " + device.DiskPath
	_, err := bashLocal(detachLoopCommand)
	if err != nil {
		msg = fmt.Sprintf("cannot detach loop device. Error %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrExecFailed, w)
		return
	}
	err = os.Remove(device.ImageName)
	if err != nil {
		msg = fmt.Sprintf("could not delete backing disk image. Error : %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrGeneral, w)
		return
	}
	deleteLoopDeviceCommand := "rm " + device.DiskPath
	output, err := bashLocal(deleteLoopDeviceCommand)
	if err != nil {
		msg = fmt.Sprintf("could not delete loop device. Error : %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrExecFailed, w)
		return
	}
	WrapResult(output, ErrNone, w)
}
