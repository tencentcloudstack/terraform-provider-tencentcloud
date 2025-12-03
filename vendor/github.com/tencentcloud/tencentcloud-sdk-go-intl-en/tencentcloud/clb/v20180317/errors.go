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

package v20180317

const (
	// error codes for specific actions

	// CAM signature/authentication error
	AUTHFAILURE = "AuthFailure"

	// Token error.
	AUTHFAILURE_TOKENFAILURE = "AuthFailure.TokenFailure"

	// DryRun operation, which means the DryRun parameter is passed in yet the request will still be successful.
	DRYRUNOPERATION = "DryRunOperation"

	// Operation failed.
	FAILEDOPERATION = "FailedOperation"

	// Instance deletion is determined to be a high-risk operation through frequency verification. Please check the business or try again later.
	FAILEDOPERATION_FREQUENCYCHECKRISK = "FailedOperation.FrequencyCheckRisk"

	// Exceptional CLB instance status
	FAILEDOPERATION_INVALIDLBSTATUS = "FailedOperation.InvalidLBStatus"

	// Instances without listeners cannot be stopped.
	FAILEDOPERATION_NOLISTENERINLB = "FailedOperation.NoListenerInLB"

	// The specified resource is being cloned. please try again later.
	FAILEDOPERATION_RESOURCEINCLONING = "FailedOperation.ResourceInCloning"

	// The specified resource is currently being operated. Please try again later.
	FAILEDOPERATION_RESOURCEINOPERATING = "FailedOperation.ResourceInOperating"

	// Instance deletion is determined to be a high-risk operation through rule quantity verification. To forcibly delete the instance, set the forced verification parameter ForceDelete to true.
	FAILEDOPERATION_TARGETNUMCHECKRISK = "FailedOperation.TargetNumCheckRisk"

	// Instance deletion is determined to be a high-risk operation through traffic verification. To forcibly delete the instance, set the forced verification parameter ForceDelete to true.
	FAILEDOPERATION_TRAFFICCHECKRISK = "FailedOperation.TrafficCheckRisk"

	// Internal error.
	INTERNALERROR = "InternalError"

	// Parameter error.
	INVALIDPARAMETER = "InvalidParameter"

	// To ensure no resource leakage and maintain the ID idempotence of created resources, ClientToken is used to create resources. If the order process has ended and shipment failed, or the order process has not been updated for a long time, a message will indicate that the current ClientToken has timed out.
	INVALIDPARAMETER_CLIENTTOKENLIMITEXCEEDED = "InvalidParameter.ClientTokenLimitExceeded"

	// Wrong parameter format.
	INVALIDPARAMETER_FORMATERROR = "InvalidParameter.FormatError"

	// Failed to query the parameter
	INVALIDPARAMETER_INVALIDFILTER = "InvalidParameter.InvalidFilter"

	// Wrong CLB instance ID.
	INVALIDPARAMETER_LBIDNOTFOUND = "InvalidParameter.LBIdNotFound"

	// Wrong listener ID.
	INVALIDPARAMETER_LISTENERIDNOTFOUND = "InvalidParameter.ListenerIdNotFound"

	// Unable to find eligible forwarding rules.
	INVALIDPARAMETER_LOCATIONNOTFOUND = "InvalidParameter.LocationNotFound"

	// Listener port checks failed due to port conflicts or other reasons.
	INVALIDPARAMETER_PORTCHECKFAILED = "InvalidParameter.PortCheckFailed"

	// Listener protocol checks failed because the protocol used is incompatible with the corresponding operation.
	INVALIDPARAMETER_PROTOCOLCHECKFAILED = "InvalidParameter.ProtocolCheckFailed"

	// Invalid region.
	INVALIDPARAMETER_REGIONNOTFOUND = "InvalidParameter.RegionNotFound"

	// The forwarding rule has already been bound to a redirection relationship.
	INVALIDPARAMETER_REWRITEALREADYEXIST = "InvalidParameter.RewriteAlreadyExist"

	// Some redirection rules do not exist.
	INVALIDPARAMETER_SOMEREWRITENOTFOUND = "InvalidParameter.SomeRewriteNotFound"

	// Incorrect parameter value.
	INVALIDPARAMETERVALUE = "InvalidParameterValue"

	// Duplicate parameter value.
	INVALIDPARAMETERVALUE_DUPLICATE = "InvalidParameterValue.Duplicate"

	// Incorrect `Filter` parameter.
	INVALIDPARAMETERVALUE_INVALIDFILTER = "InvalidParameterValue.InvalidFilter"

	// Wrong parameter length.
	INVALIDPARAMETERVALUE_LENGTH = "InvalidParameterValue.Length"

	// Wrong parameter value range.
	INVALIDPARAMETERVALUE_RANGE = "InvalidParameterValue.Range"

	// Quota exceeded.
	LIMITEXCEEDED = "LimitExceeded"

	// Missing parameter.
	MISSINGPARAMETER = "MissingParameter"

	// Operation denied.
	OPERATIONDENIED = "OperationDenied"

	// The number of requests exceeds the frequency limit.
	REQUESTLIMITEXCEEDED = "RequestLimitExceeded"

	// The resource is occupied.
	RESOURCEINUSE = "ResourceInUse"

	// Insufficient resources.
	RESOURCEINSUFFICIENT = "ResourceInsufficient"

	// Resources do not exist.
	RESOURCENOTFOUND = "ResourceNotFound"

	// The resources have been sold out.
	RESOURCESSOLDOUT = "ResourcesSoldOut"

	// Unauthorized operation.
	UNAUTHORIZEDOPERATION = "UnauthorizedOperation"

	// Unsupported operation.
	UNSUPPORTEDOPERATION = "UnsupportedOperation"
)
