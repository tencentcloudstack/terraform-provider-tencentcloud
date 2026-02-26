package tke

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tkev20180525 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudKubernetesUpgradeTasks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudKubernetesUpgradeTasksRead,
		Schema: map[string]*schema.Schema{
			"upgrade_tasks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Upgrade tasks.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Task ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Task name.",
						},
						"component": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Component name.",
						},
						"related_resources": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "Related resources.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"upgrade_impact": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Upgrade impact.",
						},
						"planed_start_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Planned start time.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
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

func dataSourceTencentCloudKubernetesUpgradeTasksRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_kubernetes_upgrade_tasks.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	var respData []*tkev20180525.UpgradeTask
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKubernetesUpgradeTasksByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	upgradeTasksList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, upgradeTasks := range respData {
			upgradeTasksMap := map[string]interface{}{}
			if upgradeTasks.ID != nil {
				upgradeTasksMap["id"] = upgradeTasks.ID
			}

			if upgradeTasks.Name != nil {
				upgradeTasksMap["name"] = upgradeTasks.Name
			}

			if upgradeTasks.Component != nil {
				upgradeTasksMap["component"] = upgradeTasks.Component
			}

			if upgradeTasks.RelatedResources != nil {
				upgradeTasksMap["related_resources"] = upgradeTasks.RelatedResources
			}

			if upgradeTasks.UpgradeImpact != nil {
				upgradeTasksMap["upgrade_impact"] = upgradeTasks.UpgradeImpact
			}

			if upgradeTasks.PlanedStartAt != nil {
				upgradeTasksMap["planed_start_at"] = upgradeTasks.PlanedStartAt
			}

			if upgradeTasks.CreatedAt != nil {
				upgradeTasksMap["created_at"] = upgradeTasks.CreatedAt
			}

			upgradeTasksList = append(upgradeTasksList, upgradeTasksMap)
		}

		_ = d.Set("upgrade_tasks", upgradeTasksList)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), upgradeTasksList); e != nil {
			return e
		}
	}

	return nil
}
