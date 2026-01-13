package mqtt

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mqttv20240516 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mqtt/v20240516"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMqttMessageEnrichmentRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMqttMessageEnrichmentRuleCreate,
		Read:   resourceTencentCloudMqttMessageEnrichmentRuleRead,
		Update: resourceTencentCloudMqttMessageEnrichmentRuleUpdate,
		Delete: resourceTencentCloudMqttMessageEnrichmentRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "MQTT instance ID.",
			},

			"rule_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Rule name.",
			},

			"condition": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Condition expression, Base64 encoded JSON string.",
			},

			"actions": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Actions, Base64 encoded JSON string.",
			},

			"status": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Rule status: 1 for enabled, 0 for disabled. Default: 1.",
			},

			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Remark.",
			},

			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Priority for rule execution order. Lower values have higher priority. Default: 1.",
			},

			// computed fields
			"rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule ID assigned by the cloud service.",
			},

			"created_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Creation time (Unix timestamp in seconds).",
			},

			"update_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Update time (Unix timestamp in seconds).",
			},
		},
	}
}

func resourceTencentCloudMqttMessageEnrichmentRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mqtt_message_enrichment_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = mqttv20240516.NewCreateMessageEnrichmentRuleRequest()
		response   = mqttv20240516.NewCreateMessageEnrichmentRuleResponse()
		instanceId string
		ruleId     int64
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("rule_name"); ok {
		request.RuleName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("condition"); ok {
		request.Condition = helper.String(v.(string))
	}

	if v, ok := d.GetOk("actions"); ok {
		request.Actions = helper.String(v.(string))
	}

	if v, ok := d.GetOk("status"); ok {
		request.Status = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	if v, ok := d.GetOk("priority"); ok {
		request.Priority = helper.IntInt64(v.(int))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMqttV20240516Client().CreateMessageEnrichmentRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.RetryableError(fmt.Errorf("Create mqtt message enrichment rule failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create mqtt message enrichment rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.Id == nil {
		return fmt.Errorf("Id is nil.")
	}

	ruleId = *response.Response.Id
	d.SetId(strings.Join([]string{instanceId, helper.Int64ToStr(ruleId)}, tccommon.FILED_SP))

	return resourceTencentCloudMqttMessageEnrichmentRuleRead(d, meta)
}

func resourceTencentCloudMqttMessageEnrichmentRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mqtt_message_enrichment_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = MqttService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	ruleId := idSplit[1]

	ruleIdInt64 := helper.StrToInt64(ruleId)

	rule, err := service.DescribeMqttMessageEnrichmentRuleById(ctx, instanceId, ruleIdInt64)
	if err != nil {
		return err
	}

	if rule == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `mqtt_message_enrichment_rule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)
	_ = d.Set("rule_id", ruleId)

	if rule.RuleName != nil {
		_ = d.Set("rule_name", rule.RuleName)
	}

	if rule.Condition != nil {
		_ = d.Set("condition", rule.Condition)
	}

	if rule.Actions != nil {
		_ = d.Set("actions", rule.Actions)
	}

	if rule.Status != nil {
		_ = d.Set("status", rule.Status)
	}

	if rule.Remark != nil {
		_ = d.Set("remark", rule.Remark)
	}

	if rule.Priority != nil {
		_ = d.Set("priority", rule.Priority)
	}

	if rule.CreatedTime != nil {
		_ = d.Set("created_time", rule.CreatedTime)
	}

	if rule.UpdateTime != nil {
		_ = d.Set("update_time", rule.UpdateTime)
	}

	return nil
}

func resourceTencentCloudMqttMessageEnrichmentRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mqtt_message_enrichment_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = mqttv20240516.NewModifyMessageEnrichmentRuleRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	ruleId := idSplit[1]

	request.InstanceId = &instanceId
	request.Id = helper.StrToInt64Point(ruleId)

	// According to API spec, all fields must be submitted for update (full update semantics)
	if v, ok := d.GetOk("rule_name"); ok {
		request.RuleName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("condition"); ok {
		request.Condition = helper.String(v.(string))
	}

	if v, ok := d.GetOk("actions"); ok {
		request.Actions = helper.String(v.(string))
	}

	if v, ok := d.GetOk("status"); ok {
		request.Status = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	if v, ok := d.GetOk("priority"); ok {
		request.Priority = helper.IntInt64(v.(int))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMqttV20240516Client().ModifyMessageEnrichmentRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update mqtt message enrichment rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudMqttMessageEnrichmentRuleRead(d, meta)
}

func resourceTencentCloudMqttMessageEnrichmentRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mqtt_message_enrichment_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = MqttService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	ruleId := idSplit[1]

	ruleIdInt64 := helper.StrToInt64(ruleId)

	if err := service.DeleteMqttMessageEnrichmentRuleById(ctx, instanceId, ruleIdInt64); err != nil {
		log.Printf("[CRITAL]%s delete mqtt message enrichment rule failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
