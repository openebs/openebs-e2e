package v1beta1

import (
	"github.com/openebs/openebs-e2e/common/e2e_config"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	PoolSchemeBuilder = runtime.NewSchemeBuilder(poolAddKnownTypes)
	PoolAddToScheme   = PoolSchemeBuilder.AddToScheme
)

func poolAddKnownTypes(scheme *runtime.Scheme) error {
	var PoolSchemeGroupVersion = schema.GroupVersion{Group: e2e_config.GetConfig().Product.CrdGroupName,
		Version: "v1beta1"}

	scheme.AddKnownTypes(PoolSchemeGroupVersion,
		&DiskPool{},
		&DiskPoolList{},
	)

	metaV1.AddToGroupVersion(scheme, PoolSchemeGroupVersion)
	return nil
}
