/*
Provides a resource to create a eb event_rule

Example Usage

```hcl
resource "tencentcloud_eb_event_rule" "event_rule" {
  event_pattern = ""
  event_bus_id = ""
  rule_name = ""
  enable =
  description = ""
}
```

Import

eb event_rule can be imported using the id, e.g.

```
terraform import tencentcloud_eb_event_rule.event_rule event_rule_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	eb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/eb/v20210416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudEbEventRule() *schema.Resource {
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
				Description: "Event bus Id.",
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
		},
	}
}

func resourceTencentCloudEbEventRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eb_event_rule.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = eb.NewCreateRuleRequest()
		response   = eb.NewCreateRuleResponse()
		ruleId     string
		eventBusId string
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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseEbClient().CreateRule(request)
		if e != nil {
			return retryError(e)
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
	d.SetId(strings.Join([]string{ruleId, eventBusId}, FILED_SP))

	return resourceTencentCloudEbEventRuleRead(d, meta)
}

func resourceTencentCloudEbEventRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eb_event_rule.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := EbService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	ruleId := idSplit[0]
	eventBusId := idSplit[1]

	eventRule, err := service.DescribeEbEventRuleById(ctx, ruleId, eventBusId)
	if err != nil {
		return err
	}

	if eventRule == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `EbEventRule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

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

	return nil
}

func resourceTencentCloudEbEventRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eb_event_rule.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := eb.NewUpdateRuleRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	ruleId := idSplit[0]
	eventBusId := idSplit[1]

	request.RuleId = &ruleId
	request.EventBusId = &eventBusId

	immutableArgs := []string{"event_pattern", "event_bus_id", "rule_name", "enable", "description"}

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

	if d.HasChange("event_bus_id") {
		if v, ok := d.GetOk("event_bus_id"); ok {
			request.EventBusId = helper.String(v.(string))
		}
	}

	if d.HasChange("rule_name") {
		if v, ok := d.GetOk("rule_name"); ok {
			request.RuleName = helper.String(v.(string))
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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseEbClient().UpdateRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update eb eventRule failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudEbEventRuleRead(d, meta)
}

func resourceTencentCloudEbEventRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eb_event_rule.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := EbService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	ruleId := idSplit[0]
	eventBusId := idSplit[1]

	if err := service.DeleteEbEventRuleById(ctx, ruleId, eventBusId); err != nil {
		return err
	}

	return nil
}
