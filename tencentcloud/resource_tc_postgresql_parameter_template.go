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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Type:        schema.TypeSet,
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
		modifyParamSet := v.(*schema.Set).List()
		for _, item := range modifyParamSet {
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
		deleteParamSet := v.(*schema.Set).List()
		for i := range deleteParamSet {
			deleteParam := deleteParamSet[i].(string)
			modifyRequest.DeleteParamSet = append(modifyRequest.DeleteParamSet, &deleteParam)
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

	// outer layer declaration to avoid the API returning null as ParamInfo
	paramInfoSetList := []interface{}{}
	if ParameterTemplate.ParamInfoSet != nil {
		for _, paramInfoSet := range ParameterTemplate.ParamInfoSet {
			paramInfoSetList = append(paramInfoSetList, map[string]interface{}{
				"name":           *paramInfoSet.Name,
				"expected_value": *paramInfoSet.CurrentValue,
			})
		}
	}
	_ = d.Set("modify_param_entry_set", paramInfoSetList)

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
			modifyParamSet := v.(*schema.Set).List()
			for _, item := range modifyParamSet {
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
