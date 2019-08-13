package tencentcloud

import (
	"context"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/ratelimit"
	"log"
)

// info  for one direct connect gateway.
type DcgInstanceInfo struct {
	dcgId             string
	name              string
	dcgIp             string
	networkType       string
	networkInstanceId string
	gatewayType       string
	cnnRouteType      string
	createTime        string
	enableBGP         bool
}

//info for direct connect gateway[ ccn type] route.
type DcgRouteInfo struct {
	dcgId     string
	routeId   string
	cidrBlock string
	asPaths   []string
}

func (me *VpcService) DescribeDirectConnectGateway(ctx context.Context, dcgId string) (info DcgInstanceInfo, has int, errRet error) {

	infos, err := me.DescribeDirectConnectGateways(ctx, dcgId, "")
	if err != nil {
		errRet = err
		return
	}
	has = len(infos)

	if has > 1 {
		errRet = fmt.Errorf("One dcgId get %d instances by api %s",
			has,
			"DescribeDirectConnectGateways")
		return
	}
	if has == 1 {
		info = infos[0]
	}
	return
}

func (me *VpcService) DescribeDirectConnectGateways(ctx context.Context, dcgId, name string) (
	infos []DcgInstanceInfo, errRet error) {

	logId := getLogId(ctx)
	request := vpc.NewDescribeDirectConnectGatewaysRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]\n",
				logId,
				request.GetAction(),
				errRet.Error())
		}
	}()

	infos = make([]DcgInstanceInfo, 0, 100)

	var offset uint64 = 0
	var limit uint64 = 100
	var total = -1
	var has = map[string]bool{}

	var filters []*vpc.Filter
	if dcgId != "" {
		filters = me.fillFilter(filters, "direct-connect-gateway-id", dcgId)
	}
	if name != "" {
		filters = me.fillFilter(filters, "direct-connect-gateway-name", name)
	}
	if len(filters) > 0 {
		request.Filters = filters
	}

getMoreData:

	if total >= 0 && int(offset) >= total {
		return
	}
	request.Limit = &limit
	request.Offset = &offset
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().DescribeDirectConnectGateways(request)

	if err != nil {
		errRet = err
		responseStr := ""
		if response != nil {
			responseStr = response.ToJsonString()
		}
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s],response body [%s], reason[%s]\n",
			logId,
			request.GetAction(),
			request.ToJsonString(),
			responseStr,
			errRet.Error())
		return
	}

	if total < 0 {
		total = int(*response.Response.TotalCount)
	}

	if len(response.Response.DirectConnectGatewaySet) > 0 {
		offset += limit
	} else {
		//get empty ,we're done
		return
	}
	for _, item := range response.Response.DirectConnectGatewaySet {
		var basicInfo DcgInstanceInfo

		basicInfo.dcgId = *item.DirectConnectGatewayId
		basicInfo.name = *item.DirectConnectGatewayName
		basicInfo.dcgIp = *item.DirectConnectGatewayIp
		basicInfo.createTime = *item.CreateTime
		basicInfo.networkType = *item.NetworkType
		basicInfo.networkInstanceId = *item.NetworkInstanceId
		basicInfo.enableBGP = *item.EnableBGP

		if basicInfo.networkType != DCG_NETWORK_TYPE_VPC &&
			basicInfo.networkType != DCG_NETWORK_TYPE_CCN {

			errRet = fmt.Errorf("%s return an unknown NetworkType: %s",
				request.GetAction(),
				*item.NetworkType)
			return

		}

		if basicInfo.networkType == DCG_NETWORK_TYPE_CCN {

			basicInfo.cnnRouteType = *item.CcnRouteType
			if basicInfo.cnnRouteType != DCG_CCN_ROUTE_TYPE_BGP &&
				basicInfo.cnnRouteType != DCG_CCN_ROUTE_TYPE_STATIC {
				errRet = fmt.Errorf("%s return an unknown CcnRouteType: %s",
					request.GetAction(),
					*item.CcnRouteType)
				return
			}
		}
		basicInfo.gatewayType = *item.GatewayType

		if basicInfo.gatewayType != DCG_GATEWAY_TYPE_NORMAL &&
			basicInfo.gatewayType != DCG_GATEWAY_TYPE_NAT {

			errRet = fmt.Errorf("%s return an unknown GatewayType: %s",
				request.GetAction(),
				*item.GatewayType)
			return

		}

		if has[basicInfo.dcgId] {

			errRet = fmt.Errorf("get repeated dcgId[%s] when doing %s",
				basicInfo.dcgId,
				request.GetAction())
			return

		}
		has[basicInfo.dcgId] = true
		infos = append(infos, basicInfo)
	}
	goto getMoreData
}

func (me *VpcService) GetCcnRouteId(ctx context.Context, dcgId, cidr string, asPaths []string) (routeId string, has int, errRet error) {

	infos, err := me.DescribeDirectConnectGatewayCcnRoutes(ctx, dcgId)
	if err != nil {
		errRet = err
		return
	}

	for _, info := range infos {
		if info.cidrBlock == cidr {
			has++
			routeId = info.routeId
		}
	}

	if has > 1 {
		errRet = fmt.Errorf("One cidr '%s'  get %d instances by api %s",
			cidr,
			has,
			"DescribeDirectConnectGatewayCcnRoutes")
		return
	}

	return
}

func (me *VpcService) DescribeDirectConnectGatewayCcnRoute(ctx context.Context, dcgId, routeId string) (infoRet DcgRouteInfo, has int, errRet error) {

	infos, err := me.DescribeDirectConnectGatewayCcnRoutes(ctx, dcgId)
	if err != nil {
		errRet = err
		return
	}

	for _, info := range infos {
		if info.routeId == routeId {
			has++
			infoRet = info
		}
	}

	if has > 1 {
		errRet = fmt.Errorf("One routeId get %d instances by api %s",
			has,
			"DescribeDirectConnectGatewayCcnRoutes")
		return
	}

	return

}

func (me *VpcService) DescribeDirectConnectGatewayCcnRoutes(ctx context.Context, dcgId string) (infos []DcgRouteInfo, errRet error) {
	logId := getLogId(ctx)
	request := vpc.NewDescribeDirectConnectGatewayCcnRoutesRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]\n",
				logId,
				request.GetAction(),
				errRet.Error())
		}
	}()

	request.DirectConnectGatewayId = &dcgId

	infos = make([]DcgRouteInfo, 0, 100)
	var offset uint64 = 0
	var limit uint64 = 100
	var total = -1
	var has = map[string]bool{}

getMoreData:
	if total >= 0 && int(offset) >= total {
		return
	}
	request.Limit = &limit
	request.Offset = &offset
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().DescribeDirectConnectGatewayCcnRoutes(request)

	if err != nil {
		if sdkErr, ok := err.(*errors.TencentCloudSDKError); ok {
			if sdkErr.Code == "ResourceNotFound" {
				return
			}
		}
	}

	if err != nil {
		errRet = err
		responseStr := ""
		if response != nil {
			responseStr = response.ToJsonString()
		}
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s],response body [%s], reason[%s]\n",
			logId,
			request.GetAction(),
			request.ToJsonString(),
			responseStr,
			errRet.Error())
		return
	}

	if total < 0 {
		total = int(*response.Response.TotalCount)
	}

	if len(response.Response.RouteSet) > 0 {
		offset += limit
	} else {
		//get empty ,we're done
		return
	}
	for _, item := range response.Response.RouteSet {
		var basicInfo DcgRouteInfo
		basicInfo.routeId = *item.RouteId
		basicInfo.dcgId = dcgId
		basicInfo.cidrBlock = *item.DestinationCidrBlock
		basicInfo.asPaths = make([]string, 0, len(item.ASPath))
		for index := range item.ASPath {
			basicInfo.asPaths = append(basicInfo.asPaths, *item.ASPath[index])
		}

		if has[basicInfo.routeId] {
			errRet = fmt.Errorf("get repeated routeId[%s] when doing %s",
				basicInfo.routeId,
				request.GetAction())
			return
		}
		has[basicInfo.routeId] = true
		infos = append(infos, basicInfo)
	}
	goto getMoreData
}

func (me *VpcService) CreateDirectConnectGateway(ctx context.Context, name, networkType, networkInstanceId, gatewayType string) (
	dcgId string, errRet error) {

	logId := getLogId(ctx)

	request := vpc.NewCreateDirectConnectGatewayRequest()

	request.DirectConnectGatewayName = &name
	request.NetworkType = &networkType
	request.NetworkInstanceId = &networkInstanceId
	request.GatewayType = &gatewayType
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().CreateDirectConnectGateway(request)

	defer func() {
		if errRet != nil {
			responseStr := ""
			if response != nil {
				responseStr = response.ToJsonString()
			}
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s],response body [%s], reason[%s]\n",
				logId,
				request.GetAction(),
				request.ToJsonString(),
				responseStr,
				errRet.Error())
		}
	}()

	if err != nil {
		errRet = err
		return
	}

	if response.Response.DirectConnectGateway == nil || response.Response.DirectConnectGateway.DirectConnectGatewayId == nil {
		errRet = fmt.Errorf("%s return empty resource id", request.GetAction())
		return
	}
	dcgId = *response.Response.DirectConnectGateway.DirectConnectGatewayId
	return
}

func (me *VpcService) ModifyDirectConnectGatewayAttribute(ctx context.Context, dcgId, name string) (errRet error) {

	if name == "" {
		return
	}

	logId := getLogId(ctx)

	request := vpc.NewModifyDirectConnectGatewayAttributeRequest()
	request.DirectConnectGatewayId = &dcgId
	request.DirectConnectGatewayName = &name
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().ModifyDirectConnectGatewayAttribute(request)

	defer func() {
		if errRet != nil {
			responseStr := ""
			if response != nil {
				responseStr = response.ToJsonString()
			}
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s],response body [%s], reason[%s]\n",
				logId,
				request.GetAction(),
				request.ToJsonString(),
				responseStr,
				errRet.Error())
		}
	}()

	errRet = err
	return
}

func (me *VpcService) DeleteDirectConnectGateway(ctx context.Context, dcgId string) (errRet error) {

	logId := getLogId(ctx)
	request := vpc.NewDeleteDirectConnectGatewayRequest()
	request.DirectConnectGatewayId = &dcgId
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().DeleteDirectConnectGateway(request)

	defer func() {
		if errRet != nil {
			responseStr := ""
			if response != nil {
				responseStr = response.ToJsonString()
			}
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s],response body [%s], reason[%s]\n",
				logId,
				request.GetAction(),
				request.ToJsonString(),
				responseStr,
				errRet.Error())
		}
	}()

	errRet = err
	return
}

func (me *VpcService) CreateDirectConnectGatewayCcnRoute(ctx context.Context, dcgId, cidr string, asPaths []string) (routeId string, errRet error) {

	logId := getLogId(ctx)

	request := vpc.NewCreateDirectConnectGatewayCcnRoutesRequest()
	request.DirectConnectGatewayId = &dcgId

	var ccnRoute vpc.DirectConnectGatewayCcnRoute
	ccnRoute.DestinationCidrBlock = &cidr
	ccnRoute.ASPath = make([]*string, 0, len(asPaths))

	for index := range asPaths {
		ccnRoute.ASPath = append(ccnRoute.ASPath, &asPaths[index])
	}
	request.Routes = []*vpc.DirectConnectGatewayCcnRoute{&ccnRoute}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().CreateDirectConnectGatewayCcnRoutes(request)

	defer func() {
		if errRet != nil {
			responseStr := ""
			if response != nil {
				responseStr = response.ToJsonString()
			}
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s],response body [%s], reason[%s]\n",
				logId,
				request.GetAction(),
				request.ToJsonString(),
				responseStr,
				errRet.Error())
		}
	}()

	if err != nil {
		errRet = err
		return
	}

	routeIdTemp, has, err := me.GetCcnRouteId(ctx, dcgId, cidr, asPaths)

	if err != nil {
		errRet = err
		return
	}

	if has == 1 {
		routeId = routeIdTemp
		return
	} else {
		errRet = fmt.Errorf("after api `CreateDirectConnectGatewayCcnRoutes`, api `DescribeDirectConnectGatewayCcnRoutes` return null route info")
	}
	return
}

func (me *VpcService) DeleteDirectConnectGatewayCcnRoute(ctx context.Context, dcgId, routeId string) (errRet error) {

	logId := getLogId(ctx)

	request := vpc.NewDeleteDirectConnectGatewayCcnRoutesRequest()
	request.DirectConnectGatewayId = &dcgId
	request.RouteIds = []*string{&routeId}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().DeleteDirectConnectGatewayCcnRoutes(request)

	defer func() {
		if errRet != nil {
			responseStr := ""
			if response != nil {
				responseStr = response.ToJsonString()
			}
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s],response body [%s], reason[%s]\n",
				logId,
				request.GetAction(),
				request.ToJsonString(),
				responseStr,
				errRet.Error())
		}
	}()

	errRet = err
	return
}
