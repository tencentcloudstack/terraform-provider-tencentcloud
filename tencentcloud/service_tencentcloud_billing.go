package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"

	billing "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/billing/v20180709"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
)

type BillingService struct {
	client *connectivity.TencentCloudClient
}

func (me *BillingService) DescribeDeals(ctx context.Context, dealId string) (deal *billing.Deal, errRet error) {
	var (
		logId    = getLogId(ctx)
		request  = billing.NewDescribeDealsByCondRequest()
		response = billing.NewDescribeDealsByCondResponse()
		dealList []*billing.Deal
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.OrderId = helper.String(dealId)
	request.Limit = helper.IntInt64(20)
	baseTime := time.Now().Local()
	startTime := baseTime.AddDate(0, 0, -7).Format("2006-01-02 15:04:05")
	endTime := baseTime.Format("2006-01-02 15:04:05")
	request.StartTime = helper.String(startTime)
	request.EndTime = helper.String(endTime)

	ratelimit.Check(request.GetAction())
	err := resource.Retry(20*readRetryTimeout, func() *resource.RetryError {
		for _, dealState := range DEAL_STATUS_CODE {
			request.Status = helper.Int64(dealState)
			response, ee := me.client.UseBillingClient().DescribeDealsByCond(request)
			if ee != nil {
				return retryError(errRet, InternalError)
			}
			dealList = response.Response.Deals
			if len(dealList) > 0 {
				break
			}
		}
		if len(dealList) != 1 {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
			return nil
		}
		deal = dealList[0]
		if in(*deal.Status, DEAL_TERMINATE_STATUS_CODE) {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
			return nil
		}
		time.Sleep(3 * time.Second)
		return resource.RetryableError(fmt.Errorf("deal status is not ready..., current is: %v", deal.Status))
	})

	if err != nil {
		errRet = err
		return
	}
	return
}

func in(target int64, intArr []int64) bool {
	for _, element := range intArr {
		if target == element {
			return true
		}
	}
	return false
}
