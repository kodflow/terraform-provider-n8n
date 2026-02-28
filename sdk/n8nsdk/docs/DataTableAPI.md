# \DataTableAPI

All URIs are relative to */api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateDataTable**](DataTableAPI.md#CreateDataTable) | **Post** /data-tables | Create a new data table
[**DeleteDataTable**](DataTableAPI.md#DeleteDataTable) | **Delete** /data-tables/{dataTableId} | Delete a data table
[**DeleteDataTableRows**](DataTableAPI.md#DeleteDataTableRows) | **Delete** /data-tables/{dataTableId}/rows/delete | Delete rows from a data table
[**GetDataTable**](DataTableAPI.md#GetDataTable) | **Get** /data-tables/{dataTableId} | Get a data table
[**GetDataTableRows**](DataTableAPI.md#GetDataTableRows) | **Get** /data-tables/{dataTableId}/rows | Retrieve rows from a data table
[**InsertDataTableRows**](DataTableAPI.md#InsertDataTableRows) | **Post** /data-tables/{dataTableId}/rows | Insert rows into a data table
[**ListDataTables**](DataTableAPI.md#ListDataTables) | **Get** /data-tables | List all data tables
[**UpdateDataTable**](DataTableAPI.md#UpdateDataTable) | **Patch** /data-tables/{dataTableId} | Update a data table
[**UpdateDataTableRows**](DataTableAPI.md#UpdateDataTableRows) | **Patch** /data-tables/{dataTableId}/rows/update | Update rows in a data table
[**UpsertDataTableRow**](DataTableAPI.md#UpsertDataTableRow) | **Post** /data-tables/{dataTableId}/rows/upsert | Upsert a row in a data table



## CreateDataTable

> DataTable CreateDataTable(ctx).CreateDataTableRequest(createDataTableRequest).Execute()

Create a new data table



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
	createDataTableRequest := *openapiclient.NewCreateDataTableRequest("Name_example", []openapiclient.CreateDataTableRequestColumnsInner{*openapiclient.NewCreateDataTableRequestColumnsInner("Name_example", "Type_example")}) // CreateDataTableRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DataTableAPI.CreateDataTable(context.Background()).CreateDataTableRequest(createDataTableRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DataTableAPI.CreateDataTable``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateDataTable`: DataTable
	fmt.Fprintf(os.Stdout, "Response from `DataTableAPI.CreateDataTable`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateDataTableRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **createDataTableRequest** | [**CreateDataTableRequest**](CreateDataTableRequest.md) |  | 

### Return type

[**DataTable**](DataTable.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteDataTable

> DeleteDataTable(ctx, dataTableId).Execute()

Delete a data table



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
	dataTableId := "dataTableId_example" // string | The ID of the data table

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DataTableAPI.DeleteDataTable(context.Background(), dataTableId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DataTableAPI.DeleteDataTable``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**dataTableId** | **string** | The ID of the data table | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteDataTableRequest struct via the builder pattern


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


## DeleteDataTableRows

> DeleteDataTableRows200Response DeleteDataTableRows(ctx, dataTableId).Filter(filter).ReturnData(returnData).DryRun(dryRun).Execute()

Delete rows from a data table



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
	dataTableId := "dataTableId_example" // string | The ID of the data table
	filter := "{"type":"and","filters":[{"columnName":"status","condition":"eq","value":"archived"}]}" // string | JSON string of filter conditions. Required to prevent accidental deletion of all data.
	returnData := true // bool | If true, return the deleted rows; if false, return true on success (optional) (default to false)
	dryRun := true // bool | If true, preview which rows would be deleted without actually deleting them (optional) (default to false)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DataTableAPI.DeleteDataTableRows(context.Background(), dataTableId).Filter(filter).ReturnData(returnData).DryRun(dryRun).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DataTableAPI.DeleteDataTableRows``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteDataTableRows`: DeleteDataTableRows200Response
	fmt.Fprintf(os.Stdout, "Response from `DataTableAPI.DeleteDataTableRows`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**dataTableId** | **string** | The ID of the data table | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteDataTableRowsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **filter** | **string** | JSON string of filter conditions. Required to prevent accidental deletion of all data. | 
 **returnData** | **bool** | If true, return the deleted rows; if false, return true on success | [default to false]
 **dryRun** | **bool** | If true, preview which rows would be deleted without actually deleting them | [default to false]

### Return type

[**DeleteDataTableRows200Response**](DeleteDataTableRows200Response.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetDataTable

> DataTable GetDataTable(ctx, dataTableId).Execute()

Get a data table



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
	dataTableId := "dataTableId_example" // string | The ID of the data table

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DataTableAPI.GetDataTable(context.Background(), dataTableId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DataTableAPI.GetDataTable``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetDataTable`: DataTable
	fmt.Fprintf(os.Stdout, "Response from `DataTableAPI.GetDataTable`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**dataTableId** | **string** | The ID of the data table | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetDataTableRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**DataTable**](DataTable.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetDataTableRows

> DataTableRowList GetDataTableRows(ctx, dataTableId).Limit(limit).Cursor(cursor).Filter(filter).SortBy(sortBy).Search(search).Execute()

Retrieve rows from a data table



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
	dataTableId := "dataTableId_example" // string | The ID of the data table
	limit := float32(100) // float32 | The maximum number of items to return. (optional) (default to 100)
	cursor := "cursor_example" // string | Paginate by setting the cursor parameter to the nextCursor attribute returned by the previous request's response. Default value fetches the first \"page\" of the collection. See pagination for more detail. (optional)
	filter := "{"type":"and","filters":[{"columnName":"status","condition":"eq","value":"active"}]}" // string | JSON string of filter conditions (optional)
	sortBy := "createdAt:desc" // string | Sort format: columnName:asc or columnName:desc (optional)
	search := "search_example" // string | Search text across all string columns (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DataTableAPI.GetDataTableRows(context.Background(), dataTableId).Limit(limit).Cursor(cursor).Filter(filter).SortBy(sortBy).Search(search).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DataTableAPI.GetDataTableRows``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetDataTableRows`: DataTableRowList
	fmt.Fprintf(os.Stdout, "Response from `DataTableAPI.GetDataTableRows`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**dataTableId** | **string** | The ID of the data table | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetDataTableRowsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **limit** | **float32** | The maximum number of items to return. | [default to 100]
 **cursor** | **string** | Paginate by setting the cursor parameter to the nextCursor attribute returned by the previous request&#39;s response. Default value fetches the first \&quot;page\&quot; of the collection. See pagination for more detail. | 
 **filter** | **string** | JSON string of filter conditions | 
 **sortBy** | **string** | Sort format: columnName:asc or columnName:desc | 
 **search** | **string** | Search text across all string columns | 

### Return type

[**DataTableRowList**](DataTableRowList.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InsertDataTableRows

> InsertDataTableRows200Response InsertDataTableRows(ctx, dataTableId).InsertRowsRequest(insertRowsRequest).Execute()

Insert rows into a data table



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
	dataTableId := "dataTableId_example" // string | The ID of the data table
	insertRowsRequest := *openapiclient.NewInsertRowsRequest([]map[string]interface{}{map[string]interface{}{"key": interface{}(123)}}) // InsertRowsRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DataTableAPI.InsertDataTableRows(context.Background(), dataTableId).InsertRowsRequest(insertRowsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DataTableAPI.InsertDataTableRows``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `InsertDataTableRows`: InsertDataTableRows200Response
	fmt.Fprintf(os.Stdout, "Response from `DataTableAPI.InsertDataTableRows`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**dataTableId** | **string** | The ID of the data table | 

### Other Parameters

Other parameters are passed through a pointer to a apiInsertDataTableRowsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **insertRowsRequest** | [**InsertRowsRequest**](InsertRowsRequest.md) |  | 

### Return type

[**InsertDataTableRows200Response**](InsertDataTableRows200Response.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListDataTables

> DataTableList ListDataTables(ctx).Limit(limit).Cursor(cursor).Filter(filter).SortBy(sortBy).Execute()

List all data tables



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
	filter := "{"name":"my-table"}" // string | JSON string of filter conditions (optional)
	sortBy := "name:asc" // string | Sort format: field:asc or field:desc (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DataTableAPI.ListDataTables(context.Background()).Limit(limit).Cursor(cursor).Filter(filter).SortBy(sortBy).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DataTableAPI.ListDataTables``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListDataTables`: DataTableList
	fmt.Fprintf(os.Stdout, "Response from `DataTableAPI.ListDataTables`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListDataTablesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **limit** | **float32** | The maximum number of items to return. | [default to 100]
 **cursor** | **string** | Paginate by setting the cursor parameter to the nextCursor attribute returned by the previous request&#39;s response. Default value fetches the first \&quot;page\&quot; of the collection. See pagination for more detail. | 
 **filter** | **string** | JSON string of filter conditions | 
 **sortBy** | **string** | Sort format: field:asc or field:desc | 

### Return type

[**DataTableList**](DataTableList.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateDataTable

> DataTable UpdateDataTable(ctx, dataTableId).UpdateDataTableRequest(updateDataTableRequest).Execute()

Update a data table



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
	dataTableId := "dataTableId_example" // string | The ID of the data table
	updateDataTableRequest := *openapiclient.NewUpdateDataTableRequest("Name_example") // UpdateDataTableRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DataTableAPI.UpdateDataTable(context.Background(), dataTableId).UpdateDataTableRequest(updateDataTableRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DataTableAPI.UpdateDataTable``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateDataTable`: DataTable
	fmt.Fprintf(os.Stdout, "Response from `DataTableAPI.UpdateDataTable`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**dataTableId** | **string** | The ID of the data table | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateDataTableRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateDataTableRequest** | [**UpdateDataTableRequest**](UpdateDataTableRequest.md) |  | 

### Return type

[**DataTable**](DataTable.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateDataTableRows

> UpdateDataTableRows200Response UpdateDataTableRows(ctx, dataTableId).UpdateRowsRequest(updateRowsRequest).Execute()

Update rows in a data table



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
	dataTableId := "dataTableId_example" // string | The ID of the data table
	updateRowsRequest := *openapiclient.NewUpdateRowsRequest(*openapiclient.NewUpdateRowsRequestFilter([]openapiclient.UpdateRowsRequestFilterFiltersInner{*openapiclient.NewUpdateRowsRequestFilterFiltersInner("ColumnName_example", "Condition_example", interface{}(123))}), map[string]interface{}{"key": interface{}(123)}) // UpdateRowsRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DataTableAPI.UpdateDataTableRows(context.Background(), dataTableId).UpdateRowsRequest(updateRowsRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DataTableAPI.UpdateDataTableRows``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateDataTableRows`: UpdateDataTableRows200Response
	fmt.Fprintf(os.Stdout, "Response from `DataTableAPI.UpdateDataTableRows`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**dataTableId** | **string** | The ID of the data table | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateDataTableRowsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateRowsRequest** | [**UpdateRowsRequest**](UpdateRowsRequest.md) |  | 

### Return type

[**UpdateDataTableRows200Response**](UpdateDataTableRows200Response.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpsertDataTableRow

> UpsertDataTableRow200Response UpsertDataTableRow(ctx, dataTableId).UpsertRowRequest(upsertRowRequest).Execute()

Upsert a row in a data table



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
	dataTableId := "dataTableId_example" // string | The ID of the data table
	upsertRowRequest := *openapiclient.NewUpsertRowRequest(*openapiclient.NewUpsertRowRequestFilter([]openapiclient.UpdateRowsRequestFilterFiltersInner{*openapiclient.NewUpdateRowsRequestFilterFiltersInner("ColumnName_example", "Condition_example", interface{}(123))}), map[string]interface{}{"key": interface{}(123)}) // UpsertRowRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DataTableAPI.UpsertDataTableRow(context.Background(), dataTableId).UpsertRowRequest(upsertRowRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DataTableAPI.UpsertDataTableRow``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpsertDataTableRow`: UpsertDataTableRow200Response
	fmt.Fprintf(os.Stdout, "Response from `DataTableAPI.UpsertDataTableRow`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**dataTableId** | **string** | The ID of the data table | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpsertDataTableRowRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **upsertRowRequest** | [**UpsertRowRequest**](UpsertRowRequest.md) |  | 

### Return type

[**UpsertDataTableRow200Response**](UpsertDataTableRow200Response.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

