/*
Provides a resource to create a css pad_rule_attachment

Example Usage

```hcl
resource "tencentcloud_css_pad_rule_attachment" "pad_rule_attachment" {
  domain_name = "177154.push.tlivecloud.com"
  template_id = 17067
  app_name    = "qqq"
  stream_name = "ppp"
}
```

Import

css pad_rule_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_css_pad_rule_attachment.pad_rule_attachment templateId#domainName
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCssPadRuleAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCssPadRuleAttachmentCreate,
		Read:   resourceTencentCloudCssPadRuleAttachmentRead,
		Delete: resourceTencentCloudCssPadRuleAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Push domain.",
			},

			"template_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Template id.",
			},

			"app_name": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Push path, must same with play path, default is live.",
			},

			"stream_name": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Stream name.",
			},
		},
	}
}

func resourceTencentCloudCssPadRuleAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_pad_rule_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = css.NewCreateLivePadRuleRequest()
		templateId int
		domainName string
	)
	if v, ok := d.GetOk("domain_name"); ok {
		domainName = v.(string)
		request.DomainName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("template_id"); ok {
		templateId = v.(int)
		request.TemplateId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("app_name"); ok {
		request.AppName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("stream_name"); ok {
		request.StreamName = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCssClient().CreateLivePadRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create css padRuleAttachment failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{helper.IntToStr(templateId), domainName}, FILED_SP))

	return resourceTencentCloudCssPadRuleAttachmentRead(d, meta)
}

func resourceTencentCloudCssPadRuleAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_pad_rule_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CssService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	templateId := idSplit[0]
	templateIdInt64, err := strconv.ParseInt(templateId, 10, 64)
	if err != nil {
		return fmt.Errorf("TemplateId format type error: %s", err.Error())
	}
	domainName := idSplit[1]

	padRuleAttachment, err := service.DescribeCssPadRuleAttachmentById(ctx, templateIdInt64, domainName)
	if err != nil {
		return err
	}

	if padRuleAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CssPadRuleAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if padRuleAttachment.DomainName != nil {
		_ = d.Set("domain_name", padRuleAttachment.DomainName)
	}

	if padRuleAttachment.TemplateId != nil {
		_ = d.Set("template_id", padRuleAttachment.TemplateId)
	}

	if padRuleAttachment.AppName != nil {
		_ = d.Set("app_name", padRuleAttachment.AppName)
	}

	if padRuleAttachment.StreamName != nil {
		_ = d.Set("stream_name", padRuleAttachment.StreamName)
	}

	return nil
}

func resourceTencentCloudCssPadRuleAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_pad_rule_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CssService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	templateId := idSplit[0]
	templateIdInt64, _ := strconv.ParseInt(templateId, 10, 64)
	domainName := idSplit[1]

	recordRule, err := service.DescribeCssPadRuleAttachmentById(ctx, templateIdInt64, domainName)
	if err != nil {
		return err
	}

	appName := ""
	if recordRule.AppName != nil {
		appName = *recordRule.AppName
	}

	streamName := ""
	if recordRule.StreamName != nil {
		streamName = *recordRule.StreamName
	}

	if err := service.DeleteCssPadRuleAttachmentById(ctx, domainName, appName, streamName, templateIdInt64); err != nil {
		return err
	}

	return nil
}
