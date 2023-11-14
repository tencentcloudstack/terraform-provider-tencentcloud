/*
Use this data source to query detailed information of cdb error_log

Example Usage

```hcl
data "tencentcloud_cdb_error_log" "error_log" {
  instance_id = ""
  start_time =
  end_time =
  key_words =
  inst_type = ""
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

func dataSourceTencentCloudCdbErrorLog() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCdbErrorLogRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Start timestamp. For example 1585142640 .",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "End timestamp. For example 1585142640 .",
			},

			"key_words": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of keywords to match, up to 15 keywords are supported.",
			},

			"inst_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Only valid when the instance is the master instance or disaster recovery instance, the optional value: slave, which means to pull the log of the slave machine.",
			},

			"items": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "The records returned.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timestamp": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The time the error occurred.",
						},
						"content": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Error details.",
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

func dataSourceTencentCloudCdbErrorLogRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cdb_error_log.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("start_time"); v != nil {
		paramMap["StartTime"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("end_time"); v != nil {
		paramMap["EndTime"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("key_words"); ok {
		keyWordsSet := v.(*schema.Set).List()
		paramMap["KeyWords"] = helper.InterfacesStringsPoint(keyWordsSet)
	}

	if v, ok := d.GetOk("inst_type"); ok {
		paramMap["InstType"] = helper.String(v.(string))
	}

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCdbErrorLogByFilter(ctx, paramMap)
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
		for _, errlogItem := range items {
			errlogItemMap := map[string]interface{}{}

			if errlogItem.Timestamp != nil {
				errlogItemMap["timestamp"] = errlogItem.Timestamp
			}

			if errlogItem.Content != nil {
				errlogItemMap["content"] = errlogItem.Content
			}

			ids = append(ids, *errlogItem.InstanceId)
			tmpList = append(tmpList, errlogItemMap)
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
