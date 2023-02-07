package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"

	billing "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/billing/v20180709"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
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

func (me *BillingService) isYunTiAccount() bool {
	val, ok := os.LookupEnv(PROVIDER_ENABLE_YUNTI)
	if ok && strings.ToLower(val) == "true" {
		return true
	}
	return false
}

//query deal by bpass
func (me *BillingService) QueryDealByBpass(ctx context.Context, dealRegx string, msg error) (resourceId *string, err error) {
	logId := getLogId(ctx)
	err = msg
	if !me.isYunTiAccount() {
		return nil, err
	}

	e, ok := msg.(*sdkErrors.TencentCloudSDKError)
	log.Printf("[DEBUG]%s query deal for PREPAID user, msg:[%s] \n", logId, e.Code)

	if ok && IsContains(TRADE_RETRYABLE_ERROR, e.Code) {
		errStr := msg.Error()

		re := regexp.MustCompile(dealRegx)
		result := re.FindStringSubmatch(errStr)
		for i, str := range result {
			log.Printf("[DEBUG] FindStringSubmatch sub[%v]:%s,\n", i, str)
		}
		dealId := re.FindStringSubmatch(errStr)[1]
		deal, billErr := me.DescribeDeals(ctx, dealId)
		if billErr != nil {
			log.Printf("[CRITAL]%s api[DescribeDeals] fail, reason[%s]\n", logId, billErr.Error())
			return nil, billErr
		}
		resourceId = deal.ResourceId[0]
		log.Printf("[DEBUG]%s query deal for PREPAID user succeed, dealId:[%s] resourceId:[%s]\n", logId, dealId, *resourceId)
		return
	}
	return nil, err
}

func in(target int64, intArr []int64) bool {
	for _, element := range intArr {
		if target == element {
			return true
		}
	}
	return false
}
