# WorkflowsIdActivatePostRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**VersionId** | Pointer to **string** | The specific version ID to activate or publish. If not provided, the latest version is used. | [optional] 
**Name** | Pointer to **string** | Optional name for the workflow version during activation. | [optional] 
**Description** | Pointer to **string** | Optional description for the workflow version during activation. | [optional] 

## Methods

### NewWorkflowsIdActivatePostRequest

`func NewWorkflowsIdActivatePostRequest() *WorkflowsIdActivatePostRequest`

NewWorkflowsIdActivatePostRequest instantiates a new WorkflowsIdActivatePostRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowsIdActivatePostRequestWithDefaults

`func NewWorkflowsIdActivatePostRequestWithDefaults() *WorkflowsIdActivatePostRequest`

NewWorkflowsIdActivatePostRequestWithDefaults instantiates a new WorkflowsIdActivatePostRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetVersionId

`func (o *WorkflowsIdActivatePostRequest) GetVersionId() string`

GetVersionId returns the VersionId field if non-nil, zero value otherwise.

### GetVersionIdOk

`func (o *WorkflowsIdActivatePostRequest) GetVersionIdOk() (*string, bool)`

GetVersionIdOk returns a tuple with the VersionId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersionId

`func (o *WorkflowsIdActivatePostRequest) SetVersionId(v string)`

SetVersionId sets VersionId field to given value.

### HasVersionId

`func (o *WorkflowsIdActivatePostRequest) HasVersionId() bool`

HasVersionId returns a boolean if a field has been set.

### GetName

`func (o *WorkflowsIdActivatePostRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *WorkflowsIdActivatePostRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *WorkflowsIdActivatePostRequest) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *WorkflowsIdActivatePostRequest) HasName() bool`

HasName returns a boolean if a field has been set.

### GetDescription

`func (o *WorkflowsIdActivatePostRequest) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *WorkflowsIdActivatePostRequest) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *WorkflowsIdActivatePostRequest) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *WorkflowsIdActivatePostRequest) HasDescription() bool`

HasDescription returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


