/*
Provides a resource to create a postgresql parameter_template

Example Usage

```hcl
resource "tencentcloud_postgresql_parameter_template" "parameter_template" {
  template_name = "your_temp_name"
  db_major_version = "13"
  db_engine = "postgresql"
  template_description = "For_tf_test"

  modify_param_entry_set {
	name = "timezone"
	expected_value = "UTC"
  }
  modify_param_entry_set {
	name = "lock_timeout"
	expected_value = "123"
  }

  delete_param_set = ["lc_time"]
}
```

Import

postgresql parameter_template can be imported using the id, e.g.

Notice: `modify_param_entry_set` and `delete_param_set` do not support import.

```
terraform import tencentcloud_postgresql_parameter_template.parameter_template parameter_template_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudPostgresqlParameterTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlParameterTemplateCreate,
		Read:   resourceTencentCloudPostgresqlParameterTemplateRead,
		Update: resourceTencentCloudPostgresqlParameterTemplateUpdate,
		Delete: resourceTencentCloudPostgresqlParameterTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"template_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Template name, which can contain 1-60 letters, digits, and symbols (-_./()+=:@).",
			},

			"db_major_version": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The major database version number, such as 11, 12, 13.",
			},

			"db_engine": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Database engine, such as postgresql, mssql_compatible.",
			},

			"template_description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Parameter template description, which can contain 1-60 letters, digits, and symbols (-_./()+=:@).",
			},

			"modify_param_entry_set": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "The set of parameters that need to be modified or added. Note: the same parameter cannot appear in the set of modifying and adding and deleting at the same time.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The parameter name.",
						},
						"expected_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Modify the parameter value. The input parameters are passed in the form of strings, for example: decimal `0.1`, integer `1000`, enumeration `replica`.",
						},
					},
				},
			},

			"delete_param_set": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "The set of parameters that need to be deleted.",
			},

			"param_info_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Parameter information contained in the parameter template. Note: This field may return null, indicating that no valid value can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "parameter ID. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "parameter name. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"param_value_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parameter value type: integer (integer), real (floating point), bool (Boolean), enum (enumeration type), mutil_enum (enumeration type, support multiple choices).When the parameter type is For integer (integer type) and real (floating point type), the value range of the parameter is determined according to the Max and Min of the return value; When the parameter type is bool (Boolean type), the value range of the parameter setting value is true | false; When the parameter type is enum (enumeration type) or mutil_enum (multi-enumeration type), the value range of the parameter is determined by the EnumValue in the return value. Note: This field may return null, indicating that it cannot be fetched to a valid value.",
						},
						"unit": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parameter Value Unit. When the parameter has no unit, the field returns empty. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"default_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parameter default value. returned as a string. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"current_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The current running value of the parameter. returned as a string. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"max": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Numerical type (integer, real) parameters, value lower bound. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"enum_value": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Enumeration type parameter, value range. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"min": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Numeric type (integer, real) parameters, value upper bound. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"param_description_ch": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parameter Chinese description. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"param_description_en": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parameter English description. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"need_reboot": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Parameter modification, whether to restart to take effect. (true is required, false is not required) Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"classification_cn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parameter Chinese classification. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"classification_en": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parameter English Classification. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"spec_related": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Is it related to the specification. (true means related, false means don&#39;t want to close) Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"advanced": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether it is a key parameter. (true is the key parameter, the modification needs to be paid attention to, which may affect the performance of the instance) Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"last_modify_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "parameter last modified time. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"standby_related": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "There are master-standby constraints on the parameters, 0: no master-standby constraint relationship, 1: the parameter value of the standby machine must be greater than that of the master machine, 2: the parameter value of the master machine must be greater than that of the standby machine. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"version_relation_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Parameter version association information, storing specific parameter information under a specific kernel version. Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "parameter name. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"db_kernel_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The kernel version of the parameter information. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The default value of the parameter in this version and this specification. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"unit": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unit of the parameter value. When the parameter has no unit, the field returns empty. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"max": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Numeric type (integer, real) parameters, value upper bound. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"min": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Numerical type (integer, real) parameters, value lower bound. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"enum_value": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Enumeration type parameter, value range. Note: This field may return null, indicating that no valid value can be obtained.",
									},
								},
							},
						},
						"spec_relation_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Parameter specification related information, storing specific parameter information under specific specifications. Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "parameter name. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"memory": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The parameter information belongs to the specification. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The default value of the parameter under this specification. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"unit": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unit of the parameter value. When the parameter has no unit, the field returns empty. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"max": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Numeric type (integer, real) parameters, value upper bound. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"min": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Numerical type (integer, real) parameters, value lower bound. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"enum_value": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Enumeration type parameter, value range. Note: This field may return null, indicating that no valid value can be obtained.",
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

func resourceTencentCloudPostgresqlParameterTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_parameter_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = postgresql.NewCreateParameterTemplateRequest()
		response = postgresql.NewCreateParameterTemplateResponse()

		modifyRequest = postgresql.NewModifyParameterTemplateRequest()
	)
	if v, ok := d.GetOk("template_name"); ok {
		request.TemplateName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("db_major_version"); ok {
		request.DBMajorVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("db_engine"); ok {
		request.DBEngine = helper.String(v.(string))
	}

	if v, ok := d.GetOk("template_description"); ok {
		request.TemplateDescription = helper.String(v.(string))
	}

	result, err := meta.(*TencentCloudClient).apiV3Conn.UsePostgresqlClient().CreateParameterTemplate(request)

	if err != nil {
		log.Printf("[CRITAL]%s create postgresql ParameterTemplate failed, reason:%+v", logId, err)
		return err
	}
	response = result
	templateId := response.Response.TemplateId

	// call ModifyParameterTemplate to set param entry
	modifyRequest.TemplateId = templateId
	if v, ok := d.GetOk("modify_param_entry_set"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			paramEntry := postgresql.ParamEntry{}
			if v, ok := dMap["name"]; ok {
				paramEntry.Name = helper.String(v.(string))
			}
			if v, ok := dMap["expected_value"]; ok {
				paramEntry.ExpectedValue = helper.String(v.(string))
			}
			modifyRequest.ModifyParamEntrySet = append(modifyRequest.ModifyParamEntrySet, &paramEntry)
		}
	}

	if v, ok := d.GetOk("delete_param_set"); ok {
		deleteParamSetSet := v.(*schema.Set).List()
		for i := range deleteParamSetSet {
			deleteParamSet := deleteParamSetSet[i].(string)
			modifyRequest.DeleteParamSet = append(modifyRequest.DeleteParamSet, &deleteParamSet)
		}
	}

	if len(modifyRequest.ModifyParamEntrySet) > 0 || len(modifyRequest.DeleteParamSet) > 0 {
		_, err = meta.(*TencentCloudClient).apiV3Conn.UsePostgresqlClient().ModifyParameterTemplate(modifyRequest)
		if err != nil {
			log.Printf("[CRITAL]%s update postgresql ParameterTemplate in create method failed, reason:%+v", logId, err)
			return err
		}
	}

	d.SetId(*templateId)

	return resourceTencentCloudPostgresqlParameterTemplateRead(d, meta)
}

func resourceTencentCloudPostgresqlParameterTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_parameter_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := PostgresqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	templateId := d.Id()

	ParameterTemplate, err := service.DescribePostgresqlParameterTemplateById(ctx, templateId)
	if err != nil {
		return err
	}

	if ParameterTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `PostgresqlParameterTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if ParameterTemplate.TemplateName != nil {
		_ = d.Set("template_name", ParameterTemplate.TemplateName)
	}

	if ParameterTemplate.DBMajorVersion != nil {
		_ = d.Set("db_major_version", ParameterTemplate.DBMajorVersion)
	}

	if ParameterTemplate.DBEngine != nil {
		_ = d.Set("db_engine", ParameterTemplate.DBEngine)
	}

	if ParameterTemplate.TemplateDescription != nil {
		_ = d.Set("template_description", ParameterTemplate.TemplateDescription)
	}

	paramInfoSetList := []interface{}{}
	if ParameterTemplate.ParamInfoSet != nil {
		for _, paramInfoSet := range ParameterTemplate.ParamInfoSet {
			paramInfoSetMap := map[string]interface{}{}

			if paramInfoSet.ID != nil {
				paramInfoSetMap["id"] = paramInfoSet.ID
			}

			if paramInfoSet.Name != nil {
				paramInfoSetMap["name"] = paramInfoSet.Name
			}

			if paramInfoSet.ParamValueType != nil {
				paramInfoSetMap["param_value_type"] = paramInfoSet.ParamValueType
			}

			if paramInfoSet.Unit != nil {
				paramInfoSetMap["unit"] = paramInfoSet.Unit
			}

			if paramInfoSet.DefaultValue != nil {
				paramInfoSetMap["default_value"] = paramInfoSet.DefaultValue
			}

			if paramInfoSet.CurrentValue != nil {
				paramInfoSetMap["current_value"] = paramInfoSet.CurrentValue
			}

			if paramInfoSet.Max != nil {
				paramInfoSetMap["max"] = paramInfoSet.Max
			}

			if paramInfoSet.EnumValue != nil {
				paramInfoSetMap["enum_value"] = paramInfoSet.EnumValue
			}

			if paramInfoSet.Min != nil {
				paramInfoSetMap["min"] = paramInfoSet.Min
			}

			if paramInfoSet.ParamDescriptionCH != nil {
				paramInfoSetMap["param_description_ch"] = paramInfoSet.ParamDescriptionCH
			}

			if paramInfoSet.ParamDescriptionEN != nil {
				paramInfoSetMap["param_description_en"] = paramInfoSet.ParamDescriptionEN
			}

			if paramInfoSet.NeedReboot != nil {
				paramInfoSetMap["need_reboot"] = paramInfoSet.NeedReboot
			}

			if paramInfoSet.ClassificationCN != nil {
				paramInfoSetMap["classification_cn"] = paramInfoSet.ClassificationCN
			}

			if paramInfoSet.ClassificationEN != nil {
				paramInfoSetMap["classification_en"] = paramInfoSet.ClassificationEN
			}

			if paramInfoSet.SpecRelated != nil {
				paramInfoSetMap["spec_related"] = paramInfoSet.SpecRelated
			}

			if paramInfoSet.Advanced != nil {
				paramInfoSetMap["advanced"] = paramInfoSet.Advanced
			}

			if paramInfoSet.LastModifyTime != nil {
				paramInfoSetMap["last_modify_time"] = paramInfoSet.LastModifyTime
			}

			if paramInfoSet.StandbyRelated != nil {
				paramInfoSetMap["standby_related"] = paramInfoSet.StandbyRelated
			}

			if paramInfoSet.VersionRelationSet != nil {
				versionRelationSetList := []interface{}{}
				for _, versionRelationSet := range paramInfoSet.VersionRelationSet {
					versionRelationSetMap := map[string]interface{}{}

					if versionRelationSet.Name != nil {
						versionRelationSetMap["name"] = versionRelationSet.Name
					}

					if versionRelationSet.DBKernelVersion != nil {
						versionRelationSetMap["db_kernel_version"] = versionRelationSet.DBKernelVersion
					}

					if versionRelationSet.Value != nil {
						versionRelationSetMap["value"] = versionRelationSet.Value
					}

					if versionRelationSet.Unit != nil {
						versionRelationSetMap["unit"] = versionRelationSet.Unit
					}

					if versionRelationSet.Max != nil {
						versionRelationSetMap["max"] = versionRelationSet.Max
					}

					if versionRelationSet.Min != nil {
						versionRelationSetMap["min"] = versionRelationSet.Min
					}

					if versionRelationSet.EnumValue != nil {
						versionRelationSetMap["enum_value"] = versionRelationSet.EnumValue
					}

					versionRelationSetList = append(versionRelationSetList, versionRelationSetMap)
				}

				paramInfoSetMap["version_relation_set"] = []interface{}{versionRelationSetList}
			}

			if paramInfoSet.SpecRelationSet != nil {
				specRelationSetList := []interface{}{}
				for _, specRelationSet := range paramInfoSet.SpecRelationSet {
					specRelationSetMap := map[string]interface{}{}

					if specRelationSet.Name != nil {
						specRelationSetMap["name"] = specRelationSet.Name
					}

					if specRelationSet.Memory != nil {
						specRelationSetMap["memory"] = specRelationSet.Memory
					}

					if specRelationSet.Value != nil {
						specRelationSetMap["value"] = specRelationSet.Value
					}

					if specRelationSet.Unit != nil {
						specRelationSetMap["unit"] = specRelationSet.Unit
					}

					if specRelationSet.Max != nil {
						specRelationSetMap["max"] = specRelationSet.Max
					}

					if specRelationSet.Min != nil {
						specRelationSetMap["min"] = specRelationSet.Min
					}

					if specRelationSet.EnumValue != nil {
						specRelationSetMap["enum_value"] = specRelationSet.EnumValue
					}

					specRelationSetList = append(specRelationSetList, specRelationSetMap)
				}

				paramInfoSetMap["spec_relation_set"] = []interface{}{specRelationSetList}
			}

			paramInfoSetList = append(paramInfoSetList, paramInfoSetMap)
		}
	}
	_ = d.Set("param_info_set", paramInfoSetList)

	return nil
}

func resourceTencentCloudPostgresqlParameterTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_parameter_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := postgresql.NewModifyParameterTemplateRequest()

	request.TemplateId = helper.String(d.Id())

	immutableArgs := []string{"db_major_version", "db_engine"}

	// do not care the param_info_set attribute
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("template_name") {
		if v, ok := d.GetOk("template_name"); ok {
			request.TemplateName = helper.String(v.(string))
		}
	}

	if d.HasChange("template_description") {
		if v, ok := d.GetOk("template_description"); ok {
			request.TemplateDescription = helper.String(v.(string))
		}
	}

	if d.HasChange("modify_param_entry_set") {
		if v, ok := d.GetOk("modify_param_entry_set"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				paramEntry := postgresql.ParamEntry{}
				if v, ok := dMap["name"]; ok {
					paramEntry.Name = helper.String(v.(string))
				}
				if v, ok := dMap["expected_value"]; ok {
					paramEntry.ExpectedValue = helper.String(v.(string))
				}
				request.ModifyParamEntrySet = append(request.ModifyParamEntrySet, &paramEntry)
			}
		}
	}

	if d.HasChange("delete_param_set") {
		if v, ok := d.GetOk("delete_param_set"); ok {
			deleteParamSetSet := v.(*schema.Set).List()
			for i := range deleteParamSetSet {
				deleteParamSet := deleteParamSetSet[i].(string)
				request.DeleteParamSet = append(request.DeleteParamSet, &deleteParamSet)
			}
		}
	}

	_, err := meta.(*TencentCloudClient).apiV3Conn.UsePostgresqlClient().ModifyParameterTemplate(request)

	if err != nil {
		log.Printf("[CRITAL]%s update postgresql ParameterTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudPostgresqlParameterTemplateRead(d, meta)
}

func resourceTencentCloudPostgresqlParameterTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_parameter_template.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := PostgresqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	templateId := d.Id()

	if err := service.DeletePostgresqlParameterTemplateById(ctx, templateId); err != nil {
		return err
	}

	return nil
}
