/*
Use this data source to query detailed information of cdb data_backup_overview

Example Usage

```hcl
data "tencentcloud_cdb_data_backup_overview" "data_backup_overview" {
  product = "mysql"
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

func dataSourceTencentCloudCdbDataBackupOverview() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCdbDataBackupOverviewRead,
		Schema: map[string]*schema.Schema{
			"product": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The type of cloud database product to be queried, currently only supports `mysql`.",
			},

			"data_backup_volume": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Total data backup capacity of the current region (including automatic backup and manual backup, in bytes).",
			},

			"data_backup_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The total number of data backups in the current region.",
			},

			"auto_backup_volume": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The total automatic backup capacity of the current region.",
			},

			"auto_backup_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The total number of automatic backups in the current region.",
			},

			"manual_backup_volume": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The total manual backup capacity of the current region.",
			},

			"manual_backup_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The total number of manual backups in the current region.",
			},

			"remote_backup_volume": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The total capacity of remote backup.",
			},

			"remote_backup_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The total number of remote backups.",
			},

			"data_backup_archive_volume": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The total capacity of the current regional archive backup.",
			},

			"data_backup_archive_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The total number of archive backups in the current region.",
			},

			"data_backup_standby_volume": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The total backup capacity of the current regional standard storage.",
			},

			"data_backup_standby_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The total number of standard storage backups in the current region.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudCdbDataBackupOverviewRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cdb_data_backup_overview.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("product"); ok {
		paramMap["Product"] = helper.String(v.(string))
	}

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCdbDataBackupOverviewByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		dataBackupVolume = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(dataBackupVolume))
	if dataBackupVolume != nil {
		_ = d.Set("data_backup_volume", dataBackupVolume)
	}

	if dataBackupCount != nil {
		_ = d.Set("data_backup_count", dataBackupCount)
	}

	if autoBackupVolume != nil {
		_ = d.Set("auto_backup_volume", autoBackupVolume)
	}

	if autoBackupCount != nil {
		_ = d.Set("auto_backup_count", autoBackupCount)
	}

	if manualBackupVolume != nil {
		_ = d.Set("manual_backup_volume", manualBackupVolume)
	}

	if manualBackupCount != nil {
		_ = d.Set("manual_backup_count", manualBackupCount)
	}

	if remoteBackupVolume != nil {
		_ = d.Set("remote_backup_volume", remoteBackupVolume)
	}

	if remoteBackupCount != nil {
		_ = d.Set("remote_backup_count", remoteBackupCount)
	}

	if dataBackupArchiveVolume != nil {
		_ = d.Set("data_backup_archive_volume", dataBackupArchiveVolume)
	}

	if dataBackupArchiveCount != nil {
		_ = d.Set("data_backup_archive_count", dataBackupArchiveCount)
	}

	if dataBackupStandbyVolume != nil {
		_ = d.Set("data_backup_standby_volume", dataBackupStandbyVolume)
	}

	if dataBackupStandbyCount != nil {
		_ = d.Set("data_backup_standby_count", dataBackupStandbyCount)
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
