/*
Provides a resource to create a ses template.

Example Usage

Create a ses text template

```hcl
resource "tencentcloud_ses_template" "example" {
  template_name = "tf_example_ses_temp""
  template_content {
    text = "example for the ses template"
  }
}

```

Create a ses html template

```hcl
resource "tencentcloud_ses_template" "example" {
  template_name = "tf_example_ses_temp"
  template_content {
    html = <<-EOT
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>mail title</title>
</head>
<body>
<div class="container">
  <h1>Welcome to our service! </h1>
  <p>Dear user,</p>
  <p>Thank you for using Tencent Cloud:</p>
  <p><a href="https://cloud.tencent.com/document/product/1653">https://cloud.tencent.com/document/product/1653</a></p>
  <p>If you did not request this email, please ignore it. </p>
  <p><strong>from the iac team</strong></p>
</div>
</body>
</html>
    EOT
  }
}

```

Import

ses template can be imported using the id, e.g.
```
$ terraform import tencentcloud_ses_template.example template_id
```
*/
package tencentcloud

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ses "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses/v20201002"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudSesTemplate() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudSesTemplateRead,
		Create: resourceTencentCloudSesTemplateCreate,
		Update: resourceTencentCloudSesTemplateUpdate,
		Delete: resourceTencentCloudSesTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"template_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "smsTemplateName, which must be required.",
			},

			"template_content": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "Sms Template Content.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"html": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Html code after base64.",
						},
						"text": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Text content after base64.",
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
		response   *ses.CreateEmailTemplateResponse
		templateId uint64
	)

	if v, ok := d.GetOk("template_name"); ok {
		request.TemplateName = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "template_content"); ok {
		prometheusTemp := ses.TemplateContent{}
		if v, ok := dMap["html"]; ok {
			html := base64.StdEncoding.EncodeToString([]byte(v.(string)))
			prometheusTemp.Html = helper.String(html)
		}
		if v, ok := dMap["text"]; ok {
			text := base64.StdEncoding.EncodeToString([]byte(v.(string)))
			prometheusTemp.Text = helper.String(text)
		}
		request.TemplateContent = &prometheusTemp
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSesClient().CreateEmailTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create ses template failed, reason:%+v", logId, err)
		return err
	}

	templateId = *response.Response.TemplateID

	d.SetId(strconv.FormatUint(templateId, 10))
	return resourceTencentCloudSesTemplateRead(d, meta)
}

func resourceTencentCloudSesTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SesService{client: meta.(*TencentCloudClient).apiV3Conn}

	templateID := d.Id()
	templateId, ee := strconv.Atoi(templateID)
	if ee != nil {
		return ee
	}

	templateResponse, err := service.DescribeSesTemplate(ctx, uint64(templateId))
	if err != nil {
		return err
	}

	if templateResponse == nil {
		d.SetId("")
		return fmt.Errorf("resource `template` %v does not exist", templateId)
	}

	template := templateResponse.TemplateContent
	if template != nil {
		var templateContents []map[string]interface{}
		templateContent := map[string]interface{}{}
		if template.Html != nil {
			html, err := base64.StdEncoding.DecodeString(*template.Html)
			if err != nil {
				return err
			}
			contentHtml := string(html)
			templateContent["html"] = &contentHtml
		}
		if template.Text != nil {
			text, err := base64.StdEncoding.DecodeString(*template.Text)
			if err != nil {
				return err
			}
			contentText := string(text)
			templateContent["text"] = &contentText
		}
		templateContents = append(templateContents, templateContent)
		err = d.Set("template_content", templateContents)
		if ee != nil {
			return fmt.Errorf("set template_content err: %v", err)
		}
	}

	_ = d.Set("template_name", templateResponse.TemplateName)

	return nil
}

func resourceTencentCloudSesTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ses_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := ses.NewUpdateEmailTemplateRequest()

	templateID, ee := strconv.Atoi(d.Id())
	if ee != nil {
		return ee
	}
	templateId := uint64(templateID)
	request.TemplateID = &templateId

	if v, ok := d.GetOk("template_name"); ok {
		request.TemplateName = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "template_content"); ok {
		prometheusTemp := ses.TemplateContent{}
		if v, ok := dMap["html"]; ok {
			html := base64.StdEncoding.EncodeToString([]byte(v.(string)))
			prometheusTemp.Html = helper.String(html)
		}
		if v, ok := dMap["text"]; ok {
			text := base64.StdEncoding.EncodeToString([]byte(v.(string)))
			prometheusTemp.Text = helper.String(text)
		}
		request.TemplateContent = &prometheusTemp
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSesClient().UpdateEmailTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create ses template failed, reason:%+v", logId, err)
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

	templateID, ee := strconv.Atoi(d.Id())
	if ee != nil {
		return ee
	}
	templateId := uint64(templateID)

	if err := service.DeleteSesTemplateById(ctx, templateId); err != nil {
		return err
	}

	return nil
}
