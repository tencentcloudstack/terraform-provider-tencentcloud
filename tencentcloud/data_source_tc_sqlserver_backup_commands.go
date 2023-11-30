package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudSqlserverBackupCommands() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSqlserverBackupCommandsRead,
		Schema: map[string]*schema.Schema{
			"backup_file_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Backup file type. Full: full backup. FULL_LOG: full backup which needs log increments. FULL_DIFF: full backup which needs differential increments. LOG: log backup. DIFF: differential backup.",
			},
			"data_base_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Database name.",
			},
			"is_recovery": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Whether restoration is required. No: not required. Yes: required.",
			},
			"local_path": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Storage path of backup files. If this parameter is left empty, the default storage path will be D:.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Command list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"command": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create backup command.",
						},
						"request_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Request ID.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudSqlserverBackupCommandsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_sqlserver_backup_commands.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		command []*sqlserver.DescribeBackupCommandResponseParams
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("backup_file_type"); ok {
		paramMap["BackupFileType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("data_base_name"); ok {
		paramMap["DataBaseName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("is_recovery"); ok {
		paramMap["IsRecovery"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("local_path"); ok {
		paramMap["LocalPath"] = helper.String(v.(string))
	}

	service := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, err := service.DescribeSqlserverBackupCommand(ctx, paramMap)
		if err != nil {
			return retryError(err)
		}

		command = result
		return nil
	})

	if err != nil {
		return err
	}

	var list []map[string]interface{}
	var ids = make([]string, len(command))
	for _, item := range command {
		mapping := map[string]interface{}{
			"command":    item.Command,
			"request_id": item.RequestId,
		}
		list = append(list, mapping)
		ids = append(ids, fmt.Sprintf("%s%s", *item.Command, *item.RequestId))
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("list", list); e != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s\n", logId, e.Error())
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), command); e != nil {
			return e
		}
	}

	return nil
}
