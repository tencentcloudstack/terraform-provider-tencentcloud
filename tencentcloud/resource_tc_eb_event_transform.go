/*
Provides a resource to create a eb eb_transform

Example Usage

```hcl
resource "tencentcloud_eb_event_bus" "foo" {
  event_bus_name = "tf-event_bus"
  description    = "event bus desc"
  enable_store   = false
  save_days      = 1
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_eb_event_rule" "foo" {
  event_bus_id = tencentcloud_eb_event_bus.foo.id
  rule_name    = "tf-event_rule"
  description  = "event rule desc"
  enable       = true
  event_pattern = jsonencode(
    {
      source = "apigw.cloud.tencent"
      type = [
        "connector:apigw",
      ]
    }
  )
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_eb_event_transform" "foo" {
    event_bus_id = tencentcloud_eb_event_bus.foo.id
    rule_id      = tencentcloud_eb_event_rule.foo.rule_id

    transformations {

        extraction {
            extraction_input_path = "$"
            format                = "JSON"
        }

        transform {
            output_structs {
                key        = "type"
                value      = "connector:ckafka"
                value_type = "STRING"
            }
            output_structs {
                key        = "source"
                value      = "ckafka.cloud.tencent"
                value_type = "STRING"
            }
            output_structs {
                key        = "region"
                value      = "ap-guangzhou"
                value_type = "STRING"
            }
            output_structs {
                key        = "datacontenttype"
                value      = "application/json;charset=utf-8"
                value_type = "STRING"
            }
            output_structs {
                key        = "status"
                value      = "-"
                value_type = "STRING"
            }
            output_structs {
                key        = "data"
                value      = jsonencode(
                    {
                        Partition = 1
                        msgBody   = "Hello from Ckafka again!"
                        msgKey    = "test"
                        offset    = 37
                        topic     = "test-topic"
                    }
                )
                value_type = "STRING"
            }
        }
    }
}
```

Import

eb eb_transform can be imported using the id, e.g.

```
terraform import tencentcloud_eb_event_transform.eb_transform eb_transform_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	eb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/eb/v20210416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
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
				Description: "event bus Id.",
			},

			"rule_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ruleId.",
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
							Description: "Describe how to extract data.",
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
										Description: "Value: `TEXT`, `JSON`.",
									},
									"text_params": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Only Text needs to be passed.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"separator": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "`Comma`, `|`, `tab`, `space`, `newline`, `%`, `#`, the limit length is 1.",
												},
												"regex": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Fill in the regular expression: length 128.",
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
							Description: "Describe how to filter data.",
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
							Description: "Describe how to convert data.",
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
													Description: "The data type of value, optional values: `STRING`, `NUMBER`, `BOOLEAN`, `NULL`, `SYS_VARIABLE`, `JSONPATH`.",
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
		eventBusId       string
		ruleId           string
		transformationId string
	)
	if v, ok := d.GetOk("event_bus_id"); ok {
		eventBusId = v.(string)
		request.EventBusId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rule_id"); ok {
		ruleId = v.(string)
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
		log.Printf("[CRITAL]%s create eb ebTransform failed, reason:%+v", logId, err)
		return err
	}

	transformationId = *response.Response.TransformationId
	d.SetId(eventBusId + FILED_SP + ruleId + FILED_SP + transformationId)

	return resourceTencentCloudEbEventTransformRead(d, meta)
}

func resourceTencentCloudEbEventTransformRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eb_event_transform.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := EbService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	eventBusId := idSplit[0]
	ruleId := idSplit[1]
	transformationId := idSplit[2]

	ebTransform, err := service.DescribeEbEventTransformById(ctx, eventBusId, ruleId, transformationId)
	if err != nil {
		return err
	}

	if ebTransform == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `EbEventTransform` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("event_bus_id", eventBusId)
	_ = d.Set("rule_id", ruleId)

	if ebTransform != nil {
		transformationsMap := map[string]interface{}{}

		if ebTransform.Extraction != nil {
			extractionMap := map[string]interface{}{}

			if ebTransform.Extraction.ExtractionInputPath != nil {
				extractionMap["extraction_input_path"] = ebTransform.Extraction.ExtractionInputPath
			}

			if ebTransform.Extraction.Format != nil {
				extractionMap["format"] = ebTransform.Extraction.Format
			}

			if ebTransform.Extraction.TextParams != nil {
				textParamsMap := map[string]interface{}{}

				if ebTransform.Extraction.TextParams.Separator != nil {
					textParamsMap["separator"] = ebTransform.Extraction.TextParams.Separator
				}

				if ebTransform.Extraction.TextParams.Regex != nil {
					textParamsMap["regex"] = ebTransform.Extraction.TextParams.Regex
				}

				extractionMap["text_params"] = []interface{}{textParamsMap}
			}

			transformationsMap["extraction"] = []interface{}{extractionMap}
		}

		if ebTransform.EtlFilter != nil {
			etlFilterMap := map[string]interface{}{}

			if ebTransform.EtlFilter.Filter != nil {
				etlFilterMap["filter"] = ebTransform.EtlFilter.Filter
			}

			transformationsMap["etl_filter"] = []interface{}{etlFilterMap}
		}

		if ebTransform.Transform != nil {
			transformMap := map[string]interface{}{}

			if ebTransform.Transform.OutputStructs != nil {
				outputStructsList := []interface{}{}
				for _, outputStructs := range ebTransform.Transform.OutputStructs {
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

				transformMap["output_structs"] = outputStructsList
			}

			transformationsMap["transform"] = []interface{}{transformMap}
		}

		_ = d.Set("transformations", []interface{}{transformationsMap})

	}

	return nil
}

func resourceTencentCloudEbEventTransformUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eb_event_transform.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := eb.NewUpdateTransformationRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	eventBusId := idSplit[0]
	ruleId := idSplit[1]
	transformationId := idSplit[2]

	request.EventBusId = &eventBusId
	request.RuleId = &ruleId
	request.TransformationId = &transformationId

	immutableArgs := []string{"event_bus_id", "rule_id"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("transformations") {
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
		log.Printf("[CRITAL]%s update eb ebTransform failed, reason:%+v", logId, err)
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
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	eventBusId := idSplit[0]
	ruleId := idSplit[1]
	transformationId := idSplit[2]

	if err := service.DeleteEbEventTransformById(ctx, eventBusId, ruleId, transformationId); err != nil {
		return err
	}

	return nil
}
