package teo

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoFunctionComponentBinding() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoFunctionComponentBindingConfigCreate,
		Read:   resourceTencentCloudTeoFunctionComponentBindingConfigRead,
		Update: resourceTencentCloudTeoFunctionComponentBindingConfigUpdate,
		Delete: resourceTencentCloudTeoFunctionComponentBindingConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Site ID.",
			},
			"function_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Function ID.",
			},
			"function_component_bindings": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Function component binding list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The type of the bound component. Valid values: `kv_namespace`.",
						},
						"variable_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The variable name used for binding.",
						},
						"kv_namespace_parameters": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "KV namespace configuration parameters. Required when type is `kv_namespace`.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zone_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The site ID to which the KV namespace belongs.",
									},
									"namespace": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "KV namespace name.",
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

func buildFunctionComponentBindings(d *schema.ResourceData) []*teov20220901.FunctionComponentBinding {
	rawBindings := d.Get("function_component_bindings").([]interface{})
	bindings := make([]*teov20220901.FunctionComponentBinding, 0, len(rawBindings))
	for _, raw := range rawBindings {
		bindingMap := raw.(map[string]interface{})
		binding := &teov20220901.FunctionComponentBinding{}
		binding.Type = helper.String(bindingMap["type"].(string))
		binding.VariableName = helper.String(bindingMap["variable_name"].(string))
		if kvParamsRaw, ok := bindingMap["kv_namespace_parameters"]; ok {
			kvParamsList := kvParamsRaw.([]interface{})
			if len(kvParamsList) > 0 {
				kvParamsMap := kvParamsList[0].(map[string]interface{})
				binding.KVNamespaceParameters = &teov20220901.KVNamespaceParameters{
					ZoneId:    helper.String(kvParamsMap["zone_id"].(string)),
					Namespace: helper.String(kvParamsMap["namespace"].(string)),
				}
			}
		}
		bindings = append(bindings, binding)
	}
	return bindings
}

func resourceTencentCloudTeoFunctionComponentBindingConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_component_binding.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		zoneId     = d.Get("zone_id").(string)
		functionId = d.Get("function_id").(string)
		request    = teov20220901.NewModifyFunctionComponentBindingsRequest()
	)

	request.ZoneId = &zoneId
	request.FunctionId = &functionId
	operation := "rebind"
	request.Operation = &operation
	request.FunctionComponentBindings = buildFunctionComponentBindings(d)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyFunctionComponentBindingsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create teo function component binding failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{zoneId, functionId}, tccommon.FILED_SP))

	return resourceTencentCloudTeoFunctionComponentBindingConfigRead(d, meta)
}

func resourceTencentCloudTeoFunctionComponentBindingConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_component_binding.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		id      = d.Id()
		idParts = strings.Split(id, tccommon.FILED_SP)
	)

	if len(idParts) != 2 {
		return fmt.Errorf("resource ID format error, expected zone_id%sfunction_id", tccommon.FILED_SP)
	}

	zoneId := idParts[0]
	functionId := idParts[1]

	var allBindings []*teov20220901.FunctionComponentBinding
	var offset int64

	for {
		request := teov20220901.NewDescribeFunctionComponentBindingsRequest()
		request.ZoneId = &zoneId
		request.FunctionId = &functionId
		request.Offset = &offset
		limit := int64(1000)
		request.Limit = &limit

		var response *teov20220901.DescribeFunctionComponentBindingsResponse
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DescribeFunctionComponentBindingsWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			response = result
			return nil
		})
		if err != nil {
			return err
		}

		if response == nil || response.Response == nil {
			break
		}

		if response.Response.FunctionComponentBindings != nil {
			allBindings = append(allBindings, response.Response.FunctionComponentBindings...)
		}

		totalCount := int64(0)
		if response.Response.TotalCount != nil {
			totalCount = *response.Response.TotalCount
		}

		offset += limit
		if offset >= totalCount {
			break
		}
	}

	_ = d.Set("zone_id", zoneId)
	_ = d.Set("function_id", functionId)

	bindingList := make([]map[string]interface{}, 0, len(allBindings))
	for _, binding := range allBindings {
		bindingMap := map[string]interface{}{}
		if binding.Type != nil {
			bindingMap["type"] = *binding.Type
		}
		if binding.VariableName != nil {
			bindingMap["variable_name"] = *binding.VariableName
		}
		if binding.KVNamespaceParameters != nil {
			kvParams := make([]map[string]interface{}, 0, 1)
			kvParamMap := map[string]interface{}{}
			if binding.KVNamespaceParameters.ZoneId != nil {
				kvParamMap["zone_id"] = *binding.KVNamespaceParameters.ZoneId
			}
			if binding.KVNamespaceParameters.Namespace != nil {
				kvParamMap["namespace"] = *binding.KVNamespaceParameters.Namespace
			}
			kvParams = append(kvParams, kvParamMap)
			bindingMap["kv_namespace_parameters"] = kvParams
		}
		bindingList = append(bindingList, bindingMap)
	}
	_ = d.Set("function_component_bindings", bindingList)

	return nil
}

func resourceTencentCloudTeoFunctionComponentBindingConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_component_binding.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		id      = d.Id()
		idParts = strings.Split(id, tccommon.FILED_SP)
	)

	if len(idParts) != 2 {
		return fmt.Errorf("resource ID format error, expected zone_id%sfunction_id", tccommon.FILED_SP)
	}

	zoneId := idParts[0]
	functionId := idParts[1]

	if d.HasChange("function_component_bindings") {
		request := teov20220901.NewModifyFunctionComponentBindingsRequest()
		request.ZoneId = &zoneId
		request.FunctionId = &functionId
		operation := "rebind"
		request.Operation = &operation
		request.FunctionComponentBindings = buildFunctionComponentBindings(d)

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyFunctionComponentBindingsWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update teo function component binding failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudTeoFunctionComponentBindingConfigRead(d, meta)
}

func resourceTencentCloudTeoFunctionComponentBindingConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_component_binding.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		id      = d.Id()
		idParts = strings.Split(id, tccommon.FILED_SP)
	)

	if len(idParts) != 2 {
		return fmt.Errorf("resource ID format error, expected zone_id%sfunction_id", tccommon.FILED_SP)
	}

	zoneId := idParts[0]
	functionId := idParts[1]

	request := teov20220901.NewModifyFunctionComponentBindingsRequest()
	request.ZoneId = &zoneId
	request.FunctionId = &functionId
	operation := "rebind"
	request.Operation = &operation
	request.FunctionComponentBindings = []*teov20220901.FunctionComponentBinding{}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyFunctionComponentBindingsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete teo function component binding failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
