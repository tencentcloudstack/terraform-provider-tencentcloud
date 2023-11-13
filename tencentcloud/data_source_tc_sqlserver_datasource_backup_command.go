/*
Use this data source to query detailed information of sqlserver datasource_backup_command

Example Usage

```hcl
data "tencentcloud_sqlserver_datasource_backup_command" "datasource_backup_command" {
  backup_file_type = "FULL"
  data_base_name = "db_name"
  is_recovery = "No"
  local_path = ""
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

func dataSourceTencentCloudSqlserverDatasourceBackupCommand() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSqlserverDatasourceBackupCommandRead,
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

			"command": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Create backup command.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudSqlserverDatasourceBackupCommandRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_sqlserver_datasource_backup_command.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

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
		result, e := service.DescribeSqlserverDatasourceBackupCommandByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		command = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(command))
	if command != nil {
		_ = d.Set("command", command)
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
