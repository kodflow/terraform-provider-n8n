# WorkflowSettings

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**SaveExecutionProgress** | Pointer to **bool** |  | [optional] 
**SaveManualExecutions** | Pointer to **bool** |  | [optional] 
**SaveDataErrorExecution** | Pointer to **string** |  | [optional] 
**SaveDataSuccessExecution** | Pointer to **string** |  | [optional] 
**ExecutionTimeout** | Pointer to **float32** |  | [optional] 
**ErrorWorkflow** | Pointer to **string** | The ID of the workflow that contains the error trigger node. | [optional] 
**Timezone** | Pointer to **string** |  | [optional] 
**ExecutionOrder** | Pointer to **string** |  | [optional] 
**CallerPolicy** | Pointer to **string** | Controls which workflows are allowed to call this workflow using the Execute Workflow node.  Available options: - &#x60;any&#x60;: Any workflow can call this workflow (no restrictions) - &#x60;none&#x60;: No other workflows can call this workflow (completely blocked) - &#x60;workflowsFromSameOwner&#x60; (default): Only workflows owned by the same project can call this workflow   * For personal projects: Only workflows created by the same user   * For team projects: Only workflows within the same team project - &#x60;workflowsFromAList&#x60;: Only specific workflows listed in the &#x60;callerIds&#x60; field can call this workflow   * Requires the &#x60;callerIds&#x60; field to specify which workflow IDs are allowed   * See &#x60;callerIds&#x60; field documentation for usage  | [optional] [default to "workflowsFromSameOwner"]
**CallerIds** | Pointer to **string** | Comma-separated list of workflow IDs allowed to call this workflow (only used with workflowsFromAList policy) | [optional] 
**TimeSavedPerExecution** | Pointer to **float32** | Estimated time saved per execution in minutes | [optional] 
**AvailableInMCP** | Pointer to **bool** | Controls whether this workflow is accessible via the Model Context Protocol (MCP).  When enabled, this workflow can be called by MCP clients (AI assistants and other tools that support MCP). This allows external AI tools to discover and execute this workflow as part of their capabilities.  Requirements for enabling MCP access: - The workflow must be active (not deactivated) - The workflow must contain at least one active Webhook node - Only webhook-triggered workflows can be exposed via MCP  Security note: When a workflow is available in MCP, it can be discovered and executed by any MCP client that has the appropriate API credentials for your n8n instance.  | [optional] [default to false]

## Methods

### NewWorkflowSettings

`func NewWorkflowSettings() *WorkflowSettings`

NewWorkflowSettings instantiates a new WorkflowSettings object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowSettingsWithDefaults

`func NewWorkflowSettingsWithDefaults() *WorkflowSettings`

NewWorkflowSettingsWithDefaults instantiates a new WorkflowSettings object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSaveExecutionProgress

`func (o *WorkflowSettings) GetSaveExecutionProgress() bool`

GetSaveExecutionProgress returns the SaveExecutionProgress field if non-nil, zero value otherwise.

### GetSaveExecutionProgressOk

`func (o *WorkflowSettings) GetSaveExecutionProgressOk() (*bool, bool)`

GetSaveExecutionProgressOk returns a tuple with the SaveExecutionProgress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSaveExecutionProgress

`func (o *WorkflowSettings) SetSaveExecutionProgress(v bool)`

SetSaveExecutionProgress sets SaveExecutionProgress field to given value.

### HasSaveExecutionProgress

`func (o *WorkflowSettings) HasSaveExecutionProgress() bool`

HasSaveExecutionProgress returns a boolean if a field has been set.

### GetSaveManualExecutions

`func (o *WorkflowSettings) GetSaveManualExecutions() bool`

GetSaveManualExecutions returns the SaveManualExecutions field if non-nil, zero value otherwise.

### GetSaveManualExecutionsOk

`func (o *WorkflowSettings) GetSaveManualExecutionsOk() (*bool, bool)`

GetSaveManualExecutionsOk returns a tuple with the SaveManualExecutions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSaveManualExecutions

`func (o *WorkflowSettings) SetSaveManualExecutions(v bool)`

SetSaveManualExecutions sets SaveManualExecutions field to given value.

### HasSaveManualExecutions

`func (o *WorkflowSettings) HasSaveManualExecutions() bool`

HasSaveManualExecutions returns a boolean if a field has been set.

### GetSaveDataErrorExecution

`func (o *WorkflowSettings) GetSaveDataErrorExecution() string`

GetSaveDataErrorExecution returns the SaveDataErrorExecution field if non-nil, zero value otherwise.

### GetSaveDataErrorExecutionOk

`func (o *WorkflowSettings) GetSaveDataErrorExecutionOk() (*string, bool)`

GetSaveDataErrorExecutionOk returns a tuple with the SaveDataErrorExecution field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSaveDataErrorExecution

`func (o *WorkflowSettings) SetSaveDataErrorExecution(v string)`

SetSaveDataErrorExecution sets SaveDataErrorExecution field to given value.

### HasSaveDataErrorExecution

`func (o *WorkflowSettings) HasSaveDataErrorExecution() bool`

HasSaveDataErrorExecution returns a boolean if a field has been set.

### GetSaveDataSuccessExecution

`func (o *WorkflowSettings) GetSaveDataSuccessExecution() string`

GetSaveDataSuccessExecution returns the SaveDataSuccessExecution field if non-nil, zero value otherwise.

### GetSaveDataSuccessExecutionOk

`func (o *WorkflowSettings) GetSaveDataSuccessExecutionOk() (*string, bool)`

GetSaveDataSuccessExecutionOk returns a tuple with the SaveDataSuccessExecution field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSaveDataSuccessExecution

`func (o *WorkflowSettings) SetSaveDataSuccessExecution(v string)`

SetSaveDataSuccessExecution sets SaveDataSuccessExecution field to given value.

### HasSaveDataSuccessExecution

`func (o *WorkflowSettings) HasSaveDataSuccessExecution() bool`

HasSaveDataSuccessExecution returns a boolean if a field has been set.

### GetExecutionTimeout

`func (o *WorkflowSettings) GetExecutionTimeout() float32`

GetExecutionTimeout returns the ExecutionTimeout field if non-nil, zero value otherwise.

### GetExecutionTimeoutOk

`func (o *WorkflowSettings) GetExecutionTimeoutOk() (*float32, bool)`

GetExecutionTimeoutOk returns a tuple with the ExecutionTimeout field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExecutionTimeout

`func (o *WorkflowSettings) SetExecutionTimeout(v float32)`

SetExecutionTimeout sets ExecutionTimeout field to given value.

### HasExecutionTimeout

`func (o *WorkflowSettings) HasExecutionTimeout() bool`

HasExecutionTimeout returns a boolean if a field has been set.

### GetErrorWorkflow

`func (o *WorkflowSettings) GetErrorWorkflow() string`

GetErrorWorkflow returns the ErrorWorkflow field if non-nil, zero value otherwise.

### GetErrorWorkflowOk

`func (o *WorkflowSettings) GetErrorWorkflowOk() (*string, bool)`

GetErrorWorkflowOk returns a tuple with the ErrorWorkflow field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErrorWorkflow

`func (o *WorkflowSettings) SetErrorWorkflow(v string)`

SetErrorWorkflow sets ErrorWorkflow field to given value.

### HasErrorWorkflow

`func (o *WorkflowSettings) HasErrorWorkflow() bool`

HasErrorWorkflow returns a boolean if a field has been set.

### GetTimezone

`func (o *WorkflowSettings) GetTimezone() string`

GetTimezone returns the Timezone field if non-nil, zero value otherwise.

### GetTimezoneOk

`func (o *WorkflowSettings) GetTimezoneOk() (*string, bool)`

GetTimezoneOk returns a tuple with the Timezone field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimezone

`func (o *WorkflowSettings) SetTimezone(v string)`

SetTimezone sets Timezone field to given value.

### HasTimezone

`func (o *WorkflowSettings) HasTimezone() bool`

HasTimezone returns a boolean if a field has been set.

### GetExecutionOrder

`func (o *WorkflowSettings) GetExecutionOrder() string`

GetExecutionOrder returns the ExecutionOrder field if non-nil, zero value otherwise.

### GetExecutionOrderOk

`func (o *WorkflowSettings) GetExecutionOrderOk() (*string, bool)`

GetExecutionOrderOk returns a tuple with the ExecutionOrder field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExecutionOrder

`func (o *WorkflowSettings) SetExecutionOrder(v string)`

SetExecutionOrder sets ExecutionOrder field to given value.

### HasExecutionOrder

`func (o *WorkflowSettings) HasExecutionOrder() bool`

HasExecutionOrder returns a boolean if a field has been set.

### GetCallerPolicy

`func (o *WorkflowSettings) GetCallerPolicy() string`

GetCallerPolicy returns the CallerPolicy field if non-nil, zero value otherwise.

### GetCallerPolicyOk

`func (o *WorkflowSettings) GetCallerPolicyOk() (*string, bool)`

GetCallerPolicyOk returns a tuple with the CallerPolicy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCallerPolicy

`func (o *WorkflowSettings) SetCallerPolicy(v string)`

SetCallerPolicy sets CallerPolicy field to given value.

### HasCallerPolicy

`func (o *WorkflowSettings) HasCallerPolicy() bool`

HasCallerPolicy returns a boolean if a field has been set.

### GetCallerIds

`func (o *WorkflowSettings) GetCallerIds() string`

GetCallerIds returns the CallerIds field if non-nil, zero value otherwise.

### GetCallerIdsOk

`func (o *WorkflowSettings) GetCallerIdsOk() (*string, bool)`

GetCallerIdsOk returns a tuple with the CallerIds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCallerIds

`func (o *WorkflowSettings) SetCallerIds(v string)`

SetCallerIds sets CallerIds field to given value.

### HasCallerIds

`func (o *WorkflowSettings) HasCallerIds() bool`

HasCallerIds returns a boolean if a field has been set.

### GetTimeSavedPerExecution

`func (o *WorkflowSettings) GetTimeSavedPerExecution() float32`

GetTimeSavedPerExecution returns the TimeSavedPerExecution field if non-nil, zero value otherwise.

### GetTimeSavedPerExecutionOk

`func (o *WorkflowSettings) GetTimeSavedPerExecutionOk() (*float32, bool)`

GetTimeSavedPerExecutionOk returns a tuple with the TimeSavedPerExecution field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimeSavedPerExecution

`func (o *WorkflowSettings) SetTimeSavedPerExecution(v float32)`

SetTimeSavedPerExecution sets TimeSavedPerExecution field to given value.

### HasTimeSavedPerExecution

`func (o *WorkflowSettings) HasTimeSavedPerExecution() bool`

HasTimeSavedPerExecution returns a boolean if a field has been set.

### GetAvailableInMCP

`func (o *WorkflowSettings) GetAvailableInMCP() bool`

GetAvailableInMCP returns the AvailableInMCP field if non-nil, zero value otherwise.

### GetAvailableInMCPOk

`func (o *WorkflowSettings) GetAvailableInMCPOk() (*bool, bool)`

GetAvailableInMCPOk returns a tuple with the AvailableInMCP field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAvailableInMCP

`func (o *WorkflowSettings) SetAvailableInMCP(v bool)`

SetAvailableInMCP sets AvailableInMCP field to given value.

### HasAvailableInMCP

`func (o *WorkflowSettings) HasAvailableInMCP() bool`

HasAvailableInMCP returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


