# \CredentialAPI

All URIs are relative to */api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CredentialsIdTransferPut**](CredentialAPI.md#CredentialsIdTransferPut) | **Put** /credentials/{id}/transfer | Transfer a credential to another project.
[**CredentialsPost**](CredentialAPI.md#CredentialsPost) | **Post** /credentials | Create a credential
[**CredentialsSchemaCredentialTypeNameGet**](CredentialAPI.md#CredentialsSchemaCredentialTypeNameGet) | **Get** /credentials/schema/{credentialTypeName} | Show credential data schema
[**DeleteCredential**](CredentialAPI.md#DeleteCredential) | **Delete** /credentials/{id} | Delete credential by ID



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


## CredentialsPost

> CreateCredentialResponse CredentialsPost(ctx).Credential(credential).Execute()

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
	credential := *openapiclient.NewCredential("Joe's Github Credentials", "github", map[string]interface{}({"token":"ada612vad6fa5df4adf5a5dsf4389adsf76da7s"})) // Credential | Credential to be created.

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CredentialAPI.CredentialsPost(context.Background()).Credential(credential).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CredentialAPI.CredentialsPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CredentialsPost`: CreateCredentialResponse
	fmt.Fprintf(os.Stdout, "Response from `CredentialAPI.CredentialsPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCredentialsPostRequest struct via the builder pattern


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

