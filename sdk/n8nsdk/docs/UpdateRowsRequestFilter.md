# UpdateRowsRequestFilter

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Type** | Pointer to **string** |  | [optional] [default to "and"]
**Filters** | [**[]UpdateRowsRequestFilterFiltersInner**](UpdateRowsRequestFilterFiltersInner.md) |  | 

## Methods

### NewUpdateRowsRequestFilter

`func NewUpdateRowsRequestFilter(filters []UpdateRowsRequestFilterFiltersInner, ) *UpdateRowsRequestFilter`

NewUpdateRowsRequestFilter instantiates a new UpdateRowsRequestFilter object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateRowsRequestFilterWithDefaults

`func NewUpdateRowsRequestFilterWithDefaults() *UpdateRowsRequestFilter`

NewUpdateRowsRequestFilterWithDefaults instantiates a new UpdateRowsRequestFilter object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetType

`func (o *UpdateRowsRequestFilter) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *UpdateRowsRequestFilter) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *UpdateRowsRequestFilter) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *UpdateRowsRequestFilter) HasType() bool`

HasType returns a boolean if a field has been set.

### GetFilters

`func (o *UpdateRowsRequestFilter) GetFilters() []UpdateRowsRequestFilterFiltersInner`

GetFilters returns the Filters field if non-nil, zero value otherwise.

### GetFiltersOk

`func (o *UpdateRowsRequestFilter) GetFiltersOk() (*[]UpdateRowsRequestFilterFiltersInner, bool)`

GetFiltersOk returns a tuple with the Filters field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFilters

`func (o *UpdateRowsRequestFilter) SetFilters(v []UpdateRowsRequestFilterFiltersInner)`

SetFilters sets Filters field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


