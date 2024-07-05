# \NexusesAPI

All URIs are relative to */v0*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DelNexus**](NexusesAPI.md#DelNexus) | **Delete** /nexuses/{nexus_id} | 
[**DelNodeNexus**](NexusesAPI.md#DelNodeNexus) | **Delete** /nodes/{node_id}/nexuses/{nexus_id} | 
[**DelNodeNexusShare**](NexusesAPI.md#DelNodeNexusShare) | **Delete** /nodes/{node_id}/nexuses/{nexus_id}/share | 
[**GetNexus**](NexusesAPI.md#GetNexus) | **Get** /nexuses/{nexus_id} | 
[**GetNexuses**](NexusesAPI.md#GetNexuses) | **Get** /nexuses | 
[**GetNodeNexus**](NexusesAPI.md#GetNodeNexus) | **Get** /nodes/{node_id}/nexuses/{nexus_id} | 
[**GetNodeNexuses**](NexusesAPI.md#GetNodeNexuses) | **Get** /nodes/{id}/nexuses | 
[**PutNodeNexus**](NexusesAPI.md#PutNodeNexus) | **Put** /nodes/{node_id}/nexuses/{nexus_id} | 
[**PutNodeNexusShare**](NexusesAPI.md#PutNodeNexusShare) | **Put** /nodes/{node_id}/nexuses/{nexus_id}/share/{protocol} | 



## DelNexus

> DelNexus(ctx, nexusId).Execute()



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
	r, err := apiClient.NexusesAPI.DelNexus(context.Background(), nexusId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NexusesAPI.DelNexus``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**nexusId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiDelNexusRequest struct via the builder pattern


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


## DelNodeNexus

> DelNodeNexus(ctx, nodeId, nexusId).Execute()



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
	r, err := apiClient.NexusesAPI.DelNodeNexus(context.Background(), nodeId, nexusId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NexusesAPI.DelNodeNexus``: %v\n", err)
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

### Other Parameters

Other parameters are passed through a pointer to a apiDelNodeNexusRequest struct via the builder pattern


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


## DelNodeNexusShare

> DelNodeNexusShare(ctx, nodeId, nexusId).Execute()



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
	r, err := apiClient.NexusesAPI.DelNodeNexusShare(context.Background(), nodeId, nexusId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NexusesAPI.DelNodeNexusShare``: %v\n", err)
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

### Other Parameters

Other parameters are passed through a pointer to a apiDelNodeNexusShareRequest struct via the builder pattern


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


## GetNexus

> Nexus GetNexus(ctx, nexusId).Execute()



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
	resp, r, err := apiClient.NexusesAPI.GetNexus(context.Background(), nexusId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NexusesAPI.GetNexus``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetNexus`: Nexus
	fmt.Fprintf(os.Stdout, "Response from `NexusesAPI.GetNexus`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**nexusId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetNexusRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**Nexus**](Nexus.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetNexuses

> []Nexus GetNexuses(ctx).Execute()



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

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.NexusesAPI.GetNexuses(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NexusesAPI.GetNexuses``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetNexuses`: []Nexus
	fmt.Fprintf(os.Stdout, "Response from `NexusesAPI.GetNexuses`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetNexusesRequest struct via the builder pattern


### Return type

[**[]Nexus**](Nexus.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetNodeNexus

> Nexus GetNodeNexus(ctx, nodeId, nexusId).Execute()



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
	resp, r, err := apiClient.NexusesAPI.GetNodeNexus(context.Background(), nodeId, nexusId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NexusesAPI.GetNodeNexus``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetNodeNexus`: Nexus
	fmt.Fprintf(os.Stdout, "Response from `NexusesAPI.GetNodeNexus`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**nodeId** | **string** |  | 
**nexusId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetNodeNexusRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**Nexus**](Nexus.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetNodeNexuses

> []Nexus GetNodeNexuses(ctx, id).Execute()



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
	resp, r, err := apiClient.NexusesAPI.GetNodeNexuses(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NexusesAPI.GetNodeNexuses``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetNodeNexuses`: []Nexus
	fmt.Fprintf(os.Stdout, "Response from `NexusesAPI.GetNodeNexuses`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetNodeNexusesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**[]Nexus**](Nexus.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PutNodeNexus

> Nexus PutNodeNexus(ctx, nodeId, nexusId).CreateNexusBody(createNexusBody).Execute()



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
	createNexusBody := *openapiclient.NewCreateNexusBody([]string{"Children_example"}, int64(123)) // CreateNexusBody | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.NexusesAPI.PutNodeNexus(context.Background(), nodeId, nexusId).CreateNexusBody(createNexusBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NexusesAPI.PutNodeNexus``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PutNodeNexus`: Nexus
	fmt.Fprintf(os.Stdout, "Response from `NexusesAPI.PutNodeNexus`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**nodeId** | **string** |  | 
**nexusId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiPutNodeNexusRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **createNexusBody** | [**CreateNexusBody**](CreateNexusBody.md) |  | 

### Return type

[**Nexus**](Nexus.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PutNodeNexusShare

> string PutNodeNexusShare(ctx, nodeId, nexusId, protocol).Execute()



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
	protocol := openapiclient.NexusShareProtocol("nvmf") // NexusShareProtocol | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.NexusesAPI.PutNodeNexusShare(context.Background(), nodeId, nexusId, protocol).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NexusesAPI.PutNodeNexusShare``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PutNodeNexusShare`: string
	fmt.Fprintf(os.Stdout, "Response from `NexusesAPI.PutNodeNexusShare`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**nodeId** | **string** |  | 
**nexusId** | **string** |  | 
**protocol** | [**NexusShareProtocol**](.md) |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiPutNodeNexusShareRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------




### Return type

**string**

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

