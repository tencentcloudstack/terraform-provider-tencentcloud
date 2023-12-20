package dcdb

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDcdbActivateHourInstanceOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcdbActivateHourInstanceOperationCreate,
		Read:   resourceTencentCloudDcdbActivateHourInstanceOperationRead,
		Delete: resourceTencentCloudDcdbActivateHourInstanceOperationDelete,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "instance ID in the format of dcdbt-ow728lmc, which can be obtained through the `DescribeDCDBInstances` API.",
			},
		},
	}
}

func resourceTencentCloudDcdbActivateHourInstanceOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dcdb_activate_hour_instance_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = dcdb.NewActiveHourDCDBInstanceRequest()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceIds = []*string{helper.String(instanceId)}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDcdbClient().ActiveHourDCDBInstance(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate dcdb activateHourInstanceOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	return resourceTencentCloudDcdbActivateHourInstanceOperationRead(d, meta)
}

func resourceTencentCloudDcdbActivateHourInstanceOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dcdb_activate_hour_instance_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDcdbActivateHourInstanceOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dcdb_activate_hour_instance_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
