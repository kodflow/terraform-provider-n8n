# UpdateRowsRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Filter** | [**UpdateRowsRequestFilter**](UpdateRowsRequestFilter.md) |  | 
**Data** | **map[string]interface{}** | Column values to update | 
**ReturnData** | Pointer to **bool** | If true, return the updated rows; if false, return true on success | [optional] [default to false]
**DryRun** | Pointer to **bool** | If true, preview changes without persisting them | [optional] [default to false]

## Methods

### NewUpdateRowsRequest

`func NewUpdateRowsRequest(filter UpdateRowsRequestFilter, data map[string]interface{}, ) *UpdateRowsRequest`

NewUpdateRowsRequest instantiates a new UpdateRowsRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateRowsRequestWithDefaults

`func NewUpdateRowsRequestWithDefaults() *UpdateRowsRequest`

NewUpdateRowsRequestWithDefaults instantiates a new UpdateRowsRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetFilter

`func (o *UpdateRowsRequest) GetFilter() UpdateRowsRequestFilter`

GetFilter returns the Filter field if non-nil, zero value otherwise.

### GetFilterOk

`func (o *UpdateRowsRequest) GetFilterOk() (*UpdateRowsRequestFilter, bool)`

GetFilterOk returns a tuple with the Filter field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFilter

`func (o *UpdateRowsRequest) SetFilter(v UpdateRowsRequestFilter)`

SetFilter sets Filter field to given value.


### GetData

`func (o *UpdateRowsRequest) GetData() map[string]interface{}`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *UpdateRowsRequest) GetDataOk() (*map[string]interface{}, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *UpdateRowsRequest) SetData(v map[string]interface{})`

SetData sets Data field to given value.


### GetReturnData

`func (o *UpdateRowsRequest) GetReturnData() bool`

GetReturnData returns the ReturnData field if non-nil, zero value otherwise.

### GetReturnDataOk

`func (o *UpdateRowsRequest) GetReturnDataOk() (*bool, bool)`

GetReturnDataOk returns a tuple with the ReturnData field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReturnData

`func (o *UpdateRowsRequest) SetReturnData(v bool)`

SetReturnData sets ReturnData field to given value.

### HasReturnData

`func (o *UpdateRowsRequest) HasReturnData() bool`

HasReturnData returns a boolean if a field has been set.

### GetDryRun

`func (o *UpdateRowsRequest) GetDryRun() bool`

GetDryRun returns the DryRun field if non-nil, zero value otherwise.

### GetDryRunOk

`func (o *UpdateRowsRequest) GetDryRunOk() (*bool, bool)`

GetDryRunOk returns a tuple with the DryRun field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDryRun

`func (o *UpdateRowsRequest) SetDryRun(v bool)`

SetDryRun sets DryRun field to given value.

### HasDryRun

`func (o *UpdateRowsRequest) HasDryRun() bool`

HasDryRun returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


