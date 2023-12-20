package dcdb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDcdbProjects() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDcdbProjectsRead,
		Schema: map[string]*schema.Schema{
			"projects": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Project list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Project ID.",
						},
						"owner_uin": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The UIN of the resource owner (root account).",
						},
						"app_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Application ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Project name.",
						},
						"creator_uin": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Creator UIN.",
						},
						"src_plat": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Source platform.",
						},
						"src_app_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Source APPID.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Project status. Valid values: `0` (normal), `-1` (disabled), `3` (default project).",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
						},
						"is_default": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether it is the default project. Valid values: `1` (yes), `0` (no).",
						},
						"info": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description.",
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

func dataSourceTencentCloudDcdbProjectsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dcdb_projects.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := DcdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var projects []*dcdb.Project

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDcdbProjectsByFilter(ctx)
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
				ids = append(ids, helper.Int64ToStr(*project.ProjectId))
			}

			if project.OwnerUin != nil {
				projectMap["owner_uin"] = project.OwnerUin
			}

			if project.AppId != nil {
				projectMap["app_id"] = project.AppId
			}

			if project.Name != nil {
				projectMap["name"] = project.Name
			}

			if project.CreatorUin != nil {
				projectMap["creator_uin"] = project.CreatorUin
			}

			if project.SrcPlat != nil {
				projectMap["src_plat"] = project.SrcPlat
			}

			if project.SrcAppId != nil {
				projectMap["src_app_id"] = project.SrcAppId
			}

			if project.Status != nil {
				projectMap["status"] = project.Status
			}

			if project.CreateTime != nil {
				projectMap["create_time"] = project.CreateTime
			}

			if project.IsDefault != nil {
				projectMap["is_default"] = project.IsDefault
			}

			if project.Info != nil {
				projectMap["info"] = project.Info
			}

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
