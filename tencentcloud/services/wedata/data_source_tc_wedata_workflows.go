package wedata

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudWedataWorkflows() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWedataWorkflowsRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Project ID.",
			},

			"keyword": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Search keywords.",
			},

			"parent_folder_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Workflow folder.",
			},

			"workflow_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Workflow type. valid values: cycle and manual.",
			},

			"bundle_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "bundleId item.",
			},

			"owner_uin": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Owner ID.",
			},

			"create_user_uin": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Creator ID.",
			},

			"modify_time": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Modification time interval yyyy-MM-dd HH:MM:ss. fill in two times in the array.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"create_time": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Creation time range yyyy-MM-dd HH:MM:ss. two times must be filled in the array.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Describes workflow pagination information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"workflow_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Workflow ID.",
						},
						"workflow_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Workflow name.",
						},
						"workflow_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Workflow type: cycle or manual.",
						},
						"owner_uin": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Owner ID.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Creation time.",
						},
						"modify_time": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Last Modification Time.",
						},
						"update_user_uin": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Last updated user ID.",
						},
						"workflow_desc": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Workflow description.",
						},
						"create_user_uin": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Creator ID.",
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

func dataSourceTencentCloudWedataWorkflowsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_wedata_workflows.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(nil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("project_id"); ok {
		paramMap["ProjectId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("keyword"); ok {
		paramMap["Keyword"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("parent_folder_path"); ok {
		paramMap["ParentFolderPath"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("workflow_type"); ok {
		paramMap["WorkflowType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("bundle_id"); ok {
		paramMap["BundleId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("owner_uin"); ok {
		paramMap["OwnerUin"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("create_user_uin"); ok {
		paramMap["CreateUserUin"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("modify_time"); ok {
		modifyTimeList := []*string{}
		modifyTimeSet := v.(*schema.Set).List()
		for i := range modifyTimeSet {
			modifyTime := modifyTimeSet[i].(string)
			modifyTimeList = append(modifyTimeList, helper.String(modifyTime))
		}
		paramMap["ModifyTime"] = modifyTimeList
	}

	if v, ok := d.GetOk("create_time"); ok {
		createTimeList := []*string{}
		createTimeSet := v.(*schema.Set).List()
		for i := range createTimeSet {
			createTime := createTimeSet[i].(string)
			createTimeList = append(createTimeList, helper.String(createTime))
		}
		paramMap["CreateTime"] = createTimeList
	}

	var respData []*wedatav20250806.WorkflowInfo
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWedataWorkflowsByFilter(ctx, paramMap)
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

		if items.WorkflowId != nil {
			itemsMap["workflow_id"] = items.WorkflowId
			ids = append(ids, *items.WorkflowId)
		}

		if items.WorkflowName != nil {
			itemsMap["workflow_name"] = items.WorkflowName
		}

		if items.WorkflowType != nil {
			itemsMap["workflow_type"] = items.WorkflowType
		}

		if items.OwnerUin != nil {
			itemsMap["owner_uin"] = items.OwnerUin
		}

		if items.CreateTime != nil {
			itemsMap["create_time"] = items.CreateTime
		}

		if items.ModifyTime != nil {
			itemsMap["modify_time"] = items.ModifyTime
		}

		if items.UpdateUserUin != nil {
			itemsMap["update_user_uin"] = items.UpdateUserUin
		}

		if items.WorkflowDesc != nil {
			itemsMap["workflow_desc"] = items.WorkflowDesc
		}

		if items.CreateUserUin != nil {
			itemsMap["create_user_uin"] = items.CreateUserUin
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
