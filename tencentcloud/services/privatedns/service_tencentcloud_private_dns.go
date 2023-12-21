package privatedns

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	privatedns "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/privatedns/v20201028"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

// basic information

type PrivateDnsService struct {
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
