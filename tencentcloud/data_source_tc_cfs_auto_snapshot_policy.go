/*
Use this data source to query detailed information of cfs auto_snapshot_policy

Example Usage

```hcl
data "tencentcloud_cfs_auto_snapshot_policy" "auto_snapshot_policy" {
  auto_snapshot_policy_id = &lt;nil&gt;
  offset = &lt;nil&gt;
  limit = &lt;nil&gt;
  filters {
		values = &lt;nil&gt;
		name = &lt;nil&gt;

  }
  order = &lt;nil&gt;
  order_field = &lt;nil&gt;
  total_count = &lt;nil&gt;
  auto_snapshot_policies {
		auto_snapshot_policy_id = &lt;nil&gt;
		policy_name = &lt;nil&gt;
		creation_time = &lt;nil&gt;
		file_system_nums = &lt;nil&gt;
		day_of_week = &lt;nil&gt;
		hour = &lt;nil&gt;
		is_activated = &lt;nil&gt;
		next_active_time = &lt;nil&gt;
		status = &lt;nil&gt;
		app_id = &lt;nil&gt;
		alive_days = &lt;nil&gt;
		region_name = &lt;nil&gt;
		file_systems {
			creation_token = &lt;nil&gt;
			file_system_id = &lt;nil&gt;
			size_byte = &lt;nil&gt;
			storage_type = &lt;nil&gt;
			total_snapshot_size = &lt;nil&gt;
			creation_time = &lt;nil&gt;
			zone_id = &lt;nil&gt;
		}

  }
  request_id = &lt;nil&gt;
}
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfs/v20190719"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCfsAutoSnapshotPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCfsAutoSnapshotPolicyRead,
		Schema: map[string]*schema.Schema{
			"auto_snapshot_policy_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Snapshot policy ID.",
			},

			"offset": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Page offset.",
			},

			"limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Page length.",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Value.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name.",
						},
					},
				},
			},

			"order": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Ascending or descending order.",
			},

			"order_field": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting field.",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Description: "Total number of snapshot policies.",
			},

			"auto_snapshot_policies": {
				Type:        schema.TypeList,
				Description: "Snapshot policy information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_snapshot_policy_id": {
							Type:        schema.TypeString,
							Description: "Snapshot policy ID.",
						},
						"policy_name": {
							Type:        schema.TypeString,
							Description: "Snapshot policy name.",
						},
						"creation_time": {
							Type:        schema.TypeString,
							Description: "Snapshot policy creation time.",
						},
						"file_system_nums": {
							Type:        schema.TypeInt,
							Description: "Number of bound file systems.",
						},
						"day_of_week": {
							Type:        schema.TypeString,
							Description: "The day of the week on which to regularly back up the snapshot.",
						},
						"hour": {
							Type:        schema.TypeString,
							Description: "The hour of a day at which to regularly back up the snapshot.",
						},
						"is_activated": {
							Type:        schema.TypeInt,
							Description: "Whether to activate the scheduled snapshot feature.",
						},
						"next_active_time": {
							Type:        schema.TypeString,
							Description: "Next time to trigger snapshot.",
						},
						"status": {
							Type:        schema.TypeString,
							Description: "Snapshot policy status.",
						},
						"app_id": {
							Type:        schema.TypeInt,
							Description: "Account ID.",
						},
						"alive_days": {
							Type:        schema.TypeInt,
							Description: "Retention period.",
						},
						"region_name": {
							Type:        schema.TypeString,
							Description: "Region.",
						},
						"file_systems": {
							Type:        schema.TypeList,
							Description: "File system information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"creation_token": {
										Type:        schema.TypeString,
										Description: "File system name.",
									},
									"file_system_id": {
										Type:        schema.TypeString,
										Description: "File system ID.",
									},
									"size_byte": {
										Type:        schema.TypeInt,
										Description: "File system size.",
									},
									"storage_type": {
										Type:        schema.TypeString,
										Description: "Storage class.",
									},
									"total_snapshot_size": {
										Type:        schema.TypeInt,
										Description: "Total snapshot size.",
									},
									"creation_time": {
										Type:        schema.TypeString,
										Description: "File system creation time.",
									},
									"zone_id": {
										Type:        schema.TypeInt,
										Description: "Region ID of the file system.",
									},
								},
							},
						},
					},
				},
			},

			"request_id": {
				Type:        schema.TypeString,
				Description: "The unique request ID, which is returned for each request. RequestId is required for locating a problem.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudCfsAutoSnapshotPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cfs_auto_snapshot_policy.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("auto_snapshot_policy_id"); ok {
		paramMap["AutoSnapshotPolicyId"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["Offset"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["Limit"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*cfs.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := cfs.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["filters"] = tmpSet
	}

	if v, ok := d.GetOk("order"); ok {
		paramMap["Order"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_field"); ok {
		paramMap["OrderField"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("total_count"); v != nil {
		paramMap["TotalCount"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("auto_snapshot_policies"); ok {
		autoSnapshotPoliciesSet := v.([]interface{})
		tmpSet := make([]*cfs.AutoSnapshotPolicyInfo, 0, len(autoSnapshotPoliciesSet))

		for _, item := range autoSnapshotPoliciesSet {
			autoSnapshotPolicyInfo := cfs.AutoSnapshotPolicyInfo{}
			autoSnapshotPolicyInfoMap := item.(map[string]interface{})

			if v, ok := autoSnapshotPolicyInfoMap["auto_snapshot_policy_id"]; ok {
				autoSnapshotPolicyInfo.AutoSnapshotPolicyId = helper.String(v.(string))
			}
			if v, ok := autoSnapshotPolicyInfoMap["policy_name"]; ok {
				autoSnapshotPolicyInfo.PolicyName = helper.String(v.(string))
			}
			if v, ok := autoSnapshotPolicyInfoMap["creation_time"]; ok {
				autoSnapshotPolicyInfo.CreationTime = helper.String(v.(string))
			}
			if v, ok := autoSnapshotPolicyInfoMap["file_system_nums"]; ok {
				autoSnapshotPolicyInfo.FileSystemNums = helper.IntUint64(v.(int))
			}
			if v, ok := autoSnapshotPolicyInfoMap["day_of_week"]; ok {
				autoSnapshotPolicyInfo.DayOfWeek = helper.String(v.(string))
			}
			if v, ok := autoSnapshotPolicyInfoMap["hour"]; ok {
				autoSnapshotPolicyInfo.Hour = helper.String(v.(string))
			}
			if v, ok := autoSnapshotPolicyInfoMap["is_activated"]; ok {
				autoSnapshotPolicyInfo.IsActivated = helper.IntUint64(v.(int))
			}
			if v, ok := autoSnapshotPolicyInfoMap["next_active_time"]; ok {
				autoSnapshotPolicyInfo.NextActiveTime = helper.String(v.(string))
			}
			if v, ok := autoSnapshotPolicyInfoMap["status"]; ok {
				autoSnapshotPolicyInfo.Status = helper.String(v.(string))
			}
			if v, ok := autoSnapshotPolicyInfoMap["app_id"]; ok {
				autoSnapshotPolicyInfo.AppId = helper.IntUint64(v.(int))
			}
			if v, ok := autoSnapshotPolicyInfoMap["alive_days"]; ok {
				autoSnapshotPolicyInfo.AliveDays = helper.IntUint64(v.(int))
			}
			if v, ok := autoSnapshotPolicyInfoMap["region_name"]; ok {
				autoSnapshotPolicyInfo.RegionName = helper.String(v.(string))
			}
			if v, ok := autoSnapshotPolicyInfoMap["file_systems"]; ok {
				for _, item := range v.([]interface{}) {
					fileSystemsMap := item.(map[string]interface{})
					fileSystemByPolicy := cfs.FileSystemByPolicy{}
					if v, ok := fileSystemsMap["creation_token"]; ok {
						fileSystemByPolicy.CreationToken = helper.String(v.(string))
					}
					if v, ok := fileSystemsMap["file_system_id"]; ok {
						fileSystemByPolicy.FileSystemId = helper.String(v.(string))
					}
					if v, ok := fileSystemsMap["size_byte"]; ok {
						fileSystemByPolicy.SizeByte = helper.IntUint64(v.(int))
					}
					if v, ok := fileSystemsMap["storage_type"]; ok {
						fileSystemByPolicy.StorageType = helper.String(v.(string))
					}
					if v, ok := fileSystemsMap["total_snapshot_size"]; ok {
						fileSystemByPolicy.TotalSnapshotSize = helper.IntUint64(v.(int))
					}
					if v, ok := fileSystemsMap["creation_time"]; ok {
						fileSystemByPolicy.CreationTime = helper.String(v.(string))
					}
					if v, ok := fileSystemsMap["zone_id"]; ok {
						fileSystemByPolicy.ZoneId = helper.IntUint64(v.(int))
					}
					autoSnapshotPolicyInfo.FileSystems = append(autoSnapshotPolicyInfo.FileSystems, &fileSystemByPolicy)
				}
			}
			tmpSet = append(tmpSet, &autoSnapshotPolicyInfo)
		}
		paramMap["auto_snapshot_policies"] = tmpSet
	}

	if v, ok := d.GetOk("request_id"); ok {
		paramMap["RequestId"] = helper.String(v.(string))
	}

	service := CfsService{client: meta.(*TencentCloudClient).apiV3Conn}

	var autoSnapshotPolicies []*cfs.AutoSnapshotPolicyInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCfsAutoSnapshotPolicyByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		autoSnapshotPolicies = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(autoSnapshotPolicies))
	tmpList := make([]map[string]interface{}, 0, len(autoSnapshotPolicies))

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
