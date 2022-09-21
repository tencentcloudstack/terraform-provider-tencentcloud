/*
Provides a resource to create a teo rule_engine_priority

Example Usage

```hcl
resource "tencentcloud_teo_rule_engine_priority" "rule_engine_priority" {
  zone_id = "zone-294v965lwmn6"

  rules_priority {
    index = 0
    value = "rule-m9jlttua"
  }
  rules_priority {
    index = 1
    value = "rule-m5l9t4k1"
  }
}

```
Import

teo rule_engine_priority can be imported using the zone_id, e.g.
```
$ terraform import tencentcloud_teo_rule_engine_priority.rule_engine_priority zone-294v965lwmn6
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
)

func resourceTencentCloudTeoRuleEnginePriority() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTeoRuleEnginePriorityRead,
		Create: resourceTencentCloudTeoRuleEnginePriorityCreate,
		Update: resourceTencentCloudTeoRuleEnginePriorityUpdate,
		Delete: resourceTencentCloudTeoRuleEnginePriorityDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Site ID.",
			},

			"rules_priority": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Priority of rules.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"index": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Priority order of rules.",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Priority of rules id.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTeoRuleEnginePriorityCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_rule_engine_priority.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var zoneId string
	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	d.SetId(zoneId)
	err := resourceTencentCloudTeoRuleEnginePriorityUpdate(d, meta)
	if err != nil {
		log.Printf("[CRITAL]%s create teo ruleEnginePriority failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTeoRuleEnginePriorityRead(d, meta)
}

func resourceTencentCloudTeoRuleEnginePriorityRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_rule_engine_priority.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}

	zoneId := d.Id()

	ruleEnginePriority, err := service.DescribeTeoRuleEnginePriority(ctx, zoneId)

	if err != nil {
		return err
	}

	if ruleEnginePriority == nil {
		d.SetId("")
		return fmt.Errorf("resource `ruleEnginePriority` %s does not exist", zoneId)
	}

	_ = d.Set("zone_id", zoneId)

	if ruleEnginePriority != nil {
		ruleEnginePriorityList := []interface{}{}
		for i, v := range ruleEnginePriority {
			ruleId := map[string]interface{}{}
			ruleId["index"] = i
			ruleId["value"] = v.RuleId
			ruleEnginePriorityList = append(ruleEnginePriorityList, ruleId)
		}
		_ = d.Set("rules_priority", ruleEnginePriorityList)
	}

	return nil
}

func resourceTencentCloudTeoRuleEnginePriorityUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_rule_engine_priority.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	request := teo.NewModifyRulePriorityRequest()
	zoneId := d.Id()
	request.ZoneId = &zoneId

	if d.HasChange("rules_priority") {
		if v, ok := d.GetOk("rules_priority"); ok {
			l := len(v.([]interface{}))
			ruleIds := make([]*string, l)
			for _, item := range v.([]interface{}) {
				rule := item.(map[string]interface{})
				var index int
				var value string
				if vv, ok := rule["index"]; ok {
					index = vv.(int)
					if index > l {
						return fmt.Errorf("index is not continuous")
					}
				}
				if vv, ok := rule["value"]; ok {
					value = vv.(string)
				}
				if ruleIds[index] == nil {
					ruleIds[index] = &value
				} else {
					return fmt.Errorf("`index` [%v] is not repeatable", index)
				}
			}
			request.RuleIds = append(request.RuleIds, ruleIds...)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyRulePriority(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create teo ruleEnginePriority failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTeoRuleEnginePriorityRead(d, meta)
}

func resourceTencentCloudTeoRuleEnginePriorityDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_rule_engine_priority.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
