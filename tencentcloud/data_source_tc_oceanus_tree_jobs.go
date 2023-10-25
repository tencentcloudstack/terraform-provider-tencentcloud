/*
Use this data source to query detailed information of oceanus tree_jobs

Example Usage

```hcl
data "tencentcloud_oceanus_tree_jobs" "example" {
  work_space_id = "space-2idq8wbr"
}
```
*/
package tencentcloud

import (
	"context"
	"encoding/json"

	oceanus "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudOceanusTreeJobs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudOceanusTreeJobsRead,
		Schema: map[string]*schema.Schema{
			"work_space_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Workspace SerialId.",
			},
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter rules.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Field to be filtered. Can only be set `Zone` or `JobType` or `JobStatus`.",
						},
						"values": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							Description: "Filter values for the field.",
						},
					},
				},
			},
			"tree_info": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Tree structure information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name.",
						},
						"parent_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parent Id.",
						},
						"job_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of jobs.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"job_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Job ID.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Job Name.",
									},
									"job_type": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Job Type.",
									},
									"running_cu": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Resources occupied by homework.",
									},
									"status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Job status.",
									},
								},
							},
						},
						"children": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subdirectory Information.",
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

func dataSourceTencentCloudOceanusTreeJobsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_oceanus_tree_jobs.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId       = getLogId(contextNil)
		ctx         = context.WithValue(context.TODO(), logIdKey, logId)
		service     = OceanusService{client: meta.(*TencentCloudClient).apiV3Conn}
		treeJobs    *oceanus.DescribeTreeJobsResponseParams
		workSpaceId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("work_space_id"); ok {
		paramMap["WorkSpaceId"] = helper.String(v.(string))
		workSpaceId = v.(string)
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*oceanus.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := oceanus.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}

			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}

			tmpSet = append(tmpSet, &filter)
		}

		paramMap["Filters"] = tmpSet
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeOceanusTreeJobsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}

		treeJobs = result
		return nil
	})

	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0)
	if treeJobs != nil {
		treeJobsMap := map[string]interface{}{}

		if treeJobs.Id != nil {
			treeJobsMap["id"] = treeJobs.Id
		}

		if treeJobs.Name != nil {
			treeJobsMap["name"] = treeJobs.Name
		}

		if treeJobs.ParentId != nil {
			treeJobsMap["parent_id"] = treeJobs.ParentId
		}

		if treeJobs.JobSet != nil {
			jobList := make([]map[string]interface{}, 0, len(treeJobs.JobSet))
			for _, item := range treeJobs.JobSet {
				jobMap := map[string]interface{}{}
				if item.JobId != nil {
					jobMap["file_name"] = item.JobId
				}

				if item.Name != nil {
					jobMap["folder_id"] = item.Name
				}

				if item.JobType != nil {
					jobMap["name"] = item.JobType
				}

				if item.RunningCu != nil {
					jobMap["remark"] = item.RunningCu
				}

				if item.Status != nil {
					jobMap["resource_id"] = item.Status
				}

				jobList = append(jobList, jobMap)
			}

			treeJobsMap["job_set"] = jobList
		}

		if treeJobs.Children != nil {
			childrenBytes, _ := json.Marshal(treeJobs.Children)
			childrenStr := string(childrenBytes)
			treeJobsMap["children"] = childrenStr
		}

		tmpList = append(tmpList, treeJobsMap)
		_ = d.Set("tree_info", tmpList)
	}

	d.SetId(workSpaceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
