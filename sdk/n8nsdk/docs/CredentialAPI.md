# \CredentialAPI

All URIs are relative to */api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateCredential**](CredentialAPI.md#CreateCredential) | **Post** /credentials | Create a credential
[**CredentialsIdTransferPut**](CredentialAPI.md#CredentialsIdTransferPut) | **Put** /credentials/{id}/transfer | Transfer a credential to another project.
[**CredentialsSchemaCredentialTypeNameGet**](CredentialAPI.md#CredentialsSchemaCredentialTypeNameGet) | **Get** /credentials/schema/{credentialTypeName} | Show credential data schema
[**DeleteCredential**](CredentialAPI.md#DeleteCredential) | **Delete** /credentials/{id} | Delete credential by ID
[**GetCredentials**](CredentialAPI.md#GetCredentials) | **Get** /credentials | List credentials
[**UpdateCredential**](CredentialAPI.md#UpdateCredential) | **Patch** /credentials/{id} | Update credential by ID



## CreateCredential

> CreateCredentialResponse CreateCredential(ctx).Credential(credential).Execute()

Create a credential



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
	credential := *openapiclient.NewCredential("Joe's Github Credentials", "githubApi", map[string]interface{}({"accessToken":"ada612vad6fa5df4adf5a5dsf4389adsf76da7s"})) // Credential | Credential to be created.

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CredentialAPI.CreateCredential(context.Background()).Credential(credential).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CredentialAPI.CreateCredential``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateCredential`: CreateCredentialResponse
	fmt.Fprintf(os.Stdout, "Response from `CredentialAPI.CreateCredential`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateCredentialRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **credential** | [**Credential**](Credential.md) | Credential to be created. | 

### Return type

[**CreateCredentialResponse**](CreateCredentialResponse.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CredentialsIdTransferPut

> CredentialsIdTransferPut(ctx, id).CredentialsIdTransferPutRequest(credentialsIdTransferPutRequest).Execute()

Transfer a credential to another project.



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
	id := "id_example" // string | The ID of the credential.
	credentialsIdTransferPutRequest := *openapiclient.NewCredentialsIdTransferPutRequest("DestinationProjectId_example") // CredentialsIdTransferPutRequest | Destination project for the credential transfer.

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.CredentialAPI.CredentialsIdTransferPut(context.Background(), id).CredentialsIdTransferPutRequest(credentialsIdTransferPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CredentialAPI.CredentialsIdTransferPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | The ID of the credential. | 

### Other Parameters

Other parameters are passed through a pointer to a apiCredentialsIdTransferPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **credentialsIdTransferPutRequest** | [**CredentialsIdTransferPutRequest**](CredentialsIdTransferPutRequest.md) | Destination project for the credential transfer. | 

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


## CredentialsSchemaCredentialTypeNameGet

> map[string]interface{} CredentialsSchemaCredentialTypeNameGet(ctx, credentialTypeName).Execute()

Show credential data schema

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
	credentialTypeName := "credentialTypeName_example" // string | The credential type name that you want to get the schema for

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CredentialAPI.CredentialsSchemaCredentialTypeNameGet(context.Background(), credentialTypeName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CredentialAPI.CredentialsSchemaCredentialTypeNameGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CredentialsSchemaCredentialTypeNameGet`: map[string]interface{}
	fmt.Fprintf(os.Stdout, "Response from `CredentialAPI.CredentialsSchemaCredentialTypeNameGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**credentialTypeName** | **string** | The credential type name that you want to get the schema for | 

### Other Parameters

Other parameters are passed through a pointer to a apiCredentialsSchemaCredentialTypeNameGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

**map[string]interface{}**

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteCredential

> Credential DeleteCredential(ctx, id).Execute()

Delete credential by ID



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
	id := "id_example" // string | The credential ID that needs to be deleted

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CredentialAPI.DeleteCredential(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CredentialAPI.DeleteCredential``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteCredential`: Credential
	fmt.Fprintf(os.Stdout, "Response from `CredentialAPI.DeleteCredential`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | The credential ID that needs to be deleted | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteCredentialRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**Credential**](Credential.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetCredentials

> CredentialList GetCredentials(ctx).Limit(limit).Cursor(cursor).Execute()

List credentials



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
	resp, r, err := apiClient.CredentialAPI.GetCredentials(context.Background()).Limit(limit).Cursor(cursor).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CredentialAPI.GetCredentials``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetCredentials`: CredentialList
	fmt.Fprintf(os.Stdout, "Response from `CredentialAPI.GetCredentials`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGetCredentialsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **limit** | **float32** | The maximum number of items to return. | [default to 100]
 **cursor** | **string** | Paginate by setting the cursor parameter to the nextCursor attribute returned by the previous request&#39;s response. Default value fetches the first \&quot;page\&quot; of the collection. See pagination for more detail. | 

### Return type

[**CredentialList**](CredentialList.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateCredential

> CreateCredentialResponse UpdateCredential(ctx, id).UpdateCredentialRequest(updateCredentialRequest).Execute()

Update credential by ID



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
	id := "id_example" // string | The credential ID that needs to be updated
	updateCredentialRequest := *openapiclient.NewUpdateCredentialRequest() // UpdateCredentialRequest | Credential data to update. All fields are optional.

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CredentialAPI.UpdateCredential(context.Background(), id).UpdateCredentialRequest(updateCredentialRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CredentialAPI.UpdateCredential``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateCredential`: CreateCredentialResponse
	fmt.Fprintf(os.Stdout, "Response from `CredentialAPI.UpdateCredential`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | The credential ID that needs to be updated | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateCredentialRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateCredentialRequest** | [**UpdateCredentialRequest**](UpdateCredentialRequest.md) | Credential data to update. All fields are optional. | 

### Return type

[**CreateCredentialResponse**](CreateCredentialResponse.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

