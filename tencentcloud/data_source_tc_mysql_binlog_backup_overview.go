/*
Use this data source to query detailed information of mysql binlog_backup_overview

Example Usage

```hcl
data "tencentcloud_mysql_binlog_backup_overview" "binlog_backup_overview" {
  product = "mysql"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudMysqlBinlogBackupOverview() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMysqlBinlogBackupOverviewRead,
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

func dataSourceTencentCloudMysqlBinlogBackupOverviewRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mysql_binlog_backup_overview.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	product := ""
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("product"); ok {
		product = v.(string)
		paramMap["Product"] = helper.String(v.(string))
	}

	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	var binlogBackupOverview *cdb.DescribeBinlogBackupOverviewResponseParams
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMysqlBinlogBackupOverviewByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		binlogBackupOverview = result
		return nil
	})
	if err != nil {
		return err
	}

	if binlogBackupOverview.BinlogBackupVolume != nil {
		_ = d.Set("binlog_backup_volume", binlogBackupOverview.BinlogBackupVolume)
	}

	if binlogBackupOverview.BinlogBackupCount != nil {
		_ = d.Set("binlog_backup_count", binlogBackupOverview.BinlogBackupCount)
	}

	if binlogBackupOverview.RemoteBinlogVolume != nil {
		_ = d.Set("remote_binlog_volume", binlogBackupOverview.RemoteBinlogVolume)
	}

	if binlogBackupOverview.RemoteBinlogCount != nil {
		_ = d.Set("remote_binlog_count", binlogBackupOverview.RemoteBinlogCount)
	}

	if binlogBackupOverview.BinlogArchiveVolume != nil {
		_ = d.Set("binlog_archive_volume", binlogBackupOverview.BinlogArchiveVolume)
	}

	if binlogBackupOverview.BinlogArchiveCount != nil {
		_ = d.Set("binlog_archive_count", binlogBackupOverview.BinlogArchiveCount)
	}

	if binlogBackupOverview.BinlogStandbyVolume != nil {
		_ = d.Set("binlog_standby_volume", binlogBackupOverview.BinlogStandbyVolume)
	}

	if binlogBackupOverview.BinlogStandbyCount != nil {
		_ = d.Set("binlog_standby_count", binlogBackupOverview.BinlogStandbyCount)
	}

	d.SetId(helper.DataResourceIdsHash([]string{product}))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), binlogBackupOverview); e != nil {
			return e
		}
	}
	return nil
}
