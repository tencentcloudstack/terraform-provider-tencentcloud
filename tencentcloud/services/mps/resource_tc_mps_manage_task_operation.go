package mps

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMpsManageTaskOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsManageTaskOperationCreate,
		Read:   resourceTencentCloudMpsManageTaskOperationRead,
		Delete: resourceTencentCloudMpsManageTaskOperationDelete,
		Schema: map[string]*schema.Schema{
			"operation_type": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Operation type. Valid values:`Abort`: task termination. Notice: If the task type is live stream processing (LiveStreamProcessTask), tasks whose task status is `WAITING` or `PROCESSING` can be terminated.For other task types, only tasks whose task status is `WAITING` can be terminated.",
			},

			"task_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Video processing task ID.",
			},
		},
	}
}

func resourceTencentCloudMpsManageTaskOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_manage_task_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request = mps.NewManageTaskRequest()
		taskId  string
	)
	if v, ok := d.GetOk("operation_type"); ok {
		request.OperationType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("task_id"); ok {
		request.TaskId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMpsClient().ManageTask(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate mps manageTaskOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(taskId)

	return resourceTencentCloudMpsManageTaskOperationRead(d, meta)
}

func resourceTencentCloudMpsManageTaskOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_manage_task_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMpsManageTaskOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_manage_task_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
