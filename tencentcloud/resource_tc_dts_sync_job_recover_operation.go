package tencentcloud

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDtsSyncJobRecoverOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDtsSyncJobRecoverOperationCreate,
		Read:   resourceTencentCloudDtsSyncJobRecoverOperationRead,
		Delete: resourceTencentCloudDtsSyncJobRecoverOperationDelete,
		Schema: map[string]*schema.Schema{
			"job_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Synchronization instance id (i.e. identifies a synchronization job).",
			},
		},
	}
}

func resourceTencentCloudDtsSyncJobRecoverOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_sync_job_recover_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request = dts.NewRecoverSyncJobRequest()
		jobId   string
	)
	if v, ok := d.GetOk("job_id"); ok {
		request.JobId = helper.String(v.(string))
		jobId = v.(string)
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDtsClient().RecoverSyncJob(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate dts syncJobRecoverOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(jobId)

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"Running", "Stopped"}, 2*readRetryTimeout, time.Second, service.DtsSyncJobStateRefreshFunc(d.Id(), "Stopped", []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	conf = BuildStateChangeConf([]string{}, []string{"Normal"}, 2*readRetryTimeout, time.Second, service.DtsSyncJobTradeStateRefreshFunc(d.Id(), "", []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudDtsSyncJobRecoverOperationRead(d, meta)
}

func resourceTencentCloudDtsSyncJobRecoverOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_sync_job_recover_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDtsSyncJobRecoverOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_sync_job_recover_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
