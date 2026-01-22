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

func ResourceTencentCloudDlcStandardEngineResourceGroupConfigInfo() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDlcStandardEngineResourceGroupConfigInfoCreate,
		Read:   resourceTencentCloudDlcStandardEngineResourceGroupConfigInfoRead,
		Update: resourceTencentCloudDlcStandardEngineResourceGroupConfigInfoUpdate,
		Delete: resourceTencentCloudDlcStandardEngineResourceGroupConfigInfoDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"engine_resource_group_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Standard engine resource group name.",
			},

			"static_conf_context": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Static config context.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"params": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Collection of bound working groups.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"config_item": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Configuration item.",
									},
									"config_value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Configuration value.",
									},
								},
							},
						},
					},
				},
			},

			"dynamic_conf_context": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Dynamic config context.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"params": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Collection of bound working groups.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"config_item": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Configuration item.",
									},
									"config_value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Configuration value.",
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

func resourceTencentCloudDlcStandardEngineResourceGroupConfigInfoCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_standard_engine_resource_group_config_info.create")()
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

	confContextList := []*dlcv20210125.UpdateConfContext{}
	if v, ok := d.GetOk("static_conf_context"); ok {
		confContext := dlcv20210125.UpdateConfContext{}
		confContext.ConfigType = helper.String("StaticConfigType")
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			if v, ok := dMap["params"]; ok {
				for _, item := range v.(*schema.Set).List() {
					param := dlcv20210125.Param{}
					dMap := item.(map[string]interface{})
					if v, ok := dMap["config_item"].(string); ok && v != "" {
						param.ConfigItem = helper.String(v)
					}

					if v, ok := dMap["config_value"].(string); ok && v != "" {
						param.ConfigValue = helper.String(v)
					}

					param.Operate = helper.String("ADD")
					confContext.Params = append(confContext.Params, &param)
				}
			}
		}

		confContextList = append(confContextList, &confContext)
	}

	if v, ok := d.GetOk("dynamic_conf_context"); ok {
		confContext := dlcv20210125.UpdateConfContext{}
		confContext.ConfigType = helper.String("DynamicConfigType")
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			if v, ok := dMap["params"]; ok {
				for _, item := range v.(*schema.Set).List() {
					param := dlcv20210125.Param{}
					dMap := item.(map[string]interface{})
					if v, ok := dMap["config_item"].(string); ok && v != "" {
						param.ConfigItem = helper.String(v)
					}

					if v, ok := dMap["config_value"].(string); ok && v != "" {
						param.ConfigValue = helper.String(v)
					}

					param.Operate = helper.String("ADD")
					confContext.Params = append(confContext.Params, &param)
				}
			}
		}

		confContextList = append(confContextList, &confContext)
	}

	request.UpdateConfContext = confContextList
	request.IsEffectiveNow = helper.Int64(0)
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().UpdateStandardEngineResourceGroupConfigInfoWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create dlc standard engine resource group config info failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create dlc standard engine resource group config info failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(engineResourceGroupName)

	// wait
	waitErr := resource.Retry(tccommon.WriteRetryTimeout*4, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().DescribeStandardEngineResourceGroupsWithContext(ctx, &dlcv20210125.DescribeStandardEngineResourceGroupsRequest{
			Filters: []*dlcv20210125.Filter{
				{
					Name:   helper.String("engine-resource-group-name-unique"),
					Values: helper.Strings([]string{engineResourceGroupName}),
				},
			},
		})
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe dlc standard engine resource groups failed, Response is nil."))
		}

		if result.Response.UserEngineResourceGroupInfos == nil || len(result.Response.UserEngineResourceGroupInfos) == 0 {
			return resource.NonRetryableError(fmt.Errorf("UserEngineResourceGroupInfos is nil."))
		}

		if len(result.Response.UserEngineResourceGroupInfos) != 1 {
			return resource.NonRetryableError(fmt.Errorf("UserEngineResourceGroupInfos is not 1."))
		}

		state := result.Response.UserEngineResourceGroupInfos[0].ResourceGroupState
		if state != nil {
			if *state == 2 {
				return nil
			}
		} else {
			return resource.NonRetryableError(fmt.Errorf("ResourceGroupState is nil."))
		}

		return resource.RetryableError(fmt.Errorf("UserEngineResourceGroupInfos is not ready, state:%d", *state))
	})

	if waitErr != nil {
		log.Printf("[CRITAL]%s wait for dlc standard engine resource group failed, reason:%+v", logId, waitErr)
		return waitErr
	}

	return resourceTencentCloudDlcStandardEngineResourceGroupConfigInfoRead(d, meta)
}

func resourceTencentCloudDlcStandardEngineResourceGroupConfigInfoRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_standard_engine_resource_group_config_info.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId                   = tccommon.GetLogId(tccommon.ContextNil)
		ctx                     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service                 = DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		engineResourceGroupName = d.Id()
	)

	respData, err := service.DescribeDlcStandardEngineResourceGroupConfigInfoById(ctx, engineResourceGroupName)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_dlc_standard_engine_resource_group_config_info` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("engine_resource_group_name", engineResourceGroupName)

	if respData.StaticConfigPairs != nil {
		dMap := make(map[string]interface{}, 0)
		tmpList := make([]map[string]interface{}, 0, 1)
		for _, configPair := range respData.StaticConfigPairs {
			configPairMap := map[string]interface{}{}
			if configPair.ConfigItem != nil {
				configPairMap["config_item"] = configPair.ConfigItem
			}

			if configPair.ConfigValue != nil {
				configPairMap["config_value"] = configPair.ConfigValue
			}

			tmpList = append(tmpList, configPairMap)
		}

		dMap["params"] = tmpList
		_ = d.Set("static_conf_context", []interface{}{dMap})
	}

	if respData.DynamicConfigPairs != nil {
		dMap := make(map[string]interface{}, 0)
		tmpList := make([]map[string]interface{}, 0, 1)
		for _, configPair := range respData.DynamicConfigPairs {
			configPairMap := map[string]interface{}{}
			if configPair.ConfigItem != nil {
				configPairMap["config_item"] = configPair.ConfigItem
			}

			if configPair.ConfigValue != nil {
				configPairMap["config_value"] = configPair.ConfigValue
			}

			tmpList = append(tmpList, configPairMap)
		}

		dMap["params"] = tmpList
		_ = d.Set("dynamic_conf_context", []interface{}{dMap})
	}

	return nil
}

func resourceTencentCloudDlcStandardEngineResourceGroupConfigInfoUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_standard_engine_resource_group_config_info.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId                   = tccommon.GetLogId(tccommon.ContextNil)
		ctx                     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		engineResourceGroupName = d.Id()
	)

	if d.HasChange("static_conf_context") {
		oldInterface, newInterface := d.GetChange("static_conf_context.0.params")
		olds := oldInterface.(*schema.Set)
		news := newInterface.(*schema.Set)
		remove := olds.Difference(news).List()
		add := news.Difference(olds).List()

		if len(remove) > 0 {
			request := dlcv20210125.NewUpdateStandardEngineResourceGroupConfigInfoRequest()
			confContext := dlcv20210125.UpdateConfContext{}
			confContext.ConfigType = helper.String("StaticConfigType")
			for _, item := range remove {
				dMap := item.(map[string]interface{})
				param := dlcv20210125.Param{}
				if v, ok := dMap["config_item"].(string); ok && v != "" {
					param.ConfigItem = helper.String(v)
				}

				if v, ok := dMap["config_value"].(string); ok && v != "" {
					param.ConfigValue = helper.String(v)
				}

				param.Operate = helper.String("DELETE")
				confContext.Params = append(confContext.Params, &param)
			}

			request.UpdateConfContext = append(request.UpdateConfContext, &confContext)
			request.EngineResourceGroupName = &engineResourceGroupName
			request.IsEffectiveNow = helper.Int64(0)
			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().UpdateStandardEngineResourceGroupConfigInfoWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				if result == nil || result.Response == nil {
					return resource.NonRetryableError(fmt.Errorf("Update dlc standard engine resource group config info failed, Response is nil."))
				}

				return nil
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s update dlc standard engine resource group config info failed, reason:%+v", logId, reqErr)
				return reqErr
			}

			// wait
			waitErr := resource.Retry(tccommon.WriteRetryTimeout*4, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().DescribeStandardEngineResourceGroupsWithContext(ctx, &dlcv20210125.DescribeStandardEngineResourceGroupsRequest{
					Filters: []*dlcv20210125.Filter{
						{
							Name:   helper.String("engine-resource-group-name-unique"),
							Values: helper.Strings([]string{engineResourceGroupName}),
						},
					},
				})
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				if result == nil || result.Response == nil {
					return resource.NonRetryableError(fmt.Errorf("Describe dlc standard engine resource groups failed, Response is nil."))
				}

				if result.Response.UserEngineResourceGroupInfos == nil || len(result.Response.UserEngineResourceGroupInfos) == 0 {
					return resource.NonRetryableError(fmt.Errorf("UserEngineResourceGroupInfos is nil."))
				}

				if len(result.Response.UserEngineResourceGroupInfos) != 1 {
					return resource.NonRetryableError(fmt.Errorf("UserEngineResourceGroupInfos is not 1."))
				}

				state := result.Response.UserEngineResourceGroupInfos[0].ResourceGroupState
				if state != nil {
					if *state == 2 {
						return nil
					}
				} else {
					return resource.NonRetryableError(fmt.Errorf("ResourceGroupState is nil."))
				}

				return resource.RetryableError(fmt.Errorf("UserEngineResourceGroupInfos is not ready, state:%d", *state))
			})

			if waitErr != nil {
				log.Printf("[CRITAL]%s wait for dlc standard engine resource group failed, reason:%+v", logId, waitErr)
				return waitErr
			}
		}

		if len(add) > 0 {
			request := dlcv20210125.NewUpdateStandardEngineResourceGroupConfigInfoRequest()
			confContext := dlcv20210125.UpdateConfContext{}
			confContext.ConfigType = helper.String("StaticConfigType")
			for _, item := range add {
				dMap := item.(map[string]interface{})
				param := dlcv20210125.Param{}
				if v, ok := dMap["config_item"].(string); ok && v != "" {
					param.ConfigItem = helper.String(v)
				}

				if v, ok := dMap["config_value"].(string); ok && v != "" {
					param.ConfigValue = helper.String(v)
				}

				param.Operate = helper.String("ADD")
				confContext.Params = append(confContext.Params, &param)
			}

			request.UpdateConfContext = append(request.UpdateConfContext, &confContext)
			request.EngineResourceGroupName = &engineResourceGroupName
			request.IsEffectiveNow = helper.Int64(0)
			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().UpdateStandardEngineResourceGroupConfigInfoWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				if result == nil || result.Response == nil {
					return resource.NonRetryableError(fmt.Errorf("Update dlc standard engine resource group config info failed, Response is nil."))
				}

				return nil
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s update dlc standard engine resource group config info failed, reason:%+v", logId, reqErr)
				return reqErr
			}

			// wait
			waitErr := resource.Retry(tccommon.WriteRetryTimeout*4, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().DescribeStandardEngineResourceGroupsWithContext(ctx, &dlcv20210125.DescribeStandardEngineResourceGroupsRequest{
					Filters: []*dlcv20210125.Filter{
						{
							Name:   helper.String("engine-resource-group-name-unique"),
							Values: helper.Strings([]string{engineResourceGroupName}),
						},
					},
				})
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				if result == nil || result.Response == nil {
					return resource.NonRetryableError(fmt.Errorf("Describe dlc standard engine resource groups failed, Response is nil."))
				}

				if result.Response.UserEngineResourceGroupInfos == nil || len(result.Response.UserEngineResourceGroupInfos) == 0 {
					return resource.NonRetryableError(fmt.Errorf("UserEngineResourceGroupInfos is nil."))
				}

				if len(result.Response.UserEngineResourceGroupInfos) != 1 {
					return resource.NonRetryableError(fmt.Errorf("UserEngineResourceGroupInfos is not 1."))
				}

				state := result.Response.UserEngineResourceGroupInfos[0].ResourceGroupState
				if state != nil {
					if *state == 2 {
						return nil
					}
				} else {
					return resource.NonRetryableError(fmt.Errorf("ResourceGroupState is nil."))
				}

				return resource.RetryableError(fmt.Errorf("UserEngineResourceGroupInfos is not ready, state:%d", *state))
			})

			if waitErr != nil {
				log.Printf("[CRITAL]%s wait for dlc standard engine resource group failed, reason:%+v", logId, waitErr)
				return waitErr
			}
		}
	}

	if d.HasChange("dynamic_conf_context") {
		oldInterface, newInterface := d.GetChange("dynamic_conf_context.0.params")
		olds := oldInterface.(*schema.Set)
		news := newInterface.(*schema.Set)
		remove := olds.Difference(news).List()
		add := news.Difference(olds).List()

		if len(remove) > 0 {
			request := dlcv20210125.NewUpdateStandardEngineResourceGroupConfigInfoRequest()
			confContext := dlcv20210125.UpdateConfContext{}
			confContext.ConfigType = helper.String("DynamicConfigType")
			for _, item := range remove {
				dMap := item.(map[string]interface{})
				param := dlcv20210125.Param{}
				if v, ok := dMap["config_item"].(string); ok && v != "" {
					param.ConfigItem = helper.String(v)
				}

				if v, ok := dMap["config_value"].(string); ok && v != "" {
					param.ConfigValue = helper.String(v)
				}

				param.Operate = helper.String("DELETE")
				confContext.Params = append(confContext.Params, &param)
			}

			request.UpdateConfContext = append(request.UpdateConfContext, &confContext)
			request.EngineResourceGroupName = &engineResourceGroupName
			request.IsEffectiveNow = helper.Int64(0)
			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().UpdateStandardEngineResourceGroupConfigInfoWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				if result == nil || result.Response == nil {
					return resource.NonRetryableError(fmt.Errorf("Update dlc standard engine resource group config info failed, Response is nil."))
				}

				return nil
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s update dlc standard engine resource group config info failed, reason:%+v", logId, reqErr)
				return reqErr
			}

			// wait
			waitErr := resource.Retry(tccommon.WriteRetryTimeout*4, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().DescribeStandardEngineResourceGroupsWithContext(ctx, &dlcv20210125.DescribeStandardEngineResourceGroupsRequest{
					Filters: []*dlcv20210125.Filter{
						{
							Name:   helper.String("engine-resource-group-name-unique"),
							Values: helper.Strings([]string{engineResourceGroupName}),
						},
					},
				})
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				if result == nil || result.Response == nil {
					return resource.NonRetryableError(fmt.Errorf("Describe dlc standard engine resource groups failed, Response is nil."))
				}

				if result.Response.UserEngineResourceGroupInfos == nil || len(result.Response.UserEngineResourceGroupInfos) == 0 {
					return resource.NonRetryableError(fmt.Errorf("UserEngineResourceGroupInfos is nil."))
				}

				if len(result.Response.UserEngineResourceGroupInfos) != 1 {
					return resource.NonRetryableError(fmt.Errorf("UserEngineResourceGroupInfos is not 1."))
				}

				state := result.Response.UserEngineResourceGroupInfos[0].ResourceGroupState
				if state != nil {
					if *state == 2 {
						return nil
					}
				} else {
					return resource.NonRetryableError(fmt.Errorf("ResourceGroupState is nil."))
				}

				return resource.RetryableError(fmt.Errorf("UserEngineResourceGroupInfos is not ready, state:%d", *state))
			})

			if waitErr != nil {
				log.Printf("[CRITAL]%s wait for dlc standard engine resource group failed, reason:%+v", logId, waitErr)
				return waitErr
			}
		}

		if len(add) > 0 {
			request := dlcv20210125.NewUpdateStandardEngineResourceGroupConfigInfoRequest()
			confContext := dlcv20210125.UpdateConfContext{}
			confContext.ConfigType = helper.String("DynamicConfigType")
			for _, item := range add {
				dMap := item.(map[string]interface{})
				param := dlcv20210125.Param{}
				if v, ok := dMap["config_item"].(string); ok && v != "" {
					param.ConfigItem = helper.String(v)
				}

				if v, ok := dMap["config_value"].(string); ok && v != "" {
					param.ConfigValue = helper.String(v)
				}

				param.Operate = helper.String("ADD")
				confContext.Params = append(confContext.Params, &param)
			}

			request.UpdateConfContext = append(request.UpdateConfContext, &confContext)
			request.EngineResourceGroupName = &engineResourceGroupName
			request.IsEffectiveNow = helper.Int64(0)
			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().UpdateStandardEngineResourceGroupConfigInfoWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				if result == nil || result.Response == nil {
					return resource.NonRetryableError(fmt.Errorf("Update dlc standard engine resource group config info failed, Response is nil."))
				}

				return nil
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s update dlc standard engine resource group config info failed, reason:%+v", logId, reqErr)
				return reqErr
			}

			// wait
			waitErr := resource.Retry(tccommon.WriteRetryTimeout*4, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().DescribeStandardEngineResourceGroupsWithContext(ctx, &dlcv20210125.DescribeStandardEngineResourceGroupsRequest{
					Filters: []*dlcv20210125.Filter{
						{
							Name:   helper.String("engine-resource-group-name-unique"),
							Values: helper.Strings([]string{engineResourceGroupName}),
						},
					},
				})
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				if result == nil || result.Response == nil {
					return resource.NonRetryableError(fmt.Errorf("Describe dlc standard engine resource groups failed, Response is nil."))
				}

				if result.Response.UserEngineResourceGroupInfos == nil || len(result.Response.UserEngineResourceGroupInfos) == 0 {
					return resource.NonRetryableError(fmt.Errorf("UserEngineResourceGroupInfos is nil."))
				}

				if len(result.Response.UserEngineResourceGroupInfos) != 1 {
					return resource.NonRetryableError(fmt.Errorf("UserEngineResourceGroupInfos is not 1."))
				}

				state := result.Response.UserEngineResourceGroupInfos[0].ResourceGroupState
				if state != nil {
					if *state == 2 {
						return nil
					}
				} else {
					return resource.NonRetryableError(fmt.Errorf("ResourceGroupState is nil."))
				}

				return resource.RetryableError(fmt.Errorf("UserEngineResourceGroupInfos is not ready, state:%d", *state))
			})

			if waitErr != nil {
				log.Printf("[CRITAL]%s wait for dlc standard engine resource group failed, reason:%+v", logId, waitErr)
				return waitErr
			}
		}
	}

	return resourceTencentCloudDlcStandardEngineResourceGroupConfigInfoRead(d, meta)
}

func resourceTencentCloudDlcStandardEngineResourceGroupConfigInfoDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_standard_engine_resource_group_config_info.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId                   = tccommon.GetLogId(tccommon.ContextNil)
		ctx                     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service                 = DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request                 = dlcv20210125.NewUpdateStandardEngineResourceGroupConfigInfoRequest()
		engineResourceGroupName = d.Id()
	)

	// get all conf context
	respData, err := service.DescribeDlcStandardEngineResourceGroupConfigInfoById(ctx, engineResourceGroupName)
	if err != nil {
		return err
	}

	if respData == nil {
		return nil
	}

	// delete all
	if respData.StaticConfigPairs != nil && len(respData.StaticConfigPairs) > 0 {
		confContext := dlcv20210125.UpdateConfContext{}
		confContext.ConfigType = helper.String("StaticConfigType")
		for _, item := range respData.StaticConfigPairs {
			param := dlcv20210125.Param{}
			if item.ConfigItem != nil {
				param.ConfigItem = item.ConfigItem
			}

			if item.ConfigValue != nil {
				param.ConfigValue = item.ConfigValue
			}

			param.Operate = helper.String("DELETE")
			confContext.Params = append(confContext.Params, &param)
		}

		request.UpdateConfContext = append(request.UpdateConfContext, &confContext)
	}

	if respData.DynamicConfigPairs != nil && len(respData.DynamicConfigPairs) > 0 {
		confContext := dlcv20210125.UpdateConfContext{}
		confContext.ConfigType = helper.String("DynamicConfigType")
		for _, item := range respData.DynamicConfigPairs {
			param := dlcv20210125.Param{}
			if item.ConfigItem != nil {
				param.ConfigItem = item.ConfigItem
			}

			if item.ConfigValue != nil {
				param.ConfigValue = item.ConfigValue
			}

			param.Operate = helper.String("DELETE")
			confContext.Params = append(confContext.Params, &param)
		}

		request.UpdateConfContext = append(request.UpdateConfContext, &confContext)
	}

	request.EngineResourceGroupName = &engineResourceGroupName
	request.IsEffectiveNow = helper.Int64(0)
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().UpdateStandardEngineResourceGroupConfigInfoWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete dlc standard engine resource group config info failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete dlc standard engine resource group config info failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	// wait
	waitErr := resource.Retry(tccommon.WriteRetryTimeout*4, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().DescribeStandardEngineResourceGroupsWithContext(ctx, &dlcv20210125.DescribeStandardEngineResourceGroupsRequest{
			Filters: []*dlcv20210125.Filter{
				{
					Name:   helper.String("engine-resource-group-name-unique"),
					Values: helper.Strings([]string{engineResourceGroupName}),
				},
			},
		})
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe dlc standard engine resource groups failed, Response is nil."))
		}

		if result.Response.UserEngineResourceGroupInfos == nil || len(result.Response.UserEngineResourceGroupInfos) == 0 {
			return resource.NonRetryableError(fmt.Errorf("UserEngineResourceGroupInfos is nil."))
		}

		if len(result.Response.UserEngineResourceGroupInfos) != 1 {
			return resource.NonRetryableError(fmt.Errorf("UserEngineResourceGroupInfos is not 1."))
		}

		state := result.Response.UserEngineResourceGroupInfos[0].ResourceGroupState
		if state != nil {
			if *state == 2 {
				return nil
			}
		} else {
			return resource.NonRetryableError(fmt.Errorf("ResourceGroupState is nil."))
		}

		return resource.RetryableError(fmt.Errorf("UserEngineResourceGroupInfos is not ready, state:%d", *state))
	})

	if waitErr != nil {
		log.Printf("[CRITAL]%s wait for dlc standard engine resource group failed, reason:%+v", logId, waitErr)
		return waitErr
	}

	return nil
}
