package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"os/exec"
)

const (
	SYSRQ_TRIGGER_FILE = "/host/proc/sysrq-trigger"
)

func Setup() error {
	var (
		err error
		cmd *exec.Cmd
	)

	e2eHostAddr := os.Getenv("E2E_HOST_ADDR")
	ports := []string{
		os.Getenv("REST_PORT"),
		os.Getenv("MAYASTOR_PORT"),
		os.Getenv("MCP_REST_PORT"),
	}
	for _, port := range ports {
		if e2eHostAddr != "" {
			cmd = exec.Command(
				"iptables", "-t", "mangle", "-i", "eth0", "-s", e2eHostAddr,
				"-I", "PREROUTING", "-p", "tcp", "--dport", port, "-j", "ACCEPT",
				"-m", "comment", "--comment", "mayastor-e2e-test")
		} else {
			cmd = exec.Command("iptables", "-t", "mangle", "-i", "eth0",
				"-I", "PREROUTING", "-p", "tcp", "--dport", port, "-j", "ACCEPT",
				"-m", "comment", "--comment", "mayastor-e2e-test")

		}
		_, err = cmd.Output()
		if err != nil {
			return err
		}
	}
	return nil
}

// UngracefulReboot crashes and reboots the host machine
func UngracefulReboot() error {
	log.Printf("Rebooting node ungracefully")
	time.Sleep(2 * time.Second)
	data := []byte("c")
	err := os.WriteFile(SYSRQ_TRIGGER_FILE, data, 0644)
	return err
}

// GracefulReboot reboots the host gracefully
// It is not yet supported
func GracefulReboot() error {
	/*
		// Send SIGTERM to all processes
		data := []byte("e")
		if err := ioutil.WriteFile(SYSRQ_TRIGGER_FILE, data, 0644); err != nil {
			return err
		}

		time.Sleep(5 * time.Second)

		// Send SIGKILL to remaining processes
		data = []byte("i")
		if err := ioutil.WriteFile(SYSRQ_TRIGGER_FILE, data, 0644); err != nil {
			return err
		}

		time.Sleep(5 * time.Second)

		// Reboot the node
		data = []byte("b")
		if err := ioutil.WriteFile(SYSRQ_TRIGGER_FILE, data, 0644); err != nil {
			return err
		}
	*/
	return nil
}

// DropConnectionsFromNodes creates rules to drop connections from other k8s nodes
func DropConnectionsFromNodes(nodes []string, networkInterface string) error {
	log.Printf("Drop Connections from %v", nodes)
	for _, node := range nodes {
		// because we specify eth0 interface we can DROP even our IP w/o problems, because self communication uses lo
		cmd := exec.Command("iptables", "-t", "mangle", "-I", "PREROUTING", "-i", networkInterface, "-s", node, "-j", "DROP", "-m", "comment", "--comment", "mayastor-e2e-test")
		_, err := cmd.Output()
		if err != nil {
			return err
		}
	}
	return nil
}

// AcceptConnectionsFromNodes removes the rules set by
// DropConnectionsFromNodes so that other k8s nodes can reach this node again
func AcceptConnectionsFromNodes(nodes []string, networkInterface string) error {
	log.Printf("Accept Connections from %v", nodes)
	for _, node := range nodes {
		cmd := exec.Command("iptables", "-t", "mangle", "-D", "PREROUTING", "-i", networkInterface, "-s", node, "-j", "DROP", "-m", "comment", "--comment", "mayastor-e2e-test")
		_, err := cmd.Output()
		if err != nil {
			return err
		}
	}
	return nil
}

// DropIncomingTrafficOnNode creates rules to drop every incoming traffic
// on the specified node, except for SSH (port 22)
func DropIncomingTrafficOnNode() error {
	log.Printf("Dropping incoming traffic on node except for SSH")
	// Flush existing rules
	if err := AcceptIncomingTrafficOnNode(); err != nil {
		return fmt.Errorf("failed to reset iptables before applying drop rules: %w", err)
	}

	// Allow SSH explicitly (for INPUT and OUTPUT)
	for _, chain := range []string{"OUTPUT", "INPUT"} {
		if err := runIptables(
			"-A", chain, "-p", "tcp", "--dport", "22",
			"-j", "ACCEPT", "-m", "comment", "--comment", "mayastor-e2e-test",
		); err != nil {
			return fmt.Errorf("failed to allow SSH on %s chain: %w", chain, err)
		}
		if err := runIptables(
			"-A", chain, "-p", "tcp", "--dport", "10012",
			"-j", "ACCEPT", "-m", "comment", "--comment", "mayastor-e2e-test",
		); err != nil {
			return fmt.Errorf("failed to allow e2e-agent rest port on %s chain: %w", chain, err)
		}
	}

	// Drop everything else
	if err := runIptables(
		"-A", "INPUT", "-j", "DROP", "-m", "comment", "--comment", "mayastor-e2e-test",
	); err != nil {
		return fmt.Errorf("failed to append DROP rule: %w", err)
	}

	return nil
}

// AcceptIncomingTrafficOnNode removes the rules set by
// DropIncomingTrafficOnNode so that the node can accept incoming traffic again
// This function flushes the OUTPUT and INPUT chains, effectively removing all rules
// including the DROP rule set by DropIncomingTrafficOnNode.
func AcceptIncomingTrafficOnNode() error {
	log.Println("Accepting incoming traffic from all nodes")
	// Chains we want to flush
	for _, chain := range []string{"OUTPUT", "INPUT"} {
		if err := runIptables("-F", chain); err != nil {
			return fmt.Errorf("failed to flush %s chain: %w", chain, err)
		}
	}

	return nil
}

func runIptables(args ...string) error {
	// Prepare the command
	if len(args) == 0 {
		return fmt.Errorf("no iptables arguments provided")
	}
	log.Printf("Running iptables with args: %v", args)
	// Use exec.Command to run iptables with the provided arguments
	cmd := exec.Command("iptables", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("iptables %v failed: %v, output: %s", args, err, string(out))
	}
	return nil
}
