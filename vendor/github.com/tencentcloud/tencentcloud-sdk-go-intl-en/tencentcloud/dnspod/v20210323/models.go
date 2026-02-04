// Copyright (c) 2017-2025 Tencent. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v20210323

import (
    tcerr "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common/errors"
    tchttp "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common/http"
    "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common/json"
)

type AddRecordBatch struct {
	// Record type. For more information, see the `DescribeRecordType` API.
	RecordType *string `json:"RecordType,omitnil,omitempty" name:"RecordType"`

	// Record value.
	Value *string `json:"Value,omitnil,omitempty" name:"Value"`

	// Subdomain (host record), which is `@` by default.
	SubDomain *string `json:"SubDomain,omitnil,omitempty" name:"SubDomain"`

	// Split zone of the DNS record. For more information, see the `DescribeRecordLineList` API. If neither `RecordLine` nor `RecordLineId` is specified, the default split zone will be used.
	RecordLine *string `json:"RecordLine,omitnil,omitempty" name:"RecordLine"`

	// Split zone ID of the DNS record. If both `RecordLine` and `RecordLineId` are specified, `RecordLineId` will be used.
	RecordLineId *string `json:"RecordLineId,omitnil,omitempty" name:"RecordLineId"`

	// The record weight (not supported).
	Weight *uint64 `json:"Weight,omitnil,omitempty" name:"Weight"`

	// MX record value. It is `0` by default for non-MX records and required for MX records.
	MX *uint64 `json:"MX,omitnil,omitempty" name:"MX"`

	// TTL value of the record, which is `600` by default.
	TTL *uint64 `json:"TTL,omitnil,omitempty" name:"TTL"`

	// Record status (not supported). Valid values: `0` (disabled); `1` (enabled). Default value: `1`.
	Enabled *uint64 `json:"Enabled,omitnil,omitempty" name:"Enabled"`

	// Record remarks (not supported).
	Remark *string `json:"Remark,omitnil,omitempty" name:"Remark"`
}

type BatchRecordInfo struct {
	// Record ID
	// Note: This field may return null, indicating that no valid values can be obtained.
	RecordId *uint64 `json:"RecordId,omitnil,omitempty" name:"RecordId"`

	// Subdomain (host record).
	// Note: This field may return null, indicating that no valid values can be obtained.
	SubDomain *string `json:"SubDomain,omitnil,omitempty" name:"SubDomain"`

	// Record type. For more information, see the `DescribeRecordType` API.
	// Note: This field may return null, indicating that no valid values can be obtained.
	RecordType *string `json:"RecordType,omitnil,omitempty" name:"RecordType"`

	// Split zone of the DNS record. For more information, see the `DescribeRecordLineList` API.
	// Note: This field may return null, indicating that no valid values can be obtained.
	RecordLine *string `json:"RecordLine,omitnil,omitempty" name:"RecordLine"`

	// Record value
	// Note: This field may return null, indicating that no valid values can be obtained.
	Value *string `json:"Value,omitnil,omitempty" name:"Value"`

	// TTL value of the record
	// Note: This field may return null, indicating that no valid values can be obtained.
	TTL *uint64 `json:"TTL,omitnil,omitempty" name:"TTL"`

	// Record adding status
	// Note: This field may return null, indicating that no valid values can be obtained.
	Status *string `json:"Status,omitnil,omitempty" name:"Status"`

	// Operation type
	// Note: This field may return null, indicating that no valid values can be obtained.
	Operation *string `json:"Operation,omitnil,omitempty" name:"Operation"`

	// Error message
	// Note: This field may return null, indicating that no valid values can be obtained.
	ErrMsg *string `json:"ErrMsg,omitnil,omitempty" name:"ErrMsg"`

	// ID of the record in the list
	Id *uint64 `json:"Id,omitnil,omitempty" name:"Id"`

	// Record status
	// Note: This field may return null, indicating that no valid values can be obtained.
	Enabled *uint64 `json:"Enabled,omitnil,omitempty" name:"Enabled"`

	// MX weight of the record
	// Note: This field may return null, indicating that no valid values can be obtained.
	MX *uint64 `json:"MX,omitnil,omitempty" name:"MX"`

	// The record weight.
	// Note: This field may return null, indicating that no valid values can be obtained.
	Weight *uint64 `json:"Weight,omitnil,omitempty" name:"Weight"`
}

// Predefined struct for user
type CreateDomainAliasRequestParams struct {
	// Domain alias
	DomainAlias *string `json:"DomainAlias,omitnil,omitempty" name:"DomainAlias"`

	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *int64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

type CreateDomainAliasRequest struct {
	*tchttp.BaseRequest
	
	// Domain alias
	DomainAlias *string `json:"DomainAlias,omitnil,omitempty" name:"DomainAlias"`

	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *int64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

func (r *CreateDomainAliasRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateDomainAliasRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "DomainAlias")
	delete(f, "Domain")
	delete(f, "DomainId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CreateDomainAliasRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateDomainAliasResponseParams struct {
	// Domain alias ID
	DomainAliasId *int64 `json:"DomainAliasId,omitnil,omitempty" name:"DomainAliasId"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type CreateDomainAliasResponse struct {
	*tchttp.BaseResponse
	Response *CreateDomainAliasResponseParams `json:"Response"`
}

func (r *CreateDomainAliasResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateDomainAliasResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type CreateDomainBatchDetail struct {
	// See `RecordInfoBatch`.
	// Note: This field may return null, indicating that no valid values can be obtained.
	RecordList []*CreateDomainBatchRecord `json:"RecordList,omitnil,omitempty" name:"RecordList"`

	// Task ID
	Id *uint64 `json:"Id,omitnil,omitempty" name:"Id"`

	// Domain
	// Note: This field may return null, indicating that no valid values can be obtained.
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// Domain level
	// Note: This field may return null, indicating that no valid values can be obtained.
	DomainGrade *string `json:"DomainGrade,omitnil,omitempty" name:"DomainGrade"`

	// Error message
	// Note: This field may return null, indicating that no valid values can be obtained.
	ErrMsg *string `json:"ErrMsg,omitnil,omitempty" name:"ErrMsg"`

	// Task running status
	// Note: This field may return null, indicating that no valid values can be obtained.
	Status *string `json:"Status,omitnil,omitempty" name:"Status"`

	// Operation type
	// Note: This field may return null, indicating that no valid values can be obtained.
	Operation *string `json:"Operation,omitnil,omitempty" name:"Operation"`
}

type CreateDomainBatchRecord struct {
	// Subdomain (host record).
	// Note: This field may return null, indicating that no valid values can be obtained.
	SubDomain *string `json:"SubDomain,omitnil,omitempty" name:"SubDomain"`

	// Record type. For more information, see the `DescribeRecordType` API.
	// Note: This field may return null, indicating that no valid values can be obtained.
	RecordType *string `json:"RecordType,omitnil,omitempty" name:"RecordType"`

	// Split zone of the DNS record. For more information, see the `DescribeRecordLineList` API.
	// Note: This field may return null, indicating that no valid values can be obtained.
	RecordLine *string `json:"RecordLine,omitnil,omitempty" name:"RecordLine"`

	// Record value
	// Note: This field may return null, indicating that no valid values can be obtained.
	Value *string `json:"Value,omitnil,omitempty" name:"Value"`

	// TTL value of the record
	// Note: This field may return null, indicating that no valid values can be obtained.
	TTL *uint64 `json:"TTL,omitnil,omitempty" name:"TTL"`

	// Record adding status
	// Note: This field may return null, indicating that no valid values can be obtained.
	Status *string `json:"Status,omitnil,omitempty" name:"Status"`

	// Operation type
	// Note: This field may return null, indicating that no valid values can be obtained.
	Operation *string `json:"Operation,omitnil,omitempty" name:"Operation"`

	// Error message
	// Note: This field may return null, indicating that no valid values can be obtained.
	ErrMsg *string `json:"ErrMsg,omitnil,omitempty" name:"ErrMsg"`

	// ID of the record in the list
	Id *uint64 `json:"Id,omitnil,omitempty" name:"Id"`
}

// Predefined struct for user
type CreateDomainBatchRequestParams struct {
	// Domain array
	DomainList []*string `json:"DomainList,omitnil,omitempty" name:"DomainList"`

	// Add A records of @ and www for each domain with the record value of the IP. If this parameter is not passed in or is set to an empty string or null, only the domain but not the records will be added.
	RecordValue *string `json:"RecordValue,omitnil,omitempty" name:"RecordValue"`
}

type CreateDomainBatchRequest struct {
	*tchttp.BaseRequest
	
	// Domain array
	DomainList []*string `json:"DomainList,omitnil,omitempty" name:"DomainList"`

	// Add A records of @ and www for each domain with the record value of the IP. If this parameter is not passed in or is set to an empty string or null, only the domain but not the records will be added.
	RecordValue *string `json:"RecordValue,omitnil,omitempty" name:"RecordValue"`
}

func (r *CreateDomainBatchRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateDomainBatchRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "DomainList")
	delete(f, "RecordValue")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CreateDomainBatchRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateDomainBatchResponseParams struct {
	// Information of the bulk added domains
	DetailList []*CreateDomainBatchDetail `json:"DetailList,omitnil,omitempty" name:"DetailList"`

	// Bulk task ID
	JobId *uint64 `json:"JobId,omitnil,omitempty" name:"JobId"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type CreateDomainBatchResponse struct {
	*tchttp.BaseResponse
	Response *CreateDomainBatchResponseParams `json:"Response"`
}

func (r *CreateDomainBatchResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateDomainBatchResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateDomainGroupRequestParams struct {
	// Domain group
	GroupName *string `json:"GroupName,omitnil,omitempty" name:"GroupName"`
}

type CreateDomainGroupRequest struct {
	*tchttp.BaseRequest
	
	// Domain group
	GroupName *string `json:"GroupName,omitnil,omitempty" name:"GroupName"`
}

func (r *CreateDomainGroupRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateDomainGroupRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "GroupName")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CreateDomainGroupRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateDomainGroupResponseParams struct {
	// Domain group ID
	GroupId *int64 `json:"GroupId,omitnil,omitempty" name:"GroupId"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type CreateDomainGroupResponse struct {
	*tchttp.BaseResponse
	Response *CreateDomainGroupResponseParams `json:"Response"`
}

func (r *CreateDomainGroupResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateDomainGroupResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateDomainRequestParams struct {
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// The group ID of the domain. You can view the group information of this domain via the `DescribeDomainGroupList` API.
	GroupId *uint64 `json:"GroupId,omitnil,omitempty" name:"GroupId"`

	// Whether the domain is starred. Valid values: yes, no.
	IsMark *string `json:"IsMark,omitnil,omitempty" name:"IsMark"`
}

type CreateDomainRequest struct {
	*tchttp.BaseRequest
	
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// The group ID of the domain. You can view the group information of this domain via the `DescribeDomainGroupList` API.
	GroupId *uint64 `json:"GroupId,omitnil,omitempty" name:"GroupId"`

	// Whether the domain is starred. Valid values: yes, no.
	IsMark *string `json:"IsMark,omitnil,omitempty" name:"IsMark"`
}

func (r *CreateDomainRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateDomainRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Domain")
	delete(f, "GroupId")
	delete(f, "IsMark")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CreateDomainRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateDomainResponseParams struct {
	// Domain information
	DomainInfo *DomainCreateInfo `json:"DomainInfo,omitnil,omitempty" name:"DomainInfo"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type CreateDomainResponse struct {
	*tchttp.BaseResponse
	Response *CreateDomainResponseParams `json:"Response"`
}

func (r *CreateDomainResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateDomainResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreatePackageOrderRequestParams struct {
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// Valid options for the package version are as follows:
	// DPG_PROFESSIONAL
	// DPG_ENTERPRISE
	// DPG_ULTIMATE
	Grade *string `json:"Grade,omitnil,omitempty" name:"Grade"`
}

type CreatePackageOrderRequest struct {
	*tchttp.BaseRequest
	
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// Valid options for the package version are as follows:
	// DPG_PROFESSIONAL
	// DPG_ENTERPRISE
	// DPG_ULTIMATE
	Grade *string `json:"Grade,omitnil,omitempty" name:"Grade"`
}

func (r *CreatePackageOrderRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreatePackageOrderRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Domain")
	delete(f, "Grade")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CreatePackageOrderRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreatePackageOrderResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type CreatePackageOrderResponse struct {
	*tchttp.BaseResponse
	Response *CreatePackageOrderResponseParams `json:"Response"`
}

func (r *CreatePackageOrderResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreatePackageOrderResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type CreateRecordBatchDetail struct {
	// See `RecordInfoBatch`.
	// Note: This field may return null, indicating that no valid values can be obtained.
	RecordList []*CreateRecordBatchRecord `json:"RecordList,omitnil,omitempty" name:"RecordList"`

	// Task ID
	Id *uint64 `json:"Id,omitnil,omitempty" name:"Id"`

	// Domain
	// Note: This field may return null, indicating that no valid values can be obtained.
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// Domain level
	// Note: This field may return null, indicating that no valid values can be obtained.
	DomainGrade *string `json:"DomainGrade,omitnil,omitempty" name:"DomainGrade"`

	// Error message
	// Note: This field may return null, indicating that no valid values can be obtained.
	ErrMsg *string `json:"ErrMsg,omitnil,omitempty" name:"ErrMsg"`

	// Task running status
	// Note: This field may return null, indicating that no valid values can be obtained.
	Status *string `json:"Status,omitnil,omitempty" name:"Status"`

	// Operation type
	// Note: This field may return null, indicating that no valid values can be obtained.
	Operation *string `json:"Operation,omitnil,omitempty" name:"Operation"`

	// Domain ID
	// Note: This field may return null, indicating that no valid values can be obtained.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

type CreateRecordBatchRecord struct {
	// Subdomain (host record).
	// Note: This field may return null, indicating that no valid values can be obtained.
	SubDomain *string `json:"SubDomain,omitnil,omitempty" name:"SubDomain"`

	// Record type. For more information, see the `DescribeRecordType` API.
	// Note: This field may return null, indicating that no valid values can be obtained.
	RecordType *string `json:"RecordType,omitnil,omitempty" name:"RecordType"`

	// Split zone of the DNS record. For more information, see the `DescribeRecordLineList` API.
	// Note: This field may return null, indicating that no valid values can be obtained.
	RecordLine *string `json:"RecordLine,omitnil,omitempty" name:"RecordLine"`

	// Record value
	// Note: This field may return null, indicating that no valid values can be obtained.
	Value *string `json:"Value,omitnil,omitempty" name:"Value"`

	// TTL value of the record
	// Note: This field may return null, indicating that no valid values can be obtained.
	TTL *uint64 `json:"TTL,omitnil,omitempty" name:"TTL"`

	// Record adding status
	// Note: This field may return null, indicating that no valid values can be obtained.
	Status *string `json:"Status,omitnil,omitempty" name:"Status"`

	// Operation type
	// Note: This field may return null, indicating that no valid values can be obtained.
	Operation *string `json:"Operation,omitnil,omitempty" name:"Operation"`

	// Error message
	// Note: This field may return null, indicating that no valid values can be obtained.
	ErrMsg *string `json:"ErrMsg,omitnil,omitempty" name:"ErrMsg"`

	// ID of the record in the list
	Id *uint64 `json:"Id,omitnil,omitempty" name:"Id"`

	// MX weight of the record
	// Note: This field may return null, indicating that no valid values can be obtained.
	MX *uint64 `json:"MX,omitnil,omitempty" name:"MX"`

	// The record weight.
	// Note: This field may return null, indicating that no valid values can be obtained.
	Weight *uint64 `json:"Weight,omitnil,omitempty" name:"Weight"`
}

// Predefined struct for user
type CreateRecordBatchRequestParams struct {
	// Domain ID. Separate multiple ones by comma.
	DomainIdList []*string `json:"DomainIdList,omitnil,omitempty" name:"DomainIdList"`

	// Record array
	RecordList []*AddRecordBatch `json:"RecordList,omitnil,omitempty" name:"RecordList"`
}

type CreateRecordBatchRequest struct {
	*tchttp.BaseRequest
	
	// Domain ID. Separate multiple ones by comma.
	DomainIdList []*string `json:"DomainIdList,omitnil,omitempty" name:"DomainIdList"`

	// Record array
	RecordList []*AddRecordBatch `json:"RecordList,omitnil,omitempty" name:"RecordList"`
}

func (r *CreateRecordBatchRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateRecordBatchRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "DomainIdList")
	delete(f, "RecordList")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CreateRecordBatchRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateRecordBatchResponseParams struct {
	// Information of the bulk added domains
	DetailList []*CreateRecordBatchDetail `json:"DetailList,omitnil,omitempty" name:"DetailList"`

	// Bulk task ID
	JobId *uint64 `json:"JobId,omitnil,omitempty" name:"JobId"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type CreateRecordBatchResponse struct {
	*tchttp.BaseResponse
	Response *CreateRecordBatchResponseParams `json:"Response"`
}

func (r *CreateRecordBatchResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateRecordBatchResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateRecordGroupRequestParams struct {
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// Group name
	GroupName *string `json:"GroupName,omitnil,omitempty" name:"GroupName"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

type CreateRecordGroupRequest struct {
	*tchttp.BaseRequest
	
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// Group name
	GroupName *string `json:"GroupName,omitnil,omitempty" name:"GroupName"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

func (r *CreateRecordGroupRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateRecordGroupRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Domain")
	delete(f, "GroupName")
	delete(f, "DomainId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CreateRecordGroupRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateRecordGroupResponseParams struct {
	// New group ID
	GroupId *uint64 `json:"GroupId,omitnil,omitempty" name:"GroupId"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type CreateRecordGroupResponse struct {
	*tchttp.BaseResponse
	Response *CreateRecordGroupResponseParams `json:"Response"`
}

func (r *CreateRecordGroupResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateRecordGroupResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateRecordRequestParams struct {
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// Record type, which is obtained through the record type API. The value contains uppercase letters, such as `A`.
	RecordType *string `json:"RecordType,omitnil,omitempty" name:"RecordType"`

	// Record split zone, which is obtained through the record split zone API.
	RecordLine *string `json:"RecordLine,omitnil,omitempty" name:"RecordLine"`

	// Record value, such as `IP : 200.200.200.200`, `CNAME : cname.dnspod.com`, and `MX : mail.dnspod.com`.
	Value *string `json:"Value,omitnil,omitempty" name:"Value"`

	// Domain ID. The `DomainId` parameter has a higher priority than `Domain`. If `DomainId` is passed in, `Domain` will be ignored.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`

	// Host record such as `www`. If it is not passed in, it will be `@` by default.
	SubDomain *string `json:"SubDomain,omitnil,omitempty" name:"SubDomain"`

	// Split zone ID, which is obtained through the record split zone API. The value is a string such as `10=1`. The `RecordLineId` parameter has a higher priority than `RecordLine`. If both of them are passed in, `RecordLineId` will be used first.
	RecordLineId *string `json:"RecordLineId,omitnil,omitempty" name:"RecordLineId"`

	// MX priority, which is required for an MX record and will take effect if the record type is MX. Value range: 1–20.
	MX *uint64 `json:"MX,omitnil,omitempty" name:"MX"`

	// TTL. Value range: 1–604800. The minimum value varies by domain level.
	TTL *uint64 `json:"TTL,omitnil,omitempty" name:"TTL"`

	// Weight information, which is an integer between 0 and 100. It is supported only for enterprise VIP domains. `0` indicates not to pass in this parameter, i.e., not to set the weight.
	Weight *uint64 `json:"Weight,omitnil,omitempty" name:"Weight"`

	// Initial status of the record. Valid values: ENABLE, DISABLE. Default value: ENABLE. If `DISABLE` is passed in, the DNS record won't take effect, and the limit on round-robin DNS won't be verified.
	Status *string `json:"Status,omitnil,omitempty" name:"Status"`
}

type CreateRecordRequest struct {
	*tchttp.BaseRequest
	
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// Record type, which is obtained through the record type API. The value contains uppercase letters, such as `A`.
	RecordType *string `json:"RecordType,omitnil,omitempty" name:"RecordType"`

	// Record split zone, which is obtained through the record split zone API.
	RecordLine *string `json:"RecordLine,omitnil,omitempty" name:"RecordLine"`

	// Record value, such as `IP : 200.200.200.200`, `CNAME : cname.dnspod.com`, and `MX : mail.dnspod.com`.
	Value *string `json:"Value,omitnil,omitempty" name:"Value"`

	// Domain ID. The `DomainId` parameter has a higher priority than `Domain`. If `DomainId` is passed in, `Domain` will be ignored.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`

	// Host record such as `www`. If it is not passed in, it will be `@` by default.
	SubDomain *string `json:"SubDomain,omitnil,omitempty" name:"SubDomain"`

	// Split zone ID, which is obtained through the record split zone API. The value is a string such as `10=1`. The `RecordLineId` parameter has a higher priority than `RecordLine`. If both of them are passed in, `RecordLineId` will be used first.
	RecordLineId *string `json:"RecordLineId,omitnil,omitempty" name:"RecordLineId"`

	// MX priority, which is required for an MX record and will take effect if the record type is MX. Value range: 1–20.
	MX *uint64 `json:"MX,omitnil,omitempty" name:"MX"`

	// TTL. Value range: 1–604800. The minimum value varies by domain level.
	TTL *uint64 `json:"TTL,omitnil,omitempty" name:"TTL"`

	// Weight information, which is an integer between 0 and 100. It is supported only for enterprise VIP domains. `0` indicates not to pass in this parameter, i.e., not to set the weight.
	Weight *uint64 `json:"Weight,omitnil,omitempty" name:"Weight"`

	// Initial status of the record. Valid values: ENABLE, DISABLE. Default value: ENABLE. If `DISABLE` is passed in, the DNS record won't take effect, and the limit on round-robin DNS won't be verified.
	Status *string `json:"Status,omitnil,omitempty" name:"Status"`
}

func (r *CreateRecordRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateRecordRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Domain")
	delete(f, "RecordType")
	delete(f, "RecordLine")
	delete(f, "Value")
	delete(f, "DomainId")
	delete(f, "SubDomain")
	delete(f, "RecordLineId")
	delete(f, "MX")
	delete(f, "TTL")
	delete(f, "Weight")
	delete(f, "Status")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CreateRecordRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateRecordResponseParams struct {
	// Record ID
	RecordId *uint64 `json:"RecordId,omitnil,omitempty" name:"RecordId"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type CreateRecordResponse struct {
	*tchttp.BaseResponse
	Response *CreateRecordResponseParams `json:"Response"`
}

func (r *CreateRecordResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateRecordResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteDomainAliasRequestParams struct {
	// Domain alias ID. You can view all domain aliases and their IDs via the `DescribeDomainAliasList` API.
	DomainAliasId *int64 `json:"DomainAliasId,omitnil,omitempty" name:"DomainAliasId"`

	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *int64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

type DeleteDomainAliasRequest struct {
	*tchttp.BaseRequest
	
	// Domain alias ID. You can view all domain aliases and their IDs via the `DescribeDomainAliasList` API.
	DomainAliasId *int64 `json:"DomainAliasId,omitnil,omitempty" name:"DomainAliasId"`

	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *int64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

func (r *DeleteDomainAliasRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteDomainAliasRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "DomainAliasId")
	delete(f, "Domain")
	delete(f, "DomainId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DeleteDomainAliasRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteDomainAliasResponseParams struct {
	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DeleteDomainAliasResponse struct {
	*tchttp.BaseResponse
	Response *DeleteDomainAliasResponseParams `json:"Response"`
}

func (r *DeleteDomainAliasResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteDomainAliasResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type DeleteDomainBatchDetail struct {
	// The domain ID.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`

	// The domain name.
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// The error message.
	// Note: This field may return null, indicating that no valid values can be obtained.
	Error *string `json:"Error,omitnil,omitempty" name:"Error"`

	// The domain deletion status.
	Status *string `json:"Status,omitnil,omitempty" name:"Status"`

	// The operation.
	Operation *string `json:"Operation,omitnil,omitempty" name:"Operation"`
}

// Predefined struct for user
type DeleteDomainBatchRequestParams struct {
	// The domain array.
	DomainList []*string `json:"DomainList,omitnil,omitempty" name:"DomainList"`
}

type DeleteDomainBatchRequest struct {
	*tchttp.BaseRequest
	
	// The domain array.
	DomainList []*string `json:"DomainList,omitnil,omitempty" name:"DomainList"`
}

func (r *DeleteDomainBatchRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteDomainBatchRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "DomainList")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DeleteDomainBatchRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteDomainBatchResponseParams struct {
	// The task ID.
	JobId *uint64 `json:"JobId,omitnil,omitempty" name:"JobId"`

	// The array of task details.
	DetailList []*DeleteDomainBatchDetail `json:"DetailList,omitnil,omitempty" name:"DetailList"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DeleteDomainBatchResponse struct {
	*tchttp.BaseResponse
	Response *DeleteDomainBatchResponseParams `json:"Response"`
}

func (r *DeleteDomainBatchResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteDomainBatchResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteDomainRequestParams struct {
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

type DeleteDomainRequest struct {
	*tchttp.BaseRequest
	
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

func (r *DeleteDomainRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteDomainRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Domain")
	delete(f, "DomainId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DeleteDomainRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteDomainResponseParams struct {
	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DeleteDomainResponse struct {
	*tchttp.BaseResponse
	Response *DeleteDomainResponseParams `json:"Response"`
}

func (r *DeleteDomainResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteDomainResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeletePackageOrderRequestParams struct {
	// Specifies the domain name of the package to be terminated.
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`
}

type DeletePackageOrderRequest struct {
	*tchttp.BaseRequest
	
	// Specifies the domain name of the package to be terminated.
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`
}

func (r *DeletePackageOrderRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeletePackageOrderRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Domain")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DeletePackageOrderRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeletePackageOrderResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DeletePackageOrderResponse struct {
	*tchttp.BaseResponse
	Response *DeletePackageOrderResponseParams `json:"Response"`
}

func (r *DeletePackageOrderResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeletePackageOrderResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteRecordGroupRequestParams struct {
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// Group ID
	GroupId *uint64 `json:"GroupId,omitnil,omitempty" name:"GroupId"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

type DeleteRecordGroupRequest struct {
	*tchttp.BaseRequest
	
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// Group ID
	GroupId *uint64 `json:"GroupId,omitnil,omitempty" name:"GroupId"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

func (r *DeleteRecordGroupRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteRecordGroupRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Domain")
	delete(f, "GroupId")
	delete(f, "DomainId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DeleteRecordGroupRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteRecordGroupResponseParams struct {
	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DeleteRecordGroupResponse struct {
	*tchttp.BaseResponse
	Response *DeleteRecordGroupResponseParams `json:"Response"`
}

func (r *DeleteRecordGroupResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteRecordGroupResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteRecordRequestParams struct {
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// The record ID. You can view all DNS records and their IDs via the `DescribeRecordList` API.
	RecordId *uint64 `json:"RecordId,omitnil,omitempty" name:"RecordId"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

type DeleteRecordRequest struct {
	*tchttp.BaseRequest
	
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// The record ID. You can view all DNS records and their IDs via the `DescribeRecordList` API.
	RecordId *uint64 `json:"RecordId,omitnil,omitempty" name:"RecordId"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

func (r *DeleteRecordRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteRecordRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Domain")
	delete(f, "RecordId")
	delete(f, "DomainId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DeleteRecordRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteRecordResponseParams struct {
	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DeleteRecordResponse struct {
	*tchttp.BaseResponse
	Response *DeleteRecordResponseParams `json:"Response"`
}

func (r *DeleteRecordResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteRecordResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteShareDomainRequestParams struct {
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// Account of the domain to be shared
	Account *string `json:"Account,omitnil,omitempty" name:"Account"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

type DeleteShareDomainRequest struct {
	*tchttp.BaseRequest
	
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// Account of the domain to be shared
	Account *string `json:"Account,omitnil,omitempty" name:"Account"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

func (r *DeleteShareDomainRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteShareDomainRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Domain")
	delete(f, "Account")
	delete(f, "DomainId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DeleteShareDomainRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteShareDomainResponseParams struct {
	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DeleteShareDomainResponse struct {
	*tchttp.BaseResponse
	Response *DeleteShareDomainResponseParams `json:"Response"`
}

func (r *DeleteShareDomainResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteShareDomainResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeDomainAliasListRequestParams struct {
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *int64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

type DescribeDomainAliasListRequest struct {
	*tchttp.BaseRequest
	
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *int64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

func (r *DescribeDomainAliasListRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeDomainAliasListRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Domain")
	delete(f, "DomainId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeDomainAliasListRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeDomainAliasListResponseParams struct {
	// List of domain aliases
	DomainAliasList []*DomainAliasInfo `json:"DomainAliasList,omitnil,omitempty" name:"DomainAliasList"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeDomainAliasListResponse struct {
	*tchttp.BaseResponse
	Response *DescribeDomainAliasListResponseParams `json:"Response"`
}

func (r *DescribeDomainAliasListResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeDomainAliasListResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeDomainGroupListRequestParams struct {

}

type DescribeDomainGroupListRequest struct {
	*tchttp.BaseRequest
	
}

func (r *DescribeDomainGroupListRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeDomainGroupListRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeDomainGroupListRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeDomainGroupListResponseParams struct {
	// Group list
	GroupList []*GroupInfo `json:"GroupList,omitnil,omitempty" name:"GroupList"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeDomainGroupListResponse struct {
	*tchttp.BaseResponse
	Response *DescribeDomainGroupListResponseParams `json:"Response"`
}

func (r *DescribeDomainGroupListResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeDomainGroupListResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeDomainListRequestParams struct {
	// The domain group type. Valid values: `ALL` (default), `MINE`, `SHARE`, `ISMARK`, `PAUSE`, `VIP`, `RECENT`, `SHARE_OUT`, and `FREE`.
	Type *string `json:"Type,omitnil,omitempty" name:"Type"`

	// Record offset starting from `0`. Default value: `0`.
	Offset *int64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// Number of domains to be obtained. For example, `20` indicates to obtain 20 domains. Default value: `3000`.
	Limit *int64 `json:"Limit,omitnil,omitempty" name:"Limit"`

	// Group ID, which can be passed in to get all domains in the specified group
	GroupId *int64 `json:"GroupId,omitnil,omitempty" name:"GroupId"`

	// Keyword for searching for a domain
	Keyword *string `json:"Keyword,omitnil,omitempty" name:"Keyword"`
}

type DescribeDomainListRequest struct {
	*tchttp.BaseRequest
	
	// The domain group type. Valid values: `ALL` (default), `MINE`, `SHARE`, `ISMARK`, `PAUSE`, `VIP`, `RECENT`, `SHARE_OUT`, and `FREE`.
	Type *string `json:"Type,omitnil,omitempty" name:"Type"`

	// Record offset starting from `0`. Default value: `0`.
	Offset *int64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// Number of domains to be obtained. For example, `20` indicates to obtain 20 domains. Default value: `3000`.
	Limit *int64 `json:"Limit,omitnil,omitempty" name:"Limit"`

	// Group ID, which can be passed in to get all domains in the specified group
	GroupId *int64 `json:"GroupId,omitnil,omitempty" name:"GroupId"`

	// Keyword for searching for a domain
	Keyword *string `json:"Keyword,omitnil,omitempty" name:"Keyword"`
}

func (r *DescribeDomainListRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeDomainListRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Type")
	delete(f, "Offset")
	delete(f, "Limit")
	delete(f, "GroupId")
	delete(f, "Keyword")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeDomainListRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeDomainListResponseParams struct {
	// Statistics on the list page
	DomainCountInfo *DomainCountInfo `json:"DomainCountInfo,omitnil,omitempty" name:"DomainCountInfo"`

	// Domain list
	DomainList []*DomainListItem `json:"DomainList,omitnil,omitempty" name:"DomainList"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeDomainListResponse struct {
	*tchttp.BaseResponse
	Response *DescribeDomainListResponseParams `json:"Response"`
}

func (r *DescribeDomainListResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeDomainListResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeDomainLogListRequestParams struct {
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`

	// Record offset starting from `0`. Default value: `0`.
	Offset *uint64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// Total number of logs to be obtained. For example, `20` indicates to obtain 20 logs. Default value: `500`. Maximum value: `500`.
	Limit *uint64 `json:"Limit,omitnil,omitempty" name:"Limit"`
}

type DescribeDomainLogListRequest struct {
	*tchttp.BaseRequest
	
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`

	// Record offset starting from `0`. Default value: `0`.
	Offset *uint64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// Total number of logs to be obtained. For example, `20` indicates to obtain 20 logs. Default value: `500`. Maximum value: `500`.
	Limit *uint64 `json:"Limit,omitnil,omitempty" name:"Limit"`
}

func (r *DescribeDomainLogListRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeDomainLogListRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Domain")
	delete(f, "DomainId")
	delete(f, "Offset")
	delete(f, "Limit")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeDomainLogListRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeDomainLogListResponseParams struct {
	// Domain information
	// Note: This field may return null, indicating that no valid values can be obtained.
	LogList []*string `json:"LogList,omitnil,omitempty" name:"LogList"`

	// Number of results per page
	PageSize *uint64 `json:"PageSize,omitnil,omitempty" name:"PageSize"`

	// Total number of logs
	TotalCount *uint64 `json:"TotalCount,omitnil,omitempty" name:"TotalCount"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeDomainLogListResponse struct {
	*tchttp.BaseResponse
	Response *DescribeDomainLogListResponseParams `json:"Response"`
}

func (r *DescribeDomainLogListResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeDomainLogListResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeDomainPurviewRequestParams struct {
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

type DescribeDomainPurviewRequest struct {
	*tchttp.BaseRequest
	
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

func (r *DescribeDomainPurviewRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeDomainPurviewRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Domain")
	delete(f, "DomainId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeDomainPurviewRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeDomainPurviewResponseParams struct {
	// Permission list of the domain
	PurviewList []*PurviewInfo `json:"PurviewList,omitnil,omitempty" name:"PurviewList"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeDomainPurviewResponse struct {
	*tchttp.BaseResponse
	Response *DescribeDomainPurviewResponseParams `json:"Response"`
}

func (r *DescribeDomainPurviewResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeDomainPurviewResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeDomainRequestParams struct {
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

type DescribeDomainRequest struct {
	*tchttp.BaseRequest
	
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

func (r *DescribeDomainRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeDomainRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Domain")
	delete(f, "DomainId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeDomainRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeDomainResponseParams struct {
	// Domain information
	DomainInfo *DomainInfo `json:"DomainInfo,omitnil,omitempty" name:"DomainInfo"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeDomainResponse struct {
	*tchttp.BaseResponse
	Response *DescribeDomainResponseParams `json:"Response"`
}

func (r *DescribeDomainResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeDomainResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeDomainShareInfoRequestParams struct {
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

type DescribeDomainShareInfoRequest struct {
	*tchttp.BaseRequest
	
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

func (r *DescribeDomainShareInfoRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeDomainShareInfoRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Domain")
	delete(f, "DomainId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeDomainShareInfoRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeDomainShareInfoResponseParams struct {
	// Domain sharing information
	ShareList []*DomainShareInfo `json:"ShareList,omitnil,omitempty" name:"ShareList"`

	// Owner account of the domain
	Owner *string `json:"Owner,omitnil,omitempty" name:"Owner"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeDomainShareInfoResponse struct {
	*tchttp.BaseResponse
	Response *DescribeDomainShareInfoResponseParams `json:"Response"`
}

func (r *DescribeDomainShareInfoResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeDomainShareInfoResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeRecordGroupListRequestParams struct {
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`

	// Pagination offset
	Offset *uint64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// Number of items per page for pagination
	Limit *uint64 `json:"Limit,omitnil,omitempty" name:"Limit"`
}

type DescribeRecordGroupListRequest struct {
	*tchttp.BaseRequest
	
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`

	// Pagination offset
	Offset *uint64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// Number of items per page for pagination
	Limit *uint64 `json:"Limit,omitnil,omitempty" name:"Limit"`
}

func (r *DescribeRecordGroupListRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeRecordGroupListRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Domain")
	delete(f, "DomainId")
	delete(f, "Offset")
	delete(f, "Limit")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeRecordGroupListRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeRecordGroupListResponseParams struct {
	// Group list
	GroupList []*RecordGroupInfo `json:"GroupList,omitnil,omitempty" name:"GroupList"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeRecordGroupListResponse struct {
	*tchttp.BaseResponse
	Response *DescribeRecordGroupListResponseParams `json:"Response"`
}

func (r *DescribeRecordGroupListResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeRecordGroupListResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeRecordLineCategoryListRequestParams struct {
	// Specifies the domain name to query the line list.
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// Domain ID for querying the route list. parameter DomainId has higher priority than parameter Domain. if you pass parameter DomainId, it will ignore parameter Domain.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

type DescribeRecordLineCategoryListRequest struct {
	*tchttp.BaseRequest
	
	// Specifies the domain name to query the line list.
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// Domain ID for querying the route list. parameter DomainId has higher priority than parameter Domain. if you pass parameter DomainId, it will ignore parameter Domain.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

func (r *DescribeRecordLineCategoryListRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeRecordLineCategoryListRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Domain")
	delete(f, "DomainId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeRecordLineCategoryListRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeRecordLineCategoryListResponseParams struct {
	// Lists lines by category.
	LineList []*LineItem `json:"LineList,omitnil,omitempty" name:"LineList"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeRecordLineCategoryListResponse struct {
	*tchttp.BaseResponse
	Response *DescribeRecordLineCategoryListResponseParams `json:"Response"`
}

func (r *DescribeRecordLineCategoryListResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeRecordLineCategoryListResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeRecordLineListRequestParams struct {
	// Domain.
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// Domain level.
	// + Original plan. Valid values: `D_FREE` (Free Plan); `D_PLUS` (Individual Plus Plan); `D_EXTRA` (Enterprise 1 Plan); `D_EXPERT` (Enterprise 2 Plan); `D_ULTRA` (Enterprise 3 Plan).
	// + New plan. Valid values: `DP_FREE` (Free Version); `DP_PLUS` (Professional); `DP_EXTRA` (Enterprise Basic); `DP_EXPERT` (Enterprise); `DP_ULTRA` (Ultimate).
	DomainGrade *string `json:"DomainGrade,omitnil,omitempty" name:"DomainGrade"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

type DescribeRecordLineListRequest struct {
	*tchttp.BaseRequest
	
	// Domain.
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// Domain level.
	// + Original plan. Valid values: `D_FREE` (Free Plan); `D_PLUS` (Individual Plus Plan); `D_EXTRA` (Enterprise 1 Plan); `D_EXPERT` (Enterprise 2 Plan); `D_ULTRA` (Enterprise 3 Plan).
	// + New plan. Valid values: `DP_FREE` (Free Version); `DP_PLUS` (Professional); `DP_EXTRA` (Enterprise Basic); `DP_EXPERT` (Enterprise); `DP_ULTRA` (Ultimate).
	DomainGrade *string `json:"DomainGrade,omitnil,omitempty" name:"DomainGrade"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

func (r *DescribeRecordLineListRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeRecordLineListRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Domain")
	delete(f, "DomainGrade")
	delete(f, "DomainId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeRecordLineListRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeRecordLineListResponseParams struct {
	// List of split zones.
	LineList []*LineInfo `json:"LineList,omitnil,omitempty" name:"LineList"`

	// List of split zone groups.
	LineGroupList []*LineGroupInfo `json:"LineGroupList,omitnil,omitempty" name:"LineGroupList"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeRecordLineListResponse struct {
	*tchttp.BaseResponse
	Response *DescribeRecordLineListResponseParams `json:"Response"`
}

func (r *DescribeRecordLineListResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeRecordLineListResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeRecordListRequestParams struct {
	// The domain for which DNS records are to be obtained.
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// The ID of the domain whose DNS records are requested. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`

	// The host header of a DNS record. If this parameter is passed in, only the DNS record corresponding to this host header will be returned.
	Subdomain *string `json:"Subdomain,omitnil,omitempty" name:"Subdomain"`

	// The type of DNS record, such as A, CNAME, NS, AAAA, explicit URL, implicit URL, CAA, or SPF record.
	RecordType *string `json:"RecordType,omitnil,omitempty" name:"RecordType"`

	// The name of the split zone for which DNS records are requested. You can view split zones allowed by this domain via the `DescribeRecordLineList` API.
	RecordLine *string `json:"RecordLine,omitnil,omitempty" name:"RecordLine"`

	// The ID of the split zone for which DNS records are requested. If `RecordLineId` is passed in, `RecordLine` is ignored. You can view split zones allowed by this domain via the `DescribeRecordLineList` API.
	RecordLineId *string `json:"RecordLineId,omitnil,omitempty" name:"RecordLineId"`

	// The group ID passed in to get DNS records in the group.
	GroupId *uint64 `json:"GroupId,omitnil,omitempty" name:"GroupId"`

	// The keyword for searching for DNS records. Host headers and record values are supported.
	Keyword *string `json:"Keyword,omitnil,omitempty" name:"Keyword"`

	// The sorting field. Available values: `name`, `line`, `type`, `value`, `weight`, `mx`, and `ttl,updated_on`.
	SortField *string `json:"SortField,omitnil,omitempty" name:"SortField"`

	// The sorting type. Valid values: `ASC` (ascending, default), `DESC` (descending).
	SortType *string `json:"SortType,omitnil,omitempty" name:"SortType"`

	// The offset. Default value: `0`.
	Offset *uint64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// The limit. It defaults to 100 and can be up to 3,000.
	Limit *uint64 `json:"Limit,omitnil,omitempty" name:"Limit"`
}

type DescribeRecordListRequest struct {
	*tchttp.BaseRequest
	
	// The domain for which DNS records are to be obtained.
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// The ID of the domain whose DNS records are requested. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`

	// The host header of a DNS record. If this parameter is passed in, only the DNS record corresponding to this host header will be returned.
	Subdomain *string `json:"Subdomain,omitnil,omitempty" name:"Subdomain"`

	// The type of DNS record, such as A, CNAME, NS, AAAA, explicit URL, implicit URL, CAA, or SPF record.
	RecordType *string `json:"RecordType,omitnil,omitempty" name:"RecordType"`

	// The name of the split zone for which DNS records are requested. You can view split zones allowed by this domain via the `DescribeRecordLineList` API.
	RecordLine *string `json:"RecordLine,omitnil,omitempty" name:"RecordLine"`

	// The ID of the split zone for which DNS records are requested. If `RecordLineId` is passed in, `RecordLine` is ignored. You can view split zones allowed by this domain via the `DescribeRecordLineList` API.
	RecordLineId *string `json:"RecordLineId,omitnil,omitempty" name:"RecordLineId"`

	// The group ID passed in to get DNS records in the group.
	GroupId *uint64 `json:"GroupId,omitnil,omitempty" name:"GroupId"`

	// The keyword for searching for DNS records. Host headers and record values are supported.
	Keyword *string `json:"Keyword,omitnil,omitempty" name:"Keyword"`

	// The sorting field. Available values: `name`, `line`, `type`, `value`, `weight`, `mx`, and `ttl,updated_on`.
	SortField *string `json:"SortField,omitnil,omitempty" name:"SortField"`

	// The sorting type. Valid values: `ASC` (ascending, default), `DESC` (descending).
	SortType *string `json:"SortType,omitnil,omitempty" name:"SortType"`

	// The offset. Default value: `0`.
	Offset *uint64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// The limit. It defaults to 100 and can be up to 3,000.
	Limit *uint64 `json:"Limit,omitnil,omitempty" name:"Limit"`
}

func (r *DescribeRecordListRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeRecordListRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Domain")
	delete(f, "DomainId")
	delete(f, "Subdomain")
	delete(f, "RecordType")
	delete(f, "RecordLine")
	delete(f, "RecordLineId")
	delete(f, "GroupId")
	delete(f, "Keyword")
	delete(f, "SortField")
	delete(f, "SortType")
	delete(f, "Offset")
	delete(f, "Limit")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeRecordListRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeRecordListResponseParams struct {
	// The record count info.
	RecordCountInfo *RecordCountInfo `json:"RecordCountInfo,omitnil,omitempty" name:"RecordCountInfo"`

	// The record list result.
	RecordList []*RecordListItem `json:"RecordList,omitnil,omitempty" name:"RecordList"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeRecordListResponse struct {
	*tchttp.BaseResponse
	Response *DescribeRecordListResponseParams `json:"Response"`
}

func (r *DescribeRecordListResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeRecordListResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeRecordRequestParams struct {
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// The record ID. You can view all DNS records and their IDs via the `DescribeRecordList` API.
	RecordId *uint64 `json:"RecordId,omitnil,omitempty" name:"RecordId"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

type DescribeRecordRequest struct {
	*tchttp.BaseRequest
	
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// The record ID. You can view all DNS records and their IDs via the `DescribeRecordList` API.
	RecordId *uint64 `json:"RecordId,omitnil,omitempty" name:"RecordId"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

func (r *DescribeRecordRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeRecordRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Domain")
	delete(f, "RecordId")
	delete(f, "DomainId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeRecordRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeRecordResponseParams struct {
	// Record information
	RecordInfo *RecordInfo `json:"RecordInfo,omitnil,omitempty" name:"RecordInfo"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeRecordResponse struct {
	*tchttp.BaseResponse
	Response *DescribeRecordResponseParams `json:"Response"`
}

func (r *DescribeRecordResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeRecordResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeRecordTypeRequestParams struct {
	// Domain level.
	// + Original plan. Valid values: `D_FREE` (Free Plan); `D_PLUS` (Individual Plus Plan); `D_EXTRA` (Enterprise 1 Plan); `D_EXPERT` (Enterprise 2 Plan); `D_ULTRA` (Enterprise 3 Plan).
	// + New plan. Valid values: `DP_FREE` (Free Version); `DP_PLUS` (Professional); `DP_EXTRA` (Enterprise Basic); `DP_EXPERT` (Enterprise); `DP_ULTRA` (Ultimate).
	DomainGrade *string `json:"DomainGrade,omitnil,omitempty" name:"DomainGrade"`
}

type DescribeRecordTypeRequest struct {
	*tchttp.BaseRequest
	
	// Domain level.
	// + Original plan. Valid values: `D_FREE` (Free Plan); `D_PLUS` (Individual Plus Plan); `D_EXTRA` (Enterprise 1 Plan); `D_EXPERT` (Enterprise 2 Plan); `D_ULTRA` (Enterprise 3 Plan).
	// + New plan. Valid values: `DP_FREE` (Free Version); `DP_PLUS` (Professional); `DP_EXTRA` (Enterprise Basic); `DP_EXPERT` (Enterprise); `DP_ULTRA` (Ultimate).
	DomainGrade *string `json:"DomainGrade,omitnil,omitempty" name:"DomainGrade"`
}

func (r *DescribeRecordTypeRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeRecordTypeRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "DomainGrade")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeRecordTypeRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeRecordTypeResponseParams struct {
	// List of record types
	TypeList []*string `json:"TypeList,omitnil,omitempty" name:"TypeList"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeRecordTypeResponse struct {
	*tchttp.BaseResponse
	Response *DescribeRecordTypeResponseParams `json:"Response"`
}

func (r *DescribeRecordTypeResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeRecordTypeResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeSubdomainAnalyticsRequestParams struct {
	// The domain of the DNS query volume to be queried
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// Query start date in the format of `YYYY-MM-DD`
	StartDate *string `json:"StartDate,omitnil,omitempty" name:"StartDate"`

	// Query end date in the format of `YYYY-MM-DD`
	EndDate *string `json:"EndDate,omitnil,omitempty" name:"EndDate"`

	// The subdomain of the DNS query volume to be queried
	Subdomain *string `json:"Subdomain,omitnil,omitempty" name:"Subdomain"`

	// `DATE`: Daily. `HOUR`: Hourly
	DnsFormat *string `json:"DnsFormat,omitnil,omitempty" name:"DnsFormat"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

type DescribeSubdomainAnalyticsRequest struct {
	*tchttp.BaseRequest
	
	// The domain of the DNS query volume to be queried
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// Query start date in the format of `YYYY-MM-DD`
	StartDate *string `json:"StartDate,omitnil,omitempty" name:"StartDate"`

	// Query end date in the format of `YYYY-MM-DD`
	EndDate *string `json:"EndDate,omitnil,omitempty" name:"EndDate"`

	// The subdomain of the DNS query volume to be queried
	Subdomain *string `json:"Subdomain,omitnil,omitempty" name:"Subdomain"`

	// `DATE`: Daily. `HOUR`: Hourly
	DnsFormat *string `json:"DnsFormat,omitnil,omitempty" name:"DnsFormat"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

func (r *DescribeSubdomainAnalyticsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeSubdomainAnalyticsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Domain")
	delete(f, "StartDate")
	delete(f, "EndDate")
	delete(f, "Subdomain")
	delete(f, "DnsFormat")
	delete(f, "DomainId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeSubdomainAnalyticsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeSubdomainAnalyticsResponseParams struct {
	// DNS query volume in the current statistical dimension
	Data []*DomainAnalyticsDetail `json:"Data,omitnil,omitempty" name:"Data"`

	// Statistics on the DNS query volume of a subdomain
	Info *SubdomainAnalyticsInfo `json:"Info,omitnil,omitempty" name:"Info"`

	// DNS query volume of the subdomain alias
	AliasData []*SubdomainAliasAnalyticsItem `json:"AliasData,omitnil,omitempty" name:"AliasData"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeSubdomainAnalyticsResponse struct {
	*tchttp.BaseResponse
	Response *DescribeSubdomainAnalyticsResponseParams `json:"Response"`
}

func (r *DescribeSubdomainAnalyticsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeSubdomainAnalyticsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type DomainAliasInfo struct {
	// Domain alias ID
	Id *int64 `json:"Id,omitnil,omitempty" name:"Id"`

	// Domain alias
	DomainAlias *string `json:"DomainAlias,omitnil,omitempty" name:"DomainAlias"`
}

type DomainAnalyticsDetail struct {
	// DNS query volume in the current statistical dimension
	Num *uint64 `json:"Num,omitnil,omitempty" name:"Num"`

	// Collection date for daily collection
	DateKey *string `json:"DateKey,omitnil,omitempty" name:"DateKey"`

	// The last hour (0–23) for hourly collection. For example, if `HourKey` is `23`, the DNS query volume in the statistical period of 22:00–23:00 will be collected.
	// Note: This field may return null, indicating that no valid values can be obtained.
	HourKey *uint64 `json:"HourKey,omitnil,omitempty" name:"HourKey"`
}

type DomainCountInfo struct {
	// Number of eligible domains
	DomainTotal *uint64 `json:"DomainTotal,omitnil,omitempty" name:"DomainTotal"`

	// Number of all domains that can be viewed by the user
	AllTotal *uint64 `json:"AllTotal,omitnil,omitempty" name:"AllTotal"`

	// Number of domains added under the user account
	MineTotal *uint64 `json:"MineTotal,omitnil,omitempty" name:"MineTotal"`

	// Number of domains shared with the user
	ShareTotal *uint64 `json:"ShareTotal,omitnil,omitempty" name:"ShareTotal"`

	// Number of paid domains
	VipTotal *uint64 `json:"VipTotal,omitnil,omitempty" name:"VipTotal"`

	// Number of suspended domains
	PauseTotal *uint64 `json:"PauseTotal,omitnil,omitempty" name:"PauseTotal"`

	// Number of domains with a DNS configuration error
	ErrorTotal *uint64 `json:"ErrorTotal,omitnil,omitempty" name:"ErrorTotal"`

	// Number of locked domains
	LockTotal *uint64 `json:"LockTotal,omitnil,omitempty" name:"LockTotal"`

	// Number of blocked domains
	SpamTotal *uint64 `json:"SpamTotal,omitnil,omitempty" name:"SpamTotal"`

	// Number of domains that will expire within 30 days
	VipExpire *uint64 `json:"VipExpire,omitnil,omitempty" name:"VipExpire"`

	// Number of domains shared with others
	ShareOutTotal *uint64 `json:"ShareOutTotal,omitnil,omitempty" name:"ShareOutTotal"`

	// Number of domains in the specified group
	GroupTotal *uint64 `json:"GroupTotal,omitnil,omitempty" name:"GroupTotal"`
}

type DomainCreateInfo struct {
	// Domain ID
	Id *uint64 `json:"Id,omitnil,omitempty" name:"Id"`

	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// Domain Punycode
	Punycode *string `json:"Punycode,omitnil,omitempty" name:"Punycode"`

	// NS list of the domain
	GradeNsList []*string `json:"GradeNsList,omitnil,omitempty" name:"GradeNsList"`
}

type DomainInfo struct {
	// Domain ID
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`

	// Domain status
	Status *string `json:"Status,omitnil,omitempty" name:"Status"`

	// DNS plan level
	Grade *string `json:"Grade,omitnil,omitempty" name:"Grade"`

	// Domain group ID
	GroupId *uint64 `json:"GroupId,omitnil,omitempty" name:"GroupId"`

	// Whether the domain is starred
	IsMark *string `json:"IsMark,omitnil,omitempty" name:"IsMark"`

	// TTL (DNS record cache time)
	TTL *uint64 `json:"TTL,omitnil,omitempty" name:"TTL"`

	// Whether CNAME flattening is enabled
	CnameSpeedup *string `json:"CnameSpeedup,omitnil,omitempty" name:"CnameSpeedup"`

	// Domain remarks
	// Note: This field may return null, indicating that no valid values can be obtained.
	Remark *string `json:"Remark,omitnil,omitempty" name:"Remark"`

	// Domain Punycode
	Punycode *string `json:"Punycode,omitnil,omitempty" name:"Punycode"`

	// DNS status of the domain
	DnsStatus *string `json:"DnsStatus,omitnil,omitempty" name:"DnsStatus"`

	// NS list of the domain
	DnspodNsList []*string `json:"DnspodNsList,omitnil,omitempty" name:"DnspodNsList"`

	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// Domain level ID
	GradeLevel *uint64 `json:"GradeLevel,omitnil,omitempty" name:"GradeLevel"`

	// Domain user ID
	UserId *uint64 `json:"UserId,omitnil,omitempty" name:"UserId"`

	// Whether the domain is a VIP domain
	IsVip *string `json:"IsVip,omitnil,omitempty" name:"IsVip"`

	// Domain owner account
	Owner *string `json:"Owner,omitnil,omitempty" name:"Owner"`

	// Domain level description
	GradeTitle *string `json:"GradeTitle,omitnil,omitempty" name:"GradeTitle"`

	// Domain creation time
	CreatedOn *string `json:"CreatedOn,omitnil,omitempty" name:"CreatedOn"`

	// Last update time
	UpdatedOn *string `json:"UpdatedOn,omitnil,omitempty" name:"UpdatedOn"`

	// Tencent Cloud account `Uin`
	Uin *string `json:"Uin,omitnil,omitempty" name:"Uin"`

	// NS list actually used by the domain
	// Note: This field may return null, indicating that no valid values can be obtained.
	ActualNsList []*string `json:"ActualNsList,omitnil,omitempty" name:"ActualNsList"`

	// Number of domain records
	RecordCount *uint64 `json:"RecordCount,omitnil,omitempty" name:"RecordCount"`

	// Alias of the domain account owner
	// Note: This field may return null, indicating that no valid values can be obtained.
	OwnerNick *string `json:"OwnerNick,omitnil,omitempty" name:"OwnerNick"`
}

type DomainListItem struct {
	// Unique ID assigned to the domain by the system
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`

	// Original format of the domain
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// Domain status. Valid values: `ENABLE` (normal), `PAUSE` (suspended), `SPAM` (blocked).
	Status *string `json:"Status,omitnil,omitempty" name:"Status"`

	// Default TTL of the default DNS record of the domain
	TTL *uint64 `json:"TTL,omitnil,omitempty" name:"TTL"`

	// Whether CNAME flattening is enabled. Valid values: `ENABLE` (enabled); `DISABLE` (disabled).
	CNAMESpeedup *string `json:"CNAMESpeedup,omitnil,omitempty" name:"CNAMESpeedup"`

	// DNS configuration status. Valid values: `DNSERROR` (error), an empty string (normal).
	DNSStatus *string `json:"DNSStatus,omitnil,omitempty" name:"DNSStatus"`

	// Plan level code of the domain
	Grade *string `json:"Grade,omitnil,omitempty" name:"Grade"`

	// Group ID of the domain
	GroupId *uint64 `json:"GroupId,omitnil,omitempty" name:"GroupId"`

	// Whether search engine push optimization is enabled. Valid values: `YES` (yes), `NO` (no).
	SearchEnginePush *string `json:"SearchEnginePush,omitnil,omitempty" name:"SearchEnginePush"`

	// Domain remarks
	Remark *string `json:"Remark,omitnil,omitempty" name:"Remark"`

	// Punycode-encoded domain format
	Punycode *string `json:"Punycode,omitnil,omitempty" name:"Punycode"`

	// Effective DNS assigned to the domain by the system
	EffectiveDNS []*string `json:"EffectiveDNS,omitnil,omitempty" name:"EffectiveDNS"`

	// Number corresponding to the plan level of the domain
	GradeLevel *uint64 `json:"GradeLevel,omitnil,omitempty" name:"GradeLevel"`

	// Plan name
	GradeTitle *string `json:"GradeTitle,omitnil,omitempty" name:"GradeTitle"`

	// Whether it is a paid plan
	IsVip *string `json:"IsVip,omitnil,omitempty" name:"IsVip"`

	// Activation time of the paid plan
	VipStartAt *string `json:"VipStartAt,omitnil,omitempty" name:"VipStartAt"`

	// Expiry time of the paid plan
	VipEndAt *string `json:"VipEndAt,omitnil,omitempty" name:"VipEndAt"`

	// Whether VIP automatic renewal is enabled for the domain. Valid values: `YES` (yes); `NO` (no). Default value: `DEFAULT`.
	VipAutoRenew *string `json:"VipAutoRenew,omitnil,omitempty" name:"VipAutoRenew"`

	// Number of records under the domain
	RecordCount *uint64 `json:"RecordCount,omitnil,omitempty" name:"RecordCount"`

	// Domain adding time
	CreatedOn *string `json:"CreatedOn,omitnil,omitempty" name:"CreatedOn"`

	// Domain update time
	UpdatedOn *string `json:"UpdatedOn,omitnil,omitempty" name:"UpdatedOn"`

	// Account of the domain
	Owner *string `json:"Owner,omitnil,omitempty" name:"Owner"`
}

type DomainShareInfo struct {
	// Account with which the domain is shared
	ShareTo *string `json:"ShareTo,omitnil,omitempty" name:"ShareTo"`

	// Sharing mode. Valid values: `rw` (read/write), `r` (read-only).
	Mode *string `json:"Mode,omitnil,omitempty" name:"Mode"`

	// Sharing status. Valid values: `enabled` (shared successfully); `pending` (the account shared to does not exist and is pending registration).
	Status *string `json:"Status,omitnil,omitempty" name:"Status"`
}

type GroupInfo struct {
	// Group ID
	GroupId *int64 `json:"GroupId,omitnil,omitempty" name:"GroupId"`

	// Group name
	GroupName *string `json:"GroupName,omitnil,omitempty" name:"GroupName"`

	// Group type
	GroupType *string `json:"GroupType,omitnil,omitempty" name:"GroupType"`

	// Number of domains in the group
	Size *int64 `json:"Size,omitnil,omitempty" name:"Size"`
}

type LineGroupInfo struct {
	// Split zone group ID
	LineId *string `json:"LineId,omitnil,omitempty" name:"LineId"`

	// Split zone group name
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// Group type
	Type *string `json:"Type,omitnil,omitempty" name:"Type"`

	// List of split zones in the split zone group
	LineList []*string `json:"LineList,omitnil,omitempty" name:"LineList"`
}

type LineInfo struct {
	// Split zone name
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// Split zone ID
	LineId *string `json:"LineId,omitnil,omitempty" name:"LineId"`
}

type LineItem struct {
	// Parses the line name.
	LineName *string `json:"LineName,omitnil,omitempty" name:"LineName"`

	// Resolution route ID.
	// Note: This field may return null, indicating that no valid values can be obtained.
	LineId *string `json:"LineId,omitnil,omitempty" name:"LineId"`

	// Indicates whether the current line is available under the current domain.
	Useful *bool `json:"Useful,omitnil,omitempty" name:"Useful"`

	// Current line min plan level requirement.
	// Note: This field may return null, indicating that no valid values can be obtained.
	Grade *string `json:"Grade,omitnil,omitempty" name:"Grade"`

	// Describes the sub-line list under the current line category.
	// Note: This field may return null, indicating that no valid values can be obtained.
	SubGroup []*LineItem `json:"SubGroup,omitnil,omitempty" name:"SubGroup"`

	// Specifies the lines included in the custom line grouping.
	// Note: This field may return null, indicating that no valid values can be obtained.
	Lines []*string `json:"Lines,omitnil,omitempty" name:"Lines"`
}

type LockInfo struct {
	// Domain ID
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`

	// Domain unlock code
	LockCode *string `json:"LockCode,omitnil,omitempty" name:"LockCode"`

	// Automatic unlock date of the domain
	LockEnd *string `json:"LockEnd,omitnil,omitempty" name:"LockEnd"`
}

// Predefined struct for user
type ModifyDomainLockRequestParams struct {
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// Number of days to lock the domain. The maximum number of locked days can be obtained by calling the API for getting the permissions of a domain.
	LockDays *uint64 `json:"LockDays,omitnil,omitempty" name:"LockDays"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

type ModifyDomainLockRequest struct {
	*tchttp.BaseRequest
	
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// Number of days to lock the domain. The maximum number of locked days can be obtained by calling the API for getting the permissions of a domain.
	LockDays *uint64 `json:"LockDays,omitnil,omitempty" name:"LockDays"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

func (r *ModifyDomainLockRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyDomainLockRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Domain")
	delete(f, "LockDays")
	delete(f, "DomainId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ModifyDomainLockRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyDomainLockResponseParams struct {
	// Domain lock information
	LockInfo *LockInfo `json:"LockInfo,omitnil,omitempty" name:"LockInfo"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ModifyDomainLockResponse struct {
	*tchttp.BaseResponse
	Response *ModifyDomainLockResponseParams `json:"Response"`
}

func (r *ModifyDomainLockResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyDomainLockResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyDomainOwnerRequestParams struct {
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// The account to which to transfer the domain, which can be an UIN or email address
	Account *string `json:"Account,omitnil,omitempty" name:"Account"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

type ModifyDomainOwnerRequest struct {
	*tchttp.BaseRequest
	
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// The account to which to transfer the domain, which can be an UIN or email address
	Account *string `json:"Account,omitnil,omitempty" name:"Account"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

func (r *ModifyDomainOwnerRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyDomainOwnerRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Domain")
	delete(f, "Account")
	delete(f, "DomainId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ModifyDomainOwnerRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyDomainOwnerResponseParams struct {
	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ModifyDomainOwnerResponse struct {
	*tchttp.BaseResponse
	Response *ModifyDomainOwnerResponseParams `json:"Response"`
}

func (r *ModifyDomainOwnerResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyDomainOwnerResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyDomainRemarkRequestParams struct {
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`

	// Domain remarks. To delete the remarks, submit empty content.
	Remark *string `json:"Remark,omitnil,omitempty" name:"Remark"`
}

type ModifyDomainRemarkRequest struct {
	*tchttp.BaseRequest
	
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`

	// Domain remarks. To delete the remarks, submit empty content.
	Remark *string `json:"Remark,omitnil,omitempty" name:"Remark"`
}

func (r *ModifyDomainRemarkRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyDomainRemarkRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Domain")
	delete(f, "DomainId")
	delete(f, "Remark")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ModifyDomainRemarkRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyDomainRemarkResponseParams struct {
	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ModifyDomainRemarkResponse struct {
	*tchttp.BaseResponse
	Response *ModifyDomainRemarkResponseParams `json:"Response"`
}

func (r *ModifyDomainRemarkResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyDomainRemarkResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyDomainStatusRequestParams struct {
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// Domain status. Valid values: enable; disable.
	Status *string `json:"Status,omitnil,omitempty" name:"Status"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

type ModifyDomainStatusRequest struct {
	*tchttp.BaseRequest
	
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// Domain status. Valid values: enable; disable.
	Status *string `json:"Status,omitnil,omitempty" name:"Status"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

func (r *ModifyDomainStatusRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyDomainStatusRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Domain")
	delete(f, "Status")
	delete(f, "DomainId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ModifyDomainStatusRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyDomainStatusResponseParams struct {
	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ModifyDomainStatusResponse struct {
	*tchttp.BaseResponse
	Response *ModifyDomainStatusResponseParams `json:"Response"`
}

func (r *ModifyDomainStatusResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyDomainStatusResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyDomainUnlockRequestParams struct {
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// Domain unlock code, which will be returned when the domain is locked.
	LockCode *string `json:"LockCode,omitnil,omitempty" name:"LockCode"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

type ModifyDomainUnlockRequest struct {
	*tchttp.BaseRequest
	
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// Domain unlock code, which will be returned when the domain is locked.
	LockCode *string `json:"LockCode,omitnil,omitempty" name:"LockCode"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

func (r *ModifyDomainUnlockRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyDomainUnlockRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Domain")
	delete(f, "LockCode")
	delete(f, "DomainId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ModifyDomainUnlockRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyDomainUnlockResponseParams struct {
	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ModifyDomainUnlockResponse struct {
	*tchttp.BaseResponse
	Response *ModifyDomainUnlockResponseParams `json:"Response"`
}

func (r *ModifyDomainUnlockResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyDomainUnlockResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type ModifyRecordBatchDetail struct {
	// See `RecordInfoBatchModify`.
	// Note: This field may return null, indicating that no valid values can be obtained.
	RecordList []*BatchRecordInfo `json:"RecordList,omitnil,omitempty" name:"RecordList"`

	// Task ID
	Id *uint64 `json:"Id,omitnil,omitempty" name:"Id"`

	// Domain
	// Note: This field may return null, indicating that no valid values can be obtained.
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// Domain level
	// Note: This field may return null, indicating that no valid values can be obtained.
	DomainGrade *string `json:"DomainGrade,omitnil,omitempty" name:"DomainGrade"`

	// Error message
	// Note: This field may return null, indicating that no valid values can be obtained.
	ErrMsg *string `json:"ErrMsg,omitnil,omitempty" name:"ErrMsg"`

	// Task running status
	// Note: This field may return null, indicating that no valid values can be obtained.
	Status *string `json:"Status,omitnil,omitempty" name:"Status"`

	// Operation type
	// Note: This field may return null, indicating that no valid values can be obtained.
	Operation *string `json:"Operation,omitnil,omitempty" name:"Operation"`

	// Domain ID
	// Note: This field may return null, indicating that no valid values can be obtained.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

// Predefined struct for user
type ModifyRecordBatchRequestParams struct {
	// Array of record IDs. You can view all DNS records and their IDs via the `DescribeRecordList` API.
	RecordIdList []*uint64 `json:"RecordIdList,omitnil,omitempty" name:"RecordIdList"`

	// The field to be modified. Valid values: `sub_domain`, `record_type`, `area`, `value`, `mx`, `ttl`, `status`.
	Change *string `json:"Change,omitnil,omitempty" name:"Change"`

	// The value to be changed to, which is required and subject to the `change` field.
	ChangeTo *string `json:"ChangeTo,omitnil,omitempty" name:"ChangeTo"`

	// The record value to be changed to, which is required only if the `change` field is `record_type`.
	Value *string `json:"Value,omitnil,omitempty" name:"Value"`

	// MX record priority, which is required only if the `Change` field is `mx`.
	MX *string `json:"MX,omitnil,omitempty" name:"MX"`
}

type ModifyRecordBatchRequest struct {
	*tchttp.BaseRequest
	
	// Array of record IDs. You can view all DNS records and their IDs via the `DescribeRecordList` API.
	RecordIdList []*uint64 `json:"RecordIdList,omitnil,omitempty" name:"RecordIdList"`

	// The field to be modified. Valid values: `sub_domain`, `record_type`, `area`, `value`, `mx`, `ttl`, `status`.
	Change *string `json:"Change,omitnil,omitempty" name:"Change"`

	// The value to be changed to, which is required and subject to the `change` field.
	ChangeTo *string `json:"ChangeTo,omitnil,omitempty" name:"ChangeTo"`

	// The record value to be changed to, which is required only if the `change` field is `record_type`.
	Value *string `json:"Value,omitnil,omitempty" name:"Value"`

	// MX record priority, which is required only if the `Change` field is `mx`.
	MX *string `json:"MX,omitnil,omitempty" name:"MX"`
}

func (r *ModifyRecordBatchRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyRecordBatchRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "RecordIdList")
	delete(f, "Change")
	delete(f, "ChangeTo")
	delete(f, "Value")
	delete(f, "MX")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ModifyRecordBatchRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyRecordBatchResponseParams struct {
	// Bulk task ID
	JobId *uint64 `json:"JobId,omitnil,omitempty" name:"JobId"`

	// See `modifyRecordBatchDetail`.
	DetailList []*ModifyRecordBatchDetail `json:"DetailList,omitnil,omitempty" name:"DetailList"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ModifyRecordBatchResponse struct {
	*tchttp.BaseResponse
	Response *ModifyRecordBatchResponseParams `json:"Response"`
}

func (r *ModifyRecordBatchResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyRecordBatchResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyRecordGroupRequestParams struct {
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// Group name
	GroupName *string `json:"GroupName,omitnil,omitempty" name:"GroupName"`

	// ID of the group to be modified
	GroupId *uint64 `json:"GroupId,omitnil,omitempty" name:"GroupId"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

type ModifyRecordGroupRequest struct {
	*tchttp.BaseRequest
	
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// Group name
	GroupName *string `json:"GroupName,omitnil,omitempty" name:"GroupName"`

	// ID of the group to be modified
	GroupId *uint64 `json:"GroupId,omitnil,omitempty" name:"GroupId"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

func (r *ModifyRecordGroupRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyRecordGroupRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Domain")
	delete(f, "GroupName")
	delete(f, "GroupId")
	delete(f, "DomainId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ModifyRecordGroupRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyRecordGroupResponseParams struct {
	// ID of the modified group
	GroupId *uint64 `json:"GroupId,omitnil,omitempty" name:"GroupId"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ModifyRecordGroupResponse struct {
	*tchttp.BaseResponse
	Response *ModifyRecordGroupResponseParams `json:"Response"`
}

func (r *ModifyRecordGroupResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyRecordGroupResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyRecordRemarkRequestParams struct {
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// The record ID. You can view all DNS records and their IDs via the `DescribeRecordList` API.
	RecordId *uint64 `json:"RecordId,omitnil,omitempty" name:"RecordId"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`

	// DNS record remarks. To delete the remarks, submit empty content.
	Remark *string `json:"Remark,omitnil,omitempty" name:"Remark"`
}

type ModifyRecordRemarkRequest struct {
	*tchttp.BaseRequest
	
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// The record ID. You can view all DNS records and their IDs via the `DescribeRecordList` API.
	RecordId *uint64 `json:"RecordId,omitnil,omitempty" name:"RecordId"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`

	// DNS record remarks. To delete the remarks, submit empty content.
	Remark *string `json:"Remark,omitnil,omitempty" name:"Remark"`
}

func (r *ModifyRecordRemarkRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyRecordRemarkRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Domain")
	delete(f, "RecordId")
	delete(f, "DomainId")
	delete(f, "Remark")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ModifyRecordRemarkRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyRecordRemarkResponseParams struct {
	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ModifyRecordRemarkResponse struct {
	*tchttp.BaseResponse
	Response *ModifyRecordRemarkResponseParams `json:"Response"`
}

func (r *ModifyRecordRemarkResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyRecordRemarkResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyRecordRequestParams struct {
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// Record type, which is obtained through the record type API. The value contains uppercase letters, such as `A`.
	RecordType *string `json:"RecordType,omitnil,omitempty" name:"RecordType"`

	// Record split zone, which is obtained through the record split zone API.
	RecordLine *string `json:"RecordLine,omitnil,omitempty" name:"RecordLine"`

	// Record value, such as `IP : 200.200.200.200`, `CNAME : cname.dnspod.com`, and `MX : mail.dnspod.com`.
	Value *string `json:"Value,omitnil,omitempty" name:"Value"`

	// The record ID. You can view all DNS records and their IDs via the `DescribeRecordList` API.
	RecordId *uint64 `json:"RecordId,omitnil,omitempty" name:"RecordId"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`

	// Host record such as `www`. If it is not passed in, it will be `@` by default.
	SubDomain *string `json:"SubDomain,omitnil,omitempty" name:"SubDomain"`

	// Split zone ID, which is obtained through the record split zone API. The value is a string such as `10=1`. The `RecordLineId` parameter has a higher priority than `RecordLine`. If both of them are passed in, `RecordLineId` will be used first.
	RecordLineId *string `json:"RecordLineId,omitnil,omitempty" name:"RecordLineId"`

	// MX priority, which is required for an MX record and will take effect if the record type is MX. Value range: 1–20.
	MX *uint64 `json:"MX,omitnil,omitempty" name:"MX"`

	// TTL. Value range: 1–604800. The minimum value varies by domain level.
	TTL *uint64 `json:"TTL,omitnil,omitempty" name:"TTL"`

	// Weight information, which is an integer between 0 and 100. It is supported only for enterprise VIP domains. `0` indicates not to pass in this parameter, i.e., not to set the weight.
	Weight *uint64 `json:"Weight,omitnil,omitempty" name:"Weight"`

	// Initial status of the record. Valid values: ENABLE, DISABLE. Default value: ENABLE. If `DISABLE` is passed in, the DNS record won't take effect, and the limit on round-robin DNS won't be verified.
	Status *string `json:"Status,omitnil,omitempty" name:"Status"`
}

type ModifyRecordRequest struct {
	*tchttp.BaseRequest
	
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// Record type, which is obtained through the record type API. The value contains uppercase letters, such as `A`.
	RecordType *string `json:"RecordType,omitnil,omitempty" name:"RecordType"`

	// Record split zone, which is obtained through the record split zone API.
	RecordLine *string `json:"RecordLine,omitnil,omitempty" name:"RecordLine"`

	// Record value, such as `IP : 200.200.200.200`, `CNAME : cname.dnspod.com`, and `MX : mail.dnspod.com`.
	Value *string `json:"Value,omitnil,omitempty" name:"Value"`

	// The record ID. You can view all DNS records and their IDs via the `DescribeRecordList` API.
	RecordId *uint64 `json:"RecordId,omitnil,omitempty" name:"RecordId"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`

	// Host record such as `www`. If it is not passed in, it will be `@` by default.
	SubDomain *string `json:"SubDomain,omitnil,omitempty" name:"SubDomain"`

	// Split zone ID, which is obtained through the record split zone API. The value is a string such as `10=1`. The `RecordLineId` parameter has a higher priority than `RecordLine`. If both of them are passed in, `RecordLineId` will be used first.
	RecordLineId *string `json:"RecordLineId,omitnil,omitempty" name:"RecordLineId"`

	// MX priority, which is required for an MX record and will take effect if the record type is MX. Value range: 1–20.
	MX *uint64 `json:"MX,omitnil,omitempty" name:"MX"`

	// TTL. Value range: 1–604800. The minimum value varies by domain level.
	TTL *uint64 `json:"TTL,omitnil,omitempty" name:"TTL"`

	// Weight information, which is an integer between 0 and 100. It is supported only for enterprise VIP domains. `0` indicates not to pass in this parameter, i.e., not to set the weight.
	Weight *uint64 `json:"Weight,omitnil,omitempty" name:"Weight"`

	// Initial status of the record. Valid values: ENABLE, DISABLE. Default value: ENABLE. If `DISABLE` is passed in, the DNS record won't take effect, and the limit on round-robin DNS won't be verified.
	Status *string `json:"Status,omitnil,omitempty" name:"Status"`
}

func (r *ModifyRecordRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyRecordRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Domain")
	delete(f, "RecordType")
	delete(f, "RecordLine")
	delete(f, "Value")
	delete(f, "RecordId")
	delete(f, "DomainId")
	delete(f, "SubDomain")
	delete(f, "RecordLineId")
	delete(f, "MX")
	delete(f, "TTL")
	delete(f, "Weight")
	delete(f, "Status")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ModifyRecordRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyRecordResponseParams struct {
	// Record ID
	RecordId *uint64 `json:"RecordId,omitnil,omitempty" name:"RecordId"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ModifyRecordResponse struct {
	*tchttp.BaseResponse
	Response *ModifyRecordResponseParams `json:"Response"`
}

func (r *ModifyRecordResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyRecordResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyRecordStatusRequestParams struct {
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// The record ID. You can view all DNS records and their IDs via the `DescribeRecordList` API.
	RecordId *uint64 `json:"RecordId,omitnil,omitempty" name:"RecordId"`

	// Record status. Valid values: `ENABLE`, `DISABLE`. If `DISABLE` is passed in, the DNS record won't take effect, and the limit on round-robin DNS won't be verified.
	Status *string `json:"Status,omitnil,omitempty" name:"Status"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

type ModifyRecordStatusRequest struct {
	*tchttp.BaseRequest
	
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// The record ID. You can view all DNS records and their IDs via the `DescribeRecordList` API.
	RecordId *uint64 `json:"RecordId,omitnil,omitempty" name:"RecordId"`

	// Record status. Valid values: `ENABLE`, `DISABLE`. If `DISABLE` is passed in, the DNS record won't take effect, and the limit on round-robin DNS won't be verified.
	Status *string `json:"Status,omitnil,omitempty" name:"Status"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored. You can view all `Domain` and `DomainId` values via the `DescribeDomainList` API.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

func (r *ModifyRecordStatusRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyRecordStatusRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Domain")
	delete(f, "RecordId")
	delete(f, "Status")
	delete(f, "DomainId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ModifyRecordStatusRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyRecordStatusResponseParams struct {
	// Record ID.
	RecordId *uint64 `json:"RecordId,omitnil,omitempty" name:"RecordId"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ModifyRecordStatusResponse struct {
	*tchttp.BaseResponse
	Response *ModifyRecordStatusResponseParams `json:"Response"`
}

func (r *ModifyRecordStatusResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyRecordStatusResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyRecordToGroupRequestParams struct {
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// Group ID
	GroupId *uint64 `json:"GroupId,omitnil,omitempty" name:"GroupId"`

	// Record ID. Separate multiple IDs by vertical bar ("|").
	RecordId *string `json:"RecordId,omitnil,omitempty" name:"RecordId"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

type ModifyRecordToGroupRequest struct {
	*tchttp.BaseRequest
	
	// Domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// Group ID
	GroupId *uint64 `json:"GroupId,omitnil,omitempty" name:"GroupId"`

	// Record ID. Separate multiple IDs by vertical bar ("|").
	RecordId *string `json:"RecordId,omitnil,omitempty" name:"RecordId"`

	// The domain ID. `DomainId` takes priority over `Domain`. If `DomainId` is passed in, `Domain` is ignored.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

func (r *ModifyRecordToGroupRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyRecordToGroupRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Domain")
	delete(f, "GroupId")
	delete(f, "RecordId")
	delete(f, "DomainId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ModifyRecordToGroupRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyRecordToGroupResponseParams struct {
	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ModifyRecordToGroupResponse struct {
	*tchttp.BaseResponse
	Response *ModifyRecordToGroupResponseParams `json:"Response"`
}

func (r *ModifyRecordToGroupResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyRecordToGroupResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type PurviewInfo struct {
	// Permission name
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// Permission value
	Value *string `json:"Value,omitnil,omitempty" name:"Value"`
}

type RecordCountInfo struct {
	// The subdomain count.
	SubdomainCount *uint64 `json:"SubdomainCount,omitnil,omitempty" name:"SubdomainCount"`

	// The count of records returned in the list.
	ListCount *uint64 `json:"ListCount,omitnil,omitempty" name:"ListCount"`

	// The total record count.
	TotalCount *uint64 `json:"TotalCount,omitnil,omitempty" name:"TotalCount"`
}

type RecordGroupInfo struct {
	// Group ID
	GroupId *uint64 `json:"GroupId,omitnil,omitempty" name:"GroupId"`

	// Group name
	GroupName *string `json:"GroupName,omitnil,omitempty" name:"GroupName"`

	// Group type. Valid values: `system`, `user`.
	GroupType *string `json:"GroupType,omitnil,omitempty" name:"GroupType"`
}

type RecordInfo struct {
	// Record ID.
	Id *uint64 `json:"Id,omitnil,omitempty" name:"Id"`

	// Subdomain (host record).
	SubDomain *string `json:"SubDomain,omitnil,omitempty" name:"SubDomain"`

	// Record type. For more information, see the `DescribeRecordType` API.
	RecordType *string `json:"RecordType,omitnil,omitempty" name:"RecordType"`

	// Split zone of the DNS record. For more information, see the `DescribeRecordLineList` API.
	RecordLine *string `json:"RecordLine,omitnil,omitempty" name:"RecordLine"`

	// Split zone ID of the DNS record. For more information, see the `DescribeRecordLineList` API.
	RecordLineId *string `json:"RecordLineId,omitnil,omitempty" name:"RecordLineId"`

	// Record value.
	Value *string `json:"Value,omitnil,omitempty" name:"Value"`

	// Record weight.
	// Note: This field may return null, indicating that no valid values can be obtained.
	Weight *uint64 `json:"Weight,omitnil,omitempty" name:"Weight"`

	// MX record value. It is 0 by default for non-MX records.
	MX *uint64 `json:"MX,omitnil,omitempty" name:"MX"`

	// TTL value of the record.
	TTL *uint64 `json:"TTL,omitnil,omitempty" name:"TTL"`

	// Record status. Valid values: 0 (disabled); 1 (enabled).
	Enabled *uint64 `json:"Enabled,omitnil,omitempty" name:"Enabled"`

	// D-Monitor status of the record.
	// "Ok" : The server is normal.
	// "Warn" : There is an alarm on this record, and the server returns 4XX.
	// "Down" : The server is down.
	// "" : D-Monitor is disabled for this record.
	MonitorStatus *string `json:"MonitorStatus,omitnil,omitempty" name:"MonitorStatus"`

	// Record remarks.
	// Note: This field may return null, indicating that no valid values can be obtained.
	Remark *string `json:"Remark,omitnil,omitempty" name:"Remark"`

	// Last update time of the record.
	UpdatedOn *string `json:"UpdatedOn,omitnil,omitempty" name:"UpdatedOn"`

	// Domain ID.
	DomainId *uint64 `json:"DomainId,omitnil,omitempty" name:"DomainId"`
}

type RecordListItem struct {
	// The record ID.
	RecordId *uint64 `json:"RecordId,omitnil,omitempty" name:"RecordId"`

	// The record value.
	Value *string `json:"Value,omitnil,omitempty" name:"Value"`

	// The record status. Valid values: `ENABLE` (enabled), `DISABLE` (disabled).
	Status *string `json:"Status,omitnil,omitempty" name:"Status"`

	// The update time.
	UpdatedOn *string `json:"UpdatedOn,omitnil,omitempty" name:"UpdatedOn"`

	// The host name.
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// The record split zone.
	Line *string `json:"Line,omitnil,omitempty" name:"Line"`

	// The split zone ID.
	LineId *string `json:"LineId,omitnil,omitempty" name:"LineId"`

	// The record type.
	Type *string `json:"Type,omitnil,omitempty" name:"Type"`

	// The record weight, which is required for round-robin DNS records.
	// Note: This field may return null, indicating that no valid values can be obtained.
	Weight *uint64 `json:"Weight,omitnil,omitempty" name:"Weight"`

	// The monitoring status of the record. Valid values: `OK` (normal), `WARN` (warning), and `DOWN` (downtime). It is empty if no monitoring is set or the monitoring is suspended.
	MonitorStatus *string `json:"MonitorStatus,omitnil,omitempty" name:"MonitorStatus"`

	// The record remarks.
	Remark *string `json:"Remark,omitnil,omitempty" name:"Remark"`

	// The record cache time.
	TTL *uint64 `json:"TTL,omitnil,omitempty" name:"TTL"`

	// The MX value, applicable to the MX record only.
	// Note: This field may return null, indicating that no valid values can be obtained.
	MX *uint64 `json:"MX,omitnil,omitempty" name:"MX"`

	// Whether it is a default NS record.
	DefaultNS *bool `json:"DefaultNS,omitnil,omitempty" name:"DefaultNS"`
}

type SubdomainAliasAnalyticsItem struct {
	// Statistics on the DNS query volume of a subdomain
	Info *SubdomainAnalyticsInfo `json:"Info,omitnil,omitempty" name:"Info"`

	// DNS query volume in the current statistical dimension
	Data []*DomainAnalyticsDetail `json:"Data,omitnil,omitempty" name:"Data"`
}

type SubdomainAnalyticsInfo struct {
	// `DATE`: Daily; `HOUR`: Hourly
	DnsFormat *string `json:"DnsFormat,omitnil,omitempty" name:"DnsFormat"`

	// Total DNS query volume for the current statistical period
	DnsTotal *uint64 `json:"DnsTotal,omitnil,omitempty" name:"DnsTotal"`

	// The queried domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// Start date of the current statistical period
	StartDate *string `json:"StartDate,omitnil,omitempty" name:"StartDate"`

	// End date of the current statistical period
	EndDate *string `json:"EndDate,omitnil,omitempty" name:"EndDate"`

	// Subdomain
	Subdomain *string `json:"Subdomain,omitnil,omitempty" name:"Subdomain"`
}