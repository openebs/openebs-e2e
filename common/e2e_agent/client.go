package e2e_agent

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/e2e_config"
	"github.com/openebs/openebs-e2e/common/k8s_portforward"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// RestPort is the port on which e2e-agent is listening
const RestPort = 10012

// NodeList is the list of nodes to be passed to e2e-agent
type NodeList struct {
	Nodes            []string `json:"nodes"`
	NetworkInterface string   `json:"networkInterface"`
}

type CmdList struct {
	Cmd string `json:"cmd"`
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

type DiskPool struct {
	Disk string `json:"disk"`
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

type Lvm struct {
	Pv                          string `json:"pv"`                          // Physical volume
	Vg                          string `json:"vg"`                          // Volume group
	ThinPoolAutoExtendThreshold int    `json:"thinPoolAutoExtendThreshold"` // thin pool auto extend threshold
	ThinPoolAutoExtendPercent   int    `json:"ThinPoolAutoExtendPercent"`   // thin pool auto extend percent
}

type LoopDevice struct {
	// Size in bytes
	Size int64 `json:"size"`
	// The backing image name
	// eg: fake123
	ImageName string `json:"imageName"`
	// Image directory
	// eg: /tmp
	ImgDir string `json:"imgDir"`
	// the disk name
	// eg: /tmp/loop9002
	DiskPath string `json:"diskPath"`
}

func sendRequest(reqType, url string, data interface{}) error {
	_, err := sendRequestGetResponse(reqType, url, data, true)
	return err
}

func getAgentAddress(ipAddress string) string {
	addrPort, err := k8s_portforward.PortForwardNode(ipAddress, RestPort)
	if err != nil {
		// port fowarding failed, fall back to using  ipAddress:port
		addrPort = fmt.Sprintf("%s:%d", ipAddress, RestPort)
	}
	return addrPort
}

func sendRequestGetResponse(reqType, url string, data interface{}, verbose bool) (string, error) {
	client := &http.Client{}
	reqData := new(bytes.Buffer)
	if err := json.NewEncoder(reqData).Encode(data); err != nil {
		return "", err
	}
	req, err := http.NewRequest(reqType, url, reqData)
	if err != nil {
		fmt.Print(err.Error())
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("request returned code %d, %v, %s", resp.StatusCode, reqType, url)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if verbose {
		fmt.Printf("resp: %s\n", bodyBytes)
	}
	return string(bodyBytes), nil
}

// UngracefulReboot crashes and reboots the host machine
func UngracefulReboot(serverAddr string) error {
	logf.Log.Info("Executing ungraceReboot", "addr", serverAddr)
	url := "http://" + getAgentAddress(serverAddr) + "/ungracefulReboot"
	return sendRequest("POST", url, nil)
}

// IsAgentReachable checks if the agent pod is in reachable
func IsAgentReachable(serverAddr string) error {
	url := "http://" + getAgentAddress(serverAddr) + "/"
	return sendRequest("GET", url, nil)
}

// GracefulReboot reboots the host gracefully
// It is not yet supported
func GracefulReboot(serverAddr string) error {
	logf.Log.Info("Executing gracefulReboot", "addr", serverAddr)
	url := "http://" + getAgentAddress(serverAddr) + "/gracefulReboot"
	return sendRequest("POST", url, nil)
}

// DropConnectionsFromNodes creates rules to drop connections from other k8s nodes
func DropConnectionsFromNodes(serverAddr string, nodes []string) error {

	url := "http://" + getAgentAddress(serverAddr) + "/dropConnectionsFromNodes"
	data := NodeList{
		Nodes:            nodes,
		NetworkInterface: e2e_config.GetConfig().NetworkInterface,
	}
	logf.Log.Info("Executing dropConnectionsFromNodes", "addr", serverAddr, "data", data)
	return sendRequest("POST", url, data)
}

// AcceptConnectionsFromNodes removes the rules set by
// DropConnectionsFromNodes so that other k8s nodes can reach this node again
func AcceptConnectionsFromNodes(serverAddr string, nodes []string) error {
	url := "http://" + getAgentAddress(serverAddr) + "/acceptConnectionsFromNodes"
	data := NodeList{
		Nodes:            nodes,
		NetworkInterface: e2e_config.GetConfig().NetworkInterface,
	}
	logf.Log.Info("Executing acceptConnectionsFromNodes", "addr", serverAddr, "data", data)
	return sendRequest("POST", url, data)
}

// DiskPartition performs operation related to disk prtitioning
func DiskPartition(serverAddr string, cmd string) error {
	url := "http://" + getAgentAddress(serverAddr) + "/parted"
	data := CmdList{
		Cmd: cmd,
	}
	logf.Log.Info("Executing parted", "addr", serverAddr, "data", data)
	return sendRequest("POST", url, data)
}

// CreateFaultyDevice creates a device which returns an error on write IOs
func CreateFaultyDevice(serverAddr, device, table string) error {
	url := "http://" + getAgentAddress(serverAddr) + "/createFaultyDevice"
	data := Device{
		Device: device,
		Table:  table,
	}
	logf.Log.Info("Executing createFaultyDevice", "addr", serverAddr, "data", data)
	return sendRequest("POST", url, data)
}

// DeleteFaultyDevice deletes a device which returns an error on write IOs
func DeleteFaultyDevice(serverAddr, device string) error {
	url := "http://" + getAgentAddress(serverAddr) + "/deleteFaultyDevice"
	data := Device{
		Device: device,
	}
	logf.Log.Info("Executing deleteFaultyDevice", "addr", serverAddr, "data", data)
	return sendRequest("POST", url, data)
}

// ControlDevice sets the specified to the specified state
// by writing to /sys/block/<device e.g. sdb>/device/state
// The only accepted states are "running" and "offline"
func ControlDevice(serverAddr string, device string, state string) (string, error) {
	url := "http://" + getAgentAddress(serverAddr) + "/devicecontrol"
	data := ControlledDevice{
		Device: device,
		State:  state,
	}
	logf.Log.Info("Executing devicecontrol", "addr", serverAddr, "data", data)
	return sendRequestGetResponse("POST", url, data, false)
}

// KillIoEngine use kill -9 against the mayastor
func KillIoEngine(serverAddr string) (string, error) {
	url := "http://" + getAgentAddress(serverAddr) + "/killioengine"
	logf.Log.Info("Executing killioengine", "addr", serverAddr)
	return sendRequestGetResponse("POST", url, nil, true)
}

// KillCsiController use kill -9 against the mayastor csi controller process
func KillCsiController(serverAddr string) (string, error) {
	url := "http://" + getAgentAddress(serverAddr) + "/killCsiController"
	logf.Log.Info("Executing killCsiController", "addr", serverAddr)
	result, err := sendRequestGetResponse("POST", url, nil, true)
	if err != nil {
		return result, fmt.Errorf("failed to send command to e2e-agent, error: %s", err.Error())
	}
	out, errCode, err := UnwrapResult(result)
	if err == nil && errCode == 0 {
		return string(out), nil
	}
	return string(out), fmt.Errorf("errCode=%v ; err=%v", errCode, err)
}

// KillCsiNode use kill -9 against the mayastor csi node process
func KillCsiNode(serverAddr string) (string, error) {
	url := "http://" + getAgentAddress(serverAddr) + "/killCsiNode"
	logf.Log.Info("Executing killCsiNode", "addr", serverAddr)
	result, err := sendRequestGetResponse("POST", url, nil, true)
	if err != nil {
		return result, fmt.Errorf("failed to send command to e2e-agent, error: %s", err.Error())
	}
	out, errCode, err := UnwrapResult(result)
	if err == nil && errCode == 0 {
		return string(out), nil
	}
	return string(out), fmt.Errorf("errCode=%v ; err=%v", errCode, err)
}

// Flushcache flush the container engine cache
func FlushDiskWriteCache(serverAddr string) (string, error) {
	url := "http://" + getAgentAddress(serverAddr) + "/flushDiskWriteCache"
	logf.Log.Info("Executing flushcache", "addr", serverAddr)
	encodedresult, err := sendRequestGetResponse("POST", url, nil, true)
	if err != nil {
		return encodedresult, fmt.Errorf("failed to send command to e2e-agent, error: %s", err.Error())
	}
	out, e2eagenterrcode, err := UnwrapResult(encodedresult)
	if err != nil {
		return encodedresult, fmt.Errorf("failed to unwrap result, error: %s", err.Error())
	}
	if e2eagenterrcode != ErrNone {
		return out, fmt.Errorf("failed to flush container engine cache, errcode %d", e2eagenterrcode)
	}
	return out, err
}

// GetDeviceState sends the disk pool device status
func GetDeviceState(serverAddr string, disk string) (string, error) {
	url := "http://" + getAgentAddress(serverAddr) + "/getdevicestate"
	data := DiskPool{
		Disk: disk,
	}
	logf.Log.Info("Executing getdevicestate", "addr", serverAddr, "data", data)
	return sendRequestGetResponse("POST", url, data, false)
}

// NvmeConnect to connect to the target
func NvmeConnect(serverAddr string, targetIp string, nqn string, hostNqn string) (string, error) {
	data := Nvme{
		TargetIp: targetIp,
		Nqn:      nqn,
		HostNqn:  hostNqn,
		HostId:   "",
	}
	logf.Log.Info("Executing nvmeconnect", "addr", serverAddr, "data", data)
	url := "http://" + getAgentAddress(serverAddr) + "/nvmeconnect"
	return sendRequestGetResponse("POST", url, data, false)
}

// NvmeDisconnect to disconnect to the nvme target
func NvmeDisconnect(serverAddr string, nqn string) (string, error) {
	data := Nvme{
		Nqn: nqn,
	}
	logf.Log.Info("Executing nvmedisconnect", "addr", serverAddr, "data", data)
	url := "http://" + getAgentAddress(serverAddr) + "/nvmedisconnect"
	return sendRequestGetResponse("POST", url, data, false)
}

// NvmeList list nvme target
func NvmeList(serverAddr string) (string, error) {
	logf.Log.Info("Executing nvmelist", "addr", serverAddr)
	url := "http://" + getAgentAddress(serverAddr) + "/nvmelist"
	return sendRequestGetResponse("POST", url, nil, false)
}

// NvmeListSubSys list nvme target
func NvmeListSubSys(serverAddr string) (string, error) {
	logf.Log.Info("Executing nvmlistsubsys", "addr", serverAddr)
	url := "http://" + getAgentAddress(serverAddr) + "/nvmelistsubsys"
	return sendRequestGetResponse("POST", url, nil, false)
}

// checksum the device
// the returned format is <checksum> <size> <device>
func ChecksumDevice(serverAddr string, devicePath string) (string, error) {
	data := Device{
		DevicePath: devicePath,
	}
	logf.Log.Info("Executing checksumdevice", "addr", serverAddr, "data", data)
	url := "http://" + getAgentAddress(serverAddr) + "/checksumdevice"
	return sendRequestGetResponse("POST", url, data, false)
}

// fsCheck the device
func FsCheckDevice(serverAddr string, devicePath string, fsType common.FileSystemType) (string, error) {
	data := Device{
		DevicePath: devicePath,
		FsType:     string(fsType),
	}
	logf.Log.Info("Executing fscheckdevice", "addr", serverAddr, "data", data)
	url := "http://" + getAgentAddress(serverAddr) + "/fscheckdevice"
	result, err := sendRequestGetResponse("POST", url, data, false)
	if err != nil {
		return result, fmt.Errorf("failed to send command to e2e-agent, error: %s", err.Error())
	}
	out, errCode, err := UnwrapResult(result)
	if err == nil && errCode == 0 {
		return string(out), nil
	}
	return string(out), fmt.Errorf("errCode=%v ; err=%v", errCode, err)
}

// Performs fscheck equivalent for XFS filesystem
func XFSCheckDevice(serverAddr string, devicePath string, fsType common.FileSystemType) (string, error) {
	data := Device{
		DevicePath: devicePath,
		FsType:     string(fsType),
	}
	logf.Log.Info("Executing fscheckdevice", "addr", serverAddr, "data", data)
	url := "http://" + getAgentAddress(serverAddr) + "/xfscheckdevice"
	result, err := sendRequestGetResponse("POST", url, data, false)
	if err != nil {
		return result, fmt.Errorf("failed to send command to e2e-agent, error: %s", err.Error())
	}
	out, errCode, err := UnwrapResult(result)
	if err == nil && errCode == 0 {
		return string(out), nil
	}
	return string(out), fmt.Errorf("errCode=%v ; err=%v", errCode, err)
}

// fsFreeze the device
func FsFreezeDevice(serverAddr string, devicePath string, fsType common.FileSystemType) (string, error) {
	data := Device{
		DevicePath: devicePath,
	}
	logf.Log.Info("Executing fsfreezedevice", "addr", serverAddr, "data", data)
	url := "http://" + getAgentAddress(serverAddr) + "/fsfreezedevice"
	result, err := sendRequestGetResponse("POST", url, data, false)
	if err != nil {
		return result, fmt.Errorf("failed to send command to e2e-agent, error: %s", err.Error())
	}
	out, errCode, err := UnwrapResult(result)
	if err == nil && errCode == 0 {
		return string(out), nil
	}
	return string(out), fmt.Errorf("errCode=%v ; err=%v", errCode, err)
}

// fsUnfreeze the device
func FsUnfreezeDevice(serverAddr string, devicePath string) (string, error) {
	data := Device{
		DevicePath: devicePath,
	}
	logf.Log.Info("Executing fsUnfreezedevice", "addr", serverAddr, "data", data)
	url := "http://" + getAgentAddress(serverAddr) + "/fsunfreezedevice"
	result, err := sendRequestGetResponse("POST", url, data, false)
	if err != nil {
		return result, fmt.Errorf("failed to send command to e2e-agent, error: %s", err.Error())
	}
	out, errCode, err := UnwrapResult(result)
	if err == nil && errCode == 0 {
		return string(out), nil
	}
	return string(out), fmt.Errorf("errCode=%v ; err=%v", errCode, err)
}

// ListDevice list device
func ListDevice(serverAddr string) (string, error) {
	logf.Log.Info("Executing listdevice", "addr", serverAddr)
	url := "http://" + getAgentAddress(serverAddr) + "/listdevice"
	return sendRequestGetResponse("POST", url, nil, false)
}

// ZeroingDisk replace every block of ‘sda’ with zeroes.
func ZeroingDisk(serverAddr string, device string, seekParam string, blockSizeParam string) (string, error) {
	data := Disk{
		Device:         device,
		SeekParam:      seekParam,
		BlockSizeParam: blockSizeParam,
	}
	logf.Log.Info("Executing zeroingdisk", "addr", serverAddr, "data", data)
	url := "http://" + getAgentAddress(serverAddr) + "/zeroingdisk"
	return sendRequestGetResponse("POST", url, data, false)
}

// ListReservation list revervation report
func ListReservation(serverAddr string, devicePath string) (string, error) {
	data := Device{
		DevicePath: devicePath,
	}
	logf.Log.Info("Executing listreservation", "addr", serverAddr, "data", data)
	url := "http://" + getAgentAddress(serverAddr) + "/listreservation"
	return sendRequestGetResponse("POST", url, data, false)
}

// NvmeConnectWithHostId to connect to the target
func NvmeConnectWithHostId(serverAddr string, targetIp string, nqn string, hostNqn string, hostId string) (string, error) {
	data := Nvme{
		TargetIp: targetIp,
		Nqn:      nqn,
		HostNqn:  hostNqn,
		HostId:   hostId,
	}
	logf.Log.Info("Executing nvmeconnectwithhostid", "addr", serverAddr, "data", data)
	url := "http://" + getAgentAddress(serverAddr) + "/nvmeconnectwithhostid"
	return sendRequestGetResponse("POST", url, data, false)
}

func BlkDiscard(serverAddr string, device string, options string) (string, error) {
	data := blkDiscard{
		Device:  device,
		Options: options,
	}
	logf.Log.Info("Executing blkdiscard on node", "addr", serverAddr, "data", data)
	url := "http://" + getAgentAddress(serverAddr) + "/blkdiscard"
	res, err := sendRequestGetResponse("POST", url, data, false)
	if err != nil {
		return res, err
	}
	return "", err
}

type EventSubscriptionRequest struct {
	EventServerAddr string `json:"eventServerAddr"`
	Subject         string `json:"subject"`
}

type EventPublishRequest struct {
	EventServerAddr string `json:"eventServerAddr"`
	Subject         string `json:"subject"`
	Data            string `json:"data"`
}

type EventRequest struct {
	Subject string `json:"subject"`
}

func EventList(agentAddr string, subject string) (string, error) {
	data := EventRequest{
		Subject: subject,
	}
	logf.Log.Info("Executing EventList on node", "addr", agentAddr, "data", data)
	url := "http://" + getAgentAddress(agentAddr) + "/event/list"
	encodedresult, err := sendRequestGetResponse("POST", url, data, false)
	if err != nil {
		return encodedresult, fmt.Errorf("failed to send command to e2e-agent, error: %s", err.Error())
	}
	out, e2eagenterrcode, err := UnwrapResult(encodedresult)
	if err != nil {
		return encodedresult, fmt.Errorf("failed to unwrap result, error: %s", err.Error())
	}
	if e2eagenterrcode != ErrNone {
		return out, fmt.Errorf("failed to list events, errcode %d", e2eagenterrcode)
	}
	return out, err
}

func EventSubscribe(agentAddr string, eventServerAddr string, subject string) (string, error) {
	data := EventSubscriptionRequest{
		EventServerAddr: eventServerAddr,
		Subject:         subject,
	}
	logf.Log.Info("Executing EventSubscribe on node", "addr", agentAddr, "data", data)
	url := "http://" + getAgentAddress(agentAddr) + "/event/subscribe"
	encodedresult, err := sendRequestGetResponse("POST", url, data, false)
	if err != nil {
		return encodedresult, err
	}
	out, e2eagenterrcode, err := UnwrapResult(encodedresult)
	if err != nil {
		return encodedresult, err
	}
	if e2eagenterrcode != ErrNone {
		return out, fmt.Errorf("failed to subscribe, errcode %d", e2eagenterrcode)
	}
	return out, err
}

func EventPublish(agentAddr string, eventServerAddr string, subject string, message string) (string, error) {
	data := EventPublishRequest{
		EventServerAddr: eventServerAddr,
		Subject:         subject,
		Data:            message,
	}
	logf.Log.Info("Executing EventPublish on node", "addr", agentAddr, "data", data)
	url := "http://" + getAgentAddress(agentAddr) + "/event/publish"
	encodedresult, err := sendRequestGetResponse("POST", url, data, false)
	if err != nil {
		return encodedresult, err
	}
	out, e2eagenterrcode, err := UnwrapResult(encodedresult)
	if err != nil {
		return encodedresult, err
	}
	if e2eagenterrcode != ErrNone {
		return out, fmt.Errorf("failed to publish, errcode %d", e2eagenterrcode)
	}
	return out, err
}

func EventUnsubscribe(agentAddr string, subject string) (string, error) {
	data := EventRequest{
		Subject: subject,
	}
	logf.Log.Info("Executing EventUnsubscribe on node", "addr", agentAddr, "data", data)
	url := "http://" + getAgentAddress(agentAddr) + "/event/unsubscribe"
	encodedresult, err := sendRequestGetResponse("POST", url, data, false)
	if err != nil {
		return encodedresult, err
	}
	out, e2eagenterrcode, err := UnwrapResult(encodedresult)
	if err != nil {
		return encodedresult, err
	}
	if e2eagenterrcode != ErrNone {
		return out, fmt.Errorf("failed to unsubscribe, errcode %d", e2eagenterrcode)
	}
	return out, err
}

func EventUnsubscribeAll(agentAddr string) (string, error) {
	logf.Log.Info("Executing EventUnsubscribeAll on node", "addr", agentAddr)
	url := "http://" + getAgentAddress(agentAddr) + "/event/unsubscribeall"
	encodedresult, err := sendRequestGetResponse("POST", url, nil, false)
	if err != nil {
		return encodedresult, err
	}

	out, e2eagenterrcode, err := UnwrapResult(encodedresult)
	if err != nil {
		return encodedresult, err
	}
	if e2eagenterrcode != ErrNone {
		return out, fmt.Errorf("failed to unsubscribe all, errcode %d", e2eagenterrcode)
	}
	return out, err
}

type Stats struct {
	ServiceAddr string `json:"serviceAddr"`
}

func GetStats(agentAddr string, serviceAddr string) (string, error) {
	logf.Log.Info("Executing GetStats", "e2eagent node", agentAddr, "service", serviceAddr)
	data := Stats{
		ServiceAddr: serviceAddr,
	}

	url := "http://" + getAgentAddress(agentAddr) + "/stats"
	encodedresult, err := sendRequestGetResponse("POST", url, data, false)
	if err != nil {
		logf.Log.Info("sendRequestGetResponse", "encodedresult", encodedresult, "error", err.Error())
		return encodedresult, err
	}

	out, e2eagenterrcode, err := UnwrapResult(encodedresult)
	if err != nil {
		logf.Log.Info("unwrap failed", "encodedresult", encodedresult, "error", err.Error())
		return encodedresult, err
	}
	if e2eagenterrcode != ErrNone {
		return out, fmt.Errorf("failed to get stats, errcode %d", e2eagenterrcode)
	}
	logf.Log.Info("GetStats succeeded", "output", out)
	return out, err
}

// ZeroHugePages sets Huge Pages value to 0 on the selected node
func ZeroHugePages(serverAddr string) (string, error) {
	url := "http://" + getAgentAddress(serverAddr) + "/hugepagezero"
	logf.Log.Info("Executing hugepagezero", "addr", serverAddr)
	return sendRequestGetResponse("POST", url, nil, false)
}

type E2eAgentErrcode int

const (
	// general errors
	ErrNone       E2eAgentErrcode = 0
	ErrGeneral    E2eAgentErrcode = 1
	ErrJsonDecode E2eAgentErrcode = 2
	ErrJsonEncode E2eAgentErrcode = 3
	ErrReadFail   E2eAgentErrcode = 4
	ErrExecFailed E2eAgentErrcode = 5

	// event errors
	ErrConnectFail          E2eAgentErrcode = 101
	ErrConnectedOther       E2eAgentErrcode = 102
	ErrNotConnected         E2eAgentErrcode = 103
	ErrSubscriptionNotFound E2eAgentErrcode = 104
	ErrSubscribedAlready    E2eAgentErrcode = 105
	ErrSubscribeFail        E2eAgentErrcode = 106
	ErrStreamAddFail        E2eAgentErrcode = 107
	ErrStreamCreateFail     E2eAgentErrcode = 108
	ErrPublishFail          E2eAgentErrcode = 109
)

type E2eAgentError struct {
	Output    string          `json:"output"`
	Errorcode E2eAgentErrcode `json:"errorcode"`
}

func UnwrapResult(res string) (string, E2eAgentErrcode, error) {
	var response E2eAgentError
	var err error
	if err = json.Unmarshal([]byte(res), &response); err != nil {
		return "", -1, err
	}
	return response.Output, response.Errorcode, err
}

// cmp 2 devices
func Cmp(serverAddr string, path1 string, path2 string) (string, error) {
	data := CmpPaths{
		Path1: path1,
		Path2: path2,
	}
	logf.Log.Info("Executing cmp", "addr", serverAddr, "data", data)
	url := "http://" + getAgentAddress(serverAddr) + "/cmp"
	result, err := sendRequestGetResponse("POST", url, data, false)
	if err != nil {
		return result, fmt.Errorf("failed to send command to e2e-agent, error: %s", err.Error())
	}
	b64out, errCode, err := UnwrapResult(result)
	out, _ := base64.StdEncoding.DecodeString(b64out)
	if err == nil && errCode == 0 {
		return string(out), nil
	}
	return string(out), fmt.Errorf("errCode=%v ; err=%v", errCode, err)
}

// LvmListVg list lvm vg
func LvmListVg(serverAddr string) (string, error) {
	logf.Log.Info("Executing LvmListVg", "addr", serverAddr)
	url := "http://" + getAgentAddress(serverAddr) + "/lvmlistvg"
	encodedresult, err := sendRequestGetResponse("POST", url, nil, false)
	if err != nil {
		logf.Log.Info("sendRequestGetResponse", "encodedresult", encodedresult, "error", err.Error())
		return encodedresult, err
	}

	out, e2eagenterrcode, err := UnwrapResult(encodedresult)
	if err != nil {
		logf.Log.Info("unwrap failed", "encodedresult", encodedresult, "error", err.Error())
		return encodedresult, err
	}
	if e2eagenterrcode != ErrNone {
		return out, fmt.Errorf("failed to get lvm vg, errcode %d", e2eagenterrcode)
	}
	logf.Log.Info("LvmListVg succeeded", "output", out)
	return out, err
}

// LvmListPv list lvm pv
func LvmListPv(serverAddr string) (string, error) {
	logf.Log.Info("Executing LvmListPv", "addr", serverAddr)
	url := "http://" + getAgentAddress(serverAddr) + "/lvmlistpv"
	encodedresult, err := sendRequestGetResponse("POST", url, nil, false)
	if err != nil {
		logf.Log.Info("sendRequestGetResponse", "encodedresult", encodedresult, "error", err.Error())
		return encodedresult, err
	}

	out, e2eagenterrcode, err := UnwrapResult(encodedresult)
	if err != nil {
		logf.Log.Info("unwrap failed", "encodedresult", encodedresult, "error", err.Error())
		return encodedresult, err
	}
	if e2eagenterrcode != ErrNone {
		return out, fmt.Errorf("failed to get lvm pv, errcode %d", e2eagenterrcode)
	}
	logf.Log.Info("LvmListPv succeeded", "output", out)
	return out, err
}

// LvmVersion get lvm version
func LvmVersion(serverAddr string) (string, error) {
	logf.Log.Info("Executing LvmVersion", "addr", serverAddr)
	url := "http://" + getAgentAddress(serverAddr) + "/lvmversion"
	encodedresult, err := sendRequestGetResponse("POST", url, nil, false)
	if err != nil {
		logf.Log.Info("sendRequestGetResponse", "encodedresult", encodedresult, "error", err.Error())
		return encodedresult, err
	}

	out, e2eagenterrcode, err := UnwrapResult(encodedresult)
	if err != nil {
		logf.Log.Info("unwrap failed", "encodedresult", encodedresult, "error", err.Error())
		return encodedresult, err
	}
	if e2eagenterrcode != ErrNone {
		return out, fmt.Errorf("failed to get lvm version, errcode %d", e2eagenterrcode)
	}
	logf.Log.Info("LvmVersion succeeded", "output", out)
	return out, err
}

// LvmCreatePv create lvm pv
func LvmCreatePv(serverAddr string, pvDiskPath string) (string, error) {
	data := Lvm{
		Pv: pvDiskPath,
	}
	logf.Log.Info("Executing lvmcreatepv", "addr", serverAddr, "data", data)
	url := "http://" + getAgentAddress(serverAddr) + "/lvmcreatepv"
	encodedresult, err := sendRequestGetResponse("POST", url, data, false)
	if err != nil {
		logf.Log.Info("sendRequestGetResponse", "encodedresult", encodedresult, "error", err.Error())
		return encodedresult, err
	}
	out, e2eagenterrcode, err := UnwrapResult(encodedresult)
	if err != nil {
		logf.Log.Info("unwrap failed", "encodedresult", encodedresult, "error", err.Error())
		return encodedresult, err
	}
	if e2eagenterrcode != ErrNone {
		return out, fmt.Errorf("failed to create lvm pv, errcode %d", e2eagenterrcode)
	}
	logf.Log.Info("LvmCreatePv succeeded", "output", out)
	return out, err
}

// LvmCreateVg create lvm vg
func LvmCreateVg(serverAddr string, pvDiskPath string, vgName string) (string, error) {
	data := Lvm{
		Pv: pvDiskPath,
		Vg: vgName,
	}
	logf.Log.Info("Executing lvmcreatevg", "addr", serverAddr, "data", data)
	url := "http://" + getAgentAddress(serverAddr) + "/lvmcreatevg"
	encodedresult, err := sendRequestGetResponse("POST", url, data, false)
	if err != nil {
		logf.Log.Info("sendRequestGetResponse", "encodedresult", encodedresult, "error", err.Error())
		return encodedresult, err
	}
	out, e2eagenterrcode, err := UnwrapResult(encodedresult)
	if err != nil {
		logf.Log.Info("unwrap failed", "encodedresult", encodedresult, "error", err.Error())
		return encodedresult, err
	}
	if e2eagenterrcode != ErrNone {
		return out, fmt.Errorf("failed to create lvm pv, errcode %d", e2eagenterrcode)
	}
	logf.Log.Info("LvmCreateVg succeeded", "output", out)
	return out, err
}

// LvmRemovePv remove lvm pv
func LvmRemovePv(serverAddr string, pvDiskPath string, vgName string) (string, error) {
	data := Lvm{
		Pv: pvDiskPath,
	}
	logf.Log.Info("Executing lvmremovepv", "addr", serverAddr, "data", data)
	url := "http://" + getAgentAddress(serverAddr) + "/lvmremovepv"
	encodedresult, err := sendRequestGetResponse("POST", url, data, false)
	if err != nil {
		logf.Log.Info("sendRequestGetResponse", "encodedresult", encodedresult, "error", err.Error())
		return encodedresult, err
	}
	out, e2eagenterrcode, err := UnwrapResult(encodedresult)
	if err != nil {
		logf.Log.Info("unwrap failed", "encodedresult", encodedresult, "error", err.Error())
		return encodedresult, err
	}
	if e2eagenterrcode != ErrNone {
		return out, fmt.Errorf("failed to remove lvm pv, errcode %d", e2eagenterrcode)
	}
	logf.Log.Info("LvmRemovePv succeeded", "output", out)
	return out, err
}

// LvmRemoveVg remove lvm vg
func LvmRemoveVg(serverAddr string, pvDiskPath string, vgName string) (string, error) {
	data := Lvm{
		Vg: vgName,
	}
	logf.Log.Info("Executing lvmremovevg", "addr", serverAddr, "data", data)
	url := "http://" + getAgentAddress(serverAddr) + "/lvmremovevg"
	encodedresult, err := sendRequestGetResponse("POST", url, data, false)
	if err != nil {
		logf.Log.Info("sendRequestGetResponse", "encodedresult", encodedresult, "error", err.Error())
		return encodedresult, err
	}
	out, e2eagenterrcode, err := UnwrapResult(encodedresult)
	if err != nil {
		logf.Log.Info("unwrap failed", "encodedresult", encodedresult, "error", err.Error())
		return encodedresult, err
	}
	if e2eagenterrcode != ErrNone {
		return out, fmt.Errorf("failed to remove lvm vg, errcode %d", e2eagenterrcode)
	}
	logf.Log.Info("LvmRemoveVg succeeded", "output", out)
	return out, err
}

// LvmThinPoolAutoExtendThreshold update lvm thin pool auto extend threshold value
func LvmThinPoolAutoExtendThreshold(serverAddr string, thinPoolAutoExtendThreshold int) (string, error) {
	data := Lvm{
		ThinPoolAutoExtendThreshold: thinPoolAutoExtendThreshold,
	}
	logf.Log.Info("Executing lvmthinpoolautoextendthreshold", "addr", serverAddr, "data", data)
	url := "http://" + getAgentAddress(serverAddr) + "/lvmthinpoolautoextendthreshold"
	encodedresult, err := sendRequestGetResponse("POST", url, data, false)
	if err != nil {
		logf.Log.Info("sendRequestGetResponse", "encodedresult", encodedresult, "error", err.Error())
		return encodedresult, err
	}
	out, e2eagenterrcode, err := UnwrapResult(encodedresult)
	if err != nil {
		logf.Log.Info("unwrap failed", "encodedresult", encodedresult, "error", err.Error())
		return encodedresult, err
	}
	if e2eagenterrcode != ErrNone {
		return out, fmt.Errorf("failed to update lvm thin pool auto extend threshold, errcode %d", e2eagenterrcode)
	}
	logf.Log.Info("LvmThinPoolAutoExtendThreshold succeeded", "output", out)
	return out, err
}

// LvmThinPoolAutoExtendPercent update lvm thin pool auto extend percent value
func LvmThinPoolAutoExtendPercent(serverAddr string, thinPoolAutoExtendPercent int) (string, error) {
	data := Lvm{
		ThinPoolAutoExtendPercent: thinPoolAutoExtendPercent,
	}
	logf.Log.Info("Executing lvmthinpoolautoextendpercent", "addr", serverAddr, "data", data)
	url := "http://" + getAgentAddress(serverAddr) + "/lvmthinpoolautoextendpercent"
	encodedresult, err := sendRequestGetResponse("POST", url, data, false)
	if err != nil {
		logf.Log.Info("sendRequestGetResponse", "encodedresult", encodedresult, "error", err.Error())
		return encodedresult, err
	}
	out, e2eagenterrcode, err := UnwrapResult(encodedresult)
	if err != nil {
		logf.Log.Info("unwrap failed", "encodedresult", encodedresult, "error", err.Error())
		return encodedresult, err
	}
	if e2eagenterrcode != ErrNone {
		return out, fmt.Errorf("failed to update lvm thin pool auto extend percent, errcode %d", e2eagenterrcode)
	}
	logf.Log.Info("LvmThinPoolAutoExtendPercent succeeded", "output", out)
	return out, err
}

// CreateLoopDevice create loop device
func CreateLoopDevice(serverAddr string, size int64, imageDir string) (string, error) {
	data := LoopDevice{
		ImgDir: imageDir,
		Size:   size,
	}
	logf.Log.Info("Executing createloopdevice", "addr", serverAddr, "data", data)
	url := "http://" + getAgentAddress(serverAddr) + "/createloopdevice"
	encodedresult, err := sendRequestGetResponse("POST", url, data, false)
	if err != nil {
		logf.Log.Info("sendRequestGetResponse", "encodedresult", encodedresult, "error", err.Error())
		return encodedresult, err
	}
	out, e2eagenterrcode, err := UnwrapResult(encodedresult)
	if err != nil {
		logf.Log.Info("unwrap failed", "encodedresult", encodedresult, "error", err.Error())
		return encodedresult, err
	}
	if e2eagenterrcode != ErrNone {
		return out, fmt.Errorf("failed to create loop device, errcode %d", e2eagenterrcode)
	}
	logf.Log.Info("CreateLoopDevice succeeded", "output", out)
	return out, err
}

// DeleteLoopDevice delete loop device
func DeleteLoopDevice(serverAddr string, dsikPath string, imageName string) (string, error) {
	data := LoopDevice{
		DiskPath:  dsikPath,
		ImageName: imageName,
	}
	logf.Log.Info("Executing deleteloopdevice", "addr", serverAddr, "data", data)
	url := "http://" + getAgentAddress(serverAddr) + "/deleteloopdevice"
	encodedresult, err := sendRequestGetResponse("POST", url, data, false)
	if err != nil {
		logf.Log.Info("sendRequestGetResponse", "encodedresult", encodedresult, "error", err.Error())
		return encodedresult, err
	}
	out, e2eagenterrcode, err := UnwrapResult(encodedresult)
	if err != nil {
		logf.Log.Info("unwrap failed", "encodedresult", encodedresult, "error", err.Error())
		return encodedresult, err
	}
	if e2eagenterrcode != ErrNone {
		return out, fmt.Errorf("failed to delete loop device, errcode %d", e2eagenterrcode)
	}
	logf.Log.Info("DeleteLoopDevice succeeded", "output", out)
	return out, err
}
