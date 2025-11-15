# \VariablesAPI

All URIs are relative to */api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**VariablesGet**](VariablesAPI.md#VariablesGet) | **Get** /variables | Retrieve variables
[**VariablesIdDelete**](VariablesAPI.md#VariablesIdDelete) | **Delete** /variables/{id} | Delete a variable
[**VariablesIdPut**](VariablesAPI.md#VariablesIdPut) | **Put** /variables/{id} | Update a variable
[**VariablesPost**](VariablesAPI.md#VariablesPost) | **Post** /variables | Create a variable



## VariablesGet

> VariableList VariablesGet(ctx).Limit(limit).Cursor(cursor).ProjectId(projectId).State(state).Execute()

Retrieve variables



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
	limit := float32(100) // float32 | The maximum number of items to return. (optional) (default to 100)
	cursor := "cursor_example" // string | Paginate by setting the cursor parameter to the nextCursor attribute returned by the previous request's response. Default value fetches the first \"page\" of the collection. See pagination for more detail. (optional)
	projectId := "VmwOO9HeTEj20kxM" // string |  (optional)
	state := "state_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.VariablesAPI.VariablesGet(context.Background()).Limit(limit).Cursor(cursor).ProjectId(projectId).State(state).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VariablesAPI.VariablesGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `VariablesGet`: VariableList
	fmt.Fprintf(os.Stdout, "Response from `VariablesAPI.VariablesGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiVariablesGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **limit** | **float32** | The maximum number of items to return. | [default to 100]
 **cursor** | **string** | Paginate by setting the cursor parameter to the nextCursor attribute returned by the previous request&#39;s response. Default value fetches the first \&quot;page\&quot; of the collection. See pagination for more detail. | 
 **projectId** | **string** |  | 
 **state** | **string** |  | 

### Return type

[**VariableList**](VariableList.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## VariablesIdDelete

> VariablesIdDelete(ctx, id).Execute()

Delete a variable



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
	id := "id_example" // string | The ID of the variable.

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.VariablesAPI.VariablesIdDelete(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VariablesAPI.VariablesIdDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | The ID of the variable. | 

### Other Parameters

Other parameters are passed through a pointer to a apiVariablesIdDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## VariablesIdPut

> VariablesIdPut(ctx, id).VariableCreate(variableCreate).Execute()

Update a variable



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
	id := "id_example" // string | The ID of the variable.
	variableCreate := *openapiclient.NewVariableCreate("Key_example", "test") // VariableCreate | Payload for variable to update.

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.VariablesAPI.VariablesIdPut(context.Background(), id).VariableCreate(variableCreate).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VariablesAPI.VariablesIdPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | The ID of the variable. | 

### Other Parameters

Other parameters are passed through a pointer to a apiVariablesIdPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **variableCreate** | [**VariableCreate**](VariableCreate.md) | Payload for variable to update. | 

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


## VariablesPost

> VariablesPost(ctx).VariableCreate(variableCreate).Execute()

Create a variable



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
	variableCreate := *openapiclient.NewVariableCreate("Key_example", "test") // VariableCreate | Payload for variable to create.

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.VariablesAPI.VariablesPost(context.Background()).VariableCreate(variableCreate).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VariablesAPI.VariablesPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiVariablesPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **variableCreate** | [**VariableCreate**](VariableCreate.md) | Payload for variable to create. | 

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

