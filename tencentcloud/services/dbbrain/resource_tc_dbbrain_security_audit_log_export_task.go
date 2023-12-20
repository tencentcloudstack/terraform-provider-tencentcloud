package dbbrain

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDbbrainSecurityAuditLogExportTask() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudDbbrainSecurityAuditLogExportTaskRead,
		Create: resourceTencentCloudDbbrainSecurityAuditLogExportTaskCreate,
		Delete: resourceTencentCloudDbbrainSecurityAuditLogExportTaskDelete,
		// Importer: &schema.ResourceImporter{
		// 	State: schema.ImportStatePassthrough,
		// },
		Schema: map[string]*schema.Schema{
			"sec_audit_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "security audit group id.",
			},

			"start_time": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "start time.",
			},

			"end_time": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "end time.",
			},

			"product": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "product, optional value is mysql.",
			},

			"danger_levels": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Optional:    true,
				ForceNew:    true,
				Description: "List of log risk levels, supported values include: 0 no risk; 1 low risk; 2 medium risk; 3 high risk.",
			},

			"async_request_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "request of async id.",
			},
		},
	}
}

func resourceTencentCloudDbbrainSecurityAuditLogExportTaskCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dbbrain_security_audit_log_export_task.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request         = dbbrain.NewCreateSecurityAuditLogExportTaskRequest()
		response        *dbbrain.CreateSecurityAuditLogExportTaskResponse
		service         = DbbrainService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		ctx             = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		secAuditGroupId string
		asyncRequestId  string
	)

	if v, ok := d.GetOk("sec_audit_group_id"); ok {
		secAuditGroupId = v.(string)
		request.SecAuditGroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("start_time"); ok {
		request.StartTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		request.EndTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product"); ok {
		request.Product = helper.String(v.(string))
	}

	if v, ok := d.GetOk("danger_levels"); ok {
		dangerLevelsSet := v.(*schema.Set).List()
		for i := range dangerLevelsSet {
			dangerLevels := dangerLevelsSet[i].(int)
			request.DangerLevels = append(request.DangerLevels, helper.IntInt64(dangerLevels))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDbbrainClient().CreateSecurityAuditLogExportTask(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dbbrain securityAuditLogExportTask failed, reason:%+v", logId, err)
		return err
	}

	asyncRequestId = helper.UInt64ToStr(*response.Response.AsyncRequestId)
	d.SetId(secAuditGroupId + tccommon.FILED_SP + asyncRequestId)

	err = resource.Retry(2*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ret, err := service.DescribeDbbrainSecurityAuditLogExportTask(ctx, helper.String(secAuditGroupId), helper.String(asyncRequestId), nil)
		if err != nil {
			return tccommon.RetryError(err)
		}
		if ret != nil {
			log.Printf("[DEBUG]%s task.Status:[%s]\n", logId, *ret.Status)
			return nil
		}
		return resource.RetryableError(errors.New("[DEBUG] describe the audit log export task is nil, retry..."))
	})
	if err != nil {
		log.Printf("[CRITAL]%s query dbbrain securityAuditLogExportTask failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudDbbrainSecurityAuditLogExportTaskRead(d, meta)
}

func resourceTencentCloudDbbrainSecurityAuditLogExportTaskRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dbbrain_security_audit_log_export_task.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := DbbrainService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	secAuditGroupId := helper.String(idSplit[0])
	asyncRequestId := helper.String(idSplit[1])

	securityAuditLogExportTask, err := service.DescribeDbbrainSecurityAuditLogExportTask(ctx, secAuditGroupId, asyncRequestId, nil)
	if err != nil {
		return err
	}

	// _ = d.Set("sec_audit_group_id", secAuditGroupId)

	if securityAuditLogExportTask == nil {
		d.SetId("")
		return fmt.Errorf("resource `securityAuditLogExportTask` %s does not exist", d.Id())
	}

	if securityAuditLogExportTask.LogStartTime != nil {
		_ = d.Set("start_time", securityAuditLogExportTask.LogStartTime)
	}

	if securityAuditLogExportTask.LogEndTime != nil {
		_ = d.Set("end_time", securityAuditLogExportTask.LogEndTime)
	}

	if securityAuditLogExportTask.DangerLevels != nil {
		_ = d.Set("danger_levels", securityAuditLogExportTask.DangerLevels)
	}

	if securityAuditLogExportTask.AsyncRequestId != nil {
		_ = d.Set("async_request_id", securityAuditLogExportTask.AsyncRequestId)
	}

	return nil
}

func resourceTencentCloudDbbrainSecurityAuditLogExportTaskDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dbbrain_security_audit_log_export_task.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := DbbrainService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	secAuditGroupId := helper.String(idSplit[0])
	asyncRequestId := helper.String(idSplit[1])

	if err := service.DeleteDbbrainSecurityAuditLogExportTaskById(ctx, secAuditGroupId, asyncRequestId, nil); err != nil {
		return err
	}

	return nil
}
