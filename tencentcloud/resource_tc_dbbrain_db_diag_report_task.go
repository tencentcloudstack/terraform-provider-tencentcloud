/*
Provides a resource to create a dbbrain db_diag_report_task

Example Usage

```hcl
resource "tencentcloud_dbbrain_db_diag_report_task" "db_diag_report_task" {
  instance_id = ""
  start_time = ""
  end_time = ""
  send_mail_flag =
  contact_person =
  contact_group =
  product = ""
}
```

Import

dbbrain db_diag_report_task can be imported using the id, e.g.

```
terraform import tencentcloud_dbbrain_db_diag_report_task.db_diag_report_task db_diag_report_task_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudDbbrainDbDiagReportTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDbbrainDbDiagReportTaskCreate,
		Read:   resourceTencentCloudDbbrainDbDiagReportTaskRead,
		Update: resourceTencentCloudDbbrainDbDiagReportTaskUpdate,
		Delete: resourceTencentCloudDbbrainDbDiagReportTaskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "instance id.",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Start time, such as &amp;quot;2020-11-08T14:00:00+08:00&amp;quot;.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "End time, such as &amp;quot;2020-11-09T14:00:00+08:00&amp;quot;.",
			},

			"send_mail_flag": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Whether to send mail: 0 - no, 1 - yes. .",
			},

			"contact_person": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "An array of contact IDs to receive emails from.",
			},

			"contact_group": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "An array of contact group IDs to receive mail from.",
			},

			"product": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Service product type, supported values include： &amp;quot;mysql&amp;quot; - cloud database MySQL, &amp;quot;cynosdb&amp;quot; - cloud database CynosDB for MySQL, the default value is &amp;quot;mysql&amp;quot;.",
			},
		},
	}
}

func resourceTencentCloudDbbrainDbDiagReportTaskCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dbbrain_db_diag_report_task.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request        = dbbrain.NewCreateDBDiagReportTaskRequest()
		response       = dbbrain.NewCreateDBDiagReportTaskResponse()
		asyncRequestId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("start_time"); ok {
		request.StartTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		request.EndTime = helper.String(v.(string))
	}

	if v, _ := d.GetOk("send_mail_flag"); v != nil {
		request.SendMailFlag = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("contact_person"); ok {
		contactPersonSet := v.(*schema.Set).List()
		for i := range contactPersonSet {
			contactPerson := contactPersonSet[i].(int)
			request.ContactPerson = append(request.ContactPerson, helper.IntInt64(contactPerson))
		}
	}

	if v, ok := d.GetOk("contact_group"); ok {
		contactGroupSet := v.(*schema.Set).List()
		for i := range contactGroupSet {
			contactGroup := contactGroupSet[i].(int)
			request.ContactGroup = append(request.ContactGroup, helper.IntInt64(contactGroup))
		}
	}

	if v, ok := d.GetOk("product"); ok {
		request.Product = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDbbrainClient().CreateDBDiagReportTask(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dbbrain dbDiagReportTask failed, reason:%+v", logId, err)
		return err
	}

	asyncRequestId = *response.Response.AsyncRequestId
	d.SetId(helper.String(asyncRequestId))

	return resourceTencentCloudDbbrainDbDiagReportTaskRead(d, meta)
}

func resourceTencentCloudDbbrainDbDiagReportTaskRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dbbrain_db_diag_report_task.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DbbrainService{client: meta.(*TencentCloudClient).apiV3Conn}

	dbDiagReportTaskId := d.Id()

	dbDiagReportTask, err := service.DescribeDbbrainDbDiagReportTaskById(ctx, asyncRequestId)
	if err != nil {
		return err
	}

	if dbDiagReportTask == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DbbrainDbDiagReportTask` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if dbDiagReportTask.InstanceId != nil {
		_ = d.Set("instance_id", dbDiagReportTask.InstanceId)
	}

	if dbDiagReportTask.StartTime != nil {
		_ = d.Set("start_time", dbDiagReportTask.StartTime)
	}

	if dbDiagReportTask.EndTime != nil {
		_ = d.Set("end_time", dbDiagReportTask.EndTime)
	}

	if dbDiagReportTask.SendMailFlag != nil {
		_ = d.Set("send_mail_flag", dbDiagReportTask.SendMailFlag)
	}

	if dbDiagReportTask.ContactPerson != nil {
		_ = d.Set("contact_person", dbDiagReportTask.ContactPerson)
	}

	if dbDiagReportTask.ContactGroup != nil {
		_ = d.Set("contact_group", dbDiagReportTask.ContactGroup)
	}

	if dbDiagReportTask.Product != nil {
		_ = d.Set("product", dbDiagReportTask.Product)
	}

	return nil
}

func resourceTencentCloudDbbrainDbDiagReportTaskUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dbbrain_db_diag_report_task.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	immutableArgs := []string{"instance_id", "start_time", "end_time", "send_mail_flag", "contact_person", "contact_group", "product"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	return resourceTencentCloudDbbrainDbDiagReportTaskRead(d, meta)
}

func resourceTencentCloudDbbrainDbDiagReportTaskDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dbbrain_db_diag_report_task.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DbbrainService{client: meta.(*TencentCloudClient).apiV3Conn}
	dbDiagReportTaskId := d.Id()

	if err := service.DeleteDbbrainDbDiagReportTaskById(ctx, asyncRequestId); err != nil {
		return err
	}

	return nil
}
