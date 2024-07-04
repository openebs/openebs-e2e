# \ChildrenAPI

All URIs are relative to */v0*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DelNexusChild**](ChildrenAPI.md#DelNexusChild) | **Delete** /nexuses/{nexus_id}/children/{child_id} | 
[**DelNodeNexusChild**](ChildrenAPI.md#DelNodeNexusChild) | **Delete** /nodes/{node_id}/nexuses/{nexus_id}/children/{child_id} | 
[**GetNexusChild**](ChildrenAPI.md#GetNexusChild) | **Get** /nexuses/{nexus_id}/children/{child_id} | 
[**GetNexusChildren**](ChildrenAPI.md#GetNexusChildren) | **Get** /nexuses/{nexus_id}/children | 
[**GetNodeNexusChild**](ChildrenAPI.md#GetNodeNexusChild) | **Get** /nodes/{node_id}/nexuses/{nexus_id}/children/{child_id} | 
[**GetNodeNexusChildren**](ChildrenAPI.md#GetNodeNexusChildren) | **Get** /nodes/{node_id}/nexuses/{nexus_id}/children | 
[**PutNexusChild**](ChildrenAPI.md#PutNexusChild) | **Put** /nexuses/{nexus_id}/children/{child_id} | 
[**PutNodeNexusChild**](ChildrenAPI.md#PutNodeNexusChild) | **Put** /nodes/{node_id}/nexuses/{nexus_id}/children/{child_id} | 



## DelNexusChild

> DelNexusChild(ctx, nexusId, childId).Execute()



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
	nexusId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	childId := "childId_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ChildrenAPI.DelNexusChild(context.Background(), nexusId, childId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ChildrenAPI.DelNexusChild``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**nexusId** | **string** |  | 
**childId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiDelNexusChildRequest struct via the builder pattern


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


## DelNodeNexusChild

> DelNodeNexusChild(ctx, nodeId, nexusId, childId).Execute()



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
	nodeId := "nodeId_example" // string | 
	nexusId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	childId := "childId_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ChildrenAPI.DelNodeNexusChild(context.Background(), nodeId, nexusId, childId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ChildrenAPI.DelNodeNexusChild``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**nodeId** | **string** |  | 
**nexusId** | **string** |  | 
**childId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiDelNodeNexusChildRequest struct via the builder pattern


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


## GetNexusChild

> Child GetNexusChild(ctx, nexusId, childId).Execute()



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
	nexusId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	childId := "childId_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ChildrenAPI.GetNexusChild(context.Background(), nexusId, childId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ChildrenAPI.GetNexusChild``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetNexusChild`: Child
	fmt.Fprintf(os.Stdout, "Response from `ChildrenAPI.GetNexusChild`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**nexusId** | **string** |  | 
**childId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetNexusChildRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**Child**](Child.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetNexusChildren

> []Child GetNexusChildren(ctx, nexusId).Execute()



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
	nexusId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ChildrenAPI.GetNexusChildren(context.Background(), nexusId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ChildrenAPI.GetNexusChildren``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetNexusChildren`: []Child
	fmt.Fprintf(os.Stdout, "Response from `ChildrenAPI.GetNexusChildren`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**nexusId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetNexusChildrenRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**[]Child**](Child.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetNodeNexusChild

> Child GetNodeNexusChild(ctx, nodeId, nexusId, childId).Execute()



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
	nodeId := "nodeId_example" // string | 
	nexusId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	childId := "childId_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ChildrenAPI.GetNodeNexusChild(context.Background(), nodeId, nexusId, childId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ChildrenAPI.GetNodeNexusChild``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetNodeNexusChild`: Child
	fmt.Fprintf(os.Stdout, "Response from `ChildrenAPI.GetNodeNexusChild`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**nodeId** | **string** |  | 
**nexusId** | **string** |  | 
**childId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetNodeNexusChildRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------




### Return type

[**Child**](Child.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetNodeNexusChildren

> []Child GetNodeNexusChildren(ctx, nodeId, nexusId).Execute()



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
	nodeId := "nodeId_example" // string | 
	nexusId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ChildrenAPI.GetNodeNexusChildren(context.Background(), nodeId, nexusId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ChildrenAPI.GetNodeNexusChildren``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetNodeNexusChildren`: []Child
	fmt.Fprintf(os.Stdout, "Response from `ChildrenAPI.GetNodeNexusChildren`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**nodeId** | **string** |  | 
**nexusId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetNodeNexusChildrenRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**[]Child**](Child.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PutNexusChild

> Child PutNexusChild(ctx, nexusId, childId).Execute()



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
	nexusId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	childId := "childId_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ChildrenAPI.PutNexusChild(context.Background(), nexusId, childId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ChildrenAPI.PutNexusChild``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PutNexusChild`: Child
	fmt.Fprintf(os.Stdout, "Response from `ChildrenAPI.PutNexusChild`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**nexusId** | **string** |  | 
**childId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiPutNexusChildRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**Child**](Child.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PutNodeNexusChild

> Child PutNodeNexusChild(ctx, nodeId, nexusId, childId).Execute()



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
	nodeId := "nodeId_example" // string | 
	nexusId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	childId := "childId_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ChildrenAPI.PutNodeNexusChild(context.Background(), nodeId, nexusId, childId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ChildrenAPI.PutNodeNexusChild``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PutNodeNexusChild`: Child
	fmt.Fprintf(os.Stdout, "Response from `ChildrenAPI.PutNodeNexusChild`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**nodeId** | **string** |  | 
**nexusId** | **string** |  | 
**childId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiPutNodeNexusChildRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------




### Return type

[**Child**](Child.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

