/*
Use this data source to query detailed information of rum project

Example Usage

```hcl
data "tencentcloud_rum_project" "project" {
	instance_id = "rum-pasZKEI3RLgakj"
}
```
*/
package tencentcloud

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	rum "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/rum/v20210622"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudRumProject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudRumProjectRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance ID.",
			},

			"project_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Project list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Project name.",
						},
						"creator": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creator ID.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance ID.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Project type.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CreateTime.",
						},
						"repo": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Project repository address.",
						},
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Project URL.",
						},
						"rate": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Project sample rate.",
						},
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique project key (12 characters).",
						},
						"enable_url_group": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether to enable URL aggregation.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance name.",
						},
						"pid": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Project ID.",
						},
						"instance_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance key.",
						},
						"desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Project description.",
						},
						"is_star": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Starred status. `1`: yes; `0`: no.",
						},
						"project_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Project status (`1`: Creating; `2`: Running; `3`: Abnormal; `4`: Restarting; `5`: Stopping; `6`: Stopped; `7`: Terminating; `8`: Terminated).",
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

func dataSourceTencentCloudRumProjectRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_rum_project.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["instance_id"] = helper.String(v.(string))
	}

	rumService := RumService{client: meta.(*TencentCloudClient).apiV3Conn}

	var projectSet []*rum.RumProject
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := rumService.DescribeRumProjectByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		projectSet = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Rum projectSet failed, reason:%+v", logId, err)
		return err
	}

	projectSetList := []interface{}{}
	ids := make([]string, 0, len(projectSet))
	if projectSet != nil {
		for _, projectSet := range projectSet {
			ids = append(ids, strconv.FormatInt(*projectSet.ID, 10))

			projectSetMap := map[string]interface{}{}
			if projectSet.Name != nil {
				projectSetMap["name"] = projectSet.Name
			}
			if projectSet.Creator != nil {
				projectSetMap["creator"] = projectSet.Creator
			}
			if projectSet.InstanceID != nil {
				projectSetMap["instance_id"] = projectSet.InstanceID
			}
			if projectSet.Type != nil {
				projectSetMap["type"] = projectSet.Type
			}
			if projectSet.CreateTime != nil {
				projectSetMap["create_time"] = projectSet.CreateTime
			}
			if projectSet.Repo != nil {
				projectSetMap["repo"] = projectSet.Repo
			}
			if projectSet.URL != nil {
				projectSetMap["url"] = projectSet.URL
			}
			if projectSet.Rate != nil {
				projectSetMap["rate"] = projectSet.Rate
			}
			if projectSet.Key != nil {
				projectSetMap["key"] = projectSet.Key
			}
			if projectSet.EnableURLGroup != nil {
				projectSetMap["enable_url_group"] = projectSet.EnableURLGroup
			}
			if projectSet.InstanceName != nil {
				projectSetMap["instance_name"] = projectSet.InstanceName
			}
			if projectSet.ID != nil {
				projectSetMap["pid"] = projectSet.ID
			}
			if projectSet.InstanceKey != nil {
				projectSetMap["instance_key"] = projectSet.InstanceKey
			}
			if projectSet.Desc != nil {
				projectSetMap["desc"] = projectSet.Desc
			}
			if projectSet.IsStar != nil {
				projectSetMap["is_star"] = projectSet.IsStar
			}
			if projectSet.ProjectStatus != nil {
				projectSetMap["project_status"] = projectSet.ProjectStatus
			}

			projectSetList = append(projectSetList, projectSetMap)
		}
		_ = d.Set("project_set", projectSetList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), projectSetList); e != nil {
			return e
		}
	}

	return nil
}
