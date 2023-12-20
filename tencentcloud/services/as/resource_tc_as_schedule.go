package as

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudAsSchedule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAsScheduleCreate,
		Read:   resourceTencentCloudAsScheduleRead,
		Update: resourceTencentCloudAsScheduleUpdate,
		Delete: resourceTencentCloudAsScheduleDelete,

		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of a scaling group.",
			},
			"schedule_action_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 60),
				Description:  "The name of this scaling action.",
			},
			"max_size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The maximum size for the Auto Scaling group.",
			},
			"min_size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The minimum size for the Auto Scaling group.",
			},
			"desired_capacity": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The desired number of CVM instances that should be running in the group.",
			},
			"start_time": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAsScheduleTimestamp,
				Description:  "The time for this action to start, in \"YYYY-MM-DDThh:mm:ss+08:00\" format (UTC+8).",
			},
			"end_time": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAsScheduleTimestamp,
				Description:  "The time for this action to end, in \"YYYY-MM-DDThh:mm:ss+08:00\" format (UTC+8).",
			},
			"recurrence": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The time when recurring future actions will start. Start time is specified by the user following the Unix cron syntax format. And this argument should be set with end_time together.",
			},
		},
	}
}

func resourceTencentCloudAsScheduleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_schedule.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := as.NewCreateScheduledActionRequest()
	request.AutoScalingGroupId = helper.String(d.Get("scaling_group_id").(string))
	request.ScheduledActionName = helper.String(d.Get("schedule_action_name").(string))
	request.MaxSize = helper.IntUint64(d.Get("max_size").(int))
	request.MinSize = helper.IntUint64(d.Get("min_size").(int))
	request.DesiredCapacity = helper.IntUint64(d.Get("desired_capacity").(int))
	request.StartTime = helper.String(d.Get("start_time").(string))

	// end_time and recurrence must be specified at the same time
	if v, ok := d.GetOk("end_time"); ok {
		request.EndTime = helper.String(v.(string))
		if vv, ok := d.GetOk("recurrence"); ok {
			request.Recurrence = helper.String(vv.(string))
		} else {
			return fmt.Errorf("end_time and recurrence must be specified at the same time.")
		}
	} else {
		if _, ok := d.GetOk("recurrence"); ok {
			return fmt.Errorf("end_time and recurrence must be specified at the same time.")
		}
	}

	response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAsClient().CreateScheduledAction(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.ScheduledActionId == nil {
		return fmt.Errorf("schedule action id is nil")
	}
	d.SetId(*response.Response.ScheduledActionId)

	return resourceTencentCloudAsScheduleRead(d, meta)
}

func resourceTencentCloudAsScheduleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_schedule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	scheduledActionId := d.Id()
	asService := AsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		scheduledAction, has, e := asService.DescribeScheduledActionById(ctx, scheduledActionId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if has == 0 {
			d.SetId("")
			return nil
		}

		_ = d.Set("scaling_group_id", *scheduledAction.AutoScalingGroupId)
		_ = d.Set("schedule_action_name", *scheduledAction.ScheduledActionName)
		_ = d.Set("max_size", *scheduledAction.MaxSize)
		_ = d.Set("min_size", *scheduledAction.MinSize)
		_ = d.Set("desired_capacity", *scheduledAction.DesiredCapacity)
		_ = d.Set("start_time", *scheduledAction.StartTime)

		if scheduledAction.EndTime != nil {
			_ = d.Set("end_time", *scheduledAction.EndTime)
		}
		if scheduledAction.Recurrence != nil {
			_ = d.Set("recurrence", *scheduledAction.Recurrence)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func resourceTencentCloudAsScheduleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_schedule.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := as.NewModifyScheduledActionRequest()
	scheduledActionId := d.Id()
	request.ScheduledActionId = &scheduledActionId
	if d.HasChange("schedule_action_name") {
		request.ScheduledActionName = helper.String(d.Get("schedule_action_name").(string))
	}
	if d.HasChange("max_size") {
		request.MaxSize = helper.IntUint64(d.Get("max_size").(int))
	}
	if d.HasChange("min_size") {
		request.MinSize = helper.IntUint64(d.Get("min_size").(int))
	}
	if d.HasChange("desired_capacity") {
		request.DesiredCapacity = helper.IntUint64(d.Get("desired_capacity").(int))
	}
	if d.HasChange("start_time") {
		request.StartTime = helper.String(d.Get("start_time").(string))
	}
	if d.HasChange("end_time") {
		request.EndTime = helper.String(d.Get("end_time").(string))
		request.Recurrence = helper.String(d.Get("recurrence").(string))
	}
	if d.HasChange("recurrence") {
		request.Recurrence = helper.String(d.Get("recurrence").(string))
		request.EndTime = helper.String(d.Get("end_time").(string))
	}

	response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAsClient().ModifyScheduledAction(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func resourceTencentCloudAsScheduleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_schedule.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	scheduledActionId := d.Id()
	asService := AsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	err := asService.DeleteScheduledAction(ctx, scheduledActionId)
	if err != nil {
		return err
	}

	return nil
}
