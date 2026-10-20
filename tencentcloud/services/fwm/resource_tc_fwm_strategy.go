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
				Description: "Product type.",
			},

			"receive_account": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Receive account list.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"pre_strategy": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Pre strategy rule group set.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Rule group ID.",
						},
						"sequence": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Priority.",
						},
					},
				},
			},

			"post_strategy": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Post strategy rule group set.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Rule group ID.",
						},
						"sequence": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Priority.",
						},
					},
				},
			},

			"receive_group": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Receive account group list.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"strategy_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Strategy ID list.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
		for _, item := range v.([]interface{}) {
			request.ReceiveAccount = append(request.ReceiveAccount, helper.String(item.(string)))
		}
	}

	if v, ok := d.GetOk("pre_strategy"); ok {
		for _, item := range v.([]interface{}) {
			preStrategyMap := item.(map[string]interface{})
			strategyReq := fwmv20250611.StrategyReq{}
			if v, ok := preStrategyMap["group_id"].(string); ok && v != "" {
				strategyReq.GroupId = helper.String(v)
			}
			if v, ok := preStrategyMap["sequence"].(int); ok {
				strategyReq.Sequence = helper.IntInt64(v)
			}
			request.PreStrategy = append(request.PreStrategy, &strategyReq)
		}
	}

	if v, ok := d.GetOk("post_strategy"); ok {
		for _, item := range v.([]interface{}) {
			postStrategyMap := item.(map[string]interface{})
			strategyReq := fwmv20250611.StrategyReq{}
			if v, ok := postStrategyMap["group_id"].(string); ok && v != "" {
				strategyReq.GroupId = helper.String(v)
			}
			if v, ok := postStrategyMap["sequence"].(int); ok {
				strategyReq.Sequence = helper.IntInt64(v)
			}
			request.PostStrategy = append(request.PostStrategy, &strategyReq)
		}
	}

	if v, ok := d.GetOk("receive_group"); ok {
		for _, item := range v.([]interface{}) {
			request.ReceiveGroup = append(request.ReceiveGroup, helper.String(item.(string)))
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
			return resource.NonRetryableError(fmt.Errorf("Create fwm_strategy failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create fwm_strategy failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	log.Printf("[DEBUG]%s create fwm_strategy id=%s", logId, d.Id())
	if response.Response.StrategyIds == nil || len(response.Response.StrategyIds) == 0 {
		return fmt.Errorf("Create fwm_strategy failed, StrategyIds is nil/empty.")
	}

	strategyIds := make([]string, 0, len(response.Response.StrategyIds))
	for _, id := range response.Response.StrategyIds {
		if id != nil {
			strategyIds = append(strategyIds, *id)
		}
	}
	if len(strategyIds) == 0 {
		return fmt.Errorf("Create fwm_strategy failed, StrategyIds is nil/empty.")
	}

	d.SetId(strings.Join(strategyIds, tccommon.FILED_SP))
	_ = d.Set("strategy_ids", helper.StringsInterfaces(response.Response.StrategyIds))
	return resourceTencentCloudFwmStrategyRead(d, meta)
}

func resourceTencentCloudFwmStrategyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_strategy.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) < 1 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	strategyIds := idSplit
	foundCount := 0
	for _, strategyId := range strategyIds {
		request := fwmv20250611.NewDescribeStrategyRequest()
		request.StrategyId = helper.String(strategyId)
		var strategyResp *fwmv20250611.StrategyResp
		reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client().DescribeStrategyWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe fwm_strategy failed, Response is nil."))
			}

			strategyResp = result.Response.Strategy
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s read fwm_strategy failed, reason:%+v", logId, reqErr)
			return reqErr
		}

		if strategyResp != nil {
			foundCount++
		}
	}

	if foundCount == 0 {
		log.Printf("[CRUD] fwm_strategy id=%s", d.Id())
		d.SetId("")
		return nil
	}

	return nil
}

func resourceTencentCloudFwmStrategyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_strategy.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	immutableArgs := []string{"receive_account", "pre_strategy", "post_strategy", "receive_group"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("fwm_strategy %s is immutable, please recreate resource", v)
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

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) < 1 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	for _, strategyId := range idSplit {
		request.StrategyIds = append(request.StrategyIds, helper.String(strategyId))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client().DeleteStrategyWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete fwm_strategy failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete fwm_strategy failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
