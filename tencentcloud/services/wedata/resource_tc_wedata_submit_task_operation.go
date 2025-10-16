package wedata

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWedataSubmitTaskOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataSubmitTaskOperationCreate,
		Read:   resourceTencentCloudWedataSubmitTaskOperationRead,
		Delete: resourceTencentCloudWedataSubmitTaskOperationDelete,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Project ID.",
			},

			"task_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Task ID.",
			},

			"version_remark": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Version remarks.",
			},
			"status": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Status.",
			},
			"version_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Version id.",
			},
		},
	}
}

func resourceTencentCloudWedataSubmitTaskOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_submit_task_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		projectId string
		taskId    string
	)
	var (
		request  = wedatav20250806.NewSubmitTaskRequest()
		response = wedatav20250806.NewSubmitTaskResponse()
	)

	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("task_id"); ok {
		request.TaskId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("version_remark"); ok {
		request.VersionRemark = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().SubmitTaskWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create wedata submit task operation failed, reason:%+v", logId, err)
		return err
	}

	if response != nil && response.Response != nil && response.Response.Data != nil {
		if response.Response.Data.Status != nil {
			_ = d.Set("status", response.Response.Data.Status)
		}
		_ = d.Set("version_id", response.Response.Data.VersionId)
	}

	d.SetId(strings.Join([]string{projectId, taskId}, tccommon.FILED_SP))

	return resourceTencentCloudWedataSubmitTaskOperationRead(d, meta)
}

func resourceTencentCloudWedataSubmitTaskOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_submit_task_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudWedataSubmitTaskOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_submit_task_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
