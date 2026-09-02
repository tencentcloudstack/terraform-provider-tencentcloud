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

func (me *BillingService) DescribeBillingInstanceById(ctx context.Context, instanceId string) (ret *billingv20180709.RenewInstance, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := billingv20180709.NewDescribeRenewInstancesRequest()
	response := billingv20180709.NewDescribeRenewInstancesResponse()
	request.MaxResults = helper.IntUint64(1)
	request.InstanceIdList = []*string{&instanceId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseBillingV20180709Client().DescribeRenewInstances(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe billing renew instance failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	if len(response.Response.InstanceList) == 0 {
		return
	}

	ret = response.Response.InstanceList[0]
	return
}

func (me *BillingService) DescribeBillingBillDetailByFilter(ctx context.Context, param map[string]interface{}) (detailSet []*billingv20180709.BillDetail, total uint64, retContext *string, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := billingv20180709.NewDescribeBillDetailRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "PeriodType" {
			request.PeriodType = v.(*string)
		}
		if k == "Month" {
			request.Month = v.(*string)
		}
		if k == "BeginTime" {
			request.BeginTime = v.(*string)
		}
		if k == "EndTime" {
			request.EndTime = v.(*string)
		}
		if k == "NeedRecordNum" {
			request.NeedRecordNum = v.(*int64)
		}
		if k == "ProductCode" {
			request.ProductCode = v.(*string)
		}
		if k == "PayMode" {
			request.PayMode = v.(*string)
		}
		if k == "ResourceId" {
			request.ResourceId = v.(*string)
		}
		if k == "ActionType" {
			request.ActionType = v.(*string)
		}
		if k == "ProjectId" {
			request.ProjectId = v.(*int64)
		}
		if k == "BusinessCode" {
			request.BusinessCode = v.(*string)
		}
		if k == "Context" {
			request.Context = v.(*string)
		}
		if k == "PayerUin" {
			request.PayerUin = v.(*string)
		}
	}

	var (
		limit  uint64 = 300
		offset uint64 = 0
	)

	for {
		request.Offset = &offset
		request.Limit = &limit
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseBillingV20180709Client().DescribeBillDetail(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe billing bill detail failed, Response is nil."))
			}

			if result.Response.DetailSet == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe billing bill detail failed, DetailSet is nil."))
			}

			detailSet = append(detailSet, result.Response.DetailSet...)
			if result.Response.Total != nil {
				total = *result.Response.Total
			}
			if result.Response.Context != nil {
				retContext = result.Response.Context
			}
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		if uint64(len(detailSet)) >= total || total == 0 {
			break
		}

		offset += limit
	}

	return
}
