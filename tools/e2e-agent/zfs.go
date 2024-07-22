package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"k8s.io/klog/v2"
)

type Zpool struct {
	PoolDiskPath string `json:"poolDiskPath"`
	PoolName     string `json:"poolName"`
}

// ZfsVersion check zfs version installed on node
func ZfsVersion(w http.ResponseWriter, r *http.Request) {
	var msg string
	klog.Info("Get zfs version installed")

	zfsVersionCommand := "zfs version"
	output, err := bashLocal(zfsVersionCommand)
	if err != nil {
		msg = fmt.Sprintf("cannot get zfs version. Error %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrExecFailed, w)
		return
	}
	WrapResult(string(output), ErrNone, w)
}

// ZfsListPool list zfs pool
func ZfsListPool(w http.ResponseWriter, r *http.Request) {
	var msg string
	klog.Info("List zfs pool")

	zfsPoolListCommand := "zpool list"
	output, err := bashLocal(zfsPoolListCommand)
	if err != nil {
		msg = fmt.Sprintf("cannot list zfs pool Error %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrExecFailed, w)
		return
	}
	WrapResult(string(output), ErrNone, w)
}

// ZfsCreatePool create zfs pool
func ZfsCreatePool(w http.ResponseWriter, r *http.Request) {
	var msg string
	var zPool Zpool
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&zPool); err != nil {
		msg = fmt.Sprintf("failed to read JSON encoded data, Error: %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrJsonDecode, w)
		return
	}
	if zPool.PoolDiskPath == "" {
		msg = "no device path passed"
		klog.Error(msg)
		WrapResult(msg, UnprocessableEntityErrorCode, w)
		return
	}
	if zPool.PoolName == "" {
		msg = "no zfs pool name passed"
		klog.Error(msg)
		WrapResult(msg, UnprocessableEntityErrorCode, w)
		return
	}
	klog.Info("creates zfs pool, data: %v", zPool)

	zfsPoolCreateCommand := fmt.Sprintf("zpool create %s %s", zPool.PoolName, zPool.PoolDiskPath)
	output, err := bashLocal(zfsPoolCreateCommand)
	if err != nil {
		msg = fmt.Sprintf("cannot create zfs pool Error %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrExecFailed, w)
		return
	}
	WrapResult(string(output), ErrNone, w)
}

// ZfsDestroyPool destroy zfs pool
func ZfsDestroyPool(w http.ResponseWriter, r *http.Request) {
	var msg string
	var zPool Zpool
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&zPool); err != nil {
		msg = fmt.Sprintf("failed to read JSON encoded data, Error: %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrJsonDecode, w)
		return
	}
	if zPool.PoolName == "" {
		msg = "no zfs pool name passed"
		klog.Error(msg)
		WrapResult(msg, UnprocessableEntityErrorCode, w)
		return
	}

	klog.Info("destroy zfs pool, data: %v", zPool)

	ZfsPoolRemoveCommand := fmt.Sprintf("zpool destroy %s", zPool.PoolName)
	output, err := bashLocal(ZfsPoolRemoveCommand)
	if err != nil {
		msg = fmt.Sprintf("cannot destroy zfs pool Error %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrExecFailed, w)
		return
	}
	WrapResult(string(output), ErrNone, w)
}
