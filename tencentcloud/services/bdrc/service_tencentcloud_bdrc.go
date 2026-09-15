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

func (me *BdrcService) DescribeBdrcInstanceCopyPairById(ctx context.Context, copyPairId string) (ret *bdrcv20260330.CopyPair, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := bdrcv20260330.NewDescribeCopyPairsRequest()
	request.CopyPairType = helper.String("INSTANCE")
	request.CopyPairIds = []*string{helper.String(copyPairId)}
	request.Limit = helper.IntInt64(100)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var copyPair *bdrcv20260330.CopyPair
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseBdrcV20260330Client().DescribeCopyPairsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil || len(result.Response.CopyPairSet) == 0 {
			return resource.NonRetryableError(fmt.Errorf("Describe bdrc_instance_copy_pair failed, CopyPairSet is empty, copyPairId is %s.", copyPairId))
		}

		for _, item := range result.Response.CopyPairSet {
			if item.CopyPairId != nil && *item.CopyPairId == copyPairId {
				copyPair = item
				break
			}
		}

		if copyPair == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe bdrc_instance_copy_pair failed, copyPairId %s is not found in CopyPairSet.", copyPairId))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRUD] bdrc_instance_copy_pair id=%s", copyPairId)
		return nil, nil
	}

	ret = copyPair
	return
}
