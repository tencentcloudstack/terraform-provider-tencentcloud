/*
Provides a resource to create a cfw block_ignore_list

Example Usage

```hcl
resource "tencentcloud_cfw_block_ignore_list" "block_ignore_list" {
  rules {
		direction = 0
		end_time = "2023-09-09 15:04:05"
		i_p = "1.1.1.1"
		domain = ""
		comment = "block ip 1.1.1.1"
		start_time = "2023-09-01 15:04:05"

  }
  rule_type = 1
}
```

Import

cfw block_ignore_list can be imported using the id, e.g.

```
terraform import tencentcloud_cfw_block_ignore_list.block_ignore_list block_ignore_list_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfw "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfw/v20190904"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudCfwBlockIgnoreList() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCfwBlockIgnoreListCreate,
		Read:   resourceTencentCloudCfwBlockIgnoreListRead,
		Update: resourceTencentCloudCfwBlockIgnoreListUpdate,
		Delete: resourceTencentCloudCfwBlockIgnoreListDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"rules": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "Rule list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"direction": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Rule direction, 0 outbound, 1 inbound, 3 intranet.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Rule end time, format: 2006-01-02 15:04:05, must be greater than the current time.",
						},
						"i_p": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Rule IP address, one of IP and Domain is required.",
						},
						"domain": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Rule domain name, one of IP and Domain is required.",
						},
						"comment": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Remarks information, length cannot exceed 50.",
						},
						"start_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Rule start time.",
						},
					},
				},
			},

			"rule_type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Rule type, 1 block, 2 ignore, domain block is not supported.",
			},
		},
	}
}

func resourceTencentCloudCfwBlockIgnoreListCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_block_ignore_list.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = cfw.NewCreateBlockIgnoreRuleListRequest()
		response  = cfw.NewCreateBlockIgnoreRuleListResponse()
		direction int
		iP        string
		domain    string
	)
	if v, ok := d.GetOk("rules"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			intrusionDefenseRule := cfw.IntrusionDefenseRule{}
			if v, ok := dMap["direction"]; ok {
				intrusionDefenseRule.Direction = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["end_time"]; ok {
				intrusionDefenseRule.EndTime = helper.String(v.(string))
			}
			if v, ok := dMap["i_p"]; ok {
				intrusionDefenseRule.IP = helper.String(v.(string))
			}
			if v, ok := dMap["domain"]; ok {
				intrusionDefenseRule.Domain = helper.String(v.(string))
			}
			if v, ok := dMap["comment"]; ok {
				intrusionDefenseRule.Comment = helper.String(v.(string))
			}
			if v, ok := dMap["start_time"]; ok {
				intrusionDefenseRule.StartTime = helper.String(v.(string))
			}
			request.Rules = append(request.Rules, &intrusionDefenseRule)
		}
	}

	if v, ok := d.GetOkExists("rule_type"); ok {
		request.RuleType = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCfwClient().CreateBlockIgnoreRuleList(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cfw blockIgnoreList failed, reason:%+v", logId, err)
		return err
	}

	direction = *response.Response.Direction
	d.SetId(strings.Join([]string{helper.Int64ToStr(direction), iP, domain}, FILED_SP))

	return resourceTencentCloudCfwBlockIgnoreListRead(d, meta)
}

func resourceTencentCloudCfwBlockIgnoreListRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_block_ignore_list.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CfwService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	direction := idSplit[0]
	iP := idSplit[1]
	domain := idSplit[2]

	blockIgnoreList, err := service.DescribeCfwBlockIgnoreListById(ctx, direction, iP, domain)
	if err != nil {
		return err
	}

	if blockIgnoreList == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CfwBlockIgnoreList` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if blockIgnoreList.Rules != nil {
		rulesList := []interface{}{}
		for _, rules := range blockIgnoreList.Rules {
			rulesMap := map[string]interface{}{}

			if blockIgnoreList.Rules.Direction != nil {
				rulesMap["direction"] = blockIgnoreList.Rules.Direction
			}

			if blockIgnoreList.Rules.EndTime != nil {
				rulesMap["end_time"] = blockIgnoreList.Rules.EndTime
			}

			if blockIgnoreList.Rules.IP != nil {
				rulesMap["i_p"] = blockIgnoreList.Rules.IP
			}

			if blockIgnoreList.Rules.Domain != nil {
				rulesMap["domain"] = blockIgnoreList.Rules.Domain
			}

			if blockIgnoreList.Rules.Comment != nil {
				rulesMap["comment"] = blockIgnoreList.Rules.Comment
			}

			if blockIgnoreList.Rules.StartTime != nil {
				rulesMap["start_time"] = blockIgnoreList.Rules.StartTime
			}

			rulesList = append(rulesList, rulesMap)
		}

		_ = d.Set("rules", rulesList)

	}

	if blockIgnoreList.RuleType != nil {
		_ = d.Set("rule_type", blockIgnoreList.RuleType)
	}

	return nil
}

func resourceTencentCloudCfwBlockIgnoreListUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_block_ignore_list.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cfw.NewModifyBlockIgnoreRuleRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	direction := idSplit[0]
	iP := idSplit[1]
	domain := idSplit[2]

	request.Direction = &direction
	request.IP = &iP
	request.Domain = &domain

	immutableArgs := []string{"rules", "rule_type"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("rule_type") {
		if v, ok := d.GetOkExists("rule_type"); ok {
			request.RuleType = helper.IntInt64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCfwClient().ModifyBlockIgnoreRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cfw blockIgnoreList failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCfwBlockIgnoreListRead(d, meta)
}

func resourceTencentCloudCfwBlockIgnoreListDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_block_ignore_list.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CfwService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	direction := idSplit[0]
	iP := idSplit[1]
	domain := idSplit[2]

	if err := service.DeleteCfwBlockIgnoreListById(ctx, direction, iP, domain); err != nil {
		return err
	}

	return nil
}
