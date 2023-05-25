package tencentcloud

import (
	"context"
	"fmt"
	"log"

	dc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dc/v20180410"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type DcService struct {
	client *connectivity.TencentCloudClient
}

/////////common
func (me *DcService) fillFilter(ins []*dc.Filter, key, value string) (outs []*dc.Filter) {
	if ins == nil {
		ins = make([]*dc.Filter, 0, 2)
	}

	var filter = dc.Filter{Name: &key, Values: []*string{&value}}
	ins = append(ins, &filter)
	outs = ins
	return
}

func (me *DcService) strPt2str(pt *string) (ret string) {
	if pt == nil {
		return
	} else {
		return *pt
	}
}

/*
func (me *DcService) intPt2int(pt *int) (ret int) {
	if pt == nil {
		return
	} else {
		return *pt
	}
}
*/

func (me *DcService) int64Pt2int64(pt *int64) (ret int64) {
	if pt == nil {
		return
	} else {
		return *pt
	}
}

func (me *DcService) DescribeDirectConnects(ctx context.Context, dcId,
	name string) (infos []dc.DirectConnect, errRet error) {

	logId := getLogId(ctx)
	request := dc.NewDescribeDirectConnectsRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var offset int64 = 0
	var limit int64 = 100
	var total int64 = -1
	var has = map[string]bool{}

	var filters []*dc.Filter
	if dcId != "" {
		filters = me.fillFilter(filters, "direct-connect-id", dcId)
	}
	if name != "" {
		filters = me.fillFilter(filters, "direct-connect-name", name)
	}
	if len(filters) > 0 {
		request.Filters = filters
	}
	infos = make([]dc.DirectConnect, 0, 10)

getMoreData:
	if total >= 0 && offset >= total {
		return
	}
	request.Limit = &limit
	request.Offset = &offset
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDcClient().DescribeDirectConnects(request)
	if err != nil {
		errRet = err
		return
	}
	if total < 0 {
		total = *response.Response.TotalCount
	}

	if len(response.Response.DirectConnectSet) > 0 {
		offset += limit
	} else {
		//get empty set,we're done
		return
	}

	for _, item := range response.Response.DirectConnectSet {
		if has[*item.DirectConnectId] {
			errRet = fmt.Errorf("get repeated dc_id[%s] when doing DescribeDirectConnects", *item.DirectConnectId)
			return
		}
		has[*item.DirectConnectId] = true
		infos = append(infos, *item)
	}
	goto getMoreData
}

func (me *DcService) DescribeDirectConnectTunnel(ctx context.Context, dcxId string) (info dc.DirectConnectTunnel, has int64, errRet error) {

	infos, err := me.DescribeDirectConnectTunnels(ctx, dcxId, "")

	if err != nil {
		errRet = err
		return
	}
	has = int64(len(infos))

	if has > 0 {
		info = infos[0]

	}
	return
}

func (me *DcService) DescribeDirectConnectTunnels(ctx context.Context, dcxId,
	name string) (infos []dc.DirectConnectTunnel, errRet error) {

	logId := getLogId(ctx)
	request := dc.NewDescribeDirectConnectTunnelsRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var offset int64 = 0
	var limit int64 = 100
	var total int64 = -1
	var has = map[string]bool{}

	var filters []*dc.Filter
	if dcxId != "" {
		filters = me.fillFilter(filters, "direct-connect-tunnel-id", dcxId)
	}
	if name != "" {
		filters = me.fillFilter(filters, "direct-connect-tunnel-name", name)
	}
	if len(filters) > 0 {
		request.Filters = filters
	}
	infos = make([]dc.DirectConnectTunnel, 0, 10)
getMoreData:
	if total >= 0 && offset >= total {
		return
	}
	request.Limit = &limit
	request.Offset = &offset
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDcClient().DescribeDirectConnectTunnels(request)
	if err != nil {
		errRet = err
		return
	}
	if total < 0 {
		total = *response.Response.TotalCount
	}

	if len(response.Response.DirectConnectTunnelSet) > 0 {
		offset += limit
	} else {
		//get empty set,we're done
		return
	}
	for _, item := range response.Response.DirectConnectTunnelSet {
		if has[*item.DirectConnectTunnelId] {
			errRet = fmt.Errorf("get repeated dcx_id[%s] when doing DescribeDirectConnectTunnels", *item.DirectConnectTunnelId)
			return
		}
		has[*item.DirectConnectTunnelId] = true
		infos = append(infos, *item)
	}
	goto getMoreData
}

func (me *DcService) CreateDirectConnectTunnel(ctx context.Context, dcId, dcxName, networkType,
	networkRegion, vpcId, routeType, bgpAuthKey,
	tencentAddress, customerAddress, dcgId string,
	bgpAsn, vlan, bandwidth int64,
	routeFilterPrefixes []string) (dcxId string, errRet error) {

	logId := getLogId(ctx)
	request := dc.NewCreateDirectConnectTunnelRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.DirectConnectId = &dcId
	request.DirectConnectTunnelName = &dcxName
	request.NetworkType = &networkType
	request.NetworkRegion = &networkRegion
	if vpcId != "" {
		request.VpcId = &vpcId
	}
	if bandwidth >= 0 {
		request.Bandwidth = &bandwidth
	}
	request.RouteType = &routeType
	request.DirectConnectGatewayId = &dcgId
	if bgpAsn >= 0 {
		var peer dc.BgpPeer
		peer.Asn = &bgpAsn
		peer.AuthKey = &bgpAuthKey
		request.BgpPeer = &peer
	}

	request.Vlan = &vlan

	if len(routeFilterPrefixes) > 0 {
		request.RouteFilterPrefixes = make([]*dc.RouteFilterPrefix, 0, len(routeFilterPrefixes))
		for index := range routeFilterPrefixes {
			var dcPrefix dc.RouteFilterPrefix
			dcPrefix.Cidr = &routeFilterPrefixes[index]
			request.RouteFilterPrefixes = append(request.RouteFilterPrefixes, &dcPrefix)
		}
	}

	if tencentAddress != "" {
		request.TencentAddress = &tencentAddress
	}

	if customerAddress != "" {
		request.CustomerAddress = &customerAddress
	}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDcClient().CreateDirectConnectTunnel(request)
	if err != nil {
		errRet = err
		return
	}

	if len(response.Response.DirectConnectTunnelIdSet) != 1 {
		errRet = fmt.Errorf("CreateDirectConnectTunnel  return %d DirectConnectTunnelIdSet",
			len(response.Response.DirectConnectTunnelIdSet))
		return
	}
	dcxId = *response.Response.DirectConnectTunnelIdSet[0]
	return
}

func (me *DcService) DeleteDirectConnectTunnel(ctx context.Context, dcxId string) (errRet error) {

	logId := getLogId(ctx)
	request := dc.NewDeleteDirectConnectTunnelRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.DirectConnectTunnelId = &dcxId
	ratelimit.Check(request.GetAction())
	_, err := me.client.UseDcClient().DeleteDirectConnectTunnel(request)
	if err != nil {
		errRet = err
	}
	return
}

func (me *DcService) ModifyDirectConnectTunnelAttribute(ctx context.Context, dcxId string,
	name, bgpAuthKey, tencentAddress, customerAddress string,
	bandwidth, bgpAsn int64,
	routeFilterPrefixes []string) (errRet error) {

	logId := getLogId(ctx)
	request := dc.NewModifyDirectConnectTunnelAttributeRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.DirectConnectTunnelId = &dcxId
	if name != "" {
		request.DirectConnectTunnelName = &name
	}
	if tencentAddress != "" {
		request.TencentAddress = &tencentAddress
	}
	if customerAddress != "" {
		request.CustomerAddress = &customerAddress
	}

	if bgpAsn >= 0 {
		var peer dc.BgpPeer
		peer.Asn = &bgpAsn
		peer.AuthKey = &bgpAuthKey
		request.BgpPeer = &peer
	}

	if bandwidth > 0 {
		request.Bandwidth = &bandwidth
	}

	if len(routeFilterPrefixes) > 0 {
		request.RouteFilterPrefixes = make([]*dc.RouteFilterPrefix, 0, len(routeFilterPrefixes))
		for index := range routeFilterPrefixes {
			var dcPrefix dc.RouteFilterPrefix
			dcPrefix.Cidr = &routeFilterPrefixes[index]
			request.RouteFilterPrefixes = append(request.RouteFilterPrefixes, &dcPrefix)
		}
	}
	ratelimit.Check(request.GetAction())
	_, err := me.client.UseDcClient().ModifyDirectConnectTunnelAttribute(request)
	if err != nil {
		errRet = err
	}
	return
}

func (me *DcService) DescribeDcShareDcxConfigById(ctx context.Context, directConnectTunnelId string) (ShareDcxConfig *dc.DirectConnectTunnel, errRet error) {
	logId := getLogId(ctx)

	request := dc.NewDescribeDirectConnectTunnelsRequest()
	request.DirectConnectTunnelIds = []*string{&directConnectTunnelId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDcClient().DescribeDirectConnectTunnels(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.DirectConnectTunnelSet) < 1 {
		return
	}

	ShareDcxConfig = response.Response.DirectConnectTunnelSet[0]
	return
}

func (me *DcService) DescribeDcInternetAddressById(ctx context.Context, instanceId string) (internetAddress *dc.InternetAddressDetail, errRet error) {
	logId := getLogId(ctx)

	request := dc.NewDescribeInternetAddressRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDcClient().DescribeInternetAddress(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	for _, address := range response.Response.Subnets {
		if *address.InstanceId == instanceId {
			internetAddress = address
			break
		}
	}
	return
}

func (me *DcService) DeleteDcInternetAddressById(ctx context.Context, instanceId string) (errRet error) {
	logId := getLogId(ctx)

	request := dc.NewReleaseInternetAddressRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDcClient().ReleaseInternetAddress(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *DcService) DescribeDcxExtraConfigById(ctx context.Context, directConnectTunnelId string) (dcxExtraConfig *dc.DirectConnectTunnelExtra, errRet error) {
	logId := getLogId(ctx)

	request := dc.NewDescribeDirectConnectTunnelExtraRequest()
	request.DirectConnectTunnelId = &directConnectTunnelId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDcClient().DescribeDirectConnectTunnelExtra(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	dcxExtraConfig = response.Response.DirectConnectTunnelExtra
	return
}

func (me *DcService) DescribeDcInternetAddressQuota(ctx context.Context) (InternetAddressQuota *dc.DescribeInternetAddressQuotaResponse, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dc.NewDescribeInternetAddressQuotaRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDcClient().DescribeInternetAddressQuota(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	InternetAddressQuota = response

	return
}

func (me *DcService) DescribeDcInternetAddressStatistics(ctx context.Context) (statistics []*dc.InternetAddressStatistics, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dc.NewDescribeInternetAddressStatisticsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDcClient().DescribeInternetAddressStatistics(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	statistics = response.Response.InternetAddressStatistics
	return
}

func (me *DcService) DescribeDcPublicDirectConnectTunnelRoutesByFilter(ctx context.Context, param map[string]interface{}) (PublicDirectConnectTunnelRoutes []*dc.DirectConnectTunnelRoute, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dc.NewDescribePublicDirectConnectTunnelRoutesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "DirectConnectTunnelId" {
			request.DirectConnectTunnelId = v.(*string)
		}
		if k == "Filters" {
			request.Filters = v.([]*dc.Filter)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseDcClient().DescribePublicDirectConnectTunnelRoutes(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Routes) < 1 {
			break
		}
		PublicDirectConnectTunnelRoutes = append(PublicDirectConnectTunnelRoutes, response.Response.Routes...)
		if len(response.Response.Routes) < int(limit) {
			break
		}

		offset += limit
	}

	return
}
