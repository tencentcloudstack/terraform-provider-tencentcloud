/*
Provides a resource to create a cdb param_template

Example Usage

```hcl
resource "tencentcloud_cdb_param_template" "param_template" {
  name = &lt;nil&gt;
  description = &lt;nil&gt;
  engine_version = &lt;nil&gt;
  template_id = &lt;nil&gt;
  param_list {
		name = &lt;nil&gt;
		current_value = &lt;nil&gt;

  }
  template_type = "HIGH_STABILITY"
  engine_type = "InnoDB"
}
```

Import

cdb param_template can be imported using the id, e.g.

```
terraform import tencentcloud_cdb_param_template.param_template param_template_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudCdbParamTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCdbParamTemplateCreate,
		Read:   resourceTencentCloudCdbParamTemplateRead,
		Update: resourceTencentCloudCdbParamTemplateUpdate,
		Delete: resourceTencentCloudCdbParamTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The name of parameter template.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The description of parameter template.",
			},

			"engine_version": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The version of MySQL.",
			},

			"template_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The ID of source parameter template.",
			},

			"param_list": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Parameter list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of parameter.",
						},
						"current_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The value of parameter.",
						},
					},
				},
			},

			"template_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The default type of parameter template, supported value is HIGH_STABILITY or HIGH_PERFORMANCE.",
			},

			"engine_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The engine type of instance, optional value is InnoDB or RocksDB, default to InnoDB.",
			},
		},
	}
}

func resourceTencentCloudCdbParamTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_param_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = cdb.NewCreateParamTemplateRequest()
		response   = cdb.NewCreateParamTemplateResponse()
		templateId int
	)
	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("engine_version"); ok {
		request.EngineVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("template_id"); ok {
		templateId = v.(int64)
		request.TemplateId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("param_list"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			parameter := cdb.Parameter{}
			if v, ok := dMap["name"]; ok {
				parameter.Name = helper.String(v.(string))
			}
			if v, ok := dMap["current_value"]; ok {
				parameter.CurrentValue = helper.String(v.(string))
			}
			request.ParamList = append(request.ParamList, &parameter)
		}
	}

	if v, ok := d.GetOk("template_type"); ok {
		request.TemplateType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("engine_type"); ok {
		request.EngineType = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCdbClient().CreateParamTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cdb paramTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = *response.Response.TemplateId
	d.SetId(helper.Int64ToStr(templateId))

	return resourceTencentCloudCdbParamTemplateRead(d, meta)
}

func resourceTencentCloudCdbParamTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_param_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	paramTemplateId := d.Id()

	paramTemplate, err := service.DescribeCdbParamTemplateById(ctx, templateId)
	if err != nil {
		return err
	}

	if paramTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CdbParamTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if paramTemplate.Name != nil {
		_ = d.Set("name", paramTemplate.Name)
	}

	if paramTemplate.Description != nil {
		_ = d.Set("description", paramTemplate.Description)
	}

	if paramTemplate.EngineVersion != nil {
		_ = d.Set("engine_version", paramTemplate.EngineVersion)
	}

	if paramTemplate.TemplateId != nil {
		_ = d.Set("template_id", paramTemplate.TemplateId)
	}

	if paramTemplate.ParamList != nil {
		paramListList := []interface{}{}
		for _, paramList := range paramTemplate.ParamList {
			paramListMap := map[string]interface{}{}

			if paramTemplate.ParamList.Name != nil {
				paramListMap["name"] = paramTemplate.ParamList.Name
			}

			if paramTemplate.ParamList.CurrentValue != nil {
				paramListMap["current_value"] = paramTemplate.ParamList.CurrentValue
			}

			paramListList = append(paramListList, paramListMap)
		}

		_ = d.Set("param_list", paramListList)

	}

	if paramTemplate.TemplateType != nil {
		_ = d.Set("template_type", paramTemplate.TemplateType)
	}

	if paramTemplate.EngineType != nil {
		_ = d.Set("engine_type", paramTemplate.EngineType)
	}

	return nil
}

func resourceTencentCloudCdbParamTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_param_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cdb.NewModifyParamTemplateRequest()

	paramTemplateId := d.Id()

	request.TemplateId = &templateId

	immutableArgs := []string{"name", "description", "engine_version", "template_id", "param_list", "template_type", "engine_type"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
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
				parameter := cdb.Parameter{}
				if v, ok := dMap["name"]; ok {
					parameter.Name = helper.String(v.(string))
				}
				if v, ok := dMap["current_value"]; ok {
					parameter.CurrentValue = helper.String(v.(string))
				}
				request.ParamList = append(request.ParamList, &parameter)
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCdbClient().ModifyParamTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cdb paramTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCdbParamTemplateRead(d, meta)
}

func resourceTencentCloudCdbParamTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_param_template.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	paramTemplateId := d.Id()

	if err := service.DeleteCdbParamTemplateById(ctx, templateId); err != nil {
		return err
	}

	return nil
}
