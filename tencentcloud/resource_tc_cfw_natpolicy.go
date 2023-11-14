/*
Provides a resource to create a cfw natpolicy

Example Usage

```hcl
resource "tencentcloud_cfw_natpolicy" "natpolicy" {
  rules {
		source_content = "192.168.0.2"
		source_type = "ip"
		target_content = "192.168.0.2"
		target_type = "ip"
		protocol = "TCP"
		rule_action = "allow"
		port = "80"
		direction = 1
		order_index = 1
		enable = "true"
		uuid = 1
		description = "test"

  }
  from = ""
}
```

Import

cfw natpolicy can be imported using the id, e.g.

```
terraform import tencentcloud_cfw_natpolicy.natpolicy natpolicy_id
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
)

func resourceTencentCloudCfwNatpolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCfwNatpolicyCreate,
		Read:   resourceTencentCloudCfwNatpolicyRead,
		Update: resourceTencentCloudCfwNatpolicyUpdate,
		Delete: resourceTencentCloudCfwNatpolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"rules": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "NAT access control rules to be added.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_content": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Access source. Example net/IP/CIDR(192.168.0.2).",
						},
						"source_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Access source type. Values for inbound rules ip, net, template, and location. Values for outbound rules ip, net, template, instance, group, and tag.",
						},
						"target_content": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Access target. Example net IP/CIDR(192.168.0.2); domain domain name rule, e.g., *.qq.comccc.",
						},
						"target_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Access target type. Values for inbound rules ip, net, template, instance, group, and tag. Values for outbound rules ip, net, domain, template, and location.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ProtocolWhen Direction=1, optional values: TCP, UDP, ANY.When Direction=0, optional values: TCP, UDP, ICMP, ANY, HTTP, HTTPS, HTTP/HTTPS, SMTP, SMTPS, SMTP/SMTPS, FTP, and DNS.",
						},
						"rule_action": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specify how the CFW instance deals with the traffic hit the access control rule. Values accept (allow), drop (reject), and log (observe).",
						},
						"port": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The port of the access control rule. Values -1/-1 (all ports) and 80 (Port 80).",
						},
						"direction": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Rule direction. Values 1 (Inbound) and 0 (Outbound).",
						},
						"order_index": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Rule sequence number.",
						},
						"enable": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Rule status. true (Enabled); false (Disabled).",
						},
						"uuid": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The unique ID of the rule, which is not required when you create a rule.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description.",
						},
					},
				},
			},

			"from": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Source of the rules to be added. Generally, this parameter is not used. The value insert_rule indicates that rules in the specified location are inserted, and the value batch_import indicates that rules are imported in batches. If the parameter is left empty, rules defined in the API request are added.",
			},
		},
	}
}

func resourceTencentCloudCfwNatpolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_natpolicy.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = cfw.NewAddNatAcRuleRequest()
		response = cfw.NewAddNatAcRuleResponse()
		uuid     int
	)
	if v, ok := d.GetOk("rules"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			createNatRuleItem := cfw.CreateNatRuleItem{}
			if v, ok := dMap["source_content"]; ok {
				createNatRuleItem.SourceContent = helper.String(v.(string))
			}
			if v, ok := dMap["source_type"]; ok {
				createNatRuleItem.SourceType = helper.String(v.(string))
			}
			if v, ok := dMap["target_content"]; ok {
				createNatRuleItem.TargetContent = helper.String(v.(string))
			}
			if v, ok := dMap["target_type"]; ok {
				createNatRuleItem.TargetType = helper.String(v.(string))
			}
			if v, ok := dMap["protocol"]; ok {
				createNatRuleItem.Protocol = helper.String(v.(string))
			}
			if v, ok := dMap["rule_action"]; ok {
				createNatRuleItem.RuleAction = helper.String(v.(string))
			}
			if v, ok := dMap["port"]; ok {
				createNatRuleItem.Port = helper.String(v.(string))
			}
			if v, ok := dMap["direction"]; ok {
				createNatRuleItem.Direction = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["order_index"]; ok {
				createNatRuleItem.OrderIndex = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["enable"]; ok {
				createNatRuleItem.Enable = helper.String(v.(string))
			}
			if v, ok := dMap["uuid"]; ok {
				createNatRuleItem.Uuid = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["description"]; ok {
				createNatRuleItem.Description = helper.String(v.(string))
			}
			request.Rules = append(request.Rules, &createNatRuleItem)
		}
	}

	if v, ok := d.GetOk("from"); ok {
		request.From = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCfwClient().AddNatAcRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cfw natpolicy failed, reason:%+v", logId, err)
		return err
	}

	uuid = *response.Response.Uuid
	d.SetId(helper.Int64ToStr(uuid))

	return resourceTencentCloudCfwNatpolicyRead(d, meta)
}

func resourceTencentCloudCfwNatpolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_natpolicy.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CfwService{client: meta.(*TencentCloudClient).apiV3Conn}

	natpolicyId := d.Id()

	natpolicy, err := service.DescribeCfwNatpolicyById(ctx, uuid)
	if err != nil {
		return err
	}

	if natpolicy == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CfwNatpolicy` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if natpolicy.Rules != nil {
		rulesList := []interface{}{}
		for _, rules := range natpolicy.Rules {
			rulesMap := map[string]interface{}{}

			if natpolicy.Rules.SourceContent != nil {
				rulesMap["source_content"] = natpolicy.Rules.SourceContent
			}

			if natpolicy.Rules.SourceType != nil {
				rulesMap["source_type"] = natpolicy.Rules.SourceType
			}

			if natpolicy.Rules.TargetContent != nil {
				rulesMap["target_content"] = natpolicy.Rules.TargetContent
			}

			if natpolicy.Rules.TargetType != nil {
				rulesMap["target_type"] = natpolicy.Rules.TargetType
			}

			if natpolicy.Rules.Protocol != nil {
				rulesMap["protocol"] = natpolicy.Rules.Protocol
			}

			if natpolicy.Rules.RuleAction != nil {
				rulesMap["rule_action"] = natpolicy.Rules.RuleAction
			}

			if natpolicy.Rules.Port != nil {
				rulesMap["port"] = natpolicy.Rules.Port
			}

			if natpolicy.Rules.Direction != nil {
				rulesMap["direction"] = natpolicy.Rules.Direction
			}

			if natpolicy.Rules.OrderIndex != nil {
				rulesMap["order_index"] = natpolicy.Rules.OrderIndex
			}

			if natpolicy.Rules.Enable != nil {
				rulesMap["enable"] = natpolicy.Rules.Enable
			}

			if natpolicy.Rules.Uuid != nil {
				rulesMap["uuid"] = natpolicy.Rules.Uuid
			}

			if natpolicy.Rules.Description != nil {
				rulesMap["description"] = natpolicy.Rules.Description
			}

			rulesList = append(rulesList, rulesMap)
		}

		_ = d.Set("rules", rulesList)

	}

	if natpolicy.From != nil {
		_ = d.Set("from", natpolicy.From)
	}

	return nil
}

func resourceTencentCloudCfwNatpolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_natpolicy.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cfw.NewModifyNatAcRuleRequest()

	natpolicyId := d.Id()

	request.Uuid = &uuid

	immutableArgs := []string{"rules", "from"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("rules") {
		if v, ok := d.GetOk("rules"); ok {
			for _, item := range v.([]interface{}) {
				createNatRuleItem := cfw.CreateNatRuleItem{}
				if v, ok := dMap["source_content"]; ok {
					createNatRuleItem.SourceContent = helper.String(v.(string))
				}
				if v, ok := dMap["source_type"]; ok {
					createNatRuleItem.SourceType = helper.String(v.(string))
				}
				if v, ok := dMap["target_content"]; ok {
					createNatRuleItem.TargetContent = helper.String(v.(string))
				}
				if v, ok := dMap["target_type"]; ok {
					createNatRuleItem.TargetType = helper.String(v.(string))
				}
				if v, ok := dMap["protocol"]; ok {
					createNatRuleItem.Protocol = helper.String(v.(string))
				}
				if v, ok := dMap["rule_action"]; ok {
					createNatRuleItem.RuleAction = helper.String(v.(string))
				}
				if v, ok := dMap["port"]; ok {
					createNatRuleItem.Port = helper.String(v.(string))
				}
				if v, ok := dMap["direction"]; ok {
					createNatRuleItem.Direction = helper.IntUint64(v.(int))
				}
				if v, ok := dMap["order_index"]; ok {
					createNatRuleItem.OrderIndex = helper.IntInt64(v.(int))
				}
				if v, ok := dMap["enable"]; ok {
					createNatRuleItem.Enable = helper.String(v.(string))
				}
				if v, ok := dMap["uuid"]; ok {
					createNatRuleItem.Uuid = helper.IntInt64(v.(int))
				}
				if v, ok := dMap["description"]; ok {
					createNatRuleItem.Description = helper.String(v.(string))
				}
				request.Rules = append(request.Rules, &createNatRuleItem)
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCfwClient().ModifyNatAcRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cfw natpolicy failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCfwNatpolicyRead(d, meta)
}

func resourceTencentCloudCfwNatpolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_natpolicy.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CfwService{client: meta.(*TencentCloudClient).apiV3Conn}
	natpolicyId := d.Id()

	if err := service.DeleteCfwNatpolicyById(ctx, uuid); err != nil {
		return err
	}

	return nil
}
