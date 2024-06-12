package cvm

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

func dataSourceTencentCloudPlacementGroupsReadPostHandleResponse0(ctx context.Context, req map[string]interface{}, resp *[]*cvm.DisasterRecoverGroup) error {
	d := tccommon.ResourceDataFromContext(ctx)
	logId := tccommon.GetLogId(tccommon.ContextNil)
	placementGroups := *resp
	var err error
	placementGroupList := make([]map[string]interface{}, 0, len(placementGroups))
	ids := make([]string, 0, len(placementGroups))
	for _, placement := range placementGroups {
		mapping := map[string]interface{}{
			"placement_group_id": placement.DisasterRecoverGroupId,
			"name":               placement.Name,
			"type":               placement.Type,
			"cvm_quota_total":    placement.CvmQuotaTotal,
			"current_num":        placement.CurrentNum,
			"instance_ids":       helper.StringsInterfaces(placement.InstanceIds),
			"create_time":        placement.CreateTime,
		}
		placementGroupList = append(placementGroupList, mapping)
		ids = append(ids, *placement.DisasterRecoverGroupId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("placement_group_list", placementGroupList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set placement group list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	context.WithValue(ctx, "placementGroupList", placementGroupList)
	return nil
}

func dataSourceTencentCloudPlacementGroupsReadOutputContent(ctx context.Context) interface{} {
	eipList := ctx.Value("placementGroupList")
	return eipList
}
