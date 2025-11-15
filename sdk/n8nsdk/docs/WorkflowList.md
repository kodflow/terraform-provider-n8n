# WorkflowList

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | Pointer to [**[]Workflow**](Workflow.md) |  | [optional] 
**NextCursor** | Pointer to **NullableString** | Paginate through workflows by setting the cursor parameter to a nextCursor attribute returned by a previous request. Default value fetches the first \&quot;page\&quot; of the collection. | [optional] 

## Methods

### NewWorkflowList

`func NewWorkflowList() *WorkflowList`

NewWorkflowList instantiates a new WorkflowList object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowListWithDefaults

`func NewWorkflowListWithDefaults() *WorkflowList`

NewWorkflowListWithDefaults instantiates a new WorkflowList object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *WorkflowList) GetData() []Workflow`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *WorkflowList) GetDataOk() (*[]Workflow, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *WorkflowList) SetData(v []Workflow)`

SetData sets Data field to given value.

### HasData

`func (o *WorkflowList) HasData() bool`

HasData returns a boolean if a field has been set.

### GetNextCursor

`func (o *WorkflowList) GetNextCursor() string`

GetNextCursor returns the NextCursor field if non-nil, zero value otherwise.

### GetNextCursorOk

`func (o *WorkflowList) GetNextCursorOk() (*string, bool)`

GetNextCursorOk returns a tuple with the NextCursor field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNextCursor

`func (o *WorkflowList) SetNextCursor(v string)`

SetNextCursor sets NextCursor field to given value.

### HasNextCursor

`func (o *WorkflowList) HasNextCursor() bool`

HasNextCursor returns a boolean if a field has been set.

### SetNextCursorNil

`func (o *WorkflowList) SetNextCursorNil(b bool)`

 SetNextCursorNil sets the value for NextCursor to be an explicit nil

### UnsetNextCursor
`func (o *WorkflowList) UnsetNextCursor()`

UnsetNextCursor ensures that no value is present for NextCursor, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


