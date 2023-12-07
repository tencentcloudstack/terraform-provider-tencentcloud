package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudMysqlDataBackupOverview() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMysqlDataBackupOverviewRead,
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

func dataSourceTencentCloudMysqlDataBackupOverviewRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mysql_data_backup_overview.read")()
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
	var dataBackupOverview *cdb.DescribeDataBackupOverviewResponseParams
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMysqlDataBackupOverviewByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		dataBackupOverview = result
		return nil
	})
	if err != nil {
		return err
	}

	if dataBackupOverview.DataBackupVolume != nil {
		_ = d.Set("data_backup_volume", dataBackupOverview.DataBackupVolume)
	}

	if dataBackupOverview.DataBackupCount != nil {
		_ = d.Set("data_backup_count", dataBackupOverview.DataBackupCount)
	}

	if dataBackupOverview.AutoBackupVolume != nil {
		_ = d.Set("auto_backup_volume", dataBackupOverview.AutoBackupVolume)
	}

	if dataBackupOverview.AutoBackupCount != nil {
		_ = d.Set("auto_backup_count", dataBackupOverview.AutoBackupCount)
	}

	if dataBackupOverview.ManualBackupVolume != nil {
		_ = d.Set("manual_backup_volume", dataBackupOverview.ManualBackupVolume)
	}

	if dataBackupOverview.ManualBackupCount != nil {
		_ = d.Set("manual_backup_count", dataBackupOverview.ManualBackupCount)
	}

	if dataBackupOverview.RemoteBackupVolume != nil {
		_ = d.Set("remote_backup_volume", dataBackupOverview.RemoteBackupVolume)
	}

	if dataBackupOverview.RemoteBackupCount != nil {
		_ = d.Set("remote_backup_count", dataBackupOverview.RemoteBackupCount)
	}

	if dataBackupOverview.DataBackupArchiveVolume != nil {
		_ = d.Set("data_backup_archive_volume", dataBackupOverview.DataBackupArchiveVolume)
	}

	if dataBackupOverview.DataBackupArchiveCount != nil {
		_ = d.Set("data_backup_archive_count", dataBackupOverview.DataBackupArchiveCount)
	}

	if dataBackupOverview.DataBackupStandbyVolume != nil {
		_ = d.Set("data_backup_standby_volume", dataBackupOverview.DataBackupStandbyVolume)
	}

	if dataBackupOverview.DataBackupStandbyCount != nil {
		_ = d.Set("data_backup_standby_count", dataBackupOverview.DataBackupStandbyCount)
	}

	d.SetId(helper.DataResourceIdsHash([]string{product}))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), dataBackupOverview); e != nil {
			return e
		}
	}
	return nil
}
