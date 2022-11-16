/*
Provides a resource to create a css watermark_rule_attachment

Example Usage

```hcl
resource "tencentcloud_css_watermark_rule_attachment" "watermark_rule_attachment" {
  domain_name = ""
  app_name = ""
  stream_name = ""
  watermark_id = ""
    }

```
Import

css watermark_rule_attachment can be imported using the id, e.g.
```
$ terraform import tencentcloud_css_watermark_rule_attachment.watermark_rule_attachment watermarkRuleAttachment_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCssWatermarkRuleAttachment() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudCssWatermarkRuleAttachmentRead,
		Create: resourceTencentCloudCssWatermarkRuleAttachmentCreate,
		Delete: resourceTencentCloudCssWatermarkRuleAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "rule domain name.",
			},

			"app_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "rule app name.",
			},

			"stream_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "rule stream name.",
			},

			"watermark_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "watermark id created by AddLiveWatermark.",
			},

			"create_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "create time.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
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
		request = css.NewCreateLiveWatermarkRuleRequest()
		// response    *css.CreateLiveWatermarkRuleResponse
		domainName  string
		appName     string
		streamName  string
		watermarkId string
	)

	if v, ok := d.GetOk("domain_name"); ok {
		domainName = v.(string)
		request.DomainName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("app_name"); ok {
		appName = v.(string)
		request.AppName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("stream_name"); ok {
		streamName = v.(string)
		request.StreamName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("watermark_id"); ok {
		watermarkId = v.(string)
		request.TemplateId = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCssClient().CreateLiveWatermarkRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create css watermarkRuleAttachment failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(domainName + FILED_SP + appName + FILED_SP + streamName + FILED_SP + watermarkId)
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
	watermarkId := idSplit[3]

	watermarkRuleAttachment, err := service.DescribeCssWatermarkRuleAttachment(ctx, domainName, appName, streamName, watermarkId)

	if err != nil {
		return err
	}

	if watermarkRuleAttachment == nil {
		d.SetId("")
		return fmt.Errorf("resource `watermarkRuleAttachment` %s does not exist", d.Id())
	}

	if watermarkRuleAttachment.DomainName != nil {
		_ = d.Set("domain_name", watermarkRuleAttachment.DomainName)
	}

	if watermarkRuleAttachment.AppName != nil {
		_ = d.Set("app_name", watermarkRuleAttachment.AppName)
	}

	if watermarkRuleAttachment.StreamName != nil {
		_ = d.Set("stream_name", watermarkRuleAttachment.StreamName)
	}

	if watermarkRuleAttachment.TemplateId != nil {
		_ = d.Set("watermark_id", watermarkRuleAttachment.TemplateId)
	}

	if watermarkRuleAttachment.CreateTime != nil {
		_ = d.Set("create_time", watermarkRuleAttachment.CreateTime)
	}

	if watermarkRuleAttachment.UpdateTime != nil {
		_ = d.Set("update_time", watermarkRuleAttachment.UpdateTime)
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
