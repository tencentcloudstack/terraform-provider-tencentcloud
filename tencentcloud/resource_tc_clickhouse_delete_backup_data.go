package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clickhouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwch/v20200915"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudClickhouseDeleteBackupData() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClickhouseDeleteBackupDataCreate,
		Read:   resourceTencentCloudClickhouseDeleteBackupDataRead,
		Delete: resourceTencentCloudClickhouseDeleteBackupDataDelete,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"back_up_job_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Back up job id.",
			},
		},
	}
}

func resourceTencentCloudClickhouseDeleteBackupDataCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clickhouse_delete_backup_data.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request     = clickhouse.NewDeleteBackUpDataRequest()
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
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCdwchClient().DeleteBackUpData(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate clickhouse deleteBackUpData failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId + helper.IntToStr(backUpJobId))

	return resourceTencentCloudClickhouseDeleteBackupDataRead(d, meta)
}

func resourceTencentCloudClickhouseDeleteBackupDataRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clickhouse_delete_backup_data.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudClickhouseDeleteBackupDataDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clickhouse_delete_backup_data.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
