package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMpsManageTaskOperation() *schema.Resource {
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
	defer logElapsed("resource.tencentcloud_mps_manage_task_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMpsClient().ManageTask(request)
		if e != nil {
			return retryError(e)
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
	defer logElapsed("resource.tencentcloud_mps_manage_task_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMpsManageTaskOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_manage_task_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
