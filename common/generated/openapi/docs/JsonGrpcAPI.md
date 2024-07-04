# \JsonGrpcAPI

All URIs are relative to */v0*

Method | HTTP request | Description
------------- | ------------- | -------------
[**PutNodeJsongrpc**](JsonGrpcAPI.md#PutNodeJsongrpc) | **Put** /nodes/{node}/jsongrpc/{method} | 



## PutNodeJsongrpc

> map[string]interface{} PutNodeJsongrpc(ctx, node, method).Body(body).Execute()



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
	method := "method_example" // string | 
	body := map[string]interface{}{ ... } // map[string]interface{} | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.JsonGrpcAPI.PutNodeJsongrpc(context.Background(), node, method).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `JsonGrpcAPI.PutNodeJsongrpc``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PutNodeJsongrpc`: map[string]interface{}
	fmt.Fprintf(os.Stdout, "Response from `JsonGrpcAPI.PutNodeJsongrpc`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**node** | **string** |  | 
**method** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiPutNodeJsongrpcRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **body** | **map[string]interface{}** |  | 

### Return type

**map[string]interface{}**

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

