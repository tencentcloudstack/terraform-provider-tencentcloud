package privatedns

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	privatednsIntlv20201028 "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/privatedns/v20201028"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	privatedns "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/privatedns/v20201028"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

// basic information

type PrivateDnsService struct {
	client *connectivity.TencentCloudClient
}

type PrivatednsService struct {
	client *connectivity.TencentCloudClient
}

// ////////api
func (me *PrivateDnsService) DescribePrivateDnsRecordByFilter(ctx context.Context, zoneId string,
	filterList []*privatedns.Filter) (recordInfos []*privatedns.PrivateZoneRecord, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := privatedns.NewDescribePrivateZoneRecordListRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	var (
		limit  int64 = 20
		offset int64 = 0
		total  int64 = -1
	)
	request.ZoneId = &zoneId

	if filterList != nil {
		request.Filters = filterList
	}

getMoreData:

	if total >= 0 {
		if offset >= total {
			return
		}
	}
	var response *privatedns.DescribePrivateZoneRecordListResponse

	ratelimit.Check(request.GetAction())
	request.Limit = &limit
	request.Offset = &offset

	if err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UsePrivateDnsClient().DescribePrivateZoneRecordList(request)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		response = result
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s read private dns failed, reason: %v", logId, err)
		return nil, err
	}
	if total < 0 {
		total = *response.Response.TotalCount
	}

	if len(response.Response.RecordSet) > 0 {
		offset = offset + limit
	} else {
		return
	}

	recordInfos = append(recordInfos, response.Response.RecordSet...)
	goto getMoreData
}

func (me *PrivateDnsService) DescribePrivateDnsZoneVpcAttachmentById(ctx context.Context, zoneId string) (ZoneVpcAttachment *privatedns.PrivateZone, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := privatedns.NewDescribePrivateZoneRequest()
	request.ZoneId = &zoneId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UsePrivateDnsClient().DescribePrivateZone(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.PrivateZone == nil {
		return
	}

	ZoneVpcAttachment = response.Response.PrivateZone

	return
}

func (me *PrivateDnsService) DeletePrivateDnsZoneVpcAttachmentById(ctx context.Context, zoneId, uniqVpcId, region, uin string) (errRet error) {
	var (
		logId        = tccommon.GetLogId(ctx)
		asyncRequest = privatedns.NewQueryAsyncBindVpcStatusRequest()
		uniqId       string
	)

	request := privatedns.NewDeleteSpecifyPrivateZoneVpcRequest()
	request.ZoneId = &zoneId
	request.Sync = common.BoolPtr(false)
	if uin == "" {
		request.VpcSet = []*privatedns.VpcInfo{
			{
				UniqVpcId: common.StringPtr(uniqVpcId),
				Region:    common.StringPtr(region),
			},
		}
	} else {
		request.AccountVpcSet = []*privatedns.AccountVpcInfo{
			{
				UniqVpcId: common.StringPtr(uniqVpcId),
				Region:    common.StringPtr(region),
				Uin:       common.StringPtr(uin),
				VpcName:   common.StringPtr(""),
			},
		}
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UsePrivateDnsClient().DeleteSpecifyPrivateZoneVpc(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response.UniqId == nil {
		return fmt.Errorf("Delete specify private zone vpc failed.")
	}

	uniqId = *response.Response.UniqId

	// wait
	asyncRequest.UniqId = &uniqId
	err = resource.Retry(tccommon.ReadRetryTimeout*5, func() *resource.RetryError {
		result, e := me.client.UsePrivateDnsClient().QueryAsyncBindVpcStatus(asyncRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, asyncRequest.GetAction(), asyncRequest.ToJsonString(), asyncRequest.ToJsonString())
		}

		if *result.Response.Status == "success" {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("query async bind vpc status is %s.", *result.Response.Status))
	})

	if err != nil {
		log.Printf("[CRITAL]%s query async bind vpc status failed, reason:%+v", logId, err)
		return err
	}

	return
}

func (me *PrivateDnsService) DescribePrivatednsPrivateZoneListByFilter(ctx context.Context, param map[string]interface{}) (privateZoneList []*privatedns.PrivateZone, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = privatedns.NewDescribePrivateZoneListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Filters" {
			request.Filters = v.([]*privatedns.Filter)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UsePrivateDnsClient().DescribePrivateZoneList(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.PrivateZoneSet) < 1 {
		return
	}

	privateZoneList = response.Response.PrivateZoneSet
	return
}

func (me *PrivatednsService) DescribePrivateDnsForwardRuleById(ctx context.Context, ruleId string) (ret *privatednsIntlv20201028.ForwardRule, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := privatednsIntlv20201028.NewDescribeForwardRuleRequest()
	request.RuleId = helper.String(ruleId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UsePrivatednsIntlV20201028Client().DescribeForwardRule(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil {
		return
	}

	ret = response.Response.ForwardRule
	return
}

func (me *PrivatednsService) DescribePrivateDnsEndPointById(ctx context.Context, endPointId string) (ret *privatednsIntlv20201028.DescribeEndPointListResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := privatednsIntlv20201028.NewDescribeEndPointListRequest()
	filter := &privatednsIntlv20201028.Filter{
		Name:   helper.String("EndPointId"),
		Values: []*string{helper.String(endPointId)},
	}
	request.Filters = append(request.Filters, filter)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UsePrivatednsIntlV20201028Client().DescribeEndPointList(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	ret = response.Response
	return
}

func (me *PrivatednsService) DescribePrivateDnsExtendEndPointById(ctx context.Context, endPointId string) (ret *privatednsIntlv20201028.DescribeExtendEndpointListResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := privatednsIntlv20201028.NewDescribeExtendEndpointListRequest()
	response := privatednsIntlv20201028.NewDescribeExtendEndpointListResponse()
	filter := &privatednsIntlv20201028.Filter{
		Name:   helper.String("EndpointId"),
		Values: []*string{helper.String(endPointId)},
	}
	request.Filters = append(request.Filters, filter)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UsePrivatednsIntlV20201028Client().DescribeExtendEndpointList(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	ret = response.Response
	return
}

func (me *PrivatednsService) DescribePrivateDnsForwardRulesByFilter(ctx context.Context, param map[string]interface{}) (ret []*privatedns.ForwardRule, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = privatedns.NewDescribeForwardRuleListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Filters" {
			request.Filters = v.([]*privatedns.Filter)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 100
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UsePrivatednsV20201028Client().DescribeForwardRuleList(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.ForwardRuleSet) < 1 {
			break
		}
		ret = append(ret, response.Response.ForwardRuleSet...)
		if len(response.Response.ForwardRuleSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *PrivatednsService) DescribePrivateDnsEndPointsByFilter(ctx context.Context, param map[string]interface{}) (ret []*privatedns.EndPointInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = privatedns.NewDescribeEndPointListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Filters" {
			request.Filters = v.([]*privatedns.Filter)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 100
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UsePrivatednsV20201028Client().DescribeEndPointList(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.EndPointSet) < 1 {
			break
		}
		ret = append(ret, response.Response.EndPointSet...)
		if len(response.Response.EndPointSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}
