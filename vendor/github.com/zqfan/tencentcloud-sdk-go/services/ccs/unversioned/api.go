package ccs

import (
	"github.com/zqfan/tencentcloud-sdk-go/common"
)

const APIVersion = ""

func NewAddClusterInstancesRequest() (request *AddClusterInstancesRequest) {
	request = &AddClusterInstancesRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("ccs", APIVersion, "AddClusterInstances")
	return
}

func NewAddClusterInstancesResponse() (response *AddClusterInstancesResponse) {
	response = &AddClusterInstancesResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) AddClusterInstances(request *AddClusterInstancesRequest) (response *AddClusterInstancesResponse, err error) {
	if request == nil {
		request = NewAddClusterInstancesRequest()
	}
	response = NewAddClusterInstancesResponse()
	err = c.Send(request, response)
	return
}

func NewAddClusterInstancesFromExistedCvmRequest() (request *AddClusterInstancesFromExistedCvmRequest) {
	request = &AddClusterInstancesFromExistedCvmRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("ccs", APIVersion, "AddClusterInstancesFromExistedCvm")
	return
}

func NewAddClusterInstancesFromExistedCvmResponse() (response *AddClusterInstancesFromExistedCvmResponse) {
	response = &AddClusterInstancesFromExistedCvmResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) AddClusterInstancesFromExistedCvm(request *AddClusterInstancesFromExistedCvmRequest) (response *AddClusterInstancesFromExistedCvmResponse, err error) {
	if request == nil {
		request = NewAddClusterInstancesFromExistedCvmRequest()
	}
	response = NewAddClusterInstancesFromExistedCvmResponse()
	err = c.Send(request, response)
	return
}

func NewCreateClusterRequest() (request *CreateClusterRequest) {
	request = &CreateClusterRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("ccs", APIVersion, "CreateCluster")
	return
}

func NewCreateClusterResponse() (response *CreateClusterResponse) {
	response = &CreateClusterResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) CreateCluster(request *CreateClusterRequest) (response *CreateClusterResponse, err error) {
	if request == nil {
		request = NewCreateClusterRequest()
	}
	response = NewCreateClusterResponse()
	err = c.Send(request, response)
	return
}

func NewDeleteClusterRequest() (request *DeleteClusterRequest) {
	request = &DeleteClusterRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("ccs", APIVersion, "DeleteCluster")
	return
}

func NewDeleteClusterResponse() (response *DeleteClusterResponse) {
	response = &DeleteClusterResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DeleteCluster(request *DeleteClusterRequest) (response *DeleteClusterResponse, err error) {
	if request == nil {
		request = NewDeleteClusterRequest()
	}
	response = NewDeleteClusterResponse()
	err = c.Send(request, response)
	return
}

func NewDeleteClusterInstancesRequest() (request *DeleteClusterInstancesRequest) {
	request = &DeleteClusterInstancesRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("ccs", APIVersion, "DeleteClusterInstances")
	return
}

func NewDeleteClusterInstancesResponse() (response *DeleteClusterInstancesResponse) {
	response = &DeleteClusterInstancesResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DeleteClusterInstances(request *DeleteClusterInstancesRequest) (response *DeleteClusterInstancesResponse, err error) {
	if request == nil {
		request = NewDeleteClusterInstancesRequest()
	}
	response = NewDeleteClusterInstancesResponse()
	err = c.Send(request, response)
	return
}

func NewDescribeClusterRequest() (request *DescribeClusterRequest) {
	request = &DescribeClusterRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("ccs", APIVersion, "DescribeCluster")
	return
}

func NewDescribeClusterResponse() (response *DescribeClusterResponse) {
	response = &DescribeClusterResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeCluster(request *DescribeClusterRequest) (response *DescribeClusterResponse, err error) {
	if request == nil {
		request = NewDescribeClusterRequest()
	}
	response = NewDescribeClusterResponse()
	err = c.Send(request, response)
	return
}

func NewDescribeClusterInstancesRequest() (request *DescribeClusterInstancesRequest) {
	request = &DescribeClusterInstancesRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("ccs", APIVersion, "DescribeClusterInstances")
	return
}

func NewDescribeClusterInstancesResponse() (response *DescribeClusterInstancesResponse) {
	response = &DescribeClusterInstancesResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeClusterInstances(request *DescribeClusterInstancesRequest) (response *DescribeClusterInstancesResponse, err error) {
	if request == nil {
		request = NewDescribeClusterInstancesRequest()
	}
	response = NewDescribeClusterInstancesResponse()
	err = c.Send(request, response)
	return
}

func NewDescribeClusterSecurityInfoRequest() (request *DescribeClusterSecurityInfoRequest) {
	request = &DescribeClusterSecurityInfoRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("ccs", APIVersion, "DescribeClusterSecurityInfo")
	return
}

func NewDescribeClusterSecurityInfoResponse() (response *DescribeClusterSecurityInfoResponse) {
	response = &DescribeClusterSecurityInfoResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeClusterSecurityInfo(request *DescribeClusterSecurityInfoRequest) (response *DescribeClusterSecurityInfoResponse, err error) {
	if request == nil {
		request = NewDescribeClusterSecurityInfoRequest()
	}
	response = NewDescribeClusterSecurityInfoResponse()
	err = c.Send(request, response)
	return
}

func NewDescribeClusterTaskResultRequest() (request *DescribeClusterTaskResultRequest) {
	request = &DescribeClusterTaskResultRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("ccs", APIVersion, "DescribeClusterTaskResult")
	return
}

func NewDescribeClusterTaskResultResponse() (response *DescribeClusterTaskResultResponse) {
	response = &DescribeClusterTaskResultResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeClusterTaskResult(request *DescribeClusterTaskResultRequest) (response *DescribeClusterTaskResultResponse, err error) {
	if request == nil {
		request = NewDescribeClusterTaskResultRequest()
	}
	response = NewDescribeClusterTaskResultResponse()
	err = c.Send(request, response)
	return
}

func NewOperateClusterVipRequest() (request *OperateClusterVipRequest) {
	request = &OperateClusterVipRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("ccs", APIVersion, "OperateClusterVip")
	return
}

func NewOperateClusterVipResponse() (response *OperateClusterVipResponse) {
	response = &OperateClusterVipResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) OperateClusterVip(request *OperateClusterVipRequest) (response *OperateClusterVipResponse, err error) {
	if request == nil {
		request = NewOperateClusterVipRequest()
	}
	response = NewOperateClusterVipResponse()
	err = c.Send(request, response)
	return
}
