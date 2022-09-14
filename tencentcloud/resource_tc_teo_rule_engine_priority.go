/*
Provides a resource to create a teo rule_engine_priority

Example Usage

```hcl
resource "tencentcloud_teo_rule_engine_priority" "rule_engine_priority" {
  rules_priority = ""
  tags = {
    "createdBy" = "terraform"
  }
}

```
Import

teo rule_engine_priority can be imported using the id, e.g.
```
$ terraform import tencentcloud_teo_rule_engine_priority.rule_engine_priority ruleEnginePriority_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220106"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
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
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Computed:    true,
				Description: "Priority of rules.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
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

	err := resourceTencentCloudTeoRuleEnginePriorityUpdate(d, meta)
	if err != nil {
		log.Printf("[CRITAL]%s create teo ruleEnginePriority failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(zoneId)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::teo:%s:uin/:zone/%s", region, zoneId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
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

	if ruleEnginePriority != nil {
		var ruleIds []string
		for _, v := range ruleEnginePriority {
			ruleIds = append(ruleIds, *v.RuleId)
		}
		_ = d.Set("rules_priority", ruleIds)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "teo", "zone", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudTeoRuleEnginePriorityUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_rule_engine_priority.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := teo.NewModifyRulePriorityRequest()

	zoneId := d.Id()
	request.ZoneId = &zoneId

	if d.HasChange("rules_priority") {
		if v, ok := d.GetOk("rules_priority"); ok {
			rulesPrioritySet := v.(*schema.Set).List()
			for i := range rulesPrioritySet {
				rulesPriority := rulesPrioritySet[i].(string)
				request.RuleIds = append(request.RuleIds, &rulesPriority)
			}
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

	if d.HasChange("tags") {
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("teo", "zone", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudTeoRuleEnginePriorityRead(d, meta)
}

func resourceTencentCloudTeoRuleEnginePriorityDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_rule_engine_priority.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
