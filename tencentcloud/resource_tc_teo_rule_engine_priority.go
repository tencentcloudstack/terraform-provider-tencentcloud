/*
Provides a resource to create a teo rule_engine_priority

Example Usage

```hcl
resource "tencentcloud_teo_rule_engine_priority" "rule_engine_priority" {
  rules_priority = &lt;nil&gt;
}
```

Import

teo rule_engine_priority can be imported using the id, e.g.

```
terraform import tencentcloud_teo_rule_engine_priority.rule_engine_priority rule_engine_priority_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	"log"
)

func resourceTencentCloudTeoRuleEnginePriority() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoRuleEnginePriorityCreate,
		Read:   resourceTencentCloudTeoRuleEnginePriorityRead,
		Update: resourceTencentCloudTeoRuleEnginePriorityUpdate,
		Delete: resourceTencentCloudTeoRuleEnginePriorityDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"rules_priority": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Priority of rules.",
			},
		},
	}
}

func resourceTencentCloudTeoRuleEnginePriorityCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_rule_engine_priority.create")()
	defer inconsistentCheck(d, meta)()

	var zoneId string
	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	d.SetId(zoneId)

	return resourceTencentCloudTeoRuleEnginePriorityUpdate(d, meta)
}

func resourceTencentCloudTeoRuleEnginePriorityRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_rule_engine_priority.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}

	ruleEnginePriorityId := d.Id()

	ruleEnginePriority, err := service.DescribeTeoRuleEnginePriorityById(ctx, zoneId)
	if err != nil {
		return err
	}

	if ruleEnginePriority == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TeoRuleEnginePriority` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if ruleEnginePriority.RulesPriority != nil {
		_ = d.Set("rules_priority", ruleEnginePriority.RulesPriority)
	}

	return nil
}

func resourceTencentCloudTeoRuleEnginePriorityUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_rule_engine_priority.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := teo.NewModifyRulePriorityRequest()

	ruleEnginePriorityId := d.Id()

	request.ZoneId = &zoneId

	immutableArgs := []string{"rules_priority"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("rules_priority") {
		if v, ok := d.GetOk("rules_priority"); ok {
			rulesPrioritySet := v.(*schema.Set).List()
			for i := range rulesPrioritySet {
				rulesPriority := rulesPrioritySet[i].(string)
				request.RulesPriority = append(request.RulesPriority, &rulesPriority)
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyRulePriority(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update teo ruleEnginePriority failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTeoRuleEnginePriorityRead(d, meta)
}

func resourceTencentCloudTeoRuleEnginePriorityDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_rule_engine_priority.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
