/*
Use this data source to query detailed information of es elasticsearch_instance_logs

Example Usage

```hcl
data "tencentcloud_elasticsearch_instance_logs" "elasticsearch_instance_logs" {
	instance_id = "es-xxxxxx"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	es "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/es/v20180416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudElasticsearchInstanceLogs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudElasticsearchInstanceLogsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"log_type": {
				Optional: true,
				Type:     schema.TypeInt,
				Description: "Log type. Log type, default is 1, Valid values:\n" +
					"- 1: master log\n" +
					"- 2: Search slow log\n" +
					"- 3: Index slow log\n" +
					"- 4: GC log.",
			},

			"search_key": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Search key. Support LUCENE syntax, such as level:WARN, ip:1.1.1.1, message:test-index, etc.",
			},

			"start_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Start time. The format is YYYY-MM-DD HH:MM:SS, such as 2019-01-22 20:15:53.",
			},

			"end_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "End time. The format is YYYY-MM-DD HH:MM:SS, such as 2019-01-22 20:15:53.",
			},

			"order_by_type": {
				Optional: true,
				Type:     schema.TypeInt,
				Description: "Order type. Time sort method. Default is 0, valid values:\n" +
					"- 0: descending;\n" +
					"- 1: ascending order.",
			},

			"instance_log_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of log details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Log time.",
						},
						"level": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Log level.",
						},
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster node ip.",
						},
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Log message.",
						},
						"node_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster node id.",
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

func dataSourceTencentCloudElasticsearchInstanceLogsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_elasticsearch_instance_logs.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var instanceId string
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		paramMap["InstanceId"] = helper.String(instanceId)
	}

	if v, ok := d.GetOkExists("log_type"); ok {
		paramMap["LogType"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("search_key"); ok {
		paramMap["SearchKey"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("start_time"); ok {
		paramMap["StartTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["EndTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("order_by_type"); ok {
		paramMap["OrderByType"] = helper.IntUint64(v.(int))
	}

	service := ElasticsearchService{client: meta.(*TencentCloudClient).apiV3Conn}

	var instanceLogList []*es.InstanceLog

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeElasticsearchInstanceLogsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		instanceLogList = result
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(instanceLogList))

	if instanceLogList != nil {
		for _, instanceLog := range instanceLogList {
			instanceLogMap := map[string]interface{}{}

			if instanceLog.Time != nil {
				instanceLogMap["time"] = instanceLog.Time
			}

			if instanceLog.Level != nil {
				instanceLogMap["level"] = instanceLog.Level
			}

			if instanceLog.Ip != nil {
				instanceLogMap["ip"] = instanceLog.Ip
			}

			if instanceLog.Message != nil {
				instanceLogMap["message"] = instanceLog.Message
			}

			if instanceLog.NodeID != nil {
				instanceLogMap["node_id"] = instanceLog.NodeID
			}

			tmpList = append(tmpList, instanceLogMap)
		}

		_ = d.Set("instance_log_list", tmpList)
	}

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
