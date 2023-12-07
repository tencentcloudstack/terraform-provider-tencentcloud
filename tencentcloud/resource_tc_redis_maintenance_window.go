package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudRedisMaintenanceWindow() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisMaintenanceWindowCreate,
		Read:   resourceTencentCloudRedisMaintenanceWindowRead,
		Update: resourceTencentCloudRedisMaintenanceWindowUpdate,
		Delete: resourceTencentCloudRedisMaintenanceWindowDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of instance.",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Maintenance window start time, e.g. 17:00.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The end time of the maintenance window, e.g. 19:00.",
			},
		},
	}
}

func resourceTencentCloudRedisMaintenanceWindowCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_maintenance_window.create")()
	defer inconsistentCheck(d, meta)()

	var (
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudRedisMaintenanceWindowUpdate(d, meta)
}

func resourceTencentCloudRedisMaintenanceWindowRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_maintenance_window.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	maintenanceWindow, err := service.DescribeRedisMaintenanceWindowById(ctx, instanceId)
	if err != nil {
		return err
	}

	if maintenanceWindow == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `RedisMaintenanceWindow` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if maintenanceWindow.StartTime != nil {
		_ = d.Set("start_time", maintenanceWindow.StartTime)
	}

	if maintenanceWindow.EndTime != nil {
		_ = d.Set("end_time", maintenanceWindow.EndTime)
	}

	return nil
}

func resourceTencentCloudRedisMaintenanceWindowUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_maintenance_window.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := redis.NewModifyMaintenanceWindowRequest()

	instanceId := d.Id()

	request.InstanceId = &instanceId

	if v, ok := d.GetOk("start_time"); ok {
		request.StartTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		request.EndTime = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseRedisClient().ModifyMaintenanceWindow(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update redis maintenanceWindow failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudRedisMaintenanceWindowRead(d, meta)
}

func resourceTencentCloudRedisMaintenanceWindowDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_maintenance_window.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
