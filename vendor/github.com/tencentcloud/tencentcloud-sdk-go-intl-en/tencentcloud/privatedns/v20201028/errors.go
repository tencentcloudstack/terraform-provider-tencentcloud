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

const (
	// error codes for specific actions

	// Error with CAM signature/authentication.
	AUTHFAILURE = "AuthFailure"

	// Token verification failed.
	AUTHFAILURE_TOKENFAILURE = "AuthFailure.TokenFailure"

	// `DryRun` Operation. It means that the request would have succeeded, but the `DryRun` parameter was used.
	DRYRUNOPERATION = "DryRunOperation"

	// Operation failed.
	FAILEDOPERATION = "FailedOperation"

	// Failed to associate the private domain with the VPC.
	FAILEDOPERATION_BINDZONEVPCFAILED = "FailedOperation.BindZoneVpcFailed"

	// Failed to create the record.
	FAILEDOPERATION_CREATERECORDFAILED = "FailedOperation.CreateRecordFailed"

	// An exception occurs when you create an endpoint.
	FAILEDOPERATION_CREATEVPCENDPOINTERROR = "FailedOperation.CreateVpcEndPointError"

	// Failed to create an endpoint.
	FAILEDOPERATION_CREATEVPCENDPOINTFAILED = "FailedOperation.CreateVpcEndPointFailed"

	// Failed to create the private domain.
	FAILEDOPERATION_CREATEZONEFAILED = "FailedOperation.CreateZoneFailed"

	// Data exception.
	FAILEDOPERATION_DATAERROR = "FailedOperation.DataError"

	// The private domain is currently associated with a VPC. Please disassociate the VPC first before clearing its records.
	FAILEDOPERATION_DELETELASTBINDVPCRECORDFAILED = "FailedOperation.DeleteLastBindVpcRecordFailed"

	// Record deletion failed.
	FAILEDOPERATION_DELETERECORDFAILED = "FailedOperation.DeleteRecordFailed"

	// Failed to delete the endpoint.
	FAILEDOPERATION_DELETEVPCENDPOINTFAILED = "FailedOperation.DeleteVpcEndPointFailed"

	// Failed to create an endpoint service.
	FAILEDOPERATION_ENDPOINTSERVICECREATEFAILED = "FailedOperation.EndPointServiceCreateFailed"

	// Failed to delete the endpoint service.
	FAILEDOPERATION_ENDPOINTSERVICEDELETEFAILED = "FailedOperation.EndPointServiceDeleteFailed"

	// Failed to obtain the endpoint service allowlist.
	FAILEDOPERATION_ENDPOINTSERVICEWHITELISTFAILED = "FailedOperation.EndPointServiceWhiteListFailed"

	// Failed to obtain the temporary credential.
	FAILEDOPERATION_GETTMPCREDFAILED = "FailedOperation.GetTmpCredFailed"

	// Failed to modify the record.
	FAILEDOPERATION_MODIFYRECORDFAILED = "FailedOperation.ModifyRecordFailed"

	// Failed to modify the private domain.
	FAILEDOPERATION_MODIFYZONEFAILED = "FailedOperation.ModifyZoneFailed"

	// Record update failed.
	FAILEDOPERATION_UPDATERECORDFAILED = "FailedOperation.UpdateRecordFailed"

	// Internal error.
	INTERNALERROR = "InternalError"

	// Undefined error.
	INTERNALERROR_UNDEFIENDERROR = "InternalError.UndefiendError"

	// Incorrect parameter.
	INVALIDPARAMETER = "InvalidParameter"

	// A bound account already exists.
	INVALIDPARAMETER_ACCOUNTEXIST = "InvalidParameter.AccountExist"

	// The endpoint has been bound to a forwarding rule.
	INVALIDPARAMETER_ENDPOINTBINDFORWARDRULE = "InvalidParameter.EndPointBindForwardRule"

	// The endpoint does not exist.
	INVALIDPARAMETER_ENDPOINTNOTEXISTS = "InvalidParameter.EndPointNotExists"

	// The endpoint service does not exist.
	INVALIDPARAMETER_ENDPOINTSERVICENOTEXIST = "InvalidParameter.EndPointServiceNotExist"

	// The forwarding rule does not exist.
	INVALIDPARAMETER_FORWARDRULENOTEXIST = "InvalidParameter.ForwardRuleNotExist"

	// The private domain has already been bound with the forwarding rule.
	INVALIDPARAMETER_FORWARDRULEZONEREPEATBIND = "InvalidParameter.ForwardRuleZoneRepeatBind"

	// Invalid CIDR.
	INVALIDPARAMETER_ILLEGALCIDR = "InvalidParameter.IllegalCidr"

	// Incorrect domain name.
	INVALIDPARAMETER_ILLEGALDOMAIN = "InvalidParameter.IllegalDomain"

	// Incorrect top-level domain name.
	INVALIDPARAMETER_ILLEGALDOMAINTLD = "InvalidParameter.IllegalDomainTld"

	// Invalid PTR record.
	INVALIDPARAMETER_ILLEGALPTRRECORD = "InvalidParameter.IllegalPTRRecord"

	// Invalid record.
	INVALIDPARAMETER_ILLEGALRECORD = "InvalidParameter.IllegalRecord"

	// Invalid record value.
	INVALIDPARAMETER_ILLEGALRECORDVALUE = "InvalidParameter.IllegalRecordValue"

	// Invalid VPC.
	INVALIDPARAMETER_ILLEGALVPCINFO = "InvalidParameter.IllegalVpcInfo"

	// The MX value must be a multiple of 5 between 5 and 50.
	INVALIDPARAMETER_INVALIDMX = "InvalidParameter.InvalidMX"

	// This record type does not support MX priority settings.
	INVALIDPARAMETER_MXNOTSUPPORTED = "InvalidParameter.MXNotSupported"

	// The number of round-robin DNS AAAA records exceeds 50.
	INVALIDPARAMETER_RECORDAAAACOUNTEXCEED = "InvalidParameter.RecordAAAACountExceed"

	// The number of round-robin DNS A records exceeds 50.
	INVALIDPARAMETER_RECORDACOUNTEXCEED = "InvalidParameter.RecordACountExceed"

	// The number of round-robin DNS CNAME records exceeds 50.
	INVALIDPARAMETER_RECORDCNAMECOUNTEXCEED = "InvalidParameter.RecordCNAMECountExceed"

	// Records conflict.
	INVALIDPARAMETER_RECORDCONFLICT = "InvalidParameter.RecordConflict"

	// The number of records exceeds the limit.
	INVALIDPARAMETER_RECORDCOUNTEXCEED = "InvalidParameter.RecordCountExceed"

	// The record already exists.
	INVALIDPARAMETER_RECORDEXIST = "InvalidParameter.RecordExist"

	// The number of record levels exceeds the limit.
	INVALIDPARAMETER_RECORDLEVELEXCEED = "InvalidParameter.RecordLevelExceed"

	// The number of round-robin DNS MX records exceeds 50.
	INVALIDPARAMETER_RECORDMXCOUNTEXCEED = "InvalidParameter.RecordMXCountExceed"

	// The record does not exist.
	INVALIDPARAMETER_RECORDNOTEXIST = "InvalidParameter.RecordNotExist"

	// The number of round-robin DNS records exceeds the limit.
	INVALIDPARAMETER_RECORDROLLLIMITCOUNTEXCEED = "InvalidParameter.RecordRolllimitCountExceed"

	// The number of round-robin DNS TXT records exceeds 10.
	INVALIDPARAMETER_RECORDTXTCOUNTEXCEED = "InvalidParameter.RecordTXTCountExceed"

	// The current record type does not support weight.
	INVALIDPARAMETER_RECORDUNSUPPORTWEIGHT = "InvalidParameter.RecordUnsupportWeight"

	// The VPC has been bound to another domain.
	INVALIDPARAMETER_VPCBINDED = "InvalidParameter.VpcBinded"

	// The VPC has been associated with the same primary domain.
	INVALIDPARAMETER_VPCBINDEDMAINDOMAIN = "InvalidParameter.VpcBindedMainDomain"

	// The number of PTR records associated with the VPC has reached the limit.
	INVALIDPARAMETER_VPCPTRZONEBINDEXCEED = "InvalidParameter.VpcPtrZoneBindExceed"

	// The domain does not exist.
	INVALIDPARAMETER_ZONENOTEXISTS = "InvalidParameter.ZoneNotExists"

	// Incorrect parameter value.
	INVALIDPARAMETERVALUE = "InvalidParameterValue"

	// The TTL value must be in the range of 1–86400.
	INVALIDPARAMETERVALUE_ILLEGALTTLVALUE = "InvalidParameterValue.IllegalTTLValue"

	// This is a reserved domain name and cannot be created.
	INVALIDPARAMETERVALUE_RESERVEDDOMAIN = "InvalidParameterValue.ReservedDomain"

	// No access permission because the UIN is not added to the allowlist.
	INVALIDPARAMETERVALUE_UINNOTINWHITELIST = "InvalidParameterValue.UinNotInWhiteList"

	// The quota limit is exceeded.
	LIMITEXCEEDED = "LimitExceeded"

	// Exceeded the custom TLD quota.
	LIMITEXCEEDED_TLDOUTOFLIMIT = "LimitExceeded.TldOutOfLimit"

	// The number of private domains using the custom TLD exceeds the total quota.
	LIMITEXCEEDED_TLDOUTOFRANGE = "LimitExceeded.TldOutOfRange"

	// Missing parameter.
	MISSINGPARAMETER = "MissingParameter"

	// The operation was denied.
	OPERATIONDENIED = "OperationDenied"

	// The number of requests exceeds the frequency limit.
	REQUESTLIMITEXCEEDED = "RequestLimitExceeded"

	// The resource is in use.
	RESOURCEINUSE = "ResourceInUse"

	// Insufficient resource.
	RESOURCEINSUFFICIENT = "ResourceInsufficient"

	// Insufficient balance.
	RESOURCEINSUFFICIENT_BALANCE = "ResourceInsufficient.Balance"

	// The resource does not exist.
	RESOURCENOTFOUND = "ResourceNotFound"

	// The Private DNS service is not activated.
	RESOURCENOTFOUND_SERVICENOTSUBSCRIBED = "ResourceNotFound.ServiceNotSubscribed"

	// The resource is unavailable.
	RESOURCEUNAVAILABLE = "ResourceUnavailable"

	// The platform account is not added to the allowlist of the current endpoint service.
	RESOURCEUNAVAILABLE_SERVICEWHITELISTNOTADDED = "ResourceUnavailable.ServiceWhiteListNotAdded"

	// The TLD value-added package expired.
	RESOURCEUNAVAILABLE_TLDPACKAGEEXPIRED = "ResourceUnavailable.TldPackageExpired"

	// The resources have been sold out.
	RESOURCESSOLDOUT = "ResourcesSoldOut"

	// Unauthorized operation.
	UNAUTHORIZEDOPERATION = "UnauthorizedOperation"

	// Role not authorized.
	UNAUTHORIZEDOPERATION_ROLEUNAUTHORIZED = "UnauthorizedOperation.RoleUnAuthorized"

	// Unverified user.
	UNAUTHORIZEDOPERATION_UNAUTHORIZEDACCOUNT = "UnauthorizedOperation.UnauthorizedAccount"

	// Unknown parameter.
	UNKNOWNPARAMETER = "UnknownParameter"

	// Unsupported operation.
	UNSUPPORTEDOPERATION = "UnsupportedOperation"

	// Account not bound.
	UNSUPPORTEDOPERATION_ACCOUNTNOTBOUND = "UnsupportedOperation.AccountNotBound"

	// The API call frequency exceeds the limit.
	UNSUPPORTEDOPERATION_FREQUENCYLIMIT = "UnsupportedOperation.FrequencyLimit"

	// Subdomain recursive DNS is not supported.
	UNSUPPORTEDOPERATION_NOTSUPPORTDNSFORWARD = "UnsupportedOperation.NotSupportDnsForward"
)
