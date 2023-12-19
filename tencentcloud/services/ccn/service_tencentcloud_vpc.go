package ccn

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

var eipUnattachLocker = &sync.Mutex{}

/* For Adun Sake please DO NOT Declare the redundant Type STRUCT!! */
// VPC basic information
type VpcBasicInfo struct {
	vpcId                string
	name                 string
	cidr                 string
	isMulticast          bool
	isDefault            bool
	dnsServers           []string
	createTime           string
	tags                 []*vpc.Tag
	assistantCidrs       []string
	dockerAssistantCidrs []string
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
	enabled         bool
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

type VpcSecurityGroupLiteRule struct {
	action                  string
	cidrIp                  string
	port                    string
	protocol                string
	addressId               string
	addressGroupId          string
	securityGroupId         string
	protocolTemplateId      string
	protocolTemplateGroupId string
}

var securityGroupIdRE = regexp.MustCompile(`^sg-\w{8}$`)
var ipAddressIdRE = regexp.MustCompile(`^ipm-\w{8}$`)
var ipAddressGroupIdRE = regexp.MustCompile(`^ipmg-\w{8}$`)
var protocolTemplateRE = regexp.MustCompile(`^ppmg?-\w{8}$`)
var protocolTemplateIdRE = regexp.MustCompile(`^ppm-\w{8}$`)
var protocolTemplateGroupIdRE = regexp.MustCompile(`^ppmg-\w{8}$`)
var portRE = regexp.MustCompile(`^(\d{1,5},)*\d{1,5}$|^\d{1,5}-\d{1,5}$`)

// acl rule
type VpcACLRule struct {
	action   string
	cidrIp   string
	port     string
	protocol string
}

type VpcEniIP struct {
	ip      net.IP
	primary bool
	desc    *string
}

func (rule VpcSecurityGroupLiteRule) String() string {

	var source string

	if rule.cidrIp != "" {
		source = rule.cidrIp
	}
	if rule.securityGroupId != "" {
		source = rule.securityGroupId
	}
	if rule.addressId != "" {
		source = rule.addressId
	}
	if rule.addressGroupId != "" {
		source = rule.addressGroupId
	}

	protocol := rule.protocol

	if protocol == "" && rule.protocolTemplateId != "" {
		protocol = rule.protocolTemplateId
	} else if protocol == "" && rule.protocolTemplateGroupId != "" {
		protocol = rule.protocolTemplateGroupId
	}

	return fmt.Sprintf("%s#%s#%s#%s", rule.action, source, rule.port, protocol)
}

func getSecurityGroupPolicies(rules []VpcSecurityGroupLiteRule) []*vpc.SecurityGroupPolicy {
	policies := make([]*vpc.SecurityGroupPolicy, 0)

	for i := range rules {
		rule := rules[i]
		policy := &vpc.SecurityGroupPolicy{
			Action: &rule.action,
		}

		if rule.securityGroupId != "" {
			policy.SecurityGroupId = &rule.securityGroupId
		} else if rule.addressId != "" || rule.addressGroupId != "" {
			policy.AddressTemplate = &vpc.AddressTemplateSpecification{}
			if rule.addressId != "" {
				policy.AddressTemplate.AddressId = &rule.addressId
			}
			if rule.addressGroupId != "" {
				policy.AddressTemplate.AddressGroupId = &rule.addressGroupId
			}
		} else {
			policy.CidrBlock = &rule.cidrIp
		}

		usingProtocolTemplate := rule.protocolTemplateId != "" || rule.protocolTemplateGroupId != ""

		if usingProtocolTemplate {
			policy.ServiceTemplate = &vpc.ServiceTemplateSpecification{}
			if rule.protocolTemplateId != "" {
				policy.ServiceTemplate.ServiceId = &rule.protocolTemplateId
			}
			if rule.protocolTemplateGroupId != "" {
				policy.ServiceTemplate.ServiceGroupId = &rule.protocolTemplateGroupId
			}
		}

		if !usingProtocolTemplate {
			policy.Protocol = &rule.protocol
		}

		if !usingProtocolTemplate && rule.port != "" {
			policy.Port = &rule.port
		}

		policies = append(policies, policy)
	}
	return policies
}

func NewVpcService(client *connectivity.TencentCloudClient) VpcService {
	return VpcService{client: client}
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
	isMulticast bool, dnsServers []string, tags map[string]string) (vpcId string, isDefault bool, errRet error) {

	logId := tccommon.GetLogId(ctx)
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

	if len(tags) > 0 {
		for tagKey, tagValue := range tags {
			tag := vpc.Tag{
				Key:   helper.String(tagKey),
				Value: helper.String(tagValue),
			}
			request.Tags = append(request.Tags, &tag)
		}
	}

	var response *vpc.CreateVpcResponse
	if err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseVpcClient().CreateVpc(request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		response = result
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s create vpc failed, reason: %v", logId, err)
		errRet = err
		return
	}
	vpcId, isDefault = *response.Response.Vpc.VpcId, *response.Response.Vpc.IsDefault
	return
}

func (me *VpcService) DescribeVpc(ctx context.Context,
	vpcId string,
	tagKey string,
	cidrBlock string) (info VpcBasicInfo, has int, errRet error) {
	infos, err := me.DescribeVpcs(ctx, vpcId, "", nil, nil, tagKey, cidrBlock)
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

func (me *VpcService) DescribeVpcs(ctx context.Context,
	vpcId, name string,
	tags map[string]string,
	isDefaultPtr *bool,
	tagKey string,
	cidrBlock string) (infos []VpcBasicInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewDescribeVpcsRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	infos = make([]VpcBasicInfo, 0, 100)

	var (
		offset  = 0
		limit   = 100
		total   = -1
		hasVpc  = map[string]bool{}
		filters []*vpc.Filter
	)

	if vpcId != "" {
		filters = me.fillFilter(filters, "vpc-id", vpcId)
	}

	if name != "" {
		filters = me.fillFilter(filters, "vpc-name", name)
	}

	if tagKey != "" {
		filters = me.fillFilter(filters, "tag-key", tagKey)
	}

	if cidrBlock != "" {
		filters = me.fillFilter(filters, "cidr-block", cidrBlock)
	}

	if isDefaultPtr != nil {
		filters = me.fillFilter(filters, "is-default", map[bool]string{true: "true", false: "false"}[*isDefaultPtr])
	}

	for k, v := range tags {
		filters = me.fillFilter(filters, "tag:"+k, v)
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
	var response *vpc.DescribeVpcsResponse
	if err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseVpcClient().DescribeVpcs(request)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		response = result
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s read vpc failed, reason: %v", logId, err)
		return nil, err
	}

	if total < 0 {
		total = int(*response.Response.TotalCount)
	}

	if len(response.Response.VpcSet) > 0 {
		offset += limit
	} else {
		// get empty VpcInfo, we're done
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

		if len(item.AssistantCidrSet) > 0 {
			for i := range item.AssistantCidrSet {
				kind := item.AssistantCidrSet[i].AssistantType
				cidr := item.AssistantCidrSet[i].CidrBlock
				if kind != nil && *kind == 0 {
					basicInfo.assistantCidrs = append(basicInfo.assistantCidrs, *cidr)
				} else {
					basicInfo.dockerAssistantCidrs = append(basicInfo.dockerAssistantCidrs, *cidr)
				}
			}
		}

		if len(item.TagSet) > 0 {
			basicInfo.tags = item.TagSet
		}

		infos = append(infos, basicInfo)
	}
	goto getMoreData

}
func (me *VpcService) DescribeSubnet(ctx context.Context,
	subnetId string,
	isRemoteVpcSNAT *bool,
	tagKey,
	cidrBlock string) (info VpcSubnetBasicInfo, has int, errRet error) {
	infos, err := me.DescribeSubnets(ctx, subnetId, "", "", "", nil, nil, isRemoteVpcSNAT, tagKey, cidrBlock)
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

func (me *VpcService) DescribeSubnets(ctx context.Context,
	subnetId,
	vpcId,
	subnetName,
	zone string,
	tags map[string]string,
	isDefaultPtr *bool,
	isRemoteVpcSNAT *bool,
	tagKey,
	cidrBlock string) (infos []VpcSubnetBasicInfo, errRet error) {

	logId := tccommon.GetLogId(ctx)
	request := vpc.NewDescribeSubnetsRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var (
		offset    = 0
		limit     = 100
		total     = -1
		hasSubnet = map[string]bool{}
		filters   []*vpc.Filter
	)

	if subnetId != "" {
		filters = me.fillFilter(filters, "subnet-id", subnetId)
	}
	if vpcId != "" {
		filters = me.fillFilter(filters, "vpc-id", vpcId)
	}
	if subnetName != "" {
		filters = me.fillFilter(filters, "subnet-name", subnetName)
	}
	if zone != "" {
		filters = me.fillFilter(filters, "zone", zone)
	}

	if isDefaultPtr != nil {
		filters = me.fillFilter(filters, "is-default", map[bool]string{true: "true", false: "false"}[*isDefaultPtr])
	}

	if isRemoteVpcSNAT != nil {
		filters = me.fillFilter(filters, "is-remote-vpc-snat", map[bool]string{true: "true", false: "false"}[*isRemoteVpcSNAT])
	}

	if tagKey != "" {
		filters = me.fillFilter(filters, "tag-key", tagKey)
	}
	if cidrBlock != "" {
		filters = me.fillFilter(filters, "cidr-block", cidrBlock)
	}

	for k, v := range tags {
		filters = me.fillFilter(filters, "tag:"+k, v)
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
	var response *vpc.DescribeSubnetsResponse
	if err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseVpcClient().DescribeSubnets(request)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		response = result
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s read subnets failed, reason: %v", logId, err)
		return nil, err
	}
	log.Printf("[DEBUG]%s api[%s] , request body [%s], response body[%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if total < 0 {
		total = int(*response.Response.TotalCount)
	}

	if len(response.Response.SubnetSet) > 0 {
		offset += limit
	} else {
		// get empty subnet, we're done
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
	logId := tccommon.GetLogId(ctx)
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

	if err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err := me.client.UseVpcClient().ModifyVpcAttribute(request)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s modify vpc failed, reason: %v", logId, err)
		return err
	}

	return
}

func (me *VpcService) DeleteVpc(ctx context.Context, vpcId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
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

	if err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err := me.client.UseVpcClient().DeleteVpc(request)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s delete vpc failed, reason: %v", logId, err)
		return err
	}
	return

}

func (me *VpcService) CreateSubnet(ctx context.Context, vpcId, name, cidr, zone string, tags map[string]string) (subnetId string, errRet error) {
	logId := tccommon.GetLogId(ctx)
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

	if len(tags) > 0 {
		for tagKey, tagValue := range tags {
			tag := vpc.Tag{
				Key:   helper.String(tagKey),
				Value: helper.String(tagValue),
			}
			request.Tags = append(request.Tags, &tag)
		}
	}

	var response *vpc.CreateSubnetResponse
	if err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseVpcClient().CreateSubnet(request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		response = result
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s create subnet failed, reason: %v", logId, err)
		return "", err
	}

	subnetId = *response.Response.Subnet.SubnetId

	return
}

func (me *VpcService) ModifySubnetAttribute(ctx context.Context, subnetId, name string, isMulticast bool) (errRet error) {
	logId := tccommon.GetLogId(ctx)
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
	if err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err := me.client.UseVpcClient().ModifySubnetAttribute(request)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s modify subnet failed, reason: %v", logId, err)
		return err
	}
	return
}

func (me *VpcService) DeleteSubnet(ctx context.Context, subnetId string) (errRet error) {

	logId := tccommon.GetLogId(ctx)
	request := vpc.NewDeleteSubnetRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.SubnetId = &subnetId
	if err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err := me.client.UseVpcClient().DeleteSubnet(request)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s delete subnet failed, reason: %v", logId, err)
		return err
	}
	return

}

func (me *VpcService) ReplaceRouteTableAssociation(ctx context.Context, subnetId string, routeTableId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
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

	infos, err := me.DescribeRouteTables(ctx, routeTableId, "", vpcId, nil, nil, "")
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

	infos, err := me.DescribeRouteTables(ctx, routeTableId, "", "", nil, nil, "")
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
func (me *VpcService) DescribeRouteTables(ctx context.Context,
	routeTableId,
	routeTableName,
	vpcId string,
	tags map[string]string,
	associationMain *bool,
	tagKey string) (infos []VpcRouteTableBasicInfo, errRet error) {

	logId := tccommon.GetLogId(ctx)
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
	if associationMain != nil {
		filters = me.fillFilter(filters, "association.main", map[bool]string{true: "true", false: "false"}[*associationMain])
	}
	if tagKey != "" {
		filters = me.fillFilter(filters, "tag-key", tagKey)
	}
	for k, v := range tags {
		filters = me.fillFilter(filters, "tag:"+k, v)
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
		// get empty Vpcinfo, we're done
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
			entry.enabled = *v.Enabled
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

func (me *VpcService) CreateRouteTable(ctx context.Context, name, vpcId string, tags map[string]string) (routeTableId string, errRet error) {

	logId := tccommon.GetLogId(ctx)
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
	if len(tags) > 0 {
		for tagKey, tagValue := range tags {
			tag := vpc.Tag{
				Key:   helper.String(tagKey),
				Value: helper.String(tagValue),
			}
			request.Tags = append(request.Tags, &tag)
		}
	}
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

	logId := tccommon.GetLogId(ctx)
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

	logId := tccommon.GetLogId(ctx)
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

	logId := tccommon.GetLogId(ctx)

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

	logId := tccommon.GetLogId(ctx)
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
	routeTableId, destinationCidrBlock, nextType, nextHub, description string, enabled bool) (entryId int64, errRet error) {

	logId := tccommon.GetLogId(ctx)
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
	route.Enabled = &enabled
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

func (me *VpcService) SwitchRouteEnabled(ctx context.Context, routeTableId string, routeId uint64, enabled bool) error {
	if enabled {
		request := vpc.NewEnableRoutesRequest()
		request.RouteTableId = &routeTableId
		request.RouteIds = []*uint64{&routeId}
		return me.EnableRoutes(ctx, request)
	} else {
		request := vpc.NewDisableRoutesRequest()
		request.RouteTableId = &routeTableId
		request.RouteIds = []*uint64{&routeId}
		return me.DisableRoutes(ctx, request)
	}
}

func (me *VpcService) EnableRoutes(ctx context.Context, request *vpc.EnableRoutesRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().EnableRoutes(request)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
func (me *VpcService) DisableRoutes(ctx context.Context, request *vpc.DisableRoutesRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().DisableRoutes(request)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
func (me *VpcService) CreateSecurityGroup(ctx context.Context, name, desc string, projectId *int, tags map[string]string) (id string, err error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewCreateSecurityGroupRequest()

	request.GroupName = &name
	request.GroupDescription = &desc

	if projectId != nil {
		request.ProjectId = helper.String(strconv.Itoa(*projectId))
	}

	if len(tags) > 0 {
		for tagKey, tagValue := range tags {
			tag := vpc.Tag{
				Key:   helper.String(tagKey),
				Value: helper.String(tagValue),
			}
			request.Tags = append(request.Tags, &tag)
		}
	}

	if err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, err := me.client.UseVpcClient().CreateSecurityGroup(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return tccommon.RetryError(err)
		}

		if response.Response.SecurityGroup == nil || response.Response.SecurityGroup.SecurityGroupId == nil {
			err := fmt.Errorf("api[%s] return security group id is nil", request.GetAction())
			log.Printf("[CRITAL]%s %v", logId, err)
			return resource.NonRetryableError(err)
		}

		id = *response.Response.SecurityGroup.SecurityGroupId
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s create security group failed, reason: %v", logId, err)
		return "", err
	}

	return
}

func (me *VpcService) DescribeSecurityGroup(ctx context.Context, id string) (sg *vpc.SecurityGroup, err error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDescribeSecurityGroupsRequest()
	request.SecurityGroupIds = []*string{&id}

	if err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, err := me.client.UseVpcClient().DescribeSecurityGroups(request)
		if err != nil {
			if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "ResourceNotFound" {
					return nil
				}
			}

			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return tccommon.RetryError(err, tccommon.InternalError)
		}

		if len(response.Response.SecurityGroupSet) == 0 {
			return nil
		}

		sg = response.Response.SecurityGroupSet[0]

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s read security group failed, reason: %v", logId, err)
		return nil, err
	}

	return
}

func (me *VpcService) ModifySecurityGroup(ctx context.Context, id string, newName, newDesc *string) error {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewModifySecurityGroupAttributeRequest()

	request.SecurityGroupId = &id
	request.GroupName = newName
	request.GroupDescription = newDesc

	if err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err := me.client.UseVpcClient().ModifySecurityGroupAttribute(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return tccommon.RetryError(err)
		}

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s modify security group failed, reason: %v", logId, err)
		return err
	}

	return nil
}

func (me *VpcService) DeleteSecurityGroup(ctx context.Context, id string) error {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDeleteSecurityGroupRequest()
	request.SecurityGroupId = &id

	if err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		if _, err := me.client.UseVpcClient().DeleteSecurityGroup(request); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return tccommon.RetryError(err)
		}

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s delete security group failed, reason: %v", logId, err)
		return err
	}

	return nil
}

func (me *VpcService) DescribeSecurityGroupsAssociate(ctx context.Context, ids []string) ([]*vpc.SecurityGroupAssociationStatistics, error) {
	logId := tccommon.GetLogId(ctx)

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

// Deprecated: the redundant type struct cause cause unnecessary mental burden, use sdk request directly
func (me *VpcService) CreateSecurityGroupPolicy(ctx context.Context, info securityGroupRuleBasicInfoWithPolicyIndex) (ruleId string, err error) {
	logId := tccommon.GetLogId(ctx)

	createRequest := vpc.NewCreateSecurityGroupPoliciesRequest()
	createRequest.SecurityGroupId = &info.SgId

	createRequest.SecurityGroupPolicySet = new(vpc.SecurityGroupPolicySet)

	policy := new(vpc.SecurityGroupPolicy)

	policy.CidrBlock = info.CidrIp
	policy.SecurityGroupId = info.SourceSgId
	policy.AddressTemplate = &vpc.AddressTemplateSpecification{}
	if info.AddressTemplateId != nil && *info.AddressTemplateId != "" {
		policy.AddressTemplate.AddressId = info.AddressTemplateId
	}
	if info.AddressTemplateGroupId != nil && *info.AddressTemplateGroupId != "" {
		policy.AddressTemplate.AddressGroupId = info.AddressTemplateGroupId
	}

	policy.ServiceTemplate = &vpc.ServiceTemplateSpecification{}
	if info.ProtocolTemplateId != nil && *info.ProtocolTemplateId != "" {
		policy.ServiceTemplate.ServiceId = info.ProtocolTemplateId
	}
	if info.ProtocolTemplateGroupId != nil && *info.ProtocolTemplateGroupId != "" {
		policy.ServiceTemplate.ServiceGroupId = info.ProtocolTemplateGroupId
	}

	if info.Protocol != nil {
		policy.Protocol = common.StringPtr(strings.ToUpper(*info.Protocol))
	}
	policy.PolicyIndex = helper.Int64(info.PolicyIndex)
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

	ruleId, err = buildSecurityGroupRuleId(info.securityGroupRuleBasicInfo)
	if err != nil {
		return "", fmt.Errorf("build rule id error, reason: %v", err)
	}

	return ruleId, nil
}

func (me *VpcService) CreateSecurityGroupPolicies(ctx context.Context, request *vpc.CreateSecurityGroupPoliciesRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().CreateSecurityGroupPolicies(request)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

// Deprecated: use DescribeSecurityGroupPolicies instead
func (me *VpcService) DescribeSecurityGroupPolicy(ctx context.Context, ruleId string) (sgId string, policyType string, policy *vpc.SecurityGroupPolicy, errRet error) {
	logId := tccommon.GetLogId(ctx)

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

func (me *VpcService) DescribeSecurityGroupPolicies(ctx context.Context, sgId string) (result *vpc.SecurityGroupPolicySet, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewDescribeSecurityGroupPoliciesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.SecurityGroupId = &sgId

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().DescribeSecurityGroupPolicies(request)

	if err != nil {
		errRet = err
		return
	}

	result = response.Response.SecurityGroupPolicySet

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *VpcService) DeleteSecurityGroupPolicy(ctx context.Context, ruleId string) error {
	logId := tccommon.GetLogId(ctx)

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

	if info.AddressTemplateGroupId != nil && *info.AddressTemplateGroupId != "" {
		policy.AddressTemplate = &vpc.AddressTemplateSpecification{}
		policy.AddressTemplate.AddressGroupId = info.AddressTemplateGroupId
	}

	if info.AddressTemplateId != nil && *info.AddressTemplateId != "" {
		policy.AddressTemplate = &vpc.AddressTemplateSpecification{}
		policy.AddressTemplate.AddressId = info.AddressTemplateId
	}

	if info.ProtocolTemplateGroupId != nil && *info.ProtocolTemplateGroupId != "" {
		policy.ServiceTemplate = &vpc.ServiceTemplateSpecification{}
		policy.ServiceTemplate.ServiceGroupId = info.ProtocolTemplateGroupId
	}

	if info.ProtocolTemplateId != nil && *info.ProtocolTemplateId != "" {
		policy.ServiceTemplate = &vpc.ServiceTemplateSpecification{}
		policy.ServiceTemplate.ServiceId = info.ProtocolTemplateId
	}

	if info.Description != nil && *info.Description != "" {
		policy.PolicyDescription = info.Description
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

func (me *VpcService) DeleteSecurityGroupPolicies(ctx context.Context, request *vpc.DeleteSecurityGroupPoliciesRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().DeleteSecurityGroupPolicies(request)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *VpcService) DeleteSecurityGroupPolicyByPolicyIndex(ctx context.Context, policyIndex int64, sgId, policyType string) error {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDeleteSecurityGroupPoliciesRequest()
	request.SecurityGroupId = helper.String(sgId)
	request.SecurityGroupPolicySet = new(vpc.SecurityGroupPolicySet)

	policy := new(vpc.SecurityGroupPolicy)
	policy.PolicyIndex = helper.Int64(policyIndex)
	switch strings.ToLower(policyType) {
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

func (me *VpcService) DeleteSecurityGroupPolicyByPolicyIndexList(ctx context.Context, sgId string, policyIndexList []*int64, policyType string) error {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewDeleteSecurityGroupPoliciesRequest()
	request.SecurityGroupId = helper.String(sgId)
	request.SecurityGroupPolicySet = new(vpc.SecurityGroupPolicySet)

	tmpList := make([]*vpc.SecurityGroupPolicy, 0)
	for _, v := range policyIndexList {
		policy := new(vpc.SecurityGroupPolicy)
		policy.PolicyIndex = v
		tmpList = append(tmpList, policy)
	}

	switch strings.ToLower(policyType) {

	case "ingress":
		request.SecurityGroupPolicySet.Ingress = tmpList

	case "egress":
		request.SecurityGroupPolicySet.Egress = tmpList
	}

	ratelimit.Check(request.GetAction())
	if _, err := me.client.UseVpcClient().DeleteSecurityGroupPolicies(request); err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
			logId, request.GetAction(), request.ToJsonString(), err)
		return err
	}
	return nil

}

// Deprecated: Use ModifySecurityGroupPolicies instead
func (me *VpcService) ModifySecurityGroupPolicy(ctx context.Context, ruleId string, desc *string) error {
	logId := tccommon.GetLogId(ctx)

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

func (me *VpcService) ModifySecurityGroupPolicies(ctx context.Context, request *vpc.ModifySecurityGroupPoliciesRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().ModifySecurityGroupPolicies(request)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *VpcService) DescribeSecurityGroups(ctx context.Context, sgId, sgName *string, projectId *int, tags map[string]string) (sgs []*vpc.SecurityGroup, err error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDescribeSecurityGroupsRequest()

	if sgId != nil {
		request.SecurityGroupIds = []*string{sgId}
	} else {
		if sgName != nil {
			request.Filters = append(request.Filters, &vpc.Filter{
				Name:   helper.String("security-group-name"),
				Values: []*string{sgName},
			})
		}

		if projectId != nil {
			request.Filters = append(request.Filters, &vpc.Filter{
				Name:   helper.String("project-id"),
				Values: []*string{helper.String(strconv.Itoa(*projectId))},
			})
		}

		for k, v := range tags {
			request.Filters = append(request.Filters, &vpc.Filter{
				Name:   helper.String("tag:" + k),
				Values: []*string{helper.String(v)},
			})
		}
	}

	request.Limit = helper.String(strconv.Itoa(DESCRIBE_SECURITY_GROUP_LIMIT))

	offset := 0
	count := DESCRIBE_SECURITY_GROUP_LIMIT
	// run loop at least once
	for count == DESCRIBE_SECURITY_GROUP_LIMIT {
		request.Offset = helper.String(strconv.Itoa(offset))

		if err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())

			response, err := me.client.UseVpcClient().DescribeSecurityGroups(request)
			if err != nil {
				count = 0

				if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
					if sdkError.Code == "ResourceNotFound" {
						return nil
					}
				}

				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
					logId, request.GetAction(), request.ToJsonString(), err)
				return tccommon.RetryError(err, tccommon.InternalError)
			}

			set := response.Response.SecurityGroupSet
			count = len(set)
			sgs = append(sgs, set...)

			return nil
		}); err != nil {
			log.Printf("[CRITAL]%s read security groups failed, reason: %v", logId, err)
			return nil, err
		}

		offset += count
	}

	return
}

func (me *VpcService) modifyLiteRulesInSecurityGroup(ctx context.Context, sgId string, ingress, egress []VpcSecurityGroupLiteRule) error {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewModifySecurityGroupPoliciesRequest()
	request.SecurityGroupId = &sgId
	request.SecurityGroupPolicySet = new(vpc.SecurityGroupPolicySet)
	request.SecurityGroupPolicySet.Egress = getSecurityGroupPolicies(egress)
	request.SecurityGroupPolicySet.Ingress = getSecurityGroupPolicies(ingress)

	return resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		if _, err := me.client.UseVpcClient().ModifySecurityGroupPolicies(request); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return tccommon.RetryError(err)
		}

		return nil
	})
}

func (me *VpcService) DeleteLiteRules(ctx context.Context, sgId string, rules []VpcSecurityGroupLiteRule, isIngress bool) error {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDeleteSecurityGroupPoliciesRequest()
	request.SecurityGroupId = &sgId
	request.SecurityGroupPolicySet = new(vpc.SecurityGroupPolicySet)

	if isIngress {
		request.SecurityGroupPolicySet.Ingress = getSecurityGroupPolicies(rules)
	} else {
		request.SecurityGroupPolicySet.Egress = getSecurityGroupPolicies(rules)
	}

	return resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		if _, err := me.client.UseVpcClient().DeleteSecurityGroupPolicies(request); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, request.GetAction(), request.ToJsonString(), err)

			return tccommon.RetryError(err)
		}

		return nil
	})
}

func (me *VpcService) AttachLiteRulesToSecurityGroup(ctx context.Context, sgId string, ingress, egress []VpcSecurityGroupLiteRule) error {
	logId := tccommon.GetLogId(ctx)

	if err := me.modifyLiteRulesInSecurityGroup(ctx, sgId, ingress, egress); err != nil {
		log.Printf("[CRITAL]%s attach lite rules to security group failed, reason: %v", logId, err)

		return err
	}

	return nil
}

func (me *VpcService) DescribeSecurityGroupPolices(ctx context.Context, sgId string) (ingress, egress []VpcSecurityGroupLiteRule, exist bool, err error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDescribeSecurityGroupPoliciesRequest()
	request.SecurityGroupId = &sgId

	if err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, err := me.client.UseVpcClient().DescribeSecurityGroupPolicies(request)
		if err != nil {
			if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "ResourceNotFound" {
					return nil
				}
			}

			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return tccommon.RetryError(err)
		}

		policySet := response.Response.SecurityGroupPolicySet

		for _, in := range policySet.Ingress {
			if nilFields := tccommon.CheckNil(in, map[string]string{
				"Action":          "action",
				"SecurityGroupId": "nested security group id",
			}); len(nilFields) > 0 {
				err := fmt.Errorf("api[%s] security group ingress %v are nil", request.GetAction(), nilFields)
				log.Printf("[CRITAL]%s %v", logId, err)
			}

			liteRule := VpcSecurityGroupLiteRule{
				//protocol:        strings.ToUpper(*in.Protocol),
				//port:            *in.Port,
				cidrIp:          *in.CidrBlock,
				action:          *in.Action,
				securityGroupId: *in.SecurityGroupId,
			}

			if in.Protocol != nil {
				liteRule.protocol = strings.ToUpper(*in.Protocol)
			}

			if in.Port != nil {
				liteRule.port = *in.Port
			}

			if in.AddressTemplate != nil {
				liteRule.addressId = *in.AddressTemplate.AddressId
				liteRule.addressGroupId = *in.AddressTemplate.AddressGroupId
			}

			if in.ServiceTemplate != nil {
				liteRule.protocolTemplateId = *in.ServiceTemplate.ServiceId
				liteRule.protocolTemplateGroupId = *in.ServiceTemplate.ServiceGroupId
			}

			ingress = append(ingress, liteRule)
		}

		for _, eg := range policySet.Egress {
			if nilFields := tccommon.CheckNil(eg, map[string]string{
				"Action":          "action",
				"SecurityGroupId": "nested security group id",
			}); len(nilFields) > 0 {
				err := fmt.Errorf("api[%s] security group egress %v are nil", request.GetAction(), nilFields)
				log.Printf("[CRITAL]%s %v", logId, err)
			}

			liteRule := VpcSecurityGroupLiteRule{
				action:          *eg.Action,
				cidrIp:          *eg.CidrBlock,
				securityGroupId: *eg.SecurityGroupId,
			}

			if eg.Port != nil {
				liteRule.port = *eg.Port
			}

			if eg.Protocol != nil {
				liteRule.protocol = strings.ToUpper(*eg.Protocol)
			}

			if eg.AddressTemplate != nil {
				liteRule.addressId = *eg.AddressTemplate.AddressId
				liteRule.addressGroupId = *eg.AddressTemplate.AddressGroupId
			}

			if eg.ServiceTemplate != nil {
				liteRule.protocolTemplateId = *eg.ServiceTemplate.ServiceId
				liteRule.protocolTemplateGroupId = *eg.ServiceTemplate.ServiceGroupId
			}

			egress = append(egress, liteRule)
		}

		exist = true

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s describe security group policies failed, rason: %v", logId, err)
		return nil, nil, false, err
	}

	return
}

func (me *VpcService) DetachAllLiteRulesFromSecurityGroup(ctx context.Context, sgId string) error {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewModifySecurityGroupPoliciesRequest()
	request.SecurityGroupId = &sgId
	request.SecurityGroupPolicySet = &vpc.SecurityGroupPolicySet{
		Version: helper.String("0"),
	}

	return resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		if _, err := me.client.UseVpcClient().ModifySecurityGroupPolicies(request); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return tccommon.RetryError(err)
		}

		return nil
	})
}

type securityGroupRuleBasicInfo struct {
	SgId                    string  `json:"sg_id"`
	PolicyType              string  `json:"policy_type"`
	CidrIp                  *string `json:"cidr_ip,omitempty"`
	Protocol                *string `json:"protocol"`
	PortRange               *string `json:"port_range"`
	Action                  string  `json:"action"`
	SourceSgId              *string `json:"source_sg_id"`
	Description             *string `json:"description,omitempty"`
	AddressTemplateId       *string `json:"address_template_id,omitempty"`
	AddressTemplateGroupId  *string `json:"address_template_group_id,omitempty"`
	ProtocolTemplateId      *string `json:"protocol_template_id,omitempty"`
	ProtocolTemplateGroupId *string `json:"protocol_template_group_id,omitempty"`
}

type securityGroupRuleBasicInfoWithPolicyIndex struct {
	securityGroupRuleBasicInfo
	PolicyIndex int64 `json:"policy_index"`
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

	// the newest version include template
	addressTemplateId, addressTemplateOk := m["address_template_id"]
	addressGroupTemplateId, addressTemplateGroupOk := m["address_template_group_id"]
	if addressTemplateOk || addressTemplateGroupOk {
		if addressTemplateGroupOk {
			info.AddressTemplateGroupId = common.StringPtr(addressGroupTemplateId)
		} else {
			info.AddressTemplateId = common.StringPtr(addressTemplateId)
		}
		info.CidrIp = common.StringPtr("")
		info.SourceSgId = common.StringPtr("")
	} else {
		if m["sourceSgid"] == "" {
			info.CidrIp = common.StringPtr(m["cidrIp"])
		} else {
			info.CidrIp = common.StringPtr("")
		}
		info.SourceSgId = common.StringPtr(m["sourceSgid"])
	}

	protocolTemplateId, protocolTemplateOk := m["protocol_template_id"]
	protocolGroupTemplateId, protocolTemplateGroupOk := m["protocol_template_group_id"]
	if protocolTemplateOk || protocolTemplateGroupOk {
		if protocolTemplateGroupOk {
			info.ProtocolTemplateGroupId = common.StringPtr(protocolGroupTemplateId)
		} else {
			info.ProtocolTemplateId = common.StringPtr(protocolTemplateId)
		}
		info.Protocol = common.StringPtr("")
		info.PortRange = common.StringPtr("")
	} else {
		info.Protocol = common.StringPtr(m["ipProtocol"])
		info.PortRange = common.StringPtr(m["portRange"])
	}

	info.Description = common.StringPtr(m["description"])

	log.Printf("[DEBUG] parseSecurityGroupRuleId after: %v", info)
	return
}

func comparePolicyAndSecurityGroupInfo(policy *vpc.SecurityGroupPolicy, info securityGroupRuleBasicInfo) bool {
	if policy.PolicyDescription != nil && *policy.PolicyDescription != "" {
		if info.Description == nil || *policy.PolicyDescription != *info.Description {
			return false
		}
	} else {
		if info.Description != nil && *info.Description != "" {
			return false
		}
	}
	// policy.CidrBlock will be nil if address template is set
	if policy.CidrBlock != nil && *policy.CidrBlock != "" {
		if info.CidrIp == nil || *policy.CidrBlock != *info.CidrIp {
			return false
		}
	} else {
		if info.CidrIp != nil && *info.CidrIp != "" {
			return false
		}
	}

	// policy.Port will be nil if protocol template is set
	if policy.Port != nil && *policy.Port != "" {
		if info.PortRange == nil || *policy.Port != *info.PortRange {
			return false
		}
	} else {
		if info.PortRange != nil && *info.PortRange != "" && *info.PortRange != "ALL" {
			return false
		}
	}

	// policy.Protocol will be nil if protocol template is set
	if policy.Protocol != nil && *policy.Protocol != "" {
		if info.Protocol == nil || !strings.EqualFold(*policy.Protocol, *info.Protocol) {
			return false
		}
	} else {
		if info.Protocol != nil && *info.Protocol != "" && *info.Protocol != "ALL" {
			return false
		}
	}

	// policy.SecurityGroupId always not nil
	if *policy.SecurityGroupId != *info.SourceSgId {
		return false
	}

	if !strings.EqualFold(*policy.Action, info.Action) {
		return false
	}

	// if template is not null it must be compared
	if info.ProtocolTemplateId != nil && *info.ProtocolTemplateId != "" {
		if policy.ServiceTemplate == nil || policy.ServiceTemplate.ServiceId == nil || *info.ProtocolTemplateId != *policy.ServiceTemplate.ServiceId {
			log.Printf("%s %v test", *info.ProtocolTemplateId, policy.ServiceTemplate)
			return false
		}
	} else {
		if policy.ServiceTemplate != nil && policy.ServiceTemplate.ServiceId != nil && *policy.ServiceTemplate.ServiceId != "" {
			return false
		}
	}

	if info.ProtocolTemplateGroupId != nil && *info.ProtocolTemplateGroupId != "" {
		if policy.ServiceTemplate == nil || policy.ServiceTemplate.ServiceGroupId == nil || *info.ProtocolTemplateGroupId != *policy.ServiceTemplate.ServiceGroupId {
			log.Printf("%s %v test", *info.ProtocolTemplateGroupId, policy.ServiceTemplate)
			return false
		}
	} else {
		if policy.ServiceTemplate != nil && policy.ServiceTemplate.ServiceGroupId != nil && *policy.ServiceTemplate.ServiceGroupId != "" {
			return false
		}
	}

	if info.AddressTemplateGroupId != nil && *info.AddressTemplateGroupId != "" {
		if policy.AddressTemplate == nil || policy.AddressTemplate.AddressGroupId == nil || *info.AddressTemplateGroupId != *policy.AddressTemplate.AddressGroupId {
			return false
		}
	} else {
		if policy.AddressTemplate != nil && policy.AddressTemplate.AddressGroupId != nil && *policy.AddressTemplate.AddressGroupId != "" {
			return false
		}
	}
	if info.AddressTemplateId != nil && *info.AddressTemplateId != "" {
		if policy.AddressTemplate == nil || policy.AddressTemplate.AddressId == nil || *info.AddressTemplateId != *policy.AddressTemplate.AddressId {
			return false
		}
	} else {
		if policy.AddressTemplate != nil && policy.AddressTemplate.AddressId != nil && *policy.AddressTemplate.AddressId != "" {
			return false
		}
	}

	return true
}

func parseRule(str string) (liteRule VpcSecurityGroupLiteRule, err error) {
	split := strings.Split(str, "#")
	if len(split) != 4 {
		err = fmt.Errorf("invalid security group rule %s", str)
		return
	}

	var (
		source   string
		port     string
		protocol string
		// source is "sg-xxxxxx" / "ipm-xxxxxx" / "ipmg-xxxxxx" formatted
		isInstanceIdSource = true
	)

	liteRule.action, source, port, protocol = split[0], split[1], split[2], split[3]

	if securityGroupIdRE.MatchString(source) {
		liteRule.securityGroupId = source
	} else if ipAddressIdRE.MatchString(source) {
		liteRule.addressId = source
	} else if ipAddressGroupIdRE.MatchString(source) {
		liteRule.addressGroupId = source
	} else {
		isInstanceIdSource = false
		liteRule.cidrIp = source
	}

	if v := liteRule.action; v != "ACCEPT" && v != "DROP" {
		err = fmt.Errorf("invalid action `%s`, available actions: `ACCEPT`, `DROP`", v)
		return
	}

	if net.ParseIP(liteRule.cidrIp) == nil && !isInstanceIdSource {
		if _, _, err = net.ParseCIDR(liteRule.cidrIp); err != nil {
			err = fmt.Errorf("invalid cidr_ip %s, allow cidr_ip format is `8.8.8.8` or `10.0.1.0/24`", liteRule.cidrIp)
			return
		}
	}

	liteRule.port = port
	if port != "ALL" && !portRE.MatchString(port) && !protocolTemplateRE.MatchString(protocol) {
		err = fmt.Errorf("invalid port %s, allow port format is `ALL`, `53`, `80,443` or `80-90`", liteRule.port)
		return
	}

	liteRule.protocol = protocol
	if protocolTemplateRE.MatchString(protocol) {
		liteRule.port = ""
		liteRule.protocol = ""
		if protocolTemplateIdRE.MatchString(protocol) {
			liteRule.protocolTemplateId = protocol
		} else if protocolTemplateGroupIdRE.MatchString(protocol) {
			liteRule.protocolTemplateGroupId = protocol
		}
	} else if protocol != "TCP" && protocol != "UDP" && protocol != "ALL" && protocol != "ICMP" {
		err = fmt.Errorf("invalid protocol %s, allow protocol is `ALL`, `TCP`, `UDP`, `ICMP` or `ppm(g?)-xxxxxxxx`", liteRule.protocol)
	} else if protocol == "ALL" || protocol == "ICMP" {
		if liteRule.port != "ALL" {
			err = fmt.Errorf("when protocol is %s, port must be ALL", protocol)
		} else {
			liteRule.port = ""
		}
	}

	if err != nil {
		return
	}

	return
}

/*
EIP
*/
func (me *VpcService) DescribeEipById(ctx context.Context, eipId string) (eip *vpc.Address, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewDescribeAddressesRequest()
	request.AddressIds = []*string{&eipId}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().DescribeAddresses(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.AddressSet) < 1 {
		return
	}
	eip = response.Response.AddressSet[0]
	return
}

func (me *VpcService) DescribeEipByFilter(ctx context.Context, filters map[string][]string) (eips []*vpc.Address, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewDescribeAddressesRequest()
	request.Filters = make([]*vpc.Filter, 0, len(filters))
	for k, v := range filters {
		filter := &vpc.Filter{
			Name:   helper.String(k),
			Values: []*string{},
		}
		for _, vv := range v {
			filter.Values = append(filter.Values, helper.String(vv))
		}
		request.Filters = append(request.Filters, filter)
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().DescribeAddresses(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	eips = response.Response.AddressSet
	return
}

func (me *VpcService) ModifyEipName(ctx context.Context, eipId, eipName string) error {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewModifyAddressAttributeRequest()
	request.AddressId = &eipId
	request.AddressName = &eipName

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().ModifyAddressAttribute(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *VpcService) ModifyEipBandwidthOut(ctx context.Context, eipId string, bandwidthOut int) error {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewModifyAddressesBandwidthRequest()
	request.AddressIds = []*string{&eipId}
	request.InternetMaxBandwidthOut = helper.IntInt64(bandwidthOut)

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().ModifyAddressesBandwidth(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *VpcService) ModifyEipInternetChargeType(ctx context.Context, eipId string, internetChargeType string, bandwidthOut, period, renewFlag int) error {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewModifyAddressInternetChargeTypeRequest()
	request.AddressId = &eipId
	request.InternetChargeType = &internetChargeType
	request.InternetMaxBandwidthOut = helper.IntUint64(bandwidthOut)

	if internetChargeType == "BANDWIDTH_PREPAID_BY_MONTH" {
		addressChargePrepaid := vpc.AddressChargePrepaid{}
		addressChargePrepaid.AutoRenewFlag = helper.IntInt64(renewFlag)
		addressChargePrepaid.Period = helper.IntInt64(period)
		request.AddressChargePrepaid = &addressChargePrepaid
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().ModifyAddressInternetChargeType(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *VpcService) RenewAddress(ctx context.Context, eipId string, period int, renewFlag int) error {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewRenewAddressesRequest()
	request.AddressIds = []*string{&eipId}
	addressChargePrepaid := vpc.AddressChargePrepaid{}
	addressChargePrepaid.AutoRenewFlag = helper.IntInt64(renewFlag)
	addressChargePrepaid.Period = helper.IntInt64(period)
	request.AddressChargePrepaid = &addressChargePrepaid
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().RenewAddresses(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *VpcService) DeleteEip(ctx context.Context, eipId string) error {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewReleaseAddressesRequest()
	request.AddressIds = []*string{&eipId}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().ReleaseAddresses(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *VpcService) AttachEip(ctx context.Context, eipId, instanceId string) error {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewAssociateAddressRequest()
	request.AddressId = &eipId
	request.InstanceId = &instanceId

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().AssociateAddress(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *VpcService) DescribeNatGatewayById(ctx context.Context, natGateWayId string) (natGateWay *vpc.NatGateway, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewDescribeNatGatewaysRequest()
	request.NatGatewayIds = []*string{&natGateWayId}
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().DescribeNatGateways(request)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.NatGatewaySet) > 0 {
		natGateWay = response.Response.NatGatewaySet[0]
	}

	return
}

func (me *VpcService) DescribeNatGatewayByFilter(ctx context.Context, filters map[string]string) (instances []*vpc.NatGateway, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewDescribeNatGatewaysRequest()
	)
	request.Filters = make([]*vpc.Filter, 0, len(filters))
	for k, v := range filters {
		filter := vpc.Filter{
			Name:   helper.String(k),
			Values: []*string{helper.String(v)},
		}
		request.Filters = append(request.Filters, &filter)
	}

	var offset uint64 = 0
	var pageSize uint64 = 100
	instances = make([]*vpc.NatGateway, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseVpcClient().DescribeNatGateways(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.NatGatewaySet) < 1 {
			break
		}
		instances = append(instances, response.Response.NatGatewaySet...)
		if len(response.Response.NatGatewaySet) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

func (me *VpcService) DeleteNatGateway(ctx context.Context, natGatewayId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewDeleteNatGatewayRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.NatGatewayId = &natGatewayId

	errRet = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, errRet = me.client.UseVpcClient().DeleteNatGateway(request)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		return nil
	})
	return
}

func (me *VpcService) DisassociateNatGatewayAddress(ctx context.Context, request *vpc.DisassociateNatGatewayAddressRequest) (result *vpc.DisassociateNatGatewayAddressResponse, errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	// Check if Nat Gateway Ip still associate
	gateway, err := me.DescribeNatGatewayById(ctx, *request.NatGatewayId)

	if err != nil {
		errRet = err
		return
	}

	if gateway == nil || len(gateway.PublicIpAddressSet) == 0 {
		return
	}

	var gatewayAddresses []string
	var candidates []*string

	for i := range gateway.PublicIpAddressSet {
		addr := gateway.PublicIpAddressSet[i].PublicIpAddress
		gatewayAddresses = append(gatewayAddresses, *addr)
	}

	for i := range request.PublicIpAddresses {
		addr := request.PublicIpAddresses[i]
		if helper.StringsContain(gatewayAddresses, *addr) {
			candidates = append(candidates, addr)
		}
	}

	if len(candidates) == 0 {
		return nil, nil
	}

	request.PublicIpAddresses = candidates

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().DisassociateNatGatewayAddress(request)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	result = response
	return
}

func (me *VpcService) UnattachEip(ctx context.Context, eipId string) error {
	eipUnattachLocker.Lock()
	defer eipUnattachLocker.Unlock()

	logId := tccommon.GetLogId(ctx)
	eip, err := me.DescribeEipById(ctx, eipId)
	if err != nil {
		return err
	}
	if eip == nil || *eip.AddressStatus == EIP_STATUS_UNBIND {
		return nil
	}

	// DisassociateAddress Doesn't support Disassociate NAT Address
	if eip.InstanceId != nil && strings.HasPrefix(*eip.InstanceId, "nat-") {
		request := vpc.NewDisassociateNatGatewayAddressRequest()
		request.NatGatewayId = eip.InstanceId
		request.PublicIpAddresses = []*string{eip.AddressIp}
		_, err := me.DisassociateNatGatewayAddress(ctx, request)
		if err != nil {
			return err
		}

		outErr := resource.Retry(tccommon.ReadRetryTimeout*3, func() *resource.RetryError {
			eip, err := me.DescribeEipById(ctx, eipId)
			if err != nil {
				return tccommon.RetryError(err)
			}
			if eip != nil && *eip.AddressStatus != EIP_STATUS_UNBIND {
				return resource.RetryableError(fmt.Errorf("eip is still %s", EIP_STATUS_UNBIND))
			}
			return nil
		})

		if outErr != nil {
			return outErr
		}
	}

	request := vpc.NewDisassociateAddressRequest()
	request.AddressId = &eipId
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().DisassociateAddress(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	if response.Response.TaskId == nil {
		return nil
	}
	taskId, err := strconv.ParseUint(*response.Response.TaskId, 10, 64)
	if err != nil {
		return nil
	}

	taskRequest := vpc.NewDescribeTaskResultRequest()
	taskRequest.TaskId = &taskId
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(taskRequest.GetAction())
		taskResponse, err := me.client.UseVpcClient().DescribeTaskResult(taskRequest)
		if err != nil {
			return tccommon.RetryError(err)
		}
		if taskResponse.Response.Result != nil && *taskResponse.Response.Result == EIP_TASK_STATUS_RUNNING {
			return resource.RetryableError(errors.New("eip task is running"))
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (me *VpcService) CreateEni(
	ctx context.Context,
	name, vpcId, subnetId, desc string,
	securityGroups []string,
	ipv4Count *int,
	ipv4s []VpcEniIP,
	tags map[string]string,
) (id string, err error) {
	logId := tccommon.GetLogId(ctx)
	client := me.client.UseVpcClient()

	createRequest := vpc.NewCreateNetworkInterfaceRequest()
	createRequest.NetworkInterfaceName = &name
	createRequest.VpcId = &vpcId
	createRequest.SubnetId = &subnetId
	createRequest.NetworkInterfaceDescription = &desc

	if len(securityGroups) > 0 {
		createRequest.SecurityGroupIds = common.StringPtrs(securityGroups)
	}

	if len(tags) > 0 {
		for tagKey, tagValue := range tags {
			tag := vpc.Tag{
				Key:   helper.String(tagKey),
				Value: helper.String(tagValue),
			}
			createRequest.Tags = append(createRequest.Tags, &tag)
		}
	}

	if ipv4Count != nil {
		// create will assign a primary ip, secondary ip count is *ipv4Count-1
		createRequest.SecondaryPrivateIpAddressCount = helper.IntUint64(*ipv4Count - 1)
	}

	var wantIpv4 []string

	for _, ipv4 := range ipv4s {
		wantIpv4 = append(wantIpv4, ipv4.ip.String())
		createRequest.PrivateIpAddresses = append(createRequest.PrivateIpAddresses, &vpc.PrivateIpAddressSpecification{
			PrivateIpAddress: helper.String(ipv4.ip.String()),
			Primary:          helper.Bool(ipv4.primary),
			Description:      ipv4.desc,
		})
	}

	if err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(createRequest.GetAction())

		response, err := client.CreateNetworkInterface(createRequest)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, createRequest.GetAction(), createRequest.ToJsonString(), err)
			return tccommon.RetryError(err)
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

		ipv4Set := eni.PrivateIpAddressSet

		if len(wantIpv4) > 0 {
			checkMap := make(map[string]bool, len(wantIpv4))
			for _, ipv4 := range wantIpv4 {
				checkMap[ipv4] = false
			}

			for _, ipv4 := range ipv4Set {
				if ipv4.PrivateIpAddress == nil {
					err := fmt.Errorf("api[%s] eni ipv4 ip is nil", createRequest.GetAction())
					log.Printf("[CRITAL]%s %v", logId, err)
					return resource.NonRetryableError(err)
				}

				checkMap[*ipv4.PrivateIpAddress] = true
			}

			for ipv4, checked := range checkMap {
				if !checked {
					err := fmt.Errorf("api[%s] doesn't assign %s ip", createRequest.GetAction(), ipv4)
					log.Printf("[CRITAL]%s %v", logId, err)
					return resource.NonRetryableError(err)
				}
			}
		} else {
			if len(ipv4Set) != *ipv4Count {
				err := fmt.Errorf("api[%s] doesn't assign enough ip", createRequest.GetAction())
				log.Printf("[CRITAL]%s %v", logId, err)
				return resource.NonRetryableError(err)
			}

			wantIpv4 = make([]string, 0, *ipv4Count)
			for _, ipv4 := range ipv4Set {
				if ipv4.PrivateIpAddress == nil {
					err := fmt.Errorf("api[%s] eni ipv4 ip is nil", createRequest.GetAction())
					log.Printf("[CRITAL]%s %v", logId, err)
					return resource.NonRetryableError(err)
				}

				wantIpv4 = append(wantIpv4, *ipv4.PrivateIpAddress)
			}
		}

		id = *eni.NetworkInterfaceId

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s create eni failed, reason: %v", logId, err)
		return "", err
	}

	if err := waitEniReady(ctx, id, client, wantIpv4, nil); err != nil {
		log.Printf("[CRITAL]%s create eni failed, reason: %v", logId, err)
		return "", err
	}

	return
}

func (me *VpcService) describeEnis(
	ctx context.Context,
	ids []string,
	vpcId, subnetId, id, cvmId, sgId, name, desc, ipv4 *string,
	tags map[string]string,
) (enis []*vpc.NetworkInterface, err error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDescribeNetworkInterfacesRequest()

	if len(ids) > 0 {
		request.NetworkInterfaceIds = common.StringPtrs(ids)
	}

	if vpcId != nil {
		request.Filters = append(request.Filters, &vpc.Filter{
			Name:   helper.String("vpc-id"),
			Values: []*string{vpcId},
		})
	}

	if subnetId != nil {
		request.Filters = append(request.Filters, &vpc.Filter{
			Name:   helper.String("subnet-id"),
			Values: []*string{subnetId},
		})
	}

	if id != nil {
		request.Filters = append(request.Filters, &vpc.Filter{
			Name:   helper.String("network-interface-id"),
			Values: []*string{id},
		})
	}

	if cvmId != nil {
		request.Filters = append(request.Filters, &vpc.Filter{
			Name:   helper.String("attachment.instance-id"),
			Values: []*string{cvmId},
		})
	}

	if sgId != nil {
		request.Filters = append(request.Filters, &vpc.Filter{
			Name:   helper.String("groups.security-group-id"),
			Values: []*string{sgId},
		})
	}

	if name != nil {
		request.Filters = append(request.Filters, &vpc.Filter{
			Name:   helper.String("network-interface-name"),
			Values: []*string{name},
		})
	}

	if desc != nil {
		request.Filters = append(request.Filters, &vpc.Filter{
			Name:   helper.String("network-interface-description"),
			Values: []*string{desc},
		})
	}

	if ipv4 != nil {
		request.Filters = append(request.Filters, &vpc.Filter{
			Name:   helper.String("address-ip"),
			Values: []*string{ipv4},
		})
	}

	for k, v := range tags {
		request.Filters = append(request.Filters, &vpc.Filter{
			Name:   helper.String("tag:" + k),
			Values: []*string{helper.String(v)},
		})
	}

	var offset uint64
	request.Offset = &offset
	request.Limit = helper.IntUint64(ENI_DESCRIBE_LIMIT)

	count := ENI_DESCRIBE_LIMIT
	for count == ENI_DESCRIBE_LIMIT {
		if err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
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
				return tccommon.RetryError(err)
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

func (me *VpcService) DescribeEniById(ctx context.Context, ids []string) (enis []*vpc.NetworkInterface, err error) {
	return me.describeEnis(ctx, ids, nil, nil, nil, nil, nil, nil, nil, nil, nil)
}

func (me *VpcService) ModifyEniAttribute(ctx context.Context, id string, name, desc *string, sgs []string) error {
	logId := tccommon.GetLogId(ctx)
	client := me.client.UseVpcClient()

	request := vpc.NewModifyNetworkInterfaceAttributeRequest()
	request.NetworkInterfaceId = &id
	request.NetworkInterfaceName = name
	request.NetworkInterfaceDescription = desc
	request.SecurityGroupIds = common.StringPtrs(sgs)

	if err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		if _, err := client.ModifyNetworkInterfaceAttribute(request); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return tccommon.RetryError(err)
		}

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s modify eni attribute failed, reason: %v", logId, err)
		return err
	}

	if err := waitEniReady(ctx, id, client, nil, nil); err != nil {
		log.Printf("[CRITAL]%s modify eni attribute failed, reason: %v", logId, err)
		return err
	}

	return nil
}

func (me *VpcService) UnAssignIpv4FromEni(ctx context.Context, id string, ipv4s []string) error {
	logId := tccommon.GetLogId(ctx)
	client := me.client.UseVpcClient()

	request := vpc.NewUnassignPrivateIpAddressesRequest()
	request.NetworkInterfaceId = &id
	request.PrivateIpAddresses = make([]*vpc.PrivateIpAddressSpecification, 0, len(ipv4s))
	for _, ipv4 := range ipv4s {
		request.PrivateIpAddresses = append(request.PrivateIpAddresses, &vpc.PrivateIpAddressSpecification{
			PrivateIpAddress: helper.String(ipv4),
		})
	}

	if err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		if _, err := client.UnassignPrivateIpAddresses(request); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return tccommon.RetryError(err)
		}

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s unassign ipv4 from eni failed, reason: %v", logId, err)
		return err
	}

	if err := waitEniReady(ctx, id, client, nil, ipv4s); err != nil {
		log.Printf("[CRITAL]%s unassign ipv4 from eni failed, reason: %v", logId, err)
		return err
	}

	return nil
}

func (me *VpcService) AssignIpv4ToEni(ctx context.Context, id string, ipv4s []VpcEniIP, ipv4Count *int) error {
	logId := tccommon.GetLogId(ctx)
	client := me.client.UseVpcClient()

	request := vpc.NewAssignPrivateIpAddressesRequest()
	request.NetworkInterfaceId = &id

	if ipv4Count != nil {
		request.SecondaryPrivateIpAddressCount = helper.IntUint64(*ipv4Count)
	}

	var wantIpv4 []string

	if len(ipv4s) > 0 {
		request.PrivateIpAddresses = make([]*vpc.PrivateIpAddressSpecification, 0, len(ipv4s))
		wantIpv4 = make([]string, 0, len(ipv4s))

		for _, ipv4 := range ipv4s {
			wantIpv4 = append(wantIpv4, ipv4.ip.String())
			request.PrivateIpAddresses = append(request.PrivateIpAddresses, &vpc.PrivateIpAddressSpecification{
				PrivateIpAddress: helper.String(ipv4.ip.String()),
				Primary:          helper.Bool(ipv4.primary),
				Description:      ipv4.desc,
			})
		}
	}

	if err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, err := client.AssignPrivateIpAddresses(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return tccommon.RetryError(err)
		}

		ipv4Set := response.Response.PrivateIpAddressSet

		if len(wantIpv4) > 0 {
			checkMap := make(map[string]bool, len(wantIpv4))
			for _, ipv4 := range wantIpv4 {
				checkMap[ipv4] = false
			}

			for _, ipv4 := range ipv4Set {
				if ipv4.PrivateIpAddress == nil {
					err := fmt.Errorf("api[%s] eni ipv4 ip is nil", request.GetAction())
					log.Printf("[CRITAL]%s %v", logId, err)
					return resource.NonRetryableError(err)
				}

				checkMap[*ipv4.PrivateIpAddress] = true
			}

			for ipv4, checked := range checkMap {
				if !checked {
					err := fmt.Errorf("api[%s] doesn't assign %s ip", request.GetAction(), ipv4)
					log.Printf("[CRITAL]%s %v", logId, err)
					return resource.NonRetryableError(err)
				}
			}
		} else {
			if len(ipv4Set) != *ipv4Count {
				err := fmt.Errorf("api[%s] doesn't assign enough ip", request.GetAction())
				log.Printf("[CRITAL]%s %v", logId, err)
				return resource.NonRetryableError(err)
			}

			wantIpv4 = make([]string, 0, *ipv4Count)
			for _, ipv4 := range ipv4Set {
				if ipv4.PrivateIpAddress == nil {
					err := fmt.Errorf("api[%s] eni ipv4 ip is nil", request.GetAction())
					log.Printf("[CRITAL]%s %v", logId, err)
					return resource.NonRetryableError(err)
				}

				wantIpv4 = append(wantIpv4, *ipv4.PrivateIpAddress)
			}
		}

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s assign ipv4 to eni failed, reason: %v", logId, err)
		return err
	}

	if err := waitEniReady(ctx, id, client, wantIpv4, nil); err != nil {
		log.Printf("[CRITAL]%s assign ipv4 to eni failed, reason: %v", logId, err)
		return err
	}

	return nil
}

func (me *VpcService) DeleteEni(ctx context.Context, id string) error {
	logId := tccommon.GetLogId(ctx)
	client := me.client.UseVpcClient()

	deleteRequest := vpc.NewDeleteNetworkInterfaceRequest()
	deleteRequest.NetworkInterfaceId = &id

	if err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(deleteRequest.GetAction())

		if _, err := client.DeleteNetworkInterface(deleteRequest); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, deleteRequest.GetAction(), deleteRequest.ToJsonString(), err)
			return tccommon.RetryError(err)
		}

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s delete eni failed, reason: %v", logId, err)
		return err
	}

	describeRequest := vpc.NewDescribeNetworkInterfacesRequest()
	describeRequest.NetworkInterfaceIds = []*string{&id}

	if err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
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
			return tccommon.RetryError(err)
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
	logId := tccommon.GetLogId(ctx)
	client := me.client.UseVpcClient()

	attachRequest := vpc.NewAttachNetworkInterfaceRequest()
	attachRequest.NetworkInterfaceId = &eniId
	attachRequest.InstanceId = &cvmId

	if err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(attachRequest.GetAction())

		if _, err := client.AttachNetworkInterface(attachRequest); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, attachRequest.GetAction(), attachRequest.ToJsonString(), err)
			return tccommon.RetryError(err)
		}

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s attach eni to instance failed, reason: %v", logId, err)
		return err
	}

	describeRequest := vpc.NewDescribeNetworkInterfacesRequest()
	describeRequest.NetworkInterfaceIds = []*string{&eniId}

	if err := resource.Retry(2*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(describeRequest.GetAction())

		response, err := client.DescribeNetworkInterfaces(describeRequest)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, describeRequest.GetAction(), describeRequest.ToJsonString(), err)
			return tccommon.RetryError(err)
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
	logId := tccommon.GetLogId(ctx)
	client := me.client.UseVpcClient()

	request := vpc.NewDetachNetworkInterfaceRequest()
	request.NetworkInterfaceId = &eniId
	request.InstanceId = &cvmId

	if err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		if _, err := client.DetachNetworkInterface(request); err != nil {
			if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				switch sdkError.Code {
				case "UnsupportedOperation.InvalidState":
					return resource.RetryableError(errors.New("cvm may still bind eni"))

				case "ResourceNotFound":
					// eni or cvm doesn't exist
					return nil
				}
			}

			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return tccommon.RetryError(err)
		}

		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s detach eni from instance failed, reason: %v", logId, err)
		return err
	}

	if err := waitEniDetach(ctx, eniId, client); err != nil {
		log.Printf("[CRITAL]%s detach eni from instance failed, reason: %v", logId, err)
		return err
	}

	return nil
}

func (me *VpcService) ModifyEniPrimaryIpv4Desc(ctx context.Context, id, ip string, desc *string) error {
	logId := tccommon.GetLogId(ctx)
	client := me.client.UseVpcClient()

	request := vpc.NewModifyPrivateIpAddressesAttributeRequest()
	request.NetworkInterfaceId = &id
	request.PrivateIpAddresses = []*vpc.PrivateIpAddressSpecification{
		{
			PrivateIpAddress: &ip,
			Description:      desc,
		},
	}

	if err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		if _, err := client.ModifyPrivateIpAddressesAttribute(request); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return tccommon.RetryError(err)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s modify eni primary ipv4 description failed, reason: %v", logId, err)
		return err
	}

	if err := waitEniReady(ctx, id, client, []string{ip}, nil); err != nil {
		log.Printf("[CRITAL]%s modify eni primary ipv4 description failed, reason: %v", logId, err)
		return err
	}

	return nil
}

func (me *VpcService) DescribeEniByFilters(
	ctx context.Context,
	vpcId, subnetId, cvmId, sgId, name, desc, ipv4 *string,
	tags map[string]string,
) (enis []*vpc.NetworkInterface, err error) {
	return me.describeEnis(ctx, nil, vpcId, subnetId, nil, cvmId, sgId, name, desc, ipv4, tags)
}

func (me *VpcService) DescribeHaVipByFilter(ctx context.Context, filters map[string]string) (instances []*vpc.HaVip, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewDescribeHaVipsRequest()
	)
	request.Filters = make([]*vpc.Filter, 0, len(filters))
	for k, v := range filters {
		filter := vpc.Filter{
			Name:   helper.String(k),
			Values: []*string{helper.String(v)},
		}
		request.Filters = append(request.Filters, &filter)
	}

	var offset uint64 = 0
	var pageSize uint64 = 100
	instances = make([]*vpc.HaVip, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseVpcClient().DescribeHaVips(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.HaVipSet) < 1 {
			break
		}
		instances = append(instances, response.Response.HaVipSet...)
		if len(response.Response.HaVipSet) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

func (me *VpcService) DescribeHaVipEipById(ctx context.Context, haVipEipAttachmentId string) (eip string, haVip string, has bool, errRet error) {
	logId := tccommon.GetLogId(ctx)
	client := me.client.UseVpcClient()

	items := strings.Split(haVipEipAttachmentId, "#")
	if len(items) != 2 {
		errRet = fmt.Errorf("decode HA VIP EIP attachment ID error %s", haVipEipAttachmentId)
		return
	}
	haVipId := items[0]
	addressIp := items[1]

	request := vpc.NewDescribeHaVipsRequest()
	request.HaVipIds = []*string{&haVipId}
	eip = ""
	haVip = ""
	if err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		if result, err := client.DescribeHaVips(request); err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return tccommon.RetryError(err)
		} else {
			length := len(result.Response.HaVipSet)
			if length != 1 {
				if length == 0 {
					return nil
				} else {
					err = fmt.Errorf("query havip %s eip %s failed, the SDK returns %d HaVips", haVipId, addressIp, length)
					return resource.NonRetryableError(err)
				}
			} else {
				eip = *result.Response.HaVipSet[0].AddressIp
				if addressIp != eip {
					return nil
				}
				has = true
				haVip = haVipId
			}
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s describe HA VIP attachment failed, reason: %v", logId, err)
		errRet = err
	}
	return eip, haVip, has, errRet
}

func (me *VpcService) DeleteHaVip(ctx context.Context, haVipId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewDeleteHaVipRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.HaVipId = &haVipId

	errRet = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, errRet = me.client.UseVpcClient().DeleteHaVip(request)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		return nil
	})
	return
}

func waitEniReady(ctx context.Context, id string, client *vpc.Client, wantIpv4s []string, dropIpv4s []string) error {
	logId := tccommon.GetLogId(ctx)

	wantCheckMap := make(map[string]bool, len(wantIpv4s))
	for _, ipv4 := range wantIpv4s {
		wantCheckMap[ipv4] = false
	}

	dropCheckMap := make(map[string]struct{}, len(dropIpv4s))
	for _, ipv4 := range dropIpv4s {
		dropCheckMap[ipv4] = struct{}{}
	}

	request := vpc.NewDescribeNetworkInterfacesRequest()
	request.NetworkInterfaceIds = []*string{helper.String(id)}

	if err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, err := client.DescribeNetworkInterfaces(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return tccommon.RetryError(err)
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
			if ipv4.PrivateIpAddress == nil {
				err := fmt.Errorf("api[%s] eni ipv4 ip is nil", request.GetAction())
				log.Printf("[CRITAL]%s %v", logId, err)
				return resource.NonRetryableError(err)
			}

			// check drop
			if _, ok := dropCheckMap[*ipv4.PrivateIpAddress]; ok {
				err := fmt.Errorf("api[%s] drop ip %s still exists", request.GetAction(), *ipv4.PrivateIpAddress)
				log.Printf("[DEBUG]%s %v", logId, err)
				return resource.RetryableError(err)
			}

			// check want
			if _, ok := wantCheckMap[*ipv4.PrivateIpAddress]; ok {
				wantCheckMap[*ipv4.PrivateIpAddress] = true
			}

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

		for ipv4, checked := range wantCheckMap {
			if !checked {
				err := fmt.Errorf("api[%s] ipv4 %s is no ready", request.GetAction(), ipv4)
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

func flattenVpnSPDList(spd []*vpc.SecurityPolicyDatabase) (mapping []*map[string]interface{}) {
	mapping = make([]*map[string]interface{}, 0, len(spd))
	for _, spg := range spd {
		item := make(map[string]interface{})
		item["local_cidr_block"] = spg.LocalCidrBlock
		item["remote_cidr_block"] = spg.RemoteCidrBlock
		mapping = append(mapping, &item)
	}
	return
}

func waitEniDetach(ctx context.Context, id string, client *vpc.Client) error {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDescribeNetworkInterfacesRequest()
	request.NetworkInterfaceIds = []*string{helper.String(id)}

	return resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, err := client.DescribeNetworkInterfaces(request)
		if err != nil {
			if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok && sdkError.Code == "ResourceNotFound" {
				return nil
			}

			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, request.GetAction(), request.ToJsonString(), err)
			return tccommon.RetryError(err)
		}

		enis := response.Response.NetworkInterfaceSet

		if len(enis) == 0 {
			return nil
		}

		eni := enis[0]

		if eni.Attachment == nil {
			return nil
		}

		if eni.Attachment.InstanceId != nil && *eni.Attachment.InstanceId != "" {
			return resource.RetryableError(fmt.Errorf("eni %s still bind in cvm %s", id, *eni.Attachment.InstanceId))
		}

		if eni.State == nil {
			return resource.NonRetryableError(fmt.Errorf("eni %s state is nil", id))
		}

		if *eni.State != ENI_STATE_AVAILABLE {
			return resource.RetryableError(errors.New("eni is not available"))
		}

		return nil
	})
}

// deal acl
func parseACLRule(str string) (liteRule VpcACLRule, err error) {
	split := strings.Split(str, "#")
	if len(split) != 4 {
		err = fmt.Errorf("invalid acl rule %s", str)
		return
	}

	liteRule.action, liteRule.cidrIp, liteRule.port, liteRule.protocol = split[0], split[1], split[2], split[3]

	switch liteRule.action {
	default:
		err = fmt.Errorf("invalid action %s, allow action is `ACCEPT` or `DROP`", liteRule.action)
		return
	case "ACCEPT", "DROP":
	}

	if net.ParseIP(liteRule.cidrIp) == nil {
		if _, _, err = net.ParseCIDR(liteRule.cidrIp); err != nil {
			err = fmt.Errorf("invalid cidr_ip %s, allow cidr_ip format is `8.8.8.8` or `10.0.1.0/24`", liteRule.cidrIp)
			return
		}
	}

	if liteRule.port != "ALL" && !regexp.MustCompile(`^(\d{1,5},)*\d{1,5}$|^\d{1,5}-\d{1,5}$`).MatchString(liteRule.port) {
		err = fmt.Errorf("invalid port %s, allow port format is `ALL`, `53`, `80,443` or `80-90`", liteRule.port)
		return
	}

	switch liteRule.protocol {
	default:
		err = fmt.Errorf("invalid protocol %s, allow protocol is `ALL`, `TCP`, `UDP` or `ICMP`", liteRule.protocol)
		return

	case "ALL", "ICMP":
		if liteRule.port != "ALL" {
			err = fmt.Errorf("when protocol is %s, port must be ALL", liteRule.protocol)
			return
		}

		// when protocol is ALL or ICMP, port should be "" to avoid sdk error
		liteRule.port = ""

	case "TCP", "UDP":
	}

	return
}

func (me *VpcService) CreateVpcNetworkAcl(ctx context.Context, vpcID string, name string, tags map[string]string) (aclID string, errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = vpc.NewCreateNetworkAclRequest()
		response *vpc.CreateNetworkAclResponse
		err      error
	)

	request.VpcId = &vpcID
	request.NetworkAclName = &name

	if len(tags) > 0 {
		for tagKey, tagValue := range tags {
			tag := vpc.Tag{
				Key:   helper.String(tagKey),
				Value: helper.String(tagValue),
			}
			request.Tags = append(request.Tags, &tag)
		}
	}

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err = me.client.UseVpcClient().CreateNetworkAcl(request)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
			logId, request.GetAction(), request.ToJsonString(), err)
		errRet = err
		return
	}

	aclID = *response.Response.NetworkAcl.NetworkAclId
	return
}

func (me *VpcService) AttachRulesToACL(ctx context.Context, aclID string, ingressParm, egressParm []VpcACLRule) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	if len(ingressParm) == 0 && len(egressParm) == 0 {
		return
	}
	if errRet = me.ModifyNetWorkAclRules(ctx, aclID, ingressParm, egressParm); errRet != nil {
		log.Printf("[CRITAL]%s attach rules to acl failed, reason: %v", logId, errRet)
	}
	return
}

func (me *VpcService) ModifyNetWorkAclRules(ctx context.Context, aclID string, ingressParm, egressParm []VpcACLRule) (errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewModifyNetworkAclEntriesRequest()
		err     error
		ingress []*vpc.NetworkAclEntry
		egress  []*vpc.NetworkAclEntry
	)

	for i := range ingressParm {
		policy := &vpc.NetworkAclEntry{
			Protocol:  &ingressParm[i].protocol,
			CidrBlock: &ingressParm[i].cidrIp,
			Action:    &ingressParm[i].action,
		}

		if ingressParm[i].port != "" {
			policy.Port = &ingressParm[i].port
		}

		ingress = append(ingress, policy)
	}

	for i := range egressParm {
		policy := &vpc.NetworkAclEntry{
			Protocol:  &egressParm[i].protocol,
			CidrBlock: &egressParm[i].cidrIp,
			Action:    &egressParm[i].action,
		}

		if egressParm[i].port != "" {
			policy.Port = &egressParm[i].port
		}

		egress = append(egress, policy)
	}

	request.NetworkAclId = &aclID
	request.NetworkAclEntrySet = &vpc.NetworkAclEntrySet{
		Ingress: ingress,
		Egress:  egress,
	}

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err = me.client.UseVpcClient().ModifyNetworkAclEntries(request)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
			logId, request.GetAction(), request.ToJsonString(), err)
		errRet = err
		return
	}

	return
}

func (me *VpcService) DescribeNetWorkByACLID(ctx context.Context, aclID string) (info *vpc.NetworkAcl, has int, errRet error) {
	results, err := me.DescribeNetWorkAcls(ctx, aclID, "", "")
	if err != nil {
		errRet = err
		return
	}

	has = len(results)
	if has == 0 {
		return
	}

	info = results[0]
	return
}

func (me *VpcService) DeleteAcl(ctx context.Context, aclID string) (errRet error) {
	var (
		logId       = tccommon.GetLogId(ctx)
		err         error
		networkAcls []*vpc.NetworkAcl
		request     = vpc.NewDeleteNetworkAclRequest()
	)

	// Disassociate Network Acl Subnets
	networkAcls, err = me.DescribeNetWorkAcls(ctx, aclID, "", "")
	if err != nil {
		errRet = err
		return
	}

	if len(networkAcls) > 0 {
		subnets := networkAcls[0].SubnetSet
		if len(subnets) > 0 {
			requestSubnet := vpc.NewDisassociateNetworkAclSubnetsRequest()
			requestSubnet.NetworkAclId = &aclID

			for i := range subnets {
				requestSubnet.SubnetIds = append(requestSubnet.SubnetIds, subnets[i].SubnetId)
			}

			err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				ratelimit.Check(request.GetAction())
				_, err = me.client.UseVpcClient().DisassociateNetworkAclSubnets(requestSubnet)
				if err != nil {
					if ee, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
						if ee.Code == VPCNotFound {
							log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
								logId, request.GetAction(), request.ToJsonString(), err)
							return nil
						}
					}
					return tccommon.RetryError(err, tccommon.InternalError)
				}
				return nil
			})
			if err != nil {
				errRet = err
				return
			}
		}
	}

	// delete acl
	request.NetworkAclId = &aclID
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err = me.client.UseVpcClient().DeleteNetworkAcl(request)

		if err != nil {
			ee, ok := err.(*sdkErrors.TencentCloudSDKError)
			if !ok {
				return tccommon.RetryError(err, tccommon.InternalError)
			}
			if ee.Code == VPCNotFound {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
					logId, request.GetAction(), request.ToJsonString(), err)
				return nil
			}
			return tccommon.RetryError(err, tccommon.InternalError)
		}

		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
			logId, request.GetAction(), request.ToJsonString(), err)
		errRet = err
		return
	}

	return
}

func (me *VpcService) ModifyVpcNetworkAcl(ctx context.Context, id *string, name *string) (errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		err     error
		request = vpc.NewModifyNetworkAclAttributeRequest()
	)

	request.NetworkAclId = id
	request.NetworkAclName = name

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err = me.client.UseVpcClient().ModifyNetworkAclAttribute(request)
		if err != nil {
			ee, ok := err.(*sdkErrors.TencentCloudSDKError)
			if !ok {
				return tccommon.RetryError(err, tccommon.InternalError)
			}
			if ee.Code == VPCNotFound {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
					logId, request.GetAction(), request.ToJsonString(), err)
				return resource.NonRetryableError(err)
			}
			return tccommon.RetryError(err, tccommon.InternalError)
		}

		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
			logId, request.GetAction(), request.ToJsonString(), err)
		errRet = err
		return
	}

	return
}

func (me *VpcService) AssociateAclSubnets(ctx context.Context, aclId string, subnetIds []string) (errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewAssociateNetworkAclSubnetsRequest()
		err     error
		subIds  []*string
	)

	for _, i := range subnetIds {
		subIds = append(subIds, &i)
	}

	request.NetworkAclId = &aclId
	request.SubnetIds = subIds

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err = me.client.UseVpcClient().AssociateNetworkAclSubnets(request)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
			logId, request.GetAction(), request.ToJsonString(), err)
		errRet = err
		return
	}
	return
}

func (me *VpcService) DescribeNetWorkAcls(ctx context.Context, aclID, vpcID, name string) (info []*vpc.NetworkAcl, errRet error) {
	var (
		logId            = tccommon.GetLogId(ctx)
		request          = vpc.NewDescribeNetworkAclsRequest()
		response         *vpc.DescribeNetworkAclsResponse
		err              error
		filters          []*vpc.Filter
		offset, pageSize uint64 = 0, 100
	)

	if vpcID != "" {
		filters = me.fillFilter(filters, "vpc-id", vpcID)
	}
	if aclID != "" {
		filters = me.fillFilter(filters, "network-acl-id", aclID)
	}
	if name != "" {
		filters = me.fillFilter(filters, "network-acl-name", name)
	}

	if len(filters) > 0 {
		request.Filters = filters
	}

	request.Offset = &offset
	request.Limit = &pageSize
	for {
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			response, err = me.client.UseVpcClient().DescribeNetworkAcls(request)
			if err != nil {
				if ee, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
					if ee.Code == VPCNotFound {
						return nil
					}
				}
				return tccommon.RetryError(err, tccommon.InternalError)
			}
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, request.GetAction(), request.ToJsonString(), err)
			errRet = err
			return
		}
		if response.Response == nil {
			return
		}

		info = append(info, response.Response.NetworkAclSet...)
		if len(response.Response.NetworkAclSet) < int(pageSize) {
			break
		}

		offset += pageSize
	}

	return
}

func (me *VpcService) DescribeByAclId(ctx context.Context, attachmentAcl string) (has bool, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewDisassociateNetworkAclSubnetsRequest()
		aclId   string
	)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, request.GetAction(), request.ToJsonString(), errRet)
		}
	}()

	if attachmentAcl == "" {
		errRet = fmt.Errorf("DisassociateNetworkAclSubnets  can not invoke by empty routeTableId.")
		return
	}

	aclId = strings.Split(attachmentAcl, "#")[0]

	results, err := me.DescribeNetWorkAcls(ctx, aclId, "", "")
	if err != nil {
		errRet = err
		return
	}
	if len(results) < 1 || len(results[0].SubnetSet) < 1 {
		return
	}

	has = true
	return
}

func (me *VpcService) DeleteAclAttachment(ctx context.Context, attachmentAcl string) (errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewDisassociateNetworkAclSubnetsRequest()
		err     error
	)

	if attachmentAcl == "" {
		errRet = fmt.Errorf("DeleteRouteTable can not invoke by empty NetworkAclId.")
		return
	}

	items := strings.Split(attachmentAcl, "#")
	request.NetworkAclId = &items[0]
	request.SubnetIds = helper.Strings(items[1:])

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err = me.client.UseVpcClient().DisassociateNetworkAclSubnets(request)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
			logId, request.GetAction(), request.ToJsonString(), err)
		errRet = err
		return
	}
	return
}

func (me *VpcService) DescribeVpngwById(ctx context.Context, vpngwId string) (has bool, gateway *vpc.VpnGateway, err error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = vpc.NewDescribeVpnGatewaysRequest()
		response *vpc.DescribeVpnGatewaysResponse
	)
	request.VpnGatewayIds = []*string{&vpngwId}
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		response, err = me.client.UseVpcClient().DescribeVpnGateways(request)
		if err != nil {
			ee, ok := err.(*sdkErrors.TencentCloudSDKError)
			if !ok {
				return tccommon.RetryError(err)
			}
			if ee.Code == VPCNotFound {
				return nil
			} else {
				return tccommon.RetryError(err)
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]", logId, request.GetAction(), request.ToJsonString(), err)
		return
	}
	if response == nil || response.Response == nil || len(response.Response.VpnGatewaySet) < 1 {
		has = false
		return
	}

	gateway = response.Response.VpnGatewaySet[0]
	has = true
	return
}

func (me *VpcService) DescribeVpnGwByFilter(ctx context.Context, filters map[string]string) (instances []*vpc.VpnGateway, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewDescribeVpnGatewaysRequest()
	)
	request.Filters = make([]*vpc.FilterObject, 0, len(filters))
	for k, v := range filters {
		filter := vpc.FilterObject{
			Name:   helper.String(k),
			Values: []*string{helper.String(v)},
		}
		request.Filters = append(request.Filters, &filter)
	}

	var offset uint64 = 0
	var pageSize uint64 = 100
	instances = make([]*vpc.VpnGateway, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseVpcClient().DescribeVpnGateways(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.VpnGatewaySet) < 1 {
			break
		}
		instances = append(instances, response.Response.VpnGatewaySet...)
		if len(response.Response.VpnGatewaySet) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

func (me *VpcService) DeleteVpnGateway(ctx context.Context, vpnGatewayId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewDeleteVpnGatewayRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.VpnGatewayId = &vpnGatewayId

	errRet = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, errRet = me.client.UseVpcClient().DeleteVpnGateway(request)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		return nil
	})
	return
}

func (me *VpcService) DescribeCustomerGatewayByFilter(ctx context.Context, filters map[string]string) (instances []*vpc.CustomerGateway, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewDescribeCustomerGatewaysRequest()
	)
	request.Filters = make([]*vpc.Filter, 0, len(filters))
	for k, v := range filters {
		filter := vpc.Filter{
			Name:   helper.String(k),
			Values: []*string{helper.String(v)},
		}
		request.Filters = append(request.Filters, &filter)
	}

	var offset uint64 = 0
	var pageSize uint64 = 100
	instances = make([]*vpc.CustomerGateway, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseVpcClient().DescribeCustomerGateways(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.CustomerGatewaySet) < 1 {
			break
		}
		instances = append(instances, response.Response.CustomerGatewaySet...)
		if len(response.Response.CustomerGatewaySet) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

func (me *VpcService) DeleteCustomerGateway(ctx context.Context, customerGatewayId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewDeleteCustomerGatewayRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.CustomerGatewayId = &customerGatewayId

	errRet = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, errRet = me.client.UseVpcClient().DeleteCustomerGateway(request)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		return nil
	})
	return
}

func (me *VpcService) CreateAddressTemplate(ctx context.Context, name string, addresses []interface{}) (templateId string, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewCreateAddressTemplateRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, request.GetAction(), request.ToJsonString(), errRet)
		}
	}()

	request.AddressTemplateName = &name
	request.Addresses = make([]*string, len(addresses))
	for i, v := range addresses {
		request.Addresses[i] = helper.String(v.(string))
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().CreateAddressTemplate(request)
	if err != nil {
		errRet = err
		return
	}

	if response == nil || response.Response == nil || response.Response.AddressTemplate == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
	}

	templateId = *response.Response.AddressTemplate.AddressTemplateId
	return
}

func (me *VpcService) DescribeAddressTemplateById(ctx context.Context, templateId string) (template *vpc.AddressTemplate, has bool, errRet error) {
	filter := vpc.Filter{Name: helper.String("address-template-id"), Values: []*string{&templateId}}
	templates, err := me.DescribeAddressTemplates(ctx, []*vpc.Filter{&filter})
	if err != nil {
		errRet = err
		return
	}

	if len(templates) == 0 {
		return
	}
	if len(templates) > 1 {
		errRet = fmt.Errorf("TencentCloud SDK return more than one templates, instanceId %s", templateId)
	}

	has = true
	template = templates[0]
	return
}

func (me *VpcService) DescribeAddressTemplates(ctx context.Context, filter []*vpc.Filter) (templateList []*vpc.AddressTemplate, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewDescribeAddressTemplatesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	var offset, limit = 0, 100
	request.Filters = filter

	for {
		request.Offset = helper.String(strconv.Itoa(offset))
		request.Limit = helper.String(strconv.Itoa(limit))

		ratelimit.Check(request.GetAction())
		response, err := me.client.UseVpcClient().DescribeAddressTemplates(request)
		if err != nil {
			errRet = err
			return
		}
		if response == nil || response.Response == nil {
			errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
		}
		templateList = append(templateList, response.Response.AddressTemplateSet...)
		if len(response.Response.AddressTemplateSet) < limit {
			return
		}
		offset += limit
	}
}

func (me *VpcService) ModifyAddressTemplate(ctx context.Context, templateId string, name string, addresses []interface{}) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewModifyAddressTemplateAttributeRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, request.GetAction(), request.ToJsonString(), errRet)
		}
	}()

	request.AddressTemplateId = &templateId
	request.AddressTemplateName = &name
	request.Addresses = make([]*string, len(addresses))
	for i, v := range addresses {
		request.Addresses[i] = helper.String(v.(string))
	}

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseVpcClient().ModifyAddressTemplateAttribute(request)
	return err
}

func (me *VpcService) DeleteAddressTemplate(ctx context.Context, templateId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewDeleteAddressTemplateRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.AddressTemplateId = &templateId

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseVpcClient().DeleteAddressTemplate(request)
	return err
}

func (me *VpcService) CreateAddressTemplateGroup(ctx context.Context, name string, addressTemplate []interface{}) (templateId string, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewCreateAddressTemplateGroupRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, request.GetAction(), request.ToJsonString(), errRet)
		}
	}()

	request.AddressTemplateGroupName = &name
	request.AddressTemplateIds = make([]*string, len(addressTemplate))
	for i, v := range addressTemplate {
		request.AddressTemplateIds[i] = helper.String(v.(string))
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().CreateAddressTemplateGroup(request)
	if err != nil {
		errRet = err
		return
	}

	if response == nil || response.Response == nil || response.Response.AddressTemplateGroup == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
	}

	templateId = *response.Response.AddressTemplateGroup.AddressTemplateGroupId
	return
}

func (me *VpcService) ModifyAddressTemplateGroup(ctx context.Context, templateGroupId string, name string, templateIds []interface{}) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewModifyAddressTemplateGroupAttributeRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, request.GetAction(), request.ToJsonString(), errRet)
		}
	}()

	request.AddressTemplateGroupId = &templateGroupId
	request.AddressTemplateGroupName = &name
	request.AddressTemplateIds = make([]*string, len(templateIds))
	for i, v := range templateIds {
		request.AddressTemplateIds[i] = helper.String(v.(string))
	}

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseVpcClient().ModifyAddressTemplateGroupAttribute(request)
	return err
}

func (me *VpcService) DescribeAddressTemplateGroupById(ctx context.Context, templateGroupId string) (templateGroup *vpc.AddressTemplateGroup, has bool, errRet error) {
	filter := vpc.Filter{Name: helper.String("address-template-group-id"), Values: []*string{&templateGroupId}}
	templateGroups, err := me.DescribeAddressTemplateGroups(ctx, []*vpc.Filter{&filter})
	if err != nil {
		errRet = err
		return
	}

	if len(templateGroups) == 0 {
		return
	}
	if len(templateGroups) > 1 {
		errRet = fmt.Errorf("TencentCloud SDK return more than one template group, instanceId %s", templateGroupId)
	}

	has = true
	templateGroup = templateGroups[0]
	return
}

func (me *VpcService) DescribeAddressTemplateGroups(ctx context.Context, filter []*vpc.Filter) (templateList []*vpc.AddressTemplateGroup, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewDescribeAddressTemplateGroupsRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	var offset, limit = 0, 100
	request.Filters = filter

	for {
		request.Offset = helper.String(strconv.Itoa(offset))
		request.Limit = helper.String(strconv.Itoa(limit))

		ratelimit.Check(request.GetAction())
		response, err := me.client.UseVpcClient().DescribeAddressTemplateGroups(request)
		if err != nil {
			errRet = err
			return
		}
		if response == nil || response.Response == nil {
			errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
		}
		templateList = append(templateList, response.Response.AddressTemplateGroupSet...)
		if len(response.Response.AddressTemplateGroupSet) < limit {
			return
		}
		offset += limit
	}
}

func (me *VpcService) DeleteAddressTemplateGroup(ctx context.Context, templateGroupId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewDeleteAddressTemplateGroupRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.AddressTemplateGroupId = &templateGroupId

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseVpcClient().DeleteAddressTemplateGroup(request)
	return err
}

func (me *VpcService) CreateServiceTemplate(ctx context.Context, name string, services []interface{}) (templateId string, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewCreateServiceTemplateRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, request.GetAction(), request.ToJsonString(), errRet)
		}
	}()

	request.ServiceTemplateName = &name
	request.Services = make([]*string, len(services))
	for i, v := range services {
		request.Services[i] = helper.String(v.(string))
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().CreateServiceTemplate(request)
	if err != nil {
		errRet = err
		return
	}

	if response == nil || response.Response == nil || response.Response.ServiceTemplate == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
	}

	templateId = *response.Response.ServiceTemplate.ServiceTemplateId
	return
}

func (me *VpcService) ModifyServiceTemplate(ctx context.Context, templateId string, name string, services []interface{}) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewModifyServiceTemplateAttributeRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, request.GetAction(), request.ToJsonString(), errRet)
		}
	}()

	request.ServiceTemplateId = &templateId
	request.ServiceTemplateName = &name
	request.Services = make([]*string, len(services))
	for i, v := range services {
		request.Services[i] = helper.String(v.(string))
	}

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseVpcClient().ModifyServiceTemplateAttribute(request)
	return err
}

func (me *VpcService) DescribeServiceTemplateById(ctx context.Context, templateId string) (template *vpc.ServiceTemplate, has bool, errRet error) {
	filter := vpc.Filter{Name: helper.String("service-template-id"), Values: []*string{&templateId}}
	templates, err := me.DescribeServiceTemplates(ctx, []*vpc.Filter{&filter})
	if err != nil {
		errRet = err
		return
	}

	if len(templates) == 0 {
		return
	}
	if len(templates) > 1 {
		errRet = fmt.Errorf("TencentCloud SDK return more than one templates, instanceId %s", templateId)
	}

	has = true
	template = templates[0]
	return
}

func (me *VpcService) DescribeServiceTemplates(ctx context.Context, filter []*vpc.Filter) (templateList []*vpc.ServiceTemplate, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewDescribeServiceTemplatesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	var offset, limit = 0, 100
	request.Filters = filter

	for {
		request.Offset = helper.String(strconv.Itoa(offset))
		request.Limit = helper.String(strconv.Itoa(limit))

		ratelimit.Check(request.GetAction())
		response, err := me.client.UseVpcClient().DescribeServiceTemplates(request)
		if err != nil {
			errRet = err
			return
		}
		if response == nil || response.Response == nil {
			errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
		}
		templateList = append(templateList, response.Response.ServiceTemplateSet...)
		if len(response.Response.ServiceTemplateSet) < limit {
			return
		}
		offset += limit
	}
}

func (me *VpcService) DeleteServiceTemplate(ctx context.Context, templateId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewDeleteServiceTemplateRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.ServiceTemplateId = &templateId

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseVpcClient().DeleteServiceTemplate(request)
	return err
}

func (me *VpcService) CreateServiceTemplateGroup(ctx context.Context, name string, serviceTemplate []interface{}) (templateId string, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewCreateServiceTemplateGroupRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, request.GetAction(), request.ToJsonString(), errRet)
		}
	}()

	request.ServiceTemplateGroupName = &name
	request.ServiceTemplateIds = make([]*string, len(serviceTemplate))
	for i, v := range serviceTemplate {
		request.ServiceTemplateIds[i] = helper.String(v.(string))
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().CreateServiceTemplateGroup(request)
	if err != nil {
		errRet = err
		return
	}

	if response == nil || response.Response == nil || response.Response.ServiceTemplateGroup == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
	}

	templateId = *response.Response.ServiceTemplateGroup.ServiceTemplateGroupId
	return
}

func (me *VpcService) DescribeServiceTemplateGroupById(ctx context.Context, templateGroupId string) (template *vpc.ServiceTemplateGroup, has bool, errRet error) {
	filter := vpc.Filter{Name: helper.String("service-template-group-id"), Values: []*string{&templateGroupId}}
	templates, err := me.DescribeServiceTemplateGroups(ctx, []*vpc.Filter{&filter})
	if err != nil {
		errRet = err
		return
	}

	if len(templates) == 0 {
		return
	}
	if len(templates) > 1 {
		errRet = fmt.Errorf("TencentCloud SDK return more than one templates, instanceId %s", templateGroupId)
	}

	has = true
	template = templates[0]
	return
}

func (me *VpcService) DescribeServiceTemplateGroups(ctx context.Context, filter []*vpc.Filter) (templateList []*vpc.ServiceTemplateGroup, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewDescribeServiceTemplateGroupsRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	var offset, limit = 0, 100
	request.Filters = filter

	for {
		request.Offset = helper.String(strconv.Itoa(offset))
		request.Limit = helper.String(strconv.Itoa(limit))

		ratelimit.Check(request.GetAction())
		response, err := me.client.UseVpcClient().DescribeServiceTemplateGroups(request)
		if err != nil {
			errRet = err
			return
		}
		if response == nil || response.Response == nil {
			errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
		}
		templateList = append(templateList, response.Response.ServiceTemplateGroupSet...)
		if len(response.Response.ServiceTemplateGroupSet) < limit {
			return
		}
		offset += limit
	}
}

func (me *VpcService) ModifyServiceTemplateGroup(ctx context.Context, serviceGroupId string, name string, templateIds []interface{}) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewModifyServiceTemplateGroupAttributeRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
				logId, request.GetAction(), request.ToJsonString(), errRet)
		}
	}()

	request.ServiceTemplateGroupId = &serviceGroupId
	request.ServiceTemplateGroupName = &name
	request.ServiceTemplateIds = make([]*string, len(templateIds))
	for i, v := range templateIds {
		request.ServiceTemplateIds[i] = helper.String(v.(string))
	}

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseVpcClient().ModifyServiceTemplateGroupAttribute(request)
	return err
}

func (me *VpcService) DeleteServiceTemplateGroup(ctx context.Context, templateGroupId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewDeleteServiceTemplateGroupRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.ServiceTemplateGroupId = &templateGroupId

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseVpcClient().DeleteServiceTemplateGroup(request)
	return err
}

func (me *VpcService) CreateVpnGatewayRoute(ctx context.Context, vpnGatewayId string, vpnGwRoutes []*vpc.VpnGatewayRoute) (errRet error, routes []*vpc.VpnGatewayRoute) {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewCreateVpnGatewayRoutesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.VpnGatewayId = &vpnGatewayId
	request.Routes = vpnGwRoutes

	var response *vpc.CreateVpnGatewayRoutesResponse
	errRet = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, errRet = me.client.UseVpcClient().CreateVpnGatewayRoutes(request)
		if errRet != nil {
			log.Printf("[CRITAL]%s create vpn gateway route failed, reason: %v", logId, errRet)
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		return nil
	})
	if errRet != nil {
		return errRet, nil
	}

	if response == nil || response.Response == nil || response.Response.Routes == nil || len(response.Response.Routes) == 0 {
		errRet = fmt.Errorf("TencentCloud SDK return nil response, %+v, %s", response, request.GetAction())
	} else {
		routes = response.Response.Routes
	}
	return
}

func (me *VpcService) ModifyVpnGatewayRoute(ctx context.Context, vpnGatewayId, routeId, status string) (errRet error, routes *vpc.VpnGatewayRoute) {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewModifyVpnGatewayRoutesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.VpnGatewayId = &vpnGatewayId
	request.Routes = []*vpc.VpnGatewayRouteModify{{
		RouteId: &routeId,
		Status:  &status,
	}}

	var response *vpc.ModifyVpnGatewayRoutesResponse
	errRet = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, errRet = me.client.UseVpcClient().ModifyVpnGatewayRoutes(request)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		return nil
	})
	if errRet != nil {
		return errRet, nil
	}

	if response == nil || response.Response == nil || response.Response.Routes == nil || len(response.Response.Routes) == 0 {
		errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
	} else {
		routes = response.Response.Routes[0]
	}
	return
}

func (me *VpcService) DeleteVpnGatewayRoutes(ctx context.Context, vpnGatewayId string, routeIds []*string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewDeleteVpnGatewayRoutesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.VpnGatewayId = &vpnGatewayId
	request.RouteIds = routeIds

	errRet = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, errRet = me.client.UseVpcClient().DeleteVpnGatewayRoutes(request)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		return nil
	})
	return
}

func (me *VpcService) DescribeVpnGatewayRoutes(ctx context.Context, vpnGatewayId string, filters []*vpc.Filter) (errRet error, result []*vpc.VpnGatewayRoute) {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewDescribeVpnGatewayRoutesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.VpnGatewayId = &vpnGatewayId
	if len(filters) > 0 {
		request.Filters = filters
	}

	offset := int64(0)
	limit := int64(VPN_DESCRIBE_LIMIT)
	for {
		request.Offset = &offset
		request.Limit = &limit
		var response *vpc.DescribeVpnGatewayRoutesResponse
		errRet = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			response, errRet = me.client.UseVpcClient().DescribeVpnGatewayRoutes(request)
			if errRet != nil {
				return tccommon.RetryError(errRet, tccommon.InternalError)
			}
			return nil
		})
		if errRet != nil {
			return errRet, nil
		}

		if response == nil || response.Response == nil {
			return fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction()), nil
		} else if len(response.Response.Routes) > 0 {
			result = append(result, response.Response.Routes...)
		} else {
			return
		}
		offset = offset + limit
	}
}

func (me *VpcService) DescribeVpcTaskResult(ctx context.Context, taskId *string) (err error) {

	logId := tccommon.GetLogId(ctx)
	request := vpc.NewDescribeVpcTaskResultRequest()
	defer func() {
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), err.Error())
		}
	}()
	request.TaskId = taskId
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseVpcClient().DescribeVpcTaskResult(request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		if response.Response.Status != nil && *response.Response.Status == VPN_TASK_STATUS_RUNNING {
			return resource.RetryableError(errors.New("VPN task is running"))
		}
		return nil
	})
	if err != nil {
		return err
	}
	return
}

func (me *VpcService) DescribeTaskResult(ctx context.Context, taskId *uint64) (result *vpc.DescribeTaskResultResponse, err error) {

	logId := tccommon.GetLogId(ctx)
	request := vpc.NewDescribeTaskResultRequest()
	defer func() {
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), err.Error())
		}
	}()
	request.TaskId = taskId
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseVpcClient().DescribeTaskResult(request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		result = response
		return nil
	})
	if err != nil {
		return nil, err
	}
	return
}

func (me *VpcService) DescribeVpnSslServerById(ctx context.Context, sslId string) (has bool, gateway *vpc.SslVpnSever, err error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = vpc.NewDescribeVpnGatewaySslServersRequest()
		response *vpc.DescribeVpnGatewaySslServersResponse
	)
	request.SslVpnServerIds = []*string{&sslId}
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		response, err = me.client.UseVpcClient().DescribeVpnGatewaySslServers(request)
		if err != nil {
			ee, ok := err.(*sdkErrors.TencentCloudSDKError)
			if !ok {
				return tccommon.RetryError(err)
			}
			if ee.Code == VPCNotFound {
				return nil
			} else {
				return tccommon.RetryError(err)
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]", logId, request.GetAction(), request.ToJsonString(), err)
		return
	}
	if response == nil || response.Response == nil || len(response.Response.SslVpnSeverSet) < 1 {
		has = false
		return
	}

	gateway = response.Response.SslVpnSeverSet[0]
	has = true
	return
}

func (me *VpcService) DescribeVpnGwSslServerByFilter(ctx context.Context, filters map[string]string) (instances []*vpc.SslVpnSever, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewDescribeVpnGatewaySslServersRequest()
	)
	request.Filters = make([]*vpc.FilterObject, 0, len(filters))
	for k, v := range filters {
		filter := vpc.FilterObject{
			Name:   helper.String(k),
			Values: []*string{helper.String(v)},
		}
		request.Filters = append(request.Filters, &filter)
	}

	var offset uint64 = 0
	var pageSize uint64 = 100
	instances = make([]*vpc.SslVpnSever, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseVpcClient().DescribeVpnGatewaySslServers(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.SslVpnSeverSet) < 1 {
			break
		}
		instances = append(instances, response.Response.SslVpnSeverSet...)
		if len(response.Response.SslVpnSeverSet) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

func (me *VpcService) DeleteVpnGatewaySslServer(ctx context.Context, SslServerId string) (taskId uint64, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewDeleteVpnGatewaySslServerRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.SslVpnServerId = &SslServerId

	errRet = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, errRet := me.client.UseVpcClient().DeleteVpnGatewaySslServer(request)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		taskId = *response.Response.TaskId
		return nil
	})
	return
}

func (me *VpcService) DescribeVpnSslClientById(ctx context.Context, sslId string) (has bool, gateway *vpc.SslVpnClient, err error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = vpc.NewDescribeVpnGatewaySslClientsRequest()
		response *vpc.DescribeVpnGatewaySslClientsResponse
	)
	request.SslVpnClientIds = []*string{&sslId}
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		response, err = me.client.UseVpcClient().DescribeVpnGatewaySslClients(request)
		if err != nil {
			ee, ok := err.(*sdkErrors.TencentCloudSDKError)
			if !ok {
				return tccommon.RetryError(err)
			}
			if ee.Code == VPCNotFound {
				return nil
			} else {
				return tccommon.RetryError(err)
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]", logId, request.GetAction(), request.ToJsonString(), err)
		return
	}
	if response == nil || response.Response == nil || len(response.Response.SslVpnClientSet) < 1 {
		has = false
		return
	}

	gateway = response.Response.SslVpnClientSet[0]
	has = true
	return
}

func (me *VpcService) DescribeVpnGwSslClientByFilter(ctx context.Context, filters map[string]string) (instances []*vpc.SslVpnClient, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewDescribeVpnGatewaySslClientsRequest()
	)
	request.Filters = make([]*vpc.Filter, 0, len(filters))
	for k, v := range filters {
		filter := vpc.Filter{
			Name:   helper.String(k),
			Values: []*string{helper.String(v)},
		}
		request.Filters = append(request.Filters, &filter)
	}

	var offset uint64 = 0
	var pageSize uint64 = 100
	instances = make([]*vpc.SslVpnClient, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseVpcClient().DescribeVpnGatewaySslClients(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.SslVpnClientSet) < 1 {
			break
		}
		instances = append(instances, response.Response.SslVpnClientSet...)
		if len(response.Response.SslVpnClientSet) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

func (me *VpcService) DeleteVpnGatewaySslClient(ctx context.Context, SslClientId string) (taskId *uint64, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewDeleteVpnGatewaySslClientRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.SslVpnClientId = &SslClientId

	errRet = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, errRet := me.client.UseVpcClient().DeleteVpnGatewaySslClient(request)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		taskId = response.Response.TaskId
		return nil
	})
	return
}

func (me *VpcService) CreateNatGatewaySnat(ctx context.Context, natGatewayId string, snat *vpc.SourceIpTranslationNatRule) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewCreateNatGatewaySourceIpTranslationNatRuleRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.NatGatewayId = &natGatewayId
	request.SourceIpTranslationNatRules = []*vpc.SourceIpTranslationNatRule{snat}

	var response *vpc.CreateNatGatewaySourceIpTranslationNatRuleResponse
	errRet = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, errRet = me.client.UseVpcClient().CreateNatGatewaySourceIpTranslationNatRule(request)
		if errRet != nil {
			log.Printf("[CRITAL]%s create nat gateway source ip translation nat rule failed, reason: %v", logId, errRet)
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		return nil
	})
	if errRet != nil {
		return errRet
	}

	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response, %+v, %s", response, request.GetAction())
	}
	return
}

func (me *VpcService) ModifyNatGatewaySnat(ctx context.Context, natGatewayId string, snat *vpc.SourceIpTranslationNatRule) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewModifyNatGatewaySourceIpTranslationNatRuleRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.NatGatewayId = &natGatewayId
	request.SourceIpTranslationNatRule = snat

	var response *vpc.ModifyNatGatewaySourceIpTranslationNatRuleResponse
	errRet = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, errRet = me.client.UseVpcClient().ModifyNatGatewaySourceIpTranslationNatRule(request)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		return nil
	})
	if errRet != nil {
		return errRet
	}

	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
	}
	return
}

func (me *VpcService) DeleteNatGatewaySnat(ctx context.Context, natGatewayId string, snatId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewDeleteNatGatewaySourceIpTranslationNatRuleRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.NatGatewayId = &natGatewayId
	request.NatGatewaySnatIds = []*string{&snatId}

	errRet = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, errRet = me.client.UseVpcClient().DeleteNatGatewaySourceIpTranslationNatRule(request)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		return nil
	})
	return
}

func (me *VpcService) DescribeNatGatewaySnats(ctx context.Context, natGatewayId string, filters []*vpc.Filter) (errRet error, result []*vpc.SourceIpTranslationNatRule) {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewDescribeNatGatewaySourceIpTranslationNatRulesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.NatGatewayId = &natGatewayId
	if len(filters) > 0 {
		request.Filters = filters
	}

	offset := int64(0)
	limit := int64(VPN_DESCRIBE_LIMIT)
	for {
		request.Offset = &offset
		request.Limit = &limit
		var response *vpc.DescribeNatGatewaySourceIpTranslationNatRulesResponse
		errRet = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			response, errRet = me.client.UseVpcClient().DescribeNatGatewaySourceIpTranslationNatRules(request)
			if errRet != nil {
				return tccommon.RetryError(errRet, tccommon.InternalError)
			}
			return nil
		})
		if errRet != nil {
			return errRet, nil
		}

		if response == nil || response.Response == nil {
			return fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction()), nil
		} else if len(response.Response.SourceIpTranslationNatRuleSet) > 0 {
			result = append(result, response.Response.SourceIpTranslationNatRuleSet...)
		} else {
			return
		}
		offset = offset + limit
	}
}

func (me *VpcService) DescribeAssistantCidr(ctx context.Context, vpcId string) (info []*vpc.AssistantCidr, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vpc.NewDescribeAssistantCidrRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.VpcIds = []*string{&vpcId}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().DescribeAssistantCidr(request)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	info = response.Response.AssistantCidrSet

	return
}

// CheckAssistantCidr used for check if cidr conflict
func (me *VpcService) CheckAssistantCidr(ctx context.Context, request *vpc.CheckAssistantCidrRequest) (info []*vpc.ConflictSource, errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().CheckAssistantCidr(request)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	info = response.Response.ConflictSourceSet

	return
}

func (me *VpcService) CreateAssistantCidr(ctx context.Context, request *vpc.CreateAssistantCidrRequest) (info []*vpc.AssistantCidr, errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().CreateAssistantCidr(request)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	info = response.Response.AssistantCidrSet

	return
}

func (me *VpcService) ModifyAssistantCidr(ctx context.Context, request *vpc.ModifyAssistantCidrRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().ModifyAssistantCidr(request)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *VpcService) DeleteAssistantCidr(ctx context.Context, request *vpc.DeleteAssistantCidrRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().DeleteAssistantCidr(request)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *VpcService) DescribeVpcBandwidthPackage(ctx context.Context, bandwidthPackageId string) (resource *vpc.BandwidthPackage, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewDescribeBandwidthPackagesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.BandwidthPackageIds = []*string{&bandwidthPackageId}
	//request.Filters = append(
	//	request.Filters,
	//	&bwp.Filter{
	//		Name:   helper.String("bandwidth-package_id"),
	//		Values: []*string{&bandwidthPackageId},
	//	},
	//)
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().DescribeBandwidthPackages(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response != nil && len(response.Response.BandwidthPackageSet) > 0 {
		resource = response.Response.BandwidthPackageSet[0]
	}

	return
}

func (me *VpcService) DeleteVpcBandwidthPackageById(ctx context.Context, bandwidthPackageId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDeleteBandwidthPackageRequest()

	request.BandwidthPackageId = &bandwidthPackageId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().DeleteBandwidthPackage(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *VpcService) DescribeVpcBandwidthPackageAttachment(ctx context.Context, bandwidthPackageId, resourceId string) (bandwidthPackageResources *vpc.Resource, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewDescribeBandwidthPackageResourcesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.BandwidthPackageId = &bandwidthPackageId
	request.Filters = append(
		request.Filters,
		&vpc.Filter{
			Name:   helper.String("resource-id"),
			Values: []*string{&resourceId},
		},
	)

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().DescribeBandwidthPackageResources(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.ResourceSet) < 1 {
		return
	}
	bandwidthPackageResources = response.Response.ResourceSet[0]

	return

}

func (me *VpcService) DeleteVpcBandwidthPackageAttachmentById(ctx context.Context, bandwidthPackageId, resourceId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewRemoveBandwidthPackageResourcesRequest()

	if strings.HasPrefix(resourceId, "eip") {
		request.ResourceType = helper.String("Address")
	} else {
		request.ResourceType = helper.String("LoadBalance")
	}

	request.BandwidthPackageId = &bandwidthPackageId
	request.ResourceIds = []*string{&resourceId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().RemoveBandwidthPackageResources(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *VpcService) DescribeFlowLogs(ctx context.Context, request *vpc.DescribeFlowLogsRequest) (result []*vpc.FlowLog, errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().DescribeFlowLogs(request)

	if err != nil {
		errRet = err
		return
	}

	result = response.Response.FlowLog

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *VpcService) DescribeVpcFlowLogById(ctx context.Context, flowLogId, vpcId string) (FlowLog *vpc.FlowLog, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDescribeFlowLogRequest()
	request.FlowLogId = &flowLogId

	if vpcId != "" {
		request.VpcId = &vpcId
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DescribeFlowLog(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.FlowLog) < 1 {
		return
	}

	FlowLog = response.Response.FlowLog[0]
	return
}

func (me *VpcService) DeleteVpcFlowLogById(ctx context.Context, flowLogId, vpcId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDeleteFlowLogRequest()
	request.FlowLogId = &flowLogId
	if vpcId != "" {
		request.VpcId = &vpcId
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DeleteFlowLog(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *VpcService) DescribeVpcEndPointServiceById(ctx context.Context, endPointServiceId string) (endPointService *vpc.EndPointService, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDescribeVpcEndPointServiceRequest()
	request.EndPointServiceIds = []*string{&endPointServiceId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	instances := make([]*vpc.EndPointService, 0)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseVpcClient().DescribeVpcEndPointService(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.EndPointServiceSet) < 1 {
			break
		}
		instances = append(instances, response.Response.EndPointServiceSet...)
		if len(response.Response.EndPointServiceSet) < int(limit) {
			break
		}

		offset += limit
	}

	if len(instances) < 1 {
		return
	}
	endPointService = instances[0]
	return
}

func (me *VpcService) DeleteVpcEndPointServiceById(ctx context.Context, endPointServiceId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDeleteVpcEndPointServiceRequest()
	request.EndPointServiceId = &endPointServiceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DeleteVpcEndPointService(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *VpcService) DescribeVpcEndPointById(ctx context.Context, endPointId string) (endPoint *vpc.EndPoint, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDescribeVpcEndPointRequest()
	request.EndPointId = []*string{&endPointId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	instances := make([]*vpc.EndPoint, 0)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseVpcClient().DescribeVpcEndPoint(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.EndPointSet) < 1 {
			break
		}
		instances = append(instances, response.Response.EndPointSet...)
		if len(response.Response.EndPointSet) < int(limit) {
			break
		}

		offset += limit
	}

	if len(instances) < 1 {
		return
	}
	endPoint = instances[0]
	return
}

func (me *VpcService) DeleteVpcEndPointById(ctx context.Context, endPointId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDeleteVpcEndPointRequest()
	request.EndPointId = &endPointId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DeleteVpcEndPoint(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *VpcService) DescribeVpcEndPointServiceWhiteListById(ctx context.Context, userUin string, endPointServiceId string) (endPointServiceWhiteList *vpc.VpcEndPointServiceUser, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDescribeVpcEndPointServiceWhiteListRequest()

	request.Filters = make([]*vpc.Filter, 0)
	if userUin != "" {
		filter := &vpc.Filter{
			Name:   helper.String("user-uin"),
			Values: []*string{&userUin},
		}
		request.Filters = append(request.Filters, filter)
	}
	if endPointServiceId != "" {
		filter := &vpc.Filter{
			Name:   helper.String("end-point-service-id"),
			Values: []*string{&endPointServiceId},
		}
		request.Filters = append(request.Filters, filter)
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	instances := make([]*vpc.VpcEndPointServiceUser, 0)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseVpcClient().DescribeVpcEndPointServiceWhiteList(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.VpcEndpointServiceUserSet) < 1 {
			break
		}
		instances = append(instances, response.Response.VpcEndpointServiceUserSet...)
		if len(response.Response.VpcEndpointServiceUserSet) < int(limit) {
			break
		}

		offset += limit
	}

	if len(instances) < 1 {
		return
	}
	endPointServiceWhiteList = instances[0]
	return
}

func (me *VpcService) DeleteVpcEndPointServiceWhiteListById(ctx context.Context, userUin string, endPointServiceId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDeleteVpcEndPointServiceWhiteListRequest()
	request.UserUin = []*string{&userUin}
	request.EndPointServiceId = &endPointServiceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DeleteVpcEndPointServiceWhiteList(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *VpcService) DescribeVpcBandwidthPackageByEip(ctx context.Context, eipId string) (resource *vpc.BandwidthPackage, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewDescribeBandwidthPackagesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.Filters = append(
		request.Filters,
		&vpc.Filter{
			Name:   helper.String("resource.resource-id"),
			Values: []*string{&eipId},
		},
	)
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().DescribeBandwidthPackages(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response != nil && len(response.Response.BandwidthPackageSet) > 0 {
		resource = response.Response.BandwidthPackageSet[0]
	}

	return
}

func (me *VpcService) DescribeVpcCcnRoutesById(ctx context.Context, ccnId string, routeId string) (ccnRoutes *vpc.CcnRoute, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDescribeCcnRoutesRequest()
	request.CcnId = &ccnId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DescribeCcnRoutes(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	for _, route := range response.Response.RouteSet {
		if *route.RouteId == routeId {
			ccnRoutes = route
			return
		}
	}

	return
}

func (me *VpcService) DescribeCcnCrossBorderComplianceByFilter(ctx context.Context, param map[string]interface{}) (crossBorderCompliance []*vpc.CrossBorderCompliance, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewDescribeCrossBorderComplianceRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "service_provider" {
			request.ServiceProvider = v.(*string)
		}
		if k == "compliance_id" {
			if *v.(*uint64) != 0 {
				request.ComplianceId = v.(*uint64)
			}
		}
		if k == "company" {
			request.Company = v.(*string)
		}
		if k == "uniform_social_credit_code" {
			request.UniformSocialCreditCode = v.(*string)
		}
		if k == "legal_person" {
			request.LegalPerson = v.(*string)
		}
		if k == "issuing_authority" {
			request.IssuingAuthority = v.(*string)
		}
		if k == "business_address" {
			request.BusinessAddress = v.(*string)
		}
		if k == "post_code" {
			if *v.(*uint64) != 0 {
				request.PostCode = v.(*uint64)
			}
		}
		if k == "manager" {
			request.Manager = v.(*string)
		}
		if k == "manager_id" {
			request.ManagerId = v.(*string)
		}
		if k == "manager_address" {
			request.ManagerAddress = v.(*string)
		}
		if k == "manager_telephone" {
			request.ManagerTelephone = v.(*string)
		}
		if k == "email" {
			request.Email = v.(*string)
		}
		if k == "service_start_date" {
			request.ServiceStartDate = v.(*string)
		}
		if k == "service_end_date" {
			request.ServiceEndDate = v.(*string)
		}
		if k == "state" {
			request.State = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseVpcClient().DescribeCrossBorderCompliance(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.CrossBorderComplianceSet) < 1 {
			break
		}
		crossBorderCompliance = append(crossBorderCompliance, response.Response.CrossBorderComplianceSet...)
		if len(response.Response.CrossBorderComplianceSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *VpcService) DescribeTenantCcnByFilter(ctx context.Context, param map[string]interface{}) (tenantCcn []*vpc.CcnInstanceInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewDescribeTenantCcnsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.Filters = make([]*vpc.Filter, 0, len(param))
	for k, v := range param {
		filter := &vpc.Filter{
			Name:   helper.String(k),
			Values: v.([]*string),
		}
		request.Filters = append(request.Filters, filter)
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseVpcClient().DescribeTenantCcns(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.CcnSet) < 1 {
			break
		}
		tenantCcn = append(tenantCcn, response.Response.CcnSet...)
		if len(response.Response.CcnSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *VpcService) DescribeCcnCrossBorderFlowMonitorByFilter(ctx context.Context, param map[string]interface{}) (crossBorderFlowMonitor []*vpc.CrossBorderFlowMonitorData, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewDescribeCrossBorderFlowMonitorRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "source_region" {
			request.SourceRegion = v.(*string)
		}
		if k == "destination_region" {
			request.DestinationRegion = v.(*string)
		}
		if k == "ccn_id" {
			request.CcnId = v.(*string)
		}
		if k == "ccn_uin" {
			request.CcnUin = v.(*string)
		}
		if k == "period" {
			if *v.(*int64) != 0 {
				request.Period = v.(*int64)
			}
		}
		if k == "start_time" {
			request.StartTime = v.(*string)
		}
		if k == "end_time" {
			request.EndTime = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DescribeCrossBorderFlowMonitor(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.CrossBorderFlowMonitorData) < 1 {
		return
	}

	crossBorderFlowMonitor = response.Response.CrossBorderFlowMonitorData

	return
}

func (me *VpcService) DescribeVpnCustomerGatewayVendors(ctx context.Context) (vpnCustomerGatewayVendors []*vpc.CustomerGatewayVendor, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewDescribeCustomerGatewayVendorsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DescribeCustomerGatewayVendors(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.CustomerGatewayVendorSet) < 1 {
		return
	}

	vpnCustomerGatewayVendors = response.Response.CustomerGatewayVendorSet
	return
}

func (me *VpcService) DescribeVpcVpnGatewayCcnRoutesById(ctx context.Context, vpnGatewayId string, routeId string) (vpnGatewayCcnRoutes *vpc.VpngwCcnRoutes, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDescribeVpnGatewayCcnRoutesRequest()
	request.VpnGatewayId = &vpnGatewayId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DescribeVpnGatewayCcnRoutes(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.RouteSet) < 1 {
		return
	}

	for _, route := range response.Response.RouteSet {
		if *route.RouteId == routeId {
			vpnGatewayCcnRoutes = route
			break
		}
	}
	return
}

func (me *VpcService) DescribeVpcIpv6AddressById(ctx context.Context, ip6AddressId string) (ipv6Address *vpc.Address, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDescribeIp6AddressesRequest()
	request.Ip6AddressIds = []*string{&ip6AddressId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DescribeIp6Addresses(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.AddressSet) < 1 {
		return
	}

	ipv6Address = response.Response.AddressSet[0]
	return
}

func (me *VpcService) DeleteVpcIpv6AddressById(ctx context.Context, ip6AddressId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewReleaseIp6AddressesBandwidthRequest()
	request.Ip6AddressIds = []*string{&ip6AddressId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().ReleaseIp6AddressesBandwidth(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *VpcService) VpcIpv6AddressStateRefreshFunc(taskId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := tccommon.ContextNil

		taskId := helper.StrToUint64Point(taskId)

		object, err := me.DescribeTaskResult(ctx, taskId)

		if err != nil {
			return nil, "", err
		}

		return object, helper.PString(object.Response.Result), nil
	}
}

func (me *VpcService) DescribeVpcCcnRegionBandwidthLimitsByFilter(ctx context.Context, param map[string]interface{}) (CcnRegionBandwidthLimits []*vpc.CcnBandwidth, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewDescribeCrossBorderCcnRegionBandwidthLimitsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "filters" {
			request.Filters = v.([]*vpc.Filter)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseVpcClient().DescribeCrossBorderCcnRegionBandwidthLimits(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.CcnBandwidthSet) < 1 {
			break
		}
		CcnRegionBandwidthLimits = append(CcnRegionBandwidthLimits, response.Response.CcnBandwidthSet...)
		if len(response.Response.CcnBandwidthSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *VpcService) DescribeNatDcRouteByFilter(ctx context.Context, param map[string]interface{}) (natDcRoute []*vpc.NatDirectConnectGatewayRoute, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewDescribeNatGatewayDirectConnectGatewayRouteRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "NatGatewayId" {
			request.NatGatewayId = v.(*string)
		}
		if k == "VpcId" {
			request.VpcId = v.(*string)
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
		response, err := me.client.UseVpcClient().DescribeNatGatewayDirectConnectGatewayRoute(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.NatDirectConnectGatewayRouteSet) < 1 {
			break
		}
		natDcRoute = append(natDcRoute, response.Response.NatDirectConnectGatewayRouteSet...)
		if len(response.Response.NatDirectConnectGatewayRouteSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *VpcService) DescribeEipAddressQuota(ctx context.Context) (addressQuota []*vpc.Quota, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewDescribeAddressQuotaRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DescribeAddressQuota(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	addressQuota = append(addressQuota, response.Response.QuotaSet...)

	return
}

func (me *VpcService) DescribeEipNetworkAccountType(ctx context.Context) (networkAccountType *string, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewDescribeNetworkAccountTypeRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DescribeNetworkAccountType(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	networkAccountType = response.Response.NetworkAccountType

	return
}

func (me *VpcService) DescribeVpcBandwidthPackageQuota(ctx context.Context) (bandwidthPackageQuota []*vpc.Quota, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewDescribeBandwidthPackageQuotaRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DescribeBandwidthPackageQuota(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	bandwidthPackageQuota = append(bandwidthPackageQuota, response.Response.QuotaSet...)

	return
}

func (me *VpcService) DescribeVpcBandwidthPackageBillUsageByFilter(ctx context.Context, param map[string]interface{}) (bandwidthPackageBillUsage []*vpc.BandwidthPackageBillBandwidth, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewDescribeBandwidthPackageBillUsageRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "BandwidthPackageId" {
			request.BandwidthPackageId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DescribeBandwidthPackageBillUsage(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	bandwidthPackageBillUsage = append(bandwidthPackageBillUsage, response.Response.BandwidthPackageBillBandwidthSet...)

	return
}

func (me *VpcService) DescribeVpcTrafficPackageById(ctx context.Context, trafficPackageId string) (TrafficPackage *vpc.TrafficPackage, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDescribeTrafficPackagesRequest()
	request.TrafficPackageIds = []*string{&trafficPackageId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DescribeTrafficPackages(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.TrafficPackageSet) < 1 {
		return
	}

	TrafficPackage = response.Response.TrafficPackageSet[0]
	return
}

func (me *VpcService) DeleteVpcTrafficPackageById(ctx context.Context, trafficPackageId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDeleteTrafficPackagesRequest()
	request.TrafficPackageIds = []*string{&trafficPackageId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DeleteTrafficPackages(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *VpcService) DescribeVpcSnapshotPoliciesById(ctx context.Context, snapshotPolicyId string) (snapshotPolices []*vpc.SnapshotPolicy, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDescribeSnapshotPoliciesRequest()
	request.SnapshotPolicyIds = []*string{&snapshotPolicyId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DescribeSnapshotPolicies(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	snapshotPolices = response.Response.SnapshotPolicySet
	return
}

func (me *VpcService) DeleteVpcSnapshotPoliciesById(ctx context.Context, snapshotPolicyId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDeleteSnapshotPoliciesRequest()
	request.SnapshotPolicyIds = []*string{&snapshotPolicyId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DeleteSnapshotPolicies(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *VpcService) DescribeVpcSnapshotPolicyAttachmentById(ctx context.Context, snapshotPolicyId string) (snapshotPolicyAttachment []*vpc.SnapshotInstance, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDescribeSnapshotAttachedInstancesRequest()
	request.SnapshotPolicyId = &snapshotPolicyId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DescribeSnapshotAttachedInstances(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.InstanceSet) < 1 {
		return
	}

	snapshotPolicyAttachment = response.Response.InstanceSet
	return
}

func (me *VpcService) DeleteVpcSnapshotPolicyAttachmentById(ctx context.Context, snapshotPolicyId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDetachSnapshotInstancesRequest()
	request.SnapshotPolicyId = &snapshotPolicyId

	snapshotInstace, err := me.DescribeVpcSnapshotPolicyAttachmentById(ctx, snapshotPolicyId)
	if err != nil {
		errRet = err
		return
	}
	request.Instances = snapshotInstace

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DetachSnapshotInstances(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *VpcService) DescribeVpcNetDetectById(ctx context.Context, netDetectId string) (netDetect *vpc.NetDetect, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDescribeNetDetectsRequest()
	request.NetDetectIds = []*string{&netDetectId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DescribeNetDetects(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.NetDetectSet) < 1 {
		return
	}

	netDetect = response.Response.NetDetectSet[0]
	return
}

func (me *VpcService) DeleteVpcNetDetectById(ctx context.Context, netDetectId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDeleteNetDetectRequest()
	request.NetDetectId = &netDetectId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DeleteNetDetect(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *VpcService) DescribeVpcClassicLinkAttachmentById(ctx context.Context, vpcId string, instanceId string) (classicLinkAttachment *vpc.ClassicLinkInstance, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDescribeClassicLinkInstancesRequest()
	filter := vpc.FilterObject{
		Name:   helper.String("vpc-id"),
		Values: []*string{&vpcId},
	}
	request.Filters = append(request.Filters, &filter)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	instances := make([]*vpc.ClassicLinkInstance, 0)
	for {
		request.Offset = helper.Int64ToStrPoint(offset)
		request.Limit = helper.Int64ToStrPoint(limit)
		response, err := me.client.UseVpcClient().DescribeClassicLinkInstances(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.ClassicLinkInstanceSet) < 1 {
			break
		}
		instances = append(instances, response.Response.ClassicLinkInstanceSet...)
		if len(response.Response.ClassicLinkInstanceSet) < int(limit) {
			break
		}
		offset += limit
	}

	if len(instances) < 1 {
		return
	}

	for _, instance := range instances {
		if *instance.InstanceId == instanceId {
			classicLinkAttachment = instance
		}
	}

	return
}

func (me *VpcService) DeleteVpcClassicLinkAttachmentById(ctx context.Context, vpcId string, instanceId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDetachClassicLinkVpcRequest()
	request.VpcId = &vpcId
	request.InstanceIds = []*string{&instanceId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DetachClassicLinkVpc(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *VpcService) DescribeVpcDhcpIpById(ctx context.Context, dhcpIpId string) (dhcpIp *vpc.DhcpIp, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDescribeDhcpIpsRequest()
	request.DhcpIpIds = []*string{&dhcpIpId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DescribeDhcpIps(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.DhcpIpSet) < 1 {
		return
	}

	dhcpIp = response.Response.DhcpIpSet[0]
	return
}

func (me *VpcService) DeleteVpcDhcpIpById(ctx context.Context, dhcpIpId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDeleteDhcpIpRequest()
	request.DhcpIpId = &dhcpIpId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DeleteDhcpIp(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *VpcService) DescribeVpcDhcpAssociateAddressById(ctx context.Context, dhcpIpId string, addressIp string) (dhcpAssociateAddress *vpc.DhcpIp, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDescribeDhcpIpsRequest()
	request.DhcpIpIds = []*string{&dhcpIpId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DescribeDhcpIps(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.DhcpIpSet) < 1 {
		return
	}

	dhcpIp := response.Response.DhcpIpSet[0]
	if *dhcpIp.AddressIp != addressIp {
		return
	}
	dhcpAssociateAddress = dhcpIp

	return
}

func (me *VpcService) DeleteVpcDhcpAssociateAddressById(ctx context.Context, dhcpIpId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDisassociateDhcpIpWithAddressIpRequest()
	request.DhcpIpId = &dhcpIpId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DisassociateDhcpIpWithAddressIp(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *VpcService) DescribeVpcById(ctx context.Context, vpcId string) (instance *vpc.Vpc, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDescribeVpcsRequest()
	request.VpcIds = []*string{&vpcId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	instances := make([]*vpc.Vpc, 0)
	for {
		request.Offset = helper.Int64ToStrPoint(offset)
		request.Limit = helper.Int64ToStrPoint(limit)
		response, err := me.client.UseVpcClient().DescribeVpcs(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.VpcSet) < 1 {
			break
		}
		instances = append(instances, response.Response.VpcSet...)
		if len(response.Response.VpcSet) < int(limit) {
			break
		}

		offset += limit
	}

	if len(instances) < 1 {
		return
	}
	instance = instances[0]
	return
}

func (me *VpcService) DeleteVpcIpv6CidrBlockById(ctx context.Context, vpcId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewUnassignIpv6CidrBlockRequest()
	request.VpcId = &vpcId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().UnassignIpv6CidrBlock(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *VpcService) DescribeSubnetById(ctx context.Context, subnetId string) (instance *vpc.Subnet, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDescribeSubnetsRequest()
	request.SubnetIds = []*string{&subnetId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	instances := make([]*vpc.Subnet, 0)
	for {
		request.Offset = helper.Int64ToStrPoint(offset)
		request.Limit = helper.Int64ToStrPoint(limit)
		response, err := me.client.UseVpcClient().DescribeSubnets(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.SubnetSet) < 1 {
			break
		}
		instances = append(instances, response.Response.SubnetSet...)
		if len(response.Response.SubnetSet) < int(limit) {
			break
		}

		offset += limit
	}

	if len(instances) < 1 {
		return
	}
	instance = instances[0]
	return
}

func (me *VpcService) DeleteVpcIpv6SubnetCidrBlockById(ctx context.Context, vpcId string, subnetId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewUnassignIpv6SubnetCidrBlockRequest()
	request.VpcId = &vpcId

	ipv6SubnetCidrBlock := vpc.Ipv6SubnetCidrBlock{}
	ipv6SubnetCidrBlock.SubnetId = &subnetId
	request.Ipv6SubnetCidrBlocks = append(request.Ipv6SubnetCidrBlocks, &ipv6SubnetCidrBlock)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().UnassignIpv6SubnetCidrBlock(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *VpcService) DescribeVpcIpv6EniAddressById(ctx context.Context, vpcId string, ipv6Address string) (ipv6EniAddress *vpc.VpcIpv6Address, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDescribeVpcIpv6AddressesRequest()
	request.VpcId = &vpcId
	request.Ipv6Addresses = []*string{&ipv6Address}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DescribeVpcIpv6Addresses(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Ipv6AddressSet) < 1 {
		return
	}

	ipv6EniAddress = response.Response.Ipv6AddressSet[0]
	return
}

func (me *VpcService) DeleteVpcIpv6EniAddressById(ctx context.Context, networkInterfaceId string, ipv6Address string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewUnassignIpv6AddressesRequest()
	request.NetworkInterfaceId = &networkInterfaceId
	address := vpc.Ipv6Address{}
	address.Address = &ipv6Address
	request.Ipv6Addresses = append(request.Ipv6Addresses, &address)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().UnassignIpv6Addresses(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *VpcService) DescribeVpcLocalGatewayById(ctx context.Context, localGatewayId string) (localGateway *vpc.LocalGateway, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDescribeLocalGatewayRequest()

	filter := vpc.Filter{
		Name:   helper.String("local-gateway-id"),
		Values: []*string{&localGatewayId},
	}

	request.Filters = append(request.Filters, &filter)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DescribeLocalGateway(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.LocalGatewaySet) < 1 {
		return
	}

	localGateway = response.Response.LocalGatewaySet[0]
	return
}

func (me *VpcService) DeleteVpcLocalGatewayById(ctx context.Context, cdcId string, localGatewayId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDeleteLocalGatewayRequest()
	request.CdcId = &cdcId
	request.LocalGatewayId = &localGatewayId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DeleteLocalGateway(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *VpcService) DescribeVpcAccountAttributes(ctx context.Context) (accountAttributes []*vpc.AccountAttribute, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewDescribeAccountAttributesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DescribeAccountAttributes(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	accountAttributes = response.Response.AccountAttributeSet

	return
}

func (me *VpcService) DescribeVpcClassicLinkInstancesByFilter(ctx context.Context, param map[string]interface{}) (classicLinkInstances []*vpc.ClassicLinkInstance, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewDescribeClassicLinkInstancesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Filters" {
			request.Filters = v.([]*vpc.FilterObject)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = helper.Int64ToStrPoint(offset)
		request.Limit = helper.Int64ToStrPoint(limit)
		response, err := me.client.UseVpcClient().DescribeClassicLinkInstances(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.ClassicLinkInstanceSet) < 1 {
			break
		}
		classicLinkInstances = append(classicLinkInstances, response.Response.ClassicLinkInstanceSet...)
		if len(response.Response.ClassicLinkInstanceSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *VpcService) DescribeVpcGatewayFlowMonitorDetailByFilter(ctx context.Context, param map[string]interface{}) (GatewayFlowMonitorDetail []*vpc.GatewayFlowMonitorDetail, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewDescribeGatewayFlowMonitorDetailRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "TimePoint" {
			request.TimePoint = v.(*string)
		}
		if k == "VpnId" {
			request.VpnId = v.(*string)
		}
		if k == "DirectConnectGatewayId" {
			request.DirectConnectGatewayId = v.(*string)
		}
		if k == "PeeringConnectionId" {
			request.PeeringConnectionId = v.(*string)
		}
		if k == "NatId" {
			request.NatId = v.(*string)
		}
		if k == "OrderField" {
			request.OrderField = v.(*string)
		}
		if k == "OrderDirection" {
			request.OrderDirection = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseVpcClient().DescribeGatewayFlowMonitorDetail(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.GatewayFlowMonitorDetailSet) < 1 {
			break
		}
		GatewayFlowMonitorDetail = append(GatewayFlowMonitorDetail, response.Response.GatewayFlowMonitorDetailSet...)
		if len(response.Response.GatewayFlowMonitorDetailSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *VpcService) DescribeVpcGatewayFlowQosByFilter(ctx context.Context, param map[string]interface{}) (GatewayFlowQos []*vpc.GatewayQos, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewDescribeGatewayFlowQosRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "GatewayId" {
			request.GatewayId = v.(*string)
		}
		if k == "IpAddresses" {
			request.IpAddresses = v.([]*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseVpcClient().DescribeGatewayFlowQos(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.GatewayQosSet) < 1 {
			break
		}
		GatewayFlowQos = append(GatewayFlowQos, response.Response.GatewayQosSet...)
		if len(response.Response.GatewayQosSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *VpcService) DescribeVpcCvmInstancesByFilter(ctx context.Context, param map[string]interface{}) (CvmInstances []*vpc.CvmInstance, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewDescribeVpcInstancesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Filters" {
			request.Filters = v.([]*vpc.Filter)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseVpcClient().DescribeVpcInstances(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.InstanceSet) < 1 {
			break
		}
		CvmInstances = append(CvmInstances, response.Response.InstanceSet...)
		if len(response.Response.InstanceSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *VpcService) DescribeVpcNetDetectStatesByFilter(ctx context.Context, param map[string]interface{}) (NetDetectStates []*vpc.NetDetectState, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewDescribeNetDetectStatesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "NetDetectIds" {
			request.NetDetectIds = v.([]*string)
		}
		if k == "Filters" {
			request.Filters = v.([]*vpc.Filter)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseVpcClient().DescribeNetDetectStates(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.NetDetectStateSet) < 1 {
			break
		}
		NetDetectStates = append(NetDetectStates, response.Response.NetDetectStateSet...)
		if len(response.Response.NetDetectStateSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *VpcService) DescribeVpcNetworkInterfaceLimit(ctx context.Context, param map[string]interface{}) (networkInterfaceLimit *vpc.DescribeNetworkInterfaceLimitResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewDescribeNetworkInterfaceLimitRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DescribeNetworkInterfaceLimit(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	networkInterfaceLimit = response.Response

	return
}

func (me *VpcService) DescribeVpcPrivateIpAddresses(ctx context.Context, param map[string]interface{}) (PrivateIpAddresses []*vpc.VpcPrivateIpAddress, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewDescribeVpcPrivateIpAddressesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "VpcId" {
			request.VpcId = v.(*string)
		}
		if k == "PrivateIpAddresses" {
			request.PrivateIpAddresses = v.([]*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DescribeVpcPrivateIpAddresses(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	PrivateIpAddresses = response.Response.VpcPrivateIpAddressSet

	return
}

func (me *VpcService) DescribeVpcProductQuota(ctx context.Context, param map[string]interface{}) (ProductQuota []*vpc.ProductQuota, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewDescribeProductQuotaRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Product" {
			request.Product = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DescribeProductQuota(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	ProductQuota = response.Response.ProductQuotaSet

	return
}

func (me *VpcService) DescribeVpcResourceDashboard(ctx context.Context, param map[string]interface{}) (ResourceDashboard []*vpc.ResourceDashboard, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewDescribeVpcResourceDashboardRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "VpcIds" {
			request.VpcIds = v.([]*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DescribeVpcResourceDashboard(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	ResourceDashboard = response.Response.ResourceDashboardSet

	return
}

func (me *VpcService) DescribeVpcRouteConflicts(ctx context.Context, param map[string]interface{}) (routeConflicts []*vpc.RouteConflict, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewDescribeRouteConflictsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "RouteTableId" {
			request.RouteTableId = v.(*string)
		}
		if k == "DestinationCidrBlocks" {
			request.DestinationCidrBlocks = v.([]*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DescribeRouteConflicts(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	routeConflicts = response.Response.RouteConflictSet

	return
}

func (me *VpcService) DescribeVpcSecurityGroupLimits(ctx context.Context, param map[string]interface{}) (securityGroupLimit *vpc.SecurityGroupLimitSet, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewDescribeSecurityGroupLimitsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DescribeSecurityGroupLimits(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	securityGroupLimit = response.Response.SecurityGroupLimitSet

	return
}

func (me *VpcService) DescribeVpcSecurityGroupReferences(ctx context.Context, param map[string]interface{}) (securityGroupReferences []*vpc.ReferredSecurityGroup, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewDescribeSecurityGroupReferencesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "SecurityGroupIds" {
			request.SecurityGroupIds = v.([]*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DescribeSecurityGroupReferences(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	securityGroupReferences = response.Response.ReferredSecurityGroupSet

	return
}

func (me *VpcService) DescribeVpcSgSnapshotFileContent(ctx context.Context, param map[string]interface{}) (sgSnapshotFileContent *vpc.DescribeSgSnapshotFileContentResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewDescribeSgSnapshotFileContentRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "SnapshotPolicyId" {
			request.SnapshotPolicyId = v.(*string)
		}
		if k == "SnapshotFileId" {
			request.SnapshotFileId = v.(*string)
		}
		if k == "SecurityGroupId" {
			request.SecurityGroupId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DescribeSgSnapshotFileContent(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	sgSnapshotFileContent = response.Response

	return
}

func (me *VpcService) DescribeVpcSnapshotFilesByFilter(ctx context.Context, param map[string]interface{}) (SnapshotFiles []*vpc.SnapshotFileInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewDescribeSnapshotFilesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "BusinessType" {
			request.BusinessType = v.(*string)
		}
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
		if k == "StartDate" {
			request.StartDate = v.(*string)
		}
		if k == "EndDate" {
			request.EndDate = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseVpcClient().DescribeSnapshotFiles(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.SnapshotFileSet) < 1 {
			break
		}
		SnapshotFiles = append(SnapshotFiles, response.Response.SnapshotFileSet...)
		if len(response.Response.SnapshotFileSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *VpcService) DescribeVpcSubnetResourceDashboardByFilter(ctx context.Context, param map[string]interface{}) (subnetResourceDashboard []*vpc.ResourceStatistics, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewDescribeSubnetResourceDashboardRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "SubnetIds" {
			request.SubnetIds = v.([]*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DescribeSubnetResourceDashboard(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	subnetResourceDashboard = response.Response.ResourceStatisticsSet

	return
}

func (me *VpcService) DescribeVpcTemplateLimits(ctx context.Context) (templateLimit *vpc.TemplateLimit, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewDescribeTemplateLimitsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DescribeTemplateLimits(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	templateLimit = response.Response.TemplateLimit

	return
}

func (me *VpcService) DescribeVpcUsedIpAddressByFilter(ctx context.Context, param map[string]interface{}) (UsedIpAddress []*vpc.IpAddressStates, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewDescribeUsedIpAddressRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "VpcId" {
			request.VpcId = v.(*string)
		}
		if k == "SubnetId" {
			request.SubnetId = v.(*string)
		}
		if k == "IpAddresses" {
			request.IpAddresses = v.([]*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseVpcClient().DescribeUsedIpAddress(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.IpAddressStates) < 1 {
			break
		}
		UsedIpAddress = append(UsedIpAddress, response.Response.IpAddressStates...)
		if len(response.Response.IpAddressStates) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *VpcService) DescribeVpcLimitsByFilter(ctx context.Context, param map[string]interface{}) (limits []*vpc.VpcLimit, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewDescribeVpcLimitsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "LimitTypes" {
			request.LimitTypes = v.([]*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DescribeVpcLimits(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	limits = response.Response.VpcLimitSet

	return
}

func (me *VpcService) DescribeVpcNetworkAclQuintupleById(ctx context.Context, networkAclId string) (networkAclQuintuples []*vpc.NetworkAclQuintupleEntry, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDescribeNetworkAclQuintupleEntriesRequest()
	request.NetworkAclId = &networkAclId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DescribeNetworkAclQuintupleEntries(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.NetworkAclQuintupleSet) < 1 {
		return
	}

	networkAclQuintuples = response.Response.NetworkAclQuintupleSet
	return
}

func (me *VpcService) DeleteVpcNetworkAclQuintupleById(ctx context.Context, networkAclId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDeleteNetworkAclQuintupleEntriesRequest()
	request.NetworkAclId = &networkAclId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DeleteNetworkAclQuintupleEntries(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *VpcService) DeleteVpcEniSgAttachmentById(ctx context.Context, networkInterfaceId string, securityGroupIds []string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDisassociateNetworkInterfaceSecurityGroupsRequest()
	request.NetworkInterfaceIds = []*string{&networkInterfaceId}
	request.SecurityGroupIds = common.StringPtrs(securityGroupIds)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DisassociateNetworkInterfaceSecurityGroups(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *VpcService) DescribeVpcNetDetectStateCheck(ctx context.Context, param map[string]interface{}) (netDetectStateCheck []*vpc.NetDetectIpState, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewCheckNetDetectStateRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "DetectDestinationIp" {
			request.DetectDestinationIp = v.([]*string)
		}
		if k == "NextHopType" {
			request.NextHopType = v.(*string)
		}
		if k == "NextHopDestination" {
			request.NextHopDestination = v.(*string)
		}
		if k == "NetDetectId" {
			request.NetDetectId = v.(*string)
		}
		if k == "VpcId" {
			request.VpcId = v.(*string)
		}
		if k == "SubnetId" {
			request.SubnetId = v.(*string)
		}
		if k == "NetDetectName" {
			request.NetDetectName = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().CheckNetDetectState(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	netDetectStateCheck = response.Response.NetDetectIpStateSet

	return
}

func (me *VpcService) DescribeVpcNotifyRoutesById(ctx context.Context, routeTableId string, routeItemId string) (notifyRoute *vpc.Route, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDescribeRouteTablesRequest()
	request.RouteTableIds = []*string{&routeTableId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().DescribeRouteTables(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.RouteTableSet) < 1 {
		return
	}

	for _, routeTable := range response.Response.RouteTableSet {
		for _, route := range routeTable.RouteSet {
			if *route.RouteItemId == routeItemId {
				notifyRoute = route
				break
			}
		}
	}
	return
}

func (me *VpcService) DeleteVpcNotifyRoutesById(ctx context.Context, routeTableId string, routeItemId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewWithdrawNotifyRoutesRequest()
	request.RouteTableId = &routeTableId
	request.RouteItemIds = []*string{&routeItemId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().WithdrawNotifyRoutes(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *VpcService) DescribeVpnDefaultHealthCheckIp(ctx context.Context, param map[string]interface{}) (defaultHealthCheck *vpc.GenerateVpnConnectionDefaultHealthCheckIpResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewGenerateVpnConnectionDefaultHealthCheckIpRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "VpnGatewayId" {
			request.VpnGatewayId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVpcClient().GenerateVpnConnectionDefaultHealthCheckIp(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	defaultHealthCheck = response.Response

	return
}
