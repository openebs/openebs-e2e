package k8s_portforward

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"strings"
	"sync"
	"syscall"
	"time"

	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8slabels "k8s.io/apimachinery/pkg/labels"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
	k8sConfig "sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

const PROXY_KEY = "e2e-proxy"
const NILSTRING = ""

var once sync.Once
var wg sync.WaitGroup
var e2eProxyPod *coreV1.Pod

// map for port forwarding
var pfMap = map[string]*PortForwardPodStatus{}

// port forwarding map write/delete lock
var pfMapLock sync.Mutex

var restConfig *rest.Config
var kubeInt *kubernetes.Clientset

func initialise() {
	once.Do(func() {
		log.SetLogger(zap.New(zap.UseDevMode(true), zap.WriteTo(os.Stdout)))
		restConfig = k8sConfig.GetConfigOrDie()
		kubeInt = kubernetes.NewForConfigOrDie(restConfig)
		log.Log.Info("port-forward initialise", "restConfig", restConfig != nil, "kubeInt", kubeInt != nil)
		if restConfig == nil {
			log.Log.Info("ERROR: failed to create *rest.Config for talking to a Kubernetes apiserver : GetConfigOrDie")
		}
		if kubeInt == nil {
			log.Log.Info("ERROR: failed to create new Clientset for the given config : NewForConfigOrDie")
		}
	})
}

type PortForwardPodStatus struct {
	// Streams configures where to write or read input from
	Streams genericiooptions.IOStreams
	// StopCh is the channel used to manage the port forward lifecycle
	StopCh chan struct{}
	// ReadyCh communicates when the tunnel is ready to receive traffic
	ReadyCh chan struct{}
	//
	pf         *portforward.PortForwarder
	RemotePort uint16
	LocalPort  uint16
	Done       bool
	Pod        coreV1.Pod
}

func (pf PortForwardPodStatus) String() string {
	return fmt.Sprintf("RemotePort:%v, LocalPort:%v, Done:%v ",
		pf.RemotePort, pf.LocalPort, pf.Done)
}

func (pf PortForwardPodStatus) PortAddress() string {
	return fmt.Sprintf(":%d", pf.LocalPort)
}

// This function MUST be called with pfMapLock locked.
func portForwardToPod(pod coreV1.Pod, localPort int, remotePort int, wg *sync.WaitGroup) (*PortForwardPodStatus, error) {
	//	log.Log.Info("portForwardToPod", "pod", pod.Name, "namespace", pod.Namespace, "localPort", localPort, "remotePort", remotePort)
	resourcePath := path.Join("api", "v1", "namespaces", pod.Namespace, "pods", pod.Name, "portforward")
	pfStatus := PortForwardPodStatus{}
	pfStatus.StopCh = make(chan struct{}, 1)
	pfStatus.ReadyCh = make(chan struct{})
	pfStatus.Streams = genericiooptions.IOStreams{In: os.Stdin, Out: io.Discard, ErrOut: io.Discard}
	pfStatus.Pod = pod

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func(pfr *PortForwardPodStatus) {
		<-sigs
		close(pfr.StopCh)
	}(&pfStatus)

	targetURL, err := url.Parse(restConfig.Host)
	if err != nil {
		return nil, err
	}
	targetURL.Path = resourcePath
	transport, upgrader, err := spdy.RoundTripperFor(restConfig)
	if err != nil {
		return nil, err
	}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, targetURL)
	//	fw, err := portforward.NewOnAddresses(dialer, []string{ipAddress}, []string{fmt.Sprintf("%d:%d", localPort, remotePort)}, pfStatus.StopCh, pfStatus.ReadyCh, pfStatus.Streams.Out, pfStatus.Streams.ErrOut)
	fw, err := portforward.New(dialer, []string{fmt.Sprintf("%d:%d", localPort, remotePort)}, pfStatus.StopCh, pfStatus.ReadyCh, pfStatus.Streams.Out, pfStatus.Streams.ErrOut)
	if err != nil {
		return nil, err
	}
	wg.Add(1)
	go func(pfstat *PortForwardPodStatus) {
		_ = fw.ForwardPorts()
		wg.Done()
		pfstat.Done = true
	}(&pfStatus)

	//	log.Log.Info("waiting ready.....")
	select {
	case <-pfStatus.ReadyCh:
		//		log.Log.Info("ready!")
		break
	case <-time.After(time.Second * 30):
		log.Log.Info("timeout waiting for port-forwarding", resourcePath, remotePort)
		close(pfStatus.StopCh)
	}

	pfStatus.pf = fw
	return &pfStatus, err
}

func (pf *PortForwardPodStatus) isAlive() bool {
	if pf.Done {
		return false
	}
	if pod, err := kubeInt.CoreV1().Pods(pf.Pod.Namespace).Get(context.TODO(), pf.Pod.Name, metaV1.GetOptions{}); err == nil {
		return pod.Status.Phase == coreV1.PodRunning
	}
	return false
}

// portForwardToPodByNames select a pod to portfoward uses names(prefix to name and namepace) as the filter
// This function MUST be called with pfMapLock locked.
func portForwardToPodByNames(key string, ipAddress string, port int, namespace string, podPrefix string) (string, *coreV1.Pod, error) {
	log.Log.Info("portForwardToPodByNames", "ipAddress", ipAddress, "port", port, "namespace", namespace, "podPrefix", podPrefix)
	var targetPod *coreV1.Pod
	// Between the invoking function checking and locking a.n.other thread could have established the required port forwarding
	// retry and delete stale entry
	pf, ok := pfMap[key]
	if ok {
		if pf.isAlive() {
			return pf.PortAddress(), targetPod, nil
		}
		// was alive but is dead now - remove stale connection
		close(pf.StopCh)
		delete(pfMap, key)
	}

	podList, err := kubeInt.CoreV1().Pods(namespace).List(context.TODO(), metaV1.ListOptions{})
	if err == nil {
		for _, pod := range podList.Items {
			if strings.HasPrefix(pod.Name, podPrefix) && (ipAddress == "" || ipAddress == pod.Status.HostIP) {
				targetPod = &pod
				log.Log.Info("SELECTED ", "name", targetPod.Name, "ipAddress", ipAddress)
				// if ipAddress was unspecified use the pod status host IP address
				ipAddress = pod.Status.HostIP
				break
			}
		}
	}
	if targetPod == nil {
		return NILSTRING, targetPod, fmt.Errorf("pod not found %s namespace:%s (ipAddress:%s)", podPrefix, namespace, ipAddress)
	}

	pf, err = portForwardToPod(coreV1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      targetPod.Name,
			Namespace: namespace,
		},
	}, 0, port, &wg)

	if err != nil {
		return NILSTRING, targetPod, err
	}
	ports, err := pf.pf.GetPorts()
	if err != nil {
		return NILSTRING, targetPod, err
	}
	pf.LocalPort = ports[0].Local
	pf.RemotePort = ports[0].Remote
	log.Log.Info("forwarding", "pod", targetPod.Name, "ports", ports)
	pfMap[key] = pf
	return pf.PortAddress(), targetPod, err
}

// PortForwardNode  establish port forwarding to an IP address - port combination,
// and return the local port.
// This establishes a 2 hop port forwarding to the target,
//
//	1 - port forward to the in cluster e2e-proxy pod
//	2 - e2e-proxy pod forwards the IP address port.
//
// Note port forwarding is only instantiated if one does not exist or the previous
// port forwarding has stopped.
func PortForwardNode(address string, port int) (string, error) {
	//	log.Log.Info("PortForwardOnNode", "address", address, "port", port)
	initialise()
	var err error
	var proxyAddress string

	pfMapLock.Lock()
	defer pfMapLock.Unlock()
	key := fmt.Sprintf("%s:%d", address, port)

	// 2 hop port-forwarding, check that forwarding to the e2e-proxy pod is alive
	proxyPf, ok := pfMap[PROXY_KEY]
	if ok && proxyPf.isAlive() {
		// 1st hop is up
		if pf, ok := pfMap[key]; ok {
			// 2nd hop exists
			return pf.PortAddress(), nil
		}
		proxyAddress = proxyPf.PortAddress()
	} else {
		// set up port forwarding to the e2e-proxy pod on the cluster if required
		if ok {
			log.Log.Info("e2e-proxy pod is down, clearing existing port-forwarding")
			// port forwarding to the proxy is down
			// delete close all forwarding to the proxy
			var fwKeys []string
			// first close existing connections
			for k, v := range pfMap {
				if v.Pod.Name == proxyPf.Pod.Name {
					close(v.StopCh)
				}
				fwKeys = append(fwKeys, k)
			}
			// then remove the entries from the map
			for _, k := range fwKeys {
				//				log.Log.Info("pfMap" "del", k)
				delete(pfMap, k)
			}
		}
		// we have a clean slate at this point
		// establish port forwarding to the e2e-proxy pod
		proxyAddress, e2eProxyPod, err = portForwardToPodByNames(PROXY_KEY, "", 8080, "e2e-agent", "e2e-proxy")
		if err != nil {
			return NILSTRING, err
		}
	}

	{
		forwarding := []map[string]int{{address: port}}
		var jsonData []byte
		jsonData, err = json.Marshal(forwarding)
		if err != nil {
			log.Log.Info("Marshalling failed", "error", err)
			return NILSTRING, err
		}
		reqBody := bytes.NewBuffer(jsonData)
		_, err = http.Post(fmt.Sprintf("http://%s/forward", proxyAddress), "application/json", reqBody)
		if err != nil {
			log.Log.Info("http post forwarding request to e2e proxy failed", "error", err)
			return NILSTRING, err
		}
	}

	// refresh port forwarding map, with the list of port forwards returned by the e2e-proxy
	var resp *http.Response
	resp, err = http.Get(fmt.Sprintf("http://%s/listforwarding", proxyAddress))
	if err != nil {
		// request to e2e-proxy failed
		log.Log.Info("request to list port forwarding on e2e proxy failed", "proxyAddress", proxyAddress, "error", err)
		return NILSTRING, err
	}
	var body []byte
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Log.Info("failed to get a valid http response when listing port forwarding on e2e proxy", "proxyAddress", proxyAddress, "error", err)
		return NILSTRING, err
	}

	m := make(map[string]int)
	err = json.Unmarshal(body, &m)
	if err != nil {
		log.Log.Info("json unmarshal failed on http response", "error", err)
		return NILSTRING, err
	}

	localPort := 0
	for target, remotePort := range m {
		var ok bool
		var pf *PortForwardPodStatus
		pf, ok = pfMap[target]
		if ok {
			if pf.Done {
				delete(pfMap, target)
			} else {
				continue
			}
		}
		pf, err = portForwardToPod(*e2eProxyPod, localPort, remotePort, &wg)
		if err == nil {
			ports, err := pf.pf.GetPorts()
			pfMap[target] = pf
			if err == nil {
				pf.LocalPort = ports[0].Local
				pf.RemotePort = ports[0].Remote
				log.Log.Info("populate: forwarding", "target", target, "ports", ports)
			} else {
				log.Log.Info("forwarding failed", "target", target, "err", err)
			}
		} else {
			log.Log.Info("forwarding failed", "target", target, "err", err)
		}
	}

	pf := pfMap[key]
	if pf == nil {
		return NILSTRING, fmt.Errorf("failed to find port forwarding entry for %s", key)
	}
	return pf.PortAddress(), nil
}

// PortForwardService establish port forwarding to a service running on the cluster
// and return the local port.
// port forwarding is only instantiated if one does not exist or the previous
// port forwarding has stopped.
func PortForwardService(svcName string, namespace string, port int) (string, error) {
	//	log.Log.Info("PortForwardService")
	var err error
	var ports []portforward.ForwardedPort
	initialise()

	key := fmt.Sprintf("service-%s-%s", svcName, namespace)
	pf, ok := pfMap[key]
	if ok {
		if pf.isAlive() {
			return pf.PortAddress(), nil
		}
	}
	pfMapLock.Lock()
	defer pfMapLock.Unlock()
	pf, ok = pfMap[key]
	if ok {
		if pf.isAlive() {
			return pf.PortAddress(), nil
		}
		// pod is dead - close and remove stale entry
		close(pf.StopCh)
		delete(pfMap, key)
	}

	var svc *coreV1.Service
	svc, err = kubeInt.CoreV1().Services(namespace).Get(context.TODO(), svcName, metaV1.GetOptions{})
	if err != nil || svc == nil {
		return NILSTRING, fmt.Errorf("failed to retrieve service %s; %v", svcName, err)
	}
	var podList *coreV1.PodList
	podList, err = kubeInt.CoreV1().Pods(namespace).List(context.TODO(),
		metaV1.ListOptions{
			LabelSelector: k8slabels.SelectorFromSet(svc.Spec.Selector).String(),
		})
	if err != nil {
		return NILSTRING, err
	}
	for _, pod := range podList.Items {
		var pf *PortForwardPodStatus
		pf, err = portForwardToPod(pod, 0, port, &wg)
		if err != nil {
			//					return NILSTRING, err
			continue
		}
		ports, err = pf.pf.GetPorts()
		if err != nil {
			//					return NILSTRING, err
			continue
		}
		pf.LocalPort = ports[0].Local
		pf.RemotePort = ports[0].Remote
		log.Log.Info("forwarding", "service", svcName, "ports", ports)
		pfMap[key] = pf
		return pf.PortAddress(), err
	}
	return NILSTRING, err
}

func TryPortForwardNode(address string, port int) string {
	addrPort, err := PortForwardNode(address, port)
	if err != nil {
		log.Log.Info(fmt.Sprintf("TryPortForwardNode: falling back to %s:%v", address, port))
		return fmt.Sprintf("%s:%d", address, port)
	}
	return addrPort
}
