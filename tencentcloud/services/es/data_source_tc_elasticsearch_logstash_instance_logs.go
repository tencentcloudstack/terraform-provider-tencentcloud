package es

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	elasticsearch "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/es/v20180416"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudElasticsearchLogstashInstanceLogs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudElasticsearchLogstashInstanceLogsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"log_type": {
				Required: true,
				Type:     schema.TypeInt,
				Description: "Log type. Default 1, Valid values:\n" +
					" - 1: Main Log\n" +
					" - 2: Slow log\n" +
					" - 3: GC Log.",
			},

			"search_key": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Search terms, support LUCENE syntax, such as level:WARN, ip:1.1.1.1, message:test-index, etc.",
			},

			"start_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Log start time, in YYYY-MM-DD HH:MM:SS format, such as 2019-01-22 20:15:53.",
			},

			"end_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Log end time, in YYYY-MM-DD HH:MM:SS format, such as 2019-01-22 20:15:53.",
			},

			"order_by_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Time sort method. Default is 0. 0: descending; 1: ascending order.",
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
							Description: "Log content.",
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

func dataSourceTencentCloudElasticsearchLogstashInstanceLogsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_elasticsearch_logstash_instance_logs.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

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

	service := ElasticsearchService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var instanceLogList []*elasticsearch.InstanceLog

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeElasticsearchLogstashInstanceLogsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
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
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
