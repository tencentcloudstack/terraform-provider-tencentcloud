package wedata

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudWedataProjects() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWedataProjectsRead,
		Schema: map[string]*schema.Schema{
			"project_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of project IDs.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Project name or unique identifier name, supports fuzzy search.",
			},

			"status": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Project status, optional values: 0 (disabled), 1 (normal).",
			},

			"project_model": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Project model, optional values: SIMPLE, STANDARD.",
			},

			"items": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of data sources.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Project ID.",
						},
						"project_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Project identifier, English name.",
						},
						"display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Project display name, can be Chinese name.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Remarks.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
						},
						"creator_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Project creator ID.",
						},
						"project_owner_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Project owner ID.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Project status: 0: disabled, 1: enabled, -3: disabling, 2: enabling.",
						},
						"project_model": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Project model, SIMPLE: simple mode, STANDARD: standard mode.",
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

func dataSourceTencentCloudWedataProjectsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_wedata_projects.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("project_ids"); ok {
		projectIdsList := []*string{}
		projectIdsSet := v.(*schema.Set).List()
		for i := range projectIdsSet {
			projectIds := projectIdsSet[i].(string)
			projectIdsList = append(projectIdsList, helper.String(projectIds))
		}

		paramMap["ProjectIds"] = projectIdsList
	}

	if v, ok := d.GetOk("project_name"); ok {
		paramMap["ProjectName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("status"); ok {
		paramMap["Status"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("project_model"); ok {
		paramMap["ProjectModel"] = helper.String(v.(string))
	}

	var respData []*wedatav20250806.Project
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWedataProjectsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	itemsList := make([]map[string]interface{}, 0, len(respData))
	for _, items := range respData {
		itemsMap := map[string]interface{}{}
		if items.ProjectId != nil {
			itemsMap["project_id"] = items.ProjectId
		}

		if items.ProjectName != nil {
			itemsMap["project_name"] = items.ProjectName
		}

		if items.DisplayName != nil {
			itemsMap["display_name"] = items.DisplayName
		}

		if items.Description != nil {
			itemsMap["description"] = items.Description
		}

		if items.CreateTime != nil {
			itemsMap["create_time"] = items.CreateTime
		}

		if items.CreatorUin != nil {
			itemsMap["creator_uin"] = items.CreatorUin
		}

		if items.ProjectOwnerUin != nil {
			itemsMap["project_owner_uin"] = items.ProjectOwnerUin
		}

		if items.Status != nil {
			itemsMap["status"] = items.Status
		}

		if items.ProjectModel != nil {
			itemsMap["project_model"] = items.ProjectModel
		}

		itemsList = append(itemsList, itemsMap)
	}

	_ = d.Set("items", itemsList)

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), itemsList); e != nil {
			return e
		}
	}

	return nil
}
