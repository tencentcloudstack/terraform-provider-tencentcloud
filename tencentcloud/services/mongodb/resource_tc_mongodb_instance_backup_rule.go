package mongodb

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mongodb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func ResourceTencentCloudMongodbInstanceBackupRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMongodbInstanceBackupRuleCreate,
		Read:   resourceTencentCloudMongodbInstanceBackupRuleRead,
		Update: resourceTencentCloudMongodbInstanceBackupRuleUpdate,
		Delete: resourceTencentCloudMongodbInstanceBackupRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"backup_method": {
				Required: true,
				Type:     schema.TypeInt,
				Description: "Set automatic backup method. Valid values:\n" +
					"- 0: Logical backup;\n" +
					"- 1: Physical backup;\n" +
					"- 3: Snapshot backup (supported only in cloud disk version).",
			},

			"backup_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Set the start time for automatic backup. The value range is: [0,23]. For example, setting this parameter to 2 means that backup starts at 02:00.",
			},

			"backup_frequency": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specify the daily automatic backup frequency. 12: Back up twice a day, approximately 12 hours apart; 24: Back up once a day (default), approximately 24 hours apart.",
			},

			"notify": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Set whether to send failure alerts when automatic backup errors occur.\n- true: Send.\n- false: Do not send.",
			},

			"backup_retention_period": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Specifies the retention period for backup data. Unit: days, default is 7 days. Value range: [7, 365].",
			},

			"active_weekdays": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specify the specific dates for automatic backups to be performed each week. Format: Enter a number between 0 and 6 to represent Sunday through Saturday (e.g., 1 represents Monday). Separate multiple dates with commas (,). Example: Entering 1,3,5 means the system will perform backups on Mondays, Wednesdays, and Fridays every week. Default: If not set, the default is a full cycle (0,1,2,3,4,5,6), meaning backups will be performed daily.",
			},

			"long_term_unit": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Long-term retention period. Supports selecting specific dates for backups on a weekly or monthly basis (e.g., backup data for the 1st and 15th of each month) to retain for a longer period. Disabled (default): Long-term retention is disabled. Weekly retention: Specify `weekly`. Monthly retention: Specify `monthly`.",
			},

			"long_term_active_days": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specify the specific backup dates to be retained long-term. This setting only takes effect when LongTermUnit is set to weekly or monthly. Weekly Retention: Enter a number between 0 and 6 to represent Sunday through Saturday. Separate multiple dates with commas. Monthly Retention: Enter a number between 1 and 31 to represent specific dates within the month. Separate multiple dates with commas.",
			},

			"long_term_expired_days": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Long-term backup retention period. Value range [30, 1075].",
			},

			"oplog_expired_days": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Incremental backup retention period. Unit: days. Default value: 7 days. Value range: [7,365].",
			},

			"backup_version": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Backup version. Old version backup is 0, advanced backup is 1. Set this value to 1 when enabling advanced backup.",
			},

			"alarm_water_level": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Sets the alarm threshold for backup dataset storage space usage. Unit: %. Default value: 100. Value range: [50, 300].",
			},
		},
	}
}

func resourceTencentCloudMongodbInstanceBackupRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance_backup_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		instanceId string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)
	return resourceTencentCloudMongodbInstanceBackupRuleUpdate(d, meta)
}

func resourceTencentCloudMongodbInstanceBackupRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance_backup_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = mongodb.NewDescribeBackupRulesRequest()
		response   = mongodb.NewDescribeBackupRulesResponse()
		instanceId = d.Id()
	)

	request.InstanceId = &instanceId
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMongodbClient().DescribeBackupRules(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe mongodb backup rules failed, Response is nil"))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s describe mongodb backup rules failed, reason:%+v", logId, err)
		return err
	}

	_ = d.Set("instance_id", instanceId)

	if response.Response.BackupMethod != nil {
		_ = d.Set("backup_method", response.Response.BackupMethod)
	}

	if response.Response.BackupTime != nil {
		_ = d.Set("backup_time", response.Response.BackupTime)
	}

	if response.Response.BackupSaveTime != nil {
		_ = d.Set("backup_retention_period", response.Response.BackupSaveTime)
	}

	return nil
}

func resourceTencentCloudMongodbInstanceBackupRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance_backup_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = mongodb.NewSetBackupRulesRequest()
		instanceId = d.Id()
	)

	if v, ok := d.GetOkExists("backup_method"); ok {
		request.BackupMethod = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("backup_time"); ok {
		request.BackupTime = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("backup_frequency"); ok {
		request.BackupFrequency = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("notify"); ok {
		request.Notify = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("backup_retention_period"); ok {
		request.BackupRetentionPeriod = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("active_weekdays"); ok {
		request.ActiveWeekdays = helper.String(v.(string))
	}

	if v, ok := d.GetOk("long_term_unit"); ok {
		request.LongTermUnit = helper.String(v.(string))
	}

	if v, ok := d.GetOk("long_term_active_days"); ok {
		request.LongTermActiveDays = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("long_term_expired_days"); ok {
		request.LongTermExpiredDays = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("oplog_expired_days"); ok {
		request.OplogExpiredDays = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("backup_version"); ok {
		request.BackupVersion = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("alarm_water_level"); ok {
		request.AlarmWaterLevel = helper.IntInt64(v.(int))
	}

	request.InstanceId = &instanceId
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMongodbClient().SetBackupRules(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate mongodb backupRule failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMongodbInstanceBackupRuleRead(d, meta)
}

func resourceTencentCloudMongodbInstanceBackupRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance_backup_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
