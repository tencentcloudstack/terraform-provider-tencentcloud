package eb

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	eb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/eb/v20210416"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudEbEventRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudEbEventRuleCreate,
		Read:   resourceTencentCloudEbEventRuleRead,
		Update: resourceTencentCloudEbEventRuleUpdate,
		Delete: resourceTencentCloudEbEventRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"event_pattern": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Reference: [Event Mode](https://cloud.tencent.com/document/product/1359/56084).",
			},

			"event_bus_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "event bus Id.",
			},

			"rule_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Event rule name, which can only contain letters, numbers, underscores, hyphens, starts with a letter and ends with a number or letter, 2~60 characters.",
			},

			"enable": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Enable switch.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Event set description, unlimited character type, description within 200 characters.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},

			"rule_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "event rule id.",
			},
		},
	}
}

func resourceTencentCloudEbEventRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eb_event_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = eb.NewCreateRuleRequest()
		response   = eb.NewCreateRuleResponse()
		eventBusId string
		ruleId     string
	)
	if v, ok := d.GetOk("event_pattern"); ok {
		request.EventPattern = helper.String(v.(string))
	}

	if v, ok := d.GetOk("event_bus_id"); ok {
		eventBusId = v.(string)
		request.EventBusId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rule_name"); ok {
		request.RuleName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("enable"); ok {
		request.Enable = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEbClient().CreateRule(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create eb eventRule failed, reason:%+v", logId, err)
		return err
	}

	ruleId = *response.Response.RuleId
	d.SetId(eventBusId + tccommon.FILED_SP + ruleId)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := fmt.Sprintf("qcs::eb:%s:uin/:ruleid/%s/%s", region, eventBusId, ruleId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudEbEventRuleRead(d, meta)
}

func resourceTencentCloudEbEventRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eb_event_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := EbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	eventBusId := idSplit[0]
	ruleId := idSplit[1]

	eventRule, err := service.DescribeEbEventRuleById(ctx, eventBusId, ruleId)
	if err != nil {
		return err
	}

	if eventRule == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `EbEventRule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("rule_id", ruleId)

	if eventRule.EventPattern != nil {
		_ = d.Set("event_pattern", eventRule.EventPattern)
	}

	if eventRule.EventBusId != nil {
		_ = d.Set("event_bus_id", eventRule.EventBusId)
	}

	if eventRule.RuleName != nil {
		_ = d.Set("rule_name", eventRule.RuleName)
	}

	if eventRule.Enable != nil {
		_ = d.Set("enable", eventRule.Enable)
	}

	if eventRule.Description != nil {
		_ = d.Set("description", eventRule.Description)
	}

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
	tags, err := tagService.DescribeResourceTags(ctx, "eb", "ruleid", tcClient.Region, eventBusId+"/"+ruleId)
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudEbEventRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eb_event_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := eb.NewUpdateRuleRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	eventBusId := idSplit[0]
	ruleId := idSplit[1]

	request.EventBusId = &eventBusId
	request.RuleId = &ruleId

	immutableArgs := []string{"event_bus_id", "rule_name"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("event_pattern") {
		if v, ok := d.GetOk("event_pattern"); ok {
			request.EventPattern = helper.String(v.(string))
		}
	}

	if d.HasChange("enable") {
		if v, ok := d.GetOkExists("enable"); ok {
			request.Enable = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEbClient().UpdateRule(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update eb eventRule failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := tccommon.BuildTagResourceName("eb", "ruleid", tcClient.Region, eventBusId+"/"+ruleId)
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudEbEventRuleRead(d, meta)
}

func resourceTencentCloudEbEventRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eb_event_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := EbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	eventBusId := idSplit[0]
	ruleId := idSplit[1]

	if err := service.DeleteEbEventRuleById(ctx, eventBusId, ruleId); err != nil {
		return err
	}

	return nil
}
