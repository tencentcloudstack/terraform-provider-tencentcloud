/*
Provides a resource to recover a clickhouse back up

Example Usage

```hcl
resource "tencentcloud_clickhouse_recover_backup_job" "recover_backup_job" {
  instance_id = "cdwch-xxxxxx"
  back_up_job_id = 1234
}
```
*/
package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clickhouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwch/v20200915"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudClickhouseRecoverBackupJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClickhouseRecoverBackupJobCreate,
		Read:   resourceTencentCloudClickhouseRecoverBackupJobRead,
		Delete: resourceTencentCloudClickhouseRecoverBackupJobDelete,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"back_up_job_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Back up job id.",
			},
		},
	}
}

func resourceTencentCloudClickhouseRecoverBackupJobCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clickhouse_recover_backup_job.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request     = clickhouse.NewRecoverBackUpJobRequest()
		instanceId  string
		backUpJobId int
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}

	if v, _ := d.GetOk("back_up_job_id"); v != nil {
		backUpJobId = v.(int)
		request.BackUpJobId = helper.IntInt64(backUpJobId)
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCdwchClient().RecoverBackUpJob(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate clickhouse RecoverBackupJob failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId + helper.IntToStr(backUpJobId))

	return resourceTencentCloudClickhouseRecoverBackupJobRead(d, meta)
}

func resourceTencentCloudClickhouseRecoverBackupJobRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clickhouse_recover_backup_job.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudClickhouseRecoverBackupJobDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clickhouse_recover_backup_job.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
