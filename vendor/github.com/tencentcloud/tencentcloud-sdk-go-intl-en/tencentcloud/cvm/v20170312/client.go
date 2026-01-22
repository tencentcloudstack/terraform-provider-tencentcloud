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
    "context"
    "errors"
    "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
    tchttp "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common/http"
    "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common/profile"
)

const APIVersion = "2017-03-12"

type Client struct {
    common.Client
}

// Deprecated
func NewClientWithSecretId(secretId, secretKey, region string) (client *Client, err error) {
    cpf := profile.NewClientProfile()
    client = &Client{}
    client.Init(region).WithSecretId(secretId, secretKey).WithProfile(cpf)
    return
}

func NewClient(credential common.CredentialIface, region string, clientProfile *profile.ClientProfile) (client *Client, err error) {
    client = &Client{}
    client.Init(region).
        WithCredential(credential).
        WithProfile(clientProfile)
    return
}


func NewAllocateHostsRequest() (request *AllocateHostsRequest) {
    request = &AllocateHostsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "AllocateHosts")
    
    
    return
}

func NewAllocateHostsResponse() (response *AllocateHostsResponse) {
    response = &AllocateHostsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// AllocateHosts
// This API is used to create CDH instances with specified configuration.
//
// error code that may be returned:
//  INVALIDACCOUNT_INSUFFICIENTBALANCE = "InvalidAccount.InsufficientBalance"
//  INVALIDHOSTID_MALFORMED = "InvalidHostId.Malformed"
//  INVALIDPARAMETERVALUE_RANGE = "InvalidParameterValue.Range"
//  INVALIDPARAMETERVALUE_ZONENOTSUPPORTED = "InvalidParameterValue.ZoneNotSupported"
//  INVALIDPERIOD = "InvalidPeriod"
//  INVALIDPROJECTID_NOTFOUND = "InvalidProjectId.NotFound"
//  INVALIDREGION_NOTFOUND = "InvalidRegion.NotFound"
//  INVALIDZONE_MISMATCHREGION = "InvalidZone.MismatchRegion"
//  RESOURCEINSUFFICIENT_ZONESOLDOUTFORSPECIFIEDINSTANCE = "ResourceInsufficient.ZoneSoldOutForSpecifiedInstance"
func (c *Client) AllocateHosts(request *AllocateHostsRequest) (response *AllocateHostsResponse, err error) {
    return c.AllocateHostsWithContext(context.Background(), request)
}

// AllocateHosts
// This API is used to create CDH instances with specified configuration.
//
// error code that may be returned:
//  INVALIDACCOUNT_INSUFFICIENTBALANCE = "InvalidAccount.InsufficientBalance"
//  INVALIDHOSTID_MALFORMED = "InvalidHostId.Malformed"
//  INVALIDPARAMETERVALUE_RANGE = "InvalidParameterValue.Range"
//  INVALIDPARAMETERVALUE_ZONENOTSUPPORTED = "InvalidParameterValue.ZoneNotSupported"
//  INVALIDPERIOD = "InvalidPeriod"
//  INVALIDPROJECTID_NOTFOUND = "InvalidProjectId.NotFound"
//  INVALIDREGION_NOTFOUND = "InvalidRegion.NotFound"
//  INVALIDZONE_MISMATCHREGION = "InvalidZone.MismatchRegion"
//  RESOURCEINSUFFICIENT_ZONESOLDOUTFORSPECIFIEDINSTANCE = "ResourceInsufficient.ZoneSoldOutForSpecifiedInstance"
func (c *Client) AllocateHostsWithContext(ctx context.Context, request *AllocateHostsRequest) (response *AllocateHostsResponse, err error) {
    if request == nil {
        request = NewAllocateHostsRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "AllocateHosts")
    
    if c.GetCredential() == nil {
        return nil, errors.New("AllocateHosts require credential")
    }

    request.SetContext(ctx)
    
    response = NewAllocateHostsResponse()
    err = c.Send(request, response)
    return
}

func NewAssociateInstancesKeyPairsRequest() (request *AssociateInstancesKeyPairsRequest) {
    request = &AssociateInstancesKeyPairsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "AssociateInstancesKeyPairs")
    
    
    return
}

func NewAssociateInstancesKeyPairsResponse() (response *AssociateInstancesKeyPairsResponse) {
    response = &AssociateInstancesKeyPairsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// AssociateInstancesKeyPairs
// This API is used to associate key pairs with instances.
//
// 
//
// * If the public key of a key pair is written to the `SSH` configuration of the instance, users will be able to log in to the instance with the private key of the key pair.
//
// * If the instance is already associated with a key, the old key will become invalid.
//
// If you currently use a password to log in, you will no longer be able to do so after you associate the instance with a key.
//
// * Batch operations are supported. The maximum number of instances in each request is 100. If any instance in the request cannot be associated with a key, you will get an error code.
//
// error code that may be returned:
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDKEYPAIRID_MALFORMED = "InvalidKeyPairId.Malformed"
//  INVALIDKEYPAIRID_NOTFOUND = "InvalidKeyPairId.NotFound"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_KEYPAIRNOTSUPPORTED = "InvalidParameterValue.KeyPairNotSupported"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  MISSINGPARAMETER = "MissingParameter"
//  MUTEXOPERATION_TASKRUNNING = "MutexOperation.TaskRunning"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  UNAUTHORIZEDOPERATION_MFAEXPIRED = "UnauthorizedOperation.MFAExpired"
//  UNAUTHORIZEDOPERATION_MFANOTFOUND = "UnauthorizedOperation.MFANotFound"
//  UNSUPPORTEDOPERATION_INSTANCEOSWINDOWS = "UnsupportedOperation.InstanceOsWindows"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERSERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateEnterServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITRESCUEMODE = "UnsupportedOperation.InstanceStateExitRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEFREEZING = "UnsupportedOperation.InstanceStateFreezing"
//  UNSUPPORTEDOPERATION_INSTANCESTATEISOLATING = "UnsupportedOperation.InstanceStateIsolating"
//  UNSUPPORTEDOPERATION_INSTANCESTATEPENDING = "UnsupportedOperation.InstanceStatePending"
//  UNSUPPORTEDOPERATION_INSTANCESTATEREBOOTING = "UnsupportedOperation.InstanceStateRebooting"
//  UNSUPPORTEDOPERATION_INSTANCESTATERESCUEMODE = "UnsupportedOperation.InstanceStateRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATERUNNING = "UnsupportedOperation.InstanceStateRunning"
//  UNSUPPORTEDOPERATION_INSTANCESTATESHUTDOWN = "UnsupportedOperation.InstanceStateShutdown"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTARTING = "UnsupportedOperation.InstanceStateStarting"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPING = "UnsupportedOperation.InstanceStateStopping"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATING = "UnsupportedOperation.InstanceStateTerminating"
//  UNSUPPORTEDOPERATION_STOPPEDMODESTOPCHARGING = "UnsupportedOperation.StoppedModeStopCharging"
func (c *Client) AssociateInstancesKeyPairs(request *AssociateInstancesKeyPairsRequest) (response *AssociateInstancesKeyPairsResponse, err error) {
    return c.AssociateInstancesKeyPairsWithContext(context.Background(), request)
}

// AssociateInstancesKeyPairs
// This API is used to associate key pairs with instances.
//
// 
//
// * If the public key of a key pair is written to the `SSH` configuration of the instance, users will be able to log in to the instance with the private key of the key pair.
//
// * If the instance is already associated with a key, the old key will become invalid.
//
// If you currently use a password to log in, you will no longer be able to do so after you associate the instance with a key.
//
// * Batch operations are supported. The maximum number of instances in each request is 100. If any instance in the request cannot be associated with a key, you will get an error code.
//
// error code that may be returned:
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDKEYPAIRID_MALFORMED = "InvalidKeyPairId.Malformed"
//  INVALIDKEYPAIRID_NOTFOUND = "InvalidKeyPairId.NotFound"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_KEYPAIRNOTSUPPORTED = "InvalidParameterValue.KeyPairNotSupported"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  MISSINGPARAMETER = "MissingParameter"
//  MUTEXOPERATION_TASKRUNNING = "MutexOperation.TaskRunning"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  UNAUTHORIZEDOPERATION_MFAEXPIRED = "UnauthorizedOperation.MFAExpired"
//  UNAUTHORIZEDOPERATION_MFANOTFOUND = "UnauthorizedOperation.MFANotFound"
//  UNSUPPORTEDOPERATION_INSTANCEOSWINDOWS = "UnsupportedOperation.InstanceOsWindows"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERSERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateEnterServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITRESCUEMODE = "UnsupportedOperation.InstanceStateExitRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEFREEZING = "UnsupportedOperation.InstanceStateFreezing"
//  UNSUPPORTEDOPERATION_INSTANCESTATEISOLATING = "UnsupportedOperation.InstanceStateIsolating"
//  UNSUPPORTEDOPERATION_INSTANCESTATEPENDING = "UnsupportedOperation.InstanceStatePending"
//  UNSUPPORTEDOPERATION_INSTANCESTATEREBOOTING = "UnsupportedOperation.InstanceStateRebooting"
//  UNSUPPORTEDOPERATION_INSTANCESTATERESCUEMODE = "UnsupportedOperation.InstanceStateRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATERUNNING = "UnsupportedOperation.InstanceStateRunning"
//  UNSUPPORTEDOPERATION_INSTANCESTATESHUTDOWN = "UnsupportedOperation.InstanceStateShutdown"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTARTING = "UnsupportedOperation.InstanceStateStarting"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPING = "UnsupportedOperation.InstanceStateStopping"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATING = "UnsupportedOperation.InstanceStateTerminating"
//  UNSUPPORTEDOPERATION_STOPPEDMODESTOPCHARGING = "UnsupportedOperation.StoppedModeStopCharging"
func (c *Client) AssociateInstancesKeyPairsWithContext(ctx context.Context, request *AssociateInstancesKeyPairsRequest) (response *AssociateInstancesKeyPairsResponse, err error) {
    if request == nil {
        request = NewAssociateInstancesKeyPairsRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "AssociateInstancesKeyPairs")
    
    if c.GetCredential() == nil {
        return nil, errors.New("AssociateInstancesKeyPairs require credential")
    }

    request.SetContext(ctx)
    
    response = NewAssociateInstancesKeyPairsResponse()
    err = c.Send(request, response)
    return
}

func NewAssociateSecurityGroupsRequest() (request *AssociateSecurityGroupsRequest) {
    request = &AssociateSecurityGroupsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "AssociateSecurityGroups")
    
    
    return
}

func NewAssociateSecurityGroupsResponse() (response *AssociateSecurityGroupsResponse) {
    response = &AssociateSecurityGroupsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// AssociateSecurityGroups
// This API is used to associate security groups with specified instances.
//
// error code that may be returned:
//  FAILEDOPERATION_SECURITYGROUPACTIONFAILED = "FailedOperation.SecurityGroupActionFailed"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDSECURITYGROUPID_NOTFOUND = "InvalidSecurityGroupId.NotFound"
//  INVALIDSGID_MALFORMED = "InvalidSgId.Malformed"
//  LIMITEXCEEDED_ASSOCIATEUSGLIMITEXCEEDED = "LimitExceeded.AssociateUSGLimitExceeded"
//  LIMITEXCEEDED_CVMSVIFSPERSECGROUPLIMITEXCEEDED = "LimitExceeded.CvmsVifsPerSecGroupLimitExceeded"
//  LIMITEXCEEDED_SINGLEUSGQUOTA = "LimitExceeded.SingleUSGQuota"
//  MUTEXOPERATION_TASKRUNNING = "MutexOperation.TaskRunning"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  SECGROUPACTIONFAILURE = "SecGroupActionFailure"
//  UNSUPPORTEDOPERATION_INSTANCESTATEISOLATING = "UnsupportedOperation.InstanceStateIsolating"
//  UNSUPPORTEDOPERATION_INSTANCESTATEPENDING = "UnsupportedOperation.InstanceStatePending"
func (c *Client) AssociateSecurityGroups(request *AssociateSecurityGroupsRequest) (response *AssociateSecurityGroupsResponse, err error) {
    return c.AssociateSecurityGroupsWithContext(context.Background(), request)
}

// AssociateSecurityGroups
// This API is used to associate security groups with specified instances.
//
// error code that may be returned:
//  FAILEDOPERATION_SECURITYGROUPACTIONFAILED = "FailedOperation.SecurityGroupActionFailed"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDSECURITYGROUPID_NOTFOUND = "InvalidSecurityGroupId.NotFound"
//  INVALIDSGID_MALFORMED = "InvalidSgId.Malformed"
//  LIMITEXCEEDED_ASSOCIATEUSGLIMITEXCEEDED = "LimitExceeded.AssociateUSGLimitExceeded"
//  LIMITEXCEEDED_CVMSVIFSPERSECGROUPLIMITEXCEEDED = "LimitExceeded.CvmsVifsPerSecGroupLimitExceeded"
//  LIMITEXCEEDED_SINGLEUSGQUOTA = "LimitExceeded.SingleUSGQuota"
//  MUTEXOPERATION_TASKRUNNING = "MutexOperation.TaskRunning"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  SECGROUPACTIONFAILURE = "SecGroupActionFailure"
//  UNSUPPORTEDOPERATION_INSTANCESTATEISOLATING = "UnsupportedOperation.InstanceStateIsolating"
//  UNSUPPORTEDOPERATION_INSTANCESTATEPENDING = "UnsupportedOperation.InstanceStatePending"
func (c *Client) AssociateSecurityGroupsWithContext(ctx context.Context, request *AssociateSecurityGroupsRequest) (response *AssociateSecurityGroupsResponse, err error) {
    if request == nil {
        request = NewAssociateSecurityGroupsRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "AssociateSecurityGroups")
    
    if c.GetCredential() == nil {
        return nil, errors.New("AssociateSecurityGroups require credential")
    }

    request.SetContext(ctx)
    
    response = NewAssociateSecurityGroupsResponse()
    err = c.Send(request, response)
    return
}

func NewConfigureChcAssistVpcRequest() (request *ConfigureChcAssistVpcRequest) {
    request = &ConfigureChcAssistVpcRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "ConfigureChcAssistVpc")
    
    
    return
}

func NewConfigureChcAssistVpcResponse() (response *ConfigureChcAssistVpcResponse) {
    response = &ConfigureChcAssistVpcResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// ConfigureChcAssistVpc
// This API is used to configure the out-of-band network and deployment network of a CHC host.
//
// error code that may be returned:
//  ACCOUNTQUALIFICATIONRESTRICTIONS = "AccountQualificationRestrictions"
//  INVALIDHOST_NOTSUPPORTED = "InvalidHost.NotSupported"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_AMOUNTNOTEQUAL = "InvalidParameterValue.AmountNotEqual"
//  INVALIDPARAMETERVALUE_CHCHOSTSNOTFOUND = "InvalidParameterValue.ChcHostsNotFound"
//  INVALIDPARAMETERVALUE_IPADDRESSMALFORMED = "InvalidParameterValue.IPAddressMalformed"
//  INVALIDPARAMETERVALUE_INCORRECTFORMAT = "InvalidParameterValue.IncorrectFormat"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_MUSTDHCPENABLEDVPC = "InvalidParameterValue.MustDhcpEnabledVpc"
//  INVALIDPARAMETERVALUE_SUBNETIDMALFORMED = "InvalidParameterValue.SubnetIdMalformed"
//  LIMITEXCEEDED_VPCSUBNETNUM = "LimitExceeded.VpcSubnetNum"
//  VPCADDRNOTINSUBNET = "VpcAddrNotInSubNet"
//  VPCIPISUSED = "VpcIpIsUsed"
func (c *Client) ConfigureChcAssistVpc(request *ConfigureChcAssistVpcRequest) (response *ConfigureChcAssistVpcResponse, err error) {
    return c.ConfigureChcAssistVpcWithContext(context.Background(), request)
}

// ConfigureChcAssistVpc
// This API is used to configure the out-of-band network and deployment network of a CHC host.
//
// error code that may be returned:
//  ACCOUNTQUALIFICATIONRESTRICTIONS = "AccountQualificationRestrictions"
//  INVALIDHOST_NOTSUPPORTED = "InvalidHost.NotSupported"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_AMOUNTNOTEQUAL = "InvalidParameterValue.AmountNotEqual"
//  INVALIDPARAMETERVALUE_CHCHOSTSNOTFOUND = "InvalidParameterValue.ChcHostsNotFound"
//  INVALIDPARAMETERVALUE_IPADDRESSMALFORMED = "InvalidParameterValue.IPAddressMalformed"
//  INVALIDPARAMETERVALUE_INCORRECTFORMAT = "InvalidParameterValue.IncorrectFormat"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_MUSTDHCPENABLEDVPC = "InvalidParameterValue.MustDhcpEnabledVpc"
//  INVALIDPARAMETERVALUE_SUBNETIDMALFORMED = "InvalidParameterValue.SubnetIdMalformed"
//  LIMITEXCEEDED_VPCSUBNETNUM = "LimitExceeded.VpcSubnetNum"
//  VPCADDRNOTINSUBNET = "VpcAddrNotInSubNet"
//  VPCIPISUSED = "VpcIpIsUsed"
func (c *Client) ConfigureChcAssistVpcWithContext(ctx context.Context, request *ConfigureChcAssistVpcRequest) (response *ConfigureChcAssistVpcResponse, err error) {
    if request == nil {
        request = NewConfigureChcAssistVpcRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "ConfigureChcAssistVpc")
    
    if c.GetCredential() == nil {
        return nil, errors.New("ConfigureChcAssistVpc require credential")
    }

    request.SetContext(ctx)
    
    response = NewConfigureChcAssistVpcResponse()
    err = c.Send(request, response)
    return
}

func NewConfigureChcDeployVpcRequest() (request *ConfigureChcDeployVpcRequest) {
    request = &ConfigureChcDeployVpcRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "ConfigureChcDeployVpc")
    
    
    return
}

func NewConfigureChcDeployVpcResponse() (response *ConfigureChcDeployVpcResponse) {
    response = &ConfigureChcDeployVpcResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// ConfigureChcDeployVpc
// This API is used to configure the deployment network of a CHC host.
//
// error code that may be returned:
//  ACCOUNTQUALIFICATIONRESTRICTIONS = "AccountQualificationRestrictions"
//  INVALIDHOST_NOTSUPPORTED = "InvalidHost.NotSupported"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_AMOUNTNOTEQUAL = "InvalidParameterValue.AmountNotEqual"
//  INVALIDPARAMETERVALUE_CHCHOSTSNOTFOUND = "InvalidParameterValue.ChcHostsNotFound"
//  INVALIDPARAMETERVALUE_DEPLOYVPCALREADYEXISTS = "InvalidParameterValue.DeployVpcAlreadyExists"
//  INVALIDPARAMETERVALUE_INCORRECTFORMAT = "InvalidParameterValue.IncorrectFormat"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_MUSTDHCPENABLEDVPC = "InvalidParameterValue.MustDhcpEnabledVpc"
//  VPCADDRNOTINSUBNET = "VpcAddrNotInSubNet"
//  VPCIPISUSED = "VpcIpIsUsed"
func (c *Client) ConfigureChcDeployVpc(request *ConfigureChcDeployVpcRequest) (response *ConfigureChcDeployVpcResponse, err error) {
    return c.ConfigureChcDeployVpcWithContext(context.Background(), request)
}

// ConfigureChcDeployVpc
// This API is used to configure the deployment network of a CHC host.
//
// error code that may be returned:
//  ACCOUNTQUALIFICATIONRESTRICTIONS = "AccountQualificationRestrictions"
//  INVALIDHOST_NOTSUPPORTED = "InvalidHost.NotSupported"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_AMOUNTNOTEQUAL = "InvalidParameterValue.AmountNotEqual"
//  INVALIDPARAMETERVALUE_CHCHOSTSNOTFOUND = "InvalidParameterValue.ChcHostsNotFound"
//  INVALIDPARAMETERVALUE_DEPLOYVPCALREADYEXISTS = "InvalidParameterValue.DeployVpcAlreadyExists"
//  INVALIDPARAMETERVALUE_INCORRECTFORMAT = "InvalidParameterValue.IncorrectFormat"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_MUSTDHCPENABLEDVPC = "InvalidParameterValue.MustDhcpEnabledVpc"
//  VPCADDRNOTINSUBNET = "VpcAddrNotInSubNet"
//  VPCIPISUSED = "VpcIpIsUsed"
func (c *Client) ConfigureChcDeployVpcWithContext(ctx context.Context, request *ConfigureChcDeployVpcRequest) (response *ConfigureChcDeployVpcResponse, err error) {
    if request == nil {
        request = NewConfigureChcDeployVpcRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "ConfigureChcDeployVpc")
    
    if c.GetCredential() == nil {
        return nil, errors.New("ConfigureChcDeployVpc require credential")
    }

    request.SetContext(ctx)
    
    response = NewConfigureChcDeployVpcResponse()
    err = c.Send(request, response)
    return
}

func NewConvertOperatingSystemsRequest() (request *ConvertOperatingSystemsRequest) {
    request = &ConvertOperatingSystemsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "ConvertOperatingSystems")
    
    
    return
}

func NewConvertOperatingSystemsResponse() (response *ConvertOperatingSystemsResponse) {
    response = &ConvertOperatingSystemsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// ConvertOperatingSystems
// This API is used to switch the operating system of an instance with CentOS 7 or CentOS 8 as the source operating system.
//
// error code that may be returned:
//  FAILEDOPERATION_GETINSTANCETATAGENTSTATUSFAILED = "FailedOperation.GetInstanceTATAgentStatusFailed"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETER_AUTOSNAPSHOTNOTSUPPORTED = "InvalidParameter.AutoSnapshotNotSupported"
//  INVALIDPARAMETER_INVALIDTARGETOSTYPE = "InvalidParameter.InvalidTargetOSType"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  UNSUPPORTEDOPERATION_INSTANCEOSCONVERTOSNOTSUPPORT = "UnsupportedOperation.InstanceOsConvertOsNotSupport"
//  UNSUPPORTEDOPERATION_INSTANCEOSWINDOWS = "UnsupportedOperation.InstanceOsWindows"
//  UNSUPPORTEDOPERATION_INSTANCESTATENOTRUNNING = "UnsupportedOperation.InstanceStateNotRunning"
//  UNSUPPORTEDOPERATION_MARKETIMAGECONVERTOSUNSUPPORTED = "UnsupportedOperation.MarketImageConvertOSUnsupported"
//  UNSUPPORTEDOPERATION_SPECIALINSTANCETYPE = "UnsupportedOperation.SpecialInstanceType"
//  UNSUPPORTEDOPERATION_TATAGENTNOTONLINE = "UnsupportedOperation.TatAgentNotOnline"
//  UNSUPPORTEDOPERATION_USERCONVERTOSNOTSUPPORT = "UnsupportedOperation.UserConvertOsNotSupport"
func (c *Client) ConvertOperatingSystems(request *ConvertOperatingSystemsRequest) (response *ConvertOperatingSystemsResponse, err error) {
    return c.ConvertOperatingSystemsWithContext(context.Background(), request)
}

// ConvertOperatingSystems
// This API is used to switch the operating system of an instance with CentOS 7 or CentOS 8 as the source operating system.
//
// error code that may be returned:
//  FAILEDOPERATION_GETINSTANCETATAGENTSTATUSFAILED = "FailedOperation.GetInstanceTATAgentStatusFailed"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETER_AUTOSNAPSHOTNOTSUPPORTED = "InvalidParameter.AutoSnapshotNotSupported"
//  INVALIDPARAMETER_INVALIDTARGETOSTYPE = "InvalidParameter.InvalidTargetOSType"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  UNSUPPORTEDOPERATION_INSTANCEOSCONVERTOSNOTSUPPORT = "UnsupportedOperation.InstanceOsConvertOsNotSupport"
//  UNSUPPORTEDOPERATION_INSTANCEOSWINDOWS = "UnsupportedOperation.InstanceOsWindows"
//  UNSUPPORTEDOPERATION_INSTANCESTATENOTRUNNING = "UnsupportedOperation.InstanceStateNotRunning"
//  UNSUPPORTEDOPERATION_MARKETIMAGECONVERTOSUNSUPPORTED = "UnsupportedOperation.MarketImageConvertOSUnsupported"
//  UNSUPPORTEDOPERATION_SPECIALINSTANCETYPE = "UnsupportedOperation.SpecialInstanceType"
//  UNSUPPORTEDOPERATION_TATAGENTNOTONLINE = "UnsupportedOperation.TatAgentNotOnline"
//  UNSUPPORTEDOPERATION_USERCONVERTOSNOTSUPPORT = "UnsupportedOperation.UserConvertOsNotSupport"
func (c *Client) ConvertOperatingSystemsWithContext(ctx context.Context, request *ConvertOperatingSystemsRequest) (response *ConvertOperatingSystemsResponse, err error) {
    if request == nil {
        request = NewConvertOperatingSystemsRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "ConvertOperatingSystems")
    
    if c.GetCredential() == nil {
        return nil, errors.New("ConvertOperatingSystems require credential")
    }

    request.SetContext(ctx)
    
    response = NewConvertOperatingSystemsResponse()
    err = c.Send(request, response)
    return
}

func NewCreateDisasterRecoverGroupRequest() (request *CreateDisasterRecoverGroupRequest) {
    request = &CreateDisasterRecoverGroupRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "CreateDisasterRecoverGroup")
    
    
    return
}

func NewCreateDisasterRecoverGroupResponse() (response *CreateDisasterRecoverGroupResponse) {
    response = &CreateDisasterRecoverGroupResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// CreateDisasterRecoverGroup
// This API is used to create a [spread placement group](https://intl.cloud.tencent.com/document/product/213/15486?from_cn_redirect=1). After you create one, you can specify it for an instance when you [create the instance](https://intl.cloud.tencent.com/document/api/213/15730?from_cn_redirect=1),
//
// error code that may be returned:
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_RANGE = "InvalidParameterValue.Range"
//  INVALIDPARAMETERVALUE_TOOLONG = "InvalidParameterValue.TooLong"
//  RESOURCEINSUFFICIENT_INSUFFICIENTGROUPQUOTA = "ResourceInsufficient.InsufficientGroupQuota"
func (c *Client) CreateDisasterRecoverGroup(request *CreateDisasterRecoverGroupRequest) (response *CreateDisasterRecoverGroupResponse, err error) {
    return c.CreateDisasterRecoverGroupWithContext(context.Background(), request)
}

// CreateDisasterRecoverGroup
// This API is used to create a [spread placement group](https://intl.cloud.tencent.com/document/product/213/15486?from_cn_redirect=1). After you create one, you can specify it for an instance when you [create the instance](https://intl.cloud.tencent.com/document/api/213/15730?from_cn_redirect=1),
//
// error code that may be returned:
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_RANGE = "InvalidParameterValue.Range"
//  INVALIDPARAMETERVALUE_TOOLONG = "InvalidParameterValue.TooLong"
//  RESOURCEINSUFFICIENT_INSUFFICIENTGROUPQUOTA = "ResourceInsufficient.InsufficientGroupQuota"
func (c *Client) CreateDisasterRecoverGroupWithContext(ctx context.Context, request *CreateDisasterRecoverGroupRequest) (response *CreateDisasterRecoverGroupResponse, err error) {
    if request == nil {
        request = NewCreateDisasterRecoverGroupRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "CreateDisasterRecoverGroup")
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateDisasterRecoverGroup require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateDisasterRecoverGroupResponse()
    err = c.Send(request, response)
    return
}

func NewCreateImageRequest() (request *CreateImageRequest) {
    request = &CreateImageRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "CreateImage")
    
    
    return
}

func NewCreateImageResponse() (response *CreateImageResponse) {
    response = &CreateImageResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// CreateImage
// This API is used to create a new image with the system disk of an instance. The image can be used to create new instances.
//
// error code that may be returned:
//  IMAGEQUOTALIMITEXCEEDED = "ImageQuotaLimitExceeded"
//  INVALIDACCOUNT_INSUFFICIENTBALANCE = "InvalidAccount.InsufficientBalance"
//  INVALIDIMAGENAME_DUPLICATE = "InvalidImageName.Duplicate"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETER_DATADISKIDCONTAINSROOTDISK = "InvalidParameter.DataDiskIdContainsRootDisk"
//  INVALIDPARAMETER_DATADISKNOTBELONGSPECIFIEDINSTANCE = "InvalidParameter.DataDiskNotBelongSpecifiedInstance"
//  INVALIDPARAMETER_DUPLICATESYSTEMSNAPSHOTS = "InvalidParameter.DuplicateSystemSnapshots"
//  INVALIDPARAMETER_INSTANCEIMAGENOTSUPPORT = "InvalidParameter.InstanceImageNotSupport"
//  INVALIDPARAMETER_INVALIDDEPENDENCE = "InvalidParameter.InvalidDependence"
//  INVALIDPARAMETER_LOCALDATADISKNOTSUPPORT = "InvalidParameter.LocalDataDiskNotSupport"
//  INVALIDPARAMETER_SNAPSHOTNOTFOUND = "InvalidParameter.SnapshotNotFound"
//  INVALIDPARAMETER_SPECIFYONEPARAMETER = "InvalidParameter.SpecifyOneParameter"
//  INVALIDPARAMETER_SWAPDISKNOTSUPPORT = "InvalidParameter.SwapDiskNotSupport"
//  INVALIDPARAMETER_SYSTEMSNAPSHOTNOTFOUND = "InvalidParameter.SystemSnapshotNotFound"
//  INVALIDPARAMETER_VALUETOOLARGE = "InvalidParameter.ValueTooLarge"
//  INVALIDPARAMETERCONFLICT = "InvalidParameterConflict"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEFAMILY = "InvalidParameterValue.InvalidImageFamily"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_PREHEATNOTSUPPORTEDINSTANCETYPE = "InvalidParameterValue.PreheatNotSupportedInstanceType"
//  INVALIDPARAMETERVALUE_PREHEATNOTSUPPORTEDZONE = "InvalidParameterValue.PreheatNotSupportedZone"
//  INVALIDPARAMETERVALUE_PREHEATUNAVAILABLEZONES = "InvalidParameterValue.PreheatUnavailableZones"
//  INVALIDPARAMETERVALUE_TAGKEYNOTFOUND = "InvalidParameterValue.TagKeyNotFound"
//  INVALIDPARAMETERVALUE_TAGQUOTALIMITEXCEEDED = "InvalidParameterValue.TagQuotaLimitExceeded"
//  INVALIDPARAMETERVALUE_TOOLARGE = "InvalidParameterValue.TooLarge"
//  INVALIDZONE_MISMATCHREGION = "InvalidZone.MismatchRegion"
//  LIMITEXCEEDED_PREHEATIMAGESNAPSHOTOUTOFQUOTA = "LimitExceeded.PreheatImageSnapshotOutOfQuota"
//  MUTEXOPERATION_TASKRUNNING = "MutexOperation.TaskRunning"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  RESOURCEINUSE = "ResourceInUse"
//  RESOURCEINUSE_DISKROLLBACKING = "ResourceInUse.DiskRollbacking"
//  RESOURCEINSUFFICIENT_CLOUDDISKUNAVAILABLE = "ResourceInsufficient.CloudDiskUnavailable"
//  RESOURCEUNAVAILABLE_SNAPSHOTCREATING = "ResourceUnavailable.SnapshotCreating"
//  UNSUPPORTEDOPERATION_ENCRYPTEDIMAGESNOTSUPPORTED = "UnsupportedOperation.EncryptedImagesNotSupported"
//  UNSUPPORTEDOPERATION_INSTANCEREINSTALLFAILED = "UnsupportedOperation.InstanceReinstallFailed"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERRESCUEMODE = "UnsupportedOperation.InstanceStateEnterRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERSERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateEnterServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITRESCUEMODE = "UnsupportedOperation.InstanceStateExitRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEFREEZING = "UnsupportedOperation.InstanceStateFreezing"
//  UNSUPPORTEDOPERATION_INSTANCESTATEREBOOTING = "UnsupportedOperation.InstanceStateRebooting"
//  UNSUPPORTEDOPERATION_INSTANCESTATERESCUEMODE = "UnsupportedOperation.InstanceStateRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATESERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTARTING = "UnsupportedOperation.InstanceStateStarting"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPING = "UnsupportedOperation.InstanceStateStopping"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATING = "UnsupportedOperation.InstanceStateTerminating"
//  UNSUPPORTEDOPERATION_INVALIDDISKFASTROLLBACK = "UnsupportedOperation.InvalidDiskFastRollback"
//  UNSUPPORTEDOPERATION_NOTSUPPORTINSTANCEIMAGE = "UnsupportedOperation.NotSupportInstanceImage"
//  UNSUPPORTEDOPERATION_PREHEATIMAGE = "UnsupportedOperation.PreheatImage"
//  UNSUPPORTEDOPERATION_SPECIALINSTANCETYPE = "UnsupportedOperation.SpecialInstanceType"
//  UNSUPPORTEDOPERATION_STOPPEDMODESTOPCHARGING = "UnsupportedOperation.StoppedModeStopCharging"
func (c *Client) CreateImage(request *CreateImageRequest) (response *CreateImageResponse, err error) {
    return c.CreateImageWithContext(context.Background(), request)
}

// CreateImage
// This API is used to create a new image with the system disk of an instance. The image can be used to create new instances.
//
// error code that may be returned:
//  IMAGEQUOTALIMITEXCEEDED = "ImageQuotaLimitExceeded"
//  INVALIDACCOUNT_INSUFFICIENTBALANCE = "InvalidAccount.InsufficientBalance"
//  INVALIDIMAGENAME_DUPLICATE = "InvalidImageName.Duplicate"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETER_DATADISKIDCONTAINSROOTDISK = "InvalidParameter.DataDiskIdContainsRootDisk"
//  INVALIDPARAMETER_DATADISKNOTBELONGSPECIFIEDINSTANCE = "InvalidParameter.DataDiskNotBelongSpecifiedInstance"
//  INVALIDPARAMETER_DUPLICATESYSTEMSNAPSHOTS = "InvalidParameter.DuplicateSystemSnapshots"
//  INVALIDPARAMETER_INSTANCEIMAGENOTSUPPORT = "InvalidParameter.InstanceImageNotSupport"
//  INVALIDPARAMETER_INVALIDDEPENDENCE = "InvalidParameter.InvalidDependence"
//  INVALIDPARAMETER_LOCALDATADISKNOTSUPPORT = "InvalidParameter.LocalDataDiskNotSupport"
//  INVALIDPARAMETER_SNAPSHOTNOTFOUND = "InvalidParameter.SnapshotNotFound"
//  INVALIDPARAMETER_SPECIFYONEPARAMETER = "InvalidParameter.SpecifyOneParameter"
//  INVALIDPARAMETER_SWAPDISKNOTSUPPORT = "InvalidParameter.SwapDiskNotSupport"
//  INVALIDPARAMETER_SYSTEMSNAPSHOTNOTFOUND = "InvalidParameter.SystemSnapshotNotFound"
//  INVALIDPARAMETER_VALUETOOLARGE = "InvalidParameter.ValueTooLarge"
//  INVALIDPARAMETERCONFLICT = "InvalidParameterConflict"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEFAMILY = "InvalidParameterValue.InvalidImageFamily"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_PREHEATNOTSUPPORTEDINSTANCETYPE = "InvalidParameterValue.PreheatNotSupportedInstanceType"
//  INVALIDPARAMETERVALUE_PREHEATNOTSUPPORTEDZONE = "InvalidParameterValue.PreheatNotSupportedZone"
//  INVALIDPARAMETERVALUE_PREHEATUNAVAILABLEZONES = "InvalidParameterValue.PreheatUnavailableZones"
//  INVALIDPARAMETERVALUE_TAGKEYNOTFOUND = "InvalidParameterValue.TagKeyNotFound"
//  INVALIDPARAMETERVALUE_TAGQUOTALIMITEXCEEDED = "InvalidParameterValue.TagQuotaLimitExceeded"
//  INVALIDPARAMETERVALUE_TOOLARGE = "InvalidParameterValue.TooLarge"
//  INVALIDZONE_MISMATCHREGION = "InvalidZone.MismatchRegion"
//  LIMITEXCEEDED_PREHEATIMAGESNAPSHOTOUTOFQUOTA = "LimitExceeded.PreheatImageSnapshotOutOfQuota"
//  MUTEXOPERATION_TASKRUNNING = "MutexOperation.TaskRunning"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  RESOURCEINUSE = "ResourceInUse"
//  RESOURCEINUSE_DISKROLLBACKING = "ResourceInUse.DiskRollbacking"
//  RESOURCEINSUFFICIENT_CLOUDDISKUNAVAILABLE = "ResourceInsufficient.CloudDiskUnavailable"
//  RESOURCEUNAVAILABLE_SNAPSHOTCREATING = "ResourceUnavailable.SnapshotCreating"
//  UNSUPPORTEDOPERATION_ENCRYPTEDIMAGESNOTSUPPORTED = "UnsupportedOperation.EncryptedImagesNotSupported"
//  UNSUPPORTEDOPERATION_INSTANCEREINSTALLFAILED = "UnsupportedOperation.InstanceReinstallFailed"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERRESCUEMODE = "UnsupportedOperation.InstanceStateEnterRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERSERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateEnterServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITRESCUEMODE = "UnsupportedOperation.InstanceStateExitRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEFREEZING = "UnsupportedOperation.InstanceStateFreezing"
//  UNSUPPORTEDOPERATION_INSTANCESTATEREBOOTING = "UnsupportedOperation.InstanceStateRebooting"
//  UNSUPPORTEDOPERATION_INSTANCESTATERESCUEMODE = "UnsupportedOperation.InstanceStateRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATESERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTARTING = "UnsupportedOperation.InstanceStateStarting"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPING = "UnsupportedOperation.InstanceStateStopping"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATING = "UnsupportedOperation.InstanceStateTerminating"
//  UNSUPPORTEDOPERATION_INVALIDDISKFASTROLLBACK = "UnsupportedOperation.InvalidDiskFastRollback"
//  UNSUPPORTEDOPERATION_NOTSUPPORTINSTANCEIMAGE = "UnsupportedOperation.NotSupportInstanceImage"
//  UNSUPPORTEDOPERATION_PREHEATIMAGE = "UnsupportedOperation.PreheatImage"
//  UNSUPPORTEDOPERATION_SPECIALINSTANCETYPE = "UnsupportedOperation.SpecialInstanceType"
//  UNSUPPORTEDOPERATION_STOPPEDMODESTOPCHARGING = "UnsupportedOperation.StoppedModeStopCharging"
func (c *Client) CreateImageWithContext(ctx context.Context, request *CreateImageRequest) (response *CreateImageResponse, err error) {
    if request == nil {
        request = NewCreateImageRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "CreateImage")
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateImage require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateImageResponse()
    err = c.Send(request, response)
    return
}

func NewCreateKeyPairRequest() (request *CreateKeyPairRequest) {
    request = &CreateKeyPairRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "CreateKeyPair")
    
    
    return
}

func NewCreateKeyPairResponse() (response *CreateKeyPairResponse) {
    response = &CreateKeyPairResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// CreateKeyPair
// This API is used to create an `OpenSSH RSA` key pair, which you can use to log in to a `Linux` instance.
//
// 
//
// * You only need to specify a name, and the system will automatically create a key pair and return its `ID` and the public and private keys.
//
// * The name of the key pair must be unique.
//
// * You can save the private key to a file and use it as an authentication method for `SSH`.
//
// * Tencent Cloud does not save users' private keys. Be sure to save it yourself.
//
// error code that may be returned:
//  FAILEDOPERATION_TAGKEYRESERVED = "FailedOperation.TagKeyReserved"
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDKEYPAIR_LIMITEXCEEDED = "InvalidKeyPair.LimitExceeded"
//  INVALIDKEYPAIRNAME_DUPLICATE = "InvalidKeyPairName.Duplicate"
//  INVALIDKEYPAIRNAMEEMPTY = "InvalidKeyPairNameEmpty"
//  INVALIDKEYPAIRNAMEINCLUDEILLEGALCHAR = "InvalidKeyPairNameIncludeIllegalChar"
//  INVALIDKEYPAIRNAMETOOLONG = "InvalidKeyPairNameTooLong"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPROJECTID_NOTFOUND = "InvalidProjectId.NotFound"
//  LIMITEXCEEDED_TAGRESOURCEQUOTA = "LimitExceeded.TagResourceQuota"
//  MISSINGPARAMETER = "MissingParameter"
func (c *Client) CreateKeyPair(request *CreateKeyPairRequest) (response *CreateKeyPairResponse, err error) {
    return c.CreateKeyPairWithContext(context.Background(), request)
}

// CreateKeyPair
// This API is used to create an `OpenSSH RSA` key pair, which you can use to log in to a `Linux` instance.
//
// 
//
// * You only need to specify a name, and the system will automatically create a key pair and return its `ID` and the public and private keys.
//
// * The name of the key pair must be unique.
//
// * You can save the private key to a file and use it as an authentication method for `SSH`.
//
// * Tencent Cloud does not save users' private keys. Be sure to save it yourself.
//
// error code that may be returned:
//  FAILEDOPERATION_TAGKEYRESERVED = "FailedOperation.TagKeyReserved"
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDKEYPAIR_LIMITEXCEEDED = "InvalidKeyPair.LimitExceeded"
//  INVALIDKEYPAIRNAME_DUPLICATE = "InvalidKeyPairName.Duplicate"
//  INVALIDKEYPAIRNAMEEMPTY = "InvalidKeyPairNameEmpty"
//  INVALIDKEYPAIRNAMEINCLUDEILLEGALCHAR = "InvalidKeyPairNameIncludeIllegalChar"
//  INVALIDKEYPAIRNAMETOOLONG = "InvalidKeyPairNameTooLong"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPROJECTID_NOTFOUND = "InvalidProjectId.NotFound"
//  LIMITEXCEEDED_TAGRESOURCEQUOTA = "LimitExceeded.TagResourceQuota"
//  MISSINGPARAMETER = "MissingParameter"
func (c *Client) CreateKeyPairWithContext(ctx context.Context, request *CreateKeyPairRequest) (response *CreateKeyPairResponse, err error) {
    if request == nil {
        request = NewCreateKeyPairRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "CreateKeyPair")
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateKeyPair require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateKeyPairResponse()
    err = c.Send(request, response)
    return
}

func NewCreateLaunchTemplateRequest() (request *CreateLaunchTemplateRequest) {
    request = &CreateLaunchTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "CreateLaunchTemplate")
    
    
    return
}

func NewCreateLaunchTemplateResponse() (response *CreateLaunchTemplateResponse) {
    response = &CreateLaunchTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// CreateLaunchTemplate
// This interface (CreateLaunchTemplate) is used for instance launch template creation.
//
// 
//
// An instance launch template is a configuration data and can be used to create instances. Its content includes configurations required to create instances, such as instance type, types and sizes of data disk and system disk, and security group and other information.
//
// 
//
// This API is used to create an instance launch template. After the initial creation of the instance template, its template version is the default version 1. A new version can be created using CreateLaunchTemplateVersion (https://intl.cloud.tencent.com/document/product/213/66326?from_cn_redirect=1), and the version number will increment. By default, when specifying an instance launch template in RunInstances (https://intl.cloud.tencent.com/document/product/213/15730?from_cn_redirect=1), if the template version number is not specified, the default version will be used.
//
// error code that may be returned:
//  ACCOUNTQUALIFICATIONRESTRICTIONS = "AccountQualificationRestrictions"
//  AUTHFAILURE_CAMROLENAMEAUTHENTICATEFAILED = "AuthFailure.CamRoleNameAuthenticateFailed"
//  FAILEDOPERATION_DISASTERRECOVERGROUPNOTFOUND = "FailedOperation.DisasterRecoverGroupNotFound"
//  FAILEDOPERATION_INQUIRYPRICEFAILED = "FailedOperation.InquiryPriceFailed"
//  FAILEDOPERATION_NOAVAILABLEIPADDRESSCOUNTINSUBNET = "FailedOperation.NoAvailableIpAddressCountInSubnet"
//  FAILEDOPERATION_SECURITYGROUPACTIONFAILED = "FailedOperation.SecurityGroupActionFailed"
//  FAILEDOPERATION_SNAPSHOTSIZELARGERTHANDATASIZE = "FailedOperation.SnapshotSizeLargerThanDataSize"
//  FAILEDOPERATION_TAGKEYRESERVED = "FailedOperation.TagKeyReserved"
//  INSTANCESQUOTALIMITEXCEEDED = "InstancesQuotaLimitExceeded"
//  INTERNALERROR_TRADEUNKNOWNERROR = "InternalError.TradeUnknownError"
//  INVALIDACCOUNT_INSUFFICIENTBALANCE = "InvalidAccount.InsufficientBalance"
//  INVALIDCLIENTTOKEN_TOOLONG = "InvalidClientToken.TooLong"
//  INVALIDHOSTID_MALFORMED = "InvalidHostId.Malformed"
//  INVALIDHOSTID_NOTFOUND = "InvalidHostId.NotFound"
//  INVALIDIMAGEID_MALFORMED = "InvalidImageId.Malformed"
//  INVALIDIMAGEID_NOTFOUND = "InvalidImageId.NotFound"
//  INVALIDINSTANCENAME_TOOLONG = "InvalidInstanceName.TooLong"
//  INVALIDINSTANCETYPE_MALFORMED = "InvalidInstanceType.Malformed"
//  INVALIDPARAMETER_INSTANCEIMAGENOTSUPPORT = "InvalidParameter.InstanceImageNotSupport"
//  INVALIDPARAMETER_INTERNETACCESSIBLENOTSUPPORTED = "InvalidParameter.InternetAccessibleNotSupported"
//  INVALIDPARAMETER_INVALIDIPFORMAT = "InvalidParameter.InvalidIpFormat"
//  INVALIDPARAMETER_LACKCORECOUNTORTHREADPERCORE = "InvalidParameter.LackCoreCountOrThreadPerCore"
//  INVALIDPARAMETER_PASSWORDNOTSUPPORTED = "InvalidParameter.PasswordNotSupported"
//  INVALIDPARAMETER_SNAPSHOTNOTFOUND = "InvalidParameter.SnapshotNotFound"
//  INVALIDPARAMETERCOMBINATION = "InvalidParameterCombination"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_CLOUDSSDDATADISKSIZETOOSMALL = "InvalidParameterValue.CloudSsdDataDiskSizeTooSmall"
//  INVALIDPARAMETERVALUE_CORECOUNTVALUE = "InvalidParameterValue.CoreCountValue"
//  INVALIDPARAMETERVALUE_ILLEGALHOSTNAME = "InvalidParameterValue.IllegalHostName"
//  INVALIDPARAMETERVALUE_INSTANCETYPENOTSUPPORTHPCCLUSTER = "InvalidParameterValue.InstanceTypeNotSupportHpcCluster"
//  INVALIDPARAMETERVALUE_INSTANCETYPEREQUIREDHPCCLUSTER = "InvalidParameterValue.InstanceTypeRequiredHpcCluster"
//  INVALIDPARAMETERVALUE_INSUFFICIENTOFFERING = "InvalidParameterValue.InsufficientOffering"
//  INVALIDPARAMETERVALUE_INSUFFICIENTPRICE = "InvalidParameterValue.InsufficientPrice"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEFORGIVENINSTANCETYPE = "InvalidParameterValue.InvalidImageForGivenInstanceType"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEID = "InvalidParameterValue.InvalidImageId"
//  INVALIDPARAMETERVALUE_INVALIDIMAGESTATE = "InvalidParameterValue.InvalidImageState"
//  INVALIDPARAMETERVALUE_INVALIDIPFORMAT = "InvalidParameterValue.InvalidIpFormat"
//  INVALIDPARAMETERVALUE_INVALIDLAUNCHTEMPLATEDESCRIPTION = "InvalidParameterValue.InvalidLaunchTemplateDescription"
//  INVALIDPARAMETERVALUE_INVALIDLAUNCHTEMPLATENAME = "InvalidParameterValue.InvalidLaunchTemplateName"
//  INVALIDPARAMETERVALUE_INVALIDLAUNCHTEMPLATEVERSIONDESCRIPTION = "InvalidParameterValue.InvalidLaunchTemplateVersionDescription"
//  INVALIDPARAMETERVALUE_INVALIDPASSWORD = "InvalidParameterValue.InvalidPassword"
//  INVALIDPARAMETERVALUE_INVALIDTIMEFORMAT = "InvalidParameterValue.InvalidTimeFormat"
//  INVALIDPARAMETERVALUE_INVALIDUSERDATAFORMAT = "InvalidParameterValue.InvalidUserDataFormat"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_MUSTDHCPENABLEDVPC = "InvalidParameterValue.MustDhcpEnabledVpc"
//  INVALIDPARAMETERVALUE_RANGE = "InvalidParameterValue.Range"
//  INVALIDPARAMETERVALUE_SNAPSHOTIDMALFORMED = "InvalidParameterValue.SnapshotIdMalformed"
//  INVALIDPARAMETERVALUE_SUBNETNOTEXIST = "InvalidParameterValue.SubnetNotExist"
//  INVALIDPARAMETERVALUE_THREADPERCOREVALUE = "InvalidParameterValue.ThreadPerCoreValue"
//  INVALIDPARAMETERVALUE_VPCIDMALFORMED = "InvalidParameterValue.VpcIdMalformed"
//  INVALIDPARAMETERVALUE_VPCIDNOTEXIST = "InvalidParameterValue.VpcIdNotExist"
//  INVALIDPARAMETERVALUE_VPCIDZONEIDNOTMATCH = "InvalidParameterValue.VpcIdZoneIdNotMatch"
//  INVALIDPARAMETERVALUE_ZONENOTSUPPORTED = "InvalidParameterValue.ZoneNotSupported"
//  INVALIDPASSWORD = "InvalidPassword"
//  INVALIDPERIOD = "InvalidPeriod"
//  INVALIDPERMISSION = "InvalidPermission"
//  INVALIDPROJECTID_NOTFOUND = "InvalidProjectId.NotFound"
//  INVALIDSECURITYGROUPID_NOTFOUND = "InvalidSecurityGroupId.NotFound"
//  INVALIDZONE_MISMATCHREGION = "InvalidZone.MismatchRegion"
//  LIMITEXCEEDED_INSTANCEQUOTA = "LimitExceeded.InstanceQuota"
//  LIMITEXCEEDED_LAUNCHTEMPLATEQUOTA = "LimitExceeded.LaunchTemplateQuota"
//  LIMITEXCEEDED_SINGLEUSGQUOTA = "LimitExceeded.SingleUSGQuota"
//  LIMITEXCEEDED_SPOTQUOTA = "LimitExceeded.SpotQuota"
//  LIMITEXCEEDED_USERSPOTQUOTA = "LimitExceeded.UserSpotQuota"
//  LIMITEXCEEDED_VPCSUBNETNUM = "LimitExceeded.VpcSubnetNum"
//  MISSINGPARAMETER = "MissingParameter"
//  MISSINGPARAMETER_DPDKINSTANCETYPEREQUIREDVPC = "MissingParameter.DPDKInstanceTypeRequiredVPC"
//  MISSINGPARAMETER_MONITORSERVICE = "MissingParameter.MonitorService"
//  RESOURCEINSUFFICIENT_AVAILABILITYZONESOLDOUT = "ResourceInsufficient.AvailabilityZoneSoldOut"
//  RESOURCEINSUFFICIENT_CLOUDDISKSOLDOUT = "ResourceInsufficient.CloudDiskSoldOut"
//  RESOURCEINSUFFICIENT_CLOUDDISKUNAVAILABLE = "ResourceInsufficient.CloudDiskUnavailable"
//  RESOURCEINSUFFICIENT_DISASTERRECOVERGROUPCVMQUOTA = "ResourceInsufficient.DisasterRecoverGroupCvmQuota"
//  RESOURCEINSUFFICIENT_SPECIFIEDINSTANCETYPE = "ResourceInsufficient.SpecifiedInstanceType"
//  RESOURCENOTFOUND_HPCCLUSTER = "ResourceNotFound.HpcCluster"
//  RESOURCENOTFOUND_KEYPAIRNOTFOUND = "ResourceNotFound.KeyPairNotFound"
//  RESOURCENOTFOUND_NODEFAULTCBS = "ResourceNotFound.NoDefaultCbs"
//  RESOURCENOTFOUND_NODEFAULTCBSWITHREASON = "ResourceNotFound.NoDefaultCbsWithReason"
//  RESOURCEUNAVAILABLE_INSTANCETYPE = "ResourceUnavailable.InstanceType"
//  RESOURCESSOLDOUT_EIPINSUFFICIENT = "ResourcesSoldOut.EipInsufficient"
//  RESOURCESSOLDOUT_SPECIFIEDINSTANCETYPE = "ResourcesSoldOut.SpecifiedInstanceType"
//  UNSUPPORTEDOPERATION_BANDWIDTHPACKAGEIDNOTSUPPORTED = "UnsupportedOperation.BandwidthPackageIdNotSupported"
//  UNSUPPORTEDOPERATION_INSTANCETYPENOTSUPPORTJUMBOFRAME = "UnsupportedOperation.InstanceTypeNotSupportJumboFrame"
//  UNSUPPORTEDOPERATION_INVALIDDISK = "UnsupportedOperation.InvalidDisk"
//  UNSUPPORTEDOPERATION_KEYPAIRUNSUPPORTEDWINDOWS = "UnsupportedOperation.KeyPairUnsupportedWindows"
//  UNSUPPORTEDOPERATION_NOINSTANCETYPESUPPORTSPOT = "UnsupportedOperation.NoInstanceTypeSupportSpot"
//  UNSUPPORTEDOPERATION_NOTSUPPORTIMPORTINSTANCESACTIONTIMER = "UnsupportedOperation.NotSupportImportInstancesActionTimer"
//  UNSUPPORTEDOPERATION_ONLYFORPREPAIDACCOUNT = "UnsupportedOperation.OnlyForPrepaidAccount"
//  VPCADDRNOTINSUBNET = "VpcAddrNotInSubNet"
//  VPCIPISUSED = "VpcIpIsUsed"
func (c *Client) CreateLaunchTemplate(request *CreateLaunchTemplateRequest) (response *CreateLaunchTemplateResponse, err error) {
    return c.CreateLaunchTemplateWithContext(context.Background(), request)
}

// CreateLaunchTemplate
// This interface (CreateLaunchTemplate) is used for instance launch template creation.
//
// 
//
// An instance launch template is a configuration data and can be used to create instances. Its content includes configurations required to create instances, such as instance type, types and sizes of data disk and system disk, and security group and other information.
//
// 
//
// This API is used to create an instance launch template. After the initial creation of the instance template, its template version is the default version 1. A new version can be created using CreateLaunchTemplateVersion (https://intl.cloud.tencent.com/document/product/213/66326?from_cn_redirect=1), and the version number will increment. By default, when specifying an instance launch template in RunInstances (https://intl.cloud.tencent.com/document/product/213/15730?from_cn_redirect=1), if the template version number is not specified, the default version will be used.
//
// error code that may be returned:
//  ACCOUNTQUALIFICATIONRESTRICTIONS = "AccountQualificationRestrictions"
//  AUTHFAILURE_CAMROLENAMEAUTHENTICATEFAILED = "AuthFailure.CamRoleNameAuthenticateFailed"
//  FAILEDOPERATION_DISASTERRECOVERGROUPNOTFOUND = "FailedOperation.DisasterRecoverGroupNotFound"
//  FAILEDOPERATION_INQUIRYPRICEFAILED = "FailedOperation.InquiryPriceFailed"
//  FAILEDOPERATION_NOAVAILABLEIPADDRESSCOUNTINSUBNET = "FailedOperation.NoAvailableIpAddressCountInSubnet"
//  FAILEDOPERATION_SECURITYGROUPACTIONFAILED = "FailedOperation.SecurityGroupActionFailed"
//  FAILEDOPERATION_SNAPSHOTSIZELARGERTHANDATASIZE = "FailedOperation.SnapshotSizeLargerThanDataSize"
//  FAILEDOPERATION_TAGKEYRESERVED = "FailedOperation.TagKeyReserved"
//  INSTANCESQUOTALIMITEXCEEDED = "InstancesQuotaLimitExceeded"
//  INTERNALERROR_TRADEUNKNOWNERROR = "InternalError.TradeUnknownError"
//  INVALIDACCOUNT_INSUFFICIENTBALANCE = "InvalidAccount.InsufficientBalance"
//  INVALIDCLIENTTOKEN_TOOLONG = "InvalidClientToken.TooLong"
//  INVALIDHOSTID_MALFORMED = "InvalidHostId.Malformed"
//  INVALIDHOSTID_NOTFOUND = "InvalidHostId.NotFound"
//  INVALIDIMAGEID_MALFORMED = "InvalidImageId.Malformed"
//  INVALIDIMAGEID_NOTFOUND = "InvalidImageId.NotFound"
//  INVALIDINSTANCENAME_TOOLONG = "InvalidInstanceName.TooLong"
//  INVALIDINSTANCETYPE_MALFORMED = "InvalidInstanceType.Malformed"
//  INVALIDPARAMETER_INSTANCEIMAGENOTSUPPORT = "InvalidParameter.InstanceImageNotSupport"
//  INVALIDPARAMETER_INTERNETACCESSIBLENOTSUPPORTED = "InvalidParameter.InternetAccessibleNotSupported"
//  INVALIDPARAMETER_INVALIDIPFORMAT = "InvalidParameter.InvalidIpFormat"
//  INVALIDPARAMETER_LACKCORECOUNTORTHREADPERCORE = "InvalidParameter.LackCoreCountOrThreadPerCore"
//  INVALIDPARAMETER_PASSWORDNOTSUPPORTED = "InvalidParameter.PasswordNotSupported"
//  INVALIDPARAMETER_SNAPSHOTNOTFOUND = "InvalidParameter.SnapshotNotFound"
//  INVALIDPARAMETERCOMBINATION = "InvalidParameterCombination"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_CLOUDSSDDATADISKSIZETOOSMALL = "InvalidParameterValue.CloudSsdDataDiskSizeTooSmall"
//  INVALIDPARAMETERVALUE_CORECOUNTVALUE = "InvalidParameterValue.CoreCountValue"
//  INVALIDPARAMETERVALUE_ILLEGALHOSTNAME = "InvalidParameterValue.IllegalHostName"
//  INVALIDPARAMETERVALUE_INSTANCETYPENOTSUPPORTHPCCLUSTER = "InvalidParameterValue.InstanceTypeNotSupportHpcCluster"
//  INVALIDPARAMETERVALUE_INSTANCETYPEREQUIREDHPCCLUSTER = "InvalidParameterValue.InstanceTypeRequiredHpcCluster"
//  INVALIDPARAMETERVALUE_INSUFFICIENTOFFERING = "InvalidParameterValue.InsufficientOffering"
//  INVALIDPARAMETERVALUE_INSUFFICIENTPRICE = "InvalidParameterValue.InsufficientPrice"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEFORGIVENINSTANCETYPE = "InvalidParameterValue.InvalidImageForGivenInstanceType"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEID = "InvalidParameterValue.InvalidImageId"
//  INVALIDPARAMETERVALUE_INVALIDIMAGESTATE = "InvalidParameterValue.InvalidImageState"
//  INVALIDPARAMETERVALUE_INVALIDIPFORMAT = "InvalidParameterValue.InvalidIpFormat"
//  INVALIDPARAMETERVALUE_INVALIDLAUNCHTEMPLATEDESCRIPTION = "InvalidParameterValue.InvalidLaunchTemplateDescription"
//  INVALIDPARAMETERVALUE_INVALIDLAUNCHTEMPLATENAME = "InvalidParameterValue.InvalidLaunchTemplateName"
//  INVALIDPARAMETERVALUE_INVALIDLAUNCHTEMPLATEVERSIONDESCRIPTION = "InvalidParameterValue.InvalidLaunchTemplateVersionDescription"
//  INVALIDPARAMETERVALUE_INVALIDPASSWORD = "InvalidParameterValue.InvalidPassword"
//  INVALIDPARAMETERVALUE_INVALIDTIMEFORMAT = "InvalidParameterValue.InvalidTimeFormat"
//  INVALIDPARAMETERVALUE_INVALIDUSERDATAFORMAT = "InvalidParameterValue.InvalidUserDataFormat"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_MUSTDHCPENABLEDVPC = "InvalidParameterValue.MustDhcpEnabledVpc"
//  INVALIDPARAMETERVALUE_RANGE = "InvalidParameterValue.Range"
//  INVALIDPARAMETERVALUE_SNAPSHOTIDMALFORMED = "InvalidParameterValue.SnapshotIdMalformed"
//  INVALIDPARAMETERVALUE_SUBNETNOTEXIST = "InvalidParameterValue.SubnetNotExist"
//  INVALIDPARAMETERVALUE_THREADPERCOREVALUE = "InvalidParameterValue.ThreadPerCoreValue"
//  INVALIDPARAMETERVALUE_VPCIDMALFORMED = "InvalidParameterValue.VpcIdMalformed"
//  INVALIDPARAMETERVALUE_VPCIDNOTEXIST = "InvalidParameterValue.VpcIdNotExist"
//  INVALIDPARAMETERVALUE_VPCIDZONEIDNOTMATCH = "InvalidParameterValue.VpcIdZoneIdNotMatch"
//  INVALIDPARAMETERVALUE_ZONENOTSUPPORTED = "InvalidParameterValue.ZoneNotSupported"
//  INVALIDPASSWORD = "InvalidPassword"
//  INVALIDPERIOD = "InvalidPeriod"
//  INVALIDPERMISSION = "InvalidPermission"
//  INVALIDPROJECTID_NOTFOUND = "InvalidProjectId.NotFound"
//  INVALIDSECURITYGROUPID_NOTFOUND = "InvalidSecurityGroupId.NotFound"
//  INVALIDZONE_MISMATCHREGION = "InvalidZone.MismatchRegion"
//  LIMITEXCEEDED_INSTANCEQUOTA = "LimitExceeded.InstanceQuota"
//  LIMITEXCEEDED_LAUNCHTEMPLATEQUOTA = "LimitExceeded.LaunchTemplateQuota"
//  LIMITEXCEEDED_SINGLEUSGQUOTA = "LimitExceeded.SingleUSGQuota"
//  LIMITEXCEEDED_SPOTQUOTA = "LimitExceeded.SpotQuota"
//  LIMITEXCEEDED_USERSPOTQUOTA = "LimitExceeded.UserSpotQuota"
//  LIMITEXCEEDED_VPCSUBNETNUM = "LimitExceeded.VpcSubnetNum"
//  MISSINGPARAMETER = "MissingParameter"
//  MISSINGPARAMETER_DPDKINSTANCETYPEREQUIREDVPC = "MissingParameter.DPDKInstanceTypeRequiredVPC"
//  MISSINGPARAMETER_MONITORSERVICE = "MissingParameter.MonitorService"
//  RESOURCEINSUFFICIENT_AVAILABILITYZONESOLDOUT = "ResourceInsufficient.AvailabilityZoneSoldOut"
//  RESOURCEINSUFFICIENT_CLOUDDISKSOLDOUT = "ResourceInsufficient.CloudDiskSoldOut"
//  RESOURCEINSUFFICIENT_CLOUDDISKUNAVAILABLE = "ResourceInsufficient.CloudDiskUnavailable"
//  RESOURCEINSUFFICIENT_DISASTERRECOVERGROUPCVMQUOTA = "ResourceInsufficient.DisasterRecoverGroupCvmQuota"
//  RESOURCEINSUFFICIENT_SPECIFIEDINSTANCETYPE = "ResourceInsufficient.SpecifiedInstanceType"
//  RESOURCENOTFOUND_HPCCLUSTER = "ResourceNotFound.HpcCluster"
//  RESOURCENOTFOUND_KEYPAIRNOTFOUND = "ResourceNotFound.KeyPairNotFound"
//  RESOURCENOTFOUND_NODEFAULTCBS = "ResourceNotFound.NoDefaultCbs"
//  RESOURCENOTFOUND_NODEFAULTCBSWITHREASON = "ResourceNotFound.NoDefaultCbsWithReason"
//  RESOURCEUNAVAILABLE_INSTANCETYPE = "ResourceUnavailable.InstanceType"
//  RESOURCESSOLDOUT_EIPINSUFFICIENT = "ResourcesSoldOut.EipInsufficient"
//  RESOURCESSOLDOUT_SPECIFIEDINSTANCETYPE = "ResourcesSoldOut.SpecifiedInstanceType"
//  UNSUPPORTEDOPERATION_BANDWIDTHPACKAGEIDNOTSUPPORTED = "UnsupportedOperation.BandwidthPackageIdNotSupported"
//  UNSUPPORTEDOPERATION_INSTANCETYPENOTSUPPORTJUMBOFRAME = "UnsupportedOperation.InstanceTypeNotSupportJumboFrame"
//  UNSUPPORTEDOPERATION_INVALIDDISK = "UnsupportedOperation.InvalidDisk"
//  UNSUPPORTEDOPERATION_KEYPAIRUNSUPPORTEDWINDOWS = "UnsupportedOperation.KeyPairUnsupportedWindows"
//  UNSUPPORTEDOPERATION_NOINSTANCETYPESUPPORTSPOT = "UnsupportedOperation.NoInstanceTypeSupportSpot"
//  UNSUPPORTEDOPERATION_NOTSUPPORTIMPORTINSTANCESACTIONTIMER = "UnsupportedOperation.NotSupportImportInstancesActionTimer"
//  UNSUPPORTEDOPERATION_ONLYFORPREPAIDACCOUNT = "UnsupportedOperation.OnlyForPrepaidAccount"
//  VPCADDRNOTINSUBNET = "VpcAddrNotInSubNet"
//  VPCIPISUSED = "VpcIpIsUsed"
func (c *Client) CreateLaunchTemplateWithContext(ctx context.Context, request *CreateLaunchTemplateRequest) (response *CreateLaunchTemplateResponse, err error) {
    if request == nil {
        request = NewCreateLaunchTemplateRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "CreateLaunchTemplate")
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateLaunchTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateLaunchTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewCreateLaunchTemplateVersionRequest() (request *CreateLaunchTemplateVersionRequest) {
    request = &CreateLaunchTemplateVersionRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "CreateLaunchTemplateVersion")
    
    
    return
}

func NewCreateLaunchTemplateVersionResponse() (response *CreateLaunchTemplateVersionResponse) {
    response = &CreateLaunchTemplateVersionResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// CreateLaunchTemplateVersion
// This API is used to create an instance launch template based on the specified template ID and the corresponding template version number. The default version number will be used when no template version numbers are specified. Each instance launch template can have up to 30 version numbers.
//
// error code that may be returned:
//  AUTHFAILURE_CAMROLENAMEAUTHENTICATEFAILED = "AuthFailure.CamRoleNameAuthenticateFailed"
//  FAILEDOPERATION_DISASTERRECOVERGROUPNOTFOUND = "FailedOperation.DisasterRecoverGroupNotFound"
//  FAILEDOPERATION_INQUIRYPRICEFAILED = "FailedOperation.InquiryPriceFailed"
//  FAILEDOPERATION_TAGKEYRESERVED = "FailedOperation.TagKeyReserved"
//  INTERNALERROR_TRADEUNKNOWNERROR = "InternalError.TradeUnknownError"
//  INVALIDCLIENTTOKEN_TOOLONG = "InvalidClientToken.TooLong"
//  INVALIDHOSTID_MALFORMED = "InvalidHostId.Malformed"
//  INVALIDHOSTID_NOTFOUND = "InvalidHostId.NotFound"
//  INVALIDIMAGEID_MALFORMED = "InvalidImageId.Malformed"
//  INVALIDIMAGEID_NOTFOUND = "InvalidImageId.NotFound"
//  INVALIDINSTANCENAME_TOOLONG = "InvalidInstanceName.TooLong"
//  INVALIDINSTANCETYPE_MALFORMED = "InvalidInstanceType.Malformed"
//  INVALIDPARAMETER_INSTANCEIMAGENOTSUPPORT = "InvalidParameter.InstanceImageNotSupport"
//  INVALIDPARAMETER_INVALIDIPFORMAT = "InvalidParameter.InvalidIpFormat"
//  INVALIDPARAMETER_PASSWORDNOTSUPPORTED = "InvalidParameter.PasswordNotSupported"
//  INVALIDPARAMETER_SNAPSHOTNOTFOUND = "InvalidParameter.SnapshotNotFound"
//  INVALIDPARAMETERCOMBINATION = "InvalidParameterCombination"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_CLOUDSSDDATADISKSIZETOOSMALL = "InvalidParameterValue.CloudSsdDataDiskSizeTooSmall"
//  INVALIDPARAMETERVALUE_ILLEGALHOSTNAME = "InvalidParameterValue.IllegalHostName"
//  INVALIDPARAMETERVALUE_INSTANCETYPENOTSUPPORTHPCCLUSTER = "InvalidParameterValue.InstanceTypeNotSupportHpcCluster"
//  INVALIDPARAMETERVALUE_INVALIDIMAGESTATE = "InvalidParameterValue.InvalidImageState"
//  INVALIDPARAMETERVALUE_INVALIDIPFORMAT = "InvalidParameterValue.InvalidIpFormat"
//  INVALIDPARAMETERVALUE_INVALIDLAUNCHTEMPLATEVERSIONDESCRIPTION = "InvalidParameterValue.InvalidLaunchTemplateVersionDescription"
//  INVALIDPARAMETERVALUE_INVALIDUSERDATAFORMAT = "InvalidParameterValue.InvalidUserDataFormat"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDMALFORMED = "InvalidParameterValue.LaunchTemplateIdMalformed"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDNOTEXISTED = "InvalidParameterValue.LaunchTemplateIdNotExisted"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDVERNOTEXISTED = "InvalidParameterValue.LaunchTemplateIdVerNotExisted"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATENOTFOUND = "InvalidParameterValue.LaunchTemplateNotFound"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEVERSION = "InvalidParameterValue.LaunchTemplateVersion"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_MUSTDHCPENABLEDVPC = "InvalidParameterValue.MustDhcpEnabledVpc"
//  INVALIDPARAMETERVALUE_RANGE = "InvalidParameterValue.Range"
//  INVALIDPARAMETERVALUE_SNAPSHOTIDMALFORMED = "InvalidParameterValue.SnapshotIdMalformed"
//  INVALIDPARAMETERVALUE_SUBNETNOTEXIST = "InvalidParameterValue.SubnetNotExist"
//  INVALIDPARAMETERVALUE_THREADPERCOREVALUE = "InvalidParameterValue.ThreadPerCoreValue"
//  INVALIDPARAMETERVALUE_VPCIDZONEIDNOTMATCH = "InvalidParameterValue.VpcIdZoneIdNotMatch"
//  INVALIDPARAMETERVALUE_VPCNOTSUPPORTIPV6ADDRESS = "InvalidParameterValue.VpcNotSupportIpv6Address"
//  INVALIDPARAMETERVALUE_ZONENOTSUPPORTED = "InvalidParameterValue.ZoneNotSupported"
//  INVALIDPASSWORD = "InvalidPassword"
//  INVALIDPERIOD = "InvalidPeriod"
//  INVALIDPERMISSION = "InvalidPermission"
//  INVALIDPROJECTID_NOTFOUND = "InvalidProjectId.NotFound"
//  INVALIDSECURITYGROUPID_NOTFOUND = "InvalidSecurityGroupId.NotFound"
//  INVALIDZONE_MISMATCHREGION = "InvalidZone.MismatchRegion"
//  LIMITEXCEEDED_INSTANCEQUOTA = "LimitExceeded.InstanceQuota"
//  LIMITEXCEEDED_LAUNCHTEMPLATEQUOTA = "LimitExceeded.LaunchTemplateQuota"
//  LIMITEXCEEDED_LAUNCHTEMPLATEVERSIONQUOTA = "LimitExceeded.LaunchTemplateVersionQuota"
//  LIMITEXCEEDED_SINGLEUSGQUOTA = "LimitExceeded.SingleUSGQuota"
//  LIMITEXCEEDED_SPOTQUOTA = "LimitExceeded.SpotQuota"
//  LIMITEXCEEDED_USERSPOTQUOTA = "LimitExceeded.UserSpotQuota"
//  LIMITEXCEEDED_VPCSUBNETNUM = "LimitExceeded.VpcSubnetNum"
//  MISSINGPARAMETER = "MissingParameter"
//  MISSINGPARAMETER_DPDKINSTANCETYPEREQUIREDVPC = "MissingParameter.DPDKInstanceTypeRequiredVPC"
//  MISSINGPARAMETER_MONITORSERVICE = "MissingParameter.MonitorService"
//  RESOURCEINSUFFICIENT_CLOUDDISKSOLDOUT = "ResourceInsufficient.CloudDiskSoldOut"
//  RESOURCEINSUFFICIENT_CLOUDDISKUNAVAILABLE = "ResourceInsufficient.CloudDiskUnavailable"
//  RESOURCEINSUFFICIENT_DISASTERRECOVERGROUPCVMQUOTA = "ResourceInsufficient.DisasterRecoverGroupCvmQuota"
//  RESOURCEINSUFFICIENT_SPECIFIEDINSTANCETYPE = "ResourceInsufficient.SpecifiedInstanceType"
//  RESOURCENOTFOUND_HPCCLUSTER = "ResourceNotFound.HpcCluster"
//  RESOURCENOTFOUND_KEYPAIRNOTFOUND = "ResourceNotFound.KeyPairNotFound"
//  RESOURCENOTFOUND_NODEFAULTCBS = "ResourceNotFound.NoDefaultCbs"
//  RESOURCENOTFOUND_NODEFAULTCBSWITHREASON = "ResourceNotFound.NoDefaultCbsWithReason"
//  RESOURCEUNAVAILABLE_INSTANCETYPE = "ResourceUnavailable.InstanceType"
//  RESOURCESSOLDOUT_EIPINSUFFICIENT = "ResourcesSoldOut.EipInsufficient"
//  RESOURCESSOLDOUT_SPECIFIEDINSTANCETYPE = "ResourcesSoldOut.SpecifiedInstanceType"
//  UNKNOWNPARAMETER = "UnknownParameter"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
//  UNSUPPORTEDOPERATION_BANDWIDTHPACKAGEIDNOTSUPPORTED = "UnsupportedOperation.BandwidthPackageIdNotSupported"
//  UNSUPPORTEDOPERATION_INSTANCETYPENOTSUPPORTJUMBOFRAME = "UnsupportedOperation.InstanceTypeNotSupportJumboFrame"
//  UNSUPPORTEDOPERATION_INVALIDDISK = "UnsupportedOperation.InvalidDisk"
//  UNSUPPORTEDOPERATION_KEYPAIRUNSUPPORTEDWINDOWS = "UnsupportedOperation.KeyPairUnsupportedWindows"
//  UNSUPPORTEDOPERATION_NOINSTANCETYPESUPPORTSPOT = "UnsupportedOperation.NoInstanceTypeSupportSpot"
//  UNSUPPORTEDOPERATION_ONLYFORPREPAIDACCOUNT = "UnsupportedOperation.OnlyForPrepaidAccount"
//  VPCADDRNOTINSUBNET = "VpcAddrNotInSubNet"
//  VPCIPISUSED = "VpcIpIsUsed"
func (c *Client) CreateLaunchTemplateVersion(request *CreateLaunchTemplateVersionRequest) (response *CreateLaunchTemplateVersionResponse, err error) {
    return c.CreateLaunchTemplateVersionWithContext(context.Background(), request)
}

// CreateLaunchTemplateVersion
// This API is used to create an instance launch template based on the specified template ID and the corresponding template version number. The default version number will be used when no template version numbers are specified. Each instance launch template can have up to 30 version numbers.
//
// error code that may be returned:
//  AUTHFAILURE_CAMROLENAMEAUTHENTICATEFAILED = "AuthFailure.CamRoleNameAuthenticateFailed"
//  FAILEDOPERATION_DISASTERRECOVERGROUPNOTFOUND = "FailedOperation.DisasterRecoverGroupNotFound"
//  FAILEDOPERATION_INQUIRYPRICEFAILED = "FailedOperation.InquiryPriceFailed"
//  FAILEDOPERATION_TAGKEYRESERVED = "FailedOperation.TagKeyReserved"
//  INTERNALERROR_TRADEUNKNOWNERROR = "InternalError.TradeUnknownError"
//  INVALIDCLIENTTOKEN_TOOLONG = "InvalidClientToken.TooLong"
//  INVALIDHOSTID_MALFORMED = "InvalidHostId.Malformed"
//  INVALIDHOSTID_NOTFOUND = "InvalidHostId.NotFound"
//  INVALIDIMAGEID_MALFORMED = "InvalidImageId.Malformed"
//  INVALIDIMAGEID_NOTFOUND = "InvalidImageId.NotFound"
//  INVALIDINSTANCENAME_TOOLONG = "InvalidInstanceName.TooLong"
//  INVALIDINSTANCETYPE_MALFORMED = "InvalidInstanceType.Malformed"
//  INVALIDPARAMETER_INSTANCEIMAGENOTSUPPORT = "InvalidParameter.InstanceImageNotSupport"
//  INVALIDPARAMETER_INVALIDIPFORMAT = "InvalidParameter.InvalidIpFormat"
//  INVALIDPARAMETER_PASSWORDNOTSUPPORTED = "InvalidParameter.PasswordNotSupported"
//  INVALIDPARAMETER_SNAPSHOTNOTFOUND = "InvalidParameter.SnapshotNotFound"
//  INVALIDPARAMETERCOMBINATION = "InvalidParameterCombination"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_CLOUDSSDDATADISKSIZETOOSMALL = "InvalidParameterValue.CloudSsdDataDiskSizeTooSmall"
//  INVALIDPARAMETERVALUE_ILLEGALHOSTNAME = "InvalidParameterValue.IllegalHostName"
//  INVALIDPARAMETERVALUE_INSTANCETYPENOTSUPPORTHPCCLUSTER = "InvalidParameterValue.InstanceTypeNotSupportHpcCluster"
//  INVALIDPARAMETERVALUE_INVALIDIMAGESTATE = "InvalidParameterValue.InvalidImageState"
//  INVALIDPARAMETERVALUE_INVALIDIPFORMAT = "InvalidParameterValue.InvalidIpFormat"
//  INVALIDPARAMETERVALUE_INVALIDLAUNCHTEMPLATEVERSIONDESCRIPTION = "InvalidParameterValue.InvalidLaunchTemplateVersionDescription"
//  INVALIDPARAMETERVALUE_INVALIDUSERDATAFORMAT = "InvalidParameterValue.InvalidUserDataFormat"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDMALFORMED = "InvalidParameterValue.LaunchTemplateIdMalformed"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDNOTEXISTED = "InvalidParameterValue.LaunchTemplateIdNotExisted"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDVERNOTEXISTED = "InvalidParameterValue.LaunchTemplateIdVerNotExisted"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATENOTFOUND = "InvalidParameterValue.LaunchTemplateNotFound"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEVERSION = "InvalidParameterValue.LaunchTemplateVersion"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_MUSTDHCPENABLEDVPC = "InvalidParameterValue.MustDhcpEnabledVpc"
//  INVALIDPARAMETERVALUE_RANGE = "InvalidParameterValue.Range"
//  INVALIDPARAMETERVALUE_SNAPSHOTIDMALFORMED = "InvalidParameterValue.SnapshotIdMalformed"
//  INVALIDPARAMETERVALUE_SUBNETNOTEXIST = "InvalidParameterValue.SubnetNotExist"
//  INVALIDPARAMETERVALUE_THREADPERCOREVALUE = "InvalidParameterValue.ThreadPerCoreValue"
//  INVALIDPARAMETERVALUE_VPCIDZONEIDNOTMATCH = "InvalidParameterValue.VpcIdZoneIdNotMatch"
//  INVALIDPARAMETERVALUE_VPCNOTSUPPORTIPV6ADDRESS = "InvalidParameterValue.VpcNotSupportIpv6Address"
//  INVALIDPARAMETERVALUE_ZONENOTSUPPORTED = "InvalidParameterValue.ZoneNotSupported"
//  INVALIDPASSWORD = "InvalidPassword"
//  INVALIDPERIOD = "InvalidPeriod"
//  INVALIDPERMISSION = "InvalidPermission"
//  INVALIDPROJECTID_NOTFOUND = "InvalidProjectId.NotFound"
//  INVALIDSECURITYGROUPID_NOTFOUND = "InvalidSecurityGroupId.NotFound"
//  INVALIDZONE_MISMATCHREGION = "InvalidZone.MismatchRegion"
//  LIMITEXCEEDED_INSTANCEQUOTA = "LimitExceeded.InstanceQuota"
//  LIMITEXCEEDED_LAUNCHTEMPLATEQUOTA = "LimitExceeded.LaunchTemplateQuota"
//  LIMITEXCEEDED_LAUNCHTEMPLATEVERSIONQUOTA = "LimitExceeded.LaunchTemplateVersionQuota"
//  LIMITEXCEEDED_SINGLEUSGQUOTA = "LimitExceeded.SingleUSGQuota"
//  LIMITEXCEEDED_SPOTQUOTA = "LimitExceeded.SpotQuota"
//  LIMITEXCEEDED_USERSPOTQUOTA = "LimitExceeded.UserSpotQuota"
//  LIMITEXCEEDED_VPCSUBNETNUM = "LimitExceeded.VpcSubnetNum"
//  MISSINGPARAMETER = "MissingParameter"
//  MISSINGPARAMETER_DPDKINSTANCETYPEREQUIREDVPC = "MissingParameter.DPDKInstanceTypeRequiredVPC"
//  MISSINGPARAMETER_MONITORSERVICE = "MissingParameter.MonitorService"
//  RESOURCEINSUFFICIENT_CLOUDDISKSOLDOUT = "ResourceInsufficient.CloudDiskSoldOut"
//  RESOURCEINSUFFICIENT_CLOUDDISKUNAVAILABLE = "ResourceInsufficient.CloudDiskUnavailable"
//  RESOURCEINSUFFICIENT_DISASTERRECOVERGROUPCVMQUOTA = "ResourceInsufficient.DisasterRecoverGroupCvmQuota"
//  RESOURCEINSUFFICIENT_SPECIFIEDINSTANCETYPE = "ResourceInsufficient.SpecifiedInstanceType"
//  RESOURCENOTFOUND_HPCCLUSTER = "ResourceNotFound.HpcCluster"
//  RESOURCENOTFOUND_KEYPAIRNOTFOUND = "ResourceNotFound.KeyPairNotFound"
//  RESOURCENOTFOUND_NODEFAULTCBS = "ResourceNotFound.NoDefaultCbs"
//  RESOURCENOTFOUND_NODEFAULTCBSWITHREASON = "ResourceNotFound.NoDefaultCbsWithReason"
//  RESOURCEUNAVAILABLE_INSTANCETYPE = "ResourceUnavailable.InstanceType"
//  RESOURCESSOLDOUT_EIPINSUFFICIENT = "ResourcesSoldOut.EipInsufficient"
//  RESOURCESSOLDOUT_SPECIFIEDINSTANCETYPE = "ResourcesSoldOut.SpecifiedInstanceType"
//  UNKNOWNPARAMETER = "UnknownParameter"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
//  UNSUPPORTEDOPERATION_BANDWIDTHPACKAGEIDNOTSUPPORTED = "UnsupportedOperation.BandwidthPackageIdNotSupported"
//  UNSUPPORTEDOPERATION_INSTANCETYPENOTSUPPORTJUMBOFRAME = "UnsupportedOperation.InstanceTypeNotSupportJumboFrame"
//  UNSUPPORTEDOPERATION_INVALIDDISK = "UnsupportedOperation.InvalidDisk"
//  UNSUPPORTEDOPERATION_KEYPAIRUNSUPPORTEDWINDOWS = "UnsupportedOperation.KeyPairUnsupportedWindows"
//  UNSUPPORTEDOPERATION_NOINSTANCETYPESUPPORTSPOT = "UnsupportedOperation.NoInstanceTypeSupportSpot"
//  UNSUPPORTEDOPERATION_ONLYFORPREPAIDACCOUNT = "UnsupportedOperation.OnlyForPrepaidAccount"
//  VPCADDRNOTINSUBNET = "VpcAddrNotInSubNet"
//  VPCIPISUSED = "VpcIpIsUsed"
func (c *Client) CreateLaunchTemplateVersionWithContext(ctx context.Context, request *CreateLaunchTemplateVersionRequest) (response *CreateLaunchTemplateVersionResponse, err error) {
    if request == nil {
        request = NewCreateLaunchTemplateVersionRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "CreateLaunchTemplateVersion")
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateLaunchTemplateVersion require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateLaunchTemplateVersionResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteDisasterRecoverGroupsRequest() (request *DeleteDisasterRecoverGroupsRequest) {
    request = &DeleteDisasterRecoverGroupsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "DeleteDisasterRecoverGroups")
    
    
    return
}

func NewDeleteDisasterRecoverGroupsResponse() (response *DeleteDisasterRecoverGroupsResponse) {
    response = &DeleteDisasterRecoverGroupsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DeleteDisasterRecoverGroups
// This API is used to delete a [spread placement group](https://intl.cloud.tencent.com/document/product/213/15486?from_cn_redirect=1). Only empty placement groups can be deleted. To delete a non-empty group, you need to terminate all the CVM instances in it first. Otherwise, the deletion will fail.
//
// error code that may be returned:
//  FAILEDOPERATION_PLACEMENTSETNOTEMPTY = "FailedOperation.PlacementSetNotEmpty"
//  INVALIDPARAMETERVALUE_DISASTERRECOVERGROUPIDMALFORMED = "InvalidParameterValue.DisasterRecoverGroupIdMalformed"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  RESOURCEINSUFFICIENT_DISASTERRECOVERGROUPCVMQUOTA = "ResourceInsufficient.DisasterRecoverGroupCvmQuota"
//  RESOURCENOTFOUND_INVALIDPLACEMENTSET = "ResourceNotFound.InvalidPlacementSet"
func (c *Client) DeleteDisasterRecoverGroups(request *DeleteDisasterRecoverGroupsRequest) (response *DeleteDisasterRecoverGroupsResponse, err error) {
    return c.DeleteDisasterRecoverGroupsWithContext(context.Background(), request)
}

// DeleteDisasterRecoverGroups
// This API is used to delete a [spread placement group](https://intl.cloud.tencent.com/document/product/213/15486?from_cn_redirect=1). Only empty placement groups can be deleted. To delete a non-empty group, you need to terminate all the CVM instances in it first. Otherwise, the deletion will fail.
//
// error code that may be returned:
//  FAILEDOPERATION_PLACEMENTSETNOTEMPTY = "FailedOperation.PlacementSetNotEmpty"
//  INVALIDPARAMETERVALUE_DISASTERRECOVERGROUPIDMALFORMED = "InvalidParameterValue.DisasterRecoverGroupIdMalformed"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  RESOURCEINSUFFICIENT_DISASTERRECOVERGROUPCVMQUOTA = "ResourceInsufficient.DisasterRecoverGroupCvmQuota"
//  RESOURCENOTFOUND_INVALIDPLACEMENTSET = "ResourceNotFound.InvalidPlacementSet"
func (c *Client) DeleteDisasterRecoverGroupsWithContext(ctx context.Context, request *DeleteDisasterRecoverGroupsRequest) (response *DeleteDisasterRecoverGroupsResponse, err error) {
    if request == nil {
        request = NewDeleteDisasterRecoverGroupsRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "DeleteDisasterRecoverGroups")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteDisasterRecoverGroups require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteDisasterRecoverGroupsResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteImagesRequest() (request *DeleteImagesRequest) {
    request = &DeleteImagesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "DeleteImages")
    
    
    return
}

func NewDeleteImagesResponse() (response *DeleteImagesResponse) {
    response = &DeleteImagesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DeleteImages
// This API is used to delete one or more images.
//
// 
//
// * If the [ImageState](https://intl.cloud.tencent.com/document/product/213/15753?from_cn_redirect=1#Image) of an image is `CREATING` or `USING`, the image cannot be deleted. Call the [DescribeImages](https://intl.cloud.tencent.com/document/api/213/15715?from_cn_redirect=1) API to query the image status.
//
// * Up to 10 custom images are allowed in each region. If you have run out of the quota, delete unused images to create new ones.
//
// * A shared image cannot be deleted.
//
// error code that may be returned:
//  INVALIDIMAGEID_INSHARED = "InvalidImageId.InShared"
//  INVALIDIMAGEID_INCORRECTSTATE = "InvalidImageId.IncorrectState"
//  INVALIDIMAGEID_MALFORMED = "InvalidImageId.Malformed"
//  INVALIDIMAGEID_NOTFOUND = "InvalidImageId.NotFound"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEID = "InvalidParameterValue.InvalidImageId"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEIDISSHARED = "InvalidParameterValue.InvalidImageIdIsShared"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  MUTEXOPERATION_TASKRUNNING = "MutexOperation.TaskRunning"
func (c *Client) DeleteImages(request *DeleteImagesRequest) (response *DeleteImagesResponse, err error) {
    return c.DeleteImagesWithContext(context.Background(), request)
}

// DeleteImages
// This API is used to delete one or more images.
//
// 
//
// * If the [ImageState](https://intl.cloud.tencent.com/document/product/213/15753?from_cn_redirect=1#Image) of an image is `CREATING` or `USING`, the image cannot be deleted. Call the [DescribeImages](https://intl.cloud.tencent.com/document/api/213/15715?from_cn_redirect=1) API to query the image status.
//
// * Up to 10 custom images are allowed in each region. If you have run out of the quota, delete unused images to create new ones.
//
// * A shared image cannot be deleted.
//
// error code that may be returned:
//  INVALIDIMAGEID_INSHARED = "InvalidImageId.InShared"
//  INVALIDIMAGEID_INCORRECTSTATE = "InvalidImageId.IncorrectState"
//  INVALIDIMAGEID_MALFORMED = "InvalidImageId.Malformed"
//  INVALIDIMAGEID_NOTFOUND = "InvalidImageId.NotFound"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEID = "InvalidParameterValue.InvalidImageId"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEIDISSHARED = "InvalidParameterValue.InvalidImageIdIsShared"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  MUTEXOPERATION_TASKRUNNING = "MutexOperation.TaskRunning"
func (c *Client) DeleteImagesWithContext(ctx context.Context, request *DeleteImagesRequest) (response *DeleteImagesResponse, err error) {
    if request == nil {
        request = NewDeleteImagesRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "DeleteImages")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteImages require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteImagesResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteKeyPairsRequest() (request *DeleteKeyPairsRequest) {
    request = &DeleteKeyPairsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "DeleteKeyPairs")
    
    
    return
}

func NewDeleteKeyPairsResponse() (response *DeleteKeyPairsResponse) {
    response = &DeleteKeyPairsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DeleteKeyPairs
// This API is used to delete the key pairs hosted in Tencent Cloud.
//
// 
//
// * You can delete multiple key pairs at the same time.
//
// * A key pair used by an instance or image cannot be deleted. Therefore, you need to verify whether all the key pairs have been deleted successfully.
//
// error code that may be returned:
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDKEYPAIR_LIMITEXCEEDED = "InvalidKeyPair.LimitExceeded"
//  INVALIDKEYPAIRID_MALFORMED = "InvalidKeyPairId.Malformed"
//  INVALIDKEYPAIRID_NOTFOUND = "InvalidKeyPairId.NotFound"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_KEYPAIRNOTSUPPORTED = "InvalidParameterValue.KeyPairNotSupported"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  MISSINGPARAMETER = "MissingParameter"
func (c *Client) DeleteKeyPairs(request *DeleteKeyPairsRequest) (response *DeleteKeyPairsResponse, err error) {
    return c.DeleteKeyPairsWithContext(context.Background(), request)
}

// DeleteKeyPairs
// This API is used to delete the key pairs hosted in Tencent Cloud.
//
// 
//
// * You can delete multiple key pairs at the same time.
//
// * A key pair used by an instance or image cannot be deleted. Therefore, you need to verify whether all the key pairs have been deleted successfully.
//
// error code that may be returned:
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDKEYPAIR_LIMITEXCEEDED = "InvalidKeyPair.LimitExceeded"
//  INVALIDKEYPAIRID_MALFORMED = "InvalidKeyPairId.Malformed"
//  INVALIDKEYPAIRID_NOTFOUND = "InvalidKeyPairId.NotFound"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_KEYPAIRNOTSUPPORTED = "InvalidParameterValue.KeyPairNotSupported"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  MISSINGPARAMETER = "MissingParameter"
func (c *Client) DeleteKeyPairsWithContext(ctx context.Context, request *DeleteKeyPairsRequest) (response *DeleteKeyPairsResponse, err error) {
    if request == nil {
        request = NewDeleteKeyPairsRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "DeleteKeyPairs")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteKeyPairs require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteKeyPairsResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteLaunchTemplateRequest() (request *DeleteLaunchTemplateRequest) {
    request = &DeleteLaunchTemplateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "DeleteLaunchTemplate")
    
    
    return
}

func NewDeleteLaunchTemplateResponse() (response *DeleteLaunchTemplateResponse) {
    response = &DeleteLaunchTemplateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DeleteLaunchTemplate
// This API is used to delete an instance launch template.
//
// error code that may be returned:
//  ACCOUNTQUALIFICATIONRESTRICTIONS = "AccountQualificationRestrictions"
//  AUTHFAILURE_CAMROLENAMEAUTHENTICATEFAILED = "AuthFailure.CamRoleNameAuthenticateFailed"
//  INVALIDPARAMETERCOMBINATION = "InvalidParameterCombination"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDMALFORMED = "InvalidParameterValue.LaunchTemplateIdMalformed"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDNOTEXISTED = "InvalidParameterValue.LaunchTemplateIdNotExisted"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATENOTFOUND = "InvalidParameterValue.LaunchTemplateNotFound"
func (c *Client) DeleteLaunchTemplate(request *DeleteLaunchTemplateRequest) (response *DeleteLaunchTemplateResponse, err error) {
    return c.DeleteLaunchTemplateWithContext(context.Background(), request)
}

// DeleteLaunchTemplate
// This API is used to delete an instance launch template.
//
// error code that may be returned:
//  ACCOUNTQUALIFICATIONRESTRICTIONS = "AccountQualificationRestrictions"
//  AUTHFAILURE_CAMROLENAMEAUTHENTICATEFAILED = "AuthFailure.CamRoleNameAuthenticateFailed"
//  INVALIDPARAMETERCOMBINATION = "InvalidParameterCombination"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDMALFORMED = "InvalidParameterValue.LaunchTemplateIdMalformed"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDNOTEXISTED = "InvalidParameterValue.LaunchTemplateIdNotExisted"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATENOTFOUND = "InvalidParameterValue.LaunchTemplateNotFound"
func (c *Client) DeleteLaunchTemplateWithContext(ctx context.Context, request *DeleteLaunchTemplateRequest) (response *DeleteLaunchTemplateResponse, err error) {
    if request == nil {
        request = NewDeleteLaunchTemplateRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "DeleteLaunchTemplate")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteLaunchTemplate require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteLaunchTemplateResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteLaunchTemplateVersionsRequest() (request *DeleteLaunchTemplateVersionsRequest) {
    request = &DeleteLaunchTemplateVersionsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "DeleteLaunchTemplateVersions")
    
    
    return
}

func NewDeleteLaunchTemplateVersionsResponse() (response *DeleteLaunchTemplateVersionsResponse) {
    response = &DeleteLaunchTemplateVersionsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DeleteLaunchTemplateVersions
// This API is used to delete one or more instance launch template versions.
//
// error code that may be returned:
//  INVALIDPARAMETERCOMBINATION = "InvalidParameterCombination"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEDEFAULTVERSION = "InvalidParameterValue.LaunchTemplateDefaultVersion"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDMALFORMED = "InvalidParameterValue.LaunchTemplateIdMalformed"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDNOTEXISTED = "InvalidParameterValue.LaunchTemplateIdNotExisted"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDVERNOTEXISTED = "InvalidParameterValue.LaunchTemplateIdVerNotExisted"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATENOTFOUND = "InvalidParameterValue.LaunchTemplateNotFound"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEVERSION = "InvalidParameterValue.LaunchTemplateVersion"
//  MISSINGPARAMETER = "MissingParameter"
//  UNKNOWNPARAMETER = "UnknownParameter"
func (c *Client) DeleteLaunchTemplateVersions(request *DeleteLaunchTemplateVersionsRequest) (response *DeleteLaunchTemplateVersionsResponse, err error) {
    return c.DeleteLaunchTemplateVersionsWithContext(context.Background(), request)
}

// DeleteLaunchTemplateVersions
// This API is used to delete one or more instance launch template versions.
//
// error code that may be returned:
//  INVALIDPARAMETERCOMBINATION = "InvalidParameterCombination"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEDEFAULTVERSION = "InvalidParameterValue.LaunchTemplateDefaultVersion"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDMALFORMED = "InvalidParameterValue.LaunchTemplateIdMalformed"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDNOTEXISTED = "InvalidParameterValue.LaunchTemplateIdNotExisted"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDVERNOTEXISTED = "InvalidParameterValue.LaunchTemplateIdVerNotExisted"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATENOTFOUND = "InvalidParameterValue.LaunchTemplateNotFound"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEVERSION = "InvalidParameterValue.LaunchTemplateVersion"
//  MISSINGPARAMETER = "MissingParameter"
//  UNKNOWNPARAMETER = "UnknownParameter"
func (c *Client) DeleteLaunchTemplateVersionsWithContext(ctx context.Context, request *DeleteLaunchTemplateVersionsRequest) (response *DeleteLaunchTemplateVersionsResponse, err error) {
    if request == nil {
        request = NewDeleteLaunchTemplateVersionsRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "DeleteLaunchTemplateVersions")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteLaunchTemplateVersions require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteLaunchTemplateVersionsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeChcDeniedActionsRequest() (request *DescribeChcDeniedActionsRequest) {
    request = &DescribeChcDeniedActionsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "DescribeChcDeniedActions")
    
    
    return
}

func NewDescribeChcDeniedActionsResponse() (response *DescribeChcDeniedActionsResponse) {
    response = &DescribeChcDeniedActionsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeChcDeniedActions
// This API is used to query the actions not allowed for the specified CHC instances.
//
// error code that may be returned:
//  ACCOUNTQUALIFICATIONRESTRICTIONS = "AccountQualificationRestrictions"
//  INVALIDPARAMETERVALUE_CHCHOSTSNOTFOUND = "InvalidParameterValue.ChcHostsNotFound"
//  INVALIDPARAMETERVALUE_INCORRECTFORMAT = "InvalidParameterValue.IncorrectFormat"
func (c *Client) DescribeChcDeniedActions(request *DescribeChcDeniedActionsRequest) (response *DescribeChcDeniedActionsResponse, err error) {
    return c.DescribeChcDeniedActionsWithContext(context.Background(), request)
}

// DescribeChcDeniedActions
// This API is used to query the actions not allowed for the specified CHC instances.
//
// error code that may be returned:
//  ACCOUNTQUALIFICATIONRESTRICTIONS = "AccountQualificationRestrictions"
//  INVALIDPARAMETERVALUE_CHCHOSTSNOTFOUND = "InvalidParameterValue.ChcHostsNotFound"
//  INVALIDPARAMETERVALUE_INCORRECTFORMAT = "InvalidParameterValue.IncorrectFormat"
func (c *Client) DescribeChcDeniedActionsWithContext(ctx context.Context, request *DescribeChcDeniedActionsRequest) (response *DescribeChcDeniedActionsResponse, err error) {
    if request == nil {
        request = NewDescribeChcDeniedActionsRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "DescribeChcDeniedActions")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeChcDeniedActions require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeChcDeniedActionsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeChcHostsRequest() (request *DescribeChcHostsRequest) {
    request = &DescribeChcHostsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "DescribeChcHosts")
    
    
    return
}

func NewDescribeChcHostsResponse() (response *DescribeChcHostsResponse) {
    response = &DescribeChcHostsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeChcHosts
// This API is used to query the details of one or more CHC host.
//
// 
//
// * You can filter the query results with the instance ID, name or device type. See `Filter` for more information.
//
// * If no parameter is defined, a certain number of instances under the current account will be returned. The number is specified by `Limit` and is `20` by default.
//
// error code that may be returned:
//  ACCOUNTQUALIFICATIONRESTRICTIONS = "AccountQualificationRestrictions"
//  INVALIDFILTER = "InvalidFilter"
//  INVALIDFILTERVALUE_LIMITEXCEEDED = "InvalidFilterValue.LimitExceeded"
//  INVALIDPARAMETER_ATMOSTONE = "InvalidParameter.AtMostOne"
//  INVALIDPARAMETERVALUE_CHCHOSTSNOTFOUND = "InvalidParameterValue.ChcHostsNotFound"
//  INVALIDPARAMETERVALUE_INCORRECTFORMAT = "InvalidParameterValue.IncorrectFormat"
//  INVALIDPARAMETERVALUE_NOTEMPTY = "InvalidParameterValue.NotEmpty"
//  INVALIDPARAMETERVALUE_RANGE = "InvalidParameterValue.Range"
//  INVALIDPARAMETERVALUE_SUBNETIDMALFORMED = "InvalidParameterValue.SubnetIdMalformed"
//  INVALIDPARAMETERVALUE_VPCIDMALFORMED = "InvalidParameterValue.VpcIdMalformed"
//  INVALIDPARAMETERVALUELIMIT = "InvalidParameterValueLimit"
//  INVALIDPARAMETERVALUEOFFSET = "InvalidParameterValueOffset"
//  INVALIDZONE_MISMATCHREGION = "InvalidZone.MismatchRegion"
//  MISSINGPARAMETER = "MissingParameter"
func (c *Client) DescribeChcHosts(request *DescribeChcHostsRequest) (response *DescribeChcHostsResponse, err error) {
    return c.DescribeChcHostsWithContext(context.Background(), request)
}

// DescribeChcHosts
// This API is used to query the details of one or more CHC host.
//
// 
//
// * You can filter the query results with the instance ID, name or device type. See `Filter` for more information.
//
// * If no parameter is defined, a certain number of instances under the current account will be returned. The number is specified by `Limit` and is `20` by default.
//
// error code that may be returned:
//  ACCOUNTQUALIFICATIONRESTRICTIONS = "AccountQualificationRestrictions"
//  INVALIDFILTER = "InvalidFilter"
//  INVALIDFILTERVALUE_LIMITEXCEEDED = "InvalidFilterValue.LimitExceeded"
//  INVALIDPARAMETER_ATMOSTONE = "InvalidParameter.AtMostOne"
//  INVALIDPARAMETERVALUE_CHCHOSTSNOTFOUND = "InvalidParameterValue.ChcHostsNotFound"
//  INVALIDPARAMETERVALUE_INCORRECTFORMAT = "InvalidParameterValue.IncorrectFormat"
//  INVALIDPARAMETERVALUE_NOTEMPTY = "InvalidParameterValue.NotEmpty"
//  INVALIDPARAMETERVALUE_RANGE = "InvalidParameterValue.Range"
//  INVALIDPARAMETERVALUE_SUBNETIDMALFORMED = "InvalidParameterValue.SubnetIdMalformed"
//  INVALIDPARAMETERVALUE_VPCIDMALFORMED = "InvalidParameterValue.VpcIdMalformed"
//  INVALIDPARAMETERVALUELIMIT = "InvalidParameterValueLimit"
//  INVALIDPARAMETERVALUEOFFSET = "InvalidParameterValueOffset"
//  INVALIDZONE_MISMATCHREGION = "InvalidZone.MismatchRegion"
//  MISSINGPARAMETER = "MissingParameter"
func (c *Client) DescribeChcHostsWithContext(ctx context.Context, request *DescribeChcHostsRequest) (response *DescribeChcHostsResponse, err error) {
    if request == nil {
        request = NewDescribeChcHostsRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "DescribeChcHosts")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeChcHosts require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeChcHostsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeDisasterRecoverGroupQuotaRequest() (request *DescribeDisasterRecoverGroupQuotaRequest) {
    request = &DescribeDisasterRecoverGroupQuotaRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "DescribeDisasterRecoverGroupQuota")
    
    
    return
}

func NewDescribeDisasterRecoverGroupQuotaResponse() (response *DescribeDisasterRecoverGroupQuotaResponse) {
    response = &DescribeDisasterRecoverGroupQuotaResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeDisasterRecoverGroupQuota
// This API is used to query the quota of [spread placement groups](https://intl.cloud.tencent.com/document/product/213/15486?from_cn_redirect=1).
//
// error code that may be returned:
//  ACCOUNTQUALIFICATIONRESTRICTIONS = "AccountQualificationRestrictions"
//  INVALIDFILTER = "InvalidFilter"
//  INVALIDFILTERVALUE_LIMITEXCEEDED = "InvalidFilterValue.LimitExceeded"
//  INVALIDPARAMETER_ATMOSTONE = "InvalidParameter.AtMostOne"
//  INVALIDPARAMETERVALUE_CHCHOSTSNOTFOUND = "InvalidParameterValue.ChcHostsNotFound"
//  INVALIDPARAMETERVALUE_INCORRECTFORMAT = "InvalidParameterValue.IncorrectFormat"
//  INVALIDPARAMETERVALUE_NOTEMPTY = "InvalidParameterValue.NotEmpty"
//  INVALIDPARAMETERVALUE_RANGE = "InvalidParameterValue.Range"
//  INVALIDPARAMETERVALUE_SUBNETIDMALFORMED = "InvalidParameterValue.SubnetIdMalformed"
//  INVALIDPARAMETERVALUE_VPCIDMALFORMED = "InvalidParameterValue.VpcIdMalformed"
//  INVALIDPARAMETERVALUELIMIT = "InvalidParameterValueLimit"
//  INVALIDPARAMETERVALUEOFFSET = "InvalidParameterValueOffset"
//  INVALIDZONE_MISMATCHREGION = "InvalidZone.MismatchRegion"
//  MISSINGPARAMETER = "MissingParameter"
func (c *Client) DescribeDisasterRecoverGroupQuota(request *DescribeDisasterRecoverGroupQuotaRequest) (response *DescribeDisasterRecoverGroupQuotaResponse, err error) {
    return c.DescribeDisasterRecoverGroupQuotaWithContext(context.Background(), request)
}

// DescribeDisasterRecoverGroupQuota
// This API is used to query the quota of [spread placement groups](https://intl.cloud.tencent.com/document/product/213/15486?from_cn_redirect=1).
//
// error code that may be returned:
//  ACCOUNTQUALIFICATIONRESTRICTIONS = "AccountQualificationRestrictions"
//  INVALIDFILTER = "InvalidFilter"
//  INVALIDFILTERVALUE_LIMITEXCEEDED = "InvalidFilterValue.LimitExceeded"
//  INVALIDPARAMETER_ATMOSTONE = "InvalidParameter.AtMostOne"
//  INVALIDPARAMETERVALUE_CHCHOSTSNOTFOUND = "InvalidParameterValue.ChcHostsNotFound"
//  INVALIDPARAMETERVALUE_INCORRECTFORMAT = "InvalidParameterValue.IncorrectFormat"
//  INVALIDPARAMETERVALUE_NOTEMPTY = "InvalidParameterValue.NotEmpty"
//  INVALIDPARAMETERVALUE_RANGE = "InvalidParameterValue.Range"
//  INVALIDPARAMETERVALUE_SUBNETIDMALFORMED = "InvalidParameterValue.SubnetIdMalformed"
//  INVALIDPARAMETERVALUE_VPCIDMALFORMED = "InvalidParameterValue.VpcIdMalformed"
//  INVALIDPARAMETERVALUELIMIT = "InvalidParameterValueLimit"
//  INVALIDPARAMETERVALUEOFFSET = "InvalidParameterValueOffset"
//  INVALIDZONE_MISMATCHREGION = "InvalidZone.MismatchRegion"
//  MISSINGPARAMETER = "MissingParameter"
func (c *Client) DescribeDisasterRecoverGroupQuotaWithContext(ctx context.Context, request *DescribeDisasterRecoverGroupQuotaRequest) (response *DescribeDisasterRecoverGroupQuotaResponse, err error) {
    if request == nil {
        request = NewDescribeDisasterRecoverGroupQuotaRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "DescribeDisasterRecoverGroupQuota")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeDisasterRecoverGroupQuota require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeDisasterRecoverGroupQuotaResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeDisasterRecoverGroupsRequest() (request *DescribeDisasterRecoverGroupsRequest) {
    request = &DescribeDisasterRecoverGroupsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "DescribeDisasterRecoverGroups")
    
    
    return
}

func NewDescribeDisasterRecoverGroupsResponse() (response *DescribeDisasterRecoverGroupsResponse) {
    response = &DescribeDisasterRecoverGroupsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeDisasterRecoverGroups
// This API is used to query the information on [spread placement groups](https://intl.cloud.tencent.com/document/product/213/15486?from_cn_redirect=1).
//
// error code that may be returned:
//  INVALIDPARAMETERVALUE_DISASTERRECOVERGROUPIDMALFORMED = "InvalidParameterValue.DisasterRecoverGroupIdMalformed"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
func (c *Client) DescribeDisasterRecoverGroups(request *DescribeDisasterRecoverGroupsRequest) (response *DescribeDisasterRecoverGroupsResponse, err error) {
    return c.DescribeDisasterRecoverGroupsWithContext(context.Background(), request)
}

// DescribeDisasterRecoverGroups
// This API is used to query the information on [spread placement groups](https://intl.cloud.tencent.com/document/product/213/15486?from_cn_redirect=1).
//
// error code that may be returned:
//  INVALIDPARAMETERVALUE_DISASTERRECOVERGROUPIDMALFORMED = "InvalidParameterValue.DisasterRecoverGroupIdMalformed"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
func (c *Client) DescribeDisasterRecoverGroupsWithContext(ctx context.Context, request *DescribeDisasterRecoverGroupsRequest) (response *DescribeDisasterRecoverGroupsResponse, err error) {
    if request == nil {
        request = NewDescribeDisasterRecoverGroupsRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "DescribeDisasterRecoverGroups")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeDisasterRecoverGroups require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeDisasterRecoverGroupsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeHostsRequest() (request *DescribeHostsRequest) {
    request = &DescribeHostsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "DescribeHosts")
    
    
    return
}

func NewDescribeHostsResponse() (response *DescribeHostsResponse) {
    response = &DescribeHostsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeHosts
// This API is used to query the details of CDH instances.
//
// error code that may be returned:
//  FAILEDOPERATION_TAGKEYRESERVED = "FailedOperation.TagKeyReserved"
//  INVALIDFILTER = "InvalidFilter"
//  INVALIDFILTERVALUE_LIMITEXCEEDED = "InvalidFilterValue.LimitExceeded"
//  INVALIDHOSTID_MALFORMED = "InvalidHostId.Malformed"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDZONE_MISMATCHREGION = "InvalidZone.MismatchRegion"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
func (c *Client) DescribeHosts(request *DescribeHostsRequest) (response *DescribeHostsResponse, err error) {
    return c.DescribeHostsWithContext(context.Background(), request)
}

// DescribeHosts
// This API is used to query the details of CDH instances.
//
// error code that may be returned:
//  FAILEDOPERATION_TAGKEYRESERVED = "FailedOperation.TagKeyReserved"
//  INVALIDFILTER = "InvalidFilter"
//  INVALIDFILTERVALUE_LIMITEXCEEDED = "InvalidFilterValue.LimitExceeded"
//  INVALIDHOSTID_MALFORMED = "InvalidHostId.Malformed"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDZONE_MISMATCHREGION = "InvalidZone.MismatchRegion"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
func (c *Client) DescribeHostsWithContext(ctx context.Context, request *DescribeHostsRequest) (response *DescribeHostsResponse, err error) {
    if request == nil {
        request = NewDescribeHostsRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "DescribeHosts")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeHosts require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeHostsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeImageFromFamilyRequest() (request *DescribeImageFromFamilyRequest) {
    request = &DescribeImageFromFamilyRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "DescribeImageFromFamily")
    
    
    return
}

func NewDescribeImageFromFamilyResponse() (response *DescribeImageFromFamilyResponse) {
    response = &DescribeImageFromFamilyResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeImageFromFamily
// This API is used to view information about available images within an image family.
//
// error code that may be returned:
//  FAILEDOPERATION_TAGKEYRESERVED = "FailedOperation.TagKeyReserved"
//  INVALIDFILTER = "InvalidFilter"
//  INVALIDFILTERVALUE_LIMITEXCEEDED = "InvalidFilterValue.LimitExceeded"
//  INVALIDHOSTID_MALFORMED = "InvalidHostId.Malformed"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDZONE_MISMATCHREGION = "InvalidZone.MismatchRegion"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
func (c *Client) DescribeImageFromFamily(request *DescribeImageFromFamilyRequest) (response *DescribeImageFromFamilyResponse, err error) {
    return c.DescribeImageFromFamilyWithContext(context.Background(), request)
}

// DescribeImageFromFamily
// This API is used to view information about available images within an image family.
//
// error code that may be returned:
//  FAILEDOPERATION_TAGKEYRESERVED = "FailedOperation.TagKeyReserved"
//  INVALIDFILTER = "InvalidFilter"
//  INVALIDFILTERVALUE_LIMITEXCEEDED = "InvalidFilterValue.LimitExceeded"
//  INVALIDHOSTID_MALFORMED = "InvalidHostId.Malformed"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDZONE_MISMATCHREGION = "InvalidZone.MismatchRegion"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
func (c *Client) DescribeImageFromFamilyWithContext(ctx context.Context, request *DescribeImageFromFamilyRequest) (response *DescribeImageFromFamilyResponse, err error) {
    if request == nil {
        request = NewDescribeImageFromFamilyRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "DescribeImageFromFamily")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeImageFromFamily require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeImageFromFamilyResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeImageQuotaRequest() (request *DescribeImageQuotaRequest) {
    request = &DescribeImageQuotaRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "DescribeImageQuota")
    
    
    return
}

func NewDescribeImageQuotaResponse() (response *DescribeImageQuotaResponse) {
    response = &DescribeImageQuotaResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeImageQuota
// This API is used to query the image quota of an user account.
//
// error code that may be returned:
//  FAILEDOPERATION_TAGKEYRESERVED = "FailedOperation.TagKeyReserved"
//  INVALIDFILTER = "InvalidFilter"
//  INVALIDFILTERVALUE_LIMITEXCEEDED = "InvalidFilterValue.LimitExceeded"
//  INVALIDHOSTID_MALFORMED = "InvalidHostId.Malformed"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDZONE_MISMATCHREGION = "InvalidZone.MismatchRegion"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
func (c *Client) DescribeImageQuota(request *DescribeImageQuotaRequest) (response *DescribeImageQuotaResponse, err error) {
    return c.DescribeImageQuotaWithContext(context.Background(), request)
}

// DescribeImageQuota
// This API is used to query the image quota of an user account.
//
// error code that may be returned:
//  FAILEDOPERATION_TAGKEYRESERVED = "FailedOperation.TagKeyReserved"
//  INVALIDFILTER = "InvalidFilter"
//  INVALIDFILTERVALUE_LIMITEXCEEDED = "InvalidFilterValue.LimitExceeded"
//  INVALIDHOSTID_MALFORMED = "InvalidHostId.Malformed"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDZONE_MISMATCHREGION = "InvalidZone.MismatchRegion"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
func (c *Client) DescribeImageQuotaWithContext(ctx context.Context, request *DescribeImageQuotaRequest) (response *DescribeImageQuotaResponse, err error) {
    if request == nil {
        request = NewDescribeImageQuotaRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "DescribeImageQuota")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeImageQuota require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeImageQuotaResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeImageSharePermissionRequest() (request *DescribeImageSharePermissionRequest) {
    request = &DescribeImageSharePermissionRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "DescribeImageSharePermission")
    
    
    return
}

func NewDescribeImageSharePermissionResponse() (response *DescribeImageSharePermissionResponse) {
    response = &DescribeImageSharePermissionResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeImageSharePermission
// This API is used to query image sharing information.
//
// error code that may be returned:
//  INVALIDACCOUNTID_NOTFOUND = "InvalidAccountId.NotFound"
//  INVALIDACCOUNTIS_YOURSELF = "InvalidAccountIs.YourSelf"
//  INVALIDIMAGEID_MALFORMED = "InvalidImageId.Malformed"
//  INVALIDIMAGEID_NOTFOUND = "InvalidImageId.NotFound"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  OVERQUOTA = "OverQuota"
//  UNAUTHORIZEDOPERATION_IMAGENOTBELONGTOACCOUNT = "UnauthorizedOperation.ImageNotBelongToAccount"
func (c *Client) DescribeImageSharePermission(request *DescribeImageSharePermissionRequest) (response *DescribeImageSharePermissionResponse, err error) {
    return c.DescribeImageSharePermissionWithContext(context.Background(), request)
}

// DescribeImageSharePermission
// This API is used to query image sharing information.
//
// error code that may be returned:
//  INVALIDACCOUNTID_NOTFOUND = "InvalidAccountId.NotFound"
//  INVALIDACCOUNTIS_YOURSELF = "InvalidAccountIs.YourSelf"
//  INVALIDIMAGEID_MALFORMED = "InvalidImageId.Malformed"
//  INVALIDIMAGEID_NOTFOUND = "InvalidImageId.NotFound"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  OVERQUOTA = "OverQuota"
//  UNAUTHORIZEDOPERATION_IMAGENOTBELONGTOACCOUNT = "UnauthorizedOperation.ImageNotBelongToAccount"
func (c *Client) DescribeImageSharePermissionWithContext(ctx context.Context, request *DescribeImageSharePermissionRequest) (response *DescribeImageSharePermissionResponse, err error) {
    if request == nil {
        request = NewDescribeImageSharePermissionRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "DescribeImageSharePermission")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeImageSharePermission require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeImageSharePermissionResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeImagesRequest() (request *DescribeImagesRequest) {
    request = &DescribeImagesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "DescribeImages")
    
    
    return
}

func NewDescribeImagesResponse() (response *DescribeImagesResponse) {
    response = &DescribeImagesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeImages
// This API is used to view the list of images.
//
// 
//
// * You specify the image ID or set filters to query the details of certain images.
//
// * You can specify `Offset` and `Limit` to select a certain part of the results. By default, the information on the first 20 matching results is returned.
//
// error code that may be returned:
//  FAILEDOPERATION_ILLEGALTAGKEY = "FailedOperation.IllegalTagKey"
//  INVALIDFILTER = "InvalidFilter"
//  INVALIDFILTERVALUE_LIMITEXCEEDED = "InvalidFilterValue.LimitExceeded"
//  INVALIDIMAGEID_MALFORMED = "InvalidImageId.Malformed"
//  INVALIDINSTANCETYPE_MALFORMED = "InvalidInstanceType.Malformed"
//  INVALIDPARAMETER_INVALIDPARAMETERCOEXISTIMAGEIDSFILTERS = "InvalidParameter.InvalidParameterCoexistImageIdsFilters"
//  INVALIDPARAMETERVALUE_INSTANCETYPENOTFOUND = "InvalidParameterValue.InstanceTypeNotFound"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_INVALIDPARAMETERVALUELIMIT = "InvalidParameterValue.InvalidParameterValueLimit"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_TAGKEYNOTFOUND = "InvalidParameterValue.TagKeyNotFound"
//  INVALIDPARAMETERVALUE_ZONENOTSUPPORTED = "InvalidParameterValue.ZoneNotSupported"
//  INVALIDREGION_NOTFOUND = "InvalidRegion.NotFound"
//  INVALIDZONE_MISMATCHREGION = "InvalidZone.MismatchRegion"
//  RESOURCESSOLDOUT_SPECIFIEDINSTANCETYPE = "ResourcesSoldOut.SpecifiedInstanceType"
//  UNAUTHORIZEDOPERATION_INVALIDTOKEN = "UnauthorizedOperation.InvalidToken"
//  UNAUTHORIZEDOPERATION_PERMISSIONDENIED = "UnauthorizedOperation.PermissionDenied"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
func (c *Client) DescribeImages(request *DescribeImagesRequest) (response *DescribeImagesResponse, err error) {
    return c.DescribeImagesWithContext(context.Background(), request)
}

// DescribeImages
// This API is used to view the list of images.
//
// 
//
// * You specify the image ID or set filters to query the details of certain images.
//
// * You can specify `Offset` and `Limit` to select a certain part of the results. By default, the information on the first 20 matching results is returned.
//
// error code that may be returned:
//  FAILEDOPERATION_ILLEGALTAGKEY = "FailedOperation.IllegalTagKey"
//  INVALIDFILTER = "InvalidFilter"
//  INVALIDFILTERVALUE_LIMITEXCEEDED = "InvalidFilterValue.LimitExceeded"
//  INVALIDIMAGEID_MALFORMED = "InvalidImageId.Malformed"
//  INVALIDINSTANCETYPE_MALFORMED = "InvalidInstanceType.Malformed"
//  INVALIDPARAMETER_INVALIDPARAMETERCOEXISTIMAGEIDSFILTERS = "InvalidParameter.InvalidParameterCoexistImageIdsFilters"
//  INVALIDPARAMETERVALUE_INSTANCETYPENOTFOUND = "InvalidParameterValue.InstanceTypeNotFound"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_INVALIDPARAMETERVALUELIMIT = "InvalidParameterValue.InvalidParameterValueLimit"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_TAGKEYNOTFOUND = "InvalidParameterValue.TagKeyNotFound"
//  INVALIDPARAMETERVALUE_ZONENOTSUPPORTED = "InvalidParameterValue.ZoneNotSupported"
//  INVALIDREGION_NOTFOUND = "InvalidRegion.NotFound"
//  INVALIDZONE_MISMATCHREGION = "InvalidZone.MismatchRegion"
//  RESOURCESSOLDOUT_SPECIFIEDINSTANCETYPE = "ResourcesSoldOut.SpecifiedInstanceType"
//  UNAUTHORIZEDOPERATION_INVALIDTOKEN = "UnauthorizedOperation.InvalidToken"
//  UNAUTHORIZEDOPERATION_PERMISSIONDENIED = "UnauthorizedOperation.PermissionDenied"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
func (c *Client) DescribeImagesWithContext(ctx context.Context, request *DescribeImagesRequest) (response *DescribeImagesResponse, err error) {
    if request == nil {
        request = NewDescribeImagesRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "DescribeImages")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeImages require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeImagesResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeImportImageOsRequest() (request *DescribeImportImageOsRequest) {
    request = &DescribeImportImageOsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "DescribeImportImageOs")
    
    
    return
}

func NewDescribeImportImageOsResponse() (response *DescribeImportImageOsResponse) {
    response = &DescribeImportImageOsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeImportImageOs
// This API is used to query the list of supported operating systems of imported images.
//
// error code that may be returned:
//  FAILEDOPERATION_ILLEGALTAGKEY = "FailedOperation.IllegalTagKey"
//  INVALIDFILTER = "InvalidFilter"
//  INVALIDFILTERVALUE_LIMITEXCEEDED = "InvalidFilterValue.LimitExceeded"
//  INVALIDIMAGEID_MALFORMED = "InvalidImageId.Malformed"
//  INVALIDINSTANCETYPE_MALFORMED = "InvalidInstanceType.Malformed"
//  INVALIDPARAMETER_INVALIDPARAMETERCOEXISTIMAGEIDSFILTERS = "InvalidParameter.InvalidParameterCoexistImageIdsFilters"
//  INVALIDPARAMETERVALUE_INSTANCETYPENOTFOUND = "InvalidParameterValue.InstanceTypeNotFound"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_INVALIDPARAMETERVALUELIMIT = "InvalidParameterValue.InvalidParameterValueLimit"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_TAGKEYNOTFOUND = "InvalidParameterValue.TagKeyNotFound"
//  INVALIDPARAMETERVALUE_ZONENOTSUPPORTED = "InvalidParameterValue.ZoneNotSupported"
//  INVALIDREGION_NOTFOUND = "InvalidRegion.NotFound"
//  INVALIDZONE_MISMATCHREGION = "InvalidZone.MismatchRegion"
//  RESOURCESSOLDOUT_SPECIFIEDINSTANCETYPE = "ResourcesSoldOut.SpecifiedInstanceType"
//  UNAUTHORIZEDOPERATION_INVALIDTOKEN = "UnauthorizedOperation.InvalidToken"
//  UNAUTHORIZEDOPERATION_PERMISSIONDENIED = "UnauthorizedOperation.PermissionDenied"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
func (c *Client) DescribeImportImageOs(request *DescribeImportImageOsRequest) (response *DescribeImportImageOsResponse, err error) {
    return c.DescribeImportImageOsWithContext(context.Background(), request)
}

// DescribeImportImageOs
// This API is used to query the list of supported operating systems of imported images.
//
// error code that may be returned:
//  FAILEDOPERATION_ILLEGALTAGKEY = "FailedOperation.IllegalTagKey"
//  INVALIDFILTER = "InvalidFilter"
//  INVALIDFILTERVALUE_LIMITEXCEEDED = "InvalidFilterValue.LimitExceeded"
//  INVALIDIMAGEID_MALFORMED = "InvalidImageId.Malformed"
//  INVALIDINSTANCETYPE_MALFORMED = "InvalidInstanceType.Malformed"
//  INVALIDPARAMETER_INVALIDPARAMETERCOEXISTIMAGEIDSFILTERS = "InvalidParameter.InvalidParameterCoexistImageIdsFilters"
//  INVALIDPARAMETERVALUE_INSTANCETYPENOTFOUND = "InvalidParameterValue.InstanceTypeNotFound"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_INVALIDPARAMETERVALUELIMIT = "InvalidParameterValue.InvalidParameterValueLimit"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_TAGKEYNOTFOUND = "InvalidParameterValue.TagKeyNotFound"
//  INVALIDPARAMETERVALUE_ZONENOTSUPPORTED = "InvalidParameterValue.ZoneNotSupported"
//  INVALIDREGION_NOTFOUND = "InvalidRegion.NotFound"
//  INVALIDZONE_MISMATCHREGION = "InvalidZone.MismatchRegion"
//  RESOURCESSOLDOUT_SPECIFIEDINSTANCETYPE = "ResourcesSoldOut.SpecifiedInstanceType"
//  UNAUTHORIZEDOPERATION_INVALIDTOKEN = "UnauthorizedOperation.InvalidToken"
//  UNAUTHORIZEDOPERATION_PERMISSIONDENIED = "UnauthorizedOperation.PermissionDenied"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
func (c *Client) DescribeImportImageOsWithContext(ctx context.Context, request *DescribeImportImageOsRequest) (response *DescribeImportImageOsResponse, err error) {
    if request == nil {
        request = NewDescribeImportImageOsRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "DescribeImportImageOs")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeImportImageOs require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeImportImageOsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeInstanceFamilyConfigsRequest() (request *DescribeInstanceFamilyConfigsRequest) {
    request = &DescribeInstanceFamilyConfigsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "DescribeInstanceFamilyConfigs")
    
    
    return
}

func NewDescribeInstanceFamilyConfigsResponse() (response *DescribeInstanceFamilyConfigsResponse) {
    response = &DescribeInstanceFamilyConfigsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeInstanceFamilyConfigs
// This API is used to query a list of model families available to the current user in the current region.
//
// error code that may be returned:
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDREGION_NOTFOUND = "InvalidRegion.NotFound"
func (c *Client) DescribeInstanceFamilyConfigs(request *DescribeInstanceFamilyConfigsRequest) (response *DescribeInstanceFamilyConfigsResponse, err error) {
    return c.DescribeInstanceFamilyConfigsWithContext(context.Background(), request)
}

// DescribeInstanceFamilyConfigs
// This API is used to query a list of model families available to the current user in the current region.
//
// error code that may be returned:
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDREGION_NOTFOUND = "InvalidRegion.NotFound"
func (c *Client) DescribeInstanceFamilyConfigsWithContext(ctx context.Context, request *DescribeInstanceFamilyConfigsRequest) (response *DescribeInstanceFamilyConfigsResponse, err error) {
    if request == nil {
        request = NewDescribeInstanceFamilyConfigsRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "DescribeInstanceFamilyConfigs")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeInstanceFamilyConfigs require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeInstanceFamilyConfigsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeInstancesRequest() (request *DescribeInstancesRequest) {
    request = &DescribeInstancesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "DescribeInstances")
    
    
    return
}

func NewDescribeInstancesResponse() (response *DescribeInstancesResponse) {
    response = &DescribeInstancesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeInstances
// This API is used to query the details of instances.
//
// 
//
// * You can filter the query results with the instance `ID`, name, or billing method. See `Filter` for more information.
//
// * If no parameter is defined, a certain number of instances under the current account will be returned. The number is specified by `Limit` and is 20 by default.
//
// error code that may be returned:
//  FAILEDOPERATION_ILLEGALTAGKEY = "FailedOperation.IllegalTagKey"
//  FAILEDOPERATION_ILLEGALTAGVALUE = "FailedOperation.IllegalTagValue"
//  FAILEDOPERATION_TAGKEYRESERVED = "FailedOperation.TagKeyReserved"
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDFILTER = "InvalidFilter"
//  INVALIDFILTERVALUE_LIMITEXCEEDED = "InvalidFilterValue.LimitExceeded"
//  INVALIDHOSTID_MALFORMED = "InvalidHostId.Malformed"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_IPADDRESSMALFORMED = "InvalidParameterValue.IPAddressMalformed"
//  INVALIDPARAMETERVALUE_IPV6ADDRESSMALFORMED = "InvalidParameterValue.IPv6AddressMalformed"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_INVALIDIPFORMAT = "InvalidParameterValue.InvalidIpFormat"
//  INVALIDPARAMETERVALUE_INVALIDTIMEFORMAT = "InvalidParameterValue.InvalidTimeFormat"
//  INVALIDPARAMETERVALUE_INVALIDVAGUENAME = "InvalidParameterValue.InvalidVagueName"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_SUBNETIDMALFORMED = "InvalidParameterValue.SubnetIdMalformed"
//  INVALIDPARAMETERVALUE_TAGKEYNOTFOUND = "InvalidParameterValue.TagKeyNotFound"
//  INVALIDPARAMETERVALUE_UUIDMALFORMED = "InvalidParameterValue.UuidMalformed"
//  INVALIDPARAMETERVALUE_VPCIDMALFORMED = "InvalidParameterValue.VpcIdMalformed"
//  INVALIDSECURITYGROUPID_NOTFOUND = "InvalidSecurityGroupId.NotFound"
//  INVALIDSGID_MALFORMED = "InvalidSgId.Malformed"
//  INVALIDZONE_MISMATCHREGION = "InvalidZone.MismatchRegion"
//  RESOURCENOTFOUND_HPCCLUSTER = "ResourceNotFound.HpcCluster"
//  UNAUTHORIZEDOPERATION_INVALIDTOKEN = "UnauthorizedOperation.InvalidToken"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
func (c *Client) DescribeInstances(request *DescribeInstancesRequest) (response *DescribeInstancesResponse, err error) {
    return c.DescribeInstancesWithContext(context.Background(), request)
}

// DescribeInstances
// This API is used to query the details of instances.
//
// 
//
// * You can filter the query results with the instance `ID`, name, or billing method. See `Filter` for more information.
//
// * If no parameter is defined, a certain number of instances under the current account will be returned. The number is specified by `Limit` and is 20 by default.
//
// error code that may be returned:
//  FAILEDOPERATION_ILLEGALTAGKEY = "FailedOperation.IllegalTagKey"
//  FAILEDOPERATION_ILLEGALTAGVALUE = "FailedOperation.IllegalTagValue"
//  FAILEDOPERATION_TAGKEYRESERVED = "FailedOperation.TagKeyReserved"
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDFILTER = "InvalidFilter"
//  INVALIDFILTERVALUE_LIMITEXCEEDED = "InvalidFilterValue.LimitExceeded"
//  INVALIDHOSTID_MALFORMED = "InvalidHostId.Malformed"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_IPADDRESSMALFORMED = "InvalidParameterValue.IPAddressMalformed"
//  INVALIDPARAMETERVALUE_IPV6ADDRESSMALFORMED = "InvalidParameterValue.IPv6AddressMalformed"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_INVALIDIPFORMAT = "InvalidParameterValue.InvalidIpFormat"
//  INVALIDPARAMETERVALUE_INVALIDTIMEFORMAT = "InvalidParameterValue.InvalidTimeFormat"
//  INVALIDPARAMETERVALUE_INVALIDVAGUENAME = "InvalidParameterValue.InvalidVagueName"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_SUBNETIDMALFORMED = "InvalidParameterValue.SubnetIdMalformed"
//  INVALIDPARAMETERVALUE_TAGKEYNOTFOUND = "InvalidParameterValue.TagKeyNotFound"
//  INVALIDPARAMETERVALUE_UUIDMALFORMED = "InvalidParameterValue.UuidMalformed"
//  INVALIDPARAMETERVALUE_VPCIDMALFORMED = "InvalidParameterValue.VpcIdMalformed"
//  INVALIDSECURITYGROUPID_NOTFOUND = "InvalidSecurityGroupId.NotFound"
//  INVALIDSGID_MALFORMED = "InvalidSgId.Malformed"
//  INVALIDZONE_MISMATCHREGION = "InvalidZone.MismatchRegion"
//  RESOURCENOTFOUND_HPCCLUSTER = "ResourceNotFound.HpcCluster"
//  UNAUTHORIZEDOPERATION_INVALIDTOKEN = "UnauthorizedOperation.InvalidToken"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
func (c *Client) DescribeInstancesWithContext(ctx context.Context, request *DescribeInstancesRequest) (response *DescribeInstancesResponse, err error) {
    if request == nil {
        request = NewDescribeInstancesRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "DescribeInstances")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeInstances require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeInstancesResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeInstancesAttributesRequest() (request *DescribeInstancesAttributesRequest) {
    request = &DescribeInstancesAttributesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "DescribeInstancesAttributes")
    
    
    return
}

func NewDescribeInstancesAttributesResponse() (response *DescribeInstancesAttributesResponse) {
    response = &DescribeInstancesAttributesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeInstancesAttributes
// This API is used to obtain the attributes of specified instances. Currently, it supports querying the custom data UserData of instances.
//
// error code that may be returned:
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
func (c *Client) DescribeInstancesAttributes(request *DescribeInstancesAttributesRequest) (response *DescribeInstancesAttributesResponse, err error) {
    return c.DescribeInstancesAttributesWithContext(context.Background(), request)
}

// DescribeInstancesAttributes
// This API is used to obtain the attributes of specified instances. Currently, it supports querying the custom data UserData of instances.
//
// error code that may be returned:
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
func (c *Client) DescribeInstancesAttributesWithContext(ctx context.Context, request *DescribeInstancesAttributesRequest) (response *DescribeInstancesAttributesResponse, err error) {
    if request == nil {
        request = NewDescribeInstancesAttributesRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "DescribeInstancesAttributes")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeInstancesAttributes require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeInstancesAttributesResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeInstancesOperationLimitRequest() (request *DescribeInstancesOperationLimitRequest) {
    request = &DescribeInstancesOperationLimitRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "DescribeInstancesOperationLimit")
    
    
    return
}

func NewDescribeInstancesOperationLimitResponse() (response *DescribeInstancesOperationLimitResponse) {
    response = &DescribeInstancesOperationLimitResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeInstancesOperationLimit
// This API is used to query limitations on operations on an instance.
//
// 
//
// * Currently you can use this API to query the maximum number of times you can modify the configuration of an instance.
//
// error code that may be returned:
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
func (c *Client) DescribeInstancesOperationLimit(request *DescribeInstancesOperationLimitRequest) (response *DescribeInstancesOperationLimitResponse, err error) {
    return c.DescribeInstancesOperationLimitWithContext(context.Background(), request)
}

// DescribeInstancesOperationLimit
// This API is used to query limitations on operations on an instance.
//
// 
//
// * Currently you can use this API to query the maximum number of times you can modify the configuration of an instance.
//
// error code that may be returned:
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
func (c *Client) DescribeInstancesOperationLimitWithContext(ctx context.Context, request *DescribeInstancesOperationLimitRequest) (response *DescribeInstancesOperationLimitResponse, err error) {
    if request == nil {
        request = NewDescribeInstancesOperationLimitRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "DescribeInstancesOperationLimit")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeInstancesOperationLimit require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeInstancesOperationLimitResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeInstancesStatusRequest() (request *DescribeInstancesStatusRequest) {
    request = &DescribeInstancesStatusRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "DescribeInstancesStatus")
    
    
    return
}

func NewDescribeInstancesStatusResponse() (response *DescribeInstancesStatusResponse) {
    response = &DescribeInstancesStatusResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeInstancesStatus
// This API is used to query the status of instances.
//
// 
//
// * You can query the status of an instance with its `ID`.
//
// * If no parameter is defined, the status of a certain number of instances under the current account will be returned. The number is specified by `Limit` and is 20 by default.
//
// error code that may be returned:
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDFILTER = "InvalidFilter"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  UNAUTHORIZEDOPERATION_INVALIDTOKEN = "UnauthorizedOperation.InvalidToken"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
func (c *Client) DescribeInstancesStatus(request *DescribeInstancesStatusRequest) (response *DescribeInstancesStatusResponse, err error) {
    return c.DescribeInstancesStatusWithContext(context.Background(), request)
}

// DescribeInstancesStatus
// This API is used to query the status of instances.
//
// 
//
// * You can query the status of an instance with its `ID`.
//
// * If no parameter is defined, the status of a certain number of instances under the current account will be returned. The number is specified by `Limit` and is 20 by default.
//
// error code that may be returned:
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDFILTER = "InvalidFilter"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  UNAUTHORIZEDOPERATION_INVALIDTOKEN = "UnauthorizedOperation.InvalidToken"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
func (c *Client) DescribeInstancesStatusWithContext(ctx context.Context, request *DescribeInstancesStatusRequest) (response *DescribeInstancesStatusResponse, err error) {
    if request == nil {
        request = NewDescribeInstancesStatusRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "DescribeInstancesStatus")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeInstancesStatus require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeInstancesStatusResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeInternetChargeTypeConfigsRequest() (request *DescribeInternetChargeTypeConfigsRequest) {
    request = &DescribeInternetChargeTypeConfigsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "DescribeInternetChargeTypeConfigs")
    
    
    return
}

func NewDescribeInternetChargeTypeConfigsResponse() (response *DescribeInternetChargeTypeConfigsResponse) {
    response = &DescribeInternetChargeTypeConfigsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeInternetChargeTypeConfigs
// This API is used to query the network billing methods.
//
// error code that may be returned:
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDFILTER = "InvalidFilter"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  UNAUTHORIZEDOPERATION_INVALIDTOKEN = "UnauthorizedOperation.InvalidToken"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
func (c *Client) DescribeInternetChargeTypeConfigs(request *DescribeInternetChargeTypeConfigsRequest) (response *DescribeInternetChargeTypeConfigsResponse, err error) {
    return c.DescribeInternetChargeTypeConfigsWithContext(context.Background(), request)
}

// DescribeInternetChargeTypeConfigs
// This API is used to query the network billing methods.
//
// error code that may be returned:
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDFILTER = "InvalidFilter"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  UNAUTHORIZEDOPERATION_INVALIDTOKEN = "UnauthorizedOperation.InvalidToken"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
func (c *Client) DescribeInternetChargeTypeConfigsWithContext(ctx context.Context, request *DescribeInternetChargeTypeConfigsRequest) (response *DescribeInternetChargeTypeConfigsResponse, err error) {
    if request == nil {
        request = NewDescribeInternetChargeTypeConfigsRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "DescribeInternetChargeTypeConfigs")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeInternetChargeTypeConfigs require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeInternetChargeTypeConfigsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeKeyPairsRequest() (request *DescribeKeyPairsRequest) {
    request = &DescribeKeyPairsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "DescribeKeyPairs")
    
    
    return
}

func NewDescribeKeyPairsResponse() (response *DescribeKeyPairsResponse) {
    response = &DescribeKeyPairsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeKeyPairs
// This API is used to query key pairs.
//
// 
//
// * A key pair is a pair of keys generated by an algorithm in which the public key is available to the public and the private key is available only to the user. You can use this API to query the public key but not the private key.
//
// error code that may be returned:
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDFILTER = "InvalidFilter"
//  INVALIDFILTERVALUE_LIMITEXCEEDED = "InvalidFilterValue.LimitExceeded"
//  INVALIDKEYPAIR_LIMITEXCEEDED = "InvalidKeyPair.LimitExceeded"
//  INVALIDKEYPAIRID_MALFORMED = "InvalidKeyPairId.Malformed"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUELIMIT = "InvalidParameterValueLimit"
//  INVALIDPARAMETERVALUEOFFSET = "InvalidParameterValueOffset"
//  UNAUTHORIZEDOPERATION_INVALIDTOKEN = "UnauthorizedOperation.InvalidToken"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
func (c *Client) DescribeKeyPairs(request *DescribeKeyPairsRequest) (response *DescribeKeyPairsResponse, err error) {
    return c.DescribeKeyPairsWithContext(context.Background(), request)
}

// DescribeKeyPairs
// This API is used to query key pairs.
//
// 
//
// * A key pair is a pair of keys generated by an algorithm in which the public key is available to the public and the private key is available only to the user. You can use this API to query the public key but not the private key.
//
// error code that may be returned:
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDFILTER = "InvalidFilter"
//  INVALIDFILTERVALUE_LIMITEXCEEDED = "InvalidFilterValue.LimitExceeded"
//  INVALIDKEYPAIR_LIMITEXCEEDED = "InvalidKeyPair.LimitExceeded"
//  INVALIDKEYPAIRID_MALFORMED = "InvalidKeyPairId.Malformed"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUELIMIT = "InvalidParameterValueLimit"
//  INVALIDPARAMETERVALUEOFFSET = "InvalidParameterValueOffset"
//  UNAUTHORIZEDOPERATION_INVALIDTOKEN = "UnauthorizedOperation.InvalidToken"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
func (c *Client) DescribeKeyPairsWithContext(ctx context.Context, request *DescribeKeyPairsRequest) (response *DescribeKeyPairsResponse, err error) {
    if request == nil {
        request = NewDescribeKeyPairsRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "DescribeKeyPairs")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeKeyPairs require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeKeyPairsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeLaunchTemplateVersionsRequest() (request *DescribeLaunchTemplateVersionsRequest) {
    request = &DescribeLaunchTemplateVersionsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "DescribeLaunchTemplateVersions")
    
    
    return
}

func NewDescribeLaunchTemplateVersionsResponse() (response *DescribeLaunchTemplateVersionsResponse) {
    response = &DescribeLaunchTemplateVersionsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeLaunchTemplateVersions
// This API is used to query the information of instance launch template versions.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERCOMBINATION = "InvalidParameterCombination"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDMALFORMED = "InvalidParameterValue.LaunchTemplateIdMalformed"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDNOTEXISTED = "InvalidParameterValue.LaunchTemplateIdNotExisted"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDVERNOTEXISTED = "InvalidParameterValue.LaunchTemplateIdVerNotExisted"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATENOTFOUND = "InvalidParameterValue.LaunchTemplateNotFound"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEVERSION = "InvalidParameterValue.LaunchTemplateVersion"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_NOTSUPPORTED = "InvalidParameterValue.NotSupported"
//  MISSINGPARAMETER = "MissingParameter"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
//  UNAUTHORIZEDOPERATION_MFAEXPIRED = "UnauthorizedOperation.MFAExpired"
//  UNAUTHORIZEDOPERATION_MFANOTFOUND = "UnauthorizedOperation.MFANotFound"
//  UNKNOWNPARAMETER = "UnknownParameter"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
func (c *Client) DescribeLaunchTemplateVersions(request *DescribeLaunchTemplateVersionsRequest) (response *DescribeLaunchTemplateVersionsResponse, err error) {
    return c.DescribeLaunchTemplateVersionsWithContext(context.Background(), request)
}

// DescribeLaunchTemplateVersions
// This API is used to query the information of instance launch template versions.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERCOMBINATION = "InvalidParameterCombination"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDMALFORMED = "InvalidParameterValue.LaunchTemplateIdMalformed"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDNOTEXISTED = "InvalidParameterValue.LaunchTemplateIdNotExisted"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDVERNOTEXISTED = "InvalidParameterValue.LaunchTemplateIdVerNotExisted"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATENOTFOUND = "InvalidParameterValue.LaunchTemplateNotFound"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEVERSION = "InvalidParameterValue.LaunchTemplateVersion"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_NOTSUPPORTED = "InvalidParameterValue.NotSupported"
//  MISSINGPARAMETER = "MissingParameter"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
//  UNAUTHORIZEDOPERATION_MFAEXPIRED = "UnauthorizedOperation.MFAExpired"
//  UNAUTHORIZEDOPERATION_MFANOTFOUND = "UnauthorizedOperation.MFANotFound"
//  UNKNOWNPARAMETER = "UnknownParameter"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
func (c *Client) DescribeLaunchTemplateVersionsWithContext(ctx context.Context, request *DescribeLaunchTemplateVersionsRequest) (response *DescribeLaunchTemplateVersionsResponse, err error) {
    if request == nil {
        request = NewDescribeLaunchTemplateVersionsRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "DescribeLaunchTemplateVersions")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeLaunchTemplateVersions require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeLaunchTemplateVersionsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeLaunchTemplatesRequest() (request *DescribeLaunchTemplatesRequest) {
    request = &DescribeLaunchTemplatesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "DescribeLaunchTemplates")
    
    
    return
}

func NewDescribeLaunchTemplatesResponse() (response *DescribeLaunchTemplatesResponse) {
    response = &DescribeLaunchTemplatesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeLaunchTemplates
// This API is used to query one or more instance launch templates.
//
// error code that may be returned:
//  INVALIDFILTER = "InvalidFilter"
//  INVALIDFILTERVALUE_LIMITEXCEEDED = "InvalidFilterValue.LimitExceeded"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_INVALIDLAUNCHTEMPLATENAME = "InvalidParameterValue.InvalidLaunchTemplateName"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDMALFORMED = "InvalidParameterValue.LaunchTemplateIdMalformed"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDNOTEXISTED = "InvalidParameterValue.LaunchTemplateIdNotExisted"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDVERNOTEXISTED = "InvalidParameterValue.LaunchTemplateIdVerNotExisted"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATENOTFOUND = "InvalidParameterValue.LaunchTemplateNotFound"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEVERSION = "InvalidParameterValue.LaunchTemplateVersion"
//  UNKNOWNPARAMETER = "UnknownParameter"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
func (c *Client) DescribeLaunchTemplates(request *DescribeLaunchTemplatesRequest) (response *DescribeLaunchTemplatesResponse, err error) {
    return c.DescribeLaunchTemplatesWithContext(context.Background(), request)
}

// DescribeLaunchTemplates
// This API is used to query one or more instance launch templates.
//
// error code that may be returned:
//  INVALIDFILTER = "InvalidFilter"
//  INVALIDFILTERVALUE_LIMITEXCEEDED = "InvalidFilterValue.LimitExceeded"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_INVALIDLAUNCHTEMPLATENAME = "InvalidParameterValue.InvalidLaunchTemplateName"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDMALFORMED = "InvalidParameterValue.LaunchTemplateIdMalformed"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDNOTEXISTED = "InvalidParameterValue.LaunchTemplateIdNotExisted"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDVERNOTEXISTED = "InvalidParameterValue.LaunchTemplateIdVerNotExisted"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATENOTFOUND = "InvalidParameterValue.LaunchTemplateNotFound"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEVERSION = "InvalidParameterValue.LaunchTemplateVersion"
//  UNKNOWNPARAMETER = "UnknownParameter"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
func (c *Client) DescribeLaunchTemplatesWithContext(ctx context.Context, request *DescribeLaunchTemplatesRequest) (response *DescribeLaunchTemplatesResponse, err error) {
    if request == nil {
        request = NewDescribeLaunchTemplatesRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "DescribeLaunchTemplates")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeLaunchTemplates require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeLaunchTemplatesResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeRegionsRequest() (request *DescribeRegionsRequest) {
    request = &DescribeRegionsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "DescribeRegions")
    
    
    return
}

func NewDescribeRegionsResponse() (response *DescribeRegionsResponse) {
    response = &DescribeRegionsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeRegions
// This API is suspended. To query the information of regions, use [DescribeZones](https://intl.cloud.tencent.com/document/product/1596/77930?from_cn_redirect=1).
//
// error code that may be returned:
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
func (c *Client) DescribeRegions(request *DescribeRegionsRequest) (response *DescribeRegionsResponse, err error) {
    return c.DescribeRegionsWithContext(context.Background(), request)
}

// DescribeRegions
// This API is suspended. To query the information of regions, use [DescribeZones](https://intl.cloud.tencent.com/document/product/1596/77930?from_cn_redirect=1).
//
// error code that may be returned:
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
func (c *Client) DescribeRegionsWithContext(ctx context.Context, request *DescribeRegionsRequest) (response *DescribeRegionsResponse, err error) {
    if request == nil {
        request = NewDescribeRegionsRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "DescribeRegions")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeRegions require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeRegionsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeReservedInstancesRequest() (request *DescribeReservedInstancesRequest) {
    request = &DescribeReservedInstancesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "DescribeReservedInstances")
    
    
    return
}

func NewDescribeReservedInstancesResponse() (response *DescribeReservedInstancesResponse) {
    response = &DescribeReservedInstancesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeReservedInstances
// This API is used to list the reserved instances purchased by the user.
//
// error code that may be returned:
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDFILTER = "InvalidFilter"
//  INVALIDFILTERVALUE_LIMITEXCEEDED = "InvalidFilterValue.LimitExceeded"
//  INVALIDINSTANCETYPE_MALFORMED = "InvalidInstanceType.Malformed"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  UNSUPPORTEDOPERATION_INVALIDPERMISSIONNONINTERNATIONALACCOUNT = "UnsupportedOperation.InvalidPermissionNonInternationalAccount"
//  UNSUPPORTEDOPERATION_RESERVEDINSTANCEINVISIBLEFORUSER = "UnsupportedOperation.ReservedInstanceInvisibleForUser"
func (c *Client) DescribeReservedInstances(request *DescribeReservedInstancesRequest) (response *DescribeReservedInstancesResponse, err error) {
    return c.DescribeReservedInstancesWithContext(context.Background(), request)
}

// DescribeReservedInstances
// This API is used to list the reserved instances purchased by the user.
//
// error code that may be returned:
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDFILTER = "InvalidFilter"
//  INVALIDFILTERVALUE_LIMITEXCEEDED = "InvalidFilterValue.LimitExceeded"
//  INVALIDINSTANCETYPE_MALFORMED = "InvalidInstanceType.Malformed"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  UNSUPPORTEDOPERATION_INVALIDPERMISSIONNONINTERNATIONALACCOUNT = "UnsupportedOperation.InvalidPermissionNonInternationalAccount"
//  UNSUPPORTEDOPERATION_RESERVEDINSTANCEINVISIBLEFORUSER = "UnsupportedOperation.ReservedInstanceInvisibleForUser"
func (c *Client) DescribeReservedInstancesWithContext(ctx context.Context, request *DescribeReservedInstancesRequest) (response *DescribeReservedInstancesResponse, err error) {
    if request == nil {
        request = NewDescribeReservedInstancesRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "DescribeReservedInstances")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeReservedInstances require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeReservedInstancesResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeReservedInstancesConfigInfosRequest() (request *DescribeReservedInstancesConfigInfosRequest) {
    request = &DescribeReservedInstancesConfigInfosRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "DescribeReservedInstancesConfigInfos")
    
    
    return
}

func NewDescribeReservedInstancesConfigInfosResponse() (response *DescribeReservedInstancesConfigInfosResponse) {
    response = &DescribeReservedInstancesConfigInfosResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeReservedInstancesConfigInfos
// This API is used to describe reserved instance (RI) offerings.
//
// error code that may be returned:
//  INVALIDFILTER = "InvalidFilter"
//  INVALIDZONE_MISMATCHREGION = "InvalidZone.MismatchRegion"
//  UNSUPPORTEDOPERATION_INVALIDPERMISSIONNONINTERNATIONALACCOUNT = "UnsupportedOperation.InvalidPermissionNonInternationalAccount"
//  UNSUPPORTEDOPERATION_RESERVEDINSTANCEINVISIBLEFORUSER = "UnsupportedOperation.ReservedInstanceInvisibleForUser"
func (c *Client) DescribeReservedInstancesConfigInfos(request *DescribeReservedInstancesConfigInfosRequest) (response *DescribeReservedInstancesConfigInfosResponse, err error) {
    return c.DescribeReservedInstancesConfigInfosWithContext(context.Background(), request)
}

// DescribeReservedInstancesConfigInfos
// This API is used to describe reserved instance (RI) offerings.
//
// error code that may be returned:
//  INVALIDFILTER = "InvalidFilter"
//  INVALIDZONE_MISMATCHREGION = "InvalidZone.MismatchRegion"
//  UNSUPPORTEDOPERATION_INVALIDPERMISSIONNONINTERNATIONALACCOUNT = "UnsupportedOperation.InvalidPermissionNonInternationalAccount"
//  UNSUPPORTEDOPERATION_RESERVEDINSTANCEINVISIBLEFORUSER = "UnsupportedOperation.ReservedInstanceInvisibleForUser"
func (c *Client) DescribeReservedInstancesConfigInfosWithContext(ctx context.Context, request *DescribeReservedInstancesConfigInfosRequest) (response *DescribeReservedInstancesConfigInfosResponse, err error) {
    if request == nil {
        request = NewDescribeReservedInstancesConfigInfosRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "DescribeReservedInstancesConfigInfos")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeReservedInstancesConfigInfos require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeReservedInstancesConfigInfosResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeReservedInstancesOfferingsRequest() (request *DescribeReservedInstancesOfferingsRequest) {
    request = &DescribeReservedInstancesOfferingsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "DescribeReservedInstancesOfferings")
    
    
    return
}

func NewDescribeReservedInstancesOfferingsResponse() (response *DescribeReservedInstancesOfferingsResponse) {
    response = &DescribeReservedInstancesOfferingsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeReservedInstancesOfferings
// This API is used to describe Reserved Instance offerings that are available for purchase.
//
// error code that may be returned:
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDFILTER = "InvalidFilter"
//  INVALIDFILTERVALUE_LIMITEXCEEDED = "InvalidFilterValue.LimitExceeded"
//  INVALIDINSTANCETYPE_MALFORMED = "InvalidInstanceType.Malformed"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  UNSUPPORTEDOPERATION_INVALIDPERMISSIONNONINTERNATIONALACCOUNT = "UnsupportedOperation.InvalidPermissionNonInternationalAccount"
//  UNSUPPORTEDOPERATION_RESERVEDINSTANCEINVISIBLEFORUSER = "UnsupportedOperation.ReservedInstanceInvisibleForUser"
func (c *Client) DescribeReservedInstancesOfferings(request *DescribeReservedInstancesOfferingsRequest) (response *DescribeReservedInstancesOfferingsResponse, err error) {
    return c.DescribeReservedInstancesOfferingsWithContext(context.Background(), request)
}

// DescribeReservedInstancesOfferings
// This API is used to describe Reserved Instance offerings that are available for purchase.
//
// error code that may be returned:
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDFILTER = "InvalidFilter"
//  INVALIDFILTERVALUE_LIMITEXCEEDED = "InvalidFilterValue.LimitExceeded"
//  INVALIDINSTANCETYPE_MALFORMED = "InvalidInstanceType.Malformed"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  UNSUPPORTEDOPERATION_INVALIDPERMISSIONNONINTERNATIONALACCOUNT = "UnsupportedOperation.InvalidPermissionNonInternationalAccount"
//  UNSUPPORTEDOPERATION_RESERVEDINSTANCEINVISIBLEFORUSER = "UnsupportedOperation.ReservedInstanceInvisibleForUser"
func (c *Client) DescribeReservedInstancesOfferingsWithContext(ctx context.Context, request *DescribeReservedInstancesOfferingsRequest) (response *DescribeReservedInstancesOfferingsResponse, err error) {
    if request == nil {
        request = NewDescribeReservedInstancesOfferingsRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "DescribeReservedInstancesOfferings")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeReservedInstancesOfferings require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeReservedInstancesOfferingsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeZoneInstanceConfigInfosRequest() (request *DescribeZoneInstanceConfigInfosRequest) {
    request = &DescribeZoneInstanceConfigInfosRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "DescribeZoneInstanceConfigInfos")
    
    
    return
}

func NewDescribeZoneInstanceConfigInfosResponse() (response *DescribeZoneInstanceConfigInfosResponse) {
    response = &DescribeZoneInstanceConfigInfosResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeZoneInstanceConfigInfos
// This API is used to query the configurations of models in an availability zone.
//
// error code that may be returned:
//  INVALIDFILTER = "InvalidFilter"
//  INVALIDFILTERVALUE_LIMITEXCEEDED = "InvalidFilterValue.LimitExceeded"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDINSTANCETYPE_MALFORMED = "InvalidInstanceType.Malformed"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_ZONENOTSUPPORTED = "InvalidParameterValue.ZoneNotSupported"
//  INVALIDREGION_NOTFOUND = "InvalidRegion.NotFound"
//  INVALIDZONE_MISMATCHREGION = "InvalidZone.MismatchRegion"
//  RESOURCEINSUFFICIENT_AVAILABILITYZONESOLDOUT = "ResourceInsufficient.AvailabilityZoneSoldOut"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
//  UNSUPPORTEDOPERATION_UNSUPPORTEDPOOL = "UnsupportedOperation.UnsupportedPool"
func (c *Client) DescribeZoneInstanceConfigInfos(request *DescribeZoneInstanceConfigInfosRequest) (response *DescribeZoneInstanceConfigInfosResponse, err error) {
    return c.DescribeZoneInstanceConfigInfosWithContext(context.Background(), request)
}

// DescribeZoneInstanceConfigInfos
// This API is used to query the configurations of models in an availability zone.
//
// error code that may be returned:
//  INVALIDFILTER = "InvalidFilter"
//  INVALIDFILTERVALUE_LIMITEXCEEDED = "InvalidFilterValue.LimitExceeded"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDINSTANCETYPE_MALFORMED = "InvalidInstanceType.Malformed"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_ZONENOTSUPPORTED = "InvalidParameterValue.ZoneNotSupported"
//  INVALIDREGION_NOTFOUND = "InvalidRegion.NotFound"
//  INVALIDZONE_MISMATCHREGION = "InvalidZone.MismatchRegion"
//  RESOURCEINSUFFICIENT_AVAILABILITYZONESOLDOUT = "ResourceInsufficient.AvailabilityZoneSoldOut"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
//  UNSUPPORTEDOPERATION_UNSUPPORTEDPOOL = "UnsupportedOperation.UnsupportedPool"
func (c *Client) DescribeZoneInstanceConfigInfosWithContext(ctx context.Context, request *DescribeZoneInstanceConfigInfosRequest) (response *DescribeZoneInstanceConfigInfosResponse, err error) {
    if request == nil {
        request = NewDescribeZoneInstanceConfigInfosRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "DescribeZoneInstanceConfigInfos")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeZoneInstanceConfigInfos require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeZoneInstanceConfigInfosResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeZonesRequest() (request *DescribeZonesRequest) {
    request = &DescribeZonesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "DescribeZones")
    
    
    return
}

func NewDescribeZonesResponse() (response *DescribeZonesResponse) {
    response = &DescribeZonesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeZones
// This API is used to query availability zones.
//
// error code that may be returned:
//  INVALIDFILTER = "InvalidFilter"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_ZONENOTSUPPORTED = "InvalidParameterValue.ZoneNotSupported"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
func (c *Client) DescribeZones(request *DescribeZonesRequest) (response *DescribeZonesResponse, err error) {
    return c.DescribeZonesWithContext(context.Background(), request)
}

// DescribeZones
// This API is used to query availability zones.
//
// error code that may be returned:
//  INVALIDFILTER = "InvalidFilter"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_ZONENOTSUPPORTED = "InvalidParameterValue.ZoneNotSupported"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
func (c *Client) DescribeZonesWithContext(ctx context.Context, request *DescribeZonesRequest) (response *DescribeZonesResponse, err error) {
    if request == nil {
        request = NewDescribeZonesRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "DescribeZones")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeZones require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeZonesResponse()
    err = c.Send(request, response)
    return
}

func NewDisassociateInstancesKeyPairsRequest() (request *DisassociateInstancesKeyPairsRequest) {
    request = &DisassociateInstancesKeyPairsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "DisassociateInstancesKeyPairs")
    
    
    return
}

func NewDisassociateInstancesKeyPairsResponse() (response *DisassociateInstancesKeyPairsResponse) {
    response = &DisassociateInstancesKeyPairsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DisassociateInstancesKeyPairs
// This API is used to unbind one or more key pairs from one or more instances.
//
// 
//
// * It only supports [`STOPPED`](https://intl.cloud.tencent.com/document/product/213/15753?from_cn_redirect=1#InstanceStatus) Linux instances.
//
// * After a key pair is disassociated from an instance, you can log in to the instance with password.
//
// * If you did not set a password for the instance, you will not be able to log in via SSH after the unbinding. In this case, you can call [ResetInstancesPassword](https://intl.cloud.tencent.com/document/api/213/15736?from_cn_redirect=1) to set a login password.
//
// * Batch operations are supported. The maximum number of instances in each request is 100. If instances not available for the operation are selected, you will get an error code.
//
// error code that may be returned:
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDKEYPAIRID_MALFORMED = "InvalidKeyPairId.Malformed"
//  INVALIDKEYPAIRID_NOTFOUND = "InvalidKeyPairId.NotFound"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  MISSINGPARAMETER = "MissingParameter"
//  MUTEXOPERATION_TASKRUNNING = "MutexOperation.TaskRunning"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  UNSUPPORTEDOPERATION_INSTANCEOSWINDOWS = "UnsupportedOperation.InstanceOsWindows"
//  UNSUPPORTEDOPERATION_INSTANCESTATEBANNING = "UnsupportedOperation.InstanceStateBanning"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERRESCUEMODE = "UnsupportedOperation.InstanceStateEnterRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEFREEZING = "UnsupportedOperation.InstanceStateFreezing"
//  UNSUPPORTEDOPERATION_INSTANCESTATEPENDING = "UnsupportedOperation.InstanceStatePending"
//  UNSUPPORTEDOPERATION_INSTANCESTATEREBOOTING = "UnsupportedOperation.InstanceStateRebooting"
//  UNSUPPORTEDOPERATION_INSTANCESTATERESCUEMODE = "UnsupportedOperation.InstanceStateRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATERUNNING = "UnsupportedOperation.InstanceStateRunning"
//  UNSUPPORTEDOPERATION_INSTANCESTATESHUTDOWN = "UnsupportedOperation.InstanceStateShutdown"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTARTING = "UnsupportedOperation.InstanceStateStarting"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPING = "UnsupportedOperation.InstanceStateStopping"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATING = "UnsupportedOperation.InstanceStateTerminating"
//  UNSUPPORTEDOPERATION_STOPPEDMODESTOPCHARGING = "UnsupportedOperation.StoppedModeStopCharging"
func (c *Client) DisassociateInstancesKeyPairs(request *DisassociateInstancesKeyPairsRequest) (response *DisassociateInstancesKeyPairsResponse, err error) {
    return c.DisassociateInstancesKeyPairsWithContext(context.Background(), request)
}

// DisassociateInstancesKeyPairs
// This API is used to unbind one or more key pairs from one or more instances.
//
// 
//
// * It only supports [`STOPPED`](https://intl.cloud.tencent.com/document/product/213/15753?from_cn_redirect=1#InstanceStatus) Linux instances.
//
// * After a key pair is disassociated from an instance, you can log in to the instance with password.
//
// * If you did not set a password for the instance, you will not be able to log in via SSH after the unbinding. In this case, you can call [ResetInstancesPassword](https://intl.cloud.tencent.com/document/api/213/15736?from_cn_redirect=1) to set a login password.
//
// * Batch operations are supported. The maximum number of instances in each request is 100. If instances not available for the operation are selected, you will get an error code.
//
// error code that may be returned:
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDKEYPAIRID_MALFORMED = "InvalidKeyPairId.Malformed"
//  INVALIDKEYPAIRID_NOTFOUND = "InvalidKeyPairId.NotFound"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  MISSINGPARAMETER = "MissingParameter"
//  MUTEXOPERATION_TASKRUNNING = "MutexOperation.TaskRunning"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  UNSUPPORTEDOPERATION_INSTANCEOSWINDOWS = "UnsupportedOperation.InstanceOsWindows"
//  UNSUPPORTEDOPERATION_INSTANCESTATEBANNING = "UnsupportedOperation.InstanceStateBanning"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERRESCUEMODE = "UnsupportedOperation.InstanceStateEnterRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEFREEZING = "UnsupportedOperation.InstanceStateFreezing"
//  UNSUPPORTEDOPERATION_INSTANCESTATEPENDING = "UnsupportedOperation.InstanceStatePending"
//  UNSUPPORTEDOPERATION_INSTANCESTATEREBOOTING = "UnsupportedOperation.InstanceStateRebooting"
//  UNSUPPORTEDOPERATION_INSTANCESTATERESCUEMODE = "UnsupportedOperation.InstanceStateRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATERUNNING = "UnsupportedOperation.InstanceStateRunning"
//  UNSUPPORTEDOPERATION_INSTANCESTATESHUTDOWN = "UnsupportedOperation.InstanceStateShutdown"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTARTING = "UnsupportedOperation.InstanceStateStarting"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPING = "UnsupportedOperation.InstanceStateStopping"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATING = "UnsupportedOperation.InstanceStateTerminating"
//  UNSUPPORTEDOPERATION_STOPPEDMODESTOPCHARGING = "UnsupportedOperation.StoppedModeStopCharging"
func (c *Client) DisassociateInstancesKeyPairsWithContext(ctx context.Context, request *DisassociateInstancesKeyPairsRequest) (response *DisassociateInstancesKeyPairsResponse, err error) {
    if request == nil {
        request = NewDisassociateInstancesKeyPairsRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "DisassociateInstancesKeyPairs")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DisassociateInstancesKeyPairs require credential")
    }

    request.SetContext(ctx)
    
    response = NewDisassociateInstancesKeyPairsResponse()
    err = c.Send(request, response)
    return
}

func NewDisassociateSecurityGroupsRequest() (request *DisassociateSecurityGroupsRequest) {
    request = &DisassociateSecurityGroupsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "DisassociateSecurityGroups")
    
    
    return
}

func NewDisassociateSecurityGroupsResponse() (response *DisassociateSecurityGroupsResponse) {
    response = &DisassociateSecurityGroupsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DisassociateSecurityGroups
// This API is used to disassociate security groups from instances.
//
// error code that may be returned:
//  FAILEDOPERATION_SECURITYGROUPACTIONFAILED = "FailedOperation.SecurityGroupActionFailed"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDSECURITYGROUPID_NOTFOUND = "InvalidSecurityGroupId.NotFound"
//  INVALIDSGID_MALFORMED = "InvalidSgId.Malformed"
//  MUTEXOPERATION_TASKRUNNING = "MutexOperation.TaskRunning"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  SECGROUPACTIONFAILURE = "SecGroupActionFailure"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
func (c *Client) DisassociateSecurityGroups(request *DisassociateSecurityGroupsRequest) (response *DisassociateSecurityGroupsResponse, err error) {
    return c.DisassociateSecurityGroupsWithContext(context.Background(), request)
}

// DisassociateSecurityGroups
// This API is used to disassociate security groups from instances.
//
// error code that may be returned:
//  FAILEDOPERATION_SECURITYGROUPACTIONFAILED = "FailedOperation.SecurityGroupActionFailed"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDSECURITYGROUPID_NOTFOUND = "InvalidSecurityGroupId.NotFound"
//  INVALIDSGID_MALFORMED = "InvalidSgId.Malformed"
//  MUTEXOPERATION_TASKRUNNING = "MutexOperation.TaskRunning"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  SECGROUPACTIONFAILURE = "SecGroupActionFailure"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
func (c *Client) DisassociateSecurityGroupsWithContext(ctx context.Context, request *DisassociateSecurityGroupsRequest) (response *DisassociateSecurityGroupsResponse, err error) {
    if request == nil {
        request = NewDisassociateSecurityGroupsRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "DisassociateSecurityGroups")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DisassociateSecurityGroups require credential")
    }

    request.SetContext(ctx)
    
    response = NewDisassociateSecurityGroupsResponse()
    err = c.Send(request, response)
    return
}

func NewEnterRescueModeRequest() (request *EnterRescueModeRequest) {
    request = &EnterRescueModeRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "EnterRescueMode")
    
    
    return
}

func NewEnterRescueModeResponse() (response *EnterRescueModeResponse) {
    response = &EnterRescueModeResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// EnterRescueMode
// This API is used to enter the rescue mode.
//
// error code that may be returned:
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_INVALIDPASSWORD = "InvalidParameterValue.InvalidPassword"
//  INVALIDPASSWORD = "InvalidPassword"
//  MUTEXOPERATION_TASKRUNNING = "MutexOperation.TaskRunning"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  UNSUPPORTEDOPERATION_ARMARCHITECTURE = "UnsupportedOperation.ArmArchitecture"
//  UNSUPPORTEDOPERATION_EDGEZONENOTSUPPORTCLOUDDISK = "UnsupportedOperation.EdgeZoneNotSupportCloudDisk"
//  UNSUPPORTEDOPERATION_INSTANCEREINSTALLFAILED = "UnsupportedOperation.InstanceReinstallFailed"
//  UNSUPPORTEDOPERATION_INSTANCESTATEBANNING = "UnsupportedOperation.InstanceStateBanning"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERRESCUEMODE = "UnsupportedOperation.InstanceStateEnterRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERSERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateEnterServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITRESCUEMODE = "UnsupportedOperation.InstanceStateExitRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITSERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateExitServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATEFREEZING = "UnsupportedOperation.InstanceStateFreezing"
//  UNSUPPORTEDOPERATION_INSTANCESTATEREBOOTING = "UnsupportedOperation.InstanceStateRebooting"
//  UNSUPPORTEDOPERATION_INSTANCESTATERESCUEMODE = "UnsupportedOperation.InstanceStateRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATESERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATESHUTDOWN = "UnsupportedOperation.InstanceStateShutdown"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPING = "UnsupportedOperation.InstanceStateStopping"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATING = "UnsupportedOperation.InstanceStateTerminating"
//  UNSUPPORTEDOPERATION_SPECIALINSTANCETYPE = "UnsupportedOperation.SpecialInstanceType"
//  UNSUPPORTEDOPERATION_STOPPEDMODESTOPCHARGING = "UnsupportedOperation.StoppedModeStopCharging"
func (c *Client) EnterRescueMode(request *EnterRescueModeRequest) (response *EnterRescueModeResponse, err error) {
    return c.EnterRescueModeWithContext(context.Background(), request)
}

// EnterRescueMode
// This API is used to enter the rescue mode.
//
// error code that may be returned:
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_INVALIDPASSWORD = "InvalidParameterValue.InvalidPassword"
//  INVALIDPASSWORD = "InvalidPassword"
//  MUTEXOPERATION_TASKRUNNING = "MutexOperation.TaskRunning"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  UNSUPPORTEDOPERATION_ARMARCHITECTURE = "UnsupportedOperation.ArmArchitecture"
//  UNSUPPORTEDOPERATION_EDGEZONENOTSUPPORTCLOUDDISK = "UnsupportedOperation.EdgeZoneNotSupportCloudDisk"
//  UNSUPPORTEDOPERATION_INSTANCEREINSTALLFAILED = "UnsupportedOperation.InstanceReinstallFailed"
//  UNSUPPORTEDOPERATION_INSTANCESTATEBANNING = "UnsupportedOperation.InstanceStateBanning"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERRESCUEMODE = "UnsupportedOperation.InstanceStateEnterRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERSERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateEnterServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITRESCUEMODE = "UnsupportedOperation.InstanceStateExitRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITSERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateExitServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATEFREEZING = "UnsupportedOperation.InstanceStateFreezing"
//  UNSUPPORTEDOPERATION_INSTANCESTATEREBOOTING = "UnsupportedOperation.InstanceStateRebooting"
//  UNSUPPORTEDOPERATION_INSTANCESTATERESCUEMODE = "UnsupportedOperation.InstanceStateRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATESERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATESHUTDOWN = "UnsupportedOperation.InstanceStateShutdown"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPING = "UnsupportedOperation.InstanceStateStopping"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATING = "UnsupportedOperation.InstanceStateTerminating"
//  UNSUPPORTEDOPERATION_SPECIALINSTANCETYPE = "UnsupportedOperation.SpecialInstanceType"
//  UNSUPPORTEDOPERATION_STOPPEDMODESTOPCHARGING = "UnsupportedOperation.StoppedModeStopCharging"
func (c *Client) EnterRescueModeWithContext(ctx context.Context, request *EnterRescueModeRequest) (response *EnterRescueModeResponse, err error) {
    if request == nil {
        request = NewEnterRescueModeRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "EnterRescueMode")
    
    if c.GetCredential() == nil {
        return nil, errors.New("EnterRescueMode require credential")
    }

    request.SetContext(ctx)
    
    response = NewEnterRescueModeResponse()
    err = c.Send(request, response)
    return
}

func NewExitRescueModeRequest() (request *ExitRescueModeRequest) {
    request = &ExitRescueModeRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "ExitRescueMode")
    
    
    return
}

func NewExitRescueModeResponse() (response *ExitRescueModeResponse) {
    response = &ExitRescueModeResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// ExitRescueMode
// This API is used to exit the rescue mode.
//
// error code that may be returned:
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  MUTEXOPERATION_TASKRUNNING = "MutexOperation.TaskRunning"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERRESCUEMODE = "UnsupportedOperation.InstanceStateEnterRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITRESCUEMODE = "UnsupportedOperation.InstanceStateExitRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEREBOOTING = "UnsupportedOperation.InstanceStateRebooting"
//  UNSUPPORTEDOPERATION_INSTANCESTATERUNNING = "UnsupportedOperation.InstanceStateRunning"
//  UNSUPPORTEDOPERATION_INSTANCESTATESERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATESHUTDOWN = "UnsupportedOperation.InstanceStateShutdown"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPED = "UnsupportedOperation.InstanceStateStopped"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPING = "UnsupportedOperation.InstanceStateStopping"
func (c *Client) ExitRescueMode(request *ExitRescueModeRequest) (response *ExitRescueModeResponse, err error) {
    return c.ExitRescueModeWithContext(context.Background(), request)
}

// ExitRescueMode
// This API is used to exit the rescue mode.
//
// error code that may be returned:
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  MUTEXOPERATION_TASKRUNNING = "MutexOperation.TaskRunning"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERRESCUEMODE = "UnsupportedOperation.InstanceStateEnterRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITRESCUEMODE = "UnsupportedOperation.InstanceStateExitRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEREBOOTING = "UnsupportedOperation.InstanceStateRebooting"
//  UNSUPPORTEDOPERATION_INSTANCESTATERUNNING = "UnsupportedOperation.InstanceStateRunning"
//  UNSUPPORTEDOPERATION_INSTANCESTATESERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATESHUTDOWN = "UnsupportedOperation.InstanceStateShutdown"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPED = "UnsupportedOperation.InstanceStateStopped"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPING = "UnsupportedOperation.InstanceStateStopping"
func (c *Client) ExitRescueModeWithContext(ctx context.Context, request *ExitRescueModeRequest) (response *ExitRescueModeResponse, err error) {
    if request == nil {
        request = NewExitRescueModeRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "ExitRescueMode")
    
    if c.GetCredential() == nil {
        return nil, errors.New("ExitRescueMode require credential")
    }

    request.SetContext(ctx)
    
    response = NewExitRescueModeResponse()
    err = c.Send(request, response)
    return
}

func NewExportImagesRequest() (request *ExportImagesRequest) {
    request = &ExportImagesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "ExportImages")
    
    
    return
}

func NewExportImagesResponse() (response *ExportImagesResponse) {
    response = &ExportImagesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// ExportImages
// This API is used to export custom images to the specified COS bucket.
//
// error code that may be returned:
//  AUTHFAILURE_CAMROLENAMEAUTHENTICATEFAILED = "AuthFailure.CamRoleNameAuthenticateFailed"
//  INVALIDIMAGEID_MALFORMED = "InvalidImageId.Malformed"
//  INVALIDIMAGEID_NOTFOUND = "InvalidImageId.NotFound"
//  INVALIDPARAMETER_IMAGEIDSSNAPSHOTIDSMUSTONE = "InvalidParameter.ImageIdsSnapshotIdsMustOne"
//  INVALIDPARAMETER_SNAPSHOTNOTFOUND = "InvalidParameter.SnapshotNotFound"
//  INVALIDPARAMETERVALUE_BUCKETNOTFOUND = "InvalidParameterValue.BucketNotFound"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_INVALIDBUCKETPERMISSIONFOREXPORT = "InvalidParameterValue.InvalidBucketPermissionForExport"
//  INVALIDPARAMETERVALUE_INVALIDFILENAMEPREFIXLIST = "InvalidParameterValue.InvalidFileNamePrefixList"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEID = "InvalidParameterValue.InvalidImageId"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEOSNAME = "InvalidParameterValue.InvalidImageOsName"
//  INVALIDPARAMETERVALUE_INVALIDIMAGESTATE = "InvalidParameterValue.InvalidImageState"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_TOOLARGE = "InvalidParameterValue.TooLarge"
//  LIMITEXCEEDED_EXPORTIMAGETASKLIMITEXCEEDED = "LimitExceeded.ExportImageTaskLimitExceeded"
//  UNSUPPORTEDOPERATION_ENCRYPTEDIMAGESNOTSUPPORTED = "UnsupportedOperation.EncryptedImagesNotSupported"
//  UNSUPPORTEDOPERATION_IMAGETOOLARGEEXPORTUNSUPPORTED = "UnsupportedOperation.ImageTooLargeExportUnsupported"
//  UNSUPPORTEDOPERATION_MARKETIMAGEEXPORTUNSUPPORTED = "UnsupportedOperation.MarketImageExportUnsupported"
//  UNSUPPORTEDOPERATION_PUBLICIMAGEEXPORTUNSUPPORTED = "UnsupportedOperation.PublicImageExportUnsupported"
//  UNSUPPORTEDOPERATION_REDHATIMAGEEXPORTUNSUPPORTED = "UnsupportedOperation.RedHatImageExportUnsupported"
//  UNSUPPORTEDOPERATION_SHAREDIMAGEEXPORTUNSUPPORTED = "UnsupportedOperation.SharedImageExportUnsupported"
//  UNSUPPORTEDOPERATION_WINDOWSIMAGEEXPORTUNSUPPORTED = "UnsupportedOperation.WindowsImageExportUnsupported"
func (c *Client) ExportImages(request *ExportImagesRequest) (response *ExportImagesResponse, err error) {
    return c.ExportImagesWithContext(context.Background(), request)
}

// ExportImages
// This API is used to export custom images to the specified COS bucket.
//
// error code that may be returned:
//  AUTHFAILURE_CAMROLENAMEAUTHENTICATEFAILED = "AuthFailure.CamRoleNameAuthenticateFailed"
//  INVALIDIMAGEID_MALFORMED = "InvalidImageId.Malformed"
//  INVALIDIMAGEID_NOTFOUND = "InvalidImageId.NotFound"
//  INVALIDPARAMETER_IMAGEIDSSNAPSHOTIDSMUSTONE = "InvalidParameter.ImageIdsSnapshotIdsMustOne"
//  INVALIDPARAMETER_SNAPSHOTNOTFOUND = "InvalidParameter.SnapshotNotFound"
//  INVALIDPARAMETERVALUE_BUCKETNOTFOUND = "InvalidParameterValue.BucketNotFound"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_INVALIDBUCKETPERMISSIONFOREXPORT = "InvalidParameterValue.InvalidBucketPermissionForExport"
//  INVALIDPARAMETERVALUE_INVALIDFILENAMEPREFIXLIST = "InvalidParameterValue.InvalidFileNamePrefixList"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEID = "InvalidParameterValue.InvalidImageId"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEOSNAME = "InvalidParameterValue.InvalidImageOsName"
//  INVALIDPARAMETERVALUE_INVALIDIMAGESTATE = "InvalidParameterValue.InvalidImageState"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_TOOLARGE = "InvalidParameterValue.TooLarge"
//  LIMITEXCEEDED_EXPORTIMAGETASKLIMITEXCEEDED = "LimitExceeded.ExportImageTaskLimitExceeded"
//  UNSUPPORTEDOPERATION_ENCRYPTEDIMAGESNOTSUPPORTED = "UnsupportedOperation.EncryptedImagesNotSupported"
//  UNSUPPORTEDOPERATION_IMAGETOOLARGEEXPORTUNSUPPORTED = "UnsupportedOperation.ImageTooLargeExportUnsupported"
//  UNSUPPORTEDOPERATION_MARKETIMAGEEXPORTUNSUPPORTED = "UnsupportedOperation.MarketImageExportUnsupported"
//  UNSUPPORTEDOPERATION_PUBLICIMAGEEXPORTUNSUPPORTED = "UnsupportedOperation.PublicImageExportUnsupported"
//  UNSUPPORTEDOPERATION_REDHATIMAGEEXPORTUNSUPPORTED = "UnsupportedOperation.RedHatImageExportUnsupported"
//  UNSUPPORTEDOPERATION_SHAREDIMAGEEXPORTUNSUPPORTED = "UnsupportedOperation.SharedImageExportUnsupported"
//  UNSUPPORTEDOPERATION_WINDOWSIMAGEEXPORTUNSUPPORTED = "UnsupportedOperation.WindowsImageExportUnsupported"
func (c *Client) ExportImagesWithContext(ctx context.Context, request *ExportImagesRequest) (response *ExportImagesResponse, err error) {
    if request == nil {
        request = NewExportImagesRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "ExportImages")
    
    if c.GetCredential() == nil {
        return nil, errors.New("ExportImages require credential")
    }

    request.SetContext(ctx)
    
    response = NewExportImagesResponse()
    err = c.Send(request, response)
    return
}

func NewImportImageRequest() (request *ImportImageRequest) {
    request = &ImportImageRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "ImportImage")
    
    
    return
}

func NewImportImageResponse() (response *ImportImageResponse) {
    response = &ImportImageResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// ImportImage
// This API is used to import an image. The image imported can be used to create instances. Currently, this API supports RAW, VHD, QCOW2, and VMDK image formats.
//
// error code that may be returned:
//  IMAGEQUOTALIMITEXCEEDED = "ImageQuotaLimitExceeded"
//  INVALIDIMAGENAME_DUPLICATE = "InvalidImageName.Duplicate"
//  INVALIDIMAGEOSTYPE_UNSUPPORTED = "InvalidImageOsType.Unsupported"
//  INVALIDIMAGEOSVERSION_UNSUPPORTED = "InvalidImageOsVersion.Unsupported"
//  INVALIDPARAMETER_INVALIDPARAMETERURLERROR = "InvalidParameter.InvalidParameterUrlError"
//  INVALIDPARAMETERVALUE_ISOMUSTIMPORTBYFORCE = "InvalidParameterValue.ISOMustImportByForce"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_INVALIDBOOTMODE = "InvalidParameterValue.InvalidBootMode"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEFAMILY = "InvalidParameterValue.InvalidImageFamily"
//  INVALIDPARAMETERVALUE_INVALIDLICENSETYPE = "InvalidParameterValue.InvalidLicenseType"
//  INVALIDPARAMETERVALUE_TAGKEYNOTFOUND = "InvalidParameterValue.TagKeyNotFound"
//  INVALIDPARAMETERVALUE_TOOLARGE = "InvalidParameterValue.TooLarge"
//  OPERATIONDENIED_INNERUSERPROHIBITACTION = "OperationDenied.InnerUserProhibitAction"
//  REGIONABILITYLIMIT_UNSUPPORTEDTOIMPORTIMAGE = "RegionAbilityLimit.UnsupportedToImportImage"
func (c *Client) ImportImage(request *ImportImageRequest) (response *ImportImageResponse, err error) {
    return c.ImportImageWithContext(context.Background(), request)
}

// ImportImage
// This API is used to import an image. The image imported can be used to create instances. Currently, this API supports RAW, VHD, QCOW2, and VMDK image formats.
//
// error code that may be returned:
//  IMAGEQUOTALIMITEXCEEDED = "ImageQuotaLimitExceeded"
//  INVALIDIMAGENAME_DUPLICATE = "InvalidImageName.Duplicate"
//  INVALIDIMAGEOSTYPE_UNSUPPORTED = "InvalidImageOsType.Unsupported"
//  INVALIDIMAGEOSVERSION_UNSUPPORTED = "InvalidImageOsVersion.Unsupported"
//  INVALIDPARAMETER_INVALIDPARAMETERURLERROR = "InvalidParameter.InvalidParameterUrlError"
//  INVALIDPARAMETERVALUE_ISOMUSTIMPORTBYFORCE = "InvalidParameterValue.ISOMustImportByForce"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_INVALIDBOOTMODE = "InvalidParameterValue.InvalidBootMode"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEFAMILY = "InvalidParameterValue.InvalidImageFamily"
//  INVALIDPARAMETERVALUE_INVALIDLICENSETYPE = "InvalidParameterValue.InvalidLicenseType"
//  INVALIDPARAMETERVALUE_TAGKEYNOTFOUND = "InvalidParameterValue.TagKeyNotFound"
//  INVALIDPARAMETERVALUE_TOOLARGE = "InvalidParameterValue.TooLarge"
//  OPERATIONDENIED_INNERUSERPROHIBITACTION = "OperationDenied.InnerUserProhibitAction"
//  REGIONABILITYLIMIT_UNSUPPORTEDTOIMPORTIMAGE = "RegionAbilityLimit.UnsupportedToImportImage"
func (c *Client) ImportImageWithContext(ctx context.Context, request *ImportImageRequest) (response *ImportImageResponse, err error) {
    if request == nil {
        request = NewImportImageRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "ImportImage")
    
    if c.GetCredential() == nil {
        return nil, errors.New("ImportImage require credential")
    }

    request.SetContext(ctx)
    
    response = NewImportImageResponse()
    err = c.Send(request, response)
    return
}

func NewImportKeyPairRequest() (request *ImportKeyPairRequest) {
    request = &ImportKeyPairRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "ImportKeyPair")
    
    
    return
}

func NewImportKeyPairResponse() (response *ImportKeyPairResponse) {
    response = &ImportKeyPairResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// ImportKeyPair
// This API is used to import key pairs.
//
// 
//
// * You can use this API to import key pairs to a user account, but the key pairs will not be automatically associated with any instance. You may use [AssociasteInstancesKeyPair](https://intl.cloud.tencent.com/document/api/213/15698?from_cn_redirect=1) to associate key pairs with instances.
//
// * You need to specify the names of the key pairs and the content of the public keys.
//
// * If you only have private keys, you can convert them to public keys with the `SSL` tool before importing them.
//
// error code that may be returned:
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDKEYPAIR_LIMITEXCEEDED = "InvalidKeyPair.LimitExceeded"
//  INVALIDKEYPAIRNAME_DUPLICATE = "InvalidKeyPairName.Duplicate"
//  INVALIDKEYPAIRNAMEEMPTY = "InvalidKeyPairNameEmpty"
//  INVALIDKEYPAIRNAMEINCLUDEILLEGALCHAR = "InvalidKeyPairNameIncludeIllegalChar"
//  INVALIDKEYPAIRNAMETOOLONG = "InvalidKeyPairNameTooLong"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPROJECTID_NOTFOUND = "InvalidProjectId.NotFound"
//  INVALIDPUBLICKEY_DUPLICATE = "InvalidPublicKey.Duplicate"
//  INVALIDPUBLICKEY_MALFORMED = "InvalidPublicKey.Malformed"
//  LIMITEXCEEDED_TAGRESOURCEQUOTA = "LimitExceeded.TagResourceQuota"
//  MISSINGPARAMETER = "MissingParameter"
func (c *Client) ImportKeyPair(request *ImportKeyPairRequest) (response *ImportKeyPairResponse, err error) {
    return c.ImportKeyPairWithContext(context.Background(), request)
}

// ImportKeyPair
// This API is used to import key pairs.
//
// 
//
// * You can use this API to import key pairs to a user account, but the key pairs will not be automatically associated with any instance. You may use [AssociasteInstancesKeyPair](https://intl.cloud.tencent.com/document/api/213/15698?from_cn_redirect=1) to associate key pairs with instances.
//
// * You need to specify the names of the key pairs and the content of the public keys.
//
// * If you only have private keys, you can convert them to public keys with the `SSL` tool before importing them.
//
// error code that may be returned:
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDKEYPAIR_LIMITEXCEEDED = "InvalidKeyPair.LimitExceeded"
//  INVALIDKEYPAIRNAME_DUPLICATE = "InvalidKeyPairName.Duplicate"
//  INVALIDKEYPAIRNAMEEMPTY = "InvalidKeyPairNameEmpty"
//  INVALIDKEYPAIRNAMEINCLUDEILLEGALCHAR = "InvalidKeyPairNameIncludeIllegalChar"
//  INVALIDKEYPAIRNAMETOOLONG = "InvalidKeyPairNameTooLong"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPROJECTID_NOTFOUND = "InvalidProjectId.NotFound"
//  INVALIDPUBLICKEY_DUPLICATE = "InvalidPublicKey.Duplicate"
//  INVALIDPUBLICKEY_MALFORMED = "InvalidPublicKey.Malformed"
//  LIMITEXCEEDED_TAGRESOURCEQUOTA = "LimitExceeded.TagResourceQuota"
//  MISSINGPARAMETER = "MissingParameter"
func (c *Client) ImportKeyPairWithContext(ctx context.Context, request *ImportKeyPairRequest) (response *ImportKeyPairResponse, err error) {
    if request == nil {
        request = NewImportKeyPairRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "ImportKeyPair")
    
    if c.GetCredential() == nil {
        return nil, errors.New("ImportKeyPair require credential")
    }

    request.SetContext(ctx)
    
    response = NewImportKeyPairResponse()
    err = c.Send(request, response)
    return
}

func NewInquirePricePurchaseReservedInstancesOfferingRequest() (request *InquirePricePurchaseReservedInstancesOfferingRequest) {
    request = &InquirePricePurchaseReservedInstancesOfferingRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "InquirePricePurchaseReservedInstancesOffering")
    
    
    return
}

func NewInquirePricePurchaseReservedInstancesOfferingResponse() (response *InquirePricePurchaseReservedInstancesOfferingResponse) {
    response = &InquirePricePurchaseReservedInstancesOfferingResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// InquirePricePurchaseReservedInstancesOffering
// This API is used to query the price of reserved instances. It only supports querying purchasable reserved instance offerings. Currently, RIs are only offered to beta users.
//
// error code that may be returned:
//  FAILEDOPERATION_INQUIRYPRICEFAILED = "FailedOperation.InquiryPriceFailed"
//  INVALIDPARAMETERVALUE_RANGE = "InvalidParameterValue.Range"
//  UNSUPPORTEDOPERATION_INVALIDPERMISSIONNONINTERNATIONALACCOUNT = "UnsupportedOperation.InvalidPermissionNonInternationalAccount"
func (c *Client) InquirePricePurchaseReservedInstancesOffering(request *InquirePricePurchaseReservedInstancesOfferingRequest) (response *InquirePricePurchaseReservedInstancesOfferingResponse, err error) {
    return c.InquirePricePurchaseReservedInstancesOfferingWithContext(context.Background(), request)
}

// InquirePricePurchaseReservedInstancesOffering
// This API is used to query the price of reserved instances. It only supports querying purchasable reserved instance offerings. Currently, RIs are only offered to beta users.
//
// error code that may be returned:
//  FAILEDOPERATION_INQUIRYPRICEFAILED = "FailedOperation.InquiryPriceFailed"
//  INVALIDPARAMETERVALUE_RANGE = "InvalidParameterValue.Range"
//  UNSUPPORTEDOPERATION_INVALIDPERMISSIONNONINTERNATIONALACCOUNT = "UnsupportedOperation.InvalidPermissionNonInternationalAccount"
func (c *Client) InquirePricePurchaseReservedInstancesOfferingWithContext(ctx context.Context, request *InquirePricePurchaseReservedInstancesOfferingRequest) (response *InquirePricePurchaseReservedInstancesOfferingResponse, err error) {
    if request == nil {
        request = NewInquirePricePurchaseReservedInstancesOfferingRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "InquirePricePurchaseReservedInstancesOffering")
    
    if c.GetCredential() == nil {
        return nil, errors.New("InquirePricePurchaseReservedInstancesOffering require credential")
    }

    request.SetContext(ctx)
    
    response = NewInquirePricePurchaseReservedInstancesOfferingResponse()
    err = c.Send(request, response)
    return
}

func NewInquiryPriceModifyInstancesChargeTypeRequest() (request *InquiryPriceModifyInstancesChargeTypeRequest) {
    request = &InquiryPriceModifyInstancesChargeTypeRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "InquiryPriceModifyInstancesChargeType")
    
    
    return
}

func NewInquiryPriceModifyInstancesChargeTypeResponse() (response *InquiryPriceModifyInstancesChargeTypeResponse) {
    response = &InquiryPriceModifyInstancesChargeTypeResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// InquiryPriceModifyInstancesChargeType
// This API is used to inquire about the price for switching billing modes of instance.
//
// 
//
// 
//
// This API is used to indicate that instances with no charge when shut down, instances of the model families Batch Computing BC1 and Batch Computing BS1, instances of scheduled termination, and spot instances do not support this operation.
//
// error code that may be returned:
//  ACCOUNTQUALIFICATIONRESTRICTIONS = "AccountQualificationRestrictions"
//  FAILEDOPERATION_INQUIRYPRICEFAILED = "FailedOperation.InquiryPriceFailed"
//  FAILEDOPERATION_INQUIRYREFUNDPRICEFAILED = "FailedOperation.InquiryRefundPriceFailed"
//  INTERNALERROR_TRADEUNKNOWNERROR = "InternalError.TradeUnknownError"
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEID = "InvalidParameterValue.InvalidImageId"
//  INVALIDPARAMETERVALUE_INVALIDINSTANCETYPEUNDERWRITE = "InvalidParameterValue.InvalidInstanceTypeUnderwrite"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPERIOD = "InvalidPeriod"
//  INVALIDPERMISSION = "InvalidPermission"
//  LIMITEXCEEDED_INSTANCETYPEBANDWIDTH = "LimitExceeded.InstanceTypeBandwidth"
//  MISSINGPARAMETER = "MissingParameter"
//  RESOURCEINSUFFICIENT_CLOUDDISKUNAVAILABLE = "ResourceInsufficient.CloudDiskUnavailable"
//  UNSUPPORTEDOPERATION_INSTANCECHARGETYPE = "UnsupportedOperation.InstanceChargeType"
//  UNSUPPORTEDOPERATION_INSTANCEMIXEDZONETYPE = "UnsupportedOperation.InstanceMixedZoneType"
func (c *Client) InquiryPriceModifyInstancesChargeType(request *InquiryPriceModifyInstancesChargeTypeRequest) (response *InquiryPriceModifyInstancesChargeTypeResponse, err error) {
    return c.InquiryPriceModifyInstancesChargeTypeWithContext(context.Background(), request)
}

// InquiryPriceModifyInstancesChargeType
// This API is used to inquire about the price for switching billing modes of instance.
//
// 
//
// 
//
// This API is used to indicate that instances with no charge when shut down, instances of the model families Batch Computing BC1 and Batch Computing BS1, instances of scheduled termination, and spot instances do not support this operation.
//
// error code that may be returned:
//  ACCOUNTQUALIFICATIONRESTRICTIONS = "AccountQualificationRestrictions"
//  FAILEDOPERATION_INQUIRYPRICEFAILED = "FailedOperation.InquiryPriceFailed"
//  FAILEDOPERATION_INQUIRYREFUNDPRICEFAILED = "FailedOperation.InquiryRefundPriceFailed"
//  INTERNALERROR_TRADEUNKNOWNERROR = "InternalError.TradeUnknownError"
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEID = "InvalidParameterValue.InvalidImageId"
//  INVALIDPARAMETERVALUE_INVALIDINSTANCETYPEUNDERWRITE = "InvalidParameterValue.InvalidInstanceTypeUnderwrite"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPERIOD = "InvalidPeriod"
//  INVALIDPERMISSION = "InvalidPermission"
//  LIMITEXCEEDED_INSTANCETYPEBANDWIDTH = "LimitExceeded.InstanceTypeBandwidth"
//  MISSINGPARAMETER = "MissingParameter"
//  RESOURCEINSUFFICIENT_CLOUDDISKUNAVAILABLE = "ResourceInsufficient.CloudDiskUnavailable"
//  UNSUPPORTEDOPERATION_INSTANCECHARGETYPE = "UnsupportedOperation.InstanceChargeType"
//  UNSUPPORTEDOPERATION_INSTANCEMIXEDZONETYPE = "UnsupportedOperation.InstanceMixedZoneType"
func (c *Client) InquiryPriceModifyInstancesChargeTypeWithContext(ctx context.Context, request *InquiryPriceModifyInstancesChargeTypeRequest) (response *InquiryPriceModifyInstancesChargeTypeResponse, err error) {
    if request == nil {
        request = NewInquiryPriceModifyInstancesChargeTypeRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "InquiryPriceModifyInstancesChargeType")
    
    if c.GetCredential() == nil {
        return nil, errors.New("InquiryPriceModifyInstancesChargeType require credential")
    }

    request.SetContext(ctx)
    
    response = NewInquiryPriceModifyInstancesChargeTypeResponse()
    err = c.Send(request, response)
    return
}

func NewInquiryPriceRenewInstancesRequest() (request *InquiryPriceRenewInstancesRequest) {
    request = &InquiryPriceRenewInstancesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "InquiryPriceRenewInstances")
    
    
    return
}

func NewInquiryPriceRenewInstancesResponse() (response *InquiryPriceRenewInstancesResponse) {
    response = &InquiryPriceRenewInstancesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// InquiryPriceRenewInstances
// This API is used to inquire about the price for renewing a monthly subscription instance.
//
// 
//
// This API is used to query the renewal price of monthly subscription instances.
//
// error code that may be returned:
//  INTERNALERROR_TRADEUNKNOWNERROR = "InternalError.TradeUnknownError"
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDACCOUNT_INSUFFICIENTBALANCE = "InvalidAccount.InsufficientBalance"
//  INVALIDACCOUNT_UNPAIDORDER = "InvalidAccount.UnpaidOrder"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_INSTANCENOTSUPPORTEDMIXPRICINGMODEL = "InvalidParameterValue.InstanceNotSupportedMixPricingModel"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_RANGE = "InvalidParameterValue.Range"
//  INVALIDPERIOD = "InvalidPeriod"
//  MISSINGPARAMETER = "MissingParameter"
//  RESOURCEINSUFFICIENT_CLOUDDISKUNAVAILABLE = "ResourceInsufficient.CloudDiskUnavailable"
//  RESOURCEUNAVAILABLE = "ResourceUnavailable"
//  UNSUPPORTEDOPERATION_INSTANCEMIXEDZONETYPE = "UnsupportedOperation.InstanceMixedZoneType"
//  UNSUPPORTEDOPERATION_INSTANCESTATEPENDING = "UnsupportedOperation.InstanceStatePending"
//  UNSUPPORTEDOPERATION_INVALIDDISKBACKUPQUOTA = "UnsupportedOperation.InvalidDiskBackupQuota"
//  UNSUPPORTEDOPERATION_UNDERWRITINGINSTANCETYPEONLYSUPPORTAUTORENEW = "UnsupportedOperation.UnderwritingInstanceTypeOnlySupportAutoRenew"
func (c *Client) InquiryPriceRenewInstances(request *InquiryPriceRenewInstancesRequest) (response *InquiryPriceRenewInstancesResponse, err error) {
    return c.InquiryPriceRenewInstancesWithContext(context.Background(), request)
}

// InquiryPriceRenewInstances
// This API is used to inquire about the price for renewing a monthly subscription instance.
//
// 
//
// This API is used to query the renewal price of monthly subscription instances.
//
// error code that may be returned:
//  INTERNALERROR_TRADEUNKNOWNERROR = "InternalError.TradeUnknownError"
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDACCOUNT_INSUFFICIENTBALANCE = "InvalidAccount.InsufficientBalance"
//  INVALIDACCOUNT_UNPAIDORDER = "InvalidAccount.UnpaidOrder"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_INSTANCENOTSUPPORTEDMIXPRICINGMODEL = "InvalidParameterValue.InstanceNotSupportedMixPricingModel"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_RANGE = "InvalidParameterValue.Range"
//  INVALIDPERIOD = "InvalidPeriod"
//  MISSINGPARAMETER = "MissingParameter"
//  RESOURCEINSUFFICIENT_CLOUDDISKUNAVAILABLE = "ResourceInsufficient.CloudDiskUnavailable"
//  RESOURCEUNAVAILABLE = "ResourceUnavailable"
//  UNSUPPORTEDOPERATION_INSTANCEMIXEDZONETYPE = "UnsupportedOperation.InstanceMixedZoneType"
//  UNSUPPORTEDOPERATION_INSTANCESTATEPENDING = "UnsupportedOperation.InstanceStatePending"
//  UNSUPPORTEDOPERATION_INVALIDDISKBACKUPQUOTA = "UnsupportedOperation.InvalidDiskBackupQuota"
//  UNSUPPORTEDOPERATION_UNDERWRITINGINSTANCETYPEONLYSUPPORTAUTORENEW = "UnsupportedOperation.UnderwritingInstanceTypeOnlySupportAutoRenew"
func (c *Client) InquiryPriceRenewInstancesWithContext(ctx context.Context, request *InquiryPriceRenewInstancesRequest) (response *InquiryPriceRenewInstancesResponse, err error) {
    if request == nil {
        request = NewInquiryPriceRenewInstancesRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "InquiryPriceRenewInstances")
    
    if c.GetCredential() == nil {
        return nil, errors.New("InquiryPriceRenewInstances require credential")
    }

    request.SetContext(ctx)
    
    response = NewInquiryPriceRenewInstancesResponse()
    err = c.Send(request, response)
    return
}

func NewInquiryPriceResetInstanceRequest() (request *InquiryPriceResetInstanceRequest) {
    request = &InquiryPriceResetInstanceRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "InquiryPriceResetInstance")
    
    
    return
}

func NewInquiryPriceResetInstanceResponse() (response *InquiryPriceResetInstanceResponse) {
    response = &InquiryPriceResetInstanceResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// InquiryPriceResetInstance
// This API is used to inquire about the price for reinstalling an instance.
//
// 
//
// * If you have specified the parameter `ImageId`, inquire about the price for reinstallation by using the specified image. Otherwise, inquire about the price for reinstallation based on the image currently used by the instance.
//
// * Currently, only instances with a [system disk type](https://intl.cloud.tencent.com/document/api/213/15753?from_cn_redirect=1#SystemDisk) of `CLOUD_BSSD`, `CLOUD_PREMIUM`, or `CLOUD_SSD` are supported for inquiring about the price for reinstallation caused by switching between `Linux` and `Windows` operating systems through this API.
//
// * Currently, instances in regions outside the Chinese mainland are not supported for inquiring about the price for reinstallation caused by switching between `Linux` and `Windows` operating systems through this API.
//
// error code that may be returned:
//  FAILEDOPERATION_INQUIRYPRICEFAILED = "FailedOperation.InquiryPriceFailed"
//  INTERNALERROR_TRADEUNKNOWNERROR = "InternalError.TradeUnknownError"
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDIMAGEID_MALFORMED = "InvalidImageId.Malformed"
//  INVALIDIMAGEID_NOTFOUND = "InvalidImageId.NotFound"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_INSTANCETYPENOTFOUND = "InvalidParameterValue.InstanceTypeNotFound"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEID = "InvalidParameterValue.InvalidImageId"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEIDFORRETSETINSTANCE = "InvalidParameterValue.InvalidImageIdForRetsetInstance"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEOSNAME = "InvalidParameterValue.InvalidImageOsName"
//  INVALIDPARAMETERVALUE_ZONENOTSUPPORTED = "InvalidParameterValue.ZoneNotSupported"
//  MISSINGPARAMETER = "MissingParameter"
//  MISSINGPARAMETER_MONITORSERVICE = "MissingParameter.MonitorService"
//  RESOURCEINSUFFICIENT_CLOUDDISKUNAVAILABLE = "ResourceInsufficient.CloudDiskUnavailable"
//  RESOURCESSOLDOUT_SPECIFIEDINSTANCETYPE = "ResourcesSoldOut.SpecifiedInstanceType"
//  UNSUPPORTEDOPERATION_INSTANCETYPENOTSUPPORTGRIDLICENCE = "UnsupportedOperation.InstanceTypeNotSupportGridLicence"
//  UNSUPPORTEDOPERATION_INVALIDIMAGELICENSETYPEFORRESET = "UnsupportedOperation.InvalidImageLicenseTypeForReset"
//  UNSUPPORTEDOPERATION_MODIFYENCRYPTIONNOTSUPPORTED = "UnsupportedOperation.ModifyEncryptionNotSupported"
//  UNSUPPORTEDOPERATION_RAWLOCALDISKINSREINSTALLTOQCOW2 = "UnsupportedOperation.RawLocalDiskInsReinstalltoQcow2"
func (c *Client) InquiryPriceResetInstance(request *InquiryPriceResetInstanceRequest) (response *InquiryPriceResetInstanceResponse, err error) {
    return c.InquiryPriceResetInstanceWithContext(context.Background(), request)
}

// InquiryPriceResetInstance
// This API is used to inquire about the price for reinstalling an instance.
//
// 
//
// * If you have specified the parameter `ImageId`, inquire about the price for reinstallation by using the specified image. Otherwise, inquire about the price for reinstallation based on the image currently used by the instance.
//
// * Currently, only instances with a [system disk type](https://intl.cloud.tencent.com/document/api/213/15753?from_cn_redirect=1#SystemDisk) of `CLOUD_BSSD`, `CLOUD_PREMIUM`, or `CLOUD_SSD` are supported for inquiring about the price for reinstallation caused by switching between `Linux` and `Windows` operating systems through this API.
//
// * Currently, instances in regions outside the Chinese mainland are not supported for inquiring about the price for reinstallation caused by switching between `Linux` and `Windows` operating systems through this API.
//
// error code that may be returned:
//  FAILEDOPERATION_INQUIRYPRICEFAILED = "FailedOperation.InquiryPriceFailed"
//  INTERNALERROR_TRADEUNKNOWNERROR = "InternalError.TradeUnknownError"
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDIMAGEID_MALFORMED = "InvalidImageId.Malformed"
//  INVALIDIMAGEID_NOTFOUND = "InvalidImageId.NotFound"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_INSTANCETYPENOTFOUND = "InvalidParameterValue.InstanceTypeNotFound"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEID = "InvalidParameterValue.InvalidImageId"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEIDFORRETSETINSTANCE = "InvalidParameterValue.InvalidImageIdForRetsetInstance"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEOSNAME = "InvalidParameterValue.InvalidImageOsName"
//  INVALIDPARAMETERVALUE_ZONENOTSUPPORTED = "InvalidParameterValue.ZoneNotSupported"
//  MISSINGPARAMETER = "MissingParameter"
//  MISSINGPARAMETER_MONITORSERVICE = "MissingParameter.MonitorService"
//  RESOURCEINSUFFICIENT_CLOUDDISKUNAVAILABLE = "ResourceInsufficient.CloudDiskUnavailable"
//  RESOURCESSOLDOUT_SPECIFIEDINSTANCETYPE = "ResourcesSoldOut.SpecifiedInstanceType"
//  UNSUPPORTEDOPERATION_INSTANCETYPENOTSUPPORTGRIDLICENCE = "UnsupportedOperation.InstanceTypeNotSupportGridLicence"
//  UNSUPPORTEDOPERATION_INVALIDIMAGELICENSETYPEFORRESET = "UnsupportedOperation.InvalidImageLicenseTypeForReset"
//  UNSUPPORTEDOPERATION_MODIFYENCRYPTIONNOTSUPPORTED = "UnsupportedOperation.ModifyEncryptionNotSupported"
//  UNSUPPORTEDOPERATION_RAWLOCALDISKINSREINSTALLTOQCOW2 = "UnsupportedOperation.RawLocalDiskInsReinstalltoQcow2"
func (c *Client) InquiryPriceResetInstanceWithContext(ctx context.Context, request *InquiryPriceResetInstanceRequest) (response *InquiryPriceResetInstanceResponse, err error) {
    if request == nil {
        request = NewInquiryPriceResetInstanceRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "InquiryPriceResetInstance")
    
    if c.GetCredential() == nil {
        return nil, errors.New("InquiryPriceResetInstance require credential")
    }

    request.SetContext(ctx)
    
    response = NewInquiryPriceResetInstanceResponse()
    err = c.Send(request, response)
    return
}

func NewInquiryPriceResetInstancesInternetMaxBandwidthRequest() (request *InquiryPriceResetInstancesInternetMaxBandwidthRequest) {
    request = &InquiryPriceResetInstancesInternetMaxBandwidthRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "InquiryPriceResetInstancesInternetMaxBandwidth")
    
    
    return
}

func NewInquiryPriceResetInstancesInternetMaxBandwidthResponse() (response *InquiryPriceResetInstancesInternetMaxBandwidthResponse) {
    response = &InquiryPriceResetInstancesInternetMaxBandwidthResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// InquiryPriceResetInstancesInternetMaxBandwidth
// This API is used to query the price for upgrading the public bandwidth cap of an instance.
//
// 
//
// * The allowed bandwidth cap varies for different models. For details, see [Purchasing Network Bandwidth](https://intl.cloud.tencent.com/document/product/213/509?from_cn_redirect=1).
//
// * For bandwidth billed by the `TRAFFIC_POSTPAID_BY_HOUR` method, changing the bandwidth cap through this API takes effect in real time. You can increase or reduce bandwidth within applicable limits.
//
// error code that may be returned:
//  FAILEDOPERATION_INQUIRYPRICEFAILED = "FailedOperation.InquiryPriceFailed"
//  FAILEDOPERATION_NOTFOUNDEIP = "FailedOperation.NotFoundEIP"
//  INTERNALERROR_TRADEUNKNOWNERROR = "InternalError.TradeUnknownError"
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDACCOUNT_INSUFFICIENTBALANCE = "InvalidAccount.InsufficientBalance"
//  INVALIDACCOUNT_UNPAIDORDER = "InvalidAccount.UnpaidOrder"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERCOMBINATION = "InvalidParameterCombination"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_RANGE = "InvalidParameterValue.Range"
//  INVALIDPERMISSION = "InvalidPermission"
//  MISSINGPARAMETER = "MissingParameter"
//  UNSUPPORTEDOPERATION_BANDWIDTHPACKAGEIDNOTSUPPORTED = "UnsupportedOperation.BandwidthPackageIdNotSupported"
func (c *Client) InquiryPriceResetInstancesInternetMaxBandwidth(request *InquiryPriceResetInstancesInternetMaxBandwidthRequest) (response *InquiryPriceResetInstancesInternetMaxBandwidthResponse, err error) {
    return c.InquiryPriceResetInstancesInternetMaxBandwidthWithContext(context.Background(), request)
}

// InquiryPriceResetInstancesInternetMaxBandwidth
// This API is used to query the price for upgrading the public bandwidth cap of an instance.
//
// 
//
// * The allowed bandwidth cap varies for different models. For details, see [Purchasing Network Bandwidth](https://intl.cloud.tencent.com/document/product/213/509?from_cn_redirect=1).
//
// * For bandwidth billed by the `TRAFFIC_POSTPAID_BY_HOUR` method, changing the bandwidth cap through this API takes effect in real time. You can increase or reduce bandwidth within applicable limits.
//
// error code that may be returned:
//  FAILEDOPERATION_INQUIRYPRICEFAILED = "FailedOperation.InquiryPriceFailed"
//  FAILEDOPERATION_NOTFOUNDEIP = "FailedOperation.NotFoundEIP"
//  INTERNALERROR_TRADEUNKNOWNERROR = "InternalError.TradeUnknownError"
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDACCOUNT_INSUFFICIENTBALANCE = "InvalidAccount.InsufficientBalance"
//  INVALIDACCOUNT_UNPAIDORDER = "InvalidAccount.UnpaidOrder"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERCOMBINATION = "InvalidParameterCombination"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_RANGE = "InvalidParameterValue.Range"
//  INVALIDPERMISSION = "InvalidPermission"
//  MISSINGPARAMETER = "MissingParameter"
//  UNSUPPORTEDOPERATION_BANDWIDTHPACKAGEIDNOTSUPPORTED = "UnsupportedOperation.BandwidthPackageIdNotSupported"
func (c *Client) InquiryPriceResetInstancesInternetMaxBandwidthWithContext(ctx context.Context, request *InquiryPriceResetInstancesInternetMaxBandwidthRequest) (response *InquiryPriceResetInstancesInternetMaxBandwidthResponse, err error) {
    if request == nil {
        request = NewInquiryPriceResetInstancesInternetMaxBandwidthRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "InquiryPriceResetInstancesInternetMaxBandwidth")
    
    if c.GetCredential() == nil {
        return nil, errors.New("InquiryPriceResetInstancesInternetMaxBandwidth require credential")
    }

    request.SetContext(ctx)
    
    response = NewInquiryPriceResetInstancesInternetMaxBandwidthResponse()
    err = c.Send(request, response)
    return
}

func NewInquiryPriceResetInstancesTypeRequest() (request *InquiryPriceResetInstancesTypeRequest) {
    request = &InquiryPriceResetInstancesTypeRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "InquiryPriceResetInstancesType")
    
    
    return
}

func NewInquiryPriceResetInstancesTypeResponse() (response *InquiryPriceResetInstancesTypeResponse) {
    response = &InquiryPriceResetInstancesTypeResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// InquiryPriceResetInstancesType
// This API is used to query the price for adjusting the instance model.
//
// 
//
// * Currently, you can only use this API to query the prices of instances whose [system disk type](https://intl.cloud.tencent.com/document/api/213/9452?from_cn_redirect=1#block_device) is `CLOUD_BASIC`, `CLOUD_PREMIUM`, or `CLOUD_SSD`.
//
// * Currently, you cannot use this API to query the prices of [CDH](https://intl.cloud.tencent.com/document/product/416?from_cn_redirect=1) instances.
//
// error code that may be returned:
//  FAILEDOPERATION_INQUIRYREFUNDPRICEFAILED = "FailedOperation.InquiryRefundPriceFailed"
//  INTERNALERROR_TRADEUNKNOWNERROR = "InternalError.TradeUnknownError"
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDACCOUNT_INSUFFICIENTBALANCE = "InvalidAccount.InsufficientBalance"
//  INVALIDACCOUNT_UNPAIDORDER = "InvalidAccount.UnpaidOrder"
//  INVALIDHOSTID_MALFORMED = "InvalidHostId.Malformed"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDINSTANCETYPE_MALFORMED = "InvalidInstanceType.Malformed"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_HOSTIDCUSTOMIZEDINSTANCETYPENOTSUPPORT = "InvalidParameter.HostIdCustomizedInstanceTypeNotSupport"
//  INVALIDPARAMETER_HOSTIDINSTANCETYPENOTSUPPORT = "InvalidParameter.HostIdInstanceTypeNotSupport"
//  INVALIDPARAMETER_HOSTIDSTANDARDINSTANCETYPENOTSUPPORT = "InvalidParameter.HostIdStandardInstanceTypeNotSupport"
//  INVALIDPARAMETER_INSTANCETYPESUPPORTEDHOSTNOTFOUND = "InvalidParameter.InstanceTypeSupportedHostNotFound"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_BASICNETWORKINSTANCEFAMILY = "InvalidParameterValue.BasicNetworkInstanceFamily"
//  INVALIDPARAMETERVALUE_GPUINSTANCEFAMILY = "InvalidParameterValue.GPUInstanceFamily"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_INVALIDINSTANCESOURCE = "InvalidParameterValue.InvalidInstanceSource"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPERIOD = "InvalidPeriod"
//  INVALIDPERMISSION = "InvalidPermission"
//  INVALIDPROJECTID_NOTFOUND = "InvalidProjectId.NotFound"
//  INVALIDZONE_MISMATCHREGION = "InvalidZone.MismatchRegion"
//  LIMITEXCEEDED_EIPNUMLIMIT = "LimitExceeded.EipNumLimit"
//  LIMITEXCEEDED_ENINUMLIMIT = "LimitExceeded.EniNumLimit"
//  LIMITEXCEEDED_INSTANCETYPEBANDWIDTH = "LimitExceeded.InstanceTypeBandwidth"
//  MISSINGPARAMETER = "MissingParameter"
//  RESOURCEINSUFFICIENT_CLOUDDISKUNAVAILABLE = "ResourceInsufficient.CloudDiskUnavailable"
//  RESOURCEINSUFFICIENT_SPECIFIEDINSTANCETYPE = "ResourceInsufficient.SpecifiedInstanceType"
//  RESOURCEUNAVAILABLE_INSTANCETYPE = "ResourceUnavailable.InstanceType"
//  RESOURCESSOLDOUT_SPECIFIEDINSTANCETYPE = "ResourcesSoldOut.SpecifiedInstanceType"
//  UNSUPPORTEDOPERATION_HETEROGENEOUSCHANGEINSTANCEFAMILY = "UnsupportedOperation.HeterogeneousChangeInstanceFamily"
//  UNSUPPORTEDOPERATION_INSTANCEMIXEDRESETINSTANCETYPE = "UnsupportedOperation.InstanceMixedResetInstanceType"
//  UNSUPPORTEDOPERATION_LOCALDATADISKCHANGEINSTANCEFAMILY = "UnsupportedOperation.LocalDataDiskChangeInstanceFamily"
//  UNSUPPORTEDOPERATION_ORIGINALINSTANCETYPEINVALID = "UnsupportedOperation.OriginalInstanceTypeInvalid"
//  UNSUPPORTEDOPERATION_STOPPEDMODESTOPCHARGING = "UnsupportedOperation.StoppedModeStopCharging"
//  UNSUPPORTEDOPERATION_STOPPEDMODESTOPCHARGINGSAMEFAMILY = "UnsupportedOperation.StoppedModeStopChargingSameFamily"
//  UNSUPPORTEDOPERATION_UNSUPPORTEDARMCHANGEINSTANCEFAMILY = "UnsupportedOperation.UnsupportedARMChangeInstanceFamily"
//  UNSUPPORTEDOPERATION_UNSUPPORTEDCHANGEINSTANCEFAMILY = "UnsupportedOperation.UnsupportedChangeInstanceFamily"
//  UNSUPPORTEDOPERATION_UNSUPPORTEDCHANGEINSTANCEFAMILYTOARM = "UnsupportedOperation.UnsupportedChangeInstanceFamilyToARM"
//  UNSUPPORTEDOPERATION_UNSUPPORTEDCHANGEINSTANCETOTHISINSTANCEFAMILY = "UnsupportedOperation.UnsupportedChangeInstanceToThisInstanceFamily"
func (c *Client) InquiryPriceResetInstancesType(request *InquiryPriceResetInstancesTypeRequest) (response *InquiryPriceResetInstancesTypeResponse, err error) {
    return c.InquiryPriceResetInstancesTypeWithContext(context.Background(), request)
}

// InquiryPriceResetInstancesType
// This API is used to query the price for adjusting the instance model.
//
// 
//
// * Currently, you can only use this API to query the prices of instances whose [system disk type](https://intl.cloud.tencent.com/document/api/213/9452?from_cn_redirect=1#block_device) is `CLOUD_BASIC`, `CLOUD_PREMIUM`, or `CLOUD_SSD`.
//
// * Currently, you cannot use this API to query the prices of [CDH](https://intl.cloud.tencent.com/document/product/416?from_cn_redirect=1) instances.
//
// error code that may be returned:
//  FAILEDOPERATION_INQUIRYREFUNDPRICEFAILED = "FailedOperation.InquiryRefundPriceFailed"
//  INTERNALERROR_TRADEUNKNOWNERROR = "InternalError.TradeUnknownError"
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDACCOUNT_INSUFFICIENTBALANCE = "InvalidAccount.InsufficientBalance"
//  INVALIDACCOUNT_UNPAIDORDER = "InvalidAccount.UnpaidOrder"
//  INVALIDHOSTID_MALFORMED = "InvalidHostId.Malformed"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDINSTANCETYPE_MALFORMED = "InvalidInstanceType.Malformed"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_HOSTIDCUSTOMIZEDINSTANCETYPENOTSUPPORT = "InvalidParameter.HostIdCustomizedInstanceTypeNotSupport"
//  INVALIDPARAMETER_HOSTIDINSTANCETYPENOTSUPPORT = "InvalidParameter.HostIdInstanceTypeNotSupport"
//  INVALIDPARAMETER_HOSTIDSTANDARDINSTANCETYPENOTSUPPORT = "InvalidParameter.HostIdStandardInstanceTypeNotSupport"
//  INVALIDPARAMETER_INSTANCETYPESUPPORTEDHOSTNOTFOUND = "InvalidParameter.InstanceTypeSupportedHostNotFound"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_BASICNETWORKINSTANCEFAMILY = "InvalidParameterValue.BasicNetworkInstanceFamily"
//  INVALIDPARAMETERVALUE_GPUINSTANCEFAMILY = "InvalidParameterValue.GPUInstanceFamily"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_INVALIDINSTANCESOURCE = "InvalidParameterValue.InvalidInstanceSource"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPERIOD = "InvalidPeriod"
//  INVALIDPERMISSION = "InvalidPermission"
//  INVALIDPROJECTID_NOTFOUND = "InvalidProjectId.NotFound"
//  INVALIDZONE_MISMATCHREGION = "InvalidZone.MismatchRegion"
//  LIMITEXCEEDED_EIPNUMLIMIT = "LimitExceeded.EipNumLimit"
//  LIMITEXCEEDED_ENINUMLIMIT = "LimitExceeded.EniNumLimit"
//  LIMITEXCEEDED_INSTANCETYPEBANDWIDTH = "LimitExceeded.InstanceTypeBandwidth"
//  MISSINGPARAMETER = "MissingParameter"
//  RESOURCEINSUFFICIENT_CLOUDDISKUNAVAILABLE = "ResourceInsufficient.CloudDiskUnavailable"
//  RESOURCEINSUFFICIENT_SPECIFIEDINSTANCETYPE = "ResourceInsufficient.SpecifiedInstanceType"
//  RESOURCEUNAVAILABLE_INSTANCETYPE = "ResourceUnavailable.InstanceType"
//  RESOURCESSOLDOUT_SPECIFIEDINSTANCETYPE = "ResourcesSoldOut.SpecifiedInstanceType"
//  UNSUPPORTEDOPERATION_HETEROGENEOUSCHANGEINSTANCEFAMILY = "UnsupportedOperation.HeterogeneousChangeInstanceFamily"
//  UNSUPPORTEDOPERATION_INSTANCEMIXEDRESETINSTANCETYPE = "UnsupportedOperation.InstanceMixedResetInstanceType"
//  UNSUPPORTEDOPERATION_LOCALDATADISKCHANGEINSTANCEFAMILY = "UnsupportedOperation.LocalDataDiskChangeInstanceFamily"
//  UNSUPPORTEDOPERATION_ORIGINALINSTANCETYPEINVALID = "UnsupportedOperation.OriginalInstanceTypeInvalid"
//  UNSUPPORTEDOPERATION_STOPPEDMODESTOPCHARGING = "UnsupportedOperation.StoppedModeStopCharging"
//  UNSUPPORTEDOPERATION_STOPPEDMODESTOPCHARGINGSAMEFAMILY = "UnsupportedOperation.StoppedModeStopChargingSameFamily"
//  UNSUPPORTEDOPERATION_UNSUPPORTEDARMCHANGEINSTANCEFAMILY = "UnsupportedOperation.UnsupportedARMChangeInstanceFamily"
//  UNSUPPORTEDOPERATION_UNSUPPORTEDCHANGEINSTANCEFAMILY = "UnsupportedOperation.UnsupportedChangeInstanceFamily"
//  UNSUPPORTEDOPERATION_UNSUPPORTEDCHANGEINSTANCEFAMILYTOARM = "UnsupportedOperation.UnsupportedChangeInstanceFamilyToARM"
//  UNSUPPORTEDOPERATION_UNSUPPORTEDCHANGEINSTANCETOTHISINSTANCEFAMILY = "UnsupportedOperation.UnsupportedChangeInstanceToThisInstanceFamily"
func (c *Client) InquiryPriceResetInstancesTypeWithContext(ctx context.Context, request *InquiryPriceResetInstancesTypeRequest) (response *InquiryPriceResetInstancesTypeResponse, err error) {
    if request == nil {
        request = NewInquiryPriceResetInstancesTypeRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "InquiryPriceResetInstancesType")
    
    if c.GetCredential() == nil {
        return nil, errors.New("InquiryPriceResetInstancesType require credential")
    }

    request.SetContext(ctx)
    
    response = NewInquiryPriceResetInstancesTypeResponse()
    err = c.Send(request, response)
    return
}

func NewInquiryPriceResizeInstanceDisksRequest() (request *InquiryPriceResizeInstanceDisksRequest) {
    request = &InquiryPriceResizeInstanceDisksRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "InquiryPriceResizeInstanceDisks")
    
    
    return
}

func NewInquiryPriceResizeInstanceDisksResponse() (response *InquiryPriceResizeInstanceDisksResponse) {
    response = &InquiryPriceResizeInstanceDisksResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// InquiryPriceResizeInstanceDisks
// This API is used to query the price for expanding data disks of an instance.
//
// 
//
// * Currently, you can only use this API to query the price of non-elastic data disks whose [disk type](https://intl.cloud.tencent.com/document/api/213/9452?from_cn_redirect=1#block_device) is `CLOUD_BASIC`, `CLOUD_PREMIUM`, or `CLOUD_SSD`. You can use [`DescribeDisks`](https://intl.cloud.tencent.com/document/api/362/16315?from_cn_redirect=1) to check whether a disk is elastic. If the `Portable` field in the response is `false`, it means that the disk is non-elastic.
//
// * Currently, you cannot use this API to query the price for [CDH](https://intl.cloud.tencent.com/document/product/416?from_cn_redirect=1) instances. *Also, you can only query the price of expanding one data disk at a time.
//
// error code that may be returned:
//  INTERNALERROR_TRADEUNKNOWNERROR = "InternalError.TradeUnknownError"
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDACCOUNT_INSUFFICIENTBALANCE = "InvalidAccount.InsufficientBalance"
//  INVALIDACCOUNT_UNPAIDORDER = "InvalidAccount.UnpaidOrder"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_SNAPSHOTIDMALFORMED = "InvalidParameterValue.SnapshotIdMalformed"
//  INVALIDPARAMETERVALUE_ZONENOTSUPPORTED = "InvalidParameterValue.ZoneNotSupported"
//  MISSINGPARAMETER = "MissingParameter"
//  MISSINGPARAMETER_ATLEASTONE = "MissingParameter.AtLeastOne"
//  RESOURCEINSUFFICIENT_CLOUDDISKUNAVAILABLE = "ResourceInsufficient.CloudDiskUnavailable"
//  UNSUPPORTEDOPERATION_INSTANCESTATECORRUPTED = "UnsupportedOperation.InstanceStateCorrupted"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPED = "UnsupportedOperation.InstanceStateStopped"
//  UNSUPPORTEDOPERATION_INVALIDDISK = "UnsupportedOperation.InvalidDisk"
//  UNSUPPORTEDOPERATION_LOCALDISKMIGRATINGTOCLOUDDISK = "UnsupportedOperation.LocalDiskMigratingToCloudDisk"
func (c *Client) InquiryPriceResizeInstanceDisks(request *InquiryPriceResizeInstanceDisksRequest) (response *InquiryPriceResizeInstanceDisksResponse, err error) {
    return c.InquiryPriceResizeInstanceDisksWithContext(context.Background(), request)
}

// InquiryPriceResizeInstanceDisks
// This API is used to query the price for expanding data disks of an instance.
//
// 
//
// * Currently, you can only use this API to query the price of non-elastic data disks whose [disk type](https://intl.cloud.tencent.com/document/api/213/9452?from_cn_redirect=1#block_device) is `CLOUD_BASIC`, `CLOUD_PREMIUM`, or `CLOUD_SSD`. You can use [`DescribeDisks`](https://intl.cloud.tencent.com/document/api/362/16315?from_cn_redirect=1) to check whether a disk is elastic. If the `Portable` field in the response is `false`, it means that the disk is non-elastic.
//
// * Currently, you cannot use this API to query the price for [CDH](https://intl.cloud.tencent.com/document/product/416?from_cn_redirect=1) instances. *Also, you can only query the price of expanding one data disk at a time.
//
// error code that may be returned:
//  INTERNALERROR_TRADEUNKNOWNERROR = "InternalError.TradeUnknownError"
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDACCOUNT_INSUFFICIENTBALANCE = "InvalidAccount.InsufficientBalance"
//  INVALIDACCOUNT_UNPAIDORDER = "InvalidAccount.UnpaidOrder"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_SNAPSHOTIDMALFORMED = "InvalidParameterValue.SnapshotIdMalformed"
//  INVALIDPARAMETERVALUE_ZONENOTSUPPORTED = "InvalidParameterValue.ZoneNotSupported"
//  MISSINGPARAMETER = "MissingParameter"
//  MISSINGPARAMETER_ATLEASTONE = "MissingParameter.AtLeastOne"
//  RESOURCEINSUFFICIENT_CLOUDDISKUNAVAILABLE = "ResourceInsufficient.CloudDiskUnavailable"
//  UNSUPPORTEDOPERATION_INSTANCESTATECORRUPTED = "UnsupportedOperation.InstanceStateCorrupted"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPED = "UnsupportedOperation.InstanceStateStopped"
//  UNSUPPORTEDOPERATION_INVALIDDISK = "UnsupportedOperation.InvalidDisk"
//  UNSUPPORTEDOPERATION_LOCALDISKMIGRATINGTOCLOUDDISK = "UnsupportedOperation.LocalDiskMigratingToCloudDisk"
func (c *Client) InquiryPriceResizeInstanceDisksWithContext(ctx context.Context, request *InquiryPriceResizeInstanceDisksRequest) (response *InquiryPriceResizeInstanceDisksResponse, err error) {
    if request == nil {
        request = NewInquiryPriceResizeInstanceDisksRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "InquiryPriceResizeInstanceDisks")
    
    if c.GetCredential() == nil {
        return nil, errors.New("InquiryPriceResizeInstanceDisks require credential")
    }

    request.SetContext(ctx)
    
    response = NewInquiryPriceResizeInstanceDisksResponse()
    err = c.Send(request, response)
    return
}

func NewInquiryPriceRunInstancesRequest() (request *InquiryPriceRunInstancesRequest) {
    request = &InquiryPriceRunInstancesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "InquiryPriceRunInstances")
    
    
    return
}

func NewInquiryPriceRunInstancesResponse() (response *InquiryPriceRunInstancesResponse) {
    response = &InquiryPriceRunInstancesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// InquiryPriceRunInstances
// This API is used to query the price of creating instances. You can only use this API for instances whose configuration is within the purchase limit. For more information, see [RunInstances](https://intl.cloud.tencent.com/document/api/213/15730?from_cn_redirect=1).
//
// error code that may be returned:
//  ACCOUNTQUALIFICATIONRESTRICTIONS = "AccountQualificationRestrictions"
//  AUTHFAILURE_CAMROLENAMEAUTHENTICATEFAILED = "AuthFailure.CamRoleNameAuthenticateFailed"
//  FAILEDOPERATION_DISASTERRECOVERGROUPNOTFOUND = "FailedOperation.DisasterRecoverGroupNotFound"
//  FAILEDOPERATION_ILLEGALTAGKEY = "FailedOperation.IllegalTagKey"
//  FAILEDOPERATION_ILLEGALTAGVALUE = "FailedOperation.IllegalTagValue"
//  FAILEDOPERATION_INQUIRYPRICEFAILED = "FailedOperation.InquiryPriceFailed"
//  FAILEDOPERATION_SNAPSHOTSIZELARGERTHANDATASIZE = "FailedOperation.SnapshotSizeLargerThanDataSize"
//  FAILEDOPERATION_TAGKEYRESERVED = "FailedOperation.TagKeyReserved"
//  FAILEDOPERATION_TATAGENTNOTSUPPORT = "FailedOperation.TatAgentNotSupport"
//  INSTANCESQUOTALIMITEXCEEDED = "InstancesQuotaLimitExceeded"
//  INTERNALERROR_TRADEUNKNOWNERROR = "InternalError.TradeUnknownError"
//  INVALIDCLIENTTOKEN_TOOLONG = "InvalidClientToken.TooLong"
//  INVALIDHOSTID_MALFORMED = "InvalidHostId.Malformed"
//  INVALIDHOSTID_NOTFOUND = "InvalidHostId.NotFound"
//  INVALIDIMAGEID_MALFORMED = "InvalidImageId.Malformed"
//  INVALIDIMAGEID_NOTFOUND = "InvalidImageId.NotFound"
//  INVALIDINSTANCENAME_TOOLONG = "InvalidInstanceName.TooLong"
//  INVALIDINSTANCETYPE_MALFORMED = "InvalidInstanceType.Malformed"
//  INVALIDPARAMETER_EDGEZONEMISSINTERNETACCESSIBLE = "InvalidParameter.EdgeZoneMissInternetAccessible"
//  INVALIDPARAMETER_HOSTIDCUSTOMIZEDINSTANCETYPENOTSUPPORT = "InvalidParameter.HostIdCustomizedInstanceTypeNotSupport"
//  INVALIDPARAMETER_HOSTIDINSTANCETYPENOTSUPPORT = "InvalidParameter.HostIdInstanceTypeNotSupport"
//  INVALIDPARAMETER_HOSTIDSTANDARDINSTANCETYPENOTSUPPORT = "InvalidParameter.HostIdStandardInstanceTypeNotSupport"
//  INVALIDPARAMETER_INSTANCEIMAGENOTSUPPORT = "InvalidParameter.InstanceImageNotSupport"
//  INVALIDPARAMETER_INSTANCETYPESUPPORTEDHOSTNOTFOUND = "InvalidParameter.InstanceTypeSupportedHostNotFound"
//  INVALIDPARAMETER_INTERNETACCESSIBLENOTSUPPORTED = "InvalidParameter.InternetAccessibleNotSupported"
//  INVALIDPARAMETER_ONLYSUPPORTFOREDGEZONE = "InvalidParameter.OnlySupportForEdgeZone"
//  INVALIDPARAMETER_SNAPSHOTNOTFOUND = "InvalidParameter.SnapshotNotFound"
//  INVALIDPARAMETER_SPECIALPARAMETERFORSPECIALACCOUNT = "InvalidParameter.SpecialParameterForSpecialAccount"
//  INVALIDPARAMETERCOMBINATION = "InvalidParameterCombination"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_BANDWIDTHPACKAGEIDMALFORMED = "InvalidParameterValue.BandwidthPackageIdMalformed"
//  INVALIDPARAMETERVALUE_BANDWIDTHPACKAGEIDNOTFOUND = "InvalidParameterValue.BandwidthPackageIdNotFound"
//  INVALIDPARAMETERVALUE_BANDWIDTHPACKAGEISPNOTMATCH = "InvalidParameterValue.BandwidthPackageIspNotMatch"
//  INVALIDPARAMETERVALUE_BANDWIDTHPACKAGEZONENOTMATCH = "InvalidParameterValue.BandwidthPackageZoneNotMatch"
//  INVALIDPARAMETERVALUE_CLOUDSSDDATADISKSIZETOOSMALL = "InvalidParameterValue.CloudSsdDataDiskSizeTooSmall"
//  INVALIDPARAMETERVALUE_DUPLICATETAGS = "InvalidParameterValue.DuplicateTags"
//  INVALIDPARAMETERVALUE_INSTANCETYPENOTFOUND = "InvalidParameterValue.InstanceTypeNotFound"
//  INVALIDPARAMETERVALUE_INSTANCETYPENOTSUPPORTHPCCLUSTER = "InvalidParameterValue.InstanceTypeNotSupportHpcCluster"
//  INVALIDPARAMETERVALUE_INSTANCETYPEREQUIREDHPCCLUSTER = "InvalidParameterValue.InstanceTypeRequiredHpcCluster"
//  INVALIDPARAMETERVALUE_INSUFFICIENTPRICE = "InvalidParameterValue.InsufficientPrice"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEFORMAT = "InvalidParameterValue.InvalidImageFormat"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEID = "InvalidParameterValue.InvalidImageId"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEOSNAME = "InvalidParameterValue.InvalidImageOsName"
//  INVALIDPARAMETERVALUE_INVALIDIMAGESTATE = "InvalidParameterValue.InvalidImageState"
//  INVALIDPARAMETERVALUE_INVALIDTIMEFORMAT = "InvalidParameterValue.InvalidTimeFormat"
//  INVALIDPARAMETERVALUE_INVALIDUSERDATAFORMAT = "InvalidParameterValue.InvalidUserDataFormat"
//  INVALIDPARAMETERVALUE_ISPNOTSUPPORTFOREDGEZONE = "InvalidParameterValue.IspNotSupportForEdgeZone"
//  INVALIDPARAMETERVALUE_ISPVALUEREPEATED = "InvalidParameterValue.IspValueRepeated"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_RANGE = "InvalidParameterValue.Range"
//  INVALIDPARAMETERVALUE_REQUIREDLOCATIONIMAGE = "InvalidParameterValue.RequiredLocationImage"
//  INVALIDPARAMETERVALUE_SNAPSHOTIDMALFORMED = "InvalidParameterValue.SnapshotIdMalformed"
//  INVALIDPARAMETERVALUE_TAGKEYNOTFOUND = "InvalidParameterValue.TagKeyNotFound"
//  INVALIDPARAMETERVALUE_TAGQUOTALIMITEXCEEDED = "InvalidParameterValue.TagQuotaLimitExceeded"
//  INVALIDPARAMETERVALUE_ZONENOTSUPPORTED = "InvalidParameterValue.ZoneNotSupported"
//  INVALIDPASSWORD = "InvalidPassword"
//  INVALIDPERIOD = "InvalidPeriod"
//  INVALIDPERMISSION = "InvalidPermission"
//  INVALIDPROJECTID_NOTFOUND = "InvalidProjectId.NotFound"
//  INVALIDZONE_MISMATCHREGION = "InvalidZone.MismatchRegion"
//  LIMITEXCEEDED_DISASTERRECOVERGROUP = "LimitExceeded.DisasterRecoverGroup"
//  LIMITEXCEEDED_INSTANCEENINUMLIMIT = "LimitExceeded.InstanceEniNumLimit"
//  MISSINGPARAMETER = "MissingParameter"
//  MISSINGPARAMETER_MONITORSERVICE = "MissingParameter.MonitorService"
//  RESOURCEINSUFFICIENT_CLOUDDISKUNAVAILABLE = "ResourceInsufficient.CloudDiskUnavailable"
//  RESOURCEINSUFFICIENT_DISASTERRECOVERGROUPCVMQUOTA = "ResourceInsufficient.DisasterRecoverGroupCvmQuota"
//  RESOURCENOTFOUND_HPCCLUSTER = "ResourceNotFound.HpcCluster"
//  RESOURCENOTFOUND_NODEFAULTCBS = "ResourceNotFound.NoDefaultCbs"
//  RESOURCENOTFOUND_NODEFAULTCBSWITHREASON = "ResourceNotFound.NoDefaultCbsWithReason"
//  RESOURCEUNAVAILABLE_INSTANCETYPE = "ResourceUnavailable.InstanceType"
//  RESOURCESSOLDOUT_SPECIFIEDINSTANCETYPE = "ResourcesSoldOut.SpecifiedInstanceType"
//  UNSUPPORTEDOPERATION_BANDWIDTHPACKAGEIDNOTSUPPORTED = "UnsupportedOperation.BandwidthPackageIdNotSupported"
//  UNSUPPORTEDOPERATION_INSTANCETYPENOTSUPPORTGRIDLICENCE = "UnsupportedOperation.InstanceTypeNotSupportGridLicence"
//  UNSUPPORTEDOPERATION_INSTANCETYPENOTSUPPORTJUMBOFRAME = "UnsupportedOperation.InstanceTypeNotSupportJumboFrame"
//  UNSUPPORTEDOPERATION_INVALIDDISK = "UnsupportedOperation.InvalidDisk"
//  UNSUPPORTEDOPERATION_INVALIDREGIONDISKENCRYPT = "UnsupportedOperation.InvalidRegionDiskEncrypt"
//  UNSUPPORTEDOPERATION_NOINSTANCETYPESUPPORTSPOT = "UnsupportedOperation.NoInstanceTypeSupportSpot"
//  UNSUPPORTEDOPERATION_NOTSUPPORTIMPORTINSTANCESACTIONTIMER = "UnsupportedOperation.NotSupportImportInstancesActionTimer"
//  UNSUPPORTEDOPERATION_ONLYFORPREPAIDACCOUNT = "UnsupportedOperation.OnlyForPrepaidAccount"
//  UNSUPPORTEDOPERATION_SPOTUNSUPPORTEDREGION = "UnsupportedOperation.SpotUnsupportedRegion"
//  UNSUPPORTEDOPERATION_UNDERWRITINGINSTANCETYPEONLYSUPPORTAUTORENEW = "UnsupportedOperation.UnderwritingInstanceTypeOnlySupportAutoRenew"
//  UNSUPPORTEDOPERATION_UNSUPPORTEDINTERNATIONALUSER = "UnsupportedOperation.UnsupportedInternationalUser"
func (c *Client) InquiryPriceRunInstances(request *InquiryPriceRunInstancesRequest) (response *InquiryPriceRunInstancesResponse, err error) {
    return c.InquiryPriceRunInstancesWithContext(context.Background(), request)
}

// InquiryPriceRunInstances
// This API is used to query the price of creating instances. You can only use this API for instances whose configuration is within the purchase limit. For more information, see [RunInstances](https://intl.cloud.tencent.com/document/api/213/15730?from_cn_redirect=1).
//
// error code that may be returned:
//  ACCOUNTQUALIFICATIONRESTRICTIONS = "AccountQualificationRestrictions"
//  AUTHFAILURE_CAMROLENAMEAUTHENTICATEFAILED = "AuthFailure.CamRoleNameAuthenticateFailed"
//  FAILEDOPERATION_DISASTERRECOVERGROUPNOTFOUND = "FailedOperation.DisasterRecoverGroupNotFound"
//  FAILEDOPERATION_ILLEGALTAGKEY = "FailedOperation.IllegalTagKey"
//  FAILEDOPERATION_ILLEGALTAGVALUE = "FailedOperation.IllegalTagValue"
//  FAILEDOPERATION_INQUIRYPRICEFAILED = "FailedOperation.InquiryPriceFailed"
//  FAILEDOPERATION_SNAPSHOTSIZELARGERTHANDATASIZE = "FailedOperation.SnapshotSizeLargerThanDataSize"
//  FAILEDOPERATION_TAGKEYRESERVED = "FailedOperation.TagKeyReserved"
//  FAILEDOPERATION_TATAGENTNOTSUPPORT = "FailedOperation.TatAgentNotSupport"
//  INSTANCESQUOTALIMITEXCEEDED = "InstancesQuotaLimitExceeded"
//  INTERNALERROR_TRADEUNKNOWNERROR = "InternalError.TradeUnknownError"
//  INVALIDCLIENTTOKEN_TOOLONG = "InvalidClientToken.TooLong"
//  INVALIDHOSTID_MALFORMED = "InvalidHostId.Malformed"
//  INVALIDHOSTID_NOTFOUND = "InvalidHostId.NotFound"
//  INVALIDIMAGEID_MALFORMED = "InvalidImageId.Malformed"
//  INVALIDIMAGEID_NOTFOUND = "InvalidImageId.NotFound"
//  INVALIDINSTANCENAME_TOOLONG = "InvalidInstanceName.TooLong"
//  INVALIDINSTANCETYPE_MALFORMED = "InvalidInstanceType.Malformed"
//  INVALIDPARAMETER_EDGEZONEMISSINTERNETACCESSIBLE = "InvalidParameter.EdgeZoneMissInternetAccessible"
//  INVALIDPARAMETER_HOSTIDCUSTOMIZEDINSTANCETYPENOTSUPPORT = "InvalidParameter.HostIdCustomizedInstanceTypeNotSupport"
//  INVALIDPARAMETER_HOSTIDINSTANCETYPENOTSUPPORT = "InvalidParameter.HostIdInstanceTypeNotSupport"
//  INVALIDPARAMETER_HOSTIDSTANDARDINSTANCETYPENOTSUPPORT = "InvalidParameter.HostIdStandardInstanceTypeNotSupport"
//  INVALIDPARAMETER_INSTANCEIMAGENOTSUPPORT = "InvalidParameter.InstanceImageNotSupport"
//  INVALIDPARAMETER_INSTANCETYPESUPPORTEDHOSTNOTFOUND = "InvalidParameter.InstanceTypeSupportedHostNotFound"
//  INVALIDPARAMETER_INTERNETACCESSIBLENOTSUPPORTED = "InvalidParameter.InternetAccessibleNotSupported"
//  INVALIDPARAMETER_ONLYSUPPORTFOREDGEZONE = "InvalidParameter.OnlySupportForEdgeZone"
//  INVALIDPARAMETER_SNAPSHOTNOTFOUND = "InvalidParameter.SnapshotNotFound"
//  INVALIDPARAMETER_SPECIALPARAMETERFORSPECIALACCOUNT = "InvalidParameter.SpecialParameterForSpecialAccount"
//  INVALIDPARAMETERCOMBINATION = "InvalidParameterCombination"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_BANDWIDTHPACKAGEIDMALFORMED = "InvalidParameterValue.BandwidthPackageIdMalformed"
//  INVALIDPARAMETERVALUE_BANDWIDTHPACKAGEIDNOTFOUND = "InvalidParameterValue.BandwidthPackageIdNotFound"
//  INVALIDPARAMETERVALUE_BANDWIDTHPACKAGEISPNOTMATCH = "InvalidParameterValue.BandwidthPackageIspNotMatch"
//  INVALIDPARAMETERVALUE_BANDWIDTHPACKAGEZONENOTMATCH = "InvalidParameterValue.BandwidthPackageZoneNotMatch"
//  INVALIDPARAMETERVALUE_CLOUDSSDDATADISKSIZETOOSMALL = "InvalidParameterValue.CloudSsdDataDiskSizeTooSmall"
//  INVALIDPARAMETERVALUE_DUPLICATETAGS = "InvalidParameterValue.DuplicateTags"
//  INVALIDPARAMETERVALUE_INSTANCETYPENOTFOUND = "InvalidParameterValue.InstanceTypeNotFound"
//  INVALIDPARAMETERVALUE_INSTANCETYPENOTSUPPORTHPCCLUSTER = "InvalidParameterValue.InstanceTypeNotSupportHpcCluster"
//  INVALIDPARAMETERVALUE_INSTANCETYPEREQUIREDHPCCLUSTER = "InvalidParameterValue.InstanceTypeRequiredHpcCluster"
//  INVALIDPARAMETERVALUE_INSUFFICIENTPRICE = "InvalidParameterValue.InsufficientPrice"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEFORMAT = "InvalidParameterValue.InvalidImageFormat"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEID = "InvalidParameterValue.InvalidImageId"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEOSNAME = "InvalidParameterValue.InvalidImageOsName"
//  INVALIDPARAMETERVALUE_INVALIDIMAGESTATE = "InvalidParameterValue.InvalidImageState"
//  INVALIDPARAMETERVALUE_INVALIDTIMEFORMAT = "InvalidParameterValue.InvalidTimeFormat"
//  INVALIDPARAMETERVALUE_INVALIDUSERDATAFORMAT = "InvalidParameterValue.InvalidUserDataFormat"
//  INVALIDPARAMETERVALUE_ISPNOTSUPPORTFOREDGEZONE = "InvalidParameterValue.IspNotSupportForEdgeZone"
//  INVALIDPARAMETERVALUE_ISPVALUEREPEATED = "InvalidParameterValue.IspValueRepeated"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_RANGE = "InvalidParameterValue.Range"
//  INVALIDPARAMETERVALUE_REQUIREDLOCATIONIMAGE = "InvalidParameterValue.RequiredLocationImage"
//  INVALIDPARAMETERVALUE_SNAPSHOTIDMALFORMED = "InvalidParameterValue.SnapshotIdMalformed"
//  INVALIDPARAMETERVALUE_TAGKEYNOTFOUND = "InvalidParameterValue.TagKeyNotFound"
//  INVALIDPARAMETERVALUE_TAGQUOTALIMITEXCEEDED = "InvalidParameterValue.TagQuotaLimitExceeded"
//  INVALIDPARAMETERVALUE_ZONENOTSUPPORTED = "InvalidParameterValue.ZoneNotSupported"
//  INVALIDPASSWORD = "InvalidPassword"
//  INVALIDPERIOD = "InvalidPeriod"
//  INVALIDPERMISSION = "InvalidPermission"
//  INVALIDPROJECTID_NOTFOUND = "InvalidProjectId.NotFound"
//  INVALIDZONE_MISMATCHREGION = "InvalidZone.MismatchRegion"
//  LIMITEXCEEDED_DISASTERRECOVERGROUP = "LimitExceeded.DisasterRecoverGroup"
//  LIMITEXCEEDED_INSTANCEENINUMLIMIT = "LimitExceeded.InstanceEniNumLimit"
//  MISSINGPARAMETER = "MissingParameter"
//  MISSINGPARAMETER_MONITORSERVICE = "MissingParameter.MonitorService"
//  RESOURCEINSUFFICIENT_CLOUDDISKUNAVAILABLE = "ResourceInsufficient.CloudDiskUnavailable"
//  RESOURCEINSUFFICIENT_DISASTERRECOVERGROUPCVMQUOTA = "ResourceInsufficient.DisasterRecoverGroupCvmQuota"
//  RESOURCENOTFOUND_HPCCLUSTER = "ResourceNotFound.HpcCluster"
//  RESOURCENOTFOUND_NODEFAULTCBS = "ResourceNotFound.NoDefaultCbs"
//  RESOURCENOTFOUND_NODEFAULTCBSWITHREASON = "ResourceNotFound.NoDefaultCbsWithReason"
//  RESOURCEUNAVAILABLE_INSTANCETYPE = "ResourceUnavailable.InstanceType"
//  RESOURCESSOLDOUT_SPECIFIEDINSTANCETYPE = "ResourcesSoldOut.SpecifiedInstanceType"
//  UNSUPPORTEDOPERATION_BANDWIDTHPACKAGEIDNOTSUPPORTED = "UnsupportedOperation.BandwidthPackageIdNotSupported"
//  UNSUPPORTEDOPERATION_INSTANCETYPENOTSUPPORTGRIDLICENCE = "UnsupportedOperation.InstanceTypeNotSupportGridLicence"
//  UNSUPPORTEDOPERATION_INSTANCETYPENOTSUPPORTJUMBOFRAME = "UnsupportedOperation.InstanceTypeNotSupportJumboFrame"
//  UNSUPPORTEDOPERATION_INVALIDDISK = "UnsupportedOperation.InvalidDisk"
//  UNSUPPORTEDOPERATION_INVALIDREGIONDISKENCRYPT = "UnsupportedOperation.InvalidRegionDiskEncrypt"
//  UNSUPPORTEDOPERATION_NOINSTANCETYPESUPPORTSPOT = "UnsupportedOperation.NoInstanceTypeSupportSpot"
//  UNSUPPORTEDOPERATION_NOTSUPPORTIMPORTINSTANCESACTIONTIMER = "UnsupportedOperation.NotSupportImportInstancesActionTimer"
//  UNSUPPORTEDOPERATION_ONLYFORPREPAIDACCOUNT = "UnsupportedOperation.OnlyForPrepaidAccount"
//  UNSUPPORTEDOPERATION_SPOTUNSUPPORTEDREGION = "UnsupportedOperation.SpotUnsupportedRegion"
//  UNSUPPORTEDOPERATION_UNDERWRITINGINSTANCETYPEONLYSUPPORTAUTORENEW = "UnsupportedOperation.UnderwritingInstanceTypeOnlySupportAutoRenew"
//  UNSUPPORTEDOPERATION_UNSUPPORTEDINTERNATIONALUSER = "UnsupportedOperation.UnsupportedInternationalUser"
func (c *Client) InquiryPriceRunInstancesWithContext(ctx context.Context, request *InquiryPriceRunInstancesRequest) (response *InquiryPriceRunInstancesResponse, err error) {
    if request == nil {
        request = NewInquiryPriceRunInstancesRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "InquiryPriceRunInstances")
    
    if c.GetCredential() == nil {
        return nil, errors.New("InquiryPriceRunInstances require credential")
    }

    request.SetContext(ctx)
    
    response = NewInquiryPriceRunInstancesResponse()
    err = c.Send(request, response)
    return
}

func NewModifyChcAttributeRequest() (request *ModifyChcAttributeRequest) {
    request = &ModifyChcAttributeRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "ModifyChcAttribute")
    
    
    return
}

func NewModifyChcAttributeResponse() (response *ModifyChcAttributeResponse) {
    response = &ModifyChcAttributeResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// ModifyChcAttribute
// This API is used to modify the CHC host attributes.
//
// error code that may be returned:
//  ACCOUNTQUALIFICATIONRESTRICTIONS = "AccountQualificationRestrictions"
//  INVALIDHOST_NOTSUPPORTED = "InvalidHost.NotSupported"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_CHCHOSTSNOTFOUND = "InvalidParameterValue.ChcHostsNotFound"
//  INVALIDPARAMETERVALUE_CHCNETWORKEMPTY = "InvalidParameterValue.ChcNetworkEmpty"
//  INVALIDPARAMETERVALUE_INCORRECTFORMAT = "InvalidParameterValue.IncorrectFormat"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_NOTEMPTY = "InvalidParameterValue.NotEmpty"
//  INVALIDPASSWORD = "InvalidPassword"
//  INVALIDSECURITYGROUPID_NOTFOUND = "InvalidSecurityGroupId.NotFound"
func (c *Client) ModifyChcAttribute(request *ModifyChcAttributeRequest) (response *ModifyChcAttributeResponse, err error) {
    return c.ModifyChcAttributeWithContext(context.Background(), request)
}

// ModifyChcAttribute
// This API is used to modify the CHC host attributes.
//
// error code that may be returned:
//  ACCOUNTQUALIFICATIONRESTRICTIONS = "AccountQualificationRestrictions"
//  INVALIDHOST_NOTSUPPORTED = "InvalidHost.NotSupported"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_CHCHOSTSNOTFOUND = "InvalidParameterValue.ChcHostsNotFound"
//  INVALIDPARAMETERVALUE_CHCNETWORKEMPTY = "InvalidParameterValue.ChcNetworkEmpty"
//  INVALIDPARAMETERVALUE_INCORRECTFORMAT = "InvalidParameterValue.IncorrectFormat"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_NOTEMPTY = "InvalidParameterValue.NotEmpty"
//  INVALIDPASSWORD = "InvalidPassword"
//  INVALIDSECURITYGROUPID_NOTFOUND = "InvalidSecurityGroupId.NotFound"
func (c *Client) ModifyChcAttributeWithContext(ctx context.Context, request *ModifyChcAttributeRequest) (response *ModifyChcAttributeResponse, err error) {
    if request == nil {
        request = NewModifyChcAttributeRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "ModifyChcAttribute")
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyChcAttribute require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyChcAttributeResponse()
    err = c.Send(request, response)
    return
}

func NewModifyDisasterRecoverGroupAttributeRequest() (request *ModifyDisasterRecoverGroupAttributeRequest) {
    request = &ModifyDisasterRecoverGroupAttributeRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "ModifyDisasterRecoverGroupAttribute")
    
    
    return
}

func NewModifyDisasterRecoverGroupAttributeResponse() (response *ModifyDisasterRecoverGroupAttributeResponse) {
    response = &ModifyDisasterRecoverGroupAttributeResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// ModifyDisasterRecoverGroupAttribute
// This API is used to modify the attributes of [spread placement groups](https://intl.cloud.tencent.com/document/product/213/15486?from_cn_redirect=1).
//
// error code that may be returned:
//  INVALIDPARAMETERVALUE_TOOLONG = "InvalidParameterValue.TooLong"
//  RESOURCEINSUFFICIENT_DISASTERRECOVERGROUPCVMQUOTA = "ResourceInsufficient.DisasterRecoverGroupCvmQuota"
//  RESOURCENOTFOUND_INVALIDPLACEMENTSET = "ResourceNotFound.InvalidPlacementSet"
func (c *Client) ModifyDisasterRecoverGroupAttribute(request *ModifyDisasterRecoverGroupAttributeRequest) (response *ModifyDisasterRecoverGroupAttributeResponse, err error) {
    return c.ModifyDisasterRecoverGroupAttributeWithContext(context.Background(), request)
}

// ModifyDisasterRecoverGroupAttribute
// This API is used to modify the attributes of [spread placement groups](https://intl.cloud.tencent.com/document/product/213/15486?from_cn_redirect=1).
//
// error code that may be returned:
//  INVALIDPARAMETERVALUE_TOOLONG = "InvalidParameterValue.TooLong"
//  RESOURCEINSUFFICIENT_DISASTERRECOVERGROUPCVMQUOTA = "ResourceInsufficient.DisasterRecoverGroupCvmQuota"
//  RESOURCENOTFOUND_INVALIDPLACEMENTSET = "ResourceNotFound.InvalidPlacementSet"
func (c *Client) ModifyDisasterRecoverGroupAttributeWithContext(ctx context.Context, request *ModifyDisasterRecoverGroupAttributeRequest) (response *ModifyDisasterRecoverGroupAttributeResponse, err error) {
    if request == nil {
        request = NewModifyDisasterRecoverGroupAttributeRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "ModifyDisasterRecoverGroupAttribute")
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyDisasterRecoverGroupAttribute require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyDisasterRecoverGroupAttributeResponse()
    err = c.Send(request, response)
    return
}

func NewModifyHostsAttributeRequest() (request *ModifyHostsAttributeRequest) {
    request = &ModifyHostsAttributeRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "ModifyHostsAttribute")
    
    
    return
}

func NewModifyHostsAttributeResponse() (response *ModifyHostsAttributeResponse) {
    response = &ModifyHostsAttributeResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// ModifyHostsAttribute
// This API is used to modify the attributes of a CDH instance, such as instance name and renewal flag. One of the two parameters, HostName and RenewFlag, must be set, but you cannot set both of them at the same time.
//
// error code that may be returned:
//  INVALIDHOST_NOTSUPPORTED = "InvalidHost.NotSupported"
//  INVALIDHOSTID_MALFORMED = "InvalidHostId.Malformed"
//  INVALIDHOSTID_NOTFOUND = "InvalidHostId.NotFound"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
func (c *Client) ModifyHostsAttribute(request *ModifyHostsAttributeRequest) (response *ModifyHostsAttributeResponse, err error) {
    return c.ModifyHostsAttributeWithContext(context.Background(), request)
}

// ModifyHostsAttribute
// This API is used to modify the attributes of a CDH instance, such as instance name and renewal flag. One of the two parameters, HostName and RenewFlag, must be set, but you cannot set both of them at the same time.
//
// error code that may be returned:
//  INVALIDHOST_NOTSUPPORTED = "InvalidHost.NotSupported"
//  INVALIDHOSTID_MALFORMED = "InvalidHostId.Malformed"
//  INVALIDHOSTID_NOTFOUND = "InvalidHostId.NotFound"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
func (c *Client) ModifyHostsAttributeWithContext(ctx context.Context, request *ModifyHostsAttributeRequest) (response *ModifyHostsAttributeResponse, err error) {
    if request == nil {
        request = NewModifyHostsAttributeRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "ModifyHostsAttribute")
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyHostsAttribute require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyHostsAttributeResponse()
    err = c.Send(request, response)
    return
}

func NewModifyImageAttributeRequest() (request *ModifyImageAttributeRequest) {
    request = &ModifyImageAttributeRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "ModifyImageAttribute")
    
    
    return
}

func NewModifyImageAttributeResponse() (response *ModifyImageAttributeResponse) {
    response = &ModifyImageAttributeResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// ModifyImageAttribute
// This API is used to modify image attributes.
//
// 
//
// * Attributes of shared images cannot be modified.
//
// error code that may be returned:
//  INVALIDIMAGEID_INCORRECTSTATE = "InvalidImageId.IncorrectState"
//  INVALIDIMAGEID_MALFORMED = "InvalidImageId.Malformed"
//  INVALIDIMAGEID_NOTFOUND = "InvalidImageId.NotFound"
//  INVALIDIMAGENAME_DUPLICATE = "InvalidImageName.Duplicate"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEFAMILY = "InvalidParameterValue.InvalidImageFamily"
//  INVALIDPARAMETERVALUE_TOOLARGE = "InvalidParameterValue.TooLarge"
//  UNAUTHORIZEDOPERATION_INVALIDTOKEN = "UnauthorizedOperation.InvalidToken"
//  UNSUPPORTEDOPERATION_SHAREDIMAGEMODIFYUNSUPPORTED = "UnsupportedOperation.SharedImageModifyUnsupported"
func (c *Client) ModifyImageAttribute(request *ModifyImageAttributeRequest) (response *ModifyImageAttributeResponse, err error) {
    return c.ModifyImageAttributeWithContext(context.Background(), request)
}

// ModifyImageAttribute
// This API is used to modify image attributes.
//
// 
//
// * Attributes of shared images cannot be modified.
//
// error code that may be returned:
//  INVALIDIMAGEID_INCORRECTSTATE = "InvalidImageId.IncorrectState"
//  INVALIDIMAGEID_MALFORMED = "InvalidImageId.Malformed"
//  INVALIDIMAGEID_NOTFOUND = "InvalidImageId.NotFound"
//  INVALIDIMAGENAME_DUPLICATE = "InvalidImageName.Duplicate"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEFAMILY = "InvalidParameterValue.InvalidImageFamily"
//  INVALIDPARAMETERVALUE_TOOLARGE = "InvalidParameterValue.TooLarge"
//  UNAUTHORIZEDOPERATION_INVALIDTOKEN = "UnauthorizedOperation.InvalidToken"
//  UNSUPPORTEDOPERATION_SHAREDIMAGEMODIFYUNSUPPORTED = "UnsupportedOperation.SharedImageModifyUnsupported"
func (c *Client) ModifyImageAttributeWithContext(ctx context.Context, request *ModifyImageAttributeRequest) (response *ModifyImageAttributeResponse, err error) {
    if request == nil {
        request = NewModifyImageAttributeRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "ModifyImageAttribute")
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyImageAttribute require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyImageAttributeResponse()
    err = c.Send(request, response)
    return
}

func NewModifyImageSharePermissionRequest() (request *ModifyImageSharePermissionRequest) {
    request = &ModifyImageSharePermissionRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "ModifyImageSharePermission")
    
    
    return
}

func NewModifyImageSharePermissionResponse() (response *ModifyImageSharePermissionResponse) {
    response = &ModifyImageSharePermissionResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// ModifyImageSharePermission
// This API is used to modify image sharing information.
//
// 
//
// * The account with shared image access can use the image to create instances.
//
// * Each custom image can be shared with a maximum of 500 accounts.
//
// * Shared images cannot have their names or description changed. They can only be used to create instances.
//
// * Sharing is only supported within the same region as the recipient's account.
//
// error code that may be returned:
//  FAILEDOPERATION_ACCOUNTALREADYEXISTS = "FailedOperation.AccountAlreadyExists"
//  FAILEDOPERATION_ACCOUNTISYOURSELF = "FailedOperation.AccountIsYourSelf"
//  FAILEDOPERATION_BYOLIMAGESHAREFAILED = "FailedOperation.BYOLImageShareFailed"
//  FAILEDOPERATION_NOTMASTERACCOUNT = "FailedOperation.NotMasterAccount"
//  FAILEDOPERATION_QIMAGESHAREFAILED = "FailedOperation.QImageShareFailed"
//  FAILEDOPERATION_RIMAGESHAREFAILED = "FailedOperation.RImageShareFailed"
//  INVALIDACCOUNTID_NOTFOUND = "InvalidAccountId.NotFound"
//  INVALIDACCOUNTIS_YOURSELF = "InvalidAccountIs.YourSelf"
//  INVALIDIMAGEID_MALFORMED = "InvalidImageId.Malformed"
//  INVALIDIMAGEID_NOTFOUND = "InvalidImageId.NotFound"
//  INVALIDPARAMETER_INSTANCEIMAGENOTSUPPORT = "InvalidParameter.InstanceImageNotSupport"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEID = "InvalidParameterValue.InvalidImageId"
//  INVALIDPARAMETERVALUE_INVALIDIMAGESTATE = "InvalidParameterValue.InvalidImageState"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  OVERQUOTA = "OverQuota"
//  UNAUTHORIZEDOPERATION_IMAGENOTBELONGTOACCOUNT = "UnauthorizedOperation.ImageNotBelongToAccount"
//  UNAUTHORIZEDOPERATION_INVALIDTOKEN = "UnauthorizedOperation.InvalidToken"
//  UNSUPPORTEDOPERATION_ENCRYPTEDIMAGESNOTSUPPORTED = "UnsupportedOperation.EncryptedImagesNotSupported"
//  UNSUPPORTEDOPERATION_LOCATIONIMAGENOTSUPPORTED = "UnsupportedOperation.LocationImageNotSupported"
func (c *Client) ModifyImageSharePermission(request *ModifyImageSharePermissionRequest) (response *ModifyImageSharePermissionResponse, err error) {
    return c.ModifyImageSharePermissionWithContext(context.Background(), request)
}

// ModifyImageSharePermission
// This API is used to modify image sharing information.
//
// 
//
// * The account with shared image access can use the image to create instances.
//
// * Each custom image can be shared with a maximum of 500 accounts.
//
// * Shared images cannot have their names or description changed. They can only be used to create instances.
//
// * Sharing is only supported within the same region as the recipient's account.
//
// error code that may be returned:
//  FAILEDOPERATION_ACCOUNTALREADYEXISTS = "FailedOperation.AccountAlreadyExists"
//  FAILEDOPERATION_ACCOUNTISYOURSELF = "FailedOperation.AccountIsYourSelf"
//  FAILEDOPERATION_BYOLIMAGESHAREFAILED = "FailedOperation.BYOLImageShareFailed"
//  FAILEDOPERATION_NOTMASTERACCOUNT = "FailedOperation.NotMasterAccount"
//  FAILEDOPERATION_QIMAGESHAREFAILED = "FailedOperation.QImageShareFailed"
//  FAILEDOPERATION_RIMAGESHAREFAILED = "FailedOperation.RImageShareFailed"
//  INVALIDACCOUNTID_NOTFOUND = "InvalidAccountId.NotFound"
//  INVALIDACCOUNTIS_YOURSELF = "InvalidAccountIs.YourSelf"
//  INVALIDIMAGEID_MALFORMED = "InvalidImageId.Malformed"
//  INVALIDIMAGEID_NOTFOUND = "InvalidImageId.NotFound"
//  INVALIDPARAMETER_INSTANCEIMAGENOTSUPPORT = "InvalidParameter.InstanceImageNotSupport"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEID = "InvalidParameterValue.InvalidImageId"
//  INVALIDPARAMETERVALUE_INVALIDIMAGESTATE = "InvalidParameterValue.InvalidImageState"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  OVERQUOTA = "OverQuota"
//  UNAUTHORIZEDOPERATION_IMAGENOTBELONGTOACCOUNT = "UnauthorizedOperation.ImageNotBelongToAccount"
//  UNAUTHORIZEDOPERATION_INVALIDTOKEN = "UnauthorizedOperation.InvalidToken"
//  UNSUPPORTEDOPERATION_ENCRYPTEDIMAGESNOTSUPPORTED = "UnsupportedOperation.EncryptedImagesNotSupported"
//  UNSUPPORTEDOPERATION_LOCATIONIMAGENOTSUPPORTED = "UnsupportedOperation.LocationImageNotSupported"
func (c *Client) ModifyImageSharePermissionWithContext(ctx context.Context, request *ModifyImageSharePermissionRequest) (response *ModifyImageSharePermissionResponse, err error) {
    if request == nil {
        request = NewModifyImageSharePermissionRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "ModifyImageSharePermission")
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyImageSharePermission require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyImageSharePermissionResponse()
    err = c.Send(request, response)
    return
}

func NewModifyInstancesAttributeRequest() (request *ModifyInstancesAttributeRequest) {
    request = &ModifyInstancesAttributeRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "ModifyInstancesAttribute")
    
    
    return
}

func NewModifyInstancesAttributeResponse() (response *ModifyInstancesAttributeResponse) {
    response = &ModifyInstancesAttributeResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// ModifyInstancesAttribute
// This API is used to modify instance attributes.
//
// 
//
// This API is used to modify one attribute of the instance per request. The attribute must be specified.
//
// The instance name is only for user convenience in management. Tencent Cloud does not use this name as the basis for online support or to perform instance management operations.
//
// This API is used to support batch operations. The maximum of 100 batch instances per request is supported.
//
// This API is used to modify the security group association. The originally associated security group of the instance will be unbound.
//
// * You can query the result of the instance operation by calling the API [DescribeInstances](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1#.E7.A4.BA.E4.BE.8B3-.E6.9F.A5.E8.AF.A2.E5.AE.9E.E4.BE.8B.E7.9A.84.E6.9C.80.E6.96.B0.E6.93.8D.E4.BD.9C.E6.83.85.E5.86.B5). If the latest operation status (LatestOperationState) of the instance is **SUCCESS**, the operation is successful.
//
// This API is used to modify the hostname. The instance restarts immediately after hostname modification, and the new hostname takes effect after restart.
//
// error code that may be returned:
//  AUTHFAILURE_CAMROLENAMEAUTHENTICATEFAILED = "AuthFailure.CamRoleNameAuthenticateFailed"
//  FAILEDOPERATION_INVALIDINSTANCEAPPLICATIONROLE = "FailedOperation.InvalidInstanceApplicationRole"
//  FAILEDOPERATION_SECURITYGROUPACTIONFAILED = "FailedOperation.SecurityGroupActionFailed"
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDINSTANCENAME_TOOLONG = "InvalidInstanceName.TooLong"
//  INVALIDPARAMETER_HOSTNAMEILLEGAL = "InvalidParameter.HostNameIllegal"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_CAMROLENAMEMALFORMED = "InvalidParameterValue.CamRoleNameMalformed"
//  INVALIDPARAMETERVALUE_ILLEGALHOSTNAME = "InvalidParameterValue.IllegalHostName"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_INVALIDUSERDATAFORMAT = "InvalidParameterValue.InvalidUserDataFormat"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDSECURITYGROUPID_NOTFOUND = "InvalidSecurityGroupId.NotFound"
//  LIMITEXCEEDED_ASSOCIATEUSGLIMITEXCEEDED = "LimitExceeded.AssociateUSGLimitExceeded"
//  LIMITEXCEEDED_CVMSVIFSPERSECGROUPLIMITEXCEEDED = "LimitExceeded.CvmsVifsPerSecGroupLimitExceeded"
//  LIMITEXCEEDED_SINGLEUSGQUOTA = "LimitExceeded.SingleUSGQuota"
//  MISSINGPARAMETER = "MissingParameter"
//  MUTEXOPERATION_TASKRUNNING = "MutexOperation.TaskRunning"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  UNAUTHORIZEDOPERATION_MFANOTFOUND = "UnauthorizedOperation.MFANotFound"
//  UNSUPPORTEDOPERATION_INSTANCEREINSTALLFAILED = "UnsupportedOperation.InstanceReinstallFailed"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERRESCUEMODE = "UnsupportedOperation.InstanceStateEnterRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERSERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateEnterServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITRESCUEMODE = "UnsupportedOperation.InstanceStateExitRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITSERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateExitServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATEFREEZING = "UnsupportedOperation.InstanceStateFreezing"
//  UNSUPPORTEDOPERATION_INSTANCESTATEISOLATING = "UnsupportedOperation.InstanceStateIsolating"
//  UNSUPPORTEDOPERATION_INSTANCESTATEPENDING = "UnsupportedOperation.InstanceStatePending"
//  UNSUPPORTEDOPERATION_INSTANCESTATEREBOOTING = "UnsupportedOperation.InstanceStateRebooting"
//  UNSUPPORTEDOPERATION_INSTANCESTATERESCUEMODE = "UnsupportedOperation.InstanceStateRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATESERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATESHUTDOWN = "UnsupportedOperation.InstanceStateShutdown"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTARTING = "UnsupportedOperation.InstanceStateStarting"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPING = "UnsupportedOperation.InstanceStateStopping"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATED = "UnsupportedOperation.InstanceStateTerminated"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATING = "UnsupportedOperation.InstanceStateTerminating"
//  UNSUPPORTEDOPERATION_INSTANCETYPENOTSUPPORTHIGHDENSITYMODESETTING = "UnsupportedOperation.InstanceTypeNotSupportHighDensityModeSetting"
//  UNSUPPORTEDOPERATION_INSTANCETYPENOTSUPPORTJUMBOFRAME = "UnsupportedOperation.InstanceTypeNotSupportJumboFrame"
//  UNSUPPORTEDOPERATION_INSTANCESENABLEJUMBOWITHOUTREBOOT = "UnsupportedOperation.InstancesEnableJumboWithoutReboot"
//  UNSUPPORTEDOPERATION_INVALIDINSTANCENOTSUPPORTEDPROTECTEDINSTANCE = "UnsupportedOperation.InvalidInstanceNotSupportedProtectedInstance"
//  UNSUPPORTEDOPERATION_STOPPEDMODESTOPCHARGING = "UnsupportedOperation.StoppedModeStopCharging"
func (c *Client) ModifyInstancesAttribute(request *ModifyInstancesAttributeRequest) (response *ModifyInstancesAttributeResponse, err error) {
    return c.ModifyInstancesAttributeWithContext(context.Background(), request)
}

// ModifyInstancesAttribute
// This API is used to modify instance attributes.
//
// 
//
// This API is used to modify one attribute of the instance per request. The attribute must be specified.
//
// The instance name is only for user convenience in management. Tencent Cloud does not use this name as the basis for online support or to perform instance management operations.
//
// This API is used to support batch operations. The maximum of 100 batch instances per request is supported.
//
// This API is used to modify the security group association. The originally associated security group of the instance will be unbound.
//
// * You can query the result of the instance operation by calling the API [DescribeInstances](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1#.E7.A4.BA.E4.BE.8B3-.E6.9F.A5.E8.AF.A2.E5.AE.9E.E4.BE.8B.E7.9A.84.E6.9C.80.E6.96.B0.E6.93.8D.E4.BD.9C.E6.83.85.E5.86.B5). If the latest operation status (LatestOperationState) of the instance is **SUCCESS**, the operation is successful.
//
// This API is used to modify the hostname. The instance restarts immediately after hostname modification, and the new hostname takes effect after restart.
//
// error code that may be returned:
//  AUTHFAILURE_CAMROLENAMEAUTHENTICATEFAILED = "AuthFailure.CamRoleNameAuthenticateFailed"
//  FAILEDOPERATION_INVALIDINSTANCEAPPLICATIONROLE = "FailedOperation.InvalidInstanceApplicationRole"
//  FAILEDOPERATION_SECURITYGROUPACTIONFAILED = "FailedOperation.SecurityGroupActionFailed"
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDINSTANCENAME_TOOLONG = "InvalidInstanceName.TooLong"
//  INVALIDPARAMETER_HOSTNAMEILLEGAL = "InvalidParameter.HostNameIllegal"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_CAMROLENAMEMALFORMED = "InvalidParameterValue.CamRoleNameMalformed"
//  INVALIDPARAMETERVALUE_ILLEGALHOSTNAME = "InvalidParameterValue.IllegalHostName"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_INVALIDUSERDATAFORMAT = "InvalidParameterValue.InvalidUserDataFormat"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDSECURITYGROUPID_NOTFOUND = "InvalidSecurityGroupId.NotFound"
//  LIMITEXCEEDED_ASSOCIATEUSGLIMITEXCEEDED = "LimitExceeded.AssociateUSGLimitExceeded"
//  LIMITEXCEEDED_CVMSVIFSPERSECGROUPLIMITEXCEEDED = "LimitExceeded.CvmsVifsPerSecGroupLimitExceeded"
//  LIMITEXCEEDED_SINGLEUSGQUOTA = "LimitExceeded.SingleUSGQuota"
//  MISSINGPARAMETER = "MissingParameter"
//  MUTEXOPERATION_TASKRUNNING = "MutexOperation.TaskRunning"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  UNAUTHORIZEDOPERATION_MFANOTFOUND = "UnauthorizedOperation.MFANotFound"
//  UNSUPPORTEDOPERATION_INSTANCEREINSTALLFAILED = "UnsupportedOperation.InstanceReinstallFailed"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERRESCUEMODE = "UnsupportedOperation.InstanceStateEnterRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERSERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateEnterServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITRESCUEMODE = "UnsupportedOperation.InstanceStateExitRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITSERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateExitServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATEFREEZING = "UnsupportedOperation.InstanceStateFreezing"
//  UNSUPPORTEDOPERATION_INSTANCESTATEISOLATING = "UnsupportedOperation.InstanceStateIsolating"
//  UNSUPPORTEDOPERATION_INSTANCESTATEPENDING = "UnsupportedOperation.InstanceStatePending"
//  UNSUPPORTEDOPERATION_INSTANCESTATEREBOOTING = "UnsupportedOperation.InstanceStateRebooting"
//  UNSUPPORTEDOPERATION_INSTANCESTATERESCUEMODE = "UnsupportedOperation.InstanceStateRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATESERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATESHUTDOWN = "UnsupportedOperation.InstanceStateShutdown"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTARTING = "UnsupportedOperation.InstanceStateStarting"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPING = "UnsupportedOperation.InstanceStateStopping"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATED = "UnsupportedOperation.InstanceStateTerminated"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATING = "UnsupportedOperation.InstanceStateTerminating"
//  UNSUPPORTEDOPERATION_INSTANCETYPENOTSUPPORTHIGHDENSITYMODESETTING = "UnsupportedOperation.InstanceTypeNotSupportHighDensityModeSetting"
//  UNSUPPORTEDOPERATION_INSTANCETYPENOTSUPPORTJUMBOFRAME = "UnsupportedOperation.InstanceTypeNotSupportJumboFrame"
//  UNSUPPORTEDOPERATION_INSTANCESENABLEJUMBOWITHOUTREBOOT = "UnsupportedOperation.InstancesEnableJumboWithoutReboot"
//  UNSUPPORTEDOPERATION_INVALIDINSTANCENOTSUPPORTEDPROTECTEDINSTANCE = "UnsupportedOperation.InvalidInstanceNotSupportedProtectedInstance"
//  UNSUPPORTEDOPERATION_STOPPEDMODESTOPCHARGING = "UnsupportedOperation.StoppedModeStopCharging"
func (c *Client) ModifyInstancesAttributeWithContext(ctx context.Context, request *ModifyInstancesAttributeRequest) (response *ModifyInstancesAttributeResponse, err error) {
    if request == nil {
        request = NewModifyInstancesAttributeRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "ModifyInstancesAttribute")
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyInstancesAttribute require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyInstancesAttributeResponse()
    err = c.Send(request, response)
    return
}

func NewModifyInstancesChargeTypeRequest() (request *ModifyInstancesChargeTypeRequest) {
    request = &ModifyInstancesChargeTypeRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "ModifyInstancesChargeType")
    
    
    return
}

func NewModifyInstancesChargeTypeResponse() (response *ModifyInstancesChargeTypeResponse) {
    response = &ModifyInstancesChargeTypeResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// ModifyInstancesChargeType
// This API is used to switch the billing mode of an instance.
//
// 
//
// This API is used to perform operations that do not support instances with no charge when shut down, instances of the model families Batch Compute BC1 and Batch Compute BS1, or instances of scheduled termination.
//
// * You can query the result of the instance operation by calling the API [DescribeInstances](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1#.E7.A4.BA.E4.BE.8B3-.E6.9F.A5.E8.AF.A2.E5.AE.9E.E4.BE.8B.E7.9A.84.E6.9C.80.E6.96.B0.E6.93.8D.E4.BD.9C.E6.83.85.E5.86.B5). If the latest operation status (LatestOperationState) of the instance is **SUCCESS**, the operation is successful.
//
// error code that may be returned:
//  ACCOUNTQUALIFICATIONRESTRICTIONS = "AccountQualificationRestrictions"
//  FAILEDOPERATION_INVALIDINSTANCEAPPLICATIONROLEEMR = "FailedOperation.InvalidInstanceApplicationRoleEmr"
//  FAILEDOPERATION_PROMOTIONALPERIORESTRICTION = "FailedOperation.PromotionalPerioRestriction"
//  INTERNALERROR_TRADEUNKNOWNERROR = "InternalError.TradeUnknownError"
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDACCOUNT_INSUFFICIENTBALANCE = "InvalidAccount.InsufficientBalance"
//  INVALIDACCOUNT_UNPAIDORDER = "InvalidAccount.UnpaidOrder"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_INVALIDINSTANCETYPEUNDERWRITE = "InvalidParameterValue.InvalidInstanceTypeUnderwrite"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPERIOD = "InvalidPeriod"
//  INVALIDPERMISSION = "InvalidPermission"
//  LIMITEXCEEDED_INSTANCEQUOTA = "LimitExceeded.InstanceQuota"
//  LIMITEXCEEDED_INSTANCETYPEBANDWIDTH = "LimitExceeded.InstanceTypeBandwidth"
//  MISSINGPARAMETER = "MissingParameter"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  RESOURCEINSUFFICIENT_CLOUDDISKUNAVAILABLE = "ResourceInsufficient.CloudDiskUnavailable"
//  UNSUPPORTEDOPERATION_COMMERCIALIMAGECHANGECHARGETYPE = "UnsupportedOperation.CommercialImageChangeChargeType"
//  UNSUPPORTEDOPERATION_INSTANCECHARGETYPE = "UnsupportedOperation.InstanceChargeType"
//  UNSUPPORTEDOPERATION_INSTANCEMIXEDZONETYPE = "UnsupportedOperation.InstanceMixedZoneType"
//  UNSUPPORTEDOPERATION_INSTANCESTATEBANNING = "UnsupportedOperation.InstanceStateBanning"
//  UNSUPPORTEDOPERATION_INSTANCESTATEFREEZING = "UnsupportedOperation.InstanceStateFreezing"
//  UNSUPPORTEDOPERATION_INSTANCESTATEREBOOTING = "UnsupportedOperation.InstanceStateRebooting"
//  UNSUPPORTEDOPERATION_INSTANCESTATERESCUEMODE = "UnsupportedOperation.InstanceStateRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATESHUTDOWN = "UnsupportedOperation.InstanceStateShutdown"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPING = "UnsupportedOperation.InstanceStateStopping"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATING = "UnsupportedOperation.InstanceStateTerminating"
//  UNSUPPORTEDOPERATION_REDHATINSTANCEUNSUPPORTED = "UnsupportedOperation.RedHatInstanceUnsupported"
//  UNSUPPORTEDOPERATION_UNDERWRITEDISCOUNTGREATERTHANPREPAIDDISCOUNT = "UnsupportedOperation.UnderwriteDiscountGreaterThanPrepaidDiscount"
//  UNSUPPORTEDOPERATION_UNDERWRITINGINSTANCETYPEONLYSUPPORTAUTORENEW = "UnsupportedOperation.UnderwritingInstanceTypeOnlySupportAutoRenew"
func (c *Client) ModifyInstancesChargeType(request *ModifyInstancesChargeTypeRequest) (response *ModifyInstancesChargeTypeResponse, err error) {
    return c.ModifyInstancesChargeTypeWithContext(context.Background(), request)
}

// ModifyInstancesChargeType
// This API is used to switch the billing mode of an instance.
//
// 
//
// This API is used to perform operations that do not support instances with no charge when shut down, instances of the model families Batch Compute BC1 and Batch Compute BS1, or instances of scheduled termination.
//
// * You can query the result of the instance operation by calling the API [DescribeInstances](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1#.E7.A4.BA.E4.BE.8B3-.E6.9F.A5.E8.AF.A2.E5.AE.9E.E4.BE.8B.E7.9A.84.E6.9C.80.E6.96.B0.E6.93.8D.E4.BD.9C.E6.83.85.E5.86.B5). If the latest operation status (LatestOperationState) of the instance is **SUCCESS**, the operation is successful.
//
// error code that may be returned:
//  ACCOUNTQUALIFICATIONRESTRICTIONS = "AccountQualificationRestrictions"
//  FAILEDOPERATION_INVALIDINSTANCEAPPLICATIONROLEEMR = "FailedOperation.InvalidInstanceApplicationRoleEmr"
//  FAILEDOPERATION_PROMOTIONALPERIORESTRICTION = "FailedOperation.PromotionalPerioRestriction"
//  INTERNALERROR_TRADEUNKNOWNERROR = "InternalError.TradeUnknownError"
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDACCOUNT_INSUFFICIENTBALANCE = "InvalidAccount.InsufficientBalance"
//  INVALIDACCOUNT_UNPAIDORDER = "InvalidAccount.UnpaidOrder"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_INVALIDINSTANCETYPEUNDERWRITE = "InvalidParameterValue.InvalidInstanceTypeUnderwrite"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPERIOD = "InvalidPeriod"
//  INVALIDPERMISSION = "InvalidPermission"
//  LIMITEXCEEDED_INSTANCEQUOTA = "LimitExceeded.InstanceQuota"
//  LIMITEXCEEDED_INSTANCETYPEBANDWIDTH = "LimitExceeded.InstanceTypeBandwidth"
//  MISSINGPARAMETER = "MissingParameter"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  RESOURCEINSUFFICIENT_CLOUDDISKUNAVAILABLE = "ResourceInsufficient.CloudDiskUnavailable"
//  UNSUPPORTEDOPERATION_COMMERCIALIMAGECHANGECHARGETYPE = "UnsupportedOperation.CommercialImageChangeChargeType"
//  UNSUPPORTEDOPERATION_INSTANCECHARGETYPE = "UnsupportedOperation.InstanceChargeType"
//  UNSUPPORTEDOPERATION_INSTANCEMIXEDZONETYPE = "UnsupportedOperation.InstanceMixedZoneType"
//  UNSUPPORTEDOPERATION_INSTANCESTATEBANNING = "UnsupportedOperation.InstanceStateBanning"
//  UNSUPPORTEDOPERATION_INSTANCESTATEFREEZING = "UnsupportedOperation.InstanceStateFreezing"
//  UNSUPPORTEDOPERATION_INSTANCESTATEREBOOTING = "UnsupportedOperation.InstanceStateRebooting"
//  UNSUPPORTEDOPERATION_INSTANCESTATERESCUEMODE = "UnsupportedOperation.InstanceStateRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATESHUTDOWN = "UnsupportedOperation.InstanceStateShutdown"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPING = "UnsupportedOperation.InstanceStateStopping"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATING = "UnsupportedOperation.InstanceStateTerminating"
//  UNSUPPORTEDOPERATION_REDHATINSTANCEUNSUPPORTED = "UnsupportedOperation.RedHatInstanceUnsupported"
//  UNSUPPORTEDOPERATION_UNDERWRITEDISCOUNTGREATERTHANPREPAIDDISCOUNT = "UnsupportedOperation.UnderwriteDiscountGreaterThanPrepaidDiscount"
//  UNSUPPORTEDOPERATION_UNDERWRITINGINSTANCETYPEONLYSUPPORTAUTORENEW = "UnsupportedOperation.UnderwritingInstanceTypeOnlySupportAutoRenew"
func (c *Client) ModifyInstancesChargeTypeWithContext(ctx context.Context, request *ModifyInstancesChargeTypeRequest) (response *ModifyInstancesChargeTypeResponse, err error) {
    if request == nil {
        request = NewModifyInstancesChargeTypeRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "ModifyInstancesChargeType")
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyInstancesChargeType require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyInstancesChargeTypeResponse()
    err = c.Send(request, response)
    return
}

func NewModifyInstancesDisasterRecoverGroupRequest() (request *ModifyInstancesDisasterRecoverGroupRequest) {
    request = &ModifyInstancesDisasterRecoverGroupRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "ModifyInstancesDisasterRecoverGroup")
    
    
    return
}

func NewModifyInstancesDisasterRecoverGroupResponse() (response *ModifyInstancesDisasterRecoverGroupResponse) {
    response = &ModifyInstancesDisasterRecoverGroupResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// ModifyInstancesDisasterRecoverGroup
// This API is used to adjust the placement group of an instance.
//
// * Currently only basic networks or Virtual Private Cloud (VPC) instances are supported.
//
// error code that may be returned:
//  FAILEDOPERATION_ALREADYINDISASTERRECOVERGROUP = "FailedOperation.AlreadyInDisasterRecoverGroup"
//  FAILEDOPERATION_DISASTERRECOVERGROUPNOTFOUND = "FailedOperation.DisasterRecoverGroupNotFound"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETERVALUE_DISASTERRECOVERGROUPIDMALFORMED = "InvalidParameterValue.DisasterRecoverGroupIdMalformed"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_NOTSUPPORTED = "InvalidParameterValue.NotSupported"
//  LIMITEXCEEDED_DISASTERRECOVERGROUP = "LimitExceeded.DisasterRecoverGroup"
//  MUTEXOPERATION_TASKRUNNING = "MutexOperation.TaskRunning"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  RESOURCEINSUFFICIENT_DISASTERRECOVERGROUPCVMQUOTA = "ResourceInsufficient.DisasterRecoverGroupCvmQuota"
//  UNSUPPORTEDOPERATION_CBSREMOTESSDNOTSUPPORT = "UnsupportedOperation.CbsRemoteSsdNotSupport"
//  UNSUPPORTEDOPERATION_INSTANCECHARGETYPE = "UnsupportedOperation.InstanceChargeType"
//  UNSUPPORTEDOPERATION_INSTANCESTATEBANNING = "UnsupportedOperation.InstanceStateBanning"
//  UNSUPPORTEDOPERATION_INSTANCESTATEFREEZING = "UnsupportedOperation.InstanceStateFreezing"
//  UNSUPPORTEDOPERATION_INSTANCESTATEREBOOTING = "UnsupportedOperation.InstanceStateRebooting"
//  UNSUPPORTEDOPERATION_INSTANCESTATERESCUEMODE = "UnsupportedOperation.InstanceStateRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATESHUTDOWN = "UnsupportedOperation.InstanceStateShutdown"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPING = "UnsupportedOperation.InstanceStateStopping"
//  UNSUPPORTEDOPERATION_NOVPCNETWORK = "UnsupportedOperation.NoVpcNetwork"
//  UNSUPPORTEDOPERATION_SPECIALINSTANCETYPE = "UnsupportedOperation.SpecialInstanceType"
//  UNSUPPORTEDOPERATION_STOPPEDMODESTOPCHARGING = "UnsupportedOperation.StoppedModeStopCharging"
func (c *Client) ModifyInstancesDisasterRecoverGroup(request *ModifyInstancesDisasterRecoverGroupRequest) (response *ModifyInstancesDisasterRecoverGroupResponse, err error) {
    return c.ModifyInstancesDisasterRecoverGroupWithContext(context.Background(), request)
}

// ModifyInstancesDisasterRecoverGroup
// This API is used to adjust the placement group of an instance.
//
// * Currently only basic networks or Virtual Private Cloud (VPC) instances are supported.
//
// error code that may be returned:
//  FAILEDOPERATION_ALREADYINDISASTERRECOVERGROUP = "FailedOperation.AlreadyInDisasterRecoverGroup"
//  FAILEDOPERATION_DISASTERRECOVERGROUPNOTFOUND = "FailedOperation.DisasterRecoverGroupNotFound"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETERVALUE_DISASTERRECOVERGROUPIDMALFORMED = "InvalidParameterValue.DisasterRecoverGroupIdMalformed"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_NOTSUPPORTED = "InvalidParameterValue.NotSupported"
//  LIMITEXCEEDED_DISASTERRECOVERGROUP = "LimitExceeded.DisasterRecoverGroup"
//  MUTEXOPERATION_TASKRUNNING = "MutexOperation.TaskRunning"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  RESOURCEINSUFFICIENT_DISASTERRECOVERGROUPCVMQUOTA = "ResourceInsufficient.DisasterRecoverGroupCvmQuota"
//  UNSUPPORTEDOPERATION_CBSREMOTESSDNOTSUPPORT = "UnsupportedOperation.CbsRemoteSsdNotSupport"
//  UNSUPPORTEDOPERATION_INSTANCECHARGETYPE = "UnsupportedOperation.InstanceChargeType"
//  UNSUPPORTEDOPERATION_INSTANCESTATEBANNING = "UnsupportedOperation.InstanceStateBanning"
//  UNSUPPORTEDOPERATION_INSTANCESTATEFREEZING = "UnsupportedOperation.InstanceStateFreezing"
//  UNSUPPORTEDOPERATION_INSTANCESTATEREBOOTING = "UnsupportedOperation.InstanceStateRebooting"
//  UNSUPPORTEDOPERATION_INSTANCESTATERESCUEMODE = "UnsupportedOperation.InstanceStateRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATESHUTDOWN = "UnsupportedOperation.InstanceStateShutdown"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPING = "UnsupportedOperation.InstanceStateStopping"
//  UNSUPPORTEDOPERATION_NOVPCNETWORK = "UnsupportedOperation.NoVpcNetwork"
//  UNSUPPORTEDOPERATION_SPECIALINSTANCETYPE = "UnsupportedOperation.SpecialInstanceType"
//  UNSUPPORTEDOPERATION_STOPPEDMODESTOPCHARGING = "UnsupportedOperation.StoppedModeStopCharging"
func (c *Client) ModifyInstancesDisasterRecoverGroupWithContext(ctx context.Context, request *ModifyInstancesDisasterRecoverGroupRequest) (response *ModifyInstancesDisasterRecoverGroupResponse, err error) {
    if request == nil {
        request = NewModifyInstancesDisasterRecoverGroupRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "ModifyInstancesDisasterRecoverGroup")
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyInstancesDisasterRecoverGroup require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyInstancesDisasterRecoverGroupResponse()
    err = c.Send(request, response)
    return
}

func NewModifyInstancesProjectRequest() (request *ModifyInstancesProjectRequest) {
    request = &ModifyInstancesProjectRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "ModifyInstancesProject")
    
    
    return
}

func NewModifyInstancesProjectResponse() (response *ModifyInstancesProjectResponse) {
    response = &ModifyInstancesProjectResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// ModifyInstancesProject
// This API is used to change the project to which an instance is assigned.
//
// 
//
// * Project is a virtual concept. You can create multiple projects under one account, manage different resources in each project, and assign different instances to different projects. You may use the [DescribeInstances](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1) API to query instances and use the project ID to filter the results.
//
// * You cannot modify the project of an instance that is bound to a load balancer. You need to unbind the load balancer from the instance by using the [DeregisterInstancesFromLoadBalancer](https://intl.cloud.tencent.com/document/api/214/1258?from_cn_redirect=1) API before using this API.
//
// * Batch operations are supported. Up to 100 instances per request is allowed.
//
// * You can use the API [DescribeInstances](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1#.E7.A4.BA.E4.BE.8B3-.E6.9F.A5.E8.AF.A2.E5.AE.9E.E4.BE.8B.E7.9A.84.E6.9C.80.E6.96.B0.E6.93.8D.E4.BD.9C.E6.83.85.E5.86.B5) to query the operation result. If the `LatestOperationState` in the response is `SUCCESS`, the operation is successful.
//
// error code that may be returned:
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPROJECTID_NOTFOUND = "InvalidProjectId.NotFound"
//  MISSINGPARAMETER = "MissingParameter"
//  MUTEXOPERATION_TASKRUNNING = "MutexOperation.TaskRunning"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  UNSUPPORTEDOPERATION_INSTANCESTATEISOLATING = "UnsupportedOperation.InstanceStateIsolating"
func (c *Client) ModifyInstancesProject(request *ModifyInstancesProjectRequest) (response *ModifyInstancesProjectResponse, err error) {
    return c.ModifyInstancesProjectWithContext(context.Background(), request)
}

// ModifyInstancesProject
// This API is used to change the project to which an instance is assigned.
//
// 
//
// * Project is a virtual concept. You can create multiple projects under one account, manage different resources in each project, and assign different instances to different projects. You may use the [DescribeInstances](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1) API to query instances and use the project ID to filter the results.
//
// * You cannot modify the project of an instance that is bound to a load balancer. You need to unbind the load balancer from the instance by using the [DeregisterInstancesFromLoadBalancer](https://intl.cloud.tencent.com/document/api/214/1258?from_cn_redirect=1) API before using this API.
//
// * Batch operations are supported. Up to 100 instances per request is allowed.
//
// * You can use the API [DescribeInstances](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1#.E7.A4.BA.E4.BE.8B3-.E6.9F.A5.E8.AF.A2.E5.AE.9E.E4.BE.8B.E7.9A.84.E6.9C.80.E6.96.B0.E6.93.8D.E4.BD.9C.E6.83.85.E5.86.B5) to query the operation result. If the `LatestOperationState` in the response is `SUCCESS`, the operation is successful.
//
// error code that may be returned:
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPROJECTID_NOTFOUND = "InvalidProjectId.NotFound"
//  MISSINGPARAMETER = "MissingParameter"
//  MUTEXOPERATION_TASKRUNNING = "MutexOperation.TaskRunning"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  UNSUPPORTEDOPERATION_INSTANCESTATEISOLATING = "UnsupportedOperation.InstanceStateIsolating"
func (c *Client) ModifyInstancesProjectWithContext(ctx context.Context, request *ModifyInstancesProjectRequest) (response *ModifyInstancesProjectResponse, err error) {
    if request == nil {
        request = NewModifyInstancesProjectRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "ModifyInstancesProject")
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyInstancesProject require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyInstancesProjectResponse()
    err = c.Send(request, response)
    return
}

func NewModifyInstancesRenewFlagRequest() (request *ModifyInstancesRenewFlagRequest) {
    request = &ModifyInstancesRenewFlagRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "ModifyInstancesRenewFlag")
    
    
    return
}

func NewModifyInstancesRenewFlagResponse() (response *ModifyInstancesRenewFlagResponse) {
    response = &ModifyInstancesRenewFlagResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// ModifyInstancesRenewFlag
// This API is used to modify the renewal flag of monthly subscription instances.
//
// 
//
// * After an instance is marked as auto-renewal, it will be automatically renewed for one month each time it expires.
//
// * Batch operations are supported. The maximum number of instances for each request is 100.* You can query the result of the instance operation by calling the API [DescribeInstances](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1#.E7.A4.BA.E4.BE.8B3-.E6.9F.A5.E8.AF.A2.E5.AE.9E.E4.BE.8B.E7.9A.84.E6.9C.80.E6.96.B0.E6.93.8D.E4.BD.9C.E6.83.85.E5.86.B5). If the latest operation status (LatestOperationState) of the instance is **SUCCESS**, the operation is successful.
//
// error code that may be returned:
//  FAILEDOPERATION_INVALIDINSTANCEAPPLICATIONROLEEMR = "FailedOperation.InvalidInstanceApplicationRoleEmr"
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  MUTEXOPERATION_TASKRUNNING = "MutexOperation.TaskRunning"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  UNAUTHORIZEDOPERATION_MFAEXPIRED = "UnauthorizedOperation.MFAExpired"
//  UNAUTHORIZEDOPERATION_MFANOTFOUND = "UnauthorizedOperation.MFANotFound"
//  UNSUPPORTEDOPERATION_INSTANCECHARGETYPE = "UnsupportedOperation.InstanceChargeType"
//  UNSUPPORTEDOPERATION_INSTANCEREINSTALLFAILED = "UnsupportedOperation.InstanceReinstallFailed"
//  UNSUPPORTEDOPERATION_INSTANCESTATEBANNING = "UnsupportedOperation.InstanceStateBanning"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERRESCUEMODE = "UnsupportedOperation.InstanceStateEnterRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITRESCUEMODE = "UnsupportedOperation.InstanceStateExitRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEFREEZING = "UnsupportedOperation.InstanceStateFreezing"
//  UNSUPPORTEDOPERATION_INSTANCESTATEISOLATING = "UnsupportedOperation.InstanceStateIsolating"
//  UNSUPPORTEDOPERATION_INSTANCESTATEPENDING = "UnsupportedOperation.InstanceStatePending"
//  UNSUPPORTEDOPERATION_INSTANCESTATEREBOOTING = "UnsupportedOperation.InstanceStateRebooting"
//  UNSUPPORTEDOPERATION_INSTANCESTATERESCUEMODE = "UnsupportedOperation.InstanceStateRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATESERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATESHUTDOWN = "UnsupportedOperation.InstanceStateShutdown"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTARTING = "UnsupportedOperation.InstanceStateStarting"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPING = "UnsupportedOperation.InstanceStateStopping"
//  UNSUPPORTEDOPERATION_SPECIALINSTANCETYPE = "UnsupportedOperation.SpecialInstanceType"
//  UNSUPPORTEDOPERATION_UNDERWRITINGINSTANCETYPEONLYSUPPORTAUTORENEW = "UnsupportedOperation.UnderwritingInstanceTypeOnlySupportAutoRenew"
func (c *Client) ModifyInstancesRenewFlag(request *ModifyInstancesRenewFlagRequest) (response *ModifyInstancesRenewFlagResponse, err error) {
    return c.ModifyInstancesRenewFlagWithContext(context.Background(), request)
}

// ModifyInstancesRenewFlag
// This API is used to modify the renewal flag of monthly subscription instances.
//
// 
//
// * After an instance is marked as auto-renewal, it will be automatically renewed for one month each time it expires.
//
// * Batch operations are supported. The maximum number of instances for each request is 100.* You can query the result of the instance operation by calling the API [DescribeInstances](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1#.E7.A4.BA.E4.BE.8B3-.E6.9F.A5.E8.AF.A2.E5.AE.9E.E4.BE.8B.E7.9A.84.E6.9C.80.E6.96.B0.E6.93.8D.E4.BD.9C.E6.83.85.E5.86.B5). If the latest operation status (LatestOperationState) of the instance is **SUCCESS**, the operation is successful.
//
// error code that may be returned:
//  FAILEDOPERATION_INVALIDINSTANCEAPPLICATIONROLEEMR = "FailedOperation.InvalidInstanceApplicationRoleEmr"
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  MUTEXOPERATION_TASKRUNNING = "MutexOperation.TaskRunning"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  UNAUTHORIZEDOPERATION_MFAEXPIRED = "UnauthorizedOperation.MFAExpired"
//  UNAUTHORIZEDOPERATION_MFANOTFOUND = "UnauthorizedOperation.MFANotFound"
//  UNSUPPORTEDOPERATION_INSTANCECHARGETYPE = "UnsupportedOperation.InstanceChargeType"
//  UNSUPPORTEDOPERATION_INSTANCEREINSTALLFAILED = "UnsupportedOperation.InstanceReinstallFailed"
//  UNSUPPORTEDOPERATION_INSTANCESTATEBANNING = "UnsupportedOperation.InstanceStateBanning"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERRESCUEMODE = "UnsupportedOperation.InstanceStateEnterRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITRESCUEMODE = "UnsupportedOperation.InstanceStateExitRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEFREEZING = "UnsupportedOperation.InstanceStateFreezing"
//  UNSUPPORTEDOPERATION_INSTANCESTATEISOLATING = "UnsupportedOperation.InstanceStateIsolating"
//  UNSUPPORTEDOPERATION_INSTANCESTATEPENDING = "UnsupportedOperation.InstanceStatePending"
//  UNSUPPORTEDOPERATION_INSTANCESTATEREBOOTING = "UnsupportedOperation.InstanceStateRebooting"
//  UNSUPPORTEDOPERATION_INSTANCESTATERESCUEMODE = "UnsupportedOperation.InstanceStateRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATESERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATESHUTDOWN = "UnsupportedOperation.InstanceStateShutdown"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTARTING = "UnsupportedOperation.InstanceStateStarting"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPING = "UnsupportedOperation.InstanceStateStopping"
//  UNSUPPORTEDOPERATION_SPECIALINSTANCETYPE = "UnsupportedOperation.SpecialInstanceType"
//  UNSUPPORTEDOPERATION_UNDERWRITINGINSTANCETYPEONLYSUPPORTAUTORENEW = "UnsupportedOperation.UnderwritingInstanceTypeOnlySupportAutoRenew"
func (c *Client) ModifyInstancesRenewFlagWithContext(ctx context.Context, request *ModifyInstancesRenewFlagRequest) (response *ModifyInstancesRenewFlagResponse, err error) {
    if request == nil {
        request = NewModifyInstancesRenewFlagRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "ModifyInstancesRenewFlag")
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyInstancesRenewFlag require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyInstancesRenewFlagResponse()
    err = c.Send(request, response)
    return
}

func NewModifyInstancesVpcAttributeRequest() (request *ModifyInstancesVpcAttributeRequest) {
    request = &ModifyInstancesVpcAttributeRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "ModifyInstancesVpcAttribute")
    
    
    return
}

func NewModifyInstancesVpcAttributeResponse() (response *ModifyInstancesVpcAttributeResponse) {
    response = &ModifyInstancesVpcAttributeResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// ModifyInstancesVpcAttribute
// This API is used to modify the VPC attributes of an instance, such as the VPC IP address.
//
// * This action will shut down the instance, and restart it after the modification is completed.
//
// * To migrate an instance to another VPC/subnet, specify the new VPC and subnet directly. Make sure that the instance to migrate is not bound to an [ENI](https://intl.cloud.tencent.com/document/product/576?from_cn_redirect=1) or [CLB](https://intl.cloud.tencent.com/document/product/214?from_cn_redirect=1) instances.
//
// * You can use the API [DescribeInstances](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1#.E7.A4.BA.E4.BE.8B3-.E6.9F.A5.E8.AF.A2.E5.AE.9E.E4.BE.8B.E7.9A.84.E6.9C.80.E6.96.B0.E6.93.8D.E4.BD.9C.E6.83.85.E5.86.B5) to query the operation result. If the `LatestOperationState` in the response is `SUCCESS`, the operation is successful.
//
// error code that may be returned:
//  ENINOTALLOWEDCHANGESUBNET = "EniNotAllowedChangeSubnet"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDINSTANCESTATE = "InvalidInstanceState"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_IPADDRESSMALFORMED = "InvalidParameterValue.IPAddressMalformed"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_INVALIDIPFORMAT = "InvalidParameterValue.InvalidIpFormat"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_MUSTDHCPENABLEDVPC = "InvalidParameterValue.MustDhcpEnabledVpc"
//  INVALIDPARAMETERVALUE_SUBNETIDMALFORMED = "InvalidParameterValue.SubnetIdMalformed"
//  INVALIDPARAMETERVALUE_VPCIDMALFORMED = "InvalidParameterValue.VpcIdMalformed"
//  MUTEXOPERATION_TASKRUNNING = "MutexOperation.TaskRunning"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  UNAUTHORIZEDOPERATION_MFAEXPIRED = "UnauthorizedOperation.MFAExpired"
//  UNAUTHORIZEDOPERATION_MFANOTFOUND = "UnauthorizedOperation.MFANotFound"
//  UNSUPPORTEDOPERATION_EDGEZONEINSTANCE = "UnsupportedOperation.EdgeZoneInstance"
//  UNSUPPORTEDOPERATION_ELASTICNETWORKINTERFACE = "UnsupportedOperation.ElasticNetworkInterface"
//  UNSUPPORTEDOPERATION_IPV6NOTSUPPORTVPCMIGRATE = "UnsupportedOperation.IPv6NotSupportVpcMigrate"
//  UNSUPPORTEDOPERATION_INSTANCECHARGETYPE = "UnsupportedOperation.InstanceChargeType"
//  UNSUPPORTEDOPERATION_INSTANCEREINSTALLFAILED = "UnsupportedOperation.InstanceReinstallFailed"
//  UNSUPPORTEDOPERATION_INSTANCESTATEBANNING = "UnsupportedOperation.InstanceStateBanning"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERRESCUEMODE = "UnsupportedOperation.InstanceStateEnterRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITRESCUEMODE = "UnsupportedOperation.InstanceStateExitRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITSERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateExitServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATEFREEZING = "UnsupportedOperation.InstanceStateFreezing"
//  UNSUPPORTEDOPERATION_INSTANCESTATEISOLATING = "UnsupportedOperation.InstanceStateIsolating"
//  UNSUPPORTEDOPERATION_INSTANCESTATEPENDING = "UnsupportedOperation.InstanceStatePending"
//  UNSUPPORTEDOPERATION_INSTANCESTATEREBOOTING = "UnsupportedOperation.InstanceStateRebooting"
//  UNSUPPORTEDOPERATION_INSTANCESTATERESCUEMODE = "UnsupportedOperation.InstanceStateRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATERUNNING = "UnsupportedOperation.InstanceStateRunning"
//  UNSUPPORTEDOPERATION_INSTANCESTATESERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATESHUTDOWN = "UnsupportedOperation.InstanceStateShutdown"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTARTING = "UnsupportedOperation.InstanceStateStarting"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPING = "UnsupportedOperation.InstanceStateStopping"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATING = "UnsupportedOperation.InstanceStateTerminating"
//  UNSUPPORTEDOPERATION_MODIFYVPCWITHCLB = "UnsupportedOperation.ModifyVPCWithCLB"
//  UNSUPPORTEDOPERATION_MODIFYVPCWITHCLASSLINK = "UnsupportedOperation.ModifyVPCWithClassLink"
//  UNSUPPORTEDOPERATION_NOVPCNETWORK = "UnsupportedOperation.NoVpcNetwork"
//  UNSUPPORTEDOPERATION_STOPPEDMODESTOPCHARGING = "UnsupportedOperation.StoppedModeStopCharging"
//  VPCADDRNOTINSUBNET = "VpcAddrNotInSubNet"
//  VPCIPISUSED = "VpcIpIsUsed"
func (c *Client) ModifyInstancesVpcAttribute(request *ModifyInstancesVpcAttributeRequest) (response *ModifyInstancesVpcAttributeResponse, err error) {
    return c.ModifyInstancesVpcAttributeWithContext(context.Background(), request)
}

// ModifyInstancesVpcAttribute
// This API is used to modify the VPC attributes of an instance, such as the VPC IP address.
//
// * This action will shut down the instance, and restart it after the modification is completed.
//
// * To migrate an instance to another VPC/subnet, specify the new VPC and subnet directly. Make sure that the instance to migrate is not bound to an [ENI](https://intl.cloud.tencent.com/document/product/576?from_cn_redirect=1) or [CLB](https://intl.cloud.tencent.com/document/product/214?from_cn_redirect=1) instances.
//
// * You can use the API [DescribeInstances](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1#.E7.A4.BA.E4.BE.8B3-.E6.9F.A5.E8.AF.A2.E5.AE.9E.E4.BE.8B.E7.9A.84.E6.9C.80.E6.96.B0.E6.93.8D.E4.BD.9C.E6.83.85.E5.86.B5) to query the operation result. If the `LatestOperationState` in the response is `SUCCESS`, the operation is successful.
//
// error code that may be returned:
//  ENINOTALLOWEDCHANGESUBNET = "EniNotAllowedChangeSubnet"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDINSTANCESTATE = "InvalidInstanceState"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_IPADDRESSMALFORMED = "InvalidParameterValue.IPAddressMalformed"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_INVALIDIPFORMAT = "InvalidParameterValue.InvalidIpFormat"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_MUSTDHCPENABLEDVPC = "InvalidParameterValue.MustDhcpEnabledVpc"
//  INVALIDPARAMETERVALUE_SUBNETIDMALFORMED = "InvalidParameterValue.SubnetIdMalformed"
//  INVALIDPARAMETERVALUE_VPCIDMALFORMED = "InvalidParameterValue.VpcIdMalformed"
//  MUTEXOPERATION_TASKRUNNING = "MutexOperation.TaskRunning"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  UNAUTHORIZEDOPERATION_MFAEXPIRED = "UnauthorizedOperation.MFAExpired"
//  UNAUTHORIZEDOPERATION_MFANOTFOUND = "UnauthorizedOperation.MFANotFound"
//  UNSUPPORTEDOPERATION_EDGEZONEINSTANCE = "UnsupportedOperation.EdgeZoneInstance"
//  UNSUPPORTEDOPERATION_ELASTICNETWORKINTERFACE = "UnsupportedOperation.ElasticNetworkInterface"
//  UNSUPPORTEDOPERATION_IPV6NOTSUPPORTVPCMIGRATE = "UnsupportedOperation.IPv6NotSupportVpcMigrate"
//  UNSUPPORTEDOPERATION_INSTANCECHARGETYPE = "UnsupportedOperation.InstanceChargeType"
//  UNSUPPORTEDOPERATION_INSTANCEREINSTALLFAILED = "UnsupportedOperation.InstanceReinstallFailed"
//  UNSUPPORTEDOPERATION_INSTANCESTATEBANNING = "UnsupportedOperation.InstanceStateBanning"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERRESCUEMODE = "UnsupportedOperation.InstanceStateEnterRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITRESCUEMODE = "UnsupportedOperation.InstanceStateExitRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITSERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateExitServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATEFREEZING = "UnsupportedOperation.InstanceStateFreezing"
//  UNSUPPORTEDOPERATION_INSTANCESTATEISOLATING = "UnsupportedOperation.InstanceStateIsolating"
//  UNSUPPORTEDOPERATION_INSTANCESTATEPENDING = "UnsupportedOperation.InstanceStatePending"
//  UNSUPPORTEDOPERATION_INSTANCESTATEREBOOTING = "UnsupportedOperation.InstanceStateRebooting"
//  UNSUPPORTEDOPERATION_INSTANCESTATERESCUEMODE = "UnsupportedOperation.InstanceStateRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATERUNNING = "UnsupportedOperation.InstanceStateRunning"
//  UNSUPPORTEDOPERATION_INSTANCESTATESERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATESHUTDOWN = "UnsupportedOperation.InstanceStateShutdown"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTARTING = "UnsupportedOperation.InstanceStateStarting"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPING = "UnsupportedOperation.InstanceStateStopping"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATING = "UnsupportedOperation.InstanceStateTerminating"
//  UNSUPPORTEDOPERATION_MODIFYVPCWITHCLB = "UnsupportedOperation.ModifyVPCWithCLB"
//  UNSUPPORTEDOPERATION_MODIFYVPCWITHCLASSLINK = "UnsupportedOperation.ModifyVPCWithClassLink"
//  UNSUPPORTEDOPERATION_NOVPCNETWORK = "UnsupportedOperation.NoVpcNetwork"
//  UNSUPPORTEDOPERATION_STOPPEDMODESTOPCHARGING = "UnsupportedOperation.StoppedModeStopCharging"
//  VPCADDRNOTINSUBNET = "VpcAddrNotInSubNet"
//  VPCIPISUSED = "VpcIpIsUsed"
func (c *Client) ModifyInstancesVpcAttributeWithContext(ctx context.Context, request *ModifyInstancesVpcAttributeRequest) (response *ModifyInstancesVpcAttributeResponse, err error) {
    if request == nil {
        request = NewModifyInstancesVpcAttributeRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "ModifyInstancesVpcAttribute")
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyInstancesVpcAttribute require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyInstancesVpcAttributeResponse()
    err = c.Send(request, response)
    return
}

func NewModifyKeyPairAttributeRequest() (request *ModifyKeyPairAttributeRequest) {
    request = &ModifyKeyPairAttributeRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "ModifyKeyPairAttribute")
    
    
    return
}

func NewModifyKeyPairAttributeResponse() (response *ModifyKeyPairAttributeResponse) {
    response = &ModifyKeyPairAttributeResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// ModifyKeyPairAttribute
// This API is used to modify attributes of a key pair.
//
// 
//
// * Modify the name and description information of the key pair specified by the key pair ID.
//
// * The key pair name should not be the same as the name of an existing key pair.
//
// * The key pair ID is the unique identifier of a key pair and cannot be modified.
//
// 
//
// * Either the key pair name or description information should be specified, and both can also be specified simultaneously.
//
// error code that may be returned:
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDKEYPAIRID_MALFORMED = "InvalidKeyPairId.Malformed"
//  INVALIDKEYPAIRID_NOTFOUND = "InvalidKeyPairId.NotFound"
//  INVALIDKEYPAIRNAME_DUPLICATE = "InvalidKeyPairName.Duplicate"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  MISSINGPARAMETER = "MissingParameter"
func (c *Client) ModifyKeyPairAttribute(request *ModifyKeyPairAttributeRequest) (response *ModifyKeyPairAttributeResponse, err error) {
    return c.ModifyKeyPairAttributeWithContext(context.Background(), request)
}

// ModifyKeyPairAttribute
// This API is used to modify attributes of a key pair.
//
// 
//
// * Modify the name and description information of the key pair specified by the key pair ID.
//
// * The key pair name should not be the same as the name of an existing key pair.
//
// * The key pair ID is the unique identifier of a key pair and cannot be modified.
//
// 
//
// * Either the key pair name or description information should be specified, and both can also be specified simultaneously.
//
// error code that may be returned:
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDKEYPAIRID_MALFORMED = "InvalidKeyPairId.Malformed"
//  INVALIDKEYPAIRID_NOTFOUND = "InvalidKeyPairId.NotFound"
//  INVALIDKEYPAIRNAME_DUPLICATE = "InvalidKeyPairName.Duplicate"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  MISSINGPARAMETER = "MissingParameter"
func (c *Client) ModifyKeyPairAttributeWithContext(ctx context.Context, request *ModifyKeyPairAttributeRequest) (response *ModifyKeyPairAttributeResponse, err error) {
    if request == nil {
        request = NewModifyKeyPairAttributeRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "ModifyKeyPairAttribute")
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyKeyPairAttribute require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyKeyPairAttributeResponse()
    err = c.Send(request, response)
    return
}

func NewModifyLaunchTemplateDefaultVersionRequest() (request *ModifyLaunchTemplateDefaultVersionRequest) {
    request = &ModifyLaunchTemplateDefaultVersionRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "ModifyLaunchTemplateDefaultVersion")
    
    
    return
}

func NewModifyLaunchTemplateDefaultVersionResponse() (response *ModifyLaunchTemplateDefaultVersionResponse) {
    response = &ModifyLaunchTemplateDefaultVersionResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// ModifyLaunchTemplateDefaultVersion
// This API is used to modify the default version of the instance launch template.
//
// error code that may be returned:
//  INVALIDPARAMETERCOMBINATION = "InvalidParameterCombination"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDMALFORMED = "InvalidParameterValue.LaunchTemplateIdMalformed"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDNOTEXISTED = "InvalidParameterValue.LaunchTemplateIdNotExisted"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDVERNOTEXISTED = "InvalidParameterValue.LaunchTemplateIdVerNotExisted"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDVERSETALREADY = "InvalidParameterValue.LaunchTemplateIdVerSetAlready"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATENOTFOUND = "InvalidParameterValue.LaunchTemplateNotFound"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEVERSION = "InvalidParameterValue.LaunchTemplateVersion"
//  MISSINGPARAMETER = "MissingParameter"
//  UNKNOWNPARAMETER = "UnknownParameter"
func (c *Client) ModifyLaunchTemplateDefaultVersion(request *ModifyLaunchTemplateDefaultVersionRequest) (response *ModifyLaunchTemplateDefaultVersionResponse, err error) {
    return c.ModifyLaunchTemplateDefaultVersionWithContext(context.Background(), request)
}

// ModifyLaunchTemplateDefaultVersion
// This API is used to modify the default version of the instance launch template.
//
// error code that may be returned:
//  INVALIDPARAMETERCOMBINATION = "InvalidParameterCombination"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDMALFORMED = "InvalidParameterValue.LaunchTemplateIdMalformed"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDNOTEXISTED = "InvalidParameterValue.LaunchTemplateIdNotExisted"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDVERNOTEXISTED = "InvalidParameterValue.LaunchTemplateIdVerNotExisted"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDVERSETALREADY = "InvalidParameterValue.LaunchTemplateIdVerSetAlready"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATENOTFOUND = "InvalidParameterValue.LaunchTemplateNotFound"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEVERSION = "InvalidParameterValue.LaunchTemplateVersion"
//  MISSINGPARAMETER = "MissingParameter"
//  UNKNOWNPARAMETER = "UnknownParameter"
func (c *Client) ModifyLaunchTemplateDefaultVersionWithContext(ctx context.Context, request *ModifyLaunchTemplateDefaultVersionRequest) (response *ModifyLaunchTemplateDefaultVersionResponse, err error) {
    if request == nil {
        request = NewModifyLaunchTemplateDefaultVersionRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "ModifyLaunchTemplateDefaultVersion")
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyLaunchTemplateDefaultVersion require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyLaunchTemplateDefaultVersionResponse()
    err = c.Send(request, response)
    return
}

func NewPurchaseReservedInstancesOfferingRequest() (request *PurchaseReservedInstancesOfferingRequest) {
    request = &PurchaseReservedInstancesOfferingRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "PurchaseReservedInstancesOffering")
    
    
    return
}

func NewPurchaseReservedInstancesOfferingResponse() (response *PurchaseReservedInstancesOfferingResponse) {
    response = &PurchaseReservedInstancesOfferingResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// PurchaseReservedInstancesOffering
// This API is used to purchase one or more specific Reserved Instances.
//
// error code that may be returned:
//  INTERNALERROR_TRADEUNKNOWNERROR = "InternalError.TradeUnknownError"
//  INVALIDACCOUNT_INSUFFICIENTBALANCE = "InvalidAccount.InsufficientBalance"
//  INVALIDCLIENTTOKEN_TOOLONG = "InvalidClientToken.TooLong"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_RANGE = "InvalidParameterValue.Range"
//  UNSUPPORTEDOPERATION_INVALIDPERMISSIONNONINTERNATIONALACCOUNT = "UnsupportedOperation.InvalidPermissionNonInternationalAccount"
//  UNSUPPORTEDOPERATION_RESERVEDINSTANCEINVISIBLEFORUSER = "UnsupportedOperation.ReservedInstanceInvisibleForUser"
//  UNSUPPORTEDOPERATION_RESERVEDINSTANCEOUTOFQUATA = "UnsupportedOperation.ReservedInstanceOutofQuata"
func (c *Client) PurchaseReservedInstancesOffering(request *PurchaseReservedInstancesOfferingRequest) (response *PurchaseReservedInstancesOfferingResponse, err error) {
    return c.PurchaseReservedInstancesOfferingWithContext(context.Background(), request)
}

// PurchaseReservedInstancesOffering
// This API is used to purchase one or more specific Reserved Instances.
//
// error code that may be returned:
//  INTERNALERROR_TRADEUNKNOWNERROR = "InternalError.TradeUnknownError"
//  INVALIDACCOUNT_INSUFFICIENTBALANCE = "InvalidAccount.InsufficientBalance"
//  INVALIDCLIENTTOKEN_TOOLONG = "InvalidClientToken.TooLong"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_RANGE = "InvalidParameterValue.Range"
//  UNSUPPORTEDOPERATION_INVALIDPERMISSIONNONINTERNATIONALACCOUNT = "UnsupportedOperation.InvalidPermissionNonInternationalAccount"
//  UNSUPPORTEDOPERATION_RESERVEDINSTANCEINVISIBLEFORUSER = "UnsupportedOperation.ReservedInstanceInvisibleForUser"
//  UNSUPPORTEDOPERATION_RESERVEDINSTANCEOUTOFQUATA = "UnsupportedOperation.ReservedInstanceOutofQuata"
func (c *Client) PurchaseReservedInstancesOfferingWithContext(ctx context.Context, request *PurchaseReservedInstancesOfferingRequest) (response *PurchaseReservedInstancesOfferingResponse, err error) {
    if request == nil {
        request = NewPurchaseReservedInstancesOfferingRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "PurchaseReservedInstancesOffering")
    
    if c.GetCredential() == nil {
        return nil, errors.New("PurchaseReservedInstancesOffering require credential")
    }

    request.SetContext(ctx)
    
    response = NewPurchaseReservedInstancesOfferingResponse()
    err = c.Send(request, response)
    return
}

func NewRebootInstancesRequest() (request *RebootInstancesRequest) {
    request = &RebootInstancesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "RebootInstances")
    
    
    return
}

func NewRebootInstancesResponse() (response *RebootInstancesResponse) {
    response = &RebootInstancesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// RebootInstances
// This API is used to restart instances.
//
// 
//
// * You can only perform this operation on instances whose status is `RUNNING`.
//
// * If the API is called successfully, the instance status will become `REBOOTING`. After the instance is restarted, its status will become `RUNNING` again.
//
// * Forced restart is supported. A forced restart is similar to switching off the power of a physical computer and starting it again. It may cause data loss or file system corruption. Be sure to only force start a CVM when it cannot be restarted normally.
//
// * Batch operations are supported. The maximum number of instances in each request is 100.
//
// error code that may be returned:
//  FAILEDOPERATION_INVALIDINSTANCEAPPLICATIONROLE = "FailedOperation.InvalidInstanceApplicationRole"
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETERCOMBINATION = "InvalidParameterCombination"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  MISSINGPARAMETER = "MissingParameter"
//  MUTEXOPERATION_TASKRUNNING = "MutexOperation.TaskRunning"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  UNAUTHORIZEDOPERATION_MFAEXPIRED = "UnauthorizedOperation.MFAExpired"
//  UNAUTHORIZEDOPERATION_MFANOTFOUND = "UnauthorizedOperation.MFANotFound"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
//  UNSUPPORTEDOPERATION_INSTANCEREINSTALLFAILED = "UnsupportedOperation.InstanceReinstallFailed"
//  UNSUPPORTEDOPERATION_INSTANCESTATECORRUPTED = "UnsupportedOperation.InstanceStateCorrupted"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERRESCUEMODE = "UnsupportedOperation.InstanceStateEnterRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERSERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateEnterServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITRESCUEMODE = "UnsupportedOperation.InstanceStateExitRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITSERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateExitServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATEFREEZING = "UnsupportedOperation.InstanceStateFreezing"
//  UNSUPPORTEDOPERATION_INSTANCESTATEISOLATING = "UnsupportedOperation.InstanceStateIsolating"
//  UNSUPPORTEDOPERATION_INSTANCESTATEPENDING = "UnsupportedOperation.InstanceStatePending"
//  UNSUPPORTEDOPERATION_INSTANCESTATEREBOOTING = "UnsupportedOperation.InstanceStateRebooting"
//  UNSUPPORTEDOPERATION_INSTANCESTATESERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATESHUTDOWN = "UnsupportedOperation.InstanceStateShutdown"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTARTING = "UnsupportedOperation.InstanceStateStarting"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPED = "UnsupportedOperation.InstanceStateStopped"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPING = "UnsupportedOperation.InstanceStateStopping"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATING = "UnsupportedOperation.InstanceStateTerminating"
//  UNSUPPORTEDOPERATION_STOPPEDMODESTOPCHARGING = "UnsupportedOperation.StoppedModeStopCharging"
func (c *Client) RebootInstances(request *RebootInstancesRequest) (response *RebootInstancesResponse, err error) {
    return c.RebootInstancesWithContext(context.Background(), request)
}

// RebootInstances
// This API is used to restart instances.
//
// 
//
// * You can only perform this operation on instances whose status is `RUNNING`.
//
// * If the API is called successfully, the instance status will become `REBOOTING`. After the instance is restarted, its status will become `RUNNING` again.
//
// * Forced restart is supported. A forced restart is similar to switching off the power of a physical computer and starting it again. It may cause data loss or file system corruption. Be sure to only force start a CVM when it cannot be restarted normally.
//
// * Batch operations are supported. The maximum number of instances in each request is 100.
//
// error code that may be returned:
//  FAILEDOPERATION_INVALIDINSTANCEAPPLICATIONROLE = "FailedOperation.InvalidInstanceApplicationRole"
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETERCOMBINATION = "InvalidParameterCombination"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  MISSINGPARAMETER = "MissingParameter"
//  MUTEXOPERATION_TASKRUNNING = "MutexOperation.TaskRunning"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  UNAUTHORIZEDOPERATION_MFAEXPIRED = "UnauthorizedOperation.MFAExpired"
//  UNAUTHORIZEDOPERATION_MFANOTFOUND = "UnauthorizedOperation.MFANotFound"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
//  UNSUPPORTEDOPERATION_INSTANCEREINSTALLFAILED = "UnsupportedOperation.InstanceReinstallFailed"
//  UNSUPPORTEDOPERATION_INSTANCESTATECORRUPTED = "UnsupportedOperation.InstanceStateCorrupted"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERRESCUEMODE = "UnsupportedOperation.InstanceStateEnterRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERSERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateEnterServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITRESCUEMODE = "UnsupportedOperation.InstanceStateExitRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITSERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateExitServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATEFREEZING = "UnsupportedOperation.InstanceStateFreezing"
//  UNSUPPORTEDOPERATION_INSTANCESTATEISOLATING = "UnsupportedOperation.InstanceStateIsolating"
//  UNSUPPORTEDOPERATION_INSTANCESTATEPENDING = "UnsupportedOperation.InstanceStatePending"
//  UNSUPPORTEDOPERATION_INSTANCESTATEREBOOTING = "UnsupportedOperation.InstanceStateRebooting"
//  UNSUPPORTEDOPERATION_INSTANCESTATESERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATESHUTDOWN = "UnsupportedOperation.InstanceStateShutdown"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTARTING = "UnsupportedOperation.InstanceStateStarting"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPED = "UnsupportedOperation.InstanceStateStopped"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPING = "UnsupportedOperation.InstanceStateStopping"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATING = "UnsupportedOperation.InstanceStateTerminating"
//  UNSUPPORTEDOPERATION_STOPPEDMODESTOPCHARGING = "UnsupportedOperation.StoppedModeStopCharging"
func (c *Client) RebootInstancesWithContext(ctx context.Context, request *RebootInstancesRequest) (response *RebootInstancesResponse, err error) {
    if request == nil {
        request = NewRebootInstancesRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "RebootInstances")
    
    if c.GetCredential() == nil {
        return nil, errors.New("RebootInstances require credential")
    }

    request.SetContext(ctx)
    
    response = NewRebootInstancesResponse()
    err = c.Send(request, response)
    return
}

func NewRemoveChcAssistVpcRequest() (request *RemoveChcAssistVpcRequest) {
    request = &RemoveChcAssistVpcRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "RemoveChcAssistVpc")
    
    
    return
}

func NewRemoveChcAssistVpcResponse() (response *RemoveChcAssistVpcResponse) {
    response = &RemoveChcAssistVpcResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// RemoveChcAssistVpc
// This API is used to remove the out-of-band network and deployment network of a CHC host.
//
// error code that may be returned:
//  ACCOUNTQUALIFICATIONRESTRICTIONS = "AccountQualificationRestrictions"
//  INVALIDHOST_NOTSUPPORTED = "InvalidHost.NotSupported"
//  INVALIDPARAMETERVALUE_CHCHOSTSNOTFOUND = "InvalidParameterValue.ChcHostsNotFound"
//  INVALIDPARAMETERVALUE_INCORRECTFORMAT = "InvalidParameterValue.IncorrectFormat"
func (c *Client) RemoveChcAssistVpc(request *RemoveChcAssistVpcRequest) (response *RemoveChcAssistVpcResponse, err error) {
    return c.RemoveChcAssistVpcWithContext(context.Background(), request)
}

// RemoveChcAssistVpc
// This API is used to remove the out-of-band network and deployment network of a CHC host.
//
// error code that may be returned:
//  ACCOUNTQUALIFICATIONRESTRICTIONS = "AccountQualificationRestrictions"
//  INVALIDHOST_NOTSUPPORTED = "InvalidHost.NotSupported"
//  INVALIDPARAMETERVALUE_CHCHOSTSNOTFOUND = "InvalidParameterValue.ChcHostsNotFound"
//  INVALIDPARAMETERVALUE_INCORRECTFORMAT = "InvalidParameterValue.IncorrectFormat"
func (c *Client) RemoveChcAssistVpcWithContext(ctx context.Context, request *RemoveChcAssistVpcRequest) (response *RemoveChcAssistVpcResponse, err error) {
    if request == nil {
        request = NewRemoveChcAssistVpcRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "RemoveChcAssistVpc")
    
    if c.GetCredential() == nil {
        return nil, errors.New("RemoveChcAssistVpc require credential")
    }

    request.SetContext(ctx)
    
    response = NewRemoveChcAssistVpcResponse()
    err = c.Send(request, response)
    return
}

func NewRemoveChcDeployVpcRequest() (request *RemoveChcDeployVpcRequest) {
    request = &RemoveChcDeployVpcRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "RemoveChcDeployVpc")
    
    
    return
}

func NewRemoveChcDeployVpcResponse() (response *RemoveChcDeployVpcResponse) {
    response = &RemoveChcDeployVpcResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// RemoveChcDeployVpc
// This API is used to remove the deployment network of a CHC host.
//
// error code that may be returned:
//  ACCOUNTQUALIFICATIONRESTRICTIONS = "AccountQualificationRestrictions"
//  INVALIDHOST_NOTSUPPORTED = "InvalidHost.NotSupported"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_CHCHOSTSNOTFOUND = "InvalidParameterValue.ChcHostsNotFound"
//  INVALIDPARAMETERVALUE_INCORRECTFORMAT = "InvalidParameterValue.IncorrectFormat"
func (c *Client) RemoveChcDeployVpc(request *RemoveChcDeployVpcRequest) (response *RemoveChcDeployVpcResponse, err error) {
    return c.RemoveChcDeployVpcWithContext(context.Background(), request)
}

// RemoveChcDeployVpc
// This API is used to remove the deployment network of a CHC host.
//
// error code that may be returned:
//  ACCOUNTQUALIFICATIONRESTRICTIONS = "AccountQualificationRestrictions"
//  INVALIDHOST_NOTSUPPORTED = "InvalidHost.NotSupported"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE_CHCHOSTSNOTFOUND = "InvalidParameterValue.ChcHostsNotFound"
//  INVALIDPARAMETERVALUE_INCORRECTFORMAT = "InvalidParameterValue.IncorrectFormat"
func (c *Client) RemoveChcDeployVpcWithContext(ctx context.Context, request *RemoveChcDeployVpcRequest) (response *RemoveChcDeployVpcResponse, err error) {
    if request == nil {
        request = NewRemoveChcDeployVpcRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "RemoveChcDeployVpc")
    
    if c.GetCredential() == nil {
        return nil, errors.New("RemoveChcDeployVpc require credential")
    }

    request.SetContext(ctx)
    
    response = NewRemoveChcDeployVpcResponse()
    err = c.Send(request, response)
    return
}

func NewRenewInstancesRequest() (request *RenewInstancesRequest) {
    request = &RenewInstancesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "RenewInstances")
    
    
    return
}

func NewRenewInstancesResponse() (response *RenewInstancesResponse) {
    response = &RenewInstancesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// RenewInstances
// This API is used to renew annual and monthly subscription instances.
//
// 
//
// This API is used to operate on monthly subscription instances only.
//
// This API is used to ensure your account balance is sufficient for renewal. You can check the balance via the DescribeAccountBalance API (https://www.tencentcloud.comom/document/product/555/20253?from_cn_redirect=1).
//
// * You can query the result of the instance operation by calling the API [DescribeInstances](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1#.E7.A4.BA.E4.BE.8B3-.E6.9F.A5.E8.AF.A2.E5.AE.9E.E4.BE.8B.E7.9A.84.E6.9C.80.E6.96.B0.E6.93.8D.E4.BD.9C.E6.83.85.E5.86.B5). If the latest operation status (LatestOperationState) of the instance is **SUCCESS**, the operation is successful.
//
// error code that may be returned:
//  FAILEDOPERATION_INVALIDINSTANCEAPPLICATIONROLE = "FailedOperation.InvalidInstanceApplicationRole"
//  FAILEDOPERATION_INVALIDINSTANCEAPPLICATIONROLEEMR = "FailedOperation.InvalidInstanceApplicationRoleEmr"
//  INTERNALERROR_TRADEUNKNOWNERROR = "InternalError.TradeUnknownError"
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDACCOUNT_INSUFFICIENTBALANCE = "InvalidAccount.InsufficientBalance"
//  INVALIDACCOUNT_UNPAIDORDER = "InvalidAccount.UnpaidOrder"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_INSTANCENOTSUPPORTEDMIXPRICINGMODEL = "InvalidParameterValue.InstanceNotSupportedMixPricingModel"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPERIOD = "InvalidPeriod"
//  MISSINGPARAMETER = "MissingParameter"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  UNSUPPORTEDOPERATION_INSTANCECHARGETYPE = "UnsupportedOperation.InstanceChargeType"
//  UNSUPPORTEDOPERATION_INSTANCEMIXEDZONETYPE = "UnsupportedOperation.InstanceMixedZoneType"
//  UNSUPPORTEDOPERATION_INSTANCESTATEBANNING = "UnsupportedOperation.InstanceStateBanning"
//  UNSUPPORTEDOPERATION_INSTANCESTATEFREEZING = "UnsupportedOperation.InstanceStateFreezing"
//  UNSUPPORTEDOPERATION_INSTANCESTATEPENDING = "UnsupportedOperation.InstanceStatePending"
//  UNSUPPORTEDOPERATION_INSTANCESTATEREBOOTING = "UnsupportedOperation.InstanceStateRebooting"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTARTING = "UnsupportedOperation.InstanceStateStarting"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPING = "UnsupportedOperation.InstanceStateStopping"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATING = "UnsupportedOperation.InstanceStateTerminating"
//  UNSUPPORTEDOPERATION_INVALIDDISKBACKUPQUOTA = "UnsupportedOperation.InvalidDiskBackupQuota"
//  UNSUPPORTEDOPERATION_SPECIALINSTANCETYPE = "UnsupportedOperation.SpecialInstanceType"
func (c *Client) RenewInstances(request *RenewInstancesRequest) (response *RenewInstancesResponse, err error) {
    return c.RenewInstancesWithContext(context.Background(), request)
}

// RenewInstances
// This API is used to renew annual and monthly subscription instances.
//
// 
//
// This API is used to operate on monthly subscription instances only.
//
// This API is used to ensure your account balance is sufficient for renewal. You can check the balance via the DescribeAccountBalance API (https://www.tencentcloud.comom/document/product/555/20253?from_cn_redirect=1).
//
// * You can query the result of the instance operation by calling the API [DescribeInstances](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1#.E7.A4.BA.E4.BE.8B3-.E6.9F.A5.E8.AF.A2.E5.AE.9E.E4.BE.8B.E7.9A.84.E6.9C.80.E6.96.B0.E6.93.8D.E4.BD.9C.E6.83.85.E5.86.B5). If the latest operation status (LatestOperationState) of the instance is **SUCCESS**, the operation is successful.
//
// error code that may be returned:
//  FAILEDOPERATION_INVALIDINSTANCEAPPLICATIONROLE = "FailedOperation.InvalidInstanceApplicationRole"
//  FAILEDOPERATION_INVALIDINSTANCEAPPLICATIONROLEEMR = "FailedOperation.InvalidInstanceApplicationRoleEmr"
//  INTERNALERROR_TRADEUNKNOWNERROR = "InternalError.TradeUnknownError"
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDACCOUNT_INSUFFICIENTBALANCE = "InvalidAccount.InsufficientBalance"
//  INVALIDACCOUNT_UNPAIDORDER = "InvalidAccount.UnpaidOrder"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_INSTANCENOTSUPPORTEDMIXPRICINGMODEL = "InvalidParameterValue.InstanceNotSupportedMixPricingModel"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPERIOD = "InvalidPeriod"
//  MISSINGPARAMETER = "MissingParameter"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  UNSUPPORTEDOPERATION_INSTANCECHARGETYPE = "UnsupportedOperation.InstanceChargeType"
//  UNSUPPORTEDOPERATION_INSTANCEMIXEDZONETYPE = "UnsupportedOperation.InstanceMixedZoneType"
//  UNSUPPORTEDOPERATION_INSTANCESTATEBANNING = "UnsupportedOperation.InstanceStateBanning"
//  UNSUPPORTEDOPERATION_INSTANCESTATEFREEZING = "UnsupportedOperation.InstanceStateFreezing"
//  UNSUPPORTEDOPERATION_INSTANCESTATEPENDING = "UnsupportedOperation.InstanceStatePending"
//  UNSUPPORTEDOPERATION_INSTANCESTATEREBOOTING = "UnsupportedOperation.InstanceStateRebooting"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTARTING = "UnsupportedOperation.InstanceStateStarting"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPING = "UnsupportedOperation.InstanceStateStopping"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATING = "UnsupportedOperation.InstanceStateTerminating"
//  UNSUPPORTEDOPERATION_INVALIDDISKBACKUPQUOTA = "UnsupportedOperation.InvalidDiskBackupQuota"
//  UNSUPPORTEDOPERATION_SPECIALINSTANCETYPE = "UnsupportedOperation.SpecialInstanceType"
func (c *Client) RenewInstancesWithContext(ctx context.Context, request *RenewInstancesRequest) (response *RenewInstancesResponse, err error) {
    if request == nil {
        request = NewRenewInstancesRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "RenewInstances")
    
    if c.GetCredential() == nil {
        return nil, errors.New("RenewInstances require credential")
    }

    request.SetContext(ctx)
    
    response = NewRenewInstancesResponse()
    err = c.Send(request, response)
    return
}

func NewResetInstanceRequest() (request *ResetInstanceRequest) {
    request = &ResetInstanceRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "ResetInstance")
    
    
    return
}

func NewResetInstanceResponse() (response *ResetInstanceResponse) {
    response = &ResetInstanceResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// ResetInstance
// This API (ResetInstance) is used to reinstall the operating system on a specified instance.
//
// 
//
// 
//
// 
//
// * If you have specified the parameter `ImageId`, use the specified image for reinstallation. Otherwise, perform reinstallation based on the image currently used by the instance.
//
// * The system disk will be formatted and reset. Ensure that there are no important files in the system disk.
//
// * If you do not specify a password, a random password will be sent via Message Center.
//
// * Currently, only instances with a [system disk type](https://intl.cloud.tencent.com/document/api/213/9452?from_cn_redirect=1#SystemDisk) of `CLOUD_BASIC`, `CLOUD_PREMIUM`, `CLOUD_SSD`, or `CLOUD_BSSD` are supported for implementing operating system switching through this API.
//
// * You can query the result of the instance operation by calling the API [DescribeInstances](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1#.E7.A4.BA.E4.BE.8B3-.E6.9F.A5.E8.AF.A2.E5.AE.9E.E4.BE.8B.E7.9A.84.E6.9C.80.E6.96.B0.E6.93.8D.E4.BD.9C.E6.83.85.E5.86.B5). If the latest operation status (LatestOperationState) of the instance is **SUCCESS**, the operation is successful.
//
// error code that may be returned:
//  FAILEDOPERATION_INVALIDINSTANCEAPPLICATIONROLEEMR = "FailedOperation.InvalidInstanceApplicationRoleEmr"
//  INTERNALERROR_TRADEUNKNOWNERROR = "InternalError.TradeUnknownError"
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDACCOUNT_INSUFFICIENTBALANCE = "InvalidAccount.InsufficientBalance"
//  INVALIDACCOUNT_UNPAIDORDER = "InvalidAccount.UnpaidOrder"
//  INVALIDIMAGEID_MALFORMED = "InvalidImageId.Malformed"
//  INVALIDIMAGEID_NOTFOUND = "InvalidImageId.NotFound"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDINSTANCETYPE_MALFORMED = "InvalidInstanceType.Malformed"
//  INVALIDKEYPAIRID_MALFORMED = "InvalidKeyPairId.Malformed"
//  INVALIDPARAMETER_HOSTNAMEILLEGAL = "InvalidParameter.HostNameIllegal"
//  INVALIDPARAMETER_INSTANCEIMAGENOTSUPPORT = "InvalidParameter.InstanceImageNotSupport"
//  INVALIDPARAMETER_PARAMETERCONFLICT = "InvalidParameter.ParameterConflict"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_ILLEGALHOSTNAME = "InvalidParameterValue.IllegalHostName"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_INSTANCETYPENOTFOUND = "InvalidParameterValue.InstanceTypeNotFound"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEFORGIVENINSTANCETYPE = "InvalidParameterValue.InvalidImageForGivenInstanceType"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEFORMAT = "InvalidParameterValue.InvalidImageFormat"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEID = "InvalidParameterValue.InvalidImageId"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEIDFORRETSETINSTANCE = "InvalidParameterValue.InvalidImageIdForRetsetInstance"
//  INVALIDPARAMETERVALUE_INVALIDIMAGESTATE = "InvalidParameterValue.InvalidImageState"
//  INVALIDPARAMETERVALUE_INVALIDPASSWORD = "InvalidParameterValue.InvalidPassword"
//  INVALIDPARAMETERVALUE_INVALIDUSERDATAFORMAT = "InvalidParameterValue.InvalidUserDataFormat"
//  INVALIDPARAMETERVALUE_KEYPAIRNOTFOUND = "InvalidParameterValue.KeyPairNotFound"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_REQUIREDLOCATIONIMAGE = "InvalidParameterValue.RequiredLocationImage"
//  INVALIDPARAMETERVALUE_ZONENOTSUPPORTED = "InvalidParameterValue.ZoneNotSupported"
//  MISSINGPARAMETER = "MissingParameter"
//  MISSINGPARAMETER_MONITORSERVICE = "MissingParameter.MonitorService"
//  MUTEXOPERATION_TASKRUNNING = "MutexOperation.TaskRunning"
//  OPERATIONDENIED_CHCINSTALLCLOUDIMAGEWITHOUTDEPLOYNETWORK = "OperationDenied.ChcInstallCloudImageWithoutDeployNetwork"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  RESOURCEINUSE = "ResourceInUse"
//  RESOURCEINSUFFICIENT_CLOUDDISKSOLDOUT = "ResourceInsufficient.CloudDiskSoldOut"
//  RESOURCENOTFOUND_KEYPAIRNOTFOUND = "ResourceNotFound.KeyPairNotFound"
//  RESOURCESSOLDOUT_SPECIFIEDINSTANCETYPE = "ResourcesSoldOut.SpecifiedInstanceType"
//  UNAUTHORIZEDOPERATION_INVALIDTOKEN = "UnauthorizedOperation.InvalidToken"
//  UNAUTHORIZEDOPERATION_MFAEXPIRED = "UnauthorizedOperation.MFAExpired"
//  UNAUTHORIZEDOPERATION_MFANOTFOUND = "UnauthorizedOperation.MFANotFound"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
//  UNSUPPORTEDOPERATION_INSTANCECHARGETYPE = "UnsupportedOperation.InstanceChargeType"
//  UNSUPPORTEDOPERATION_INSTANCESTATECORRUPTED = "UnsupportedOperation.InstanceStateCorrupted"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERRESCUEMODE = "UnsupportedOperation.InstanceStateEnterRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERSERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateEnterServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITRESCUEMODE = "UnsupportedOperation.InstanceStateExitRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITSERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateExitServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATEFREEZING = "UnsupportedOperation.InstanceStateFreezing"
//  UNSUPPORTEDOPERATION_INSTANCESTATEISOLATING = "UnsupportedOperation.InstanceStateIsolating"
//  UNSUPPORTEDOPERATION_INSTANCESTATEPENDING = "UnsupportedOperation.InstanceStatePending"
//  UNSUPPORTEDOPERATION_INSTANCESTATEREBOOTING = "UnsupportedOperation.InstanceStateRebooting"
//  UNSUPPORTEDOPERATION_INSTANCESTATERESCUEMODE = "UnsupportedOperation.InstanceStateRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATESERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATESHUTDOWN = "UnsupportedOperation.InstanceStateShutdown"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTARTING = "UnsupportedOperation.InstanceStateStarting"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPING = "UnsupportedOperation.InstanceStateStopping"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATED = "UnsupportedOperation.InstanceStateTerminated"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATING = "UnsupportedOperation.InstanceStateTerminating"
//  UNSUPPORTEDOPERATION_INSTANCETYPENOTSUPPORTGRIDLICENCE = "UnsupportedOperation.InstanceTypeNotSupportGridLicence"
//  UNSUPPORTEDOPERATION_INVALIDIMAGELICENSETYPEFORRESET = "UnsupportedOperation.InvalidImageLicenseTypeForReset"
//  UNSUPPORTEDOPERATION_KEYPAIRUNSUPPORTEDWINDOWS = "UnsupportedOperation.KeyPairUnsupportedWindows"
//  UNSUPPORTEDOPERATION_MODIFYENCRYPTIONNOTSUPPORTED = "UnsupportedOperation.ModifyEncryptionNotSupported"
//  UNSUPPORTEDOPERATION_RAWLOCALDISKINSREINSTALLTOQCOW2 = "UnsupportedOperation.RawLocalDiskInsReinstalltoQcow2"
//  UNSUPPORTEDOPERATION_STOPPEDMODESTOPCHARGING = "UnsupportedOperation.StoppedModeStopCharging"
func (c *Client) ResetInstance(request *ResetInstanceRequest) (response *ResetInstanceResponse, err error) {
    return c.ResetInstanceWithContext(context.Background(), request)
}

// ResetInstance
// This API (ResetInstance) is used to reinstall the operating system on a specified instance.
//
// 
//
// 
//
// 
//
// * If you have specified the parameter `ImageId`, use the specified image for reinstallation. Otherwise, perform reinstallation based on the image currently used by the instance.
//
// * The system disk will be formatted and reset. Ensure that there are no important files in the system disk.
//
// * If you do not specify a password, a random password will be sent via Message Center.
//
// * Currently, only instances with a [system disk type](https://intl.cloud.tencent.com/document/api/213/9452?from_cn_redirect=1#SystemDisk) of `CLOUD_BASIC`, `CLOUD_PREMIUM`, `CLOUD_SSD`, or `CLOUD_BSSD` are supported for implementing operating system switching through this API.
//
// * You can query the result of the instance operation by calling the API [DescribeInstances](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1#.E7.A4.BA.E4.BE.8B3-.E6.9F.A5.E8.AF.A2.E5.AE.9E.E4.BE.8B.E7.9A.84.E6.9C.80.E6.96.B0.E6.93.8D.E4.BD.9C.E6.83.85.E5.86.B5). If the latest operation status (LatestOperationState) of the instance is **SUCCESS**, the operation is successful.
//
// error code that may be returned:
//  FAILEDOPERATION_INVALIDINSTANCEAPPLICATIONROLEEMR = "FailedOperation.InvalidInstanceApplicationRoleEmr"
//  INTERNALERROR_TRADEUNKNOWNERROR = "InternalError.TradeUnknownError"
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDACCOUNT_INSUFFICIENTBALANCE = "InvalidAccount.InsufficientBalance"
//  INVALIDACCOUNT_UNPAIDORDER = "InvalidAccount.UnpaidOrder"
//  INVALIDIMAGEID_MALFORMED = "InvalidImageId.Malformed"
//  INVALIDIMAGEID_NOTFOUND = "InvalidImageId.NotFound"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDINSTANCETYPE_MALFORMED = "InvalidInstanceType.Malformed"
//  INVALIDKEYPAIRID_MALFORMED = "InvalidKeyPairId.Malformed"
//  INVALIDPARAMETER_HOSTNAMEILLEGAL = "InvalidParameter.HostNameIllegal"
//  INVALIDPARAMETER_INSTANCEIMAGENOTSUPPORT = "InvalidParameter.InstanceImageNotSupport"
//  INVALIDPARAMETER_PARAMETERCONFLICT = "InvalidParameter.ParameterConflict"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_ILLEGALHOSTNAME = "InvalidParameterValue.IllegalHostName"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_INSTANCETYPENOTFOUND = "InvalidParameterValue.InstanceTypeNotFound"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEFORGIVENINSTANCETYPE = "InvalidParameterValue.InvalidImageForGivenInstanceType"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEFORMAT = "InvalidParameterValue.InvalidImageFormat"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEID = "InvalidParameterValue.InvalidImageId"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEIDFORRETSETINSTANCE = "InvalidParameterValue.InvalidImageIdForRetsetInstance"
//  INVALIDPARAMETERVALUE_INVALIDIMAGESTATE = "InvalidParameterValue.InvalidImageState"
//  INVALIDPARAMETERVALUE_INVALIDPASSWORD = "InvalidParameterValue.InvalidPassword"
//  INVALIDPARAMETERVALUE_INVALIDUSERDATAFORMAT = "InvalidParameterValue.InvalidUserDataFormat"
//  INVALIDPARAMETERVALUE_KEYPAIRNOTFOUND = "InvalidParameterValue.KeyPairNotFound"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_REQUIREDLOCATIONIMAGE = "InvalidParameterValue.RequiredLocationImage"
//  INVALIDPARAMETERVALUE_ZONENOTSUPPORTED = "InvalidParameterValue.ZoneNotSupported"
//  MISSINGPARAMETER = "MissingParameter"
//  MISSINGPARAMETER_MONITORSERVICE = "MissingParameter.MonitorService"
//  MUTEXOPERATION_TASKRUNNING = "MutexOperation.TaskRunning"
//  OPERATIONDENIED_CHCINSTALLCLOUDIMAGEWITHOUTDEPLOYNETWORK = "OperationDenied.ChcInstallCloudImageWithoutDeployNetwork"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  RESOURCEINUSE = "ResourceInUse"
//  RESOURCEINSUFFICIENT_CLOUDDISKSOLDOUT = "ResourceInsufficient.CloudDiskSoldOut"
//  RESOURCENOTFOUND_KEYPAIRNOTFOUND = "ResourceNotFound.KeyPairNotFound"
//  RESOURCESSOLDOUT_SPECIFIEDINSTANCETYPE = "ResourcesSoldOut.SpecifiedInstanceType"
//  UNAUTHORIZEDOPERATION_INVALIDTOKEN = "UnauthorizedOperation.InvalidToken"
//  UNAUTHORIZEDOPERATION_MFAEXPIRED = "UnauthorizedOperation.MFAExpired"
//  UNAUTHORIZEDOPERATION_MFANOTFOUND = "UnauthorizedOperation.MFANotFound"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
//  UNSUPPORTEDOPERATION_INSTANCECHARGETYPE = "UnsupportedOperation.InstanceChargeType"
//  UNSUPPORTEDOPERATION_INSTANCESTATECORRUPTED = "UnsupportedOperation.InstanceStateCorrupted"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERRESCUEMODE = "UnsupportedOperation.InstanceStateEnterRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERSERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateEnterServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITRESCUEMODE = "UnsupportedOperation.InstanceStateExitRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITSERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateExitServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATEFREEZING = "UnsupportedOperation.InstanceStateFreezing"
//  UNSUPPORTEDOPERATION_INSTANCESTATEISOLATING = "UnsupportedOperation.InstanceStateIsolating"
//  UNSUPPORTEDOPERATION_INSTANCESTATEPENDING = "UnsupportedOperation.InstanceStatePending"
//  UNSUPPORTEDOPERATION_INSTANCESTATEREBOOTING = "UnsupportedOperation.InstanceStateRebooting"
//  UNSUPPORTEDOPERATION_INSTANCESTATERESCUEMODE = "UnsupportedOperation.InstanceStateRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATESERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATESHUTDOWN = "UnsupportedOperation.InstanceStateShutdown"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTARTING = "UnsupportedOperation.InstanceStateStarting"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPING = "UnsupportedOperation.InstanceStateStopping"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATED = "UnsupportedOperation.InstanceStateTerminated"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATING = "UnsupportedOperation.InstanceStateTerminating"
//  UNSUPPORTEDOPERATION_INSTANCETYPENOTSUPPORTGRIDLICENCE = "UnsupportedOperation.InstanceTypeNotSupportGridLicence"
//  UNSUPPORTEDOPERATION_INVALIDIMAGELICENSETYPEFORRESET = "UnsupportedOperation.InvalidImageLicenseTypeForReset"
//  UNSUPPORTEDOPERATION_KEYPAIRUNSUPPORTEDWINDOWS = "UnsupportedOperation.KeyPairUnsupportedWindows"
//  UNSUPPORTEDOPERATION_MODIFYENCRYPTIONNOTSUPPORTED = "UnsupportedOperation.ModifyEncryptionNotSupported"
//  UNSUPPORTEDOPERATION_RAWLOCALDISKINSREINSTALLTOQCOW2 = "UnsupportedOperation.RawLocalDiskInsReinstalltoQcow2"
//  UNSUPPORTEDOPERATION_STOPPEDMODESTOPCHARGING = "UnsupportedOperation.StoppedModeStopCharging"
func (c *Client) ResetInstanceWithContext(ctx context.Context, request *ResetInstanceRequest) (response *ResetInstanceResponse, err error) {
    if request == nil {
        request = NewResetInstanceRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "ResetInstance")
    
    if c.GetCredential() == nil {
        return nil, errors.New("ResetInstance require credential")
    }

    request.SetContext(ctx)
    
    response = NewResetInstanceResponse()
    err = c.Send(request, response)
    return
}

func NewResetInstancesInternetMaxBandwidthRequest() (request *ResetInstancesInternetMaxBandwidthRequest) {
    request = &ResetInstancesInternetMaxBandwidthRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "ResetInstancesInternetMaxBandwidth")
    
    
    return
}

func NewResetInstancesInternetMaxBandwidthResponse() (response *ResetInstancesInternetMaxBandwidthResponse) {
    response = &ResetInstancesInternetMaxBandwidthResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// ResetInstancesInternetMaxBandwidth
// This API is used to change the public bandwidth cap of an instance.
//
// 
//
// * The allowed bandwidth cap varies for different models. For details, see [Purchasing Network Bandwidth](https://intl.cloud.tencent.com/document/product/213/509?from_cn_redirect=1).
//
// * For bandwidth billed by the `TRAFFIC_POSTPAID_BY_HOUR` method, changing the bandwidth cap through this API takes effect in real time. Users can increase or reduce bandwidth within applicable limits.
//
// error code that may be returned:
//  FAILEDOPERATION_NOTFOUNDEIP = "FailedOperation.NotFoundEIP"
//  INTERNALERROR_TRADEUNKNOWNERROR = "InternalError.TradeUnknownError"
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDACCOUNT_INSUFFICIENTBALANCE = "InvalidAccount.InsufficientBalance"
//  INVALIDACCOUNT_UNPAIDORDER = "InvalidAccount.UnpaidOrder"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_BANDWIDTHPACKAGEIDMALFORMED = "InvalidParameterValue.BandwidthPackageIdMalformed"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_RANGE = "InvalidParameterValue.Range"
//  INVALIDPERMISSION = "InvalidPermission"
//  MISSINGPARAMETER = "MissingParameter"
//  MUTEXOPERATION_TASKRUNNING = "MutexOperation.TaskRunning"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  UNSUPPORTEDOPERATION_BANDWIDTHPACKAGEIDNOTSUPPORTED = "UnsupportedOperation.BandwidthPackageIdNotSupported"
//  UNSUPPORTEDOPERATION_EDGEZONEINSTANCE = "UnsupportedOperation.EdgeZoneInstance"
//  UNSUPPORTEDOPERATION_INSTANCECHARGETYPE = "UnsupportedOperation.InstanceChargeType"
//  UNSUPPORTEDOPERATION_INSTANCEREINSTALLFAILED = "UnsupportedOperation.InstanceReinstallFailed"
//  UNSUPPORTEDOPERATION_INSTANCESTATEBANNING = "UnsupportedOperation.InstanceStateBanning"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERSERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateEnterServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITSERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateExitServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATEFREEZING = "UnsupportedOperation.InstanceStateFreezing"
//  UNSUPPORTEDOPERATION_INSTANCESTATEISOLATING = "UnsupportedOperation.InstanceStateIsolating"
//  UNSUPPORTEDOPERATION_INSTANCESTATEPENDING = "UnsupportedOperation.InstanceStatePending"
//  UNSUPPORTEDOPERATION_INSTANCESTATEREBOOTING = "UnsupportedOperation.InstanceStateRebooting"
//  UNSUPPORTEDOPERATION_INSTANCESTATERESCUEMODE = "UnsupportedOperation.InstanceStateRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATESHUTDOWN = "UnsupportedOperation.InstanceStateShutdown"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTARTING = "UnsupportedOperation.InstanceStateStarting"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPING = "UnsupportedOperation.InstanceStateStopping"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATING = "UnsupportedOperation.InstanceStateTerminating"
//  UNSUPPORTEDOPERATION_STOPPEDMODESTOPCHARGING = "UnsupportedOperation.StoppedModeStopCharging"
func (c *Client) ResetInstancesInternetMaxBandwidth(request *ResetInstancesInternetMaxBandwidthRequest) (response *ResetInstancesInternetMaxBandwidthResponse, err error) {
    return c.ResetInstancesInternetMaxBandwidthWithContext(context.Background(), request)
}

// ResetInstancesInternetMaxBandwidth
// This API is used to change the public bandwidth cap of an instance.
//
// 
//
// * The allowed bandwidth cap varies for different models. For details, see [Purchasing Network Bandwidth](https://intl.cloud.tencent.com/document/product/213/509?from_cn_redirect=1).
//
// * For bandwidth billed by the `TRAFFIC_POSTPAID_BY_HOUR` method, changing the bandwidth cap through this API takes effect in real time. Users can increase or reduce bandwidth within applicable limits.
//
// error code that may be returned:
//  FAILEDOPERATION_NOTFOUNDEIP = "FailedOperation.NotFoundEIP"
//  INTERNALERROR_TRADEUNKNOWNERROR = "InternalError.TradeUnknownError"
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDACCOUNT_INSUFFICIENTBALANCE = "InvalidAccount.InsufficientBalance"
//  INVALIDACCOUNT_UNPAIDORDER = "InvalidAccount.UnpaidOrder"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_BANDWIDTHPACKAGEIDMALFORMED = "InvalidParameterValue.BandwidthPackageIdMalformed"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_RANGE = "InvalidParameterValue.Range"
//  INVALIDPERMISSION = "InvalidPermission"
//  MISSINGPARAMETER = "MissingParameter"
//  MUTEXOPERATION_TASKRUNNING = "MutexOperation.TaskRunning"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  UNSUPPORTEDOPERATION_BANDWIDTHPACKAGEIDNOTSUPPORTED = "UnsupportedOperation.BandwidthPackageIdNotSupported"
//  UNSUPPORTEDOPERATION_EDGEZONEINSTANCE = "UnsupportedOperation.EdgeZoneInstance"
//  UNSUPPORTEDOPERATION_INSTANCECHARGETYPE = "UnsupportedOperation.InstanceChargeType"
//  UNSUPPORTEDOPERATION_INSTANCEREINSTALLFAILED = "UnsupportedOperation.InstanceReinstallFailed"
//  UNSUPPORTEDOPERATION_INSTANCESTATEBANNING = "UnsupportedOperation.InstanceStateBanning"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERSERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateEnterServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITSERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateExitServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATEFREEZING = "UnsupportedOperation.InstanceStateFreezing"
//  UNSUPPORTEDOPERATION_INSTANCESTATEISOLATING = "UnsupportedOperation.InstanceStateIsolating"
//  UNSUPPORTEDOPERATION_INSTANCESTATEPENDING = "UnsupportedOperation.InstanceStatePending"
//  UNSUPPORTEDOPERATION_INSTANCESTATEREBOOTING = "UnsupportedOperation.InstanceStateRebooting"
//  UNSUPPORTEDOPERATION_INSTANCESTATERESCUEMODE = "UnsupportedOperation.InstanceStateRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATESHUTDOWN = "UnsupportedOperation.InstanceStateShutdown"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTARTING = "UnsupportedOperation.InstanceStateStarting"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPING = "UnsupportedOperation.InstanceStateStopping"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATING = "UnsupportedOperation.InstanceStateTerminating"
//  UNSUPPORTEDOPERATION_STOPPEDMODESTOPCHARGING = "UnsupportedOperation.StoppedModeStopCharging"
func (c *Client) ResetInstancesInternetMaxBandwidthWithContext(ctx context.Context, request *ResetInstancesInternetMaxBandwidthRequest) (response *ResetInstancesInternetMaxBandwidthResponse, err error) {
    if request == nil {
        request = NewResetInstancesInternetMaxBandwidthRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "ResetInstancesInternetMaxBandwidth")
    
    if c.GetCredential() == nil {
        return nil, errors.New("ResetInstancesInternetMaxBandwidth require credential")
    }

    request.SetContext(ctx)
    
    response = NewResetInstancesInternetMaxBandwidthResponse()
    err = c.Send(request, response)
    return
}

func NewResetInstancesPasswordRequest() (request *ResetInstancesPasswordRequest) {
    request = &ResetInstancesPasswordRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "ResetInstancesPassword")
    
    
    return
}

func NewResetInstancesPasswordResponse() (response *ResetInstancesPasswordResponse) {
    response = &ResetInstancesPasswordResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// ResetInstancesPassword
// This API is used to reset the password of the operating system instances to a user-specified password.
//
// 
//
// * To modify the password of the administrator account: the name of the administrator account varies with the operating system. In Windows, it is `Administrator`; in Ubuntu, it is `ubuntu`; in Linux, it is `root`.
//
// * To reset the password of a running instance, you need to set the parameter `ForceStop` to `True` for a forced shutdown. If not, only passwords of stopped instances can be reset.
//
// * Batch operations are supported. You can reset the passwords of up to 100 instances to the same value once.
//
// * You can call the [DescribeInstances](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1#.E7.A4.BA.E4.BE.8B3-.E6.9F.A5.E8.AF.A2.E5.AE.9E.E4.BE.8B.E7.9A.84.E6.9C.80.E6.96.B0.E6.93.8D.E4.BD.9C.E6.83.85.E5.86.B5) API and find the result of the operation in the response parameter `LatestOperationState`. If the value is `SUCCESS`, the operation is successful.
//
// error code that may be returned:
//  FAILEDOPERATION_INVALIDINSTANCEAPPLICATIONROLE = "FailedOperation.InvalidInstanceApplicationRole"
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_INVALIDPASSWORD = "InvalidParameterValue.InvalidPassword"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_TOOLONG = "InvalidParameterValue.TooLong"
//  MISSINGPARAMETER = "MissingParameter"
//  MUTEXOPERATION_TASKRUNNING = "MutexOperation.TaskRunning"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  UNAUTHORIZEDOPERATION_MFAEXPIRED = "UnauthorizedOperation.MFAExpired"
//  UNAUTHORIZEDOPERATION_MFANOTFOUND = "UnauthorizedOperation.MFANotFound"
//  UNSUPPORTEDOPERATION_INSTANCEREINSTALLFAILED = "UnsupportedOperation.InstanceReinstallFailed"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERRESCUEMODE = "UnsupportedOperation.InstanceStateEnterRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITRESCUEMODE = "UnsupportedOperation.InstanceStateExitRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEFREEZING = "UnsupportedOperation.InstanceStateFreezing"
//  UNSUPPORTEDOPERATION_INSTANCESTATEISOLATING = "UnsupportedOperation.InstanceStateIsolating"
//  UNSUPPORTEDOPERATION_INSTANCESTATEPENDING = "UnsupportedOperation.InstanceStatePending"
//  UNSUPPORTEDOPERATION_INSTANCESTATEREBOOTING = "UnsupportedOperation.InstanceStateRebooting"
//  UNSUPPORTEDOPERATION_INSTANCESTATERESCUEMODE = "UnsupportedOperation.InstanceStateRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATERUNNING = "UnsupportedOperation.InstanceStateRunning"
//  UNSUPPORTEDOPERATION_INSTANCESTATESERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATESHUTDOWN = "UnsupportedOperation.InstanceStateShutdown"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTARTING = "UnsupportedOperation.InstanceStateStarting"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPING = "UnsupportedOperation.InstanceStateStopping"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATING = "UnsupportedOperation.InstanceStateTerminating"
//  UNSUPPORTEDOPERATION_SPECIALINSTANCETYPE = "UnsupportedOperation.SpecialInstanceType"
//  UNSUPPORTEDOPERATION_STOPPEDMODESTOPCHARGING = "UnsupportedOperation.StoppedModeStopCharging"
func (c *Client) ResetInstancesPassword(request *ResetInstancesPasswordRequest) (response *ResetInstancesPasswordResponse, err error) {
    return c.ResetInstancesPasswordWithContext(context.Background(), request)
}

// ResetInstancesPassword
// This API is used to reset the password of the operating system instances to a user-specified password.
//
// 
//
// * To modify the password of the administrator account: the name of the administrator account varies with the operating system. In Windows, it is `Administrator`; in Ubuntu, it is `ubuntu`; in Linux, it is `root`.
//
// * To reset the password of a running instance, you need to set the parameter `ForceStop` to `True` for a forced shutdown. If not, only passwords of stopped instances can be reset.
//
// * Batch operations are supported. You can reset the passwords of up to 100 instances to the same value once.
//
// * You can call the [DescribeInstances](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1#.E7.A4.BA.E4.BE.8B3-.E6.9F.A5.E8.AF.A2.E5.AE.9E.E4.BE.8B.E7.9A.84.E6.9C.80.E6.96.B0.E6.93.8D.E4.BD.9C.E6.83.85.E5.86.B5) API and find the result of the operation in the response parameter `LatestOperationState`. If the value is `SUCCESS`, the operation is successful.
//
// error code that may be returned:
//  FAILEDOPERATION_INVALIDINSTANCEAPPLICATIONROLE = "FailedOperation.InvalidInstanceApplicationRole"
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_INVALIDPASSWORD = "InvalidParameterValue.InvalidPassword"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_TOOLONG = "InvalidParameterValue.TooLong"
//  MISSINGPARAMETER = "MissingParameter"
//  MUTEXOPERATION_TASKRUNNING = "MutexOperation.TaskRunning"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  UNAUTHORIZEDOPERATION_MFAEXPIRED = "UnauthorizedOperation.MFAExpired"
//  UNAUTHORIZEDOPERATION_MFANOTFOUND = "UnauthorizedOperation.MFANotFound"
//  UNSUPPORTEDOPERATION_INSTANCEREINSTALLFAILED = "UnsupportedOperation.InstanceReinstallFailed"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERRESCUEMODE = "UnsupportedOperation.InstanceStateEnterRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITRESCUEMODE = "UnsupportedOperation.InstanceStateExitRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEFREEZING = "UnsupportedOperation.InstanceStateFreezing"
//  UNSUPPORTEDOPERATION_INSTANCESTATEISOLATING = "UnsupportedOperation.InstanceStateIsolating"
//  UNSUPPORTEDOPERATION_INSTANCESTATEPENDING = "UnsupportedOperation.InstanceStatePending"
//  UNSUPPORTEDOPERATION_INSTANCESTATEREBOOTING = "UnsupportedOperation.InstanceStateRebooting"
//  UNSUPPORTEDOPERATION_INSTANCESTATERESCUEMODE = "UnsupportedOperation.InstanceStateRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATERUNNING = "UnsupportedOperation.InstanceStateRunning"
//  UNSUPPORTEDOPERATION_INSTANCESTATESERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATESHUTDOWN = "UnsupportedOperation.InstanceStateShutdown"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTARTING = "UnsupportedOperation.InstanceStateStarting"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPING = "UnsupportedOperation.InstanceStateStopping"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATING = "UnsupportedOperation.InstanceStateTerminating"
//  UNSUPPORTEDOPERATION_SPECIALINSTANCETYPE = "UnsupportedOperation.SpecialInstanceType"
//  UNSUPPORTEDOPERATION_STOPPEDMODESTOPCHARGING = "UnsupportedOperation.StoppedModeStopCharging"
func (c *Client) ResetInstancesPasswordWithContext(ctx context.Context, request *ResetInstancesPasswordRequest) (response *ResetInstancesPasswordResponse, err error) {
    if request == nil {
        request = NewResetInstancesPasswordRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "ResetInstancesPassword")
    
    if c.GetCredential() == nil {
        return nil, errors.New("ResetInstancesPassword require credential")
    }

    request.SetContext(ctx)
    
    response = NewResetInstancesPasswordResponse()
    err = c.Send(request, response)
    return
}

func NewResetInstancesTypeRequest() (request *ResetInstancesTypeRequest) {
    request = &ResetInstancesTypeRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "ResetInstancesType")
    
    
    return
}

func NewResetInstancesTypeResponse() (response *ResetInstancesTypeResponse) {
    response = &ResetInstancesTypeResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// ResetInstancesType
// This API is used to change the model of an instance.
//
// * You can only use this API to change the models of instances whose [system disk type](https://intl.cloud.tencent.com/document/api/213/9452?from_cn_redirect=1#block_device) is `CLOUD_BASIC`, `CLOUD_PREMIUM`, or `CLOUD_SSD`.
//
// * Currently, you cannot use this API to change the models of [CDH](https://intl.cloud.tencent.com/document/product/416?from_cn_redirect=1) instances.
//
// error code that may be returned:
//  FAILEDOPERATION_INVALIDINSTANCEAPPLICATIONROLE = "FailedOperation.InvalidInstanceApplicationRole"
//  FAILEDOPERATION_INVALIDINSTANCEAPPLICATIONROLEEMR = "FailedOperation.InvalidInstanceApplicationRoleEmr"
//  FAILEDOPERATION_PROMOTIONALPERIORESTRICTION = "FailedOperation.PromotionalPerioRestriction"
//  INTERNALERROR_TRADEUNKNOWNERROR = "InternalError.TradeUnknownError"
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDACCOUNT_INSUFFICIENTBALANCE = "InvalidAccount.InsufficientBalance"
//  INVALIDACCOUNT_UNPAIDORDER = "InvalidAccount.UnpaidOrder"
//  INVALIDHOSTID_MALFORMED = "InvalidHostId.Malformed"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDINSTANCETYPE_MALFORMED = "InvalidInstanceType.Malformed"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_HOSTIDSTATUSNOTSUPPORT = "InvalidParameter.HostIdStatusNotSupport"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_BASICNETWORKINSTANCEFAMILY = "InvalidParameterValue.BasicNetworkInstanceFamily"
//  INVALIDPARAMETERVALUE_GPUINSTANCEFAMILY = "InvalidParameterValue.GPUInstanceFamily"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_INVALIDGPUFAMILYCHANGE = "InvalidParameterValue.InvalidGPUFamilyChange"
//  INVALIDPARAMETERVALUE_INVALIDINSTANCESOURCE = "InvalidParameterValue.InvalidInstanceSource"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_ZONENOTSUPPORTED = "InvalidParameterValue.ZoneNotSupported"
//  INVALIDPERIOD = "InvalidPeriod"
//  INVALIDPERMISSION = "InvalidPermission"
//  INVALIDZONE_MISMATCHREGION = "InvalidZone.MismatchRegion"
//  LIMITEXCEEDED_EIPNUMLIMIT = "LimitExceeded.EipNumLimit"
//  LIMITEXCEEDED_ENILIMITINSTANCETYPE = "LimitExceeded.EniLimitInstanceType"
//  LIMITEXCEEDED_ENINUMLIMIT = "LimitExceeded.EniNumLimit"
//  LIMITEXCEEDED_INSTANCETYPEBANDWIDTH = "LimitExceeded.InstanceTypeBandwidth"
//  LIMITEXCEEDED_SPOTQUOTA = "LimitExceeded.SpotQuota"
//  MISSINGPARAMETER = "MissingParameter"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  RESOURCEINSUFFICIENT_CLOUDDISKSOLDOUT = "ResourceInsufficient.CloudDiskSoldOut"
//  RESOURCEINSUFFICIENT_CLOUDDISKUNAVAILABLE = "ResourceInsufficient.CloudDiskUnavailable"
//  RESOURCEINSUFFICIENT_SPECIFIEDINSTANCETYPE = "ResourceInsufficient.SpecifiedInstanceType"
//  RESOURCENOTFOUND_INVALIDZONEINSTANCETYPE = "ResourceNotFound.InvalidZoneInstanceType"
//  RESOURCEUNAVAILABLE_INSTANCETYPE = "ResourceUnavailable.InstanceType"
//  RESOURCESSOLDOUT_AVAILABLEZONE = "ResourcesSoldOut.AvailableZone"
//  RESOURCESSOLDOUT_SPECIFIEDINSTANCETYPE = "ResourcesSoldOut.SpecifiedInstanceType"
//  UNAUTHORIZEDOPERATION_MFAEXPIRED = "UnauthorizedOperation.MFAExpired"
//  UNAUTHORIZEDOPERATION_MFANOTFOUND = "UnauthorizedOperation.MFANotFound"
//  UNSUPPORTEDOPERATION_DISKSNAPCREATETIMETOOOLD = "UnsupportedOperation.DiskSnapCreateTimeTooOld"
//  UNSUPPORTEDOPERATION_HETEROGENEOUSCHANGEINSTANCEFAMILY = "UnsupportedOperation.HeterogeneousChangeInstanceFamily"
//  UNSUPPORTEDOPERATION_INSTANCEREINSTALLFAILED = "UnsupportedOperation.InstanceReinstallFailed"
//  UNSUPPORTEDOPERATION_INSTANCESTATEBANNING = "UnsupportedOperation.InstanceStateBanning"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITSERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateExitServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATEFREEZING = "UnsupportedOperation.InstanceStateFreezing"
//  UNSUPPORTEDOPERATION_INSTANCESTATEREBOOTING = "UnsupportedOperation.InstanceStateRebooting"
//  UNSUPPORTEDOPERATION_INSTANCESTATERESCUEMODE = "UnsupportedOperation.InstanceStateRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATERUNNING = "UnsupportedOperation.InstanceStateRunning"
//  UNSUPPORTEDOPERATION_INSTANCESTATESHUTDOWN = "UnsupportedOperation.InstanceStateShutdown"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTARTING = "UnsupportedOperation.InstanceStateStarting"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPING = "UnsupportedOperation.InstanceStateStopping"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATING = "UnsupportedOperation.InstanceStateTerminating"
//  UNSUPPORTEDOPERATION_INVALIDINSTANCEWITHSWAPDISK = "UnsupportedOperation.InvalidInstanceWithSwapDisk"
//  UNSUPPORTEDOPERATION_LOCALDATADISKCHANGEINSTANCEFAMILY = "UnsupportedOperation.LocalDataDiskChangeInstanceFamily"
//  UNSUPPORTEDOPERATION_LOCALDISKMIGRATINGTOCLOUDDISK = "UnsupportedOperation.LocalDiskMigratingToCloudDisk"
//  UNSUPPORTEDOPERATION_NOINSTANCETYPESUPPORTSPOT = "UnsupportedOperation.NoInstanceTypeSupportSpot"
//  UNSUPPORTEDOPERATION_NOTSUPPORTUNPAIDORDER = "UnsupportedOperation.NotSupportUnpaidOrder"
//  UNSUPPORTEDOPERATION_ORIGINALINSTANCETYPEINVALID = "UnsupportedOperation.OriginalInstanceTypeInvalid"
//  UNSUPPORTEDOPERATION_REDHATINSTANCEUNSUPPORTED = "UnsupportedOperation.RedHatInstanceUnsupported"
//  UNSUPPORTEDOPERATION_SPECIALINSTANCETYPE = "UnsupportedOperation.SpecialInstanceType"
//  UNSUPPORTEDOPERATION_STOPPEDMODESTOPCHARGING = "UnsupportedOperation.StoppedModeStopCharging"
//  UNSUPPORTEDOPERATION_STOPPEDMODESTOPCHARGINGSAMEFAMILY = "UnsupportedOperation.StoppedModeStopChargingSameFamily"
//  UNSUPPORTEDOPERATION_SYSTEMDISKTYPE = "UnsupportedOperation.SystemDiskType"
//  UNSUPPORTEDOPERATION_UNSUPPORTEDARMCHANGEINSTANCEFAMILY = "UnsupportedOperation.UnsupportedARMChangeInstanceFamily"
//  UNSUPPORTEDOPERATION_UNSUPPORTEDCHANGEINSTANCEFAMILY = "UnsupportedOperation.UnsupportedChangeInstanceFamily"
//  UNSUPPORTEDOPERATION_UNSUPPORTEDCHANGEINSTANCEFAMILYTOARM = "UnsupportedOperation.UnsupportedChangeInstanceFamilyToARM"
//  UNSUPPORTEDOPERATION_UNSUPPORTEDCHANGEINSTANCETOTHISINSTANCEFAMILY = "UnsupportedOperation.UnsupportedChangeInstanceToThisInstanceFamily"
func (c *Client) ResetInstancesType(request *ResetInstancesTypeRequest) (response *ResetInstancesTypeResponse, err error) {
    return c.ResetInstancesTypeWithContext(context.Background(), request)
}

// ResetInstancesType
// This API is used to change the model of an instance.
//
// * You can only use this API to change the models of instances whose [system disk type](https://intl.cloud.tencent.com/document/api/213/9452?from_cn_redirect=1#block_device) is `CLOUD_BASIC`, `CLOUD_PREMIUM`, or `CLOUD_SSD`.
//
// * Currently, you cannot use this API to change the models of [CDH](https://intl.cloud.tencent.com/document/product/416?from_cn_redirect=1) instances.
//
// error code that may be returned:
//  FAILEDOPERATION_INVALIDINSTANCEAPPLICATIONROLE = "FailedOperation.InvalidInstanceApplicationRole"
//  FAILEDOPERATION_INVALIDINSTANCEAPPLICATIONROLEEMR = "FailedOperation.InvalidInstanceApplicationRoleEmr"
//  FAILEDOPERATION_PROMOTIONALPERIORESTRICTION = "FailedOperation.PromotionalPerioRestriction"
//  INTERNALERROR_TRADEUNKNOWNERROR = "InternalError.TradeUnknownError"
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDACCOUNT_INSUFFICIENTBALANCE = "InvalidAccount.InsufficientBalance"
//  INVALIDACCOUNT_UNPAIDORDER = "InvalidAccount.UnpaidOrder"
//  INVALIDHOSTID_MALFORMED = "InvalidHostId.Malformed"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDINSTANCETYPE_MALFORMED = "InvalidInstanceType.Malformed"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_HOSTIDSTATUSNOTSUPPORT = "InvalidParameter.HostIdStatusNotSupport"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_BASICNETWORKINSTANCEFAMILY = "InvalidParameterValue.BasicNetworkInstanceFamily"
//  INVALIDPARAMETERVALUE_GPUINSTANCEFAMILY = "InvalidParameterValue.GPUInstanceFamily"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_INVALIDGPUFAMILYCHANGE = "InvalidParameterValue.InvalidGPUFamilyChange"
//  INVALIDPARAMETERVALUE_INVALIDINSTANCESOURCE = "InvalidParameterValue.InvalidInstanceSource"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_ZONENOTSUPPORTED = "InvalidParameterValue.ZoneNotSupported"
//  INVALIDPERIOD = "InvalidPeriod"
//  INVALIDPERMISSION = "InvalidPermission"
//  INVALIDZONE_MISMATCHREGION = "InvalidZone.MismatchRegion"
//  LIMITEXCEEDED_EIPNUMLIMIT = "LimitExceeded.EipNumLimit"
//  LIMITEXCEEDED_ENILIMITINSTANCETYPE = "LimitExceeded.EniLimitInstanceType"
//  LIMITEXCEEDED_ENINUMLIMIT = "LimitExceeded.EniNumLimit"
//  LIMITEXCEEDED_INSTANCETYPEBANDWIDTH = "LimitExceeded.InstanceTypeBandwidth"
//  LIMITEXCEEDED_SPOTQUOTA = "LimitExceeded.SpotQuota"
//  MISSINGPARAMETER = "MissingParameter"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  RESOURCEINSUFFICIENT_CLOUDDISKSOLDOUT = "ResourceInsufficient.CloudDiskSoldOut"
//  RESOURCEINSUFFICIENT_CLOUDDISKUNAVAILABLE = "ResourceInsufficient.CloudDiskUnavailable"
//  RESOURCEINSUFFICIENT_SPECIFIEDINSTANCETYPE = "ResourceInsufficient.SpecifiedInstanceType"
//  RESOURCENOTFOUND_INVALIDZONEINSTANCETYPE = "ResourceNotFound.InvalidZoneInstanceType"
//  RESOURCEUNAVAILABLE_INSTANCETYPE = "ResourceUnavailable.InstanceType"
//  RESOURCESSOLDOUT_AVAILABLEZONE = "ResourcesSoldOut.AvailableZone"
//  RESOURCESSOLDOUT_SPECIFIEDINSTANCETYPE = "ResourcesSoldOut.SpecifiedInstanceType"
//  UNAUTHORIZEDOPERATION_MFAEXPIRED = "UnauthorizedOperation.MFAExpired"
//  UNAUTHORIZEDOPERATION_MFANOTFOUND = "UnauthorizedOperation.MFANotFound"
//  UNSUPPORTEDOPERATION_DISKSNAPCREATETIMETOOOLD = "UnsupportedOperation.DiskSnapCreateTimeTooOld"
//  UNSUPPORTEDOPERATION_HETEROGENEOUSCHANGEINSTANCEFAMILY = "UnsupportedOperation.HeterogeneousChangeInstanceFamily"
//  UNSUPPORTEDOPERATION_INSTANCEREINSTALLFAILED = "UnsupportedOperation.InstanceReinstallFailed"
//  UNSUPPORTEDOPERATION_INSTANCESTATEBANNING = "UnsupportedOperation.InstanceStateBanning"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITSERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateExitServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATEFREEZING = "UnsupportedOperation.InstanceStateFreezing"
//  UNSUPPORTEDOPERATION_INSTANCESTATEREBOOTING = "UnsupportedOperation.InstanceStateRebooting"
//  UNSUPPORTEDOPERATION_INSTANCESTATERESCUEMODE = "UnsupportedOperation.InstanceStateRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATERUNNING = "UnsupportedOperation.InstanceStateRunning"
//  UNSUPPORTEDOPERATION_INSTANCESTATESHUTDOWN = "UnsupportedOperation.InstanceStateShutdown"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTARTING = "UnsupportedOperation.InstanceStateStarting"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPING = "UnsupportedOperation.InstanceStateStopping"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATING = "UnsupportedOperation.InstanceStateTerminating"
//  UNSUPPORTEDOPERATION_INVALIDINSTANCEWITHSWAPDISK = "UnsupportedOperation.InvalidInstanceWithSwapDisk"
//  UNSUPPORTEDOPERATION_LOCALDATADISKCHANGEINSTANCEFAMILY = "UnsupportedOperation.LocalDataDiskChangeInstanceFamily"
//  UNSUPPORTEDOPERATION_LOCALDISKMIGRATINGTOCLOUDDISK = "UnsupportedOperation.LocalDiskMigratingToCloudDisk"
//  UNSUPPORTEDOPERATION_NOINSTANCETYPESUPPORTSPOT = "UnsupportedOperation.NoInstanceTypeSupportSpot"
//  UNSUPPORTEDOPERATION_NOTSUPPORTUNPAIDORDER = "UnsupportedOperation.NotSupportUnpaidOrder"
//  UNSUPPORTEDOPERATION_ORIGINALINSTANCETYPEINVALID = "UnsupportedOperation.OriginalInstanceTypeInvalid"
//  UNSUPPORTEDOPERATION_REDHATINSTANCEUNSUPPORTED = "UnsupportedOperation.RedHatInstanceUnsupported"
//  UNSUPPORTEDOPERATION_SPECIALINSTANCETYPE = "UnsupportedOperation.SpecialInstanceType"
//  UNSUPPORTEDOPERATION_STOPPEDMODESTOPCHARGING = "UnsupportedOperation.StoppedModeStopCharging"
//  UNSUPPORTEDOPERATION_STOPPEDMODESTOPCHARGINGSAMEFAMILY = "UnsupportedOperation.StoppedModeStopChargingSameFamily"
//  UNSUPPORTEDOPERATION_SYSTEMDISKTYPE = "UnsupportedOperation.SystemDiskType"
//  UNSUPPORTEDOPERATION_UNSUPPORTEDARMCHANGEINSTANCEFAMILY = "UnsupportedOperation.UnsupportedARMChangeInstanceFamily"
//  UNSUPPORTEDOPERATION_UNSUPPORTEDCHANGEINSTANCEFAMILY = "UnsupportedOperation.UnsupportedChangeInstanceFamily"
//  UNSUPPORTEDOPERATION_UNSUPPORTEDCHANGEINSTANCEFAMILYTOARM = "UnsupportedOperation.UnsupportedChangeInstanceFamilyToARM"
//  UNSUPPORTEDOPERATION_UNSUPPORTEDCHANGEINSTANCETOTHISINSTANCEFAMILY = "UnsupportedOperation.UnsupportedChangeInstanceToThisInstanceFamily"
func (c *Client) ResetInstancesTypeWithContext(ctx context.Context, request *ResetInstancesTypeRequest) (response *ResetInstancesTypeResponse, err error) {
    if request == nil {
        request = NewResetInstancesTypeRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "ResetInstancesType")
    
    if c.GetCredential() == nil {
        return nil, errors.New("ResetInstancesType require credential")
    }

    request.SetContext(ctx)
    
    response = NewResetInstancesTypeResponse()
    err = c.Send(request, response)
    return
}

func NewResizeInstanceDisksRequest() (request *ResizeInstanceDisksRequest) {
    request = &ResizeInstanceDisksRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "ResizeInstanceDisks")
    
    
    return
}

func NewResizeInstanceDisksResponse() (response *ResizeInstanceDisksResponse) {
    response = &ResizeInstanceDisksResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// ResizeInstanceDisks
// This API (ResizeInstanceDisks) is used to expand the data disks of an instance.
//
// 
//
// * Currently, you can only use the API to expand non-elastic data disks whose [disk type](https://intl.cloud.tencent.com/document/api/213/9452?from_cn_redirect=1#block_device) is `CLOUD_BASIC`, `CLOUD_PREMIUM`, or `CLOUD_SSD`. You can use [`DescribeDisks`](https://intl.cloud.tencent.com/document/api/362/16315?from_cn_redirect=1) to check whether a disk is elastic. If the `Portable` field in the response is `false`, it means that the disk is non-elastic.
//
// * Currently, this API does not support [CDH](https://intl.cloud.tencent.com/document/product/416?from_cn_redirect=1) instances.
//
// * Currently, only one data disk can be expanded at a time.
//
// error code that may be returned:
//  FAILEDOPERATION_INVALIDINSTANCEAPPLICATIONROLE = "FailedOperation.InvalidInstanceApplicationRole"
//  INTERNALERROR = "InternalError"
//  INTERNALERROR_TRADEUNKNOWNERROR = "InternalError.TradeUnknownError"
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDACCOUNT_INSUFFICIENTBALANCE = "InvalidAccount.InsufficientBalance"
//  INVALIDACCOUNT_UNPAIDORDER = "InvalidAccount.UnpaidOrder"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_CDHONLYLOCALDATADISKRESIZE = "InvalidParameterValue.CdhOnlyLocalDataDiskResize"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_ZONENOTSUPPORTED = "InvalidParameterValue.ZoneNotSupported"
//  MISSINGPARAMETER = "MissingParameter"
//  MISSINGPARAMETER_ATLEASTONE = "MissingParameter.AtLeastOne"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  UNSUPPORTEDOPERATION_INSTANCECHARGETYPE = "UnsupportedOperation.InstanceChargeType"
//  UNSUPPORTEDOPERATION_INSTANCESTATEBANNING = "UnsupportedOperation.InstanceStateBanning"
//  UNSUPPORTEDOPERATION_INSTANCESTATECORRUPTED = "UnsupportedOperation.InstanceStateCorrupted"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERSERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateEnterServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITRESCUEMODE = "UnsupportedOperation.InstanceStateExitRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEFREEZING = "UnsupportedOperation.InstanceStateFreezing"
//  UNSUPPORTEDOPERATION_INSTANCESTATEPENDING = "UnsupportedOperation.InstanceStatePending"
//  UNSUPPORTEDOPERATION_INSTANCESTATEREBOOTING = "UnsupportedOperation.InstanceStateRebooting"
//  UNSUPPORTEDOPERATION_INSTANCESTATERESCUEMODE = "UnsupportedOperation.InstanceStateRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATERUNNING = "UnsupportedOperation.InstanceStateRunning"
//  UNSUPPORTEDOPERATION_INSTANCESTATESERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATESHUTDOWN = "UnsupportedOperation.InstanceStateShutdown"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTARTING = "UnsupportedOperation.InstanceStateStarting"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPED = "UnsupportedOperation.InstanceStateStopped"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPING = "UnsupportedOperation.InstanceStateStopping"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATING = "UnsupportedOperation.InstanceStateTerminating"
//  UNSUPPORTEDOPERATION_INVALIDDATADISK = "UnsupportedOperation.InvalidDataDisk"
//  UNSUPPORTEDOPERATION_SPECIALINSTANCETYPE = "UnsupportedOperation.SpecialInstanceType"
//  UNSUPPORTEDOPERATION_STOPPEDMODESTOPCHARGING = "UnsupportedOperation.StoppedModeStopCharging"
func (c *Client) ResizeInstanceDisks(request *ResizeInstanceDisksRequest) (response *ResizeInstanceDisksResponse, err error) {
    return c.ResizeInstanceDisksWithContext(context.Background(), request)
}

// ResizeInstanceDisks
// This API (ResizeInstanceDisks) is used to expand the data disks of an instance.
//
// 
//
// * Currently, you can only use the API to expand non-elastic data disks whose [disk type](https://intl.cloud.tencent.com/document/api/213/9452?from_cn_redirect=1#block_device) is `CLOUD_BASIC`, `CLOUD_PREMIUM`, or `CLOUD_SSD`. You can use [`DescribeDisks`](https://intl.cloud.tencent.com/document/api/362/16315?from_cn_redirect=1) to check whether a disk is elastic. If the `Portable` field in the response is `false`, it means that the disk is non-elastic.
//
// * Currently, this API does not support [CDH](https://intl.cloud.tencent.com/document/product/416?from_cn_redirect=1) instances.
//
// * Currently, only one data disk can be expanded at a time.
//
// error code that may be returned:
//  FAILEDOPERATION_INVALIDINSTANCEAPPLICATIONROLE = "FailedOperation.InvalidInstanceApplicationRole"
//  INTERNALERROR = "InternalError"
//  INTERNALERROR_TRADEUNKNOWNERROR = "InternalError.TradeUnknownError"
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDACCOUNT_INSUFFICIENTBALANCE = "InvalidAccount.InsufficientBalance"
//  INVALIDACCOUNT_UNPAIDORDER = "InvalidAccount.UnpaidOrder"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_CDHONLYLOCALDATADISKRESIZE = "InvalidParameterValue.CdhOnlyLocalDataDiskResize"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_ZONENOTSUPPORTED = "InvalidParameterValue.ZoneNotSupported"
//  MISSINGPARAMETER = "MissingParameter"
//  MISSINGPARAMETER_ATLEASTONE = "MissingParameter.AtLeastOne"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  UNSUPPORTEDOPERATION_INSTANCECHARGETYPE = "UnsupportedOperation.InstanceChargeType"
//  UNSUPPORTEDOPERATION_INSTANCESTATEBANNING = "UnsupportedOperation.InstanceStateBanning"
//  UNSUPPORTEDOPERATION_INSTANCESTATECORRUPTED = "UnsupportedOperation.InstanceStateCorrupted"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERSERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateEnterServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITRESCUEMODE = "UnsupportedOperation.InstanceStateExitRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEFREEZING = "UnsupportedOperation.InstanceStateFreezing"
//  UNSUPPORTEDOPERATION_INSTANCESTATEPENDING = "UnsupportedOperation.InstanceStatePending"
//  UNSUPPORTEDOPERATION_INSTANCESTATEREBOOTING = "UnsupportedOperation.InstanceStateRebooting"
//  UNSUPPORTEDOPERATION_INSTANCESTATERESCUEMODE = "UnsupportedOperation.InstanceStateRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATERUNNING = "UnsupportedOperation.InstanceStateRunning"
//  UNSUPPORTEDOPERATION_INSTANCESTATESERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATESHUTDOWN = "UnsupportedOperation.InstanceStateShutdown"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTARTING = "UnsupportedOperation.InstanceStateStarting"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPED = "UnsupportedOperation.InstanceStateStopped"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPING = "UnsupportedOperation.InstanceStateStopping"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATING = "UnsupportedOperation.InstanceStateTerminating"
//  UNSUPPORTEDOPERATION_INVALIDDATADISK = "UnsupportedOperation.InvalidDataDisk"
//  UNSUPPORTEDOPERATION_SPECIALINSTANCETYPE = "UnsupportedOperation.SpecialInstanceType"
//  UNSUPPORTEDOPERATION_STOPPEDMODESTOPCHARGING = "UnsupportedOperation.StoppedModeStopCharging"
func (c *Client) ResizeInstanceDisksWithContext(ctx context.Context, request *ResizeInstanceDisksRequest) (response *ResizeInstanceDisksResponse, err error) {
    if request == nil {
        request = NewResizeInstanceDisksRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "ResizeInstanceDisks")
    
    if c.GetCredential() == nil {
        return nil, errors.New("ResizeInstanceDisks require credential")
    }

    request.SetContext(ctx)
    
    response = NewResizeInstanceDisksResponse()
    err = c.Send(request, response)
    return
}

func NewRunInstancesRequest() (request *RunInstancesRequest) {
    request = &RunInstancesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "RunInstances")
    
    
    return
}

func NewRunInstancesResponse() (response *RunInstancesResponse) {
    response = &RunInstancesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// RunInstances
// This API is used to create one or more instances with a specified configuration.
//
// 
//
// * After an instance is created successfully, it will start up automatically, and the [instance status](https://intl.cloud.tencent.com/document/api/213/9452?from_cn_redirect=1#instance_state) will become "Running".
//
// * If you create a pay-as-you-go instance billed on an hourly basis, an amount equivalent to the hourly rate will be frozen. Make sure your account balance is sufficient before calling this API.
//
// * The number of instances you can purchase through this API is subject to the [Quota for CVM Instances](https://intl.cloud.tencent.com/document/product/213/2664?from_cn_redirect=1). Instances created through this API and in the CVM console are counted toward the quota.
//
// * This API is an async API. An instance ID list is returned after the creation request is sent. However, it does not mean the creation has been completed. The status of the instance will be `Creating` during the creation. You can use [DescribeInstances](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1) to query the status of the instance. If the status changes from `Creating` to `Running`, it means that the instance has been created successfully.
//
// error code that may be returned:
//  ACCOUNTQUALIFICATIONRESTRICTIONS = "AccountQualificationRestrictions"
//  AUTHFAILURE_CAMROLENAMEAUTHENTICATEFAILED = "AuthFailure.CamRoleNameAuthenticateFailed"
//  FAILEDOPERATION_DISASTERRECOVERGROUPNOTFOUND = "FailedOperation.DisasterRecoverGroupNotFound"
//  FAILEDOPERATION_ILLEGALTAGKEY = "FailedOperation.IllegalTagKey"
//  FAILEDOPERATION_ILLEGALTAGVALUE = "FailedOperation.IllegalTagValue"
//  FAILEDOPERATION_INQUIRYPRICEFAILED = "FailedOperation.InquiryPriceFailed"
//  FAILEDOPERATION_NOAVAILABLEIPADDRESSCOUNTINSUBNET = "FailedOperation.NoAvailableIpAddressCountInSubnet"
//  FAILEDOPERATION_PROMOTIONALREGIONRESTRICTION = "FailedOperation.PromotionalRegionRestriction"
//  FAILEDOPERATION_SECURITYGROUPACTIONFAILED = "FailedOperation.SecurityGroupActionFailed"
//  FAILEDOPERATION_SNAPSHOTSIZELARGERTHANDATASIZE = "FailedOperation.SnapshotSizeLargerThanDataSize"
//  FAILEDOPERATION_SNAPSHOTSIZELESSTHANDATASIZE = "FailedOperation.SnapshotSizeLessThanDataSize"
//  FAILEDOPERATION_TAGKEYRESERVED = "FailedOperation.TagKeyReserved"
//  INSTANCESQUOTALIMITEXCEEDED = "InstancesQuotaLimitExceeded"
//  INTERNALERROR_TRADEUNKNOWNERROR = "InternalError.TradeUnknownError"
//  INVALIDACCOUNT_INSUFFICIENTBALANCE = "InvalidAccount.InsufficientBalance"
//  INVALIDCLIENTTOKEN_TOOLONG = "InvalidClientToken.TooLong"
//  INVALIDHOSTID_MALFORMED = "InvalidHostId.Malformed"
//  INVALIDHOSTID_NOTFOUND = "InvalidHostId.NotFound"
//  INVALIDIMAGEID_MALFORMED = "InvalidImageId.Malformed"
//  INVALIDIMAGEID_NOTFOUND = "InvalidImageId.NotFound"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCENAME_TOOLONG = "InvalidInstanceName.TooLong"
//  INVALIDINSTANCETYPE_MALFORMED = "InvalidInstanceType.Malformed"
//  INVALIDKEYPAIRID_MALFORMED = "InvalidKeyPairId.Malformed"
//  INVALIDPARAMETER_CDCNOTSUPPORTED = "InvalidParameter.CdcNotSupported"
//  INVALIDPARAMETER_EDGEZONEMISSINTERNETACCESSIBLE = "InvalidParameter.EdgeZoneMissInternetAccessible"
//  INVALIDPARAMETER_HOSTIDCUSTOMIZEDINSTANCETYPENOTSUPPORT = "InvalidParameter.HostIdCustomizedInstanceTypeNotSupport"
//  INVALIDPARAMETER_HOSTIDINSTANCETYPENOTSUPPORT = "InvalidParameter.HostIdInstanceTypeNotSupport"
//  INVALIDPARAMETER_HOSTIDSTANDARDINSTANCETYPENOTSUPPORT = "InvalidParameter.HostIdStandardInstanceTypeNotSupport"
//  INVALIDPARAMETER_HOSTIDSTATUSNOTSUPPORT = "InvalidParameter.HostIdStatusNotSupport"
//  INVALIDPARAMETER_INSTANCEIMAGENOTSUPPORT = "InvalidParameter.InstanceImageNotSupport"
//  INVALIDPARAMETER_INSTANCETYPESUPPORTEDHOSTNOTFOUND = "InvalidParameter.InstanceTypeSupportedHostNotFound"
//  INVALIDPARAMETER_INTERNETACCESSIBLENOTSUPPORTED = "InvalidParameter.InternetAccessibleNotSupported"
//  INVALIDPARAMETER_INVALIDIPFORMAT = "InvalidParameter.InvalidIpFormat"
//  INVALIDPARAMETER_LACKCORECOUNTORTHREADPERCORE = "InvalidParameter.LackCoreCountOrThreadPerCore"
//  INVALIDPARAMETER_ONLYSUPPORTFOREDGEZONE = "InvalidParameter.OnlySupportForEdgeZone"
//  INVALIDPARAMETER_PARAMETERCONFLICT = "InvalidParameter.ParameterConflict"
//  INVALIDPARAMETER_SNAPSHOTNOTFOUND = "InvalidParameter.SnapshotNotFound"
//  INVALIDPARAMETER_SPECIALPARAMETERFORSPECIALACCOUNT = "InvalidParameter.SpecialParameterForSpecialAccount"
//  INVALIDPARAMETERCOMBINATION = "InvalidParameterCombination"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_BANDWIDTHPACKAGEIDMALFORMED = "InvalidParameterValue.BandwidthPackageIdMalformed"
//  INVALIDPARAMETERVALUE_BANDWIDTHPACKAGEIDNOTFOUND = "InvalidParameterValue.BandwidthPackageIdNotFound"
//  INVALIDPARAMETERVALUE_BANDWIDTHPACKAGEISPNOTMATCH = "InvalidParameterValue.BandwidthPackageIspNotMatch"
//  INVALIDPARAMETERVALUE_BANDWIDTHPACKAGEZONENOTMATCH = "InvalidParameterValue.BandwidthPackageZoneNotMatch"
//  INVALIDPARAMETERVALUE_CAMROLENAMEMALFORMED = "InvalidParameterValue.CamRoleNameMalformed"
//  INVALIDPARAMETERVALUE_CHCHOSTSNOTFOUND = "InvalidParameterValue.ChcHostsNotFound"
//  INVALIDPARAMETERVALUE_CLOUDSSDDATADISKSIZETOOSMALL = "InvalidParameterValue.CloudSsdDataDiskSizeTooSmall"
//  INVALIDPARAMETERVALUE_CORECOUNTVALUE = "InvalidParameterValue.CoreCountValue"
//  INVALIDPARAMETERVALUE_DEDICATEDCLUSTERNOTSUPPORTEDCHARGETYPE = "InvalidParameterValue.DedicatedClusterNotSupportedChargeType"
//  INVALIDPARAMETERVALUE_DUPLICATE = "InvalidParameterValue.Duplicate"
//  INVALIDPARAMETERVALUE_DUPLICATETAGS = "InvalidParameterValue.DuplicateTags"
//  INVALIDPARAMETERVALUE_ELASTICNETWORKNOTEXIST = "InvalidParameterValue.ElasticNetworkNotExist"
//  INVALIDPARAMETERVALUE_ELASTICNETWORKVPCSUBNETMISMATCH = "InvalidParameterValue.ElasticNetworkVpcSubnetMismatch"
//  INVALIDPARAMETERVALUE_EXTERNALIPQUOTALIMITED = "InvalidParameterValue.ExternalIpQuotaLimited"
//  INVALIDPARAMETERVALUE_HPCCLUSTERIDZONEIDNOTMATCH = "InvalidParameterValue.HpcClusterIdZoneIdNotMatch"
//  INVALIDPARAMETERVALUE_IPADDRESSMALFORMED = "InvalidParameterValue.IPAddressMalformed"
//  INVALIDPARAMETERVALUE_ILLEGALHOSTNAME = "InvalidParameterValue.IllegalHostName"
//  INVALIDPARAMETERVALUE_INCORRECTFORMAT = "InvalidParameterValue.IncorrectFormat"
//  INVALIDPARAMETERVALUE_INSTANCETYPENOTFOUND = "InvalidParameterValue.InstanceTypeNotFound"
//  INVALIDPARAMETERVALUE_INSTANCETYPENOTSUPPORTELASTICNETWORKS = "InvalidParameterValue.InstanceTypeNotSupportElasticNetworks"
//  INVALIDPARAMETERVALUE_INSTANCETYPENOTSUPPORTHPCCLUSTER = "InvalidParameterValue.InstanceTypeNotSupportHpcCluster"
//  INVALIDPARAMETERVALUE_INSTANCETYPEREQUIREDHPCCLUSTER = "InvalidParameterValue.InstanceTypeRequiredHpcCluster"
//  INVALIDPARAMETERVALUE_INSUFFICIENTOFFERING = "InvalidParameterValue.InsufficientOffering"
//  INVALIDPARAMETERVALUE_INSUFFICIENTPRICE = "InvalidParameterValue.InsufficientPrice"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEFORGIVENINSTANCETYPE = "InvalidParameterValue.InvalidImageForGivenInstanceType"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEFORMAT = "InvalidParameterValue.InvalidImageFormat"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEID = "InvalidParameterValue.InvalidImageId"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEOSNAME = "InvalidParameterValue.InvalidImageOsName"
//  INVALIDPARAMETERVALUE_INVALIDIMAGESTATE = "InvalidParameterValue.InvalidImageState"
//  INVALIDPARAMETERVALUE_INVALIDIPFORMAT = "InvalidParameterValue.InvalidIpFormat"
//  INVALIDPARAMETERVALUE_INVALIDNETWORKINTERFACEID = "InvalidParameterValue.InvalidNetworkInterfaceId"
//  INVALIDPARAMETERVALUE_INVALIDPARAMETERMINCOUNT = "InvalidParameterValue.InvalidParameterMinCount"
//  INVALIDPARAMETERVALUE_INVALIDPASSWORD = "InvalidParameterValue.InvalidPassword"
//  INVALIDPARAMETERVALUE_INVALIDTIMEFORMAT = "InvalidParameterValue.InvalidTimeFormat"
//  INVALIDPARAMETERVALUE_INVALIDUSERDATAFORMAT = "InvalidParameterValue.InvalidUserDataFormat"
//  INVALIDPARAMETERVALUE_INVALIDVPCIDSUBNETIDNOTFOUND = "InvalidParameterValue.InvalidVpcIdSubnetIdNotFound"
//  INVALIDPARAMETERVALUE_ISPNOTSUPPORTFOREDGEZONE = "InvalidParameterValue.IspNotSupportForEdgeZone"
//  INVALIDPARAMETERVALUE_ISPVALUEREPEATED = "InvalidParameterValue.IspValueRepeated"
//  INVALIDPARAMETERVALUE_KEYPAIRNOTFOUND = "InvalidParameterValue.KeyPairNotFound"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDMALFORMED = "InvalidParameterValue.LaunchTemplateIdMalformed"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDNOTEXISTED = "InvalidParameterValue.LaunchTemplateIdNotExisted"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATENOTFOUND = "InvalidParameterValue.LaunchTemplateNotFound"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEVERSION = "InvalidParameterValue.LaunchTemplateVersion"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_MUSTDHCPENABLEDVPC = "InvalidParameterValue.MustDhcpEnabledVpc"
//  INVALIDPARAMETERVALUE_MUSTENABLEDISRDMA = "InvalidParameterValue.MustEnabledIsRdma"
//  INVALIDPARAMETERVALUE_NOTCDCSUBNET = "InvalidParameterValue.NotCdcSubnet"
//  INVALIDPARAMETERVALUE_RANGE = "InvalidParameterValue.Range"
//  INVALIDPARAMETERVALUE_SNAPSHOTIDMALFORMED = "InvalidParameterValue.SnapshotIdMalformed"
//  INVALIDPARAMETERVALUE_SUBNETIDMALFORMED = "InvalidParameterValue.SubnetIdMalformed"
//  INVALIDPARAMETERVALUE_SUBNETIDZONEIDNOTMATCH = "InvalidParameterValue.SubnetIdZoneIdNotMatch"
//  INVALIDPARAMETERVALUE_SUBNETNOTEXIST = "InvalidParameterValue.SubnetNotExist"
//  INVALIDPARAMETERVALUE_TAGKEYNOTFOUND = "InvalidParameterValue.TagKeyNotFound"
//  INVALIDPARAMETERVALUE_TAGQUOTALIMITEXCEEDED = "InvalidParameterValue.TagQuotaLimitExceeded"
//  INVALIDPARAMETERVALUE_THREADPERCOREVALUE = "InvalidParameterValue.ThreadPerCoreValue"
//  INVALIDPARAMETERVALUE_TOOLARGE = "InvalidParameterValue.TooLarge"
//  INVALIDPARAMETERVALUE_VPCIDMALFORMED = "InvalidParameterValue.VpcIdMalformed"
//  INVALIDPARAMETERVALUE_VPCIDNOTEXIST = "InvalidParameterValue.VpcIdNotExist"
//  INVALIDPARAMETERVALUE_VPCIDSUBNETIDNOTMATCH = "InvalidParameterValue.VpcIdSubnetIdNotMatch"
//  INVALIDPARAMETERVALUE_VPCIDZONEIDNOTMATCH = "InvalidParameterValue.VpcIdZoneIdNotMatch"
//  INVALIDPARAMETERVALUE_VPCNOTSUPPORTIPV6ADDRESS = "InvalidParameterValue.VpcNotSupportIpv6Address"
//  INVALIDPARAMETERVALUE_ZONENOTSUPPORTED = "InvalidParameterValue.ZoneNotSupported"
//  INVALIDPASSWORD = "InvalidPassword"
//  INVALIDPERIOD = "InvalidPeriod"
//  INVALIDPERMISSION = "InvalidPermission"
//  INVALIDPROJECTID_NOTFOUND = "InvalidProjectId.NotFound"
//  INVALIDSECURITYGROUPID_NOTFOUND = "InvalidSecurityGroupId.NotFound"
//  INVALIDZONE_MISMATCHREGION = "InvalidZone.MismatchRegion"
//  LIMITEXCEEDED_CVMINSTANCEQUOTA = "LimitExceeded.CvmInstanceQuota"
//  LIMITEXCEEDED_CVMSVIFSPERSECGROUPLIMITEXCEEDED = "LimitExceeded.CvmsVifsPerSecGroupLimitExceeded"
//  LIMITEXCEEDED_DISASTERRECOVERGROUP = "LimitExceeded.DisasterRecoverGroup"
//  LIMITEXCEEDED_IPV6ADDRESSNUM = "LimitExceeded.IPv6AddressNum"
//  LIMITEXCEEDED_INSTANCEENINUMLIMIT = "LimitExceeded.InstanceEniNumLimit"
//  LIMITEXCEEDED_INSTANCEQUOTA = "LimitExceeded.InstanceQuota"
//  LIMITEXCEEDED_PREPAYQUOTA = "LimitExceeded.PrepayQuota"
//  LIMITEXCEEDED_PREPAYUNDERWRITEQUOTA = "LimitExceeded.PrepayUnderwriteQuota"
//  LIMITEXCEEDED_SINGLEUSGQUOTA = "LimitExceeded.SingleUSGQuota"
//  LIMITEXCEEDED_SPOTQUOTA = "LimitExceeded.SpotQuota"
//  LIMITEXCEEDED_USERSPOTQUOTA = "LimitExceeded.UserSpotQuota"
//  LIMITEXCEEDED_VPCSUBNETNUM = "LimitExceeded.VpcSubnetNum"
//  MISSINGPARAMETER = "MissingParameter"
//  MISSINGPARAMETER_DPDKINSTANCETYPEREQUIREDVPC = "MissingParameter.DPDKInstanceTypeRequiredVPC"
//  MISSINGPARAMETER_MONITORSERVICE = "MissingParameter.MonitorService"
//  OPERATIONDENIED_ACCOUNTNOTSUPPORTED = "OperationDenied.AccountNotSupported"
//  OPERATIONDENIED_CHCINSTALLCLOUDIMAGEWITHOUTDEPLOYNETWORK = "OperationDenied.ChcInstallCloudImageWithoutDeployNetwork"
//  OPERATIONDENIED_NOTBANDWIDTHSHIFTUPUSERAPPLYEDGEZONEEIP = "OperationDenied.NotBandwidthShiftUpUserApplyEdgeZoneEip"
//  RESOURCEINSUFFICIENT_AVAILABILITYZONESOLDOUT = "ResourceInsufficient.AvailabilityZoneSoldOut"
//  RESOURCEINSUFFICIENT_CIDRBLOCK = "ResourceInsufficient.CidrBlock"
//  RESOURCEINSUFFICIENT_CLOUDDISKSOLDOUT = "ResourceInsufficient.CloudDiskSoldOut"
//  RESOURCEINSUFFICIENT_CLOUDDISKUNAVAILABLE = "ResourceInsufficient.CloudDiskUnavailable"
//  RESOURCEINSUFFICIENT_DISASTERRECOVERGROUPCVMQUOTA = "ResourceInsufficient.DisasterRecoverGroupCvmQuota"
//  RESOURCEINSUFFICIENT_INSUFFICIENTOFFERINGMINIMUM = "ResourceInsufficient.InsufficientOfferingMinimum"
//  RESOURCEINSUFFICIENT_SPECIFIEDINSTANCETYPE = "ResourceInsufficient.SpecifiedInstanceType"
//  RESOURCENOTFOUND_HPCCLUSTER = "ResourceNotFound.HpcCluster"
//  RESOURCENOTFOUND_KEYPAIRNOTFOUND = "ResourceNotFound.KeyPairNotFound"
//  RESOURCENOTFOUND_NODEFAULTCBS = "ResourceNotFound.NoDefaultCbs"
//  RESOURCENOTFOUND_NODEFAULTCBSWITHREASON = "ResourceNotFound.NoDefaultCbsWithReason"
//  RESOURCEUNAVAILABLE_INSTANCETYPE = "ResourceUnavailable.InstanceType"
//  RESOURCESSOLDOUT_EIPINSUFFICIENT = "ResourcesSoldOut.EipInsufficient"
//  RESOURCESSOLDOUT_SPECIFIEDINSTANCETYPE = "ResourcesSoldOut.SpecifiedInstanceType"
//  UNAUTHORIZEDOPERATION_INVALIDTOKEN = "UnauthorizedOperation.InvalidToken"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
//  UNSUPPORTEDOPERATION_BANDWIDTHPACKAGEIDNOTSUPPORTED = "UnsupportedOperation.BandwidthPackageIdNotSupported"
//  UNSUPPORTEDOPERATION_HIBERNATIONOSVERSION = "UnsupportedOperation.HibernationOsVersion"
//  UNSUPPORTEDOPERATION_INSTANCESTATEISOLATING = "UnsupportedOperation.InstanceStateIsolating"
//  UNSUPPORTEDOPERATION_INSTANCETYPENOTSUPPORTCONFIDENTIALITY = "UnsupportedOperation.InstanceTypeNotSupportConfidentiality"
//  UNSUPPORTEDOPERATION_INSTANCETYPENOTSUPPORTGRIDLICENCE = "UnsupportedOperation.InstanceTypeNotSupportGridLicence"
//  UNSUPPORTEDOPERATION_INSTANCETYPENOTSUPPORTJUMBOFRAME = "UnsupportedOperation.InstanceTypeNotSupportJumboFrame"
//  UNSUPPORTEDOPERATION_INSTANCESENABLEJUMBOWITHOUTREBOOT = "UnsupportedOperation.InstancesEnableJumboWithoutReboot"
//  UNSUPPORTEDOPERATION_INVALIDDISK = "UnsupportedOperation.InvalidDisk"
//  UNSUPPORTEDOPERATION_INVALIDREGIONDISKENCRYPT = "UnsupportedOperation.InvalidRegionDiskEncrypt"
//  UNSUPPORTEDOPERATION_KEYPAIRUNSUPPORTEDWINDOWS = "UnsupportedOperation.KeyPairUnsupportedWindows"
//  UNSUPPORTEDOPERATION_MINCOUNTUNSUPPORTEDCHARGETYPE = "UnsupportedOperation.MinCountUnsupportedChargeType"
//  UNSUPPORTEDOPERATION_MINCOUNTUNSUPPORTEDREGION = "UnsupportedOperation.MinCountUnsupportedRegion"
//  UNSUPPORTEDOPERATION_NOINSTANCETYPESUPPORTSPOT = "UnsupportedOperation.NoInstanceTypeSupportSpot"
//  UNSUPPORTEDOPERATION_NOTSUPPORTIMPORTINSTANCESACTIONTIMER = "UnsupportedOperation.NotSupportImportInstancesActionTimer"
//  UNSUPPORTEDOPERATION_ONLYFORPREPAIDACCOUNT = "UnsupportedOperation.OnlyForPrepaidAccount"
//  UNSUPPORTEDOPERATION_PERIODICCONTRACTNOTSUPPORTMANUALRENEW = "UnsupportedOperation.PeriodicContractNotSupportManualRenew"
//  UNSUPPORTEDOPERATION_RAWLOCALDISKINSREINSTALLTOQCOW2 = "UnsupportedOperation.RawLocalDiskInsReinstalltoQcow2"
//  UNSUPPORTEDOPERATION_SPOTUNSUPPORTEDREGION = "UnsupportedOperation.SpotUnsupportedRegion"
//  UNSUPPORTEDOPERATION_UNDERWRITINGINSTANCETYPEONLYSUPPORTAUTORENEW = "UnsupportedOperation.UnderwritingInstanceTypeOnlySupportAutoRenew"
//  UNSUPPORTEDOPERATION_UNSUPPORTEDINTERNATIONALUSER = "UnsupportedOperation.UnsupportedInternationalUser"
//  VPCADDRNOTINSUBNET = "VpcAddrNotInSubNet"
//  VPCIPISUSED = "VpcIpIsUsed"
func (c *Client) RunInstances(request *RunInstancesRequest) (response *RunInstancesResponse, err error) {
    return c.RunInstancesWithContext(context.Background(), request)
}

// RunInstances
// This API is used to create one or more instances with a specified configuration.
//
// 
//
// * After an instance is created successfully, it will start up automatically, and the [instance status](https://intl.cloud.tencent.com/document/api/213/9452?from_cn_redirect=1#instance_state) will become "Running".
//
// * If you create a pay-as-you-go instance billed on an hourly basis, an amount equivalent to the hourly rate will be frozen. Make sure your account balance is sufficient before calling this API.
//
// * The number of instances you can purchase through this API is subject to the [Quota for CVM Instances](https://intl.cloud.tencent.com/document/product/213/2664?from_cn_redirect=1). Instances created through this API and in the CVM console are counted toward the quota.
//
// * This API is an async API. An instance ID list is returned after the creation request is sent. However, it does not mean the creation has been completed. The status of the instance will be `Creating` during the creation. You can use [DescribeInstances](https://intl.cloud.tencent.com/document/api/213/15728?from_cn_redirect=1) to query the status of the instance. If the status changes from `Creating` to `Running`, it means that the instance has been created successfully.
//
// error code that may be returned:
//  ACCOUNTQUALIFICATIONRESTRICTIONS = "AccountQualificationRestrictions"
//  AUTHFAILURE_CAMROLENAMEAUTHENTICATEFAILED = "AuthFailure.CamRoleNameAuthenticateFailed"
//  FAILEDOPERATION_DISASTERRECOVERGROUPNOTFOUND = "FailedOperation.DisasterRecoverGroupNotFound"
//  FAILEDOPERATION_ILLEGALTAGKEY = "FailedOperation.IllegalTagKey"
//  FAILEDOPERATION_ILLEGALTAGVALUE = "FailedOperation.IllegalTagValue"
//  FAILEDOPERATION_INQUIRYPRICEFAILED = "FailedOperation.InquiryPriceFailed"
//  FAILEDOPERATION_NOAVAILABLEIPADDRESSCOUNTINSUBNET = "FailedOperation.NoAvailableIpAddressCountInSubnet"
//  FAILEDOPERATION_PROMOTIONALREGIONRESTRICTION = "FailedOperation.PromotionalRegionRestriction"
//  FAILEDOPERATION_SECURITYGROUPACTIONFAILED = "FailedOperation.SecurityGroupActionFailed"
//  FAILEDOPERATION_SNAPSHOTSIZELARGERTHANDATASIZE = "FailedOperation.SnapshotSizeLargerThanDataSize"
//  FAILEDOPERATION_SNAPSHOTSIZELESSTHANDATASIZE = "FailedOperation.SnapshotSizeLessThanDataSize"
//  FAILEDOPERATION_TAGKEYRESERVED = "FailedOperation.TagKeyReserved"
//  INSTANCESQUOTALIMITEXCEEDED = "InstancesQuotaLimitExceeded"
//  INTERNALERROR_TRADEUNKNOWNERROR = "InternalError.TradeUnknownError"
//  INVALIDACCOUNT_INSUFFICIENTBALANCE = "InvalidAccount.InsufficientBalance"
//  INVALIDCLIENTTOKEN_TOOLONG = "InvalidClientToken.TooLong"
//  INVALIDHOSTID_MALFORMED = "InvalidHostId.Malformed"
//  INVALIDHOSTID_NOTFOUND = "InvalidHostId.NotFound"
//  INVALIDIMAGEID_MALFORMED = "InvalidImageId.Malformed"
//  INVALIDIMAGEID_NOTFOUND = "InvalidImageId.NotFound"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCENAME_TOOLONG = "InvalidInstanceName.TooLong"
//  INVALIDINSTANCETYPE_MALFORMED = "InvalidInstanceType.Malformed"
//  INVALIDKEYPAIRID_MALFORMED = "InvalidKeyPairId.Malformed"
//  INVALIDPARAMETER_CDCNOTSUPPORTED = "InvalidParameter.CdcNotSupported"
//  INVALIDPARAMETER_EDGEZONEMISSINTERNETACCESSIBLE = "InvalidParameter.EdgeZoneMissInternetAccessible"
//  INVALIDPARAMETER_HOSTIDCUSTOMIZEDINSTANCETYPENOTSUPPORT = "InvalidParameter.HostIdCustomizedInstanceTypeNotSupport"
//  INVALIDPARAMETER_HOSTIDINSTANCETYPENOTSUPPORT = "InvalidParameter.HostIdInstanceTypeNotSupport"
//  INVALIDPARAMETER_HOSTIDSTANDARDINSTANCETYPENOTSUPPORT = "InvalidParameter.HostIdStandardInstanceTypeNotSupport"
//  INVALIDPARAMETER_HOSTIDSTATUSNOTSUPPORT = "InvalidParameter.HostIdStatusNotSupport"
//  INVALIDPARAMETER_INSTANCEIMAGENOTSUPPORT = "InvalidParameter.InstanceImageNotSupport"
//  INVALIDPARAMETER_INSTANCETYPESUPPORTEDHOSTNOTFOUND = "InvalidParameter.InstanceTypeSupportedHostNotFound"
//  INVALIDPARAMETER_INTERNETACCESSIBLENOTSUPPORTED = "InvalidParameter.InternetAccessibleNotSupported"
//  INVALIDPARAMETER_INVALIDIPFORMAT = "InvalidParameter.InvalidIpFormat"
//  INVALIDPARAMETER_LACKCORECOUNTORTHREADPERCORE = "InvalidParameter.LackCoreCountOrThreadPerCore"
//  INVALIDPARAMETER_ONLYSUPPORTFOREDGEZONE = "InvalidParameter.OnlySupportForEdgeZone"
//  INVALIDPARAMETER_PARAMETERCONFLICT = "InvalidParameter.ParameterConflict"
//  INVALIDPARAMETER_SNAPSHOTNOTFOUND = "InvalidParameter.SnapshotNotFound"
//  INVALIDPARAMETER_SPECIALPARAMETERFORSPECIALACCOUNT = "InvalidParameter.SpecialParameterForSpecialAccount"
//  INVALIDPARAMETERCOMBINATION = "InvalidParameterCombination"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_BANDWIDTHPACKAGEIDMALFORMED = "InvalidParameterValue.BandwidthPackageIdMalformed"
//  INVALIDPARAMETERVALUE_BANDWIDTHPACKAGEIDNOTFOUND = "InvalidParameterValue.BandwidthPackageIdNotFound"
//  INVALIDPARAMETERVALUE_BANDWIDTHPACKAGEISPNOTMATCH = "InvalidParameterValue.BandwidthPackageIspNotMatch"
//  INVALIDPARAMETERVALUE_BANDWIDTHPACKAGEZONENOTMATCH = "InvalidParameterValue.BandwidthPackageZoneNotMatch"
//  INVALIDPARAMETERVALUE_CAMROLENAMEMALFORMED = "InvalidParameterValue.CamRoleNameMalformed"
//  INVALIDPARAMETERVALUE_CHCHOSTSNOTFOUND = "InvalidParameterValue.ChcHostsNotFound"
//  INVALIDPARAMETERVALUE_CLOUDSSDDATADISKSIZETOOSMALL = "InvalidParameterValue.CloudSsdDataDiskSizeTooSmall"
//  INVALIDPARAMETERVALUE_CORECOUNTVALUE = "InvalidParameterValue.CoreCountValue"
//  INVALIDPARAMETERVALUE_DEDICATEDCLUSTERNOTSUPPORTEDCHARGETYPE = "InvalidParameterValue.DedicatedClusterNotSupportedChargeType"
//  INVALIDPARAMETERVALUE_DUPLICATE = "InvalidParameterValue.Duplicate"
//  INVALIDPARAMETERVALUE_DUPLICATETAGS = "InvalidParameterValue.DuplicateTags"
//  INVALIDPARAMETERVALUE_ELASTICNETWORKNOTEXIST = "InvalidParameterValue.ElasticNetworkNotExist"
//  INVALIDPARAMETERVALUE_ELASTICNETWORKVPCSUBNETMISMATCH = "InvalidParameterValue.ElasticNetworkVpcSubnetMismatch"
//  INVALIDPARAMETERVALUE_EXTERNALIPQUOTALIMITED = "InvalidParameterValue.ExternalIpQuotaLimited"
//  INVALIDPARAMETERVALUE_HPCCLUSTERIDZONEIDNOTMATCH = "InvalidParameterValue.HpcClusterIdZoneIdNotMatch"
//  INVALIDPARAMETERVALUE_IPADDRESSMALFORMED = "InvalidParameterValue.IPAddressMalformed"
//  INVALIDPARAMETERVALUE_ILLEGALHOSTNAME = "InvalidParameterValue.IllegalHostName"
//  INVALIDPARAMETERVALUE_INCORRECTFORMAT = "InvalidParameterValue.IncorrectFormat"
//  INVALIDPARAMETERVALUE_INSTANCETYPENOTFOUND = "InvalidParameterValue.InstanceTypeNotFound"
//  INVALIDPARAMETERVALUE_INSTANCETYPENOTSUPPORTELASTICNETWORKS = "InvalidParameterValue.InstanceTypeNotSupportElasticNetworks"
//  INVALIDPARAMETERVALUE_INSTANCETYPENOTSUPPORTHPCCLUSTER = "InvalidParameterValue.InstanceTypeNotSupportHpcCluster"
//  INVALIDPARAMETERVALUE_INSTANCETYPEREQUIREDHPCCLUSTER = "InvalidParameterValue.InstanceTypeRequiredHpcCluster"
//  INVALIDPARAMETERVALUE_INSUFFICIENTOFFERING = "InvalidParameterValue.InsufficientOffering"
//  INVALIDPARAMETERVALUE_INSUFFICIENTPRICE = "InvalidParameterValue.InsufficientPrice"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEFORGIVENINSTANCETYPE = "InvalidParameterValue.InvalidImageForGivenInstanceType"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEFORMAT = "InvalidParameterValue.InvalidImageFormat"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEID = "InvalidParameterValue.InvalidImageId"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEOSNAME = "InvalidParameterValue.InvalidImageOsName"
//  INVALIDPARAMETERVALUE_INVALIDIMAGESTATE = "InvalidParameterValue.InvalidImageState"
//  INVALIDPARAMETERVALUE_INVALIDIPFORMAT = "InvalidParameterValue.InvalidIpFormat"
//  INVALIDPARAMETERVALUE_INVALIDNETWORKINTERFACEID = "InvalidParameterValue.InvalidNetworkInterfaceId"
//  INVALIDPARAMETERVALUE_INVALIDPARAMETERMINCOUNT = "InvalidParameterValue.InvalidParameterMinCount"
//  INVALIDPARAMETERVALUE_INVALIDPASSWORD = "InvalidParameterValue.InvalidPassword"
//  INVALIDPARAMETERVALUE_INVALIDTIMEFORMAT = "InvalidParameterValue.InvalidTimeFormat"
//  INVALIDPARAMETERVALUE_INVALIDUSERDATAFORMAT = "InvalidParameterValue.InvalidUserDataFormat"
//  INVALIDPARAMETERVALUE_INVALIDVPCIDSUBNETIDNOTFOUND = "InvalidParameterValue.InvalidVpcIdSubnetIdNotFound"
//  INVALIDPARAMETERVALUE_ISPNOTSUPPORTFOREDGEZONE = "InvalidParameterValue.IspNotSupportForEdgeZone"
//  INVALIDPARAMETERVALUE_ISPVALUEREPEATED = "InvalidParameterValue.IspValueRepeated"
//  INVALIDPARAMETERVALUE_KEYPAIRNOTFOUND = "InvalidParameterValue.KeyPairNotFound"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDMALFORMED = "InvalidParameterValue.LaunchTemplateIdMalformed"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDNOTEXISTED = "InvalidParameterValue.LaunchTemplateIdNotExisted"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATENOTFOUND = "InvalidParameterValue.LaunchTemplateNotFound"
//  INVALIDPARAMETERVALUE_LAUNCHTEMPLATEVERSION = "InvalidParameterValue.LaunchTemplateVersion"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_MUSTDHCPENABLEDVPC = "InvalidParameterValue.MustDhcpEnabledVpc"
//  INVALIDPARAMETERVALUE_MUSTENABLEDISRDMA = "InvalidParameterValue.MustEnabledIsRdma"
//  INVALIDPARAMETERVALUE_NOTCDCSUBNET = "InvalidParameterValue.NotCdcSubnet"
//  INVALIDPARAMETERVALUE_RANGE = "InvalidParameterValue.Range"
//  INVALIDPARAMETERVALUE_SNAPSHOTIDMALFORMED = "InvalidParameterValue.SnapshotIdMalformed"
//  INVALIDPARAMETERVALUE_SUBNETIDMALFORMED = "InvalidParameterValue.SubnetIdMalformed"
//  INVALIDPARAMETERVALUE_SUBNETIDZONEIDNOTMATCH = "InvalidParameterValue.SubnetIdZoneIdNotMatch"
//  INVALIDPARAMETERVALUE_SUBNETNOTEXIST = "InvalidParameterValue.SubnetNotExist"
//  INVALIDPARAMETERVALUE_TAGKEYNOTFOUND = "InvalidParameterValue.TagKeyNotFound"
//  INVALIDPARAMETERVALUE_TAGQUOTALIMITEXCEEDED = "InvalidParameterValue.TagQuotaLimitExceeded"
//  INVALIDPARAMETERVALUE_THREADPERCOREVALUE = "InvalidParameterValue.ThreadPerCoreValue"
//  INVALIDPARAMETERVALUE_TOOLARGE = "InvalidParameterValue.TooLarge"
//  INVALIDPARAMETERVALUE_VPCIDMALFORMED = "InvalidParameterValue.VpcIdMalformed"
//  INVALIDPARAMETERVALUE_VPCIDNOTEXIST = "InvalidParameterValue.VpcIdNotExist"
//  INVALIDPARAMETERVALUE_VPCIDSUBNETIDNOTMATCH = "InvalidParameterValue.VpcIdSubnetIdNotMatch"
//  INVALIDPARAMETERVALUE_VPCIDZONEIDNOTMATCH = "InvalidParameterValue.VpcIdZoneIdNotMatch"
//  INVALIDPARAMETERVALUE_VPCNOTSUPPORTIPV6ADDRESS = "InvalidParameterValue.VpcNotSupportIpv6Address"
//  INVALIDPARAMETERVALUE_ZONENOTSUPPORTED = "InvalidParameterValue.ZoneNotSupported"
//  INVALIDPASSWORD = "InvalidPassword"
//  INVALIDPERIOD = "InvalidPeriod"
//  INVALIDPERMISSION = "InvalidPermission"
//  INVALIDPROJECTID_NOTFOUND = "InvalidProjectId.NotFound"
//  INVALIDSECURITYGROUPID_NOTFOUND = "InvalidSecurityGroupId.NotFound"
//  INVALIDZONE_MISMATCHREGION = "InvalidZone.MismatchRegion"
//  LIMITEXCEEDED_CVMINSTANCEQUOTA = "LimitExceeded.CvmInstanceQuota"
//  LIMITEXCEEDED_CVMSVIFSPERSECGROUPLIMITEXCEEDED = "LimitExceeded.CvmsVifsPerSecGroupLimitExceeded"
//  LIMITEXCEEDED_DISASTERRECOVERGROUP = "LimitExceeded.DisasterRecoverGroup"
//  LIMITEXCEEDED_IPV6ADDRESSNUM = "LimitExceeded.IPv6AddressNum"
//  LIMITEXCEEDED_INSTANCEENINUMLIMIT = "LimitExceeded.InstanceEniNumLimit"
//  LIMITEXCEEDED_INSTANCEQUOTA = "LimitExceeded.InstanceQuota"
//  LIMITEXCEEDED_PREPAYQUOTA = "LimitExceeded.PrepayQuota"
//  LIMITEXCEEDED_PREPAYUNDERWRITEQUOTA = "LimitExceeded.PrepayUnderwriteQuota"
//  LIMITEXCEEDED_SINGLEUSGQUOTA = "LimitExceeded.SingleUSGQuota"
//  LIMITEXCEEDED_SPOTQUOTA = "LimitExceeded.SpotQuota"
//  LIMITEXCEEDED_USERSPOTQUOTA = "LimitExceeded.UserSpotQuota"
//  LIMITEXCEEDED_VPCSUBNETNUM = "LimitExceeded.VpcSubnetNum"
//  MISSINGPARAMETER = "MissingParameter"
//  MISSINGPARAMETER_DPDKINSTANCETYPEREQUIREDVPC = "MissingParameter.DPDKInstanceTypeRequiredVPC"
//  MISSINGPARAMETER_MONITORSERVICE = "MissingParameter.MonitorService"
//  OPERATIONDENIED_ACCOUNTNOTSUPPORTED = "OperationDenied.AccountNotSupported"
//  OPERATIONDENIED_CHCINSTALLCLOUDIMAGEWITHOUTDEPLOYNETWORK = "OperationDenied.ChcInstallCloudImageWithoutDeployNetwork"
//  OPERATIONDENIED_NOTBANDWIDTHSHIFTUPUSERAPPLYEDGEZONEEIP = "OperationDenied.NotBandwidthShiftUpUserApplyEdgeZoneEip"
//  RESOURCEINSUFFICIENT_AVAILABILITYZONESOLDOUT = "ResourceInsufficient.AvailabilityZoneSoldOut"
//  RESOURCEINSUFFICIENT_CIDRBLOCK = "ResourceInsufficient.CidrBlock"
//  RESOURCEINSUFFICIENT_CLOUDDISKSOLDOUT = "ResourceInsufficient.CloudDiskSoldOut"
//  RESOURCEINSUFFICIENT_CLOUDDISKUNAVAILABLE = "ResourceInsufficient.CloudDiskUnavailable"
//  RESOURCEINSUFFICIENT_DISASTERRECOVERGROUPCVMQUOTA = "ResourceInsufficient.DisasterRecoverGroupCvmQuota"
//  RESOURCEINSUFFICIENT_INSUFFICIENTOFFERINGMINIMUM = "ResourceInsufficient.InsufficientOfferingMinimum"
//  RESOURCEINSUFFICIENT_SPECIFIEDINSTANCETYPE = "ResourceInsufficient.SpecifiedInstanceType"
//  RESOURCENOTFOUND_HPCCLUSTER = "ResourceNotFound.HpcCluster"
//  RESOURCENOTFOUND_KEYPAIRNOTFOUND = "ResourceNotFound.KeyPairNotFound"
//  RESOURCENOTFOUND_NODEFAULTCBS = "ResourceNotFound.NoDefaultCbs"
//  RESOURCENOTFOUND_NODEFAULTCBSWITHREASON = "ResourceNotFound.NoDefaultCbsWithReason"
//  RESOURCEUNAVAILABLE_INSTANCETYPE = "ResourceUnavailable.InstanceType"
//  RESOURCESSOLDOUT_EIPINSUFFICIENT = "ResourcesSoldOut.EipInsufficient"
//  RESOURCESSOLDOUT_SPECIFIEDINSTANCETYPE = "ResourcesSoldOut.SpecifiedInstanceType"
//  UNAUTHORIZEDOPERATION_INVALIDTOKEN = "UnauthorizedOperation.InvalidToken"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
//  UNSUPPORTEDOPERATION_BANDWIDTHPACKAGEIDNOTSUPPORTED = "UnsupportedOperation.BandwidthPackageIdNotSupported"
//  UNSUPPORTEDOPERATION_HIBERNATIONOSVERSION = "UnsupportedOperation.HibernationOsVersion"
//  UNSUPPORTEDOPERATION_INSTANCESTATEISOLATING = "UnsupportedOperation.InstanceStateIsolating"
//  UNSUPPORTEDOPERATION_INSTANCETYPENOTSUPPORTCONFIDENTIALITY = "UnsupportedOperation.InstanceTypeNotSupportConfidentiality"
//  UNSUPPORTEDOPERATION_INSTANCETYPENOTSUPPORTGRIDLICENCE = "UnsupportedOperation.InstanceTypeNotSupportGridLicence"
//  UNSUPPORTEDOPERATION_INSTANCETYPENOTSUPPORTJUMBOFRAME = "UnsupportedOperation.InstanceTypeNotSupportJumboFrame"
//  UNSUPPORTEDOPERATION_INSTANCESENABLEJUMBOWITHOUTREBOOT = "UnsupportedOperation.InstancesEnableJumboWithoutReboot"
//  UNSUPPORTEDOPERATION_INVALIDDISK = "UnsupportedOperation.InvalidDisk"
//  UNSUPPORTEDOPERATION_INVALIDREGIONDISKENCRYPT = "UnsupportedOperation.InvalidRegionDiskEncrypt"
//  UNSUPPORTEDOPERATION_KEYPAIRUNSUPPORTEDWINDOWS = "UnsupportedOperation.KeyPairUnsupportedWindows"
//  UNSUPPORTEDOPERATION_MINCOUNTUNSUPPORTEDCHARGETYPE = "UnsupportedOperation.MinCountUnsupportedChargeType"
//  UNSUPPORTEDOPERATION_MINCOUNTUNSUPPORTEDREGION = "UnsupportedOperation.MinCountUnsupportedRegion"
//  UNSUPPORTEDOPERATION_NOINSTANCETYPESUPPORTSPOT = "UnsupportedOperation.NoInstanceTypeSupportSpot"
//  UNSUPPORTEDOPERATION_NOTSUPPORTIMPORTINSTANCESACTIONTIMER = "UnsupportedOperation.NotSupportImportInstancesActionTimer"
//  UNSUPPORTEDOPERATION_ONLYFORPREPAIDACCOUNT = "UnsupportedOperation.OnlyForPrepaidAccount"
//  UNSUPPORTEDOPERATION_PERIODICCONTRACTNOTSUPPORTMANUALRENEW = "UnsupportedOperation.PeriodicContractNotSupportManualRenew"
//  UNSUPPORTEDOPERATION_RAWLOCALDISKINSREINSTALLTOQCOW2 = "UnsupportedOperation.RawLocalDiskInsReinstalltoQcow2"
//  UNSUPPORTEDOPERATION_SPOTUNSUPPORTEDREGION = "UnsupportedOperation.SpotUnsupportedRegion"
//  UNSUPPORTEDOPERATION_UNDERWRITINGINSTANCETYPEONLYSUPPORTAUTORENEW = "UnsupportedOperation.UnderwritingInstanceTypeOnlySupportAutoRenew"
//  UNSUPPORTEDOPERATION_UNSUPPORTEDINTERNATIONALUSER = "UnsupportedOperation.UnsupportedInternationalUser"
//  VPCADDRNOTINSUBNET = "VpcAddrNotInSubNet"
//  VPCIPISUSED = "VpcIpIsUsed"
func (c *Client) RunInstancesWithContext(ctx context.Context, request *RunInstancesRequest) (response *RunInstancesResponse, err error) {
    if request == nil {
        request = NewRunInstancesRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "RunInstances")
    
    if c.GetCredential() == nil {
        return nil, errors.New("RunInstances require credential")
    }

    request.SetContext(ctx)
    
    response = NewRunInstancesResponse()
    err = c.Send(request, response)
    return
}

func NewStartInstancesRequest() (request *StartInstancesRequest) {
    request = &StartInstancesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "StartInstances")
    
    
    return
}

func NewStartInstancesResponse() (response *StartInstancesResponse) {
    response = &StartInstancesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// StartInstances
// This API is used to start instances.
//
// 
//
// * You can only perform this operation on instances whose status is `STOPPED`.
//
// * The instance status will become `STARTING` when the API is called successfully and `RUNNING` when the instance is successfully started.
//
// * Batch operations are supported. The maximum number of instances in each request is 100.
//
// error code that may be returned:
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_ZONENOTSUPPORTED = "InvalidParameterValue.ZoneNotSupported"
//  MISSINGPARAMETER = "MissingParameter"
//  MUTEXOPERATION_TASKRUNNING = "MutexOperation.TaskRunning"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  UNSUPPORTEDOPERATION_INSTANCEREINSTALLFAILED = "UnsupportedOperation.InstanceReinstallFailed"
//  UNSUPPORTEDOPERATION_INSTANCESTATECORRUPTED = "UnsupportedOperation.InstanceStateCorrupted"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERRESCUEMODE = "UnsupportedOperation.InstanceStateEnterRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITRESCUEMODE = "UnsupportedOperation.InstanceStateExitRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEFREEZING = "UnsupportedOperation.InstanceStateFreezing"
//  UNSUPPORTEDOPERATION_INSTANCESTATEISOLATING = "UnsupportedOperation.InstanceStateIsolating"
//  UNSUPPORTEDOPERATION_INSTANCESTATEPENDING = "UnsupportedOperation.InstanceStatePending"
//  UNSUPPORTEDOPERATION_INSTANCESTATEREBOOTING = "UnsupportedOperation.InstanceStateRebooting"
//  UNSUPPORTEDOPERATION_INSTANCESTATERESCUEMODE = "UnsupportedOperation.InstanceStateRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATERUNNING = "UnsupportedOperation.InstanceStateRunning"
//  UNSUPPORTEDOPERATION_INSTANCESTATESERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATESHUTDOWN = "UnsupportedOperation.InstanceStateShutdown"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTARTING = "UnsupportedOperation.InstanceStateStarting"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPED = "UnsupportedOperation.InstanceStateStopped"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPING = "UnsupportedOperation.InstanceStateStopping"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATED = "UnsupportedOperation.InstanceStateTerminated"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATING = "UnsupportedOperation.InstanceStateTerminating"
func (c *Client) StartInstances(request *StartInstancesRequest) (response *StartInstancesResponse, err error) {
    return c.StartInstancesWithContext(context.Background(), request)
}

// StartInstances
// This API is used to start instances.
//
// 
//
// * You can only perform this operation on instances whose status is `STOPPED`.
//
// * The instance status will become `STARTING` when the API is called successfully and `RUNNING` when the instance is successfully started.
//
// * Batch operations are supported. The maximum number of instances in each request is 100.
//
// error code that may be returned:
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_ZONENOTSUPPORTED = "InvalidParameterValue.ZoneNotSupported"
//  MISSINGPARAMETER = "MissingParameter"
//  MUTEXOPERATION_TASKRUNNING = "MutexOperation.TaskRunning"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  UNSUPPORTEDOPERATION_INSTANCEREINSTALLFAILED = "UnsupportedOperation.InstanceReinstallFailed"
//  UNSUPPORTEDOPERATION_INSTANCESTATECORRUPTED = "UnsupportedOperation.InstanceStateCorrupted"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERRESCUEMODE = "UnsupportedOperation.InstanceStateEnterRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITRESCUEMODE = "UnsupportedOperation.InstanceStateExitRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEFREEZING = "UnsupportedOperation.InstanceStateFreezing"
//  UNSUPPORTEDOPERATION_INSTANCESTATEISOLATING = "UnsupportedOperation.InstanceStateIsolating"
//  UNSUPPORTEDOPERATION_INSTANCESTATEPENDING = "UnsupportedOperation.InstanceStatePending"
//  UNSUPPORTEDOPERATION_INSTANCESTATEREBOOTING = "UnsupportedOperation.InstanceStateRebooting"
//  UNSUPPORTEDOPERATION_INSTANCESTATERESCUEMODE = "UnsupportedOperation.InstanceStateRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATERUNNING = "UnsupportedOperation.InstanceStateRunning"
//  UNSUPPORTEDOPERATION_INSTANCESTATESERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATESHUTDOWN = "UnsupportedOperation.InstanceStateShutdown"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTARTING = "UnsupportedOperation.InstanceStateStarting"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPED = "UnsupportedOperation.InstanceStateStopped"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPING = "UnsupportedOperation.InstanceStateStopping"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATED = "UnsupportedOperation.InstanceStateTerminated"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATING = "UnsupportedOperation.InstanceStateTerminating"
func (c *Client) StartInstancesWithContext(ctx context.Context, request *StartInstancesRequest) (response *StartInstancesResponse, err error) {
    if request == nil {
        request = NewStartInstancesRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "StartInstances")
    
    if c.GetCredential() == nil {
        return nil, errors.New("StartInstances require credential")
    }

    request.SetContext(ctx)
    
    response = NewStartInstancesResponse()
    err = c.Send(request, response)
    return
}

func NewStopInstancesRequest() (request *StopInstancesRequest) {
    request = &StopInstancesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "StopInstances")
    
    
    return
}

func NewStopInstancesResponse() (response *StopInstancesResponse) {
    response = &StopInstancesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// StopInstances
// This API is used to shut down instances.
//
// 
//
// * You can only perform this operation on instances whose status is `RUNNING`.
//
// * The instance status will become `STOPPING` when the API is called successfully and `STOPPED` when the instance is successfully shut down.
//
// * Forced shutdown is supported. A forced shutdown is similar to switching off the power of a physical computer. It may cause data loss or file system corruption. Be sure to only force shut down a CVM when it cannot be sht down normally.
//
// * Batch operations are supported. The maximum number of instances in each request is 100.
//
// error code that may be returned:
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETERCOMBINATION = "InvalidParameterCombination"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_ZONENOTSUPPORTED = "InvalidParameterValue.ZoneNotSupported"
//  MISSINGPARAMETER = "MissingParameter"
//  MUTEXOPERATION_TASKRUNNING = "MutexOperation.TaskRunning"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  UNAUTHORIZEDOPERATION_MFAEXPIRED = "UnauthorizedOperation.MFAExpired"
//  UNAUTHORIZEDOPERATION_MFANOTFOUND = "UnauthorizedOperation.MFANotFound"
//  UNSUPPORTEDOPERATION_HIBERNATIONFORNORMALINSTANCE = "UnsupportedOperation.HibernationForNormalInstance"
//  UNSUPPORTEDOPERATION_INSTANCEREINSTALLFAILED = "UnsupportedOperation.InstanceReinstallFailed"
//  UNSUPPORTEDOPERATION_INSTANCESTATEBANNING = "UnsupportedOperation.InstanceStateBanning"
//  UNSUPPORTEDOPERATION_INSTANCESTATECORRUPTED = "UnsupportedOperation.InstanceStateCorrupted"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERRESCUEMODE = "UnsupportedOperation.InstanceStateEnterRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERSERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateEnterServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITRESCUEMODE = "UnsupportedOperation.InstanceStateExitRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITSERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateExitServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATEFREEZING = "UnsupportedOperation.InstanceStateFreezing"
//  UNSUPPORTEDOPERATION_INSTANCESTATEISOLATING = "UnsupportedOperation.InstanceStateIsolating"
//  UNSUPPORTEDOPERATION_INSTANCESTATEPENDING = "UnsupportedOperation.InstanceStatePending"
//  UNSUPPORTEDOPERATION_INSTANCESTATEREBOOTING = "UnsupportedOperation.InstanceStateRebooting"
//  UNSUPPORTEDOPERATION_INSTANCESTATERESCUEMODE = "UnsupportedOperation.InstanceStateRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATESERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATESHUTDOWN = "UnsupportedOperation.InstanceStateShutdown"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTARTING = "UnsupportedOperation.InstanceStateStarting"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPED = "UnsupportedOperation.InstanceStateStopped"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPING = "UnsupportedOperation.InstanceStateStopping"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATED = "UnsupportedOperation.InstanceStateTerminated"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATING = "UnsupportedOperation.InstanceStateTerminating"
//  UNSUPPORTEDOPERATION_STOPPEDMODESTOPCHARGING = "UnsupportedOperation.StoppedModeStopCharging"
func (c *Client) StopInstances(request *StopInstancesRequest) (response *StopInstancesResponse, err error) {
    return c.StopInstancesWithContext(context.Background(), request)
}

// StopInstances
// This API is used to shut down instances.
//
// 
//
// * You can only perform this operation on instances whose status is `RUNNING`.
//
// * The instance status will become `STOPPING` when the API is called successfully and `STOPPED` when the instance is successfully shut down.
//
// * Forced shutdown is supported. A forced shutdown is similar to switching off the power of a physical computer. It may cause data loss or file system corruption. Be sure to only force shut down a CVM when it cannot be sht down normally.
//
// * Batch operations are supported. The maximum number of instances in each request is 100.
//
// error code that may be returned:
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDPARAMETERCOMBINATION = "InvalidParameterCombination"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_ZONENOTSUPPORTED = "InvalidParameterValue.ZoneNotSupported"
//  MISSINGPARAMETER = "MissingParameter"
//  MUTEXOPERATION_TASKRUNNING = "MutexOperation.TaskRunning"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  UNAUTHORIZEDOPERATION_MFAEXPIRED = "UnauthorizedOperation.MFAExpired"
//  UNAUTHORIZEDOPERATION_MFANOTFOUND = "UnauthorizedOperation.MFANotFound"
//  UNSUPPORTEDOPERATION_HIBERNATIONFORNORMALINSTANCE = "UnsupportedOperation.HibernationForNormalInstance"
//  UNSUPPORTEDOPERATION_INSTANCEREINSTALLFAILED = "UnsupportedOperation.InstanceReinstallFailed"
//  UNSUPPORTEDOPERATION_INSTANCESTATEBANNING = "UnsupportedOperation.InstanceStateBanning"
//  UNSUPPORTEDOPERATION_INSTANCESTATECORRUPTED = "UnsupportedOperation.InstanceStateCorrupted"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERRESCUEMODE = "UnsupportedOperation.InstanceStateEnterRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERSERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateEnterServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITRESCUEMODE = "UnsupportedOperation.InstanceStateExitRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITSERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateExitServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATEFREEZING = "UnsupportedOperation.InstanceStateFreezing"
//  UNSUPPORTEDOPERATION_INSTANCESTATEISOLATING = "UnsupportedOperation.InstanceStateIsolating"
//  UNSUPPORTEDOPERATION_INSTANCESTATEPENDING = "UnsupportedOperation.InstanceStatePending"
//  UNSUPPORTEDOPERATION_INSTANCESTATEREBOOTING = "UnsupportedOperation.InstanceStateRebooting"
//  UNSUPPORTEDOPERATION_INSTANCESTATERESCUEMODE = "UnsupportedOperation.InstanceStateRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATESERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATESHUTDOWN = "UnsupportedOperation.InstanceStateShutdown"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTARTING = "UnsupportedOperation.InstanceStateStarting"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPED = "UnsupportedOperation.InstanceStateStopped"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPING = "UnsupportedOperation.InstanceStateStopping"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATED = "UnsupportedOperation.InstanceStateTerminated"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATING = "UnsupportedOperation.InstanceStateTerminating"
//  UNSUPPORTEDOPERATION_STOPPEDMODESTOPCHARGING = "UnsupportedOperation.StoppedModeStopCharging"
func (c *Client) StopInstancesWithContext(ctx context.Context, request *StopInstancesRequest) (response *StopInstancesResponse, err error) {
    if request == nil {
        request = NewStopInstancesRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "StopInstances")
    
    if c.GetCredential() == nil {
        return nil, errors.New("StopInstances require credential")
    }

    request.SetContext(ctx)
    
    response = NewStopInstancesResponse()
    err = c.Send(request, response)
    return
}

func NewSyncImagesRequest() (request *SyncImagesRequest) {
    request = &SyncImagesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "SyncImages")
    
    
    return
}

func NewSyncImagesResponse() (response *SyncImagesResponse) {
    response = &SyncImagesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// SyncImages
// This API is used to synchronize custom images to other regions.
//
// 
//
// * This API only supports synchronizing one image per call.
//
// * This API supports multiple synchronization regions.
//
// * A single account can have a maximum of 500 custom images in each region.
//
// error code that may be returned:
//  FAILEDOPERATION_INVALIDIMAGESTATE = "FailedOperation.InvalidImageState"
//  IMAGEQUOTALIMITEXCEEDED = "ImageQuotaLimitExceeded"
//  INVALIDIMAGEID_INCORRECTSTATE = "InvalidImageId.IncorrectState"
//  INVALIDIMAGEID_MALFORMED = "InvalidImageId.Malformed"
//  INVALIDIMAGEID_NOTFOUND = "InvalidImageId.NotFound"
//  INVALIDIMAGEID_TOOLARGE = "InvalidImageId.TooLarge"
//  INVALIDIMAGENAME_DUPLICATE = "InvalidImageName.Duplicate"
//  INVALIDPARAMETER_INSTANCEIMAGENOTSUPPORT = "InvalidParameter.InstanceImageNotSupport"
//  INVALIDPARAMETER_INVALIDKMSKEYID = "InvalidParameter.InvalidKmsKeyId"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEID = "InvalidParameterValue.InvalidImageId"
//  INVALIDPARAMETERVALUE_INVALIDIMAGESTATE = "InvalidParameterValue.InvalidImageState"
//  INVALIDPARAMETERVALUE_INVALIDREGION = "InvalidParameterValue.InvalidRegion"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDREGION_NOTFOUND = "InvalidRegion.NotFound"
//  INVALIDREGION_UNAVAILABLE = "InvalidRegion.Unavailable"
//  UNSUPPORTEDOPERATION_ENCRYPTEDIMAGESNOTSUPPORTED = "UnsupportedOperation.EncryptedImagesNotSupported"
//  UNSUPPORTEDOPERATION_LOCATIONIMAGENOTSUPPORTED = "UnsupportedOperation.LocationImageNotSupported"
//  UNSUPPORTEDOPERATION_REGION = "UnsupportedOperation.Region"
//  UNSUPPORTEDOPERATION_SYNCENCRYPTIMAGENOTSUPPORT = "UnsupportedOperation.SyncEncryptImageNotSupport"
func (c *Client) SyncImages(request *SyncImagesRequest) (response *SyncImagesResponse, err error) {
    return c.SyncImagesWithContext(context.Background(), request)
}

// SyncImages
// This API is used to synchronize custom images to other regions.
//
// 
//
// * This API only supports synchronizing one image per call.
//
// * This API supports multiple synchronization regions.
//
// * A single account can have a maximum of 500 custom images in each region.
//
// error code that may be returned:
//  FAILEDOPERATION_INVALIDIMAGESTATE = "FailedOperation.InvalidImageState"
//  IMAGEQUOTALIMITEXCEEDED = "ImageQuotaLimitExceeded"
//  INVALIDIMAGEID_INCORRECTSTATE = "InvalidImageId.IncorrectState"
//  INVALIDIMAGEID_MALFORMED = "InvalidImageId.Malformed"
//  INVALIDIMAGEID_NOTFOUND = "InvalidImageId.NotFound"
//  INVALIDIMAGEID_TOOLARGE = "InvalidImageId.TooLarge"
//  INVALIDIMAGENAME_DUPLICATE = "InvalidImageName.Duplicate"
//  INVALIDPARAMETER_INSTANCEIMAGENOTSUPPORT = "InvalidParameter.InstanceImageNotSupport"
//  INVALIDPARAMETER_INVALIDKMSKEYID = "InvalidParameter.InvalidKmsKeyId"
//  INVALIDPARAMETERVALUE_INVALIDIMAGEID = "InvalidParameterValue.InvalidImageId"
//  INVALIDPARAMETERVALUE_INVALIDIMAGESTATE = "InvalidParameterValue.InvalidImageState"
//  INVALIDPARAMETERVALUE_INVALIDREGION = "InvalidParameterValue.InvalidRegion"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDREGION_NOTFOUND = "InvalidRegion.NotFound"
//  INVALIDREGION_UNAVAILABLE = "InvalidRegion.Unavailable"
//  UNSUPPORTEDOPERATION_ENCRYPTEDIMAGESNOTSUPPORTED = "UnsupportedOperation.EncryptedImagesNotSupported"
//  UNSUPPORTEDOPERATION_LOCATIONIMAGENOTSUPPORTED = "UnsupportedOperation.LocationImageNotSupported"
//  UNSUPPORTEDOPERATION_REGION = "UnsupportedOperation.Region"
//  UNSUPPORTEDOPERATION_SYNCENCRYPTIMAGENOTSUPPORT = "UnsupportedOperation.SyncEncryptImageNotSupport"
func (c *Client) SyncImagesWithContext(ctx context.Context, request *SyncImagesRequest) (response *SyncImagesResponse, err error) {
    if request == nil {
        request = NewSyncImagesRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "SyncImages")
    
    if c.GetCredential() == nil {
        return nil, errors.New("SyncImages require credential")
    }

    request.SetContext(ctx)
    
    response = NewSyncImagesResponse()
    err = c.Send(request, response)
    return
}

func NewTerminateInstancesRequest() (request *TerminateInstancesRequest) {
    request = &TerminateInstancesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("cvm", APIVersion, "TerminateInstances")
    
    
    return
}

func NewTerminateInstancesResponse() (response *TerminateInstancesResponse) {
    response = &TerminateInstancesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// TerminateInstances
// This API is used to return instances.
//
// 
//
// * Use this API to return instances that are no longer required.
//
// * Pay-as-you-go instances can be returned directly through this API.
//
// * Batch operations are supported. The allowed maximum number of instances in each request is 100.
//
// error code that may be returned:
//  FAILEDOPERATION_INVALIDINSTANCEAPPLICATIONROLE = "FailedOperation.InvalidInstanceApplicationRole"
//  FAILEDOPERATION_INVALIDINSTANCEAPPLICATIONROLEEMR = "FailedOperation.InvalidInstanceApplicationRoleEmr"
//  FAILEDOPERATION_UNRETURNABLE = "FailedOperation.Unreturnable"
//  INTERNALERROR_TRADEUNKNOWNERROR = "InternalError.TradeUnknownError"
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDACCOUNT_UNPAIDORDER = "InvalidAccount.UnpaidOrder"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDINSTANCENOTSUPPORTEDPREPAIDINSTANCE = "InvalidInstanceNotSupportedPrepaidInstance"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_NOTSUPPORTED = "InvalidParameterValue.NotSupported"
//  LIMITEXCEEDED_USERRETURNQUOTA = "LimitExceeded.UserReturnQuota"
//  MISSINGPARAMETER = "MissingParameter"
//  MUTEXOPERATION_TASKRUNNING = "MutexOperation.TaskRunning"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  UNAUTHORIZEDOPERATION_INVALIDTOKEN = "UnauthorizedOperation.InvalidToken"
//  UNAUTHORIZEDOPERATION_MFAEXPIRED = "UnauthorizedOperation.MFAExpired"
//  UNAUTHORIZEDOPERATION_MFANOTFOUND = "UnauthorizedOperation.MFANotFound"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
//  UNSUPPORTEDOPERATION_INSTANCECHARGETYPE = "UnsupportedOperation.InstanceChargeType"
//  UNSUPPORTEDOPERATION_INSTANCEMIXEDPRICINGMODEL = "UnsupportedOperation.InstanceMixedPricingModel"
//  UNSUPPORTEDOPERATION_INSTANCESTATEBANNING = "UnsupportedOperation.InstanceStateBanning"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERRESCUEMODE = "UnsupportedOperation.InstanceStateEnterRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERSERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateEnterServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITRESCUEMODE = "UnsupportedOperation.InstanceStateExitRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITSERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateExitServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATEFREEZING = "UnsupportedOperation.InstanceStateFreezing"
//  UNSUPPORTEDOPERATION_INSTANCESTATEISOLATING = "UnsupportedOperation.InstanceStateIsolating"
//  UNSUPPORTEDOPERATION_INSTANCESTATELAUNCHFAILED = "UnsupportedOperation.InstanceStateLaunchFailed"
//  UNSUPPORTEDOPERATION_INSTANCESTATEPENDING = "UnsupportedOperation.InstanceStatePending"
//  UNSUPPORTEDOPERATION_INSTANCESTATEREBOOTING = "UnsupportedOperation.InstanceStateRebooting"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTARTING = "UnsupportedOperation.InstanceStateStarting"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPING = "UnsupportedOperation.InstanceStateStopping"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATED = "UnsupportedOperation.InstanceStateTerminated"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATING = "UnsupportedOperation.InstanceStateTerminating"
//  UNSUPPORTEDOPERATION_INSTANCESPROTECTED = "UnsupportedOperation.InstancesProtected"
//  UNSUPPORTEDOPERATION_INVALIDINSTANCESOWNER = "UnsupportedOperation.InvalidInstancesOwner"
//  UNSUPPORTEDOPERATION_REDHATINSTANCETERMINATEUNSUPPORTED = "UnsupportedOperation.RedHatInstanceTerminateUnsupported"
//  UNSUPPORTEDOPERATION_REDHATINSTANCEUNSUPPORTED = "UnsupportedOperation.RedHatInstanceUnsupported"
//  UNSUPPORTEDOPERATION_REGION = "UnsupportedOperation.Region"
//  UNSUPPORTEDOPERATION_SPECIALINSTANCETYPE = "UnsupportedOperation.SpecialInstanceType"
//  UNSUPPORTEDOPERATION_USERLIMITOPERATIONEXCEEDQUOTA = "UnsupportedOperation.UserLimitOperationExceedQuota"
func (c *Client) TerminateInstances(request *TerminateInstancesRequest) (response *TerminateInstancesResponse, err error) {
    return c.TerminateInstancesWithContext(context.Background(), request)
}

// TerminateInstances
// This API is used to return instances.
//
// 
//
// * Use this API to return instances that are no longer required.
//
// * Pay-as-you-go instances can be returned directly through this API.
//
// * Batch operations are supported. The allowed maximum number of instances in each request is 100.
//
// error code that may be returned:
//  FAILEDOPERATION_INVALIDINSTANCEAPPLICATIONROLE = "FailedOperation.InvalidInstanceApplicationRole"
//  FAILEDOPERATION_INVALIDINSTANCEAPPLICATIONROLEEMR = "FailedOperation.InvalidInstanceApplicationRoleEmr"
//  FAILEDOPERATION_UNRETURNABLE = "FailedOperation.Unreturnable"
//  INTERNALERROR_TRADEUNKNOWNERROR = "InternalError.TradeUnknownError"
//  INTERNALSERVERERROR = "InternalServerError"
//  INVALIDACCOUNT_UNPAIDORDER = "InvalidAccount.UnpaidOrder"
//  INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"
//  INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"
//  INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"
//  INVALIDINSTANCENOTSUPPORTEDPREPAIDINSTANCE = "InvalidInstanceNotSupportedPrepaidInstance"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"
//  INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"
//  INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"
//  INVALIDPARAMETERVALUE_NOTSUPPORTED = "InvalidParameterValue.NotSupported"
//  LIMITEXCEEDED_USERRETURNQUOTA = "LimitExceeded.UserReturnQuota"
//  MISSINGPARAMETER = "MissingParameter"
//  MUTEXOPERATION_TASKRUNNING = "MutexOperation.TaskRunning"
//  OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"
//  UNAUTHORIZEDOPERATION_INVALIDTOKEN = "UnauthorizedOperation.InvalidToken"
//  UNAUTHORIZEDOPERATION_MFAEXPIRED = "UnauthorizedOperation.MFAExpired"
//  UNAUTHORIZEDOPERATION_MFANOTFOUND = "UnauthorizedOperation.MFANotFound"
//  UNSUPPORTEDOPERATION = "UnsupportedOperation"
//  UNSUPPORTEDOPERATION_INSTANCECHARGETYPE = "UnsupportedOperation.InstanceChargeType"
//  UNSUPPORTEDOPERATION_INSTANCEMIXEDPRICINGMODEL = "UnsupportedOperation.InstanceMixedPricingModel"
//  UNSUPPORTEDOPERATION_INSTANCESTATEBANNING = "UnsupportedOperation.InstanceStateBanning"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERRESCUEMODE = "UnsupportedOperation.InstanceStateEnterRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEENTERSERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateEnterServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITRESCUEMODE = "UnsupportedOperation.InstanceStateExitRescueMode"
//  UNSUPPORTEDOPERATION_INSTANCESTATEEXITSERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateExitServiceLiveMigrate"
//  UNSUPPORTEDOPERATION_INSTANCESTATEFREEZING = "UnsupportedOperation.InstanceStateFreezing"
//  UNSUPPORTEDOPERATION_INSTANCESTATEISOLATING = "UnsupportedOperation.InstanceStateIsolating"
//  UNSUPPORTEDOPERATION_INSTANCESTATELAUNCHFAILED = "UnsupportedOperation.InstanceStateLaunchFailed"
//  UNSUPPORTEDOPERATION_INSTANCESTATEPENDING = "UnsupportedOperation.InstanceStatePending"
//  UNSUPPORTEDOPERATION_INSTANCESTATEREBOOTING = "UnsupportedOperation.InstanceStateRebooting"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTARTING = "UnsupportedOperation.InstanceStateStarting"
//  UNSUPPORTEDOPERATION_INSTANCESTATESTOPPING = "UnsupportedOperation.InstanceStateStopping"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATED = "UnsupportedOperation.InstanceStateTerminated"
//  UNSUPPORTEDOPERATION_INSTANCESTATETERMINATING = "UnsupportedOperation.InstanceStateTerminating"
//  UNSUPPORTEDOPERATION_INSTANCESPROTECTED = "UnsupportedOperation.InstancesProtected"
//  UNSUPPORTEDOPERATION_INVALIDINSTANCESOWNER = "UnsupportedOperation.InvalidInstancesOwner"
//  UNSUPPORTEDOPERATION_REDHATINSTANCETERMINATEUNSUPPORTED = "UnsupportedOperation.RedHatInstanceTerminateUnsupported"
//  UNSUPPORTEDOPERATION_REDHATINSTANCEUNSUPPORTED = "UnsupportedOperation.RedHatInstanceUnsupported"
//  UNSUPPORTEDOPERATION_REGION = "UnsupportedOperation.Region"
//  UNSUPPORTEDOPERATION_SPECIALINSTANCETYPE = "UnsupportedOperation.SpecialInstanceType"
//  UNSUPPORTEDOPERATION_USERLIMITOPERATIONEXCEEDQUOTA = "UnsupportedOperation.UserLimitOperationExceedQuota"
func (c *Client) TerminateInstancesWithContext(ctx context.Context, request *TerminateInstancesRequest) (response *TerminateInstancesResponse, err error) {
    if request == nil {
        request = NewTerminateInstancesRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "cvm", APIVersion, "TerminateInstances")
    
    if c.GetCredential() == nil {
        return nil, errors.New("TerminateInstances require credential")
    }

    request.SetContext(ctx)
    
    response = NewTerminateInstancesResponse()
    err = c.Send(request, response)
    return
}
