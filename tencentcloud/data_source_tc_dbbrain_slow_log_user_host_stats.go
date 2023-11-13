/*
Use this data source to query detailed information of dbbrain slow_log_user_host_stats

Example Usage

```hcl
data "tencentcloud_dbbrain_slow_log_user_host_stats" "slow_log_user_host_stats" {
  instance_id = ""
  start_time = ""
  end_time = ""
  product = ""
  md5 = ""
    }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDbbrainSlowLogUserHostStats() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainSlowLogUserHostStatsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Start time of the query range, time format such as：2019-09-10 12:13:14。.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "EndTime time of the query range, time format such as：2019-09-10 12:13:14。.",
			},

			"product": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Types of service products, supported values：mysql - Cloud Database MySQL; cynosdb - Cloud Database TDSQL-C for MySQL, defaults to mysql.",
			},

			"md5": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "MD5 value of SOL template.",
			},

			"total_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Number of source addresses.",
			},

			"items": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Detailed list of the slow log proportion for each source address.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Source address.",
						},
						"ratio": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "The ratio of the number of slow logs of the source address to the total, in %.",
						},
						"count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of slow logs for this source address.",
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

func dataSourceTencentCloudDbbrainSlowLogUserHostStatsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dbbrain_slow_log_user_host_stats.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("start_time"); ok {
		paramMap["StartTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["EndTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product"); ok {
		paramMap["Product"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("md5"); ok {
		paramMap["Md5"] = helper.String(v.(string))
	}

	service := DbbrainService{client: meta.(*TencentCloudClient).apiV3Conn}

	var items []*dbbrain.SlowLogHost

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDbbrainSlowLogUserHostStatsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		items = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(items))
	tmpList := make([]map[string]interface{}, 0, len(items))

	if totalCount != nil {
		_ = d.Set("total_count", totalCount)
	}

	if items != nil {
		for _, slowLogHost := range items {
			slowLogHostMap := map[string]interface{}{}

			if slowLogHost.UserHost != nil {
				slowLogHostMap["user_host"] = slowLogHost.UserHost
			}

			if slowLogHost.Ratio != nil {
				slowLogHostMap["ratio"] = slowLogHost.Ratio
			}

			if slowLogHost.Count != nil {
				slowLogHostMap["count"] = slowLogHost.Count
			}

			ids = append(ids, *slowLogHost.InstanceId)
			tmpList = append(tmpList, slowLogHostMap)
		}

		_ = d.Set("items", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
