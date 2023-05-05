/*
Provides a resource to create a dts compare_task_stop_operation

Example Usage

```hcl
resource "tencentcloud_dts_compare_task_stop_operation" "compare_task_stop_operation" {
  job_id = "dts-8yv4w2i1"
  compare_task_id = "dts-8yv4w2i1-cmp-37skmii9"
}
```

*/
package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDtsCompareTaskStopOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDtsCompareTaskStopOperationCreate,
		Read:   resourceTencentCloudDtsCompareTaskStopOperationRead,
		Delete: resourceTencentCloudDtsCompareTaskStopOperationDelete,
		Schema: map[string]*schema.Schema{
			"job_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "job id.",
			},

			"compare_task_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Compare task id.",
			},
		},
	}
}

func resourceTencentCloudDtsCompareTaskStopOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_compare_task_stop_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request       = dts.NewStopCompareRequest()
		jobId         string
		compareTaskId string
	)
	if v, ok := d.GetOk("job_id"); ok {
		jobId = v.(string)
		request.JobId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("compare_task_id"); ok {
		compareTaskId = v.(string)
		request.CompareTaskId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("job_id"); ok {
		request.JobId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDtsClient().StopCompare(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate dts compareTaskStopOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(jobId + FILED_SP + compareTaskId)

	return resourceTencentCloudDtsCompareTaskStopOperationRead(d, meta)
}

func resourceTencentCloudDtsCompareTaskStopOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_compare_task_stop_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDtsCompareTaskStopOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_compare_task_stop_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
