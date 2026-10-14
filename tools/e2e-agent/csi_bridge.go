package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"

	"k8s.io/klog/v2"
)

// csiBridges tracks the running TCP<->unix-socket proxies started by
// StartCsiSocketBridge, keyed by the TCP port they listen on, so
// StopCsiSocketBridge can find and stop the right one.
var (
	csiBridgesMu sync.Mutex
	csiBridges   = map[int]context.CancelFunc{}
)

// StartCsiSocketBridge listens on a TCP port and proxies every connection
// to a unix socket -- specifically, a node plugin's CSI socket, which
// kubelet requires to live on the host filesystem (reachable here via this
// agent's own hostPath mount) rather than inside the plugin pod alone. The
// agent itself runs with hostNetwork, so the listening port is reachable
// from outside the cluster exactly the way this agent's own HTTP port
// already is -- no socat, no kubectl port-forward.
func StartCsiSocketBridge(socketPath string, port int) error {
	csiBridgesMu.Lock()
	defer csiBridgesMu.Unlock()
	if _, exists := csiBridges[port]; exists {
		return fmt.Errorf("a CSI socket bridge is already listening on port %d", port)
	}

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("listening on port %d: %w", port, err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	csiBridges[port] = cancel

	go func() {
		<-ctx.Done()
		_ = ln.Close()
	}()
	go acceptCsiBridgeConns(ln, socketPath)

	klog.Info("Started CSI socket bridge on port ", port, " to ", socketPath)
	return nil
}

func acceptCsiBridgeConns(ln net.Listener, socketPath string) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			// Listener closed by StopCsiSocketBridge, or a fatal accept
			// error -- either way, nothing more to serve.
			return
		}
		go proxyCsiBridgeConn(conn, socketPath)
	}
}

func proxyCsiBridgeConn(conn net.Conn, socketPath string) {
	defer func() { _ = conn.Close() }()
	upstream, err := net.Dial("unix", socketPath)
	if err != nil {
		klog.Error("CSI socket bridge: dialing ", socketPath, " Error: ", err)
		return
	}
	defer func() { _ = upstream.Close() }()

	done := make(chan struct{}, 2)
	go func() { _, _ = io.Copy(upstream, conn); done <- struct{}{} }()
	go func() { _, _ = io.Copy(conn, upstream); done <- struct{}{} }()
	<-done
}

// StopCsiSocketBridge stops the bridge listening on port, if any. Stopping
// a port with no bridge is not an error -- the caller may be cleaning up
// after a test that never got as far as starting one.
func StopCsiSocketBridge(port int) {
	csiBridgesMu.Lock()
	defer csiBridgesMu.Unlock()
	cancel, ok := csiBridges[port]
	if !ok {
		return
	}
	cancel()
	delete(csiBridges, port)
	klog.Info("Stopped CSI socket bridge on port ", port)
}
