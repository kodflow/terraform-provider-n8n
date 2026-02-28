# CreateDataTableRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Name of the data table | 
**Columns** | [**[]CreateDataTableRequestColumnsInner**](CreateDataTableRequestColumnsInner.md) | Column definitions for the table | 

## Methods

### NewCreateDataTableRequest

`func NewCreateDataTableRequest(name string, columns []CreateDataTableRequestColumnsInner, ) *CreateDataTableRequest`

NewCreateDataTableRequest instantiates a new CreateDataTableRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateDataTableRequestWithDefaults

`func NewCreateDataTableRequestWithDefaults() *CreateDataTableRequest`

NewCreateDataTableRequestWithDefaults instantiates a new CreateDataTableRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *CreateDataTableRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *CreateDataTableRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *CreateDataTableRequest) SetName(v string)`

SetName sets Name field to given value.


### GetColumns

`func (o *CreateDataTableRequest) GetColumns() []CreateDataTableRequestColumnsInner`

GetColumns returns the Columns field if non-nil, zero value otherwise.

### GetColumnsOk

`func (o *CreateDataTableRequest) GetColumnsOk() (*[]CreateDataTableRequestColumnsInner, bool)`

GetColumnsOk returns a tuple with the Columns field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetColumns

`func (o *CreateDataTableRequest) SetColumns(v []CreateDataTableRequestColumnsInner)`

SetColumns sets Columns field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


