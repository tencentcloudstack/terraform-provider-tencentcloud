/*
Use this data source to query detailed information of cdb binlog_backup_overview

Example Usage

```hcl
data "tencentcloud_cdb_binlog_backup_overview" "binlog_backup_overview" {
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

func dataSourceTencentCloudCdbBinlogBackupOverview() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCdbBinlogBackupOverviewRead,
		Schema: map[string]*schema.Schema{
			"product": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The type of cloud database product to be queried, currently only supports `mysql`.",
			},

			"binlog_backup_volume": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Total log backup capacity, including off-site log backup (unit is byte).",
			},

			"binlog_backup_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The total number of log backups, including remote log backups.",
			},

			"remote_binlog_volume": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Remote log backup capacity (in bytes).",
			},

			"remote_binlog_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The number of remote log backups.",
			},

			"binlog_archive_volume": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Archived log backup capacity (in bytes).",
			},

			"binlog_archive_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The number of archived log backups.",
			},

			"binlog_standby_volume": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Standard storage log backup capacity (in bytes).",
			},

			"binlog_standby_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The number of standard storage log backups.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudCdbBinlogBackupOverviewRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cdb_binlog_backup_overview.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("product"); ok {
		paramMap["Product"] = helper.String(v.(string))
	}

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCdbBinlogBackupOverviewByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		binlogBackupVolume = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(binlogBackupVolume))
	if binlogBackupVolume != nil {
		_ = d.Set("binlog_backup_volume", binlogBackupVolume)
	}

	if binlogBackupCount != nil {
		_ = d.Set("binlog_backup_count", binlogBackupCount)
	}

	if remoteBinlogVolume != nil {
		_ = d.Set("remote_binlog_volume", remoteBinlogVolume)
	}

	if remoteBinlogCount != nil {
		_ = d.Set("remote_binlog_count", remoteBinlogCount)
	}

	if binlogArchiveVolume != nil {
		_ = d.Set("binlog_archive_volume", binlogArchiveVolume)
	}

	if binlogArchiveCount != nil {
		_ = d.Set("binlog_archive_count", binlogArchiveCount)
	}

	if binlogStandbyVolume != nil {
		_ = d.Set("binlog_standby_volume", binlogStandbyVolume)
	}

	if binlogStandbyCount != nil {
		_ = d.Set("binlog_standby_count", binlogStandbyCount)
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
