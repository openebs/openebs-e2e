package k8stest

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetConfigMap(name string, nameSpace string) (*v1.ConfigMap, error) {
	return gTestEnv.KubeInt.CoreV1().ConfigMaps(nameSpace).Get(context.TODO(), name, metaV1.GetOptions{})
}

func DeleteConfigMap(name string, nameSpace string) error {
	return gTestEnv.KubeInt.CoreV1().ConfigMaps(nameSpace).Delete(context.TODO(), name, metaV1.DeleteOptions{})
}

func ReplaceConfigMapData(name string, nameSpace string, data map[string]string) error {
	cm, err := GetConfigMap(name, nameSpace)
	if err != nil {
		return err
	}
	cm.Data = data
	_, err = gTestEnv.KubeInt.CoreV1().ConfigMaps(nameSpace).Update(context.TODO(), cm, metaV1.UpdateOptions{})
	return err
}
