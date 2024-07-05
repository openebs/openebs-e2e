# \AppNodesAPI

All URIs are relative to */v0*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DeregisterAppNode**](AppNodesAPI.md#DeregisterAppNode) | **Delete** /app-nodes/{app_node_id} | 
[**GetAppNode**](AppNodesAPI.md#GetAppNode) | **Get** /app-nodes/{app_node_id} | 
[**GetAppNodes**](AppNodesAPI.md#GetAppNodes) | **Get** /app-nodes | 
[**RegisterAppNode**](AppNodesAPI.md#RegisterAppNode) | **Put** /app-nodes/{app_node_id} | 



## DeregisterAppNode

> DeregisterAppNode(ctx, appNodeId).Execute()



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
	appNodeId := "appNodeId_example" // string | ID of the app node to be deregistered.

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.AppNodesAPI.DeregisterAppNode(context.Background(), appNodeId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AppNodesAPI.DeregisterAppNode``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**appNodeId** | **string** | ID of the app node to be deregistered. | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeregisterAppNodeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetAppNode

> AppNode GetAppNode(ctx, appNodeId).Execute()



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
	appNodeId := "appNodeId_example" // string | Id of the app node.

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AppNodesAPI.GetAppNode(context.Background(), appNodeId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AppNodesAPI.GetAppNode``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAppNode`: AppNode
	fmt.Fprintf(os.Stdout, "Response from `AppNodesAPI.GetAppNode`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**appNodeId** | **string** | Id of the app node. | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetAppNodeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**AppNode**](AppNode.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetAppNodes

> AppNodes GetAppNodes(ctx).MaxEntries(maxEntries).StartingToken(startingToken).Execute()



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
	maxEntries := int32(56) // int32 | The maximum number of results to return. (default to 0)
	startingToken := int32(56) // int32 | The offset to start pagination from. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AppNodesAPI.GetAppNodes(context.Background()).MaxEntries(maxEntries).StartingToken(startingToken).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AppNodesAPI.GetAppNodes``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAppNodes`: AppNodes
	fmt.Fprintf(os.Stdout, "Response from `AppNodesAPI.GetAppNodes`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGetAppNodesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **maxEntries** | **int32** | The maximum number of results to return. | [default to 0]
 **startingToken** | **int32** | The offset to start pagination from. | 

### Return type

[**AppNodes**](AppNodes.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RegisterAppNode

> RegisterAppNode(ctx, appNodeId).RegisterAppNode(registerAppNode).Execute()



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
	appNodeId := "appNodeId_example" // string | 
	registerAppNode := *openapiclient.NewRegisterAppNode("Endpoint_example") // RegisterAppNode | Contents of the app node to be registered.

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.AppNodesAPI.RegisterAppNode(context.Background(), appNodeId).RegisterAppNode(registerAppNode).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AppNodesAPI.RegisterAppNode``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**appNodeId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiRegisterAppNodeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **registerAppNode** | [**RegisterAppNode**](RegisterAppNode.md) | Contents of the app node to be registered. | 

### Return type

 (empty response body)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

