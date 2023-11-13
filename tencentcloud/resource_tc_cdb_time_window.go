/*
Provides a resource to create a cdb time_window

Example Usage

```hcl
resource "tencentcloud_cdb_time_window" "time_window" {
  instance_id = "cdb-c1nl9rpv"
  time_ranges =
  weekdays =
  max_delay_time =
}
```

Import

cdb time_window can be imported using the id, e.g.

```
terraform import tencentcloud_cdb_time_window.time_window time_window_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"log"
)

func resourceTencentCloudCdbTimeWindow() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCdbTimeWindowCreate,
		Read:   resourceTencentCloudCdbTimeWindowRead,
		Update: resourceTencentCloudCdbTimeWindowUpdate,
		Delete: resourceTencentCloudCdbTimeWindowDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID in the format of cdb-c1nl9rpv or cdbro-c1nl9rpv. It is the same as the instance ID displayed on the TencentDB Console page.",
			},

			"time_ranges": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Time period available for maintenance after modification in the format of 10:00-12:00. Each period lasts from half an hour to three hours, with the start time and end time aligned by half-hour. Up to two time periods can be set. Start and end time range: [00:00, 24:00].",
			},

			"weekdays": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Specifies for which day to modify the time period. Value range: Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday. If it is not specified or is left blank, the time period will be modified for every day by default.",
			},

			"max_delay_time": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Data delay threshold. It takes effect only for source instance and disaster recovery instance. Default value: 10.",
			},
		},
	}
}

func resourceTencentCloudCdbTimeWindowCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_time_window.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudCdbTimeWindowUpdate(d, meta)
}

func resourceTencentCloudCdbTimeWindowRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_time_window.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	timeWindowId := d.Id()

	timeWindow, err := service.DescribeCdbTimeWindowById(ctx, instanceId)
	if err != nil {
		return err
	}

	if timeWindow == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CdbTimeWindow` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if timeWindow.InstanceId != nil {
		_ = d.Set("instance_id", timeWindow.InstanceId)
	}

	if timeWindow.TimeRanges != nil {
		_ = d.Set("time_ranges", timeWindow.TimeRanges)
	}

	if timeWindow.Weekdays != nil {
		_ = d.Set("weekdays", timeWindow.Weekdays)
	}

	if timeWindow.MaxDelayTime != nil {
		_ = d.Set("max_delay_time", timeWindow.MaxDelayTime)
	}

	return nil
}

func resourceTencentCloudCdbTimeWindowUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_time_window.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cdb.NewModifyTimeWindowRequest()

	timeWindowId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_id", "time_ranges", "weekdays", "max_delay_time"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCdbClient().ModifyTimeWindow(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cdb timeWindow failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCdbTimeWindowRead(d, meta)
}

func resourceTencentCloudCdbTimeWindowDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_time_window.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
