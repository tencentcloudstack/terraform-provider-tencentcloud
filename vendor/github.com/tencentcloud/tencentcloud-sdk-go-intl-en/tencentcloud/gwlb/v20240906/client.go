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

import (
    "context"
    "errors"
    "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
    tchttp "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common/http"
    "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common/profile"
)

const APIVersion = "2024-09-06"

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


func NewAssociateTargetGroupsRequest() (request *AssociateTargetGroupsRequest) {
    request = &AssociateTargetGroupsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("gwlb", APIVersion, "AssociateTargetGroups")
    
    
    return
}

func NewAssociateTargetGroupsResponse() (response *AssociateTargetGroupsResponse) {
    response = &AssociateTargetGroupsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// AssociateTargetGroups
// This API is used to bind target groups to a CLB.This is an async API. After the API return succeeds, you can call the DescribeTaskStatus API with the returned RequestID as an input parameter to check whether this task is successful.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_RESOURCEINOPERATING = "FailedOperation.ResourceInOperating"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  LIMITEXCEEDED = "LimitExceeded"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) AssociateTargetGroups(request *AssociateTargetGroupsRequest) (response *AssociateTargetGroupsResponse, err error) {
    return c.AssociateTargetGroupsWithContext(context.Background(), request)
}

// AssociateTargetGroups
// This API is used to bind target groups to a CLB.This is an async API. After the API return succeeds, you can call the DescribeTaskStatus API with the returned RequestID as an input parameter to check whether this task is successful.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_RESOURCEINOPERATING = "FailedOperation.ResourceInOperating"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  LIMITEXCEEDED = "LimitExceeded"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) AssociateTargetGroupsWithContext(ctx context.Context, request *AssociateTargetGroupsRequest) (response *AssociateTargetGroupsResponse, err error) {
    if request == nil {
        request = NewAssociateTargetGroupsRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("AssociateTargetGroups require credential")
    }

    request.SetContext(ctx)
    
    response = NewAssociateTargetGroupsResponse()
    err = c.Send(request, response)
    return
}

func NewCreateGatewayLoadBalancerRequest() (request *CreateGatewayLoadBalancerRequest) {
    request = &CreateGatewayLoadBalancerRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("gwlb", APIVersion, "CreateGatewayLoadBalancer")
    
    
    return
}

func NewCreateGatewayLoadBalancerResponse() (response *CreateGatewayLoadBalancerResponse) {
    response = &CreateGatewayLoadBalancerResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// CreateGatewayLoadBalancer
// This API is used to create a GWLB instance. To use the GWLB service, you must purchase one or more GWLB instances. After this API is called successfully, a unique ID for the GWLB instance will be returned.Note: The default purchase quota for each account in each region is 10.This is an async API. After the API is called successfully, you can use the DescribeGatewayLoadBalancers API to query the status of the GWLB instance (such as creating and normal) to determine whether the creation is successful.
//
// error code that may be returned:
//  INVALIDPARAMETER_FORMATERROR = "InvalidParameter.FormatError"
//  INVALIDPARAMETERVALUE_LENGTH = "InvalidParameterValue.Length"
func (c *Client) CreateGatewayLoadBalancer(request *CreateGatewayLoadBalancerRequest) (response *CreateGatewayLoadBalancerResponse, err error) {
    return c.CreateGatewayLoadBalancerWithContext(context.Background(), request)
}

// CreateGatewayLoadBalancer
// This API is used to create a GWLB instance. To use the GWLB service, you must purchase one or more GWLB instances. After this API is called successfully, a unique ID for the GWLB instance will be returned.Note: The default purchase quota for each account in each region is 10.This is an async API. After the API is called successfully, you can use the DescribeGatewayLoadBalancers API to query the status of the GWLB instance (such as creating and normal) to determine whether the creation is successful.
//
// error code that may be returned:
//  INVALIDPARAMETER_FORMATERROR = "InvalidParameter.FormatError"
//  INVALIDPARAMETERVALUE_LENGTH = "InvalidParameterValue.Length"
func (c *Client) CreateGatewayLoadBalancerWithContext(ctx context.Context, request *CreateGatewayLoadBalancerRequest) (response *CreateGatewayLoadBalancerResponse, err error) {
    if request == nil {
        request = NewCreateGatewayLoadBalancerRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateGatewayLoadBalancer require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateGatewayLoadBalancerResponse()
    err = c.Send(request, response)
    return
}

func NewCreateTargetGroupRequest() (request *CreateTargetGroupRequest) {
    request = &CreateTargetGroupRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("gwlb", APIVersion, "CreateTargetGroup")
    
    
    return
}

func NewCreateTargetGroupResponse() (response *CreateTargetGroupResponse) {
    response = &CreateTargetGroupResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// CreateTargetGroup
// This API is used to create a target group. This feature is in beta testing. If you need to use it, please [submit a ticket](https://console.cloud.tencent.com/workorder/category?level1_id=6&level2_id=163&source=0&data_title=%E8%B4%9F%E8%BD%BD%E5%9D%87%E8%A1%A1%20LB&step=1).
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_FORMATERROR = "InvalidParameter.FormatError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  LIMITEXCEEDED = "LimitExceeded"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CreateTargetGroup(request *CreateTargetGroupRequest) (response *CreateTargetGroupResponse, err error) {
    return c.CreateTargetGroupWithContext(context.Background(), request)
}

// CreateTargetGroup
// This API is used to create a target group. This feature is in beta testing. If you need to use it, please [submit a ticket](https://console.cloud.tencent.com/workorder/category?level1_id=6&level2_id=163&source=0&data_title=%E8%B4%9F%E8%BD%BD%E5%9D%87%E8%A1%A1%20LB&step=1).
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_FORMATERROR = "InvalidParameter.FormatError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  LIMITEXCEEDED = "LimitExceeded"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) CreateTargetGroupWithContext(ctx context.Context, request *CreateTargetGroupRequest) (response *CreateTargetGroupResponse, err error) {
    if request == nil {
        request = NewCreateTargetGroupRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("CreateTargetGroup require credential")
    }

    request.SetContext(ctx)
    
    response = NewCreateTargetGroupResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteGatewayLoadBalancerRequest() (request *DeleteGatewayLoadBalancerRequest) {
    request = &DeleteGatewayLoadBalancerRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("gwlb", APIVersion, "DeleteGatewayLoadBalancer")
    
    
    return
}

func NewDeleteGatewayLoadBalancerResponse() (response *DeleteGatewayLoadBalancerResponse) {
    response = &DeleteGatewayLoadBalancerResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DeleteGatewayLoadBalancer
// This API is used to delete one or more specified GWLB instances. After successful deletion, the GWLB instances will be unbound from the backend service.This is an async API. After the API return succeeds, you can call the DescribeTaskStatus API with the returned RequestId as an input parameter to check whether this task is successful.
//
// error code that may be returned:
//  INVALIDPARAMETER_FORMATERROR = "InvalidParameter.FormatError"
func (c *Client) DeleteGatewayLoadBalancer(request *DeleteGatewayLoadBalancerRequest) (response *DeleteGatewayLoadBalancerResponse, err error) {
    return c.DeleteGatewayLoadBalancerWithContext(context.Background(), request)
}

// DeleteGatewayLoadBalancer
// This API is used to delete one or more specified GWLB instances. After successful deletion, the GWLB instances will be unbound from the backend service.This is an async API. After the API return succeeds, you can call the DescribeTaskStatus API with the returned RequestId as an input parameter to check whether this task is successful.
//
// error code that may be returned:
//  INVALIDPARAMETER_FORMATERROR = "InvalidParameter.FormatError"
func (c *Client) DeleteGatewayLoadBalancerWithContext(ctx context.Context, request *DeleteGatewayLoadBalancerRequest) (response *DeleteGatewayLoadBalancerResponse, err error) {
    if request == nil {
        request = NewDeleteGatewayLoadBalancerRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteGatewayLoadBalancer require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteGatewayLoadBalancerResponse()
    err = c.Send(request, response)
    return
}

func NewDeleteTargetGroupsRequest() (request *DeleteTargetGroupsRequest) {
    request = &DeleteTargetGroupsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("gwlb", APIVersion, "DeleteTargetGroups")
    
    
    return
}

func NewDeleteTargetGroupsResponse() (response *DeleteTargetGroupsResponse) {
    response = &DeleteTargetGroupsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DeleteTargetGroups
// This API is used to delete a target group.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeleteTargetGroups(request *DeleteTargetGroupsRequest) (response *DeleteTargetGroupsResponse, err error) {
    return c.DeleteTargetGroupsWithContext(context.Background(), request)
}

// DeleteTargetGroups
// This API is used to delete a target group.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeleteTargetGroupsWithContext(ctx context.Context, request *DeleteTargetGroupsRequest) (response *DeleteTargetGroupsResponse, err error) {
    if request == nil {
        request = NewDeleteTargetGroupsRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeleteTargetGroups require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeleteTargetGroupsResponse()
    err = c.Send(request, response)
    return
}

func NewDeregisterTargetGroupInstancesRequest() (request *DeregisterTargetGroupInstancesRequest) {
    request = &DeregisterTargetGroupInstancesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("gwlb", APIVersion, "DeregisterTargetGroupInstances")
    
    
    return
}

func NewDeregisterTargetGroupInstancesResponse() (response *DeregisterTargetGroupInstancesResponse) {
    response = &DeregisterTargetGroupInstancesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DeregisterTargetGroupInstances
// This API is used to unbind a server from a target group.
//
// This is an async API. After the API return succeeds, you can call the DescribeTaskStatus API with the returned RequestID as an input parameter to check whether this task is successful.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_RESOURCEINOPERATING = "FailedOperation.ResourceInOperating"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeregisterTargetGroupInstances(request *DeregisterTargetGroupInstancesRequest) (response *DeregisterTargetGroupInstancesResponse, err error) {
    return c.DeregisterTargetGroupInstancesWithContext(context.Background(), request)
}

// DeregisterTargetGroupInstances
// This API is used to unbind a server from a target group.
//
// This is an async API. After the API return succeeds, you can call the DescribeTaskStatus API with the returned RequestID as an input parameter to check whether this task is successful.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_RESOURCEINOPERATING = "FailedOperation.ResourceInOperating"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DeregisterTargetGroupInstancesWithContext(ctx context.Context, request *DeregisterTargetGroupInstancesRequest) (response *DeregisterTargetGroupInstancesResponse, err error) {
    if request == nil {
        request = NewDeregisterTargetGroupInstancesRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DeregisterTargetGroupInstances require credential")
    }

    request.SetContext(ctx)
    
    response = NewDeregisterTargetGroupInstancesResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeGatewayLoadBalancersRequest() (request *DescribeGatewayLoadBalancersRequest) {
    request = &DescribeGatewayLoadBalancersRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("gwlb", APIVersion, "DescribeGatewayLoadBalancers")
    
    
    return
}

func NewDescribeGatewayLoadBalancersResponse() (response *DescribeGatewayLoadBalancersResponse) {
    response = &DescribeGatewayLoadBalancersResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeGatewayLoadBalancers
// This API is used to query the list of GWLB instances in a region.
//
// error code that may be returned:
//  INVALIDPARAMETER_FORMATERROR = "InvalidParameter.FormatError"
//  INVALIDPARAMETERVALUE_LENGTH = "InvalidParameterValue.Length"
func (c *Client) DescribeGatewayLoadBalancers(request *DescribeGatewayLoadBalancersRequest) (response *DescribeGatewayLoadBalancersResponse, err error) {
    return c.DescribeGatewayLoadBalancersWithContext(context.Background(), request)
}

// DescribeGatewayLoadBalancers
// This API is used to query the list of GWLB instances in a region.
//
// error code that may be returned:
//  INVALIDPARAMETER_FORMATERROR = "InvalidParameter.FormatError"
//  INVALIDPARAMETERVALUE_LENGTH = "InvalidParameterValue.Length"
func (c *Client) DescribeGatewayLoadBalancersWithContext(ctx context.Context, request *DescribeGatewayLoadBalancersRequest) (response *DescribeGatewayLoadBalancersResponse, err error) {
    if request == nil {
        request = NewDescribeGatewayLoadBalancersRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeGatewayLoadBalancers require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeGatewayLoadBalancersResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeTargetGroupInstanceStatusRequest() (request *DescribeTargetGroupInstanceStatusRequest) {
    request = &DescribeTargetGroupInstanceStatusRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("gwlb", APIVersion, "DescribeTargetGroupInstanceStatus")
    
    
    return
}

func NewDescribeTargetGroupInstanceStatusResponse() (response *DescribeTargetGroupInstanceStatusResponse) {
    response = &DescribeTargetGroupInstanceStatusResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeTargetGroupInstanceStatus
// This API is used to query the backend service status of a target group. Currently, only GWLB type target groups support querying backend service status.
//
// error code that may be returned:
//  INVALIDPARAMETER_FORMATERROR = "InvalidParameter.FormatError"
//  INVALIDPARAMETERVALUE_LENGTH = "InvalidParameterValue.Length"
func (c *Client) DescribeTargetGroupInstanceStatus(request *DescribeTargetGroupInstanceStatusRequest) (response *DescribeTargetGroupInstanceStatusResponse, err error) {
    return c.DescribeTargetGroupInstanceStatusWithContext(context.Background(), request)
}

// DescribeTargetGroupInstanceStatus
// This API is used to query the backend service status of a target group. Currently, only GWLB type target groups support querying backend service status.
//
// error code that may be returned:
//  INVALIDPARAMETER_FORMATERROR = "InvalidParameter.FormatError"
//  INVALIDPARAMETERVALUE_LENGTH = "InvalidParameterValue.Length"
func (c *Client) DescribeTargetGroupInstanceStatusWithContext(ctx context.Context, request *DescribeTargetGroupInstanceStatusRequest) (response *DescribeTargetGroupInstanceStatusResponse, err error) {
    if request == nil {
        request = NewDescribeTargetGroupInstanceStatusRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeTargetGroupInstanceStatus require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeTargetGroupInstanceStatusResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeTargetGroupInstancesRequest() (request *DescribeTargetGroupInstancesRequest) {
    request = &DescribeTargetGroupInstancesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("gwlb", APIVersion, "DescribeTargetGroupInstances")
    
    
    return
}

func NewDescribeTargetGroupInstancesResponse() (response *DescribeTargetGroupInstancesResponse) {
    response = &DescribeTargetGroupInstancesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeTargetGroupInstances
// This API is used to obtain information on servers bound to a target group.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeTargetGroupInstances(request *DescribeTargetGroupInstancesRequest) (response *DescribeTargetGroupInstancesResponse, err error) {
    return c.DescribeTargetGroupInstancesWithContext(context.Background(), request)
}

// DescribeTargetGroupInstances
// This API is used to obtain information on servers bound to a target group.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeTargetGroupInstancesWithContext(ctx context.Context, request *DescribeTargetGroupInstancesRequest) (response *DescribeTargetGroupInstancesResponse, err error) {
    if request == nil {
        request = NewDescribeTargetGroupInstancesRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeTargetGroupInstances require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeTargetGroupInstancesResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeTargetGroupListRequest() (request *DescribeTargetGroupListRequest) {
    request = &DescribeTargetGroupListRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("gwlb", APIVersion, "DescribeTargetGroupList")
    
    
    return
}

func NewDescribeTargetGroupListResponse() (response *DescribeTargetGroupListResponse) {
    response = &DescribeTargetGroupListResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeTargetGroupList
// This API is used to obtain a target group list.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_FORMATERROR = "InvalidParameter.FormatError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeTargetGroupList(request *DescribeTargetGroupListRequest) (response *DescribeTargetGroupListResponse, err error) {
    return c.DescribeTargetGroupListWithContext(context.Background(), request)
}

// DescribeTargetGroupList
// This API is used to obtain a target group list.
//
// error code that may be returned:
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_FORMATERROR = "InvalidParameter.FormatError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeTargetGroupListWithContext(ctx context.Context, request *DescribeTargetGroupListRequest) (response *DescribeTargetGroupListResponse, err error) {
    if request == nil {
        request = NewDescribeTargetGroupListRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeTargetGroupList require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeTargetGroupListResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeTargetGroupsRequest() (request *DescribeTargetGroupsRequest) {
    request = &DescribeTargetGroupsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("gwlb", APIVersion, "DescribeTargetGroups")
    
    
    return
}

func NewDescribeTargetGroupsResponse() (response *DescribeTargetGroupsResponse) {
    response = &DescribeTargetGroupsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeTargetGroups
// This API is used to query target group information.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeTargetGroups(request *DescribeTargetGroupsRequest) (response *DescribeTargetGroupsResponse, err error) {
    return c.DescribeTargetGroupsWithContext(context.Background(), request)
}

// DescribeTargetGroups
// This API is used to query target group information.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DescribeTargetGroupsWithContext(ctx context.Context, request *DescribeTargetGroupsRequest) (response *DescribeTargetGroupsResponse, err error) {
    if request == nil {
        request = NewDescribeTargetGroupsRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeTargetGroups require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeTargetGroupsResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeTaskStatusRequest() (request *DescribeTaskStatusRequest) {
    request = &DescribeTaskStatusRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("gwlb", APIVersion, "DescribeTaskStatus")
    
    
    return
}

func NewDescribeTaskStatusResponse() (response *DescribeTaskStatusResponse) {
    response = &DescribeTaskStatusResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DescribeTaskStatus
// This API is used to query the execution status of an async task. After non-query APIs (for example, used to create/delete CLB instances) are called successfully, this API needs to be used to query whether the task is successful.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_REGIONNOTFOUND = "InvalidParameter.RegionNotFound"
func (c *Client) DescribeTaskStatus(request *DescribeTaskStatusRequest) (response *DescribeTaskStatusResponse, err error) {
    return c.DescribeTaskStatusWithContext(context.Background(), request)
}

// DescribeTaskStatus
// This API is used to query the execution status of an async task. After non-query APIs (for example, used to create/delete CLB instances) are called successfully, this API needs to be used to query whether the task is successful.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETER_REGIONNOTFOUND = "InvalidParameter.RegionNotFound"
func (c *Client) DescribeTaskStatusWithContext(ctx context.Context, request *DescribeTaskStatusRequest) (response *DescribeTaskStatusResponse, err error) {
    if request == nil {
        request = NewDescribeTaskStatusRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DescribeTaskStatus require credential")
    }

    request.SetContext(ctx)
    
    response = NewDescribeTaskStatusResponse()
    err = c.Send(request, response)
    return
}

func NewDisassociateTargetGroupsRequest() (request *DisassociateTargetGroupsRequest) {
    request = &DisassociateTargetGroupsRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("gwlb", APIVersion, "DisassociateTargetGroups")
    
    
    return
}

func NewDisassociateTargetGroupsResponse() (response *DisassociateTargetGroupsResponse) {
    response = &DisassociateTargetGroupsResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// DisassociateTargetGroups
// This API is used to disassociate a CLB from a target group.This is an async API. After the API return succeeds, you can call the DescribeTaskStatus API with the returned RequestID as an input parameter to check whether this task is successful.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_RESOURCEINOPERATING = "FailedOperation.ResourceInOperating"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DisassociateTargetGroups(request *DisassociateTargetGroupsRequest) (response *DisassociateTargetGroupsResponse, err error) {
    return c.DisassociateTargetGroupsWithContext(context.Background(), request)
}

// DisassociateTargetGroups
// This API is used to disassociate a CLB from a target group.This is an async API. After the API return succeeds, you can call the DescribeTaskStatus API with the returned RequestID as an input parameter to check whether this task is successful.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_RESOURCEINOPERATING = "FailedOperation.ResourceInOperating"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) DisassociateTargetGroupsWithContext(ctx context.Context, request *DisassociateTargetGroupsRequest) (response *DisassociateTargetGroupsResponse, err error) {
    if request == nil {
        request = NewDisassociateTargetGroupsRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("DisassociateTargetGroups require credential")
    }

    request.SetContext(ctx)
    
    response = NewDisassociateTargetGroupsResponse()
    err = c.Send(request, response)
    return
}

func NewInquirePriceCreateGatewayLoadBalancerRequest() (request *InquirePriceCreateGatewayLoadBalancerRequest) {
    request = &InquirePriceCreateGatewayLoadBalancerRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("gwlb", APIVersion, "InquirePriceCreateGatewayLoadBalancer")
    
    
    return
}

func NewInquirePriceCreateGatewayLoadBalancerResponse() (response *InquirePriceCreateGatewayLoadBalancerResponse) {
    response = &InquirePriceCreateGatewayLoadBalancerResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// InquirePriceCreateGatewayLoadBalancer
// This API is used to query the price for creating a GWLB.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_RESOURCEINOPERATING = "FailedOperation.ResourceInOperating"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) InquirePriceCreateGatewayLoadBalancer(request *InquirePriceCreateGatewayLoadBalancerRequest) (response *InquirePriceCreateGatewayLoadBalancerResponse, err error) {
    return c.InquirePriceCreateGatewayLoadBalancerWithContext(context.Background(), request)
}

// InquirePriceCreateGatewayLoadBalancer
// This API is used to query the price for creating a GWLB.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_RESOURCEINOPERATING = "FailedOperation.ResourceInOperating"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) InquirePriceCreateGatewayLoadBalancerWithContext(ctx context.Context, request *InquirePriceCreateGatewayLoadBalancerRequest) (response *InquirePriceCreateGatewayLoadBalancerResponse, err error) {
    if request == nil {
        request = NewInquirePriceCreateGatewayLoadBalancerRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("InquirePriceCreateGatewayLoadBalancer require credential")
    }

    request.SetContext(ctx)
    
    response = NewInquirePriceCreateGatewayLoadBalancerResponse()
    err = c.Send(request, response)
    return
}

func NewModifyGatewayLoadBalancerAttributeRequest() (request *ModifyGatewayLoadBalancerAttributeRequest) {
    request = &ModifyGatewayLoadBalancerAttributeRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("gwlb", APIVersion, "ModifyGatewayLoadBalancerAttribute")
    
    
    return
}

func NewModifyGatewayLoadBalancerAttributeResponse() (response *ModifyGatewayLoadBalancerAttributeResponse) {
    response = &ModifyGatewayLoadBalancerAttributeResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// ModifyGatewayLoadBalancerAttribute
// This API is used to modify the attributes of a CLB instance. It supports modifying the name and bandwidth cap of the CLB instance.
//
// error code that may be returned:
//  INVALIDPARAMETER_FORMATERROR = "InvalidParameter.FormatError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_DUPLICATE = "InvalidParameterValue.Duplicate"
func (c *Client) ModifyGatewayLoadBalancerAttribute(request *ModifyGatewayLoadBalancerAttributeRequest) (response *ModifyGatewayLoadBalancerAttributeResponse, err error) {
    return c.ModifyGatewayLoadBalancerAttributeWithContext(context.Background(), request)
}

// ModifyGatewayLoadBalancerAttribute
// This API is used to modify the attributes of a CLB instance. It supports modifying the name and bandwidth cap of the CLB instance.
//
// error code that may be returned:
//  INVALIDPARAMETER_FORMATERROR = "InvalidParameter.FormatError"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  INVALIDPARAMETERVALUE_DUPLICATE = "InvalidParameterValue.Duplicate"
func (c *Client) ModifyGatewayLoadBalancerAttributeWithContext(ctx context.Context, request *ModifyGatewayLoadBalancerAttributeRequest) (response *ModifyGatewayLoadBalancerAttributeResponse, err error) {
    if request == nil {
        request = NewModifyGatewayLoadBalancerAttributeRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyGatewayLoadBalancerAttribute require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyGatewayLoadBalancerAttributeResponse()
    err = c.Send(request, response)
    return
}

func NewModifyTargetGroupAttributeRequest() (request *ModifyTargetGroupAttributeRequest) {
    request = &ModifyTargetGroupAttributeRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("gwlb", APIVersion, "ModifyTargetGroupAttribute")
    
    
    return
}

func NewModifyTargetGroupAttributeResponse() (response *ModifyTargetGroupAttributeResponse) {
    response = &ModifyTargetGroupAttributeResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// ModifyTargetGroupAttribute
// This API is used to modify the name, health check, and other attributes of the target group.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifyTargetGroupAttribute(request *ModifyTargetGroupAttributeRequest) (response *ModifyTargetGroupAttributeResponse, err error) {
    return c.ModifyTargetGroupAttributeWithContext(context.Background(), request)
}

// ModifyTargetGroupAttribute
// This API is used to modify the name, health check, and other attributes of the target group.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifyTargetGroupAttributeWithContext(ctx context.Context, request *ModifyTargetGroupAttributeRequest) (response *ModifyTargetGroupAttributeResponse, err error) {
    if request == nil {
        request = NewModifyTargetGroupAttributeRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyTargetGroupAttribute require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyTargetGroupAttributeResponse()
    err = c.Send(request, response)
    return
}

func NewModifyTargetGroupInstancesWeightRequest() (request *ModifyTargetGroupInstancesWeightRequest) {
    request = &ModifyTargetGroupInstancesWeightRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("gwlb", APIVersion, "ModifyTargetGroupInstancesWeight")
    
    
    return
}

func NewModifyTargetGroupInstancesWeightResponse() (response *ModifyTargetGroupInstancesWeightResponse) {
    response = &ModifyTargetGroupInstancesWeightResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// ModifyTargetGroupInstancesWeight
// This API is used to modify the server weight of a target group.This is an async API. After the API return succeeds, you can call the DescribeTaskStatus API with the returned RequestID as an input parameter to check whether this task is successful.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_RESOURCEINOPERATING = "FailedOperation.ResourceInOperating"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifyTargetGroupInstancesWeight(request *ModifyTargetGroupInstancesWeightRequest) (response *ModifyTargetGroupInstancesWeightResponse, err error) {
    return c.ModifyTargetGroupInstancesWeightWithContext(context.Background(), request)
}

// ModifyTargetGroupInstancesWeight
// This API is used to modify the server weight of a target group.This is an async API. After the API return succeeds, you can call the DescribeTaskStatus API with the returned RequestID as an input parameter to check whether this task is successful.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_RESOURCEINOPERATING = "FailedOperation.ResourceInOperating"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) ModifyTargetGroupInstancesWeightWithContext(ctx context.Context, request *ModifyTargetGroupInstancesWeightRequest) (response *ModifyTargetGroupInstancesWeightResponse, err error) {
    if request == nil {
        request = NewModifyTargetGroupInstancesWeightRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("ModifyTargetGroupInstancesWeight require credential")
    }

    request.SetContext(ctx)
    
    response = NewModifyTargetGroupInstancesWeightResponse()
    err = c.Send(request, response)
    return
}

func NewRegisterTargetGroupInstancesRequest() (request *RegisterTargetGroupInstancesRequest) {
    request = &RegisterTargetGroupInstancesRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    
    request.Init().WithApiInfo("gwlb", APIVersion, "RegisterTargetGroupInstances")
    
    
    return
}

func NewRegisterTargetGroupInstancesResponse() (response *RegisterTargetGroupInstancesResponse) {
    response = &RegisterTargetGroupInstancesResponse{
        BaseResponse: &tchttp.BaseResponse{},
    } 
    return

}

// RegisterTargetGroupInstances
// This API is used to register servers to a target group.This is an async API. After the API return succeeds, you can call the DescribeTaskStatus API with the returned RequestID as an input parameter to check whether this task is successful.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_RESOURCEINOPERATING = "FailedOperation.ResourceInOperating"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  LIMITEXCEEDED = "LimitExceeded"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) RegisterTargetGroupInstances(request *RegisterTargetGroupInstancesRequest) (response *RegisterTargetGroupInstancesResponse, err error) {
    return c.RegisterTargetGroupInstancesWithContext(context.Background(), request)
}

// RegisterTargetGroupInstances
// This API is used to register servers to a target group.This is an async API. After the API return succeeds, you can call the DescribeTaskStatus API with the returned RequestID as an input parameter to check whether this task is successful.
//
// error code that may be returned:
//  FAILEDOPERATION = "FailedOperation"
//  FAILEDOPERATION_RESOURCEINOPERATING = "FailedOperation.ResourceInOperating"
//  INTERNALERROR = "InternalError"
//  INVALIDPARAMETER = "InvalidParameter"
//  INVALIDPARAMETERVALUE = "InvalidParameterValue"
//  LIMITEXCEEDED = "LimitExceeded"
//  UNAUTHORIZEDOPERATION = "UnauthorizedOperation"
func (c *Client) RegisterTargetGroupInstancesWithContext(ctx context.Context, request *RegisterTargetGroupInstancesRequest) (response *RegisterTargetGroupInstancesResponse, err error) {
    if request == nil {
        request = NewRegisterTargetGroupInstancesRequest()
    }
    
    if c.GetCredential() == nil {
        return nil, errors.New("RegisterTargetGroupInstances require credential")
    }

    request.SetContext(ctx)
    
    response = NewRegisterTargetGroupInstancesResponse()
    err = c.Send(request, response)
    return
}
