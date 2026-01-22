package dlc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlcv20210125 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDlcDataMaskStrategy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDlcDataMaskStrategyCreate,
		Read:   resourceTencentCloudDlcDataMaskStrategyRead,
		Update: resourceTencentCloudDlcDataMaskStrategyUpdate,
		Delete: resourceTencentCloudDlcDataMaskStrategyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"strategy": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "The data masking strategy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"strategy_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the data masking strategy.",
						},
						"strategy_desc": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The description of the data masking strategy.",
						},
						"groups": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Collection of bound working groups.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"work_group_id": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The unique ID of the work group.",
									},
									"strategy_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The type of the data masking strategy. Supported value: MASK/MASK_NONE/MASK_NULL/MASK_HASH/MASK_SHOW_LAST_4/MASK_SHOW_FIRST_4/MASK_DATE_SHOW_YEAR.",
									},
								},
							},
						},
						"strategy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the data masking strategy.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudDlcDataMaskStrategyCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_data_mask_strategy.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = dlcv20210125.NewCreateDataMaskStrategyRequest()
		response   = dlcv20210125.NewCreateDataMaskStrategyResponse()
		strategyId string
	)

	if strategyMap, ok := helper.InterfacesHeadMap(d, "strategy"); ok {
		dataMaskStrategyInfo := dlcv20210125.DataMaskStrategyInfo{}
		if v, ok := strategyMap["strategy_name"].(string); ok && v != "" {
			dataMaskStrategyInfo.StrategyName = helper.String(v)
		}

		if v, ok := strategyMap["strategy_desc"].(string); ok && v != "" {
			dataMaskStrategyInfo.StrategyDesc = helper.String(v)
		}

		if v, ok := strategyMap["groups"]; ok {
			for _, item := range v.([]interface{}) {
				groupsMap := item.(map[string]interface{})
				groupInfo := dlcv20210125.GroupInfo{}
				if v, ok := groupsMap["work_group_id"].(int); ok {
					groupInfo.WorkGroupId = helper.IntInt64(v)
				}

				if v, ok := groupsMap["strategy_type"].(string); ok && v != "" {
					groupInfo.StrategyType = helper.String(v)
				}

				dataMaskStrategyInfo.Groups = append(dataMaskStrategyInfo.Groups, &groupInfo)
			}
		}

		request.Strategy = &dataMaskStrategyInfo
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().CreateDataMaskStrategyWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create dlc data mask strategy failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create dlc data mask strategy failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.StrategyId == nil {
		return fmt.Errorf("StrategyId is nil.")
	}

	strategyId = *response.Response.StrategyId
	d.SetId(strategyId)
	return resourceTencentCloudDlcDataMaskStrategyRead(d, meta)
}

func resourceTencentCloudDlcDataMaskStrategyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_data_mask_strategy.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service    = DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		strategyId = d.Id()
	)

	respData, err := service.DescribeDlcDataMaskStrategyById(ctx, strategyId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_dlc_data_mask_strategy` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	strategiesList := make([]map[string]interface{}, 0, 1)
	strategiesMap := map[string]interface{}{}
	if respData.StrategyName != nil {
		strategiesMap["strategy_name"] = respData.StrategyName
	}

	if respData.StrategyDesc != nil {
		strategiesMap["strategy_desc"] = respData.StrategyDesc
	}

	if respData.Groups != nil {
		groupsList := make([]map[string]interface{}, 0, len(respData.Groups))
		for _, groups := range respData.Groups {
			groupsMap := map[string]interface{}{}
			if groups.WorkGroupId != nil {
				groupsMap["work_group_id"] = groups.WorkGroupId
			}

			if groups.StrategyType != nil {
				groupsMap["strategy_type"] = groups.StrategyType
			}

			groupsList = append(groupsList, groupsMap)
		}

		strategiesMap["groups"] = groupsList
	}

	if respData.StrategyId != nil {
		strategiesMap["strategy_id"] = respData.StrategyId
	}

	strategiesList = append(strategiesList, strategiesMap)
	_ = d.Set("strategy", strategiesList)

	return nil
}

func resourceTencentCloudDlcDataMaskStrategyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_data_mask_strategy.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		strategyId = d.Id()
	)

	needChange := false
	mutableArgs := []string{"strategy"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := dlcv20210125.NewUpdateDataMaskStrategyRequest()
		if strategyMap, ok := helper.InterfacesHeadMap(d, "strategy"); ok {
			dataMaskStrategyInfo := dlcv20210125.DataMaskStrategyInfo{}
			if v, ok := strategyMap["strategy_name"].(string); ok && v != "" {
				dataMaskStrategyInfo.StrategyName = helper.String(v)
			}

			if v, ok := strategyMap["strategy_desc"].(string); ok && v != "" {
				dataMaskStrategyInfo.StrategyDesc = helper.String(v)
			}

			if v, ok := strategyMap["groups"]; ok {
				for _, item := range v.([]interface{}) {
					groupsMap := item.(map[string]interface{})
					groupInfo := dlcv20210125.GroupInfo{}
					if v, ok := groupsMap["work_group_id"].(int); ok {
						groupInfo.WorkGroupId = helper.IntInt64(v)
					}

					if v, ok := groupsMap["strategy_type"].(string); ok && v != "" {
						groupInfo.StrategyType = helper.String(v)
					}

					dataMaskStrategyInfo.Groups = append(dataMaskStrategyInfo.Groups, &groupInfo)
				}
			}

			dataMaskStrategyInfo.StrategyId = &strategyId
			request.Strategy = &dataMaskStrategyInfo
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().UpdateDataMaskStrategyWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update dlc data mask strategy failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudDlcDataMaskStrategyRead(d, meta)
}

func resourceTencentCloudDlcDataMaskStrategyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_data_mask_strategy.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = dlcv20210125.NewDeleteDataMaskStrategyRequest()
		strategyId = d.Id()
	)

	request.StrategyId = &strategyId
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().DeleteDataMaskStrategyWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete dlc data mask strategy failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
