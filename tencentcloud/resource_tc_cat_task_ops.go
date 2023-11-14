/*
Provides a resource to create a cat task_ops

Example Usage

```hcl
resource "tencentcloud_cat_task_ops" "task_ops" {
  task_ids =
}
```

Import

cat task_ops can be imported using the id, e.g.

```
terraform import tencentcloud_cat_task_ops.task_ops task_ops_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cat/v20180409"
	"log"
)

func resourceTencentCloudCatTaskOps() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCatTaskOpsCreate,
		Read:   resourceTencentCloudCatTaskOpsRead,
		Delete: resourceTencentCloudCatTaskOpsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"task_ids": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Task Id.",
			},
		},
	}
}

func resourceTencentCloudCatTaskOpsCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cat_task_ops.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = cat.NewSuspendProbeTaskRequest()
		response = cat.NewSuspendProbeTaskResponse()
		taskId   string
	)
	if v, ok := d.GetOk("task_ids"); ok {
		taskIdsSet := v.(*schema.Set).List()
		for i := range taskIdsSet {
			taskIds := taskIdsSet[i].(string)
			request.TaskIds = append(request.TaskIds, &taskIds)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCatClient().SuspendProbeTask(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cat taskOps failed, reason:%+v", logId, err)
		return err
	}

	taskId = *response.Response.TaskId
	d.SetId(taskId)

	return resourceTencentCloudCatTaskOpsRead(d, meta)
}

func resourceTencentCloudCatTaskOpsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cat_task_ops.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CatService{client: meta.(*TencentCloudClient).apiV3Conn}

	taskOpsId := d.Id()

	taskOps, err := service.DescribeCatTaskOpsById(ctx, taskId)
	if err != nil {
		return err
	}

	if taskOps == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CatTaskOps` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if taskOps.TaskIds != nil {
		_ = d.Set("task_ids", taskOps.TaskIds)
	}

	return nil
}

func resourceTencentCloudCatTaskOpsDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cat_task_ops.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CatService{client: meta.(*TencentCloudClient).apiV3Conn}
	taskOpsId := d.Id()

	if err := service.DeleteCatTaskOpsById(ctx, taskId); err != nil {
		return err
	}

	return nil
}
