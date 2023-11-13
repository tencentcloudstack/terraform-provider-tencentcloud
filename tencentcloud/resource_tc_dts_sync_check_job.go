/*
Provides a resource to create a dts sync_check_job

Example Usage

```hcl
resource "tencentcloud_dts_sync_check_job" "sync_check_job" {
  job_id = &lt;nil&gt;
}
```

Import

dts sync_check_job can be imported using the id, e.g.

```
terraform import tencentcloud_dts_sync_check_job.sync_check_job sync_check_job_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

func resourceTencentCloudDtsSyncCheckJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDtsSyncCheckJobCreate,
		Read:   resourceTencentCloudDtsSyncCheckJobRead,
		Update: resourceTencentCloudDtsSyncCheckJobUpdate,
		Delete: resourceTencentCloudDtsSyncCheckJobDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"job_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Synchronization task id.",
			},
		},
	}
}

func resourceTencentCloudDtsSyncCheckJobCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_sync_check_job.create")()
	defer inconsistentCheck(d, meta)()

	var jobId string
	if v, ok := d.GetOk("job_id"); ok {
		jobId = v.(string)
	}

	d.SetId(jobId)

	return resourceTencentCloudDtsSyncCheckJobUpdate(d, meta)
}

func resourceTencentCloudDtsSyncCheckJobRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_sync_check_job.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	syncCheckJobId := d.Id()

	syncCheckJob, err := service.DescribeDtsSyncCheckJobById(ctx, jobId)
	if err != nil {
		return err
	}

	if syncCheckJob == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DtsSyncCheckJob` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if syncCheckJob.JobId != nil {
		_ = d.Set("job_id", syncCheckJob.JobId)
	}

	return nil
}

func resourceTencentCloudDtsSyncCheckJobUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_sync_check_job.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	immutableArgs := []string{"job_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	return resourceTencentCloudDtsSyncCheckJobRead(d, meta)
}

func resourceTencentCloudDtsSyncCheckJobDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_sync_check_job.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
