package k8stest

import (
	"context"

	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
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

func GetService(serviceName string, nameSpace string) (*v1.Service, error) {
	service, err := gTestEnv.KubeInt.CoreV1().Services(nameSpace).Get(context.TODO(), serviceName, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}
	logf.Log.Info("Found service", "name", serviceName, "IP", service.Spec.ClusterIP)
	return service, nil
}

func IsServiceExists(serviceName string, nameSpace string) (bool, error) {
	_, err := gTestEnv.KubeInt.CoreV1().Services(nameSpace).Get(context.TODO(), serviceName, metaV1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	logf.Log.Info("Found service", "name", serviceName)
	return true, nil
}
