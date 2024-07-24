package main

import (
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

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome home!\n")
}

type CmdList struct {
	Cmd string `json:"cmd"`
}

var Version = "undefined"

const (
	InternalServerErrorCode      = 500
	UnprocessableEntityErrorCode = 422
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
	router.HandleFunc("/devicecontrol", controlDevice).Methods("POST")
	router.HandleFunc("/killioengine", killIoEngine).Methods("POST")
	router.HandleFunc("/killCsiController", killCsiController).Methods("POST")
	router.HandleFunc("/killCsiNode", killCsiNode).Methods("POST")
	router.HandleFunc("/getdevicestate", getDeviceState).Methods("POST")
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
	fmt.Fprint(w, "Graceful reboots are not yet supported")
}

func dropConnectionsFromNodes(w http.ResponseWriter, r *http.Request) {
	var list NodeList
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&list); err != nil {
		fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
	}
	klog.Info("Dropping connection from nodes ", list.Nodes, list.NetworkInterface)
	if err := DropConnectionsFromNodes(list.Nodes, list.NetworkInterface); err != nil {
		w.WriteHeader(InternalServerErrorCode)
		fmt.Fprint(w, err.Error())
		klog.Error("failed to drop connection from nodes:", list.Nodes, "Error: ", err)
		return
	}
	klog.Info("Successfully stopped network services")
}

func acceptConnectionsFromNodes(w http.ResponseWriter, r *http.Request) {
	var list NodeList
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&list); err != nil {
		fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
	}
	klog.Info("Accept connection from nodes ", list.Nodes, list.NetworkInterface)
	err := AcceptConnectionsFromNodes(list.Nodes, list.NetworkInterface)
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		fmt.Fprint(w, err.Error())
		klog.Error("failed to accept connection from nodes:", list.Nodes, "Error: ", err)
		return
	}
	fmt.Fprint(w, "Successfully started network services\n")
	klog.Info("Successfully started network services")
}

func createFaultyDevice(w http.ResponseWriter, r *http.Request) {
	var (
		device Device
		cmd    *exec.Cmd
	)
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&device); err != nil {
		fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
	}
	if len(device.Device) == 0 {
		w.WriteHeader(UnprocessableEntityErrorCode)
		fmt.Fprint(w, "no device passed")
		klog.Error("no device passed")
		return
	}
	if len(device.Table) == 0 {
		w.WriteHeader(UnprocessableEntityErrorCode)
		fmt.Fprint(w, "no table passed")
		klog.Error("no table passed")
		return
	}
	klog.Info("create faulty device ", device)
	f, err := os.Create("table")
	if err != nil {
		klog.Error(err)
		return
	}
	_, err = f.WriteString(device.Table)
	if err != nil {
		klog.Error(err)
		f.Close()
		return
	}
	err = f.Close()
	if err != nil {
		klog.Error(err)
		return
	}
	devName := strings.Split(device.Device, "/")

	cmdStr := "dmsetup create" + " " + devName[2] + " " + "table"
	cmdArgs := strings.Split(cmdStr, " ")
	cmdName := cmdArgs[0]
	if len(cmdArgs) > 1 {
		cmd = exec.Command(cmdName, cmdArgs[1:]...)
	} else {
		cmd = exec.Command(cmdName)
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		fmt.Fprint(w, err.Error())
		klog.Error(err)
	} else {
		fmt.Fprint(w, string(output))
		klog.Info(string(output))
	}
}

func deleteFaultyDevice(w http.ResponseWriter, r *http.Request) {
	var (
		device Device
		cmd    *exec.Cmd
	)
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&device); err != nil {
		fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
	}
	if len(device.Device) == 0 {
		w.WriteHeader(UnprocessableEntityErrorCode)
		fmt.Fprint(w, "no device passed")
		klog.Error("no device passed")
		return
	}
	klog.Info("delete faulty device ", device)
	devName := strings.Split(device.Device, "/")

	cmdStr := "dmsetup remove" + " " + devName[2]
	cmdArgs := strings.Split(cmdStr, " ")
	cmdName := cmdArgs[0]
	if len(cmdArgs) > 1 {
		cmd = exec.Command(cmdName, cmdArgs[1:]...)
	} else {
		cmd = exec.Command(cmdName)
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(InternalServerErrorCode)
		fmt.Fprint(w, err.Error())
		klog.Error("failed to delete faulty deice ", device, "Error: ", err)
	} else {
		fmt.Fprint(w, string(output))
		klog.Info(string(output))
	}
}

func controlDevice(w http.ResponseWriter, r *http.Request) {
	var device ControlledDevice
	var err error
	params := make([]string, 2)

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&device); err != nil {
		fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
	}
	if len(device.Device) == 0 {
		w.WriteHeader(UnprocessableEntityErrorCode)
		fmt.Fprint(w, "no device passed")
		klog.Error("no device passed")
		return
	}
	if device.State != "offline" && device.State != "running" {
		w.WriteHeader(UnprocessableEntityErrorCode)
		fmt.Fprint(w, "invalid state")
		klog.Error("invalid state")
		return
	}

	// device link which e2e-agent will receive will be like /disk/by-id/scsi-0HC_Volume_29805493 or sdb
	device.Device = fmt.Sprintf("/dev/%s", device.Device)
	klog.Info("Resolving device ", device.Device)
	if device.Device, err = filepath.EvalSymlinks(device.Device); err != nil {
		w.WriteHeader(InternalServerErrorCode)
		fmt.Fprint(w, err.Error())
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
		fmt.Fprint(w, err.Error())
		klog.Error("failed to run command ", params[1], "Error: ", err)
		return
	}
	fmt.Fprint(w, string(output))
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
		fmt.Fprint(w, err.Error())
		klog.Error("failed to run command ", params[1], "Error", err)
		return
	}
	fmt.Fprint(w, string(output))
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
		fmt.Fprint(w, err.Error())
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
		fmt.Fprint(w, err.Error())
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
		fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
	}
	if disk.Disk == "" {
		w.WriteHeader(UnprocessableEntityErrorCode)
		fmt.Fprint(w, "no disk pool passed")
		klog.Error("no disk pool passed")
		return
	}
	disk.Disk = fmt.Sprintf("/dev/%s", disk.Disk)
	klog.Info("Resolving device ", disk.Disk)
	if disk.Disk, err = filepath.EvalSymlinks(disk.Disk); err != nil {
		w.WriteHeader(InternalServerErrorCode)
		fmt.Fprint(w, err.Error())
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
		fmt.Fprint(w, err.Error())
		klog.Error("failed to run command ", params[1], "Error: ", err)
		return
	}
	fmt.Fprint(w, string(output))
	klog.Info(string(output))
}

func NvmeConnect(w http.ResponseWriter, r *http.Request) {
	var nvme Nvme
	params := make([]string, 2)

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&nvme); err != nil {
		fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
	}
	if nvme.TargetIp == "" || nvme.Nqn == "" {
		w.WriteHeader(UnprocessableEntityErrorCode)
		fmt.Fprint(w, "no nvme target or nqn passed")
		klog.Error("no nvme target or nqn passed")
		return
	}

	cmdStr := "bash"
	params[0] = "-c"
	if nvme.HostNqn == "" {
		params[1] = fmt.Sprintf("nvme connect -a %s -t tcp -s 8420 -n %s", nvme.TargetIp, nvme.Nqn)
	} else {
		params[1] = fmt.Sprintf("nvme connect -a %s -t tcp -s 8420 -n %s -q %s", nvme.TargetIp, nvme.Nqn, nvme.HostNqn)
	}
	klog.Info("running command ", params[1])
	cmd := exec.Command(cmdStr, params...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		fmt.Fprint(w, err.Error())
		klog.Error("failed to execute command ", params[1], "Error: ", err)
		return
	}
	fmt.Fprint(w, string(output))
	klog.Info(string(output))
}

func NvmeDisconnect(w http.ResponseWriter, r *http.Request) {
	var nvme Nvme
	params := make([]string, 2)

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&nvme); err != nil {
		fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
	}
	if nvme.Nqn == "" {
		w.WriteHeader(UnprocessableEntityErrorCode)
		fmt.Fprint(w, "no nvme nqn passed")
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
		fmt.Fprint(w, err.Error())
		klog.Error("failed to execute command ", params[1], "Error: ", err)
		return
	}
	fmt.Fprint(w, string(output))
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
		fmt.Fprint(w, err.Error())
		klog.Error("failed to execute command ", params[1], "Error: ", err)
		return
	}
	fmt.Fprint(w, string(output))
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
		fmt.Fprint(w, err.Error())
		klog.Error("failed to execute command ", params[1], "Error: ", err)
		return
	}
	fmt.Fprint(w, string(output))
	klog.Info(string(output))
}

func ChecksumDevice(w http.ResponseWriter, r *http.Request) {
	var device Device
	params := make([]string, 2)

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&device); err != nil {
		fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
	}
	klog.Info("Running cksum on device, data: %v", device)
	if device.DevicePath == "" {
		w.WriteHeader(UnprocessableEntityErrorCode)
		fmt.Fprint(w, "no device path passed")
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
		fmt.Fprint(w, err.Error())
		klog.Error("failed to execute command ", params[1], "Error: ", err)
		return
	}
	fmt.Fprint(w, string(output))
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
	params := fmt.Sprintf("losetup -o 5M %s %s", loDevice, devicePath)
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
		fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
	}
	klog.Info("Running fsck on device, data: %v", device)
	if device.DevicePath == "" {
		w.WriteHeader(UnprocessableEntityErrorCode)
		fmt.Fprint(w, "no loop device passed")
		klog.Error("no loop device passed")
		return
	}
	if device.FsType == "" {
		w.WriteHeader(UnprocessableEntityErrorCode)
		fmt.Fprint(w, "no fsType passed")
		klog.Error("no fsType passed")
		return
	}

	loDevice, err := freeLoopDevice()
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		fmt.Fprint(w, err.Error())
		klog.Error("failed to get free loop device ", loDevice, "Error: ", err)
		return
	}
	_, err = setupLoopDevice(loDevice, "5M", device.DevicePath)
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		fmt.Fprint(w, err.Error())
		klog.Error("failed to setup loop device for", device.DevicePath, "Error: ", err)
		return
	}

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
		fmt.Fprint(w, err.Error())
		klog.Error("failed to execute command ", params, "Error: ", err)
		return
	}
	klog.Info(outputString)

	_, err = detachLoopDevice(loDevice)
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		fmt.Fprint(w, err.Error())
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
		fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
	}
	klog.Info("device data: %v", device)
	if device.DevicePath == "" {
		w.WriteHeader(UnprocessableEntityErrorCode)
		fmt.Fprint(w, "no loop device passed")
		klog.Error("no loop device passed")
		return
	}
	if device.FsType == "" {
		w.WriteHeader(UnprocessableEntityErrorCode)
		fmt.Fprint(w, "no fsType passed")
		klog.Error("no fsType passed")
		return
	}

	tmpDir, err := createTempDir()
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		fmt.Fprint(w, err.Error())
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
		fmt.Fprint(w, err.Error())
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
		fmt.Fprint(w, err.Error())
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
		fmt.Fprint(w, err.Error())
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
		fmt.Fprint(w, err.Error())
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
		fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
	}
	klog.Info("Running fs freeze on device, data: %v", device)
	if device.DevicePath == "" {
		w.WriteHeader(UnprocessableEntityErrorCode)
		fmt.Fprint(w, "no device path passed")
		klog.Error("no device path passed")
		return
	}

	mountPath, err := listMountPoint(device.DevicePath)
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		fmt.Fprint(w, err.Error())
		klog.Error("failed to list mountpoint for ", device.DevicePath, "Error: ", err)
		return
	}

	params = fmt.Sprintf("echo $(fsfreeze -f %s; echo $?)", mountPath)
	outputString, err := bashLocal(params)
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		fmt.Fprint(w, err.Error())
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
		fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
	}
	klog.Info("Running fs unfreeze on device, data: %v", device)
	if device.DevicePath == "" {
		w.WriteHeader(UnprocessableEntityErrorCode)
		fmt.Fprint(w, "no device path passed")
		klog.Error("no device path passed")
		return
	}

	mountPath, err := listMountPoint(device.DevicePath)
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		fmt.Fprint(w, err.Error())
		klog.Error("failed to list mountpoint for ", device.DevicePath, "Error: ", err)
		return
	}

	params = fmt.Sprintf("echo $(fsfreeze -u %s; echo $?)", mountPath)
	outputString, err := bashLocal(params)
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		fmt.Fprint(w, err.Error())
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
		fmt.Fprint(w, err.Error())
		klog.Error("failed to execute command ", params[1], "Error: ", err)
		return
	}
	fmt.Fprint(w, string(output))
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
		fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
	}
	if disk.Device == "" || disk.SeekParam == "" || disk.BlockSizeParam == "" {
		w.WriteHeader(UnprocessableEntityErrorCode)
		fmt.Fprint(w, "no device or seek or block param passed")
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
		fmt.Fprint(w, err.Error())
		klog.Error("failed to execute command ", params[1], "Error", err)
		return
	}
	fmt.Fprint(w, string(output))
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
		fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
	}
	if len(cmdline.Cmd) == 0 {
		w.WriteHeader(UnprocessableEntityErrorCode)
		fmt.Fprint(w, "no command passed")
		klog.Error("no command passed")
		return
	}

	cmdArgs := strings.Split(cmdline.Cmd, " ")
	klog.Info("executing command ", cmdArgs)
	cmd = exec.Command("parted", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		w.WriteHeader(InternalServerErrorCode)
		fmt.Fprint(w, err.Error())
		klog.Error("failed to execute command", cmdArgs, "Error: ", err)
	} else {
		fmt.Fprint(w, string(output))
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
		fmt.Fprint(w, err.Error())
		klog.Error("failed to execute command ", command, "Error", err)
	} else {
		fmt.Fprint(w, string(output))
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
		fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
	}
	if device.DevicePath == "" {
		w.WriteHeader(UnprocessableEntityErrorCode)
		fmt.Fprint(w, "no device path passed")
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
		fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
	}
	if nvme.TargetIp == "" || nvme.Nqn == "" || nvme.HostId == "" {
		w.WriteHeader(UnprocessableEntityErrorCode)
		fmt.Fprint(w, "no nvme target or nqn or host id passed")
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
		fmt.Fprint(w, err.Error())
		klog.Error("failed to execute command ", params[1], "Error", err)
		return
	}
	fmt.Fprint(w, string(output))
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
		fmt.Fprint(w, err.Error())
		klog.Error("failed to read JSON encoded data, Error: ", err)
	}
	klog.Info("Running cmp on device, data: %v", paths)
	if paths.Path1 == "" {
		w.WriteHeader(UnprocessableEntityErrorCode)
		fmt.Fprint(w, "no path1 passed")
		klog.Error("no device path1 passed")
		return
	}
	if paths.Path2 == "" {
		w.WriteHeader(UnprocessableEntityErrorCode)
		fmt.Fprint(w, "no path2 passed")
		klog.Error("no path2 passed")
		return
	}
	cmdStr := "bash"
	params[0] = "-c"
	params[1] = fmt.Sprintf("cmp -b %s %s", paths.Path1, paths.Path2)

	var errCode E2eAgentErrcode = ErrNone

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
