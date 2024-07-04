package v1alpha1Ext

import (
	"context"
	"time"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/e2e_config"

	v1alpha12 "github.com/openebs/openebs-e2e/common/custom_resources/api/types/v1alpha1Ext"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

// DiskPoolInterface has methods to work with Mayastor pool resources.
type DiskPoolInterface interface {
	Create(ctxt context.Context, diskpool *v1alpha12.DiskPool, opts metav1.CreateOptions) (*v1alpha12.DiskPool, error)
	Get(ctxt context.Context, name string, opts metav1.GetOptions) (*v1alpha12.DiskPool, error)
	List(ctxt context.Context, opts metav1.ListOptions) (*v1alpha12.DiskPoolList, error)
	Update(ctxt context.Context, diskpool *v1alpha12.DiskPool, opts metav1.UpdateOptions) (*v1alpha12.DiskPool, error)
	Delete(ctxt context.Context, name string, opts metav1.DeleteOptions) error
	Watch(ctxt context.Context, opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

// dspClient implements DiskPoolInterface
type dspClient struct {
	restClient rest.Interface
}

// Create takes the representation of a Mayastor pool and creates it. Returns the server's representation of the pool, and an error if one occurred.
func (c *dspClient) Create(ctxt context.Context, diskpool *v1alpha12.DiskPool, opts metav1.CreateOptions) (*v1alpha12.DiskPool, error) {
	result := v1alpha12.DiskPool{}
	err := c.restClient.
		Post().
		Namespace(common.NSMayastor()).
		Resource(e2e_config.GetConfig().Product.CrdPoolsResourceName).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(diskpool).
		Do(ctxt).
		Into(&result)

	return &result, err
}

// Get takes the name of the Mayastor pool and returns the server's representation of it, and an error if one occurred.
func (c *dspClient) Get(ctxt context.Context, name string, opts metav1.GetOptions) (*v1alpha12.DiskPool, error) {
	result := v1alpha12.DiskPool{}
	err := c.restClient.
		Get().
		Namespace(common.NSMayastor()).
		Resource(e2e_config.GetConfig().Product.CrdPoolsResourceName).
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctxt).
		Into(&result)

	return &result, err
}

// List takes the label and field selectors, and returns a list matching Mayastor Pool, and an error if one occurred.
func (c *dspClient) List(ctxt context.Context, opts metav1.ListOptions) (*v1alpha12.DiskPoolList, error) {
	result := v1alpha12.DiskPoolList{}
	err := c.restClient.
		Get().
		Namespace(common.NSMayastor()).
		Resource(e2e_config.GetConfig().Product.CrdPoolsResourceName).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctxt).
		Into(&result)

	return &result, err
}

// Update takes the representation of a Mayastor pool and updates it. Returns the server's representation of the pool, and an error if one occurred.
func (c *dspClient) Update(ctxt context.Context, diskpool *v1alpha12.DiskPool, opts metav1.UpdateOptions) (*v1alpha12.DiskPool, error) {
	result := v1alpha12.DiskPool{}
	err := c.restClient.
		Put().
		Namespace(common.NSMayastor()).
		Resource(e2e_config.GetConfig().Product.CrdPoolsResourceName).
		Name(diskpool.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(diskpool).
		Do(ctxt).
		Into(&result)

	return &result, err
}

// Delete takes the name of the Mayastor pool and deletes it. Returns error if one occurred.
func (c *dspClient) Delete(ctxt context.Context, name string, opts metav1.DeleteOptions) error {
	return c.restClient.
		Delete().
		Namespace(common.NSMayastor()).
		Resource(e2e_config.GetConfig().Product.CrdPoolsResourceName).
		Name(name).
		Body(&opts).
		Do(ctxt).
		Error()
}

// Watch takes the label and field selectors, and returns a watch.Interface the watches matching Mayastor pools, and an error if one occurred.
func (c *dspClient) Watch(ctxt context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(common.NSMayastor()).
		Resource(e2e_config.GetConfig().Product.CrdPoolsResourceName).
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctxt)
}
