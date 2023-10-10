/*
Use this data source to query detailed information of cat metric_data
Example Usage
```hcl
data "tencentcloud_cat_metric_data" "metric_data" {
  analyze_task_type = "AnalyzeTaskType_Network"
  metric_type = "gauge"
  field = "avg(\"ping_time\")"
  filters = [
    "\"host\" = 'www.qq.com'",
    "time >= now()-1h",
  ]
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cat/v20180409"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCatMetricData() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCatMetricDataRead,
		Schema: map[string]*schema.Schema{
			"analyze_task_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Analysis of task type, supported types: `AnalyzeTaskType_Network`: network quality, `AnalyzeTaskType_Browse`: page performance, `AnalyzeTaskType_Transport`: port performance, `AnalyzeTaskType_UploadDownload`: file transport, `AnalyzeTaskType_MediaStream`: audiovisual experience.",
			},

			"metric_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Metric type, metrics queries are passed with gauge by default.",
			},

			"field": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Detailed fields of metrics, specified metrics can be passed or aggregate metrics, such as avg(ping_time) means entire delay.",
			},

			"filter": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Filter conditions can be passed as a single filter or multiple parameters concatenated together.",
			},

			"group_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Aggregation time, such as 1m, 1d, 30d, and so on.",
			},

			"filters": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Multiple condition filtering, supports combining multiple filtering conditions for query.",
			},

			"metric_set": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Return JSON string.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudCatMetricDataRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cat_metric_data.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("analyze_task_type"); ok {
		paramMap["AnalyzeTaskType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("metric_type"); ok {
		paramMap["MetricType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("field"); ok {
		paramMap["Field"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filter"); ok {
		paramMap["Filter"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_by"); ok {
		paramMap["GroupBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.(*schema.Set).List()
		paramMap["Filters"] = helper.InterfacesStringsPoint(filtersSet)
	}

	service := CatService{client: meta.(*TencentCloudClient).apiV3Conn}

	var metric *cat.DescribeProbeMetricDataResponseParams
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCatMetricDataByFilter(ctx, paramMap)
		if e != nil {
			if sdkError, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "FailedOperation.DbQueryFailed" {
					return resource.NonRetryableError(e)
				}
			}
			return retryError(e)
		}
		metric = result
		return nil
	})
	if err != nil {
		return err
	}

	var metricSet string
	if metric != nil && metric.MetricSet != nil {
		metricSet = *metric.MetricSet
		_ = d.Set("metric_set", metric.MetricSet)
	}

	d.SetId(helper.DataResourceIdsHash([]string{metricSet}))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), metricSet); e != nil {
			return e
		}
	}
	return nil
}
