/*
Provides a resource to create a live callback_template

Example Usage

```hcl
resource "tencentcloud_live_callback_template" "callback_template" {
  template_name = "demo"
  description = "this is demo"
  stream_begin_notify_url = "http://www.yourdomain.com/api/notify?action=streamBegin"
  stream_end_notify_url = "http://www.yourdomain.com/api/notify?action=streamEnd"
  record_notify_url = "http://www.yourdomain.com/api/notify?action=record"
  snapshot_notify_url = "http://www.yourdomain.com/api/notify?action=snapshot"
  porn_censorship_notify_url = "http://www.yourdomain.com/api/notify?action=porn"
  callback_key = "adasda131312"
  push_exception_notify_url = "http://www.yourdomain.com/api/notify?action=pushException"
}
```

Import

live callback_template can be imported using the id, e.g.

```
terraform import tencentcloud_live_callback_template.callback_template callback_template_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	live "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudLiveCallbackTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLiveCallbackTemplateCreate,
		Read:   resourceTencentCloudLiveCallbackTemplateRead,
		Update: resourceTencentCloudLiveCallbackTemplateUpdate,
		Delete: resourceTencentCloudLiveCallbackTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"template_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Template name.Maximum length: 255 bytes. Only Chinese, English, numbers, &amp;amp;#39;_&amp;amp;#39;, &amp;amp;#39;-&amp;amp;#39; are supported.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Description information.Maximum length: 1024 bytes.Only Chinese, English, numbers, &amp;amp;#39;_&amp;amp;#39;, &amp;amp;#39;-&amp;amp;#39; are supported.",
			},

			"stream_begin_notify_url": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Launch callback URL.",
			},

			"stream_end_notify_url": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Cutoff callback URL.",
			},

			"record_notify_url": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Recording callback URL.",
			},

			"snapshot_notify_url": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Snapshot callback URL.",
			},

			"porn_censorship_notify_url": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "PornCensorship callback URL.",
			},

			"callback_key": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Callback Key, public callback URL.",
			},

			"push_exception_notify_url": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: " Streaming Exception Callback URL .",
			},
		},
	}
}

func resourceTencentCloudLiveCallbackTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_callback_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = live.NewCreateLiveCallbackTemplateRequest()
		response   = live.NewCreateLiveCallbackTemplateResponse()
		templateId int
	)
	if v, ok := d.GetOk("template_name"); ok {
		request.TemplateName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("stream_begin_notify_url"); ok {
		request.StreamBeginNotifyUrl = helper.String(v.(string))
	}

	if v, ok := d.GetOk("stream_end_notify_url"); ok {
		request.StreamEndNotifyUrl = helper.String(v.(string))
	}

	if v, ok := d.GetOk("record_notify_url"); ok {
		request.RecordNotifyUrl = helper.String(v.(string))
	}

	if v, ok := d.GetOk("snapshot_notify_url"); ok {
		request.SnapshotNotifyUrl = helper.String(v.(string))
	}

	if v, ok := d.GetOk("porn_censorship_notify_url"); ok {
		request.PornCensorshipNotifyUrl = helper.String(v.(string))
	}

	if v, ok := d.GetOk("callback_key"); ok {
		request.CallbackKey = helper.String(v.(string))
	}

	if v, ok := d.GetOk("push_exception_notify_url"); ok {
		request.PushExceptionNotifyUrl = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLiveClient().CreateLiveCallbackTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create live callbackTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = *response.Response.TemplateId
	d.SetId(helper.Int64ToStr(templateId))

	return resourceTencentCloudLiveCallbackTemplateRead(d, meta)
}

func resourceTencentCloudLiveCallbackTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_callback_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LiveService{client: meta.(*TencentCloudClient).apiV3Conn}

	callbackTemplateId := d.Id()

	callbackTemplate, err := service.DescribeLiveCallbackTemplateById(ctx, templateId)
	if err != nil {
		return err
	}

	if callbackTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `LiveCallbackTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if callbackTemplate.TemplateName != nil {
		_ = d.Set("template_name", callbackTemplate.TemplateName)
	}

	if callbackTemplate.Description != nil {
		_ = d.Set("description", callbackTemplate.Description)
	}

	if callbackTemplate.StreamBeginNotifyUrl != nil {
		_ = d.Set("stream_begin_notify_url", callbackTemplate.StreamBeginNotifyUrl)
	}

	if callbackTemplate.StreamEndNotifyUrl != nil {
		_ = d.Set("stream_end_notify_url", callbackTemplate.StreamEndNotifyUrl)
	}

	if callbackTemplate.RecordNotifyUrl != nil {
		_ = d.Set("record_notify_url", callbackTemplate.RecordNotifyUrl)
	}

	if callbackTemplate.SnapshotNotifyUrl != nil {
		_ = d.Set("snapshot_notify_url", callbackTemplate.SnapshotNotifyUrl)
	}

	if callbackTemplate.PornCensorshipNotifyUrl != nil {
		_ = d.Set("porn_censorship_notify_url", callbackTemplate.PornCensorshipNotifyUrl)
	}

	if callbackTemplate.CallbackKey != nil {
		_ = d.Set("callback_key", callbackTemplate.CallbackKey)
	}

	if callbackTemplate.PushExceptionNotifyUrl != nil {
		_ = d.Set("push_exception_notify_url", callbackTemplate.PushExceptionNotifyUrl)
	}

	return nil
}

func resourceTencentCloudLiveCallbackTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_callback_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := live.NewModifyLiveCallbackTemplateRequest()

	callbackTemplateId := d.Id()

	request.TemplateId = &templateId

	immutableArgs := []string{"template_name", "description", "stream_begin_notify_url", "stream_end_notify_url", "record_notify_url", "snapshot_notify_url", "porn_censorship_notify_url", "callback_key", "push_exception_notify_url"}

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

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	if d.HasChange("stream_begin_notify_url") {
		if v, ok := d.GetOk("stream_begin_notify_url"); ok {
			request.StreamBeginNotifyUrl = helper.String(v.(string))
		}
	}

	if d.HasChange("stream_end_notify_url") {
		if v, ok := d.GetOk("stream_end_notify_url"); ok {
			request.StreamEndNotifyUrl = helper.String(v.(string))
		}
	}

	if d.HasChange("record_notify_url") {
		if v, ok := d.GetOk("record_notify_url"); ok {
			request.RecordNotifyUrl = helper.String(v.(string))
		}
	}

	if d.HasChange("snapshot_notify_url") {
		if v, ok := d.GetOk("snapshot_notify_url"); ok {
			request.SnapshotNotifyUrl = helper.String(v.(string))
		}
	}

	if d.HasChange("porn_censorship_notify_url") {
		if v, ok := d.GetOk("porn_censorship_notify_url"); ok {
			request.PornCensorshipNotifyUrl = helper.String(v.(string))
		}
	}

	if d.HasChange("callback_key") {
		if v, ok := d.GetOk("callback_key"); ok {
			request.CallbackKey = helper.String(v.(string))
		}
	}

	if d.HasChange("push_exception_notify_url") {
		if v, ok := d.GetOk("push_exception_notify_url"); ok {
			request.PushExceptionNotifyUrl = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLiveClient().ModifyLiveCallbackTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update live callbackTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudLiveCallbackTemplateRead(d, meta)
}

func resourceTencentCloudLiveCallbackTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_callback_template.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LiveService{client: meta.(*TencentCloudClient).apiV3Conn}
	callbackTemplateId := d.Id()

	if err := service.DeleteLiveCallbackTemplateById(ctx, templateId); err != nil {
		return err
	}

	return nil
}
