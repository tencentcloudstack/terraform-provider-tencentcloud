package postgresql

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudPostgresqlModifySwitchTimePeriodOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlModifySwitchTimePeriodOperationCreate,
		Read:   resourceTencentCloudPostgresqlModifySwitchTimePeriodOperationRead,
		Delete: resourceTencentCloudPostgresqlModifySwitchTimePeriodOperationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The ID of the instance waiting for a switch.",
			},

			"switch_tag": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Valid value: `0` (switch immediately).",
			},
		},
	}
}

func resourceTencentCloudPostgresqlModifySwitchTimePeriodOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_modify_switch_time_period_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request      = postgresql.NewModifySwitchTimePeriodRequest()
		dBInstanceId string
	)
	if v, ok := d.GetOk("db_instance_id"); ok {
		request.DBInstanceId = helper.String(v.(string))
		dBInstanceId = v.(string)
	}

	if v, _ := d.GetOk("switch_tag"); v != nil {
		request.SwitchTag = helper.IntUint64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().ModifySwitchTimePeriod(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate postgresql ModifySwitchTimePeriodOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(dBInstanceId)

	return resourceTencentCloudPostgresqlModifySwitchTimePeriodOperationRead(d, meta)
}

func resourceTencentCloudPostgresqlModifySwitchTimePeriodOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_modify_switch_time_period_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudPostgresqlModifySwitchTimePeriodOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_modify_switch_time_period_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
