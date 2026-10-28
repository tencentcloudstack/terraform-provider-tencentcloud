package waf

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// apiSecApiNameOpSchema returns the schema of the `api_name_op` block (SDK ApiNameOp), shared by the
// privilege/custom event/white rule resources.
func apiSecApiNameOpSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "API match list.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"value": {
					Type:        schema.TypeSet,
					Optional:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Description: "Match value list.",
				},
				"op": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Match method, such as belong and regex.",
				},
				"api_name_method": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "When manually filtering, this structure should be passed.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"api_name": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "API name.",
							},
							"method": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "API request method.",
							},
							"count": {
								Type:        schema.TypeInt,
								Computed:    true,
								Description: "API request count in the last 30 days.",
							},
							"label": {
								Type:        schema.TypeSet,
								Computed:    true,
								Elem:        &schema.Schema{Type: schema.TypeString},
								Description: "API label.",
							},
						},
					},
				},
			},
		},
	}
}

func buildApiSecApiNameOpList(list []interface{}) []*waf.ApiNameOp {
	result := make([]*waf.ApiNameOp, 0, len(list))
	for _, item := range list {
		dMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		apiNameOp := waf.ApiNameOp{}
		if v, ok := dMap["value"]; ok {
			apiNameOp.Value = helper.InterfacesStringsPoint(v.(*schema.Set).List())
		}

		if v, ok := dMap["op"]; ok && v.(string) != "" {
			apiNameOp.Op = helper.String(v.(string))
		}

		if v, ok := dMap["api_name_method"]; ok {
			for _, methodItem := range v.([]interface{}) {
				methodMap, ok := methodItem.(map[string]interface{})
				if !ok {
					continue
				}

				apiNameMethod := waf.ApiNameMethod{}
				if vv, ok := methodMap["api_name"]; ok && vv.(string) != "" {
					apiNameMethod.ApiName = helper.String(vv.(string))
				}

				if vv, ok := methodMap["method"]; ok && vv.(string) != "" {
					apiNameMethod.Method = helper.String(vv.(string))
				}

				apiNameOp.ApiNameMethod = append(apiNameOp.ApiNameMethod, &apiNameMethod)
			}
		}

		result = append(result, &apiNameOp)
	}

	return result
}

func flattenApiSecApiNameOpList(list []*waf.ApiNameOp) []interface{} {
	result := make([]interface{}, 0, len(list))
	for _, item := range list {
		if item == nil {
			continue
		}

		dMap := map[string]interface{}{}
		if item.Value != nil {
			dMap["value"] = item.Value
		}

		if item.Op != nil {
			dMap["op"] = item.Op
		}

		if item.ApiNameMethod != nil {
			methods := make([]interface{}, 0, len(item.ApiNameMethod))
			for _, methodItem := range item.ApiNameMethod {
				if methodItem == nil {
					continue
				}

				methodMap := map[string]interface{}{}
				if methodItem.ApiName != nil {
					methodMap["api_name"] = methodItem.ApiName
				}

				if methodItem.Method != nil {
					methodMap["method"] = methodItem.Method
				}

				if methodItem.Count != nil {
					methodMap["count"] = methodItem.Count
				}

				if methodItem.Label != nil {
					methodMap["label"] = methodItem.Label
				}

				methods = append(methods, methodMap)
			}

			dMap["api_name_method"] = methods
		}

		result = append(result, dMap)
	}

	return result
}

// apiSecSceneRuleEntrySchema returns the schema of an `ApiSecSceneRuleEntry` list block, shared by
// the scene/custom event rule resources.
func apiSecSceneRuleEntrySchema(description string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: description,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"key": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Match field.",
				},
				"value": {
					Type:        schema.TypeSet,
					Optional:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Description: "Match value.",
				},
				"operate": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Operator.",
				},
				"name": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "When the match field is get parameter value, post parameter value, cookie parameter value, header parameter value or rsp parameter value, this field can be filled.",
				},
			},
		},
	}
}

func buildApiSecSceneRuleEntryList(list []interface{}) []*waf.ApiSecSceneRuleEntry {
	result := make([]*waf.ApiSecSceneRuleEntry, 0, len(list))
	for _, item := range list {
		dMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		entry := waf.ApiSecSceneRuleEntry{}
		if v, ok := dMap["key"]; ok && v.(string) != "" {
			entry.Key = helper.String(v.(string))
		}

		if v, ok := dMap["value"]; ok {
			entry.Value = helper.InterfacesStringsPoint(v.(*schema.Set).List())
		}

		if v, ok := dMap["operate"]; ok && v.(string) != "" {
			entry.Operate = helper.String(v.(string))
		}

		if v, ok := dMap["name"]; ok && v.(string) != "" {
			entry.Name = helper.String(v.(string))
		}

		result = append(result, &entry)
	}

	return result
}

func flattenApiSecSceneRuleEntryList(list []*waf.ApiSecSceneRuleEntry) []interface{} {
	result := make([]interface{}, 0, len(list))
	for _, item := range list {
		if item == nil {
			continue
		}

		dMap := map[string]interface{}{}
		if item.Key != nil {
			dMap["key"] = item.Key
		}

		if item.Value != nil {
			dMap["value"] = item.Value
		}

		if item.Operate != nil {
			dMap["operate"] = item.Operate
		}

		if item.Name != nil {
			dMap["name"] = item.Name
		}

		result = append(result, dMap)
	}

	return result
}
