# ExecutionsIdRetryPostRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**LoadWorkflow** | Pointer to **bool** | Whether to load the currently saved workflow to execute instead of the one saved at the time of the execution. If set to true, it will retry with the latest version of the workflow. | [optional] 

## Methods

### NewExecutionsIdRetryPostRequest

`func NewExecutionsIdRetryPostRequest() *ExecutionsIdRetryPostRequest`

NewExecutionsIdRetryPostRequest instantiates a new ExecutionsIdRetryPostRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewExecutionsIdRetryPostRequestWithDefaults

`func NewExecutionsIdRetryPostRequestWithDefaults() *ExecutionsIdRetryPostRequest`

NewExecutionsIdRetryPostRequestWithDefaults instantiates a new ExecutionsIdRetryPostRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetLoadWorkflow

`func (o *ExecutionsIdRetryPostRequest) GetLoadWorkflow() bool`

GetLoadWorkflow returns the LoadWorkflow field if non-nil, zero value otherwise.

### GetLoadWorkflowOk

`func (o *ExecutionsIdRetryPostRequest) GetLoadWorkflowOk() (*bool, bool)`

GetLoadWorkflowOk returns a tuple with the LoadWorkflow field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLoadWorkflow

`func (o *ExecutionsIdRetryPostRequest) SetLoadWorkflow(v bool)`

SetLoadWorkflow sets LoadWorkflow field to given value.

### HasLoadWorkflow

`func (o *ExecutionsIdRetryPostRequest) HasLoadWorkflow() bool`

HasLoadWorkflow returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


