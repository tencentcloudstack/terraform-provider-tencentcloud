package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCynosdbDescribeInstanceSlowQueries() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCynosdbDescribeInstanceSlowQueriesRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},
			"start_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "start time.",
			},
			"end_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "End time.",
			},
			"binlogs": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Note to the Binlog list: This field may return null, indicating that a valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"file_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Binlog file name.",
						},
						"file_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "File size in bytes.",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Earliest transaction time.",
						},
						"finish_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Latest transaction time.",
						},
						"binlog_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Binlog file ID.",
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

func dataSourceTencentCloudCynosdbDescribeInstanceSlowQueriesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cynosdb_describe_instance_slow_queries.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId     = getLogId(contextNil)
		ctx       = context.WithValue(context.TODO(), logIdKey, logId)
		service   = CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
		binlogs   []*cynosdb.BinlogItem
		clusterId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["ClusterId"] = helper.String(v.(string))
		clusterId = v.(string)
	}

	if v, ok := d.GetOk("start_time"); ok {
		paramMap["StartTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["EndTime"] = helper.String(v.(string))
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCynosdbDescribeInstanceSlowQueriesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}

		binlogs = result
		return nil
	})

	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(binlogs))

	if binlogs != nil {
		for _, binlogItem := range binlogs {
			binlogItemMap := map[string]interface{}{}

			if binlogItem.FileName != nil {
				binlogItemMap["file_name"] = binlogItem.FileName
			}

			if binlogItem.FileSize != nil {
				binlogItemMap["file_size"] = binlogItem.FileSize
			}

			if binlogItem.StartTime != nil {
				binlogItemMap["start_time"] = binlogItem.StartTime
			}

			if binlogItem.FinishTime != nil {
				binlogItemMap["finish_time"] = binlogItem.FinishTime
			}

			if binlogItem.BinlogId != nil {
				binlogItemMap["binlog_id"] = binlogItem.BinlogId
			}

			tmpList = append(tmpList, binlogItemMap)
		}

		_ = d.Set("binlogs", tmpList)
	}

	d.SetId(clusterId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
