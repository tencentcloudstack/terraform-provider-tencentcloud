package sqlserver

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudSqlserverRollbackTime() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSqlserverRollbackTimeRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
			"dbs": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of databases to be queried.",
			},
			"details": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Information of time range available for database rollback.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database name.",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Start time of time range available for rollback.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "End time of time range available for rollback.",
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

func dataSourceTencentCloudSqlserverRollbackTimeRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_sqlserver_rollback_time.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		details    []*sqlserver.DbRollbackTimeInfo
		instanceId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("dbs"); ok {
		dBsSet := v.(*schema.Set).List()
		paramMap["DBs"] = helper.InterfacesStringsPoint(dBsSet)
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSqlserverRollbackTimeByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		details = result
		return nil
	})

	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(details))

	if details != nil {
		for _, dbRollbackTimeInfo := range details {
			dbRollbackTimeInfoMap := map[string]interface{}{}

			if dbRollbackTimeInfo.DBName != nil {
				dbRollbackTimeInfoMap["db_name"] = dbRollbackTimeInfo.DBName
			}

			if dbRollbackTimeInfo.StartTime != nil {
				dbRollbackTimeInfoMap["start_time"] = dbRollbackTimeInfo.StartTime
			}

			if dbRollbackTimeInfo.EndTime != nil {
				dbRollbackTimeInfoMap["end_time"] = dbRollbackTimeInfo.EndTime
			}

			tmpList = append(tmpList, dbRollbackTimeInfoMap)
		}

		_ = d.Set("details", tmpList)
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
