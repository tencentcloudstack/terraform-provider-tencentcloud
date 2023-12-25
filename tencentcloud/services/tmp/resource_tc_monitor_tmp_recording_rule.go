package tmp

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcmonitor "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/monitor"

	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMonitorTmpRecordingRule() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudMonitorTmpRecordingRuleRead,
		Create: resourceTencentCloudMonitorTmpRecordingRuleCreate,
		Update: resourceTencentCloudMonitorTmpRecordingRuleUpdate,
		Delete: resourceTencentCloudMonitorTmpRecordingRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Recording rule name.",
			},
			"group": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Recording rule group.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance id.",
			},
			"rule_state": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Rule state.",
			},
		},
	}
}

func resourceTencentCloudMonitorTmpRecordingRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_recording_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = monitor.NewCreateRecordingRuleRequest()
		response *monitor.CreateRecordingRuleResponse
	)

	var instanceId string

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}
	if v, ok := d.GetOk("group"); ok {
		request.Group = helper.String(tccommon.StringToBase64(v.(string)))
	}
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}
	if v, ok := d.GetOk("rule_state"); ok {
		request.RuleState = helper.IntInt64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorClient().CreateRecordingRule(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create monitor recordingRule failed, reason:%+v", logId, err)
		return err
	}

	recordingRuleId := *response.Response.RuleId
	d.SetId(strings.Join([]string{instanceId, recordingRuleId}, tccommon.FILED_SP))

	return resourceTencentCloudMonitorTmpRecordingRuleRead(d, meta)
}

func resourceTencentCloudMonitorTmpRecordingRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tmp_monitor_recording_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	ids := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(ids) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	recordingRule, err := service.DescribeMonitorRecordingRuleById(ctx, ids[0], ids[1])

	if err != nil {
		return err
	}

	if recordingRule == nil {
		d.SetId("")
		return fmt.Errorf("resource `recordingRule` %s does not exist", ids[1])
	}

	_ = d.Set("instance_id", ids[0])
	if recordingRule.Name != nil {
		_ = d.Set("name", recordingRule.Name)
	}
	if recordingRule.Group != nil {
		group, err := tccommon.Base64ToString(*recordingRule.Group)
		if err != nil {
			return fmt.Errorf("`recordingRule.Group` %s does not be decoded to yaml", *recordingRule.Group)
		}
		_ = d.Set("group", &group)
	}

	if recordingRule.RuleState != nil {
		_ = d.Set("rule_state", recordingRule.RuleState)
	}

	return nil
}

func resourceTencentCloudMonitorTmpRecordingRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tmp_monitor_recording_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := monitor.NewUpdateRecordingRuleRequest()

	ids := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(ids) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	request.InstanceId = helper.String(ids[0])
	request.RuleId = helper.String(ids[1])

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group"); ok {
		request.Group = helper.String(v.(string))
	}

	if d.HasChange("rule_state") {
		if v, ok := d.GetOk("rule_state"); ok {
			request.RuleState = helper.IntInt64(v.(int))
		}
	}
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorClient().UpdateRecordingRule(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudMonitorTmpRecordingRuleRead(d, meta)
}

func resourceTencentCloudMonitorTmpRecordingRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_recording_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	ids := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(ids) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	service := svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	if err := service.DeleteMonitorRecordingRule(ctx, ids[0], ids[1]); err != nil {
		return err
	}

	return nil
}
