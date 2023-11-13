/*
Use this data source to query detailed information of sqlserver datasource_ins_attribute

Example Usage

```hcl
data "tencentcloud_sqlserver_datasource_ins_attribute" "datasource_ins_attribute" {
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

func dataSourceTencentCloudSqlserverDatasourceInsAttribute() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSqlserverDatasourceInsAttributeRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"regular_backup_enable": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Archive backup status. Valid values: enable (enabled), disable (disabled).",
			},

			"regular_backup_save_days": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Archive backup retention period: [90-3650] days.",
			},

			"regular_backup_strategy": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Archive backup policy. Valid values: years (yearly); quarters (quarterly);months` (monthly).",
			},

			"regular_backup_counts": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The number of retained archive backups.",
			},

			"regular_backup_start_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Archive backup start date in YYYY-MM-DD format, which is the current time by default.",
			},

			"blocked_threshold": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Block process threshold in milliseconds.",
			},

			"event_save_days": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Retention period for the files of slow SQL, blocking, deadlock, and extended events.",
			},

			"t_d_e_config": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "TDE transparent data encryption configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"encryption": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether TDE encryption has been enabled, enable-enabled, disable-not enabled.",
						},
						"certificate_attribution": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Certificate attribution. self- means to use the account&amp;#39;s own certificate, others- means to refer to the certificate of other accounts, none- means no certificate.",
						},
						"quote_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Other main account IDs referenced when TDE encryption is activated Note: This field may return null, indicating that no valid value can be obtained.",
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

func dataSourceTencentCloudSqlserverDatasourceInsAttributeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_sqlserver_datasource_ins_attribute.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	service := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSqlserverDatasourceInsAttributeByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		instanceId = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(instanceId))
	if instanceId != nil {
		_ = d.Set("instance_id", instanceId)
	}

	if regularBackupEnable != nil {
		_ = d.Set("regular_backup_enable", regularBackupEnable)
	}

	if regularBackupSaveDays != nil {
		_ = d.Set("regular_backup_save_days", regularBackupSaveDays)
	}

	if regularBackupStrategy != nil {
		_ = d.Set("regular_backup_strategy", regularBackupStrategy)
	}

	if regularBackupCounts != nil {
		_ = d.Set("regular_backup_counts", regularBackupCounts)
	}

	if regularBackupStartTime != nil {
		_ = d.Set("regular_backup_start_time", regularBackupStartTime)
	}

	if blockedThreshold != nil {
		_ = d.Set("blocked_threshold", blockedThreshold)
	}

	if eventSaveDays != nil {
		_ = d.Set("event_save_days", eventSaveDays)
	}

	if tDEConfig != nil {
		tDEConfigAttributeMap := map[string]interface{}{}

		if tDEConfig.Encryption != nil {
			tDEConfigAttributeMap["encryption"] = tDEConfig.Encryption
		}

		if tDEConfig.CertificateAttribution != nil {
			tDEConfigAttributeMap["certificate_attribution"] = tDEConfig.CertificateAttribution
		}

		if tDEConfig.QuoteUin != nil {
			tDEConfigAttributeMap["quote_uin"] = tDEConfig.QuoteUin
		}

		ids = append(ids, *tDEConfig.InstanceId)
		_ = d.Set("t_d_e_config", tDEConfigAttributeMap)
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
