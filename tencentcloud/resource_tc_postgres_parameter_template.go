/*
Provides a resource to create a postgres parameter_template

Example Usage

```hcl
resource "tencentcloud_postgres_parameter_template" "parameter_template" {
  template_name = "test_param_template"
  d_b_major_version = "13"
  d_b_engine = "postgresql"
  template_description = "test use"
}
```

Import

postgres parameter_template can be imported using the id, e.g.

```
terraform import tencentcloud_postgres_parameter_template.parameter_template parameter_template_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgres "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudPostgresParameterTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresParameterTemplateCreate,
		Read:   resourceTencentCloudPostgresParameterTemplateRead,
		Update: resourceTencentCloudPostgresParameterTemplateUpdate,
		Delete: resourceTencentCloudPostgresParameterTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"template_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Template name, which can contain 1-60 letters, digits, and symbols (-_./()+=:@).",
			},

			"d_b_major_version": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The major database version number, such as 11, 12, 13.",
			},

			"d_b_engine": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Database engine, such as postgresql, mssql_compatible.",
			},

			"template_description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Parameter template description, which can contain 1-60 letters, digits, and symbols (-_./()+=:@).",
			},
		},
	}
}

func resourceTencentCloudPostgresParameterTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_parameter_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = postgres.NewCreateParameterTemplateRequest()
		response   = postgres.NewCreateParameterTemplateResponse()
		templateId string
	)
	if v, ok := d.GetOk("template_name"); ok {
		request.TemplateName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("d_b_major_version"); ok {
		request.DBMajorVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("d_b_engine"); ok {
		request.DBEngine = helper.String(v.(string))
	}

	if v, ok := d.GetOk("template_description"); ok {
		request.TemplateDescription = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresClient().CreateParameterTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create postgres ParameterTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = *response.Response.TemplateId
	d.SetId(templateId)

	return resourceTencentCloudPostgresParameterTemplateRead(d, meta)
}

func resourceTencentCloudPostgresParameterTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_parameter_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := PostgresService{client: meta.(*TencentCloudClient).apiV3Conn}

	parameterTemplateId := d.Id()

	ParameterTemplate, err := service.DescribePostgresParameterTemplateById(ctx, templateId)
	if err != nil {
		return err
	}

	if ParameterTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `PostgresParameterTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if ParameterTemplate.TemplateName != nil {
		_ = d.Set("template_name", ParameterTemplate.TemplateName)
	}

	if ParameterTemplate.DBMajorVersion != nil {
		_ = d.Set("d_b_major_version", ParameterTemplate.DBMajorVersion)
	}

	if ParameterTemplate.DBEngine != nil {
		_ = d.Set("d_b_engine", ParameterTemplate.DBEngine)
	}

	if ParameterTemplate.TemplateDescription != nil {
		_ = d.Set("template_description", ParameterTemplate.TemplateDescription)
	}

	return nil
}

func resourceTencentCloudPostgresParameterTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_parameter_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := postgres.NewModifyParameterTemplateRequest()

	parameterTemplateId := d.Id()

	request.TemplateId = &templateId

	immutableArgs := []string{"template_name", "d_b_major_version", "d_b_engine", "template_description"}

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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresClient().ModifyParameterTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update postgres ParameterTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudPostgresParameterTemplateRead(d, meta)
}

func resourceTencentCloudPostgresParameterTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_parameter_template.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := PostgresService{client: meta.(*TencentCloudClient).apiV3Conn}
	parameterTemplateId := d.Id()

	if err := service.DeletePostgresParameterTemplateById(ctx, templateId); err != nil {
		return err
	}

	return nil
}
