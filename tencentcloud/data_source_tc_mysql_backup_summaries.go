/*
Use this data source to query detailed information of mysql backup_summaries

Example Usage

```hcl
data "tencentcloud_mysql_backup_summaries" "backup_summaries" {
  product = "mysql"
  order_by = "BackupVolume"
  order_direction = "ASC"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudMysqlBackupSummaries() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMysqlBackupSummariesRead,
		Schema: map[string]*schema.Schema{
			"product": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The type of cloud database product to be queried, currently only supports `mysql`.",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Specify to sort by a certain item, the optional values include: BackupVolume: backup volume, DataBackupVolume: data backup volume, BinlogBackupVolume: log backup volume, AutoBackupVolume: automatic backup volume, ManualBackupVolume: manual backup volume. By default, they are sorted by BackupVolume.",
			},

			"order_direction": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Specify the sorting direction, optional values include: ASC: forward order, DESC: reverse order. The default is ASC.",
			},

			"items": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Instance backup statistics entries.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance ID.",
						},
						"auto_backup_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of automatic data backups for this instance.",
						},
						"auto_backup_volume": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The automatic data backup capacity of this instance.",
						},
						"manual_backup_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of manual data backups for this instance.",
						},
						"manual_backup_volume": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The capacity of manual data backup for this instance.",
						},
						"data_backup_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total number of data backups (including automatic backups and manual backups) of the instance.",
						},
						"data_backup_volume": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total data backup capacity of this instance.",
						},
						"binlog_backup_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of log backups for this instance.",
						},
						"binlog_backup_volume": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The capacity of the instance log backup.",
						},
						"backup_volume": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total backup (including data backup and log backup) of the instance occupies capacity.",
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

func dataSourceTencentCloudMysqlBackupSummariesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mysql_backup_summaries.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("product"); ok {
		paramMap["Product"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_direction"); ok {
		paramMap["OrderDirection"] = helper.String(v.(string))
	}

	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	var backupSummaries []*cdb.BackupSummaryItem
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMysqlBackupSummariesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		backupSummaries = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(backupSummaries))
	tmpList := make([]map[string]interface{}, 0, len(backupSummaries))

	if backupSummaries != nil {
		for _, backupSummaryItem := range backupSummaries {
			backupSummaryItemMap := map[string]interface{}{}

			if backupSummaryItem.InstanceId != nil {
				backupSummaryItemMap["instance_id"] = backupSummaryItem.InstanceId
			}

			if backupSummaryItem.AutoBackupCount != nil {
				backupSummaryItemMap["auto_backup_count"] = backupSummaryItem.AutoBackupCount
			}

			if backupSummaryItem.AutoBackupVolume != nil {
				backupSummaryItemMap["auto_backup_volume"] = backupSummaryItem.AutoBackupVolume
			}

			if backupSummaryItem.ManualBackupCount != nil {
				backupSummaryItemMap["manual_backup_count"] = backupSummaryItem.ManualBackupCount
			}

			if backupSummaryItem.ManualBackupVolume != nil {
				backupSummaryItemMap["manual_backup_volume"] = backupSummaryItem.ManualBackupVolume
			}

			if backupSummaryItem.DataBackupCount != nil {
				backupSummaryItemMap["data_backup_count"] = backupSummaryItem.DataBackupCount
			}

			if backupSummaryItem.DataBackupVolume != nil {
				backupSummaryItemMap["data_backup_volume"] = backupSummaryItem.DataBackupVolume
			}

			if backupSummaryItem.BinlogBackupCount != nil {
				backupSummaryItemMap["binlog_backup_count"] = backupSummaryItem.BinlogBackupCount
			}

			if backupSummaryItem.BinlogBackupVolume != nil {
				backupSummaryItemMap["binlog_backup_volume"] = backupSummaryItem.BinlogBackupVolume
			}

			if backupSummaryItem.BackupVolume != nil {
				backupSummaryItemMap["backup_volume"] = backupSummaryItem.BackupVolume
			}

			ids = append(ids, *backupSummaryItem.InstanceId)
			tmpList = append(tmpList, backupSummaryItemMap)
		}

		_ = d.Set("items", tmpList)
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
