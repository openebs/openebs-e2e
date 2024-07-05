# \VolumesAPI

All URIs are relative to */v0*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DelShare**](VolumesAPI.md#DelShare) | **Delete** /volumes{volume_id}/share | 
[**DelVolume**](VolumesAPI.md#DelVolume) | **Delete** /volumes/{volume_id} | 
[**DelVolumeShutdownTargets**](VolumesAPI.md#DelVolumeShutdownTargets) | **Delete** /volumes/{volume_id}/shutdown_targets | 
[**DelVolumeTarget**](VolumesAPI.md#DelVolumeTarget) | **Delete** /volumes/{volume_id}/target | 
[**GetRebuildHistory**](VolumesAPI.md#GetRebuildHistory) | **Get** /volumes/{volume_id}/rebuild-history | 
[**GetVolume**](VolumesAPI.md#GetVolume) | **Get** /volumes/{volume_id} | 
[**GetVolumes**](VolumesAPI.md#GetVolumes) | **Get** /volumes | 
[**PutSnapshotVolume**](VolumesAPI.md#PutSnapshotVolume) | **Put** /snapshots/{snapshot_id}/volumes/{volume_id} | 
[**PutVolume**](VolumesAPI.md#PutVolume) | **Put** /volumes/{volume_id} | 
[**PutVolumeProperty**](VolumesAPI.md#PutVolumeProperty) | **Put** /volumes/{volume_id}/property | 
[**PutVolumeReplicaCount**](VolumesAPI.md#PutVolumeReplicaCount) | **Put** /volumes/{volume_id}/replica_count/{replica_count} | 
[**PutVolumeShare**](VolumesAPI.md#PutVolumeShare) | **Put** /volumes/{volume_id}/share/{protocol} | 
[**PutVolumeSize**](VolumesAPI.md#PutVolumeSize) | **Put** /volumes/{volume_id}/size | 
[**PutVolumeTarget**](VolumesAPI.md#PutVolumeTarget) | **Put** /volumes/{volume_id}/target | 



## DelShare

> DelShare(ctx, volumeId).Execute()



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
	r, err := apiClient.VolumesAPI.DelShare(context.Background(), volumeId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VolumesAPI.DelShare``: %v\n", err)
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

Other parameters are passed through a pointer to a apiDelShareRequest struct via the builder pattern


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


## DelVolume

> DelVolume(ctx, volumeId).Execute()



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
	r, err := apiClient.VolumesAPI.DelVolume(context.Background(), volumeId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VolumesAPI.DelVolume``: %v\n", err)
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

Other parameters are passed through a pointer to a apiDelVolumeRequest struct via the builder pattern


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


## DelVolumeShutdownTargets

> DelVolumeShutdownTargets(ctx, volumeId).Execute()



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
	r, err := apiClient.VolumesAPI.DelVolumeShutdownTargets(context.Background(), volumeId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VolumesAPI.DelVolumeShutdownTargets``: %v\n", err)
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

Other parameters are passed through a pointer to a apiDelVolumeShutdownTargetsRequest struct via the builder pattern


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


## DelVolumeTarget

> Volume DelVolumeTarget(ctx, volumeId).Force(force).Execute()



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
	force := true // bool | Force unpublish if the node is not online. This should only be used when it is safe to do so, eg: when the node is not coming back up. (optional) (default to false)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.VolumesAPI.DelVolumeTarget(context.Background(), volumeId).Force(force).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VolumesAPI.DelVolumeTarget``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DelVolumeTarget`: Volume
	fmt.Fprintf(os.Stdout, "Response from `VolumesAPI.DelVolumeTarget`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**volumeId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiDelVolumeTargetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **force** | **bool** | Force unpublish if the node is not online. This should only be used when it is safe to do so, eg: when the node is not coming back up. | [default to false]

### Return type

[**Volume**](Volume.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetRebuildHistory

> RebuildHistory GetRebuildHistory(ctx, volumeId).Execute()



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
	resp, r, err := apiClient.VolumesAPI.GetRebuildHistory(context.Background(), volumeId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VolumesAPI.GetRebuildHistory``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetRebuildHistory`: RebuildHistory
	fmt.Fprintf(os.Stdout, "Response from `VolumesAPI.GetRebuildHistory`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**volumeId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetRebuildHistoryRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**RebuildHistory**](RebuildHistory.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetVolume

> Volume GetVolume(ctx, volumeId).Execute()



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
	resp, r, err := apiClient.VolumesAPI.GetVolume(context.Background(), volumeId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VolumesAPI.GetVolume``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetVolume`: Volume
	fmt.Fprintf(os.Stdout, "Response from `VolumesAPI.GetVolume`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**volumeId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetVolumeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**Volume**](Volume.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetVolumes

> Volumes GetVolumes(ctx).MaxEntries(maxEntries).VolumeId(volumeId).StartingToken(startingToken).Execute()



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
	volumeId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | The uuid of a volume to retrieve. This can be used to \"bypass\" the 404 error when a volume does not exist. (optional)
	startingToken := int32(56) // int32 | the offset to start pagination from (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.VolumesAPI.GetVolumes(context.Background()).MaxEntries(maxEntries).VolumeId(volumeId).StartingToken(startingToken).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VolumesAPI.GetVolumes``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetVolumes`: Volumes
	fmt.Fprintf(os.Stdout, "Response from `VolumesAPI.GetVolumes`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGetVolumesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **maxEntries** | **int32** | the maximum number of results to return | [default to 0]
 **volumeId** | **string** | The uuid of a volume to retrieve. This can be used to \&quot;bypass\&quot; the 404 error when a volume does not exist. | 
 **startingToken** | **int32** | the offset to start pagination from | 

### Return type

[**Volumes**](Volumes.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PutSnapshotVolume

> Volume PutSnapshotVolume(ctx, snapshotId, volumeId).CreateVolumeBody(createVolumeBody).Execute()



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
	volumeId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	createVolumeBody := *openapiclient.NewCreateVolumeBody(*openapiclient.NewVolumePolicy(false), int32(123), int64(123), false) // CreateVolumeBody | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.VolumesAPI.PutSnapshotVolume(context.Background(), snapshotId, volumeId).CreateVolumeBody(createVolumeBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VolumesAPI.PutSnapshotVolume``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PutSnapshotVolume`: Volume
	fmt.Fprintf(os.Stdout, "Response from `VolumesAPI.PutSnapshotVolume`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**snapshotId** | **string** |  | 
**volumeId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiPutSnapshotVolumeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **createVolumeBody** | [**CreateVolumeBody**](CreateVolumeBody.md) |  | 

### Return type

[**Volume**](Volume.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PutVolume

> Volume PutVolume(ctx, volumeId).CreateVolumeBody(createVolumeBody).Execute()



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
	createVolumeBody := *openapiclient.NewCreateVolumeBody(*openapiclient.NewVolumePolicy(false), int32(123), int64(123), false) // CreateVolumeBody | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.VolumesAPI.PutVolume(context.Background(), volumeId).CreateVolumeBody(createVolumeBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VolumesAPI.PutVolume``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PutVolume`: Volume
	fmt.Fprintf(os.Stdout, "Response from `VolumesAPI.PutVolume`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**volumeId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiPutVolumeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **createVolumeBody** | [**CreateVolumeBody**](CreateVolumeBody.md) |  | 

### Return type

[**Volume**](Volume.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PutVolumeProperty

> Volume PutVolumeProperty(ctx, volumeId).SetVolumePropertyBody(setVolumePropertyBody).Execute()



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
	setVolumePropertyBody := *openapiclient.NewSetVolumePropertyBody() // SetVolumePropertyBody | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.VolumesAPI.PutVolumeProperty(context.Background(), volumeId).SetVolumePropertyBody(setVolumePropertyBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VolumesAPI.PutVolumeProperty``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PutVolumeProperty`: Volume
	fmt.Fprintf(os.Stdout, "Response from `VolumesAPI.PutVolumeProperty`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**volumeId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiPutVolumePropertyRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **setVolumePropertyBody** | [**SetVolumePropertyBody**](SetVolumePropertyBody.md) |  | 

### Return type

[**Volume**](Volume.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PutVolumeReplicaCount

> Volume PutVolumeReplicaCount(ctx, volumeId, replicaCount).Execute()



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
	replicaCount := int32(56) // int32 | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.VolumesAPI.PutVolumeReplicaCount(context.Background(), volumeId, replicaCount).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VolumesAPI.PutVolumeReplicaCount``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PutVolumeReplicaCount`: Volume
	fmt.Fprintf(os.Stdout, "Response from `VolumesAPI.PutVolumeReplicaCount`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**volumeId** | **string** |  | 
**replicaCount** | **int32** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiPutVolumeReplicaCountRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**Volume**](Volume.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PutVolumeShare

> string PutVolumeShare(ctx, volumeId, protocol).FrontendHost(frontendHost).Execute()



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
	protocol := openapiclient.VolumeShareProtocol("nvmf") // VolumeShareProtocol | 
	frontendHost := "frontendHost_example" // string | Host if specified, is allowed to connect the target. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.VolumesAPI.PutVolumeShare(context.Background(), volumeId, protocol).FrontendHost(frontendHost).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VolumesAPI.PutVolumeShare``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PutVolumeShare`: string
	fmt.Fprintf(os.Stdout, "Response from `VolumesAPI.PutVolumeShare`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**volumeId** | **string** |  | 
**protocol** | [**VolumeShareProtocol**](.md) |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiPutVolumeShareRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **frontendHost** | **string** | Host if specified, is allowed to connect the target. | 

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


## PutVolumeSize

> Volume PutVolumeSize(ctx, volumeId).ResizeVolumeBody(resizeVolumeBody).Execute()





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
	resizeVolumeBody := *openapiclient.NewResizeVolumeBody(int32(123)) // ResizeVolumeBody | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.VolumesAPI.PutVolumeSize(context.Background(), volumeId).ResizeVolumeBody(resizeVolumeBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VolumesAPI.PutVolumeSize``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PutVolumeSize`: Volume
	fmt.Fprintf(os.Stdout, "Response from `VolumesAPI.PutVolumeSize`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**volumeId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiPutVolumeSizeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **resizeVolumeBody** | [**ResizeVolumeBody**](ResizeVolumeBody.md) |  | 

### Return type

[**Volume**](Volume.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PutVolumeTarget

> Volume PutVolumeTarget(ctx, volumeId).PublishVolumeBody(publishVolumeBody).Execute()





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
	publishVolumeBody := *openapiclient.NewPublishVolumeBody(map[string]string{"key": "Inner_example"}, openapiclient.VolumeShareProtocol("nvmf")) // PublishVolumeBody | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.VolumesAPI.PutVolumeTarget(context.Background(), volumeId).PublishVolumeBody(publishVolumeBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VolumesAPI.PutVolumeTarget``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PutVolumeTarget`: Volume
	fmt.Fprintf(os.Stdout, "Response from `VolumesAPI.PutVolumeTarget`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**volumeId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiPutVolumeTargetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **publishVolumeBody** | [**PublishVolumeBody**](PublishVolumeBody.md) |  | 

### Return type

[**Volume**](Volume.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

