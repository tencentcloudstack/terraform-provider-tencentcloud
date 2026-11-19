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

			"receive_account": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of receive account UINs.",
			},

			"receive_group": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of receive account group IDs.",
			},

			"pre_strategy": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"post_strategy"},
				Description:   "Pre Rule Group.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Rule group ID.",
						},
						"sequence": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Priority.",
						},
					},
				},
			},

			"post_strategy": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"pre_strategy"},
				Description:   "Post Rule Group.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Rule group ID.",
						},
						"sequence": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Priority.",
						},
					},
				},
			},

			// computed
			"strategy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Strategy ID.",
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
		accounts := v.(*schema.Set).List()
		if len(accounts) > 0 {
			accountList := make([]string, 0, len(accounts))
			for _, account := range accounts {
				accountList = append(accountList, account.(string))
			}
			request.ReceiveAccount = helper.Strings(accountList)
		}
	}

	if v, ok := d.GetOk("receive_group"); ok {
		groups := v.(*schema.Set).List()
		if len(groups) > 0 {
			groupList := make([]string, 0, len(groups))
			for _, group := range groups {
				groupList = append(groupList, group.(string))
			}
			request.ReceiveGroup = helper.Strings(groupList)
		}
	}

	if v, ok := d.GetOk("pre_strategy"); ok {
		strategies := v.([]interface{})
		if len(strategies) > 0 {
			preStrategy := strategies[0].(map[string]interface{})
			request.PreStrategy = []*fwmv20250611.StrategyReq{
				{
					GroupId:  helper.String(preStrategy["group_id"].(string)),
					Sequence: helper.IntInt64(1),
				},
			}
		}
	}

	if v, ok := d.GetOk("post_strategy"); ok {
		strategies := v.([]interface{})
		if len(strategies) > 0 {
			preStrategy := strategies[0].(map[string]interface{})
			request.PostStrategy = []*fwmv20250611.StrategyReq{
				{
					GroupId:  helper.String(preStrategy["group_id"].(string)),
					Sequence: helper.IntInt64(1),
				},
			}
		}
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
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = fwmv20250611.NewDescribeStrategyRequest()
		strategyId = d.Id()
	)

	request.StrategyId = helper.String(strategyId)
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
		log.Printf("[WARN]%s resource `tencentcloud_fwm_strategy` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if strategy.StrategyId != nil {
		_ = d.Set("strategy_id", strategy.StrategyId)
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

	if strategy.ExecArea != nil {
		tmpMap := map[string]interface{}{}
		if strategy.GroupId != nil {
			tmpMap["group_id"] = strategy.GroupId
		}

		if strategy.Sequence != nil {
			tmpMap["sequence"] = strategy.Sequence
		}

		if *strategy.ExecArea == "pre" {
			_ = d.Set("pre_strategy", []interface{}{tmpMap})
		}

		if *strategy.ExecArea == "post" {
			_ = d.Set("post_strategy", []interface{}{tmpMap})
		}
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
	mutableArgs := []string{"group_id", "receive_account", "receive_group"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := fwmv20250611.NewModifyStrategyRequest()
		request.StrategyId = helper.String(d.Id())

		if v, ok := d.GetOk("group_id"); ok {
			request.GroupId = helper.String(v.(string))
		}

		if v, ok := d.GetOk("receive_account"); ok {
			accounts := v.(*schema.Set).List()
			if len(accounts) > 0 {
				accountList := make([]string, 0, len(accounts))
				for _, account := range accounts {
					accountList = append(accountList, account.(string))
				}
				request.ReceiveAccount = helper.Strings(accountList)
			}
		}

		if v, ok := d.GetOk("receive_group"); ok {
			groups := v.(*schema.Set).List()
			if len(groups) > 0 {
				groupList := make([]string, 0, len(groups))
				for _, group := range groups {
					groupList = append(groupList, group.(string))
				}
				request.ReceiveGroup = helper.Strings(groupList)
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
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = fwmv20250611.NewDeleteStrategyRequest()
		strategyId = d.Id()
	)

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
