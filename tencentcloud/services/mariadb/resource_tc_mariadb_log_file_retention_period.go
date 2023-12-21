package mariadb

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMariadbLogFileRetentionPeriod() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudMariadbLogFileRetentionPeriodRead,
		Create: resourceTencentCloudMariadbLogFileRetentionPeriodCreate,
		Update: resourceTencentCloudMariadbLogFileRetentionPeriodUpdate,
		Delete: resourceTencentCloudMariadbLogFileRetentionPeriodDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "instance id.",
			},

			"days": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The number of days to save, cannot exceed 30.",
			},
		},
	}
}

func resourceTencentCloudMariadbLogFileRetentionPeriodCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mariadb_log_file_retention_period.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)
	return resourceTencentCloudMariadbLogFileRetentionPeriodUpdate(d, meta)
}

func resourceTencentCloudMariadbLogFileRetentionPeriodRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mariadb_log_file_retention_period.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MariadbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	instanceId := d.Id()

	logFileRetentionPeriod, err := service.DescribeMariadbLogFileRetentionPeriod(ctx, instanceId)

	if err != nil {
		return err
	}

	if logFileRetentionPeriod == nil {
		d.SetId("")
		return fmt.Errorf("resource `logFileRetentionPeriod` %s does not exist", instanceId)
	}

	if logFileRetentionPeriod.InstanceId != nil {
		_ = d.Set("instance_id", logFileRetentionPeriod.InstanceId)
	}

	if logFileRetentionPeriod.Days != nil {
		_ = d.Set("days", int(*logFileRetentionPeriod.Days))
	}

	return nil
}

func resourceTencentCloudMariadbLogFileRetentionPeriodUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mariadb_log_file_retention_period.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := mariadb.NewModifyLogFileRetentionPeriodRequest()

	instanceId := d.Id()

	request.InstanceId = &instanceId

	if v, ok := d.GetOk("days"); ok {
		request.Days = helper.Uint64(uint64(v.(int)))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMariadbClient().ModifyLogFileRetentionPeriod(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create mariadb logFileRetentionPeriod failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMariadbLogFileRetentionPeriodRead(d, meta)
}

func resourceTencentCloudMariadbLogFileRetentionPeriodDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mariadb_log_file_retention_period.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
