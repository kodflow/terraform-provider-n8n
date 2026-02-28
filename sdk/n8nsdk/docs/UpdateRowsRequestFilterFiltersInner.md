# UpdateRowsRequestFilterFiltersInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ColumnName** | **string** |  | 
**Condition** | **string** |  | 
**Value** | **interface{}** |  | 

## Methods

### NewUpdateRowsRequestFilterFiltersInner

`func NewUpdateRowsRequestFilterFiltersInner(columnName string, condition string, value interface{}, ) *UpdateRowsRequestFilterFiltersInner`

NewUpdateRowsRequestFilterFiltersInner instantiates a new UpdateRowsRequestFilterFiltersInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateRowsRequestFilterFiltersInnerWithDefaults

`func NewUpdateRowsRequestFilterFiltersInnerWithDefaults() *UpdateRowsRequestFilterFiltersInner`

NewUpdateRowsRequestFilterFiltersInnerWithDefaults instantiates a new UpdateRowsRequestFilterFiltersInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetColumnName

`func (o *UpdateRowsRequestFilterFiltersInner) GetColumnName() string`

GetColumnName returns the ColumnName field if non-nil, zero value otherwise.

### GetColumnNameOk

`func (o *UpdateRowsRequestFilterFiltersInner) GetColumnNameOk() (*string, bool)`

GetColumnNameOk returns a tuple with the ColumnName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetColumnName

`func (o *UpdateRowsRequestFilterFiltersInner) SetColumnName(v string)`

SetColumnName sets ColumnName field to given value.


### GetCondition

`func (o *UpdateRowsRequestFilterFiltersInner) GetCondition() string`

GetCondition returns the Condition field if non-nil, zero value otherwise.

### GetConditionOk

`func (o *UpdateRowsRequestFilterFiltersInner) GetConditionOk() (*string, bool)`

GetConditionOk returns a tuple with the Condition field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCondition

`func (o *UpdateRowsRequestFilterFiltersInner) SetCondition(v string)`

SetCondition sets Condition field to given value.


### GetValue

`func (o *UpdateRowsRequestFilterFiltersInner) GetValue() interface{}`

GetValue returns the Value field if non-nil, zero value otherwise.

### GetValueOk

`func (o *UpdateRowsRequestFilterFiltersInner) GetValueOk() (*interface{}, bool)`

GetValueOk returns a tuple with the Value field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetValue

`func (o *UpdateRowsRequestFilterFiltersInner) SetValue(v interface{})`

SetValue sets Value field to given value.


### SetValueNil

`func (o *UpdateRowsRequestFilterFiltersInner) SetValueNil(b bool)`

 SetValueNil sets the value for Value to be an explicit nil

### UnsetValue
`func (o *UpdateRowsRequestFilterFiltersInner) UnsetValue()`

UnsetValue ensures that no value is present for Value, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


