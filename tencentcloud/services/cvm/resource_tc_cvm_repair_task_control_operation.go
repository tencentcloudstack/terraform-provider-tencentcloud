package cvm

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCvmRepairTaskControlOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCvmRepairTaskControlOperationCreate,
		Read:   resourceTencentCloudCvmRepairTaskControlOperationRead,
		Delete: resourceTencentCloudCvmRepairTaskControlOperationDelete,
		Schema: map[string]*schema.Schema{
			"product": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Product type the pending-authorization task instance belongs to. Valid values: `CVM` (Cloud Virtual Machine), `CDH` (Cloud Dedicated Host), `CPM2.0` (Cloud Physical Machine 2.0).",
			},
			"instance_ids": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of instance IDs to operate on. Only repair tasks related to these instance IDs are authorized. Can be obtained from `InstanceId` in the `DescribeTaskInfo` API response.",
			},
			"task_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the repair task to operate on. Can be obtained from `TaskId` in the `DescribeTaskInfo` API response.",
			},
			"operate": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Operation type. Currently only `AuthorizeRepair` is supported.",
			},
			"order_auth_time": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Scheduled authorization time, format `YYYY-MM-DD HH:MM:SS`. The scheduled time must be at least 5 minutes later than the current time and within 48 hours.",
			},
			"task_sub_method": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Additional authorization handling strategy. When empty, the default authorization is used. For repair tasks supporting lossy migration, set to `LossyLocal` to allow lossy local-disk migration. WARNING: when `LossyLocal` is used on a local-disk instance, all local disk data will be wiped, equivalent to redeploying the local-disk instance.",
			},
		},
	}
}

func resourceTencentCloudCvmRepairTaskControlOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_repair_task_control_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request = cvm.NewRepairTaskControlRequest()
		taskId  string
	)

	if v, ok := d.GetOk("product"); ok {
		request.Product = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsList := v.([]interface{})
		for _, instanceId := range instanceIdsList {
			request.InstanceIds = append(request.InstanceIds, helper.String(instanceId.(string)))
		}
	}

	if v, ok := d.GetOk("task_id"); ok {
		request.TaskId = helper.String(v.(string))
		taskId = v.(string)
	}

	if v, ok := d.GetOk("operate"); ok {
		request.Operate = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_auth_time"); ok {
		request.OrderAuthTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("task_sub_method"); ok {
		request.TaskSubMethod = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().RepairTaskControl(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			if result != nil && result.Response != nil && result.Response.TaskId != nil {
				taskId = *result.Response.TaskId
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cvm repairTaskControlOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(taskId)

	return resourceTencentCloudCvmRepairTaskControlOperationRead(d, meta)
}

func resourceTencentCloudCvmRepairTaskControlOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_repair_task_control_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCvmRepairTaskControlOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_repair_task_control_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
