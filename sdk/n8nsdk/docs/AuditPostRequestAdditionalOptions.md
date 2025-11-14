# AuditPostRequestAdditionalOptions

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**DaysAbandonedWorkflow** | Pointer to **int32** | Days for a workflow to be considered abandoned if not executed | [optional] 
**Categories** | Pointer to **[]string** |  | [optional] 

## Methods

### NewAuditPostRequestAdditionalOptions

`func NewAuditPostRequestAdditionalOptions() *AuditPostRequestAdditionalOptions`

NewAuditPostRequestAdditionalOptions instantiates a new AuditPostRequestAdditionalOptions object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAuditPostRequestAdditionalOptionsWithDefaults

`func NewAuditPostRequestAdditionalOptionsWithDefaults() *AuditPostRequestAdditionalOptions`

NewAuditPostRequestAdditionalOptionsWithDefaults instantiates a new AuditPostRequestAdditionalOptions object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDaysAbandonedWorkflow

`func (o *AuditPostRequestAdditionalOptions) GetDaysAbandonedWorkflow() int32`

GetDaysAbandonedWorkflow returns the DaysAbandonedWorkflow field if non-nil, zero value otherwise.

### GetDaysAbandonedWorkflowOk

`func (o *AuditPostRequestAdditionalOptions) GetDaysAbandonedWorkflowOk() (*int32, bool)`

GetDaysAbandonedWorkflowOk returns a tuple with the DaysAbandonedWorkflow field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDaysAbandonedWorkflow

`func (o *AuditPostRequestAdditionalOptions) SetDaysAbandonedWorkflow(v int32)`

SetDaysAbandonedWorkflow sets DaysAbandonedWorkflow field to given value.

### HasDaysAbandonedWorkflow

`func (o *AuditPostRequestAdditionalOptions) HasDaysAbandonedWorkflow() bool`

HasDaysAbandonedWorkflow returns a boolean if a field has been set.

### GetCategories

`func (o *AuditPostRequestAdditionalOptions) GetCategories() []string`

GetCategories returns the Categories field if non-nil, zero value otherwise.

### GetCategoriesOk

`func (o *AuditPostRequestAdditionalOptions) GetCategoriesOk() (*[]string, bool)`

GetCategoriesOk returns a tuple with the Categories field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCategories

`func (o *AuditPostRequestAdditionalOptions) SetCategories(v []string)`

SetCategories sets Categories field to given value.

### HasCategories

`func (o *AuditPostRequestAdditionalOptions) HasCategories() bool`

HasCategories returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


