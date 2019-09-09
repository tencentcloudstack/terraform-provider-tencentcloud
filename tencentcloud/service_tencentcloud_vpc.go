package tencentcloud

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

// VPC basic information
type VpcBasicInfo struct {
	vpcId       string
	name        string
	cidr        string
	isMulticast bool
	isDefault   bool
	dnsServers  []string
	createTime  string
}

// subnet basic information
type VpcSubnetBasicInfo struct {
	vpcId            string
	subnetId         string
	routeTableId     string
	name             string
	cidr             string
	isMulticast      bool
	isDefault        bool
	zone             string
	availableIpCount int64
	createTime       string
}

// route entry basic information
type VpcRouteEntryBasicInfo struct {
	routeEntryId    int64
	destinationCidr string
	nextType        string
	nextBub         string
	description     string
	entryType       string
}

// route table basic information
type VpcRouteTableBasicInfo struct {
	routeTableId string
	name         string
	vpcId        string
	isDefault    bool
	subnetIds    []string
	entryInfos   []VpcRouteEntryBasicInfo
	createTime   string
}

type VpcEniIP struct {
	ip      net.IP
	primary bool
	desc    *string
}

type VpcService struct {
	client *connectivity.TencentCloudClient
}

// ///////common
func (me *VpcService) fillFilter(ins []*vpc.Filter, key, value string) (outs []*vpc.Filter) {
	if ins == nil {
		ins = make([]*vpc.Filter, 0, 2)
	}

	var filter = vpc.Filter{Name: &key, Values: []*string{&value}}
	ins = append(ins, &filter)
	outs = ins
	return
}

// ////////api
func (me *VpcService) CreateVpc(ctx context.Context, name, cidr string,
	isMulticast bool, dnsServers []string) (vpcId string, isDefault bool, errRet error) {

	logId := getLogId(ctx)
	request := vpc.NewCreateVpcRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.VpcName = &name
	request.CidrBlock = &cidr

	var enableMulticast = map[bool]string{true: "true", false: "false"}[isMulticast]
	request.EnableMulticast = &enableMulticast

	if len(dnsServers) > 0 {
		request.DnsServers = make([]*string, 0, len(dnsServers))
		for index := range dnsServers {
			request.DnsServers = append(request.DnsServers, &dnsServers[index])
		}
	}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().CreateVpc(request)
	if err == nil {
		log.Printf("[DEBUG]%s api[%s] , request body [%s], response body[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		vpcId, isDefault = *response.Response.Vpc.VpcId, *response.Response.Vpc.IsDefault
		return
	}

	errRet = err
	return
}

func (me *VpcService) DescribeVpc(ctx context.Context, vpcId string) (info VpcBasicInfo, has int, errRet error) {
	infos, err := me.DescribeVpcs(ctx, vpcId, "")
	if err != nil {
		errRet = err
		return
	}
	has = len(infos)
	if has > 0 {
		info = infos[0]
	}
	return
}

func (me *VpcService) DescribeVpcs(ctx context.Context, vpcId, name string) (infos []VpcBasicInfo, errRet error) {
	logId := getLogId(ctx)
	request := vpc.NewDescribeVpcsRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	infos = make([]VpcBasicInfo, 0, 100)

	var offset = 0
	var limit = 100
	var total = -1
	var hasVpc = map[string]bool{}

	var filters []*vpc.Filter
	if vpcId != "" {
		filters = me.fillFilter(filters, "vpc-id", vpcId)
	}
	if name != "" {
		filters = me.fillFilter(filters, "vpc-name", name)
	}
	if len(filters) > 0 {
		request.Filters = filters
	}

getMoreData:

	if total >= 0 {
		if offset >= total {
			return
		}
	}
	var strLimit = fmt.Sprintf("%d", limit)
	request.Limit = &strLimit

	var strOffset = fmt.Sprintf("%d", offset)
	request.Offset = &strOffset
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().DescribeVpcs(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] , request body [%s], response body[%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if total < 0 {
		total = int(*response.Response.TotalCount)
	}

	if len(response.Response.VpcSet) > 0 {
		offset += limit
	} else {
		// get empty Vpcinfo,we're done
		return
	}
	for _, item := range response.Response.VpcSet {
		var basicInfo VpcBasicInfo
		basicInfo.cidr = *item.CidrBlock
		basicInfo.createTime = *item.CreatedTime
		basicInfo.dnsServers = make([]string, 0, len(item.DnsServerSet))

		for _, v := range item.DnsServerSet {
			basicInfo.dnsServers = append(basicInfo.dnsServers, *v)
		}
		basicInfo.isDefault = *item.IsDefault
		basicInfo.isMulticast = *item.EnableMulticast
		basicInfo.name = *item.VpcName
		basicInfo.vpcId = *item.VpcId

		if hasVpc[basicInfo.vpcId] {
			errRet = fmt.Errorf("get repeated vpc_id[%s] when doing DescribeVpcs", basicInfo.vpcId)
			return
		}
		hasVpc[basicInfo.vpcId] = true
		infos = append(infos, basicInfo)
	}
	goto getMoreData

}
func (me *VpcService) DescribeSubnet(ctx context.Context, subnetId string) (info VpcSubnetBasicInfo, has int, errRet error) {
	infos, err := me.DescribeSubnets(ctx, subnetId, "", "", "")
	if err != nil {
		errRet = err
		return
	}
	has = len(infos)
	if has > 0 {
		info = infos[0]
	}
	return
}

func (me *VpcService) DescribeSubnets(ctx context.Context, subnet_id, vpc_id, subnet_name, zone string) (infos []VpcSubnetBasicInfo, errRet error) {

	logId := getLogId(ctx)
	request := vpc.NewDescribeSubnetsRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	var offset = 0
	var limit = 100
	var total = -1
	var hasSubnet = map[string]bool{}

	var filters []*vpc.Filter
	if subnet_id != "" {
		filters = me.fillFilter(filters, "subnet-id", subnet_id)
	}
	if vpc_id != "" {
		filters = me.fillFilter(filters, "vpc-id", vpc_id)
	}
	if subnet_name != "" {
		filters = me.fillFilter(filters, "subnet-name", subnet_name)
	}
	if zone != "" {
		filters = me.fillFilter(filters, "zone", zone)
	}

	if len(filters) > 0 {
		request.Filters = filters
	}

getMoreData:
	if total >= 0 {
		if offset >= total {
			return
		}
	}
	var strLimit = fmt.Sprintf("%d", limit)
	request.Limit = &strLimit

	var strOffset = fmt.Sprintf("%d", offset)
	request.Offset = &strOffset
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().DescribeSubnets(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] , request body [%s], response body[%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if total < 0 {
		total = int(*response.Response.TotalCount)
	}

	if len(response.Response.SubnetSet) > 0 {
		offset += limit
	} else {
		// get empty subnet ,we're done
		return
	}
	for _, item := range response.Response.SubnetSet {
		var basicInfo VpcSubnetBasicInfo

		basicInfo.cidr = *item.CidrBlock
		basicInfo.createTime = *item.CreatedTime
		basicInfo.vpcId = *item.VpcId
		basicInfo.subnetId = *item.SubnetId
		basicInfo.routeTableId = *item.RouteTableId

		basicInfo.name = *item.SubnetName
		basicInfo.isDefault = *item.IsDefault
		basicInfo.isMulticast = *item.EnableBroadcast

		basicInfo.zone = *item.Zone
		basicInfo.availableIpCount = int64(*item.AvailableIpAddressCount)

		if hasSubnet[basicInfo.subnetId] {
			errRet = fmt.Errorf("get repeated subnetId[%s] when doing DescribeSubnets", basicInfo.subnetId)
			return
		}
		hasSubnet[basicInfo.subnetId] = true
		infos = append(infos, basicInfo)
	}
	goto getMoreData
}

func (me *VpcService) ModifyVpcAttribute(ctx context.Context, vpcId, name string, isMulticast bool, dnsServers []string) (errRet error) {
	logId := getLogId(ctx)
	request := vpc.NewModifyVpcAttributeRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.VpcId = &vpcId
	request.VpcName = &name

	if len(dnsServers) > 0 {
		request.DnsServers = make([]*string, 0, len(dnsServers))
		for index := range dnsServers {
			request.DnsServers = append(request.DnsServers, &dnsServers[index])
		}
	}
	var enableMulticast = map[bool]string{true: "true", false: "false"}[isMulticast]
	request.EnableMulticast = &enableMulticast
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().ModifyVpcAttribute(request)
	if err == nil {
		log.Printf("[DEBUG]%s api[%s] , request body [%s], response body[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	}
	errRet = err
	return
}

func (me *VpcService) DeleteVpc(ctx context.Context, vpcId string) (errRet error) {
	logId := getLogId(ctx)
	request := vpc.NewDeleteVpcRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	if vpcId == "" {
		errRet = fmt.Errorf("DeleteVpc can not delete empty vpc_id.")
		return
	}

	request.VpcId = &vpcId
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().DeleteVpc(request)
	if err == nil {
		log.Printf("[DEBUG]%s api[%s] , request body [%s], response body[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	}
	errRet = err
	return

}

func (me *VpcService) CreateSubnet(ctx context.Context, vpcId, name, cidr, zone string) (subnetId string, errRet error) {
	logId := getLogId(ctx)
	request := vpc.NewCreateSubnetRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	if vpcId == "" {
		errRet = fmt.Errorf("CreateSubnet can not invoke by empty vpc_id.")
		return
	}
	request.VpcId = &vpcId
	request.SubnetName = &name
	request.CidrBlock = &cidr
	request.Zone = &zone
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().CreateSubnet(request)
	errRet = err
	if err == nil {
		log.Printf("[DEBUG]%s api[%s] , request body [%s], response body[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		subnetId = *response.Response.Subnet.SubnetId
	}
	return
}

func (me *VpcService) ModifySubnetAttribute(ctx context.Context, subnetId, name string, isMulticast bool) (errRet error) {
	logId := getLogId(ctx)
	request := vpc.NewModifySubnetAttributeRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var enableMulticast = map[bool]string{true: "true", false: "false"}[isMulticast]

	request.SubnetId = &subnetId
	request.SubnetName = &name
	request.EnableBroadcast = &enableMulticast
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().ModifySubnetAttribute(request)

	errRet = err
	if err == nil {
		log.Printf("[DEBUG]%s api[%s] , request body [%s], response body[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	}
	return
}

func (me *VpcService) DeleteSubnet(ctx context.Context, subnetId string) (errRet error) {

	logId := getLogId(ctx)
	request := vpc.NewDeleteSubnetRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.SubnetId = &subnetId
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().DeleteSubnet(request)

	errRet = err
	if err == nil {
		log.Printf("[DEBUG]%s api[%s] , request body [%s], response body[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	}
	return

}

func (me *VpcService) ReplaceRouteTableAssociation(ctx context.Context, subnetId string, routeTableId string) (errRet error) {
	logId := getLogId(ctx)
	request := vpc.NewReplaceRouteTableAssociationRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.SubnetId = &subnetId
	request.RouteTableId = &routeTableId
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().ReplaceRouteTableAssociation(request)

	errRet = err
	if err == nil {
		log.Printf("[DEBUG]%s api[%s] , request body [%s], response body[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	}
	return
}

func (me *VpcService) IsRouteTableInVpc(ctx context.Context, routeTableId, vpcId string) (info VpcRouteTableBasicInfo, has int, errRet error) {

	infos, err := me.DescribeRouteTables(ctx, routeTableId, "", vpcId)
	if err != nil {
		errRet = err
		return
	}
	has = len(infos)
	if has > 0 {
		info = infos[0]
	}
	return

}

func (me *VpcService) DescribeRouteTable(ctx context.Context, routeTableId string) (info VpcRouteTableBasicInfo, has int, errRet error) {

	infos, err := me.DescribeRouteTables(ctx, routeTableId, "", "")
	if err != nil {
		errRet = err
		return
	}

	has = len(infos)

	if has == 0 {
		return
	}
	info = infos[0]
	return
}
func (me *VpcService) DescribeRouteTables(ctx context.Context, routeTableId, routeTableName, vpcId string) (infos []VpcRouteTableBasicInfo, errRet error) {

	logId := getLogId(ctx)
	request := vpc.NewDescribeRouteTablesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	infos = make([]VpcRouteTableBasicInfo, 0, 100)
	var offset = 0
	var limit = 100
	var total = -1
	var hasTableMap = map[string]bool{}

	var filters []*vpc.Filter
	if routeTableId != "" {
		filters = me.fillFilter(filters, "route-table-id", routeTableId)
	}
	if vpcId != "" {
		filters = me.fillFilter(filters, "vpc-id", vpcId)
	}
	if routeTableName != "" {
		filters = me.fillFilter(filters, "route-table-name", routeTableName)
	}
	if len(filters) > 0 {
		request.Filters = filters
	}

getMoreData:
	if total >= 0 {
		if offset >= total {
			return
		}
	}
	var strLimit = fmt.Sprintf("%d", limit)
	request.Limit = &strLimit

	var strOffset = fmt.Sprintf("%d", offset)
	request.Offset = &strOffset
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().DescribeRouteTables(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] , request body [%s], response body[%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if total < 0 {
		total = int(*response.Response.TotalCount)
	}

	if len(response.Response.RouteTableSet) > 0 {
		offset += limit
	} else {
		// get empty Vpcinfo,we're done
		return
	}
	for _, item := range response.Response.RouteTableSet {
		var basicInfo VpcRouteTableBasicInfo
		basicInfo.createTime = *item.CreatedTime
		basicInfo.isDefault = *item.Main
		basicInfo.name = *item.RouteTableName
		basicInfo.routeTableId = *item.RouteTableId
		basicInfo.vpcId = *item.VpcId

		basicInfo.subnetIds = make([]string, 0, len(item.AssociationSet))
		for _, v := range item.AssociationSet {
			basicInfo.subnetIds = append(basicInfo.subnetIds, *v.SubnetId)
		}

		basicInfo.entryInfos = make([]VpcRouteEntryBasicInfo, 0, len(item.RouteSet))

		for _, v := range item.RouteSet {
			var entry VpcRouteEntryBasicInfo
			entry.destinationCidr = *v.DestinationCidrBlock
			entry.nextBub = *v.GatewayId
			entry.nextType = *v.GatewayType
			entry.description = *v.RouteDescription
			entry.routeEntryId = int64(*v.RouteId)
			entry.entryType = *v.RouteType
			basicInfo.entryInfos = append(basicInfo.entryInfos, entry)
		}
		if hasTableMap[basicInfo.routeTableId] {
			errRet = fmt.Errorf("get repeated route_table_id[%s] when doing DescribeRouteTables", basicInfo.routeTableId)
			return
		}
		hasTableMap[basicInfo.routeTableId] = true
		infos = append(infos, basicInfo)
	}
	goto getMoreData

}

func (me *VpcService) CreateRouteTable(ctx context.Context, name, vpcId string) (routeTableId string, errRet error) {

	logId := getLogId(ctx)
	request := vpc.NewCreateRouteTableRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	if vpcId == "" {
		errRet = fmt.Errorf("CreateRouteTable can not invoke by empty vpc_id.")
		return
	}
	request.VpcId = &vpcId
	request.RouteTableName = &name
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().CreateRouteTable(request)
	errRet = err
	if err == nil {
		log.Printf("[DEBUG]%s api[%s] , request body [%s], response body[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		routeTableId = *response.Response.RouteTable.RouteTableId
	}
	return
}

func (me *VpcService) DeleteRouteTable(ctx context.Context, routeTableId string) (errRet error) {

	logId := getLogId(ctx)
	request := vpc.NewDeleteRouteTableRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	if routeTableId == "" {
		errRet = fmt.Errorf("DeleteRouteTable can not invoke by empty routeTableId.")
		return
	}
	request.RouteTableId = &routeTableId
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().DeleteRouteTable(request)
	errRet = err
	if err == nil {
		log.Printf("[DEBUG]%s api[%s] , request body [%s], response body[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	}

	return
}

func (me *VpcService) ModifyRouteTableAttribute(ctx context.Context, routeTableId string, name string) (errRet error) {

	logId := getLogId(ctx)
	request := vpc.NewModifyRouteTableAttributeRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	if routeTableId == "" {
		errRet = fmt.Errorf("ModifyRouteTableAttribute can not invoke by empty routeTableId.")
		return
	}
	request.RouteTableId = &routeTableId
	request.RouteTableName = &name
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().ModifyRouteTableAttribute(request)
	errRet = err
	if err == nil {
		log.Printf("[DEBUG]%s api[%s] , request body [%s], response body[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	}

	return
}

func (me *VpcService) GetRouteId(ctx context.Context,
	routeTableId, destinationCidrBlock, nextType, nextHub, description string) (entryId int64, errRet error) {

	logId := getLogId(ctx)

	info, has, err := me.DescribeRouteTable(ctx, routeTableId)
	if err != nil {
		errRet = err
		return
	}
	if has == 0 {
		errRet = fmt.Errorf("not fonud the  route table of this  route entry")
		return
	}

	if has != 1 {
		errRet = fmt.Errorf("one routeTableId id get %d routeTableId infos", has)
		return
	}

	for _, v := range info.entryInfos {

		if v.destinationCidr == destinationCidrBlock && v.nextType == nextType && v.nextBub == nextHub {
			entryId = v.routeEntryId
			return
		}
	}
	errRet = fmt.Errorf("not found  route entry id from route table [%s]", routeTableId)

	for _, v := range info.entryInfos {
		log.Printf("%s[WARN] GetRouteId [%+v] vs [%+v],[%+v] vs [%+v],[%+v] vs [%+v]   %+v\n",
			logId,
			v.destinationCidr,
			destinationCidrBlock,
			v.nextType,
			nextType,
			v.nextBub,
			nextHub,
			v.destinationCidr == destinationCidrBlock && v.nextType == nextType && v.nextBub == nextHub)
	}

	return

}

func (me *VpcService) DeleteRoutes(ctx context.Context, routeTableId string, entryId uint64) (errRet error) {

	logId := getLogId(ctx)
	request := vpc.NewDeleteRoutesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	if routeTableId == "" {
		errRet = fmt.Errorf("DeleteRoutes can not invoke by empty routeTableId.")
		return
	}

	request.RouteTableId = &routeTableId
	var route vpc.Route
	route.RouteId = &entryId
	request.Routes = []*vpc.Route{&route}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().DeleteRoutes(request)
	errRet = err
	if err == nil {
		log.Printf("[DEBUG]%s api[%s] , request body [%s], response body[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	}
	return
}

func (me *VpcService) CreateRoutes(ctx context.Context,
	routeTableId, destinationCidrBlock, nextType, nextHub, description string) (entryId int64, errRet error) {

	logId := getLogId(ctx)
	request := vpc.NewCreateRoutesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	if routeTableId == "" {
		errRet = fmt.Errorf("CreateRoutes can not invoke by empty routeTableId.")
		return
	}
	request.RouteTableId = &routeTableId
	var route vpc.Route
	route.DestinationCidrBlock = &destinationCidrBlock
	route.RouteDescription = &description
	route.GatewayType = &nextType
	route.GatewayId = &nextHub
	request.Routes = []*vpc.Route{&route}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().CreateRoutes(request)
	errRet = err
	if err == nil {
		log.Printf("[DEBUG]%s api[%s] , request body [%s], response body[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	} else {
		return
	}

	entryId, errRet = me.GetRouteId(ctx, routeTableId, destinationCidrBlock, nextType, nextHub, description)

	if errRet != nil {
		time.Sleep(3 * time.Second)
		entryId, errRet = me.GetRouteId(ctx, routeTableId, destinationCidrBlock, nextType, nextHub, description)
	}

	if errRet != nil {
		time.Sleep(5 * time.Second)
		entryId, errRet = me.GetRouteId(ctx, routeTableId, destinationCidrBlock, nextType, nextHub, description)
	}

	/*
		if *(response.Response.TotalCount) != 1 {
			errRet = fmt.Errorf("CreateRoutes  return %d routeTable . but we only request 1.", *response.Response.TotalCount)
			return
		}

		if len(response.Response.RouteTableSet) != 1 {
			errRet = fmt.Errorf("CreateRoutes  return %d routeTable  info . but we only request 1.", len(response.Response.RouteTableSet))
			return
		}

		if len(response.Response.RouteTableSet[0].RouteSet) != 1 {
			errRet = fmt.Errorf("CreateRoutes  return %d routeTableSet  info . but we only create 1.", len(response.Response.RouteTableSet[0].RouteSet))
			return
		}

		entryId = int64(*response.Response.RouteTableSet[0].RouteSet[0].RouteId)
	*/

	return
}

func (me *VpcService) CreateSecurityGroup(ctx context.Context, name, desc string, projectId *int) (id string, err error) {
	logId := getLogId(ctx)

	request := vpc.NewCreateSecurityGroupRequest()

	request.GroupName = &name
	request.GroupDescription = &desc

	if projectId != nil {
		request.ProjectId = common.StringPtr(strconv.Itoa(*projectId))
	}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().CreateSecurityGroup(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
			logId, request.GetAction(), request.ToJsonString(), err)
		return "", err
	}

	if response.Response.SecurityGroup == nil || response.Response.SecurityGroup.SecurityGroupId == nil {
		return "", errors.New("response is nil")
	}

	return *response.Response.SecurityGroup.SecurityGroupId, nil
}

func (me *VpcService) DescribeSecurityGroup(ctx context.Context, id string) (sg *vpc.SecurityGroup, has int, err error) {
	logId := getLogId(ctx)

	request := vpc.NewDescribeSecurityGroupsRequest()

	request.SecurityGroupIds = []*string{&id}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().DescribeSecurityGroups(request)
	if err != nil {
		if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
			if sdkError.Code == "ResourceNotFound" {
				return nil, 0, nil
			}
		}

		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
			logId, request.GetAction(), request.ToJsonString(), err)
		return nil, 0, err
	}

	if len(response.Response.SecurityGroupSet) == 0 {
		return nil, 0, nil
	}

	return response.Response.SecurityGroupSet[0], len(response.Response.SecurityGroupSet), nil
}

func (me *VpcService) ModifySecurityGroup(ctx context.Context, id string, newName, newDesc *string) error {
	logId := getLogId(ctx)

	request := vpc.NewModifySecurityGroupAttributeRequest()

	request.SecurityGroupId = &id
	request.GroupName = newName
	request.GroupDescription = newDesc
	ratelimit.Check(request.GetAction())
	_, err := me.client.UseVpcClient().ModifySecurityGroupAttribute(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
			logId, request.GetAction(), request.ToJsonString(), err)
		return err
	}

	return nil
}

func (me *VpcService) DeleteSecurityGroup(ctx context.Context, id string) error {
	logId := getLogId(ctx)

	request := vpc.NewDeleteSecurityGroupRequest()

	request.SecurityGroupId = &id
	ratelimit.Check(request.GetAction())
	if _, err := me.client.UseVpcClient().DeleteSecurityGroup(request); err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
			logId, request.GetAction(), request.ToJsonString(), err)
		return err
	}

	return nil
}

func (me *VpcService) DescribeSecurityGroupsAssociate(ctx context.Context, ids []string) ([]*vpc.SecurityGroupAssociationStatistics, error) {
	logId := getLogId(ctx)

	request := vpc.NewDescribeSecurityGroupAssociationStatisticsRequest()
	request.SecurityGroupIds = common.StringPtrs(ids)
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().DescribeSecurityGroupAssociationStatistics(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
			logId, request.GetAction(), request.ToJsonString(), err)
		return nil, err
	}

	return response.Response.SecurityGroupAssociationStatisticsSet, nil
}

func (me *VpcService) CreateSecurityGroupPolicy(ctx context.Context, info securityGroupRuleBasicInfo) (ruleId string, err error) {
	logId := getLogId(ctx)

	createRequest := vpc.NewCreateSecurityGroupPoliciesRequest()
	createRequest.SecurityGroupId = &info.SgId

	createRequest.SecurityGroupPolicySet = new(vpc.SecurityGroupPolicySet)

	policy := new(vpc.SecurityGroupPolicy)

	policy.CidrBlock = info.CidrIp
	policy.SecurityGroupId = info.SourceSgId
	if info.Protocol != nil {
		policy.Protocol = common.StringPtr(strings.ToUpper(*info.Protocol))
	}

	policy.Port = info.PortRange
	policy.PolicyDescription = info.Description
	policy.Action = common.StringPtr(strings.ToUpper(info.Action))

	switch strings.ToLower(info.PolicyType) {
	case "ingress":
		createRequest.SecurityGroupPolicySet.Ingress = []*vpc.SecurityGroupPolicy{policy}

	case "egress":
		createRequest.SecurityGroupPolicySet.Egress = []*vpc.SecurityGroupPolicy{policy}
	}
	ratelimit.Check(createRequest.GetAction())
	if _, err := me.client.UseVpcClient().CreateSecurityGroupPolicies(createRequest); err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
			logId, createRequest.GetAction(), createRequest.ToJsonString(), err)
		return "", err
	}

	if info.CidrIp == nil {
		info.CidrIp = common.StringPtr("")
	}
	if info.Protocol == nil {
		info.Protocol = common.StringPtr("ALL")
	}
	if info.PortRange == nil {
		info.PortRange = common.StringPtr("ALL")
	}
	if info.SourceSgId == nil {
		info.SourceSgId = common.StringPtr("")
	}

	ruleId, err = buildSecurityGroupRuleId(info)
	if err != nil {
		return "", fmt.Errorf("build rule id error, reason: %v", err)
	}

	return ruleId, nil
}

func (me *VpcService) DescribeSecurityGroupPolicy(ctx context.Context, ruleId string) (sgId string, policyType string, policy *vpc.SecurityGroupPolicy, errRet error) {
	logId := getLogId(ctx)

	info, err := parseSecurityGroupRuleId(ruleId)
	if err != nil {
		errRet = err
		return
	}

	request := vpc.NewDescribeSecurityGroupPoliciesRequest()
	request.SecurityGroupId = &info.SgId
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().DescribeSecurityGroupPolicies(request)
	if err != nil {
		if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
			// if security group does not exist, security group rule does not exist too
			if sdkError.Code == "ResourceNotFound" {
				return
			}
		}

		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
			logId, request.GetAction(), request.ToJsonString(), err)
		errRet = err
		return
	}

	policySet := response.Response.SecurityGroupPolicySet

	if policySet == nil {
		log.Printf("[DEBUG]%s policy set is nil", logId)
		return
	}

	var policies []*vpc.SecurityGroupPolicy

	switch strings.ToLower(info.PolicyType) {
	case "ingress":
		policies = policySet.Ingress

	case "egress":
		policies = policySet.Egress
	}

	for _, pl := range policies {
		if comparePolicyAndSecurityGroupInfo(pl, info) {
			policy = pl
			break
		}
	}

	if policy == nil {
		log.Printf("[DEBUG]%s can't find security group rule, maybe user modify rules on web console", logId)
		return
	}

	return info.SgId, info.PolicyType, policy, nil
}

func (me *VpcService) DeleteSecurityGroupPolicy(ctx context.Context, ruleId string) error {
	logId := getLogId(ctx)

	info, err := parseSecurityGroupRuleId(ruleId)
	if err != nil {
		return err
	}

	request := vpc.NewDeleteSecurityGroupPoliciesRequest()
	request.SecurityGroupId = &info.SgId
	request.SecurityGroupPolicySet = new(vpc.SecurityGroupPolicySet)

	policy := new(vpc.SecurityGroupPolicy)
	policy.Action = common.StringPtr(strings.ToUpper(info.Action))

	if *info.CidrIp != "" {
		policy.CidrBlock = info.CidrIp
	}

	if *info.Protocol != "ALL" {
		policy.Protocol = common.StringPtr(strings.ToUpper(*info.Protocol))
	}

	if *info.PortRange != "ALL" {
		policy.Port = info.PortRange
	}

	if *info.SourceSgId != "" {
		policy.SecurityGroupId = info.SourceSgId
	}

	switch strings.ToLower(info.PolicyType) {
	case "ingress":
		request.SecurityGroupPolicySet.Ingress = []*vpc.SecurityGroupPolicy{policy}

	case "egress":
		request.SecurityGroupPolicySet.Egress = []*vpc.SecurityGroupPolicy{policy}
	}
	ratelimit.Check(request.GetAction())
	if _, err := me.client.UseVpcClient().DeleteSecurityGroupPolicies(request); err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
			logId, request.GetAction(), request.ToJsonString(), err)
		return err
	}

	return nil
}

func (me *VpcService) ModifySecurityGroupPolicy(ctx context.Context, ruleId string, desc *string) error {
	logId := getLogId(ctx)

	info, err := parseSecurityGroupRuleId(ruleId)
	if err != nil {
		return err
	}

	request := vpc.NewReplaceSecurityGroupPolicyRequest()
	request.SecurityGroupId = &info.SgId
	request.SecurityGroupPolicySet = new(vpc.SecurityGroupPolicySet)

	policy := &vpc.SecurityGroupPolicy{
		Action:            &info.Action,
		CidrBlock:         info.CidrIp,
		Protocol:          info.Protocol,
		Port:              info.PortRange,
		SecurityGroupId:   info.SourceSgId,
		PolicyDescription: desc,
	}

	switch info.PolicyType {
	case "ingress":
		request.SecurityGroupPolicySet.Ingress = []*vpc.SecurityGroupPolicy{policy}

	case "egress":
		request.SecurityGroupPolicySet.Egress = []*vpc.SecurityGroupPolicy{policy}
	}
	ratelimit.Check(request.GetAction())
	if _, err := me.client.UseVpcClient().ReplaceSecurityGroupPolicy(request); err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
			logId, request.GetAction(), request.ToJsonString(), err)
		return err
	}

	return nil
}

func (me *VpcService) DescribeSecurityGroups(ctx context.Context, sgId, sgName *string, projectId *int) (sgs []*vpc.SecurityGroup, err error) {
	logId := getLogId(ctx)

	request := vpc.NewDescribeSecurityGroupsRequest()

	if sgId != nil {
		request.SecurityGroupIds = []*string{sgId}
	} else {
		if sgName != nil {
			request.Filters = append(request.Filters, &vpc.Filter{
				Name:   common.StringPtr("security-group-name"),
				Values: []*string{sgName},
			})
		}

		if projectId != nil {
			request.Filters = append(request.Filters, &vpc.Filter{
				Name:   common.StringPtr("project-id"),
				Values: []*string{common.StringPtr(strconv.Itoa(*projectId))},
			})
		}
	}

	offset := 0
	pageSize := 50
	sgs = make([]*vpc.SecurityGroup, 0)

	for {
		request.Offset = common.StringPtr(strconv.Itoa(offset))
		request.Limit = common.StringPtr(strconv.Itoa(pageSize))
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseVpcClient().DescribeSecurityGroups(request)
		if err != nil {
			if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "ResourceNotFound" {
					break
				}
			}

			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return nil, err
		}

		set := response.Response.SecurityGroupSet
		if len(set) == 0 {
			break
		}
		sgs = append(sgs, set...)

		if len(set) < pageSize {
			break
		}

		offset += pageSize
	}

	return sgs, nil
}

type securityGroupRuleBasicInfo struct {
	SgId        string  `json:"sg_id"`
	PolicyType  string  `json:"policy_type"`
	CidrIp      *string `json:"cidr_ip,omitempty"`
	Protocol    *string `json:"protocol"`
	PortRange   *string `json:"port_range"`
	Action      string  `json:"action"`
	SourceSgId  *string `json:"source_sg_id"`
	Description *string `json:"description,omitempty"`
}

// Build an ID for a Security Group Rule (new version)
func buildSecurityGroupRuleId(info securityGroupRuleBasicInfo) (ruleId string, err error) {
	b, err := json.Marshal(info)
	if err != nil {
		return "", err
	}

	log.Printf("[DEBUG] build rule is %s", string(b))

	return base64.StdEncoding.EncodeToString(b), nil
}

// Parse Security Group Rule ID
func parseSecurityGroupRuleId(ruleId string) (info securityGroupRuleBasicInfo, errRet error) {
	log.Printf("[DEBUG] parseSecurityGroupRuleId before: %v", ruleId)

	// new version ID
	if b, err := base64.StdEncoding.DecodeString(ruleId); err == nil {
		errRet = json.Unmarshal(b, &info)
		return
	}

	// old version ID
	m := make(map[string]string)
	ruleQueryStrings := strings.Split(ruleId, "&")
	if len(ruleQueryStrings) == 0 {
		errRet = errors.New("ruleId is invalid")
		return
	}
	for _, str := range ruleQueryStrings {
		arr := strings.Split(str, "=")
		if len(arr) != 2 {
			errRet = errors.New("ruleId is invalid")
			return
		}
		m[arr[0]] = arr[1]
	}

	info.SgId = m["sgId"]
	info.PolicyType = m["direction"]
	info.Action = m["action"]
	if m["sourceSgid"] == "" {
		info.CidrIp = common.StringPtr(m["cidrIp"])
	} else {
		info.CidrIp = common.StringPtr("")
	}
	info.SourceSgId = common.StringPtr(m["sourceSgid"])
	info.Protocol = common.StringPtr(m["ipProtocol"])
	info.PortRange = common.StringPtr(m["portRange"])
	info.Description = common.StringPtr(m["description"])

	log.Printf("[DEBUG] parseSecurityGroupRuleId after: %v", info)
	return
}

func comparePolicyAndSecurityGroupInfo(policy *vpc.SecurityGroupPolicy, info securityGroupRuleBasicInfo) bool {
	// policy.CidrBlock always not nil
	if *policy.CidrBlock != *info.CidrIp {
		return false
	}

	// policy.Port always not nil
	if *policy.Port != *info.PortRange {
		return false
	}

	// policy.Protocol always not nil
	if !strings.EqualFold(*policy.Protocol, *info.Protocol) {
		return false
	}

	// policy.SecurityGroupId always not nil
	if *policy.SecurityGroupId != *info.SourceSgId {
		return false
	}

	if !strings.EqualFold(*policy.Action, info.Action) {
		return false
	}

	return true
}

func (me *VpcService) CreateEni(
	ctx context.Context,
	name, vpcId, subnetId, desc string,
	securityGroups []string,
	ipv4Count *int,
	ipv4s []VpcEniIP,
) (id string, err error) {
	logId := getLogId(ctx)
	client := me.client.UseVpcClient()

	createRequest := vpc.NewCreateNetworkInterfaceRequest()
	createRequest.NetworkInterfaceName = &name
	createRequest.VpcId = &vpcId
	createRequest.SubnetId = &subnetId
	createRequest.NetworkInterfaceDescription = &desc

	if len(securityGroups) > 0 {
		createRequest.SecurityGroupIds = common.StringPtrs(securityGroups)
	}

	if ipv4Count != nil {
		createRequest.SecondaryPrivateIpAddressCount = intToPointer(*ipv4Count - 1)
	}

	for _, ipv4 := range ipv4s {
		createRequest.PrivateIpAddresses = append(createRequest.PrivateIpAddresses, &vpc.PrivateIpAddressSpecification{
			PrivateIpAddress: stringToPointer(ipv4.ip.String()),
			Primary:          boolToPointer(ipv4.primary),
			Description:      ipv4.desc,
		})
	}

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(createRequest.GetAction())

		response, err := client.CreateNetworkInterface(createRequest)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, createRequest.GetAction(), createRequest.ToJsonString(), err)
			return retryError(err)
		}

		eni := response.Response.NetworkInterface

		if eni == nil {
			err := fmt.Errorf("api[%s] eni is nil", createRequest.GetAction())
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}

		if eni.NetworkInterfaceId == nil {
			err := fmt.Errorf("api[%s] eni id is nil", createRequest.GetAction())
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}

		id = *eni.NetworkInterfaceId

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s create eni failed, reason: %v", logId, err)
		return "", err
	}

	if err := waitEniReady(ctx, id, client); err != nil {
		log.Printf("[CRITAL]%s create eni failed, reason: %v", logId, err)
		return "", err
	}

	return
}

func (me *VpcService) AssignIpv6ToEni(ctx context.Context, id string, ipv6s []VpcEniIP, ipv6Count *int) error {
	logId := getLogId(ctx)
	client := me.client.UseVpcClient()

	assignRequest := vpc.NewAssignIpv6AddressesRequest()
	assignRequest.NetworkInterfaceId = &id

	if ipv6Count != nil {
		assignRequest.Ipv6AddressCount = intToPointer(*ipv6Count)
	}

	for _, ipv6 := range ipv6s {
		assignRequest.Ipv6Addresses = append(assignRequest.Ipv6Addresses, &vpc.Ipv6Address{
			Address:     stringToPointer(ipv6.ip.String()),
			Primary:     boolToPointer(ipv6.primary),
			Description: ipv6.desc,
		})
	}

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(assignRequest.GetAction())

		if _, err := client.AssignIpv6Addresses(assignRequest); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, assignRequest.GetAction(), assignRequest.ToJsonString(), err)
			return retryError(err)
		}

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s assign ipv6 to eni failed, reason: %v", logId, err)
		return err
	}

	if err := waitEniReady(ctx, id, client); err != nil {
		log.Printf("[CRITAL]%s assign ipv6 to eni failed, reason: %v", logId, err)
		return err
	}

	return nil
}

func (me *VpcService) describeEnis(
	ctx context.Context,
	ids []string,
	vpcId, subnetId, id, cvmId, sgId, name, desc, ipv4 *string,
	tags map[string]string,
) (enis []*vpc.NetworkInterface, err error) {
	logId := getLogId(ctx)

	request := vpc.NewDescribeNetworkInterfacesRequest()

	if len(ids) > 0 {
		request.NetworkInterfaceIds = common.StringPtrs(ids)
	}

	if vpcId != nil {
		request.Filters = append(request.Filters, &vpc.Filter{
			Name:   stringToPointer("vpc-id"),
			Values: []*string{vpcId},
		})
	}

	if subnetId != nil {
		request.Filters = append(request.Filters, &vpc.Filter{
			Name:   stringToPointer("subnet-id"),
			Values: []*string{subnetId},
		})
	}

	if id != nil {
		request.Filters = append(request.Filters, &vpc.Filter{
			Name:   stringToPointer("network-interface-id"),
			Values: []*string{id},
		})
	}

	if cvmId != nil {
		request.Filters = append(request.Filters, &vpc.Filter{
			Name:   stringToPointer("attachment.instance-id"),
			Values: []*string{cvmId},
		})
	}

	if sgId != nil {
		request.Filters = append(request.Filters, &vpc.Filter{
			Name:   stringToPointer("groups.security-group-id"),
			Values: []*string{sgId},
		})
	}

	if name != nil {
		request.Filters = append(request.Filters, &vpc.Filter{
			Name:   stringToPointer("network-interface-name"),
			Values: []*string{name},
		})
	}

	if desc != nil {
		request.Filters = append(request.Filters, &vpc.Filter{
			Name:   stringToPointer("network-interface-description"),
			Values: []*string{desc},
		})
	}

	if ipv4 != nil {
		request.Filters = append(request.Filters, &vpc.Filter{
			Name:   stringToPointer("address-ip"),
			Values: []*string{ipv4},
		})
	}

	for k, v := range tags {
		request.Filters = append(request.Filters, &vpc.Filter{
			Name:   stringToPointer("tag:" + k),
			Values: []*string{stringToPointer(v)},
		})
	}

	request.Limit = intToPointer(100)

	var offset uint64

	request.Offset = &offset

	count := 100
	for count == 100 {
		if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())

			response, err := me.client.UseVpcClient().DescribeNetworkInterfaces(request)
			if err != nil {
				count = 0

				if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
					if sdkError.Code == "ResourceNotFound" {
						return nil
					}
				}

				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
					logId, request.GetAction(), request.ToJsonString(), err)
				return retryError(err)
			}

			eniSet := response.Response.NetworkInterfaceSet
			count = len(eniSet)
			enis = append(enis, eniSet...)

			return nil
		}); err != nil {
			log.Printf("[CRITAL]%s read eni list failed, reason: %v", logId, err)
			return nil, err
		}

		offset += uint64(count)
	}

	return
}

func (me *VpcService) DescribeEniById(ctx context.Context, id string) (enis []*vpc.NetworkInterface, err error) {
	return me.describeEnis(ctx, []string{id}, nil, nil, nil, nil, nil, nil, nil, nil, nil)
}

func (me *VpcService) ModifyEniAttribute(ctx context.Context, id string, name, desc *string, sgs []string) error {
	logId := getLogId(ctx)
	client := me.client.UseVpcClient()

	request := vpc.NewModifyNetworkInterfaceAttributeRequest()
	request.NetworkInterfaceId = &id
	request.NetworkInterfaceName = name
	request.NetworkInterfaceDescription = desc
	request.SecurityGroupIds = common.StringPtrs(sgs)

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		if _, err := client.ModifyNetworkInterfaceAttribute(request); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return retryError(err)
		}

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s modify eni attribute failed, reason: %v", logId, err)
		return err
	}

	if err := waitEniReady(ctx, id, client); err != nil {
		log.Printf("[CRITAL]%s modify eni attribute failed, reason: %v", logId, err)
		return err
	}

	return nil
}

func (me *VpcService) UnAssignIpv4FromEni(ctx context.Context, id string, ipv4s []string) error {
	logId := getLogId(ctx)
	client := me.client.UseVpcClient()

	sort.Strings(ipv4s)

	unAssignRequest := vpc.NewUnassignPrivateIpAddressesRequest()
	unAssignRequest.NetworkInterfaceId = &id
	unAssignRequest.PrivateIpAddresses = make([]*vpc.PrivateIpAddressSpecification, 0, len(ipv4s))
	for _, ipv4 := range ipv4s {
		unAssignRequest.PrivateIpAddresses = append(unAssignRequest.PrivateIpAddresses, &vpc.PrivateIpAddressSpecification{
			PrivateIpAddress: stringToPointer(ipv4),
		})
	}

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(unAssignRequest.GetAction())

		if _, err := client.UnassignPrivateIpAddresses(unAssignRequest); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, unAssignRequest.GetAction(), unAssignRequest.ToJsonString(), err)
			return retryError(err)
		}

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s unassign ipv4 from eni failed, reason: %v", logId, err)
		return err
	}

	describeRequest := vpc.NewDescribeNetworkInterfacesRequest()
	describeRequest.NetworkInterfaceIds = []*string{&id}

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(describeRequest.GetAction())

		response, err := client.DescribeNetworkInterfaces(describeRequest)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, unAssignRequest.GetAction(), unAssignRequest.ToJsonString(), err)
			return retryError(err)
		}

		var eni *vpc.NetworkInterface

		for _, e := range response.Response.NetworkInterfaceSet {
			if e.NetworkInterfaceId == nil {
				err := fmt.Errorf("api[%s] eni id is nil", describeRequest.GetAction())
				log.Printf("[CRITAL]%s %v", logId, err)
				return resource.NonRetryableError(err)
			}

			if *e.NetworkInterfaceId == id {
				eni = e
				break
			}
		}

		if eni == nil {
			err := fmt.Errorf("api[%s] eni not found", describeRequest.GetAction())
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}

		// check unassigned ipv4 exists or not
		for _, ipv4 := range eni.PrivateIpAddressSet {
			if ipv4.PrivateIpAddress == nil {
				err := fmt.Errorf("api[%s] eni ipv4 ip is nil", describeRequest.GetAction())
				log.Printf("[CRITAL]%s %v", logId, err)
				return resource.NonRetryableError(err)
			}

			index := sort.SearchStrings(ipv4s, *ipv4.PrivateIpAddress)
			if index < len(ipv4s) && ipv4s[index] == *ipv4.PrivateIpAddress {
				err := errors.New("eni unassigned ipv4 still exists")
				log.Printf("[DEBUG]%s %v", logId, err)
				return resource.RetryableError(err)
			}
		}

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s unassign ipv4 from eni failed, reason: %v", logId, err)
		return err
	}

	if err := waitEniReady(ctx, id, client); err != nil {
		log.Printf("[CRITAL]%s unassign ipv4 from eni failed, reason: %v", logId, err)
		return err
	}

	return nil
}

func (me *VpcService) AssignIpv4ToEni(ctx context.Context, id string, ipv4s []VpcEniIP, ipv4Count *int) error {
	logId := getLogId(ctx)
	client := me.client.UseVpcClient()

	request := vpc.NewAssignPrivateIpAddressesRequest()
	request.NetworkInterfaceId = &id

	if ipv4Count != nil {
		request.SecondaryPrivateIpAddressCount = intToPointer(*ipv4Count)
	}

	if len(ipv4s) > 0 {
		request.PrivateIpAddresses = make([]*vpc.PrivateIpAddressSpecification, 0, len(ipv4s))
		for _, ipv4 := range ipv4s {
			request.PrivateIpAddresses = append(request.PrivateIpAddresses, &vpc.PrivateIpAddressSpecification{
				PrivateIpAddress: stringToPointer(ipv4.ip.String()),
				Primary:          boolToPointer(ipv4.primary),
				Description:      ipv4.desc,
			})
		}
	}

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		if _, err := client.AssignPrivateIpAddresses(request); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return retryError(err)
		}

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s assign ipv4 to eni failed, reason: %v", logId, err)
		return err
	}

	if err := waitEniReady(ctx, id, client); err != nil {
		log.Printf("[CRITAL]%s assign ipv4 to eni failed, reason: %v", logId, err)
		return err
	}

	return nil
}

func (me *VpcService) UnAssignIpv6FromEni(ctx context.Context, id string, ipv6s []string) error {
	logId := getLogId(ctx)
	client := me.client.UseVpcClient()

	sort.Strings(ipv6s)

	unAssignRequest := vpc.NewUnassignIpv6AddressesRequest()
	unAssignRequest.NetworkInterfaceId = &id
	unAssignRequest.Ipv6Addresses = make([]*vpc.Ipv6Address, 0, len(ipv6s))
	for _, ipv6 := range ipv6s {
		unAssignRequest.Ipv6Addresses = append(unAssignRequest.Ipv6Addresses, &vpc.Ipv6Address{
			Address: stringToPointer(ipv6),
		})
	}

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(unAssignRequest.GetAction())

		if _, err := client.UnassignIpv6Addresses(unAssignRequest); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, unAssignRequest.GetAction(), unAssignRequest.ToJsonString(), err)
			return retryError(err)
		}

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s unassign ipv6 from eni failed, reason: %v", logId, err)
		return err
	}

	describeRequest := vpc.NewDescribeNetworkInterfacesRequest()
	describeRequest.NetworkInterfaceIds = []*string{&id}

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(describeRequest.GetAction())

		response, err := client.DescribeNetworkInterfaces(describeRequest)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, describeRequest.GetAction(), describeRequest.ToJsonString(), err)
			return retryError(err)
		}

		var eni *vpc.NetworkInterface

		for _, e := range response.Response.NetworkInterfaceSet {
			if e.NetworkInterfaceId == nil {
				err := fmt.Errorf("api[%s] eni id is nil", describeRequest.GetAction())
				log.Printf("[CRITAL]%s %v", logId, err)
				return resource.NonRetryableError(err)
			}

			if *e.NetworkInterfaceId == id {
				eni = e
				break
			}
		}

		if eni == nil {
			err := fmt.Errorf("api[%s] eni not exists", describeRequest.GetAction())
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}

		for _, ipv6 := range eni.Ipv6AddressSet {
			if ipv6.Address == nil {
				err := fmt.Errorf("api[%s] eni ipv6 address is nil", describeRequest.GetAction())
				log.Printf("[CRITAL]%s %v", logId, err)
				return resource.NonRetryableError(err)
			}

			index := sort.SearchStrings(ipv6s, *ipv6.Address)
			if index < len(ipv6s) && ipv6s[index] == *ipv6.Address {
				err := errors.New("unassigned ipv6 still exists")
				log.Printf("[DEBUG]%s %v", logId, err)
				return resource.RetryableError(err)
			}
		}

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s unassign ipv6 from eni failed, reason: %v", logId, err)
		return err
	}

	return nil
}

func (me *VpcService) DeleteEni(ctx context.Context, id string) error {
	logId := getLogId(ctx)
	client := me.client.UseVpcClient()

	deleteRequest := vpc.NewDeleteNetworkInterfaceRequest()
	deleteRequest.NetworkInterfaceId = &id

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(deleteRequest.GetAction())

		if _, err := client.DeleteNetworkInterface(deleteRequest); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, deleteRequest.GetAction(), deleteRequest.ToJsonString(), err)
			return retryError(err)
		}

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s delete eni failed, reason: %v", logId, err)
		return err
	}

	describeRequest := vpc.NewDescribeNetworkInterfacesRequest()
	describeRequest.NetworkInterfaceIds = []*string{&id}

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(describeRequest.GetAction())

		response, err := client.DescribeNetworkInterfaces(describeRequest)
		if err != nil {
			if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "ResourceNotFound" {
					return nil
				}
			}

			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, describeRequest.GetAction(), describeRequest.ToJsonString(), err)
			return retryError(err)
		}

		for _, eni := range response.Response.NetworkInterfaceSet {
			if eni.NetworkInterfaceId == nil {
				err := fmt.Errorf("api[%s] eni id is nil", describeRequest.GetAction())
				log.Printf("[CRITAL]%s %v", logId, err)
				return resource.NonRetryableError(err)
			}

			if *eni.NetworkInterfaceId == id {
				err := errors.New("eni still exists")
				log.Printf("[DEBUG]%s %v", logId, err)
				return resource.RetryableError(err)
			}
		}

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s delete eni failed, reason: %v", logId, err)
		return err
	}

	return nil
}

func (me *VpcService) AttachEniToCvm(ctx context.Context, eniId, cvmId string) error {
	logId := getLogId(ctx)
	client := me.client.UseVpcClient()

	attachRequest := vpc.NewAttachNetworkInterfaceRequest()
	attachRequest.NetworkInterfaceId = &eniId
	attachRequest.InstanceId = &cvmId

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(attachRequest.GetAction())

		if _, err := client.AttachNetworkInterface(attachRequest); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, attachRequest.GetAction(), attachRequest.ToJsonString(), err)
			return retryError(err)
		}

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s attach eni to instance failed, reason: %v", logId, err)
		return err
	}

	describeRequest := vpc.NewDescribeNetworkInterfacesRequest()
	describeRequest.NetworkInterfaceIds = []*string{&eniId}

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(describeRequest.GetAction())

		response, err := client.DescribeNetworkInterfaces(describeRequest)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, describeRequest.GetAction(), describeRequest.ToJsonString(), err)
			return retryError(err)
		}

		var eni *vpc.NetworkInterface
		for _, e := range response.Response.NetworkInterfaceSet {
			if e.NetworkInterfaceId == nil {
				err := fmt.Errorf("api[%s] eni id is nil", describeRequest.GetAction())
				log.Printf("[CRITAL]%s %v", logId, err)
				return resource.NonRetryableError(err)
			}

			if *e.NetworkInterfaceId == eniId {
				eni = e
				break
			}
		}

		if eni == nil {
			err := fmt.Errorf("api[%s] eni not found", describeRequest.GetAction())
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}

		if eni.Attachment == nil {
			err := fmt.Errorf("api[%s] eni attachment is not ready", describeRequest.GetAction())
			log.Printf("[DEBUG]%s %v", logId, err)
			return resource.RetryableError(err)
		}

		if eni.Attachment.InstanceId == nil {
			err := fmt.Errorf("api[%s] eni attach instance id is nil", describeRequest.GetAction())
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}

		if *eni.Attachment.InstanceId != cvmId {
			err := fmt.Errorf("api[%s] eni attach instance id is not right", describeRequest.GetAction())
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}

		if eni.State == nil {
			err := fmt.Errorf("api[%s] eni state is nil", describeRequest.GetAction())
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}

		if *eni.State != ENI_STATE_AVAILABLE {
			err := errors.New("eni is not ready")
			log.Printf("[DEBUG]%s %v", logId, err)
			return resource.RetryableError(err)
		}

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s attach eni to instance failed, reason: %v", logId, err)
		return err
	}

	return nil
}

func (me *VpcService) DetachEniFromCvm(ctx context.Context, eniId, cvmId string) error {
	logId := getLogId(ctx)
	client := me.client.UseVpcClient()

	request := vpc.NewDetachNetworkInterfaceRequest()
	request.NetworkInterfaceId = &eniId
	request.InstanceId = &cvmId

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		if _, err := client.DetachNetworkInterface(request); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return retryError(err)
		}

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s detach eni from instance failed, reason: %v", logId, err)
		return err
	}

	if err := waitEniReady(ctx, eniId, client); err != nil {
		log.Printf("[CRITAL]%s detach eni from instance failed, reason: %v", logId, err)
		return err
	}

	return nil
}

func waitEniReady(ctx context.Context, id string, client *vpc.Client) error {
	logId := getLogId(ctx)

	request := vpc.NewDescribeNetworkInterfacesRequest()
	request.NetworkInterfaceIds = []*string{stringToPointer(id)}

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, err := client.DescribeNetworkInterfaces(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return retryError(err)
		}

		var eni *vpc.NetworkInterface
		for _, networkInterface := range response.Response.NetworkInterfaceSet {
			if networkInterface.NetworkInterfaceId == nil {
				err := fmt.Errorf("api[%s] eni id is nil", request.GetAction())
				log.Printf("[CRITAL]%s %v", logId, err)
				return resource.NonRetryableError(err)
			}

			if *networkInterface.NetworkInterfaceId == id {
				eni = networkInterface
				break
			}
		}

		if eni == nil {
			err := fmt.Errorf("api[%s] eni not exist", request.GetAction())
			log.Printf("[DEBUG]%s %v", logId, err)
			return resource.RetryableError(err)
		}

		if eni.State == nil {
			err := fmt.Errorf("api[%s] eni state is nil", request.GetAction())
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}

		if *eni.State != ENI_STATE_AVAILABLE {
			err := errors.New("eni is not available")
			log.Printf("[DEBUG]%s %v", logId, err)
			return resource.RetryableError(err)
		}

		for _, ipv4 := range eni.PrivateIpAddressSet {
			if ipv4.State == nil {
				err := fmt.Errorf("api[%s] eni ipv4 state is nil", request.GetAction())
				log.Printf("[CRITAL]%s %v", logId, err)
				return resource.NonRetryableError(err)
			}

			if *ipv4.State != ENI_IP_AVAILABLE {
				err := errors.New("eni ipv4 is not available")
				log.Printf("[DEBUG]%s %v", logId, err)
				return resource.RetryableError(err)
			}
		}

		for _, ipv6 := range eni.Ipv6AddressSet {
			if ipv6.State == nil {
				err := fmt.Errorf("api[%s] eni ipv6 state is nil", request.GetAction())
				log.Printf("[CRITAL]%s %v", logId, err)
				return resource.NonRetryableError(err)
			}

			if *ipv6.State != ENI_IP_AVAILABLE {
				err := errors.New("eni ipv6 is not available")
				log.Printf("[DEBUG]%s %v", logId, err)
				return resource.RetryableError(err)
			}
		}

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s eni is not available failed, reason: %v", logId, err)
		return err
	}

	return nil
}
