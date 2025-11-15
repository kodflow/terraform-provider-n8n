# \ExecutionAPI

All URIs are relative to */api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ExecutionsGet**](ExecutionAPI.md#ExecutionsGet) | **Get** /executions | Retrieve all executions
[**ExecutionsIdDelete**](ExecutionAPI.md#ExecutionsIdDelete) | **Delete** /executions/{id} | Delete an execution
[**ExecutionsIdGet**](ExecutionAPI.md#ExecutionsIdGet) | **Get** /executions/{id} | Retrieve an execution
[**ExecutionsIdRetryPost**](ExecutionAPI.md#ExecutionsIdRetryPost) | **Post** /executions/{id}/retry | Retry an execution



## ExecutionsGet

> ExecutionList ExecutionsGet(ctx).IncludeData(includeData).Status(status).WorkflowId(workflowId).ProjectId(projectId).Limit(limit).Cursor(cursor).Execute()

Retrieve all executions



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
	includeData := true // bool | Whether or not to include the execution's detailed data. (optional)
	status := "status_example" // string | Status to filter the executions by. (optional)
	workflowId := "1000" // string | Workflow to filter the executions by. (optional)
	projectId := "VmwOO9HeTEj20kxM" // string |  (optional)
	limit := float32(100) // float32 | The maximum number of items to return. (optional) (default to 100)
	cursor := "cursor_example" // string | Paginate by setting the cursor parameter to the nextCursor attribute returned by the previous request's response. Default value fetches the first \"page\" of the collection. See pagination for more detail. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ExecutionAPI.ExecutionsGet(context.Background()).IncludeData(includeData).Status(status).WorkflowId(workflowId).ProjectId(projectId).Limit(limit).Cursor(cursor).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ExecutionAPI.ExecutionsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ExecutionsGet`: ExecutionList
	fmt.Fprintf(os.Stdout, "Response from `ExecutionAPI.ExecutionsGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiExecutionsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **includeData** | **bool** | Whether or not to include the execution&#39;s detailed data. | 
 **status** | **string** | Status to filter the executions by. | 
 **workflowId** | **string** | Workflow to filter the executions by. | 
 **projectId** | **string** |  | 
 **limit** | **float32** | The maximum number of items to return. | [default to 100]
 **cursor** | **string** | Paginate by setting the cursor parameter to the nextCursor attribute returned by the previous request&#39;s response. Default value fetches the first \&quot;page\&quot; of the collection. See pagination for more detail. | 

### Return type

[**ExecutionList**](ExecutionList.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ExecutionsIdDelete

> Execution ExecutionsIdDelete(ctx, id).Execute()

Delete an execution



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
	id := float32(8.14) // float32 | The ID of the execution.

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ExecutionAPI.ExecutionsIdDelete(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ExecutionAPI.ExecutionsIdDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ExecutionsIdDelete`: Execution
	fmt.Fprintf(os.Stdout, "Response from `ExecutionAPI.ExecutionsIdDelete`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **float32** | The ID of the execution. | 

### Other Parameters

Other parameters are passed through a pointer to a apiExecutionsIdDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**Execution**](Execution.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ExecutionsIdGet

> Execution ExecutionsIdGet(ctx, id).IncludeData(includeData).Execute()

Retrieve an execution



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
	id := float32(8.14) // float32 | The ID of the execution.
	includeData := true // bool | Whether or not to include the execution's detailed data. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ExecutionAPI.ExecutionsIdGet(context.Background(), id).IncludeData(includeData).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ExecutionAPI.ExecutionsIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ExecutionsIdGet`: Execution
	fmt.Fprintf(os.Stdout, "Response from `ExecutionAPI.ExecutionsIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **float32** | The ID of the execution. | 

### Other Parameters

Other parameters are passed through a pointer to a apiExecutionsIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **includeData** | **bool** | Whether or not to include the execution&#39;s detailed data. | 

### Return type

[**Execution**](Execution.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ExecutionsIdRetryPost

> Execution ExecutionsIdRetryPost(ctx, id).ExecutionsIdRetryPostRequest(executionsIdRetryPostRequest).Execute()

Retry an execution



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
	id := float32(8.14) // float32 | The ID of the execution.
	executionsIdRetryPostRequest := *openapiclient.NewExecutionsIdRetryPostRequest() // ExecutionsIdRetryPostRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ExecutionAPI.ExecutionsIdRetryPost(context.Background(), id).ExecutionsIdRetryPostRequest(executionsIdRetryPostRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ExecutionAPI.ExecutionsIdRetryPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ExecutionsIdRetryPost`: Execution
	fmt.Fprintf(os.Stdout, "Response from `ExecutionAPI.ExecutionsIdRetryPost`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **float32** | The ID of the execution. | 

### Other Parameters

Other parameters are passed through a pointer to a apiExecutionsIdRetryPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **executionsIdRetryPostRequest** | [**ExecutionsIdRetryPostRequest**](ExecutionsIdRetryPostRequest.md) |  | 

### Return type

[**Execution**](Execution.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

