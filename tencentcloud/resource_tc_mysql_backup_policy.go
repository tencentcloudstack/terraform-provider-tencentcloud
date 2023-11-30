package tencentcloud

import (
	"bytes"
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTencentCloudMysqlBackupPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlBackupPolicyCreate,
		Read:   resourceTencentCloudMysqlBackupPolicyRead,
		Update: resourceTencentCloudMysqlBackupPolicyUpdate,
		Delete: resourceTencentCloudMysqlBackupPolicyDelete,

		Schema: map[string]*schema.Schema{
			"mysql_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Instance ID to which policies will be applied.",
			},
			"retention_period": {
				Type:         schema.TypeInt,
				ValidateFunc: validateIntegerInRange(7, 1830),
				Optional:     true,
				Default:      7,
				Description:  "The retention time of backup files, in days. The minimum value is 7 days and the maximum value is 1830 days. And default value is `7`.",
			},
			"backup_model": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      MYSQL_ALLOW_BACKUP_MODEL[1],
				ValidateFunc: validateAllowedStringValue(MYSQL_ALLOW_BACKUP_MODEL),
				Description:  "Backup method. Supported values include: `physical` - physical backup.",
			},
			"backup_time": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      MYSQL_ALLOW_BACKUP_TIME[0],
				ValidateFunc: validateAllowedStringValue(MYSQL_ALLOW_BACKUP_TIME),
				Description:  "Instance backup time, in the format of 'HH:mm-HH:mm'. Time setting interval is four hours. Default to `02:00-06:00`. The following value can be supported: `02:00-06:00`, `06:00-10:00`, `10:00-14:00`, `14:00-18:00`, `18:00-22:00`, and `22:00-02:00`.",
			},

			"binlog_period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Binlog retention time, in days. The minimum value is 7 days and the maximum value is 1830 days. This value cannot be set greater than the backup file retention time.",
			},

			"enable_binlog_standby": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "off",
				Description: "Whether to enable the log backup standard storage policy, `off` - close, `on` - open, the default is off.",
			},

			"binlog_standby_days": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The standard starting number of days for log backup storage. The log backup will be converted when it reaches the standard starting number of days for storage. The minimum is 30 days and must not be greater than the number of days for log backup retention.",
			},
		},
	}
}

func resourceTencentCloudMysqlBackupPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_backup_policy.create")()

	d.SetId(d.Get("mysql_id").(string))

	return resourceTencentCloudMysqlBackupPolicyUpdate(d, meta)
}

func resourceTencentCloudMysqlBackupPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_backup_policy.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	mysqlService := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		desResponse, e := mysqlService.DescribeBackupConfigByMysqlId(ctx, d.Id())
		if e != nil {
			if mysqlService.NotFoundMysqlInstance(e) {
				d.SetId("")
				return nil
			}
			return retryError(e)
		}
		_ = d.Set("mysql_id", d.Id())
		_ = d.Set("retention_period", int(*desResponse.Response.BackupExpireDays))
		_ = d.Set("backup_model", *desResponse.Response.BackupMethod)
		var buf bytes.Buffer

		if *desResponse.Response.StartTimeMin < 10 {
			buf.WriteString("0")
		}
		buf.WriteString(fmt.Sprintf("%d:00-", *desResponse.Response.StartTimeMin))

		if *desResponse.Response.StartTimeMax < 10 {
			buf.WriteString("0")
		}
		buf.WriteString(fmt.Sprintf("%d:00", *desResponse.Response.StartTimeMax))
		_ = d.Set("backup_time", buf.String())
		_ = d.Set("binlog_period", int(*desResponse.Response.BinlogExpireDays))

		if desResponse.Response.EnableBinlogStandby != nil {
			_ = d.Set("enable_binlog_standby", desResponse.Response.EnableBinlogStandby)
		}

		if desResponse.Response.BinlogStandbyDays != nil {
			_ = d.Set("binlog_standby_days", desResponse.Response.BinlogStandbyDays)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("[API]Describe mysql backup policy fail,reason:%s", err.Error())
	}
	return nil
}

func resourceTencentCloudMysqlBackupPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_backup_policy.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	mysqlService := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		isUpdate = false

		mysqlId         = d.Get("mysql_id").(string)
		retentionPeriod = int64(d.Get("retention_period").(int))
		backupModel     = d.Get("backup_model").(string)
		backupTime      = d.Get("backup_time").(string)

		binlogExpireDays    int64
		enableBinlogStandby string
		binlogStandbyDays   int64
	)

	if v, ok := d.GetOkExists("binlog_period"); ok {
		binlogExpireDays = int64(v.(int))
	}

	if v, ok := d.GetOk("enable_binlog_standby"); ok {
		enableBinlogStandby = v.(string)
	}

	if v, ok := d.GetOkExists("binlog_standby_days"); ok {
		binlogStandbyDays = int64(v.(int))
	}

	if d.HasChange("retention_period") || d.HasChange("backup_model") || d.HasChange("backup_time") ||
		d.HasChange("binlog_period") || d.HasChange("enable_binlog_standby") || d.HasChange("binlog_standby_days") {
		if backupModel != "physical" {
			return fmt.Errorf("`backup_model` only support 'physical'")
		}
		isUpdate = true
	}

	if isUpdate {
		err := mysqlService.ModifyBackupConfigByMysqlId(ctx, mysqlId, retentionPeriod, backupModel, backupTime, binlogExpireDays, enableBinlogStandby, binlogStandbyDays)
		if err != nil {
			return err
		}
	}
	return resourceTencentCloudMysqlBackupPolicyRead(d, meta)
}

//set all config to default
func resourceTencentCloudMysqlBackupPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_backup_policy.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	mysqlService := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		retentionPeriod     int64 = 7
		backupModel               = MYSQL_ALLOW_BACKUP_MODEL[1]
		backupTime                = MYSQL_ALLOW_BACKUP_TIME[0]
		binlogExpireDays    int64 = 7
		enableBinlogStandby       = "off"
		binlogStandbyDays   int64 = 180
	)
	err := mysqlService.ModifyBackupConfigByMysqlId(ctx, d.Id(), retentionPeriod, backupModel, backupTime, binlogExpireDays, enableBinlogStandby, binlogStandbyDays)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
