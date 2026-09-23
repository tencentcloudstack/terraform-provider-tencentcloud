package config

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	configv20220802 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudConfigAggregator() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudConfigAggregatorCreate,
		Read:   resourceTencentCloudConfigAggregatorRead,
		Update: resourceTencentCloudConfigAggregatorUpdate,
		Delete: resourceTencentCloudConfigAggregatorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Aggregator name.",
			},

			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Aggregator description.",
			},

			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Aggregator type. Valid values: `RD` (global aggregator), `CUSTOM` (custom aggregator).",
			},

			"owner_uin": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Creator UIN of the aggregator. Required by Describe/Update/Delete APIs. Create does not return it.",
			},

			"aggregator_accounts": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Member account list of the aggregator, up to 100 entries.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"member_uin": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Member account ID.",
						},
						"member_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Member account name.",
						},
					},
				},
			},

			"account_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Aggregator ID.",
			},

			"aggregator_status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Aggregator creation status.",
			},
		},
	}
}

func resourceTencentCloudConfigAggregatorCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_config_aggregator.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = configv20220802.NewCreateAggregatorRequest()
		response = configv20220802.NewCreateAggregatorResponse()
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("aggregator_accounts"); ok {
		for _, item := range v.([]interface{}) {
			accountMap := item.(map[string]interface{})
			account := configv20220802.AggregatorAccount{}
			if v, ok := accountMap["member_uin"].(int); ok {
				account.MemberUin = helper.IntUint64(v)
			}

			if v, ok := accountMap["member_name"].(string); ok && v != "" {
				account.MemberName = helper.String(v)
			}

			request.AggregatorAccounts = append(request.AggregatorAccounts, &account)
		}
	}

	ownerUin := d.Get("owner_uin").(string)

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseConfigV20220802Client().CreateAggregatorWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("create config aggregator failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create config aggregator failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	log.Printf("[INFO]%s create config aggregator success, logId:%s, d.Id():%s", logId, logId, d.Id())

	if response.Response.AccountGroupId == nil || *response.Response.AccountGroupId == "" {
		return fmt.Errorf("AccountGroupId is nil or empty.")
	}

	accountGroupId := *response.Response.AccountGroupId
	d.SetId(strings.Join([]string{accountGroupId, ownerUin}, tccommon.FILED_SP))
	return resourceTencentCloudConfigAggregatorRead(d, meta)
}

func resourceTencentCloudConfigAggregatorRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_config_aggregator.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = ConfigService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	accountGroupId := idSplit[0]
	ownerUin := idSplit[1]

	respData, err := service.DescribeConfigAggregatorById(ctx, accountGroupId, ownerUin)
	if err != nil {
		return err
	}

	if respData == nil || respData.Name == nil {
		log.Printf("[CRUD] tencentcloud_config_aggregator id=%s", d.Id())
		d.SetId("")
		return nil
	}

	if respData.Name != nil {
		_ = d.Set("name", respData.Name)
	}

	if respData.Description != nil {
		_ = d.Set("description", respData.Description)
	}

	if respData.Type != nil {
		_ = d.Set("type", respData.Type)
	}

	if respData.AggregatorStatus != nil {
		_ = d.Set("aggregator_status", respData.AggregatorStatus)
	}

	_ = d.Set("account_group_id", accountGroupId)
	_ = d.Set("owner_uin", ownerUin)

	if respData.AggregatorAccounts != nil && len(respData.AggregatorAccounts) > 0 {
		aggregatorAccountsList := make([]map[string]interface{}, 0, len(respData.AggregatorAccounts))
		for _, account := range respData.AggregatorAccounts {
			accountMap := map[string]interface{}{}
			if account.MemberUin != nil {
				accountMap["member_uin"] = account.MemberUin
			}

			if account.MemberName != nil {
				accountMap["member_name"] = account.MemberName
			}

			aggregatorAccountsList = append(aggregatorAccountsList, accountMap)
		}

		_ = d.Set("aggregator_accounts", aggregatorAccountsList)
	}

	return nil
}

func resourceTencentCloudConfigAggregatorUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_config_aggregator.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	accountGroupId := idSplit[0]
	ownerUin := idSplit[1]

	needChange := false
	mutableArgs := []string{"name", "description", "aggregator_accounts"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := configv20220802.NewUpdateAggregatorRequest()
		request.AccountGroupId = &accountGroupId
		request.OwnerUin = helper.StrToUint64Point(ownerUin)

		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}

		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}

		if v, ok := d.GetOk("aggregator_accounts"); ok {
			for _, item := range v.([]interface{}) {
				accountMap := item.(map[string]interface{})
				account := configv20220802.AggregatorAccount{}
				if v, ok := accountMap["member_uin"].(int); ok {
					account.MemberUin = helper.IntUint64(v)
				}

				if v, ok := accountMap["member_name"].(string); ok && v != "" {
					account.MemberName = helper.String(v)
				}

				request.AggregatorAccounts = append(request.AggregatorAccounts, &account)
			}
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseConfigV20220802Client().UpdateAggregatorWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("update config aggregator failed, Response is nil."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update config aggregator failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudConfigAggregatorRead(d, meta)
}

func resourceTencentCloudConfigAggregatorDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_config_aggregator.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = configv20220802.NewDeleteAggregatorsRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	accountGroupId := idSplit[0]
	ownerUin := idSplit[1]

	request.AccountGroupId = &accountGroupId
	request.OwnerUin = helper.StrToUint64Point(ownerUin)

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseConfigV20220802Client().DeleteAggregatorsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("delete config aggregator failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete config aggregator failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
