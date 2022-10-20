package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

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
	recordId string) (recordInfos []*privatedns.PrivateZoneRecord, errRet error) {
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
	request.Filters = make([]*privatedns.Filter, 0)

	if recordId != "" {
		filter := privatedns.Filter{
			Name:   helper.String("RecordId"),
			Values: []*string{&recordId},
		}
		request.Filters = append(request.Filters, &filter)
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
