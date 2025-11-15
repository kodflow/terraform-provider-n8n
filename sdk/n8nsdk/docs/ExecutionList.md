# ExecutionList

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | Pointer to [**[]Execution**](Execution.md) |  | [optional] 
**NextCursor** | Pointer to **NullableString** | Paginate through executions by setting the cursor parameter to a nextCursor attribute returned by a previous request. Default value fetches the first \&quot;page\&quot; of the collection. | [optional] 

## Methods

### NewExecutionList

`func NewExecutionList() *ExecutionList`

NewExecutionList instantiates a new ExecutionList object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewExecutionListWithDefaults

`func NewExecutionListWithDefaults() *ExecutionList`

NewExecutionListWithDefaults instantiates a new ExecutionList object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *ExecutionList) GetData() []Execution`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *ExecutionList) GetDataOk() (*[]Execution, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *ExecutionList) SetData(v []Execution)`

SetData sets Data field to given value.

### HasData

`func (o *ExecutionList) HasData() bool`

HasData returns a boolean if a field has been set.

### GetNextCursor

`func (o *ExecutionList) GetNextCursor() string`

GetNextCursor returns the NextCursor field if non-nil, zero value otherwise.

### GetNextCursorOk

`func (o *ExecutionList) GetNextCursorOk() (*string, bool)`

GetNextCursorOk returns a tuple with the NextCursor field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNextCursor

`func (o *ExecutionList) SetNextCursor(v string)`

SetNextCursor sets NextCursor field to given value.

### HasNextCursor

`func (o *ExecutionList) HasNextCursor() bool`

HasNextCursor returns a boolean if a field has been set.

### SetNextCursorNil

`func (o *ExecutionList) SetNextCursorNil(b bool)`

 SetNextCursorNil sets the value for NextCursor to be an explicit nil

### UnsetNextCursor
`func (o *ExecutionList) UnsetNextCursor()`

UnsetNextCursor ensures that no value is present for NextCursor, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


