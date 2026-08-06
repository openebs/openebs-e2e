package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"k8s.io/klog/v2"
)

type NodeList struct {
	Nodes            []string `json:"nodes"`
	NetworkInterface string   `json:"networkInterface"`
}

type DiskPool struct {
	Disk string `json:"disk"`
}

type Device struct {
	Device     string `json:"device"`
	Table      string `json:"table"`
	DevicePath string `json:"devicePath"`
	Uuid       string `json:"uuid"`
	FsType     string `json:"fsType"`
}

type ControlledDevice struct {
	Device string `json:"device"`
	State  string `json:"state"`
}

type Product struct {
	Product string `json:"product"`
	Pid     string `json:"pid"`
}

type Nvme struct {
	TargetIp string `json:"targetIp"`
	Nqn      string `json:"nqn"`
	HostNqn  string `json:"hostNqn"`
	HostId   string `json:"hostId"`
}

type Disk struct {
	Device         string `json:"device"`
	SeekParam      string `json:"seekParam"`
	BlockSizeParam string `json:"blockSizeParam"`
}

type blkDiscard struct {
	Device  string `json:"device"`
	Options string `json:"options"`
}

type CmpPaths struct {
	Path1 string `json:"path1"`
	Path2 string `json:"path2"`
}

type NetworkInterface struct {
	NetworkInterface string `json:"networkInterface"`
}

type DevLinkPort struct {
	DevLinkPort string `json:"devLinkPort"`
}

func homePage(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, "Welcome home!\n")
}

type CmdList struct {
	Cmd string `json:"cmd"`
}

type KernelModule struct {
	Name           string `json:"name"`
	PersistentPath string `json:"persistentPath"`
}

type DmDevice struct {
	Device  string `json:"device"`
	Sectors uint64 `json:"sectors"`
}

var Version = "undefined"

const (
	InternalServerErrorCode      = 500
	UnprocessableEntityErrorCode = 422
	rpcGssdServiceName           = "rpc-gssd"
	dmsetupTimeout               = 30 * time.Second
)

type Command string

const (
	Start   Command = "start"
	Stop    Command = "stop"
	Restart Command = "restart"
)

func main() {
	// Following works when both glog & klog have been used
	// ignore errors - not much we can do here
	_ = flag.Set("logtostderr", "true")
	_ = flag.Set("alsologtostderr", "true")
	flag.Parse()

	klogFlags := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(klogFlags)

	// Sync the glog and klog flags.
	flag.CommandLine.VisitAll(func(f1 *flag.Flag) {
		f2 := klogFlags.Lookup(f1.Name)
		if f2 != nil {
			value := f1.Value.String()
			// ignore errors - not much we can do here
			_ = f2.Value.Set(value)
		}
	})
	defer klog.Flush()
	klog.Info("Starting e2e agent, version: ", Version)
	if err := Setup(); err != nil {
		klog.Fatal(err)
		log.Fatal(err)
	}
	handleRequests()
}

func handleRequests() {
	podIP := os.Getenv("MY_POD_IP")
	restPort := os.Getenv("REST_PORT")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage)
	router.HandleFunc("/ungracefulReboot", ungracefulReboot).Methods("POST")
	router.HandleFunc("/gracefulReboot", gracefulReboot).Methods("POST")
	router.HandleFunc("/dropConnectionsFromNodes", dropConnectionsFromNodes).Methods("POST")
	router.HandleFunc("/acceptConnectionsFromNodes", acceptConnectionsFromNodes).Methods("POST")
	router.HandleFunc("/createFaultyDevice", createFaultyDevice).Methods("POST")
	router.HandleFunc("/deleteFaultyDevice", deleteFaultyDevice).Methods("POST")
	router.HandleFunc("/reloadDevice", reloadDevice).Methods("POST")
	router.HandleFunc("/devicecontrol", controlDevice).Methods("POST")
	router.HandleFunc("/killioengine", killIoEngine).Methods("POST")
	router.HandleFunc("/killCsiController", killCsiController).Methods("POST")
	router.HandleFunc("/killCsiNode", killCsiNode).Methods("POST")
	router.HandleFunc("/getdevicestate", getDeviceState).Methods("POST")
	router.HandleFunc("/gethostid", getHostID).Methods("POST")
	router.HandleFunc("/nvmeconnect", NvmeConnect).Methods("POST")
	router.HandleFunc("/nvmedisconnect", NvmeDisconnect).Methods("POST")
	router.HandleFunc("/nvmelist", NvmeList).Methods("POST")
	router.HandleFunc("/nvmelistsubsys", NvmeListSubSys).Methods("POST")
	router.HandleFunc("/checksumdevice", ChecksumDevice).Methods("POST")
	router.HandleFunc("/fscheckdevice", FsCheckDevice).Methods("POST")
	router.HandleFunc("/xfscheckdevice", XFSCheckDevice).Methods("POST")
	router.HandleFunc("/fsfreezedevice", FsFreezeDevice).Methods("POST")
	router.HandleFunc("/fsunfreezedevice", FsUnfreezeDevice).Methods("POST")
	router.HandleFunc("/listdevice", ListDevice).Methods("POST")
	router.HandleFunc("/flushDiskWriteCache", flushDiskWriteCache).Methods("POST")
	router.HandleFunc("/zeroingdisk", ZeroingDisk).Methods("POST")
	router.HandleFunc("/parted", Parted).Methods("POST")
	router.HandleFunc("/findmnt", Findmnt).Methods("POST")
	router.HandleFunc("/lsblk", Lsblk).Methods("POST")
	router.HandleFunc("/dmesg", Dmesg).Methods("POST")
	router.HandleFunc("/syslog", Getsyslog).Methods("POST")
	router.HandleFunc("/listreservation", ListReservation).Methods("POST")
	router.HandleFunc("/nvmeconnectwithhostid", NvmeConnectWithHostId).Methods("POST")
	router.HandleFunc("/blkdiscard", BlkDiscard).Methods("POST")
	router.HandleFunc("/event/list", EventList).Methods("POST")
	router.HandleFunc("/event/publish", EventPublish).Methods("POST")
	router.HandleFunc("/event/subscribe", EventSubscribe).Methods("POST")
	router.HandleFunc("/event/unsubscribe", EventUnsubscribe).Methods("POST")
	router.HandleFunc("/event/unsubscribeall", EventUnsubscribeAll).Methods("POST")
	router.HandleFunc("/stats", GetStats).Methods("POST")
	router.HandleFunc("/cmp", Cmp).Methods("POST")
	router.HandleFunc("/hugepagezero", ZeroingHugePages).Methods("POST")
	router.HandleFunc("/dropIncomingTrafficOnNode", dropIncomingTrafficOnNode).Methods("POST")
	router.HandleFunc("/acceptIncomingTrafficOnNode", acceptIncomingTrafficOnNode).Methods("POST")
	router.HandleFunc("/loadKernelModule", LoadKernelModule).Methods("POST")
	router.HandleFunc("/unloadKernelModule", UnloadKernelModule).Methods("POST")
	router.HandleFunc("/getProcessID", GetProcessID).Methods("POST")
	router.HandleFunc("/isKernelModuleLoaded", IsKernelModuleLoaded).Methods("POST")
	router.HandleFunc("/isPackageInstalled", IsPackageInstalled).Methods("POST")
	router.HandleFunc("/isKernelModulePersistent", IsKernelModulePersistent).Methods("POST")
	router.HandleFunc("/configureNonPersistentHugePages", ConfigureNonPersistentHugePages).Methods("POST")
	router.HandleFunc("/isHugePagesPersistent", IsHugePagesPersistent).Methods("POST")
	router.HandleFunc("/isHugePagesConfigured", IsHugePagesConfigured).Methods("POST")
	router.HandleFunc("/restartService", RestartService).Methods("POST")
	router.HandleFunc("/startRpcGssdService", StartRpcGssdService).Methods("POST")
	router.HandleFunc("/stopRpcGssdService", StopRpcGssdService).Methods("POST")
	router.HandleFunc("/dm/createPassThrough", createPassThroughDevice).Methods("POST")
	router.HandleFunc("/dm/suspend", suspendDevice).Methods("POST")
	router.HandleFunc("/dm/resume", resumeDevice).Methods("POST")
	router.HandleFunc("/dm/remove", removeDevice).Methods("POST")
	router.HandleFunc("/dm/getDeviceSectors", getDeviceSizeInSectors).Methods("POST")

	//LVM
	router.HandleFunc("/lvmversion", LvmVersion).Methods("POST")
	router.HandleFunc("/lvmlistvg", LvmListVg).Methods("POST")
	router.HandleFunc("/lvmlistpv", LvmListPv).Methods("POST")
	router.HandleFunc("/lvmcreatepv", LvmCreatePv).Methods("POST")
	router.HandleFunc("/lvmcreatevg", LvmCreateVg).Methods("POST")
	router.HandleFunc("/lvmremovepv", LvmRemovePv).Methods("POST")
	router.HandleFunc("/lvmremovevg", LvmRemoveVg).Methods("POST")
	router.HandleFunc("/lvmthinpoolautoextendthreshold", LvmThinPoolAutoExtendThreshold).Methods("POST")
	router.HandleFunc("/lvmthinpoolautoextendpercent", LvmThinPoolAutoExtendPercent).Methods("POST")
	router.HandleFunc("/lvmlvchangemonitor", LvmLvChangeMonitor).Methods("POST")
	router.HandleFunc("/lvmlvremovethinpool", LvmLvRemoveThinPool).Methods("POST")
	//loop device
	router.HandleFunc("/createloopdevice", CreateLoopDevice).Methods("POST")
	router.HandleFunc("/deleteloopdevice", DeleteLoopDevice).Methods("POST")
	//ZFS
	router.HandleFunc("/zfsversion", ZfsVersion).Methods("POST")
	router.HandleFunc("/zfslistpool", ZfsListPool).Methods("POST")
	router.HandleFunc("/zfscreatepool", ZfsCreatePool).Methods("POST")
	router.HandleFunc("/zfsdestroypool", ZfsDestroyPool).Methods("POST")
	//localPV
	router.HandleFunc("/createhostpathdisk", CreateHostPathDisk).Methods("POST")
	router.HandleFunc("/removehostpathdisk", RemoveHostPathDisk).Methods("POST")
	//RDMA
	router.HandleFunc("/listrdmadevice", ListRdmaDevice).Methods("POST")
	router.HandleFunc("/createrdmadevice", CreateRdmaDevice).Methods("POST")
	router.HandleFunc("/deleterdmadevice", DeleteRdmaDevice).Methods("POST")
	router.HandleFunc("/enablenetworkinterface", EnableNetworkInterface).Methods("POST")
	router.HandleFunc("/disablenetworkinterface", DisableNetworkInterface).Methods("POST")
	router.HandleFunc("/enabledevlink", EnableDevLink).Methods("POST")
	router.HandleFunc("/disabledevlink", DisableDevLink).Methods("POST")
	router.HandleFunc("/listdevlink", ListDevLink).Methods("POST")
	log.Fatal(http.ListenAndServe(podIP+":"+restPort, router))
}

func ungracefulReboot(w http.ResponseWriter, r *http.Request) {
	go func() {
		if err := UngracefulReboot(); err != nil {
			klog.Info(err)
		}
	}()
}

func gracefulReboot(w http.ResponseWriter, r *http.Request) {
	klog.Info("Graceful reboots are not yet supported")
	_, _ = fmt.Fprint(w, "Graceful reboots are not yet supported")
}

func dropConnectionsFromNodes(w http.ResponseWriter, r *http.Request) {
	var list NodeList
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&list); err != nil {
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
	}
	klog.Info("Dropping connection from nodes ", list.Nodes, list.NetworkInterface)
	if err := DropConnectionsFromNodes(list.Nodes, list.NetworkInterface); err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to drop connection from nodes:", list.Nodes, "Error: ", err)
		return
	}
	klog.Info("Successfully stopped network services")
}

func acceptConnectionsFromNodes(w http.ResponseWriter, r *http.Request) {
	var list NodeList
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&list); err != nil {
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
	}
	klog.Info("Accept connection from nodes ", list.Nodes, list.NetworkInterface)
	err := AcceptConnectionsFromNodes(list.Nodes, list.NetworkInterface)
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to accept connection from nodes:", list.Nodes, "Error: ", err)
		return
	}
	_, _ = fmt.Fprint(w, "Successfully started network services\n")
	klog.Info("Successfully started network services")
}

func dropIncomingTrafficOnNode(w http.ResponseWriter, r *http.Request) {

	klog.Info("Dropping incoming traffic on node")
	if err := DropIncomingTrafficOnNode(); err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to drop incoming traffic from node", "Error: ", err)
		return
	}
	klog.Info("Successfully stopped network services")
}

func acceptIncomingTrafficOnNode(w http.ResponseWriter, r *http.Request) {
	klog.Info("Accept incoming traffic on node ")
	err := AcceptIncomingTrafficOnNode()
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to accept incoming traffic on node", "Error: ", err)
		return
	}
	klog.Info("Successfully started network services")
}

func LoadKernelModule(w http.ResponseWriter, r *http.Request) {
	var module KernelModule
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&module); err != nil {
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
		return
	}
	klog.Info("Loading kernel module ", module.Name)
	params := fmt.Sprintf("chroot /host ; /sbin/modprobe -d /host %s ; echo $?", module.Name)
	output, err := bashLocal(params)
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to load kernel module:", module, "Error: ", err)
		return
	}
	klog.Info("Successfully loaded kernel module", module.Name)
	WrapResult(output, ErrNone, w)
}

func UnloadKernelModule(w http.ResponseWriter, r *http.Request) {
	var module KernelModule
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&module); err != nil {
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
		return
	}
	klog.Info("Unloading kernel module ", module.Name)
	params := fmt.Sprintf("chroot /host ; /sbin/modprobe -r -d /host %s ; echo $?", module.Name)
	klog.Info("Unloading kernel module")
	output, err := bashLocal(params)
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to unload kernel module:", module, "Error: ", err)
		return
	}
	klog.Info("Successfully unloaded kernel module", module.Name)
	WrapResult(output, ErrNone, w)
}

func IsKernelModuleLoaded(w http.ResponseWriter, r *http.Request) {
	var module KernelModule
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&module); err != nil {
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
		return
	}
	klog.Info("Checking if module is loaded ", module)
	params := fmt.Sprintf("lsmod | grep -w '^%s' | wc -l", module.Name)
	klog.Info("Checking if module is loaded")
	output, err := bashLocal(params)
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to check if module is loaded:", module, "Error: ", err)
		return
	}
	klog.Info("Successfully checked if module is loaded")
	WrapResult(output, ErrNone, w)
}

func IsPackageInstalled(w http.ResponseWriter, r *http.Request) {
	var packageName string
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&packageName); err != nil {
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
		return
	}
	klog.Info("Checking if package is installed ", packageName)
	params := fmt.Sprintf("nsenter --mount=/proc/1/ns/mnt dpkg -s %s >/dev/null 2>&1 && echo 0 || echo 1", packageName)
	klog.Info("Checking if package is installed")
	output, err := bashLocal(params)
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to check if package is installed:", packageName, "Error: ", err)
		return
	}
	klog.Info("Successfully checked if package is installed")
	WrapResult(output, ErrNone, w)
}

func IsKernelModulePersistent(w http.ResponseWriter, r *http.Request) {
	var module KernelModule
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&module); err != nil {
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
		return
	}
	klog.Info("Checking if module is persistent ", module)
	params := fmt.Sprintf(
		"if [ -f /host%s ]; then "+
			"grep -w '^%s$' /host%s | wc -l; "+
			"else "+
			"echo 0; "+
			"fi",
		module.PersistentPath,
		module.Name,
		module.PersistentPath,
	)
	klog.Info("Checking if module is persistent")
	output, err := bashLocal(params)
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to check if module is persistent:", module, "Error: ", err)
		return
	}
	klog.Info("Successfully checked if module is persistent")
	WrapResult(output, ErrNone, w)
}

func ConfigureNonPersistentHugePages(w http.ResponseWriter, r *http.Request) {
	params := "echo 1024 > /sys/kernel/mm/hugepages/hugepages-2048kB/nr_hugepages; echo $?"
	klog.Info("Configuring non-persistent hugepages")
	output, err := bashLocal(params)
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to configure non-persistent hugepages", "Error: ", err)
		return
	}
	klog.Info("Successfully configured non-persistent hugepages")
	WrapResult(output, ErrNone, w)
}

func IsHugePagesConfigured(w http.ResponseWriter, r *http.Request) {
	params := "cat /sys/kernel/mm/hugepages/hugepages-2048kB/nr_hugepages"
	klog.Info("Checking if hugepages are configured")
	output, err := bashLocal(params)
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to check if hugepages are configured", "Error: ", err)
		return
	}
	klog.Info("Successfully checked if hugepages are configured")
	WrapResult(output, ErrNone, w)
}

func IsHugePagesPersistent(w http.ResponseWriter, r *http.Request) {
	params := `grep -h -E '^vm\.nr_hugepages' /host/etc/sysctl.conf /host/etc/sysctl.d/*.conf 2>/dev/null | awk -F'=' '{gsub(/ /,"",$2); print $2}' | tail -n1 || echo 0`
	klog.Info("Checking if hugepages are persistent")
	output, err := bashLocal(params)
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to check if hugepages are persistent", "Error: ", err)
		return
	}
	klog.Info("Successfully checked if hugepages are persistent")
	WrapResult(output, ErrNone, w)
}

func RestartService(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	var service string
	if err := d.Decode(&service); err != nil {
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
		return
	}
	klog.Info("Restarting service ", service)
	output, err := runSystemctlCommand(service, Restart)
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to restart service:", service, "Error: ", err)
		return
	}
	klog.Info("Successfully restarted service", service)
	WrapResult(output, ErrNone, w)
}

func StartRpcGssdService(w http.ResponseWriter, r *http.Request) {
	klog.Info("Starting service ", rpcGssdServiceName)
	output, err := runSystemctlCommand(rpcGssdServiceName, Start)
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to start service:", rpcGssdServiceName, "Error: ", err)
		return
	}
	klog.Info("Successfully started service", rpcGssdServiceName)
	WrapResult(output, ErrNone, w)
}

func StopRpcGssdService(w http.ResponseWriter, r *http.Request) {
	klog.Info("Stopping service ", rpcGssdServiceName)
	output, err := runSystemctlCommand(rpcGssdServiceName, Stop)
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to stop service:", rpcGssdServiceName, "Error: ", err)
		return
	}
	klog.Info("Successfully stopped service", rpcGssdServiceName)
	WrapResult(output, ErrNone, w)
}

func runSystemctlCommand(service string, command Command) (string, error) {
	params := fmt.Sprintf("nsenter --target 1 --mount --uts --ipc --net --pid -- systemctl %s %s ; echo $?", command, service)
	klog.Info("Executing command ", params, " on service ", service)
	output, err := bashLocal(params)
	if err != nil {
		klog.Error("failed to execute command:", params, " on service:", service, "Error: ", err)
		return output, err
	}
	klog.Info("Successfully executed command ", params, " on service ", service)
	return output, nil
}

func GetProcessID(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	var process string
	if err := d.Decode(&process); err != nil {
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
		return
	}
	params := fmt.Sprintf("pidof %s", process)
	klog.Info("Retrieving process ID for ", process)
	output, err := bashLocal(params)
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to retrieve process ID for:", process, "Error: ", err)
		return
	}
	klog.Info("Successfully retrieved process ID")
	WrapResult(output, ErrNone, w)
}

func createFaultyDevice(w http.ResponseWriter, r *http.Request) {

	var device Device

	// Decode request
	if err := json.NewDecoder(r.Body).Decode(&device); err != nil {
		klog.Error("decode failed:", err)
		WrapResult(err.Error(), ErrJsonDecode, w)
		return
	}

	// Resolve symlink
	backing, err := filepath.EvalSymlinks(device.Device)
	if err != nil {
		klog.Error("symlink resolve failed:", err)
		WrapResult(err.Error(), ErrExecFailed, w)
		return
	}

	// Prepare DM table
	table := strings.ReplaceAll(device.Table, device.Device, backing)
	dmName := filepath.Base(backing) + "-faulty"

	klog.Infof("Creating faulty device: %s", dmName)

	// Execute dmsetup
	cmd := exec.Command(
		"chroot", "/host",
		"dmsetup", "--noudevsync",
		"create", dmName,
		"--table", table,
	)

	cmd.Stdin = nil

	out, err := cmd.CombinedOutput()
	if err != nil {
		klog.Error("dmsetup create failed:", string(out))
		WrapResult(string(out), ErrExecFailed, w)
		return
	}

	devicePath := "/dev/mapper/" + dmName

	klog.Infof("Created faulty device: %s (%s)", dmName, devicePath)

	WrapResult(devicePath, ErrNone, w)
}

func deleteFaultyDevice(w http.ResponseWriter, r *http.Request) {
	var device Device

	if err := json.NewDecoder(r.Body).Decode(&device); err != nil {
		klog.Error("decode failed:", err)
		WrapResult(err.Error(), ErrJsonDecode, w)
		return
	}

	resolved, err := filepath.EvalSymlinks(device.Device)
	if err != nil {
		klog.Error("resolve failed:", err)
		WrapResult(err.Error(), ErrExecFailed, w)
		return
	}

	base := filepath.Base(resolved)
	dmName := base + "-faulty"

	klog.Infof("Removing DM device: %s", dmName)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cmd := exec.CommandContext(
		ctx,
		"chroot", "/host",
		"dmsetup", "remove", "--retry",
		dmName,
	)

	out, err := cmd.CombinedOutput()

	// Timeout handling
	if ctx.Err() == context.DeadlineExceeded {

		// Check if device actually disappeared
		if _, statErr := os.Stat("/host/dev/mapper/" + dmName); os.IsNotExist(statErr) {

			klog.Warning("dm device removed:", dmName)

			WrapResult("removed "+dmName, ErrNone, w)
			return
		}

		msg := "dmsetup hung and device still exists"
		klog.Error(msg)

		WrapResult(msg, ErrExecFailed, w)
		return
	}

	// check if dmsetup returned error
	if err != nil {

		if strings.Contains(string(out), "No such device") {

			klog.Info("Device already removed:", dmName)

			WrapResult("already removed "+dmName, ErrNone, w)
			return
		}

		klog.Error("dmsetup remove failed:", string(out))

		WrapResult(string(out), ErrExecFailed, w)
		return
	}

	klog.Infof("Removed faulty device: %s", dmName)

	WrapResult("removed "+dmName, ErrNone, w)
}

func controlDevice(w http.ResponseWriter, r *http.Request) {
	var device ControlledDevice
	var err error
	params := make([]string, 2)

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&device); err != nil {
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
	}
	if len(device.Device) == 0 {
		w.WriteHeader(UnprocessableEntityErrorCode)
		_, _ = fmt.Fprint(w, "no device passed")
		klog.Error("no device passed")
		return
	}
	if device.State != "offline" && device.State != "running" {
		w.WriteHeader(UnprocessableEntityErrorCode)
		_, _ = fmt.Fprint(w, "invalid state")
		klog.Error("invalid state")
		return
	}

	// device link which e2e-agent will receive will be like /disk/by-id/scsi-0HC_Volume_29805493 or sdb
	device.Device = fmt.Sprintf("/dev/%s", device.Device)
	klog.Info("Resolving device ", device.Device)
	if device.Device, err = filepath.EvalSymlinks(device.Device); err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to get disk name for dev link:", device.Device, "Error: ", err)
		return
	}
	symLink := strings.Split(device.Device, "/")
	device.Device = symLink[len(symLink)-1]

	klog.Info("Successfully got device name ", device.Device)
	cmdStr := "bash"
	params[0] = "-c"
	params[1] = "echo " + device.State + " > /host/sys/block/" + device.Device + "/device/state"
	klog.Info("running command ", params[1])
	cmd := exec.Command(cmdStr, params...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to run command ", params[1], "Error: ", err)
		return
	}
	_, _ = fmt.Fprint(w, string(output))
	klog.Info(string(output))
}

func killIoEngine(w http.ResponseWriter, r *http.Request) {
	params := make([]string, 2)

	cmdStr := "bash"
	params[0] = "-c"
	params[1] = "MS=$(pidof io-engine) && kill -9 $MS"

	klog.Info("kill io engine")
	cmd := exec.Command(cmdStr, params...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to run command ", params[1], "Error", err)
		return
	}
	_, _ = fmt.Fprint(w, string(output))
	klog.Info(string(output))
}

func killCsiController(w http.ResponseWriter, r *http.Request) {
	params := make([]string, 2)

	cmdStr := "bash"
	params[0] = "-c"
	params[1] = "MS=$(pidof csi-controller) && date && kill -9 $MS"

	klog.Info("kill csi-controller engine")
	cmd := exec.Command(cmdStr, params...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to run command ", params[1], "Error", err)
		return
	}
	klog.Info(string(output))
	WrapResult(string(output), ErrNone, w)
}

func killCsiNode(w http.ResponseWriter, r *http.Request) {
	params := make([]string, 2)

	cmdStr := "bash"
	params[0] = "-c"
	params[1] = "MS=$(pidof csi-node) && date && kill -9 $MS"

	klog.Info("kill csi-node")
	cmd := exec.Command(cmdStr, params...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to run command ", params[1], "Error", err)
		return
	}
	klog.Info(string(output))
	WrapResult(string(output), ErrNone, w)
}

func getDeviceState(w http.ResponseWriter, r *http.Request) {
	var disk DiskPool
	var err error
	params := make([]string, 2)

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&disk); err != nil {
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
	}
	if disk.Disk == "" {
		w.WriteHeader(UnprocessableEntityErrorCode)
		_, _ = fmt.Fprint(w, "no disk pool passed")
		klog.Error("no disk pool passed")
		return
	}
	disk.Disk = fmt.Sprintf("/dev/%s", disk.Disk)
	klog.Info("Resolving device ", disk.Disk)
	if disk.Disk, err = filepath.EvalSymlinks(disk.Disk); err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to get disk name for dev link:", disk.Disk, "Error: ", err)
		return
	}
	symLink := strings.Split(disk.Disk, "/")
	disk.Disk = symLink[len(symLink)-1]

	klog.Info("Successfully got device name ", disk.Disk)

	cmdStr := "bash"
	params[0] = "-c"
	params[1] = "cat /sys/block/" + disk.Disk + "/device/state"
	klog.Info("running command ", params[1])
	cmd := exec.Command(cmdStr, params...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to run command ", params[1], "Error: ", err)
		return
	}
	_, _ = fmt.Fprint(w, string(output))
	klog.Info(string(output))
}

func getHostID(w http.ResponseWriter, r *http.Request) {
	params := make([]string, 2)

	cmdStr := "bash"
	params[0] = "-c"
	params[1] = "cat /sys/class/dmi/id/product_uuid"
	klog.Info("running command ", params[1])
	cmd := exec.Command(cmdStr, params...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to run command ", params[1], "Error: ", err)
		return
	}
	_, _ = fmt.Fprint(w, strings.TrimSpace(string(output)))
	klog.Info(strings.TrimSpace(string(output)))
}

func NvmeConnect(w http.ResponseWriter, r *http.Request) {
	var nvme Nvme
	params := make([]string, 2)

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&nvme); err != nil {
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
	}
	if nvme.TargetIp == "" || nvme.Nqn == "" {
		w.WriteHeader(UnprocessableEntityErrorCode)
		_, _ = fmt.Fprint(w, "no nvme target or nqn passed")
		klog.Error("no nvme target or nqn passed")
		return
	}

	cmdStr := "bash"
	params[0] = "-c"
	if nvme.HostNqn == "" {
		params[1] = fmt.Sprintf("nvme connect -a %s -t tcp -s 8420 -n %s", nvme.TargetIp, nvme.Nqn)
	} else if nvme.HostId != "" {
		params[1] = fmt.Sprintf("nvme connect -a %s -t tcp -s 8420 -n %s -q %s -I %s", nvme.TargetIp, nvme.Nqn, nvme.HostNqn, nvme.HostId)
	} else {
		params[1] = fmt.Sprintf("nvme connect -a %s -t tcp -s 8420 -n %s -q %s", nvme.TargetIp, nvme.Nqn, nvme.HostNqn)
	}
	klog.Info("running command ", params[1])
	cmd := exec.Command(cmdStr, params...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to execute command ", params[1], "Error: ", err)
		return
	}
	_, _ = fmt.Fprint(w, string(output))
	klog.Info(string(output))
}

func NvmeDisconnect(w http.ResponseWriter, r *http.Request) {
	var nvme Nvme
	params := make([]string, 2)

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&nvme); err != nil {
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
	}
	if nvme.Nqn == "" {
		w.WriteHeader(UnprocessableEntityErrorCode)
		_, _ = fmt.Fprint(w, "no nvme nqn passed")
		klog.Error("no nvme nqn passed")
		return
	}

	cmdStr := "bash"
	params[0] = "-c"
	params[1] = fmt.Sprintf("nvme disconnect -n %s", nvme.Nqn)
	klog.Info("executing command ", params[1])
	cmd := exec.Command(cmdStr, params...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to execute command ", params[1], "Error: ", err)
		return
	}
	_, _ = fmt.Fprint(w, string(output))
	klog.Info(string(output))
}

func NvmeList(w http.ResponseWriter, r *http.Request) {
	params := make([]string, 2)

	cmdStr := "bash"
	params[0] = "-c"
	params[1] = "nvme list --output-format=json"

	klog.Info("executing command ", params[1])
	cmd := exec.Command(cmdStr, params...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to execute command ", params[1], "Error: ", err)
		return
	}
	_, _ = fmt.Fprint(w, string(output))
	klog.Info(string(output))
}

func NvmeListSubSys(w http.ResponseWriter, r *http.Request) {
	params := make([]string, 2)

	cmdStr := "bash"
	params[0] = "-c"
	params[1] = "nvme list-subsys --output-format=json"

	klog.Info("executing command ", params[1])
	cmd := exec.Command(cmdStr, params...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to execute command ", params[1], "Error: ", err)
		return
	}
	_, _ = fmt.Fprint(w, string(output))
	klog.Info(string(output))
}

func ChecksumDevice(w http.ResponseWriter, r *http.Request) {
	var device Device
	params := make([]string, 2)

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&device); err != nil {
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
	}
	klog.Info("Running cksum on device, data: %v", device)
	if device.DevicePath == "" {
		w.WriteHeader(UnprocessableEntityErrorCode)
		_, _ = fmt.Fprint(w, "no device path passed")
		klog.Error("no device path passed")
		return
	}

	cmdStr := "bash"
	params[0] = "-c"
	params[1] = fmt.Sprintf("cksum %s", device.DevicePath)

	klog.Info("executing command", params[1])
	cmd := exec.Command(cmdStr, params...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to execute command ", params[1], "Error: ", err)
		return
	}
	_, _ = fmt.Fprint(w, string(output))
	klog.Info(string(output))
}

// Executes bash command on the node and returns value to caller function
func bashLocal(command string) (string, error) {
	var cmd *exec.Cmd
	klog.Info("executing command ", command)
	cmd = exec.Command("bash", []string{"-c", command}...)
	output, err := cmd.CombinedOutput()
	outputString := strings.TrimSpace(string(output))
	if err != nil {
		klog.Error("failed to execute command ", command, "Error: ", err)
		return outputString, fmt.Errorf("run failed, cmd={%s}, error={%v}, output={%s}", cmd, err, string(outputString))
	}
	return outputString, err
}

func freeLoopDevice() (string, error) {
	params := "losetup -f"
	return bashLocal(params)
}

func setupLoopDevice(loDevice string, offset string, devicePath string) (string, error) {
	params := fmt.Sprintf("losetup -o %s %s %s", offset, loDevice, devicePath)
	return bashLocal(params)
}

func detachLoopDevice(loDevice string) (string, error) {
	params := fmt.Sprintf("losetup -d %s", loDevice)
	return bashLocal(params)
}

func FsCheckDevice(w http.ResponseWriter, r *http.Request) {
	var device Device
	var params string

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&device); err != nil {
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
	}
	klog.Info("Running fsck on device, data: %v", device)
	if device.DevicePath == "" {
		w.WriteHeader(UnprocessableEntityErrorCode)
		_, _ = fmt.Fprint(w, "no loop device passed")
		klog.Error("no loop device passed")
		return
	}
	if device.FsType == "" {
		w.WriteHeader(UnprocessableEntityErrorCode)
		_, _ = fmt.Fprint(w, "no fsType passed")
		klog.Error("no fsType passed")
		return
	}

	loDevice, err := freeLoopDevice()
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to get free loop device ", loDevice, "Error: ", err)
		return
	}
	_, err = setupLoopDevice(loDevice, "0", device.DevicePath)
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to setup loop device for", device.DevicePath, "Error: ", err)
		return
	}

	//nolint:staticcheck // QF1003: if-chain kept for simplicity with few filesystem types
	if device.FsType == "ext4" {
		params = fmt.Sprintf("echo $(fsck -n -f %s; echo $?)", loDevice)
	} else if device.FsType == "xfs" {
		params = fmt.Sprintf("echo $(xfs_repair -n %s; echo $?)", loDevice)
	} else if device.FsType == "btrfs" {
		params = fmt.Sprintf("echo $(btrfs check --readonly %s; echo $?)", loDevice)
	} else {
		klog.Error("not a supported filesystem for fscheck", device.FsType, "Error: ", device.FsType)
		return
	}

	outputString, err := bashLocal(params)
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to execute command ", params, "Error: ", err)
		return
	}
	klog.Info(outputString)

	_, err = detachLoopDevice(loDevice)
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to setup loop device for", device.DevicePath, "Error: ", err)
		return
	}
	WrapResult(outputString, ErrNone, w)
}

func createTempDir() (string, error) {
	params := "mktemp -d"
	return bashLocal(params)
}

// This function mounts the xfs device on a temporary directory
//
//	and unmounts the directory
//
// This step is needed to replay metadata from log
func XFSCheckDevice(w http.ResponseWriter, r *http.Request) {
	var device Device
	var params string

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&device); err != nil {
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
	}
	klog.Info("device data: %v", device)
	if device.DevicePath == "" {
		w.WriteHeader(UnprocessableEntityErrorCode)
		_, _ = fmt.Fprint(w, "no loop device passed")
		klog.Error("no loop device passed")
		return
	}
	if device.FsType == "" {
		w.WriteHeader(UnprocessableEntityErrorCode)
		_, _ = fmt.Fprint(w, "no fsType passed")
		klog.Error("no fsType passed")
		return
	}

	tmpDir, err := createTempDir()
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to create temporary directory", "Error: ", err)
		return
	}

	if device.FsType == "xfs" {
		params = fmt.Sprintf("mount %s %s -t xfs; echo $?", device.DevicePath, tmpDir)
	} else {
		klog.Error("not a supported filesystem for fscheck", device.FsType, "Error: ", device.FsType)
		return
	}

	outputString, err := bashLocal(params)
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to execute command ", params, "Error: ", err)
		return
	}
	klog.Info(outputString)

	if device.FsType == "xfs" {
		params = fmt.Sprintf("umount %s; echo $?", tmpDir)
	} else {
		klog.Error("not a supported filesystem for fscheck", device.FsType, "Error: ", device.FsType)
		return
	}

	outputString, err = bashLocal(params)
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to execute command ", params, "Error: ", err)
		return
	}
	klog.Info(outputString)

	if device.FsType == "xfs" {
		params = fmt.Sprintf("rm -rf %s; echo $?", tmpDir)
	} else {
		klog.Error("not a supported filesystem for fscheck", device.FsType, "Error: ", device.FsType)
		return
	}

	outputString, err = bashLocal(params)
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to execute command ", params, "Error: ", err)
		return
	}
	klog.Info(outputString)

	if device.FsType == "xfs" {
		params = fmt.Sprintf("echo $(xfs_repair -n %s; echo $?)", device.DevicePath)
	} else {
		klog.Error("not a supported filesystem for fscheck", device.FsType, "Error: ", device.FsType)
		return
	}

	outputString, err = bashLocal(params)
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to execute command ", params, "Error: ", err)
		return
	}
	klog.Info(outputString)
	WrapResult(outputString, ErrNone, w)
}

func listMountPoint(devicePath string) (string, error) {
	nvmeName := strings.TrimPrefix(devicePath, "/dev/")
	params := fmt.Sprintf("lsblk -fa|grep %s|awk '{print $NF}'", nvmeName)
	return bashLocal(params)
}

func FsFreezeDevice(w http.ResponseWriter, r *http.Request) {
	var device Device
	var params string

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&device); err != nil {
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
	}
	klog.Info("Running fs freeze on device, data: %v", device)
	if device.DevicePath == "" {
		w.WriteHeader(UnprocessableEntityErrorCode)
		_, _ = fmt.Fprint(w, "no device path passed")
		klog.Error("no device path passed")
		return
	}

	mountPath, err := listMountPoint(device.DevicePath)
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to list mountpoint for ", device.DevicePath, "Error: ", err)
		return
	}

	params = fmt.Sprintf("echo $(fsfreeze -f %s; echo $?)", mountPath)
	outputString, err := bashLocal(params)
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to execute command ", params, "Error: ", err)
		return
	}
	klog.Info(outputString)
	WrapResult(outputString, ErrNone, w)
}

func FsUnfreezeDevice(w http.ResponseWriter, r *http.Request) {
	var device Device
	var params string

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&device); err != nil {
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
	}
	klog.Info("Running fs unfreeze on device, data: %v", device)
	if device.DevicePath == "" {
		w.WriteHeader(UnprocessableEntityErrorCode)
		_, _ = fmt.Fprint(w, "no device path passed")
		klog.Error("no device path passed")
		return
	}

	mountPath, err := listMountPoint(device.DevicePath)
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to list mountpoint for ", device.DevicePath, "Error: ", err)
		return
	}

	params = fmt.Sprintf("echo $(fsfreeze -u %s; echo $?)", mountPath)
	outputString, err := bashLocal(params)
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to execute command ", params, "Error: ", err)
		return
	}
	klog.Info(outputString)
	WrapResult(outputString, ErrNone, w)
}

func ListDevice(w http.ResponseWriter, r *http.Request) {
	params := make([]string, 2)

	cmdStr := "bash"
	params[0] = "-c"
	params[1] = "ls /dev/"
	klog.Info("executing command ", params[1])
	cmd := exec.Command(cmdStr, params...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to execute command ", params[1], "Error: ", err)
		return
	}
	_, _ = fmt.Fprint(w, string(output))
}

func flushDiskWriteCache(w http.ResponseWriter, r *http.Request) {
	params := make([]string, 2)

	cmdStr := "bash"
	params[0] = "-c"
	params[1] = "sync"
	klog.Info("executing command ", params[1])
	cmd := exec.Command(cmdStr, params...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		klog.Error("failed to execute command ", params[1], "Error: ", err)
		WrapResult(fmt.Sprintf("failed to execute: %s, got error: %v", params[1], err), ErrJsonEncode, w)
		return
	}
	klog.Info(string(output))
	WrapResult(string(output), ErrNone, w)
}

func ZeroingDisk(w http.ResponseWriter, r *http.Request) {
	var disk Disk
	params := make([]string, 2)

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&disk); err != nil {
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
	}
	if disk.Device == "" || disk.SeekParam == "" || disk.BlockSizeParam == "" {
		w.WriteHeader(UnprocessableEntityErrorCode)
		_, _ = fmt.Fprint(w, "no device or seek or block param passed")
		klog.Error("no device or seek or block param passed")
		return
	}

	cmdStr := "bash"
	params[0] = "-c"
	params[1] = fmt.Sprintf("dd if=/dev/zero of=/host%s count=1 oflag=direct %s %s && sync",
		disk.Device,
		disk.SeekParam,
		disk.BlockSizeParam,
	)
	klog.Info("executing command ", params[1])
	cmd := exec.Command(cmdStr, params...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to execute command ", params[1], "Error", err)
		return
	}
	_, _ = fmt.Fprint(w, string(output))
	klog.Info(string(output))
}

func ZeroingHugePages(w http.ResponseWriter, r *http.Request) {
	params := make([]string, 2)
	cmdStr := "bash"
	params[0] = "-c"
	params[1] = "echo 0 | sudo tee /proc/sys/vm/nr_hugepages"

	klog.Info("executing command ", params[1])
	cmd := exec.Command(cmdStr, params...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		klog.Error("failed to execute command ", params[1], "Error", err)
		msg := fmt.Sprintf("failed to execute command %s, err: %v", params[1], err)
		WrapResult(msg, ErrGeneral, w)
		return
	}
	klog.Info(string(output))
	WrapResult(string(output), ErrNone, w)
}

func Parted(w http.ResponseWriter, r *http.Request) {
	var cmdline CmdList
	var cmd *exec.Cmd
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&cmdline); err != nil {
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
	}
	if len(cmdline.Cmd) == 0 {
		w.WriteHeader(UnprocessableEntityErrorCode)
		_, _ = fmt.Fprint(w, "no command passed")
		klog.Error("no command passed")
		return
	}

	cmdArgs := strings.Split(cmdline.Cmd, " ")
	klog.Info("executing command ", cmdArgs)
	cmd = exec.Command("parted", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to execute command", cmdArgs, "Error: ", err)
	} else {
		_, _ = fmt.Fprint(w, string(output))
		klog.Info(string(output))
	}
}

func bashit(command string, w http.ResponseWriter, r *http.Request) {
	var cmd *exec.Cmd
	klog.Info("executing command ", command)
	cmd = exec.Command("bash", []string{"-c", command}...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to execute command ", command, "Error", err)
	} else {
		_, _ = fmt.Fprint(w, string(output))
		klog.Info(string(output))
	}
}

func Findmnt(w http.ResponseWriter, r *http.Request) {
	bashit("findmnt --json", w, r)
}

func Lsblk(w http.ResponseWriter, r *http.Request) {
	bashit("lsblk --json", w, r)
}

func Dmesg(w http.ResponseWriter, r *http.Request) {
	bashit("dmesg -T", w, r)
}

func Getsyslog(w http.ResponseWriter, r *http.Request) {
	bashit("cat /host/var/log/syslog", w, r)
}

func ListReservation(w http.ResponseWriter, r *http.Request) {
	var device Device
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&device); err != nil {
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
	}
	if device.DevicePath == "" {
		w.WriteHeader(UnprocessableEntityErrorCode)
		_, _ = fmt.Fprint(w, "no device path passed")
		klog.Error("no device path passed")
		return
	}
	cmd := fmt.Sprintf("nvme resv-report %s -c 1 --output-format=json", device.DevicePath)
	bashit(cmd, w, r)
}

func NvmeConnectWithHostId(w http.ResponseWriter, r *http.Request) {
	var nvme Nvme
	params := make([]string, 2)

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&nvme); err != nil {
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
	}
	if nvme.TargetIp == "" || nvme.Nqn == "" || nvme.HostId == "" {
		w.WriteHeader(UnprocessableEntityErrorCode)
		_, _ = fmt.Fprint(w, "no nvme target or nqn or host id passed")
		klog.Error("no nvme target or nqn or host id passed")
		return
	}

	cmdStr := "bash"
	params[0] = "-c"
	params[1] = fmt.Sprintf("nvme connect -a %s -t tcp -s 8420 -n %s -q %s -I %s", nvme.TargetIp, nvme.Nqn, nvme.HostNqn, nvme.HostId)
	klog.Info("executing command ", params[1])
	cmd := exec.Command(cmdStr, params...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to execute command ", params[1], "Error", err)
		return
	}
	_, _ = fmt.Fprint(w, string(output))
	klog.Info(string(output))
}

func BlkDiscard(w http.ResponseWriter, r *http.Request) {
	var data blkDiscard
	params := make([]string, 2)

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&data); err != nil {
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
		return
	}
	if data.Device == "" {
		w.WriteHeader(UnprocessableEntityErrorCode)
		_, _ = fmt.Fprint(w, "no device param passed")
		klog.Error("no device  param passed")
		return
	}

	cmdStr := "bash"
	params[0] = "-c"
	params[1] = fmt.Sprintf("blkdiscard %s %s", data.Device, data.Options)

	klog.Info("executing command ", params[1])
	cmd := exec.Command(cmdStr, params...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to execute command ", params[1], "Error", err)
		return
	}
	_, _ = fmt.Fprint(w, string(output))
	klog.Info(string(output))
}

func Cmp(w http.ResponseWriter, r *http.Request) {
	var paths CmpPaths
	params := make([]string, 2)

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&paths); err != nil {
		_, _ = fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
	}
	klog.Info("Running cmp on device, data: %v", paths)
	if paths.Path1 == "" {
		w.WriteHeader(UnprocessableEntityErrorCode)
		_, _ = fmt.Fprint(w, "no path1 passed")
		klog.Error("no device path1 passed")
		return
	}
	if paths.Path2 == "" {
		w.WriteHeader(UnprocessableEntityErrorCode)
		_, _ = fmt.Fprint(w, "no path2 passed")
		klog.Error("no path2 passed")
		return
	}
	cmdStr := "bash"
	params[0] = "-c"
	params[1] = fmt.Sprintf("cmp -b %s %s", paths.Path1, paths.Path2)

	var errCode = ErrNone

	klog.Info("executing command ", params[1])
	cmd := exec.Command(cmdStr, params...)
	output, err := cmd.CombinedOutput()
	b64out := base64.StdEncoding.EncodeToString(output)
	if err != nil {
		errCode = ErrGeneral
		klog.Error("failed command ", params[1], " Error: ", err)
	}
	klog.Info(string(output))
	WrapResult(b64out, errCode, w)
}

// EnableNetworkInterface enable network interface
func EnableNetworkInterface(w http.ResponseWriter, r *http.Request) {
	var msg string
	var iface NetworkInterface
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&iface); err != nil {
		msg = fmt.Sprintf("failed to read JSON encoded data, Error: %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrJsonDecode, w)
		return
	}

	if iface.NetworkInterface == "" {
		msg = "no network interface name passed"
		klog.Error(msg)
		WrapResult(msg, UnprocessableEntityErrorCode, w)
		return
	}
	klog.Info("enable network interface, data: %v", iface)

	iFaceUpCommand := fmt.Sprintf("ifconfig %s up", iface.NetworkInterface)
	output, err := bashLocal(iFaceUpCommand)
	if err != nil {
		msg = fmt.Sprintf("cannot enable network interface, Error %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrExecFailed, w)
		return
	}
	WrapResult(string(output), ErrNone, w)
}

// DisableNetworkInterface disable network interface
func DisableNetworkInterface(w http.ResponseWriter, r *http.Request) {
	var msg string
	var iface NetworkInterface
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&iface); err != nil {
		msg = fmt.Sprintf("failed to read JSON encoded data, Error: %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrJsonDecode, w)
		return
	}
	if iface.NetworkInterface == "" {
		msg = "no interface name passed"
		klog.Error(msg)
		WrapResult(msg, UnprocessableEntityErrorCode, w)
		return
	}
	klog.Info("disable network interface, data: %v", iface)

	iFaceDownCommand := fmt.Sprintf("ifconfig %s down", iface.NetworkInterface)
	output, err := bashLocal(iFaceDownCommand)
	if err != nil {
		msg = fmt.Sprintf("cannot disable network interface, Error %s", err.Error())
		klog.Error(msg)
		WrapResult(msg, ErrExecFailed, w)
		return
	}
	WrapResult(string(output), ErrNone, w)
}

func resolveBlockDevice(dev string) (string, error) {
	resolved, err := filepath.EvalSymlinks(dev)
	if err != nil {
		return "", err
	}
	return resolved, nil
}

type DmCreateResponse struct {
	Device string `json:"device"`
}

func createPassThroughDevice(w http.ResponseWriter, r *http.Request) {
	var req DmDevice

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		klog.Error("decode failed:", err)
		WrapResult(err.Error(), ErrJsonDecode, w)
		return
	}

	backing, err := resolveBlockDevice(req.Device)
	if err != nil {
		klog.Error("resolve failed:", err)
		WrapResult(err.Error(), ErrExecFailed, w)
		return
	}

	name := filepath.Base(backing) + "-timeout"
	table := fmt.Sprintf("0 %d linear %s 0", req.Sectors, backing)

	klog.Infof("Creating passthrough device: %s", name)

	cmd := exec.Command("dmsetup", "--noudevsync", "create", name, "--table", table)
	out, err := cmd.CombinedOutput()
	if err != nil {
		klog.Error("dmsetup create failed:", string(out))
		WrapResult(string(out), ErrExecFailed, w)
		return
	}

	klog.Infof("Created passthrough device: %s", name)

	resp := DmCreateResponse{
		Device: name,
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		klog.Error("json marshal failed:", err)
		WrapResult(err.Error(), ErrJsonDecode, w)
		return
	}

	WrapResult(string(respBytes), ErrNone, w)
}

func suspendDevice(w http.ResponseWriter, r *http.Request) {

	var req DmDevice

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		klog.Error("decode failed:", err)
		WrapResult(err.Error(), ErrJsonDecode, w)
		return
	}
	name := filepath.Base(req.Device)

	klog.Infof("Suspending device: %s", name)

	cmd := exec.Command(
		"chroot", "/host",
		"dmsetup", "suspend",
		name,
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		klog.Error("dmsetup suspend failed:", string(out))
		WrapResult(string(out), ErrExecFailed, w)
		return
	}

	klog.Infof("Suspended device: %s", name)
	WrapResult("suspended "+name, ErrNone, w)
}

func isDeviceActive(name string) bool {
	out, err := exec.Command(
		"chroot", "/host",
		"dmsetup", "info",
		"-c",
		"--noheadings",
		"-o", "name,attr",
		name,
	).CombinedOutput()

	if err != nil {
		return false
	}

	line := strings.TrimSpace(string(out))
	if line == "" {
		return false
	}

	// Extract attr (works for both "name:attr" and "name attr")
	attr := ""
	if i := strings.Index(line, ":"); i != -1 {
		attr = line[i+1:]
	} else {
		parts := strings.Fields(line)
		if len(parts) > 1 {
			attr = parts[1]
		}
	}

	// active if no 's'
	return attr != "" && !strings.Contains(attr, "s")
}

func resumeDevice(w http.ResponseWriter, r *http.Request) {
	var req DmDevice

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		klog.Error("decode failed:", err)
		WrapResult(err.Error(), ErrJsonDecode, w)
		return
	}

	name := filepath.Base(req.Device)
	klog.Infof("Resuming device (noflush): %s", name)

	// Timeout to prevent hanging forever
	ctx, cancel := context.WithTimeout(context.Background(), dmsetupTimeout)
	defer cancel()

	start := time.Now()

	cmd := exec.CommandContext(
		ctx,
		"chroot", "/host",
		"dmsetup", "resume",
		"--noflush",
		name,
	)

	out, err := cmd.CombinedOutput()
	duration := time.Since(start)

	// Handle timeout and verify actual device state
	if ctx.Err() == context.DeadlineExceeded {

		if isDeviceActive(name) {
			klog.Infof("Device %s is active ", name)
			WrapResult("resumed "+name+"", ErrNone, w)
			return
		}

		klog.Errorf("dmsetup resume timed out and device still suspended: %s", string(out))
		WrapResult("dmsetup resume timeout", ErrExecFailed, w)
		return
	}

	// Handle execution error → still verify state
	if err != nil {
		klog.Warningf("dmsetup resume returned error for %s: %s, verifying state...", name, string(out))

		if isDeviceActive(name) {
			klog.Infof("Device %s is active despite error", name)
			WrapResult("resumed "+name+"", ErrNone, w)
			return
		}

		klog.Errorf("dmsetup resume failed for %s: %s", name, string(out))
		WrapResult(string(out), ErrExecFailed, w)
		return
	}

	klog.Infof("Resumed device: %s (took %v)", name, duration)
	WrapResult("resumed "+name, ErrNone, w)
}

func removeDevice(w http.ResponseWriter, r *http.Request) {
	var req DmDevice

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		klog.Error("decode failed:", err)
		WrapResult(err.Error(), ErrJsonDecode, w)
		return
	}

	name := filepath.Base(req.Device)
	klog.Infof("Removing passthrough device: %s", name)

	// Use timeout to avoid blocking
	ctx, cancel := context.WithTimeout(context.Background(), dmsetupTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "dmsetup", "remove", "--deferred", name)
	err := cmd.Run()

	if ctx.Err() == context.DeadlineExceeded {
		klog.Warningf("dmsetup remove timed out for %s", name)
		WrapResult("remove triggered "+name+" (timeout)", ErrNone, w)
		return
	}

	if err != nil {
		klog.Warningf("dmsetup remove error for %s: %v", name, err)
		WrapResult("remove attempted "+name, ErrNone, w)
		return
	}

	klog.Infof("Removed passthrough device: %s", name)
	WrapResult("removed "+name, ErrNone, w)
}

func getDeviceSizeInSectors(w http.ResponseWriter, r *http.Request) {
	var req DmDevice

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		klog.Error("decode failed:", err)
		WrapResult(err.Error(), ErrJsonDecode, w)
		return
	}

	backing, err := resolveBlockDevice(req.Device)
	if err != nil {
		klog.Error("resolve failed:", err)
		WrapResult(err.Error(), ErrExecFailed, w)
		return
	}

	cmd := exec.Command("blockdev", "--getsz", backing)

	out, err := cmd.CombinedOutput()
	if err != nil {
		klog.Error("blockdev failed:", string(out))
		WrapResult(string(out), ErrExecFailed, w)
		return
	}

	size := strings.TrimSpace(string(out))

	klog.Infof("Device %s size(sectors): %s", backing, size)
	WrapResult(size, ErrNone, w)
}

func dmSuspend(dmName string) error {
	cmd := exec.Command("chroot", "/host", "dmsetup", "--noudevsync", "suspend", dmName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("suspend failed: %s", string(out))
	}
	return nil
}

func dmResume(dmName string) error {
	cmd := exec.Command("chroot", "/host", "dmsetup", "--noudevsync", "resume", dmName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("resume failed: %s", string(out))
	}
	return nil
}

func dmReload(dmName, table string) error {
	cmd := exec.Command("chroot", "/host", "dmsetup", "--noudevsync", "reload", dmName, "--table", table)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("reload failed: %s", string(out))
	}
	return nil
}

func reloadDevice(w http.ResponseWriter, r *http.Request) {
	var device Device

	// Decode request
	if err := json.NewDecoder(r.Body).Decode(&device); err != nil {
		klog.Error("decode failed:", err)
		WrapResult(err.Error(), ErrJsonDecode, w)
		return
	}

	dmName := filepath.Base(device.Device)
	table := strings.TrimSpace(device.Table)

	klog.Infof("Reloading device: %s with table: %s", dmName, table)

	cmd := exec.Command("chroot", "/host", "dmsetup", "info", dmName)
	if out, err := cmd.CombinedOutput(); err != nil {
		klog.Error("dmsetup info failed:", string(out))
		WrapResult(string(out), ErrExecFailed, w)
		return
	}

	if err := dmSuspend(dmName); err != nil {
		klog.Error(err, "suspend failed", "device", dmName)
		WrapResult(err.Error(), ErrExecFailed, w)
		return
	}

	if err := dmReload(dmName, table); err != nil {
		klog.Error(err, "reload failed", "device", dmName)

		if resumeErr := dmResume(dmName); resumeErr != nil {
			klog.Error(resumeErr, "resume failed after reload failure", "device", dmName)
		}

		WrapResult(err.Error(), ErrExecFailed, w)
		return
	}

	if err := dmResume(dmName); err != nil {
		klog.Error(err, "resume failed", "device", dmName)
		WrapResult(err.Error(), ErrExecFailed, w)
		return
	}

	klog.Infof("Reloaded device successfully: %s", dmName)

	resp := DmCreateResponse{
		Device: device.Device,
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		klog.Error("json marshal failed:", err)
		WrapResult(err.Error(), ErrJsonDecode, w)
		return
	}

	WrapResult(string(respBytes), ErrNone, w)
}
