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
		Create: resourceTencentCloudTeoFunctionComponentBindingCreate,
		Read:   resourceTencentCloudTeoFunctionComponentBindingRead,
		Update: resourceTencentCloudTeoFunctionComponentBindingUpdate,
		Delete: resourceTencentCloudTeoFunctionComponentBindingDelete,
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
							Description: "The type of the bound component. Valid value: `kv_namespace`.",
						},
						"variable_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The variable name used for binding, limited to 1-50 characters.",
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
										Description: "The site ID of the KV namespace.",
									},
									"namespace": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The KV namespace name.",
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

func resourceTencentCloudTeoFunctionComponentBindingCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_component_binding.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		zoneId     string
		functionId string
	)
	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}
	if v, ok := d.GetOk("function_id"); ok {
		functionId = v.(string)
	}

	d.SetId(strings.Join([]string{zoneId, functionId}, tccommon.FILED_SP))

	return resourceTencentCloudTeoFunctionComponentBindingUpdate(d, meta)
}

func resourceTencentCloudTeoFunctionComponentBindingRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_component_binding.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	functionId := idSplit[1]

	_ = d.Set("zone_id", zoneId)
	_ = d.Set("function_id", functionId)

	var allBindings []*teov20220901.FunctionComponentBinding
	var offset int64

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		allBindings = nil
		offset = 0
		for {
			request := teov20220901.NewDescribeFunctionComponentBindingsRequest()
			request.ZoneId = helper.String(zoneId)
			request.FunctionId = helper.String(functionId)
			request.Offset = helper.Int64(offset)
			request.Limit = helper.Int64(1000)

			response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DescribeFunctionComponentBindingsWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

			if response == nil || response.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("read teo_function_component_binding failed, response is nil, id: %s", d.Id()))
			}

			if response.Response.FunctionComponentBindings != nil {
				allBindings = append(allBindings, response.Response.FunctionComponentBindings...)
			}

			if response.Response.TotalCount == nil || int64(len(allBindings)) >= *response.Response.TotalCount {
				break
			}
			offset += 1000
		}
		return nil
	})
	if err != nil {
		return err
	}

	bindingsList := make([]map[string]interface{}, 0, len(allBindings))
	for _, binding := range allBindings {
		bindingMap := map[string]interface{}{}

		if binding.Type != nil {
			bindingMap["type"] = binding.Type
		}

		if binding.VariableName != nil {
			bindingMap["variable_name"] = binding.VariableName
		}

		if binding.KVNamespaceParameters != nil {
			kvParams := map[string]interface{}{}
			if binding.KVNamespaceParameters.ZoneId != nil {
				kvParams["zone_id"] = binding.KVNamespaceParameters.ZoneId
			}
			if binding.KVNamespaceParameters.Namespace != nil {
				kvParams["namespace"] = binding.KVNamespaceParameters.Namespace
			}
			bindingMap["kv_namespace_parameters"] = []interface{}{kvParams}
		}

		bindingsList = append(bindingsList, bindingMap)
	}

	_ = d.Set("function_component_bindings", bindingsList)

	return nil
}

func resourceTencentCloudTeoFunctionComponentBindingUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_component_binding.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	functionId := idSplit[1]

	needChange := false
	mutableArgs := []string{"function_component_bindings"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := teov20220901.NewModifyFunctionComponentBindingsRequest()
		request.ZoneId = helper.String(zoneId)
		request.FunctionId = helper.String(functionId)
		request.Operation = helper.String("rebind")

		if v, ok := d.GetOk("function_component_bindings"); ok {
			for _, item := range v.([]interface{}) {
				bindingMap := item.(map[string]interface{})
				binding := teov20220901.FunctionComponentBinding{}
				if v, ok := bindingMap["type"]; ok {
					binding.Type = helper.String(v.(string))
				}
				if v, ok := bindingMap["variable_name"]; ok {
					binding.VariableName = helper.String(v.(string))
				}
				if v, ok := bindingMap["kv_namespace_parameters"]; ok {
					kvParamsList := v.([]interface{})
					if len(kvParamsList) > 0 {
						kvParamsMap := kvParamsList[0].(map[string]interface{})
						kvParams := teov20220901.KVNamespaceParameters{}
						if v, ok := kvParamsMap["zone_id"]; ok {
							kvParams.ZoneId = helper.String(v.(string))
						}
						if v, ok := kvParamsMap["namespace"]; ok {
							kvParams.Namespace = helper.String(v.(string))
						}
						binding.KVNamespaceParameters = &kvParams
					}
				}
				request.FunctionComponentBindings = append(request.FunctionComponentBindings, &binding)
			}
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyFunctionComponentBindingsWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update teo_function_component_binding failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudTeoFunctionComponentBindingRead(d, meta)
}

func resourceTencentCloudTeoFunctionComponentBindingDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_component_binding.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
