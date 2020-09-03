/*
Use this data source to query the list of SQL Server backups.

Example Usage

```hcl
data "tencentcloud_sqlserver_backups" "foo" {
  instance_id           = "mssql-3cdq7kx5"
  start_time         = "2020-06-17 00:00:00"
  end_time			= "2020-06-22 00:00:00"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudSqlserverBackups() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceTencentSqlserverBackupsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance ID.",
			},
			"start_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Start time of the instance list, like yyyy-MM-dd HH:mm:ss.",
			},
			"end_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "End time of the instance list, like yyyy-MM-dd HH:mm:ss.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to store results.",
			},
			// Computed values
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of SQL Server backup. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the backup.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance ID.",
						},
						"file_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "File name of the backup.",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Start time of the backup.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "End time of the backup.",
						},
						"db_list": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "Database name list of the backup.",
						},
						"strategy": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Strategy of the backup. 0 for instance backup, 1 for multi-databases backup.",
						},
						"trigger_model": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The way to trigger backup. 0 for timed trigger, 1 for manual trigger.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Status of the backup. 1 for creating, 2 for successfully created, 3 for failed.",
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The size of backup file. Unit is KB.",
						},
						"intranet_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "URL for downloads internally.",
						},
						"internet_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "URL for downloads externally.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentSqlserverBackupsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_sqlserver_backups.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	instanceId := d.Get("instance_id").(string)
	startTime := d.Get("start_time").(string)
	endTime := d.Get("end_time").(string)
	sqlserverService := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	backInfoItems, err := sqlserverService.DescribeSqlserverBackups(ctx, instanceId, startTime, endTime)

	if err != nil {
		return fmt.Errorf("api[DescribeBackups]fail, return %s", err.Error())
	}

	var list []map[string]interface{}
	var ids = make([]string, len(backInfoItems))

	for _, item := range backInfoItems {
		mapping := map[string]interface{}{
			"start_time":    item.StartTime,
			"end_time":      item.EndTime,
			"size":          item.Size,
			"trigger_model": item.BackupWay,
			"intranet_url":  item.InternalAddr,
			"internet_url":  item.ExternalAddr,
			"status":        item.Status,
			"file_name":     item.FileName,
			"instance_id":   instanceId,
			"id":            strconv.Itoa(int(*item.Id)),
			"db_list":       helper.StringsInterfaces(item.DBs),
		}
		list = append(list, mapping)
		ids = append(ids, fmt.Sprintf("%d", *item.Id))
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("list", list); e != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s\n", logId, e.Error())
		return e
	}
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		return writeToFile(output.(string), list)
	}

	return nil
}
