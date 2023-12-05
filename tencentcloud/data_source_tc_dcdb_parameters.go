package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDcdbParameters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDcdbParametersRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "instance id.",
			},

			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "parameter list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"param": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "parameter name.",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "parameter value.",
						},
						"default": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "default value.",
						},
						"constraint": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "params constraint.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "type.",
									},
									"enum": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "a list of optional values of type num.",
									},
									"range": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "range constraint.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"min": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "min value.",
												},
												"max": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "max value.",
												},
											},
										},
									},
									"string": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "constraint type is string.",
									},
								},
							},
						},
						"have_set_value": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "have set value.",
						},
						"need_restart": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "need restart.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudDcdbParametersRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dcdb_parameters.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["instance_id"] = helper.String(v.(string))
	}

	dcdbService := DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var params []*dcdb.ParamDesc
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := dcdbService.DescribeDcdbParametersByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		params = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Dcdb params failed, reason:%+v", logId, err)
		return err
	}

	paramList := []interface{}{}
	ids := make([]string, 0, len(params))
	if params != nil {
		for _, param := range params {
			paramMap := map[string]interface{}{}
			if param.Param != nil {
				paramMap["param"] = param.Param
			}
			if param.Value != nil {
				paramMap["value"] = param.Value
			}
			if param.Default != nil {
				paramMap["default"] = param.Default
			}
			if param.Constraint != nil {
				constraintMap := map[string]interface{}{}
				if param.Constraint.Type != nil {
					constraintMap["type"] = param.Constraint.Type
				}
				if param.Constraint.Enum != nil {
					constraintMap["enum"] = param.Constraint.Enum
				}
				if param.Constraint.Range != nil {
					rangeMap := map[string]interface{}{}
					if param.Constraint.Range.Min != nil {
						rangeMap["min"] = param.Constraint.Range.Min
					}
					if param.Constraint.Range.Max != nil {
						rangeMap["max"] = param.Constraint.Range.Max
					}

					constraintMap["range"] = []interface{}{rangeMap}
				}
				if param.Constraint.String != nil {
					constraintMap["string"] = param.Constraint.String
				}

				paramMap["constraint"] = []interface{}{constraintMap}
			}
			if param.HaveSetValue != nil {
				paramMap["have_set_value"] = param.HaveSetValue
			}
			if param.NeedRestart != nil {
				paramMap["need_restart"] = param.NeedRestart
			}
			ids = append(ids, *param.Param+FILED_SP+*param.Value)
			paramList = append(paramList, paramMap)
		}
		d.SetId(helper.DataResourceIdsHash(ids))
		_ = d.Set("list", paramList)
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), paramList); e != nil {
			return e
		}
	}

	return nil
}
