package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMpsEnableScheduleConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsEnableScheduleConfigCreate,
		Read:   resourceTencentCloudMpsEnableScheduleConfigRead,
		Update: resourceTencentCloudMpsEnableScheduleConfigUpdate,
		Delete: resourceTencentCloudMpsEnableScheduleConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"schedule_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "The scheme ID.",
			},

			"enabled": {
				Required:    true,
				Type:        schema.TypeBool,
				Description: "true: enable; false: disable.",
			},
		},
	}
}

func resourceTencentCloudMpsEnableScheduleConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_enable_schedule_config.create")()
	defer inconsistentCheck(d, meta)()

	var scheduleId int
	if v, ok := d.GetOkExists("schedule_id"); ok {
		scheduleId = v.(int)
	}
	d.SetId(helper.IntToStr(scheduleId))

	return resourceTencentCloudMpsEnableScheduleConfigUpdate(d, meta)
}

func resourceTencentCloudMpsEnableScheduleConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_enable_schedule_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}

	scheduleId := d.Id()

	schedules, err := service.DescribeMpsScheduleById(ctx, &scheduleId)
	if err != nil {
		return err
	}

	if len(schedules) == 0 {
		d.SetId("")
		log.Printf("[WARN]%s resource `MpsEnableScheduleConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	enableScheduleConfig := schedules[0]

	if enableScheduleConfig.ScheduleId != nil {
		_ = d.Set("schedule_id", enableScheduleConfig.ScheduleId)
	}

	status := enableScheduleConfig.Status
	if status != nil {
		if *status == "Enabled" {
			_ = d.Set("enabled", true)
		} else {
			_ = d.Set("enabled", false)
		}
	}

	return nil
}

func resourceTencentCloudMpsEnableScheduleConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_enable_schedule_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		enableRequest  = mps.NewEnableScheduleRequest()
		disableRequest = mps.NewDisableScheduleRequest()
		scheduleId     *int64
		enabled        bool
	)

	scheduleId = helper.StrToInt64Point(d.Id())

	if v, ok := d.GetOkExists("enabled"); ok && v != nil {
		enabled = v.(bool)

		if enabled {
			enableRequest.ScheduleId = scheduleId
			err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseMpsClient().EnableSchedule(enableRequest)
				if e != nil {
					return retryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, enableRequest.GetAction(), enableRequest.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s operate mps enableScheduleConfig failed, reason:%+v", logId, err)
				return err
			}
		} else {
			disableRequest.ScheduleId = scheduleId
			err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseMpsClient().DisableSchedule(disableRequest)
				if e != nil {
					return retryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, disableRequest.GetAction(), disableRequest.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s operate mps disableScheduleConfig failed, reason:%+v", logId, err)
				return err
			}
		}

	}
	return resourceTencentCloudMpsEnableScheduleConfigRead(d, meta)
}

func resourceTencentCloudMpsEnableScheduleConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_enable_schedule_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
