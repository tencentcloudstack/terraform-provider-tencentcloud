package config

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	configv20220802 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudConfigAlarmPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudConfigAlarmPolicyCreate,
		Read:   resourceTencentCloudConfigAlarmPolicyRead,
		Update: resourceTencentCloudConfigAlarmPolicyUpdate,
		Delete: resourceTencentCloudConfigAlarmPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Alarm policy name.",
			},

			"event_scope": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Event scope. Valid values: 1 (current account), 2 (multi-account).",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"risk_level": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Risk level list. Valid values: 1 (high risk), 2 (medium risk), 3 (low risk).",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"notice_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Notification time range, format: HH:mm:ss~HH:mm:ss (e.g. 09:30:00~23:30:00).",
			},

			"notification_mechanism": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Notification mechanism.",
			},

			"status": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Policy status. Valid values: 1 (enabled), 2 (disabled).",
			},

			"notice_period": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Notification weekday list, 1-7 represent Monday to Sunday.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Policy description.",
			},

			// Computed
			"alarm_policy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Alarm policy unique ID.",
			},
		},
	}
}

func resourceTencentCloudConfigAlarmPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_config_alarm_policy.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = configv20220802.NewAddAlarmPolicyRequest()
		response = configv20220802.NewAddAlarmPolicyResponse()
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("event_scope"); ok {
		rawList := v.([]interface{})
		for _, item := range rawList {
			request.EventScope = append(request.EventScope, helper.IntInt64(item.(int)))
		}
	}

	if v, ok := d.GetOk("risk_level"); ok {
		rawList := v.([]interface{})
		for _, item := range rawList {
			request.RiskLevel = append(request.RiskLevel, helper.IntInt64(item.(int)))
		}
	}

	if v, ok := d.GetOk("notice_time"); ok {
		request.NoticeTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("notification_mechanism"); ok {
		request.NotificationMechanism = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("status"); ok {
		request.Status = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("notice_period"); ok {
		rawList := v.([]interface{})
		for _, item := range rawList {
			request.NoticePeriod = append(request.NoticePeriod, helper.IntInt64(item.(int)))
		}
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseConfigV20220802Client().AddAlarmPolicyWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("create config alarm policy failed, Response is nil"))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create config alarm policy failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.AlarmPolicyId == nil {
		return fmt.Errorf("AlarmPolicyId is nil")
	}

	alarmPolicyId := *response.Response.AlarmPolicyId
	d.SetId(helper.UInt64ToStr(alarmPolicyId))
	return resourceTencentCloudConfigAlarmPolicyRead(d, meta)
}

func resourceTencentCloudConfigAlarmPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_config_alarm_policy.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = ConfigService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idUint, err := strconv.ParseUint(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("invalid alarm policy id: %s", d.Id())
	}

	respData, err := service.DescribeConfigAlarmPolicyById(ctx, idUint)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_config_alarm_policy` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.Name != nil {
		_ = d.Set("name", respData.Name)
	}

	if respData.EventScope != nil {
		eventScope := make([]int, 0, len(respData.EventScope))
		for _, v := range respData.EventScope {
			if v != nil {
				eventScope = append(eventScope, int(*v))
			}
		}

		_ = d.Set("event_scope", eventScope)
	}

	if respData.RiskLevel != nil {
		riskLevel := make([]int, 0, len(respData.RiskLevel))
		for _, v := range respData.RiskLevel {
			if v != nil {
				riskLevel = append(riskLevel, int(*v))
			}
		}

		_ = d.Set("risk_level", riskLevel)
	}

	if respData.NoticeTime != nil {
		_ = d.Set("notice_time", respData.NoticeTime)
	}

	if respData.NotificationMechanism != nil {
		_ = d.Set("notification_mechanism", respData.NotificationMechanism)
	}

	if respData.Status != nil {
		_ = d.Set("status", int(*respData.Status))
	}

	if respData.NoticePeriod != nil {
		noticePeriod := make([]int, 0, len(respData.NoticePeriod))
		for _, v := range respData.NoticePeriod {
			if v != nil {
				noticePeriod = append(noticePeriod, int(*v))
			}
		}

		_ = d.Set("notice_period", noticePeriod)
	}

	if respData.Description != nil {
		_ = d.Set("description", respData.Description)
	}

	if respData.AlarmPolicyId != nil {
		_ = d.Set("alarm_policy_id", helper.UInt64ToStr(*respData.AlarmPolicyId))
	}

	return nil
}

func resourceTencentCloudConfigAlarmPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_config_alarm_policy.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idUint, err := strconv.ParseUint(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("invalid alarm policy id: %s", d.Id())
	}

	mutableArgs := []string{"name", "event_scope", "risk_level", "notice_time", "notification_mechanism", "status", "notice_period", "description"}
	needChange := false
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := configv20220802.NewUpdateAlarmPolicyRequest()
		request.AlarmPolicyId = &idUint

		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}

		if v, ok := d.GetOk("event_scope"); ok {
			rawList := v.([]interface{})
			for _, item := range rawList {
				request.EventScope = append(request.EventScope, helper.IntInt64(item.(int)))
			}
		}

		if v, ok := d.GetOk("risk_level"); ok {
			rawList := v.([]interface{})
			for _, item := range rawList {
				request.RiskLevel = append(request.RiskLevel, helper.IntInt64(item.(int)))
			}
		}

		if v, ok := d.GetOk("notice_time"); ok {
			request.NoticeTime = helper.String(v.(string))
		}

		if v, ok := d.GetOk("notification_mechanism"); ok {
			request.NotificationMechanism = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("status"); ok {
			request.Status = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOk("notice_period"); ok {
			rawList := v.([]interface{})
			for _, item := range rawList {
				request.NoticePeriod = append(request.NoticePeriod, helper.IntInt64(item.(int)))
			}
		}

		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseConfigV20220802Client().UpdateAlarmPolicyWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update config alarm policy failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudConfigAlarmPolicyRead(d, meta)
}

func resourceTencentCloudConfigAlarmPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_config_alarm_policy.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = configv20220802.NewDeleteAlarmPolicyRequest()
	)

	idUint, err := strconv.ParseUint(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("invalid alarm policy id: %s", d.Id())
	}

	request.AlarmPolicyId = &idUint

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseConfigV20220802Client().DeleteAlarmPolicyWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete config alarm policy failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
