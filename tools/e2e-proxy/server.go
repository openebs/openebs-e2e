package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

var portFwdMutex sync.Mutex
var forwardedPort map[string]int = map[string]int{}

func relay(src net.Conn, dst net.Conn) {
	defer func() { _ = src.Close() }()
	defer func() { _ = dst.Close() }()
	_, _ = io.Copy(dst, src)
}

func forward(conn net.Conn, remote string) {
	client, err := net.Dial("tcp", remote)
	if err != nil {
		log.Printf("port_forward: dial failed %s", remote)
		return
	}
	go relay(conn, client)
	go relay(client, conn)
}

func listenAndForward(listener net.Listener, remote string) {
	defer func() { _ = listener.Close() }()
	for {
		conn, err := listener.Accept()
		if err == nil {
			go forward(conn, remote)
		} else {
			log.Printf("accept failed %v\n", err)
		}
	}
}

func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err == nil {
		var tcpListener *net.TCPListener
		tcpListener, err = net.ListenTCP("tcp", addr)
		if err == nil {
			defer func() { _ = tcpListener.Close() }()
			return tcpListener.Addr().(*net.TCPAddr).Port, err
		} else {
			log.Printf("listen TCP failed %v\n", err)
		}
	} else {
		log.Printf("resolve TCP address failed %v\n", err)
	}
	return 0, err
}

func StartForwarding(remote string, portIn int) (int, error) {
	var port int
	var err error
	var ok bool
	portFwdMutex.Lock()
	defer portFwdMutex.Unlock()
	port, ok = forwardedPort[remote]
	if !ok {
		port = portIn
		if port == 0 {
			port, err = getFreePort()
			log.Printf("got port %d for %s\n", port, remote)
		}
		if err == nil {
			var listener net.Listener
			listener, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
			if err == nil {
				go listenAndForward(listener, remote)
				forwardedPort[remote] = port
				log.Printf("fowarding %s to port %d\n", remote, port)
			} else {
				log.Printf("listen failed\n")
			}
		}
	} else {
		log.Printf("already fowarding on port %d\n", port)
	}
	return port, err
}

var portNum = 9000

func forwardHandler(w http.ResponseWriter, r *http.Request) {
	fwdReqsList := make([]map[string]int, 0)
	body, err := io.ReadAll(r.Body)
	if err == nil {
		err = json.Unmarshal(body, &fwdReqsList)
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println("Received forwarding requests:")
	for _, fwdReq := range fwdReqsList {
		fmt.Printf("%v\n", fwdReq)
	}
	for _, fwdReq := range fwdReqsList {
		for hostAddress, port := range fwdReq {
			target := fmt.Sprintf("%s:%d", hostAddress, port)
			portSpec := 0
			{
				portFwdMutex.Lock()
				portNum += 1
				portSpec = portNum
				portFwdMutex.Unlock()
			}
			port, fwErr := StartForwarding(target, portSpec)
			if fwErr != nil {
				fmt.Printf("Failed     : {\"%s\":%d} %v\n", target, port, fwErr)
			}
		}
	}
}

func mapHandler(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.Marshal(forwardedPort)
	if err == nil {
		_, _ = w.Write(bytes)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	// Create Server and Route Handlers
	r := mux.NewRouter()
	r.HandleFunc("/listforwarding", mapHandler)
	r.HandleFunc("/forward", forwardHandler)
	r.HandleFunc("/health", healthHandler)
	r.HandleFunc("/readiness", readinessHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start Server
	go func() {
		log.Println("Starting Server")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// Graceful Shutdown
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_ = srv.Shutdown(ctx)

	log.Println("Shutting down")
	os.Exit(0)
}
