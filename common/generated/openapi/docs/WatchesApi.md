# \WatchesAPI

All URIs are relative to */v0*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DelWatchVolume**](WatchesAPI.md#DelWatchVolume) | **Delete** /watches/volumes/{volume_id} | 
[**GetWatchVolume**](WatchesAPI.md#GetWatchVolume) | **Get** /watches/volumes/{volume_id} | 
[**PutWatchVolume**](WatchesAPI.md#PutWatchVolume) | **Put** /watches/volumes/{volume_id} | 



## DelWatchVolume

> DelWatchVolume(ctx, volumeId).Callback(callback).Execute()



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
	volumeId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	callback := "callback_example" // string | URL callback

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.WatchesAPI.DelWatchVolume(context.Background(), volumeId).Callback(callback).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WatchesAPI.DelWatchVolume``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**volumeId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiDelWatchVolumeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **callback** | **string** | URL callback | 

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


## GetWatchVolume

> []RestWatch GetWatchVolume(ctx, volumeId).Execute()



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
	volumeId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.WatchesAPI.GetWatchVolume(context.Background(), volumeId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WatchesAPI.GetWatchVolume``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetWatchVolume`: []RestWatch
	fmt.Fprintf(os.Stdout, "Response from `WatchesAPI.GetWatchVolume`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**volumeId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetWatchVolumeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**[]RestWatch**](RestWatch.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PutWatchVolume

> PutWatchVolume(ctx, volumeId).Callback(callback).Execute()



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
	volumeId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	callback := "callback_example" // string | URL callback

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.WatchesAPI.PutWatchVolume(context.Background(), volumeId).Callback(callback).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WatchesAPI.PutWatchVolume``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**volumeId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiPutWatchVolumeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **callback** | **string** | URL callback | 

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

