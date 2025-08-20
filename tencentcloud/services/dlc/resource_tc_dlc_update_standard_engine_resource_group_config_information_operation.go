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

func ResourceTencentCloudDlcUpdateStandardEngineResourceGroupConfigInformationOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDlcUpdateStandardEngineResourceGroupConfigInformationOperationCreate,
		Read:   resourceTencentCloudDlcUpdateStandardEngineResourceGroupConfigInformationOperationRead,
		Delete: resourceTencentCloudDlcUpdateStandardEngineResourceGroupConfigInformationOperationDelete,
		Schema: map[string]*schema.Schema{
			"engine_resource_group_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Engine resource group name.",
			},

			"update_conf_context": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "Configuration that needs to be updated.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_type": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Parameter type, optional: StaticConfigType, DynamicConfigType.",
						},
						"params": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Configuration array of parameters.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"config_item": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Parameter key.",
									},
									"config_value": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Parameter value.",
									},
									"operate": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Send operations, support: ADD, DELETE, MODIFY.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudDlcUpdateStandardEngineResourceGroupConfigInformationOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_update_standard_engine_resource_group_config_information_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId                   = tccommon.GetLogId(tccommon.ContextNil)
		ctx                     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request                 = dlcv20210125.NewUpdateStandardEngineResourceGroupConfigInfoRequest()
		engineResourceGroupName string
	)

	if v, ok := d.GetOk("engine_resource_group_name"); ok {
		request.EngineResourceGroupName = helper.String(v.(string))
		engineResourceGroupName = v.(string)
	}

	if v, ok := d.GetOk("update_conf_context"); ok {
		for _, item := range v.([]interface{}) {
			updateConfContextMap := item.(map[string]interface{})
			updateConfContext := dlcv20210125.UpdateConfContext{}
			if v, ok := updateConfContextMap["config_type"].(string); ok && v != "" {
				updateConfContext.ConfigType = helper.String(v)
			}

			if v, ok := updateConfContextMap["params"]; ok {
				for _, item := range v.([]interface{}) {
					paramsMap := item.(map[string]interface{})
					param := dlcv20210125.Param{}
					if v, ok := paramsMap["config_item"].(string); ok && v != "" {
						param.ConfigItem = helper.String(v)
					}

					if v, ok := paramsMap["config_value"].(string); ok && v != "" {
						param.ConfigValue = helper.String(v)
					}

					if v, ok := paramsMap["operate"].(string); ok && v != "" {
						param.Operate = helper.String(v)
					}

					updateConfContext.Params = append(updateConfContext.Params, &param)
				}
			}

			request.UpdateConfContext = append(request.UpdateConfContext, &updateConfContext)
		}
	}

	request.IsEffectiveNow = helper.IntInt64(0)
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().UpdateStandardEngineResourceGroupConfigInfoWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create dlc update standard engine resource group config information operation failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(engineResourceGroupName)

	// wait
	waitReq := dlcv20210125.NewDescribeStandardEngineResourceGroupsRequest()
	waitReq.Filters = []*dlcv20210125.Filter{
		&dlcv20210125.Filter{
			Name:   helper.String("engine-resource-group-name-unique"),
			Values: helper.Strings([]string{engineResourceGroupName}),
		},
	}
	reqErr = resource.Retry(tccommon.ReadRetryTimeout*7, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().DescribeStandardEngineResourceGroupsWithContext(ctx, waitReq)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe standard engine resource groups failed, Response is nil."))
		}

		if result.Response.UserEngineResourceGroupInfos == nil || len(result.Response.UserEngineResourceGroupInfos) == 0 {
			return resource.NonRetryableError(fmt.Errorf("UserEngineResourceGroupInfos is nil."))
		}

		if result.Response.UserEngineResourceGroupInfos[0].ResourceGroupState == nil {
			return resource.NonRetryableError(fmt.Errorf("ResourceGroupState is nil."))
		}

		resourceGroupState := result.Response.UserEngineResourceGroupInfos[0].ResourceGroupState
		if *resourceGroupState == 2 {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("ResourceGroupState is not running, current state is %d...", *resourceGroupState))
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s describe standard engine resource groups failed,, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudDlcUpdateStandardEngineResourceGroupConfigInformationOperationRead(d, meta)
}

func resourceTencentCloudDlcUpdateStandardEngineResourceGroupConfigInformationOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_update_standard_engine_resource_group_config_information_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDlcUpdateStandardEngineResourceGroupConfigInformationOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_update_standard_engine_resource_group_config_information_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
