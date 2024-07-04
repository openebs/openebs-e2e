# \SnapshotsAPI

All URIs are relative to */v0*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DelSnapshot**](SnapshotsAPI.md#DelSnapshot) | **Delete** /volumes/snapshots/{snapshot_id} | 
[**DelVolumeSnapshot**](SnapshotsAPI.md#DelVolumeSnapshot) | **Delete** /volumes/{volume_id}/snapshots/{snapshot_id} | 
[**GetVolumeSnapshot**](SnapshotsAPI.md#GetVolumeSnapshot) | **Get** /volumes/{volume_id}/snapshots/{snapshot_id} | 
[**GetVolumeSnapshots**](SnapshotsAPI.md#GetVolumeSnapshots) | **Get** /volumes/{volume_id}/snapshots | 
[**GetVolumesSnapshot**](SnapshotsAPI.md#GetVolumesSnapshot) | **Get** /volumes/snapshots/{snapshot_id} | 
[**GetVolumesSnapshots**](SnapshotsAPI.md#GetVolumesSnapshots) | **Get** /volumes/snapshots | 
[**PutVolumeSnapshot**](SnapshotsAPI.md#PutVolumeSnapshot) | **Put** /volumes/{volume_id}/snapshots/{snapshot_id} | 



## DelSnapshot

> DelSnapshot(ctx, snapshotId).Execute()



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
	snapshotId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.SnapshotsAPI.DelSnapshot(context.Background(), snapshotId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SnapshotsAPI.DelSnapshot``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**snapshotId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiDelSnapshotRequest struct via the builder pattern


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


## DelVolumeSnapshot

> DelVolumeSnapshot(ctx, volumeId, snapshotId).Execute()



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
	snapshotId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.SnapshotsAPI.DelVolumeSnapshot(context.Background(), volumeId, snapshotId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SnapshotsAPI.DelVolumeSnapshot``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**volumeId** | **string** |  | 
**snapshotId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiDelVolumeSnapshotRequest struct via the builder pattern


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


## GetVolumeSnapshot

> VolumeSnapshot GetVolumeSnapshot(ctx, volumeId, snapshotId).Execute()



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
	snapshotId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SnapshotsAPI.GetVolumeSnapshot(context.Background(), volumeId, snapshotId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SnapshotsAPI.GetVolumeSnapshot``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetVolumeSnapshot`: VolumeSnapshot
	fmt.Fprintf(os.Stdout, "Response from `SnapshotsAPI.GetVolumeSnapshot`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**volumeId** | **string** |  | 
**snapshotId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetVolumeSnapshotRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**VolumeSnapshot**](VolumeSnapshot.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetVolumeSnapshots

> VolumeSnapshots GetVolumeSnapshots(ctx, volumeId).MaxEntries(maxEntries).StartingToken(startingToken).Execute()



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
	maxEntries := int32(56) // int32 | the maximum number of results to return (default to 0)
	startingToken := int32(56) // int32 | the offset to start pagination from (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SnapshotsAPI.GetVolumeSnapshots(context.Background(), volumeId).MaxEntries(maxEntries).StartingToken(startingToken).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SnapshotsAPI.GetVolumeSnapshots``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetVolumeSnapshots`: VolumeSnapshots
	fmt.Fprintf(os.Stdout, "Response from `SnapshotsAPI.GetVolumeSnapshots`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**volumeId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetVolumeSnapshotsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **maxEntries** | **int32** | the maximum number of results to return | [default to 0]
 **startingToken** | **int32** | the offset to start pagination from | 

### Return type

[**VolumeSnapshots**](VolumeSnapshots.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetVolumesSnapshot

> VolumeSnapshot GetVolumesSnapshot(ctx, snapshotId).Execute()



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
	snapshotId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SnapshotsAPI.GetVolumesSnapshot(context.Background(), snapshotId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SnapshotsAPI.GetVolumesSnapshot``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetVolumesSnapshot`: VolumeSnapshot
	fmt.Fprintf(os.Stdout, "Response from `SnapshotsAPI.GetVolumesSnapshot`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**snapshotId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetVolumesSnapshotRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**VolumeSnapshot**](VolumeSnapshot.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetVolumesSnapshots

> VolumeSnapshots GetVolumesSnapshots(ctx).MaxEntries(maxEntries).SnapshotId(snapshotId).VolumeId(volumeId).StartingToken(startingToken).Execute()



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
	maxEntries := int32(56) // int32 | the maximum number of results to return (default to 0)
	snapshotId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | The uuid of the snapshot to retrieve. (optional)
	volumeId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | The uuid of the snapshots source volume. (optional)
	startingToken := int32(56) // int32 | the offset to start pagination from (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SnapshotsAPI.GetVolumesSnapshots(context.Background()).MaxEntries(maxEntries).SnapshotId(snapshotId).VolumeId(volumeId).StartingToken(startingToken).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SnapshotsAPI.GetVolumesSnapshots``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetVolumesSnapshots`: VolumeSnapshots
	fmt.Fprintf(os.Stdout, "Response from `SnapshotsAPI.GetVolumesSnapshots`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGetVolumesSnapshotsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **maxEntries** | **int32** | the maximum number of results to return | [default to 0]
 **snapshotId** | **string** | The uuid of the snapshot to retrieve. | 
 **volumeId** | **string** | The uuid of the snapshots source volume. | 
 **startingToken** | **int32** | the offset to start pagination from | 

### Return type

[**VolumeSnapshots**](VolumeSnapshots.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PutVolumeSnapshot

> VolumeSnapshot PutVolumeSnapshot(ctx, volumeId, snapshotId).Execute()



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
	snapshotId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SnapshotsAPI.PutVolumeSnapshot(context.Background(), volumeId, snapshotId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SnapshotsAPI.PutVolumeSnapshot``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PutVolumeSnapshot`: VolumeSnapshot
	fmt.Fprintf(os.Stdout, "Response from `SnapshotsAPI.PutVolumeSnapshot`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**volumeId** | **string** |  | 
**snapshotId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiPutVolumeSnapshotRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**VolumeSnapshot**](VolumeSnapshot.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

