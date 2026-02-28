# CredentialList

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | Pointer to [**[]CredentialListItem**](CredentialListItem.md) |  | [optional] 
**NextCursor** | Pointer to **NullableString** | Paginate through credentials by setting the cursor parameter to a nextCursor attribute returned by a previous request. Default value fetches the first \&quot;page\&quot; of the collection. | [optional] 

## Methods

### NewCredentialList

`func NewCredentialList() *CredentialList`

NewCredentialList instantiates a new CredentialList object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCredentialListWithDefaults

`func NewCredentialListWithDefaults() *CredentialList`

NewCredentialListWithDefaults instantiates a new CredentialList object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *CredentialList) GetData() []CredentialListItem`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *CredentialList) GetDataOk() (*[]CredentialListItem, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *CredentialList) SetData(v []CredentialListItem)`

SetData sets Data field to given value.

### HasData

`func (o *CredentialList) HasData() bool`

HasData returns a boolean if a field has been set.

### GetNextCursor

`func (o *CredentialList) GetNextCursor() string`

GetNextCursor returns the NextCursor field if non-nil, zero value otherwise.

### GetNextCursorOk

`func (o *CredentialList) GetNextCursorOk() (*string, bool)`

GetNextCursorOk returns a tuple with the NextCursor field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNextCursor

`func (o *CredentialList) SetNextCursor(v string)`

SetNextCursor sets NextCursor field to given value.

### HasNextCursor

`func (o *CredentialList) HasNextCursor() bool`

HasNextCursor returns a boolean if a field has been set.

### SetNextCursorNil

`func (o *CredentialList) SetNextCursorNil(b bool)`

 SetNextCursorNil sets the value for NextCursor to be an explicit nil

### UnsetNextCursor
`func (o *CredentialList) UnsetNextCursor()`

UnsetNextCursor ensures that no value is present for NextCursor, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


