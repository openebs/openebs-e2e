# \SpecsAPI

All URIs are relative to */v0*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetSpecs**](SpecsAPI.md#GetSpecs) | **Get** /specs | 



## GetSpecs

> Specs GetSpecs(ctx).Execute()



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
	resp, r, err := apiClient.SpecsAPI.GetSpecs(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SpecsAPI.GetSpecs``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetSpecs`: Specs
	fmt.Fprintf(os.Stdout, "Response from `SpecsAPI.GetSpecs`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetSpecsRequest struct via the builder pattern


### Return type

[**Specs**](Specs.md)

### Authorization

[JWT](../README.md#JWT)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

