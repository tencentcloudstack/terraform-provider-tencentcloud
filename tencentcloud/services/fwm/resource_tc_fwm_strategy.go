package fwm

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudFwmStrategy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudFwmStrategyCreate,
		Read:   resourceTencentCloudFwmStrategyRead,
		Update: resourceTencentCloudFwmStrategyUpdate,
		Delete: resourceTencentCloudFwmStrategyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"product": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Product type, e.g. enterprise_sg / cfw_vpc_acl / cfw_nat_acl / cfw_edge_acl.",
			},

			"sub_account_uin": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Sub account UIN. Declarative parameter: the FWM cloud API does not expose a corresponding input field, so it is not sent on Create/Modify and is not read back. The value is retained in state.",
			},

			"exec_area": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Execution area, pre or post. Determines whether the strategy is created in PreStrategy or PostStrategy.",
			},

			"group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Rule group ID.",
			},

			"sequence": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Priority, starting from 1.",
			},

			"receive_account": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of receive account UINs.",
			},

			"receive_group": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of receive account group IDs.",
			},

			"strategy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Strategy ID.",
			},

			"group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule group name.",
			},

			"rule_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Rule count.",
			},

			"rule_status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Strategy status.",
			},

			"err_msg": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Dispatch failure reason.",
			},

			"error_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Dispatch failure reason type.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Update time.",
			},

			"create_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creator.",
			},

			"create_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creator name.",
			},

			"update_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Updater.",
			},

			"update_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Updater name.",
			},
		},
	}
}

func resourceTencentCloudFwmStrategyCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_strategy.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = fwmv20250611.NewCreateStrategyRequest()
		response = fwmv20250611.NewCreateStrategyResponse()
	)

	if v, ok := d.GetOk("product"); ok {
		request.Product = helper.String(v.(string))
	}

	if v, ok := d.GetOk("receive_account"); ok {
		accounts := v.([]interface{})
		if len(accounts) > 0 {
			accountList := make([]string, 0, len(accounts))
			for _, account := range accounts {
				accountList = append(accountList, account.(string))
			}
			request.ReceiveAccount = helper.Strings(accountList)
		}
	}

	if v, ok := d.GetOk("receive_group"); ok {
		groups := v.([]interface{})
		if len(groups) > 0 {
			groupList := make([]string, 0, len(groups))
			for _, group := range groups {
				groupList = append(groupList, group.(string))
			}
			request.ReceiveGroup = helper.Strings(groupList)
		}
	}

	execArea := d.Get("exec_area").(string)
	strategyReq := fwmv20250611.StrategyReq{}
	if v, ok := d.GetOk("group_id"); ok {
		strategyReq.GroupId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("sequence"); ok {
		strategyReq.Sequence = helper.IntInt64(v.(int))
	}
	switch execArea {
	case "pre":
		request.PreStrategy = []*fwmv20250611.StrategyReq{&strategyReq}
	case "post":
		request.PostStrategy = []*fwmv20250611.StrategyReq{&strategyReq}
	default:
		return fmt.Errorf("exec_area must be pre or post, got %s", execArea)
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client().CreateStrategyWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create fwm strategy failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create fwm strategy failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	log.Printf("[DEBUG]%s check StrategyIds before SetId, d.Id()=%s", logId, d.Id())
	if response.Response.StrategyIds == nil || len(response.Response.StrategyIds) == 0 || *response.Response.StrategyIds[0] == "" {
		return fmt.Errorf("Create fwm strategy failed, StrategyIds is empty. logId=%s", logId)
	}

	strategyId := *response.Response.StrategyIds[0]
	d.SetId(strategyId)
	return resourceTencentCloudFwmStrategyRead(d, meta)
}

func resourceTencentCloudFwmStrategyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_strategy.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = fwmv20250611.NewDescribeStrategyRequest()
	)

	request.StrategyId = helper.String(d.Id())

	var strategy *fwmv20250611.StrategyResp
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client().DescribeStrategyWithContext(ctx, request)
		if e != nil {
			if strings.Contains(e.Error(), "ResourceNotFound") {
				return nil
			}
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Read fwm strategy failed, Response is nil."))
		}

		strategy = result.Response.Strategy
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s read fwm strategy failed, reason:%+v", logId, err)
		return err
	}

	if strategy == nil {
		log.Printf("[CRUD] fwm_strategy id=%s", d.Id())
		d.SetId("")
		return nil
	}

	if strategy.StrategyId != nil {
		_ = d.Set("strategy_id", strategy.StrategyId)
	}

	if strategy.GroupId != nil {
		_ = d.Set("group_id", strategy.GroupId)
	}

	if strategy.GroupName != nil {
		_ = d.Set("group_name", strategy.GroupName)
	}

	if strategy.RuleCount != nil {
		_ = d.Set("rule_count", strategy.RuleCount)
	}

	if strategy.RuleStatus != nil {
		_ = d.Set("rule_status", strategy.RuleStatus)
	}

	if strategy.Sequence != nil {
		_ = d.Set("sequence", strategy.Sequence)
	}

	if strategy.ErrMsg != nil {
		_ = d.Set("err_msg", strategy.ErrMsg)
	}

	if strategy.ErrorType != nil {
		_ = d.Set("error_type", strategy.ErrorType)
	}

	if strategy.CreateTime != nil {
		_ = d.Set("create_time", strategy.CreateTime)
	}

	if strategy.UpdateTime != nil {
		_ = d.Set("update_time", strategy.UpdateTime)
	}

	if strategy.CreateBy != nil {
		_ = d.Set("create_by", strategy.CreateBy)
	}

	if strategy.CreateName != nil {
		_ = d.Set("create_name", strategy.CreateName)
	}

	if strategy.UpdateBy != nil {
		_ = d.Set("update_by", strategy.UpdateBy)
	}

	if strategy.UpdateName != nil {
		_ = d.Set("update_name", strategy.UpdateName)
	}

	if strategy.ReceiveAccount != nil {
		receiveAccountList := make([]string, 0)
		receiveGroupList := make([]string, 0)
		for _, item := range strategy.ReceiveAccount {
			if item == nil {
				continue
			}
			if item.ReceiverType != nil {
				if *item.ReceiverType == 0 {
					if item.Uin != nil {
						receiveAccountList = append(receiveAccountList, *item.Uin)
					}
				} else if *item.ReceiverType == 1 {
					if item.Uin != nil {
						receiveGroupList = append(receiveGroupList, *item.Uin)
					}
				}
			}
		}
		_ = d.Set("receive_account", receiveAccountList)
		_ = d.Set("receive_group", receiveGroupList)
	}

	return nil
}

func resourceTencentCloudFwmStrategyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_strategy.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	needChange := false
	mutableArgs := []string{"group_id", "sequence", "receive_account", "receive_group"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := fwmv20250611.NewModifyStrategyRequest()
		request.StrategyId = helper.String(d.Id())

		if d.HasChange("group_id") {
			if v, ok := d.GetOk("group_id"); ok {
				request.GroupId = helper.String(v.(string))
			}
		}

		if d.HasChange("sequence") {
			if v, ok := d.GetOk("sequence"); ok {
				request.Sequence = helper.IntInt64(v.(int))
			}
		}

		if d.HasChange("receive_account") {
			if v, ok := d.GetOk("receive_account"); ok {
				accounts := v.([]interface{})
				if len(accounts) > 0 {
					accountList := make([]string, 0, len(accounts))
					for _, account := range accounts {
						accountList = append(accountList, account.(string))
					}
					request.ReceiveAccount = helper.Strings(accountList)
				}
			}
		}

		if d.HasChange("receive_group") {
			if v, ok := d.GetOk("receive_group"); ok {
				groups := v.([]interface{})
				if len(groups) > 0 {
					groupList := make([]string, 0, len(groups))
					for _, group := range groups {
						groupList = append(groupList, group.(string))
					}
					request.ReceiveGroup = helper.Strings(groupList)
				}
			}
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client().ModifyStrategyWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Modify fwm strategy failed, Response is nil."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update fwm strategy failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudFwmStrategyRead(d, meta)
}

func resourceTencentCloudFwmStrategyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_strategy.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = fwmv20250611.NewDeleteStrategyRequest()
	)

	strategyId := d.Id()
	request.StrategyIds = []*string{helper.String(strategyId)}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client().DeleteStrategyWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete fwm strategy failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete fwm strategy failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
