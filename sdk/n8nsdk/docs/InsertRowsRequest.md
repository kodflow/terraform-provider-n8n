# InsertRowsRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | **[]map[string]interface{}** | Array of rows to insert. Each row is an object with column names as keys. | 
**ReturnType** | Pointer to **string** | - count: Return only the number of rows inserted - id: Return an array of inserted row IDs - all: Return the full row data for all inserted rows  | [optional] [default to "count"]

## Methods

### NewInsertRowsRequest

`func NewInsertRowsRequest(data []map[string]interface{}, ) *InsertRowsRequest`

NewInsertRowsRequest instantiates a new InsertRowsRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInsertRowsRequestWithDefaults

`func NewInsertRowsRequestWithDefaults() *InsertRowsRequest`

NewInsertRowsRequestWithDefaults instantiates a new InsertRowsRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *InsertRowsRequest) GetData() []map[string]interface{}`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *InsertRowsRequest) GetDataOk() (*[]map[string]interface{}, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *InsertRowsRequest) SetData(v []map[string]interface{})`

SetData sets Data field to given value.


### GetReturnType

`func (o *InsertRowsRequest) GetReturnType() string`

GetReturnType returns the ReturnType field if non-nil, zero value otherwise.

### GetReturnTypeOk

`func (o *InsertRowsRequest) GetReturnTypeOk() (*string, bool)`

GetReturnTypeOk returns a tuple with the ReturnType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReturnType

`func (o *InsertRowsRequest) SetReturnType(v string)`

SetReturnType sets ReturnType field to given value.

### HasReturnType

`func (o *InsertRowsRequest) HasReturnType() bool`

HasReturnType returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


