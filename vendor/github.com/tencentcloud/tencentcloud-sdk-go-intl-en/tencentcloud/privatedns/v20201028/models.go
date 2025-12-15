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

package v20201028

import (
    tcerr "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common/errors"
    tchttp "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common/http"
    "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common/json"
)

type AccountVpcInfo struct {
	// VpcId: vpc-xadsafsdasd
	UniqVpcId *string `json:"UniqVpcId,omitnil,omitempty" name:"UniqVpcId"`

	// VPC region: ap-guangzhou, ap-shanghai
	// Note: this field may return `null`, indicating that no valid values can be obtained.
	Region *string `json:"Region,omitnil,omitempty" name:"Region"`

	// VPC account: 123456789
	// Note: this field may return `null`, indicating that no valid values can be obtained.
	Uin *string `json:"Uin,omitnil,omitempty" name:"Uin"`

	// VPC name: testname
	// Note: this field may return `null`, indicating that no valid values can be obtained.
	VpcName *string `json:"VpcName,omitnil,omitempty" name:"VpcName"`
}

type AccountVpcInfoOut struct {
	// VpcId: vpc-xadsafsdasd
	VpcId *string `json:"VpcId,omitnil,omitempty" name:"VpcId"`

	// Region: ap-guangzhou, ap-shanghai
	Region *string `json:"Region,omitnil,omitempty" name:"Region"`

	// VPC ID: 123456789
	Uin *string `json:"Uin,omitnil,omitempty" name:"Uin"`

	// VPC name: testname
	VpcName *string `json:"VpcName,omitnil,omitempty" name:"VpcName"`
}

type AccountVpcInfoOutput struct {
	// UIN of the VPC account
	Uin *string `json:"Uin,omitnil,omitempty" name:"Uin"`

	// VPC ID
	UniqVpcId *string `json:"UniqVpcId,omitnil,omitempty" name:"UniqVpcId"`

	// Region
	Region *string `json:"Region,omitnil,omitempty" name:"Region"`
}

type AuditLog struct {
	// Log type
	Resource *string `json:"Resource,omitnil,omitempty" name:"Resource"`

	// Log table name
	Metric *string `json:"Metric,omitnil,omitempty" name:"Metric"`

	// Total number of logs
	TotalCount *int64 `json:"TotalCount,omitnil,omitempty" name:"TotalCount"`

	// List of logs
	DataSet []*AuditLogInfo `json:"DataSet,omitnil,omitempty" name:"DataSet"`
}

type AuditLogInfo struct {
	// Time
	Date *string `json:"Date,omitnil,omitempty" name:"Date"`

	// Operator UIN
	OperatorUin *string `json:"OperatorUin,omitnil,omitempty" name:"OperatorUin"`

	// Log content
	Content *string `json:"Content,omitnil,omitempty" name:"Content"`
}

// Predefined struct for user
type CreateEndPointAndEndPointServiceRequestParams struct {
	// VPC instance ID.
	VpcId *string `json:"VpcId,omitnil,omitempty" name:"VpcId"`

	// Whether automatic forwarding is supported.
	AutoAcceptFlag *bool `json:"AutoAcceptFlag,omitnil,omitempty" name:"AutoAcceptFlag"`

	// Backend service ID.
	ServiceInstanceId *string `json:"ServiceInstanceId,omitnil,omitempty" name:"ServiceInstanceId"`

	// Endpoint name.
	EndPointName *string `json:"EndPointName,omitnil,omitempty" name:"EndPointName"`

	// Endpoint region, which should be consistent with the region of the endpoint service.
	EndPointRegion *string `json:"EndPointRegion,omitnil,omitempty" name:"EndPointRegion"`

	// Endpoint service name.
	EndPointServiceName *string `json:"EndPointServiceName,omitnil,omitempty" name:"EndPointServiceName"`

	// Mounted PaaS service type. Valid values: CLB, CDB, and CRS.
	ServiceType *string `json:"ServiceType,omitnil,omitempty" name:"ServiceType"`

	// Number of endpoint IP addresses.
	IpNum *int64 `json:"IpNum,omitnil,omitempty" name:"IpNum"`
}

type CreateEndPointAndEndPointServiceRequest struct {
	*tchttp.BaseRequest
	
	// VPC instance ID.
	VpcId *string `json:"VpcId,omitnil,omitempty" name:"VpcId"`

	// Whether automatic forwarding is supported.
	AutoAcceptFlag *bool `json:"AutoAcceptFlag,omitnil,omitempty" name:"AutoAcceptFlag"`

	// Backend service ID.
	ServiceInstanceId *string `json:"ServiceInstanceId,omitnil,omitempty" name:"ServiceInstanceId"`

	// Endpoint name.
	EndPointName *string `json:"EndPointName,omitnil,omitempty" name:"EndPointName"`

	// Endpoint region, which should be consistent with the region of the endpoint service.
	EndPointRegion *string `json:"EndPointRegion,omitnil,omitempty" name:"EndPointRegion"`

	// Endpoint service name.
	EndPointServiceName *string `json:"EndPointServiceName,omitnil,omitempty" name:"EndPointServiceName"`

	// Mounted PaaS service type. Valid values: CLB, CDB, and CRS.
	ServiceType *string `json:"ServiceType,omitnil,omitempty" name:"ServiceType"`

	// Number of endpoint IP addresses.
	IpNum *int64 `json:"IpNum,omitnil,omitempty" name:"IpNum"`
}

func (r *CreateEndPointAndEndPointServiceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateEndPointAndEndPointServiceRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "VpcId")
	delete(f, "AutoAcceptFlag")
	delete(f, "ServiceInstanceId")
	delete(f, "EndPointName")
	delete(f, "EndPointRegion")
	delete(f, "EndPointServiceName")
	delete(f, "ServiceType")
	delete(f, "IpNum")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CreateEndPointAndEndPointServiceRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateEndPointAndEndPointServiceResponseParams struct {
	// Endpoint ID.
	EndPointId *string `json:"EndPointId,omitnil,omitempty" name:"EndPointId"`

	// Endpoint name.
	EndPointName *string `json:"EndPointName,omitnil,omitempty" name:"EndPointName"`

	// Endpoint service ID.
	EndPointServiceId *string `json:"EndPointServiceId,omitnil,omitempty" name:"EndPointServiceId"`

	// IP address list of the endpoint.
	EndPointVipSet []*string `json:"EndPointVipSet,omitnil,omitempty" name:"EndPointVipSet"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type CreateEndPointAndEndPointServiceResponse struct {
	*tchttp.BaseResponse
	Response *CreateEndPointAndEndPointServiceResponseParams `json:"Response"`
}

func (r *CreateEndPointAndEndPointServiceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateEndPointAndEndPointServiceResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateEndPointRequestParams struct {
	// Endpoint name.
	EndPointName *string `json:"EndPointName,omitnil,omitempty" name:"EndPointName"`

	// Endpoint service ID (namely, VPC endpoint service ID).
	EndPointServiceId *string `json:"EndPointServiceId,omitnil,omitempty" name:"EndPointServiceId"`

	// Endpoint region, which should be consistent with the region of the endpoint service.
	EndPointRegion *string `json:"EndPointRegion,omitnil,omitempty" name:"EndPointRegion"`

	// Number of endpoint IP addresses.
	IpNum *int64 `json:"IpNum,omitnil,omitempty" name:"IpNum"`
}

type CreateEndPointRequest struct {
	*tchttp.BaseRequest
	
	// Endpoint name.
	EndPointName *string `json:"EndPointName,omitnil,omitempty" name:"EndPointName"`

	// Endpoint service ID (namely, VPC endpoint service ID).
	EndPointServiceId *string `json:"EndPointServiceId,omitnil,omitempty" name:"EndPointServiceId"`

	// Endpoint region, which should be consistent with the region of the endpoint service.
	EndPointRegion *string `json:"EndPointRegion,omitnil,omitempty" name:"EndPointRegion"`

	// Number of endpoint IP addresses.
	IpNum *int64 `json:"IpNum,omitnil,omitempty" name:"IpNum"`
}

func (r *CreateEndPointRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateEndPointRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "EndPointName")
	delete(f, "EndPointServiceId")
	delete(f, "EndPointRegion")
	delete(f, "IpNum")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CreateEndPointRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateEndPointResponseParams struct {
	// Endpoint ID.
	EndPointId *string `json:"EndPointId,omitnil,omitempty" name:"EndPointId"`

	// Endpoint name.
	EndPointName *string `json:"EndPointName,omitnil,omitempty" name:"EndPointName"`

	// Endpoint service ID.
	EndPointServiceId *string `json:"EndPointServiceId,omitnil,omitempty" name:"EndPointServiceId"`

	// IP address list of the endpoint.
	EndPointVipSet []*string `json:"EndPointVipSet,omitnil,omitempty" name:"EndPointVipSet"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type CreateEndPointResponse struct {
	*tchttp.BaseResponse
	Response *CreateEndPointResponseParams `json:"Response"`
}

func (r *CreateEndPointResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateEndPointResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateExtendEndpointRequestParams struct {
	// Outbound endpoint name.
	EndpointName *string `json:"EndpointName,omitnil,omitempty" name:"EndpointName"`

	// The region of the outbound endpoint must be consistent with the region of the forwarding target VIP.
	EndpointRegion *string `json:"EndpointRegion,omitnil,omitempty" name:"EndpointRegion"`

	// Forwarding target.
	ForwardIp *ForwardIp `json:"ForwardIp,omitnil,omitempty" name:"ForwardIp"`
}

type CreateExtendEndpointRequest struct {
	*tchttp.BaseRequest
	
	// Outbound endpoint name.
	EndpointName *string `json:"EndpointName,omitnil,omitempty" name:"EndpointName"`

	// The region of the outbound endpoint must be consistent with the region of the forwarding target VIP.
	EndpointRegion *string `json:"EndpointRegion,omitnil,omitempty" name:"EndpointRegion"`

	// Forwarding target.
	ForwardIp *ForwardIp `json:"ForwardIp,omitnil,omitempty" name:"ForwardIp"`
}

func (r *CreateExtendEndpointRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateExtendEndpointRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "EndpointName")
	delete(f, "EndpointRegion")
	delete(f, "ForwardIp")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CreateExtendEndpointRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateExtendEndpointResponseParams struct {
	// Endpoint ID.
	EndpointId *string `json:"EndpointId,omitnil,omitempty" name:"EndpointId"`

	// Endpoint name.
	EndpointName *string `json:"EndpointName,omitnil,omitempty" name:"EndpointName"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type CreateExtendEndpointResponse struct {
	*tchttp.BaseResponse
	Response *CreateExtendEndpointResponseParams `json:"Response"`
}

func (r *CreateExtendEndpointResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateExtendEndpointResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateForwardRuleRequestParams struct {
	// Forwarding rule name.
	RuleName *string `json:"RuleName,omitnil,omitempty" name:"RuleName"`

	// Forwarding rule type. DOWN: From cloud to off-cloud; UP: From off-cloud to cloud.
	RuleType *string `json:"RuleType,omitnil,omitempty" name:"RuleType"`

	// Private domain ID, which can be viewed on the private domain list page.
	ZoneId *string `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`

	// Endpoint ID.
	EndPointId *string `json:"EndPointId,omitnil,omitempty" name:"EndPointId"`
}

type CreateForwardRuleRequest struct {
	*tchttp.BaseRequest
	
	// Forwarding rule name.
	RuleName *string `json:"RuleName,omitnil,omitempty" name:"RuleName"`

	// Forwarding rule type. DOWN: From cloud to off-cloud; UP: From off-cloud to cloud.
	RuleType *string `json:"RuleType,omitnil,omitempty" name:"RuleType"`

	// Private domain ID, which can be viewed on the private domain list page.
	ZoneId *string `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`

	// Endpoint ID.
	EndPointId *string `json:"EndPointId,omitnil,omitempty" name:"EndPointId"`
}

func (r *CreateForwardRuleRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateForwardRuleRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "RuleName")
	delete(f, "RuleType")
	delete(f, "ZoneId")
	delete(f, "EndPointId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CreateForwardRuleRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateForwardRuleResponseParams struct {
	// Forwarding rule ID.
	RuleId *string `json:"RuleId,omitnil,omitempty" name:"RuleId"`

	// Forwarding rule name.
	RuleName *string `json:"RuleName,omitnil,omitempty" name:"RuleName"`

	// Forwarding rule type.
	RuleType *string `json:"RuleType,omitnil,omitempty" name:"RuleType"`

	// Private domain ID.
	ZoneId *string `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`

	// Endpoint ID.
	EndPointId *string `json:"EndPointId,omitnil,omitempty" name:"EndPointId"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type CreateForwardRuleResponse struct {
	*tchttp.BaseResponse
	Response *CreateForwardRuleResponseParams `json:"Response"`
}

func (r *CreateForwardRuleResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateForwardRuleResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreatePrivateDNSAccountRequestParams struct {
	// Private DNS account
	Account *PrivateDNSAccount `json:"Account,omitnil,omitempty" name:"Account"`
}

type CreatePrivateDNSAccountRequest struct {
	*tchttp.BaseRequest
	
	// Private DNS account
	Account *PrivateDNSAccount `json:"Account,omitnil,omitempty" name:"Account"`
}

func (r *CreatePrivateDNSAccountRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreatePrivateDNSAccountRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Account")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CreatePrivateDNSAccountRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreatePrivateDNSAccountResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type CreatePrivateDNSAccountResponse struct {
	*tchttp.BaseResponse
	Response *CreatePrivateDNSAccountResponseParams `json:"Response"`
}

func (r *CreatePrivateDNSAccountResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreatePrivateDNSAccountResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreatePrivateZoneRecordRequestParams struct {
	// Private domain ID
	ZoneId *string `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`

	// Record type. Valid values: "A", "AAAA", "CNAME", "MX", "TXT", "PTR"
	RecordType *string `json:"RecordType,omitnil,omitempty" name:"RecordType"`

	// Subdomain, such as "www", "m", and "@"
	SubDomain *string `json:"SubDomain,omitnil,omitempty" name:"SubDomain"`

	// Record value, such as IP: 192.168.10.2, CNAME: cname.qcloud.com., and MX: mail.qcloud.com.
	RecordValue *string `json:"RecordValue,omitnil,omitempty" name:"RecordValue"`

	// Record weight. Value range: 1–100
	Weight *int64 `json:"Weight,omitnil,omitempty" name:"Weight"`

	// MX priority, which is required when the record type is MX. Valid values: 5, 10, 15, 20, 30, 40, 50
	MX *int64 `json:"MX,omitnil,omitempty" name:"MX"`

	// Record cache time. The smaller the value, the faster the record will take effect. Value range: 1–86400s. Default value: 600
	TTL *int64 `json:"TTL,omitnil,omitempty" name:"TTL"`
}

type CreatePrivateZoneRecordRequest struct {
	*tchttp.BaseRequest
	
	// Private domain ID
	ZoneId *string `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`

	// Record type. Valid values: "A", "AAAA", "CNAME", "MX", "TXT", "PTR"
	RecordType *string `json:"RecordType,omitnil,omitempty" name:"RecordType"`

	// Subdomain, such as "www", "m", and "@"
	SubDomain *string `json:"SubDomain,omitnil,omitempty" name:"SubDomain"`

	// Record value, such as IP: 192.168.10.2, CNAME: cname.qcloud.com., and MX: mail.qcloud.com.
	RecordValue *string `json:"RecordValue,omitnil,omitempty" name:"RecordValue"`

	// Record weight. Value range: 1–100
	Weight *int64 `json:"Weight,omitnil,omitempty" name:"Weight"`

	// MX priority, which is required when the record type is MX. Valid values: 5, 10, 15, 20, 30, 40, 50
	MX *int64 `json:"MX,omitnil,omitempty" name:"MX"`

	// Record cache time. The smaller the value, the faster the record will take effect. Value range: 1–86400s. Default value: 600
	TTL *int64 `json:"TTL,omitnil,omitempty" name:"TTL"`
}

func (r *CreatePrivateZoneRecordRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreatePrivateZoneRecordRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ZoneId")
	delete(f, "RecordType")
	delete(f, "SubDomain")
	delete(f, "RecordValue")
	delete(f, "Weight")
	delete(f, "MX")
	delete(f, "TTL")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CreatePrivateZoneRecordRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreatePrivateZoneRecordResponseParams struct {
	// Record ID
	RecordId *string `json:"RecordId,omitnil,omitempty" name:"RecordId"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type CreatePrivateZoneRecordResponse struct {
	*tchttp.BaseResponse
	Response *CreatePrivateZoneRecordResponseParams `json:"Response"`
}

func (r *CreatePrivateZoneRecordResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreatePrivateZoneRecordResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreatePrivateZoneRequestParams struct {
	// Domain name, which must be in the format of standard TLD
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// Tags the private domain when it is created
	TagSet []*TagInfo `json:"TagSet,omitnil,omitempty" name:"TagSet"`

	// Associates the private domain to a VPC when it is created
	VpcSet []*VpcInfo `json:"VpcSet,omitnil,omitempty" name:"VpcSet"`

	// Remarks
	Remark *string `json:"Remark,omitnil,omitempty" name:"Remark"`

	// Whether to enable subdomain recursive DNS. Valid values: `ENABLED` (default) and `DISABLED`.
	DnsForwardStatus *string `json:"DnsForwardStatus,omitnil,omitempty" name:"DnsForwardStatus"`

	// Associates the private domain to a VPC when it is created
	//
	// Deprecated: Vpcs is deprecated.
	Vpcs []*VpcInfo `json:"Vpcs,omitnil,omitempty" name:"Vpcs"`

	// List of authorized accounts' VPCs to associate with the private domain
	AccountVpcSet []*AccountVpcInfo `json:"AccountVpcSet,omitnil,omitempty" name:"AccountVpcSet"`

	// Whether to enable CNAME flattening. Valid values: `ENABLED` (default) and `DISABLED`.
	CnameSpeedupStatus *string `json:"CnameSpeedupStatus,omitnil,omitempty" name:"CnameSpeedupStatus"`
}

type CreatePrivateZoneRequest struct {
	*tchttp.BaseRequest
	
	// Domain name, which must be in the format of standard TLD
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// Tags the private domain when it is created
	TagSet []*TagInfo `json:"TagSet,omitnil,omitempty" name:"TagSet"`

	// Associates the private domain to a VPC when it is created
	VpcSet []*VpcInfo `json:"VpcSet,omitnil,omitempty" name:"VpcSet"`

	// Remarks
	Remark *string `json:"Remark,omitnil,omitempty" name:"Remark"`

	// Whether to enable subdomain recursive DNS. Valid values: `ENABLED` (default) and `DISABLED`.
	DnsForwardStatus *string `json:"DnsForwardStatus,omitnil,omitempty" name:"DnsForwardStatus"`

	// Associates the private domain to a VPC when it is created
	Vpcs []*VpcInfo `json:"Vpcs,omitnil,omitempty" name:"Vpcs"`

	// List of authorized accounts' VPCs to associate with the private domain
	AccountVpcSet []*AccountVpcInfo `json:"AccountVpcSet,omitnil,omitempty" name:"AccountVpcSet"`

	// Whether to enable CNAME flattening. Valid values: `ENABLED` (default) and `DISABLED`.
	CnameSpeedupStatus *string `json:"CnameSpeedupStatus,omitnil,omitempty" name:"CnameSpeedupStatus"`
}

func (r *CreatePrivateZoneRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreatePrivateZoneRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Domain")
	delete(f, "TagSet")
	delete(f, "VpcSet")
	delete(f, "Remark")
	delete(f, "DnsForwardStatus")
	delete(f, "Vpcs")
	delete(f, "AccountVpcSet")
	delete(f, "CnameSpeedupStatus")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CreatePrivateZoneRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreatePrivateZoneResponseParams struct {
	// Private domain ID, such as zone-xxxxxx
	ZoneId *string `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`

	// Private domain
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type CreatePrivateZoneResponse struct {
	*tchttp.BaseResponse
	Response *CreatePrivateZoneResponseParams `json:"Response"`
}

func (r *CreatePrivateZoneResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreatePrivateZoneResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type DatePoint struct {
	// Time
	Date *string `json:"Date,omitnil,omitempty" name:"Date"`

	// Value
	Value *int64 `json:"Value,omitnil,omitempty" name:"Value"`
}

// Predefined struct for user
type DeleteEndPointRequestParams struct {
	// Endpoint ID.
	EndPointId *string `json:"EndPointId,omitnil,omitempty" name:"EndPointId"`
}

type DeleteEndPointRequest struct {
	*tchttp.BaseRequest
	
	// Endpoint ID.
	EndPointId *string `json:"EndPointId,omitnil,omitempty" name:"EndPointId"`
}

func (r *DeleteEndPointRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteEndPointRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "EndPointId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DeleteEndPointRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteEndPointResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DeleteEndPointResponse struct {
	*tchttp.BaseResponse
	Response *DeleteEndPointResponseParams `json:"Response"`
}

func (r *DeleteEndPointResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteEndPointResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteForwardRuleRequestParams struct {
	// Array of forwarding rule IDs.
	RuleIdSet []*string `json:"RuleIdSet,omitnil,omitempty" name:"RuleIdSet"`
}

type DeleteForwardRuleRequest struct {
	*tchttp.BaseRequest
	
	// Array of forwarding rule IDs.
	RuleIdSet []*string `json:"RuleIdSet,omitnil,omitempty" name:"RuleIdSet"`
}

func (r *DeleteForwardRuleRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteForwardRuleRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "RuleIdSet")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DeleteForwardRuleRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteForwardRuleResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DeleteForwardRuleResponse struct {
	*tchttp.BaseResponse
	Response *DeleteForwardRuleResponseParams `json:"Response"`
}

func (r *DeleteForwardRuleResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteForwardRuleResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeletePrivateZoneRecordRequestParams struct {
	// Private domain ID
	ZoneId *string `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`

	// Record ID
	RecordId *string `json:"RecordId,omitnil,omitempty" name:"RecordId"`

	// Array of record IDs. `RecordId` takes precedence.
	RecordIdSet []*string `json:"RecordIdSet,omitnil,omitempty" name:"RecordIdSet"`
}

type DeletePrivateZoneRecordRequest struct {
	*tchttp.BaseRequest
	
	// Private domain ID
	ZoneId *string `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`

	// Record ID
	RecordId *string `json:"RecordId,omitnil,omitempty" name:"RecordId"`

	// Array of record IDs. `RecordId` takes precedence.
	RecordIdSet []*string `json:"RecordIdSet,omitnil,omitempty" name:"RecordIdSet"`
}

func (r *DeletePrivateZoneRecordRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeletePrivateZoneRecordRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ZoneId")
	delete(f, "RecordId")
	delete(f, "RecordIdSet")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DeletePrivateZoneRecordRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeletePrivateZoneRecordResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DeletePrivateZoneRecordResponse struct {
	*tchttp.BaseResponse
	Response *DeletePrivateZoneRecordResponseParams `json:"Response"`
}

func (r *DeletePrivateZoneRecordResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeletePrivateZoneRecordResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeAccountVpcListRequestParams struct {
	// UIN of account
	AccountUin *string `json:"AccountUin,omitnil,omitempty" name:"AccountUin"`

	// Pagination offset, starting from 0
	Offset *int64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// Number of entries per page. Maximum value: `100`. Default value: `20`
	Limit *int64 `json:"Limit,omitnil,omitempty" name:"Limit"`

	// Filter parameters
	Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`
}

type DescribeAccountVpcListRequest struct {
	*tchttp.BaseRequest
	
	// UIN of account
	AccountUin *string `json:"AccountUin,omitnil,omitempty" name:"AccountUin"`

	// Pagination offset, starting from 0
	Offset *int64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// Number of entries per page. Maximum value: `100`. Default value: `20`
	Limit *int64 `json:"Limit,omitnil,omitempty" name:"Limit"`

	// Filter parameters
	Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`
}

func (r *DescribeAccountVpcListRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeAccountVpcListRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "AccountUin")
	delete(f, "Offset")
	delete(f, "Limit")
	delete(f, "Filters")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeAccountVpcListRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeAccountVpcListResponseParams struct {
	// Number of VPCs
	TotalCount *int64 `json:"TotalCount,omitnil,omitempty" name:"TotalCount"`

	// VPC list
	VpcSet []*AccountVpcInfoOut `json:"VpcSet,omitnil,omitempty" name:"VpcSet"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeAccountVpcListResponse struct {
	*tchttp.BaseResponse
	Response *DescribeAccountVpcListResponseParams `json:"Response"`
}

func (r *DescribeAccountVpcListResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeAccountVpcListResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeAuditLogRequestParams struct {
	// Request volume statistics start time
	TimeRangeBegin *string `json:"TimeRangeBegin,omitnil,omitempty" name:"TimeRangeBegin"`

	// Filter parameter. Valid values: ZoneId (private domain ID), Domain (private domain), OperatorUin (operator account ID)
	Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`

	// Request volume statistics end time
	TimeRangeEnd *string `json:"TimeRangeEnd,omitnil,omitempty" name:"TimeRangeEnd"`

	// Pagination offset, starting from 0
	Offset *int64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// Number of entries per page. Maximum value: 100. Default value: 20
	Limit *int64 `json:"Limit,omitnil,omitempty" name:"Limit"`
}

type DescribeAuditLogRequest struct {
	*tchttp.BaseRequest
	
	// Request volume statistics start time
	TimeRangeBegin *string `json:"TimeRangeBegin,omitnil,omitempty" name:"TimeRangeBegin"`

	// Filter parameter. Valid values: ZoneId (private domain ID), Domain (private domain), OperatorUin (operator account ID)
	Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`

	// Request volume statistics end time
	TimeRangeEnd *string `json:"TimeRangeEnd,omitnil,omitempty" name:"TimeRangeEnd"`

	// Pagination offset, starting from 0
	Offset *int64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// Number of entries per page. Maximum value: 100. Default value: 20
	Limit *int64 `json:"Limit,omitnil,omitempty" name:"Limit"`
}

func (r *DescribeAuditLogRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeAuditLogRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "TimeRangeBegin")
	delete(f, "Filters")
	delete(f, "TimeRangeEnd")
	delete(f, "Offset")
	delete(f, "Limit")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeAuditLogRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeAuditLogResponseParams struct {
	// List of operation logs
	Data []*AuditLog `json:"Data,omitnil,omitempty" name:"Data"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeAuditLogResponse struct {
	*tchttp.BaseResponse
	Response *DescribeAuditLogResponseParams `json:"Response"`
}

func (r *DescribeAuditLogResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeAuditLogResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeDashboardRequestParams struct {

}

type DescribeDashboardRequest struct {
	*tchttp.BaseRequest
	
}

func (r *DescribeDashboardRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeDashboardRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeDashboardRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeDashboardResponseParams struct {
	// Total number of private domain DNS records
	ZoneTotal *int64 `json:"ZoneTotal,omitnil,omitempty" name:"ZoneTotal"`

	// Number of VPCs associated with private domain
	ZoneVpcCount *int64 `json:"ZoneVpcCount,omitnil,omitempty" name:"ZoneVpcCount"`

	// Total number of historical requests
	RequestTotalCount *int64 `json:"RequestTotalCount,omitnil,omitempty" name:"RequestTotalCount"`

	// Traffic package usage
	FlowUsage []*FlowUsage `json:"FlowUsage,omitnil,omitempty" name:"FlowUsage"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeDashboardResponse struct {
	*tchttp.BaseResponse
	Response *DescribeDashboardResponseParams `json:"Response"`
}

func (r *DescribeDashboardResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeDashboardResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeEndPointListRequestParams struct {
	// Pagination offset, starting from 0.
	Offset *int64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// Pagination limit. Maximum value: 100. Default value: 20.
	Limit *int64 `json:"Limit,omitnil,omitempty" name:"Limit"`

	// Filter parameters. Valid values: EndPointName, EndPointId, EndPointServiceId, and EndPointVip.
	Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`
}

type DescribeEndPointListRequest struct {
	*tchttp.BaseRequest
	
	// Pagination offset, starting from 0.
	Offset *int64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// Pagination limit. Maximum value: 100. Default value: 20.
	Limit *int64 `json:"Limit,omitnil,omitempty" name:"Limit"`

	// Filter parameters. Valid values: EndPointName, EndPointId, EndPointServiceId, and EndPointVip.
	Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`
}

func (r *DescribeEndPointListRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeEndPointListRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Offset")
	delete(f, "Limit")
	delete(f, "Filters")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeEndPointListRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeEndPointListResponseParams struct {
	// Total number of endpoints.
	TotalCount *int64 `json:"TotalCount,omitnil,omitempty" name:"TotalCount"`

	// Endpoint list.
	// Note: This field may return null, indicating that no valid values can be obtained.
	EndPointSet []*EndPointInfo `json:"EndPointSet,omitnil,omitempty" name:"EndPointSet"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeEndPointListResponse struct {
	*tchttp.BaseResponse
	Response *DescribeEndPointListResponseParams `json:"Response"`
}

func (r *DescribeEndPointListResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeEndPointListResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeEndPointRegionRequestParams struct {

}

type DescribeEndPointRegionRequest struct {
	*tchttp.BaseRequest
	
}

func (r *DescribeEndPointRegionRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeEndPointRegionRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeEndPointRegionRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeEndPointRegionResponseParams struct {
	// Region array.
	RegionSet []*RegionInfo `json:"RegionSet,omitnil,omitempty" name:"RegionSet"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeEndPointRegionResponse struct {
	*tchttp.BaseResponse
	Response *DescribeEndPointRegionResponseParams `json:"Response"`
}

func (r *DescribeEndPointRegionResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeEndPointRegionResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeExtendEndpointListRequestParams struct {
	// Pagination offset, starting from 0.
	Offset *int64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// Pagination limit. Maximum value: 100. Default value: 20.
	Limit *int64 `json:"Limit,omitnil,omitempty" name:"Limit"`

	// Filter parameters. Valid values: EndpointName, EndpointId.
	Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`
}

type DescribeExtendEndpointListRequest struct {
	*tchttp.BaseRequest
	
	// Pagination offset, starting from 0.
	Offset *int64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// Pagination limit. Maximum value: 100. Default value: 20.
	Limit *int64 `json:"Limit,omitnil,omitempty" name:"Limit"`

	// Filter parameters. Valid values: EndpointName, EndpointId.
	Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`
}

func (r *DescribeExtendEndpointListRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeExtendEndpointListRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Offset")
	delete(f, "Limit")
	delete(f, "Filters")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeExtendEndpointListRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeExtendEndpointListResponseParams struct {
	// Total number of endpoints.
	TotalCount *int64 `json:"TotalCount,omitnil,omitempty" name:"TotalCount"`

	// Endpoint list.
	OutboundEndpointSet []*OutboundEndpoint `json:"OutboundEndpointSet,omitnil,omitempty" name:"OutboundEndpointSet"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeExtendEndpointListResponse struct {
	*tchttp.BaseResponse
	Response *DescribeExtendEndpointListResponseParams `json:"Response"`
}

func (r *DescribeExtendEndpointListResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeExtendEndpointListResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeForwardRuleListRequestParams struct {
	// Pagination offset, starting from 0.
	Offset *int64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// Pagination limit. Maximum value: 100. Default value: 20.
	Limit *int64 `json:"Limit,omitnil,omitempty" name:"Limit"`

	// Filter parameters.
	Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`
}

type DescribeForwardRuleListRequest struct {
	*tchttp.BaseRequest
	
	// Pagination offset, starting from 0.
	Offset *int64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// Pagination limit. Maximum value: 100. Default value: 20.
	Limit *int64 `json:"Limit,omitnil,omitempty" name:"Limit"`

	// Filter parameters.
	Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`
}

func (r *DescribeForwardRuleListRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeForwardRuleListRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Offset")
	delete(f, "Limit")
	delete(f, "Filters")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeForwardRuleListRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeForwardRuleListResponseParams struct {
	// Number of private domains.
	TotalCount *int64 `json:"TotalCount,omitnil,omitempty" name:"TotalCount"`

	// Private domain list.
	ForwardRuleSet []*ForwardRule `json:"ForwardRuleSet,omitnil,omitempty" name:"ForwardRuleSet"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeForwardRuleListResponse struct {
	*tchttp.BaseResponse
	Response *DescribeForwardRuleListResponseParams `json:"Response"`
}

func (r *DescribeForwardRuleListResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeForwardRuleListResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeForwardRuleRequestParams struct {
	// Forwarding rule ID.
	RuleId *string `json:"RuleId,omitnil,omitempty" name:"RuleId"`
}

type DescribeForwardRuleRequest struct {
	*tchttp.BaseRequest
	
	// Forwarding rule ID.
	RuleId *string `json:"RuleId,omitnil,omitempty" name:"RuleId"`
}

func (r *DescribeForwardRuleRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeForwardRuleRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "RuleId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeForwardRuleRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeForwardRuleResponseParams struct {
	// Forwarding rule details.
	ForwardRule *ForwardRule `json:"ForwardRule,omitnil,omitempty" name:"ForwardRule"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeForwardRuleResponse struct {
	*tchttp.BaseResponse
	Response *DescribeForwardRuleResponseParams `json:"Response"`
}

func (r *DescribeForwardRuleResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeForwardRuleResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribePrivateDNSAccountListRequestParams struct {
	// Pagination offset, starting from `0`
	Offset *int64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// Number of entries per page. Maximum value: `100`. Default value: `20`
	Limit *int64 `json:"Limit,omitnil,omitempty" name:"Limit"`

	// Filter parameters
	Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`
}

type DescribePrivateDNSAccountListRequest struct {
	*tchttp.BaseRequest
	
	// Pagination offset, starting from `0`
	Offset *int64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// Number of entries per page. Maximum value: `100`. Default value: `20`
	Limit *int64 `json:"Limit,omitnil,omitempty" name:"Limit"`

	// Filter parameters
	Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`
}

func (r *DescribePrivateDNSAccountListRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribePrivateDNSAccountListRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Offset")
	delete(f, "Limit")
	delete(f, "Filters")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribePrivateDNSAccountListRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribePrivateDNSAccountListResponseParams struct {
	// Number of Private DNS accounts
	TotalCount *int64 `json:"TotalCount,omitnil,omitempty" name:"TotalCount"`

	// List of Private DNS accounts
	AccountSet []*PrivateDNSAccount `json:"AccountSet,omitnil,omitempty" name:"AccountSet"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribePrivateDNSAccountListResponse struct {
	*tchttp.BaseResponse
	Response *DescribePrivateDNSAccountListResponseParams `json:"Response"`
}

func (r *DescribePrivateDNSAccountListResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribePrivateDNSAccountListResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribePrivateZoneListRequestParams struct {
	// Pagination offset, starting from 0.
	Offset *int64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// Pagination limit. Maximum value: 100. Default value: 20.
	Limit *int64 `json:"Limit,omitnil,omitempty" name:"Limit"`

	// Filter parameters.
	Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`
}

type DescribePrivateZoneListRequest struct {
	*tchttp.BaseRequest
	
	// Pagination offset, starting from 0.
	Offset *int64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// Pagination limit. Maximum value: 100. Default value: 20.
	Limit *int64 `json:"Limit,omitnil,omitempty" name:"Limit"`

	// Filter parameters.
	Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`
}

func (r *DescribePrivateZoneListRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribePrivateZoneListRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Offset")
	delete(f, "Limit")
	delete(f, "Filters")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribePrivateZoneListRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribePrivateZoneListResponseParams struct {
	// Number of private domains.
	TotalCount *int64 `json:"TotalCount,omitnil,omitempty" name:"TotalCount"`

	// Private domain list.
	PrivateZoneSet []*PrivateZone `json:"PrivateZoneSet,omitnil,omitempty" name:"PrivateZoneSet"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribePrivateZoneListResponse struct {
	*tchttp.BaseResponse
	Response *DescribePrivateZoneListResponseParams `json:"Response"`
}

func (r *DescribePrivateZoneListResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribePrivateZoneListResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribePrivateZoneRecordListRequestParams struct {
	// Private domain ID: zone-xxxxxx
	ZoneId *string `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`

	// Filter parameter
	Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`

	// Pagination offset, starting from 0
	Offset *int64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// Number of entries per page. Maximum value: 100. Default value: 20
	Limit *int64 `json:"Limit,omitnil,omitempty" name:"Limit"`
}

type DescribePrivateZoneRecordListRequest struct {
	*tchttp.BaseRequest
	
	// Private domain ID: zone-xxxxxx
	ZoneId *string `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`

	// Filter parameter
	Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`

	// Pagination offset, starting from 0
	Offset *int64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// Number of entries per page. Maximum value: 100. Default value: 20
	Limit *int64 `json:"Limit,omitnil,omitempty" name:"Limit"`
}

func (r *DescribePrivateZoneRecordListRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribePrivateZoneRecordListRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ZoneId")
	delete(f, "Filters")
	delete(f, "Offset")
	delete(f, "Limit")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribePrivateZoneRecordListRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribePrivateZoneRecordListResponseParams struct {
	// Number of DNS records
	TotalCount *int64 `json:"TotalCount,omitnil,omitempty" name:"TotalCount"`

	// List of DNS records
	RecordSet []*PrivateZoneRecord `json:"RecordSet,omitnil,omitempty" name:"RecordSet"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribePrivateZoneRecordListResponse struct {
	*tchttp.BaseResponse
	Response *DescribePrivateZoneRecordListResponseParams `json:"Response"`
}

func (r *DescribePrivateZoneRecordListResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribePrivateZoneRecordListResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribePrivateZoneServiceRequestParams struct {

}

type DescribePrivateZoneServiceRequest struct {
	*tchttp.BaseRequest
	
}

func (r *DescribePrivateZoneServiceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribePrivateZoneServiceRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribePrivateZoneServiceRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribePrivateZoneServiceResponseParams struct {
	// Private DNS service activation status. Valid values: ENABLED, DISABLED
	ServiceStatus *string `json:"ServiceStatus,omitnil,omitempty" name:"ServiceStatus"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribePrivateZoneServiceResponse struct {
	*tchttp.BaseResponse
	Response *DescribePrivateZoneServiceResponseParams `json:"Response"`
}

func (r *DescribePrivateZoneServiceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribePrivateZoneServiceResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeQuotaUsageRequestParams struct {

}

type DescribeQuotaUsageRequest struct {
	*tchttp.BaseRequest
	
}

func (r *DescribeQuotaUsageRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeQuotaUsageRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeQuotaUsageRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeQuotaUsageResponseParams struct {
	// TLD quota usage
	TldQuota *TldQuota `json:"TldQuota,omitnil,omitempty" name:"TldQuota"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeQuotaUsageResponse struct {
	*tchttp.BaseResponse
	Response *DescribeQuotaUsageResponseParams `json:"Response"`
}

func (r *DescribeQuotaUsageResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeQuotaUsageResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeRecordRequestParams struct {
	// Private domain ID.
	ZoneId *string `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`

	// Record ID
	RecordId *string `json:"RecordId,omitnil,omitempty" name:"RecordId"`
}

type DescribeRecordRequest struct {
	*tchttp.BaseRequest
	
	// Private domain ID.
	ZoneId *string `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`

	// Record ID
	RecordId *string `json:"RecordId,omitnil,omitempty" name:"RecordId"`
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
	delete(f, "ZoneId")
	delete(f, "RecordId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeRecordRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeRecordResponseParams struct {
	// Record information.
	RecordInfo *RecordInfo `json:"RecordInfo,omitnil,omitempty" name:"RecordInfo"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
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
type DescribeRequestDataRequestParams struct {
	// Request volume statistics start time in the format of 2020-11-22 00:00:00
	TimeRangeBegin *string `json:"TimeRangeBegin,omitnil,omitempty" name:"TimeRangeBegin"`

	// Filter parameter:
	Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`

	// Request volume statistics end time in the format of 2020-11-22 23:59:59
	TimeRangeEnd *string `json:"TimeRangeEnd,omitnil,omitempty" name:"TimeRangeEnd"`
}

type DescribeRequestDataRequest struct {
	*tchttp.BaseRequest
	
	// Request volume statistics start time in the format of 2020-11-22 00:00:00
	TimeRangeBegin *string `json:"TimeRangeBegin,omitnil,omitempty" name:"TimeRangeBegin"`

	// Filter parameter:
	Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`

	// Request volume statistics end time in the format of 2020-11-22 23:59:59
	TimeRangeEnd *string `json:"TimeRangeEnd,omitnil,omitempty" name:"TimeRangeEnd"`
}

func (r *DescribeRequestDataRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeRequestDataRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "TimeRangeBegin")
	delete(f, "Filters")
	delete(f, "TimeRangeEnd")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeRequestDataRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeRequestDataResponseParams struct {
	// Request volume statistics table
	Data []*MetricData `json:"Data,omitnil,omitempty" name:"Data"`

	// Request volume unit time. Valid values: Day, Hour
	Interval *string `json:"Interval,omitnil,omitempty" name:"Interval"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeRequestDataResponse struct {
	*tchttp.BaseResponse
	Response *DescribeRequestDataResponseParams `json:"Response"`
}

func (r *DescribeRequestDataResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeRequestDataResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type EndPointInfo struct {
	// Endpoint ID.
	EndPointId *string `json:"EndPointId,omitnil,omitempty" name:"EndPointId"`

	// Endpoint name.
	EndPointName *string `json:"EndPointName,omitnil,omitempty" name:"EndPointName"`

	// Endpoint service ID.
	EndPointServiceId *string `json:"EndPointServiceId,omitnil,omitempty" name:"EndPointServiceId"`

	// VIP list of the endpoint.
	EndPointVipSet []*string `json:"EndPointVipSet,omitnil,omitempty" name:"EndPointVipSet"`

	// ap-guangzhou
	// Note: This field may return null, indicating that no valid values can be obtained.
	RegionCode *string `json:"RegionCode,omitnil,omitempty" name:"RegionCode"`

	// Tag key-value pair collection.
	// Note: This field may return null, indicating that no valid values can be obtained.
	Tags []*TagInfo `json:"Tags,omitnil,omitempty" name:"Tags"`
}

type EndpointService struct {
	// Specifies the forwarding target IP network access type.
	// CLB: Specifies that the forwarding IP is the private CLB VIP.
	// CCN: Specifies forwarding IP through CCN routing.
	AccessType *string `json:"AccessType,omitnil,omitempty" name:"AccessType"`

	// Specifies the forwarding target IP address.
	Pip *string `json:"Pip,omitnil,omitempty" name:"Pip"`

	// Specifies the forwarding IP port number.
	Pport *int64 `json:"Pport,omitnil,omitempty" name:"Pport"`

	// Specifies the unique VPC ID.
	VpcId *string `json:"VpcId,omitnil,omitempty" name:"VpcId"`

	// Specifies the forwarding target IP proxy IP.
	Vip *string `json:"Vip,omitnil,omitempty" name:"Vip"`

	// Specifies the forwarding target IP proxy port.
	Vport *int64 `json:"Vport,omitnil,omitempty" name:"Vport"`

	// Specifies the forwarding target IP protocol.
	Proto *string `json:"Proto,omitnil,omitempty" name:"Proto"`

	// Specifies the unique subnet ID.
	// Required if the access type is CCN.
	SubnetId *string `json:"SubnetId,omitnil,omitempty" name:"SubnetId"`

	// ccn id
	// Required if the access type is CCN.
	AccessGatewayId *string `json:"AccessGatewayId,omitnil,omitempty" name:"AccessGatewayId"`

	// The SNAT CIDR block of the outbound endpoint.
	SnatVipCidr *string `json:"SnatVipCidr,omitnil,omitempty" name:"SnatVipCidr"`

	// The SNAT IP list of the outbound endpoint.
	SnatVipSet *string `json:"SnatVipSet,omitnil,omitempty" name:"SnatVipSet"`

	// The region of the outbound endpoint service.
	Region *string `json:"Region,omitnil,omitempty" name:"Region"`
}

type Filter struct {
	// Parameter name
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// Array of parameter values
	Values []*string `json:"Values,omitnil,omitempty" name:"Values"`
}

type FlowUsage struct {
	// Traffic package type, Valid values: ZONE (private domain); TRAFFIC (DNS traffic package)
	FlowType *string `json:"FlowType,omitnil,omitempty" name:"FlowType"`

	// Traffic package quota
	TotalQuantity *int64 `json:"TotalQuantity,omitnil,omitempty" name:"TotalQuantity"`

	// Available quota of traffic package
	AvailableQuantity *int64 `json:"AvailableQuantity,omitnil,omitempty" name:"AvailableQuantity"`
}

type ForwardIp struct {
	// Forwarding target IP network access type.
	// CLB: The forwarding IP is the internal CLB VIP.
	// CCN: Forwarding IP through CCN routing.
	AccessType *string `json:"AccessType,omitnil,omitempty" name:"AccessType"`

	// Forwarding target IP address.
	Host *string `json:"Host,omitnil,omitempty" name:"Host"`

	// Specifies the forwarding IP port number.
	Port *int64 `json:"Port,omitnil,omitempty" name:"Port"`

	// Specifies the number of outbound endpoints.
	// Minimum 1, maximum 6.
	IpNum *int64 `json:"IpNum,omitnil,omitempty" name:"IpNum"`

	// Unique VPC ID.
	VpcId *string `json:"VpcId,omitnil,omitempty" name:"VpcId"`

	// Unique subnet ID.
	// Required when the access type is CCN.
	SubnetId *string `json:"SubnetId,omitnil,omitempty" name:"SubnetId"`

	// ccn id
	// Required when the access type is CCN.
	AccessGatewayId *string `json:"AccessGatewayId,omitnil,omitempty" name:"AccessGatewayId"`
}

type ForwardRule struct {
	// Private domain name.
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// Forwarding rule name.
	RuleName *string `json:"RuleName,omitnil,omitempty" name:"RuleName"`

	// Rule ID
	RuleId *string `json:"RuleId,omitnil,omitempty" name:"RuleId"`

	// Forwarding rule type. DOWN: From cloud to off-cloud; UP: From off-cloud to cloud.
	RuleType *string `json:"RuleType,omitnil,omitempty" name:"RuleType"`

	// Creation time
	CreatedAt *string `json:"CreatedAt,omitnil,omitempty" name:"CreatedAt"`

	// Update time
	UpdatedAt *string `json:"UpdatedAt,omitnil,omitempty" name:"UpdatedAt"`

	// Endpoint name.
	EndPointName *string `json:"EndPointName,omitnil,omitempty" name:"EndPointName"`

	// Endpoint ID.
	EndPointId *string `json:"EndPointId,omitnil,omitempty" name:"EndPointId"`

	// Forwarding address.
	ForwardAddress []*string `json:"ForwardAddress,omitnil,omitempty" name:"ForwardAddress"`

	// List of VPCs bound to the private domain.
	VpcSet []*VpcInfo `json:"VpcSet,omitnil,omitempty" name:"VpcSet"`

	// ID of the bound private domain.
	ZoneId *string `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`

	// Tag
	Tags []*TagInfo `json:"Tags,omitnil,omitempty" name:"Tags"`
}

type MetricData struct {
	// Resource description
	Resource *string `json:"Resource,omitnil,omitempty" name:"Resource"`

	// Table name
	Metric *string `json:"Metric,omitnil,omitempty" name:"Metric"`

	// Table data
	DataSet []*DatePoint `json:"DataSet,omitnil,omitempty" name:"DataSet"`

	// The total number of requests within the query scope.
	// Note: This field may return null, indicating that no valid value can be obtained.
	MetricCount *int64 `json:"MetricCount,omitnil,omitempty" name:"MetricCount"`
}

// Predefined struct for user
type ModifyForwardRuleRequestParams struct {
	// Forwarding rule ID.
	RuleId *string `json:"RuleId,omitnil,omitempty" name:"RuleId"`

	// Forwarding rule name.
	RuleName *string `json:"RuleName,omitnil,omitempty" name:"RuleName"`

	// Endpoint ID.
	EndPointId *string `json:"EndPointId,omitnil,omitempty" name:"EndPointId"`
}

type ModifyForwardRuleRequest struct {
	*tchttp.BaseRequest
	
	// Forwarding rule ID.
	RuleId *string `json:"RuleId,omitnil,omitempty" name:"RuleId"`

	// Forwarding rule name.
	RuleName *string `json:"RuleName,omitnil,omitempty" name:"RuleName"`

	// Endpoint ID.
	EndPointId *string `json:"EndPointId,omitnil,omitempty" name:"EndPointId"`
}

func (r *ModifyForwardRuleRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyForwardRuleRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "RuleId")
	delete(f, "RuleName")
	delete(f, "EndPointId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ModifyForwardRuleRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyForwardRuleResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ModifyForwardRuleResponse struct {
	*tchttp.BaseResponse
	Response *ModifyForwardRuleResponseParams `json:"Response"`
}

func (r *ModifyForwardRuleResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyForwardRuleResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyPrivateZoneRecordRequestParams struct {
	// Private domain ID
	ZoneId *string `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`

	// Record ID
	RecordId *string `json:"RecordId,omitnil,omitempty" name:"RecordId"`

	// Record type. Valid values: "A", "AAAA", "CNAME", "MX", "TXT", "PTR"
	RecordType *string `json:"RecordType,omitnil,omitempty" name:"RecordType"`

	// Subdomain, such as "www", "m", and "@"
	SubDomain *string `json:"SubDomain,omitnil,omitempty" name:"SubDomain"`

	// Record value, such as IP: 192.168.10.2, CNAME: cname.qcloud.com., and MX: mail.qcloud.com.
	RecordValue *string `json:"RecordValue,omitnil,omitempty" name:"RecordValue"`

	// Record weight. Value range: 1–100
	Weight *int64 `json:"Weight,omitnil,omitempty" name:"Weight"`

	// MX priority, which is required when the record type is MX. Valid values: 5, 10, 15, 20, 30, 40, 50
	MX *int64 `json:"MX,omitnil,omitempty" name:"MX"`

	// Record cache time. The smaller the value, the faster the record will take effect. Value range: 1–86400s. Default value: 600
	TTL *int64 `json:"TTL,omitnil,omitempty" name:"TTL"`
}

type ModifyPrivateZoneRecordRequest struct {
	*tchttp.BaseRequest
	
	// Private domain ID
	ZoneId *string `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`

	// Record ID
	RecordId *string `json:"RecordId,omitnil,omitempty" name:"RecordId"`

	// Record type. Valid values: "A", "AAAA", "CNAME", "MX", "TXT", "PTR"
	RecordType *string `json:"RecordType,omitnil,omitempty" name:"RecordType"`

	// Subdomain, such as "www", "m", and "@"
	SubDomain *string `json:"SubDomain,omitnil,omitempty" name:"SubDomain"`

	// Record value, such as IP: 192.168.10.2, CNAME: cname.qcloud.com., and MX: mail.qcloud.com.
	RecordValue *string `json:"RecordValue,omitnil,omitempty" name:"RecordValue"`

	// Record weight. Value range: 1–100
	Weight *int64 `json:"Weight,omitnil,omitempty" name:"Weight"`

	// MX priority, which is required when the record type is MX. Valid values: 5, 10, 15, 20, 30, 40, 50
	MX *int64 `json:"MX,omitnil,omitempty" name:"MX"`

	// Record cache time. The smaller the value, the faster the record will take effect. Value range: 1–86400s. Default value: 600
	TTL *int64 `json:"TTL,omitnil,omitempty" name:"TTL"`
}

func (r *ModifyPrivateZoneRecordRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyPrivateZoneRecordRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ZoneId")
	delete(f, "RecordId")
	delete(f, "RecordType")
	delete(f, "SubDomain")
	delete(f, "RecordValue")
	delete(f, "Weight")
	delete(f, "MX")
	delete(f, "TTL")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ModifyPrivateZoneRecordRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyPrivateZoneRecordResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ModifyPrivateZoneRecordResponse struct {
	*tchttp.BaseResponse
	Response *ModifyPrivateZoneRecordResponseParams `json:"Response"`
}

func (r *ModifyPrivateZoneRecordResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyPrivateZoneRecordResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyPrivateZoneRequestParams struct {
	// Private domain ID
	ZoneId *string `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`

	// Remarks
	Remark *string `json:"Remark,omitnil,omitempty" name:"Remark"`

	// Whether to enable subdomain recursive DNS. Valid values: ENABLED, DISABLED
	DnsForwardStatus *string `json:"DnsForwardStatus,omitnil,omitempty" name:"DnsForwardStatus"`

	// Whether to enable CNAME flattening. Valid values: `ENABLED` and `DISABLED`.
	CnameSpeedupStatus *string `json:"CnameSpeedupStatus,omitnil,omitempty" name:"CnameSpeedupStatus"`
}

type ModifyPrivateZoneRequest struct {
	*tchttp.BaseRequest
	
	// Private domain ID
	ZoneId *string `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`

	// Remarks
	Remark *string `json:"Remark,omitnil,omitempty" name:"Remark"`

	// Whether to enable subdomain recursive DNS. Valid values: ENABLED, DISABLED
	DnsForwardStatus *string `json:"DnsForwardStatus,omitnil,omitempty" name:"DnsForwardStatus"`

	// Whether to enable CNAME flattening. Valid values: `ENABLED` and `DISABLED`.
	CnameSpeedupStatus *string `json:"CnameSpeedupStatus,omitnil,omitempty" name:"CnameSpeedupStatus"`
}

func (r *ModifyPrivateZoneRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyPrivateZoneRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ZoneId")
	delete(f, "Remark")
	delete(f, "DnsForwardStatus")
	delete(f, "CnameSpeedupStatus")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ModifyPrivateZoneRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyPrivateZoneResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ModifyPrivateZoneResponse struct {
	*tchttp.BaseResponse
	Response *ModifyPrivateZoneResponseParams `json:"Response"`
}

func (r *ModifyPrivateZoneResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyPrivateZoneResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyPrivateZoneVpcRequestParams struct {
	// Private domain ID
	ZoneId *string `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`

	// List of all VPCs associated with private domain
	VpcSet []*VpcInfo `json:"VpcSet,omitnil,omitempty" name:"VpcSet"`

	// List of authorized accounts' VPCs to associate with the private domain
	AccountVpcSet []*AccountVpcInfo `json:"AccountVpcSet,omitnil,omitempty" name:"AccountVpcSet"`
}

type ModifyPrivateZoneVpcRequest struct {
	*tchttp.BaseRequest
	
	// Private domain ID
	ZoneId *string `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`

	// List of all VPCs associated with private domain
	VpcSet []*VpcInfo `json:"VpcSet,omitnil,omitempty" name:"VpcSet"`

	// List of authorized accounts' VPCs to associate with the private domain
	AccountVpcSet []*AccountVpcInfo `json:"AccountVpcSet,omitnil,omitempty" name:"AccountVpcSet"`
}

func (r *ModifyPrivateZoneVpcRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyPrivateZoneVpcRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ZoneId")
	delete(f, "VpcSet")
	delete(f, "AccountVpcSet")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ModifyPrivateZoneVpcRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyPrivateZoneVpcResponseParams struct {
	// Private domain ID, such as zone-xxxxxx
	ZoneId *string `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`

	// List of VPCs associated with domain
	VpcSet []*VpcInfo `json:"VpcSet,omitnil,omitempty" name:"VpcSet"`

	// List of authorized accounts' VPCs associated with the private domain
	AccountVpcSet []*AccountVpcInfoOutput `json:"AccountVpcSet,omitnil,omitempty" name:"AccountVpcSet"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ModifyPrivateZoneVpcResponse struct {
	*tchttp.BaseResponse
	Response *ModifyPrivateZoneVpcResponseParams `json:"Response"`
}

func (r *ModifyPrivateZoneVpcResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyPrivateZoneVpcResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyRecordsStatusRequestParams struct {
	// The private domain ID
	ZoneId *string `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`

	// The DNS record IDs.
	RecordIds []*int64 `json:"RecordIds,omitnil,omitempty" name:"RecordIds"`

	// `enabled`: Enable; `disabled`: Disable.
	Status *string `json:"Status,omitnil,omitempty" name:"Status"`
}

type ModifyRecordsStatusRequest struct {
	*tchttp.BaseRequest
	
	// The private domain ID
	ZoneId *string `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`

	// The DNS record IDs.
	RecordIds []*int64 `json:"RecordIds,omitnil,omitempty" name:"RecordIds"`

	// `enabled`: Enable; `disabled`: Disable.
	Status *string `json:"Status,omitnil,omitempty" name:"Status"`
}

func (r *ModifyRecordsStatusRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyRecordsStatusRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ZoneId")
	delete(f, "RecordIds")
	delete(f, "Status")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ModifyRecordsStatusRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyRecordsStatusResponseParams struct {
	// The private domain ID
	ZoneId *string `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`

	// The DNS record IDs.
	RecordIds []*int64 `json:"RecordIds,omitnil,omitempty" name:"RecordIds"`

	// `enabled`: Enabled; `disabled`: Disabled.
	Status *string `json:"Status,omitnil,omitempty" name:"Status"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ModifyRecordsStatusResponse struct {
	*tchttp.BaseResponse
	Response *ModifyRecordsStatusResponseParams `json:"Response"`
}

func (r *ModifyRecordsStatusResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyRecordsStatusResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type OutboundEndpoint struct {
	// Outbound endpoint ID.
	EndpointId *string `json:"EndpointId,omitnil,omitempty" name:"EndpointId"`

	// Outbound endpoint name.
	EndpointName *string `json:"EndpointName,omitnil,omitempty" name:"EndpointName"`

	// The region of the outbound endpoint.
	Region *string `json:"Region,omitnil,omitempty" name:"Region"`

	// Tag
	Tags []*TagInfo `json:"Tags,omitnil,omitempty" name:"Tags"`

	// Outbound endpoint information.
	// Returned only when the forwarding architecture is V2R.
	EndpointServiceSet []*EndpointService `json:"EndpointServiceSet,omitnil,omitempty" name:"EndpointServiceSet"`

	// Forwarding link architecture.
	// V2V: privatelink
	// V2R: jnsgw
	ForwardLinkArch *string `json:"ForwardLinkArch,omitnil,omitempty" name:"ForwardLinkArch"`

	// Endpoint service ID.
	// 
	// Returned only when the forwarding architecture is V2V.
	EndPointServiceId *string `json:"EndPointServiceId,omitnil,omitempty" name:"EndPointServiceId"`

	// VIP list of the endpoint.
	// 
	// Returned only when the forwarding architecture is V2V.
	EndPointVipSet []*string `json:"EndPointVipSet,omitnil,omitempty" name:"EndPointVipSet"`
}

type PrivateDNSAccount struct {
	// Root account UIN
	Uin *string `json:"Uin,omitnil,omitempty" name:"Uin"`

	// Root account name
	Account *string `json:"Account,omitnil,omitempty" name:"Account"`

	// Account name
	Nickname *string `json:"Nickname,omitnil,omitempty" name:"Nickname"`
}

type PrivateZone struct {
	// Private domain ID, which is in zone-xxxxxxxx format.
	ZoneId *string `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`

	// UIN of the domain name owner.
	OwnerUin *int64 `json:"OwnerUin,omitnil,omitempty" name:"OwnerUin"`

	// Private domain name.
	Domain *string `json:"Domain,omitnil,omitempty" name:"Domain"`

	// Creation time
	CreatedOn *string `json:"CreatedOn,omitnil,omitempty" name:"CreatedOn"`

	// Modification time
	UpdatedOn *string `json:"UpdatedOn,omitnil,omitempty" name:"UpdatedOn"`

	// Number of records.
	RecordCount *int64 `json:"RecordCount,omitnil,omitempty" name:"RecordCount"`

	// Remarks.
	// Note: This field may return null, indicating that no valid values can be obtained.
	Remark *string `json:"Remark,omitnil,omitempty" name:"Remark"`

	// List of bound VPCs.
	VpcSet []*VpcInfo `json:"VpcSet,omitnil,omitempty" name:"VpcSet"`

	// Status of the VPC bound with the private domain. SUSPEND: The VPC is not associated; ENABLED: the VPC has been associated.
	// , FAILED: the VPC fails to be associated.
	Status *string `json:"Status,omitnil,omitempty" name:"Status"`

	// Recursive resolution status of the domain name. ENABLED: enabled; DISABLED: disabled.
	DnsForwardStatus *string `json:"DnsForwardStatus,omitnil,omitempty" name:"DnsForwardStatus"`

	// Tag key-value pair collection.
	Tags []*TagInfo `json:"Tags,omitnil,omitempty" name:"Tags"`

	// List of bound VPCs of the associated account.
	// Note: This field may return null, indicating that no valid values can be obtained.
	AccountVpcSet []*AccountVpcInfoOutput `json:"AccountVpcSet,omitnil,omitempty" name:"AccountVpcSet"`

	// Whether the TLD is a custom one.
	// Note: This field may return null, indicating that no valid values can be obtained.
	IsCustomTld *bool `json:"IsCustomTld,omitnil,omitempty" name:"IsCustomTld"`

	// CNAME acceleration status. ENABLED: enabled; DISABLED: disabled.
	CnameSpeedupStatus *string `json:"CnameSpeedupStatus,omitnil,omitempty" name:"CnameSpeedupStatus"`

	// Forwarding rule name.
	// Note: This field may return null, indicating that no valid values can be obtained.
	ForwardRuleName *string `json:"ForwardRuleName,omitnil,omitempty" name:"ForwardRuleName"`

	// Forwarding rule type. DOWN: from cloud to off-cloud; UP: from off-cloud to cloud. Currently, only DOWN is supported.
	// Note: This field may return null, indicating that no valid values can be obtained.
	ForwardRuleType *string `json:"ForwardRuleType,omitnil,omitempty" name:"ForwardRuleType"`

	// Forwarding address.
	// Note: This field may return null, indicating that no valid values can be obtained.
	ForwardAddress *string `json:"ForwardAddress,omitnil,omitempty" name:"ForwardAddress"`

	// Endpoint name.
	// Note: This field may return null, indicating that no valid values can be obtained.
	EndPointName *string `json:"EndPointName,omitnil,omitempty" name:"EndPointName"`

	// Deleted VPC.
	// Note: This field may return null, indicating that no valid values can be obtained.
	DeletedVpcSet []*VpcInfo `json:"DeletedVpcSet,omitnil,omitempty" name:"DeletedVpcSet"`
}

type PrivateZoneRecord struct {
	// Record ID
	RecordId *string `json:"RecordId,omitnil,omitempty" name:"RecordId"`

	// Private domain ID: zone-xxxxxxxx
	ZoneId *string `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`

	// Subdomain
	SubDomain *string `json:"SubDomain,omitnil,omitempty" name:"SubDomain"`

	// Record type. Valid values: "A", "AAAA", "CNAME", "MX", "TXT", "PTR"
	RecordType *string `json:"RecordType,omitnil,omitempty" name:"RecordType"`

	// Record value
	RecordValue *string `json:"RecordValue,omitnil,omitempty" name:"RecordValue"`

	// Record cache time. The smaller the value, the faster the record will take effect. Value range: 1–86400s. Default value: 600
	TTL *int64 `json:"TTL,omitnil,omitempty" name:"TTL"`

	// MX priority, which is required when the record type is MX. Valid values: 5, 10, 15, 20, 30, 40, 50
	// Note: this field may return null, indicating that no valid values can be obtained.
	MX *int64 `json:"MX,omitnil,omitempty" name:"MX"`

	// Record status: ENABLED
	Status *string `json:"Status,omitnil,omitempty" name:"Status"`

	// Record weight. Value range: 1–100
	// Note: this field may return null, indicating that no valid values can be obtained.
	Weight *int64 `json:"Weight,omitnil,omitempty" name:"Weight"`

	// Record creation time
	CreatedOn *string `json:"CreatedOn,omitnil,omitempty" name:"CreatedOn"`

	// Record update time
	UpdatedOn *string `json:"UpdatedOn,omitnil,omitempty" name:"UpdatedOn"`

	// Additional information
	// Note: this field may return null, indicating that no valid values can be obtained.
	Extra *string `json:"Extra,omitnil,omitempty" name:"Extra"`
}

type RecordInfo struct {
	// Record ID.
	RecordId *string `json:"RecordId,omitnil,omitempty" name:"RecordId"`

	// Private domain ID, which is in zone-xxxxxxxx format.
	ZoneId *string `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`

	// Subdomain name.
	SubDomain *string `json:"SubDomain,omitnil,omitempty" name:"SubDomain"`

	// Record type. Valid values: A, AAAA, CNAME, MX, TXT, and PTR.
	RecordType *string `json:"RecordType,omitnil,omitempty" name:"RecordType"`

	// Record value.
	RecordValue *string `json:"RecordValue,omitnil,omitempty" name:"RecordValue"`

	// Record cache time. The smaller the value, the faster the record will take effect. Value range: 1-86,400s. Default value: 600.
	TTL *int64 `json:"TTL,omitnil,omitempty" name:"TTL"`

	// MX priority, which is required when the record type is mx. valid values: 5, 10, 15, 20, 30, 40, and 50.
	MX *int64 `json:"MX,omitnil,omitempty" name:"MX"`

	// Record weight. valid values: 1–100.
	Weight *int64 `json:"Weight,omitnil,omitempty" name:"Weight"`

	// Record creation time.
	CreatedOn *string `json:"CreatedOn,omitnil,omitempty" name:"CreatedOn"`

	// Record update time.
	UpdatedOn *string `json:"UpdatedOn,omitnil,omitempty" name:"UpdatedOn"`

	// 0 suspend 1 enable.
	Enabled *uint64 `json:"Enabled,omitnil,omitempty" name:"Enabled"`

	// Remarks
	Remark *string `json:"Remark,omitnil,omitempty" name:"Remark"`
}

type RegionInfo struct {
	// Region encoding
	RegionCode *string `json:"RegionCode,omitnil,omitempty" name:"RegionCode"`

	// Region name
	// 
	// Note: This field may return null, indicating that no valid values can be obtained.
	CnName *string `json:"CnName,omitnil,omitempty" name:"CnName"`

	// English name of the region
	EnName *string `json:"EnName,omitnil,omitempty" name:"EnName"`

	// Region ID
	// 
	// Note: This field may return null, indicating that no valid values can be obtained.
	RegionId *uint64 `json:"RegionId,omitnil,omitempty" name:"RegionId"`

	// Number of AZs
	// 
	// Note: This field may return null, indicating that no valid values can be obtained.
	AvailableZoneNum *uint64 `json:"AvailableZoneNum,omitnil,omitempty" name:"AvailableZoneNum"`
}

// Predefined struct for user
type SubscribePrivateZoneServiceRequestParams struct {

}

type SubscribePrivateZoneServiceRequest struct {
	*tchttp.BaseRequest
	
}

func (r *SubscribePrivateZoneServiceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *SubscribePrivateZoneServiceRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "SubscribePrivateZoneServiceRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type SubscribePrivateZoneServiceResponseParams struct {
	// Private DNS service activation status
	ServiceStatus *string `json:"ServiceStatus,omitnil,omitempty" name:"ServiceStatus"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type SubscribePrivateZoneServiceResponse struct {
	*tchttp.BaseResponse
	Response *SubscribePrivateZoneServiceResponseParams `json:"Response"`
}

func (r *SubscribePrivateZoneServiceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *SubscribePrivateZoneServiceResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type TagInfo struct {
	// Tag key
	TagKey *string `json:"TagKey,omitnil,omitempty" name:"TagKey"`

	// Tag value
	TagValue *string `json:"TagValue,omitnil,omitempty" name:"TagValue"`
}

type TldQuota struct {
	// Total quota
	Total *int64 `json:"Total,omitnil,omitempty" name:"Total"`

	// Used quota
	Used *int64 `json:"Used,omitnil,omitempty" name:"Used"`

	// Available quota
	Stock *int64 `json:"Stock,omitnil,omitempty" name:"Stock"`

	// User’s quota
	Quota *int64 `json:"Quota,omitnil,omitempty" name:"Quota"`
}

type VpcInfo struct {
	// VpcId: vpc-xadsafsdasd
	UniqVpcId *string `json:"UniqVpcId,omitnil,omitempty" name:"UniqVpcId"`

	// VPC region: ap-guangzhou, ap-shanghai
	Region *string `json:"Region,omitnil,omitempty" name:"Region"`
}