package v1beta3

import (
	"github.com/openebs/openebs-e2e/common/e2e_config"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

// Pool APIx

type DiskPoolV1Beta3Interface interface {
	DiskPools() DiskPoolInterface
}

type DiskPoolV1Beta3Client struct {
	restClient rest.Interface
}

func DspNewForConfig(c *rest.Config) (*DiskPoolV1Beta3Client, error) {
	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: e2e_config.GetConfig().Product.CrdGroupName,
		Version: "v1beta3"}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &DiskPoolV1Beta3Client{restClient: client}, nil
}

func (c *DiskPoolV1Beta3Client) DiskPools() DiskPoolInterface {
	return &dspClient{
		restClient: c.restClient,
	}
}
