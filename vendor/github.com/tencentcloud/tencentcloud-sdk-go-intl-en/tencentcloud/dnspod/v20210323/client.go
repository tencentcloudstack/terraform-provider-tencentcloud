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
    "context"
    "errors"
    "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
    tchttp "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common/http"
    "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common/profile"
)

const APIVersion = "2021-03-23"

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


func NewCreateDomainRequest() (request *CreateDomainRequest) {
    request = &CreateDomainRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("dnspod", APIVersion, "CreateDomain")
    
    
    return
}

func NewCreateDomainResponse() (response *CreateDomainResponse) {
    response = &CreateDomainResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// CreateDomain
// This API is used to add a domain.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_DOMAINEXISTS = "FailedOperation.DomainExists"
//  FAILEDOPERATION_DOMAININENTERPRISEMAILACCOUNT = "FailedOperation.DomainInEnterpriseMailAccount"
//  FAILEDOPERATION_DOMAINOWNEDBYOTHERUSER = "FailedOperation.DomainOwnedByOtherUser"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_ACCOUNTISBANNED = "InvalidParameter.AccountIsBanned"
//  INVALIDPARAMETER_DOMAININBLACKLIST = "InvalidParameter.DomainInBlackList"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETER_DOMAINISMYALIAS = "InvalidParameter.DomainIsMyAlias"
//  INVALIDPARAMETER_DOMAINNOTREGED = "InvalidParameter.DomainNotReged"
//  INVALIDPARAMETER_EMAILNOTVERIFIED = "InvalidParameter.EmailNotVerified"
//  INVALIDPARAMETER_OPERATEFAILED = "InvalidParameter.OperateFailed"
//  INVALIDPARAMETER_QUHUITXTNOTMATCH = "InvalidParameter.QuhuiTxtNotMatch"
//  INVALIDPARAMETER_QUHUITXTRECORDWAIT = "InvalidParameter.QuhuiTxtRecordWait"
//  INVALIDPARAMETER_TOOLSDOMAININVALID = "InvalidParameter.ToolsDomainInvalid"
//  INVALIDPARAMETER_UNREALNAMEUSER = "InvalidParameter.UnrealNameUser"
func (c *Client) CreateDomain(request *CreateDomainRequest) (response *CreateDomainResponse, err error) {
    return c.CreateDomainWithContext(context.Background(), request)
}

// CreateDomain
// This API is used to add a domain.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_DOMAINEXISTS = "FailedOperation.DomainExists"
//  FAILEDOPERATION_DOMAININENTERPRISEMAILACCOUNT = "FailedOperation.DomainInEnterpriseMailAccount"
//  FAILEDOPERATION_DOMAINOWNEDBYOTHERUSER = "FailedOperation.DomainOwnedByOtherUser"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_ACCOUNTISBANNED = "InvalidParameter.AccountIsBanned"
//  INVALIDPARAMETER_DOMAININBLACKLIST = "InvalidParameter.DomainInBlackList"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETER_DOMAINISMYALIAS = "InvalidParameter.DomainIsMyAlias"
//  INVALIDPARAMETER_DOMAINNOTREGED = "InvalidParameter.DomainNotReged"
//  INVALIDPARAMETER_EMAILNOTVERIFIED = "InvalidParameter.EmailNotVerified"
//  INVALIDPARAMETER_OPERATEFAILED = "InvalidParameter.OperateFailed"
//  INVALIDPARAMETER_QUHUITXTNOTMATCH = "InvalidParameter.QuhuiTxtNotMatch"
//  INVALIDPARAMETER_QUHUITXTRECORDWAIT = "InvalidParameter.QuhuiTxtRecordWait"
//  INVALIDPARAMETER_TOOLSDOMAININVALID = "InvalidParameter.ToolsDomainInvalid"
//  INVALIDPARAMETER_UNREALNAMEUSER = "InvalidParameter.UnrealNameUser"
func (c *Client) CreateDomainWithContext(ctx context.Context, request *CreateDomainRequest) (response *CreateDomainResponse, err error) {
    if request == nil {
        request = NewCreateDomainRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "dnspod", APIVersion, "CreateDomain")
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateDomain require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateDomainResponse()
    err = c.Send(request, response)
    return
}

func NewCreateDomainAliasRequest() (request *CreateDomainAliasRequest) {
    request = &CreateDomainAliasRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("dnspod", APIVersion, "CreateDomainAlias")
    
    
    return
}

func NewCreateDomainAliasResponse() (response *CreateDomainAliasResponse) {
    response = &CreateDomainAliasResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// CreateDomainAlias
// This API is used to create a domain alias.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_DOMAINEXISTS = "FailedOperation.DomainExists"
//  FAILEDOPERATION_DOMAINISLOCKED = "FailedOperation.DomainIsLocked"
//  FAILEDOPERATION_DOMAINISSPAM = "FailedOperation.DomainIsSpam"
//  FAILEDOPERATION_DOMAINOWNEDBYOTHERUSER = "FailedOperation.DomainOwnedByOtherUser"
//  FAILEDOPERATION_NOTREALNAMEDUSER = "FailedOperation.NotRealNamedUser"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_DOMAINALIASEXISTS = "InvalidParameter.DomainAliasExists"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININEFFECTORINVALIDATED = "InvalidParameter.DomainInEffectOrInvalidated"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  LIMITEXCEEDED_DOMAINALIASCOUNTEXCEEDED = "LimitExceeded.DomainAliasCountExceeded"
//  LIMITEXCEEDED_DOMAINALIASNUMBERLIMIT = "LimitExceeded.DomainAliasNumberLimit"
//  MISSINGPARAMETER = "MissingParameter"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
//  RESOURCENOTFOUND = "ResourceNotFound"
func (c *Client) CreateDomainAlias(request *CreateDomainAliasRequest) (response *CreateDomainAliasResponse, err error) {
    return c.CreateDomainAliasWithContext(context.Background(), request)
}

// CreateDomainAlias
// This API is used to create a domain alias.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_DOMAINEXISTS = "FailedOperation.DomainExists"
//  FAILEDOPERATION_DOMAINISLOCKED = "FailedOperation.DomainIsLocked"
//  FAILEDOPERATION_DOMAINISSPAM = "FailedOperation.DomainIsSpam"
//  FAILEDOPERATION_DOMAINOWNEDBYOTHERUSER = "FailedOperation.DomainOwnedByOtherUser"
//  FAILEDOPERATION_NOTREALNAMEDUSER = "FailedOperation.NotRealNamedUser"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_DOMAINALIASEXISTS = "InvalidParameter.DomainAliasExists"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININEFFECTORINVALIDATED = "InvalidParameter.DomainInEffectOrInvalidated"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  LIMITEXCEEDED_DOMAINALIASCOUNTEXCEEDED = "LimitExceeded.DomainAliasCountExceeded"
//  LIMITEXCEEDED_DOMAINALIASNUMBERLIMIT = "LimitExceeded.DomainAliasNumberLimit"
//  MISSINGPARAMETER = "MissingParameter"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
//  RESOURCENOTFOUND = "ResourceNotFound"
func (c *Client) CreateDomainAliasWithContext(ctx context.Context, request *CreateDomainAliasRequest) (response *CreateDomainAliasResponse, err error) {
    if request == nil {
        request = NewCreateDomainAliasRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "dnspod", APIVersion, "CreateDomainAlias")
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateDomainAlias require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateDomainAliasResponse()
    err = c.Send(request, response)
    return
}

func NewCreateDomainBatchRequest() (request *CreateDomainBatchRequest) {
    request = &CreateDomainBatchRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("dnspod", APIVersion, "CreateDomainBatch")
    
    
    return
}

func NewCreateDomainBatchResponse() (response *CreateDomainBatchResponse) {
    response = &CreateDomainBatchResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// CreateDomainBatch
// This API is used to bulk add domains.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_DOMAINISLOCKED = "FailedOperation.DomainIsLocked"
//  FAILEDOPERATION_DOMAINISSPAM = "FailedOperation.DomainIsSpam"
//  INVALIDPARAMETER_BATCHDOMAINCREATEACTIONERROR = "InvalidParameter.BatchDomainCreateActionError"
//  INVALIDPARAMETER_BATCHTASKCOUNTLIMIT = "InvalidParameter.BatchTaskCountLimit"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETER_DOMAINSEMPTY = "InvalidParameter.DomainsEmpty"
//  INVALIDPARAMETER_JOBGREATERTHANLIMIT = "InvalidParameter.JobGreaterThanLimit"
//  INVALIDPARAMETER_RECORDVALUEINVALID = "InvalidParameter.RecordValueInvalid"
//  INVALIDPARAMETER_RECORDVALUELENGTHINVALID = "InvalidParameter.RecordValueLengthInvalid"
//  INVALIDPARAMETER_TOOMANYINVALIDDOMAINS = "InvalidParameter.TooManyInvalidDomains"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  OPERATIONDENIED_IPINBLACKLISTNOTALLOWED = "OperationDenied.IPInBlacklistNotAllowed"
//  REQUESTLIMITEXCEEDED_BATCHTASKLIMIT = "RequestLimitExceeded.BatchTaskLimit"
//  REQUESTLIMITEXCEEDED_CREATEDOMAINLIMIT = "RequestLimitExceeded.CreateDomainLimit"
//  REQUESTLIMITEXCEEDED_REQUESTLIMITEXCEEDED = "RequestLimitExceeded.RequestLimitExceeded"
func (c *Client) CreateDomainBatch(request *CreateDomainBatchRequest) (response *CreateDomainBatchResponse, err error) {
    return c.CreateDomainBatchWithContext(context.Background(), request)
}

// CreateDomainBatch
// This API is used to bulk add domains.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_DOMAINISLOCKED = "FailedOperation.DomainIsLocked"
//  FAILEDOPERATION_DOMAINISSPAM = "FailedOperation.DomainIsSpam"
//  INVALIDPARAMETER_BATCHDOMAINCREATEACTIONERROR = "InvalidParameter.BatchDomainCreateActionError"
//  INVALIDPARAMETER_BATCHTASKCOUNTLIMIT = "InvalidParameter.BatchTaskCountLimit"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETER_DOMAINSEMPTY = "InvalidParameter.DomainsEmpty"
//  INVALIDPARAMETER_JOBGREATERTHANLIMIT = "InvalidParameter.JobGreaterThanLimit"
//  INVALIDPARAMETER_RECORDVALUEINVALID = "InvalidParameter.RecordValueInvalid"
//  INVALIDPARAMETER_RECORDVALUELENGTHINVALID = "InvalidParameter.RecordValueLengthInvalid"
//  INVALIDPARAMETER_TOOMANYINVALIDDOMAINS = "InvalidParameter.TooManyInvalidDomains"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  OPERATIONDENIED_IPINBLACKLISTNOTALLOWED = "OperationDenied.IPInBlacklistNotAllowed"
//  REQUESTLIMITEXCEEDED_BATCHTASKLIMIT = "RequestLimitExceeded.BatchTaskLimit"
//  REQUESTLIMITEXCEEDED_CREATEDOMAINLIMIT = "RequestLimitExceeded.CreateDomainLimit"
//  REQUESTLIMITEXCEEDED_REQUESTLIMITEXCEEDED = "RequestLimitExceeded.RequestLimitExceeded"
func (c *Client) CreateDomainBatchWithContext(ctx context.Context, request *CreateDomainBatchRequest) (response *CreateDomainBatchResponse, err error) {
    if request == nil {
        request = NewCreateDomainBatchRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "dnspod", APIVersion, "CreateDomainBatch")
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateDomainBatch require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateDomainBatchResponse()
    err = c.Send(request, response)
    return
}

func NewCreateDomainGroupRequest() (request *CreateDomainGroupRequest) {
    request = &CreateDomainGroupRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("dnspod", APIVersion, "CreateDomainGroup")
    
    
    return
}

func NewCreateDomainGroupResponse() (response *CreateDomainGroupResponse) {
    response = &CreateDomainGroupResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// CreateDomainGroup
// This API is used to create a domain group.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_NOTREALNAMEDUSER = "FailedOperation.NotRealNamedUser"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_GROUPNAMEEXISTS = "InvalidParameter.GroupNameExists"
//  INVALIDPARAMETER_GROUPNAMEINVALID = "InvalidParameter.GroupNameInvalid"
//  LIMITEXCEEDED_GROUPNUMBERLIMIT = "LimitExceeded.GroupNumberLimit"
//  MISSINGPARAMETER = "MissingParameter"
//  RESOURCENOTFOUND = "ResourceNotFound"
func (c *Client) CreateDomainGroup(request *CreateDomainGroupRequest) (response *CreateDomainGroupResponse, err error) {
    return c.CreateDomainGroupWithContext(context.Background(), request)
}

// CreateDomainGroup
// This API is used to create a domain group.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_NOTREALNAMEDUSER = "FailedOperation.NotRealNamedUser"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_GROUPNAMEEXISTS = "InvalidParameter.GroupNameExists"
//  INVALIDPARAMETER_GROUPNAMEINVALID = "InvalidParameter.GroupNameInvalid"
//  LIMITEXCEEDED_GROUPNUMBERLIMIT = "LimitExceeded.GroupNumberLimit"
//  MISSINGPARAMETER = "MissingParameter"
//  RESOURCENOTFOUND = "ResourceNotFound"
func (c *Client) CreateDomainGroupWithContext(ctx context.Context, request *CreateDomainGroupRequest) (response *CreateDomainGroupResponse, err error) {
    if request == nil {
        request = NewCreateDomainGroupRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "dnspod", APIVersion, "CreateDomainGroup")
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateDomainGroup require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateDomainGroupResponse()
    err = c.Send(request, response)
    return
}

func NewCreatePackageOrderRequest() (request *CreatePackageOrderRequest) {
    request = &CreatePackageOrderRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("dnspod", APIVersion, "CreatePackageOrder")
    
    
    return
}

func NewCreatePackageOrderResponse() (response *CreatePackageOrderResponse) {
    response = &CreatePackageOrderResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// CreatePackageOrder
// This API is used to enable a paid plan on the international website.
//
// error code that may be returned:
//  FAILEDOPERATION_NOTBINDCREDITCARD = "FailedOperation.NotBindCreditCard"
//  FAILEDOPERATION_POSTPAYPAYMENTNOTOPEN = "FailedOperation.PostPayPaymentNotOpen"
func (c *Client) CreatePackageOrder(request *CreatePackageOrderRequest) (response *CreatePackageOrderResponse, err error) {
    return c.CreatePackageOrderWithContext(context.Background(), request)
}

// CreatePackageOrder
// This API is used to enable a paid plan on the international website.
//
// error code that may be returned:
//  FAILEDOPERATION_NOTBINDCREDITCARD = "FailedOperation.NotBindCreditCard"
//  FAILEDOPERATION_POSTPAYPAYMENTNOTOPEN = "FailedOperation.PostPayPaymentNotOpen"
func (c *Client) CreatePackageOrderWithContext(ctx context.Context, request *CreatePackageOrderRequest) (response *CreatePackageOrderResponse, err error) {
    if request == nil {
        request = NewCreatePackageOrderRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "dnspod", APIVersion, "CreatePackageOrder")
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreatePackageOrder require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreatePackageOrderResponse()
    err = c.Send(request, response)
    return
}

func NewCreateRecordRequest() (request *CreateRecordRequest) {
    request = &CreateRecordRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("dnspod", APIVersion, "CreateRecord")
    
    
    return
}

func NewCreateRecordResponse() (response *CreateRecordResponse) {
    response = &CreateRecordResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// CreateRecord
// This API is used to add a record.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_DOMAINISLOCKED = "FailedOperation.DomainIsLocked"
//  FAILEDOPERATION_DOMAINISSPAM = "FailedOperation.DomainIsSpam"
//  FAILEDOPERATION_LOGINAREANOTALLOWED = "FailedOperation.LoginAreaNotAllowed"
//  FAILEDOPERATION_LOGINFAILED = "FailedOperation.LoginFailed"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_ACCOUNTISBANNED = "InvalidParameter.AccountIsBanned"
//  INVALIDPARAMETER_CUSTOMMESSAGE = "InvalidParameter.CustomMessage"
//  INVALIDPARAMETER_DNSSECADDCNAMEERROR = "InvalidParameter.DnssecAddCnameError"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETER_DOMAINNOTALLOWEDMODIFYRECORDS = "InvalidParameter.DomainNotAllowedModifyRecords"
//  INVALIDPARAMETER_DOMAINNOTBEIAN = "InvalidParameter.DomainNotBeian"
//  INVALIDPARAMETER_DOMAINRECORDEXIST = "InvalidParameter.DomainRecordExist"
//  INVALIDPARAMETER_EMAILNOTVERIFIED = "InvalidParameter.EmailNotVerified"
//  INVALIDPARAMETER_INVALIDWEIGHT = "InvalidParameter.InvalidWeight"
//  INVALIDPARAMETER_LOGINTOKENIDERROR = "InvalidParameter.LoginTokenIdError"
//  INVALIDPARAMETER_LOGINTOKENNOTEXISTS = "InvalidParameter.LoginTokenNotExists"
//  INVALIDPARAMETER_LOGINTOKENVALIDATEFAILED = "InvalidParameter.LoginTokenValidateFailed"
//  INVALIDPARAMETER_MOBILENOTVERIFIED = "InvalidParameter.MobileNotVerified"
//  INVALIDPARAMETER_MXINVALID = "InvalidParameter.MxInvalid"
//  INVALIDPARAMETER_RECORDLINEINVALID = "InvalidParameter.RecordLineInvalid"
//  INVALIDPARAMETER_RECORDTYPEINVALID = "InvalidParameter.RecordTypeInvalid"
//  INVALIDPARAMETER_RECORDVALUEINVALID = "InvalidParameter.RecordValueInvalid"
//  INVALIDPARAMETER_RECORDVALUELENGTHINVALID = "InvalidParameter.RecordValueLengthInvalid"
//  INVALIDPARAMETER_REQUESTIPLIMITED = "InvalidParameter.RequestIpLimited"
//  INVALIDPARAMETER_SUBDOMAININVALID = "InvalidParameter.SubdomainInvalid"
//  INVALIDPARAMETER_UNREALNAMEUSER = "InvalidParameter.UnrealNameUser"
//  INVALIDPARAMETER_URLVALUEILLEGAL = "InvalidParameter.UrlValueIllegal"
//  INVALIDPARAMETER_USERNOTEXISTS = "InvalidParameter.UserNotExists"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  INVALIDPARAMETERVALUE_USERIDINVALID = "InvalidParameterValue.UserIdInvalid"
//  LIMITEXCEEDED_AAAACOUNTLIMIT = "LimitExceeded.AAAACountLimit"
//  LIMITEXCEEDED_ATNSRECORDLIMIT = "LimitExceeded.AtNsRecordLimit"
//  LIMITEXCEEDED_FAILEDLOGINLIMITEXCEEDED = "LimitExceeded.FailedLoginLimitExceeded"
//  LIMITEXCEEDED_HIDDENURLEXCEEDED = "LimitExceeded.HiddenUrlExceeded"
//  LIMITEXCEEDED_NSCOUNTLIMIT = "LimitExceeded.NsCountLimit"
//  LIMITEXCEEDED_RECORDTTLLIMIT = "LimitExceeded.RecordTtlLimit"
//  LIMITEXCEEDED_SRVCOUNTLIMIT = "LimitExceeded.SrvCountLimit"
//  LIMITEXCEEDED_SUBDOMAINLEVELLIMIT = "LimitExceeded.SubdomainLevelLimit"
//  LIMITEXCEEDED_SUBDOMAINROLLLIMIT = "LimitExceeded.SubdomainRollLimit"
//  LIMITEXCEEDED_SUBDOMAINWCARDLIMIT = "LimitExceeded.SubdomainWcardLimit"
//  LIMITEXCEEDED_URLCOUNTLIMIT = "LimitExceeded.UrlCountLimit"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_IPINBLACKLISTNOTALLOWED = "OperationDenied.IPInBlacklistNotAllowed"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
//  OPERATIONDENIED_NOTADMIN = "OperationDenied.NotAdmin"
//  OPERATIONDENIED_NOTAGENT = "OperationDenied.NotAgent"
//  OPERATIONDENIED_NOTMANAGEDUSER = "OperationDenied.NotManagedUser"
//  REQUESTLIMITEXCEEDED_REQUESTLIMITEXCEEDED = "RequestLimitExceeded.RequestLimitExceeded"
func (c *Client) CreateRecord(request *CreateRecordRequest) (response *CreateRecordResponse, err error) {
    return c.CreateRecordWithContext(context.Background(), request)
}

// CreateRecord
// This API is used to add a record.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_DOMAINISLOCKED = "FailedOperation.DomainIsLocked"
//  FAILEDOPERATION_DOMAINISSPAM = "FailedOperation.DomainIsSpam"
//  FAILEDOPERATION_LOGINAREANOTALLOWED = "FailedOperation.LoginAreaNotAllowed"
//  FAILEDOPERATION_LOGINFAILED = "FailedOperation.LoginFailed"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_ACCOUNTISBANNED = "InvalidParameter.AccountIsBanned"
//  INVALIDPARAMETER_CUSTOMMESSAGE = "InvalidParameter.CustomMessage"
//  INVALIDPARAMETER_DNSSECADDCNAMEERROR = "InvalidParameter.DnssecAddCnameError"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETER_DOMAINNOTALLOWEDMODIFYRECORDS = "InvalidParameter.DomainNotAllowedModifyRecords"
//  INVALIDPARAMETER_DOMAINNOTBEIAN = "InvalidParameter.DomainNotBeian"
//  INVALIDPARAMETER_DOMAINRECORDEXIST = "InvalidParameter.DomainRecordExist"
//  INVALIDPARAMETER_EMAILNOTVERIFIED = "InvalidParameter.EmailNotVerified"
//  INVALIDPARAMETER_INVALIDWEIGHT = "InvalidParameter.InvalidWeight"
//  INVALIDPARAMETER_LOGINTOKENIDERROR = "InvalidParameter.LoginTokenIdError"
//  INVALIDPARAMETER_LOGINTOKENNOTEXISTS = "InvalidParameter.LoginTokenNotExists"
//  INVALIDPARAMETER_LOGINTOKENVALIDATEFAILED = "InvalidParameter.LoginTokenValidateFailed"
//  INVALIDPARAMETER_MOBILENOTVERIFIED = "InvalidParameter.MobileNotVerified"
//  INVALIDPARAMETER_MXINVALID = "InvalidParameter.MxInvalid"
//  INVALIDPARAMETER_RECORDLINEINVALID = "InvalidParameter.RecordLineInvalid"
//  INVALIDPARAMETER_RECORDTYPEINVALID = "InvalidParameter.RecordTypeInvalid"
//  INVALIDPARAMETER_RECORDVALUEINVALID = "InvalidParameter.RecordValueInvalid"
//  INVALIDPARAMETER_RECORDVALUELENGTHINVALID = "InvalidParameter.RecordValueLengthInvalid"
//  INVALIDPARAMETER_REQUESTIPLIMITED = "InvalidParameter.RequestIpLimited"
//  INVALIDPARAMETER_SUBDOMAININVALID = "InvalidParameter.SubdomainInvalid"
//  INVALIDPARAMETER_UNREALNAMEUSER = "InvalidParameter.UnrealNameUser"
//  INVALIDPARAMETER_URLVALUEILLEGAL = "InvalidParameter.UrlValueIllegal"
//  INVALIDPARAMETER_USERNOTEXISTS = "InvalidParameter.UserNotExists"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  INVALIDPARAMETERVALUE_USERIDINVALID = "InvalidParameterValue.UserIdInvalid"
//  LIMITEXCEEDED_AAAACOUNTLIMIT = "LimitExceeded.AAAACountLimit"
//  LIMITEXCEEDED_ATNSRECORDLIMIT = "LimitExceeded.AtNsRecordLimit"
//  LIMITEXCEEDED_FAILEDLOGINLIMITEXCEEDED = "LimitExceeded.FailedLoginLimitExceeded"
//  LIMITEXCEEDED_HIDDENURLEXCEEDED = "LimitExceeded.HiddenUrlExceeded"
//  LIMITEXCEEDED_NSCOUNTLIMIT = "LimitExceeded.NsCountLimit"
//  LIMITEXCEEDED_RECORDTTLLIMIT = "LimitExceeded.RecordTtlLimit"
//  LIMITEXCEEDED_SRVCOUNTLIMIT = "LimitExceeded.SrvCountLimit"
//  LIMITEXCEEDED_SUBDOMAINLEVELLIMIT = "LimitExceeded.SubdomainLevelLimit"
//  LIMITEXCEEDED_SUBDOMAINROLLLIMIT = "LimitExceeded.SubdomainRollLimit"
//  LIMITEXCEEDED_SUBDOMAINWCARDLIMIT = "LimitExceeded.SubdomainWcardLimit"
//  LIMITEXCEEDED_URLCOUNTLIMIT = "LimitExceeded.UrlCountLimit"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_IPINBLACKLISTNOTALLOWED = "OperationDenied.IPInBlacklistNotAllowed"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
//  OPERATIONDENIED_NOTADMIN = "OperationDenied.NotAdmin"
//  OPERATIONDENIED_NOTAGENT = "OperationDenied.NotAgent"
//  OPERATIONDENIED_NOTMANAGEDUSER = "OperationDenied.NotManagedUser"
//  REQUESTLIMITEXCEEDED_REQUESTLIMITEXCEEDED = "RequestLimitExceeded.RequestLimitExceeded"
func (c *Client) CreateRecordWithContext(ctx context.Context, request *CreateRecordRequest) (response *CreateRecordResponse, err error) {
    if request == nil {
        request = NewCreateRecordRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "dnspod", APIVersion, "CreateRecord")
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateRecord require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateRecordResponse()
    err = c.Send(request, response)
    return
}

func NewCreateRecordBatchRequest() (request *CreateRecordBatchRequest) {
    request = &CreateRecordBatchRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("dnspod", APIVersion, "CreateRecordBatch")
    
    
    return
}

func NewCreateRecordBatchResponse() (response *CreateRecordBatchResponse) {
    response = &CreateRecordBatchResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// CreateRecordBatch
// This API is used to bulk add records.
//
// error code that may be returned:
//  FAILEDOPERATION_DOMAINISLOCKED = "FailedOperation.DomainIsLocked"
//  FAILEDOPERATION_DOMAINISSPAM = "FailedOperation.DomainIsSpam"
//  INVALIDPARAMETER_BATCHRECORDCREATEACTIONERROR = "InvalidParameter.BatchRecordCreateActionError"
//  INVALIDPARAMETER_BATCHTASKCOUNTLIMIT = "InvalidParameter.BatchTaskCountLimit"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETER_DOMAINSEMPTY = "InvalidParameter.DomainsEmpty"
//  INVALIDPARAMETER_JOBGREATERTHANLIMIT = "InvalidParameter.JobGreaterThanLimit"
//  INVALIDPARAMETER_PARAMSILLEGAL = "InvalidParameter.ParamsIllegal"
//  INVALIDPARAMETER_RECORDVALUEINVALID = "InvalidParameter.RecordValueInvalid"
//  INVALIDPARAMETER_RECORDVALUELENGTHINVALID = "InvalidParameter.RecordValueLengthInvalid"
//  INVALIDPARAMETER_RECORDSEMPTY = "InvalidParameter.RecordsEmpty"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  OPERATIONDENIED_IPINBLACKLISTNOTALLOWED = "OperationDenied.IPInBlacklistNotAllowed"
//  REQUESTLIMITEXCEEDED_BATCHTASKLIMIT = "RequestLimitExceeded.BatchTaskLimit"
//  REQUESTLIMITEXCEEDED_CREATEDOMAINLIMIT = "RequestLimitExceeded.CreateDomainLimit"
//  REQUESTLIMITEXCEEDED_REQUESTLIMITEXCEEDED = "RequestLimitExceeded.RequestLimitExceeded"
func (c *Client) CreateRecordBatch(request *CreateRecordBatchRequest) (response *CreateRecordBatchResponse, err error) {
    return c.CreateRecordBatchWithContext(context.Background(), request)
}

// CreateRecordBatch
// This API is used to bulk add records.
//
// error code that may be returned:
//  FAILEDOPERATION_DOMAINISLOCKED = "FailedOperation.DomainIsLocked"
//  FAILEDOPERATION_DOMAINISSPAM = "FailedOperation.DomainIsSpam"
//  INVALIDPARAMETER_BATCHRECORDCREATEACTIONERROR = "InvalidParameter.BatchRecordCreateActionError"
//  INVALIDPARAMETER_BATCHTASKCOUNTLIMIT = "InvalidParameter.BatchTaskCountLimit"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETER_DOMAINSEMPTY = "InvalidParameter.DomainsEmpty"
//  INVALIDPARAMETER_JOBGREATERTHANLIMIT = "InvalidParameter.JobGreaterThanLimit"
//  INVALIDPARAMETER_PARAMSILLEGAL = "InvalidParameter.ParamsIllegal"
//  INVALIDPARAMETER_RECORDVALUEINVALID = "InvalidParameter.RecordValueInvalid"
//  INVALIDPARAMETER_RECORDVALUELENGTHINVALID = "InvalidParameter.RecordValueLengthInvalid"
//  INVALIDPARAMETER_RECORDSEMPTY = "InvalidParameter.RecordsEmpty"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  OPERATIONDENIED_IPINBLACKLISTNOTALLOWED = "OperationDenied.IPInBlacklistNotAllowed"
//  REQUESTLIMITEXCEEDED_BATCHTASKLIMIT = "RequestLimitExceeded.BatchTaskLimit"
//  REQUESTLIMITEXCEEDED_CREATEDOMAINLIMIT = "RequestLimitExceeded.CreateDomainLimit"
//  REQUESTLIMITEXCEEDED_REQUESTLIMITEXCEEDED = "RequestLimitExceeded.RequestLimitExceeded"
func (c *Client) CreateRecordBatchWithContext(ctx context.Context, request *CreateRecordBatchRequest) (response *CreateRecordBatchResponse, err error) {
    if request == nil {
        request = NewCreateRecordBatchRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "dnspod", APIVersion, "CreateRecordBatch")
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateRecordBatch require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateRecordBatchResponse()
    err = c.Send(request, response)
    return
}

func NewCreateRecordGroupRequest() (request *CreateRecordGroupRequest) {
    request = &CreateRecordGroupRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("dnspod", APIVersion, "CreateRecordGroup")
    
    
    return
}

func NewCreateRecordGroupResponse() (response *CreateRecordGroupResponse) {
    response = &CreateRecordGroupResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// CreateRecordGroup
// This API is used to add a record group.
//
// error code that may be returned:
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
func (c *Client) CreateRecordGroup(request *CreateRecordGroupRequest) (response *CreateRecordGroupResponse, err error) {
    return c.CreateRecordGroupWithContext(context.Background(), request)
}

// CreateRecordGroup
// This API is used to add a record group.
//
// error code that may be returned:
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
func (c *Client) CreateRecordGroupWithContext(ctx context.Context, request *CreateRecordGroupRequest) (response *CreateRecordGroupResponse, err error) {
    if request == nil {
        request = NewCreateRecordGroupRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "dnspod", APIVersion, "CreateRecordGroup")
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateRecordGroup require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateRecordGroupResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteDomainRequest() (request *DeleteDomainRequest) {
    request = &DeleteDomainRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("dnspod", APIVersion, "DeleteDomain")
    
    
    return
}

func NewDeleteDomainResponse() (response *DeleteDomainResponse) {
    response = &DeleteDomainResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DeleteDomain
// This API is used to delete a domain.
//
// error code that may be returned:
//  FAILEDOPERATION_DOMAINISKEYDOMAIN = "FailedOperation.DomainIsKeyDomain"
//  FAILEDOPERATION_DOMAINISLOCKED = "FailedOperation.DomainIsLocked"
//  FAILEDOPERATION_DOMAINISSPAM = "FailedOperation.DomainIsSpam"
//  FAILEDOPERATION_DOMAINISVIP = "FailedOperation.DomainIsVip"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_ACCOUNTISBANNED = "InvalidParameter.AccountIsBanned"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININEFFECTORINVALIDATED = "InvalidParameter.DomainInEffectOrInvalidated"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
func (c *Client) DeleteDomain(request *DeleteDomainRequest) (response *DeleteDomainResponse, err error) {
    return c.DeleteDomainWithContext(context.Background(), request)
}

// DeleteDomain
// This API is used to delete a domain.
//
// error code that may be returned:
//  FAILEDOPERATION_DOMAINISKEYDOMAIN = "FailedOperation.DomainIsKeyDomain"
//  FAILEDOPERATION_DOMAINISLOCKED = "FailedOperation.DomainIsLocked"
//  FAILEDOPERATION_DOMAINISSPAM = "FailedOperation.DomainIsSpam"
//  FAILEDOPERATION_DOMAINISVIP = "FailedOperation.DomainIsVip"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_ACCOUNTISBANNED = "InvalidParameter.AccountIsBanned"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININEFFECTORINVALIDATED = "InvalidParameter.DomainInEffectOrInvalidated"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
func (c *Client) DeleteDomainWithContext(ctx context.Context, request *DeleteDomainRequest) (response *DeleteDomainResponse, err error) {
    if request == nil {
        request = NewDeleteDomainRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "dnspod", APIVersion, "DeleteDomain")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteDomain require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteDomainResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteDomainAliasRequest() (request *DeleteDomainAliasRequest) {
    request = &DeleteDomainAliasRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("dnspod", APIVersion, "DeleteDomainAlias")
    
    
    return
}

func NewDeleteDomainAliasResponse() (response *DeleteDomainAliasResponse) {
    response = &DeleteDomainAliasResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DeleteDomainAlias
// This API is used to delete a domain alias.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_DOMAINISLOCKED = "FailedOperation.DomainIsLocked"
//  FAILEDOPERATION_DOMAINISSPAM = "FailedOperation.DomainIsSpam"
//  FAILEDOPERATION_NOTREALNAMEDUSER = "FailedOperation.NotRealNamedUser"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_DOMAINALIASIDINVALID = "InvalidParameter.DomainAliasIdInvalid"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININEFFECTORINVALIDATED = "InvalidParameter.DomainInEffectOrInvalidated"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  MISSINGPARAMETER = "MissingParameter"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
//  RESOURCENOTFOUND = "ResourceNotFound"
func (c *Client) DeleteDomainAlias(request *DeleteDomainAliasRequest) (response *DeleteDomainAliasResponse, err error) {
    return c.DeleteDomainAliasWithContext(context.Background(), request)
}

// DeleteDomainAlias
// This API is used to delete a domain alias.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_DOMAINISLOCKED = "FailedOperation.DomainIsLocked"
//  FAILEDOPERATION_DOMAINISSPAM = "FailedOperation.DomainIsSpam"
//  FAILEDOPERATION_NOTREALNAMEDUSER = "FailedOperation.NotRealNamedUser"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_DOMAINALIASIDINVALID = "InvalidParameter.DomainAliasIdInvalid"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININEFFECTORINVALIDATED = "InvalidParameter.DomainInEffectOrInvalidated"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  MISSINGPARAMETER = "MissingParameter"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
//  RESOURCENOTFOUND = "ResourceNotFound"
func (c *Client) DeleteDomainAliasWithContext(ctx context.Context, request *DeleteDomainAliasRequest) (response *DeleteDomainAliasResponse, err error) {
    if request == nil {
        request = NewDeleteDomainAliasRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "dnspod", APIVersion, "DeleteDomainAlias")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteDomainAlias require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteDomainAliasResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteDomainBatchRequest() (request *DeleteDomainBatchRequest) {
    request = &DeleteDomainBatchRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("dnspod", APIVersion, "DeleteDomainBatch")
    
    
    return
}

func NewDeleteDomainBatchResponse() (response *DeleteDomainBatchResponse) {
    response = &DeleteDomainBatchResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DeleteDomainBatch
// This API is used to batch delete domains.
//
// error code that may be returned:
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_BATCHRECORDREMOVEACTIONERROR = "InvalidParameter.BatchRecordRemoveActionError"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETER_DOMAINSEMPTY = "InvalidParameter.DomainsEmpty"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
func (c *Client) DeleteDomainBatch(request *DeleteDomainBatchRequest) (response *DeleteDomainBatchResponse, err error) {
    return c.DeleteDomainBatchWithContext(context.Background(), request)
}

// DeleteDomainBatch
// This API is used to batch delete domains.
//
// error code that may be returned:
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_BATCHRECORDREMOVEACTIONERROR = "InvalidParameter.BatchRecordRemoveActionError"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETER_DOMAINSEMPTY = "InvalidParameter.DomainsEmpty"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
func (c *Client) DeleteDomainBatchWithContext(ctx context.Context, request *DeleteDomainBatchRequest) (response *DeleteDomainBatchResponse, err error) {
    if request == nil {
        request = NewDeleteDomainBatchRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "dnspod", APIVersion, "DeleteDomainBatch")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteDomainBatch require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteDomainBatchResponse()
    err = c.Send(request, response)
    return
}

func NewDeletePackageOrderRequest() (request *DeletePackageOrderRequest) {
    request = &DeletePackageOrderRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("dnspod", APIVersion, "DeletePackageOrder")
    
    
    return
}

func NewDeletePackageOrderResponse() (response *DeletePackageOrderResponse) {
    response = &DeletePackageOrderResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DeletePackageOrder
// This API is used to disable the paid plan on the international website.
//
// error code that may be returned:
//  FAILEDOPERATION_CANNOTCLOSEFREEDOMAIN = "FailedOperation.CanNotCloseFreeDomain"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
func (c *Client) DeletePackageOrder(request *DeletePackageOrderRequest) (response *DeletePackageOrderResponse, err error) {
    return c.DeletePackageOrderWithContext(context.Background(), request)
}

// DeletePackageOrder
// This API is used to disable the paid plan on the international website.
//
// error code that may be returned:
//  FAILEDOPERATION_CANNOTCLOSEFREEDOMAIN = "FailedOperation.CanNotCloseFreeDomain"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
func (c *Client) DeletePackageOrderWithContext(ctx context.Context, request *DeletePackageOrderRequest) (response *DeletePackageOrderResponse, err error) {
    if request == nil {
        request = NewDeletePackageOrderRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "dnspod", APIVersion, "DeletePackageOrder")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeletePackageOrder require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeletePackageOrderResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteRecordRequest() (request *DeleteRecordRequest) {
    request = &DeleteRecordRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("dnspod", APIVersion, "DeleteRecord")
    
    
    return
}

func NewDeleteRecordResponse() (response *DeleteRecordResponse) {
    response = &DeleteRecordResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DeleteRecord
// This API is used to delete a record.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_DOMAINISLOCKED = "FailedOperation.DomainIsLocked"
//  FAILEDOPERATION_DOMAINISSPAM = "FailedOperation.DomainIsSpam"
//  FAILEDOPERATION_LOGINAREANOTALLOWED = "FailedOperation.LoginAreaNotAllowed"
//  FAILEDOPERATION_LOGINFAILED = "FailedOperation.LoginFailed"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETER_DOMAINNOTALLOWEDMODIFYRECORDS = "InvalidParameter.DomainNotAllowedModifyRecords"
//  INVALIDPARAMETER_LOGINTOKENIDERROR = "InvalidParameter.LoginTokenIdError"
//  INVALIDPARAMETER_LOGINTOKENNOTEXISTS = "InvalidParameter.LoginTokenNotExists"
//  INVALIDPARAMETER_LOGINTOKENVALIDATEFAILED = "InvalidParameter.LoginTokenValidateFailed"
//  INVALIDPARAMETER_RECORDIDINVALID = "InvalidParameter.RecordIdInvalid"
//  INVALIDPARAMETER_REQUESTIPLIMITED = "InvalidParameter.RequestIpLimited"
//  INVALIDPARAMETER_UNREALNAMEUSER = "InvalidParameter.UnrealNameUser"
//  INVALIDPARAMETER_USERNOTEXISTS = "InvalidParameter.UserNotExists"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  INVALIDPARAMETERVALUE_USERIDINVALID = "InvalidParameterValue.UserIdInvalid"
//  LIMITEXCEEDED_FAILEDLOGINLIMITEXCEEDED = "LimitExceeded.FailedLoginLimitExceeded"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
//  OPERATIONDENIED_NOTADMIN = "OperationDenied.NotAdmin"
//  OPERATIONDENIED_NOTAGENT = "OperationDenied.NotAgent"
//  OPERATIONDENIED_NOTMANAGEDUSER = "OperationDenied.NotManagedUser"
//  REQUESTLIMITEXCEEDED_REQUESTLIMITEXCEEDED = "RequestLimitExceeded.RequestLimitExceeded"
func (c *Client) DeleteRecord(request *DeleteRecordRequest) (response *DeleteRecordResponse, err error) {
    return c.DeleteRecordWithContext(context.Background(), request)
}

// DeleteRecord
// This API is used to delete a record.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_DOMAINISLOCKED = "FailedOperation.DomainIsLocked"
//  FAILEDOPERATION_DOMAINISSPAM = "FailedOperation.DomainIsSpam"
//  FAILEDOPERATION_LOGINAREANOTALLOWED = "FailedOperation.LoginAreaNotAllowed"
//  FAILEDOPERATION_LOGINFAILED = "FailedOperation.LoginFailed"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETER_DOMAINNOTALLOWEDMODIFYRECORDS = "InvalidParameter.DomainNotAllowedModifyRecords"
//  INVALIDPARAMETER_LOGINTOKENIDERROR = "InvalidParameter.LoginTokenIdError"
//  INVALIDPARAMETER_LOGINTOKENNOTEXISTS = "InvalidParameter.LoginTokenNotExists"
//  INVALIDPARAMETER_LOGINTOKENVALIDATEFAILED = "InvalidParameter.LoginTokenValidateFailed"
//  INVALIDPARAMETER_RECORDIDINVALID = "InvalidParameter.RecordIdInvalid"
//  INVALIDPARAMETER_REQUESTIPLIMITED = "InvalidParameter.RequestIpLimited"
//  INVALIDPARAMETER_UNREALNAMEUSER = "InvalidParameter.UnrealNameUser"
//  INVALIDPARAMETER_USERNOTEXISTS = "InvalidParameter.UserNotExists"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  INVALIDPARAMETERVALUE_USERIDINVALID = "InvalidParameterValue.UserIdInvalid"
//  LIMITEXCEEDED_FAILEDLOGINLIMITEXCEEDED = "LimitExceeded.FailedLoginLimitExceeded"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
//  OPERATIONDENIED_NOTADMIN = "OperationDenied.NotAdmin"
//  OPERATIONDENIED_NOTAGENT = "OperationDenied.NotAgent"
//  OPERATIONDENIED_NOTMANAGEDUSER = "OperationDenied.NotManagedUser"
//  REQUESTLIMITEXCEEDED_REQUESTLIMITEXCEEDED = "RequestLimitExceeded.RequestLimitExceeded"
func (c *Client) DeleteRecordWithContext(ctx context.Context, request *DeleteRecordRequest) (response *DeleteRecordResponse, err error) {
    if request == nil {
        request = NewDeleteRecordRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "dnspod", APIVersion, "DeleteRecord")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteRecord require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteRecordResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteRecordGroupRequest() (request *DeleteRecordGroupRequest) {
    request = &DeleteRecordGroupRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("dnspod", APIVersion, "DeleteRecordGroup")
    
    
    return
}

func NewDeleteRecordGroupResponse() (response *DeleteRecordGroupResponse) {
    response = &DeleteRecordGroupResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DeleteRecordGroup
// This API is used to delete a record group.
//
// error code that may be returned:
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
func (c *Client) DeleteRecordGroup(request *DeleteRecordGroupRequest) (response *DeleteRecordGroupResponse, err error) {
    return c.DeleteRecordGroupWithContext(context.Background(), request)
}

// DeleteRecordGroup
// This API is used to delete a record group.
//
// error code that may be returned:
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
func (c *Client) DeleteRecordGroupWithContext(ctx context.Context, request *DeleteRecordGroupRequest) (response *DeleteRecordGroupResponse, err error) {
    if request == nil {
        request = NewDeleteRecordGroupRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "dnspod", APIVersion, "DeleteRecordGroup")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteRecordGroup require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteRecordGroupResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteShareDomainRequest() (request *DeleteShareDomainRequest) {
    request = &DeleteShareDomainRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("dnspod", APIVersion, "DeleteShareDomain")
    
    
    return
}

func NewDeleteShareDomainResponse() (response *DeleteShareDomainResponse) {
    response = &DeleteShareDomainResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DeleteShareDomain
// This API is used to unshare a domain.
//
// error code that may be returned:
//  FAILEDOPERATION_DOMAINISLOCKED = "FailedOperation.DomainIsLocked"
//  FAILEDOPERATION_DOMAINISSPAM = "FailedOperation.DomainIsSpam"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETER_EMAILINVALID = "InvalidParameter.EmailInvalid"
//  INVALIDPARAMETER_EMAILORQQINVALID = "InvalidParameter.EmailOrQqInvalid"
//  INVALIDPARAMETER_USERNOTEXISTS = "InvalidParameter.UserNotExists"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
func (c *Client) DeleteShareDomain(request *DeleteShareDomainRequest) (response *DeleteShareDomainResponse, err error) {
    return c.DeleteShareDomainWithContext(context.Background(), request)
}

// DeleteShareDomain
// This API is used to unshare a domain.
//
// error code that may be returned:
//  FAILEDOPERATION_DOMAINISLOCKED = "FailedOperation.DomainIsLocked"
//  FAILEDOPERATION_DOMAINISSPAM = "FailedOperation.DomainIsSpam"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETER_EMAILINVALID = "InvalidParameter.EmailInvalid"
//  INVALIDPARAMETER_EMAILORQQINVALID = "InvalidParameter.EmailOrQqInvalid"
//  INVALIDPARAMETER_USERNOTEXISTS = "InvalidParameter.UserNotExists"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
func (c *Client) DeleteShareDomainWithContext(ctx context.Context, request *DeleteShareDomainRequest) (response *DeleteShareDomainResponse, err error) {
    if request == nil {
        request = NewDeleteShareDomainRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "dnspod", APIVersion, "DeleteShareDomain")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteShareDomain require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteShareDomainResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeDomainRequest() (request *DescribeDomainRequest) {
    request = &DescribeDomainRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("dnspod", APIVersion, "DescribeDomain")
    
    
    return
}

func NewDescribeDomainResponse() (response *DescribeDomainResponse) {
    response = &DescribeDomainResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeDomain
// This API is used to get the information of a domain.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
func (c *Client) DescribeDomain(request *DescribeDomainRequest) (response *DescribeDomainResponse, err error) {
    return c.DescribeDomainWithContext(context.Background(), request)
}

// DescribeDomain
// This API is used to get the information of a domain.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
func (c *Client) DescribeDomainWithContext(ctx context.Context, request *DescribeDomainRequest) (response *DescribeDomainResponse, err error) {
    if request == nil {
        request = NewDescribeDomainRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "dnspod", APIVersion, "DescribeDomain")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeDomain require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeDomainResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeDomainAliasListRequest() (request *DescribeDomainAliasListRequest) {
    request = &DescribeDomainAliasListRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("dnspod", APIVersion, "DescribeDomainAliasList")
    
    
    return
}

func NewDescribeDomainAliasListResponse() (response *DescribeDomainAliasListResponse) {
    response = &DescribeDomainAliasListResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeDomainAliasList
// This API is used to get the list of domain aliases.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_NOTREALNAMEDUSER = "FailedOperation.NotRealNamedUser"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  MISSINGPARAMETER = "MissingParameter"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
//  RESOURCENOTFOUND = "ResourceNotFound"
//  RESOURCENOTFOUND_NODATAOFDOMAINALIAS = "ResourceNotFound.NoDataOfDomainAlias"
func (c *Client) DescribeDomainAliasList(request *DescribeDomainAliasListRequest) (response *DescribeDomainAliasListResponse, err error) {
    return c.DescribeDomainAliasListWithContext(context.Background(), request)
}

// DescribeDomainAliasList
// This API is used to get the list of domain aliases.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_NOTREALNAMEDUSER = "FailedOperation.NotRealNamedUser"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  MISSINGPARAMETER = "MissingParameter"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
//  RESOURCENOTFOUND = "ResourceNotFound"
//  RESOURCENOTFOUND_NODATAOFDOMAINALIAS = "ResourceNotFound.NoDataOfDomainAlias"
func (c *Client) DescribeDomainAliasListWithContext(ctx context.Context, request *DescribeDomainAliasListRequest) (response *DescribeDomainAliasListResponse, err error) {
    if request == nil {
        request = NewDescribeDomainAliasListRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "dnspod", APIVersion, "DescribeDomainAliasList")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeDomainAliasList require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeDomainAliasListResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeDomainGroupListRequest() (request *DescribeDomainGroupListRequest) {
    request = &DescribeDomainGroupListRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("dnspod", APIVersion, "DescribeDomainGroupList")
    
    
    return
}

func NewDescribeDomainGroupListResponse() (response *DescribeDomainGroupListResponse) {
    response = &DescribeDomainGroupListResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeDomainGroupList
// This API is used to get the list of domain groups.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_NOTREALNAMEDUSER = "FailedOperation.NotRealNamedUser"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  MISSINGPARAMETER = "MissingParameter"
//  RESOURCENOTFOUND = "ResourceNotFound"
func (c *Client) DescribeDomainGroupList(request *DescribeDomainGroupListRequest) (response *DescribeDomainGroupListResponse, err error) {
    return c.DescribeDomainGroupListWithContext(context.Background(), request)
}

// DescribeDomainGroupList
// This API is used to get the list of domain groups.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_NOTREALNAMEDUSER = "FailedOperation.NotRealNamedUser"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  MISSINGPARAMETER = "MissingParameter"
//  RESOURCENOTFOUND = "ResourceNotFound"
func (c *Client) DescribeDomainGroupListWithContext(ctx context.Context, request *DescribeDomainGroupListRequest) (response *DescribeDomainGroupListResponse, err error) {
    if request == nil {
        request = NewDescribeDomainGroupListRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "dnspod", APIVersion, "DescribeDomainGroupList")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeDomainGroupList require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeDomainGroupListResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeDomainListRequest() (request *DescribeDomainListRequest) {
    request = &DescribeDomainListRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("dnspod", APIVersion, "DescribeDomainList")
    
    
    return
}

func NewDescribeDomainListResponse() (response *DescribeDomainListResponse) {
    response = &DescribeDomainListResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeDomainList
// This API is used to get the list of domains.
//
// error code that may be returned:
//  AUTHFAILURE = "AuthFailure"
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_NOTREALNAMEDUSER = "FailedOperation.NotRealNamedUser"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_GROUPIDINVALID = "InvalidParameter.GroupIdInvalid"
//  INVALIDPARAMETER_OFFSETINVALID = "InvalidParameter.OffsetInvalid"
//  INVALIDPARAMETER_OPERATEFAILED = "InvalidParameter.OperateFailed"
//  INVALIDPARAMETER_RESULTMORETHAN500 = "InvalidParameter.ResultMoreThan500"
//  INVALIDPARAMETERVALUE_LIMITINVALID = "InvalidParameterValue.LimitInvalid"
//  OPERATIONDENIED_ACCESSDENIED = "OperationDenied.AccessDenied"
//  REQUESTLIMITEXCEEDED = "RequestLimitExceeded"
//  REQUESTLIMITEXCEEDED_REQUESTLIMITEXCEEDED = "RequestLimitExceeded.RequestLimitExceeded"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeDomainList(request *DescribeDomainListRequest) (response *DescribeDomainListResponse, err error) {
    return c.DescribeDomainListWithContext(context.Background(), request)
}

// DescribeDomainList
// This API is used to get the list of domains.
//
// error code that may be returned:
//  AUTHFAILURE = "AuthFailure"
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_NOTREALNAMEDUSER = "FailedOperation.NotRealNamedUser"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_GROUPIDINVALID = "InvalidParameter.GroupIdInvalid"
//  INVALIDPARAMETER_OFFSETINVALID = "InvalidParameter.OffsetInvalid"
//  INVALIDPARAMETER_OPERATEFAILED = "InvalidParameter.OperateFailed"
//  INVALIDPARAMETER_RESULTMORETHAN500 = "InvalidParameter.ResultMoreThan500"
//  INVALIDPARAMETERVALUE_LIMITINVALID = "InvalidParameterValue.LimitInvalid"
//  OPERATIONDENIED_ACCESSDENIED = "OperationDenied.AccessDenied"
//  REQUESTLIMITEXCEEDED = "RequestLimitExceeded"
//  REQUESTLIMITEXCEEDED_REQUESTLIMITEXCEEDED = "RequestLimitExceeded.RequestLimitExceeded"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeDomainListWithContext(ctx context.Context, request *DescribeDomainListRequest) (response *DescribeDomainListResponse, err error) {
    if request == nil {
        request = NewDescribeDomainListRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "dnspod", APIVersion, "DescribeDomainList")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeDomainList require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeDomainListResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeDomainLogListRequest() (request *DescribeDomainLogListRequest) {
    request = &DescribeDomainLogListRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("dnspod", APIVersion, "DescribeDomainLogList")
    
    
    return
}

func NewDescribeDomainLogListResponse() (response *DescribeDomainLogListResponse) {
    response = &DescribeDomainLogListResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeDomainLogList
// This API is used to get the log of a domain.
//
// error code that may be returned:
//  FAILEDOPERATION_DOMAINISSPAM = "FailedOperation.DomainIsSpam"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
func (c *Client) DescribeDomainLogList(request *DescribeDomainLogListRequest) (response *DescribeDomainLogListResponse, err error) {
    return c.DescribeDomainLogListWithContext(context.Background(), request)
}

// DescribeDomainLogList
// This API is used to get the log of a domain.
//
// error code that may be returned:
//  FAILEDOPERATION_DOMAINISSPAM = "FailedOperation.DomainIsSpam"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
func (c *Client) DescribeDomainLogListWithContext(ctx context.Context, request *DescribeDomainLogListRequest) (response *DescribeDomainLogListResponse, err error) {
    if request == nil {
        request = NewDescribeDomainLogListRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "dnspod", APIVersion, "DescribeDomainLogList")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeDomainLogList require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeDomainLogListResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeDomainPurviewRequest() (request *DescribeDomainPurviewRequest) {
    request = &DescribeDomainPurviewRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("dnspod", APIVersion, "DescribeDomainPurview")
    
    
    return
}

func NewDescribeDomainPurviewResponse() (response *DescribeDomainPurviewResponse) {
    response = &DescribeDomainPurviewResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeDomainPurview
// This API is used to get the permissions of a domain.
//
// error code that may be returned:
//  FAILEDOPERATION_LOGINAREANOTALLOWED = "FailedOperation.LoginAreaNotAllowed"
//  FAILEDOPERATION_LOGINFAILED = "FailedOperation.LoginFailed"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETER_LOGINTOKENIDERROR = "InvalidParameter.LoginTokenIdError"
//  INVALIDPARAMETER_LOGINTOKENNOTEXISTS = "InvalidParameter.LoginTokenNotExists"
//  INVALIDPARAMETER_LOGINTOKENVALIDATEFAILED = "InvalidParameter.LoginTokenValidateFailed"
//  INVALIDPARAMETER_REQUESTIPLIMITED = "InvalidParameter.RequestIpLimited"
//  INVALIDPARAMETER_UNREALNAMEUSER = "InvalidParameter.UnrealNameUser"
//  INVALIDPARAMETER_USERNOTEXISTS = "InvalidParameter.UserNotExists"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  INVALIDPARAMETERVALUE_USERIDINVALID = "InvalidParameterValue.UserIdInvalid"
//  LIMITEXCEEDED_FAILEDLOGINLIMITEXCEEDED = "LimitExceeded.FailedLoginLimitExceeded"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
//  OPERATIONDENIED_NOTAGENT = "OperationDenied.NotAgent"
//  OPERATIONDENIED_NOTMANAGEDUSER = "OperationDenied.NotManagedUser"
//  REQUESTLIMITEXCEEDED_REQUESTLIMITEXCEEDED = "RequestLimitExceeded.RequestLimitExceeded"
func (c *Client) DescribeDomainPurview(request *DescribeDomainPurviewRequest) (response *DescribeDomainPurviewResponse, err error) {
    return c.DescribeDomainPurviewWithContext(context.Background(), request)
}

// DescribeDomainPurview
// This API is used to get the permissions of a domain.
//
// error code that may be returned:
//  FAILEDOPERATION_LOGINAREANOTALLOWED = "FailedOperation.LoginAreaNotAllowed"
//  FAILEDOPERATION_LOGINFAILED = "FailedOperation.LoginFailed"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETER_LOGINTOKENIDERROR = "InvalidParameter.LoginTokenIdError"
//  INVALIDPARAMETER_LOGINTOKENNOTEXISTS = "InvalidParameter.LoginTokenNotExists"
//  INVALIDPARAMETER_LOGINTOKENVALIDATEFAILED = "InvalidParameter.LoginTokenValidateFailed"
//  INVALIDPARAMETER_REQUESTIPLIMITED = "InvalidParameter.RequestIpLimited"
//  INVALIDPARAMETER_UNREALNAMEUSER = "InvalidParameter.UnrealNameUser"
//  INVALIDPARAMETER_USERNOTEXISTS = "InvalidParameter.UserNotExists"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  INVALIDPARAMETERVALUE_USERIDINVALID = "InvalidParameterValue.UserIdInvalid"
//  LIMITEXCEEDED_FAILEDLOGINLIMITEXCEEDED = "LimitExceeded.FailedLoginLimitExceeded"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
//  OPERATIONDENIED_NOTAGENT = "OperationDenied.NotAgent"
//  OPERATIONDENIED_NOTMANAGEDUSER = "OperationDenied.NotManagedUser"
//  REQUESTLIMITEXCEEDED_REQUESTLIMITEXCEEDED = "RequestLimitExceeded.RequestLimitExceeded"
func (c *Client) DescribeDomainPurviewWithContext(ctx context.Context, request *DescribeDomainPurviewRequest) (response *DescribeDomainPurviewResponse, err error) {
    if request == nil {
        request = NewDescribeDomainPurviewRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "dnspod", APIVersion, "DescribeDomainPurview")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeDomainPurview require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeDomainPurviewResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeDomainShareInfoRequest() (request *DescribeDomainShareInfoRequest) {
    request = &DescribeDomainShareInfoRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("dnspod", APIVersion, "DescribeDomainShareInfo")
    
    
    return
}

func NewDescribeDomainShareInfoResponse() (response *DescribeDomainShareInfoResponse) {
    response = &DescribeDomainShareInfoResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeDomainShareInfo
// This API is used to get the domain sharing information.
//
// error code that may be returned:
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
func (c *Client) DescribeDomainShareInfo(request *DescribeDomainShareInfoRequest) (response *DescribeDomainShareInfoResponse, err error) {
    return c.DescribeDomainShareInfoWithContext(context.Background(), request)
}

// DescribeDomainShareInfo
// This API is used to get the domain sharing information.
//
// error code that may be returned:
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
func (c *Client) DescribeDomainShareInfoWithContext(ctx context.Context, request *DescribeDomainShareInfoRequest) (response *DescribeDomainShareInfoResponse, err error) {
    if request == nil {
        request = NewDescribeDomainShareInfoRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "dnspod", APIVersion, "DescribeDomainShareInfo")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeDomainShareInfo require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeDomainShareInfoResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeRecordRequest() (request *DescribeRecordRequest) {
    request = &DescribeRecordRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("dnspod", APIVersion, "DescribeRecord")
    
    
    return
}

func NewDescribeRecordResponse() (response *DescribeRecordResponse) {
    response = &DescribeRecordResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeRecord
// This API is used to get the information of a record.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_LOGINAREANOTALLOWED = "FailedOperation.LoginAreaNotAllowed"
//  FAILEDOPERATION_LOGINFAILED = "FailedOperation.LoginFailed"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETER_LOGINTOKENIDERROR = "InvalidParameter.LoginTokenIdError"
//  INVALIDPARAMETER_LOGINTOKENNOTEXISTS = "InvalidParameter.LoginTokenNotExists"
//  INVALIDPARAMETER_LOGINTOKENVALIDATEFAILED = "InvalidParameter.LoginTokenValidateFailed"
//  INVALIDPARAMETER_RECORDIDINVALID = "InvalidParameter.RecordIdInvalid"
//  INVALIDPARAMETER_REQUESTIPLIMITED = "InvalidParameter.RequestIpLimited"
//  INVALIDPARAMETER_UNREALNAMEUSER = "InvalidParameter.UnrealNameUser"
//  INVALIDPARAMETER_USERNOTEXISTS = "InvalidParameter.UserNotExists"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  INVALIDPARAMETERVALUE_USERIDINVALID = "InvalidParameterValue.UserIdInvalid"
//  LIMITEXCEEDED_FAILEDLOGINLIMITEXCEEDED = "LimitExceeded.FailedLoginLimitExceeded"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
//  OPERATIONDENIED_NOTADMIN = "OperationDenied.NotAdmin"
//  OPERATIONDENIED_NOTAGENT = "OperationDenied.NotAgent"
//  OPERATIONDENIED_NOTMANAGEDUSER = "OperationDenied.NotManagedUser"
//  REQUESTLIMITEXCEEDED_REQUESTLIMITEXCEEDED = "RequestLimitExceeded.RequestLimitExceeded"
func (c *Client) DescribeRecord(request *DescribeRecordRequest) (response *DescribeRecordResponse, err error) {
    return c.DescribeRecordWithContext(context.Background(), request)
}

// DescribeRecord
// This API is used to get the information of a record.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_LOGINAREANOTALLOWED = "FailedOperation.LoginAreaNotAllowed"
//  FAILEDOPERATION_LOGINFAILED = "FailedOperation.LoginFailed"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETER_LOGINTOKENIDERROR = "InvalidParameter.LoginTokenIdError"
//  INVALIDPARAMETER_LOGINTOKENNOTEXISTS = "InvalidParameter.LoginTokenNotExists"
//  INVALIDPARAMETER_LOGINTOKENVALIDATEFAILED = "InvalidParameter.LoginTokenValidateFailed"
//  INVALIDPARAMETER_RECORDIDINVALID = "InvalidParameter.RecordIdInvalid"
//  INVALIDPARAMETER_REQUESTIPLIMITED = "InvalidParameter.RequestIpLimited"
//  INVALIDPARAMETER_UNREALNAMEUSER = "InvalidParameter.UnrealNameUser"
//  INVALIDPARAMETER_USERNOTEXISTS = "InvalidParameter.UserNotExists"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  INVALIDPARAMETERVALUE_USERIDINVALID = "InvalidParameterValue.UserIdInvalid"
//  LIMITEXCEEDED_FAILEDLOGINLIMITEXCEEDED = "LimitExceeded.FailedLoginLimitExceeded"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
//  OPERATIONDENIED_NOTADMIN = "OperationDenied.NotAdmin"
//  OPERATIONDENIED_NOTAGENT = "OperationDenied.NotAgent"
//  OPERATIONDENIED_NOTMANAGEDUSER = "OperationDenied.NotManagedUser"
//  REQUESTLIMITEXCEEDED_REQUESTLIMITEXCEEDED = "RequestLimitExceeded.RequestLimitExceeded"
func (c *Client) DescribeRecordWithContext(ctx context.Context, request *DescribeRecordRequest) (response *DescribeRecordResponse, err error) {
    if request == nil {
        request = NewDescribeRecordRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "dnspod", APIVersion, "DescribeRecord")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeRecord require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeRecordResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeRecordGroupListRequest() (request *DescribeRecordGroupListRequest) {
    request = &DescribeRecordGroupListRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("dnspod", APIVersion, "DescribeRecordGroupList")
    
    
    return
}

func NewDescribeRecordGroupListResponse() (response *DescribeRecordGroupListResponse) {
    response = &DescribeRecordGroupListResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeRecordGroupList
// This API is used to query the list of DNS record groups.
//
// error code that may be returned:
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
func (c *Client) DescribeRecordGroupList(request *DescribeRecordGroupListRequest) (response *DescribeRecordGroupListResponse, err error) {
    return c.DescribeRecordGroupListWithContext(context.Background(), request)
}

// DescribeRecordGroupList
// This API is used to query the list of DNS record groups.
//
// error code that may be returned:
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
func (c *Client) DescribeRecordGroupListWithContext(ctx context.Context, request *DescribeRecordGroupListRequest) (response *DescribeRecordGroupListResponse, err error) {
    if request == nil {
        request = NewDescribeRecordGroupListRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "dnspod", APIVersion, "DescribeRecordGroupList")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeRecordGroupList require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeRecordGroupListResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeRecordLineCategoryListRequest() (request *DescribeRecordLineCategoryListRequest) {
    request = &DescribeRecordLineCategoryListRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("dnspod", APIVersion, "DescribeRecordLineCategoryList")
    
    
    return
}

func NewDescribeRecordLineCategoryListResponse() (response *DescribeRecordLineCategoryListResponse) {
    response = &DescribeRecordLineCategoryListResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeRecordLineCategoryList
// This API is used to return a line list by category.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
func (c *Client) DescribeRecordLineCategoryList(request *DescribeRecordLineCategoryListRequest) (response *DescribeRecordLineCategoryListResponse, err error) {
    return c.DescribeRecordLineCategoryListWithContext(context.Background(), request)
}

// DescribeRecordLineCategoryList
// This API is used to return a line list by category.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
func (c *Client) DescribeRecordLineCategoryListWithContext(ctx context.Context, request *DescribeRecordLineCategoryListRequest) (response *DescribeRecordLineCategoryListResponse, err error) {
    if request == nil {
        request = NewDescribeRecordLineCategoryListRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "dnspod", APIVersion, "DescribeRecordLineCategoryList")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeRecordLineCategoryList require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeRecordLineCategoryListResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeRecordLineListRequest() (request *DescribeRecordLineListRequest) {
    request = &DescribeRecordLineListRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("dnspod", APIVersion, "DescribeRecordLineList")
    
    
    return
}

func NewDescribeRecordLineListResponse() (response *DescribeRecordLineListResponse) {
    response = &DescribeRecordLineListResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeRecordLineList
// This API is used to get the split zones allowed by the domain level.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_LOGINAREANOTALLOWED = "FailedOperation.LoginAreaNotAllowed"
//  FAILEDOPERATION_LOGINFAILED = "FailedOperation.LoginFailed"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETER_LOGINTOKENIDERROR = "InvalidParameter.LoginTokenIdError"
//  INVALIDPARAMETER_LOGINTOKENNOTEXISTS = "InvalidParameter.LoginTokenNotExists"
//  INVALIDPARAMETER_LOGINTOKENVALIDATEFAILED = "InvalidParameter.LoginTokenValidateFailed"
//  INVALIDPARAMETER_REQUESTIPLIMITED = "InvalidParameter.RequestIpLimited"
//  INVALIDPARAMETER_UNREALNAMEUSER = "InvalidParameter.UnrealNameUser"
//  INVALIDPARAMETER_USERNOTEXISTS = "InvalidParameter.UserNotExists"
//  INVALIDPARAMETERVALUE_DOMAINGRADEINVALID = "InvalidParameterValue.DomainGradeInvalid"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  INVALIDPARAMETERVALUE_USERIDINVALID = "InvalidParameterValue.UserIdInvalid"
//  LIMITEXCEEDED_FAILEDLOGINLIMITEXCEEDED = "LimitExceeded.FailedLoginLimitExceeded"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
//  OPERATIONDENIED_NOTAGENT = "OperationDenied.NotAgent"
//  OPERATIONDENIED_NOTMANAGEDUSER = "OperationDenied.NotManagedUser"
//  REQUESTLIMITEXCEEDED_REQUESTLIMITEXCEEDED = "RequestLimitExceeded.RequestLimitExceeded"
func (c *Client) DescribeRecordLineList(request *DescribeRecordLineListRequest) (response *DescribeRecordLineListResponse, err error) {
    return c.DescribeRecordLineListWithContext(context.Background(), request)
}

// DescribeRecordLineList
// This API is used to get the split zones allowed by the domain level.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_LOGINAREANOTALLOWED = "FailedOperation.LoginAreaNotAllowed"
//  FAILEDOPERATION_LOGINFAILED = "FailedOperation.LoginFailed"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETER_LOGINTOKENIDERROR = "InvalidParameter.LoginTokenIdError"
//  INVALIDPARAMETER_LOGINTOKENNOTEXISTS = "InvalidParameter.LoginTokenNotExists"
//  INVALIDPARAMETER_LOGINTOKENVALIDATEFAILED = "InvalidParameter.LoginTokenValidateFailed"
//  INVALIDPARAMETER_REQUESTIPLIMITED = "InvalidParameter.RequestIpLimited"
//  INVALIDPARAMETER_UNREALNAMEUSER = "InvalidParameter.UnrealNameUser"
//  INVALIDPARAMETER_USERNOTEXISTS = "InvalidParameter.UserNotExists"
//  INVALIDPARAMETERVALUE_DOMAINGRADEINVALID = "InvalidParameterValue.DomainGradeInvalid"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  INVALIDPARAMETERVALUE_USERIDINVALID = "InvalidParameterValue.UserIdInvalid"
//  LIMITEXCEEDED_FAILEDLOGINLIMITEXCEEDED = "LimitExceeded.FailedLoginLimitExceeded"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
//  OPERATIONDENIED_NOTAGENT = "OperationDenied.NotAgent"
//  OPERATIONDENIED_NOTMANAGEDUSER = "OperationDenied.NotManagedUser"
//  REQUESTLIMITEXCEEDED_REQUESTLIMITEXCEEDED = "RequestLimitExceeded.RequestLimitExceeded"
func (c *Client) DescribeRecordLineListWithContext(ctx context.Context, request *DescribeRecordLineListRequest) (response *DescribeRecordLineListResponse, err error) {
    if request == nil {
        request = NewDescribeRecordLineListRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "dnspod", APIVersion, "DescribeRecordLineList")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeRecordLineList require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeRecordLineListResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeRecordListRequest() (request *DescribeRecordListRequest) {
    request = &DescribeRecordListRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("dnspod", APIVersion, "DescribeRecordList")
    
    
    return
}

func NewDescribeRecordListResponse() (response *DescribeRecordListResponse) {
    response = &DescribeRecordListResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeRecordList
// This API is used to get the DNS records of a domain.
//
// error code that may be returned:
//  AUTHFAILURE = "AuthFailure"
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_NOTDOMAINOWNER = "FailedOperation.NotDomainOwner"
//  FAILEDOPERATION_NOTREALNAMEDUSER = "FailedOperation.NotRealNamedUser"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETER_OPERATEFAILED = "InvalidParameter.OperateFailed"
//  INVALIDPARAMETER_PARAMINVALID = "InvalidParameter.ParamInvalid"
//  INVALIDPARAMETER_RECORDLINEINVALID = "InvalidParameter.RecordLineInvalid"
//  INVALIDPARAMETER_RECORDTYPEINVALID = "InvalidParameter.RecordTypeInvalid"
//  INVALIDPARAMETER_RESULTMORETHAN500 = "InvalidParameter.ResultMoreThan500"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  INVALIDPARAMETERVALUE_LIMITINVALID = "InvalidParameterValue.LimitInvalid"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
//  REQUESTLIMITEXCEEDED = "RequestLimitExceeded"
//  REQUESTLIMITEXCEEDED_REQUESTLIMITEXCEEDED = "RequestLimitExceeded.RequestLimitExceeded"
//  RESOURCENOTFOUND_NODATAOFRECORD = "ResourceNotFound.NoDataOfRecord"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeRecordList(request *DescribeRecordListRequest) (response *DescribeRecordListResponse, err error) {
    return c.DescribeRecordListWithContext(context.Background(), request)
}

// DescribeRecordList
// This API is used to get the DNS records of a domain.
//
// error code that may be returned:
//  AUTHFAILURE = "AuthFailure"
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_NOTDOMAINOWNER = "FailedOperation.NotDomainOwner"
//  FAILEDOPERATION_NOTREALNAMEDUSER = "FailedOperation.NotRealNamedUser"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETER_OPERATEFAILED = "InvalidParameter.OperateFailed"
//  INVALIDPARAMETER_PARAMINVALID = "InvalidParameter.ParamInvalid"
//  INVALIDPARAMETER_RECORDLINEINVALID = "InvalidParameter.RecordLineInvalid"
//  INVALIDPARAMETER_RECORDTYPEINVALID = "InvalidParameter.RecordTypeInvalid"
//  INVALIDPARAMETER_RESULTMORETHAN500 = "InvalidParameter.ResultMoreThan500"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  INVALIDPARAMETERVALUE_LIMITINVALID = "InvalidParameterValue.LimitInvalid"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
//  REQUESTLIMITEXCEEDED = "RequestLimitExceeded"
//  REQUESTLIMITEXCEEDED_REQUESTLIMITEXCEEDED = "RequestLimitExceeded.RequestLimitExceeded"
//  RESOURCENOTFOUND_NODATAOFRECORD = "ResourceNotFound.NoDataOfRecord"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeRecordListWithContext(ctx context.Context, request *DescribeRecordListRequest) (response *DescribeRecordListResponse, err error) {
    if request == nil {
        request = NewDescribeRecordListRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "dnspod", APIVersion, "DescribeRecordList")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeRecordList require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeRecordListResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeRecordTypeRequest() (request *DescribeRecordTypeRequest) {
    request = &DescribeRecordTypeRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("dnspod", APIVersion, "DescribeRecordType")
    
    
    return
}

func NewDescribeRecordTypeResponse() (response *DescribeRecordTypeResponse) {
    response = &DescribeRecordTypeResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeRecordType
// This API is used to get the record type allowed by the domain level.
//
// error code that may be returned:
//  FAILEDOPERATION_LOGINAREANOTALLOWED = "FailedOperation.LoginAreaNotAllowed"
//  FAILEDOPERATION_LOGINFAILED = "FailedOperation.LoginFailed"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_LOGINTOKENIDERROR = "InvalidParameter.LoginTokenIdError"
//  INVALIDPARAMETER_LOGINTOKENNOTEXISTS = "InvalidParameter.LoginTokenNotExists"
//  INVALIDPARAMETER_LOGINTOKENVALIDATEFAILED = "InvalidParameter.LoginTokenValidateFailed"
//  INVALIDPARAMETER_REQUESTIPLIMITED = "InvalidParameter.RequestIpLimited"
//  INVALIDPARAMETER_UNREALNAMEUSER = "InvalidParameter.UnrealNameUser"
//  INVALIDPARAMETER_USERNOTEXISTS = "InvalidParameter.UserNotExists"
//  INVALIDPARAMETERVALUE_DOMAINGRADEINVALID = "InvalidParameterValue.DomainGradeInvalid"
//  INVALIDPARAMETERVALUE_USERIDINVALID = "InvalidParameterValue.UserIdInvalid"
//  LIMITEXCEEDED_FAILEDLOGINLIMITEXCEEDED = "LimitExceeded.FailedLoginLimitExceeded"
//  MISSINGPARAMETER = "MissingParameter"
//  OPERATIONDENIED_NOTAGENT = "OperationDenied.NotAgent"
//  OPERATIONDENIED_NOTMANAGEDUSER = "OperationDenied.NotManagedUser"
//  REQUESTLIMITEXCEEDED_REQUESTLIMITEXCEEDED = "RequestLimitExceeded.RequestLimitExceeded"
//  RESOURCENOTFOUND = "ResourceNotFound"
func (c *Client) DescribeRecordType(request *DescribeRecordTypeRequest) (response *DescribeRecordTypeResponse, err error) {
    return c.DescribeRecordTypeWithContext(context.Background(), request)
}

// DescribeRecordType
// This API is used to get the record type allowed by the domain level.
//
// error code that may be returned:
//  FAILEDOPERATION_LOGINAREANOTALLOWED = "FailedOperation.LoginAreaNotAllowed"
//  FAILEDOPERATION_LOGINFAILED = "FailedOperation.LoginFailed"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_LOGINTOKENIDERROR = "InvalidParameter.LoginTokenIdError"
//  INVALIDPARAMETER_LOGINTOKENNOTEXISTS = "InvalidParameter.LoginTokenNotExists"
//  INVALIDPARAMETER_LOGINTOKENVALIDATEFAILED = "InvalidParameter.LoginTokenValidateFailed"
//  INVALIDPARAMETER_REQUESTIPLIMITED = "InvalidParameter.RequestIpLimited"
//  INVALIDPARAMETER_UNREALNAMEUSER = "InvalidParameter.UnrealNameUser"
//  INVALIDPARAMETER_USERNOTEXISTS = "InvalidParameter.UserNotExists"
//  INVALIDPARAMETERVALUE_DOMAINGRADEINVALID = "InvalidParameterValue.DomainGradeInvalid"
//  INVALIDPARAMETERVALUE_USERIDINVALID = "InvalidParameterValue.UserIdInvalid"
//  LIMITEXCEEDED_FAILEDLOGINLIMITEXCEEDED = "LimitExceeded.FailedLoginLimitExceeded"
//  MISSINGPARAMETER = "MissingParameter"
//  OPERATIONDENIED_NOTAGENT = "OperationDenied.NotAgent"
//  OPERATIONDENIED_NOTMANAGEDUSER = "OperationDenied.NotManagedUser"
//  REQUESTLIMITEXCEEDED_REQUESTLIMITEXCEEDED = "RequestLimitExceeded.RequestLimitExceeded"
//  RESOURCENOTFOUND = "ResourceNotFound"
func (c *Client) DescribeRecordTypeWithContext(ctx context.Context, request *DescribeRecordTypeRequest) (response *DescribeRecordTypeResponse, err error) {
    if request == nil {
        request = NewDescribeRecordTypeRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "dnspod", APIVersion, "DescribeRecordType")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeRecordType require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeRecordTypeResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeSubdomainAnalyticsRequest() (request *DescribeSubdomainAnalyticsRequest) {
    request = &DescribeSubdomainAnalyticsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("dnspod", APIVersion, "DescribeSubdomainAnalytics")
    
    
    return
}

func NewDescribeSubdomainAnalyticsResponse() (response *DescribeSubdomainAnalyticsResponse) {
    response = &DescribeSubdomainAnalyticsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeSubdomainAnalytics
// This API is used to collect statistics on the DNS query volume of a subdomain. It helps you understand the traffic and time period distribution and allows you to view statistics in the last three months. It is available only for domains under a paid plan.
//
// error code that may be returned:
//  FAILEDOPERATION_DOMAINISSPAM = "FailedOperation.DomainIsSpam"
//  FAILEDOPERATION_DOMAINNOTINSERVICE = "FailedOperation.DomainNotInService"
//  FAILEDOPERATION_TEMPORARYERROR = "FailedOperation.TemporaryError"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETERVALUE_DOMAINGRADEINVALID = "InvalidParameterValue.DomainGradeInvalid"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
func (c *Client) DescribeSubdomainAnalytics(request *DescribeSubdomainAnalyticsRequest) (response *DescribeSubdomainAnalyticsResponse, err error) {
    return c.DescribeSubdomainAnalyticsWithContext(context.Background(), request)
}

// DescribeSubdomainAnalytics
// This API is used to collect statistics on the DNS query volume of a subdomain. It helps you understand the traffic and time period distribution and allows you to view statistics in the last three months. It is available only for domains under a paid plan.
//
// error code that may be returned:
//  FAILEDOPERATION_DOMAINISSPAM = "FailedOperation.DomainIsSpam"
//  FAILEDOPERATION_DOMAINNOTINSERVICE = "FailedOperation.DomainNotInService"
//  FAILEDOPERATION_TEMPORARYERROR = "FailedOperation.TemporaryError"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETERVALUE_DOMAINGRADEINVALID = "InvalidParameterValue.DomainGradeInvalid"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
func (c *Client) DescribeSubdomainAnalyticsWithContext(ctx context.Context, request *DescribeSubdomainAnalyticsRequest) (response *DescribeSubdomainAnalyticsResponse, err error) {
    if request == nil {
        request = NewDescribeSubdomainAnalyticsRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "dnspod", APIVersion, "DescribeSubdomainAnalytics")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeSubdomainAnalytics require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeSubdomainAnalyticsResponse()
    err = c.Send(request, response)
    return
}

func NewModifyDomainLockRequest() (request *ModifyDomainLockRequest) {
    request = &ModifyDomainLockRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("dnspod", APIVersion, "ModifyDomainLock")
    
    
    return
}

func NewModifyDomainLockResponse() (response *ModifyDomainLockResponse) {
    response = &ModifyDomainLockResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// ModifyDomainLock
// This API is used to lock a domain.
//
// error code that may be returned:
//  FAILEDOPERATION_DOMAINISLOCKED = "FailedOperation.DomainIsLocked"
//  FAILEDOPERATION_DOMAINISSPAM = "FailedOperation.DomainIsSpam"
//  FAILEDOPERATION_LOGINAREANOTALLOWED = "FailedOperation.LoginAreaNotAllowed"
//  FAILEDOPERATION_LOGINFAILED = "FailedOperation.LoginFailed"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETER_DOMAINNOTALLOWEDLOCK = "InvalidParameter.DomainNotAllowedLock"
//  INVALIDPARAMETER_LOCKDAYSINVALID = "InvalidParameter.LockDaysInvalid"
//  INVALIDPARAMETER_LOGINTOKENIDERROR = "InvalidParameter.LoginTokenIdError"
//  INVALIDPARAMETER_LOGINTOKENNOTEXISTS = "InvalidParameter.LoginTokenNotExists"
//  INVALIDPARAMETER_LOGINTOKENVALIDATEFAILED = "InvalidParameter.LoginTokenValidateFailed"
//  INVALIDPARAMETER_REQUESTIPLIMITED = "InvalidParameter.RequestIpLimited"
//  INVALIDPARAMETER_UNREALNAMEUSER = "InvalidParameter.UnrealNameUser"
//  INVALIDPARAMETER_USERNOTEXISTS = "InvalidParameter.UserNotExists"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  INVALIDPARAMETERVALUE_USERIDINVALID = "InvalidParameterValue.UserIdInvalid"
//  LIMITEXCEEDED_FAILEDLOGINLIMITEXCEEDED = "LimitExceeded.FailedLoginLimitExceeded"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
//  OPERATIONDENIED_NOTAGENT = "OperationDenied.NotAgent"
//  OPERATIONDENIED_NOTMANAGEDUSER = "OperationDenied.NotManagedUser"
//  REQUESTLIMITEXCEEDED_REQUESTLIMITEXCEEDED = "RequestLimitExceeded.RequestLimitExceeded"
func (c *Client) ModifyDomainLock(request *ModifyDomainLockRequest) (response *ModifyDomainLockResponse, err error) {
    return c.ModifyDomainLockWithContext(context.Background(), request)
}

// ModifyDomainLock
// This API is used to lock a domain.
//
// error code that may be returned:
//  FAILEDOPERATION_DOMAINISLOCKED = "FailedOperation.DomainIsLocked"
//  FAILEDOPERATION_DOMAINISSPAM = "FailedOperation.DomainIsSpam"
//  FAILEDOPERATION_LOGINAREANOTALLOWED = "FailedOperation.LoginAreaNotAllowed"
//  FAILEDOPERATION_LOGINFAILED = "FailedOperation.LoginFailed"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETER_DOMAINNOTALLOWEDLOCK = "InvalidParameter.DomainNotAllowedLock"
//  INVALIDPARAMETER_LOCKDAYSINVALID = "InvalidParameter.LockDaysInvalid"
//  INVALIDPARAMETER_LOGINTOKENIDERROR = "InvalidParameter.LoginTokenIdError"
//  INVALIDPARAMETER_LOGINTOKENNOTEXISTS = "InvalidParameter.LoginTokenNotExists"
//  INVALIDPARAMETER_LOGINTOKENVALIDATEFAILED = "InvalidParameter.LoginTokenValidateFailed"
//  INVALIDPARAMETER_REQUESTIPLIMITED = "InvalidParameter.RequestIpLimited"
//  INVALIDPARAMETER_UNREALNAMEUSER = "InvalidParameter.UnrealNameUser"
//  INVALIDPARAMETER_USERNOTEXISTS = "InvalidParameter.UserNotExists"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  INVALIDPARAMETERVALUE_USERIDINVALID = "InvalidParameterValue.UserIdInvalid"
//  LIMITEXCEEDED_FAILEDLOGINLIMITEXCEEDED = "LimitExceeded.FailedLoginLimitExceeded"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
//  OPERATIONDENIED_NOTAGENT = "OperationDenied.NotAgent"
//  OPERATIONDENIED_NOTMANAGEDUSER = "OperationDenied.NotManagedUser"
//  REQUESTLIMITEXCEEDED_REQUESTLIMITEXCEEDED = "RequestLimitExceeded.RequestLimitExceeded"
func (c *Client) ModifyDomainLockWithContext(ctx context.Context, request *ModifyDomainLockRequest) (response *ModifyDomainLockResponse, err error) {
    if request == nil {
        request = NewModifyDomainLockRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "dnspod", APIVersion, "ModifyDomainLock")
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyDomainLock require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyDomainLockResponse()
    err = c.Send(request, response)
    return
}

func NewModifyDomainOwnerRequest() (request *ModifyDomainOwnerRequest) {
    request = &ModifyDomainOwnerRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("dnspod", APIVersion, "ModifyDomainOwner")
    
    
    return
}

func NewModifyDomainOwnerResponse() (response *ModifyDomainOwnerResponse) {
    response = &ModifyDomainOwnerResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// ModifyDomainOwner
// This API is used to transfer ownership of a domain.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_DOMAINISLOCKED = "FailedOperation.DomainIsLocked"
//  FAILEDOPERATION_DOMAINISSPAM = "FailedOperation.DomainIsSpam"
//  FAILEDOPERATION_DOMAINISVIP = "FailedOperation.DomainIsVip"
//  FAILEDOPERATION_TRANSFERTOENTERPRISEDENIED = "FailedOperation.TransferToEnterpriseDenied"
//  FAILEDOPERATION_TRANSFERTOPERSONDENIED = "FailedOperation.TransferToPersonDenied"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETER_EMAILINVALID = "InvalidParameter.EmailInvalid"
//  INVALIDPARAMETER_EMAILORQQINVALID = "InvalidParameter.EmailOrQqInvalid"
//  INVALIDPARAMETER_EMAILSAME = "InvalidParameter.EmailSame"
//  INVALIDPARAMETER_OTHERACCOUNTUNREALNAME = "InvalidParameter.OtherAccountUnrealName"
//  INVALIDPARAMETER_QCLOUDUININVALID = "InvalidParameter.QcloudUinInvalid"
//  INVALIDPARAMETER_USERAREAINVALID = "InvalidParameter.UserAreaInvalid"
//  INVALIDPARAMETER_USERNOTEXISTS = "InvalidParameter.UserNotExists"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
func (c *Client) ModifyDomainOwner(request *ModifyDomainOwnerRequest) (response *ModifyDomainOwnerResponse, err error) {
    return c.ModifyDomainOwnerWithContext(context.Background(), request)
}

// ModifyDomainOwner
// This API is used to transfer ownership of a domain.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_DOMAINISLOCKED = "FailedOperation.DomainIsLocked"
//  FAILEDOPERATION_DOMAINISSPAM = "FailedOperation.DomainIsSpam"
//  FAILEDOPERATION_DOMAINISVIP = "FailedOperation.DomainIsVip"
//  FAILEDOPERATION_TRANSFERTOENTERPRISEDENIED = "FailedOperation.TransferToEnterpriseDenied"
//  FAILEDOPERATION_TRANSFERTOPERSONDENIED = "FailedOperation.TransferToPersonDenied"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETER_EMAILINVALID = "InvalidParameter.EmailInvalid"
//  INVALIDPARAMETER_EMAILORQQINVALID = "InvalidParameter.EmailOrQqInvalid"
//  INVALIDPARAMETER_EMAILSAME = "InvalidParameter.EmailSame"
//  INVALIDPARAMETER_OTHERACCOUNTUNREALNAME = "InvalidParameter.OtherAccountUnrealName"
//  INVALIDPARAMETER_QCLOUDUININVALID = "InvalidParameter.QcloudUinInvalid"
//  INVALIDPARAMETER_USERAREAINVALID = "InvalidParameter.UserAreaInvalid"
//  INVALIDPARAMETER_USERNOTEXISTS = "InvalidParameter.UserNotExists"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
func (c *Client) ModifyDomainOwnerWithContext(ctx context.Context, request *ModifyDomainOwnerRequest) (response *ModifyDomainOwnerResponse, err error) {
    if request == nil {
        request = NewModifyDomainOwnerRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "dnspod", APIVersion, "ModifyDomainOwner")
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyDomainOwner require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyDomainOwnerResponse()
    err = c.Send(request, response)
    return
}

func NewModifyDomainRemarkRequest() (request *ModifyDomainRemarkRequest) {
    request = &ModifyDomainRemarkRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("dnspod", APIVersion, "ModifyDomainRemark")
    
    
    return
}

func NewModifyDomainRemarkResponse() (response *ModifyDomainRemarkResponse) {
    response = &ModifyDomainRemarkResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// ModifyDomainRemark
// This API is used to set the remarks of a domain.
//
// error code that may be returned:
//  FAILEDOPERATION_DOMAINISLOCKED = "FailedOperation.DomainIsLocked"
//  FAILEDOPERATION_DOMAINISSPAM = "FailedOperation.DomainIsSpam"
//  FAILEDOPERATION_LOGINAREANOTALLOWED = "FailedOperation.LoginAreaNotAllowed"
//  FAILEDOPERATION_LOGINFAILED = "FailedOperation.LoginFailed"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETER_LOGINTOKENIDERROR = "InvalidParameter.LoginTokenIdError"
//  INVALIDPARAMETER_LOGINTOKENNOTEXISTS = "InvalidParameter.LoginTokenNotExists"
//  INVALIDPARAMETER_LOGINTOKENVALIDATEFAILED = "InvalidParameter.LoginTokenValidateFailed"
//  INVALIDPARAMETER_REMARKTOOLONG = "InvalidParameter.RemarkTooLong"
//  INVALIDPARAMETER_REQUESTIPLIMITED = "InvalidParameter.RequestIpLimited"
//  INVALIDPARAMETER_UNREALNAMEUSER = "InvalidParameter.UnrealNameUser"
//  INVALIDPARAMETER_USERNOTEXISTS = "InvalidParameter.UserNotExists"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  INVALIDPARAMETERVALUE_USERIDINVALID = "InvalidParameterValue.UserIdInvalid"
//  LIMITEXCEEDED_FAILEDLOGINLIMITEXCEEDED = "LimitExceeded.FailedLoginLimitExceeded"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
//  OPERATIONDENIED_NOTAGENT = "OperationDenied.NotAgent"
//  OPERATIONDENIED_NOTMANAGEDUSER = "OperationDenied.NotManagedUser"
//  REQUESTLIMITEXCEEDED_REQUESTLIMITEXCEEDED = "RequestLimitExceeded.RequestLimitExceeded"
func (c *Client) ModifyDomainRemark(request *ModifyDomainRemarkRequest) (response *ModifyDomainRemarkResponse, err error) {
    return c.ModifyDomainRemarkWithContext(context.Background(), request)
}

// ModifyDomainRemark
// This API is used to set the remarks of a domain.
//
// error code that may be returned:
//  FAILEDOPERATION_DOMAINISLOCKED = "FailedOperation.DomainIsLocked"
//  FAILEDOPERATION_DOMAINISSPAM = "FailedOperation.DomainIsSpam"
//  FAILEDOPERATION_LOGINAREANOTALLOWED = "FailedOperation.LoginAreaNotAllowed"
//  FAILEDOPERATION_LOGINFAILED = "FailedOperation.LoginFailed"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETER_LOGINTOKENIDERROR = "InvalidParameter.LoginTokenIdError"
//  INVALIDPARAMETER_LOGINTOKENNOTEXISTS = "InvalidParameter.LoginTokenNotExists"
//  INVALIDPARAMETER_LOGINTOKENVALIDATEFAILED = "InvalidParameter.LoginTokenValidateFailed"
//  INVALIDPARAMETER_REMARKTOOLONG = "InvalidParameter.RemarkTooLong"
//  INVALIDPARAMETER_REQUESTIPLIMITED = "InvalidParameter.RequestIpLimited"
//  INVALIDPARAMETER_UNREALNAMEUSER = "InvalidParameter.UnrealNameUser"
//  INVALIDPARAMETER_USERNOTEXISTS = "InvalidParameter.UserNotExists"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  INVALIDPARAMETERVALUE_USERIDINVALID = "InvalidParameterValue.UserIdInvalid"
//  LIMITEXCEEDED_FAILEDLOGINLIMITEXCEEDED = "LimitExceeded.FailedLoginLimitExceeded"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
//  OPERATIONDENIED_NOTAGENT = "OperationDenied.NotAgent"
//  OPERATIONDENIED_NOTMANAGEDUSER = "OperationDenied.NotManagedUser"
//  REQUESTLIMITEXCEEDED_REQUESTLIMITEXCEEDED = "RequestLimitExceeded.RequestLimitExceeded"
func (c *Client) ModifyDomainRemarkWithContext(ctx context.Context, request *ModifyDomainRemarkRequest) (response *ModifyDomainRemarkResponse, err error) {
    if request == nil {
        request = NewModifyDomainRemarkRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "dnspod", APIVersion, "ModifyDomainRemark")
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyDomainRemark require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyDomainRemarkResponse()
    err = c.Send(request, response)
    return
}

func NewModifyDomainStatusRequest() (request *ModifyDomainStatusRequest) {
    request = &ModifyDomainStatusRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("dnspod", APIVersion, "ModifyDomainStatus")
    
    
    return
}

func NewModifyDomainStatusResponse() (response *ModifyDomainStatusResponse) {
    response = &ModifyDomainStatusResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// ModifyDomainStatus
// This API is used to modify the status of a domain.
//
// error code that may be returned:
//  FAILEDOPERATION_DOMAINISKEYDOMAIN = "FailedOperation.DomainIsKeyDomain"
//  FAILEDOPERATION_DOMAINISLOCKED = "FailedOperation.DomainIsLocked"
//  FAILEDOPERATION_DOMAINISSPAM = "FailedOperation.DomainIsSpam"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETER_TOOLSDOMAININVALID = "InvalidParameter.ToolsDomainInvalid"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
func (c *Client) ModifyDomainStatus(request *ModifyDomainStatusRequest) (response *ModifyDomainStatusResponse, err error) {
    return c.ModifyDomainStatusWithContext(context.Background(), request)
}

// ModifyDomainStatus
// This API is used to modify the status of a domain.
//
// error code that may be returned:
//  FAILEDOPERATION_DOMAINISKEYDOMAIN = "FailedOperation.DomainIsKeyDomain"
//  FAILEDOPERATION_DOMAINISLOCKED = "FailedOperation.DomainIsLocked"
//  FAILEDOPERATION_DOMAINISSPAM = "FailedOperation.DomainIsSpam"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETER_TOOLSDOMAININVALID = "InvalidParameter.ToolsDomainInvalid"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
func (c *Client) ModifyDomainStatusWithContext(ctx context.Context, request *ModifyDomainStatusRequest) (response *ModifyDomainStatusResponse, err error) {
    if request == nil {
        request = NewModifyDomainStatusRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "dnspod", APIVersion, "ModifyDomainStatus")
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyDomainStatus require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyDomainStatusResponse()
    err = c.Send(request, response)
    return
}

func NewModifyDomainUnlockRequest() (request *ModifyDomainUnlockRequest) {
    request = &ModifyDomainUnlockRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("dnspod", APIVersion, "ModifyDomainUnlock")
    
    
    return
}

func NewModifyDomainUnlockResponse() (response *ModifyDomainUnlockResponse) {
    response = &ModifyDomainUnlockResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// ModifyDomainUnlock
// This API is used to unlock a domain.
//
// error code that may be returned:
//  FAILEDOPERATION_DOMAINISSPAM = "FailedOperation.DomainIsSpam"
//  FAILEDOPERATION_LOGINAREANOTALLOWED = "FailedOperation.LoginAreaNotAllowed"
//  FAILEDOPERATION_LOGINFAILED = "FailedOperation.LoginFailed"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETER_DOMAINISNOTLOCKED = "InvalidParameter.DomainIsNotlocked"
//  INVALIDPARAMETER_LOGINTOKENIDERROR = "InvalidParameter.LoginTokenIdError"
//  INVALIDPARAMETER_LOGINTOKENNOTEXISTS = "InvalidParameter.LoginTokenNotExists"
//  INVALIDPARAMETER_LOGINTOKENVALIDATEFAILED = "InvalidParameter.LoginTokenValidateFailed"
//  INVALIDPARAMETER_REQUESTIPLIMITED = "InvalidParameter.RequestIpLimited"
//  INVALIDPARAMETER_UNLOCKCODEEXPIRED = "InvalidParameter.UnLockCodeExpired"
//  INVALIDPARAMETER_UNLOCKCODEINVALID = "InvalidParameter.UnLockCodeInvalid"
//  INVALIDPARAMETER_UNREALNAMEUSER = "InvalidParameter.UnrealNameUser"
//  INVALIDPARAMETER_USERNOTEXISTS = "InvalidParameter.UserNotExists"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  INVALIDPARAMETERVALUE_USERIDINVALID = "InvalidParameterValue.UserIdInvalid"
//  LIMITEXCEEDED_FAILEDLOGINLIMITEXCEEDED = "LimitExceeded.FailedLoginLimitExceeded"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
//  OPERATIONDENIED_NOTAGENT = "OperationDenied.NotAgent"
//  OPERATIONDENIED_NOTMANAGEDUSER = "OperationDenied.NotManagedUser"
//  REQUESTLIMITEXCEEDED_REQUESTLIMITEXCEEDED = "RequestLimitExceeded.RequestLimitExceeded"
func (c *Client) ModifyDomainUnlock(request *ModifyDomainUnlockRequest) (response *ModifyDomainUnlockResponse, err error) {
    return c.ModifyDomainUnlockWithContext(context.Background(), request)
}

// ModifyDomainUnlock
// This API is used to unlock a domain.
//
// error code that may be returned:
//  FAILEDOPERATION_DOMAINISSPAM = "FailedOperation.DomainIsSpam"
//  FAILEDOPERATION_LOGINAREANOTALLOWED = "FailedOperation.LoginAreaNotAllowed"
//  FAILEDOPERATION_LOGINFAILED = "FailedOperation.LoginFailed"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETER_DOMAINISNOTLOCKED = "InvalidParameter.DomainIsNotlocked"
//  INVALIDPARAMETER_LOGINTOKENIDERROR = "InvalidParameter.LoginTokenIdError"
//  INVALIDPARAMETER_LOGINTOKENNOTEXISTS = "InvalidParameter.LoginTokenNotExists"
//  INVALIDPARAMETER_LOGINTOKENVALIDATEFAILED = "InvalidParameter.LoginTokenValidateFailed"
//  INVALIDPARAMETER_REQUESTIPLIMITED = "InvalidParameter.RequestIpLimited"
//  INVALIDPARAMETER_UNLOCKCODEEXPIRED = "InvalidParameter.UnLockCodeExpired"
//  INVALIDPARAMETER_UNLOCKCODEINVALID = "InvalidParameter.UnLockCodeInvalid"
//  INVALIDPARAMETER_UNREALNAMEUSER = "InvalidParameter.UnrealNameUser"
//  INVALIDPARAMETER_USERNOTEXISTS = "InvalidParameter.UserNotExists"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  INVALIDPARAMETERVALUE_USERIDINVALID = "InvalidParameterValue.UserIdInvalid"
//  LIMITEXCEEDED_FAILEDLOGINLIMITEXCEEDED = "LimitExceeded.FailedLoginLimitExceeded"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
//  OPERATIONDENIED_NOTAGENT = "OperationDenied.NotAgent"
//  OPERATIONDENIED_NOTMANAGEDUSER = "OperationDenied.NotManagedUser"
//  REQUESTLIMITEXCEEDED_REQUESTLIMITEXCEEDED = "RequestLimitExceeded.RequestLimitExceeded"
func (c *Client) ModifyDomainUnlockWithContext(ctx context.Context, request *ModifyDomainUnlockRequest) (response *ModifyDomainUnlockResponse, err error) {
    if request == nil {
        request = NewModifyDomainUnlockRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "dnspod", APIVersion, "ModifyDomainUnlock")
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyDomainUnlock require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyDomainUnlockResponse()
    err = c.Send(request, response)
    return
}

func NewModifyRecordRequest() (request *ModifyRecordRequest) {
    request = &ModifyRecordRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("dnspod", APIVersion, "ModifyRecord")
    
    
    return
}

func NewModifyRecordResponse() (response *ModifyRecordResponse) {
    response = &ModifyRecordResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// ModifyRecord
// This API is used to modify a record.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_DOMAINISLOCKED = "FailedOperation.DomainIsLocked"
//  FAILEDOPERATION_DOMAINISSPAM = "FailedOperation.DomainIsSpam"
//  FAILEDOPERATION_LOGINAREANOTALLOWED = "FailedOperation.LoginAreaNotAllowed"
//  FAILEDOPERATION_LOGINFAILED = "FailedOperation.LoginFailed"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_ACCOUNTISBANNED = "InvalidParameter.AccountIsBanned"
//  INVALIDPARAMETER_CUSTOMMESSAGE = "InvalidParameter.CustomMessage"
//  INVALIDPARAMETER_DNSSECADDCNAMEERROR = "InvalidParameter.DnssecAddCnameError"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETER_DOMAINNOTALLOWEDMODIFYRECORDS = "InvalidParameter.DomainNotAllowedModifyRecords"
//  INVALIDPARAMETER_DOMAINNOTBEIAN = "InvalidParameter.DomainNotBeian"
//  INVALIDPARAMETER_DOMAINRECORDEXIST = "InvalidParameter.DomainRecordExist"
//  INVALIDPARAMETER_EMAILNOTVERIFIED = "InvalidParameter.EmailNotVerified"
//  INVALIDPARAMETER_INVALIDWEIGHT = "InvalidParameter.InvalidWeight"
//  INVALIDPARAMETER_LOGINTOKENIDERROR = "InvalidParameter.LoginTokenIdError"
//  INVALIDPARAMETER_LOGINTOKENNOTEXISTS = "InvalidParameter.LoginTokenNotExists"
//  INVALIDPARAMETER_LOGINTOKENVALIDATEFAILED = "InvalidParameter.LoginTokenValidateFailed"
//  INVALIDPARAMETER_MOBILENOTVERIFIED = "InvalidParameter.MobileNotVerified"
//  INVALIDPARAMETER_MXINVALID = "InvalidParameter.MxInvalid"
//  INVALIDPARAMETER_RECORDIDINVALID = "InvalidParameter.RecordIdInvalid"
//  INVALIDPARAMETER_RECORDLINEINVALID = "InvalidParameter.RecordLineInvalid"
//  INVALIDPARAMETER_RECORDTYPEINVALID = "InvalidParameter.RecordTypeInvalid"
//  INVALIDPARAMETER_RECORDVALUEINVALID = "InvalidParameter.RecordValueInvalid"
//  INVALIDPARAMETER_RECORDVALUELENGTHINVALID = "InvalidParameter.RecordValueLengthInvalid"
//  INVALIDPARAMETER_REQUESTIPLIMITED = "InvalidParameter.RequestIpLimited"
//  INVALIDPARAMETER_SUBDOMAININVALID = "InvalidParameter.SubdomainInvalid"
//  INVALIDPARAMETER_UNREALNAMEUSER = "InvalidParameter.UnrealNameUser"
//  INVALIDPARAMETER_URLVALUEILLEGAL = "InvalidParameter.UrlValueIllegal"
//  INVALIDPARAMETER_USERNOTEXISTS = "InvalidParameter.UserNotExists"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  INVALIDPARAMETERVALUE_USERIDINVALID = "InvalidParameterValue.UserIdInvalid"
//  LIMITEXCEEDED_AAAACOUNTLIMIT = "LimitExceeded.AAAACountLimit"
//  LIMITEXCEEDED_ATNSRECORDLIMIT = "LimitExceeded.AtNsRecordLimit"
//  LIMITEXCEEDED_FAILEDLOGINLIMITEXCEEDED = "LimitExceeded.FailedLoginLimitExceeded"
//  LIMITEXCEEDED_HIDDENURLEXCEEDED = "LimitExceeded.HiddenUrlExceeded"
//  LIMITEXCEEDED_NSCOUNTLIMIT = "LimitExceeded.NsCountLimit"
//  LIMITEXCEEDED_RECORDTTLLIMIT = "LimitExceeded.RecordTtlLimit"
//  LIMITEXCEEDED_SRVCOUNTLIMIT = "LimitExceeded.SrvCountLimit"
//  LIMITEXCEEDED_SUBDOMAINLEVELLIMIT = "LimitExceeded.SubdomainLevelLimit"
//  LIMITEXCEEDED_SUBDOMAINROLLLIMIT = "LimitExceeded.SubdomainRollLimit"
//  LIMITEXCEEDED_SUBDOMAINWCARDLIMIT = "LimitExceeded.SubdomainWcardLimit"
//  LIMITEXCEEDED_URLCOUNTLIMIT = "LimitExceeded.UrlCountLimit"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_IPINBLACKLISTNOTALLOWED = "OperationDenied.IPInBlacklistNotAllowed"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
//  OPERATIONDENIED_NOTADMIN = "OperationDenied.NotAdmin"
//  OPERATIONDENIED_NOTAGENT = "OperationDenied.NotAgent"
//  OPERATIONDENIED_NOTMANAGEDUSER = "OperationDenied.NotManagedUser"
//  REQUESTLIMITEXCEEDED_REQUESTLIMITEXCEEDED = "RequestLimitExceeded.RequestLimitExceeded"
func (c *Client) ModifyRecord(request *ModifyRecordRequest) (response *ModifyRecordResponse, err error) {
    return c.ModifyRecordWithContext(context.Background(), request)
}

// ModifyRecord
// This API is used to modify a record.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_DOMAINISLOCKED = "FailedOperation.DomainIsLocked"
//  FAILEDOPERATION_DOMAINISSPAM = "FailedOperation.DomainIsSpam"
//  FAILEDOPERATION_LOGINAREANOTALLOWED = "FailedOperation.LoginAreaNotAllowed"
//  FAILEDOPERATION_LOGINFAILED = "FailedOperation.LoginFailed"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_ACCOUNTISBANNED = "InvalidParameter.AccountIsBanned"
//  INVALIDPARAMETER_CUSTOMMESSAGE = "InvalidParameter.CustomMessage"
//  INVALIDPARAMETER_DNSSECADDCNAMEERROR = "InvalidParameter.DnssecAddCnameError"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETER_DOMAINNOTALLOWEDMODIFYRECORDS = "InvalidParameter.DomainNotAllowedModifyRecords"
//  INVALIDPARAMETER_DOMAINNOTBEIAN = "InvalidParameter.DomainNotBeian"
//  INVALIDPARAMETER_DOMAINRECORDEXIST = "InvalidParameter.DomainRecordExist"
//  INVALIDPARAMETER_EMAILNOTVERIFIED = "InvalidParameter.EmailNotVerified"
//  INVALIDPARAMETER_INVALIDWEIGHT = "InvalidParameter.InvalidWeight"
//  INVALIDPARAMETER_LOGINTOKENIDERROR = "InvalidParameter.LoginTokenIdError"
//  INVALIDPARAMETER_LOGINTOKENNOTEXISTS = "InvalidParameter.LoginTokenNotExists"
//  INVALIDPARAMETER_LOGINTOKENVALIDATEFAILED = "InvalidParameter.LoginTokenValidateFailed"
//  INVALIDPARAMETER_MOBILENOTVERIFIED = "InvalidParameter.MobileNotVerified"
//  INVALIDPARAMETER_MXINVALID = "InvalidParameter.MxInvalid"
//  INVALIDPARAMETER_RECORDIDINVALID = "InvalidParameter.RecordIdInvalid"
//  INVALIDPARAMETER_RECORDLINEINVALID = "InvalidParameter.RecordLineInvalid"
//  INVALIDPARAMETER_RECORDTYPEINVALID = "InvalidParameter.RecordTypeInvalid"
//  INVALIDPARAMETER_RECORDVALUEINVALID = "InvalidParameter.RecordValueInvalid"
//  INVALIDPARAMETER_RECORDVALUELENGTHINVALID = "InvalidParameter.RecordValueLengthInvalid"
//  INVALIDPARAMETER_REQUESTIPLIMITED = "InvalidParameter.RequestIpLimited"
//  INVALIDPARAMETER_SUBDOMAININVALID = "InvalidParameter.SubdomainInvalid"
//  INVALIDPARAMETER_UNREALNAMEUSER = "InvalidParameter.UnrealNameUser"
//  INVALIDPARAMETER_URLVALUEILLEGAL = "InvalidParameter.UrlValueIllegal"
//  INVALIDPARAMETER_USERNOTEXISTS = "InvalidParameter.UserNotExists"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  INVALIDPARAMETERVALUE_USERIDINVALID = "InvalidParameterValue.UserIdInvalid"
//  LIMITEXCEEDED_AAAACOUNTLIMIT = "LimitExceeded.AAAACountLimit"
//  LIMITEXCEEDED_ATNSRECORDLIMIT = "LimitExceeded.AtNsRecordLimit"
//  LIMITEXCEEDED_FAILEDLOGINLIMITEXCEEDED = "LimitExceeded.FailedLoginLimitExceeded"
//  LIMITEXCEEDED_HIDDENURLEXCEEDED = "LimitExceeded.HiddenUrlExceeded"
//  LIMITEXCEEDED_NSCOUNTLIMIT = "LimitExceeded.NsCountLimit"
//  LIMITEXCEEDED_RECORDTTLLIMIT = "LimitExceeded.RecordTtlLimit"
//  LIMITEXCEEDED_SRVCOUNTLIMIT = "LimitExceeded.SrvCountLimit"
//  LIMITEXCEEDED_SUBDOMAINLEVELLIMIT = "LimitExceeded.SubdomainLevelLimit"
//  LIMITEXCEEDED_SUBDOMAINROLLLIMIT = "LimitExceeded.SubdomainRollLimit"
//  LIMITEXCEEDED_SUBDOMAINWCARDLIMIT = "LimitExceeded.SubdomainWcardLimit"
//  LIMITEXCEEDED_URLCOUNTLIMIT = "LimitExceeded.UrlCountLimit"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_IPINBLACKLISTNOTALLOWED = "OperationDenied.IPInBlacklistNotAllowed"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
//  OPERATIONDENIED_NOTADMIN = "OperationDenied.NotAdmin"
//  OPERATIONDENIED_NOTAGENT = "OperationDenied.NotAgent"
//  OPERATIONDENIED_NOTMANAGEDUSER = "OperationDenied.NotManagedUser"
//  REQUESTLIMITEXCEEDED_REQUESTLIMITEXCEEDED = "RequestLimitExceeded.RequestLimitExceeded"
func (c *Client) ModifyRecordWithContext(ctx context.Context, request *ModifyRecordRequest) (response *ModifyRecordResponse, err error) {
    if request == nil {
        request = NewModifyRecordRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "dnspod", APIVersion, "ModifyRecord")
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyRecord require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyRecordResponse()
    err = c.Send(request, response)
    return
}

func NewModifyRecordBatchRequest() (request *ModifyRecordBatchRequest) {
    request = &ModifyRecordBatchRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("dnspod", APIVersion, "ModifyRecordBatch")
    
    
    return
}

func NewModifyRecordBatchResponse() (response *ModifyRecordBatchResponse) {
    response = &ModifyRecordBatchResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// ModifyRecordBatch
// This API is used to bulk modify records.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_DOMAINISLOCKED = "FailedOperation.DomainIsLocked"
//  FAILEDOPERATION_DOMAINISSPAM = "FailedOperation.DomainIsSpam"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_BATCHRECORDMODIFYACTIONERROR = "InvalidParameter.BatchRecordModifyActionError"
//  INVALIDPARAMETER_BATCHRECORDMODIFYACTIONINVALIDVALUE = "InvalidParameter.BatchRecordModifyActionInvalidValue"
//  INVALIDPARAMETER_BATCHRECORDREPLACEACTIONERROR = "InvalidParameter.BatchRecordReplaceActionError"
//  INVALIDPARAMETER_BATCHTASKCOUNTLIMIT = "InvalidParameter.BatchTaskCountLimit"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETER_DOMAINSEMPTY = "InvalidParameter.DomainsEmpty"
//  INVALIDPARAMETER_JOBGREATERTHANLIMIT = "InvalidParameter.JobGreaterThanLimit"
//  INVALIDPARAMETER_MXINVALID = "InvalidParameter.MxInvalid"
//  INVALIDPARAMETER_PARAMSILLEGAL = "InvalidParameter.ParamsIllegal"
//  INVALIDPARAMETER_PARAMSMISSING = "InvalidParameter.ParamsMissing"
//  INVALIDPARAMETER_RECORDTYPEINVALID = "InvalidParameter.RecordTypeInvalid"
//  INVALIDPARAMETER_RECORDVALUEINVALID = "InvalidParameter.RecordValueInvalid"
//  INVALIDPARAMETER_RECORDVALUELENGTHINVALID = "InvalidParameter.RecordValueLengthInvalid"
//  INVALIDPARAMETER_RECORDSEMPTY = "InvalidParameter.RecordsEmpty"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  OPERATIONDENIED_IPINBLACKLISTNOTALLOWED = "OperationDenied.IPInBlacklistNotAllowed"
//  REQUESTLIMITEXCEEDED_BATCHTASKLIMIT = "RequestLimitExceeded.BatchTaskLimit"
//  REQUESTLIMITEXCEEDED_CREATEDOMAINLIMIT = "RequestLimitExceeded.CreateDomainLimit"
//  REQUESTLIMITEXCEEDED_REQUESTLIMITEXCEEDED = "RequestLimitExceeded.RequestLimitExceeded"
func (c *Client) ModifyRecordBatch(request *ModifyRecordBatchRequest) (response *ModifyRecordBatchResponse, err error) {
    return c.ModifyRecordBatchWithContext(context.Background(), request)
}

// ModifyRecordBatch
// This API is used to bulk modify records.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_DOMAINISLOCKED = "FailedOperation.DomainIsLocked"
//  FAILEDOPERATION_DOMAINISSPAM = "FailedOperation.DomainIsSpam"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_BATCHRECORDMODIFYACTIONERROR = "InvalidParameter.BatchRecordModifyActionError"
//  INVALIDPARAMETER_BATCHRECORDMODIFYACTIONINVALIDVALUE = "InvalidParameter.BatchRecordModifyActionInvalidValue"
//  INVALIDPARAMETER_BATCHRECORDREPLACEACTIONERROR = "InvalidParameter.BatchRecordReplaceActionError"
//  INVALIDPARAMETER_BATCHTASKCOUNTLIMIT = "InvalidParameter.BatchTaskCountLimit"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETER_DOMAINSEMPTY = "InvalidParameter.DomainsEmpty"
//  INVALIDPARAMETER_JOBGREATERTHANLIMIT = "InvalidParameter.JobGreaterThanLimit"
//  INVALIDPARAMETER_MXINVALID = "InvalidParameter.MxInvalid"
//  INVALIDPARAMETER_PARAMSILLEGAL = "InvalidParameter.ParamsIllegal"
//  INVALIDPARAMETER_PARAMSMISSING = "InvalidParameter.ParamsMissing"
//  INVALIDPARAMETER_RECORDTYPEINVALID = "InvalidParameter.RecordTypeInvalid"
//  INVALIDPARAMETER_RECORDVALUEINVALID = "InvalidParameter.RecordValueInvalid"
//  INVALIDPARAMETER_RECORDVALUELENGTHINVALID = "InvalidParameter.RecordValueLengthInvalid"
//  INVALIDPARAMETER_RECORDSEMPTY = "InvalidParameter.RecordsEmpty"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  OPERATIONDENIED_IPINBLACKLISTNOTALLOWED = "OperationDenied.IPInBlacklistNotAllowed"
//  REQUESTLIMITEXCEEDED_BATCHTASKLIMIT = "RequestLimitExceeded.BatchTaskLimit"
//  REQUESTLIMITEXCEEDED_CREATEDOMAINLIMIT = "RequestLimitExceeded.CreateDomainLimit"
//  REQUESTLIMITEXCEEDED_REQUESTLIMITEXCEEDED = "RequestLimitExceeded.RequestLimitExceeded"
func (c *Client) ModifyRecordBatchWithContext(ctx context.Context, request *ModifyRecordBatchRequest) (response *ModifyRecordBatchResponse, err error) {
    if request == nil {
        request = NewModifyRecordBatchRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "dnspod", APIVersion, "ModifyRecordBatch")
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyRecordBatch require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyRecordBatchResponse()
    err = c.Send(request, response)
    return
}

func NewModifyRecordGroupRequest() (request *ModifyRecordGroupRequest) {
    request = &ModifyRecordGroupRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("dnspod", APIVersion, "ModifyRecordGroup")
    
    
    return
}

func NewModifyRecordGroupResponse() (response *ModifyRecordGroupResponse) {
    response = &ModifyRecordGroupResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// ModifyRecordGroup
// This API is used to modify a record group.
//
// error code that may be returned:
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
func (c *Client) ModifyRecordGroup(request *ModifyRecordGroupRequest) (response *ModifyRecordGroupResponse, err error) {
    return c.ModifyRecordGroupWithContext(context.Background(), request)
}

// ModifyRecordGroup
// This API is used to modify a record group.
//
// error code that may be returned:
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
func (c *Client) ModifyRecordGroupWithContext(ctx context.Context, request *ModifyRecordGroupRequest) (response *ModifyRecordGroupResponse, err error) {
    if request == nil {
        request = NewModifyRecordGroupRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "dnspod", APIVersion, "ModifyRecordGroup")
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyRecordGroup require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyRecordGroupResponse()
    err = c.Send(request, response)
    return
}

func NewModifyRecordRemarkRequest() (request *ModifyRecordRemarkRequest) {
    request = &ModifyRecordRemarkRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("dnspod", APIVersion, "ModifyRecordRemark")
    
    
    return
}

func NewModifyRecordRemarkResponse() (response *ModifyRecordRemarkResponse) {
    response = &ModifyRecordRemarkResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// ModifyRecordRemark
// This API is used to set the remarks of a record.
//
// error code that may be returned:
//  FAILEDOPERATION_DOMAINISLOCKED = "FailedOperation.DomainIsLocked"
//  FAILEDOPERATION_DOMAINISSPAM = "FailedOperation.DomainIsSpam"
//  FAILEDOPERATION_LOGINAREANOTALLOWED = "FailedOperation.LoginAreaNotAllowed"
//  FAILEDOPERATION_LOGINFAILED = "FailedOperation.LoginFailed"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETER_DOMAINNOTALLOWEDMODIFYRECORDS = "InvalidParameter.DomainNotAllowedModifyRecords"
//  INVALIDPARAMETER_LOGINTOKENIDERROR = "InvalidParameter.LoginTokenIdError"
//  INVALIDPARAMETER_LOGINTOKENNOTEXISTS = "InvalidParameter.LoginTokenNotExists"
//  INVALIDPARAMETER_LOGINTOKENVALIDATEFAILED = "InvalidParameter.LoginTokenValidateFailed"
//  INVALIDPARAMETER_RECORDIDINVALID = "InvalidParameter.RecordIdInvalid"
//  INVALIDPARAMETER_REMARKLENGTHEXCEEDED = "InvalidParameter.RemarkLengthExceeded"
//  INVALIDPARAMETER_REQUESTIPLIMITED = "InvalidParameter.RequestIpLimited"
//  INVALIDPARAMETER_UNREALNAMEUSER = "InvalidParameter.UnrealNameUser"
//  INVALIDPARAMETER_USERNOTEXISTS = "InvalidParameter.UserNotExists"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  INVALIDPARAMETERVALUE_USERIDINVALID = "InvalidParameterValue.UserIdInvalid"
//  LIMITEXCEEDED_FAILEDLOGINLIMITEXCEEDED = "LimitExceeded.FailedLoginLimitExceeded"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
//  OPERATIONDENIED_NOTADMIN = "OperationDenied.NotAdmin"
//  OPERATIONDENIED_NOTAGENT = "OperationDenied.NotAgent"
//  OPERATIONDENIED_NOTMANAGEDUSER = "OperationDenied.NotManagedUser"
//  REQUESTLIMITEXCEEDED_REQUESTLIMITEXCEEDED = "RequestLimitExceeded.RequestLimitExceeded"
func (c *Client) ModifyRecordRemark(request *ModifyRecordRemarkRequest) (response *ModifyRecordRemarkResponse, err error) {
    return c.ModifyRecordRemarkWithContext(context.Background(), request)
}

// ModifyRecordRemark
// This API is used to set the remarks of a record.
//
// error code that may be returned:
//  FAILEDOPERATION_DOMAINISLOCKED = "FailedOperation.DomainIsLocked"
//  FAILEDOPERATION_DOMAINISSPAM = "FailedOperation.DomainIsSpam"
//  FAILEDOPERATION_LOGINAREANOTALLOWED = "FailedOperation.LoginAreaNotAllowed"
//  FAILEDOPERATION_LOGINFAILED = "FailedOperation.LoginFailed"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETER_DOMAINNOTALLOWEDMODIFYRECORDS = "InvalidParameter.DomainNotAllowedModifyRecords"
//  INVALIDPARAMETER_LOGINTOKENIDERROR = "InvalidParameter.LoginTokenIdError"
//  INVALIDPARAMETER_LOGINTOKENNOTEXISTS = "InvalidParameter.LoginTokenNotExists"
//  INVALIDPARAMETER_LOGINTOKENVALIDATEFAILED = "InvalidParameter.LoginTokenValidateFailed"
//  INVALIDPARAMETER_RECORDIDINVALID = "InvalidParameter.RecordIdInvalid"
//  INVALIDPARAMETER_REMARKLENGTHEXCEEDED = "InvalidParameter.RemarkLengthExceeded"
//  INVALIDPARAMETER_REQUESTIPLIMITED = "InvalidParameter.RequestIpLimited"
//  INVALIDPARAMETER_UNREALNAMEUSER = "InvalidParameter.UnrealNameUser"
//  INVALIDPARAMETER_USERNOTEXISTS = "InvalidParameter.UserNotExists"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  INVALIDPARAMETERVALUE_USERIDINVALID = "InvalidParameterValue.UserIdInvalid"
//  LIMITEXCEEDED_FAILEDLOGINLIMITEXCEEDED = "LimitExceeded.FailedLoginLimitExceeded"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
//  OPERATIONDENIED_NOTADMIN = "OperationDenied.NotAdmin"
//  OPERATIONDENIED_NOTAGENT = "OperationDenied.NotAgent"
//  OPERATIONDENIED_NOTMANAGEDUSER = "OperationDenied.NotManagedUser"
//  REQUESTLIMITEXCEEDED_REQUESTLIMITEXCEEDED = "RequestLimitExceeded.RequestLimitExceeded"
func (c *Client) ModifyRecordRemarkWithContext(ctx context.Context, request *ModifyRecordRemarkRequest) (response *ModifyRecordRemarkResponse, err error) {
    if request == nil {
        request = NewModifyRecordRemarkRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "dnspod", APIVersion, "ModifyRecordRemark")
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyRecordRemark require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyRecordRemarkResponse()
    err = c.Send(request, response)
    return
}

func NewModifyRecordStatusRequest() (request *ModifyRecordStatusRequest) {
    request = &ModifyRecordStatusRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("dnspod", APIVersion, "ModifyRecordStatus")
    
    
    return
}

func NewModifyRecordStatusResponse() (response *ModifyRecordStatusResponse) {
    response = &ModifyRecordStatusResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// ModifyRecordStatus
// This API is used to modify the DNS record status.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_DOMAINISLOCKED = "FailedOperation.DomainIsLocked"
//  FAILEDOPERATION_DOMAINISSPAM = "FailedOperation.DomainIsSpam"
//  FAILEDOPERATION_LOGINAREANOTALLOWED = "FailedOperation.LoginAreaNotAllowed"
//  FAILEDOPERATION_LOGINFAILED = "FailedOperation.LoginFailed"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_CUSTOMMESSAGE = "InvalidParameter.CustomMessage"
//  INVALIDPARAMETER_DNSSECADDCNAMEERROR = "InvalidParameter.DnssecAddCnameError"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETER_DOMAINNOTALLOWEDMODIFYRECORDS = "InvalidParameter.DomainNotAllowedModifyRecords"
//  INVALIDPARAMETER_DOMAINNOTBEIAN = "InvalidParameter.DomainNotBeian"
//  INVALIDPARAMETER_DOMAINRECORDEXIST = "InvalidParameter.DomainRecordExist"
//  INVALIDPARAMETER_EMAILNOTVERIFIED = "InvalidParameter.EmailNotVerified"
//  INVALIDPARAMETER_LOGINTOKENIDERROR = "InvalidParameter.LoginTokenIdError"
//  INVALIDPARAMETER_LOGINTOKENNOTEXISTS = "InvalidParameter.LoginTokenNotExists"
//  INVALIDPARAMETER_LOGINTOKENVALIDATEFAILED = "InvalidParameter.LoginTokenValidateFailed"
//  INVALIDPARAMETER_MOBILENOTVERIFIED = "InvalidParameter.MobileNotVerified"
//  INVALIDPARAMETER_RECORDIDINVALID = "InvalidParameter.RecordIdInvalid"
//  INVALIDPARAMETER_RECORDVALUEINVALID = "InvalidParameter.RecordValueInvalid"
//  INVALIDPARAMETER_REQUESTIPLIMITED = "InvalidParameter.RequestIpLimited"
//  INVALIDPARAMETER_SUBDOMAININVALID = "InvalidParameter.SubdomainInvalid"
//  INVALIDPARAMETER_UNREALNAMEUSER = "InvalidParameter.UnrealNameUser"
//  INVALIDPARAMETER_URLVALUEILLEGAL = "InvalidParameter.UrlValueIllegal"
//  INVALIDPARAMETER_USERNOTEXISTS = "InvalidParameter.UserNotExists"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  INVALIDPARAMETERVALUE_USERIDINVALID = "InvalidParameterValue.UserIdInvalid"
//  LIMITEXCEEDED_AAAACOUNTLIMIT = "LimitExceeded.AAAACountLimit"
//  LIMITEXCEEDED_ATNSRECORDLIMIT = "LimitExceeded.AtNsRecordLimit"
//  LIMITEXCEEDED_FAILEDLOGINLIMITEXCEEDED = "LimitExceeded.FailedLoginLimitExceeded"
//  LIMITEXCEEDED_HIDDENURLEXCEEDED = "LimitExceeded.HiddenUrlExceeded"
//  LIMITEXCEEDED_NSCOUNTLIMIT = "LimitExceeded.NsCountLimit"
//  LIMITEXCEEDED_SRVCOUNTLIMIT = "LimitExceeded.SrvCountLimit"
//  LIMITEXCEEDED_SUBDOMAINLEVELLIMIT = "LimitExceeded.SubdomainLevelLimit"
//  LIMITEXCEEDED_SUBDOMAINROLLLIMIT = "LimitExceeded.SubdomainRollLimit"
//  LIMITEXCEEDED_SUBDOMAINWCARDLIMIT = "LimitExceeded.SubdomainWcardLimit"
//  LIMITEXCEEDED_URLCOUNTLIMIT = "LimitExceeded.UrlCountLimit"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
//  OPERATIONDENIED_NOTADMIN = "OperationDenied.NotAdmin"
//  OPERATIONDENIED_NOTAGENT = "OperationDenied.NotAgent"
//  OPERATIONDENIED_NOTMANAGEDUSER = "OperationDenied.NotManagedUser"
//  REQUESTLIMITEXCEEDED_REQUESTLIMITEXCEEDED = "RequestLimitExceeded.RequestLimitExceeded"
func (c *Client) ModifyRecordStatus(request *ModifyRecordStatusRequest) (response *ModifyRecordStatusResponse, err error) {
    return c.ModifyRecordStatusWithContext(context.Background(), request)
}

// ModifyRecordStatus
// This API is used to modify the DNS record status.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_DOMAINISLOCKED = "FailedOperation.DomainIsLocked"
//  FAILEDOPERATION_DOMAINISSPAM = "FailedOperation.DomainIsSpam"
//  FAILEDOPERATION_LOGINAREANOTALLOWED = "FailedOperation.LoginAreaNotAllowed"
//  FAILEDOPERATION_LOGINFAILED = "FailedOperation.LoginFailed"
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_CUSTOMMESSAGE = "InvalidParameter.CustomMessage"
//  INVALIDPARAMETER_DNSSECADDCNAMEERROR = "InvalidParameter.DnssecAddCnameError"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETER_DOMAINNOTALLOWEDMODIFYRECORDS = "InvalidParameter.DomainNotAllowedModifyRecords"
//  INVALIDPARAMETER_DOMAINNOTBEIAN = "InvalidParameter.DomainNotBeian"
//  INVALIDPARAMETER_DOMAINRECORDEXIST = "InvalidParameter.DomainRecordExist"
//  INVALIDPARAMETER_EMAILNOTVERIFIED = "InvalidParameter.EmailNotVerified"
//  INVALIDPARAMETER_LOGINTOKENIDERROR = "InvalidParameter.LoginTokenIdError"
//  INVALIDPARAMETER_LOGINTOKENNOTEXISTS = "InvalidParameter.LoginTokenNotExists"
//  INVALIDPARAMETER_LOGINTOKENVALIDATEFAILED = "InvalidParameter.LoginTokenValidateFailed"
//  INVALIDPARAMETER_MOBILENOTVERIFIED = "InvalidParameter.MobileNotVerified"
//  INVALIDPARAMETER_RECORDIDINVALID = "InvalidParameter.RecordIdInvalid"
//  INVALIDPARAMETER_RECORDVALUEINVALID = "InvalidParameter.RecordValueInvalid"
//  INVALIDPARAMETER_REQUESTIPLIMITED = "InvalidParameter.RequestIpLimited"
//  INVALIDPARAMETER_SUBDOMAININVALID = "InvalidParameter.SubdomainInvalid"
//  INVALIDPARAMETER_UNREALNAMEUSER = "InvalidParameter.UnrealNameUser"
//  INVALIDPARAMETER_URLVALUEILLEGAL = "InvalidParameter.UrlValueIllegal"
//  INVALIDPARAMETER_USERNOTEXISTS = "InvalidParameter.UserNotExists"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  INVALIDPARAMETERVALUE_USERIDINVALID = "InvalidParameterValue.UserIdInvalid"
//  LIMITEXCEEDED_AAAACOUNTLIMIT = "LimitExceeded.AAAACountLimit"
//  LIMITEXCEEDED_ATNSRECORDLIMIT = "LimitExceeded.AtNsRecordLimit"
//  LIMITEXCEEDED_FAILEDLOGINLIMITEXCEEDED = "LimitExceeded.FailedLoginLimitExceeded"
//  LIMITEXCEEDED_HIDDENURLEXCEEDED = "LimitExceeded.HiddenUrlExceeded"
//  LIMITEXCEEDED_NSCOUNTLIMIT = "LimitExceeded.NsCountLimit"
//  LIMITEXCEEDED_SRVCOUNTLIMIT = "LimitExceeded.SrvCountLimit"
//  LIMITEXCEEDED_SUBDOMAINLEVELLIMIT = "LimitExceeded.SubdomainLevelLimit"
//  LIMITEXCEEDED_SUBDOMAINROLLLIMIT = "LimitExceeded.SubdomainRollLimit"
//  LIMITEXCEEDED_SUBDOMAINWCARDLIMIT = "LimitExceeded.SubdomainWcardLimit"
//  LIMITEXCEEDED_URLCOUNTLIMIT = "LimitExceeded.UrlCountLimit"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
//  OPERATIONDENIED_NOTADMIN = "OperationDenied.NotAdmin"
//  OPERATIONDENIED_NOTAGENT = "OperationDenied.NotAgent"
//  OPERATIONDENIED_NOTMANAGEDUSER = "OperationDenied.NotManagedUser"
//  REQUESTLIMITEXCEEDED_REQUESTLIMITEXCEEDED = "RequestLimitExceeded.RequestLimitExceeded"
func (c *Client) ModifyRecordStatusWithContext(ctx context.Context, request *ModifyRecordStatusRequest) (response *ModifyRecordStatusResponse, err error) {
    if request == nil {
        request = NewModifyRecordStatusRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "dnspod", APIVersion, "ModifyRecordStatus")
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyRecordStatus require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyRecordStatusResponse()
    err = c.Send(request, response)
    return
}

func NewModifyRecordToGroupRequest() (request *ModifyRecordToGroupRequest) {
    request = &ModifyRecordToGroupRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("dnspod", APIVersion, "ModifyRecordToGroup")
    
    
    return
}

func NewModifyRecordToGroupResponse() (response *ModifyRecordToGroupResponse) {
    response = &ModifyRecordToGroupResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// ModifyRecordToGroup
// This API is used to add a record to a group.
//
// error code that may be returned:
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
func (c *Client) ModifyRecordToGroup(request *ModifyRecordToGroupRequest) (response *ModifyRecordToGroupResponse, err error) {
    return c.ModifyRecordToGroupWithContext(context.Background(), request)
}

// ModifyRecordToGroup
// This API is used to add a record to a group.
//
// error code that may be returned:
//  FAILEDOPERATION_UNKNOWERROR = "FailedOperation.UnknowError"
//  INVALIDPARAMETER_DOMAINIDINVALID = "InvalidParameter.DomainIdInvalid"
//  INVALIDPARAMETER_DOMAININVALID = "InvalidParameter.DomainInvalid"
//  INVALIDPARAMETER_DOMAINISALIASER = "InvalidParameter.DomainIsAliaser"
//  INVALIDPARAMETERVALUE_DOMAINNOTEXISTS = "InvalidParameterValue.DomainNotExists"
//  OPERATIONDENIED_DOMAINOWNERALLOWEDONLY = "OperationDenied.DomainOwnerAllowedOnly"
//  OPERATIONDENIED_NOPERMISSIONTOOPERATEDOMAIN = "OperationDenied.NoPermissionToOperateDomain"
func (c *Client) ModifyRecordToGroupWithContext(ctx context.Context, request *ModifyRecordToGroupRequest) (response *ModifyRecordToGroupResponse, err error) {
    if request == nil {
        request = NewModifyRecordToGroupRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "dnspod", APIVersion, "ModifyRecordToGroup")
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyRecordToGroup require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyRecordToGroupResponse()
    err = c.Send(request, response)
    return
}
