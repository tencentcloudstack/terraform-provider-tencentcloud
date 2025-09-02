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
