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

func (me *BdrcService) DescribeSecurityGroupMappingById(ctx context.Context, sitePairId, securityGroupMappingId string) (ret *bdrcv20260330.SecurityGroupMapping, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := bdrcv20260330.NewDescribeSecurityGroupMappingsRequest()
	request.SitePairId = helper.String(sitePairId)
	request.Limit = helper.Int64(500)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseBdrcV20260330Client().DescribeSecurityGroupMappingsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe bdrc security_group_mapping failed, Response is nil."))
		}

		if result.Response.SecurityGroupMappingSet == nil || len(result.Response.SecurityGroupMappingSet) == 0 {
			return nil
		}

		for _, item := range result.Response.SecurityGroupMappingSet {
			if item.SecurityGroupMappingId != nil && *item.SecurityGroupMappingId == securityGroupMappingId {
				ret = item
				return nil
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

func (me *BdrcService) DescribeSecurityGroupMappingByFilter(ctx context.Context, sitePairId, srcSecurityGroupId, targetSecurityGroupId string) (ret *bdrcv20260330.SecurityGroupMapping, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := bdrcv20260330.NewDescribeSecurityGroupMappingsRequest()
	request.SitePairId = helper.String(sitePairId)
	request.Limit = helper.Int64(500)

	srcFilter := &bdrcv20260330.FilterModel{
		Name:   helper.String("src-security-group-id"),
		Values: []*string{helper.String(srcSecurityGroupId)},
	}
	targetFilter := &bdrcv20260330.FilterModel{
		Name:   helper.String("target-security-group-id"),
		Values: []*string{helper.String(targetSecurityGroupId)},
	}
	request.Filters = []*bdrcv20260330.FilterModel{srcFilter, targetFilter}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseBdrcV20260330Client().DescribeSecurityGroupMappingsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe bdrc security_group_mapping by filter failed, Response is nil."))
		}

		if result.Response.SecurityGroupMappingSet == nil || len(result.Response.SecurityGroupMappingSet) == 0 {
			return resource.RetryableError(fmt.Errorf("bdrc security_group_mapping not found by filter, please retry."))
		}

		for _, item := range result.Response.SecurityGroupMappingSet {
			if item.SourceSecurityGroupId != nil && *item.SourceSecurityGroupId == srcSecurityGroupId &&
				item.TargetSecurityGroupId != nil && *item.TargetSecurityGroupId == targetSecurityGroupId {
				ret = item
				return nil
			}
		}

		return resource.RetryableError(fmt.Errorf("bdrc security_group_mapping not matched by src/target id, please retry."))
	})

	if err != nil {
		errRet = err
		return
	}

	return
}
