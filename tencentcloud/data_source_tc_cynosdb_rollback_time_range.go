/*
Use this data source to query detailed information of cynosdb rollback_time_range

Example Usage

```hcl
data "tencentcloud_cynosdb_rollback_time_range" "rollback_time_range" {
  cluster_id = "cynosdbmysql-bws8h88b"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCynosdbRollbackTimeRange() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCynosdbRollbackTimeRangeRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},
			"time_range_start": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Effective regression time range start time point (obsolete) Note: This field may return null, indicating that a valid value cannot be obtained.",
			},
			"time_range_end": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Effective regression time range end time point (obsolete) Note: This field may return null, indicating that a valid value cannot be obtained.",
			},
			"rollback_time_ranges": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Reversible time range.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time_range_start": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "start time.",
						},
						"time_range_end": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "End time.",
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

func dataSourceTencentCloudCynosdbRollbackTimeRangeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cynosdb_rollback_time_range.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		service           = CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
		rollbackTimeRange *cynosdb.DescribeRollbackTimeRangeResponseParams
		clusterId         string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["ClusterId"] = helper.String(v.(string))
		clusterId = v.(string)
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCynosdbRollbackTimeRangeByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}

		rollbackTimeRange = result
		return nil
	})

	if err != nil {
		return err
	}

	if rollbackTimeRange.TimeRangeStart != nil {
		_ = d.Set("time_range_start", rollbackTimeRange.TimeRangeStart)
	}

	if rollbackTimeRange.TimeRangeEnd != nil {
		_ = d.Set("time_range_end", rollbackTimeRange.TimeRangeEnd)
	}

	if rollbackTimeRange.RollbackTimeRanges != nil {
		tmpList := []interface{}{}
		for _, timeRange := range rollbackTimeRange.RollbackTimeRanges {
			rollbackTimeRangeMap := map[string]interface{}{}

			if timeRange.TimeRangeStart != nil {
				rollbackTimeRangeMap["time_range_start"] = timeRange.TimeRangeStart
			}

			if timeRange.TimeRangeEnd != nil {
				rollbackTimeRangeMap["time_range_end"] = timeRange.TimeRangeEnd
			}

			tmpList = append(tmpList, rollbackTimeRangeMap)
		}

		_ = d.Set("rollback_time_ranges", tmpList)
	}

	d.SetId(clusterId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
