/*
Provides a resource to create a dts migrate_job_config

Example Usage

```hcl
resource "tencentcloud_dts_migrate_job_config" "migrate_job_config" {
  job_id = "dts-ekmhr27i"
  complete_mode = "immediately"
  action = "dts-ekmhr27i"
}
```

*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	dts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDtsMigrateJobConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDtsMigrateJobConfigCreate,
		Read:   resourceTencentCloudDtsMigrateJobConfigRead,
		Update: resourceTencentCloudDtsMigrateJobConfigUpdate,
		Delete: resourceTencentCloudDtsMigrateJobConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"job_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "job id.",
			},

			"complete_mode": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "complete mode, optional value is waitForSync or immediately.",
			},

			"action": {
				Required:     true,
				Type:         schema.TypeString,
				Description:  "The operation want to perform. Valid values are: `pause`, `continue`, `complete`, `recover`,`stop`.",
				ValidateFunc: validateAllowedStringValue([]string{DTS_MIGRATE_ACTION_PAUSE, DTS_MIGRATE_ACTION_CONTINUE, DTS_MIGRATE_ACTION_COMPLETE, DTS_MIGRATE_ACTION_RECOVER, DTS_MIGRATE_ACTION_STOP, DTS_MIGRATE_ACTION_ISOLATE}),
			},
		},
	}
}

func resourceTencentCloudDtsMigrateJobConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_job_config.create")()
	defer inconsistentCheck(d, meta)()

	var jobId string

	if v, ok := d.GetOk("job_id"); ok {
		jobId = v.(string)
	}
	d.SetId(jobId)

	return resourceTencentCloudDtsMigrateJobConfigUpdate(d, meta)
}

func resourceTencentCloudDtsMigrateJobConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_job_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	jobId := d.Id()

	migrateJobConfig, err := service.DescribeDtsMigrateJobById(ctx, jobId)
	if err != nil {
		return err
	}

	if migrateJobConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DtsMigrateJobConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if migrateJobConfig.JobId != nil {
		_ = d.Set("job_id", migrateJobConfig.JobId)
	}

	// if migrateJobConfig.RunMode != nil {
	// 	_ = d.Set("complete_mode", migrateJobConfig.RunMode)
	// }

	return nil
}

func resourceTencentCloudDtsMigrateJobConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_job_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		action string
	)

	jobId := d.Id()

	if d.HasChange("action") {
		if v, ok := d.GetOk("action"); ok {
			action = v.(string)
			var inErr error
			switch action {
			case DTS_MIGRATE_ACTION_PAUSE:
				inErr = handlePauseMigrate(d, meta, logId, jobId)
				if inErr != nil {
					return inErr
				}
			case DTS_MIGRATE_ACTION_CONTINUE:
				inErr = handleContinueMigrate(d, meta, logId, jobId)
				if inErr != nil {
					return inErr
				}
			case DTS_MIGRATE_ACTION_COMPLETE:
				inErr = handleCompleteMigrate(d, meta, logId, jobId)
				if inErr != nil {
					return inErr
				}
			case DTS_MIGRATE_ACTION_RECOVER:
				inErr = handleRecoverMigrate(d, meta, logId, jobId)
				if inErr != nil {
					return inErr
				}
			case DTS_MIGRATE_ACTION_STOP:
				inErr = handleStopMigrate(d, meta, logId, jobId)
				if inErr != nil {
					return inErr
				}
			case DTS_MIGRATE_ACTION_ISOLATE:
				inErr = handleIsolateMigrate(d, meta, logId, jobId)
				if inErr != nil {
					return inErr
				}
			default:
				return fmt.Errorf("invalid action: %s", action)
			} // switch end
		}
	}

	return resourceTencentCloudDtsMigrateJobConfigRead(d, meta)
}

func resourceTencentCloudDtsMigrateJobConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_job_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func handlePauseMigrate(d *schema.ResourceData, meta interface{}, logId, jobId string) error {
	request := dts.NewPauseMigrateJobRequest()
	request.JobId = helper.String(jobId)
	// response = dts.NewPauseMigrateJobResponse()

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
		log.Printf("[CRITAL]%s update dts migrateJobConfig failed, reason:%+v", logId, err)
		return err
	}

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"manualPaused"}, 2*readRetryTimeout, time.Second, service.DtsMigrateJobConfigStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}
	return nil
}

func handleContinueMigrate(d *schema.ResourceData, meta interface{}, logId, jobId string) error {
	request := dts.NewContinueMigrateJobRequest()
	request.JobId = helper.String(jobId)
	// response = dts.NewPauseMigrateJobResponse()

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDtsClient().ContinueMigrateJob(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update dts migrateJobConfig failed, reason:%+v", logId, err)
		return err
	}

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"running"}, 2*readRetryTimeout, time.Second, service.DtsMigrateJobConfigStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}

func handleCompleteMigrate(d *schema.ResourceData, meta interface{}, logId, jobId string) error {
	request := dts.NewCompleteMigrateJobRequest()
	request.JobId = helper.String(jobId)
	// response = dts.NewPauseMigrateJobResponse()

	if d.HasChange("complete_mode") {
		if v, ok := d.GetOk("complete_mode"); ok {
			request.CompleteMode = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDtsClient().CompleteMigrateJob(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update dts migrateJobConfig failed, reason:%+v", logId, err)
		return err
	}

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"success"}, 3*readRetryTimeout, time.Second, service.DtsMigrateJobConfigStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}

func handleRecoverMigrate(d *schema.ResourceData, meta interface{}, logId, jobId string) error {
	request := dts.NewRecoverMigrateJobRequest()
	request.JobId = helper.String(jobId)
	// response = dts.NewPauseMigrateJobResponse()

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDtsClient().RecoverMigrateJob(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update dts migrateJobConfig failed, reason:%+v", logId, err)
		return err
	}

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"running"}, 2*readRetryTimeout, time.Second, service.DtsMigrateJobConfigStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}

func handleStopMigrate(d *schema.ResourceData, meta interface{}, logId, jobId string) error {
	request := dts.NewStopMigrateJobRequest()
	request.JobId = helper.String(jobId)
	// response = dts.NewPauseMigrateJobResponse()

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDtsClient().StopMigrateJob(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update dts migrateJobConfig failed, reason:%+v", logId, err)
		return err
	}

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"canceled"}, 2*readRetryTimeout, time.Second, service.DtsMigrateJobConfigStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}

func handleIsolateMigrate(d *schema.ResourceData, meta interface{}, logId, jobId string) error {
	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	err := service.IsolateDtsMigrateJobById(ctx, jobId)

	if err != nil {
		log.Printf("[CRITAL]%s isolate dts migrateJobConfig failed, reason:%+v", logId, err)
		return err
	}

	conf := BuildStateChangeConf([]string{}, []string{"isolated"}, 2*readRetryTimeout, time.Second, service.DtsMigrateJobConfigStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}
