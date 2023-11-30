package tencentcloud

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oceanus "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudOceanusTreeResources() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudOceanusTreeResourcesRead,
		Schema: map[string]*schema.Schema{
			"work_space_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Workspace SerialId.",
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
						"items": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of items.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"file_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "File name.",
									},
									"folder_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Folder id.",
									},
									"ref_job_status_count_set": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Counting the number of associated tasks by state.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"job_status": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Job status.",
												},
												"count": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Job count.",
												},
											},
										},
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name.",
									},
									"remark": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Remark.",
									},
									"resource_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Resource Id.",
									},
									"resource_type": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Resource Type.",
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

func dataSourceTencentCloudOceanusTreeResourcesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_oceanus_tree_resources.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId         = getLogId(contextNil)
		ctx           = context.WithValue(context.TODO(), logIdKey, logId)
		service       = OceanusService{client: meta.(*TencentCloudClient).apiV3Conn}
		treeResources *oceanus.DescribeTreeResourcesResponseParams
		workSpaceId   string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("work_space_id"); ok {
		paramMap["WorkSpaceId"] = helper.String(v.(string))
		workSpaceId = v.(string)
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeOceanusTreeResourcesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}

		treeResources = result
		return nil
	})

	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0)
	if treeResources != nil {
		treeResourceMap := map[string]interface{}{}
		if treeResources.Id != nil {
			treeResourceMap["id"] = treeResources.Id
		}

		if treeResources.Name != nil {
			treeResourceMap["name"] = treeResources.Name
		}

		if treeResources.ParentId != nil {
			treeResourceMap["parent_id"] = treeResources.ParentId
		}

		if treeResources.Items != nil {
			itemList := make([]map[string]interface{}, 0, len(treeResources.Items))
			for _, item := range treeResources.Items {
				itemMap := map[string]interface{}{}
				if item.FileName != nil {
					itemMap["file_name"] = item.FileName
				}

				if item.FolderId != nil {
					itemMap["folder_id"] = item.FolderId
				}

				if item.RefJobStatusCountSet != nil {
					jobList := make([]map[string]interface{}, 0, len(item.RefJobStatusCountSet))
					for _, job := range item.RefJobStatusCountSet {
						jobMap := map[string]interface{}{}
						if job.JobStatus != nil {
							jobMap["job_status"] = job.JobStatus
						}

						if job.Count != nil {
							jobMap["count"] = job.Count
						}

						jobList = append(jobList, jobMap)
					}

					itemMap["ref_job_status_count_set"] = jobList
				}

				if item.Name != nil {
					itemMap["name"] = item.Name
				}

				if item.Remark != nil {
					itemMap["remark"] = item.Remark
				}

				if item.ResourceId != nil {
					itemMap["resource_id"] = item.ResourceId
				}

				if item.ResourceType != nil {
					itemMap["resource_type"] = item.ResourceType
				}

				itemList = append(itemList, itemMap)
			}

			treeResourceMap["items"] = itemList
		}

		if treeResources.Children != nil {
			childrenBytes, _ := json.Marshal(treeResources.Children)
			childrenStr := string(childrenBytes)
			treeResourceMap["children"] = childrenStr
		}

		tmpList = append(tmpList, treeResourceMap)
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
