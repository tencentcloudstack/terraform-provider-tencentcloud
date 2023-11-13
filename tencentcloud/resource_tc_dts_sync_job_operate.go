/*
Provides a resource to create a dts sync_job_operate

Example Usage

```hcl
resource "tencentcloud_dts_sync_job_operate" "sync_job_operate" {
  job_id = "sync-werwfs23"
  new_instance_class = "large"
}
```

Import

dts sync_job_operate can be imported using the id, e.g.

```
terraform import tencentcloud_dts_sync_job_operate.sync_job_operate sync_job_operate_id
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

func resourceTencentCloudDtsSyncJobOperate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDtsSyncJobOperateCreate,
		Read:   resourceTencentCloudDtsSyncJobOperateRead,
		Update: resourceTencentCloudDtsSyncJobOperateUpdate,
		Delete: resourceTencentCloudDtsSyncJobOperateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"job_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Synchronization instance id (i.e. identifies a synchronization job).",
			},

			"new_instance_class": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Task specification.",
			},
		},
	}
}

func resourceTencentCloudDtsSyncJobOperateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_sync_job_operate.create")()
	defer inconsistentCheck(d, meta)()

	var jobId string
	if v, ok := d.GetOk("job_id"); ok {
		jobId = v.(string)
	}

	d.SetId(jobId)

	return resourceTencentCloudDtsSyncJobOperateUpdate(d, meta)
}

func resourceTencentCloudDtsSyncJobOperateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_sync_job_operate.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	syncJobOperateId := d.Id()

	syncJobOperate, err := service.DescribeDtsSyncJobOperateById(ctx, jobId)
	if err != nil {
		return err
	}

	if syncJobOperate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DtsSyncJobOperate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if syncJobOperate.JobId != nil {
		_ = d.Set("job_id", syncJobOperate.JobId)
	}

	if syncJobOperate.NewInstanceClass != nil {
		_ = d.Set("new_instance_class", syncJobOperate.NewInstanceClass)
	}

	return nil
}

func resourceTencentCloudDtsSyncJobOperateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_sync_job_operate.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		pauseSyncJobRequest  = dts.NewPauseSyncJobRequest()
		pauseSyncJobResponse = dts.NewPauseSyncJobResponse()
	)

	syncJobOperateId := d.Id()

	request.JobId = &jobId

	immutableArgs := []string{"job_id", "new_instance_class"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDtsClient().PauseSyncJob(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update dts syncJobOperate failed, reason:%+v", logId, err)
		return err
	}

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"Paused"}, 2*readRetryTimeout, time.Second, service.DtsSyncJobOperateStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"Running"}, 2*readRetryTimeout, time.Second, service.DtsSyncJobOperateStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"Running"}, 2*readRetryTimeout, time.Second, service.DtsSyncJobOperateStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"Running"}, 2*readRetryTimeout, time.Second, service.DtsSyncJobOperateStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"Running"}, 2*readRetryTimeout, time.Second, service.DtsSyncJobOperateStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"Stopped"}, 2*readRetryTimeout, time.Second, service.DtsSyncJobOperateStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudDtsSyncJobOperateRead(d, meta)
}

func resourceTencentCloudDtsSyncJobOperateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_sync_job_operate.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
