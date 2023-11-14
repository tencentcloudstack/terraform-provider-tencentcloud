/*
Provides a resource to create a ses template

Example Usage

```hcl
resource "tencentcloud_ses_template" "template" {
  template_name = "smsTemplateName"
  template_content {
		html = &lt;nil&gt;
		text = &lt;nil&gt;

  }
}
```

Import

ses template can be imported using the id, e.g.

```
terraform import tencentcloud_ses_template.template template_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ses "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses/v20201002"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudSesTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSesTemplateCreate,
		Read:   resourceTencentCloudSesTemplateRead,
		Update: resourceTencentCloudSesTemplateUpdate,
		Delete: resourceTencentCloudSesTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"template_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "SmsTemplateName, which must be required.",
			},

			"template_content": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Sms Template Content.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"html": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "HTML code after base64 encoding.",
						},
						"text": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Text content after base64 encoding.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudSesTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = ses.NewCreateEmailTemplateRequest()
		response   = ses.NewCreateEmailTemplateResponse()
		templateID int
	)
	if v, ok := d.GetOk("template_name"); ok {
		request.TemplateName = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "template_content"); ok {
		templateContent := ses.TemplateContent{}
		if v, ok := dMap["html"]; ok {
			templateContent.Html = helper.String(v.(string))
		}
		if v, ok := dMap["text"]; ok {
			templateContent.Text = helper.String(v.(string))
		}
		request.TemplateContent = &templateContent
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSesClient().CreateEmailTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ses template failed, reason:%+v", logId, err)
		return err
	}

	templateID = *response.Response.TemplateID
	d.SetId(helper.Int64ToStr(int64(templateID)))

	return resourceTencentCloudSesTemplateRead(d, meta)
}

func resourceTencentCloudSesTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SesService{client: meta.(*TencentCloudClient).apiV3Conn}

	templateId := d.Id()

	template, err := service.DescribeSesTemplateById(ctx, templateID)
	if err != nil {
		return err
	}

	if template == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SesTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if template.TemplateName != nil {
		_ = d.Set("template_name", template.TemplateName)
	}

	if template.TemplateContent != nil {
		templateContentMap := map[string]interface{}{}

		if template.TemplateContent.Html != nil {
			templateContentMap["html"] = template.TemplateContent.Html
		}

		if template.TemplateContent.Text != nil {
			templateContentMap["text"] = template.TemplateContent.Text
		}

		_ = d.Set("template_content", []interface{}{templateContentMap})
	}

	return nil
}

func resourceTencentCloudSesTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := ses.NewUpdateEmailTemplateRequest()

	templateId := d.Id()

	request.TemplateID = &templateID

	immutableArgs := []string{"template_name", "template_content"}

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

	if d.HasChange("template_content") {
		if dMap, ok := helper.InterfacesHeadMap(d, "template_content"); ok {
			templateContent := ses.TemplateContent{}
			if v, ok := dMap["html"]; ok {
				templateContent.Html = helper.String(v.(string))
			}
			if v, ok := dMap["text"]; ok {
				templateContent.Text = helper.String(v.(string))
			}
			request.TemplateContent = &templateContent
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSesClient().UpdateEmailTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update ses template failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSesTemplateRead(d, meta)
}

func resourceTencentCloudSesTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_template.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SesService{client: meta.(*TencentCloudClient).apiV3Conn}
	templateId := d.Id()

	if err := service.DeleteSesTemplateById(ctx, templateID); err != nil {
		return err
	}

	return nil
}
