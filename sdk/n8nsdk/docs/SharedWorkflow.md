# SharedWorkflow

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** | Shared workflow ID | [optional] 
**ProjectId** | Pointer to **string** | Project ID | [optional] 
**UserId** | Pointer to **string** | User ID | [optional] 
**Role** | Pointer to **string** | User role | [optional] 
**WorkflowId** | Pointer to **string** |  | [optional] 
**Project** | Pointer to [**SharedWorkflowProject1**](SharedWorkflowProject1.md) |  | [optional] 
**CreatedAt** | Pointer to **time.Time** |  | [optional] [readonly] 
**UpdatedAt** | Pointer to **time.Time** |  | [optional] [readonly] 

## Methods

### NewSharedWorkflow

`func NewSharedWorkflow() *SharedWorkflow`

NewSharedWorkflow instantiates a new SharedWorkflow object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSharedWorkflowWithDefaults

`func NewSharedWorkflowWithDefaults() *SharedWorkflow`

NewSharedWorkflowWithDefaults instantiates a new SharedWorkflow object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *SharedWorkflow) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *SharedWorkflow) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *SharedWorkflow) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *SharedWorkflow) HasId() bool`

HasId returns a boolean if a field has been set.

### GetProjectId

`func (o *SharedWorkflow) GetProjectId() string`

GetProjectId returns the ProjectId field if non-nil, zero value otherwise.

### GetProjectIdOk

`func (o *SharedWorkflow) GetProjectIdOk() (*string, bool)`

GetProjectIdOk returns a tuple with the ProjectId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProjectId

`func (o *SharedWorkflow) SetProjectId(v string)`

SetProjectId sets ProjectId field to given value.

### HasProjectId

`func (o *SharedWorkflow) HasProjectId() bool`

HasProjectId returns a boolean if a field has been set.

### GetUserId

`func (o *SharedWorkflow) GetUserId() string`

GetUserId returns the UserId field if non-nil, zero value otherwise.

### GetUserIdOk

`func (o *SharedWorkflow) GetUserIdOk() (*string, bool)`

GetUserIdOk returns a tuple with the UserId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserId

`func (o *SharedWorkflow) SetUserId(v string)`

SetUserId sets UserId field to given value.

### HasUserId

`func (o *SharedWorkflow) HasUserId() bool`

HasUserId returns a boolean if a field has been set.

### GetRole

`func (o *SharedWorkflow) GetRole() string`

GetRole returns the Role field if non-nil, zero value otherwise.

### GetRoleOk

`func (o *SharedWorkflow) GetRoleOk() (*string, bool)`

GetRoleOk returns a tuple with the Role field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRole

`func (o *SharedWorkflow) SetRole(v string)`

SetRole sets Role field to given value.

### HasRole

`func (o *SharedWorkflow) HasRole() bool`

HasRole returns a boolean if a field has been set.

### GetWorkflowId

`func (o *SharedWorkflow) GetWorkflowId() string`

GetWorkflowId returns the WorkflowId field if non-nil, zero value otherwise.

### GetWorkflowIdOk

`func (o *SharedWorkflow) GetWorkflowIdOk() (*string, bool)`

GetWorkflowIdOk returns a tuple with the WorkflowId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowId

`func (o *SharedWorkflow) SetWorkflowId(v string)`

SetWorkflowId sets WorkflowId field to given value.

### HasWorkflowId

`func (o *SharedWorkflow) HasWorkflowId() bool`

HasWorkflowId returns a boolean if a field has been set.

### GetProject

`func (o *SharedWorkflow) GetProject() SharedWorkflowProject1`

GetProject returns the Project field if non-nil, zero value otherwise.

### GetProjectOk

`func (o *SharedWorkflow) GetProjectOk() (*SharedWorkflowProject1, bool)`

GetProjectOk returns a tuple with the Project field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProject

`func (o *SharedWorkflow) SetProject(v SharedWorkflowProject1)`

SetProject sets Project field to given value.

### HasProject

`func (o *SharedWorkflow) HasProject() bool`

HasProject returns a boolean if a field has been set.

### GetCreatedAt

`func (o *SharedWorkflow) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *SharedWorkflow) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *SharedWorkflow) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *SharedWorkflow) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetUpdatedAt

`func (o *SharedWorkflow) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *SharedWorkflow) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *SharedWorkflow) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *SharedWorkflow) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


