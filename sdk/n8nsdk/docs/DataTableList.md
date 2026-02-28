# DataTableList

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | Pointer to [**[]DataTable**](DataTable.md) |  | [optional] 
**NextCursor** | Pointer to **NullableString** | Paginate through data tables by setting the cursor parameter to a nextCursor attribute returned by a previous request. Default value fetches the first \&quot;page\&quot; of the collection. | [optional] 

## Methods

### NewDataTableList

`func NewDataTableList() *DataTableList`

NewDataTableList instantiates a new DataTableList object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDataTableListWithDefaults

`func NewDataTableListWithDefaults() *DataTableList`

NewDataTableListWithDefaults instantiates a new DataTableList object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *DataTableList) GetData() []DataTable`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *DataTableList) GetDataOk() (*[]DataTable, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *DataTableList) SetData(v []DataTable)`

SetData sets Data field to given value.

### HasData

`func (o *DataTableList) HasData() bool`

HasData returns a boolean if a field has been set.

### GetNextCursor

`func (o *DataTableList) GetNextCursor() string`

GetNextCursor returns the NextCursor field if non-nil, zero value otherwise.

### GetNextCursorOk

`func (o *DataTableList) GetNextCursorOk() (*string, bool)`

GetNextCursorOk returns a tuple with the NextCursor field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNextCursor

`func (o *DataTableList) SetNextCursor(v string)`

SetNextCursor sets NextCursor field to given value.

### HasNextCursor

`func (o *DataTableList) HasNextCursor() bool`

HasNextCursor returns a boolean if a field has been set.

### SetNextCursorNil

`func (o *DataTableList) SetNextCursorNil(b bool)`

 SetNextCursorNil sets the value for NextCursor to be an explicit nil

### UnsetNextCursor
`func (o *DataTableList) UnsetNextCursor()`

UnsetNextCursor ensures that no value is present for NextCursor, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


