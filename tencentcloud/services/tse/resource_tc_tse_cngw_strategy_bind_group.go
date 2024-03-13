package tse

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

func ResourceTencentCloudTseCngwStrategyBindGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTseCngwStrategyBindGroupCreate,
		Read:   resourceTencentCloudTseCngwStrategyBindGroupRead,
		Update: resourceTencentCloudTseCngwStrategyBindGroupUpdate,
		Delete: resourceTencentCloudTseCngwStrategyBindGroupDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				defaultValues := map[string]interface{}{
					"option": "bind",
				}

				for k, v := range defaultValues {
					_ = d.Set(k, v)
				}
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "gateway ID.",
			},

			"strategy_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "strategy ID.",
			},

			"group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "group ID.",
			},
			"option": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "`bind` or `unbind`.",
			},

			"status": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Binding status.",
			},
		},
	}
}

func resourceTencentCloudTseCngwStrategyBindGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tse_cngw_strategy_bind_group.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		gatewayId  string
		strategyId string
		groupId    string
	)
	if v, ok := d.GetOk("gateway_id"); ok {
		gatewayId = v.(string)
	}
	if v, ok := d.GetOk("strategy_id"); ok {
		strategyId = v.(string)
	}
	if v, ok := d.GetOk("group_id"); ok {
		groupId = v.(string)
	}

	d.SetId(gatewayId + tccommon.FILED_SP + strategyId + tccommon.FILED_SP + groupId)

	return resourceTencentCloudTseCngwStrategyBindGroupUpdate(d, meta)
}

func resourceTencentCloudTseCngwStrategyBindGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tse_cngw_strategy_bind_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayId := idSplit[0]
	strategyId := idSplit[1]
	groupId := idSplit[2]

	cngwStrategyBindGroup, err := service.DescribeTseCngwStrategyBindGroupById(ctx, gatewayId, strategyId, groupId)
	if err != nil {
		return err
	}

	if cngwStrategyBindGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `cngwStrategyBindGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("gateway_id", gatewayId)
	_ = d.Set("strategy_id", strategyId)
	_ = d.Set("group_id", groupId)

	if cngwStrategyBindGroup.Status != nil {
		_ = d.Set("status", cngwStrategyBindGroup.Status)
	}

	return nil
}

func resourceTencentCloudTseCngwStrategyBindGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tse_cngw_strategy_bind_group.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayId := idSplit[0]
	strategyId := idSplit[1]
	groupId := idSplit[2]

	if v, ok := d.GetOk("option"); ok {
		if v.(string) == "bind" {
			request := tse.NewBindAutoScalerResourceStrategyToGroupsRequest()
			request.GatewayId = &gatewayId
			request.StrategyId = &strategyId
			request.GroupIds = []*string{&groupId}

			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTseClient().BindAutoScalerResourceStrategyToGroups(request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s bind tse cngwStrategyBindGroup failed, reason:%+v", logId, err)
				return err
			}
		} else if v.(string) == "unbind" {
			request := tse.NewUnbindAutoScalerResourceStrategyFromGroupsRequest()
			request.GatewayId = &gatewayId
			request.StrategyId = &strategyId
			request.GroupIds = []*string{&groupId}

			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTseClient().UnbindAutoScalerResourceStrategyFromGroups(request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s unbind tse cngwStrategyBindGroup failed, reason:%+v", logId, err)
				return err
			}
		}
	} else {
		return fmt.Errorf("The option is incorrectly filled in. The optional value is bind or unbind.")
	}

	return resourceTencentCloudTseCngwStrategyBindGroupRead(d, meta)
}

func resourceTencentCloudTseCngwStrategyBindGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tse_cngw_strategy_bind_group.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
