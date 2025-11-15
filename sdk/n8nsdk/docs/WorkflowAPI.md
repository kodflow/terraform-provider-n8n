# \WorkflowAPI

All URIs are relative to */api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**WorkflowsGet**](WorkflowAPI.md#WorkflowsGet) | **Get** /workflows | Retrieve all workflows
[**WorkflowsIdActivatePost**](WorkflowAPI.md#WorkflowsIdActivatePost) | **Post** /workflows/{id}/activate | Activate a workflow
[**WorkflowsIdDeactivatePost**](WorkflowAPI.md#WorkflowsIdDeactivatePost) | **Post** /workflows/{id}/deactivate | Deactivate a workflow
[**WorkflowsIdDelete**](WorkflowAPI.md#WorkflowsIdDelete) | **Delete** /workflows/{id} | Delete a workflow
[**WorkflowsIdGet**](WorkflowAPI.md#WorkflowsIdGet) | **Get** /workflows/{id} | Retrieves a workflow
[**WorkflowsIdPut**](WorkflowAPI.md#WorkflowsIdPut) | **Put** /workflows/{id} | Update a workflow
[**WorkflowsIdTagsGet**](WorkflowAPI.md#WorkflowsIdTagsGet) | **Get** /workflows/{id}/tags | Get workflow tags
[**WorkflowsIdTagsPut**](WorkflowAPI.md#WorkflowsIdTagsPut) | **Put** /workflows/{id}/tags | Update tags of a workflow
[**WorkflowsIdTransferPut**](WorkflowAPI.md#WorkflowsIdTransferPut) | **Put** /workflows/{id}/transfer | Transfer a workflow to another project.
[**WorkflowsPost**](WorkflowAPI.md#WorkflowsPost) | **Post** /workflows | Create a workflow



## WorkflowsGet

> WorkflowList WorkflowsGet(ctx).Active(active).Tags(tags).Name(name).ProjectId(projectId).ExcludePinnedData(excludePinnedData).Limit(limit).Cursor(cursor).Execute()

Retrieve all workflows



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
	active := true // bool |  (optional)
	tags := "test,production" // string |  (optional)
	name := "My Workflow" // string |  (optional)
	projectId := "VmwOO9HeTEj20kxM" // string |  (optional)
	excludePinnedData := true // bool | Set this to avoid retrieving pinned data (optional)
	limit := float32(100) // float32 | The maximum number of items to return. (optional) (default to 100)
	cursor := "cursor_example" // string | Paginate by setting the cursor parameter to the nextCursor attribute returned by the previous request's response. Default value fetches the first \"page\" of the collection. See pagination for more detail. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.WorkflowAPI.WorkflowsGet(context.Background()).Active(active).Tags(tags).Name(name).ProjectId(projectId).ExcludePinnedData(excludePinnedData).Limit(limit).Cursor(cursor).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WorkflowAPI.WorkflowsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `WorkflowsGet`: WorkflowList
	fmt.Fprintf(os.Stdout, "Response from `WorkflowAPI.WorkflowsGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiWorkflowsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **active** | **bool** |  | 
 **tags** | **string** |  | 
 **name** | **string** |  | 
 **projectId** | **string** |  | 
 **excludePinnedData** | **bool** | Set this to avoid retrieving pinned data | 
 **limit** | **float32** | The maximum number of items to return. | [default to 100]
 **cursor** | **string** | Paginate by setting the cursor parameter to the nextCursor attribute returned by the previous request&#39;s response. Default value fetches the first \&quot;page\&quot; of the collection. See pagination for more detail. | 

### Return type

[**WorkflowList**](WorkflowList.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## WorkflowsIdActivatePost

> Workflow WorkflowsIdActivatePost(ctx, id).Execute()

Activate a workflow



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
	id := "id_example" // string | The ID of the workflow.

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.WorkflowAPI.WorkflowsIdActivatePost(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WorkflowAPI.WorkflowsIdActivatePost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `WorkflowsIdActivatePost`: Workflow
	fmt.Fprintf(os.Stdout, "Response from `WorkflowAPI.WorkflowsIdActivatePost`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | The ID of the workflow. | 

### Other Parameters

Other parameters are passed through a pointer to a apiWorkflowsIdActivatePostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**Workflow**](Workflow.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## WorkflowsIdDeactivatePost

> Workflow WorkflowsIdDeactivatePost(ctx, id).Execute()

Deactivate a workflow



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
	id := "id_example" // string | The ID of the workflow.

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.WorkflowAPI.WorkflowsIdDeactivatePost(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WorkflowAPI.WorkflowsIdDeactivatePost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `WorkflowsIdDeactivatePost`: Workflow
	fmt.Fprintf(os.Stdout, "Response from `WorkflowAPI.WorkflowsIdDeactivatePost`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | The ID of the workflow. | 

### Other Parameters

Other parameters are passed through a pointer to a apiWorkflowsIdDeactivatePostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**Workflow**](Workflow.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## WorkflowsIdDelete

> Workflow WorkflowsIdDelete(ctx, id).Execute()

Delete a workflow



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
	id := "id_example" // string | The ID of the workflow.

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.WorkflowAPI.WorkflowsIdDelete(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WorkflowAPI.WorkflowsIdDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `WorkflowsIdDelete`: Workflow
	fmt.Fprintf(os.Stdout, "Response from `WorkflowAPI.WorkflowsIdDelete`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | The ID of the workflow. | 

### Other Parameters

Other parameters are passed through a pointer to a apiWorkflowsIdDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**Workflow**](Workflow.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## WorkflowsIdGet

> Workflow WorkflowsIdGet(ctx, id).ExcludePinnedData(excludePinnedData).Execute()

Retrieves a workflow



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
	id := "id_example" // string | The ID of the workflow.
	excludePinnedData := true // bool | Set this to avoid retrieving pinned data (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.WorkflowAPI.WorkflowsIdGet(context.Background(), id).ExcludePinnedData(excludePinnedData).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WorkflowAPI.WorkflowsIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `WorkflowsIdGet`: Workflow
	fmt.Fprintf(os.Stdout, "Response from `WorkflowAPI.WorkflowsIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | The ID of the workflow. | 

### Other Parameters

Other parameters are passed through a pointer to a apiWorkflowsIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **excludePinnedData** | **bool** | Set this to avoid retrieving pinned data | 

### Return type

[**Workflow**](Workflow.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## WorkflowsIdPut

> Workflow WorkflowsIdPut(ctx, id).Workflow(workflow).Execute()

Update a workflow



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
	id := "id_example" // string | The ID of the workflow.
	workflow := *openapiclient.NewWorkflow("Workflow 1", []openapiclient.Node{*openapiclient.NewNode()}, map[string]interface{}({"Jira":{"main":[[{"node":"Jira","type":"main","index":0}]]}}), *openapiclient.NewWorkflowSettings()) // Workflow | Updated workflow object.

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.WorkflowAPI.WorkflowsIdPut(context.Background(), id).Workflow(workflow).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WorkflowAPI.WorkflowsIdPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `WorkflowsIdPut`: Workflow
	fmt.Fprintf(os.Stdout, "Response from `WorkflowAPI.WorkflowsIdPut`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | The ID of the workflow. | 

### Other Parameters

Other parameters are passed through a pointer to a apiWorkflowsIdPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **workflow** | [**Workflow**](Workflow.md) | Updated workflow object. | 

### Return type

[**Workflow**](Workflow.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## WorkflowsIdTagsGet

> []Tag WorkflowsIdTagsGet(ctx, id).Execute()

Get workflow tags



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
	id := "id_example" // string | The ID of the workflow.

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.WorkflowAPI.WorkflowsIdTagsGet(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WorkflowAPI.WorkflowsIdTagsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `WorkflowsIdTagsGet`: []Tag
	fmt.Fprintf(os.Stdout, "Response from `WorkflowAPI.WorkflowsIdTagsGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | The ID of the workflow. | 

### Other Parameters

Other parameters are passed through a pointer to a apiWorkflowsIdTagsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**[]Tag**](Tag.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## WorkflowsIdTagsPut

> []Tag WorkflowsIdTagsPut(ctx, id).TagIdsInner(tagIdsInner).Execute()

Update tags of a workflow



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
	id := "id_example" // string | The ID of the workflow.
	tagIdsInner := []openapiclient.TagIdsInner{*openapiclient.NewTagIdsInner("2tUt1wbLX592XDdX")} // []TagIdsInner | List of tags

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.WorkflowAPI.WorkflowsIdTagsPut(context.Background(), id).TagIdsInner(tagIdsInner).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WorkflowAPI.WorkflowsIdTagsPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `WorkflowsIdTagsPut`: []Tag
	fmt.Fprintf(os.Stdout, "Response from `WorkflowAPI.WorkflowsIdTagsPut`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | The ID of the workflow. | 

### Other Parameters

Other parameters are passed through a pointer to a apiWorkflowsIdTagsPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **tagIdsInner** | [**[]TagIdsInner**](TagIdsInner.md) | List of tags | 

### Return type

[**[]Tag**](Tag.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## WorkflowsIdTransferPut

> WorkflowsIdTransferPut(ctx, id).WorkflowsIdTransferPutRequest(workflowsIdTransferPutRequest).Execute()

Transfer a workflow to another project.



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
	id := "id_example" // string | The ID of the workflow.
	workflowsIdTransferPutRequest := *openapiclient.NewWorkflowsIdTransferPutRequest("DestinationProjectId_example") // WorkflowsIdTransferPutRequest | Destination project information for the workflow transfer.

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.WorkflowAPI.WorkflowsIdTransferPut(context.Background(), id).WorkflowsIdTransferPutRequest(workflowsIdTransferPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WorkflowAPI.WorkflowsIdTransferPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | The ID of the workflow. | 

### Other Parameters

Other parameters are passed through a pointer to a apiWorkflowsIdTransferPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **workflowsIdTransferPutRequest** | [**WorkflowsIdTransferPutRequest**](WorkflowsIdTransferPutRequest.md) | Destination project information for the workflow transfer. | 

### Return type

 (empty response body)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## WorkflowsPost

> Workflow WorkflowsPost(ctx).Workflow(workflow).Execute()

Create a workflow



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
	workflow := *openapiclient.NewWorkflow("Workflow 1", []openapiclient.Node{*openapiclient.NewNode()}, map[string]interface{}({"Jira":{"main":[[{"node":"Jira","type":"main","index":0}]]}}), *openapiclient.NewWorkflowSettings()) // Workflow | Created workflow object.

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.WorkflowAPI.WorkflowsPost(context.Background()).Workflow(workflow).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WorkflowAPI.WorkflowsPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `WorkflowsPost`: Workflow
	fmt.Fprintf(os.Stdout, "Response from `WorkflowAPI.WorkflowsPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiWorkflowsPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **workflow** | [**Workflow**](Workflow.md) | Created workflow object. | 

### Return type

[**Workflow**](Workflow.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

