package fwm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudFwmRuleGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudFwmRuleGroupCreate,
		Read:   resourceTencentCloudFwmRuleGroupRead,
		Update: resourceTencentCloudFwmRuleGroupUpdate,
		Delete: resourceTencentCloudFwmRuleGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Rule group name, cannot be empty, up to 50 characters.",
			},

			"product": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"enterprise_sg", "cfw_nat_acl", "cfw_vpc_acl", "cfw_edge_acl"}),
				Description:  "Product the rule group belongs to. Valid values: `enterprise_sg` (enterprise security group), `cfw_vpc_acl` (VPC border firewall), `cfw_nat_acl` (NAT border firewall), `cfw_edge_acl` (internet border firewall).",
			},

			// computed
			"group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule group ID.",
			},

			"app_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Tenant APP ID.",
			},

			"rule_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of rules in the rule group.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Update time.",
			},

			"create_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creator sub account UIN.",
			},

			"create_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creator name.",
			},

			"update_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Updater sub account UIN.",
			},

			"update_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Updater name.",
			},

			"receive_account": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Associated delivery accounts.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Tenant UIN.",
						},
						"nickname": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Tenant name.",
						},
						"receiver_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Receiver type. `0`: account UIN, `1`: account group.",
						},
						"members": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Members of the account group, only returned when `receiver_type` is `1`.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"app_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Member APP ID.",
									},
									"uin": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Member UIN.",
									},
									"nickname": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Member name.",
									},
									"member_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Member ID.",
									},
									"node_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Member node name.",
									},
								},
							},
						},
					},
				},
			},

			"exec_areas": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Execution stages.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"direction": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Direction. `1`: inbound, `0`: outbound, `-1`: no direction.",
			},
		},
	}
}

func resourceTencentCloudFwmRuleGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_rule_group.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = NewFwmService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		request = fwmv20250611.NewCreateRuleGroupRequest()
	)

	if v, ok := d.GetOk("group_name"); ok {
		request.GroupName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product"); ok {
		request.Product = helper.String(v.(string))
	}

	var groupId string
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client().CreateRuleGroupWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create tencentcloud_fwm_rule_group failed, Response is nil."))
		}

		if result.Response.GroupId == nil || *result.Response.GroupId == "" {
			return resource.NonRetryableError(fmt.Errorf("Create tencentcloud_fwm_rule_group failed, GroupId is empty."))
		}

		groupId = *result.Response.GroupId
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create tencentcloud_fwm_rule_group failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	log.Printf("[DEBUG]%s create tencentcloud_fwm_rule_group success, groupId: %s", logId, groupId)
	d.SetId(groupId)

	pollErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ruleGroup, e := service.DescribeFwmRuleGroupById(ctx, groupId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if ruleGroup == nil {
			return resource.RetryableError(fmt.Errorf("tencentcloud_fwm_rule_group %s is not visible yet", groupId))
		}
		return nil
	})

	if pollErr != nil {
		log.Printf("[CRITAL]%s wait tencentcloud_fwm_rule_group %s visible failed, reason:%+v", logId, groupId, pollErr)
		return pollErr
	}

	return resourceTencentCloudFwmRuleGroupRead(d, meta)
}

func resourceTencentCloudFwmRuleGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_rule_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = NewFwmService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		groupId = d.Id()
	)

	respData, err := service.DescribeFwmRuleGroupById(ctx, groupId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[CRUD] tencentcloud_fwm_rule_group id=%s", d.Id())
		d.SetId("")
		return nil
	}

	if respData.GroupId != nil {
		_ = d.Set("group_id", respData.GroupId)
	}

	if respData.GroupName != nil {
		_ = d.Set("group_name", respData.GroupName)
	}

	if respData.Product != nil {
		_ = d.Set("product", fwmRuleGroupProductCodeToName(*respData.Product))
	}

	if respData.AppId != nil {
		_ = d.Set("app_id", respData.AppId)
	}

	if respData.RuleCount != nil {
		_ = d.Set("rule_count", respData.RuleCount)
	}

	if respData.CreateTime != nil {
		_ = d.Set("create_time", respData.CreateTime)
	}

	if respData.UpdateTime != nil {
		_ = d.Set("update_time", respData.UpdateTime)
	}

	if respData.CreateBy != nil {
		_ = d.Set("create_by", respData.CreateBy)
	}

	if respData.CreateName != nil {
		_ = d.Set("create_name", respData.CreateName)
	}

	if respData.UpdateBy != nil {
		_ = d.Set("update_by", respData.UpdateBy)
	}

	if respData.UpdateName != nil {
		_ = d.Set("update_name", respData.UpdateName)
	}

	if respData.ReceiveAccount != nil {
		_ = d.Set("receive_account", fwmRuleGroupReceiveAccountsToMaps(respData.ReceiveAccount))
	}

	if respData.ExecAreas != nil {
		execAreas := make([]string, 0, len(respData.ExecAreas))
		for _, execArea := range respData.ExecAreas {
			if execArea == nil {
				continue
			}
			execAreas = append(execAreas, *execArea)
		}
		_ = d.Set("exec_areas", execAreas)
	}

	if respData.Direction != nil {
		_ = d.Set("direction", respData.Direction)
	}

	return nil
}

func resourceTencentCloudFwmRuleGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_rule_group.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	if d.HasChange("group_name") {
		request := fwmv20250611.NewModifyRuleGroupRequest()
		request.GroupId = helper.String(d.Id())

		if v, ok := d.GetOk("group_name"); ok {
			request.GroupName = helper.String(v.(string))
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client().ModifyRuleGroupWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Modify tencentcloud_fwm_rule_group failed, Response is nil."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update tencentcloud_fwm_rule_group failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudFwmRuleGroupRead(d, meta)
}

func resourceTencentCloudFwmRuleGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_rule_group.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = fwmv20250611.NewDeleteRuleGroupRequest()
		groupId = d.Id()
	)

	request.GroupIds = []*string{&groupId}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client().DeleteRuleGroupWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete tencentcloud_fwm_rule_group failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete tencentcloud_fwm_rule_group failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}

func fwmRuleGroupProductCodeToName(v string) string {
	switch v {
	case "1":
		return "enterprise_sg"
	case "2":
		return "cfw_nat_acl"
	case "3":
		return "cfw_vpc_acl"
	case "4":
		return "cfw_edge_acl"
	default:
		return v
	}
}

func fwmRuleGroupReceiveAccountsToMaps(receiveAccounts []*fwmv20250611.ReceiveAccount) []map[string]interface{} {
	receiveAccountList := make([]map[string]interface{}, 0, len(receiveAccounts))
	for _, receiveAccount := range receiveAccounts {
		if receiveAccount == nil {
			continue
		}

		receiveAccountMap := map[string]interface{}{}
		if receiveAccount.Uin != nil {
			receiveAccountMap["uin"] = receiveAccount.Uin
		}
		if receiveAccount.Nickname != nil {
			receiveAccountMap["nickname"] = receiveAccount.Nickname
		}
		if receiveAccount.ReceiverType != nil {
			receiveAccountMap["receiver_type"] = receiveAccount.ReceiverType
		}
		if receiveAccount.Members != nil {
			membersList := make([]map[string]interface{}, 0, len(receiveAccount.Members))
			for _, member := range receiveAccount.Members {
				if member == nil {
					continue
				}

				memberMap := map[string]interface{}{}
				if member.AppId != nil {
					memberMap["app_id"] = member.AppId
				}
				if member.Uin != nil {
					memberMap["uin"] = member.Uin
				}
				if member.Nickname != nil {
					memberMap["nickname"] = member.Nickname
				}
				if member.MemberId != nil {
					memberMap["member_id"] = member.MemberId
				}
				if member.NodeName != nil {
					memberMap["node_name"] = member.NodeName
				}
				membersList = append(membersList, memberMap)
			}

			receiveAccountMap["members"] = membersList
		}

		receiveAccountList = append(receiveAccountList, receiveAccountMap)
	}

	return receiveAccountList
}
