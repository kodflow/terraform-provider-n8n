# UsersPost201Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**User** | Pointer to [**UsersPost201ResponseUser**](UsersPost201ResponseUser.md) |  | [optional] 
**Error** | Pointer to **string** |  | [optional] 

## Methods

### NewUsersPost201Response

`func NewUsersPost201Response() *UsersPost201Response`

NewUsersPost201Response instantiates a new UsersPost201Response object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUsersPost201ResponseWithDefaults

`func NewUsersPost201ResponseWithDefaults() *UsersPost201Response`

NewUsersPost201ResponseWithDefaults instantiates a new UsersPost201Response object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUser

`func (o *UsersPost201Response) GetUser() UsersPost201ResponseUser`

GetUser returns the User field if non-nil, zero value otherwise.

### GetUserOk

`func (o *UsersPost201Response) GetUserOk() (*UsersPost201ResponseUser, bool)`

GetUserOk returns a tuple with the User field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUser

`func (o *UsersPost201Response) SetUser(v UsersPost201ResponseUser)`

SetUser sets User field to given value.

### HasUser

`func (o *UsersPost201Response) HasUser() bool`

HasUser returns a boolean if a field has been set.

### GetError

`func (o *UsersPost201Response) GetError() string`

GetError returns the Error field if non-nil, zero value otherwise.

### GetErrorOk

`func (o *UsersPost201Response) GetErrorOk() (*string, bool)`

GetErrorOk returns a tuple with the Error field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetError

`func (o *UsersPost201Response) SetError(v string)`

SetError sets Error field to given value.

### HasError

`func (o *UsersPost201Response) HasError() bool`

HasError returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


