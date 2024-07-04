package v1alpha1Ext

import (
	"github.com/openebs/openebs-e2e/common/e2e_config"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

// Pool API

type DiskPoolV1Alpha1Interface interface {
	DiskPools() DiskPoolInterface
}

type DiskPoolV1Alpha1Client struct {
	restClient rest.Interface
}

func DspNewForConfig(c *rest.Config) (*DiskPoolV1Alpha1Client, error) {
	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: e2e_config.GetConfig().Product.CrdGroupName,
		Version: "v1alpha1"}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &DiskPoolV1Alpha1Client{restClient: client}, nil
}

func (c *DiskPoolV1Alpha1Client) DiskPools() DiskPoolInterface {
	return &dspClient{
		restClient: c.restClient,
	}
}
