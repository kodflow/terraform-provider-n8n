# Execution

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **float32** |  | [optional] 
**Data** | Pointer to **map[string]interface{}** |  | [optional] 
**Finished** | Pointer to **bool** | Deprecated - use status instead | [optional] 
**Mode** | Pointer to **string** |  | [optional] 
**RetryOf** | Pointer to **float32** |  | [optional] 
**CreatedAt** | Pointer to **NullableTime** | The time at which the execution was created | [optional] [readonly] 
**RetrySuccessId** | Pointer to **NullableFloat32** |  | [optional] 
**StartedAt** | Pointer to **NullableTime** | The time at which the execution started | [optional] 
**StoppedAt** | Pointer to **NullableTime** | The time at which the execution stopped. Will only be null for executions that still have the status &#39;running&#39;. | [optional] 
**WorkflowId** | Pointer to **float32** |  | [optional] 
**WaitTill** | Pointer to **NullableTime** |  | [optional] 
**CustomData** | Pointer to **map[string]interface{}** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 

## Methods

### NewExecution

`func NewExecution() *Execution`

NewExecution instantiates a new Execution object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewExecutionWithDefaults

`func NewExecutionWithDefaults() *Execution`

NewExecutionWithDefaults instantiates a new Execution object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *Execution) GetId() float32`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Execution) GetIdOk() (*float32, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Execution) SetId(v float32)`

SetId sets Id field to given value.

### HasId

`func (o *Execution) HasId() bool`

HasId returns a boolean if a field has been set.

### GetData

`func (o *Execution) GetData() map[string]interface{}`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *Execution) GetDataOk() (*map[string]interface{}, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *Execution) SetData(v map[string]interface{})`

SetData sets Data field to given value.

### HasData

`func (o *Execution) HasData() bool`

HasData returns a boolean if a field has been set.

### GetFinished

`func (o *Execution) GetFinished() bool`

GetFinished returns the Finished field if non-nil, zero value otherwise.

### GetFinishedOk

`func (o *Execution) GetFinishedOk() (*bool, bool)`

GetFinishedOk returns a tuple with the Finished field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFinished

`func (o *Execution) SetFinished(v bool)`

SetFinished sets Finished field to given value.

### HasFinished

`func (o *Execution) HasFinished() bool`

HasFinished returns a boolean if a field has been set.

### GetMode

`func (o *Execution) GetMode() string`

GetMode returns the Mode field if non-nil, zero value otherwise.

### GetModeOk

`func (o *Execution) GetModeOk() (*string, bool)`

GetModeOk returns a tuple with the Mode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMode

`func (o *Execution) SetMode(v string)`

SetMode sets Mode field to given value.

### HasMode

`func (o *Execution) HasMode() bool`

HasMode returns a boolean if a field has been set.

### GetRetryOf

`func (o *Execution) GetRetryOf() float32`

GetRetryOf returns the RetryOf field if non-nil, zero value otherwise.

### GetRetryOfOk

`func (o *Execution) GetRetryOfOk() (*float32, bool)`

GetRetryOfOk returns a tuple with the RetryOf field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRetryOf

`func (o *Execution) SetRetryOf(v float32)`

SetRetryOf sets RetryOf field to given value.

### HasRetryOf

`func (o *Execution) HasRetryOf() bool`

HasRetryOf returns a boolean if a field has been set.

### GetCreatedAt

`func (o *Execution) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *Execution) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *Execution) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *Execution) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### SetCreatedAtNil

`func (o *Execution) SetCreatedAtNil(b bool)`

 SetCreatedAtNil sets the value for CreatedAt to be an explicit nil

### UnsetCreatedAt
`func (o *Execution) UnsetCreatedAt()`

UnsetCreatedAt ensures that no value is present for CreatedAt, not even an explicit nil
### GetRetrySuccessId

`func (o *Execution) GetRetrySuccessId() float32`

GetRetrySuccessId returns the RetrySuccessId field if non-nil, zero value otherwise.

### GetRetrySuccessIdOk

`func (o *Execution) GetRetrySuccessIdOk() (*float32, bool)`

GetRetrySuccessIdOk returns a tuple with the RetrySuccessId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRetrySuccessId

`func (o *Execution) SetRetrySuccessId(v float32)`

SetRetrySuccessId sets RetrySuccessId field to given value.

### HasRetrySuccessId

`func (o *Execution) HasRetrySuccessId() bool`

HasRetrySuccessId returns a boolean if a field has been set.

### SetRetrySuccessIdNil

`func (o *Execution) SetRetrySuccessIdNil(b bool)`

 SetRetrySuccessIdNil sets the value for RetrySuccessId to be an explicit nil

### UnsetRetrySuccessId
`func (o *Execution) UnsetRetrySuccessId()`

UnsetRetrySuccessId ensures that no value is present for RetrySuccessId, not even an explicit nil
### GetStartedAt

`func (o *Execution) GetStartedAt() time.Time`

GetStartedAt returns the StartedAt field if non-nil, zero value otherwise.

### GetStartedAtOk

`func (o *Execution) GetStartedAtOk() (*time.Time, bool)`

GetStartedAtOk returns a tuple with the StartedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStartedAt

`func (o *Execution) SetStartedAt(v time.Time)`

SetStartedAt sets StartedAt field to given value.

### HasStartedAt

`func (o *Execution) HasStartedAt() bool`

HasStartedAt returns a boolean if a field has been set.

### SetStartedAtNil

`func (o *Execution) SetStartedAtNil(b bool)`

 SetStartedAtNil sets the value for StartedAt to be an explicit nil

### UnsetStartedAt
`func (o *Execution) UnsetStartedAt()`

UnsetStartedAt ensures that no value is present for StartedAt, not even an explicit nil
### GetStoppedAt

`func (o *Execution) GetStoppedAt() time.Time`

GetStoppedAt returns the StoppedAt field if non-nil, zero value otherwise.

### GetStoppedAtOk

`func (o *Execution) GetStoppedAtOk() (*time.Time, bool)`

GetStoppedAtOk returns a tuple with the StoppedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStoppedAt

`func (o *Execution) SetStoppedAt(v time.Time)`

SetStoppedAt sets StoppedAt field to given value.

### HasStoppedAt

`func (o *Execution) HasStoppedAt() bool`

HasStoppedAt returns a boolean if a field has been set.

### SetStoppedAtNil

`func (o *Execution) SetStoppedAtNil(b bool)`

 SetStoppedAtNil sets the value for StoppedAt to be an explicit nil

### UnsetStoppedAt
`func (o *Execution) UnsetStoppedAt()`

UnsetStoppedAt ensures that no value is present for StoppedAt, not even an explicit nil
### GetWorkflowId

`func (o *Execution) GetWorkflowId() float32`

GetWorkflowId returns the WorkflowId field if non-nil, zero value otherwise.

### GetWorkflowIdOk

`func (o *Execution) GetWorkflowIdOk() (*float32, bool)`

GetWorkflowIdOk returns a tuple with the WorkflowId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowId

`func (o *Execution) SetWorkflowId(v float32)`

SetWorkflowId sets WorkflowId field to given value.

### HasWorkflowId

`func (o *Execution) HasWorkflowId() bool`

HasWorkflowId returns a boolean if a field has been set.

### GetWaitTill

`func (o *Execution) GetWaitTill() time.Time`

GetWaitTill returns the WaitTill field if non-nil, zero value otherwise.

### GetWaitTillOk

`func (o *Execution) GetWaitTillOk() (*time.Time, bool)`

GetWaitTillOk returns a tuple with the WaitTill field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWaitTill

`func (o *Execution) SetWaitTill(v time.Time)`

SetWaitTill sets WaitTill field to given value.

### HasWaitTill

`func (o *Execution) HasWaitTill() bool`

HasWaitTill returns a boolean if a field has been set.

### SetWaitTillNil

`func (o *Execution) SetWaitTillNil(b bool)`

 SetWaitTillNil sets the value for WaitTill to be an explicit nil

### UnsetWaitTill
`func (o *Execution) UnsetWaitTill()`

UnsetWaitTill ensures that no value is present for WaitTill, not even an explicit nil
### GetCustomData

`func (o *Execution) GetCustomData() map[string]interface{}`

GetCustomData returns the CustomData field if non-nil, zero value otherwise.

### GetCustomDataOk

`func (o *Execution) GetCustomDataOk() (*map[string]interface{}, bool)`

GetCustomDataOk returns a tuple with the CustomData field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCustomData

`func (o *Execution) SetCustomData(v map[string]interface{})`

SetCustomData sets CustomData field to given value.

### HasCustomData

`func (o *Execution) HasCustomData() bool`

HasCustomData returns a boolean if a field has been set.

### GetStatus

`func (o *Execution) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *Execution) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *Execution) SetStatus(v string)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *Execution) HasStatus() bool`

HasStatus returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


