/*
Provides a resource to create a dbbrain security_audit_log_export_task

Example Usage

```hcl
resource "tencentcloud_dbbrain_security_audit_log_export_task" "security_audit_log_export_task" {
  sec_audit_group_id = &lt;nil&gt;
  start_time = &lt;nil&gt;
  end_time = &lt;nil&gt;
  product = &lt;nil&gt;
  danger_levels = &lt;nil&gt;
}
```

Import

dbbrain security_audit_log_export_task can be imported using the id, e.g.

```
terraform import tencentcloud_dbbrain_security_audit_log_export_task.security_audit_log_export_task security_audit_log_export_task_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudDbbrainSecurityAuditLogExportTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDbbrainSecurityAuditLogExportTaskCreate,
		Read:   resourceTencentCloudDbbrainSecurityAuditLogExportTaskRead,
		Update: resourceTencentCloudDbbrainSecurityAuditLogExportTaskUpdate,
		Delete: resourceTencentCloudDbbrainSecurityAuditLogExportTaskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"sec_audit_group_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Security audit group id.",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Start time.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "End time.",
			},

			"product": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Product, optional value is mysql.",
			},

			"danger_levels": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Danger level list.",
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
		response        = dbbrain.NewCreateSecurityAuditLogExportTaskResponse()
		secAuditGroupId string
		asyncRequestId  int
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
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dbbrain securityAuditLogExportTask failed, reason:%+v", logId, err)
		return err
	}

	secAuditGroupId = *response.Response.SecAuditGroupId
	d.SetId(strings.Join([]string{secAuditGroupId, helper.Int64ToStr(asyncRequestId)}, FILED_SP))

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
	secAuditGroupId := idSplit[0]
	asyncRequestId := idSplit[1]

	securityAuditLogExportTask, err := service.DescribeDbbrainSecurityAuditLogExportTaskById(ctx, secAuditGroupId, asyncRequestId)
	if err != nil {
		return err
	}

	if securityAuditLogExportTask == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DbbrainSecurityAuditLogExportTask` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if securityAuditLogExportTask.SecAuditGroupId != nil {
		_ = d.Set("sec_audit_group_id", securityAuditLogExportTask.SecAuditGroupId)
	}

	if securityAuditLogExportTask.StartTime != nil {
		_ = d.Set("start_time", securityAuditLogExportTask.StartTime)
	}

	if securityAuditLogExportTask.EndTime != nil {
		_ = d.Set("end_time", securityAuditLogExportTask.EndTime)
	}

	if securityAuditLogExportTask.Product != nil {
		_ = d.Set("product", securityAuditLogExportTask.Product)
	}

	if securityAuditLogExportTask.DangerLevels != nil {
		_ = d.Set("danger_levels", securityAuditLogExportTask.DangerLevels)
	}

	return nil
}

func resourceTencentCloudDbbrainSecurityAuditLogExportTaskUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dbbrain_security_audit_log_export_task.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	immutableArgs := []string{"sec_audit_group_id", "start_time", "end_time", "product", "danger_levels"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	return resourceTencentCloudDbbrainSecurityAuditLogExportTaskRead(d, meta)
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
	secAuditGroupId := idSplit[0]
	asyncRequestId := idSplit[1]

	if err := service.DeleteDbbrainSecurityAuditLogExportTaskById(ctx, secAuditGroupId, asyncRequestId); err != nil {
		return err
	}

	return nil
}
