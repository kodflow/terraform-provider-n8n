# Node

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**WebhookId** | Pointer to **string** |  | [optional] 
**Disabled** | Pointer to **bool** |  | [optional] 
**NotesInFlow** | Pointer to **bool** |  | [optional] 
**Notes** | Pointer to **string** |  | [optional] 
**Type** | Pointer to **string** |  | [optional] 
**TypeVersion** | Pointer to **float32** |  | [optional] 
**ExecuteOnce** | Pointer to **bool** |  | [optional] 
**AlwaysOutputData** | Pointer to **bool** |  | [optional] 
**RetryOnFail** | Pointer to **bool** |  | [optional] 
**MaxTries** | Pointer to **float32** |  | [optional] 
**WaitBetweenTries** | Pointer to **float32** |  | [optional] 
**ContinueOnFail** | Pointer to **bool** | use onError instead | [optional] 
**OnError** | Pointer to **string** |  | [optional] 
**Position** | Pointer to **[]float32** |  | [optional] 
**Parameters** | Pointer to **map[string]interface{}** |  | [optional] 
**Credentials** | Pointer to **map[string]interface{}** |  | [optional] 
**CreatedAt** | Pointer to **time.Time** |  | [optional] [readonly] 
**UpdatedAt** | Pointer to **time.Time** |  | [optional] [readonly] 

## Methods

### NewNode

`func NewNode() *Node`

NewNode instantiates a new Node object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewNodeWithDefaults

`func NewNodeWithDefaults() *Node`

NewNodeWithDefaults instantiates a new Node object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *Node) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Node) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Node) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *Node) HasId() bool`

HasId returns a boolean if a field has been set.

### GetName

`func (o *Node) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *Node) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *Node) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *Node) HasName() bool`

HasName returns a boolean if a field has been set.

### GetWebhookId

`func (o *Node) GetWebhookId() string`

GetWebhookId returns the WebhookId field if non-nil, zero value otherwise.

### GetWebhookIdOk

`func (o *Node) GetWebhookIdOk() (*string, bool)`

GetWebhookIdOk returns a tuple with the WebhookId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWebhookId

`func (o *Node) SetWebhookId(v string)`

SetWebhookId sets WebhookId field to given value.

### HasWebhookId

`func (o *Node) HasWebhookId() bool`

HasWebhookId returns a boolean if a field has been set.

### GetDisabled

`func (o *Node) GetDisabled() bool`

GetDisabled returns the Disabled field if non-nil, zero value otherwise.

### GetDisabledOk

`func (o *Node) GetDisabledOk() (*bool, bool)`

GetDisabledOk returns a tuple with the Disabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDisabled

`func (o *Node) SetDisabled(v bool)`

SetDisabled sets Disabled field to given value.

### HasDisabled

`func (o *Node) HasDisabled() bool`

HasDisabled returns a boolean if a field has been set.

### GetNotesInFlow

`func (o *Node) GetNotesInFlow() bool`

GetNotesInFlow returns the NotesInFlow field if non-nil, zero value otherwise.

### GetNotesInFlowOk

`func (o *Node) GetNotesInFlowOk() (*bool, bool)`

GetNotesInFlowOk returns a tuple with the NotesInFlow field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNotesInFlow

`func (o *Node) SetNotesInFlow(v bool)`

SetNotesInFlow sets NotesInFlow field to given value.

### HasNotesInFlow

`func (o *Node) HasNotesInFlow() bool`

HasNotesInFlow returns a boolean if a field has been set.

### GetNotes

`func (o *Node) GetNotes() string`

GetNotes returns the Notes field if non-nil, zero value otherwise.

### GetNotesOk

`func (o *Node) GetNotesOk() (*string, bool)`

GetNotesOk returns a tuple with the Notes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNotes

`func (o *Node) SetNotes(v string)`

SetNotes sets Notes field to given value.

### HasNotes

`func (o *Node) HasNotes() bool`

HasNotes returns a boolean if a field has been set.

### GetType

`func (o *Node) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *Node) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *Node) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *Node) HasType() bool`

HasType returns a boolean if a field has been set.

### GetTypeVersion

`func (o *Node) GetTypeVersion() float32`

GetTypeVersion returns the TypeVersion field if non-nil, zero value otherwise.

### GetTypeVersionOk

`func (o *Node) GetTypeVersionOk() (*float32, bool)`

GetTypeVersionOk returns a tuple with the TypeVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTypeVersion

`func (o *Node) SetTypeVersion(v float32)`

SetTypeVersion sets TypeVersion field to given value.

### HasTypeVersion

`func (o *Node) HasTypeVersion() bool`

HasTypeVersion returns a boolean if a field has been set.

### GetExecuteOnce

`func (o *Node) GetExecuteOnce() bool`

GetExecuteOnce returns the ExecuteOnce field if non-nil, zero value otherwise.

### GetExecuteOnceOk

`func (o *Node) GetExecuteOnceOk() (*bool, bool)`

GetExecuteOnceOk returns a tuple with the ExecuteOnce field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExecuteOnce

`func (o *Node) SetExecuteOnce(v bool)`

SetExecuteOnce sets ExecuteOnce field to given value.

### HasExecuteOnce

`func (o *Node) HasExecuteOnce() bool`

HasExecuteOnce returns a boolean if a field has been set.

### GetAlwaysOutputData

`func (o *Node) GetAlwaysOutputData() bool`

GetAlwaysOutputData returns the AlwaysOutputData field if non-nil, zero value otherwise.

### GetAlwaysOutputDataOk

`func (o *Node) GetAlwaysOutputDataOk() (*bool, bool)`

GetAlwaysOutputDataOk returns a tuple with the AlwaysOutputData field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAlwaysOutputData

`func (o *Node) SetAlwaysOutputData(v bool)`

SetAlwaysOutputData sets AlwaysOutputData field to given value.

### HasAlwaysOutputData

`func (o *Node) HasAlwaysOutputData() bool`

HasAlwaysOutputData returns a boolean if a field has been set.

### GetRetryOnFail

`func (o *Node) GetRetryOnFail() bool`

GetRetryOnFail returns the RetryOnFail field if non-nil, zero value otherwise.

### GetRetryOnFailOk

`func (o *Node) GetRetryOnFailOk() (*bool, bool)`

GetRetryOnFailOk returns a tuple with the RetryOnFail field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRetryOnFail

`func (o *Node) SetRetryOnFail(v bool)`

SetRetryOnFail sets RetryOnFail field to given value.

### HasRetryOnFail

`func (o *Node) HasRetryOnFail() bool`

HasRetryOnFail returns a boolean if a field has been set.

### GetMaxTries

`func (o *Node) GetMaxTries() float32`

GetMaxTries returns the MaxTries field if non-nil, zero value otherwise.

### GetMaxTriesOk

`func (o *Node) GetMaxTriesOk() (*float32, bool)`

GetMaxTriesOk returns a tuple with the MaxTries field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaxTries

`func (o *Node) SetMaxTries(v float32)`

SetMaxTries sets MaxTries field to given value.

### HasMaxTries

`func (o *Node) HasMaxTries() bool`

HasMaxTries returns a boolean if a field has been set.

### GetWaitBetweenTries

`func (o *Node) GetWaitBetweenTries() float32`

GetWaitBetweenTries returns the WaitBetweenTries field if non-nil, zero value otherwise.

### GetWaitBetweenTriesOk

`func (o *Node) GetWaitBetweenTriesOk() (*float32, bool)`

GetWaitBetweenTriesOk returns a tuple with the WaitBetweenTries field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWaitBetweenTries

`func (o *Node) SetWaitBetweenTries(v float32)`

SetWaitBetweenTries sets WaitBetweenTries field to given value.

### HasWaitBetweenTries

`func (o *Node) HasWaitBetweenTries() bool`

HasWaitBetweenTries returns a boolean if a field has been set.

### GetContinueOnFail

`func (o *Node) GetContinueOnFail() bool`

GetContinueOnFail returns the ContinueOnFail field if non-nil, zero value otherwise.

### GetContinueOnFailOk

`func (o *Node) GetContinueOnFailOk() (*bool, bool)`

GetContinueOnFailOk returns a tuple with the ContinueOnFail field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContinueOnFail

`func (o *Node) SetContinueOnFail(v bool)`

SetContinueOnFail sets ContinueOnFail field to given value.

### HasContinueOnFail

`func (o *Node) HasContinueOnFail() bool`

HasContinueOnFail returns a boolean if a field has been set.

### GetOnError

`func (o *Node) GetOnError() string`

GetOnError returns the OnError field if non-nil, zero value otherwise.

### GetOnErrorOk

`func (o *Node) GetOnErrorOk() (*string, bool)`

GetOnErrorOk returns a tuple with the OnError field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOnError

`func (o *Node) SetOnError(v string)`

SetOnError sets OnError field to given value.

### HasOnError

`func (o *Node) HasOnError() bool`

HasOnError returns a boolean if a field has been set.

### GetPosition

`func (o *Node) GetPosition() []float32`

GetPosition returns the Position field if non-nil, zero value otherwise.

### GetPositionOk

`func (o *Node) GetPositionOk() (*[]float32, bool)`

GetPositionOk returns a tuple with the Position field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPosition

`func (o *Node) SetPosition(v []float32)`

SetPosition sets Position field to given value.

### HasPosition

`func (o *Node) HasPosition() bool`

HasPosition returns a boolean if a field has been set.

### GetParameters

`func (o *Node) GetParameters() map[string]interface{}`

GetParameters returns the Parameters field if non-nil, zero value otherwise.

### GetParametersOk

`func (o *Node) GetParametersOk() (*map[string]interface{}, bool)`

GetParametersOk returns a tuple with the Parameters field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetParameters

`func (o *Node) SetParameters(v map[string]interface{})`

SetParameters sets Parameters field to given value.

### HasParameters

`func (o *Node) HasParameters() bool`

HasParameters returns a boolean if a field has been set.

### GetCredentials

`func (o *Node) GetCredentials() map[string]interface{}`

GetCredentials returns the Credentials field if non-nil, zero value otherwise.

### GetCredentialsOk

`func (o *Node) GetCredentialsOk() (*map[string]interface{}, bool)`

GetCredentialsOk returns a tuple with the Credentials field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCredentials

`func (o *Node) SetCredentials(v map[string]interface{})`

SetCredentials sets Credentials field to given value.

### HasCredentials

`func (o *Node) HasCredentials() bool`

HasCredentials returns a boolean if a field has been set.

### GetCreatedAt

`func (o *Node) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *Node) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *Node) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *Node) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetUpdatedAt

`func (o *Node) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *Node) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *Node) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *Node) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


