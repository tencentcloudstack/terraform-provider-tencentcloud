package cvm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCvmActionTimer() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCvmActionTimerCreate,
		Read:   resourceTencentCloudCvmActionTimerRead,
		//Update: resourceTencentCloudCvmActionTimerUpdate,
		Delete: resourceTencentCloudCvmActionTimerDelete,
		//Importer: &schema.ResourceImporter{
		//	State: schema.ImportStatePassthrough,
		//},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
			"action_timer": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Scheduled tasks. This parameter can be used to specify scheduled tasks for instances, and currently only supports scheduled destruction.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timer_action": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Timer action, currently only supports destroying one value: TerminateInstances.",
						},
						"action_time": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Execution time, expressed according to ISO8601 standard and using UTC time. The format is YYYY-MM-DDThh:mm:ssZ. For example, 2018-05-29T11:26:40Z, the execution time must be 5 minutes longer than the current time.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCvmActionTimerCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_action_timer.create")()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		request       = cvm.NewImportInstancesActionTimerRequest()
		response      = cvm.NewImportInstancesActionTimerResponse()
		actionTimerId string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceIds = []*string{helper.String(v.(string))}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "action_timer"); ok {
		actionTimer := cvm.ActionTimer{}
		if v, ok := dMap["timer_action"]; ok {
			actionTimer.TimerAction = helper.String(v.(string))
		}

		if v, ok := dMap["action_time"]; ok {
			actionTimer.ActionTime = helper.String(v.(string))
		}

		request.ActionTimer = &actionTimer
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().ImportInstancesActionTimer(request)
		if e != nil {
			return resource.RetryableError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || len(result.Response.ActionTimerIds) != 1 {
			e = fmt.Errorf("create cvm InstanceActionTimer failed")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cvm InstanceActionTimer failed, reason:%+v", logId, err)
		return err
	}

	actionTimerId = *response.Response.ActionTimerIds[0]
	d.SetId(actionTimerId)

	return resourceTencentCloudCvmActionTimerRead(d, meta)
}

func resourceTencentCloudCvmActionTimerRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_action_timer.read")()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service       = CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		actionTimerId = d.Id()
	)

	InstanceActionTimer, err := service.DescribeCvmInstanceActionTimerById(ctx, actionTimerId)
	if err != nil {
		return err
	}

	if InstanceActionTimer == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CvmInstanceActionTimer` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	actionTimerMap := map[string]interface{}{}
	if InstanceActionTimer.TimerAction != nil {
		actionTimerMap["timer_action"] = InstanceActionTimer.TimerAction
	}

	if InstanceActionTimer.ActionTime != nil {
		actionTimerMap["action_time"] = InstanceActionTimer.ActionTime
	}

	_ = d.Set("action_timer", []interface{}{actionTimerMap})

	return nil
}

// func resourceTencentCloudCvmActionTimerUpdate(d *schema.ResourceData, meta interface{}) error {
// 	defer tccommon.LogElapsed("resource.tencentcloud_cvm_action_timer.update")()
// 	defer tccommon.InconsistentCheck(d, meta)()

// 	return resourceTencentCloudCvmActionTimerRead(d, meta)
// }

func resourceTencentCloudCvmActionTimerDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_action_timer.delete")()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service       = CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		actionTimerId = d.Id()
	)

	if err := service.DeleteCvmInstanceActionTimerById(ctx, actionTimerId); err != nil {
		return err
	}

	return nil
}
