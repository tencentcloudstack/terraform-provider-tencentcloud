/*
Provides a resource to create a dts compare_task_operate

Example Usage

```hcl
resource "tencentcloud_dts_compare_task_operate" "compare_task_operate" {
  job_id = "dts-8yv4w2i1"
  compare_task_id = "dts-8yv4w2i1-cmp-37skmii9"
}
```

Import

dts compare_task_operate can be imported using the id, e.g.

```
terraform import tencentcloud_dts_compare_task_operate.compare_task_operate compare_task_operate_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206"
	"log"
	"strings"
)

func resourceTencentCloudDtsCompareTaskOperate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDtsCompareTaskOperateCreate,
		Read:   resourceTencentCloudDtsCompareTaskOperateRead,
		Update: resourceTencentCloudDtsCompareTaskOperateUpdate,
		Delete: resourceTencentCloudDtsCompareTaskOperateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"job_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Job id.",
			},

			"compare_task_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Compare task id.",
			},
		},
	}
}

func resourceTencentCloudDtsCompareTaskOperateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_compare_task_operate.create")()
	defer inconsistentCheck(d, meta)()

	var jobId string
	if v, ok := d.GetOk("job_id"); ok {
		jobId = v.(string)
	}

	var compareTaskId string
	if v, ok := d.GetOk("compare_task_id"); ok {
		compareTaskId = v.(string)
	}

	d.SetId(strings.Join([]string{jobId, compareTaskId}, FILED_SP))

	return resourceTencentCloudDtsCompareTaskOperateUpdate(d, meta)
}

func resourceTencentCloudDtsCompareTaskOperateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_compare_task_operate.read")()
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

	compareTaskOperate, err := service.DescribeDtsCompareTaskOperateById(ctx, jobId, compareTaskId)
	if err != nil {
		return err
	}

	if compareTaskOperate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DtsCompareTaskOperate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if compareTaskOperate.JobId != nil {
		_ = d.Set("job_id", compareTaskOperate.JobId)
	}

	if compareTaskOperate.CompareTaskId != nil {
		_ = d.Set("compare_task_id", compareTaskOperate.CompareTaskId)
	}

	return nil
}

func resourceTencentCloudDtsCompareTaskOperateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_compare_task_operate.update")()
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
		log.Printf("[CRITAL]%s update dts compareTaskOperate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudDtsCompareTaskOperateRead(d, meta)
}

func resourceTencentCloudDtsCompareTaskOperateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_compare_task_operate.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
