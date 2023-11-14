/*
Provides a resource to create a wedata integration_offline_task

Example Usage

```hcl
resource "tencentcloud_wedata_integration_offline_task" "integration_offline_task" {
  project_id = "1455251608631480391"
  cycle_step = 1
  delay_time = 0
  end_time = "2099-12-31 00:00:00"
  notes = "Task for test"
  start_time = "2023-12-31 00:00:00"
  task_name = "TaskTest_10"
  type_id = 27
  task_action = "0,3,4"
  task_mode = "1"
}
```

Import

wedata integration_offline_task can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_integration_offline_task.integration_offline_task integration_offline_task_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedata "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20210820"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudWedataIntegration_offline_task() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataIntegration_offline_taskCreate,
		Read:   resourceTencentCloudWedataIntegration_offline_taskRead,
		Update: resourceTencentCloudWedataIntegration_offline_taskUpdate,
		Delete: resourceTencentCloudWedataIntegration_offline_taskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Project ID.",
			},

			"cycle_step": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Interval time of scheduling, the minimum value: 1.",
			},

			"delay_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Execution time, unit is minutes, only available for day/week/month/year scheduling. For example, daily scheduling is executed once every day at 02:00, and the delayTime is 120 minutes.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Effective end time, the format is yyyy-MM-dd HH:mm:ss.",
			},

			"notes": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Description information.",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Effective start time, the format is yyyy-MM-dd HH:mm:ss.",
			},

			"task_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Task name.",
			},

			"type_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Task type ID, for intgration task the value is 27.",
			},

			"task_action": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Scheduling configuration: flexible period configuration, only available for hourly/weekly/monthly/yearly scheduling. If the hourly task is specified to run at 0:00, 3:00 and 4:00 every day, it is &amp;amp;#39;0,3,4&amp;amp;#39;.",
			},

			"task_mode": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Task display mode, 0: canvas mode, 1: form mode.",
			},
		},
	}
}

func resourceTencentCloudWedataIntegration_offline_taskCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_integration_offline_task.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = wedata.NewCreateOfflineTaskRequest()
		response = wedata.NewCreateOfflineTaskResponse()
		taskId   string
	)
	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("cycle_step"); ok {
		request.CycleStep = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("delay_time"); ok {
		request.DelayTime = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("end_time"); ok {
		request.EndTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("notes"); ok {
		request.Notes = helper.String(v.(string))
	}

	if v, ok := d.GetOk("start_time"); ok {
		request.StartTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("task_name"); ok {
		request.TaskName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("type_id"); ok {
		request.TypeId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("task_action"); ok {
		request.TaskAction = helper.String(v.(string))
	}

	if v, ok := d.GetOk("task_mode"); ok {
		request.TaskMode = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWedataClient().CreateOfflineTask(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create wedata integration_offline_task failed, reason:%+v", logId, err)
		return err
	}

	taskId = *response.Response.TaskId
	d.SetId(taskId)

	return resourceTencentCloudWedataIntegration_offline_taskRead(d, meta)
}

func resourceTencentCloudWedataIntegration_offline_taskRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_integration_offline_task.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := WedataService{client: meta.(*TencentCloudClient).apiV3Conn}

	integration_offline_taskId := d.Id()

	integration_offline_task, err := service.DescribeWedataIntegration_offline_taskById(ctx, taskId)
	if err != nil {
		return err
	}

	if integration_offline_task == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `WedataIntegration_offline_task` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if integration_offline_task.ProjectId != nil {
		_ = d.Set("project_id", integration_offline_task.ProjectId)
	}

	if integration_offline_task.CycleStep != nil {
		_ = d.Set("cycle_step", integration_offline_task.CycleStep)
	}

	if integration_offline_task.DelayTime != nil {
		_ = d.Set("delay_time", integration_offline_task.DelayTime)
	}

	if integration_offline_task.EndTime != nil {
		_ = d.Set("end_time", integration_offline_task.EndTime)
	}

	if integration_offline_task.Notes != nil {
		_ = d.Set("notes", integration_offline_task.Notes)
	}

	if integration_offline_task.StartTime != nil {
		_ = d.Set("start_time", integration_offline_task.StartTime)
	}

	if integration_offline_task.TaskName != nil {
		_ = d.Set("task_name", integration_offline_task.TaskName)
	}

	if integration_offline_task.TypeId != nil {
		_ = d.Set("type_id", integration_offline_task.TypeId)
	}

	if integration_offline_task.TaskAction != nil {
		_ = d.Set("task_action", integration_offline_task.TaskAction)
	}

	if integration_offline_task.TaskMode != nil {
		_ = d.Set("task_mode", integration_offline_task.TaskMode)
	}

	return nil
}

func resourceTencentCloudWedataIntegration_offline_taskUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_integration_offline_task.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := wedata.NewModifyIntegrationTaskRequest()

	integration_offline_taskId := d.Id()

	request.TaskId = &taskId

	immutableArgs := []string{"project_id", "cycle_step", "delay_time", "end_time", "notes", "start_time", "task_name", "type_id", "task_action", "task_mode"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("project_id") {
		if v, ok := d.GetOk("project_id"); ok {
			request.ProjectId = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWedataClient().ModifyIntegrationTask(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update wedata integration_offline_task failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudWedataIntegration_offline_taskRead(d, meta)
}

func resourceTencentCloudWedataIntegration_offline_taskDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_integration_offline_task.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := WedataService{client: meta.(*TencentCloudClient).apiV3Conn}
	integration_offline_taskId := d.Id()

	if err := service.DeleteWedataIntegration_offline_taskById(ctx, taskId); err != nil {
		return err
	}

	return nil
}
