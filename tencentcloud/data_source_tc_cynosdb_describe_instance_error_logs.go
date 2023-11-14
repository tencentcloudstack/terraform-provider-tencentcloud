/*
Use this data source to query detailed information of cynosdb describe_instance_error_logs

Example Usage

```hcl
data "tencentcloud_cynosdb_describe_instance_error_logs" "describe_instance_error_logs" {
  instance_id = "cynosdbmysql-ins-4senc2fm"
  start_time = "2022-01-02 15:04:05"
  end_time = "2022-02-02 15:04:05"
  order_by = "Timestamp"
  order_by_type = "ASC"
  log_levels =
  key_words =
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

func dataSourceTencentCloudCynosdbDescribeInstanceErrorLogs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCynosdbDescribeInstanceErrorLogsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance Id.",
			},

			"start_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Start time.",
			},

			"end_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "End time.",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sort fields with Timestamp enumeration values.",
			},

			"order_by_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sort type, with ASC and DESC enumeration values.",
			},

			"log_levels": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Log levels, including error, warning, and note, support simultaneous search of multiple levels.",
			},

			"key_words": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Keywords, supports fuzzy search.",
			},

			"error_logs": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Error log list note: This field may return null, indicating that a valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timestamp": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Log timestamp note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"level": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Log level note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"content": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Note to log content: This field may return null, indicating that a valid value cannot be obtained.",
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

func dataSourceTencentCloudCynosdbDescribeInstanceErrorLogsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cynosdb_describe_instance_error_logs.read")()
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

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by_type"); ok {
		paramMap["OrderByType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("log_levels"); ok {
		logLevelsSet := v.(*schema.Set).List()
		paramMap["LogLevels"] = helper.InterfacesStringsPoint(logLevelsSet)
	}

	if v, ok := d.GetOk("key_words"); ok {
		keyWordsSet := v.(*schema.Set).List()
		paramMap["KeyWords"] = helper.InterfacesStringsPoint(keyWordsSet)
	}

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var errorLogs []*cynosdb.CynosdbErrorLogItem

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCynosdbDescribeInstanceErrorLogsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		errorLogs = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(errorLogs))
	tmpList := make([]map[string]interface{}, 0, len(errorLogs))

	if errorLogs != nil {
		for _, cynosdbErrorLogItem := range errorLogs {
			cynosdbErrorLogItemMap := map[string]interface{}{}

			if cynosdbErrorLogItem.Timestamp != nil {
				cynosdbErrorLogItemMap["timestamp"] = cynosdbErrorLogItem.Timestamp
			}

			if cynosdbErrorLogItem.Level != nil {
				cynosdbErrorLogItemMap["level"] = cynosdbErrorLogItem.Level
			}

			if cynosdbErrorLogItem.Content != nil {
				cynosdbErrorLogItemMap["content"] = cynosdbErrorLogItem.Content
			}

			ids = append(ids, *cynosdbErrorLogItem.Content)
			tmpList = append(tmpList, cynosdbErrorLogItemMap)
		}

		_ = d.Set("error_logs", tmpList)
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
