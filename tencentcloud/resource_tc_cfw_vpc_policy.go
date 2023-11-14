/*
Provides a resource to create a cfw vpc_policy

Example Usage

```hcl
resource "tencentcloud_cfw_vpc_policy" "vpc_policy" {
  rules {
		source_content = "0.0.0.0/0"
		source_type = "net"
		dest_content = "192.168.0.2"
		dest_type = "net"
		protocol = "ANY"
		rule_action = "log"
		port = "-1/-1"
		description = "test vpc rule"
		order_index = 28
		uuid =
		enable = "true"
		edge_id = "ALL"
		detected_times = 0
		edge_name = ""
		internal_uuid = 0
		deleted = 0
		fw_group_id = ""
		fw_group_name = ""
		beta_list {
			task_id =
			task_name = ""
			last_time = ""
		}

  }
  from = ""
}
```

Import

cfw vpc_policy can be imported using the id, e.g.

```
terraform import tencentcloud_cfw_vpc_policy.vpc_policy vpc_policy_id
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

func resourceTencentCloudCfwVpcPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCfwVpcPolicyCreate,
		Read:   resourceTencentCloudCfwVpcPolicyRead,
		Update: resourceTencentCloudCfwVpcPolicyUpdate,
		Delete: resourceTencentCloudCfwVpcPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"rules": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "List of rules between vpc intranets that need to be added.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_content": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Access source examplnet:IP/CIDR(192.168.0.2).",
						},
						"source_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Access source type, the type can be: net.",
						},
						"dest_content": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Access purpose example:net:IP/CIDR(192.168.0.2)domain:domain rule,for example*.qq.com.",
						},
						"dest_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Access purpose type, the type can be: net, domain.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Protocol, optional value:TCPUDPICMPANYHTTPHTTPSHTTP/HTTPSSMTPSMTPSSMTP/SMTPSFTPDNSTLS/SSLNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"rule_action": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "How traffic set in the access control policy passes through the cloud firewall. Value:accept:acceptdrop:droplog:log.",
						},
						"port": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The port for the access control policy. Value: -1/-1: All ports;80: port 80.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"description": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Describe.",
						},
						"order_index": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Rule order, -1 means lowest, 1 means highest.",
						},
						"uuid": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The unique id corresponding to the rule.",
						},
						"enable": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Rule status, true means enabled, false means disabled.",
						},
						"edge_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The scope of the rule&amp;#39;s effectiveness, whether it is between a pair of vpcs or for all vpcs.",
						},
						"detected_times": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The number of hits of the rule. There is no need to pass this parameter when adding, deleting, modifying or querying rules. It is mainly used to return query result data.",
						},
						"edge_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description of the pair of inter-VPC firewalls corresponding to EdgeId.",
						},
						"internal_uuid": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Uuid used internally, this field is generally not used.",
						},
						"deleted": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The rule is deleted: 1, deleted; 0, not deleted.",
						},
						"fw_group_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Firewall instance ID where the rule takes effectNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"fw_group_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Firewall nameNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"beta_list": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Beta mission detailsNote: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"task_id": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Task idNote: This field may return null, indicating that no valid value can be obtained.",
									},
									"task_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Mission nameNote: This field may return null, indicating that no valid value can be obtained.",
									},
									"last_time": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Last execution timeNote: This field may return null, indicating that no valid value can be obtained.",
									},
								},
							},
						},
					},
				},
			},

			"from": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The source of adding rules. Generally, there is no need to use it. The value insert_rule means inserting the rules at the specified position; the value batch_import means batch importing rules; when it is empty, it means adding rules.",
			},
		},
	}
}

func resourceTencentCloudCfwVpcPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_vpc_policy.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = cfw.NewAddVpcAcRuleRequest()
		response = cfw.NewAddVpcAcRuleResponse()
		uuid     int
	)
	if v, ok := d.GetOk("rules"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			vpcRuleItem := cfw.VpcRuleItem{}
			if v, ok := dMap["source_content"]; ok {
				vpcRuleItem.SourceContent = helper.String(v.(string))
			}
			if v, ok := dMap["source_type"]; ok {
				vpcRuleItem.SourceType = helper.String(v.(string))
			}
			if v, ok := dMap["dest_content"]; ok {
				vpcRuleItem.DestContent = helper.String(v.(string))
			}
			if v, ok := dMap["dest_type"]; ok {
				vpcRuleItem.DestType = helper.String(v.(string))
			}
			if v, ok := dMap["protocol"]; ok {
				vpcRuleItem.Protocol = helper.String(v.(string))
			}
			if v, ok := dMap["rule_action"]; ok {
				vpcRuleItem.RuleAction = helper.String(v.(string))
			}
			if v, ok := dMap["port"]; ok {
				vpcRuleItem.Port = helper.String(v.(string))
			}
			if v, ok := dMap["description"]; ok {
				vpcRuleItem.Description = helper.String(v.(string))
			}
			if v, ok := dMap["order_index"]; ok {
				vpcRuleItem.OrderIndex = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["uuid"]; ok {
				vpcRuleItem.Uuid = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["enable"]; ok {
				vpcRuleItem.Enable = helper.String(v.(string))
			}
			if v, ok := dMap["edge_id"]; ok {
				vpcRuleItem.EdgeId = helper.String(v.(string))
			}
			if v, ok := dMap["detected_times"]; ok {
				vpcRuleItem.DetectedTimes = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["edge_name"]; ok {
				vpcRuleItem.EdgeName = helper.String(v.(string))
			}
			if v, ok := dMap["internal_uuid"]; ok {
				vpcRuleItem.InternalUuid = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["deleted"]; ok {
				vpcRuleItem.Deleted = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["fw_group_id"]; ok {
				vpcRuleItem.FwGroupId = helper.String(v.(string))
			}
			if v, ok := dMap["fw_group_name"]; ok {
				vpcRuleItem.FwGroupName = helper.String(v.(string))
			}
			if v, ok := dMap["beta_list"]; ok {
				for _, item := range v.([]interface{}) {
					betaListMap := item.(map[string]interface{})
					betaInfoByACL := cfw.BetaInfoByACL{}
					if v, ok := betaListMap["task_id"]; ok {
						betaInfoByACL.TaskId = helper.IntInt64(v.(int))
					}
					if v, ok := betaListMap["task_name"]; ok {
						betaInfoByACL.TaskName = helper.String(v.(string))
					}
					if v, ok := betaListMap["last_time"]; ok {
						betaInfoByACL.LastTime = helper.String(v.(string))
					}
					vpcRuleItem.BetaList = append(vpcRuleItem.BetaList, &betaInfoByACL)
				}
			}
			request.Rules = append(request.Rules, &vpcRuleItem)
		}
	}

	if v, ok := d.GetOk("from"); ok {
		request.From = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCfwClient().AddVpcAcRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cfw vpcPolicy failed, reason:%+v", logId, err)
		return err
	}

	uuid = *response.Response.Uuid
	d.SetId(helper.Int64ToStr(uuid))

	return resourceTencentCloudCfwVpcPolicyRead(d, meta)
}

func resourceTencentCloudCfwVpcPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_vpc_policy.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CfwService{client: meta.(*TencentCloudClient).apiV3Conn}

	vpcPolicyId := d.Id()

	vpcPolicy, err := service.DescribeCfwVpcPolicyById(ctx, uuid)
	if err != nil {
		return err
	}

	if vpcPolicy == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CfwVpcPolicy` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if vpcPolicy.Rules != nil {
		rulesList := []interface{}{}
		for _, rules := range vpcPolicy.Rules {
			rulesMap := map[string]interface{}{}

			if vpcPolicy.Rules.SourceContent != nil {
				rulesMap["source_content"] = vpcPolicy.Rules.SourceContent
			}

			if vpcPolicy.Rules.SourceType != nil {
				rulesMap["source_type"] = vpcPolicy.Rules.SourceType
			}

			if vpcPolicy.Rules.DestContent != nil {
				rulesMap["dest_content"] = vpcPolicy.Rules.DestContent
			}

			if vpcPolicy.Rules.DestType != nil {
				rulesMap["dest_type"] = vpcPolicy.Rules.DestType
			}

			if vpcPolicy.Rules.Protocol != nil {
				rulesMap["protocol"] = vpcPolicy.Rules.Protocol
			}

			if vpcPolicy.Rules.RuleAction != nil {
				rulesMap["rule_action"] = vpcPolicy.Rules.RuleAction
			}

			if vpcPolicy.Rules.Port != nil {
				rulesMap["port"] = vpcPolicy.Rules.Port
			}

			if vpcPolicy.Rules.Description != nil {
				rulesMap["description"] = vpcPolicy.Rules.Description
			}

			if vpcPolicy.Rules.OrderIndex != nil {
				rulesMap["order_index"] = vpcPolicy.Rules.OrderIndex
			}

			if vpcPolicy.Rules.Uuid != nil {
				rulesMap["uuid"] = vpcPolicy.Rules.Uuid
			}

			if vpcPolicy.Rules.Enable != nil {
				rulesMap["enable"] = vpcPolicy.Rules.Enable
			}

			if vpcPolicy.Rules.EdgeId != nil {
				rulesMap["edge_id"] = vpcPolicy.Rules.EdgeId
			}

			if vpcPolicy.Rules.DetectedTimes != nil {
				rulesMap["detected_times"] = vpcPolicy.Rules.DetectedTimes
			}

			if vpcPolicy.Rules.EdgeName != nil {
				rulesMap["edge_name"] = vpcPolicy.Rules.EdgeName
			}

			if vpcPolicy.Rules.InternalUuid != nil {
				rulesMap["internal_uuid"] = vpcPolicy.Rules.InternalUuid
			}

			if vpcPolicy.Rules.Deleted != nil {
				rulesMap["deleted"] = vpcPolicy.Rules.Deleted
			}

			if vpcPolicy.Rules.FwGroupId != nil {
				rulesMap["fw_group_id"] = vpcPolicy.Rules.FwGroupId
			}

			if vpcPolicy.Rules.FwGroupName != nil {
				rulesMap["fw_group_name"] = vpcPolicy.Rules.FwGroupName
			}

			if vpcPolicy.Rules.BetaList != nil {
				betaListList := []interface{}{}
				for _, betaList := range vpcPolicy.Rules.BetaList {
					betaListMap := map[string]interface{}{}

					if betaList.TaskId != nil {
						betaListMap["task_id"] = betaList.TaskId
					}

					if betaList.TaskName != nil {
						betaListMap["task_name"] = betaList.TaskName
					}

					if betaList.LastTime != nil {
						betaListMap["last_time"] = betaList.LastTime
					}

					betaListList = append(betaListList, betaListMap)
				}

				rulesMap["beta_list"] = []interface{}{betaListList}
			}

			rulesList = append(rulesList, rulesMap)
		}

		_ = d.Set("rules", rulesList)

	}

	if vpcPolicy.From != nil {
		_ = d.Set("from", vpcPolicy.From)
	}

	return nil
}

func resourceTencentCloudCfwVpcPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_vpc_policy.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cfw.NewModifyVpcAcRuleRequest()

	vpcPolicyId := d.Id()

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
				vpcRuleItem := cfw.VpcRuleItem{}
				if v, ok := dMap["source_content"]; ok {
					vpcRuleItem.SourceContent = helper.String(v.(string))
				}
				if v, ok := dMap["source_type"]; ok {
					vpcRuleItem.SourceType = helper.String(v.(string))
				}
				if v, ok := dMap["dest_content"]; ok {
					vpcRuleItem.DestContent = helper.String(v.(string))
				}
				if v, ok := dMap["dest_type"]; ok {
					vpcRuleItem.DestType = helper.String(v.(string))
				}
				if v, ok := dMap["protocol"]; ok {
					vpcRuleItem.Protocol = helper.String(v.(string))
				}
				if v, ok := dMap["rule_action"]; ok {
					vpcRuleItem.RuleAction = helper.String(v.(string))
				}
				if v, ok := dMap["port"]; ok {
					vpcRuleItem.Port = helper.String(v.(string))
				}
				if v, ok := dMap["description"]; ok {
					vpcRuleItem.Description = helper.String(v.(string))
				}
				if v, ok := dMap["order_index"]; ok {
					vpcRuleItem.OrderIndex = helper.IntInt64(v.(int))
				}
				if v, ok := dMap["uuid"]; ok {
					vpcRuleItem.Uuid = helper.IntInt64(v.(int))
				}
				if v, ok := dMap["enable"]; ok {
					vpcRuleItem.Enable = helper.String(v.(string))
				}
				if v, ok := dMap["edge_id"]; ok {
					vpcRuleItem.EdgeId = helper.String(v.(string))
				}
				if v, ok := dMap["detected_times"]; ok {
					vpcRuleItem.DetectedTimes = helper.IntInt64(v.(int))
				}
				if v, ok := dMap["edge_name"]; ok {
					vpcRuleItem.EdgeName = helper.String(v.(string))
				}
				if v, ok := dMap["internal_uuid"]; ok {
					vpcRuleItem.InternalUuid = helper.IntInt64(v.(int))
				}
				if v, ok := dMap["deleted"]; ok {
					vpcRuleItem.Deleted = helper.IntInt64(v.(int))
				}
				if v, ok := dMap["fw_group_id"]; ok {
					vpcRuleItem.FwGroupId = helper.String(v.(string))
				}
				if v, ok := dMap["fw_group_name"]; ok {
					vpcRuleItem.FwGroupName = helper.String(v.(string))
				}
				if v, ok := dMap["beta_list"]; ok {
					for _, item := range v.([]interface{}) {
						betaListMap := item.(map[string]interface{})
						betaInfoByACL := cfw.BetaInfoByACL{}
						if v, ok := betaListMap["task_id"]; ok {
							betaInfoByACL.TaskId = helper.IntInt64(v.(int))
						}
						if v, ok := betaListMap["task_name"]; ok {
							betaInfoByACL.TaskName = helper.String(v.(string))
						}
						if v, ok := betaListMap["last_time"]; ok {
							betaInfoByACL.LastTime = helper.String(v.(string))
						}
						vpcRuleItem.BetaList = append(vpcRuleItem.BetaList, &betaInfoByACL)
					}
				}
				request.Rules = append(request.Rules, &vpcRuleItem)
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCfwClient().ModifyVpcAcRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cfw vpcPolicy failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCfwVpcPolicyRead(d, meta)
}

func resourceTencentCloudCfwVpcPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_vpc_policy.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CfwService{client: meta.(*TencentCloudClient).apiV3Conn}
	vpcPolicyId := d.Id()

	if err := service.DeleteCfwVpcPolicyById(ctx, uuid); err != nil {
		return err
	}

	return nil
}
