/*
Provides a resource to create a mysql time_window

Example Usage

```hcl
resource "tencentcloud_mysql_time_window" "time_window" {
  instance_id    = "cdb-lw71b6ar"
  max_delay_time = 10
  time_ranges    = [
    "01:00-02:01"
  ]
  weekdays       = [
    "friday",
    "monday",
    "saturday",
    "thursday",
    "tuesday",
    "wednesday",
  ]
}
```

Import

mysql time_window can be imported using the id, e.g.

```
terraform import tencentcloud_mysql_time_window.time_window instanceId
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMysqlTimeWindow() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlTimeWindowCreate,
		Read:   resourceTencentCloudMysqlTimeWindowRead,
		Update: resourceTencentCloudMysqlTimeWindowUpdate,
		Delete: resourceTencentCloudMysqlTimeWindowDelete,
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

func resourceTencentCloudMysqlTimeWindowCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_time_window.create")()
	defer inconsistentCheck(d, meta)()

	d.SetId(d.Get("instance_id").(string))

	return resourceTencentCloudMysqlTimeWindowUpdate(d, meta)
}

func resourceTencentCloudMysqlTimeWindowRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_time_window.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	timeWindow, err := service.DescribeMysqlTimeWindowById(ctx, instanceId)
	if err != nil {
		return err
	}

	if timeWindow == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `tencentcloud_mysql_time_window` [%s] not found, please check if it has been deleted.",
			logId, instanceId,
		)
		return nil
	}

	var timeRanges []*string
	var weekdays []*string

	_ = d.Set("instance_id", instanceId)

	if *timeWindow.Response.Monday[0] != "00:00-00:00" {
		timeRanges = timeWindow.Response.Monday
		weekdays = append(weekdays, helper.String("monday"))
	}
	if *timeWindow.Response.Tuesday[0] != "00:00-00:00" {
		timeRanges = timeWindow.Response.Tuesday
		weekdays = append(weekdays, helper.String("tuesday"))
	}
	if *timeWindow.Response.Wednesday[0] != "00:00-00:00" {
		timeRanges = timeWindow.Response.Wednesday
		weekdays = append(weekdays, helper.String("wednesday"))
	}
	if *timeWindow.Response.Thursday[0] != "00:00-00:00" {
		timeRanges = timeWindow.Response.Thursday
		weekdays = append(weekdays, helper.String("thursday"))
	}
	if *timeWindow.Response.Friday[0] != "00:00-00:00" {
		timeRanges = timeWindow.Response.Friday
		weekdays = append(weekdays, helper.String("friday"))
	}
	if *timeWindow.Response.Saturday[0] != "00:00-00:00" {
		timeRanges = timeWindow.Response.Saturday
		weekdays = append(weekdays, helper.String("saturday"))
	}
	if *timeWindow.Response.Wednesday[0] != "00:00-00:00" {
		timeRanges = timeWindow.Response.Wednesday
		weekdays = append(weekdays, helper.String("wednesday"))
	}

	if timeRanges != nil {
		_ = d.Set("time_ranges", timeRanges)
	}

	if weekdays != nil {
		_ = d.Set("weekdays", weekdays)
	}

	if timeWindow.Response.MaxDelayTime != nil {
		_ = d.Set("max_delay_time", timeWindow.Response.MaxDelayTime)
	}

	return nil
}

func resourceTencentCloudMysqlTimeWindowUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_time_window.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := mysql.NewModifyTimeWindowRequest()

	instanceId := d.Id()

	request.InstanceId = &instanceId

	if v, ok := d.GetOk("time_ranges"); ok {
		timeRangesSet := v.(*schema.Set).List()
		for i := range timeRangesSet {
			timeRange := timeRangesSet[i].(string)
			request.TimeRanges = append(request.TimeRanges, &timeRange)
		}
	}

	if v, ok := d.GetOk("weekdays"); ok {
		weekdaysSet := v.(*schema.Set).List()
		for i := range weekdaysSet {
			weekday := weekdaysSet[i].(string)
			request.Weekdays = append(request.Weekdays, &weekday)
		}
	}

	if d.HasChange("max_delay_time") {
		if v, _ := d.GetOk("max_delay_time"); v != nil {
			request.MaxDelayTime = helper.IntUint64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMysqlClient().ModifyTimeWindow(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		d.SetId("")
		log.Printf("[CRITAL]%s update mysql timeWindow failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMysqlTimeWindowRead(d, meta)
}

func resourceTencentCloudMysqlTimeWindowDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_time_window.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	instanceId := d.Id()

	if err := service.DeleteMysqlTimeWindowById(ctx, instanceId); err != nil {
		return err
	}

	return nil
}
