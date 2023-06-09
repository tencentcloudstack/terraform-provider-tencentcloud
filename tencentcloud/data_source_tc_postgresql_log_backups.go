/*
Use this data source to query detailed information of postgresql log_backups

Example Usage

```hcl
data "tencentcloud_postgresql_log_backups" "log_backups" {
  min_finish_time = "%s"
  max_finish_time = "%s"
  filters {
		name = "db-instance-id"
		values = [local.pgsql_id]
  }
  order_by = "StartTime"
  order_by_type = "desc"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudPostgresqlLogBackups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudPostgresqlLogBackupsRead,
		Schema: map[string]*schema.Schema{
			"min_finish_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Minimum end time of a backup in the format of `2018-01-01 00:00:00`. It is 7 days ago by default.",
			},

			"max_finish_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Maximum end time of a backup in the format of `2018-01-01 00:00:00`. It is the current time by default.",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter instances using one or more criteria. Valid filter names:db-instance-id: Filter by instance ID (in string format).db-instance-name: Filter by instance name (in string format).db-instance-ip: Filter by instance VPC IP (in string format).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Filter name.",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "One or more filter values.",
						},
					},
				},
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting field. Valid values: `StartTime`, `FinishTime`, `Size`.",
			},

			"order_by_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting order. Valid values: `asc` (ascending), `desc` (descending).",
			},

			"log_backup_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of log backup details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance ID.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique ID of a backup file.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup file name.",
						},
						"backup_method": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup method, including physical and logical.",
						},
						"backup_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup mode, including automatic and manual.",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup task status.",
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Backup set size in bytes.",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup start time.",
						},
						"finish_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup end time.",
						},
						"expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup expiration time.",
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

func dataSourceTencentCloudPostgresqlLogBackupsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_postgresql_log_backups.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("min_finish_time"); ok {
		paramMap["MinFinishTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("max_finish_time"); ok {
		paramMap["MaxFinishTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*postgresql.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := postgresql.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}
			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["filters"] = tmpSet
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by_type"); ok {
		paramMap["OrderByType"] = helper.String(v.(string))
	}

	service := PostgresqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	var logBackupSet []*postgresql.LogBackup

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribePostgresqlLogBackupsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		logBackupSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(logBackupSet))
	tmpList := make([]map[string]interface{}, 0, len(logBackupSet))

	if logBackupSet != nil {
		for _, logBackup := range logBackupSet {
			logBackupMap := map[string]interface{}{}

			if logBackup.DBInstanceId != nil {
				logBackupMap["db_instance_id"] = logBackup.DBInstanceId
			}

			if logBackup.Id != nil {
				logBackupMap["id"] = logBackup.Id
			}

			if logBackup.Name != nil {
				logBackupMap["name"] = logBackup.Name
			}

			if logBackup.BackupMethod != nil {
				logBackupMap["backup_method"] = logBackup.BackupMethod
			}

			if logBackup.BackupMode != nil {
				logBackupMap["backup_mode"] = logBackup.BackupMode
			}

			if logBackup.State != nil {
				logBackupMap["state"] = logBackup.State
			}

			if logBackup.Size != nil {
				logBackupMap["size"] = logBackup.Size
			}

			if logBackup.StartTime != nil {
				logBackupMap["start_time"] = logBackup.StartTime
			}

			if logBackup.FinishTime != nil {
				logBackupMap["finish_time"] = logBackup.FinishTime
			}

			if logBackup.ExpireTime != nil {
				logBackupMap["expire_time"] = logBackup.ExpireTime
			}

			ids = append(ids, *logBackup.DBInstanceId)
			tmpList = append(tmpList, logBackupMap)
		}

		_ = d.Set("log_backup_set", tmpList)
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
