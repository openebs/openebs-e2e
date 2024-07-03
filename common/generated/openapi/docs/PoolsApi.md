# \PoolsAPI

All URIs are relative to */v0*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DelNodePool**](PoolsAPI.md#DelNodePool) | **Delete** /nodes/{node_id}/pools/{pool_id} | 
[**DelPool**](PoolsAPI.md#DelPool) | **Delete** /pools/{pool_id} | 
[**GetNodePool**](PoolsAPI.md#GetNodePool) | **Get** /nodes/{node_id}/pools/{pool_id} | 
[**GetNodePools**](PoolsAPI.md#GetNodePools) | **Get** /nodes/{id}/pools | 
[**GetPool**](PoolsAPI.md#GetPool) | **Get** /pools/{pool_id} | 
[**GetPools**](PoolsAPI.md#GetPools) | **Get** /pools | 
[**PutNodePool**](PoolsAPI.md#PutNodePool) | **Put** /nodes/{node_id}/pools/{pool_id} | 



## DelNodePool

> DelNodePool(ctx, nodeId, poolId).Execute()



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
	poolId := "poolId_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.PoolsAPI.DelNodePool(context.Background(), nodeId, poolId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PoolsAPI.DelNodePool``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**nodeId** | **string** |  | 
**poolId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiDelNodePoolRequest struct via the builder pattern


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


## DelPool

> DelPool(ctx, poolId).Execute()



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
	poolId := "poolId_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.PoolsAPI.DelPool(context.Background(), poolId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PoolsAPI.DelPool``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**poolId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiDelPoolRequest struct via the builder pattern


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


## GetNodePool

> Pool GetNodePool(ctx, nodeId, poolId).Execute()



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
	poolId := "poolId_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PoolsAPI.GetNodePool(context.Background(), nodeId, poolId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PoolsAPI.GetNodePool``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetNodePool`: Pool
	fmt.Fprintf(os.Stdout, "Response from `PoolsAPI.GetNodePool`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**nodeId** | **string** |  | 
**poolId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetNodePoolRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**Pool**](Pool.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetNodePools

> []Pool GetNodePools(ctx, id).Execute()



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
	id := "id_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PoolsAPI.GetNodePools(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PoolsAPI.GetNodePools``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetNodePools`: []Pool
	fmt.Fprintf(os.Stdout, "Response from `PoolsAPI.GetNodePools`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetNodePoolsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**[]Pool**](Pool.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetPool

> Pool GetPool(ctx, poolId).Execute()



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
	poolId := "poolId_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PoolsAPI.GetPool(context.Background(), poolId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PoolsAPI.GetPool``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetPool`: Pool
	fmt.Fprintf(os.Stdout, "Response from `PoolsAPI.GetPool`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**poolId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetPoolRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**Pool**](Pool.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetPools

> []Pool GetPools(ctx).VolumeId(volumeId).Execute()



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
	volumeId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PoolsAPI.GetPools(context.Background()).VolumeId(volumeId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PoolsAPI.GetPools``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetPools`: []Pool
	fmt.Fprintf(os.Stdout, "Response from `PoolsAPI.GetPools`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGetPoolsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **volumeId** | **string** |  | 

### Return type

[**[]Pool**](Pool.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PutNodePool

> Pool PutNodePool(ctx, nodeId, poolId).CreatePoolBody(createPoolBody).Execute()



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
	poolId := "poolId_example" // string | 
	createPoolBody := *openapiclient.NewCreatePoolBody([]string{"malloc:///disk?size_mb=100"}) // CreatePoolBody | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PoolsAPI.PutNodePool(context.Background(), nodeId, poolId).CreatePoolBody(createPoolBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PoolsAPI.PutNodePool``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PutNodePool`: Pool
	fmt.Fprintf(os.Stdout, "Response from `PoolsAPI.PutNodePool`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**nodeId** | **string** |  | 
**poolId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiPutNodePoolRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **createPoolBody** | [**CreatePoolBody**](CreatePoolBody.md) |  | 

### Return type

[**Pool**](Pool.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

