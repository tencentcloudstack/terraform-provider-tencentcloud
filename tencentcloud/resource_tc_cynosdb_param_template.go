/*
Provides a resource to create a cynosdb param_template

Example Usage

```hcl
resource "tencentcloud_cynosdb_param_template" "param_template" {
  template_name = ""
  engine_version = "5.7"
  template_description = ""
  template_id = 1000
  db_mode = "NORMAL"
  param_list {
		param_name = ""
		current_value = ""
		old_value = ""

  }
}
```

Import

cynosdb param_template can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_param_template.param_template param_template_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudCynosdbParamTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbParamTemplateCreate,
		Read:   resourceTencentCloudCynosdbParamTemplateRead,
		Update: resourceTencentCloudCynosdbParamTemplateUpdate,
		Delete: resourceTencentCloudCynosdbParamTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"template_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Template Name.",
			},

			"engine_version": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "MySQL version number.",
			},

			"template_description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Template Description.",
			},

			"template_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Optional parameter, template ID to be copied.",
			},

			"db_mode": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Database type, optional values: NORMAL (default), SERVERLESS.",
			},

			"param_list": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Parameter list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"param_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Parameter Name.",
						},
						"current_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Current value.",
						},
						"old_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Original value.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCynosdbParamTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_param_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = cynosdb.NewCreateParamTemplateRequest()
		response   = cynosdb.NewCreateParamTemplateResponse()
		templateId int
	)
	if v, ok := d.GetOk("template_name"); ok {
		request.TemplateName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("engine_version"); ok {
		request.EngineVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("template_description"); ok {
		request.TemplateDescription = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("template_id"); ok {
		templateId = v.(int64)
		request.TemplateId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("db_mode"); ok {
		request.DbMode = helper.String(v.(string))
	}

	if v, ok := d.GetOk("param_list"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			paramItem := cynosdb.ParamItem{}
			if v, ok := dMap["param_name"]; ok {
				paramItem.ParamName = helper.String(v.(string))
			}
			if v, ok := dMap["current_value"]; ok {
				paramItem.CurrentValue = helper.String(v.(string))
			}
			if v, ok := dMap["old_value"]; ok {
				paramItem.OldValue = helper.String(v.(string))
			}
			request.ParamList = append(request.ParamList, &paramItem)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().CreateParamTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cynosdb paramTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = *response.Response.TemplateId
	d.SetId(helper.Int64ToStr(templateId))

	return resourceTencentCloudCynosdbParamTemplateRead(d, meta)
}

func resourceTencentCloudCynosdbParamTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_param_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	paramTemplateId := d.Id()

	paramTemplate, err := service.DescribeCynosdbParamTemplateById(ctx, templateId)
	if err != nil {
		return err
	}

	if paramTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CynosdbParamTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if paramTemplate.TemplateName != nil {
		_ = d.Set("template_name", paramTemplate.TemplateName)
	}

	if paramTemplate.EngineVersion != nil {
		_ = d.Set("engine_version", paramTemplate.EngineVersion)
	}

	if paramTemplate.TemplateDescription != nil {
		_ = d.Set("template_description", paramTemplate.TemplateDescription)
	}

	if paramTemplate.TemplateId != nil {
		_ = d.Set("template_id", paramTemplate.TemplateId)
	}

	if paramTemplate.DbMode != nil {
		_ = d.Set("db_mode", paramTemplate.DbMode)
	}

	if paramTemplate.ParamList != nil {
		paramListList := []interface{}{}
		for _, paramList := range paramTemplate.ParamList {
			paramListMap := map[string]interface{}{}

			if paramTemplate.ParamList.ParamName != nil {
				paramListMap["param_name"] = paramTemplate.ParamList.ParamName
			}

			if paramTemplate.ParamList.CurrentValue != nil {
				paramListMap["current_value"] = paramTemplate.ParamList.CurrentValue
			}

			if paramTemplate.ParamList.OldValue != nil {
				paramListMap["old_value"] = paramTemplate.ParamList.OldValue
			}

			paramListList = append(paramListList, paramListMap)
		}

		_ = d.Set("param_list", paramListList)

	}

	return nil
}

func resourceTencentCloudCynosdbParamTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_param_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cynosdb.NewModifyParamTemplateRequest()

	paramTemplateId := d.Id()

	request.TemplateId = &templateId

	immutableArgs := []string{"template_name", "engine_version", "template_description", "template_id", "db_mode", "param_list"}

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

	if d.HasChange("template_id") {
		if v, ok := d.GetOkExists("template_id"); ok {
			request.TemplateId = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("param_list") {
		if v, ok := d.GetOk("param_list"); ok {
			for _, item := range v.([]interface{}) {
				paramItem := cynosdb.ParamItem{}
				if v, ok := dMap["param_name"]; ok {
					paramItem.ParamName = helper.String(v.(string))
				}
				if v, ok := dMap["current_value"]; ok {
					paramItem.CurrentValue = helper.String(v.(string))
				}
				if v, ok := dMap["old_value"]; ok {
					paramItem.OldValue = helper.String(v.(string))
				}
				request.ParamList = append(request.ParamList, &paramItem)
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().ModifyParamTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cynosdb paramTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCynosdbParamTemplateRead(d, meta)
}

func resourceTencentCloudCynosdbParamTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_param_template.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	paramTemplateId := d.Id()

	if err := service.DeleteCynosdbParamTemplateById(ctx, templateId); err != nil {
		return err
	}

	return nil
}
