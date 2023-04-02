/*
Provides a resource to create a dts sync_check_job_operation

Example Usage

```hcl
resource "tencentcloud_dts_sync_check_job_operation" "sync_check_job_operation" {
  job_id = ""
  }
```
*/
package tencentcloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	dts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDtsSyncCheckJobOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDtsSyncCheckJobOperationCreate,
		Read:   resourceTencentCloudDtsSyncCheckJobOperationRead,
		Update: resourceTencentCloudDtsSyncCheckJobOperationUpdate,
		Delete: resourceTencentCloudDtsSyncCheckJobOperationDelete,
		Schema: map[string]*schema.Schema{
			"job_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Sync job id.",
			},

			// "status": {
			// 	Computed:    true,
			// 	Type:        schema.TypeString,
			// 	Description: "The execution status of the verification task, such as: notStarted (not started), running (verifying), failed (verification task failed), success (task successful).",
			// },
		},
	}
}

func resourceTencentCloudDtsSyncCheckJobOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_sync_check_job_operation.create")()
	defer inconsistentCheck(d, meta)()

	var jobId string

	if v, ok := d.GetOk("job_id"); ok {
		jobId = v.(string)
	}
	d.SetId(jobId)

	return resourceTencentCloudDtsSyncCheckJobOperationUpdate(d, meta)
}

func resourceTencentCloudDtsSyncCheckJobOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_sync_check_job_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDtsSyncCheckJobOperationUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_sync_check_job_operation.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := dts.NewCreateCheckSyncJobRequest()

	request.JobId = helper.String(d.Id())

	immutableArgs := []string{"job_id", "status"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDtsClient().CreateCheckSyncJob(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update dts syncCheckJobOperation failed, reason:%+v", logId, err)
		return err
	}

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"failed,success"}, 0*readRetryTimeout, time.Second, service.DtsSyncCheckJobOperationStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudDtsSyncCheckJobOperationRead(d, meta)
}

func resourceTencentCloudDtsSyncCheckJobOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_sync_check_job_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
