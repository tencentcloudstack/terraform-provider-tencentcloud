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
				Description: "Instance id.",
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

			"backup_retention_period": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Specify the number of days to save backup data. The default is 7 days, and the support settings are 7, 30, 90, 180, 365.",
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

	if v, ok := d.GetOkExists("backup_retention_period"); ok {
		request.BackupRetentionPeriod = helper.IntUint64(v.(int))
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
