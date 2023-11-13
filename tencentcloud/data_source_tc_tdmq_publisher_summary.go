/*
Use this data source to query detailed information of tdmq publisher_summary

Example Usage

```hcl
data "tencentcloud_tdmq_publisher_summary" "publisher_summary" {
  cluster_id = ""
  namespace = ""
  topic = ""
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

func dataSourceTencentCloudTdmqPublisherSummary() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTdmqPublisherSummaryRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},

			"namespace": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Namespace name.",
			},

			"topic": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Subject name.",
			},

			"msg_rate_in": {
				Computed:    true,
				Type:        schema.TypeFloat,
				Description: "Production rate (units per second)Note: This field may return null, indicating that no valid value can be obtained.",
			},

			"msg_throughput_in": {
				Computed:    true,
				Type:        schema.TypeFloat,
				Description: "Production rate (bytes per second)Note: This field may return null, indicating that no valid value can be obtained.",
			},

			"publisher_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Number of producersNote: This field may return null, indicating that no valid value can be obtained.",
			},

			"storage_size": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Message store size in bytesNote: This field may return null, indicating that no valid value can be obtained.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudTdmqPublisherSummaryRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tdmq_publisher_summary.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["ClusterId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace"); ok {
		paramMap["Namespace"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("topic"); ok {
		paramMap["Topic"] = helper.String(v.(string))
	}

	service := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTdmqPublisherSummaryByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		msgRateIn = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(msgRateIn))
	if msgRateIn != nil {
		_ = d.Set("msg_rate_in", msgRateIn)
	}

	if msgThroughputIn != nil {
		_ = d.Set("msg_throughput_in", msgThroughputIn)
	}

	if publisherCount != nil {
		_ = d.Set("publisher_count", publisherCount)
	}

	if storageSize != nil {
		_ = d.Set("storage_size", storageSize)
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
