# {{classname}}

All URIs are relative to *https://rpc.cosmos.directory/cosmoshub*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AbciInfo**](ABCIApi.md#AbciInfo) | **Get** /abci_info | Get info about the application.
[**AbciQuery**](ABCIApi.md#AbciQuery) | **Get** /abci_query | Query the application for some information.

# **AbciInfo**
> AbciInfoResponse AbciInfo(ctx, )
Get info about the application.

Get info about the application.  Upon success, the `Cache-Control` header will be set with the default maximum age. 

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**AbciInfoResponse**](ABCIInfoResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **AbciQuery**
> AbciQueryResponse AbciQuery(ctx, path, data, optional)
Query the application for some information.

Query the application for some information. 

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **path** | **string**| Path to the data (\&quot;/a/b/c\&quot;) | 
  **data** | **string**| Data | 
 **optional** | ***ABCIApiAbciQueryOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ABCIApiAbciQueryOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **height** | **optional.Int32**| Height (0 means latest) | [default to 0]
 **prove** | **optional.Bool**| Include proofs of the transactions inclusion in the block | [default to false]

### Return type

[**AbciQueryResponse**](ABCIQueryResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

