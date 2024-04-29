# {{classname}}

All URIs are relative to *https://rpc.cosmos.directory/cosmoshub*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DialPeers**](UnsafeApi.md#DialPeers) | **Get** /dial_peers | Add Peers/Persistent Peers (unsafe)
[**DialSeeds**](UnsafeApi.md#DialSeeds) | **Get** /dial_seeds | Dial Seeds (Unsafe)

# **DialPeers**
> DialResp DialPeers(ctx, optional)
Add Peers/Persistent Peers (unsafe)

Set a persistent peer, this route in under unsafe, and has to manually enabled to use.  **Example:** curl 'localhost:26657/dial_peers?peers=\\[\"f9baeaa15fedf5e1ef7448dd60f46c01f1a9e9c4@1.2.3.4:26656\",\"0491d373a8e0fcf1023aaf18c51d6a1d0d4f31bd@5.6.7.8:26656\"\\]&persistent=false' 

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***UnsafeApiDialPeersOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a UnsafeApiDialPeersOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **persistent** | **optional.Bool**| Have the peers you are dialing be persistent | 
 **unconditional** | **optional.Bool**| Have the peers you are dialing be unconditional | 
 **private** | **optional.Bool**| Have the peers you are dialing be private | 
 **peers** | [**optional.Interface of []string**](string.md)| array of peers to dial | 

### Return type

[**DialResp**](dialResp.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DialSeeds**
> DialResp DialSeeds(ctx, optional)
Dial Seeds (Unsafe)

Dial a peer, this route in under unsafe, and has to manually enabled to use    **Example:** curl 'localhost:26657/dial_seeds?seeds=\\[\"f9baeaa15fedf5e1ef7448dd60f46c01f1a9e9c4@1.2.3.4:26656\",\"0491d373a8e0fcf1023aaf18c51d6a1d0d4f31bd@5.6.7.8:26656\"\\]' 

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***UnsafeApiDialSeedsOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a UnsafeApiDialSeedsOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **peers** | [**optional.Interface of []string**](string.md)| list of seed nodes to dial | 

### Return type

[**DialResp**](dialResp.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

