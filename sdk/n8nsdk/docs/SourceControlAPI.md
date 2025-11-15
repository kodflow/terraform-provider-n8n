# \SourceControlAPI

All URIs are relative to */api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**SourceControlPullPost**](SourceControlAPI.md#SourceControlPullPost) | **Post** /source-control/pull | Pull changes from the remote repository



## SourceControlPullPost

> ImportResult SourceControlPullPost(ctx).Pull(pull).Execute()

Pull changes from the remote repository



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/n8nsdk"
)

func main() {
	pull := *openapiclient.NewPull() // Pull | Pull options

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SourceControlAPI.SourceControlPullPost(context.Background()).Pull(pull).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SourceControlAPI.SourceControlPullPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `SourceControlPullPost`: ImportResult
	fmt.Fprintf(os.Stdout, "Response from `SourceControlAPI.SourceControlPullPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSourceControlPullPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **pull** | [**Pull**](Pull.md) | Pull options | 

### Return type

[**ImportResult**](ImportResult.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

