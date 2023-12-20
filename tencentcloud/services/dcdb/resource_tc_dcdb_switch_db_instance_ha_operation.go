package dcdb

import (
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDcdbSwitchDbInstanceHaOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcdbSwitchDbInstanceHaOperationCreate,
		Read:   resourceTencentCloudDcdbSwitchDbInstanceHaOperationRead,
		Delete: resourceTencentCloudDcdbSwitchDbInstanceHaOperationDelete,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID in the format of tdsqlshard-ow728lmc.",
			},

			"zone": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Target AZ. The node with the lowest delay in the target AZ will be automatically promoted to primary node.",
			},
		},
	}
}

func resourceTencentCloudDcdbSwitchDbInstanceHaOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dcdb_switch_db_instance_ha_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = dcdb.NewSwitchDBInstanceHARequest()
		service    = DcdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId string
		flowId     *uint64
	)
	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("zone"); ok {
		request.Zone = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDcdbClient().SwitchDBInstanceHA(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		flowId = result.Response.FlowId
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate dcdb switchDbInstanceHaOperation failed, reason:%+v", logId, err)
		return err
	}

	if flowId != nil {
		// need to wait init operation success
		// 0:success; 1:failed, 2:running
		conf := tccommon.BuildStateChangeConf([]string{}, []string{"0"}, 3*tccommon.ReadRetryTimeout, time.Second, service.DcdbDbInstanceStateRefreshFunc(helper.UInt64Int64(*flowId), []string{}))
		if _, e := conf.WaitForState(); e != nil {
			return e
		}
	}

	d.SetId(instanceId)

	return resourceTencentCloudDcdbSwitchDbInstanceHaOperationRead(d, meta)
}

func resourceTencentCloudDcdbSwitchDbInstanceHaOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dcdb_switch_db_instance_ha_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDcdbSwitchDbInstanceHaOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dcdb_switch_db_instance_ha_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
