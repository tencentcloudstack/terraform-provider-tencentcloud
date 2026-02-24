package tke

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tkev20180525 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudKubernetesUpgradeTaskDetail() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudKubernetesUpgradeTaskDetailRead,
		Schema: map[string]*schema.Schema{
			"task_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Upgrade task ID.",
			},

			"upgrade_plans": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Upgrade plans.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Upgrade plan ID.",
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster ID.",
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster name.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster region.",
						},
						"planed_start_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Planned start time.",
						},
						"upgrade_start_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Upgrade start time.",
						},
						"upgrade_end_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Upgrade end time.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Upgrade status.",
						},
						"reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Reason.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudKubernetesUpgradeTaskDetailRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_kubernetes_upgrade_task_detail.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		taskId  string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOkExists("task_id"); ok {
		paramMap["ID"] = helper.IntInt64(v.(int))
		taskId = helper.IntToStr(v.(int))
	}

	var respData []*tkev20180525.UpgradePlan
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKubernetesUpgradeTaskDetailByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	upgradePlansList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, upgradePlans := range respData {
			upgradePlansMap := map[string]interface{}{}
			if upgradePlans.ID != nil {
				upgradePlansMap["id"] = upgradePlans.ID
			}

			if upgradePlans.ClusterID != nil {
				upgradePlansMap["cluster_id"] = upgradePlans.ClusterID
			}

			if upgradePlans.ClusterName != nil {
				upgradePlansMap["cluster_name"] = upgradePlans.ClusterName
			}

			if upgradePlans.Region != nil {
				upgradePlansMap["region"] = upgradePlans.Region
			}

			if upgradePlans.PlanedStartAt != nil {
				upgradePlansMap["planed_start_at"] = upgradePlans.PlanedStartAt
			}

			if upgradePlans.UpgradeStartAt != nil {
				upgradePlansMap["upgrade_start_at"] = upgradePlans.UpgradeStartAt
			}

			if upgradePlans.UpgradeEndAt != nil {
				upgradePlansMap["upgrade_end_at"] = upgradePlans.UpgradeEndAt
			}

			if upgradePlans.Status != nil {
				upgradePlansMap["status"] = upgradePlans.Status
			}

			if upgradePlans.Reason != nil {
				upgradePlansMap["reason"] = upgradePlans.Reason
			}

			upgradePlansList = append(upgradePlansList, upgradePlansMap)
		}

		_ = d.Set("upgrade_plans", upgradePlansList)
	}

	d.SetId(taskId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), upgradePlansList); e != nil {
			return e
		}
	}

	return nil
}
