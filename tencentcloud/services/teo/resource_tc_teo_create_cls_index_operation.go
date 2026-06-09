package teo

import (
	"log"

	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudTeoCreateCLSIndexOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoCreateCLSIndexOperationCreate,
		Read:   resourceTencentCloudTeoCreateCLSIndexOperationRead,
		Delete: resourceTencentCloudTeoCreateCLSIndexOperationDelete,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone ID.",
			},
			"task_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Realtime log delivery task ID.",
			},
		},
	}
}

func resourceTencentCloudTeoCreateCLSIndexOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_create_cls_index_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := teov20220901.NewCreateCLSIndexRequest()
	zoneId := d.Get("zone_id").(string)
	request.ZoneId = helper.String(zoneId)
	taskId := d.Get("task_id").(string)
	request.TaskId = helper.String(taskId)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateCLSIndex(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(zoneId)

	return resourceTencentCloudTeoCreateCLSIndexOperationRead(d, meta)
}

func resourceTencentCloudTeoCreateCLSIndexOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_create_cls_index_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudTeoCreateCLSIndexOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_create_cls_index_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
