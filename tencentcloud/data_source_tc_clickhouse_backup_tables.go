/*
Use this data source to query detailed information of clickhouse backup tables

Example Usage

```hcl
data "tencentcloud_clickhouse_backup_tables" "backup_tables" {
  instance_id = "cdwch-xxxxxx"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clickhouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwch/v20200915"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudClickhouseBackupTables() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClickhouseBackupTablesRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"available_tables": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Available tables.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database.",
						},
						"table": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Table.",
						},
						"total_bytes": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Table total bytes.",
						},
						"v_cluster": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Virtual cluster.",
						},
						"ips": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Table ips.",
						},
						"zoo_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Zk path.",
						},
						"rip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Ip address of cvm.",
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

func dataSourceTencentCloudClickhouseBackupTablesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_clickhouse_backup_tables.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	instanceId := d.Get("instance_id").(string)

	service := CdwchService{client: meta.(*TencentCloudClient).apiV3Conn}

	var availableTables []*clickhouse.BackupTableContent

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeClickhouseBackupTablesByFilter(ctx, instanceId)
		if e != nil {
			return retryError(e)
		}
		availableTables = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(availableTables))
	tmpList := make([]map[string]interface{}, 0, len(availableTables))

	if availableTables != nil {
		for _, backupTableContent := range availableTables {
			backupTableContentMap := map[string]interface{}{}

			if backupTableContent.Database != nil {
				backupTableContentMap["database"] = backupTableContent.Database
			}

			if backupTableContent.Table != nil {
				backupTableContentMap["table"] = backupTableContent.Table
			}

			if backupTableContent.TotalBytes != nil {
				backupTableContentMap["total_bytes"] = backupTableContent.TotalBytes
			}

			if backupTableContent.VCluster != nil {
				backupTableContentMap["v_cluster"] = backupTableContent.VCluster
			}

			if backupTableContent.Ips != nil {
				backupTableContentMap["ips"] = backupTableContent.Ips
			}

			if backupTableContent.ZooPath != nil {
				backupTableContentMap["zoo_path"] = backupTableContent.ZooPath
			}

			if backupTableContent.Rip != nil {
				backupTableContentMap["rip"] = backupTableContent.Rip
			}

			ids = append(ids, *backupTableContent.Database+FILED_SP+*backupTableContent.Table)
			tmpList = append(tmpList, backupTableContentMap)
		}

		_ = d.Set("available_tables", tmpList)
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
