/*
Provides a resource to create a css watermark_rule

Example Usage

Binding watermark rule with a css stream

```hcl
resource "tencentcloud_css_pull_stream_task" "example" {
  stream_name = "tf_example_stream_name"
  source_type = "PullLivePushLive"
  source_urls = ["rtmp://xxx.com/live/stream"]
  domain_name = "test.domain.com"
  app_name    = "live"
  start_time  = "2023-09-27T10:28:21Z"
  end_time    = "2023-09-27T17:28:21Z"
  operator    = "tf_admin"
  comment     = "This is a e2e test case."
}

resource "tencentcloud_css_watermark" "example" {
  picture_url    = "picture_url"
  watermark_name = "watermark_name"
  x_position     = 0
  y_position     = 0
  width          = 0
  height         = 0
}

resource "tencentcloud_css_watermark_rule_attachment" "watermark_rule_attachment" {
  domain_name = tencentcloud_css_pull_stream_task.example.domain_name
  app_name    = tencentcloud_css_pull_stream_task.example.app_name
  stream_name = tencentcloud_css_pull_stream_task.example.stream_name
  template_id = tencentcloud_css_watermark.example.id
}
```

Import

css watermark_rule_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_css_watermark_rule_attachment.watermark_rule domain_name#app_name#stream_name#template_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCssWatermarkRuleAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCssWatermarkRuleAttachmentCreate,
		Read:   resourceTencentCloudCssWatermarkRuleAttachmentRead,
		Delete: resourceTencentCloudCssWatermarkRuleAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "rule domain name.",
			},

			"app_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "rule app name.",
			},

			"stream_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "rule stream name.",
			},

			"template_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "The template Id can be acquired by the Id of `tencentcloud_css_watermark`.",
			},

			"create_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "create time.",
			},

			"update_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "update time.",
			},
		},
	}
}

func resourceTencentCloudCssWatermarkRuleAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_watermark_rule_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = css.NewCreateLiveWatermarkRuleRequest()
		domainName string
		appName    string
		streamName string
		templateId int
	)
	if v, ok := d.GetOk("domain_name"); ok {
		domainName = v.(string)
		request.DomainName = helper.String(domainName)
	}

	if v, ok := d.GetOk("app_name"); ok {
		appName = v.(string)
		request.AppName = helper.String(appName)
	}

	if v, ok := d.GetOk("stream_name"); ok {
		streamName = v.(string)
		request.StreamName = helper.String(streamName)
	}

	if v, ok := d.GetOkExists("template_id"); ok {
		templateId = v.(int)
		request.TemplateId = helper.IntInt64(templateId)
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCssClient().CreateLiveWatermarkRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create css watermarkRule failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{domainName, appName, streamName, helper.IntToStr(templateId)}, FILED_SP))

	return resourceTencentCloudCssWatermarkRuleAttachmentRead(d, meta)
}

func resourceTencentCloudCssWatermarkRuleAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_watermark_rule_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CssService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	domainName := idSplit[0]
	appName := idSplit[1]
	streamName := idSplit[2]
	templateId := idSplit[3]

	watermarkRule, err := service.DescribeCssWatermarkRuleAttachment(ctx, domainName, appName, streamName, templateId)
	if err != nil {
		return err
	}

	if watermarkRule == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CssWatermarkRule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if watermarkRule.DomainName != nil {
		_ = d.Set("domain_name", watermarkRule.DomainName)
	}

	if watermarkRule.AppName != nil {
		_ = d.Set("app_name", watermarkRule.AppName)
	}

	if watermarkRule.StreamName != nil {
		_ = d.Set("stream_name", watermarkRule.StreamName)
	}

	if watermarkRule.TemplateId != nil {
		_ = d.Set("template_id", watermarkRule.TemplateId)
	}

	if watermarkRule.CreateTime != nil {
		_ = d.Set("create_time", watermarkRule.CreateTime)
	}

	if watermarkRule.UpdateTime != nil {
		_ = d.Set("update_time", watermarkRule.UpdateTime)
	}

	return nil
}

func resourceTencentCloudCssWatermarkRuleAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_watermark_rule_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CssService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	domainName := idSplit[0]
	appName := idSplit[1]
	streamName := idSplit[2]

	if err := service.DetachCssWatermarkRuleAttachment(ctx, domainName, appName, streamName); err != nil {
		return err
	}

	return nil
}
