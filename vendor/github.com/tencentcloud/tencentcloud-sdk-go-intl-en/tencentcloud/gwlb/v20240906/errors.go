// Copyright (c) 2017-2018 THL A29 Limited, a Tencent company. All Rights Reserved.
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

package v20240906

const (
	// error codes for specific actions

	// Operation failed.
	FAILEDOPERATION = "FailedOperation"

	// FailedOperation.ResourceInOperating
	FAILEDOPERATION_RESOURCEINOPERATING = "FailedOperation.ResourceInOperating"

	// Internal error.
	INTERNALERROR = "InternalError"

	// Parameter error.
	INVALIDPARAMETER = "InvalidParameter"

	// InvalidParameter.FormatError
	INVALIDPARAMETER_FORMATERROR = "InvalidParameter.FormatError"

	// InvalidParameter.RegionNotFound
	INVALIDPARAMETER_REGIONNOTFOUND = "InvalidParameter.RegionNotFound"

	// Invalid parameter value.
	INVALIDPARAMETERVALUE = "InvalidParameterValue"

	// The same value exists.
	INVALIDPARAMETERVALUE_DUPLICATE = "InvalidParameterValue.Duplicate"

	// InvalidParameterValue.Length
	INVALIDPARAMETERVALUE_LENGTH = "InvalidParameterValue.Length"

	// The quota limit is exceeded.
	LIMITEXCEEDED = "LimitExceeded"

	// Unauthorized operation.
	UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
)
