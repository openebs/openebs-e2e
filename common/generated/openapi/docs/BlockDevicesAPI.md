# \BlockDevicesAPI

All URIs are relative to */v0*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetNodeBlockDevices**](BlockDevicesAPI.md#GetNodeBlockDevices) | **Get** /nodes/{node}/block_devices | 



## GetNodeBlockDevices

> []BlockDevice GetNodeBlockDevices(ctx, node).All(all).Execute()



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	node := "node_example" // string | 
	all := true // bool | specifies whether to list all devices or only usable ones (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.BlockDevicesAPI.GetNodeBlockDevices(context.Background(), node).All(all).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `BlockDevicesAPI.GetNodeBlockDevices``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetNodeBlockDevices`: []BlockDevice
	fmt.Fprintf(os.Stdout, "Response from `BlockDevicesAPI.GetNodeBlockDevices`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**node** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetNodeBlockDevicesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **all** | **bool** | specifies whether to list all devices or only usable ones | 

### Return type

[**[]BlockDevice**](BlockDevice.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

