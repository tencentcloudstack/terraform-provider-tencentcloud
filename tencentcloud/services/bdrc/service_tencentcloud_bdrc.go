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

func (me *BdrcService) DescribeSecurityGroupMappingById(ctx context.Context, sitePairId, securityGroupMappingId string) (ret *bdrcv20260330.SecurityGroupMapping, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := bdrcv20260330.NewDescribeSecurityGroupMappingsRequest()
	request.SitePairId = helper.String(sitePairId)
	request.Limit = helper.Int64(100)

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
	request.Limit = helper.Int64(100)

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

func (me *BdrcService) RunCopyPairTasks(ctx context.Context, copyPairIds []*string, copyPairType string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := bdrcv20260330.NewRunCopyPairTasksRequest()
	request.CopyPairIds = copyPairIds
	request.CopyPairType = &copyPairType

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseBdrcV20260330Client().RunCopyPairTasksWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if err != nil {
		errRet = err
		return err
	}

	return
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
