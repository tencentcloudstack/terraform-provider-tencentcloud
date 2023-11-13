/*
Use this data source to query detailed information of cdb backup_overview

Example Usage

```hcl
data "tencentcloud_cdb_backup_overview" "backup_overview" {
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

func dataSourceTencentCloudCdbBackupOverview() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCdbBackupOverviewRead,
		Schema: map[string]*schema.Schema{
			"product": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The type of cloud database product to be queried, currently only supports `mysql`.",
			},

			"backup_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The total number of user backups in the current region (including data backups and log backups).",
			},

			"backup_volume": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The total backup capacity of the user in the current region.",
			},

			"billing_volume": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The billable capacity of the user&amp;amp;#39;s backup in the current region, that is, the part that exceeds the gifted capacity.",
			},

			"free_volume": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The free backup capacity obtained by the user in the current region.",
			},

			"remote_backup_volume": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The total capacity of off-site backup of the user in the current region. Note: This field may return null, indicating that no valid value can be obtained.",
			},

			"backup_archive_volume": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Archive backup capacity, including data backup and log backup. Note: This field may return null, indicating that no valid value can be obtained.",
			},

			"backup_standby_volume": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Standard storage backup capacity, including data backup and log backup. Note: This field may return null, indicating that no valid value can be obtained.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudCdbBackupOverviewRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cdb_backup_overview.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("product"); ok {
		paramMap["Product"] = helper.String(v.(string))
	}

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCdbBackupOverviewByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		backupCount = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(backupCount))
	if backupCount != nil {
		_ = d.Set("backup_count", backupCount)
	}

	if backupVolume != nil {
		_ = d.Set("backup_volume", backupVolume)
	}

	if billingVolume != nil {
		_ = d.Set("billing_volume", billingVolume)
	}

	if freeVolume != nil {
		_ = d.Set("free_volume", freeVolume)
	}

	if remoteBackupVolume != nil {
		_ = d.Set("remote_backup_volume", remoteBackupVolume)
	}

	if backupArchiveVolume != nil {
		_ = d.Set("backup_archive_volume", backupArchiveVolume)
	}

	if backupStandbyVolume != nil {
		_ = d.Set("backup_standby_volume", backupStandbyVolume)
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
