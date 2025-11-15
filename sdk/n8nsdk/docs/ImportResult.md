# ImportResult

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Variables** | Pointer to [**ImportResultVariables**](ImportResultVariables.md) |  | [optional] 
**Credentials** | Pointer to [**[]ImportResultCredentialsInner**](ImportResultCredentialsInner.md) |  | [optional] 
**Workflows** | Pointer to [**[]ImportResultWorkflowsInner**](ImportResultWorkflowsInner.md) |  | [optional] 
**Tags** | Pointer to [**ImportResultTags**](ImportResultTags.md) |  | [optional] 

## Methods

### NewImportResult

`func NewImportResult() *ImportResult`

NewImportResult instantiates a new ImportResult object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewImportResultWithDefaults

`func NewImportResultWithDefaults() *ImportResult`

NewImportResultWithDefaults instantiates a new ImportResult object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetVariables

`func (o *ImportResult) GetVariables() ImportResultVariables`

GetVariables returns the Variables field if non-nil, zero value otherwise.

### GetVariablesOk

`func (o *ImportResult) GetVariablesOk() (*ImportResultVariables, bool)`

GetVariablesOk returns a tuple with the Variables field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVariables

`func (o *ImportResult) SetVariables(v ImportResultVariables)`

SetVariables sets Variables field to given value.

### HasVariables

`func (o *ImportResult) HasVariables() bool`

HasVariables returns a boolean if a field has been set.

### GetCredentials

`func (o *ImportResult) GetCredentials() []ImportResultCredentialsInner`

GetCredentials returns the Credentials field if non-nil, zero value otherwise.

### GetCredentialsOk

`func (o *ImportResult) GetCredentialsOk() (*[]ImportResultCredentialsInner, bool)`

GetCredentialsOk returns a tuple with the Credentials field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCredentials

`func (o *ImportResult) SetCredentials(v []ImportResultCredentialsInner)`

SetCredentials sets Credentials field to given value.

### HasCredentials

`func (o *ImportResult) HasCredentials() bool`

HasCredentials returns a boolean if a field has been set.

### GetWorkflows

`func (o *ImportResult) GetWorkflows() []ImportResultWorkflowsInner`

GetWorkflows returns the Workflows field if non-nil, zero value otherwise.

### GetWorkflowsOk

`func (o *ImportResult) GetWorkflowsOk() (*[]ImportResultWorkflowsInner, bool)`

GetWorkflowsOk returns a tuple with the Workflows field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflows

`func (o *ImportResult) SetWorkflows(v []ImportResultWorkflowsInner)`

SetWorkflows sets Workflows field to given value.

### HasWorkflows

`func (o *ImportResult) HasWorkflows() bool`

HasWorkflows returns a boolean if a field has been set.

### GetTags

`func (o *ImportResult) GetTags() ImportResultTags`

GetTags returns the Tags field if non-nil, zero value otherwise.

### GetTagsOk

`func (o *ImportResult) GetTagsOk() (*ImportResultTags, bool)`

GetTagsOk returns a tuple with the Tags field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTags

`func (o *ImportResult) SetTags(v ImportResultTags)`

SetTags sets Tags field to given value.

### HasTags

`func (o *ImportResult) HasTags() bool`

HasTags returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


