package k8stest

import (
	"context"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

func GetServiceIp(serviceName string, nameSpace string) (string, error) {
	service, err := gTestEnv.KubeInt.CoreV1().Services(nameSpace).Get(context.TODO(), serviceName, metaV1.GetOptions{})
	if err != nil {
		return "", err
	}
	logf.Log.Info("Found service", "name", serviceName, "IP", service.Spec.ClusterIP)
	return service.Spec.ClusterIP, err
}
