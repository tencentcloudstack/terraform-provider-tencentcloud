package unversioned

import (
	"github.com/zqfan/tencentcloud-sdk-go/common"
)

const APIVersion = ""

func NewCreateLoadBalancerRequest() (request *CreateLoadBalancerRequest) {
	request = &CreateLoadBalancerRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("lb", APIVersion, "CreateLoadBalancer")
	return
}

func NewCreateLoadBalancerResponse() (response *CreateLoadBalancerResponse) {
	response = &CreateLoadBalancerResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) CreateLoadBalancer(request *CreateLoadBalancerRequest) (response *CreateLoadBalancerResponse, err error) {
	if request == nil {
		request = NewCreateLoadBalancerRequest()
	}
	response = NewCreateLoadBalancerResponse()
	err = c.Send(request, response)
	return
}

func NewDeleteLoadBalancersRequest() (request *DeleteLoadBalancersRequest) {
	request = &DeleteLoadBalancersRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("lb", APIVersion, "DeleteLoadBalancers")
	return
}

func NewDeleteLoadBalancersResponse() (response *DeleteLoadBalancersResponse) {
	response = &DeleteLoadBalancersResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DeleteLoadBalancers(request *DeleteLoadBalancersRequest) (response *DeleteLoadBalancersResponse, err error) {
	if request == nil {
		request = NewDeleteLoadBalancersRequest()
	}
	response = NewDeleteLoadBalancersResponse()
	err = c.Send(request, response)
	return
}

func NewDescribeForwardLBBackendsRequest() (request *DescribeForwardLBBackendsRequest) {
	request = &DescribeForwardLBBackendsRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("lb", APIVersion, "DescribeForwardLBBackends")
	return
}

func NewDescribeForwardLBBackendsResponse() (response *DescribeForwardLBBackendsResponse) {
	response = &DescribeForwardLBBackendsResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeForwardLBBackends(request *DescribeForwardLBBackendsRequest) (response *DescribeForwardLBBackendsResponse, err error) {
	if request == nil {
		request = NewDescribeForwardLBBackendsRequest()
	}
	response = NewDescribeForwardLBBackendsResponse()
	err = c.Send(request, response)
	return
}

func NewDescribeForwardLBListenersRequest() (request *DescribeForwardLBListenersRequest) {
	request = &DescribeForwardLBListenersRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("lb", APIVersion, "DescribeForwardLBListeners")
	return
}

func NewDescribeForwardLBListenersResponse() (response *DescribeForwardLBListenersResponse) {
	response = &DescribeForwardLBListenersResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeForwardLBListeners(request *DescribeForwardLBListenersRequest) (response *DescribeForwardLBListenersResponse, err error) {
	if request == nil {
		request = NewDescribeForwardLBListenersRequest()
	}
	response = NewDescribeForwardLBListenersResponse()
	err = c.Send(request, response)
	return
}

func NewDescribeLoadBalancersRequest() (request *DescribeLoadBalancersRequest) {
	request = &DescribeLoadBalancersRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("lb", APIVersion, "DescribeLoadBalancers")
	return
}

func NewDescribeLoadBalancersResponse() (response *DescribeLoadBalancersResponse) {
	response = &DescribeLoadBalancersResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeLoadBalancers(request *DescribeLoadBalancersRequest) (response *DescribeLoadBalancersResponse, err error) {
	if request == nil {
		request = NewDescribeLoadBalancersRequest()
	}
	response = NewDescribeLoadBalancersResponse()
	err = c.Send(request, response)
	return
}

func NewDescribeLoadBalancersTaskResultRequest() (request *DescribeLoadBalancersTaskResultRequest) {
	request = &DescribeLoadBalancersTaskResultRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("lb", APIVersion, "DescribeLoadBalancersTaskResult")
	return
}

func NewDescribeLoadBalancersTaskResultResponse() (response *DescribeLoadBalancersTaskResultResponse) {
	response = &DescribeLoadBalancersTaskResultResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeLoadBalancersTaskResult(request *DescribeLoadBalancersTaskResultRequest) (response *DescribeLoadBalancersTaskResultResponse, err error) {
	if request == nil {
		request = NewDescribeLoadBalancersTaskResultRequest()
	}
	response = NewDescribeLoadBalancersTaskResultResponse()
	err = c.Send(request, response)
	return
}

func NewDeregisterInstancesFromForwardLBRequest() (request *DeregisterInstancesFromForwardLBRequest) {
	request = &DeregisterInstancesFromForwardLBRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("lb", APIVersion, "DeregisterInstancesFromForwardLB")
	return
}

func NewDeregisterInstancesFromForwardLBResponse() (response *DeregisterInstancesFromForwardLBResponse) {
	response = &DeregisterInstancesFromForwardLBResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DeregisterInstancesFromForwardLB(request *DeregisterInstancesFromForwardLBRequest) (response *DeregisterInstancesFromForwardLBResponse, err error) {
	if request == nil {
		request = NewDeregisterInstancesFromForwardLBRequest()
	}
	response = NewDeregisterInstancesFromForwardLBResponse()
	err = c.Send(request, response)
	return
}

func NewModifyForwardLBNameRequest() (request *ModifyForwardLBNameRequest) {
	request = &ModifyForwardLBNameRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("lb", APIVersion, "ModifyForwardLBName")
	return
}

func NewModifyForwardLBNameResponse() (response *ModifyForwardLBNameResponse) {
	response = &ModifyForwardLBNameResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) ModifyForwardLBName(request *ModifyForwardLBNameRequest) (response *ModifyForwardLBNameResponse, err error) {
	if request == nil {
		request = NewModifyForwardLBNameRequest()
	}
	response = NewModifyForwardLBNameResponse()
	err = c.Send(request, response)
	return
}

func NewModifyLoadBalancerAttributesRequest() (request *ModifyLoadBalancerAttributesRequest) {
	request = &ModifyLoadBalancerAttributesRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("lb", APIVersion, "ModifyLoadBalancerAttributes")
	return
}

func NewModifyLoadBalancerAttributesResponse() (response *ModifyLoadBalancerAttributesResponse) {
	response = &ModifyLoadBalancerAttributesResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) ModifyLoadBalancerAttributes(request *ModifyLoadBalancerAttributesRequest) (response *ModifyLoadBalancerAttributesResponse, err error) {
	if request == nil {
		request = NewModifyLoadBalancerAttributesRequest()
	}
	response = NewModifyLoadBalancerAttributesResponse()
	err = c.Send(request, response)
	return
}

func NewRegisterInstancesWithForwardLBSeventhListenerRequest() (request *RegisterInstancesWithForwardLBSeventhListenerRequest) {
	request = &RegisterInstancesWithForwardLBSeventhListenerRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("lb", APIVersion, "RegisterInstancesWithForwardLBSeventhListener")
	return
}

func NewRegisterInstancesWithForwardLBSeventhListenerResponse() (response *RegisterInstancesWithForwardLBSeventhListenerResponse) {
	response = &RegisterInstancesWithForwardLBSeventhListenerResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) RegisterInstancesWithForwardLBSeventhListener(request *RegisterInstancesWithForwardLBSeventhListenerRequest) (response *RegisterInstancesWithForwardLBSeventhListenerResponse, err error) {
	if request == nil {
		request = NewRegisterInstancesWithForwardLBSeventhListenerRequest()
	}
	response = NewRegisterInstancesWithForwardLBSeventhListenerResponse()
	err = c.Send(request, response)
	return
}
