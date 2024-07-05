# \ReplicasAPI

All URIs are relative to */v0*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DelNodePoolReplica**](ReplicasAPI.md#DelNodePoolReplica) | **Delete** /nodes/{node_id}/pools/{pool_id}/replicas/{replica_id} | 
[**DelNodePoolReplicaShare**](ReplicasAPI.md#DelNodePoolReplicaShare) | **Delete** /nodes/{node_id}/pools/{pool_id}/replicas/{replica_id}/share | 
[**DelPoolReplica**](ReplicasAPI.md#DelPoolReplica) | **Delete** /pools/{pool_id}/replicas/{replica_id} | 
[**DelPoolReplicaShare**](ReplicasAPI.md#DelPoolReplicaShare) | **Delete** /pools/{pool_id}/replicas/{replica_id}/share | 
[**GetNodePoolReplica**](ReplicasAPI.md#GetNodePoolReplica) | **Get** /nodes/{node_id}/pools/{pool_id}/replicas/{replica_id} | 
[**GetNodePoolReplicas**](ReplicasAPI.md#GetNodePoolReplicas) | **Get** /nodes/{node_id}/pools/{pool_id}/replicas | 
[**GetNodeReplicas**](ReplicasAPI.md#GetNodeReplicas) | **Get** /nodes/{id}/replicas | 
[**GetReplica**](ReplicasAPI.md#GetReplica) | **Get** /replicas/{id} | 
[**GetReplicas**](ReplicasAPI.md#GetReplicas) | **Get** /replicas | 
[**PutNodePoolReplica**](ReplicasAPI.md#PutNodePoolReplica) | **Put** /nodes/{node_id}/pools/{pool_id}/replicas/{replica_id} | 
[**PutNodePoolReplicaShare**](ReplicasAPI.md#PutNodePoolReplicaShare) | **Put** /nodes/{node_id}/pools/{pool_id}/replicas/{replica_id}/share/nvmf | 
[**PutPoolReplica**](ReplicasAPI.md#PutPoolReplica) | **Put** /pools/{pool_id}/replicas/{replica_id} | 
[**PutPoolReplicaShare**](ReplicasAPI.md#PutPoolReplicaShare) | **Put** /pools/{pool_id}/replicas/{replica_id}/share/nvmf | 



## DelNodePoolReplica

> DelNodePoolReplica(ctx, nodeId, poolId, replicaId).Execute()



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
	replicaId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ReplicasAPI.DelNodePoolReplica(context.Background(), nodeId, poolId, replicaId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ReplicasAPI.DelNodePoolReplica``: %v\n", err)
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
**replicaId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiDelNodePoolReplicaRequest struct via the builder pattern


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


## DelNodePoolReplicaShare

> DelNodePoolReplicaShare(ctx, nodeId, poolId, replicaId).Execute()



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
	replicaId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ReplicasAPI.DelNodePoolReplicaShare(context.Background(), nodeId, poolId, replicaId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ReplicasAPI.DelNodePoolReplicaShare``: %v\n", err)
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
**replicaId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiDelNodePoolReplicaShareRequest struct via the builder pattern


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


## DelPoolReplica

> DelPoolReplica(ctx, poolId, replicaId).Execute()



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
	replicaId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ReplicasAPI.DelPoolReplica(context.Background(), poolId, replicaId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ReplicasAPI.DelPoolReplica``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**poolId** | **string** |  | 
**replicaId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiDelPoolReplicaRequest struct via the builder pattern


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


## DelPoolReplicaShare

> DelPoolReplicaShare(ctx, poolId, replicaId).Execute()



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
	replicaId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ReplicasAPI.DelPoolReplicaShare(context.Background(), poolId, replicaId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ReplicasAPI.DelPoolReplicaShare``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**poolId** | **string** |  | 
**replicaId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiDelPoolReplicaShareRequest struct via the builder pattern


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


## GetNodePoolReplica

> Replica GetNodePoolReplica(ctx, nodeId, poolId, replicaId).Execute()



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
	replicaId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ReplicasAPI.GetNodePoolReplica(context.Background(), nodeId, poolId, replicaId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ReplicasAPI.GetNodePoolReplica``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetNodePoolReplica`: Replica
	fmt.Fprintf(os.Stdout, "Response from `ReplicasAPI.GetNodePoolReplica`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**nodeId** | **string** |  | 
**poolId** | **string** |  | 
**replicaId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetNodePoolReplicaRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------




### Return type

[**Replica**](Replica.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetNodePoolReplicas

> []Replica GetNodePoolReplicas(ctx, nodeId, poolId).Execute()



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
	resp, r, err := apiClient.ReplicasAPI.GetNodePoolReplicas(context.Background(), nodeId, poolId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ReplicasAPI.GetNodePoolReplicas``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetNodePoolReplicas`: []Replica
	fmt.Fprintf(os.Stdout, "Response from `ReplicasAPI.GetNodePoolReplicas`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**nodeId** | **string** |  | 
**poolId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetNodePoolReplicasRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**[]Replica**](Replica.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetNodeReplicas

> []Replica GetNodeReplicas(ctx, id).Execute()



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
	resp, r, err := apiClient.ReplicasAPI.GetNodeReplicas(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ReplicasAPI.GetNodeReplicas``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetNodeReplicas`: []Replica
	fmt.Fprintf(os.Stdout, "Response from `ReplicasAPI.GetNodeReplicas`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetNodeReplicasRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**[]Replica**](Replica.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetReplica

> Replica GetReplica(ctx, id).Execute()



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
	id := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ReplicasAPI.GetReplica(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ReplicasAPI.GetReplica``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetReplica`: Replica
	fmt.Fprintf(os.Stdout, "Response from `ReplicasAPI.GetReplica`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetReplicaRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**Replica**](Replica.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetReplicas

> []Replica GetReplicas(ctx).Execute()



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
	resp, r, err := apiClient.ReplicasAPI.GetReplicas(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ReplicasAPI.GetReplicas``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetReplicas`: []Replica
	fmt.Fprintf(os.Stdout, "Response from `ReplicasAPI.GetReplicas`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetReplicasRequest struct via the builder pattern


### Return type

[**[]Replica**](Replica.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PutNodePoolReplica

> Replica PutNodePoolReplica(ctx, nodeId, poolId, replicaId).CreateReplicaBody(createReplicaBody).Execute()



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
	replicaId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	createReplicaBody := *openapiclient.NewCreateReplicaBody(int64(123), false) // CreateReplicaBody | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ReplicasAPI.PutNodePoolReplica(context.Background(), nodeId, poolId, replicaId).CreateReplicaBody(createReplicaBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ReplicasAPI.PutNodePoolReplica``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PutNodePoolReplica`: Replica
	fmt.Fprintf(os.Stdout, "Response from `ReplicasAPI.PutNodePoolReplica`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**nodeId** | **string** |  | 
**poolId** | **string** |  | 
**replicaId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiPutNodePoolReplicaRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



 **createReplicaBody** | [**CreateReplicaBody**](CreateReplicaBody.md) |  | 

### Return type

[**Replica**](Replica.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PutNodePoolReplicaShare

> string PutNodePoolReplicaShare(ctx, nodeId, poolId, replicaId).AllowedHosts(allowedHosts).Execute()



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
	replicaId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	allowedHosts := []string{"nqn.2014-08.org.nvmexpress:uuid:804b1e8c-b42d-4d15-92b4-7c4e4d0f507"} // []string | NQNs of hosts allowed to connect to this replica (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ReplicasAPI.PutNodePoolReplicaShare(context.Background(), nodeId, poolId, replicaId).AllowedHosts(allowedHosts).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ReplicasAPI.PutNodePoolReplicaShare``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PutNodePoolReplicaShare`: string
	fmt.Fprintf(os.Stdout, "Response from `ReplicasAPI.PutNodePoolReplicaShare`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**nodeId** | **string** |  | 
**poolId** | **string** |  | 
**replicaId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiPutNodePoolReplicaShareRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



 **allowedHosts** | **[]string** | NQNs of hosts allowed to connect to this replica | 

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


## PutPoolReplica

> Replica PutPoolReplica(ctx, poolId, replicaId).CreateReplicaBody(createReplicaBody).Execute()



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
	replicaId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	createReplicaBody := *openapiclient.NewCreateReplicaBody(int64(123), false) // CreateReplicaBody | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ReplicasAPI.PutPoolReplica(context.Background(), poolId, replicaId).CreateReplicaBody(createReplicaBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ReplicasAPI.PutPoolReplica``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PutPoolReplica`: Replica
	fmt.Fprintf(os.Stdout, "Response from `ReplicasAPI.PutPoolReplica`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**poolId** | **string** |  | 
**replicaId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiPutPoolReplicaRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **createReplicaBody** | [**CreateReplicaBody**](CreateReplicaBody.md) |  | 

### Return type

[**Replica**](Replica.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PutPoolReplicaShare

> string PutPoolReplicaShare(ctx, poolId, replicaId).AllowedHosts(allowedHosts).Execute()



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
	replicaId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	allowedHosts := []string{"nqn.2014-08.org.nvmexpress:uuid:804b1e8c-b42d-4d15-92b4-7c4e4d0f507"} // []string | NQNs of hosts allowed to connect to this replica (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ReplicasAPI.PutPoolReplicaShare(context.Background(), poolId, replicaId).AllowedHosts(allowedHosts).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ReplicasAPI.PutPoolReplicaShare``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PutPoolReplicaShare`: string
	fmt.Fprintf(os.Stdout, "Response from `ReplicasAPI.PutPoolReplicaShare`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**poolId** | **string** |  | 
**replicaId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiPutPoolReplicaShareRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **allowedHosts** | **[]string** | NQNs of hosts allowed to connect to this replica | 

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

