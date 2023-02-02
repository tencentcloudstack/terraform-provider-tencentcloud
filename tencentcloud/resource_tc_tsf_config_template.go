/*
Provides a resource to create a tsf config_template

Example Usage

```hcl
resource "tencentcloud_tsf_config_template" "config_template" {
  config_template_name = ""
  config_template_type = ""
  config_template_value = ""
  config_template_desc = ""
  program_id_list =
}
```

Import

tsf config_template can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_config_template.config_template config_template_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTsfConfigTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfConfigTemplateCreate,
		Read:   resourceTencentCloudTsfConfigTemplateRead,
		Update: resourceTencentCloudTsfConfigTemplateUpdate,
		Delete: resourceTencentCloudTsfConfigTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"config_template_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Configuration template name.",
			},

			"config_template_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Configure the microservice framework corresponding to the template.",
			},

			"config_template_value": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Configure template data.",
			},

			"config_template_desc": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Configuration template description.",
			},

			"program_id_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Program id list.",
			},

			"create_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "creation time.",
			},

			"update_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "update time.",
			},

			"config_template_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Template Id.",
			},
		},
	}
}

func resourceTencentCloudTsfConfigTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_config_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request = tsf.NewCreateConfigTemplateRequest()
		// response   = tsf.NewCreateConfigTemplateResponse()
		templateId string
	)
	if v, ok := d.GetOk("config_template_name"); ok {
		request.ConfigTemplateName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("config_template_type"); ok {
		request.ConfigTemplateType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("config_template_value"); ok {
		request.ConfigTemplateValue = helper.String(v.(string))
	}

	if v, ok := d.GetOk("config_template_desc"); ok {
		request.ConfigTemplateDesc = helper.String(v.(string))
	}

	if v, ok := d.GetOk("program_id_list"); ok {
		programIdListSet := v.(*schema.Set).List()
		for i := range programIdListSet {
			programIdList := programIdListSet[i].(string)
			request.ProgramIdList = append(request.ProgramIdList, &programIdList)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().CreateConfigTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		// response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tsf configTemplate failed, reason:%+v", logId, err)
		return err
	}

	// templateId = *response.Response.templateId
	d.SetId(templateId)

	return resourceTencentCloudTsfConfigTemplateRead(d, meta)
}

func resourceTencentCloudTsfConfigTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_config_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	templateId := d.Id()

	configTemplate, err := service.DescribeTsfConfigTemplateById(ctx, templateId)
	if err != nil {
		return err
	}

	if configTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfConfigTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if configTemplate.ConfigTemplateName != nil {
		_ = d.Set("config_template_name", configTemplate.ConfigTemplateName)
	}

	if configTemplate.ConfigTemplateType != nil {
		_ = d.Set("config_template_type", configTemplate.ConfigTemplateType)
	}

	if configTemplate.ConfigTemplateValue != nil {
		_ = d.Set("config_template_value", configTemplate.ConfigTemplateValue)
	}

	if configTemplate.ConfigTemplateDesc != nil {
		_ = d.Set("config_template_desc", configTemplate.ConfigTemplateDesc)
	}

	// if configTemplate.ProgramIdList != nil {
	// 	_ = d.Set("program_id_list", configTemplate.ProgramIdList)
	// }

	if configTemplate.CreateTime != nil {
		_ = d.Set("create_time", configTemplate.CreateTime)
	}

	if configTemplate.UpdateTime != nil {
		_ = d.Set("update_time", configTemplate.UpdateTime)
	}

	if configTemplate.ConfigTemplateId != nil {
		_ = d.Set("config_template_id", configTemplate.ConfigTemplateId)
	}

	return nil
}

func resourceTencentCloudTsfConfigTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_config_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tsf.NewUpdateConfigTemplateRequest()

	templateId := d.Id()

	request.ConfigTemplateId = &templateId

	immutableArgs := []string{"config_template_name", "config_template_type", "config_template_value", "config_template_desc", "program_id_list", "create_time", "update_time", "config_template_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("config_template_name") {
		if v, ok := d.GetOk("config_template_name"); ok {
			request.ConfigTemplateName = helper.String(v.(string))
		}
	}

	if d.HasChange("config_template_type") {
		if v, ok := d.GetOk("config_template_type"); ok {
			request.ConfigTemplateType = helper.String(v.(string))
		}
	}

	if d.HasChange("config_template_value") {
		if v, ok := d.GetOk("config_template_value"); ok {
			request.ConfigTemplateValue = helper.String(v.(string))
		}
	}

	if d.HasChange("config_template_desc") {
		if v, ok := d.GetOk("config_template_desc"); ok {
			request.ConfigTemplateDesc = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().UpdateConfigTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tsf configTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTsfConfigTemplateRead(d, meta)
}

func resourceTencentCloudTsfConfigTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_config_template.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}
	templateId := d.Id()

	if err := service.DeleteTsfConfigTemplateById(ctx, templateId); err != nil {
		return err
	}

	return nil
}
