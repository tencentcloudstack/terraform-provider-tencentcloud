package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMariadbLogFileRetentionPeriod() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudMariadbLogFileRetentionPeriodRead,
		Create: resourceTencentCloudMariadbLogFileRetentionPeriodCreate,
		Update: resourceTencentCloudMariadbLogFileRetentionPeriodUpdate,
		Delete: resourceTencentCloudMariadbLogFileRetentionPeriodDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "instance id.",
			},

			"days": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The number of days to save, cannot exceed 30.",
			},
		},
	}
}

func resourceTencentCloudMariadbLogFileRetentionPeriodCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_log_file_retention_period.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)
	return resourceTencentCloudMariadbLogFileRetentionPeriodUpdate(d, meta)
}

func resourceTencentCloudMariadbLogFileRetentionPeriodRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_log_file_retention_period.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	logFileRetentionPeriod, err := service.DescribeMariadbLogFileRetentionPeriod(ctx, instanceId)

	if err != nil {
		return err
	}

	if logFileRetentionPeriod == nil {
		d.SetId("")
		return fmt.Errorf("resource `logFileRetentionPeriod` %s does not exist", instanceId)
	}

	if logFileRetentionPeriod.InstanceId != nil {
		_ = d.Set("instance_id", logFileRetentionPeriod.InstanceId)
	}

	if logFileRetentionPeriod.Days != nil {
		_ = d.Set("days", int(*logFileRetentionPeriod.Days))
	}

	return nil
}

func resourceTencentCloudMariadbLogFileRetentionPeriodUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_log_file_retention_period.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := mariadb.NewModifyLogFileRetentionPeriodRequest()

	instanceId := d.Id()

	request.InstanceId = &instanceId

	if v, ok := d.GetOk("days"); ok {
		request.Days = helper.Uint64(uint64(v.(int)))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().ModifyLogFileRetentionPeriod(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create mariadb logFileRetentionPeriod failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMariadbLogFileRetentionPeriodRead(d, meta)
}

func resourceTencentCloudMariadbLogFileRetentionPeriodDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_log_file_retention_period.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
