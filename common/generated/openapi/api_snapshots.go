/*
IoEngine RESTful API

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: v0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)


// SnapshotsAPIService SnapshotsAPI service
type SnapshotsAPIService service

type ApiDelSnapshotRequest struct {
	ctx context.Context
	ApiService *SnapshotsAPIService
	snapshotId string
}

func (r ApiDelSnapshotRequest) Execute() (*http.Response, error) {
	return r.ApiService.DelSnapshotExecute(r)
}

/*
DelSnapshot Method for DelSnapshot

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param snapshotId
 @return ApiDelSnapshotRequest
*/
func (a *SnapshotsAPIService) DelSnapshot(ctx context.Context, snapshotId string) ApiDelSnapshotRequest {
	return ApiDelSnapshotRequest{
		ApiService: a,
		ctx: ctx,
		snapshotId: snapshotId,
	}
}

// Execute executes the request
func (a *SnapshotsAPIService) DelSnapshotExecute(r ApiDelSnapshotRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodDelete
		localVarPostBody     interface{}
		formFiles            []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "SnapshotsAPIService.DelSnapshot")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/volumes/snapshots/{snapshot_id}"
	localVarPath = strings.Replace(localVarPath, "{"+"snapshot_id"+"}", url.PathEscape(parameterValueToString(r.snapshotId, "snapshotId")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		if localVarHTTPResponse.StatusCode >= 400 && localVarHTTPResponse.StatusCode < 500 {
			var v RestJsonError
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarHTTPResponse, newErr
			}
					newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
					newErr.model = v
			return localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode >= 500 {
			var v RestJsonError
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarHTTPResponse, newErr
			}
					newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
					newErr.model = v
		}
		return localVarHTTPResponse, newErr
	}

	return localVarHTTPResponse, nil
}

type ApiDelVolumeSnapshotRequest struct {
	ctx context.Context
	ApiService *SnapshotsAPIService
	volumeId string
	snapshotId string
}

func (r ApiDelVolumeSnapshotRequest) Execute() (*http.Response, error) {
	return r.ApiService.DelVolumeSnapshotExecute(r)
}

/*
DelVolumeSnapshot Method for DelVolumeSnapshot

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param volumeId
 @param snapshotId
 @return ApiDelVolumeSnapshotRequest
*/
func (a *SnapshotsAPIService) DelVolumeSnapshot(ctx context.Context, volumeId string, snapshotId string) ApiDelVolumeSnapshotRequest {
	return ApiDelVolumeSnapshotRequest{
		ApiService: a,
		ctx: ctx,
		volumeId: volumeId,
		snapshotId: snapshotId,
	}
}

// Execute executes the request
func (a *SnapshotsAPIService) DelVolumeSnapshotExecute(r ApiDelVolumeSnapshotRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodDelete
		localVarPostBody     interface{}
		formFiles            []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "SnapshotsAPIService.DelVolumeSnapshot")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/volumes/{volume_id}/snapshots/{snapshot_id}"
	localVarPath = strings.Replace(localVarPath, "{"+"volume_id"+"}", url.PathEscape(parameterValueToString(r.volumeId, "volumeId")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"snapshot_id"+"}", url.PathEscape(parameterValueToString(r.snapshotId, "snapshotId")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		if localVarHTTPResponse.StatusCode >= 400 && localVarHTTPResponse.StatusCode < 500 {
			var v RestJsonError
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarHTTPResponse, newErr
			}
					newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
					newErr.model = v
			return localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode >= 500 {
			var v RestJsonError
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarHTTPResponse, newErr
			}
					newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
					newErr.model = v
		}
		return localVarHTTPResponse, newErr
	}

	return localVarHTTPResponse, nil
}

type ApiGetVolumeSnapshotRequest struct {
	ctx context.Context
	ApiService *SnapshotsAPIService
	volumeId string
	snapshotId string
}

func (r ApiGetVolumeSnapshotRequest) Execute() (*VolumeSnapshot, *http.Response, error) {
	return r.ApiService.GetVolumeSnapshotExecute(r)
}

/*
GetVolumeSnapshot Method for GetVolumeSnapshot

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param volumeId
 @param snapshotId
 @return ApiGetVolumeSnapshotRequest
*/
func (a *SnapshotsAPIService) GetVolumeSnapshot(ctx context.Context, volumeId string, snapshotId string) ApiGetVolumeSnapshotRequest {
	return ApiGetVolumeSnapshotRequest{
		ApiService: a,
		ctx: ctx,
		volumeId: volumeId,
		snapshotId: snapshotId,
	}
}

// Execute executes the request
//  @return VolumeSnapshot
func (a *SnapshotsAPIService) GetVolumeSnapshotExecute(r ApiGetVolumeSnapshotRequest) (*VolumeSnapshot, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *VolumeSnapshot
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "SnapshotsAPIService.GetVolumeSnapshot")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/volumes/{volume_id}/snapshots/{snapshot_id}"
	localVarPath = strings.Replace(localVarPath, "{"+"volume_id"+"}", url.PathEscape(parameterValueToString(r.volumeId, "volumeId")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"snapshot_id"+"}", url.PathEscape(parameterValueToString(r.snapshotId, "snapshotId")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		if localVarHTTPResponse.StatusCode >= 400 && localVarHTTPResponse.StatusCode < 500 {
			var v RestJsonError
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
					newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
					newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode >= 500 {
			var v RestJsonError
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
					newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
					newErr.model = v
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ApiGetVolumeSnapshotsRequest struct {
	ctx context.Context
	ApiService *SnapshotsAPIService
	volumeId string
	maxEntries *int32
	startingToken *int32
}

// the maximum number of results to return
func (r ApiGetVolumeSnapshotsRequest) MaxEntries(maxEntries int32) ApiGetVolumeSnapshotsRequest {
	r.maxEntries = &maxEntries
	return r
}

// the offset to start pagination from
func (r ApiGetVolumeSnapshotsRequest) StartingToken(startingToken int32) ApiGetVolumeSnapshotsRequest {
	r.startingToken = &startingToken
	return r
}

func (r ApiGetVolumeSnapshotsRequest) Execute() (*VolumeSnapshots, *http.Response, error) {
	return r.ApiService.GetVolumeSnapshotsExecute(r)
}

/*
GetVolumeSnapshots Method for GetVolumeSnapshots

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param volumeId
 @return ApiGetVolumeSnapshotsRequest
*/
func (a *SnapshotsAPIService) GetVolumeSnapshots(ctx context.Context, volumeId string) ApiGetVolumeSnapshotsRequest {
	return ApiGetVolumeSnapshotsRequest{
		ApiService: a,
		ctx: ctx,
		volumeId: volumeId,
	}
}

// Execute executes the request
//  @return VolumeSnapshots
func (a *SnapshotsAPIService) GetVolumeSnapshotsExecute(r ApiGetVolumeSnapshotsRequest) (*VolumeSnapshots, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *VolumeSnapshots
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "SnapshotsAPIService.GetVolumeSnapshots")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/volumes/{volume_id}/snapshots"
	localVarPath = strings.Replace(localVarPath, "{"+"volume_id"+"}", url.PathEscape(parameterValueToString(r.volumeId, "volumeId")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.maxEntries == nil {
		return localVarReturnValue, nil, reportError("maxEntries is required and must be specified")
	}

	parameterAddToHeaderOrQuery(localVarQueryParams, "max_entries", r.maxEntries, "")
	if r.startingToken != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "starting_token", r.startingToken, "")
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		if localVarHTTPResponse.StatusCode >= 400 && localVarHTTPResponse.StatusCode < 500 {
			var v RestJsonError
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
					newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
					newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode >= 500 {
			var v RestJsonError
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
					newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
					newErr.model = v
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ApiGetVolumesSnapshotRequest struct {
	ctx context.Context
	ApiService *SnapshotsAPIService
	snapshotId string
}

func (r ApiGetVolumesSnapshotRequest) Execute() (*VolumeSnapshot, *http.Response, error) {
	return r.ApiService.GetVolumesSnapshotExecute(r)
}

/*
GetVolumesSnapshot Method for GetVolumesSnapshot

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param snapshotId
 @return ApiGetVolumesSnapshotRequest
*/
func (a *SnapshotsAPIService) GetVolumesSnapshot(ctx context.Context, snapshotId string) ApiGetVolumesSnapshotRequest {
	return ApiGetVolumesSnapshotRequest{
		ApiService: a,
		ctx: ctx,
		snapshotId: snapshotId,
	}
}

// Execute executes the request
//  @return VolumeSnapshot
func (a *SnapshotsAPIService) GetVolumesSnapshotExecute(r ApiGetVolumesSnapshotRequest) (*VolumeSnapshot, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *VolumeSnapshot
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "SnapshotsAPIService.GetVolumesSnapshot")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/volumes/snapshots/{snapshot_id}"
	localVarPath = strings.Replace(localVarPath, "{"+"snapshot_id"+"}", url.PathEscape(parameterValueToString(r.snapshotId, "snapshotId")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		if localVarHTTPResponse.StatusCode >= 400 && localVarHTTPResponse.StatusCode < 500 {
			var v RestJsonError
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
					newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
					newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode >= 500 {
			var v RestJsonError
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
					newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
					newErr.model = v
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ApiGetVolumesSnapshotsRequest struct {
	ctx context.Context
	ApiService *SnapshotsAPIService
	maxEntries *int32
	snapshotId *string
	volumeId *string
	startingToken *int32
}

// the maximum number of results to return
func (r ApiGetVolumesSnapshotsRequest) MaxEntries(maxEntries int32) ApiGetVolumesSnapshotsRequest {
	r.maxEntries = &maxEntries
	return r
}

// The uuid of the snapshot to retrieve.
func (r ApiGetVolumesSnapshotsRequest) SnapshotId(snapshotId string) ApiGetVolumesSnapshotsRequest {
	r.snapshotId = &snapshotId
	return r
}

// The uuid of the snapshots source volume.
func (r ApiGetVolumesSnapshotsRequest) VolumeId(volumeId string) ApiGetVolumesSnapshotsRequest {
	r.volumeId = &volumeId
	return r
}

// the offset to start pagination from
func (r ApiGetVolumesSnapshotsRequest) StartingToken(startingToken int32) ApiGetVolumesSnapshotsRequest {
	r.startingToken = &startingToken
	return r
}

func (r ApiGetVolumesSnapshotsRequest) Execute() (*VolumeSnapshots, *http.Response, error) {
	return r.ApiService.GetVolumesSnapshotsExecute(r)
}

/*
GetVolumesSnapshots Method for GetVolumesSnapshots

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @return ApiGetVolumesSnapshotsRequest
*/
func (a *SnapshotsAPIService) GetVolumesSnapshots(ctx context.Context) ApiGetVolumesSnapshotsRequest {
	return ApiGetVolumesSnapshotsRequest{
		ApiService: a,
		ctx: ctx,
	}
}

// Execute executes the request
//  @return VolumeSnapshots
func (a *SnapshotsAPIService) GetVolumesSnapshotsExecute(r ApiGetVolumesSnapshotsRequest) (*VolumeSnapshots, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *VolumeSnapshots
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "SnapshotsAPIService.GetVolumesSnapshots")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/volumes/snapshots"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.maxEntries == nil {
		return localVarReturnValue, nil, reportError("maxEntries is required and must be specified")
	}

	if r.snapshotId != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "snapshot_id", r.snapshotId, "")
	}
	if r.volumeId != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "volume_id", r.volumeId, "")
	}
	parameterAddToHeaderOrQuery(localVarQueryParams, "max_entries", r.maxEntries, "")
	if r.startingToken != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "starting_token", r.startingToken, "")
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		if localVarHTTPResponse.StatusCode >= 400 && localVarHTTPResponse.StatusCode < 500 {
			var v RestJsonError
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
					newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
					newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode >= 500 {
			var v RestJsonError
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
					newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
					newErr.model = v
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ApiPutVolumeSnapshotRequest struct {
	ctx context.Context
	ApiService *SnapshotsAPIService
	volumeId string
	snapshotId string
}

func (r ApiPutVolumeSnapshotRequest) Execute() (*VolumeSnapshot, *http.Response, error) {
	return r.ApiService.PutVolumeSnapshotExecute(r)
}

/*
PutVolumeSnapshot Method for PutVolumeSnapshot

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param volumeId
 @param snapshotId
 @return ApiPutVolumeSnapshotRequest
*/
func (a *SnapshotsAPIService) PutVolumeSnapshot(ctx context.Context, volumeId string, snapshotId string) ApiPutVolumeSnapshotRequest {
	return ApiPutVolumeSnapshotRequest{
		ApiService: a,
		ctx: ctx,
		volumeId: volumeId,
		snapshotId: snapshotId,
	}
}

// Execute executes the request
//  @return VolumeSnapshot
func (a *SnapshotsAPIService) PutVolumeSnapshotExecute(r ApiPutVolumeSnapshotRequest) (*VolumeSnapshot, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPut
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *VolumeSnapshot
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "SnapshotsAPIService.PutVolumeSnapshot")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/volumes/{volume_id}/snapshots/{snapshot_id}"
	localVarPath = strings.Replace(localVarPath, "{"+"volume_id"+"}", url.PathEscape(parameterValueToString(r.volumeId, "volumeId")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"snapshot_id"+"}", url.PathEscape(parameterValueToString(r.snapshotId, "snapshotId")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		if localVarHTTPResponse.StatusCode >= 400 && localVarHTTPResponse.StatusCode < 500 {
			var v RestJsonError
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
					newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
					newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode >= 500 {
			var v RestJsonError
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
					newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
					newErr.model = v
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}
