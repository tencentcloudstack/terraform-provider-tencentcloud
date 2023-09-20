package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	privatedns "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/privatedns/v20201028"
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
	logId := getLogId(ctx)
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

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UsePrivateDnsClient().DescribePrivateZoneRecordList(request)
		if err != nil {
			return retryError(err, InternalError)
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
	logId := getLogId(ctx)

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
	logId := getLogId(ctx)

	request := privatedns.NewDeleteSpecifyPrivateZoneVpcRequest()
	request.ZoneId = &zoneId
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

	return
}
