/*
Provides a resource to create a dbbrain security_audit_log_export_task

Example Usage

```hcl
resource "tencentcloud_dbbrain_security_audit_log_export_task" "security_audit_log_export_task" {
  sec_audit_group_id = ""
  start_time = ""
  end_time = ""
  product = ""
  danger_levels = ""
}

```
Import

dbbrain security_audit_log_export_task can be imported using the id, e.g.
```
$ terraform import tencentcloud_dbbrain_security_audit_log_export_task.security_audit_log_export_task securityAuditLogExportTask_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDbbrainSecurityAuditLogExportTask() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudDbbrainSecurityAuditLogExportTaskRead,
		Create: resourceTencentCloudDbbrainSecurityAuditLogExportTaskCreate,
		Delete: resourceTencentCloudDbbrainSecurityAuditLogExportTaskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"sec_audit_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "security audit group id.",
			},

			"start_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "start time.",
			},

			"end_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "end time.",
			},

			"product": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "product, optional value is mysql.",
			},

			"danger_levels": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Optional:    true,
				Description: "List of log risk levels, supported values include: 0 no risk; 1 low risk; 2 medium risk; 3 high risk.",
			},
		},
	}
}

func resourceTencentCloudDbbrainSecurityAuditLogExportTaskCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dbbrain_security_audit_log_export_task.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request         = dbbrain.NewCreateSecurityAuditLogExportTaskRequest()
		response        *dbbrain.CreateSecurityAuditLogExportTaskResponse
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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDbbrainClient().CreateSecurityAuditLogExportTask(request)
		if e != nil {
			return retryError(e)
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

	d.SetId(secAuditGroupId + FILED_SP + asyncRequestId)
	return resourceTencentCloudDbbrainSecurityAuditLogExportTaskRead(d, meta)
}

func resourceTencentCloudDbbrainSecurityAuditLogExportTaskRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dbbrain_security_audit_log_export_task.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DbbrainService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	secAuditGroupId := helper.String(idSplit[0])
	asyncRequestId := helper.String(idSplit[1])

	securityAuditLogExportTask, err := service.DescribeDbbrainSecurityAuditLogExportTask(ctx, secAuditGroupId, asyncRequestId, nil)

	if err != nil {
		return err
	}

	if securityAuditLogExportTask == nil {
		d.SetId("")
		return fmt.Errorf("resource `securityAuditLogExportTask` %s does not exist", d.Id())
	}

	if securityAuditLogExportTask.StartTime != nil {
		_ = d.Set("start_time", securityAuditLogExportTask.StartTime)
	}

	if securityAuditLogExportTask.EndTime != nil {
		_ = d.Set("end_time", securityAuditLogExportTask.EndTime)
	}

	if securityAuditLogExportTask.DangerLevels != nil {
		_ = d.Set("danger_levels", securityAuditLogExportTask.DangerLevels)
	}

	return nil
}

func resourceTencentCloudDbbrainSecurityAuditLogExportTaskDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dbbrain_security_audit_log_export_task.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DbbrainService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
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
