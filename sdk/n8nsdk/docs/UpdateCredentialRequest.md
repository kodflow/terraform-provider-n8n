# UpdateCredentialRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | The name of the credential | [optional] 
**Type** | Pointer to **string** | The credential type. If changing type, data must also be provided. | [optional] 
**Data** | Pointer to **map[string]interface{}** | The credential data. Required when changing credential type. | [optional] 
**IsGlobal** | Pointer to **bool** | Whether this credential is available globally | [optional] 
**IsResolvable** | Pointer to **bool** | Whether this credential has resolvable fields | [optional] 
**IsPartialData** | Pointer to **bool** | If true, unredacts and merges existing credential data with the provided data. If false, replaces the entire data object. | [optional] [default to false]

## Methods

### NewUpdateCredentialRequest

`func NewUpdateCredentialRequest() *UpdateCredentialRequest`

NewUpdateCredentialRequest instantiates a new UpdateCredentialRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateCredentialRequestWithDefaults

`func NewUpdateCredentialRequestWithDefaults() *UpdateCredentialRequest`

NewUpdateCredentialRequestWithDefaults instantiates a new UpdateCredentialRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *UpdateCredentialRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *UpdateCredentialRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *UpdateCredentialRequest) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *UpdateCredentialRequest) HasName() bool`

HasName returns a boolean if a field has been set.

### GetType

`func (o *UpdateCredentialRequest) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *UpdateCredentialRequest) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *UpdateCredentialRequest) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *UpdateCredentialRequest) HasType() bool`

HasType returns a boolean if a field has been set.

### GetData

`func (o *UpdateCredentialRequest) GetData() map[string]interface{}`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *UpdateCredentialRequest) GetDataOk() (*map[string]interface{}, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *UpdateCredentialRequest) SetData(v map[string]interface{})`

SetData sets Data field to given value.

### HasData

`func (o *UpdateCredentialRequest) HasData() bool`

HasData returns a boolean if a field has been set.

### GetIsGlobal

`func (o *UpdateCredentialRequest) GetIsGlobal() bool`

GetIsGlobal returns the IsGlobal field if non-nil, zero value otherwise.

### GetIsGlobalOk

`func (o *UpdateCredentialRequest) GetIsGlobalOk() (*bool, bool)`

GetIsGlobalOk returns a tuple with the IsGlobal field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsGlobal

`func (o *UpdateCredentialRequest) SetIsGlobal(v bool)`

SetIsGlobal sets IsGlobal field to given value.

### HasIsGlobal

`func (o *UpdateCredentialRequest) HasIsGlobal() bool`

HasIsGlobal returns a boolean if a field has been set.

### GetIsResolvable

`func (o *UpdateCredentialRequest) GetIsResolvable() bool`

GetIsResolvable returns the IsResolvable field if non-nil, zero value otherwise.

### GetIsResolvableOk

`func (o *UpdateCredentialRequest) GetIsResolvableOk() (*bool, bool)`

GetIsResolvableOk returns a tuple with the IsResolvable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsResolvable

`func (o *UpdateCredentialRequest) SetIsResolvable(v bool)`

SetIsResolvable sets IsResolvable field to given value.

### HasIsResolvable

`func (o *UpdateCredentialRequest) HasIsResolvable() bool`

HasIsResolvable returns a boolean if a field has been set.

### GetIsPartialData

`func (o *UpdateCredentialRequest) GetIsPartialData() bool`

GetIsPartialData returns the IsPartialData field if non-nil, zero value otherwise.

### GetIsPartialDataOk

`func (o *UpdateCredentialRequest) GetIsPartialDataOk() (*bool, bool)`

GetIsPartialDataOk returns a tuple with the IsPartialData field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsPartialData

`func (o *UpdateCredentialRequest) SetIsPartialData(v bool)`

SetIsPartialData sets IsPartialData field to given value.

### HasIsPartialData

`func (o *UpdateCredentialRequest) HasIsPartialData() bool`

HasIsPartialData returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


