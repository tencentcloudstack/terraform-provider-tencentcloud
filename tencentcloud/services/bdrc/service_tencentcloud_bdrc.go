package bdrc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	bdrcv20260330 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bdrc/v20260330"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewBdrcService(client *connectivity.TencentCloudClient) BdrcService {
	return BdrcService{client: client}
}

type BdrcService struct {
	client *connectivity.TencentCloudClient
}

func (me *BdrcService) DescribeDisasterRecoverySitePairById(ctx context.Context, sitePairId, sitePairType string) (ret *bdrcv20260330.SitePair, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := bdrcv20260330.NewDescribeDisasterRecoverySitePairsRequest()
	request.SitePairIds = helper.Strings([]string{sitePairId})
	request.SitePairType = helper.String(sitePairType)
	request.Limit = helper.Int64(100)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseBdrcV20260330Client().DescribeDisasterRecoverySitePairsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe bdrc disaster recovery site pair failed, Response is nil."))
		}

		for _, sitePair := range result.Response.SitePairSet {
			if sitePair.SitePairId != nil && *sitePair.SitePairId == sitePairId {
				ret = sitePair
				break
			}
		}

		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}
