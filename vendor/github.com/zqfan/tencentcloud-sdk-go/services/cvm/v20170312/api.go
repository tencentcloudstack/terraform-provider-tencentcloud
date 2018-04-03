package cvm

import (
	"github.com/zqfan/tencentcloud-sdk-go/common"
)

const APIVersion = "2017-03-12"

func NewDescribeAddressesRequest() (request *DescribeAddressesRequest) {
	request = &DescribeAddressesRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("eip", APIVersion, "DescribeAddresses")
	return
}

func NewDescribeAddressesResponse() (response *DescribeAddressesResponse) {
	response = &DescribeAddressesResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeAddresses(request *DescribeAddressesRequest) (response *DescribeAddressesResponse, err error) {
	if request == nil {
		request = NewDescribeAddressesRequest()
	}
	response = NewDescribeAddressesResponse()
	err = c.Send(request, response)
	return
}

func NewReleaseAddressesRequest() (request *ReleaseAddressesRequest) {
	request = &ReleaseAddressesRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("eip", APIVersion, "ReleaseAddresses")
	return
}

func NewReleaseAddressesResponse() (response *ReleaseAddressesResponse) {
	response = &ReleaseAddressesResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) ReleaseAddresses(request *ReleaseAddressesRequest) (response *ReleaseAddressesResponse, err error) {
	if request == nil {
		request = NewReleaseAddressesRequest()
	}
	response = NewReleaseAddressesResponse()
	err = c.Send(request, response)
	return
}

func NewModifyAddressAttributeRequest() (request *ModifyAddressAttributeRequest) {
	request = &ModifyAddressAttributeRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("eip", APIVersion, "ModifyAddressAttribute")
	return
}

func NewModifyAddressAttributeResponse() (response *ModifyAddressAttributeResponse) {
	response = &ModifyAddressAttributeResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) ModifyAddressAttribute(request *ModifyAddressAttributeRequest) (response *ModifyAddressAttributeResponse, err error) {
	if request == nil {
		request = NewModifyAddressAttributeRequest()
	}
	response = NewModifyAddressAttributeResponse()
	err = c.Send(request, response)
	return
}

func NewAllocateAddressesRequest() (request *AllocateAddressesRequest) {
	request = &AllocateAddressesRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("eip", APIVersion, "AllocateAddresses")
	return
}

func NewAllocateAddressesResponse() (response *AllocateAddressesResponse) {
	response = &AllocateAddressesResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) AllocateAddresses(request *AllocateAddressesRequest) (response *AllocateAddressesResponse, err error) {
	if request == nil {
		request = NewAllocateAddressesRequest()
	}
	response = NewAllocateAddressesResponse()
	err = c.Send(request, response)
	return
}

func NewAssociateAddressRequest() (request *AssociateAddressRequest) {
	request = &AssociateAddressRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("eip", APIVersion, "AssociateAddress")
	return
}

func NewAssociateAddressResponse() (response *AssociateAddressResponse) {
	response = &AssociateAddressResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) AssociateAddress(request *AssociateAddressRequest) (response *AssociateAddressResponse, err error) {
	if request == nil {
		request = NewAssociateAddressRequest()
	}
	response = NewAssociateAddressResponse()
	err = c.Send(request, response)
	return
}

func NewDisassociateAddressRequest() (request *DisassociateAddressRequest) {
	request = &DisassociateAddressRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("eip", APIVersion, "DisassociateAddress")
	return
}

func NewDisassociateAddressResponse() (response *DisassociateAddressResponse) {
	response = &DisassociateAddressResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DisassociateAddress(request *DisassociateAddressRequest) (response *DisassociateAddressResponse, err error) {
	if request == nil {
		request = NewDisassociateAddressRequest()
	}
	response = NewDisassociateAddressResponse()
	err = c.Send(request, response)
	return
}

func NewRunInstancesRequest() (request *RunInstancesRequest) {
	request = &RunInstancesRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("cvm", APIVersion, "RunInstances")
	return
}

func NewRunInstancesResponse() (response *RunInstancesResponse) {
	response = &RunInstancesResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) RunInstances(request *RunInstancesRequest) (response *RunInstancesResponse, err error) {
	if request == nil {
		request = NewRunInstancesRequest()
	}
	response = NewRunInstancesResponse()
	err = c.Send(request, response)
	return
}

func NewDescribeInstancesRequest() (request *DescribeInstancesRequest) {
	request = &DescribeInstancesRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("cvm", APIVersion, "DescribeInstances")
	return
}

func NewDescribeInstancesResponse() (response *DescribeInstancesResponse) {
	response = &DescribeInstancesResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeInstances(request *DescribeInstancesRequest) (response *DescribeInstancesResponse, err error) {
	if request == nil {
		request = NewDescribeInstancesRequest()
	}
	response = NewDescribeInstancesResponse()
	err = c.Send(request, response)
	return
}

func NewTerminateInstancesRequest() (request *TerminateInstancesRequest) {
	request = &TerminateInstancesRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("cvm", APIVersion, "TerminateInstances")
	return
}

func NewTerminateInstancesResponse() (response *TerminateInstancesResponse) {
	response = &TerminateInstancesResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) TerminateInstances(request *TerminateInstancesRequest) (response *TerminateInstancesResponse, err error) {
	if request == nil {
		request = NewTerminateInstancesRequest()
	}
	response = NewTerminateInstancesResponse()
	err = c.Send(request, response)
	return
}
