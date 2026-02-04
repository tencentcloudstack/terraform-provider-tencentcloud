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

package v20170312

import (
    tcerr "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common/errors"
    tchttp "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common/http"
    "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common/json"
)

type ActionTimer struct {
	// Timer action currently only supports terminating one value: TerminateInstances.
	TimerAction *string `json:"TimerAction,omitnil,omitempty" name:"TimerAction"`

	// Execution time, in standard ISO8601 representation and using UTC time. format: YYYY-MM-DDThh:MM:ssZ. for example, 2018-05-29T11:26:40Z. the execution time must be later than the current time by 5 minutes.
	ActionTime *string `json:"ActionTime,omitnil,omitempty" name:"ActionTime"`

	// Extension data. only used as output usage.
	Externals *Externals `json:"Externals,omitnil,omitempty" name:"Externals"`

	// Timer ID. only used as output usage.
	ActionTimerId *string `json:"ActionTimerId,omitnil,omitempty" name:"ActionTimerId"`

	// Timer status, for output usage only. value ranges from: <li>UNDO: unexecuted</li> <li>DOING: executing</li> <li>DONE: execution completed.</li>.
	Status *string `json:"Status,omitnil,omitempty" name:"Status"`

	// Instance ID corresponding to a timer. used only for output.
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`
}

// Predefined struct for user
type AllocateHostsRequestParams struct {
	// Instance location. This parameter is used to specify the attributes of an instance, such as its availability zone and project.
	Placement *Placement `json:"Placement,omitnil,omitempty" name:"Placement"`

	// A string used to ensure the idempotency of the request.
	ClientToken *string `json:"ClientToken,omitnil,omitempty" name:"ClientToken"`

	// Not supported. Configuration of prepaid instances. You can use the parameter to specify the attributes of prepaid instances, such as the subscription period and the auto-renewal plan. This parameter is required for prepaid instances.
	HostChargePrepaid *ChargePrepaid `json:"HostChargePrepaid,omitnil,omitempty" name:"HostChargePrepaid"`

	// Instance [billing type](https://intl.cloud.tencent.com/document/product/213/2180?from_cn_redirect=1). <br><li>`POSTPAID_BY_HOUR`: Hourly-based pay-as-you-go <br>
	HostChargeType *string `json:"HostChargeType,omitnil,omitempty" name:"HostChargeType"`

	// CDH instance model. Default value: `HS1`.
	HostType *string `json:"HostType,omitnil,omitempty" name:"HostType"`

	// Quantity of CDH instances purchased. Default value: 1.
	HostCount *uint64 `json:"HostCount,omitnil,omitempty" name:"HostCount"`

	// Tag description. You can specify the parameter to associate a tag with an instance.
	TagSpecification []*TagSpecification `json:"TagSpecification,omitnil,omitempty" name:"TagSpecification"`
}

type AllocateHostsRequest struct {
	*tchttp.BaseRequest
	
	// Instance location. This parameter is used to specify the attributes of an instance, such as its availability zone and project.
	Placement *Placement `json:"Placement,omitnil,omitempty" name:"Placement"`

	// A string used to ensure the idempotency of the request.
	ClientToken *string `json:"ClientToken,omitnil,omitempty" name:"ClientToken"`

	// Not supported. Configuration of prepaid instances. You can use the parameter to specify the attributes of prepaid instances, such as the subscription period and the auto-renewal plan. This parameter is required for prepaid instances.
	HostChargePrepaid *ChargePrepaid `json:"HostChargePrepaid,omitnil,omitempty" name:"HostChargePrepaid"`

	// Instance [billing type](https://intl.cloud.tencent.com/document/product/213/2180?from_cn_redirect=1). <br><li>`POSTPAID_BY_HOUR`: Hourly-based pay-as-you-go <br>
	HostChargeType *string `json:"HostChargeType,omitnil,omitempty" name:"HostChargeType"`

	// CDH instance model. Default value: `HS1`.
	HostType *string `json:"HostType,omitnil,omitempty" name:"HostType"`

	// Quantity of CDH instances purchased. Default value: 1.
	HostCount *uint64 `json:"HostCount,omitnil,omitempty" name:"HostCount"`

	// Tag description. You can specify the parameter to associate a tag with an instance.
	TagSpecification []*TagSpecification `json:"TagSpecification,omitnil,omitempty" name:"TagSpecification"`
}

func (r *AllocateHostsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *AllocateHostsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Placement")
	delete(f, "ClientToken")
	delete(f, "HostChargePrepaid")
	delete(f, "HostChargeType")
	delete(f, "HostType")
	delete(f, "HostCount")
	delete(f, "TagSpecification")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "AllocateHostsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type AllocateHostsResponseParams struct {
	// IDs of created instances
	HostIdSet []*string `json:"HostIdSet,omitnil,omitempty" name:"HostIdSet"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type AllocateHostsResponse struct {
	*tchttp.BaseResponse
	Response *AllocateHostsResponseParams `json:"Response"`
}

func (r *AllocateHostsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *AllocateHostsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type AssociateInstancesKeyPairsRequestParams struct {
	// Instance ID(s). The maximum number of instances in each request is 100. <br>You can obtain the available instance IDs in two ways: <br><li>Log in to the [console](https://console.cloud.tencent.com/cvm/index) to query the instance IDs. <br><li>Call [DescribeInstances](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1) and look for `InstanceId` in the response.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// Key ID(s). The maximum number of key pairs in each request is 100. The key pair ID is in the format of `skey-3glfot13`. <br>You can obtain the available key pair IDs in two ways: <br><li>Log in to the [console](https://console.cloud.tencent.com/cvm/sshkey) to query the key pair IDs. <br><li>Call [DescribeKeyPairs](https://intl.cloud.tencent.com/document/api/213/15699?from_cn_redirect=1) and look for `KeyId` in the response.
	KeyIds []*string `json:"KeyIds,omitnil,omitempty" name:"KeyIds"`

	// Whether to force shut down a running instances. It is recommended to manually shut down a running instance before associating a key pair with it. Valid values: <br><li>TRUE: force shut down an instance after a normal shutdown fails. <br><li>FALSE: do not force shut down an instance after a normal shutdown fails. <br><br>Default value: FALSE.
	ForceStop *bool `json:"ForceStop,omitnil,omitempty" name:"ForceStop"`
}

type AssociateInstancesKeyPairsRequest struct {
	*tchttp.BaseRequest
	
	// Instance ID(s). The maximum number of instances in each request is 100. <br>You can obtain the available instance IDs in two ways: <br><li>Log in to the [console](https://console.cloud.tencent.com/cvm/index) to query the instance IDs. <br><li>Call [DescribeInstances](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1) and look for `InstanceId` in the response.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// Key ID(s). The maximum number of key pairs in each request is 100. The key pair ID is in the format of `skey-3glfot13`. <br>You can obtain the available key pair IDs in two ways: <br><li>Log in to the [console](https://console.cloud.tencent.com/cvm/sshkey) to query the key pair IDs. <br><li>Call [DescribeKeyPairs](https://intl.cloud.tencent.com/document/api/213/15699?from_cn_redirect=1) and look for `KeyId` in the response.
	KeyIds []*string `json:"KeyIds,omitnil,omitempty" name:"KeyIds"`

	// Whether to force shut down a running instances. It is recommended to manually shut down a running instance before associating a key pair with it. Valid values: <br><li>TRUE: force shut down an instance after a normal shutdown fails. <br><li>FALSE: do not force shut down an instance after a normal shutdown fails. <br><br>Default value: FALSE.
	ForceStop *bool `json:"ForceStop,omitnil,omitempty" name:"ForceStop"`
}

func (r *AssociateInstancesKeyPairsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *AssociateInstancesKeyPairsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceIds")
	delete(f, "KeyIds")
	delete(f, "ForceStop")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "AssociateInstancesKeyPairsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type AssociateInstancesKeyPairsResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type AssociateInstancesKeyPairsResponse struct {
	*tchttp.BaseResponse
	Response *AssociateInstancesKeyPairsResponseParams `json:"Response"`
}

func (r *AssociateInstancesKeyPairsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *AssociateInstancesKeyPairsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type AssociateSecurityGroupsRequestParams struct {
	// ID of the security group to be associated, such as `sg-efil73jd`. Only one security group can be associated.
	SecurityGroupIds []*string `json:"SecurityGroupIds,omitnil,omitempty" name:"SecurityGroupIds"`

	// ID of the instance bound in the format of ins-lesecurk. You can specify up to 100 instances in each request.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`
}

type AssociateSecurityGroupsRequest struct {
	*tchttp.BaseRequest
	
	// ID of the security group to be associated, such as `sg-efil73jd`. Only one security group can be associated.
	SecurityGroupIds []*string `json:"SecurityGroupIds,omitnil,omitempty" name:"SecurityGroupIds"`

	// ID of the instance bound in the format of ins-lesecurk. You can specify up to 100 instances in each request.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`
}

func (r *AssociateSecurityGroupsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *AssociateSecurityGroupsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "SecurityGroupIds")
	delete(f, "InstanceIds")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "AssociateSecurityGroupsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type AssociateSecurityGroupsResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type AssociateSecurityGroupsResponse struct {
	*tchttp.BaseResponse
	Response *AssociateSecurityGroupsResponseParams `json:"Response"`
}

func (r *AssociateSecurityGroupsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *AssociateSecurityGroupsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type Attribute struct {
	// Custom data of instances.
	UserData *string `json:"UserData,omitnil,omitempty" name:"UserData"`
}

type ChargePrepaid struct {
	// Purchased usage period, in month. Valid values: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36
	Period *uint64 `json:"Period,omitnil,omitempty" name:"Period"`

	// Auto renewal flag. Valid values: <br><li>NOTIFY_AND_AUTO_RENEW: notify upon expiration and renew automatically <br><li>NOTIFY_AND_MANUAL_RENEW: notify upon expiration but do not renew automatically <br><li>DISABLE_NOTIFY_AND_MANUAL_RENEW: neither notify upon expiration nor renew automatically <br><br>Default value: NOTIFY_AND_AUTO_RENEW. If this parameter is specified as NOTIFY_AND_AUTO_RENEW, the instance will be automatically renewed on a monthly basis if the account balance is sufficient.
	RenewFlag *string `json:"RenewFlag,omitnil,omitempty" name:"RenewFlag"`
}

type ChcDeployExtraConfig struct {

	MiniOsType *string `json:"MiniOsType,omitnil,omitempty" name:"MiniOsType"`


	BootType *string `json:"BootType,omitnil,omitempty" name:"BootType"`


	BootFile *string `json:"BootFile,omitnil,omitempty" name:"BootFile"`


	NextServerAddress *string `json:"NextServerAddress,omitnil,omitempty" name:"NextServerAddress"`
}

type ChcHost struct {
	// CHC host ID
	ChcId *string `json:"ChcId,omitnil,omitempty" name:"ChcId"`

	// Instance name
	InstanceName *string `json:"InstanceName,omitnil,omitempty" name:"InstanceName"`

	// Server serial number
	SerialNumber *string `json:"SerialNumber,omitnil,omitempty" name:"SerialNumber"`

	// CHC host status<br/>
	// <ul>
	// <li>REGISTERED: The CHC host is registered, but the out-of-band network and deployment network are not configured.</li>
	// <li>VPC_READY: The out-of-band network and deployment network are configured.</li>
	// <li>PREPARED: It’s ready and can be associated with a CVM.</li>
	// <li>ONLINE: It’s already associated with a CVM.</li>
	// </ul>
	InstanceState *string `json:"InstanceState,omitnil,omitempty" name:"InstanceState"`

	// Device type
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	DeviceType *string `json:"DeviceType,omitnil,omitempty" name:"DeviceType"`

	// Availability zone
	Placement *Placement `json:"Placement,omitnil,omitempty" name:"Placement"`

	// Out-of-band network
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	BmcVirtualPrivateCloud *VirtualPrivateCloud `json:"BmcVirtualPrivateCloud,omitnil,omitempty" name:"BmcVirtualPrivateCloud"`

	// Out-of-band network IP
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	BmcIp *string `json:"BmcIp,omitnil,omitempty" name:"BmcIp"`

	// Out-of-band network security group ID
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	BmcSecurityGroupIds []*string `json:"BmcSecurityGroupIds,omitnil,omitempty" name:"BmcSecurityGroupIds"`

	// Deployment network
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	DeployVirtualPrivateCloud *VirtualPrivateCloud `json:"DeployVirtualPrivateCloud,omitnil,omitempty" name:"DeployVirtualPrivateCloud"`

	// Deployment network IP
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	DeployIp *string `json:"DeployIp,omitnil,omitempty" name:"DeployIp"`

	// Deployment network security group ID
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	DeploySecurityGroupIds []*string `json:"DeploySecurityGroupIds,omitnil,omitempty" name:"DeploySecurityGroupIds"`

	// ID of the associated CVM
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	CvmInstanceId *string `json:"CvmInstanceId,omitnil,omitempty" name:"CvmInstanceId"`

	// Server creation time
	CreatedTime *string `json:"CreatedTime,omitnil,omitempty" name:"CreatedTime"`

	// Instance hardware description, including CPU cores, memory capacity and disk capacity.
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	HardwareDescription *string `json:"HardwareDescription,omitnil,omitempty" name:"HardwareDescription"`

	// CPU cores of the CHC host
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	CPU *int64 `json:"CPU,omitnil,omitempty" name:"CPU"`

	// Memory capacity of the CHC host (unit: GB)
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	Memory *int64 `json:"Memory,omitnil,omitempty" name:"Memory"`

	// Disk capacity of the CHC host
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	Disk *string `json:"Disk,omitnil,omitempty" name:"Disk"`

	// MAC address assigned under the out-of-band network
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	BmcMAC *string `json:"BmcMAC,omitnil,omitempty" name:"BmcMAC"`

	// MAC address assigned under the deployment network
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	DeployMAC *string `json:"DeployMAC,omitnil,omitempty" name:"DeployMAC"`

	// Management type
	// HOSTING: Hosting
	// TENANT: Leasing
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	TenantType *string `json:"TenantType,omitnil,omitempty" name:"TenantType"`

	// CHC DHCP option, which is used for MiniOS debugging.
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	DeployExtraConfig *ChcDeployExtraConfig `json:"DeployExtraConfig,omitnil,omitempty" name:"DeployExtraConfig"`
}

type ChcHostDeniedActions struct {
	// CHC instance ID
	ChcId *string `json:"ChcId,omitnil,omitempty" name:"ChcId"`

	// CHC instance status
	State *string `json:"State,omitnil,omitempty" name:"State"`

	// Actions not allowed for the current CHC instance
	DenyActions []*string `json:"DenyActions,omitnil,omitempty" name:"DenyActions"`
}

// Predefined struct for user
type ConfigureChcAssistVpcRequestParams struct {
	// CHC host IDs
	ChcIds []*string `json:"ChcIds,omitnil,omitempty" name:"ChcIds"`

	// Out-of-band network information
	BmcVirtualPrivateCloud *VirtualPrivateCloud `json:"BmcVirtualPrivateCloud,omitnil,omitempty" name:"BmcVirtualPrivateCloud"`

	// Out-of-band network security group list
	BmcSecurityGroupIds []*string `json:"BmcSecurityGroupIds,omitnil,omitempty" name:"BmcSecurityGroupIds"`

	// Deployment network information
	DeployVirtualPrivateCloud *VirtualPrivateCloud `json:"DeployVirtualPrivateCloud,omitnil,omitempty" name:"DeployVirtualPrivateCloud"`

	// Deployment network security group list
	DeploySecurityGroupIds []*string `json:"DeploySecurityGroupIds,omitnil,omitempty" name:"DeploySecurityGroupIds"`


	ChcDeployExtraConfig *ChcDeployExtraConfig `json:"ChcDeployExtraConfig,omitnil,omitempty" name:"ChcDeployExtraConfig"`
}

type ConfigureChcAssistVpcRequest struct {
	*tchttp.BaseRequest
	
	// CHC host IDs
	ChcIds []*string `json:"ChcIds,omitnil,omitempty" name:"ChcIds"`

	// Out-of-band network information
	BmcVirtualPrivateCloud *VirtualPrivateCloud `json:"BmcVirtualPrivateCloud,omitnil,omitempty" name:"BmcVirtualPrivateCloud"`

	// Out-of-band network security group list
	BmcSecurityGroupIds []*string `json:"BmcSecurityGroupIds,omitnil,omitempty" name:"BmcSecurityGroupIds"`

	// Deployment network information
	DeployVirtualPrivateCloud *VirtualPrivateCloud `json:"DeployVirtualPrivateCloud,omitnil,omitempty" name:"DeployVirtualPrivateCloud"`

	// Deployment network security group list
	DeploySecurityGroupIds []*string `json:"DeploySecurityGroupIds,omitnil,omitempty" name:"DeploySecurityGroupIds"`

	ChcDeployExtraConfig *ChcDeployExtraConfig `json:"ChcDeployExtraConfig,omitnil,omitempty" name:"ChcDeployExtraConfig"`
}

func (r *ConfigureChcAssistVpcRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ConfigureChcAssistVpcRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ChcIds")
	delete(f, "BmcVirtualPrivateCloud")
	delete(f, "BmcSecurityGroupIds")
	delete(f, "DeployVirtualPrivateCloud")
	delete(f, "DeploySecurityGroupIds")
	delete(f, "ChcDeployExtraConfig")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ConfigureChcAssistVpcRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ConfigureChcAssistVpcResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ConfigureChcAssistVpcResponse struct {
	*tchttp.BaseResponse
	Response *ConfigureChcAssistVpcResponseParams `json:"Response"`
}

func (r *ConfigureChcAssistVpcResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ConfigureChcAssistVpcResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ConfigureChcDeployVpcRequestParams struct {
	// CHC instance ID
	ChcIds []*string `json:"ChcIds,omitnil,omitempty" name:"ChcIds"`

	// Deployment network information
	DeployVirtualPrivateCloud *VirtualPrivateCloud `json:"DeployVirtualPrivateCloud,omitnil,omitempty" name:"DeployVirtualPrivateCloud"`

	// Deployment network security group list
	DeploySecurityGroupIds []*string `json:"DeploySecurityGroupIds,omitnil,omitempty" name:"DeploySecurityGroupIds"`


	ChcDeployExtraConfig *ChcDeployExtraConfig `json:"ChcDeployExtraConfig,omitnil,omitempty" name:"ChcDeployExtraConfig"`
}

type ConfigureChcDeployVpcRequest struct {
	*tchttp.BaseRequest
	
	// CHC instance ID
	ChcIds []*string `json:"ChcIds,omitnil,omitempty" name:"ChcIds"`

	// Deployment network information
	DeployVirtualPrivateCloud *VirtualPrivateCloud `json:"DeployVirtualPrivateCloud,omitnil,omitempty" name:"DeployVirtualPrivateCloud"`

	// Deployment network security group list
	DeploySecurityGroupIds []*string `json:"DeploySecurityGroupIds,omitnil,omitempty" name:"DeploySecurityGroupIds"`

	ChcDeployExtraConfig *ChcDeployExtraConfig `json:"ChcDeployExtraConfig,omitnil,omitempty" name:"ChcDeployExtraConfig"`
}

func (r *ConfigureChcDeployVpcRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ConfigureChcDeployVpcRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ChcIds")
	delete(f, "DeployVirtualPrivateCloud")
	delete(f, "DeploySecurityGroupIds")
	delete(f, "ChcDeployExtraConfig")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ConfigureChcDeployVpcRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ConfigureChcDeployVpcResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ConfigureChcDeployVpcResponse struct {
	*tchttp.BaseResponse
	Response *ConfigureChcDeployVpcResponseParams `json:"Response"`
}

func (r *ConfigureChcDeployVpcResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ConfigureChcDeployVpcResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ConvertOperatingSystemsRequestParams struct {
	// ID of an instance to execute operating system switching.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// Whether it is a minimum scale switching.
	MinimalConversion *bool `json:"MinimalConversion,omitnil,omitempty" name:"MinimalConversion"`

	// Whether it is pre-check only.
	DryRun *bool `json:"DryRun,omitnil,omitempty" name:"DryRun"`

	// Target operating system type for switching. Only TencentOS is supported.
	TargetOSType *string `json:"TargetOSType,omitnil,omitempty" name:"TargetOSType"`
}

type ConvertOperatingSystemsRequest struct {
	*tchttp.BaseRequest
	
	// ID of an instance to execute operating system switching.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// Whether it is a minimum scale switching.
	MinimalConversion *bool `json:"MinimalConversion,omitnil,omitempty" name:"MinimalConversion"`

	// Whether it is pre-check only.
	DryRun *bool `json:"DryRun,omitnil,omitempty" name:"DryRun"`

	// Target operating system type for switching. Only TencentOS is supported.
	TargetOSType *string `json:"TargetOSType,omitnil,omitempty" name:"TargetOSType"`
}

func (r *ConvertOperatingSystemsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ConvertOperatingSystemsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceIds")
	delete(f, "MinimalConversion")
	delete(f, "DryRun")
	delete(f, "TargetOSType")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ConvertOperatingSystemsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ConvertOperatingSystemsResponseParams struct {
	// Information about the target operating system, which is returned only when the input parameter DryRun is true.
	// Note: This field may return null, indicating that no valid value is found.
	SupportTargetOSList []*TargetOS `json:"SupportTargetOSList,omitnil,omitempty" name:"SupportTargetOSList"`

	// Task ID for operating system switching.
	// Note: This field may return null, indicating that no valid value is found.
	TaskId *string `json:"TaskId,omitnil,omitempty" name:"TaskId"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ConvertOperatingSystemsResponse struct {
	*tchttp.BaseResponse
	Response *ConvertOperatingSystemsResponseParams `json:"Response"`
}

func (r *ConvertOperatingSystemsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ConvertOperatingSystemsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type CpuTopology struct {
	// Number of physical CPU cores to enable.
	CoreCount *int64 `json:"CoreCount,omitnil,omitempty" name:"CoreCount"`

	// Number of threads per core. This parameter determines whether to enable or disable hyper-threading.<br><li>1: Disable hyper-threading.</li><br><li>2: Enable hyper-threading.</li>
	// If not set, an instance uses the default hyper-threading policy. For how to enable or disable hyper-threading, refer to [Enabling and Disabling Hyper-Threading](https://intl.cloud.tencent.com/document/product/213/103798?from_cn_redirect=1).
	ThreadPerCore *int64 `json:"ThreadPerCore,omitnil,omitempty" name:"ThreadPerCore"`
}

// Predefined struct for user
type CreateDisasterRecoverGroupRequestParams struct {
	// Name of the spread placement group. The name must be 1-60 characters long and can contain both Chinese characters and English letters.
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// Type of a spread placement group. Valid values:<br><li>HOST: physical machine.</li><li>SW: switch.</li><li>RACK: rack.</li>
	Type *string `json:"Type,omitnil,omitempty" name:"Type"`

	// A string used to ensure the idempotency of the request, which is generated by the user and must be unique to each request. The maximum length is 64 ASCII characters. If this parameter is not specified, the idempotency of the request cannot be guaranteed. <br>For more information, see 'How to ensure idempotency'.
	ClientToken *string `json:"ClientToken,omitnil,omitempty" name:"ClientToken"`


	Affinity *int64 `json:"Affinity,omitnil,omitempty" name:"Affinity"`

	// List of tag description. By specifying this parameter, the tag can be bound to the placement group.
	TagSpecification []*TagSpecification `json:"TagSpecification,omitnil,omitempty" name:"TagSpecification"`
}

type CreateDisasterRecoverGroupRequest struct {
	*tchttp.BaseRequest
	
	// Name of the spread placement group. The name must be 1-60 characters long and can contain both Chinese characters and English letters.
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// Type of a spread placement group. Valid values:<br><li>HOST: physical machine.</li><li>SW: switch.</li><li>RACK: rack.</li>
	Type *string `json:"Type,omitnil,omitempty" name:"Type"`

	// A string used to ensure the idempotency of the request, which is generated by the user and must be unique to each request. The maximum length is 64 ASCII characters. If this parameter is not specified, the idempotency of the request cannot be guaranteed. <br>For more information, see 'How to ensure idempotency'.
	ClientToken *string `json:"ClientToken,omitnil,omitempty" name:"ClientToken"`

	Affinity *int64 `json:"Affinity,omitnil,omitempty" name:"Affinity"`

	// List of tag description. By specifying this parameter, the tag can be bound to the placement group.
	TagSpecification []*TagSpecification `json:"TagSpecification,omitnil,omitempty" name:"TagSpecification"`
}

func (r *CreateDisasterRecoverGroupRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateDisasterRecoverGroupRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Name")
	delete(f, "Type")
	delete(f, "ClientToken")
	delete(f, "Affinity")
	delete(f, "TagSpecification")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CreateDisasterRecoverGroupRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateDisasterRecoverGroupResponseParams struct {
	// List of spread placement group IDs.
	DisasterRecoverGroupId *string `json:"DisasterRecoverGroupId,omitnil,omitempty" name:"DisasterRecoverGroupId"`

	// Type of a spread placement group. Valid values:<br><li>HOST: physical machine.</li><li>SW: switch.</li><li>RACK: rack.</li>
	Type *string `json:"Type,omitnil,omitempty" name:"Type"`

	// Name of the spread placement group. The name must be 1-60 characters long and can contain both Chinese characters and English letters.
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// The maximum number of CVMs in a placement group.
	CvmQuotaTotal *int64 `json:"CvmQuotaTotal,omitnil,omitempty" name:"CvmQuotaTotal"`

	// The current number of CVMs in a placement group.
	CurrentNum *int64 `json:"CurrentNum,omitnil,omitempty" name:"CurrentNum"`

	// Creation time of the placement group.
	CreateTime *string `json:"CreateTime,omitnil,omitempty" name:"CreateTime"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type CreateDisasterRecoverGroupResponse struct {
	*tchttp.BaseResponse
	Response *CreateDisasterRecoverGroupResponseParams `json:"Response"`
}

func (r *CreateDisasterRecoverGroupResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateDisasterRecoverGroupResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateImageRequestParams struct {
	// Image name.
	ImageName *string `json:"ImageName,omitnil,omitempty" name:"ImageName"`

	// ID of the instance from which an image will be created. This parameter is required when using instance to create an image.
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// Image description.
	ImageDescription *string `json:"ImageDescription,omitnil,omitempty" name:"ImageDescription"`

	// Whether to perform forced power-off operation to create an image.
	// Valid values:<br><li>true: indicates that an image is created after forced power-off operation</li><br><li>false: indicates that an image is created in the power-on state</li><br><br>Default value: false.<br><br>Creating an image in the power-on state may result in some unbacked-up data, affecting data security.
	ForcePoweroff *string `json:"ForcePoweroff,omitnil,omitempty" name:"ForcePoweroff"`

	// Whether to enable Sysprep when creating a Windows image.
	// Valid values: `TRUE` and `FALSE`; default value: `FALSE`.
	// 
	// Click [here](https://intl.cloud.tencent.com/document/product/213/43498?from_cn_redirect=1) to learn more about Sysprep.
	Sysprep *string `json:"Sysprep,omitnil,omitempty" name:"Sysprep"`

	// IDs of data disks included in the image. 
	DataDiskIds []*string `json:"DataDiskIds,omitnil,omitempty" name:"DataDiskIds"`

	// Specified snapshot ID used to create an image. A system disk snapshot must be included. It cannot be passed together with `InstanceId`.
	SnapshotIds []*string `json:"SnapshotIds,omitnil,omitempty" name:"SnapshotIds"`

	// Success status of this request, without affecting the resources involved
	DryRun *bool `json:"DryRun,omitnil,omitempty" name:"DryRun"`

	// Tag description list. This parameter is used to bind a tag to a custom image.
	TagSpecification []*TagSpecification `json:"TagSpecification,omitnil,omitempty" name:"TagSpecification"`

	// Image family
	ImageFamily *string `json:"ImageFamily,omitnil,omitempty" name:"ImageFamily"`
}

type CreateImageRequest struct {
	*tchttp.BaseRequest
	
	// Image name.
	ImageName *string `json:"ImageName,omitnil,omitempty" name:"ImageName"`

	// ID of the instance from which an image will be created. This parameter is required when using instance to create an image.
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// Image description.
	ImageDescription *string `json:"ImageDescription,omitnil,omitempty" name:"ImageDescription"`

	// Whether to perform forced power-off operation to create an image.
	// Valid values:<br><li>true: indicates that an image is created after forced power-off operation</li><br><li>false: indicates that an image is created in the power-on state</li><br><br>Default value: false.<br><br>Creating an image in the power-on state may result in some unbacked-up data, affecting data security.
	ForcePoweroff *string `json:"ForcePoweroff,omitnil,omitempty" name:"ForcePoweroff"`

	// Whether to enable Sysprep when creating a Windows image.
	// Valid values: `TRUE` and `FALSE`; default value: `FALSE`.
	// 
	// Click [here](https://intl.cloud.tencent.com/document/product/213/43498?from_cn_redirect=1) to learn more about Sysprep.
	Sysprep *string `json:"Sysprep,omitnil,omitempty" name:"Sysprep"`

	// IDs of data disks included in the image. 
	DataDiskIds []*string `json:"DataDiskIds,omitnil,omitempty" name:"DataDiskIds"`

	// Specified snapshot ID used to create an image. A system disk snapshot must be included. It cannot be passed together with `InstanceId`.
	SnapshotIds []*string `json:"SnapshotIds,omitnil,omitempty" name:"SnapshotIds"`

	// Success status of this request, without affecting the resources involved
	DryRun *bool `json:"DryRun,omitnil,omitempty" name:"DryRun"`

	// Tag description list. This parameter is used to bind a tag to a custom image.
	TagSpecification []*TagSpecification `json:"TagSpecification,omitnil,omitempty" name:"TagSpecification"`

	// Image family
	ImageFamily *string `json:"ImageFamily,omitnil,omitempty" name:"ImageFamily"`
}

func (r *CreateImageRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateImageRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ImageName")
	delete(f, "InstanceId")
	delete(f, "ImageDescription")
	delete(f, "ForcePoweroff")
	delete(f, "Sysprep")
	delete(f, "DataDiskIds")
	delete(f, "SnapshotIds")
	delete(f, "DryRun")
	delete(f, "TagSpecification")
	delete(f, "ImageFamily")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CreateImageRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateImageResponseParams struct {
	// Image ID.
	// Note: This field may return null, indicating that no valid value was found.
	ImageId *string `json:"ImageId,omitnil,omitempty" name:"ImageId"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type CreateImageResponse struct {
	*tchttp.BaseResponse
	Response *CreateImageResponseParams `json:"Response"`
}

func (r *CreateImageResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateImageResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateKeyPairRequestParams struct {
	// Name of the key pair, which can contain numbers, letters, and underscores, with a maximum length of 25 characters.
	KeyName *string `json:"KeyName,omitnil,omitempty" name:"KeyName"`

	// ID of the project to which the created key pair belongs.
	// 
	// You can obtain a project ID in the following ways:
	// <li>Query the project ID through the project list.</li>
	// <li>Call the [DescribeProjects](https://intl.cloud.tencent.com/document/api/651/78725?from_cn_redirect=1) API and obtain the `projectId` from the return information.</li>
	ProjectId *int64 `json:"ProjectId,omitnil,omitempty" name:"ProjectId"`

	// Tag description list. This parameter is used to bind a tag to a key pair.
	TagSpecification []*TagSpecification `json:"TagSpecification,omitnil,omitempty" name:"TagSpecification"`
}

type CreateKeyPairRequest struct {
	*tchttp.BaseRequest
	
	// Name of the key pair, which can contain numbers, letters, and underscores, with a maximum length of 25 characters.
	KeyName *string `json:"KeyName,omitnil,omitempty" name:"KeyName"`

	// ID of the project to which the created key pair belongs.
	// 
	// You can obtain a project ID in the following ways:
	// <li>Query the project ID through the project list.</li>
	// <li>Call the [DescribeProjects](https://intl.cloud.tencent.com/document/api/651/78725?from_cn_redirect=1) API and obtain the `projectId` from the return information.</li>
	ProjectId *int64 `json:"ProjectId,omitnil,omitempty" name:"ProjectId"`

	// Tag description list. This parameter is used to bind a tag to a key pair.
	TagSpecification []*TagSpecification `json:"TagSpecification,omitnil,omitempty" name:"TagSpecification"`
}

func (r *CreateKeyPairRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateKeyPairRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "KeyName")
	delete(f, "ProjectId")
	delete(f, "TagSpecification")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CreateKeyPairRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateKeyPairResponseParams struct {
	// Key pair information.
	KeyPair *KeyPair `json:"KeyPair,omitnil,omitempty" name:"KeyPair"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type CreateKeyPairResponse struct {
	*tchttp.BaseResponse
	Response *CreateKeyPairResponseParams `json:"Response"`
}

func (r *CreateKeyPairResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateKeyPairResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateLaunchTemplateRequestParams struct {
	// Name of an instance launch template. It contains 2 to 128 English or Chinese characters.
	LaunchTemplateName *string `json:"LaunchTemplateName,omitnil,omitempty" name:"LaunchTemplateName"`

	// Location of the instance. You can specify attributes such as availability zone, project, and host (specified when creating a instance on the CDH) to which the instance belongs through this parameter.
	Placement *Placement `json:"Placement,omitnil,omitempty" name:"Placement"`

	// Specify an effective [mirror](https://intl.cloud.tencent.com/document/product/213/4940?from_cn_redirect=1) ID in the format of `img-xxx`. there are four image types: <li>PUBLIC image</li> <li>custom image</li> <li>shared image</li> <li>service market image</li>  you can obtain available mirror ids in the following ways: <li>the mirror ids of `PUBLIC image`, `custom image` and `shared image` can be queried by logging in to the [console](https://console.cloud.tencent.com/cvm/image?rid=1&imageType=PUBLIC_image); the mirror ids of `service market image` can be queried through the [cloud market](https://market.cloud.tencent.com/list).</li> <li>call the api [DescribeImages](https://intl.cloud.tencent.com/document/api/213/15715?from_cn_redirect=1), input InstanceType to obtain the list of images supported by the current model, and take the `ImageId` field from the return information.</li>.
	ImageId *string `json:"ImageId,omitnil,omitempty" name:"ImageId"`

	// Version description of an instance launch template. It contains 2 to 256 English or Chinese characters.
	LaunchTemplateVersionDescription *string `json:"LaunchTemplateVersionDescription,omitnil,omitempty" name:"LaunchTemplateVersionDescription"`

	// Instance model. Different instance models specify different resource specifications.
	// 
	// <br><li>For instances created with the payment modes PREPAID or POSTPAID_BY_HOUR, the specific values can be obtained BY calling the [DescribeInstanceTypeConfigs](https://intl.cloud.tencent.com/document/api/213/15749?from_cn_redirect=1) api to get the latest specification table or referring to the [instance specifications](https://intl.cloud.tencent.com/document/product/213/11518?from_cn_redirect=1) description. if this parameter is not specified, the system will dynamically assign a default model based on the current resource sales situation in a region.</li><li>for instances created with the payment mode CDHPAID, this parameter has the prefix "CDH_" and is generated based on the CPU and memory configuration. the specific format is: CDH_XCXG. for example, for creating a CDH instance with 1 CPU core and 1 gb memory, this parameter should be CDH_1C1G.</li>.
	InstanceType *string `json:"InstanceType,omitnil,omitempty" name:"InstanceType"`

	// Instance system disk configuration information. If this parameter is not specified, it will be assigned based on the system default values.
	SystemDisk *SystemDisk `json:"SystemDisk,omitnil,omitempty" name:"SystemDisk"`

	// Instance data disk configuration information. if not specified, no data disks are purchased by default. support specifying 21 data disks at the time of purchase, among which a maximum of 1 LOCAL_BASIC data disk or LOCAL_SSD data disk can be included, and a maximum of 20 CLOUD_BASIC data disks, CLOUD_PREMIUM data disks or CLOUD_SSD data disks can be included.
	DataDisks []*DataDisk `json:"DataDisks,omitnil,omitempty" name:"DataDisks"`

	// VPC-related information configuration. You can specify information such as VPC ID and subnet ID through this parameter. If this parameter is not specified, use the basic network by default. If a VPC IP is specified in this parameter, it indicates the primary network interface card IP of each instance. In addition, the number of the InstanceCount parameter should be consistent with the number of the VPC IP and should not exceed 20.
	VirtualPrivateCloud *VirtualPrivateCloud `json:"VirtualPrivateCloud,omitnil,omitempty" name:"VirtualPrivateCloud"`

	// Public bandwidth-related settings. If this parameter is not specified, the public bandwidth is 0 Mbps by default.
	InternetAccessible *InternetAccessible `json:"InternetAccessible,omitnil,omitempty" name:"InternetAccessible"`

	// Number of instances to purchase. value range for monthly subscription instances: [1, 300]. value range for pay-as-you-go instances: [1, 100]. default value: 1. the number of instances to purchase must not exceed the remaining user quota. for specific quota limitations, see [CVM instance purchase limitations](https://intl.cloud.tencent.com/document/product/213/2664?from_cn_redirect=1).
	InstanceCount *int64 `json:"InstanceCount,omitnil,omitempty" name:"InstanceCount"`

	// Instance display name. <li>if the instance display name is not specified, it will display by default as 'unnamed'.</li> <li>when purchasing multiple instances, if the pattern string `{R:x}` is specified, it indicates generating numbers `[x, x+n-1]`, where `n` represents the number of purchased instances. for example, `server_{R:3}` will result in instance display names as `server_3` when purchasing 1 instance; when purchasing 2 instances, the instance display names will be `server_3` and `server_4` respectively. it supports specifying multiple pattern strings `{R:x}`.</li> <li>when purchasing multiple instances, if no pattern string is specified, a suffix `1, 2...n` will be added to the instance display name, where `n` represents the number of purchased instances. for example, for `server_`, when purchasing 2 instances, the instance display names will be `server_1` and `server_2` respectively.</li> <li>it supports up to 128 characters (including pattern strings).</li>.
	InstanceName *string `json:"InstanceName,omitnil,omitempty" name:"InstanceName"`

	// Instance login settings. this parameter allows you to set the instance login method to key or maintain the original login settings of the image.
	LoginSettings *LoginSettings `json:"LoginSettings,omitnil,omitempty" name:"LoginSettings"`

	// Security group to which an instance belongs. this parameter can be obtained by calling the sgId field in the returned value of [DescribeSecurityGroups](https://intl.cloud.tencent.com/document/api/215/15808?from_cn_redirect=1). if not specified, the default security group is bound.
	SecurityGroupIds []*string `json:"SecurityGroupIds,omitnil,omitempty" name:"SecurityGroupIds"`

	// Enhanced services. You can specify whether to enable services such as Cloud Security and Cloud Monitor through this parameter. If this parameter is not specified, Cloud Monitor and Cloud Security are enabled for public images by default, but not enabled for custom images and marketplace images by default. Instead, they use services retained in the images.
	EnhancedService *EnhancedService `json:"EnhancedService,omitnil,omitempty" name:"EnhancedService"`

	// A string used to ensure the idempotence of the request. This string is generated by the customer and should be unique across different requests, with a maximum length of 64 ASCII characters. If this parameter is not specified, the idempotence of the request cannot be guaranteed.
	ClientToken *string `json:"ClientToken,omitnil,omitempty" name:"ClientToken"`

	// Specifies the HostName of the cloud virtual machine.<br><li>the dot (.) and hyphen (-) cannot be used at the beginning or end of the HostName, and cannot be used consecutively.</li><li>for Windows instances: the character length is 2 to 15. it consists of letters (case insensitive), digits, and hyphens (-). dots (.) are not supported, and it cannot be all digits.</li><li>for other types (such as Linux) instances: the character length is 2 to 60. multiple dots are allowed. each segment between dots can include letters (case insensitive), digits, and hyphens (-).</li>.
	HostName *string `json:"HostName,omitnil,omitempty" name:"HostName"`

	// Scheduled task. You can specify a scheduled task for the instance through this parameter. Currently, only scheduled termination is supported.
	ActionTimer *ActionTimer `json:"ActionTimer,omitnil,omitempty" name:"ActionTimer"`

	// Placement group ID. Only one can be specified.
	DisasterRecoverGroupIds []*string `json:"DisasterRecoverGroupIds,omitnil,omitempty" name:"DisasterRecoverGroupIds"`

	// List of tag descriptions. You can bind tags to corresponding resource instances at the same time by specifying this parameter. Currently, only binding tags to the CVM is supported.
	TagSpecification []*TagSpecification `json:"TagSpecification,omitnil,omitempty" name:"TagSpecification"`

	// Market-Related options of the instance, such as relevant parameters of the bidding instance. this parameter is required if the payment mode of the specified instance is spot payment.
	InstanceMarketOptions *InstanceMarketOptionsRequest `json:"InstanceMarketOptions,omitnil,omitempty" name:"InstanceMarketOptions"`

	// User data provided for an instance must be encoded in base64. valid values for maximum data size are up to 16 KB. for details on obtaining this parameter, see running commands at startup for Windows (https://intl.cloud.tencent.com/document/product/213/17526?from_cn_redirect=1) and Linux (https://intl.cloud.tencent.com/document/product/213/17525?from_cn_redirect=1).
	UserData *string `json:"UserData,omitnil,omitempty" name:"UserData"`

	// Whether it is a pre-check for this request only.
	// true: sends a check request without creating an instance. check items include whether required parameters are filled in, request format, service limits, and cvm inventory.
	// If the check fails, return the corresponding error code.
	// If the check passed, return RequestId.
	// false (default): sends a normal request. after passing the check, creates an instance directly.
	DryRun *bool `json:"DryRun,omitnil,omitempty" name:"DryRun"`

	// CAM role name. it can be obtained through the roleName in the return value from the API DescribeRoleList.
	CamRoleName *string `json:"CamRoleName,omitnil,omitempty" name:"CamRoleName"`

	// High-performance computing cluster ID. If the created instance is a high-performance computing instance, the cluster where the instance is placed should be specified. Otherwise, it cannot be specified.
	HpcClusterId *string `json:"HpcClusterId,omitnil,omitempty" name:"HpcClusterId"`

	// Instance [billing mode](https://intl.cloud.tencent.com/document/product/213/2180?from_cn_redirect=1).<br><li>PREPAID: prepaid, that is, monthly subscription.</li><li>POSTPAID_BY_HOUR: pay-as-you-go by hour.</li><li>CDHPAID: CDH instance (created based on CDH; the resources of the host are free of charge).</li><li>SPOTPAID: spot payment.</li>Default value: POSTPAID_BY_HOUR.
	InstanceChargeType *string `json:"InstanceChargeType,omitnil,omitempty" name:"InstanceChargeType"`

	// Prepaid mode, that is, annual and monthly subscription related parameter settings. Through this parameter, you can specify the purchase duration of annual and monthly subscription instances, whether to set auto-renewal, etc. If the specified instance's billing mode is the prepaid mode, this parameter must be passed.
	InstanceChargePrepaid *InstanceChargePrepaid `json:"InstanceChargePrepaid,omitnil,omitempty" name:"InstanceChargePrepaid"`

	// Instance destruction protection flag: indicates whether an instance is allowed to be deleted through an api. value ranges from: - **TRUE**: indicates that instance protection is enabled, deletion through apis is not allowed. - **FALSE**: indicates that instance protection is disabled, deletion through apis is allowed.  default value: FALSE.
	DisableApiTermination *bool `json:"DisableApiTermination,omitnil,omitempty" name:"DisableApiTermination"`


	EnableJumboFrame *bool `json:"EnableJumboFrame,omitnil,omitempty" name:"EnableJumboFrame"`

	// Description list of tags. by specifying this parameter, tags can be bound to the instance launch template.
	LaunchTemplateTagSpecification []*TagSpecification `json:"LaunchTemplateTagSpecification,omitnil,omitempty" name:"LaunchTemplateTagSpecification"`

	// Custom metadata. specifies that custom metadata key-value pairs can be added when creating a CVM.
	// Note: this field is in beta test.
	Metadata *Metadata `json:"Metadata,omitnil,omitempty" name:"Metadata"`

	// Specifies that only the Update and Replace parameters are allowed. this parameter is valid only when custom Metadata is used in the template and Metadata is also transmitted in RunInstances. defaults to Replace.
	// 
	// -Update: if template t contains this parameter with a value of Update and metadata=[k1:v1, k2:v2], then RunInstances (with metadata=[k2:v3]) + t creates a cvm using metadata=[k1:v1, k2:v3]. 
	// -Replace: if the template t contains this parameter with a value of Replace and metadata=[k1:v1, k2:v2], then when creating a cvm using RunInstances (with metadata=[k2:v3]) + t, the created cvm will use metadata=[k2:v3]. 
	// Note: this field is in beta test.
	TemplateDataModifyAction *string `json:"TemplateDataModifyAction,omitnil,omitempty" name:"TemplateDataModifyAction"`
}

type CreateLaunchTemplateRequest struct {
	*tchttp.BaseRequest
	
	// Name of an instance launch template. It contains 2 to 128 English or Chinese characters.
	LaunchTemplateName *string `json:"LaunchTemplateName,omitnil,omitempty" name:"LaunchTemplateName"`

	// Location of the instance. You can specify attributes such as availability zone, project, and host (specified when creating a instance on the CDH) to which the instance belongs through this parameter.
	Placement *Placement `json:"Placement,omitnil,omitempty" name:"Placement"`

	// Specify an effective [mirror](https://intl.cloud.tencent.com/document/product/213/4940?from_cn_redirect=1) ID in the format of `img-xxx`. there are four image types: <li>PUBLIC image</li> <li>custom image</li> <li>shared image</li> <li>service market image</li>  you can obtain available mirror ids in the following ways: <li>the mirror ids of `PUBLIC image`, `custom image` and `shared image` can be queried by logging in to the [console](https://console.cloud.tencent.com/cvm/image?rid=1&imageType=PUBLIC_image); the mirror ids of `service market image` can be queried through the [cloud market](https://market.cloud.tencent.com/list).</li> <li>call the api [DescribeImages](https://intl.cloud.tencent.com/document/api/213/15715?from_cn_redirect=1), input InstanceType to obtain the list of images supported by the current model, and take the `ImageId` field from the return information.</li>.
	ImageId *string `json:"ImageId,omitnil,omitempty" name:"ImageId"`

	// Version description of an instance launch template. It contains 2 to 256 English or Chinese characters.
	LaunchTemplateVersionDescription *string `json:"LaunchTemplateVersionDescription,omitnil,omitempty" name:"LaunchTemplateVersionDescription"`

	// Instance model. Different instance models specify different resource specifications.
	// 
	// <br><li>For instances created with the payment modes PREPAID or POSTPAID_BY_HOUR, the specific values can be obtained BY calling the [DescribeInstanceTypeConfigs](https://intl.cloud.tencent.com/document/api/213/15749?from_cn_redirect=1) api to get the latest specification table or referring to the [instance specifications](https://intl.cloud.tencent.com/document/product/213/11518?from_cn_redirect=1) description. if this parameter is not specified, the system will dynamically assign a default model based on the current resource sales situation in a region.</li><li>for instances created with the payment mode CDHPAID, this parameter has the prefix "CDH_" and is generated based on the CPU and memory configuration. the specific format is: CDH_XCXG. for example, for creating a CDH instance with 1 CPU core and 1 gb memory, this parameter should be CDH_1C1G.</li>.
	InstanceType *string `json:"InstanceType,omitnil,omitempty" name:"InstanceType"`

	// Instance system disk configuration information. If this parameter is not specified, it will be assigned based on the system default values.
	SystemDisk *SystemDisk `json:"SystemDisk,omitnil,omitempty" name:"SystemDisk"`

	// Instance data disk configuration information. if not specified, no data disks are purchased by default. support specifying 21 data disks at the time of purchase, among which a maximum of 1 LOCAL_BASIC data disk or LOCAL_SSD data disk can be included, and a maximum of 20 CLOUD_BASIC data disks, CLOUD_PREMIUM data disks or CLOUD_SSD data disks can be included.
	DataDisks []*DataDisk `json:"DataDisks,omitnil,omitempty" name:"DataDisks"`

	// VPC-related information configuration. You can specify information such as VPC ID and subnet ID through this parameter. If this parameter is not specified, use the basic network by default. If a VPC IP is specified in this parameter, it indicates the primary network interface card IP of each instance. In addition, the number of the InstanceCount parameter should be consistent with the number of the VPC IP and should not exceed 20.
	VirtualPrivateCloud *VirtualPrivateCloud `json:"VirtualPrivateCloud,omitnil,omitempty" name:"VirtualPrivateCloud"`

	// Public bandwidth-related settings. If this parameter is not specified, the public bandwidth is 0 Mbps by default.
	InternetAccessible *InternetAccessible `json:"InternetAccessible,omitnil,omitempty" name:"InternetAccessible"`

	// Number of instances to purchase. value range for monthly subscription instances: [1, 300]. value range for pay-as-you-go instances: [1, 100]. default value: 1. the number of instances to purchase must not exceed the remaining user quota. for specific quota limitations, see [CVM instance purchase limitations](https://intl.cloud.tencent.com/document/product/213/2664?from_cn_redirect=1).
	InstanceCount *int64 `json:"InstanceCount,omitnil,omitempty" name:"InstanceCount"`

	// Instance display name. <li>if the instance display name is not specified, it will display by default as 'unnamed'.</li> <li>when purchasing multiple instances, if the pattern string `{R:x}` is specified, it indicates generating numbers `[x, x+n-1]`, where `n` represents the number of purchased instances. for example, `server_{R:3}` will result in instance display names as `server_3` when purchasing 1 instance; when purchasing 2 instances, the instance display names will be `server_3` and `server_4` respectively. it supports specifying multiple pattern strings `{R:x}`.</li> <li>when purchasing multiple instances, if no pattern string is specified, a suffix `1, 2...n` will be added to the instance display name, where `n` represents the number of purchased instances. for example, for `server_`, when purchasing 2 instances, the instance display names will be `server_1` and `server_2` respectively.</li> <li>it supports up to 128 characters (including pattern strings).</li>.
	InstanceName *string `json:"InstanceName,omitnil,omitempty" name:"InstanceName"`

	// Instance login settings. this parameter allows you to set the instance login method to key or maintain the original login settings of the image.
	LoginSettings *LoginSettings `json:"LoginSettings,omitnil,omitempty" name:"LoginSettings"`

	// Security group to which an instance belongs. this parameter can be obtained by calling the sgId field in the returned value of [DescribeSecurityGroups](https://intl.cloud.tencent.com/document/api/215/15808?from_cn_redirect=1). if not specified, the default security group is bound.
	SecurityGroupIds []*string `json:"SecurityGroupIds,omitnil,omitempty" name:"SecurityGroupIds"`

	// Enhanced services. You can specify whether to enable services such as Cloud Security and Cloud Monitor through this parameter. If this parameter is not specified, Cloud Monitor and Cloud Security are enabled for public images by default, but not enabled for custom images and marketplace images by default. Instead, they use services retained in the images.
	EnhancedService *EnhancedService `json:"EnhancedService,omitnil,omitempty" name:"EnhancedService"`

	// A string used to ensure the idempotence of the request. This string is generated by the customer and should be unique across different requests, with a maximum length of 64 ASCII characters. If this parameter is not specified, the idempotence of the request cannot be guaranteed.
	ClientToken *string `json:"ClientToken,omitnil,omitempty" name:"ClientToken"`

	// Specifies the HostName of the cloud virtual machine.<br><li>the dot (.) and hyphen (-) cannot be used at the beginning or end of the HostName, and cannot be used consecutively.</li><li>for Windows instances: the character length is 2 to 15. it consists of letters (case insensitive), digits, and hyphens (-). dots (.) are not supported, and it cannot be all digits.</li><li>for other types (such as Linux) instances: the character length is 2 to 60. multiple dots are allowed. each segment between dots can include letters (case insensitive), digits, and hyphens (-).</li>.
	HostName *string `json:"HostName,omitnil,omitempty" name:"HostName"`

	// Scheduled task. You can specify a scheduled task for the instance through this parameter. Currently, only scheduled termination is supported.
	ActionTimer *ActionTimer `json:"ActionTimer,omitnil,omitempty" name:"ActionTimer"`

	// Placement group ID. Only one can be specified.
	DisasterRecoverGroupIds []*string `json:"DisasterRecoverGroupIds,omitnil,omitempty" name:"DisasterRecoverGroupIds"`

	// List of tag descriptions. You can bind tags to corresponding resource instances at the same time by specifying this parameter. Currently, only binding tags to the CVM is supported.
	TagSpecification []*TagSpecification `json:"TagSpecification,omitnil,omitempty" name:"TagSpecification"`

	// Market-Related options of the instance, such as relevant parameters of the bidding instance. this parameter is required if the payment mode of the specified instance is spot payment.
	InstanceMarketOptions *InstanceMarketOptionsRequest `json:"InstanceMarketOptions,omitnil,omitempty" name:"InstanceMarketOptions"`

	// User data provided for an instance must be encoded in base64. valid values for maximum data size are up to 16 KB. for details on obtaining this parameter, see running commands at startup for Windows (https://intl.cloud.tencent.com/document/product/213/17526?from_cn_redirect=1) and Linux (https://intl.cloud.tencent.com/document/product/213/17525?from_cn_redirect=1).
	UserData *string `json:"UserData,omitnil,omitempty" name:"UserData"`

	// Whether it is a pre-check for this request only.
	// true: sends a check request without creating an instance. check items include whether required parameters are filled in, request format, service limits, and cvm inventory.
	// If the check fails, return the corresponding error code.
	// If the check passed, return RequestId.
	// false (default): sends a normal request. after passing the check, creates an instance directly.
	DryRun *bool `json:"DryRun,omitnil,omitempty" name:"DryRun"`

	// CAM role name. it can be obtained through the roleName in the return value from the API DescribeRoleList.
	CamRoleName *string `json:"CamRoleName,omitnil,omitempty" name:"CamRoleName"`

	// High-performance computing cluster ID. If the created instance is a high-performance computing instance, the cluster where the instance is placed should be specified. Otherwise, it cannot be specified.
	HpcClusterId *string `json:"HpcClusterId,omitnil,omitempty" name:"HpcClusterId"`

	// Instance [billing mode](https://intl.cloud.tencent.com/document/product/213/2180?from_cn_redirect=1).<br><li>PREPAID: prepaid, that is, monthly subscription.</li><li>POSTPAID_BY_HOUR: pay-as-you-go by hour.</li><li>CDHPAID: CDH instance (created based on CDH; the resources of the host are free of charge).</li><li>SPOTPAID: spot payment.</li>Default value: POSTPAID_BY_HOUR.
	InstanceChargeType *string `json:"InstanceChargeType,omitnil,omitempty" name:"InstanceChargeType"`

	// Prepaid mode, that is, annual and monthly subscription related parameter settings. Through this parameter, you can specify the purchase duration of annual and monthly subscription instances, whether to set auto-renewal, etc. If the specified instance's billing mode is the prepaid mode, this parameter must be passed.
	InstanceChargePrepaid *InstanceChargePrepaid `json:"InstanceChargePrepaid,omitnil,omitempty" name:"InstanceChargePrepaid"`

	// Instance destruction protection flag: indicates whether an instance is allowed to be deleted through an api. value ranges from: - **TRUE**: indicates that instance protection is enabled, deletion through apis is not allowed. - **FALSE**: indicates that instance protection is disabled, deletion through apis is allowed.  default value: FALSE.
	DisableApiTermination *bool `json:"DisableApiTermination,omitnil,omitempty" name:"DisableApiTermination"`

	EnableJumboFrame *bool `json:"EnableJumboFrame,omitnil,omitempty" name:"EnableJumboFrame"`

	// Description list of tags. by specifying this parameter, tags can be bound to the instance launch template.
	LaunchTemplateTagSpecification []*TagSpecification `json:"LaunchTemplateTagSpecification,omitnil,omitempty" name:"LaunchTemplateTagSpecification"`

	// Custom metadata. specifies that custom metadata key-value pairs can be added when creating a CVM.
	// Note: this field is in beta test.
	Metadata *Metadata `json:"Metadata,omitnil,omitempty" name:"Metadata"`

	// Specifies that only the Update and Replace parameters are allowed. this parameter is valid only when custom Metadata is used in the template and Metadata is also transmitted in RunInstances. defaults to Replace.
	// 
	// -Update: if template t contains this parameter with a value of Update and metadata=[k1:v1, k2:v2], then RunInstances (with metadata=[k2:v3]) + t creates a cvm using metadata=[k1:v1, k2:v3]. 
	// -Replace: if the template t contains this parameter with a value of Replace and metadata=[k1:v1, k2:v2], then when creating a cvm using RunInstances (with metadata=[k2:v3]) + t, the created cvm will use metadata=[k2:v3]. 
	// Note: this field is in beta test.
	TemplateDataModifyAction *string `json:"TemplateDataModifyAction,omitnil,omitempty" name:"TemplateDataModifyAction"`
}

func (r *CreateLaunchTemplateRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateLaunchTemplateRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "LaunchTemplateName")
	delete(f, "Placement")
	delete(f, "ImageId")
	delete(f, "LaunchTemplateVersionDescription")
	delete(f, "InstanceType")
	delete(f, "SystemDisk")
	delete(f, "DataDisks")
	delete(f, "VirtualPrivateCloud")
	delete(f, "InternetAccessible")
	delete(f, "InstanceCount")
	delete(f, "InstanceName")
	delete(f, "LoginSettings")
	delete(f, "SecurityGroupIds")
	delete(f, "EnhancedService")
	delete(f, "ClientToken")
	delete(f, "HostName")
	delete(f, "ActionTimer")
	delete(f, "DisasterRecoverGroupIds")
	delete(f, "TagSpecification")
	delete(f, "InstanceMarketOptions")
	delete(f, "UserData")
	delete(f, "DryRun")
	delete(f, "CamRoleName")
	delete(f, "HpcClusterId")
	delete(f, "InstanceChargeType")
	delete(f, "InstanceChargePrepaid")
	delete(f, "DisableApiTermination")
	delete(f, "EnableJumboFrame")
	delete(f, "LaunchTemplateTagSpecification")
	delete(f, "Metadata")
	delete(f, "TemplateDataModifyAction")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CreateLaunchTemplateRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateLaunchTemplateResponseParams struct {
	// Specifies the ID of the successfully created instance launch template when this parameter is returned by creating an instance launch template through this interface.
	LaunchTemplateId *string `json:"LaunchTemplateId,omitnil,omitempty" name:"LaunchTemplateId"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type CreateLaunchTemplateResponse struct {
	*tchttp.BaseResponse
	Response *CreateLaunchTemplateResponseParams `json:"Response"`
}

func (r *CreateLaunchTemplateResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateLaunchTemplateResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateLaunchTemplateVersionRequestParams struct {
	// Location of the instance. You can use this parameter to specify the attributes of the instance, such as its availability zone, project, and CDH (for dedicated CVMs)
	Placement *Placement `json:"Placement,omitnil,omitempty" name:"Placement"`

	// Instance launch template ID. This parameter is used as a basis for creating new template versions.
	LaunchTemplateId *string `json:"LaunchTemplateId,omitnil,omitempty" name:"LaunchTemplateId"`

	// This parameter, when specified, is used to create instance launch templates. If this parameter is not specified, the default version will be used.
	LaunchTemplateVersion *int64 `json:"LaunchTemplateVersion,omitnil,omitempty" name:"LaunchTemplateVersion"`

	// Description of instance launch template versions. This parameter can contain 2-256 characters.
	LaunchTemplateVersionDescription *string `json:"LaunchTemplateVersionDescription,omitnil,omitempty" name:"LaunchTemplateVersionDescription"`

	// Instance model. Different instance models specify different resource specifications.
	// 
	// <br><li>For instances created with the payment modes PREPAID or POSTPAID_BY_HOUR, the specific values can be obtained by calling the [DescribeInstanceTypeConfigs](https://intl.cloud.tencent.com/document/api/213/15749?from_cn_redirect=1) API for the latest specification table or referring to [Instance Specifications](https://intl.cloud.tencent.com/document/product/213/11518?from_cn_redirect=1). If this parameter is not specified, the system will dynamically assign a default model based on the current resource sales situation in a region.</li><br><li>For instances created with the payment mode CDHPAID, this parameter has the prefix "CDH_" and is generated based on the CPU and memory configuration. The specific format is: CDH_XCXG. For example, for creating a CDH instance with 1 CPU core and 1 GB memory, this parameter should be CDH_1C1G.</li>
	InstanceType *string `json:"InstanceType,omitnil,omitempty" name:"InstanceType"`

	// Assign an effective [mirror](https://intl.cloud.tencent.com/document/product/213/4940?from_cn_redirect=1) ID in the format like `img-xxx`. there are four image types: <br/><li>PUBLIC image</li><li>custom image</li><li>shared image</li><li>cloud image market</li><br/>you can obtain available mirror ids in the following ways: <br/><li>the mirror ids of `PUBLIC image`, `custom image` and `shared image` can be queried by logging in to the [console](https://console.cloud.tencent.com/cvm/image?rid=1&imageType=PUBLIC_image); the mirror ID of `cloud image market` can be queried through the [cloud market](https://market.cloud.tencent.com/list).</li><li>call the api [DescribeImages](https://intl.cloud.tencent.com/document/api/213/15715?from_cn_redirect=1), input InstanceType to obtain the list of images supported by the current model, and take the `ImageId` field from the return information.</li>.
	ImageId *string `json:"ImageId,omitnil,omitempty" name:"ImageId"`

	// System disk configuration of the instance. If this parameter is not specified, the default value will be used.
	SystemDisk *SystemDisk `json:"SystemDisk,omitnil,omitempty" name:"SystemDisk"`

	// The configuration information of instance data disks. If this parameter is not specified, no data disk will be purchased by default. When purchasing, you can specify 21 data disks, which can contain at most 1 LOCAL_BASIC or LOCAL_SSD data disk, and at most 20 CLOUD_BASIC, CLOUD_PREMIUM, or CLOUD_SSD data disks.
	DataDisks []*DataDisk `json:"DataDisks,omitnil,omitempty" name:"DataDisks"`

	// Configuration information of VPC. This parameter is used to specify VPC ID and subnet ID, etc. If this parameter is not specified, the classic network is used by default. If a VPC IP is specified in this parameter, it indicates the primary ENI IP of each instance. The value of parameter InstanceCount must be same as the number of VPC IPs, which cannot be greater than 20.
	VirtualPrivateCloud *VirtualPrivateCloud `json:"VirtualPrivateCloud,omitnil,omitempty" name:"VirtualPrivateCloud"`

	// Configuration of public network bandwidth. If this parameter is not specified, 0 Mbps will be used by default.
	InternetAccessible *InternetAccessible `json:"InternetAccessible,omitnil,omitempty" name:"InternetAccessible"`

	// Number of instances to be purchased. Value range for monthly-subscribed instances: [1, 300]. Value range for pay-as-you-go instances: [1, 100]. Default value: 1. The specified number of instances to be purchased cannot exceed the remaining quota allowed for the user. For more information on quota, see CVM instance [Purchase Limits](https://intl.cloud.tencent.com/document/product/213/2664).
	InstanceCount *int64 `json:"InstanceCount,omitnil,omitempty" name:"InstanceCount"`

	// Instance display name. <li>if the instance display name is not specified, it will display by default as 'unnamed'.</li> <li>when purchasing multiple instances, if the pattern string `{R:x}` is specified, it indicates generating numbers `[x, x+n-1]`, where `n` represents the number of purchased instances. for example, `server_{R:3}` will result in instance display names as `server_3` when purchasing 1 instance; when purchasing 2 instances, the instance display names will be `server_3` and `server_4` respectively. it supports specifying multiple pattern strings `{R:x}`.</li> <li>when purchasing multiple instances, if no pattern string is specified, a suffix `1, 2...n` will be added to the instance display name, where `n` represents the number of purchased instances. for example, for `server_`, when purchasing 2 instances, the instance display names will be `server_1` and `server_2` respectively.</li> <li>it supports up to 128 characters (including pattern strings).</li>.
	InstanceName *string `json:"InstanceName,omitnil,omitempty" name:"InstanceName"`

	// Instance login settings. you can use this parameter to set the instance's login method to key or keep the original login settings of the image.
	LoginSettings *LoginSettings `json:"LoginSettings,omitnil,omitempty" name:"LoginSettings"`

	// Security groups to which the instance belongs. To obtain the security group IDs, you can call [DescribeSecurityGroups](https://intl.cloud.tencent.com/document/api/215/15808) and look for the `sgld` fields in the response. If this parameter is not specified, the instance will be associated with default security groups.
	SecurityGroupIds []*string `json:"SecurityGroupIds,omitnil,omitempty" name:"SecurityGroupIds"`

	// Enhanced service. this parameter can be used to specify whether to enable services such as cloud security and cloud monitoring. if this parameter is not specified, cloud monitor and cloud security services are enabled for public images by default; for custom images and marketplace images, cloud monitor and cloud security services are not enabled by default, and the services retained in the image will be used.
	EnhancedService *EnhancedService `json:"EnhancedService,omitnil,omitempty" name:"EnhancedService"`

	// A unique string supplied by the client to ensure that the request is idempotent. Its maximum length is 64 ASCII characters. If this parameter is not specified, the idempotency of the request cannot be guaranteed.
	ClientToken *string `json:"ClientToken,omitnil,omitempty" name:"ClientToken"`

	// Hostname of Cloud Virtual Machine.<br><li>Period (.) and hyphen (-) should not be used as the first or last character of the hostname, and should not be used consecutively.</li><br><li>Windows instances: The hostname should contain 2 to 15 characters, including letters (case insensitive), digits, and hyphens (-), does not support periods (.), and should not be all digits.</li><br><li>Instances of other types (such as Linux instances): The hostname should contain 2 to 60 characters, including multiple periods (.), with each segment between periods considered as one section. Each section can contain letters (case insensitive), digits, and hyphens (-).</li>
	HostName *string `json:"HostName,omitnil,omitempty" name:"HostName"`

	// Scheduled tasks. You can use this parameter to specify scheduled tasks for the instance. Only scheduled termination is supported.
	ActionTimer *ActionTimer `json:"ActionTimer,omitnil,omitempty" name:"ActionTimer"`

	// Placement group ID. You can only specify one.
	DisasterRecoverGroupIds []*string `json:"DisasterRecoverGroupIds,omitnil,omitempty" name:"DisasterRecoverGroupIds"`

	// The tag description list. This parameter is used to bind a tag to a resource instance. A tag can only be bound to CVM instances.
	TagSpecification []*TagSpecification `json:"TagSpecification,omitnil,omitempty" name:"TagSpecification"`

	// Market options of the instance, such as parameters related to spot instances. This parameter is required for spot instances.
	InstanceMarketOptions *InstanceMarketOptionsRequest `json:"InstanceMarketOptions,omitnil,omitempty" name:"InstanceMarketOptions"`

	// User data provided to the instance. This parameter needs to be encoded in base64 format with the maximum size of 16KB. For more information on how to get the value of this parameter, see the commands you need to execute on startup for [Windows](https://intl.cloud.tencent.com/document/product/213/17526) or [Linux](https://intl.cloud.tencent.com/document/product/213/17525).
	UserData *string `json:"UserData,omitnil,omitempty" name:"UserData"`

	// Whether the request is a dry run only.
	// `true`: dry run only. The request will not create instance(s). A dry run can check whether all the required parameters are specified, whether the request format is right, whether the request exceeds service limits, and whether the specified CVMs are available.
	// If the dry run fails, the corresponding error code will be returned.
	// If the dry run succeeds, the RequestId will be returned.
	// `false` (default value): send a normal request and create instance(s) if all the requirements are met.
	DryRun *bool `json:"DryRun,omitnil,omitempty" name:"DryRun"`

	// CAM role name, which can be obtained from the `roleName` field in the response of the [`DescribeRoleList`](https://intl.cloud.tencent.com/document/product/598/36223?from_cn_redirect=1) API.
	CamRoleName *string `json:"CamRoleName,omitnil,omitempty" name:"CamRoleName"`

	// HPC cluster ID. The HPC cluster must and can only be specified for a high-performance computing instance.
	HpcClusterId *string `json:"HpcClusterId,omitnil,omitempty" name:"HpcClusterId"`

	// Instance [billing mode](https://intl.cloud.tencent.com/document/product/213/2180?from_cn_redirect=1).<br><li>PREPAID: prepaid, that is, monthly subscription.</li><li>POSTPAID_BY_HOUR: pay-as-you-go by hour.</li><li>CDHPAID: CDH instance (created based on CDH; the resources of the host are free of charge).</li><li>SPOTPAID: spot payment.</li>Default value: POSTPAID_BY_HOUR.
	InstanceChargeType *string `json:"InstanceChargeType,omitnil,omitempty" name:"InstanceChargeType"`

	// Details of the monthly subscription, including the purchase period, auto-renewal. It is required if the `InstanceChargeType` is `PREPAID`.
	InstanceChargePrepaid *InstanceChargePrepaid `json:"InstanceChargePrepaid,omitnil,omitempty" name:"InstanceChargePrepaid"`

	// Instance termination protection flag, indicating whether an instance is allowed to be deleted through an API. Valid values:<br><li>TRUE: Instance protection is enabled, and the instance is not allowed to be deleted through the API.</li><br><li>FALSE: Instance protection is disabled, and the instance is allowed to be deleted through the API.</li><br><br>Default value: FALSE.
	DisableApiTermination *bool `json:"DisableApiTermination,omitnil,omitempty" name:"DisableApiTermination"`


	EnableJumboFrame *bool `json:"EnableJumboFrame,omitnil,omitempty" name:"EnableJumboFrame"`

	// Custom metadata. specifies that custom metadata key-value pairs can be added when creating a CVM.
	// Note: this field is in beta test.
	Metadata *Metadata `json:"Metadata,omitnil,omitempty" name:"Metadata"`

	// Specifies that only the Update and Replace parameters are allowed. this parameter is valid only when custom Metadata is used in the template and Metadata is also transmitted in RunInstances. defaults to Replace.
	// 
	// -Update: if template t contains this parameter with a value of Update and metadata=[k1:v1, k2:v2], then RunInstances (with metadata=[k2:v3]) + t creates a cvm using metadata=[k1:v1, k2:v3]. 
	// -Replace: if the template t contains this parameter with a value of Replace and metadata=[k1:v1, k2:v2], then when creating a cvm using RunInstances (with metadata=[k2:v3]) + t, the created cvm will use metadata=[k2:v3]. 
	// Note: this field is in beta test.
	TemplateDataModifyAction *string `json:"TemplateDataModifyAction,omitnil,omitempty" name:"TemplateDataModifyAction"`
}

type CreateLaunchTemplateVersionRequest struct {
	*tchttp.BaseRequest
	
	// Location of the instance. You can use this parameter to specify the attributes of the instance, such as its availability zone, project, and CDH (for dedicated CVMs)
	Placement *Placement `json:"Placement,omitnil,omitempty" name:"Placement"`

	// Instance launch template ID. This parameter is used as a basis for creating new template versions.
	LaunchTemplateId *string `json:"LaunchTemplateId,omitnil,omitempty" name:"LaunchTemplateId"`

	// This parameter, when specified, is used to create instance launch templates. If this parameter is not specified, the default version will be used.
	LaunchTemplateVersion *int64 `json:"LaunchTemplateVersion,omitnil,omitempty" name:"LaunchTemplateVersion"`

	// Description of instance launch template versions. This parameter can contain 2-256 characters.
	LaunchTemplateVersionDescription *string `json:"LaunchTemplateVersionDescription,omitnil,omitempty" name:"LaunchTemplateVersionDescription"`

	// Instance model. Different instance models specify different resource specifications.
	// 
	// <br><li>For instances created with the payment modes PREPAID or POSTPAID_BY_HOUR, the specific values can be obtained by calling the [DescribeInstanceTypeConfigs](https://intl.cloud.tencent.com/document/api/213/15749?from_cn_redirect=1) API for the latest specification table or referring to [Instance Specifications](https://intl.cloud.tencent.com/document/product/213/11518?from_cn_redirect=1). If this parameter is not specified, the system will dynamically assign a default model based on the current resource sales situation in a region.</li><br><li>For instances created with the payment mode CDHPAID, this parameter has the prefix "CDH_" and is generated based on the CPU and memory configuration. The specific format is: CDH_XCXG. For example, for creating a CDH instance with 1 CPU core and 1 GB memory, this parameter should be CDH_1C1G.</li>
	InstanceType *string `json:"InstanceType,omitnil,omitempty" name:"InstanceType"`

	// Assign an effective [mirror](https://intl.cloud.tencent.com/document/product/213/4940?from_cn_redirect=1) ID in the format like `img-xxx`. there are four image types: <br/><li>PUBLIC image</li><li>custom image</li><li>shared image</li><li>cloud image market</li><br/>you can obtain available mirror ids in the following ways: <br/><li>the mirror ids of `PUBLIC image`, `custom image` and `shared image` can be queried by logging in to the [console](https://console.cloud.tencent.com/cvm/image?rid=1&imageType=PUBLIC_image); the mirror ID of `cloud image market` can be queried through the [cloud market](https://market.cloud.tencent.com/list).</li><li>call the api [DescribeImages](https://intl.cloud.tencent.com/document/api/213/15715?from_cn_redirect=1), input InstanceType to obtain the list of images supported by the current model, and take the `ImageId` field from the return information.</li>.
	ImageId *string `json:"ImageId,omitnil,omitempty" name:"ImageId"`

	// System disk configuration of the instance. If this parameter is not specified, the default value will be used.
	SystemDisk *SystemDisk `json:"SystemDisk,omitnil,omitempty" name:"SystemDisk"`

	// The configuration information of instance data disks. If this parameter is not specified, no data disk will be purchased by default. When purchasing, you can specify 21 data disks, which can contain at most 1 LOCAL_BASIC or LOCAL_SSD data disk, and at most 20 CLOUD_BASIC, CLOUD_PREMIUM, or CLOUD_SSD data disks.
	DataDisks []*DataDisk `json:"DataDisks,omitnil,omitempty" name:"DataDisks"`

	// Configuration information of VPC. This parameter is used to specify VPC ID and subnet ID, etc. If this parameter is not specified, the classic network is used by default. If a VPC IP is specified in this parameter, it indicates the primary ENI IP of each instance. The value of parameter InstanceCount must be same as the number of VPC IPs, which cannot be greater than 20.
	VirtualPrivateCloud *VirtualPrivateCloud `json:"VirtualPrivateCloud,omitnil,omitempty" name:"VirtualPrivateCloud"`

	// Configuration of public network bandwidth. If this parameter is not specified, 0 Mbps will be used by default.
	InternetAccessible *InternetAccessible `json:"InternetAccessible,omitnil,omitempty" name:"InternetAccessible"`

	// Number of instances to be purchased. Value range for monthly-subscribed instances: [1, 300]. Value range for pay-as-you-go instances: [1, 100]. Default value: 1. The specified number of instances to be purchased cannot exceed the remaining quota allowed for the user. For more information on quota, see CVM instance [Purchase Limits](https://intl.cloud.tencent.com/document/product/213/2664).
	InstanceCount *int64 `json:"InstanceCount,omitnil,omitempty" name:"InstanceCount"`

	// Instance display name. <li>if the instance display name is not specified, it will display by default as 'unnamed'.</li> <li>when purchasing multiple instances, if the pattern string `{R:x}` is specified, it indicates generating numbers `[x, x+n-1]`, where `n` represents the number of purchased instances. for example, `server_{R:3}` will result in instance display names as `server_3` when purchasing 1 instance; when purchasing 2 instances, the instance display names will be `server_3` and `server_4` respectively. it supports specifying multiple pattern strings `{R:x}`.</li> <li>when purchasing multiple instances, if no pattern string is specified, a suffix `1, 2...n` will be added to the instance display name, where `n` represents the number of purchased instances. for example, for `server_`, when purchasing 2 instances, the instance display names will be `server_1` and `server_2` respectively.</li> <li>it supports up to 128 characters (including pattern strings).</li>.
	InstanceName *string `json:"InstanceName,omitnil,omitempty" name:"InstanceName"`

	// Instance login settings. you can use this parameter to set the instance's login method to key or keep the original login settings of the image.
	LoginSettings *LoginSettings `json:"LoginSettings,omitnil,omitempty" name:"LoginSettings"`

	// Security groups to which the instance belongs. To obtain the security group IDs, you can call [DescribeSecurityGroups](https://intl.cloud.tencent.com/document/api/215/15808) and look for the `sgld` fields in the response. If this parameter is not specified, the instance will be associated with default security groups.
	SecurityGroupIds []*string `json:"SecurityGroupIds,omitnil,omitempty" name:"SecurityGroupIds"`

	// Enhanced service. this parameter can be used to specify whether to enable services such as cloud security and cloud monitoring. if this parameter is not specified, cloud monitor and cloud security services are enabled for public images by default; for custom images and marketplace images, cloud monitor and cloud security services are not enabled by default, and the services retained in the image will be used.
	EnhancedService *EnhancedService `json:"EnhancedService,omitnil,omitempty" name:"EnhancedService"`

	// A unique string supplied by the client to ensure that the request is idempotent. Its maximum length is 64 ASCII characters. If this parameter is not specified, the idempotency of the request cannot be guaranteed.
	ClientToken *string `json:"ClientToken,omitnil,omitempty" name:"ClientToken"`

	// Hostname of Cloud Virtual Machine.<br><li>Period (.) and hyphen (-) should not be used as the first or last character of the hostname, and should not be used consecutively.</li><br><li>Windows instances: The hostname should contain 2 to 15 characters, including letters (case insensitive), digits, and hyphens (-), does not support periods (.), and should not be all digits.</li><br><li>Instances of other types (such as Linux instances): The hostname should contain 2 to 60 characters, including multiple periods (.), with each segment between periods considered as one section. Each section can contain letters (case insensitive), digits, and hyphens (-).</li>
	HostName *string `json:"HostName,omitnil,omitempty" name:"HostName"`

	// Scheduled tasks. You can use this parameter to specify scheduled tasks for the instance. Only scheduled termination is supported.
	ActionTimer *ActionTimer `json:"ActionTimer,omitnil,omitempty" name:"ActionTimer"`

	// Placement group ID. You can only specify one.
	DisasterRecoverGroupIds []*string `json:"DisasterRecoverGroupIds,omitnil,omitempty" name:"DisasterRecoverGroupIds"`

	// The tag description list. This parameter is used to bind a tag to a resource instance. A tag can only be bound to CVM instances.
	TagSpecification []*TagSpecification `json:"TagSpecification,omitnil,omitempty" name:"TagSpecification"`

	// Market options of the instance, such as parameters related to spot instances. This parameter is required for spot instances.
	InstanceMarketOptions *InstanceMarketOptionsRequest `json:"InstanceMarketOptions,omitnil,omitempty" name:"InstanceMarketOptions"`

	// User data provided to the instance. This parameter needs to be encoded in base64 format with the maximum size of 16KB. For more information on how to get the value of this parameter, see the commands you need to execute on startup for [Windows](https://intl.cloud.tencent.com/document/product/213/17526) or [Linux](https://intl.cloud.tencent.com/document/product/213/17525).
	UserData *string `json:"UserData,omitnil,omitempty" name:"UserData"`

	// Whether the request is a dry run only.
	// `true`: dry run only. The request will not create instance(s). A dry run can check whether all the required parameters are specified, whether the request format is right, whether the request exceeds service limits, and whether the specified CVMs are available.
	// If the dry run fails, the corresponding error code will be returned.
	// If the dry run succeeds, the RequestId will be returned.
	// `false` (default value): send a normal request and create instance(s) if all the requirements are met.
	DryRun *bool `json:"DryRun,omitnil,omitempty" name:"DryRun"`

	// CAM role name, which can be obtained from the `roleName` field in the response of the [`DescribeRoleList`](https://intl.cloud.tencent.com/document/product/598/36223?from_cn_redirect=1) API.
	CamRoleName *string `json:"CamRoleName,omitnil,omitempty" name:"CamRoleName"`

	// HPC cluster ID. The HPC cluster must and can only be specified for a high-performance computing instance.
	HpcClusterId *string `json:"HpcClusterId,omitnil,omitempty" name:"HpcClusterId"`

	// Instance [billing mode](https://intl.cloud.tencent.com/document/product/213/2180?from_cn_redirect=1).<br><li>PREPAID: prepaid, that is, monthly subscription.</li><li>POSTPAID_BY_HOUR: pay-as-you-go by hour.</li><li>CDHPAID: CDH instance (created based on CDH; the resources of the host are free of charge).</li><li>SPOTPAID: spot payment.</li>Default value: POSTPAID_BY_HOUR.
	InstanceChargeType *string `json:"InstanceChargeType,omitnil,omitempty" name:"InstanceChargeType"`

	// Details of the monthly subscription, including the purchase period, auto-renewal. It is required if the `InstanceChargeType` is `PREPAID`.
	InstanceChargePrepaid *InstanceChargePrepaid `json:"InstanceChargePrepaid,omitnil,omitempty" name:"InstanceChargePrepaid"`

	// Instance termination protection flag, indicating whether an instance is allowed to be deleted through an API. Valid values:<br><li>TRUE: Instance protection is enabled, and the instance is not allowed to be deleted through the API.</li><br><li>FALSE: Instance protection is disabled, and the instance is allowed to be deleted through the API.</li><br><br>Default value: FALSE.
	DisableApiTermination *bool `json:"DisableApiTermination,omitnil,omitempty" name:"DisableApiTermination"`

	EnableJumboFrame *bool `json:"EnableJumboFrame,omitnil,omitempty" name:"EnableJumboFrame"`

	// Custom metadata. specifies that custom metadata key-value pairs can be added when creating a CVM.
	// Note: this field is in beta test.
	Metadata *Metadata `json:"Metadata,omitnil,omitempty" name:"Metadata"`

	// Specifies that only the Update and Replace parameters are allowed. this parameter is valid only when custom Metadata is used in the template and Metadata is also transmitted in RunInstances. defaults to Replace.
	// 
	// -Update: if template t contains this parameter with a value of Update and metadata=[k1:v1, k2:v2], then RunInstances (with metadata=[k2:v3]) + t creates a cvm using metadata=[k1:v1, k2:v3]. 
	// -Replace: if the template t contains this parameter with a value of Replace and metadata=[k1:v1, k2:v2], then when creating a cvm using RunInstances (with metadata=[k2:v3]) + t, the created cvm will use metadata=[k2:v3]. 
	// Note: this field is in beta test.
	TemplateDataModifyAction *string `json:"TemplateDataModifyAction,omitnil,omitempty" name:"TemplateDataModifyAction"`
}

func (r *CreateLaunchTemplateVersionRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateLaunchTemplateVersionRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Placement")
	delete(f, "LaunchTemplateId")
	delete(f, "LaunchTemplateVersion")
	delete(f, "LaunchTemplateVersionDescription")
	delete(f, "InstanceType")
	delete(f, "ImageId")
	delete(f, "SystemDisk")
	delete(f, "DataDisks")
	delete(f, "VirtualPrivateCloud")
	delete(f, "InternetAccessible")
	delete(f, "InstanceCount")
	delete(f, "InstanceName")
	delete(f, "LoginSettings")
	delete(f, "SecurityGroupIds")
	delete(f, "EnhancedService")
	delete(f, "ClientToken")
	delete(f, "HostName")
	delete(f, "ActionTimer")
	delete(f, "DisasterRecoverGroupIds")
	delete(f, "TagSpecification")
	delete(f, "InstanceMarketOptions")
	delete(f, "UserData")
	delete(f, "DryRun")
	delete(f, "CamRoleName")
	delete(f, "HpcClusterId")
	delete(f, "InstanceChargeType")
	delete(f, "InstanceChargePrepaid")
	delete(f, "DisableApiTermination")
	delete(f, "EnableJumboFrame")
	delete(f, "Metadata")
	delete(f, "TemplateDataModifyAction")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CreateLaunchTemplateVersionRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateLaunchTemplateVersionResponseParams struct {
	// Version number of the new instance launch template.
	LaunchTemplateVersionNumber *int64 `json:"LaunchTemplateVersionNumber,omitnil,omitempty" name:"LaunchTemplateVersionNumber"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type CreateLaunchTemplateVersionResponse struct {
	*tchttp.BaseResponse
	Response *CreateLaunchTemplateVersionResponseParams `json:"Response"`
}

func (r *CreateLaunchTemplateVersionResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateLaunchTemplateVersionResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type DataDisk struct {
	// Data disk size, unit: GiB. the minimum adjustment step size is 10 GiB. the value ranges of different data disk types vary. for specific limitations, see the storage overview (https://intl.cloud.tencent.com/document/product/213/4952?from_cn_redirect=1). the default value is 0, which means no data disk purchase. for more restrictions, see the product document.
	DiskSize *int64 `json:"DiskSize,omitnil,omitempty" name:"DiskSize"`

	// Specifies the data disk type. for restrictions on data disk types, refer to [storage overview](https://www.tencentcloud.com/document/product/213/4952?from_cn_redirect=1). valid values: <br /><li>LOCAL_BASIC: LOCAL disk</li> <li>LOCAL_SSD: LOCAL SSD</li><li>LOCAL_NVME: LOCAL NVME disk, which is closely related to InstanceType and cannot be specified</li><li>LOCAL_PRO: LOCAL HDD, which is closely related to InstanceType and cannot be specified</li><li>cloud_BASIC: BASIC cloud disk</li><li>cloud_PREMIUM: high-performance cloud block storage</li><li>cloud_SSD: SSD cloud disk</li><li>cloud_HSSD: enhanced SSD cloud disk</li> <li>cloud_TSSD: ultra-fast SSD cbs</li><li>cloud_BSSD: universal SSD cloud disk</li><br />default: LOCAL_BASIC.<br/><br />this parameter is invalid for the `ResizeInstanceDisk` api.
	DiskType *string `json:"DiskType,omitnil,omitempty" name:"DiskType"`

	// Specifies the data disk ID.
	// This parameter currently only serves as a response parameter for query apis such as `DescribeInstances`, and cannot be used as an input parameter for write apis such as `RunInstances`.
	DiskId *string `json:"DiskId,omitnil,omitempty" name:"DiskId"`

	// Whether the data disk is terminated with the instance. value range: <li>true: when the instance is terminated, the data disk is also terminated. only hourly postpaid cloud disks are supported. <li>false: when the instance is terminated, the data disk is retained. <br>default value: true <br>currently, this parameter is only used for the API `RunInstances`.
	DeleteWithInstance *bool `json:"DeleteWithInstance,omitnil,omitempty" name:"DeleteWithInstance"`

	// Data disk snapshot ID. the size of the selected data disk snapshot must be less than the data disk size.
	SnapshotId *string `json:"SnapshotId,omitnil,omitempty" name:"SnapshotId"`

	// Specifies whether the data disk is encrypted. value range: <li>true: encrypted</li> <li>false: unencrypted</li><br/> default value: false<br/> this parameter is currently only used for the `RunInstances` api.
	Encrypt *bool `json:"Encrypt,omitnil,omitempty" name:"Encrypt"`

	// Custom CMK ID, value is UUID or similar to kms-abcd1234. used for encrypted cloud disk.
	// 
	// This parameter is currently only used for the `RunInstances` api.
	// Note: This field may return null, indicating that no valid value is found.
	KmsKeyId *string `json:"KmsKeyId,omitnil,omitempty" name:"KmsKeyId"`

	// Cloud disk performance (unit: MiB/s). specifies additional performance for cloud disks.
	// Currently only supports extreme cbs (CLOUD_TSSD) and enhanced SSD CLOUD disk (CLOUD_HSSD).
	// Note: This field may return null, indicating that no valid value is found.
	ThroughputPerformance *int64 `json:"ThroughputPerformance,omitnil,omitempty" name:"ThroughputPerformance"`

	// Specifies the dedicated cluster ID belonging to.
	// Note: This field may return null, indicating that no valid value is found.
	CdcId *string `json:"CdcId,omitnil,omitempty" name:"CdcId"`

	// Burstable performance.
	// 
	// <B>Note: this field is in beta test.</b>.
	BurstPerformance *bool `json:"BurstPerformance,omitnil,omitempty" name:"BurstPerformance"`

	// Disk name, with a length not exceeding 128 characters.
	DiskName *string `json:"DiskName,omitnil,omitempty" name:"DiskName"`
}

// Predefined struct for user
type DeleteDisasterRecoverGroupsRequestParams struct {
	// ID list of spread placement groups, obtainable via the [DescribeDisasterRecoverGroups](https://intl.cloud.tencent.com/document/api/213/17810?from_cn_redirect=1) API. You can operate up to 10 spread placement groups in each request.
	DisasterRecoverGroupIds []*string `json:"DisasterRecoverGroupIds,omitnil,omitempty" name:"DisasterRecoverGroupIds"`
}

type DeleteDisasterRecoverGroupsRequest struct {
	*tchttp.BaseRequest
	
	// ID list of spread placement groups, obtainable via the [DescribeDisasterRecoverGroups](https://intl.cloud.tencent.com/document/api/213/17810?from_cn_redirect=1) API. You can operate up to 10 spread placement groups in each request.
	DisasterRecoverGroupIds []*string `json:"DisasterRecoverGroupIds,omitnil,omitempty" name:"DisasterRecoverGroupIds"`
}

func (r *DeleteDisasterRecoverGroupsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteDisasterRecoverGroupsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "DisasterRecoverGroupIds")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DeleteDisasterRecoverGroupsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteDisasterRecoverGroupsResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DeleteDisasterRecoverGroupsResponse struct {
	*tchttp.BaseResponse
	Response *DeleteDisasterRecoverGroupsResponseParams `json:"Response"`
}

func (r *DeleteDisasterRecoverGroupsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteDisasterRecoverGroupsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteImagesRequestParams struct {
	// List of the IDs of the instances to be deleted.
	ImageIds []*string `json:"ImageIds,omitnil,omitempty" name:"ImageIds"`

	// Whether to delete the snapshot associated with the image
	DeleteBindedSnap *bool `json:"DeleteBindedSnap,omitnil,omitempty" name:"DeleteBindedSnap"`

	// Check whether deleting an image is supported
	DryRun *bool `json:"DryRun,omitnil,omitempty" name:"DryRun"`
}

type DeleteImagesRequest struct {
	*tchttp.BaseRequest
	
	// List of the IDs of the instances to be deleted.
	ImageIds []*string `json:"ImageIds,omitnil,omitempty" name:"ImageIds"`

	// Whether to delete the snapshot associated with the image
	DeleteBindedSnap *bool `json:"DeleteBindedSnap,omitnil,omitempty" name:"DeleteBindedSnap"`

	// Check whether deleting an image is supported
	DryRun *bool `json:"DryRun,omitnil,omitempty" name:"DryRun"`
}

func (r *DeleteImagesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteImagesRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ImageIds")
	delete(f, "DeleteBindedSnap")
	delete(f, "DryRun")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DeleteImagesRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteImagesResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DeleteImagesResponse struct {
	*tchttp.BaseResponse
	Response *DeleteImagesResponseParams `json:"Response"`
}

func (r *DeleteImagesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteImagesResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteKeyPairsRequestParams struct {
	// One or more key pair IDs to be operated. The maximum number of key pairs per request is 100.<br>You can obtain an available key pair ID in the following ways:<br><li>Log in to the [console](https://console.cloud.tencent.com/cvm/sshkey) to query the key pair ID.</li><br><li>Call the [DescribeKeyPairs](https://intl.cloud.tencent.com/document/api/213/15699?from_cn_redirect=1) API and obtain the key pair ID from the `KeyId` in the response.</li>
	KeyIds []*string `json:"KeyIds,omitnil,omitempty" name:"KeyIds"`
}

type DeleteKeyPairsRequest struct {
	*tchttp.BaseRequest
	
	// One or more key pair IDs to be operated. The maximum number of key pairs per request is 100.<br>You can obtain an available key pair ID in the following ways:<br><li>Log in to the [console](https://console.cloud.tencent.com/cvm/sshkey) to query the key pair ID.</li><br><li>Call the [DescribeKeyPairs](https://intl.cloud.tencent.com/document/api/213/15699?from_cn_redirect=1) API and obtain the key pair ID from the `KeyId` in the response.</li>
	KeyIds []*string `json:"KeyIds,omitnil,omitempty" name:"KeyIds"`
}

func (r *DeleteKeyPairsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteKeyPairsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "KeyIds")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DeleteKeyPairsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteKeyPairsResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DeleteKeyPairsResponse struct {
	*tchttp.BaseResponse
	Response *DeleteKeyPairsResponseParams `json:"Response"`
}

func (r *DeleteKeyPairsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteKeyPairsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteLaunchTemplateRequestParams struct {
	// The launch template ID
	LaunchTemplateId *string `json:"LaunchTemplateId,omitnil,omitempty" name:"LaunchTemplateId"`
}

type DeleteLaunchTemplateRequest struct {
	*tchttp.BaseRequest
	
	// The launch template ID
	LaunchTemplateId *string `json:"LaunchTemplateId,omitnil,omitempty" name:"LaunchTemplateId"`
}

func (r *DeleteLaunchTemplateRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteLaunchTemplateRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "LaunchTemplateId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DeleteLaunchTemplateRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteLaunchTemplateResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DeleteLaunchTemplateResponse struct {
	*tchttp.BaseResponse
	Response *DeleteLaunchTemplateResponseParams `json:"Response"`
}

func (r *DeleteLaunchTemplateResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteLaunchTemplateResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteLaunchTemplateVersionsRequestParams struct {
	// The launch template ID
	LaunchTemplateId *string `json:"LaunchTemplateId,omitnil,omitempty" name:"LaunchTemplateId"`

	// List of instance launch template versions.
	LaunchTemplateVersions []*int64 `json:"LaunchTemplateVersions,omitnil,omitempty" name:"LaunchTemplateVersions"`
}

type DeleteLaunchTemplateVersionsRequest struct {
	*tchttp.BaseRequest
	
	// The launch template ID
	LaunchTemplateId *string `json:"LaunchTemplateId,omitnil,omitempty" name:"LaunchTemplateId"`

	// List of instance launch template versions.
	LaunchTemplateVersions []*int64 `json:"LaunchTemplateVersions,omitnil,omitempty" name:"LaunchTemplateVersions"`
}

func (r *DeleteLaunchTemplateVersionsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteLaunchTemplateVersionsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "LaunchTemplateId")
	delete(f, "LaunchTemplateVersions")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DeleteLaunchTemplateVersionsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteLaunchTemplateVersionsResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DeleteLaunchTemplateVersionsResponse struct {
	*tchttp.BaseResponse
	Response *DeleteLaunchTemplateVersionsResponseParams `json:"Response"`
}

func (r *DeleteLaunchTemplateVersionsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteLaunchTemplateVersionsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeChcDeniedActionsRequestParams struct {
	// CHC instance ID
	ChcIds []*string `json:"ChcIds,omitnil,omitempty" name:"ChcIds"`
}

type DescribeChcDeniedActionsRequest struct {
	*tchttp.BaseRequest
	
	// CHC instance ID
	ChcIds []*string `json:"ChcIds,omitnil,omitempty" name:"ChcIds"`
}

func (r *DescribeChcDeniedActionsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeChcDeniedActionsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ChcIds")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeChcDeniedActionsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeChcDeniedActionsResponseParams struct {
	// Actions not allowed for the CHC instance
	ChcHostDeniedActionSet []*ChcHostDeniedActions `json:"ChcHostDeniedActionSet,omitnil,omitempty" name:"ChcHostDeniedActionSet"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeChcDeniedActionsResponse struct {
	*tchttp.BaseResponse
	Response *DescribeChcDeniedActionsResponseParams `json:"Response"`
}

func (r *DescribeChcDeniedActionsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeChcDeniedActionsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeChcHostsRequestParams struct {
	// CHC host ID. Up to 100 instances per request is allowed. `ChcIds` and `Filters` cannot be specified at the same time.
	ChcIds []*string `json:"ChcIds,omitnil,omitempty" name:"ChcIds"`

	// <li><strong>zone</strong></li>
	// <p style="padding-left: 30px;">Filter by the <strong>availability zone</strong>, such as `ap-guangzhou-6`.</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Optional</p><p style="padding-left: 30px;">Valid values: See <a href="https://intl.cloud.tencent.com/document/product/213/6091?from_cn_redirect=1">Regions and Availability Zones</a></p>
	// <li><strong>instance-name</strong></li>
	// <p style="padding-left: 30px;">Filter by the <strong>instance name</strong>.</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Optional</p>
	// <li><strong>instance-state</strong></li>
	// <p style="padding-left: 30px;">Filter by the <strong>instance status</strong>. For status details, see [InstanceStatus](https://intl.cloud.tencent.com/document/api/213/15753?from_cn_redirect=1#InstanceStatus).</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Optional</p>
	// <li><strong>device-type</strong></li>
	// <p style="padding-left: 30px;">Filter by the <strong>device type</strong>.</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Optional</p>
	// <li><strong>vpc-id</strong></li>
	// <p style="padding-left: 30px;">Filter by the <strong>unique VPC ID</strong>.</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Optional</p>
	// <li><strong>subnet-id</strong></li>
	// <p style="padding-left: 30px;">Filter by the <strong>unique VPC subnet ID</strong>.</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Optional</p>
	Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`

	// The offset. Default value: `0`. For more information on `Offset`, see the relevant sections in API [Introduction](https://intl.cloud.tencent.com/document/api/213/15688?from_cn_redirect=1).
	Offset *int64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// The number of returned results. Default value: `20`. Maximum value: `100`. For more information on `Limit`, see the relevant sections in API [Introduction](https://intl.cloud.tencent.com/document/api/213/15688?from_cn_redirect=1).
	Limit *int64 `json:"Limit,omitnil,omitempty" name:"Limit"`
}

type DescribeChcHostsRequest struct {
	*tchttp.BaseRequest
	
	// CHC host ID. Up to 100 instances per request is allowed. `ChcIds` and `Filters` cannot be specified at the same time.
	ChcIds []*string `json:"ChcIds,omitnil,omitempty" name:"ChcIds"`

	// <li><strong>zone</strong></li>
	// <p style="padding-left: 30px;">Filter by the <strong>availability zone</strong>, such as `ap-guangzhou-6`.</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Optional</p><p style="padding-left: 30px;">Valid values: See <a href="https://intl.cloud.tencent.com/document/product/213/6091?from_cn_redirect=1">Regions and Availability Zones</a></p>
	// <li><strong>instance-name</strong></li>
	// <p style="padding-left: 30px;">Filter by the <strong>instance name</strong>.</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Optional</p>
	// <li><strong>instance-state</strong></li>
	// <p style="padding-left: 30px;">Filter by the <strong>instance status</strong>. For status details, see [InstanceStatus](https://intl.cloud.tencent.com/document/api/213/15753?from_cn_redirect=1#InstanceStatus).</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Optional</p>
	// <li><strong>device-type</strong></li>
	// <p style="padding-left: 30px;">Filter by the <strong>device type</strong>.</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Optional</p>
	// <li><strong>vpc-id</strong></li>
	// <p style="padding-left: 30px;">Filter by the <strong>unique VPC ID</strong>.</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Optional</p>
	// <li><strong>subnet-id</strong></li>
	// <p style="padding-left: 30px;">Filter by the <strong>unique VPC subnet ID</strong>.</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Optional</p>
	Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`

	// The offset. Default value: `0`. For more information on `Offset`, see the relevant sections in API [Introduction](https://intl.cloud.tencent.com/document/api/213/15688?from_cn_redirect=1).
	Offset *int64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// The number of returned results. Default value: `20`. Maximum value: `100`. For more information on `Limit`, see the relevant sections in API [Introduction](https://intl.cloud.tencent.com/document/api/213/15688?from_cn_redirect=1).
	Limit *int64 `json:"Limit,omitnil,omitempty" name:"Limit"`
}

func (r *DescribeChcHostsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeChcHostsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ChcIds")
	delete(f, "Filters")
	delete(f, "Offset")
	delete(f, "Limit")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeChcHostsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeChcHostsResponseParams struct {
	// Number of eligible instances
	TotalCount *int64 `json:"TotalCount,omitnil,omitempty" name:"TotalCount"`

	// List of returned instances
	ChcHostSet []*ChcHost `json:"ChcHostSet,omitnil,omitempty" name:"ChcHostSet"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeChcHostsResponse struct {
	*tchttp.BaseResponse
	Response *DescribeChcHostsResponseParams `json:"Response"`
}

func (r *DescribeChcHostsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeChcHostsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeDisasterRecoverGroupQuotaRequestParams struct {

}

type DescribeDisasterRecoverGroupQuotaRequest struct {
	*tchttp.BaseRequest
	
}

func (r *DescribeDisasterRecoverGroupQuotaRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeDisasterRecoverGroupQuotaRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeDisasterRecoverGroupQuotaRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeDisasterRecoverGroupQuotaResponseParams struct {
	// The maximum number of placement groups that can be created.
	GroupQuota *int64 `json:"GroupQuota,omitnil,omitempty" name:"GroupQuota"`

	// The number of placement groups that have been created by the current user.
	CurrentNum *int64 `json:"CurrentNum,omitnil,omitempty" name:"CurrentNum"`

	// Quota on instances in a physical-machine-type disaster recovery group.
	CvmInHostGroupQuota *int64 `json:"CvmInHostGroupQuota,omitnil,omitempty" name:"CvmInHostGroupQuota"`

	// Quota on instances in a switch-type disaster recovery group.
	//
	// Deprecated: CvmInSwGroupQuota is deprecated.
	CvmInSwGroupQuota *int64 `json:"CvmInSwGroupQuota,omitnil,omitempty" name:"CvmInSwGroupQuota"`

	// Quota on instances in a rack-type disaster recovery group.
	CvmInRackGroupQuota *int64 `json:"CvmInRackGroupQuota,omitnil,omitempty" name:"CvmInRackGroupQuota"`


	CvmInSwitchGroupQuota *int64 `json:"CvmInSwitchGroupQuota,omitnil,omitempty" name:"CvmInSwitchGroupQuota"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeDisasterRecoverGroupQuotaResponse struct {
	*tchttp.BaseResponse
	Response *DescribeDisasterRecoverGroupQuotaResponseParams `json:"Response"`
}

func (r *DescribeDisasterRecoverGroupQuotaResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeDisasterRecoverGroupQuotaResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeDisasterRecoverGroupsRequestParams struct {
	// ID list of spread placement groups. You can operate up to 10 spread placement groups in each request.
	DisasterRecoverGroupIds []*string `json:"DisasterRecoverGroupIds,omitnil,omitempty" name:"DisasterRecoverGroupIds"`

	// Name of a spread placement group. Fuzzy match is supported.
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// Offset; default value: 0. For more information on `Offset`, see the corresponding section in API [Introduction](https://intl.cloud.tencent.com/document/product/377).
	Offset *int64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// Number of results returned; default value: 20; maximum: 100. For more information on `Limit`, see the corresponding section in API [Introduction](https://intl.cloud.tencent.com/document/product/377). 
	Limit *int64 `json:"Limit,omitnil,omitempty" name:"Limit"`

	// <li> `tag-key` - String - Optional - Filter by the tag key.</li>
	// <li>`tag-value` - String - Optional - Filter by the tag value.</li>
	// <li> `tag:tag-key` - String - Optional - Filter by the tag key-value pair. Replace `tag-key` with the actual tag keys.</li>
	// Each request can have up to 10 `Filters` and 5 `Filters.Values`.
	Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`
}

type DescribeDisasterRecoverGroupsRequest struct {
	*tchttp.BaseRequest
	
	// ID list of spread placement groups. You can operate up to 10 spread placement groups in each request.
	DisasterRecoverGroupIds []*string `json:"DisasterRecoverGroupIds,omitnil,omitempty" name:"DisasterRecoverGroupIds"`

	// Name of a spread placement group. Fuzzy match is supported.
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// Offset; default value: 0. For more information on `Offset`, see the corresponding section in API [Introduction](https://intl.cloud.tencent.com/document/product/377).
	Offset *int64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// Number of results returned; default value: 20; maximum: 100. For more information on `Limit`, see the corresponding section in API [Introduction](https://intl.cloud.tencent.com/document/product/377). 
	Limit *int64 `json:"Limit,omitnil,omitempty" name:"Limit"`

	// <li> `tag-key` - String - Optional - Filter by the tag key.</li>
	// <li>`tag-value` - String - Optional - Filter by the tag value.</li>
	// <li> `tag:tag-key` - String - Optional - Filter by the tag key-value pair. Replace `tag-key` with the actual tag keys.</li>
	// Each request can have up to 10 `Filters` and 5 `Filters.Values`.
	Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`
}

func (r *DescribeDisasterRecoverGroupsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeDisasterRecoverGroupsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "DisasterRecoverGroupIds")
	delete(f, "Name")
	delete(f, "Offset")
	delete(f, "Limit")
	delete(f, "Filters")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeDisasterRecoverGroupsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeDisasterRecoverGroupsResponseParams struct {
	// Information on spread placement groups.
	DisasterRecoverGroupSet []*DisasterRecoverGroup `json:"DisasterRecoverGroupSet,omitnil,omitempty" name:"DisasterRecoverGroupSet"`

	// Total number of placement groups of the user.
	TotalCount *int64 `json:"TotalCount,omitnil,omitempty" name:"TotalCount"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeDisasterRecoverGroupsResponse struct {
	*tchttp.BaseResponse
	Response *DescribeDisasterRecoverGroupsResponseParams `json:"Response"`
}

func (r *DescribeDisasterRecoverGroupsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeDisasterRecoverGroupsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeHostsRequestParams struct {
	// <li><strong>zone</strong></li>
	// <p style="padding-left: 30px;">Filter by the availability zone, such as `ap-guangzhou-6`.</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Optional</p><p style="padding-left: 30px;">Valid values: See <a href="https://intl.cloud.tencent.com/document/product/213/6091?from_cn_redirect=1">Regions and Availability Zones</a></p>
	// <li><strong>project-id</strong></li>
	// <p style="padding-left: 30px;">Filter by the project ID. You can query the list of created projects by calling `DescribeProject` or logging in to the [CVM console](https://console.cloud.tencent.com/cvm/index). You can also call `AddProject` to create projects. The project ID is like 1002189. </p><p style="padding-left: 30px;">Type: Integer</p><p style="padding-left: 30px;">Optional</p>
	// <li><strong>host-id</strong></li>
	// <p style="padding-left: 30px;">Filter by the CDH instance ID. Format: host-xxxxxxxx. </p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Optional</p>
	// <li><strong>host-name</strong></li>
	// <p style="padding-left: 30px;">Filter by the host name.</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Optional</p>
	// <li><strong>host-state</strong></li>
	// <p style="padding-left: 30px;">Filter by the CDH instance status. (`PENDING`: Creating | `LAUNCH_FAILURE`: Failed to create | `RUNNING`: Running | `EXPIRED`: Expired)</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Optional</p>
	// Each request can have up to 10 `Filters` and 5 `Filter.Values`.
	Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`

	// Offset; default value: 0.
	Offset *uint64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// Number of results returned; default value: 20; maximum: 100.
	Limit *uint64 `json:"Limit,omitnil,omitempty" name:"Limit"`
}

type DescribeHostsRequest struct {
	*tchttp.BaseRequest
	
	// <li><strong>zone</strong></li>
	// <p style="padding-left: 30px;">Filter by the availability zone, such as `ap-guangzhou-6`.</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Optional</p><p style="padding-left: 30px;">Valid values: See <a href="https://intl.cloud.tencent.com/document/product/213/6091?from_cn_redirect=1">Regions and Availability Zones</a></p>
	// <li><strong>project-id</strong></li>
	// <p style="padding-left: 30px;">Filter by the project ID. You can query the list of created projects by calling `DescribeProject` or logging in to the [CVM console](https://console.cloud.tencent.com/cvm/index). You can also call `AddProject` to create projects. The project ID is like 1002189. </p><p style="padding-left: 30px;">Type: Integer</p><p style="padding-left: 30px;">Optional</p>
	// <li><strong>host-id</strong></li>
	// <p style="padding-left: 30px;">Filter by the CDH instance ID. Format: host-xxxxxxxx. </p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Optional</p>
	// <li><strong>host-name</strong></li>
	// <p style="padding-left: 30px;">Filter by the host name.</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Optional</p>
	// <li><strong>host-state</strong></li>
	// <p style="padding-left: 30px;">Filter by the CDH instance status. (`PENDING`: Creating | `LAUNCH_FAILURE`: Failed to create | `RUNNING`: Running | `EXPIRED`: Expired)</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Optional</p>
	// Each request can have up to 10 `Filters` and 5 `Filter.Values`.
	Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`

	// Offset; default value: 0.
	Offset *uint64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// Number of results returned; default value: 20; maximum: 100.
	Limit *uint64 `json:"Limit,omitnil,omitempty" name:"Limit"`
}

func (r *DescribeHostsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeHostsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Filters")
	delete(f, "Offset")
	delete(f, "Limit")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeHostsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeHostsResponseParams struct {
	// Total number of CDH instances meeting the query conditions
	TotalCount *uint64 `json:"TotalCount,omitnil,omitempty" name:"TotalCount"`

	// Information on CDH instances
	HostSet []*HostItem `json:"HostSet,omitnil,omitempty" name:"HostSet"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeHostsResponse struct {
	*tchttp.BaseResponse
	Response *DescribeHostsResponseParams `json:"Response"`
}

func (r *DescribeHostsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeHostsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeImageFromFamilyRequestParams struct {
	// Image family
	ImageFamily *string `json:"ImageFamily,omitnil,omitempty" name:"ImageFamily"`
}

type DescribeImageFromFamilyRequest struct {
	*tchttp.BaseRequest
	
	// Image family
	ImageFamily *string `json:"ImageFamily,omitnil,omitempty" name:"ImageFamily"`
}

func (r *DescribeImageFromFamilyRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeImageFromFamilyRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ImageFamily")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeImageFromFamilyRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeImageFromFamilyResponseParams struct {
	// Image information. Null is returned when there is no available image.
	// Note: This field may return null, indicating that no valid value is found.
	Image *Image `json:"Image,omitnil,omitempty" name:"Image"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeImageFromFamilyResponse struct {
	*tchttp.BaseResponse
	Response *DescribeImageFromFamilyResponseParams `json:"Response"`
}

func (r *DescribeImageFromFamilyResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeImageFromFamilyResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeImageQuotaRequestParams struct {

}

type DescribeImageQuotaRequest struct {
	*tchttp.BaseRequest
	
}

func (r *DescribeImageQuotaRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeImageQuotaRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeImageQuotaRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeImageQuotaResponseParams struct {
	// The image quota of an account
	ImageNumQuota *int64 `json:"ImageNumQuota,omitnil,omitempty" name:"ImageNumQuota"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeImageQuotaResponse struct {
	*tchttp.BaseResponse
	Response *DescribeImageQuotaResponseParams `json:"Response"`
}

func (r *DescribeImageQuotaResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeImageQuotaResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeImageSharePermissionRequestParams struct {
	// The ID of the image to be shared
	ImageId *string `json:"ImageId,omitnil,omitempty" name:"ImageId"`
}

type DescribeImageSharePermissionRequest struct {
	*tchttp.BaseRequest
	
	// The ID of the image to be shared
	ImageId *string `json:"ImageId,omitnil,omitempty" name:"ImageId"`
}

func (r *DescribeImageSharePermissionRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeImageSharePermissionRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ImageId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeImageSharePermissionRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeImageSharePermissionResponseParams struct {
	// Information on image sharing.
	SharePermissionSet []*SharePermission `json:"SharePermissionSet,omitnil,omitempty" name:"SharePermissionSet"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeImageSharePermissionResponse struct {
	*tchttp.BaseResponse
	Response *DescribeImageSharePermissionResponseParams `json:"Response"`
}

func (r *DescribeImageSharePermissionResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeImageSharePermissionResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeImagesRequestParams struct {
	// List of image IDs, such as `img-gvbnzy6f`. For the format of array-type parameters, see [API Introduction](https://intl.cloud.tencent.com/document/api/213/15688?from_cn_redirect=1). You can obtain the image IDs in two ways: <br><li>Call [DescribeImages](https://intl.cloud.tencent.com/document/api/213/15715?from_cn_redirect=1) and look for `ImageId` in the response. <br><li>View the image IDs in the [Image Console](https://console.cloud.tencent.com/cvm/image).
	ImageIds []*string `json:"ImageIds,omitnil,omitempty" name:"ImageIds"`

	// Filters. Each request can have up to 10 `Filters`, and 5 `Filters.Values` for each filter. `ImageIds` and `Filters` cannot be specified at the same time. See details:
	// 
	// <li><strong>image-id</strong></li>
	// <p style="padding-left: 30px;">Filter by the <strong>image ID</strong>.</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Optional</p>
	// <li><strong>image-type</strong></li>
	// <p style="padding-left: 30px;">Filter by the <strong>image type</strong>.</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Optional</p><p style="padding-left: 30px;">Options:</p><p style="padding-left: 30px;">`PRIVATE_IMAGE`: Private images (images created by this account)</p><p style="padding-left: 30px;">`PUBLIC_IMAGE`: Public images (Tencent Cloud official images)</p><p style="padding-left: 30px;">`SHARED_IMAGE`: Shared images (images shared by other accounts to this account)</p>
	// <li><strong>image-name</strong></li>
	// <p style="padding-left: 30px;">Filter by the <strong>image name</strong>.</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Optional</p>
	// <li><strong>platform</strong></li>
	// <p style="padding-left: 30px;">Filter by the <strong>image operating system</strong>, such as `CentOS`.</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Optional</p>
	// <li><strong>tag-key</strong></li>
	// <p style="padding-left: 30px;">Filter by the <strong>tag key</strong>.</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Optional</p>
	// <li><strong>tag-value</strong></li>
	// <p style="padding-left: 30px;">Filter by the <strong>tag value</strong>.</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Optional</p>
	// <li><strong>tag:tag-key</strong></li>
	// <p style="padding-left: 30px;">Filter by the <strong>tag key-value pair</strong>. Replace "tag-key" with the actual value. </p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Optional</p>
	Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`

	// Offset; default value: 0. For more information on `Offset`, see [API Introduction](https://intl.cloud.tencent.com/document/api/213/568?from_cn_redirect=1#.E8.BE.93.E5.85.A5.E5.8F.82.E6.95.B0.E4.B8.8E.E8.BF.94.E5.9B.9E.E5.8F.82.E6.95.B0.E9.87.8A.E4.B9.89).
	Offset *uint64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// Number of results returned; default value: 20; maximum: 100. For more information on `Limit`, see [API Introduction](https://intl.cloud.tencent.com/document/api/213/568?from_cn_redirect=1#.E8.BE.93.E5.85.A5.E5.8F.82.E6.95.B0.E4.B8.8E.E8.BF.94.E5.9B.9E.E5.8F.82.E6.95.B0.E9.87.8A.E4.B9.89).
	Limit *uint64 `json:"Limit,omitnil,omitempty" name:"Limit"`

	// Instance type, e.g. `SA5.MEDIUM2`
	InstanceType *string `json:"InstanceType,omitnil,omitempty" name:"InstanceType"`
}

type DescribeImagesRequest struct {
	*tchttp.BaseRequest
	
	// List of image IDs, such as `img-gvbnzy6f`. For the format of array-type parameters, see [API Introduction](https://intl.cloud.tencent.com/document/api/213/15688?from_cn_redirect=1). You can obtain the image IDs in two ways: <br><li>Call [DescribeImages](https://intl.cloud.tencent.com/document/api/213/15715?from_cn_redirect=1) and look for `ImageId` in the response. <br><li>View the image IDs in the [Image Console](https://console.cloud.tencent.com/cvm/image).
	ImageIds []*string `json:"ImageIds,omitnil,omitempty" name:"ImageIds"`

	// Filters. Each request can have up to 10 `Filters`, and 5 `Filters.Values` for each filter. `ImageIds` and `Filters` cannot be specified at the same time. See details:
	// 
	// <li><strong>image-id</strong></li>
	// <p style="padding-left: 30px;">Filter by the <strong>image ID</strong>.</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Optional</p>
	// <li><strong>image-type</strong></li>
	// <p style="padding-left: 30px;">Filter by the <strong>image type</strong>.</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Optional</p><p style="padding-left: 30px;">Options:</p><p style="padding-left: 30px;">`PRIVATE_IMAGE`: Private images (images created by this account)</p><p style="padding-left: 30px;">`PUBLIC_IMAGE`: Public images (Tencent Cloud official images)</p><p style="padding-left: 30px;">`SHARED_IMAGE`: Shared images (images shared by other accounts to this account)</p>
	// <li><strong>image-name</strong></li>
	// <p style="padding-left: 30px;">Filter by the <strong>image name</strong>.</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Optional</p>
	// <li><strong>platform</strong></li>
	// <p style="padding-left: 30px;">Filter by the <strong>image operating system</strong>, such as `CentOS`.</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Optional</p>
	// <li><strong>tag-key</strong></li>
	// <p style="padding-left: 30px;">Filter by the <strong>tag key</strong>.</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Optional</p>
	// <li><strong>tag-value</strong></li>
	// <p style="padding-left: 30px;">Filter by the <strong>tag value</strong>.</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Optional</p>
	// <li><strong>tag:tag-key</strong></li>
	// <p style="padding-left: 30px;">Filter by the <strong>tag key-value pair</strong>. Replace "tag-key" with the actual value. </p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Optional</p>
	Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`

	// Offset; default value: 0. For more information on `Offset`, see [API Introduction](https://intl.cloud.tencent.com/document/api/213/568?from_cn_redirect=1#.E8.BE.93.E5.85.A5.E5.8F.82.E6.95.B0.E4.B8.8E.E8.BF.94.E5.9B.9E.E5.8F.82.E6.95.B0.E9.87.8A.E4.B9.89).
	Offset *uint64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// Number of results returned; default value: 20; maximum: 100. For more information on `Limit`, see [API Introduction](https://intl.cloud.tencent.com/document/api/213/568?from_cn_redirect=1#.E8.BE.93.E5.85.A5.E5.8F.82.E6.95.B0.E4.B8.8E.E8.BF.94.E5.9B.9E.E5.8F.82.E6.95.B0.E9.87.8A.E4.B9.89).
	Limit *uint64 `json:"Limit,omitnil,omitempty" name:"Limit"`

	// Instance type, e.g. `SA5.MEDIUM2`
	InstanceType *string `json:"InstanceType,omitnil,omitempty" name:"InstanceType"`
}

func (r *DescribeImagesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeImagesRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ImageIds")
	delete(f, "Filters")
	delete(f, "Offset")
	delete(f, "Limit")
	delete(f, "InstanceType")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeImagesRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeImagesResponseParams struct {
	// Information on an image, including its state and attributes.
	ImageSet []*Image `json:"ImageSet,omitnil,omitempty" name:"ImageSet"`

	// Number of images meeting the filtering conditions.
	TotalCount *int64 `json:"TotalCount,omitnil,omitempty" name:"TotalCount"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeImagesResponse struct {
	*tchttp.BaseResponse
	Response *DescribeImagesResponseParams `json:"Response"`
}

func (r *DescribeImagesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeImagesResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeImportImageOsRequestParams struct {

}

type DescribeImportImageOsRequest struct {
	*tchttp.BaseRequest
	
}

func (r *DescribeImportImageOsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeImportImageOsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeImportImageOsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeImportImageOsResponseParams struct {
	// Supported operating system types of imported images.
	ImportImageOsListSupported *ImageOsList `json:"ImportImageOsListSupported,omitnil,omitempty" name:"ImportImageOsListSupported"`

	// Supported operating system versions of imported images. 
	ImportImageOsVersionSet []*OsVersion `json:"ImportImageOsVersionSet,omitnil,omitempty" name:"ImportImageOsVersionSet"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeImportImageOsResponse struct {
	*tchttp.BaseResponse
	Response *DescribeImportImageOsResponseParams `json:"Response"`
}

func (r *DescribeImportImageOsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeImportImageOsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeInstanceFamilyConfigsRequestParams struct {

}

type DescribeInstanceFamilyConfigsRequest struct {
	*tchttp.BaseRequest
	
}

func (r *DescribeInstanceFamilyConfigsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeInstanceFamilyConfigsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeInstanceFamilyConfigsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeInstanceFamilyConfigsResponseParams struct {
	// List of instance model families
	InstanceFamilyConfigSet []*InstanceFamilyConfig `json:"InstanceFamilyConfigSet,omitnil,omitempty" name:"InstanceFamilyConfigSet"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeInstanceFamilyConfigsResponse struct {
	*tchttp.BaseResponse
	Response *DescribeInstanceFamilyConfigsResponseParams `json:"Response"`
}

func (r *DescribeInstanceFamilyConfigsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeInstanceFamilyConfigsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeInstancesAttributesRequestParams struct {
	// Instance attributes to be obtained. Valid value(s): 
	// UserData: Custom data of instances.
	Attributes []*string `json:"Attributes,omitnil,omitempty" name:"Attributes"`

	// Instance ID list.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`
}

type DescribeInstancesAttributesRequest struct {
	*tchttp.BaseRequest
	
	// Instance attributes to be obtained. Valid value(s): 
	// UserData: Custom data of instances.
	Attributes []*string `json:"Attributes,omitnil,omitempty" name:"Attributes"`

	// Instance ID list.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`
}

func (r *DescribeInstancesAttributesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeInstancesAttributesRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Attributes")
	delete(f, "InstanceIds")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeInstancesAttributesRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeInstancesAttributesResponseParams struct {
	// List of attributes of specified instances.
	InstanceSet []*InstanceAttribute `json:"InstanceSet,omitnil,omitempty" name:"InstanceSet"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeInstancesAttributesResponse struct {
	*tchttp.BaseResponse
	Response *DescribeInstancesAttributesResponseParams `json:"Response"`
}

func (r *DescribeInstancesAttributesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeInstancesAttributesResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeInstancesOperationLimitRequestParams struct {
	// Query by instance ID(s). You can obtain the instance IDs from the value of `InstanceId` returned by the [DescribeInstances](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1) API. For example, instance ID: ins-xxxxxxxx. (For the specific format, refer to section `ids.N` of the API [Introduction](https://intl.cloud.tencent.com/document/api/213/15688?from_cn_redirect=1).) You can query up to 100 instances in each request.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// Operation on the instance(s).
	// <li> INSTANCE_DEGRADE: downgrade the instance configurations</li>
	Operation *string `json:"Operation,omitnil,omitempty" name:"Operation"`
}

type DescribeInstancesOperationLimitRequest struct {
	*tchttp.BaseRequest
	
	// Query by instance ID(s). You can obtain the instance IDs from the value of `InstanceId` returned by the [DescribeInstances](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1) API. For example, instance ID: ins-xxxxxxxx. (For the specific format, refer to section `ids.N` of the API [Introduction](https://intl.cloud.tencent.com/document/api/213/15688?from_cn_redirect=1).) You can query up to 100 instances in each request.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// Operation on the instance(s).
	// <li> INSTANCE_DEGRADE: downgrade the instance configurations</li>
	Operation *string `json:"Operation,omitnil,omitempty" name:"Operation"`
}

func (r *DescribeInstancesOperationLimitRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeInstancesOperationLimitRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceIds")
	delete(f, "Operation")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeInstancesOperationLimitRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeInstancesOperationLimitResponseParams struct {
	// The maximum number of times you can modify the instance configurations (degrading the configurations)
	InstanceOperationLimitSet []*OperationCountLimit `json:"InstanceOperationLimitSet,omitnil,omitempty" name:"InstanceOperationLimitSet"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeInstancesOperationLimitResponse struct {
	*tchttp.BaseResponse
	Response *DescribeInstancesOperationLimitResponseParams `json:"Response"`
}

func (r *DescribeInstancesOperationLimitResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeInstancesOperationLimitResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeInstancesRequestParams struct {
	// Query by instance ID(s). For example, instance ID: `ins-xxxxxxxx`. For the specific format, refer to section `Ids.N` of the API [Introduction](https://intl.cloud.tencent.com/document/api/213/15688?from_cn_redirect=1). You can query up to 100 instances in each request. However, `InstanceIds` and `Filters` cannot be specified at the same time.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// Filters
	// <li> `zone` - String - Optional - Filter by the availability zone.</li>
	// <li> `project-id` - Integer - Optional - Filter by the project ID. You can query the list of created projects by calling `DescribeProject` or logging in to the [CVM console](https://console.cloud.tencent.com/cvm/index). You can also call `AddProject` to create projects. </li>
	// <li> `host-id` - String - Optional - Filter by the CDH instance ID. Format: `host-xxxxxxxx`.</li>
	// </li>`vpc-id` - String - Optional - Filter by the VPC ID. Format: `vpc-xxxxxxxx`.</li>
	// <li> `subnet-id` - String - Optional - Filter by the subnet ID. Format: `subnet-xxxxxxxx`.</li>
	// </li>`instance-id` - String - Optional - Filter by the instance ID. Format: `ins-xxxxxxxx`.</li>
	// </li>`security-group-id` - String - Optional - Filter by the security group ID. Format: `sg-8jlk3f3r`.</li>
	// </li>`instance-name` - String - Optional - Filter by the instance name.</li>
	// </li>`instance-charge-type` - String - Optional - Filter by the instance billing method. `POSTPAID_BY_HOUR`: pay-as-you-go | `CDHPAID`: You are only billed for [CDH](https://intl.cloud.tencent.com/document/product/416?from_cn_redirect=1) instances, not the CVMs running on the [CDH](https://intl.cloud.tencent.com/document/product/416?from_cn_redirect=1) instances.</li>
	// </li>`private-ip-address` - String - Optional - Filter by the private IP address of the instance's primary ENI.</li>
	// </li>`public-ip-address` - String - Optional - Filter by the public IP address of the instance's primary ENI, including the IP addresses automatically assigned during the instance creation and the EIPs manually associated after the instance creation.</li>
	// <li> `tag-key` - String - Optional - Filter by the tag key.</li>
	// </li>`tag-value` - String - Optional - Filter by the tag value.</li>
	// <li> `tag:tag-key` - String - Optional - Filter by the tag key-value pair. Replace `tag-key` with the actual tag keys. See example 2.</li>
	// Each request can have up to 10 `Filters` and 5 `Filters.Values`. You cannot specify `InstanceIds` and `Filters` at the same time.
	Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`

	// Offset; default value: 0. For more information on `Offset`, see the corresponding section in API [Introduction](https://intl.cloud.tencent.com/document/product/377).
	Offset *int64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// Number of results returned; default value: 20; maximum: 100. For more information on `Limit`, see the corresponding section in API [Introduction](https://intl.cloud.tencent.com/document/product/377). 
	Limit *int64 `json:"Limit,omitnil,omitempty" name:"Limit"`
}

type DescribeInstancesRequest struct {
	*tchttp.BaseRequest
	
	// Query by instance ID(s). For example, instance ID: `ins-xxxxxxxx`. For the specific format, refer to section `Ids.N` of the API [Introduction](https://intl.cloud.tencent.com/document/api/213/15688?from_cn_redirect=1). You can query up to 100 instances in each request. However, `InstanceIds` and `Filters` cannot be specified at the same time.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// Filters
	// <li> `zone` - String - Optional - Filter by the availability zone.</li>
	// <li> `project-id` - Integer - Optional - Filter by the project ID. You can query the list of created projects by calling `DescribeProject` or logging in to the [CVM console](https://console.cloud.tencent.com/cvm/index). You can also call `AddProject` to create projects. </li>
	// <li> `host-id` - String - Optional - Filter by the CDH instance ID. Format: `host-xxxxxxxx`.</li>
	// </li>`vpc-id` - String - Optional - Filter by the VPC ID. Format: `vpc-xxxxxxxx`.</li>
	// <li> `subnet-id` - String - Optional - Filter by the subnet ID. Format: `subnet-xxxxxxxx`.</li>
	// </li>`instance-id` - String - Optional - Filter by the instance ID. Format: `ins-xxxxxxxx`.</li>
	// </li>`security-group-id` - String - Optional - Filter by the security group ID. Format: `sg-8jlk3f3r`.</li>
	// </li>`instance-name` - String - Optional - Filter by the instance name.</li>
	// </li>`instance-charge-type` - String - Optional - Filter by the instance billing method. `POSTPAID_BY_HOUR`: pay-as-you-go | `CDHPAID`: You are only billed for [CDH](https://intl.cloud.tencent.com/document/product/416?from_cn_redirect=1) instances, not the CVMs running on the [CDH](https://intl.cloud.tencent.com/document/product/416?from_cn_redirect=1) instances.</li>
	// </li>`private-ip-address` - String - Optional - Filter by the private IP address of the instance's primary ENI.</li>
	// </li>`public-ip-address` - String - Optional - Filter by the public IP address of the instance's primary ENI, including the IP addresses automatically assigned during the instance creation and the EIPs manually associated after the instance creation.</li>
	// <li> `tag-key` - String - Optional - Filter by the tag key.</li>
	// </li>`tag-value` - String - Optional - Filter by the tag value.</li>
	// <li> `tag:tag-key` - String - Optional - Filter by the tag key-value pair. Replace `tag-key` with the actual tag keys. See example 2.</li>
	// Each request can have up to 10 `Filters` and 5 `Filters.Values`. You cannot specify `InstanceIds` and `Filters` at the same time.
	Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`

	// Offset; default value: 0. For more information on `Offset`, see the corresponding section in API [Introduction](https://intl.cloud.tencent.com/document/product/377).
	Offset *int64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// Number of results returned; default value: 20; maximum: 100. For more information on `Limit`, see the corresponding section in API [Introduction](https://intl.cloud.tencent.com/document/product/377). 
	Limit *int64 `json:"Limit,omitnil,omitempty" name:"Limit"`
}

func (r *DescribeInstancesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeInstancesRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceIds")
	delete(f, "Filters")
	delete(f, "Offset")
	delete(f, "Limit")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeInstancesRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeInstancesResponseParams struct {
	// Number of instances meeting the filtering conditions.
	TotalCount *int64 `json:"TotalCount,omitnil,omitempty" name:"TotalCount"`

	// Detailed instance information.
	InstanceSet []*Instance `json:"InstanceSet,omitnil,omitempty" name:"InstanceSet"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeInstancesResponse struct {
	*tchttp.BaseResponse
	Response *DescribeInstancesResponseParams `json:"Response"`
}

func (r *DescribeInstancesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeInstancesResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeInstancesStatusRequestParams struct {
	// Query by instance ID(s). For example, instance ID: `ins-xxxxxxxx`. For the specific format, refer to section `Ids.N` of the API [Introduction](https://intl.cloud.tencent.com/document/api/213/15688?from_cn_redirect=1). You can query up to 100 instances in each request.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// Offset; default value: 0. For more information on `Offset`, see the corresponding section in API [Introduction](https://intl.cloud.tencent.com/document/product/377).
	Offset *int64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// Number of results returned; default value: 20; maximum: 100. For more information on `Limit`, see the corresponding section in API [Introduction](https://intl.cloud.tencent.com/document/product/377).
	Limit *int64 `json:"Limit,omitnil,omitempty" name:"Limit"`
}

type DescribeInstancesStatusRequest struct {
	*tchttp.BaseRequest
	
	// Query by instance ID(s). For example, instance ID: `ins-xxxxxxxx`. For the specific format, refer to section `Ids.N` of the API [Introduction](https://intl.cloud.tencent.com/document/api/213/15688?from_cn_redirect=1). You can query up to 100 instances in each request.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// Offset; default value: 0. For more information on `Offset`, see the corresponding section in API [Introduction](https://intl.cloud.tencent.com/document/product/377).
	Offset *int64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// Number of results returned; default value: 20; maximum: 100. For more information on `Limit`, see the corresponding section in API [Introduction](https://intl.cloud.tencent.com/document/product/377).
	Limit *int64 `json:"Limit,omitnil,omitempty" name:"Limit"`
}

func (r *DescribeInstancesStatusRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeInstancesStatusRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceIds")
	delete(f, "Offset")
	delete(f, "Limit")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeInstancesStatusRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeInstancesStatusResponseParams struct {
	// Number of instance states meeting the filtering conditions.
	TotalCount *int64 `json:"TotalCount,omitnil,omitempty" name:"TotalCount"`

	// [Instance status](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1) list.
	InstanceStatusSet []*InstanceStatus `json:"InstanceStatusSet,omitnil,omitempty" name:"InstanceStatusSet"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeInstancesStatusResponse struct {
	*tchttp.BaseResponse
	Response *DescribeInstancesStatusResponseParams `json:"Response"`
}

func (r *DescribeInstancesStatusResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeInstancesStatusResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeInternetChargeTypeConfigsRequestParams struct {

}

type DescribeInternetChargeTypeConfigsRequest struct {
	*tchttp.BaseRequest
	
}

func (r *DescribeInternetChargeTypeConfigsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeInternetChargeTypeConfigsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeInternetChargeTypeConfigsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeInternetChargeTypeConfigsResponseParams struct {
	// List of network billing methods.
	InternetChargeTypeConfigSet []*InternetChargeTypeConfig `json:"InternetChargeTypeConfigSet,omitnil,omitempty" name:"InternetChargeTypeConfigSet"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeInternetChargeTypeConfigsResponse struct {
	*tchttp.BaseResponse
	Response *DescribeInternetChargeTypeConfigsResponseParams `json:"Response"`
}

func (r *DescribeInternetChargeTypeConfigsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeInternetChargeTypeConfigsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeKeyPairsRequestParams struct {
	// Key pair ID(s) in the format of `skey-11112222`. This API supports using multiple IDs as filters at the same time. For more information on the format of this parameter, see the `id.N` section in [API Introduction](https://intl.cloud.tencent.com/document/api/213/15688?from_cn_redirect=1). You cannot specify `KeyIds` and `Filters` at the same time. You can log in to the [console](https://console.cloud.tencent.com/cvm/index) to query the key pair IDs.
	KeyIds []*string `json:"KeyIds,omitnil,omitempty" name:"KeyIds"`

	// Filters
	// <li> `project-id` - Integer - Optional - Filter by project ID. To view the list of project IDs, you can go to [Project Management](https://console.cloud.tencent.com/project), or call the `DescribeProject` API and look for `projectId` in the response.</li>
	// <li> `key-name` - String - Optional - Filter by key pair name.</li>You cannot specify `KeyIds` and `Filters` at the same time.
	Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`

	// Offset. The default value is `0`. For more information on `Offset`, see the relevant section in API [Introduction](https://intl.cloud.tencent.com/document/api/213/15688?from_cn_redirect=1).
	Offset *int64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// Number of results returned; default value: 20; maximum: 100. For more information on `Limit`, see the corresponding section in API [Introduction](https://intl.cloud.tencent.com/document/product/377). 
	Limit *int64 `json:"Limit,omitnil,omitempty" name:"Limit"`
}

type DescribeKeyPairsRequest struct {
	*tchttp.BaseRequest
	
	// Key pair ID(s) in the format of `skey-11112222`. This API supports using multiple IDs as filters at the same time. For more information on the format of this parameter, see the `id.N` section in [API Introduction](https://intl.cloud.tencent.com/document/api/213/15688?from_cn_redirect=1). You cannot specify `KeyIds` and `Filters` at the same time. You can log in to the [console](https://console.cloud.tencent.com/cvm/index) to query the key pair IDs.
	KeyIds []*string `json:"KeyIds,omitnil,omitempty" name:"KeyIds"`

	// Filters
	// <li> `project-id` - Integer - Optional - Filter by project ID. To view the list of project IDs, you can go to [Project Management](https://console.cloud.tencent.com/project), or call the `DescribeProject` API and look for `projectId` in the response.</li>
	// <li> `key-name` - String - Optional - Filter by key pair name.</li>You cannot specify `KeyIds` and `Filters` at the same time.
	Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`

	// Offset. The default value is `0`. For more information on `Offset`, see the relevant section in API [Introduction](https://intl.cloud.tencent.com/document/api/213/15688?from_cn_redirect=1).
	Offset *int64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// Number of results returned; default value: 20; maximum: 100. For more information on `Limit`, see the corresponding section in API [Introduction](https://intl.cloud.tencent.com/document/product/377). 
	Limit *int64 `json:"Limit,omitnil,omitempty" name:"Limit"`
}

func (r *DescribeKeyPairsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeKeyPairsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "KeyIds")
	delete(f, "Filters")
	delete(f, "Offset")
	delete(f, "Limit")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeKeyPairsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeKeyPairsResponseParams struct {
	// Number of key pairs meeting the filtering conditions.
	TotalCount *int64 `json:"TotalCount,omitnil,omitempty" name:"TotalCount"`

	// Detailed information on key pairs.
	KeyPairSet []*KeyPair `json:"KeyPairSet,omitnil,omitempty" name:"KeyPairSet"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeKeyPairsResponse struct {
	*tchttp.BaseResponse
	Response *DescribeKeyPairsResponseParams `json:"Response"`
}

func (r *DescribeKeyPairsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeKeyPairsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeLaunchTemplateVersionsRequestParams struct {
	// The launch template ID
	LaunchTemplateId *string `json:"LaunchTemplateId,omitnil,omitempty" name:"LaunchTemplateId"`

	// List of instance launch templates.
	LaunchTemplateVersions []*uint64 `json:"LaunchTemplateVersions,omitnil,omitempty" name:"LaunchTemplateVersions"`

	// The minimum version number specified, which defaults to 0.
	MinVersion *uint64 `json:"MinVersion,omitnil,omitempty" name:"MinVersion"`

	// The maximum version number specified, which defaults to 30.
	MaxVersion *uint64 `json:"MaxVersion,omitnil,omitempty" name:"MaxVersion"`

	// The offset. Default value: 0. For more information on `Offset`, see the relevant sections in API [Introduction](https://intl.cloud.tencent.com/document/api/213/15688?from_cn_redirect=1).
	Offset *uint64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// The number of returned results. Default value: 20. Maximum value: 100. For more information on `Limit`, see the relevant sections in API [Introduction](https://intl.cloud.tencent.com/document/api/213/15688?from_cn_redirect=1).
	Limit *uint64 `json:"Limit,omitnil,omitempty" name:"Limit"`

	// Specify whether to query the default version. This parameter and `LaunchTemplateVersions` cannot be specified at the same time.
	DefaultVersion *bool `json:"DefaultVersion,omitnil,omitempty" name:"DefaultVersion"`
}

type DescribeLaunchTemplateVersionsRequest struct {
	*tchttp.BaseRequest
	
	// The launch template ID
	LaunchTemplateId *string `json:"LaunchTemplateId,omitnil,omitempty" name:"LaunchTemplateId"`

	// List of instance launch templates.
	LaunchTemplateVersions []*uint64 `json:"LaunchTemplateVersions,omitnil,omitempty" name:"LaunchTemplateVersions"`

	// The minimum version number specified, which defaults to 0.
	MinVersion *uint64 `json:"MinVersion,omitnil,omitempty" name:"MinVersion"`

	// The maximum version number specified, which defaults to 30.
	MaxVersion *uint64 `json:"MaxVersion,omitnil,omitempty" name:"MaxVersion"`

	// The offset. Default value: 0. For more information on `Offset`, see the relevant sections in API [Introduction](https://intl.cloud.tencent.com/document/api/213/15688?from_cn_redirect=1).
	Offset *uint64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// The number of returned results. Default value: 20. Maximum value: 100. For more information on `Limit`, see the relevant sections in API [Introduction](https://intl.cloud.tencent.com/document/api/213/15688?from_cn_redirect=1).
	Limit *uint64 `json:"Limit,omitnil,omitempty" name:"Limit"`

	// Specify whether to query the default version. This parameter and `LaunchTemplateVersions` cannot be specified at the same time.
	DefaultVersion *bool `json:"DefaultVersion,omitnil,omitempty" name:"DefaultVersion"`
}

func (r *DescribeLaunchTemplateVersionsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeLaunchTemplateVersionsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "LaunchTemplateId")
	delete(f, "LaunchTemplateVersions")
	delete(f, "MinVersion")
	delete(f, "MaxVersion")
	delete(f, "Offset")
	delete(f, "Limit")
	delete(f, "DefaultVersion")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeLaunchTemplateVersionsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeLaunchTemplateVersionsResponseParams struct {
	// Total number of instance launch templates.
	TotalCount *uint64 `json:"TotalCount,omitnil,omitempty" name:"TotalCount"`

	// Set of instance launch template versions.
	LaunchTemplateVersionSet []*LaunchTemplateVersionInfo `json:"LaunchTemplateVersionSet,omitnil,omitempty" name:"LaunchTemplateVersionSet"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeLaunchTemplateVersionsResponse struct {
	*tchttp.BaseResponse
	Response *DescribeLaunchTemplateVersionsResponseParams `json:"Response"`
}

func (r *DescribeLaunchTemplateVersionsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeLaunchTemplateVersionsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeLaunchTemplatesRequestParams struct {
	// Instance launch template ID. ID of one or more instance launch templates. If not specified, all templates of the user will be displayed.
	LaunchTemplateIds []*string `json:"LaunchTemplateIds,omitnil,omitempty" name:"LaunchTemplateIds"`

	// <p style="padding-left: 30px;">Filter by [<strong>LaunchTemplateName</strong>].</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Optional</p>
	// The maximum number of `Filters` in each request is 10. The upper limit for `Filter.Values` is 5. This parameter cannot specify `LaunchTemplateIds` and `Filters` at the same time.
	Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`

	// The offset. Default value: 0. For more information on `Offset`, see the relevant sections in API [Introduction](https://intl.cloud.tencent.com/document/api/213/15688?from_cn_redirect=1).
	Offset *int64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// The number of returned results. Default value: 20. Maximum value: 100. For more information on `Limit`, see the relevant sections in API [Introduction](https://intl.cloud.tencent.com/document/api/213/15688?from_cn_redirect=1).
	Limit *int64 `json:"Limit,omitnil,omitempty" name:"Limit"`
}

type DescribeLaunchTemplatesRequest struct {
	*tchttp.BaseRequest
	
	// Instance launch template ID. ID of one or more instance launch templates. If not specified, all templates of the user will be displayed.
	LaunchTemplateIds []*string `json:"LaunchTemplateIds,omitnil,omitempty" name:"LaunchTemplateIds"`

	// <p style="padding-left: 30px;">Filter by [<strong>LaunchTemplateName</strong>].</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Optional</p>
	// The maximum number of `Filters` in each request is 10. The upper limit for `Filter.Values` is 5. This parameter cannot specify `LaunchTemplateIds` and `Filters` at the same time.
	Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`

	// The offset. Default value: 0. For more information on `Offset`, see the relevant sections in API [Introduction](https://intl.cloud.tencent.com/document/api/213/15688?from_cn_redirect=1).
	Offset *int64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// The number of returned results. Default value: 20. Maximum value: 100. For more information on `Limit`, see the relevant sections in API [Introduction](https://intl.cloud.tencent.com/document/api/213/15688?from_cn_redirect=1).
	Limit *int64 `json:"Limit,omitnil,omitempty" name:"Limit"`
}

func (r *DescribeLaunchTemplatesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeLaunchTemplatesRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "LaunchTemplateIds")
	delete(f, "Filters")
	delete(f, "Offset")
	delete(f, "Limit")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeLaunchTemplatesRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeLaunchTemplatesResponseParams struct {
	// Number of eligible instance templates.
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	TotalCount *int64 `json:"TotalCount,omitnil,omitempty" name:"TotalCount"`

	// List of instance details
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	LaunchTemplateSet []*LaunchTemplateInfo `json:"LaunchTemplateSet,omitnil,omitempty" name:"LaunchTemplateSet"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeLaunchTemplatesResponse struct {
	*tchttp.BaseResponse
	Response *DescribeLaunchTemplatesResponseParams `json:"Response"`
}

func (r *DescribeLaunchTemplatesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeLaunchTemplatesResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeRegionsRequestParams struct {

}

type DescribeRegionsRequest struct {
	*tchttp.BaseRequest
	
}

func (r *DescribeRegionsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeRegionsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeRegionsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeRegionsResponseParams struct {
	// Number of regions
	TotalCount *uint64 `json:"TotalCount,omitnil,omitempty" name:"TotalCount"`

	// List of regions
	RegionSet []*RegionInfo `json:"RegionSet,omitnil,omitempty" name:"RegionSet"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeRegionsResponse struct {
	*tchttp.BaseResponse
	Response *DescribeRegionsResponseParams `json:"Response"`
}

func (r *DescribeRegionsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeRegionsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeReservedInstancesConfigInfosRequestParams struct {
	// <li><strong>zone</li></strong>
	// Filters by the availability zones in which the reserved instance can be purchased, such as `ap-guangzhou-1`.
	// Type: String
	// Required: no
	// Valid values: list of regions/availability zones
	// 
	// <li><strong>product-description</li></strong>
	// Filters by the platform description (operating system) of the reserved instance, such as `linux`.
	// Type: String
	// Required: no
	// Valid value: linux
	// 
	// <li><strong>duration</li></strong>
	// Filters by the **validity** of the reserved instance, which is the purchased usage period. For example, `31536000`.
	// Type: Integer
	// Unit: second
	// Required: no
	// Valid value: 31536000 (1 year)
	Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`
}

type DescribeReservedInstancesConfigInfosRequest struct {
	*tchttp.BaseRequest
	
	// <li><strong>zone</li></strong>
	// Filters by the availability zones in which the reserved instance can be purchased, such as `ap-guangzhou-1`.
	// Type: String
	// Required: no
	// Valid values: list of regions/availability zones
	// 
	// <li><strong>product-description</li></strong>
	// Filters by the platform description (operating system) of the reserved instance, such as `linux`.
	// Type: String
	// Required: no
	// Valid value: linux
	// 
	// <li><strong>duration</li></strong>
	// Filters by the **validity** of the reserved instance, which is the purchased usage period. For example, `31536000`.
	// Type: Integer
	// Unit: second
	// Required: no
	// Valid value: 31536000 (1 year)
	Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`
}

func (r *DescribeReservedInstancesConfigInfosRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeReservedInstancesConfigInfosRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Filters")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeReservedInstancesConfigInfosRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeReservedInstancesConfigInfosResponseParams struct {
	// Static configurations of the reserved instance.
	ReservedInstanceConfigInfos []*ReservedInstanceConfigInfoItem `json:"ReservedInstanceConfigInfos,omitnil,omitempty" name:"ReservedInstanceConfigInfos"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeReservedInstancesConfigInfosResponse struct {
	*tchttp.BaseResponse
	Response *DescribeReservedInstancesConfigInfosResponseParams `json:"Response"`
}

func (r *DescribeReservedInstancesConfigInfosResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeReservedInstancesConfigInfosResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeReservedInstancesOfferingsRequestParams struct {
	// Dry run. Default value: false.
	DryRun *bool `json:"DryRun,omitnil,omitempty" name:"DryRun"`

	// The offset. Default value: 0. For more information on `Offset`, see the relevant sections in API [Introduction](https://intl.cloud.tencent.com/document/api/213/15688?from_cn_redirect=1).
	Offset *int64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// The number of returned results. Default value: 20. Maximum value: 100. For more information on `Limit`, see the relevant sections in API [Introduction](https://intl.cloud.tencent.com/document/api/213/15688?from_cn_redirect=1).
	Limit *int64 `json:"Limit,omitnil,omitempty" name:"Limit"`

	// The maximum duration as a filter, 
	// in seconds.
	// Default value: 94608000.
	MaxDuration *int64 `json:"MaxDuration,omitnil,omitempty" name:"MaxDuration"`

	// The minimum duration as a filter, 
	// in seconds.
	// Default value: 2592000.
	MinDuration *int64 `json:"MinDuration,omitnil,omitempty" name:"MinDuration"`

	// <li><strong>zone</strong></li>
	// <p style="padding-left: 30px;">Filters by the <strong>availability zones</strong> in which the Reserved Instances can be purchased, such as ap-guangzhou-1.</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Required: no</p><p style="padding-left: 30px;">Valid values: please see <a href="https://intl.cloud.tencent.com/document/product/213/6091?from_cn_redirect=1">Availability Zones</a></p>
	// <li><strong>duration</strong></li>
	// <p style="padding-left: 30px;">Filters by the <strong>duration</strong> of the Reserved Instance, in seconds. For example, 31536000.</p><p style="padding-left: 30px;">Type: Integer</p><p style="padding-left: 30px;">Unit: second</p><p style="padding-left: 30px;">Required: no</p><p style="padding-left: 30px;">Valid values: 31536000 (1 year) | 94608000 (3 years)</p>
	// <li><strong>instance-type</strong></li>
	// <p style="padding-left: 30px;">Filters by <strong>type of the Reserved Instance</strong>, such as `S3.MEDIUM4`.</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Required: no</p><p style="padding-left: 30px;">Valid values: please see <a href="https://intl.cloud.tencent.com/document/product/213/11518?from_cn_redirect=1">Instance Types</a></p>
	// <li><strong>offering-type</strong></li>
	// <p style="padding-left: 30px;">Filters by **<strong>payment term</strong>**, such as `All Upfront`.</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Required: no</p><p style="padding-left: 30px;">Valid value: All Upfront</p>
	// <li><strong>product-description</strong></li>
	// <p style="padding-left: 30px;">Filters by the <strong>platform description</strong> (operating system) of the Reserved Instance, such as `linux`.</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Required: no</p><p style="padding-left: 30px;">Valid value: linux</p>
	// <li><strong>reserved-instances-offering-id</strong></li>
	// <p style="padding-left: 30px;">Filters by <strong>Reserved Instance ID</strong>, in the form of 650c138f-ae7e-4750-952a-96841d6e9fc1.</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Required: no</p>
	// Each request can have up to 10 `Filters` and 5 `Filter.Values`.
	Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`
}

type DescribeReservedInstancesOfferingsRequest struct {
	*tchttp.BaseRequest
	
	// Dry run. Default value: false.
	DryRun *bool `json:"DryRun,omitnil,omitempty" name:"DryRun"`

	// The offset. Default value: 0. For more information on `Offset`, see the relevant sections in API [Introduction](https://intl.cloud.tencent.com/document/api/213/15688?from_cn_redirect=1).
	Offset *int64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// The number of returned results. Default value: 20. Maximum value: 100. For more information on `Limit`, see the relevant sections in API [Introduction](https://intl.cloud.tencent.com/document/api/213/15688?from_cn_redirect=1).
	Limit *int64 `json:"Limit,omitnil,omitempty" name:"Limit"`

	// The maximum duration as a filter, 
	// in seconds.
	// Default value: 94608000.
	MaxDuration *int64 `json:"MaxDuration,omitnil,omitempty" name:"MaxDuration"`

	// The minimum duration as a filter, 
	// in seconds.
	// Default value: 2592000.
	MinDuration *int64 `json:"MinDuration,omitnil,omitempty" name:"MinDuration"`

	// <li><strong>zone</strong></li>
	// <p style="padding-left: 30px;">Filters by the <strong>availability zones</strong> in which the Reserved Instances can be purchased, such as ap-guangzhou-1.</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Required: no</p><p style="padding-left: 30px;">Valid values: please see <a href="https://intl.cloud.tencent.com/document/product/213/6091?from_cn_redirect=1">Availability Zones</a></p>
	// <li><strong>duration</strong></li>
	// <p style="padding-left: 30px;">Filters by the <strong>duration</strong> of the Reserved Instance, in seconds. For example, 31536000.</p><p style="padding-left: 30px;">Type: Integer</p><p style="padding-left: 30px;">Unit: second</p><p style="padding-left: 30px;">Required: no</p><p style="padding-left: 30px;">Valid values: 31536000 (1 year) | 94608000 (3 years)</p>
	// <li><strong>instance-type</strong></li>
	// <p style="padding-left: 30px;">Filters by <strong>type of the Reserved Instance</strong>, such as `S3.MEDIUM4`.</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Required: no</p><p style="padding-left: 30px;">Valid values: please see <a href="https://intl.cloud.tencent.com/document/product/213/11518?from_cn_redirect=1">Instance Types</a></p>
	// <li><strong>offering-type</strong></li>
	// <p style="padding-left: 30px;">Filters by **<strong>payment term</strong>**, such as `All Upfront`.</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Required: no</p><p style="padding-left: 30px;">Valid value: All Upfront</p>
	// <li><strong>product-description</strong></li>
	// <p style="padding-left: 30px;">Filters by the <strong>platform description</strong> (operating system) of the Reserved Instance, such as `linux`.</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Required: no</p><p style="padding-left: 30px;">Valid value: linux</p>
	// <li><strong>reserved-instances-offering-id</strong></li>
	// <p style="padding-left: 30px;">Filters by <strong>Reserved Instance ID</strong>, in the form of 650c138f-ae7e-4750-952a-96841d6e9fc1.</p><p style="padding-left: 30px;">Type: String</p><p style="padding-left: 30px;">Required: no</p>
	// Each request can have up to 10 `Filters` and 5 `Filter.Values`.
	Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`
}

func (r *DescribeReservedInstancesOfferingsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeReservedInstancesOfferingsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "DryRun")
	delete(f, "Offset")
	delete(f, "Limit")
	delete(f, "MaxDuration")
	delete(f, "MinDuration")
	delete(f, "Filters")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeReservedInstancesOfferingsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeReservedInstancesOfferingsResponseParams struct {
	// The number of Reserved Instances that meet the condition.
	TotalCount *int64 `json:"TotalCount,omitnil,omitempty" name:"TotalCount"`

	// The list of Reserved Instances that meet the condition.
	ReservedInstancesOfferingsSet []*ReservedInstancesOffering `json:"ReservedInstancesOfferingsSet,omitnil,omitempty" name:"ReservedInstancesOfferingsSet"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeReservedInstancesOfferingsResponse struct {
	*tchttp.BaseResponse
	Response *DescribeReservedInstancesOfferingsResponseParams `json:"Response"`
}

func (r *DescribeReservedInstancesOfferingsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeReservedInstancesOfferingsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeReservedInstancesRequestParams struct {
	// Trial run. Default value: false.
	DryRun *bool `json:"DryRun,omitnil,omitempty" name:"DryRun"`

	// Offset. Default value: 0. For more information on `Offset`, see the relevant section in API [Introduction](https://intl.cloud.tencent.com/document/api/213/15688?from_cn_redirect=1).
	Offset *int64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// Number of returned instances. Default value: 20. Maximum value: 100. For more information on `Limit`, see the relevant section in API [Introduction](https://intl.cloud.tencent.com/document/api/213/15688?from_cn_redirect=1).
	Limit *int64 `json:"Limit,omitnil,omitempty" name:"Limit"`

	// <li><strong>zone</strong></li>
	// <p style="padding-left: 30px;">Filter by [<strong>availability zones</strong>] in which reserved instances can be purchased. For example, ap-guangzhou-1.</p><p style="padding-left: 30px;">Type: String.</p><p style="padding-left: 30px;">Required: No.</p><p style="padding-left: 30px;">Valid values: See the <a href="https://intl.cloud.tencent.com/document/product/213/6091?from_cn_redirect=1">availability zone list</a>.</p>
	// <li><strong>duration</strong></li>
	// <p style="padding-left: 30px;">Filter by [<strong>validity periods</strong>] of reserved instances, which is the instance purchase duration. For example, 31536000.</p><p style="padding-left: 30px;">Type: Integer.</p><p style="padding-left: 30px;">Unit: Second.</p><p style="padding-left: 30px;">Required: No.</p><p style="padding-left: 30px;">Valid values: 31536000 (1 year) | 94608000 (3 years).</p>
	// <li><strong>instance-type</strong></li>
	// <p style="padding-left: 30px;">Filter by [<strong>specifications of reserved instances</strong>]. For example, S3.MEDIUM4.</p><p style="padding-left: 30px;">Type: String.</p><p style="padding-left: 30px;">Required: No.</p><p style="padding-left: 30px;">Valid values: See the <a href="https://intl.cloud.tencent.com/document/product/213/11518?from_cn_redirect=1">reserved instance specification list</a>.</p>
	// <li><strong>instance-family</strong></li>
	// <p style="padding-left: 30px;">Filter by [<strong>types of reserved instances</strong>]. For example, S3.</p><p style="padding-left: 30px;">Type: String.</p><p style="padding-left: 30px;">Required: No.</p><p style="padding-left: 30px;">Valid values: See the <a href="https://intl.cloud.tencent.com/document/product/213/11518?from_cn_redirect=1">reserved instance type list</a>.</p>
	// <li><strong>offering-type</strong></li>
	// <p style="padding-left: 30px;">Filter by <strong>payment types</strong>]. For example, All Upfront (fully prepaid).</p><p style="padding-left: 30px;">Type: String.</p><p style="padding-left: 30px;">Required: No.</p><p style="padding-left: 30px;">Valid values: All Upfront (fully prepaid) | Partial Upfront (partially prepaid) | No Upfront (non-prepaid).</p>
	// <li><strong>product-description</strong></li>
	// <p style="padding-left: 30px;">Filter by [<strong>platform descriptions</strong>] (operating system) of reserved instances. For example, linux.</p><p style="padding-left: 30px;">Type: String.</p><p style="padding-left: 30px;">Required: No.</p><p style="padding-left: 30px;">Valid value: linux.</p>
	// <li><strong>reserved-instances-id</strong></li>
	// <p style="padding-left: 30px;">Filter by [<strong>IDs of purchased reserved instances</strong>]. For example, 650c138f-ae7e-4750-952a-96841d6e9fc1.</p><p style="padding-left: 30px;">Type: String.</p><p style="padding-left: 30px;">Required: No.</p>
	// <li><strong>state</strong></li>
	// <p style="padding-left: 30px;">Filter by [<strong>statuses of purchased reserved instances</strong>]. For example, active.</p><p style="padding-left: 30px;">Type: String.</p><p style="padding-left: 30px;">Required: No.</p><p style="padding-left: 30px;">Valid values: active (created) | pending (waiting to be created) | retired (expired).</p>
	// Each request can have up to 10 filters, and each filter can have up to 5 values.
	Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`
}

type DescribeReservedInstancesRequest struct {
	*tchttp.BaseRequest
	
	// Trial run. Default value: false.
	DryRun *bool `json:"DryRun,omitnil,omitempty" name:"DryRun"`

	// Offset. Default value: 0. For more information on `Offset`, see the relevant section in API [Introduction](https://intl.cloud.tencent.com/document/api/213/15688?from_cn_redirect=1).
	Offset *int64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// Number of returned instances. Default value: 20. Maximum value: 100. For more information on `Limit`, see the relevant section in API [Introduction](https://intl.cloud.tencent.com/document/api/213/15688?from_cn_redirect=1).
	Limit *int64 `json:"Limit,omitnil,omitempty" name:"Limit"`

	// <li><strong>zone</strong></li>
	// <p style="padding-left: 30px;">Filter by [<strong>availability zones</strong>] in which reserved instances can be purchased. For example, ap-guangzhou-1.</p><p style="padding-left: 30px;">Type: String.</p><p style="padding-left: 30px;">Required: No.</p><p style="padding-left: 30px;">Valid values: See the <a href="https://intl.cloud.tencent.com/document/product/213/6091?from_cn_redirect=1">availability zone list</a>.</p>
	// <li><strong>duration</strong></li>
	// <p style="padding-left: 30px;">Filter by [<strong>validity periods</strong>] of reserved instances, which is the instance purchase duration. For example, 31536000.</p><p style="padding-left: 30px;">Type: Integer.</p><p style="padding-left: 30px;">Unit: Second.</p><p style="padding-left: 30px;">Required: No.</p><p style="padding-left: 30px;">Valid values: 31536000 (1 year) | 94608000 (3 years).</p>
	// <li><strong>instance-type</strong></li>
	// <p style="padding-left: 30px;">Filter by [<strong>specifications of reserved instances</strong>]. For example, S3.MEDIUM4.</p><p style="padding-left: 30px;">Type: String.</p><p style="padding-left: 30px;">Required: No.</p><p style="padding-left: 30px;">Valid values: See the <a href="https://intl.cloud.tencent.com/document/product/213/11518?from_cn_redirect=1">reserved instance specification list</a>.</p>
	// <li><strong>instance-family</strong></li>
	// <p style="padding-left: 30px;">Filter by [<strong>types of reserved instances</strong>]. For example, S3.</p><p style="padding-left: 30px;">Type: String.</p><p style="padding-left: 30px;">Required: No.</p><p style="padding-left: 30px;">Valid values: See the <a href="https://intl.cloud.tencent.com/document/product/213/11518?from_cn_redirect=1">reserved instance type list</a>.</p>
	// <li><strong>offering-type</strong></li>
	// <p style="padding-left: 30px;">Filter by <strong>payment types</strong>]. For example, All Upfront (fully prepaid).</p><p style="padding-left: 30px;">Type: String.</p><p style="padding-left: 30px;">Required: No.</p><p style="padding-left: 30px;">Valid values: All Upfront (fully prepaid) | Partial Upfront (partially prepaid) | No Upfront (non-prepaid).</p>
	// <li><strong>product-description</strong></li>
	// <p style="padding-left: 30px;">Filter by [<strong>platform descriptions</strong>] (operating system) of reserved instances. For example, linux.</p><p style="padding-left: 30px;">Type: String.</p><p style="padding-left: 30px;">Required: No.</p><p style="padding-left: 30px;">Valid value: linux.</p>
	// <li><strong>reserved-instances-id</strong></li>
	// <p style="padding-left: 30px;">Filter by [<strong>IDs of purchased reserved instances</strong>]. For example, 650c138f-ae7e-4750-952a-96841d6e9fc1.</p><p style="padding-left: 30px;">Type: String.</p><p style="padding-left: 30px;">Required: No.</p>
	// <li><strong>state</strong></li>
	// <p style="padding-left: 30px;">Filter by [<strong>statuses of purchased reserved instances</strong>]. For example, active.</p><p style="padding-left: 30px;">Type: String.</p><p style="padding-left: 30px;">Required: No.</p><p style="padding-left: 30px;">Valid values: active (created) | pending (waiting to be created) | retired (expired).</p>
	// Each request can have up to 10 filters, and each filter can have up to 5 values.
	Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`
}

func (r *DescribeReservedInstancesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeReservedInstancesRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "DryRun")
	delete(f, "Offset")
	delete(f, "Limit")
	delete(f, "Filters")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeReservedInstancesRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeReservedInstancesResponseParams struct {
	// Number of reserved instances that meet the conditions.
	TotalCount *int64 `json:"TotalCount,omitnil,omitempty" name:"TotalCount"`

	// List of reserved instances that meet the conditions.
	ReservedInstancesSet []*ReservedInstances `json:"ReservedInstancesSet,omitnil,omitempty" name:"ReservedInstancesSet"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeReservedInstancesResponse struct {
	*tchttp.BaseResponse
	Response *DescribeReservedInstancesResponseParams `json:"Response"`
}

func (r *DescribeReservedInstancesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeReservedInstancesResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeZoneInstanceConfigInfosRequestParams struct {
	// <li> instance-charge-type-String-required: no-(filter) billing mode of instances. (POSTPAID_BY_HOUR: pay-as-you-go billing by hour | SPOTPAID: spot billing, which is suitable for a [spot instance] (https://intl.cloud.Tencent.com/document/product/213/17817) | CDHPAID: CDH billing, that is, billing only for CDH, but not for the instances on CDH. )  </li>
	Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`
}

type DescribeZoneInstanceConfigInfosRequest struct {
	*tchttp.BaseRequest
	
	// <li> instance-charge-type-String-required: no-(filter) billing mode of instances. (POSTPAID_BY_HOUR: pay-as-you-go billing by hour | SPOTPAID: spot billing, which is suitable for a [spot instance] (https://intl.cloud.Tencent.com/document/product/213/17817) | CDHPAID: CDH billing, that is, billing only for CDH, but not for the instances on CDH. )  </li>
	Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`
}

func (r *DescribeZoneInstanceConfigInfosRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeZoneInstanceConfigInfosRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Filters")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeZoneInstanceConfigInfosRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeZoneInstanceConfigInfosResponseParams struct {
	// List of model configurations for the availability zone.
	InstanceTypeQuotaSet []*InstanceTypeQuotaItem `json:"InstanceTypeQuotaSet,omitnil,omitempty" name:"InstanceTypeQuotaSet"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeZoneInstanceConfigInfosResponse struct {
	*tchttp.BaseResponse
	Response *DescribeZoneInstanceConfigInfosResponseParams `json:"Response"`
}

func (r *DescribeZoneInstanceConfigInfosResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeZoneInstanceConfigInfosResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeZonesRequestParams struct {

}

type DescribeZonesRequest struct {
	*tchttp.BaseRequest
	
}

func (r *DescribeZonesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeZonesRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeZonesRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeZonesResponseParams struct {
	// Number of availability zones.
	TotalCount *uint64 `json:"TotalCount,omitnil,omitempty" name:"TotalCount"`

	// List of availability zones.
	ZoneSet []*ZoneInfo `json:"ZoneSet,omitnil,omitempty" name:"ZoneSet"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeZonesResponse struct {
	*tchttp.BaseResponse
	Response *DescribeZonesResponseParams `json:"Response"`
}

func (r *DescribeZonesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeZonesResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DisassociateInstancesKeyPairsRequestParams struct {
	// Instance ID(s). The maximum number of instances in each request is 100. <br><br>You can obtain the available instance IDs in two ways: <br><li>Log in to the [console](https://console.cloud.tencent.com/cvm/index) to query the instance IDs. <br><li>Call [DescribeInstances](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1) and look for `InstanceId` in the response.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// List of key pair IDs. The maximum number of key pairs in each request is 100. The key pair ID is in the format of `skey-11112222`. <br><br>You can obtain the available key pair IDs in two ways: <br><li>Log in to the [console](https://console.cloud.tencent.com/cvm/sshkey) to query the key pair IDs. <br><li>Call [DescribeKeyPairs](https://intl.cloud.tencent.com/document/api/213/15699?from_cn_redirect=1) and look for `KeyId` in the response.
	KeyIds []*string `json:"KeyIds,omitnil,omitempty" name:"KeyIds"`

	// Whether to force shut down a running instances. It is recommended to manually shut down a running instance before disassociating a key pair from it. Valid values: <br><li>TRUE: force shut down an instance after a normal shutdown fails. <br><li>FALSE: do not force shut down an instance after a normal shutdown fails. <br><br>Default value: FALSE.
	ForceStop *bool `json:"ForceStop,omitnil,omitempty" name:"ForceStop"`
}

type DisassociateInstancesKeyPairsRequest struct {
	*tchttp.BaseRequest
	
	// Instance ID(s). The maximum number of instances in each request is 100. <br><br>You can obtain the available instance IDs in two ways: <br><li>Log in to the [console](https://console.cloud.tencent.com/cvm/index) to query the instance IDs. <br><li>Call [DescribeInstances](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1) and look for `InstanceId` in the response.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// List of key pair IDs. The maximum number of key pairs in each request is 100. The key pair ID is in the format of `skey-11112222`. <br><br>You can obtain the available key pair IDs in two ways: <br><li>Log in to the [console](https://console.cloud.tencent.com/cvm/sshkey) to query the key pair IDs. <br><li>Call [DescribeKeyPairs](https://intl.cloud.tencent.com/document/api/213/15699?from_cn_redirect=1) and look for `KeyId` in the response.
	KeyIds []*string `json:"KeyIds,omitnil,omitempty" name:"KeyIds"`

	// Whether to force shut down a running instances. It is recommended to manually shut down a running instance before disassociating a key pair from it. Valid values: <br><li>TRUE: force shut down an instance after a normal shutdown fails. <br><li>FALSE: do not force shut down an instance after a normal shutdown fails. <br><br>Default value: FALSE.
	ForceStop *bool `json:"ForceStop,omitnil,omitempty" name:"ForceStop"`
}

func (r *DisassociateInstancesKeyPairsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DisassociateInstancesKeyPairsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceIds")
	delete(f, "KeyIds")
	delete(f, "ForceStop")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DisassociateInstancesKeyPairsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DisassociateInstancesKeyPairsResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DisassociateInstancesKeyPairsResponse struct {
	*tchttp.BaseResponse
	Response *DisassociateInstancesKeyPairsResponseParams `json:"Response"`
}

func (r *DisassociateInstancesKeyPairsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DisassociateInstancesKeyPairsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DisassociateSecurityGroupsRequestParams struct {
	// ID of the security group to be disassociated, such as `sg-efil73jd`. Only one security group can be disassociated.
	SecurityGroupIds []*string `json:"SecurityGroupIds,omitnil,omitempty" name:"SecurityGroupIds"`

	// ID(s) of the instance(s) to be disassociated,such as `ins-lesecurk`. You can specify multiple instances.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`
}

type DisassociateSecurityGroupsRequest struct {
	*tchttp.BaseRequest
	
	// ID of the security group to be disassociated, such as `sg-efil73jd`. Only one security group can be disassociated.
	SecurityGroupIds []*string `json:"SecurityGroupIds,omitnil,omitempty" name:"SecurityGroupIds"`

	// ID(s) of the instance(s) to be disassociated,such as `ins-lesecurk`. You can specify multiple instances.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`
}

func (r *DisassociateSecurityGroupsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DisassociateSecurityGroupsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "SecurityGroupIds")
	delete(f, "InstanceIds")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DisassociateSecurityGroupsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DisassociateSecurityGroupsResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DisassociateSecurityGroupsResponse struct {
	*tchttp.BaseResponse
	Response *DisassociateSecurityGroupsResponseParams `json:"Response"`
}

func (r *DisassociateSecurityGroupsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DisassociateSecurityGroupsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type DisasterRecoverGroup struct {
	// ID of a spread placement group.
	DisasterRecoverGroupId *string `json:"DisasterRecoverGroupId,omitnil,omitempty" name:"DisasterRecoverGroupId"`

	// Name of a spread placement group, which must be 1-60 characters long.
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// Type of a spread placement group. Valid values:<br>
	// <li>HOST: physical machine.<br></li>
	// <li>SW: switch.<br></li>
	// <li>RACK: rack.</li>
	Type *string `json:"Type,omitnil,omitempty" name:"Type"`

	// The maximum number of CVMs that can be hosted in a spread placement group.
	CvmQuotaTotal *int64 `json:"CvmQuotaTotal,omitnil,omitempty" name:"CvmQuotaTotal"`

	// The current number of CVMs in a spread placement group.
	CurrentNum *int64 `json:"CurrentNum,omitnil,omitempty" name:"CurrentNum"`

	// The list of CVM IDs in a spread placement group.
	// Note: This field may return null, indicating that no valid value was found.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// Creation time of a spread placement group.
	// Note: This field may return null, indicating that no valid value is found.
	CreateTime *string `json:"CreateTime,omitnil,omitempty" name:"CreateTime"`


	Affinity *int64 `json:"Affinity,omitnil,omitempty" name:"Affinity"`

	// List of tags associated with the placement group.
	Tags []*Tag `json:"Tags,omitnil,omitempty" name:"Tags"`
}

type EnhancedService struct {
	// Enables cloud security service. If this parameter is not specified, the cloud security service will be enabled by default.
	SecurityService *RunSecurityServiceEnabled `json:"SecurityService,omitnil,omitempty" name:"SecurityService"`

	// Enables cloud monitor service. If this parameter is not specified, the cloud monitor service will be enabled by default.
	MonitorService *RunMonitorServiceEnabled `json:"MonitorService,omitnil,omitempty" name:"MonitorService"`

	// Whether to enable the TAT service. If this parameter is not specified, the TAT service is enabled for public images and disabled for other images by default.
	AutomationService *RunAutomationServiceEnabled `json:"AutomationService,omitnil,omitempty" name:"AutomationService"`
}

// Predefined struct for user
type EnterRescueModeRequestParams struct {
	// Instance ID Needs to Enter Rescue Mode
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// System Password in Rescue Mode
	Password *string `json:"Password,omitnil,omitempty" name:"Password"`

	// System Username in Rescue Mode
	Username *string `json:"Username,omitnil,omitempty" name:"Username"`

	// Whether to perform forced shutdown.
	//
	// Deprecated: ForceStop is deprecated.
	ForceStop *bool `json:"ForceStop,omitnil,omitempty" name:"ForceStop"`


	StopType *string `json:"StopType,omitnil,omitempty" name:"StopType"`
}

type EnterRescueModeRequest struct {
	*tchttp.BaseRequest
	
	// Instance ID Needs to Enter Rescue Mode
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// System Password in Rescue Mode
	Password *string `json:"Password,omitnil,omitempty" name:"Password"`

	// System Username in Rescue Mode
	Username *string `json:"Username,omitnil,omitempty" name:"Username"`

	// Whether to perform forced shutdown.
	ForceStop *bool `json:"ForceStop,omitnil,omitempty" name:"ForceStop"`

	StopType *string `json:"StopType,omitnil,omitempty" name:"StopType"`
}

func (r *EnterRescueModeRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *EnterRescueModeRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceId")
	delete(f, "Password")
	delete(f, "Username")
	delete(f, "ForceStop")
	delete(f, "StopType")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "EnterRescueModeRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type EnterRescueModeResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type EnterRescueModeResponse struct {
	*tchttp.BaseResponse
	Response *EnterRescueModeResponseParams `json:"Response"`
}

func (r *EnterRescueModeResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *EnterRescueModeResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ExitRescueModeRequestParams struct {
	// Instance ID Exiting Rescue Mode
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`
}

type ExitRescueModeRequest struct {
	*tchttp.BaseRequest
	
	// Instance ID Exiting Rescue Mode
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`
}

func (r *ExitRescueModeRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ExitRescueModeRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ExitRescueModeRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ExitRescueModeResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ExitRescueModeResponse struct {
	*tchttp.BaseResponse
	Response *ExitRescueModeResponseParams `json:"Response"`
}

func (r *ExitRescueModeResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ExitRescueModeResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ExportImagesRequestParams struct {
	// COS bucket name
	BucketName *string `json:"BucketName,omitnil,omitempty" name:"BucketName"`

	// List of image IDs
	ImageIds []*string `json:"ImageIds,omitnil,omitempty" name:"ImageIds"`

	// Format of the exported image file. Valid values: `RAW`, `QCOW2`, `VHD` and `VMDK`. Default value: `RAW`.
	ExportFormat *string `json:"ExportFormat,omitnil,omitempty" name:"ExportFormat"`

	// Prefix list of the name of exported files
	FileNamePrefixList []*string `json:"FileNamePrefixList,omitnil,omitempty" name:"FileNamePrefixList"`

	// Whether to export only the system disk
	OnlyExportRootDisk *bool `json:"OnlyExportRootDisk,omitnil,omitempty" name:"OnlyExportRootDisk"`

	// Check whether the image can be exported
	DryRun *bool `json:"DryRun,omitnil,omitempty" name:"DryRun"`

	// Role name (Default: `CVM_QcsRole`). Before exporting the images, make sure the role exists, and it has write permission to COS.
	RoleName *string `json:"RoleName,omitnil,omitempty" name:"RoleName"`
}

type ExportImagesRequest struct {
	*tchttp.BaseRequest
	
	// COS bucket name
	BucketName *string `json:"BucketName,omitnil,omitempty" name:"BucketName"`

	// List of image IDs
	ImageIds []*string `json:"ImageIds,omitnil,omitempty" name:"ImageIds"`

	// Format of the exported image file. Valid values: `RAW`, `QCOW2`, `VHD` and `VMDK`. Default value: `RAW`.
	ExportFormat *string `json:"ExportFormat,omitnil,omitempty" name:"ExportFormat"`

	// Prefix list of the name of exported files
	FileNamePrefixList []*string `json:"FileNamePrefixList,omitnil,omitempty" name:"FileNamePrefixList"`

	// Whether to export only the system disk
	OnlyExportRootDisk *bool `json:"OnlyExportRootDisk,omitnil,omitempty" name:"OnlyExportRootDisk"`

	// Check whether the image can be exported
	DryRun *bool `json:"DryRun,omitnil,omitempty" name:"DryRun"`

	// Role name (Default: `CVM_QcsRole`). Before exporting the images, make sure the role exists, and it has write permission to COS.
	RoleName *string `json:"RoleName,omitnil,omitempty" name:"RoleName"`
}

func (r *ExportImagesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ExportImagesRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "BucketName")
	delete(f, "ImageIds")
	delete(f, "ExportFormat")
	delete(f, "FileNamePrefixList")
	delete(f, "OnlyExportRootDisk")
	delete(f, "DryRun")
	delete(f, "RoleName")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ExportImagesRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ExportImagesResponseParams struct {
	// ID of the image export task
	TaskId *uint64 `json:"TaskId,omitnil,omitempty" name:"TaskId"`

	// List of COS filenames of the exported images
	CosPaths []*string `json:"CosPaths,omitnil,omitempty" name:"CosPaths"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ExportImagesResponse struct {
	*tchttp.BaseResponse
	Response *ExportImagesResponseParams `json:"Response"`
}

func (r *ExportImagesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ExportImagesResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type Externals struct {
	// Release Address
	ReleaseAddress *bool `json:"ReleaseAddress,omitnil,omitempty" name:"ReleaseAddress"`

	// Unsupported network type. valid values: <br><li>BASIC: BASIC network</li><li>VPC1.0: private network VPC1.0</li>.
	UnsupportNetworks []*string `json:"UnsupportNetworks,omitnil,omitempty" name:"UnsupportNetworks"`

	// Specifies the HDD local storage attributes.
	StorageBlockAttr *StorageBlock `json:"StorageBlockAttr,omitnil,omitempty" name:"StorageBlockAttr"`
}

type Filter struct {
	// Filters.
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// Filter values.
	Values []*string `json:"Values,omitnil,omitempty" name:"Values"`
}

type GPUInfo struct {
	// Specifies the GPU count of the instance. a value less than 1 indicates VGPU type, and a value larger than 1 indicates GPU passthrough type.
	GPUCount *float64 `json:"GPUCount,omitnil,omitempty" name:"GPUCount"`

	// Specifies the GPU address of the instance.
	GPUId []*string `json:"GPUId,omitnil,omitempty" name:"GPUId"`

	// Specifies the GPU type of the instance.
	GPUType *string `json:"GPUType,omitnil,omitempty" name:"GPUType"`
}

type HostItem struct {
	// CDH instance location. This parameter is used to specify the AZ, project, and other attributes of the instance.
	Placement *Placement `json:"Placement,omitnil,omitempty" name:"Placement"`

	// CDH instance ID
	HostId *string `json:"HostId,omitnil,omitempty" name:"HostId"`

	// CDH instance type
	HostType *string `json:"HostType,omitnil,omitempty" name:"HostType"`

	// CDH instance name
	HostName *string `json:"HostName,omitnil,omitempty" name:"HostName"`

	// CDH instance billing mode
	HostChargeType *string `json:"HostChargeType,omitnil,omitempty" name:"HostChargeType"`

	// CDH instance renewal flag
	RenewFlag *string `json:"RenewFlag,omitnil,omitempty" name:"RenewFlag"`

	// CDH instance creation time
	CreatedTime *string `json:"CreatedTime,omitnil,omitempty" name:"CreatedTime"`

	// CDH instance expiry time
	ExpiredTime *string `json:"ExpiredTime,omitnil,omitempty" name:"ExpiredTime"`

	// List of IDs of CVMs created on a CDH instance
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// CDH instance status
	HostState *string `json:"HostState,omitnil,omitempty" name:"HostState"`

	// CDH instance IP
	HostIp *string `json:"HostIp,omitnil,omitempty" name:"HostIp"`

	// CDH instance resource information
	HostResource *HostResource `json:"HostResource,omitnil,omitempty" name:"HostResource"`

	// Cage ID of the CDH instance. This parameter is only valid for CDH instances in the cages of finance availability zones.
	// Note: This field may return null, indicating that no valid value is found.
	CageId *string `json:"CageId,omitnil,omitempty" name:"CageId"`
}

type HostResource struct {
	// Total number of CPU cores in the CDH instance
	CpuTotal *uint64 `json:"CpuTotal,omitnil,omitempty" name:"CpuTotal"`

	// Number of available CPU cores in the CDH instance
	CpuAvailable *uint64 `json:"CpuAvailable,omitnil,omitempty" name:"CpuAvailable"`

	// Total memory size of the CDH instance (unit: GiB)
	MemTotal *float64 `json:"MemTotal,omitnil,omitempty" name:"MemTotal"`

	// Available memory size of the CDH instance (unit: GiB)
	MemAvailable *float64 `json:"MemAvailable,omitnil,omitempty" name:"MemAvailable"`

	// Total disk size of the CDH instance (unit: GiB)
	DiskTotal *uint64 `json:"DiskTotal,omitnil,omitempty" name:"DiskTotal"`

	// Available disk size of the CDH instance (unit: GiB)
	DiskAvailable *uint64 `json:"DiskAvailable,omitnil,omitempty" name:"DiskAvailable"`

	// Disk type of the CDH instance
	DiskType *string `json:"DiskType,omitnil,omitempty" name:"DiskType"`

	// Total number of GPU cards in the CDH instance
	GpuTotal *uint64 `json:"GpuTotal,omitnil,omitempty" name:"GpuTotal"`

	// Number of available GPU cards in the CDH instance
	GpuAvailable *uint64 `json:"GpuAvailable,omitnil,omitempty" name:"GpuAvailable"`
}

type Image struct {
	// Image ID
	ImageId *string `json:"ImageId,omitnil,omitempty" name:"ImageId"`

	// Operating system of the image
	OsName *string `json:"OsName,omitnil,omitempty" name:"OsName"`

	// Image type
	ImageType *string `json:"ImageType,omitnil,omitempty" name:"ImageType"`

	// Creation time of the image
	CreatedTime *string `json:"CreatedTime,omitnil,omitempty" name:"CreatedTime"`

	// Image name
	ImageName *string `json:"ImageName,omitnil,omitempty" name:"ImageName"`

	// Image description
	ImageDescription *string `json:"ImageDescription,omitnil,omitempty" name:"ImageDescription"`

	// Image size
	ImageSize *int64 `json:"ImageSize,omitnil,omitempty" name:"ImageSize"`

	// Image architecture
	Architecture *string `json:"Architecture,omitnil,omitempty" name:"Architecture"`

	// Image state
	ImageState *string `json:"ImageState,omitnil,omitempty" name:"ImageState"`

	// Source platform of the image
	Platform *string `json:"Platform,omitnil,omitempty" name:"Platform"`

	// Image creator
	ImageCreator *string `json:"ImageCreator,omitnil,omitempty" name:"ImageCreator"`

	// Image source
	ImageSource *string `json:"ImageSource,omitnil,omitempty" name:"ImageSource"`

	// Synchronization percentage
	// Note: This field may return null, indicating that no valid value is found.
	SyncPercent *int64 `json:"SyncPercent,omitnil,omitempty" name:"SyncPercent"`

	// Whether the image supports cloud-init
	// Note: This field may return null, indicating that no valid value is found.
	IsSupportCloudinit *bool `json:"IsSupportCloudinit,omitnil,omitempty" name:"IsSupportCloudinit"`

	// Information on the snapshots associated with the image
	// Note: This field may return null, indicating that no valid value is found.
	SnapshotSet []*Snapshot `json:"SnapshotSet,omitnil,omitempty" name:"SnapshotSet"`

	// The list of tags bound to the image.
	// Note: This field may return `null`, indicating that no valid value was found.
	Tags []*Tag `json:"Tags,omitnil,omitempty" name:"Tags"`

	// Image license type
	LicenseType *string `json:"LicenseType,omitnil,omitempty" name:"LicenseType"`

	// Image family, Note: This field may return empty
	ImageFamily *string `json:"ImageFamily,omitnil,omitempty" name:"ImageFamily"`

	// Whether the image is deprecated
	ImageDeprecated *bool `json:"ImageDeprecated,omitnil,omitempty" name:"ImageDeprecated"`
}

type ImageOsList struct {
	// Supported Windows OS
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	Windows []*string `json:"Windows,omitnil,omitempty" name:"Windows"`

	// Supported Linux OS
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	Linux []*string `json:"Linux,omitnil,omitempty" name:"Linux"`
}

// Predefined struct for user
type ImportImageRequestParams struct {
	// OS architecture of the image to be imported, `x86_64` or `i386`.
	Architecture *string `json:"Architecture,omitnil,omitempty" name:"Architecture"`

	// OS type of the image to be imported. You can call `DescribeImportImageOs` to obtain the list of supported operating systems.
	OsType *string `json:"OsType,omitnil,omitempty" name:"OsType"`

	// OS version of the image to be imported. You can call `DescribeImportImageOs` to obtain the list of supported operating systems.
	OsVersion *string `json:"OsVersion,omitnil,omitempty" name:"OsVersion"`

	// Address on COS where the image to be imported is stored.
	ImageUrl *string `json:"ImageUrl,omitnil,omitempty" name:"ImageUrl"`

	// Image name
	ImageName *string `json:"ImageName,omitnil,omitempty" name:"ImageName"`

	// Image description
	ImageDescription *string `json:"ImageDescription,omitnil,omitempty" name:"ImageDescription"`

	// Dry run to check the parameters without performing the operation
	DryRun *bool `json:"DryRun,omitnil,omitempty" name:"DryRun"`

	// Whether to force import the image. For more information, see [Forcibly Import Image](https://intl.cloud.tencent.com/document/product/213/12849).
	Force *bool `json:"Force,omitnil,omitempty" name:"Force"`

	// Tag description list. This parameter is used to bind a tag to a custom image.
	TagSpecification []*TagSpecification `json:"TagSpecification,omitnil,omitempty" name:"TagSpecification"`

	// The license type used to activate the OS after importing an image.
	// Valid values:
	// `TencentCloud`: Tencent Cloud official license
	// `BYOL`: Bring Your Own License
	LicenseType *string `json:"LicenseType,omitnil,omitempty" name:"LicenseType"`

	// Boot mode
	BootMode *string `json:"BootMode,omitnil,omitempty" name:"BootMode"`

	// Image family
	ImageFamily *string `json:"ImageFamily,omitnil,omitempty" name:"ImageFamily"`
}

type ImportImageRequest struct {
	*tchttp.BaseRequest
	
	// OS architecture of the image to be imported, `x86_64` or `i386`.
	Architecture *string `json:"Architecture,omitnil,omitempty" name:"Architecture"`

	// OS type of the image to be imported. You can call `DescribeImportImageOs` to obtain the list of supported operating systems.
	OsType *string `json:"OsType,omitnil,omitempty" name:"OsType"`

	// OS version of the image to be imported. You can call `DescribeImportImageOs` to obtain the list of supported operating systems.
	OsVersion *string `json:"OsVersion,omitnil,omitempty" name:"OsVersion"`

	// Address on COS where the image to be imported is stored.
	ImageUrl *string `json:"ImageUrl,omitnil,omitempty" name:"ImageUrl"`

	// Image name
	ImageName *string `json:"ImageName,omitnil,omitempty" name:"ImageName"`

	// Image description
	ImageDescription *string `json:"ImageDescription,omitnil,omitempty" name:"ImageDescription"`

	// Dry run to check the parameters without performing the operation
	DryRun *bool `json:"DryRun,omitnil,omitempty" name:"DryRun"`

	// Whether to force import the image. For more information, see [Forcibly Import Image](https://intl.cloud.tencent.com/document/product/213/12849).
	Force *bool `json:"Force,omitnil,omitempty" name:"Force"`

	// Tag description list. This parameter is used to bind a tag to a custom image.
	TagSpecification []*TagSpecification `json:"TagSpecification,omitnil,omitempty" name:"TagSpecification"`

	// The license type used to activate the OS after importing an image.
	// Valid values:
	// `TencentCloud`: Tencent Cloud official license
	// `BYOL`: Bring Your Own License
	LicenseType *string `json:"LicenseType,omitnil,omitempty" name:"LicenseType"`

	// Boot mode
	BootMode *string `json:"BootMode,omitnil,omitempty" name:"BootMode"`

	// Image family
	ImageFamily *string `json:"ImageFamily,omitnil,omitempty" name:"ImageFamily"`
}

func (r *ImportImageRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ImportImageRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Architecture")
	delete(f, "OsType")
	delete(f, "OsVersion")
	delete(f, "ImageUrl")
	delete(f, "ImageName")
	delete(f, "ImageDescription")
	delete(f, "DryRun")
	delete(f, "Force")
	delete(f, "TagSpecification")
	delete(f, "LicenseType")
	delete(f, "BootMode")
	delete(f, "ImageFamily")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ImportImageRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ImportImageResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ImportImageResponse struct {
	*tchttp.BaseResponse
	Response *ImportImageResponseParams `json:"Response"`
}

func (r *ImportImageResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ImportImageResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ImportKeyPairRequestParams struct {
	// Key pair name, which can contain numbers, letters, and underscores, with a maximum length of 25 characters.
	KeyName *string `json:"KeyName,omitnil,omitempty" name:"KeyName"`

	// The project ID to which the key pair belongs after it is created. <br><br>You can obtain the project ID in the following ways: <br><li>Check the project list in the [Project management](https://console.cloud.tencent.com/project) page.<br><li>Call the `DescribeProject` API and view the `projectId` in the response.
	// 
	// If you want to use the default project, specify 0 for the parameter.
	ProjectId *int64 `json:"ProjectId,omitnil,omitempty" name:"ProjectId"`

	// Content of the public key in the key pair in the `OpenSSH RSA` format.
	PublicKey *string `json:"PublicKey,omitnil,omitempty" name:"PublicKey"`

	// Tag description list. This parameter is used to bind a tag to a key pair.
	TagSpecification []*TagSpecification `json:"TagSpecification,omitnil,omitempty" name:"TagSpecification"`
}

type ImportKeyPairRequest struct {
	*tchttp.BaseRequest
	
	// Key pair name, which can contain numbers, letters, and underscores, with a maximum length of 25 characters.
	KeyName *string `json:"KeyName,omitnil,omitempty" name:"KeyName"`

	// The project ID to which the key pair belongs after it is created. <br><br>You can obtain the project ID in the following ways: <br><li>Check the project list in the [Project management](https://console.cloud.tencent.com/project) page.<br><li>Call the `DescribeProject` API and view the `projectId` in the response.
	// 
	// If you want to use the default project, specify 0 for the parameter.
	ProjectId *int64 `json:"ProjectId,omitnil,omitempty" name:"ProjectId"`

	// Content of the public key in the key pair in the `OpenSSH RSA` format.
	PublicKey *string `json:"PublicKey,omitnil,omitempty" name:"PublicKey"`

	// Tag description list. This parameter is used to bind a tag to a key pair.
	TagSpecification []*TagSpecification `json:"TagSpecification,omitnil,omitempty" name:"TagSpecification"`
}

func (r *ImportKeyPairRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ImportKeyPairRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "KeyName")
	delete(f, "ProjectId")
	delete(f, "PublicKey")
	delete(f, "TagSpecification")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ImportKeyPairRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ImportKeyPairResponseParams struct {
	// Key pair ID
	KeyId *string `json:"KeyId,omitnil,omitempty" name:"KeyId"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ImportKeyPairResponse struct {
	*tchttp.BaseResponse
	Response *ImportKeyPairResponseParams `json:"Response"`
}

func (r *ImportKeyPairResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ImportKeyPairResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type InquirePricePurchaseReservedInstancesOfferingRequestParams struct {
	// The number of the reserved instances you are purchasing.
	InstanceCount *uint64 `json:"InstanceCount,omitnil,omitempty" name:"InstanceCount"`

	// The ID of the reserved instance offering.
	ReservedInstancesOfferingId *string `json:"ReservedInstancesOfferingId,omitnil,omitempty" name:"ReservedInstancesOfferingId"`

	// Dry run.
	DryRun *bool `json:"DryRun,omitnil,omitempty" name:"DryRun"`

	// A unique string supplied by the client to ensure that the request is idempotent. Its maximum length is 64 ASCII characters. If this parameter is not specified, the idempotency of the request cannot be guaranteed.<br>For more information, see Ensuring Idempotency.
	ClientToken *string `json:"ClientToken,omitnil,omitempty" name:"ClientToken"`

	// Reserved instance name.<br><li>The RI name defaults to "unnamed" if this parameter is left empty.</li><li>You can enter any name within 60 characters (including the pattern string).</li>
	ReservedInstanceName *string `json:"ReservedInstanceName,omitnil,omitempty" name:"ReservedInstanceName"`
}

type InquirePricePurchaseReservedInstancesOfferingRequest struct {
	*tchttp.BaseRequest
	
	// The number of the reserved instances you are purchasing.
	InstanceCount *uint64 `json:"InstanceCount,omitnil,omitempty" name:"InstanceCount"`

	// The ID of the reserved instance offering.
	ReservedInstancesOfferingId *string `json:"ReservedInstancesOfferingId,omitnil,omitempty" name:"ReservedInstancesOfferingId"`

	// Dry run.
	DryRun *bool `json:"DryRun,omitnil,omitempty" name:"DryRun"`

	// A unique string supplied by the client to ensure that the request is idempotent. Its maximum length is 64 ASCII characters. If this parameter is not specified, the idempotency of the request cannot be guaranteed.<br>For more information, see Ensuring Idempotency.
	ClientToken *string `json:"ClientToken,omitnil,omitempty" name:"ClientToken"`

	// Reserved instance name.<br><li>The RI name defaults to "unnamed" if this parameter is left empty.</li><li>You can enter any name within 60 characters (including the pattern string).</li>
	ReservedInstanceName *string `json:"ReservedInstanceName,omitnil,omitempty" name:"ReservedInstanceName"`
}

func (r *InquirePricePurchaseReservedInstancesOfferingRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *InquirePricePurchaseReservedInstancesOfferingRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceCount")
	delete(f, "ReservedInstancesOfferingId")
	delete(f, "DryRun")
	delete(f, "ClientToken")
	delete(f, "ReservedInstanceName")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "InquirePricePurchaseReservedInstancesOfferingRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type InquirePricePurchaseReservedInstancesOfferingResponseParams struct {
	// Price of the reserved instance with specified configuration.
	Price *ReservedInstancePrice `json:"Price,omitnil,omitempty" name:"Price"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type InquirePricePurchaseReservedInstancesOfferingResponse struct {
	*tchttp.BaseResponse
	Response *InquirePricePurchaseReservedInstancesOfferingResponseParams `json:"Response"`
}

func (r *InquirePricePurchaseReservedInstancesOfferingResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *InquirePricePurchaseReservedInstancesOfferingResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type InquiryPriceModifyInstancesChargeTypeRequestParams struct {
	// One or more instance IDs to be operated. You can obtain the instance ID through the `InstanceId` in the return value from the API [DescribeInstances](https://www.tencentcloud.com/document/api/213/15728?from_cn_redirect=1). The maximum number of instances per request is 20.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// Modified instance [billing type](https://www.tencentcloud.com/document/product/213/2180?from_cn_redirect=1). <br><li>`PREPAID`: monthly subscription.</li>
	// 
	// **Note:** Only supports converting pay-as-you-go instances to annual and monthly subscription instances.
	// 
	// default value: `PREPAID`
	InstanceChargeType *string `json:"InstanceChargeType,omitnil,omitempty" name:"InstanceChargeType"`

	// Prepaid mode, parameter settings related to monthly/annual subscription. through this parameter, specify the purchase duration of annual and monthly subscription instances, whether to enable auto-renewal, and other attributes. 
	InstanceChargePrepaid *InstanceChargePrepaid `json:"InstanceChargePrepaid,omitnil,omitempty" name:"InstanceChargePrepaid"`

	// Whether to switch the billing mode of elastic data cloud disks simultaneously. valid values: <br><li>true: means switching the billing mode of elastic data cloud disks</li><li>false: means not switching the billing mode of elastic data cloud disks</li><br>default value: false.
	ModifyPortableDataDisk *bool `json:"ModifyPortableDataDisk,omitnil,omitempty" name:"ModifyPortableDataDisk"`
}

type InquiryPriceModifyInstancesChargeTypeRequest struct {
	*tchttp.BaseRequest
	
	// One or more instance IDs to be operated. You can obtain the instance ID through the `InstanceId` in the return value from the API [DescribeInstances](https://www.tencentcloud.com/document/api/213/15728?from_cn_redirect=1). The maximum number of instances per request is 20.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// Modified instance [billing type](https://www.tencentcloud.com/document/product/213/2180?from_cn_redirect=1). <br><li>`PREPAID`: monthly subscription.</li>
	// 
	// **Note:** Only supports converting pay-as-you-go instances to annual and monthly subscription instances.
	// 
	// default value: `PREPAID`
	InstanceChargeType *string `json:"InstanceChargeType,omitnil,omitempty" name:"InstanceChargeType"`

	// Prepaid mode, parameter settings related to monthly/annual subscription. through this parameter, specify the purchase duration of annual and monthly subscription instances, whether to enable auto-renewal, and other attributes. 
	InstanceChargePrepaid *InstanceChargePrepaid `json:"InstanceChargePrepaid,omitnil,omitempty" name:"InstanceChargePrepaid"`

	// Whether to switch the billing mode of elastic data cloud disks simultaneously. valid values: <br><li>true: means switching the billing mode of elastic data cloud disks</li><li>false: means not switching the billing mode of elastic data cloud disks</li><br>default value: false.
	ModifyPortableDataDisk *bool `json:"ModifyPortableDataDisk,omitnil,omitempty" name:"ModifyPortableDataDisk"`
}

func (r *InquiryPriceModifyInstancesChargeTypeRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *InquiryPriceModifyInstancesChargeTypeRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceIds")
	delete(f, "InstanceChargeType")
	delete(f, "InstanceChargePrepaid")
	delete(f, "ModifyPortableDataDisk")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "InquiryPriceModifyInstancesChargeTypeRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type InquiryPriceModifyInstancesChargeTypeResponseParams struct {
	// This parameter indicates the price for switching the billing mode of the corresponding configuration instance.
	Price *Price `json:"Price,omitnil,omitempty" name:"Price"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type InquiryPriceModifyInstancesChargeTypeResponse struct {
	*tchttp.BaseResponse
	Response *InquiryPriceModifyInstancesChargeTypeResponseParams `json:"Response"`
}

func (r *InquiryPriceModifyInstancesChargeTypeResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *InquiryPriceModifyInstancesChargeTypeResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type InquiryPriceRenewInstancesRequestParams struct {
	// One or more instance IDs to be operated. You can obtain the instance ID through the `InstanceId` in the return value from the API [DescribeInstances](https://www.tencentcloud.com/zh/document/api/213/33258). The maximum number of instances per request is 100.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// Prepaid mode, that is, monthly subscription-related parameter settings. You can specify attributes of a monthly subscription instance, such as purchase duration and whether to set auto-renewal, through this parameter.
	InstanceChargePrepaid *InstanceChargePrepaid `json:"InstanceChargePrepaid,omitnil,omitempty" name:"InstanceChargePrepaid"`

	// Trial run, for testing purposes, does not execute specific logic. valid values: <li>`true`: skip execution logic</li> <li>`false`: execute logic</li>  default value: `false`.
	DryRun *bool `json:"DryRun,omitnil,omitempty" name:"DryRun"`

	// Whether to renew the elastic data disk. valid values:<br><li>true: indicates renewing the annual and monthly subscription instance and its mounted elastic data disk simultaneously</li><li>false: indicates renewing the annual and monthly subscription instance while no longer renewing its mounted elastic data disk</li><br>default value: true.
	RenewPortableDataDisk *bool `json:"RenewPortableDataDisk,omitnil,omitempty" name:"RenewPortableDataDisk"`
}

type InquiryPriceRenewInstancesRequest struct {
	*tchttp.BaseRequest
	
	// One or more instance IDs to be operated. You can obtain the instance ID through the `InstanceId` in the return value from the API [DescribeInstances](https://www.tencentcloud.com/zh/document/api/213/33258). The maximum number of instances per request is 100.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// Prepaid mode, that is, monthly subscription-related parameter settings. You can specify attributes of a monthly subscription instance, such as purchase duration and whether to set auto-renewal, through this parameter.
	InstanceChargePrepaid *InstanceChargePrepaid `json:"InstanceChargePrepaid,omitnil,omitempty" name:"InstanceChargePrepaid"`

	// Trial run, for testing purposes, does not execute specific logic. valid values: <li>`true`: skip execution logic</li> <li>`false`: execute logic</li>  default value: `false`.
	DryRun *bool `json:"DryRun,omitnil,omitempty" name:"DryRun"`

	// Whether to renew the elastic data disk. valid values:<br><li>true: indicates renewing the annual and monthly subscription instance and its mounted elastic data disk simultaneously</li><li>false: indicates renewing the annual and monthly subscription instance while no longer renewing its mounted elastic data disk</li><br>default value: true.
	RenewPortableDataDisk *bool `json:"RenewPortableDataDisk,omitnil,omitempty" name:"RenewPortableDataDisk"`
}

func (r *InquiryPriceRenewInstancesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *InquiryPriceRenewInstancesRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceIds")
	delete(f, "InstanceChargePrepaid")
	delete(f, "DryRun")
	delete(f, "RenewPortableDataDisk")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "InquiryPriceRenewInstancesRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type InquiryPriceRenewInstancesResponseParams struct {
	// This parameter indicates the price for the corresponding configuration instance.
	Price *Price `json:"Price,omitnil,omitempty" name:"Price"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type InquiryPriceRenewInstancesResponse struct {
	*tchttp.BaseResponse
	Response *InquiryPriceRenewInstancesResponseParams `json:"Response"`
}

func (r *InquiryPriceRenewInstancesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *InquiryPriceRenewInstancesResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type InquiryPriceResetInstanceRequestParams struct {
	// Instance ID. To obtain the instance IDs, you can call [`DescribeInstances`](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1) and look for `InstanceId` in the response.
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// [Image](https://intl.cloud.tencent.com/document/product/213/4940?from_cn_redirect=1) ID in the format of `img-xxx`. There are four types of images: <br/><li>Public images </li><li>Custom images </li><li>Shared images </li><li>Marketplace images </li><br/>You can obtain the available image IDs in the following ways: <br/><li>For IDs of `public images`, `custom images`, and `shared images`, log in to the [console](https://console.cloud.tencent.com/cvm/image?rid=1&imageType=PUBLIC_IMAGE) to query the information; for IDs of `marketplace images`, go to [Cloud Marketplace](https://market.cloud.tencent.com/list). </li><li>Call [DescribeImages](https://intl.cloud.tencent.com/document/api/213/15715?from_cn_redirect=1) and look for `ImageId` in the response.</li>
	ImageId *string `json:"ImageId,omitnil,omitempty" name:"ImageId"`

	// Configuration of the system disk of the instance. For instances with a cloud disk as the system disk, you can expand the system disk by using this parameter to specify the new capacity after reinstallation. If the parameter is not specified, the system disk capacity remains unchanged by default. You can only expand the capacity of the system disk; reducing its capacity is not supported. When reinstalling the system, you can only modify the capacity of the system disk, not the type.
	SystemDisk *SystemDisk `json:"SystemDisk,omitnil,omitempty" name:"SystemDisk"`

	// Login settings of the instance. You can use this parameter to set the login method, password, and key of the instance or keep the login settings of the original image. By default, a random password will be generated and sent to you via the Message Center.
	LoginSettings *LoginSettings `json:"LoginSettings,omitnil,omitempty" name:"LoginSettings"`

	// Enhanced services. You can use this parameter to specify whether to enable services such as Cloud Monitor and Cloud Security. If this parameter is not specified, Cloud Monitor and Cloud Security will be enabled by default.
	EnhancedService *EnhancedService `json:"EnhancedService,omitnil,omitempty" name:"EnhancedService"`


	UserData *string `json:"UserData,omitnil,omitempty" name:"UserData"`
}

type InquiryPriceResetInstanceRequest struct {
	*tchttp.BaseRequest
	
	// Instance ID. To obtain the instance IDs, you can call [`DescribeInstances`](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1) and look for `InstanceId` in the response.
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// [Image](https://intl.cloud.tencent.com/document/product/213/4940?from_cn_redirect=1) ID in the format of `img-xxx`. There are four types of images: <br/><li>Public images </li><li>Custom images </li><li>Shared images </li><li>Marketplace images </li><br/>You can obtain the available image IDs in the following ways: <br/><li>For IDs of `public images`, `custom images`, and `shared images`, log in to the [console](https://console.cloud.tencent.com/cvm/image?rid=1&imageType=PUBLIC_IMAGE) to query the information; for IDs of `marketplace images`, go to [Cloud Marketplace](https://market.cloud.tencent.com/list). </li><li>Call [DescribeImages](https://intl.cloud.tencent.com/document/api/213/15715?from_cn_redirect=1) and look for `ImageId` in the response.</li>
	ImageId *string `json:"ImageId,omitnil,omitempty" name:"ImageId"`

	// Configuration of the system disk of the instance. For instances with a cloud disk as the system disk, you can expand the system disk by using this parameter to specify the new capacity after reinstallation. If the parameter is not specified, the system disk capacity remains unchanged by default. You can only expand the capacity of the system disk; reducing its capacity is not supported. When reinstalling the system, you can only modify the capacity of the system disk, not the type.
	SystemDisk *SystemDisk `json:"SystemDisk,omitnil,omitempty" name:"SystemDisk"`

	// Login settings of the instance. You can use this parameter to set the login method, password, and key of the instance or keep the login settings of the original image. By default, a random password will be generated and sent to you via the Message Center.
	LoginSettings *LoginSettings `json:"LoginSettings,omitnil,omitempty" name:"LoginSettings"`

	// Enhanced services. You can use this parameter to specify whether to enable services such as Cloud Monitor and Cloud Security. If this parameter is not specified, Cloud Monitor and Cloud Security will be enabled by default.
	EnhancedService *EnhancedService `json:"EnhancedService,omitnil,omitempty" name:"EnhancedService"`

	UserData *string `json:"UserData,omitnil,omitempty" name:"UserData"`
}

func (r *InquiryPriceResetInstanceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *InquiryPriceResetInstanceRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceId")
	delete(f, "ImageId")
	delete(f, "SystemDisk")
	delete(f, "LoginSettings")
	delete(f, "EnhancedService")
	delete(f, "UserData")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "InquiryPriceResetInstanceRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type InquiryPriceResetInstanceResponseParams struct {
	// Price of reinstalling the instance with the specified configuration.
	Price *Price `json:"Price,omitnil,omitempty" name:"Price"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type InquiryPriceResetInstanceResponse struct {
	*tchttp.BaseResponse
	Response *InquiryPriceResetInstanceResponseParams `json:"Response"`
}

func (r *InquiryPriceResetInstanceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *InquiryPriceResetInstanceResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type InquiryPriceResetInstancesInternetMaxBandwidthRequestParams struct {
	// Instance ID(s). To obtain the instance IDs, you can call [`DescribeInstances`](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1) and look for `InstanceId` in the response. The maximum number of instances in each request is 100. When changing the bandwidth of instances with `BANDWIDTH_PREPAID` or `BANDWIDTH_POSTPAID_BY_HOUR` as the network billing method, you can only specify one instance at a time.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// Configuration of public network egress bandwidth. The maximum bandwidth varies among different models. For more information, see the documentation on bandwidth limits. Currently only the `InternetMaxBandwidthOut` parameter is supported.
	InternetAccessible *InternetAccessible `json:"InternetAccessible,omitnil,omitempty" name:"InternetAccessible"`

	// Date from which the new bandwidth takes effect. Format: `YYYY-MM-DD`, such as `2016-10-30`. The starting date cannot be earlier than the current date. If the starting date is the current date, the new bandwidth takes effect immediately. This parameter is only valid for prepaid bandwidth. If you specify the parameter for bandwidth with other network billing methods, an error code will be returned.
	StartTime *string `json:"StartTime,omitnil,omitempty" name:"StartTime"`

	// Date until which the bandwidth takes effect, in the format of `YYYY-MM-DD`, such as `2016-10-30`. The validity period of the new bandwidth covers the end date. The end date should not be later than the expiration date of a monthly subscription instance. You can obtain the expiration date of an instance through the `ExpiredTime` in the return value from the API [DescribeInstances](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1). This parameter is only valid for monthly subscription bandwidth, and is not supported for bandwidth billed by other modes. Otherwise, the API will return a corresponding error code.
	EndTime *string `json:"EndTime,omitnil,omitempty" name:"EndTime"`
}

type InquiryPriceResetInstancesInternetMaxBandwidthRequest struct {
	*tchttp.BaseRequest
	
	// Instance ID(s). To obtain the instance IDs, you can call [`DescribeInstances`](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1) and look for `InstanceId` in the response. The maximum number of instances in each request is 100. When changing the bandwidth of instances with `BANDWIDTH_PREPAID` or `BANDWIDTH_POSTPAID_BY_HOUR` as the network billing method, you can only specify one instance at a time.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// Configuration of public network egress bandwidth. The maximum bandwidth varies among different models. For more information, see the documentation on bandwidth limits. Currently only the `InternetMaxBandwidthOut` parameter is supported.
	InternetAccessible *InternetAccessible `json:"InternetAccessible,omitnil,omitempty" name:"InternetAccessible"`

	// Date from which the new bandwidth takes effect. Format: `YYYY-MM-DD`, such as `2016-10-30`. The starting date cannot be earlier than the current date. If the starting date is the current date, the new bandwidth takes effect immediately. This parameter is only valid for prepaid bandwidth. If you specify the parameter for bandwidth with other network billing methods, an error code will be returned.
	StartTime *string `json:"StartTime,omitnil,omitempty" name:"StartTime"`

	// Date until which the bandwidth takes effect, in the format of `YYYY-MM-DD`, such as `2016-10-30`. The validity period of the new bandwidth covers the end date. The end date should not be later than the expiration date of a monthly subscription instance. You can obtain the expiration date of an instance through the `ExpiredTime` in the return value from the API [DescribeInstances](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1). This parameter is only valid for monthly subscription bandwidth, and is not supported for bandwidth billed by other modes. Otherwise, the API will return a corresponding error code.
	EndTime *string `json:"EndTime,omitnil,omitempty" name:"EndTime"`
}

func (r *InquiryPriceResetInstancesInternetMaxBandwidthRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *InquiryPriceResetInstancesInternetMaxBandwidthRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceIds")
	delete(f, "InternetAccessible")
	delete(f, "StartTime")
	delete(f, "EndTime")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "InquiryPriceResetInstancesInternetMaxBandwidthRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type InquiryPriceResetInstancesInternetMaxBandwidthResponseParams struct {
	// Price of the new bandwidth
	Price *Price `json:"Price,omitnil,omitempty" name:"Price"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type InquiryPriceResetInstancesInternetMaxBandwidthResponse struct {
	*tchttp.BaseResponse
	Response *InquiryPriceResetInstancesInternetMaxBandwidthResponseParams `json:"Response"`
}

func (r *InquiryPriceResetInstancesInternetMaxBandwidthResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *InquiryPriceResetInstancesInternetMaxBandwidthResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type InquiryPriceResetInstancesTypeRequestParams struct {
	// One or more instance IDs to be operated. You can obtain the instance ID through the `InstanceId` in the return value from the API [DescribeInstances](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1). The maximum number of instances per request is 1.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// Instance model. Resources vary with the instance model. Specific values can be found in the tables of [Instance Types] (https://intl.cloud.tencent.com/document/product/213/11518?from_cn_redirect=1) or in the latest specifications via the [DescribeInstanceTypeConfigs] (https://intl.cloud.tencent.com/document/product/213/15749?from_cn_redirect=1) API.
	InstanceType *string `json:"InstanceType,omitnil,omitempty" name:"InstanceType"`
}

type InquiryPriceResetInstancesTypeRequest struct {
	*tchttp.BaseRequest
	
	// One or more instance IDs to be operated. You can obtain the instance ID through the `InstanceId` in the return value from the API [DescribeInstances](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1). The maximum number of instances per request is 1.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// Instance model. Resources vary with the instance model. Specific values can be found in the tables of [Instance Types] (https://intl.cloud.tencent.com/document/product/213/11518?from_cn_redirect=1) or in the latest specifications via the [DescribeInstanceTypeConfigs] (https://intl.cloud.tencent.com/document/product/213/15749?from_cn_redirect=1) API.
	InstanceType *string `json:"InstanceType,omitnil,omitempty" name:"InstanceType"`
}

func (r *InquiryPriceResetInstancesTypeRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *InquiryPriceResetInstancesTypeRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceIds")
	delete(f, "InstanceType")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "InquiryPriceResetInstancesTypeRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type InquiryPriceResetInstancesTypeResponseParams struct {
	// Price of the instance using the specified model
	Price *Price `json:"Price,omitnil,omitempty" name:"Price"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type InquiryPriceResetInstancesTypeResponse struct {
	*tchttp.BaseResponse
	Response *InquiryPriceResetInstancesTypeResponseParams `json:"Response"`
}

func (r *InquiryPriceResetInstancesTypeResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *InquiryPriceResetInstancesTypeResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type InquiryPriceResizeInstanceDisksRequestParams struct {
	// Instance ID. To obtain the instance IDs, you can call [`DescribeInstances`](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1) and look for `InstanceId` in the response.
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// Configuration information of a data disk to be expanded. Only inelastic data disks (with `Portable` being `false` in the return values of [DescribeDisks](https://intl.cloud.tencent.com/document/api/362/16315?from_cn_redirect=1)) can be expanded. The unit of data disk capacity is GB. The minimum expansion step is 10 GB. For more information about data disk types, refer to Disk Product Introduction. The available data disk type is restricted by the instance type `InstanceType`. Additionally, the maximum allowable capacity for expansion varies by data disk type.
	// <dx-alert infotype="explain" title="">You should specify either DataDisks or SystemDisk, but you cannot specify both at the same time.</dx-alert>
	DataDisks []*DataDisk `json:"DataDisks,omitnil,omitempty" name:"DataDisks"`

	// Whether to forcibly shut down a running instance. It is recommended to manually shut down a running instance first and then reset the user password. Valid values:<br><li>true: Forcibly shut down an instance after a normal shutdown fails.</li><br><li>false: Do not forcibly shut down an instance after a normal shutdown fails.</li><br><br>Default value: false.<br><br>Forced shutdown is equivalent to turning off a physical computer's power switch. Forced shutdown may cause data loss or file system corruption and should only be used when a server cannot be shut down normally.
	ForceStop *bool `json:"ForceStop,omitnil,omitempty" name:"ForceStop"`
}

type InquiryPriceResizeInstanceDisksRequest struct {
	*tchttp.BaseRequest
	
	// Instance ID. To obtain the instance IDs, you can call [`DescribeInstances`](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1) and look for `InstanceId` in the response.
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// Configuration information of a data disk to be expanded. Only inelastic data disks (with `Portable` being `false` in the return values of [DescribeDisks](https://intl.cloud.tencent.com/document/api/362/16315?from_cn_redirect=1)) can be expanded. The unit of data disk capacity is GB. The minimum expansion step is 10 GB. For more information about data disk types, refer to Disk Product Introduction. The available data disk type is restricted by the instance type `InstanceType`. Additionally, the maximum allowable capacity for expansion varies by data disk type.
	// <dx-alert infotype="explain" title="">You should specify either DataDisks or SystemDisk, but you cannot specify both at the same time.</dx-alert>
	DataDisks []*DataDisk `json:"DataDisks,omitnil,omitempty" name:"DataDisks"`

	// Whether to forcibly shut down a running instance. It is recommended to manually shut down a running instance first and then reset the user password. Valid values:<br><li>true: Forcibly shut down an instance after a normal shutdown fails.</li><br><li>false: Do not forcibly shut down an instance after a normal shutdown fails.</li><br><br>Default value: false.<br><br>Forced shutdown is equivalent to turning off a physical computer's power switch. Forced shutdown may cause data loss or file system corruption and should only be used when a server cannot be shut down normally.
	ForceStop *bool `json:"ForceStop,omitnil,omitempty" name:"ForceStop"`
}

func (r *InquiryPriceResizeInstanceDisksRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *InquiryPriceResizeInstanceDisksRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceId")
	delete(f, "DataDisks")
	delete(f, "ForceStop")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "InquiryPriceResizeInstanceDisksRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type InquiryPriceResizeInstanceDisksResponseParams struct {
	// Price of the disks after being expanded to the specified configurations
	Price *Price `json:"Price,omitnil,omitempty" name:"Price"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type InquiryPriceResizeInstanceDisksResponse struct {
	*tchttp.BaseResponse
	Response *InquiryPriceResizeInstanceDisksResponseParams `json:"Response"`
}

func (r *InquiryPriceResizeInstanceDisksResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *InquiryPriceResizeInstanceDisksResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type InquiryPriceRunInstancesRequestParams struct {
	// Location of the instance. You can use this parameter to specify the attributes of the instance, such as its availability zone and project.
	//  <b>Note: `Placement` is required when `LaunchTemplate` is not specified. If both the parameters are passed in, `Placement` prevails.</b>
	Placement *Placement `json:"Placement,omitnil,omitempty" name:"Placement"`

	// [Image](https://intl.cloud.tencent.com/document/product/213/4940?from_cn_redirect=1) ID in the format of `img-xxx`. There are three types of images: <br/><li>Public images</li><li>Custom images</li><li>Shared images</li><br/>You can obtain the available image IDs in the following ways: <br/><li>For IDs of `public images`, `custom images`, and `shared images`, log in to the [CVM console](https://console.tencentcloud.com/cvm/image/index?rid=1&tab=PUBLIC_IMAGE&imageType=PUBLIC_IMAGE) to query the information. </li><li>Call [DescribeImages](https://intl.cloud.tencent.com/document/api/213/15715?from_cn_redirect=1) and look for `ImageId` in the response.</li>
	//  <b>Note: `ImageId` is required when `LaunchTemplate` is not specified. If both the parameters are passed in, `ImageId` prevails.</b>
	ImageId *string `json:"ImageId,omitnil,omitempty" name:"ImageId"`

	// The instance [billing method](https://intl.cloud.tencent.com/document/product/213/2180?from_cn_redirect=1).<br><li>POSTPAID_BY_HOUR: Pay-as-you-go on an hourly basis<br>Default value: POSTPAID_BY_HOUR.
	InstanceChargeType *string `json:"InstanceChargeType,omitnil,omitempty" name:"InstanceChargeType"`

	// Details of the monthly subscription, including the purchase period, auto-renewal. It is required if the `InstanceChargeType` is `PREPAID`. 
	InstanceChargePrepaid *InstanceChargePrepaid `json:"InstanceChargePrepaid,omitnil,omitempty" name:"InstanceChargePrepaid"`

	// Instance model. Different instance models specify different resource specifications.
	// <br><li>For instances created with the payment modes PREPAID or POSTPAID_BY_HOUR, specifies the specific values obtained BY calling the [DescribeInstanceTypeConfig](https://www.tencentcloud.com/document/product/1119/45686?lang=en) api for the latest specification table or referring to [instance specifications](https://www.tencentcloud.com/document/product/213/11518). if not specified, the system will dynamically assign a default model based on the current resource sales situation in a region.</li><br><li>for instances created with the payment mode CDHPAID, indicates this parameter uses the prefix "CDH_" and is generated based on CPU and memory configuration. the specific format is: CDH_XCXG. for example, for creating a CDH instance with 1 CPU core and 1 gb memory, this parameter should be CDH_1C1G.</li>.
	InstanceType *string `json:"InstanceType,omitnil,omitempty" name:"InstanceType"`

	// System disk configuration of the instance. If this parameter is not specified, the default value will be used.
	SystemDisk *SystemDisk `json:"SystemDisk,omitnil,omitempty" name:"SystemDisk"`

	// Data disk configuration of the instance. If the parameter is not specified, no data disk will be purchased by default. If you want to purchase data disks, you can specify 21 data disks, including up to 1 `LOCAL_BASIC` data disk or `LOCAL_SSD` data disk and up to 20 `CLOUD_BASIC` data disks, `CLOUD_PREMIUM` data disks, or `CLOUD_SSD` data disks.
	DataDisks []*DataDisk `json:"DataDisks,omitnil,omitempty" name:"DataDisks"`

	// VPC configurations (VPC ID, subnet ID, etc). If It's not specified, the classic network will be used by default. If a VPC IP is specified in this parameter, the `InstanceCount` can only be 1.
	VirtualPrivateCloud *VirtualPrivateCloud `json:"VirtualPrivateCloud,omitnil,omitempty" name:"VirtualPrivateCloud"`

	// Configuration of public network bandwidth. If it's not specified, 0 Mbps is used by default.
	InternetAccessible *InternetAccessible `json:"InternetAccessible,omitnil,omitempty" name:"InternetAccessible"`

	// Number of instances to purchase. Value range: 1 (default) to 100. It cannot exceed the remaining CVM quota of the user. For more information on quota, see [Restrictions on CVM Instance Purchase](https://intl.cloud.tencent.com/document/product/213/2664).
	InstanceCount *int64 `json:"InstanceCount,omitnil,omitempty" name:"InstanceCount"`

	// Instance display name. <li>if no instance display name is specified, it will display 'unnamed' by default.</li> <li>when purchasing multiple instances, if the pattern string `{R:x}` is specified, it means generating numbers `[x, x+n-1]`, where `n` represents the number of purchased instances. for example, `server_{R:3}`: when purchasing 1 instance, the instance display name is `server_3`; when purchasing 2 instances, the instance display names are `server_3` and `server_4` respectively. supports specifying multiple pattern strings `{R:x}`.</li> <li>when purchasing multiple instances without specifying a pattern string, suffixes `1, 2...n` will be added to the instance display name, where `n` represents the number of purchased instances. for example, `server_`: when purchasing 2 instances, the instance display names are `server_1` and `server_2` respectively.</li> <li>supports up to 128 characters (including pattern strings).</li>.
	InstanceName *string `json:"InstanceName,omitnil,omitempty" name:"InstanceName"`

	// Login settings of the instance. You can use this parameter to set the login method, password, and key of the instance, or keep the original login settings of the image. By default, a random password will be generated and sent to you via the Message Center.
	LoginSettings *LoginSettings `json:"LoginSettings,omitnil,omitempty" name:"LoginSettings"`

	// Security group to which an instance belongs. obtain this parameter by calling the `SecurityGroupId` field in the return value of [DescribeSecurityGroups](https://www.tencentcloud.com/document/product/215/15808). if not specified, bind the default security group under the designated project. if the default security group does not exist, automatically create it.
	SecurityGroupIds []*string `json:"SecurityGroupIds,omitnil,omitempty" name:"SecurityGroupIds"`

	// Enhanced services. You can use this parameter to specify whether to enable services such as Cloud Security and Cloud Monitor. If this parameter is not specified, Cloud Monitor and Cloud Security will be enabled by default.
	EnhancedService *EnhancedService `json:"EnhancedService,omitnil,omitempty" name:"EnhancedService"`

	// A unique string supplied by the client to ensure that the request is idempotent. Its maximum length is 64 ASCII characters. If this parameter is not specified, the idempotency of the request cannot be guaranteed.<br>For more information, see Ensuring Idempotency.
	ClientToken *string `json:"ClientToken,omitnil,omitempty" name:"ClientToken"`

	// Hostname of Cloud Virtual Machine.<br><li>Period (.) and hyphen (-) should not be used as the first or last character of the hostname, and should not be used consecutively.</li><br><li>Windows instances: The hostname should contain 2 to 15 characters, including letters (case insensitive), digits, and hyphens (-), does not support periods (.), and should not be all digits.</li><br><li>Instances of other types (such as Linux instances): The hostname should contain 2 to 30 characters, including multiple periods (.), with each segment between periods considered as one section. Each section can contain letters (case insensitive), digits, and hyphens (-).</li>
	HostName *string `json:"HostName,omitnil,omitempty" name:"HostName"`

	// The tag description list. This parameter is used to bind a tag to a resource instance. A tag can only be bound to CVM instances.
	TagSpecification []*TagSpecification `json:"TagSpecification,omitnil,omitempty" name:"TagSpecification"`

	// The market options of the instance.
	InstanceMarketOptions *InstanceMarketOptionsRequest `json:"InstanceMarketOptions,omitnil,omitempty" name:"InstanceMarketOptions"`

	// Custom metadata, supports creating key-value pairs of custom metadata when creating a CVM.
	// 
	// **Note: this field is in beta test.**.
	Metadata *Metadata `json:"Metadata,omitnil,omitempty" name:"Metadata"`

	// HPC cluster ID.
	HpcClusterId *string `json:"HpcClusterId,omitnil,omitempty" name:"HpcClusterId"`

	// Information about the CPU topology of an instance. If not specified, it is determined by system resources.
	CpuTopology *CpuTopology `json:"CpuTopology,omitnil,omitempty" name:"CpuTopology"`


	LaunchTemplate *LaunchTemplate `json:"LaunchTemplate,omitnil,omitempty" name:"LaunchTemplate"`
}

type InquiryPriceRunInstancesRequest struct {
	*tchttp.BaseRequest
	
	// Location of the instance. You can use this parameter to specify the attributes of the instance, such as its availability zone and project.
	//  <b>Note: `Placement` is required when `LaunchTemplate` is not specified. If both the parameters are passed in, `Placement` prevails.</b>
	Placement *Placement `json:"Placement,omitnil,omitempty" name:"Placement"`

	// [Image](https://intl.cloud.tencent.com/document/product/213/4940?from_cn_redirect=1) ID in the format of `img-xxx`. There are three types of images: <br/><li>Public images</li><li>Custom images</li><li>Shared images</li><br/>You can obtain the available image IDs in the following ways: <br/><li>For IDs of `public images`, `custom images`, and `shared images`, log in to the [CVM console](https://console.tencentcloud.com/cvm/image/index?rid=1&tab=PUBLIC_IMAGE&imageType=PUBLIC_IMAGE) to query the information. </li><li>Call [DescribeImages](https://intl.cloud.tencent.com/document/api/213/15715?from_cn_redirect=1) and look for `ImageId` in the response.</li>
	//  <b>Note: `ImageId` is required when `LaunchTemplate` is not specified. If both the parameters are passed in, `ImageId` prevails.</b>
	ImageId *string `json:"ImageId,omitnil,omitempty" name:"ImageId"`

	// The instance [billing method](https://intl.cloud.tencent.com/document/product/213/2180?from_cn_redirect=1).<br><li>POSTPAID_BY_HOUR: Pay-as-you-go on an hourly basis<br>Default value: POSTPAID_BY_HOUR.
	InstanceChargeType *string `json:"InstanceChargeType,omitnil,omitempty" name:"InstanceChargeType"`

	// Details of the monthly subscription, including the purchase period, auto-renewal. It is required if the `InstanceChargeType` is `PREPAID`. 
	InstanceChargePrepaid *InstanceChargePrepaid `json:"InstanceChargePrepaid,omitnil,omitempty" name:"InstanceChargePrepaid"`

	// Instance model. Different instance models specify different resource specifications.
	// <br><li>For instances created with the payment modes PREPAID or POSTPAID_BY_HOUR, specifies the specific values obtained BY calling the [DescribeInstanceTypeConfig](https://www.tencentcloud.com/document/product/1119/45686?lang=en) api for the latest specification table or referring to [instance specifications](https://www.tencentcloud.com/document/product/213/11518). if not specified, the system will dynamically assign a default model based on the current resource sales situation in a region.</li><br><li>for instances created with the payment mode CDHPAID, indicates this parameter uses the prefix "CDH_" and is generated based on CPU and memory configuration. the specific format is: CDH_XCXG. for example, for creating a CDH instance with 1 CPU core and 1 gb memory, this parameter should be CDH_1C1G.</li>.
	InstanceType *string `json:"InstanceType,omitnil,omitempty" name:"InstanceType"`

	// System disk configuration of the instance. If this parameter is not specified, the default value will be used.
	SystemDisk *SystemDisk `json:"SystemDisk,omitnil,omitempty" name:"SystemDisk"`

	// Data disk configuration of the instance. If the parameter is not specified, no data disk will be purchased by default. If you want to purchase data disks, you can specify 21 data disks, including up to 1 `LOCAL_BASIC` data disk or `LOCAL_SSD` data disk and up to 20 `CLOUD_BASIC` data disks, `CLOUD_PREMIUM` data disks, or `CLOUD_SSD` data disks.
	DataDisks []*DataDisk `json:"DataDisks,omitnil,omitempty" name:"DataDisks"`

	// VPC configurations (VPC ID, subnet ID, etc). If It's not specified, the classic network will be used by default. If a VPC IP is specified in this parameter, the `InstanceCount` can only be 1.
	VirtualPrivateCloud *VirtualPrivateCloud `json:"VirtualPrivateCloud,omitnil,omitempty" name:"VirtualPrivateCloud"`

	// Configuration of public network bandwidth. If it's not specified, 0 Mbps is used by default.
	InternetAccessible *InternetAccessible `json:"InternetAccessible,omitnil,omitempty" name:"InternetAccessible"`

	// Number of instances to purchase. Value range: 1 (default) to 100. It cannot exceed the remaining CVM quota of the user. For more information on quota, see [Restrictions on CVM Instance Purchase](https://intl.cloud.tencent.com/document/product/213/2664).
	InstanceCount *int64 `json:"InstanceCount,omitnil,omitempty" name:"InstanceCount"`

	// Instance display name. <li>if no instance display name is specified, it will display 'unnamed' by default.</li> <li>when purchasing multiple instances, if the pattern string `{R:x}` is specified, it means generating numbers `[x, x+n-1]`, where `n` represents the number of purchased instances. for example, `server_{R:3}`: when purchasing 1 instance, the instance display name is `server_3`; when purchasing 2 instances, the instance display names are `server_3` and `server_4` respectively. supports specifying multiple pattern strings `{R:x}`.</li> <li>when purchasing multiple instances without specifying a pattern string, suffixes `1, 2...n` will be added to the instance display name, where `n` represents the number of purchased instances. for example, `server_`: when purchasing 2 instances, the instance display names are `server_1` and `server_2` respectively.</li> <li>supports up to 128 characters (including pattern strings).</li>.
	InstanceName *string `json:"InstanceName,omitnil,omitempty" name:"InstanceName"`

	// Login settings of the instance. You can use this parameter to set the login method, password, and key of the instance, or keep the original login settings of the image. By default, a random password will be generated and sent to you via the Message Center.
	LoginSettings *LoginSettings `json:"LoginSettings,omitnil,omitempty" name:"LoginSettings"`

	// Security group to which an instance belongs. obtain this parameter by calling the `SecurityGroupId` field in the return value of [DescribeSecurityGroups](https://www.tencentcloud.com/document/product/215/15808). if not specified, bind the default security group under the designated project. if the default security group does not exist, automatically create it.
	SecurityGroupIds []*string `json:"SecurityGroupIds,omitnil,omitempty" name:"SecurityGroupIds"`

	// Enhanced services. You can use this parameter to specify whether to enable services such as Cloud Security and Cloud Monitor. If this parameter is not specified, Cloud Monitor and Cloud Security will be enabled by default.
	EnhancedService *EnhancedService `json:"EnhancedService,omitnil,omitempty" name:"EnhancedService"`

	// A unique string supplied by the client to ensure that the request is idempotent. Its maximum length is 64 ASCII characters. If this parameter is not specified, the idempotency of the request cannot be guaranteed.<br>For more information, see Ensuring Idempotency.
	ClientToken *string `json:"ClientToken,omitnil,omitempty" name:"ClientToken"`

	// Hostname of Cloud Virtual Machine.<br><li>Period (.) and hyphen (-) should not be used as the first or last character of the hostname, and should not be used consecutively.</li><br><li>Windows instances: The hostname should contain 2 to 15 characters, including letters (case insensitive), digits, and hyphens (-), does not support periods (.), and should not be all digits.</li><br><li>Instances of other types (such as Linux instances): The hostname should contain 2 to 30 characters, including multiple periods (.), with each segment between periods considered as one section. Each section can contain letters (case insensitive), digits, and hyphens (-).</li>
	HostName *string `json:"HostName,omitnil,omitempty" name:"HostName"`

	// The tag description list. This parameter is used to bind a tag to a resource instance. A tag can only be bound to CVM instances.
	TagSpecification []*TagSpecification `json:"TagSpecification,omitnil,omitempty" name:"TagSpecification"`

	// The market options of the instance.
	InstanceMarketOptions *InstanceMarketOptionsRequest `json:"InstanceMarketOptions,omitnil,omitempty" name:"InstanceMarketOptions"`

	// Custom metadata, supports creating key-value pairs of custom metadata when creating a CVM.
	// 
	// **Note: this field is in beta test.**.
	Metadata *Metadata `json:"Metadata,omitnil,omitempty" name:"Metadata"`

	// HPC cluster ID.
	HpcClusterId *string `json:"HpcClusterId,omitnil,omitempty" name:"HpcClusterId"`

	// Information about the CPU topology of an instance. If not specified, it is determined by system resources.
	CpuTopology *CpuTopology `json:"CpuTopology,omitnil,omitempty" name:"CpuTopology"`

	LaunchTemplate *LaunchTemplate `json:"LaunchTemplate,omitnil,omitempty" name:"LaunchTemplate"`
}

func (r *InquiryPriceRunInstancesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *InquiryPriceRunInstancesRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Placement")
	delete(f, "ImageId")
	delete(f, "InstanceChargeType")
	delete(f, "InstanceChargePrepaid")
	delete(f, "InstanceType")
	delete(f, "SystemDisk")
	delete(f, "DataDisks")
	delete(f, "VirtualPrivateCloud")
	delete(f, "InternetAccessible")
	delete(f, "InstanceCount")
	delete(f, "InstanceName")
	delete(f, "LoginSettings")
	delete(f, "SecurityGroupIds")
	delete(f, "EnhancedService")
	delete(f, "ClientToken")
	delete(f, "HostName")
	delete(f, "TagSpecification")
	delete(f, "InstanceMarketOptions")
	delete(f, "Metadata")
	delete(f, "HpcClusterId")
	delete(f, "CpuTopology")
	delete(f, "LaunchTemplate")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "InquiryPriceRunInstancesRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type InquiryPriceRunInstancesResponseParams struct {
	// Price of the instance with the specified configurations.
	Price *Price `json:"Price,omitnil,omitempty" name:"Price"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type InquiryPriceRunInstancesResponse struct {
	*tchttp.BaseResponse
	Response *InquiryPriceRunInstancesResponseParams `json:"Response"`
}

func (r *InquiryPriceRunInstancesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *InquiryPriceRunInstancesResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type Instance struct {
	// Location of the instance
	Placement *Placement `json:"Placement,omitnil,omitempty" name:"Placement"`

	// Instance `ID`
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// Instance model
	InstanceType *string `json:"InstanceType,omitnil,omitempty" name:"InstanceType"`

	// Number of CPU cores of the instance; unit: core
	CPU *int64 `json:"CPU,omitnil,omitempty" name:"CPU"`

	// Instance memory capacity. unit: GiB.
	Memory *int64 `json:"Memory,omitnil,omitempty" name:"Memory"`

	// Instance business status. valid values:<br><li>NORMAL: indicates an instance in the NORMAL state</li><li>EXPIRED: indicates an EXPIRED instance</li><li>PROTECTIVELY_ISOLATED: PROTECTIVELY ISOLATED instance</li>.
	RestrictState *string `json:"RestrictState,omitnil,omitempty" name:"RestrictState"`

	// Instance name
	InstanceName *string `json:"InstanceName,omitnil,omitempty" name:"InstanceName"`

	// Instance billing plan. Valid values:<br><li>`POSTPAID_BY_HOUR`: pay after use. You are billed by the hour, by traffic.<br><li>`CDHPAID`: `CDH` billing plan. Applicable to `CDH` only, not the instances on the host.<br>
	InstanceChargeType *string `json:"InstanceChargeType,omitnil,omitempty" name:"InstanceChargeType"`

	// Information on the system disk of the instance
	SystemDisk *SystemDisk `json:"SystemDisk,omitnil,omitempty" name:"SystemDisk"`

	// Information of the instance data disks.
	DataDisks []*DataDisk `json:"DataDisks,omitnil,omitempty" name:"DataDisks"`

	// List of private IPs of the instance's primary ENI.
	PrivateIpAddresses []*string `json:"PrivateIpAddresses,omitnil,omitempty" name:"PrivateIpAddresses"`

	// List of public IPs of the instance's primary ENI.
	// Note: This field may return null, indicating that no valid value is found.
	PublicIpAddresses []*string `json:"PublicIpAddresses,omitnil,omitempty" name:"PublicIpAddresses"`

	// Information on instance bandwidth.
	InternetAccessible *InternetAccessible `json:"InternetAccessible,omitnil,omitempty" name:"InternetAccessible"`

	// Information on the VPC where the instance resides.
	VirtualPrivateCloud *VirtualPrivateCloud `json:"VirtualPrivateCloud,omitnil,omitempty" name:"VirtualPrivateCloud"`

	// `ID` of the image used to create the instance.
	ImageId *string `json:"ImageId,omitnil,omitempty" name:"ImageId"`

	// AUTO-Renewal flag. valid values:<br><li>`NOTIFY_AND_MANUAL_RENEW`: indicates that a notification of impending expiration is made but AUTO-renewal is not performed</li><li>`NOTIFY_AND_AUTO_RENEW`: indicates that a notification of impending expiration is made AND AUTO-renewal is performed</li><li>`DISABLE_NOTIFY_AND_MANUAL_RENEW`: indicates that notification that it is about to expire is not made AND AUTO-renewal is not performed.
	// Note: this field is null in postpaid mode.
	RenewFlag *string `json:"RenewFlag,omitnil,omitempty" name:"RenewFlag"`

	// Creation time following the `ISO8601` standard and using `UTC` time in the format of `YYYY-MM-DDThh:mm:ssZ`.
	CreatedTime *string `json:"CreatedTime,omitnil,omitempty" name:"CreatedTime"`

	// Expiration time in UTC format following the `ISO8601` standard: `YYYY-MM-DDThh:mm:ssZ`. Note: this parameter is `null` for postpaid instances.
	ExpiredTime *string `json:"ExpiredTime,omitnil,omitempty" name:"ExpiredTime"`

	// Operating system name.
	OsName *string `json:"OsName,omitnil,omitempty" name:"OsName"`

	// Security groups to which the instance belongs. To obtain the security group IDs, you can call [DescribeSecurityGroups](https://intl.cloud.tencent.com/document/api/215/15808) and look for the `sgld` fields in the response.
	SecurityGroupIds []*string `json:"SecurityGroupIds,omitnil,omitempty" name:"SecurityGroupIds"`

	// Login settings of the instance. Currently only the key associated with the instance is returned.
	LoginSettings *LoginSettings `json:"LoginSettings,omitnil,omitempty" name:"LoginSettings"`

	// Instance status. for specific status types, see the  [instance status table](https://www.tencentcloud.com/document/product/213/15753#instancestatus)
	InstanceState *string `json:"InstanceState,omitnil,omitempty" name:"InstanceState"`

	// List of tags associated with the instance.
	Tags []*Tag `json:"Tags,omitnil,omitempty" name:"Tags"`

	// Shutdown billing mode of an instance.
	// 
	// Valid values: <br><li>KEEP_CHARGING: billing continues after shutdown</li><li>STOP_CHARGING: billing stops after shutdown</li><li>NOT_APPLICABLE: the instance is NOT shut down or stopping billing after shutdown is NOT APPLICABLE to the instance</li>.
	StopChargingMode *string `json:"StopChargingMode,omitnil,omitempty" name:"StopChargingMode"`

	// Globally unique ID of the instance.
	Uuid *string `json:"Uuid,omitnil,omitempty" name:"Uuid"`

	// Last operation of the instance, such as StopInstances or ResetInstance.
	LatestOperation *string `json:"LatestOperation,omitnil,omitempty" name:"LatestOperation"`

	// The latest operation status of the instance. valid values:<br><li>SUCCESS: indicates that the operation succeeded</li><li>OPERATING: indicates that the operation is in progress</li><li>FAILED: indicates that the operation FAILED</li>.
	// Note: This field may return null, indicating that no valid value is found.
	LatestOperationState *string `json:"LatestOperationState,omitnil,omitempty" name:"LatestOperationState"`

	// Unique request ID for the last operation of the instance.
	LatestOperationRequestId *string `json:"LatestOperationRequestId,omitnil,omitempty" name:"LatestOperationRequestId"`

	// Spread placement group ID.
	DisasterRecoverGroupId *string `json:"DisasterRecoverGroupId,omitnil,omitempty" name:"DisasterRecoverGroupId"`

	// IPv6 address of the instance.
	// Note: this field may return null, indicating that no valid value is obtained.
	IPv6Addresses []*string `json:"IPv6Addresses,omitnil,omitempty" name:"IPv6Addresses"`

	// CAM role name.
	// Note: this field may return null, indicating that no valid value is obtained.
	CamRoleName *string `json:"CamRoleName,omitnil,omitempty" name:"CamRoleName"`

	// High-performance computing cluster ID.
	HpcClusterId *string `json:"HpcClusterId,omitnil,omitempty" name:"HpcClusterId"`

	// IP list of HPC cluster.
	// Note: this field may return null, indicating that no valid value was found.
	RdmaIpAddresses []*string `json:"RdmaIpAddresses,omitnil,omitempty" name:"RdmaIpAddresses"`

	// Dedicated cluster ID where the instance is located.
	DedicatedClusterId *string `json:"DedicatedClusterId,omitnil,omitempty" name:"DedicatedClusterId"`

	// Instance isolation type. valid values:<br><li>ARREAR: indicates arrears isolation<br></li><li>EXPIRE: indicates isolation upon expiration<br></li><li>MANMADE: indicates voluntarily refunded isolation<br></li><li>NOTISOLATED: indicates unisolated<br></li>.
	IsolatedSource *string `json:"IsolatedSource,omitnil,omitempty" name:"IsolatedSource"`

	// GPU information. if it is a gpu type instance, this value will be returned. for other type instances, it does not return.
	GPUInfo *GPUInfo `json:"GPUInfo,omitnil,omitempty" name:"GPUInfo"`

	// Instance OS license type. Default value: `TencentCloud`
	LicenseType *string `json:"LicenseType,omitnil,omitempty" name:"LicenseType"`

	// Instance destruction protection flag indicates whether an instance is allowed to be deleted through an api. value ranges from:<br><li>true: indicates that instance protection is enabled, deletion through apis is not allowed</li><li>false: indicates that instance protection is disabled, allow passage</li><br>default value: false.
	DisableApiTermination *bool `json:"DisableApiTermination,omitnil,omitempty" name:"DisableApiTermination"`

	// Default login user
	DefaultLoginUser *string `json:"DefaultLoginUser,omitnil,omitempty" name:"DefaultLoginUser"`

	// Default login port
	DefaultLoginPort *int64 `json:"DefaultLoginPort,omitnil,omitempty" name:"DefaultLoginPort"`

	// Latest operation errors of the instance.
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	LatestOperationErrorMsg *string `json:"LatestOperationErrorMsg,omitnil,omitempty" name:"LatestOperationErrorMsg"`

	// Custom metadata. this parameter corresponds to the metadata information specified when creating a CVM. **note: in beta test**.
	Metadata *Metadata `json:"Metadata,omitnil,omitempty" name:"Metadata"`

	// Specifies the public IPv6 address bound to the instance.
	PublicIPv6Addresses []*string `json:"PublicIPv6Addresses,omitnil,omitempty" name:"PublicIPv6Addresses"`
}

type InstanceAttribute struct {
	// Instance ID.
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// Instance attribute information.
	Attributes *Attribute `json:"Attributes,omitnil,omitempty" name:"Attributes"`
}

type InstanceChargePrepaid struct {
	// Subscription period in months. value range: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36, 48, 60.
	Period *int64 `json:"Period,omitnil,omitempty" name:"Period"`

	// AUTO-Renewal flag. value ranges:<br><li>NOTIFY_AND_AUTO_RENEW: NOTIFY of expiration AND automatically RENEW.</li><br><li>NOTIFY_AND_MANUAL_RENEW: NOTIFY of expiration but do not automatically RENEW.</li><br><li>DISABLE_NOTIFY_AND_MANUAL_RENEW: do not NOTIFY of expiration AND do not automatically RENEW.</li><br><br>default value: NOTIFY_AND_MANUAL_RENEW. if this parameter is set to NOTIFY_AND_AUTO_RENEW, the instance will be automatically renewed monthly after expiration, provided that the account balance is sufficient.
	RenewFlag *string `json:"RenewFlag,omitnil,omitempty" name:"RenewFlag"`
}

type InstanceFamilyConfig struct {
	// Full name of the model family.
	InstanceFamilyName *string `json:"InstanceFamilyName,omitnil,omitempty" name:"InstanceFamilyName"`

	// Acronym of the model family.
	InstanceFamily *string `json:"InstanceFamily,omitnil,omitempty" name:"InstanceFamily"`
}

type InstanceMarketOptionsRequest struct {
	// Relevant options for spot instances.
	SpotOptions *SpotMarketOptions `json:"SpotOptions,omitnil,omitempty" name:"SpotOptions"`

	// Market option type. The value can only be spot currently.
	MarketType *string `json:"MarketType,omitnil,omitempty" name:"MarketType"`
}

type InstanceStatus struct {
	// Instance `ID`.
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// Instance status. Valid values:<br><li>PENDING: Creating.<br></li><li>LAUNCH_FAILED: Creation failed.<br></li><li>RUNNING: Running.<br></li><li>STOPPED: Shut down.<br></li><li>STARTING: Starting up.<br></li><li>STOPPING: Shutting down.<br></li><li>REBOOTING: Restarting.<br></li><li>SHUTDOWN: Shut down and to be terminated.<br></li><li>TERMINATING: Terminating.<br></li><li>ENTER_RESCUE_MODE: Entering the rescue mode.<br></li><li>RESCUE_MODE: In the rescue mode.<br></li><li>EXIT_RESCUE_MODE: Exiting the rescue mode.<br></li><li>ENTER_SERVICE_LIVE_MIGRATE: Entering online service migration.<br></li><li>SERVICE_LIVE_MIGRATE: In online service migration.<br></li><li>EXIT_SERVICE_LIVE_MIGRATE: Exiting online service migration.<br></li>
	InstanceState *string `json:"InstanceState,omitnil,omitempty" name:"InstanceState"`
}

type InstanceTypeQuotaItem struct {
	// Availability zone.
	Zone *string `json:"Zone,omitnil,omitempty" name:"Zone"`

	// Instance model.
	InstanceType *string `json:"InstanceType,omitnil,omitempty" name:"InstanceType"`

	// Instance billing plan. Valid values: <br><li>POSTPAID_BY_HOUR: pay after use. You are billed for your traffic by the hour.<br><li>`CDHPAID`: [`CDH`](https://intl.cloud.tencent.com/document/product/416?from_cn_redirect=1) billing plan. Applicable to `CDH` only, not the instances on the host.
	InstanceChargeType *string `json:"InstanceChargeType,omitnil,omitempty" name:"InstanceChargeType"`

	// ENI type. For example, 25 represents an ENI of 25 GB.
	NetworkCard *int64 `json:"NetworkCard,omitnil,omitempty" name:"NetworkCard"`

	// Additional data.
	// Note: This field may return null, indicating that no valid value is found.
	Externals *Externals `json:"Externals,omitnil,omitempty" name:"Externals"`

	// Number of CPU cores of an instance model.
	Cpu *int64 `json:"Cpu,omitnil,omitempty" name:"Cpu"`

	// Instance memory capacity; unit: `GB`.
	Memory *int64 `json:"Memory,omitnil,omitempty" name:"Memory"`

	// Instance model family.
	InstanceFamily *string `json:"InstanceFamily,omitnil,omitempty" name:"InstanceFamily"`

	// Model name.
	TypeName *string `json:"TypeName,omitnil,omitempty" name:"TypeName"`

	// List of local disk specifications. If the parameter returns null, it means that local disks cannot be created.
	LocalDiskTypeList []*LocalDiskType `json:"LocalDiskTypeList,omitnil,omitempty" name:"LocalDiskTypeList"`

	// Whether an instance is for sale. Valid values:<br><li>SELL: The instance is available for purchase.<br></li>SOLD_OUT: The instance has been sold out.
	Status *string `json:"Status,omitnil,omitempty" name:"Status"`

	// Price of an instance model.
	Price *ItemPrice `json:"Price,omitnil,omitempty" name:"Price"`

	// Details of out-of-stock items
	// Note: this field may return null, indicating that no valid value is obtained.
	SoldOutReason *string `json:"SoldOutReason,omitnil,omitempty" name:"SoldOutReason"`

	// Private network bandwidth, in Gbps.
	InstanceBandwidth *float64 `json:"InstanceBandwidth,omitnil,omitempty" name:"InstanceBandwidth"`

	// The max packet sending and receiving capability (in 10k PPS).
	InstancePps *int64 `json:"InstancePps,omitnil,omitempty" name:"InstancePps"`

	// Number of local storage blocks.
	StorageBlockAmount *int64 `json:"StorageBlockAmount,omitnil,omitempty" name:"StorageBlockAmount"`

	// CPU type.
	CpuType *string `json:"CpuType,omitnil,omitempty" name:"CpuType"`

	// Number of GPUs of the instance.
	Gpu *int64 `json:"Gpu,omitnil,omitempty" name:"Gpu"`

	// Number of FPGAs of the instance.
	Fpga *int64 `json:"Fpga,omitnil,omitempty" name:"Fpga"`

	// Descriptive information of the instance.
	Remark *string `json:"Remark,omitnil,omitempty" name:"Remark"`


	GpuCount *float64 `json:"GpuCount,omitnil,omitempty" name:"GpuCount"`

	// CPU clock rate of the instance
	Frequency *string `json:"Frequency,omitnil,omitempty" name:"Frequency"`

	// Inventory status. Valid values:
	// <li>EnoughStock: Inventory is sufficient.</li> 
	// <li>NormalStock: Supply is guaranteed.</li>
	// <li>UnderStock: Inventory is about to sell out.</li> 
	// <li>WithoutStock: Inventory is already sold out.</li>
	// Note: This field may return null, indicating that no valid value is found.
	StatusCategory *string `json:"StatusCategory,omitnil,omitempty" name:"StatusCategory"`
}

type InternetAccessible struct {
	// Network connection billing plan. Valid value:
	// 
	// <li>TRAFFIC_POSTPAID_BY_HOUR: pay after use. You are billed for your traffic, by the hour. </li>
	// <li>BANDWIDTH_PACKAGE: Bandwidth package user. </li>
	InternetChargeType *string `json:"InternetChargeType,omitnil,omitempty" name:"InternetChargeType"`

	// The maximum outbound bandwidth of the public network, in Mbps. The default value is 0 Mbps. The upper limit of bandwidth varies for different models. For more information, see [Purchase Network Bandwidth](https://intl.cloud.tencent.com/document/product/213/12523?from_cn_redirect=1).
	InternetMaxBandwidthOut *int64 `json:"InternetMaxBandwidthOut,omitnil,omitempty" name:"InternetMaxBandwidthOut"`

	// Whether to allocate a public IP address. Valid values:<br><li>true: Allocate a public IP address.</li><li>false: Do not allocate a public IP address.</li><br>When the public network bandwidth is greater than 0 Mbps, you can choose whether to enable the public IP address. The public IP address is enabled by default. When the public network bandwidth is 0, allocating the public IP address is not supported. This parameter is only used as an input parameter in the RunInstances API.
	PublicIpAssigned *bool `json:"PublicIpAssigned,omitnil,omitempty" name:"PublicIpAssigned"`

	// Bandwidth package ID. it can be obtained through the `BandwidthPackageId` in the return value from the DescribeBandwidthPackages api. this parameter is used as an input parameter only in the RunInstances api.
	BandwidthPackageId *string `json:"BandwidthPackageId,omitnil,omitempty" name:"BandwidthPackageId"`

	// Line type. for details on various types of lines and supported regions, refer to [EIP IP address types](https://cloud.tencent.com/document/product/1199/41646). default value: BGP.
	// <Li>BGP: specifies the general bgp line.</li>.
	// For a user with static single-line IP allowlist enabled, valid values include:.
	// <Li>CMCC: china mobile.</li>.
	// <Li>CTCC: china telecom</li>.
	// <Li>CUCC: china unicom</li>.
	// Note: The static single-line IP is only supported in some regions.
	// 
	InternetServiceProvider *string `json:"InternetServiceProvider,omitnil,omitempty" name:"InternetServiceProvider"`

	// Specifies the public IP type.
	// 
	// <Li>WanIP: specifies the public ip address.</li>.
	// <Li>HighQualityEIP: specifies a high quality ip. high quality ip is only supported in Singapore and hong kong (china).</li>.
	// <li> AntiDDoSEIP: specifies the anti-ddos eip. only partial regions support anti-ddos eip. details visible in the [elastic IP product overview](https://www.tencentcloud.comom/document/product/1199/41646?from_cn_redirect=1).</li>.
	// If needed, assign a public IPv4 address to the resource by specifying the IPv4 address type.
	// 
	// This feature is in beta test in selected regions. submit a ticket for consultation (https://console.cloud.tencent.com/workorder/category) if needed.
	IPv4AddressType *string `json:"IPv4AddressType,omitnil,omitempty" name:"IPv4AddressType"`

	// Indicates the type of elastic public IPv6.
	// <Li>EIPv6: elastic ip version 6.</li>.
	// <Li>HighQualityEIPv6: specifies the high quality ipv6. highqualityeipv6 is only supported in hong kong (china).</li>.
	// If needed, assign an elastic IPv6 address for resource allocation.
	// 
	// This feature is in beta test in selected regions. submit a ticket for consultation (https://console.cloud.tencent.com/workorder/category) if needed.
	IPv6AddressType *string `json:"IPv6AddressType,omitnil,omitempty" name:"IPv6AddressType"`

	// DDoS protection package unique ID. this field is required when applying for a ddos protection IP.
	AntiDDoSPackageId *string `json:"AntiDDoSPackageId,omitnil,omitempty" name:"AntiDDoSPackageId"`
}

type InternetChargeTypeConfig struct {
	// Network billing method.
	InternetChargeType *string `json:"InternetChargeType,omitnil,omitempty" name:"InternetChargeType"`

	// Description of the network billing method.
	Description *string `json:"Description,omitnil,omitempty" name:"Description"`
}

type ItemPrice struct {
	// Original price of subsequent total costs, postpaid billing mode usage, unit: usd. <li>if other time interval items are returned, such as UnitPriceSecondStep, this item represents the time interval (0, 96) hr. if no other time interval items are returned, this item represents the full period (0, ∞) hr.
	UnitPrice *float64 `json:"UnitPrice,omitnil,omitempty" name:"UnitPrice"`

	// Billing unit for pay-as-you-go mode. valid values: <br><li>HOUR: billed on an hourly basis. it's used for hourly POSTPAID instances (`POSTPAID_BY_HOUR`). <br><li>GB: bill BY TRAFFIC in GB. it's used for POSTPAID products that are billed BY the hourly TRAFFIC (`TRAFFIC_POSTPAID_BY_HOUR`).
	ChargeUnit *string `json:"ChargeUnit,omitnil,omitempty" name:"ChargeUnit"`

	// Original price of total prepaid costs. measurement unit: usd.
	OriginalPrice *float64 `json:"OriginalPrice,omitnil,omitempty" name:"OriginalPrice"`

	// Discount price of total prepaid costs. unit: usd.
	DiscountPrice *float64 `json:"DiscountPrice,omitnil,omitempty" name:"DiscountPrice"`

	// Discount, such as 20.0 representing 80% off.
	Discount *float64 `json:"Discount,omitnil,omitempty" name:"Discount"`

	// Discounted price of subsequent total cost, postpaid billing mode usage, unit: usd <li>if other time interval items are returned, such as UnitPriceDiscountSecondStep, this item represents the time interval (0, 96) hr; if no other time interval items are returned, this item represents the full period (0, ∞) hr.
	UnitPriceDiscount *float64 `json:"UnitPriceDiscount,omitnil,omitempty" name:"UnitPriceDiscount"`

	// Original price of subsequent total costs for usage time range (96, 360) hr in postpaid billing mode. unit: usd.
	UnitPriceSecondStep *float64 `json:"UnitPriceSecondStep,omitnil,omitempty" name:"UnitPriceSecondStep"`

	// Discounted price of subsequent total cost for usage time interval (96, 360) hr in postpaid billing mode. unit: usd.
	UnitPriceDiscountSecondStep *float64 `json:"UnitPriceDiscountSecondStep,omitnil,omitempty" name:"UnitPriceDiscountSecondStep"`

	// Specifies the original price of subsequent total costs with a usage time interval exceeding 360 hr in postpaid billing mode. measurement unit: usd.
	UnitPriceThirdStep *float64 `json:"UnitPriceThirdStep,omitnil,omitempty" name:"UnitPriceThirdStep"`

	// Discounted price of subsequent total cost for usage time interval exceeding 360 hr in postpaid billing mode. measurement unit: usd.
	UnitPriceDiscountThirdStep *float64 `json:"UnitPriceDiscountThirdStep,omitnil,omitempty" name:"UnitPriceDiscountThirdStep"`

	// Specifies the original price of total 3-year prepaid costs in prepaid billing mode. measurement unit: usd.
	// Note: This field may return null, indicating that no valid value is found.
	OriginalPriceThreeYear *float64 `json:"OriginalPriceThreeYear,omitnil,omitempty" name:"OriginalPriceThreeYear"`

	// Specifies the discount price for an advance payment of the total fee for three years, prepaid mode usage, measurement unit: usd.
	// Note: This field may return null, indicating that no valid value is found.
	DiscountPriceThreeYear *float64 `json:"DiscountPriceThreeYear,omitnil,omitempty" name:"DiscountPriceThreeYear"`

	// Specifies the discount for a 3-year advance payment, for example 20.0 represents 80% off.
	// Note: This field may return null, indicating that no valid value is found.
	DiscountThreeYear *float64 `json:"DiscountThreeYear,omitnil,omitempty" name:"DiscountThreeYear"`

	// Specifies the original price of the 5-year total cost with advance payment, using prepaid billing mode. measurement unit: usd.
	// Note: This field may return null, indicating that no valid value is found.
	OriginalPriceFiveYear *float64 `json:"OriginalPriceFiveYear,omitnil,omitempty" name:"OriginalPriceFiveYear"`

	// Prepaid 5-year total cost discount price, prepaid billing mode usage. unit: usd.
	// Note: This field may return null, indicating that no valid value is found.
	DiscountPriceFiveYear *float64 `json:"DiscountPriceFiveYear,omitnil,omitempty" name:"DiscountPriceFiveYear"`

	// Specifies the discount for 5-year advance payment, such as 20.0 for 80% off.
	// Note: This field may return null, indicating that no valid value is found.
	DiscountFiveYear *float64 `json:"DiscountFiveYear,omitnil,omitempty" name:"DiscountFiveYear"`

	// Original price of one-year advance payment total cost. prepaid mode usage. unit: usd.
	// Note: This field may return null, indicating that no valid value is found.
	OriginalPriceOneYear *float64 `json:"OriginalPriceOneYear,omitnil,omitempty" name:"OriginalPriceOneYear"`

	// Discount price for total advance payment for one year. specifies prepaid mode usage. measurement unit: usd.
	// Note: This field may return null, indicating that no valid value is found.
	DiscountPriceOneYear *float64 `json:"DiscountPriceOneYear,omitnil,omitempty" name:"DiscountPriceOneYear"`

	// Specifies the discount for a one-year advance payment, such as 20.0 for 80% off.
	// Note: This field may return null, indicating that no valid value is found.
	DiscountOneYear *float64 `json:"DiscountOneYear,omitnil,omitempty" name:"DiscountOneYear"`
}

type KeyPair struct {
	// Key pair `ID`, the unique identifier of a key pair.
	KeyId *string `json:"KeyId,omitnil,omitempty" name:"KeyId"`

	// Key pair name.
	KeyName *string `json:"KeyName,omitnil,omitempty" name:"KeyName"`

	// `ID` of the project to which a key pair belongs.
	ProjectId *int64 `json:"ProjectId,omitnil,omitempty" name:"ProjectId"`

	// Key pair description.
	Description *string `json:"Description,omitnil,omitempty" name:"Description"`

	// Content of public key in a key pair.
	PublicKey *string `json:"PublicKey,omitnil,omitempty" name:"PublicKey"`

	// Content of private key in a key pair. Tencent Cloud do not keep private keys. Please keep it properly.
	PrivateKey *string `json:"PrivateKey,omitnil,omitempty" name:"PrivateKey"`

	// `ID` list of instances associated with a key.
	AssociatedInstanceIds []*string `json:"AssociatedInstanceIds,omitnil,omitempty" name:"AssociatedInstanceIds"`

	// Creation time, which follows the `ISO8601` standard and uses `UTC` time in the format of `YYYY-MM-DDThh:mm:ssZ`.
	CreatedTime *string `json:"CreatedTime,omitnil,omitempty" name:"CreatedTime"`

	// The list of tags bound to the key.
	// Note: This field may return `null`, indicating that no valid value can be obtained.
	Tags []*Tag `json:"Tags,omitnil,omitempty" name:"Tags"`
}

type LaunchTemplate struct {
	// Instance launch template ID. This parameter enables you to create an instance using the preset parameters in the template.
	LaunchTemplateId *string `json:"LaunchTemplateId,omitnil,omitempty" name:"LaunchTemplateId"`

	// Instance launch template version number. If specified, this parameter will be used to create a new instance launch template.
	LaunchTemplateVersion *uint64 `json:"LaunchTemplateVersion,omitnil,omitempty" name:"LaunchTemplateVersion"`
}

type LaunchTemplateInfo struct {
	// Instance launch template version number.
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	LatestVersionNumber *uint64 `json:"LatestVersionNumber,omitnil,omitempty" name:"LatestVersionNumber"`

	// Instance launch template ID.
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	LaunchTemplateId *string `json:"LaunchTemplateId,omitnil,omitempty" name:"LaunchTemplateId"`

	// Instance launch template name.
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	LaunchTemplateName *string `json:"LaunchTemplateName,omitnil,omitempty" name:"LaunchTemplateName"`

	// Default instance launch template version number.
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	DefaultVersionNumber *uint64 `json:"DefaultVersionNumber,omitnil,omitempty" name:"DefaultVersionNumber"`

	// Total number of versions that the instance launch template contains.
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	LaunchTemplateVersionCount *uint64 `json:"LaunchTemplateVersionCount,omitnil,omitempty" name:"LaunchTemplateVersionCount"`

	// UIN of the user who created the template.
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	CreatedBy *string `json:"CreatedBy,omitnil,omitempty" name:"CreatedBy"`

	// Creation time of the template.
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	CreationTime *string `json:"CreationTime,omitnil,omitempty" name:"CreationTime"`
}

type LaunchTemplateVersionData struct {
	// Location of the instance.
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	Placement *Placement `json:"Placement,omitnil,omitempty" name:"Placement"`

	// Instance model.
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	InstanceType *string `json:"InstanceType,omitnil,omitempty" name:"InstanceType"`

	// Instance name.
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	InstanceName *string `json:"InstanceName,omitnil,omitempty" name:"InstanceName"`

	// Instance billing mode. Valid values: <br><li>`POSTPAID_BY_HOUR`: postpaid for pay-as-you-go instances <br><li>`CDHPAID`: billed for CDH instances, not the CVMs running on the CDHs.<br>
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	InstanceChargeType *string `json:"InstanceChargeType,omitnil,omitempty" name:"InstanceChargeType"`

	// Instance system disk information.
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	SystemDisk *SystemDisk `json:"SystemDisk,omitnil,omitempty" name:"SystemDisk"`

	// Instance data disk information. This parameter only covers the data disks purchased together with the instance.
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	DataDisks []*DataDisk `json:"DataDisks,omitnil,omitempty" name:"DataDisks"`

	// Instance bandwidth information.
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	InternetAccessible *InternetAccessible `json:"InternetAccessible,omitnil,omitempty" name:"InternetAccessible"`

	// Information of the VPC where the instance resides.
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	VirtualPrivateCloud *VirtualPrivateCloud `json:"VirtualPrivateCloud,omitnil,omitempty" name:"VirtualPrivateCloud"`

	// `ID` of the image used to create the instance.
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	ImageId *string `json:"ImageId,omitnil,omitempty" name:"ImageId"`

	// Security groups to which the instance belongs. To obtain the security group IDs, you can call [DescribeSecurityGroups](https://intl.cloud.tencent.com/document/api/215/15808) and look for the `sgld` fields in the response.
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	SecurityGroupIds []*string `json:"SecurityGroupIds,omitnil,omitempty" name:"SecurityGroupIds"`

	// Login settings of the instance. Currently, only the key associated with the instance is returned.
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	LoginSettings *LoginSettings `json:"LoginSettings,omitnil,omitempty" name:"LoginSettings"`

	// CAM role name.
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	CamRoleName *string `json:"CamRoleName,omitnil,omitempty" name:"CamRoleName"`

	// HPC cluster `ID`.
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	HpcClusterId *string `json:"HpcClusterId,omitnil,omitempty" name:"HpcClusterId"`

	// Number of instances purchased.
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	InstanceCount *uint64 `json:"InstanceCount,omitnil,omitempty" name:"InstanceCount"`

	// Enhanced service.
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	EnhancedService *EnhancedService `json:"EnhancedService,omitnil,omitempty" name:"EnhancedService"`

	// User data provided to the instance. This parameter needs to be encoded in base64 format with the maximum size of 16KB.
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	UserData *string `json:"UserData,omitnil,omitempty" name:"UserData"`

	// Placement group ID. You can only specify one.
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	DisasterRecoverGroupIds []*string `json:"DisasterRecoverGroupIds,omitnil,omitempty" name:"DisasterRecoverGroupIds"`

	// Scheduled tasks. You can use this parameter to specify scheduled tasks for the instance. Only scheduled termination is supported.
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	ActionTimer *ActionTimer `json:"ActionTimer,omitnil,omitempty" name:"ActionTimer"`

	// Market options of the instance, such as parameters related to spot instances. This parameter is required for spot instances.
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	InstanceMarketOptions *InstanceMarketOptionsRequest `json:"InstanceMarketOptions,omitnil,omitempty" name:"InstanceMarketOptions"`

	// Hostname of a CVM.
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	HostName *string `json:"HostName,omitnil,omitempty" name:"HostName"`

	// A string used to ensure the idempotency of the request.
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	ClientToken *string `json:"ClientToken,omitnil,omitempty" name:"ClientToken"`

	// Prepaid mode. This parameter indicates relevant parameter settings for monthly-subscribed instances.
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	InstanceChargePrepaid *InstanceChargePrepaid `json:"InstanceChargePrepaid,omitnil,omitempty" name:"InstanceChargePrepaid"`

	// List of tag description. By specifying this parameter, the tag can be bound to the corresponding CVM and CBS instances at the same time.
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	TagSpecification []*TagSpecification `json:"TagSpecification,omitnil,omitempty" name:"TagSpecification"`

	// Whether to enable termination protection. Valid values:
	// 
	// TRUE: Termination protection is enabled.
	// FALSE: Termination protection is disabled.
	// 
	// Default value: `FALSE`.
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	DisableApiTermination *bool `json:"DisableApiTermination,omitnil,omitempty" name:"DisableApiTermination"`
}

type LaunchTemplateVersionInfo struct {
	// Instance launch template version number.
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	LaunchTemplateVersion *uint64 `json:"LaunchTemplateVersion,omitnil,omitempty" name:"LaunchTemplateVersion"`

	// Details of instance launch template versions.
	LaunchTemplateVersionData *LaunchTemplateVersionData `json:"LaunchTemplateVersionData,omitnil,omitempty" name:"LaunchTemplateVersionData"`

	// Creation time of the instance launch template version.
	CreationTime *string `json:"CreationTime,omitnil,omitempty" name:"CreationTime"`

	// Instance launch template ID.
	LaunchTemplateId *string `json:"LaunchTemplateId,omitnil,omitempty" name:"LaunchTemplateId"`

	// Specifies whether it’s the default launch template version.
	IsDefaultVersion *bool `json:"IsDefaultVersion,omitnil,omitempty" name:"IsDefaultVersion"`

	// Information of instance launch template version description.
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	LaunchTemplateVersionDescription *string `json:"LaunchTemplateVersionDescription,omitnil,omitempty" name:"LaunchTemplateVersionDescription"`

	// Creator account
	CreatedBy *string `json:"CreatedBy,omitnil,omitempty" name:"CreatedBy"`
}

type LocalDiskType struct {
	// Type of a local disk.
	Type *string `json:"Type,omitnil,omitempty" name:"Type"`

	// Attributes of a local disk.
	PartitionType *string `json:"PartitionType,omitnil,omitempty" name:"PartitionType"`

	// Minimum size of a local disk.
	MinSize *int64 `json:"MinSize,omitnil,omitempty" name:"MinSize"`

	// Maximum size of a local disk.
	MaxSize *int64 `json:"MaxSize,omitnil,omitempty" name:"MaxSize"`

	// Whether a local disk is required during purchase. Valid values:<br><li>REQUIRED: required<br><li>OPTIONAL: optional
	Required *string `json:"Required,omitnil,omitempty" name:"Required"`
}

type LoginSettings struct {
	// Instance login password. The password complexity limits vary with the operating system type as follows: <br><li>The Linux instance password must be 8 to 30 characters long and include at least two of the following: [a-z], [A-Z], [0-9], and special characters of [( ) \` ~ ! @ # $ % ^ & * - + = | { } [ ] : ; ' , . ? / ]. <br><li>The Windows instance password must be 12 to 30 characters long and include at least three of the following: [a-z], [A-Z], [0-9], and special characters [( ) \` ~ ! @ # $ % ^ & * - + = | { } [ ] : ; ' , . ? /]. <br><br>If this parameter is not specified, you need to set it before login by using the console to "reset password" or by calling the ResetInstancesPassword API.
	Password *string `json:"Password,omitnil,omitempty" name:"Password"`

	// List of key IDs. After an instance is associated with a key, you can access the instance with the private key in the key pair. You can call [`DescribeKeyPairs`](https://intl.cloud.tencent.com/document/api/213/15699?from_cn_redirect=1) to obtain `KeyId`. You cannot specify a key and a password at the same time. Windows instances do not support keys.
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	KeyIds []*string `json:"KeyIds,omitnil,omitempty" name:"KeyIds"`

	// Retain the original settings of the image. this parameter cannot be specified simultaneously with Password or KeyIds.N. it can be set to true only when an instance is created with a custom image, shared image, or externally imported image. value ranges from true to false: <li>true: indicates that the login settings of the image are retained</li><li>false: indicates that the login settings of the image are not retained</li>. default value: false.
	KeepImageLogin *string `json:"KeepImageLogin,omitnil,omitempty" name:"KeepImageLogin"`
}

type Metadata struct {
	// A list of custom metadata key-value pairs.
	Items []*MetadataItem `json:"Items,omitnil,omitempty" name:"Items"`
}

type MetadataItem struct {
	// Custom metadata key, consisting of uppercase letters (A-Z), lowercase letters (A-Z), digits (0-9), underscores (_), or hyphens (-), with a size limit of 128 bytes.
	Key *string `json:"Key,omitnil,omitempty" name:"Key"`

	// Custom metadata value. The upper limit of message size is 256 KB.
	Value *string `json:"Value,omitnil,omitempty" name:"Value"`
}

// Predefined struct for user
type ModifyChcAttributeRequestParams struct {
	// CHC host IDs
	ChcIds []*string `json:"ChcIds,omitnil,omitempty" name:"ChcIds"`

	// CHC host name
	InstanceName *string `json:"InstanceName,omitnil,omitempty" name:"InstanceName"`

	// Server type
	DeviceType *string `json:"DeviceType,omitnil,omitempty" name:"DeviceType"`

	// Valid characters: Letters, numbers, hyphens and underscores
	BmcUser *string `json:"BmcUser,omitnil,omitempty" name:"BmcUser"`

	// The password can contain 8 to 16 characters, including letters, numbers and special symbols (()`~!@#$%^&*-+=_|{}).
	Password *string `json:"Password,omitnil,omitempty" name:"Password"`

	// BMC network security group list
	BmcSecurityGroupIds []*string `json:"BmcSecurityGroupIds,omitnil,omitempty" name:"BmcSecurityGroupIds"`
}

type ModifyChcAttributeRequest struct {
	*tchttp.BaseRequest
	
	// CHC host IDs
	ChcIds []*string `json:"ChcIds,omitnil,omitempty" name:"ChcIds"`

	// CHC host name
	InstanceName *string `json:"InstanceName,omitnil,omitempty" name:"InstanceName"`

	// Server type
	DeviceType *string `json:"DeviceType,omitnil,omitempty" name:"DeviceType"`

	// Valid characters: Letters, numbers, hyphens and underscores
	BmcUser *string `json:"BmcUser,omitnil,omitempty" name:"BmcUser"`

	// The password can contain 8 to 16 characters, including letters, numbers and special symbols (()`~!@#$%^&*-+=_|{}).
	Password *string `json:"Password,omitnil,omitempty" name:"Password"`

	// BMC network security group list
	BmcSecurityGroupIds []*string `json:"BmcSecurityGroupIds,omitnil,omitempty" name:"BmcSecurityGroupIds"`
}

func (r *ModifyChcAttributeRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyChcAttributeRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ChcIds")
	delete(f, "InstanceName")
	delete(f, "DeviceType")
	delete(f, "BmcUser")
	delete(f, "Password")
	delete(f, "BmcSecurityGroupIds")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ModifyChcAttributeRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyChcAttributeResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ModifyChcAttributeResponse struct {
	*tchttp.BaseResponse
	Response *ModifyChcAttributeResponseParams `json:"Response"`
}

func (r *ModifyChcAttributeResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyChcAttributeResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyDisasterRecoverGroupAttributeRequestParams struct {
	// Spread placement group ID, which can be obtained by calling the [DescribeDisasterRecoverGroups](https://intl.cloud.tencent.com/document/api/213/17810?from_cn_redirect=1) API.
	DisasterRecoverGroupId *string `json:"DisasterRecoverGroupId,omitnil,omitempty" name:"DisasterRecoverGroupId"`

	// Name of a spread placement group. The name must be 1-60 characters long and can contain both Chinese characters and English letters.
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`
}

type ModifyDisasterRecoverGroupAttributeRequest struct {
	*tchttp.BaseRequest
	
	// Spread placement group ID, which can be obtained by calling the [DescribeDisasterRecoverGroups](https://intl.cloud.tencent.com/document/api/213/17810?from_cn_redirect=1) API.
	DisasterRecoverGroupId *string `json:"DisasterRecoverGroupId,omitnil,omitempty" name:"DisasterRecoverGroupId"`

	// Name of a spread placement group. The name must be 1-60 characters long and can contain both Chinese characters and English letters.
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`
}

func (r *ModifyDisasterRecoverGroupAttributeRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyDisasterRecoverGroupAttributeRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "DisasterRecoverGroupId")
	delete(f, "Name")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ModifyDisasterRecoverGroupAttributeRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyDisasterRecoverGroupAttributeResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ModifyDisasterRecoverGroupAttributeResponse struct {
	*tchttp.BaseResponse
	Response *ModifyDisasterRecoverGroupAttributeResponseParams `json:"Response"`
}

func (r *ModifyDisasterRecoverGroupAttributeResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyDisasterRecoverGroupAttributeResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyHostsAttributeRequestParams struct {
	// CDH instance ID(s).
	HostIds []*string `json:"HostIds,omitnil,omitempty" name:"HostIds"`

	// CDH instance name to be displayed. You can specify any name you like, but its length cannot exceed 60 characters.
	HostName *string `json:"HostName,omitnil,omitempty" name:"HostName"`

	// Auto renewal flag. Valid values: <br><li>NOTIFY_AND_AUTO_RENEW: notify upon expiration and renew automatically <br><li>NOTIFY_AND_MANUAL_RENEW: notify upon expiration but do not renew automatically <br><li>DISABLE_NOTIFY_AND_MANUAL_RENEW: neither notify upon expiration nor renew automatically <br><br>If this parameter is specified as NOTIFY_AND_AUTO_RENEW, the instance will be automatically renewed on a monthly basis if the account balance is sufficient.
	RenewFlag *string `json:"RenewFlag,omitnil,omitempty" name:"RenewFlag"`

	// Project ID. You can use the `AddProject` API to create projects, and obtain the `projectId` field in the response of the `DescribeProject` API. When using the [DescribeHosts](https://intl.cloud.tencent.com/document/api/213/16474?from_cn_redirect=1) API to query instances later, you can filter the results by the project ID.
	ProjectId *uint64 `json:"ProjectId,omitnil,omitempty" name:"ProjectId"`
}

type ModifyHostsAttributeRequest struct {
	*tchttp.BaseRequest
	
	// CDH instance ID(s).
	HostIds []*string `json:"HostIds,omitnil,omitempty" name:"HostIds"`

	// CDH instance name to be displayed. You can specify any name you like, but its length cannot exceed 60 characters.
	HostName *string `json:"HostName,omitnil,omitempty" name:"HostName"`

	// Auto renewal flag. Valid values: <br><li>NOTIFY_AND_AUTO_RENEW: notify upon expiration and renew automatically <br><li>NOTIFY_AND_MANUAL_RENEW: notify upon expiration but do not renew automatically <br><li>DISABLE_NOTIFY_AND_MANUAL_RENEW: neither notify upon expiration nor renew automatically <br><br>If this parameter is specified as NOTIFY_AND_AUTO_RENEW, the instance will be automatically renewed on a monthly basis if the account balance is sufficient.
	RenewFlag *string `json:"RenewFlag,omitnil,omitempty" name:"RenewFlag"`

	// Project ID. You can use the `AddProject` API to create projects, and obtain the `projectId` field in the response of the `DescribeProject` API. When using the [DescribeHosts](https://intl.cloud.tencent.com/document/api/213/16474?from_cn_redirect=1) API to query instances later, you can filter the results by the project ID.
	ProjectId *uint64 `json:"ProjectId,omitnil,omitempty" name:"ProjectId"`
}

func (r *ModifyHostsAttributeRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyHostsAttributeRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "HostIds")
	delete(f, "HostName")
	delete(f, "RenewFlag")
	delete(f, "ProjectId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ModifyHostsAttributeRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyHostsAttributeResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ModifyHostsAttributeResponse struct {
	*tchttp.BaseResponse
	Response *ModifyHostsAttributeResponseParams `json:"Response"`
}

func (r *ModifyHostsAttributeResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyHostsAttributeResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyImageAttributeRequestParams struct {
	// Image ID, such as `img-gvbnzy6f`. You can obtain the image ID in the following ways:<li>Call the [DescribeImages](https://intl.cloud.tencent.com/document/api/213/15715?from_cn_redirect=1) API and find the value of `ImageId` in the response.</li><li>Obtain it in the [Image console](https://console.cloud.tencent.com/cvm/image).</li>
	ImageId *string `json:"ImageId,omitnil,omitempty" name:"ImageId"`

	// New image name, which should meet the following requirements:<li>It should not exceed 60 characters.</li><li>It should be unique.</li>
	ImageName *string `json:"ImageName,omitnil,omitempty" name:"ImageName"`

	// New image description, which should meet the following requirement:<li>It should not exceed 256 characters.</li>
	ImageDescription *string `json:"ImageDescription,omitnil,omitempty" name:"ImageDescription"`

	// Sets the image family;
	ImageFamily *string `json:"ImageFamily,omitnil,omitempty" name:"ImageFamily"`

	// Sets whether the image is deprecated;
	ImageDeprecated *bool `json:"ImageDeprecated,omitnil,omitempty" name:"ImageDeprecated"`
}

type ModifyImageAttributeRequest struct {
	*tchttp.BaseRequest
	
	// Image ID, such as `img-gvbnzy6f`. You can obtain the image ID in the following ways:<li>Call the [DescribeImages](https://intl.cloud.tencent.com/document/api/213/15715?from_cn_redirect=1) API and find the value of `ImageId` in the response.</li><li>Obtain it in the [Image console](https://console.cloud.tencent.com/cvm/image).</li>
	ImageId *string `json:"ImageId,omitnil,omitempty" name:"ImageId"`

	// New image name, which should meet the following requirements:<li>It should not exceed 60 characters.</li><li>It should be unique.</li>
	ImageName *string `json:"ImageName,omitnil,omitempty" name:"ImageName"`

	// New image description, which should meet the following requirement:<li>It should not exceed 256 characters.</li>
	ImageDescription *string `json:"ImageDescription,omitnil,omitempty" name:"ImageDescription"`

	// Sets the image family;
	ImageFamily *string `json:"ImageFamily,omitnil,omitempty" name:"ImageFamily"`

	// Sets whether the image is deprecated;
	ImageDeprecated *bool `json:"ImageDeprecated,omitnil,omitempty" name:"ImageDeprecated"`
}

func (r *ModifyImageAttributeRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyImageAttributeRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ImageId")
	delete(f, "ImageName")
	delete(f, "ImageDescription")
	delete(f, "ImageFamily")
	delete(f, "ImageDeprecated")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ModifyImageAttributeRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyImageAttributeResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ModifyImageAttributeResponse struct {
	*tchttp.BaseResponse
	Response *ModifyImageAttributeResponseParams `json:"Response"`
}

func (r *ModifyImageAttributeResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyImageAttributeResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyImageSharePermissionRequestParams struct {
	// Image ID, such as `img-gvbnzy6f`. You can obtain the image ID in the following ways:<br><li>Call the [DescribeImages](https://intl.cloud.tencent.com/document/api/213/15715?from_cn_redirect=1) API and find the value of `ImageId` in the response.</li><br><li>Obtain it in the [Image console](https://console.cloud.tencent.com/cvm/image).</li><br>The image ID should correspond to an image in the `NORMAL` state. For more information on image status, see [Image Data Table](https://intl.cloud.tencent.com/document/product/213/15753?from_cn_redirect=1#Image).
	ImageId *string `json:"ImageId,omitnil,omitempty" name:"ImageId"`

	// ID list of root accounts receiving shared images. For the format of array-type parameters, see [API Introduction](https://intl.cloud.tencent.com/document/api/213/568?from_cn_redirect=1). An account ID is different from a QQ number. For details on root account IDs, refer to the account ID section in [Account Information](https://console.cloud.tencent.com/developer).
	AccountIds []*string `json:"AccountIds,omitnil,omitempty" name:"AccountIds"`

	// Operations. Valid values: `SHARE`, sharing an image; `CANCEL`, cancelling an image sharing. 
	Permission *string `json:"Permission,omitnil,omitempty" name:"Permission"`
}

type ModifyImageSharePermissionRequest struct {
	*tchttp.BaseRequest
	
	// Image ID, such as `img-gvbnzy6f`. You can obtain the image ID in the following ways:<br><li>Call the [DescribeImages](https://intl.cloud.tencent.com/document/api/213/15715?from_cn_redirect=1) API and find the value of `ImageId` in the response.</li><br><li>Obtain it in the [Image console](https://console.cloud.tencent.com/cvm/image).</li><br>The image ID should correspond to an image in the `NORMAL` state. For more information on image status, see [Image Data Table](https://intl.cloud.tencent.com/document/product/213/15753?from_cn_redirect=1#Image).
	ImageId *string `json:"ImageId,omitnil,omitempty" name:"ImageId"`

	// ID list of root accounts receiving shared images. For the format of array-type parameters, see [API Introduction](https://intl.cloud.tencent.com/document/api/213/568?from_cn_redirect=1). An account ID is different from a QQ number. For details on root account IDs, refer to the account ID section in [Account Information](https://console.cloud.tencent.com/developer).
	AccountIds []*string `json:"AccountIds,omitnil,omitempty" name:"AccountIds"`

	// Operations. Valid values: `SHARE`, sharing an image; `CANCEL`, cancelling an image sharing. 
	Permission *string `json:"Permission,omitnil,omitempty" name:"Permission"`
}

func (r *ModifyImageSharePermissionRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyImageSharePermissionRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ImageId")
	delete(f, "AccountIds")
	delete(f, "Permission")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ModifyImageSharePermissionRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyImageSharePermissionResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ModifyImageSharePermissionResponse struct {
	*tchttp.BaseResponse
	Response *ModifyImageSharePermissionResponseParams `json:"Response"`
}

func (r *ModifyImageSharePermissionResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyImageSharePermissionResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyInstancesAttributeRequestParams struct {
	// Instance ID(s). To obtain the instance IDs, you can call [DescribeInstances](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1) and look for `InstanceId` in the response. The maximum number of instances in each request is 100.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// Modified instance name. can be named as required but should not exceed 128 characters.
	InstanceName *string `json:"InstanceName,omitnil,omitempty" name:"InstanceName"`

	// User data provided to an instance, which needs to be encoded in Base64 format with a maximum size of 16 KB. For details on obtaining this parameter, refer to the startup commands for [Windows](https://intl.cloud.tencent.com/document/product/213/17526?from_cn_redirect=1) and [Linux](https://intl.cloud.tencent.com/document/product/213/17525?from_cn_redirect=1).
	UserData *string `json:"UserData,omitnil,omitempty" name:"UserData"`

	// Specifies the security group Id list of the specified instance after modification. the instance will reassociate with the security groups in the specified list, and the associated security group will be unbound.
	SecurityGroups []*string `json:"SecurityGroups,omitnil,omitempty" name:"SecurityGroups"`

	// The role bound with the instance. If it is not specified, it indicates to unbind the current role of the CVM.
	CamRoleName *string `json:"CamRoleName,omitnil,omitempty" name:"CamRoleName"`

	// Modified hostname of an instance.<li>Period (.) and hyphen (-) should not be used as the first or last character of the hostname, and should not be used consecutively.</li><li>Windows instances: The hostname should contain 2 to 15 characters, including letters (case insensitive), digits, and hyphens (-), does not support periods (.), and should not be all digits.</li><li>Instances of other types (such as Linux instances): The hostname should contain 2 to 60 characters, including multiple periods (.), with each segment between periods considered as one section. Each section can contain letters (case insensitive), digits, and hyphens (-).</li>Note: After the hostname is modified, the instance will restart immediately, and the new hostname will take effect after the restart.
	HostName *string `json:"HostName,omitnil,omitempty" name:"HostName"`

	// Instance termination protection flag, indicating whether an instance is allowed to be deleted through an API. Valid values:<li>true: Instance protection is enabled, and the instance is not allowed to be deleted through the API.</li><li>false: Instance protection is disabled, and the instance is allowed to be deleted through the API.</li>Default value: false.
	DisableApiTermination *bool `json:"DisableApiTermination,omitnil,omitempty" name:"DisableApiTermination"`

	// Role type, used in conjunction with CamRoleName. this value can be obtained from the RoleType field in the API response of CAM [DescribeRoleList](https://www.tencentcloud.com/document/product/1219/67889) or [GetRole](https://www.tencentcloud.com/document/product/598/33557?lang=en). currently, only user, system, and service_linked types are accepted.
	// For example, when CamRoleName contains "LinkedRoleIn" (such as TKE_QCSLinkedRoleInPrometheusService), DescribeRoleList and GetRole return RoleType as service_linked, this parameter must also transmit service_linked.
	// The parameter default value is user. this parameter can be omitted if CameRoleName is not of the service_linked kind.
	CamRoleType *string `json:"CamRoleType,omitnil,omitempty" name:"CamRoleType"`

	// Whether to automatically restart an instance when modifying a hostname. If not specified, the instance will automatically restart by default.
	// - true: Modify the hostname and automatically restart the instance.
	// - false: Modify the hostname without automatically restarting the instance. A manual restart is required for the new hostname to take effect.
	// Note: This parameter is valid only when a hostname is modified.
	AutoReboot *bool `json:"AutoReboot,omitnil,omitempty" name:"AutoReboot"`
}

type ModifyInstancesAttributeRequest struct {
	*tchttp.BaseRequest
	
	// Instance ID(s). To obtain the instance IDs, you can call [DescribeInstances](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1) and look for `InstanceId` in the response. The maximum number of instances in each request is 100.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// Modified instance name. can be named as required but should not exceed 128 characters.
	InstanceName *string `json:"InstanceName,omitnil,omitempty" name:"InstanceName"`

	// User data provided to an instance, which needs to be encoded in Base64 format with a maximum size of 16 KB. For details on obtaining this parameter, refer to the startup commands for [Windows](https://intl.cloud.tencent.com/document/product/213/17526?from_cn_redirect=1) and [Linux](https://intl.cloud.tencent.com/document/product/213/17525?from_cn_redirect=1).
	UserData *string `json:"UserData,omitnil,omitempty" name:"UserData"`

	// Specifies the security group Id list of the specified instance after modification. the instance will reassociate with the security groups in the specified list, and the associated security group will be unbound.
	SecurityGroups []*string `json:"SecurityGroups,omitnil,omitempty" name:"SecurityGroups"`

	// The role bound with the instance. If it is not specified, it indicates to unbind the current role of the CVM.
	CamRoleName *string `json:"CamRoleName,omitnil,omitempty" name:"CamRoleName"`

	// Modified hostname of an instance.<li>Period (.) and hyphen (-) should not be used as the first or last character of the hostname, and should not be used consecutively.</li><li>Windows instances: The hostname should contain 2 to 15 characters, including letters (case insensitive), digits, and hyphens (-), does not support periods (.), and should not be all digits.</li><li>Instances of other types (such as Linux instances): The hostname should contain 2 to 60 characters, including multiple periods (.), with each segment between periods considered as one section. Each section can contain letters (case insensitive), digits, and hyphens (-).</li>Note: After the hostname is modified, the instance will restart immediately, and the new hostname will take effect after the restart.
	HostName *string `json:"HostName,omitnil,omitempty" name:"HostName"`

	// Instance termination protection flag, indicating whether an instance is allowed to be deleted through an API. Valid values:<li>true: Instance protection is enabled, and the instance is not allowed to be deleted through the API.</li><li>false: Instance protection is disabled, and the instance is allowed to be deleted through the API.</li>Default value: false.
	DisableApiTermination *bool `json:"DisableApiTermination,omitnil,omitempty" name:"DisableApiTermination"`

	// Role type, used in conjunction with CamRoleName. this value can be obtained from the RoleType field in the API response of CAM [DescribeRoleList](https://www.tencentcloud.com/document/product/1219/67889) or [GetRole](https://www.tencentcloud.com/document/product/598/33557?lang=en). currently, only user, system, and service_linked types are accepted.
	// For example, when CamRoleName contains "LinkedRoleIn" (such as TKE_QCSLinkedRoleInPrometheusService), DescribeRoleList and GetRole return RoleType as service_linked, this parameter must also transmit service_linked.
	// The parameter default value is user. this parameter can be omitted if CameRoleName is not of the service_linked kind.
	CamRoleType *string `json:"CamRoleType,omitnil,omitempty" name:"CamRoleType"`

	// Whether to automatically restart an instance when modifying a hostname. If not specified, the instance will automatically restart by default.
	// - true: Modify the hostname and automatically restart the instance.
	// - false: Modify the hostname without automatically restarting the instance. A manual restart is required for the new hostname to take effect.
	// Note: This parameter is valid only when a hostname is modified.
	AutoReboot *bool `json:"AutoReboot,omitnil,omitempty" name:"AutoReboot"`
}

func (r *ModifyInstancesAttributeRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyInstancesAttributeRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceIds")
	delete(f, "InstanceName")
	delete(f, "UserData")
	delete(f, "SecurityGroups")
	delete(f, "CamRoleName")
	delete(f, "HostName")
	delete(f, "DisableApiTermination")
	delete(f, "CamRoleType")
	delete(f, "AutoReboot")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ModifyInstancesAttributeRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyInstancesAttributeResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ModifyInstancesAttributeResponse struct {
	*tchttp.BaseResponse
	Response *ModifyInstancesAttributeResponseParams `json:"Response"`
}

func (r *ModifyInstancesAttributeResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyInstancesAttributeResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyInstancesChargeTypeRequestParams struct {
	// One or more instance ids to be operated. you can obtain the instance ID through the `InstanceId` in the return value from the api [DescribeInstances](https://www.tencentcloud.com/document/api/213/15728?from_cn_redirect=1). the maximum number of instances per request is 20.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// Modified instance [billing type](https://www.tencentcloud.com/document/product/213/2180?from_cn_redirect=1). <li>`PREPAID`: monthly subscription.</li> 
	// **Note:** Only supports converting pay-as-you-go instances to annual and monthly subscription instances.
	// 
	// default value: `PREPAID`
	InstanceChargeType *string `json:"InstanceChargeType,omitnil,omitempty" name:"InstanceChargeType"`

	// Modified prepaid mode, parameter settings related to monthly/annual subscription. by specifying this parameter, you can specify the purchase duration of annual and monthly subscription instances, whether to enable auto-renewal, and other attributes. 
	InstanceChargePrepaid *InstanceChargePrepaid `json:"InstanceChargePrepaid,omitnil,omitempty" name:"InstanceChargePrepaid"`

	// Whether to switch the billing mode of elastic data cloud disks simultaneously. valid values: <li> true: means switching the billing mode of elastic data cloud disks</li><li> false: means not switching the billing mode of elastic data cloud disks</li>default value: false.
	ModifyPortableDataDisk *bool `json:"ModifyPortableDataDisk,omitnil,omitempty" name:"ModifyPortableDataDisk"`
}

type ModifyInstancesChargeTypeRequest struct {
	*tchttp.BaseRequest
	
	// One or more instance ids to be operated. you can obtain the instance ID through the `InstanceId` in the return value from the api [DescribeInstances](https://www.tencentcloud.com/document/api/213/15728?from_cn_redirect=1). the maximum number of instances per request is 20.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// Modified instance [billing type](https://www.tencentcloud.com/document/product/213/2180?from_cn_redirect=1). <li>`PREPAID`: monthly subscription.</li> 
	// **Note:** Only supports converting pay-as-you-go instances to annual and monthly subscription instances.
	// 
	// default value: `PREPAID`
	InstanceChargeType *string `json:"InstanceChargeType,omitnil,omitempty" name:"InstanceChargeType"`

	// Modified prepaid mode, parameter settings related to monthly/annual subscription. by specifying this parameter, you can specify the purchase duration of annual and monthly subscription instances, whether to enable auto-renewal, and other attributes. 
	InstanceChargePrepaid *InstanceChargePrepaid `json:"InstanceChargePrepaid,omitnil,omitempty" name:"InstanceChargePrepaid"`

	// Whether to switch the billing mode of elastic data cloud disks simultaneously. valid values: <li> true: means switching the billing mode of elastic data cloud disks</li><li> false: means not switching the billing mode of elastic data cloud disks</li>default value: false.
	ModifyPortableDataDisk *bool `json:"ModifyPortableDataDisk,omitnil,omitempty" name:"ModifyPortableDataDisk"`
}

func (r *ModifyInstancesChargeTypeRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyInstancesChargeTypeRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceIds")
	delete(f, "InstanceChargeType")
	delete(f, "InstanceChargePrepaid")
	delete(f, "ModifyPortableDataDisk")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ModifyInstancesChargeTypeRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyInstancesChargeTypeResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ModifyInstancesChargeTypeResponse struct {
	*tchttp.BaseResponse
	Response *ModifyInstancesChargeTypeResponseParams `json:"Response"`
}

func (r *ModifyInstancesChargeTypeResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyInstancesChargeTypeResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyInstancesDisasterRecoverGroupRequestParams struct {
	// One or more instance ids to be operated. you can obtain the instance ID through the `InstanceId` in the return value from the api [DescribeInstances](https://www.tencentcloud.com/zh/document/api/213/33258). the maximum number of instances per request is 100.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// Spread placement group ID. obtain through the api [DescribeDisasterRecoverGroups](https://www.tencentcloud.com/zh/document/api/213/33261).
	DisasterRecoverGroupId *string `json:"DisasterRecoverGroupId,omitnil,omitempty" name:"DisasterRecoverGroupId"`

	// Whether to forcibly change instance hosts. value range:<br><li>true: indicates that the instance is allowed to change hosts, allowing the instance to reboot. local disk machine does not support specifying this parameter.</li><br><li>false: does not allow the instance to change hosts. instances can only be added to the placement group on the current host. this may result in a failure to change the placement group.</li><br><br>default value: false.
	Force *bool `json:"Force,omitnil,omitempty" name:"Force"`
}

type ModifyInstancesDisasterRecoverGroupRequest struct {
	*tchttp.BaseRequest
	
	// One or more instance ids to be operated. you can obtain the instance ID through the `InstanceId` in the return value from the api [DescribeInstances](https://www.tencentcloud.com/zh/document/api/213/33258). the maximum number of instances per request is 100.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// Spread placement group ID. obtain through the api [DescribeDisasterRecoverGroups](https://www.tencentcloud.com/zh/document/api/213/33261).
	DisasterRecoverGroupId *string `json:"DisasterRecoverGroupId,omitnil,omitempty" name:"DisasterRecoverGroupId"`

	// Whether to forcibly change instance hosts. value range:<br><li>true: indicates that the instance is allowed to change hosts, allowing the instance to reboot. local disk machine does not support specifying this parameter.</li><br><li>false: does not allow the instance to change hosts. instances can only be added to the placement group on the current host. this may result in a failure to change the placement group.</li><br><br>default value: false.
	Force *bool `json:"Force,omitnil,omitempty" name:"Force"`
}

func (r *ModifyInstancesDisasterRecoverGroupRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyInstancesDisasterRecoverGroupRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceIds")
	delete(f, "DisasterRecoverGroupId")
	delete(f, "Force")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ModifyInstancesDisasterRecoverGroupRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyInstancesDisasterRecoverGroupResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ModifyInstancesDisasterRecoverGroupResponse struct {
	*tchttp.BaseResponse
	Response *ModifyInstancesDisasterRecoverGroupResponseParams `json:"Response"`
}

func (r *ModifyInstancesDisasterRecoverGroupResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyInstancesDisasterRecoverGroupResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyInstancesProjectRequestParams struct {
	// Instance IDs. To obtain the instance IDs, you can call [`DescribeInstances`](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1) and look for `InstanceId` in the response. You can operate up to 100 instances in each request.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// Project ID. You can use the API `AddProject` to create projects, and obtain the `projectId` field in the response of the `DescribeProject` API. When using the [DescribeInstances](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1) API to query instances later, you can filter the results by the project ID.
	ProjectId *int64 `json:"ProjectId,omitnil,omitempty" name:"ProjectId"`
}

type ModifyInstancesProjectRequest struct {
	*tchttp.BaseRequest
	
	// Instance IDs. To obtain the instance IDs, you can call [`DescribeInstances`](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1) and look for `InstanceId` in the response. You can operate up to 100 instances in each request.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// Project ID. You can use the API `AddProject` to create projects, and obtain the `projectId` field in the response of the `DescribeProject` API. When using the [DescribeInstances](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1) API to query instances later, you can filter the results by the project ID.
	ProjectId *int64 `json:"ProjectId,omitnil,omitempty" name:"ProjectId"`
}

func (r *ModifyInstancesProjectRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyInstancesProjectRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceIds")
	delete(f, "ProjectId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ModifyInstancesProjectRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyInstancesProjectResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ModifyInstancesProjectResponse struct {
	*tchttp.BaseResponse
	Response *ModifyInstancesProjectResponseParams `json:"Response"`
}

func (r *ModifyInstancesProjectResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyInstancesProjectResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyInstancesRenewFlagRequestParams struct {
	// For one or more instance IDs to be operated. You can obtain the instance ID through the `instanceid` in the returned value from the API [DescribeInstances](https://cloud.tencent.com/document/api/213/15728). The maximum number of instances that can be operated for each request is 100.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// Auto-renewal flag. Valid values: <br><li>NOTIFY_AND_AUTO_RENEW: Notifies of expiration and performs auto-renewal.</li><li>NOTIFY_AND_MANUAL_RENEW: Notifies of expiration but does not perform auto-renewal.</li><li>DISABLE_NOTIFY_AND_MANUAL_RENEW: Not notifies of expiration nor perform auto-renewal.</li><br>If this parameter is specified to NOTIFY_AND_AUTO_RENEW, the instance will be automatically renewed on a monthly basis after it expires when there is sufficient account balance.
	RenewFlag *string `json:"RenewFlag,omitnil,omitempty" name:"RenewFlag"`
}

type ModifyInstancesRenewFlagRequest struct {
	*tchttp.BaseRequest
	
	// For one or more instance IDs to be operated. You can obtain the instance ID through the `instanceid` in the returned value from the API [DescribeInstances](https://cloud.tencent.com/document/api/213/15728). The maximum number of instances that can be operated for each request is 100.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// Auto-renewal flag. Valid values: <br><li>NOTIFY_AND_AUTO_RENEW: Notifies of expiration and performs auto-renewal.</li><li>NOTIFY_AND_MANUAL_RENEW: Notifies of expiration but does not perform auto-renewal.</li><li>DISABLE_NOTIFY_AND_MANUAL_RENEW: Not notifies of expiration nor perform auto-renewal.</li><br>If this parameter is specified to NOTIFY_AND_AUTO_RENEW, the instance will be automatically renewed on a monthly basis after it expires when there is sufficient account balance.
	RenewFlag *string `json:"RenewFlag,omitnil,omitempty" name:"RenewFlag"`
}

func (r *ModifyInstancesRenewFlagRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyInstancesRenewFlagRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceIds")
	delete(f, "RenewFlag")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ModifyInstancesRenewFlagRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyInstancesRenewFlagResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ModifyInstancesRenewFlagResponse struct {
	*tchttp.BaseResponse
	Response *ModifyInstancesRenewFlagResponseParams `json:"Response"`
}

func (r *ModifyInstancesRenewFlagResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyInstancesRenewFlagResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyInstancesVpcAttributeRequestParams struct {
	// Instance ID(s). To obtain the instance IDs, you can call [`DescribeInstances`](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1) and look for `InstanceId` in the response.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// VPC configurations. You can use this parameter to specify the VPC ID, subnet ID, VPC IP, etc. If the specified VPC ID and subnet ID (the subnet must be in the same availability zone as the instance) are different from the VPC where the specified instance resides, the instance will be migrated to a subnet of the specified VPC. You can use `PrivateIpAddresses` to specify the VPC subnet IP. If you want to specify the subnet IP, you will need to specify a subnet IP for each of the specified instances, and each `InstanceIds` will match a `PrivateIpAddresses`. If `PrivateIpAddresses` is not specified, the VPC subnet IP will be assigned randomly.
	VirtualPrivateCloud *VirtualPrivateCloud `json:"VirtualPrivateCloud,omitnil,omitempty" name:"VirtualPrivateCloud"`

	// Whether to force shut down a running instances. Default value: TRUE.
	ForceStop *bool `json:"ForceStop,omitnil,omitempty" name:"ForceStop"`

	// Whether to keep the host name. Default value: FALSE.
	ReserveHostName *bool `json:"ReserveHostName,omitnil,omitempty" name:"ReserveHostName"`
}

type ModifyInstancesVpcAttributeRequest struct {
	*tchttp.BaseRequest
	
	// Instance ID(s). To obtain the instance IDs, you can call [`DescribeInstances`](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1) and look for `InstanceId` in the response.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// VPC configurations. You can use this parameter to specify the VPC ID, subnet ID, VPC IP, etc. If the specified VPC ID and subnet ID (the subnet must be in the same availability zone as the instance) are different from the VPC where the specified instance resides, the instance will be migrated to a subnet of the specified VPC. You can use `PrivateIpAddresses` to specify the VPC subnet IP. If you want to specify the subnet IP, you will need to specify a subnet IP for each of the specified instances, and each `InstanceIds` will match a `PrivateIpAddresses`. If `PrivateIpAddresses` is not specified, the VPC subnet IP will be assigned randomly.
	VirtualPrivateCloud *VirtualPrivateCloud `json:"VirtualPrivateCloud,omitnil,omitempty" name:"VirtualPrivateCloud"`

	// Whether to force shut down a running instances. Default value: TRUE.
	ForceStop *bool `json:"ForceStop,omitnil,omitempty" name:"ForceStop"`

	// Whether to keep the host name. Default value: FALSE.
	ReserveHostName *bool `json:"ReserveHostName,omitnil,omitempty" name:"ReserveHostName"`
}

func (r *ModifyInstancesVpcAttributeRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyInstancesVpcAttributeRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceIds")
	delete(f, "VirtualPrivateCloud")
	delete(f, "ForceStop")
	delete(f, "ReserveHostName")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ModifyInstancesVpcAttributeRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyInstancesVpcAttributeResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ModifyInstancesVpcAttributeResponse struct {
	*tchttp.BaseResponse
	Response *ModifyInstancesVpcAttributeResponseParams `json:"Response"`
}

func (r *ModifyInstancesVpcAttributeResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyInstancesVpcAttributeResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyKeyPairAttributeRequestParams struct {
	// Key pair ID in the format of `skey-xxxxxxxx`. <br><br>You can obtain the available key pair IDs in two ways: <br><li>Log in to the [console](https://console.cloud.tencent.com/cvm/sshkey) to query the key pair IDs. <br><li>Call [DescribeKeyPairs](https://intl.cloud.tencent.com/document/api/213/15699?from_cn_redirect=1) and look for `KeyId` in the response.
	KeyId *string `json:"KeyId,omitnil,omitempty" name:"KeyId"`

	// New key pair name, which can contain numbers, letters, and underscores, with a maximum length of 25 characters.
	KeyName *string `json:"KeyName,omitnil,omitempty" name:"KeyName"`

	// New key pair description. You can specify any name you like, but its length cannot exceed 60 characters.
	Description *string `json:"Description,omitnil,omitempty" name:"Description"`
}

type ModifyKeyPairAttributeRequest struct {
	*tchttp.BaseRequest
	
	// Key pair ID in the format of `skey-xxxxxxxx`. <br><br>You can obtain the available key pair IDs in two ways: <br><li>Log in to the [console](https://console.cloud.tencent.com/cvm/sshkey) to query the key pair IDs. <br><li>Call [DescribeKeyPairs](https://intl.cloud.tencent.com/document/api/213/15699?from_cn_redirect=1) and look for `KeyId` in the response.
	KeyId *string `json:"KeyId,omitnil,omitempty" name:"KeyId"`

	// New key pair name, which can contain numbers, letters, and underscores, with a maximum length of 25 characters.
	KeyName *string `json:"KeyName,omitnil,omitempty" name:"KeyName"`

	// New key pair description. You can specify any name you like, but its length cannot exceed 60 characters.
	Description *string `json:"Description,omitnil,omitempty" name:"Description"`
}

func (r *ModifyKeyPairAttributeRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyKeyPairAttributeRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "KeyId")
	delete(f, "KeyName")
	delete(f, "Description")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ModifyKeyPairAttributeRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyKeyPairAttributeResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ModifyKeyPairAttributeResponse struct {
	*tchttp.BaseResponse
	Response *ModifyKeyPairAttributeResponseParams `json:"Response"`
}

func (r *ModifyKeyPairAttributeResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyKeyPairAttributeResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyLaunchTemplateDefaultVersionRequestParams struct {
	// The launch template ID
	LaunchTemplateId *string `json:"LaunchTemplateId,omitnil,omitempty" name:"LaunchTemplateId"`

	// The number of the version that you want to set as the default version
	DefaultVersion *int64 `json:"DefaultVersion,omitnil,omitempty" name:"DefaultVersion"`
}

type ModifyLaunchTemplateDefaultVersionRequest struct {
	*tchttp.BaseRequest
	
	// The launch template ID
	LaunchTemplateId *string `json:"LaunchTemplateId,omitnil,omitempty" name:"LaunchTemplateId"`

	// The number of the version that you want to set as the default version
	DefaultVersion *int64 `json:"DefaultVersion,omitnil,omitempty" name:"DefaultVersion"`
}

func (r *ModifyLaunchTemplateDefaultVersionRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyLaunchTemplateDefaultVersionRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "LaunchTemplateId")
	delete(f, "DefaultVersion")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ModifyLaunchTemplateDefaultVersionRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyLaunchTemplateDefaultVersionResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ModifyLaunchTemplateDefaultVersionResponse struct {
	*tchttp.BaseResponse
	Response *ModifyLaunchTemplateDefaultVersionResponseParams `json:"Response"`
}

func (r *ModifyLaunchTemplateDefaultVersionResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyLaunchTemplateDefaultVersionResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type OperationCountLimit struct {
	// Instance operation. Valid values: <br><li>`INSTANCE_DEGRADE`: downgrade an instance<br><li>`INTERNET_CHARGE_TYPE_CHANGE`: modify the billing plan of the network connection
	Operation *string `json:"Operation,omitnil,omitempty" name:"Operation"`

	// Instance ID.
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// Number of operations already performed. If it returns `-1`, it means there is no limit on the times of the operation.
	CurrentCount *int64 `json:"CurrentCount,omitnil,omitempty" name:"CurrentCount"`

	// Maximum number of times you can perform an operation. If it returns `-1`, it means there is no limit on the times of the operation. If it returns `0`, it means that configuration modification is not supported.
	LimitCount *int64 `json:"LimitCount,omitnil,omitempty" name:"LimitCount"`
}

type OsVersion struct {
	// Operating system type
	OsName *string `json:"OsName,omitnil,omitempty" name:"OsName"`

	// Supported operating system versions
	OsVersions []*string `json:"OsVersions,omitnil,omitempty" name:"OsVersions"`

	// Supported operating system architecture
	Architecture []*string `json:"Architecture,omitnil,omitempty" name:"Architecture"`
}

type Placement struct {
	// ID of the availability zone where the instance resides. You can call the [DescribeZones](https://intl.cloud.tencent.com/document/product/213/35071) API and obtain the ID in the returned `Zone` field.
	Zone *string `json:"Zone,omitnil,omitempty" name:"Zone"`

	// Instance'S project ID. obtain this parameter by calling the `ProjectId` field in the return value of [DescribeProjects](https://www.tencentcloud.com/document/product/651/54679). default value 0 means default project.
	ProjectId *int64 `json:"ProjectId,omitnil,omitempty" name:"ProjectId"`

	// Specifies the dedicated host ID list for instance ownership, only used for input parameters. if you purchase a dedicated host and specify this parameter, instances you purchase will be randomly deployed on these dedicated hosts. obtain this parameter by calling the `HostId` field in the return value of [DescribeHosts](https://www.tencentcloud.com/document/product/213/33279?lang=en).
	HostIds []*string `json:"HostIds,omitnil,omitempty" name:"HostIds"`

	// The ID of the CDH to which the instance belongs, only used as an output parameter.
	HostId *string `json:"HostId,omitnil,omitempty" name:"HostId"`
}

type Price struct {
	// Instance price.
	InstancePrice *ItemPrice `json:"InstancePrice,omitnil,omitempty" name:"InstancePrice"`

	// Network price.
	BandwidthPrice *ItemPrice `json:"BandwidthPrice,omitnil,omitempty" name:"BandwidthPrice"`
}

// Predefined struct for user
type PurchaseReservedInstancesOfferingRequestParams struct {
	// The number of the Reserved Instance you are purchasing.
	InstanceCount *int64 `json:"InstanceCount,omitnil,omitempty" name:"InstanceCount"`

	// The ID of the Reserved Instance.
	ReservedInstancesOfferingId *string `json:"ReservedInstancesOfferingId,omitnil,omitempty" name:"ReservedInstancesOfferingId"`

	// Dry run
	DryRun *bool `json:"DryRun,omitnil,omitempty" name:"DryRun"`

	// A unique string supplied by the client to ensure that the request is idempotent. Its maximum length is 64 ASCII characters. If this parameter is not specified, the idempotency of the request cannot be guaranteed.<br>For more information, see Ensuring Idempotency.
	ClientToken *string `json:"ClientToken,omitnil,omitempty" name:"ClientToken"`

	// Reserved instance name.<br><li>The RI name defaults to "unnamed" if this parameter is left empty.</li><li>You can enter any name within 60 characters (including the pattern string).</li>
	ReservedInstanceName *string `json:"ReservedInstanceName,omitnil,omitempty" name:"ReservedInstanceName"`
}

type PurchaseReservedInstancesOfferingRequest struct {
	*tchttp.BaseRequest
	
	// The number of the Reserved Instance you are purchasing.
	InstanceCount *int64 `json:"InstanceCount,omitnil,omitempty" name:"InstanceCount"`

	// The ID of the Reserved Instance.
	ReservedInstancesOfferingId *string `json:"ReservedInstancesOfferingId,omitnil,omitempty" name:"ReservedInstancesOfferingId"`

	// Dry run
	DryRun *bool `json:"DryRun,omitnil,omitempty" name:"DryRun"`

	// A unique string supplied by the client to ensure that the request is idempotent. Its maximum length is 64 ASCII characters. If this parameter is not specified, the idempotency of the request cannot be guaranteed.<br>For more information, see Ensuring Idempotency.
	ClientToken *string `json:"ClientToken,omitnil,omitempty" name:"ClientToken"`

	// Reserved instance name.<br><li>The RI name defaults to "unnamed" if this parameter is left empty.</li><li>You can enter any name within 60 characters (including the pattern string).</li>
	ReservedInstanceName *string `json:"ReservedInstanceName,omitnil,omitempty" name:"ReservedInstanceName"`
}

func (r *PurchaseReservedInstancesOfferingRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *PurchaseReservedInstancesOfferingRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceCount")
	delete(f, "ReservedInstancesOfferingId")
	delete(f, "DryRun")
	delete(f, "ClientToken")
	delete(f, "ReservedInstanceName")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "PurchaseReservedInstancesOfferingRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type PurchaseReservedInstancesOfferingResponseParams struct {
	// The ID of the Reserved Instance purchased.
	ReservedInstanceId *string `json:"ReservedInstanceId,omitnil,omitempty" name:"ReservedInstanceId"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type PurchaseReservedInstancesOfferingResponse struct {
	*tchttp.BaseResponse
	Response *PurchaseReservedInstancesOfferingResponseParams `json:"Response"`
}

func (r *PurchaseReservedInstancesOfferingResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *PurchaseReservedInstancesOfferingResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type RebootInstancesRequestParams struct {
	// Instance IDs. To obtain the instance IDs, you can call [`DescribeInstances`](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1) and look for `InstanceId` in the response. You can operate up to 100 instances in each request.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// Whether to forcibly restart an instance after a normal restart fails. Valid values: <br><li>`TRUE`: yes;<br><li>`FALSE`: no<br><br>Default value: `FALSE`. This parameter has been disused. We recommend using `StopType` instead. Note that `ForceReboot` and `StopType` parameters cannot be specified at the same time.
	ForceReboot *bool `json:"ForceReboot,omitnil,omitempty" name:"ForceReboot"`

	// Shutdown type. Valid values: <br><li>SOFT: soft shutdown<br><li>HARD: hard shutdown<br><li>SOFT_FIRST: perform a soft shutdown first, and perform a hard shutdown if the soft shutdown fails<br><br>Default value: SOFT.
	StopType *string `json:"StopType,omitnil,omitempty" name:"StopType"`
}

type RebootInstancesRequest struct {
	*tchttp.BaseRequest
	
	// Instance IDs. To obtain the instance IDs, you can call [`DescribeInstances`](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1) and look for `InstanceId` in the response. You can operate up to 100 instances in each request.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// Whether to forcibly restart an instance after a normal restart fails. Valid values: <br><li>`TRUE`: yes;<br><li>`FALSE`: no<br><br>Default value: `FALSE`. This parameter has been disused. We recommend using `StopType` instead. Note that `ForceReboot` and `StopType` parameters cannot be specified at the same time.
	ForceReboot *bool `json:"ForceReboot,omitnil,omitempty" name:"ForceReboot"`

	// Shutdown type. Valid values: <br><li>SOFT: soft shutdown<br><li>HARD: hard shutdown<br><li>SOFT_FIRST: perform a soft shutdown first, and perform a hard shutdown if the soft shutdown fails<br><br>Default value: SOFT.
	StopType *string `json:"StopType,omitnil,omitempty" name:"StopType"`
}

func (r *RebootInstancesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *RebootInstancesRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceIds")
	delete(f, "ForceReboot")
	delete(f, "StopType")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "RebootInstancesRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type RebootInstancesResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type RebootInstancesResponse struct {
	*tchttp.BaseResponse
	Response *RebootInstancesResponseParams `json:"Response"`
}

func (r *RebootInstancesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *RebootInstancesResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type RegionInfo struct {
	// Region name, such as `ap-guangzhou`
	Region *string `json:"Region,omitnil,omitempty" name:"Region"`

	// Region description, such as South China (Guangzhou)
	RegionName *string `json:"RegionName,omitnil,omitempty" name:"RegionName"`

	// Whether the region is available
	RegionState *string `json:"RegionState,omitnil,omitempty" name:"RegionState"`
}

// Predefined struct for user
type RemoveChcAssistVpcRequestParams struct {
	// CHC ID
	ChcIds []*string `json:"ChcIds,omitnil,omitempty" name:"ChcIds"`
}

type RemoveChcAssistVpcRequest struct {
	*tchttp.BaseRequest
	
	// CHC ID
	ChcIds []*string `json:"ChcIds,omitnil,omitempty" name:"ChcIds"`
}

func (r *RemoveChcAssistVpcRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *RemoveChcAssistVpcRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ChcIds")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "RemoveChcAssistVpcRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type RemoveChcAssistVpcResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type RemoveChcAssistVpcResponse struct {
	*tchttp.BaseResponse
	Response *RemoveChcAssistVpcResponseParams `json:"Response"`
}

func (r *RemoveChcAssistVpcResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *RemoveChcAssistVpcResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type RemoveChcDeployVpcRequestParams struct {
	// CHC ID
	ChcIds []*string `json:"ChcIds,omitnil,omitempty" name:"ChcIds"`
}

type RemoveChcDeployVpcRequest struct {
	*tchttp.BaseRequest
	
	// CHC ID
	ChcIds []*string `json:"ChcIds,omitnil,omitempty" name:"ChcIds"`
}

func (r *RemoveChcDeployVpcRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *RemoveChcDeployVpcRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ChcIds")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "RemoveChcDeployVpcRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type RemoveChcDeployVpcResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type RemoveChcDeployVpcResponse struct {
	*tchttp.BaseResponse
	Response *RemoveChcDeployVpcResponseParams `json:"Response"`
}

func (r *RemoveChcDeployVpcResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *RemoveChcDeployVpcResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type RenewInstancesRequestParams struct {
	// One or more instance IDs to be operated. You can obtain the instance ID through the `InstanceId` in the return value from the API [DescribeInstances](https://www.tencentcloud.com/document/api/213/33258). The maximum number of instances per request is 100.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// Prepaid mode, that is, parameter settings related to monthly/annual subscription. specifies attributes of a monthly subscription instance, such as renewal duration and whether to enable auto-renewal, by specifying this parameter. <dx-alert infotype="explain" title="">.
	// Annual and monthly subscription instances. this parameter is a required parameter.</dx-alert>.
	InstanceChargePrepaid *InstanceChargePrepaid `json:"InstanceChargePrepaid,omitnil,omitempty" name:"InstanceChargePrepaid"`

	// Whether to renew the elastic data disk. valid values:<br><li>true: indicates renewing the annual and monthly subscription instance and its mounted elastic data disk simultaneously</li><li>false: indicates renewing the annual and monthly subscription instance while no longer renewing its mounted elastic data disk</li><br>default value: true.
	RenewPortableDataDisk *bool `json:"RenewPortableDataDisk,omitnil,omitempty" name:"RenewPortableDataDisk"`
}

type RenewInstancesRequest struct {
	*tchttp.BaseRequest
	
	// One or more instance IDs to be operated. You can obtain the instance ID through the `InstanceId` in the return value from the API [DescribeInstances](https://www.tencentcloud.com/document/api/213/33258). The maximum number of instances per request is 100.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// Prepaid mode, that is, parameter settings related to monthly/annual subscription. specifies attributes of a monthly subscription instance, such as renewal duration and whether to enable auto-renewal, by specifying this parameter. <dx-alert infotype="explain" title="">.
	// Annual and monthly subscription instances. this parameter is a required parameter.</dx-alert>.
	InstanceChargePrepaid *InstanceChargePrepaid `json:"InstanceChargePrepaid,omitnil,omitempty" name:"InstanceChargePrepaid"`

	// Whether to renew the elastic data disk. valid values:<br><li>true: indicates renewing the annual and monthly subscription instance and its mounted elastic data disk simultaneously</li><li>false: indicates renewing the annual and monthly subscription instance while no longer renewing its mounted elastic data disk</li><br>default value: true.
	RenewPortableDataDisk *bool `json:"RenewPortableDataDisk,omitnil,omitempty" name:"RenewPortableDataDisk"`
}

func (r *RenewInstancesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *RenewInstancesRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceIds")
	delete(f, "InstanceChargePrepaid")
	delete(f, "RenewPortableDataDisk")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "RenewInstancesRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type RenewInstancesResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type RenewInstancesResponse struct {
	*tchttp.BaseResponse
	Response *RenewInstancesResponseParams `json:"Response"`
}

func (r *RenewInstancesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *RenewInstancesResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type ReservedInstanceConfigInfoItem struct {
	// Abbreviation name of the instance type.
	Type *string `json:"Type,omitnil,omitempty" name:"Type"`

	// Full name of the instance type.
	TypeName *string `json:"TypeName,omitnil,omitempty" name:"TypeName"`

	// Priority.
	Order *int64 `json:"Order,omitnil,omitempty" name:"Order"`

	// List of instance families.
	InstanceFamilies []*ReservedInstanceFamilyItem `json:"InstanceFamilies,omitnil,omitempty" name:"InstanceFamilies"`
}

type ReservedInstanceFamilyItem struct {
	// Instance family.
	InstanceFamily *string `json:"InstanceFamily,omitnil,omitempty" name:"InstanceFamily"`

	// Priority.
	Order *int64 `json:"Order,omitnil,omitempty" name:"Order"`

	// List of instance types.
	InstanceTypes []*ReservedInstanceTypeItem `json:"InstanceTypes,omitnil,omitempty" name:"InstanceTypes"`
}

type ReservedInstancePrice struct {
	// Original upfront payment, in USD.
	OriginalFixedPrice *float64 `json:"OriginalFixedPrice,omitnil,omitempty" name:"OriginalFixedPrice"`

	// Discounted upfront payment, in USD.
	DiscountFixedPrice *float64 `json:"DiscountFixedPrice,omitnil,omitempty" name:"DiscountFixedPrice"`

	// Original subsequent unit price, in USD/hr.
	OriginalUsagePrice *float64 `json:"OriginalUsagePrice,omitnil,omitempty" name:"OriginalUsagePrice"`

	// Discounted subsequent unit price, in USD/hr.
	DiscountUsagePrice *float64 `json:"DiscountUsagePrice,omitnil,omitempty" name:"DiscountUsagePrice"`

	// Discount on upfront cost. For example, 20.0 indicates 80% off. Note: This field may return null, indicating that no valid value is found.
	// Note: This field may return null, indicating that no valid value is found.
	FixedPriceDiscount *float64 `json:"FixedPriceDiscount,omitnil,omitempty" name:"FixedPriceDiscount"`

	// Discount on subsequent cost. For example, 20.0 indicates 80% off. Note: This field may return null, indicating that no valid value is found.
	// Note: This field may return null, indicating that no valid value is found.
	UsagePriceDiscount *float64 `json:"UsagePriceDiscount,omitnil,omitempty" name:"UsagePriceDiscount"`
}

type ReservedInstancePriceItem struct {
	// Payment method. Valid values: All Upfront, Partial Upfront, and No Upfront.
	OfferingType *string `json:"OfferingType,omitnil,omitempty" name:"OfferingType"`

	// Upfront payment, in USD.
	FixedPrice *float64 `json:"FixedPrice,omitnil,omitempty" name:"FixedPrice"`

	// Subsequent unit price, in USD/hr.
	UsagePrice *float64 `json:"UsagePrice,omitnil,omitempty" name:"UsagePrice"`

	// The ID of the reserved instance offering.
	ReservedInstancesOfferingId *string `json:"ReservedInstancesOfferingId,omitnil,omitempty" name:"ReservedInstancesOfferingId"`

	// The availability zone in which the reserved instance can be purchased.
	Zone *string `json:"Zone,omitnil,omitempty" name:"Zone"`

	// The **validity** of the reserved instance in seconds, which is the purchased usage period. For example, `31536000`.
	// Unit: second
	Duration *uint64 `json:"Duration,omitnil,omitempty" name:"Duration"`

	// The operating system of the reserved instance, such as `Linux`.
	// Valid value: `Linux`.
	ProductDescription *string `json:"ProductDescription,omitnil,omitempty" name:"ProductDescription"`

	// Discount price for subsequent total cost, in USD/hr.
	DiscountUsagePrice *float64 `json:"DiscountUsagePrice,omitnil,omitempty" name:"DiscountUsagePrice"`

	// Discount price for upfront total cost, in USD.
	DiscountFixedPrice *float64 `json:"DiscountFixedPrice,omitnil,omitempty" name:"DiscountFixedPrice"`
}

type ReservedInstanceTypeItem struct {
	// Instance type.
	InstanceType *string `json:"InstanceType,omitnil,omitempty" name:"InstanceType"`

	// Number of CPU cores.
	Cpu *uint64 `json:"Cpu,omitnil,omitempty" name:"Cpu"`

	// Memory size.
	Memory *uint64 `json:"Memory,omitnil,omitempty" name:"Memory"`

	// Number of GPUs.
	Gpu *uint64 `json:"Gpu,omitnil,omitempty" name:"Gpu"`

	// Number of FPGAs.
	Fpga *uint64 `json:"Fpga,omitnil,omitempty" name:"Fpga"`

	// Number of local storage blocks.
	StorageBlock *uint64 `json:"StorageBlock,omitnil,omitempty" name:"StorageBlock"`

	// Number of NICs.
	NetworkCard *uint64 `json:"NetworkCard,omitnil,omitempty" name:"NetworkCard"`

	// Maximum bandwidth.
	MaxBandwidth *float64 `json:"MaxBandwidth,omitnil,omitempty" name:"MaxBandwidth"`

	// CPU frequency.
	Frequency *string `json:"Frequency,omitnil,omitempty" name:"Frequency"`

	// CPU type.
	CpuModelName *string `json:"CpuModelName,omitnil,omitempty" name:"CpuModelName"`

	// Packet forwarding rate.
	Pps *uint64 `json:"Pps,omitnil,omitempty" name:"Pps"`

	// Other information.
	Externals *Externals `json:"Externals,omitnil,omitempty" name:"Externals"`

	// Remarks.
	Remark *string `json:"Remark,omitnil,omitempty" name:"Remark"`

	// Price information about the reserved instance.
	Prices []*ReservedInstancePriceItem `json:"Prices,omitnil,omitempty" name:"Prices"`
}

type ReservedInstances struct {
	// (This field has been deprecated. ReservedInstanceId is recommended.) IDs of purchased reserved instances. For example, ri-rtbh4han.
	//
	// Deprecated: ReservedInstancesId is deprecated.
	ReservedInstancesId *string `json:"ReservedInstancesId,omitnil,omitempty" name:"ReservedInstancesId"`

	// Specifications of reserved instances. For example, S3.MEDIUM4.
	// Return values: <a href="https://intl.cloud.tencent.com/document/product/213/11518?from_cn_redirect=1">Reserved instance specification list.</a>
	InstanceType *string `json:"InstanceType,omitnil,omitempty" name:"InstanceType"`

	// Availability zones in which reserved instances can be purchased. For example, ap-guangzhou-1.
	// Return values: <a href="https://intl.cloud.tencent.com/document/product/213/6091?from_cn_redirect=1">Availability zone list.</a>
	Zone *string `json:"Zone,omitnil,omitempty" name:"Zone"`

	// Billing start time of reserved instances. For example, 1949-10-01 00:00:00.
	StartTime *string `json:"StartTime,omitnil,omitempty" name:"StartTime"`

	// Billing end time of reserved instances. For example, 1949-10-01 00:00:00.
	EndTime *string `json:"EndTime,omitnil,omitempty" name:"EndTime"`

	// Validity periods of reserved instances, which is the purchase duration of reserved instances. For example, 31536000.
	// Unit: second.
	Duration *int64 `json:"Duration,omitnil,omitempty" name:"Duration"`

	// Number of purchased reserved instances. For example, 10.
	InstanceCount *int64 `json:"InstanceCount,omitnil,omitempty" name:"InstanceCount"`

	// Platform descriptions (operating systems) of reserved instances. For example, linux.
	// Return value: linux.
	ProductDescription *string `json:"ProductDescription,omitnil,omitempty" name:"ProductDescription"`

	// Statuses of purchased reserved instances. For example: active.
	// Return values: active (created) | pending (waiting to be created) | retired (expired).
	State *string `json:"State,omitnil,omitempty" name:"State"`

	// Billing currencies of purchasable reserved instances. Use standard currency codes defined in ISO 4217. For example, USD.
	// Return value: USD.
	CurrencyCode *string `json:"CurrencyCode,omitnil,omitempty" name:"CurrencyCode"`

	// Payment types of reserved instances. For example, All Upfront.
	// Return value: All Upfront (fully prepaid).
	OfferingType *string `json:"OfferingType,omitnil,omitempty" name:"OfferingType"`

	// Types of reserved instances. For example, S3.
	// Return values: <a href="https://intl.cloud.tencent.com/document/product/213/11518?from_cn_redirect=1">Reserved instance type list.</a>
	InstanceFamily *string `json:"InstanceFamily,omitnil,omitempty" name:"InstanceFamily"`

	// IDs of purchased reserved instances. For example, ri-rtbh4han.
	ReservedInstanceId *string `json:"ReservedInstanceId,omitnil,omitempty" name:"ReservedInstanceId"`

	// Display names of reserved instances. For example, riname-01.
	ReservedInstanceName *string `json:"ReservedInstanceName,omitnil,omitempty" name:"ReservedInstanceName"`
}

type ReservedInstancesOffering struct {
	// The availability zones in which the Reserved Instance can be purchased, such as ap-guangzhou-1.
	// Valid value: <a href="https://intl.cloud.tencent.com/document/product/213/6091?from_cn_redirect=1">Availability Zones</a>
	Zone *string `json:"Zone,omitnil,omitempty" name:"Zone"`

	// The billing currency of the Reserved Instance you are purchasing. It's specified using ISO 4217 standard currency.
	// Value: USD.
	CurrencyCode *string `json:"CurrencyCode,omitnil,omitempty" name:"CurrencyCode"`

	// The **validity** of the Reserved Instance in seconds, which is the purchased usage period. For example, 31536000.
	// Unit: second
	Duration *int64 `json:"Duration,omitnil,omitempty" name:"Duration"`

	// The purchase price of the Reserved Instance, such as 4000.0.
	// Unit: this field uses the currency code specified in `currencyCode`, and only supports “USD” at this time.
	FixedPrice *float64 `json:"FixedPrice,omitnil,omitempty" name:"FixedPrice"`

	// The instance model of the Reserved Instance, such as S3.MEDIUM4.
	// Valid values: please see <a href="https://intl.cloud.tencent.com/document/product/213/11518">Reserved Instance Types</a>
	InstanceType *string `json:"InstanceType,omitnil,omitempty" name:"InstanceType"`

	// The payment term of the Reserved Instance, such as **All Upfront**.
	// Valid value: All Upfront.
	OfferingType *string `json:"OfferingType,omitnil,omitempty" name:"OfferingType"`

	// The ID of the Reserved Instance offering, such as 650c138f-ae7e-4750-952a-96841d6e9fc1.
	ReservedInstancesOfferingId *string `json:"ReservedInstancesOfferingId,omitnil,omitempty" name:"ReservedInstancesOfferingId"`

	// The operating system of the Reserved Instance, such as **linux**.
	// Valid value: linux.
	ProductDescription *string `json:"ProductDescription,omitnil,omitempty" name:"ProductDescription"`

	// The hourly usage price of the Reserved Instance, such as 0.0.
	// Currently, the only supported payment mode is **All Upfront**, so the default value of `UsagePrice` is 0 USD/hr.
	// Unit: USD/hr. This field uses the currency code specified in `currencyCode`, and only supports “USD” at this time.
	UsagePrice *float64 `json:"UsagePrice,omitnil,omitempty" name:"UsagePrice"`
}

// Predefined struct for user
type ResetInstanceRequestParams struct {
	// Instance ID. To obtain the instance IDs, you can call [`DescribeInstances`](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1) and look for `InstanceId` in the response.
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// Specified effective [image](https://intl.cloud.tencent.com/document/product/213/4940?from_cn_redirect=1) ID in the format of `img-xxx`. There are four types of images:<br/><li>Public images</li><li>Custom images</li><li>Shared images</li><li>Marketplace images </li><br/>You can obtain the available image IDs in the following ways:<br/><li>for IDs of `public images`, `custom images`, and `shared images`, log in to the [CVM console](https://console.cloud.tencent.com/cvm/image?rid=1&imageType=PUBLIC_IMAGE); for IDs of `marketplace images`, go to [Cloud Marketplace](https://market.cloud.tencent.com/list).</li><li>Call the API [DescribeImages](https://intl.cloud.tencent.com/document/api/213/15715?from_cn_redirect=1) and look for `ImageId` in the response.</li>
	// <br>Default value: current image.
	ImageId *string `json:"ImageId,omitnil,omitempty" name:"ImageId"`

	// Configurations of the system disk. For an instance whose system disk is a cloud disk, this parameter can be used to expand the system disk by specifying a new capacity after reinstallation. The system disk capacity can only be expanded but not shrunk. Reinstalling the system can only resize rather than changing the type of the system disk.
	SystemDisk *SystemDisk `json:"SystemDisk,omitnil,omitempty" name:"SystemDisk"`

	// Login settings of the instance. You can use this parameter to set the login method, password, and key of the instance or keep the login settings of the original image. By default, a random password will be generated and sent to you via the Message Center.
	LoginSettings *LoginSettings `json:"LoginSettings,omitnil,omitempty" name:"LoginSettings"`

	// Enhanced services. You can specify whether to enable services such as Cloud Security and Cloud Monitor through this parameter. If this parameter is not specified, Cloud Monitor and Cloud Security are enabled for public images by default, but not enabled for custom images and marketplace images by default. Instead, they use services retained in the images.
	EnhancedService *EnhancedService `json:"EnhancedService,omitnil,omitempty" name:"EnhancedService"`

	// When reinstalling a system, you can modify an instance's hostname.<br><li>Period (.) and hyphen (-) should not be used as the first or last character of the hostname, and should not be used consecutively.</li><br><li>Windows instances: The hostname should contain 2 to 15 characters, including letters (case insensitive), digits, and hyphens (-), does not support periods (.), and should not be all digits.</li><br><li>Instances of other types (such as Linux instances): The hostname should contain 2 to 60 characters, including multiple periods (.), with each segment between periods considered as one section. Each section can contain letters (case insensitive), digits, and hyphens (-).</li>
	HostName *string `json:"HostName,omitnil,omitempty" name:"HostName"`

	// User data provided to the instance. This parameter needs to be encoded in base64 format with the maximum size of 16 KB. For more information on how to get the value of this parameter, see the commands you need to execute on startup for [Windows](https://intl.cloud.tencent.com/document/product/213/17526) or [Linux](https://intl.cloud.tencent.com/document/product/213/17525).
	UserData *string `json:"UserData,omitnil,omitempty" name:"UserData"`
}

type ResetInstanceRequest struct {
	*tchttp.BaseRequest
	
	// Instance ID. To obtain the instance IDs, you can call [`DescribeInstances`](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1) and look for `InstanceId` in the response.
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// Specified effective [image](https://intl.cloud.tencent.com/document/product/213/4940?from_cn_redirect=1) ID in the format of `img-xxx`. There are four types of images:<br/><li>Public images</li><li>Custom images</li><li>Shared images</li><li>Marketplace images </li><br/>You can obtain the available image IDs in the following ways:<br/><li>for IDs of `public images`, `custom images`, and `shared images`, log in to the [CVM console](https://console.cloud.tencent.com/cvm/image?rid=1&imageType=PUBLIC_IMAGE); for IDs of `marketplace images`, go to [Cloud Marketplace](https://market.cloud.tencent.com/list).</li><li>Call the API [DescribeImages](https://intl.cloud.tencent.com/document/api/213/15715?from_cn_redirect=1) and look for `ImageId` in the response.</li>
	// <br>Default value: current image.
	ImageId *string `json:"ImageId,omitnil,omitempty" name:"ImageId"`

	// Configurations of the system disk. For an instance whose system disk is a cloud disk, this parameter can be used to expand the system disk by specifying a new capacity after reinstallation. The system disk capacity can only be expanded but not shrunk. Reinstalling the system can only resize rather than changing the type of the system disk.
	SystemDisk *SystemDisk `json:"SystemDisk,omitnil,omitempty" name:"SystemDisk"`

	// Login settings of the instance. You can use this parameter to set the login method, password, and key of the instance or keep the login settings of the original image. By default, a random password will be generated and sent to you via the Message Center.
	LoginSettings *LoginSettings `json:"LoginSettings,omitnil,omitempty" name:"LoginSettings"`

	// Enhanced services. You can specify whether to enable services such as Cloud Security and Cloud Monitor through this parameter. If this parameter is not specified, Cloud Monitor and Cloud Security are enabled for public images by default, but not enabled for custom images and marketplace images by default. Instead, they use services retained in the images.
	EnhancedService *EnhancedService `json:"EnhancedService,omitnil,omitempty" name:"EnhancedService"`

	// When reinstalling a system, you can modify an instance's hostname.<br><li>Period (.) and hyphen (-) should not be used as the first or last character of the hostname, and should not be used consecutively.</li><br><li>Windows instances: The hostname should contain 2 to 15 characters, including letters (case insensitive), digits, and hyphens (-), does not support periods (.), and should not be all digits.</li><br><li>Instances of other types (such as Linux instances): The hostname should contain 2 to 60 characters, including multiple periods (.), with each segment between periods considered as one section. Each section can contain letters (case insensitive), digits, and hyphens (-).</li>
	HostName *string `json:"HostName,omitnil,omitempty" name:"HostName"`

	// User data provided to the instance. This parameter needs to be encoded in base64 format with the maximum size of 16 KB. For more information on how to get the value of this parameter, see the commands you need to execute on startup for [Windows](https://intl.cloud.tencent.com/document/product/213/17526) or [Linux](https://intl.cloud.tencent.com/document/product/213/17525).
	UserData *string `json:"UserData,omitnil,omitempty" name:"UserData"`
}

func (r *ResetInstanceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ResetInstanceRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceId")
	delete(f, "ImageId")
	delete(f, "SystemDisk")
	delete(f, "LoginSettings")
	delete(f, "EnhancedService")
	delete(f, "HostName")
	delete(f, "UserData")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ResetInstanceRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ResetInstanceResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ResetInstanceResponse struct {
	*tchttp.BaseResponse
	Response *ResetInstanceResponseParams `json:"Response"`
}

func (r *ResetInstanceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ResetInstanceResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ResetInstancesInternetMaxBandwidthRequestParams struct {
	// Instance ID(s). To obtain the instance IDs, you can call [`DescribeInstances`](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1) and look for `InstanceId` in the response. The maximum number of instances in each request is 100. When changing the bandwidth of instances with `BANDWIDTH_PREPAID` or `BANDWIDTH_POSTPAID_BY_HOUR` as the network billing method, you can only specify one instance at a time.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// Configuration of public network egress bandwidth. The maximum bandwidth varies among different models. For more information, see the documentation on bandwidth limits. Currently only the `InternetMaxBandwidthOut` parameter is supported.
	InternetAccessible *InternetAccessible `json:"InternetAccessible,omitnil,omitempty" name:"InternetAccessible"`

	// Date from which the new bandwidth takes effect. Format: `YYYY-MM-DD`, such as `2016-10-30`. The starting date cannot be earlier than the current date. If the starting date is the current date, the new bandwidth takes effect immediately. This parameter is only valid for prepaid bandwidth. If you specify the parameter for bandwidth with other network billing methods, an error code will be returned.
	StartTime *string `json:"StartTime,omitnil,omitempty" name:"StartTime"`

	// Date until which the bandwidth takes effect, in the format of `YYYY-MM-DD`, such as `2016-10-30`. The validity period of the new bandwidth covers the end date. The end date should not be later than the expiration date of a monthly subscription instance. You can obtain the expiration date of an instance through the `ExpiredTime` in the return value from the API [DescribeInstances](https://intl.cloud.tencent.com/document/api/213/9388?from_cn_redirect=1). This parameter is only valid for monthly subscription bandwidth, and is not supported for bandwidth billed by other modes. Otherwise, the API will return a corresponding error code.
	EndTime *string `json:"EndTime,omitnil,omitempty" name:"EndTime"`
}

type ResetInstancesInternetMaxBandwidthRequest struct {
	*tchttp.BaseRequest
	
	// Instance ID(s). To obtain the instance IDs, you can call [`DescribeInstances`](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1) and look for `InstanceId` in the response. The maximum number of instances in each request is 100. When changing the bandwidth of instances with `BANDWIDTH_PREPAID` or `BANDWIDTH_POSTPAID_BY_HOUR` as the network billing method, you can only specify one instance at a time.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// Configuration of public network egress bandwidth. The maximum bandwidth varies among different models. For more information, see the documentation on bandwidth limits. Currently only the `InternetMaxBandwidthOut` parameter is supported.
	InternetAccessible *InternetAccessible `json:"InternetAccessible,omitnil,omitempty" name:"InternetAccessible"`

	// Date from which the new bandwidth takes effect. Format: `YYYY-MM-DD`, such as `2016-10-30`. The starting date cannot be earlier than the current date. If the starting date is the current date, the new bandwidth takes effect immediately. This parameter is only valid for prepaid bandwidth. If you specify the parameter for bandwidth with other network billing methods, an error code will be returned.
	StartTime *string `json:"StartTime,omitnil,omitempty" name:"StartTime"`

	// Date until which the bandwidth takes effect, in the format of `YYYY-MM-DD`, such as `2016-10-30`. The validity period of the new bandwidth covers the end date. The end date should not be later than the expiration date of a monthly subscription instance. You can obtain the expiration date of an instance through the `ExpiredTime` in the return value from the API [DescribeInstances](https://intl.cloud.tencent.com/document/api/213/9388?from_cn_redirect=1). This parameter is only valid for monthly subscription bandwidth, and is not supported for bandwidth billed by other modes. Otherwise, the API will return a corresponding error code.
	EndTime *string `json:"EndTime,omitnil,omitempty" name:"EndTime"`
}

func (r *ResetInstancesInternetMaxBandwidthRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ResetInstancesInternetMaxBandwidthRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceIds")
	delete(f, "InternetAccessible")
	delete(f, "StartTime")
	delete(f, "EndTime")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ResetInstancesInternetMaxBandwidthRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ResetInstancesInternetMaxBandwidthResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ResetInstancesInternetMaxBandwidthResponse struct {
	*tchttp.BaseResponse
	Response *ResetInstancesInternetMaxBandwidthResponseParams `json:"Response"`
}

func (r *ResetInstancesInternetMaxBandwidthResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ResetInstancesInternetMaxBandwidthResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ResetInstancesPasswordRequestParams struct {
	// Instance ID(s). To obtain the instance IDs, you can call [`DescribeInstances`](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1) and look for `InstanceId` in the response. The maximum number of instances in each request is 100.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// Login password of the instance. The password requirements vary among different operating systems:
	// For a Linux instance, the password must be 8 to 30 characters in length; password with more than 12 characters is recommended. It cannot begin with "/", and must contain at least one character from three of the following categories: <br><li>Lowercase letters: [a-z]<br><li>Uppercase letters: [A-Z]<br><li>Numbers: 0-9<br><li>Special characters: ()\`\~!@#$%^&\*-+=\_|{}[]:;'<>,.?/
	// For a Windows CVM, the password must be 12 to 30 characters in length. It cannot begin with "/" or contain a username. It must contain at least one character from three of the following categories: <br><li>Lowercase letters: [a-z]<br><li>Uppercase letters: [A-Z]<br><li>Numbers: 0-9<br><li>Special characters: ()\`\~!@#$%^&\*-+=\_|{}[]:;' <>,.?/<br><li>If the specified instances include both `Linux` and `Windows` instances, you will need to follow the password requirements for `Windows` instances.
	Password *string `json:"Password,omitnil,omitempty" name:"Password"`

	// Username of the instance operating system for which the password needs to be reset. This parameter is limited to 64 characters.
	UserName *string `json:"UserName,omitnil,omitempty" name:"UserName"`

	// Whether to force shut down a running instances. It is recommended to manually shut down a running instance before resetting the user password. Valid values: <br><li>TRUE: force shut down an instance after a normal shutdown fails. <br><li>FALSE: do not force shut down an instance after a normal shutdown fails. <br><br>Default value: FALSE. <br><br>A forced shutdown is similar to switching off the power of a physical computer. It may cause data loss or file system corruption. Be sure to only force shut down a CVM when it cannot be shut down normally.
	ForceStop *bool `json:"ForceStop,omitnil,omitempty" name:"ForceStop"`
}

type ResetInstancesPasswordRequest struct {
	*tchttp.BaseRequest
	
	// Instance ID(s). To obtain the instance IDs, you can call [`DescribeInstances`](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1) and look for `InstanceId` in the response. The maximum number of instances in each request is 100.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// Login password of the instance. The password requirements vary among different operating systems:
	// For a Linux instance, the password must be 8 to 30 characters in length; password with more than 12 characters is recommended. It cannot begin with "/", and must contain at least one character from three of the following categories: <br><li>Lowercase letters: [a-z]<br><li>Uppercase letters: [A-Z]<br><li>Numbers: 0-9<br><li>Special characters: ()\`\~!@#$%^&\*-+=\_|{}[]:;'<>,.?/
	// For a Windows CVM, the password must be 12 to 30 characters in length. It cannot begin with "/" or contain a username. It must contain at least one character from three of the following categories: <br><li>Lowercase letters: [a-z]<br><li>Uppercase letters: [A-Z]<br><li>Numbers: 0-9<br><li>Special characters: ()\`\~!@#$%^&\*-+=\_|{}[]:;' <>,.?/<br><li>If the specified instances include both `Linux` and `Windows` instances, you will need to follow the password requirements for `Windows` instances.
	Password *string `json:"Password,omitnil,omitempty" name:"Password"`

	// Username of the instance operating system for which the password needs to be reset. This parameter is limited to 64 characters.
	UserName *string `json:"UserName,omitnil,omitempty" name:"UserName"`

	// Whether to force shut down a running instances. It is recommended to manually shut down a running instance before resetting the user password. Valid values: <br><li>TRUE: force shut down an instance after a normal shutdown fails. <br><li>FALSE: do not force shut down an instance after a normal shutdown fails. <br><br>Default value: FALSE. <br><br>A forced shutdown is similar to switching off the power of a physical computer. It may cause data loss or file system corruption. Be sure to only force shut down a CVM when it cannot be shut down normally.
	ForceStop *bool `json:"ForceStop,omitnil,omitempty" name:"ForceStop"`
}

func (r *ResetInstancesPasswordRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ResetInstancesPasswordRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceIds")
	delete(f, "Password")
	delete(f, "UserName")
	delete(f, "ForceStop")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ResetInstancesPasswordRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ResetInstancesPasswordResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ResetInstancesPasswordResponse struct {
	*tchttp.BaseResponse
	Response *ResetInstancesPasswordResponseParams `json:"Response"`
}

func (r *ResetInstancesPasswordResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ResetInstancesPasswordResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ResetInstancesTypeRequestParams struct {
	// Instance ID(s). To obtain the instance IDs, you can call the [`DescribeInstances`](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1) API and find the value `InstanceId` in the response. The maximum number of instances in each request is 1.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// Instance model. Different resource specifications are specified for different models. For specific values, call [DescribeInstanceTypeConfigs](https://intl.cloud.tencent.com/document/api/213/15749?from_cn_redirect=1) to get the latest specification list or refer to [Instance Types](https://intl.cloud.tencent.com/document/product/213/11518?from_cn_redirect=1).
	InstanceType *string `json:"InstanceType,omitnil,omitempty" name:"InstanceType"`

	// Forced shutdown of a running instances. We recommend you firstly try to shut down a running instance manually. Valid values: <br><li>TRUE: forced shutdown of an instance after a normal shutdown fails.<br><li>FALSE: no forced shutdown of an instance after a normal shutdown fails.<br><br>Default value: FALSE.<br><br>A forced shutdown is similar to switching off the power of a physical computer. It may cause data loss or file system corruption. Be sure to only force a CVM to shut off if the normal shutdown fails.
	ForceStop *bool `json:"ForceStop,omitnil,omitempty" name:"ForceStop"`
}

type ResetInstancesTypeRequest struct {
	*tchttp.BaseRequest
	
	// Instance ID(s). To obtain the instance IDs, you can call the [`DescribeInstances`](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1) API and find the value `InstanceId` in the response. The maximum number of instances in each request is 1.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// Instance model. Different resource specifications are specified for different models. For specific values, call [DescribeInstanceTypeConfigs](https://intl.cloud.tencent.com/document/api/213/15749?from_cn_redirect=1) to get the latest specification list or refer to [Instance Types](https://intl.cloud.tencent.com/document/product/213/11518?from_cn_redirect=1).
	InstanceType *string `json:"InstanceType,omitnil,omitempty" name:"InstanceType"`

	// Forced shutdown of a running instances. We recommend you firstly try to shut down a running instance manually. Valid values: <br><li>TRUE: forced shutdown of an instance after a normal shutdown fails.<br><li>FALSE: no forced shutdown of an instance after a normal shutdown fails.<br><br>Default value: FALSE.<br><br>A forced shutdown is similar to switching off the power of a physical computer. It may cause data loss or file system corruption. Be sure to only force a CVM to shut off if the normal shutdown fails.
	ForceStop *bool `json:"ForceStop,omitnil,omitempty" name:"ForceStop"`
}

func (r *ResetInstancesTypeRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ResetInstancesTypeRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceIds")
	delete(f, "InstanceType")
	delete(f, "ForceStop")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ResetInstancesTypeRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ResetInstancesTypeResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ResetInstancesTypeResponse struct {
	*tchttp.BaseResponse
	Response *ResetInstancesTypeResponseParams `json:"Response"`
}

func (r *ResetInstancesTypeResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ResetInstancesTypeResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ResizeInstanceDisksRequestParams struct {
	// Instance ID to be operated. can be obtained from the `InstanceId` in the return value from the DescribeInstances api (https://www.tencentcloud.com/document/api/213/15728?from_cn_redirect=1).
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// Specifies the configuration information of the data disk to be expanded, only supporting specifying the target capacity of the disk to be expanded. only non-elastic data disks (with `Portable` being `false` in the return values of [DescribeDisks](https://www.tencentcloud.com/document/api/362/16315?from_cn_redirect=1)) can be expanded. the unit of data disk capacity is GiB. the minimum expansion step is 10 GiB. for data disk type selection, refer to [disk product introduction](https://www.tencentcloud.com/document/product/362/2353?from_cn_redirect=1). the available data disk type is restricted by the instance type `InstanceType`. additionally, the maximum allowable capacity for expansion varies by data disk type.
	// <dx-alert infotype="explain" title="">You should specify either DataDisks or SystemDisk, but you cannot specify both at the same time.</dx-alert>
	DataDisks []*DataDisk `json:"DataDisks,omitnil,omitempty" name:"DataDisks"`

	// Specifies whether to forcibly shut down a running instance. it is recommended to manually shut down a running instance first and then expand the instance disk. valid values:<br><li>true: forcibly shut down an instance after a normal shutdown fails.</li><br><li>false: do not forcibly shut down an instance after a normal shutdown fails.</li><br><br>default value: false.<br><br>forced shutdown is equivalent to turning off a physical computer's power switch. forced shutdown may cause data loss or file system corruption and should only be used when a server cannot be shut down normally.
	ForceStop *bool `json:"ForceStop,omitnil,omitempty" name:"ForceStop"`

	// System disk configuration information to be expanded. only supports specifying the purpose capacity of the disk to be expanded. only supports cloud disk expansion.
	// <dx-alert infotype="explain" title="">You should specify either DataDisks or SystemDisk, but you cannot specify both at the same time.</dx-alert>
	SystemDisk *SystemDisk `json:"SystemDisk,omitnil,omitempty" name:"SystemDisk"`

	// Whether the cloud disk is expanded online.
	ResizeOnline *bool `json:"ResizeOnline,omitnil,omitempty" name:"ResizeOnline"`
}

type ResizeInstanceDisksRequest struct {
	*tchttp.BaseRequest
	
	// Instance ID to be operated. can be obtained from the `InstanceId` in the return value from the DescribeInstances api (https://www.tencentcloud.com/document/api/213/15728?from_cn_redirect=1).
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// Specifies the configuration information of the data disk to be expanded, only supporting specifying the target capacity of the disk to be expanded. only non-elastic data disks (with `Portable` being `false` in the return values of [DescribeDisks](https://www.tencentcloud.com/document/api/362/16315?from_cn_redirect=1)) can be expanded. the unit of data disk capacity is GiB. the minimum expansion step is 10 GiB. for data disk type selection, refer to [disk product introduction](https://www.tencentcloud.com/document/product/362/2353?from_cn_redirect=1). the available data disk type is restricted by the instance type `InstanceType`. additionally, the maximum allowable capacity for expansion varies by data disk type.
	// <dx-alert infotype="explain" title="">You should specify either DataDisks or SystemDisk, but you cannot specify both at the same time.</dx-alert>
	DataDisks []*DataDisk `json:"DataDisks,omitnil,omitempty" name:"DataDisks"`

	// Specifies whether to forcibly shut down a running instance. it is recommended to manually shut down a running instance first and then expand the instance disk. valid values:<br><li>true: forcibly shut down an instance after a normal shutdown fails.</li><br><li>false: do not forcibly shut down an instance after a normal shutdown fails.</li><br><br>default value: false.<br><br>forced shutdown is equivalent to turning off a physical computer's power switch. forced shutdown may cause data loss or file system corruption and should only be used when a server cannot be shut down normally.
	ForceStop *bool `json:"ForceStop,omitnil,omitempty" name:"ForceStop"`

	// System disk configuration information to be expanded. only supports specifying the purpose capacity of the disk to be expanded. only supports cloud disk expansion.
	// <dx-alert infotype="explain" title="">You should specify either DataDisks or SystemDisk, but you cannot specify both at the same time.</dx-alert>
	SystemDisk *SystemDisk `json:"SystemDisk,omitnil,omitempty" name:"SystemDisk"`

	// Whether the cloud disk is expanded online.
	ResizeOnline *bool `json:"ResizeOnline,omitnil,omitempty" name:"ResizeOnline"`
}

func (r *ResizeInstanceDisksRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ResizeInstanceDisksRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceId")
	delete(f, "DataDisks")
	delete(f, "ForceStop")
	delete(f, "SystemDisk")
	delete(f, "ResizeOnline")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ResizeInstanceDisksRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ResizeInstanceDisksResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ResizeInstanceDisksResponse struct {
	*tchttp.BaseResponse
	Response *ResizeInstanceDisksResponseParams `json:"Response"`
}

func (r *ResizeInstanceDisksResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ResizeInstanceDisksResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type RunAutomationServiceEnabled struct {
	// Whether to enable the TAT service. Valid values: <br><li>`TRUE`: yes;<br><li>`FALSE`: no<br><br>Default: `FALSE`.
	Enabled *bool `json:"Enabled,omitnil,omitempty" name:"Enabled"`
}

// Predefined struct for user
type RunInstancesRequestParams struct {
	// Instance [billing type](https://intl.cloud.tencent.com/document/product/213/2180?from_cn_redirect=1). <br><li>`PREPAID`: Monthly Subscription, used for at least one month <br><li>`POSTPAID_BY_HOUR`: Hourly-based pay-as-you-go <br><li>`CDHPAID`: [Dedicated CVM](https://www.tencentcloud.com/document/product/416/5068?lang=en&pg=) (associated with a dedicated host. Resource usage of the dedicated host is free of charge.) <br><li>`SPOTPAID`: [Spot instance](https://intl.cloud.tencent.com/document/product/213/17817)<br>Default value: `POSTPAID_BY_HOUR`.
	InstanceChargeType *string `json:"InstanceChargeType,omitnil,omitempty" name:"InstanceChargeType"`

	// Details of the monthly subscription, including the purchase period, auto-renewal. It is required if the `InstanceChargeType` is `PREPAID`.
	InstanceChargePrepaid *InstanceChargePrepaid `json:"InstanceChargePrepaid,omitnil,omitempty" name:"InstanceChargePrepaid"`

	// Location of the instance. You can use this parameter to specify the availability zone, project, and CDH (for dedicated CVMs).
	//  <b>Note: `Placement` is required when `LaunchTemplate` is not specified. If both the parameters are passed in, `Placement` prevails.</b>
	Placement *Placement `json:"Placement,omitnil,omitempty" name:"Placement"`

	// The instance model. 
	// <br><li>To view specific values for `POSTPAID_BY_HOUR` instances, you can call [DescribeInstanceTypeConfigs](https://intl.cloud.tencent.com/document/api/213/15749?from_cn_redirect=1) or refer to [Instance Types](https://intl.cloud.tencent.com/document/product/213/11518?from_cn_redirect=1). <br><li>For `CDHPAID` instances, the value of this parameter is in the format of `CDH_XCXG` based on the number of CPU cores and memory capacity. For example, if you want to create a CDH instance with a single-core CPU and 1 GB memory, specify this parameter as `CDH_1C1G`.
	InstanceType *string `json:"InstanceType,omitnil,omitempty" name:"InstanceType"`

	// The [image](https://intl.cloud.tencent.com/document/product/213/4940?from_cn_redirect=1) ID in the format of `img-xxx`. There are three types of images:<br/><li>Public images</li><li>Custom images</li><li>Shared images</li><br/>To check the image ID:<br/><li>For public images, custom images, and shared images, go to the [CVM console](https://console.cloud.tencent.com/cvm/image?rid=1&imageType=PUBLIC_IMAGE). </li><li>Call [DescribeImages](https://intl.cloud.tencent.com/document/api/213/15715?from_cn_redirect=1), pass in `InstanceType` to retrieve the list of images supported by the current model, and then find the `ImageId` in the response.</li>
	//  <b>Note: `ImageId` is required when `LaunchTemplate` is not specified. If both the parameters are passed in, `ImageId` prevails.</b>
	ImageId *string `json:"ImageId,omitnil,omitempty" name:"ImageId"`

	// System disk configuration of the instance. If this parameter is not specified, the default value will be used.
	SystemDisk *SystemDisk `json:"SystemDisk,omitnil,omitempty" name:"SystemDisk"`

	// The configuration information of instance data disks. If this parameter is not specified, no data disk will be purchased by default. When purchasing, you can specify 21 data disks, which can contain at most 1 LOCAL_BASIC or LOCAL_SSD data disk, and at most 20 CLOUD_BASIC, CLOUD_PREMIUM, or CLOUD_SSD data disks.
	DataDisks []*DataDisk `json:"DataDisks,omitnil,omitempty" name:"DataDisks"`

	// Configuration information of VPC. This parameter is used to specify VPC ID and subnet ID, etc. If a VPC IP is specified in this parameter, it indicates the primary ENI IP of each instance. The value of parameter InstanceCount must be the same as the number of VPC IPs, which cannot be greater than 20.
	VirtualPrivateCloud *VirtualPrivateCloud `json:"VirtualPrivateCloud,omitnil,omitempty" name:"VirtualPrivateCloud"`

	// Configuration of public network bandwidth. If this parameter is not specified, 0 Mbps will be used by default.
	InternetAccessible *InternetAccessible `json:"InternetAccessible,omitnil,omitempty" name:"InternetAccessible"`

	// The number of instances to be purchased. Value range for pay-as-you-go instances: [1, 100]. Default value: `1`. The specified number of instances to be purchased cannot exceed the remaining quota allowed for the user. For more information on the quota, see [Quota for CVM Instances](https://intl.cloud.tencent.com/document/product/213/2664).
	InstanceCount *int64 `json:"InstanceCount,omitnil,omitempty" name:"InstanceCount"`

	// Specifies the minimum number of instances to create. value range: positive integer not greater than InstanceCount.
	// Specifies the minimum purchasable quantity, guarantees to create at least MinCount instances, and creates InstanceCount instances as much as possible.
	// Insufficient inventory to meet the minimum purchasable quantity will trigger an error info indicating insufficient stock.
	// Only applicable to accounts, regions, and billing modes (annual/monthly subscription, pay-as-you-go, spot instance, exclusive sales) with partial support.
	MinCount *int64 `json:"MinCount,omitnil,omitempty" name:"MinCount"`

	// Instance display name.<br><li>If no instance display name is specified, it will display 'unnamed' by default.</li><li>Supports up to 128 characters (including pattern strings).</li><li>When purchasing multiple instances:.
	// - Specify a pattern string {R:x}: Generates a numeric sequence [x, x+n-1], where n represents the number of purchased instances. For example: Input server_{R:3}. When purchasing 1 instance, the display name is server_3; when purchasing 2 instances, the display names are server_3 and server_4.
	// - Specify a pattern string {R:x,F:y}: y indicates fixed digit (optional), value range [0,8], default value 0 means no fixed digit (equivalent to {R:x}). Automatically pads with zeros when digits are insufficient, for example: input server_{R:3,F:3}, when purchasing 2 instances, the instance display name is server_003, server_004. If digit count exceeds y (such as {R:99,F:2}), the actual number is used, for example: app_{R:99,F:2}, when purchasing 2 instances, the instance display name is app_99, app_100.
	// - Pattern strings must strictly follow the format {R:x,F:y} or {R:x}. Invalid formats (such as {}) are treated as plain text. Multiple pattern strings are supported.
	// - No pattern string specified: The display name is appended with suffix 1, 2...n, where n indicates the number of instances purchased. For example server_. When purchasing 2 instances, generates server_1 and server_2.
	InstanceName *string `json:"InstanceName,omitnil,omitempty" name:"InstanceName"`

	// Instance login settings. You can use this parameter to set the login method, password and key of the instance, or keep the original login settings of the image. If it's not specified, the user needs to set the login password using the "Reset password" option in the CVM console or calling the API `ResetInstancesPassword` to complete the creation of the CVM instance(s).
	LoginSettings *LoginSettings `json:"LoginSettings,omitnil,omitempty" name:"LoginSettings"`

	// Security group to which an instance belongs. obtain this parameter by calling the `SecurityGroupId` field in the return value of [DescribeSecurityGroups](https://www.tencentcloud.com/document/product/215/15808?from_search=1). if not specified, bind the default security group under the designated project. if the default security group does not exist, automatically create it.
	SecurityGroupIds []*string `json:"SecurityGroupIds,omitnil,omitempty" name:"SecurityGroupIds"`

	// Enhanced service. You can use this parameter to specify whether to enable services such as Anti-DDoS and Cloud Monitor. If this parameter is not specified, Cloud Monitor and Anti-DDoS are enabled for public images by default. However, for custom images and images from the marketplace, Anti-DDoS and Cloud Monitor are not enabled by default. The original services in the image will be retained.
	EnhancedService *EnhancedService `json:"EnhancedService,omitnil,omitempty" name:"EnhancedService"`

	// A unique string supplied by the client to ensure that the request is idempotent. Its maximum length is 64 ASCII characters. If this parameter is not specified, the idem-potency of the request cannot be guaranteed.
	ClientToken *string `json:"ClientToken,omitnil,omitempty" name:"ClientToken"`

	// Instance HostName.<br><li>Period (.) and hyphen (-) should not be used as the first or last character of the hostname, and should not be used consecutively.</li><br><li>Windows instances: The hostname should contain 2 to 15 characters, including letters (case insensitive), digits, and hyphens (-), does not support periods (.), and should not be all digits.</li><br><li>Instances of other types (such as Linux instances): The hostname should contain 2 to 60 characters, including multiple periods (.), with each segment between periods considered as one section. Each section can contain letters (case insensitive), digits, and hyphens (-).</li><br><li>When purchasing multiple instances:
	// - Specify a pattern string {R:x}: Generates a numeric sequence [x, x+n-1], where n represents the number of purchased instances. For example: Input server_{R:3}. When purchasing 1 instance, the hostname is server_3; when purchasing 2 instances, the hostnames are server_3 and server_4.
	// - Specify a pattern string {R:x,F:y}: y indicates fixed digit width (optional), valid range [0,8]. Default value 0 means no fixed width (equivalent to {R:x}). Insufficient digits are automatically padded with zeros. For example: Input server{R:3,F:3}. When purchasing 2 instances, the hostnames are server_003 and server_004. If the number of digits exceeds y (e.g., {R:99,F:2}), the actual digit count is used. For example: app{R:99,F:2}. When purchasing 2 instances, the hostnames are app_99 and app_100.
	// - Specify a pattern string {IP}: Automatically replace with the private IP address of the instance. For example: input node-{IP}, the instance hostname is node-10.0.12.8. Supports mixed use with serial number pattern strings, such as: input web-{IP}-{R:1}, when purchasing 2 instances, the instance hostnames are web-10.0.12.8-1 and web-10.0.12.9-2 respectively.
	// - Pattern strings must strictly follow the format {R:x,F:y}, {R:x}, or {IP}. Invalid formats (such as {}) are treated as plain text. Multiple pattern strings are supported.
	// - No specified pattern string: add suffix 1, 2...n to instance hostname, where n means the number of purchased instances, such as server_1, server_2 when purchasing 2 instances.
	// </li>
	HostName *string `json:"HostName,omitnil,omitempty" name:"HostName"`

	// Scheduled tasks. You can use this parameter to specify scheduled tasks for the instance. Only scheduled termination is supported.
	ActionTimer *ActionTimer `json:"ActionTimer,omitnil,omitempty" name:"ActionTimer"`

	// Placement group ID. You can only specify one.
	DisasterRecoverGroupIds []*string `json:"DisasterRecoverGroupIds,omitnil,omitempty" name:"DisasterRecoverGroupIds"`

	// List of tag description. By specifying this parameter, the tag can be bound to the corresponding CVM and CBS instances at the same time.
	TagSpecification []*TagSpecification `json:"TagSpecification,omitnil,omitempty" name:"TagSpecification"`

	// The market options of the instance.
	InstanceMarketOptions *InstanceMarketOptionsRequest `json:"InstanceMarketOptions,omitnil,omitempty" name:"InstanceMarketOptions"`

	// User data provided to the instance. This parameter needs to be encoded in base64 format with the maximum size of 16 KB. For more information on how to get the value of this parameter, see the commands you need to execute on startup for [Windows](https://intl.cloud.tencent.com/document/product/213/17526) or [Linux](https://intl.cloud.tencent.com/document/product/213/17525).
	UserData *string `json:"UserData,omitnil,omitempty" name:"UserData"`

	// Custom metadata. specifies the support for creating custom metadata key-value pairs when creating a CVM.
	// **Note: this field is in beta test.**.
	Metadata *Metadata `json:"Metadata,omitnil,omitempty" name:"Metadata"`

	// Whether the request is a dry run only.
	// `true`: dry run only. The request will not create instance(s). A dry run can check whether all the required parameters are specified, whether the request format is right, whether the request exceeds service limits, and whether the specified CVMs are available.
	// If the dry run fails, the corresponding error code will be returned.
	// If the dry run succeeds, the RequestId will be returned.
	// `false` (default value): Send a normal request and create instance(s) if all the requirements are met.
	DryRun *bool `json:"DryRun,omitnil,omitempty" name:"DryRun"`

	// Information about the CPU topology of an instance. If not specified, it is determined by system resources.
	CpuTopology *CpuTopology `json:"CpuTopology,omitnil,omitempty" name:"CpuTopology"`

	// CAM role name, which can be obtained from the `roleName` field in the response of the [`DescribeRoleList`](https://intl.cloud.tencent.com/document/product/598/36223?from_cn_redirect=1) API.
	CamRoleName *string `json:"CamRoleName,omitnil,omitempty" name:"CamRoleName"`

	// HPC cluster ID. The HPC cluster must and can only be specified for a high-performance computing instance.
	HpcClusterId *string `json:"HpcClusterId,omitnil,omitempty" name:"HpcClusterId"`

	// Instance launch template.
	LaunchTemplate *LaunchTemplate `json:"LaunchTemplate,omitnil,omitempty" name:"LaunchTemplate"`

	// Specify the ID of the dedicated cluster where the CVM is created.
	DedicatedClusterId *string `json:"DedicatedClusterId,omitnil,omitempty" name:"DedicatedClusterId"`

	// Specify the CHC physical server that used to create the CHC CVM.
	ChcIds []*string `json:"ChcIds,omitnil,omitempty" name:"ChcIds"`

	// Instance termination protection flag, indicating whether an instance is allowed to be deleted through an API. Valid values:<br><li>true: Instance protection is enabled, and the instance is not allowed to be deleted through the API.</li><br><li>false: Instance protection is disabled, and the instance is allowed to be deleted through the API.</li><br><br>Default value: false.
	DisableApiTermination *bool `json:"DisableApiTermination,omitnil,omitempty" name:"DisableApiTermination"`

	// Whether the instance enables jumbo frames. valid values:<br><li/> true: means the instance enables jumbo frames. only models supporting jumbo frames can be set to true.<br><li/> false: means the instance disables jumbo frames. only models supporting jumbo frames can be set to false.<br> instance specifications supporting jumbo frames: [instance specifications](https://www.tencentcloud.com/document/product/213/11518?lang=en&pg=).
	EnableJumboFrame *bool `json:"EnableJumboFrame,omitnil,omitempty" name:"EnableJumboFrame"`
}

type RunInstancesRequest struct {
	*tchttp.BaseRequest
	
	// Instance [billing type](https://intl.cloud.tencent.com/document/product/213/2180?from_cn_redirect=1). <br><li>`PREPAID`: Monthly Subscription, used for at least one month <br><li>`POSTPAID_BY_HOUR`: Hourly-based pay-as-you-go <br><li>`CDHPAID`: [Dedicated CVM](https://www.tencentcloud.com/document/product/416/5068?lang=en&pg=) (associated with a dedicated host. Resource usage of the dedicated host is free of charge.) <br><li>`SPOTPAID`: [Spot instance](https://intl.cloud.tencent.com/document/product/213/17817)<br>Default value: `POSTPAID_BY_HOUR`.
	InstanceChargeType *string `json:"InstanceChargeType,omitnil,omitempty" name:"InstanceChargeType"`

	// Details of the monthly subscription, including the purchase period, auto-renewal. It is required if the `InstanceChargeType` is `PREPAID`.
	InstanceChargePrepaid *InstanceChargePrepaid `json:"InstanceChargePrepaid,omitnil,omitempty" name:"InstanceChargePrepaid"`

	// Location of the instance. You can use this parameter to specify the availability zone, project, and CDH (for dedicated CVMs).
	//  <b>Note: `Placement` is required when `LaunchTemplate` is not specified. If both the parameters are passed in, `Placement` prevails.</b>
	Placement *Placement `json:"Placement,omitnil,omitempty" name:"Placement"`

	// The instance model. 
	// <br><li>To view specific values for `POSTPAID_BY_HOUR` instances, you can call [DescribeInstanceTypeConfigs](https://intl.cloud.tencent.com/document/api/213/15749?from_cn_redirect=1) or refer to [Instance Types](https://intl.cloud.tencent.com/document/product/213/11518?from_cn_redirect=1). <br><li>For `CDHPAID` instances, the value of this parameter is in the format of `CDH_XCXG` based on the number of CPU cores and memory capacity. For example, if you want to create a CDH instance with a single-core CPU and 1 GB memory, specify this parameter as `CDH_1C1G`.
	InstanceType *string `json:"InstanceType,omitnil,omitempty" name:"InstanceType"`

	// The [image](https://intl.cloud.tencent.com/document/product/213/4940?from_cn_redirect=1) ID in the format of `img-xxx`. There are three types of images:<br/><li>Public images</li><li>Custom images</li><li>Shared images</li><br/>To check the image ID:<br/><li>For public images, custom images, and shared images, go to the [CVM console](https://console.cloud.tencent.com/cvm/image?rid=1&imageType=PUBLIC_IMAGE). </li><li>Call [DescribeImages](https://intl.cloud.tencent.com/document/api/213/15715?from_cn_redirect=1), pass in `InstanceType` to retrieve the list of images supported by the current model, and then find the `ImageId` in the response.</li>
	//  <b>Note: `ImageId` is required when `LaunchTemplate` is not specified. If both the parameters are passed in, `ImageId` prevails.</b>
	ImageId *string `json:"ImageId,omitnil,omitempty" name:"ImageId"`

	// System disk configuration of the instance. If this parameter is not specified, the default value will be used.
	SystemDisk *SystemDisk `json:"SystemDisk,omitnil,omitempty" name:"SystemDisk"`

	// The configuration information of instance data disks. If this parameter is not specified, no data disk will be purchased by default. When purchasing, you can specify 21 data disks, which can contain at most 1 LOCAL_BASIC or LOCAL_SSD data disk, and at most 20 CLOUD_BASIC, CLOUD_PREMIUM, or CLOUD_SSD data disks.
	DataDisks []*DataDisk `json:"DataDisks,omitnil,omitempty" name:"DataDisks"`

	// Configuration information of VPC. This parameter is used to specify VPC ID and subnet ID, etc. If a VPC IP is specified in this parameter, it indicates the primary ENI IP of each instance. The value of parameter InstanceCount must be the same as the number of VPC IPs, which cannot be greater than 20.
	VirtualPrivateCloud *VirtualPrivateCloud `json:"VirtualPrivateCloud,omitnil,omitempty" name:"VirtualPrivateCloud"`

	// Configuration of public network bandwidth. If this parameter is not specified, 0 Mbps will be used by default.
	InternetAccessible *InternetAccessible `json:"InternetAccessible,omitnil,omitempty" name:"InternetAccessible"`

	// The number of instances to be purchased. Value range for pay-as-you-go instances: [1, 100]. Default value: `1`. The specified number of instances to be purchased cannot exceed the remaining quota allowed for the user. For more information on the quota, see [Quota for CVM Instances](https://intl.cloud.tencent.com/document/product/213/2664).
	InstanceCount *int64 `json:"InstanceCount,omitnil,omitempty" name:"InstanceCount"`

	// Specifies the minimum number of instances to create. value range: positive integer not greater than InstanceCount.
	// Specifies the minimum purchasable quantity, guarantees to create at least MinCount instances, and creates InstanceCount instances as much as possible.
	// Insufficient inventory to meet the minimum purchasable quantity will trigger an error info indicating insufficient stock.
	// Only applicable to accounts, regions, and billing modes (annual/monthly subscription, pay-as-you-go, spot instance, exclusive sales) with partial support.
	MinCount *int64 `json:"MinCount,omitnil,omitempty" name:"MinCount"`

	// Instance display name.<br><li>If no instance display name is specified, it will display 'unnamed' by default.</li><li>Supports up to 128 characters (including pattern strings).</li><li>When purchasing multiple instances:.
	// - Specify a pattern string {R:x}: Generates a numeric sequence [x, x+n-1], where n represents the number of purchased instances. For example: Input server_{R:3}. When purchasing 1 instance, the display name is server_3; when purchasing 2 instances, the display names are server_3 and server_4.
	// - Specify a pattern string {R:x,F:y}: y indicates fixed digit (optional), value range [0,8], default value 0 means no fixed digit (equivalent to {R:x}). Automatically pads with zeros when digits are insufficient, for example: input server_{R:3,F:3}, when purchasing 2 instances, the instance display name is server_003, server_004. If digit count exceeds y (such as {R:99,F:2}), the actual number is used, for example: app_{R:99,F:2}, when purchasing 2 instances, the instance display name is app_99, app_100.
	// - Pattern strings must strictly follow the format {R:x,F:y} or {R:x}. Invalid formats (such as {}) are treated as plain text. Multiple pattern strings are supported.
	// - No pattern string specified: The display name is appended with suffix 1, 2...n, where n indicates the number of instances purchased. For example server_. When purchasing 2 instances, generates server_1 and server_2.
	InstanceName *string `json:"InstanceName,omitnil,omitempty" name:"InstanceName"`

	// Instance login settings. You can use this parameter to set the login method, password and key of the instance, or keep the original login settings of the image. If it's not specified, the user needs to set the login password using the "Reset password" option in the CVM console or calling the API `ResetInstancesPassword` to complete the creation of the CVM instance(s).
	LoginSettings *LoginSettings `json:"LoginSettings,omitnil,omitempty" name:"LoginSettings"`

	// Security group to which an instance belongs. obtain this parameter by calling the `SecurityGroupId` field in the return value of [DescribeSecurityGroups](https://www.tencentcloud.com/document/product/215/15808?from_search=1). if not specified, bind the default security group under the designated project. if the default security group does not exist, automatically create it.
	SecurityGroupIds []*string `json:"SecurityGroupIds,omitnil,omitempty" name:"SecurityGroupIds"`

	// Enhanced service. You can use this parameter to specify whether to enable services such as Anti-DDoS and Cloud Monitor. If this parameter is not specified, Cloud Monitor and Anti-DDoS are enabled for public images by default. However, for custom images and images from the marketplace, Anti-DDoS and Cloud Monitor are not enabled by default. The original services in the image will be retained.
	EnhancedService *EnhancedService `json:"EnhancedService,omitnil,omitempty" name:"EnhancedService"`

	// A unique string supplied by the client to ensure that the request is idempotent. Its maximum length is 64 ASCII characters. If this parameter is not specified, the idem-potency of the request cannot be guaranteed.
	ClientToken *string `json:"ClientToken,omitnil,omitempty" name:"ClientToken"`

	// Instance HostName.<br><li>Period (.) and hyphen (-) should not be used as the first or last character of the hostname, and should not be used consecutively.</li><br><li>Windows instances: The hostname should contain 2 to 15 characters, including letters (case insensitive), digits, and hyphens (-), does not support periods (.), and should not be all digits.</li><br><li>Instances of other types (such as Linux instances): The hostname should contain 2 to 60 characters, including multiple periods (.), with each segment between periods considered as one section. Each section can contain letters (case insensitive), digits, and hyphens (-).</li><br><li>When purchasing multiple instances:
	// - Specify a pattern string {R:x}: Generates a numeric sequence [x, x+n-1], where n represents the number of purchased instances. For example: Input server_{R:3}. When purchasing 1 instance, the hostname is server_3; when purchasing 2 instances, the hostnames are server_3 and server_4.
	// - Specify a pattern string {R:x,F:y}: y indicates fixed digit width (optional), valid range [0,8]. Default value 0 means no fixed width (equivalent to {R:x}). Insufficient digits are automatically padded with zeros. For example: Input server{R:3,F:3}. When purchasing 2 instances, the hostnames are server_003 and server_004. If the number of digits exceeds y (e.g., {R:99,F:2}), the actual digit count is used. For example: app{R:99,F:2}. When purchasing 2 instances, the hostnames are app_99 and app_100.
	// - Specify a pattern string {IP}: Automatically replace with the private IP address of the instance. For example: input node-{IP}, the instance hostname is node-10.0.12.8. Supports mixed use with serial number pattern strings, such as: input web-{IP}-{R:1}, when purchasing 2 instances, the instance hostnames are web-10.0.12.8-1 and web-10.0.12.9-2 respectively.
	// - Pattern strings must strictly follow the format {R:x,F:y}, {R:x}, or {IP}. Invalid formats (such as {}) are treated as plain text. Multiple pattern strings are supported.
	// - No specified pattern string: add suffix 1, 2...n to instance hostname, where n means the number of purchased instances, such as server_1, server_2 when purchasing 2 instances.
	// </li>
	HostName *string `json:"HostName,omitnil,omitempty" name:"HostName"`

	// Scheduled tasks. You can use this parameter to specify scheduled tasks for the instance. Only scheduled termination is supported.
	ActionTimer *ActionTimer `json:"ActionTimer,omitnil,omitempty" name:"ActionTimer"`

	// Placement group ID. You can only specify one.
	DisasterRecoverGroupIds []*string `json:"DisasterRecoverGroupIds,omitnil,omitempty" name:"DisasterRecoverGroupIds"`

	// List of tag description. By specifying this parameter, the tag can be bound to the corresponding CVM and CBS instances at the same time.
	TagSpecification []*TagSpecification `json:"TagSpecification,omitnil,omitempty" name:"TagSpecification"`

	// The market options of the instance.
	InstanceMarketOptions *InstanceMarketOptionsRequest `json:"InstanceMarketOptions,omitnil,omitempty" name:"InstanceMarketOptions"`

	// User data provided to the instance. This parameter needs to be encoded in base64 format with the maximum size of 16 KB. For more information on how to get the value of this parameter, see the commands you need to execute on startup for [Windows](https://intl.cloud.tencent.com/document/product/213/17526) or [Linux](https://intl.cloud.tencent.com/document/product/213/17525).
	UserData *string `json:"UserData,omitnil,omitempty" name:"UserData"`

	// Custom metadata. specifies the support for creating custom metadata key-value pairs when creating a CVM.
	// **Note: this field is in beta test.**.
	Metadata *Metadata `json:"Metadata,omitnil,omitempty" name:"Metadata"`

	// Whether the request is a dry run only.
	// `true`: dry run only. The request will not create instance(s). A dry run can check whether all the required parameters are specified, whether the request format is right, whether the request exceeds service limits, and whether the specified CVMs are available.
	// If the dry run fails, the corresponding error code will be returned.
	// If the dry run succeeds, the RequestId will be returned.
	// `false` (default value): Send a normal request and create instance(s) if all the requirements are met.
	DryRun *bool `json:"DryRun,omitnil,omitempty" name:"DryRun"`

	// Information about the CPU topology of an instance. If not specified, it is determined by system resources.
	CpuTopology *CpuTopology `json:"CpuTopology,omitnil,omitempty" name:"CpuTopology"`

	// CAM role name, which can be obtained from the `roleName` field in the response of the [`DescribeRoleList`](https://intl.cloud.tencent.com/document/product/598/36223?from_cn_redirect=1) API.
	CamRoleName *string `json:"CamRoleName,omitnil,omitempty" name:"CamRoleName"`

	// HPC cluster ID. The HPC cluster must and can only be specified for a high-performance computing instance.
	HpcClusterId *string `json:"HpcClusterId,omitnil,omitempty" name:"HpcClusterId"`

	// Instance launch template.
	LaunchTemplate *LaunchTemplate `json:"LaunchTemplate,omitnil,omitempty" name:"LaunchTemplate"`

	// Specify the ID of the dedicated cluster where the CVM is created.
	DedicatedClusterId *string `json:"DedicatedClusterId,omitnil,omitempty" name:"DedicatedClusterId"`

	// Specify the CHC physical server that used to create the CHC CVM.
	ChcIds []*string `json:"ChcIds,omitnil,omitempty" name:"ChcIds"`

	// Instance termination protection flag, indicating whether an instance is allowed to be deleted through an API. Valid values:<br><li>true: Instance protection is enabled, and the instance is not allowed to be deleted through the API.</li><br><li>false: Instance protection is disabled, and the instance is allowed to be deleted through the API.</li><br><br>Default value: false.
	DisableApiTermination *bool `json:"DisableApiTermination,omitnil,omitempty" name:"DisableApiTermination"`

	// Whether the instance enables jumbo frames. valid values:<br><li/> true: means the instance enables jumbo frames. only models supporting jumbo frames can be set to true.<br><li/> false: means the instance disables jumbo frames. only models supporting jumbo frames can be set to false.<br> instance specifications supporting jumbo frames: [instance specifications](https://www.tencentcloud.com/document/product/213/11518?lang=en&pg=).
	EnableJumboFrame *bool `json:"EnableJumboFrame,omitnil,omitempty" name:"EnableJumboFrame"`
}

func (r *RunInstancesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *RunInstancesRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceChargeType")
	delete(f, "InstanceChargePrepaid")
	delete(f, "Placement")
	delete(f, "InstanceType")
	delete(f, "ImageId")
	delete(f, "SystemDisk")
	delete(f, "DataDisks")
	delete(f, "VirtualPrivateCloud")
	delete(f, "InternetAccessible")
	delete(f, "InstanceCount")
	delete(f, "MinCount")
	delete(f, "InstanceName")
	delete(f, "LoginSettings")
	delete(f, "SecurityGroupIds")
	delete(f, "EnhancedService")
	delete(f, "ClientToken")
	delete(f, "HostName")
	delete(f, "ActionTimer")
	delete(f, "DisasterRecoverGroupIds")
	delete(f, "TagSpecification")
	delete(f, "InstanceMarketOptions")
	delete(f, "UserData")
	delete(f, "Metadata")
	delete(f, "DryRun")
	delete(f, "CpuTopology")
	delete(f, "CamRoleName")
	delete(f, "HpcClusterId")
	delete(f, "LaunchTemplate")
	delete(f, "DedicatedClusterId")
	delete(f, "ChcIds")
	delete(f, "DisableApiTermination")
	delete(f, "EnableJumboFrame")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "RunInstancesRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type RunInstancesResponseParams struct {
	// If you use this API to create instance(s), this parameter will be returned, representing one or more instance IDs. Retuning the instance ID list does not necessarily mean that the instance(s) were created successfully. To check whether the instance(s) were created successfully, you can call [DescribeInstances](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1) and check the status of the instances in `InstancesSet` in the response. If the status of an instance changes from "PENDING" to "RUNNING", it means that the instance has been created successfully.
	InstanceIdSet []*string `json:"InstanceIdSet,omitnil,omitempty" name:"InstanceIdSet"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type RunInstancesResponse struct {
	*tchttp.BaseResponse
	Response *RunInstancesResponseParams `json:"Response"`
}

func (r *RunInstancesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *RunInstancesResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type RunMonitorServiceEnabled struct {
	// Whether to enable the cloud monitor service. value ranges from: <li>true: indicates enabling the cloud monitor service</li> <li>false: indicates disabling the cloud monitor service</li> default value: true.
	Enabled *bool `json:"Enabled,omitnil,omitempty" name:"Enabled"`
}

type RunSecurityServiceEnabled struct {
	// Whether to enable [Cloud Security](https://intl.cloud.tencent.com/document/product/296?from_cn_redirect=1). Valid values: <br><li>TRUE: enable Cloud Security <br><li>FALSE: do not enable Cloud Security <br><br>Default value: TRUE.
	Enabled *bool `json:"Enabled,omitnil,omitempty" name:"Enabled"`
}

type SharePermission struct {
	// Time when an image was shared.
	CreatedTime *string `json:"CreatedTime,omitnil,omitempty" name:"CreatedTime"`

	// ID of the account with which the image is shared.
	AccountId *string `json:"AccountId,omitnil,omitempty" name:"AccountId"`
}

type Snapshot struct {
	// Snapshot ID.
	SnapshotId *string `json:"SnapshotId,omitnil,omitempty" name:"SnapshotId"`

	// Type of the cloud disk used to create the snapshot. Valid values:
	// SYSTEM_DISK: system disk
	// DATA_DISK: data disk
	DiskUsage *string `json:"DiskUsage,omitnil,omitempty" name:"DiskUsage"`

	// Size of the cloud disk used to create the snapshot; unit: GB.
	DiskSize *int64 `json:"DiskSize,omitnil,omitempty" name:"DiskSize"`
}

type SpotMarketOptions struct {
	// Bid price.
	MaxPrice *string `json:"MaxPrice,omitnil,omitempty" name:"MaxPrice"`

	// Bid request type. valid values: one-time. currently, only the one-time type is supported.
	SpotInstanceType *string `json:"SpotInstanceType,omitnil,omitempty" name:"SpotInstanceType"`
}

// Predefined struct for user
type StartInstancesRequestParams struct {
	// Instance ID(s). To obtain the instance IDs, you can call [`DescribeInstances`](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1) and look for `InstanceId` in the response. The maximum number of instances in each request is 100.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`
}

type StartInstancesRequest struct {
	*tchttp.BaseRequest
	
	// Instance ID(s). To obtain the instance IDs, you can call [`DescribeInstances`](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1) and look for `InstanceId` in the response. The maximum number of instances in each request is 100.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`
}

func (r *StartInstancesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *StartInstancesRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceIds")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "StartInstancesRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type StartInstancesResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type StartInstancesResponse struct {
	*tchttp.BaseResponse
	Response *StartInstancesResponseParams `json:"Response"`
}

func (r *StartInstancesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *StartInstancesResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type StopInstancesRequestParams struct {
	// Instance ID(s). To obtain the instance IDs, you can call [`DescribeInstances`](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1) and look for `InstanceId` in the response. The maximum number of instances in each request is 100.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// (Disused. Please use `StopType` instead.) Whether to forcibly shut down an instance after a normal shutdown fails. Valid values: <br><li>`TRUE`: yes;<br><li>`FALSE`: no<br><br>Default value: `FALSE`. 
	//
	// Deprecated: ForceStop is deprecated.
	ForceStop *bool `json:"ForceStop,omitnil,omitempty" name:"ForceStop"`

	// Instance shutdown mode. Valid values: <br><li>SOFT_FIRST: perform a soft shutdown first, and force shut down the instance if the soft shutdown fails <br><li>HARD: force shut down the instance directly <br><li>SOFT: soft shutdown only <br>Default value: SOFT.
	StopType *string `json:"StopType,omitnil,omitempty" name:"StopType"`

	// Billing method of a pay-as-you-go instance after shutdown.
	// Valid values: <br><li>KEEP_CHARGING: billing continues after shutdown <br><li>STOP_CHARGING: billing stops after shutdown <br>Default value: KEEP_CHARGING.
	// This parameter is only valid for some pay-as-you-go instances using cloud disks. For more information, see [No charges when shut down for pay-as-you-go instances](https://intl.cloud.tencent.com/document/product/213/19918).
	StoppedMode *string `json:"StoppedMode,omitnil,omitempty" name:"StoppedMode"`
}

type StopInstancesRequest struct {
	*tchttp.BaseRequest
	
	// Instance ID(s). To obtain the instance IDs, you can call [`DescribeInstances`](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1) and look for `InstanceId` in the response. The maximum number of instances in each request is 100.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// (Disused. Please use `StopType` instead.) Whether to forcibly shut down an instance after a normal shutdown fails. Valid values: <br><li>`TRUE`: yes;<br><li>`FALSE`: no<br><br>Default value: `FALSE`. 
	ForceStop *bool `json:"ForceStop,omitnil,omitempty" name:"ForceStop"`

	// Instance shutdown mode. Valid values: <br><li>SOFT_FIRST: perform a soft shutdown first, and force shut down the instance if the soft shutdown fails <br><li>HARD: force shut down the instance directly <br><li>SOFT: soft shutdown only <br>Default value: SOFT.
	StopType *string `json:"StopType,omitnil,omitempty" name:"StopType"`

	// Billing method of a pay-as-you-go instance after shutdown.
	// Valid values: <br><li>KEEP_CHARGING: billing continues after shutdown <br><li>STOP_CHARGING: billing stops after shutdown <br>Default value: KEEP_CHARGING.
	// This parameter is only valid for some pay-as-you-go instances using cloud disks. For more information, see [No charges when shut down for pay-as-you-go instances](https://intl.cloud.tencent.com/document/product/213/19918).
	StoppedMode *string `json:"StoppedMode,omitnil,omitempty" name:"StoppedMode"`
}

func (r *StopInstancesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *StopInstancesRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceIds")
	delete(f, "ForceStop")
	delete(f, "StopType")
	delete(f, "StoppedMode")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "StopInstancesRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type StopInstancesResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type StopInstancesResponse struct {
	*tchttp.BaseResponse
	Response *StopInstancesResponseParams `json:"Response"`
}

func (r *StopInstancesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *StopInstancesResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type StorageBlock struct {
	// HDD LOCAL storage type specifies the value: LOCAL_PRO.
	Type *string `json:"Type,omitnil,omitempty" name:"Type"`

	// Specifies the minimum HDD local storage capacity. measurement unit: GiB.
	MinSize *int64 `json:"MinSize,omitnil,omitempty" name:"MinSize"`

	// Specifies the maximum capacity of HDD local storage. measurement unit: GiB.
	MaxSize *int64 `json:"MaxSize,omitnil,omitempty" name:"MaxSize"`
}

type SyncImage struct {
	// Image ID
	ImageId *string `json:"ImageId,omitnil,omitempty" name:"ImageId"`

	// Region
	Region *string `json:"Region,omitnil,omitempty" name:"Region"`
}

// Predefined struct for user
type SyncImagesRequestParams struct {
	// Image ID list. You can obtain the image IDs in the following ways:<br><li>Call the [DescribeImages](https://intl.cloud.tencent.com/document/api/213/15715?from_cn_redirect=1) API and find the value of `ImageId` in the response.</li><li>Obtain the image IDs in the [Image console](https://console.cloud.tencent.com/cvm/image).<br>The image IDs should meet the following requirement:</li><li>The image ID should correspond to an image in the `NORMAL` state.</li>For more information on image status, see [Image Data Table](https://intl.cloud.tencent.com/document/product/213/15753?from_cn_redirect=1#Image).
	ImageIds []*string `json:"ImageIds,omitnil,omitempty" name:"ImageIds"`

	// List of target synchronization regions, which should meet the following requirements:<br><li>It should be a valid region.</li><li>If it is a custom image, the target synchronization region cannot be the source region.</li><li>If it is a shared image, the target synchronization region only supports the source region, meaning the shared image will be copied as a custom image in the source region.</li><li>Partial region synchronization is not supported currently. For details, see [Copying Images](https://intl.cloud.tencent.com/document/product/213/4943?from_cn_redirect=1#.E5.A4.8D.E5.88.B6.E8.AF.B4.E6.98.8E).</li>For specific regional parameters, refer to [Region](https://intl.cloud.tencent.com/document/product/213/6091?from_cn_redirect=1).
	DestinationRegions []*string `json:"DestinationRegions,omitnil,omitempty" name:"DestinationRegions"`

	// Checks whether image synchronization can be initiated.
	// 
	// Default value: false.
	DryRun *bool `json:"DryRun,omitnil,omitempty" name:"DryRun"`

	// Target image name. By default, the source image name is used.
	ImageName *string `json:"ImageName,omitnil,omitempty" name:"ImageName"`

	// Whether to return the ID of the image created in the target region.
	// 
	// Default value: false.
	ImageSetRequired *bool `json:"ImageSetRequired,omitnil,omitempty" name:"ImageSetRequired"`

	// Whether to synchronize as an encrypted custom image.
	// Default value is `false`.
	// Synchronization to an encrypted custom image is only supported within the same region.
	Encrypt *bool `json:"Encrypt,omitnil,omitempty" name:"Encrypt"`

	// KMS key ID used when synchronizing to an encrypted custom image. 
	// This parameter is valid only synchronizing to an encrypted image.
	// If KmsKeyId is not specified, the default CBS cloud product KMS key is used.
	KmsKeyId *string `json:"KmsKeyId,omitnil,omitempty" name:"KmsKeyId"`
}

type SyncImagesRequest struct {
	*tchttp.BaseRequest
	
	// Image ID list. You can obtain the image IDs in the following ways:<br><li>Call the [DescribeImages](https://intl.cloud.tencent.com/document/api/213/15715?from_cn_redirect=1) API and find the value of `ImageId` in the response.</li><li>Obtain the image IDs in the [Image console](https://console.cloud.tencent.com/cvm/image).<br>The image IDs should meet the following requirement:</li><li>The image ID should correspond to an image in the `NORMAL` state.</li>For more information on image status, see [Image Data Table](https://intl.cloud.tencent.com/document/product/213/15753?from_cn_redirect=1#Image).
	ImageIds []*string `json:"ImageIds,omitnil,omitempty" name:"ImageIds"`

	// List of target synchronization regions, which should meet the following requirements:<br><li>It should be a valid region.</li><li>If it is a custom image, the target synchronization region cannot be the source region.</li><li>If it is a shared image, the target synchronization region only supports the source region, meaning the shared image will be copied as a custom image in the source region.</li><li>Partial region synchronization is not supported currently. For details, see [Copying Images](https://intl.cloud.tencent.com/document/product/213/4943?from_cn_redirect=1#.E5.A4.8D.E5.88.B6.E8.AF.B4.E6.98.8E).</li>For specific regional parameters, refer to [Region](https://intl.cloud.tencent.com/document/product/213/6091?from_cn_redirect=1).
	DestinationRegions []*string `json:"DestinationRegions,omitnil,omitempty" name:"DestinationRegions"`

	// Checks whether image synchronization can be initiated.
	// 
	// Default value: false.
	DryRun *bool `json:"DryRun,omitnil,omitempty" name:"DryRun"`

	// Target image name. By default, the source image name is used.
	ImageName *string `json:"ImageName,omitnil,omitempty" name:"ImageName"`

	// Whether to return the ID of the image created in the target region.
	// 
	// Default value: false.
	ImageSetRequired *bool `json:"ImageSetRequired,omitnil,omitempty" name:"ImageSetRequired"`

	// Whether to synchronize as an encrypted custom image.
	// Default value is `false`.
	// Synchronization to an encrypted custom image is only supported within the same region.
	Encrypt *bool `json:"Encrypt,omitnil,omitempty" name:"Encrypt"`

	// KMS key ID used when synchronizing to an encrypted custom image. 
	// This parameter is valid only synchronizing to an encrypted image.
	// If KmsKeyId is not specified, the default CBS cloud product KMS key is used.
	KmsKeyId *string `json:"KmsKeyId,omitnil,omitempty" name:"KmsKeyId"`
}

func (r *SyncImagesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *SyncImagesRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ImageIds")
	delete(f, "DestinationRegions")
	delete(f, "DryRun")
	delete(f, "ImageName")
	delete(f, "ImageSetRequired")
	delete(f, "Encrypt")
	delete(f, "KmsKeyId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "SyncImagesRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type SyncImagesResponseParams struct {
	// ID of the image created in the destination region
	ImageSet []*SyncImage `json:"ImageSet,omitnil,omitempty" name:"ImageSet"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type SyncImagesResponse struct {
	*tchttp.BaseResponse
	Response *SyncImagesResponseParams `json:"Response"`
}

func (r *SyncImagesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *SyncImagesResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type SystemDisk struct {
	// Specifies the system disk type. for the restrictions on the system disk type, refer to [storage overview](https://www.tencentcloud.com/document/product/362/31636). value range:<br>
	// <li>LOCAL_BASIC: Local SATA disk</li>
	// <li>LOCAL_SSD: Local NVMe SSD</li>
	// <li>CLOUD_BASIC: Cloud SATA disk</li>
	// <li>CLOUD_SSD: Cloud SSD</li>
	// <li>CLOUD_PREMIUM: Premium SSD</li>
	// <li>CLOUD_BSSD: Balanced SSD</li>
	// <li>CLOUD_HSSD: Enhanced SSD</li>
	// <li>CLOUD_TSSD: Tremendous SSD</li>
	// Default value: Current disk types with inventory available.
	DiskType *string `json:"DiskType,omitnil,omitempty" name:"DiskType"`

	// System disk ID.
	// Currently, this parameter is only used for response parameters in query apis such as [DescribeInstances](https://www.tencentcloud.com/document/api/213/33258) and is not applicable to request parameters in write apis such as [RunInstances](https://www.tencentcloud.com/document/api/213/33237).
	DiskId *string `json:"DiskId,omitnil,omitempty" name:"DiskId"`

	// System disk size; unit: GiB; default value: 50 GiB.
	DiskSize *int64 `json:"DiskSize,omitnil,omitempty" name:"DiskSize"`

	// Specifies the dedicated cluster ID belonging to.
	// Note: This field may return null, indicating that no valid value is found.
	CdcId *string `json:"CdcId,omitnil,omitempty" name:"CdcId"`

	// Disk name, which specifies a length not exceeding 128 characters.
	DiskName *string `json:"DiskName,omitnil,omitempty" name:"DiskName"`
}

type Tag struct {
	// Tag key
	Key *string `json:"Key,omitnil,omitempty" name:"Key"`

	// Tag value
	Value *string `json:"Value,omitnil,omitempty" name:"Value"`
}

type TagSpecification struct {
	// Specifies the resource type for Tag binding. valid values: "instance" (cloud virtual machine), "host" (cdh), "image" (mirror), "keypair" (key), "ps" (placement group), "hpc" (hyper computing cluster).
	ResourceType *string `json:"ResourceType,omitnil,omitempty" name:"ResourceType"`

	// Tag pair list
	Tags []*Tag `json:"Tags,omitnil,omitempty" name:"Tags"`
}

type TargetOS struct {
	// Target operating system type.
	TargetOSType *string `json:"TargetOSType,omitnil,omitempty" name:"TargetOSType"`

	// Target operating system version.
	TargetOSVersion *string `json:"TargetOSVersion,omitnil,omitempty" name:"TargetOSVersion"`
}

// Predefined struct for user
type TerminateInstancesRequestParams struct {
	// One or more instance IDs to be operated. You can obtain the instance ID through the `InstanceId` in the return value from the API [DescribeInstances](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1). The maximum number of instances per request is 100.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// Release an Elastic IP. Under EIP 2.0, only the first EIP on the primary network interface can be released, and currently supported release types are limited to HighQualityEIP, AntiDDoSEIP, EIPv6, and HighQualityEIPv6.
	// Default value:  `false`.
	// 
	// This feature is currently in gradually released phase. To access it, please contact us.
	ReleaseAddress *bool `json:"ReleaseAddress,omitnil,omitempty" name:"ReleaseAddress"`

	// Whether to release a monthly subscription data disk mounted on an instance. true: Release the data disk along with termination of the instance. false: Only terminate the instance.
	// Default value: `false`.
	ReleasePrepaidDataDisks *bool `json:"ReleasePrepaidDataDisks,omitnil,omitempty" name:"ReleasePrepaidDataDisks"`
}

type TerminateInstancesRequest struct {
	*tchttp.BaseRequest
	
	// One or more instance IDs to be operated. You can obtain the instance ID through the `InstanceId` in the return value from the API [DescribeInstances](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1). The maximum number of instances per request is 100.
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// Release an Elastic IP. Under EIP 2.0, only the first EIP on the primary network interface can be released, and currently supported release types are limited to HighQualityEIP, AntiDDoSEIP, EIPv6, and HighQualityEIPv6.
	// Default value:  `false`.
	// 
	// This feature is currently in gradually released phase. To access it, please contact us.
	ReleaseAddress *bool `json:"ReleaseAddress,omitnil,omitempty" name:"ReleaseAddress"`

	// Whether to release a monthly subscription data disk mounted on an instance. true: Release the data disk along with termination of the instance. false: Only terminate the instance.
	// Default value: `false`.
	ReleasePrepaidDataDisks *bool `json:"ReleasePrepaidDataDisks,omitnil,omitempty" name:"ReleasePrepaidDataDisks"`
}

func (r *TerminateInstancesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *TerminateInstancesRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceIds")
	delete(f, "ReleaseAddress")
	delete(f, "ReleasePrepaidDataDisks")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "TerminateInstancesRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type TerminateInstancesResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type TerminateInstancesResponse struct {
	*tchttp.BaseResponse
	Response *TerminateInstancesResponseParams `json:"Response"`
}

func (r *TerminateInstancesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *TerminateInstancesResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type VirtualPrivateCloud struct {
	// vpc ID, such as `vpc-xxx`. valid vpc ids can be queried by logging in to the console (https://console.cloud.tencent.com/vpc/vpc?rid=1) or by calling the API [DescribeVpcs](https://www.tencentcloud.com/document/product/215/15778?lang=en) and obtaining the `VpcId` field from the API response. if both VpcId and SubnetId are input as `DEFAULT` when creating an instance, the DEFAULT vpc network will be forcibly used.
	VpcId *string `json:"VpcId,omitnil,omitempty" name:"VpcId"`

	// vpc subnet ID, in the form of `subnet-xxx`. valid vpc subnet ids can be queried by logging in to the [console](https://console.tencentcloud.com/vpc/subnet); or they can be obtained from the `SubnetId` field in the API response by calling the [DescribeSubnets](https://www.tencentcloud.com/document/product/215/15784) API . if SubnetId and VpcId are both input as `DEFAULT` when creating an instance, the DEFAULT vpc network will be forcibly used.
	SubnetId *string `json:"SubnetId,omitnil,omitempty" name:"SubnetId"`

	// Whether it is used as a public gateway. A public gateway can only be used normally when an instance has a public IP address and is in a VPC. Valid values:<li>true: It is used as a public gateway.</li><li>false: It is not used as a public gateway.</li>Default value: false.
	AsVpcGateway *bool `json:"AsVpcGateway,omitnil,omitempty" name:"AsVpcGateway"`

	// Array of VPC subnet IPs. You can use this parameter when creating instances or modifying VPC attributes of instances. Currently you can specify multiple IPs in one subnet only when creating multiple instances at the same time.
	PrivateIpAddresses []*string `json:"PrivateIpAddresses,omitnil,omitempty" name:"PrivateIpAddresses"`

	// Number of IPv6 addresses randomly generated for the ENI.
	// If IPv6AddressType is specified under InternetAccessible, this parameter must not be set to 0.
	Ipv6AddressCount *uint64 `json:"Ipv6AddressCount,omitnil,omitempty" name:"Ipv6AddressCount"`
}

type ZoneInfo struct {
	// Availability zone name, for example, ap-guangzhou-3.
	// 
	// The names of availability zones across the network are as follows:
	// <li> ap-chongqing-1 </li>
	// <li> ap-seoul-1 </li>
	// <li> ap-seoul-2 </li>
	// <li> ap-chengdu-1 </li>
	// <li> ap-chengdu-2 </li>
	// <li> ap-hongkong-1 (sold out)</li>
	// <li> ap-hongkong-2 </li>
	// <li> ap-hongkong-3 </li>
	// <li> ap-shenzhen-fsi-1 </li>
	// <li> ap-shenzhen-fsi-2 </li>
	// <li> ap-shenzhen-fsi-3 (sold out)</li>
	// <li> ap-guangzhou-1 (sold out)</li>
	// <li> ap-guangzhou-3 </li>
	// <li> ap-guangzhou-4 </li>
	// <li> ap-guangzhou-6 </li>
	// <li> ap-guangzhou-7 </li>
	// <li> ap-tokyo-1 </li>
	// <li> ap-tokyo-2 </li>
	// <li> ap-singapore-1 </li>
	// <li> ap-singapore-2 </li>
	// <li> ap-singapore-3 </li>
	// <li>ap-singapore-4 </li>
	// <li> ap-shanghai-fsi-1 </li>
	// <li> ap-shanghai-fsi-2 </li>
	// <li> ap-shanghai-fsi-3 </li>
	// <li> ap-bangkok-1 </li>
	// <li> ap-bangkok-2 </li>
	// <li> ap-shanghai-2 </li>
	// <li> ap-shanghai-3 </li>
	// <li> ap-shanghai-4 </li>
	// <li> ap-shanghai-5 </li>
	// <li> ap-shanghai-8 </li>
	// <li> ap-mumbai-1 </li>
	// <li> ap-mumbai-2 </li>
	// <li> ap-beijing-1 (sold out)</li>
	// <li> ap-beijing-3 </li>
	// <li> ap-beijing-4 </li>
	// <li> ap-beijing-5 </li>
	// <li> ap-beijing-6 </li>
	// <li> ap-beijing-7 </li>
	// <li> na-siliconvalley-1 </li>
	// <li> na-siliconvalley-2 </li>
	// <li> eu-frankfurt-1 </li>
	// <li> eu-frankfurt-2 </li>
	// <li> na-ashburn-1 </li>
	// <li> na-ashburn-2 </li>
	// <li> ap-nanjing-1 </li>
	// <li> ap-nanjing-2 </li>
	// <li> ap-nanjing-3 </li>
	// <li> sa-saopaulo-1</li>
	// <li> ap-jakarta-1 </li>
	// <li> ap-jakarta-2 </li>
	Zone *string `json:"Zone,omitnil,omitempty" name:"Zone"`

	// Availability zone description, such as Guangzhou Zone 3.
	ZoneName *string `json:"ZoneName,omitnil,omitempty" name:"ZoneName"`

	// Availability zone ID.
	ZoneId *string `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`

	// Availability zone status. Valid values: `AVAILABLE`: available; `UNAVAILABLE`: unavailable.
	ZoneState *string `json:"ZoneState,omitnil,omitempty" name:"ZoneState"`
}