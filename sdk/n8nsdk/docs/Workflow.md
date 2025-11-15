# Workflow

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** |  | [optional] [readonly] 
**Name** | **string** |  | 
**Active** | Pointer to **bool** |  | [optional] [readonly] 
**CreatedAt** | Pointer to **time.Time** |  | [optional] [readonly] 
**UpdatedAt** | Pointer to **time.Time** |  | [optional] [readonly] 
**Nodes** | [**[]Node**](Node.md) |  | 
**Connections** | **map[string]interface{}** |  | 
**Settings** | [**WorkflowSettings**](WorkflowSettings.md) |  | 
**StaticData** | Pointer to [**WorkflowStaticData**](WorkflowStaticData.md) |  | [optional] 
**Tags** | Pointer to [**[]Tag**](Tag.md) |  | [optional] [readonly] 
**Shared** | Pointer to [**[]SharedWorkflow**](SharedWorkflow.md) |  | [optional] 

## Methods

### NewWorkflow

`func NewWorkflow(name string, nodes []Node, connections map[string]interface{}, settings WorkflowSettings, ) *Workflow`

NewWorkflow instantiates a new Workflow object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowWithDefaults

`func NewWorkflowWithDefaults() *Workflow`

NewWorkflowWithDefaults instantiates a new Workflow object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *Workflow) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Workflow) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Workflow) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *Workflow) HasId() bool`

HasId returns a boolean if a field has been set.

### GetName

`func (o *Workflow) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *Workflow) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *Workflow) SetName(v string)`

SetName sets Name field to given value.


### GetActive

`func (o *Workflow) GetActive() bool`

GetActive returns the Active field if non-nil, zero value otherwise.

### GetActiveOk

`func (o *Workflow) GetActiveOk() (*bool, bool)`

GetActiveOk returns a tuple with the Active field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetActive

`func (o *Workflow) SetActive(v bool)`

SetActive sets Active field to given value.

### HasActive

`func (o *Workflow) HasActive() bool`

HasActive returns a boolean if a field has been set.

### GetCreatedAt

`func (o *Workflow) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *Workflow) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *Workflow) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *Workflow) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetUpdatedAt

`func (o *Workflow) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *Workflow) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *Workflow) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *Workflow) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.

### GetNodes

`func (o *Workflow) GetNodes() []Node`

GetNodes returns the Nodes field if non-nil, zero value otherwise.

### GetNodesOk

`func (o *Workflow) GetNodesOk() (*[]Node, bool)`

GetNodesOk returns a tuple with the Nodes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNodes

`func (o *Workflow) SetNodes(v []Node)`

SetNodes sets Nodes field to given value.


### GetConnections

`func (o *Workflow) GetConnections() map[string]interface{}`

GetConnections returns the Connections field if non-nil, zero value otherwise.

### GetConnectionsOk

`func (o *Workflow) GetConnectionsOk() (*map[string]interface{}, bool)`

GetConnectionsOk returns a tuple with the Connections field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConnections

`func (o *Workflow) SetConnections(v map[string]interface{})`

SetConnections sets Connections field to given value.


### GetSettings

`func (o *Workflow) GetSettings() WorkflowSettings`

GetSettings returns the Settings field if non-nil, zero value otherwise.

### GetSettingsOk

`func (o *Workflow) GetSettingsOk() (*WorkflowSettings, bool)`

GetSettingsOk returns a tuple with the Settings field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSettings

`func (o *Workflow) SetSettings(v WorkflowSettings)`

SetSettings sets Settings field to given value.


### GetStaticData

`func (o *Workflow) GetStaticData() WorkflowStaticData`

GetStaticData returns the StaticData field if non-nil, zero value otherwise.

### GetStaticDataOk

`func (o *Workflow) GetStaticDataOk() (*WorkflowStaticData, bool)`

GetStaticDataOk returns a tuple with the StaticData field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStaticData

`func (o *Workflow) SetStaticData(v WorkflowStaticData)`

SetStaticData sets StaticData field to given value.

### HasStaticData

`func (o *Workflow) HasStaticData() bool`

HasStaticData returns a boolean if a field has been set.

### GetTags

`func (o *Workflow) GetTags() []Tag`

GetTags returns the Tags field if non-nil, zero value otherwise.

### GetTagsOk

`func (o *Workflow) GetTagsOk() (*[]Tag, bool)`

GetTagsOk returns a tuple with the Tags field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTags

`func (o *Workflow) SetTags(v []Tag)`

SetTags sets Tags field to given value.

### HasTags

`func (o *Workflow) HasTags() bool`

HasTags returns a boolean if a field has been set.

### GetShared

`func (o *Workflow) GetShared() []SharedWorkflow`

GetShared returns the Shared field if non-nil, zero value otherwise.

### GetSharedOk

`func (o *Workflow) GetSharedOk() (*[]SharedWorkflow, bool)`

GetSharedOk returns a tuple with the Shared field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetShared

`func (o *Workflow) SetShared(v []SharedWorkflow)`

SetShared sets Shared field to given value.

### HasShared

`func (o *Workflow) HasShared() bool`

HasShared returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


