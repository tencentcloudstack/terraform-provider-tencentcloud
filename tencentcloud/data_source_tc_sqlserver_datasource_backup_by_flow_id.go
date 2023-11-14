/*
Use this data source to query detailed information of sqlserver datasource_backup_by_flow_id

Example Usage

```hcl
data "tencentcloud_sqlserver_datasource_backup_by_flow_id" "datasource_backup_by_flow_id" {
  instance_id = "mssql-i1z41iwd"
  flow_id = ""
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

func dataSourceTencentCloudSqlserverDatasourceBackupByFlowId() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSqlserverDatasourceBackupByFlowIdRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"flow_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Create a backup process ID, which can be obtained through the [CreateBackup](https://cloud.tencent.com/document/product/238/19946) interface.",
			},

			"file_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "File name. For a single-database backup file, only the file name of the first record is returned; for a single-database backup file, the file names of all records need to be obtained through the DescribeBackupFiles interface.",
			},

			"backup_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Backup task name, customizable.",
			},

			"start_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Backup start time.",
			},

			"end_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Backup end time.",
			},

			"strategy": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Backup strategy, 0-instance backup; 1-multi-database backup; when the instance status is 0-creating, this field is the default value 0, meaningless.",
			},

			"status": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Backup file status, 0-creating; 1-success; 2-failure.",
			},

			"backup_way": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Backup method, 0-scheduled backup; 1-manual temporary backup; instance status is 0-creating, this field is the default value 0, meaningless.",
			},

			"d_bs": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "For the DB list, only the library name contained in the first record is returned for a single-database backup file; for a single-database backup file, the library names of all records need to be obtained through the DescribeBackupFiles interface.",
			},

			"internal_addr": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Intranet download address, for a single database backup file, only the intranet download address of the first record is returned; single database backup files need to obtain the download addresses of all records through the DescribeBackupFiles interface.",
			},

			"external_addr": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "External network download address, for a single database backup file, only the external network download address of the first record is returned; single database backup files need to obtain the download addresses of all records through the DescribeBackupFiles interface.",
			},

			"group_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Aggregate Id, this value is not returned for packaged backup files. Use this value to call the DescribeBackupFiles interface to obtain the detailed information of a single database backup file.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudSqlserverDatasourceBackupByFlowIdRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_sqlserver_datasource_backup_by_flow_id.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("flow_id"); ok {
		paramMap["FlowId"] = helper.String(v.(string))
	}

	service := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSqlserverDatasourceBackupByFlowIdByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		fileName = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(fileName))
	if fileName != nil {
		_ = d.Set("file_name", fileName)
	}

	if backupName != nil {
		_ = d.Set("backup_name", backupName)
	}

	if startTime != nil {
		_ = d.Set("start_time", startTime)
	}

	if endTime != nil {
		_ = d.Set("end_time", endTime)
	}

	if strategy != nil {
		_ = d.Set("strategy", strategy)
	}

	if status != nil {
		_ = d.Set("status", status)
	}

	if backupWay != nil {
		_ = d.Set("backup_way", backupWay)
	}

	if dBs != nil {
		_ = d.Set("d_bs", dBs)
	}

	if internalAddr != nil {
		_ = d.Set("internal_addr", internalAddr)
	}

	if externalAddr != nil {
		_ = d.Set("external_addr", externalAddr)
	}

	if groupId != nil {
		_ = d.Set("group_id", groupId)
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
