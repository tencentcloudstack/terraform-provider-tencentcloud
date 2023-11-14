/*
Use this data source to query detailed information of waf attack_log_histogram

Example Usage

```hcl
data "tencentcloud_waf_attack_log_histogram" "attack_log_histogram" {
  domain = ""
  start_time = ""
  end_time = ""
  query_string = ""
    }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudWafAttackLogHistogram() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWafAttackLogHistogramRead,
		Schema: map[string]*schema.Schema{
			"domain": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Domain for query , all domain use &amp;amp;#39;all&amp;amp;#39;.",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Begin time.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "End time.",
			},

			"query_string": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Lucene grammar.",
			},

			"data": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "The statistics detail.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The count of logs.",
						},
						"time_stamp": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Timestamp.",
						},
					},
				},
			},

			"period": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Period.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudWafAttackLogHistogramRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_waf_attack_log_histogram.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("domain"); ok {
		paramMap["Domain"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("start_time"); ok {
		paramMap["StartTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["EndTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("query_string"); ok {
		paramMap["QueryString"] = helper.String(v.(string))
	}

	service := WafService{client: meta.(*TencentCloudClient).apiV3Conn}

	var data []*waf.LogHistogramInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWafAttackLogHistogramByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		data = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(data))
	tmpList := make([]map[string]interface{}, 0, len(data))

	if data != nil {
		for _, logHistogramInfo := range data {
			logHistogramInfoMap := map[string]interface{}{}

			if logHistogramInfo.Count != nil {
				logHistogramInfoMap["count"] = logHistogramInfo.Count
			}

			if logHistogramInfo.TimeStamp != nil {
				logHistogramInfoMap["time_stamp"] = logHistogramInfo.TimeStamp
			}

			ids = append(ids, *logHistogramInfo.RequestId)
			tmpList = append(tmpList, logHistogramInfoMap)
		}

		_ = d.Set("data", tmpList)
	}

	if period != nil {
		_ = d.Set("period", period)
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
