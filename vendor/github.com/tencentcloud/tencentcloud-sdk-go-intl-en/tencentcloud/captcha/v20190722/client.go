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

package v20190722

import (
    "context"
    "errors"
    "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
    tchttp "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common/http"
    "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common/profile"
)

const APIVersion = "2019-07-22"

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


func NewCreateCaptchaInfoInternationalRequest() (request *CreateCaptchaInfoInternationalRequest) {
    request = &CreateCaptchaInfoInternationalRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("captcha", APIVersion, "CreateCaptchaInfoInternational")
    
    
    return
}

func NewCreateCaptchaInfoInternationalResponse() (response *CreateCaptchaInfoInternationalResponse) {
    response = &CreateCaptchaInfoInternationalResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// CreateCaptchaInfoInternational
// Create a captcha: You can create multiple Captcha based on different business needs. Each verification has different client types and security policies. The limit for new Captcha is 50.
//
// error code that may be returned:
//  AUTHFAILURE_UNAUTHORIZEDOPERATION = "AuthFailure.UnauthorizedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  UNAUTHORIZEDOPERATION_ERRAUTH = "UnauthorizedOperation.ErrAuth"
//  UNAUTHORIZEDOPERATION_UNAUTHORIZED = "UnauthorizedOperation.Unauthorized"
func (c *Client) CreateCaptchaInfoInternational(request *CreateCaptchaInfoInternationalRequest) (response *CreateCaptchaInfoInternationalResponse, err error) {
    return c.CreateCaptchaInfoInternationalWithContext(context.Background(), request)
}

// CreateCaptchaInfoInternational
// Create a captcha: You can create multiple Captcha based on different business needs. Each verification has different client types and security policies. The limit for new Captcha is 50.
//
// error code that may be returned:
//  AUTHFAILURE_UNAUTHORIZEDOPERATION = "AuthFailure.UnauthorizedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  UNAUTHORIZEDOPERATION_ERRAUTH = "UnauthorizedOperation.ErrAuth"
//  UNAUTHORIZEDOPERATION_UNAUTHORIZED = "UnauthorizedOperation.Unauthorized"
func (c *Client) CreateCaptchaInfoInternationalWithContext(ctx context.Context, request *CreateCaptchaInfoInternationalRequest) (response *CreateCaptchaInfoInternationalResponse, err error) {
    if request == nil {
        request = NewCreateCaptchaInfoInternationalRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "captcha", APIVersion, "CreateCaptchaInfoInternational")
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateCaptchaInfoInternational require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateCaptchaInfoInternationalResponse()
    err = c.Send(request, response)
    return
}

func NewCreateIpWhiteListInternationalRequest() (request *CreateIpWhiteListInternationalRequest) {
    request = &CreateIpWhiteListInternationalRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("captcha", APIVersion, "CreateIpWhiteListInternational")
    
    
    return
}

func NewCreateIpWhiteListInternationalResponse() (response *CreateIpWhiteListInternationalResponse) {
    response = &CreateIpWhiteListInternationalResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// CreateIpWhiteListInternational
// Create an IP allowlist: You can create an IP allowlist based on different business needs.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  UNAUTHORIZEDOPERATION_ERRAUTH = "UnauthorizedOperation.ErrAuth"
//  UNAUTHORIZEDOPERATION_UNAUTHORIZED = "UnauthorizedOperation.Unauthorized"
func (c *Client) CreateIpWhiteListInternational(request *CreateIpWhiteListInternationalRequest) (response *CreateIpWhiteListInternationalResponse, err error) {
    return c.CreateIpWhiteListInternationalWithContext(context.Background(), request)
}

// CreateIpWhiteListInternational
// Create an IP allowlist: You can create an IP allowlist based on different business needs.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  UNAUTHORIZEDOPERATION_ERRAUTH = "UnauthorizedOperation.ErrAuth"
//  UNAUTHORIZEDOPERATION_UNAUTHORIZED = "UnauthorizedOperation.Unauthorized"
func (c *Client) CreateIpWhiteListInternationalWithContext(ctx context.Context, request *CreateIpWhiteListInternationalRequest) (response *CreateIpWhiteListInternationalResponse, err error) {
    if request == nil {
        request = NewCreateIpWhiteListInternationalRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "captcha", APIVersion, "CreateIpWhiteListInternational")
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateIpWhiteListInternational require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateIpWhiteListInternationalResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteIpWhiteListInternationalRequest() (request *DeleteIpWhiteListInternationalRequest) {
    request = &DeleteIpWhiteListInternationalRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("captcha", APIVersion, "DeleteIpWhiteListInternational")
    
    
    return
}

func NewDeleteIpWhiteListInternationalResponse() (response *DeleteIpWhiteListInternationalResponse) {
    response = &DeleteIpWhiteListInternationalResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DeleteIpWhiteListInternational
// Delete an IP allowlist: You can delete an IP allowlist based on different business needs.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  UNAUTHORIZEDOPERATION_ERRAUTH = "UnauthorizedOperation.ErrAuth"
//  UNAUTHORIZEDOPERATION_UNAUTHORIZED = "UnauthorizedOperation.Unauthorized"
func (c *Client) DeleteIpWhiteListInternational(request *DeleteIpWhiteListInternationalRequest) (response *DeleteIpWhiteListInternationalResponse, err error) {
    return c.DeleteIpWhiteListInternationalWithContext(context.Background(), request)
}

// DeleteIpWhiteListInternational
// Delete an IP allowlist: You can delete an IP allowlist based on different business needs.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  UNAUTHORIZEDOPERATION_ERRAUTH = "UnauthorizedOperation.ErrAuth"
//  UNAUTHORIZEDOPERATION_UNAUTHORIZED = "UnauthorizedOperation.Unauthorized"
func (c *Client) DeleteIpWhiteListInternationalWithContext(ctx context.Context, request *DeleteIpWhiteListInternationalRequest) (response *DeleteIpWhiteListInternationalResponse, err error) {
    if request == nil {
        request = NewDeleteIpWhiteListInternationalRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "captcha", APIVersion, "DeleteIpWhiteListInternational")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteIpWhiteListInternational require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteIpWhiteListInternationalResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeCaptchaInfoListInternationalRequest() (request *DescribeCaptchaInfoListInternationalRequest) {
    request = &DescribeCaptchaInfoListInternationalRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("captcha", APIVersion, "DescribeCaptchaInfoListInternational")
    
    
    return
}

func NewDescribeCaptchaInfoListInternationalResponse() (response *DescribeCaptchaInfoListInternationalResponse) {
    response = &DescribeCaptchaInfoListInternationalResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeCaptchaInfoListInternational
// Query the Captcha list to obtain all verification CaptchaAppIds, verification names, and other information internationally.
//
// error code that may be returned:
//  AUTHFAILURE = "AuthFailure"
//  AUTHFAILURE_UNAUTHORIZEDOPERATION = "AuthFailure.UnauthorizedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  UNAUTHORIZEDOPERATION_ERRAUTH = "UnauthorizedOperation.ErrAuth"
//  UNAUTHORIZEDOPERATION_UNAUTHORIZED = "UnauthorizedOperation.Unauthorized"
func (c *Client) DescribeCaptchaInfoListInternational(request *DescribeCaptchaInfoListInternationalRequest) (response *DescribeCaptchaInfoListInternationalResponse, err error) {
    return c.DescribeCaptchaInfoListInternationalWithContext(context.Background(), request)
}

// DescribeCaptchaInfoListInternational
// Query the Captcha list to obtain all verification CaptchaAppIds, verification names, and other information internationally.
//
// error code that may be returned:
//  AUTHFAILURE = "AuthFailure"
//  AUTHFAILURE_UNAUTHORIZEDOPERATION = "AuthFailure.UnauthorizedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  UNAUTHORIZEDOPERATION_ERRAUTH = "UnauthorizedOperation.ErrAuth"
//  UNAUTHORIZEDOPERATION_UNAUTHORIZED = "UnauthorizedOperation.Unauthorized"
func (c *Client) DescribeCaptchaInfoListInternationalWithContext(ctx context.Context, request *DescribeCaptchaInfoListInternationalRequest) (response *DescribeCaptchaInfoListInternationalResponse, err error) {
    if request == nil {
        request = NewDescribeCaptchaInfoListInternationalRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "captcha", APIVersion, "DescribeCaptchaInfoListInternational")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeCaptchaInfoListInternational require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeCaptchaInfoListInternationalResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeCaptchaResultRequest() (request *DescribeCaptchaResultRequest) {
    request = &DescribeCaptchaResultRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("captcha", APIVersion, "DescribeCaptchaResult")
    
    
    return
}

func NewDescribeCaptchaResultResponse() (response *DescribeCaptchaResultResponse) {
    response = &DescribeCaptchaResultResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeCaptchaResult
// This API is used to query the result of CAPTCHA ticket verification (web and app).
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  MISSINGPARAMETER = "MissingParameter"
//  UNAUTHORIZEDOPERATION_ERRAUTH = "UnauthorizedOperation.ErrAuth"
//  UNAUTHORIZEDOPERATION_UNAUTHORIZED = "UnauthorizedOperation.Unauthorized"
func (c *Client) DescribeCaptchaResult(request *DescribeCaptchaResultRequest) (response *DescribeCaptchaResultResponse, err error) {
    return c.DescribeCaptchaResultWithContext(context.Background(), request)
}

// DescribeCaptchaResult
// This API is used to query the result of CAPTCHA ticket verification (web and app).
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  MISSINGPARAMETER = "MissingParameter"
//  UNAUTHORIZEDOPERATION_ERRAUTH = "UnauthorizedOperation.ErrAuth"
//  UNAUTHORIZEDOPERATION_UNAUTHORIZED = "UnauthorizedOperation.Unauthorized"
func (c *Client) DescribeCaptchaResultWithContext(ctx context.Context, request *DescribeCaptchaResultRequest) (response *DescribeCaptchaResultResponse, err error) {
    if request == nil {
        request = NewDescribeCaptchaResultRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "captcha", APIVersion, "DescribeCaptchaResult")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeCaptchaResult require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeCaptchaResultResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeIpWhiteListInternationalRequest() (request *DescribeIpWhiteListInternationalRequest) {
    request = &DescribeIpWhiteListInternationalRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("captcha", APIVersion, "DescribeIpWhiteListInternational")
    
    
    return
}

func NewDescribeIpWhiteListInternationalResponse() (response *DescribeIpWhiteListInternationalResponse) {
    response = &DescribeIpWhiteListInternationalResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeIpWhiteListInternational
// IP allowlist list: You can query the IP whitelist list based on different business needs.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  UNAUTHORIZEDOPERATION_ERRAUTH = "UnauthorizedOperation.ErrAuth"
//  UNAUTHORIZEDOPERATION_UNAUTHORIZED = "UnauthorizedOperation.Unauthorized"
func (c *Client) DescribeIpWhiteListInternational(request *DescribeIpWhiteListInternationalRequest) (response *DescribeIpWhiteListInternationalResponse, err error) {
    return c.DescribeIpWhiteListInternationalWithContext(context.Background(), request)
}

// DescribeIpWhiteListInternational
// IP allowlist list: You can query the IP whitelist list based on different business needs.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  UNAUTHORIZEDOPERATION_ERRAUTH = "UnauthorizedOperation.ErrAuth"
//  UNAUTHORIZEDOPERATION_UNAUTHORIZED = "UnauthorizedOperation.Unauthorized"
func (c *Client) DescribeIpWhiteListInternationalWithContext(ctx context.Context, request *DescribeIpWhiteListInternationalRequest) (response *DescribeIpWhiteListInternationalResponse, err error) {
    if request == nil {
        request = NewDescribeIpWhiteListInternationalRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "captcha", APIVersion, "DescribeIpWhiteListInternational")
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeIpWhiteListInternational require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeIpWhiteListInternationalResponse()
    err = c.Send(request, response)
    return
}

func NewModifyCaptchaInfoInternationalRequest() (request *ModifyCaptchaInfoInternationalRequest) {
    request = &ModifyCaptchaInfoInternationalRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("captcha", APIVersion, "ModifyCaptchaInfoInternational")
    
    
    return
}

func NewModifyCaptchaInfoInternationalResponse() (response *ModifyCaptchaInfoInternationalResponse) {
    response = &ModifyCaptchaInfoInternationalResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// ModifyCaptchaInfoInternational
// Change the captcha configuration, including basic, appearance, and security settings such as captcha name, prompt language, and validation type.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  UNAUTHORIZEDOPERATION_ERRAUTH = "UnauthorizedOperation.ErrAuth"
//  UNAUTHORIZEDOPERATION_UNAUTHORIZED = "UnauthorizedOperation.Unauthorized"
func (c *Client) ModifyCaptchaInfoInternational(request *ModifyCaptchaInfoInternationalRequest) (response *ModifyCaptchaInfoInternationalResponse, err error) {
    return c.ModifyCaptchaInfoInternationalWithContext(context.Background(), request)
}

// ModifyCaptchaInfoInternational
// Change the captcha configuration, including basic, appearance, and security settings such as captcha name, prompt language, and validation type.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  UNAUTHORIZEDOPERATION_ERRAUTH = "UnauthorizedOperation.ErrAuth"
//  UNAUTHORIZEDOPERATION_UNAUTHORIZED = "UnauthorizedOperation.Unauthorized"
func (c *Client) ModifyCaptchaInfoInternationalWithContext(ctx context.Context, request *ModifyCaptchaInfoInternationalRequest) (response *ModifyCaptchaInfoInternationalResponse, err error) {
    if request == nil {
        request = NewModifyCaptchaInfoInternationalRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "captcha", APIVersion, "ModifyCaptchaInfoInternational")
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyCaptchaInfoInternational require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyCaptchaInfoInternationalResponse()
    err = c.Send(request, response)
    return
}

func NewModifyIpWhiteListInternationalRequest() (request *ModifyIpWhiteListInternationalRequest) {
    request = &ModifyIpWhiteListInternationalRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("captcha", APIVersion, "ModifyIpWhiteListInternational")
    
    
    return
}

func NewModifyIpWhiteListInternationalResponse() (response *ModifyIpWhiteListInternationalResponse) {
    response = &ModifyIpWhiteListInternationalResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// ModifyIpWhiteListInternational
// Edit IP allowlist: You can edit the IP allowlist based on different business needs.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  UNAUTHORIZEDOPERATION_ERRAUTH = "UnauthorizedOperation.ErrAuth"
//  UNAUTHORIZEDOPERATION_UNAUTHORIZED = "UnauthorizedOperation.Unauthorized"
func (c *Client) ModifyIpWhiteListInternational(request *ModifyIpWhiteListInternationalRequest) (response *ModifyIpWhiteListInternationalResponse, err error) {
    return c.ModifyIpWhiteListInternationalWithContext(context.Background(), request)
}

// ModifyIpWhiteListInternational
// Edit IP allowlist: You can edit the IP allowlist based on different business needs.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  UNAUTHORIZEDOPERATION_ERRAUTH = "UnauthorizedOperation.ErrAuth"
//  UNAUTHORIZEDOPERATION_UNAUTHORIZED = "UnauthorizedOperation.Unauthorized"
func (c *Client) ModifyIpWhiteListInternationalWithContext(ctx context.Context, request *ModifyIpWhiteListInternationalRequest) (response *ModifyIpWhiteListInternationalResponse, err error) {
    if request == nil {
        request = NewModifyIpWhiteListInternationalRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "captcha", APIVersion, "ModifyIpWhiteListInternational")
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyIpWhiteListInternational require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyIpWhiteListInternationalResponse()
    err = c.Send(request, response)
    return
}

func NewRemoveCaptchaInfoInternationalRequest() (request *RemoveCaptchaInfoInternationalRequest) {
    request = &RemoveCaptchaInfoInternationalRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("captcha", APIVersion, "RemoveCaptchaInfoInternational")
    
    
    return
}

func NewRemoveCaptchaInfoInternationalResponse() (response *RemoveCaptchaInfoInternationalResponse) {
    response = &RemoveCaptchaInfoInternationalResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// RemoveCaptchaInfoInternational
// Delete a captcha: once deleted, verification scenarios using this CaptchaAppId will fail to load the verification code on the frontend, and invoice verification will report an error on the backend. Proceed with caution.
//
// error code that may be returned:
//  AUTHFAILURE_UNAUTHORIZEDOPERATION = "AuthFailure.UnauthorizedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  UNAUTHORIZEDOPERATION_ERRAUTH = "UnauthorizedOperation.ErrAuth"
//  UNAUTHORIZEDOPERATION_UNAUTHORIZED = "UnauthorizedOperation.Unauthorized"
func (c *Client) RemoveCaptchaInfoInternational(request *RemoveCaptchaInfoInternationalRequest) (response *RemoveCaptchaInfoInternationalResponse, err error) {
    return c.RemoveCaptchaInfoInternationalWithContext(context.Background(), request)
}

// RemoveCaptchaInfoInternational
// Delete a captcha: once deleted, verification scenarios using this CaptchaAppId will fail to load the verification code on the frontend, and invoice verification will report an error on the backend. Proceed with caution.
//
// error code that may be returned:
//  AUTHFAILURE_UNAUTHORIZEDOPERATION = "AuthFailure.UnauthorizedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  UNAUTHORIZEDOPERATION_ERRAUTH = "UnauthorizedOperation.ErrAuth"
//  UNAUTHORIZEDOPERATION_UNAUTHORIZED = "UnauthorizedOperation.Unauthorized"
func (c *Client) RemoveCaptchaInfoInternationalWithContext(ctx context.Context, request *RemoveCaptchaInfoInternationalRequest) (response *RemoveCaptchaInfoInternationalResponse, err error) {
    if request == nil {
        request = NewRemoveCaptchaInfoInternationalRequest()
    }
    c.InitBaseRequest(&request.BaseRequest, "captcha", APIVersion, "RemoveCaptchaInfoInternational")
    
    if c.GetCredential() == nil {
        return nil, errors.New("RemoveCaptchaInfoInternational require credential")
    }

    request.SetContext(ctx)
    
    response = NewRemoveCaptchaInfoInternationalResponse()
    err = c.Send(request, response)
    return
}
