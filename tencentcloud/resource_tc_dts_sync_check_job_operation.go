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
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDtsSyncCheckJobOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDtsSyncCheckJobOperationCreate,
		Read:   resourceTencentCloudDtsSyncCheckJobOperationRead,
		Delete: resourceTencentCloudDtsSyncCheckJobOperationDelete,
		Schema: map[string]*schema.Schema{
			"job_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Sync job id.",
			},
		},
	}
}

func resourceTencentCloudDtsSyncCheckJobOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_sync_check_job_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request = dts.NewCreateCheckSyncJobRequest()
		jobId   string
	)
	if v, ok := d.GetOk("job_id"); ok {
		jobId = v.(string)
		request.JobId = helper.String(v.(string))
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
		log.Printf("[CRITAL]%s operate dts syncCheckJobOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(jobId)

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"failed", "success"}, readRetryTimeout, time.Second, service.DtsSyncCheckJobOperationStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudDtsSyncCheckJobOperationRead(d, meta)
}

func resourceTencentCloudDtsSyncCheckJobOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_sync_check_job_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDtsSyncCheckJobOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_sync_check_job_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
