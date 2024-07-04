package stats

import (
	"fmt"
	"strings"

	"github.com/openebs/openebs-e2e/common/e2e_agent"
	"github.com/openebs/openebs-e2e/common/k8stest"

	"github.com/prometheus/common/expfmt"
)

type StatsContext struct {
	E2eAgentAddress string
	StatsServiceIp  string
}

type StatsType int
type StatsAction int

const (
	POOL   StatsType = 1
	VOLUME StatsType = 2
)

const (
	CREATED StatsAction = 1
	DELETED StatsAction = 2
)

func NewStatsContext(statsService string, port string, namespace string) (StatsContext, error) {
	var err error
	var context StatsContext
	nodeIPs := k8stest.GetMayastorNodeIPAddresses()

	if len(nodeIPs) < 1 {
		return context, fmt.Errorf("no mayastor nodes")
	}
	context.E2eAgentAddress = nodeIPs[0]

	statsServiceIp, err := k8stest.GetServiceIp(statsService, namespace)
	if err != nil {
		return context, fmt.Errorf("failed to get service %s, namespace %s, error %s", statsService, namespace, err.Error())
	}
	statsServiceIp += ":" + port
	context.StatsServiceIp = statsServiceIp
	return context, err
}

func (context *StatsContext) GetStats(statsType StatsType, statsAction StatsAction) (int, error) {
	out, err := e2e_agent.GetStats(context.E2eAgentAddress, context.StatsServiceIp)
	if err != nil {
		return 0, fmt.Errorf("failed to get stats type: %s, action: %s, error %s", fmt.Sprint(statsType), fmt.Sprint(statsAction), err.Error())
	}
	return context.ParseStats(out, statsType, statsAction)
}

func (context *StatsContext) ParseStats(raw string, statsType StatsType, statsAction StatsAction) (int, error) {
	r := strings.NewReader(raw)

	var parser expfmt.TextParser
	mf, err := parser.TextToMetricFamilies(r)
	if err != nil {
		fmt.Printf("er, error: %s", err.Error())
		return 0, fmt.Errorf("failed to parse stats text, error: %s", err.Error())
	}
	var typekey string
	var actionkey string
	switch statsType {
	case POOL:
		typekey = "pool"
	case VOLUME:
		typekey = "volume"
	}
	switch statsAction {
	case CREATED:
		actionkey = "created"
	case DELETED:
		actionkey = "deleted"
	}

	item, exists := mf[typekey]
	if !exists {
		return 0, fmt.Errorf("Failed to find typekey %s", typekey)
	}

	for _, m := range item.Metric {
		for _, l := range m.Label {
			if *l.Name == "action" && *l.Value == actionkey {
				if m.Counter != nil {
					counter := m.Counter
					return int(*counter.Value), nil
				} else {
					fmt.Printf("nil counter value in %v\n", m)
				}
			}
		}
	}
	return 0, fmt.Errorf("could not find metric item: %s action: %s", typekey, actionkey)
}
