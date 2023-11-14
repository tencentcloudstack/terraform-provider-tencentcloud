/*
Use this data source to query detailed information of cdb slow_log

Example Usage

```hcl
data "tencentcloud_cdb_slow_log" "slow_log" {
  instance_id = ""
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

func dataSourceTencentCloudCdbSlowLog() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCdbSlowLogRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID, in the format: cdb-c1nl9rpv. Same instance ID as displayed in the ApsaraDB for Console page.",
			},

			"items": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Details of slow query logs that meet the query conditions.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup file name.",
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Backup file size, unit: Byte.",
						},
						"date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup snapshot time, time format: 2016-03-17 02:10:37.",
						},
						"intranet_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Intranet download address.",
						},
						"internet_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "External network download address.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Log specific type, possible values: slowlog - slow log.",
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

func dataSourceTencentCloudCdbSlowLogRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cdb_slow_log.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCdbSlowLogByFilter(ctx, paramMap)
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
		for _, slowLogInfo := range items {
			slowLogInfoMap := map[string]interface{}{}

			if slowLogInfo.Name != nil {
				slowLogInfoMap["name"] = slowLogInfo.Name
			}

			if slowLogInfo.Size != nil {
				slowLogInfoMap["size"] = slowLogInfo.Size
			}

			if slowLogInfo.Date != nil {
				slowLogInfoMap["date"] = slowLogInfo.Date
			}

			if slowLogInfo.IntranetUrl != nil {
				slowLogInfoMap["intranet_url"] = slowLogInfo.IntranetUrl
			}

			if slowLogInfo.InternetUrl != nil {
				slowLogInfoMap["internet_url"] = slowLogInfo.InternetUrl
			}

			if slowLogInfo.Type != nil {
				slowLogInfoMap["type"] = slowLogInfo.Type
			}

			ids = append(ids, *slowLogInfo.InstanceId)
			tmpList = append(tmpList, slowLogInfoMap)
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
