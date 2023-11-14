/*
Provides a resource to create a sms template

Example Usage

```hcl
resource "tencentcloud_sms_template" "template" {
  template_name = "TemplateName"
  template_content = "Template test content"
  international = 0
  sms_type = 0
  remark = "sms test"
      }
```

Import

sms template can be imported using the id, e.g.

```
terraform import tencentcloud_sms_template.template template_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudSmsTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSmsTemplateCreate,
		Read:   resourceTencentCloudSmsTemplateRead,
		Update: resourceTencentCloudSmsTemplateUpdate,
		Delete: resourceTencentCloudSmsTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"template_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Message Template name, which must be unique.",
			},

			"template_content": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Message Template Content.",
			},

			"international": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Whether it is Global SMS: 0: Mainland China SMS; 1: Global SMS.",
			},

			"sms_type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "SMS type. 0: regular SMS, 1: marketing SMS.",
			},

			"remark": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Template remarks, such as reason for application and use case.",
			},

			"review_reply": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Review reply, i.e., response given by the reviewer, which is usually the reason for rejection.",
			},

			"create_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Application submission time in the format of UNIX timestamp in seconds.",
			},

			"status_code": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Signature application status. Valid values: 0: approved; 1: under review; -1: application rejected or failed.",
			},
		},
	}
}

func resourceTencentCloudSmsTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sms_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = sms.NewAddSmsTemplateRequest()
		response   = sms.NewAddSmsTemplateResponse()
		templateId int
	)
	if v, ok := d.GetOk("template_name"); ok {
		request.TemplateName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("template_content"); ok {
		request.TemplateContent = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("international"); ok {
		request.International = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("sms_type"); ok {
		request.SmsType = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSmsClient().AddSmsTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create sms template failed, reason:%+v", logId, err)
		return err
	}

	templateId = *response.Response.TemplateId
	d.SetId(helper.Int64ToStr(int64(templateId)))

	return resourceTencentCloudSmsTemplateRead(d, meta)
}

func resourceTencentCloudSmsTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sms_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SmsService{client: meta.(*TencentCloudClient).apiV3Conn}

	templateId := d.Id()

	template, err := service.DescribeSmsTemplateById(ctx, templateId)
	if err != nil {
		return err
	}

	if template == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SmsTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if template.TemplateName != nil {
		_ = d.Set("template_name", template.TemplateName)
	}

	if template.TemplateContent != nil {
		_ = d.Set("template_content", template.TemplateContent)
	}

	if template.International != nil {
		_ = d.Set("international", template.International)
	}

	if template.SmsType != nil {
		_ = d.Set("sms_type", template.SmsType)
	}

	if template.Remark != nil {
		_ = d.Set("remark", template.Remark)
	}

	if template.ReviewReply != nil {
		_ = d.Set("review_reply", template.ReviewReply)
	}

	if template.CreateTime != nil {
		_ = d.Set("create_time", template.CreateTime)
	}

	if template.StatusCode != nil {
		_ = d.Set("status_code", template.StatusCode)
	}

	return nil
}

func resourceTencentCloudSmsTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sms_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := sms.NewModifySmsTemplateRequest()

	templateId := d.Id()

	request.TemplateId = &templateId

	immutableArgs := []string{"template_name", "template_content", "international", "sms_type", "remark", "review_reply", "create_time", "status_code"}

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
		if v, ok := d.GetOk("template_content"); ok {
			request.TemplateContent = helper.String(v.(string))
		}
	}

	if d.HasChange("international") {
		if v, ok := d.GetOkExists("international"); ok {
			request.International = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("sms_type") {
		if v, ok := d.GetOkExists("sms_type"); ok {
			request.SmsType = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("remark") {
		if v, ok := d.GetOk("remark"); ok {
			request.Remark = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSmsClient().ModifySmsTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update sms template failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSmsTemplateRead(d, meta)
}

func resourceTencentCloudSmsTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sms_template.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SmsService{client: meta.(*TencentCloudClient).apiV3Conn}
	templateId := d.Id()

	if err := service.DeleteSmsTemplateById(ctx, templateId); err != nil {
		return err
	}

	return nil
}
