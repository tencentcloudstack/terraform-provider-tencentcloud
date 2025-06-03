package ccn

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"

	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

// Ccn basic information
type CcnBasicInfo struct {
	ccnId             string
	name              string
	description       string
	state             string
	qos               string
	chargeType        string
	bandWithLimitType string
	instanceCount     int64
	createTime        string
	ecmpFlag          bool
	overlapFlag       bool
}

type CcnInstanceBind struct {
	instanceId   string
	instanceType string
}

type CcnInstanceBindResponseData struct {
	InstanceBindSet []struct {
		AliasInstanceId  string `json:"AliasInstanceId"`
		AliasType        string `json:"AliasType"`
		CcnId            string `json:"CcnId"`
		InstanceBindTime string `json:"InstanceBindTime"`
		InstanceId       string `json:"InstanceId"`
		InstanceName     string `json:"InstanceName"`
		InstanceRegion   string `json:"InstanceRegion"`
		InstanceType     string `json:"InstanceType"`
		InstanceUin      string `json:"InstanceUin"`
		RouteTableId     string `json:"RouteTableId"`
		State            string `json:"State"`
	} `json:"InstanceBindSet"`
	RequestId  string `json:"RequestId"`
	TotalCount int    `json:"TotalCount"`
}

type CcnInstanceBindApiResponse struct {
	Response CcnInstanceBindResponseData `json:"Response"`
}

func (info CcnBasicInfo) CcnId() string {
	return info.ccnId
}

func (info CcnBasicInfo) Name() string {
	return info.name
}

func (info CcnBasicInfo) BandWithLimitType() string {
	return info.bandWithLimitType
}

func (info CcnBasicInfo) CreateTime() string {
	return info.createTime
}

type CcnAttachedInstanceInfo struct {
	ccnUin         string
	ccnId          string
	instanceType   string
	instanceRegion string
	instanceId     string
	state          string
	attachedTime   string
	cidrBlock      []string
	description    string
	routeTableId   string
}

type CcnBandwidthLimit struct {
	region string
	limit  int64
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

func (me *VpcService) DescribeCcn(ctx context.Context, ccnId string) (info CcnBasicInfo, has int, errRet error) {
	infos, err := me.DescribeCcns(ctx, ccnId, "")
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

func (me *VpcService) DescribeCcns(ctx context.Context, ccnId, name string) (infos []CcnBasicInfo, errRet error) {

	logId := tccommon.GetLogId(ctx)
	request := vpc.NewDescribeCcnsRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	infos = make([]CcnBasicInfo, 0, 100)

	var offset uint64 = 0
	var limit uint64 = 100
	var total = -1
	var has = map[string]bool{}

	var filters []*vpc.Filter
	if ccnId != "" {
		filters = me.fillFilter(filters, "ccn-id", ccnId)
	}
	if name != "" {
		filters = me.fillFilter(filters, "ccn-name", name)
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
	response, err := me.client.UseVpcClient().DescribeCcns(request)

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

	log.Printf("[DEBUG]%s api[%s] , request body [%s], response body[%s]\n",
		logId,
		request.GetAction(),
		request.ToJsonString(),
		response.ToJsonString())

	if total < 0 {
		total = int(*response.Response.TotalCount)
	}

	if len(response.Response.CcnSet) > 0 {
		offset += limit
	} else {
		//get empty ,we're done
		return
	}
	for _, item := range response.Response.CcnSet {
		var basicInfo CcnBasicInfo

		basicInfo.ccnId = *item.CcnId
		basicInfo.name = *item.CcnName
		basicInfo.createTime = *item.CreateTime
		basicInfo.description = *item.CcnDescription
		basicInfo.instanceCount = int64(*item.InstanceCount)
		basicInfo.qos = *item.QosLevel
		basicInfo.state = *item.State
		basicInfo.chargeType = *item.InstanceChargeType
		basicInfo.bandWithLimitType = *item.BandwidthLimitType
		basicInfo.ecmpFlag = *item.RouteECMPFlag
		basicInfo.overlapFlag = *item.RouteOverlapFlag

		if has[basicInfo.ccnId] {
			errRet = fmt.Errorf("get repeated ccn_id[%s] when doing DescribeCcns", basicInfo.ccnId)
			return
		}
		has[basicInfo.ccnId] = true
		infos = append(infos, basicInfo)
	}
	goto getMoreData

}

func (me *VpcService) DescribeCcnRegionBandwidthLimits(ctx context.Context, ccnId string) (infos []CcnBandwidthLimit, errRet error) {

	logId := tccommon.GetLogId(ctx)
	request := vpc.NewDescribeCcnRegionBandwidthLimitsRequest()

	infos = make([]CcnBandwidthLimit, 0, 100)

	request.CcnId = &ccnId
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().DescribeCcnRegionBandwidthLimits(request)

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
	log.Printf("[DEBUG]%s api[%s] , request body [%s], response body[%s]\n",
		logId,
		request.GetAction(),
		request.ToJsonString(),
		response.ToJsonString())

	for _, item := range response.Response.CcnRegionBandwidthLimitSet {

		var ccnBandwidthLimit CcnBandwidthLimit
		ccnBandwidthLimit.region = *item.Region
		ccnBandwidthLimit.limit = int64(*item.BandwidthLimit)
		infos = append(infos, ccnBandwidthLimit)
	}
	return
}

func (me *VpcService) CreateCcn(ctx context.Context, name, description,
	qos, chargeType, bandWithLimitType string) (basicInfo CcnBasicInfo, errRet error) {

	logId := tccommon.GetLogId(ctx)
	request := vpc.NewCreateCcnRequest()

	request.CcnName = &name
	request.CcnDescription = &description
	request.QosLevel = &qos
	request.InstanceChargeType = &chargeType
	request.BandwidthLimitType = &bandWithLimitType
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().CreateCcn(request)

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
	log.Printf("[DEBUG]%s api[%s] , request body [%s], response body[%s]\n",
		logId,
		request.GetAction(),
		request.ToJsonString(),
		response.ToJsonString())

	if response.Response.Ccn == nil || response.Response.Ccn.CcnId == nil || *response.Response.Ccn.CcnId == "" {
		errRet = fmt.Errorf("CreateCcn return empty response.Response.Ccn ")
		return
	}

	item := response.Response.Ccn
	basicInfo.ccnId = *item.CcnId
	basicInfo.name = *item.CcnName
	basicInfo.createTime = *item.CreateTime

	basicInfo.description = *item.CcnDescription
	basicInfo.instanceCount = int64(*item.InstanceCount)
	basicInfo.qos = *item.QosLevel
	basicInfo.state = *item.State
	return
}

func (me *VpcService) DeleteCcn(ctx context.Context, ccnId string) (errRet error) {

	logId := tccommon.GetLogId(ctx)
	request := vpc.NewDeleteCcnRequest()
	request.CcnId = &ccnId
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().DeleteCcn(request)

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
	log.Printf("[DEBUG]%s api[%s] , request body [%s], response body[%s]\n",
		logId,
		request.GetAction(),
		request.ToJsonString(),
		response.ToJsonString())
	return
}

func (me *VpcService) ModifyCcnAttribute(ctx context.Context, ccnId, name, description string, ecmpFlag, overlapFlag bool) (errRet error) {

	logId := tccommon.GetLogId(ctx)
	request := vpc.NewModifyCcnAttributeRequest()
	request.CcnId = &ccnId

	if name != "" {
		request.CcnName = &name
	}

	if description != "" {
		request.CcnDescription = &description
	}

	request.RouteECMPFlag = &ecmpFlag
	request.RouteOverlapFlag = &overlapFlag

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().ModifyCcnAttribute(request)

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
	log.Printf("[DEBUG]%s api[%s] , request body [%s], response body[%s]\n",
		logId,
		request.GetAction(),
		request.ToJsonString(),
		response.ToJsonString())
	return
}

func (me *VpcService) DescribeCcnAttachedInstance(ctx context.Context, ccnId,
	instanceRegion, instanceType, instanceId string) (info CcnAttachedInstanceInfo, has int, errRet error) {

	infos, err := me.DescribeCcnAttachedInstances(ctx, ccnId)

	if err != nil {
		errRet = err
		return
	}

	for _, item := range infos {
		if item.instanceId == instanceId &&
			item.instanceRegion == instanceRegion &&
			strings.EqualFold(item.instanceType, instanceType) {
			has = 1
			info = item
			return
		}
	}
	return
}

func (me *VpcService) DescribeCcnAttachedInstanceByFilter(ctx context.Context, ccnId, instanceType, instanceRegion, instanceId string) (info *vpc.CcnAttachedInstance, errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = vpc.NewDescribeCcnAttachedInstancesRequest()
		response = vpc.NewDescribeCcnAttachedInstancesResponse()
		result   []*vpc.CcnAttachedInstance
		limit    uint64 = 20
		offset   uint64 = 0
	)

	request.CcnId = &ccnId
	request.Filters = []*vpc.Filter{
		{
			Name:   helper.String("instance-type"),
			Values: helper.Strings([]string{instanceType}),
		},

		{
			Name:   helper.String("instance-region"),
			Values: helper.Strings([]string{instanceRegion}),
		},

		{
			Name:   helper.String("instance-id"),
			Values: helper.Strings([]string{instanceId}),
		},
	}

	for {
		request.Limit = &limit
		request.Offset = &offset
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseVpcClient().DescribeCcnAttachedInstancesWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe attach ccn instance failed, Response is nil."))
			}

			response = result
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		if response == nil || len(response.Response.InstanceSet) < 1 {
			break
		}

		result = append(result, response.Response.InstanceSet...)
		if len(response.Response.InstanceSet) < int(limit) {
			break
		}

		offset += limit
	}

	if len(result) != 1 {
		return
	}

	info = result[0]
	return
}

func (me *VpcService) DescribeCcnAttachedInstances(ctx context.Context, ccnId string) (infos []CcnAttachedInstanceInfo, errRet error) {

	logId := tccommon.GetLogId(ctx)
	var (
		request  = vpc.NewDescribeCcnAttachedInstancesRequest()
		response = vpc.NewDescribeCcnAttachedInstancesResponse()
		result   []*vpc.CcnAttachedInstance
		err      error
		limit    uint64 = 20
		offset   uint64 = 0
	)

	request.CcnId = &ccnId

	ratelimit.Check(request.GetAction())

	for {
		request.Limit = &limit
		request.Offset = &offset
		response, err = me.client.UseVpcClient().DescribeCcnAttachedInstances(request)
		if err != nil {
			errRet = err
			return
		}

		if response == nil || len(response.Response.InstanceSet) < 1 {
			break
		}
		result = append(result, response.Response.InstanceSet...)
		if len(response.Response.InstanceSet) < int(limit) {
			break
		}
		offset += limit
	}
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
	log.Printf("[DEBUG]%s api[%s] , request body [%s], response body[%s]\n",
		logId,
		request.GetAction(),
		request.ToJsonString(),
		response.ToJsonString())

	infos = make([]CcnAttachedInstanceInfo, 0, len(result))

	for _, item := range result {

		var info CcnAttachedInstanceInfo

		info.attachedTime = *item.AttachedTime
		info.cidrBlock = make([]string, 0, len(item.CidrBlock))

		for _, v := range item.CidrBlock {
			info.cidrBlock = append(info.cidrBlock, *v)
		}

		info.ccnId = ccnId
		info.instanceId = *item.InstanceId
		info.instanceRegion = *item.InstanceRegion
		info.instanceType = *item.InstanceType
		info.state = *item.State
		info.description = *item.Description
		info.routeTableId = *item.RouteTableId
		info.ccnUin = *item.CcnUin
		infos = append(infos, info)
	}
	return
}

func (me *VpcService) DescribeVpcDescribeCcnRoutesByFilter(ctx context.Context, paramMap map[string]interface{}) (RouteSet []*vpc.CcnRoute, errRet error) {
	var (
		logId           = tccommon.GetLogId(ctx)
		request         = vpc.NewDescribeCcnRoutesRequest()
		response        = vpc.NewDescribeCcnRoutesResponse()
		limit    uint64 = 20
		offset   uint64 = 0
		err      error
	)

	for k, v := range paramMap {
		if k == "CcnId" {
			request.CcnId = v.(*string)
		}

		if k == "Filters" {
			request.Filters = v.([]*vpc.Filter)
		}
	}

	ratelimit.Check(request.GetAction())

	for {
		request.Limit = &limit
		request.Offset = &offset
		response, err = me.client.UseVpcClient().DescribeCcnRoutes(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]", logId, request.GetAction(), request.ToJsonString(), err)
			errRet = err
			return
		}

		if response == nil || len(response.Response.RouteSet) < 1 {
			break
		}

		RouteSet = append(RouteSet, response.Response.RouteSet...)
		if len(response.Response.RouteSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *VpcService) DescribeCcnAttachmentsByInstance(ctx context.Context, instanceType string, instanceId string, instanceRegion string) (infos []vpc.CcnAttachedInstance, errRet error) {

	logId := tccommon.GetLogId(ctx)
	request := vpc.NewDescribeCcnAttachedInstancesRequest()
	request.Filters = make([]*vpc.Filter, 0, 3)
	request.Filters = append(request.Filters, &vpc.Filter{Name: helper.String("instance-type"), Values: []*string{&instanceType}})
	request.Filters = append(request.Filters, &vpc.Filter{Name: helper.String("instance-id"), Values: []*string{&instanceId}})
	request.Filters = append(request.Filters, &vpc.Filter{Name: helper.String("instance-region"), Values: []*string{&instanceRegion}})

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().DescribeCcnAttachedInstances(request)

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
	log.Printf("[DEBUG]%s api[%s] , request body [%s], response body[%s]\n",
		logId,
		request.GetAction(),
		request.ToJsonString(),
		response.ToJsonString())

	infos = make([]vpc.CcnAttachedInstance, 0, len(response.Response.InstanceSet))

	for _, item := range response.Response.InstanceSet {
		infos = append(infos, *item)
	}
	return
}

func (me *VpcService) AttachCcnInstances(ctx context.Context, ccnId, instanceRegion, instanceType, instanceId, ccnUin, description, routeTableId string) (errRet error) {

	logId := tccommon.GetLogId(ctx)
	request := vpc.NewAttachCcnInstancesRequest()
	request.CcnId = &ccnId

	if ccnUin != "" {
		request.CcnUin = &ccnUin
	}

	var ccnInstance vpc.CcnInstance
	ccnInstance.InstanceId = &instanceId
	ccnInstance.InstanceRegion = &instanceRegion
	ccnInstance.InstanceType = &instanceType
	if description != "" {
		ccnInstance.Description = &description
	}

	if routeTableId != "" {
		ccnInstance.RouteTableId = &routeTableId
	}

	request.Instances = []*vpc.CcnInstance{&ccnInstance}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().AttachCcnInstances(request)

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
	log.Printf("[DEBUG]%s api[%s] , request body [%s], response body[%s]\n",
		logId,
		request.GetAction(),
		request.ToJsonString(),
		response.ToJsonString())
	return
}

func (me *VpcService) DetachCcnInstances(ctx context.Context, ccnId, instanceRegion, instanceType, instanceId string) (errRet error) {

	logId := tccommon.GetLogId(ctx)
	request := vpc.NewDetachCcnInstancesRequest()
	request.CcnId = &ccnId

	var ccnInstance vpc.CcnInstance
	ccnInstance.InstanceId = &instanceId
	ccnInstance.InstanceRegion = &instanceRegion
	ccnInstance.InstanceType = &instanceType

	request.Instances = []*vpc.CcnInstance{&ccnInstance}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().DetachCcnInstances(request)

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
	log.Printf("[DEBUG]%s api[%s] , request body [%s], response body[%s]\n",
		logId,
		request.GetAction(),
		request.ToJsonString(),
		response.ToJsonString())
	return

}

func (me *VpcService) DescribeCcnRegionBandwidthLimit(ctx context.Context, ccnId,
	region string) (bandwidth int64, errRet error) {

	infos, err := me.DescribeCcnRegionBandwidthLimits(ctx, ccnId)
	if err != nil {
		errRet = err
		return
	}
	for _, v := range infos {
		if v.region == region {
			bandwidth = v.limit
			break
		}
	}
	return
}

func (me *VpcService) GetCcnRegionBandwidthLimit(ctx context.Context, ccnId,
	region, dstRegion, limitType string) (int64, error) {
	infos, err := me.GetCcnRegionBandwidthLimits(ctx, ccnId)
	if err != nil {
		return 0, err
	}
	for _, v := range infos {
		if v.Region != nil {
			switch limitType {
			case OuterRegionLimit:
				if *v.Region == region {
					return int64(*v.BandwidthLimit), nil
				}
			case InterRegionLimit:
				if v.DstRegion != nil && *v.DstRegion == dstRegion && *v.Region == region {
					return int64(*v.BandwidthLimit), nil
				}
			default:
				return 0, fmt.Errorf("unknown type of band with limit type")
			}
		}
	}
	return 0, nil
}

func (me *VpcService) GetCcnRegionBandwidthLimits(ctx context.Context,
	ccnID string) (infos []vpc.CcnRegionBandwidthLimit, errRet error) {
	var (
		request  = vpc.NewGetCcnRegionBandwidthLimitsRequest()
		response *vpc.GetCcnRegionBandwidthLimitsResponse
		err      error
		limit    uint64 = 100
		offset   uint64 = 0
	)
	request.CcnId = &ccnID
	request.Limit = &limit
	request.Offset = &offset

	ratelimit.Check(request.GetAction())
	for {
		response, err = me.client.UseVpcClient().GetCcnRegionBandwidthLimits(request)
		if err != nil {
			errRet = err
			return
		}
		if response.Response == nil || response.Response.CcnBandwidthSet == nil {
			errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
			return
		}

		for _, item := range response.Response.CcnBandwidthSet {
			if item.CcnRegionBandwidthLimit != nil {
				infos = append(infos, *item.CcnRegionBandwidthLimit)
			}
		}
		if len(response.Response.CcnBandwidthSet) < int(limit) {
			break
		}
		offset += limit
	}
	return
}

func (me *VpcService) SetCcnRegionBandwidthLimits(ctx context.Context, ccnId, region, dstRegion string,
	bandwidth int64, setFlag bool) (errRet error) {

	logId := tccommon.GetLogId(ctx)
	request := vpc.NewSetCcnRegionBandwidthLimitsRequest()
	request.CcnId = &ccnId

	var uint64bandwidth = uint64(bandwidth)
	var ccnRegionBandwidthLimit vpc.CcnRegionBandwidthLimit
	ccnRegionBandwidthLimit.BandwidthLimit = &uint64bandwidth
	ccnRegionBandwidthLimit.Region = &region
	if dstRegion != "" {
		ccnRegionBandwidthLimit.DstRegion = &dstRegion
	}

	request.CcnRegionBandwidthLimits = []*vpc.CcnRegionBandwidthLimit{&ccnRegionBandwidthLimit}

	request.SetDefaultLimitFlag = helper.Bool(setFlag)
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().SetCcnRegionBandwidthLimits(request)

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
	log.Printf("[DEBUG]%s api[%s] , request body [%s], response body[%s]\n",
		logId,
		request.GetAction(),
		request.ToJsonString(),
		response.ToJsonString())
	return
}

func (me *VpcService) ModifyCcnRegionBandwidthLimitsType(ctx context.Context, ccnID, limitType string) error {
	request := vpc.NewModifyCcnRegionBandwidthLimitsTypeRequest()
	request.CcnId = &ccnID
	request.BandwidthLimitType = &limitType
	_, err := me.client.UseVpcClient().ModifyCcnRegionBandwidthLimitsType(request)
	if err != nil {
		return err
	}
	return nil
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

func (me *VpcService) DescribeRouteTableAssociatedInstancesById(ctx context.Context, meta interface{}, ccnId, routeTableId string) (instanceBindList []CcnInstanceBind, errRet error) {
	body := map[string]interface{}{
		"Filters": []map[string]interface{}{
			{
				"Name":   "ccn-id",
				"Values": []string{ccnId},
			},
			{
				"Name":   "ccn-route-table-id",
				"Values": []string{routeTableId},
			},
		},
		"Offset": 0,
		"Limit":  100,
	}

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOmitNilClient("vpc")
	request := tchttp.NewCommonRequest("vpc", "2017-03-12", "DescribeRouteTableAssociatedInstances")
	err := request.SetActionParameters(body)
	if err != nil {
		errRet = err
		return
	}

	response := tchttp.NewCommonResponse()
	err = client.Send(request, response)
	if err != nil {
		fmt.Printf("describe vpc DescribeRouteTableAssociatedInstances failed: %v \n", err)
		errRet = err
		return
	}

	resultStr := string(response.GetBody())
	var ccnInstanceBindApiResponse CcnInstanceBindApiResponse
	err = json.Unmarshal([]byte(resultStr), &ccnInstanceBindApiResponse)
	if err != nil {
		errRet = err
		return
	}

	InstanceBindSet := ccnInstanceBindApiResponse.Response.InstanceBindSet
	for _, BindSet := range InstanceBindSet {
		var tmpBind CcnInstanceBind
		if BindSet.InstanceId != "" {
			tmpBind.instanceId = BindSet.InstanceId
		}

		if BindSet.InstanceType != "" {
			tmpBind.instanceType = BindSet.InstanceType
		}

		instanceBindList = append(instanceBindList, tmpBind)
	}

	return

	//logId := tccommon.GetLogId(ctx)
	//
	//request := vpc.NewDescribeRouteTableAssociatedInstancesRequest()
	//request.Filters = []*vpc.Filter{
	//	{
	//		Name:   helper.String("ccn-id"),
	//		Values: helper.Strings([]string{ccnId}),
	//	},
	//	{
	//		Name:   helper.String("ccn-route-table-id"),
	//		Values: helper.Strings([]string{routeTableId}),
	//	},
	//}
	//
	//defer func() {
	//	if errRet != nil {
	//		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
	//	}
	//}()
	//
	//ratelimit.Check(request.GetAction())
	//response, err := me.client.UseVpcClient().DescribeRouteTableAssociatedInstances(request)
	//if err != nil {
	//	errRet = err
	//	return
	//}
	//
	//log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	//if response.Response == nil || len(response.Response.InstanceBindSet) < 1 {
	//	return
	//}
	//
	//instanceBindList = response.Response.InstanceBindSet
	//return
}

func (me *VpcService) DescribeVpcReplaceCcnRouteTableInputPolicysById(ctx context.Context, ccnId, routeTableId string) (policySet *vpc.CcnRouteTableInputPolicys, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDescribeCcnRouteTableInputPolicysRequest()
	request.CcnId = helper.String(ccnId)
	request.RouteTableId = helper.String(routeTableId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().DescribeCcnRouteTableInputPolicys(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if response.Response == nil || len(response.Response.PolicySet) < 1 {
		return
	}

	policySet = response.Response.PolicySet[0]
	return
}

func (me *VpcService) DescribeVpcReplaceCcnRouteTableBroadcastPolicysById(ctx context.Context, ccnId, routeTableId string) (policySet *vpc.CcnRouteTableBroadcastPolicys, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDescribeCcnRouteTableBroadcastPolicysRequest()
	request.CcnId = helper.String(ccnId)
	request.RouteTableId = helper.String(routeTableId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().DescribeCcnRouteTableBroadcastPolicys(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if response.Response == nil || len(response.Response.PolicySet) < 1 {
		return
	}

	policySet = response.Response.PolicySet[0]
	return
}

func (me *VpcService) DescribeVpcReplaceCcnRouteTableSelectionPolicysById(ctx context.Context, ccnId string) (routeSelectionPolicySet []*vpc.RouteSelectionPolicy, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDescribeRouteTableSelectionPoliciesRequest()
	request.Filters = []*vpc.Filter{
		{
			Name:   helper.String("ccn-id"),
			Values: helper.Strings([]string{ccnId}),
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().DescribeRouteTableSelectionPolicies(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if response.Response == nil || len(response.Response.RouteSelectionPolicySet) == 0 {
		return
	}

	routeSelectionPolicySet = response.Response.RouteSelectionPolicySet
	return
}

func (me *VpcService) DescribeVpcCcnRouteTablesById(ctx context.Context, routeTableId string) (ccnRouteTable *vpc.CcnRouteTable, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDescribeCcnRouteTablesRequest()
	request.Filters = []*vpc.Filter{
		{
			Name:   helper.String("route-table-id"),
			Values: helper.Strings([]string{routeTableId}),
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseVpcClient().DescribeCcnRouteTables(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if response.Response == nil || len(response.Response.CcnRouteTableSet) != 1 {
		return
	}

	ccnRouteTable = response.Response.CcnRouteTableSet[0]
	return
}

func (me *VpcService) DescribeVpcCcnRoutesById(ctx context.Context, ccnId string, routeId string) (ccnRoutes *vpc.CcnRoute, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vpc.NewDescribeCcnRoutesRequest()
	request.CcnId = &ccnId
	request.Filters = []*vpc.Filter{
		{
			Name:   helper.String("route-id"),
			Values: []*string{helper.String(routeId)},
		},
	}
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

	if response != nil && response.Response != nil && len(response.Response.RouteSet) > 0 {
		ccnRoutes = response.Response.RouteSet[0]
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

func (me *VpcService) DescribeCcnRouteTableInputPoliciesByFilter(ctx context.Context, param map[string]interface{}) (ret []*vpc.CcnRouteTableInputPolicys, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = vpc.NewDescribeCcnRouteTableInputPolicysRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "CcnId" {
			request.CcnId = v.(*string)
		}
		if k == "RouteTableId" {
			request.RouteTableId = v.(*string)
		}
		if k == "PolicyVersion" {
			request.PolicyVersion = v.(*uint64)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCcnV20170312Client().DescribeCcnRouteTableInputPolicys(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.PolicySet) < 1 {
		return
	}

	ret = response.Response.PolicySet
	return
}
