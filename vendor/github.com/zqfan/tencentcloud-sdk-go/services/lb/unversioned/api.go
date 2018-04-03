package lb

import (
	"github.com/zqfan/tencentcloud-sdk-go/common"
)

const APIVersion = ""

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
