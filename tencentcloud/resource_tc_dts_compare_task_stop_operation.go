/*
Provides a resource to create a dts compare_task_stop_operation

Example Usage

```hcl
resource "tencentcloud_dts_compare_task_stop_operation" "compare_task_stop_operation" {
  job_id = "dts-8yv4w2i1"
  compare_task_id = "dts-8yv4w2i1-cmp-37skmii9"
}
```

Import

dts compare_task_stop_operation can be imported using the id, e.g.

```
terraform import tencentcloud_dts_compare_task_stop_operation.compare_task_stop_operation compare_task_stop_operation_id
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
	dts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206"
)

func resourceTencentCloudDtsCompareTaskStopOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDtsCompareTaskStopOperationCreate,
		Read:   resourceTencentCloudDtsCompareTaskStopOperationRead,
		Update: resourceTencentCloudDtsCompareTaskStopOperationUpdate,
		Delete: resourceTencentCloudDtsCompareTaskStopOperationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"job_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "job id.",
			},

			"compare_task_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Compare task id.",
			},
		},
	}
}

func resourceTencentCloudDtsCompareTaskStopOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_compare_task_stop_operation.create")()
	defer inconsistentCheck(d, meta)()

	var (
		jobId         string
		compareTaskId string
	)

	if v, ok := d.GetOk("job_id"); ok {
		jobId = v.(string)
	}

	if v, ok := d.GetOk("compare_task_id"); ok {
		compareTaskId = v.(string)
	}

	d.SetId(jobId + FILED_SP + compareTaskId)

	return resourceTencentCloudDtsCompareTaskStopOperationUpdate(d, meta)
}

func resourceTencentCloudDtsCompareTaskStopOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_compare_task_stop_operation.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	jobId := idSplit[0]
	compareTaskId := idSplit[1]

	ret, err := service.DescribeDtsCompareTaskStopOperationById(ctx, jobId, compareTaskId)
	if err != nil {
		return err
	}

	if ret == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DtsCompareTaskStopOperation` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	if ret.Abstract.Status == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DescribeCompareReportResponseParams.Abstract.Status` [%s] not found, please check it.\n", logId, d.Id())
		return nil
	}

	return nil
}

func resourceTencentCloudDtsCompareTaskStopOperationUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_compare_task_stop_operation.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := dts.NewStopCompareRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	jobId := idSplit[0]
	compareTaskId := idSplit[1]

	request.JobId = &jobId
	request.CompareTaskId = &compareTaskId

	immutableArgs := []string{"job_id", "compare_task_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
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
		log.Printf("[CRITAL]%s update dts compareTaskStopOperation failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudDtsCompareTaskStopOperationRead(d, meta)
}

func resourceTencentCloudDtsCompareTaskStopOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_compare_task_stop_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
