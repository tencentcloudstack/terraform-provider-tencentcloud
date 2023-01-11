/*
Provides a resource to create a mysql param template

Example Usage

```hcl
resource "tencentcloud_mysql_param_template" "param_template" {
  name           = "terraform-test"
  description    = "terraform-test"
  engine_version = "8.0"
  template_type  = "HIGH_STABILITY"
  engine_type    = "InnoDB"
}
```

Import

mysql param template can be imported using the id, e.g.

```
terraform import tencentcloud_mysql_param_template.param_template template_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMysqlParamTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlParamTemplateCreate,
		Read:   resourceTencentCloudMysqlParamTemplateRead,
		Update: resourceTencentCloudMysqlParamTemplateUpdate,
		Delete: resourceTencentCloudMysqlParamTemplateDelete,
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

			"template_id": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The ID of source parameter template.",
			},

			"param_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "parameter list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of parameter.",
						},
						"current_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The value of parameter.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudMysqlParamTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_param_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = mysql.NewCreateParamTemplateRequest()
		response   = mysql.NewCreateParamTemplateResponse()
		templateId int64
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

	if v, ok := d.GetOk("template_type"); ok {
		request.TemplateType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("engine_type"); ok {
		request.EngineType = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMysqlClient().CreateParamTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mysql paramTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = *response.Response.TemplateId
	d.SetId(helper.Int64ToStr(templateId))

	return resourceTencentCloudMysqlParamTemplateRead(d, meta)
}

func resourceTencentCloudMysqlParamTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_param_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	templateId := d.Id()

	paramTemplate, err := service.DescribeMysqlParamTemplateById(ctx, templateId)
	if err != nil {
		return err
	}

	if paramTemplate == nil {
		d.SetId("")
		return fmt.Errorf("resource `MysqlParamTemplate` %s does not exist", d.Id())
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

	if paramTemplate.Items != nil {
		paramListList := []interface{}{}
		for _, paramList := range paramTemplate.Items {
			paramListMap := map[string]interface{}{}

			if paramList.Name != nil {
				paramListMap["name"] = paramList.Name
			}

			if paramList.CurrentValue != nil {
				paramListMap["current_value"] = paramList.CurrentValue
			}

			paramListList = append(paramListList, paramListMap)
		}

		_ = d.Set("param_list", paramListList)

	}

	if paramTemplate.TemplateType != nil {
		_ = d.Set("template_type", paramTemplate.TemplateType)
	}

	return nil
}

func resourceTencentCloudMysqlParamTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_param_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := mysql.NewModifyParamTemplateRequest()

	templateId := d.Id()

	request.TemplateId = helper.StrToInt64Point(templateId)

	immutableArgs := []string{"engine_version", "template_type", "engine_type"}

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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMysqlClient().ModifyParamTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update mysql paramTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMysqlParamTemplateRead(d, meta)
}

func resourceTencentCloudMysqlParamTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_param_template.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	templateId := d.Id()

	if err := service.DeleteMysqlParamTemplateById(ctx, templateId); err != nil {
		return err
	}

	return nil
}
