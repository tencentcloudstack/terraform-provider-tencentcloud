/*
Provides a resource to create a eb event_transform

Example Usage

```hcl
resource "tencentcloud_eb_event_transform" "event_transform" {
  event_bus_id = ""
  rule_id = ""
  transformations {
		extraction {
			extraction_input_path = ""
			format = ""
			text_params {
				separator = ""
				regex = ""
			}
		}
		etl_filter {
			filter = ""
		}
		transform {
			output_structs {
				key = ""
				value = ""
				value_type = ""
			}
		}

  }
}
```

Import

eb event_transform can be imported using the id, e.g.

```
terraform import tencentcloud_eb_event_transform.event_transform event_transform_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	eb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/eb/v20210416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudEbEventTransform() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudEbEventTransformCreate,
		Read:   resourceTencentCloudEbEventTransformRead,
		Update: resourceTencentCloudEbEventTransformUpdate,
		Delete: resourceTencentCloudEbEventTransformDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"event_bus_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Event bus Id.",
			},

			"rule_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "RuleId.",
			},

			"transformations": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "A list of transformation rules, currently only one.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"extraction": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Describe how to extract data. Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"extraction_input_path": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "JsonPath, if not specified, the default value $.",
									},
									"format": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "&amp;#39;Value: TEXT/JSON&amp;#39;.",
									},
									"text_params": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Only Text needs to be passed. Note: this field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"separator": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Comma, | , tab, space, newline, %, #, the limit length is 1. Note: This field may return null, indicating that no valid value can be obtained.",
												},
												"regex": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Fill in the regular expression: length 128, note: this field may return null, indicating that no valid value can be obtained.",
												},
											},
										},
									},
								},
							},
						},
						"etl_filter": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Describe how to filter data, note: this field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"filter": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Grammatical Rules are consistent.",
									},
								},
							},
						},
						"transform": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Describe how to convert data. Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"output_structs": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Describe how the data is transformed.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Corresponding to the key in the output json.",
												},
												"value": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "You can fill in the json-path and also support constants or built-in keyword date types.",
												},
												"value_type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "&amp;#39;value data type, optional values: STRING, NUMBER, BOOLEAN, NULL, SYS_VARIABLE, JSONPATH&amp;#39;.",
												},
											},
										},
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

func resourceTencentCloudEbEventTransformCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eb_event_transform.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request          = eb.NewCreateTransformationRequest()
		response         = eb.NewCreateTransformationResponse()
		transformationId string
		eventBusId       string
	)
	if v, ok := d.GetOk("event_bus_id"); ok {
		eventBusId = v.(string)
		request.EventBusId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rule_id"); ok {
		request.RuleId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("transformations"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			transformation := eb.Transformation{}
			if extractionMap, ok := helper.InterfaceToMap(dMap, "extraction"); ok {
				extraction := eb.Extraction{}
				if v, ok := extractionMap["extraction_input_path"]; ok {
					extraction.ExtractionInputPath = helper.String(v.(string))
				}
				if v, ok := extractionMap["format"]; ok {
					extraction.Format = helper.String(v.(string))
				}
				if textParamsMap, ok := helper.InterfaceToMap(extractionMap, "text_params"); ok {
					textParams := eb.TextParams{}
					if v, ok := textParamsMap["separator"]; ok {
						textParams.Separator = helper.String(v.(string))
					}
					if v, ok := textParamsMap["regex"]; ok {
						textParams.Regex = helper.String(v.(string))
					}
					extraction.TextParams = &textParams
				}
				transformation.Extraction = &extraction
			}
			if etlFilterMap, ok := helper.InterfaceToMap(dMap, "etl_filter"); ok {
				etlFilter := eb.EtlFilter{}
				if v, ok := etlFilterMap["filter"]; ok {
					etlFilter.Filter = helper.String(v.(string))
				}
				transformation.EtlFilter = &etlFilter
			}
			if transformMap, ok := helper.InterfaceToMap(dMap, "transform"); ok {
				transform := eb.Transform{}
				if v, ok := transformMap["output_structs"]; ok {
					for _, item := range v.([]interface{}) {
						outputStructsMap := item.(map[string]interface{})
						outputStructParam := eb.OutputStructParam{}
						if v, ok := outputStructsMap["key"]; ok {
							outputStructParam.Key = helper.String(v.(string))
						}
						if v, ok := outputStructsMap["value"]; ok {
							outputStructParam.Value = helper.String(v.(string))
						}
						if v, ok := outputStructsMap["value_type"]; ok {
							outputStructParam.ValueType = helper.String(v.(string))
						}
						transform.OutputStructs = append(transform.OutputStructs, &outputStructParam)
					}
				}
				transformation.Transform = &transform
			}
			request.Transformations = append(request.Transformations, &transformation)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseEbClient().CreateTransformation(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create eb eventTransform failed, reason:%+v", logId, err)
		return err
	}

	transformationId = *response.Response.TransformationId
	d.SetId(strings.Join([]string{transformationId, eventBusId}, FILED_SP))

	return resourceTencentCloudEbEventTransformRead(d, meta)
}

func resourceTencentCloudEbEventTransformRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eb_event_transform.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := EbService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	transformationId := idSplit[0]
	eventBusId := idSplit[1]

	eventTransform, err := service.DescribeEbEventTransformById(ctx, transformationId, eventBusId)
	if err != nil {
		return err
	}

	if eventTransform == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `EbEventTransform` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if eventTransform.EventBusId != nil {
		_ = d.Set("event_bus_id", eventTransform.EventBusId)
	}

	if eventTransform.RuleId != nil {
		_ = d.Set("rule_id", eventTransform.RuleId)
	}

	if eventTransform.Transformations != nil {
		transformationsList := []interface{}{}
		for _, transformations := range eventTransform.Transformations {
			transformationsMap := map[string]interface{}{}

			if eventTransform.Transformations.Extraction != nil {
				extractionMap := map[string]interface{}{}

				if eventTransform.Transformations.Extraction.ExtractionInputPath != nil {
					extractionMap["extraction_input_path"] = eventTransform.Transformations.Extraction.ExtractionInputPath
				}

				if eventTransform.Transformations.Extraction.Format != nil {
					extractionMap["format"] = eventTransform.Transformations.Extraction.Format
				}

				if eventTransform.Transformations.Extraction.TextParams != nil {
					textParamsMap := map[string]interface{}{}

					if eventTransform.Transformations.Extraction.TextParams.Separator != nil {
						textParamsMap["separator"] = eventTransform.Transformations.Extraction.TextParams.Separator
					}

					if eventTransform.Transformations.Extraction.TextParams.Regex != nil {
						textParamsMap["regex"] = eventTransform.Transformations.Extraction.TextParams.Regex
					}

					extractionMap["text_params"] = []interface{}{textParamsMap}
				}

				transformationsMap["extraction"] = []interface{}{extractionMap}
			}

			if eventTransform.Transformations.EtlFilter != nil {
				etlFilterMap := map[string]interface{}{}

				if eventTransform.Transformations.EtlFilter.Filter != nil {
					etlFilterMap["filter"] = eventTransform.Transformations.EtlFilter.Filter
				}

				transformationsMap["etl_filter"] = []interface{}{etlFilterMap}
			}

			if eventTransform.Transformations.Transform != nil {
				transformMap := map[string]interface{}{}

				if eventTransform.Transformations.Transform.OutputStructs != nil {
					outputStructsList := []interface{}{}
					for _, outputStructs := range eventTransform.Transformations.Transform.OutputStructs {
						outputStructsMap := map[string]interface{}{}

						if outputStructs.Key != nil {
							outputStructsMap["key"] = outputStructs.Key
						}

						if outputStructs.Value != nil {
							outputStructsMap["value"] = outputStructs.Value
						}

						if outputStructs.ValueType != nil {
							outputStructsMap["value_type"] = outputStructs.ValueType
						}

						outputStructsList = append(outputStructsList, outputStructsMap)
					}

					transformMap["output_structs"] = []interface{}{outputStructsList}
				}

				transformationsMap["transform"] = []interface{}{transformMap}
			}

			transformationsList = append(transformationsList, transformationsMap)
		}

		_ = d.Set("transformations", transformationsList)

	}

	return nil
}

func resourceTencentCloudEbEventTransformUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eb_event_transform.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := eb.NewUpdateTransformationRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	transformationId := idSplit[0]
	eventBusId := idSplit[1]

	request.TransformationId = &transformationId
	request.EventBusId = &eventBusId

	immutableArgs := []string{"event_bus_id", "rule_id", "transformations"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("event_bus_id") {
		if v, ok := d.GetOk("event_bus_id"); ok {
			request.EventBusId = helper.String(v.(string))
		}
	}

	if d.HasChange("rule_id") {
		if v, ok := d.GetOk("rule_id"); ok {
			request.RuleId = helper.String(v.(string))
		}
	}

	if d.HasChange("transformations") {
		if v, ok := d.GetOk("transformations"); ok {
			for _, item := range v.([]interface{}) {
				transformation := eb.Transformation{}
				if extractionMap, ok := helper.InterfaceToMap(dMap, "extraction"); ok {
					extraction := eb.Extraction{}
					if v, ok := extractionMap["extraction_input_path"]; ok {
						extraction.ExtractionInputPath = helper.String(v.(string))
					}
					if v, ok := extractionMap["format"]; ok {
						extraction.Format = helper.String(v.(string))
					}
					if textParamsMap, ok := helper.InterfaceToMap(extractionMap, "text_params"); ok {
						textParams := eb.TextParams{}
						if v, ok := textParamsMap["separator"]; ok {
							textParams.Separator = helper.String(v.(string))
						}
						if v, ok := textParamsMap["regex"]; ok {
							textParams.Regex = helper.String(v.(string))
						}
						extraction.TextParams = &textParams
					}
					transformation.Extraction = &extraction
				}
				if etlFilterMap, ok := helper.InterfaceToMap(dMap, "etl_filter"); ok {
					etlFilter := eb.EtlFilter{}
					if v, ok := etlFilterMap["filter"]; ok {
						etlFilter.Filter = helper.String(v.(string))
					}
					transformation.EtlFilter = &etlFilter
				}
				if transformMap, ok := helper.InterfaceToMap(dMap, "transform"); ok {
					transform := eb.Transform{}
					if v, ok := transformMap["output_structs"]; ok {
						for _, item := range v.([]interface{}) {
							outputStructsMap := item.(map[string]interface{})
							outputStructParam := eb.OutputStructParam{}
							if v, ok := outputStructsMap["key"]; ok {
								outputStructParam.Key = helper.String(v.(string))
							}
							if v, ok := outputStructsMap["value"]; ok {
								outputStructParam.Value = helper.String(v.(string))
							}
							if v, ok := outputStructsMap["value_type"]; ok {
								outputStructParam.ValueType = helper.String(v.(string))
							}
							transform.OutputStructs = append(transform.OutputStructs, &outputStructParam)
						}
					}
					transformation.Transform = &transform
				}
				request.Transformations = append(request.Transformations, &transformation)
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseEbClient().UpdateTransformation(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update eb eventTransform failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudEbEventTransformRead(d, meta)
}

func resourceTencentCloudEbEventTransformDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eb_event_transform.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := EbService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	transformationId := idSplit[0]
	eventBusId := idSplit[1]

	if err := service.DeleteEbEventTransformById(ctx, transformationId, eventBusId); err != nil {
		return err
	}

	return nil
}
