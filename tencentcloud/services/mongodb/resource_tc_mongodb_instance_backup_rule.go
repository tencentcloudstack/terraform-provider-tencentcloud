package mongodb

import (
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

	logId := tccommon.GetLogId(tccommon.ContextNil)
	var (
		request = mongodb.NewSetBackupRulesRequest()
	)
	instanceId := d.Get("instance_id").(string)
	request.InstanceId = helper.String(instanceId)

	if v, _ := d.GetOk("backup_method"); v != nil {
		request.BackupMethod = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("backup_time"); v != nil {
		request.BackupTime = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("backup_retention_period"); ok {
		request.BackupRetentionPeriod = helper.IntUint64(v.(int))
	}

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

	d.SetId(instanceId)

	return resourceTencentCloudMongodbInstanceBackupRuleRead(d, meta)
}

func resourceTencentCloudMongodbInstanceBackupRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance_backup_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	request := mongodb.NewDescribeBackupRulesRequest()
	request.InstanceId = helper.String(d.Id())
	ratelimit.Check(request.GetAction())
	response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMongodbClient().DescribeBackupRules(request)
	if err != nil {
		return err
	}
	_ = d.Set("instance_id", d.Id())
	_ = d.Set("backup_method", response.Response.BackupMethod)
	_ = d.Set("backup_time", response.Response.BackupTime)
	_ = d.Set("backup_retention_period", response.Response.BackupSaveTime)
	return nil
}

func resourceTencentCloudMongodbInstanceBackupRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance_transparent_data_encryption.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	return resourceTencentCloudMongodbInstanceBackupRuleCreate(d, meta)
}

func resourceTencentCloudMongodbInstanceBackupRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance_backup_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
