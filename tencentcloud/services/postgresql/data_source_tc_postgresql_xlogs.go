package postgresql

import (
	"context"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudPostgresqlXlogs() *schema.Resource {
	return &schema.Resource{
		Read: datasourceTencentCloudPostgresqlXlogsRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "PostgreSQL instance id.",
			},
			"start_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Xlog start time, format `yyyy-MM-dd hh:mm:ss`, start time cannot before 7 days ago.",
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Xlog end time, format `yyyy-MM-dd hh:mm:ss`.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used for save results.",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of Xlog query result.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Xlog id.",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Xlog file created start time.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Xlog file created end time.",
						},
						"internal_addr": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Xlog internal download address.",
						},
						"external_addr": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Xlog external download address.",
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Xlog file size.",
						},
					},
				},
			},
		},
	}
}

func datasourceTencentCloudPostgresqlXlogsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("datasource.tencentcloud_postgresql_xlogs.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := PostgresqlService{client}

	request := postgresql.NewDescribeDBXlogsRequest()

	id := d.Get("instance_id").(string)
	defaultEndTime := time.Now()

	if endTime, ok := d.GetOk("end_time"); ok && endTime != "" {
		request.EndTime = helper.String(endTime.(string))
	} else {
		endTime := defaultEndTime.Format("2006-01-02 15:04:05")
		request.EndTime = &endTime
	}

	if startTime, ok := d.GetOk("start_time"); ok && startTime != "" {
		request.StartTime = helper.String(startTime.(string))
	} else {
		defaultStartTime := defaultEndTime.AddDate(0, 0, -7)
		startTime := defaultStartTime.Format("2006-01-02 15:04:05")
		request.StartTime = &startTime
	}

	request.DBInstanceId = &id

	result, err := service.DescribeDBXlogs(ctx, request)

	if err != nil {
		d.SetId("")
		return err
	}

	list := make([]interface{}, 0, len(result))

	for i := range result {
		item := result[i]
		xlog := map[string]interface{}{
			"id":            item.Id,
			"start_time":    item.StartTime,
			"end_time":      item.EndTime,
			"internal_addr": item.InternalAddr,
			"external_addr": item.ExternalAddr,
			"size":          item.Size,
		}

		list = append(list, xlog)
	}

	d.SetId("postgres-xlog-" + id)

	if err := d.Set("list", list); err != nil {
		return err
	}

	if output, ok := d.GetOk("result_output_file"); ok {
		return tccommon.WriteToFile(output.(string), list)
	}

	return nil
}
