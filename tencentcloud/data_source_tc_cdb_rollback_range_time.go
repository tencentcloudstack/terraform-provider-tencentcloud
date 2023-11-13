/*
Use this data source to query detailed information of cdb rollback_range_time

Example Usage

```hcl
data "tencentcloud_cdb_rollback_range_time" "rollback_range_time" {
  instance_ids =
  is_remote_zone = ""
  backup_region = ""
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCdbRollbackRangeTime() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCdbRollbackRangeTimeRead,
		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of instance IDs, the format of a single instance ID is: cdb-c1nl9rpv. Same instance ID as displayed in the ApsaraDB for Console page.",
			},

			"is_remote_zone": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Whether the clone instance is in the same zone as the source instance, yes: `false`, no: `true`.",
			},

			"backup_region": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "If the clone instance is not in the same region as the source instance, fill in the region where the clone instance is located, for example: `ap-guangzhou`.",
			},

			"items": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Returned parameter information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Query database error code.",
						},
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Query database error information.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A list of instance IDs. The format of a single instance ID is: cdb-c1nl9rpv. Same as the instance ID displayed in the cloud database console page.",
						},
						"times": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Retrievable time range.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"begin": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance rollback start time, time format: 2016-10-29 01:06:04.",
									},
									"end": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "End time of instance rollback, time format: 2016-11-02 11:44:47.",
									},
								},
							},
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

func dataSourceTencentCloudCdbRollbackRangeTimeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cdb_rollback_range_time.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		paramMap["InstanceIds"] = helper.InterfacesStringsPoint(instanceIdsSet)
	}

	if v, ok := d.GetOk("is_remote_zone"); ok {
		paramMap["IsRemoteZone"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("backup_region"); ok {
		paramMap["BackupRegion"] = helper.String(v.(string))
	}

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCdbRollbackRangeTimeByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		totalCount = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(totalCount))
	if items != nil {
		for _, instanceRollbackRangeTime := range items {
			instanceRollbackRangeTimeMap := map[string]interface{}{}

			if instanceRollbackRangeTime.Code != nil {
				instanceRollbackRangeTimeMap["code"] = instanceRollbackRangeTime.Code
			}

			if instanceRollbackRangeTime.Message != nil {
				instanceRollbackRangeTimeMap["message"] = instanceRollbackRangeTime.Message
			}

			if instanceRollbackRangeTime.InstanceId != nil {
				instanceRollbackRangeTimeMap["instance_id"] = instanceRollbackRangeTime.InstanceId
			}

			if instanceRollbackRangeTime.Times != nil {
				timesList := []interface{}{}
				for _, times := range instanceRollbackRangeTime.Times {
					timesMap := map[string]interface{}{}

					if times.Begin != nil {
						timesMap["begin"] = times.Begin
					}

					if times.End != nil {
						timesMap["end"] = times.End
					}

					timesList = append(timesList, timesMap)
				}

				instanceRollbackRangeTimeMap["times"] = []interface{}{timesList}
			}

			ids = append(ids, *instanceRollbackRangeTime.IdsHash)
			tmpList = append(tmpList, instanceRollbackRangeTimeMap)
		}

		_ = d.Set("items", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string)); e != nil {
			return e
		}
	}
	return nil
}
