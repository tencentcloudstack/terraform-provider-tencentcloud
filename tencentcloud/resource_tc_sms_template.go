/*
Provides a resource to create a sms template

Example Usage

```hcl
resource "tencentcloud_sms_template" "template" {
  template_name = "TemplateName"
  template_content = "Template test content"
  international = 0
  sms_type = 0
  remark = "短信tf测试"
}

```
Import

sms template can be imported using the id, e.g.
```
$ terraform import tencentcloud_sms_template.template template_id
```
*/

package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudSmsTemplate() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudSmsTemplateRead,
		Create: resourceTencentCloudSmsTemplateCreate,

		Update: resourceTencentCloudSmsTemplateUpdate,
		Delete: resourceTencentCloudSmsTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"template_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Message Template name, which must be unique.",
			},

			"template_content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Message Template Content.",
			},

			"international": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Whether it is Global SMS: 0: Mainland China SMS; 1: Global SMS.",
			},

			"sms_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "SMS type. 0: regular SMS, 1: marketing SMS.",
			},

			"remark": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Template remarks, such as reason for application and use case.",
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
		response   *sms.AddSmsTemplateResponse
		//未使用
		// templateId uint64
	)

	if v, ok := d.GetOk("template_name"); ok {
		request.TemplateName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("template_content"); ok {
		request.TemplateContent = helper.String(v.(string))
	}

	// if v, ok := d.GetOk("international"); ok {
	// 	request.International = helper.IntUint64(v.(int))
	// }
	
	// 更换函数 d.GetOk()  -> d.Get()  在uint64，值为0时会 ok会为false 导致参数没传进入
	international := d.Get("international")
	request.International = helper.IntUint64(international.(int))
	
	sms_type := d.Get("sms_type")
	request.SmsType = helper.IntUint64(sms_type.(int))

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSmsClient().AddSmsTemplate(request)
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
		log.Printf("[CRITAL]%s create sms template failed, reason:%+v", logId, err)
		return err
	}

	templateId_string := *response.Response.AddTemplateStatus.TemplateId //数据结构修改
	// templateId, _ = strconv.ParseUint(templateId_string, 10, 64) //类型修改，未使用

	d.SetId(templateId_string)  //函数只接受字符串id
	d.Set("international", international) //
	return resourceTencentCloudSmsTemplateRead(d, meta)
}

func resourceTencentCloudSmsTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sms_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SmsService{client: meta.(*TencentCloudClient).apiV3Conn}

	//read api 传参
	templateId := d.Id()
	international := d.Get("international") 
	template, err := service.DescribeSmsTemplate(ctx, templateId, helper.IntUint64(international.(int)))

	if err != nil {
		return err
	}

	if template == nil {
		d.SetId("")
		return fmt.Errorf("resource `template` %s does not exist", templateId)
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

	if template.ReviewReply != nil {  //修改字段 SmsType -> ReviewReply
		_ = d.Set("review_reply", template.ReviewReply)
	}

	if template.CreateTime != nil {   //修改字段 Remark -> CreateTime
		_ = d.Set("create_time", template.CreateTime)
	}
	//添加字段 StatusCode
	if template.StatusCode != nil {   
		_ = d.Set("status_code", template.StatusCode)
	}

	return nil
}

func resourceTencentCloudSmsTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sms_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	// ctx定义但未使用
	// ctx := context.WithValue(context.TODO(), logIdKey, logId) 

	request := sms.NewModifySmsTemplateRequest()

	templateId_string := d.Id()
	templateId, _ := strconv.ParseUint(templateId_string, 10, 64)

	request.TemplateId = &templateId

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
		if v, ok := d.GetOk("international"); ok {
			request.International = helper.IntUint64(v.(int))
		}

	}

	if d.HasChange("sms_type") {
		if v, ok := d.GetOk("sms_type"); ok {
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
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create sms template failed, reason:%+v", logId, err)
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
