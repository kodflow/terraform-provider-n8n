# ProjectMember

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** | The user&#39;s unique identifier. | [optional] [readonly] 
**Email** | Pointer to **string** | The user&#39;s email address. | [optional] [readonly] 
**FirstName** | Pointer to **string** | The user&#39;s first name. | [optional] [readonly] 
**LastName** | Pointer to **string** | The user&#39;s last name. | [optional] [readonly] 
**CreatedAt** | Pointer to **time.Time** | When the user was created. | [optional] [readonly] 
**UpdatedAt** | Pointer to **time.Time** | When the user was last updated. | [optional] [readonly] 
**Role** | Pointer to **string** | The user&#39;s role in the project (e.g. project:admin, project:viewer). | [optional] [readonly] 

## Methods

### NewProjectMember

`func NewProjectMember() *ProjectMember`

NewProjectMember instantiates a new ProjectMember object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewProjectMemberWithDefaults

`func NewProjectMemberWithDefaults() *ProjectMember`

NewProjectMemberWithDefaults instantiates a new ProjectMember object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *ProjectMember) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *ProjectMember) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *ProjectMember) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *ProjectMember) HasId() bool`

HasId returns a boolean if a field has been set.

### GetEmail

`func (o *ProjectMember) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *ProjectMember) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *ProjectMember) SetEmail(v string)`

SetEmail sets Email field to given value.

### HasEmail

`func (o *ProjectMember) HasEmail() bool`

HasEmail returns a boolean if a field has been set.

### GetFirstName

`func (o *ProjectMember) GetFirstName() string`

GetFirstName returns the FirstName field if non-nil, zero value otherwise.

### GetFirstNameOk

`func (o *ProjectMember) GetFirstNameOk() (*string, bool)`

GetFirstNameOk returns a tuple with the FirstName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFirstName

`func (o *ProjectMember) SetFirstName(v string)`

SetFirstName sets FirstName field to given value.

### HasFirstName

`func (o *ProjectMember) HasFirstName() bool`

HasFirstName returns a boolean if a field has been set.

### GetLastName

`func (o *ProjectMember) GetLastName() string`

GetLastName returns the LastName field if non-nil, zero value otherwise.

### GetLastNameOk

`func (o *ProjectMember) GetLastNameOk() (*string, bool)`

GetLastNameOk returns a tuple with the LastName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastName

`func (o *ProjectMember) SetLastName(v string)`

SetLastName sets LastName field to given value.

### HasLastName

`func (o *ProjectMember) HasLastName() bool`

HasLastName returns a boolean if a field has been set.

### GetCreatedAt

`func (o *ProjectMember) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *ProjectMember) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *ProjectMember) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *ProjectMember) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetUpdatedAt

`func (o *ProjectMember) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *ProjectMember) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *ProjectMember) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *ProjectMember) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.

### GetRole

`func (o *ProjectMember) GetRole() string`

GetRole returns the Role field if non-nil, zero value otherwise.

### GetRoleOk

`func (o *ProjectMember) GetRoleOk() (*string, bool)`

GetRoleOk returns a tuple with the Role field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRole

`func (o *ProjectMember) SetRole(v string)`

SetRole sets Role field to given value.

### HasRole

`func (o *ProjectMember) HasRole() bool`

HasRole returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


