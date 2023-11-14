/*
Provides a resource to create a dts migrate_job_operate

Example Usage

```hcl
resource "tencentcloud_dts_migrate_job_operate" "migrate_job_operate" {
  job_id = "dts-ekmhr27i"
  complete_mode = "immediately"
}
```

Import

dts migrate_job_operate can be imported using the id, e.g.

```
terraform import tencentcloud_dts_migrate_job_operate.migrate_job_operate migrate_job_operate_id
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
	"time"
)

func resourceTencentCloudDtsMigrateJobOperate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDtsMigrateJobOperateCreate,
		Read:   resourceTencentCloudDtsMigrateJobOperateRead,
		Update: resourceTencentCloudDtsMigrateJobOperateUpdate,
		Delete: resourceTencentCloudDtsMigrateJobOperateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"job_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Job id.",
			},

			"complete_mode": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Complete mode, optional value is waitForSync or immediately.",
			},
		},
	}
}

func resourceTencentCloudDtsMigrateJobOperateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_job_operate.create")()
	defer inconsistentCheck(d, meta)()

	var jobId string
	if v, ok := d.GetOk("job_id"); ok {
		jobId = v.(string)
	}

	d.SetId(jobId)

	return resourceTencentCloudDtsMigrateJobOperateUpdate(d, meta)
}

func resourceTencentCloudDtsMigrateJobOperateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_job_operate.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	migrateJobOperateId := d.Id()

	migrateJobOperate, err := service.DescribeDtsMigrateJobOperateById(ctx, jobId)
	if err != nil {
		return err
	}

	if migrateJobOperate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DtsMigrateJobOperate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if migrateJobOperate.JobId != nil {
		_ = d.Set("job_id", migrateJobOperate.JobId)
	}

	if migrateJobOperate.CompleteMode != nil {
		_ = d.Set("complete_mode", migrateJobOperate.CompleteMode)
	}

	return nil
}

func resourceTencentCloudDtsMigrateJobOperateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_job_operate.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		pauseMigrateJobRequest  = dts.NewPauseMigrateJobRequest()
		pauseMigrateJobResponse = dts.NewPauseMigrateJobResponse()
	)

	migrateJobOperateId := d.Id()

	request.JobId = &jobId

	immutableArgs := []string{"job_id", "complete_mode"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDtsClient().PauseMigrateJob(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update dts migrateJobOperate failed, reason:%+v", logId, err)
		return err
	}

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"manualPaused"}, 2*readRetryTimeout, time.Second, service.DtsMigrateJobOperateStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"running"}, 2*readRetryTimeout, time.Second, service.DtsMigrateJobOperateStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"success"}, 3*readRetryTimeout, time.Second, service.DtsMigrateJobOperateStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"running"}, 2*readRetryTimeout, time.Second, service.DtsMigrateJobOperateStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"canceled"}, 2*readRetryTimeout, time.Second, service.DtsMigrateJobOperateStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudDtsMigrateJobOperateRead(d, meta)
}

func resourceTencentCloudDtsMigrateJobOperateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_job_operate.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
