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

func (me *BdrcService) DescribeDisasterRecoveryProtectGroupById(ctx context.Context, protectGroupId, protectGroupType string) (ret *bdrcv20260330.ProtectGroup, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = bdrcv20260330.NewDescribeDisasterRecoveryProtectGroupsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.ProtectGroupIds = []*string{&protectGroupId}
	if protectGroupType != "" {
		request.ProtectGroupType = helper.String(protectGroupType)
	}

	var (
		offset int64 = 0
		limit  int64 = 100
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		var response *bdrcv20260330.DescribeDisasterRecoveryProtectGroupsResponse
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseBdrcV20260330Client().DescribeDisasterRecoveryProtectGroupsWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe bdrc disaster_recovery_protect_group failed, Response is nil."))
			}

			response = result
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		if response.Response.ProtectGroupSet == nil || len(response.Response.ProtectGroupSet) == 0 {
			return
		}

		for _, item := range response.Response.ProtectGroupSet {
			if item.ProtectGroupId != nil && *item.ProtectGroupId == protectGroupId {
				ret = item
				return
			}
		}

		if len(response.Response.ProtectGroupSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
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
