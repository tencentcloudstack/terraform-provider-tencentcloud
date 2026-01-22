package cdwpg

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdwpgv20201230 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwpg/v20201230"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCdwpgDbconfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCdwpgDbconfigCreate,
		Read:   resourceTencentCloudCdwpgDbconfigRead,
		Update: resourceTencentCloudCdwpgDbconfigUpdate,
		Delete: resourceTencentCloudCdwpgDbconfigDelete,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance id.",
			},
			"node_config_params": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Node config parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Node type.",
						},
						"parameter_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Parameter name.",
						},
						"parameter_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Parameter value.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCdwpgDbconfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdwpg_dbconfig.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	instanceId := d.Get("instance_id").(string)
	d.SetId(instanceId)

	return resourceTencentCloudCdwpgDbconfigUpdate(d, meta)
}

func resourceTencentCloudCdwpgDbconfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdwpg_dbconfig.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	instanceId := d.Id()

	var respData []*cdwpgv20201230.ParamItem
	paramMap := make(map[string]interface{})
	paramMap["InstanceId"] = helper.String(instanceId)
	paramMap["NodeTypes"] = helper.Strings([]string{"cn", "dn"})
	service := CdwpgService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCdwpgDbconfigByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		respData = result
		return nil
	})
	if err != nil {
		return err
	}

	nodeConfigParamKeyMap := make(map[string]struct{})
	if v, ok := d.GetOk("node_config_params"); ok {
		nodeConfigParams := v.(*schema.Set).List()
		for _, nodeConfigParam := range nodeConfigParams {
			nodeConfigParamMap := nodeConfigParam.(map[string]interface{})
			key := fmt.Sprintf("%s-%s", nodeConfigParamMap["node_type"], nodeConfigParamMap["parameter_name"])
			nodeConfigParamKeyMap[key] = struct{}{}
		}
	}
	nodeConfigParams := make([]interface{}, 0)
	for _, nodeConfigParam := range respData {
		nodeType := nodeConfigParam.NodeType
		for _, configParam := range nodeConfigParam.Details {
			if nodeType != nil && configParam != nil && configParam.ParamName != nil {
				key := fmt.Sprintf("%s-%s", *nodeType, *configParam.ParamName)
				if _, ok := nodeConfigParamKeyMap[key]; ok {
					nodeConfigParams = append(nodeConfigParams, map[string]interface{}{
						"node_type":       nodeType,
						"parameter_name":  configParam.ParamName,
						"parameter_value": configParam.RunningValue,
					})
				}
			}
		}
	}
	_ = d.Set("instance_id", instanceId)
	_ = d.Set("node_config_params", nodeConfigParams)
	_ = instanceId
	_ = ctx
	return nil
}

func resourceTencentCloudCdwpgDbconfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdwpg_dbconfig.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	instanceId := d.Id()

	needChange := false
	mutableArgs := []string{"node_config_params"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		var respData []*cdwpgv20201230.ParamItem
		paramMap := make(map[string]interface{})
		paramMap["InstanceId"] = helper.String(instanceId)
		paramMap["NodeTypes"] = helper.Strings([]string{"cn", "dn"})
		service := CdwpgService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := service.DescribeCdwpgDbconfigByFilter(ctx, paramMap)
			if e != nil {
				return tccommon.RetryError(e)
			}
			respData = result
			return nil
		})
		if err != nil {
			return err
		}

		nodeConfigParamMap := make(map[string]string)
		for _, item := range respData {
			nodeType := item.NodeType
			for _, param := range item.Details {
				key := fmt.Sprintf("%s-%s", *nodeType, *param.ParamName)
				nodeConfigParamMap[key] = *param.RunningValue
			}
		}

		request := cdwpgv20201230.NewModifyDBParametersRequest()
		request.InstanceId = helper.String(instanceId)

		if v, ok := d.GetOk("node_config_params"); ok {
			params := v.(*schema.Set).List()
			for _, item := range params {
				nodeConfigParamsMap := item.(map[string]interface{})
				nodeConfigParams := cdwpgv20201230.NodeConfigParams{}
				var nodeType, parameterName, parameterValue string
				if v, ok := nodeConfigParamsMap["node_type"]; ok {
					nodeType = v.(string)
					nodeConfigParams.NodeType = helper.String(v.(string))
				}
				if v, ok := nodeConfigParamsMap["parameter_name"]; ok {
					parameterName = v.(string)
				}
				if v, ok := nodeConfigParamsMap["parameter_value"]; ok {
					parameterValue = v.(string)
				}

				key := fmt.Sprintf("%s-%s", nodeType, parameterName)

				configParams := cdwpgv20201230.ConfigParams{}
				configParams.ParameterName = helper.String(parameterName)
				configParams.ParameterOldValue = helper.String(nodeConfigParamMap[key])
				configParams.ParameterValue = helper.String(parameterValue)
				nodeConfigParams.ConfigParams = append(nodeConfigParams.ConfigParams, &configParams)
				request.NodeConfigParams = append(request.NodeConfigParams, &nodeConfigParams)
			}
		}

		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCdwpgV20201230Client().ModifyDBParametersWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update cdwpg dbconfig failed, reason:%+v", logId, err)
			return err
		}

		conf := tccommon.BuildStateChangeConf([]string{}, []string{"Serving"}, 10*tccommon.ReadRetryTimeout, time.Second, service.InstanceStateRefreshFunc(instanceId, []string{}))

		if _, e := conf.WaitForState(); e != nil {
			return e
		}

	}

	return resourceTencentCloudCdwpgDbconfigRead(d, meta)
}

func resourceTencentCloudCdwpgDbconfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdwpg_dbconfig.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
