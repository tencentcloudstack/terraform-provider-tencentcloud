package fwm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

var fwmRuleGroupProductToName = map[string]string{
	"1": "enterprise_sg",
	"2": "cfw_nat_acl",
	"3": "cfw_vpc_acl",
	"4": "cfw_edge_acl",
}

func fwmRuleGroupNormalizeProduct(product string) string {
	if name, ok := fwmRuleGroupProductToName[product]; ok {
		return name
	}
	return product
}

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
				Description: "Rule group name, cannot be empty and the maximum length is 50 characters.",
			},

			"product": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Product the rule group belongs to. Valid values: `enterprise_sg`, `cfw_nat_acl`, `cfw_vpc_acl`, `cfw_edge_acl`.",
				ValidateFunc: validation.StringInSlice([]string{
					"enterprise_sg", "cfw_nat_acl", "cfw_vpc_acl", "cfw_edge_acl",
				}, false),
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
				Description: "Tenant AppId.",
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
				Description: "Creator SubUin.",
			},

			"create_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creator name.",
			},

			"update_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Updater SubUin.",
			},

			"update_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Updater name.",
			},

			"receive_account": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Associated accounts the rule group is delivered to.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Account Uin.",
						},
						"nickname": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Account nickname.",
						},
						"receiver_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Receiver type. `0`: account uin, `1`: account group.",
						},
						"members": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Members of the account group, returned when `receiver_type` is `1`.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"app_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Member AppId.",
									},
									"uin": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Member Uin.",
									},
									"nickname": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Member nickname.",
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
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = fwmv20250611.NewCreateRuleGroupRequest()
		response = fwmv20250611.NewCreateRuleGroupResponse()
	)

	if v, ok := d.GetOk("group_name"); ok {
		request.GroupName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product"); ok {
		request.Product = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client().CreateRuleGroupWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create fwm rule group failed, Response is nil."))
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result.Response.GroupId == nil || *result.Response.GroupId == "" {
			return resource.NonRetryableError(fmt.Errorf("Create fwm rule group failed, GroupId is empty."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create fwm rule group failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	log.Printf("[DEBUG]%s create fwm rule group succeed, logId: %s, group id: %s", logId, logId, *response.Response.GroupId)
	d.SetId(*response.Response.GroupId)
	return resourceTencentCloudFwmRuleGroupRead(d, meta)
}

func resourceTencentCloudFwmRuleGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_rule_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = FwmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	respData, err := service.DescribeFwmRuleGroupById(ctx, d.Id())
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[CRUD] tencentcloud_fwm_rule_group id=%s", d.Id())
		log.Printf("[WARN]%s resource `tencentcloud_fwm_rule_group` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
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
		_ = d.Set("product", fwmRuleGroupNormalizeProduct(*respData.Product))
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

	if respData.ReceiveAccount != nil && len(respData.ReceiveAccount) > 0 {
		receiveAccountList := make([]map[string]interface{}, 0, len(respData.ReceiveAccount))
		for _, receiveAccount := range respData.ReceiveAccount {
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

			if receiveAccount.Members != nil && len(receiveAccount.Members) > 0 {
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

		_ = d.Set("receive_account", receiveAccountList)
	}

	if respData.ExecAreas != nil && len(respData.ExecAreas) > 0 {
		execAreasList := make([]interface{}, 0, len(respData.ExecAreas))
		for _, execArea := range respData.ExecAreas {
			if execArea == nil {
				continue
			}
			execAreasList = append(execAreasList, *execArea)
		}

		_ = d.Set("exec_areas", execAreasList)
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

	immutableArgs := []string{"product"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("tencentcloud_fwm_rule_group: argument `%s` is immutable, please recreate the resource", v)
		}
	}

	mutableArgs := []string{"group_name"}
	needChange := false
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := fwmv20250611.NewModifyRuleGroupRequest()
		request.GroupId = helper.String(d.Id())

		if v, ok := d.GetOk("group_name"); ok {
			request.GroupName = helper.String(v.(string))
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client().ModifyRuleGroupWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Modify fwm rule group failed, Response is nil."))
			}

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s modify fwm rule group failed, reason:%+v", logId, reqErr)
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
	)

	request.GroupIds = []*string{helper.String(d.Id())}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client().DeleteRuleGroupWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete fwm rule group failed, Response is nil."))
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete fwm rule group failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
