# ProjectMemberList

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | Pointer to [**[]ProjectMember**](ProjectMember.md) |  | [optional] 
**NextCursor** | Pointer to **NullableString** | Paginate through project members by setting the cursor parameter to the nextCursor attribute returned by a previous request. Default value fetches the first page of the collection. | [optional] 

## Methods

### NewProjectMemberList

`func NewProjectMemberList() *ProjectMemberList`

NewProjectMemberList instantiates a new ProjectMemberList object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewProjectMemberListWithDefaults

`func NewProjectMemberListWithDefaults() *ProjectMemberList`

NewProjectMemberListWithDefaults instantiates a new ProjectMemberList object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *ProjectMemberList) GetData() []ProjectMember`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *ProjectMemberList) GetDataOk() (*[]ProjectMember, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *ProjectMemberList) SetData(v []ProjectMember)`

SetData sets Data field to given value.

### HasData

`func (o *ProjectMemberList) HasData() bool`

HasData returns a boolean if a field has been set.

### GetNextCursor

`func (o *ProjectMemberList) GetNextCursor() string`

GetNextCursor returns the NextCursor field if non-nil, zero value otherwise.

### GetNextCursorOk

`func (o *ProjectMemberList) GetNextCursorOk() (*string, bool)`

GetNextCursorOk returns a tuple with the NextCursor field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNextCursor

`func (o *ProjectMemberList) SetNextCursor(v string)`

SetNextCursor sets NextCursor field to given value.

### HasNextCursor

`func (o *ProjectMemberList) HasNextCursor() bool`

HasNextCursor returns a boolean if a field has been set.

### SetNextCursorNil

`func (o *ProjectMemberList) SetNextCursorNil(b bool)`

 SetNextCursorNil sets the value for NextCursor to be an explicit nil

### UnsetNextCursor
`func (o *ProjectMemberList) UnsetNextCursor()`

UnsetNextCursor ensures that no value is present for NextCursor, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


