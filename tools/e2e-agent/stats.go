package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"k8s.io/klog/v2"
)

type Stats struct {
	ServiceAddr string `json:"serviceAddr"`
}

func GetStats(w http.ResponseWriter, r *http.Request) {
	var stats_req Stats
	var msg string

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&stats_req); err != nil {
		msg = fmt.Sprintf("failed to read JSON encoded data, Error: %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrJsonDecode, w)
		return
	}

	resp, err := http.Get("http://" + stats_req.ServiceAddr + "/stats")
	if err != nil {
		msg = fmt.Sprintf("failed to connect to stats module, Error: %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrConnectFail, w)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		msg = fmt.Sprintf("failed to connect to read response, Error: %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrReadFail, w)
		return
	}

	sb := string(body)
	klog.Info("stats response: ", sb)
	WrapResult(sb, ErrNone, w)
}
