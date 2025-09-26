package billing

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	billingv20180709 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/billing/v20180709"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewBillingService(client *connectivity.TencentCloudClient) BillingService {
	return BillingService{client: client}
}

type BillingService struct {
	client *connectivity.TencentCloudClient
}

func (me *BillingService) DescribeBillingAllocationTagById(ctx context.Context, tagKey string) (ret *billingv20180709.TagDataInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := billingv20180709.NewDescribeTagListRequest()
	response := billingv20180709.NewDescribeTagListResponse()
	request.TagKey = helper.String(tagKey)
	request.Limit = helper.IntUint64(1)
	request.Offset = helper.IntUint64(0)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var (
		limit    uint64 = 1000
		offset   uint64 = 0
		dataList []*billingv20180709.TagDataInfo
	)

	for {
		request.Offset = &offset
		request.Limit = &limit
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseBillingV20180709Client().DescribeTagList(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe billing allocation tag failed, Response is nil."))
			}

			response = result
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		if response == nil || len(response.Response.Data) < 1 {
			break

		}
		dataList = append(dataList, response.Response.Data...)
		if len(response.Response.Data) < int(limit) {
			break
		}

		offset += limit
	}

	for _, item := range dataList {
		if item.TagKey != nil && *item.TagKey == tagKey {
			ret = item
			break
		}
	}

	return
}

func (me *BillingService) DescribeBillingBudgetById(ctx context.Context, budgetId string) (ret *billingv20180709.DescribeBudgetResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := billingv20180709.NewDescribeBudgetRequest()
	request.BudgetId = helper.String(budgetId)
	request.PageNo = helper.IntInt64(1)
	request.PageSize = helper.IntInt64(10)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseBillingV20180709Client().DescribeBudget(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	ret = response.Response
	return
}

func (me *BillingService) DescribeBillingBudgetOperationLogByFilter(ctx context.Context, param map[string]interface{}) (records []*billingv20180709.BudgetOperationLogEntity, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := billingv20180709.NewDescribeBudgetOperationLogRequest()
	response := billingv20180709.NewDescribeBudgetOperationLogResponse()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "BudgetId" {
			request.BudgetId = v.(*string)
		}
	}
	var (
		pageNo   int64 = 1
		pageSize int64 = 100
	)

	for {
		request.PageNo = helper.Int64(pageNo)
		request.PageSize = helper.Int64(pageSize)
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseBillingV20180709Client().DescribeBudgetOperationLog(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe billing allocation tag failed, Response is nil."))
			}

			response = result
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		if response == nil || response.Response.Data == nil || len(response.Response.Data.Records) == 0 {
			break

		}
		records = append(records, response.Response.Data.Records...)
		if len(response.Response.Data.Records) < int(pageSize) {
			break
		}

		pageNo += 1
	}
	return
}
