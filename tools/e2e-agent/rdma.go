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

// ListDevLink list dev links
func ListDevLink(w http.ResponseWriter, r *http.Request) {
	var msg string
	klog.Info("List available dev link")

	devLinkPortCommand := "devlink port show -j"
	output, err := bashLocal(devLinkPortCommand)
	if err != nil {
		msg = fmt.Sprintf("cannot list dev links. Error %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrExecFailed, w)
		return
	}
	WrapResult(string(output), ErrNone, w)
}

// EnableDevLink enable dev link
func EnableDevLink(w http.ResponseWriter, r *http.Request) {
	setDevLinkState(w, r, true)
}

// DisableDevLink disable dev link
func DisableDevLink(w http.ResponseWriter, r *http.Request) {
	setDevLinkState(w, r, false)
}

func setDevLinkState(w http.ResponseWriter, r *http.Request, enable bool) {
	var msg string
	var devLink DevLinkPort
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&devLink); err != nil {
		msg = fmt.Sprintf("failed to read JSON encoded data, Error: %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrJsonDecode, w)
		return
	}
	if devLink.DevLinkPort == "" {
		msg = "no dev link name passed"
		klog.Error(msg)
		WrapResult(msg, UnprocessableEntityErrorCode, w)
		return
	}
	klog.Info("set dev link, data: %v, enable_rdma vale: %s", devLink, enable)

	devLinkDownCommand := fmt.Sprintf("devlink dev param set %s name enable_rdma value %v cmode driverinit", devLink.DevLinkPort, enable)
	output, err := bashLocal(devLinkDownCommand)
	if err != nil {
		msg = fmt.Sprintf("cannot set %s dev link enable_rdma value to %v, Error %s", devLink.DevLinkPort, enable, err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrExecFailed, w)
		return
	}
	// reload dev link
	err = devLinkDriverReload(devLink.DevLinkPort)
	if err != nil {
		msg = fmt.Sprintf("cannot reload dev link %s, Error %s", devLink.DevLinkPort, err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrExecFailed, w)
		return
	}
	WrapResult(string(output), ErrNone, w)
}

func devLinkDriverReload(devLink string) error {
	klog.Info("reload dev link, link: %v", devLink)
	devLinkDriverReload := fmt.Sprintf("devlink dev reload %s", devLink)
	output, err := bashLocal(devLinkDriverReload)
	if err != nil {
		return fmt.Errorf("failed to reload dev link %s, output: %s, error: %v", devLink, output, err)
	}
	return nil
}
