/*
Provides a resource to create a tsf enable_unit_rule

Example Usage

```hcl
resource "tencentcloud_tsf_enable_unit_rule" "enable_unit_rule" {
  rule_id = "unit-rl-is9m4nxz"
  switch = "enabled"
}
```

Import

tsf enable_unit_rule can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_enable_unit_rule.enable_unit_rule enable_unit_rule_id
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
)

func resourceTencentCloudTsfEnableUnitRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfEnableUnitRuleCreate,
		Read:   resourceTencentCloudTsfEnableUnitRuleRead,
		Update: resourceTencentCloudTsfEnableUnitRuleUpdate,
		Delete: resourceTencentCloudTsfEnableUnitRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"rule_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "api ID.",
			},

			"switch": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "switch, on: `enabled`, off: `disabled`.",
			},
		},
	}
}

func resourceTencentCloudTsfEnableUnitRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_enable_unit_rule.create")()
	defer inconsistentCheck(d, meta)()

	var id string
	if v, ok := d.GetOk("rule_id"); ok {
		id = v.(string)
	}

	d.SetId(id)

	return resourceTencentCloudTsfEnableUnitRuleUpdate(d, meta)
}

func resourceTencentCloudTsfEnableUnitRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_enable_unit_rule.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	id := d.Id()

	enableUnitRule, err := service.DescribeTsfEnableUnitRuleById(ctx, id)
	if err != nil {
		return err
	}

	if enableUnitRule == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfEnableUnitRule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if enableUnitRule.Id != nil {
		_ = d.Set("rule_id", enableUnitRule.Id)
	}

	if enableUnitRule.Status != nil {
		_ = d.Set("switch", enableUnitRule.Status)
	}

	return nil
}

func resourceTencentCloudTsfEnableUnitRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_enable_unit_rule.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	id := d.Id()
	if v, ok := d.GetOk("switch"); ok {
		if v.(string) == "enabled" {
			request := tsf.NewEnableUnitRuleRequest()
			request.Id = &id
			err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().EnableUnitRule(request)
				if e != nil {
					return retryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s update tsf enableUnitRule failed, reason:%+v", logId, err)
				return err
			}
		}

		if v.(string) == "disabled" {
			request := tsf.NewDisableUnitRuleRequest()
			request.Id = &id
			err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().DisableUnitRule(request)
				if e != nil {
					return retryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s update tsf disableUnitRule failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	return resourceTencentCloudTsfEnableUnitRuleRead(d, meta)
}

func resourceTencentCloudTsfEnableUnitRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_enable_unit_rule.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
