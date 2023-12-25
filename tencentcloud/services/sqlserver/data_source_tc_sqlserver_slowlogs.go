package sqlserver

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudSqlserverSlowlogs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSqlserverSlowlogsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Query start time.",
			},
			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Query end time.",
			},
			"slowlogs": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Information list of slow query logs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Unique ID of slow query log file.",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "File generation start time.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "File generation end time.",
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "File size in KB.",
						},
						"count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of logs in file.",
						},
						"internal_addr": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Download address for private network.",
						},
						"external_addr": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Download address for public network.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Status (1: success, 2: failure) Note: this field may return null, indicating that no valid values can be obtained.",
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

func dataSourceTencentCloudSqlserverSlowlogsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_sqlserver_slowlogs.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("start_time"); ok {
		paramMap["StartTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["EndTime"] = helper.String(v.(string))
	}

	var slowlogs []*sqlserver.SlowlogInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSqlserverSlowlogsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		slowlogs = result
		return nil
	})

	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(slowlogs))

	if slowlogs != nil {
		for _, slowlogInfo := range slowlogs {
			slowlogInfoMap := map[string]interface{}{}

			if slowlogInfo.Id != nil {
				slowlogInfoMap["id"] = slowlogInfo.Id
			}

			if slowlogInfo.StartTime != nil {
				slowlogInfoMap["start_time"] = slowlogInfo.StartTime
			}

			if slowlogInfo.EndTime != nil {
				slowlogInfoMap["end_time"] = slowlogInfo.EndTime
			}

			if slowlogInfo.Size != nil {
				slowlogInfoMap["size"] = slowlogInfo.Size
			}

			if slowlogInfo.Count != nil {
				slowlogInfoMap["count"] = slowlogInfo.Count
			}

			if slowlogInfo.InternalAddr != nil {
				slowlogInfoMap["internal_addr"] = slowlogInfo.InternalAddr
			}

			if slowlogInfo.ExternalAddr != nil {
				slowlogInfoMap["external_addr"] = slowlogInfo.ExternalAddr
			}

			if slowlogInfo.Status != nil {
				slowlogInfoMap["status"] = slowlogInfo.Status
			}

			tmpList = append(tmpList, slowlogInfoMap)
		}

		_ = d.Set("slowlogs", tmpList)
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
