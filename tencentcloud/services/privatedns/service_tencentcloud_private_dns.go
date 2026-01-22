package privatedns

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	privatednsIntlv20201028 "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/privatedns/v20201028"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	privatedns "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/privatedns/v20201028"

	intlSdkError "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common/errors"
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
		limit  int64 = 200
		offset int64 = 0
		total  int64 = -1
	)
	request.ZoneId = &zoneId

	if filterList != nil {
		request.Filters = filterList
	}

	var tmpRetry = PRIVATEDNS_CUSTOM_RETRY_SDK_ERROR
	tmpRetry = append(tmpRetry, tccommon.InternalError)

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
			return tccommon.RetryError(err, tmpRetry...)
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
	response := privatedns.NewDescribePrivateZoneResponse()
	request.ZoneId = &zoneId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UsePrivateDnsClient().DescribePrivateZone(request)
		if e != nil {
			return tccommon.RetryError(e, PRIVATEDNS_CUSTOM_RETRY_SDK_ERROR...)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe PrivateDns zone failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

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
	response := privatedns.NewDeleteSpecifyPrivateZoneVpcResponse()
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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UsePrivateDnsClient().DeleteSpecifyPrivateZoneVpc(request)
		if e != nil {
			return tccommon.RetryError(e, PRIVATEDNS_CUSTOM_RETRY_SDK_ERROR...)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.UniqId == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe PrivateDns zone failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	// wait
	uniqId = *response.Response.UniqId
	asyncRequest.UniqId = &uniqId
	err = resource.Retry(tccommon.ReadRetryTimeout*5, func() *resource.RetryError {
		result, e := me.client.UsePrivateDnsClient().QueryAsyncBindVpcStatus(asyncRequest)
		if e != nil {
			return tccommon.RetryError(e, PRIVATEDNS_CUSTOM_RETRY_SDK_ERROR...)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, asyncRequest.GetAction(), asyncRequest.ToJsonString(), asyncRequest.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Query async bind vpc status failed, Response is nil."))
		}

		if result.Response.Status == nil {
			return resource.NonRetryableError(fmt.Errorf("Status is nil."))
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
		logId    = tccommon.GetLogId(ctx)
		request  = privatedns.NewDescribePrivateZoneListRequest()
		response = privatedns.NewDescribePrivateZoneListResponse()
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

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UsePrivateDnsClient().DescribePrivateZoneList(request)
		if e != nil {
			return tccommon.RetryError(e, PRIVATEDNS_CUSTOM_RETRY_SDK_ERROR...)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe PrivateDns zone list failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	if len(response.Response.PrivateZoneSet) < 1 {
		return
	}

	privateZoneList = response.Response.PrivateZoneSet
	return
}

func (me *PrivatednsService) DescribePrivateDnsForwardRuleById(ctx context.Context, ruleId string) (ret *privatednsIntlv20201028.ForwardRule, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := privatednsIntlv20201028.NewDescribeForwardRuleRequest()
	response := privatednsIntlv20201028.NewDescribeForwardRuleResponse()
	request.RuleId = helper.String(ruleId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UsePrivatednsIntlV20201028Client().DescribeForwardRule(request)
		if e != nil {
			return tccommon.RetryError(e, PRIVATEDNS_CUSTOM_RETRY_SDK_ERROR...)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe forward rule failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response.ForwardRule
	return
}

func (me *PrivatednsService) DescribePrivateDnsEndPointById(ctx context.Context, endPointId string) (ret *privatednsIntlv20201028.DescribeEndPointListResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := privatednsIntlv20201028.NewDescribeEndPointListRequest()
	response := privatednsIntlv20201028.NewDescribeEndPointListResponse()

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

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UsePrivatednsIntlV20201028Client().DescribeEndPointList(request)
		if e != nil {
			return tccommon.RetryError(e, PRIVATEDNS_CUSTOM_RETRY_SDK_ERROR...)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe end point list failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

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

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UsePrivatednsIntlV20201028Client().DescribeExtendEndpointList(request)
		if e != nil {
			return tccommon.RetryError(e, PRIVATEDNS_CUSTOM_RETRY_SDK_ERROR...)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe extend end point list failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

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

func (me *PrivatednsService) DescribePrivateDnsEndPointsByFilter(ctx context.Context, param map[string]interface{}) (ret []*privatednsIntlv20201028.EndPointInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = privatednsIntlv20201028.NewDescribeEndPointListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Filters" {
			request.Filters = v.([]*privatednsIntlv20201028.Filter)
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
		response, err := me.client.UsePrivatednsIntlV20201028Client().DescribeEndPointList(request)
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

func (me *PrivateDnsService) DescribePrivateDnsRecordById(ctx context.Context, zoneId, recordId string) (recordInfo *privatednsIntlv20201028.RecordInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := privatednsIntlv20201028.NewDescribeRecordRequest()
	response := privatednsIntlv20201028.NewDescribeRecordResponse()
	request.ZoneId = &zoneId
	request.RecordId = &recordId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UsePrivatednsIntlV20201028Client().DescribeRecord(request)
		if e != nil {
			if sdkError, ok := e.(*intlSdkError.TencentCloudSDKError); ok {
				if sdkError.Code == "InvalidParameter.RecordNotExist" {
					return nil
				}
			}

			return tccommon.RetryError(e, PRIVATEDNS_CUSTOM_RETRY_SDK_ERROR...)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe PrivateDns record %s failed, Response is nil.", recordId))
		}

		if result.Response.RecordInfo != nil && result.Response.RecordInfo.RecordId != nil {
			respRecordId := *result.Response.RecordInfo.RecordId
			if respRecordId == recordId {
				response = result
				return nil
			} else {
				return resource.NonRetryableError(fmt.Errorf("Describe PrivateDns record %s does not meet expectations, Response is %s.", recordId, respRecordId))
			}
		}

		return resource.RetryableError(fmt.Errorf("Record %s is still creating...", recordId))
	})

	if err != nil {
		errRet = err
		return
	}

	if response != nil && response.Response != nil {
		recordInfo = response.Response.RecordInfo
	}

	return
}

func (me *PrivatednsService) DescribePrivateDnsInboundEndpointById(ctx context.Context, endpointId string) (ret *privatedns.InboundEndpointSet, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := privatedns.NewDescribeInboundEndpointListRequest()
	response := privatedns.NewDescribeInboundEndpointListResponse()
	request.Filters = []*privatedns.Filter{
		{
			Name:   common.StringPtr("EndPointId"),
			Values: common.StringPtrs([]string{endpointId}),
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UsePrivatednsV20201028Client().DescribeInboundEndpointList(request)
		if e != nil {
			return tccommon.RetryError(e, PRIVATEDNS_CUSTOM_RETRY_SDK_ERROR...)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.InboundEndpointSet == nil || len(result.Response.InboundEndpointSet) == 0 {
			return resource.NonRetryableError(fmt.Errorf("Describe inbound endpoint list failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response.InboundEndpointSet[0]
	return
}
