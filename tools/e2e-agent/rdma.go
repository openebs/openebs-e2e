package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"k8s.io/klog/v2"
)

type Rdma struct {
	DeviceName    string `json:"deviceName"`
	InterfaceName string `json:"interfaceName"`
}

// ListRdmaDevice list RDMA device
func ListRdmaDevice(w http.ResponseWriter, r *http.Request) {
	var msg string
	klog.Info("List available RDMA device")

	rdmaDeviceListCommand := "rdma link -j"
	output, err := bashLocal(rdmaDeviceListCommand)
	if err != nil {
		msg = fmt.Sprintf("cannot list RDMA device. Error %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrExecFailed, w)
		return
	}
	WrapResult(string(output), ErrNone, w)
}

// CreateRdmaDevice create rdma device
func CreateRdmaDevice(w http.ResponseWriter, r *http.Request) {
	var msg string
	var rdma Rdma
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&rdma); err != nil {
		msg = fmt.Sprintf("failed to read JSON encoded data, Error: %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrJsonDecode, w)
		return
	}
	if rdma.DeviceName == "" {
		msg = "no device name passed"
		klog.Error(msg)
		WrapResult(msg, UnprocessableEntityErrorCode, w)
		return
	}
	if rdma.InterfaceName == "" {
		msg = "no interface name passed"
		klog.Error(msg)
		WrapResult(msg, UnprocessableEntityErrorCode, w)
		return
	}
	klog.Info("creates rdma device, data: %v", rdma)

	rdmaDeviceCreateCommand := fmt.Sprintf("rdma link add %s type rxe netdev %s", rdma.DeviceName, rdma.InterfaceName)
	output, err := bashLocal(rdmaDeviceCreateCommand)
	if err != nil {
		msg = fmt.Sprintf("cannot create rdma device, Error %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrExecFailed, w)
		return
	}
	WrapResult(string(output), ErrNone, w)
}

// DeleteRdmaDevice destroy rdma device
func DeleteRdmaDevice(w http.ResponseWriter, r *http.Request) {
	var msg string
	var rdma Rdma
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&rdma); err != nil {
		msg = fmt.Sprintf("failed to read JSON encoded data, Error: %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrJsonDecode, w)
		return
	}
	if rdma.DeviceName == "" {
		msg = "no rdma device name passed"
		klog.Error(msg)
		WrapResult(msg, UnprocessableEntityErrorCode, w)
		return
	}

	klog.Info("delete rdma device, data: %v", rdma)
	rdmaDeviceDeleteCommand := fmt.Sprintf("rdma link delete %s", rdma.DeviceName)
	output, err := bashLocal(rdmaDeviceDeleteCommand)
	if err != nil {
		msg = fmt.Sprintf("cannot delete rdma device, Error %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrExecFailed, w)
		return
	}
	WrapResult(string(output), ErrNone, w)
}
