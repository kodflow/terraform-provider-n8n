# ActiveVersion

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**VersionId** | Pointer to **string** | Unique identifier for this workflow version | [optional] [readonly] 
**WorkflowId** | Pointer to **string** | The workflow this version belongs to | [optional] [readonly] 
**Nodes** | Pointer to [**[]Node**](Node.md) |  | [optional] [readonly] 
**Connections** | Pointer to **map[string]interface{}** |  | [optional] [readonly] 
**Authors** | Pointer to **string** | Comma-separated list of author IDs who contributed to this version | [optional] [readonly] 
**CreatedAt** | Pointer to **time.Time** |  | [optional] [readonly] 
**UpdatedAt** | Pointer to **time.Time** |  | [optional] [readonly] 

## Methods

### NewActiveVersion

`func NewActiveVersion() *ActiveVersion`

NewActiveVersion instantiates a new ActiveVersion object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewActiveVersionWithDefaults

`func NewActiveVersionWithDefaults() *ActiveVersion`

NewActiveVersionWithDefaults instantiates a new ActiveVersion object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetVersionId

`func (o *ActiveVersion) GetVersionId() string`

GetVersionId returns the VersionId field if non-nil, zero value otherwise.

### GetVersionIdOk

`func (o *ActiveVersion) GetVersionIdOk() (*string, bool)`

GetVersionIdOk returns a tuple with the VersionId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersionId

`func (o *ActiveVersion) SetVersionId(v string)`

SetVersionId sets VersionId field to given value.

### HasVersionId

`func (o *ActiveVersion) HasVersionId() bool`

HasVersionId returns a boolean if a field has been set.

### GetWorkflowId

`func (o *ActiveVersion) GetWorkflowId() string`

GetWorkflowId returns the WorkflowId field if non-nil, zero value otherwise.

### GetWorkflowIdOk

`func (o *ActiveVersion) GetWorkflowIdOk() (*string, bool)`

GetWorkflowIdOk returns a tuple with the WorkflowId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowId

`func (o *ActiveVersion) SetWorkflowId(v string)`

SetWorkflowId sets WorkflowId field to given value.

### HasWorkflowId

`func (o *ActiveVersion) HasWorkflowId() bool`

HasWorkflowId returns a boolean if a field has been set.

### GetNodes

`func (o *ActiveVersion) GetNodes() []Node`

GetNodes returns the Nodes field if non-nil, zero value otherwise.

### GetNodesOk

`func (o *ActiveVersion) GetNodesOk() (*[]Node, bool)`

GetNodesOk returns a tuple with the Nodes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNodes

`func (o *ActiveVersion) SetNodes(v []Node)`

SetNodes sets Nodes field to given value.

### HasNodes

`func (o *ActiveVersion) HasNodes() bool`

HasNodes returns a boolean if a field has been set.

### GetConnections

`func (o *ActiveVersion) GetConnections() map[string]interface{}`

GetConnections returns the Connections field if non-nil, zero value otherwise.

### GetConnectionsOk

`func (o *ActiveVersion) GetConnectionsOk() (*map[string]interface{}, bool)`

GetConnectionsOk returns a tuple with the Connections field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConnections

`func (o *ActiveVersion) SetConnections(v map[string]interface{})`

SetConnections sets Connections field to given value.

### HasConnections

`func (o *ActiveVersion) HasConnections() bool`

HasConnections returns a boolean if a field has been set.

### GetAuthors

`func (o *ActiveVersion) GetAuthors() string`

GetAuthors returns the Authors field if non-nil, zero value otherwise.

### GetAuthorsOk

`func (o *ActiveVersion) GetAuthorsOk() (*string, bool)`

GetAuthorsOk returns a tuple with the Authors field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthors

`func (o *ActiveVersion) SetAuthors(v string)`

SetAuthors sets Authors field to given value.

### HasAuthors

`func (o *ActiveVersion) HasAuthors() bool`

HasAuthors returns a boolean if a field has been set.

### GetCreatedAt

`func (o *ActiveVersion) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *ActiveVersion) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *ActiveVersion) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *ActiveVersion) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetUpdatedAt

`func (o *ActiveVersion) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *ActiveVersion) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *ActiveVersion) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *ActiveVersion) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


