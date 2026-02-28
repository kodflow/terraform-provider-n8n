# UpsertDataTableRow200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int32** | The row ID (auto-generated) | [optional] 
**CreatedAt** | Pointer to **time.Time** | The date and time the row was created | [optional] 
**UpdatedAt** | Pointer to **time.Time** | The date and time the row was last updated | [optional] 

## Methods

### NewUpsertDataTableRow200Response

`func NewUpsertDataTableRow200Response() *UpsertDataTableRow200Response`

NewUpsertDataTableRow200Response instantiates a new UpsertDataTableRow200Response object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpsertDataTableRow200ResponseWithDefaults

`func NewUpsertDataTableRow200ResponseWithDefaults() *UpsertDataTableRow200Response`

NewUpsertDataTableRow200ResponseWithDefaults instantiates a new UpsertDataTableRow200Response object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *UpsertDataTableRow200Response) GetId() int32`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *UpsertDataTableRow200Response) GetIdOk() (*int32, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *UpsertDataTableRow200Response) SetId(v int32)`

SetId sets Id field to given value.

### HasId

`func (o *UpsertDataTableRow200Response) HasId() bool`

HasId returns a boolean if a field has been set.

### GetCreatedAt

`func (o *UpsertDataTableRow200Response) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *UpsertDataTableRow200Response) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *UpsertDataTableRow200Response) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *UpsertDataTableRow200Response) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetUpdatedAt

`func (o *UpsertDataTableRow200Response) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *UpsertDataTableRow200Response) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *UpsertDataTableRow200Response) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *UpsertDataTableRow200Response) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


