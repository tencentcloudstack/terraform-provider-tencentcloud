package wedata

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudWedataTaskVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWedataTaskVersionsRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Project ID.",
			},

			"task_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Task ID.",
			},

			"task_version_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SAVE version.\nSUBMIT version.\nDefaults to SAVE.",
			},

			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Task version list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"create_time": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Creation time.",
						},
						"version_num": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Version number.",
						},
						"create_user_uin": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Creator ID.",
						},
						"version_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Saved version ID.",
						},
						"version_remark": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Version description.",
						},
						"approve_status": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Approval status (only for submit version).",
						},
						"status": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Production status (only for submit version).",
						},
						"approve_user_uin": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Approver (only for submit version).",
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

func dataSourceTencentCloudWedataTaskVersionsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_wedata_task_versions.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(nil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("project_id"); ok {
		paramMap["ProjectId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("task_id"); ok {
		paramMap["TaskId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("task_version_type"); ok {
		paramMap["TaskVersionType"] = helper.String(v.(string))
	}

	var respData []*wedatav20250806.TaskVersion
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWedataTaskVersionsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		respData = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(respData))
	itemsList := make([]map[string]interface{}, 0, len(respData))

	for _, items := range respData {
		itemsMap := map[string]interface{}{}

		if items.CreateTime != nil {
			itemsMap["create_time"] = items.CreateTime
		}

		if items.VersionNum != nil {
			itemsMap["version_num"] = items.VersionNum
		}

		if items.CreateUserUin != nil {
			itemsMap["create_user_uin"] = items.CreateUserUin
		}

		if items.VersionId != nil {
			itemsMap["version_id"] = items.VersionId
		}

		if items.VersionRemark != nil {
			itemsMap["version_remark"] = items.VersionRemark
		}

		if items.ApproveStatus != nil {
			itemsMap["approve_status"] = items.ApproveStatus
		}

		if items.Status != nil {
			itemsMap["status"] = items.Status
		}

		if items.ApproveUserUin != nil {
			itemsMap["approve_user_uin"] = items.ApproveUserUin
		}

		itemsList = append(itemsList, itemsMap)
	}

	_ = d.Set("data", itemsList)

	d.SetId(helper.DataResourceIdsHash(ids))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), itemsList); e != nil {
			return e
		}
	}

	return nil
}
