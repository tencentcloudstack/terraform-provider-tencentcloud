/*
Provides a resource to create a mps schedule_config

Example Usage

```hcl
resource "tencentcloud_mps_schedule_config" "schedule_config" {
  schedule_id =
}
```

Import

mps schedule_config can be imported using the id, e.g.

```
terraform import tencentcloud_mps_schedule_config.schedule_config schedule_config_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudMpsScheduleConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsScheduleConfigCreate,
		Read:   resourceTencentCloudMpsScheduleConfigRead,
		Update: resourceTencentCloudMpsScheduleConfigUpdate,
		Delete: resourceTencentCloudMpsScheduleConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"schedule_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "The scheme ID.",
			},
		},
	}
}

func resourceTencentCloudMpsScheduleConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_schedule_config.create")()
	defer inconsistentCheck(d, meta)()

	var scheduleId int64
	if v, ok := d.GetOkExists("schedule_id"); ok {
		scheduleId = v.(int64)
	}

	d.SetId(helper.Int64ToStr(scheduleId))

	return resourceTencentCloudMpsScheduleConfigUpdate(d, meta)
}

func resourceTencentCloudMpsScheduleConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_schedule_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}

	scheduleConfigId := d.Id()

	scheduleConfig, err := service.DescribeMpsScheduleConfigById(ctx, scheduleId)
	if err != nil {
		return err
	}

	if scheduleConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MpsScheduleConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if scheduleConfig.ScheduleId != nil {
		_ = d.Set("schedule_id", scheduleConfig.ScheduleId)
	}

	return nil
}

func resourceTencentCloudMpsScheduleConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_schedule_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := mps.NewEnableScheduleRequest()

	scheduleConfigId := d.Id()

	request.ScheduleId = &scheduleId

	immutableArgs := []string{"schedule_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMpsClient().EnableSchedule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update mps scheduleConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMpsScheduleConfigRead(d, meta)
}

func resourceTencentCloudMpsScheduleConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_schedule_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
