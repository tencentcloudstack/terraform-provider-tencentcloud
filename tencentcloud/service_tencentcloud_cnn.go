package tencentcloud

import (
	"context"
	"fmt"
	"log"

	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

//Cnn basic information
type CnnBasicInfo struct {
	cnnId         string
	name          string
	description   string
	state         string
	qos           string
	instanceCount int64
	createTime    string
}

type CnnAttachedInstanceInfo struct {
	cnnId          string
	instanceType   string
	instanceRegion string
	instanceId     string
	state          string
	attachedTime   string
	cidrBlock      []string
}

type CnnBandwidthLimit struct {
	region string
	limit  int64
}

func (me *VpcService) DescribeCcn(ctx context.Context, cnnId string) (info CnnBasicInfo, has int, errRet error) {
	infos, err := me.DescribeCcns(ctx, cnnId, "")
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

func (me *VpcService) DescribeCcns(ctx context.Context, cnnId, name string) (infos []CnnBasicInfo, errRet error) {

	logId := GetLogId(ctx)
	request := vpc.NewDescribeCcnsRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	infos = make([]CnnBasicInfo, 0, 100)

	var offset uint64 = 0
	var limit uint64 = 100
	var total = -1
	var has = map[string]bool{}

	var filters []*vpc.Filter
	if cnnId != "" {
		filters = me.fillFilter(filters, "ccn-id", cnnId)
	}
	if name != "" {
		filters = me.fillFilter(filters, "ccn-name", name)
	}
	if len(filters) > 0 {
		request.Filters = filters
	}

getMoreData:

	if total >= 0 {
		if int(offset) >= total {
			return
		}
	}
	request.Limit = &limit
	request.Offset = &offset

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
		var basicInfo CnnBasicInfo

		basicInfo.cnnId = *item.CcnId
		basicInfo.name = *item.CcnName
		basicInfo.createTime = *item.CreateTime

		basicInfo.description = *item.CcnDescription
		basicInfo.instanceCount = int64(*item.InstanceCount)
		basicInfo.qos = *item.QosLevel
		basicInfo.state = *item.State

		if has[basicInfo.cnnId] {
			errRet = fmt.Errorf("get repeated cnn_id[%s] when doing DescribeCcns", basicInfo.cnnId)
			return
		}
		has[basicInfo.cnnId] = true
		infos = append(infos, basicInfo)
	}
	goto getMoreData

}

func (me *VpcService) DescribeCcnRegionBandwidthLimits(ctx context.Context, cnnId string) (infos []CnnBandwidthLimit, errRet error) {

	logId := GetLogId(ctx)
	request := vpc.NewDescribeCcnRegionBandwidthLimitsRequest()

	infos = make([]CnnBandwidthLimit, 0, 100)

	request.CcnId = &cnnId

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

		var cnnBandwidthLimit CnnBandwidthLimit
		cnnBandwidthLimit.region = *item.Region
		cnnBandwidthLimit.limit = int64(*item.BandwidthLimit)
		infos = append(infos, cnnBandwidthLimit)
	}
	return
}

func (me *VpcService) CreateCcn(ctx context.Context, name, description, qos string) (basicInfo CnnBasicInfo, errRet error) {

	logId := GetLogId(ctx)
	request := vpc.NewCreateCcnRequest()

	request.CcnName = &name
	request.CcnDescription = &description
	request.QosLevel = &qos

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
	basicInfo.cnnId = *item.CcnId
	basicInfo.name = *item.CcnName
	basicInfo.createTime = *item.CreateTime

	basicInfo.description = *item.CcnDescription
	basicInfo.instanceCount = int64(*item.InstanceCount)
	basicInfo.qos = *item.QosLevel
	basicInfo.state = *item.State
	return
}

func (me *VpcService) DeleteCcn(ctx context.Context, cnnId string) (errRet error) {

	logId := GetLogId(ctx)
	request := vpc.NewDeleteCcnRequest()
	request.CcnId = &cnnId

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

func (me *VpcService) ModifyCcnAttribute(ctx context.Context, cnnId, name, description string) (errRet error) {

	logId := GetLogId(ctx)
	request := vpc.NewModifyCcnAttributeRequest()
	request.CcnId = &cnnId

	if name != "" {
		request.CcnName = &name
	}
	if description != "" {
		request.CcnDescription = &description
	}

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
	return
}

func (me *VpcService) DescribeCcnAttachedInstances(ctx context.Context, cnnId string) (infos []CnnAttachedInstanceInfo, errRet error) {

	logId := GetLogId(ctx)
	request := vpc.NewDescribeCcnAttachedInstancesRequest()
	request.CcnId = &cnnId

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

	infos = make([]CnnAttachedInstanceInfo, 0, len(response.Response.InstanceSet))

	for _, item := range response.Response.InstanceSet {

		var info CnnAttachedInstanceInfo

		info.attachedTime = *item.AttachedTime
		info.cidrBlock = make([]string, 0, len(item.CidrBlock))

		for _, v := range item.CidrBlock {
			info.cidrBlock = append(info.cidrBlock, *v)
		}

		info.cnnId = cnnId
		info.instanceId = *item.InstanceId
		info.instanceRegion = *item.InstanceRegion
		info.instanceType = *item.InstanceType
		info.state = *item.State
		infos = append(infos, info)
	}
	return
}
