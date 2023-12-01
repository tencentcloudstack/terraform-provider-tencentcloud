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

package v20200326

import (
    "context"
    "errors"
    "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
    tchttp "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common/http"
    "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common/profile"
)

const APIVersion = "2020-03-26"

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


func NewCreateStreamLiveChannelRequest() (request *CreateStreamLiveChannelRequest) {
    request = &CreateStreamLiveChannelRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("mdl", APIVersion, "CreateStreamLiveChannel")
    
    
    return
}

func NewCreateStreamLiveChannelResponse() (response *CreateStreamLiveChannelResponse) {
    response = &CreateStreamLiveChannelResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// CreateStreamLiveChannel
// This API is used to create a StreamLive channel.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER_AVTEMPLATES = "InvalidParameter.AVTemplates"
//  INVALIDPARAMETER_ALREADYASSOCIATEDCHANNEL = "InvalidParameter.AlreadyAssociatedChannel"
//  INVALIDPARAMETER_ATTACHEDINPUTS = "InvalidParameter.AttachedInputs"
//  INVALIDPARAMETER_AUDIOTEMPLATES = "InvalidParameter.AudioTemplates"
//  INVALIDPARAMETER_EXCEEDEDQUANTITYLIMIT = "InvalidParameter.ExceededQuantityLimit"
//  INVALIDPARAMETER_ID = "InvalidParameter.Id"
//  INVALIDPARAMETER_NAME = "InvalidParameter.Name"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
//  INVALIDPARAMETER_NOTIFYKEY = "InvalidParameter.NotifyKey"
//  INVALIDPARAMETER_NOTIFYURL = "InvalidParameter.NotifyUrl"
//  INVALIDPARAMETER_OUTPUTGROUPS = "InvalidParameter.OutputGroups"
//  INVALIDPARAMETER_VIDEOTEMPLATES = "InvalidParameter.VideoTemplates"
func (c *Client) CreateStreamLiveChannel(request *CreateStreamLiveChannelRequest) (response *CreateStreamLiveChannelResponse, err error) {
    return c.CreateStreamLiveChannelWithContext(context.Background(), request)
}

// CreateStreamLiveChannel
// This API is used to create a StreamLive channel.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER_AVTEMPLATES = "InvalidParameter.AVTemplates"
//  INVALIDPARAMETER_ALREADYASSOCIATEDCHANNEL = "InvalidParameter.AlreadyAssociatedChannel"
//  INVALIDPARAMETER_ATTACHEDINPUTS = "InvalidParameter.AttachedInputs"
//  INVALIDPARAMETER_AUDIOTEMPLATES = "InvalidParameter.AudioTemplates"
//  INVALIDPARAMETER_EXCEEDEDQUANTITYLIMIT = "InvalidParameter.ExceededQuantityLimit"
//  INVALIDPARAMETER_ID = "InvalidParameter.Id"
//  INVALIDPARAMETER_NAME = "InvalidParameter.Name"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
//  INVALIDPARAMETER_NOTIFYKEY = "InvalidParameter.NotifyKey"
//  INVALIDPARAMETER_NOTIFYURL = "InvalidParameter.NotifyUrl"
//  INVALIDPARAMETER_OUTPUTGROUPS = "InvalidParameter.OutputGroups"
//  INVALIDPARAMETER_VIDEOTEMPLATES = "InvalidParameter.VideoTemplates"
func (c *Client) CreateStreamLiveChannelWithContext(ctx context.Context, request *CreateStreamLiveChannelRequest) (response *CreateStreamLiveChannelResponse, err error) {
    if request == nil {
        request = NewCreateStreamLiveChannelRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateStreamLiveChannel require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateStreamLiveChannelResponse()
    err = c.Send(request, response)
    return
}

func NewCreateStreamLiveInputRequest() (request *CreateStreamLiveInputRequest) {
    request = &CreateStreamLiveInputRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("mdl", APIVersion, "CreateStreamLiveInput")
    
    
    return
}

func NewCreateStreamLiveInputResponse() (response *CreateStreamLiveInputResponse) {
    response = &CreateStreamLiveInputResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// CreateStreamLiveInput
// This API is used to create a StreamLive input.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_EXCEEDEDQUANTITYLIMIT = "InvalidParameter.ExceededQuantityLimit"
//  INVALIDPARAMETER_INPUTSETTINGS = "InvalidParameter.InputSettings"
//  INVALIDPARAMETER_NAME = "InvalidParameter.Name"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
//  INVALIDPARAMETER_SECURITYGROUPS = "InvalidParameter.SecurityGroups"
//  INVALIDPARAMETER_TYPE = "InvalidParameter.Type"
func (c *Client) CreateStreamLiveInput(request *CreateStreamLiveInputRequest) (response *CreateStreamLiveInputResponse, err error) {
    return c.CreateStreamLiveInputWithContext(context.Background(), request)
}

// CreateStreamLiveInput
// This API is used to create a StreamLive input.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_EXCEEDEDQUANTITYLIMIT = "InvalidParameter.ExceededQuantityLimit"
//  INVALIDPARAMETER_INPUTSETTINGS = "InvalidParameter.InputSettings"
//  INVALIDPARAMETER_NAME = "InvalidParameter.Name"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
//  INVALIDPARAMETER_SECURITYGROUPS = "InvalidParameter.SecurityGroups"
//  INVALIDPARAMETER_TYPE = "InvalidParameter.Type"
func (c *Client) CreateStreamLiveInputWithContext(ctx context.Context, request *CreateStreamLiveInputRequest) (response *CreateStreamLiveInputResponse, err error) {
    if request == nil {
        request = NewCreateStreamLiveInputRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateStreamLiveInput require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateStreamLiveInputResponse()
    err = c.Send(request, response)
    return
}

func NewCreateStreamLiveInputSecurityGroupRequest() (request *CreateStreamLiveInputSecurityGroupRequest) {
    request = &CreateStreamLiveInputSecurityGroupRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("mdl", APIVersion, "CreateStreamLiveInputSecurityGroup")
    
    
    return
}

func NewCreateStreamLiveInputSecurityGroupResponse() (response *CreateStreamLiveInputSecurityGroupResponse) {
    response = &CreateStreamLiveInputSecurityGroupResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// CreateStreamLiveInputSecurityGroup
// This API is used to create an input security group. Up to 5 security groups are allowed.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_EXCEEDEDQUANTITYLIMIT = "InvalidParameter.ExceededQuantityLimit"
//  INVALIDPARAMETER_NAME = "InvalidParameter.Name"
//  INVALIDPARAMETER_WHITELIST = "InvalidParameter.Whitelist"
func (c *Client) CreateStreamLiveInputSecurityGroup(request *CreateStreamLiveInputSecurityGroupRequest) (response *CreateStreamLiveInputSecurityGroupResponse, err error) {
    return c.CreateStreamLiveInputSecurityGroupWithContext(context.Background(), request)
}

// CreateStreamLiveInputSecurityGroup
// This API is used to create an input security group. Up to 5 security groups are allowed.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_EXCEEDEDQUANTITYLIMIT = "InvalidParameter.ExceededQuantityLimit"
//  INVALIDPARAMETER_NAME = "InvalidParameter.Name"
//  INVALIDPARAMETER_WHITELIST = "InvalidParameter.Whitelist"
func (c *Client) CreateStreamLiveInputSecurityGroupWithContext(ctx context.Context, request *CreateStreamLiveInputSecurityGroupRequest) (response *CreateStreamLiveInputSecurityGroupResponse, err error) {
    if request == nil {
        request = NewCreateStreamLiveInputSecurityGroupRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateStreamLiveInputSecurityGroup require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateStreamLiveInputSecurityGroupResponse()
    err = c.Send(request, response)
    return
}

func NewCreateStreamLivePlanRequest() (request *CreateStreamLivePlanRequest) {
    request = &CreateStreamLivePlanRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("mdl", APIVersion, "CreateStreamLivePlan")
    
    
    return
}

func NewCreateStreamLivePlanResponse() (response *CreateStreamLivePlanResponse) {
    response = &CreateStreamLivePlanResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// CreateStreamLivePlan
// This API is used to create an event in the plan.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
//  INVALIDPARAMETER_PLAN = "InvalidParameter.Plan"
func (c *Client) CreateStreamLivePlan(request *CreateStreamLivePlanRequest) (response *CreateStreamLivePlanResponse, err error) {
    return c.CreateStreamLivePlanWithContext(context.Background(), request)
}

// CreateStreamLivePlan
// This API is used to create an event in the plan.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
//  INVALIDPARAMETER_PLAN = "InvalidParameter.Plan"
func (c *Client) CreateStreamLivePlanWithContext(ctx context.Context, request *CreateStreamLivePlanRequest) (response *CreateStreamLivePlanResponse, err error) {
    if request == nil {
        request = NewCreateStreamLivePlanRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateStreamLivePlan require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateStreamLivePlanResponse()
    err = c.Send(request, response)
    return
}

func NewCreateStreamLiveWatermarkRequest() (request *CreateStreamLiveWatermarkRequest) {
    request = &CreateStreamLiveWatermarkRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("mdl", APIVersion, "CreateStreamLiveWatermark")
    
    
    return
}

func NewCreateStreamLiveWatermarkResponse() (response *CreateStreamLiveWatermarkResponse) {
    response = &CreateStreamLiveWatermarkResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// CreateStreamLiveWatermark
// This API is used to add a watermark.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER_EXCEEDEDQUANTITYLIMIT = "InvalidParameter.ExceededQuantityLimit"
//  INVALIDPARAMETER_IMAGESETTINGS = "InvalidParameter.ImageSettings"
//  INVALIDPARAMETER_NAME = "InvalidParameter.Name"
//  INVALIDPARAMETER_TEXTSETTINGS = "InvalidParameter.TextSettings"
//  INVALIDPARAMETER_TYPE = "InvalidParameter.Type"
func (c *Client) CreateStreamLiveWatermark(request *CreateStreamLiveWatermarkRequest) (response *CreateStreamLiveWatermarkResponse, err error) {
    return c.CreateStreamLiveWatermarkWithContext(context.Background(), request)
}

// CreateStreamLiveWatermark
// This API is used to add a watermark.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER_EXCEEDEDQUANTITYLIMIT = "InvalidParameter.ExceededQuantityLimit"
//  INVALIDPARAMETER_IMAGESETTINGS = "InvalidParameter.ImageSettings"
//  INVALIDPARAMETER_NAME = "InvalidParameter.Name"
//  INVALIDPARAMETER_TEXTSETTINGS = "InvalidParameter.TextSettings"
//  INVALIDPARAMETER_TYPE = "InvalidParameter.Type"
func (c *Client) CreateStreamLiveWatermarkWithContext(ctx context.Context, request *CreateStreamLiveWatermarkRequest) (response *CreateStreamLiveWatermarkResponse, err error) {
    if request == nil {
        request = NewCreateStreamLiveWatermarkRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateStreamLiveWatermark require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateStreamLiveWatermarkResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteStreamLiveChannelRequest() (request *DeleteStreamLiveChannelRequest) {
    request = &DeleteStreamLiveChannelRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("mdl", APIVersion, "DeleteStreamLiveChannel")
    
    
    return
}

func NewDeleteStreamLiveChannelResponse() (response *DeleteStreamLiveChannelResponse) {
    response = &DeleteStreamLiveChannelResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DeleteStreamLiveChannel
// This API is used to delete a StreamLive channel.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_ID = "InvalidParameter.Id"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
//  INVALIDPARAMETER_STATE = "InvalidParameter.State"
//  INVALIDPARAMETER_STATEERROR = "InvalidParameter.StateError"
func (c *Client) DeleteStreamLiveChannel(request *DeleteStreamLiveChannelRequest) (response *DeleteStreamLiveChannelResponse, err error) {
    return c.DeleteStreamLiveChannelWithContext(context.Background(), request)
}

// DeleteStreamLiveChannel
// This API is used to delete a StreamLive channel.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_ID = "InvalidParameter.Id"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
//  INVALIDPARAMETER_STATE = "InvalidParameter.State"
//  INVALIDPARAMETER_STATEERROR = "InvalidParameter.StateError"
func (c *Client) DeleteStreamLiveChannelWithContext(ctx context.Context, request *DeleteStreamLiveChannelRequest) (response *DeleteStreamLiveChannelResponse, err error) {
    if request == nil {
        request = NewDeleteStreamLiveChannelRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteStreamLiveChannel require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteStreamLiveChannelResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteStreamLiveInputRequest() (request *DeleteStreamLiveInputRequest) {
    request = &DeleteStreamLiveInputRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("mdl", APIVersion, "DeleteStreamLiveInput")
    
    
    return
}

func NewDeleteStreamLiveInputResponse() (response *DeleteStreamLiveInputResponse) {
    response = &DeleteStreamLiveInputResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DeleteStreamLiveInput
// This API is used to delete a StreamLive input.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_ALREADYASSOCIATEDCHANNEL = "InvalidParameter.AlreadyAssociatedChannel"
//  INVALIDPARAMETER_ID = "InvalidParameter.Id"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
func (c *Client) DeleteStreamLiveInput(request *DeleteStreamLiveInputRequest) (response *DeleteStreamLiveInputResponse, err error) {
    return c.DeleteStreamLiveInputWithContext(context.Background(), request)
}

// DeleteStreamLiveInput
// This API is used to delete a StreamLive input.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_ALREADYASSOCIATEDCHANNEL = "InvalidParameter.AlreadyAssociatedChannel"
//  INVALIDPARAMETER_ID = "InvalidParameter.Id"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
func (c *Client) DeleteStreamLiveInputWithContext(ctx context.Context, request *DeleteStreamLiveInputRequest) (response *DeleteStreamLiveInputResponse, err error) {
    if request == nil {
        request = NewDeleteStreamLiveInputRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteStreamLiveInput require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteStreamLiveInputResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteStreamLiveInputSecurityGroupRequest() (request *DeleteStreamLiveInputSecurityGroupRequest) {
    request = &DeleteStreamLiveInputSecurityGroupRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("mdl", APIVersion, "DeleteStreamLiveInputSecurityGroup")
    
    
    return
}

func NewDeleteStreamLiveInputSecurityGroupResponse() (response *DeleteStreamLiveInputSecurityGroupResponse) {
    response = &DeleteStreamLiveInputSecurityGroupResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DeleteStreamLiveInputSecurityGroup
// This API is used to delete an input security group.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_ALREADYASSOCIATEDINPUT = "InvalidParameter.AlreadyAssociatedInput"
//  INVALIDPARAMETER_ID = "InvalidParameter.Id"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
func (c *Client) DeleteStreamLiveInputSecurityGroup(request *DeleteStreamLiveInputSecurityGroupRequest) (response *DeleteStreamLiveInputSecurityGroupResponse, err error) {
    return c.DeleteStreamLiveInputSecurityGroupWithContext(context.Background(), request)
}

// DeleteStreamLiveInputSecurityGroup
// This API is used to delete an input security group.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_ALREADYASSOCIATEDINPUT = "InvalidParameter.AlreadyAssociatedInput"
//  INVALIDPARAMETER_ID = "InvalidParameter.Id"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
func (c *Client) DeleteStreamLiveInputSecurityGroupWithContext(ctx context.Context, request *DeleteStreamLiveInputSecurityGroupRequest) (response *DeleteStreamLiveInputSecurityGroupResponse, err error) {
    if request == nil {
        request = NewDeleteStreamLiveInputSecurityGroupRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteStreamLiveInputSecurityGroup require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteStreamLiveInputSecurityGroupResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteStreamLivePlanRequest() (request *DeleteStreamLivePlanRequest) {
    request = &DeleteStreamLivePlanRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("mdl", APIVersion, "DeleteStreamLivePlan")
    
    
    return
}

func NewDeleteStreamLivePlanResponse() (response *DeleteStreamLivePlanResponse) {
    response = &DeleteStreamLivePlanResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DeleteStreamLivePlan
// This API is used to delete a StreamLive event.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
//  INVALIDPARAMETER_PLAN = "InvalidParameter.Plan"
func (c *Client) DeleteStreamLivePlan(request *DeleteStreamLivePlanRequest) (response *DeleteStreamLivePlanResponse, err error) {
    return c.DeleteStreamLivePlanWithContext(context.Background(), request)
}

// DeleteStreamLivePlan
// This API is used to delete a StreamLive event.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
//  INVALIDPARAMETER_PLAN = "InvalidParameter.Plan"
func (c *Client) DeleteStreamLivePlanWithContext(ctx context.Context, request *DeleteStreamLivePlanRequest) (response *DeleteStreamLivePlanResponse, err error) {
    if request == nil {
        request = NewDeleteStreamLivePlanRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteStreamLivePlan require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteStreamLivePlanResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteStreamLiveWatermarkRequest() (request *DeleteStreamLiveWatermarkRequest) {
    request = &DeleteStreamLiveWatermarkRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("mdl", APIVersion, "DeleteStreamLiveWatermark")
    
    
    return
}

func NewDeleteStreamLiveWatermarkResponse() (response *DeleteStreamLiveWatermarkResponse) {
    response = &DeleteStreamLiveWatermarkResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DeleteStreamLiveWatermark
// This API is used to delete a watermark.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER_ALREADYASSOCIATEDCHANNEL = "InvalidParameter.AlreadyAssociatedChannel"
//  INVALIDPARAMETER_ID = "InvalidParameter.Id"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
func (c *Client) DeleteStreamLiveWatermark(request *DeleteStreamLiveWatermarkRequest) (response *DeleteStreamLiveWatermarkResponse, err error) {
    return c.DeleteStreamLiveWatermarkWithContext(context.Background(), request)
}

// DeleteStreamLiveWatermark
// This API is used to delete a watermark.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER_ALREADYASSOCIATEDCHANNEL = "InvalidParameter.AlreadyAssociatedChannel"
//  INVALIDPARAMETER_ID = "InvalidParameter.Id"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
func (c *Client) DeleteStreamLiveWatermarkWithContext(ctx context.Context, request *DeleteStreamLiveWatermarkRequest) (response *DeleteStreamLiveWatermarkResponse, err error) {
    if request == nil {
        request = NewDeleteStreamLiveWatermarkRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteStreamLiveWatermark require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteStreamLiveWatermarkResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeStreamLiveChannelRequest() (request *DescribeStreamLiveChannelRequest) {
    request = &DescribeStreamLiveChannelRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("mdl", APIVersion, "DescribeStreamLiveChannel")
    
    
    return
}

func NewDescribeStreamLiveChannelResponse() (response *DescribeStreamLiveChannelResponse) {
    response = &DescribeStreamLiveChannelResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeStreamLiveChannel
// This API is used to query a StreamLive channel.
//
// error code that may be returned:
//  INVALIDPARAMETER_ID = "InvalidParameter.Id"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
func (c *Client) DescribeStreamLiveChannel(request *DescribeStreamLiveChannelRequest) (response *DescribeStreamLiveChannelResponse, err error) {
    return c.DescribeStreamLiveChannelWithContext(context.Background(), request)
}

// DescribeStreamLiveChannel
// This API is used to query a StreamLive channel.
//
// error code that may be returned:
//  INVALIDPARAMETER_ID = "InvalidParameter.Id"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
func (c *Client) DescribeStreamLiveChannelWithContext(ctx context.Context, request *DescribeStreamLiveChannelRequest) (response *DescribeStreamLiveChannelResponse, err error) {
    if request == nil {
        request = NewDescribeStreamLiveChannelRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeStreamLiveChannel require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeStreamLiveChannelResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeStreamLiveChannelAlertsRequest() (request *DescribeStreamLiveChannelAlertsRequest) {
    request = &DescribeStreamLiveChannelAlertsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("mdl", APIVersion, "DescribeStreamLiveChannelAlerts")
    
    
    return
}

func NewDescribeStreamLiveChannelAlertsResponse() (response *DescribeStreamLiveChannelAlertsResponse) {
    response = &DescribeStreamLiveChannelAlertsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeStreamLiveChannelAlerts
// This API is used to query the alarm information of a StreamLive channel.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
func (c *Client) DescribeStreamLiveChannelAlerts(request *DescribeStreamLiveChannelAlertsRequest) (response *DescribeStreamLiveChannelAlertsResponse, err error) {
    return c.DescribeStreamLiveChannelAlertsWithContext(context.Background(), request)
}

// DescribeStreamLiveChannelAlerts
// This API is used to query the alarm information of a StreamLive channel.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
func (c *Client) DescribeStreamLiveChannelAlertsWithContext(ctx context.Context, request *DescribeStreamLiveChannelAlertsRequest) (response *DescribeStreamLiveChannelAlertsResponse, err error) {
    if request == nil {
        request = NewDescribeStreamLiveChannelAlertsRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeStreamLiveChannelAlerts require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeStreamLiveChannelAlertsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeStreamLiveChannelInputStatisticsRequest() (request *DescribeStreamLiveChannelInputStatisticsRequest) {
    request = &DescribeStreamLiveChannelInputStatisticsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("mdl", APIVersion, "DescribeStreamLiveChannelInputStatistics")
    
    
    return
}

func NewDescribeStreamLiveChannelInputStatisticsResponse() (response *DescribeStreamLiveChannelInputStatisticsResponse) {
    response = &DescribeStreamLiveChannelInputStatisticsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeStreamLiveChannelInputStatistics
// This API is used to query input statistics.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER_ID = "InvalidParameter.Id"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
//  INVALIDPARAMETER_STARTTIME = "InvalidParameter.StartTime"
func (c *Client) DescribeStreamLiveChannelInputStatistics(request *DescribeStreamLiveChannelInputStatisticsRequest) (response *DescribeStreamLiveChannelInputStatisticsResponse, err error) {
    return c.DescribeStreamLiveChannelInputStatisticsWithContext(context.Background(), request)
}

// DescribeStreamLiveChannelInputStatistics
// This API is used to query input statistics.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER_ID = "InvalidParameter.Id"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
//  INVALIDPARAMETER_STARTTIME = "InvalidParameter.StartTime"
func (c *Client) DescribeStreamLiveChannelInputStatisticsWithContext(ctx context.Context, request *DescribeStreamLiveChannelInputStatisticsRequest) (response *DescribeStreamLiveChannelInputStatisticsResponse, err error) {
    if request == nil {
        request = NewDescribeStreamLiveChannelInputStatisticsRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeStreamLiveChannelInputStatistics require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeStreamLiveChannelInputStatisticsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeStreamLiveChannelLogsRequest() (request *DescribeStreamLiveChannelLogsRequest) {
    request = &DescribeStreamLiveChannelLogsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("mdl", APIVersion, "DescribeStreamLiveChannelLogs")
    
    
    return
}

func NewDescribeStreamLiveChannelLogsResponse() (response *DescribeStreamLiveChannelLogsResponse) {
    response = &DescribeStreamLiveChannelLogsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeStreamLiveChannelLogs
// This API is used to query StreamLive channel logs, such as push event logs.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER_ENDTIME = "InvalidParameter.EndTime"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
//  INVALIDPARAMETER_STARTTIME = "InvalidParameter.StartTime"
func (c *Client) DescribeStreamLiveChannelLogs(request *DescribeStreamLiveChannelLogsRequest) (response *DescribeStreamLiveChannelLogsResponse, err error) {
    return c.DescribeStreamLiveChannelLogsWithContext(context.Background(), request)
}

// DescribeStreamLiveChannelLogs
// This API is used to query StreamLive channel logs, such as push event logs.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER_ENDTIME = "InvalidParameter.EndTime"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
//  INVALIDPARAMETER_STARTTIME = "InvalidParameter.StartTime"
func (c *Client) DescribeStreamLiveChannelLogsWithContext(ctx context.Context, request *DescribeStreamLiveChannelLogsRequest) (response *DescribeStreamLiveChannelLogsResponse, err error) {
    if request == nil {
        request = NewDescribeStreamLiveChannelLogsRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeStreamLiveChannelLogs require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeStreamLiveChannelLogsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeStreamLiveChannelOutputStatisticsRequest() (request *DescribeStreamLiveChannelOutputStatisticsRequest) {
    request = &DescribeStreamLiveChannelOutputStatisticsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("mdl", APIVersion, "DescribeStreamLiveChannelOutputStatistics")
    
    
    return
}

func NewDescribeStreamLiveChannelOutputStatisticsResponse() (response *DescribeStreamLiveChannelOutputStatisticsResponse) {
    response = &DescribeStreamLiveChannelOutputStatisticsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeStreamLiveChannelOutputStatistics
// This API is used to query the output statistics of a StreamLive channel.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER_ENDTIME = "InvalidParameter.EndTime"
//  INVALIDPARAMETER_ID = "InvalidParameter.Id"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
//  INVALIDPARAMETER_STARTTIME = "InvalidParameter.StartTime"
func (c *Client) DescribeStreamLiveChannelOutputStatistics(request *DescribeStreamLiveChannelOutputStatisticsRequest) (response *DescribeStreamLiveChannelOutputStatisticsResponse, err error) {
    return c.DescribeStreamLiveChannelOutputStatisticsWithContext(context.Background(), request)
}

// DescribeStreamLiveChannelOutputStatistics
// This API is used to query the output statistics of a StreamLive channel.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER_ENDTIME = "InvalidParameter.EndTime"
//  INVALIDPARAMETER_ID = "InvalidParameter.Id"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
//  INVALIDPARAMETER_STARTTIME = "InvalidParameter.StartTime"
func (c *Client) DescribeStreamLiveChannelOutputStatisticsWithContext(ctx context.Context, request *DescribeStreamLiveChannelOutputStatisticsRequest) (response *DescribeStreamLiveChannelOutputStatisticsResponse, err error) {
    if request == nil {
        request = NewDescribeStreamLiveChannelOutputStatisticsRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeStreamLiveChannelOutputStatistics require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeStreamLiveChannelOutputStatisticsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeStreamLiveChannelsRequest() (request *DescribeStreamLiveChannelsRequest) {
    request = &DescribeStreamLiveChannelsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("mdl", APIVersion, "DescribeStreamLiveChannels")
    
    
    return
}

func NewDescribeStreamLiveChannelsResponse() (response *DescribeStreamLiveChannelsResponse) {
    response = &DescribeStreamLiveChannelsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeStreamLiveChannels
// This API is used to query StreamLive channels in batches.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
func (c *Client) DescribeStreamLiveChannels(request *DescribeStreamLiveChannelsRequest) (response *DescribeStreamLiveChannelsResponse, err error) {
    return c.DescribeStreamLiveChannelsWithContext(context.Background(), request)
}

// DescribeStreamLiveChannels
// This API is used to query StreamLive channels in batches.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
func (c *Client) DescribeStreamLiveChannelsWithContext(ctx context.Context, request *DescribeStreamLiveChannelsRequest) (response *DescribeStreamLiveChannelsResponse, err error) {
    if request == nil {
        request = NewDescribeStreamLiveChannelsRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeStreamLiveChannels require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeStreamLiveChannelsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeStreamLiveInputRequest() (request *DescribeStreamLiveInputRequest) {
    request = &DescribeStreamLiveInputRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("mdl", APIVersion, "DescribeStreamLiveInput")
    
    
    return
}

func NewDescribeStreamLiveInputResponse() (response *DescribeStreamLiveInputResponse) {
    response = &DescribeStreamLiveInputResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeStreamLiveInput
// This API is used to query a StreamLive input.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_ID = "InvalidParameter.Id"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
func (c *Client) DescribeStreamLiveInput(request *DescribeStreamLiveInputRequest) (response *DescribeStreamLiveInputResponse, err error) {
    return c.DescribeStreamLiveInputWithContext(context.Background(), request)
}

// DescribeStreamLiveInput
// This API is used to query a StreamLive input.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_ID = "InvalidParameter.Id"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
func (c *Client) DescribeStreamLiveInputWithContext(ctx context.Context, request *DescribeStreamLiveInputRequest) (response *DescribeStreamLiveInputResponse, err error) {
    if request == nil {
        request = NewDescribeStreamLiveInputRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeStreamLiveInput require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeStreamLiveInputResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeStreamLiveInputSecurityGroupRequest() (request *DescribeStreamLiveInputSecurityGroupRequest) {
    request = &DescribeStreamLiveInputSecurityGroupRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("mdl", APIVersion, "DescribeStreamLiveInputSecurityGroup")
    
    
    return
}

func NewDescribeStreamLiveInputSecurityGroupResponse() (response *DescribeStreamLiveInputSecurityGroupResponse) {
    response = &DescribeStreamLiveInputSecurityGroupResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeStreamLiveInputSecurityGroup
// This API is used to query an input security group.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_ID = "InvalidParameter.Id"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
func (c *Client) DescribeStreamLiveInputSecurityGroup(request *DescribeStreamLiveInputSecurityGroupRequest) (response *DescribeStreamLiveInputSecurityGroupResponse, err error) {
    return c.DescribeStreamLiveInputSecurityGroupWithContext(context.Background(), request)
}

// DescribeStreamLiveInputSecurityGroup
// This API is used to query an input security group.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_ID = "InvalidParameter.Id"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
func (c *Client) DescribeStreamLiveInputSecurityGroupWithContext(ctx context.Context, request *DescribeStreamLiveInputSecurityGroupRequest) (response *DescribeStreamLiveInputSecurityGroupResponse, err error) {
    if request == nil {
        request = NewDescribeStreamLiveInputSecurityGroupRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeStreamLiveInputSecurityGroup require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeStreamLiveInputSecurityGroupResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeStreamLiveInputSecurityGroupsRequest() (request *DescribeStreamLiveInputSecurityGroupsRequest) {
    request = &DescribeStreamLiveInputSecurityGroupsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("mdl", APIVersion, "DescribeStreamLiveInputSecurityGroups")
    
    
    return
}

func NewDescribeStreamLiveInputSecurityGroupsResponse() (response *DescribeStreamLiveInputSecurityGroupsResponse) {
    response = &DescribeStreamLiveInputSecurityGroupsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeStreamLiveInputSecurityGroups
// This API is used to query input security groups in batches.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
func (c *Client) DescribeStreamLiveInputSecurityGroups(request *DescribeStreamLiveInputSecurityGroupsRequest) (response *DescribeStreamLiveInputSecurityGroupsResponse, err error) {
    return c.DescribeStreamLiveInputSecurityGroupsWithContext(context.Background(), request)
}

// DescribeStreamLiveInputSecurityGroups
// This API is used to query input security groups in batches.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
func (c *Client) DescribeStreamLiveInputSecurityGroupsWithContext(ctx context.Context, request *DescribeStreamLiveInputSecurityGroupsRequest) (response *DescribeStreamLiveInputSecurityGroupsResponse, err error) {
    if request == nil {
        request = NewDescribeStreamLiveInputSecurityGroupsRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeStreamLiveInputSecurityGroups require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeStreamLiveInputSecurityGroupsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeStreamLiveInputsRequest() (request *DescribeStreamLiveInputsRequest) {
    request = &DescribeStreamLiveInputsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("mdl", APIVersion, "DescribeStreamLiveInputs")
    
    
    return
}

func NewDescribeStreamLiveInputsResponse() (response *DescribeStreamLiveInputsResponse) {
    response = &DescribeStreamLiveInputsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeStreamLiveInputs
// This API is used to query StreamLive inputs in batches.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
func (c *Client) DescribeStreamLiveInputs(request *DescribeStreamLiveInputsRequest) (response *DescribeStreamLiveInputsResponse, err error) {
    return c.DescribeStreamLiveInputsWithContext(context.Background(), request)
}

// DescribeStreamLiveInputs
// This API is used to query StreamLive inputs in batches.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
func (c *Client) DescribeStreamLiveInputsWithContext(ctx context.Context, request *DescribeStreamLiveInputsRequest) (response *DescribeStreamLiveInputsResponse, err error) {
    if request == nil {
        request = NewDescribeStreamLiveInputsRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeStreamLiveInputs require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeStreamLiveInputsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeStreamLivePlansRequest() (request *DescribeStreamLivePlansRequest) {
    request = &DescribeStreamLivePlansRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("mdl", APIVersion, "DescribeStreamLivePlans")
    
    
    return
}

func NewDescribeStreamLivePlansResponse() (response *DescribeStreamLivePlansResponse) {
    response = &DescribeStreamLivePlansResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeStreamLivePlans
// This API is used to query the events in the plan in batches.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
func (c *Client) DescribeStreamLivePlans(request *DescribeStreamLivePlansRequest) (response *DescribeStreamLivePlansResponse, err error) {
    return c.DescribeStreamLivePlansWithContext(context.Background(), request)
}

// DescribeStreamLivePlans
// This API is used to query the events in the plan in batches.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
func (c *Client) DescribeStreamLivePlansWithContext(ctx context.Context, request *DescribeStreamLivePlansRequest) (response *DescribeStreamLivePlansResponse, err error) {
    if request == nil {
        request = NewDescribeStreamLivePlansRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeStreamLivePlans require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeStreamLivePlansResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeStreamLiveRegionsRequest() (request *DescribeStreamLiveRegionsRequest) {
    request = &DescribeStreamLiveRegionsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("mdl", APIVersion, "DescribeStreamLiveRegions")
    
    
    return
}

func NewDescribeStreamLiveRegionsResponse() (response *DescribeStreamLiveRegionsResponse) {
    response = &DescribeStreamLiveRegionsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeStreamLiveRegions
// This API is used to query all StreamLive regions.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
func (c *Client) DescribeStreamLiveRegions(request *DescribeStreamLiveRegionsRequest) (response *DescribeStreamLiveRegionsResponse, err error) {
    return c.DescribeStreamLiveRegionsWithContext(context.Background(), request)
}

// DescribeStreamLiveRegions
// This API is used to query all StreamLive regions.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
func (c *Client) DescribeStreamLiveRegionsWithContext(ctx context.Context, request *DescribeStreamLiveRegionsRequest) (response *DescribeStreamLiveRegionsResponse, err error) {
    if request == nil {
        request = NewDescribeStreamLiveRegionsRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeStreamLiveRegions require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeStreamLiveRegionsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeStreamLiveTranscodeDetailRequest() (request *DescribeStreamLiveTranscodeDetailRequest) {
    request = &DescribeStreamLiveTranscodeDetailRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("mdl", APIVersion, "DescribeStreamLiveTranscodeDetail")
    
    
    return
}

func NewDescribeStreamLiveTranscodeDetailResponse() (response *DescribeStreamLiveTranscodeDetailResponse) {
    response = &DescribeStreamLiveTranscodeDetailResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeStreamLiveTranscodeDetail
// This API is used to query the transcoding information of StreamLive streams.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER_CHANNELID = "InvalidParameter.ChannelId"
//  INVALIDPARAMETER_ENDTIME = "InvalidParameter.EndTime"
//  INVALIDPARAMETER_ID = "InvalidParameter.Id"
//  INVALIDPARAMETER_PAGENUM = "InvalidParameter.PageNum"
//  INVALIDPARAMETER_STARTTIME = "InvalidParameter.StartTime"
func (c *Client) DescribeStreamLiveTranscodeDetail(request *DescribeStreamLiveTranscodeDetailRequest) (response *DescribeStreamLiveTranscodeDetailResponse, err error) {
    return c.DescribeStreamLiveTranscodeDetailWithContext(context.Background(), request)
}

// DescribeStreamLiveTranscodeDetail
// This API is used to query the transcoding information of StreamLive streams.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER_CHANNELID = "InvalidParameter.ChannelId"
//  INVALIDPARAMETER_ENDTIME = "InvalidParameter.EndTime"
//  INVALIDPARAMETER_ID = "InvalidParameter.Id"
//  INVALIDPARAMETER_PAGENUM = "InvalidParameter.PageNum"
//  INVALIDPARAMETER_STARTTIME = "InvalidParameter.StartTime"
func (c *Client) DescribeStreamLiveTranscodeDetailWithContext(ctx context.Context, request *DescribeStreamLiveTranscodeDetailRequest) (response *DescribeStreamLiveTranscodeDetailResponse, err error) {
    if request == nil {
        request = NewDescribeStreamLiveTranscodeDetailRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeStreamLiveTranscodeDetail require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeStreamLiveTranscodeDetailResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeStreamLiveWatermarkRequest() (request *DescribeStreamLiveWatermarkRequest) {
    request = &DescribeStreamLiveWatermarkRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("mdl", APIVersion, "DescribeStreamLiveWatermark")
    
    
    return
}

func NewDescribeStreamLiveWatermarkResponse() (response *DescribeStreamLiveWatermarkResponse) {
    response = &DescribeStreamLiveWatermarkResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeStreamLiveWatermark
// This API is used to query a watermark.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER_ID = "InvalidParameter.Id"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
func (c *Client) DescribeStreamLiveWatermark(request *DescribeStreamLiveWatermarkRequest) (response *DescribeStreamLiveWatermarkResponse, err error) {
    return c.DescribeStreamLiveWatermarkWithContext(context.Background(), request)
}

// DescribeStreamLiveWatermark
// This API is used to query a watermark.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER_ID = "InvalidParameter.Id"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
func (c *Client) DescribeStreamLiveWatermarkWithContext(ctx context.Context, request *DescribeStreamLiveWatermarkRequest) (response *DescribeStreamLiveWatermarkResponse, err error) {
    if request == nil {
        request = NewDescribeStreamLiveWatermarkRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeStreamLiveWatermark require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeStreamLiveWatermarkResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeStreamLiveWatermarksRequest() (request *DescribeStreamLiveWatermarksRequest) {
    request = &DescribeStreamLiveWatermarksRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("mdl", APIVersion, "DescribeStreamLiveWatermarks")
    
    
    return
}

func NewDescribeStreamLiveWatermarksResponse() (response *DescribeStreamLiveWatermarksResponse) {
    response = &DescribeStreamLiveWatermarksResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// DescribeStreamLiveWatermarks
// This API is used to query multiple watermarks at a time.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER_ID = "InvalidParameter.Id"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
func (c *Client) DescribeStreamLiveWatermarks(request *DescribeStreamLiveWatermarksRequest) (response *DescribeStreamLiveWatermarksResponse, err error) {
    return c.DescribeStreamLiveWatermarksWithContext(context.Background(), request)
}

// DescribeStreamLiveWatermarks
// This API is used to query multiple watermarks at a time.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER_ID = "InvalidParameter.Id"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
func (c *Client) DescribeStreamLiveWatermarksWithContext(ctx context.Context, request *DescribeStreamLiveWatermarksRequest) (response *DescribeStreamLiveWatermarksResponse, err error) {
    if request == nil {
        request = NewDescribeStreamLiveWatermarksRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeStreamLiveWatermarks require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeStreamLiveWatermarksResponse()
    err = c.Send(request, response)
    return
}

func NewModifyStreamLiveChannelRequest() (request *ModifyStreamLiveChannelRequest) {
    request = &ModifyStreamLiveChannelRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("mdl", APIVersion, "ModifyStreamLiveChannel")
    
    
    return
}

func NewModifyStreamLiveChannelResponse() (response *ModifyStreamLiveChannelResponse) {
    response = &ModifyStreamLiveChannelResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ModifyStreamLiveChannel
// This API is used to modify a StreamLive channel.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_AVTEMPLATES = "InvalidParameter.AVTemplates"
//  INVALIDPARAMETER_ALREADYASSOCIATEDCHANNEL = "InvalidParameter.AlreadyAssociatedChannel"
//  INVALIDPARAMETER_ATTACHEDINPUTS = "InvalidParameter.AttachedInputs"
//  INVALIDPARAMETER_AUDIOTEMPLATES = "InvalidParameter.AudioTemplates"
//  INVALIDPARAMETER_ID = "InvalidParameter.Id"
//  INVALIDPARAMETER_NAME = "InvalidParameter.Name"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
//  INVALIDPARAMETER_OUTPUTGROUPS = "InvalidParameter.OutputGroups"
//  INVALIDPARAMETER_VIDEOTEMPLATES = "InvalidParameter.VideoTemplates"
func (c *Client) ModifyStreamLiveChannel(request *ModifyStreamLiveChannelRequest) (response *ModifyStreamLiveChannelResponse, err error) {
    return c.ModifyStreamLiveChannelWithContext(context.Background(), request)
}

// ModifyStreamLiveChannel
// This API is used to modify a StreamLive channel.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_AVTEMPLATES = "InvalidParameter.AVTemplates"
//  INVALIDPARAMETER_ALREADYASSOCIATEDCHANNEL = "InvalidParameter.AlreadyAssociatedChannel"
//  INVALIDPARAMETER_ATTACHEDINPUTS = "InvalidParameter.AttachedInputs"
//  INVALIDPARAMETER_AUDIOTEMPLATES = "InvalidParameter.AudioTemplates"
//  INVALIDPARAMETER_ID = "InvalidParameter.Id"
//  INVALIDPARAMETER_NAME = "InvalidParameter.Name"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
//  INVALIDPARAMETER_OUTPUTGROUPS = "InvalidParameter.OutputGroups"
//  INVALIDPARAMETER_VIDEOTEMPLATES = "InvalidParameter.VideoTemplates"
func (c *Client) ModifyStreamLiveChannelWithContext(ctx context.Context, request *ModifyStreamLiveChannelRequest) (response *ModifyStreamLiveChannelResponse, err error) {
    if request == nil {
        request = NewModifyStreamLiveChannelRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyStreamLiveChannel require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyStreamLiveChannelResponse()
    err = c.Send(request, response)
    return
}

func NewModifyStreamLiveInputRequest() (request *ModifyStreamLiveInputRequest) {
    request = &ModifyStreamLiveInputRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("mdl", APIVersion, "ModifyStreamLiveInput")
    
    
    return
}

func NewModifyStreamLiveInputResponse() (response *ModifyStreamLiveInputResponse) {
    response = &ModifyStreamLiveInputResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ModifyStreamLiveInput
// This API is used to modify a StreamLive input.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_ID = "InvalidParameter.Id"
//  INVALIDPARAMETER_INPUTSETTINGS = "InvalidParameter.InputSettings"
//  INVALIDPARAMETER_NAME = "InvalidParameter.Name"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
//  INVALIDPARAMETER_SECURITYGROUPS = "InvalidParameter.SecurityGroups"
//  INVALIDPARAMETER_STATEERROR = "InvalidParameter.StateError"
func (c *Client) ModifyStreamLiveInput(request *ModifyStreamLiveInputRequest) (response *ModifyStreamLiveInputResponse, err error) {
    return c.ModifyStreamLiveInputWithContext(context.Background(), request)
}

// ModifyStreamLiveInput
// This API is used to modify a StreamLive input.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_ID = "InvalidParameter.Id"
//  INVALIDPARAMETER_INPUTSETTINGS = "InvalidParameter.InputSettings"
//  INVALIDPARAMETER_NAME = "InvalidParameter.Name"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
//  INVALIDPARAMETER_SECURITYGROUPS = "InvalidParameter.SecurityGroups"
//  INVALIDPARAMETER_STATEERROR = "InvalidParameter.StateError"
func (c *Client) ModifyStreamLiveInputWithContext(ctx context.Context, request *ModifyStreamLiveInputRequest) (response *ModifyStreamLiveInputResponse, err error) {
    if request == nil {
        request = NewModifyStreamLiveInputRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyStreamLiveInput require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyStreamLiveInputResponse()
    err = c.Send(request, response)
    return
}

func NewModifyStreamLiveInputSecurityGroupRequest() (request *ModifyStreamLiveInputSecurityGroupRequest) {
    request = &ModifyStreamLiveInputSecurityGroupRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("mdl", APIVersion, "ModifyStreamLiveInputSecurityGroup")
    
    
    return
}

func NewModifyStreamLiveInputSecurityGroupResponse() (response *ModifyStreamLiveInputSecurityGroupResponse) {
    response = &ModifyStreamLiveInputSecurityGroupResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ModifyStreamLiveInputSecurityGroup
// This API is used to modify an input security group.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_ID = "InvalidParameter.Id"
//  INVALIDPARAMETER_NAME = "InvalidParameter.Name"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
//  INVALIDPARAMETER_WHITELIST = "InvalidParameter.Whitelist"
func (c *Client) ModifyStreamLiveInputSecurityGroup(request *ModifyStreamLiveInputSecurityGroupRequest) (response *ModifyStreamLiveInputSecurityGroupResponse, err error) {
    return c.ModifyStreamLiveInputSecurityGroupWithContext(context.Background(), request)
}

// ModifyStreamLiveInputSecurityGroup
// This API is used to modify an input security group.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_ID = "InvalidParameter.Id"
//  INVALIDPARAMETER_NAME = "InvalidParameter.Name"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
//  INVALIDPARAMETER_WHITELIST = "InvalidParameter.Whitelist"
func (c *Client) ModifyStreamLiveInputSecurityGroupWithContext(ctx context.Context, request *ModifyStreamLiveInputSecurityGroupRequest) (response *ModifyStreamLiveInputSecurityGroupResponse, err error) {
    if request == nil {
        request = NewModifyStreamLiveInputSecurityGroupRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyStreamLiveInputSecurityGroup require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyStreamLiveInputSecurityGroupResponse()
    err = c.Send(request, response)
    return
}

func NewModifyStreamLiveWatermarkRequest() (request *ModifyStreamLiveWatermarkRequest) {
    request = &ModifyStreamLiveWatermarkRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("mdl", APIVersion, "ModifyStreamLiveWatermark")
    
    
    return
}

func NewModifyStreamLiveWatermarkResponse() (response *ModifyStreamLiveWatermarkResponse) {
    response = &ModifyStreamLiveWatermarkResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// ModifyStreamLiveWatermark
// This API is used to modify a watermark.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER_ID = "InvalidParameter.Id"
//  INVALIDPARAMETER_IMAGESETTINGS = "InvalidParameter.ImageSettings"
//  INVALIDPARAMETER_NAME = "InvalidParameter.Name"
//  INVALIDPARAMETER_TEXTSETTINGS = "InvalidParameter.TextSettings"
//  INVALIDPARAMETER_TYPE = "InvalidParameter.Type"
func (c *Client) ModifyStreamLiveWatermark(request *ModifyStreamLiveWatermarkRequest) (response *ModifyStreamLiveWatermarkResponse, err error) {
    return c.ModifyStreamLiveWatermarkWithContext(context.Background(), request)
}

// ModifyStreamLiveWatermark
// This API is used to modify a watermark.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER_ID = "InvalidParameter.Id"
//  INVALIDPARAMETER_IMAGESETTINGS = "InvalidParameter.ImageSettings"
//  INVALIDPARAMETER_NAME = "InvalidParameter.Name"
//  INVALIDPARAMETER_TEXTSETTINGS = "InvalidParameter.TextSettings"
//  INVALIDPARAMETER_TYPE = "InvalidParameter.Type"
func (c *Client) ModifyStreamLiveWatermarkWithContext(ctx context.Context, request *ModifyStreamLiveWatermarkRequest) (response *ModifyStreamLiveWatermarkResponse, err error) {
    if request == nil {
        request = NewModifyStreamLiveWatermarkRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyStreamLiveWatermark require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyStreamLiveWatermarkResponse()
    err = c.Send(request, response)
    return
}

func NewQueryInputStreamStateRequest() (request *QueryInputStreamStateRequest) {
    request = &QueryInputStreamStateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("mdl", APIVersion, "QueryInputStreamState")
    
    
    return
}

func NewQueryInputStreamStateResponse() (response *QueryInputStreamStateResponse) {
    response = &QueryInputStreamStateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// QueryInputStreamState
// This API is used to query the stream status of a StreamLive input.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_ID = "InvalidParameter.Id"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
func (c *Client) QueryInputStreamState(request *QueryInputStreamStateRequest) (response *QueryInputStreamStateResponse, err error) {
    return c.QueryInputStreamStateWithContext(context.Background(), request)
}

// QueryInputStreamState
// This API is used to query the stream status of a StreamLive input.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_ID = "InvalidParameter.Id"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
func (c *Client) QueryInputStreamStateWithContext(ctx context.Context, request *QueryInputStreamStateRequest) (response *QueryInputStreamStateResponse, err error) {
    if request == nil {
        request = NewQueryInputStreamStateRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("QueryInputStreamState require credential")
    }

    request.SetContext(ctx)
    
    response = NewQueryInputStreamStateResponse()
    err = c.Send(request, response)
    return
}

func NewStartStreamLiveChannelRequest() (request *StartStreamLiveChannelRequest) {
    request = &StartStreamLiveChannelRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("mdl", APIVersion, "StartStreamLiveChannel")
    
    
    return
}

func NewStartStreamLiveChannelResponse() (response *StartStreamLiveChannelResponse) {
    response = &StartStreamLiveChannelResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// StartStreamLiveChannel
// This API is used to start a StreamLive channel.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_ID = "InvalidParameter.Id"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
//  INVALIDPARAMETER_STATEERROR = "InvalidParameter.StateError"
func (c *Client) StartStreamLiveChannel(request *StartStreamLiveChannelRequest) (response *StartStreamLiveChannelResponse, err error) {
    return c.StartStreamLiveChannelWithContext(context.Background(), request)
}

// StartStreamLiveChannel
// This API is used to start a StreamLive channel.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_ID = "InvalidParameter.Id"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
//  INVALIDPARAMETER_STATEERROR = "InvalidParameter.StateError"
func (c *Client) StartStreamLiveChannelWithContext(ctx context.Context, request *StartStreamLiveChannelRequest) (response *StartStreamLiveChannelResponse, err error) {
    if request == nil {
        request = NewStartStreamLiveChannelRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("StartStreamLiveChannel require credential")
    }

    request.SetContext(ctx)
    
    response = NewStartStreamLiveChannelResponse()
    err = c.Send(request, response)
    return
}

func NewStopStreamLiveChannelRequest() (request *StopStreamLiveChannelRequest) {
    request = &StopStreamLiveChannelRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("mdl", APIVersion, "StopStreamLiveChannel")
    
    
    return
}

func NewStopStreamLiveChannelResponse() (response *StopStreamLiveChannelResponse) {
    response = &StopStreamLiveChannelResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// StopStreamLiveChannel
// This API is used to stop a StreamLive channel.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_ID = "InvalidParameter.Id"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
//  INVALIDPARAMETER_STATEERROR = "InvalidParameter.StateError"
func (c *Client) StopStreamLiveChannel(request *StopStreamLiveChannelRequest) (response *StopStreamLiveChannelResponse, err error) {
    return c.StopStreamLiveChannelWithContext(context.Background(), request)
}

// StopStreamLiveChannel
// This API is used to stop a StreamLive channel.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_ID = "InvalidParameter.Id"
//  INVALIDPARAMETER_NOTFOUND = "InvalidParameter.NotFound"
//  INVALIDPARAMETER_STATEERROR = "InvalidParameter.StateError"
func (c *Client) StopStreamLiveChannelWithContext(ctx context.Context, request *StopStreamLiveChannelRequest) (response *StopStreamLiveChannelResponse, err error) {
    if request == nil {
        request = NewStopStreamLiveChannelRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("StopStreamLiveChannel require credential")
    }

    request.SetContext(ctx)
    
    response = NewStopStreamLiveChannelResponse()
    err = c.Send(request, response)
    return
}
