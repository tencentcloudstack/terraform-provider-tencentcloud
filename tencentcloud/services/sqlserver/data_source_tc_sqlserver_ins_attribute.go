package sqlserver

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudSqlserverInsAttribute() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSqlserverInsAttributeRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
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
			"tde_config": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "TDE Transparent Data Encryption Configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"certificate_attribution": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Certificate ownership. Self - indicates using the account's own certificate, others - indicates referencing certificates from other accounts, and none - indicates no certificate.",
						},
						"encryption": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "TDE encryption, 'enable' - enabled, 'disable' - not enabled.",
						},
						"quote_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Other primary account IDs referenced when activating TDE encryption\nNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
					},
				},
			},
			"ssl_config": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "SSL encryption.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"encryption": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SSL encryption status, enable - turned on, disable-not turned on, enable_doing - enabling, disable_doing-closing, renew_doing-updating, wait_doing-wait for execution within maintenance time Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"ssl_validity_period": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SSL certificate validity period, time format YYYY-MM-DD HH:MM:SS Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"ssl_validity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "SSL certificate validity, 0-invalid, 1-valid Note: This field may return null, indicating that no valid value can be obtained.",
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

func dataSourceTencentCloudSqlserverInsAttributeRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_sqlserver_ins_attribute.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service      = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		insAttribute *sqlserver.DescribeDBInstancesAttributeResponseParams
		instanceId   string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSqlserverInsAttributeByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		insAttribute = result
		return nil
	})

	if err != nil {
		return err
	}

	if insAttribute.InstanceId != nil {
		_ = d.Set("instance_id", instanceId)
	}

	if insAttribute.RegularBackupEnable != nil {
		_ = d.Set("regular_backup_enable", insAttribute.RegularBackupEnable)
	}

	if insAttribute.RegularBackupSaveDays != nil {
		_ = d.Set("regular_backup_save_days", insAttribute.RegularBackupSaveDays)
	}

	if insAttribute.RegularBackupStrategy != nil {
		_ = d.Set("regular_backup_strategy", insAttribute.RegularBackupStrategy)
	}

	if insAttribute.RegularBackupCounts != nil {
		_ = d.Set("regular_backup_counts", insAttribute.RegularBackupCounts)
	}

	if insAttribute.RegularBackupStartTime != nil {
		_ = d.Set("regular_backup_start_time", insAttribute.RegularBackupStartTime)
	}

	if insAttribute.BlockedThreshold != nil {
		_ = d.Set("blocked_threshold", insAttribute.BlockedThreshold)
	}

	if insAttribute.EventSaveDays != nil {
		_ = d.Set("event_save_days", insAttribute.EventSaveDays)
	}

	if insAttribute.TDEConfig != nil {
		tmpList := make([]map[string]interface{}, 0)
		configMap := map[string]interface{}{}
		if insAttribute.TDEConfig.CertificateAttribution != nil {
			configMap["certificate_attribution"] = insAttribute.TDEConfig.CertificateAttribution
		}

		if insAttribute.TDEConfig.Encryption != nil {
			configMap["encryption"] = insAttribute.TDEConfig.Encryption
		}

		if insAttribute.TDEConfig.QuoteUin != nil {
			configMap["quote_uin"] = insAttribute.TDEConfig.QuoteUin
		}

		tmpList = append(tmpList, configMap)

		_ = d.Set("tde_config", tmpList)
	}

	if insAttribute.SSLConfig != nil {
		tmpList := make([]map[string]interface{}, 0)
		configMap := map[string]interface{}{}
		if insAttribute.SSLConfig.Encryption != nil {
			configMap["encryption"] = insAttribute.SSLConfig.Encryption
		}

		if insAttribute.SSLConfig.SSLValidityPeriod != nil {
			configMap["ssl_validity_period"] = insAttribute.SSLConfig.SSLValidityPeriod
		}

		if insAttribute.SSLConfig.SSLValidity != nil {
			configMap["ssl_validity"] = insAttribute.SSLConfig.SSLValidity
		}

		tmpList = append(tmpList, configMap)

		_ = d.Set("ssl_config", tmpList)
	}

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
