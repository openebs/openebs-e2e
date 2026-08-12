package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"k8s.io/klog/v2"
)

type E2eAgentErrcode int

const (
	// general errors
	ErrNone         E2eAgentErrcode = 0
	ErrGeneral      E2eAgentErrcode = 1
	ErrJsonDecode   E2eAgentErrcode = 2
	ErrJsonEncode   E2eAgentErrcode = 3
	ErrReadFail     E2eAgentErrcode = 4
	ErrExecFailed   E2eAgentErrcode = 5
	ErrFileNotExist E2eAgentErrcode = 6

	// stats errors
	ErrConnectFail E2eAgentErrcode = 101
)

type E2eAgentError struct {
	Output    string          `json:"output"`
	Errorcode E2eAgentErrcode `json:"errorcode"`
}

func WrapResult(output string, errorcode E2eAgentErrcode, w http.ResponseWriter) {
	e2eagenterr := E2eAgentError{
		Output:    output,
		Errorcode: errorcode,
	}
	jsn, err := json.Marshal(e2eagenterr)
	if err != nil {
		klog.Error("WrapResult: failed to marshal error, Error: ", err)
		//nolint:errcheck // secondary write failure is not actionable here
		fmt.Fprint(w, "WrapResult: failed to marshal error")
		w.WriteHeader(InternalServerErrorCode)
		return
	}
	if errorcode != ErrNone {
		klog.Error("output: ", output, " errorcode: ", errorcode)
	} else {
		klog.Info("output: ", output, " errorcode: ", errorcode)
	}
	//nolint:errcheck // secondary write failure is not actionable here
	fmt.Fprint(w, string(jsn))
}
