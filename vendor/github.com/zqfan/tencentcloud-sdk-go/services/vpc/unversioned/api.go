package vpc

import (
	"github.com/zqfan/tencentcloud-sdk-go/common"
)

const APIVersion = ""

func NewAddDnaptRuleRequest() (request *AddDnaptRuleRequest) {
	request = &AddDnaptRuleRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("vpc", APIVersion, "AddDnaptRule")
	return
}

func NewAddDnaptRuleResponse() (response *AddDnaptRuleResponse) {
	response = &AddDnaptRuleResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) AddDnaptRule(request *AddDnaptRuleRequest) (response *AddDnaptRuleResponse, err error) {
	if request == nil {
		request = NewAddDnaptRuleRequest()
	}
	response = NewAddDnaptRuleResponse()
	err = c.Send(request, response)
	return
}

func NewCreateNatGatewayRequest() (request *CreateNatGatewayRequest) {
	request = &CreateNatGatewayRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("vpc", APIVersion, "CreateNatGateway")
	return
}

func NewCreateNatGatewayResponse() (response *CreateNatGatewayResponse) {
	response = &CreateNatGatewayResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) CreateNatGateway(request *CreateNatGatewayRequest) (response *CreateNatGatewayResponse, err error) {
	if request == nil {
		request = NewCreateNatGatewayRequest()
	}
	response = NewCreateNatGatewayResponse()
	err = c.Send(request, response)
	return
}

func NewDeleteDnaptRuleRequest() (request *DeleteDnaptRuleRequest) {
	request = &DeleteDnaptRuleRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("vpc", APIVersion, "DeleteDnaptRule")
	return
}

func NewDeleteDnaptRuleResponse() (response *DeleteDnaptRuleResponse) {
	response = &DeleteDnaptRuleResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DeleteDnaptRule(request *DeleteDnaptRuleRequest) (response *DeleteDnaptRuleResponse, err error) {
	if request == nil {
		request = NewDeleteDnaptRuleRequest()
	}
	response = NewDeleteDnaptRuleResponse()
	err = c.Send(request, response)
	return
}

func NewDeleteNatGatewayRequest() (request *DeleteNatGatewayRequest) {
	request = &DeleteNatGatewayRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("vpc", APIVersion, "DeleteNatGateway")
	return
}

func NewDeleteNatGatewayResponse() (response *DeleteNatGatewayResponse) {
	response = &DeleteNatGatewayResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DeleteNatGateway(request *DeleteNatGatewayRequest) (response *DeleteNatGatewayResponse, err error) {
	if request == nil {
		request = NewDeleteNatGatewayRequest()
	}
	response = NewDeleteNatGatewayResponse()
	err = c.Send(request, response)
	return
}

func NewDescribeNatGatewayRequest() (request *DescribeNatGatewayRequest) {
	request = &DescribeNatGatewayRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("vpc", APIVersion, "DescribeNatGateway")
	return
}

func NewDescribeNatGatewayResponse() (response *DescribeNatGatewayResponse) {
	response = &DescribeNatGatewayResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeNatGateway(request *DescribeNatGatewayRequest) (response *DescribeNatGatewayResponse, err error) {
	if request == nil {
		request = NewDescribeNatGatewayRequest()
	}
	response = NewDescribeNatGatewayResponse()
	err = c.Send(request, response)
	return
}

func NewDescribeNetworkInterfacesRequest() (request *DescribeNetworkInterfacesRequest) {
	request = &DescribeNetworkInterfacesRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("vpc", APIVersion, "DescribeNetworkInterfaces")
	return
}

func NewDescribeNetworkInterfacesResponse() (response *DescribeNetworkInterfacesResponse) {
	response = &DescribeNetworkInterfacesResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeNetworkInterfaces(request *DescribeNetworkInterfacesRequest) (response *DescribeNetworkInterfacesResponse, err error) {
	if request == nil {
		request = NewDescribeNetworkInterfacesRequest()
	}
	response = NewDescribeNetworkInterfacesResponse()
	err = c.Send(request, response)
	return
}

func NewDescribeVpcExRequest() (request *DescribeVpcExRequest) {
	request = &DescribeVpcExRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("vpc", APIVersion, "DescribeVpcEx")
	return
}

func NewDescribeVpcExResponse() (response *DescribeVpcExResponse) {
	response = &DescribeVpcExResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeVpcEx(request *DescribeVpcExRequest) (response *DescribeVpcExResponse, err error) {
	if request == nil {
		request = NewDescribeVpcExRequest()
	}
	response = NewDescribeVpcExResponse()
	err = c.Send(request, response)
	return
}

func NewDescribeVpcTaskResultRequest() (request *DescribeVpcTaskResultRequest) {
	request = &DescribeVpcTaskResultRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("vpc", APIVersion, "DescribeVpcTaskResult")
	return
}

func NewDescribeVpcTaskResultResponse() (response *DescribeVpcTaskResultResponse) {
	response = &DescribeVpcTaskResultResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeVpcTaskResult(request *DescribeVpcTaskResultRequest) (response *DescribeVpcTaskResultResponse, err error) {
	if request == nil {
		request = NewDescribeVpcTaskResultRequest()
	}
	response = NewDescribeVpcTaskResultResponse()
	err = c.Send(request, response)
	return
}

func NewEipBindNatGatewayRequest() (request *EipBindNatGatewayRequest) {
	request = &EipBindNatGatewayRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("vpc", APIVersion, "EipBindNatGateway")
	return
}

func NewEipBindNatGatewayResponse() (response *EipBindNatGatewayResponse) {
	response = &EipBindNatGatewayResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) EipBindNatGateway(request *EipBindNatGatewayRequest) (response *EipBindNatGatewayResponse, err error) {
	if request == nil {
		request = NewEipBindNatGatewayRequest()
	}
	response = NewEipBindNatGatewayResponse()
	err = c.Send(request, response)
	return
}

func NewEipUnBindNatGatewayRequest() (request *EipUnBindNatGatewayRequest) {
	request = &EipUnBindNatGatewayRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("vpc", APIVersion, "EipUnBindNatGateway")
	return
}

func NewEipUnBindNatGatewayResponse() (response *EipUnBindNatGatewayResponse) {
	response = &EipUnBindNatGatewayResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) EipUnBindNatGateway(request *EipUnBindNatGatewayRequest) (response *EipUnBindNatGatewayResponse, err error) {
	if request == nil {
		request = NewEipUnBindNatGatewayRequest()
	}
	response = NewEipUnBindNatGatewayResponse()
	err = c.Send(request, response)
	return
}

func NewGetDnaptRuleRequest() (request *GetDnaptRuleRequest) {
	request = &GetDnaptRuleRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("vpc", APIVersion, "GetDnaptRule")
	return
}

func NewGetDnaptRuleResponse() (response *GetDnaptRuleResponse) {
	response = &GetDnaptRuleResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) GetDnaptRule(request *GetDnaptRuleRequest) (response *GetDnaptRuleResponse, err error) {
	if request == nil {
		request = NewGetDnaptRuleRequest()
	}
	response = NewGetDnaptRuleResponse()
	err = c.Send(request, response)
	return
}

func NewModifyDnaptRuleRequest() (request *ModifyDnaptRuleRequest) {
	request = &ModifyDnaptRuleRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("vpc", APIVersion, "ModifyDnaptRule")
	return
}

func NewModifyDnaptRuleResponse() (response *ModifyDnaptRuleResponse) {
	response = &ModifyDnaptRuleResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) ModifyDnaptRule(request *ModifyDnaptRuleRequest) (response *ModifyDnaptRuleResponse, err error) {
	if request == nil {
		request = NewModifyDnaptRuleRequest()
	}
	response = NewModifyDnaptRuleResponse()
	err = c.Send(request, response)
	return
}

func NewModifyNatGatewayRequest() (request *ModifyNatGatewayRequest) {
	request = &ModifyNatGatewayRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("vpc", APIVersion, "ModifyNatGateway")
	return
}

func NewModifyNatGatewayResponse() (response *ModifyNatGatewayResponse) {
	response = &ModifyNatGatewayResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) ModifyNatGateway(request *ModifyNatGatewayRequest) (response *ModifyNatGatewayResponse, err error) {
	if request == nil {
		request = NewModifyNatGatewayRequest()
	}
	response = NewModifyNatGatewayResponse()
	err = c.Send(request, response)
	return
}

func NewQueryNatGatewayProductionStatusRequest() (request *QueryNatGatewayProductionStatusRequest) {
	request = &QueryNatGatewayProductionStatusRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("vpc", APIVersion, "QueryNatGatewayProductionStatus")
	return
}

func NewQueryNatGatewayProductionStatusResponse() (response *QueryNatGatewayProductionStatusResponse) {
	response = &QueryNatGatewayProductionStatusResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) QueryNatGatewayProductionStatus(request *QueryNatGatewayProductionStatusRequest) (response *QueryNatGatewayProductionStatusResponse, err error) {
	if request == nil {
		request = NewQueryNatGatewayProductionStatusRequest()
	}
	response = NewQueryNatGatewayProductionStatusResponse()
	err = c.Send(request, response)
	return
}

func NewUpgradeNatGatewayRequest() (request *UpgradeNatGatewayRequest) {
	request = &UpgradeNatGatewayRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("vpc", APIVersion, "UpgradeNatGateway")
	return
}

func NewUpgradeNatGatewayResponse() (response *UpgradeNatGatewayResponse) {
	response = &UpgradeNatGatewayResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) UpgradeNatGateway(request *UpgradeNatGatewayRequest) (response *UpgradeNatGatewayResponse, err error) {
	if request == nil {
		request = NewUpgradeNatGatewayRequest()
	}
	response = NewUpgradeNatGatewayResponse()
	err = c.Send(request, response)
	return
}
