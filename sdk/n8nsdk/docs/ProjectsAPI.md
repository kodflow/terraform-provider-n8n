# \ProjectsAPI

All URIs are relative to */api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ProjectsGet**](ProjectsAPI.md#ProjectsGet) | **Get** /projects | Retrieve projects
[**ProjectsPost**](ProjectsAPI.md#ProjectsPost) | **Post** /projects | Create a project
[**ProjectsProjectIdDelete**](ProjectsAPI.md#ProjectsProjectIdDelete) | **Delete** /projects/{projectId} | Delete a project
[**ProjectsProjectIdPut**](ProjectsAPI.md#ProjectsProjectIdPut) | **Put** /projects/{projectId} | Update a project
[**ProjectsProjectIdUsersPost**](ProjectsAPI.md#ProjectsProjectIdUsersPost) | **Post** /projects/{projectId}/users | Add one or more users to a project
[**ProjectsProjectIdUsersUserIdDelete**](ProjectsAPI.md#ProjectsProjectIdUsersUserIdDelete) | **Delete** /projects/{projectId}/users/{userId} | Delete a user from a project
[**ProjectsProjectIdUsersUserIdPatch**](ProjectsAPI.md#ProjectsProjectIdUsersUserIdPatch) | **Patch** /projects/{projectId}/users/{userId} | Change a user&#39;s role in a project



## ProjectsGet

> ProjectList ProjectsGet(ctx).Limit(limit).Cursor(cursor).Execute()

Retrieve projects



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

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ProjectsAPI.ProjectsGet(context.Background()).Limit(limit).Cursor(cursor).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ProjectsAPI.ProjectsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ProjectsGet`: ProjectList
	fmt.Fprintf(os.Stdout, "Response from `ProjectsAPI.ProjectsGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiProjectsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **limit** | **float32** | The maximum number of items to return. | [default to 100]
 **cursor** | **string** | Paginate by setting the cursor parameter to the nextCursor attribute returned by the previous request&#39;s response. Default value fetches the first \&quot;page\&quot; of the collection. See pagination for more detail. | 

### Return type

[**ProjectList**](ProjectList.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ProjectsPost

> ProjectsPost(ctx).Project(project).Execute()

Create a project



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
	project := *openapiclient.NewProject("Name_example") // Project | Payload for project to create.

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ProjectsAPI.ProjectsPost(context.Background()).Project(project).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ProjectsAPI.ProjectsPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiProjectsPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **project** | [**Project**](Project.md) | Payload for project to create. | 

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


## ProjectsProjectIdDelete

> ProjectsProjectIdDelete(ctx, projectId).Execute()

Delete a project



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
	projectId := "projectId_example" // string | The ID of the project.

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ProjectsAPI.ProjectsProjectIdDelete(context.Background(), projectId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ProjectsAPI.ProjectsProjectIdDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**projectId** | **string** | The ID of the project. | 

### Other Parameters

Other parameters are passed through a pointer to a apiProjectsProjectIdDeleteRequest struct via the builder pattern


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


## ProjectsProjectIdPut

> ProjectsProjectIdPut(ctx, projectId).Project(project).Execute()

Update a project



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
	projectId := "projectId_example" // string | The ID of the project.
	project := *openapiclient.NewProject("Name_example") // Project | Updated project object.

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ProjectsAPI.ProjectsProjectIdPut(context.Background(), projectId).Project(project).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ProjectsAPI.ProjectsProjectIdPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**projectId** | **string** | The ID of the project. | 

### Other Parameters

Other parameters are passed through a pointer to a apiProjectsProjectIdPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **project** | [**Project**](Project.md) | Updated project object. | 

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


## ProjectsProjectIdUsersPost

> ProjectsProjectIdUsersPost(ctx, projectId).ProjectsProjectIdUsersPostRequest(projectsProjectIdUsersPostRequest).Execute()

Add one or more users to a project



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
	projectId := "projectId_example" // string | The ID of the project.
	projectsProjectIdUsersPostRequest := *openapiclient.NewProjectsProjectIdUsersPostRequest([]openapiclient.ProjectsProjectIdUsersPostRequestRelationsInner{*openapiclient.NewProjectsProjectIdUsersPostRequestRelationsInner("91765f0d-3b29-45df-adb9-35b23937eb92", "project:viewer")}) // ProjectsProjectIdUsersPostRequest | Payload containing an array of one or more users to add to the project. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ProjectsAPI.ProjectsProjectIdUsersPost(context.Background(), projectId).ProjectsProjectIdUsersPostRequest(projectsProjectIdUsersPostRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ProjectsAPI.ProjectsProjectIdUsersPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**projectId** | **string** | The ID of the project. | 

### Other Parameters

Other parameters are passed through a pointer to a apiProjectsProjectIdUsersPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **projectsProjectIdUsersPostRequest** | [**ProjectsProjectIdUsersPostRequest**](ProjectsProjectIdUsersPostRequest.md) | Payload containing an array of one or more users to add to the project. | 

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


## ProjectsProjectIdUsersUserIdDelete

> ProjectsProjectIdUsersUserIdDelete(ctx, projectId, userId).Execute()

Delete a user from a project



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
	projectId := "projectId_example" // string | The ID of the project.
	userId := "userId_example" // string | The ID of the user.

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ProjectsAPI.ProjectsProjectIdUsersUserIdDelete(context.Background(), projectId, userId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ProjectsAPI.ProjectsProjectIdUsersUserIdDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**projectId** | **string** | The ID of the project. | 
**userId** | **string** | The ID of the user. | 

### Other Parameters

Other parameters are passed through a pointer to a apiProjectsProjectIdUsersUserIdDeleteRequest struct via the builder pattern


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


## ProjectsProjectIdUsersUserIdPatch

> ProjectsProjectIdUsersUserIdPatch(ctx, projectId, userId).ProjectsProjectIdUsersUserIdPatchRequest(projectsProjectIdUsersUserIdPatchRequest).Execute()

Change a user's role in a project



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
	projectId := "projectId_example" // string | The ID of the project.
	userId := "userId_example" // string | The ID of the user.
	projectsProjectIdUsersUserIdPatchRequest := *openapiclient.NewProjectsProjectIdUsersUserIdPatchRequest("project:viewer") // ProjectsProjectIdUsersUserIdPatchRequest | Payload containing the new role to assign to the project user. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ProjectsAPI.ProjectsProjectIdUsersUserIdPatch(context.Background(), projectId, userId).ProjectsProjectIdUsersUserIdPatchRequest(projectsProjectIdUsersUserIdPatchRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ProjectsAPI.ProjectsProjectIdUsersUserIdPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**projectId** | **string** | The ID of the project. | 
**userId** | **string** | The ID of the user. | 

### Other Parameters

Other parameters are passed through a pointer to a apiProjectsProjectIdUsersUserIdPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **projectsProjectIdUsersUserIdPatchRequest** | [**ProjectsProjectIdUsersUserIdPatchRequest**](ProjectsProjectIdUsersUserIdPatchRequest.md) | Payload containing the new role to assign to the project user. | 

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

