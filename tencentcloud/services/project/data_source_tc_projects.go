package project

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tag "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudProjects() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudProjectsRead,
		Schema: map[string]*schema.Schema{
			"all_list": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "1 means to list all project, 0 means to list visible project.",
			},

			"projects": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of projects.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID of Project.",
						},
						"project_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of Project.",
						},
						"creator_uin": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Uin of Creator.",
						},
						"project_info": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of project.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time.",
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

func dataSourceTencentCloudProjectsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tag_project.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("all_list"); v != nil {
		paramMap["AllList"] = helper.IntUint64(v.(int))
	}

	service := TagService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var projects []*tag.Project

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeProjects(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		projects = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(projects))
	tmpList := make([]map[string]interface{}, 0, len(projects))

	if projects != nil {
		for _, project := range projects {
			projectMap := map[string]interface{}{}

			if project.ProjectId != nil {
				projectMap["project_id"] = project.ProjectId
			}

			if project.ProjectName != nil {
				projectMap["project_name"] = project.ProjectName
			}

			if project.CreatorUin != nil {
				projectMap["creator_uin"] = project.CreatorUin
			}

			if project.ProjectInfo != nil {
				projectMap["project_info"] = project.ProjectInfo
			}

			if project.CreateTime != nil {
				projectMap["create_time"] = project.CreateTime
			}

			ids = append(ids, helper.UInt64ToStr(*project.ProjectId))
			tmpList = append(tmpList, projectMap)
		}

		_ = d.Set("projects", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
