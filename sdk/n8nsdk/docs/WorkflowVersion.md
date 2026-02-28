# WorkflowVersion

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**VersionId** | **string** | The version ID of this workflow snapshot | [readonly] 
**WorkflowId** | **string** | The workflow ID this version belongs to | [readonly] 
**Nodes** | [**[]Node**](Node.md) | Nodes as they were in this version | [readonly] 
**Connections** | **map[string]interface{}** | Connections as they were in this version | [readonly] 
**Authors** | **string** | Authors who created this version | [readonly] 
**Name** | Pointer to **NullableString** | Workflow name at this version | [optional] 
**Description** | Pointer to **NullableString** | Workflow description at this version | [optional] 
**CreatedAt** | Pointer to **time.Time** | When this version was created | [optional] [readonly] 
**UpdatedAt** | Pointer to **time.Time** | When this version was last updated | [optional] [readonly] 

## Methods

### NewWorkflowVersion

`func NewWorkflowVersion(versionId string, workflowId string, nodes []Node, connections map[string]interface{}, authors string, ) *WorkflowVersion`

NewWorkflowVersion instantiates a new WorkflowVersion object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowVersionWithDefaults

`func NewWorkflowVersionWithDefaults() *WorkflowVersion`

NewWorkflowVersionWithDefaults instantiates a new WorkflowVersion object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetVersionId

`func (o *WorkflowVersion) GetVersionId() string`

GetVersionId returns the VersionId field if non-nil, zero value otherwise.

### GetVersionIdOk

`func (o *WorkflowVersion) GetVersionIdOk() (*string, bool)`

GetVersionIdOk returns a tuple with the VersionId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersionId

`func (o *WorkflowVersion) SetVersionId(v string)`

SetVersionId sets VersionId field to given value.


### GetWorkflowId

`func (o *WorkflowVersion) GetWorkflowId() string`

GetWorkflowId returns the WorkflowId field if non-nil, zero value otherwise.

### GetWorkflowIdOk

`func (o *WorkflowVersion) GetWorkflowIdOk() (*string, bool)`

GetWorkflowIdOk returns a tuple with the WorkflowId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowId

`func (o *WorkflowVersion) SetWorkflowId(v string)`

SetWorkflowId sets WorkflowId field to given value.


### GetNodes

`func (o *WorkflowVersion) GetNodes() []Node`

GetNodes returns the Nodes field if non-nil, zero value otherwise.

### GetNodesOk

`func (o *WorkflowVersion) GetNodesOk() (*[]Node, bool)`

GetNodesOk returns a tuple with the Nodes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNodes

`func (o *WorkflowVersion) SetNodes(v []Node)`

SetNodes sets Nodes field to given value.


### GetConnections

`func (o *WorkflowVersion) GetConnections() map[string]interface{}`

GetConnections returns the Connections field if non-nil, zero value otherwise.

### GetConnectionsOk

`func (o *WorkflowVersion) GetConnectionsOk() (*map[string]interface{}, bool)`

GetConnectionsOk returns a tuple with the Connections field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConnections

`func (o *WorkflowVersion) SetConnections(v map[string]interface{})`

SetConnections sets Connections field to given value.


### GetAuthors

`func (o *WorkflowVersion) GetAuthors() string`

GetAuthors returns the Authors field if non-nil, zero value otherwise.

### GetAuthorsOk

`func (o *WorkflowVersion) GetAuthorsOk() (*string, bool)`

GetAuthorsOk returns a tuple with the Authors field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthors

`func (o *WorkflowVersion) SetAuthors(v string)`

SetAuthors sets Authors field to given value.


### GetName

`func (o *WorkflowVersion) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *WorkflowVersion) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *WorkflowVersion) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *WorkflowVersion) HasName() bool`

HasName returns a boolean if a field has been set.

### SetNameNil

`func (o *WorkflowVersion) SetNameNil(b bool)`

 SetNameNil sets the value for Name to be an explicit nil

### UnsetName
`func (o *WorkflowVersion) UnsetName()`

UnsetName ensures that no value is present for Name, not even an explicit nil
### GetDescription

`func (o *WorkflowVersion) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *WorkflowVersion) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *WorkflowVersion) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *WorkflowVersion) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### SetDescriptionNil

`func (o *WorkflowVersion) SetDescriptionNil(b bool)`

 SetDescriptionNil sets the value for Description to be an explicit nil

### UnsetDescription
`func (o *WorkflowVersion) UnsetDescription()`

UnsetDescription ensures that no value is present for Description, not even an explicit nil
### GetCreatedAt

`func (o *WorkflowVersion) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *WorkflowVersion) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *WorkflowVersion) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *WorkflowVersion) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetUpdatedAt

`func (o *WorkflowVersion) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *WorkflowVersion) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *WorkflowVersion) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *WorkflowVersion) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


